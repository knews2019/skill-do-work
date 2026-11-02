---
id: REQ-709
title: 'recover ignores loose do-work files that no request record can own, so do-work commit can commit inbox notes'
status: completed
created_at: 2026-11-02T19:41:34Z
user_request: UR-162
domain: backend
prime_files: []
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-11-02T19:46:07Z
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-02T19:49:44Z
  basis:
    - Route A
    - 4-file write set
    - 4 acceptance criteria
dispatch_at: 2026-11-02T19:51:31Z
builder_handback_at: 2026-11-02T20:01:32Z
integration_at: 2026-11-02T20:02:04Z
review_at: 2026-11-02T20:13:24Z
kb_status: pending
completed_at: 2026-11-02T20:20:23Z
commit: 6e2ae9a4e770050c002c75377b520f7c26f310bb
release_at: 2026-11-02T20:20:23Z
---
# Recover Ignores Loose do-work Files That No Request Record Can Own
## What
The recovery check (`recover`, which `do-work commit` and the work loop run first) refuses with `FINALIZATION-DISCOVERY-AMBIGUOUS` when the only uncommitted file is a note under `do-work/` that no request record can own, such as an untracked `do-work/inbox/<note>.md`. It should leave such a file alone, so `do-work commit` can commit it.
## Why
`do-work commit` stops on the refusal, so it can never commit an inbox note, which is exactly the leftover file it exists for.
## Detailed Requirements
- A file under `do-work/` outside the folders that hold request records (the queue, working, archive and user-request folders, REQ reservation markers, and run records) is not finalization evidence. Recovery neither refuses on it nor associates it with a request.
- Key the rule on that condition, stated once in code: where request records live. Do not list note folders such as `inbox/`. The user chose this scope so that every note (inbox notes, handoff notes such as `do-work/HANDDOWN-*.md`) is covered.
- Real unfinished finalization work next to such a note still refuses or settles exactly as today.
- `do-work commit` then groups the note as an ordinary change with no owning request.
## Context
Seen twice on 2026-11-02. During the UR-161 run, `recover` refused on the untracked note `do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md`; the run agent hid it with a temporary `.git/info/exclude` line. Then `do-work commit` refused on the same note and committed it by judgment (commit `14878ab6`). The refusal is raised in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`, where the note lands in the unresolved shared-state list. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## Red-Green Proof
**RED prompt/case:** A test repository whose only uncommitted file is an untracked `do-work/inbox/note.md`. Run `recover`, then the `do-work commit` preflight.
**Why RED now:** `recover` refuses with `FINALIZATION-DISCOVERY-AMBIGUOUS` and names the note, so `do-work commit` stops.
**GREEN when:** `recover` returns success, and `do-work commit` commits the note as a group with no owning request. A real unfinished request next to the note is still handled exactly as today.
**Validation:** User confirmed.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23666 tokens per `do-work/lessons-index.md` at claim, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers semantic recovery completeness and condition-based classifiers; families `closed-enumeration-for-a-condition`, `observed-subset-is-not-semantic-completeness` match this work.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Narrow `sharedFinalizationPath` to the condition "where request records live", keyed on the first path segment under `do-work/`; add `unownedDoWorkNote` and skip such paths (path and rename origin judged alone) when discovery builds its `dirty` set. Every other user of the predicate, including the apply-time check in `finalization_apply.go`, inherits the narrowing. No new package, no exported API, refusal codes and messages unchanged. (Builder hand-back, transcribed by the orchestrator.)
- [x] **[APPLY]:** Three files: `finalization_discovery.go`, the new `finalization_unowned_notes_test.go`, and `recovery_commands_test.go`. No `do-work/` path written in the worktree; `corehelpers/inventory.go` untouched (D-02).
- [x] **[UNIFY]:** `git diff --stat main...HEAD`: 3 files, 179 insertions, 3 deletions. `gofmt -l` empty, `go vet` clean on both packages. Checked `finalization_discovery.go` (whole diff, every caller of the predicate, doc comment states the condition, no debug output), `finalization_unowned_notes_test.go` (one-line failure comment per test, existing helpers only, about 2.5 s), `recovery_commands_test.go` (calls `handleRecover` with no flags as `actions/commit.md` Step 1 does, asserts the inventory rows exactly).
*Source: "recover should not treat untracked do-work/inbox notes as unfinished finalization work, so do-work commit can commit them"*

## Triage

**Route: A** - Simple

**Reasoning:** A bug fix with a clear reproduction that names the file: `recover` refuses `FINALIZATION-DISCOVERY-AMBIGUOUS` on an untracked `do-work/inbox/` note because `sharedFinalizationPath` in `finalization_discovery.go` reads every path under `do-work/` as shared finalization state. Reproduced at claim time with a temporary `do-work/inbox/` note in this repository (refused, the note named as the only blocked path; note deleted after). The fix is one predicate stated once in the finalization package plus focused tests.

**Planning:** Not required

**Open questions (Step 3.5):** none in the REQ. Three design calls are settled by orchestrator judgment and recorded as D-01 to D-03 in `## Decisions`: the lifecycle record files `do-work/CHECKPOINT.md` and `do-work/calibration-log.tsv` stay evidence; the rule lives in the finalization package, not the inventory chokepoint; every other file under `do-work/` outside the record locations (for example `lessons-index.md`, handoff notes) is covered by the same condition.

