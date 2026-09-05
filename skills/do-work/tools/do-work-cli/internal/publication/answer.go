package publication

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
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

// The separator BuildAnswerPlan appends to a resolved question line, and the two disposition
// prefixes it can place immediately after that separator. Writer and reader share these exact
// literals so the line this package composes and the line its terminal-status verdict reads
// cannot drift apart.
const (
	resolvedQuestionSeparator = " → "
	confirmedLabelPrefix      = "Confirmed: "
	discardedLabelPrefix      = "Discarded: "
)

// The delimiters a `## Blocked`, `## Reports`, or `## Implementation` evidence payload is
// composed from. Every byte of those payloads is written by the command's caller, so the only
// positions a reader may attribute are the ones these delimiters open: a history entry's own
// list bullet, the field after the bracketed date that bullet may carry, and the field after
// the single em-dash separator an entry uses to close its subject.
const (
	historyFieldSeparator = " — "
	// bulletListMarkers is CommonMark's own set of bullet list markers, taken wholesale rather
	// than narrowed to the one spelling this project's writers happen to use. A marker must be
	// followed by a space or a tab, which is what separates the "*" of a bullet from the "**"
	// of emphasis.
	bulletListMarkers = "-+*"
	// indentedCodeColumns is CommonMark's threshold for an indented code block. Anything at or
	// beyond it is a picture of a record rather than a record.
	indentedCodeColumns = 4
	// atxHeadingMaximumHashes and blockRunMinimumLength are the two lengths CommonMark fixes for
	// the constructs that end a section: an ATX heading of one to six "#", and a run of three or
	// more of one punctuation mark, which is how every fence, thematic break and frontmatter
	// delimiter is built.
	atxHeadingMaximumHashes = 6
	blockRunMinimumLength   = 3
)

// The markers a terminal stakeholder disposition must carry, at the position the writer places
// them. These are spellings, not a condition, and that is deliberate: the set is closed by the
// writers, not by this reader. `actions/stakeholder-answers.md` Step 5, `actions/clarify.md`
// Step 4 and `internal/doctor`'s HOLLOW-COMPLETION exception all write and read the first
// Implementation spelling; the second is the earlier form this package's fixtures carry. Adding
// a spelling here without adding it there splits two readers of one marker, so a new spelling
// belongs in the writers first. A missing spelling refuses a genuine terminal disposition —
// visible, pre-mutation, correctable — where a missing *position* would complete and archive a
// REQ on the caller's own narrative, so this list may only ever be too short, never too
// permissive. Positions are stated as conditions below for exactly the opposite reason.
const blockedResolutionMarker = "resolved"

