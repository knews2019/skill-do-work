// Package repositorymodel discovers do-work repositories into one typed
// snapshot and allocates durable collision-free request reservations.
package repositorymodel

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
)

var ErrRepositoryNotFound = errors.New("do-work repository not found")

var requestNumberPattern = regexp.MustCompile(`(?i)^REQ-0*([0-9]+)`)
var reservationNumberPattern = regexp.MustCompile(`^REQ-([0-9]+)$`)
var checkpointClaimPattern = regexp.MustCompile(`^\s*-\s+(REQ-0*[0-9]+)\s*:`)
var checkpointWriterPattern = regexp.MustCompile(`\s+—\s+writer:\s*(\S(?:.*\S)?)\s*$`)
var checkpointClaimedAtPattern = regexp.MustCompile(`\s+—\s+claimed\s+(.*?)(?:\s+—\s+writer:|\s*$)`)

// beforeReservationMarkerCreate is a test seam for a deterministic directory swap.
var beforeReservationMarkerCreate = func(string) {}

// RequestFile is one discovered REQ document and its exact path evidence.
type RequestFile struct {
	AbsolutePath   string
	RelativePath   string
	TreeSection    string
	FilenameID     string
	TypedRecord    requestmodel.RequestRecord
	ParsedDocument *requestmodel.RequestDocument
	ContentBytes   []byte
	ModifiedAt     time.Time
	ParseFailure   string
}

// UserRequestFile is one active or archived UR input document.
type UserRequestFile struct {
	AbsolutePath   string
	RelativePath   string
	TreeSection    string
	TypedRecord    requestmodel.RequestRecord
	ParsedDocument *requestmodel.RequestDocument
	ContentBytes   []byte
	ParseFailure   string
}

// DamagedRecordFile is a contained REQ or UR record whose bytes cannot be
// projected into frontmatter. The bytes and parse failure stay first-class so
// diagnostics never have to recover them from a warning string.
type DamagedRecordFile struct {
	AbsolutePath string
	RelativePath string
	RecordKind   string
	ContentBytes []byte
	ParseFailure string
}

// RunManifestFile is one root run manifest. Cleanup consumes only manifests
// whose exact status is "consumed"; every other run remains durable evidence.
type RunManifestFile struct {
	AbsolutePath string
	RelativePath string
	RunDirectory string
	Status       string
}

// ReservationFile is one durable request-number marker.
type ReservationFile struct {
	RequestNumber int
	RequestID     string
	AbsolutePath  string
}

// CollisionEvidence retains every exact path claiming one request id.
type CollisionEvidence struct {
	RequestID  string
	ClaimPaths []string
}

// CheckpointClaimEvidence is one checkpoint claim header. It is
// policy-free discovery evidence; callers decide whether a claim blocks work.
type CheckpointClaimEvidence struct {
	RequestID    string
	ClaimedAt    string
	Writer       string
	HasWriter    bool
	RelativePath string
	SourceLine   int
	HeaderText   string
}

// RepositorySnapshot is one complete typed read of a do-work tree.
type RepositorySnapshot struct {
	RepositoryRoot       string
	DoWorkRoot           string
	RequestFiles         []*RequestFile
	RequestsByID         map[string][]*RequestFile
	CheckpointClaimsByID map[string][]CheckpointClaimEvidence
	CheckpointPath       string
	CheckpointBytes      []byte
	CheckpointExists     bool
	UserRequestFiles     []*UserRequestFile
	RunManifestFiles     []RunManifestFile
	ReservationFiles     []ReservationFile
	CollisionEntries     []CollisionEvidence
	DamagedRecords       []DamagedRecordFile
	StrayRequestPaths    []string
	WarningMessages      []string
}

