---
id: REQ-688
title: 'capture-files --example prints one valid manifest with payload templates, and the capture-reference fence example is fixed'
status: claimed
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-658, REQ-659, REQ-660]
write_set: ["skills/do-work/tools/do-work-cli/internal/publication/capture_files_example.go", "skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go", "skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go", "skills/do-work/actions/capture-reference.md", "skills/do-work/actions/capture.md"]
batch: cli-ergonomics
claimed_at: 2026-10-10T13:14:50Z
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-10-10T13:19:39Z
  basis:
    - trivial short-circuit
---
# capture-files --example Prints One Valid Manifest, and the capture-reference Fence Example Is Fixed
## What
`do-work-cli capture-files --example [--raw-input <file>]` prints to stdout one valid `capture-files` manifest with the UR and REQ payload templates it refers to, every required key present with a placeholder, and, when `--raw-input` is given, the UR template already carrying the byte-derived verbatim block. In the same change the `actions/capture-reference.md` UR example is fixed to the fence shape the command actually accepts. This is the fallback the source report offered in place of a full `capture-files init`; the maintainer chose it on 2026-10-10 because the report itself called this item's evidence the thinnest of its four.
## Why
On 2026-10-05 in this checkout the first `capture-files --dry-run` refused with `CAPTURE-RAW-INPUT-NOT-CONTAINED` because the UR payload copied the four-backtick `text` fence from the reference example. On 2026-10-09 a consumer session grepped `json:"` in `capture_files.go` to recover the manifest schema; the same hunt appears in two other repos. The schema hunt and the fence mismatch are concrete; a filled-in example answers both.
## Verified Facts (checked at 0.305.101)
- `skills/do-work/actions/capture.md` Step 5 describes the manifest in prose only; the `capture-files --manifest ... --commit` line is the only command shown.
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` holds the JSON tags for `Manifest`, `CaptureManifest`, `CaptureRequest`, `PublishedFile` and `PayloadFile`. No example manifest ships.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` refuses a UR payload that does not contain `containedOutsideBytes(raw)`; per `publication_manifest.go` that block is a fence of max(longest backtick run + 1, 3) backticks with no info string, every line prefixed `> `.
- `skills/do-work/actions/capture-reference.md` → UR input.md still shows `> ````text` with a four-backtick close, and the paragraph after it says the fence is sized for the sample; the `text` info string alone still fails containment for that one-line input.
## Detailed Requirements
1. `capture-files --example` prints one valid manifest for a one-UR, one-REQ capture: every required key present, ids and reservation path as placeholders in the canonical stored-id spelling (`REQ-NNN`, `do-work/.req-reservations/REQ-NNN`), payload paths as placeholders.
2. The output also contains the UR and REQ payload templates the manifest names, matching the shapes in `actions/capture-reference.md`, delimited so a session can split them into files.
3. With `--raw-input <file>`, the UR template's Full Verbatim Input block is derived from the file's bytes by `containedOutsideBytes`, and the manifest's `raw_input` points at that file.
4. The example, once its placeholders are replaced with real ids and paths, passes `capture-files --dry-run` on the first try, including a raw input that contains a triple-backtick run.
5. Fix the `capture-reference.md` UR example: for its one-line sample input the derived fence is three backticks with no info string, so both fence lines become three bare backticks. Keep the sizing paragraph consistent with the fixed example.
6. `actions/capture.md` Step 5 names `capture-files --example` as the way to get the manifest shape; the prose description stays.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- `--example` writes nothing to disk and never creates a reservation marker; the `capture-files` transaction stays the only writer of markers.
- No new frontmatter field and no new status.
- Do not change the containment rule itself. If another shipped example carries a `text` info string inside a byte-checked verbatim block, fix that example in the same change.
- This REQ is its own release; do not fold REQ-658 (finalize --auto-manifest), REQ-659 (frontmatter set) or REQ-660 (worktree lifecycle) into it.
## Assumptions (recorded at capture)
- The full `init` subcommand from the cancelled REQ-661 (capture-files init writes a filled manifest skeleton) is not wanted; the maintainer chose the report's fallback on 2026-10-10. If the builder finds that `--example` cannot carry the derived verbatim block without writing a file, printing the block to stdout as a separate delimited part is acceptable.
- Every id in the example is a placeholder. The id scan stays prose in `actions/capture.md`; `--example` does not scan.
## Dependencies
None. Independent of REQ-658, REQ-659 and REQ-660.
## Builder Guidance
Certainty is high on the doc fix and the dry-run acceptance target. Latitude: flag spelling, placeholder spelling, output delimiters.
## Red-Green Proof
**RED prompt/case:** `capture-files --example --raw-input raw.md` where `raw.md` contains a triple-backtick run; replace the placeholders; run `capture-files --manifest <m> --dry-run`.
**Why RED now:** `capture-files` has no `--example` form; sessions copy the `capture-reference.md` example and the first dry run refuses with `CAPTURE-RAW-INPUT-NOT-CONTAINED`.
**GREEN when:** the dry run passes on the first try; the `capture-reference.md` UR example uses a bare three-backtick fence; a Go test proves the example block equals `containedOutsideBytes` of its sample input and that the printed example, placeholders filled, passes the dry run.
**Validation:** Inferred during capture (from the source report's Acceptance check and the maintainer's 2026-10-10 decision).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: family `alternate-writer-contract-drift` fits an example that must match the capture-files contract byte for byte.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over budget; `slugged: partial`). Matching reason: family `restated-mechanism-unchecked` fits a doc example that restates the fence mechanism with no check.
## Full Context
See `do-work/user-requests/UR-153/input.md` for the decision record. The cancelled original is `do-work/archive/REQ-661-capture-files-init-manifest-skeleton.md`; its source report is `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 4.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: maintainer decision of 2026-10-10 (UR-153), replacing REQ-661 with the source report's fallback: "If a new subcommand is too much, a `capture-files --example` that prints one valid manifest covers most of the gap."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (a new `--example` form on an existing command, one doc example fix, one Step 5 pointer), but the output shape, the option routing inside the shared publication option parser, and the test seam needed discovery in `internal/publication/`. Two subsystems (Go CLI, action prose), no architectural change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Read by the pre-dispatch agent at main `bd56c4b0` (0.305.101 tree). All Verified Facts in the REQ body were re-checked against the code and hold.

**Required-lessons consult.** Index rows matched: `_dev/primes/lessons-releases.md` (666 tokens, `slugged: full`; the REQ ships a release, requirement 7, and `prime-releases.md` is in `prime_files`), added to `required_lessons` and read. Its families (`canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`) do not bear on the build itself; the first matters to the integrator only if a changelog or lesson link names an archive path. The two capture-time drops stay dropped (both whole satellites over budget and `slugged: partial`, so no targeted form). The builder brief quotes the four bullets of those two satellites that bear on this change (family `alternate-writer-contract-drift` in `lessons-do-work-cli.md`, family `restated-mechanism-unchecked` in `lessons-action-files.md`), as reading guidance only.

**Where the change goes (Go, `skills/do-work/tools/do-work-cli/internal/publication/`):**
- `publication_commands.go:34-43` `Handlers()` registers one closure per publication operation onto `handlePublicationCommand` (`:45`). `parseCommandOptions` (`:80-118`) is shared by all four operations and refuses any unknown option (`unknown publication option "--example"`, observed at base: exit 2, `PUBLICATION-USAGE`) and requires `--manifest`. The `--example` form must branch before that parser, for `capture-files` only, so `answer`, `release` and `defer-gate` keep refusing it.
- Printing: commands that print plain text set `resultmodel.CommandResult.ExactTextOutput` (`internal/resultmodel/result_model.go:637-640`, `json:"-"`); `renderText` (`:1023-1028`) writes it verbatim under the default `--format text`. Under `--format json` the field is dropped, so a session that copies the `--format json` from capture Step 5's command line would get a success with no example. Precedents: `corehelpers/commands.go:69-75` (`now`), `toolboxcommands/commands.go`. `ExecutionContext.Format` (`internal/commandruntime/command_runtime.go:14-17`) tells the handler which format was asked for.
- Containment: `containedOutsideBytes(contents, lineEnding)` (`publication_manifest.go:107-136`) is the derivation to reuse, never restate: fence = max(longest backtick run + 1, 3) backticks, no info string, every line prefixed `> `, trailing empty line dropped. `readPayload` (`:64-96`) resolves a relative path against the repository root (absolute paths pass through) and refuses symlinks; `validateOutsideBytes` (`:98-105`) refuses C0 controls and DEL. `capture_files.go:77-92` is the check the example must satisfy: it reads `raw_input` through `readPayload`, then requires the UR bytes to contain `containedOutsideBytes(raw, "\n")` (or `\r\n` when the UR uses CRLF), else `CAPTURE-RAW-INPUT-NOT-CONTAINED`.
- Schema the templates must satisfy: `validateCaptureURSchema` (`capture_files.go:181-211`: `id`, single-quoted `title`, whole-second UTC `created_at`, `requests` list equal to the manifest REQ ids, integer `word_count`) and `validateCaptureREQSchema` (`:213-283`: `id`, quoted `title`, `status: pending`, `created_at`, `user_request`, canonical `domain`, `prime_files` list, `tdd`, `maintenance`; `impact-user-visible` default must not be mirrored in the title). Paths: UR at exactly `do-work/user-requests/UR-NNN/input.md` (`:44-48`), REQ at `do-work/queue/REQ-NNN-<slug>.md` (`:96-99`), reservation at the canonical stored-id path (`:116-120`, `canonicalReservationPath`).
- Manifest JSON tags: `publication_types.go:19-60` (`Manifest`, `CaptureManifest`, `CaptureRequest`, `PublishedFile`, `PayloadFile`). The decoder is strict (`DisallowUnknownFields`, `publication_manifest.go:18-50`), so the example must be marshalled from these types (`json.MarshalIndent` of a `Manifest` value), never hand-written JSON text.
- Dry-run needs a git repository (`ApplyPlan`, `publication_commands.go:134` runs `gittransaction.ExecuteTransaction` even for `--dry-run`); the test helper `initializedGitRepository` (`publication_commands_test.go:226-236`) makes one with five git subprocesses.

**Tests (`capture_files_test.go`):** `TestBuildCapturePlanAcceptsPublishedCaptureExamples` (`:244-287`) already extracts the `### UR input.md` example from `capture-reference.md` and runs `BuildCapturePlan` on it, but its "UR input" case (`:275`) passes no `RawInput`, which is why the four-backtick `text` fence shipped green. Giving that case a raw input of the sample text (`add keyboard shortcuts`) turns it RED at base and GREEN after the doc fix: this is the Go test the REQ's GREEN line asks for ("the example block equals `containedOutsideBytes` of its sample input"). Helpers to reuse: `writeFixture` (`:380`), `canonicalCaptureManifest` (`:391`). `TestBuildCapturePlanBootstrapsAbsentDoWorkAndPinsRawInput` (`:78-100`) shows the raw-input manifest shape. Focused run time at base: under 1 s after compile.

