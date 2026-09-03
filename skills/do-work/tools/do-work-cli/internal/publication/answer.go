package publication

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

var openQuestionPattern = regexp.MustCompile(`(?m)^- \[ \] `)

// orderedListMarkerPattern matches the one Markdown block opener that starts with a character
// prose also starts with: a digit run followed by "." or ")". Every other block construct
// opens with block-significant leading whitespace or with an ASCII punctuation mark, which
// summaryRequiresContainment tests directly.
var orderedListMarkerPattern = regexp.MustCompile(`^[0-9]+[.)]`)

// summaryRequiresContainment reports whether a one-line answer summary must be carried as a
// file-backed raw payload under canonical containment instead of being written inline into the
// request document.
//
// The rule is the condition, not a list of examples: a summary may be inlined only when no
// Markdown reader can take it for the document's own delimiters or structure. Markdown builds
// every block construct — headings, setext underlines, thematic breaks, code fences, block
// quotes, bullet and ordered list markers, HTML blocks, link reference and footnote
// definitions, tables, frontmatter fences, and whatever a dialect adds next — out of exactly
// three ingredients at a line start: block-significant leading whitespace, an ASCII
// punctuation mark, or a digit run forming an ordered-list marker. Testing those three
// ingredients catches future syntax built from them with no example list to maintain.
//
// Doubt resolves toward containment. Containing a summary that did not need it costs one
// file-backed payload and loses no bytes; inlining one that did writes an unescaped delimiter
// into the document carrying it.
func summaryRequiresContainment(summary string) bool {
	if strings.TrimSpace(summary) == "" {
		return false
	}
	firstByte := summary[0]
	if firstByte == ' ' || firstByte == '\t' {
		return true
	}
	if isMarkdownBlockPunctuation(firstByte) {
		return true
	}
	return orderedListMarkerPattern.MatchString(summary)
}

// isMarkdownBlockPunctuation reports whether a leading byte is one of the ASCII punctuation
// marks Markdown can build a block opener from. The ranges are the whole ASCII punctuation
// block — CommonMark's own definition of ASCII punctuation — taken wholesale rather than
// narrowed to the marks today's block syntax happens to use: the narrower set would be an
// enumeration to revisit every time a dialect adds a construct, which is the defect this
// predicate replaced. A leading letter, digit, or non-ASCII rune reaches no block construct
// except the ordered-list marker its caller tests.
func isMarkdownBlockPunctuation(value byte) bool {
	return value >= '!' && value <= '/' ||
		value >= ':' && value <= '@' ||
		value >= '[' && value <= '`' ||
		value >= '{' && value <= '~'
}

