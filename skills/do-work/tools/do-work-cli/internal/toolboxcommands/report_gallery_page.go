package toolboxcommands

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

type reportGalleryLink struct {
	Label string
	Href  string
}

type reportGalleryCard struct {
	AnchorID string
	NameID   string
	Title    string
	CopyName string
	UseWhen  string
	Starter  reportGalleryStarter
	Preview  string
	Open     reportGalleryLink
	More     []reportGalleryLink
}

type reportGalleryStarter struct {
	ID, Prompt string
	Borrow     []string
}

type reportGalleryRepositoryCard struct {
	Title, UseWhen, Report, Preview string
	Starter                         reportGalleryStarter
	Review                          *reportGalleryReview
	Worked                          *workedReportCard
}

type reportGalleryKind struct {
	Label   string
	Command string
	Latest  reportGalleryLink
}

type reportGalleryPageData struct {
	Stylesheet       string
	Samples          []reportGalleryCard
	Reports          []reportGalleryRepositoryCard
	HasWorkedReports bool
	Kinds            []reportGalleryKind
	Presets          []reportGalleryLink
	Live             *reportGalleryLivePage
	workedNames      map[string]bool
}

func (page reportGalleryPageData) UnclaimedID(id string) string {
	if page.workedNames[id] {
		return "gallery_" + id // The underscore keeps this ID outside the worked name format.
	}
	return id
}

func renderReportGalleryPage(layout *template.Template, repositoryRoot string, repository *os.Root, toolboxDocs string, samples []reportGallerySample, bundles []reportBundle, live *reportGalleryLivePage) ([]byte, error) {
	page := reportGalleryPageData{Kinds: reportGalleryKinds(bundles), Live: live}
	if toolboxDocs != "" {
		stylesheet := path.Join(toolboxDocs, "ai-report-kit", "presentation.css")
		page.Stylesheet = reportGalleryFileLink(repository, stylesheet)
		var err error
		page.Presets, err = reportGalleryPresets(repository, path.Join(toolboxDocs, secondOpinionPresetsFolder))
		if err != nil {
			return nil, err
		}
	}
	worked, err := readWorkedReports(repository)
	if err != nil {
		return nil, err
	}
	page.HasWorkedReports = len(worked) != 0
	if live != nil {
		if err := live.discover(repositoryRoot, worked, bundles); err != nil {
			return nil, err
		}
	}
	workedByReport := make(map[string]*workedReportCard, len(worked))
	page.workedNames = make(map[string]bool, len(worked))
	for index := range worked {
		card := &worked[index]
		workedByReport[card.ReportPath] = card
		page.workedNames[card.Name] = true
	}
	generatedByReport := make(map[string]reportGalleryGeneratedSource, len(bundles))
	for _, source := range reportGalleryGeneratedSources(repository, bundles) {
		generatedByReport[source.bundle.Entry] = source
	}
	seenReports := make(map[string]bool, len(bundles))
	for index, bundle := range reportGalleryBundles(bundles, worked) {
		if bundle.Entry == "" || seenReports[bundle.Entry] {
			continue
		}
		seenReports[bundle.Entry] = true
		card := reportGalleryRepositoryCard{Worked: workedByReport[bundle.Entry]}
		if live != nil {
			card.Review = live.Reviews[bundle.Entry]
		}
		source, found := generatedByReport[bundle.Entry]
		if card.Worked == nil && !found && card.Review == nil {
			continue
		}
		if card.Worked == nil {
			card.Title = bundle.Title
			card.Report = reportGalleryHref(bundle.Entry)
			card.Preview = reportGalleryFileLink(repository, reportGalleryPreviewPath(bundle.Entry))
		}
		if found {
			if design, err := source.readDesign(repository); err == nil {
				starterID := "generated-prompt-" + strconv.Itoa(index)
				if card.Worked != nil {
					starterID = "worked-prompt-" + card.Worked.Name
				} else if design.recordPath != bundle.Entry {
					card.Title, card.UseWhen = design.title, design.useWhen
				}
				card.Starter = newReportGalleryStarter(starterID, source.identity, design.borrow)
			}
		}
		page.Reports = append(page.Reports, card)
	}
	copyNames := reportGalleryCopyNames(samples)
	for _, sample := range samples {
		card := reportGalleryCard{
			AnchorID: page.UnclaimedID("sample-" + sample.Name),
			NameID:   page.UnclaimedID("name-" + sample.Name),
			Title:    sample.Title,
			CopyName: copyNames[sample.Name],
			UseWhen:  sample.UseWhen,
		}
		card.Starter = newReportGalleryStarter(page.UnclaimedID("prompt-"+sample.Name), card.CopyName, sample.Borrow)
		folder := path.Dir(sample.RecordPath)
		var links []reportGalleryLink
		if href := reportGalleryFileLink(repository, path.Join(folder, reportGalleryExample)); href != "" {
			links = append(links, reportGalleryLink{"Worked example", href})
		} else if href := reportGallerySourceLink(repositoryRoot, repository, sample.Source); href != "" {
			links = append(links, reportGalleryLink{"Source", href})
		}
		card.Preview = reportGalleryFileLink(repository, path.Join(folder, "preview.png"))
		if sample.Image != "" {
			if image := reportGalleryFileLink(repository, path.Join(folder, sample.Image)); image != "" {
				links = append(links, reportGalleryLink{"Reference image", image})
				if card.Preview == "" {
					card.Preview = image
				}
			}
		}
		if figure := reportGalleryFileLink(repository, path.Join(folder, "preview.svg")); figure != "" {
			links = append(links, reportGalleryLink{"Figure preview", figure})
			if card.Preview == "" {
				card.Preview = figure
			}
		}
		links = append(links, reportGalleryLink{"Sample record", reportGalleryHref(sample.RecordPath)})
		card.Open, card.More = links[0], links[1:]
		page.Samples = append(page.Samples, card)
	}
	sort.SliceStable(page.Samples, func(i, j int) bool {
		return strings.ToLower(page.Samples[i].CopyName) < strings.ToLower(page.Samples[j].CopyName)
	})
	var output bytes.Buffer
	output.WriteString(reportGalleryPageHeader)
	if err := layout.Execute(&output, page); err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(output.Bytes()[len(reportGalleryPageHeader):])) == 0 {
		return nil, fmt.Errorf("gallery layout %q produced an empty page", layout.Name())
	}
	return output.Bytes(), nil
}

