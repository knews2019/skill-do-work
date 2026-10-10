---
id: REQ-658
title: 'finalize --auto-manifest builds the mechanical fields of the finalization manifest and preflights the tree'
status: claimed
created_at: 2026-10-09T21:14:04Z
user_request: UR-145
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-659, REQ-660, REQ-661]
batch: cli-ergonomics
claimed_at: 2026-10-10T12:52:19Z
---
# finalize --auto-manifest Builds the Mechanical Fields of the Finalization Manifest and Preflights the Tree
## What
`do-work-cli finalize --auto-manifest REQ-N` builds the finalization manifest from the REQ, the live files and the finalization planner. The action still supplies every judged field as an input. The command fills the mechanical fields, preflights the tree before writing anything, and prints a short outcome table. `--emit <path>` writes the manifest and stops, which is the validate mode sessions are missing today.
## Why
Report item 1 (UR-145 input, "What happened"): 171 `advance ... --finalization-manifest <hand-built json>` commands in about 30 sessions over 15 days. Sessions wrote helpers named `make_manifest.py`, `finalize.py`, `fin.sh` and `adv.sh`, or inlined the same logic. A per-machine note says `commit_paths` is found by "submit a first guess and read the refusal, which names the exact missing paths", although the planner already computes that set.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/actions/work-reference.md:719`: "The action judges and authors one strict finalization manifest..." (the report cites `:810` at 0.305.84). Judged and mechanical fields sit in one hand-built file.
- `skills/do-work/actions/work.md:467`: "Run the exact `advance` continuation with the selected request path and the single action-authored finalization manifest" (report: `:475`).
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:264-265`: the missing-evidence hint is the literal placeholder `<action-authored-finalization-manifest>`.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go:32-49`: the `Manifest` struct is the only schema. No example manifest ships.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go:318-334`: `finalize` accepts only `--manifest <file>`. No dry run, no validate mode.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:158-174`: the planner computes the required commit paths (`statePlan.TargetPaths` plus release postimages) and only reports them in the refusal "commit_paths omits planned lifecycle or release targets".
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go:648-662`: refused exits 1, rolled back 3, risk 4, failure 2.
## Detailed Requirements
1. Add `--auto-manifest REQ-N` to `finalize`. Judged inputs: `--transition`, `--terminal-status`, `--message-file`, and provenance as either `--provenance supplied_commit --implementation-hash <merge>` or `--provenance primary_commit`, plus optional `--release-manifest`.
2. Fill the mechanical fields: request path from the working REQ; `expected_request_sha256` and `expected_checkpoint_sha256` from the live files; `completed_at` from the same source as `do-work-cli now` (Timestamp rule, `actions/work-reference.md:59`); the writer label; and `commit_paths` as the planner's required set from `finalization_prepare.go:158-174` plus any `--extra-path`.
3. Preflight before writing anything: when a release manifest is given, the project version still matches its old version; the index is clean; nothing staged lies outside `commit_paths`.
4. `--emit <path>` writes the manifest and stops.
5. On success or refusal print one table: outcome, phase, archived path, commit hash, then one line per finding (code and evidence) and one per gate.
6. Exit codes stay as `result_model.go:648-662` defines them.
7. The judged fields are never invented. A missing judged input (for example `--message-file`) refuses.
8. Change the action lines that describe the hand-built manifest (`actions/work-reference.md:719` area, `actions/work.md:467`) to name the new command, and keep the hand path described as valid.
9. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No new frontmatter field, no new status, no change to who judges the commit message, release payload or transition. This command fills mechanical fields only.
- `finalize --manifest <file>` and `advance ... --finalization-manifest <file>` keep working unchanged.
- Out of scope: a gate runner with failure-only summaries and flake repeats.
- This REQ is its own release; do not fold REQ-659, REQ-660 or REQ-661 into it.
## Assumptions (recorded at capture, no questions asked)
- The work loop runs finalization through `advance ... --finalization-manifest`, not bare `finalize`. The `--emit` output must be accepted unchanged by both `finalize --manifest` and `advance --finalization-manifest`, so the action can build with `--auto-manifest --emit` and continue through `advance` as today. Whether `advance` also gains a direct auto mode is the builder's call; record it in Decisions.
- The `fail` transition's `failure_type` and `failure_error` are judged fields, so they are inputs (`--failure-type`, `--failure-error-file` or similar), never derived.
- The writer label comes from an input flag. If the CLI already has a default writer label elsewhere, reuse it; never invent a new label vocabulary.
- `release_at` is filled from the same instant as `completed_at` when a release manifest is given, unless the existing finalize contract says otherwise.
- "The working REQ" means the one file for REQ-N in `do-work/working/`. Zero or several matches refuse.
- The version preflight runs only when `--release-manifest` is given.
## Dependencies
None. Independent of REQ-659, REQ-660 and REQ-661.
## Builder Guidance
Certainty is high on behaviour; the report lists inputs, filled fields, preflight and output. Latitude: flag spelling for the failure fields and the writer label, the table layout, and whether the planner's required-path computation is exposed through a shared helper or called directly.
## Red-Green Proof
**RED prompt/case:** On a fixture REQ ready for finalization, run `finalize --auto-manifest REQ-N --message-file msg.txt --transition complete --terminal-status completed --provenance primary_commit --emit m.json`.
**Why RED now:** `finalization_commands.go:318-334` rejects any option other than `--manifest` ("unknown finalize option").
**GREEN when:** The command writes `m.json`, and `finalize --manifest m.json` accepts it unchanged. A staged unrelated path makes the command refuse with exit 1 before any journal is written. Omitting `--message-file` refuses.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers commit preflight ordering and do-work-cli internals; family `commit-preflight-before-side-effect` fits a preflight that must refuse before any journal is written.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: its index row covers prescribed command blocks, which this REQ rewrites in `work.md` and `work-reference.md`.
## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 1, What happened, Where the behaviour lives today Item 1, Proposed direction 1, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 1: "`finalize --auto-manifest REQ-N`: build the finalization manifest's mechanical fields from the REQ and the planner, preflight the tree, print a short outcome table."*
