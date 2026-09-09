---
id: REQ-622
title: 'Correct verified prose drift'
status: completed
route: B
heavy_verified_at: 2026-09-08T22:22:12Z
heavy_verified_revision: 628c0e34ef499e34e8b123297d1d8965bad04c7e
review_at: 2026-09-08T22:18:53Z
kb_status: pending
commit: 628c0e34ef499e34e8b123297d1d8965bad04c7e
integration_at: 2026-09-08T22:13:03Z
builder_handback_at: 2026-09-08T22:12:22Z
dispatch_at: 2026-09-08T22:01:19Z
estimate:
  p50_active_minutes: 60
  confidence: medium
  calculated_at: 2026-09-08T21:54:29Z
  basis:
    - Route B
    - 21-file write set
    - 3 subsystems involved
    - 19 acceptance criteria
    - cross-route regression gates
created_at: 2026-09-08T21:49:05Z
user_request: UR-130
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-kanban-board.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: []
related: ["REQ-623", "REQ-624"]
batch: instruction-history-cleanup
write_set: ["skills/do-work-board/tools/queue-kanban/open_work.go", "skills/do-work-board/tools/queue-kanban/testing.go", "skills/do-work-board/actions/board.md", "_dev/primes/prime-action-files.md", "skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go", "skills/do-work-board/tools/queue-kanban/web/board-core.js", "skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md", "skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/citations.go", "skills/do-work/actions/abandon.md", "skills/do-work-board/tools/queue-kanban/generate.go", "README.md", "skills/do-work/actions/clarify.md", "skills/do-work/tools/do-work-cli/internal/publication/answer.go", "skills/do-work/actions/verify-requests.md", "_dev/lessons/validated-runtime-boundaries.md", "do-work/prose-backlog.md", "do-work/lessons-index.md"]
claimed_at: 2026-09-08T21:53:21Z
completed_at: 2026-09-08T22:22:36Z
release_at: 2026-09-08T22:22:36Z
---

# Correct Verified Prose Drift

## What
Reconcile the 19 open entries in do-work/prose-backlog.md against current code, canonical contracts, and recorded intent. Correct surviving prose discrepancies and record which entries are resolved, obsolete, or still need an intent decision. Prefer removing duplicate statements; preserve runtime behavior.

## Detailed Requirements
- Revalidate the 19 open entries preserved in `do-work/user-requests/UR-130/assets/prose-backlog-at-capture.md`, including the stale isolation lesson. The snapshot identifies this drain's items; it does not assert they remain valid today.
- Give every entry an evidence-backed disposition: corrected, already obsolete, or unresolved intent. Tick resolved/obsolete lines in the live backlog with evidence and `(drained by REQ-622)`; retain unresolved lines and explain the conflict.

## Constraints
Preserve runtime behavior. Current code, canonical prose, and recorded intent must be reconciled; a disagreement is not permission to silently change a rule. Remove redundant wording where that resolves the discrepancy.

## Dependencies
First in the approved batch. REQ-623 (Reduce orchestration instruction duplication) follows this request; REQ-624 (Addendum: Trim the shipped changelog to 50 entries) follows that request.

## Builder Guidance
The backlog is the bounded scope, not a mandate to rewrite surrounding mechanisms. No pending or pending-answers candidate in any UR shares this request's root cause; the seven existing queued requests address different release, commit-integrity, and diagnostics defects.

## Completion Proof
Account for all 19 source items and verify changed references and affected contracts against their canonical owners. Record the evidence behind each disposition; code-comment edits preserve executable behavior. An unresolved intent decision remains explicit rather than being reported as a fixed discrepancy.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens; action/prime prose and condition-preserving extraction. Partial slug coverage prevents targeted selection.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens; timestamp script and prescribed containment wording. Partial slug coverage prevents targeted selection.
- `_dev/primes/lessons-kanban-board.md` — 5912 tokens; the board prime governs named board paths. Partial slug coverage prevents targeted selection.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 6888 tokens; the board's own prime governs the named comments, docs, and testing files. Partial slug coverage prevents targeted selection.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 15905 tokens; alternate answer-writer contract drift and publication comments. Partial slug coverage prevents targeted selection.

## Full Context
See `do-work/user-requests/UR-130/input.md` for the user instruction, adopted report, and batch decisions. The approved prompt excerpts are in `do-work/user-requests/UR-130/assets/rev-v2-approved-scope.md`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Revalidate the 19 items; correct the 16 surviving prose discrepancies, tick three obsolete findings with evidence, then verify unchanged executable behavior and references.
- [x] **[APPLY]:** Integrated 16 source corrections and all 19 backlog dispositions, plus the lesson-index size correction.
- [x] **[UNIFY]:** Read the full hand-back and merged diff; builder reference and focused CLI/board checks passed. Merged-tree checks follow.

