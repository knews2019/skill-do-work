package lifecycleadvance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/publication"
)

type advanceFinalizationResult struct {
	Outcome       string                        `json:"outcome"`
	Advance       *advanceCommandEvidence       `json:"advance"`
	Findings      []advanceCommandFinding       `json:"findings"`
	Changes       []advanceRecordedChange       `json:"changes"`
	Rollback      advanceRollbackResult         `json:"rollback"`
	Finalization  *advanceFinalizationEvidence  `json:"finalization"`
	Finalizations []advanceFinalizationEvidence `json:"finalizations"`
}

type advanceRecordedChange struct {
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type advanceRollbackResult struct {
	Status  string   `json:"status"`
	Actions []string `json:"actions"`
	Errors  []string `json:"errors"`
}

type advanceFinalizationEvidence struct {
	RequestID             string   `json:"request_id"`
	RequestPath           string   `json:"request_path"`
	ArchivePath           string   `json:"archive_path"`
	JournalPath           string   `json:"journal_path"`
	Phase                 string   `json:"phase"`
	TerminalStatus        string   `json:"terminal_status"`
	Resumed               bool     `json:"resumed"`
	Discovered            bool     `json:"discovered"`
	CommitPaths           []string `json:"commit_paths"`
	PrimaryCommit         string   `json:"primary_commit"`
	MetadataCommit        string   `json:"metadata_commit"`
	CreatedPrimaryCommit  string   `json:"created_primary_commit"`
	CreatedMetadataCommit string   `json:"created_metadata_commit"`
	BlockedPaths          []string `json:"blocked_paths"`
	ReasonCodes           []string `json:"reason_codes"`
	NextArgv              []string `json:"next_argv"`
	VerificationArgv      []string `json:"verification_argv"`
	CollectionArgv        []string `json:"collection_argv"`
}

type advanceFinalizationManifest struct {
	RequestID                string   `json:"request_id"`
	RequestPath              string   `json:"request_path"`
	WriterLabel              string   `json:"writer_label"`
	Transition               string   `json:"transition"`
	TerminalStatus           string   `json:"terminal_status"`
	CompletedAt              string   `json:"completed_at"`
	ExpectedRequestSHA256    string   `json:"expected_request_sha256"`
	ExpectedCheckpointSHA256 string   `json:"expected_checkpoint_sha256"`
	CommitPaths              []string `json:"commit_paths"`
	CommitMessage            string   `json:"commit_message"`
	ProvenanceMode           string   `json:"provenance_mode"`
	ImplementationHash       string   `json:"implementation_hash,omitempty"`
	ReleaseManifestPath      string   `json:"release_manifest_path,omitempty"`
	ReleaseAt                string   `json:"release_at,omitempty"`
}

func TestAdvanceFinalizationRunsTerminalPathMatrix(t *testing.T) {
	tests := []struct {
		name                   string
		requestID              string
		terminalStatus         string
		provenanceMode         string
		implementationChanged  bool
		wantMetadataCommit     bool
		wantArchivedCommitHash bool
		withRelease            bool
	}{
		{name: "serial", requestID: "REQ-780", terminalStatus: "completed", provenanceMode: "primary_commit", implementationChanged: true, wantMetadataCommit: true, wantArchivedCommitHash: true, withRelease: true},
		{name: "supplied worktree", requestID: "REQ-781", terminalStatus: "completed", provenanceMode: "supplied_commit", implementationChanged: true, wantArchivedCommitHash: true},
		{name: "completed with issues", requestID: "REQ-782", terminalStatus: "completed-with-issues", provenanceMode: "primary_commit", implementationChanged: true, wantMetadataCommit: true, wantArchivedCommitHash: true},
		{name: "already green without release", requestID: "REQ-783", terminalStatus: "completed", provenanceMode: "primary_commit", wantMetadataCommit: true, wantArchivedCommitHash: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot, requestPath, manifestPath, seedCommit := seedAdvanceFinalization(t, test.requestID, test.terminalStatus, test.provenanceMode, test.implementationChanged, test.withRelease)
			result, status, output := runAdvanceFinalization(t, repositoryRoot, "json", test.requestID, "--request-path", requestPath, "--finalization-manifest", manifestPath)
			if status != 0 || result.Outcome != "success" || result.Advance == nil || len(result.Finalizations) != 1 || result.Finalization == nil {
				t.Fatalf("status=%d result=%#v\n%s", status, result, output)
			}
			if result.Advance.RequestID != test.requestID || result.Advance.RequestPath != requestPath || result.Advance.Phase != "finalize" || result.Advance.PhaseKind != "mechanical" {
				t.Fatalf("advance identity/phase = %#v", result.Advance)
			}
			record := result.Finalizations[0]
			if !reflect.DeepEqual(record, *result.Finalization) || record.RequestID != test.requestID || record.RequestPath != requestPath || record.Phase != "cleanup_complete" || record.TerminalStatus != test.terminalStatus {
				t.Fatalf("finalization record = %#v singular=%#v", record, result.Finalization)
			}
			if record.ArchivePath == "" || record.JournalPath == "" || record.PrimaryCommit == "" || record.CreatedPrimaryCommit == "" || len(record.CommitPaths) == 0 || record.BlockedPaths == nil || record.ReasonCodes == nil || record.NextArgv == nil || record.VerificationArgv == nil || record.CollectionArgv == nil {
				t.Fatalf("finalization record lost typed fields: %#v", record)
			}
			if got := record.MetadataCommit != "" && record.CreatedMetadataCommit != ""; got != test.wantMetadataCommit {
				t.Fatalf("metadata commit presence=%t, want %t: %#v", got, test.wantMetadataCommit, record)
			}
			if record.PrimaryCommit != record.CreatedPrimaryCommit || test.wantMetadataCommit && record.MetadataCommit != record.CreatedMetadataCommit {
				t.Fatalf("settled and created hashes differ: %#v", record)
			}
			archiveBytes, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(record.ArchivePath)))
			if err != nil {
				t.Fatal(err)
			}
			archiveText := string(archiveBytes)
			if !strings.Contains(archiveText, "status: "+test.terminalStatus) {
				t.Fatalf("archive status mismatch:\n%s", archiveText)
			}
			if test.wantArchivedCommitHash && !strings.Contains(archiveText, "commit:") {
				t.Fatalf("archive has no commit provenance:\n%s", archiveText)
			}
			if test.provenanceMode == "supplied_commit" && !strings.Contains(archiveText, "commit: "+seedCommit[:12]) {
				t.Fatalf("archive does not carry supplied commit %s:\n%s", seedCommit, archiveText)
			}
			if test.provenanceMode == "supplied_commit" {
				for _, path := range record.CommitPaths {
					if path == "implementation.txt" {
						t.Fatalf("supplied implementation leaked into finalization commit allowlist: %#v", record.CommitPaths)
					}
				}
				implementationAtSuppliedHash := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "show", seedCommit+":implementation.txt")))
				if implementationAtSuppliedHash != "after" {
					t.Fatalf("supplied hash %s does not own implementation bytes: %q", seedCommit, implementationAtSuppliedHash)
				}
			}
			versionBytes, err := os.ReadFile(filepath.Join(repositoryRoot, "VERSION"))
			if err != nil {
				t.Fatal(err)
			}
			changelogBytes, err := os.ReadFile(filepath.Join(repositoryRoot, "CHANGELOG.md"))
			if err != nil {
				t.Fatal(err)
			}
			if test.withRelease {
				if string(versionBytes) != "1.0.1\n" || !strings.Contains(string(changelogBytes), "## 1.0.1 — Advance Finalization") || !strings.Contains(archiveText, "release_at: 2026-09-04T13:05:00Z") {
					t.Fatalf("release payload was not published exactly: version=%q changelog=%q archive=%q", versionBytes, changelogBytes, archiveText)
				}
				for _, wantPath := range []string{"VERSION", "CHANGELOG.md"} {
					found := false
					for _, path := range record.CommitPaths {
						found = found || path == wantPath
					}
					if !found {
						t.Fatalf("release path %q missing from finalization commit allowlist: %#v", wantPath, record.CommitPaths)
					}
				}
			} else if string(versionBytes) != "1.0.0\n" || string(changelogBytes) != "# Changelog\n\n## 1.0.0 — Seed\n" || strings.Contains(archiveText, "release_at:") {
				t.Fatalf("no-release path changed release bytes: version=%q changelog=%q archive=%q", versionBytes, changelogBytes, archiveText)
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))); !os.IsNotExist(err) {
				t.Fatalf("working request remains: %v", err)
			}
			if gitStatus := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=all"))); gitStatus != "" {
				t.Fatalf("finalization left Git dirt: %s", gitStatus)
			}
		})
	}
}

