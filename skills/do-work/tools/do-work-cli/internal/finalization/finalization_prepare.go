package finalization

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/publication"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
)

var (
	requestIDPattern = regexp.MustCompile(`^REQ-[0-9]+$`)
	digestPattern    = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
)

type requestBindingError struct {
	reason string
}

func (err requestBindingError) Error() string {
	return err.reason
}

func prepareJournal(ctx context.Context, repositoryRoot, manifestPath string) (*Journal, bool, error) {
	return prepareBoundJournal(ctx, repositoryRoot, manifestPath, "", "", "")
}

func prepareBoundJournal(ctx context.Context, repositoryRoot, manifestPath, expectedRequestID, expectedRequestPath, requiredTransition string) (*Journal, bool, error) {
	manifest, manifestBytes, err := decodeManifest(repositoryRoot, manifestPath)
	if err != nil {
		return nil, false, err
	}
	if expectedRequestID != "" && (manifest.RequestID != expectedRequestID || manifest.RequestPath != expectedRequestPath) {
		return nil, false, requestBindingError{reason: fmt.Sprintf("finalization manifest identifies %s at %s, expected %s at %s", manifest.RequestID, manifest.RequestPath, expectedRequestID, expectedRequestPath)}
	}
	if requiredTransition != "" && manifest.Transition != requiredTransition {
		return nil, false, fmt.Errorf("this lifecycle phase only permits finalization transition %q", requiredTransition)
	}
	if err := exec.Command("git", "-C", repositoryRoot, "diff", "--cached", "--quiet", "--exit-code").Run(); err != nil {
		return nil, false, fmt.Errorf("finalization requires an empty existing index")
	}
	journalPath, payloadDirectory, err := journalLocations(repositoryRoot, manifest.RequestID)
	if err != nil {
		return nil, false, err
	}
	manifestDigest := digestBytes(manifestBytes)
	if _, err := os.Lstat(journalPath); err == nil {
		journal, readError := readJournal(repositoryRoot, journalPath)
		if readError != nil {
			return nil, false, readError
		}
		if journal.ManifestSHA256 != manifestDigest {
			return nil, false, fmt.Errorf("an unfinished journal exists for %s with a different manifest", manifest.RequestID)
		}
		return journal, true, nil
	} else if !os.IsNotExist(err) {
		return nil, false, err
	}

	requestBytes, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(manifest.RequestPath)))
	if err != nil || digestBytes(requestBytes) != strings.ToLower(manifest.ExpectedRequestSHA256) {
		return nil, false, fmt.Errorf("request preimage does not match expected_request_sha256")
	}
	checkpointBytes, err := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "CHECKPOINT.md"))
	if err != nil || digestBytes(checkpointBytes) != strings.ToLower(manifest.ExpectedCheckpointSHA256) {
		return nil, false, fmt.Errorf("checkpoint preimage does not match expected_checkpoint_sha256")
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		return nil, false, err
	}
	completedAt, _ := time.Parse(time.RFC3339, manifest.CompletedAt)
	stateOptions := requeststate.StateOptions{
		RequestID: manifest.RequestID, RequestPath: manifest.RequestPath, WriterLabel: manifest.WriterLabel,
		Now: completedAt, ImplementationHash: manifest.ImplementationHash,
	}
	if manifest.Transition == "complete" {
		stateOptions.Transition = requeststate.TransitionComplete
		stateOptions.TerminalStatus = manifest.TerminalStatus
	} else {
		stateOptions.Transition = requeststate.TransitionFail
		stateOptions.FailureError = manifest.FailureError
		stateOptions.FailureType = manifest.FailureType
	}
	statePlan := requeststate.BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), stateOptions)
	if !statePlan.Runnable() {
		if statePlan.Refusal != nil {
			return nil, false, fmt.Errorf("lifecycle plan refused: %s: %s", statePlan.Refusal.Code, statePlan.Refusal.Reason)
		}
		return nil, false, fmt.Errorf("lifecycle plan is not runnable")
	}
	lifecyclePostimages, err := requeststate.PlannedPostimages(statePlan)
	if err != nil {
		return nil, false, err
	}
	lifecyclePreimages, err := snapshotImages(repositoryRoot, statePlan.TargetPaths)
	if err != nil {
		return nil, false, err
	}
	journalLifecyclePostimages := make([]FileImage, 0, len(lifecyclePostimages))
	for _, image := range lifecyclePostimages {
		journalLifecyclePostimages = append(journalLifecyclePostimages, FileImage{Path: image.Path, Exists: image.Exists, Bytes: image.Bytes, Mode: image.Mode})
	}

	var releaseManifest *publication.Manifest
	var releasePreimages, releasePostimages []FileImage
	if manifest.ReleaseManifestPath != "" {
		if err := os.MkdirAll(payloadDirectory, 0o700); err != nil {
			return nil, false, err
		}
		preparedRelease, prepareError := adoptReleaseManifest(repositoryRoot, manifest.ReleaseManifestPath, payloadDirectory)
		if prepareError != nil {
			_ = os.RemoveAll(payloadDirectory)
			return nil, false, prepareError
		}
		releaseManifest = &preparedRelease
		releasePlan := publication.BuildReleasePlan(repositoryRoot, preparedRelease)
		if !releasePlan.Runnable() {
			_ = os.RemoveAll(payloadDirectory)
			if releasePlan.Refusal != nil {
				return nil, false, fmt.Errorf("release plan refused: %s: %s", releasePlan.Refusal.Code, releasePlan.Refusal.Reason)
			}
			return nil, false, fmt.Errorf("release plan is not runnable")
		}
		releasePreimages, err = snapshotImages(repositoryRoot, releasePlan.TargetPaths)
		if err != nil {
			_ = os.RemoveAll(payloadDirectory)
			return nil, false, err
		}
		releasePostimages = publicationPostimages(releasePlan, releasePreimages)
		archiveBefore, archiveAfter, stampError := releaseStampImages(statePlan.DestinationPath, journalLifecyclePostimages, manifest.ReleaseAt)
		if stampError != nil {
			_ = os.RemoveAll(payloadDirectory)
			return nil, false, stampError
		}
		releasePreimages = append(releasePreimages, archiveBefore)
		releasePostimages = append(releasePostimages, archiveAfter)
		releasePreimages = sortedUniqueImages(releasePreimages)
		releasePostimages = sortedUniqueImages(releasePostimages)
	}

	effectiveCommitPaths, err := normalizeRepositoryPaths(manifest.CommitPaths)
	if err != nil {
		_ = os.RemoveAll(payloadDirectory)
		return nil, false, err
	}
	requiredCommitPaths := append([]string(nil), statePlan.TargetPaths...)
	for _, image := range releasePostimages {
		requiredCommitPaths = append(requiredCommitPaths, image.Path)
	}
	missing := subtractPaths(uniqueSorted(requiredCommitPaths), effectiveCommitPaths)
	if len(missing) > 0 {
		_ = os.RemoveAll(payloadDirectory)
		return nil, false, fmt.Errorf("commit_paths omits planned lifecycle or release targets: %s", strings.Join(missing, ", "))
	}

	now := time.Now().UTC().Truncate(time.Second)
	journal := &Journal{
		Version: journalVersion, CreatedAt: now, UpdatedAt: now, Phase: PhasePrepared,
		Manifest: manifest, ManifestSHA256: manifestDigest, JournalPath: journalPath,
		ArchivedPath: statePlan.DestinationPath, LifecyclePreimages: lifecyclePreimages,
		LifecyclePostimages: journalLifecyclePostimages, ReleaseManifest: releaseManifest,
		ReleasePreimages: releasePreimages, ReleasePostimages: releasePostimages,
		EffectiveCommitPaths: effectiveCommitPaths,
	}
	if releaseManifest != nil {
		journal.PayloadDirectory = payloadDirectory
	}
	if err := writeJournal(journal); err != nil {
		_ = os.RemoveAll(payloadDirectory)
		return nil, false, err
	}
	return journal, false, nil
}

