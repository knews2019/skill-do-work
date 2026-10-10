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
route: B
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go", "skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go", "skills/do-work/actions/work.md", "skills/do-work/actions/work-reference.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
estimate:
  p50_active_minutes: 30
  confidence: medium
  basis:
  - Route B
  - 7-file write set
  - 2 subsystems involved
  - 8 acceptance criteria
  calculated_at: 2026-10-10T13:19:06Z
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
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens at claim, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers commit preflight ordering and do-work-cli internals; families `commit-preflight-before-side-effect` and `guard-on-one-entry-path` fit a preflight that must refuse before any journal is written and an emit path that must run the same checks as the real one.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: its index row covers prescribed command blocks, which this REQ rewrites in `work.md` and `work-reference.md`.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over budget; `slugged: partial`). Matching reason: its index row covers status contracts and downstream readers of action prose, which the Step 8/9 wording change touches.
## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 1, What happened, Where the behaviour lives today Item 1, Proposed direction 1, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 1: "`finalize --auto-manifest REQ-N`: build the finalization manifest's mechanical fields from the REQ and the planner, preflight the tree, print a short outcome table."*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (inputs, filled fields, preflight, emit mode, exit codes, prose lines), and the code lives in one Go package (`internal/finalization/`) plus two action paragraphs. Exploration confirmed the planner seam (`prepareBoundJournal`) and the contract checks that guard the action prose, so a separate Route C plan adds no information.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Pre-dispatch exploration at main `bd56c4b0` (0.305.87 lines still hold).

**Required-lessons consult.** `required_lessons` refreshed to `_dev/primes/lessons-releases.md` (666 tokens, `slugged: full`, taken whole; matches release payloads, since this REQ is its own release). Three over-budget satellites are recorded in the dropped section above. The builder brief points at five specific bullets of the dropped do-work-cli satellite (families `commit-preflight-before-side-effect`, `guard-on-one-entry-path`) as recommended reading.

