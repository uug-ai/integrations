package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/slack-go/slack"
	"github.com/uug-ai/integrations/pkg/message"
)

func TestSlackSendWebhookCompatibility(t *testing.T) {
	for _, withMedia := range []bool{false, true} {
		t.Run(strconv.FormatBool(withMedia), func(t *testing.T) {
			msg := message.Message{Body: "Motion detected"}
			wantText := msg.Body
			if withMedia {
				msg.Media = []message.Media{{
					Type: "video", Url: "https://media.example/recording.mp4",
					ThumbnailUrl: "https://media.example/thumbnail.jpg",
				}}
				wantText += "\r\n" + msg.Media[0].Url
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var payload slack.WebhookMessage
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Error(err)
				}
				if r.Method != http.MethodPost || payload.Username != "Alerts" || payload.Text != wantText {
					t.Errorf("unexpected webhook request: %s %+v", r.Method, payload)
				}
				if withMedia {
					if len(payload.Attachments) != 1 || payload.Attachments[0].ImageURL != msg.Media[0].ThumbnailUrl {
						t.Errorf("thumbnail attachment = %+v", payload.Attachments)
					}
				} else if len(payload.Attachments) != 0 {
					t.Errorf("unexpected attachments: %+v", payload.Attachments)
				}
				_, _ = io.WriteString(w, "ok")
			}))
			defer server.Close()
			if err := (Slack{Hook: server.URL, Username: "Alerts"}).Send(msg); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSlackUploadVideo(t *testing.T) {
	for _, test := range []struct {
		name    string
		chunked bool
		title   string
	}{
		{name: "known length", title: "Camera alert"},
		{name: "chunked download", chunked: true, title: "Camera alert"},
		{name: "default title"},
	} {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			t.Setenv("TMPDIR", tempDir)
			video := "\x00\x00\x00\x18ftypmp42-recording-bytes"
			wantTitle := test.title
			if wantTitle == "" {
				wantTitle = "Recording"
			}
			requests := make(chan string, 8)
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r.URL.Path
				if r.URL.Path != "/recording.mp4" && r.Method != http.MethodPost {
					t.Errorf("Slack API method = %s", r.Method)
				}
				switch r.URL.Path {
				case "/recording.mp4":
					if r.Header.Get("Authorization") != "" {
						t.Error("Slack token leaked to media download")
					}
					if test.chunked {
						w.(http.Flusher).Flush()
					}
					_, _ = io.WriteString(w, video)
				case "/files.getUploadURLExternal":
					if err := r.ParseForm(); err != nil {
						t.Error(err)
					}
					if r.Form.Get("token") != "test-token" || r.Form.Get("filename") != "recording.mp4" || r.Form.Get("length") != strconv.Itoa(len(video)) {
						t.Errorf("upload parameters = %v", r.Form)
					}
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"ok": true, "file_id": "F123", "upload_url": server.URL + "/upload",
					})
				case "/upload":
					if r.Header.Get("Authorization") != "Bearer test-token" {
						t.Error("missing upload authorization")
					}
					if err := r.ParseMultipartForm(1 << 20); err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					defer r.MultipartForm.RemoveAll()
					file, header, err := r.FormFile("file")
					if err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					defer file.Close()
					content, err := io.ReadAll(file)
					if err != nil || string(content) != video || header.Filename != "recording.mp4" {
						t.Errorf("uploaded recording = %q, filename = %q, error = %v", content, header.Filename, err)
					}
					_, _ = io.WriteString(w, "ok")
				case "/files.completeUploadExternal":
					if err := r.ParseForm(); err != nil {
						t.Error(err)
					}
					var files []slack.FileSummary
					if err := json.Unmarshal([]byte(r.Form.Get("files")), &files); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(files, []slack.FileSummary{{ID: "F123", Title: wantTitle}}) {
						t.Errorf("shared files = %+v", files)
					}
					if r.Form.Get("channel_id") != "C123" || r.Form.Get("initial_comment") != "Motion detected" {
						t.Errorf("share parameters = %v", r.Form)
					}
					_, _ = io.WriteString(w, `{"ok":true,"files":[{"id":"F123"}]}`)
				default:
					t.Errorf("unexpected request: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			client := server.Client()
			api := slack.New("test-token", slack.OptionAPIURL(server.URL+"/"), slack.OptionHTTPClient(client))
			sender := Slack{DeliveryMode: "bot", Token: "test-token", Channel: "C123", Hook: server.URL + "/hook"}
			err := sender.send(context.Background(), message.Message{
				Title: test.title, Body: "Motion detected",
				Media: []message.Media{{Type: "video", Url: server.URL + "/recording.mp4?signature=private"}},
			}, client, api)
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"/recording.mp4", "/files.getUploadURLExternal", "/upload", "/files.completeUploadExternal"}
			for _, path := range want {
				select {
				case got := <-requests:
					if got != path {
						t.Errorf("request = %s, want %s", got, path)
					}
				default:
					t.Fatalf("missing request: %s", path)
				}
			}
			if len(requests) != 0 {
				t.Fatal("sent extra requests, including a possible duplicate webhook")
			}
			assertSlackTempFilesCleaned(t, tempDir)
		})
	}
}

