package gateevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	gateEvidenceSchemaVersion = 1
	gateEvidenceDirectoryName = "do-work-green-gates"
	reportedGateProvenance    = "reported_direct_zero_exit"
	persistedGateProvenance   = "persisted_green_run"
	gateLogPathPrefix         = "_dev/gate-runs/"
)

type storedGateEvidence struct {
	SchemaVersion      int      `json:"schema_version"`
	RepositoryIdentity string   `json:"repository_identity"`
	GateCommand        []string `json:"gate_command"`
	GateCommandSHA256  string   `json:"gate_command_sha256"`
	RecordProvenance   string   `json:"record_provenance"`
	GateExitStatus     int      `json:"gate_exit_status"`
	RecordedRevision   string   `json:"recorded_revision"`
}

type repositoryEvidenceContext struct {
	repositoryIdentity string
	headRevision       string
	recordPath         string
	gateCommandSHA256  string
}

func RecordGreenGate(repositoryRoot string, gateCommand []string) (resultmodel.GateEvidenceResult, error) {
	context, err := resolveEvidenceContext(repositoryRoot, gateCommand)
	if err != nil {
		return invalidEvidence(gateCommand), err
	}
	result := evidenceResult(context, gateCommand)
	result.State = resultmodel.GateEvidenceRecorded
	result.RecordProvenance = reportedGateProvenance
	result.RecordedRevision = context.headRevision
	result.BaselineRevision = context.headRevision
	result.TargetRevision = context.headRevision

	record := storedGateEvidence{
		SchemaVersion: gateEvidenceSchemaVersion, RepositoryIdentity: context.repositoryIdentity,
		GateCommand: append([]string(nil), gateCommand...), GateCommandSHA256: context.gateCommandSHA256,
		RecordProvenance: reportedGateProvenance, GateExitStatus: 0, RecordedRevision: context.headRevision,
	}
	contents, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return invalidEvidence(gateCommand), fmt.Errorf("encode green-gate evidence: %w", err)
	}
	contents = append(contents, '\n')
	if err := publishEvidenceRecord(context.recordPath, record, contents); err != nil {
		result.State = resultmodel.GateEvidenceInvalidRecord
		result.BaselineRevision = ""
		return result, err
	}
	return result, nil
}

// CheckGreenGate answers whether the recorded green run for this exact gate argv
// still covers the current HEAD. This is the HEAD-bound rule every ordinary REQ uses.
func CheckGreenGate(repositoryRoot string, gateCommand []string) (resultmodel.GateEvidenceResult, error) {
	return CheckGreenGateAtRevision(repositoryRoot, gateCommand, "")
}

