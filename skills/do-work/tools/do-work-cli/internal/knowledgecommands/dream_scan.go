package knowledgecommands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var dreamFindingCodes = []string{
	"DREAM-MISSING-FROM-INDEX", "DREAM-DANGLING-INDEX", "DREAM-BROKEN-WIKI-LINK", "DREAM-ORPHAN-PAGE", "DREAM-STALE-PAGE", "DREAM-RELATIVE-DATE", "DREAM-LIKELY-DUPLICATE",
}

var markdownLinkPattern = regexp.MustCompile(`\]\(([^)]+\.md)\)`)
var relativeDatePattern = regexp.MustCompile(`(?i)\b(yesterday|today|tomorrow|tonight|last\s+(?:night|week|month|year|monday|tuesday|wednesday|thursday|friday|saturday|sunday)|next\s+(?:week|month|year|monday|tuesday|wednesday|thursday|friday|saturday|sunday)|this\s+(?:week|month|year|morning|afternoon|evening)|a\s+(?:few\s+)?(?:days?|weeks?|months?)\s+ago|recently|just\s+now|earlier\s+today|the\s+other\s+day)\b`)

type dreamPage struct {
	path, stem, title, body string
	frontmatter             map[string]string
	links                   []string
}

func handleDreamScan(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseBKBOptions(arguments, false)
	if err != nil {
		return usageResult(CommandDreamScan, err)
	}
	target, err := ensureSafeTarget(executionContext.RepositoryRoot, options.target)
	if err != nil {
		return scanFailure(CommandDreamScan, options.target, err)
	}
	indexPath, err := findDreamIndex(target)
	if err != nil {
		return scanFailure(CommandDreamScan, options.target, err)
	}
	wikiRoot, err := findDreamWiki(target, indexPath)
	if err != nil {
		return scanFailure(CommandDreamScan, options.target, err)
	}
	physicalRoot, err := ensureSafeTarget(executionContext.RepositoryRoot, ".")
	if err != nil {
		return scanFailure(CommandDreamScan, options.target, err)
	}
	pages, err := loadDreamPages(wikiRoot, indexPath, physicalRoot)
	if err != nil {
		return scanFailure(CommandDreamScan, options.target, err)
	}
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		return scanFailure(CommandDreamScan, options.target, err)
	}
	findings := scanDreamWorklist(options.target, evidencePath(physicalRoot, indexPath), pages, string(indexBytes))
	sortFindings(findings)
	if len(findings) == 0 {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{scanFinding(CommandDreamScan, options.target, "DREAM-SCAN-CLEAN", resultmodel.SeverityInfo, nil, fmt.Sprintf("all seven deterministic scans passed for %d pages", len(pages)), resultmodel.FixabilityManual, "semantic consolidation remains action-owned")}}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: findings}
}

func scanDreamWorklist(target, indexEvidence string, pages []dreamPage, indexText string) []resultmodel.CommandFinding {
	findings := []resultmodel.CommandFinding{}
	pageByStem, inbound := map[string]dreamPage{}, map[string]int{}
	for _, page := range pages {
		pageByStem[page.stem] = page
		inbound[page.stem] = 0
	}
	indexStems := map[string]bool{}
	for _, link := range normalizedLinks(indexText) {
		indexStems[link] = true
	}
	for _, match := range markdownLinkPattern.FindAllStringSubmatch(indexText, -1) {
		indexStems[normalizeLink(match[1])] = true
	}
	for _, page := range pages {
		if !indexStems[page.stem] {
			findings = append(findings, dreamFinding(target, "DREAM-MISSING-FROM-INDEX", []string{page.path}, "add to index: "+page.stem, "index placement requires action judgment"))
		}
		for _, link := range page.links {
			if link == page.stem {
				continue
			}
			if _, exists := pageByStem[link]; exists {
				inbound[link]++
			} else {
				findings = append(findings, dreamFinding(target, "DREAM-BROKEN-WIKI-LINK", []string{page.path}, fmt.Sprintf("%s -> [[%s]]", page.stem, link), "the action chooses rename, stub, or drop"))
			}
		}
		if date, ok := preferredDreamDate(page.frontmatter); ok {
			days := int(nowUTC().Sub(date).Hours() / 24)
			if days > 90 {
				findings = append(findings, dreamFinding(target, "DREAM-STALE-PAGE", []string{page.path}, fmt.Sprintf("%s updated %d days ago", page.stem, days), "staleness is evidence; truth and edits require judgment"))
			}
		}
		matches := relativeDatePattern.FindAllString(page.body, -1)
		if len(matches) > 0 {
			unique := map[string]bool{}
			for _, match := range matches {
				unique[strings.ToLower(match)] = true
			}
			values := sortedMapKeys(unique)
			findings = append(findings, dreamFinding(target, "DREAM-RELATIVE-DATE", []string{page.path}, fmt.Sprintf("%s — %s", page.stem, strings.Join(values, ", ")), "the action must establish or deliberately remove false date precision"))
		}
	}
	for stem := range indexStems {
		if _, exists := pageByStem[stem]; !exists {
			findings = append(findings, dreamFinding(target, "DREAM-DANGLING-INDEX", []string{indexEvidence}, "remove from index (dangling): "+stem, "index editing remains action-owned"))
		}
	}
	for stem, count := range inbound {
		if count == 0 {
			findings = append(findings, dreamFinding(target, "DREAM-ORPHAN-PAGE", []string{pageByStem[stem].path}, stem+" has no inbound links", "the action decides whether the orphan is a valid top-level page"))
		}
	}
	for first := 0; first < len(pages); first++ {
		for second := first + 1; second < len(pages); second++ {
			if likelyDuplicateTitle(pages[first].title, pages[second].title) {
				findings = append(findings, dreamFinding(target, "DREAM-LIKELY-DUPLICATE", []string{pages[first].path, pages[second].path}, fmt.Sprintf("%s ~ %s (similar title)", pages[first].stem, pages[second].stem), "the action reads both bodies and chooses merge or deliberate preservation"))
			}
		}
	}
	return findings
}

