package internal

// NewTemplateGetter creates a TemplateFunc from a map of template names to content.
// This is used to simplify template handling in bundlers by converting embedded
// templates into a standard lookup function.
//
// Example usage:
//
//	//go:embed templates/README.md.tmpl
//	var readmeTemplate string
//
//	var GetTemplate = NewTemplateGetter(map[string]string{
//	    "README.md": readmeTemplate,
//	})
func NewTemplateGetter(templates map[string]string) TemplateFunc {
	return func(name string) (string, bool) {
		tmpl, ok := templates[name]
		return tmpl, ok
	}
}

// StandardTemplates returns a TemplateFunc for components that only have a README template.
// This is the most common case and reduces boilerplate further.
//
// Example usage:
//
//	//go:embed templates/README.md.tmpl
//	var readmeTemplate string
//
//	var GetTemplate = StandardTemplates(readmeTemplate)
func StandardTemplates(readmeTemplate string) TemplateFunc {
	return NewTemplateGetter(map[string]string{
		"README.md": readmeTemplate,
	})
}
