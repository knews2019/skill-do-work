---
id: REQ-722
title: 'Cleanup Pass 6 skips assets folders when scanning for blanked records'
status: completed
created_at: 2026-11-14T21:47:51Z
user_request: UR-173
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: [skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go, skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go]
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-14T22:19:27Z
  basis:
    - trivial short-circuit
claimed_at: 2026-11-14T22:17:58Z
route: A
dispatch_at: 2026-11-14T22:19:40Z
builder_handback_at: 2026-11-14T22:20:25Z
integration_at: 2026-11-14T22:21:00Z
review_at: 2026-11-14T22:26:30Z
kb_status: pending
commit: 0c3c94abe6aead467f61665467a2e28f323565f1
heavy_verified_at: 2026-11-14T22:29:53Z
heavy_verified_revision: 0c3c94abe6aead467f61665467a2e28f323565f1
completed_at: 2026-11-14T22:30:17Z
release_at: 2026-11-14T22:30:17Z
---
# Cleanup Pass 6 Skips Assets Folders When Scanning for Blanked Records
## What
Cleanup Pass 6 (`AddBlankedRepairs`, `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go:117`) treats every `REQ-*.md` or `UR-*.md` file under `do-work/archive`, `do-work/queue` and `do-work/working` as a request record, including evidence files inside `assets/` folders. Such a file has no frontmatter, so the pass reports `BLANKED-RECORD-UNRECOVERABLE` and `cleanup` exits 1. Skip `assets/` folders in the Pass 6 walk and add one test for it.
## Why
A consumer repo got 49 false warnings and exit 1 after archiving two user requests. The cleanup action stops on a nonzero exit, so the false warnings block a clean cleanup.
## Context
- Finding 1 of a validate-feedback triage, source: pasted text (consumer feedback). Original severity: none given. Verdict: Accept. Surface-cost: N/A (direct fix, no new guard).
- Evidence: the walk at `cleanup_git.go:105-119` filters on the file name only. `isBlankedDocument` (`cleanup_git.go:159`) returns true when `requestmodel.ParseDocument` fails, which it does for any file without frontmatter. `repositorymodel` already skips these folders (`internal/repositorymodel/repository_model.go:192`, `directoryEntry.Name() == "assets"` returns `fs.SkipDir`).
- No queued REQ in any UR shares this root cause (the queue was empty at capture).
- Check whether `tools/checks/blanked-req-scan.sh` (named by `actions/cleanup.md` Pass 6) has the same gap. If it does, the builder decides whether it belongs in this fix.
## Red-Green Proof
**RED prompt/case:** A repo whose only do-work file is `do-work/archive/UR-001/assets/REQ-001-evidence.md` (plain text, no frontmatter), committed. Run `do-work-cli --format json cleanup --dry-run`.
**Why RED now:** The output holds one `BLANKED-RECORD-UNRECOVERABLE` finding for that path and the exit code is 1. Reproduced during triage on 2026-11-14.
**GREEN when:** A new Go test in `internal/cleanup/cleanup_git_test.go` calls `AddBlankedRepairs` on that fixture and gets no finding and no group, and the existing test `TestBlankedRestorationRequiresExactConsentAndUsesGitHistory` still passes (a truly blanked record outside `assets/` is still reported).
**Validation:** Inferred during capture
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24277 tokens, over the budget, and its row is `slugged: partial` so no targeted entry is legal. Matched because the satellite's prime governs `internal/cleanup` (families `closed-enumeration-for-a-condition`, `guard-on-one-entry-path`).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the prime and the walk. Add an exact-name assets directory skip (filepath.SkipDir) to the AddBlankedRepairs walk, matching repositorymodel; pin it with one Go test on the RED fixture, test first.
- [x] **[APPLY]:** Test written first and seen RED (BLANKED-RECORD-UNRECOVERABLE), then the 4-line skip in cleanup_git.go; only the two planned files changed.
- [x] **[UNIFY]:** git diff --stat ef4f4dc6..0355ac11: 2 files, 18 insertions. Checked cleanup_git.go (skip placed after the symlink check, before the name filter) and cleanup_git_test.go (fixture matches the REQ RED case). gofmt -l clean, go vet ./internal/cleanup clean, go test ./internal/cleanup ok. No debug artifacts.
*Source: do-work capture-request: cleanup Pass 6 (AddBlankedRepairs, cleanup_git.go:117) flags assets/REQ-\*.md evidence files as BLANKED-RECORD-UNRECOVERABLE and exits 1.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the walk, the two files and the existing skip rule to copy (`repositorymodel` already skips `assets` folders). The fix is one directory check plus one Go test, and the Red-Green Proof gives the exact fixture. The `blanked-req-scan.sh` question was answered during triage by reading the code: that launcher runs `doctor`, which reads records through `repositorymodel` and already skips `assets` folders, so it has no gap (recorded as D-01 in Decisions).

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- D-01 (DECIDE & STATE): `tools/checks/blanked-req-scan.sh` stays out of this fix because it has no gap. The launcher runs `do-work-cli blanked-req-scan`, which is the `doctor` handler (`internal/doctor/doctor_commands.go:22`). `doctor` finds damaged records through `repositorymodel.DiscoverRepository`, and that walk already skips every `assets` folder (`internal/repositorymodel/repository_model.go:192`).
- D-02 (DECIDE & STATE): The new check matches the folder name `assets` exactly, the same rule `repositorymodel` uses, and it applies in all three scan roots (archive, queue, working). Records never live inside an `assets` folder, so no real record is hidden.

