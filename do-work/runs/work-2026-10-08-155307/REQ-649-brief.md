# Builder brief — REQ-649 (Board Earmarked badge tooltip names Needs input · Blocked as the home for operator-gated work)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-649-board-earmarked-tooltip-names-blocked-column
- Branch: worktree-agent-REQ-649-board-earmarked-tooltip-names-blocked-column, created from main HEAD dd54712a with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-649-board-earmarked-tooltip-names-blocked-column.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-5, Constraints, Builder Guidance, Red-Green Proof. Requirement 5 (release) belongs to the integrator; you PROPOSE the changelog entry and the lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-142/input.md.
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-08-155307/REQ-649-handback.md
- Route A, tdd: false, impact-user-visible, effort-mechanical, domain frontend. One string literal in JavaScript.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, frontend.md — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-kanban-board.md and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (read both before editing; versioning is folded into the skill and the release step owns VERSION/CHANGELOG, not you; never commit build outputs). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md read the top entries for 0.305.79 (REQ-643, this tooltip's own history) plus the prime's Traps.

## The change (decided at capture, do not reopen)
File: `skills/do-work-board/tools/queue-kanban/web/board-cards.js`, the `assignedBadge.title` string (lines ~234-239).
1. Extend the last sentence so the title ends: "… The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." Keep the badge text, class, placement, and every other sentence unchanged. Match the file's string-concatenation style for line wrapping.
2. Add no group-header hint (`makePendingGroup` stays a name and a count).
3. Touch nothing else: not `model.go`, `generate.go`, `board.css`, the data payload, or the three other badges that carry "Display only: the board never reorders, blocks, or hides on this." (lines ~269, ~299, ~321).

## Write boundary
Exactly `skills/do-work-board/tools/queue-kanban/web/board-cards.js`. Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `grep -c "belongs under Needs input" skills/do-work-board/tools/queue-kanban/web/board-cards.js` prints 1.
- `cd skills/do-work-board/tools/queue-kanban && go build ./... && go vet ./... && go test ./... -count=1` (whole package, ≈60 s; the package tests read the web files); `gofmt -l .` prints nothing.
- If `node` is available: `node --check skills/do-work-board/tools/queue-kanban/web/board-cards.js`.
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: dd54712a.
- File manifest: the file with (modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, go build/vet/test and gofmt results, and the file checked.
- Red-green evidence: the grep count before (0) and after (1), and the full title string after.
- `## Decisions` (D-01 onwards; likely none beyond line wrapping).
- `## Discovered Tasks` (each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for `lessons-do-kanban.md` in the file's top-entry shape (`- [family: <slug>] <version>: **bold one-liner.** body`, self-contained, no relative archive link) only if the change taught something; "none" is a valid answer. A proposed CHANGELOG entry (descriptive title saying what shipped, plain words). The integrator writes both.
- Integration seams (none expected). Test wall times.
