package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// requestIdNumberPattern pulls the numeric part out of a "REQ-NNN" id or a
// "REQ-NNN-slug.md" filename. Zero-padding width is irrelevant — the value is
// what orders ids, so REQ-0142 and REQ-142 are the same number.
var requestIdNumberPattern = regexp.MustCompile(`^REQ-0*(\d+)`)

// requestReservationDirectoryName is the queue-local durable coordination
// store used by next-req. The board walk skips hidden directories, so markers
// never become cards or stray-REQ warnings; capture commits them as queue
// metadata alongside the UR/REQ they reserved the number for.
const requestReservationDirectoryName = ".req-reservations"

// nextRequestNumber atomically reserves and returns the next free REQ number for
// a repo: one past the highest number already in use across do-work/queue/,
// do-work/working/, the whole do-work/archive/ subtree (including nested
// archive/UR-NNN/ folders), and prior reservation markers; 1 when neither the
// tree nor the reservation store holds a number.
//
// This is actions/capture.md's existing allocation rule, executed instead of
// eyeballed. It reuses enumerateDoWorkTree (walk.go) — the same walk the board
// builds on, with the same pruning of deliverables/, runs/, and assets/ — rather
// than introducing a second scan that could drift from it.
//
// Gaps are deliberately fine (REQ-072 requirement 7): reservations are durable,
// so a capture that stops after allocation consumes a number without creating a
// REQ. Nothing in the skill requires a contiguous sequence; a gap is safer than
// handing the abandoned number to another capture.
//
// Atomicity comes from O_CREATE|O_EXCL on a per-number marker. Concurrent callers
// can compute the same candidate, but exactly one creates that marker; every
// loser advances and retries. No global lock or stale-lock recovery is needed.
func nextRequestNumber(repoRootOverride string) (int, error) {
	repoRoot, resolveError := resolveRepoRootOrDefault(repoRootOverride)
	if resolveError != nil {
		return 0, resolveError
	}
	discovered, walkError := enumerateDoWorkTree(repoRoot)
	if walkError != nil {
		return 0, walkError
	}

	highestNumberInUse := 0
	for _, requestFile := range discovered.RequestFiles {
		for _, candidateNumber := range requestNumbersInFile(requestFile.AbsolutePath) {
			if candidateNumber > highestNumberInUse {
				highestNumberInUse = candidateNumber
			}
		}
	}

	reservationStore, reservationDirectoryError := ensureRequestReservationDirectory(repoRoot)
	if reservationDirectoryError != nil {
		return 0, reservationDirectoryError
	}
	defer reservationStore.root.Close()
	reservationEntries, reservationReadError := fs.ReadDir(reservationStore.root.FS(), ".")
	if reservationReadError != nil {
		return 0, fmt.Errorf("queue-kanban: reading REQ reservations at %s: %w", reservationStore.directoryPath, reservationReadError)
	}
	for _, reservationEntry := range reservationEntries {
		if reservationEntry.IsDir() {
			continue
		}
		if reservedNumber, parsedOk := requestNumberFromText(reservationEntry.Name()); parsedOk && reservedNumber > highestNumberInUse {
			highestNumberInUse = reservedNumber
		}
	}

	for candidateNumber := highestNumberInUse + 1; ; candidateNumber++ {
		reservationFileName := requestReservationFileName(candidateNumber)
		reservationPath := filepath.Join(reservationStore.directoryPath, reservationFileName)
		reservationFile, createError := reservationStore.root.OpenFile(reservationFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if createError == nil {
			if closeError := reservationFile.Close(); closeError != nil {
				return 0, fmt.Errorf("queue-kanban: closing REQ reservation %s: %w", reservationPath, closeError)
			}
			return candidateNumber, nil
		}
		if os.IsExist(createError) {
			continue
		}
		return 0, fmt.Errorf("queue-kanban: reserving REQ-%d at %s: %w", candidateNumber, reservationPath, createError)
	}
}

type requestReservationStore struct {
	directoryPath string
	root          *os.Root
}

// ensureRequestReservationDirectory resolves the repo and do-work roots before
// creating the marker directory. It returns a rooted filesystem handle so a
// later symlink swap cannot turn the checked path into an arbitrary write.
func ensureRequestReservationDirectory(repoRoot string) (*requestReservationStore, error) {
	absoluteRepoRoot, absoluteRootError := filepath.Abs(repoRoot)
	if absoluteRootError != nil {
		return nil, fmt.Errorf("queue-kanban: resolving repository root %s: %w", repoRoot, absoluteRootError)
	}
	resolvedRepoRoot, repoResolveError := filepath.EvalSymlinks(absoluteRepoRoot)
	if repoResolveError != nil {
		return nil, fmt.Errorf("queue-kanban: resolving repository root %s: %w", absoluteRepoRoot, repoResolveError)
	}
	resolvedDoWorkRoot, doWorkResolveError := filepath.EvalSymlinks(filepath.Join(absoluteRepoRoot, "do-work"))
	if doWorkResolveError != nil {
		return nil, fmt.Errorf("queue-kanban: resolving do-work root: %w", doWorkResolveError)
	}
	relativeDoWorkPath, relativeError := filepath.Rel(resolvedRepoRoot, resolvedDoWorkRoot)
	if relativeError != nil || relativeDoWorkPath == ".." || strings.HasPrefix(relativeDoWorkPath, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("queue-kanban: do-work root resolves outside the repository — refusing to reserve a REQ number")
	}
	repositoryRoot, repositoryOpenError := os.OpenRoot(resolvedRepoRoot)
	if repositoryOpenError != nil {
		return nil, fmt.Errorf("queue-kanban: opening repository root %s: %w", resolvedRepoRoot, repositoryOpenError)
	}
	defer repositoryRoot.Close()
	doWorkRoot, doWorkOpenError := repositoryRoot.OpenRoot(relativeDoWorkPath)
	if doWorkOpenError != nil {
		return nil, fmt.Errorf("queue-kanban: opening contained do-work root %s: %w", resolvedDoWorkRoot, doWorkOpenError)
	}
	defer doWorkRoot.Close()

	reservationDirectory := filepath.Join(resolvedDoWorkRoot, requestReservationDirectoryName)
	for {
		reservationInfo, lstatError := doWorkRoot.Lstat(requestReservationDirectoryName)
		switch {
		case lstatError == nil:
			if !reservationInfo.IsDir() {
				return nil, fmt.Errorf("queue-kanban: REQ reservation path %s is not a directory", reservationDirectory)
			}
			reservationRoot, reservationOpenError := doWorkRoot.OpenRoot(requestReservationDirectoryName)
			if reservationOpenError != nil {
				return nil, fmt.Errorf("queue-kanban: opening contained REQ reservation directory %s: %w", reservationDirectory, reservationOpenError)
			}
			return &requestReservationStore{directoryPath: reservationDirectory, root: reservationRoot}, nil
		case !os.IsNotExist(lstatError):
			return nil, fmt.Errorf("queue-kanban: checking REQ reservation directory %s: %w", reservationDirectory, lstatError)
		}

		mkdirError := doWorkRoot.Mkdir(requestReservationDirectoryName, 0o755)
		if mkdirError == nil {
			continue
		}
		if !os.IsExist(mkdirError) {
			return nil, fmt.Errorf("queue-kanban: creating REQ reservation directory %s: %w", reservationDirectory, mkdirError)
		}
		// Another allocator created the path between Lstat and Mkdir. Re-check
		// its type instead of trusting that the colliding entry is a directory.
	}
}

// requestReservationFileName uses fixed-width decimal names so directory
// listings remain naturally ordered while numeric parsing stays width-agnostic.
func requestReservationFileName(requestNumber int) string {
	return fmt.Sprintf("REQ-%06d", requestNumber)
}

// requestNumbersInFile returns every REQ number a single file lays claim to: the
// one in its filename and, when the frontmatter carries a different `id:`, that
// one too. Both count — a file renamed away from its id still owns its id's
// number, and handing that number out again would manufacture the duplicate this
// tool is meant to detect.
func requestNumbersInFile(absolutePath string) []int {
	var claimedNumbers []int
	if fileNameNumber, parsedOk := requestNumberFromText(baseNameOf(absolutePath)); parsedOk {
		claimedNumbers = append(claimedNumbers, fileNameNumber)
	}
	if frontmatterNumber, parsedOk := requestNumberFromFrontmatterId(absolutePath); parsedOk {
		claimedNumbers = append(claimedNumbers, frontmatterNumber)
	}
	return claimedNumbers
}

// requestNumberFromText parses the numeric part of a REQ id or REQ filename.
func requestNumberFromText(text string) (int, bool) {
	match := requestIdNumberPattern.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return 0, false
	}
	parsedNumber, parseError := strconv.Atoi(match[1])
	if parseError != nil {
		return 0, false
	}
	return parsedNumber, true
}

// requestNumberFromFrontmatterId reads just the frontmatter `id:` field. It uses
// the existing lenient frontmatter reader so a REQ with slightly malformed YAML
// still contributes its number instead of silently dropping out of the max.
func requestNumberFromFrontmatterId(absolutePath string) (int, bool) {
	fileBytes, readError := os.ReadFile(absolutePath)
	if readError != nil {
		return 0, false
	}
	yamlText, _, _, hasFrontmatter := splitFrontmatter(string(fileBytes))
	if !hasFrontmatter {
		return 0, false
	}
	fields := lenientFrontmatterFields(yamlText)
	return requestNumberFromText(coerceScalarToString(fields["id"]))
}

// baseNameOf returns the last path element. Kept local and tiny so allocate.go
// does not pull filepath in for one call.
func baseNameOf(path string) string {
	if slashIndex := strings.LastIndexAny(path, `/\`); slashIndex >= 0 {
		return path[slashIndex+1:]
	}
	return path
}