**Test policy for this run (user decision, 2026-11-02):** focused tests only. The gate argv for every evidence step is `bash do-work/runs/work-2026-11-02-194601/REQ-709-gate.sh`; `_dev/tests/maintainer-verify.sh` and heavy lanes are not run.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_unowned_notes_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modified)

**What was done:** `sharedFinalizationPath` now answers true only where request records live: `do-work/queue/`, `working/`, `archive/`, `user-requests/`, `.req-reservations/`, `runs/`, `CHECKPOINT.md` and `calibration-log.tsv` (plus the bare `do-work` path). Its doc comment states that condition and why. A new `unownedDoWorkNote` helper marks every other path under `do-work/`, and finalization discovery leaves such paths out of its evidence set, judging a rename's two sides one at a time. The apply-time commit-safety check for a discovered journal (`finalization_apply.go`) uses the same narrowed predicate, so a legacy tail now settles around a note instead of refusing. Refusal codes and messages are unchanged. Builder commit `4634893c`, merged as `5a1d9387` (range `c5ae2d8e..5a1d9387`).

**Review fix (builder commit `6582e1d4`, re-merged as `6e2ae9a4` with the same `<pre>`, cumulative range `c5ae2d8e..6e2ae9a4`):** tests only. Two new subtests pin the sole-releaser legacy tail and the claim-only topology beside a note, and both legacy-tail subtests now compare `commit_paths` exactly (review F1, decisions D-10 and D-11).

## Decisions

(D-01 to D-03 are orchestrator design calls given in the builder brief; D-04 to D-11 come from the builder hand-back `do-work/runs/work-2026-11-02-194601/REQ-709-handback.md` and its review-fix addendum.)

- **D-01** (DECIDE & STATE, orchestrator): the condition is "where request records live", stated once in `sharedFinalizationPath`: the queue, working, archive and user-request folders, REQ reservation markers, run records, and the two lifecycle record files every finalization writes, `do-work/CHECKPOINT.md` (claim records) and `do-work/calibration-log.tsv` (one row per finished request). The REQ's own list names the folders but not those two files; they stay evidence because the REQ requires real unfinished finalization work to refuse or settle exactly as today, and a dirty checkpoint or calibration log is that work. No note folder is listed.
- **D-02** (DECIDE & STATE, orchestrator): the rule lives in the finalization package, not in the inventory chokepoint `untrackedDoWorkMetadata` that handles Finder `.DS_Store` files, because `do-work commit` reads that inventory and must still see the note to commit it.
- **D-03** (DECIDE & STATE, orchestrator): discovery neither refuses on such a file nor associates it with a request, even when a legacy Implementation Summary names it. Consequence: other non-record files directly under `do-work/` (for example `do-work/lessons-index.md`) also stop being evidence, so a legacy tail next to a dirty one now settles and leaves it dirty for `do-work commit`, where it refused before. This is the user's chosen scope ("every note"), and no byte is lost: the file stays in the tree.
- **D-04** (DECIDE & STATE, builder): a path whose first segment under `do-work/` is a record location counts as evidence even with no further segment (for example `do-work/runs` as a file), and the bare `do-work` path stays evidence. Conservative reading of "exactly as today".
- **D-05** (DECIDE & STATE, builder): the name `sharedFinalizationPath` is kept so its callers are unchanged; the doc comment carries the new meaning.
- **D-06** (DECIDE & STATE, builder): the finalization tests sit in a new file `finalization_unowned_notes_test.go` (about 2.5 s), and the public test calls the registered `handleRecover` handler as the neighbouring tests do.
- **D-07** (DECIDE & STATE, builder): fixtures use the real UR-161 inbox note filename.
- **D-08** (DECIDE & STATE, builder): tests and fix land in one commit (no red-only commit); the RED text is recorded in `## Testing`.
- **D-09** (DECIDE & STATE, builder): no prose change. No shipped action or doc states the old "every do-work path is shared state" rule, and `actions/commit.md` Step 1 is still true.
- **D-10** (DECIDE & STATE, builder, review fix F1): the claim-only subtest asserts exactly one `FINALIZATION-NONE` finding, which separates "succeeded with nothing to do" from any other success shape.
- **D-11** (DECIDE & STATE, builder, review fix F1): the sole-releaser subtest uses `seedSoleReleaserLegacyTail`, the tail that only settles under `--assume-sole-releaser`, so it reaches the `soleReleaserSharedPath` branch; its expected attributed set is copied from the existing REQ-499 test.

