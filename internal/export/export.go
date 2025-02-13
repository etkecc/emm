package export

import (
	"os"
	"strings"
	"text/template"

	"github.com/etkecc/emm/internal/matrix"
	"maunium.net/go/mautrix/id"
)

// DefaultTemplate text
const DefaultTemplate = `
id={{ .ID }}
title={{ .Title }}
replace={{ .Replace }}
author={{ .Author }}
text={{ .Text }}
html={{ .HTML }}
created_at={{ .CreatedAt }}
created_at_full={{ .CreatedAtFull }}
`

// Run export
func Run(templatePath, output string, messages map[id.EventID]*matrix.Message) error {
	templatedOutput := strings.Contains(output, "{{")
	tpl, err := createTemplate(templatePath)
	if err != nil {
		return err
	}
	for _, message := range messages {
		// edge case for templated output: if the message is a replacement, we need to actually replace the original message
		if message.Replace != "" && messages[message.Replace] != nil && templatedOutput {
			err = save(tpl, output, message, messages[message.Replace])
		} else {
			err = save(tpl, output, message, nil)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func save(tpl *template.Template, path string, message, replaces *matrix.Message) error {
	var file *os.File
	var err error
	if replaces != nil {
		file, err = getOutput(path, replaces.Vars())
	} else {
		file, err = getOutput(path, message.Vars())
	}
	if err != nil {
		return err
	}

	return tpl.Execute(file, message)
}
