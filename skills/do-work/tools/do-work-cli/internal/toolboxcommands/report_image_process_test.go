//go:build !windows

package toolboxcommands

import (
	"context"
	"testing"
	"time"
)

func TestOwnedProcessCancellationTerminatesAndReaps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	oldGrace := reportImageGracePeriod
	reportImageGracePeriod = 100 * time.Millisecond
	t.Cleanup(func() { reportImageGracePeriod = oldGrace })
	done := make(chan ownedProcessResult, 1)
	go func() { done <- runOwnedProcess(ctx, "", "sh", "-c", "trap '' TERM; sleep 30 & wait") }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case result := <-done:
		if !result.Interrupted {
			t.Fatalf("result=%+v", result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("owned process tree was not terminated and reaped")
	}
}