// FindRepositoryRoot walks upward from a file or directory to a queue root.
func FindRepositoryRoot(startPath string) (string, error) {
	absoluteStart, absoluteError := filepath.Abs(startPath)
	if absoluteError != nil {
		return "", fmt.Errorf("resolving repository search path %s: %w", startPath, absoluteError)
	}
	currentDirectory := absoluteStart
	if startInfo, statError := os.Stat(absoluteStart); statError == nil && !startInfo.IsDir() {
		currentDirectory = filepath.Dir(absoluteStart)
	}
	for {
		doWorkPath := filepath.Join(currentDirectory, "do-work")
		if doWorkInfo, statError := os.Stat(doWorkPath); statError == nil && doWorkInfo.IsDir() {
			if _, skillError := os.Stat(filepath.Join(doWorkPath, "SKILL.md")); os.IsNotExist(skillError) {
				return currentDirectory, nil
			}
		}
		parentDirectory := filepath.Dir(currentDirectory)
		if parentDirectory == currentDirectory {
			return "", fmt.Errorf("%w walking up from %s", ErrRepositoryNotFound, startPath)
		}
		currentDirectory = parentDirectory
	}
}

// DiscoverRepository walks the do-work tree once and returns typed evidence.
func DiscoverRepository(repositoryRoot string) (*RepositorySnapshot, error) {
	absoluteRoot, absoluteError := filepath.Abs(repositoryRoot)
	if absoluteError != nil {
		return nil, fmt.Errorf("resolving repository root %s: %w", repositoryRoot, absoluteError)
	}
	doWorkRoot := filepath.Join(absoluteRoot, "do-work")
	doWorkInfo, statError := os.Lstat(doWorkRoot)
	if statError != nil {
		return nil, fmt.Errorf("%w at %s: %v", ErrRepositoryNotFound, doWorkRoot, statError)
	}
	if !doWorkInfo.IsDir() || doWorkInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("do-work root %s is not a real directory", doWorkRoot)
	}
	snapshot := &RepositorySnapshot{
		RepositoryRoot:       absoluteRoot,
		DoWorkRoot:           doWorkRoot,
		RequestsByID:         map[string][]*RequestFile{},
		CheckpointClaimsByID: map[string][]CheckpointClaimEvidence{},
	}
	discoveryRoot, rootError := os.OpenRoot(doWorkRoot)
	if rootError != nil {
		return nil, fmt.Errorf("opening rooted do-work tree %s: %w", doWorkRoot, rootError)
	}
	defer discoveryRoot.Close()
	openedRootInfo, openedRootStatError := discoveryRoot.Stat(".")
	if openedRootStatError != nil || !os.SameFile(doWorkInfo, openedRootInfo) {
		return nil, fmt.Errorf("do-work root %s changed while opening its discovery handle", doWorkRoot)
	}
	walkError := filepath.WalkDir(doWorkRoot, func(absolutePath string, directoryEntry fs.DirEntry, entryError error) error {
		if entryError != nil {
			snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot inspect %s: %v", absolutePath, entryError))
			return nil
		}
		relativePath, relativeError := filepath.Rel(doWorkRoot, absolutePath)
		if relativeError != nil {
			snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot make %s relative to do-work root: %v", absolutePath, relativeError))
			return nil
		}
		relativeSlashPath := filepath.ToSlash(relativePath)
		pathParts := strings.Split(relativeSlashPath, "/")
		topSection := pathParts[0]
		if directoryEntry.IsDir() {
			if absolutePath == doWorkRoot {
				return nil
			}
			if topSection == "deliverables" || directoryEntry.Name() == "assets" {
				return fs.SkipDir
			}
			if strings.HasPrefix(directoryEntry.Name(), ".") && relativeSlashPath != ".req-reservations" {
				return fs.SkipDir
			}
			return nil
		}

		if topSection == ".req-reservations" && len(pathParts) == 2 {
			if requestNumber, parsed := requestNumberFromReservationName(directoryEntry.Name()); parsed {
				snapshot.ReservationFiles = append(snapshot.ReservationFiles, ReservationFile{
					RequestNumber: requestNumber, RequestID: formatRequestID(requestNumber), AbsolutePath: absolutePath,
				})
			}
			return nil
		}
		if relativeSlashPath == "CHECKPOINT.md" {
			fileBytes, _, readError := readContainedRegularFile(discoveryRoot, relativeSlashPath, absolutePath)
			if readError != nil {
				snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot read checkpoint %s: %v", absolutePath, readError))
				return nil
			}
			snapshot.CheckpointPath = relativeSlashPath
			snapshot.CheckpointBytes = append([]byte(nil), fileBytes...)
			snapshot.CheckpointExists = true
			projectCheckpointClaims(snapshot, relativeSlashPath, fileBytes)
			return nil
		}

		baseName := directoryEntry.Name()
		if topSection == "runs" && len(pathParts) == 3 && baseName == "manifest.md" {
			fileBytes, _, readError := readContainedRegularFile(discoveryRoot, relativeSlashPath, absolutePath)
			if readError != nil {
				snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot read run manifest %s: %v", absolutePath, readError))
				return nil
			}
			snapshot.RunManifestFiles = append(snapshot.RunManifestFiles, RunManifestFile{
				AbsolutePath: absolutePath, RelativePath: relativeSlashPath,
				RunDirectory: filepath.ToSlash(filepath.Dir(relativeSlashPath)), Status: manifestStatus(fileBytes),
			})
			return nil
		}
		if topSection == "runs" {
			return nil
		}
		if strings.HasPrefix(strings.ToUpper(baseName), "REQ-") && strings.HasSuffix(strings.ToLower(baseName), ".md") {
			switch topSection {
			case "queue", "working", "archive":
				requestFile := loadRequestFile(snapshot, discoveryRoot, absolutePath, relativeSlashPath, topSection)
				snapshot.RequestFiles = append(snapshot.RequestFiles, requestFile)
			default:
				snapshot.StrayRequestPaths = append(snapshot.StrayRequestPaths, absolutePath)
			}
			return nil
		}
		if strings.HasPrefix(strings.ToUpper(baseName), "UR-") && strings.HasSuffix(strings.ToLower(baseName), ".md") &&
			(topSection == "queue" || topSection == "working" || topSection == "archive") {
			loadStandaloneUserRequestRecord(snapshot, discoveryRoot, absolutePath, relativeSlashPath)
			return nil
		}
		if baseName == "input.md" && (topSection == "user-requests" || topSection == "archive") {
			userRequest := loadUserRequestFile(snapshot, discoveryRoot, absolutePath, relativeSlashPath, topSection)
			snapshot.UserRequestFiles = append(snapshot.UserRequestFiles, userRequest)
		}
		return nil
	})
	if walkError != nil {
		return nil, fmt.Errorf("walking do-work tree %s: %w", doWorkRoot, walkError)
	}

	for _, requestFile := range snapshot.RequestFiles {
		requestID := requestFile.TypedRecord.RequestID
		if requestID == "" {
			requestID = requestFile.FilenameID
		}
		if requestID == "" {
			snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("request path %s has no usable id", requestFile.AbsolutePath))
			continue
		}
		snapshot.RequestsByID[requestID] = append(snapshot.RequestsByID[requestID], requestFile)
	}
	claimedPathsByID := map[string]map[string]bool{}
	for _, requestFile := range snapshot.RequestFiles {
		claimIDs := []string{requestFile.FilenameID, requestFile.TypedRecord.RequestID}
		if claimIDs[0] != "" && claimIDs[1] != "" && claimIDs[0] != claimIDs[1] {
			snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("request path %s claims filename id %s but frontmatter id %s", requestFile.AbsolutePath, claimIDs[0], claimIDs[1]))
		}
		for _, claimID := range claimIDs {
			if claimID == "" {
				continue
			}
			if normalizedID := requestIDFromText(claimID); normalizedID != "" {
				claimID = normalizedID
			}
			if claimedPathsByID[claimID] == nil {
				claimedPathsByID[claimID] = map[string]bool{}
			}
			claimedPathsByID[claimID][requestFile.AbsolutePath] = true
		}
	}
	for requestID, claimedPaths := range claimedPathsByID {
		if len(claimedPaths) < 2 {
			continue
		}
		paths := make([]string, 0, len(claimedPaths))
		for claimedPath := range claimedPaths {
			paths = append(paths, claimedPath)
		}
		sort.Strings(paths)
		snapshot.CollisionEntries = append(snapshot.CollisionEntries, CollisionEvidence{RequestID: requestID, ClaimPaths: paths})
	}
	sort.Slice(snapshot.CollisionEntries, func(leftIndex, rightIndex int) bool {
		return requestIDLess(snapshot.CollisionEntries[leftIndex].RequestID, snapshot.CollisionEntries[rightIndex].RequestID)
	})
	sort.Slice(snapshot.RequestFiles, func(leftIndex, rightIndex int) bool {
		leftID := snapshot.RequestFiles[leftIndex].TypedRecord.RequestID
		if leftID == "" {
			leftID = snapshot.RequestFiles[leftIndex].FilenameID
		}
		rightID := snapshot.RequestFiles[rightIndex].TypedRecord.RequestID
		if rightID == "" {
			rightID = snapshot.RequestFiles[rightIndex].FilenameID
		}
		if leftID == rightID {
			return snapshot.RequestFiles[leftIndex].AbsolutePath < snapshot.RequestFiles[rightIndex].AbsolutePath
		}
		return requestIDLess(leftID, rightID)
	})
	sort.Slice(snapshot.UserRequestFiles, func(leftIndex, rightIndex int) bool {
		return snapshot.UserRequestFiles[leftIndex].AbsolutePath < snapshot.UserRequestFiles[rightIndex].AbsolutePath
	})
	sort.Slice(snapshot.ReservationFiles, func(leftIndex, rightIndex int) bool {
		return snapshot.ReservationFiles[leftIndex].RequestNumber < snapshot.ReservationFiles[rightIndex].RequestNumber
	})
	sort.Slice(snapshot.RunManifestFiles, func(leftIndex, rightIndex int) bool {
		return snapshot.RunManifestFiles[leftIndex].RunDirectory < snapshot.RunManifestFiles[rightIndex].RunDirectory
	})
	sort.Strings(snapshot.StrayRequestPaths)
	sort.Slice(snapshot.DamagedRecords, func(leftIndex, rightIndex int) bool {
		return snapshot.DamagedRecords[leftIndex].RelativePath < snapshot.DamagedRecords[rightIndex].RelativePath
	})
	return snapshot, nil
}

