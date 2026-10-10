package toolboxcommands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// The catalog is derived: every run rebuilds it from the bundles, so it never needs hand edits and
// never breaks the rule that published bundles are immutable.
const (
	reportIndexGenerator = "do-work-cli ai-report-index"
	reportIndexMarker    = `<meta name="generator" content="` + reportIndexGenerator + `">`
	reportsFolder        = "ai-reports"
	reportCatalogFile    = reportsFolder + "/catalog.json"
	reportIndexPageFile  = reportsFolder + "/index.html"
)

// Entry files in pick order; after these come the first .html, then the first .md, by name.
var reportEntryNames = []string{"index.html", "index.md", "README.md", "prompt.md", "report.md"}

// Kinds recognised in a folder name when the entry has no ai-report-kind meta. Illustrative, not a
// closed list: the meta always wins, and a kind missing here only falls back to "unknown".
var reportKnownKinds = []string{"proposal", "root-cause", "architecture-report"}

var (
	reportDatePattern     = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})(?:_(\d{4}|\d{6}))?(?:\D|$)`)
	reportLinkedIDPattern = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9])((?:UR|REQ)-\d+)`)
	reportTopicIDPattern  = regexp.MustCompile(`(?i)^(UR|REQ)-(\d+)$`)
	reportTitleTag        = regexp.MustCompile(`(?is)<title(?:\s[^>]*)?>(.*?)</title>`)
	reportHeadingTag      = regexp.MustCompile(`(?is)<h[1-6](?:\s[^>]*)?>(.*?)</h[1-6]>`)
	reportMarkdownHeading = regexp.MustCompile(`(?m)^#{1,6}[ \t]+(.+)$`)
	reportMetaTag         = regexp.MustCompile(`(?is)<meta\b[^>]*>`)
	reportMetaAttribute   = regexp.MustCompile(`(?is)\b(name|content)\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	reportMarkupTag       = regexp.MustCompile(`<[^>]*>`)
)

type reportCatalog struct {
	Generator string         `json:"generator"`
	Bundles   []reportBundle `json:"bundles"`
}

// reportBundle field names are read by ai-report revise (path, date, supersedes, superseded_by).
type reportBundle struct {
	Path         string   `json:"path"`
	Entry        string   `json:"entry"`
	Title        string   `json:"title"`
	Date         string   `json:"date"`
	Kind         string   `json:"kind"`
	LinkedIDs    []string `json:"linked_ids"`
	Verdict      string   `json:"verdict"`
	Supersedes   string   `json:"supersedes"`
	SupersededBy string   `json:"superseded_by"`
}

func handleAIReportIndex(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	if len(args) == 0 {
		return reportIndexWrite(ctx.RepositoryRoot)
	}
	if len(args) >= 2 && args[0] == "--find" {
		if topic := strings.TrimSpace(strings.Join(args[1:], " ")); topic != "" {
			return reportIndexFind(ctx.RepositoryRoot, topic)
		}
	}
	return usageResult(CommandAIReportIndex, "Usage: ai-report-index [--find <topic...>]")
}

// reportIndexWalk treats every directory directly under ai-reports/ whose name does not start with
// "." as a bundle, whatever its naming style. Loose files and symbolic links are never bundles.
func reportIndexWalk(repositoryRoot string) (bundles []reportBundle, present bool, err error) {
	entries, err := os.ReadDir(filepath.Join(repositoryRoot, reportsFolder))
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, true, err
	}
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			bundles = append(bundles, readReportBundle(repositoryRoot, entry.Name()))
		}
	}
	sort.SliceStable(bundles, func(i, j int) bool {
		left, right := bundles[i], bundles[j]
		if left.Date != right.Date {
			return right.Date == "" || (left.Date != "" && left.Date > right.Date)
		}
		return left.Path < right.Path
	})
	byFolder := map[string]int{}
	for index, bundle := range bundles {
		byFolder[path.Base(bundle.Path)] = index
	}
	for index := range bundles {
		target, found := byFolder[path.Base(strings.TrimRight(filepath.ToSlash(bundles[index].Supersedes), "/"))]
		if bundles[index].Supersedes == "" || !found || target == index {
			continue
		}
		bundles[index].Supersedes = bundles[target].Path
		// Bundles are sorted newest first, so the first successor seen is the newest one.
		if bundles[target].SupersededBy == "" {
			bundles[target].SupersededBy = bundles[index].Path
		}
	}
	return bundles, true, nil
}

func readReportBundle(repositoryRoot, folder string) reportBundle {
	bundle := reportBundle{Path: reportsFolder + "/" + folder, Title: folder, Kind: "unknown", LinkedIDs: []string{}}
	if match := reportDatePattern.FindStringSubmatch(folder); match != nil {
		bundle.Date = match[1]
		if len(match[2]) >= 4 {
			bundle.Date += "T" + match[2][0:2] + ":" + match[2][2:4]
		}
		if len(match[2]) == 6 {
			bundle.Date += ":" + match[2][4:6]
		}
	}
	entryName := pickReportEntry(filepath.Join(repositoryRoot, reportsFolder, folder))
	metas := map[string]string{}
	if entryName != "" {
		bundle.Entry = bundle.Path + "/" + entryName
		content, err := os.ReadFile(filepath.Join(repositoryRoot, reportsFolder, folder, entryName))
		if err == nil {
			title := ""
			if strings.EqualFold(filepath.Ext(entryName), ".html") {
				metas = readReportMetas(content)
				if match := reportTitleTag.FindSubmatch(content); match != nil {
					title = cleanReportText(string(match[1]))
				}
				if match := reportHeadingTag.FindSubmatch(content); title == "" && match != nil {
					title = cleanReportText(string(match[1]))
				}
			} else if match := reportMarkdownHeading.FindSubmatch(content); match != nil {
				title = cleanReportText(strings.TrimRight(string(match[1]), " \t#"))
			}
			if title != "" {
				bundle.Title = title
			}
		}
	}
	bundle.Kind = reportKind(folder, metas["ai-report-kind"])
	bundle.Verdict = metas["ai-report-verdict"]
	bundle.Supersedes = metas["ai-report-supersedes"]
	seen := map[string]bool{}
	for _, match := range reportLinkedIDPattern.FindAllStringSubmatch(folder+" "+bundle.Title, -1) {
		id := strings.ToUpper(match[1])
		if !seen[id] {
			seen[id] = true
			bundle.LinkedIDs = append(bundle.LinkedIDs, id)
		}
	}
	return bundle
}

func pickReportEntry(bundleDirectory string) string {
	entries, err := os.ReadDir(bundleDirectory)
	if err != nil {
		return ""
	}
	regular := map[string]bool{}
	firstHTML, firstMarkdown := "", ""
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		name := entry.Name()
		regular[name] = true
		extension := strings.ToLower(filepath.Ext(name))
		if extension == ".html" && firstHTML == "" {
			firstHTML = name
		}
		if extension == ".md" && firstMarkdown == "" {
			firstMarkdown = name
		}
	}
	for _, name := range reportEntryNames {
		if regular[name] {
			return name
		}
	}
	if firstHTML != "" {
		return firstHTML
	}
	return firstMarkdown
}

// readReportMetas reads <meta name=... content=...> tags in either attribute order and quote style.
func readReportMetas(content []byte) map[string]string {
	metas := map[string]string{}
	for _, tag := range reportMetaTag.FindAll(content, -1) {
		attributes := map[string]string{}
		for _, attribute := range reportMetaAttribute.FindAllSubmatch(tag, -1) {
			attributes[strings.ToLower(string(attribute[1]))] = string(attribute[2]) + string(attribute[3])
		}
		name := strings.ToLower(strings.TrimSpace(attributes["name"]))
		if _, seen := metas[name]; name != "" && !seen {
			metas[name] = strings.TrimSpace(html.UnescapeString(attributes["content"]))
		}
	}
	return metas
}

func cleanReportText(markup string) string {
	return strings.Join(strings.Fields(html.UnescapeString(reportMarkupTag.ReplaceAllString(markup, " "))), " ")
}

// reportKind prefers the meta; otherwise a known kind must appear as a whole run of kebab words in
// the folder name (date digits can never form a kind word).
func reportKind(folder, meta string) string {
	if meta != "" {
		return meta
	}
	words := "-" + strings.Join(strings.FieldsFunc(strings.ToLower(folder), func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	}), "-") + "-"
	for _, kind := range reportKnownKinds {
		if strings.Contains(words, "-"+kind+"-") {
			return kind
		}
	}
	return "unknown"
}

func reportIndexWrite(repositoryRoot string) resultmodel.CommandResult {
	bundles, present, err := reportIndexWalk(repositoryRoot)
	if err != nil {
		return reportIndexFailure(reportsFolder, err)
	}
	if !present {
		return exactOutputResult(reportsFolder+"/ does not exist; nothing to index\n", nil)
	}
	outputs := []struct {
		relative, marker string
		data             []byte
		existed          bool
	}{
		{relative: reportCatalogFile, marker: `"generator": "` + reportIndexGenerator + `"`, data: renderReportCatalog(bundles)},
		{relative: reportIndexPageFile, marker: reportIndexMarker, data: renderReportIndexPage(bundles)},
	}
	refused := resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused}
	for index := range outputs {
		output := &outputs[index]
		info, err := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(output.relative)))
		if os.IsNotExist(err) {
			continue
		}
		output.existed = true
		existing := []byte{}
		if err == nil && info.Mode().IsRegular() {
			existing, _ = os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(output.relative)))
		}
		if !bytes.Contains(existing, []byte(output.marker)) {
			refused.Findings = append(refused.Findings, toolboxFinding(CommandAIReportIndex, "AI-REPORT-INDEX-HAND-MADE", resultmodel.SeverityError,
				[]string{output.relative}, output.relative+" was not written by ai-report-index (no generator marker); move it aside by hand, then rerun",
				resultmodel.FixabilityManual, "a hand-made catalog file is never overwritten"))
		}
	}
	if len(refused.Findings) > 0 {
		return refused
	}
	changes := []resultmodel.RecordedChange{}
	for _, output := range outputs {
		if err := rootedPublishFile(repositoryRoot, output.relative, output.data, 0o644, output.existed); err != nil {
			return reportIndexFailure(output.relative, err)
		}
		kind := "created"
		if output.existed {
			kind = "modified"
		}
		changes = append(changes, resultmodel.RecordedChange{Path: output.relative, Kind: kind, Detail: "regenerated derived report catalog"})
	}
	return exactOutputResult(fmt.Sprintf("%s: %d bundles\n%s\n", reportCatalogFile, len(bundles), reportIndexPageFile), changes)
}

func reportIndexFailure(relative string, err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandAIReportIndex, "AI-REPORT-INDEX-FAILED", resultmodel.SeverityError, []string{relative}, err.Error(),
			resultmodel.FixabilityManual, "the report catalog could not be read or written"),
	}}
}

func renderReportCatalog(bundles []reportBundle) []byte {
	if bundles == nil {
		bundles = []reportBundle{}
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(reportCatalog{Generator: reportIndexGenerator, Bundles: bundles})
	return buffer.Bytes()
}

// reportPageLink turns a catalog path into an href relative to ai-reports/index.html.
func reportPageLink(bundle reportBundle) string {
	target := bundle.Entry
	if target == "" {
		target = bundle.Path + "/"
	}
	relative := strings.TrimPrefix(target, reportsFolder+"/")
	return html.EscapeString((&url.URL{Path: relative}).String())
}

func renderReportIndexPage(bundles []reportBundle) []byte {
	byPath := map[string]reportBundle{}
	groups := map[string][]reportBundle{}
	kinds := []string{}
	for _, bundle := range bundles {
		byPath[bundle.Path] = bundle
		if groups[bundle.Kind] == nil {
			kinds = append(kinds, bundle.Kind)
		}
		groups[bundle.Kind] = append(groups[bundle.Kind], bundle)
	}
	sort.Slice(kinds, func(i, j int) bool {
		if (kinds[i] == "unknown") != (kinds[j] == "unknown") {
			return kinds[j] == "unknown"
		}
		return kinds[i] < kinds[j]
	})
	var page strings.Builder
	page.WriteString(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
` + reportIndexMarker + `
<title>AI Report Catalog</title>
<style>
:root{color-scheme:light dark;--text:#1f2328;--muted:#656d76;--line:#d0d7de;--page:#ffffff;--link:#0969da}
@media (prefers-color-scheme: dark){:root{--text:#e6edf3;--muted:#8d96a0;--line:#30363d;--page:#0d1117;--link:#4493f8}}
body{margin:0 auto;max-width:64rem;padding:1.5rem 1rem;font:15px/1.5 system-ui,sans-serif;color:var(--text);background:var(--page)}
a{color:var(--link)}
table{width:100%;border-collapse:collapse;margin-bottom:1.5rem}
td{padding:.35rem .5rem;border-top:1px solid var(--line);vertical-align:top}
td.report-date,td.report-ids{white-space:nowrap;color:var(--muted)}
tr.superseded{opacity:.5}
</style>
</head>
<body>
<h1>AI Report Catalog</h1>
<p>Generated by do-work-cli ai-report-index from the bundles under ai-reports/. Do not edit by hand; rerun the command.</p>
`)
	for _, kind := range kinds {
		fmt.Fprintf(&page, "<h2>%s</h2>\n<table>\n", html.EscapeString(kind))
		for _, bundle := range groups[kind] {
			rowClass, successor := "", ""
			if bundle.SupersededBy != "" {
				rowClass = ` class="superseded"`
				successor = fmt.Sprintf(` <span>superseded by <a href="%s">%s</a></span>`,
					reportPageLink(byPath[bundle.SupersededBy]), html.EscapeString(bundle.SupersededBy))
			}
			fmt.Fprintf(&page, "<tr%s><td class=\"report-date\">%s</td><td><a href=\"%s\">%s</a>%s</td><td class=\"report-ids\">%s</td></tr>\n",
				rowClass, html.EscapeString(bundle.Date), reportPageLink(bundle), html.EscapeString(bundle.Title), successor,
				html.EscapeString(strings.Join(bundle.LinkedIDs, " ")))
		}
		page.WriteString("</table>\n")
	}
	page.WriteString("</body>\n</html>\n")
	return []byte(page.String())
}

