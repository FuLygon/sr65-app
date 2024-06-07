package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed bin/*
var embeddedBinFS embed.FS

// ExtractBinaries extract embedded binaries into a temporary directory and adds it to PATH
func ExtractBinaries() (string, error) {
	// create temporary directory
	tmpDir, err := os.MkdirTemp("", "sr65-app-*")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %w", err)
	}

	// extract embedded binaries
	err = extractEmbedFiles(embeddedBinFS, tmpDir, "bin")
	if err != nil {
		return "", fmt.Errorf("error extracting embedded binaries: %w", err)
	}

	// add temporary directory to PATH
	pathEnv := os.Getenv("PATH")
	pathEnv = fmt.Sprintf("%s%c%s", tmpDir, os.PathListSeparator, pathEnv)
	err = os.Setenv("PATH", pathEnv)
	if err != nil {
		return "", fmt.Errorf("error updating PATH: %w", err)
	}

	return tmpDir, nil
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
