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

// Fast-stage evidence lets the routine gate skip one expensive stage whose
// complete relevant inputs are unchanged. It reuses this package's coverage
// rules, disposition and reason constants, evidence-store mechanics, four-hour
// ceiling and bounded toolchain probe, and differs from the heavy lane in the
// two ways the fast gate demands:
//
//   - The seal is taken over WORKING-TREE BYTES, not committed object ids. The
//     heavy lane may seal committed objects only because it refuses a dirty tree
//     first; the fast gate exists to run on a dirty tree, so committed objects
//     would be a false green the moment anything is uncommitted.
//   - Records live in their own key space with their own schema version. A fast
//     green is computed on a possibly-dirty tree and is not attributable to a
//     revision, so it must never be readable as heavy-lane evidence, which is.
const (
	fastStageEvidenceSchemaVersion = 1
	fastStageManifestSchemaVersion = 1
	fastStageEvidenceDirectoryName = "do-work-fast-stages"
)

// fastStageManifest declares the gate stages that may reuse evidence, the trees
// no stage covers, and the paths no stage seals. Everything the manifest
// classifies nowhere is sealed into EVERY stage, so an unclassified change
// forces all of them.
//
// SealExclusions is the only rule that beats a stage's own coverage: a path it
// names is sealed nowhere, even into a stage whose coverage matches it. Admit a
// path only when BOTH halves hold — the gate or the orchestrator writes it WHILE
// a gate runs, and no stage reads its BYTES, NOT ITS EXISTENCE. That condition
// is the contract, not the entries: the list in _dev/tests/fast-stages.json is
// only today's set of paths satisfying it, and a new churn directory earns a
// place by passing the same test rather than by resembling one already listed.
// The first half is what makes the concept necessary at all: the fingerprint the
// gate records is the PRE-run one, so a seal over a file the run itself writes
// can never match the evidence it authorized, and that stage would execute
// forever.
//
// "Bytes, not existence" is a real limit of the seal, not careful wording. The
// board stats every repo-relative path mentioned in any REQ or UR body
// (skills/do-work-board/tools/queue-kanban/filementions.go), and those mentions
// reach do-work/runs and do-work/deliverables, so creating or deleting a
// mentioned file inside an excluded subtree changes the board JSON without
// moving any seal. No fast-stage assertion reads that map today, which is the
// only reason today's entries are safe; a path whose existence a stage asserts
// on cannot be excluded until that read is gone.
//
// Three of today's entries pass the second half only because
// skills/do-work-board/tools/queue-kanban/walk.go's isSkippedSection prunes
// runs, deliverables and dotted directories from the board walk — which is what
// keeps do-work/runs, do-work/deliverables and do-work/.req-reservations unread.
// That is one set enumerated in two modules, so read isSkippedSection before
// adding an entry here, and re-check these three whenever it changes.
type fastStageManifest struct {
	SchemaVersion    int                 `json:"schema_version"`
	Stages           []manifestFastStage `json:"stages"`
	NonStageCoverage []coverageRule      `json:"non_stage_coverage"`
	SealExclusions   []coverageRule      `json:"seal_exclusions"`
}

type manifestFastStage struct {
	ID       string         `json:"id"`
	Argv     []string       `json:"argv"`
	Coverage []coverageRule `json:"coverage"`
	// Fingerprint is required in practice: a stage that declares none can never
	// be reused, because an undeclared toolchain is an input this command cannot
	// determine. It stays a pointer so the omission is a decodable manifest that
	// always executes rather than an unreadable one.
	Fingerprint *laneFingerprintInputs `json:"fingerprint,omitempty"`
}

// readFastStageManifest reads the manifest from the WORKING TREE. The heavy
// lane reads its manifest at a revision because it runs at that revision; a
// fast stage runs against the files on disk, including uncommitted ones, so the
// manifest that decides its coverage has to be the one on disk too.
func readFastStageManifest(repositoryRoot, manifestPath string) (fastStageManifest, error) {
	manifestAbsolutePath, _, err := resolveManifestPath(repositoryRoot, manifestPath)
	if err != nil {
		return fastStageManifest{}, err
	}
	contents, err := os.ReadFile(manifestAbsolutePath)
	if err != nil {
		return fastStageManifest{}, fmt.Errorf("read fast-stage manifest: %w", err)
	}
	return decodeFastStageManifest(contents)
}

