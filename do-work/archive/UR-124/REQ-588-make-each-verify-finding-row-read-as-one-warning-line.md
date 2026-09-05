---
id: REQ-588
title: 'Addendum: make each verify-finding row read as one warning line, not a paragraph'
status: completed
created_at: 2026-09-05T14:45:41Z
user_request: UR-124
addendum_to: REQ-579
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-579]
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
status_changed_at: 2026-09-05T17:00:29Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T17:04:59Z
  basis:
    - trivial short-circuit
route: A
dispatch_at: 2026-09-05T17:05:11Z
builder_handback_at: 2026-09-05T17:24:18Z
integration_at: 2026-09-05T17:24:19Z
review_at: 2026-09-05T17:33:00Z
commit: ab251f240f8218559b75138f48b9c517ea977c23
heavy_verified_at: 2026-09-05T17:36:23Z
heavy_verified_revision: ab251f240f8218559b75138f48b9c517ea977c23
claimed_at: 2026-09-05T17:02:32Z
completed_at: 2026-09-05T17:36:24Z
release_at: 2026-09-05T17:36:24Z
---

# Addendum: Make Each Verify-Finding Row Read as One Warning Line, Not a Paragraph

## What

REQ-579 (render verify findings and skipped probes as compact rows in one list) replaced the finding cards with rows, but each row now renders as one long grey paragraph: the detail sentence, an arrow and the remedy sentence flow together and wrap under the chip, with one orphaned word on the second line. The subject heading, the chip and the row text also sit at three unrelated sizes and weights. Change the row layout so each finding reads as one warning line with its remedy visibly separated, using the layout the user picks from the mock-up report in `ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/`.

## Prior Implementation

REQ-579 shipped in commit b169396e (builder branch commit 1dc13ef7, board 0.236.20). The producer (`verify.go`) gained a `Subject` on `VerifyFinding`, mirrored as `subject` in the board payload by `generate.go`. The client (`renderVerifyFindingsStrip` in `web/board-cards.js`) groups findings by exact subject, prints a `.board-findings-subject` heading per group, then one `.board-findings-row` per finding: a `.board-findings-chip` span followed by a `.board-findings-text` span holding the detail, the optional "cleanup can fix" tag and the remedy prefixed with an arrow, all inline. Skipped probes are rows with a "not checked" chip and the muted class. `web/board.css` lays the row out as `display: flex; align-items: baseline` with the text span wrapping beside the chip. The Node behavior lane (`javascript_behavior_c_test.go`, the `board-findings` cases around lines 2600 and 3030) asserts the list shape, the subject heading, and which rows carry the muted class.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Port mock-up M1's five CSS rules onto REQ-579's `.board-findings-*` block, drop the arrow in the renderer, rewrite every subject-bearing detail in the producer, and write both RED tests before any edit. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-170800/REQ-588-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`79107190`) touching the five write-set files plus the one test file D4 forced (D-07).
- [x] **[UNIFY]:** `git diff --stat` reviewed (6 files, 311+/88-); `gofmt -l .` empty after one reformat; `go vet ./...` clean; debug-artifact scan over added lines empty; per-file checks listed in the hand-back.

## Why

The user's words: "verify findings styling is broken". The strip exists to be scanned as a list of warnings; a row that wraps as a paragraph has to be read, not scanned, and the remedy (what to do) is the part that gets buried at the end of the wrap.

## Requirements

- **D1, remedy separated from detail.** The remedy never continues the detail sentence inline. It sits on its own line under the detail (M1, M2) or is revealed on demand (M3), per the answer to Q1 below. The arrow prefix goes away when the remedy has its own line; a line does not need a pointer to say it follows.
- **D2, chip and text on one grid.** The row is a two-column grid (chip, text) so wrapped text stays in the text column and never lands under the chip. In M2 the grid has a third leading column for the subject, and every chip in the strip aligns in one column.
- **D3, one type scale.** Subject heading, chip and row text use one scale: subject in the board's mono face at the row text size and semibold ink-base (it is an identifier, like a REQ id on a card); chip one step smaller, uppercase, letter-spaced; row text at the size REQ-579 set; remedy in ink-soft. No fourth size.
- **D4, the producer stops repeating the subject.** When `Subject` is set, the detail sentence no longer starts with that same name (today: "worktree-agent-REQ-573-activity-drawer exists — .git/...", under a heading that already says the name). Change the wording in `verify.go` for the subject-bearing findings and update `verify_test.go`. Grouping stays by the payload field, never by parsing text.
- **D5, everything else from REQ-579 stands.** One list, weight only from `fixable` and "not checked", grouping by subject, hide rules, no card classes. The Node behavior lane's existing assertions keep passing; the assertion for the chosen layout is added beside them (see Red-Green Proof).
- Board changes follow `_dev/primes/prime-kanban-board.md` (versioning, parser lock-step, build outputs).