## Discovered Tasks

(from the builder hand-back; its "impact-low" is restamped to the canonical `impact-negligible`)

- **impact-negligible** The record-location list in `sharedFinalizationPath` repeats the `do-work/` layout that `internal/repositorymodel/repository_model.go`'s walker encodes with its own literals; there is no shared declaration. A record folder added to the walker but not here would be treated as a note, and recovery would leave it dirty instead of refusing. The eight-row table in `TestRecoverFinalizationStillRefusesRecordStateNextToANote` pins today's set. → report only

## Qualification

**Gate records:** `advance --diff-range c5ae2d8e..5a1d9387` → `qualify` satisfied (no debug artifact, no output-primitive finding, P-A-U boxes ticked from the hand-back). One warning, `QUALIFY-NEW-FILE-UNWIRED` on `finalization_unowned_notes_test.go`: judged an entry-point exception, because `go test` discovers `*_test.go` files by name and the file's two tests ran (RED and GREEN) in the builder's and the orchestrator's runs. Route A has no Scope, so no scope-drift record.

**Requirement trace (read against `git diff c5ae2d8e..5a1d9387` and the merged file):**
- A file under `do-work/` outside the record folders is not finalization evidence, and recovery neither refuses on it nor associates it: `discoverFinalizationJournals` skips every `unownedDoWorkNote` path (and rename origin) when it builds `dirty`, so the note reaches neither the shared loop, the ordinary-owner loop, `coherentClaimOnlyTopology` nor `ambiguousSharedRemainder`; the apply-time check in `finalization_apply.go` uses the narrowed predicate. Pinned by `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns` (no candidate, legacy tail with an untracked inbox note, legacy tail with a modified tracked handoff note that its own Implementation Summary names: the note is in no `commit_paths` and keeps its porcelain status).
- Keyed on the condition, stated once: `sharedFinalizationPath` is the only place the record locations appear, as the first path segment under `do-work/`; no note folder is named. `do-work/working-notes.md` (a prefix sibling of `working/`) is correctly a note.
- Real unfinished finalization work still refuses exactly as today: `TestRecoverFinalizationStillRefusesRecordStateNextToANote` (a foreign checkpoint claim line, and an untracked file in each of the eight record locations, each next to the inbox note) refuses `FINALIZATION-DISCOVERY-AMBIGUOUS` naming only the record path. The existing Finder-metadata and unowned-refusal tests still pass.
- `do-work commit` groups the note as an ordinary change: `TestRecoverLeavesAnUntrackedInboxNoteForCommit` runs the public `recover` (the commit preflight) on a repository whose only uncommitted file is the inbox note, gets success, and reads the note back from `corehelpers.ReadProtectedInventory` as an `A` row. Association was read in code, not tested: the commit association check never claims a `do-work/` path, so Step 4 groups it as unassociated.

**Judgment:** substantive and small (36 production lines changed, 143 test lines), matches D-01 to D-09; the hand-back's "Klarna check" in its UNIFY note is the Klarna Test from `crew-members/shared-principles.md` (check the user's intended outcome), applied correctly (review F5 corrected this record, which first called it a stray word). Qualified.

## Testing

