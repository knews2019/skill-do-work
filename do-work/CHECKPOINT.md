---
session_ended: 2026-08-21T09:08:21Z
last_completed: REQ-308
queue_state: 1 pending, 6 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 22
session_depth: heavy
---

# Session Checkpoint

## Completed This Session

`do-work run` — drained every claimable REQ in the queue, serial mode, one commit per REQ.

- **REQ-279**: Exclude archive assets from blanked-REQ scans — `f487f04`
- **REQ-283**: Route board verification through its skill — `2308afd`
- **REQ-295**: Clarify maintainability audit impact score wording — `0b2261b`
- **REQ-262**: Govern the prompt-kit templates' date headers (Route B) — `24587e5`
- **REQ-269**: Draw the cross-package citation class by what a citation is (Route C) — `f71dfee`
- **REQ-274**: Retire the "the SessionStart hook exits nonzero" framing (Route B) — `0efefa6`
- **REQ-280**: Probe timestamp ordering, point Check 12 at the shipped repair (Route B) — `5e180d0`
- **REQ-281**: Reconcile the calibration log against its own frontmatter (Route B) — `a868827`
- **REQ-284**: Emit every verify finding from the board's Go producer (Route B) — `8f61f69`
- **REQ-285**: Render a verify-findings strip on the board (Route B) — `fed89c9`
- **REQ-288**: Fix the three unfiled contradictions in clarify's Step 4 (Route C) — `c25ee71`
- **REQ-291**: Browser behavior probe lane beside the Node behavior lane (Route B) — `6fa130d`
- **REQ-292**: Move Durations label placement into the browser (Route C) — `ce28510`
- **REQ-266**: Name builds beside the JS renderer's measured face numbers (Route A) — `8fe9eb6`
- **REQ-277**: State the mark-label face constant's real scope at its home (Route B) — `54282b0`
- **REQ-293**: Make the impact/effort lock-in checks pin the property (Route B) — `df976d9`
- **REQ-298**: Sweep the unchecked-exit-status primitive across shipped scripts (Route C) — `01abc28`
- **REQ-299**: Carry builder-authored sections past Step 8 (Route C, 94%) — `3469b39`, **0.233.0**
- **REQ-303**: Run the pinned live-archive assertions only in a suite checkout (Route B, 96%) — `69f3319`, **0.233.1**
- **REQ-304**: Draw a reversed wait as a break, not as a valid bar (Route B, 95%) — `5e08a31`, **0.233.2**
- **REQ-305**: Say the timeline forecast describes the whole queue (Route B, 96%) — `ef0cc55`, **0.233.3**
- **REQ-308**: Judge effort_estimate on every capture, as impact already is (Route B, 95%) — `9bce005`, **0.234.0**

`maintainer-verify.sh` exited 0 immediately before every commit, and each implementation hash was
confirmed with `record-commit-hash.sh --verify`. Serial mode throughout — no worktrees created,
none to clean up.

## In Progress (interrupted)

None. `do-work/working/` is empty.

## Still Queued

- REQ-307: Standing prose-reconciliation sweep (pending, `standing: true` — the default scan never selects it)
- REQ-309: Run the repo's canonical gate before hand-back, not only the changed area's tests (pending-answers — 1 question)
- REQ-310: Check a template payload's citations against where the payload lands (pending-answers — 1 question)
- REQ-311: Resolve the nine calibration-log rows that disagree with their frontmatter (pending-answers — 1 question)
- REQ-312: Resolve same-package citations in the shipped reference contract (pending-answers — 1 question)
- REQ-313: Count the breaks the timeline actually draws (pending-answers — 1 question)
- REQ-314: Judge effort_estimate on review-minted follow-ups too (pending-answers — 1 question)

Nothing is claimable. Six REQs await `do-work clarify`; the seventh is a standing sweep that only
runs on demand.

## Session Notes

- **The canonical gate was red on arrival.** `_dev/tests/record-commit-hash-guards.sh:428` carried
  `\\` where a line continuation was meant, so a `printf` never ran and a bare redirect created its
  fixture at 0 bytes. Red since `2308afd` (REQ-283). Repaired first, because every hand-back in the
  run depends on that exit code.
- **Environment had to be built.** Go 1.26.1 (`/usr/local/go-1.26`) and ShellCheck 0.11.0
  (`/usr/local/bin`) were installed to meet the repo's minimums; `just` came from apt. The standard
  invocation for the full gate is
  `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium-1194/chrome-linux/chrome bash _dev/tests/maintainer-verify.sh`
  with `/usr/local/go-1.26/bin` and `/usr/local/bin` on PATH.
- **Working standard adopted for the run, and kept:** reproduce the captured RED on the untouched
  tree before writing the fix; mutation-test every new check by reverting the fix it guards; measure
  against the real corpus rather than fixtures alone; drive a real browser for anything visual.
- **Cleanup archived seven closed URs** (051, 056, 058, 059, 060, 061, 063), moving 46 loose archive
  REQs into their folders. That broke 38 citations pointing at the old flat paths; all were rewritten
  mechanically and every archive citation in the repo was then verified to resolve. Eleven of the
  repaired files are archived REQ bodies — a path substitution only, no other content touched.
- **`do-work/runs/` holds eight run directories** from 2026-08-18. None were produced by this
  session and none carry a `Status: consumed` marker, so Pass 4 left them alone.

## Context Summary (heavy sessions only)

**Re-read before starting:** `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`,
and `_dev/primes/prime-kanban-board.md`. All three gained lessons this session and their `## Lessons`
lists are the compressed record of what went wrong.

**Decisions worth carrying:**

- **REQ-299 D-01** — a check that needs a list of sections can be turned into a check that every
  section mention states who writes it. The list moves into the prose, co-located with what it
  describes.
- **REQ-303 D-01 / REQ-305 D-03** — the same lesson from two directions: a test that calls the
  function under test directly cannot hold its call site, and an inline `t.Fatalf` chain cannot be
  proven to still bite. Extract findings; drive the real entry point.
- **REQ-308 D-01/D-02** — two write-set extensions, both because a sweep found sites the REQ had not
  listed. Enumerating sites in a REQ is a starting point, never the set.

**Patterns established this session:** the browser behavior lane (REQ-291) and the timeline DOM
probe (REQ-304) are both new harnesses that later work should extend rather than duplicate.

**Open thread:** the six `pending-answers` REQs are the run's output as much as the code is. Three
were created by this session's own reviews (312, 313, 314) and each names why the decision is the
user's rather than the agent's.
