package toolboxcommands

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

const (
	reportGalleryFolder   = "docs/report-gallery"
	ReportGalleryPageFile = reportGalleryFolder + "/index.html"
	reportGalleryRecord   = "sample.md"
	reportGalleryExample  = "example.html"
)

var (
	reportGalleryNamePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)+$`)
	errHandMadeGallery       = errors.New("missing the ai-report-gallery generator header")
)

type reportGallerySample struct {
	Name       string
	RecordPath string
	Title      string
	UseWhen    string
	Source     string
	Image      string
	Borrow     []string
}

func handleAIReportGallery(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	if len(args) > 0 {
		switch args[0] {
		case "--prepare-generation", "--generation-digest", "--record-generation", "--record-factual-review":
			return handleReportGalleryGeneration(ctx.RepositoryRoot, args)
		}
	}
	listOnly := len(args) == 1 && args[0] == "--list"
	resolveOnly := len(args) == 2 && args[0] == "--resolve"
	serve := len(args) == 1 && args[0] == "--serve"
	serveWithoutBrowser := len(args) == 2 && args[0] == "--serve" && args[1] == "--no-open"
	capturePreviews := len(args) >= 1 && len(args) <= 2 && args[0] == "--capture-previews"
	if len(args) != 0 && !listOnly && !resolveOnly && !serve && !serveWithoutBrowser && !capturePreviews {
		return usageResult(CommandAIReportGallery, "Usage: ai-report-gallery [--list | --resolve <name> | --serve [--no-open] | --capture-previews [<report entry>] | "+reportGenerationUsage+"]")
	}
	repository, err := os.OpenRoot(ctx.RepositoryRoot)
	if err != nil {
		return reportGalleryFailure(".", err)
	}
	defer repository.Close()

	if capturePreviews {
		return captureReportGalleryPreviews(ctx.RepositoryRoot, repository, args[1:])
	}
	if resolveOnly && strings.HasPrefix(strings.TrimSpace(args[1]), reportGalleryGeneratedPrefix) {
		return resolveReportGalleryGenerated(ctx.RepositoryRoot, repository, strings.TrimSpace(args[1]))
	}

	// Reject linked output parents before discovery. The publisher also confines the final write.
	if len(args) == 0 {
		if err := validateNoLinkedAncestors(ctx.RepositoryRoot, ReportGalleryPageFile, false); err != nil {
			return reportGalleryRefusal(err)
		}
	}
	toolboxDocs, err := reportGalleryToolboxDocs(repository)
	if err != nil {
		return reportGalleryFailure("do-work-toolbox/docs", err)
	}
	var shippedFolder string
	if toolboxDocs != "" {
		shippedFolder = path.Join(toolboxDocs, "report-gallery")
	}
	samples, findings, err := discoverReportGallerySamples(repository, shippedFolder)
	if err != nil {
		return reportGalleryFailure(reportGalleryFolder, err)
	}

	var result resultmodel.CommandResult
	switch {
	case listOnly:
		var output strings.Builder
		for _, sample := range samples {
			fmt.Fprintf(&output, "%s\t%s\t%s\n", sample.Name, sample.UseWhen, sample.RecordPath)
		}
		result = exactOutputResult(output.String(), nil)
	case resolveOnly:
		result = resolveReportGallerySample(samples, args[1])
	case serve || serveWithoutBrowser:
		result = serveReportGallery(ctx.RepositoryRoot, repository, toolboxDocs, samples, serve)
	default:
		result = writeReportGallery(ctx.RepositoryRoot, repository, toolboxDocs, samples)
	}
	result.Findings = append(result.Findings, findings...)
	return result
}

// matchReportGalleryName gives canonical folder slugs precedence over titles.
func matchReportGalleryName(samples []reportGallerySample, name string) []reportGallerySample {
	name = strings.TrimSpace(name)
	var matches []reportGallerySample
	for _, sample := range samples {
		if sample.Name == name {
			return []reportGallerySample{sample}
		}
		if strings.EqualFold(sample.Title, name) {
			matches = append(matches, sample)
		}
	}
	return matches
}

// reportGalleryCopyNames chooses the visible title whenever it identifies its own sample.
func reportGalleryCopyNames(samples []reportGallerySample) map[string]string {
	names := make(map[string]string, len(samples))
	for _, sample := range samples {
		matches := matchReportGalleryName(samples, sample.Title)
		if len(matches) == 1 && matches[0].Name == sample.Name && !strings.HasPrefix(sample.Title, reportGalleryGeneratedPrefix) {
			names[sample.Name] = sample.Title
		} else {
			names[sample.Name] = sample.Name
		}
	}
	return names
}

func resolveReportGallerySample(samples []reportGallerySample, name string) resultmodel.CommandResult {
	matches := matchReportGalleryName(samples, name)
	switch len(matches) {
	case 0:
		return reportGalleryFailure(reportGalleryFolder, fmt.Errorf(
			"no gallery design matches %q. Browse the report gallery for a title or run ai-report-gallery --list for folder slugs", name))
	case 1:
		sample := matches[0]
		return exactOutputResult(fmt.Sprintf("%s\t%s\t%s\n", sample.Name, sample.UseWhen, sample.RecordPath), nil)
	default:
		slugs := make([]string, 0, len(matches))
		for _, sample := range matches {
			slugs = append(slugs, sample.Name)
		}
		return reportGalleryFailure(reportGalleryFolder, fmt.Errorf(
			"gallery design %q is ambiguous. Use one of these folder slugs: %s", name, strings.Join(slugs, ", ")))
	}
}

// Installed files take precedence when a project also contains a suite checkout.
func reportGalleryToolboxDocs(repository *os.Root) (string, error) {
	for _, candidate := range []string{".claude/skills/do-work-toolbox/docs", "skills/do-work-toolbox/docs"} {
		info, err := repository.Lstat(filepath.FromSlash(candidate))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if !info.IsDir() {
			return "", fmt.Errorf("%s must be a directory", candidate)
		}
		return candidate, nil
	}
	return "", nil
}

func discoverReportGallerySamples(repository *os.Root, shippedFolder string) ([]reportGallerySample, []resultmodel.CommandFinding, error) {
	byName := map[string]reportGallerySample{}
	var findings []resultmodel.CommandFinding
	for _, folder := range []string{shippedFolder, reportGalleryFolder} {
		if folder == "" {
			continue
		}
		entries, err := fs.ReadDir(repository.FS(), folder)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, findings, err
		}
		for _, entry := range entries {
			name := entry.Name()
			if !entry.IsDir() || !reportGalleryNamePattern.MatchString(name) {
				continue
			}
			record := path.Join(folder, name, reportGalleryRecord)
			sample, err := readReportGallerySample(repository, name, record)
			if err != nil {
				findings = append(findings, toolboxFinding(CommandAIReportGallery, "AI-REPORT-GALLERY-SAMPLE-SKIPPED", resultmodel.SeverityWarning,
					[]string{record}, fmt.Sprintf("skipped sample %s: %v", name, err), resultmodel.FixabilityManual, ""))
				continue
			}
			byName[name] = sample
		}
	}
	samples := make([]reportGallerySample, 0, len(byName))
	for _, sample := range byName {
		samples = append(samples, sample)
	}
	sort.Slice(samples, func(left, right int) bool { return samples[left].Name < samples[right].Name })
	return samples, findings, nil
}

func readReportGallerySample(repository *os.Root, name, record string) (reportGallerySample, error) {
	sample := reportGallerySample{Name: name, RecordPath: record}
	filename := filepath.FromSlash(record)
	document, err := readReportGalleryRecord(repository, filename)
	if err != nil {
		return sample, err
	}
	var added string
	for _, field := range []struct {
		name  string
		value *string
	}{
		{"title", &sample.Title}, {"use_when", &sample.UseWhen}, {"source", &sample.Source}, {"added", &added}, {"image", &sample.Image},
	} {
		if _, present := document.FieldValue(field.name); !present && field.name == "image" {
			continue
		}
		*field.value, err = reportGalleryScalar(document, field.name)
		if err != nil {
			return sample, err
		}
	}
	if _, err := time.Parse(time.DateOnly, added); err != nil {
		return sample, errors.New("added must be a valid YYYY-MM-DD date")
	}
	if sample.Image != "" && (!fs.ValidPath(sample.Image) || sample.Image == "." || strings.ContainsAny(sample.Image, `\:`)) {
		return sample, errors.New("image must be a path relative to the sample folder")
	}
	sample.Borrow, err = readReportGalleryBorrow(document.BodyBytes())
	return sample, err
}

func readReportGalleryRecord(repository *os.Root, filename string) (*requestmodel.RequestDocument, error) {
	info, err := repository.Lstat(filename)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s must be a regular file", filepath.Base(filename))
	}
	content, err := rootedfs.ReadFile(repository, filename)
	if err != nil {
		return nil, err
	}
	document, err := requestmodel.ParseDocument(content)
	if err != nil {
		return nil, err
	}
	if warnings := document.ParseWarnings(); len(warnings) != 0 {
		return nil, errors.New(warnings[0])
	}
	return document, nil
}

func reportGalleryScalar(document *requestmodel.RequestDocument, name string) (string, error) {
	evidence, _ := document.FieldValue(name)
	value := strings.TrimSpace(evidence.ScalarValue)
	raw := evidence.RawValue
	quoted := strings.HasPrefix(raw, `"`) || strings.HasPrefix(raw, "'")
	if value == "" || (quoted && evidence.ScalarValue == raw) || evidence.ListValues != nil || strings.HasPrefix(raw, "{") ||
		strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return "", fmt.Errorf("%s must be a nonempty one-line scalar", name)
	}
	return value, nil
}

