//go:build !windows

package toolboxcommands

import (
	"os/exec"
	"syscall"
)

func configureOwnedProcess(command *exec.Cmd) error {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return nil
}

func terminateOwnedProcess(processID int) error {
	return syscall.Kill(-processID, syscall.SIGTERM)
}

func killOwnedProcess(processID int) error {
	return syscall.Kill(-processID, syscall.SIGKILL)
}
