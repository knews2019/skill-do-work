//go:build unix

package ownedprocess

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const groupPollInterval = 20 * time.Millisecond

func configureOwnedGroup(command *exec.Cmd) bool {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return true
}

// groupMember is one process observed in the owned group. A zombie member has
// exited but has not been reaped, which matters in both directions: it can no
// longer execute, and it still answers kill(pid, 0) as though it were alive.
type groupMember struct {
	processID int
	parentID  int
	zombie    bool
}

// terminateOwnedGroup ends the group from the bottom up: descendants before
// their own parents, and the leader last.
//
// The order carries the correctness, because a process is reaped by its parent
// and by nobody else. A group-wide signal reaches the leader first, which
// orphans a still-running hook to init and leaves that hook an unreaped zombie
// once it dies — and a zombie still satisfies kill(pid, 0), so a caller that
// returns there has proved nothing about the hook being dead. Ending one level
// at a time keeps every dying process's parent alive to waitpid() it, which
// closes the orphan and the zombie together. git commit then fails on its own
// because its hook died, which the transaction's existing commit-failure path
// already turns into a rollback.
//
// PR_SET_CHILD_SUBREAPER would be the other way to own an orphan, but it is
// not in the standard library, and this module's dependencies are.
func terminateOwnedGroup(leaderPID int, grace time.Duration) error {
	if grace <= 0 {
		grace = DefaultGracePeriod
	}
	if processGroup, err := syscall.Getpgid(leaderPID); err != nil || processGroup != leaderPID {
		// The module's standing runtime boundary: a leader whose own process
		// group cannot be proved isolated is signalled by bare pid and never by
		// group, because -leaderPID could then name the caller's own group.
		return escalateOnProcess(leaderPID, grace)
	}
	members, enumerated := ownedGroupMembers(leaderPID)
	if !enumerated {
		return terminateWholeGroup(leaderPID, grace)
	}
	if len(members) == 0 {
		return os.ErrProcessDone
	}
	if terminateDescendants(leaderPID, grace) {
		// A leader whose descendants just died has the cleanest exit
		// available: git releases index.lock and reports the hook failure
		// itself. Signalling it mid-commit would leave that lock for the
		// rollback to trip over, so spend the grace window waiting first.
		awaitCondition(func() bool { return !leaderRuns(leaderPID) }, grace)
	}
	for _, escalation := range []syscall.Signal{syscall.SIGTERM, syscall.SIGKILL} {
		if !leaderRuns(leaderPID) {
			break
		}
		_ = syscall.Kill(leaderPID, escalation)
		awaitCondition(func() bool { return !leaderRuns(leaderPID) }, grace)
	}
	// Killing the leader orphans whatever it had not reaped, so sweep once more.
	terminateDescendants(leaderPID, grace)
	return nil
}

// terminateDescendants ends the group's descendants deepest level first and
// reports whether any signal was delivered.
//
// The levels come from one snapshot rather than being re-derived each round,
// and that is what bounds this. A parent that forks a replacement as fast as
// its children are signalled always has a live child, so a rule of "only
// signal a process with no live children" would stay at the bottom of the tree
// forever and never return. Working through a snapshot's levels signals such a
// parent on its own level regardless of what it has spawned since.
func terminateDescendants(leaderPID int, grace time.Duration) bool {
	delivered := false
	// The first pass ends the tree that existed when cancellation arrived. A
	// parent that forks while its own level is being signalled can leave one
	// more generation behind as an orphan, so a second pass ends that. A third
	// would be chasing a process respawning faster than it can be signalled,
	// which the leader's own termination settles instead of this loop.
	for pass := 0; pass < 2; pass++ {
		tree, orphans, anyPresent := descendantLevels(leaderPID)
		if !anyPresent {
			break
		}
		for index := len(tree) - 1; index >= 0; index-- {
			if escalateOnMembers(leaderPID, tree[index], grace, true) {
				delivered = true
			}
		}
		if escalateOnMembers(leaderPID, orphans, grace, false) {
			delivered = true
		}
	}
	return delivered
}

// escalateOnMembers signals one set of group members from SIGTERM to SIGKILL
// and waits for them to go.
//
// requireReaped separates the two things "gone" can mean. A tree member's
// parent is still alive, so it must actually disappear: a zombie is exactly
// what an external caller checking "is the hook dead" with kill(pid, 0) still
// sees. An orphan belongs to init, so the most that can be proved about it is
// that it stopped running, and waiting for its reap would only be waiting on
// init's schedule.
func escalateOnMembers(leaderPID int, processIDs []int, grace time.Duration, requireReaped bool) bool {
	if len(processIDs) == 0 {
		return false
	}
	gone := func() bool { return !anyMember(leaderPID, processIDs, requireReaped) }
	delivered := false
	for _, escalation := range []syscall.Signal{syscall.SIGTERM, syscall.SIGKILL} {
		if gone() {
			break
		}
		for _, processID := range processIDs {
			// A pid can be recycled between reading the process table and
			// signalling it. Confirming group membership immediately before the
			// signal narrows that to one syscall, so the signal cannot land on
			// an unrelated process that inherited the number.
			if processGroup, err := syscall.Getpgid(processID); err != nil || processGroup != leaderPID {
				continue
			}
			if syscall.Kill(processID, escalation) == nil {
				delivered = true
			}
		}
		if awaitCondition(gone, grace) {
			break
		}
	}
	return delivered
}

