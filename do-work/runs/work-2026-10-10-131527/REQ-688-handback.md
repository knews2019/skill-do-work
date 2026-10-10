# Hand-back: REQ-688 (capture-files --example prints a valid manifest with payload templates; capture-reference fence fix)

- Branch: `worktree-agent-REQ-688-capture-files-example-fence-fix`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-688-capture-files-example-fence-fix`
- Base commit: `bd56c4b0`
- Commits: `db502d32` `[REQ-688] add capture-files --example and fix the capture-reference fence` (one commit)

## File manifest
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_example.go` (new): `handleCaptureFilesExample` prints the manifest (encoded from `Manifest`), the UR template and the REQ template, each after a `==> <path> <==` line; `--raw-input` reads through `readPayload` + `validateOutsideBytes` and fills the UR verbatim block with `containedOutsideBytes(raw, "\n")`; refuses under `--format json`.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (modified): 3 lines at the top of `handlePublicationCommand` route `capture-files` with `--example` to the new handler before `parseCommandOptions`; `slices` import.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` (modified): the `UR input` case of `TestBuildCapturePlanAcceptsPublishedCaptureExamples` now supplies raw input `add keyboard shortcuts\n` and sets `RawInput`; new `TestCaptureFilesExampleFilledInPassesDryRun` (two subtests: with and without `--raw-input`).
- `skills/do-work/actions/capture-reference.md` (modified): UR example fence lines become `> ```` (bare); the sizing sentence now states the real derivation once and points at `capture-files --example --raw-input <file>`.
- `skills/do-work/actions/capture.md` (modified): one sentence after the Step 5 manifest paragraph naming the wrapper line with `capture-files --example --raw-input "<raw-input-path>"`.

## P-A-U
- **[PLAN]:** Branch `--example` for capture-files only before the shared option parser. Build the manifest as a `Manifest` value and encode it; keep both payload templates as Go text (PD-6). Reuse `readPayload`, `validateOutsideBytes`, `containedOutsideBytes`, `commandFailure`, `refusalResult`, `refusedPlan`. Pin with two tests: the published UR example under raw-input containment, and the printed example filled and dry-run.
- **[APPLY]:** Done as planned in the five Write-boundary files. No other file touched; no file under `do-work/` staged.
- **[UNIFY]:** `git diff bd56c4b0 --stat`:
  ```
   skills/do-work/actions/capture-reference.md        |   6 +-
   skills/do-work/actions/capture.md                  |   2 +
   .../internal/publication/capture_files_example.go  | 112 +++++++++++++++++++++
   .../internal/publication/capture_files_test.go     | 105 ++++++++++++++++++-
   .../internal/publication/publication_commands.go   |   4 +
   5 files changed, 221 insertions(+), 8 deletions(-)
  ```
  Checks (from the worktree root): `gofmt -l .../internal/publication` printed nothing, exit 0. `go vet -C skills/do-work/tools/do-work-cli ./internal/publication/` exit 0 (under 1 s). `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/publication/` exit 0, 38.3 s (machine under load from sibling builders, single run, passed). `REQ-688-probe.sh` exit 0, 2 s, printed `REQ-688 GREEN probe: ok`. `git diff --check` and `git diff --cached --check` exit 0. Files checked: all five above, read in full in the diff; no debug output, no stray files (`git status --short` empty after the commit).

## Proof record
1. **RED, CLI at base `bd56c4b0`:** `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" capture-files --example` printed `finding PUBLICATION-USAGE [error]: unknown publication option "--example"`, exit 2.
2. **RED then GREEN, doc example:** with only the `UR input` test change applied, `go test -count=1 -v -run '^TestBuildCapturePlanAcceptsPublishedCaptureExamples$' ./internal/publication/` failed: `--- FAIL: .../UR_input` with `Refusal{Code:"CAPTURE-RAW-INPUT-NOT-CONTAINED", ...}`; the other three cases passed. After the `capture-reference.md` fix: all four subtests PASS.
3. **RED then GREEN, example test:** `TestCaptureFilesExampleFilledInPassesDryRun` written first; both subtests failed with `example refused: ... Evidence:["unknown publication option \"--example\""]`. After `capture_files_example.go` and the routing: both subtests PASS (0.31 s).
4. **GREEN, by hand:** in a `mktemp -d` git repo with `raw.md` = `add a fenced sample:\n```bash\necho hi\n```\n## Not a heading\n`, ran the worktree wrapper `--repo-root <scratch> capture-files --example --raw-input raw.md` (exit 0, nothing written to the scratch repo), split at the `==>` lines into a second `mktemp -d`, replaced `<payload-dir>`, `UR-NNN`→`UR-001`, `REQ-NNN`→`REQ-001`, `<slug>`→`fenced-sample`, `<created-at>`→`2026-10-10T13:30:00Z`, `<word-count>`→`9`, then `capture-files --manifest <payload>/manifest.json --dry-run`: first run `capture-files: success`, `PUBLICATION-DRY-RUN`, exit 0. Both scratch dirs removed.
   Also checked by CLI: `--format json capture-files --example` refuses `PUBLICATION-USAGE` naming `--format text`; `--example --dry-run` refuses `PUBLICATION-USAGE`; `--raw-input x` alone refuses `unknown publication option "--raw-input"`; `answer --example` refuses `unknown publication option "--example"`; a missing raw file refuses `CAPTURE-RAW-INPUT-INVALID`.

