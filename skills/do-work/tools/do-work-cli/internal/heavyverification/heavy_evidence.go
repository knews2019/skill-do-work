package heavyverification

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
)

const (
	laneEvidenceSchemaVersion = 2
	laneEvidenceDirectoryName = "do-work-heavy-lanes"
	// laneEvidenceMaximumAge is a ceiling, never a guarantee. It is checked
	// independently of the fingerprint, so either condition alone forces a
	// rerun, and a reuse never refreshes the record that authorized it.
	laneEvidenceMaximumAge = 4 * time.Hour
)

// Lane dispositions. Every lane in a run carries exactly one of these so a
// reader can tell which greens were measured now and which were inherited.
const (
	LaneDispositionExecuted = "executed"
	LaneDispositionReused   = "reused"
)

// Reasons an executed lane was not reused, and the single reason a lane was.
// A caller reports the code; it never has to infer the decision from silence.
const (
	laneReasonFingerprintMatch     = "fingerprint_match"
	laneReasonReuseDisabled        = "reuse_disabled"
	laneReasonFingerprintUncertain = "fingerprint_uncertain"
	laneReasonEvidenceUnavailable  = "evidence_store_unavailable"
	laneReasonNoPriorEvidence      = "no_prior_evidence"
	laneReasonEvidenceUnusable     = "evidence_unusable"
	laneReasonFingerprintMismatch  = "fingerprint_mismatch"
	laneReasonEvidenceExpired      = "evidence_expired"
)

// storedLaneEvidence is one successful lane result held in Git-private state.
// Only a green, unskipped lane is ever stored, so a record that says anything
// else is corrupt rather than a negative result to trust.
type storedLaneEvidence struct {
	SchemaVersion      int      `json:"schema_version"`
	RepositoryIdentity string   `json:"repository_identity"`
	LaneID             string   `json:"lane_id"`
	CommandArgv        []string `json:"command_argv"`
	FingerprintSHA256  string   `json:"fingerprint_sha256"`
	ExitStatus         int      `json:"exit_status"`
	Skipped            bool     `json:"skipped"`
	WallSeconds        int      `json:"wall_seconds"`
	ExecutionRevision  string   `json:"execution_revision"`
	RecordedAt         string   `json:"recorded_at"`
}

// laneEvidenceStore is the private per-repository directory holding lane
// records. It lives beside the green-gate records in the Git common directory,
// so linked worktrees of one repository share the same evidence.
type laneEvidenceStore struct {
	directoryPath      string
	repositoryIdentity string
}

func openLaneEvidenceStore(repositoryRoot string) (*laneEvidenceStore, error) {
	commonDirectoryBytes, err := runGit(repositoryRoot, "rev-parse", "--git-common-dir")
	if err != nil {
		return nil, fmt.Errorf("resolve Git common directory: %w", err)
	}
	commonDirectory := strings.TrimSpace(string(commonDirectoryBytes))
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(repositoryRoot, commonDirectory)
	}
	commonDirectory, err = filepath.Abs(commonDirectory)
	if err != nil {
		return nil, fmt.Errorf("canonicalize Git common directory: %w", err)
	}
	resolvedDirectory, err := filepath.EvalSymlinks(commonDirectory)
	if err != nil {
		return nil, fmt.Errorf("resolve Git common directory identity: %w", err)
	}
	return &laneEvidenceStore{
		directoryPath:      filepath.Join(resolvedDirectory, laneEvidenceDirectoryName),
		repositoryIdentity: filepath.Clean(resolvedDirectory),
	}, nil
}

func (store *laneEvidenceStore) recordPath(laneID string) string {
	digest := sha256.Sum256([]byte(laneID))
	return filepath.Join(store.directoryPath, hex.EncodeToString(digest[:])+".json")
}