// CheckGreenGateAtRevision answers the same question about targetRevisionSpec instead
// of HEAD; an empty spec means HEAD. An already-green repository-gate repair no-op
// changes no project path, so an unrelated commit moving HEAD can never make it the
// cause of a red gate — its evidence stays verifiable at its own recorded revision.
func CheckGreenGateAtRevision(repositoryRoot string, gateCommand []string, targetRevisionSpec string) (resultmodel.GateEvidenceResult, error) {
	context, err := resolveEvidenceContext(repositoryRoot, gateCommand)
	if err != nil {
		return invalidEvidence(gateCommand), err
	}
	result := evidenceResult(context, gateCommand)
	targetRevision := context.headRevision
	if targetRevisionSpec != "" {
		targetRevision, err = resolveCommitRevision(repositoryRoot, targetRevisionSpec)
		if err != nil {
			result.State = resultmodel.GateEvidenceInvalidRecord
			return result, err
		}
	}
	result.TargetRevision = targetRevision
	record, exists, err := readEvidenceRecord(context.recordPath)
	if err != nil {
		result.State = resultmodel.GateEvidenceInvalidRecord
		return result, err
	}
	if !exists {
		result.State = resultmodel.GateEvidenceMissing
		return result, nil
	}
	result.RecordProvenance = persistedGateProvenance
	result.GateExitStatus = record.GateExitStatus
	result.RecordedRevision = record.RecordedRevision
	if record.SchemaVersion != gateEvidenceSchemaVersion || record.RecordProvenance != reportedGateProvenance || record.GateExitStatus != 0 || record.RecordedRevision == "" {
		result.State = resultmodel.GateEvidenceInvalidRecord
		return result, fmt.Errorf("green-gate evidence has an unsupported or incomplete record schema")
	}
	if record.RepositoryIdentity != context.repositoryIdentity {
		result.State = resultmodel.GateEvidenceDifferentRepository
		return result, nil
	}
	if record.GateCommandSHA256 != context.gateCommandSHA256 || !equalArgv(record.GateCommand, gateCommand) {
		result.State = resultmodel.GateEvidenceDifferentArgv
		return result, nil
	}

	recordedRevisionExists, err := commitResolvesExactly(repositoryRoot, record.RecordedRevision)
	if err != nil {
		result.State = resultmodel.GateEvidenceInvalidRecord
		return result, err
	}
	if !recordedRevisionExists {
		result.State = resultmodel.GateEvidenceRecordedRevisionMissing
		return result, nil
	}
	if record.RecordedRevision == targetRevision {
		result.State = resultmodel.GateEvidenceExactRevisionMatch
		result.Matches = true
		result.MatchBasis = "exact_revision"
		result.BaselineRevision = targetRevision
		return result, nil
	}

	isAncestor, err := recordedRevisionIsAncestor(repositoryRoot, record.RecordedRevision, targetRevision)
	if err != nil {
		result.State = resultmodel.GateEvidenceInvalidRecord
		return result, err
	}
	if !isAncestor {
		result.State = resultmodel.GateEvidenceRecordedRevisionNotAncestor
		return result, nil
	}
	logOnly, err := interveningCommitsAreGateLogs(repositoryRoot, record.RecordedRevision, targetRevision)
	if err != nil {
		result.State = resultmodel.GateEvidenceInvalidRecord
		return result, err
	}
	if !logOnly {
		result.State = resultmodel.GateEvidenceInvalidated
		return result, nil
	}
	result.State = resultmodel.GateEvidenceLogDescendantMatch
	result.Matches = true
	result.MatchBasis = "gate_log_only_descendant"
	result.BaselineRevision = targetRevision
	return result, nil
}

// checkGreenGateAtRevision preserves the package-private test seam while the
// exported function is shared with repair validation.
func checkGreenGateAtRevision(repositoryRoot string, gateCommand []string, targetRevisionSpec string) (resultmodel.GateEvidenceResult, error) {
	return CheckGreenGateAtRevision(repositoryRoot, gateCommand, targetRevisionSpec)
}

func GateCommandSHA256(gateCommand []string) string {
	encoded, _ := json.Marshal(gateCommand)
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

func resolveEvidenceContext(repositoryRoot string, gateCommand []string) (repositoryEvidenceContext, error) {
	if len(gateCommand) == 0 {
		return repositoryEvidenceContext{}, fmt.Errorf("gate command must contain at least one argv token")
	}
	commonDirectoryBytes, err := runGitOutput(repositoryRoot, "rev-parse", "--git-common-dir")
	if err != nil {
		return repositoryEvidenceContext{}, fmt.Errorf("resolve Git common directory: %w", err)
	}
	commonDirectory := strings.TrimSpace(string(commonDirectoryBytes))
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(repositoryRoot, commonDirectory)
	}
	commonDirectory, err = filepath.Abs(commonDirectory)
	if err != nil {
		return repositoryEvidenceContext{}, fmt.Errorf("canonicalize Git common directory: %w", err)
	}
	if resolvedDirectory, resolveError := filepath.EvalSymlinks(commonDirectory); resolveError == nil {
		commonDirectory = resolvedDirectory
	} else {
		return repositoryEvidenceContext{}, fmt.Errorf("resolve Git common directory identity: %w", resolveError)
	}
	headBytes, err := runGitOutput(repositoryRoot, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return repositoryEvidenceContext{}, fmt.Errorf("resolve HEAD revision: %w", err)
	}
	headRevision := strings.TrimSpace(string(headBytes))
	gateDigest := GateCommandSHA256(gateCommand)
	recordPath := filepath.Join(commonDirectory, gateEvidenceDirectoryName, gateDigest+".json")
	return repositoryEvidenceContext{
		repositoryIdentity: filepath.Clean(commonDirectory), headRevision: headRevision,
		recordPath: recordPath, gateCommandSHA256: gateDigest,
	}, nil
}

