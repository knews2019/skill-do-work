package repositorymodel

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func writeRepositoryFixture(t *testing.T, repositoryRoot string, relativePath string, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func requestFixture(requestID string, status string) string {
	return "---\nid: " + requestID + "\nstatus: " + status + "\n---\n\nBody\n"
}

func TestDiscoverRepositoryCoversLiveArchiveReservationAndExcludedLayouts(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-001-queue.md", requestFixture("REQ-001", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/working/REQ-002-working.md", requestFixture("REQ-002", "claimed"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/REQ-003-legacy.md", requestFixture("REQ-003", "completed"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-080/REQ-004-nested.md", requestFixture("REQ-004", "done"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-001-100/UR-081/REQ-005-banded.md", requestFixture("REQ-005", "completed"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/user-requests/UR-081/input.md", "---\nid: UR-081\nrequests: [REQ-005]\n---\nInput\n")
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-001-100/UR-080/input.md", "---\nid: UR-080\n---\nArchived input\n")
	writeRepositoryFixture(t, repositoryRoot, "do-work/.req-reservations/REQ-000009", "")
	writeRepositoryFixture(t, repositoryRoot, "do-work/.req-reservations/REQ-010", "")
	writeRepositoryFixture(t, repositoryRoot, "do-work/.req-reservations/REQ-999-copy", "")
	writeRepositoryFixture(t, repositoryRoot, "do-work/deliverables/REQ-090-copy.md", requestFixture("REQ-090", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/runs/run/REQ-091-copy.md", requestFixture("REQ-091", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/user-requests/REQ-006-stray.md", requestFixture("REQ-006", "pending"))

	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatalf("DiscoverRepository: %v", err)
	}
	if len(snapshot.RequestFiles) != 5 || len(snapshot.UserRequestFiles) != 2 || len(snapshot.ReservationFiles) != 2 {
		t.Fatalf("snapshot counts: requests=%d URs=%d reservations=%d", len(snapshot.RequestFiles), len(snapshot.UserRequestFiles), len(snapshot.ReservationFiles))
	}
	if snapshot.ReservationFiles[0].RequestNumber != 9 || snapshot.ReservationFiles[1].RequestNumber != 10 {
		t.Fatalf("reservation identities = %#v", snapshot.ReservationFiles)
	}
	if len(snapshot.StrayRequestPaths) != 1 || !strings.HasSuffix(filepath.ToSlash(snapshot.StrayRequestPaths[0]), "user-requests/REQ-006-stray.md") {
		t.Fatalf("stray requests = %v", snapshot.StrayRequestPaths)
	}
	if snapshot.RequestsByID["REQ-004"][0].TypedRecord.RequestStatus != "completed" {
		t.Fatalf("legacy status was not normalized: %#v", snapshot.RequestsByID["REQ-004"][0].TypedRecord)
	}
	root, err := FindRepositoryRoot(filepath.Join(repositoryRoot, "do-work", "archive", "UR-080"))
	if err != nil || root != repositoryRoot {
		t.Fatalf("FindRepositoryRoot = %q, %v; want %q", root, err, repositoryRoot)
	}
}

func TestDiscoverRepositoryRetainsRequestModificationTime(t *testing.T) {
	repositoryRoot := t.TempDir()
	relativePath := "do-work/working/REQ-001-working.md"
	writeRepositoryFixture(t, repositoryRoot, relativePath, requestFixture("REQ-001", "claimed"))
	modifiedAt := time.Date(2026, 8, 30, 19, 15, 12, 0, time.FixedZone("UTC+2", 2*60*60))
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.Chtimes(absolutePath, modifiedAt, modifiedAt); err != nil {
		t.Fatal(err)
	}

	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	requestFile := snapshot.RequestsByID["REQ-001"][0]
	wantModifiedAt := modifiedAt.UTC().Truncate(time.Second)
	if !requestFile.ModifiedAt.Equal(wantModifiedAt) {
		t.Fatalf("modified_at = %s, want %s", requestFile.ModifiedAt, wantModifiedAt)
	}
}

func TestDiscoverRepositoryProjectsCheckpointClaimsInSourceOrder(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-001-first.md", requestFixture("REQ-001", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/CHECKPOINT.md", strings.Join([]string{
		"# Session Checkpoint",
		"",
		"## Arbitrary placement",
		"- REQ-000001: First claim — claimed 2026-09-01T10:00:00Z — writer: host-a:/repo",
		"- REQ-999: Unrelated — claimed earlier — writer: host-z:/repo",
		"- REQ-001: Duplicate claim — claimed 2026-09-01T10:01:00Z — writer: host-b:/repo",
		"- REQ-001: Missing writer — claimed 2026-09-01T10:02:00Z",
		"ordinary prose mentioning REQ-001 — writer: host-c:/repo",
	}, "\n"))

	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	claims := snapshot.CheckpointClaimsByID["REQ-001"]
	if len(claims) != 2 {
		t.Fatalf("REQ-001 checkpoint claims = %#v, want two writer-bearing headers", claims)
	}
	if claims[0].Writer != "host-a:/repo" || claims[0].ClaimedAt != "2026-09-01T10:00:00Z" || claims[0].SourceLine != 4 || claims[0].RelativePath != "CHECKPOINT.md" {
		t.Fatalf("first checkpoint claim = %#v", claims[0])
	}
	if claims[1].Writer != "host-b:/repo" || claims[1].ClaimedAt != "2026-09-01T10:01:00Z" || claims[1].SourceLine != 6 {
		t.Fatalf("second checkpoint claim = %#v", claims[1])
	}
	if claims[0].HeaderText != "- REQ-000001: First claim — claimed 2026-09-01T10:00:00Z — writer: host-a:/repo" {
		t.Fatalf("header text = %q", claims[0].HeaderText)
	}
	if unrelated := snapshot.CheckpointClaimsByID["REQ-999"]; len(unrelated) != 1 || unrelated[0].Writer != "host-z:/repo" {
		t.Fatalf("unrelated checkpoint projection = %#v", unrelated)
	}
}

func TestDiscoverRepositoryRetainsCollisionEvidence(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-010-first.md", requestFixture("REQ-011", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/REQ-011-second.md", requestFixture("REQ-011", "completed"))
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.RequestsByID["REQ-011"]) != 2 || len(snapshot.CollisionEntries) != 1 {
		t.Fatalf("collisions = %#v index=%v", snapshot.CollisionEntries, snapshot.RequestsByID)
	}
	if snapshot.CollisionEntries[0].RequestID != "REQ-011" || len(snapshot.CollisionEntries[0].ClaimPaths) != 2 {
		t.Fatalf("collision evidence = %#v", snapshot.CollisionEntries[0])
	}
}