func TestSlackSelectedModeDoesNotFallbackOnProviderFailure(t *testing.T) {
	for _, mode := range []string{"webhook", "bot"} {
		t.Run(mode, func(t *testing.T) {
			calls := make(chan string, 2)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls <- r.URL.Path
				if r.URL.Path == "/hook" {
					w.WriteHeader(http.StatusServiceUnavailable)
				} else {
					_, _ = io.WriteString(w, `{"ok":false,"error":"missing_scope"}`)
				}
			}))
			defer server.Close()
			client := server.Client()
			api := slack.New("saved-token", slack.OptionAPIURL(server.URL+"/"), slack.OptionHTTPClient(client))
			sender := Slack{DeliveryMode: mode, Token: "saved-token", Channel: "C123", Hook: server.URL + "/hook"}
			if err := sender.send(context.Background(), message.Message{Body: "Test"}, client, api); err == nil {
				t.Fatal("provider failure was reported as success")
			}
			wantPath := "/hook"
			if mode == "bot" {
				wantPath = "/chat.postMessage"
			}
			if len(calls) != 1 || <-calls != wantPath {
				t.Fatal("provider failure used an inactive destination")
			}
		})
	}
}

func TestSlackBotNonVideoMessages(t *testing.T) {
	for _, media := range [][]message.Media{
		nil,
		{{Type: "image", Url: "https://media.example/photo.jpg", ThumbnailUrl: "https://media.example/thumb.jpg"}},
	} {
		t.Run(strconv.Itoa(len(media)), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				wantText := "Motion detected"
				if len(media) != 0 {
					wantText += "\r\n" + media[0].Url
					var attachments []slack.Attachment
					if err := json.Unmarshal([]byte(r.Form.Get("attachments")), &attachments); err != nil {
						t.Error(err)
					}
					if len(attachments) != 1 || attachments[0].ImageURL != media[0].ThumbnailUrl {
						t.Errorf("attachments = %+v", attachments)
					}
				}
				if r.URL.Path != "/chat.postMessage" || r.Form.Get("channel") != "C123" || r.Form.Get("text") != wantText {
					t.Errorf("unexpected text notification: %s %v", r.URL.Path, r.Form)
				}
				_, _ = io.WriteString(w, `{"ok":true,"channel":"C123","ts":"123.456"}`)
			}))
			defer server.Close()
			client := server.Client()
			api := slack.New("test-token", slack.OptionAPIURL(server.URL+"/"), slack.OptionHTTPClient(client))
			err := (Slack{Token: "test-token", Channel: "C123"}).send(context.Background(),
				message.Message{Body: "Motion detected", Media: media}, client, api)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSlackUploadFailuresDoNotFallBackToWebhook(t *testing.T) {
	for _, test := range []struct {
		name       string
		failPath   string
		status     int
		body       string
		wantError  string
		wantCalls  int
		wantStatus int
	}{
		{name: "expired media", failPath: "/video", status: 403, body: "private media error", wantError: "Forbidden", wantCalls: 1, wantStatus: 403},
		{name: "empty media", failPath: "/video", status: 200, wantError: "empty", wantCalls: 1},
		{name: "missing scope", failPath: "/files.getUploadURLExternal", status: 200, body: `{"ok":false,"error":"missing_scope"}`, wantError: "missing_scope", wantCalls: 2},
		{name: "upload rejected", failPath: "/upload", status: 500, body: "private upload error", wantError: "upload Slack video", wantCalls: 3, wantStatus: 500},
		{name: "share rejected", failPath: "/files.completeUploadExternal", status: 200, body: `{"ok":false,"error":"not_in_channel"}`, wantError: "not_in_channel", wantCalls: 4},
	} {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			t.Setenv("TMPDIR", tempDir)
			requests := make(chan string, 8)
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r.URL.Path
				_, _ = io.Copy(io.Discard, r.Body)
				if r.URL.Path == test.failPath {
					w.WriteHeader(test.status)
					_, _ = io.WriteString(w, test.body)
					return
				}
				switch r.URL.Path {
				case "/video":
					_, _ = io.WriteString(w, "video bytes")
				case "/files.getUploadURLExternal":
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"ok": true, "file_id": "F123", "upload_url": server.URL + "/upload",
					})
				case "/upload":
					_, _ = io.WriteString(w, "ok")
				default:
					t.Errorf("unexpected request after failure: %s", r.URL.Path)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))
			defer server.Close()
			client := server.Client()
			api := slack.New("test-token", slack.OptionAPIURL(server.URL+"/"), slack.OptionHTTPClient(client))
			err := (Slack{Token: "test-token", Channel: "C123", Hook: server.URL + "/hook"}).send(context.Background(),
				message.Message{Media: []message.Media{{Type: "video", Url: server.URL + "/video?signature=private"}}}, client, api)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("error = %v, want %q", err, test.wantError)
			}
			if strings.Contains(err.Error(), "private") || strings.Contains(err.Error(), server.URL) {
				t.Fatalf("failure exposed private media information: %v", err)
			}
			if test.wantStatus != 0 {
				var status interface{ HTTPStatusCode() int }
				if !errors.As(err, &status) || status.HTTPStatusCode() != test.wantStatus {
					t.Fatalf("error = %v, want HTTP status %d", err, test.wantStatus)
				}
			}
			if len(requests) != test.wantCalls {
				t.Fatalf("made %d requests, want %d", len(requests), test.wantCalls)
			}
			assertSlackTempFilesCleaned(t, tempDir)
		})
	}
}

