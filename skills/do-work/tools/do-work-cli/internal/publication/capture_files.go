package publication

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
)

var (
	userRequestPattern = regexp.MustCompile(`^UR-[0-9]+$`)
	requestPattern     = regexp.MustCompile(`^REQ-[0-9]+$`)
)

func BuildCapturePlan(repositoryRoot string, manifest Manifest) PublicationPlan {
	plan := PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: repositoryRoot, CommitMessage: manifest.CommitMessage}
	capture := manifest.Capture
	if capture == nil {
		return refusedPlan(plan, "CAPTURE-MANIFEST-MISSING", "capture body is required", nil)
	}
	if !userRequestPattern.MatchString(capture.UserRequestID) {
		return refusedPlan(plan, "CAPTURE-UR-ID-INVALID", "user request id is not canonical", []string{capture.UserRequestID})
	}
	urPath, err := containedPath(capture.UserRequest.Path)
	if err != nil {
		return refusedPlan(plan, "CAPTURE-PATH-UNSAFE", err.Error(), []string{capture.UserRequestID}, capture.UserRequest.Path)
	}
	if !strings.Contains(filepath.Base(urPath), "input.md") || !strings.Contains(urPath, "/"+capture.UserRequestID+"/") {
		return refusedPlan(plan, "CAPTURE-UR-PATH-MISMATCH", "UR input path must name its UR directory and input.md", []string{capture.UserRequestID}, urPath)
	}
	urBytes, urMode, err := readPayload(repositoryRoot, capture.UserRequest.Payload)
	if err != nil {
		return refusedPlan(plan, "CAPTURE-PAYLOAD-INVALID", err.Error(), []string{capture.UserRequestID}, capture.UserRequest.Payload.SourcePath)
	}
	urDocument, err := requestmodel.ParseDocument(urBytes)
	if err != nil {
		return refusedPlan(plan, "CAPTURE-UR-INVALID", err.Error(), []string{capture.UserRequestID}, urPath)
	}
	if urDocument.TypedRecord().RequestID != capture.UserRequestID {
		return refusedPlan(plan, "CAPTURE-UR-ID-MISMATCH", "UR frontmatter id does not match manifest", []string{capture.UserRequestID}, urPath)
	}
	requestMembership := map[string]bool{}
	if field, found := urDocument.FieldValue("requests"); found {
		for _, requestID := range field.ListValues {
			requestMembership[requestID] = true
		}
	}
	if capture.RawInput != nil {
		rawBytes, _, rawError := readPayload(repositoryRoot, *capture.RawInput)
		if rawError != nil {
			return refusedPlan(plan, "CAPTURE-RAW-INPUT-INVALID", rawError.Error(), []string{capture.UserRequestID}, capture.RawInput.SourcePath)
		}
		if controlError := validateOutsideBytes(rawBytes); controlError != nil {
			return refusedPlan(plan, "CAPTURE-RAW-INPUT-UNSAFE", controlError.Error(), []string{capture.UserRequestID}, capture.RawInput.SourcePath)
		}
		lineEnding := "\n"
		if bytes.Contains(urBytes, []byte("\r\n")) {
			lineEnding = "\r\n"
		}
		if !bytes.Contains(urBytes, containedOutsideBytes(rawBytes, lineEnding)) {
			return refusedPlan(plan, "CAPTURE-RAW-INPUT-NOT-CONTAINED", "UR payload does not contain the canonical byte-derived outside-text block", []string{capture.UserRequestID}, urPath)
		}
	}
	plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationCreate, Path: urPath, Contents: urBytes, Mode: selectedMode(capture.UserRequest.Mode, urMode)})
	seenIDs := map[string]bool{}
	for _, request := range capture.Requests {
		if !requestPattern.MatchString(request.ID) || request.UserRequestID != capture.UserRequestID || seenIDs[request.ID] {
			return refusedPlan(plan, "CAPTURE-REQ-LINKAGE-INVALID", "REQ ids must be unique, canonical, and linked to the manifest UR", []string{request.ID, capture.UserRequestID}, request.File.Path)
		}
		seenIDs[request.ID] = true
		requestPath, pathError := containedPath(request.File.Path)
		if pathError != nil || !strings.HasPrefix(filepath.Base(requestPath), request.ID+"-") {
			return refusedPlan(plan, "CAPTURE-REQ-PATH-MISMATCH", "REQ filename must begin with its exact id", []string{request.ID}, request.File.Path)
		}
		if !requestMembership[request.ID] {
			return refusedPlan(plan, "CAPTURE-UR-MEMBERSHIP-MISSING", "UR requests field does not contain REQ id", []string{request.ID, capture.UserRequestID}, urPath)
		}
		requestBytes, requestMode, payloadError := readPayload(repositoryRoot, request.File.Payload)
		if payloadError != nil {
			return refusedPlan(plan, "CAPTURE-PAYLOAD-INVALID", payloadError.Error(), []string{request.ID}, request.File.Payload.SourcePath)
		}
		document, parseError := requestmodel.ParseDocument(requestBytes)
		if parseError != nil {
			return refusedPlan(plan, "CAPTURE-REQ-INVALID", parseError.Error(), []string{request.ID}, requestPath)
		}
		record := document.TypedRecord()
		if record.RequestID != request.ID || record.UserRequestID != capture.UserRequestID {
			return refusedPlan(plan, "CAPTURE-REQ-LINKAGE-INVALID", "REQ frontmatter id and user_request must match manifest", []string{request.ID, capture.UserRequestID}, requestPath)
		}
		reservationPath, reservationError := containedPath(request.ReservationPath)
		if reservationError != nil || reservationPath != "do-work/.req-reservations/"+request.ID {
			return refusedPlan(plan, "CAPTURE-RESERVATION-MISMATCH", "reservation path must exactly match the REQ id", []string{request.ID}, request.ReservationPath)
		}
		plan.Mutations = append(plan.Mutations,
			PlannedMutation{Kind: MutationCreate, Path: reservationPath, Contents: []byte(request.ID + "\n"), Mode: 0o644},
			PlannedMutation{Kind: MutationCreate, Path: requestPath, Contents: requestBytes, Mode: selectedMode(request.File.Mode, requestMode)},
		)
	}
	for _, asset := range capture.Assets {
		assetPath, pathError := containedPath(asset.Path)
		if pathError != nil || !strings.HasPrefix(assetPath, "do-work/user-requests/"+capture.UserRequestID+"/assets/") {
			return refusedPlan(plan, "CAPTURE-ASSET-PATH-INVALID", "asset destination must be inside the exact UR assets directory", []string{capture.UserRequestID}, asset.Path)
		}
		assetBytes, assetMode, payloadError := readPayload(repositoryRoot, asset.Payload)
		if payloadError != nil {
			return refusedPlan(plan, "CAPTURE-ASSET-SOURCE-INVALID", payloadError.Error(), []string{capture.UserRequestID}, asset.Payload.SourcePath)
		}
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationCreate, Path: assetPath, Contents: assetBytes, Mode: selectedMode(asset.Mode, assetMode)})
	}
	for _, fold := range capture.Folds {
		foldPath, pathError := containedPath(fold.Path)
		if pathError != nil || !strings.HasPrefix(foldPath, "do-work/") {
			return refusedPlan(plan, "CAPTURE-FOLD-PATH-INVALID", "fold target must be contained in do-work", nil, fold.Path)
		}
		expectedBytes, _, expectedError := readPayload(repositoryRoot, fold.ExpectedPayload)
		newBytes, _, newError := readPayload(repositoryRoot, fold.NewPayload)
		if expectedError != nil || newError != nil {
			return refusedPlan(plan, "CAPTURE-FOLD-PAYLOAD-INVALID", firstError(expectedError, newError).Error(), nil, foldPath)
		}
		currentBytes, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(foldPath)))
		if readError != nil || !bytes.Equal(currentBytes, expectedBytes) {
			return refusedPlan(plan, "CAPTURE-FOLD-STALE", "fold target does not match expected bytes", nil, foldPath)
		}
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationReplace, Path: foldPath, ExpectedBytes: expectedBytes, Contents: newBytes, AllowUntracked: fold.AllowUntracked})
	}
	plan = finalizePlan(plan)
	if plan.Refusal != nil {
		return plan
	}
	for _, mutation := range plan.Mutations {
		if mutation.Kind == MutationCreate {
			if _, statError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(mutation.Path))); statError == nil {
				return refusedPlan(plan, "CAPTURE-COLLISION", "capture destination already exists", nil, mutation.Path)
			} else if !os.IsNotExist(statError) {
				return refusedPlan(plan, "CAPTURE-INSPECTION-FAILED", statError.Error(), nil, mutation.Path)
			}
		}
	}
	directories, topologyError := planCreatedDirectories(repositoryRoot, plan.TargetPaths)
	if topologyError != nil {
		return refusedPlan(plan, "CAPTURE-TOPOLOGY-UNSAFE", topologyError.Error(), nil, "do-work")
	}
	plan.CreatedDirectoryPaths = directories
	return plan
}