func TestCollisionEvidenceIncludesFilenameAndFrontmatterClaims(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-020-first.md", requestFixture("REQ-021", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/REQ-021-second.md", requestFixture("REQ-022", "completed"))
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.CollisionEntries) != 1 || snapshot.CollisionEntries[0].RequestID != "REQ-021" || len(snapshot.CollisionEntries[0].ClaimPaths) != 2 {
		t.Fatalf("filename/frontmatter collision evidence = %#v", snapshot.CollisionEntries)
	}
}

func TestCollisionEvidenceNormalizesNumericFrontmatterIdentity(t *testing.T) {
	repositoryRoot := t.TempDir()
	firstPath := "do-work/queue/REQ-900-first.md"
	secondPath := "do-work/queue/REQ-901-second.md"
	writeRepositoryFixture(t, repositoryRoot, firstPath, requestFixture("REQ-452", "pending"))
	writeRepositoryFixture(t, repositoryRoot, secondPath, requestFixture("REQ-0452", "pending"))
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}

	if len(snapshot.CollisionEntries) != 1 || snapshot.CollisionEntries[0].RequestID != "REQ-452" {
		t.Fatalf("normalized collision evidence = %#v", snapshot.CollisionEntries)
	}
	wantPaths := []string{
		filepath.Join(repositoryRoot, filepath.FromSlash(firstPath)),
		filepath.Join(repositoryRoot, filepath.FromSlash(secondPath)),
	}
	claimPaths := snapshot.CollisionEntries[0].ClaimPaths
	if len(claimPaths) != len(wantPaths) || claimPaths[0] != wantPaths[0] || claimPaths[1] != wantPaths[1] {
		t.Fatalf("collision paths = %#v, want %#v", snapshot.CollisionEntries[0].ClaimPaths, wantPaths)
	}
}

