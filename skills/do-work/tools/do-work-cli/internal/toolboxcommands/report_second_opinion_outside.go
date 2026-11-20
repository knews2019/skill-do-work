package toolboxcommands

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type secondOpinionOutsideFile struct {
	mode      os.FileMode
	signature string
}

// Observe the isolated checkout itself, including ignored paths. Git's mutable
// HEAD, index and ignore rules cannot prove what the backend changed. The root
// .git control entry is excluded; no repository links are followed.
func secondOpinionOutsideSnapshot(ctx context.Context, repository, bundle string) (map[string]secondOpinionOutsideFile, error) {
	files := map[string]secondOpinionOutsideFile{}
	err := filepath.WalkDir(repository, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if path == bundle || path == filepath.Join(repository, ".git") {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if path == repository {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := secondOpinionOutsideFile{mode: info.Mode()}
		switch {
		case info.Mode().IsRegular():
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			opened, err := file.Stat()
			if err != nil {
				return err
			}
			if !os.SameFile(info, opened) {
				return fmt.Errorf("outside path changed while reading: %s", path)
			}
			digest := sha256.New()
			if _, err := io.Copy(digest, secondOpinionContextReader{ctx, file}); err != nil {
				return fmt.Errorf("read outside path %s: %w", path, err)
			}
			item.signature = fmt.Sprintf("%x", digest.Sum(nil))
		case info.Mode()&os.ModeSymlink != 0:
			item.signature, err = os.Readlink(path)
			if err != nil {
				return err
			}
			// Directories and special files need only type/mode evidence. Never
			// open a device or FIFO just to check for an outside write.
		}
		relative, err := filepath.Rel(repository, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(relative)] = item
		return nil
	})
	return files, err
}

type secondOpinionContextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (reader secondOpinionContextReader) Read(buffer []byte) (int, error) {
	if err := reader.ctx.Err(); err != nil {
		return 0, err
	}
	return reader.reader.Read(buffer)
}

func secondOpinionOutsideCheck(ctx context.Context, repository, bundle string, before map[string]secondOpinionOutsideFile, row *resultmodel.AIReportSecondOpinionRun) error {
	after, err := secondOpinionOutsideSnapshot(ctx, repository, bundle)
	if err != nil {
		return err
	}
	paths := []string{}
	for name, old := range before {
		if current, ok := after[name]; !ok || current != old {
			paths = append(paths, name)
		}
	}
	for name := range after {
		if _, existed := before[name]; !existed {
			paths = append(paths, name)
		}
	}
	sort.Strings(paths)
	row.OutsidePaths = paths
	if len(paths) > 0 {
		return fmt.Errorf("outside-bundle changes: %s", strings.Join(paths, ", "))
	}
	return nil
}