func selectedMode(manifestMode uint32, sourceMode os.FileMode) os.FileMode {
	if manifestMode != 0 {
		return os.FileMode(manifestMode)
	}
	return sourceMode
}

func firstError(first, second error) error {
	if first != nil {
		return first
	}
	return second
}

func planCreatedDirectories(repositoryRoot string, targetPaths []string) ([]string, error) {
	directorySet := map[string]bool{}
	for _, targetPath := range targetPaths {
		parent := filepath.Dir(filepath.FromSlash(targetPath))
		parts := strings.Split(filepath.ToSlash(parent), "/")
		current := ""
		missingParent := false
		for _, part := range parts {
			if part == "." || part == "" {
				continue
			}
			current = filepath.ToSlash(filepath.Join(current, part))
			if missingParent {
				directorySet[current] = true
				continue
			}
			info, err := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(current)))
			if os.IsNotExist(err) {
				missingParent = true
				directorySet[current] = true
				continue
			}
			if err != nil {
				return nil, err
			}
			if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("path component %s is not a real directory", current)
			}
		}
	}
	directories := make([]string, 0, len(directorySet))
	for path := range directorySet {
		directories = append(directories, path)
	}
	sort.Slice(directories, func(i, j int) bool {
		firstDepth, secondDepth := strings.Count(directories[i], "/"), strings.Count(directories[j], "/")
		if firstDepth == secondDepth {
			return directories[i] < directories[j]
		}
		return firstDepth < secondDepth
	})
	return directories, nil
}