func manifestStatus(fileBytes []byte) string {
	for _, line := range strings.Split(string(fileBytes), "\n") {
		name, value, found := strings.Cut(strings.TrimSpace(line), ":")
		if found && strings.EqualFold(strings.TrimSpace(name), "status") {
			return strings.ToLower(strings.TrimSpace(value))
		}
	}
	return ""
}

func projectCheckpointClaims(snapshot *RepositorySnapshot, relativePath string, fileBytes []byte) {
	lines := strings.Split(string(fileBytes), "\n")
	headingLine, sectionEnd, found := checkpointSectionBounds(lines)
	if !found {
		// Older checkpoints predate the section heading. Keep their structural
		// claim headers visible to selection and recovery until the explicit
		// checkpoint writer upgrades the document.
		headingLine = -1
		sectionEnd = len(lines)
	}
	for lineIndex := headingLine + 1; lineIndex < sectionEnd; lineIndex++ {
		headerText := lines[lineIndex]
		headerText = strings.TrimSuffix(headerText, "\r")
		claimMatch := checkpointClaimPattern.FindStringSubmatch(headerText)
		if claimMatch == nil {
			continue
		}
		requestID := requestIDFromText(claimMatch[1])
		if requestID == "" {
			continue
		}
		claimedAt := ""
		if claimedMatch := checkpointClaimedAtPattern.FindStringSubmatch(headerText); claimedMatch != nil {
			claimedAt = strings.TrimSpace(claimedMatch[1])
		}
		writer := ""
		hasWriter := false
		if writerMatch := checkpointWriterPattern.FindStringSubmatch(headerText); writerMatch != nil {
			writer = strings.TrimSpace(writerMatch[1])
			hasWriter = writer != ""
		}
		snapshot.CheckpointClaimsByID[requestID] = append(snapshot.CheckpointClaimsByID[requestID], CheckpointClaimEvidence{
			RequestID: requestID, ClaimedAt: claimedAt, Writer: writer, HasWriter: hasWriter,
			RelativePath: relativePath, SourceLine: lineIndex + 1, HeaderText: headerText,
		})
	}
}