**Where the code goes.**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go:33-42` `handleFinalize` and `:315-331` `parseFinalizeArguments` (only `--manifest`). The new flag set dispatches here.
- `finalization_prepare.go:39-188` `prepareBoundJournal` is the one planner path: decode (`finalization_journal.go:18-44`), index-empty check (`:51-53`), journal resume (`:54-69`), digest checks (`:71-78`), lifecycle plan (`:80-112`), release adopt and plan (`:114-153`, which writes payload copies into the Git-private payload directory), required commit paths (`:155-171`), journal write (`:173-187`). The required set is `statePlan.TargetPaths` plus every release postimage path.
- `finalization_journal.go:46-99` `validateManifest` holds every field rule the emitted manifest must pass: writer non-empty, complete needs `completed` or `completed-with-issues`, fail needs `failure_error` plus `failure_type` in intent/spec/code/environment, RFC3339 times, digests, non-empty message and commit paths, provenance rules (`supplied_commit` runs `validateSuppliedCommit`), and `release_manifest_path` with `release_at` together.
- `finalization_types.go:32-49` `Manifest`: the emitted JSON is `json.Marshal` of this struct, so `finalize --manifest` (DisallowUnknownFields) accepts it by construction.
- Timestamp: `internal/corehelpers/commands.go:73` formats `time.Now().UTC().Truncate(time.Second)` as `2006-01-02T15:04:05Z`.
- Writer default: `internal/requeststate/state_plan.go:342-346` already defaults an empty writer to `hostname + ":" + repositoryRoot` (the same shape the live CHECKPOINT claims carry). `internal/lifecycleadvance/queue_commands.go:424` has a near copy, but `lifecycleadvance` imports `finalization`, so `finalization` cannot call it.
- Release version: `internal/publication/release.go:46-49` refuses `RELEASE-PREIMAGE-STALE` when a version target no longer matches its expected bytes, which is the "project version still matches the old version" check.
- Exit codes: `internal/resultmodel/result_model.go:648-662`. Refused is 1; `commandFailure` (`finalization_commands.go:335`) is Failure, exit 2, which is what every `prepareBoundJournal` error returns today.
- Text output: `result_model.go:1023-1110` `renderText` already prints the outcome line, every finding (code and evidence), and per finalization record the phase, archive path, and primary and metadata commits.
- `advance` binds the manifest through `internal/lifecycleadvance/finalization_gate.go:35` `FinalizeBound`; its missing-evidence hint is `advance_commands.go:263-265`.
- Test fixtures: `finalization_commands_test.go:231-277` (`newFinalizationRepository`, `writeFinalizationFile`, `runFinalizationGit`, `digestFile`); `finalization_req499_test.go:638-675` `seedPlannedReleaseFinalization` builds a claimed REQ-760 with a VERSION 1.0.0 to 1.0.1 release manifest on disk, and `finalize --manifest` accepts its hand-built manifest.

**Contract checks that guard the prose.** `_dev/tests/contracts/core-checks.sh:752-788` captures work.md Step 8, Step 9, and work-reference's Changelog and Commit sections. It requires the token "author exactly one strict manifest" (Step 8); the tokens "exact `advance` continuation", "ordered `finalizations`", "`phase: cleanup_complete`" and "empty `blocked_paths`/`reason_codes`" (Step 9); and the sentence "Pass that one manifest to the current `advance` continuation." (Commit section); it caps Step 9 at 8 lines; and it rejects, per line, the regexes `do-work-cli\.sh.*finalize` (which also matches "finalization" later on a line naming the wrapper), `finalize --manifest`, `git add`, `git commit`, `record-commit-hash`. The script runs in about 8 s.

**Seam check against the run table.** `work.md`: REQ-659 (stamp lines near Step 6/7) and REQ-689 (coordinator prose) edit other steps; this REQ edits only Step 8 item 5 and Step 9. `work-reference.md`: REQ-659 and REQ-690 edit other sections; this REQ edits only the Commit & Metadata-Commit paragraph at `:719`. `lessons-do-work-cli.md`: REQ-659 and REQ-660 add bullets; this REQ touches only the `internal/finalization/` module bullet at `:121`. `internal/requeststate/state_plan.go` is not in any sibling's known set. Missed by the table: none found.

**Pre-dispatch decisions (D-01 to D-10; the builder continues at D-11).**
- D-01: `advance` gets no direct auto mode. The action runs `finalize --auto-manifest ... --emit <path>`, then the existing `advance ... --finalization-manifest <path>`. Value: no `lifecycleadvance` change, one new surface. Risk: two commands per finalization; reversible by adding the flag later.
- D-02: Writer label is an optional `--writer`, defaulting to the existing `hostname:repositoryRoot` expression. Export it once as `requeststate.DefaultWriterLabel(repositoryRoot)` from `state_plan.go:342-346` and call it there and in the new code. No new label vocabulary (capture assumption).
- D-03: The preflight "index clean; nothing staged outside commit_paths" is one check, an empty index, because the finalizer already refuses any non-empty index (`finalization_prepare.go:51-53`), so the second clause can never be the deciding one. The refusal lists the staged paths.
- D-04: The version preflight is the existing release planner, run through the shared prepare core; its `RELEASE-PREIMAGE-STALE` refusal is the version mismatch. No new version reader.
- D-05: `--emit` validation runs the same prepare core as the real finalize, in a dry mode that adopts release payloads into a temporary directory, writes no journal and creates nothing under `.git`, so an emitted manifest has passed every check `finalize --manifest` would make (families `guard-on-one-entry-path`, `commit-preflight-before-side-effect`).
- D-06: Preflight refusals and missing judged inputs return OutcomeRefused (exit 1). Unknown options stay `FINALIZATION-USAGE` failure (exit 2), as today. The plain `finalize --manifest` path keeps its current outcomes.
- D-07: The outcome table is the existing text rendering, which already prints outcome, findings, and each record's phase, archive path, and commit hashes. No new renderer. Emit mode adds one info finding with the emitted path and the commit_paths. Value: no output code to maintain. Risk: the text is about 15 lines per record instead of a compact table; reversible.
- D-08: When an unfinished journal exists for the REQ, auto mode refuses and names `recover-finalization`, because a freshly built manifest (new `completed_at`) can never match the journalled manifest digest.
- D-09: The fail transition's judged inputs are `--failure-type` and `--failure-error-file`; `--terminal-status` is required for `complete` and refused for `fail`.
- D-10: The `advance` missing-evidence hint (`advance_commands.go:263-265`) stays unchanged; changing `advance` output is out of scope.

*Generated by pre-dispatch exploration*

## Scope

**Files I will touch:**
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest.go (new) — flag parsing, working-REQ lookup, mechanical fields, preflight, emit
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go (new) — the four Red-Green tests
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go (modify) — split prepareBoundJournal into decode plus a shared core with a dry mode that reports the required commit paths
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go (modify) — dispatch the auto-manifest flag set, usage text
- skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go (modify) — export DefaultWriterLabel
- skills/do-work/actions/work.md (modify) — Step 8 item 5 and Step 9 name the command
- skills/do-work/actions/work-reference.md (modify) — the Commit and Metadata-Commit paragraph names the command, hand path stays valid
- skills/do-work/tools/do-work-cli/lessons-do-work-cli.md (modify) — one clause on the finalization module bullet

**Files I will NOT touch:** `internal/lifecycleadvance/` (no advance auto mode, hint unchanged), `internal/resultmodel/` (no new renderer), `internal/publication/`, `_dev/tests/contracts/core-checks.sh`, CHANGELOG, VERSION and mirrors (the integrator's release).

**Acceptance criteria (restated from REQ):**
- [ ] `finalize --auto-manifest REQ-N` takes the judged inputs (transition, terminal status, message file, provenance with implementation hash, optional release manifest, failure fields) and refuses with exit 1 when one is missing.
- [ ] It fills request path (exactly one working file for REQ-N, else refuse), both digests from the live files, `completed_at` from the `now` source, the writer label, `release_at` equal to `completed_at` when a release manifest is given, and `commit_paths` as the planner's required set plus every `--extra-path`.
- [ ] Before writing anything it refuses (exit 1, no journal, no payload directory) on a non-empty index and on a release version that no longer matches.
- [ ] `--emit <path>` writes the manifest and stops; `finalize --manifest` and `advance --finalization-manifest` accept that file unchanged.
- [ ] Success and refusal print outcome, phase, archived path, commit hash, and one line per finding.
- [ ] Exit codes stay as `result_model.go:648-662` defines them.
- [ ] `actions/work.md` Step 8/9 and `actions/work-reference.md` Commit paragraph name the command, keep the hand path valid, and `_dev/tests/contracts/core-checks.sh` still passes.
- [ ] `finalize --manifest <file>` and `advance ... --finalization-manifest <file>` keep working unchanged.