func BuildAnswerPlan(repositoryRoot string, manifest Manifest, answerTime time.Time) PublicationPlan {
	plan := PublicationPlan{Operation: OperationAnswer, RepositoryRoot: repositoryRoot, CommitMessage: manifest.CommitMessage}
	answer := manifest.Answer
	if answer == nil || len(answer.Answers) == 0 {
		return refusedPlan(plan, "ANSWER-MANIFEST-MISSING", "answer body and at least one answer are required", nil)
	}
	requestPath, pathError := containedPath(answer.RequestPath)
	if pathError != nil || !strings.HasPrefix(requestPath, "do-work/") {
		return refusedPlan(plan, "ANSWER-PATH-UNSAFE", "request path must be contained in do-work", nil, answer.RequestPath)
	}
	if strings.HasPrefix(requestPath, "do-work/archive/") {
		return refusedPlan(plan, "ANSWER-ARCHIVED-READ-ONLY", "stale archived replies cannot rewrite their original REQ", nil, requestPath)
	}
	snapshot, discoveryError := repositorymodel.DiscoverRepository(repositoryRoot)
	if discoveryError != nil {
		return refusedPlan(plan, "ANSWER-DISCOVERY-FAILED", discoveryError.Error(), nil, requestPath)
	}
	requestBytes, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(requestPath)))
	if readError != nil {
		return refusedPlan(plan, "ANSWER-REQUEST-MISSING", readError.Error(), nil, requestPath)
	}
	document, parseError := requestmodel.ParseDocument(requestBytes)
	if parseError != nil {
		return refusedPlan(plan, "ANSWER-REQUEST-INVALID", parseError.Error(), nil, requestPath)
	}
	record := document.TypedRecord()
	matchingRecords := 0
	for _, requestFile := range snapshot.RequestFiles {
		discoveredPath := filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(requestFile.RelativePath)))
		if discoveredPath == requestPath && requestFile.TypedRecord.RequestID == record.RequestID {
			matchingRecords++
		}
	}
	if matchingRecords != 1 {
		return refusedPlan(plan, "ANSWER-REQUEST-IDENTITY-MISMATCH", "exact path and frontmatter id must resolve to one repository record", []string{record.RequestID}, requestPath)
	}
	if answer.ExpectedStatus == "" || record.RequestStatus != answer.ExpectedStatus {
		return refusedPlan(plan, "ANSWER-SNAPSHOT-STALE", fmt.Sprintf("expected status %q, found %q", answer.ExpectedStatus, record.RequestStatus), []string{record.RequestID}, requestPath)
	}
	if answerTime.IsZero() {
		answerTime = time.Now().UTC()
	}
	date := answerTime.UTC().Format("2006-01-02")
	canonicalTimestamp := requestmodel.CanonicalTimestamp(answerTime)
	lineEnding := "\n"
	if bytes.Contains(requestBytes, []byte("\r\n")) {
		lineEnding = "\r\n"
	}
	var answerNotes []byte
	seenIdentities := map[string]bool{}
	allSubmittedDiscarded, allSubmittedConfirmed := true, true
	for _, question := range answer.Answers {
		identity := question.QuestionID
		if identity == "" {
			identity = question.ExpectedLine
		}
		if identity == "" || seenIdentities[identity] {
			return refusedPlan(plan, "ANSWER-QUESTION-IDENTITY-INVALID", "each answer needs one unique Q-ID or exact expected line", []string{record.RequestID}, requestPath)
		}
		seenIdentities[identity] = true
		if question.Outcome != "answered" && question.Outcome != "confirmed" && question.Outcome != "discarded" {
			return refusedPlan(plan, "ANSWER-OUTCOME-INVALID", "answer outcome must be answered, confirmed, or discarded", []string{record.RequestID}, requestPath)
		}
		allSubmittedDiscarded = allSubmittedDiscarded && question.Outcome == "discarded"
		allSubmittedConfirmed = allSubmittedConfirmed && question.Outcome == "confirmed"
		if strings.TrimSpace(question.Summary) == "" || strings.ContainsAny(question.Summary, "\r\n") {
			return refusedPlan(plan, "ANSWER-SUMMARY-INVALID", "each answer requires one non-empty summary line", []string{record.RequestID}, requestPath)
		}
		if controlError := validateOutsideBytes([]byte(question.Summary)); controlError != nil {
			return refusedPlan(plan, "ANSWER-TEXT-UNSAFE", controlError.Error(), []string{record.RequestID}, requestPath)
		}
		var rawAnswerBytes []byte
		if question.RawAnswer != nil {
			rawBytes, _, rawError := readPayload(repositoryRoot, *question.RawAnswer)
			if rawError != nil {
				return refusedPlan(plan, "ANSWER-RAW-PAYLOAD-INVALID", rawError.Error(), []string{record.RequestID}, question.RawAnswer.SourcePath)
			}
			if controlError := validateOutsideBytes(rawBytes); controlError != nil {
				return refusedPlan(plan, "ANSWER-TEXT-UNSAFE", controlError.Error(), []string{record.RequestID}, question.RawAnswer.SourcePath)
			}
			rawAnswerBytes = rawBytes
		}
		containedSummary := summaryRequiresContainment(question.Summary)
		if containedSummary && question.RawAnswer == nil {
			return refusedPlan(plan, "ANSWER-RAW-PAYLOAD-REQUIRED", "delimiter-shaped answers require an exact file-backed raw payload", []string{record.RequestID}, requestPath)
		}
		if containedSummary && !bytes.Equal(rawAnswerBytes, []byte(question.Summary)) {
			return refusedPlan(plan, "ANSWER-RAW-PAYLOAD-MISMATCH", "delimiter-shaped summary must byte-match its raw payload", []string{record.RequestID}, question.RawAnswer.SourcePath)
		}
		body := document.BodyBytes()
		lineStart, lineEnd, identityError := findQuestionLine(body, question)
		if identityError != nil && question.InsertQuestion && question.ExpectedLine != "" {
			insertedLine := question.ExpectedLine
			if !strings.HasPrefix(insertedLine, "- [ ] ") {
				return refusedPlan(plan, "ANSWER-QUESTION-IDENTITY-INVALID", "inserted question must be an exact open checkbox line", []string{record.RequestID}, requestPath)
			}
			bodyLength := len(body)
			prefix := ""
			if bodyLength > 0 && !bytes.HasSuffix(body, []byte("\n")) {
				prefix = lineEnding
			}
			if insertError := document.ReplaceBodySpan(bodyLength, bodyLength, []byte(prefix+insertedLine+lineEnding)); insertError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", insertError.Error(), []string{record.RequestID}, requestPath)
			}
			body = document.BodyBytes()
			lineStart, lineEnd, identityError = findQuestionLine(body, question)
		}
		if identityError != nil {
			return refusedPlan(plan, "ANSWER-QUESTION-NOT-UNIQUE", identityError.Error(), []string{record.RequestID, question.QuestionID}, requestPath)
		}
		originalLine := body[lineStart:lineEnd]
		if !bytes.HasPrefix(originalLine, []byte("- [ ] ")) {
			return refusedPlan(plan, "ANSWER-QUESTION-ALREADY-RESOLVED", "matched question is not open", []string{record.RequestID, question.QuestionID}, requestPath)
		}
		answerLabel := question.Summary
		if containedSummary {
			answerLabel = "See contained answer note"
		}
		if question.Outcome == "confirmed" {
			answerLabel = "Confirmed: " + answerLabel
		}
		if question.Outcome == "discarded" {
			answerLabel = "Discarded: " + answerLabel
		}
		resolvedLine := append([]byte("- [x] "), originalLine[len("- [ ] "):]...)
		resolvedLine = append(resolvedLine, []byte(" → "+answerLabel)...)
		if replaceError := document.ReplaceBodySpan(lineStart, lineEnd, resolvedLine); replaceError != nil {
			return refusedPlan(plan, "ANSWER-EDIT-FAILED", replaceError.Error(), []string{record.RequestID}, requestPath)
		}
		answerNotes = append(answerNotes, []byte("- "+date+" "+identity+": "+answerLabel+lineEnding)...)
		if question.RawAnswer != nil {
			answerNotes = append(answerNotes, containedOutsideBytes(rawAnswerBytes, lineEnding)...)
			answerNotes = append(answerNotes, []byte(lineEnding)...)
		}
	}
	if len(answerNotes) > 0 {
		if replaceError := appendAnswerNotes(document, answerNotes, lineEnding); replaceError != nil {
			return refusedPlan(plan, "ANSWER-EDIT-FAILED", replaceError.Error(), []string{record.RequestID}, requestPath)
		}
	}
	questionsAfterEdits := questionSectionBytes(document.BodyBytes())
	remainingOpen := openQuestionPattern.Match(questionsAfterEdits)
	destinationPath := requestPath
	resultStatus := record.RequestStatus
	switch answer.Mode {
	case "clarify":
		status := "pending-answers"
		if !remainingOpen {
			status = "pending"
			if allSubmittedDiscarded && allResolvedQuestionsMatch(questionsAfterEdits, "→ Discarded:") {
				status = "cancelled"
			}
			if allSubmittedConfirmed && record.BuilderDecidedValue == "true" && allResolvedQuestionsMatch(questionsAfterEdits, "→ Confirmed:") {
				status = "completed"
			}
		}
		if setError := document.SetScalar("status", status); setError != nil {
			return refusedPlan(plan, "ANSWER-EDIT-FAILED", setError.Error(), []string{record.RequestID}, requestPath)
		}
		resultStatus = status
		if status == "completed" || status == "cancelled" {
			if timestampError := document.SetScalar("completed_at", canonicalTimestamp); timestampError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", timestampError.Error(), []string{record.RequestID}, requestPath)
			}
		} else if status == "pending" {
			if timestampError := document.SetScalar("status_changed_at", canonicalTimestamp); timestampError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", timestampError.Error(), []string{record.RequestID}, requestPath)
			}
		}
	case "stakeholder":
		status := "blocked"
		if !remainingOpen {
			status = "completed"
		}
		if setError := document.SetScalar("status", status); setError != nil {
			return refusedPlan(plan, "ANSWER-EDIT-FAILED", setError.Error(), []string{record.RequestID}, requestPath)
		}
		resultStatus = status
		if status == "completed" {
			if answer.StakeholderTerminal == nil {
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-TERMINAL-EVIDENCE-MISSING", "terminal stakeholder disposition requires Blocked history and Implementation evidence payloads", []string{record.RequestID}, requestPath)
			}
			blockedHistory, _, blockedError := readPayload(repositoryRoot, answer.StakeholderTerminal.BlockedHistory)
			implementation, _, implementationError := readPayload(repositoryRoot, answer.StakeholderTerminal.Implementation)
			if blockedError != nil || implementationError != nil {
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-EVIDENCE-INVALID", firstError(blockedError, implementationError).Error(), []string{record.RequestID}, requestPath)
			}
			if !bytes.Contains(bytes.ToLower(blockedHistory), []byte("resolved")) || !bytes.Contains(bytes.ToLower(implementation), []byte("no code")) {
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-EVIDENCE-INVALID", "terminal evidence must carry resolved Blocked history and an Implementation no-code marker", []string{record.RequestID}, requestPath)
			}
			if appendError := appendSectionEvidence(document, "## Blocked", blockedHistory, lineEnding); appendError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", appendError.Error(), []string{record.RequestID}, requestPath)
			}
			if appendError := appendSectionEvidence(document, "## Implementation", implementation, lineEnding); appendError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", appendError.Error(), []string{record.RequestID}, requestPath)
			}
			if timestampError := document.SetScalar("completed_at", canonicalTimestamp); timestampError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", timestampError.Error(), []string{record.RequestID}, requestPath)
			}
			_ = document.DeleteField("blocked_by")
			_ = document.DeleteField("blocked_at")
		} else {
			if answer.Report == nil || answer.StakeholderReport == nil {
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-REPORT-EVIDENCE-MISSING", "partial stakeholder disposition requires the fresh report and its linkage/history evidence", []string{record.RequestID}, requestPath)
			}
			reportPath, reportPathError := containedPath(answer.Report.Path)
			if reportPathError != nil || answer.StakeholderReport.BlockedBy != reportPath {
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-REPORT-LINKAGE-INVALID", "blocked_by must exactly match the fresh report path", []string{record.RequestID}, answer.StakeholderReport.BlockedBy)
			}
			reportsHistory, _, reportsError := readPayload(repositoryRoot, answer.StakeholderReport.ReportsHistory)
			if reportsError != nil || !bytes.Contains(reportsHistory, []byte(reportPath)) {
				reason := "Reports history must name the fresh report path"
				if reportsError != nil {
					reason = reportsError.Error()
				}
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-REPORT-EVIDENCE-INVALID", reason, []string{record.RequestID}, reportPath)
			}
			if setError := document.SetScalar("blocked_by", reportPath); setError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", setError.Error(), []string{record.RequestID}, requestPath)
			}
			if appendError := appendSectionEvidence(document, "## Reports", reportsHistory, lineEnding); appendError != nil {
				return refusedPlan(plan, "ANSWER-EDIT-FAILED", appendError.Error(), []string{record.RequestID}, requestPath)
			}
		}
	case "verify-repair":
		// Verify owns prose judgment; the command owns only the exact question edit.
	default:
		return refusedPlan(plan, "ANSWER-MODE-INVALID", "answer mode must be clarify, stakeholder, or verify-repair", []string{record.RequestID}, requestPath)
	}
	terminal := resultStatus == "completed" || resultStatus == "cancelled"
	projectedURClosure := false
	archivedURFallback := false
	if terminal && record.UserRequestID != "" {
		projectedURClosure = true
		for _, requestFile := range snapshot.RequestFiles {
			if requestFile.TypedRecord.RequestID == record.RequestID || requestFile.TypedRecord.UserRequestID != record.UserRequestID {
				continue
			}
			if !schemanormalization.IsTerminalResolved(requestFile.TypedRecord.RequestStatus) {
				projectedURClosure = false
				break
			}
		}
	}
	if answer.CloseUserRequest != projectedURClosure {
		return refusedPlan(plan, "ANSWER-UR-CLOSURE-MISMATCH", "close_user_request must match the whole repository record set", []string{record.RequestID, record.UserRequestID}, requestPath)
	}
	if terminal {
		destinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", filepath.Base(requestPath)))
		if projectedURClosure {
			destinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", record.UserRequestID, filepath.Base(requestPath)))
		}
		if answer.ArchivePath != "" {
			declaredArchive, archiveError := containedPath(answer.ArchivePath)
			if archiveError != nil || declaredArchive != destinationPath {
				return refusedPlan(plan, "ANSWER-ARCHIVE-PATH-MISMATCH", "declared archive path does not match the command-derived disposition", []string{record.RequestID}, answer.ArchivePath)
			}
		}
		if projectedURClosure {
			derivedUserRequestPath := filepath.ToSlash(filepath.Join("do-work", "user-requests", record.UserRequestID))
			derivedArchiveDirectory := filepath.ToSlash(filepath.Join("do-work", "archive", record.UserRequestID))
			if answer.UserRequestPath != "" && answer.UserRequestPath != derivedUserRequestPath || answer.ArchiveDirectory != "" && answer.ArchiveDirectory != derivedArchiveDirectory {
				return refusedPlan(plan, "ANSWER-UR-CLOSURE-PATH-INVALID", "declared UR closure paths do not match the command-derived UR", []string{record.RequestID, record.UserRequestID})
			}
			answer.UserRequestPath, answer.ArchiveDirectory = derivedUserRequestPath, derivedArchiveDirectory
			activeInfo, activeError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(derivedUserRequestPath)))
			archiveInfo, archiveError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(derivedArchiveDirectory)))
			if os.IsNotExist(activeError) && archiveError == nil && archiveInfo.IsDir() && archiveInfo.Mode()&os.ModeSymlink == 0 {
				archivedURFallback = true
			} else if activeError != nil || !activeInfo.IsDir() || activeInfo.Mode()&os.ModeSymlink != 0 {
				return refusedPlan(plan, "ANSWER-UR-CLOSURE-SOURCE-INVALID", "UR closure needs either the exact active UR or an already archived real UR directory", []string{record.RequestID, record.UserRequestID}, derivedUserRequestPath, derivedArchiveDirectory)
			}
		}
	}
	updatedBytes := document.DocumentBytes()
	if destinationPath == requestPath {
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationReplace, Path: requestPath, ExpectedBytes: requestBytes, Contents: updatedBytes})
	} else {
		if collision := pathExists(repositoryRoot, destinationPath); collision {
			return refusedPlan(plan, "ANSWER-ARCHIVE-COLLISION", "archive destination already exists", []string{record.RequestID}, destinationPath)
		}
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationMove, Path: requestPath, DestinationPath: destinationPath, ExpectedBytes: requestBytes, Contents: updatedBytes})
	}
	if answer.Report != nil {
		if refusal := appendCreate(repositoryRoot, &plan, *answer.Report, "ANSWER-REPORT"); refusal != nil {
			return *refusal
		}
	}
	if len(answer.OverrideCreates) > 0 || len(answer.OverrideFolds) > 0 {
		return refusedPlan(plan, "ANSWER-OVERRIDE-UNSTRUCTURED", "override publication must use override_capture with full capture identity, linkage, membership, reservation, raw, asset, topology, and fold validation", []string{record.RequestID})
	}
	if answer.OverrideCapture != nil {
		overridePlan := BuildCapturePlan(repositoryRoot, Manifest{Operation: OperationCaptureFiles, Capture: answer.OverrideCapture})
		if overridePlan.Refusal != nil {
			return refusedPlan(plan, "ANSWER-OVERRIDE-CAPTURE-"+overridePlan.Refusal.Code, overridePlan.Refusal.Reason, overridePlan.Refusal.IDs, overridePlan.Refusal.Paths...)
		}
		plan.Mutations = append(plan.Mutations, overridePlan.Mutations...)
	}
	if answer.CloseUserRequest && !archivedURFallback {
		closurePlan := appendUserRequestClosure(repositoryRoot, plan, *answer)
		if closurePlan.Refusal != nil {
			return closurePlan
		}
		plan = closurePlan
	}
	plan = finalizePlan(plan)
	if plan.Refusal != nil {
		return plan
	}
	directories, topologyError := planCreatedDirectories(repositoryRoot, plan.TargetPaths)
	if topologyError != nil {
		return refusedPlan(plan, "ANSWER-TOPOLOGY-UNSAFE", topologyError.Error(), []string{record.RequestID})
	}
	plan.CreatedDirectoryPaths = directories
	return plan
}