func TestReserveNextRequestIDIsCollisionFreeAcrossEvidenceAndRaces(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-008-file.md", requestFixture("REQ-012", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/REQ-010-old.md", requestFixture("REQ-010", "completed"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/.req-reservations/REQ-000011", "")
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	reservations := make(chan ReservationFile, 2)
	errorsSeen := make(chan error, 2)
	var waitGroup sync.WaitGroup
	for workerIndex := 0; workerIndex < 2; workerIndex++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			<-start
			reservation, reserveError := ReserveNextRequestID(snapshot)
			reservations <- reservation
			errorsSeen <- reserveError
		}()
	}
	close(start)
	waitGroup.Wait()
	close(reservations)
	close(errorsSeen)
	for reserveError := range errorsSeen {
		if reserveError != nil {
			t.Fatalf("ReserveNextRequestID: %v", reserveError)
		}
	}
	var numbers []int
	for reservation := range reservations {
		numbers = append(numbers, reservation.RequestNumber)
		if _, err := os.Stat(reservation.AbsolutePath); err != nil {
			t.Errorf("reservation missing: %v", err)
		}
		if basename := filepath.Base(reservation.AbsolutePath); basename != "REQ-013" && basename != "REQ-014" {
			t.Errorf("reservation basename = %q, want REQ-013 or REQ-014", basename)
		}
	}
	sort.Ints(numbers)
	if len(numbers) != 2 || numbers[0] != 13 || numbers[1] != 14 {
		t.Fatalf("reserved numbers = %v, want [13 14]", numbers)
	}
}

