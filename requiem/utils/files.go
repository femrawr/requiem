package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

		rel, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		file, err := writer.Create(rel)
		if err != nil {
			return err
		}

		open, err := os.Open(path)
		if err != nil {
			return err
		}

		defer open.Close()

		_, err = io.Copy(file, open)
		return err
	})
}

func GenFileTree(root string, maxDepth int) (string, error) {
	var tree strings.Builder

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		depth := strings.Count(rel, string(os.PathSeparator))
		if depth > maxDepth {
			if entry.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		indent := strings.Repeat("  ", depth)

		if entry.IsDir() {
			fmt.Fprintf(&tree, "%s📁 %s\n", indent, entry.Name())
		} else {
			fmt.Fprintf(&tree, "%s📄 %s\n", indent, entry.Name())
		}

		return nil
	})

	return tree.String(), err
}

func HideFile(path string) error {
	cmd := exec.Command("attrib", "+h", "+s", path)
	return cmd.Run()
}