**Docs:** `skills/do-work/actions/capture-reference.md:182-196` (UR input.md example: `> ````text` / `> ````` and the sizing paragraph at `:196`, which says "a fence longer than the longest backtick run" and "four-backtick fence is sized for its sample"). The example sits inside a ```` ```markdown ```` fence; a `> ```` line cannot close that outer fence (a closing fence may not start with `>`), and the test's extractor cuts at `\n```\n`, which `> ```` does not match. `skills/do-work/actions/capture.md` Step 5 (`:226-243`) describes the manifest in prose (`:228`) and shows only the `--manifest ... --commit` line (`:241`).

**Other shipped `text`-fenced verbatim blocks (REQ Constraint check):** a grep for a line-leading `> ` plus three or more backticks across `skills/` finds only `capture-reference.md:191` and `capture.md:136`. The second is the queued-addendum example: an addendum is appended to a pending REQ by hand or through a capture fold (`ReplacementFile`), and neither path runs the byte containment check, so it is not a byte-checked block and stays out of scope (PD-5 below).

**Pre-dispatch decisions** (this REQ has no Open Questions; these settle builder-latitude items so eleven parallel builds do not each re-derive them; labelled PD so they do not take the builder's D-XX numbers):
- **PD-1 Output channel:** `--example` returns `ExactTextOutput` (default `--format text`). Under `--format json` it refuses with a `PUBLICATION-USAGE` finding that says to use `--format text`, so it never reports a silent success with no example. Value: one `if`, and no empty success. Risk: low; reversible.
- **PD-2 Routing:** `--example` is recognised only for `capture-files`, before `parseCommandOptions`; it combines only with `--raw-input <file>` and refuses `--manifest`, `--dry-run`, `--commit` and `--at` beside it. `--raw-input` without `--example` stays an unknown option.
- **PD-3 Raw input path:** read with `readPayload` and `validateOutsideBytes` (same resolution and refusals as `capture-files` itself: relative to `--repo-root`, absolute allowed, symlink refused), and written into the manifest's `raw_input.source_path` exactly as given, so the printed manifest resolves to the same file. Without `--raw-input`, the manifest omits `raw_input` and the UR template's verbatim block is a bare three-backtick fence around one placeholder line; the output says to rerun with `--raw-input` for the byte-derived block.
- **PD-4 No disk writes:** the example handler never calls `BuildCapturePlan` or `ApplyPlan`; the test asserts no `do-work/` path appears in the temp repository after `--example`.
- **PD-5 Scope of the fence fix:** only `capture-reference.md`'s UR example and its sizing paragraph. `capture.md:136` (queued-addendum example) stays as is (not byte-checked; see above).
- **PD-6 Templates live in Go, not read from the doc at runtime:** the binary must not depend on finding `actions/capture-reference.md` beside it. The REQ template is the Simple REQ shape (with `## What` and the P-A-U block); the UR template is the simple UR shape. Validity is pinned by the dry-run test, not by a byte-equality test against the doc.

*Generated by pre-dispatch agent (Explore step done in-session)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_example.go` (new): the example builder and printer for capture-files --example and --raw-input
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (modify): route the example form for capture-files before the shared option parser
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` (modify): raw input on the published UR example case, and one new test that fills the printed example and dry-runs it
- `skills/do-work/actions/capture-reference.md` (modify): UR example fence becomes three bare backticks, sizing paragraph matches the derivation
- `skills/do-work/actions/capture.md` (modify): Step 5 names the example command as the way to get the manifest shape

**Files I will NOT touch:** `capture_files.go`, `publication_manifest.go` (the containment rule and the validator do not change), `publication_types.go` (no new manifest field), `skills/do-work/actions/clarify.md`, the queued-addendum example in `capture.md` (PD-5), `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` and every version mirror (the release is the integrator's), anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] `capture-files --example` prints one valid manifest for a one-UR, one-REQ capture with every required key present; ids and the reservation path use the canonical stored-id placeholder spelling (`REQ-NNN`, `do-work/.req-reservations/REQ-NNN`); payload paths are placeholders.
- [ ] The output also contains the UR and REQ payload templates the manifest names, in the capture-reference shapes, delimited so a session can split them into files.
- [ ] With `--raw-input <file>`, the UR template's Full Verbatim Input block is `containedOutsideBytes` of the file's bytes, and the manifest's `raw_input` names that file.
- [ ] The printed example, placeholders replaced with real ids and paths, passes `capture-files --dry-run` on the first try, including a raw input with a triple-backtick run.
- [ ] The `capture-reference.md` UR example uses a bare three-backtick fence on both fence lines, and the sizing paragraph agrees with it.
- [ ] `capture.md` Step 5 names `capture-files --example` as the way to get the manifest shape; the prose description stays.
- [ ] `--example` writes nothing to disk and never creates a reservation marker; no new frontmatter field or status; the containment rule is unchanged.
