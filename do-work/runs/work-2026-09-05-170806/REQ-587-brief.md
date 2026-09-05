# Builder brief — REQ-587

## Where you work

- **Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-587-give-the-timeline-view-one-scroll-surface`
- **Branch:** `worktree-agent-REQ-587-give-the-timeline-view-one-scroll-surface`, already checked out there and reset to the integration tip `d6b6adb1`. Clean.
- **Commit on that branch.** Do not touch the main tree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` — with exactly one exception, the hand-back file below.

## Why the worktree matters more than usual here

A sibling session is editing `skills/do-work-board/tools/queue-kanban/model.go` and `dependency_test.go` uncommitted in the shared main tree — the same package you are changing. Your worktree is cut from a commit, so those bytes cannot reach you. Never read, measure, or "fix" anything in the main tree. If a test fails, it failed in your tree.

## Never touch

- Any path under `do-work/` except the one hand-back file.
- Anything outside your declared `## Scope` list. If you need a file outside it, stop and report to the orchestrator with the exact line and where it goes.

## Hand-back

Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-170806/REQ-587-handback.md` (absolute path, main tree — the one exception). Never stage or commit it.

Sections it must carry:
- `## File manifest` — every file created/modified/deleted with the verb, plus your branch head commit.
- `## Measured Evidence` — every number you measured in a real engine, each naming the browser and build that produced it. The REQ's GREEN condition is a measurement, not a claim.
- `## P-A-U` — `[PLAN]`, `[APPLY]`, `[UNIFY]` content. `[UNIFY]` means `git diff --stat`, review every changed file, run `go vet` and `gofmt -l`, confirm no debug artifacts, and list each file you checked and what you checked.
- `## Decisions` — continue numbering from `D-14` (the plan already used D-01 through D-13). Sort each by the decide-vs-escalate gate in `crew-members/coding-guardrails.md` § Think Before Coding; an ESCALATE entry carries `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out of scope finds, not fixed inline.
- `## Lesson evidence` — which lesson satellites you read and any that were missing.
- `## Integration seams` — any line belonging in a file outside your write set: the exact line and where it goes.

## Commands

Board module root: `skills/do-work-board/tools/queue-kanban`. Clean baseline before you start: `go -C skills/do-work-board/tools/queue-kanban test -count=1 ./...` exits 0 in about 42s.

`web/` is **embedded, not read at runtime** — CSS and JS edits are invisible until you rebuild the binary. Rebuild before you serve or measure anything.

Browser probes and the browser lane need the engine named explicitly, because Chrome is installed but not on `PATH` under any name the probe looks for:

    QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" QUEUE_KANBAN_BROWSER_PROBES=on go test -count=1 -run 'TestBrowserBehavior...' .

Without that variable the probe reports **skipped**, and a skip is not a pass. Check for the skip line in the output every time.
