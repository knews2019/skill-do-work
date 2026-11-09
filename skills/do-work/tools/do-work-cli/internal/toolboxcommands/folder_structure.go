package toolboxcommands

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// folder-structure-check is the mechanical half of the toolbox folder-structure action. It is
// read-only and reports where a repository differs from that action's Target Layout. Every gap
// is a warning finding and the outcome is "findings" (exit 1), so a project can make it blocking
// by calling it from CI; the judgment (exceptions, moves, REQs) stays in the action.

const CommandFolderStructureCheck = "folder-structure-check"

// layoutRootFolders are the folders the Target Layout defines at the root. This is the layout's
// own definition, not a sample, so it changes only together with the action file.
var layoutRootFolders = map[string]bool{
	"cmd": true, "src": true, "data": true, "docs": true, "tests": true, "ai-reports": true, "do-work": true,
}

// layoutRootDocuments are compared by the name before its first dot, case-insensitively, so
// README.md, LICENSE and CHANGELOG.md all match.
var layoutRootDocuments = map[string]bool{
	"readme": true, "license": true, "changelog": true, "claude": true, "agents": true, "justfile": true,
}

var appFolderName = regexp.MustCompile(`^(be|fe|mobile)-[a-z0-9]+(-[a-z0-9]+)*$`)

// sqliteHeader is the first 16 bytes of every SQLite 3 database file. Matching the content, not
// a file extension, finds a database whatever its name.
var sqliteHeader = []byte("SQLite format 3\x00")

type folderStructureGap struct {
	code, path, detail string
}

func handleFolderStructureCheck(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	allowed := map[string]bool{}
	for index := 0; index < len(arguments); index++ {
		if arguments[index] != "--allow" && !strings.HasPrefix(arguments[index], "--allow=") {
			return usageResult(CommandFolderStructureCheck, "Usage: folder-structure-check [--allow <root-entry>]...")
		}
		value, err := optionValue(arguments, &index, "--allow")
		if err != nil {
			return usageResult(CommandFolderStructureCheck, err.Error())
		}
		allowed[strings.TrimSuffix(value, "/")] = true
	}
	root := executionContext.RepositoryRoot
	listing, err := auditGit(root, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	if err != nil {
		return folderStructureFailure(err)
	}
	trackedListing, err := auditGit(root, "ls-files", "-z", "--cached")
	if err != nil {
		return folderStructureFailure(err)
	}
	visible := splitNul(listing)
	tracked := map[string]bool{}
	for _, path := range splitNul(trackedListing) {
		tracked[path] = true
	}

	gaps := rootAndSourceGaps(visible, allowed)
	gaps = append(gaps, dataFolderGaps(visible)...)
	databaseGaps, err := databasesOutsideData(root, tracked)
	if err != nil {
		return folderStructureFailure(err)
	}
	gaps = append(gaps, databaseGaps...)
	gaps = append(gaps, justfileLogicGaps(root)...)
	inbox := inboxFiles(visible)
	return folderStructureResult(gaps, inbox)
}

func splitNul(listing string) []string {
	paths := []string{}
	for _, path := range strings.Split(listing, "\x00") {
		if path != "" {
			paths = append(paths, path)
		}
	}
	return paths
}

// rootAndSourceGaps reads only Git-visible paths, so an ignored local folder such as
// node_modules/ is never reported as root clutter.
func rootAndSourceGaps(visible []string, allowed map[string]bool) []folderStructureGap {
	seenRoot, seenSource := map[string]bool{}, map[string]bool{}
	gaps := []folderStructureGap{}
	for _, path := range visible {
		segments := strings.Split(path, "/")
		entry := segments[0]
		if !seenRoot[entry] {
			seenRoot[entry] = true
			if !rootEntryAllowed(entry, len(segments) > 1, allowed) {
				display := entry
				if len(segments) > 1 {
					display += "/"
				}
				gaps = append(gaps, folderStructureGap{"FOLDER-STRUCTURE-ROOT-ENTRY", display,
					"not part of the root layout (cmd, src, data, docs, tests, ai-reports, do-work, root documents, dot-config)"})
			}
		}
		if entry == "src" && len(segments) > 1 && !seenSource[segments[1]] {
			seenSource[segments[1]] = true
			if len(segments) == 2 || !appFolderName.MatchString(segments[1]) {
				gaps = append(gaps, folderStructureGap{"FOLDER-STRUCTURE-SRC-ENTRY", "src/" + segments[1],
					"src/ holds only app folders named <role>-<tech>, with role be, fe or mobile (for example src/be-golang)"})
			}
		}
	}
	return gaps
}

func rootEntryAllowed(entry string, isFolder bool, allowed map[string]bool) bool {
	if allowed[entry] || strings.HasPrefix(entry, ".") {
		return true
	}
	if isFolder {
		return layoutRootFolders[entry]
	}
	stem, _, _ := strings.Cut(strings.ToLower(entry), ".")
	return layoutRootDocuments[stem]
}

// dataFolderGaps flags anything under data/ that Git would commit: data/ is gitignored except
// data/README.md, which names each app's files.
func dataFolderGaps(visible []string) []folderStructureGap {
	gaps := []folderStructureGap{}
	for _, path := range visible {
		if strings.HasPrefix(path, "data/") && path != "data/README.md" {
			gaps = append(gaps, folderStructureGap{"FOLDER-STRUCTURE-DATA-NOT-IGNORED", path,
				"data/ must be gitignored except data/README.md; add data/* and !data/README.md to .gitignore"})
		}
	}
	return gaps
}

// databasesOutsideData walks the working tree for SQLite files that Git does not track. A tracked
// database is a deliberate fixture; an untracked one was written at runtime and belongs under
// data/<app>/. The walk skips .git, the root data/ folder, symlinks, the do-work/ and ai-reports/
// record folders (their evidence may copy a database on purpose), and node_modules folders, which
// hold package-manager downloads rather than app data.
func databasesOutsideData(root string, tracked map[string]bool) ([]folderStructureGap, error) {
	gaps := []folderStructureGap{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if entry.IsDir() {
			if relative == ".git" || relative == "data" || relative == "do-work" || relative == "ai-reports" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() || tracked[relative] {
			return nil
		}
		if isSQLiteFile(path) {
			gaps = append(gaps, folderStructureGap{"FOLDER-STRUCTURE-DATA-OUTSIDE", relative,
				"SQLite database outside data/; store it under data/<app>/ and point the app at it with an env var"})
		}
		return nil
	})
	return gaps, err
}

func isSQLiteFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	header := make([]byte, len(sqliteHeader))
	read, _ := file.Read(header)
	return read == len(sqliteHeader) && bytes.Equal(header, sqliteHeader)
}