func TestAdvanceFinalizationRequiresOneBoundManifestWithoutMutation(t *testing.T) {
	tests := []struct {
		name           string
		arguments      func(string, string) []string
		mutateManifest func(*advanceFinalizationManifest)
		wantOutcome    string
	}{
		{name: "missing manifest", arguments: func(requestPath, _ string) []string { return []string{"--request-path", requestPath} }, wantOutcome: "success"},
		{name: "duplicate manifest", arguments: func(requestPath, manifestPath string) []string {
			return []string{"--request-path", requestPath, "--finalization-manifest", manifestPath, "--finalization-manifest", manifestPath}
		}, wantOutcome: "refused"},
		{name: "empty manifest", arguments: func(requestPath, _ string) []string {
			return []string{"--request-path", requestPath, "--finalization-manifest="}
		}, wantOutcome: "refused"},
		{name: "hostile token remains data", arguments: func(requestPath, manifestPath string) []string {
			return []string{"--request-path", requestPath, "--finalization-manifest", manifestPath, "$(touch hostile-marker)"}
		}, wantOutcome: "refused"},
		{name: "outer path mismatch", arguments: func(_ string, manifestPath string) []string {
			return []string{"--request-path", "do-work/working/REQ-999-wrong.md", "--finalization-manifest", manifestPath}
		}, wantOutcome: "refused"},
		{name: "manifest id mismatch", arguments: func(requestPath, manifestPath string) []string {
			return []string{"--request-path", requestPath, "--finalization-manifest", manifestPath}
		}, mutateManifest: func(manifest *advanceFinalizationManifest) { manifest.RequestID = "REQ-999" }, wantOutcome: "refused"},
		{name: "manifest path mismatch", arguments: func(requestPath, manifestPath string) []string {
			return []string{"--request-path", requestPath, "--finalization-manifest", manifestPath}
		}, mutateManifest: func(manifest *advanceFinalizationManifest) { manifest.RequestPath = "do-work/working/REQ-784-other.md" }, wantOutcome: "refused"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot, requestPath, manifestPath, _ := seedAdvanceFinalization(t, "REQ-784", "completed", "primary_commit", true, false)
			if test.mutateManifest != nil {
				var manifest advanceFinalizationManifest
				manifestBytes, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
					t.Fatal(err)
				}
				test.mutateManifest(&manifest)
				writeAdvanceJSON(t, manifestPath, manifest)
			}
			beforeTree := advanceTreeDigest(t, repositoryRoot)
			beforeHead := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))
			result, status, output := runAdvanceFinalization(t, repositoryRoot, "json", "REQ-784", test.arguments(requestPath, manifestPath)...)
			if result.Outcome != test.wantOutcome {
				t.Fatalf("status=%d outcome=%q, want %q\n%s", status, result.Outcome, test.wantOutcome, output)
			}
			if test.wantOutcome == "success" {
				if status != 0 || result.Advance == nil || result.Advance.Phase != "finalize" || len(result.Advance.NextArgv) == 0 || !reflect.DeepEqual(result.Advance.NextArgv[len(result.Advance.NextArgv)-2:], []string{"--finalization-manifest", "<action-authored-finalization-manifest>"}) {
					t.Fatalf("missing-input projection = status %d %#v", status, result)
				}
			} else if status == 0 || len(result.Findings) == 0 {
				t.Fatalf("refusal/failure lost typed finding: status=%d result=%#v", status, result)
			}
			if afterTree := advanceTreeDigest(t, repositoryRoot); afterTree != beforeTree {
				t.Fatalf("repository bytes changed: before=%s after=%s", beforeTree, afterTree)
			}
			if afterHead := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD"))); afterHead != beforeHead {
				t.Fatalf("HEAD changed: before=%s after=%s", beforeHead, afterHead)
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, "hostile-marker")); !os.IsNotExist(err) {
				t.Fatalf("hostile token executed: %v", err)
			}
		})
	}
}