## Discovered Tasks

- The standalone `archive-collision` helper (`skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:164-176`, launched by `tools/checks/archive-collision.sh`) walks `do-work/archive` by file name with no `assets` skip, so an evidence file such as `UR-001/assets/REQ-001-evidence.md` would be listed as a collision for REQ-001. The `advance` hold uses the `repositorymodel` snapshot and is not affected, and no action prescribes the helper. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (modified)

**What was done:** The cleanup Pass 6 walk in `AddBlankedRepairs` now returns `filepath.SkipDir` for every directory named `assets`, so evidence files such as `assets/REQ-001-evidence.md` are no longer read as blanked records. The new test `TestBlankedScanSkipsEvidenceFilesInAssetsFolders` commits that exact fixture plus a blanked record in the same UR folder, and expects one finding for the record only; it failed with `BLANKED-RECORD-UNRECOVERABLE` before the fix. After review F2 the blanked sibling was added, so a skip widened to the whole UR folder would also fail the test.

*Completed by work action (Route A)*

## Qualification

**Range:** `ef4f4dc6..0c3c94ab` (two merges of the REQ branch: build 627a3be3 merged at 0355ac11, review fix F2 0a8742ad merged at 0c3c94ab)
**Mechanical gates:** qualify ✓ satisfied on the first call, at `ef4f4dc6..0355ac11`. The second merge changes only the test file.

**Requirement trace (orchestrator, read against the diff):**
- Skip `assets/` folders in the Pass 6 walk: `cleanup_git.go` returns `filepath.SkipDir` for a directory named `assets`, inside the one walk callback that serves `do-work/archive`, `do-work/queue` and `do-work/working`.
- One test for it: `TestBlankedScanSkipsEvidenceFilesInAssetsFolders` builds the REQ's RED fixture (`do-work/archive/UR-001/assets/REQ-001-evidence.md`, plain text, committed) and also commits a blanked record beside it (`do-work/archive/UR-001/REQ-002-done.md`), and expects exactly one `BLANKED-RECORD-REQUIRES-CONSENT` finding, for that record only, and no group.
- A truly blanked record outside `assets/` is still reported: `TestBlankedRestorationRequiresExactConsentAndUsesGitHistory` passes unchanged.
- `blanked-req-scan.sh` check: done, no gap (D-01).

**Live flow:** the REQ's RED case run end to end on a scratch Git repository holding only that evidence file. `do-work-cli --format json cleanup --dry-run` exited 1 with one `BLANKED-RECORD-UNRECOVERABLE` from the main tree before the merge, and exited 0 with no blanked finding from the build worktree.

*Qualified by work action*

## Testing

**Tests run:** focused probe `go test ./internal/cleanup -run 'TestBlankedScanSkipsEvidenceFilesInAssetsFolders|BlankedRestoration' -count=1` (run by `advance` through the test gate), then `bash _dev/tests/maintainer-verify.sh` at the merge revision 0c3c94ab from a detached checkout
**Result:** ✓ Focused probe green through the test gate; repository gate exit 0, gate wall 148s (do-work-cli module 69s, 981 tests; slowest file `internal/finalization/finalization_recovery_test.go` 20.7s, under the 30s budget). `gofmt -l` and `go vet ./internal/cleanup` clean; `go test ./internal/cleanup` ok.

**Repository gate runs:** one completed run, green. A first run at 0355ac11 was stopped before it finished, because review finding F2 changed the test; the gate then ran once at the second merge 0c3c94ab.

**Red-green validation (`tdd: true`):**
- RED at 627a3be3's parent (test written first, fix not yet applied): `--- FAIL: TestBlankedScanSkipsEvidenceFilesInAssetsFolders` with `cleanup_git_test.go:151: assets evidence file was scanned as a blanked record`, finding code `BLANKED-RECORD-UNRECOVERABLE` for `do-work/archive/UR-001/assets/REQ-001-evidence.md`.
- GREEN at 627a3be3: `--- PASS: TestBlankedScanSkipsEvidenceFilesInAssetsFolders` and `--- PASS: TestBlankedRestorationRequiresExactConsentAndUsesGitHistory`.
- After F2 (0a8742ad), the widened test was run against the parent's `cleanup_git.go`: `--- FAIL ... want one consent finding for do-work/archive/UR-001/REQ-002-done.md and none for the assets evidence file`; with the fix it passes.
- End to end, the REQ's RED case: a scratch Git repository holding only `do-work/archive/UR-001/assets/REQ-001-evidence.md`. `do-work-cli --format json cleanup --dry-run` exited 1 with one `BLANKED-RECORD-UNRECOVERABLE` before the fix and exited 0 with no blanked finding after it.

