package lifecycleadvance

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const checkpointRepositoryPath = "do-work/CHECKPOINT.md"

func handleAdvanceCheckpoint(executionContext commandruntime.ExecutionContext) resultmodel.CommandResult {
	snapshot, discoveryError := discoverAdvanceRepository(executionContext.RepositoryRoot)
	if discoveryError != nil {
		return advanceFailure("ADVANCE-CHECKPOINT-DISCOVERY", discoveryError.Error())
	}
	writtenAt := requestmodel.CanonicalTimestamp(time.Now().UTC().Truncate(time.Second))
	checkpointBytes := checkpointSessionBytes(snapshot.CheckpointBytes, writtenAt, checkpointQueueState(snapshot))
	targetOptions := checkpointTransactionOptions(executionContext.RepositoryRoot, snapshot.CheckpointExists)
	transaction := gittransaction.ExecuteTransaction(context.Background(), targetOptions, func(recorder *gittransaction.MutationRecorder) error {
		absolutePath := filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(checkpointRepositoryPath))
		if snapshot.CheckpointExists {
			currentBytes, readError := os.ReadFile(absolutePath)
			if readError != nil || !bytes.Equal(currentBytes, snapshot.CheckpointBytes) {
				return fmt.Errorf("checkpoint changed after discovery")
			}
			if replaceError := atomicfile.ReplaceExisting(absolutePath, checkpointBytes); replaceError != nil {
				return replaceError
			}
			return recorder.RecordTouched(checkpointRepositoryPath)
		}
		if makeError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); makeError != nil {
			return makeError
		}
		if createError := atomicfile.CreateExclusive(absolutePath, checkpointBytes, 0o644); createError != nil {
			return createError
		}
		return recorder.RecordCreated(checkpointRepositoryPath)
	})
	result := gittransaction.BuildCommandResult(CommandAdvance, transaction)
	result.RepositoryRoot = executionContext.RepositoryRoot
	result.Checkpoint = &resultmodel.CheckpointResult{
		CheckpointPath: checkpointRepositoryPath, PreservedClaims: checkpointClaimCount(snapshot), WrittenAt: writtenAt,
	}
	if transaction.Failure == nil {
		result.Outcome = resultmodel.OutcomeSuccess
		result.Findings = append(result.Findings, resultmodel.CommandFinding{
			Code: "ADVANCE-CHECKPOINT-WRITTEN", Severity: resultmodel.SeverityInfo,
			AffectedPaths: []string{checkpointRepositoryPath}, Evidence: []string{"session checkpoint refreshed from one repository snapshot"},
			Fixability:       resultmodel.FixabilityAutomatic,
			NextArgv:         []string{"do-work-cli", "--format", "json", "next"},
			VerificationArgv: []string{"git", "diff", "--check", "--", checkpointRepositoryPath},
		})
	}
	return result
}

func checkpointTransactionOptions(repositoryRoot string, checkpointExists bool) gittransaction.TransactionOptions {
	options := gittransaction.TransactionOptions{RepositoryRoot: repositoryRoot, TargetPaths: []string{checkpointRepositoryPath}}
	if !checkpointExists {
		return options
	}
	tracked := exec.Command("git", "-C", repositoryRoot, "ls-files", "--error-unmatch", "--", checkpointRepositoryPath).Run() == nil
	if !tracked {
		options.ExistingUntrackedTargetPaths = []string{checkpointRepositoryPath}
		return options
	}
	if exec.Command("git", "-C", repositoryRoot, "diff", "--quiet", "--", checkpointRepositoryPath).Run() != nil {
		options.ExistingDirtyTargetPaths = []string{checkpointRepositoryPath}
	}
	return options
}

func checkpointSessionBytes(existing []byte, writtenAt, queueState string) []byte {
	body := append([]byte(nil), existing...)
	frontmatter := ""
	if bytes.HasPrefix(body, []byte("---\n")) {
		if end := bytes.Index(body[4:], []byte("\n---\n")); end >= 0 {
			end += 4
			frontmatter = string(body[4:end])
			body = body[end+5:]
		}
	}
	frontmatter = setCheckpointScalar(frontmatter, "session_ended", writtenAt)
	frontmatter = setCheckpointScalar(frontmatter, "queue_state", queueState)
	if len(bytes.TrimSpace(body)) == 0 {
		body = []byte("# Session Checkpoint\n\n## In Progress (interrupted)\n")
	}
	return []byte("---\n" + strings.TrimSpace(frontmatter) + "\n---\n\n" + strings.TrimLeft(string(body), "\n"))
}

func setCheckpointScalar(frontmatter, name, value string) string {
	lines := strings.Split(strings.TrimSpace(frontmatter), "\n")
	prefix := name + ":"
	for lineIndex, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[lineIndex] = prefix + " " + value
			return strings.Join(lines, "\n")
		}
	}
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	return strings.Join(append(lines, prefix+" "+value), "\n")
}

func checkpointQueueState(snapshot *repositorymodel.RepositorySnapshot) string {
	counts := map[string]int{}
	for _, request := range snapshot.RequestFiles {
		if request.TreeSection == "working" {
			counts["in-progress"]++
			continue
		}
		if request.TreeSection == "queue" {
			counts[request.TypedRecord.RequestStatus]++
		}
	}
	return fmt.Sprintf("[%d pending, %d pending-answers, %d pending-heavy-testing, %d blocked, %d blocked-archive-collision, %d blocked-dependency-cycle, %d in-progress]",
		counts["pending"], counts["pending-answers"], counts["pending-heavy-testing"], counts["blocked"],
		counts["blocked-archive-collision"], counts["blocked-dependency-cycle"], counts["in-progress"])
}

func checkpointClaimCount(snapshot *repositorymodel.RepositorySnapshot) int {
	claimCount := 0
	for _, claims := range snapshot.CheckpointClaimsByID {
		claimCount += len(claims)
	}
	return claimCount
}
