---
id: REQ-560
title: '[impact-rule-change] Hand-back and finalize check cleanliness only on the REQ''s own paths'
status: claimed
priority: now
created_at: 2026-09-03T20:05:46Z
user_request: UR-106
domain: backend
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-531, REQ-559, REQ-503]
batch: lifecycle-overhead
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/internal/finalization/
claimed_at: 2026-09-04T18:15:54Z
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-04T18:17:56Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
dispatch_at: 2026-09-04T18:23:51Z
builder_handback_at: 2026-09-04T18:51:24Z
integration_at: 2026-09-04T18:51:24Z
review_at: 2026-09-04T19:05:25Z
remediation_at: 2026-09-04T19:05:25Z
re_review_at: 2026-09-04T19:05:25Z
---

# Hand-Back and Finalize Check Cleanliness Only on the REQ's Own Paths

## What

A path the active REQ does not own, whether untracked, modified, or staged by another session, is never a reason to stop, never surfaced as a blocker, and never committed by the pipeline. Step 7's hand-back settlement and Step 9's finalization check the index and tree only for the REQ's own paths: its run artifacts, its lifecycle files, its write set, and its release paths. Everything else is left exactly as found and named in one progress line. The "preserve in a separate unrelated-work commit" behaviour goes away.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both primes plus the crew rules; traced the refusal to `commitSafety`'s shared-remainder check in `finalization_apply.go` and found the two hand-back step-0 texts plus the one sentence authorizing a pipeline-authored preserve commit.
- [x] **[APPLY]:** Shared-remainder refusal narrowed to discovered recovery groups; step 0's third category rewritten to leave-alone-and-name in both action files; the preserve-commit sentence rewritten. One new Go test pins the narrowing.
- [x] **[UNIFY]:** `git diff --stat` reviewed across 3 modified files and 1 new test; `gofmt -l` empty, `go vet` clean, whole `do-work-cli` module green; no debug artifacts, no stray prints, no commented-out code.

## Why

"Finalization was briefly blocked by an unrelated existing calibration change, which had to be preserved separately." The pipeline authored two commits of other sessions' work to get past its own cleanliness check: 061b7dbf "Preserve corrected calibration history" and 83594c5e "Preserve concurrent maintainability audit draft".

## Context

