package export

import (
	"fmt"
	"os"
	"text/template"

	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/emm/matrix"
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
func Run(templatePath string, output string, messages []*matrix.Message) error {
	tpl, err := createTemplate(templatePath)
	if err != nil {
		return err
	}
	for _, message := range messages {
		if message.Replace != "" {
			remove(output, message.Replace)
		}
		err = save(tpl, output, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func remove(output string, eventID id.EventID) {
	if !isMulti(output) {
		return
	}

	// nolint // in 99.99% cases it will be "no such file or directory" error, which is ok.
	os.Remove(fmt.Sprintf(output, eventID))
}

func save(tpl *template.Template, path string, message *matrix.Message) error {
	file, err := getOutput(path, message.ID)
	if err != nil {
		return err
	}

	return tpl.Execute(file, message)
}
