// Package ownedprocess launches subprocesses as the leaders of their own
// process groups and terminates those groups as a unit, so cancelling a
// command reaches every descendant it spawned — a Git hook, and whatever that
// hook started, included.
//
// Callers differ in what they may do when a platform cannot prove that
// ownership, so ConfigureGroup reports it instead of deciding: an image
// backend fails closed, while a Git transaction degrades to os/exec's default
// single-process cancellation because refusing would disable Git entirely.
package ownedprocess

import (
	"os/exec"
	"time"
)

// DefaultGracePeriod is the pause between the graceful signal and the
// escalation when a caller does not name its own budget.
const DefaultGracePeriod = 250 * time.Millisecond

// ConfigureGroup asks the platform to start command as the leader of a new
// process group and reports whether that ownership was established. A caller
// whose safety depends on reaching descendants must fail closed on false; a
// caller for which single-process cancellation is an acceptable degradation
// may ignore the result.
func ConfigureGroup(command *exec.Cmd) bool {
	return configureOwnedGroup(command)
}

// TerminateGroup ends every process in leaderPID's owned process group and
// blocks until they are gone, so a caller that returns afterwards has proved
// the group is dead rather than merely signalled. grace is the pause before
// each escalation; a non-positive value means DefaultGracePeriod. It returns
// os.ErrProcessDone when the group held nothing to signal.
func TerminateGroup(leaderPID int, grace time.Duration) error {
	return terminateOwnedGroup(leaderPID, grace)
}
