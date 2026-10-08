## Review: REQ-643

**Approve** — the board now puts an earmarked pending REQ under Pending → Earmarked, and the Ready counts in the browser, `summary`, and `open-work` agree with the run's scan. Several restatements of the old "display only, never buckets" contract were left behind, one of them a tooltip users see.
Route B | merge `1e6d9909` (builder commit `3393d267`), range `fd5dbd3b..1e6d9909`

### What's built
- A `pending` REQ with a non-empty `assigned_to` and no unmet dependency moves from Ready to a new Earmarked group in the Pending column. An assigned REQ with an unmet dependency stays in Waiting. Payload, summary, digest, renderer, and two docs carry the third group.
- Still open: eight restatements of the old meaning in files outside or inside Scope (F1–F8 below), and requirement 8 (release) belongs to finalization.

### Decisions / risks for you
- F1 (the assigned badge tooltip still says "Display only: the board never reorders, blocks, or hides on this") contradicts the new placement on the same card. The REQ said "nothing about the badge changes", so the builder correctly left it alone, but the tooltip text is now false. Fixing it means a tooltip-only edit; the badge value and visible text can stay as they are.
- Builder decisions D-01 to D-04 are sound (see Domain Review).

### Findings

**Important:**
- F1 `skills/do-work-board/tools/queue-kanban/web/board-cards.js:233-237` (badge tooltip "Display only: the board never reorders, blocks, or hides on this.") and `:224-225` (comment "Display only: the board never buckets, orders, or hides a card on it"): both say the board never acts on `assigned_to`, yet the card now sits under Earmarked because of it. The tooltip is user-facing. — impact-user-visible → report only
- F2 `skills/do-work-board/actions/board.md:120` (Rules, the parser lock-step paragraph): "`depends_on` drives the Ready/Waiting presentation. Pending REQs partition between those groups" and "stable-sorts Pending Ready and Pending Waiting independently". This is the rule an agent reads before changing the parser. It no longer names `assigned_to` as a placement input or the third group. Lines 111-112 of the same file were updated; line 120 was not. — impact-rule-change → report only
- F3 `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:48` (Stakes: "with one derived presentation rule"), `:39` (Trap: "`model.go` buckets by normalized `status` plus the already-derived unmet-dependency set"), `:15` (Read first: "column bucketing driven by normalized REQ status and derived dependency readiness"). The prime is read before every board change and now omits the `assigned_to` placement input. — impact-rule-change → report only

**Minor:**
- F4 `skills/do-work-board/tools/queue-kanban/generate.go:184-185`: the `AssignedTo` payload field comment says "Verbatim and display only — a card badge and a drawer row, never column or dispatch meaning". The file is in this diff. — impact-negligible → report only
- F5 `skills/do-work-board/tools/queue-kanban/web/board.css:1037-1038` ("Pending sub-groups (Ready / Waiting on dependencies). Rendered only when at least one pending REQ is waiting") and `:1302` (`.badge-assigned`: "Display only — nothing on the board branches on it."). — impact-negligible → report only
- F6 Priority restatements name only two groups: `skills/do-work/actions/work-reference.md:120` ("sorts Pending Ready and Pending Waiting independently"), `:248` ("Pending Ready/Waiting order"), `skills/do-work-board/tools/queue-kanban/model.go:224-225` (Priority comment). The code sorts all three groups. — impact-negligible → report only
- F7 `skills/do-work/docs/work-guide.md:119`: the user guide's earmarking paragraph says only "the board shows an `assigned` badge". It does not say the card moves to Earmarked, so a user may still look under Ready. — impact-user-visible → report only
- F8 Test gap: no test pins a nonzero earmarked count in `summary` or the digest (both fixtures assert 0), and no JS behavior or browser test pins the Earmarked group in `fillPendingColumn` (`priority_browser_probe_test.go` reads `.pending-group` only for priority). The builder flagged the renderer half. — impact-negligible → report only

**Nit:**
- N1 `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md:88` still calls the parse "display-only". It is a decision record, so the change is optional. — impact-negligible → report only
- N2 `model.go:1705` comment line is not wrapped (about 115 characters). The lock-step pair matches in meaning but not in wording: `work-reference.md:114` says "placement only, with **no scheduling**" after "Everything else is display", while `model.go` says "never orders across groups, schedules, or gates". — impact-negligible → report only

### Requirements Checklist