var implementationNoCodeMarkers = []string{"no changes needed", "no code changes"}

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
		if label := answeredSummaryDispositionLabel(question); label != "" {
			return refusedPlan(plan, "ANSWER-SUMMARY-INVALID", fmt.Sprintf("an answered summary must not open with the disposition label %q: at that position no reader can tell it from the label this writer places there, so carry that meaning in the outcome field instead", label), []string{record.RequestID}, requestPath)
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
			answerLabel = confirmedLabelPrefix + answerLabel
		}
		if question.Outcome == "discarded" {
			answerLabel = discardedLabelPrefix + answerLabel
		}
		resolvedLine := append([]byte("- [x] "), originalLine[len("- [ ] "):]...)
		resolvedLine = append(resolvedLine, []byte(resolvedQuestionSeparator+answerLabel)...)
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
			if allSubmittedDiscarded && allResolvedQuestionsMatch(questionsAfterEdits, discardedLabelPrefix) {
				status = "cancelled"
			}
			if allSubmittedConfirmed && record.BuilderDecidedValue == "true" && allResolvedQuestionsMatch(questionsAfterEdits, confirmedLabelPrefix) {
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
			if !blockedHistoryRecordsResolution(blockedHistory) || !implementationRecordsNoCodeCompletion(implementation) {
				return refusedPlan(plan, "ANSWER-STAKEHOLDER-EVIDENCE-INVALID", "terminal evidence must carry a resolved Blocked history entry and an Implementation no-change note, each at the position its writer places it: "+terminalEvidencePositionEvidence(blockedHistory, implementation), []string{record.RequestID}, requestPath)
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
			if reportsError != nil || !reportsHistoryNamesReportPath(reportsHistory, reportPath, requestPath) {
				reason := "Reports history must name the fresh report path as the path field of one history entry, not merely somewhere in its text"
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
		reason := "close_user_request must match the whole repository record set"
		if answer.CloseUserRequest && !terminal {
			reason += "; " + nonTerminalDispositionEvidence(resultStatus, questionsAfterEdits)
		}
		return refusedPlan(plan, "ANSWER-UR-CLOSURE-MISMATCH", reason, []string{record.RequestID, record.UserRequestID}, requestPath)
	}
	if terminal {
		destinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", filepath.Base(requestPath)))
		if projectedURClosure {
			destinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", record.UserRequestID, filepath.Base(requestPath)))
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
	// A declared archive path is the caller's own computed disposition, so it is checked against
	// every derived verdict rather than only against a terminal one. Validating it inside the
	// terminal branch alone silently dropped the disagreement that matters most: a REQ whose
	// questions the caller read as all discarded going back to the queue as pending.
	if answer.ArchivePath != "" {
		declaredArchive, archiveError := containedPath(answer.ArchivePath)
		if archiveError != nil || declaredArchive != destinationPath {
			reason := "declared archive path does not match the command-derived disposition"
			if !terminal {
				reason += "; " + nonTerminalDispositionEvidence(resultStatus, questionsAfterEdits)
			}
			return refusedPlan(plan, "ANSWER-ARCHIVE-PATH-MISMATCH", reason, []string{record.RequestID}, answer.ArchivePath)
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

// answeredSummaryDispositionLabel returns the disposition label an "answered" summary would
// open with, or "" when the summary is safe to write.
//
// This is the write-side half of the disposition contract, and the reader cannot stand without
// it. An answered summary is written at the one position a disposition label can occupy, and
// the published line records nothing about which of the two supplied those bytes — so anchoring
// the reader to that position cannot tell a label from a summary opening with one. Refusing the
// collision here is what makes the position mean something: a label at the writer's position
// can only have come from the writer's own outcome field. This check reads the same constants
// the reader prefix-tests, so the two cannot drift.
func answeredSummaryDispositionLabel(question QuestionAnswer) string {
	if question.Outcome != "answered" {
		// Any other outcome puts the writer's own label at the position first, so the summary
		// behind it cannot change what a reader attributes to the line.
		return ""
	}
	for _, labelPrefix := range []string{confirmedLabelPrefix, discardedLabelPrefix} {
		if strings.HasPrefix(question.Summary, labelPrefix) {
			return labelPrefix
		}
	}
	return ""
}

// lineIndentColumns returns the width in columns of a line's leading whitespace and the byte
// offset where its content begins, with a tab advancing to the next four-column stop as
// Markdown counts it. Every caller measures before it trims: a trim run first destroys the very
// whitespace that constitutes the structure, which is the trap REQ-460 recorded and which is
// how an indented code block read as a history entry here.
func lineIndentColumns(line []byte) (int, int) {
	columns, offset := 0, 0
	for offset < len(line) {
		switch line[offset] {
		case ' ':
			columns++
		case '\t':
			columns += indentedCodeColumns - columns%indentedCodeColumns
		default:
			return columns, offset
		}
		offset++
	}
	return columns, offset
}

// blockConstructOpensLine reports whether a line opens a Markdown block construct, which ends
// the region of the payload that records anything.
//
// This is the condition, not a list of spellings. Markdown builds every construct that can end
// a section out of two shapes at a line start: an ATX heading, which is one to six "#" followed
// by a space, a tab, or nothing; and a run of three or more of one ASCII punctuation mark,
// which is how backtick fences, tilde fences, thematic breaks, frontmatter delimiters, setext
// underlines and a dialect's own container fences are all spelled. Taking the whole ASCII
// punctuation class from isMarkdownBlockPunctuation — the same wholesale move that predicate
// already makes for summaries — is what stops a fence character this file has never heard of
// from silently reopening the hole: a new spelling is still punctuation, so it still ends the
// region.
//
// Ending the region rather than toggling a fence pair is deliberate. A toggle has to find a
// matching close, and an unclosed or mismatched fence then leaves everything after it readable
// again — the failure this treatment cannot have. Ending instead means a payload that opens any
// such construct records nothing after it, which refuses rather than accepts.
func blockConstructOpensLine(line []byte) bool {
	_, offset := lineIndentColumns(line)
	content := line[offset:]
	if len(content) == 0 {
		return false
	}
	if content[0] == '#' {
		hashes := 0
		for hashes < len(content) && content[hashes] == '#' {
			hashes++
		}
		if hashes <= atxHeadingMaximumHashes && (hashes == len(content) || content[hashes] == ' ' || content[hashes] == '\t') {
			return true
		}
	}
	if isMarkdownBlockPunctuation(content[0]) {
		run := 0
		for run < len(content) && content[run] == content[0] {
			run++
		}
		if run >= blockRunMinimumLength {
			return true
		}
	}
	return false
}

// activeEvidenceLines returns the lines of a caller-supplied evidence payload that belong to the
// section the payload is published into. The optional heading the payload opens with is dropped;
// everything from the first line that opens another block construct onwards belongs to something
// other than this section's records, so none of it may supply a marker that decides a lifecycle
// write. Whether a surviving line records anything is recordLineContent's question, not this
// one's.
func activeEvidenceLines(evidence []byte, heading string) [][]byte {
	var activeLines [][]byte
	for lineIndex, rawLine := range bytes.Split(evidence, []byte("\n")) {
		line := bytes.TrimRight(bytes.TrimSuffix(rawLine, []byte("\r")), " \t")
		if lineIndex == 0 && bytes.Equal(bytes.TrimSpace(line), []byte(heading)) {
			continue
		}
		if blockConstructOpensLine(line) {
			break
		}
		activeLines = append(activeLines, line)
	}
	return activeLines
}

// recordLineContent returns a line's content when the line can record something. Four or more
// columns of leading whitespace is an indented code block, so its bytes are a picture of a
// record rather than one — the single place that condition is decided, so that both the history
// readers and the Implementation reader answer it the same way.
func recordLineContent(line []byte) ([]byte, bool) {
	columns, offset := lineIndentColumns(line)
	if columns >= indentedCodeColumns {
		return nil, false
	}
	return line[offset:], true
}

// historyEntryContent returns the content of a history entry line: the bytes after the list
// bullet every entry opens with, and after the bracketed date field when the entry carries one.
// A line that is not an entry offers no field boundary at all and is reported as such, so a
// sentence of narrative can never stand in for a recorded entry.
//
// A bracketed field is only a date when its "]" is not immediately followed by "(". That one
// byte is what separates "- [2026-09-01] …", whose bracket closes a date, from
// "- [Title](path)", whose bracket opens a link whose destination is the entry's path field.
func historyEntryContent(line []byte) ([]byte, bool) {
	entry, isRecord := recordLineContent(line)
	if !isRecord || len(entry) < 2 || !strings.ContainsRune(bulletListMarkers, rune(entry[0])) {
		return nil, false
	}
	if entry[1] != ' ' && entry[1] != '\t' {
		return nil, false
	}
	content := bytes.TrimLeft(entry[1:], " \t")
	if bytes.HasPrefix(content, []byte("[")) {
		if closingIndex := bytes.IndexByte(content, ']'); closingIndex >= 0 &&
			(closingIndex+1 == len(content) || content[closingIndex+1] != '(') {
			content = bytes.TrimLeft(content[closingIndex+1:], " \t")
		}
	}
	return content, true
}

// historyEntryTrailingField returns the field an entry's own separator opens, and whether that
// position is identifiable. Both sides of the separator are human text — a quoted stakeholder
// name on one side, a disposition on the other — so an entry carrying more than one separator
// has no field a reader can attribute, exactly as a resolved question line does not.
func historyEntryTrailingField(content []byte) ([]byte, bool) {
	separator := []byte(historyFieldSeparator)
	if bytes.Count(content, separator) != 1 {
		return nil, false
	}
	return content[bytes.Index(content, separator)+len(separator):], true
}

// historyEntryPathField returns the path a history entry names in its first field, and whether
// the entry opens with one at all.
//
// The same path is written three ways in this repository — bare, wrapped in a Markdown link
// whose destination is the path, and fenced in backticks — and all three are the same evidence
// wearing different skins, so all three are read. What matters is that each skin carries its own
// terminator: a link destination ends at its ")", a backticked path at its closing backtick, a
// bare path at the entry's field separator or the end of the line. Returning the whole field
// rather than testing a prefix is what keeps a neighbouring bundle whose path merely starts with
// this one out, in every skin at once.
func historyEntryPathField(content []byte) ([]byte, bool) {
	if bytes.HasPrefix(content, []byte("[")) {
		linkStart := bytes.Index(content, []byte("]("))
		if linkStart < 0 {
			return nil, false
		}
		destination := content[linkStart+len("]("):]
		closingIndex := bytes.IndexByte(destination, ')')
		if closingIndex < 0 {
			return nil, false
		}
		return destination[:closingIndex], true
	}
	if bytes.HasPrefix(content, []byte("`")) {
		closingIndex := bytes.IndexByte(content[1:], '`')
		if closingIndex < 0 {
			return nil, false
		}
		return content[1 : 1+closingIndex], true
	}
	if separatorIndex := bytes.Index(content, []byte(historyFieldSeparator)); separatorIndex >= 0 {
		return content[:separatorIndex], true
	}
	return content, true
}

// markerOpensField reports whether a field opens with a marker as a whole word. Leading emphasis
// is the writer's formatting rather than content, so it is stepped over; a letter or digit
// immediately after the marker means the field opens with a different word that merely starts
// with the same bytes, which is how "no code review yet" passed for the no-code marker.
func markerOpensField(field []byte, marker string) bool {
	candidate := bytes.TrimLeft(field, "*_ \t")
	if len(candidate) < len(marker) || !strings.EqualFold(string(candidate[:len(marker)]), marker) {
		return false
	}
	remainder := candidate[len(marker):]
	if len(remainder) == 0 {
		return true
	}
	next := remainder[0]
	return !(next >= 'a' && next <= 'z' || next >= 'A' && next <= 'Z' || next >= '0' && next <= '9')
}

// blockedHistoryRecordsResolution reports whether a `## Blocked` payload records the resolution
// its caller's terminal disposition claims. The marker is read only where a history entry can
// place one — opening the entry, or opening the field the entry's own separator closes its
// subject with — because the payload is caller-authored prose in which "still not resolved"
// carries the same bytes as a resolution.
func blockedHistoryRecordsResolution(evidence []byte) bool {
	for _, line := range activeEvidenceLines(evidence, "## Blocked") {
		content, isEntry := historyEntryContent(line)
		if !isEntry {
			continue
		}
		if markerOpensField(content, blockedResolutionMarker) {
			return true
		}
		if trailingField, identifiable := historyEntryTrailingField(content); identifiable && markerOpensField(trailingField, blockedResolutionMarker) {
			return true
		}
	}
	return false
}

// implementationRecordsNoCodeCompletion reports whether an `## Implementation` payload opens its
// note with the no-change marker its caller's terminal disposition claims. The marker states the
// whole note, so it must open one — as a paragraph or as the list item some writers use, which
// are the same statement in different skins. The same words inside a sentence describe something
// else, which is what "no code review yet" and "no code changes were needed in the CLI" both are.
func implementationRecordsNoCodeCompletion(evidence []byte) bool {
	for _, line := range activeEvidenceLines(evidence, "## Implementation") {
		statement, isRecord := recordLineContent(line)
		if !isRecord {
			continue
		}
		if entryContent, isEntry := historyEntryContent(line); isEntry {
			statement = entryContent
		}
		for _, marker := range implementationNoCodeMarkers {
			if markerOpensField(statement, marker) {
				return true
			}
		}
	}
	return false
}

// reportsHistoryNamesReportPath reports whether a `## Reports` payload records the fresh report
// as the path field of one history entry, in whichever skin that entry writes it.
func reportsHistoryNamesReportPath(evidence []byte, reportPath, requestPath string) bool {
	for _, line := range activeEvidenceLines(evidence, "## Reports") {
		content, isEntry := historyEntryContent(line)
		if !isEntry {
			continue
		}
		namedPath, hasPath := historyEntryPathField(content)
		if hasPath && repositoryRelativeEvidencePath(string(namedPath), requestPath) == reportPath {
			return true
		}
	}
	return false
}

// repositoryRelativeEvidencePath resolves the path a history entry names into the
// repository-relative form the manifest declares. A link inside a request document is written
// relative to that document — this repository's own archive does exactly that — so the request's
// own directory is the only frame in which such an entry can be read.
func repositoryRelativeEvidencePath(namedPath, requestPath string) string {
	namedPath = strings.TrimSpace(namedPath)
	if namedPath == "" {
		return ""
	}
	if strings.HasPrefix(namedPath, "./") || strings.HasPrefix(namedPath, "../") {
		namedPath = path.Join(path.Dir(requestPath), namedPath)
	}
	// No escape guard sits here on purpose. reportPath has already been through containedPath,
	// so it is repository-relative and can never equal a path that climbs out of the repository
	// or starts at the filesystem root — those resolve to something this comparison rejects on
	// its own, and a guard that cannot change a verdict is a branch no test can pin.
	return path.Clean(namedPath)
}

// terminalEvidencePositionEvidence names which half of the terminal evidence failed and where
// the marker it wanted has to sit. A caller cannot correct what the refusal does not name, and
// the two payloads fail for different reasons often enough that one shared sentence sends the
// caller to the wrong file.
func terminalEvidencePositionEvidence(blockedHistory, implementation []byte) string {
	var missing []string
	if !blockedHistoryRecordsResolution(blockedHistory) {
		missing = append(missing, fmt.Sprintf("no Blocked history entry opens with %q, or carries it in the field after its %q separator",
			blockedResolutionMarker, historyFieldSeparator))
	}
	if !implementationRecordsNoCodeCompletion(implementation) {
		missing = append(missing, "no Implementation paragraph or list item opens with "+quotedAlternatives(implementationNoCodeMarkers))
	}
	return strings.Join(missing, "; ")
}

// quotedAlternatives renders a set of accepted spellings as prose a caller can act on. Printing
// the slice with %q hands them Go syntax instead, which reads as a bug in the tool rather than
// as the two spellings they may choose between.
func quotedAlternatives(alternatives []string) string {
	quoted := make([]string, 0, len(alternatives))
	for _, alternative := range alternatives {
		quoted = append(quoted, fmt.Sprintf("%q", alternative))
	}
	return strings.Join(quoted, " or ")
}

// resolvedQuestionDisposition returns the text a resolved question line carries at the one
// position BuildAnswerPlan can place a disposition, and whether that position is identifiable
// on this line at all.
//
// Every resolved line is composed as `- [x] <question text> → [Confirmed: |Discarded: ]<summary>`,
// so a disposition begins immediately after the single resolvedQuestionSeparator the writer
// appended. Both fields around that separator are human text — the question and the answer
// summary — and either may contain the separator itself, so on a line carrying more than one
// separator no reader can tell which one the writer wrote. That line has no identifiable
// disposition and is not read as carrying one.
//
// What the position achieves and what it does not: it stops text elsewhere on the line from
// supplying a disposition, and it makes an ambiguous line fail closed. It cannot tell a real
// label from an answer summary that opens with those same bytes, because both occupy this one
// position and the document distinguishes them nowhere. That case is closed on the write side,
// by answeredSummaryDispositionLabel; without that refusal a user's own answer text could still
// decide a terminal status here.
func resolvedQuestionDisposition(line []byte) ([]byte, bool) {
	separator := []byte(resolvedQuestionSeparator)
	if bytes.Count(line, separator) != 1 {
		return nil, false
	}
	return line[bytes.Index(line, separator)+len(separator):], true
}

// resolvedQuestionLines returns every resolved question line in a question section, with the
// trailing carriage return of a CRLF document trimmed so a line reads identically in both
// formats. One iteration serves the verdict and its evidence, so they cannot disagree about
// which lines are being judged.
func resolvedQuestionLines(body []byte) [][]byte {
	var resolvedLines [][]byte
	for _, rawLine := range bytes.Split(body, []byte("\n")) {
		line := bytes.TrimSuffix(rawLine, []byte("\r"))
		if bytes.HasPrefix(line, []byte("- [x] ")) {
			resolvedLines = append(resolvedLines, line)
		}
	}
	return resolvedLines
}

// allResolvedQuestionsMatch reports whether every resolved question in the section carries the
// required disposition prefix at the writer's position. Its callers turn this verdict straight
// into a terminal status and an archive move, so a resolved line whose disposition is not
// identifiable fails the check: a section holding one lands on the non-terminal status instead
// of being cancelled or completed on evidence that cannot be attributed to the writer.
func allResolvedQuestionsMatch(body []byte, labelPrefix string) bool {
	resolvedLines := resolvedQuestionLines(body)
	for _, line := range resolvedLines {
		disposition, identifiable := resolvedQuestionDisposition(line)
		if !identifiable || !bytes.HasPrefix(disposition, []byte(labelPrefix)) {
			return false
		}
	}
	return len(resolvedLines) > 0
}

// nonTerminalDispositionEvidence explains a derived non-terminal status to a caller that
// declared a terminal one, naming the resolved lines whose disposition cannot be attributed to
// the writer. Those lines are why an otherwise uniform section fails the check, and a caller
// cannot correct what the refusal does not name.
func nonTerminalDispositionEvidence(status string, questionSection []byte) string {
	evidence := fmt.Sprintf("the command derived the non-terminal status %q", status)
	var unattributableLines []string
	for _, line := range resolvedQuestionLines(questionSection) {
		if _, identifiable := resolvedQuestionDisposition(line); !identifiable {
			unattributableLines = append(unattributableLines, string(line))
		}
	}
	if len(unattributableLines) == 0 {
		return evidence
	}
	return fmt.Sprintf("%s; %d resolved question line(s) carry no disposition at the writer's position, first %q",
		evidence, len(unattributableLines), unattributableLines[0])
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