## Open Questions

- [x] Which row layout from the mock-up report should ship? → M1, remedy on its own line under the detail, chip and text on a two-column grid, subject headings kept. Picked from the mock-up report ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/ over M2 (subject as a row-label column) and M3 (remedy behind a toggle); both stay out of scope.
  Recommended: M1 (remedy on its own line under the detail, chip and text on a two-column grid) — the smallest change that fixes all three defects, and it keeps REQ-579's subject heading and row order untouched.
  Also: M2 (subject as a row label in a leading column, chips aligned in one column, remedy under the detail) — the most scannable, costs a small renderer change; M3 (one-line rows, remedy revealed by a per-row toggle) — the smallest strip, costs a click to read what to do.

## Red-Green Proof
**RED prompt/case:** Serve the board with two verify findings under two subjects (the state in Screenshot 1) and look at the strip, or render the same payload in the Node behavior lane and inspect one `.board-findings-row`.
**Why RED now:** The row's text span holds detail, arrow and remedy as inline siblings; on a 2000 px wide board the second finding wraps to a second line with one word on it, and the remedy has no separation from the detail.
**GREEN when:** For the chosen layout, the Node lane asserts the row's structure (M1/M2: the remedy is a block-level element after the detail with no arrow text; M2: the subject is a cell of the row, not a heading above it; M3: the remedy is inside a toggle element that is closed by default) and the board screenshot shows each finding as a chip, one detail line and one separated remedy line, with subject, chip and text at the scale in D3.
**Validation:** User confirmed the defect from Screenshot 1 ("verify findings styling is broken") and asked for mock-ups to choose the layout; the choice is recorded as the answer to Q1.

## Assets

- Screenshot 1 (not saved: the attachment cache had already expired when capture ran): the board at 127.0.0.1:8090, Board view, light theme, 14:25 UTC on 2026-09-05. Strip header "VERIFY FINDINGS 2 findings queue and process problems queue-kanban verify detects — each names what to do about it". Two subject headings (worktree-agent-REQ-573-activity-drawer, worktree-agent-REQ-582-arrow-citations), each with one row: an uppercase chip (UNMERGED-WORKTREE-LEFTOVER, WORKTREE-PRESENT-RUN-IN-FLIGHT), then detail, arrow and remedy as one grey paragraph wrapping to a second line with one orphaned word. The mock-up report's M0 page rebuilds this state from the same payload text.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (6224 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4959 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "verify findings styling is broken, which req should fix it?" / "capture it and fix it" / "also it's make do-work ai-report with mock-ups so I have options to choose from"*


## Answer Notes

- 2026-09-05 - [ ] Which row layout from the mock-up report should ship?: M1, remedy on its own line under the detail, chip and text on a two-column grid, subject headings kept. Picked from the mock-up report ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/ over M2 (subject as a row-label column) and M3 (remedy behind a toggle); both stay out of scope.
> ```
> M1 remedy under detail (Recommended)
> ```

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the files, the mock-up report fixes the exact CSS (M1), and the renderer and producer edits are each a few lines: a styling change plus two small wording/markup edits, with the answered question closing the only design choice. `effort_estimate: effort-mechanical`.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified)

**What was done:** The finding row is now a chip-and-text grid: the detail is one line, the remedy sits under it as its own block with no arrow, and subject, chip and row text share one type scale with the subject in the board's mono face (mock-up M1). The producer no longer opens a subject-bearing detail with its own subject; the terminal verify report prints the subject before the detail so its lines keep naming the REQ or worktree. Twelve test assertions that found a finding by a substring of its detail now compare the subject field. Merge range 0174cb65..ab251f24; builder branch head 79107190. Builder-authored Decisions (D-01 to D-06), Discovered Tasks and render evidence live in the hand-back at do-work/runs/work-2026-09-05-170800/REQ-588-handback.md.

## Decisions

- **D-07 (orchestrator) — write set extended by one test file.** The builder reported that D4 (the producer stops repeating the subject) breaks one assertion in timeline_test.go that located a claim finding by a substring of its detail. The REQ's own requirement demands that file class, so the touch is accepted, the file is added to the write set, and the Scope contradiction is closed here rather than as a follow-up. Reversible, one line.

## Qualification

**Passed.** Read from the merge range 0174cb65..ab251f24, and the builder's render evidence (four screenshots in both themes) inspected by the orchestrator.

