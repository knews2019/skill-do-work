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
	"strings"
)

type reportGalleryLink struct {
	Label string
	Href  string
}

type reportGalleryCard struct {
	Title   string
	UseWhen string
	Borrow  []string
	Preview string
	Open    reportGalleryLink
	More    []reportGalleryLink
}

type reportGalleryKind struct {
	Label   string
	Command string
	Latest  reportGalleryLink
}

type reportGalleryPageData struct {
	Stylesheet string
	Samples    []reportGalleryCard
	Kinds      []reportGalleryKind
	Presets    []reportGalleryLink
}

func renderReportGalleryPage(repositoryRoot string, repository *os.Root, toolboxDocs string, samples []reportGallerySample, bundles []reportBundle) ([]byte, error) {
	page := reportGalleryPageData{Kinds: reportGalleryKinds(bundles)}
	if toolboxDocs != "" {
		stylesheet := path.Join(toolboxDocs, "ai-report-kit", "presentation.css")
		page.Stylesheet = reportGalleryFileLink(repository, stylesheet)
		var err error
		page.Presets, err = reportGalleryPresets(repository, path.Join(toolboxDocs, secondOpinionPresetsFolder))
		if err != nil {
			return nil, err
		}
	}
	for _, sample := range samples {
		card := reportGalleryCard{Title: sample.Title, UseWhen: sample.UseWhen, Borrow: sample.Borrow}
		folder := path.Dir(sample.RecordPath)
		var links []reportGalleryLink
		if href := reportGalleryFileLink(repository, path.Join(folder, reportGalleryExample)); href != "" {
			links = append(links, reportGalleryLink{"Worked example", href})
		}
		if sample.Image != "" {
			card.Preview = reportGalleryFileLink(repository, path.Join(folder, sample.Image))
			if card.Preview != "" {
				links = append(links, reportGalleryLink{"Reference image", card.Preview})
			}
		}
		if card.Preview == "" {
			card.Preview = reportGalleryFileLink(repository, path.Join(folder, "preview.svg"))
		}
		if len(links) == 0 {
			if href := reportGallerySourceLink(repositoryRoot, repository, sample.Source); href != "" {
				links = append(links, reportGalleryLink{"Source", href})
			}
		}
		links = append(links, reportGalleryLink{"Sample record", reportGalleryHref(sample.RecordPath)})
		card.Open, card.More = links[0], links[1:]
		page.Samples = append(page.Samples, card)
	}
	var output bytes.Buffer
	if err := reportGalleryPageTemplate.Execute(&output, page); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
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

var reportGalleryPageTemplate = template.Must(template.New("report-gallery").Parse(reportGalleryPageHeader + `<meta charset="utf-8">
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
.gallery-nav a:hover, .sample-open:hover h3 { text-decoration: underline; }
.gallery-nav .add-link, .gallery-copy button { background: var(--text); color: var(--ground); border: 1px solid var(--text); border-radius: .4rem; padding: .65rem 1rem; font: 600 .95rem var(--sans); cursor: pointer; }
.gallery section { scroll-margin-top: 1rem; }
.samples { padding: 1.1rem 0 2rem; }
.gallery-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: clamp(14px, 1.55vw, 32px); align-items: start; }
.gallery-card { min-width: 0; }
.sample-open { display: block; text-decoration: none; }
.sample-preview { display: block; aspect-ratio: 480 / 284; width: 100%; height: auto; object-fit: contain; background: var(--panel); border: 1px solid var(--edge); border-radius: .5rem; }
.sample-open:hover .sample-preview { border-color: var(--accent); }
.sample-open h3 { display: flex; justify-content: space-between; align-items: baseline; gap: .5rem; margin: .8rem 0 .35rem; font-size: clamp(1.42rem, 1.8vw, 2.2rem); letter-spacing: -.025em; }
.sample-open h3 span { font: 400 1.15rem var(--sans); color: var(--accent); }
.gallery-card p { margin: 0 0 .45rem; color: var(--subtle); line-height: 1.45; max-width: 52ch; }
.gallery-card details { background: none; border: 0; padding: 0; border-radius: 0; }
.gallery-card summary { display: list-item; font: 600 .9rem var(--sans); width: fit-content; color: var(--accent); padding: .65rem 0; cursor: pointer; }
.gallery-card details[open] summary { border: 0; margin: 0; }
.gallery-card ul { padding-left: 1.2rem; margin: .25rem 0 .6rem; }
.gallery-card details a { font-size: .9rem; }
.gallery-tools { display: grid; grid-template-columns: 3fr 2fr; gap: clamp(1.5rem, 4vw, 5rem); border-block: 1px solid var(--edge); padding: 1.65rem 0 2rem; }
.gallery-kinds, .gallery-presets { list-style: none; margin: 0; padding: 0; }
.gallery-kinds li { border-top: 1px solid var(--edge); padding: .55rem 0; display: grid; grid-template-columns: minmax(7rem, 27%) minmax(0, 1fr); gap: .4rem 1rem; }
.gallery-kinds code, .gallery-copy code { background: none; color: var(--text); padding: 0; overflow-wrap: anywhere; }
.gallery-presets { display: flex; flex-wrap: wrap; gap: .45rem; }
.gallery-presets a { display: block; padding: .4rem .75rem; background: var(--wash); border-radius: .35rem; font-size: .9rem; text-decoration: none; }
.gallery-presets a:hover { text-decoration: underline; }
.gallery-add { padding: 1.6rem 0; }
.gallery-add p { color: var(--subtle); margin: 0 0 .75rem; }
.gallery-copy { display: flex; align-items: center; flex-wrap: wrap; gap: .8rem; }
.gallery-copy button[hidden] { display: none; }
@media (max-width: 1100px) { .gallery-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 720px) { .gallery-tools { grid-template-columns: 1fr; } }
@media (max-width: 560px) {
  .gallery-grid { grid-template-columns: minmax(0, 1fr); gap: 1.6rem; }
  .gallery-header { align-items: start; flex-direction: column; gap: .5rem; }
  .gallery-nav { gap: .5rem .9rem; }
  .gallery-kinds li { grid-template-columns: 1fr; }
}
</style>
</head>
<body>
<main class="gallery">
  <header class="gallery-header">
    <h1>AI Report Gallery</h1>
    <nav class="gallery-nav" aria-label="Gallery sections">
      <a href="#report-kinds">Report kinds</a>
      <a href="#presets">Presets</a>
      <a class="add-link" href="#add-sample">Add a sample</a>
    </nav>
  </header>
  <section class="samples" id="samples" aria-labelledby="samples-title">
    <h2 id="samples-title">Presentation samples</h2>
    {{if .Samples}}<div class="gallery-grid">
      {{range .Samples}}<article class="gallery-card">
        <a class="sample-open" href="{{.Open.Href}}" aria-label="{{.Title}} — {{.Open.Label}}">
          {{if .Preview}}<img class="sample-preview" src="{{.Preview}}" alt="" width="480" height="284">{{end}}
          <h3>{{.Title}}<span aria-hidden="true">→</span></h3>
        </a>
        <p>{{.UseWhen}}</p>
        <details><summary>What to borrow</summary>
          <ul aria-label="Borrow">{{range .Borrow}}<li>{{.}}</li>{{end}}</ul>
          {{if .More}}<div class="gallery-links">{{range .More}}<a href="{{.Href}}">{{.Label}}</a>{{end}}</div>{{end}}
        </details>
      </article>{{end}}
    </div>{{else}}<p>No samples yet.</p>{{end}}
  </section>
  <div class="gallery-tools">
  <section id="report-kinds" aria-labelledby="report-kinds-title">
    <h2 id="report-kinds-title">Report kinds</h2>
    <ul class="gallery-kinds">
      {{range .Kinds}}<li><strong>{{.Label}}</strong>
        {{if .Latest.Href}}<a href="{{.Latest.Href}}">{{.Latest.Label}}</a>{{else}}<span>Create one with <code>{{.Command}}</code></span>{{end}}
      </li>{{end}}
    </ul>
  </section>
  <section id="presets" aria-labelledby="presets-title">
    <h2 id="presets-title">Second-opinion presets</h2>
    {{if .Presets}}<ul class="gallery-links gallery-presets">
      {{range .Presets}}<li><a href="{{.Href}}">{{.Label}}</a></li>{{end}}
    </ul>{{else}}<p>No shipped presets found.</p>{{end}}
  </section>
  </div>
  <section class="gallery-add" id="add-sample" aria-labelledby="add-sample-title">
    <h2 id="add-sample-title">Add a sample</h2>
    <p>Paste this line into your agent, replacing the placeholders with the source and what you liked.</p>
    <div class="gallery-copy">
      <code id="add-line">do-work-toolbox ai-report gallery add &lt;image path | link | html path&gt; [what the user liked]</code>
      <button class="btn" id="copy-add-line" type="button" hidden>Copy command</button>
      <span id="copy-status" role="status"></span>
    </div>
  </section>
</main>
<script>
const addLine = document.getElementById('add-line');
const copyButton = document.getElementById('copy-add-line');
const copyStatus = document.getElementById('copy-status');
copyButton.hidden = false;
copyButton.addEventListener('click', async () => {
  if (location.protocol !== 'file:' && navigator.clipboard) {
    try {
      await navigator.clipboard.writeText(addLine.textContent);
      copyStatus.textContent = 'Copied.';
      return;
    } catch {}
  }
  const range = document.createRange();
  range.selectNodeContents(addLine);
  const selection = window.getSelection();
  selection.removeAllRanges();
  selection.addRange(range);
  let copied = false;
  try { copied = document.execCommand('copy'); } catch {}
  if (copied) selection.removeAllRanges();
  copyStatus.textContent = copied ? 'Copied.' : 'Command selected. Use your keyboard to copy.';
});
</script>
</body>
</html>
`))