var justLeadingName = regexp.MustCompile(`^@?([A-Za-z_][A-Za-z0-9_-]*)`)

// justRecipeName returns the recipe name when an unindented justfile line is a recipe header, and
// "" otherwise. A header is a name, optional parameters, then a ":" that is not ":=". The scan
// skips quoted text, because a parameter default may hold "=" or ":" (run-kanban $port="8090":,
// url='http://localhost'); stopping at the first "=" used to hide such recipes. A line whose first
// unquoted colon is ":=" is a variable assignment.
func justRecipeName(line string) string {
	match := justLeadingName.FindStringSubmatch(line)
	if match == nil || justDirective(match[1]) {
		return ""
	}
	rest := line[len(match[0]):]
	var quote rune
	for index, character := range rest {
		switch {
		case quote != 0:
			if character == quote {
				quote = 0
			}
		case character == '"' || character == '\'' || character == '`':
			quote = character
		case character == ':':
			if strings.HasPrefix(rest[index+1:], "=") {
				return ""
			}
			return match[1]
		}
	}
	return ""
}

// justfileLogicGaps flags a recipe whose body has more than one command line. Logic belongs in
// cmd/<name>.sh and the recipe calls it in one line.
func justfileLogicGaps(root string) []folderStructureGap {
	for _, name := range []string{"justfile", "Justfile", ".justfile"} {
		content, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			continue
		}
		gaps := []folderStructureGap{}
		recipe, bodyLines := "", 0
		flush := func() {
			if recipe != "" && bodyLines > 1 {
				gaps = append(gaps, folderStructureGap{"FOLDER-STRUCTURE-JUST-LOGIC", name + ": " + recipe,
					fmt.Sprintf("recipe has %d body lines; move them to cmd/%s.sh and call it in one line", bodyLines, recipe)})
			}
			recipe, bodyLines = "", 0
		}
		for _, line := range strings.Split(string(content), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if line[0] == ' ' || line[0] == '\t' {
				if recipe != "" {
					bodyLines++
				}
				continue
			}
			flush()
			recipe = justRecipeName(line)
		}
		flush()
		return gaps
	}
	return nil
}

func justDirective(word string) bool {
	switch word {
	case "set", "alias", "export", "import", "mod":
		return true
	}
	return false
}

func inboxFiles(visible []string) []string {
	files := []string{}
	for _, path := range visible {
		if strings.HasPrefix(path, "docs/inbox/") {
			files = append(files, path)
		}
	}
	return files
}

func folderStructureResult(gaps []folderStructureGap, inbox []string) resultmodel.CommandResult {
	sort.SliceStable(gaps, func(left, right int) bool {
		if gaps[left].code != gaps[right].code {
			return gaps[left].code < gaps[right].code
		}
		return gaps[left].path < gaps[right].path
	})
	var text strings.Builder
	findings := []resultmodel.CommandFinding{}
	for _, gap := range gaps {
		fmt.Fprintf(&text, "%s  %s  %s\n", gap.code, gap.path, gap.detail)
		findings = append(findings, toolboxFinding(CommandFolderStructureCheck, gap.code, resultmodel.SeverityWarning,
			[]string{gap.path}, gap.detail, resultmodel.FixabilityManual, "a layout move needs a tidy-repo plan or a REQ"))
	}
	if len(inbox) > 0 {
		fmt.Fprintf(&text, "docs/inbox/ has %d file(s) waiting to be sorted\n", len(inbox))
		findings = append(findings, toolboxFinding(CommandFolderStructureCheck, "FOLDER-STRUCTURE-INBOX", resultmodel.SeverityInfo,
			inbox, fmt.Sprintf("%d file(s) in docs/inbox/", len(inbox)), resultmodel.FixabilityManual, "sorting the inbox is a judgment call"))
	}
	outcome := resultmodel.OutcomeSuccess
	if len(gaps) > 0 {
		outcome = resultmodel.OutcomeFindings
		fmt.Fprintf(&text, "folder-structure-check: %d gap(s)\n", len(gaps))
	} else {
		text.WriteString("folder-structure-check: layout matches\n")
	}
	output := text.String()
	return resultmodel.CommandResult{Outcome: outcome, Findings: findings, ExactTextOutput: &output}
}

func folderStructureFailure(err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandFolderStructureCheck, "FOLDER-STRUCTURE-FAILED", resultmodel.SeverityError, nil, err.Error(),
			resultmodel.FixabilityManual, "run the check inside a Git working tree"),
	}}
}