Exact printed output for the scratch raw file:
````text
# capture-files example. Write each part below to the file its ==> line names, then replace every placeholder:
# UR-NNN, REQ-NNN (the ids capture reserved), <payload-dir> (a private temporary directory), <slug>, <title>,
# <created-at> (whole-second UTC, e.g. 2026-01-26T10:00:00Z), <word-count>, and the angle-bracket body text.
# Then check it: capture-files --manifest <payload-dir>/manifest.json --dry-run
==> <payload-dir>/manifest.json <==
{
  "operation": "capture-files",
  "commit_message": "[UR-NNN] captured <title> (1 REQs)",
  "capture": {
    "user_request_id": "UR-NNN",
    "user_request": {
      "path": "do-work/user-requests/UR-NNN/input.md",
      "payload": {
        "source_path": "<payload-dir>/ur-input.md"
      }
    },
    "raw_input": {
      "source_path": "raw.md"
    },
    "requests": [
      {
        "id": "REQ-NNN",
        "user_request_id": "UR-NNN",
        "file": {
          "path": "do-work/queue/REQ-NNN-<slug>.md",
          "payload": {
            "source_path": "<payload-dir>/req.md"
          }
        },
        "reservation_path": "do-work/.req-reservations/REQ-NNN"
      }
    ]
  }
}
==> <payload-dir>/ur-input.md <==
---
id: UR-NNN
title: '<title>'
created_at: <created-at>
requests: [REQ-NNN]
word_count: <word-count>
---
## Full Verbatim Input
> ````
> add a fenced sample:
> ```bash
> echo hi
> ```
> ## Not a heading
> ````
==> <payload-dir>/req.md <==
---
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: <original verbatim request>*
````
Without `--raw-input` the manifest has no `raw_input`, the note gains `# Rerun with --raw-input <file> to get the Full Verbatim Input block derived from the raw input bytes.`, and the UR block is `> ```` / `> <verbatim input>` / `> ````.

## Decisions
- **D-01 Delimiter:** DECIDE & STATE. `==> <path> <==`, where `<path>` is the payload file the part belongs in (`<payload-dir>/manifest.json`, `<payload-dir>/ur-input.md`, `<payload-dir>/req.md`), so the split also says where to write. Safe because raw-input lines in the UR part are `> `-prefixed (code comment says so).
- **D-02 Placeholders:** DECIDE & STATE. `UR-NNN`, `REQ-NNN`, `do-work/.req-reservations/REQ-NNN` (fixed by the REQ), `<payload-dir>`, `<slug>`, `<title>`, `<created-at>`, `<word-count>`, plus free body text `<what ...>`, `<original verbatim request>`, `<verbatim input>`. Only `UR-NNN`, `REQ-NNN`, `<payload-dir>`, `<slug>`, `<created-at>`, `<word-count>` must be replaced to pass the dry run; `<title>` and body text pass as is.
- **D-03 Valid defaults instead of placeholders for enum fields:** DECIDE & STATE. `domain: general`, `tdd: false`, `maintenance: false`, `prime_files: []`, `status: pending` are canonical values, so the validator passes them; a `<domain>` placeholder would fail the canonical-value check and add one more required replacement. Capture Step 1 still owns the real choice.
- **D-04 Encoder instead of `json.MarshalIndent`:** DECIDE & STATE. `MarshalIndent` escapes `<` and `>` as `<`/`>`, which makes the placeholders unreadable and unreplaceable by plain text substitution. `json.NewEncoder` with `SetEscapeHTML(false)` and `SetIndent("", "  ")` marshals the same `Manifest` value. Still never hand-written JSON.
- **D-05 Routing location:** DECIDE & STATE. Top of `handlePublicationCommand`, `operation == OperationCaptureFiles && slices.Contains(arguments, "--example")`. So `--manifest/--dry-run/--commit/--at` beside `--example` reach `parseCaptureExampleOptions` and refuse `PUBLICATION-USAGE`; `--raw-input` without `--example` still reaches `parseCommandOptions` and refuses as unknown (PD-2).
- **D-06 Usage refusal before the format check:** DECIDE & STATE. A bad option combination is reported first; then `--format json` refuses. Both are `PUBLICATION-USAGE`.
- **D-07 Raw-input read refusals:** DECIDE & STATE. Returned as refusals (`refusalResult(refusedPlan(...))`) with `CAPTURE-RAW-INPUT-INVALID` / `CAPTURE-RAW-INPUT-UNSAFE`, the same outcome and codes `BuildCapturePlan` gives.
- **D-08 Step 5 sentence placement:** DECIDE & STATE. Its own short paragraph right after the manifest paragraph (`capture.md:230`), inline wrapper command in the same style as `verify-requests.md:168`; it says plain text only, because the Step 5 command line carries `--format json` and a session copying it would get the refusal.
- **D-09 Test splits with a regexp:** DECIDE & STATE. The test parses `^==> (.+) <==\n` itself rather than importing a production splitter; no splitter exists in production because the CLI never reads its own example back.

