package corehelpers

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var requestIDPattern = regexp.MustCompile(`^REQ-[0-9]+$`)

// handleFrontmatterSet writes one scalar field of an active REQ. A field ending
// in _at is an append-only stamp (actions/work-reference.md → Stamps are
// append-only): it is written once, in the Timestamp rule's form, and never
// replaced. Quoting is SetScalar's, so this writer adds no encoding of its own.
func handleFrontmatterSet(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	stampRequested := len(arguments) == 4 && arguments[2] == "--at"
	if len(arguments) != 3 && !stampRequested {
		return frontmatterUsageResult("usage: frontmatter set <file|REQ-N> <field> (<value> | --at now)")
	}
	target, field, value := arguments[0], arguments[1], arguments[2]
	isStampField := strings.HasSuffix(field, "_at")
	if stampRequested {
		if !isStampField || arguments[3] != "now" {
			return writerRefusal("FRONTMATTER-AT-INVALID", target, "--at accepts only the word now, and only on a field ending in _at")
		}
		value = requestmodel.CanonicalTimestamp(time.Now())
	} else if isStampField {
		if parsed, err := time.Parse(time.RFC3339, value); err != nil || requestmodel.CanonicalTimestamp(parsed) != value {
			return writerRefusal("FRONTMATTER-STAMP-NOT-CANONICAL", target, "a stamp value must be UTC YYYY-MM-DDTHH:MM:SSZ: "+value)
		}
	}
	requestPath, refusal := resolveActiveRequest(executionContext.RepositoryRoot, target)
	if refusal != nil {
		return *refusal
	}
	document, contents, err := readRequestDocument(executionContext.RepositoryRoot, requestPath)
	if err != nil {
		return writerRefusal("REQUEST-UNREADABLE", requestPath, err.Error())
	}
	evidence, found := document.FieldValue(field)
	if found && (evidence.ListValues != nil || evidence.NestedValues != nil) {
		return writerRefusal("FRONTMATTER-FIELD-STRUCTURED", requestPath, field+" holds a list or nested block; set writes scalars only")
	}
	if found && isStampField && strings.TrimSpace(evidence.ScalarValue) != "" {
		return writerRefusal("FRONTMATTER-STAMP-EXISTS", requestPath, field+" is already stamped; stamps are append-only")
	}
	if err := document.SetScalar(field, value); err != nil {
		return writerRefusal("FRONTMATTER-SET-INVALID", requestPath, err.Error())
	}
	if refusal := replaceRequestFile(executionContext.RepositoryRoot, requestPath, contents, document.DocumentBytes()); refusal != nil {
		return *refusal
	}
	output := value + "\n"
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &output}
}

// handleRequest owns `req append-section`: it writes one lifecycle section in
// requestmodel.CanonicalSectionOrder, the order `advance` checks, so the writer
// can never place a section where advance would refuse it.
func handleRequest(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	usage := "usage: req append-section REQ-N --section <name> --from <file>"
	if len(arguments) < 2 || arguments[0] != "append-section" || !requestIDPattern.MatchString(arguments[1]) {
		return usageResult(CommandRequest, usage)
	}
	sectionName, sourcePath := "", ""
	for index := 2; index < len(arguments); index++ {
		var err error
		switch name, _, _ := strings.Cut(arguments[index], "="); name {
		case "--section":
			sectionName, err = optionValue(arguments, &index, name)
		case "--from":
			sourcePath, err = optionValue(arguments, &index, name)
		default:
			return usageResult(CommandRequest, "unknown option "+arguments[index])
		}
		if err != nil {
			return usageResult(CommandRequest, err.Error())
		}
	}
	if sectionName == "" || sourcePath == "" {
		return usageResult(CommandRequest, usage)
	}
	sourceBytes, err := os.ReadFile(absoluteFromRoot(executionContext.RepositoryRoot, sourcePath))
	if err != nil {
		return usageResult(CommandRequest, err.Error())
	}
	requestID := arguments[1]
	if !slices.Contains(requestmodel.CanonicalSectionOrder, sectionName) {
		return writerRefusal("SECTION-NOT-CANONICAL", requestID, sectionName+" is not a canonical REQ section")
	}
	requestPath, refusal := resolveActiveRequest(executionContext.RepositoryRoot, requestID)
	if refusal != nil {
		return *refusal
	}
	document, contents, err := readRequestDocument(executionContext.RepositoryRoot, requestPath)
	if err != nil {
		return writerRefusal("REQUEST-UNREADABLE", requestPath, err.Error())
	}
	sectionBody := trimBlankLines(string(sourceBytes))
	body := document.BodyBytes()
	laterSections := requestmodel.SectionsAfter(sectionName)
	existing := []requestmodel.VisibleSection{}
	insertAt := -1
	// Only column-0 headings are sections, matching advance's advanceSections;
	// fenced, commented or indented headings are the user's text.
	for _, section := range requestmodel.VisibleSections(body) {
		if section.HeadingIndent != 0 {
			continue
		}
		if section.Name == sectionName {
			existing = append(existing, section)
		} else if insertAt < 0 && slices.Contains(laterSections, section.Name) {
			insertAt = section.Start
		}
	}
	success := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: new(string)}
	switch len(existing) {
	case 0:
	case 1:
		sectionText := string(body[existing[0].Start:existing[0].End])
		_, existingBody, _ := strings.Cut(sectionText, "\n")
		if trimBlankLines(existingBody) == sectionBody {
			return success
		}
		return writerRefusal("SECTION-CONFLICT", requestPath, "## "+sectionName+" already exists with a different body")
	default:
		return writerRefusal("SECTION-DUPLICATE", requestPath, "## "+sectionName+" appears more than once")
	}
	// New text is LF: requestmodel does not expose the document's line ending.
	insertion := "## " + sectionName + "\n\n" + sectionBody + "\n"
	if insertAt >= 0 {
		insertion += "\n"
	} else {
		insertAt = len(body)
		switch {
		case len(body) == 0 || strings.HasSuffix(string(body), "\n\n"):
		case strings.HasSuffix(string(body), "\n"):
			insertion = "\n" + insertion
		default:
			insertion = "\n\n" + insertion
		}
	}
	if err := document.ReplaceBodySpan(insertAt, insertAt, []byte(insertion)); err != nil {
		return writerRefusal("SECTION-WRITE-FAILED", requestPath, err.Error())
	}
	// An unclosed fence or comment at the end of the file would hide the new
	// heading, and every retry would then append another invisible copy.
	visibleCopies := 0
	for _, section := range requestmodel.VisibleSections(document.BodyBytes()) {
		if section.HeadingIndent == 0 && section.Name == sectionName {
			visibleCopies++
		}
	}
	if visibleCopies != 1 {
		return writerRefusal("SECTION-WRITE-FAILED", requestPath, "the inserted ## "+sectionName+" would not be a visible section")
	}
	if refusal := replaceRequestFile(executionContext.RepositoryRoot, requestPath, contents, document.DocumentBytes()); refusal != nil {
		return *refusal
	}
	return success
}

