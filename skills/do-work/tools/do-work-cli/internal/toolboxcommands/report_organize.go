package toolboxcommands

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

// ai-report-organize files every report of a finished month into ai-reports/YYYY-MM_N-reports/
// (YYYY-MM_1-report for one report), keeps the current month's reports directly under
// ai-reports/, and rewrites every link to a moved report so it still opens. It is the only
// command that moves reports or rewrites links, and nothing runs it automatically.

var (
	reportBundleDatePrefix = regexp.MustCompile(`^[0-9]{4}-(0[1-9]|1[0-2])-[0-9]{2}`)
	reportLinkScheme       = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
	reportLinkAttribute    = regexp.MustCompile(`(?i)\b(?:href|src)\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	reportMarkdownTarget   = regexp.MustCompile(`\]\(([^)\s]+)`)
)

// A repository path mention ends at whitespace or at one of these bytes.
const reportMentionStopBytes = "\"'`<>()[]{}|,;*:#?\\"

// reportOrganizeHintLine is the one line suite install, update, `ai-report index` and every new
// report print while organize has work to do. It recommends; nothing runs organize for the user.
const reportOrganizeHintLine = "ai-reports/ needs organizing by month: run do-work-toolbox ai-report organize (add --dry-run to preview)."

// ReportOrganizeHint returns reportOrganizeHintLine when organize would move, rename or remove
// anything under repositoryRoot's ai-reports/, else "". It reads only ai-reports/ and writes
// nothing. A check error (a symlinked ai-reports, say) returns "", so the hint never changes a
// caller's outcome.
func ReportOrganizeHint(repositoryRoot string) string {
	plan, err := planReportOrganizeLayout(repositoryRoot, toolboxNow().Format("2006-01"))
	if err != nil || !plan.hasWork() {
		return ""
	}
	return reportOrganizeHintLine
}

// reportOrganizeStep is one mkdir, move, rename or remove. Paths are repository-relative slash
// paths as they are when the step runs: steps run mkdirs, then moves, then renames, then removes.
type reportOrganizeStep struct {
	Action, From, To string
}

type reportOrganizeLeft struct {
	Path, Reason string
}

// reportOrganizeEdit is one file whose links change. Path is where it is now, FinalPath where it
// ends up after the moves.
type reportOrganizeEdit struct {
	Path, FinalPath string
	Before, After   []byte
	Links           int
}

type reportOrganizePlan struct {
	Steps []reportOrganizeStep
	Left  []reportOrganizeLeft
	Edits []reportOrganizeEdit
	// CommitPaths are the step and edit paths Git knows a file under, so `git add -A -- <paths>`
	// never meets a pathspec that matches nothing.
	CommitPaths []string
}

// hasWork reports whether organize would move, rename or remove anything. It is the whole check
// a caller needs to decide whether to recommend organize.
func (plan reportOrganizePlan) hasWork() bool {
	return len(plan.Steps) > 0
}

// finalPath returns where relative ends up once every move and rename has run, matching each step
// at a path segment boundary only, so ai-reports/X-rev1 never follows a move of ai-reports/X.
func (plan reportOrganizePlan) finalPath(relative string) string {
	for _, step := range plan.Steps {
		if step.Action != "move" && step.Action != "rename" {
			continue
		}
		if relative == step.From || strings.HasPrefix(relative, step.From+"/") {
			relative = step.To + relative[len(step.From):]
		}
	}
	return relative
}

func handleAIReportOrganize(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	// --check is the report writers' seam: the hint line or nothing, never a write.
	if len(arguments) == 1 && arguments[0] == "--check" {
		hint := ReportOrganizeHint(executionContext.RepositoryRoot)
		if hint != "" {
			hint += "\n"
		}
		return exactOutputResult(hint, nil)
	}
	dryRun := len(arguments) == 1 && arguments[0] == "--dry-run"
	if len(arguments) > 0 && !dryRun {
		return usageResult(CommandAIReportOrganize, "Usage: ai-report-organize [--dry-run | --check]")
	}
	repositoryRoot := executionContext.RepositoryRoot
	plan, err := planReportOrganizeLayout(repositoryRoot, toolboxNow().Format("2006-01"))
	if err != nil {
		return reportOrganizeFailure(nil, err)
	}
	var output strings.Builder
	skippedWork := []resultmodel.SkippedWork{}
	leftLines := ""
	for _, left := range plan.Left {
		line := fmt.Sprintf("%s (%s)", left.Path, left.Reason)
		leftLines += "left " + line + "\n"
		skippedWork = append(skippedWork, resultmodel.SkippedWork{Code: "AI-REPORT-ORGANIZE-LEFT", Reason: line})
	}
	if !plan.hasWork() {
		text := "nothing to organize\n" + leftLines
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, SkippedWork: skippedWork, ExactTextOutput: &text}
	}
	// The catalog is regenerated last, so a hand-made one must refuse before the first write.
	if findings := reportCatalogHandMadeFindings(repositoryRoot, CommandAIReportOrganize); len(findings) > 0 {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: findings}
	}
	if err := planReportOrganizeLinks(repositoryRoot, &plan); err != nil {
		return reportOrganizeFailure(nil, err)
	}
	dirtyPaths, err := reportOrganizeDirtyPaths(repositoryRoot, plan)
	if err != nil {
		return reportOrganizeFailure(nil, err)
	}
	if len(dirtyPaths) > 0 {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{
			toolboxFinding(CommandAIReportOrganize, "AI-REPORT-ORGANIZE-DIRTY", resultmodel.SeverityError, dirtyPaths,
				"organize would move or rewrite files that have uncommitted changes; commit or restore them, then rerun",
				resultmodel.FixabilityManual, "the organize commit takes exact paths and must not sweep in another edit"),
		}}
	}
	if dryRun {
		output.WriteString("dry run: nothing written\n")
	}
	for _, step := range plan.Steps {
		switch step.Action {
		case "move":
			fmt.Fprintf(&output, "moved %s -> %s\n", step.From, plan.finalPath(step.From))
		case "rename":
			fmt.Fprintf(&output, "renamed %s -> %s\n", step.From, step.To)
		case "remove":
			fmt.Fprintf(&output, "removed %s (empty month folder)\n", step.From)
		}
	}
	for _, edit := range plan.Edits {
		fmt.Fprintf(&output, "links %s (%d)\n", edit.FinalPath, edit.Links)
	}
	output.WriteString(leftLines)
	if dryRun {
		text := output.String()
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, SkippedWork: skippedWork, ExactTextOutput: &text}
	}
	changes, err := applyReportOrganize(repositoryRoot, plan)
	if err != nil {
		return reportOrganizeFailure(changes, err)
	}
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, SkippedWork: skippedWork}
	commitPaths := append([]string(nil), plan.CommitPaths...)
	if catalog := reportIndexWrite(repositoryRoot); catalog.Outcome == resultmodel.OutcomeSuccess {
		output.WriteString(*catalog.ExactTextOutput)
		changes = append(changes, catalog.Changes...)
		// The same rule as CommitPaths: a catalog file Git ignores would make `git add -A` reject
		// the whole list, so only the files Git would take are listed.
		arguments := []string{"ls-files", "-z", "--cached", "--others", "--exclude-standard", "--"}
		for _, change := range catalog.Changes {
			arguments = append(arguments, change.Path)
		}
		listed, err := auditGit(repositoryRoot, arguments...)
		if err != nil {
			return reportOrganizeFailure(changes, err)
		}
		for _, catalogPath := range strings.Split(listed, "\x00") {
			if catalogPath != "" {
				commitPaths = append(commitPaths, catalogPath)
			}
		}
	} else {
		// A hand-made catalog was refused before any write, so this is a write failure after the
		// moves. They stay, and their commit paths are still printed; `ai-report index` reruns
		// the catalog once the cause is fixed.
		result.Outcome, result.Findings = resultmodel.OutcomeFindings, catalog.Findings
	}
	output.WriteString("commit paths:\n")
	sort.Strings(commitPaths)
	for index, commitPath := range commitPaths {
		if index == 0 || commitPath != commitPaths[index-1] {
			output.WriteString(commitPath + "\n")
		}
	}
	text := output.String()
	result.Changes, result.ExactTextOutput = changes, &text
	return result
}

