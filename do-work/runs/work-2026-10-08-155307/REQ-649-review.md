## Review: REQ-649

**Approve** — the `assigned` badge tooltip now ends with the exact sentence the REQ quotes, and nothing else changed.
Route A | merge range `ce5cfbe7..5eb615be` (merge `5eb615be`)

### What's built
- Hovering the `assigned` badge on a board card now ends: "... The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." (`skills/do-work-board/tools/queue-kanban/web/board-cards.js:239-240`).
- The diff is one file, +2/-1. Badge text, class, placement, the three earlier tooltip sentences, the other badges, `makePendingGroup`, `model.go`, `generate.go`, `board.css` and the board payload are untouched.

### Decisions / risks for you
- None. The browser heavy lane (`queue-kanban-browser`) runs after this review. It is pending, not a gap.

### Findings

**Important:**
- None

**Minor:**
- F1 `skills/do-work-board/docs/board-guide.md:18` explains the Earmarked group (pending, no unmet dependency, non-empty `assigned_to`) but never says that work waiting on the operator belongs under Needs input · Blocked. The tooltip, `docs/work-guide.md:119`, `actions/capture.md:107`, `actions/work-reference.md:114` and `actions/work.md:523` all say it now. This is an omission, not a contradiction, and the file is outside the wave. Suggested text to append after "...an assigned REQ that is still waiting on a dependency stays under Waiting.": "Earmarked is for work reserved for another session or checkout; a REQ waiting on you as operator is `status: blocked` and shows under Needs input · Blocked." — impact-negligible → report only

**Nit:**
- F2 `board-cards.js:240` says an operator-gated REQ "belongs under Needs input · Blocked (status: blocked)". A `blocked` REQ that also has an unmet `depends_on` is displayed under Pending → Waiting until the dependency clears (`board-guide.md:18`). The "not here" part stays true, and the maintainer accepted the wording as quoted, so no change is suggested. — impact-negligible → report only

### Requirements Checklist

- [x] R1 Tooltip ends with the quoted sentence; every other sentence, badge text, class and placement unchanged — delivered (`board-cards.js:234-240`; `grep -c "belongs under Needs input"` prints 1)
- [x] R2 No group-header hint; `makePendingGroup` stays name + count — delivered (not in diff)
- [x] R3 No change to `model.go`, `generate.go`, payload, `board.css`, other badges — delivered (diff stat: one file)
- [x] R4 Package build/vet/test/gofmt green — delivered (builder hand-back; repository gate exit 0 at `5eb615be`; reviewer rerun below)
- [x] R5 Release — N/A at review time (finalization writes the CHANGELOG entry and bump)
- [x] UR-142 ask §3 (the badge should say what Earmarked is not for) — delivered

### Acceptance Testing

**Result: Pass**
- `go test . -run 'AssignedTo|AssignedPending' -count=1` in `skills/do-work-board/tools/queue-kanban` at `5eb615be`: 4 tests PASS, `ok` in 0.27 s.
- `node --check web/board-cards.js`: exit 0.
- `grep -c "belongs under Needs input" web/board-cards.js`: 1.
- The named column matches the rendered title (`web/template.html:281`: "Needs input &middot; Blocked").
- Rendered-page check is the `queue-kanban-browser` heavy lane, pending after this review.

### Restatement Sweep

This diff redefines nothing. It adds one more restatement of REQ-648's rule (earmark is for another session or checkout; operator-gated work is `status: blocked`). The trigger set is inherited at the wave end:
- REQ-648 (`assigned_to` audience): `actions/work-reference.md:114`, `actions/capture.md:107`, `actions/work.md:523`, `docs/work-guide.md:119`, and now `board-cards.js:240` agree. `capture.md:143` and `capture-reference.md:41` state mechanics only. Board-side Earmarked wording in `actions/board.md:112,120`, `prime-do-kanban.md:15,39,48`, `model.go:176-190`, `board.css:1297-1303` and the `board-cards.js:220-228` comment describe placement only and do not conflict. `board-guide.md:18` is silent on the operator case (F1). REQ-648's own `work.md:523` "hides" nit is already recorded in its review and is not repeated here.
- REQ-647 (clarify Step 5.5 re-entry now conditional on `assigned_to`): `clarify.md:192`, `docs/work-guide.md:119` (last sentence) agree. `clarify.md:13` is the general "when to use" line and does not conflict.
- No member is unread for the sweep.

### Suggested Additional Testing

- Hover the `assigned` badge on a rendered board and confirm the full sentence shows in the native tooltip (covered by the pending browser heavy lane).

### Scores (on the record — not the headline)

**Overall: 100%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All five requirements met |
| Code Quality | 100% | Matches the file's concatenation style; D-01 documents the trailing-space choice |
| Test Adequacy | N/A | Copy change; no test pins tooltip text, by design |
| Scope | 100% | Touched files = `write_set` exactly |
| Risk | None | Display string only |
| Acceptance | Pass | Focused tests green, grep 1, node check clean |

### Follow-ups created
None (2 findings report only)

## Review

**Overall: 100%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** F1 `skills/do-work-board/docs/board-guide.md:18` explains the Earmarked group but never says operator-gated work belongs under Needs input · Blocked, which five other homes now state (omission, file outside the wave) — impact-negligible → report only. Nit F2 `board-cards.js:240` says an operator-gated REQ belongs under Needs input · Blocked, but a `blocked` REQ with an unmet `depends_on` displays under Pending → Waiting; "not here" still holds and the wording was accepted as quoted — impact-negligible → report only
**Acceptance:** Pass — `go test . -run 'AssignedTo|AssignedPending'` green at `5eb615be`, `node --check` clean, `grep -c "belongs under Needs input"` = 1; browser heavy lane pending
**Restatement sweep:** nothing redefined (the tooltip restates REQ-648's rule without changing it); inherited REQ-648 `assigned_to` audience and REQ-647 Step 5.5 conditional re-entry swept: `work-reference.md:114`, `capture.md:107`, `work.md:523`, `work-guide.md:119`, `clarify.md:192`, `board-cards.js:240` agree; `board.md:112,120`, `prime-do-kanban.md:15,39,48`, `model.go:176-190`, `board.css:1297-1303` placement-only, no conflict; `board-guide.md:18` silent (F1)
**Suggested testing:** 1 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*
