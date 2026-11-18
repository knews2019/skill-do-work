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
	Links   []reportGalleryLink
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
		if href := reportGalleryFileLink(repository, path.Join(folder, reportGalleryExample)); href != "" {
			card.Links = append(card.Links, reportGalleryLink{"Worked example", href})
		}
		if sample.Image != "" {
			if href := reportGalleryFileLink(repository, path.Join(folder, sample.Image)); href != "" {
				card.Links = append(card.Links, reportGalleryLink{"Reference image", href})
			}
		}
		if len(card.Links) == 0 {
			if href := reportGallerySourceLink(repositoryRoot, repository, sample.Source); href != "" {
				card.Links = append(card.Links, reportGalleryLink{"Source", href})
			}
		}
		card.Links = append(card.Links, reportGalleryLink{"Sample record", reportGalleryHref(sample.RecordPath)})
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
html { color-scheme: light dark; }
.gallery { max-width: 72rem; margin: 0 auto; padding: 0 clamp(1rem, 4vw, 3rem) 4rem; }
.gallery .band { padding: 2rem 0; }
.gallery-header { padding: 3rem 0 1rem; }
.gallery-links { display: flex; flex-wrap: wrap; gap: .5rem 1.25rem; }
.gallery-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 19rem), 1fr)); gap: 1rem; }
.gallery-card { display: flex; flex-direction: column; }
.gallery-card ul { margin: .25rem 0 1rem; padding-inline-start: 1.25em; }
.gallery-card .gallery-links { margin: auto 0 0; font-size: .95rem; }
.gallery-kinds { list-style: none; margin: 0; padding: 0; }
.gallery-kinds li { display: grid; grid-template-columns: minmax(9rem, 13rem) 1fr; gap: .25rem 1rem; padding: .65rem 0; border-bottom: 1px solid var(--line, GrayText); }
.gallery-presets { list-style: none; padding: 0; }
.gallery-copy { display: flex; align-items: center; flex-wrap: wrap; gap: .75rem; }
.gallery-copy code { padding: .35em .5em; overflow-wrap: anywhere; }
.gallery-copy button[hidden] { display: none; }
@media (max-width: 40rem) { .gallery-kinds li { grid-template-columns: 1fr; } }
</style>
</head>
<body>
<main class="gallery">
  <header class="gallery-header">
    <h1>AI Report Gallery</h1>
    <nav class="gallery-links" aria-label="Gallery sections">
      <a href="#samples">Samples</a>
      <a href="#report-kinds">Report kinds</a>
      <a href="#presets">Second-opinion presets</a>
      <a href="#add-sample">Add a sample</a>
    </nav>
  </header>
  <section class="band" id="samples" aria-labelledby="samples-title">
    <h2 id="samples-title">Samples</h2>
    {{if .Samples}}<div class="gallery-grid">
      {{range .Samples}}<article class="node gallery-card">
        <h3>{{.Title}}</h3>
        <p>{{.UseWhen}}</p>
        <strong>Borrow</strong>
        <ul>{{range .Borrow}}<li>{{.}}</li>{{end}}</ul>
        <p class="gallery-links">{{range .Links}}<a href="{{.Href}}">{{.Label}}</a>{{end}}</p>
      </article>{{end}}
    </div>{{else}}<p>No samples yet.</p>{{end}}
  </section>
  <section class="band" id="report-kinds" aria-labelledby="report-kinds-title">
    <h2 id="report-kinds-title">Report kinds</h2>
    <ul class="gallery-kinds">
      {{range .Kinds}}<li><strong>{{.Label}}</strong>
        {{if .Latest.Href}}<a href="{{.Latest.Href}}">{{.Latest.Label}}</a>{{else}}<span>Create one with <code>{{.Command}}</code></span>{{end}}
      </li>{{end}}
    </ul>
  </section>
  <section class="band" id="presets" aria-labelledby="presets-title">
    <h2 id="presets-title">Second-opinion presets</h2>
    {{if .Presets}}<ul class="gallery-links gallery-presets">
      {{range .Presets}}<li><a href="{{.Href}}">{{.Label}}</a></li>{{end}}
    </ul>{{else}}<p>No shipped presets found.</p>{{end}}
  </section>
  <section class="band" id="add-sample" aria-labelledby="add-sample-title">
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
