---
id: REQ-661
title: 'capture-files init writes a filled manifest skeleton and payload templates, and the capture-reference fence example is fixed'
status: claimed
created_at: 2026-10-09T21:14:04Z
user_request: UR-145
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-658, REQ-659, REQ-660]
batch: cli-ergonomics
claimed_at: 2026-10-10T12:52:20Z
---
# capture-files init Writes a Filled Manifest Skeleton and Payload Templates, and the capture-reference Fence Example Is Fixed
## What
`do-work-cli capture-files init <dir> [--raw-input <file>]` writes a `manifest.json` with every required key filled by a placeholder or a scanned value (next UR id, next REQ id and its reservation path), plus UR and REQ payload templates. With `--raw-input`, the UR template already carries the byte-derived verbatim block. In the same change, the `actions/capture-reference.md` UR example is fixed to the fence shape the command actually accepts.
## Why
Report item 4 (UR-145 input, "What happened"): on 2026-10-05 in this checkout the first `capture-files --dry-run` refused with `CAPTURE-RAW-INPUT-NOT-CONTAINED` because the UR payload copied the four-backtick `text` fence from the reference example. On 2026-10-09 a consumer session grepped `json:"` in `capture_files.go` to recover the manifest schema; the same hunt appears in two other repos, and the tags are in a different file. The report says this item has the thinnest count of the four and is included because the schema hunt and the fence mismatch are concrete.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/actions/capture.md:228` describes the manifest in prose only; `:241` is the only command line shown.
- `skills/do-work/actions/capture.md:78-80`: the ID scan is prose ("The scan proposes IDs but writes nothing").
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go:19-26` (`Manifest`) and `:53-60` (`CaptureManifest`) hold the JSON tags.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go:89-90` refuses a UR payload that does not contain `containedOutsideBytes(raw)`. Per `publication_manifest.go`, that block is a fence of max(longest backtick run + 1, 3) backticks with no info string, every line prefixed `> `.
- `skills/do-work/actions/capture-reference.md:191-193` still shows `> ````text` and a four-backtick close. The paragraph at `:195` now says the four-backtick size is sized for the sample, but the `text` info string alone still fails containment for that one-line input.
## Detailed Requirements
1. `capture-files init <dir>` writes `manifest.json` with every required key filled with a placeholder or a scanned value: next UR id, next REQ id and its canonical `do-work/.req-reservations/REQ-NNN` path, from the same read-only scan `actions/capture.md:78-80` describes.
2. It writes UR and REQ payload templates matching the shapes in `actions/capture-reference.md`.
3. With `--raw-input <file>`, the UR template already carries the byte-derived verbatim block from `containedOutsideBytes`, and the manifest's `raw_input` points at that file.
4. The output, after the placeholders are filled, passes `capture-files --dry-run` on the first try, including a raw input that contains a triple-backtick run.
5. Fix the `actions/capture-reference.md:191-193` example: for that one-line input the command derives a three-backtick fence with no info string, so both fence lines become three backticks with no `text`. Keep `:195`'s sizing explanation consistent with the fixed example.
6. Change `actions/capture.md` Step 5 (`:228` area) to name `capture-files init`, keeping the hand path valid.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- `init` writes only into `<dir>`. It never writes into `do-work/` and never creates a reservation marker; the `capture-files` transaction stays the only writer of markers (`actions/capture.md:80`).
- No new frontmatter field and no new status.
- This REQ is its own release; do not fold REQ-658, REQ-659 or REQ-660 into it.
## Assumptions (recorded at capture, no questions asked)
- The report offers a fallback: "If a new subcommand is too much, a `capture-files --example` that prints one valid manifest covers most of the gap." The captured intent is the full `init`. The builder takes the fallback only if `init` proves unworkable, and records why in Decisions.
- `init` writes one REQ template by default. A `--requests N` flag (or similar) writes N templates with consecutive ids and reservation paths.
- `init` refuses a non-empty `<dir>` so it never overwrites a session's payloads.
- Every scanned id is a proposal, as the prose scan already says. A concurrent capture can still win the id; the transaction's existing refusal covers that, and `init` adds no lock.
- Payload `source_path` values are written as paths inside `<dir>`, since `capture-files` accepts absolute scratch paths.
- `actions/clarify.md` Step 4's Outside-text containment (`:106`) says "a code fence longer than the longest backtick run" and shows no info string requirement. If the builder finds any other shipped example with a `text` info string inside a byte-checked verbatim block, fix it in the same change; do not change the containment rule itself.
## Dependencies
None. Independent of REQ-658, REQ-659 and REQ-660.
## Builder Guidance
Certainty is high on the doc fix and on the dry-run acceptance target. Medium on the template content: reuse the capture-reference shapes, do not author new ones. Latitude: flag names, placeholder spelling, and whether the scan is shared with an existing Go helper.
## Red-Green Proof
**RED prompt/case:** Run `capture-files init /tmp/x --raw-input raw.md` where `raw.md` contains a triple-backtick run, fill the placeholders, then run `capture-files --manifest /tmp/x/manifest.json --dry-run`.
**Why RED now:** `capture-files` has no `init` form; sessions copy the `capture-reference.md:191` example and the first dry run refuses with `CAPTURE-RAW-INPUT-NOT-CONTAINED`.
**GREEN when:** The dry run passes on the first try. The `capture-reference.md` UR example uses a bare three-backtick fence, and a test or check proves the example block equals `containedOutsideBytes` of its sample input.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers do-work-cli internals and exact filesystem-basename authority; family `alternate-writer-contract-drift` fits a template writer that must match the capture-files contract byte for byte.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over budget; `slugged: partial`). Matching reason: family `restated-mechanism-unchecked` fits a doc example that restates the fence mechanism with no check.
## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 4, What happened, Where the behaviour lives today Item 4, Proposed direction 4, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 4: "`capture-files init <dir>`: write a filled `capture-files` manifest skeleton and the UR/REQ payload templates."*
