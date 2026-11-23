package integrations

import (
	"strings"
	"testing"

	"github.com/uug-ai/integrations/pkg/message"
)

// TestReplaceValues_UserPrefersDisplayName guards the v1.3.8 contract: the
// {{user}} email variable must show the display name from Data["user"] when
// present, falling back to the routing username otherwise.
func TestReplaceValues_UserPrefersDisplayName(t *testing.T) {
	tmpl := "Hello {{user}}"

	withDisplay := ReplaceValues(tmpl, message.Message{
		User: "cedricve",
		Data: map[string]string{"user": "Cedric Verstraeten"},
	})
	if withDisplay != "Hello Cedric Verstraeten" {
		t.Fatalf("display name not used: %q", withDisplay)
	}

	fallback := ReplaceValues(tmpl, message.Message{User: "cedricve"})
	if fallback != "Hello cedricve" {
		t.Fatalf("username fallback not used: %q", fallback)
	}
}

// TestReplaceValues_LinkPrefersDataLink ensures the tenant-aware Data["link"]
// (from BuildAlertMediaLink) overrides the raw media URL for {{link}}, matching
// the Mailgun channel so the same alert links to the same place regardless of
// email provider.
func TestReplaceValues_LinkPrefersDataLink(t *testing.T) {
	tmpl := "Watch: {{link}}"
	out := ReplaceValues(tmpl, message.Message{
		Media: []message.Media{{Url: "https://raw/media.mp4"}},
		Data:  map[string]string{"link": "https://acme.kerberos.io/media/123"},
	})
	if out != "Watch: https://acme.kerberos.io/media/123" {
		t.Fatalf("Data[link] did not override media url: %q", out)
	}
}

// TestReplaceValues_LinkFallsBackToMedia ensures {{link}} still uses the raw
// media URL when no tenant-aware Data["link"] is present, so the link is never
// left empty.
func TestReplaceValues_LinkFallsBackToMedia(t *testing.T) {
	tmpl := "Watch: {{link}}"
	out := ReplaceValues(tmpl, message.Message{
		Media: []message.Media{{Url: "https://raw/media.mp4"}},
	})
	if out != "Watch: https://raw/media.mp4" {
		t.Fatalf("media url fallback not used: %q", out)
	}
}

// TestReplaceValues_MediaLink builds the app media link from MediaId.
func TestReplaceValues_MediaLink(t *testing.T) {
	out := ReplaceValues("Open {{medialink}}", message.Message{MediaId: "abc123"})
	if out != "Open https://app.kerberos.io/media/abc123" {
		t.Fatalf("medialink not built: %q", out)
	}
}

// TestReplaceValues_NoPanicOnEmptyMedia ensures the email body builder is safe
// when no media is attached (a panic here would lose the email notification).
func TestReplaceValues_NoPanicOnEmptyMedia(t *testing.T) {
	out := ReplaceValues("{{title}} {{text}}", message.Message{
		Title: "Alert",
		Body:  "Something happened",
	}) // must not panic
	if !strings.Contains(out, "Alert") || !strings.Contains(out, "Something happened") {
		t.Fatalf("title/text not substituted: %q", out)
	}
}

// TestReplaceValues_NoLeftoverCommonPlaceholders asserts the variables that
// always have a value are fully substituted, so recipients never receive a
// template with raw {{user}}/{{text}}/{{link}}/{{medialink}} tokens.
func TestReplaceValues_NoLeftoverCommonPlaceholders(t *testing.T) {
	tmpl := "{{user}} | {{text}} | {{link}} | {{medialink}}"
	out := ReplaceValues(tmpl, message.Message{
		User:    "cedricve",
		Body:    "body",
		MediaId: "abc123",
		Media:   []message.Media{{Url: "https://raw/media.mp4"}},
		Data:    map[string]string{"user": "Cedric Verstraeten", "link": "https://acme.kerberos.io/media/abc123"},
	})
	for _, leftover := range []string{"{{user}}", "{{text}}", "{{link}}", "{{medialink}}"} {
		if strings.Contains(out, leftover) {
			t.Fatalf("unsubstituted placeholder %s in %q", leftover, out)
		}
	}
}
