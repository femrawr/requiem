package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadFile(url string, outPath string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if outPath == "" {
		outPath = os.TempDir()
	}

	name := strings.Split(url, "?")[0]
	extension := filepath.Ext(name)

	filePath := filepath.Join(
		outPath,
		fmt.Sprintf("%d%s", time.Now().UnixNano(), extension),
	)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func DownloadFiles(urls []string, outPath string) ([]string, error) {
	var paths []string

	for _, url := range urls {
		path, err := DownloadFile(url, outPath)
		if err != nil {
			return paths, err
		}

		paths = append(paths, path)
	}

	return paths, nil
}