func reportOrganizeFailure(changes []resultmodel.RecordedChange, err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Changes: changes, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandAIReportOrganize, "AI-REPORT-ORGANIZE-FAILED", resultmodel.SeverityError, nil, err.Error(),
			resultmodel.FixabilityManual, "organize stopped; any change listed was made: commit or restore it with Git, then rerun"),
	}}
}

// reportBundleMonth returns the YYYY-MM of a bundle folder name's date prefix, else fallback.
func reportBundleMonth(name, fallback string) string {
	if reportBundleDatePrefix.MatchString(name) {
		return name[:7]
	}
	return fallback
}

// planReportOrganizeLayout reads ai-reports/ only (no Git, no other files) and plans the steps
// that put every finished month in one counted month folder and every other report at the top.
func planReportOrganizeLayout(repositoryRoot, currentMonth string) (reportOrganizePlan, error) {
	var plan reportOrganizePlan
	info, err := os.Lstat(filepath.Join(repositoryRoot, reportsFolder))
	if os.IsNotExist(err) {
		return plan, nil
	}
	if err != nil {
		return plan, err
	}
	if !info.IsDir() {
		return plan, fmt.Errorf("%s is a symbolic link or a file, not a folder; organize moves reports only inside a real folder", reportsFolder)
	}
	entries, err := os.ReadDir(filepath.Join(repositoryRoot, reportsFolder))
	if err != nil {
		return plan, err
	}
	type reportSpot struct{ folder, name string } // folder is "" for a bundle directly under ai-reports/
	bundlesByMonth := map[string][]reportSpot{}
	foldersByMonth := map[string][]string{}
	folderChildren := map[string][]os.DirEntry{}
	var monthFolders []string
	topLevel := map[string]bool{}
	for _, entry := range entries {
		name, relative := entry.Name(), reportsFolder+"/"+entry.Name()
		topLevel[name] = true
		switch {
		case strings.HasPrefix(name, "."), !entry.IsDir() && (relative == reportCatalogFile || relative == reportIndexPageFile):
		case entry.Type()&os.ModeSymlink != 0:
			plan.Left = append(plan.Left, reportOrganizeLeft{relative, "symbolic link"})
		case !entry.IsDir():
			plan.Left = append(plan.Left, reportOrganizeLeft{relative, "loose file"})
		case reportMonthFolder(name):
			children, err := os.ReadDir(filepath.Join(repositoryRoot, reportsFolder, name))
			if err != nil {
				return plan, err
			}
			folderMonth := name[:7]
			monthFolders = append(monthFolders, name)
			foldersByMonth[folderMonth] = append(foldersByMonth[folderMonth], name)
			folderChildren[name] = children
			for _, child := range children {
				if child.IsDir() && !strings.HasPrefix(child.Name(), ".") {
					month := reportBundleMonth(child.Name(), folderMonth)
					bundlesByMonth[month] = append(bundlesByMonth[month], reportSpot{name, child.Name()})
				}
			}
		case reportBundleMonth(name, "") != "":
			bundlesByMonth[name[:7]] = append(bundlesByMonth[name[:7]], reportSpot{"", name})
		default:
			plan.Left = append(plan.Left, reportOrganizeLeft{relative, "no date prefix"})
		}
	}
	months := []string{}
	for month := range bundlesByMonth {
		months = append(months, month)
	}
	for month := range foldersByMonth {
		if bundlesByMonth[month] == nil {
			months = append(months, month)
		}
	}
	sort.Strings(months)

	var mkdirs, moves, renames, removes []reportOrganizeStep
	keepers, skippedMonths, moving, plannedTargets := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, month := range months {
		bundles := bundlesByMonth[month]
		finished := month < currentMonth
		target, keeper, blocker, reason := "", "", "", ""
		if finished && len(bundles) > 0 {
			target = fmt.Sprintf("%s_%d-reports", month, len(bundles))
			if len(bundles) == 1 {
				target = month + "_1-report"
			}
			for _, folder := range foldersByMonth[month] {
				if folder == target {
					keeper = folder
				}
			}
			if keeper == "" && len(foldersByMonth[month]) > 0 {
				keeper = foldersByMonth[month][0]
			}
			if keeper != target && topLevel[target] {
				blocker, reason = reportsFolder+"/"+target, "already exists and is not a month folder"
			}
		}
		// Moves land in the keeper under its current name, in a new counted folder, or at the top.
		destination := keeper
		if finished && keeper == "" {
			destination = target
		}
		var monthMoves []reportOrganizeStep
		names := map[string]bool{}
		for _, spot := range bundles {
			if blocker != "" {
				break
			}
			from := path.Join(reportsFolder, spot.folder, spot.name)
			if names[spot.name] {
				blocker, reason = from, "another report of this month has the same name"
				break
			}
			names[spot.name] = true
			if spot.folder == destination {
				continue
			}
			to := path.Join(reportsFolder, destination, spot.name)
			occupied := topLevel[spot.name]
			if destination != "" {
				occupied = false
				for _, child := range folderChildren[destination] {
					occupied = occupied || child.Name() == spot.name
				}
			}
			if occupied {
				blocker, reason = to, "already exists"
				break
			}
			if plannedTargets[to] {
				// Two months' folders can hold the same undated name; the later month waits.
				blocker, reason = to, "another month's report moves here"
				break
			}
			monthMoves = append(monthMoves, reportOrganizeStep{Action: "move", From: from, To: to})
		}
		if blocker != "" {
			skippedMonths[month] = true
			plan.Left = append(plan.Left, reportOrganizeLeft{blocker, fmt.Sprintf("%s; month %s not organized", reason, month)})
			continue
		}
		if keeper != "" {
			keepers[keeper] = true
		}
		if finished && keeper == "" && len(monthMoves) > 0 {
			mkdirs = append(mkdirs, reportOrganizeStep{Action: "mkdir", To: reportsFolder + "/" + target})
		}
		for _, move := range monthMoves {
			moving[move.From], plannedTargets[move.To] = true, true
		}
		moves = append(moves, monthMoves...)
		if keeper != "" && keeper != target {
			renames = append(renames, reportOrganizeStep{Action: "rename", From: reportsFolder + "/" + keeper, To: reportsFolder + "/" + target})
		}
	}
	for _, folder := range monthFolders {
		if keepers[folder] {
			continue
		}
		remaining := 0
		for _, child := range folderChildren[folder] {
			if !moving[reportsFolder+"/"+folder+"/"+child.Name()] {
				remaining++
			}
		}
		if remaining == 0 {
			removes = append(removes, reportOrganizeStep{Action: "remove", From: reportsFolder + "/" + folder})
		} else if !skippedMonths[folder[:7]] {
			plan.Left = append(plan.Left, reportOrganizeLeft{reportsFolder + "/" + folder, "month folder holds other files"})
		}
	}
	plan.Steps = append(append(append(mkdirs, moves...), renames...), removes...)
	sort.Slice(plan.Left, func(i, j int) bool { return plan.Left[i].Path < plan.Left[j].Path })
	return plan, nil
}

