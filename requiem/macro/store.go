package macro

import "strings"

const (
	MACRO_FILE_HEADER string = "RMF"
)

var symbols = map[string]string{
	"CMD":  "0",
	"WAIT": "1",
}

type Line struct {
	Line   string
	Symbol string
	ID     string
	Value  string
	Args   []string
}

type Macro struct {
	Lines []Line
}

func (line Line) Encode() string {
	encoded := line.ID + "." + line.Value

	if len(line.Args) > 0 {
		encoded += "." + strings.Join(line.Args, ";")
	}

	return encoded
}

func (macro Macro) Encode() string {
	parts := make([]string, len(macro.Lines))

	for i, line := range macro.Lines {
		parts[i] = line.Encode()
	}

	return strings.Join(parts, "+")
}
