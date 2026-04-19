package macro

import "strings"

const (
	_MACRO_FILE_HEADER string = "RMF"
)

var symbols = map[string]string{
	"CMD":  "0",
	"WAIT": "1",
}

type macroDataLine struct {
	Line   string
	Symbol string
	ID     string
	Value  string
	Args   []string
}

type macroData struct {
	Lines []macroDataLine
}

func (line macroDataLine) Encode() string {
	encoded := line.ID + "." + line.Value

	if len(line.Args) > 0 {
		encoded += "." + strings.Join(line.Args, ";")
	}

	return encoded
}

func (macro macroData) Encode() string {
	parts := make([]string, len(macro.Lines))

	for i, line := range macro.Lines {
		parts[i] = line.Encode()
	}

	return strings.Join(parts, "+")
}
