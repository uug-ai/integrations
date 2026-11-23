package templates

import (
	"strings"
	"testing"
)

func TestRecordingTemplates(t *testing.T) {
	html := GetTemplate("recording")
	text := GetTextTemplate("recording")
	if html == "" || text == "" {
		t.Fatal("recording templates must be registered")
	}
	for _, variable := range []string{"{{user}}", "{{text}}", "{{sections}}", "{{link}}"} {
		if !strings.Contains(html, variable) {
			t.Fatalf("recording HTML is missing %s", variable)
		}
	}
	for _, variable := range []string{"{{user}}", "{{text}}", "{{sectionsText}}", "{{link}}"} {
		if !strings.Contains(text, variable) {
			t.Fatalf("recording text is missing %s", variable)
		}
	}
	for _, template := range []string{html, text} {
		if strings.Contains(template, "go offline") || strings.Contains(template, "changed status") {
			t.Fatal("the recording template must not word the alert as a status change")
		}
		if !strings.Contains(template, "upload limit") {
			t.Fatal("the recording template must point at the upload limit")
		}
	}
	if html == GetTemplate("device") || html == GetTemplate("device_group") {
		t.Fatal("recording must be a template of its own")
	}
}

func TestDeviceTextTemplateIsNotTheWelcomeEmail(t *testing.T) {
	text := GetTextTemplate("device")
	if text == GetTextTemplate("welcome") || strings.Contains(text, "Activate your account") {
		t.Fatal("the device text template must describe the device change, not welcome the user")
	}
	for _, variable := range []string{"{{user}}", "{{text}}", "{{link}}"} {
		if !strings.Contains(text, variable) {
			t.Fatalf("device text is missing %s", variable)
		}
	}
}