func findQuestionLine(body []byte, question QuestionAnswer) (int, int, error) {
	lines := bytes.SplitAfter(body, []byte("\n"))
	offset := 0
	matches := [][2]int{}
	for _, rawLine := range lines {
		line := bytes.TrimSuffix(rawLine, []byte("\n"))
		line = bytes.TrimSuffix(line, []byte("\r"))
		match := question.ExpectedLine != "" && bytes.Equal(line, []byte(question.ExpectedLine))
		if question.QuestionID != "" && bytes.Contains(line, []byte("**"+question.QuestionID+"**")) {
			match = true
		}
		if match {
			matches = append(matches, [2]int{offset, offset + len(line)})
		}
		offset += len(rawLine)
	}
	if len(matches) != 1 {
		return 0, 0, fmt.Errorf("question identity matched %d lines", len(matches))
	}
	return matches[0][0], matches[0][1], nil
}

func allResolvedQuestionsMatch(body []byte, marker string) bool {
	found := false
	for _, line := range bytes.Split(body, []byte("\n")) {
		if bytes.HasPrefix(bytes.TrimSuffix(line, []byte("\r")), []byte("- [x] ")) {
			found = true
			if !bytes.Contains(line, []byte(marker)) {
				return false
			}
		}
	}
	return found
}

