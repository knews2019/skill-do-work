---
id: REQ-658
title: 'finalize --auto-manifest builds the mechanical fields of the finalization manifest and preflights the tree'
status: completed
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
builder_handback_at: 2026-10-10T13:33:33Z
integration_at: 2026-10-10T15:05:54Z
review_at: 2026-10-10T15:18:24Z
remediation_at: 2026-10-10T15:19:32Z
kb_status: pending
re_review_at: 2026-10-10T15:26:38Z
commit: fe788bf4c2b864ace38bc4687fa0ac59dfeff7b9
heavy_verified_at: 2026-10-10T15:26:51Z
heavy_verified_revision: fe788bf4c2b864ace38bc4687fa0ac59dfeff7b9
completed_at: 2026-10-10T15:27:35Z
release_at: 2026-10-10T15:27:35Z
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
- [x] **[PLAN]:** (from the builder hand-back) Split `prepareBoundJournal` into decode plus a shared core with a dry mode (D-05). Auto mode: parse flags, refuse missing judged inputs (exit 1), preflight journal and index, find the one working REQ from the `DiscoverRepository` snapshot (`RequestsByID[id]` filtered by `TreeSection == "working"`), fill digests, `completed_at`, writer, then run the core dry on a draft (commit_paths = extra paths) to learn the required set, fill commit_paths, `validateManifest`, then the core again (dry for `--emit`, real otherwise followed by `advanceJournal`). Tests first.
- [x] **[APPLY]:** (from the builder hand-back) Implemented as planned within the eight write-boundary paths. RED recorded before any production edit.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff bd56c4b0 --stat`: 8 files, 474 insertions, 32 deletions (the eight `write_set` paths). `REQ-658-probe.sh` exit 0 (2 s); `REQ-658-preflight-probe.sh` exit 0 (13 s); `go test -count=1 ./internal/finalization/ ./internal/requeststate/` ok; `bash _dev/tests/contracts/core-checks.sh` exit 0; `gofmt -l` empty, `go vet` exit 0, `git diff --check` exit 0. All eight files checked: no debug output, no stray files, nothing under `do-work/` committed.
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
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest.go` (new) — flag parsing, working-REQ lookup, mechanical fields, preflight, emit
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go` (new) — the four Red-Green tests
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modify) — split prepareBoundJournal into decode plus a shared core with a dry mode that reports the required commit paths
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modify) — dispatch the auto-manifest flag set, usage text
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modify) — export DefaultWriterLabel
- `skills/do-work/actions/work.md` (modify) — Step 8 item 5 and Step 9 name the command
- `skills/do-work/actions/work-reference.md` (modify) — the Commit and Metadata-Commit paragraph names the command, hand path stays valid
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify) — one clause on the finalization module bullet

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

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest.go` (new)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified)

**What was done:** Added `do-work-cli finalize --auto-manifest REQ-NNN`. It takes the judged inputs (`--transition`, `--terminal-status` or `--failure-type` with `--failure-error-file`, `--message-file`, `--provenance` with `--implementation-hash` for `supplied_commit`, optional `--release-manifest`, `--extra-path`, `--writer`) and fills the request path, both digests, `completed_at`, `release_at`, the writer label and the planner's required `commit_paths`. It refuses with exit 1 before writing anything on a missing judged input, an unfinished journal, a staged index, or a planner refusal such as `RELEASE-PREIMAGE-STALE`. `--emit <path>` writes the manifest after the full check without a journal. `prepareBoundJournal` now decodes and calls one shared core, `prepareManifestJournal`, which has a dry mode (temporary payload directory, no journal) and reports the required commit paths. `requeststate.DefaultWriterLabel` is exported for the writer default. `work.md` Step 8 item 5 and Step 9 and the work-reference Commit paragraph name the command and keep the hand-built manifest valid. The hand-back merge needed no conflict resolution and no seam edits.

After review (re-merge `fe788bf4`, builder-branch commit `e09c25d8` by the integrator): `--emit <path>` is required, so the command never finalizes by itself and the `advance` phase gate always applies (review F1); the working REQ is resolved through `requeststate.ResolveTarget` plus a working-tree check (F2); the staged-index preflight runs before the journal lookup, so an index refusal creates nothing under `.git` (F4); `work.md` Step 8 item 5, the work-reference Commit paragraph and the `lessons-do-work-cli.md` module bullet now say which inputs are judged and which fields the command fills (F3).

## Decisions

(from the builder hand-back, verbatim; D-01 to D-10 are in `## Exploration`)

