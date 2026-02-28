package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func CopyFile(srcPath string, outPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	defer src.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return err
	}

	return out.Sync()
}

func ZipDir(dirPath string) (string, error) {
	filePath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("%s.zip", filepath.Base(dirPath)),
	)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	return filePath, filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		target, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		created, err := writer.Create(target)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(created, file)
		return err
	})
}

func HideFile(path string) error {
	cmd := exec.Command("attrib", "+h", "+s", path)
	return cmd.Run()
}
