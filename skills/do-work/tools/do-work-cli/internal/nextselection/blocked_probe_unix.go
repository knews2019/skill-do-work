//go:build unix

package nextselection

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func runOwnedProbe(probeBytes []byte, timeout time.Duration) (int, error) {
	command := exec.Command("sh", "-c", string(probeBytes))
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := command.Start(); err != nil {
		return BlockedProbeLaunchStatus, err
	}
	processGroup, err := syscall.Getpgid(command.Process.Pid)
	if err != nil || processGroup != command.Process.Pid {
		_ = command.Process.Kill()
		_ = command.Wait()
		if err == nil {
			err = errors.New("probe process group was not isolated")
		}
		return BlockedProbeLaunchStatus, err
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalChannel)
	select {
	case waitError := <-done:
		status := probeExitStatus(waitError)
		cleanupReapedProcessGroup(processGroup)
		return status, nil
	case <-timer.C:
		terminateOwnedProcessGroup(processGroup, syscall.SIGTERM, done)
		return BlockedProbeTimeoutStatus, nil
	case received := <-signalChannel:
		forwarded, ok := received.(syscall.Signal)
		if !ok {
			forwarded = syscall.SIGTERM
		}
		terminateOwnedProcessGroup(processGroup, forwarded, done)
		return 128 + int(forwarded), nil
	}
}

func cleanupReapedProcessGroup(processGroup int) {
	if syscall.Kill(-processGroup, 0) != nil {
		return
	}
	_ = syscall.Kill(-processGroup, syscall.SIGTERM)
	time.Sleep(100 * time.Millisecond)
	if syscall.Kill(-processGroup, 0) == nil {
		_ = syscall.Kill(-processGroup, syscall.SIGKILL)
	}
}

func terminateOwnedProcessGroup(processGroup int, initialSignal syscall.Signal, done <-chan error) {
	_ = syscall.Kill(-processGroup, initialSignal)
	grace := time.NewTimer(500 * time.Millisecond)
	defer grace.Stop()
	leaderReaped := false
	select {
	case <-done:
		leaderReaped = true
	case <-grace.C:
	}
	if leaderReaped {
		select {
		case <-grace.C:
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
	if syscall.Kill(-processGroup, 0) == nil {
		_ = syscall.Kill(-processGroup, syscall.SIGKILL)
	}
	if !leaderReaped {
		<-done
	}
}

func probeExitStatus(err error) int {
	if err == nil {
		return 0
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if status.Signaled() {
				return 128 + int(status.Signal())
			}
			return status.ExitStatus()
		}
	}
	return BlockedProbeLaunchStatus
}