func releaseStampImages(archivedPath string, lifecyclePostimages []FileImage, releaseAt string) (FileImage, FileImage, error) {
	for _, image := range lifecyclePostimages {
		if image.Path != archivedPath || !image.Exists {
			continue
		}
		document, err := requestmodel.ParseDocument(image.Bytes)
		if err != nil {
			return FileImage{}, FileImage{}, err
		}
		if err := document.SetScalar("release_at", releaseAt); err != nil {
			return FileImage{}, FileImage{}, err
		}
		after := image
		after.Bytes = document.DocumentBytes()
		return image, after, nil
	}
	return FileImage{}, FileImage{}, fmt.Errorf("release_at archive postimage is missing")
}

func sortedUniqueImages(images []FileImage) []FileImage {
	byPath := imagesByPath(images)
	paths := make([]string, 0, len(byPath))
	for path := range byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	result := make([]FileImage, 0, len(paths))
	for _, path := range paths {
		result = append(result, byPath[path])
	}
	return result
}

func snapshotImages(repositoryRoot string, paths []string) ([]FileImage, error) {
	normalized, err := normalizeRepositoryPaths(paths)
	if err != nil {
		return nil, err
	}
	images := make([]FileImage, 0, len(normalized))
	for _, path := range normalized {
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		info, err := os.Lstat(absolutePath)
		if os.IsNotExist(err) {
			images = append(images, FileImage{Path: path})
			continue
		}
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("target preimage is not a regular file: %s", path)
		}
		contents, err := os.ReadFile(absolutePath)
		if err != nil {
			return nil, err
		}
		images = append(images, FileImage{Path: path, Exists: true, Bytes: contents, Mode: uint32(info.Mode())})
	}
	return images, nil
}

