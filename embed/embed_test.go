package embed

import (
	"os"
	"strings"
	"testing"
)

func TestExtractBinaries(t *testing.T) {
	cleanup, err := ExtractBinaries()
	defer cleanup()

	if err != nil {
		t.Errorf("ExtractBinaries() error = %v", err)
		return
	}

	// check if the temporary directory is added to PATH
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	var tmpDirAdded bool
	for _, path := range paths {
		if strings.HasPrefix(path, "/tmp/sr65-software-") {
			tmpDirAdded = true
			break
		}
	}
	if !tmpDirAdded {
		t.Errorf("temporary directory was not added to PATH")
	}

	// check if the temporary directory is removed after cleanup
	cleanup()
	tmpDir := strings.Split(pathEnv, string(os.PathListSeparator))[0]
	_, err = os.ReadDir(tmpDir)
	if !os.IsNotExist(err) {
		t.Errorf("temporary directory was not removed after cleanup")
	}
}
