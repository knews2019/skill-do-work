//go:build !windows

package toolboxcommands

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
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

func ownedProcessGroupAlive(processID int) bool {
	output, err := exec.Command("ps", "-eo", "pgid=,stat=").Output()
	if err != nil {
		return syscall.Kill(-processID, 0) == nil
	}
	for _, line := range bytes.Split(output, []byte{'\n'}) {
		fields := strings.Fields(string(line))
		if len(fields) < 2 {
			continue
		}
		group, parseErr := strconv.Atoi(fields[0])
		if parseErr == nil && group == processID && !strings.HasPrefix(fields[1], "Z") {
			return true
		}
	}
	return false
}