- `work-reference.md` already states the rule as judgment ("Current-REQ relevance": preserve it, exclude it from this REQ's staging, continue, spend no time on it). Step 7 step 0 contradicts it mechanically: "stop and surface every other `do-work/` path" and "step 0 ends with a clean index". Under four concurrent sessions that condition is rarely true, so the orchestrator invents a way to make it true, which on 2026-09-03 meant committing another session's untracked audit draft under the pipeline's own name.
- The claim transaction already treats "a dirty claim target or index" as shared-target dirt with a typed refusal. This REQ keeps that for the REQ's own target and the index, and drops it for everything else.
- Committing a foreign file is worse than leaving it: it strips the owner's chance to finish it, attributes it to the wrong REQ, and can land a half-written draft on main. Leaving it costs nothing.
- REQ-503 to REQ-510 (the advance command) will absorb the hand-back mechanics later; this REQ changes the rule now so the chain inherits it.

## Detailed Requirements

- Step 7 step 0: the stage category and the allow category stay; the "stop and surface" category becomes "leave alone and name": the path is listed once in the progress output and excluded from staging. Only a dirty path the REQ itself owns still stops.
- The index must be clean of the REQ's own paths before the merge, not of every path; a foreign staged path is unstaged from the index only if it was staged by this run, otherwise it is left and named.
- Finalization: the manifest's exact-allowlist validation is unchanged; a foreign modified or untracked path outside the manifest never refuses the transaction. If a do-work-cli command today refuses on tree dirt outside its declared paths, narrow that check to its declared paths and pin the narrowing with a test in that package.
- Delete the sentences that authorize or describe a pipeline-authored "unrelated work" or "preserve" commit; the pipeline commits only what the REQ declares.
- One line in the run's progress output per foreign path left alone, so the maintainer sees what was skipped without the run stopping.

## Constraints

- Mechanics in programs, judgment in prose (CLAUDE.md). A do-work-cli change carries its own test; no sentence pins.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.
- The claim transaction's refusal on a dirty claim target is unchanged.

## Red-Green Proof
**RED prompt/case:** With an untracked file from another session under `do-work/audits/` and a modified `do-work/calibration-log.tsv` that this REQ did not touch, run a REQ through Step 7 hand-back and Step 9 finalize.
**Why RED now:** step 0 stops on the foreign paths or the orchestrator commits them under the pipeline's name to reach a clean index (061b7dbf, 83594c5e on 2026-09-03).
**GREEN when:** hand-back and finalize complete; the foreign paths are still untracked or modified in the working tree exactly as before; every commit the run made contains only the REQ's declared paths; the progress output names each foreign path once; and a dirty path the REQ itself owns still stops with the typed refusal.
**Validation:** Inferred during capture; the maintainer approved the capture ("do 1, 2 and 3").

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a pipeline step contract.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over budget and `slugged: partial`. Matched because this REQ may change finalization validation in do-work-cli internals.

## Full Context
See `do-work/user-requests/UR-106/input.md` for complete verbatim input.

---

## Triage

**Route: B** - Medium

**Reasoning:** The rule change is stated exactly, but the finalization tree-dirt check in do-work-cli has to be found before it can be narrowed. Outcome clear, location needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ⚠ `bash _dev/tests/maintainer-verify.sh` red BEFORE any change — one pre-existing failure, `_dev/tests/session-start-hook-behavior.sh took 44s; each test file must finish under 30s`. A wall-clock budget miss on a slow container, no assertion failed. Recorded in `do-work/working/baseline-failures.txt` so Step 6.5 separates it from new regressions; not attributable to this REQ and not deferred to a repair REQ.
**Dependencies:** ✓ Go 1.26.1 and ShellCheck 0.11.0 provisioned for this session (container shipped Go 1.24.7 / no ShellCheck)

*Checked by work action*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work-reference.md` (modify) — hand-back step 0's third category, the index rule, and the sentence authorizing a pipeline-authored preserve commit
- `skills/do-work/actions/work.md` (modify) — Step 6's condensed copy of step 0, mirrored to the same rule
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modify) — narrow the shared-remainder refusal to discovered recovery groups
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req560_test.go` (new) — pin the narrowing

**Files I will NOT touch:** `internal/gittransaction/` (its empty-index refusal is out of scope), the manifest allowlist validation, the claim transaction's dirty-target refusal, `CHANGELOG.md` and `VERSION` (Step 9 finalization).

**Acceptance criteria (restated from REQ):**
- [x] A foreign untracked or modified path never stops hand-back or finalization
- [x] Foreign paths are left byte-for-byte as found and named once in the progress output
- [x] Every commit the run makes contains only the REQ's declared paths
- [x] A dirty path the REQ itself owns still stops with the typed refusal
- [x] The manifest's exact-allowlist validation is unchanged
- [x] The narrowing is pinned by a test in the same Go package

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req560_test.go` (new)

**What was done:** `commitSafety`'s shared-remainder refusal (`FINALIZATION-AMBIGUOUS-SHARED-STATE`) now applies only to discovered, tree-inferred recovery groups; a journaled `finalize --manifest` declares its exact write set up front, so a path it never declared cannot belong to the transaction and is left alone. The check that fires first, on paths the journal does declare, is untouched. In prose, hand-back step 0's third category changed from "stop and surface" to "leave alone and name" in both the canonical text and its condensed copy, and the sentence authorizing a pipeline-authored unrelated-work commit was rewritten to say: keep every byte, leave it where it is, name it, carry on.

## Decisions

**D-01 — DECIDE & STATE. A foreign *staged* path is unstaged, not left in the index.** The REQ says a foreign staged path is unstaged only if this run staged it, otherwise left and named. Taken literally that produces the outcome the REQ forbids: every commit in the hand-back sequence is a whole-index commit, so a foreign entry left in the index is not left — it is silently adopted into one of this REQ's commits, breaking the REQ's own GREEN condition and repeating the exact harm the Why is about. Of the three possible outcomes for a staged foreign path — stop, commit it, take it out of the index — the What rules out the first two in one sentence, so the third is the only one left. `git restore --staged` removes the index entry and leaves the file's bytes untouched: the owner loses a staging flag, not a character of work, and the path is named so they can re-stage it.

**D-02 — DECIDE & STATE. The narrowing is keyed on `journal.Discovered`, not on a path list.** The shared-remainder check is correct for `recover-finalization --discover`, where the group is inferred from the tree and an unattributed shared path really could be a torn tail of the same interrupted transaction. It is meaningless for `finalize --manifest`, where the journal binds the exact write set up front. Keying on the condition (inferred or declared?) rather than on tolerated path spellings follows CLAUDE.md's "state conditions, not lists".

**D-03 — DECIDE & STATE. The manifest allowlist, both empty-index checks, and discovery's own refusal are untouched.** `gittransaction.CommitExactPaths` commits the whole index and refuses a non-empty index itself, in a package outside this REQ's write set; narrowing finalization's check alone would only move the same refusal later and report it worse. The REQ's own Context also keeps the shared-target rule for the REQ's own target and the index, and scopes the never-refuse rule to foreign modified or untracked paths.

**D-04 — DECIDE & STATE. The preserve-commit sentence was rewritten, not deleted.** The only sentence authorizing a pipeline-authored preserve commit sat in "Stuck Runs Hand Off to Judgment". Deleting the bullet outright would leave that dirt class with no guidance, so it now says: keep every byte, leave it where it is, name it in the progress output, carry on — and if a canonical command still refuses because of that path, fall through to the shared-state bullet below it, which already names the resolving verb.

**D-05 — DECIDE & STATE. No changelog or version write from the builder.** Both are outside the declared write set and belong to finalization.

## Discovered Tasks

- **impact-moderate, report only** — the two hand-back step-0 texts (`work-reference.md` and `work.md` Step 6's condensed copy) are a hand-maintained duplicate pair with no mechanical check pinning them together. Any future change to the sequence has to be made in both or they drift; `prime-action-files.md` names exactly this shape as `alternate-writer-contract-drift`.
- **impact-moderate, report only** — `_dev/tests/session-start-hook-behavior.sh` fails inside a builder worktree because its launcher subshell resolves the system Go instead of the pinned toolchain, making the hook probes useless as a signal for any builder.
- **impact-noncritical, report only** — `internal/publication`'s hostile-argv probe hard-fails when `just` is absent rather than skipping, turning a missing optional binary into a red module.

## Qualification

Passed — 4 files verified in the merge range `6adb8b9..cb3a831`, 6 acceptance criteria traced, P-A-U confirmed. `qualify.sh` returned `OK: mechanical qualification passed`. Its one WARN, that the new test file has no static reference, is the documented test-file exception. Judgment checks: the `finalization_apply.go` change is real control flow, not a stub; the declared-path refusal above it is untouched, which is the safety property the REQ promises to keep.

## Testing

**Tests run:** `go test ./internal/finalization/ -count=1` and `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Finalization package green (22.1s). Repository gate **exit 0, no failures**.

**Red-green validation:**
- `TestFinalizeIgnoresForeignTreeDirtOutsideTheManifest`: ✗ before implementation → ✓ after. The RED refusal, reproduced independently by the orchestrator against the pre-merge revision, is `FINALIZATION-AMBIGUOUS-SHARED-STATE` with `AffectedPaths: [do-work/audits/maintainability-draft.md, do-work/calibration-log.tsv]` and outcome `rolled_back`.

**Field evidence from this same run.** Earlier, finalizing REQ-559 was refused by this exact code path:

```
FINALIZATION-AMBIGUOUS-SHARED-STATE
paths: do-work/runs/work-2026-09-04-182017/REQ-515-handback.md,
       do-work/runs/work-2026-09-04-182017/REQ-559-review.md,
       do-work/runs/work-2026-09-04-182017/REQ-560-handback.md
```

Three files belonging to the other two REQs in this wave, none of them owned or declared by REQ-559, and none of them a reason that REQ could not complete. The orchestrator moved them aside, finalized, and moved them back — the manual version of what the pipeline used to reach by committing foreign work under its own name. The REQ was captured from a 2026-09-03 incident and the defect reproduced unprompted during its own fix's verification.

**New tests added:**
- `TestFinalizeIgnoresForeignTreeDirtOutsideTheManifest` — with a modified tracked `do-work/calibration-log.tsv` and an untracked `do-work/audits/maintainability-draft.md` that the REQ neither owns nor declares: `finalize --manifest` succeeds, both foreign paths are still dirty afterwards exactly as found, and the primary commit's file list contains neither.

**Not weakened:** the whole finalization package passes unchanged, including `TestFinalizeAcceptsWorkingRequestDirtWrittenByThePipeline` and the recovery/discovery suite, so discovery keeps its original refusal.

## Review

**Verdict: Pass** — independent review, overall 92%, acceptance Pass. Full record: `do-work/runs/work-2026-09-04-182017/REQ-560-review.md`. Eight findings, none blocking.

The reviewer fired both refusals this REQ deliberately keeps rather than assuming they survived: discovery still refuses on shared dirt outside its inferred group, and the manifest's exact-allowlist validation still refuses a `commit_paths` missing a planned target. No sentence authorizing a pipeline-authored preserve commit survives in the shipped tree.

**D-01 upheld, premise checked rather than accepted.** The reviewer read `internal/gittransaction/exact_commit.go` and confirmed `CommitExactPaths` refuses a non-empty index, then runs `git add -A -- <dirty declared paths>` and a plain `git commit` that takes the whole index. Step 0's own prose ends the same way. So a foreign staged entry left in the index really is committed, the REQ's `## What` forbids both stopping and committing, and unstaging is the only outcome left. Faithful to intent, not an overreach.

**Findings fixed rather than deferred (F1, F2, F3, F7).** Four of the eight sat inside this REQ's own write set and three of them weakened the REQ's own promise, so they were fixed and re-merged rather than filed:

- **F1 (Important).** The prose this REQ added claimed `git restore --staged` "leaves the file's bytes exactly as they are so nothing the other writer typed is lost". That is false for a path staged and then modified again: the staged snapshot survives only as a dangling blob reachable through `git fsck --lost-found`. In a rule whose entire purpose is protecting another writer's work, an agent reading that sentence would believe the operation is always lossless and would not warn the owner. Both copies now state the working-tree guarantee accurately and name the `MM` case as the exception the progress line must call out.
- **F2 (Important).** `grep -rn "AMBIGUOUS-SHARED-STATE" --include=*.go` returned exactly one hit, the production string — nothing pinned either refusal this REQ kept, before or after the change. A new test, `TestCommitSafetyStillRefusesSharedDirtForADiscoveredGroup`, drives `commitSafety` directly over one tree in both modes; it fails on the pre-change code, so it catches a silent widening, which is exactly what this REQ's Constraints forbid.
- **F3.** The condensed copy said a foreign path is "another session's"; the canonical text says "another session's or the user's". The user's own uncommitted work is precisely the class the originating incident was about, and the two texts are a hand-maintained pair with no mechanical check.
- **F7.** The new test accepted `"M  "` as well as `" M "`. `"M  "` means staged, the opposite of leaving the path alone.

**Findings recorded, not fixed:** F4 (the prose promises per-entry checkpoint granularity that finalize's whole-file digest does not have), F5 (`finalize --manifest` now passes foreign dirt silently where it used to name it in the refusal; step 0 already named the same paths earlier in the run), F6 (this REQ's own Implementation Summary described the surviving check too broadly — corrected in place above), F8 (the touched line adopted the pre-existing `sharedFinalizationPath` helper, a behaviour-free de-duplication on a line already being edited).

**Remediation:** commit `0e1b687` on the REQ's branch, re-merged once. `<pre>` held at `6adb8b9`, so the final merge range is `6adb8b9..3430117` and covers the original work plus the fixes.

## Lessons Learned

**What worked:**
- Checking a builder's premise instead of its conclusion. D-01 deviated from the REQ's literal sentence; the reviewer settled it by reading `CommitExactPaths` and confirming the commit really is whole-index. The deviation was right, and it is right for a reason now on the record rather than an assertion.
- Keying the narrowing on `journal.Discovered` rather than on a list of tolerated path spellings. One condition answers the real question — was this group inferred from the tree or declared in a manifest — and it cannot go stale the way a path list does.

**What didn't:**
- The first version shipped a false safety claim in the same sentence that introduced the safety behaviour. "Unstaging leaves the bytes untouched" is true of the working tree and false of a staged-then-modified snapshot. A sentence that reassures is worth more scrutiny than one that warns, because nobody re-checks reassurance.
- Pinning only the behaviour that changed. The two refusals this REQ deliberately *kept* had no test at all, so the change protected its own new path while leaving the paths it promised not to touch open to silent deletion. Kept behaviour needs a lock-in exactly when a change walks past it.

**Worth knowing:**
- `finalize --manifest` now passes foreign dirt with no finding and no path list, where it used to name those paths in its refusal. Step 0 names the same paths earlier in the run, so the information exists, but the CLI alone no longer reports what it skipped.
- The prose says only a dirty path the REQ owns stops, listing its checkpoint entry. Finalize's actual guard is a whole-file `expected_checkpoint_sha256`, so any foreign edit to `do-work/CHECKPOINT.md` still refuses. The prose promises per-entry granularity the code does not have.
- The two hand-back step-0 texts remain a hand-maintained duplicate pair with no mechanical check. This REQ had to edit both, and the review found them already drifting by one clause.

## Orientation

Another session's or the user's uncommitted files no longer stop a run or get swept into its commits. Lives in the work pipeline's hand-back and finalization boundary: the rule in the work action and its reference at hand-back step 0, the mechanism in the do-work CLI's finalization commit-safety check, which now applies its shared-remainder refusal only to a recovery group inferred from the tree.

Neither prime this REQ lists was made stale: `_dev/primes/prime-action-files.md` and `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` both still resolve every path they reference.

## Post-Merge Verification

`bash _dev/tests/maintainer-verify.sh` against the merged tree at `3430117`: **exit 0, no failures.** Run twice for this REQ — once after the first merge and once after the remediation merge — green both times.
