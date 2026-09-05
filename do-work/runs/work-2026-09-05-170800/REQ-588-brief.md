# Builder brief — REQ-588

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1708/worktree-agent-REQ-588-one-warning-line`
- **Your branch (already checked out there):** `worktree-agent-REQ-588-one-warning-line`
- **Route:** A
- **Base commit:** 93f856ee (main)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work-board/tools/queue-kanban/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ (listed under "Write set" below). If you need one, stop and report it in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it. Run only the focused tests named below.
- Do not build or serve the board on port 8090: a live board owned by the user is running there. Use another port if you serve one at all.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/frontend.md`

Also read the REQ's `prime_files` entry `_dev/primes/prime-kanban-board.md`, the shipped prime it points at (`skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`), and in the two lesson satellites (`skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, `_dev/primes/lessons-kanban-board.md`) the bullets that touch the web assets, the Node behaviour lane and `verify.go`. The REQ's `## Required Lessons — Dropped for Budget` names both as over budget; that is a record, not a prohibition.

## The REQ

Read it in full in the main tree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-588-make-each-verify-finding-row-read-as-one-warning-line.md` (read only). Its `## Prior Implementation` says what REQ-579 built; its `## Requirements` D1–D5 are the contract; its `## Open Questions` is answered: **M1** ships.

## What to build (M1, from the mock-up report)

The approved mock-up is `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/mockups/m1-remedy-under-detail.html` (open `index.html` beside it for the reasoning; `mockups/shared.css` is the shipped rules copied verbatim, the page's own `<style>` block is the diff). Read that `<style>` block; it is the CSS you are shipping, adapted into `web/board.css` where the `.board-findings-*` rules from REQ-579 live (around lines 607–692):

1. **D2, chip | text grid.** `.board-findings-row` becomes `display: grid; grid-template-columns: max-content minmax(0, 1fr); column-gap: 10px; align-items: baseline; font-size: 0.8rem`. Wrapped text stays in the text column.
2. **D1, remedy on its own line.** `.board-findings-text { display: block }`, `.board-findings-remedy { display: block; margin: 1px 0 0 }` (drop its `margin-left`). In `web/board-cards.js` `makeFindingRow`, the remedy span gets `finding.remedy` with **no** `"→ "` prefix, and update the comment that says the arrow matches the terminal report. The "cleanup can fix" tag stays inline after the detail, inside the detail line.
3. **D3, one type scale.** `.board-findings-subject`: `font-family: var(--font-mono); font-size: 0.8rem; font-weight: 600; letter-spacing: 0; color: var(--ink-base)` (keep its margins and `overflow-wrap`). `.board-findings-chip`: `font-size: 0.66rem; padding: 1px 7px; line-height: 1.6`. Update the CSS comments so they describe the new rules, not the old flex row.
4. **D4, the producer stops repeating the subject.** In `verify.go`, every finding that sets `Subject` must not start its `Detail` with that same name. Today the worktree probes write `"<name> exists — <path>"`; the REQ-level probes write `"<REQ-id> has terminal status …"`, and so on. Rewrite each such detail so it reads correctly under a heading that already says the subject (the mock-up uses `exists at <path>` and `has terminal status "completed" but still sits in do-work/working/`). Grep `Subject:` in `verify.go` to find every site; keep the terminal `verify` report readable too (it prints detail without the heading — check how the terminal report renders a finding and, if it does not print the subject, print it there, or keep the subject in the terminal line only; decide and log it as a D-XX decision). Update `verify_test.go` expectations and any generate/snapshot fixture that embeds those strings.
5. **D5.** Everything else from REQ-579 stands: one list, weight only from `fixable`/"not checked", grouping by subject, `.board-findings-row-detached`, the two pass-through hosts, hide rules. No card classes.

## Write set (your write boundary)

- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go`
- `skills/do-work-board/tools/queue-kanban/verify.go`
- `skills/do-work-board/tools/queue-kanban/verify_test.go`

If D4 forces a fixture or another test file (for example a `generate_test.go` expectation carrying a detail string), that is the REQ's own requirement demanding the file class: flag it in the hand-back and proceed.

Board version and parser lock-step follow `_dev/primes/prime-kanban-board.md`; do not bump `VERSION` (finalization owns it).

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` in the queue-kanban module), verify no debug artifacts (`console.log`, `debugger`, stray `TODO`) in added lines, and list each file you checked and what you checked.

## Focused tests and proof

Every test-file invocation must finish in under 30 seconds. From the repo root of your worktree:
- Node lane: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...`
- Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`

`tdd: false`, but the REQ's `## Red-Green Proof` names the assertion for M1 and the Node lane can run it: extend `TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList` (javascript_behavior_c_test.go around line 2942) or add a sibling so it asserts that a finding row's remedy is an element whose text does not start with the arrow, sitting after the detail inside the text span; write it first, watch it fail on the current renderer, then make it pass. For D4, a `verify_test.go` assertion that a subject-bearing finding's detail does not begin with its subject, RED first.

The layout itself (grid, remedy on its own line, one scale) is a render fact: generate the board into a scratch directory (`go run . generate --out <dir>` from the module, or the shipped recipe the prime names), serve it on a port other than 8090, screenshot the strip with headless Chrome at 1600 px in light and dark (`"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless=new --disable-gpu --hide-scrollbars --window-size=1600,440 --screenshot=<png> <url>`), look at the images, and record what you saw and where the PNGs are in the hand-back. The live queue has at least one finding right now (a leftover worktree), so the real board shows the strip.

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-170800/REQ-588-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it. Sections, each under its own `##` heading: `Branch` (branch, head commit, base, worktree, tree clean), `File manifest` (verb, path, what changed, for every file), `P-A-U` (with `### [PLAN]`, `### [APPLY]`, `### [UNIFY]`), `Test evidence` (exact commands, exit codes, RED then GREEN output excerpts, render screenshots and what they show), `Lesson evidence` (which satellite bullets you read, which entries were missing), `Decisions` (D-XX entries, numbered from D-01, with reasoning; add Value/Risk lines on any you would escalate), `Discovered Tasks` (out-of-scope finds, or "none"), `Integration seams` (shared-file lines the orchestrator must apply, or "none").

Commit on your branch with a message starting `[REQ-588] `, then write the hand-back and stop. Do not merge, do not touch main.
