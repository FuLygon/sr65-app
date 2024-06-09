//go:build !windows

package internal

import (
	"errors"
	"os/exec"
	"runtime"
	"syscall"
)

func OpenInExplorer(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		return errors.New("unsupported platform")
	}

	// detach process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