func questionSectionBytes(body []byte) []byte {
	heading := []byte("## Open Questions")
	start := bytes.Index(body, heading)
	if start < 0 {
		return body
	}
	section := body[start+len(heading):]
	if next := bytes.Index(section, []byte("\n## ")); next >= 0 {
		section = section[:next]
	}
	return section
}

func appendAnswerNotes(document *requestmodel.RequestDocument, notes []byte, lineEnding string) error {
	body := document.BodyBytes()
	heading := []byte("## Answer Notes")
	headingStart := bytes.Index(body, heading)
	if headingStart < 0 {
		block := []byte(lineEnding + lineEnding + "## Answer Notes" + lineEnding + lineEnding)
		block = append(block, notes...)
		return document.ReplaceBodySpan(len(body), len(body), block)
	}
	sectionEnd := len(body)
	afterHeading := headingStart + len(heading)
	if nextHeading := bytes.Index(body[afterHeading:], []byte("\n## ")); nextHeading >= 0 {
		sectionEnd = afterHeading + nextHeading
	}
	prefix := []byte{}
	if sectionEnd > 0 && body[sectionEnd-1] != '\n' {
		prefix = []byte(lineEnding)
	}
	insertion := append(prefix, notes...)
	return document.ReplaceBodySpan(sectionEnd, sectionEnd, insertion)
}

