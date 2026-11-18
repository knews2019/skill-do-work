package toolboxcommands

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

const (
	reportGenerationUsage = "--prepare-generation <manifest.json> [<parent/index.html>] | --generation-digest <report/index.html> | --record-generation <report/index.html> <receipt.json> | --record-factual-review <report/index.html> <review.json>"
	reportGenerationFile  = "gallery-generation.json"
	reportFactReviewFile  = "gallery-factual-review.json"
	reportGenerationModel = "gpt-6.2-astra"
)

type reportEvidenceManifest struct {
	SourceID string   `json:"source_id"`
	Files    []string `json:"files"`
}

type reportEvidenceFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type reportEvidencePin struct {
	SourceID string               `json:"source_id"`
	Files    []reportEvidenceFile `json:"files"`
}

type reportGenerationPlan struct {
	Evidence     reportEvidencePin    `json:"evidence"`
	ParentReport string               `json:"parent_report"`
	ParentSHA256 string               `json:"parent_sha256"`
	Marks        []reportGalleryEntry `json:"marks"`
}

type reportGenerationExecution struct {
	ProviderInstanceID string `json:"provider_instance_id"`
	ExecutedModel      string `json:"executed_model"`
	Status             string `json:"status"`
	TaskID             string `json:"task_id"`
	ThreadID           string `json:"thread_id"`
	RunID              string `json:"run_id"`
}

func (execution reportGenerationExecution) completed() bool {
	return execution.Status == "completed" && execution.ExecutedModel == reportGenerationModel &&
		reportGenerationText(execution.ProviderInstanceID) && reportGenerationText(execution.TaskID) &&
		reportGenerationText(execution.ThreadID) && reportGenerationText(execution.RunID)
}

type reportGenerationChange struct {
	Before      string `json:"before"`
	After       string `json:"after"`
	Description string `json:"description"`
}

type reportGenerationReceipt struct {
	Plan      reportGenerationPlan      `json:"plan"`
	Execution reportGenerationExecution `json:"execution"`
	Changes   []reportGenerationChange  `json:"changes"`
}

type reportGenerationRecord struct {
	Version        int    `json:"version"`
	ReportPath     string `json:"report_path"`
	DesignIdentity string `json:"design_identity"`
	reportGenerationReceipt
}

type reportFactFinding struct {
	Anchor string `json:"anchor"`
	Reason string `json:"reason"`
}

type reportFactReview struct {
	BundleSHA256  string              `json:"bundle_sha256"`
	RenderVerdict string              `json:"render_verdict"`
	Verdict       string              `json:"verdict"`
	Findings      []reportFactFinding `json:"findings"`
}

type reportFactRecord struct {
	Version int `json:"version"`
	reportFactReview
}

func handleReportGalleryGeneration(repositoryRoot string, args []string) resultmodel.CommandResult {
	option := args[0]
	valid := len(args) == 3
	if option == "--prepare-generation" {
		valid = len(args) == 2 || (len(args) == 3 && args[2] != "")
	} else if option == "--generation-digest" {
		valid = len(args) == 2
	}
	if !valid {
		return usageResult(CommandAIReportGallery, "Usage: ai-report-gallery "+reportGenerationUsage)
	}
	repository, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return reportGalleryFailure(".", err)
	}
	defer repository.Close()

	target := args[1]
	var output string
	var changes []resultmodel.RecordedChange
	switch option {
	case "--prepare-generation":
		var parent string
		if len(args) == 3 {
			parent = args[2]
		}
		output, err = prepareReportGeneration(repository, target, parent)
	case "--generation-digest":
		var record reportGenerationRecord
		if record, err = readReportGenerationRecord(repository, target); err == nil {
			if err = verifyReportEvidence(repository, record.Plan.Evidence); err == nil {
				output, err = reportGenerationDigest(repository, target)
			}
		}
	case "--record-generation":
		changes, err = saveReportGeneration(repository, target, args[2])
	case "--record-factual-review":
		changes, err = saveReportFactReview(repository, target, args[2])
	}
	if err != nil {
		return reportGalleryFailure(target, err)
	}
	if len(changes) != 0 {
		output = changes[0].Path
	}
	return exactOutputResult(output+"\n", changes)
}