// resolveActiveRequest maps a REQ id (through requeststate.ResolveTarget) or a
// repository path to one request file in do-work/working/ or do-work/queue/.
// Archived requests are history and are never written here.
func resolveActiveRequest(repositoryRoot, target string) (string, *resultmodel.CommandResult) {
	if !requestIDPattern.MatchString(target) {
		relative, err := filepath.Rel(repositoryRoot, absoluteFromRoot(repositoryRoot, target))
		relative = filepath.ToSlash(relative)
		if err != nil || !(strings.HasPrefix(relative, "do-work/working/") || strings.HasPrefix(relative, "do-work/queue/")) {
			refusal := writerRefusal("REQUEST-NOT-ACTIVE", target, "the file is not under do-work/working/ or do-work/queue/")
			return "", &refusal
		}
		return relative, nil
	}
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		refusal := writerRefusal("REQUEST-NOT-FOUND", target, err.Error())
		return "", &refusal
	}
	requestFile, stateRefusal := requeststate.ResolveTarget(snapshot, target, "")
	if stateRefusal != nil {
		refusal := writerRefusal(stateRefusal.Code, target, stateRefusal.Reason)
		return "", &refusal
	}
	requestPath := filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(requestFile.RelativePath)))
	if requestFile.TreeSection != "working" && requestFile.TreeSection != "queue" {
		refusal := writerRefusal("REQUEST-NOT-ACTIVE", requestPath, target+" is in "+requestFile.TreeSection+", not working or queue")
		return "", &refusal
	}
	return requestPath, nil
}

func readRequestDocument(repositoryRoot, requestPath string) (*requestmodel.RequestDocument, []byte, error) {
	contents, err := os.ReadFile(absoluteFromRoot(repositoryRoot, requestPath))
	if err != nil {
		return nil, nil, err
	}
	document, err := requestmodel.ParseDocument(contents)
	return document, contents, err
}

// replaceRequestFile publishes the edited bytes atomically; it writes nothing
// when the edit changed no byte.
func replaceRequestFile(repositoryRoot, requestPath string, before, after []byte) *resultmodel.CommandResult {
	if string(before) == string(after) {
		return nil
	}
	if err := atomicfile.ReplaceExisting(absoluteFromRoot(repositoryRoot, requestPath), after); err != nil {
		refusal := writerRefusal("REQUEST-WRITE-FAILED", requestPath, err.Error())
		return &refusal
	}
	return nil
}

// writerRefusal is the one refusal shape for both writers: outcome refused
// (exit 1), the request file untouched.
func writerRefusal(code, path, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{
		helperFinding(code, resultmodel.SeverityError, []string{path}, evidence, resultmodel.FixabilityRefused,
			"the request file was left unchanged", nil, nil),
	}}
}

// trimBlankLines drops whitespace-only lines from both ends, so a section body
// compares equal however many blank lines surround it.
func trimBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