// planReportOrganizeLinks reads every file Git lists (tracked plus untracked, not ignored) and
// fills plan.Edits and plan.CommitPaths. It writes nothing.
func planReportOrganizeLinks(repositoryRoot string, plan *reportOrganizePlan) error {
	listing, err := auditGit(repositoryRoot, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	if err != nil {
		return fmt.Errorf("organize needs a Git repository: %w", err)
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	files := strings.Split(listing, "\x00")
	sort.Strings(files)
	for index, relative := range files {
		if relative == "" || (index > 0 && relative == files[index-1]) || relative == reportCatalogFile || relative == reportIndexPageFile {
			continue
		}
		info, err := root.Lstat(filepath.FromSlash(relative))
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		content, err := rootedfs.ReadFile(root, filepath.FromSlash(relative))
		if err != nil {
			return err
		}
		// Git's own binary test: a NUL byte in the first 8000 bytes. Binary files are never edited.
		if bytes.IndexByte(content[:min(len(content), 8000)], 0) >= 0 {
			continue
		}
		inReports := strings.HasPrefix(relative, reportsFolder+"/")
		if !inReports && !bytes.Contains(content, []byte(reportsFolder+"/")) {
			continue
		}
		if after, links := plan.rewriteReportLinks(root, relative, content, inReports); links > 0 {
			plan.Edits = append(plan.Edits, reportOrganizeEdit{Path: relative, FinalPath: plan.finalPath(relative), Before: content, After: after, Links: links})
		}
	}
	for _, step := range plan.Steps {
		if step.From == "" {
			continue
		}
		known := false
		for _, file := range files {
			known = known || file == step.From || strings.HasPrefix(file, step.From+"/")
		}
		if known {
			plan.CommitPaths = append(plan.CommitPaths, step.From, plan.finalPath(step.From))
		}
	}
	for _, edit := range plan.Edits {
		plan.CommitPaths = append(plan.CommitPaths, edit.FinalPath)
	}
	return nil
}

// rewriteReportLinks applies the two link rules to one file and returns the new bytes and how
// many links changed.
//
// Rule B (files inside ai-reports/ only): a relative href, src or Markdown link target that
// opens now and would not open from the file's final place is rewritten relative to that place,
// keeping its query, fragment and trailing slash. A link that is broken now, climbs above the
// repository, or still opens keeps its bytes.
//
// Rule A (every file): a repository path mention `ai-reports/...` not glued to a longer name and
// not inside a Rule B link is replaced by its final path.
func (plan reportOrganizePlan) rewriteReportLinks(root *os.Root, relative string, content []byte, inReports bool) ([]byte, int) {
	type replacement struct {
		start, end int
		text       string
	}
	var replacements []replacement
	var relativeLinks [][2]int
	overlapsRelativeLink := func(start, end int) bool {
		for _, span := range relativeLinks {
			if start < span[1] && span[0] < end {
				return true
			}
		}
		return false
	}
	if inReports {
		finalDirectory := path.Dir(plan.finalPath(relative))
		for _, pattern := range []*regexp.Regexp{reportLinkAttribute, reportMarkdownTarget} {
			for _, match := range pattern.FindAllSubmatchIndex(content, -1) {
				start, end := match[2], match[3]
				if start < 0 {
					start, end = match[4], match[5]
				}
				link := string(content[start:end])
				if link == "" || strings.ContainsAny(link[:1], "/#?") || reportLinkScheme.MatchString(link) || overlapsRelativeLink(start, end) {
					continue
				}
				relativeLinks = append(relativeLinks, [2]int{start, end})
				linkPath, suffix := link, ""
				if cut := strings.IndexAny(link, "?#"); cut >= 0 {
					linkPath, suffix = link[:cut], link[cut:]
				}
				target := path.Join(path.Dir(relative), linkPath)
				if target == ".." || strings.HasPrefix(target, "../") {
					continue
				}
				if _, err := root.Lstat(filepath.FromSlash(target)); err != nil {
					continue
				}
				finalTarget := plan.finalPath(target)
				if path.Join(finalDirectory, linkPath) == finalTarget {
					continue
				}
				rewritten, err := filepath.Rel(filepath.FromSlash(finalDirectory), filepath.FromSlash(finalTarget))
				if err != nil {
					continue
				}
				rewritten = filepath.ToSlash(rewritten)
				if strings.HasSuffix(linkPath, "/") && !strings.HasSuffix(rewritten, "/") {
					rewritten += "/"
				}
				replacements = append(replacements, replacement{start, end, rewritten + suffix})
			}
		}
	}
	marker := []byte(reportsFolder + "/")
	for offset := 0; ; {
		found := bytes.Index(content[offset:], marker)
		if found < 0 {
			break
		}
		start := offset + found
		end := start + len(marker)
		offset = end
		if start > 0 && reportNameByte(content[start-1]) {
			continue
		}
		for end < len(content) && content[end] > ' ' && bytes.IndexByte([]byte(reportMentionStopBytes), content[end]) < 0 {
			end++
		}
		for end > start+len(marker) && content[end-1] == '.' {
			end--
		}
		offset = end
		if overlapsRelativeLink(start, end) {
			continue
		}
		mention := string(content[start:end])
		if final := plan.finalPath(mention); final != mention {
			replacements = append(replacements, replacement{start, end, final})
		}
	}
	if len(replacements) == 0 {
		return content, 0
	}
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].start < replacements[j].start })
	var rewritten bytes.Buffer
	previous := 0
	for _, change := range replacements {
		rewritten.Write(content[previous:change.start])
		rewritten.WriteString(change.text)
		previous = change.end
	}
	rewritten.Write(content[previous:])
	return rewritten.Bytes(), len(replacements)
}

