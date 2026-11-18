package toolboxcommands

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type reportGalleryLink struct {
	Label string
	Href  string
}

type reportGalleryCard struct {
	Slug     string
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

type reportGalleryGeneratedCard struct {
	Title, UseWhen, Report string
	Starter                reportGalleryStarter
	Review                 *reportGalleryReview
}

type reportGalleryKind struct {
	Label   string
	Command string
	Latest  reportGalleryLink
}

type workedReportGroup struct {
	Label string
	Cards []workedReportCard
}

func groupWorkedReports(cards []workedReportCard) []workedReportGroup {
	var groups []workedReportGroup
	groupByKind := make(map[string]int)
	for _, card := range cards {
		index, exists := groupByKind[card.Kind]
		if !exists {
			label := card.Kind
			switch card.Kind {
			case "Proposal":
				label = "Proposals"
			case "Root cause":
				label = "Diagnose a problem"
			}
			index = len(groups)
			groupByKind[card.Kind] = index
			groups = append(groups, workedReportGroup{Label: label})
		}
		groups[index].Cards = append(groups[index].Cards, card)
	}
	return groups
}

type reportGalleryPageData struct {
	Stylesheet string
	Samples    []reportGalleryCard
	Worked     []workedReportGroup
	Generated  []reportGalleryGeneratedCard
	Older      []reportGalleryGeneratedCard
	Kinds      []reportGalleryKind
	Presets    []reportGalleryLink
	Live       *reportGalleryLivePage
}

func renderReportGalleryPage(repositoryRoot string, repository *os.Root, toolboxDocs string, samples []reportGallerySample, bundles []reportBundle, live *reportGalleryLivePage) ([]byte, error) {
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
	if live != nil {
		if err := live.discover(repositoryRoot, worked, bundles); err != nil {
			return nil, err
		}
	}
	workedByReport := make(map[string]*workedReportCard, len(worked))
	for index, card := range worked {
		workedByReport[card.ReportPath] = &worked[index]
	}
	generatedByReport := make(map[string]reportGalleryGeneratedSource, len(bundles))
	for _, source := range reportGalleryGeneratedSources(repository, bundles) {
		generatedByReport[source.bundle.Entry] = source
	}
	seenReports := make(map[string]bool, len(bundles))
	for index, bundle := range bundles {
		if bundle.Entry == "" || seenReports[bundle.Entry] {
			continue
		}
		seenReports[bundle.Entry] = true
		source, found := generatedByReport[bundle.Entry]
		var design reportGalleryGeneratedDesign
		if found {
			design, err = source.readDesign(repository)
			found = err == nil
		}
		if workedCard := workedByReport[bundle.Entry]; workedCard != nil {
			if found {
				workedCard.Starter = newReportGalleryStarter("worked-prompt-"+workedCard.Name, source.identity, design.borrow)
			}
			continue
		}
		card := reportGalleryGeneratedCard{Title: bundle.Title, Report: reportGalleryHref(bundle.Entry)}
		if live != nil {
			card.Review = live.Reviews[bundle.Entry]
		}
		if found {
			card.Starter = newReportGalleryStarter("generated-prompt-"+strconv.Itoa(index), source.identity, design.borrow)
			if design.recordPath != bundle.Entry {
				card.Title, card.UseWhen = design.title, design.useWhen
				page.Generated = append(page.Generated, card)
				continue
			}
		}
		if found || card.Review != nil {
			page.Older = append(page.Older, card)
		}
	}
	page.Worked = groupWorkedReports(worked)
	copyNames := reportGalleryCopyNames(samples)
	for _, sample := range samples {
		card := reportGalleryCard{Slug: sample.Name, Title: sample.Title, CopyName: copyNames[sample.Name], UseWhen: sample.UseWhen}
		card.Starter = newReportGalleryStarter("prompt-"+card.Slug, card.CopyName, sample.Borrow)
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
	if err := reportGalleryPageTemplate.Execute(&output, page); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
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

var reportGalleryPageTemplate = template.Must(template.Must(template.New("report-gallery").Parse(reportGalleryFeedbackTemplates)).Parse(reportGalleryPageHeader + `<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>AI Report Gallery</title>
{{if .Stylesheet}}<link rel="stylesheet" href="{{.Stylesheet}}">{{end}}
<style>
* { box-sizing: border-box; }
:root {
  color-scheme: light dark;
  --ground: #faf9f5; --panel: #fffefa; --text: #232b28; --subtle: #535d56;
  --edge: #cccfc5; --accent: #285949; --wash: #e9eee7;
  --sans: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --serif: "Iowan Old Style", "Baskerville", Georgia, serif;
}
@media (prefers-color-scheme: dark) {
  :root { --ground: #171c1a; --panel: #1e2522; --text: #f0f2eb; --subtle: #b7c0b9; --edge: #46524b; --accent: #b1dfc7; --wash: #293a31; }
}
html, body { background: var(--ground); color: var(--text); }
body { margin: 0; font: 16px/1.5 var(--sans); }
.gallery { width: 100%; padding: 0 clamp(16px, 3.2vw, 76px) 3rem; }
.gallery h1, .gallery h2, .gallery h3 { color: var(--text); font-family: var(--serif); font-weight: 500; line-height: 1.1; }
.gallery h1 { margin: 0; font-size: clamp(2.3rem, 3.4vw, 3.6rem); letter-spacing: -.045em; }
.gallery h2 { margin: 0 0 1rem; font-size: clamp(1.6rem, 2vw, 2.1rem); letter-spacing: -.025em; }
.gallery a { color: var(--accent); text-underline-offset: .2em; }
.gallery :focus-visible { outline: 3px solid var(--accent); outline-offset: 4px; }
.gallery-header { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 1rem; padding: clamp(16px, 1.6vw, 30px) 0 1rem; border-bottom: 2px solid var(--text); }
.gallery-nav, .gallery-links { display: flex; flex-wrap: wrap; align-items: center; gap: .5rem 1.15rem; }
.gallery-nav a { font-size: .95rem; font-weight: 600; text-decoration: none; padding-block: .45rem; }
.gallery-nav a:hover { text-decoration: underline; }
.gallery-nav .add-link, .gallery-copy button { background: var(--text); color: var(--ground); border: 1px solid var(--text); border-radius: .4rem; padding: .65rem 1rem; font: 600 .95rem var(--sans); cursor: pointer; }
.gallery section { scroll-margin-top: 1rem; }
.samples { padding: 1.1rem 0 2rem; }
{{if .Worked}}.report-group + .report-group { margin-top: 1.75rem; }
.report-group > h3 { margin: 0 0 .75rem; font: 700 1rem/1.4 var(--sans); color: var(--accent); }
{{end}}.gallery-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 1.8rem clamp(16px, 1.55vw, 32px); align-items: start; }
.gallery-card { min-width: 0; scroll-margin-top: 1rem; }
.gallery-card:target { outline: 2px solid var(--accent); outline-offset: 6px; }
.sample-cover { display: block; position: relative; background: var(--panel); border: 1px solid var(--edge); border-radius: .4rem; overflow: hidden; }
.sample-cover:hover { border-color: var(--accent); }
.sample-cover span { position: absolute; bottom: .5rem; right: .5rem; border: 1px solid #c4c9c3; background: #fff; color: #20382c; padding: .25rem .55rem; border-radius: .25rem; font: 600 .78rem/1.4 var(--sans); }
.sample-preview { display: block; width: 100%; height: auto; aspect-ratio: 1040 / 740; object-fit: cover; object-position: top; }
.card-heading { display: flex; justify-content: space-between; gap: .5rem; align-items: baseline; margin: .55rem 0 .1rem; }
.gallery .card-heading h3 { margin: 0; min-width: 0; font-size: clamp(1.5rem, 2vw, 2.15rem); letter-spacing: -.025em; }
{{if or .Generated .Older}}.generated-design .card-heading h3 { overflow-wrap: anywhere; }
{{end}}{{if .Older}}.report-archive { margin-top: 1rem; }
.report-archive > summary { cursor: pointer; }
{{end}}.name-copy { display: inline-flex; align-items: center; gap: .4rem; max-width: 100%; background: none; color: inherit; font: inherit; text-align: left; border: 0; padding: .15rem 0; cursor: copy; overflow-wrap: anywhere; }
.name-copy:hover { color: var(--accent); text-decoration: underline; text-underline-offset: .14em; }
.name-copy:disabled { opacity: 1; cursor: text; }
.copy-icon { position: relative; display: inline-block; flex: 0 0 14px; height: 15px; color: var(--subtle); }
.copy-icon::before, .copy-icon::after { content: ""; position: absolute; border: 1.5px solid; border-radius: 2px; width: 8px; height: 10px; top: 0; left: 0; }
.copy-icon::after { left: 4px; top: 3px; background: var(--ground); }
.report-link { flex: none; font: 600 .85rem/1.4 var(--sans); padding: .55rem 0; }
.gallery-card p { margin: 0 0 .45rem; color: var(--subtle); line-height: 1.45; max-width: 52ch; }
{{if not .Worked}}.worked-report .worked-kind { font: 700 .84rem/1.4 var(--sans); color: var(--accent); }
{{end}}.gallery-card details { background: none; border: 0; padding: 0; border-radius: 0; }
.gallery-card summary { display: list-item; font: 600 .9rem var(--sans); width: fit-content; color: var(--accent); padding: .65rem 0; cursor: pointer; }
.gallery-card details[open] summary { border: 0; margin: 0; }
.gallery-card ul { padding-left: 1.2rem; margin: .25rem 0 .6rem; }
.gallery-card details a { font-size: .9rem; }
.section-note { display: flex; align-items: baseline; justify-content: space-between; flex-wrap: wrap; gap: .25rem 1rem; margin-bottom: .85rem; }
.section-note h2 { margin: 0; }
.section-note p { margin: 0; color: var(--subtle); font-size: .9rem; }
.design-prompt { display: block; margin: .55rem 0; padding: .6rem; background: var(--wash); color: var(--text); overflow-wrap: anywhere; }
.prompt-copy { font: 600 .85rem var(--sans); padding: .45rem .65rem; background: var(--panel); border: 1px solid var(--edge); color: var(--accent); border-radius: .3rem; cursor: pointer; }
{{if .Samples}}.gallery-dialog { padding: 0; width: min(68rem, calc(100% - 24px)); max-width: none; max-height: 90vh; background: var(--panel); color: var(--text); border: 1px solid var(--edge); border-radius: .5rem; }
.gallery-dialog[open] { display: flex; flex-direction: column; }
.gallery-dialog::backdrop { background: rgb(14 22 18 / .78); }
.dialog-bar { flex: none; display: flex; align-items: center; justify-content: space-between; gap: .5rem; padding: .55rem .9rem; border-bottom: 1px solid var(--edge); }
.dialog-bar h2 { font-size: 1.15rem; margin: 0; }
.dialog-bar a { margin-left: auto; }
.preview-scroll { overflow: auto; overscroll-behavior: contain; }
.preview-scroll img { display: block; width: 100%; height: auto; }
.preview-scroll.actual-size img { width: auto; max-width: none; }
#preview-size[aria-pressed="true"] { background: var(--text); color: var(--ground); }
@media (max-width: 520px) { .dialog-bar { flex-wrap: wrap; } .dialog-bar h2 { width: 100%; } .dialog-bar a { margin-left: 0; margin-right: auto; } }
{{end}}.copy-status:not(:empty) { pointer-events: none; position: fixed; bottom: 1rem; left: 50%; transform: translateX(-50%); z-index: 2; color: var(--ground); background: var(--text); border-radius: .4rem; padding: .6rem 1rem; max-width: calc(100% - 2rem); width: max-content; }
.gallery-tools { display: grid; grid-template-columns: 3fr 2fr; gap: clamp(1.5rem, 4vw, 5rem); border-block: 1px solid var(--edge); padding: 1.65rem 0 2rem; }
{{if not .Worked}}.gallery-kinds, {{end}}.gallery-presets { list-style: none; margin: 0; padding: 0; }
{{if not .Worked}}.gallery-kinds li { border-top: 1px solid var(--edge); padding: .55rem 0; display: grid; grid-template-columns: minmax(7rem, 27%) minmax(0, 1fr); gap: .4rem 1rem; }
{{end}}{{if not .Worked}}.gallery-kinds code, {{end}}.gallery-copy code { background: none; color: var(--text); padding: 0; overflow-wrap: anywhere; }
.gallery-presets { display: flex; flex-wrap: wrap; gap: .45rem; }
.gallery-presets a { display: block; padding: .4rem .75rem; background: var(--wash); border-radius: .35rem; font-size: .9rem; text-decoration: none; }
.gallery-presets a:hover { text-decoration: underline; }
.gallery-add { padding: 1.6rem 0; }
.gallery-add p { color: var(--subtle); margin: 0 0 .75rem; }
.gallery-copy { display: flex; align-items: center; flex-wrap: wrap; gap: .8rem; }
.gallery-copy button[hidden] { display: none; }
@media (min-width: 1700px) { .gallery-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); } }
@media (max-width: 1000px) { .gallery-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 720px) { .gallery-tools { grid-template-columns: 1fr; } }
@media (max-width: 560px) {
  .gallery-grid { grid-template-columns: minmax(0, 1fr); gap: 1.6rem; }
  .gallery-header { align-items: start; flex-direction: column; gap: .5rem; }
  .gallery-nav { gap: .5rem .9rem; }
{{if not .Worked}}  .gallery-kinds li { grid-template-columns: 1fr; }
{{end}}}
</style>{{if .Live}}{{template "galleryFeedbackStyles"}}{{end}}
</head>
<body>
<main class="gallery">
  <header class="gallery-header">
    <h1>AI Report Gallery</h1>
    <nav class="gallery-nav" aria-label="Gallery sections">
{{if not .Worked}}      <a href="#report-kinds">Report kinds</a>
{{end}}{{if and .Worked .Samples}}      <a href="#samples">Samples</a>
{{end}}{{if or .Generated .Older}}      <a href="#generated-designs">Generated designs</a>
{{end}}      <a href="#presets">Presets</a>
      <a class="add-link" href="#add-sample">Add a sample</a>
    </nav>
  </header>
  {{if .Worked}}<section class="samples" id="worked-reports" aria-labelledby="worked-title">
    <div class="section-note"><h2 id="worked-title">Worked reports</h2>
      <p>Open a report from its preview. Click a name to copy it.</p></div>
    {{range $index, $group := .Worked}}{{if $index}}    {{end}}<section class="report-group" aria-label="{{.Label}}">
      <h3>{{.Label}}</h3>
      <div class="gallery-grid">
{{range .Cards}}        <article class="gallery-card worked-report" data-gallery-card id="{{.Name}}">
        <a class="sample-cover" data-gallery-preview href="{{.Report}}" aria-label="Open {{.Subject}}">
          <img class="sample-preview" src="{{.Image}}" loading="lazy" alt="" width="{{.Width}}" height="{{.Height}}" style="aspect-ratio: {{.Width}} / {{.Height}}">
        </a>
        <div class="card-heading"><h3><button class="name-copy" data-gallery-name type="button" disabled data-copy="worked-name-{{.Name}}" title="Copy report name"><span id="worked-name-{{.Name}}">{{.Name}}</span><i class="copy-icon" aria-hidden="true"></i></button></h3>
          <a class="report-link" href="#{{.Name}}" aria-label="Permalink to {{.Name}}">#</a></div>
        <p>{{.Subject}}</p>
        {{if .Starter.Prompt}}<details><summary>Use this design</summary>
          {{template "starter-prompt" .Starter}}
        </details>{{else}}<p>Design unavailable.</p>{{end}}{{if $.Live}}{{with index $.Live.Reviews .ReportPath}}{{template "galleryFeedbackReview" .}}{{end}}{{end}}
      </article>
{{end}}      </div>
    </section>
{{end}}  </section>{{end}}
  {{if or (not .Worked) .Samples}}<section class="samples" id="samples" aria-labelledby="samples-title">
    <div class="section-note"><h2 id="samples-title">{{if .Worked}}Sample designs{{else}}{{len .Samples}} report designs{{end}}</h2>
      {{if .Samples}}<p>Select a page to preview it. Click a design name to copy it.</p>{{end}}</div>
    {{if .Samples}}<div class="gallery-grid">
      {{range .Samples}}<article class="gallery-card" id="sample-{{.Slug}}">
        {{if .Preview}}<a class="sample-cover" href="{{.Preview}}" aria-label="Preview {{.Title}}" data-preview>
          <img class="sample-preview" src="{{.Preview}}" loading="lazy" alt="" width="1040" height="740">
          <span aria-hidden="true">Enlarge ↗</span></a>{{end}}
        <div class="card-heading"><h3><button class="name-copy" type="button" disabled data-copy="name-{{.Slug}}" title="Copy design name" aria-label="Copy design name {{.CopyName}}"><span id="name-{{.Slug}}">{{.CopyName}}</span><i class="copy-icon" aria-hidden="true"></i></button></h3>
          <a class="report-link" href="{{.Open.Href}}" aria-label="{{.Open.Label}} — {{.Title}}">{{if eq .Open.Label "Worked example"}}Open report{{else}}{{.Open.Label}}{{end}} ↗</a></div>
        {{if ne .CopyName .Title}}<p>{{.Title}}</p>{{end}}<p>{{.UseWhen}}</p>
        <details><summary>Use this design</summary>
          {{template "starter-prompt" .Starter}}
          <div class="gallery-links"><a href="#sample-{{.Slug}}">Permalink</a>{{range .More}}<a href="{{.Href}}">{{.Label}}</a>{{end}}</div>
        </details>
      </article>{{end}}
    </div>{{else}}<p>No samples yet.</p>{{end}}
  </section>
  {{end}}{{if or .Generated .Older}}<section class="samples" id="generated-designs" aria-labelledby="generated-designs-title">
    <div class="section-note"><h2 id="generated-designs-title">Generated report designs</h2></div>
    <div{{if .Live}} id="project-reports"{{end}}>
{{if .Generated}}      <div class="gallery-grid">{{range .Generated}}{{template "generated-report-card" .}}{{end}}</div>
{{end}}{{if .Older}}      <details class="report-archive">
        <summary>Browse other reports ({{len .Older}})</summary>
        <div class="gallery-grid">{{range .Older}}{{template "generated-report-card" .}}{{end}}</div>
      </details>{{end}}
    </div>
  </section>
  {{end}}{{if not .Worked}}<div class="gallery-tools">
  <section id="report-kinds" aria-labelledby="report-kinds-title">
    <h2 id="report-kinds-title">Report kinds</h2>
    <ul class="gallery-kinds">
      {{range .Kinds}}<li><strong>{{.Label}}</strong>
        {{if .Latest.Href}}<a href="{{.Latest.Href}}">{{.Latest.Label}}</a>{{else}}<span>Create one with <code>{{.Command}}</code></span>{{end}}
      </li>{{end}}
    </ul>
  </section>
{{end}}{{if not .Worked}}  {{end}}<section id="presets" aria-labelledby="presets-title">
    <h2 id="presets-title">Second-opinion presets</h2>
    {{if .Presets}}<ul class="gallery-links gallery-presets">
      {{range .Presets}}<li><a href="{{.Href}}">{{.Label}}</a></li>{{end}}
    </ul>{{else}}<p>No shipped presets found.</p>{{end}}
  </section>
{{/* Keep the closing div byte from the scored worked-report page. */}}  </div>
  <section class="gallery-add" id="add-sample" aria-labelledby="add-sample-title">
    <h2 id="add-sample-title">Add a sample</h2>
    <p>Paste this line into your agent, replacing the placeholders with the source and what you liked.</p>
    <div class="gallery-copy">
      <code id="add-line">do-work-toolbox ai-report gallery add &lt;image path | link | html path&gt; [what the user liked]</code>
      <button class="btn" id="copy-add-line" data-copy="add-line" type="button" hidden>Copy command</button>
    </div>
  </section>
  <div id="copy-status" class="copy-status" role="status"></div>
{{if .Samples}}  <dialog class="gallery-dialog" aria-labelledby="preview-title">
    <div class="dialog-bar"><h2 id="preview-title">Page preview</h2><a class="report-link" id="preview-open"></a><button class="prompt-copy" type="button" id="preview-size" aria-pressed="false">Actual size</button><button class="prompt-copy" type="button" id="preview-close" autofocus>Close</button></div>
    <div class="preview-scroll" tabindex="0" aria-label="Page image"></div>
  </dialog>
{{end}}</main>
<script>
const copyStatus = document.getElementById('copy-status');
let statusTimer;
function announce(message) {
  clearTimeout(statusTimer);
  copyStatus.textContent = message;
  statusTimer = setTimeout(() => { copyStatus.textContent = ''; }, 3200);
}
async function copyText(element) {
  if (location.protocol !== 'file:' && navigator.clipboard) {
    try { await navigator.clipboard.writeText(element.textContent); return true; } catch {}
  }
  const range = document.createRange();
  range.selectNodeContents(element);
  const selection = window.getSelection();
  selection.removeAllRanges();
  selection.addRange(range);
  let copied = false;
  try { copied = document.execCommand('copy'); } catch {}
  if (copied) selection.removeAllRanges();
  return copied;
}
for (const button of document.querySelectorAll('[data-copy]')) {
  button.hidden = false;
  button.disabled = false;
  button.addEventListener('click', async () => {
    const target = document.getElementById(button.dataset.copy);
    const success = await copyText(target);
    const message = button.classList.contains('name-copy') ? 'Copied “' + target.textContent + '”.' : 'Copied.';
    announce(success ? message : 'Text selected. Use your keyboard to copy.');
  });
}
{{if .Samples}}const dialog = document.querySelector('.gallery-dialog');
if (typeof dialog.showModal === 'function') {
  const image = document.createElement('img');
  image.id = 'preview-image';
  const open = document.getElementById('preview-open');
  const scroller = dialog.querySelector('.preview-scroll');
  const size = document.getElementById('preview-size');
  size.addEventListener('click', () => {
    size.setAttribute('aria-pressed', scroller.classList.toggle('actual-size'));
  });
  for (const cover of document.querySelectorAll('[data-preview]')) {
    cover.addEventListener('click', event => {
      if (event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
      event.preventDefault();
      const card = cover.closest('article');
      const report = card.querySelector('.report-link');
      document.getElementById('preview-title').textContent = card.querySelector('h3 span').textContent;
      open.href = report.href;
      open.textContent = report.textContent;
      image.src = cover.href;
      image.alt = cover.getAttribute('aria-label');
      if (!image.isConnected) scroller.append(image);
      scroller.classList.remove('actual-size');
      size.setAttribute('aria-pressed', 'false');
      dialog.showModal();
      scroller.scrollTo(0, 0);
    });
  }
  document.getElementById('preview-close').addEventListener('click', () => dialog.close());
  dialog.addEventListener('click', event => {
    if (event.target !== dialog) return;
    const box = dialog.getBoundingClientRect();
    if (event.clientX < box.left || event.clientX > box.right || event.clientY < box.top || event.clientY > box.bottom) dialog.close();
  });
}
{{end}}</script>{{if .Live}}{{template "galleryFeedbackScript"}}{{end}}
</body>
</html>
{{define "starter-prompt" -}}
<code class="design-prompt" id="{{.ID}}">{{.Prompt}}</code>
          <button class="prompt-copy" type="button" data-copy="{{.ID}}" hidden>Copy starter prompt</button>{{if .Borrow}}
          <ul aria-label="Borrow">{{range .Borrow}}<li>{{.}}</li>{{end}}</ul>{{end}}
{{- end -}}
{{define "generated-report-card"}}
      <article class="gallery-card generated-design">
        <div class="card-heading"><h3>{{.Title}}</h3>
          <a class="report-link" href="{{.Report}}" aria-label="Open report — {{.Title}}">Open report ↗</a></div>
{{if .UseWhen}}        <p>{{.UseWhen}}</p>
{{end -}}
        {{if .Starter.Prompt}}<details><summary>Use this design</summary>
          {{template "starter-prompt" .Starter}}
        </details>{{else}}<p>Design unavailable.</p>{{end}}{{with .Review}}{{template "galleryFeedbackReview" .}}{{end}}
      </article>
{{end -}}
`))

const reportGalleryFeedbackTemplates = `
{{define "galleryFeedbackStyles"}}
<style>
.gallery .review-details { margin-top: .55rem; border-top: 1px solid var(--edge); }
.gallery .review-details summary { min-height: 36px; }
.review-form { margin: 0 0 .8rem; }
.review-form fieldset { min-width: 0; margin: 0; padding: .5rem .65rem; border: 1px solid var(--edge); border-radius: .35rem; }
.review-form legend { max-width: 100%; padding: 0 .2rem; font: 600 .9rem/1.35 var(--sans); overflow-wrap: anywhere; }
.review-mark { display: flex; flex-wrap: wrap; column-gap: .75rem; margin-bottom: .25rem; }
.review-mark label { display: inline-flex; gap: .28rem; align-items: center; min-height: 32px; cursor: pointer; }
.review-mark input { margin: 0; accent-color: var(--accent); }
.review-note { display: grid; gap: .25rem; font-size: .85rem; }
.review-note textarea { width: 100%; min-height: 3rem; resize: vertical; padding: .4rem; font: inherit; color: var(--text); background: var(--panel); border: 1px solid var(--edge); border-radius: .25rem; }
.review-save { display: flex; flex-wrap: wrap; align-items: baseline; gap: .35rem .65rem; margin-top: .5rem; }
.review-save button { min-height: 32px; }
.review-status { color: var(--subtle); font-size: .85rem; overflow-wrap: anywhere; }
.review-status[data-error] { font-weight: 600; }
</style>
{{end}}
{{define "galleryFeedbackReview"}}
        <details class="review-details"><summary>Review report and figures</summary>
          {{template "galleryFeedbackForm" .}}
          {{range .Figures}}
            <form class="review-form" data-report-feedback data-target-kind="figure" data-report-path="{{$.ReportPath}}" data-figure-id="{{.ID}}">
              <fieldset disabled>
                <legend>{{if .Href}}<a href="{{.Href}}">{{.Label}}</a>{{else}}{{.Label}}{{end}}</legend>
                {{template "galleryFeedbackFields"}}
              </fieldset>
              <div class="review-status" role="status">Loading feedback…</div>
            </form>
          {{end}}
        </details>{{end}}
{{define "galleryFeedbackForm"}}
          <form class="review-form" data-report-feedback data-target-kind="report" data-report-path="{{.ReportPath}}">
            <fieldset disabled>
              <legend>Whole report</legend>
              {{template "galleryFeedbackFields"}}
            </fieldset>
            <div class="review-status" role="status">Loading feedback…</div>
          </form>{{end}}
{{define "galleryFeedbackFields"}}
                <div class="review-mark">
                  <label><input type="radio" name="mark" value="like" required> Like</label>
                  <label><input type="radio" name="mark" value="dislike" required> Dislike</label>
                </div>
                <label class="review-note">Note (optional)<textarea name="note" rows="2" maxlength="500"></textarea></label>
                <div class="review-save"><button class="prompt-copy" type="submit">Save feedback</button></div>
{{end}}
{{define "galleryFeedbackScript"}}
<script>
(() => {
  const forms = [...document.querySelectorAll('[data-report-feedback]')];
  if (!forms.length) return;
  const endpoint = "/api/report-gallery-feedback";
  const identity = entry => JSON.stringify([entry.target_kind, entry.report_path, entry.figure_id || '']);
  async function readResponse(response) {
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error || 'The gallery could not save or load feedback.');
    return payload;
  }
  fetch(endpoint, {mode: 'same-origin', cache: 'no-store'}).then(readResponse).then(feedback => {
    const savedEntries = new Map(feedback.entries.map(entry => [identity(entry), entry]));
    for (const form of forms) {
      const target = {target_kind: form.dataset.targetKind, report_path: form.dataset.reportPath};
      if (form.dataset.figureId) target.figure_id = form.dataset.figureId;
      const saved = savedEntries.get(identity(target));
      const fields = form.querySelector('fieldset');
      const status = form.querySelector('[role="status"]');
      form.elements.mark.value = saved?.mark || '';
      form.elements.note.value = saved?.note || '';
      status.textContent = saved ? 'Saved ' + saved.mark + '.' : '';
      fields.disabled = false;
      form.addEventListener('submit', async event => {
        event.preventDefault();
        const entry = {...target, mark: form.elements.mark.value, note: form.elements.note.value};
        fields.disabled = true;
        delete status.dataset.error;
        status.textContent = 'Saving…';
        try {
          const stored = await readResponse(await fetch(endpoint, {
            method: 'PUT', mode: 'same-origin',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(entry),
          }));
          if (identity(stored) !== identity(entry) || stored.mark !== entry.mark || (stored.note || '') !== entry.note) {
            throw new Error('The gallery did not confirm this entry.');
          }
          status.textContent = 'Saved ' + stored.mark + '.';
        } catch (error) {
          status.dataset.error = '';
          status.textContent = 'Feedback could not be saved. ' + error.message + ' Use Save feedback to retry.';
        } finally {
          fields.disabled = false;
        }
      });
    }
  }).catch(error => {
    for (const form of forms) {
      const status = form.querySelector('[role="status"]');
      status.dataset.error = '';
      status.textContent = 'Feedback could not be loaded. ' + error.message + ' Reopen or reload the gallery to retry.';
    }
  });
})();
</script>{{end}}
`