func evidenceResult(context repositoryEvidenceContext, gateCommand []string) resultmodel.GateEvidenceResult {
	return resultmodel.GateEvidenceResult{
		RepositoryIdentity: context.repositoryIdentity, GateCommand: append([]string(nil), gateCommand...),
		GateCommandSHA256: context.gateCommandSHA256, RecordPath: context.recordPath,
		GateExitStatus: 0, HeadRevision: context.headRevision, MatchBasis: "none",
	}
}

func invalidEvidence(gateCommand []string) resultmodel.GateEvidenceResult {
	return resultmodel.GateEvidenceResult{
		GateCommand: append([]string(nil), gateCommand...), GateCommandSHA256: GateCommandSHA256(gateCommand),
		State: resultmodel.GateEvidenceInvalidRecord, MatchBasis: "none",
	}
}

func publishEvidenceRecord(recordPath string, expected storedGateEvidence, contents []byte) error {
	recordDirectory := filepath.Dir(recordPath)
	if err := ensurePrivateDirectory(recordDirectory); err != nil {
		return err
	}
	recordInfo, statError := os.Lstat(recordPath)
	switch {
	case statError == nil:
		if !recordInfo.Mode().IsRegular() {
			return fmt.Errorf("green-gate evidence target is not a regular file: %s", recordPath)
		}
		current, exists, readError := readEvidenceRecord(recordPath)
		if readError != nil {
			return fmt.Errorf("validate existing green-gate evidence: %w", readError)
		}
		if !exists {
			return fmt.Errorf("validate existing green-gate evidence: target disappeared before replacement")
		}
		if current.SchemaVersion != gateEvidenceSchemaVersion || current.RecordProvenance != reportedGateProvenance || current.GateExitStatus != 0 || current.RecordedRevision == "" {
			return fmt.Errorf("existing green-gate evidence has an unsupported or incomplete record schema")
		}
		if current.RepositoryIdentity != expected.RepositoryIdentity || current.GateCommandSHA256 != expected.GateCommandSHA256 || !equalArgv(current.GateCommand, expected.GateCommand) {
			return fmt.Errorf("existing green-gate evidence does not match this repository and argv")
		}
		return atomicfile.ReplaceExisting(recordPath, contents)
	case os.IsNotExist(statError):
		return atomicfile.CreateExclusive(recordPath, contents, 0o600)
	default:
		return fmt.Errorf("inspect green-gate evidence target: %w", statError)
	}
}

func ensurePrivateDirectory(directoryPath string) error {
	info, err := os.Lstat(directoryPath)
	if os.IsNotExist(err) {
		if createError := os.Mkdir(directoryPath, 0o700); createError != nil {
			if !os.IsExist(createError) {
				return fmt.Errorf("create green-gate evidence directory: %w", createError)
			}
			info, err = os.Lstat(directoryPath)
		} else {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("inspect green-gate evidence directory: %w", err)
	}
	if !info.Mode().IsDir() {
		return fmt.Errorf("green-gate evidence directory is not a real directory: %s", directoryPath)
	}
	if info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("green-gate evidence directory is not private: %s", directoryPath)
	}
	return nil
}