// ReadLaneEvidence returns the stored record for a lane, whether one exists,
// and any reason the stored bytes could not be trusted.
func (store *laneEvidenceStore) ReadLaneEvidence(laneID string) (storedLaneEvidence, bool, error) {
	recordPath := store.recordPath(laneID)
	info, err := os.Lstat(recordPath)
	if os.IsNotExist(err) {
		return storedLaneEvidence{}, false, nil
	}
	if err != nil {
		return storedLaneEvidence{}, false, fmt.Errorf("inspect heavy-lane evidence: %w", err)
	}
	if !info.Mode().IsRegular() {
		return storedLaneEvidence{}, false, fmt.Errorf("heavy-lane evidence is not a regular file: %s", recordPath)
	}
	if info.Mode().Perm() != 0o600 {
		return storedLaneEvidence{}, false, fmt.Errorf("heavy-lane evidence does not have private 0600 permissions: %s", recordPath)
	}
	file, err := os.Open(recordPath)
	if err != nil {
		return storedLaneEvidence{}, false, fmt.Errorf("open heavy-lane evidence: %w", err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		return storedLaneEvidence{}, false, fmt.Errorf("heavy-lane evidence identity changed before read")
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		return storedLaneEvidence{}, false, fmt.Errorf("read heavy-lane evidence: %w", err)
	}
	var record storedLaneEvidence
	if err := json.Unmarshal(contents, &record); err != nil {
		return storedLaneEvidence{}, false, fmt.Errorf("decode heavy-lane evidence: %w", err)
	}
	return record, true, nil
}

// WriteLaneEvidence replaces this lane's record with the result just measured.
// Recording is separate from reuse: a run with reuse disabled still refreshes
// the record it deliberately ignored.
func (store *laneEvidenceStore) WriteLaneEvidence(record storedLaneEvidence) error {
	contents, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("encode heavy-lane evidence: %w", err)
	}
	contents = append(contents, '\n')
	if err := ensurePrivateEvidenceDirectory(store.directoryPath); err != nil {
		return err
	}
	recordPath := store.recordPath(record.LaneID)
	info, statError := os.Lstat(recordPath)
	switch {
	case statError == nil:
		if !info.Mode().IsRegular() {
			return fmt.Errorf("heavy-lane evidence target is not a regular file: %s", recordPath)
		}
		return atomicfile.ReplaceExisting(recordPath, contents)
	case os.IsNotExist(statError):
		return atomicfile.CreateExclusive(recordPath, contents, 0o600)
	default:
		return fmt.Errorf("inspect heavy-lane evidence target: %w", statError)
	}
}

func ensurePrivateEvidenceDirectory(directoryPath string) error {
	info, err := os.Lstat(directoryPath)
	if os.IsNotExist(err) {
		if createError := os.Mkdir(directoryPath, 0o700); createError != nil {
			if !os.IsExist(createError) {
				return fmt.Errorf("create heavy-lane evidence directory: %w", createError)
			}
			info, err = os.Lstat(directoryPath)
		} else {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("inspect heavy-lane evidence directory: %w", err)
	}
	if !info.Mode().IsDir() {
		return fmt.Errorf("heavy-lane evidence directory is not a real directory: %s", directoryPath)
	}
	if info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("heavy-lane evidence directory is not private: %s", directoryPath)
	}
	return nil
}

// treeEntry is one path the execution revision commits, with the exact object
// the lane would read at that path.
type treeEntry struct {
	Mode       string
	ObjectType string
	ObjectID   string
	Path       string
}

