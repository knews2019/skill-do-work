// Package nextselection selects queue work from one typed repository snapshot.
// It is read-only: claiming, state transitions, and release remain pipeline work.
package nextselection

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

const (
	ProvenanceDefault     = "default"
	ProvenanceExplicit    = "explicit-req"
	ProvenanceUserRequest = "ur-expanded"
	ProvenanceSimple      = "simple-selector"
)

const (
	PriorityRepositoryGateRepair = "repository-gate-repair"
	PriorityDeferredParent       = "gate-deferred"
	PriorityOrdinary             = "ordinary"
)

const (
	RequestPriorityNow   = "now"
	RequestPriorityNext  = "next"
	RequestPriorityLater = "later"
)

type SelectionOptions struct {
	TargetTokens         []string
	WaveDepth            *int
	FanOutLimit          *int
	SkipImpactNegligible bool
	SimpleOnly           bool
}

type selectionCandidate struct {
	RequestFile *repositorymodel.RequestFile
	RequestID   string
	Provenance  string
	SourceToken string
	Priority    string
}

type ProbeRunner func(probeBytes []byte, timeoutSeconds int) (int, error)

func requestID(requestFile *repositorymodel.RequestFile) string {
	if requestFile == nil {
		return ""
	}
	if requestFile.TypedRecord.RequestID != "" {
		return requestFile.TypedRecord.RequestID
	}
	return requestFile.FilenameID
}

func numericID(identifier, prefix string) (int, bool) {
	trimmed := strings.TrimSpace(identifier)
	if len(trimmed) <= len(prefix) || !strings.EqualFold(trimmed[:len(prefix)], prefix) {
		return 0, false
	}
	digits := trimmed[len(prefix):]
	for _, character := range digits {
		if character < '0' || character > '9' {
			return 0, false
		}
	}
	value, err := strconv.Atoi(digits)
	return value, err == nil
}

func canonicalToken(identifier, prefix string) (string, error) {
	number, valid := numericID(identifier, prefix)
	if !valid {
		return "", fmt.Errorf("%s must be %s followed by digits", identifier, prefix)
	}
	return fmt.Sprintf("%s%03d", prefix, number), nil
}
