package export

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/etkecc/emm/internal/utils"
)

var outputFileSingle *os.File

func getOutput(output string, vars map[string]string) (*os.File, error) {
	if isMulti(output) {
		return getOutputMulti(output, vars)
	}

	return getOutputSingle(output)
}

// isMulti check if output is a single file
func isMulti(output string) bool {
	return strings.Contains(output, "%s") || strings.Contains(output, "{{")
}

func getOutputSingle(output string) (*os.File, error) {
	var err error
	if outputFileSingle == nil {
		outputFileSingle, err = os.OpenFile(output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	}

	return outputFileSingle, err
}

func getOutputMulti(output string, vars map[string]string) (*os.File, error) {
	// old mode, where %s is used
	if strings.Contains(output, "%s") {
		output = strings.ReplaceAll(output, "%s", vars["ID"])
	}
	// new mode, with template support
	if strings.Contains(output, "{{") {
		parsed, err := parseTemplate(output, vars)
		if err != nil {
			return nil, err
		}
		output = parsed
	}
	return os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
}

// parseTemplate parse template with vars
func parseTemplate(output string, rawVars map[string]string) (string, error) {
	urlSafeVars := make(map[string]string, len(rawVars))
	for k, v := range rawVars {
		urlSafeVars[k] = utils.MakeURLSafe(v)
	}
	urlSafeVars["ID"] = rawVars["ID"] // special case - ID is, in fact, url-safe, but may use additional characters, like $, and we'd like to keep them, for backward compatibility

	tpl, err := template.New("output").Parse(output)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tpl.Execute(&buf, urlSafeVars); err != nil {
		return "", err
	}
	result := buf.String()
	// get file extension from the result
	ext := filepath.Ext(result)

	// we assume the template is in the last part of the path, so we should make it url-safe
	parts := strings.Split(result, "/")
	lastidx := len(parts) - 1
	last := parts[lastidx]
	if last == "" || last == ext {
		last = urlSafeVars["ID"] + ext
	}
	if len(last) > 100 {
		last = last[:100] + ext
	}
	parts[lastidx] = last
	result = strings.Join(parts, "/")
	return result, nil
}
