package toolboxcommands

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

type ownedProcessResult struct {
	Status      int
	Interrupted bool
	Err         error
}

var reportImageGracePeriod = time.Second

func runOwnedProcess(ctx context.Context, directory string, argv ...string) ownedProcessResult {
	if len(argv) == 0 {
		return ownedProcessResult{Status: 2, Err: errors.New("an executable is required")}
	}
	command := exec.Command(argv[0], argv[1:]...)
	command.Dir = directory
	if err := configureOwnedProcess(command); err != nil {
		return ownedProcessResult{Status: 1, Err: err}
	}
	if err := command.Start(); err != nil {
		return ownedProcessResult{Status: 1, Err: err}
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	select {
	case err := <-done:
		return completedProcessResult(err)
	case <-ctx.Done():
		_ = terminateOwnedProcess(command.Process.Pid)
		timer := time.NewTimer(reportImageGracePeriod)
		defer timer.Stop()
		select {
		case <-done:
		case <-timer.C:
			_ = killOwnedProcess(command.Process.Pid)
			<-done
		}
		return ownedProcessResult{Status: 1, Interrupted: true, Err: ctx.Err()}
	}
}

func completedProcessResult(err error) ownedProcessResult {
	if err == nil {
		return ownedProcessResult{}
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return ownedProcessResult{Status: exitError.ExitCode(), Err: err}
	}
	return ownedProcessResult{Status: 1, Err: err}
}
