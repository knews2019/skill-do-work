# Board

Renders the `do-work/` queue as a Kanban board in your browser, plus a queue activity calendar and a testing track. Read-only toward the work pipeline — it never claims a REQ or edits a `status`. The one thing it writes is the testing record (see Testing view).

> **Needs the Go toolchain.** The board and the core deterministic command platform are compiled Go programs. Without `go` on your `PATH`, the requested command reports the missing prerequisite and stops instead of falling back to prose mutation.

## Modes

| Mode | What you get |
|------|--------------|
| **serve** (default) | Live board at `http://localhost:8090` (`--port N` to change). Reload the page to see new state — it re-reads the queue per request rather than pushing updates. Ctrl-C stops it. |
| **static** | Self-contained HTML under `build/queue-kanban-board/` (`--out DIR` to change) — opens from `file://`, no server, shareable. Throwaway; safe to delete. |
| **summary** | Column counts in the terminal. No browser. |
| **open-work** | What's in flight, in the terminal: the open total (pending / claimed / needs-input-blocked), every claimed REQ with its title, and every needs-input/blocked REQ with the status parking it there. Nothing finished, no browser — the two-second check when opening a board tab is more than the question deserves. |

## Board view

Four columns: **Pending**, **Claimed**, **Needs input · Blocked**, **Recently done**. Pending splits into *Ready* and *Waiting*. A `pending` REQ with an unsatisfied `depends_on` stays in Pending and sorts under Waiting. A bare `blocked` REQ with an unmet dependency also displays there until every dependency is source-ready, while retaining its blocked status and badges; once its dependencies are source-ready, its external blocked condition is actionable and the card returns to Needs input · Blocked. Source-ready means terminally successful, or `claimed` while carrying a `commit:` — a REQ held for heavy lanes has landed its implementation, so its dependents are already selectable. A REQ whose `status` isn't in the schema vocabulary parks under Needs input · Blocked with an `invalid` flag.

Two strips sit above the columns and stay visible in every view:

- **Notes** — your `do-work-toolbox note` lines, verbatim and collapsible. Hidden when there are none.
- **Completion anomalies** — finished REQs whose completion bookkeeping is broken: an unresolvable completion instant, or an impossible one (`completed_at` earlier than `claimed_at`). Bookkeeping to fix, not recent work, so the strip ignores the filters and the window.

The toolbar carries a text filter (id or title), domain and status selects, a **Recently done** window (24h / 48h / 7d), a **Lens** toggle (three readings: flat **Columns**, grouped **By UR**, and **URs only**), and the **Board / Calendar / Testing** switch. In the By UR lens, **Active** shows URs with open work or a REQ inside the selected Recently done window; **All** also browses older resolved URs. Calendar carries every REQ, colour-keyed by status: work that has not started sits in a band at the top, claimed REQs sit on the day they were claimed, and completed, cancelled, and failed REQs sit on the day they resolved. Each day label counts what happened on it, and the summary line above doubles as the colour key.

### The By UR and URs only lenses

**By UR** groups every REQ under its user request and starts with each group's cards visible. **URs only** is the same reading with the cards folded away, for scanning many user requests at once; its groups start collapsed. In both, the group header is a real button that folds its own card grid — several groups can be open or closed at once — and **Details** beside it is a separate button that opens the drawer. Folding is view state only: it lives in the page, resets when the board re-renders or reloads, and never touches the queue.

The By UR header also carries a progress strip for the whole user request. It reads the UR's complete membership across queue, working and archive, so the text filter, the domain and status selects, the Recently done window and the Active/All scope change which cards you see and never change these five figures or their denominators:

- **Grouped REQs** — how many REQs belong to the user request.
- **Active** — time already spent. It sums the spans the board already accepted for finished members (measured from each REQ's FIRST recorded lifecycle stamp, the same rule behind the `took …` badge) plus the live elapsed time of every member claimed right now, and it ticks with the rest of the board while the page is open.
- **Remaining** — an approximate forecast, always marked with `~`. Each unfinished member contributes its saved `estimate.p50_active_minutes` when it has one, otherwise the Timeline's median for its effort class, and only while the Timeline has enough history to call that median confident. A claimed member spends its estimate down and stops at zero rather than going negative.
- **Successful** — `completed` plus `completed-with-issues`, with the count and the total.
- **Resolved** — successful plus `cancelled`, with the count and the total. A failed REQ counts toward neither percentage.

It says so when it does not know. `at least` means the figure is a floor because something is missing; `N excluded` counts finished members whose span the board refused (an assumed pause, or reversed stamps); `N unmeasured` counts members whose work has ended with no usable span at all, including every cancellation; `N unknown` counts unfinished members with neither a saved estimate nor a confident fallback; `⚠ clock skew` means a claim is stamped ahead of your clock; and a user request with no members reads `unavailable` rather than dividing by zero. Missing evidence is never counted as zero. The **URs only** reading deliberately shows no strip — it exists to see many user requests at once — so there the header keeps its plain REQ count.

## Badges

| Badge | Meaning |
|-------|---------|
| domain, UR id, `route` | frontmatter, as declared |
| `blocked by …` | the external condition named in `blocked_by` |
| `unblocks N` | how many other REQs this one releases when it lands |
| `overlaps …` | declared write sets could collide — see below |
| `anomaly`, `⚠ future stamp` | broken completion bookkeeping (unresolvable or reversed span), or a timestamp later than now |
| `took …` | wall-clock span from `claimed_at` to `completed_at`; informational, not a workflow state |
| `over 4h · assumed pause` | the span crossed the board's single-session ceiling, so it is assumed to include a pause and is excluded from duration medians; the REQ remains completed |
| `reversed stamps` | `completed_at` is earlier than `claimed_at`, so the card refuses to state a duration; use the `anomaly` badge to find the stamp to repair |
| `testing …` | the card carries a testing record |

### Reading the `overlaps` badge

It names the other pending/claimed REQs whose declared `write_set` could touch the same files — an informational heads-up, never a block. The badge schedules nothing, however many builders are running: when you hand several REQs out at once, the declared set is advisory input to *your* pick, and the merge refusing is what proves two builders didn't collide (`../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). It just surfaces declared file contention.

It under-reports, so read it as a floor, not a guarantee:

- **No badge doesn't mean safe.** A REQ that never declared a `write_set` can't overlap anything as far as the board can tell.
- **Globs match one path segment at a time.** `*` never crosses `/`, and `**` is not recursive — it behaves like a single `*`.
- **A malformed pattern matches nothing** in that direction — except that identical text short-circuits, so two REQs declaring the *same* malformed pattern still badge each other.
- **A directory entry never badges a file inside it.** `actions/` and `actions/board.md` don't see each other.

## Card drawer

Clicking a card opens a docked panel: the REQ's frontmatter as labelled rows — status, dependencies, write set and its overlaps, route, timestamps, testing record, file location — above the rendered REQ body. Claim and block rows tick as live timers while the hold is open. REQ and UR ids are links, **Copy** grabs the ticket file with its frontmatter fence byte-for-byte as it sits on disk — so the paste can be saved straight back as a valid REQ or UR — and adds two things to the body so the paste is not a wall of bare numbers: the first mention of each id it cites gains that ticket's title, and a **Referenced requests** list at the end names every id the body mentioned, with its title and status. Quoted text is left alone — a code fence, a code span or a blockquoted verbatim block keeps its exact words. The divider drags to resize (double-click resets). Three drawer rows can't come along, because they aren't in the file: **Tree** (which directory it's in), **Overlapping write sets** and **Unblocks** are all worked out while the board is built. Relative times don't carry either — you get the raw stamps they were computed from. If a board was generated by an older version and its source-text bundle is missing, Copy falls back to the rendered text under a `# REQ-042: <title>` heading rather than inventing a frontmatter block.

Opening a **UR** shows the same five progress figures as the By UR header, as labelled rows above the rest — the two surfaces render one rollup, so they cannot disagree. Below them, **REQ ids** lists every grouped REQ as a link. That row starts open, is height-capped so a user request with dozens of members cannot push `input.md` and the body out of the panel, and its label is a button that folds the list away entirely. Folding hides the ids and nothing else; the figures and the body stay put. Like the group folds, this is per-open view state and resets the next time you open the drawer.

## Testing view

Tracks who tested which finished REQ, in four columns: **Ready to test → In testing → Returned with feedback → Tested**. Pick your name in the toolbar or add it — profiles are one bullet each in `do-work/testers.md`, hand-editable. Then use the card buttons: *Start testing*, *Mark tested*, *Return with feedback* (prompts for the note), *Clear*. All but Clear need a tester selected.

Each action writes the record into the REQ's own frontmatter, so `git log` on the REQ file answers "who tested this, and when" — and the main board shows a `testing` badge without you switching views. There's no locking: changes land in your working tree for you to review and commit. A static snapshot renders this view read-only.

## Usage

```
do-work-board board
do-work-board static
do-work-board summary
do-work-board cli
do-work-board board
```

The full-suite installer publishes a managed flat recipe surface for board, core, knowledge, and toolbox commands. Run `just --list` for the live inventory. Board shortcuts remain `just run-kanban`, `run-kanban-cli`, `kanban-static`, and `kanban-summary`; `just do-work-update` is canonical, while `just run-do-work-update` remains a compatibility alias over the same transaction.
