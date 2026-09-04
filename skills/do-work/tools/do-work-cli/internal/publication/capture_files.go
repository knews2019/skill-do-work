package publication

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
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
	manifestRequestIDs := make([]string, 0, len(capture.Requests))
	seenManifestIDs := map[string]bool{}
	for _, request := range capture.Requests {
		if !requestPattern.MatchString(request.ID) || request.UserRequestID != capture.UserRequestID || seenManifestIDs[request.ID] {
			return refusedPlan(plan, "CAPTURE-REQ-LINKAGE-INVALID", "REQ ids must be unique, canonical, and linked to the manifest UR", []string{request.ID, capture.UserRequestID}, request.File.Path)
		}
		seenManifestIDs[request.ID] = true
		manifestRequestIDs = append(manifestRequestIDs, request.ID)
	}
	urPath, err := containedPath(capture.UserRequest.Path)
	if err != nil {
		return refusedPlan(plan, "CAPTURE-PATH-UNSAFE", err.Error(), []string{capture.UserRequestID}, capture.UserRequest.Path)
	}
	wantedURPath := "do-work/user-requests/" + capture.UserRequestID + "/input.md"
	if urPath != wantedURPath {
		return refusedPlan(plan, "CAPTURE-UR-PATH-MISMATCH", "UR input path must exactly match the canonical active UR input path", []string{capture.UserRequestID}, urPath)
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
	if requestsField, found := urDocument.FieldValue("requests"); found && requestsField.ListValues != nil {
		for _, requestID := range manifestRequestIDs {
			if !requestMembership[requestID] {
				return refusedPlan(plan, "CAPTURE-UR-MEMBERSHIP-MISSING", "UR requests field does not contain REQ id", []string{requestID, capture.UserRequestID}, urPath)
			}
		}
	}
	if schemaError := validateCaptureURSchema(urDocument, manifestRequestIDs); schemaError != nil {
		return refusedPlan(plan, "CAPTURE-UR-SCHEMA-INVALID", schemaError.Error(), []string{capture.UserRequestID}, urPath)
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
	var requestMutations []PlannedMutation
	for _, request := range capture.Requests {
		requestPath, pathError := containedPath(request.File.Path)
		requestBase := filepath.Base(requestPath)
		if pathError != nil || filepath.ToSlash(filepath.Dir(requestPath)) != "do-work/queue" || !strings.HasPrefix(requestBase, request.ID+"-") || !strings.HasSuffix(requestBase, ".md") || requestBase == request.ID+"-.md" {
			return refusedPlan(plan, "CAPTURE-REQ-PATH-MISMATCH", "REQ path must exactly name a slugged Markdown file in do-work/queue", []string{request.ID}, request.File.Path)
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
		if schemaError := validateCaptureREQSchema(document, record); schemaError != nil {
			return refusedPlan(plan, "CAPTURE-REQ-SCHEMA-INVALID", schemaError.Error(), []string{request.ID}, requestPath)
		}
		reservationPath, reservationError := containedPath(request.ReservationPath)
		canonicalPath, canonicalError := canonicalReservationPath(request.ID)
		if reservationError != nil || canonicalError != nil || reservationPath != canonicalPath {
			return refusedPlan(plan, "CAPTURE-RESERVATION-MISMATCH", "reservation path must use the canonical stored-REQ-id spelling: "+canonicalPath, []string{request.ID}, request.ReservationPath)
		}
		existingReservations, inspectionError := existingReservationMarkers(repositoryRoot, request.ID)
		if inspectionError != nil {
			return refusedPlan(plan, "CAPTURE-INSPECTION-FAILED", inspectionError.Error(), []string{request.ID}, publicationReservationDirectory)
		}
		if len(existingReservations) > 0 {
			return refusedPlan(plan, "CAPTURE-COLLISION", "reservation number already exists under a canonical or legacy spelling", []string{request.ID}, reservationMarkerPaths(existingReservations)...)
		}
		requestMutations = append(requestMutations,
			PlannedMutation{Kind: MutationCreate, Path: reservationPath, Contents: []byte(request.ID + "\n"), Mode: 0o644},
			PlannedMutation{Kind: MutationCreate, Path: requestPath, Contents: requestBytes, Mode: selectedMode(request.File.Mode, requestMode)},
		)
	}
	plan.Mutations = append([]PlannedMutation{{Kind: MutationCreate, Path: urPath, Contents: urBytes, Mode: selectedMode(capture.UserRequest.Mode, urMode)}}, requestMutations...)
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

func validateCaptureURSchema(document *requestmodel.RequestDocument, manifestRequestIDs []string) error {
	if duplicateField := duplicateCaptureField(document); duplicateField != "" {
		return fmt.Errorf("duplicate top-level field %s", duplicateField)
	}
	if aliasField, canonicalField := captureAliasField(document); aliasField != "" {
		return fmt.Errorf("read-only alias %s must be authored as %s", aliasField, canonicalField)
	}
	for _, fieldName := range []string{"id", "title", "created_at", "word_count"} {
		if field, found := document.FieldValue(fieldName); !found || field.ScalarValue == "" || field.ListValues != nil || field.NestedValues != nil {
			return fmt.Errorf("required scalar field %s is missing or malformed", fieldName)
		}
	}
	requestsField, found := document.FieldValue("requests")
	if !found || requestsField.ListValues == nil {
		return fmt.Errorf("required field requests must be a list")
	}
	if !exactUniqueStrings(requestsField.ListValues, manifestRequestIDs) {
		return fmt.Errorf("requests must exactly match the duplicate-free manifest REQ list")
	}
	if !isCanonicalCaptureTimestamp(document, "created_at") {
		return fmt.Errorf("created_at must be a whole-second UTC timestamp")
	}
	if !isSafelyAuthoredUserText(document, "title") {
		return fmt.Errorf("title must use a single-quoted or literal-block scalar")
	}
	wordCount, conversionError := strconv.Atoi(document.TypedRecord().FieldEvidenceByName["word_count"].ScalarValue)
	if conversionError != nil || wordCount < 0 {
		return fmt.Errorf("word_count must be a non-negative integer")
	}
	return nil
}

func validateCaptureREQSchema(document *requestmodel.RequestDocument, record requestmodel.RequestRecord) error {
	if duplicateField := duplicateCaptureField(document); duplicateField != "" {
		return fmt.Errorf("duplicate top-level field %s", duplicateField)
	}
	if aliasField, canonicalField := captureAliasField(document); aliasField != "" {
		return fmt.Errorf("read-only alias %s must be authored as %s", aliasField, canonicalField)
	}
	for _, fieldName := range []string{"id", "title", "status", "created_at", "user_request", "domain", "tdd", "maintenance"} {
		if field, found := document.FieldValue(fieldName); !found || field.ScalarValue == "" || field.ListValues != nil || field.NestedValues != nil {
			return fmt.Errorf("required scalar field %s is missing or malformed", fieldName)
		}
	}
	for _, fieldName := range []string{"prime_files", "required_lessons", "depends_on", "related", "write_set"} {
		if field, found := document.FieldValue(fieldName); found && field.ListValues == nil {
			return fmt.Errorf("field %s must be a list", fieldName)
		}
	}
	if _, found := document.FieldValue("prime_files"); !found {
		return fmt.Errorf("required field prime_files is missing")
	}
	for _, fieldName := range schemanormalization.SchemaFieldNames() {
		field, found := document.FieldValue(fieldName)
		if found && (field.ListValues != nil || field.NestedValues != nil || !schemanormalization.NormalizeField(fieldName, field.ScalarValue).IsCanonicalAuthoringValue) {
			return fmt.Errorf("field %s must use an exact canonical value", fieldName)
		}
	}
	if record.RequestStatus != "pending" && record.RequestStatus != "pending-answers" && record.RequestStatus != "blocked" {
		return fmt.Errorf("new REQ status must be pending, pending-answers, or blocked")
	}
	if !isCanonicalCaptureTimestamp(document, "created_at") {
		return fmt.Errorf("created_at must be a whole-second UTC timestamp")
	}
	for _, fieldName := range []string{"title", "blocked_by", "blocked_check", "assigned_to", "stakeholder"} {
		if _, found := document.FieldValue(fieldName); found && !isSafelyAuthoredUserText(document, fieldName) {
			return fmt.Errorf("%s must use a single-quoted or literal-block scalar", fieldName)
		}
	}
	blockedBy, hasBlockedBy := document.FieldValue("blocked_by")
	_, hasBlockedAt := document.FieldValue("blocked_at")
	_, hasBlockedCheck := document.FieldValue("blocked_check")
	if record.RequestStatus == "blocked" {
		if !hasBlockedBy || blockedBy.ScalarValue == "" || !hasBlockedAt || !isCanonicalCaptureTimestamp(document, "blocked_at") {
			return fmt.Errorf("blocked status requires blocked_by and a whole-second UTC blocked_at")
		}
	} else if hasBlockedBy || hasBlockedAt || hasBlockedCheck {
		return fmt.Errorf("blocked metadata is only valid with blocked status")
	}
	for _, fieldName := range []string{"depends_on", "related"} {
		if field, found := document.FieldValue(fieldName); found {
			for _, requestID := range field.ListValues {
				if !requestPattern.MatchString(requestID) {
					return fmt.Errorf("field %s must contain canonical REQ ids", fieldName)
				}
			}
		}
	}
	if addendumField, found := document.FieldValue("addendum_to"); found && (addendumField.ListValues != nil || !requestPattern.MatchString(addendumField.ScalarValue)) {
		return fmt.Errorf("addendum_to must be one canonical REQ id")
	}
	impactValue := record.ImpactValue
	impactPrefix := "[" + impactValue + "] "
	hasImpactPrefix := regexp.MustCompile(`^\[impact-[^]]+\] `).MatchString(record.RequestTitle)
	if impactValue == "impact-user-visible" {
		if hasImpactPrefix {
			return fmt.Errorf("default impact must not be mirrored in title")
		}
	} else if !strings.HasPrefix(record.RequestTitle, impactPrefix) {
		return fmt.Errorf("non-default impact must be mirrored in title")
	}
	return nil
}

func duplicateCaptureField(document *requestmodel.RequestDocument) string {
	firstDuplicate := ""
	for fieldName, field := range document.TypedRecord().FieldEvidenceByName {
		if field.DuplicateCount > 1 && (firstDuplicate == "" || fieldName < firstDuplicate) {
			firstDuplicate = fieldName
		}
	}
	return firstDuplicate
}

func captureAliasField(document *requestmodel.RequestDocument) (string, string) {
	firstAlias, canonicalKey := "", ""
	for fieldName := range document.TypedRecord().FieldEvidenceByName {
		if canonicalField, isAlias := schemanormalization.CanonicalFieldForAlias(fieldName); isAlias && (firstAlias == "" || fieldName < firstAlias) {
			firstAlias, canonicalKey = fieldName, canonicalField
		}
	}
	return firstAlias, canonicalKey
}

func exactUniqueStrings(actualValues, expectedValues []string) bool {
	if len(actualValues) != len(expectedValues) {
		return false
	}
	expectedSet := make(map[string]bool, len(expectedValues))
	for _, expectedValue := range expectedValues {
		expectedSet[expectedValue] = true
	}
	for _, actualValue := range actualValues {
		if !expectedSet[actualValue] {
			return false
		}
		delete(expectedSet, actualValue)
	}
	return len(expectedSet) == 0
}

func isCanonicalCaptureTimestamp(document *requestmodel.RequestDocument, fieldName string) bool {
	field, found := document.FieldValue(fieldName)
	if !found || field.ListValues != nil || field.NestedValues != nil {
		return false
	}
	parsedTimestamp, parseError := time.Parse("2006-01-02T15:04:05Z", field.ScalarValue)
	return parseError == nil && parsedTimestamp.Format("2006-01-02T15:04:05Z") == field.ScalarValue
}

func isSafelyAuthoredUserText(document *requestmodel.RequestDocument, fieldName string) bool {
	field, found := document.FieldValue(fieldName)
	if !found || field.ListValues != nil || field.NestedValues != nil {
		return false
	}
	rawValue := strings.TrimSpace(field.RawValue)
	if rawValue == "|" || rawValue == "|-" || rawValue == "|+" {
		return true
	}
	if len(rawValue) < 2 || rawValue[0] != '\'' || rawValue[len(rawValue)-1] != '\'' {
		return false
	}
	for byteIndex := 1; byteIndex < len(rawValue)-1; byteIndex++ {
		if rawValue[byteIndex] != '\'' {
			continue
		}
		if byteIndex+1 >= len(rawValue)-1 || rawValue[byteIndex+1] != '\'' {
			return false
		}
		byteIndex++
	}
	return true
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