func checkpointSectionBounds(lines []string) (int, int, bool) {
	for lineIndex, line := range lines {
		if strings.TrimSuffix(line, "\r") != "## In Progress (interrupted)" {
			continue
		}
		for sectionEnd := lineIndex + 1; sectionEnd < len(lines); sectionEnd++ {
			if strings.HasPrefix(strings.TrimSuffix(lines[sectionEnd], "\r"), "## ") {
				return lineIndex, sectionEnd, true
			}
		}
		return lineIndex, len(lines), true
	}
	return 0, 0, false
}

// ReserveNextRequestID exclusively creates and returns the next marker.
func ReserveNextRequestID(snapshot *RepositorySnapshot) (ReservationFile, error) {
	if snapshot == nil || snapshot.DoWorkRoot == "" {
		return ReservationFile{}, fmt.Errorf("repository snapshot is required")
	}
	highestNumber := 0
	for _, requestFile := range snapshot.RequestFiles {
		for _, requestText := range []string{requestFile.FilenameID, requestFile.TypedRecord.RequestID} {
			if requestNumber, parsed := requestNumberFromText(requestText); parsed && requestNumber > highestNumber {
				highestNumber = requestNumber
			}
		}
	}
	for _, reservation := range snapshot.ReservationFiles {
		if reservation.RequestNumber > highestNumber {
			highestNumber = reservation.RequestNumber
		}
	}
	reservationDirectory := filepath.Join(snapshot.DoWorkRoot, ".req-reservations")
	reservationStore, storeError := openReservationStore(snapshot.DoWorkRoot, reservationDirectory)
	if storeError != nil {
		return ReservationFile{}, storeError
	}
	defer reservationStore.closeStore()
	for candidateNumber := highestNumber + 1; ; candidateNumber++ {
		markerName := formatRequestID(candidateNumber)
		markerPath := filepath.Join(reservationDirectory, markerName)
		beforeReservationMarkerCreate(reservationDirectory)
		createError := atomicfile.CreateExclusiveAt(reservationStore.directoryRoot, markerName, nil, 0o644)
		if createError == nil {
			if !reservationStore.isCurrent() {
				reservationStore.directoryRoot.Remove(markerName)
				return ReservationFile{}, fmt.Errorf("reservation directory %s changed before marker publication", reservationDirectory)
			}
			return ReservationFile{RequestNumber: candidateNumber, RequestID: formatRequestID(candidateNumber), AbsolutePath: markerPath}, nil
		}
		if errors.Is(createError, os.ErrExist) {
			continue
		}
		return ReservationFile{}, fmt.Errorf("reserving %s: %w", formatRequestID(candidateNumber), createError)
	}
}

