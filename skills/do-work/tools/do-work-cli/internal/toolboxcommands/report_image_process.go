package toolboxcommands

import (
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/ownedprocess"
)

type ownedProcessResult struct {
	Status      int
	Interrupted bool
	Err         error
}

// reportImageGracePeriod is the pause before each escalation while a cancelled
// image backend's process group is torn down. It is a variable so a test can
// shorten it.
var reportImageGracePeriod = time.Second

// configureOwnedProcess fails closed for image backends: a backend that keeps
// running after cancellation can write into scratch this command is about to
// remove, so a platform that cannot prove descendant ownership must not run
// one at all. runGit makes the opposite choice for the same shared API,
// because Git cannot be given up.
func configureOwnedProcess(command *exec.Cmd) error {
	if !ownedprocess.ConfigureGroup(command) {
		return errors.New("report image generation is unavailable: descendant process ownership cannot be proved on this platform")
	}
	return nil
}

func runOwnedProcess(ctx context.Context, directory string, argv ...string) ownedProcessResult {
	if len(argv) == 0 {
		return ownedProcessResult{Status: 2, Err: errors.New("an executable is required")}
	}
	command := exec.Command(argv[0], argv[1:]...)
	command.Dir = directory
	if err := ctx.Err(); err != nil {
		return ownedProcessResult{Status: 1, Interrupted: true, Err: err}
	}
	if err := configureOwnedProcess(command); err != nil {
		return ownedProcessResult{Status: 1, Err: err}
	}
	if err := ctx.Err(); err != nil {
		return ownedProcessResult{Status: 1, Interrupted: true, Err: err}
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
		// TerminateGroup returns only once the group is gone, so the leader is
		// exited by the time Wait is collected and scratch cleanup cannot race
		// a backend that is still writing.
		_ = ownedprocess.TerminateGroup(command.Process.Pid, reportImageGracePeriod)
		<-done
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