- D-11 (DECIDE & STATE): Message-file and failure-error-file handling: trailing `\r` and `\n` are trimmed, everything else kept byte for byte. This matches what a hand-built manifest carries and git strips trailing newlines anyway.
- D-12 (DECIDE & STATE): The dry core still calls `journalLocations`, which creates the empty Git-private `.git/do-work-finalization/` directory, exactly as every `finalize --manifest` attempt (including refused ones) already does. No journal file and no payload directory are created; in dry mode the core blanks the Git-private payload path so its cleanup calls cannot reach it. A non-creating lookup would need a new helper for no observable gain.
- D-13 (DECIDE & STATE): The core returns the required commit paths alongside the "commit_paths omits planned lifecycle or release targets" refusal. The draft pass treats a non-nil set as "the planners passed" and ignores that one refusal; every other refusal returns a nil set and becomes the auto-mode refusal. The second pass with the filled manifest runs the full check, so nothing is skipped. No typed error was added.
- D-14 (DECIDE & STATE): The draft manifest carries `--extra-path` values as its commit_paths, so the release shipped-change guard (which reads commit_paths under `primary_commit`) sees the implementation paths in the draft pass as well.
- D-15 (DECIDE & STATE): The draft pass passes nil manifest bytes; their digest can never equal an existing journal's, so the dry core can never resume a journal (preflight (a) already refuses first).
- D-16 (DECIDE & STATE): Up-front judged-input checks cover presence and transition pairing only (`--transition` complete|fail, `--terminal-status` for complete only, `--failure-type` and `--failure-error-file` for fail only, `--message-file`, `--provenance`, `--implementation-hash` with `supplied_commit`). Value rules (terminal status words, failure type words, hash shape, primary_commit forbids a hash) stay in `validateManifest`, which runs on the filled manifest; its refusal is also OutcomeRefused with code `FINALIZATION-AUTO-INPUT`.
- D-17 (DECIDE & STATE): Refusal codes: `FINALIZATION-AUTO-INPUT` (missing or invalid judged input), `FINALIZATION-JOURNAL-UNFINISHED` (next argv `do-work-cli recover-finalization`), `FINALIZATION-INDEX-NOT-EMPTY` (staged paths in affected paths and evidence), `FINALIZATION-REQUEST-NOT-WORKING` (zero or several working files), `FINALIZATION-PREPARE-REFUSED` (dry core refusal, carries e.g. `RELEASE-PREIMAGE-STALE`). Success with `--emit` adds info finding `FINALIZATION-MANIFEST-EMITTED` with the path and commit_paths. Internal faults (git or discovery errors, emit write failure `FINALIZATION-EMIT`) stay `commandFailure`, exit 2, as today. The non-emit real core refusing after a clean dry pass keeps `FINALIZATION-PREPARE` failure, identical to `finalize --manifest`.
- D-18 (DECIDE & STATE): `--emit` overwrites an existing file at the path (mode 0o600). The action passes a fresh scratch path; refusing on existence would add a check nobody asked for.
- D-19 (DECIDE & STATE): Routing is `slices.Contains(arguments, "--auto-manifest")` in `handleFinalize`; the auto parser refuses `--manifest` as mutually exclusive (usage failure, exit 2). `parseFinalizeArguments` is unchanged.
- D-20 (integrator, DECIDE & STATE): The pre-dispatch `## Scope` list held bare paths, which the scope-drift gate does not read as declared paths (every file reported `SCOPE-UNDECLARED-TOUCH`). The integrator wrapped each of the eight paths in backticks; the declared set did not change. The integrator also stamped `builder_handback_at` from the builder commit's committer date (13:33:33Z) and wrote `integration_at` and this REQ's `## Implementation Summary` and `## Qualification` with the merged REQ-659 writers.
- D-21 (integrator, after review, DECIDE & STATE): Review F1 showed that `finalize --auto-manifest` without `--emit` archived a REQ that `advance` refused (a Triage-only REQ), and `work.md` now names the command, so dropping one flag reached an ungated finalize. The REQ text (requirements 4 and 5) implies a direct mode; the integrator challenged that and made `--emit` required instead, because the work loop always emits and then runs `advance --finalization-manifest` (pre-dispatch D-01), and `finalization` cannot import the `advance` phase check (import direction). Requirement 5's archived path and commit hash are printed by that `advance` call. F2 (reuse `ResolveTarget`), F3 (state the judged/mechanical split) and F4 (index check first) were fixed in the same builder-branch commit `e09c25d8` and re-merged with the same `<pre>` as `fe788bf4`. Correction to D-16 (review F5): a bad `--terminal-status` or `--extra-path` refuses `FINALIZATION-PREPARE-REFUSED`, not `FINALIZATION-AUTO-INPUT`, because the draft core runs before `validateManifest`; both exit 1. Correction to D-12 (review F4): it is no longer true that an index refusal creates `.git/do-work-finalization/`. F5, F6 and F7 stay report only. After the fix, D-17's `FINALIZATION-REQUEST-NOT-WORKING` also carries `ResolveTarget`'s not-found, ambiguous and identity-mismatch codes in its evidence, and the direct-mode `FINALIZATION-PREPARE` failure no longer exists (re-review N2).

