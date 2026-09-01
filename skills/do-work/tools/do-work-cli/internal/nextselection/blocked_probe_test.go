//go:build unix

package nextselection

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestBlockedProbePreservesRawStatus(t *testing.T) {
	status, err := RunBlockedProbe([]byte("exit 37"), 2)
	if err != nil || status != 37 {
		t.Fatalf("status=%d err=%v", status, err)
	}
}
func TestBlockedProbeTimeoutKillsDescendantGroup(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("(trap '' TERM; sleep 30) & echo $! > %q; wait", pidPath)
	status, err := RunBlockedProbe([]byte(probe), 1)
	if err != nil || status != BlockedProbeTimeoutStatus {
		t.Fatalf("status=%d err=%v", status, err)
	}
	pidBytes, readError := os.ReadFile(pidPath)
	if readError != nil {
		t.Fatal(readError)
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("descendant %d survived timeout", pid)
}

func TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("sleep 30 & echo $! > %q; exit 0", pidPath)
	status, err := RunBlockedProbe([]byte(probe), 2)
	if err != nil || status != 0 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	pidBytes, readError := os.ReadFile(pidPath)
	if readError != nil {
		t.Fatal(readError)
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("background descendant %d survived successful probe", pid)
}