**Tests run:** focused repository gate `bash do-work/runs/work-2026-11-02-194601/REQ-709-gate.sh` at merge `5a1d9387` through `run-timed-command` (2026-11-02T20:03:31Z to 20:05:23Z, 112 s, 1-minute load rising to 19 with the reviewer's test runs alongside): `gofmt -l`, `go vet` and `go test -count=1` on `./internal/finalization/` (71.8 s) and `./internal/lifecycleadvance/` (33.0 s), then `_dev/tests/shipped-package-reference-contract.sh`, `contract-regressions.sh`, `action-shell-blocks.sh`, `prescribed-shell-canonicalization.sh`, `suite-manifest-contract.sh`. Then `advance --gate-arg bash --gate-arg do-work/runs/work-2026-11-02-194601/REQ-709-gate.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-11-02-194601/REQ-709-probe.sh`: `green-gate` and `run-blocked-check` satisfied (probe status 0: the three new tests plus the Finder-metadata and unowned-refusal tests PASS).
**Result:** ✓ All passing (gate exit 0 on the first run; both packages pass in full; the new finalization test file runs in about 2.5 s, under the 30 s budget)

**Red-green validation:**
- `TestRecoverLeavesAnUntrackedInboxNoteForCommit` (`lifecycleadvance/recovery_commands_test.go`), the `## Red-Green Proof` RED case: ✗ before (`recover refused on a note no request record owns`: `Outcome:"refused"`, reason `FINALIZATION-DISCOVERY-AMBIGUOUS`, blocked path the inbox note) → ✓ after (public `recover` succeeds and the note is still an `A` inventory row for `do-work commit`).
- `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns` (`finalization/finalization_unowned_notes_test.go`, 3 subtests): ✗ before (each refused `FINALIZATION-DISCOVERY-AMBIGUOUS` naming the note; the no-candidate branch also named `do-work/working-notes.md`) → ✓ after (success; a legacy tail settles to `cleanup_complete` and the note keeps its `??` or ` M` status and is in no `commit_paths`). The legacy-tail subtest is the branch this repository hit, because its archive holds legacy candidates.
- `TestRecoverFinalizationStillRefusesRecordStateNextToANote` (9 subtests), the GREEN clause "a real unfinished request next to the note is still handled exactly as today": ✗ before (each refusal listed the record path and the note) → ✓ after (each still refuses `FINALIZATION-DISCOVERY-AMBIGUOUS` and names only the record path).
- Builder differential: with only the narrowed predicate, the handoff subtest fails (the note joins REQ-750's group); with only the `dirty` filter, both legacy-tail subtests fail at the apply-time check with `FINALIZATION-AMBIGUOUS-SHARED-STATE`. Both halves are needed.

**New tests added:**
- `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns`, `TestRecoverFinalizationStillRefusesRecordStateNextToANote` in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_unowned_notes_test.go`
- `TestRecoverLeavesAnUntrackedInboxNoteForCommit` in `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go`

**Heavy verification plan:** heavy lanes not run: user has not granted heavy-lane permission for this run; focused tests only (2026-11-02). No heavy hold is written. Re-planned on the cumulative range after review fix F1; the planner would have selected:
- Range: c5ae2d8ecbd8ec9fec263407480311d4b0cac51a..6e2ae9a4e770050c002c75377b520f7c26f310bb
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go matched subtree skills/do-work/tools/do-work-cli (and 2 more paths)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go matched subtree skills (and 2 more paths)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go matched subtree skills/do-work/tools/do-work-cli (and 2 more paths)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go matched subtree skills/do-work/tools/do-work-cli (and 2 more paths)

**Repository gate after review fix F1:** the same focused gate argv at `6e2ae9a4` through `run-timed-command` (2026-11-02T20:16:56Z to 20:18:15Z, 78 s, 1-minute load 6.1) exit 0: finalization 51.0 s, lifecycleadvance 24.5 s, all five contract scripts green; recorded with `record-green-gate` at `6e2ae9a4`. `REQ-709-probe.sh` exits 0 there. The two new subtests were RED on the pre-change discovery file (sole-releaser refused naming only the note; claim-only refused naming the checkpoint, the note and both REQ-699 paths) and GREEN on the branch; the reviewer re-confirmed both.

*Verified by work action*

## Review

**Overall: 95%** | 2026-11-02T20:13:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. Claim-only topology and `--assume-sole-releaser` changed from refusing to proceeding when a note sits beside them (correct, verified old vs new in scratch), but no test pinned either entry; the legacy-tail subtests also checked only the phase and the note's absence, not exact `commit_paths`. — impact-negligible → report only (closed in `6582e1d4` / merge `6e2ae9a4`: sole-releaser and claim-only subtests plus exact `commit_paths`, RED-GREEN re-confirmed)
- F2. The record-location list in `sharedFinalizationPath` (`finalization_discovery.go`) repeats the repository walker's layout (`repository_model.go`) with no shared declaration. A record folder added to the walker alone would be treated as a note, so its dirt would pass recovery (fail-open). Agrees with the builder's Discovered Task, restamped. — impact-rule-change → report only
- F3. A REQ-shaped file outside the record folders (for example an untracked `do-work/REQ-800-stray.md`) no longer stops `recover`, and `do-work commit` would commit it as an unowned change. This follows the REQ's location rule and D-03; `doctor` still reports `STRAY-REQUEST`. — impact-negligible → report only
- F4 (Nit). `owners` in `discoverFinalizationJournals` now includes `do-work/` note paths from Implementation Summaries; keeping a note out of a request's group rests on the `dirty` filter alone (pinned by the handoff subtest and its exact `commit_paths`). — impact-negligible → report only
- F5 (Nit). `## Qualification` first called the hand-back's "Klarna check" a stray word; it is the Klarna Test that `crew-members/shared-principles.md` requires, applied correctly. — impact-negligible → report only (corrected in Qualification)
**Acceptance:** Pass — implementation and integration: focused tests and the full `finalization` and `lifecycleadvance` packages green, RED confirmed on pre-change code by builder and reviewer, and a CLI built from the merge on a scratch repository shows `recover` succeeding with an untracked inbox note and `protected-inventory associate` reporting the note unowned (the pre-change binary refuses). Deployment and live acceptance unassessed.
**Restatement sweep:** redefined which `do-work/` paths are finalization evidence (`sharedFinalizationPath`). Swept `skills/` prose and Go comments, `README.md`, `skills/*/docs/`. No stale restatement.
**Suggested testing:** 3 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

**Re-check (`5a1d9387..6e2ae9a4`, recorded 2026-11-02T20:19:30Z):** F1 closed, tests-only range, no new finding; overall 93% → 95%, Test Adequacy 80% → 90%. Full report and re-check: `do-work/runs/work-2026-11-02-194601/REQ-709-review.md`.

**Orchestrator disposition:** a Pass at 93% needs no remediation. F1 was fixed because the REQ's core promise is "exactly as today" at every entry, and the prime's `guard-on-one-entry-path` trap asks for each entry to be pinned; it cost one test-only commit, one re-merge with the same `<pre>` and one focused gate. F2 to F4 stay report only (none is impact-critical, so none is queued).

## Lessons Learned

**What worked:** Reproducing the refusal at claim time with a temporary inbox note in this repository showed which discovery branch it hit (legacy candidates exist here, so the note fell into the shared loop, not the no-candidate remainder); the tests cover both branches. Testing each half of the fix alone (predicate only, filter only) showed both are needed: the filter keeps a note out of a request's commit group, the narrowed predicate clears the apply-time commit-safety check.
**What didn't:** The first legacy-tail tests asserted only the phase and the note's absence, so "settles exactly as today" was not pinned, and two entry points (claim-only topology, `--assume-sole-releaser`) that changed outcome had no test; the review caught it (F1).
**Worth knowing:** The 0.266.4 Finder fix filters hidden files at the inventory chokepoint, which every inventory reader shares; that was right for metadata, but a note must stay visible to `do-work commit`, so this rule lives only in finalization discovery. The record-location list in `sharedFinalizationPath` duplicates the repository walker's layout literals (F2): a new record folder must be added to both, or its dirt passes recovery. A stray REQ-shaped file outside the record folders now passes `recover` too (F3); `doctor` still reports it.

## Orientation

Now `recover`, and with it `do-work commit` and the work loop's queue boundary, leaves alone any file under `do-work/` that no request record can own (an inbox note, a handoff note), so `do-work commit` can commit it; unfinished finalization work in a record location still refuses or settles as before. Lives in the do-work-cli finalization recovery subsystem (`internal/finalization`, `sharedFinalizationPath` in `finalization_discovery.go`). `prime_files` is empty; spot-check of `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`: every path it names still exists and it does not restate which `do-work/` paths are finalization evidence; not stale.

## Timing

Observed 2026-11-02T19:46:07Z to 2026-11-02T20:18:15Z: 32m 08s total, 32m 21s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 10m 09s | 1 |
| builder-work | 10m 02s | 1 |
| planning | 5m 28s | 1 |
| remediation | 3m 26s | 1 |
| verification-gate | 3m 10s | 2 |
| handback-merge | 6s | 1 |

Slowest stage: review / independent review with restatement sweep, 10m 09s, outcome success.
Slowest command: verification-gate / focused repository gate, 1m 52s, exit 0, bash (2 argv tokens).
