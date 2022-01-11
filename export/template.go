package export

import (
	"os"
	"text/template"
)

func createTemplate(templatePath string) (*template.Template, error) {
	var err error
	templateContent := DefaultTemplate
	if templatePath != "" {
		templateContent, err = loadTemplate(templatePath)
		if err != nil {
			return nil, err
		}
	}

	return template.New("emm").Parse(templateContent)
}

func loadTemplate(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
