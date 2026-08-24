---
id: REQ-366
title: 'Keep dependency-gated blocked REQs out of Needs Input · Blocked'
status: claimed
claimed_at: 2026-08-24T20:17:59Z
status_changed_at: 2026-08-24T20:17:59Z
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
write_set: [skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/model_test.go, skills/do-work-board/actions/board.md]
---

# Keep Dependency-Gated Blocked REQs Out of Needs Input · Blocked

## What

Change the board's column bucketing so a `status: blocked` ticket with at least one unmet `depends_on` renders in the PENDING column's waiting-on-dependencies group (keeping its blocked badge) instead of NEEDS INPUT · BLOCKED, and enters NEEDS INPUT · BLOCKED only once every dependency has completed. Presentation and counting only — no frontmatter, probe, or scheduling change.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

**Route: B** — The user-visible rule and three-file target are firm, but implementation must trace the existing dependency annotation into shared column bucketing, preserve status-specific exceptions, and prove inherited counters plus rendered badges with focused and browser evidence.