func appendSectionEvidence(document *requestmodel.RequestDocument, heading string, evidence []byte, lineEnding string) error {
	body := document.BodyBytes()
	evidence = bytes.TrimPrefix(evidence, []byte(heading))
	evidence = bytes.TrimLeft(evidence, "\r\n")
	if len(evidence) == 0 {
		return fmt.Errorf("%s evidence is empty", heading)
	}
	headingBytes := []byte(heading)
	headingStart := bytes.Index(body, headingBytes)
	if headingStart < 0 {
		block := []byte{}
		if len(body) > 0 && !bytes.HasSuffix(body, []byte("\n")) {
			block = append(block, []byte(lineEnding)...)
		}
		block = append(block, []byte(lineEnding+heading+lineEnding+lineEnding)...)
		block = append(block, evidence...)
		if !bytes.HasSuffix(block, []byte("\n")) {
			block = append(block, []byte(lineEnding)...)
		}
		return document.ReplaceBodySpan(len(body), len(body), block)
	}
	sectionEnd := len(body)
	afterHeading := headingStart + len(headingBytes)
	if nextHeading := bytes.Index(body[afterHeading:], []byte("\n## ")); nextHeading >= 0 {
		sectionEnd = afterHeading + nextHeading
	}
	insertion := []byte{}
	if sectionEnd > 0 && body[sectionEnd-1] != '\n' {
		insertion = append(insertion, []byte(lineEnding)...)
	}
	insertion = append(insertion, evidence...)
	if !bytes.HasSuffix(insertion, []byte("\n")) {
		insertion = append(insertion, []byte(lineEnding)...)
	}
	return document.ReplaceBodySpan(sectionEnd, sectionEnd, insertion)
}

