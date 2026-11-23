package templates

import (
	"strings"
	"testing"
)

func TestDeviceGroupTemplates(t *testing.T) {
	html := GetTemplate("device_group")
	text := GetTextTemplate("device_group")
	if html == "" || text == "" {
		t.Fatal("device_group templates must be registered")
	}
	for _, variable := range []string{"{{user}}", "{{text}}", "{{sections}}", "{{link}}"} {
		if !strings.Contains(html, variable) {
			t.Fatalf("device_group HTML is missing %s", variable)
		}
	}
	for _, variable := range []string{"{{user}}", "{{text}}", "{{sectionsText}}", "{{link}}"} {
		if !strings.Contains(text, variable) {
			t.Fatalf("device_group text is missing %s", variable)
		}
	}
	if GetTemplate("device") == html {
		t.Fatal("device and device_group must be different templates")
	}
}