// descendantLevels groups the group's live descendants by distance from the
// leader, shallowest level first, and separates those whose parent has already
// left the group. The split decides who can reap them: a tree member's parent
// is alive and reaps it in milliseconds, while an orphan belongs to init. The
// third result reports whether any descendant is present at all, zombies
// included, which is a different question from whether any can be signalled.
func descendantLevels(leaderPID int) ([][]int, []int, bool) {
	members, enumerated := ownedGroupMembers(leaderPID)
	if !enumerated {
		return nil, nil, false
	}
	isMember := make(map[int]bool, len(members))
	for _, member := range members {
		isMember[member.processID] = true
	}
	liveChildren := make(map[int][]int, len(members))
	liveDescendants := make([]int, 0, len(members))
	anyPresent := false
	for _, member := range members {
		if member.processID == leaderPID {
			continue
		}
		anyPresent = true
		if member.zombie {
			// A zombie needs its parent to reap it, not a signal.
			continue
		}
		liveDescendants = append(liveDescendants, member.processID)
		if isMember[member.parentID] {
			liveChildren[member.parentID] = append(liveChildren[member.parentID], member.processID)
		}
	}
	var tree [][]int
	placed := make(map[int]bool, len(liveDescendants))
	for level := []int{leaderPID}; ; {
		next := make([]int, 0, len(liveDescendants))
		for _, parentID := range level {
			for _, childID := range liveChildren[parentID] {
				if placed[childID] {
					continue
				}
				placed[childID] = true
				next = append(next, childID)
			}
		}
		if len(next) == 0 {
			break
		}
		tree = append(tree, next)
		level = next
	}
	var orphans []int
	for _, processID := range liveDescendants {
		if !placed[processID] {
			orphans = append(orphans, processID)
		}
	}
	return tree, orphans, anyPresent
}

// anyMember reports whether any of the given pids is still in the group.
// includeZombies picks between "still visible to kill(pid, 0)" and "still
// executing".
func anyMember(leaderPID int, processIDs []int, includeZombies bool) bool {
	members, enumerated := ownedGroupMembers(leaderPID)
	if !enumerated {
		return false
	}
	wanted := make(map[int]struct{}, len(processIDs))
	for _, processID := range processIDs {
		wanted[processID] = struct{}{}
	}
	for _, member := range members {
		if _, found := wanted[member.processID]; !found {
			continue
		}
		if member.zombie && !includeZombies {
			continue
		}
		return true
	}
	return false
}

// leaderRuns reports whether the leader is still executing. A zombie leader is
// not: it is the caller's own child, and the caller's Cmd.Wait reaps it.
func leaderRuns(leaderPID int) bool {
	members, enumerated := ownedGroupMembers(leaderPID)
	if !enumerated {
		return syscall.Kill(leaderPID, 0) == nil
	}
	for _, member := range members {
		if member.processID == leaderPID {
			return !member.zombie
		}
	}
	return false
}

// escalateOnProcess signals one process from SIGTERM to SIGKILL, touching no
// group. It is the degraded path for a leader whose group ownership was never
// established, so it can reach only the launched process itself.
func escalateOnProcess(processID int, grace time.Duration) error {
	if syscall.Kill(processID, 0) != nil {
		return os.ErrProcessDone
	}
	for _, escalation := range []syscall.Signal{syscall.SIGTERM, syscall.SIGKILL} {
		if syscall.Kill(processID, 0) != nil {
			break
		}
		_ = syscall.Kill(processID, escalation)
		awaitCondition(func() bool { return syscall.Kill(processID, 0) != nil }, grace)
	}
	return nil
}

// terminateWholeGroup is the fallback for a Unix host whose process table
// cannot be read: without member pids there is no way to spare the leader or
// order the levels, so this is the group-wide escalation, with the orphan
// window that ordering exists to avoid.
func terminateWholeGroup(leaderPID int, grace time.Duration) error {
	if syscall.Kill(-leaderPID, 0) != nil {
		return os.ErrProcessDone
	}
	for _, escalation := range []syscall.Signal{syscall.SIGTERM, syscall.SIGKILL} {
		if syscall.Kill(-leaderPID, 0) != nil {
			break
		}
		_ = syscall.Kill(-leaderPID, escalation)
		awaitCondition(func() bool { return syscall.Kill(-leaderPID, 0) != nil }, grace)
	}
	return nil
}

// ownedGroupMembers lists the process table entries whose process group is
// leaderPID. The boolean reports whether the table could be read at all, which
// is a different answer from an empty group.
func ownedGroupMembers(leaderPID int) ([]groupMember, bool) {
	output, err := exec.Command("ps", "-eo", "pgid=,pid=,ppid=,stat=").Output()
	if err != nil {
		return nil, false
	}
	members := make([]groupMember, 0, 4)
	for _, line := range bytes.Split(output, []byte{'\n'}) {
		fields := strings.Fields(string(line))
		if len(fields) < 4 {
			continue
		}
		group, groupErr := strconv.Atoi(fields[0])
		processID, processErr := strconv.Atoi(fields[1])
		parentID, parentErr := strconv.Atoi(fields[2])
		if groupErr != nil || processErr != nil || parentErr != nil || group != leaderPID {
			continue
		}
		members = append(members, groupMember{
			processID: processID,
			parentID:  parentID,
			zombie:    strings.HasPrefix(fields[3], "Z"),
		})
	}
	return members, true
}

func awaitCondition(satisfied func() bool, budget time.Duration) bool {
	deadline := time.Now().Add(budget)
	for {
		if satisfied() {
			return true
		}
		if !time.Now().Before(deadline) {
			return false
		}
		time.Sleep(groupPollInterval)
	}
}
