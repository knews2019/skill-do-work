---
id: REQ-366
title: 'Keep dependency-gated blocked REQs out of Needs Input · Blocked'
status: completed
claimed_at: 2026-08-24T20:17:59Z
completed_at: 2026-08-24T20:56:19Z
commit: c18deb8d575988e614492dc73242b96f6b56cb1c
status_changed_at: 2026-08-24T20:56:19Z
route: B
created_at: 2026-08-24T14:03:59Z
user_request: UR-070
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-24T20:17:59Z
  basis:
    - Route B
    - 3-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/actions/board.md
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/docs/board-guide.md
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
---

# Keep Dependency-Gated Blocked REQs Out of Needs Input · Blocked

## What

Change the board's column bucketing so a `status: blocked` ticket with at least one unmet `depends_on` renders in the PENDING column's waiting-on-dependencies group (keeping its blocked badge) instead of NEEDS INPUT · BLOCKED, and enters NEEDS INPUT · BLOCKED only once every dependency has completed. Presentation and counting only — no frontmatter, probe, or scheduling change.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Traced authoritative dependency annotation into the single column-bucketing seam and mapped every exact-status exception plus inherited counter/render consumer before editing.
- [x] **[APPLY]:** Routed only bare blocked tickets with unmet dependencies into Pending/Waiting, added mutation-sensitive model coverage, documented the actionable-inbox invariant, and corrected one directly stale badge comment.
- [x] **[UNIFY]:** Reviewed the exact four-file diff; focused tests, four mutations, full Go module, vet, live CLI/browser transition, canonical verification, formatting, and artifact checks passed.

## Why

The NEEDS INPUT · BLOCKED column is the operator's inbox: "you read it, you do something, it leaves." A blocked REQ whose external condition cannot be satisfied until an upstream REQ lands (e.g. `blocked_by: "owner explicitly approved the architecture decision report from REQ-A"` with `depends_on: [REQ-A]` and REQ-A still pending) is not actionable, yet it sits in the column — possibly for weeks — training the operator to ignore the column.

## Context

Bucketing today is status-only: `isNeedsInputOrBlockedStatus` (`skills/do-work-board/tools/queue-kanban/model.go:1018`) routes `pending-answers` / `blocked` / `blocked-archive-collision` / `blocked-dependency-cycle` / `failed` into the column in `bucketColumns` (`model.go:1589`), and the unrecognized-status default lands there too. Dependency readiness is consulted only for the Pending ready/waiting split (`model.go:1596`), even though `UnmetDependencies` is already computed for every ticket by `annotateDependencyState` (`model.go:1679`, dangling and cancelled deps count as unmet).

Inheritance the builder can lean on:

- Every counter reads `board.Columns`, so one routing change in `bucketColumns` moves them all: the `summary` counts (`main.go:122-126`), the open-work digest (`open_work.go:44-65`), and the Columns-lens JSON payload (`generate.go:607-611`).
- The frontend already renders both badges from ticket data, independent of column: the `blocked by` badge keys on `request.status === "blocked"` (`web/board-cards.js:100`) and the `needs` dep chips key on `request.dependsOn` / `request.unmetDependencies` (`web/board-cards.js:284`). A blocked ticket routed into `PendingWaiting` shows both with no JS badge work expected.
- Timeline and calendar color by status, not by column membership — verify, expect no change there.

## Detailed Requirements

In the user's words:

- "a blocked ticket with unmet depends_on is waiting on its dependency first, not on me. Render it in the PENDING column with the existing waiting-on-deps treatment (NEEDS badge) plus its blocked badge, and promote it to NEEDS INPUT · BLOCKED only when every dependency is completed — the moment its external condition becomes the sole gate and therefore actionable."
- "Apply the same rule everywhere the bucket is counted or rendered: the Board lens, the summary command's needs-input/blocked counter (such tickets count under pending → waiting on deps), the open-work digest line, and the timeline/calendar state coloring if they distinguish the bucket."
- "Invariant to state in the board's docs: NEEDS INPUT · BLOCKED contains only tickets the operator can act on right now." State it in `skills/do-work-board/actions/board.md` and keep the `BoardColumns` comments (`model.go:347-351`) in step.

## Constraints

Scope guards, in the user's words:

- "pending-answers tickets keep their current placement (their questions are answerable regardless of deps — that's why answering them is called clarify, not unblock)."
- "Stakeholder-question REQs keep their current placement (getting answers from the named person is always actionable)." They carry no `depends_on` by construction (`actions/work-reference.md` → stakeholder schema), so the blocked+unmet-deps rule never matches them — no special-casing needed, but do not add any.
- "The unrecognized-status catch-all stays in NEEDS INPUT so nothing ever goes invisible."
- "blocked_check probe semantics are untouched — this is presentation and counting only, no frontmatter or scheduling change, so verify and work behavior are unaffected."
- The rule keys on status `blocked` exactly — `blocked-dependency-cycle` and `blocked-archive-collision` stay in NEEDS INPUT · BLOCKED (breaking a cycle or resolving a collision is operator-actionable).
- Dependency "completed" means terminal success, exactly as `annotateDependencyState` already judges it: `cancelled` and dangling ids stay unmet, matching the pending split's behavior.

## Dependencies

None.

## Builder Guidance

Certainty: Firm — rule, scope guards, and acceptance were stated by the user. The prime file's render-evidence rule applies: generate a board with a dependency-gated blocked ticket and look at it, don't stop at the Go suite.

## Red-Green Proof

**RED prompt/case:** A synthetic queue holds REQ-A (`status: pending`) and REQ-B (`status: blocked`, `blocked_by` naming a condition on REQ-A's output, `depends_on: [REQ-A]`). Build the board: REQ-B is bucketed into NeedsInputOrBlocked; `summary` counts it under needs-input/blocked.
**Why RED now:** `bucketColumns` routes on status alone; `UnmetDependencies` is never consulted for the blocked statuses.
**GREEN when:** The same synthetic queue puts REQ-B in Pending → Waiting on dependencies, rendering both its `blocked by` badge and its unmet `needs` chip, and `summary`/open-work count it under pending → waiting. Flip REQ-A to `completed`: REQ-B now appears in NEEDS INPUT · BLOCKED and the counts move with it. A blocked REQ with no `depends_on`, or with all deps completed, behaves exactly as today.
**Validation:** User confirmed — the acceptance criteria above are the user's own words.

## Assets

- `do-work/user-requests/UR-070/assets/REQ-366-screenshot-1-needs-input-column-with-dependency-gated-blocked-card.png` — the board of a consuming project (glw-game-find-the-difference). NEEDS INPUT · BLOCKED holds one card, REQ-1653 "[impact-critical] Implement the approved Go player backend", status blocked, badge "BLOCKED BY owner explicitly approved the architecture deci…", NEEDS REQ-1652; the detail drawer shows DEPENDS ON REQ-1652 (pending), blocked since 2h 32m. PENDING holds 13 cards, CLAIMED is empty. The card is the concrete non-actionable case: the report to approve does not exist yet.

## Full Context

See `do-work/user-requests/UR-070/input.md` for complete verbatim input.

---
*Source: "is it crazy that I don't want to see anything in that column that I can not act on? how can we make that happen?" plus the follow-up "Board wish" spec — verbatim in UR-070.*

## Triage

**Route: B** — The user-visible rule and three behavioral files are firm, but implementation must trace the existing dependency annotation into shared column bucketing, preserve status-specific exceptions, prove inherited counters plus rendered badges, and keep the directly affected frontend commentary honest.

## Plan

**Planning not required** — Route B: exploration-guided implementation at the shared bucketing seam with focused status/dependency contracts and rendered transition evidence.

## Exploration

- `annotateDependencyState` already owns the unmet set and terminal-success meaning; only completed and completed-with-issues satisfy a dependency, while cancelled and dangling ids remain unmet.
- `bucketColumns` is the single presentation/counting seam. Summary, open-work, generated payloads, and Board columns all inherit its `BoardColumns` result without parallel rules.
- Blocked badges and dependency chips read retained ticket status/metadata rather than column membership, while Timeline and Calendar coloring remain status-driven. No frontend behavior or scheduling/probe semantics need to change.
- The exact-status exception matrix must keep pending-answers, blocked-archive-collision, blocked-dependency-cycle, failed, and unknown statuses in Needs Input even when dependency metadata exists.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/model.go`
- `skills/do-work-board/tools/queue-kanban/model_test.go`
- `skills/do-work-board/actions/board.md`
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/docs/board-guide.md`
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`

**Acceptance criteria:**

- A ticket with exact normalized status blocked and non-empty authoritative unmet dependencies appears in Pending/Waiting while retaining its blocked status and metadata.
- The same ticket moves to Needs Input only after every dependency reaches terminal success; cancelled and dangling targets remain unmet.
- Pending-answers, blocked-archive-collision, blocked-dependency-cycle, failed, unknown, blocked-with-no-deps, and blocked-with-satisfied-deps preserve their specified placements.
- Summary, open-work, generated payloads, docs/comments, and the rendered badges/counts all inherit the same actionable-inbox invariant without changing scheduling, probes, Timeline, or Calendar semantics.
- Every shipped bucketing guide and Pending empty-state message remains true for both pending-status and dependency-gated blocked waiting groups.

## Decisions

- **D-01 — Route at the existing `BoardColumns` seam.** One exact branch makes every current counter and renderer inherit the rule without duplicating dependency logic.
- **D-02 — Key the detour on exact status `blocked`.** The broader Needs-input family remains operator-actionable and must not be swept into dependency waiting.
- **D-03 — Reuse authoritative dependency annotation.** This REQ does not redefine terminal success or inspect raw dependency ids inside bucketing.
- **D-04 — Keep the accepted shared-file extension comment-only.** `board-cards.js` behavior already follows status metadata across columns; only its stale column-specific comment changes.
- **D-05 — Remediate directly caused wording drift before release.** Review found two shipped guides and the Ready empty-state still asserting the old status-only model; these are part of the contract changed here, not unrelated cleanup.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/model.go` (modified): routes exact blocked-plus-unmet tickets into Pending/Waiting and keeps the column contract comments synchronized.
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified): covers pending/cancelled/dangling/satisfied/no-dependency transitions, exact-status exceptions, retained metadata, and inherited summary/open-work/generated-payload counts with four mutation axes.
- `skills/do-work-board/actions/board.md` (modified): defines Needs Input as the presently operator-actionable inbox and documents dependency-gated blocked placement.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified): corrects the blocked-badge comment and keeps the Ready empty-state truthful when only dependency-gated blocked work is waiting; routing behavior is unchanged.
- `skills/do-work-board/docs/board-guide.md` (modified): documents the blocked-plus-unmet Pending/Waiting exception to ordinary status placement.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified): qualifies normalized-status and bare-blocked placement guidance with the authoritative unmet-dependency rule.