func readReportGalleryBorrow(body []byte) ([]string, error) {
	var bullets []string
	for _, section := range requestmodel.VisibleSections(body) {
		if section.Name != "Borrow" || section.HeadingIndent != 0 {
			continue
		}
		if bullets != nil {
			return nil, errors.New("record has more than one ## Borrow section")
		}
		bullets = []string{}
		_, remainder, _ := strings.Cut(string(body[section.Start:section.End]), "\n")
		for _, line := range strings.Split(remainder, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if len(line) < 3 || !strings.ContainsRune("-*+", rune(line[0])) || line[1] != ' ' {
				return nil, errors.New("## Borrow must contain only Markdown bullets")
			}
			bullets = append(bullets, strings.TrimSpace(line[2:]))
		}
	}
	if len(bullets) < 2 || len(bullets) > 3 {
		return nil, errors.New("## Borrow must contain two or three Markdown bullets")
	}
	return bullets, nil
}

func writeReportGallery(repositoryRoot string, repository *os.Root, toolboxDocs string, samples []reportGallerySample) resultmodel.CommandResult {
	info, err := repository.Lstat(filepath.FromSlash(ReportGalleryPageFile))
	exists := err == nil
	if exists && !info.Mode().IsRegular() {
		return reportGalleryRefusal(errors.New("the output is not a regular file"))
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return reportGalleryFailure(ReportGalleryPageFile, err)
	}

	bundles, _, err := reportIndexWalk(repositoryRoot)
	if err != nil {
		return reportGalleryFailure(reportsFolder, err)
	}
	page, err := renderReportGalleryPage(repositoryRoot, repository, toolboxDocs, samples, bundles, nil)
	if err != nil {
		return reportGalleryFailure(ReportGalleryPageFile, err)
	}
	if !exists {
		if err := rootedMkdirAll(repositoryRoot, reportGalleryFolder, 0o755); err != nil {
			return reportGalleryFailure(ReportGalleryPageFile, err)
		}
	}
	checkGenerator := func(original []byte) error {
		if !bytes.HasPrefix(original, []byte(reportGalleryPageHeader)) {
			return errHandMadeGallery
		}
		return nil
	}
	if err := rootedPublishFileChecked(repositoryRoot, ReportGalleryPageFile, page, 0o644, exists, checkGenerator); err != nil {
		if errors.Is(err, errHandMadeGallery) {
			return reportGalleryRefusal(err)
		}
		return reportGalleryFailure(ReportGalleryPageFile, err)
	}
	kind := "created"
	if exists {
		kind = "modified"
	}
	return exactOutputResult(ReportGalleryPageFile+"\n", []resultmodel.RecordedChange{
		{Path: ReportGalleryPageFile, Kind: kind, Detail: "regenerated derived report gallery"},
	})
}

func reportGalleryRefusal(err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandAIReportGallery, "AI-REPORT-GALLERY-OUTPUT-REFUSED", resultmodel.SeverityError,
			[]string{ReportGalleryPageFile}, fmt.Sprintf("refusing to replace %s: %v", ReportGalleryPageFile, err),
			resultmodel.FixabilityManual, "the gallery can only replace its own generated page at this path"),
	}}
}

func reportGalleryFailure(relative string, err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandAIReportGallery, "AI-REPORT-GALLERY-FAILED", resultmodel.SeverityError, []string{relative},
			err.Error(), resultmodel.FixabilityManual, "the report gallery could not be read or written"),
	}}
}
