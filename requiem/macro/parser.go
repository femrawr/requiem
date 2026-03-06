package macro

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ParseMacro(filePath string) (*Macro, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan() // header
	scanner.Scan() // separator

	macro := &Macro{}

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

func parseLine(line string) (Line, error) {
	cleaned := removeComment(line)

	fields := strings.Fields(cleaned)
	if len(fields) == 0 {
		return Line{}, errors.New("the line is empty")
	}

	symbol := fields[0]
	id, ok := symbols[symbol]
	if !ok {
		return Line{}, fmt.Errorf("invalid symbol: \"%s\"", symbol)
	}

	if len(fields) < 2 {
		return Line{}, fmt.Errorf("symbol \"%s\" does not have a trailing value", symbol)
	}

	return Line{
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
