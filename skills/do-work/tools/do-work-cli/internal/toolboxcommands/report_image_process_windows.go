//go:build windows

package toolboxcommands

import (
	"errors"
	"os/exec"
)

func configureOwnedProcess(*exec.Cmd) error {
	return errors.New("report image generation is unavailable: descendant process ownership cannot be proved on Windows")
}

func terminateOwnedProcess(int) error { return nil }
func killOwnedProcess(int) error      { return nil }
func ownedProcessGroupAlive(int) bool { return false }