func seedAdvanceFinalization(t *testing.T, requestID, terminalStatus, provenanceMode string, implementationChanged, withRelease bool) (string, string, string, string) {
	t.Helper()
	repositoryRoot := t.TempDir()
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Finalization Fixture")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance-finalization@example.invalid")
	requestPath := "do-work/working/" + requestID + "-finalize.md"
	archivePath := "do-work/archive/" + requestID + "-finalize.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	requestStatus := "claimed"
	if terminalStatus == "completed-with-issues" {
		requestStatus = terminalStatus
	}
	requestFrontmatter := "route: C\nplanning_at: 2026-09-04T12:00:00Z\nclaimed_at: 2026-09-04T11:00:00Z\nwrite_set: [implementation.txt]\nestimate:\n  p50_active_minutes: 30\ncommit:\n"
	writeAdvanceRequest(t, repositoryRoot, "working", requestID, requestStatus, requestFrontmatter, routeCBodyThrough("Orientation"))
	defaultPath := "do-work/working/" + requestID + "-fixture.md"
	if defaultPath != requestPath {
		if err := os.Rename(filepath.Join(repositoryRoot, filepath.FromSlash(defaultPath)), filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))); err != nil {
			t.Fatal(err)
		}
	}
	writeAdvanceFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- "+requestID+": fixture — claimed now — writer: host:/repo\n")
	writeAdvanceFile(t, repositoryRoot, "implementation.txt", "before\n")
	writeAdvanceFile(t, repositoryRoot, "VERSION", "1.0.0\n")
	writeAdvanceFile(t, repositoryRoot, "CHANGELOG.md", "# Changelog\n\n## 1.0.0 — Seed\n")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "seed")
	seedCommit := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))
	if implementationChanged {
		writeAdvanceFile(t, repositoryRoot, "implementation.txt", "after\n")
		if provenanceMode == "supplied_commit" {
			runAdvanceGit(t, repositoryRoot, "add", "implementation.txt")
			runAdvanceGit(t, repositoryRoot, "commit", "-qm", "implementation")
			seedCommit = strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))
		}
	}
	manifest := advanceFinalizationManifest{
		RequestID: requestID, RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete", TerminalStatus: terminalStatus,
		CompletedAt: "2026-09-04T13:00:00Z", ExpectedRequestSHA256: advanceFileDigest(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: advanceFileDigest(t, repositoryRoot, checkpointPath),
		CommitPaths: []string{requestPath, archivePath, checkpointPath, "do-work/calibration-log.tsv"}, CommitMessage: "[" + requestID + "] finalize through advance", ProvenanceMode: provenanceMode,
	}
	if implementationChanged && provenanceMode != "supplied_commit" {
		manifest.CommitPaths = append(manifest.CommitPaths, "implementation.txt")
	}
	if provenanceMode == "supplied_commit" {
		manifest.ImplementationHash = seedCommit
	}
	if withRelease {
		payloadRoot := t.TempDir()
		versionOldPath := filepath.Join(payloadRoot, "version-old")
		versionNewPath := filepath.Join(payloadRoot, "version-new")
		changelogOldPath := filepath.Join(payloadRoot, "changelog-old")
		changelogNewPath := filepath.Join(payloadRoot, "changelog-new")
		writeAdvancePayload(t, versionOldPath, "1.0.0\n")
		writeAdvancePayload(t, versionNewPath, "1.0.1\n")
		writeAdvancePayload(t, changelogOldPath, "# Changelog\n\n## 1.0.0 — Seed\n")
		writeAdvancePayload(t, changelogNewPath, "# Changelog\n\n## 1.0.1 — Advance Finalization\n\nThe terminal transaction now runs through advance.\n\n## 1.0.0 — Seed\n")
		releaseManifest := publication.Manifest{Operation: publication.OperationRelease, Release: &publication.ReleaseManifest{
			OldVersion: "1.0.0", NewVersion: "1.0.1", ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
			Targets:    []publication.ReleaseTarget{{Path: "VERSION", ExpectedPayload: publication.PayloadFile{SourcePath: versionOldPath}, NewPayload: publication.PayloadFile{SourcePath: versionNewPath}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
			Changelogs: []publication.ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: publication.PayloadFile{SourcePath: changelogOldPath}, NewPayload: publication.PayloadFile{SourcePath: changelogNewPath}, InsertionAnchor: "## 1.0.0", EntryKey: "1.0.1", EntryTitle: "Advance Finalization"}},
		}}
		releaseManifestPath := filepath.Join(t.TempDir(), "release-manifest.json")
		writeAdvanceJSON(t, releaseManifestPath, releaseManifest)
		manifest.ReleaseManifestPath = releaseManifestPath
		manifest.ReleaseAt = "2026-09-04T13:05:00Z"
		manifest.CommitPaths = append(manifest.CommitPaths, "VERSION", "CHANGELOG.md")
	}
	manifestPath := filepath.Join(t.TempDir(), "finalization-manifest.json")
	writeAdvanceJSON(t, manifestPath, manifest)
	return repositoryRoot, requestPath, manifestPath, seedCommit
}

func writeAdvancePayload(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeAdvanceJSON(t *testing.T, path string, value any) {
	t.Helper()
	contents, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
}

func advanceFileDigest(t *testing.T, repositoryRoot, relativePath string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(contents)
	return hex.EncodeToString(digest[:])
}

func runAdvanceFinalization(t *testing.T, repositoryRoot, format, requestID string, arguments ...string) (advanceFinalizationResult, int, string) {
	t.Helper()
	commandArguments := []string{"--repo-root", repositoryRoot, "--format", format, "advance", requestID}
	commandArguments = append(commandArguments, arguments...)
	command := exec.Command(advanceCLIBinary(t), commandArguments...)
	output, err := command.CombinedOutput()
	status := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
		} else {
			t.Fatalf("advance launch: %v", err)
		}
	}
	var result advanceFinalizationResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode advance finalization result: %v\n%s", err, output)
	}
	return result, status, string(output)
}
