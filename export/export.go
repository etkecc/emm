package export

import (
	"text/template"

	"gitlab.com/etke.cc/tools/emm/matrix"
)

// DefaultTemplate text
const DefaultTemplate = `
id={{ .ID }}
replace={{ .Replace }}
author={{ .Author }}
text={{ .Text }}
html={{ .HTML }}
created_at={{ .CreatedAt }}
created_at_full={{ .CreatedAtFull }}
`

// Run export
func Run(templatePath, output string, messages []*matrix.Message) error {
	tpl, err := createTemplate(templatePath)
	if err != nil {
		return err
	}
	for _, message := range messages {
		err = save(tpl, output, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func save(tpl *template.Template, path string, message *matrix.Message) error {
	file, err := getOutput(path, message.ID)
	if err != nil {
		return err
	}

	return tpl.Execute(file, message)
}
