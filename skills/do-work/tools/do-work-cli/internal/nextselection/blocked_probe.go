package nextselection

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	BlockedProbeTimeoutStatus   = 124
	BlockedProbeLaunchStatus    = 125
	blockedProbeDiagnosticLimit = 64 * 1024
)

type BlockedProbeEvidence struct {
	ExitStatus       int
	Launched         bool
	TimedOut         bool
	Diagnostic       string
	DiagnosticSHA256 string
}

type boundedProbeWriter struct {
	mutex     sync.Mutex
	bytes     []byte
	truncated bool
}

func (writer *boundedProbeWriter) Write(input []byte) (int, error) {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()
	remaining := blockedProbeDiagnosticLimit - len(writer.bytes)
	if remaining > 0 {
		if len(input) < remaining {
			remaining = len(input)
		}
		writer.bytes = append(writer.bytes, input[:remaining]...)
	}
	if remaining < len(input) {
		writer.truncated = true
	}
	return len(input), nil
}

func (writer *boundedProbeWriter) String() string {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()
	diagnostic := string(writer.bytes)
	if writer.truncated {
		diagnostic += "\n[diagnostic truncated]"
	}
	return diagnostic
}

// BlockedProbeInterruption distinguishes an invocation signal from a probe's
// own non-zero status so queue selection can stop instead of excluding one REQ.
type BlockedProbeInterruption struct {
	ExitStatus int
}

func (interruption BlockedProbeInterruption) Error() string {
	return fmt.Sprintf("blocked probe interrupted with exit %d", interruption.ExitStatus)
}

func (interruption BlockedProbeInterruption) InterruptionExitStatus() int {
	return interruption.ExitStatus
}

// RunBlockedProbe executes one materialized probe while the platform-specific
// implementation owns its complete descendant tree. Raw legacy status is returned
// as evidence; public command rendering remains in the standard 0-4 envelope.
func RunBlockedProbe(probeBytes []byte, timeoutSeconds int) (int, error) {
	repositoryRoot, err := os.Getwd()
	if err != nil {
		return BlockedProbeLaunchStatus, err
	}
	return RunBlockedProbeAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
}

// RunBlockedProbeAtRoot executes one materialized probe relative to the selected repository.
func RunBlockedProbeAtRoot(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error) {
	evidence, err := RunBlockedProbeEvidenceAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
	return evidence.ExitStatus, err
}

// RunBlockedProbeEvidenceAtRoot retains bounded diagnostics while preserving
// the owned process-tree and raw-status behavior used by queue selection.
func RunBlockedProbeEvidenceAtRoot(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (BlockedProbeEvidence, error) {
	if repositoryRoot == "" {
		return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus}, fmt.Errorf("repository root is empty")
	}
	if timeoutSeconds <= 0 {
		return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus}, fmt.Errorf("timeout must be positive")
	}
	if len(probeBytes) == 0 {
		return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus}, fmt.Errorf("probe is empty")
	}
	diagnosticWriter := &boundedProbeWriter{}
	status, err := runOwnedProbe(repositoryRoot, probeBytes, time.Duration(timeoutSeconds)*time.Second, diagnosticWriter)
	diagnostic, diagnosticSHA256 := BlockedProbeDiagnosticIdentity(diagnosticWriter.String(), repositoryRoot)
	return BlockedProbeEvidence{
		ExitStatus: status, Launched: status != BlockedProbeLaunchStatus,
		TimedOut: status == BlockedProbeTimeoutStatus, Diagnostic: diagnostic,
		DiagnosticSHA256: diagnosticSHA256,
	}, err
}

// BlockedProbeDiagnosticIdentity applies the same bounded-probe normalization
// to saved baseline bytes before their identity is compared.
func BlockedProbeDiagnosticIdentity(diagnostic, repositoryRoot string) (string, string) {
	const truncatedMarker = "\n[diagnostic truncated]"
	if len(diagnostic) > blockedProbeDiagnosticLimit && !strings.HasSuffix(diagnostic, truncatedMarker) {
		diagnostic = diagnostic[:blockedProbeDiagnosticLimit] + truncatedMarker
	}
	diagnostic = strings.ReplaceAll(diagnostic, "\r\n", "\n")
	diagnostic = strings.ReplaceAll(diagnostic, repositoryRoot, "<repo-root>")
	lines := strings.Split(strings.TrimSpace(diagnostic), "\n")
	for index := range lines {
		lines[index] = strings.TrimSpace(lines[index])
	}
	normalized := strings.Join(lines, "\n")
	digest := sha256.Sum256([]byte(normalized))
	return normalized, hex.EncodeToString(digest[:])
}

func blockedProbeInterruptionStatus(err error) (int, bool) {
	var interruption interface {
		error
		InterruptionExitStatus() int
	}
	if !errors.As(err, &interruption) {
		return 0, false
	}
	status := interruption.InterruptionExitStatus()
	return status, status == 129 || status == 130 || status == 143
}
