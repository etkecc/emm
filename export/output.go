package export

import (
	"fmt"
	"os"
	"strings"

	"maunium.net/go/mautrix/id"
)

var outputFileSingle *os.File

func getOutput(output string, eventID id.EventID) (*os.File, error) {
	if isMulti(output) {
		return getOutputMulti(fmt.Sprintf(output, eventID))
	}

	return getOutputSingle(output)
}

// isMulti check if output is a single file
func isMulti(output string) bool {
	return strings.Contains(output, "%s")
}

func getOutputSingle(output string) (*os.File, error) {
	var err error
	if outputFileSingle == nil {
		outputFileSingle, err = os.OpenFile(output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	}

	return outputFileSingle, err
}

func getOutputMulti(output string) (*os.File, error) {
	return os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0o644)
}
