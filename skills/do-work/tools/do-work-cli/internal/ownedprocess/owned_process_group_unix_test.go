//go:build unix

package ownedprocess

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The escalation is SIGTERM then SIGKILL, and the pause between them is the
// only thing that makes the SIGTERM mean anything. Nothing else in the module
// notices if that pause disappears: sending both signals back to back still
// ends the group, and every cancellation test stays green while a subprocess
// loses its chance to clean up after itself. This pins the window with a
// process that can only exit gracefully.
func TestTerminateGroupLetsTheGracefulSignalRunFirst(t *testing.T) {
	markerPath := filepath.Join(t.TempDir(), "graceful")
	// `read` is a builtin, so this process has no children and is itself the
	// process the escalation signals. Its stdin never closes, so only the
	// signal can end it.
	command := exec.Command("sh", "-c", "trap 'echo graceful > "+markerPath+"; exit 0' TERM; read line")
	stdinReader, stdinWriter, pipeError := os.Pipe()
	if pipeError != nil {
		t.Fatal(pipeError)
	}
	defer stdinWriter.Close()
	command.Stdin = stdinReader
	if !ConfigureGroup(command) {
		t.Skip("this platform cannot own a process group")
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	stdinReader.Close()
	waited := make(chan error, 1)
	go func() { waited <- command.Wait() }()

	if err := TerminateGroup(command.Process.Pid, time.Second); err != nil {
		t.Fatalf("TerminateGroup = %v", err)
	}
	select {
	case <-waited:
	case <-time.After(3 * time.Second):
		t.Fatal("the terminated leader was never reaped")
	}

	contents, readError := os.ReadFile(markerPath)
	if readError != nil {
		t.Fatalf("the graceful signal had no window to run its handler: %v", readError)
	}
	if strings.TrimSpace(string(contents)) != "graceful" {
		t.Fatalf("marker = %q, want the handler's own bytes", contents)
	}
	if status, ok := command.ProcessState.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		t.Errorf("the leader was killed by %v rather than exiting from its own handler", status.Signal())
	}
}

// A group with nothing in it is a different answer from a group that could not
// be read, and the caller distinguishes them: os.ErrProcessDone tells os/exec
// not to inject a cancellation error over a command that already finished.
func TestTerminateGroupReportsAnAlreadyFinishedGroup(t *testing.T) {
	command := exec.Command("sh", "-c", "exit 0")
	if !ConfigureGroup(command) {
		t.Skip("this platform cannot own a process group")
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	leaderPID := command.Process.Pid
	if err := command.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := TerminateGroup(leaderPID, time.Second); err != os.ErrProcessDone {
		t.Fatalf("TerminateGroup on a finished group = %v, want %v", err, os.ErrProcessDone)
	}
}

// A backend that respawns a helper in a loop always has a live child. A sweep
// that only signals members with no live children therefore never climbs past
// the bottom of the tree, never signals the respawning parent, and never
// returns — the caller hangs until something outside kills it. The parent here
// runs its own respawn loop and starts a nested one, so both levels have to be
// ended from a snapshot of the tree rather than from whatever is childless at
// the moment each signal is sent.
func TestTerminateGroupEndsParentsThatKeepForkingChildren(t *testing.T) {
	nestedMarker := filepath.Join(t.TempDir(), "nested.pid")
	// Each respawning shell owns a deadline even if the test process is killed.
	// Thirty seconds outlasts all failure assertions below; expiry cannot prove cleanup.
	respawnLoop := "respawn_deadline=$((SECONDS + 30)); while (( SECONDS < respawn_deadline )); do sleep 0.05; done"
	command := exec.Command("bash", "-c",
		"( "+respawnLoop+" ) & printf '%s\\n' \"$!\" > "+nestedMarker+
			"; "+respawnLoop)
	if !ConfigureGroup(command) {
		t.Skip("this platform cannot own a process group")
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	waited := make(chan error, 1)
	go func() { waited <- command.Wait() }()
	nestedPID := awaitRespawningChild(t, nestedMarker)

	returned := make(chan error, 1)
	go func() { returned <- TerminateGroup(command.Process.Pid, 200*time.Millisecond) }()
	select {
	case err := <-returned:
		if err != nil {
			t.Fatalf("TerminateGroup = %v", err)
		}
	case <-time.After(5 * time.Second):
		_ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
		t.Fatal("TerminateGroup never returned: the sweep is chasing respawned children instead of ending their parent")
	}
	select {
	case <-waited:
	case <-time.After(3 * time.Second):
		t.Fatal("the terminated leader was never reaped")
	}

	if syscall.Kill(nestedPID, 0) == nil {
		_ = syscall.Kill(nestedPID, syscall.SIGKILL)
		t.Errorf("the nested respawn loop %d survived cancellation", nestedPID)
	}
	// A helper forked in the last instant before its parent died is an orphan
	// that init reaps on its own schedule, so allow it to drain.
	if !awaitCondition(func() bool { return !anyLiveGroupMember(command.Process.Pid) }, 2*time.Second) {
		members, _ := ownedGroupMembers(command.Process.Pid)
		t.Errorf("group members were still running after cancellation: %+v", members)
	}
}

func anyLiveGroupMember(leaderPID int) bool {
	members, enumerated := ownedGroupMembers(leaderPID)
	if !enumerated {
		return false
	}
	for _, member := range members {
		if !member.zombie {
			return true
		}
	}
	return false
}

func awaitRespawningChild(t *testing.T, markerPath string) int {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		contents, err := os.ReadFile(markerPath)
		if err == nil {
			if processID, convertError := strconv.Atoi(strings.TrimSpace(string(contents))); convertError == nil && processID > 0 {
				return processID
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("no nested process id was recorded at %s", markerPath)
	return 0
}
