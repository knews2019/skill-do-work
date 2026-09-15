// Package rootedfs backports the os.Root operations do-work-cli uses that the standard
// library added only in Go 1.25: Rename, Link, MkdirAll, ReadFile and RemoveAll. The module floor is
// Go 1.24 so hosts whose toolchain cannot be upgraded (a pinned CI image, a proxy that
// blocks toolchain downloads) still build the command.
//
// MkdirAll, ReadFile and RemoveAll are composed entirely from Go 1.24 root methods and
// keep the full rooted guarantee. Rename and Link cannot be: the Go 1.25 methods run
// against the root's directory descriptor, so a path component swapped for a symlink
// mid-operation cannot escape the root, and the Go 1.24 standard library offers no
// portable way to do that (syscall exposes renameat on Linux only and linkat nowhere;
// golang.org/x/sys would add a module download on exactly the hosts that cannot fetch a
// toolchain). Those two resolve their paths through the root first — root.Lstat refuses
// any component that escapes it — and then perform the operating-system call on the
// joined absolute path. The residual window is between that check and the call, and
// using it needs write access inside the root. Delete this package and call the os.Root
// methods directly once go.mod returns to 1.25 or newer.
package rootedfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Rename renames oldPath to newPath, both relative to root. A symlink at oldPath is
// renamed itself, not followed, matching os.Rename and Root.Rename. The root must have
// been opened on an absolute path, or through OpenRoot below; see rootDirectory.
func Rename(root *os.Root, oldPath, newPath string) error {
	directory, err := rootDirectory(root)
	if err == nil {
		err = resolveSource(root, oldPath)
	}
	if err == nil {
		err = resolveTargetParent(root, newPath)
	}
	if err != nil {
		return &os.LinkError{Op: "renameat", Old: oldPath, New: newPath, Err: err}
	}
	return os.Rename(filepath.Join(directory, oldPath), filepath.Join(directory, newPath))
}

// Link creates newPath as a hard link to oldPath, both relative to root. It fails when
// newPath already exists, which is what makes it an exclusive publish. The root must have
// been opened on an absolute path, or through OpenRoot below; see rootDirectory.
func Link(root *os.Root, oldPath, newPath string) error {
	directory, err := rootDirectory(root)
	if err == nil {
		err = resolveSource(root, oldPath)
	}
	if err == nil {
		err = resolveTargetParent(root, newPath)
	}
	if err != nil {
		return &os.LinkError{Op: "linkat", Old: oldPath, New: newPath, Err: err}
	}
	return os.Link(filepath.Join(directory, oldPath), filepath.Join(directory, newPath))
}

// OpenRoot opens name as a nested root inside parent, for callers that will Rename or
// Link inside it. Older Go versions name a nested root only by the relative name
// passed to (*Root).OpenRoot, which rootDirectory may not resolve. This opens the
// nested root the guarded way first, proving name lies inside parent, then
// reopens the same directory by absolute path and refuses unless both opens landed on one
// directory.
func OpenRoot(parent *os.Root, name string) (*os.Root, error) {
	parentDirectory, err := rootDirectory(parent)
	if err != nil {
		return nil, &os.PathError{Op: "openat", Path: name, Err: err}
	}
	guarded, err := parent.OpenRoot(name)
	if err != nil {
		return nil, err
	}
	guardedInfo, err := guarded.Stat(".")
	guarded.Close()
	if err != nil {
		return nil, err
	}
	nested, err := os.OpenRoot(filepath.Join(parentDirectory, name))
	if err != nil {
		return nil, err
	}
	nestedInfo, err := nested.Stat(".")
	if err != nil || !os.SameFile(guardedInfo, nestedInfo) {
		nested.Close()
		if err == nil {
			err = errRootUnresolvable
		}
		return nil, &os.PathError{Op: "openat", Path: name, Err: err}
	}
	return nested, nil
}