func readEvidenceRecord(recordPath string) (storedGateEvidence, bool, error) {
	info, err := os.Lstat(recordPath)
	if os.IsNotExist(err) {
		return storedGateEvidence{}, false, nil
	}
	if err != nil {
		return storedGateEvidence{}, false, fmt.Errorf("inspect green-gate evidence: %w", err)
	}
	if !info.Mode().IsRegular() {
		return storedGateEvidence{}, false, fmt.Errorf("green-gate evidence is not a regular file: %s", recordPath)
	}
	if info.Mode().Perm() != 0o600 {
		return storedGateEvidence{}, false, fmt.Errorf("green-gate evidence does not have private 0600 permissions: %s", recordPath)
	}
	file, err := os.Open(recordPath)
	if err != nil {
		return storedGateEvidence{}, false, fmt.Errorf("open green-gate evidence: %w", err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		return storedGateEvidence{}, false, fmt.Errorf("green-gate evidence identity changed before read")
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		return storedGateEvidence{}, false, fmt.Errorf("read green-gate evidence: %w", err)
	}
	afterInfo, err := file.Stat()
	if err != nil || !os.SameFile(openedInfo, afterInfo) || openedInfo.Size() != afterInfo.Size() || !openedInfo.ModTime().Equal(afterInfo.ModTime()) {
		return storedGateEvidence{}, false, fmt.Errorf("green-gate evidence changed during read")
	}
	var record storedGateEvidence
	if err := json.Unmarshal(contents, &record); err != nil {
		return storedGateEvidence{}, false, fmt.Errorf("decode green-gate evidence: %w", err)
	}
	return record, true, nil
}

// resolveCommitRevision turns a caller-supplied revision specification into the exact
// commit it names. An unresolvable specification is unverifiable evidence, not a miss.
func resolveCommitRevision(repositoryRoot, revisionSpec string) (string, error) {
	output, status, err := runGitStatus(repositoryRoot, "rev-parse", "--verify", revisionSpec+"^{commit}")
	if err != nil {
		return "", err
	}
	if status != 0 {
		return "", fmt.Errorf("target revision %q does not name a commit in this repository", revisionSpec)
	}
	return strings.TrimSpace(string(output)), nil
}

func commitResolvesExactly(repositoryRoot, revision string) (bool, error) {
	output, status, err := runGitStatus(repositoryRoot, "rev-parse", "--verify", revision+"^{commit}")
	if err != nil {
		return false, err
	}
	if status != 0 {
		return false, nil
	}
	return strings.TrimSpace(string(output)) == revision, nil
}

func recordedRevisionIsAncestor(repositoryRoot, recordedRevision, headRevision string) (bool, error) {
	_, status, err := runGitStatus(repositoryRoot, "merge-base", "--is-ancestor", recordedRevision, headRevision)
	if err != nil {
		return false, err
	}
	switch status {
	case 0:
		return true, nil
	case 1:
		return false, nil
	default:
		return false, fmt.Errorf("git merge-base --is-ancestor exited %d", status)
	}
}

func interveningCommitsAreGateLogs(repositoryRoot, recordedRevision, headRevision string) (bool, error) {
	output, err := runGitOutput(repositoryRoot, "rev-list", "--reverse", recordedRevision+".."+headRevision)
	if err != nil {
		return false, fmt.Errorf("enumerate revisions after green gate: %w", err)
	}
	for _, commit := range strings.Fields(string(output)) {
		paths, err := runGitOutput(repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", "-m", "-z", commit)
		if err != nil {
			return false, fmt.Errorf("inspect revision %s after green gate: %w", commit, err)
		}
		for _, pathBytes := range bytes.Split(paths, []byte{0}) {
			if len(pathBytes) > 0 && !strings.HasPrefix(string(pathBytes), gateLogPathPrefix) {
				return false, nil
			}
		}
	}
	return true, nil
}

func runGitOutput(repositoryRoot string, arguments ...string) ([]byte, error) {
	output, status, err := runGitStatus(repositoryRoot, arguments...)
	if err != nil {
		return output, err
	}
	if status != 0 {
		return output, fmt.Errorf("git %s exited %d", strings.Join(arguments, " "), status)
	}
	return output, nil
}

func runGitStatus(repositoryRoot string, arguments ...string) ([]byte, int, error) {
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.Output()
	if err == nil {
		return output, 0, nil
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return output, exitError.ExitCode(), nil
	}
	return output, -1, fmt.Errorf("launch git %s: %w", strings.Join(arguments, " "), err)
}

func equalArgv(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
