package lifecycleadvance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExplicitQueueContinuationSurvivesFirstTargetLeavingQueue(t *testing.T) {
	for _, archiveFirst := range []bool{false, true} {
		t.Run(map[bool]string{false: "working", true: "archived"}[archiveFirst], func(t *testing.T) {
			repositoryRoot := newAdvanceQueueRepository(t)
			for _, requestID := range []string{"REQ-801", "REQ-802"} {
				writeAdvanceRequest(t, repositoryRoot, "queue", requestID, "pending", "", "")
			}
			commitAdvanceQueueFixture(t, repositoryRoot)
			firstResult, firstStatus := runAdvanceQueueJSON(t, repositoryRoot, "REQ-801", "REQ-802")
			if firstStatus != 0 || firstResult.QueueAdvance == nil || len(firstResult.QueueAdvance.Claimed) != 1 || firstResult.QueueAdvance.Claimed[0].RequestID != "REQ-801" {
				t.Fatalf("first claim: status=%d result=%#v", firstStatus, firstResult)
			}
			if archiveFirst {
				completeClaimedFixture(t, repositoryRoot, "REQ-801")
			}
			writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-803", "pending", "", "")
			runAdvanceGit(t, repositoryRoot, "add", "do-work")
			runAdvanceGit(t, repositoryRoot, "commit", "-qm", "later queue state")
			secondResult, secondStatus := runAdvanceContinuationJSON(t, repositoryRoot, firstResult.QueueAdvance.ContinuationArgv)
			if secondStatus != 0 || secondResult.QueueAdvance == nil || len(secondResult.QueueAdvance.Claimed) != 1 || secondResult.QueueAdvance.Claimed[0].RequestID != "REQ-802" || len(secondResult.QueueAdvance.ContinuationArgv) != 0 {
				t.Fatalf("continuation: status=%d result=%#v", secondStatus, secondResult)
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-803-fixture.md")); err != nil {
				t.Fatalf("later arrival was consumed: %v", err)
			}
			replayResult, replayStatus := runAdvanceContinuationJSON(t, repositoryRoot, firstResult.QueueAdvance.ContinuationArgv)
			if replayResult.QueueAdvance != nil && len(replayResult.QueueAdvance.Claimed) != 0 {
				t.Fatalf("replay claimed a member twice: status=%d result=%#v", replayStatus, replayResult)
			}
		})
	}
}