func appendCreate(repositoryRoot string, plan *PublicationPlan, published PublishedFile, code string) *PublicationPlan {
	path, pathError := containedPath(published.Path)
	if pathError != nil {
		refused := refusedPlan(*plan, code+"-PATH-UNSAFE", pathError.Error(), nil, published.Path)
		return &refused
	}
	if pathExists(repositoryRoot, path) {
		refused := refusedPlan(*plan, code+"-COLLISION", "destination already exists", nil, path)
		return &refused
	}
	contents, mode, payloadError := readPayload(repositoryRoot, published.Payload)
	if payloadError != nil {
		refused := refusedPlan(*plan, code+"-PAYLOAD-INVALID", payloadError.Error(), nil, published.Payload.SourcePath)
		return &refused
	}
	plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationCreate, Path: path, Contents: contents, Mode: selectedMode(published.Mode, mode)})
	return nil
}

func appendUserRequestClosure(repositoryRoot string, plan PublicationPlan, answer AnswerManifest) PublicationPlan {
	sourceRoot, sourceError := containedPath(answer.UserRequestPath)
	archiveRoot, archiveError := containedPath(answer.ArchiveDirectory)
	if sourceError != nil || archiveError != nil || !strings.HasPrefix(sourceRoot, "do-work/user-requests/UR-") || !strings.HasPrefix(archiveRoot, "do-work/archive/UR-") {
		return refusedPlan(plan, "ANSWER-UR-CLOSURE-PATH-INVALID", "closure requires exact active and archive UR directories", nil, answer.UserRequestPath, answer.ArchiveDirectory)
	}
	absoluteSource := filepath.Join(repositoryRoot, filepath.FromSlash(sourceRoot))
	sourceInfo, statError := os.Lstat(absoluteSource)
	if statError != nil || !sourceInfo.IsDir() || sourceInfo.Mode()&os.ModeSymlink != 0 {
		return refusedPlan(plan, "ANSWER-UR-CLOSURE-SOURCE-INVALID", "active UR path is not a real directory", nil, sourceRoot)
	}
	var moves []PlannedMutation
	walkError := filepath.WalkDir(absoluteSource, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == absoluteSource {
			return nil
		}
		info, infoError := entry.Info()
		if infoError != nil {
			return infoError
		}
		if entry.Type()&os.ModeSymlink != 0 || !entry.IsDir() && !info.Mode().IsRegular() {
			return fmt.Errorf("UR closure contains unsafe object %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		relative, relError := filepath.Rel(absoluteSource, path)
		if relError != nil {
			return relError
		}
		destination := filepath.ToSlash(filepath.Join(archiveRoot, relative))
		if pathExists(repositoryRoot, destination) {
			return fmt.Errorf("UR closure destination exists: %s", destination)
		}
		contents, readError := os.ReadFile(path)
		if readError != nil {
			return readError
		}
		moves = append(moves, PlannedMutation{Kind: MutationMove, Path: filepath.ToSlash(filepath.Join(sourceRoot, relative)), DestinationPath: destination, ExpectedBytes: contents, Mode: info.Mode()})
		return nil
	})
	if walkError != nil {
		return refusedPlan(plan, "ANSWER-UR-CLOSURE-INVALID", walkError.Error(), nil, sourceRoot)
	}
	sort.Slice(moves, func(i, j int) bool { return moves[i].Path < moves[j].Path })
	plan.Mutations = append(plan.Mutations, moves...)
	return plan
}

func pathExists(repositoryRoot, path string) bool {
	_, err := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
	return err == nil
}