func prepareReportGeneration(repository *os.Root, manifestPath, parent string) (string, error) {
	var manifest reportEvidenceManifest
	if err := readReportGenerationInput(repository, manifestPath, &manifest); err != nil {
		return "", err
	}
	evidence, err := pinReportEvidence(repository, manifest.SourceID, manifest.Files)
	if err != nil {
		return "", err
	}
	plan := reportGenerationPlan{Evidence: evidence, ParentReport: parent, Marks: []reportGalleryEntry{}}
	if parent != "" {
		digest, err := eligibleReportParent(repository, parent, evidence)
		if err != nil {
			return "", err
		}
		plan.ParentSHA256 = digest
		page, err := readReportGenerationFile(repository, parent)
		if err != nil {
			return "", err
		}
		feedback, _, err := readReportGalleryFeedback(repository.Name(), repository)
		if err != nil {
			return "", err
		}
		figures := reportGenerationFigures(page)
		for _, mark := range feedback.Entries {
			if mark.ReportPath == parent && (mark.Kind == "report" || figures[mark.FigureID]) {
				plan.Marks = append(plan.Marks, mark)
			}
		}
		if len(plan.Marks) == 0 {
			return "", fmt.Errorf("excluded %s: no current feedback on the report or its figures", parent)
		}
	}
	content, err := json.Marshal(plan)
	return string(content), err
}

func eligibleReportParent(repository *os.Root, parent string, evidence reportEvidencePin) (string, error) {
	digest, err := reviewedReportParent(repository, parent, evidence)
	if err != nil {
		return "", fmt.Errorf("excluded %s: %w", parent, err)
	}
	return digest, nil
}

func reviewedReportParent(repository *os.Root, parent string, evidence reportEvidencePin) (string, error) {
	record, err := readReportGenerationRecord(repository, parent)
	if err != nil {
		return "", err
	}
	if !sameReportEvidence(record.Plan.Evidence, evidence) {
		return "", errors.New("evidence pin differs from the parent generation")
	}
	var review reportFactRecord
	_, err = readReportSidecar(repository, parent, reportFactReviewFile, &review)
	if errors.Is(err, fs.ErrNotExist) {
		return "", errors.New("no factual review")
	}
	if err != nil {
		return "", err
	}
	page, err := readReportGenerationFile(repository, parent)
	if err != nil {
		return "", err
	}
	digest, err := reportGenerationDigest(repository, parent)
	if err != nil {
		return "", err
	}
	if review.BundleSHA256 != digest {
		return "", errors.New("stale factual review: report bundle changed")
	}
	if review.Version != 1 {
		return "", errors.New("unsupported factual review version")
	}
	if err := validateReportFactReview(review.reportFactReview, record.Plan.ParentReport, page); err != nil {
		return "", err
	}
	if review.Verdict != "pass" {
		reasons := make([]string, 0, len(review.Findings))
		for _, finding := range review.Findings {
			reasons = append(reasons, finding.Anchor+": "+finding.Reason)
		}
		return "", fmt.Errorf("factual review %s: %s", review.Verdict, strings.Join(reasons, "; "))
	}
	return digest, nil
}

func saveReportGeneration(repository *os.Root, report string, receiptPath string) ([]resultmodel.RecordedChange, error) {
	page, err := readReportGenerationPage(repository, report)
	if err != nil {
		return nil, err
	}
	var receipt reportGenerationReceipt
	if err := readReportGenerationInput(repository, receiptPath, &receipt); err != nil {
		return nil, err
	}
	if !receipt.Execution.completed() {
		return nil, errors.New("a completed T3 execution of " + reportGenerationModel + " with provider instance, task, thread and run IDs is required")
	}
	source := reportGalleryGeneratedSource{bundle: reportBundle{Path: path.Dir(report), Entry: report}}
	design, err := source.readDesign(repository)
	if err != nil {
		return nil, err
	}
	if design.recordPath == report {
		return nil, errors.New("the generated report requires gallery-design.md")
	}
	parentPage, err := validateReportGenerationPlan(repository, report, receipt.Plan)
	if err != nil {
		return nil, err
	}
	if err := validateReportGenerationChanges(report, page, receipt.Plan.ParentReport, parentPage, receipt.Changes); err != nil {
		return nil, err
	}
	record := reportGenerationRecord{Version: 1, ReportPath: report, DesignIdentity: source.qualifiedIdentity(), reportGenerationReceipt: receipt}
	return publishReportSidecar(repository, report, reportGenerationFile, record, nil)
}

