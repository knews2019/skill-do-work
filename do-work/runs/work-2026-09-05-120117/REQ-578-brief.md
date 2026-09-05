# Builder brief — REQ-578

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-578-findings-strip`
- **Your branch (already checked out there):** `worktree-agent-REQ-578-findings-strip`
- **Route:** A
- **Base commit:** 09a13839 (main)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ below. If you need one, stop and report it to the orchestrator in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (the REQ is `tdd: true`)

Also read every path in the REQ's `prime_files`, and the `lessons-<name>.md` satellite beside each prime whose Read-first or Traps entries your change touches.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the declared write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes), verify no debug artifacts (`console.log`, `debugger`, stray `TODO`) in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. Use:
- Node lane: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...`
- Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-120117/REQ-578-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, and for a `tdd: true` REQ the RED observation (test name + failure) and the GREEN observation.
- `## Lesson evidence` — each `required_lessons` entry you read (whole-satellite or family-targeted) and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN`, each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings. Do not fix them inline.
- `## Integration seams` — any exact line that belongs in a file outside your write set, with where it goes. The orchestrator applies it.

`tdd: true`. The `## Red-Green Proof` in the REQ is the captured RED/GREEN pair. The Node behavior lane only runs with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on` and node on PATH; run it that way to observe RED and GREEN.

---

# The request

---
id: REQ-578
title: 'Hide the verify-findings strip on the Activity view'
status: claimed
created_at: 2026-09-04T23:58:59Z
user_request: UR-117
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-573]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
claimed_at: 2026-09-05T12:00:56Z
---

# Hide the Verify-Findings Strip on the Activity View

## What

The Verify Findings strip (`#board-findings`, added by REQ-285) sits outside the view panels so it stays visible on every view. On the Activity view it pushes the transitions table down and is not what that view is for. Hide the strip while the Activity view is active and show it again when the reader switches to any other view. The strip's content and its behavior on the other views do not change.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- With the Activity view selected, `#board-findings` is hidden (the `hidden` attribute, matching how the strip already hides itself when there are no findings) even when findings exist.
- Switching to Board, Calendar, Durations, Timeline or Testing shows the strip again exactly as today; the "probe(s) could not run" disclosure under it follows the strip.
- Only the Verify Findings strip is affected. The completion-anomalies strip above it is not part of this request.
- The rule lives in the view-switching code (`board-controls.js`), not in the Activity renderer, so a re-render of the Activity table never touches the strip. Update the template comment that says the strip stays visible in every view.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane, render with two verify findings in `boardData`, switch the view to `activity`, and read `document.getElementById("board-findings").hidden`.
**Why RED now:** The strip is outside the view panels by design and nothing in the view switch touches it, so it stays visible on the Activity view (screenshot 3).
**GREEN when:** `hidden` is true while the Activity view is active and false again after switching back to the Board view with the same findings.
**Validation:** User request from the live board; proof inferred during capture.

## Builder Guidance

The user is certain about the outcome. Keep it to the view switch plus one test; do not restructure the strips.

## Assets

- `do-work/user-requests/UR-117/assets/REQ-578-screenshot-3-activity-view-with-verify-strip.png`: the Activity view at 24h with "175 transitions across 38 REQs in the last 24 hours" and rows for REQ-576, 575, 574, 572 (four rows: work merged, builder handed back, builder dispatched, captured), 506, 570, 573, 505 and others; above the table a Verify Findings strip with two cards (WORKTREE-MERGE-STATE-UNDETERMINED for a REQ-506 worktree, WORKTREE-PRESENT-RUN-IN-FLIGHT for the REQ-570 worktree) and a "1 probe(s) could not run" disclosure.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "remove verify finding from this view"*

---

## Triage

**Route: A** - Simple

**Reasoning:** One view-switch rule in a named file (`board-controls.js`), one template comment, one Node behavior test, all three declared in the write set, with a captured RED/GREEN pair. No discovery needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: direct to builder

*Skipped by work action*