type reservationStore struct {
	doWorkPath    string
	doWorkInfo    fs.FileInfo
	doWorkRoot    *os.Root
	directoryPath string
	directoryInfo fs.FileInfo
	directoryRoot *os.Root
}

func openReservationStore(doWorkPath string, reservationDirectory string) (*reservationStore, error) {
	doWorkInfo, doWorkStatError := os.Lstat(doWorkPath)
	if doWorkStatError != nil || !doWorkInfo.IsDir() || doWorkInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("do-work path %s is not a real directory", doWorkPath)
	}
	doWorkRoot, doWorkRootError := os.OpenRoot(doWorkPath)
	if doWorkRootError != nil {
		return nil, fmt.Errorf("opening rooted do-work directory %s: %w", doWorkPath, doWorkRootError)
	}
	openedDoWorkInfo, openedDoWorkStatError := doWorkRoot.Stat(".")
	if openedDoWorkStatError != nil || !os.SameFile(doWorkInfo, openedDoWorkInfo) {
		doWorkRoot.Close()
		return nil, fmt.Errorf("do-work directory %s changed while opening its rooted handle", doWorkPath)
	}
	for {
		directoryInfo, statError := doWorkRoot.Lstat(".req-reservations")
		switch {
		case statError == nil:
			if !directoryInfo.IsDir() || directoryInfo.Mode()&os.ModeSymlink != 0 {
				doWorkRoot.Close()
				return nil, fmt.Errorf("reservation path %s is not a real directory", reservationDirectory)
			}
		case !os.IsNotExist(statError):
			doWorkRoot.Close()
			return nil, fmt.Errorf("checking reservation directory %s: %w", reservationDirectory, statError)
		default:
			mkdirError := doWorkRoot.Mkdir(".req-reservations", 0o755)
			if mkdirError != nil && !os.IsExist(mkdirError) {
				doWorkRoot.Close()
				return nil, fmt.Errorf("creating reservation directory %s: %w", reservationDirectory, mkdirError)
			}
			continue
		}
		break
	}
	directoryInfo, statError := doWorkRoot.Lstat(".req-reservations")
	if statError != nil {
		doWorkRoot.Close()
		return nil, fmt.Errorf("checking reservation directory %s: %w", reservationDirectory, statError)
	}
	directoryRoot, rootError := doWorkRoot.OpenRoot(".req-reservations")
	if rootError != nil {
		doWorkRoot.Close()
		return nil, fmt.Errorf("opening rooted reservation directory %s: %w", reservationDirectory, rootError)
	}
	openedInfo, openedStatError := directoryRoot.Stat(".")
	if openedStatError != nil || !os.SameFile(directoryInfo, openedInfo) {
		directoryRoot.Close()
		doWorkRoot.Close()
		return nil, fmt.Errorf("reservation directory %s changed while opening its rooted handle", reservationDirectory)
	}
	return &reservationStore{
		doWorkPath: doWorkPath, doWorkInfo: doWorkInfo, doWorkRoot: doWorkRoot,
		directoryPath: reservationDirectory, directoryInfo: directoryInfo, directoryRoot: directoryRoot,
	}, nil
}