// MkdirAll creates path and every missing parent inside root, all with perm. An existing
// directory anywhere on the way is fine; an existing non-directory is an error, as with
// os.MkdirAll. Every Mkdir goes through the root, so this helper keeps the full rooted
// guarantee.
func MkdirAll(root *os.Root, path string, perm os.FileMode) error {
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return nil
	}
	if filepath.IsAbs(cleaned) {
		return &os.PathError{Op: "mkdirat", Path: path, Err: errPathEscapesRoot}
	}
	prefix := ""
	for _, component := range strings.Split(cleaned, string(filepath.Separator)) {
		prefix = filepath.Join(prefix, component)
		err := root.Mkdir(prefix, perm)
		if err == nil || !errors.Is(err, fs.ErrExist) {
			if err != nil {
				return err
			}
			continue
		}
		info, statError := root.Stat(prefix)
		if statError != nil {
			return statError
		}
		if !info.IsDir() {
			return &os.PathError{Op: "mkdirat", Path: prefix, Err: errNotDirectory}
		}
	}
	return nil
}

// ReadFile returns the contents of path inside root, opened through the root.
func ReadFile(root *os.Root, path string) ([]byte, error) {
	file, err := root.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

// RemoveAll removes path and, when it is a directory, everything below it, every step
// through the root. Symlinks are removed, never followed. A missing path is not an error,
// as with os.RemoveAll.
func RemoveAll(root *os.Root, path string) error {
	info, err := root.Lstat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		directoryRoot, openError := root.OpenRoot(path)
		if openError != nil {
			return openError
		}
		defer directoryRoot.Close()
		openedInfo, statError := directoryRoot.Stat(".")
		if statError != nil {
			return statError
		}
		if !os.SameFile(info, openedInfo) {
			return fmt.Errorf("directory changed before recursive removal: %s", path)
		}
		directory, openError := directoryRoot.Open(".")
		if openError != nil {
			return openError
		}
		entries, readError := directory.ReadDir(-1)
		directory.Close()
		if readError != nil {
			return readError
		}
		if afterRemoveAllEnumeration != nil {
			afterRemoveAllEnumeration(path)
		}
		for _, entry := range entries {
			if removeError := RemoveAll(directoryRoot, entry.Name()); removeError != nil {
				return removeError
			}
		}
	}
	return root.Remove(path)
}

var afterRemoveAllEnumeration func(string)

var (
	errPathEscapesRoot  = errors.New("path escapes from parent")
	errNotDirectory     = errors.New("not a directory")
	errRootUnresolvable = errors.New("root directory could not be resolved to an absolute path; open it with os.OpenRoot on an absolute path or rootedfs.OpenRoot")
)

// rootDirectory returns the absolute directory a root was opened on. root.Name is the
// string given to os.OpenRoot, so a relative one is resolved against the working
// directory, and the result is accepted only when it is the very directory the root
// holds open: on older Go versions a nested root may carry only its relative name
// and otherwise resolve to an unrelated path.
func rootDirectory(root *os.Root) (string, error) {
	directory := root.Name()
	if !filepath.IsAbs(directory) {
		absolute, err := filepath.Abs(directory)
		if err != nil {
			return "", err
		}
		directory = absolute
	}
	rootInfo, err := root.Stat(".")
	if err != nil {
		return "", err
	}
	directoryInfo, err := os.Stat(directory)
	if err != nil || !os.SameFile(rootInfo, directoryInfo) {
		return "", errRootUnresolvable
	}
	return directory, nil
}

// resolveSource proves the source exists inside the root without following a final
// symlink; an intermediate component that escapes the root fails here.
func resolveSource(root *os.Root, path string) error {
	if path == "" || filepath.IsAbs(path) {
		return errPathEscapesRoot
	}
	_, err := root.Lstat(path)
	return err
}

// resolveTargetParent proves the target's parent is a directory inside the root, so the
// final component is created where the root says it lives.
func resolveTargetParent(root *os.Root, path string) error {
	if path == "" || filepath.IsAbs(path) {
		return errPathEscapesRoot
	}
	info, err := root.Stat(filepath.Dir(path))
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s: %w", filepath.Dir(path), errNotDirectory)
	}
	return nil
}