// reportNameByte reports whether b can be part of a file name, so `my-ai-reports/` is not a
// mention of ai-reports/.
func reportNameByte(b byte) bool {
	return b == '_' || b == '.' || b == '-' || '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}

// reportOrganizeDirtyPaths lists uncommitted changes, untracked files included, under every path
// organize would move, remove or rewrite. Changes anywhere else are not organize's business.
func reportOrganizeDirtyPaths(repositoryRoot string, plan reportOrganizePlan) ([]string, error) {
	arguments := []string{"--literal-pathspecs", "status", "--porcelain=v1", "-z", "--untracked-files=all", "--"}
	for _, step := range plan.Steps {
		if step.From != "" {
			arguments = append(arguments, step.From)
		}
	}
	for _, edit := range plan.Edits {
		arguments = append(arguments, edit.Path)
	}
	status, err := auditGit(repositoryRoot, arguments...)
	if err != nil {
		return nil, err
	}
	var dirty []string
	records := strings.Split(status, "\x00")
	for index := 0; index < len(records); index++ {
		record := records[index]
		if len(record) < 4 {
			continue
		}
		dirty = append(dirty, record[3:])
		if record[0] == 'R' || record[0] == 'C' {
			index++ // the next record is the rename's source path
		}
	}
	sort.Strings(dirty)
	return dirty, nil
}

// applyReportOrganize writes the edits at the files' current paths first, then runs the steps.
// It stops at the first error and returns the changes made so far; there is no rollback, Git is
// the undo.
func applyReportOrganize(repositoryRoot string, plan reportOrganizePlan) ([]resultmodel.RecordedChange, error) {
	changes := []resultmodel.RecordedChange{}
	for _, edit := range plan.Edits {
		info, err := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(edit.Path)))
		if err != nil {
			return changes, err
		}
		current, err := rootedReadFile(repositoryRoot, edit.Path)
		if err != nil {
			return changes, err
		}
		if !bytes.Equal(current, edit.Before) {
			return changes, fmt.Errorf("%s changed after organize read it", edit.Path)
		}
		if err := rootedPublishFile(repositoryRoot, edit.Path, edit.After, info.Mode().Perm(), true); err != nil {
			return changes, err
		}
		changes = append(changes, resultmodel.RecordedChange{Path: edit.FinalPath, Kind: "modified", Detail: fmt.Sprintf("rewrote %d report links", edit.Links)})
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return changes, err
	}
	defer root.Close()
	for _, step := range plan.Steps {
		if err := applyReportOrganizeStep(repositoryRoot, root, step); err != nil {
			return changes, fmt.Errorf("%s %q -> %q: %w", step.Action, step.From, step.To, err)
		}
		switch step.Action {
		case "move", "rename":
			final := plan.finalPath(step.From)
			changes = append(changes,
				resultmodel.RecordedChange{Path: step.From, Kind: "deleted", Detail: "moved to " + final},
				resultmodel.RecordedChange{Path: final, Kind: "created", Detail: "moved from " + step.From})
		case "remove":
			changes = append(changes, resultmodel.RecordedChange{Path: step.From, Kind: "deleted", Detail: "removed empty month folder"})
		}
	}
	return changes, nil
}

func applyReportOrganizeStep(repositoryRoot string, root *os.Root, step reportOrganizeStep) error {
	switch step.Action {
	case "mkdir":
		_, err := rootedMkdirExclusive(repositoryRoot, step.To, 0o755)
		return err
	case "remove":
		return root.Remove(filepath.FromSlash(step.From))
	}
	// Bind the move to the real directory planned, right before the rename: a symlink swapped in,
	// or anything created at the destination (POSIX rename silently replaces an empty directory),
	// stops the run.
	if err := validateNoLinkedAncestors(repositoryRoot, step.From, true); err != nil {
		return err
	}
	if err := validateNoLinkedAncestors(repositoryRoot, step.To, false); err != nil {
		return err
	}
	if source, err := root.Lstat(filepath.FromSlash(step.From)); err != nil || !source.IsDir() {
		return fmt.Errorf("%s is no longer a real folder", step.From)
	}
	if _, err := root.Lstat(filepath.FromSlash(step.To)); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", step.To)
	}
	return rootedfs.Rename(root, filepath.FromSlash(step.From), filepath.FromSlash(step.To))
}
