package templates

import (
	"strings"
	"testing"
)

func TestInvitationTemplates(t *testing.T) {
	html := GetTemplate("invitation")
	text := GetTextTemplate("invitation")
	if html == "" || text == "" {
		t.Fatal("invitation templates must be registered")
	}
	for _, template := range []string{html, text} {
		for _, variable := range []string{"{{user}}", "{{organisation}}", "{{url}}", "{{expiry}}"} {
			if !strings.Contains(template, variable) {
				t.Fatalf("invitation template is missing %s", variable)
			}
		}
		if strings.Contains(strings.ToLower(template), "recording") {
			t.Fatal("the invitation template must not word the email as a shared recording")
		}
	}
	if html == GetTemplate("share") {
		t.Fatal("invitation must be a template of its own")
	}
}