func decodeFastStageManifest(contents []byte) (fastStageManifest, error) {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	var manifest fastStageManifest
	if err := decoder.Decode(&manifest); err != nil {
		return fastStageManifest{}, fmt.Errorf("decode fast-stage manifest: %w", err)
	}
	var trailingValue any
	trailingError := decoder.Decode(&trailingValue)
	if trailingError == nil {
		return fastStageManifest{}, fmt.Errorf("decode fast-stage manifest: trailing JSON value")
	}
	if trailingError != io.EOF {
		return fastStageManifest{}, fmt.Errorf("decode fast-stage manifest: %w", trailingError)
	}
	if manifest.SchemaVersion != fastStageManifestSchemaVersion {
		return fastStageManifest{}, fmt.Errorf("fast-stage manifest schema_version must be %d", fastStageManifestSchemaVersion)
	}
	if len(manifest.Stages) == 0 {
		return fastStageManifest{}, fmt.Errorf("fast-stage manifest must declare at least one stage")
	}
	seenStageIDs := make(map[string]bool)
	for stageIndex, stage := range manifest.Stages {
		if strings.TrimSpace(stage.ID) == "" {
			return fastStageManifest{}, fmt.Errorf("stage %d has an empty id", stageIndex)
		}
		if seenStageIDs[stage.ID] {
			return fastStageManifest{}, fmt.Errorf("stage id %q is duplicated", stage.ID)
		}
		seenStageIDs[stage.ID] = true
		if len(stage.Argv) == 0 {
			return fastStageManifest{}, fmt.Errorf("stage %q has empty argv", stage.ID)
		}
		for _, argument := range stage.Argv {
			if argument == "" {
				return fastStageManifest{}, fmt.Errorf("stage %q argv contains an empty token", stage.ID)
			}
		}
		if len(stage.Coverage) == 0 {
			return fastStageManifest{}, fmt.Errorf("stage %q has no coverage rules", stage.ID)
		}
		for _, rule := range stage.Coverage {
			if err := rule.validate(); err != nil {
				return fastStageManifest{}, fmt.Errorf("stage %q coverage: %w", stage.ID, err)
			}
		}
		if stage.Fingerprint != nil {
			if err := stage.Fingerprint.validate(); err != nil {
				return fastStageManifest{}, fmt.Errorf("stage %q: %w", stage.ID, err)
			}
		}
	}
	for _, rule := range manifest.NonStageCoverage {
		if err := rule.validate(); err != nil {
			return fastStageManifest{}, fmt.Errorf("non-stage coverage: %w", err)
		}
	}
	for _, rule := range manifest.SealExclusions {
		if err := rule.validate(); err != nil {
			return fastStageManifest{}, fmt.Errorf("seal exclusion: %w", err)
		}
	}
	return manifest, nil
}

func selectFastStage(manifest fastStageManifest, stageID string) (manifestFastStage, error) {
	for _, stage := range manifest.Stages {
		if stage.ID == stageID {
			return stage, nil
		}
	}
	return manifestFastStage{}, fmt.Errorf("fast-stage manifest declares no stage %q", stageID)
}

// fastStageSealExcludesPath answers the question that outranks stage coverage:
// may this path be sealed into any stage at all? See fastStageManifest for the
// condition an entry has to satisfy.
func fastStageSealExcludesPath(manifest fastStageManifest, path string) bool {
	return laneCoversPath(manifest.SealExclusions, path)
}

func fastStageManifestClassifiesPath(manifest fastStageManifest, path string) bool {
	if laneCoversPath(manifest.NonStageCoverage, path) {
		return true
	}
	for _, stage := range manifest.Stages {
		if laneCoversPath(stage.Coverage, path) {
			return true
		}
	}
	return false
}

