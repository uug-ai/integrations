package integrations

import (
	"io/fs"
	"os"
)

// DefaultTemplateDir is where a deployment puts email templates that replace
// a service's built-in ones, as <name>.html and <name>.txt.
const DefaultTemplateDir = "/mail/templates"

// EmailTemplate is the HTML and plain-text body of an email. Both may hold
// {{variables}}, which ReplaceValues fills from the message.
type EmailTemplate struct {
	HTML string
	Text string
}

// Templates finds email templates by name in <name>.html and <name>.txt
// files. A file in Override replaces the same file in Builtin; either may be
// nil. A file found in neither gives an empty body.
type Templates struct {
	Override fs.FS
	Builtin  fs.FS
}

// NewTemplates returns a service's built-in templates, each replaced by the
// file of the same name in DefaultTemplateDir when there is one.
func NewTemplates(builtin fs.FS) Templates {
	return Templates{Override: os.DirFS(DefaultTemplateDir), Builtin: builtin}
}

// Get returns the template with the given name.
func (templates Templates) Get(name string) EmailTemplate {
	if name == "" {
		return EmailTemplate{}
	}
	return EmailTemplate{
		HTML: templates.read(name + ".html"),
		Text: templates.read(name + ".txt"),
	}
}

func (templates Templates) read(file string) string {
	for _, source := range []fs.FS{templates.Override, templates.Builtin} {
		if source == nil {
			continue
		}
		if content, err := fs.ReadFile(source, file); err == nil {
			return string(content)
		}
	}
	return ""
}