func (store *reservationStore) isCurrent() bool {
	currentDoWorkInfo, doWorkStatError := os.Lstat(store.doWorkPath)
	if doWorkStatError != nil || !currentDoWorkInfo.IsDir() || currentDoWorkInfo.Mode()&os.ModeSymlink != 0 || !os.SameFile(store.doWorkInfo, currentDoWorkInfo) {
		return false
	}
	currentDirectoryInfo, directoryStatError := os.Lstat(store.directoryPath)
	if directoryStatError != nil || !currentDirectoryInfo.IsDir() || currentDirectoryInfo.Mode()&os.ModeSymlink != 0 || !os.SameFile(store.directoryInfo, currentDirectoryInfo) {
		return false
	}
	rootedDirectoryInfo, rootedStatError := store.doWorkRoot.Lstat(".req-reservations")
	return rootedStatError == nil && rootedDirectoryInfo.IsDir() && rootedDirectoryInfo.Mode()&os.ModeSymlink == 0 && os.SameFile(store.directoryInfo, rootedDirectoryInfo)
}

func (store *reservationStore) closeStore() {
	store.directoryRoot.Close()
	store.doWorkRoot.Close()
}

func loadRequestFile(snapshot *RepositorySnapshot, discoveryRoot *os.Root, absolutePath string, relativePath string, section string) *RequestFile {
	requestFile := &RequestFile{
		AbsolutePath: absolutePath, RelativePath: relativePath, TreeSection: section,
		FilenameID: requestIDFromText(filepath.Base(absolutePath)),
	}
	fileBytes, fileInfo, readError := readContainedRegularFile(discoveryRoot, relativePath, absolutePath)
	if readError != nil {
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot read request %s: %v", absolutePath, readError))
		return requestFile
	}
	requestFile.ContentBytes = append([]byte(nil), fileBytes...)
	requestFile.ModifiedAt = fileInfo.ModTime().UTC().Truncate(time.Second)
	document, parseError := requestmodel.ParseDocument(fileBytes)
	if parseError != nil {
		requestFile.ParseFailure = parseError.Error()
		snapshot.DamagedRecords = append(snapshot.DamagedRecords, DamagedRecordFile{
			AbsolutePath: absolutePath, RelativePath: relativePath, RecordKind: "REQ",
			ContentBytes: append([]byte(nil), fileBytes...), ParseFailure: parseError.Error(),
		})
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot parse request %s: %v", absolutePath, parseError))
		return requestFile
	}
	requestFile.ParsedDocument = document
	requestFile.TypedRecord = document.TypedRecord()
	for _, warning := range document.ParseWarnings() {
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("%s: %s", absolutePath, warning))
	}
	if requestFile.TypedRecord.StatusEvidence.WarningMessage != "" {
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("%s: %s", absolutePath, requestFile.TypedRecord.StatusEvidence.WarningMessage))
	}
	return requestFile
}