- [x] R1 `PendingEarmarked` bucket; assigned+no-unmet-dep leaves Ready; union is Ready, Waiting, Earmarked; priority sort covers the new group; assigned+unmet-dep stays in Waiting; `BoardColumns` and `AssignedTo` comments updated — delivered (`model.go:1722-1763`)
- [x] R2 `work-reference.md:114` schema line rewritten in the same commit as `model.go` (`3393d267`), no heading changed — delivered
- [x] R3 `pendingEarmarked` in the payload struct and copy site — delivered (`generate.go:148,779`)
- [x] R4 renderer: Earmarked group after Waiting; flat list only when Waiting and Earmarked are both empty; badge value and visible text unchanged — delivered (tooltip text now stale, F1)
- [x] R5 summary `ready to work` excludes earmarked plus a new `earmarked` line; digest `(N ready, N waiting, N earmarked)`; `countOpenWork.PendingEarmarked` — delivered
- [x] R6 old Status-only test replaced by `TestAssignedPendingRequestIsEarmarkedNotReady` through `buildBoard`; verbatim-read and absent-reads-as-empty tests kept — delivered
- [x] R7 `board-guide.md:18` and `board.md:111-112` name the third group — delivered (`board.md:120` missed, F2)
- [ ] R8 version bump and changelog — N/A here; finalization owns it
- [x] Constraints: placement only, no alias map or case folding, no filter/sort/schedule on `assigned_to` — delivered. The only new Go reader is `model.go:1728`; the JS reads `columns.pendingEarmarked`, not `assignedTo`, for placement.
- [x] UR batch constraint: badge stays verbatim-read and its value is untouched — delivered

### Acceptance Testing