func TestReservationDirectorySymlinkIsRefused(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-001.md", requestFixture("REQ-001", "pending"))
	outsideDirectory := t.TempDir()
	reservationPath := filepath.Join(repositoryRoot, "do-work", ".req-reservations")
	if err := os.Symlink(outsideDirectory, reservationPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ReserveNextRequestID(snapshot); err == nil || errors.Is(err, os.ErrExist) {
		t.Fatalf("symlinked reservation directory error = %v", err)
	}
}

func TestReservationDirectorySwapCannotRedirectMarker(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-001.md", requestFixture("REQ-001", "pending"))
	writeRepositoryFixture(t, repositoryRoot, "do-work/.req-reservations/REQ-000001", "")
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	reservationDirectory := filepath.Join(repositoryRoot, "do-work", ".req-reservations")
	originalDirectory := reservationDirectory + ".original"
	outsideDirectory := t.TempDir()
	previousHook := beforeReservationMarkerCreate
	beforeReservationMarkerCreate = func(path string) {
		beforeReservationMarkerCreate = previousHook
		if renameError := os.Rename(path, originalDirectory); renameError != nil {
			t.Fatalf("rename reservation directory: %v", renameError)
		}
		if symlinkError := os.Symlink(outsideDirectory, path); symlinkError != nil {
			t.Fatalf("swap reservation directory for symlink: %v", symlinkError)
		}
	}
	t.Cleanup(func() { beforeReservationMarkerCreate = previousHook })
	if _, reserveError := ReserveNextRequestID(snapshot); reserveError == nil {
		t.Fatal("reservation succeeded after its rooted directory was replaced")
	}
	outsideEntries, err := os.ReadDir(outsideDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if len(outsideEntries) != 0 {
		t.Fatalf("reservation escaped into symlink target: %v", outsideEntries)
	}
	originalEntries, err := os.ReadDir(originalDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if len(originalEntries) != 1 || originalEntries[0].Name() != "REQ-000001" {
		t.Fatalf("failed reservation left a marker in the displaced directory: %v", originalEntries)
	}
}

func TestDiscoveryRejectsSymlinkedRequestAndUserRequestFiles(t *testing.T) {
	repositoryRoot := t.TempDir()
	outsideRequestPath := filepath.Join(t.TempDir(), "outside-request.md")
	if err := os.WriteFile(outsideRequestPath, []byte(requestFixture("REQ-999", "pending")), 0o644); err != nil {
		t.Fatal(err)
	}
	outsideUserRequestPath := filepath.Join(t.TempDir(), "outside-input.md")
	if err := os.WriteFile(outsideUserRequestPath, []byte("---\nid: UR-999\n---\noutside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	requestLinkPath := filepath.Join(repositoryRoot, "do-work", "queue", "REQ-010-link.md")
	userRequestLinkPath := filepath.Join(repositoryRoot, "do-work", "user-requests", "UR-010", "input.md")
	if err := os.MkdirAll(filepath.Dir(requestLinkPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(userRequestLinkPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideRequestPath, requestLinkPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.Symlink(outsideUserRequestPath, userRequestLinkPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.RequestFiles) != 1 || snapshot.RequestFiles[0].ParsedDocument != nil || snapshot.RequestFiles[0].TypedRecord.RequestID != "" {
		t.Fatalf("symlinked REQ content entered snapshot: %#v", snapshot.RequestFiles)
	}
	if len(snapshot.UserRequestFiles) != 1 || snapshot.UserRequestFiles[0].ParsedDocument != nil || snapshot.UserRequestFiles[0].TypedRecord.RequestID != "" {
		t.Fatalf("symlinked UR content entered snapshot: %#v", snapshot.UserRequestFiles)
	}
	warningText := strings.Join(snapshot.WarningMessages, "\n")
	if !strings.Contains(warningText, requestLinkPath) || !strings.Contains(warningText, userRequestLinkPath) || strings.Count(warningText, "symlink") < 2 {
		t.Fatalf("symlink warnings = %q", warningText)
	}
}

func TestDiscoverRepositoryRetainsRootRunManifestEvidence(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/runs/finished/manifest.md", "# Run\nStatus: consumed\n")
	writeRepositoryFixture(t, repositoryRoot, "do-work/runs/live/manifest.md", "Status: in-progress\n")
	writeRepositoryFixture(t, repositoryRoot, "do-work/runs/live/nested/manifest.md", "Status: consumed\n")
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.RunManifestFiles) != 2 {
		t.Fatalf("run manifests = %#v", snapshot.RunManifestFiles)
	}
	if snapshot.RunManifestFiles[0].RunDirectory != "runs/finished" || snapshot.RunManifestFiles[0].Status != "consumed" {
		t.Fatalf("first run manifest = %#v", snapshot.RunManifestFiles[0])
	}
	if snapshot.RunManifestFiles[1].RunDirectory != "runs/live" || snapshot.RunManifestFiles[1].Status != "in-progress" {
		t.Fatalf("second run manifest = %#v", snapshot.RunManifestFiles[1])
	}
}

func TestDiscoverRepositoryRetainsBlankAndMalformedRecordEvidence(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-040-blank.md", "")
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-040.md", "header destroyed\n")
	writeRepositoryFixture(t, repositoryRoot, "do-work/user-requests/UR-041/input.md", "---\nid: UR-041\n")
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.DamagedRecords) != 2 {
		t.Fatalf("damaged records = %#v", snapshot.DamagedRecords)
	}
	for _, damaged := range snapshot.DamagedRecords {
		if damaged.RelativePath == "" || damaged.ParseFailure == "" || damaged.RecordKind == "" {
			t.Fatalf("incomplete damage evidence = %#v", damaged)
		}
	}
}

func TestDiscoverRepositoryDoesNotClassifyLegacyInputAsBlankedRecord(t *testing.T) {
	repositoryRoot := t.TempDir()
	legacyInput := "# Legacy request\n\nThis predates frontmatter and remains valid source material.\n"
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-003/input.md", legacyInput)
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-004.md", "header destroyed\n")
	writeRepositoryFixture(t, repositoryRoot, "do-work/queue/REQ-041-blank.md", "")

	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.DamagedRecords) != 2 {
		t.Fatalf("damaged records = %#v", snapshot.DamagedRecords)
	}
	for _, damaged := range snapshot.DamagedRecords {
		if damaged.RelativePath == "archive/UR-003/input.md" {
			t.Fatalf("legacy input was classified as damaged: %#v", damaged)
		}
	}
	for _, warning := range snapshot.WarningMessages {
		if strings.Contains(warning, "archive/UR-003/input.md") {
			t.Fatalf("valid legacy input was classified as an inspection warning: %s", warning)
		}
	}
	if len(snapshot.UserRequestFiles) != 1 || snapshot.UserRequestFiles[0].ParseFailure == "" {
		t.Fatalf("legacy input inspection evidence was not retained: %#v", snapshot.UserRequestFiles)
	}
}

func TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass(t *testing.T) {
	productionPath := filepath.Join("..", "..", "..", "..", "..", "..", "do-work", "archive", "UR-003", "input.md")
	productionBytes, err := os.ReadFile(productionPath)
	if os.IsNotExist(err) {
		t.Skip("production do-work corpus is not present")
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(productionBytes) != 5608 {
		t.Fatalf("production legacy fixture changed size: got %d bytes", len(productionBytes))
	}
	repositoryRoot := t.TempDir()
	writeRepositoryFixture(t, repositoryRoot, "do-work/archive/UR-003/input.md", string(productionBytes))
	snapshot, err := DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.DamagedRecords) != 0 {
		t.Fatalf("production legacy input class was damaged: %#v", snapshot.DamagedRecords)
	}
}