## Triage

**Route: B** — Medium.

**Reasoning:** The 19 recorded items name bounded locations, but current code and historical intent need revalidation before selecting safe prose corrections.

**Planning:** Not required; exploration-guided implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

Revalidated all 19 recorded items: 16 surviving corrections and 3 obsolete entries. The timestamp pointers and obsolete forensics table are gone; the exact-once containment assertion was removed. REQ-460 and REQ-543 resolve the apparent containment/isolation questions without changing behavior. The Copy lesson now lives in its satellite. Durable evidence for every source item is retained in `do-work/prose-backlog.md`; independent review confirmed all 19 dispositions.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/open_work.go` (modify) — verified prose/comment correction
- `skills/do-work-board/tools/queue-kanban/testing.go` (modify) — verified prose/comment correction
- `skills/do-work-board/actions/board.md` (modify) — verified prose/comment correction
- `_dev/primes/prime-action-files.md` (modify) — verified prose/comment correction
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — verified prose/comment correction
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modify) — verified prose/comment correction
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modify) — verified prose/comment correction
- `skills/do-work-board/docs/board-guide.md` (modify) — verified prose/comment correction
- `skills/do-work-board/tools/queue-kanban/citations.go` (modify) — verified prose/comment correction
- `skills/do-work/actions/abandon.md` (modify) — verified prose/comment correction
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — verified prose/comment correction
- `README.md` (modify) — verified prose/comment correction
- `skills/do-work/actions/clarify.md` (modify) — verified prose/comment correction
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — verified prose/comment correction
- `skills/do-work/actions/verify-requests.md` (modify) — verified prose/comment correction
- `_dev/lessons/validated-runtime-boundaries.md` (modify) — verified prose/comment correction
- `do-work/prose-backlog.md` (modify) — backlog disposition / lesson-index integration seam
- `do-work/lessons-index.md` (modify) — backlog disposition / lesson-index integration seam

**Acceptance criteria (restated from REQ):**
- [x] Account for all 19 snapshot entries with current evidence.
- [x] Correct surviving discrepancies while preserving runtime behavior, canonical contracts, and recorded intent.
- [x] Tick resolved/obsolete records with evidence; leave any unresolved conflict explicit.
- [x] Verify affected references and focused contracts.

## Pre-Flight

Canonical `advance` returned satisfied request-bound preflight and green-gate records. Baseline: `bash _dev/tests/shipped-package-reference-contract.sh` passed. Canonical gate: `bash _dev/tests/maintainer-verify.sh` passed (128s; 402 board tests, 815 CLI tests; all test-file budgets below 30s). Only owner-authored REQ/run evidence is dirty; implementation sources are unchanged.

## Implementation Summary

**Files changed:**
- `README.md` (modified)
- `_dev/lessons/validated-runtime-boundaries.md` (modified)
- `_dev/primes/prime-action-files.md` (modified)
- `skills/do-work-board/actions/board.md` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/citations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modified)
- `skills/do-work-board/tools/queue-kanban/open_work.go` (modified)
- `skills/do-work-board/tools/queue-kanban/testing.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modified)
- `skills/do-work/actions/abandon.md` (modified)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/verify-requests.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modified)

**What was done:** Corrected 16 surviving prose discrepancies without changing executable code or assertions. Accounted for all 19 captured backlog entries (three already obsolete); the ordinary cancellation example changed, while the already-archived example was confirmed correct. Source text decreased by 1,636 bytes. Backlog and lesson-index updates were integrated in the same merge.

**Integration range:** `8b6f4f153b11012b3758b3f668c95f6301f39e71..628c0e34ef499e34e8b123297d1d8965bad04c7e`.

## Qualification

Passed — canonical advance returned satisfied qualify and scope-drift records for the retained merge range, with no findings. The 18-path merge matches the declared scope; all 19 backlog dispositions map to the captured snapshot. Independent Go AST comparison passed for all six changed Go files; their executable syntax is identical. The JavaScript diff changes only the helper doc comment, and no semantic Go directives changed.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on merged main via canonical timing wrapper; focused shipped-package-reference check via advance; builder focused publication and board tests.
**Result:** All passed. Merged canonical gate: 106 seconds, 402 board tests and 815 CLI tests. All per-file test budgets below 30 seconds. Advance returned satisfied focused-test, scope-drift and green-gate records at `628c0e34ef499e34e8b123297d1d8965bad04c7e`. No new tests or executable changes.