func loadUserRequestFile(snapshot *RepositorySnapshot, discoveryRoot *os.Root, absolutePath string, relativePath string, section string) *UserRequestFile {
	userRequest := &UserRequestFile{AbsolutePath: absolutePath, RelativePath: relativePath, TreeSection: section}
	fileBytes, _, readError := readContainedRegularFile(discoveryRoot, relativePath, absolutePath)
	if readError != nil {
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot read user request %s: %v", absolutePath, readError))
		return userRequest
	}
	userRequest.ContentBytes = append([]byte(nil), fileBytes...)
	document, parseError := requestmodel.ParseDocument(fileBytes)
	if parseError != nil {
		userRequest.ParseFailure = parseError.Error()
		if section != "archive" {
			snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot parse user request %s: %v", absolutePath, parseError))
		}
		return userRequest
	}
	userRequest.ParsedDocument = document
	userRequest.TypedRecord = document.TypedRecord()
	return userRequest
}

func loadStandaloneUserRequestRecord(snapshot *RepositorySnapshot, discoveryRoot *os.Root, absolutePath, relativePath string) {
	fileBytes, _, readError := readContainedRegularFile(discoveryRoot, relativePath, absolutePath)
	if readError != nil {
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot read user request record %s: %v", absolutePath, readError))
		return
	}
	if _, parseError := requestmodel.ParseDocument(fileBytes); parseError != nil {
		snapshot.DamagedRecords = append(snapshot.DamagedRecords, DamagedRecordFile{
			AbsolutePath: absolutePath, RelativePath: relativePath, RecordKind: "UR",
			ContentBytes: append([]byte(nil), fileBytes...), ParseFailure: parseError.Error(),
		})
		snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot parse user request record %s: %v", absolutePath, parseError))
	}
}

func readContainedRegularFile(discoveryRoot *os.Root, relativePath string, absolutePath string) ([]byte, fs.FileInfo, error) {
	rootedPath := filepath.FromSlash(relativePath)
	pathInfo, lstatError := discoveryRoot.Lstat(rootedPath)
	if lstatError != nil {
		return nil, nil, lstatError
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 {
		return nil, nil, fmt.Errorf("%s is a symlink; refusing to read linked content", absolutePath)
	}
	if !pathInfo.Mode().IsRegular() {
		return nil, nil, fmt.Errorf("%s is not a regular file", absolutePath)
	}
	openedFile, openError := discoveryRoot.Open(rootedPath)
	if openError != nil {
		return nil, nil, openError
	}
	defer openedFile.Close()
	openedInfo, openedStatError := openedFile.Stat()
	if openedStatError != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return nil, nil, fmt.Errorf("%s changed before its contained read", absolutePath)
	}
	fileBytes, readError := io.ReadAll(openedFile)
	if readError != nil {
		return nil, nil, readError
	}
	afterInfo, afterStatError := openedFile.Stat()
	if afterStatError != nil || !os.SameFile(openedInfo, afterInfo) || openedInfo.Size() != afterInfo.Size() || !openedInfo.ModTime().Equal(afterInfo.ModTime()) {
		return nil, nil, fmt.Errorf("%s changed during its contained read", absolutePath)
	}
	return fileBytes, afterInfo, nil
}

func requestNumberFromText(requestText string) (int, bool) {
	match := requestNumberPattern.FindStringSubmatch(strings.TrimSpace(requestText))
	if match == nil {
		return 0, false
	}
	requestNumber, parseError := strconv.Atoi(match[1])
	return requestNumber, parseError == nil
}

func requestNumberFromReservationName(name string) (int, bool) {
	match := reservationNumberPattern.FindStringSubmatch(name)
	if match == nil {
		return 0, false
	}
	requestNumber, parseError := strconv.Atoi(match[1])
	return requestNumber, parseError == nil && requestNumber > 0
}

func requestIDFromText(requestText string) string {
	requestNumber, parsed := requestNumberFromText(requestText)
	if !parsed {
		return ""
	}
	return formatRequestID(requestNumber)
}

func formatRequestID(requestNumber int) string {
	return fmt.Sprintf("REQ-%03d", requestNumber)
}

func requestIDLess(leftID string, rightID string) bool {
	leftNumber, leftParsed := requestNumberFromText(leftID)
	rightNumber, rightParsed := requestNumberFromText(rightID)
	if leftParsed && rightParsed && leftNumber != rightNumber {
		return leftNumber < rightNumber
	}
	return leftID < rightID
}
