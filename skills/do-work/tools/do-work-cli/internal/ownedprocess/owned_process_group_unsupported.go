//go:build !unix

package ownedprocess

import (
	"os"
	"os/exec"
	"time"
)

// configureOwnedGroup declines: without a process group there is no way to
// prove which descendants a command spawned, and reporting false is what lets
// each caller pick between failing closed and degrading.
//
// The file is named for the condition rather than for windows so its !unix
// constraint is the only one that applies. A _windows.go name would add an
// implicit GOOS=windows constraint and leave every other non-unix target with
// no implementation at all.
func configureOwnedGroup(*exec.Cmd) bool { return false }

// terminateOwnedGroup has nothing to terminate, because no group was ever
// established. A caller that reached here never received a true from
// ConfigureGroup and is relying on os/exec's default cancellation.
func terminateOwnedGroup(int, time.Duration) error { return os.ErrProcessDone }