**Heavy verification plan:**
- Range: `8b6f4f153b11012b3758b3f668c95f6301f39e71..628c0e34ef499e34e8b123297d1d8965bad04c7e`
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — skills/do-work-board/tools/queue-kanban/citations.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/generate.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/open_work.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/testing.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/web/board-core.js matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — skills/do-work-board/tools/queue-kanban/citations.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/generate.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/open_work.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/testing.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/web/board-core.js matched subtree skills/do-work-board/tools/queue-kanban
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-board/actions/board.md matched subtree skills; skills/do-work-board/docs/board-guide.md matched subtree skills; skills/do-work-board/tools/queue-kanban/citations.go matched subtree skills; skills/do-work-board/tools/queue-kanban/generate.go matched subtree skills; skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md matched subtree skills; skills/do-work-board/tools/queue-kanban/open_work.go matched subtree skills; skills/do-work-board/tools/queue-kanban/testing.go matched subtree skills; skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go matched subtree skills; skills/do-work-board/tools/queue-kanban/web/board-core.js matched subtree skills; skills/do-work/actions/abandon.md matched subtree skills; skills/do-work/actions/clarify.md matched subtree skills; skills/do-work/actions/verify-requests.md matched subtree skills; skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — README.md matched exact path README.md; skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills/do-work/tools/do-work-cli

## Review

**Overall: 100%** | 2026-09-08T22:16:18Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):** None.

**Minor findings:** None.

**Nit findings:**
- **F1** — `skills/do-work-board/tools/queue-kanban/citations.go:94` still says “three values” above four named surfaces. This pre-existing cosmetic count was disclosed by the builder and is outside the captured link-suppression discrepancy; no score deduction. — impact-negligible → report only

**Acceptance:** Pass — all 19 dispositions independently checked; executable-source edits are comments; focused fixtures, merged canonical gate, and canonical test-gate passed. Heavy verification remains the orchestrator's separate completion gate.
**Suggested testing:** 0 items beyond already selected heavy lanes.
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Revalidated each captured finding against current code and recorded decisions; all 19 now carry evidence in `do-work/prose-backlog.md`.
**What didn't:**
- The original cancellation finding overgeneralized across active and already-archived targets. Reading the early archive branch prevented an incorrect prose change.
**Worth knowing:**
- Historical Copy guidance describes source bytes; later clipboard annotation is a separate stage. No new reusable lesson rule is needed beyond the corrected existing satellites. KB handoff remains pending; no KB write was requested.

## Orientation

Core workflow instructions and board guidance now match current behavior. The touched action and board primes still resolve their referenced owners; no system map changed.

## Heavy Verification Plan

Base revision: `8b6f4f153b11012b3758b3f668c95f6301f39e71`
Target revision: `628c0e34ef499e34e8b123297d1d8965bad04c7e`

- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — skills/do-work-board/tools/queue-kanban/citations.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/generate.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/open_work.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/testing.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/web/board-core.js matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — skills/do-work-board/tools/queue-kanban/citations.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/generate.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/open_work.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/testing.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/web/board-core.js matched subtree skills/do-work-board/tools/queue-kanban
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-board/actions/board.md matched subtree skills; skills/do-work-board/docs/board-guide.md matched subtree skills; skills/do-work-board/tools/queue-kanban/citations.go matched subtree skills; skills/do-work-board/tools/queue-kanban/generate.go matched subtree skills; skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md matched subtree skills; skills/do-work-board/tools/queue-kanban/open_work.go matched subtree skills; skills/do-work-board/tools/queue-kanban/testing.go matched subtree skills; skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go matched subtree skills; skills/do-work-board/tools/queue-kanban/web/board-core.js matched subtree skills; skills/do-work/actions/abandon.md matched subtree skills; skills/do-work/actions/clarify.md matched subtree skills; skills/do-work/actions/verify-requests.md matched subtree skills; skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — README.md matched exact path README.md; skills/do-work/tools/do-work-cli/internal/publication/answer.go matched subtree skills/do-work/tools/do-work-cli

## Disposition Evidence

Numbers follow the capture snapshot's open-entry order.

