package embed

import (
	"embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sr65-software/logger"
)

//go:embed bin/*
var embeddedBinFS embed.FS

func ExtractingBinaries() (string, func(), error) {
	// create temporary directory
	tempDir, err := os.MkdirTemp("", "sr65-software-*")
	if err != nil {
		logger.LogError(logrus.ErrorLevel, "error creating temporary directory", err)
		return "", func() {}, err
	}

	// extract embedded binaries
	err = extractEmbedFiles(embeddedBinFS, tempDir, "bin")
	if err != nil {
		logger.LogError(logrus.ErrorLevel, "error extracting embedded binaries", err)
		return "", func() {}, err
	}

	// add temporary directory to PATH
	pathEnv := os.Getenv("PATH")
	pathEnv = fmt.Sprintf("%s%c%s", tempDir, os.PathListSeparator, pathEnv)
	err = os.Setenv("PATH", pathEnv)
	if err != nil {
		logger.Log.Error("error updating PATH: ", err)
		return "", func() {}, err
	}

	return tempDir, func() { os.RemoveAll(tempDir) }, nil
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

func ffmpeg() {
	// Now you can execute ffmpeg
	cmd := exec.Command("ffmpeg", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
