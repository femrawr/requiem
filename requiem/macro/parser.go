package macro

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ParseMacro(filePath string) (*macroData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan() // header
	scanner.Scan() // separator

	macro := &macroData{}

	atLine := 2

	for scanner.Scan() {
		atLine += 1

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		parsed, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("error on line %d: %s", atLine, err)
		}

		macro.Lines = append(macro.Lines, parsed)
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return macro, nil
}

func parseLine(line string) (macroDataLine, error) {
	cleaned := removeComment(line)

	fields := strings.Fields(cleaned)
	if len(fields) == 0 {
		return macroDataLine{}, errors.New("the line is empty")
	}

	symbol := fields[0]
	id, ok := symbols[symbol]
	if !ok {
		return macroDataLine{}, fmt.Errorf("invalid symbol: %q", symbol)
	}

	if len(fields) < 2 {
		return macroDataLine{}, fmt.Errorf("symbol %q does not have a trailing value", symbol)
	}

	return macroDataLine{
		Line:   line,
		Symbol: symbol,
		ID:     id,
		Value:  fields[1],
		Args:   fields[2:],
	}, nil
}

func removeComment(line string) string {
	things, _, _ := strings.Cut(line, "#")
	return strings.TrimSpace(things)
}
