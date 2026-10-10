## Review: REQ-688

**Approve, with one one-line fix before release.** `capture-files --example` works. The filled example passes the dry run, and the reference UR example now uses the fence the command accepts.
Route B | range `008cc74c..125c62c2` (merge `125c62c2`)

### What's built
- `do-work-cli capture-files --example [--raw-input <file>]` prints a manifest encoded from the real `Manifest` type, a UR template and a REQ template. Each part follows a `==> <path> <==` line. With `--raw-input`, the UR block is `containedOutsideBytes` of the file. It writes nothing.
- `actions/capture-reference.md` UR example: bare three-backtick fence. The paragraph after it states the real fence rule.
- `actions/capture.md` Step 5 names the command.
- Still to do: the release (requirement 7). It belongs to finalization, after this review.

### Decisions / risks for you
- F1 is the only thing I would fix before release. It costs one line of Go and no test.
- PD-6 keeps the UR and REQ templates as Go text, so they are a second copy of the `capture-reference.md` shapes. The test proves the copy is valid, not that it matches the doc. It can drift and still pass (F6).

### Findings

**Minor:**
- F1 `capture_files_example.go:76-81`: the header says "replace every placeholder", and the UR part repeats the raw input word for word. If the raw input contains `UR-NNN` or `REQ-NNN` (common in captures about do-work itself), a global replace also changes the verbatim block, and the dry run refuses `CAPTURE-RAW-INPUT-NOT-CONTAINED`. I reproduced this in a scratch repo: raw input `rename REQ-NNN files and the UR-NNN folder`, `sed` replacement over all three parts, dry run refused. With a raw input that has no placeholder text, the same steps passed. The test uses a global `strings.Replacer` too, and passes only because its raw input has no placeholder text. — impact-user-visible. **Fix before release.** In `skills/do-work/tools/do-work-cli/internal/publication/capture_files_example.go`, replace lines 79-81:
  ```go
  	if rawInputPath == "" {
  		output.WriteString("# Rerun with --raw-input <file> to get the Full Verbatim Input block derived from the raw input bytes.\n")
  	}
  ```
  with:
  ```go
  	if rawInputPath == "" {
  		output.WriteString("# Rerun with --raw-input <file> to get the Full Verbatim Input block derived from the raw input bytes.\n")
  	} else {
  		output.WriteString("# Leave the \"> \" lines under ## Full Verbatim Input as printed, even where they contain UR-NNN or REQ-NNN: they must match the raw input byte for byte.\n")
  	}
  ```
- F2 `skills/do-work/actions/clarify.md:106` (stale restatement): it says "a code fence longer than the longest backtick run". It does not say "never shorter than three" or "no info string". All three Go writers apply both rules: `containedOutsideBytes` (`publication_manifest.go:107`, used by capture and by `answer.go:238`) and `containedOutsideText` (`requeststate/state_apply.go:1052`, used by abandon). `capture-reference.md:196` now says "It uses clarify.md Step 4's ... form" and then states the stricter rule, so the source it cites and the citation disagree. The builder's reason to leave it ("answer notes are not byte-checked") is inaccurate: `answer` derives the note with the same function. The outcome is the same, because no session writes those notes by hand. — impact-rule-change → report only. Fix for a later REQ that touches clarify.md: replace "open a code fence longer than the longest backtick run anywhere in the text." with "open a code fence one backtick longer than the longest backtick run anywhere in the text, never shorter than three, with no info string." `capture-reference.md:196` could then point at that line and stop restating the rule.
- F3 `skills/do-work/actions/capture.md:136,138` (stale example): the queued-addendum example still uses `> ````text` / `> ````. That is four backticks and an info string for a sample with no backticks. Nothing byte-checks the addendum (PD-5), so nothing fails. But it is the same shape that caused the 2026-10-05 refusal, and the new capture-reference text says "no info string". A session copying it into a UR repeats the failure. — impact-negligible → report only. Fix: replace both lines with `> ```` (bare three backticks).