// reportIndexFind walks the folder in memory on every call; it never reads or writes catalog.json.
func reportIndexFind(repositoryRoot, topic string) resultmodel.CommandResult {
	bundles, _, err := reportIndexWalk(repositoryRoot)
	if err != nil {
		return reportIndexFailure(reportsFolder, err)
	}
	needle := strings.ToLower(topic)
	topicID := normalizeReportID(topic)
	var output strings.Builder
	for _, bundle := range bundles {
		// Match the folder name, not the full path: every path starts with "ai-reports/", so topics
		// such as "ai" or "reports" would otherwise match every bundle.
		matched := false
		for _, field := range append([]string{path.Base(bundle.Path), bundle.Title, bundle.Kind, bundle.Verdict}, bundle.LinkedIDs...) {
			matched = matched || strings.Contains(strings.ToLower(field), needle)
		}
		for _, id := range bundle.LinkedIDs {
			matched = matched || (topicID != "" && normalizeReportID(id) == topicID)
		}
		if !matched {
			continue
		}
		fmt.Fprintf(&output, "%s\t%s", bundle.Path, bundle.Title)
		if bundle.SupersededBy != "" {
			fmt.Fprintf(&output, " (superseded by %s)", bundle.SupersededBy)
		}
		output.WriteString("\n")
	}
	if output.Len() == 0 {
		return exactOutputResult("no report matches "+topic+"\n", nil)
	}
	return exactOutputResult(output.String(), nil)
}

// normalizeReportID returns "REQ-412" for "req-0412", or "" when the text is not a UR/REQ id.
func normalizeReportID(text string) string {
	match := reportTopicIDPattern.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return ""
	}
	number := strings.TrimLeft(match[2], "0")
	if number == "" {
		number = "0"
	}
	return strings.ToUpper(match[1]) + "-" + number
}
