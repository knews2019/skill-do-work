package requeststate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

func ResolveTarget(snapshot *repositorymodel.RepositorySnapshot, requestID, requestPath string) (*repositorymodel.RequestFile, *StateRefusal) {
	if snapshot == nil {
		return nil, refuse("STATE-DISCOVERY-FAILED", "repository snapshot is required")
	}
	if requestID == "" {
		return nil, refuse("STATE-USAGE", "one REQ id is required")
	}
	candidates := snapshot.RequestsByID[requestID]
	if len(candidates) != 1 {
		code := "REQUEST-NOT-FOUND"
		if len(candidates) > 1 {
			code = "REQUEST-AMBIGUOUS"
		}
		return nil, refuse(code, fmt.Sprintf("%s resolves to %d repository records", requestID, len(candidates)))
	}
	target := candidates[0]
	if target.ParseFailure != "" || target.TypedRecord.RequestID != requestID || target.FilenameID != requestID {
		return nil, refuse("REQUEST-IDENTITY-MISMATCH", "filename, frontmatter id, and parsed record must agree", requestPathFor(target))
	}
	if requestPath != "" {
		cleaned := filepath.Clean(requestPath)
		if filepath.IsAbs(requestPath) || cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
			return nil, refuse("REQUEST-PATH-UNSAFE", "request path must be contained and repository-relative", requestPath)
		}
		cleaned = filepath.ToSlash(cleaned)
		if cleaned != requestPathFor(target) {
			return nil, refuse("REQUEST-SNAPSHOT-STALE", fmt.Sprintf("returned path %s no longer identifies %s", requestPath, requestID), requestPath)
		}
	}
	return target, nil
}

func requestPathFor(requestFile *repositorymodel.RequestFile) string {
	if requestFile == nil || requestFile.RelativePath == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(requestFile.RelativePath)))
}

func refuse(code, reason string, paths ...string) *StateRefusal {
	return &StateRefusal{Code: code, Reason: reason, Paths: paths}
}