## Discovered Tasks

(from the builder hand-back; impact tokens added by the integrator per `actions/review-work.md` Step 10)

- **impact-negligible** `prepareManifestJournal` in non-dry mode still runs `os.RemoveAll(payloadDirectory)` on early refusals before any payload was adopted, so a stray `REQ-N.payloads` folder with no journal is deleted by any refused `finalize --manifest` (pre-existing behaviour, unchanged) → report only
- **impact-negligible** `lifecycleadvance/queue_commands.go:424` `queueWriterLabel` is a near copy of the writer default that differs on purpose-or-drift: it substitutes `unknown-host` for an empty hostname, while `requeststate.DefaultWriterLabel` keeps `":"+root`. `lifecycleadvance` already imports `requeststate`, so the two could be unified; left untouched per the brief → report only

## Qualification

**Gate records (`advance --diff-range 84d56197..457608cc`):** `qualify` satisfied (success), `scope-drift` satisfied (success). Two `QUALIFY-NEW-FILE-UNWIRED` warnings, judged false positives: `finalization_auto_manifest.go` is a same-package Go file whose `handleAutoManifest` is called from `handleFinalize` in `finalization_commands.go`, and `finalization_auto_manifest_test.go` is a `_test.go` file the Go test runner finds by convention. No debug-artifact, P-A-U or output-primitive findings. The first run reported `SCOPE-UNDECLARED-TOUCH` for all eight files because the pre-dispatch `## Scope` list held bare paths; the integrator wrapped each path in backticks (no path added or removed, D-20), after which scope-drift passed.

**Scope:** declared `write_set` (8 paths) equals the touched set in `git diff 84d56197..457608cc --stat` (8 files, 474 insertions, 32 deletions). No `do-work/` path on the builder branch (the queue guard printed nothing before the merge). The merge auto-resolved `work.md`, `work-reference.md` and `lessons-do-work-cli.md` against REQ-659's and REQ-660's landed edits with no conflict; the merged tree builds (`go build ./...`, `go vet` on both packages).

**Requirement trace (read against the merged files):**
1. Flag set: `parseAutoManifestArguments` accepts `--auto-manifest`, `--transition`, `--terminal-status`, `--message-file`, `--provenance`, `--implementation-hash`, `--release-manifest`, `--failure-type`, `--failure-error-file`, `--writer`, `--emit`, and repeated `--extra-path`; `handleFinalize` routes there when `--auto-manifest` is present.
2. Mechanical fields: the request path is the one `RequestsByID[id]` entry with `TreeSection == "working"` (zero or several refuse `FINALIZATION-REQUEST-NOT-WORKING`); both digests come from the live REQ and `do-work/CHECKPOINT.md` bytes; `completed_at` uses the same `time.Now().UTC().Truncate(time.Second)` format as `do-work-cli now`; `release_at` equals `completed_at` when `--release-manifest` is given; the writer defaults to `requeststate.DefaultWriterLabel` (the expression `planCheckpoint` used before, now shared); `commit_paths` is the required set the core returns plus each `--extra-path`, normalized.
3. Preflight before any write: unfinished journal (`FINALIZATION-JOURNAL-UNFINISHED`), staged index (`FINALIZATION-INDEX-NOT-EMPTY`, names the paths), then the dry core, which surfaces a stale release version as `RELEASE-PREIMAGE-STALE` inside `FINALIZATION-PREPARE-REFUSED`. "Nothing staged outside commit_paths" is subsumed by the empty-index rule the finalizer already enforces (pre-dispatch D-03).
4. `--emit`: the filled manifest passes `validateManifest` and a second dry core pass, then is written; no journal. The test then runs `finalize --manifest` on the emitted file to `cleanup_complete`. `advance --finalization-manifest` binds through `FinalizeBound` → `prepareBoundJournal` → the same `prepareManifestJournal`, so it reads the same bytes the same way.
5. Output: the existing text renderer prints outcome, findings, and each finalization record's phase, archive path and commits (pre-dispatch D-07); emit mode adds the info finding `FINALIZATION-MANIFEST-EMITTED` with the path and `commit_paths`.
6. Exit codes: refusals are `OutcomeRefused` (exit 1, asserted by the tests through `resultmodel.ExitCode`); usage and internal faults stay `commandFailure` (exit 2), as for `finalize --manifest`.
7. Judged fields never invented: `judgedInputProblem` refuses a missing transition, terminal status, failure pair, message file, provenance or `supplied_commit` hash.
8. Prose: `work.md` Step 8 item 5 names `finalize --auto-manifest REQ-NNN` with its inputs and `--emit`, Step 9 accepts "the Step 8 `--emit` file or a hand-built one", and the work-reference Commit paragraph names the command and keeps the hand-built manifest valid. The `core-checks.sh` tokens ("author exactly one strict manifest", "Pass that one manifest to the current `advance` continuation.") are kept; the gate confirms it.
9. Release: handled at finalization (patch bump, own changelog entry; this REQ closes UR-145).

