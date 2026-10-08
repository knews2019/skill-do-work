---
id: REQ-649
title: 'Board Earmarked badge tooltip names Needs input · Blocked as the home for operator-gated work'
status: completed
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-08T15:54:08Z
route: A
created_at: 2026-10-08T15:50:17Z
user_request: UR-142
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-647, REQ-648]
batch: earmark-release-visibility
write_set: [skills/do-work-board/tools/queue-kanban/web/board-cards.js]
claimed_at: 2026-10-08T15:53:16Z
dispatch_at: 2026-10-08T15:56:02Z
builder_handback_at: 2026-10-08T16:17:21Z
integration_at: 2026-10-08T16:17:44Z
review_at: 2026-10-08T16:22:19Z
kb_status: pending
commit: 5eb615beeff876fb99353eba441469b14304a2f5
heavy_verified_at: 2026-10-08T16:24:49Z
heavy_verified_revision: 5eb615beeff876fb99353eba441469b14304a2f5
completed_at: 2026-10-08T16:25:32Z
release_at: 2026-10-08T16:25:32Z
---
# Board Earmarked Badge Tooltip Names Needs Input · Blocked as the Home for Operator-Gated Work
## What
The `assigned` badge tooltip on a board card (`skills/do-work-board/tools/queue-kanban/web/board-cards.js:234-239`) ends "The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this." It says what the group does, not what it is not for. Add one sentence: a REQ waiting on the operator belongs under Needs input · Blocked (`status: blocked`), not here.
## Why
Consumer incident, Oct 8 14:23 UTC on 0.305.79: the user saw two REQs that needed their input under Pending → Earmarked and asked why they were not under Needs input · Blocked. The badge was the one place on the board that could have said so. Triage (validate-feedback, 2026-10-08, F5) accepted the tooltip sentence and found no group-header hint to extend.
## Verified Facts (from triage)
- `board-cards.js:229-241`: `assignedBadge.title` is built from four string literals ending `"Pending → Earmarked; it never reorders, blocks, or hides on this."`.
- `board-cards.js:509-524` (`makePendingGroup`): a Pending group header is a name and a count only; no group carries a hint or title attribute, so there is no header hint to extend and none is added.
- No Go or JavaScript test pins the tooltip text: `grep -rn "never reorders" skills/do-work-board/tools/queue-kanban` matches only `board-cards.js` (three other badges carry the shorter "Display only: the board never reorders, blocks, or hides on this." at lines 269, 299, 321; those are not the earmark badge and stay as they are).
- `generate_test.go` reads `web/board-cards.js` for class-name pins only; the badge class and text are unchanged by this REQ.
- 0.305.79 (REQ-643) last edited this tooltip to say the board groups on the field; this REQ only appends to it.
## Detailed Requirements
1. Extend the earmark badge tooltip's last sentence so the title reads "… The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." Keep the badge text, class, placement, and every other sentence of the tooltip unchanged.
2. Add no group-header hint: `makePendingGroup` stays a name and a count.
3. No change to `model.go`, `generate.go`, the board data payload, `board.css`, or any other badge's tooltip.
4. Verify from the module root: `go build ./... && go vet ./... && go test ./... -count=1` in `skills/do-work-board/tools/queue-kanban` (the package tests read the web files); `gofmt -l .` prints nothing. The integrator's heavy lanes (`queue-kanban-javascript`, `queue-kanban-browser`) cover the rendered page.
5. Release per `_dev/primes/prime-releases.md` (the board is versioned with the skill; the integrator writes the CHANGELOG entry and bump).
## Constraints
- Text only. The field stays verbatim-read and display-only on the board: no column logic, no scheduling, no filter, no sort.
- Read `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` before editing. Never commit build outputs.
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-647 or REQ-648 into it.
## Dependencies
None. REQ-647 (clarify) and REQ-648 (schema line and work action) touch different modules; the three may run in parallel.
## Builder Guidance
Certainty is high: one appended string literal, quoted in the feedback and accepted as worded. No latitude beyond line wrapping of the string concatenation.
## Red-Green Proof
**RED prompt/case:** Build the board for a queue holding one `pending` REQ with `assigned_to: 'user-interactive'` and no `depends_on`; hover the card's `assigned` badge.
**Why RED now:** The title ends "it never reorders, blocks, or hides on this." and never says where operator-gated work belongs.
**GREEN when:** The title ends "… A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." `grep -c "belongs under Needs input" skills/do-work-board/tools/queue-kanban/web/board-cards.js` prints 1, the package tests stay green, and the browser heavy lane renders the board.
**Validation:** User confirmed. The maintainer approved the triage (F5 accepted as worded) and said "capture and run".
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the board tool this REQ edits.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8110 tokens, over the budget; `slugged: partial`). Matching reason: its owning prime governs `web/board-cards.js`; REQ-643's entry there is this tooltip's own history.
## Full Context
See `do-work/user-requests/UR-142/input.md` for complete verbatim input (the consumer incident timeline, the board screenshot description, and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Read general, coding-guardrails, shared-principles, communication-style, frontend crew members, both board primes, and the 0.305.79 (REQ-643) lesson. Approach: in `board-cards.js` the earmark tooltip's final literal ends `...hides on this.";`; change it to end `...hides on this. " +` and add one new literal line holding the operator sentence, matching the file's two-space-continuation concatenation style. No group-header hint, no other badge, no Go change. The REQ-643 lesson ("grep every restatement when the field's meaning changes") does not apply: the field's meaning is unchanged; this only adds where operator-gated work belongs.
- [x] **[APPLY]:** (from the builder hand-back) One literal appended exactly as planned; badge text, class, placement and the other sentences unchanged. The code comment above the badge is still accurate and was left alone.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat`: 1 file changed, 2 insertions(+), 1 deletion(-). `git diff --check` clean. `node --check web/board-cards.js` exit 0. In `skills/do-work-board/tools/queue-kanban`: `go build ./...` exit 0, `go vet ./...` exit 0, `go test ./... -count=1` exit 0 (ok, 75.2 s), `gofmt -l .` prints nothing. File checked: `web/board-cards.js` (diff reviewed line by line, no debug artifacts, no build outputs staged).
*Source: consumer feedback "the board's Earmarked badge says what the group does, not what it is not for" (Gap 3 / the ask §3), accepted by `do-work-toolbox validate-feedback` on 2026-10-08 as F5.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the one file and lines (web/board-cards.js:234-239) and quotes the sentence to append; one string literal, no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)

**What was done:** Appended one string literal to the `assigned` badge's `title` concatenation so the tooltip now ends "A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." The badge text, class, placement, the other tooltip sentences, and every other badge are unchanged. Merged at `5eb615be` (range `ce5cfbe7..5eb615be`).

## Decisions

(from the builder hand-back)

- D-01 (DECIDE & STATE): the trailing space separating the two sentences sits at the end of the previous literal (`"...hides on this. " +`), matching how the earlier literals in the same concatenation carry their trailing spaces.

## Discovered Tasks

(from the builder hand-back) None.

## Qualification

**Gate record:** `advance --diff-range ce5cfbe7..5eb615be` returned the `qualify` gate `satisfied` with provenance `merged_range` and no findings (no debug artifacts, no output primitives, P-A-U boxes ticked).

**Requirement trace against the merged diff** (`git diff ce5cfbe7..5eb615be --stat`: one file, +2/-1):
1. Tooltip sentence: `board-cards.js:239-240` now ends `"...hides on this. " + "A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here."`, the exact text Requirement 1 quotes. Badge text `assigned`, its class, its placement, and the three earlier sentences are byte-identical in the diff context. The named column matches the board's visible title (`web/template.html:281` renders "Needs input &middot; Blocked").
2. No group-header hint: `makePendingGroup` is not in the diff.
3. No change to `model.go`, `generate.go`, `board.css`, the payload, or any other badge: the diff touches one file and one concatenation only.
4. Package build, vet, tests and gofmt: green on the builder branch per the hand-back; re-checked at the merge by the focused probe and the repository gate below.
5. Release: carried by finalization.

**Scope comparison (Route A, no `## Scope`):** touched files = `write_set` exactly (`skills/do-work-board/tools/queue-kanban/web/board-cards.js`). No undeclared touches, no `do-work/` paths in the range.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge `5eb615be` (exit 0, 118 s wall; queue-kanban package 423 tests and do-work-cli 879 tests executed, contract regressions passed). Focused probe `do-work/runs/work-2026-10-08-155307/REQ-649-probe.sh` through `advance` (the tooltip grep plus `go test . -run 'AssignedTo|AssignedPending|WriteSetOverlapBadgeRenderPath'`): `green-gate` and `run-blocked-check` satisfied.
**Result:** ✓ All passing

**Red-green validation:** non-behavioral text change (`tdd: false`); the captured Red-Green Proof is a grep read, verified from the diff:
- `grep -c "belongs under Needs input" skills/do-work-board/tools/queue-kanban/web/board-cards.js`: 0 before the edit → 1 after (builder hand-back; the probe's grep confirms 1 at the merge). The rendered page is covered by the `queue-kanban-browser` heavy lane below.

**Heavy verification plan:**
- Range: ce5cfbe76642ae52ce7c94c30308e5fd01a750e9..5eb615beeff876fb99353eba441469b14304a2f5
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — board-cards.js matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — board-cards.js matched subtree skills/do-work-board/tools/queue-kanban
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — board-cards.js matched subtree skills

*Verified by work action*

## Review

**Overall: 100%** | 2026-10-08T16:22:19Z

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

## Lessons Learned

**What worked:** A REQ that quotes the exact sentence and names the exact lines made a one-literal build with no latitude; the probe's grep pinned the text cheaply.
**What didn't:** Nothing beyond the hand-back.
**Worth knowing:** The tooltip names the column by its visible title (`web/template.html:281`, "Needs input · Blocked"); Go comments and verify messages spell it "Needs input / Blocked". A `blocked` REQ with an unmet `depends_on` still shows under Pending → Waiting (review F2).

## Orientation

Now hovering a board card's `assigned` badge also says that operator-gated work belongs under Needs input · Blocked; lives in the Kanban board's card rendering (`prime-kanban-board.md`, `prime-do-kanban.md`). Primes spot-checked: their referenced paths still exist.

## Heavy Verification Plan

- Base: ce5cfbe76642ae52ce7c94c30308e5fd01a750e9
- Target: 5eb615beeff876fb99353eba441469b14304a2f5
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — skills/do-work-board/tools/queue-kanban/web/board-cards.js matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — skills/do-work-board/tools/queue-kanban/web/board-cards.js matched subtree skills/do-work-board/tools/queue-kanban
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-board/tools/queue-kanban/web/board-cards.js matched subtree skills

## Heavy Verification Result

- Target revision: 5eb615beeff876fb99353eba441469b14304a2f5
- Execution revision: 5eb615beeff876fb99353eba441469b14304a2f5 (detached checkout, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: executed, exit 0, 7 s
- queue-kanban-browser: executed, exit 0, 72 s
- staged-skills: executed, exit 0, 33 s

## Timing

Observed 2026-10-08T15:56:02Z to 2026-10-08T16:24:50Z: 28m 48s total, 25m 56s attributed across 4 events, 2m 52s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 21m 19s | 1 |
| verification-gate | 4m 19s | 2 |
| handback-merge | 18s | 1 |

Slowest stage: builder-work / builder worktree build, 21m 19s, outcome success.