func reportGalleryBundles(bundles []reportBundle, worked []workedReportCard) []reportBundle {
	if len(worked) == 0 {
		return bundles
	}
	catalogEntries := make(map[string]bool, len(bundles))
	for _, bundle := range bundles {
		catalogEntries[bundle.Entry] = true
	}
	// Manifests can reference reports outside the catalog's directory layout.
	var manifestOnly []reportBundle
	for _, card := range worked {
		if !catalogEntries[card.ReportPath] {
			folder := path.Dir(card.ReportPath)
			manifestOnly = append(manifestOnly, reportBundle{
				Path: folder, Entry: card.ReportPath, Date: reportBundleDate(path.Base(folder)),
			})
		}
	}
	if len(manifestOnly) == 0 {
		return bundles
	}
	combined := slices.Concat(bundles, manifestOnly)
	sortReportBundles(combined)
	return combined
}

func newReportGalleryStarter(id, identity string, borrow []string) reportGalleryStarter {
	return reportGalleryStarter{ID: id, Prompt: "do-work-toolbox ai-report --design " + quoteDesignName(identity), Borrow: borrow}
}

// Escape shell interpolation inside the quoted design argument in a pasted starter.
func quoteDesignName(name string) string {
	escape := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "$", `\$`, "`", "\\`")
	return `"` + escape.Replace(name) + `"`
}

