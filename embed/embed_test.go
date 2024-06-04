package embed

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractBinaries(t *testing.T) {
	tmpDir, err := ExtractBinaries()
	if err != nil {
		t.Errorf("error extracting embedded binaries %v", err)

		// remove temporary directory whether it was created
		err = os.Remove(tmpDir)
		if err != nil {
			t.Logf("error removing temporary directory %v", err)
		}
		t.FailNow()
	}

	// check if the temporary directory is added to PATH
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	var tmpDirAdded bool
	for _, path := range paths {
		if strings.HasPrefix(path, filepath.Join(os.TempDir(), "sr65-app-")) {
			tmpDirAdded = true
			break
		}
	}
	if !tmpDirAdded {
		t.Errorf("temporary directory was not added to PATH")
	}

	// check if the temporary directory is removed after cleanup
	err = os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("error removing temporary directory %v", err)
		t.FailNow()
	}
	tmpDir = strings.Split(pathEnv, string(os.PathListSeparator))[0]
	_, err = os.ReadDir(tmpDir)
	if !os.IsNotExist(err) {
		t.Errorf("temporary directory was not removed after cleanup")
	}
}
