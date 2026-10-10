package publication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const captureExampleREQTemplate = `---
id: REQ-NNN
title: '<title>'
status: pending
created_at: <created-at>
user_request: UR-NNN
domain: general
prime_files: []
tdd: false
maintenance: false
---
# <title>
## What
<what is being requested, in 1-3 sentences>
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed ` + "`prime_files`" + ` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run ` + "`git diff --stat`" + ` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: <original verbatim request>*
`

// handleCaptureFilesExample prints one capture-files manifest plus the UR and REQ payloads it names,
// for a session to split into files and fill in. It never plans, applies, or writes anything.
func handleCaptureFilesExample(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	rawInputPath, usageError := parseCaptureExampleOptions(arguments)
	if usageError != nil {
		return commandFailure(executionContext.RepositoryRoot, OperationCaptureFiles, "PUBLICATION-USAGE", usageError.Error())
	}
	if executionContext.Format == resultmodel.FormatJSON {
		return commandFailure(executionContext.RepositoryRoot, OperationCaptureFiles, "PUBLICATION-USAGE", "--example prints plain text; use --format text")
	}
	verbatimBlock := "> ```\n> <verbatim input>\n> ```"
	var rawInput *PayloadFile
	if rawInputPath != "" {
		rawInput = &PayloadFile{SourcePath: rawInputPath}
		rawBytes, _, readError := readPayload(executionContext.RepositoryRoot, *rawInput)
		if readError != nil {
			return refusalResult(refusedPlan(PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: executionContext.RepositoryRoot}, "CAPTURE-RAW-INPUT-INVALID", readError.Error(), nil, rawInputPath))
		}
		if controlError := validateOutsideBytes(rawBytes); controlError != nil {
			return refusalResult(refusedPlan(PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: executionContext.RepositoryRoot}, "CAPTURE-RAW-INPUT-UNSAFE", controlError.Error(), nil, rawInputPath))
		}
		verbatimBlock = string(containedOutsideBytes(rawBytes, "\n"))
	}

	manifestPath, urPayloadPath, reqPayloadPath := "<payload-dir>/manifest.json", "<payload-dir>/ur-input.md", "<payload-dir>/req.md"
	manifest := Manifest{Operation: OperationCaptureFiles, CommitMessage: "[UR-NNN] captured <title> (1 REQs)", Capture: &CaptureManifest{
		UserRequestID: "UR-NNN",
		UserRequest:   PublishedFile{Path: "do-work/user-requests/UR-NNN/input.md", Payload: PayloadFile{SourcePath: urPayloadPath}},
		RawInput:      rawInput,
		Requests: []CaptureRequest{{ID: "REQ-NNN", UserRequestID: "UR-NNN", ReservationPath: "do-work/.req-reservations/REQ-NNN",
			File: PublishedFile{Path: "do-work/queue/REQ-NNN-<slug>.md", Payload: PayloadFile{SourcePath: reqPayloadPath}}}},
	}}
	var manifestJSON bytes.Buffer
	encoder := json.NewEncoder(&manifestJSON)
	encoder.SetEscapeHTML(false) // keep <placeholders> readable instead of < escapes
	encoder.SetIndent("", "  ")
	if encodeError := encoder.Encode(manifest); encodeError != nil {
		return commandFailure(executionContext.RepositoryRoot, OperationCaptureFiles, "PUBLICATION-USAGE", encodeError.Error())
	}
	urTemplate := "---\nid: UR-NNN\ntitle: '<title>'\ncreated_at: <created-at>\nrequests: [REQ-NNN]\nword_count: <word-count>\n---\n## Full Verbatim Input\n" + verbatimBlock + "\n"

	var output strings.Builder
	output.WriteString("# capture-files example. Write each part below to the file its ==> line names, then replace every placeholder:\n")
	output.WriteString("# UR-NNN, REQ-NNN (the ids capture reserved), <payload-dir> (a private temporary directory), <slug>, <title>,\n")
	output.WriteString("# <created-at> (whole-second UTC, e.g. 2026-01-26T10:00:00Z), <word-count>, and the angle-bracket body text.\n")
	if rawInputPath == "" {
		output.WriteString("# Rerun with --raw-input <file> to get the Full Verbatim Input block derived from the raw input bytes.\n")
	}
	output.WriteString("# Then check it: capture-files --manifest " + manifestPath + " --dry-run\n")
	// The ==> delimiter lines are safe to split on: every raw-input line inside the UR part is "> "-prefixed by containment.
	fmt.Fprintf(&output, "==> %s <==\n%s", manifestPath, manifestJSON.String())
	fmt.Fprintf(&output, "==> %s <==\n%s", urPayloadPath, urTemplate)
	fmt.Fprintf(&output, "==> %s <==\n%s", reqPayloadPath, captureExampleREQTemplate)
	text := output.String()
	return resultmodel.CommandResult{Command: string(OperationCaptureFiles), Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: executionContext.RepositoryRoot, ExactTextOutput: &text}
}

// parseCaptureExampleOptions accepts only --example and an optional --raw-input <file>.
func parseCaptureExampleOptions(arguments []string) (string, error) {
	rawInputPath, exampleSeen := "", false
	for index := 0; index < len(arguments); index++ {
		switch arguments[index] {
		case "--example":
			if exampleSeen {
				return "", fmt.Errorf("--example may be supplied once")
			}
			exampleSeen = true
		case "--raw-input":
			index++
			if index >= len(arguments) || arguments[index] == "" || rawInputPath != "" {
				return "", fmt.Errorf("--raw-input requires one file path, supplied once")
			}
			rawInputPath = arguments[index]
		default:
			return "", fmt.Errorf("--example accepts only --raw-input <file>, not %q", arguments[index])
		}
	}
	return rawInputPath, nil
}
