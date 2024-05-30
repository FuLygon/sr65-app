package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sr65-software/logger"
)

//go:embed bin/*
var embeddedBinFS embed.FS

func ExtractBinaries() (func(), error) {
	// create temporary directory
	tmpDir, err := os.MkdirTemp("", "sr65-software-*")
	if err != nil {
		logger.Error("error creating temporary directory", err)
		return func() {}, err
	}

	// extract embedded binaries
	err = extractEmbedFiles(embeddedBinFS, tmpDir, "bin")
	if err != nil {
		logger.Error("error extracting embedded binaries", err)
		return func() {}, err
	}

	// add temporary directory to PATH
	pathEnv := os.Getenv("PATH")
	pathEnv = fmt.Sprintf("%s%c%s", tmpDir, os.PathListSeparator, pathEnv)
	err = os.Setenv("PATH", pathEnv)
	if err != nil {
		logger.Error("error updating PATH", err)
		return func() {}, err
	}

	return func() { os.RemoveAll(tmpDir) }, nil
}

func extractEmbedFiles(embedFS embed.FS, dir string, root string) error {
	return fs.WalkDir(embedFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skipping directories or .gitkeep
		if d.IsDir() || d.Name() == ".gitkeep" {
			return nil
		}

		// read files from embedded fs
		data, err := embedFS.ReadFile(path)
		if err != nil {
			return err
		}

		// create file path
		targetPath := filepath.Join(dir, path[len(root):])
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}

		// write file
		return os.WriteFile(targetPath, data, 0755)
	})
}