// workingTreeSeals hashes the bytes a stage would actually read right now, for
// every tracked and untracked path the stage covers, plus every path the
// manifest classifies nowhere, MINUS every path SealExclusions names — an
// excluded path is skipped even when the stage's own coverage matches it. An
// input that cannot be determined — a tracked path missing from the worktree, a
// symlink or other non-regular file where a file is expected, an unreadable
// file — is an error, never a weaker seal.
func fastStageCoverageRoots(rules []coverageRule) []string {
	var roots []string
	seen := map[string]bool{}
	for _, rule := range rules {
		var path string
		switch rule.Kind {
		case "exact", "subtree":
			path = rule.Path
		case "suffix-under":
			path = rule.Root
		default:
			return nil
		}
		if path == "" || path == "." {
			return nil
		}
		cleanPath := filepath.Clean(filepath.FromSlash(path))
		if cleanPath == "." || strings.HasPrefix(cleanPath, "..") {
			return nil
		}
		if !seen[path] {
			seen[path] = true
			roots = append(roots, path)
		}
	}
	sort.Strings(roots)
	return roots
}

// workingTreeSeals hashes the bytes a stage would actually read right now, for
// every tracked and untracked path the stage covers, plus every path the
// manifest classifies nowhere, MINUS every path SealExclusions names — an
// excluded path is skipped even when the stage's own coverage matches it. An
// input that cannot be determined — a tracked path missing from the worktree, a
// symlink or other non-regular file where a file is expected, an unreadable
// file — is an error, never a weaker seal.
func workingTreeSeals(repositoryRoot string, stage manifestFastStage, manifest fastStageManifest) ([]string, bool, error) {
	seals := make([]string, 0, 2048)
	sealedPaths := make(map[string]bool, 2048)
	hasCoveredInput := false
	sealPath := func(path string) error {
		if sealedPaths[path] {
			return nil
		}
		sealedPaths[path] = true
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		info, err := os.Lstat(absolutePath)
		if err != nil {
			return fmt.Errorf("inspect stage input %s: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("stage input is not a regular file: %s", path)
		}
		contents, err := os.ReadFile(absolutePath)
		if err != nil {
			return fmt.Errorf("read stage input %s: %w", path, err)
		}
		digest := sha256.Sum256(contents)
		seals = append(seals, fmt.Sprintf("worktree %o %x %s", info.Mode().Perm(), digest, path))
		return nil
	}

	worktreeOutput, err := runGit(repositoryRoot, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	if err != nil {
		return nil, false, fmt.Errorf("inspect worktree stage inputs: %w", err)
	}
	for _, path := range strings.Split(string(worktreeOutput), "\x00") {
		// A seal-excluded path is nobody's input: the gate or the orchestrator
		// writes it while a gate runs, so sealing it would compare a stage
		// against bytes its own run produced. This test comes BEFORE the
		// coverage test below, which is what lets an exclusion beat coverage.
		if path == "" || fastStageSealExcludesPath(manifest, path) {
			continue
		}
		stageCovered := laneCoversPath(stage.Coverage, path)
		hasCoveredInput = hasCoveredInput || stageCovered
		// An unclassified input has to force every stage. Sealing it into each
		// of them is what makes that conservative selection survive reuse.
		if stageCovered || !fastStageManifestClassifiesPath(manifest, path) {
			if err := sealPath(path); err != nil {
				return nil, false, err
			}
		}
	}

	// Go and shell stages read untracked files, Git-ignored build outputs
	// included, so their bytes are sealed as well. An ignored file the manifest
	// classifies nowhere is skipped: refusing all of them would disable reuse on
	// any normal checkout, which carries ignored generated artifacts. Being
	// ignored is no protection once a stage covers the path, so an ignored file
	// the gate writes itself needs a seal exclusion, same as a tracked one.
	ignoredArgv := []string{"ls-files", "-z", "--others", "--exclude-standard", "--ignored"}
	if coverageRoots := fastStageCoverageRoots(stage.Coverage); len(coverageRoots) > 0 {
		ignoredArgv = append(ignoredArgv, "--")
		ignoredArgv = append(ignoredArgv, coverageRoots...)
	}
	ignoredOutput, err := runGit(repositoryRoot, ignoredArgv...)
	if err != nil {
		return nil, false, fmt.Errorf("inspect ignored stage inputs: %w", err)
	}
	for _, path := range strings.Split(string(ignoredOutput), "\x00") {
		if path == "" || fastStageSealExcludesPath(manifest, path) {
			continue
		}
		stageCovered := laneCoversPath(stage.Coverage, path)
		if !stageCovered {
			continue
		}
		hasCoveredInput = true
		if err := sealPath(path); err != nil {
			return nil, false, err
		}
	}
	sort.Strings(seals)
	return seals, hasCoveredInput, nil
}

// fastStageFingerprintDocument is the exact preimage the stage digest is taken
// over. It is a struct rather than a map so the encoding is byte-stable.
type fastStageFingerprintDocument struct {
	SchemaVersion int                  `json:"schema_version"`
	StageID       string               `json:"stage_id"`
	CommandArgv   []string             `json:"command_argv"`
	SealedFiles   []string             `json:"sealed_files"`
	Toolchain     []toolchainProbeSeal `json:"toolchain"`
	Environment   []environmentSeal    `json:"environment"`
}

// fastStageFingerprint digests the stage's command, the working-tree bytes of
// every path it could read, the output of every declared toolchain probe, and
// the whole inherited environment. Any input that cannot be determined is an
// error, which the decision below turns into execution.
func fastStageFingerprint(repositoryRoot string, stage manifestFastStage, manifest fastStageManifest) (string, error) {
	if stage.Fingerprint == nil {
		return "", fmt.Errorf("stage %s declares no fingerprint inputs", stage.ID)
	}
	sealedFiles, hasCoveredInput, err := workingTreeSeals(repositoryRoot, stage, manifest)
	if err != nil {
		return "", err
	}
	if !hasCoveredInput {
		// A stage whose declared coverage matches nothing on disk has no inputs
		// to compare, so no record can be trusted against it.
		return "", fmt.Errorf("stage %s covers no path in this working tree", stage.ID)
	}
	document := fastStageFingerprintDocument{
		SchemaVersion: fastStageEvidenceSchemaVersion, StageID: stage.ID,
		CommandArgv: append([]string(nil), stage.Argv...), SealedFiles: sealedFiles,
		Toolchain: []toolchainProbeSeal{}, Environment: []environmentSeal{},
	}
	for _, probe := range stage.Fingerprint.ToolchainProbes {
		output, err := runFingerprintProbe(repositoryRoot, probe, 5*time.Second)
		if err != nil {
			return "", fmt.Errorf("stage %s toolchain probe %s: %v", stage.ID, strings.Join(probe, " "), err)
		}
		outputDigest := sha256.Sum256(output)
		document.Toolchain = append(document.Toolchain, toolchainProbeSeal{
			ProbeArgv: append([]string(nil), probe...), OutputSHA256: hex.EncodeToString(outputDigest[:]),
		})
	}
	document.Environment = sealedEnvironmentVariables(stage.Fingerprint.EnvironmentVariables)
	encoded, err := json.Marshal(document)
	if err != nil {
		return "", fmt.Errorf("encode stage %s fingerprint: %w", stage.ID, err)
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

// sealedEnvironmentVariables seals every name in the inherited environment, not
// a declared subset: a hand-maintained subset cannot prove that an omitted
// selector or tool setting stayed unchanged. Declared names are unioned in so a
// variable that decides the stage is sealed as absent when it is unset.
func sealedEnvironmentVariables(declaredNames []string) []environmentSeal {
	environmentSet := map[string]bool{}
	for _, name := range declaredNames {
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
	seals := make([]environmentSeal, 0, len(environmentNames))
	for _, variableName := range environmentNames {
		value, isSet := os.LookupEnv(variableName)
		seal := environmentSeal{VariableName: variableName, Set: isSet}
		if isSet {
			valueDigest := sha256.Sum256([]byte(value))
			seal.ValueSHA256 = hex.EncodeToString(valueDigest[:])
		}
		seals = append(seals, seal)
	}
	return seals
}

// storedFastStageEvidence is one successful stage result held in Git-private
// state. Only a stage that exited zero is ever stored, so a record that says
// anything else is corrupt rather than a negative result to trust. It carries
// no revision on purpose: a fast green is computed on a possibly-dirty tree and
// is attributable to that tree, which the fingerprint seals, not to a commit.
type storedFastStageEvidence struct {
	SchemaVersion      int      `json:"schema_version"`
	RepositoryIdentity string   `json:"repository_identity"`
	WorkingTreeRoot    string   `json:"working_tree_root"`
	StageID            string   `json:"stage_id"`
	CommandArgv        []string `json:"command_argv"`
	FingerprintSHA256  string   `json:"fingerprint_sha256"`
	ExitStatus         int      `json:"exit_status"`
	RecordedAt         string   `json:"recorded_at"`
}

// fastStageEvidenceStore is the private per-working-tree record directory. It
// lives in the Git common directory beside the heavy-lane and green-gate
// records, under its own name, and each record is keyed by stage id AND working
// tree: linked worktrees share one common directory but are different trees, and
// a green measured in one of them is not evidence about another.
type fastStageEvidenceStore struct {
	directoryPath      string
	repositoryIdentity string
	workingTreeRoot    string
}

func openFastStageEvidenceStore(repositoryRoot string) (*fastStageEvidenceStore, error) {
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
	resolvedWorkingTree, err := filepath.EvalSymlinks(repositoryRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve working tree identity: %w", err)
	}
	return &fastStageEvidenceStore{
		directoryPath:      filepath.Join(resolvedDirectory, fastStageEvidenceDirectoryName),
		repositoryIdentity: filepath.Clean(resolvedDirectory),
		workingTreeRoot:    filepath.Clean(resolvedWorkingTree),
	}, nil
}

func (store *fastStageEvidenceStore) recordPath(stageID string) string {
	digest := sha256.Sum256([]byte(stageID + "\x00" + store.workingTreeRoot))
	return filepath.Join(store.directoryPath, hex.EncodeToString(digest[:])+".json")
}

// ReadFastStageEvidence returns the stored record for a stage, whether one
// exists, and any reason the stored bytes could not be trusted.
func (store *fastStageEvidenceStore) ReadFastStageEvidence(stageID string) (storedFastStageEvidence, bool, error) {
	recordPath := store.recordPath(stageID)
	info, err := os.Lstat(recordPath)
	if os.IsNotExist(err) {
		return storedFastStageEvidence{}, false, nil
	}
	if err != nil {
		return storedFastStageEvidence{}, false, fmt.Errorf("inspect fast-stage evidence: %w", err)
	}
	if !info.Mode().IsRegular() {
		return storedFastStageEvidence{}, false, fmt.Errorf("fast-stage evidence is not a regular file: %s", recordPath)
	}
	if info.Mode().Perm() != 0o600 {
		return storedFastStageEvidence{}, false, fmt.Errorf("fast-stage evidence does not have private 0600 permissions: %s", recordPath)
	}
	file, err := os.Open(recordPath)
	if err != nil {
		return storedFastStageEvidence{}, false, fmt.Errorf("open fast-stage evidence: %w", err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		return storedFastStageEvidence{}, false, fmt.Errorf("fast-stage evidence identity changed before read")
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		return storedFastStageEvidence{}, false, fmt.Errorf("read fast-stage evidence: %w", err)
	}
	var record storedFastStageEvidence
	if err := json.Unmarshal(contents, &record); err != nil {
		return storedFastStageEvidence{}, false, fmt.Errorf("decode fast-stage evidence: %w", err)
	}
	return record, true, nil
}

func (store *fastStageEvidenceStore) WriteFastStageEvidence(record storedFastStageEvidence) error {
	contents, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("encode fast-stage evidence: %w", err)
	}
	contents = append(contents, '\n')
	if err := ensurePrivateEvidenceDirectory(store.directoryPath); err != nil {
		return err
	}
	recordPath := store.recordPath(record.StageID)
	info, statError := os.Lstat(recordPath)
	switch {
	case statError == nil:
		if !info.Mode().IsRegular() {
			return fmt.Errorf("fast-stage evidence target is not a regular file: %s", recordPath)
		}
		return atomicfile.ReplaceExisting(recordPath, contents)
	case os.IsNotExist(statError):
		return atomicfile.CreateExclusive(recordPath, contents, 0o600)
	default:
		return fmt.Errorf("inspect fast-stage evidence target: %w", statError)
	}
}

// InvalidateFastStageEvidence removes a previous success before any new attempt.
// A failed, interrupted, or unlaunchable stage must never leave an older green
// reusable, so a caller that cannot establish the removal refuses to execute.
func (store *fastStageEvidenceStore) InvalidateFastStageEvidence(stageID string) error {
	if err := ensurePrivateEvidenceDirectory(store.directoryPath); err != nil {
		return err
	}
	recordPath := store.recordPath(stageID)
	info, err := os.Lstat(recordPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("fast-stage evidence target is not a regular file: %s", recordPath)
	}
	return os.Remove(recordPath)
}

// fastStageReuseDecision is what the gate does with one stage and why.
type fastStageReuseDecision struct {
	Disposition string
	Reason      string
	Record      storedFastStageEvidence
}

// decideFastStageReuse answers whether this stage's stored evidence still
// stands. Expiry and fingerprint equality are checked independently, so either
// alone forces a rerun; the reported reason names the first condition that
// failed. Every path except the last returns "executed", which is what makes an
// unknown impact, incomplete evidence, or an unverifiable input select the
// broader verification instead of a false green.
func decideFastStageReuse(store *fastStageEvidenceStore, stage manifestFastStage, fingerprint string, fingerprintError error, evaluatedAt time.Time) fastStageReuseDecision {
	rerun := func(reason string) fastStageReuseDecision {
		return fastStageReuseDecision{Disposition: LaneDispositionExecuted, Reason: reason}
	}
	if fingerprintError != nil {
		return rerun(laneReasonFingerprintUncertain)
	}
	if store == nil {
		return rerun(laneReasonEvidenceUnavailable)
	}
	record, exists, err := store.ReadFastStageEvidence(stage.ID)
	if err != nil {
		return rerun(laneReasonEvidenceUnusable)
	}
	if !exists {
		return rerun(laneReasonNoPriorEvidence)
	}
	if record.SchemaVersion != fastStageEvidenceSchemaVersion || record.RepositoryIdentity != store.repositoryIdentity ||
		record.WorkingTreeRoot != store.workingTreeRoot || record.StageID != stage.ID ||
		!equalLaneArgv(record.CommandArgv, stage.Argv) || record.ExitStatus != 0 || record.FingerprintSHA256 == "" {
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
	return fastStageReuseDecision{Disposition: LaneDispositionReused, Reason: laneReasonFingerprintMatch, Record: record}
}

// FastStageDecisionRequest is one stage's reuse question. SuppliedArgv is the
// command the caller is about to run: the gate and the manifest must agree on
// what a stage IS, and a caller that would run something else is refused rather
// than answered.
type FastStageDecisionRequest struct {
	RepositoryRoot string
	ManifestPath   string
	StageID        string
	SuppliedArgv   []string
	// EvaluatedAt freezes the clock for deterministic tests. Production leaves
	// it zero so each stage decides against its own current instant.
	EvaluatedAt time.Time
}

// DecideFastStage reports whether the named stage may be reported from stored
// evidence instead of executed. The reported line is the gate's whole input, so
// it carries the reason and the fingerprint the caller must supply when it
// later records a green.
func DecideFastStage(request FastStageDecisionRequest) (string, error) {
	stage, manifest, err := resolveFastStage(request.RepositoryRoot, request.ManifestPath, request.StageID, request.SuppliedArgv)
	if err != nil {
		return "", err
	}
	evaluatedAt := request.EvaluatedAt
	if evaluatedAt.IsZero() {
		evaluatedAt = time.Now()
	}
	fingerprint, fingerprintError := fastStageFingerprint(request.RepositoryRoot, stage, manifest)
	// An unreachable store is a reason to execute, never a reason to refuse the
	// question: the gate still has a stage to run.
	store, storeError := openFastStageEvidenceStore(request.RepositoryRoot)
	if storeError != nil {
		store = nil
	}
	decision := decideFastStageReuse(store, stage, fingerprint, fingerprintError, evaluatedAt)
	return formatFastStageDecision(decision, fingerprint), nil
}

// formatFastStageDecision renders the one line the gate parses:
// "<disposition> <reason> <fingerprint-or-dash> <recorded-at-or-dash>". Every
// field is a single token, so a shell reads it with `read` and nothing else.
func formatFastStageDecision(decision fastStageReuseDecision, fingerprint string) string {
	fingerprintField := fingerprint
	if fingerprintField == "" {
		fingerprintField = "-"
	}
	recordedAtField := decision.Record.RecordedAt
	if recordedAtField == "" {
		recordedAtField = "-"
	}
	return fmt.Sprintf("%s %s %s %s\n", decision.Disposition, decision.Reason, fingerprintField, recordedAtField)
}

// FastStageRecordRequest is one green stage result offered for storage.
type FastStageRecordRequest struct {
	RepositoryRoot string
	ManifestPath   string
	StageID        string
	SuppliedArgv   []string
	// SuppliedFingerprint is the digest the decision reported before the stage
	// ran. Recording recomputes it and refuses on a difference, which catches a
	// stage that modified its own inputs while running and would otherwise
	// record a green against a tree that no longer exists.
	SuppliedFingerprint string
	StageExitStatus     int
	RecordedAt          time.Time
}

// RecordFastStage stores a stage's success. A non-zero status writes nothing: a
// skipped, failed, or interrupted stage supplies no evidence.
func RecordFastStage(request FastStageRecordRequest) error {
	if request.StageExitStatus != 0 {
		return fmt.Errorf("stage %s exited %d; only a zero exit records evidence", request.StageID, request.StageExitStatus)
	}
	stage, manifest, err := resolveFastStage(request.RepositoryRoot, request.ManifestPath, request.StageID, request.SuppliedArgv)
	if err != nil {
		return err
	}
	recomputedFingerprint, err := fastStageFingerprint(request.RepositoryRoot, stage, manifest)
	if err != nil {
		return fmt.Errorf("recompute stage %s fingerprint before recording: %w", stage.ID, err)
	}
	if recomputedFingerprint != request.SuppliedFingerprint {
		return fmt.Errorf("stage %s changed its own inputs while running; recorded nothing", stage.ID)
	}
	store, err := openFastStageEvidenceStore(request.RepositoryRoot)
	if err != nil {
		return err
	}
	recordedAt := request.RecordedAt
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}
	return store.WriteFastStageEvidence(storedFastStageEvidence{
		SchemaVersion: fastStageEvidenceSchemaVersion, RepositoryIdentity: store.repositoryIdentity,
		WorkingTreeRoot: store.workingTreeRoot, StageID: stage.ID,
		CommandArgv: append([]string(nil), stage.Argv...), FingerprintSHA256: recomputedFingerprint,
		ExitStatus: 0, RecordedAt: recordedAt.UTC().Format(time.RFC3339),
	})
}

// InvalidateFastStage revokes a stage's prior success before it is attempted.
func InvalidateFastStage(repositoryRoot, stageID string) error {
	if strings.TrimSpace(stageID) == "" {
		return fmt.Errorf("stage id must not be empty")
	}
	store, err := openFastStageEvidenceStore(repositoryRoot)
	if err != nil {
		return err
	}
	return store.InvalidateFastStageEvidence(stageID)
}

// resolveFastStage reads the manifest, finds the stage, and refuses when the
// caller's own command differs from the one the manifest declares. Without that
// check the gate could execute one command and seal another.
func resolveFastStage(repositoryRoot, manifestPath, stageID string, suppliedArgv []string) (manifestFastStage, fastStageManifest, error) {
	manifest, err := readFastStageManifest(repositoryRoot, manifestPath)
	if err != nil {
		return manifestFastStage{}, fastStageManifest{}, err
	}
	stage, err := selectFastStage(manifest, stageID)
	if err != nil {
		return manifestFastStage{}, fastStageManifest{}, err
	}
	if !equalLaneArgv(suppliedArgv, stage.Argv) {
		return manifestFastStage{}, fastStageManifest{}, fmt.Errorf(
			"stage %s is declared as [%s] but the caller supplied [%s]",
			stage.ID, strings.Join(stage.Argv, " "), strings.Join(suppliedArgv, " "))
	}
	return stage, manifest, nil
}
