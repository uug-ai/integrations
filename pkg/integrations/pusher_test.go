package integrations

import (
	"strings"
	"testing"

	"github.com/uug-ai/integrations/pkg/message"
)

// TestValidPusherChannel guards the contract that the Pusher routing key is the
// user's username (no spaces / Pusher-illegal characters), not a display name.
// A regression here means browser/real-time notifications are silently dropped.
func TestValidPusherChannel(t *testing.T) {
	cases := []struct {
		name      string
		recipient string
		want      bool
	}{
		{"username", "cedricve", true},
		{"username with dot", "cedric.ve", true},
		{"domain scoped", "acme@cedricve", true},
		{"display name with space", "Cedric Verstraeten", false},
		{"empty", "", false},
		{"tab", "cedric\tve", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := validPusherChannel(tc.recipient); got != tc.want {
				t.Fatalf("validPusherChannel(%q) = %v, want %v", tc.recipient, got, tc.want)
			}
		})
	}
}

// TestDisplayNameRejectedAsRoutingKey documents the exact v1.3.8 regression:
// a "First Last" display name must never pass the routing-key guard, so the
// routing key has to come from the username instead.
func TestDisplayNameRejectedAsRoutingKey(t *testing.T) {
	username := "cedricve"
	displayName := strings.TrimSpace("Cedric" + " " + "Verstraeten") // "Cedric Verstraeten"

	if !validPusherChannel(username) {
		t.Fatalf("username %q should be a valid pusher channel", username)
	}
	if validPusherChannel(displayName) {
		t.Fatalf("display name %q must not be usable as a routing key", displayName)
	}
}

// TestBuildPusherPayloadEmptyMedia guards against a panic when Media is empty.
// A panic here would be caught by the handler's recover and dead-letter the
// event, silently losing the in-app notification.
func TestBuildPusherPayloadEmptyMedia(t *testing.T) {
	msg := message.Message{Title: "t", Body: "b"} // no Media
	got := buildPusherPayload(msg)                // must not panic
	if len(got.Sequence.Media) != 0 {
		t.Fatalf("expected no media entries, got %d", len(got.Sequence.Media))
	}
	if got.Sequence.Title != "t" || got.Sequence.Text != "b" {
		t.Fatalf("title/text not mapped: %+v", got.Sequence)
	}
}

// TestBuildPusherPayloadMapsMedia asserts the first media item is mapped into
// the payload so the in-app notification keeps its thumbnail/link.
func TestBuildPusherPayloadMapsMedia(t *testing.T) {
	msg := message.Message{
		Title: "t",
		Body:  "b",
		Media: []message.Media{{Timestamp: 42, Type: "video", Url: "https://x/y.mp4"}},
	}
	got := buildPusherPayload(msg)
	if len(got.Sequence.Media) != 1 {
		t.Fatalf("expected 1 media entry, got %d", len(got.Sequence.Media))
	}
	m := got.Sequence.Media[0]
	if m.Title != "42" || m.Type != "video" || m.Media != "https://x/y.mp4" {
		t.Fatalf("media not mapped correctly: %+v", m)
	}
}