func validateReportGenerationPlan(repository *os.Root, report string, plan reportGenerationPlan) ([]byte, error) {
	if err := verifyReportEvidence(repository, plan.Evidence); err != nil {
		return nil, err
	}
	if plan.Marks == nil {
		return nil, errors.New("marks must be an array")
	}
	if plan.ParentReport == "" {
		if len(plan.Marks) != 0 || plan.ParentSHA256 != "" {
			return nil, errors.New("an initial report must have an empty parent digest and no marks")
		}
		return nil, nil
	}
	if plan.ParentReport == report || len(plan.Marks) == 0 {
		return nil, errors.New("a child needs a different parent and at least one saved mark")
	}
	digest, err := eligibleReportParent(repository, plan.ParentReport, plan.Evidence)
	if err != nil {
		return nil, err
	}
	if plan.ParentSHA256 != digest {
		return nil, errors.New("parent bundle changed after generation was prepared")
	}
	page, err := readReportGenerationFile(repository, plan.ParentReport)
	if err != nil {
		return nil, err
	}
	figures := reportGenerationFigures(page)
	seen := make(map[reportGalleryTarget]bool, len(plan.Marks))
	for _, mark := range plan.Marks {
		if !validReportGalleryEntry(mark) || mark.ReportPath != plan.ParentReport ||
			(mark.Kind == "figure" && !figures[mark.FigureID]) || seen[mark.reportGalleryTarget] {
			return nil, errors.New("a snapshotted mark has an invalid, duplicate or unrelated target")
		}
		seen[mark.reportGalleryTarget] = true
	}
	return page, nil
}

func reportGenerationFigures(page []byte) map[string]bool {
	figures := make(map[string]bool)
	for _, figure := range reportGalleryFigures(page, "") {
		figures[figure.ID] = true
	}
	return figures
}

func validateReportGenerationChanges(report string, page []byte, parent string, parentPage []byte, changes []reportGenerationChange) error {
	if len(changes) == 0 {
		return errors.New("at least one anchored change is required")
	}
	after := reportGenerationAnchors(page)
	before := reportGenerationAnchors(parentPage)
	for _, change := range changes {
		if !reportGenerationText(change.Description) {
			return errors.New("every change needs a description")
		}
		if !reportGenerationLink(report, change.After, after) {
			return fmt.Errorf("change has a missing or ambiguous after anchor: %q", change.After)
		}
		if (parent == "" && change.Before != "") || (parent != "" && !reportGenerationLink(parent, change.Before, before)) {
			return fmt.Errorf("change has a missing or invalid before anchor for parent %q: %q", parent, change.Before)
		}
	}
	return nil
}

var reportGenerationTag = regexp.MustCompile(`(?is)<[a-z][a-z0-9:-]*\b(` + reportGalleryTagAttrs + `)>`)

func reportGenerationAnchors(page []byte) map[string]int {
	anchors := make(map[string]int)
	for _, tag := range reportGenerationTag.FindAllStringSubmatch(reportGalleryHidden.ReplaceAllString(string(page), " "), -1) {
		if id := galleryVisualAttrs(tag[1])["id"]; id != "" {
			anchors[id]++
		}
	}
	return anchors
}

func reportGenerationLink(page, link string, anchors map[string]int) bool {
	file, anchor, ok := strings.Cut(link, "#")
	return ok && file == page && anchors[anchor] == 1
}

func readReportGenerationRecord(repository *os.Root, report string) (reportGenerationRecord, error) {
	var record reportGenerationRecord
	if _, err := readReportGenerationPage(repository, report); err != nil {
		return record, err
	}
	_, err := readReportSidecar(repository, report, reportGenerationFile, &record)
	if errors.Is(err, fs.ErrNotExist) {
		return record, errors.New("no generation receipt")
	}
	if err != nil {
		return record, err
	}
	if record.Version != 1 || record.ReportPath != report || !record.Execution.completed() {
		return record, errors.New("invalid or relocated generation receipt")
	}
	return record, nil
}