| Item | Disposition | Verified owner / consequence |
|---|---|---|
| 1 | Already obsolete | `repair-req-timestamps.sh` is a CLI launcher; `forensics.md` delegates future/out-of-order timestamps to doctor. The old work-reference Check 11 pointer is absent. |
| 2 | Corrected | `testing.go` uses its mutex around the two Testing operations. `open_work.go` now states only its own read-only contract, leaving the global count to the prime. |
| 3 | Corrected | `board.md` checklist points to its Input table, which includes verify, instead of repeating a stale mode list. |
| 4 | Corrected | `shipped-package-reference-contract.sh` contains both same-package and cross-package resolution and fixtures. The prime describes both. |
| 5 | Corrected | The timeline test comment names rendered state and shared filters without a property count. Test assertions are unchanged. |
| 6 | Already obsolete | `forensics.md` Check 14 maps findings, skips, and not-applicable output classes; the old finite probe table is gone. |
| 7 | Corrected | `formatElapsedDuration` callers include live timers, completed wall time in `board-cards.js`, and phase spans in `board-detail.js`. The revised interval description covers all. |
| 8 | Corrected at moved home | The REQ-089 lesson in `lessons-do-kanban.md` preserves the source-byte rule and qualifies later clipboard annotation. The maintainer board prime already states this distinction. |
| 9 | Corrected | `board-filters.js` retains partial ID/title matching and resolves whole cited IDs; `board-cards.js` emits the “cites” badge for citation-only matches. |
| 10 | Corrected | `citations.go` protects clipboard link labels; `board-detail.js` decorates authored anchors. The revised comments retain ancestor-link precedence without claiming both surfaces suppress anchors. |
| 11 | Corrected in part | `requeststate/state_plan.go:64–68` confirms the active/archived target distinction. Only the ordinary cancellation example needed changing. |
| 12 | Corrected | `publishStaticSiteOutputs` calls `Lstat` and refuses non-regular targets before allocating/writing staged output. |
| 13 | Corrected | `hookcommands/session_start.go` emits the unfinished-finalization warning now included in README. |
| 14 | Corrected | `summaryRequiresContainment` inspects the summary itself. REQ-460's recorded implementation preserves that byte-class policy; prose now delegates to canonical answer independently of position. |
| 15 | Corrected | `publication/answer.go` substitutes `See contained answer note` for a contained summary, while safe summaries stay inline. |
| 16 | Already obsolete | REQ-539 commit `b450a857` replaced the monolithic contract suite with an owner dispatcher; the old once-only/position assertion is absent. No guard was recreated. |
| 17 | Corrected | `answer.go` comments describe the selected digit-led and byte-class checks, removing the false claim to cover every Markdown dialect. |
| 18 | Corrected | Canonical writer/reader use the `Discarded: ` prefix. Both named action statements now include the colon; classification behavior is unchanged. |
| 19 | Corrected | `ownedprocess/owned_process_group_unix.go` confirms descendants-first teardown, verified-group fallback, and bare-PID escalation. `gittransaction.configureCancellableProcessGroup` and REQ-543 D-02 confirm runGit's unsupported-platform default cancellation. |


## Decisions

- **D-01:** Preserve the Testing mutex's actual two-operation scope instead of replacing the global surface count.
- **D-02:** Preserve recorded summary classification and process fallback policies; no new classifier, guard, or rule was authorized.
- **D-03:** Qualify the relocated historical Copy lesson around source bytes before clipboard annotations.
- **D-04:** Correct only the active cancellation example; the archive-target early branch makes the neighboring cleanup example correct.
- **D-05:** Leave the adjacent mentionSurface count outside this captured drain; F1 remains report only.

## Heavy Verification Result

Target and execution revision: `628c0e34ef499e34e8b123297d1d8965bad04c7e`. Stored lane plan reproduced exactly. All six selected lanes executed and passed, with no skipped lanes.
- queue-kanban-javascript: exit 0, 8 seconds; executed (fingerprint_mismatch).
- queue-kanban-browser: exit 0, 122 seconds; executed (fingerprint_uncertain).
- do-work-cli-integrations: exit 0, 59 seconds; executed (fingerprint_mismatch).
- staged-skills: exit 0, 43 seconds; executed (fingerprint_mismatch).
- updater: exit 0, 69 seconds; executed (fingerprint_mismatch).
- installer: exit 0, 31 seconds; executed (fingerprint_mismatch).

## Timing

Observed 2026-09-08T22:01:19Z to 2026-09-08T22:15:16Z: 13m 57s total, 12m 35s attributed across 2 events, 1m 22s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 10m 49s | 1 |
| verification-gate | 1m 46s | 1 |

Slowest stage: builder-work / implementation, 10m 49s, outcome success.
Slowest command: verification-gate / merged-tree, 1m 46s, exit 0, bash (2 argv tokens).
