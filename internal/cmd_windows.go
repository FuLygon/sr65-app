//go:build windows

package internal

import (
	"os/exec"
	"syscall"
)

func OpenInExplorer(path string) error {
	cmd := exec.Command("explorer", path)
	err := cmd.Start()
	if err != nil {
		return err
	}

	// wait for process to finish
	err = cmd.Wait()
	if err != nil {
		return err
	}

	// set focus to explorer
	user32 := syscall.NewLazyDLL("user32.dll")
	setForegroundWindow := user32.NewProc("SetForegroundWindow")
	hwnd, _, _ := setForegroundWindow.Call(uintptr(cmd.Process.Pid))
	setForegroundWindow.Call(hwnd)

	return nil
}