func saveReportFactReview(repository *os.Root, report, reviewPath string) ([]resultmodel.RecordedChange, error) {
	record, err := readReportGenerationRecord(repository, report)
	if err != nil {
		return nil, err
	}
	if err := verifyReportEvidence(repository, record.Plan.Evidence); err != nil {
		return nil, err
	}
	var review reportFactReview
	if err := readReportGenerationInput(repository, reviewPath, &review); err != nil {
		return nil, err
	}
	page, err := readReportGenerationFile(repository, report)
	if err != nil {
		return nil, err
	}
	digest, err := reportGenerationDigest(repository, report)
	if err != nil {
		return nil, err
	}
	if digest != review.BundleSHA256 {
		return nil, errors.New("the report bundle changed since the review snapshot")
	}
	if err := validateReportFactReview(review, record.Plan.ParentReport, page); err != nil {
		return nil, err
	}
	var previous reportFactRecord
	original, err := readReportSidecar(repository, report, reportFactReviewFile, &previous)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	// An incomplete render may be retried. A completed factual verdict on the
	// same bytes is immutable, including a failure on a liked report.
	if previous.BundleSHA256 == digest && previous.Verdict != "unverified" {
		return nil, errors.New("a factual verdict already exists for these bytes")
	}
	return publishReportSidecar(repository, report, reportFactReviewFile, reportFactRecord{Version: 1, reportFactReview: review}, original)
}

func validateReportFactReview(review reportFactReview, parent string, page []byte) error {
	switch review.RenderVerdict {
	case "pass", "fail", "skipped", "error":
	case "parent-links-only":
		if parent == "" {
			return errors.New("parent-links-only requires a parent")
		}
	default:
		return errors.New("unknown render verdict")
	}
	switch review.Verdict {
	case "pass":
		if review.RenderVerdict != "pass" && review.RenderVerdict != "parent-links-only" {
			return errors.New("a passing factual review requires a verified render")
		}
	case "fail", "unverified":
	default:
		return errors.New("factual verdict must be pass, fail or unverified")
	}
	if len(review.Findings) == 0 {
		return errors.New("factual review requires concrete findings")
	}
	anchors := reportGenerationAnchors(page)
	for _, finding := range review.Findings {
		if anchors[finding.Anchor] != 1 || strings.TrimSpace(finding.Reason) == "" {
			return fmt.Errorf("factual finding needs a unique page anchor and a reason: %q", finding.Anchor)
		}
	}
	return nil
}

func pinReportEvidence(repository *os.Root, sourceID string, files []string) (reportEvidencePin, error) {
	pin := reportEvidencePin{SourceID: sourceID}
	if !reportGenerationText(sourceID) || len(files) == 0 {
		return pin, errors.New("the evidence manifest requires a source_id and at least one file")
	}
	files = slices.Clone(files)
	slices.Sort(files)
	for index, file := range files {
		if index > 0 && file == files[index-1] {
			return pin, fmt.Errorf("duplicate evidence file %q", file)
		}
		if err := validateReportGenerationFile(repository, file); err != nil {
			return pin, err
		}
		digest, err := hashReportGenerationFile(repository, file)
		if err != nil {
			return pin, err
		}
		pin.Files = append(pin.Files, reportEvidenceFile{Path: file, SHA256: digest})
	}
	return pin, nil
}

func sameReportEvidence(first, second reportEvidencePin) bool {
	return first.SourceID == second.SourceID && slices.Equal(first.Files, second.Files)
}

func verifyReportEvidence(repository *os.Root, pin reportEvidencePin) error {
	names := make([]string, 0, len(pin.Files))
	for _, file := range pin.Files {
		names = append(names, file.Path)
	}
	current, err := pinReportEvidence(repository, pin.SourceID, names)
	if err != nil {
		return fmt.Errorf("evidence pin: %w", err)
	}
	if !sameReportEvidence(pin, current) {
		return errors.New("evidence pin changed after generation was prepared")
	}
	return nil
}