- **D1, remedy separated.** The renderer appends the remedy with no arrow, and the stylesheet makes the text column and the remedy blocks; the screenshots show each finding as one detail line and one remedy line under it. The "cleanup can fix" tag stays inline at the end of the detail line, as the mock-up drew it.
- **D2, chip and text on one grid.** `.board-findings-row` is a two-column grid (`max-content minmax(0, 1fr)`); a long remedy wraps inside the text column and nothing lands under the chip. The now-inert `flex: none` on the chip was removed with a reason (D-05).
- **D3, one scale.** Subject in the mono face at 0.8rem semibold ink-base, chip at 0.66rem, row text at 0.8rem, remedy in ink-soft. The subject heading no longer sits at a fourth size. The change is five rules and their rewritten comments, byte-for-byte the mock-up's diff except for the comment prose.
- **D4, no repeated subject.** Sixteen detail strings rewritten; the rule "no detail starts with its own subject" is pinned by a producer test across eleven categories, with category coverage asserted first so an empty fixture cannot pass silently. The terminal report prints the subject before the detail (D-01), so the CLI line still names the REQ or worktree. Two mid-sentence mentions that carry the evidence were deliberately kept (D-02); read and agreed.
- **D5, everything else stands.** REQ-579's one-row-list test and REQ-578's hide-on-Activity test pass unchanged; the two pass-through hosts, grouping and the detached-row spacing are untouched in the diff.
- The one file outside the captured write set (timeline_test.go, one assertion) is a consequence of D4 and is now in the write set (D-07). Twelve assertions in prior REQs' tests moved from a substring of the detail to the subject field; the behavior change is intentional, so the tests were updated rather than deleted.

Scope: six files, five declared plus the D4-forced test file. No drift beyond D-07.

*Checked by work action*

## Testing