// readCommittedTree lists every path the revision commits. The runner already
// refuses a tracked tree that differs from HEAD outside the queue directory,
// so committed object ids describe the bytes a lane will actually read.
func readCommittedTree(repositoryRoot, revision string) ([]treeEntry, error) {
	output, err := runGit(repositoryRoot, "ls-tree", "-r", "-z", "--full-tree", revision)
	if err != nil {
		return nil, fmt.Errorf("list the committed tree at %s: %w", revision, err)
	}
	entries := []treeEntry{}
	for _, record := range bytes.Split(output, []byte{0}) {
		if len(record) == 0 {
			continue
		}
		prefix, path, found := bytes.Cut(record, []byte{'\t'})
		fields := bytes.Fields(prefix)
		if !found || len(fields) != 3 || len(path) == 0 {
			return nil, fmt.Errorf("list the committed tree at %s: malformed entry %q", revision, record)
		}
		entries = append(entries, treeEntry{
			Mode: string(fields[0]), ObjectType: string(fields[1]),
			ObjectID: string(fields[2]), Path: string(path),
		})
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("list the committed tree at %s: the revision commits no paths", revision)
	}
	return entries, nil
}

// laneFingerprintInputs declares what, beyond the lane's own covered files,
// decides that lane's result. A lane that declares none is never reusable:
// an undeclared toolchain is an input this command cannot determine, and the
// contract fails closed rather than treating a recent green as authorization.
type laneFingerprintInputs struct {
	ToolchainProbes      [][]string `json:"toolchain_probes"`
	EnvironmentVariables []string   `json:"environment_variables"`
}

func (inputs laneFingerprintInputs) validate() error {
	if len(inputs.ToolchainProbes) == 0 {
		return fmt.Errorf("fingerprint requires at least one toolchain probe")
	}
	for _, probe := range inputs.ToolchainProbes {
		if len(probe) == 0 {
			return fmt.Errorf("fingerprint toolchain probe has empty argv")
		}
		for _, argument := range probe {
			if argument == "" {
				return fmt.Errorf("fingerprint toolchain probe argv contains an empty token")
			}
		}
	}
	seenVariableNames := map[string]bool{}
	for _, variableName := range inputs.EnvironmentVariables {
		if strings.TrimSpace(variableName) == "" {
			return fmt.Errorf("fingerprint environment variable name is empty")
		}
		if seenVariableNames[variableName] {
			return fmt.Errorf("fingerprint environment variable %q is declared twice", variableName)
		}
		seenVariableNames[variableName] = true
	}
	return nil
}

// fingerprintDocument is the exact preimage the lane digest is taken over. It
// is a struct rather than a map so the encoding is byte-stable across runs.
type fingerprintDocument struct {
	SchemaVersion int                  `json:"schema_version"`
	LaneID        string               `json:"lane_id"`
	CommandArgv   []string             `json:"command_argv"`
	CoveredFiles  []string             `json:"covered_files"`
	Toolchain     []toolchainProbeSeal `json:"toolchain"`
	Environment   []environmentSeal    `json:"environment"`
}

type toolchainProbeSeal struct {
	ProbeArgv    []string `json:"probe_argv"`
	OutputSHA256 string   `json:"output_sha256"`
}

type environmentSeal struct {
	VariableName string `json:"variable_name"`
	Set          bool   `json:"set"`
	ValueSHA256  string `json:"value_sha256,omitempty"`
}