**Result: Pass**
- Built the tool at `1e6d9909` into the scratchpad. Red-Green Proof case (one pending REQ, no `depends_on`, `assigned_to: 'user-interactive'`): `summary` shows `ready to work : 0` and `earmarked : 1`; `open-work` shows `pending 1 (0 ready, 0 waiting, 1 earmarked)`; `generate` payload has `"pendingReady":[],"pendingWaiting":[],"pendingEarmarked":["REQ-001"]`.
- Rendered the generated static board in headless Chrome (`--dump-dom`). The Pending column shows groups `Ready 0` with "Nothing ready — everything here is waiting or earmarked", then `Earmarked 1` with REQ-001 and the `assigned user-interactive` badge. No Waiting group (D-03).
- Edge cases, one queue: whitespace-only `assigned_to` → Ready (trim matches the core selector's `strings.TrimSpace`); `blocked` + assigned + no deps → Needs input · Blocked (unchanged); `blocked` + assigned + unmet dep → Waiting; `pending` + assigned + a dependency that is already met → Earmarked; `claimed` + assigned → Claimed; priority `now` sorts before `later` inside Earmarked. Rendered groups: Ready 1, Waiting 1, Earmarked 3.
- `go test . -run 'AssignedTo|Earmark|PendingColumns|Bucket|OpenWork|BlockedDependencyGate|Priority'` → ok; `go vet .` clean; `gofmt -l .` empty.
- Lock-step pair (check 2): `model.go` `AssignedTo` comment and `work-reference.md:114` both changed in `3393d267` and both state the Earmarked placement with no scheduling (wording difference, N2).
- Badge (check 3): value, visible text, truncation, and the drawer row are byte-identical to before; the verbatim-read tests are unchanged.
- Nothing schedules, filters, or sorts on `assigned_to` (check 4): `git grep AssignedTo` in the board package finds the parse, the placement case, the payload copy, and `verify.go`'s existing probe only.

### Suggested Additional Testing

- Browser lane (`queue-kanban-browser`): with Ready, Waiting, and Earmarked all non-empty, check the three group headers and their order, the column count, and the column copy button count. Also check the case with only Earmarked cards (Ready 0 with its empty text, no Waiting header). Then check the case with a text filter that hides every earmarked card (the group should disappear, no empty Earmarked header).
- JavaScript lane: confirm no existing behavior test calls `fillPendingColumn` with the old three-argument signature (none found by grep).

### Domain Review

- D-01 (order Ready, Waiting, Earmarked): matches the REQ; pinned by the union-order assertion.
- D-02 (Ready empty text "… waiting or earmarked"): correct. With only Earmarked cards, the old text "everything here is waiting" would be false.
- D-03 (Waiting renders only when non-empty): correct and needed. The grouped branch is now reachable with Waiting empty, and an empty Waiting header with no text would be noise. Existing boards are unaffected because before this change that branch always had a non-empty Waiting.
- D-04 (partition test includes earmarked): correct. It keeps `ready + waiting + earmarked = pending` true as an invariant, which is the paired-predicate-drift guard.
- Test placement (check 5): the replacement test writes three files into a temp `do-work/queue`, runs `buildBoard`, and asserts `PendingReady`, `PendingEarmarked`, `PendingWaiting`, and the `Pending` union order. It pins columns, not status. RED was observed at assertion level by the builder and again by the integrator.
- Testing crew rule "test at the caller seam": satisfied for Go (`buildBoard`); not for the renderer (F8).

### Scores (on the record — not the headline)

**Overall: 92%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | R1–R7 delivered; R8 belongs to finalization |
| Code Quality | 85% | Clean threading of the existing pattern; stale comments left in two touched files (F1 comment, F4) |
| Test Adequacy | 85% | Placement test goes through the board build with real RED; no nonzero-count or renderer pin (F8) |
| Scope | 100% | Exactly the ten declared files |
| Risk | Low | Display-only change; the payload adds a field; old clients ignore it |
| Acceptance | Pass | CLI, payload, and headless render match the Red-Green Proof |

### Restatement Sweep

Elements this diff redefines: (E1) what `assigned_to` does on the board (now a placement input: Pending → Earmarked); (E2) what the Pending column and "Ready" mean (three groups; Ready excludes assigned REQs); (E3) the `summary` and `open-work` count shapes. Searched with `git grep` outside `do-work/`, CHANGELOG files, `ai-reports/`, and `kb/` (history). Stale: F1, F2, F3, F4, F5, F6, F7, N1. No test or tool parses the old digest shape `(N ready, N waiting)` (`_dev/tests` has no match).

Inherited from REQ-644 (capture separates a session earmark from work that needs the user at the keyboard): the earmark-vs-blocked boundary and the External-condition look-alike list. `work-reference.md:114,511`, `capture-reference.md:41`, `capture.md:107,108,143`, and `work-guide.md:119` (its last sentence) still agree. REQ-643's change is consistent with them: an earmarked REQ goes to Pending → Earmarked, and a REQ that needs the operator is still captured `blocked` and lands in Needs input · Blocked.

### Follow-ups created
None (10 findings report only)

## Review

**Overall: 92%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `web/board-cards.js:233-237` assigned-badge tooltip ("Display only: the board never reorders, blocks, or hides on this") and comment `:224-225` ("never buckets") contradict the new Earmarked placement — impact-user-visible → report only
- F2 `skills/do-work-board/actions/board.md:120` parser lock-step rule still says `depends_on` drives a Ready/Waiting split, Pending partitions between those two groups, and the priority sort covers two groups — impact-rule-change → report only
- F3 `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:48` Stakes ("one derived presentation rule"), `:39` Trap, `:15` Read-first omit `assigned_to` as a placement input — impact-rule-change → report only

**Minor findings:** F4 `generate.go:184-185` `AssignedTo` payload comment "display only … never column … meaning" — impact-negligible → report only; F5 `web/board.css:1037-1038` (sub-groups "Ready / Waiting", rendered only when something waits) and `:1302` ("nothing on the board branches on it") — impact-negligible → report only; F6 priority restatements name two groups (`work-reference.md:120,248`, `model.go:224-225`) — impact-negligible → report only; F7 `skills/do-work/docs/work-guide.md:119` says only "the board shows an `assigned` badge", not the Earmarked placement — impact-user-visible → report only; F8 no test pins a nonzero earmarked count in summary/digest or the Earmarked group in the renderer — impact-negligible → report only. Nits: N1 ADR-018:88 "display-only" parse (decision record) — impact-negligible → report only; N2 `model.go:1705` unwrapped comment line, and the lock-step pair matches in meaning but not wording — impact-negligible → report only
**Acceptance:** Pass — the built CLI gives `ready to work : 0`, `earmarked : 1`, digest `(0 ready, 0 waiting, 1 earmarked)`, payload `pendingEarmarked: ["REQ-001"]`; headless Chrome renders Ready 0 + Earmarked 1; edge cases (whitespace value, blocked, unmet and met deps, priority) place correctly
**Restatement sweep:** redefined `assigned_to`'s board meaning (placement in Pending → Earmarked), the Pending column / Ready meaning (three groups, Ready excludes assigned), and the summary/digest count shapes — stale: F1–F7, N1. Inherited REQ-644 elements (earmark-vs-blocked boundary, External-condition look-alike list) re-checked at `work-reference.md:114,511`, `capture-reference.md:41`, `capture.md:107,108,143`, `work-guide.md:119`: still agree
**Suggested testing:** 2 items
**Follow-ups created:** None (10 findings report only)

*Reviewed by review-work action*