func dreamFinding(target, code string, paths []string, evidence, stopReason string) resultmodel.CommandFinding {
	return scanFinding(CommandDreamScan, target, code, resultmodel.SeverityWarning, paths, evidence, resultmodel.FixabilityManual, stopReason)
}

func findDreamIndex(target string) (string, error) {
	for _, name := range []string{"MEMORY.md", "_master_index.md", "index.md"} {
		path := filepath.Join(target, name)
		if info, err := os.Lstat(path); err == nil && info.Mode().IsRegular() {
			return path, nil
		}
	}
	return "", fmt.Errorf("target directory has no index file (MEMORY.md / _master_index.md / index.md)")
}

func findDreamWiki(target, indexPath string) (string, error) {
	for _, candidate := range []string{filepath.Join(target, "wiki"), filepath.Join(target, "pages"), target} {
		entries, err := os.ReadDir(candidate)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.Type()&os.ModeSymlink != 0 {
				return "", fmt.Errorf("scan path %q is a symlink", filepath.Join(candidate, entry.Name()))
			}
			if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".md") && filepath.Join(candidate, entry.Name()) != indexPath {
				return candidate, nil
			}
		}
	}
	return "", fmt.Errorf("target has no wiki pages")
}

func loadDreamPages(wikiRoot, indexPath, repositoryRoot string) ([]dreamPage, error) {
	pages := []dreamPage{}
	entries, err := os.ReadDir(wikiRoot)
	if err != nil {
		return nil, err
	}
	stems := map[string]string{}
	for _, entry := range entries {
		path := filepath.Join(wikiRoot, entry.Name())
		if entry.Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("scan path %q is a symlink", path)
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") || path == indexPath {
			continue
		}
		base := filepath.Base(path)
		if base == "MEMORY.md" || base == "_master_index.md" || base == "index.md" || base == "log.md" {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		frontmatter, body, err := parseFrontmatter(string(data))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		stem := strings.TrimSuffix(base, filepath.Ext(base))
		key := strings.ToLower(normalizeLink(stem))
		if previous, exists := stems[key]; exists {
			return nil, fmt.Errorf("duplicate page stem %q in %s and %s", stem, previous, path)
		}
		stems[key] = path
		page := dreamPage{path: evidencePath(repositoryRoot, path), stem: stem, title: frontmatter["name"], body: body, frontmatter: frontmatter, links: normalizedLinks(body)}
		if page.title == "" {
			page.title = frontmatter["title"]
		}
		if page.title == "" {
			page.title = firstHeading(body, page.stem)
		}
		pages = append(pages, page)
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].path < pages[j].path })
	return pages, nil
}

func evidencePath(repositoryRoot, path string) string {
	relative, err := filepath.Rel(repositoryRoot, path)
	if err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(relative)
	}
	return filepath.ToSlash(path)
}

func preferredDreamDate(frontmatter map[string]string) (time.Time, bool) {
	for _, field := range []string{"last_updated", "updated", "created"} {
		value := frontmatter[field]
		if value == "" {
			continue
		}
		for _, layout := range []string{"2006-01-02", "2006/01/02"} {
			if parsed, err := time.ParseInLocation(layout, value, time.UTC); err == nil {
				return parsed, true
			}
		}
	}
	return time.Time{}, false
}

func likelyDuplicateTitle(first, second string) bool {
	first = normalizeTitle(first)
	second = normalizeTitle(second)
	if first == "" || second == "" || first == second {
		return first != ""
	}
	if strings.Contains(first, second) || strings.Contains(second, first) {
		return true
	}
	firstTokens, secondTokens := tokenSet(first), tokenSet(second)
	shorter := len(firstTokens)
	if len(secondTokens) < shorter {
		shorter = len(secondTokens)
	}
	if shorter > 0 {
		shared := 0
		for token := range firstTokens {
			if secondTokens[token] {
				shared++
			}
		}
		if float64(shared)/float64(shorter) >= 0.8 {
			return true
		}
	}
	return normalizeTitleSuffix(first) == normalizeTitleSuffix(second)
}
func normalizeTitle(value string) string {
	value = strings.ToLower(strings.TrimSpace(strings.TrimRightFunc(value, unicode.IsPunct)))
	return strings.Join(strings.Fields(value), " ")
}
func normalizeTitleSuffix(value string) string {
	tokens := strings.FieldsFunc(value, func(r rune) bool { return !(unicode.IsLetter(r) || unicode.IsDigit(r)) })
	if len(tokens) == 0 {
		return ""
	}
	last := tokens[len(tokens)-1]
	if regexp.MustCompile(`^v\d+$`).MatchString(last) {
		tokens = tokens[:len(tokens)-1]
	} else if strings.HasSuffix(last, "s") && len(last) > 1 {
		tokens[len(tokens)-1] = strings.TrimSuffix(last, "s")
	}
	return strings.Join(tokens, " ")
}
func tokenSet(value string) map[string]bool {
	tokens := map[string]bool{}
	for _, token := range strings.FieldsFunc(value, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		tokens[token] = true
	}
	return tokens
}
func sortedMapKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
