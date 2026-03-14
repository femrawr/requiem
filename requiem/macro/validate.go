package macro

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ValidateFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = hasValidHeader(file)
	if err != nil {
		return err
	}

	err = hasValidFormat(file)
	if err != nil {
		return err
	}

	return nil
}

func hasValidHeader(file *os.File) error {
	buffer := make([]byte, len(MACRO_FILE_HEADER))

	_, err := io.ReadFull(file, buffer)
	if err != nil {
		return err
	}

	if string(buffer) != MACRO_FILE_HEADER {
		return errors.New("the file does not contain a valid macro file header")
	}

	return nil
}

func hasValidFormat(file *os.File) error {
	scanner := bufio.NewScanner(file)

	// this is the header
	scanner.Scan()

	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "" {
		return errors.New("the file does not contain valid content separator")
	}

	empty := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		empty = false

		clean, _, _ := strings.Cut(line, "#")

		fields := strings.Fields(clean)
		if len(fields) == 0 {
			continue
		}

		symbol := fields[0]
		if symbol != strings.ToUpper(symbol) {
			return fmt.Errorf("the file contains an invalid symbol: %q", symbol)
		}
	}

	err := scanner.Err()
	if err != nil {
		return err
	}

	if empty {
		return errors.New("the file does not contain any content")
	}

	return nil
}
