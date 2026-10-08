# Hand-back — REQ-649 (Board Earmarked badge tooltip names Needs input · Blocked)

- Branch: `worktree-agent-REQ-649-board-earmarked-tooltip-names-blocked-column`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-649-board-earmarked-tooltip-names-blocked-column`
- Base: `dd54712a`
- Commits: `d1b142aa` [REQ-649] Board Earmarked badge tooltip names Needs input · Blocked (single commit)

## File manifest

- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified): the `assignedBadge.title` concatenation gains one string literal after its last sentence: "A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here."

Nothing else changed. Write boundary held.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read general, coding-guardrails, shared-principles, communication-style, frontend crew members, both board primes, and the 0.305.79 (REQ-643) lesson. Approach: in `board-cards.js` the earmark tooltip's final literal ends `...hides on this.";`; change it to end `...hides on this. " +` and add one new literal line holding the operator sentence, matching the file's two-space-continuation concatenation style. No group-header hint, no other badge, no Go change. The REQ-643 lesson ("grep every restatement when the field's meaning changes") does not apply: the field's meaning is unchanged; this only adds where operator-gated work belongs.
- [x] **[APPLY]:** One literal appended exactly as planned; badge text, class, placement and the other sentences unchanged. The code comment above the badge is still accurate and was left alone.
- [x] **[UNIFY]:** `git diff --stat`: 1 file changed, 2 insertions(+), 1 deletion(-). `git diff --check` clean. `node --check web/board-cards.js` exit 0. In `skills/do-work-board/tools/queue-kanban`: `go build ./...` exit 0, `go vet ./...` exit 0, `go test ./... -count=1` exit 0 (ok, 75.2 s), `gofmt -l .` prints nothing. File checked: `web/board-cards.js` (diff reviewed line by line, no debug artifacts, no build outputs staged).

## Red-green evidence

- RED: `grep -c "belongs under Needs input" skills/do-work-board/tools/queue-kanban/web/board-cards.js` printed 0 before the edit.
- GREEN: same command prints 1 after the edit.
- Title after (with `request.assignedTo` interpolated):

```
Earmarked for <assignedTo> — an advisory claim marker, not a lock. Another session's default run skips and reports it; naming it explicitly overrides that and clears the field. The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here.
```

- Not exercised by me: the rendered hover in a browser. The integrator's `queue-kanban-browser` heavy lane covers the rendered page.

## Decisions

- D-01 (DECIDE & STATE): the trailing space separating the two sentences sits at the end of the previous literal (`"...hides on this. " +`), matching how the earlier literals in the same concatenation carry their trailing spaces.

## Discovered Tasks

None.

## Lessons read

- Satellite: `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, the 0.305.79 entry `[family: paired-predicate-drift]` (REQ-643, this tooltip's history), plus the prime's Traps in `prime-do-kanban.md` and the Conventions in `_dev/primes/prime-kanban-board.md`.

## Proposed lesson bullet

none. The change was a scoped text append and taught nothing new.

## Proposed CHANGELOG entry

**Board Earmarked badge says where operator-gated work belongs.** Hovering the `assigned` badge on a board card now ends with "A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." A user saw REQs that needed their input sitting under Pending → Earmarked and could not tell from the board that they were filed wrong. The tooltip now says it. Text only: no column, order, or data change.

## Integration seams

None. REQ-647 and REQ-648 touch other modules.

## Wall times

| Step | Exit | Seconds |
|---|---|---|
| go build ./... | 0 | 0 |
| go vet ./... | 0 | 1 |
| go test ./... -count=1 | 0 | 76 |
| gofmt -l . | empty | under 1 |
| node --check | 0 | under 1 |