// laneFingerprint digests the lane's command, the exact committed bytes of
// every path its manifest coverage declares, the output of every declared
// toolchain probe, and the whole inherited environment. Any input that
// cannot be determined is an error, never a weaker fingerprint.
func laneFingerprint(repositoryRoot string, lane manifestLane, manifest laneManifest, committedTree []treeEntry) (string, error) {
	if lane.Fingerprint == nil {
		return "", fmt.Errorf("lane %s declares no fingerprint inputs", lane.ID)
	}
	document := fingerprintDocument{
		SchemaVersion: laneEvidenceSchemaVersion, LaneID: lane.ID,
		CommandArgv:  append([]string(nil), lane.Argv...),
		CoveredFiles: []string{}, Toolchain: []toolchainProbeSeal{}, Environment: []environmentSeal{},
	}
	hasCoveredInput := false
	for _, entry := range committedTree {
		laneCovered := laneCoversPath(lane.Coverage, entry.Path)
		hasCoveredInput = hasCoveredInput || laneCovered
		// An unclassified input forces all lanes in the planner. Include it
		// here too, or reuse could silently undo that conservative selection.
		if laneCovered || !manifestClassifiesPath(manifest, entry.Path) {
			if entry.Mode != "100644" && entry.Mode != "100755" {
				return "", fmt.Errorf("lane %s input is a symlink, submodule, or unsupported tree mode", lane.ID)
			}
			document.CoveredFiles = append(document.CoveredFiles, fmt.Sprintf("%s %s %s %s", entry.Mode, entry.ObjectType, entry.ObjectID, entry.Path))
		}
	}
	if !hasCoveredInput {
		// A lane whose declared coverage matches nothing at this revision has
		// no test inputs to compare, so no record can be trusted against it.
		return "", fmt.Errorf("lane %s covers no committed path at this revision", lane.ID)
	}
	untrackedFiles, err := fingerprintUntrackedFiles(repositoryRoot, lane, manifest)
	if err != nil {
		return "", err
	}
	document.CoveredFiles = append(document.CoveredFiles, untrackedFiles...)
	sort.Strings(document.CoveredFiles)
	for _, probe := range lane.Fingerprint.ToolchainProbes {
		output, err := runFingerprintProbe(repositoryRoot, probe, 5*time.Second)
		if err != nil {
			return "", fmt.Errorf("lane %s toolchain probe %s: %v", lane.ID, strings.Join(probe, " "), err)
		}
		outputDigest := sha256.Sum256(output)
		document.Toolchain = append(document.Toolchain, toolchainProbeSeal{
			ProbeArgv: append([]string(nil), probe...), OutputSHA256: hex.EncodeToString(outputDigest[:]),
		})
	}
	// Lanes inherit the whole process environment. A hand-maintained subset
	// cannot prove that an omitted selector or tool setting stayed unchanged.
	environmentSet := map[string]bool{}
	for _, name := range lane.Fingerprint.EnvironmentVariables {
		environmentSet[name] = true
	}
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		environmentSet[name] = true
	}
	environmentNames := make([]string, 0, len(environmentSet))
	for name := range environmentSet {
		environmentNames = append(environmentNames, name)
	}
	sort.Strings(environmentNames)
	for _, variableName := range environmentNames {
		value, isSet := os.LookupEnv(variableName)
		seal := environmentSeal{VariableName: variableName, Set: isSet}
		if isSet {
			// The value is digested rather than stored: an environment value
			// can carry any bytes, and only equality decides reuse.
			valueDigest := sha256.Sum256([]byte(value))
			seal.ValueSHA256 = hex.EncodeToString(valueDigest[:])
		}
		document.Environment = append(document.Environment, seal)
	}
	encoded, err := json.Marshal(document)
	if err != nil {
		return "", fmt.Errorf("encode lane %s fingerprint: %w", lane.ID, err)
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

// laneReuseDecision is what the runner does with one lane and why.
type laneReuseDecision struct {
	Disposition string
	Reason      string
	Record      storedLaneEvidence
}

// decideLaneReuse answers whether this lane's stored evidence still stands.
// Expiry and fingerprint equality are checked independently, so either alone
// forces a rerun; the reported reason names the first condition that failed.
func decideLaneReuse(store *laneEvidenceStore, lane manifestLane, fingerprint string, fingerprintError error, evaluatedAt time.Time) laneReuseDecision {
	rerun := func(reason string) laneReuseDecision {
		return laneReuseDecision{Disposition: LaneDispositionExecuted, Reason: reason}
	}
	if fingerprintError != nil {
		return rerun(laneReasonFingerprintUncertain)
	}
	if store == nil {
		return rerun(laneReasonEvidenceUnavailable)
	}
	record, exists, err := store.ReadLaneEvidence(lane.ID)
	if err != nil {
		return rerun(laneReasonEvidenceUnusable)
	}
	if !exists {
		return rerun(laneReasonNoPriorEvidence)
	}
	if record.SchemaVersion != laneEvidenceSchemaVersion || record.RepositoryIdentity != store.repositoryIdentity ||
		record.LaneID != lane.ID || !equalLaneArgv(record.CommandArgv, lane.Argv) ||
		record.ExitStatus != 0 || record.Skipped || record.FingerprintSHA256 == "" || record.ExecutionRevision == "" {
		return rerun(laneReasonEvidenceUnusable)
	}
	recordedAt, err := time.Parse(time.RFC3339, record.RecordedAt)
	if err != nil || recordedAt.After(evaluatedAt) {
		// A record stamped in the future cannot be aged against this clock.
		return rerun(laneReasonEvidenceUnusable)
	}
	if record.FingerprintSHA256 != fingerprint {
		return rerun(laneReasonFingerprintMismatch)
	}
	if evaluatedAt.Sub(recordedAt) > laneEvidenceMaximumAge {
		return rerun(laneReasonEvidenceExpired)
	}
	return laneReuseDecision{Disposition: LaneDispositionReused, Reason: laneReasonFingerprintMatch, Record: record}
}

func equalLaneArgv(left, right []string) bool {
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

// InvalidateLaneEvidence removes a previous success before launching any new
// attempt. A failed, skipped, interrupted, or unlaunchable run must never leave
// an older green reusable. Refuse execution if that cannot be established.
func (store *laneEvidenceStore) InvalidateLaneEvidence(laneID string) error {
	if err := ensurePrivateEvidenceDirectory(store.directoryPath); err != nil {
		return err
	}
	recordPath := store.recordPath(laneID)
	info, err := os.Lstat(recordPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("evidence target is not a regular file: %s", recordPath)
	}
	return os.Remove(recordPath)
}

func laneCoversPath(rules []coverageRule, path string) bool {
	for _, rule := range rules {
		if rule.matches(path) {
			return true
		}
	}
	return false
}

func manifestClassifiesPath(manifest laneManifest, path string) bool {
	if laneCoversPath(manifest.NonHeavyCoverage, path) {
		return true
	}
	for _, lane := range manifest.Lanes {
		if laneCoversPath(lane.Coverage, path) {
			return true
		}
	}
	return false
}

// Go and shell consumers can read untracked files, even Git-ignored sources.
// Seal their bytes as well as committed objects. Generated launcher binaries
// then remain reusable while stable; ignoring them or refusing all of them
// would respectively weaken evidence or disable caching on normal installs.
func fingerprintUntrackedFiles(repositoryRoot string, lane manifestLane, manifest laneManifest) ([]string, error) {
	seals := []string{}
	for _, ignored := range []bool{false, true} {
		argv := []string{"ls-files", "-z", "--others", "--exclude-standard"}
		if ignored {
			argv = append(argv, "--ignored")
		}
		output, err := runGit(repositoryRoot, argv...)
		if err != nil {
			return nil, fmt.Errorf("inspect untracked lane inputs: %w", err)
		}
		for _, path := range strings.Split(string(output), "\x00") {
			if path == "" || strings.HasPrefix(path, queueStatePrefix) {
				continue
			}
			if !laneCoversPath(lane.Coverage, path) && (ignored || manifestClassifiesPath(manifest, path)) {
				continue
			}
			absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
			info, err := os.Lstat(absolutePath)
			if err != nil {
				return nil, fmt.Errorf("inspect untracked lane input: %w", err)
			}
			if !info.Mode().IsRegular() {
				return nil, fmt.Errorf("untracked lane input is not a regular file")
			}
			contents, err := os.ReadFile(absolutePath)
			if err != nil {
				return nil, fmt.Errorf("read untracked lane input: %w", err)
			}
			digest := sha256.Sum256(contents)
			seals = append(seals, fmt.Sprintf("untracked %o %x %s", info.Mode().Perm(), digest, path))
		}
	}
	return seals, nil
}