**Constraints:** no new frontmatter field or status; `parseFinalizeArguments` and the `finalize --manifest` path are unchanged apart from calling the shared core, whose non-dry behaviour is identical (same checks in the same order, same journal write); `internal/lifecycleadvance/` is untouched.

**Red-Green Proof trace:** `TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts` (emit, then `finalize --manifest` accepts the file), `TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal` (exit 1, names `unrelated.txt`, no journal, payload directory or emit file), `TestFinalizeAutoManifestRefusesWithoutAMessageFile` (exit 1, names `--message-file`), `TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting` (`RELEASE-PREIMAGE-STALE`, nothing written).

**Remediation re-merge (review F1-F4), cumulative range `84d56197..fe788bf4`:** `advance` no longer takes qualify input at this phase, so the integrator ran the same `qualify --request-path <P> --diff-range 84d56197..fe788bf4` handler directly: success, only the two judged `QUALIFY-NEW-FILE-UNWIRED` warnings above. The touched set is still the 8 declared files (`git diff 84d56197..fe788bf4 --stat`: 477 insertions, 32 deletions); the delta touches `finalization_auto_manifest.go`, its test, `work.md`, `work-reference.md` and `lessons-do-work-cli.md`. The queue guard before the re-merge printed nothing. Requirement 4 and 5 now read: `--emit` is required and the emitted file goes to `advance`, which prints the outcome, phase, archive path and commits (D-21).

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `457608cc` (machine quiet before launch: 1-minute load 4.65, no other gate running; load 10.26 at the end), then `advance REQ-658 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-658-probe.sh`.
**Result:** ✓ Repository gate passed on the first run (exit 0, gate wall 173 s, do-work-cli 897 tests in 78 s, slowest file `internal/finalization/finalization_recovery_test.go` 27.37 s under the 30 s limit; queue-kanban 421 tests, slowest 21.45 s). Probe exit 0. Gate records `green-gate`, `scope-drift` and `run-blocked-check` satisfied.

**After the review-fix re-merge (`fe788bf4`):** the same gate argv ran again (load 2.44 before launch, 8.24 at the end): exit 0, gate wall 153 s, 897 do-work-cli tests, slowest file `internal/finalization/finalization_recovery_test.go` 24.11 s. `advance` refuses gate input past this phase, so the green record was written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` at `fe788bf4`, and `REQ-658-probe.sh` was run directly: exit 0 ("4 tests, both action files name the command").

**Red-green validation:** (from the builder hand-back, traced to `## Red-Green Proof`)
- `TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts`, `TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal`, `TestFinalizeAutoManifestRefusesWithoutAMessageFile`, `TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting` (`internal/finalization/finalization_auto_manifest_test.go`): ✗ before the production edit (each `Outcome:"failure"`, `FINALIZATION-USAGE`, `unknown finalize option "--auto-manifest"`) → ✓ after (1.20 s, 0.09 s, 0.04 s, 0.16 s).
- `TestFinalizeAutoManifestRefusesWithoutAMessageFile`, new assertion for review F1 (no `--emit` must refuse and write nothing): ✗ with the `--emit` check disabled (`missing --emit result = ... Outcome:"failure"`) → ✓ with it (exit 0).
- Existing `finalize --manifest` and `FinalizeBound` tests in `internal/finalization/` pass unchanged after `prepareBoundJournal` was split into the shared core.