func readReportGenerationPage(repository *os.Root, report string) ([]byte, error) {
	parts := strings.Split(report, "/")
	var direct, monthly bool
	if len(parts) >= 3 {
		direct = len(parts) == 3 && !reportMonthFolder(parts[1])
		monthly = len(parts) == 4 && reportMonthFolder(parts[1])
	}
	if (!direct && !monthly) || parts[0] != reportsFolder || path.Base(report) != "index.html" ||
		strings.HasPrefix(path.Base(path.Dir(report)), ".") || strings.ContainsAny(report, ":#?") {
		return nil, fmt.Errorf("%q must be the exact project-relative ai-reports bundle index.html path", report)
	}
	return readReportGenerationFile(repository, report)
}

// Hash every regular bundle file, including local assets and the generation
// receipt. The factual review is the only excluded file because it stores this hash.
func reportGenerationDigest(repository *os.Root, report string) (string, error) {
	if _, err := readReportGenerationPage(repository, report); err != nil {
		return "", err
	}
	bundle, err := rootedfs.OpenRoot(repository, filepath.FromSlash(path.Dir(report)))
	if err != nil {
		return "", err
	}
	defer bundle.Close()
	digest := sha256.New()
	err = fs.WalkDir(bundle.FS(), ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == reportFactReviewFile {
			if !entry.Type().IsRegular() {
				return errors.New("factual review sidecar must be a regular file")
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("bundle file %q must be regular", name)
		}
		fileDigest, err := hashReportGenerationFile(bundle, name)
		if err != nil {
			return err
		}
		fmt.Fprintf(digest, "%s\x00%s\n", name, fileDigest)
		return nil
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", digest.Sum(nil)), nil
}

func hashReportGenerationFile(repository *os.Root, name string) (string, error) {
	handle, err := repository.Open(filepath.FromSlash(name))
	if err != nil {
		return "", err
	}
	defer handle.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, handle); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", digest.Sum(nil)), nil
}

func reportGenerationText(value string) bool {
	return strings.TrimSpace(value) != "" && strings.IndexFunc(value, unicode.IsControl) == -1
}

func validateReportGenerationFile(repository *os.Root, name string) error {
	if !fs.ValidPath(name) || name == "." || strings.Contains(name, `\`) || !reportGenerationText(name) {
		return fmt.Errorf("%q must be a project-relative file", name)
	}
	if err := validateNoLinkedAncestors(repository.Name(), name, true); err != nil {
		return err
	}
	info, err := repository.Lstat(filepath.FromSlash(name))
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%q must be a regular file", name)
	}
	return nil
}

func readReportGenerationFile(repository *os.Root, name string) ([]byte, error) {
	if err := validateReportGenerationFile(repository, name); err != nil {
		return nil, err
	}
	return rootedfs.ReadFile(repository, filepath.FromSlash(name))
}

func readReportGenerationInput(repository *os.Root, name string, target any) error {
	var content []byte
	var err error
	if filepath.IsAbs(name) {
		// Explicit caller input can live in scratch outside the project.
		content, err = os.ReadFile(name)
	} else {
		content, err = readReportGenerationFile(repository, name)
	}
	if err != nil {
		return err
	}
	return decodeReportGenerationJSON(name, content, target)
}

func readReportSidecar(repository *os.Root, report, leaf string, target any) ([]byte, error) {
	name := path.Join(path.Dir(report), leaf)
	content, err := readReportGenerationFile(repository, name)
	if err != nil {
		return nil, err
	}
	return content, decodeReportGenerationJSON(name, content, target)
}

func decodeReportGenerationJSON(name string, content []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("%s must contain one JSON document", name)
	}
	return nil
}

func publishReportSidecar(repository *os.Root, report, leaf string, record any, original []byte) ([]resultmodel.RecordedChange, error) {
	name := path.Join(path.Dir(report), leaf)
	content, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, err
	}
	err = rootedPublishFileChecked(repository.Name(), name, append(content, '\n'), 0o644, original != nil, func(current []byte) error {
		if !bytes.Equal(current, original) {
			return fmt.Errorf("%s changed during the save", name)
		}
		return nil
	})
	if errors.Is(err, fs.ErrExist) {
		return nil, fmt.Errorf("%s already exists", name)
	}
	if err != nil {
		return nil, err
	}
	kind := "created"
	if original != nil {
		kind = "modified"
	}
	return []resultmodel.RecordedChange{{Path: name, Kind: kind, Detail: "saved report gallery record"}}, nil
}