func publicationPostimages(plan publication.PublicationPlan, preimageSets ...[]FileImage) []FileImage {
	images := map[string]FileImage{}
	preimages := map[string]FileImage{}
	if len(preimageSets) > 0 {
		preimages = imagesByPath(preimageSets[0])
	}
	for _, mutation := range plan.Mutations {
		switch mutation.Kind {
		case publication.MutationCreate, publication.MutationReplace:
			mode := uint32(mutation.Mode)
			if mutation.Kind == publication.MutationReplace && mode == 0 {
				mode = preimages[mutation.Path].Mode
			}
			images[mutation.Path] = FileImage{Path: mutation.Path, Exists: true, Bytes: append([]byte(nil), mutation.Contents...), Mode: mode}
		case publication.MutationMove:
			images[mutation.Path] = FileImage{Path: mutation.Path}
			contents := mutation.ExpectedBytes
			if len(mutation.Contents) > 0 {
				contents = mutation.Contents
			}
			images[mutation.DestinationPath] = FileImage{Path: mutation.DestinationPath, Exists: true, Bytes: append([]byte(nil), contents...), Mode: uint32(mutation.Mode)}
		}
	}
	paths := make([]string, 0, len(images))
	for path := range images {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	result := make([]FileImage, 0, len(paths))
	for _, path := range paths {
		result = append(result, images[path])
	}
	return result
}

func adoptReleaseManifest(repositoryRoot, manifestPath, payloadDirectory string) (publication.Manifest, error) {
	absolutePath, err := containedOrAbsolute(repositoryRoot, manifestPath)
	if err != nil {
		return publication.Manifest{}, err
	}
	file, err := os.Open(absolutePath)
	if err != nil {
		return publication.Manifest{}, err
	}
	manifest, err := publication.DecodeManifest(file, publication.OperationRelease)
	_ = file.Close()
	if err != nil {
		return publication.Manifest{}, err
	}
	manifest.CommitMessage = ""
	counter := 0
	copyPayload := func(payload *publication.PayloadFile) error {
		if payload == nil || payload.SourcePath == "" {
			return nil
		}
		sourcePath, err := containedOrAbsolute(repositoryRoot, payload.SourcePath)
		if err != nil {
			return err
		}
		info, err := os.Lstat(sourcePath)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("release payload is not a regular non-symlink file: %s", payload.SourcePath)
		}
		contents, err := os.ReadFile(sourcePath)
		if err != nil {
			return err
		}
		if payload.SHA256 != "" && !strings.EqualFold(payload.SHA256, digestBytes(contents)) {
			return fmt.Errorf("release payload digest is stale: %s", payload.SourcePath)
		}
		counter++
		destination := filepath.Join(payloadDirectory, fmt.Sprintf("payload-%03d", counter))
		if err := os.WriteFile(destination, contents, 0o600); err != nil {
			return err
		}
		payload.SourcePath = destination
		payload.SHA256 = digestBytes(contents)
		return nil
	}
	for index := range manifest.Release.Targets {
		if err := copyPayload(&manifest.Release.Targets[index].ExpectedPayload); err != nil {
			return publication.Manifest{}, err
		}
		if err := copyPayload(&manifest.Release.Targets[index].NewPayload); err != nil {
			return publication.Manifest{}, err
		}
	}
	for index := range manifest.Release.Changelogs {
		if !manifest.Release.Changelogs[index].Create {
			if err := copyPayload(&manifest.Release.Changelogs[index].ExpectedPayload); err != nil {
				return publication.Manifest{}, err
			}
		}
		if err := copyPayload(&manifest.Release.Changelogs[index].NewPayload); err != nil {
			return publication.Manifest{}, err
		}
	}
	return manifest, nil
}

func normalizeRepositoryPaths(paths []string) ([]string, error) {
	set := map[string]bool{}
	for _, path := range paths {
		if path == "" || filepath.IsAbs(path) {
			return nil, fmt.Errorf("path must be non-empty and repository-relative: %s", path)
		}
		cleaned := filepath.Clean(path)
		if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("path escapes repository: %s", path)
		}
		set[filepath.ToSlash(cleaned)] = true
	}
	result := make([]string, 0, len(set))
	for path := range set {
		result = append(result, path)
	}
	sort.Strings(result)
	return result, nil
}

func subtractPaths(left, right []string) []string {
	set := map[string]bool{}
	for _, path := range right {
		set[path] = true
	}
	result := []string{}
	for _, path := range left {
		if !set[path] {
			result = append(result, path)
		}
	}
	return result
}

func uniqueSorted(paths []string) []string {
	result, _ := normalizeRepositoryPaths(paths)
	return result
}