## Discovered Tasks
- `skills/do-work/actions/capture.md:316` (final checklist) says the commit format is `[UR-NNN] captured: ...` with a colon, while `:271` (Step 7) defines `[UR-NNN] captured {title} ({N} REQs)` without one. Two spellings of one format → report only.
- `skills/do-work/actions/clarify.md:106` (Outside-text containment) says "a code fence longer than the longest backtick run", without the three-backtick minimum or the no-info-string rule that `containedOutsideBytes` applies. Capture byte-checks the UR block, so capture-reference now states the exact derivation; clarify's answer notes are not byte-checked, so this is a precision gap, not a failure → report only.

## Lessons read
- `_dev/primes/lessons-releases.md` (whole; families `canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`; only bear on the integrator's release and lesson link).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: the three `[family: alternate-writer-contract-drift]` bullets (REQ-489, REQ-547, REQ-504). Applied: the example is a second writer of the capture-files contract, so the test feeds its printed output to the real handler's dry run instead of checking its shape.
- `_dev/primes/lessons-action-files.md`: the `[family: restated-mechanism-unchecked]` bullet (REQ-639). Applied: the new capture-reference sentence was written from `containedOutsideBytes` and `capture_files.go:77-92`, and is pinned by the `UR input` test case.

## Anti-bloat check
`git diff --stat` is in [UNIFY] above. Added names the brief did not name:
- `handleCaptureFilesExample` (func): the handler the brief describes; it needs a name.
- `parseCaptureExampleOptions` (func): the example form's option parser; `parseCommandOptions` requires `--manifest`, so it cannot be reused.
- `captureExampleREQTemplate` (const): the REQ template text (PD-6, templates in Go); a const keeps the P-A-U backticks out of the handler body.
- `slices` import in `publication_commands.go`: for the one routing check.
Nothing else added: no new flag beyond `--example`/`--raw-input`, no manifest or frontmatter field, no file writes, no test beyond the one named (with its two subtests).

## Proposed CHANGELOG entry
**capture-files --example prints a valid manifest and payload templates, and the capture reference shows the fence the command accepts**

Sessions writing a capture by hand had to grep the Go source for the manifest keys, and the reference UR example's four-backtick `text` fence failed the raw-input check on the first dry run. Now one command prints a manifest and both payloads that pass the dry run once the placeholders are filled.

- New `do-work-cli capture-files --example [--raw-input <file>]`: prints the manifest, the UR template and the REQ template, each after a `==> <path> <==` line. With `--raw-input`, the UR template already holds the exact verbatim block derived from the file's bytes, and `raw_input` names the file. It writes nothing and refuses under `--format json`.
- `actions/capture-reference.md`: the UR example now uses a bare three-backtick fence, and the paragraph after it states the real fence rule once and points at the new command.
- `actions/capture.md` Step 5 names the command as the way to get a valid manifest.
- The published-example test now runs the UR example through the raw-input containment check, so a wrong fence can no longer ship green.

## Proposed lesson bullet
Satellite: `_dev/primes/lessons-action-files.md`.
`- [family: restated-mechanism-unchecked] [REQ-688: a test that "accepts the published examples" skipped the one check the example existed for (no RawInput on the UR case), so a four-backtick text fence that every real capture refused shipped green; when a doc example is pinned by a test, feed the test the same inputs the real command gets](../../do-work/archive/UR-153/REQ-688-capture-files-example-and-fence-fix.md#lessons-learned)` (the integrator fixes the archive path).

## Integration seams
- `capture.md`: my only edit is the new paragraph at `:230-231` (after the manifest paragraph at `:228`). Nothing inserted above `:226`, so REQ-691's line citations (`:43-44, :78, :113, :115, :125, :162, :222`) hold, and REQ-692's possible clause near `:43-44` does not overlap.
- `capture-reference.md`: edits at `:191`, `:193`, `:196` only; REQ-691's citations `:113, :125` unchanged. Line count unchanged.
- `internal/publication/`: no sibling touches it. `internal/resultmodel/result_model.go` not touched; `ExactTextOutput` used as it stands.
- Test wall times: focused `TestBuildCapturePlanAcceptsPublishedCaptureExamples` ~0.4 s, `TestCaptureFilesExampleFilledInPassesDryRun` 0.31 s (two `initializedGitRepository` repos), whole package 38.3 s under load, probe 2 s.
