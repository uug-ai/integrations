package integrations

import (
	"os"
	"testing"
	"testing/fstest"
)

func TestTemplatesOverrideReplacesBuiltinFileByFile(t *testing.T) {
	templates := Templates{
		Override: fstest.MapFS{"detection.html": {Data: []byte("custom html")}},
		Builtin: fstest.MapFS{
			"detection.html": {Data: []byte("builtin html")},
			"detection.txt":  {Data: []byte("builtin text")},
		},
	}
	got := templates.Get("detection")
	if got.HTML != "custom html" || got.Text != "builtin text" {
		t.Fatalf("Get = %+v", got)
	}
}

func TestTemplatesMissingGivesEmptyBodies(t *testing.T) {
	templates := Templates{Builtin: fstest.MapFS{"detection.html": {Data: []byte("x")}}}
	if got := templates.Get("unknown"); got != (EmailTemplate{}) {
		t.Fatalf("unknown template = %+v", got)
	}
	if got := templates.Get(""); got != (EmailTemplate{}) {
		t.Fatalf("empty name = %+v", got)
	}
	if got := (Templates{}).Get("detection"); got != (EmailTemplate{}) {
		t.Fatalf("no sources = %+v", got)
	}
}

func TestTemplatesRejectPathsOutsideTheSource(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/secret.html", []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	templates := Templates{Override: os.DirFS(dir + "/sub")}
	if got := templates.Get("../secret"); got.HTML != "" {
		t.Fatalf("template name escaped the directory: %q", got.HTML)
	}
}

func TestSMTPFromEnvKeepsTemplates(t *testing.T) {
	t.Setenv("SMTP_SERVER", "smtp.example.com")
	t.Setenv("SMTP_PORT", "25")
	t.Setenv("EMAIL_FROM", "hub@example.com")
	templates := Templates{Builtin: os.DirFS(t.TempDir())}
	email, err := EmailFromEnv(EmailProviderSMTP, "to@example.com", "detection", templates)
	if err != nil {
		t.Fatal(err)
	}
	if got := email.(*SMTP).options.Templates; got != templates {
		t.Fatalf("templates = %+v, want %+v", got, templates)
	}
}