**Nit:**
- F4 `capture_files_example.go:70-72`: the `encodeError` branch cannot fail, because the encoder gets a fixed struct of strings. Three dead lines. — impact-negligible → report only.
- F5 `capture_files_example.go:93-100`: the duplicate-`--example` check (`exampleSeen`) is a refusal nobody asked for. Also, `--raw-input --example` takes `--example` as the file path and refuses with `inspect payload --example` (harmless: the answer is still a refusal). — impact-negligible → report only.
- F6 `capture_files_example.go:13-32`: the REQ template is a second copy of the `capture-reference.md` Simple REQ shape (accepted decision PD-6). The test pins validity, not equality, so the copy can drift while the test stays green. — impact-negligible → report only.
- F7 `capture_files_test.go:327`: the `--format json` refusal check runs inside both subtests, so the same assertion runs twice. Once is enough. — impact-negligible → report only.

### Requirements Checklist
- [x] R1 One valid manifest with every required key; `UR-NNN`, `REQ-NNN`, `do-work/.req-reservations/REQ-NNN`; placeholder payload paths. Encoded from `Manifest`, so the keys cannot drift from the strict decoder.
- [x] R2 UR and REQ payload templates, delimited by `==> <path> <==`. The delimiter is safe because every raw-input line is `> `-prefixed.
- [x] R3 `--raw-input`: block is `containedOutsideBytes(raw, "\n")`; `raw_input.source_path` is the path as given; read through `readPayload` and `validateOutsideBytes`.
- [x] R4 Filled example passes `--dry-run` on the first try, including a triple-backtick run (test plus my by-hand run). F1 is the one input class where a natural fill fails.
- [x] R5 `capture-reference.md:191,193` bare three-backtick fences; the paragraph at `:196` matches `containedOutsideBytes`.
- [x] R6 `capture.md:230` names `capture-files --example --raw-input`; the Step 5 prose stays.
- [ ] R7 Release: pending finalization (not this review's stage).
- [x] Constraints: no disk writes, no reservation marker (checked: scratch repo tree unchanged after 14 CLI calls; test asserts no `do-work/`); no frontmatter field or status; containment rule unchanged (`publication_manifest.go`, `capture_files.go` untouched).

**Doc accuracy (check 2).** `capture-reference.md:196`: "one backtick longer than the longest backtick run ... never shorter than three, with no info string, and every line prefixed `> `" matches `publication_manifest.go:107-136`. "refuses a UR that does not contain it" matches `capture_files.go:77-92` (it applies when `raw_input` is supplied, which `capture.md:236` requires). `capture.md:230` is accurate: plain text only (JSON refuses), writes nothing, `==>` lines name the files.

### Acceptance Testing

**Result: Pass** (stages covered: focused Go tests, CLI by hand in a scratch git repo. Release and installed-consumer stages: not exercised.)
- `go test -count=1 -run 'TestCaptureFilesExampleFilledInPassesDryRun|TestBuildCapturePlanAcceptsPublishedCaptureExamples' ./internal/publication/`: PASS (all 4 + 2 subtests, 0.7 s).
- Built the CLI to scratch. End to end: `--example --raw-input <abs path>` with a ```` ```bash ```` run and a `## ` line, split at `==>`, placeholders replaced, `--manifest ... --dry-run` → `success`, `PUBLICATION-DRY-RUN`.
- Refusals checked: `--format json` → `PUBLICATION-USAGE`; `--example --dry-run`, `--example --manifest m.json`, `--manifest --example` → `PUBLICATION-USAGE`; `--raw-input` with no value, `--example --example` → `PUBLICATION-USAGE`; missing file, `../` path, symlink → `CAPTURE-RAW-INPUT-INVALID`; control byte → `CAPTURE-RAW-INPUT-UNSAFE`; `--raw-input` without `--example` and `answer --example` → `unknown publication option`.
- No `do-work/` or reservation marker appeared in the scratch repo after any call.
- Test judgment: `TestCaptureFilesExampleFilledInPassesDryRun` feeds the printed output to the real handler's dry run. It would catch a renamed placeholder (the id then fails the pattern), a dropped `raw_input`, a hand-rolled verbatim block, and any validator tightening that the templates miss. It is not decorative. Its gap is F1. The `UR input` case of `TestBuildCapturePlanAcceptsPublishedCaptureExamples` now pins the doc example under real containment, which is the check that should have caught the old fence.

### Suggested Additional Testing
- Release stage: unassessed. After finalization, run the installed wrapper `capture-files --example` once in a consumer repo.
- Edge case: raw input with CRLF line endings. `containedOutsideBytes` normalizes it to LF, and the example uses LF, so it should pass. Not run.

### Restatement Sweep
Redefined element: the Full Verbatim Input fence rule (longest backtick run + 1, minimum three, no info string, `> ` prefix) and the capture-reference UR example that shows it.
- `skills/do-work/actions/capture-reference.md:191-196`: consistent (the new home).
- `skills/do-work/actions/clarify.md:106`: stale. No minimum of three, no "no info string" (F2).
- `skills/do-work/actions/capture.md:136,138`: stale example, `text` info string and four backticks (F3).
- `skills/do-work/actions/capture.md:129`, `capture.md:236`, `stakeholder-answers.md:56`, `abandon.md:66,99`, `fan-out-reference.md:174`, `clarify.md:194`, `work-reference.md:77`: consistent. They cite the clarify contract or the capture-reference template by name and do not restate the fence shape.
- Go: `publication_manifest.go:107` and `requeststate/state_apply.go:1052` (abandon's own copy): consistent with each other and with the new text.

### Anti-bloat (check 4)
Items the REQ did not name: `parseCaptureExampleOptions` (needed, because the shared parser requires `--manifest`); the `captureExampleREQTemplate` const (follows PD-6); the `slices` import; the duplicate-`--example` refusal (F5); the unreachable encode-error branch (F4); the repeated JSON assertion in the test (F7). The new file, both flags and the one new test were named by the REQ or by pre-dispatch. No constraint is breached. Total bloat is about 10 lines.

### Scores (on the record — not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | 6 of 6 build requirements delivered; release pending finalization |
| Code Quality | 88% | Clean reuse of the real decoder types and containment function; F1 header gap, F4/F5 small extras |
| Test Adequacy | 90% | Real dry-run round trip, RED-GREEN recorded; F1 input class not covered |
| Scope | 100% | 5 declared files, 5 touched |
| Risk | Low | Read-only command; refusals reuse capture's own codes |
| Acceptance | Pass | Focused tests and by-hand CLI in a scratch repo |

### Follow-ups created
None (7 findings report only; F1 is recommended as a fix before release)

## Review

**Overall: 94%** | 2026-10-10T16:56:08Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. Fix F1 (one line) before release.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `capture_files_example.go:79-81` header tells sessions to replace every placeholder, which corrupts the verbatim block when the raw input contains `UR-NNN`/`REQ-NNN` (dry run refuses; reproduced). Add the "leave the `> ` lines as printed" line in an else branch — impact-user-visible → fix before release (exact text in the review report); F2 `clarify.md:106` fence rule lacks the three-backtick minimum and the no-info-string rule that all Go writers apply, while `capture-reference.md:196` cites it with the stricter rule — impact-rule-change → report only; F3 `capture.md:136,138` addendum example still uses a four-backtick `text` fence — impact-negligible → report only; F4 unreachable encode-error branch `capture_files_example.go:70-72` — impact-negligible → report only; F5 duplicate-`--example` refusal and `--raw-input --example` path capture `capture_files_example.go:93-100` — impact-negligible → report only; F6 Go REQ template is a second copy of the doc shape, pinned for validity only — impact-negligible → report only; F7 JSON-refusal assertion repeated in both subtests `capture_files_test.go:327` — impact-negligible → report only.
**Acceptance:** Pass — focused Go tests plus by-hand CLI in a scratch git repo (end-to-end dry run, 13 refusal paths, no writes); release and consumer stages not exercised.
**Restatement sweep:** redefined the Full Verbatim Input fence rule (longest run + 1, minimum three, no info string, `> ` prefix) and the capture-reference UR example; stale: `clarify.md:106` (F2), `capture.md:136,138` (F3); consistent: `capture.md:129,236`, `stakeholder-answers.md:56`, `abandon.md:66,99`, `fan-out-reference.md:174`, `clarify.md:194`, `work-reference.md:77`, Go `publication_manifest.go:107`, `state_apply.go:1052`.
**Suggested testing:** 2 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*