**New tests added:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go` (four tests through `Handlers()`, on the existing REQ-760 planned-release fixture)

**Existing tests updated:** none.

**Heavy verification plan:**
- Range: 84d56197bb78fa5b91f2ece328bbbb21ec8fdcd1..fe788bf4c2b864ace38bc4687fa0ac59dfeff7b9 (re-planned after the review-fix re-merge; same four lanes as the first plan at `457608cc`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 88%** | 2026-10-10T15:18:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 93% |
| Code Quality | 85% |
| Test Adequacy | 78% |
| Scope | 96% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `finalize --auto-manifest` without `--emit` finalizes outside the `advance` phase gate. Reproduced: a Triage-only REQ was archived as completed while advance refused the same manifest (`finalization_auto_manifest.go:200-206`, named at `work.md:462`). Fix: refuse without `--emit`, or apply the advance phase check. — impact-user-visible → report only

**Minor findings:**
- F2 second REQ-id resolver (`finalization_auto_manifest.go:148-156`) instead of `requeststate.ResolveTarget` + a `TreeSection` check; the outcome converges, only the refusal code differs — impact-negligible → report only
- F3 `work.md:462` and `work-reference.md:719` list the mechanical fields as authored and then as filled; Step 8 item 5 grew from 486 to 1080 characters; the judged/mechanical split is never stated — impact-negligible → report only
- F4 preflight (a) creates `.git/do-work-finalization/` before the index refusal that says "nothing was written"; `finalize --manifest` does not; run the index check first — impact-negligible → report only
- F5 a bad `--terminal-status` or `--extra-path` refuses `FINALIZATION-PREPARE-REFUSED`, not `FINALIZATION-AUTO-INPUT` as D-16 says; the 7 new codes have no reader — impact-negligible → report only
- F6 no repository test for `advance --finalization-manifest <emitted>`, non-emit mode, `supplied_commit`, `fail`, or the journal and not-working refusals; add one lifecycleadvance emit-then-advance test — impact-negligible → report only
- F7 inherited from REQ-660 (worktree lifecycle command) F3: `handoff.go:87`, `cleanup_git.go:345`, queue-kanban `verify.go` and Crash Recovery still read owned worktree links as dirt; untouched by this merge — impact-user-visible → report only

**Acceptance:** Pass — implementation and integration stages: focused finalization tests ok; the emitted manifest passes `finalize --manifest` and (scratch copy) `advance --finalization-manifest` for primary+release, supplied and already-green fixtures; 7 refusals exit 1; deployment and live acceptance unassessed
**Restatement sweep:** redefined (E1) how the finalization manifest is built (`finalize --auto-manifest ... --emit`), (E2) which manifest fields are mechanical vs judged, (E3) `work.md` Step 8 item 5 / Step 9 and `work-reference.md` Commit wording. Stale: `work.md:462` and `work-reference.md:719` first clauses still frame mechanical fields as action-authored (F3). Consistent: `work.md:256`, `:467`, `:500`; `work-reference.md:352`, `:388`, `:715`; `lessons-do-work-cli.md:127`; `advance_commands.go:264-265` placeholder `<action-authored-finalization-manifest>` (kept by D-10, still valid); no shipped text describes finding `commit_paths` by refusal; `docs/` and `skills/do-work/docs/` carry no `finalize` usage. Inherited: REQ-659 (frontmatter set and req append-section) F4 `work-reference.md:61` and F5 `:73` fixed and still fixed; REQ-660 (worktree lifecycle command) F1 and F2 (`fan-out-reference.md:78`) fixed and still fixed, F3 dirt readers still stale (F7); this merge made none of them stale again (it touched only `work.md:462`, `:467`, `work-reference.md:719`, `lessons-do-work-cli.md:127`)
**Suggested testing:** 5 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*

Full report: `do-work/runs/work-2026-10-10-131527/REQ-658-review.md`.

### Re-review (delta 457608cc..fe788bf4) | 2026-10-10T15:26:38Z

F1 to F4 were fixed on the builder branch and re-merged with the same `<pre>` (D-21).

- Overall: 93% | Acceptance: Pass (focused AutoManifest tests ok; core-checks.sh exit 0; scratch tests prove missing `--emit` → `FINALIZATION-AUTO-INPUT` exit 1 with nothing written, the index refusal creates no `.git/do-work-finalization/`, and a duplicate id refuses `REQUEST-AMBIGUOUS`)
- F1: closed (`--emit` is required at `finalization_auto_manifest.go:94`, the non-emit finalize branch and the `context` import are deleted, and the new assertion fails under the old code). F2: closed (`ResolveTarget` + working-section check at `:150-157`). F3: closed (`work.md:462`, `work-reference.md:719` and `lessons-do-work-cli.md:127` state the judged/mechanical split and "never finalizes", Step 9 is consistent, no other shipped text describes the direct mode). F4: closed (the index check runs before `journalLocations`, proven in a scratch repo).
- N1 F4 is not pinned by a repository test, because `assertNothingWritten` creates `.git/do-work-finalization/` itself through `journalLocations` — impact-negligible → report only
- N2 `FINALIZATION-REQUEST-NOT-WORKING` now wraps not-found, ambiguous and identity-mismatch refusals and drops their paths; REQ decision D-17 still describes the old resolver and the deleted non-emit refusal — impact-negligible → report only
- N3 `work-reference.md:719` judged-input list omits the failure type and failure error (`work.md:462` includes them) — impact-negligible → report only

Full re-review report: `do-work/runs/work-2026-10-10-131527/REQ-658-rereview.md`.

## Lessons Learned

**What worked:** Splitting `prepareBoundJournal` into decode plus one shared core with a dry mode meant `--emit` validates through the same checks `finalize --manifest` and `advance` run, with no second copy. Letting the core return the required commit paths alongside its "commit_paths omits" refusal let a draft manifest learn the set in one pass. The reviewer proved the emitted file through `advance --finalization-manifest` in three fixtures, and this REQ's own finalization used the merged command.
**What didn't:** The first cut kept a direct mode (no `--emit`) that finalized through the journal without the `advance` phase gate, so one dropped flag archived a Triage-only REQ (review F1, fixed by requiring `--emit`). The working-REQ lookup was a second resolver beside `ResolveTarget` (F2, fixed). The pre-dispatch `## Scope` list used bare paths, which the scope-drift gate does not read as declared paths.
**Worth knowing:** A command that builds input for a gated step must not also perform that step: the gate lives in the caller (`advance`), and the builder package cannot import it. Write `## Scope` paths in backticks, or scope-drift reports every touched file as undeclared.