**Untested:** the consumer repository that reported the 49 warnings was not available; the fix is proven on the reduced fixture only.

**Heavy verification plan:**
- Range: ef4f4dc612a9617dc2ec5ee3db71e134ffe1ba88..0c3c94abe6aead467f61665467a2e28f323565f1
- do-work-cli-integrations, staged-skills, updater, installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane>` — both changed paths matched subtree `skills/do-work/tools/do-work-cli` (and `skills` for staged-skills)

*Verified by work action*

## Review

**Overall: 93%** | 2026-11-14T22:26:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Reviewer:** one independent foreground reviewer read `ef4f4dc6..0355ac11` against the code, re-ran RED with the parent's `cleanup_git.go` in a scratch copy (fails with exactly `BLANKED-RECORD-UNRECOVERABLE`), and verified D-01 (`blanked-req-scan` maps to `handleDoctor`, which reads through `repositorymodel.DiscoverRepository`).

**Minor findings:**
- F1: the standalone `archive-collision` helper (`skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:164-176`) walks `do-work/archive` by file name with no `assets` skip, so `archive-collision REQ-001` would list `UR-001/assets/REQ-001-evidence.md` as a collision. The `advance` archive-collision hold is not affected: it reads the `repositorymodel` snapshot (`internal/lifecycleadvance/queue_commands.go:92`). No action prescribes the helper. Outside this REQ's write set; recorded in Discovered Tasks. — impact-negligible → report only
- F2 (Nit): the first test could not tell a narrow skip from one widened to the whole UR folder. Fixed in 0a8742ad (merged at 0c3c94ab): the fixture adds a blanked `do-work/archive/UR-001/REQ-002-done.md` and the test expects exactly one `BLANKED-RECORD-REQUIRES-CONSENT` finding, for that path only (`skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go:144`). — impact-negligible → report only

**Acceptance:** Pass — the REQ's RED case is green in the unit test and end to end through `cleanup --dry-run` (exit 1 before, exit 0 after); a real blanked record, inside a UR folder and flat in the archive, is still reported.
**Restatement sweep:** no prose names the Pass 6 scan scope; `actions/cleanup.md` Pass 6 describes blanked records, which `assets` evidence files are not, so nothing else changes.
**Suggested testing:** run `do-work cleanup --dry-run` in the consumer repo that reported the 49 warnings and expect no `BLANKED-RECORD-*` finding for its `assets/` files.
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: ef4f4dc612a9617dc2ec5ee3db71e134ffe1ba88
- Target revision: 0c3c94abe6aead467f61665467a2e28f323565f1
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — both changed paths matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both changed paths matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — both changed paths matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — both changed paths matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: 0c3c94abe6aead467f61665467a2e28f323565f1
- Execution revision: 0c3c94abe6aead467f61665467a2e28f323565f1 (detached checkout)
- do-work-cli-integrations: exit 0, 72s, executed
- staged-skills: exit 0, 45s, executed
- updater: exit 0, 69s, executed
- installer: exit 0, 26s, executed
- No `HEAVY-RUN-LANE-SKIPPED` finding.

## Lessons Learned

**What worked:** Writing the test from the REQ's RED case first. It failed with the exact finding the consumer reported, so the fix was proven against the real failure and not a guess. Reading `blanked-req-scan` down to its handler answered the REQ's side question in two minutes: it runs `doctor`, which reads through `repositorymodel` and already skips `assets` folders. The independent reviewer, launched right after qualify, found that the first test could not tell a narrow skip from one that hides the whole UR folder; adding a blanked sibling record fixed that before the gate ran to the end.
**What didn't:** The first gate run had to be stopped for the review fix, so the gate ran to the end once, at the second merge.
**Worth knowing:** do-work-cli has more than one walker over `do-work/archive`, `queue` and `working`. `repositorymodel.DiscoverRepository` owns the folder skips (`assets`, `deliverables`, dot folders); a walker of its own that picks records by file name must copy those skips or use the snapshot. The standalone `archive-collision` helper is one more walker without the skip (Discovered Tasks).

## Orientation

Now `do-work cleanup` leaves evidence files in a user request's `assets/` folder alone. Pass 6, which looks for request records that were emptied by a bad write, skips every `assets` folder, as the rest of do-work-cli already did, so a healthy archive no longer gives `BLANKED-RECORD-UNRECOVERABLE` warnings and exit 1. This lives in the cleanup package of do-work-cli (`internal/cleanup/cleanup_git.go`); no action prose changed. Prime staleness check: `prime-do-work-cli.md` names no path this change moved or removed.
