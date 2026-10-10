---
id: REQ-688
title: 'capture-files --example prints one valid manifest with payload templates, and the capture-reference fence example is fixed'
status: pending
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-658, REQ-659, REQ-660]
batch: cli-ergonomics
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
