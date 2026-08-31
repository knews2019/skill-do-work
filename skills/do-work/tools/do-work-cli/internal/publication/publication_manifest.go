package publication

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func DecodeManifest(reader io.Reader, expectedOperation OperationName) (Manifest, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return manifest, fmt.Errorf("decode publication manifest: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return manifest, fmt.Errorf("publication manifest must contain one JSON object")
	}
	if manifest.Operation != expectedOperation {
		return manifest, fmt.Errorf("manifest operation %q does not match command %q", manifest.Operation, expectedOperation)
	}
	bodyCount := 0
	if manifest.Capture != nil {
		bodyCount++
	}
	if manifest.Answer != nil {
		bodyCount++
	}
	if manifest.Release != nil {
		bodyCount++
	}
	if bodyCount != 1 {
		return manifest, fmt.Errorf("publication manifest requires exactly one operation body")
	}
	if expectedOperation == OperationCaptureFiles && manifest.Capture == nil || expectedOperation == OperationAnswer && manifest.Answer == nil || expectedOperation == OperationRelease && manifest.Release == nil {
		return manifest, fmt.Errorf("publication manifest body does not match operation %q", expectedOperation)
	}
	return manifest, nil
}

func containedPath(path string) (string, error) {
	if path == "" || filepath.IsAbs(path) {
		return "", fmt.Errorf("path %q must be non-empty and repository-relative", path)
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes or names the repository root", path)
	}
	return filepath.ToSlash(cleaned), nil
}

func readPayload(repositoryRoot string, payload PayloadFile) ([]byte, os.FileMode, error) {
	if payload.SourcePath == "" {
		return nil, 0, fmt.Errorf("payload source path is required")
	}
	path := payload.SourcePath
	absolutePath := path
	if !filepath.IsAbs(path) {
		var err error
		path, err = containedPath(path)
		if err != nil {
			return nil, 0, err
		}
		absolutePath = filepath.Join(repositoryRoot, filepath.FromSlash(path))
	}
	info, err := os.Lstat(absolutePath)
	if err != nil {
		return nil, 0, fmt.Errorf("inspect payload %s: %w", path, err)
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, 0, fmt.Errorf("payload %s must be a regular non-symlink file", path)
	}
	contents, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, 0, err
	}
	if payload.SHA256 != "" {
		digest := sha256.Sum256(contents)
		if !strings.EqualFold(payload.SHA256, hex.EncodeToString(digest[:])) {
			return nil, 0, fmt.Errorf("payload %s digest is stale", path)
		}
	}
	return contents, info.Mode(), nil
}

func validateOutsideBytes(contents []byte) error {
	for _, value := range contents {
		if value == 0x7f || value < 0x20 && value != '\n' && value != '\t' {
			return fmt.Errorf("outside text contains unsupported control byte 0x%02x", value)
		}
	}
	return nil
}

func containedOutsideBytes(contents []byte, lineEnding string) []byte {
	longestRun, currentRun := 0, 0
	for _, value := range contents {
		if value == '`' {
			currentRun++
			if currentRun > longestRun {
				longestRun = currentRun
			}
		} else {
			currentRun = 0
		}
	}
	fenceLength := longestRun + 1
	if fenceLength < 3 {
		fenceLength = 3
	}
	fence := bytes.Repeat([]byte("`"), fenceLength)
	normalized := bytes.ReplaceAll(contents, []byte("\r\n"), []byte("\n"))
	lines := bytes.Split(normalized, []byte("\n"))
	var result bytes.Buffer
	result.WriteString("> ")
	result.Write(fence)
	result.WriteString(lineEnding)
	for index, line := range lines {
		if index == len(lines)-1 && len(line) == 0 {
			break
		}
		result.WriteString("> ")
		result.Write(line)
		result.WriteString(lineEnding)
	}
	result.WriteString("> ")
	result.Write(fence)
	return result.Bytes()
}

func finalizePlan(plan PublicationPlan) PublicationPlan {
	sort.SliceStable(plan.Mutations, func(first, second int) bool {
		firstPath, secondPath := plan.Mutations[first].Path, plan.Mutations[second].Path
		if plan.Mutations[first].Kind == MutationMove {
			firstPath = plan.Mutations[first].DestinationPath
		}
		if plan.Mutations[second].Kind == MutationMove {
			secondPath = plan.Mutations[second].DestinationPath
		}
		return firstPath < secondPath
	})
	seen := map[string]bool{}
	targets := []string{}
	untracked := []string{}
	for _, mutation := range plan.Mutations {
		paths := []string{mutation.Path}
		if mutation.Kind == MutationMove {
			paths = append(paths, mutation.DestinationPath)
		}
		for _, rawPath := range paths {
			path, err := containedPath(rawPath)
			if err != nil {
				return refusedPlan(plan, "PUBLICATION-PATH-UNSAFE", err.Error(), nil, rawPath)
			}
			if seen[path] {
				return refusedPlan(plan, "PUBLICATION-DUPLICATE-TARGET", "one path has multiple planned operations", nil, path)
			}
			seen[path] = true
			targets = append(targets, path)
		}
		if mutation.AllowUntracked {
			untracked = append(untracked, mutation.Path)
		}
	}
	sort.Strings(targets)
	sort.Strings(untracked)
	plan.TargetPaths = targets
	plan.ExistingUntrackedTargetPaths = untracked
	plan.Changes = make([]resultmodel.RecordedChange, 0, len(plan.Mutations))
	for _, mutation := range plan.Mutations {
		path, kind := mutation.Path, string(mutation.Kind)
		if mutation.Kind == MutationMove {
			path = mutation.DestinationPath
		}
		plan.Changes = append(plan.Changes, resultmodel.RecordedChange{Path: path, Kind: kind, Detail: "planned"})
	}
	sort.Slice(plan.Changes, func(i, j int) bool { return plan.Changes[i].Path < plan.Changes[j].Path })
	return plan
}

func refusedPlan(plan PublicationPlan, code, reason string, ids []string, paths ...string) PublicationPlan {
	plan.Refusal = &Refusal{Code: code, Reason: reason, IDs: ids, Paths: paths}
	return plan
}
