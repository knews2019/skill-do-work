package publication

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const publicationReservationDirectory = "do-work/.req-reservations"

var publicationReservationIDPattern = regexp.MustCompile(`^REQ-([0-9]+)$`)

func canonicalReservationPath(requestID string) (string, error) {
	number, ok := publicationReservationNumber(requestID)
	if !ok {
		return "", fmt.Errorf("request id %q has no positive numeric reservation identity", requestID)
	}
	return fmt.Sprintf("%s/REQ-%03d", publicationReservationDirectory, number), nil
}

func publicationReservationNumber(text string) (int, bool) {
	match := publicationReservationIDPattern.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return 0, false
	}
	number, err := strconv.Atoi(match[1])
	return number, err == nil && number > 0
}

// existingReservationPaths returns every concrete marker spelling that claims
// requestID's numeric identity. New writers use canonicalReservationPath, while
// this read remains width-agnostic until legacy fixed-six markers have drained.
func existingReservationPaths(repositoryRoot, requestID string) ([]string, error) {
	wantedNumber, ok := publicationReservationNumber(requestID)
	if !ok {
		return nil, fmt.Errorf("request id %q has no positive numeric reservation identity", requestID)
	}
	absoluteDirectory := filepath.Join(repositoryRoot, filepath.FromSlash(publicationReservationDirectory))
	directoryInfo, err := os.Lstat(absoluteDirectory)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !directoryInfo.IsDir() || directoryInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("reservation directory is not a real directory")
	}
	directoryRoot, err := os.OpenRoot(absoluteDirectory)
	if err != nil {
		return nil, err
	}
	defer directoryRoot.Close()
	openedInfo, err := directoryRoot.Stat(".")
	if err != nil || !os.SameFile(directoryInfo, openedInfo) {
		return nil, fmt.Errorf("reservation directory changed while opening its rooted handle")
	}
	entries, err := fs.ReadDir(directoryRoot.FS(), ".")
	if err != nil {
		return nil, err
	}
	currentInfo, err := os.Lstat(absoluteDirectory)
	if err != nil || !currentInfo.IsDir() || currentInfo.Mode()&os.ModeSymlink != 0 || !os.SameFile(directoryInfo, currentInfo) {
		return nil, fmt.Errorf("reservation directory changed during inspection")
	}
	paths := []string{}
	for _, entry := range entries {
		number, parsed := publicationReservationNumber(entry.Name())
		if parsed && number == wantedNumber {
			paths = append(paths, publicationReservationDirectory+"/"+entry.Name())
		}
	}
	sort.Strings(paths)
	return paths, nil
}