## Orientation

Now the finalization manifest's mechanical fields are built by `do-work-cli finalize --auto-manifest REQ-NNN --emit <path>` in `internal/finalization/finalization_auto_manifest.go`, which runs the shared planner core `prepareManifestJournal` (`finalization_prepare.go`) in dry mode; the action still judges transition, message, release payload, extra paths and provenance, and passes the emitted file to `advance --finalization-manifest` (`actions/work.md` Step 8 item 5 and Step 9, `actions/work-reference.md` Commit & Metadata-Commit Procedure). [MAP CHANGED]: one new `finalize` mode and one shared prepare core. Prime spot-check: `prime-shell-commands.md`, `prime-action-files.md` and `prime-releases.md` name no path this change moved or removed.

## Heavy Verification Plan

- Base: 84d56197bb78fa5b91f2ece328bbbb21ec8fdcd1
- Target: fe788bf4c2b864ace38bc4687fa0ac59dfeff7b9
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: fe788bf4c2b864ace38bc4687fa0ac59dfeff7b9
- Execution revision: fe788bf4c2b864ace38bc4687fa0ac59dfeff7b9 (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-658`, `QUEUE_KANBAN_BROWSER` set, removed afterwards)
- do-work-cli-integrations: exit 0, executed, 81 s
- staged-skills: exit 0, executed, 41 s
- updater: exit 0, executed, 73 s
- installer: exit 0, executed, 35 s
- Earlier drain at the first merge `457608cc` (before the review fix): all four lanes exit 0, executed (80, 44, 71, 42 s).

## Timing

Observed 2026-10-10T15:05:27Z to 2026-10-10T15:26:39Z: 21m 12s total, 25m 34s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 14m 10s | 4 |
| review | 10m 56s | 2 |
| handback-merge | 28s | 1 |

Slowest stage: review / independent review of 84d56197..457608cc, 7m 18s, outcome success.

Note: no builder-work event was recorded. The hand-back had landed before this integrator started (dispatch 13:25:57Z, builder commit 13:33:33Z, about 8 minutes of build), so per fan-out-reference.md → Landed hand-back the event is skipped rather than charging the builder with the wait. Review and verification overlap in wall time because each heavy drain ran while a reviewer worked, so attributed time exceeds the observed span.
