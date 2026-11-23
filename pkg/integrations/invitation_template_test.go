package integrations

import (
	"strings"
	"testing"

	"github.com/uug-ai/integrations/pkg/message"
	"github.com/uug-ai/integrations/pkg/templates"
)

// hub-api sends invitations with only Title and Data set; every variable the
// invitation template uses must be filled from them.
func TestInvitationTemplateRendersFromHubAPIData(t *testing.T) {
	m := message.Message{
		Title: "You're invited to join an organisation on Kerberos Hub",
		User:  "owner@example.com",
		Data: map[string]string{
			"user":         "Ada Owner",
			"organisation": "Acme Security",
			"url":          "https://app.example.com/invitations/token123",
			"expiry":       "25 November 2026 13:00 UTC",
		},
	}
	for name, body := range map[string]string{
		"html": ReplaceValues(templates.GetTemplate("invitation"), m),
		"text": ReplaceValues(templates.GetTextTemplate("invitation"), m),
	} {
		for _, want := range []string{"Ada Owner", "Acme Security", "https://app.example.com/invitations/token123", "25 November 2026 13:00 UTC"} {
			if !strings.Contains(body, want) {
				t.Fatalf("%s invitation is missing %q", name, want)
			}
		}
		for _, left := range []string{"{{user}}", "{{organisation}}", "{{url}}", "{{expiry}}", "{{title}}", "{{body}}"} {
			if strings.Contains(body, left) {
				t.Fatalf("%s invitation still contains %s", name, left)
			}
		}
	}
}