## Discovered Tasks

None.

## Testing

- TDD RED placed all blocked fixtures in Needs Input and produced pending 1 / needs 6 instead of pending 4 / needs 3. GREEN passed the focused exact-routing and exception tests.
- Four independent mutations—restored status-only routing, broadened family routing, wrong PendingReady placement, and missing Pending union membership—each failed the intended contract and were restored.
- Full queue-kanban tests, Go vet, formatting, diff checks, and the builder canonical maintainer gate passed, including all 109 prescribed shell cases and strict JavaScript; its optional browser lane skipped because live evidence ran separately.
- A generated fixture showed blocked REQ-101 in Pending/Waiting with blocked and needs badges while REQ-100 was pending, then moved it to Needs Input after REQ-100 completed. CLI summary/open-work changed in step and no application console errors appeared.
- Initial independent review approved behavior at 96/100 with no Important findings and two Minor wording findings: adjacent guides still stated status-only placement, and the Ready empty-state could falsely refer only to pending REQs. Both were accepted for narrow remediation before release.
- Remediation updated both guides and neutralized the empty-state sentence without changing routing. Post-remediation focused routing/exception tests, strict JavaScript, the full queue-kanban module, and the final main-tree canonical gate passed; the canonical browser lane made its standard no-browser skip after separate live evidence.

## Qualification

- Exact cumulative merge range `11574d4..c18deb8d575988e614492dc73242b96f6b56cb1c` passed mechanical qualification.
- Scope drift passed: all six implementation files exactly match the accepted Scope and Implementation Summary.
- Orchestrator judgment confirmed substantive shared-seam routing, complete status/dependency requirement tracing, authoritative unmet-dependency data flow, inherited counts/rendering, and no generated/debug artifacts.

## Review

Independent review first approved behavior at 96/100 with two Minor wording findings and no Important findings. Narrow remediation closed the stale adjacent guides and false Ready empty-state. Final re-review approved with no Important, Minor, or Nit findings: behavior 100, tests/mutation 99, docs/UI 100, overall 99/100, low risk, acceptance complete.

## Lessons Learned

Moving one status shape between presentation buckets changes more than the routing branch: empty states and nearby conceptual guides can encode the old partition too. Review the negative space—the words shown when a subgroup is empty—as part of the state model.

## Orientation

Released in 0.236.55. Exact blocked tickets with unmet dependencies now wait in Pending until their prerequisites finish, keeping Needs Input limited to work the operator can act on now.

## Scope Extensions

- **Pre-freeze comment extension:** `skills/do-work-board/tools/queue-kanban/web/board-cards.js` joins the scope only to correct its blocked-badge comment, which otherwise would falsely claim every blocked card shares the Needs Input column. REQ-367 edits the same file in line-disjoint column-rendering code, so integration must explicitly inspect the overlap.
- **Review wording extension:** `skills/do-work-board/docs/board-guide.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` join to remove direct status-only placement claims made false by this REQ. The already-scoped `board-cards.js` may also neutralize its Ready empty-state sentence so a waiting group containing only blocked cards is described truthfully.
