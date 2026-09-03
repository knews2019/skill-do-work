package publication

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

const publicationReservationDirectory = "do-work/.req-reservations"

var publicationReservationIDPattern = regexp.MustCompile(`^REQ-([0-9]+)$`)
var beforeReservationMarkerOpen = func(string) {}

type existingReservationMarker struct {
	path              string
	name              string
	directoryIdentity os.FileInfo
	fileIdentity      os.FileInfo
}

func canonicalReservationPath(requestID string) (string, error) {
	number, ok := publicationReservationNumber(requestID)
	if !ok {
		return "", fmt.Errorf("request id %q has no positive numeric reservation identity", requestID)
	}
	return fmt.Sprintf("%s/REQ-%03d", publicationReservationDirectory, number), nil
}

func publicationReservationNumber(text string) (int, bool) {
	match := publicationReservationIDPattern.FindStringSubmatch(text)
	if match == nil {
		return 0, false
	}
	number, err := strconv.Atoi(match[1])
	return number, err == nil && number > 0
}

// existingReservationMarkers returns every concrete marker spelling that claims
// requestID's numeric identity. New writers use canonicalReservationPath, while
// this read remains width-agnostic until legacy fixed-six markers have drained.
func existingReservationMarkers(repositoryRoot, requestID string) ([]existingReservationMarker, error) {
	wantedNumber, ok := publicationReservationNumber(requestID)
	if !ok {
		return nil, fmt.Errorf("request id %q has no positive numeric reservation identity", requestID)
	}
	directoryRoot, directoryInfo, absoluteDirectory, err := openPublicationReservationDirectory(repositoryRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer directoryRoot.Close()
	entries, err := fs.ReadDir(directoryRoot.FS(), ".")
	if err != nil {
		return nil, err
	}
	markers := []existingReservationMarker{}
	for _, entry := range entries {
		number, parsed := publicationReservationNumber(entry.Name())
		if parsed && number == wantedNumber {
			identity, identityError := directoryRoot.Lstat(entry.Name())
			if identityError != nil {
				return nil, fmt.Errorf("inspect reservation marker %s: %w", entry.Name(), identityError)
			}
			markers = append(markers, existingReservationMarker{
				path: publicationReservationDirectory + "/" + entry.Name(), name: entry.Name(),
				directoryIdentity: directoryInfo, fileIdentity: identity,
			})
		}
	}
	currentInfo, err := os.Lstat(absoluteDirectory)
	if err != nil || !currentInfo.IsDir() || currentInfo.Mode()&os.ModeSymlink != 0 || !os.SameFile(directoryInfo, currentInfo) {
		return nil, fmt.Errorf("reservation directory changed during inspection")
	}
	sort.Slice(markers, func(left, right int) bool { return markers[left].path < markers[right].path })
	return markers, nil
}

func reservationMarkerPaths(markers []existingReservationMarker) []string {
	paths := make([]string, len(markers))
	for index, marker := range markers {
		paths[index] = marker.path
	}
	return paths
}

func openPublicationReservationDirectory(repositoryRoot string) (*os.Root, os.FileInfo, string, error) {
	absoluteDirectory := filepath.Join(repositoryRoot, filepath.FromSlash(publicationReservationDirectory))
	directoryInfo, err := os.Lstat(absoluteDirectory)
	if err != nil {
		return nil, nil, absoluteDirectory, err
	}
	if !directoryInfo.IsDir() || directoryInfo.Mode()&os.ModeSymlink != 0 {
		return nil, nil, absoluteDirectory, fmt.Errorf("reservation directory is not a real directory")
	}
	directoryRoot, err := os.OpenRoot(absoluteDirectory)
	if err != nil {
		return nil, nil, absoluteDirectory, err
	}
	openedInfo, err := directoryRoot.Stat(".")
	if err != nil || !os.SameFile(directoryInfo, openedInfo) {
		directoryRoot.Close()
		return nil, nil, absoluteDirectory, fmt.Errorf("reservation directory changed while opening its rooted handle")
	}
	return directoryRoot, directoryInfo, absoluteDirectory, nil
}

// readReservationMarker reads only the concrete regular file found by the
// rooted discovery pass. Symlinks are never accepted as marker-byte authority,
// and both the directory and file identity are revalidated around the read.
func readReservationMarker(repositoryRoot string, marker existingReservationMarker) ([]byte, error) {
	directoryRoot, directoryInfo, absoluteDirectory, err := openPublicationReservationDirectory(repositoryRoot)
	if err != nil {
		return nil, err
	}
	defer directoryRoot.Close()
	if !os.SameFile(marker.directoryIdentity, directoryInfo) {
		return nil, fmt.Errorf("reservation directory changed after inspection")
	}
	beforeOpen, err := directoryRoot.Lstat(marker.name)
	if err != nil || !beforeOpen.Mode().IsRegular() || !os.SameFile(marker.fileIdentity, beforeOpen) {
		return nil, fmt.Errorf("reservation marker changed or is not a regular non-symlink file")
	}
	beforeReservationMarkerOpen(marker.path)
	markerHandle, err := directoryRoot.Open(marker.name)
	if err != nil {
		return nil, err
	}
	defer markerHandle.Close()
	openedInfo, openedError := markerHandle.Stat()
	afterOpen, pathError := directoryRoot.Lstat(marker.name)
	currentDirectory, directoryError := os.Lstat(absoluteDirectory)
	if openedError != nil || pathError != nil || directoryError != nil ||
		!openedInfo.Mode().IsRegular() || !afterOpen.Mode().IsRegular() ||
		!os.SameFile(marker.fileIdentity, openedInfo) || !os.SameFile(openedInfo, afterOpen) ||
		!currentDirectory.IsDir() || currentDirectory.Mode()&os.ModeSymlink != 0 || !os.SameFile(directoryInfo, currentDirectory) {
		return nil, fmt.Errorf("reservation marker changed while opening its rooted handle")
	}
	contents, err := io.ReadAll(markerHandle)
	if err != nil {
		return nil, err
	}
	afterReadHandle, handleError := markerHandle.Stat()
	afterReadPath, pathError := directoryRoot.Lstat(marker.name)
	currentDirectory, directoryError = os.Lstat(absoluteDirectory)
	if handleError != nil || pathError != nil || directoryError != nil ||
		!afterReadHandle.Mode().IsRegular() || !afterReadPath.Mode().IsRegular() ||
		!os.SameFile(openedInfo, afterReadHandle) || !os.SameFile(afterReadHandle, afterReadPath) ||
		openedInfo.Size() != afterReadHandle.Size() || !openedInfo.ModTime().Equal(afterReadHandle.ModTime()) ||
		!currentDirectory.IsDir() || currentDirectory.Mode()&os.ModeSymlink != 0 || !os.SameFile(directoryInfo, currentDirectory) {
		return nil, fmt.Errorf("reservation marker changed during read")
	}
	return contents, nil
}
