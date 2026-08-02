# Board

Renders the `do-work/` queue as a Kanban board in your browser, plus a completion calendar and a testing track. Read-only toward the work pipeline — it never claims a REQ or edits a `status`. The one thing it writes is the testing record (see Testing view).

> **Needs the Go toolchain** — the only do-work capability that does. Without `go` on your `PATH` the board reports and stops; nothing else in the skill is affected.

## Modes

| Mode | What you get |
|------|--------------|
| **serve** (default) | Live board at `http://localhost:8090` (`--port N` to change). Reload the page to see new state — it re-reads the queue per request rather than pushing updates. Ctrl-C stops it. |
| **static** | Self-contained HTML under `build/queue-kanban-board/` (`--out DIR` to change) — opens from `file://`, no server, shareable. Throwaway; safe to delete. |
| **summary** | Column counts in the terminal. No browser. |

## Board view

Four columns: **Pending**, **Claimed**, **Needs input · Blocked**, **Recently done**. Pending splits into *Ready* and *Waiting on dependencies* — an unsatisfied `depends_on` doesn't move a REQ out of Pending, it just sorts it under Waiting. A REQ whose `status` isn't in the schema vocabulary parks under Needs input · Blocked with an `invalid` flag.

Two strips sit above the columns and stay visible in every view:

- **Notes** — your `do-work note` lines, verbatim and collapsible. Hidden when there are none.
- **Completion anomalies** — finished REQs whose completion instant can't be resolved. Bookkeeping to fix, not recent work, so the strip ignores the filters and the window.

The toolbar carries a text filter (id or title), domain and status selects, a **Recently done** window (24h / 48h / 7d), a **Lens** toggle (flat Columns vs. grouped **By UR**), and the **Board / Calendar / Testing** switch. Calendar shows completed work day by day.

## Badges

| Badge | Meaning |
|-------|---------|
| domain, UR id, `route` | frontmatter, as declared |
| `blocked by …` | the external condition named in `blocked_by` |
| `unblocks N` | how many other REQs this one releases when it lands |
| `overlaps …` | declared write sets could collide — see below |
| `anomaly`, `⚠ future stamp` | unresolvable completion instant, or a timestamp later than now |
| `testing …` | the card carries a testing record |

### Reading the `overlaps` badge

It names the other pending/claimed REQs whose declared `write_set` could touch the same files — an informational heads-up, never a block. Under the exclusive-session model `do-work run` builds one REQ at a time, so the badge schedules nothing; it just surfaces declared file contention.

It under-reports, so read it as a floor, not a guarantee:

- **No badge doesn't mean safe.** A REQ that never declared a `write_set` can't overlap anything as far as the board can tell.
- **Globs match one path segment at a time.** `*` never crosses `/`, and `**` is not recursive — it behaves like a single `*`.
- **A malformed pattern matches nothing** in that direction — except that identical text short-circuits, so two REQs declaring the *same* malformed pattern still badge each other.
- **A directory entry never badges a file inside it.** `actions/` and `actions/board.md` don't see each other.

## Card drawer

Clicking a card opens a docked panel: the REQ's frontmatter as labelled rows — status, dependencies, write set and its overlaps, route, timestamps, testing record, file location — above the rendered REQ body. Claim and block rows tick as live timers while the hold is open. REQ and UR ids are links, **Copy** grabs the source Markdown under a `# REQ-042: <title>` heading so a paste carries its own identity, and the divider drags to resize (double-click resets).

## Testing view

Tracks who tested which finished REQ, in four columns: **Ready to test → In testing → Returned with feedback → Tested**. Pick your name in the toolbar or add it — profiles are one bullet each in `do-work/testers.md`, hand-editable. Then use the card buttons: *Start testing*, *Mark tested*, *Return with feedback* (prompts for the note), *Clear*. All but Clear need a tester selected.

Each action writes the record into the REQ's own frontmatter, so `git log` on the REQ file answers "who tested this, and when" — and the main board shows a `testing` badge without you switching views. There's no locking: changes land in your working tree for you to review and commit. A static snapshot renders this view read-only.

## Usage

```
do-work board
do-work board static
do-work board summary
do-work kanban
```

`do-work install just-kanban` adds `just run-kanban` / `kanban-static` / `kanban-summary` recipes if you'd rather run the board without the agent, plus `just run-do-work-update` for the guarded project-local skill updater.