func TestSlackBotRequiresTokenAndChannel(t *testing.T) {
	for _, sender := range []Slack{{Token: "test-token"}, {Channel: "C123"}, {DeliveryMode: "bot", Hook: "https://unused.example/hook"}} {
		if err := sender.Send(message.Message{}); err == nil || !strings.Contains(err.Error(), "both a bot token and channel ID") {
			t.Fatalf("error = %v, want configuration error", err)
		}
	}
}

func TestSlackSelectedDeliveryMode(t *testing.T) {
	for _, mode := range []string{"webhook", "bot"} {
		t.Run(mode, func(t *testing.T) {
			calls := make(chan string, 2)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls <- r.URL.Path
				if r.URL.Path == "/hook" {
					_, _ = io.WriteString(w, "ok")
				} else if r.URL.Path == "/chat.postMessage" {
					_, _ = io.WriteString(w, `{"ok":true,"channel":"C123","ts":"123.456"}`)
				} else {
					t.Errorf("unexpected Slack request %s", r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()
			client := server.Client()
			api := slack.New("saved-token", slack.OptionAPIURL(server.URL+"/"), slack.OptionHTTPClient(client))
			sender := Slack{DeliveryMode: mode, Token: "saved-token", Channel: "C123", Hook: server.URL + "/hook"}
			err := sender.send(context.Background(), message.Message{Body: "Test"}, client, api)
			if err != nil {
				t.Fatal(err)
			}
			wantPath := "/hook"
			if mode == "bot" {
				wantPath = "/chat.postMessage"
			}
			if len(calls) != 1 || <-calls != wantPath {
				t.Fatal("Slack sent to an inactive or duplicate destination")
			}
		})
	}
	for _, sender := range []Slack{
		{DeliveryMode: "unknown", Token: "saved-token", Channel: "C123"},
		{DeliveryMode: "webhook", Token: "saved-token", Channel: "C123"},
	} {
		if err := sender.Send(message.Message{}); err == nil {
			t.Fatal("invalid selected mode silently fell back to another destination")
		}
	}
}

func TestSlackSendCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, sender := range []Slack{
		{Hook: "http://127.0.0.1:1/hook"},
		{Token: "test-token", Channel: "C123"},
	} {
		err := sender.SendWithContext(ctx, message.Message{
			Media: []message.Media{{Type: "video", Url: "http://127.0.0.1:1/video?signature=private"}},
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want cancellation", err)
		}
		if sender.Token != "" && strings.Contains(err.Error(), "signature") {
			t.Fatalf("download error exposed a signed URL: %v", err)
		}
	}
}

func TestSlackRejectsOversizedDownload(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TMPDIR", tempDir)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", strconv.FormatInt(slackMaxVideoBytes+1, 10))
	}))
	defer server.Close()
	err := (Slack{Token: "test-token", Channel: "C123"}).Send(message.Message{
		Media: []message.Media{{Type: "video", Url: server.URL}},
	})
	if err == nil || !strings.Contains(err.Error(), "1 GB") {
		t.Fatalf("error = %v, want size limit", err)
	}
	assertSlackTempFilesCleaned(t, tempDir)
}