func reportGalleryKinds(bundles []reportBundle) []reportGalleryKind {
	latest := map[string]reportGalleryLink{}
	for _, bundle := range bundles { // The catalog walk has already sorted newest first.
		if bundle.Entry == "" {
			continue
		}
		kind := bundle.Kind
		if kind == "unknown" {
			kind = "completed-work" // Pre-kind reports and the default ai-report form have no kind meta.
		}
		if _, found := latest[kind]; !found {
			latest[kind] = reportGalleryLink{bundle.Title, reportGalleryHref(bundle.Entry)}
		}
	}
	// These are the report-producing routes with a creation command. The
	// catalog also accepts custom kind metadata, which is included when present.
	var kinds []reportGalleryKind
	for _, kind := range []struct{ name, label, command string }{
		{"completed-work", "Completed work", "do-work-toolbox ai-report <REQ-NNN | UR-NNN>"},
		{"proposal", "Proposal", "do-work-toolbox ai-report --kind proposal <topic | REQ-NNN | UR-NNN>"},
		{"root-cause", "Root cause", "do-work-toolbox ai-report --kind root-cause <topic | REQ-NNN | UR-NNN>"},
		{"architecture-report", "Architecture report", "do-work-toolbox architecture-report"},
	} {
		kinds = append(kinds, reportGalleryKind{Label: kind.label, Command: kind.command, Latest: latest[kind.name]})
		delete(latest, kind.name)
	}
	extraNames := make([]string, 0, len(latest))
	for name := range latest {
		extraNames = append(extraNames, name)
	}
	sort.Strings(extraNames)
	for _, name := range extraNames {
		kinds = append(kinds, reportGalleryKind{Label: name, Latest: latest[name]})
	}
	return kinds
}

func reportGalleryPresets(repository *os.Root, folder string) ([]reportGalleryLink, error) {
	entries, err := fs.ReadDir(repository.FS(), folder)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var links []reportGalleryLink
	for _, entry := range entries {
		name, markdown := strings.CutSuffix(entry.Name(), ".md")
		if entry.Type().IsRegular() && markdown && name != "" && !strings.HasPrefix(name, ".") {
			links = append(links, reportGalleryLink{name, reportGalleryHref(path.Join(folder, entry.Name()))})
		}
	}
	return links, nil
}

func reportGallerySourceLink(repositoryRoot string, repository *os.Root, source string) string {
	if source == "shipped" {
		return ""
	}
	filename := filepath.FromSlash(source)
	if filepath.IsAbs(filename) {
		relative, err := filepath.Rel(repositoryRoot, filename)
		if err != nil {
			return ""
		}
		if !filepath.IsLocal(relative) {
			// Keep an explicitly supplied absolute source usable from disk
			// without reading files outside the repository.
			return reportGalleryHref(filepath.ToSlash(relative))
		}
		filename = relative
	} else if reference, err := url.Parse(source); err == nil && reference.Scheme != "" {
		if (reference.Scheme == "https" || reference.Scheme == "http") && reference.Host != "" {
			return reference.String()
		}
		return ""
	}
	relative := filepath.ToSlash(filepath.Clean(filename))
	if !fs.ValidPath(relative) || relative == "." {
		return ""
	}
	return reportGalleryFileLink(repository, relative)
}

func reportGalleryFileLink(repository *os.Root, relative string) string {
	info, err := repository.Lstat(filepath.FromSlash(relative))
	if err != nil || !info.Mode().IsRegular() {
		return ""
	}
	return reportGalleryHref(relative)
}

func reportGalleryHref(target string) string {
	relative, err := filepath.Rel(filepath.FromSlash(reportGalleryFolder), filepath.FromSlash(target))
	if err != nil {
		return ""
	}
	return (&url.URL{Path: filepath.ToSlash(relative)}).String()
}

// The exact header puts the ownership marker in an actual meta element.
const reportGalleryPageHeader = `<!doctype html>
<html lang="en">
<head>
<meta name="generator" content="do-work-cli ai-report-gallery">
`

const reportGalleryTemplateFile = "report_gallery_page.gohtml"

//go:embed report_gallery_page.gohtml
var reportGalleryPageLayout string

var reportGalleryPageTemplate = template.Must(template.New(reportGalleryTemplateFile).Parse(reportGalleryPageLayout))

func installedReportGalleryTemplate(repository *os.Root) (*template.Template, error) {
	const skillFolder = ".claude/skills/do-work"
	const filename = "tools/do-work-cli/internal/toolboxcommands/" + reportGalleryTemplateFile
	if err := validateNoLinkedAncestors(repository.Name(), skillFolder, true); err != nil {
		return nil, err
	}
	skill, err := rootedfs.OpenRoot(repository, filepath.FromSlash(skillFolder))
	if err != nil {
		return nil, err
	}
	defer skill.Close()
	content, err := readReportGenerationFile(skill, filename)
	if err != nil {
		return nil, err
	}
	return template.New(reportGalleryTemplateFile).Parse(string(content))
}
