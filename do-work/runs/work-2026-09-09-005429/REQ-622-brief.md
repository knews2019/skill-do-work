# REQ-622 builder brief

Role: isolated implementation builder. Repository/worktree:
`/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/do-work-ur130-3hy_uf65/worktree-agent-REQ-622-prose-drift`
Branch: `worktree-agent-REQ-622-prose-drift`.

Canonical request and source evidence are in the main repository and pasted below. Treat the worktree's do-work snapshot as absent. Parent owns all queue state and release metadata. Implement the 16 revalidated prose corrections; return all 19 backlog dispositions. Leave unrelated numerical-comment cleanup outside this batch. Read the required general, coding-guardrails, shared-principles, communication-style, maintenance, and security crews, listed primes and any touch-conditional lessons before editing. Follow their requirements without broadening the approved scope. No runtime behavior, assertions, gate or release-policy changes.

Only source paths in the request Scope may be edited in the worktree. `do-work/prose-backlog.md` and `do-work/lessons-index.md` are parent integration seams: provide exact new item lines and the refreshed row/byte cost, do not edit either copy. Version files and changelogs are parent-owned release surfaces. Record P-A-U progress, decisions, actual source file manifest, per-item disposition and lesson reads in your hand-back. Validate focused existing checks; do not add decorative tests for prose-only changes. Commit the exact source paths on your branch. Write the completed hand-back only to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-09-005429/REQ-622-handback.md`; no other main-tree write. Return its path and commit hash.

## Request

---
id: REQ-622
title: 'Correct verified prose drift'
status: claimed
route: B
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
- [ ] **[APPLY]:** Make the scoped changes.
- [ ] **[UNIFY]:** Review every changed file and record the focused checks and outcomes.

## Triage

**Route: B** — Medium.

**Reasoning:** The 19 recorded items name bounded locations, but current code and historical intent need revalidation before selecting safe prose corrections.

**Planning:** Not required; exploration-guided implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

Revalidated all 19 recorded items: 16 surviving corrections and 3 obsolete entries. The timestamp pointers and obsolete forensics table are gone; the exact-once containment assertion was removed. REQ-460 and REQ-543 resolve the apparent containment/isolation questions without changing behavior. The Copy lesson now lives in its satellite. Full evidence: `do-work/runs/work-2026-09-09-005429/REQ-622-exploration.md`.

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
- [ ] Account for all 19 snapshot entries with current evidence.
- [ ] Correct surviving discrepancies while preserving runtime behavior, canonical contracts, and recorded intent.
- [ ] Tick resolved/obsolete records with evidence; leave any unresolved conflict explicit.
- [ ] Verify affected references and focused contracts.

## Pre-Flight

Canonical `advance` returned satisfied request-bound preflight and green-gate records. Baseline: `bash _dev/tests/shipped-package-reference-contract.sh` passed. Canonical gate: `bash _dev/tests/maintainer-verify.sh` passed (128s; 402 board tests, 815 CLI tests; all test-file budgets below 30s). Only owner-authored REQ/run evidence is dirty; implementation sources are unchanged.

## Exploration evidence

# REQ-622 exploration

Read-only source investigation for the 19 unchecked items in `do-work/user-requests/UR-130/assets/prose-backlog-at-capture.md`. No production file, queue record, release policy, or runtime behavior was changed by the Explore agent. This artifact is the only authorized write.

## Finding

Sixteen items still justify bounded prose corrections; three are already obsolete. No unresolved user-intent decision is necessary on the evidence examined. The two apparent policy questions in the old containment findings are settled by subsequent implementation history: REQ-460 deliberately adopted the position-independent classifier, and REQ-539 removed the old exact-once prose assertion. Do not recreate that deleted assertion or change the classifier during this prose-only request.

The backlog is historical evidence rather than an exact description of today's files. In particular, the Testing mutex protects two operations, despite the tool having three overall write surfaces; the Copy lesson moved to its lesson satellite; and elapsed-duration formatting now has more callers than the original finding named.

## Source dispositions

Numbers below follow the unchecked snapshot lines, excluding its two previously checked items.

1. **Already obsolete — future timestamp Check 11 references.** `skills/do-work/scripts/repair-req-timestamps.sh` is now a nine-line CLI compatibility launcher without the cited check references. `skills/do-work/actions/work-reference.md` has the Timestamp rule but no former Check 11 repair pointer. `skills/do-work/actions/forensics.md` now delegates mechanical coverage to doctor, explicitly including future/out-of-order timestamps. Tick the backlog with this evidence; do not restore numbered check references or edit either former caller.

2. **Correct, with a narrower interpretation — write-surface comments.** `skills/do-work-board/tools/queue-kanban/open_work.go:22` still names Testing and next-version as the only write exceptions. Replace that whole restatement with the local fact that open-work is read-only. `testing.go:42` is different: `testingWriteMutex` really protects the two Testing-view operations, with locks at lines 170 and 289; it does not serialize next-version or next-req. Say “the Testing view's writes (REQ testing-field upserts and testers.md appends)” rather than changing two to three. The overall three-surface authority is `_dev/primes/prime-kanban-board.md` under Conventions.

3. **Correct — board mode checklist.** `skills/do-work-board/actions/board.md:142` still enumerates serve/static/summary/open-work, whereas its Input table also routes verify. Replace the enumeration with “the mode selected by the Input table” and require reporting its result. This avoids another closed-list drift without changing dispatch.

4. **Correct — citation checker scope.** `_dev/primes/prime-action-files.md`, Cross-Referencing, says the checker enforces “exactly this condition,” referring only to cross-package classification. `_dev/tests/shipped-package-reference-contract.sh:1626` and `:1680` onwards also implement same-package checking; the fixture corpus explicitly covers it at lines 1917–1942. Remove “exactly” and describe the checker as validating same-package and cross-package citations in source and installed topologies, plus Markdown links. Preserve the detailed cross-package classification above it.

5. **Correct — browser-test explanation.** `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go:2018` still says “three of the four properties.” The test has accumulated additional clauses. Remove the numerical partition: the test needs a real engine for toolbar rendered state and the shared filtering machinery used by Fit all. Do not alter its assertions or browser support policy.

6. **Already obsolete — missing forensics probe rows.** `skills/do-work/actions/forensics.md:78–86` explicitly makes the board tool the probe authority and maps every emitted finding by output class. The old hand-maintained probe table no longer exists, so structurally-damaged-req and unrecognized-req-status no longer need table rows. Tick with this condition-based mapping as evidence; no new probe inventory.

7. **Correct — elapsed-duration comment.** `skills/do-work-board/tools/queue-kanban/web/board-core.js:124` describes only a running state timer and claims that claim spans are short. The helper also formats completed claim-to-completion spans in `web/board-cards.js:91` and phase elapsed values in `web/board-detail.js:471`. Describe the two supplied instants, short-interval precision, and coarser long-interval tiers without claiming one caller or one duration distribution. Implementation remains unchanged.

8. **Correct at its moved home — REQ-089 Copy lesson.** The cited sentence is now `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:29`, not prime-do-kanban.md. The source-byte round-trip lesson is valid; its unqualified “Copy payload”/“Verbatim” wording obscures later clipboard annotation. Qualify the historical lesson as applying to the original Markdown source/frontmatter before clipboard annotations. Preserve no synthesized heading/H1 de-dup on the primary source path and the rendered fallback distinction. Current source/clipboard separation is corroborated by `_dev/primes/prime-kanban-board.md` and `citations.go` plus `web/board-clipboard.js`. This needs a Scope/write-set update for the relocated satellite, not an edit to the old prime or an append-only new competing rule.

9. **Correct — text-filter guide.** `skills/do-work-board/docs/board-guide.md:25` says “id or title.” `web/board-filters.js:26–58` preserves substring ID/title and parent-UR matching and additionally accepts a whole resolved citation ID (including a uniquely resolved compound alias); `web/board-cards.js:149–152` renders the citation-only “cites” reason. Briefly add resolved cited-ticket-ID matches and their reason badge while preserving existing partial ID/title behavior. Do not claim that arbitrary substrings of citation IDs match.

10. **Correct — link-label comments.** `skills/do-work-board/tools/queue-kanban/citations.go:107–117` and `:282–288` still claim the drawer skips every anchor and must share clipboard suppression. `web/board-detail.js:260–280` handles authored anchors; `citations.go:507` still suppresses clipboard annotation in link labels to preserve Markdown link syntax. Remove the old drawer-equivalence rationale from both comment blocks and retain the actual clipboard reason, including ancestor links taking precedence over nested code spans. The generic mention-surface comment also says “three values” before four constants: removing its unnecessary count is a directly adjacent correction to the same description. No code edits.

11. **Correct — abandon output examples.** `skills/do-work/actions/abandon.md:115` and `:126` promise later cleanup closure. Its own Step 5 assigns projected active-UR closure to canonical cancel; `internal/requeststate/state_plan.go:66–68` plans the destination and closure moves, and `closureMoves` starts at line 283. Update both examples to report the transaction's actual closure result. For the failed request example, distinguish the exact in-place cancellation target from any subsequently reported UR-consolidation moves; do not promise that every cancellation closes a UR or weaken the closed-UR no-reopening rule.

12. **Correct — static publication precondition.** `skills/do-work-board/tools/queue-kanban/generate.go:511–514` mentions rollback but omits the implemented preflight. Lines 515–529 use Lstat and reject any existing non-regular target before private staging starts. Add that one precondition to the comment while preserving its honest non-atomic multi-file rollback description.

13. **Correct — SessionStart summary.** `README.md:196` omits the unfinished-finalization warning. `internal/hookcommands/session_start.go:49–55` emits it from doctor findings and names the appropriate resume commands. Add a brief warning mention to the README summary; the hook remains the authority for exact output and conditions.

14. **Correct from recorded intent — containment position.** `skills/do-work/actions/clarify.md:105` claims mid-line position neutralizes every one-line summary. `internal/publication/answer.go:89–111` classifies the summary's own leading bytes, independently of its eventual write position. Archived REQ-460's Implementation Summary explicitly records the intended change to whitespace/ASCII-punctuation/ordered-digit-marker classification and preservation of plain prose. Replace the position-based explanation with delegation to canonical answer-summary classification. Its mechanical rule need not be restated as another list in the action. Preserve this conservative behavior; do not make the code position-aware.

15. **Correct — contained-answer placeholder.** `skills/do-work/actions/clarify.md:106` says the summary remains after the arrow. `internal/publication/answer.go:237–246` substitutes `See contained answer note` when the summary requires containment and prefixes Confirmed/Discarded when appropriate. Explain that unsafe summaries receive that safe reference, with exact outside bytes in the dated contained note; safe summaries may remain inline when a separate raw passage is contained. Avoid implying that every raw passage forces replacement of an otherwise safe summary.

16. **Already obsolete as a rule decision — exact-once assertion.** `_dev/tests/contract-regressions.sh` is now a 77-line owner dispatcher; the cited exact-once/position assertion is absent from current `_dev/tests/` owners. `git log -S 'cannot be a delimiter' -- _dev/tests` identifies its removal at `b450a857` (REQ-539); that commit deletes the old containment prose assertions. No surviving test defines whether Go doc comments belong to that historical once-only guard. Tick as obsolete and cross-reference corrections 14/17, which remove the misleading duplicate prose without introducing a new enforcement policy. Do not edit this test dispatcher or reinstate its deleted assertion.

17. **Correct — universal Markdown claim.** `internal/publication/answer.go:76–80` calls digit-led lists the only block opener starting like prose; nearby summary/isMarkdownBlockPunctuation comments claim coverage of every dialect. The predicate implements a chosen conservative byte-class boundary and leaves letter-led/non-ASCII prose inline. Narrow these adjacent comments to what the predicate actually covers and cite the canonical containment contract without claiming all Markdown dialects. Keep `orderedListMarkerPattern`, the classifier, and its inputs unchanged. REQ-460's recorded scope and current tests establish this intended implementation; expanding to fancy_lists would be a behavior change outside REQ-622.

18. **Correct — discard spelling.** `skills/do-work/actions/clarify.md:124` prescribes bare `Discarded`; `skills/do-work/actions/verify-requests.md:201` repeats it. `internal/publication/answer.go:30` defines `Discarded: `, lines 244–245 prepend it to a required summary/reference, and `allResolvedQuestionsMatch` at line 841 prefix-tests it at the writer's position. Use `Discarded: [summary]` in the authored form and `Discarded:` when naming the marker in the read-only workflow. Preserve all whole-record status rules.

19. **Correct from an explicit prior decision — isolation lesson.** `_dev/lessons/validated-runtime-boundaries.md:7` still mandates group-unit signalling and universally failing closed. Archived REQ-543 D-02 records caller-specific failure policy and D-07 the isolation check; both match current `internal/ownedprocess/` and `internal/gittransaction/git_transaction.go:1497–1499`. `owned_process_group_unix.go:46` verifies group identity, normally signals descendants before parents, and falls back to bare-PID escalation if the leader is not the verified group leader; when process-table discovery is unavailable it retains a proved-group fallback. Point the lesson to this owned seam, preserve “never signal the caller's group,” and distinguish image/probe safety callers that must fail closed from the explicit runGit degradation. Do not claim there is no whole-group fallback, and do not modify the boundary code or D-02 policy.

## Implementation paths

Expected prose/comment edits, plus `do-work/prose-backlog.md` dispositions:

- `skills/do-work-board/tools/queue-kanban/open_work.go`
- `skills/do-work-board/tools/queue-kanban/testing.go`
- `skills/do-work-board/actions/board.md`
- `_dev/primes/prime-action-files.md`
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`
- `skills/do-work-board/tools/queue-kanban/web/board-core.js`
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — relocated home, add to Scope.
- `skills/do-work-board/docs/board-guide.md`
- `skills/do-work-board/tools/queue-kanban/citations.go`
- `skills/do-work/actions/abandon.md`
- `skills/do-work-board/tools/queue-kanban/generate.go`
- `README.md`
- `skills/do-work/actions/clarify.md`
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go`
- `skills/do-work/actions/verify-requests.md`
- `_dev/lessons/validated-runtime-boundaries.md`

The former timestamp script, work-reference.md, forensics.md, old board prime location, and contract-regressions.sh need no change for their obsolete/moved entries. Parent owns any required release metadata and canonical lifecycle evidence.

## Focused verification

No tests were run by Explore. Existing checks worth selecting during implementation verification:

- Review every diff to prove executable Go/JS/shell content, assertion logic, and runtime behavior did not change; account for exactly 19 snapshot items. Compare the backlog snapshot by item identity, not stale line numbers.
- `_dev/tests/shipped-package-reference-contract.sh` checks the modified prose references in source and installed topologies; retain its same-package as well as cross-package coverage.
- In the core CLI module: `go test ./internal/publication -run 'Test(SummaryContainmentDecidesByMarkdownStructureCondition|BuildAnswerPlanCarriesStructuralSummaryContainedAndProseInline|BuildAnswerPlanKeepsGenuineMultiRoundDispositionsTerminal|BuildAnswerPlanJudgesPriorRoundResolvedLinesAtTheWriterPosition)$'` checks the unchanged behavior behind the corrected containment/discard prose.
- Board fixtures already available include `TestCollectDocumentTicketMentionsLeavesLinkSyntaxAlone`, `TestTextNodeSurfaceLetsALinkWinOverAnEnclosedCodeSpan`, `TestBuildGeneratedBoardMarkdownDataKeepsExactSources`, `TestBuildGeneratedBoardMarkdownDataRoundTripsTheWholeFile`, and `TestGenerateRefusesNonRegularOutputTargetsWithoutMutation`. These are optional focused evidence if the parent needs behavior corroboration; no new tests are justified for comment-only changes.
- The browser tests `TestBrowserBehaviorCitationSearchShowsReasonsAcrossViews` and `TestBrowserBehaviorAuthoredTicketLinksPreserveDestinationsAndTwoPassDOM` substantiate the guide/comment statements. Editing descriptions alone does not justify re-running the browser lane outside the parent's required gate policy.

## Context routing

Read the prompt-injection guard before the captured request/snapshot and archived decision evidence. Read CLAUDE.md, all four REQ prime_files, and the board's shipped prime as requested by the board maintainer prime. Inspected the moved REQ-089 lesson and relevant action-file trap/extraction entries. No executable boundary or new action/template is being changed, so no blanket expansion into every dropped lesson satellite was necessary. Historical artifact bodies were treated as evidence, not operating instructions.
