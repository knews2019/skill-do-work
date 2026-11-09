package knowledgecommands

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var wikiLinkPattern = regexp.MustCompile(`\[\[([^\]]+?)\]\]`)
var datePattern = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`)

type knowledgePage struct {
	path             string
	stem             string
	title            string
	body             string
	frontmatter      map[string]string
	frontmatterError string
	links            []string
	relations        []pageRelation
}
type pageRelation struct{ page, relation string }

func handleBKBStatus(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseBKBOptions(arguments, false)
	if err != nil {
		return usageResult(CommandBKBStatus, err)
	}
	location := resolveKnowledgeBaseLocation(executionContext.RepositoryRoot, options.target)
	target, relative, err := locateKnowledgeBase(executionContext.RepositoryRoot, options.target)
	if err != nil {
		finding := knowledgeFinding(CommandBKBStatus, "BKB-NOT-FOUND", resultmodel.SeverityWarning, []string{relative}, err.Error(), resultmodel.FixabilityAutomatic, "initialize a knowledge base or supply its repository-relative path")
		finding.NextArgv = []string{"do-work-cli", CommandBKBInit, "--kb", relative}
		finding.NextJustRecipe = CommandBKBInit + " " + quoteRecipeArgument(relative)
		finding.VerificationArgv = []string{"do-work-cli", "--format", "json", CommandBKBStatus, "--kb", relative}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{finding}}
	}
	pages, scanError := loadKnowledgePages(target, relative)
	if scanError != nil {
		return scanFailure(CommandBKBStatus, relative, scanError)
	}
	topicIndexes, _ := filepath.Glob(filepath.Join(target, "wiki", "topics", "_index_*.md"))
	masterPath := filepath.Join(target, "wiki", "_master_index.md")
	masterData, masterError := os.ReadFile(masterPath)
	articleCounts := inspectNamedCounts(string(masterData), "Total articles:")
	topicCounts := inspectNamedCounts(string(masterData), "Topic clusters:")
	masterArticles, masterTopics := articleCounts.value, topicCounts.value
	inboxCount := countRegularEntries(filepath.Join(target, "raw", "inbox"))
	queueReady := countLinesContaining(filepath.Join(target, "raw", "_inbox_queue.md"), "ready")
	latestActivity := latestNamedDate(filepath.Join(target, "wiki", "daily"))
	logBytes, _ := os.ReadFile(filepath.Join(target, "wiki", "log.md"))
	lastLint := latestLogDate(string(logBytes), "lint |")
	lastDefrag := latestLogDate(string(logBytes), "defrag |")
	lastGarden := latestLogDate(string(logBytes), "garden |")
	customAgents := countCustomAgents(filepath.Join(target, "agents"))
	evidence := fmt.Sprintf("location=%s articles=%d topic_clusters=%d disk_articles=%d disk_topic_clusters=%d inbox=%d ready=%d last_activity=%s last_lint=%s last_defrag=%s last_garden=%s agents=8 built-in+%d custom", relative, masterArticles, masterTopics, len(pages), len(topicIndexes), inboxCount, queueReady, displayDate(latestActivity), displayDate(lastLint), displayDate(lastDefrag), displayDate(lastGarden), customAgents)
	findings := []resultmodel.CommandFinding{scanFinding(CommandBKBStatus, relative, "BKB-STATUS", resultmodel.SeverityInfo, nil, evidence, resultmodel.FixabilityManual, "")}
	masterRelative := filepath.ToSlash(filepath.Join(relative, "wiki/_master_index.md"))
	if masterError != nil {
		findings = append(findings, scanFinding(CommandBKBStatus, relative, "BKB-STATUS-MASTER-MISSING", resultmodel.SeverityError, []string{masterRelative}, masterError.Error(), resultmodel.FixabilityManual, "restore the canonical master index before trusting status counts"))
	} else {
		findings = append(findings, namedCountFindings(relative, masterRelative, "ARTICLE", "Total articles", articleCounts, len(pages))...)
		findings = append(findings, namedCountFindings(relative, masterRelative, "TOPIC", "Topic clusters", topicCounts, len(topicIndexes))...)
	}
	for _, maintenance := range []struct{ name, date string }{{"defrag", lastDefrag}, {"garden", lastGarden}} {
		days := daysSince(maintenance.date)
		if maintenance.date == "" || days >= 14 {
			observed := maintenance.name + " has never run"
			if maintenance.date != "" {
				observed = fmt.Sprintf("%s last ran %d days ago on %s", maintenance.name, days, maintenance.date)
			}
			findings = append(findings, scanFinding(CommandBKBStatus, relative, "BKB-"+strings.ToUpper(maintenance.name)+"-OVERDUE", resultmodel.SeverityWarning, []string{filepath.ToSlash(filepath.Join(relative, "wiki/log.md"))}, observed, resultmodel.FixabilityManual, "maintenance scheduling and edits remain action-owned"))
		}
	}
	findings = append(findings, location.notices(CommandBKBStatus)...)
	sortFindings(findings)
	outcome := resultmodel.OutcomeSuccess
	for _, finding := range findings {
		if finding.Severity != resultmodel.SeverityInfo {
			outcome = resultmodel.OutcomeFindings
		}
	}
	return resultmodel.CommandResult{Outcome: outcome, Findings: findings}
}

func handleBKBLint(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseBKBOptions(arguments, false)
	if err != nil {
		return usageResult(CommandBKBLintStructure, err)
	}
	location := resolveKnowledgeBaseLocation(executionContext.RepositoryRoot, options.target)
	target, relative, err := locateKnowledgeBase(executionContext.RepositoryRoot, options.target)
	if err != nil {
		return scanFailure(CommandBKBLintStructure, relative, err)
	}
	pages, scanError := loadKnowledgePages(target, relative)
	if scanError != nil {
		return scanFailure(CommandBKBLintStructure, relative, scanError)
	}
	findings := lintKnowledgeBase(target, relative, pages)
	if len(findings) == 0 {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-STRUCTURE-CLEAN", resultmodel.SeverityInfo, nil, "all deterministic structural checks passed", resultmodel.FixabilityManual, "semantic lint remains action-owned"))
	}
	findings = append(findings, location.notices(CommandBKBLintStructure)...)
	sortFindings(findings)
	outcome := resultmodel.OutcomeSuccess
	for _, finding := range findings {
		if finding.Severity != resultmodel.SeverityInfo {
			outcome = resultmodel.OutcomeFindings
		}
	}
	return resultmodel.CommandResult{Outcome: outcome, Findings: findings}
}

func lintKnowledgeBase(target, relative string, pages []knowledgePage) []resultmodel.CommandFinding {
	findings := []resultmodel.CommandFinding{}
	pageByStem, inbound := map[string]knowledgePage{}, map[string]int{}
	for _, page := range pages {
		pageByStem[page.stem] = page
		inbound[page.stem] = 0
	}
	for _, page := range pages {
		for _, link := range page.links {
			if link == page.stem {
				continue
			}
			if _, exists := pageByStem[link]; exists {
				inbound[link]++
			} else {
				findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-BROKEN-LINK", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("%s links to missing [[%s]]", page.stem, link), resultmodel.FixabilityManual, "the action must choose rename, stub, or removal"))
			}
		}
	}
	for stem, count := range inbound {
		if count == 0 {
			page := pageByStem[stem]
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-ORPHAN-PAGE", resultmodel.SeverityWarning, []string{page.path}, stem+" has no inbound wiki links", resultmodel.FixabilityManual, "the action decides whether a top-level orphan is intentional"))
		}
	}

	indexMembership := map[string][]string{}
	topicIndexes, _ := filepath.Glob(filepath.Join(target, "wiki", "topics", "_index_*.md"))
	sort.Strings(topicIndexes)
	masterBytes, masterError := os.ReadFile(filepath.Join(target, "wiki", "_master_index.md"))
	masterText := string(masterBytes)
	if masterError != nil {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-MASTER-INDEX-MISSING", resultmodel.SeverityError, []string{filepath.ToSlash(filepath.Join(relative, "wiki/_master_index.md"))}, masterError.Error(), resultmodel.FixabilityManual, "the canonical master index is required"))
	}
	for _, indexPath := range topicIndexes {
		data, readError := os.ReadFile(indexPath)
		if readError != nil {
			continue
		}
		indexRelative := relativePath(target, relative, indexPath)
		indexStem := strings.TrimSuffix(filepath.Base(indexPath), ".md")
		if !containsNormalizedLink(masterText, indexStem) {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-TOPIC-NOT-IN-MASTER", resultmodel.SeverityWarning, []string{indexRelative}, indexStem+" is absent from the master index", resultmodel.FixabilityManual, "cluster placement is action-owned"))
		}
		links := normalizedLinks(string(data))
		if declared, ok := extractNamedCount(string(data), "Total articles:"); ok && declared != len(links) {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-TOPIC-ARTICLE-COUNT", resultmodel.SeverityWarning, []string{indexRelative}, fmt.Sprintf("%s declares %d articles but links %d", indexStem, declared, len(links)), resultmodel.FixabilityManual, "the action owns index rebuilding"))
		}
		if len(links) > 40 {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-TOPIC-SPLIT-THRESHOLD", resultmodel.SeverityWarning, []string{indexRelative}, fmt.Sprintf("%s contains %d article links; maximum is 40", indexStem, len(links)), resultmodel.FixabilityManual, "cluster split design requires semantic judgment"))
		}
		for _, stem := range links {
			if _, exists := pageByStem[stem]; exists {
				indexMembership[stem] = append(indexMembership[stem], indexRelative)
			} else {
				findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-TOPIC-DANGLING-ENTRY", resultmodel.SeverityWarning, []string{indexRelative}, fmt.Sprintf("%s links to missing article [[%s]]", indexStem, stem), resultmodel.FixabilityManual, "the action chooses rename, stub, or removal from the topic index"))
			}
		}
	}
	for _, page := range pages {
		memberships := indexMembership[page.stem]
		if len(memberships) != 1 {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-INDEX-MEMBERSHIP", resultmodel.SeverityWarning, append([]string{page.path}, memberships...), fmt.Sprintf("%s appears in %d topic indexes; required exactly one", page.stem, len(memberships)), resultmodel.FixabilityManual, "the action chooses the correct topic cluster"))
		}
		findings = append(findings, frontmatterFindings(relative, page, pageByStem)...)
	}
	if masterError == nil {
		if declared, ok := extractNamedCount(masterText, "Total articles:"); ok && declared != len(pages) {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-MASTER-ARTICLE-COUNT", resultmodel.SeverityWarning, []string{filepath.ToSlash(filepath.Join(relative, "wiki/_master_index.md"))}, fmt.Sprintf("master declares %d articles but disk has %d", declared, len(pages)), resultmodel.FixabilityManual, "the action owns index rebuilding"))
		}
		if declared, ok := extractNamedCount(masterText, "Topic clusters:"); ok && declared != len(topicIndexes) {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-MASTER-TOPIC-COUNT", resultmodel.SeverityWarning, []string{filepath.ToSlash(filepath.Join(relative, "wiki/_master_index.md"))}, fmt.Sprintf("master declares %d clusters but disk has %d", declared, len(topicIndexes)), resultmodel.FixabilityManual, "the action owns index rebuilding"))
		}
	}
	logBytes, _ := os.ReadFile(filepath.Join(target, "wiki", "log.md"))
	for _, line := range strings.Split(string(logBytes), "\n") {
		if strings.Contains(strings.ToLower(line), "ingest") {
			if date := datePattern.FindString(line); date != "" {
				dailyPath := filepath.Join(target, "wiki", "daily", date+".md")
				if _, err := os.Stat(dailyPath); os.IsNotExist(err) {
					findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-DAILY-LOG-MISSING", resultmodel.SeverityWarning, []string{filepath.ToSlash(filepath.Join(relative, "wiki/daily", date+".md"))}, "ingestion activity on "+date+" has no daily log", resultmodel.FixabilityManual, "the action reconstructs the semantic daily report"))
				}
			}
		}
	}
	findings = append(findings, agentStalenessFindings(target, relative)...)
	return findings
}

func frontmatterFindings(relative string, page knowledgePage, pageByStem map[string]knowledgePage) []resultmodel.CommandFinding {
	findings := []resultmodel.CommandFinding{}
	required := []string{"title", "type", "sources", "related", "created", "updated", "confidence"}
	for _, field := range required {
		if _, exists := page.frontmatter[field]; !exists {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-FRONTMATTER-MISSING", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("%s is missing required field %s", page.stem, field), resultmodel.FixabilityManual, "the action must infer a truthful value"))
		}
	}
	if page.frontmatter["topic_cluster"] == "" && page.frontmatter["topic"] == "" && page.frontmatter["topic_category"] == "" {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-FRONTMATTER-MISSING", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("%s is missing required field topic_cluster (or accepted read alias topic/topic_category)", page.stem), resultmodel.FixabilityManual, "the action must infer a truthful value"))
	}
	allowedTypes := map[string]bool{"concept": true, "entity": true, "source-summary": true, "comparison": true, "daily-log": true, "monthly-rollup": true}
	if value := normalizeTypeEnum(page.frontmatter["type"]); value != "" && !allowedTypes[value] {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-FRONTMATTER-ENUM", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("type=%q is not accepted", value), resultmodel.FixabilityManual, "the action chooses the correct normalized enum"))
	}
	if value := page.frontmatter["confidence"]; value != "" && value != "high" && value != "medium" && value != "low" {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-FRONTMATTER-ENUM", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("confidence=%q is not accepted", value), resultmodel.FixabilityManual, "the action chooses the evidence-backed confidence"))
	}
	if len(page.relations) > 8 {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-RELATIONSHIP-DENSITY", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("%s has %d relationships; maximum is 8", page.stem, len(page.relations)), resultmodel.FixabilityManual, "the action chooses the weakest relationship to remove"))
	}
	allowedRelations := map[string]bool{"extends": true, "contradicts": true, "evidence-for": true, "complements": true, "supersedes": true, "depends-on": true}
	for _, relation := range page.relations {
		normalizedRelation := normalizeRelationshipEnum(relation.relation)
		if !allowedRelations[normalizedRelation] {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-RELATIONSHIP-ENUM", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("%s uses invalid rel=%q", page.stem, relation.relation), resultmodel.FixabilityManual, "the action chooses the semantic relationship"))
		}
		targetStem := normalizeLink(relation.page)
		target, exists := pageByStem[targetStem]
		if !exists {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-RELATIONSHIP-TARGET", resultmodel.SeverityWarning, []string{page.path}, fmt.Sprintf("%s relates to missing page %q", page.stem, relation.page), resultmodel.FixabilityManual, "the action chooses repair or removal"))
			continue
		}
		reciprocal := false
		for _, candidate := range target.relations {
			if normalizeLink(candidate.page) == page.stem {
				reciprocal = true
				break
			}
		}
		if !reciprocal {
			findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-RELATIONSHIP-RECIPROCITY", resultmodel.SeverityWarning, []string{page.path, target.path}, fmt.Sprintf("%s relates to %s without a reciprocal typed relationship", page.stem, targetStem), resultmodel.FixabilityManual, "the action chooses the truthful reciprocal relationship or removes the one-way edge"))
		}
	}
	if page.frontmatterError != "" {
		findings = append(findings, scanFinding(CommandBKBLintStructure, relative, "BKB-FRONTMATTER-MALFORMED", resultmodel.SeverityWarning, []string{page.path}, page.frontmatterError, resultmodel.FixabilityManual, "repair the malformed YAML scalar without changing unrelated frontmatter"))
	}
	return findings
}

func loadKnowledgePages(target, relativeTarget string) ([]knowledgePage, error) {
	wikiRoot := filepath.Join(target, "wiki")
	pages := []knowledgePage{}
	err := filepath.WalkDir(wikiRoot, func(path string, entry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("scan path %q is a symlink", path)
		}
		if entry.IsDir() || strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}
		relative, _ := filepath.Rel(wikiRoot, path)
		relative = filepath.ToSlash(relative)
		base := filepath.Base(path)
		if base == "_master_index.md" || base == "log.md" || base == "overview.md" || base == "agent.md" || strings.HasPrefix(base, "_index_") || strings.HasPrefix(relative, "daily/") || strings.HasPrefix(relative, "monthly/") {
			return nil
		}
		data, readError := os.ReadFile(path)
		if readError != nil {
			return readError
		}
		frontmatter, body, parseError := parseFrontmatter(string(data))
		if frontmatter == nil {
			frontmatter = map[string]string{}
		}
		page := knowledgePage{path: filepath.ToSlash(filepath.Join(relativeTarget, "wiki", relative)), stem: strings.TrimSuffix(base, ".md"), title: frontmatter["title"], body: body, frontmatter: frontmatter, links: normalizedLinks(body), relations: parseRelations(string(data))}
		if parseError != nil {
			page.frontmatterError = parseError.Error()
		}
		if page.title == "" {
			page.title = firstHeading(body, page.stem)
		}
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].path < pages[j].path })
	return pages, nil
}

func parseFrontmatter(content string) (map[string]string, string, error) {
	fields := map[string]string{}
	content = strings.TrimPrefix(content, "\ufeff")
	content = strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.HasPrefix(content, "---\n") {
		return fields, content, nil
	}
	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return nil, "", fmt.Errorf("unterminated YAML frontmatter")
	}
	frontmatterText := content[4 : 4+end]
	var parseError error
	for lineIndex, line := range strings.Split(frontmatterText, "\n") {
		if strings.HasPrefix(line, " ") || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		fieldName := strings.TrimSpace(parts[0])
		value, valueError := parseFrontmatterScalar(strings.TrimSpace(parts[1]))
		fields[fieldName] = value
		if valueError != nil && parseError == nil {
			parseError = fmt.Errorf("frontmatter line %d field %s: %w", lineIndex+2, fieldName, valueError)
		}
	}
	bodyStart := 4 + end + 4
	if bodyStart < len(content) && content[bodyStart] == '\n' {
		bodyStart++
	}
	return fields, content[bodyStart:], parseError
}

func parseFrontmatterScalar(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	if strings.HasPrefix(value, "\"") {
		return parseDoubleQuotedYAMLScalar(value)
	}
	if strings.HasPrefix(value, "'") {
		return parseSingleQuotedYAMLScalar(value)
	}
	return strings.Trim(value, "[]"), nil
}

func parseDoubleQuotedYAMLScalar(value string) (string, error) {
	for index := 1; index < len(value); index++ {
		if value[index] == '\\' {
			index++
			continue
		}
		if value[index] != '"' {
			continue
		}
		remainder := strings.TrimSpace(value[index+1:])
		if remainder != "" && !strings.HasPrefix(remainder, "#") {
			return strings.ReplaceAll(value[1:index], `\"`, `"`), fmt.Errorf("unexpected content after closing double quote")
		}
		return strings.ReplaceAll(value[1:index], `\"`, `"`), nil
	}
	return strings.TrimPrefix(value, "\""), fmt.Errorf("unterminated double-quoted scalar")
}

func parseSingleQuotedYAMLScalar(value string) (string, error) {
	for index := 1; index < len(value); index++ {
		if value[index] != '\'' {
			continue
		}
		if index+1 < len(value) && value[index+1] == '\'' {
			index++
			continue
		}
		remainder := strings.TrimSpace(value[index+1:])
		if remainder != "" && !strings.HasPrefix(remainder, "#") {
			return strings.ReplaceAll(value[1:index], "''", "'"), fmt.Errorf("unexpected content after closing single quote")
		}
		return strings.ReplaceAll(value[1:index], "''", "'"), nil
	}
	return strings.TrimPrefix(value, "'"), fmt.Errorf("unterminated single-quoted scalar")
}

func parseRelations(content string) []pageRelation {
	lines, relations := strings.Split(content, "\n"), []pageRelation{}
	currentPage := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-"))
		if strings.HasPrefix(trimmed, "page:") {
			currentPage = strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "page:")), "\"'")
		} else if strings.HasPrefix(trimmed, "rel:") && currentPage != "" {
			relations = append(relations, pageRelation{page: currentPage, relation: strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "rel:")), "\"'")})
			currentPage = ""
		}
	}
	return relations
}

func normalizedLinks(content string) []string {
	seen := map[string]bool{}
	links := []string{}
	for _, match := range wikiLinkPattern.FindAllStringSubmatch(content, -1) {
		link := normalizeLink(match[1])
		if link != "" && !seen[link] {
			seen[link] = true
			links = append(links, link)
		}
	}
	sort.Strings(links)
	return links
}
func normalizeLink(link string) string {
	link = strings.TrimSpace(strings.SplitN(link, "|", 2)[0])
	link = strings.TrimSuffix(link, ".md")
	return filepath.Base(filepath.FromSlash(link))
}

func normalizeTypeEnum(value string) string {
	return strings.NewReplacer("source_summary", "source-summary", "daily_log", "daily-log", "monthly_rollup", "monthly-rollup").Replace(value)
}

func normalizeRelationshipEnum(value string) string {
	value = strings.ReplaceAll(value, "_", "-")
	if value == "extend" {
		return "extends"
	}
	if value == "supersede" {
		return "supersedes"
	}
	return value
}
func firstHeading(body, fallback string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return fallback
}

// defaultKnowledgeBase is where a knowledge base lives when no path is
// supplied. legacyKnowledgeBases are the earlier root-folder defaults: they are
// still found, so an existing consumer never gets a second, empty scaffold.
const defaultKnowledgeBase = "docs/kb"

var legacyKnowledgeBases = []string{"kb", "knowledge-base"}

// knowledgeBaseLocation is the repository-relative knowledge base one command
// run uses. legacy names a root legacy folder the default order met: equal to
// relative when that folder was chosen, and different when docs/kb won over it
// (both are knowledge bases). defaultOccupied marks a chosen legacy folder whose
// move target, docs/kb, already exists as something that is not a knowledge base.
type knowledgeBaseLocation struct {
	relative        string
	legacy          string
	defaultOccupied bool
}

// resolveKnowledgeBaseLocation applies the one location order every knowledge
// command shares. A supplied path is taken literally. With none: docs/kb when
// it is a knowledge base, else a legacy root folder that is one (or that holds
// raw/ or wiki/ and so was filled in place by the old default), else docs/kb
// (to create or to report as missing). A docs/kb folder that is not a
// knowledge base never hides a legacy one.
func resolveKnowledgeBaseLocation(root, requested string) knowledgeBaseLocation {
	if requested != "" {
		return knowledgeBaseLocation{relative: requested}
	}
	completeLegacy, partialLegacy := "", ""
	for _, candidate := range legacyKnowledgeBases {
		path := filepath.Join(root, candidate)
		if isKnowledgeBase(path) {
			completeLegacy = candidate
			break
		}
		if partialLegacy == "" && holdsKnowledgeBaseFolder(path) {
			partialLegacy = candidate
		}
	}
	defaultPath := filepath.Join(root, filepath.FromSlash(defaultKnowledgeBase))
	if isKnowledgeBase(defaultPath) {
		return knowledgeBaseLocation{relative: defaultKnowledgeBase, legacy: completeLegacy}
	}
	legacy := completeLegacy
	if legacy == "" {
		legacy = partialLegacy
	}
	if legacy == "" {
		return knowledgeBaseLocation{relative: defaultKnowledgeBase}
	}
	_, statError := os.Lstat(defaultPath)
	return knowledgeBaseLocation{relative: legacy, legacy: legacy, defaultOccupied: statError == nil}
}

// holdsKnowledgeBaseFolder reports a folder that has raw/ or wiki/ but is not
// yet a complete knowledge base.
func holdsKnowledgeBaseFolder(path string) bool {
	for _, name := range []string{"raw", "wiki"} {
		if info, err := os.Stat(filepath.Join(path, name)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

// notices tells the user when the default order met a legacy root folder: the
// one-line move when that folder is in use, a warning when docs/kb and the
// legacy folder are both knowledge bases.
func (location knowledgeBaseLocation) notices(commandName string) []resultmodel.CommandFinding {
	if location.legacy == "" {
		return nil
	}
	if location.relative == location.legacy {
		evidence := fmt.Sprintf("using %s/, the earlier default location; the default is now %s/. Move it with: mkdir -p docs && git mv %s %s", location.legacy, defaultKnowledgeBase, location.legacy, defaultKnowledgeBase)
		if location.defaultOccupied {
			evidence = fmt.Sprintf("using %s/, the earlier default location; the default is now %s/, but %s/ exists and is not a knowledge base. Remove that folder first (move its files elsewhere), then: mkdir -p docs && git mv %s %s", location.legacy, defaultKnowledgeBase, defaultKnowledgeBase, location.legacy, defaultKnowledgeBase)
		}
		finding := knowledgeFinding(commandName, "BKB-LEGACY-LOCATION", resultmodel.SeverityInfo, nil, evidence, resultmodel.FixabilityManual, "moving the folder is the user's decision")
		finding.AffectedPaths = []string{location.legacy}
		return []resultmodel.CommandFinding{finding}
	}
	evidence := fmt.Sprintf("both %s/ and %s/ are knowledge bases; %s/ is used. Merge %s/ into it, or pass --kb %s to keep using the old folder", location.relative, location.legacy, location.relative, location.legacy, location.legacy)
	finding := knowledgeFinding(commandName, "BKB-DUPLICATE-LOCATION", resultmodel.SeverityWarning, nil, evidence, resultmodel.FixabilityManual, "which folder holds the real knowledge base is the user's decision")
	finding.AffectedPaths = []string{location.relative, location.legacy}
	return []resultmodel.CommandFinding{finding}
}

// locateKnowledgeBase finds the knowledge base to read. An empty request means
// the default order. On failure the second value is the path that was expected.
func locateKnowledgeBase(root, requested string) (string, string, error) {
	checked := requested
	if requested == "" {
		requested = resolveKnowledgeBaseLocation(root, "").relative
		checked = strings.Join(append([]string{defaultKnowledgeBase}, legacyKnowledgeBases...), ", ")
	}
	path, err := ensureSafeTarget(root, requested)
	if err != nil {
		return "", requested, err
	}
	if isKnowledgeBase(path) {
		return path, filepath.ToSlash(requested), nil
	}
	current, _ := filepath.Abs(root)
	for depth := 0; depth < 3; depth++ {
		if isKnowledgeBase(current) {
			relative, _ := filepath.Rel(root, current)
			return current, filepath.ToSlash(relative), nil
		}
		current = filepath.Dir(current)
	}
	return "", requested, fmt.Errorf("no knowledge base found; checked %s", checked)
}

func scanFinding(commandName, target, code string, severity resultmodel.FindingSeverity, affected []string, evidence string, fixability resultmodel.FindingFixability, stopReason string) resultmodel.CommandFinding {
	finding := knowledgeFinding(commandName, code, severity, []string{target}, evidence, fixability, stopReason)
	finding.AffectedPaths = append([]string(nil), affected...)
	return finding
}
func scanFailure(commandName, target string, err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{scanFinding(commandName, target, "KNOWLEDGE-SCAN-FAILED", resultmodel.SeverityError, []string{target}, err.Error(), resultmodel.FixabilityManual, "the canonical scan could not complete")}}
}
func relativePath(target, relative, path string) string {
	suffix, _ := filepath.Rel(target, path)
	return filepath.ToSlash(filepath.Join(relative, suffix))
}
func containsNormalizedLink(content, stem string) bool {
	for _, link := range normalizedLinks(content) {
		if link == stem {
			return true
		}
	}
	return false
}
func extractNamedCount(content, label string) (int, bool) {
	inspection := inspectNamedCounts(content, label)
	return inspection.value, inspection.occurrences > 0 && !inspection.malformed
}

type namedCountInspection struct {
	value       int
	values      []int
	occurrences int
	malformed   bool
}

func inspectNamedCounts(content, label string) namedCountInspection {
	inspection := namedCountInspection{}
	remaining := content
	for {
		index := strings.Index(remaining, label)
		if index < 0 {
			break
		}
		inspection.occurrences++
		rest := strings.TrimSpace(remaining[index+len(label):])
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			inspection.malformed = true
		} else if value, err := strconv.Atoi(strings.Trim(fields[0], "|")); err != nil || value < 0 {
			inspection.malformed = true
		} else {
			inspection.values = append(inspection.values, value)
			if len(inspection.values) == 1 {
				inspection.value = value
			}
		}
		remaining = remaining[index+len(label):]
	}
	return inspection
}

func namedCountFindings(relative, masterPath, kind, label string, inspection namedCountInspection, diskCount int) []resultmodel.CommandFinding {
	findings := []resultmodel.CommandFinding{}
	prefix := "BKB-STATUS-" + kind + "-COUNT-"
	if inspection.occurrences == 0 {
		return []resultmodel.CommandFinding{scanFinding(CommandBKBStatus, relative, prefix+"MISSING", resultmodel.SeverityWarning, []string{masterPath}, "master index has no "+label+" declaration", resultmodel.FixabilityManual, "add one canonical master-index count declaration")}
	}
	if inspection.malformed || len(inspection.values) != inspection.occurrences {
		findings = append(findings, scanFinding(CommandBKBStatus, relative, prefix+"MALFORMED", resultmodel.SeverityWarning, []string{masterPath}, "master index has an unparseable "+label+" declaration", resultmodel.FixabilityManual, "repair every declared master-index count"))
	}
	if inspection.occurrences > 1 {
		findings = append(findings, scanFinding(CommandBKBStatus, relative, prefix+"DUPLICATE", resultmodel.SeverityWarning, []string{masterPath}, fmt.Sprintf("master index has %d %s declarations", inspection.occurrences, label), resultmodel.FixabilityManual, "retain exactly one canonical master-index count declaration"))
	}
	if len(inspection.values) > 1 {
		for _, value := range inspection.values[1:] {
			if value != inspection.values[0] {
				findings = append(findings, scanFinding(CommandBKBStatus, relative, prefix+"INCONSISTENT", resultmodel.SeverityWarning, []string{masterPath}, fmt.Sprintf("master-index %s declarations disagree: %v", label, inspection.values), resultmodel.FixabilityManual, "reconcile duplicate declarations to one truthful count"))
				break
			}
		}
	}
	if !inspection.malformed && len(inspection.values) > 0 && inspection.value != diskCount {
		findings = append(findings, scanFinding(CommandBKBStatus, relative, prefix+"DISK-MISMATCH", resultmodel.SeverityWarning, []string{masterPath}, fmt.Sprintf("master declares %d but disk inventory has %d", inspection.value, diskCount), resultmodel.FixabilityManual, "rebuild the master index or reconcile the disk inventory"))
	}
	return findings
}
func countRegularEntries(path string) int {
	entries, _ := os.ReadDir(path)
	count := 0
	for _, entry := range entries {
		if info, err := entry.Info(); err == nil && info.Mode().IsRegular() {
			count++
		}
	}
	return count
}
func countLinesContaining(path, needle string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(needle)) {
			count++
		}
	}
	return count
}
func latestNamedDate(path string) string {
	entries, _ := os.ReadDir(path)
	latest := ""
	for _, entry := range entries {
		date := datePattern.FindString(entry.Name())
		if date > latest {
			latest = date
		}
	}
	return latest
}
func latestLogDate(content, marker string) string {
	latest := ""
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, marker) {
			if date := datePattern.FindString(line); date > latest {
				latest = date
			}
		}
	}
	return latest
}
func displayDate(value string) string {
	if value == "" {
		return "never"
	}
	return value
}
func daysSince(value string) int {
	if value == "" {
		return 1 << 30
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	if err != nil {
		return 1 << 30
	}
	return int(nowUTC().Sub(parsed).Hours() / 24)
}
func countCustomAgents(path string) int {
	entries, _ := filepath.Glob(filepath.Join(path, "*.md"))
	count := 0
	for _, entry := range entries {
		data, _ := os.ReadFile(entry)
		if strings.Contains(string(data), "## Custom Agent") {
			count++
		}
	}
	return count
}
func agentStalenessFindings(target, relative string) []resultmodel.CommandFinding {
	data, err := os.ReadFile(filepath.Join(target, "wiki", "agent.md"))
	if err != nil {
		return nil
	}
	text := string(data)
	queryIndex := strings.Index(text, "## Query Log")
	if queryIndex < 0 {
		return nil
	}
	hotIndex := strings.Index(text, "## Hot Topics")
	queries := 0
	for _, line := range strings.Split(text[queryIndex:], "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "|") && !strings.Contains(line, "---") && !strings.Contains(line, "Date | Query") {
			queries++
		}
	}
	if queries < 10 {
		return nil
	}
	hot := ""
	if hotIndex >= 0 && hotIndex < queryIndex {
		hot = text[hotIndex:queryIndex]
	}
	if strings.Contains(hot, "(none yet") || strings.TrimSpace(strings.TrimPrefix(hot, "## Hot Topics")) == "" {
		path := filepath.ToSlash(filepath.Join(relative, "wiki/agent.md"))
		return []resultmodel.CommandFinding{scanFinding(CommandBKBLintStructure, relative, "BKB-AGENT-QUERY-STALE", resultmodel.SeverityWarning, []string{path}, fmt.Sprintf("Query Log has %d entries but Hot Topics were not regenerated", queries), resultmodel.FixabilityManual, "the action synthesizes hot topics")}
	}
	return nil
}
