# Builder brief — REQ-486, increment 1 of 2: the folds

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-486-collapsible-ur-progress-summaries`
**Branch:** `worktree-agent-REQ-486-collapsible-ur-progress-summaries`, based on `a2497c6`
**Route:** C. **TDD: yes.** **Impact: user-visible.** **Estimate P50: 85 active minutes for the whole REQ.**

## Read first

1. `do-work/runs/work-2026-09-05-231943/REQ-486-plan.md` in the MAIN checkout — the synthesized plan. Tasks T1 and T2,
   decisions D2, D3, D6, D9, D15, and the testing section are yours. It is the authority; this brief
   is a pointer.
2. `do-work/runs/work-2026-09-05-231943/REQ-486-exploration.md` — the evidence, with file:line anchors.
3. `do-work/working/REQ-486-collapsible-ur-progress-summaries.md` in the worktree — the request, its
   `## Plan` and its `## Scope`.
4. `_dev/primes/prime-kanban-board.md` and
   `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — binding for anything touching the
   board: versioning, parser lock-step, build outputs.
5. Crew rules in `skills/do-work/crew-members/`: general.md, shared-principles.md,
   coding-guardrails.md, communication-style.md, testing.md, ui-design.md, frontend.md.

## This increment: T1 and T2 only

**T1 — make By UR groups foldable by deleting the head's two-branch shape.** `renderUserRequestLens`
in `web/board-cards.js` builds two shapes of group head today. Delete the branch: always build the
`div.ur-group-row` wrapper, always move `data-detail-kind`/`data-detail-id` onto the sibling
`button.ur-group-detail`, always attach the same click listener. The only thing the folded flag still
decides is the initial state. This is a net line deletion.

**T2 — make the drawer's REQ-id list foldable and height-capped.** `openUserRequestDetail` in
`web/board-detail.js`. Add one sibling helper beside `appendMetaRow` that puts a real `<button>`
carrying the label inside the `<dt>` and toggles the `<dd>`'s `hidden` property. Only the "REQ ids"
row becomes foldable. Add `max-height: 40vh; overflow-y: auto` to the id list (D9) — the list starts
open, so a fold alone reproduces the reported problem on every open.

**Do NOT start T3, T4 or T5.** No Go change, no payload field, no rollup, no clock change, no CSS for
the summary strip, no doc edits, no release. A second builder takes those against your merged branch.

## The one assertion you must invert, not delete

`javascript_behavior_a_test.go:1652` asserts the By UR head does NOT carry `aria-expanded`. It becomes
red the moment T1 lands. **Invert it in place, in the same commit** (D3). A builder who resolves the
red by not setting `aria-expanded` has silently dropped a Detailed Requirement; a builder who deletes
it removes the only thing watching the seam between the two readings.

## Two selector repairs that are forced, not optional

`user_request_clipboard_browser_probe_test.go:142` and `generate_test.go:1037` select the group head
as the drawer trigger. Moving `data-detail-*` to the sibling button breaks them. They must land in the
same commit as the markup change, or the heavy browser lane breaks while the fast stage stays green.
Two of the three planning panels missed these entirely.

## Testing — three lanes, and a lane that skips reports success

Record each lane's own exit line, not the gate's.

```
# fast stage (Go payload only; both client lanes excluded by prefix here)
cd skills/do-work-board/tools/queue-kanban && \
  QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off \
  DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior go test -count=1 ./...

# heavy JavaScript lane (every semantic claim; slices the shipped source)
QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on \
  QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run '^TestJavaScriptBehavior' -v .

# heavy browser lane (pixels and trusted input only)
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=on \
  QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  go test -count=1 -run '^TestBrowserBehavior' -v .
```

Take a green baseline of both heavy lanes BEFORE your first edit, so a later red has a known-green
predecessor.

## Environment — read this before running anything

Fresh cloud container. Wrap every gate and every Go test:

```
env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL=/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gitconfig-gate \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  <command>
```

`NODE_OPTIONS` and the `GIT_CONFIG_*` triple are refused by the runtime fingerprint as opaque; the
global `commit.gpgsign` points at an empty key, so every `git commit` inside a test fixture fails.
Node is at `/opt/node22/bin/node`, Chromium at `/opt/pw-browsers/chromium`. Capture exit status from
`$?` directly — never pipe a gate to `tail`, or a red gate reads as green. This machine has 4 CPUs:
run the full canonical gate at most twice.

## Hand-back

Commit on your branch in the worktree, message ending
`Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>`. Write your report to the ABSOLUTE path
`/home/user/skill-do-work/do-work/runs/work-2026-09-05-231943/REQ-486-handback-folds.md` in the MAIN checkout. Do not
stage or commit it, do not push, do not touch any other worktree.
