package finalization

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func decodeManifest(repositoryRoot, manifestPath string) (Manifest, []byte, error) {
	absolutePath, err := containedOrAbsolute(repositoryRoot, manifestPath)
	if err != nil {
		return Manifest{}, nil, err
	}
	info, err := os.Lstat(absolutePath)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return Manifest{}, nil, fmt.Errorf("finalization manifest must be a regular non-symlink file")
	}
	contents, err := os.ReadFile(absolutePath)
	if err != nil {
		return Manifest{}, nil, err
	}
	decoder := json.NewDecoder(strings.NewReader(string(contents)))
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, nil, fmt.Errorf("decode finalization manifest: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return Manifest{}, nil, fmt.Errorf("finalization manifest must contain one JSON object")
	}
	if err := validateManifest(manifest); err != nil {
		return Manifest{}, nil, err
	}
	return manifest, contents, nil
}

func validateManifest(manifest Manifest) error {
	if !requestIDPattern.MatchString(manifest.RequestID) {
		return fmt.Errorf("request_id must be an exact REQ-NNN id")
	}
	if manifest.RequestPath == "" || filepath.IsAbs(manifest.RequestPath) {
		return fmt.Errorf("request_path must be repository-relative")
	}
	if manifest.WriterLabel == "" {
		return fmt.Errorf("writer_label is required")
	}
	if manifest.Transition != "complete" && manifest.Transition != "fail" {
		return fmt.Errorf("transition must be complete or fail")
	}
	if manifest.Transition == "complete" {
		if manifest.TerminalStatus != "completed" && manifest.TerminalStatus != "completed-with-issues" {
			return fmt.Errorf("complete requires terminal_status completed or completed-with-issues")
		}
	} else if strings.TrimSpace(manifest.FailureError) == "" || !canonicalFailureType(manifest.FailureType) {
		return fmt.Errorf("fail requires failure_error and failure_type intent, spec, code, or environment")
	}
	if _, err := time.Parse(time.RFC3339, manifest.CompletedAt); err != nil {
		return fmt.Errorf("completed_at requires RFC3339: %w", err)
	}
	if !digestPattern.MatchString(strings.ToLower(manifest.ExpectedRequestSHA256)) || !digestPattern.MatchString(strings.ToLower(manifest.ExpectedCheckpointSHA256)) {
		return fmt.Errorf("expected request and checkpoint SHA-256 values are required")
	}
	if strings.TrimSpace(manifest.CommitMessage) == "" || len(manifest.CommitPaths) == 0 {
		return fmt.Errorf("commit_message and a non-empty commit_paths allowlist are required")
	}
	if manifest.ReleaseManifestPath == "" && manifest.ReleaseAt != "" || manifest.ReleaseManifestPath != "" && manifest.ReleaseAt == "" {
		return fmt.Errorf("release_manifest_path and release_at must be supplied together")
	}
	if manifest.ReleaseAt != "" {
		if _, err := time.Parse(time.RFC3339, manifest.ReleaseAt); err != nil {
			return fmt.Errorf("release_at requires RFC3339: %w", err)
		}
	}
	return nil
}

func journalLocations(repositoryRoot, requestID string) (string, string, error) {
	output, err := exec.Command("git", "-C", repositoryRoot, "rev-parse", "--git-path", "do-work-finalization").Output()
	if err != nil {
		return "", "", fmt.Errorf("resolve finalization journal directory: %w", err)
	}
	directory := strings.TrimSpace(string(output))
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(repositoryRoot, directory)
	}
	directory = filepath.Clean(directory)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return "", "", err
	}
	return filepath.Join(directory, requestID+".json"), filepath.Join(directory, requestID+".payloads"), nil
}

func writeJournal(journal *Journal) error {
	journal.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	contents, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return err
	}
	contents = append(contents, '\n')
	parent := filepath.Dir(journal.JournalPath)
	temporary, err := os.CreateTemp(parent, ".finalization-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, journal.JournalPath); err != nil {
		return err
	}
	directory, err := os.Open(parent)
	if err == nil {
		_ = directory.Sync()
		_ = directory.Close()
	}
	return nil
}

func readJournal(path string) (*Journal, error) {
	path = filepath.Clean(path)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("journal is not a regular non-symlink file: %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var journal Journal
	if err := decoder.Decode(&journal); err != nil {
		return nil, err
	}
	if journal.Version != journalVersion || filepath.Clean(journal.JournalPath) != path || !requestIDPattern.MatchString(journal.Manifest.RequestID) {
		return nil, fmt.Errorf("journal identity or version is invalid")
	}
	if filepath.Base(path) != journal.Manifest.RequestID+".json" {
		return nil, fmt.Errorf("journal filename does not match request identity")
	}
	if err := validateManifest(journal.Manifest); err != nil {
		return nil, fmt.Errorf("journal manifest is invalid: %w", err)
	}
	expectedPayloadDirectory := filepath.Join(filepath.Dir(path), journal.Manifest.RequestID+".payloads")
	if journal.PayloadDirectory != "" && filepath.Clean(journal.PayloadDirectory) != expectedPayloadDirectory {
		return nil, fmt.Errorf("journal payload directory is outside its private request slot")
	}
	switch journal.Phase {
	case PhasePrepared, PhaseLifecycleApplied, PhaseReleaseApplied, PhasePrimaryCommitted, PhaseMetadataCommitted:
	default:
		return nil, fmt.Errorf("journal phase is invalid: %s", journal.Phase)
	}
	return &journal, nil
}

func listJournals(repositoryRoot string) ([]string, error) {
	probePath, _, err := journalLocations(repositoryRoot, "REQ-0")
	if err != nil {
		return nil, err
	}
	directory := filepath.Dir(probePath)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	type journalEntry struct {
		path      string
		createdAt time.Time
	}
	journals := []journalEntry{}
	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.HasPrefix(entry.Name(), "REQ-") && strings.HasSuffix(entry.Name(), ".json") {
			path := filepath.Join(directory, entry.Name())
			journal, readError := readJournal(path)
			if readError != nil {
				return nil, readError
			}
			journals = append(journals, journalEntry{path: path, createdAt: journal.CreatedAt})
		}
	}
	sort.Slice(journals, func(left, right int) bool {
		if journals[left].createdAt.Equal(journals[right].createdAt) {
			return journals[left].path < journals[right].path
		}
		return journals[left].createdAt.Before(journals[right].createdAt)
	})
	paths := make([]string, 0, len(journals))
	for _, journal := range journals {
		paths = append(paths, journal.path)
	}
	return paths, nil
}

func removeJournal(journal *Journal) error {
	if journal.PayloadDirectory != "" {
		expected := filepath.Join(filepath.Dir(journal.JournalPath), journal.Manifest.RequestID+".payloads")
		if filepath.Clean(journal.PayloadDirectory) != expected {
			return fmt.Errorf("refusing to remove a noncanonical payload directory")
		}
		if err := os.RemoveAll(journal.PayloadDirectory); err != nil {
			return err
		}
	}
	return os.Remove(journal.JournalPath)
}

func digestBytes(contents []byte) string {
	digest := sha256.Sum256(contents)
	return hex.EncodeToString(digest[:])
}

func containedOrAbsolute(repositoryRoot, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes repository: %s", path)
	}
	return filepath.Join(repositoryRoot, cleaned), nil
}

func canonicalFailureType(value string) bool {
	switch value {
	case "intent", "spec", "code", "environment":
		return true
	default:
		return false
	}
}