func TestCopySlackVideoSizeLimit(t *testing.T) {
	for _, size := range []int64{0, slackMaxVideoBytes, slackMaxVideoBytes + 1} {
		t.Run(strconv.FormatInt(size, 10), func(t *testing.T) {
			src := io.LimitReader(slackTestVideoReader{}, size)
			copied, err := copySlackVideo(io.Discard, src)
			if size == slackMaxVideoBytes {
				if err != nil || copied != size {
					t.Fatalf("copied %d bytes, error = %v", copied, err)
				}
			} else if err == nil {
				t.Fatalf("accepted invalid video size %d", size)
			}
		})
	}
}

type slackTestVideoReader struct{}

func (slackTestVideoReader) Read(p []byte) (int, error) {
	return len(p), nil
}

func assertSlackTempFilesCleaned(t *testing.T, dir string) {
	t.Helper()
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("Slack video temporary files remain: %v", files)
	}
}

// Setup and recording problems cannot be fixed by retrying, so the delivery
// classifies them as rejected.
func TestSlackSetupErrorsAreRejections(t *testing.T) {
	var rejected interface{ DeliveryRejected() bool }
	for _, sender := range []Slack{
		{Token: "xoxb-token"},
		{DeliveryMode: "unknown", Token: "xoxb-token", Channel: "C123"},
		{DeliveryMode: "webhook", Token: "xoxb-token", Channel: "C123"},
		{DeliveryMode: "bot", Hook: "https://unused.example/hook"},
	} {
		err := sender.Send(message.Message{})
		if !errors.As(err, &rejected) || !rejected.DeliveryRejected() {
			t.Fatalf("invalid configuration error = %#v, want a rejected delivery", err)
		}
	}
	if _, err := copySlackVideo(io.Discard, strings.NewReader("")); !errors.As(err, &rejected) {
		t.Fatalf("empty video error = %#v, want a rejected delivery", err)
	}
}
