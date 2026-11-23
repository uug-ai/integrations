package integrations

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendTwilioSMS(t *testing.T) {
	var gotPath, gotUser, gotPass, gotContentType string
	var gotForm map[string][]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotUser, gotPass, _ = r.BasicAuth()
		gotContentType = r.Header.Get("Content-Type")
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		gotForm = r.PostForm
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	if err := sendTwilioSMS(server.URL, "AC123", "secret", "+3200000001", "+3200000002", "hello"); err != nil {
		t.Fatalf("sendTwilioSMS returned error: %v", err)
	}

	if gotPath != "/Accounts/AC123/Messages.json" {
		t.Errorf("path = %q", gotPath)
	}
	if gotUser != "AC123" || gotPass != "secret" {
		t.Errorf("basic auth = %q/%q", gotUser, gotPass)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Errorf("content type = %q", gotContentType)
	}
	for key, want := range map[string]string{"From": "+3200000001", "To": "+3200000002", "Body": "hello"} {
		if got := strings.Join(gotForm[key], ","); got != want {
			t.Errorf("form %s = %q, want %q", key, got, want)
		}
	}
}

func TestSendTwilioSMSErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":20003,"message":"Authenticate"}`))
	}))
	defer server.Close()

	err := sendTwilioSMS(server.URL, "AC123", "wrong", "+3200000001", "+3200000002", "hello")
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("expected 401 error, got %v", err)
	}
}