**Tests run (post-merge, detached worktree at ab251f24, the main tree carrying a sibling session's uncommitted do-work/ edits):**
- Focused probe do-work/runs/work-2026-09-05-170800/REQ-588-probe.sh, recorded through advance as the test gate: the Node-lane cases for this REQ, REQ-579 and REQ-578, the verify suite, the timeline claim assertion D4 forced, and the generate payload tests — exit 0, wall 7s.
- Repository gate `bash _dev/tests/maintainer-verify.sh` — exit 0, 2m24s; queue-kanban module 390 tests in 23s, do-work-cli module 762 tests in 65s, slowest file 25.34s against the 30s per-file budget.
- Builder, on the branch before the merge: Node lane 62 tests exit 0 in 6s; full queue-kanban module 398 tests exit 0 in 39s; `gofmt -l .` and `go vet ./...` clean.

**Repository gate retry:** first run exited 1, rerun exited 0. The first run's only failures were the 30s per-file budget on three do-work-cli test files this REQ does not touch (finalization_recovery_test.go 32.13s, defer_gate_test.go 30.89s, finalization_req499_test.go 30.07s) at load average 30 with the sibling session and the review agent active; the queue-kanban module passed in 22s in that same run. The rerun brought every file under budget.

**Red-green validation** (traced to `## Red-Green Proof`):
- javascript_behavior_c_test.go `TestJavaScriptBehaviorVerifyFindingRemedyIsItsOwnLineAfterTheDetail`: ✗ before implementation (remedy text still began with the arrow; remedy not a block; row not a grid; subject off the row's scale) → ✓ after. This is the M1 assertion the proof named: the remedy is a block element after the detail with no arrow text, and the three CSS rules the layout rests on are read from the generated page.
- verify_test.go `TestVerifyFindingDetailsDoNotRepeatTheirSubject`: ✗ before (ten violations across eight categories) → ✓ after; its closing assertion pins that the terminal report still names the subjects (D-01).
- Render fact (grid, one detail line, remedy below, one scale): builder screenshots at 1600 px in light and dark, plus a payload with a fixable finding and a skipped probe, inspected by the builder and by the orchestrator; matches mock-up M1.

**New tests added:**
- `TestJavaScriptBehaviorVerifyFindingRemedyIsItsOwnLineAfterTheDetail` (javascript_behavior_c_test.go)
- `TestVerifyFindingDetailsDoNotRepeatTheirSubject` (verify_test.go)

**Existing tests updated (cross-REQ impact):**
- verify_test.go (REQ-083, REQ-285, REQ-458, REQ-579 and others): twelve assertions that located a finding by a substring of its detail now compare the subject field — intentional, the id moved to `Subject` under D4.
- timeline_test.go (REQ-331 era): one assertion in `TestTimelineProjectionReservesTimeForUntimedClaimedWork`, same move (D-07).

**Heavy verification plan:**
- Range: 0174cb65..ab251f24
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — six changed paths under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — same subtree

*Verified by work action*

## Review

**Overall: 94%** | 2026-09-05T17:33:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- Three worktree-leftover assertions (`verify_test.go:781`, `:829`, `:957`) still match `finding.Detail` for the leftover name and now pass only because a live worktree's path embeds its branch name; the property belongs to `Subject` and a branch-only leftover would fail them — `impact-rule-change` → report only
- The terminal report's unconditional `Subject + " " + Detail` join reads as a stutter or juxtaposition for the checkpoint-ghost, release and calibration categories that D-02 deliberately left naming their subject mid-sentence — `impact-user-visible` → report only

**Nit findings:**
- The calibration detail's `its frontmatter` lost the antecedent the old wording gave it, and a log line has no frontmatter — `impact-negligible` → report only
- Both structural-damage details keep `(the id named here was recovered from the filename)` in a sentence that no longer names an id — `impact-negligible` → report only
- The new producer test's `CHANGELOG.md` case would pass with the subject prefix deleted, because the release details name the file themselves; `REQ-914` and `REQ-915` carry that assertion alone — `impact-negligible` → report only

**Acceptance:** Pass — both focused lanes green at `ab251f24`, and the rendered strip shows the M1 layout with the remedy on its own line and no text under the chip.
**Suggested testing:** 4 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Deciding the design before the build: the mock-up report's M1 `<style>` block was the whole CSS spec, so the builder ported five rules instead of inventing them, and the review could check the shipped values against the approved page line by line.
- Both RED tests asserted coverage before the rule (eleven producer categories; three CSS rules read from the generated page), so a fixture that stopped producing findings or a DOM probe that cannot see `display: block` could not pass by accident.

**What didn't:**
- The first repository-gate run went red on the 30s per-file budget of three do-work-cli test files this REQ never touched, at load average 30 with a sibling session and the review agent running; the rerun was green with every file under 26s. A budget failure under load is not a code failure, and the one-retry rule exists for it.
- Three orchestrator argv mistakes each cost a round trip: `--dispatch-bound` alone makes queue-mode advance a continuation and refuses; a `## Plan` note written before the estimate block trips the section-order classifier; `--diff-range` on the test-gate call is refused as irrelevant input.

**Worth knowing:**
- A verify finding is located by `finding.Subject`, never by a substring of `finding.Detail`: the detail no longer names the subject. Three worktree-leftover assertions in verify_test.go still match the detail and pass only because a live worktree path embeds its branch name (review F1, report only); a branch-only leftover would expose them.
- The terminal `verify` report prints the subject before the detail (D-01), so the board and the terminal say the same thing from one string; the join reads awkwardly for the three categories whose detail names the subject mid-sentence (review F2, report only).

## Orientation

Now the board's Verify Findings strip reads as one warning line per finding — chip, one detail line, and the remedy on its own line under it — and a finding's detail no longer opens by repeating the subject its heading already names; lives in the queue-kanban board tool (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). Prime spot-check: every path both primes name still exists; neither is made stale. No map change.

## Heavy Verification Plan

- Base revision: 0174cb65f218dc99bd51713f1bc3752f9d3c0fb8
- Target revision: ab251f240f8218559b75138f48b9c517ea977c23 (the merge commit, recorded in `commit:`)
- Planned with `plan-heavy-verification --manifest _dev/tests/heavy-lanes.json --base-revision 0174cb65 --target-revision ab251f24`; six changed paths, all under skills/do-work-board/tools/queue-kanban; no uncovered paths.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — every changed path matched subtree skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same reasons
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — same reasons

## Heavy Verification Result

- Target revision: ab251f240f8218559b75138f48b9c517ea977c23
- Execution revision: ab251f240f8218559b75138f48b9c517ea977c23 (detached worktree `.git/work-run-20260905-1708/gate-ab251f24`; the main tree's HEAD had moved on to a sibling session's REQ-583 merge, and the runner refuses a dirty tracked tree, so the lanes ran at this REQ's own merge revision)
- Run: `run-heavy-verification --manifest _dev/tests/heavy-lanes.json --lane queue-kanban-javascript --lane queue-kanban-browser --lane staged-skills` with `QUEUE_KANBAN_BROWSER` naming Chrome, 17:33:08Z to 17:35:43Z
- queue-kanban-javascript: exit 0, executed (fingerprint mismatch), 7s
- queue-kanban-browser: exit 0, executed (fingerprint uncertain, browser runtime), 106s
- staged-skills: exit 0, executed (fingerprint mismatch), 41s

## Timing

Observed 2026-09-05T17:05:11Z to 2026-09-05T17:36:24Z: 31m 13s total, 28m 40s attributed across 3 events, 2m 33s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 19m 08s | 1 |
| verification-gate | 9m 32s | 2 |

Slowest stage: builder-work / implementation, 19m 08s, outcome success.
Slowest command: verification-gate / maintainer-verify, 6m 16s, exit 0, .
