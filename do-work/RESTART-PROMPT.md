```
do-work run --fan-out 2

This command is sufficient; everything below it is context.

Two REQs are ready with disjoint write sets, so a fan-out of two is safe.
Nothing is held awaiting an answer.
```

---

## Reference

### Where the work stands

Branch `claude/request-id-autocomplete-xnhkow`, restarted from `main` after
PR [#169](https://github.com/knews2019/skill-do-work/pull/169) **merged**. Clean tree, no worktrees,
no builder branches, nothing claimed. Follow-up work rides on
PR [#173](https://github.com/knews2019/skill-do-work/pull/173); no CI is configured on this repo.

Two REQs shipped this session, both archived under `do-work/archive/UR-075/` with their commit
hashes recorded:

| REQ | What | Acceptance |
|---|---|---|
| REQ-378 | Ticket ids in the drawer carry their titles; dead ids flagged; drawer glossary | Pass |
| REQ-379 | The clipboard payload carries titles and a referenced-requests glossary | Pass, 91% |

Released as **0.239.0** and **0.240.0**. `main` was merged at `a290cd6`, which reconciled a second
version collision — main took 0.238.0 for its architecture-report feature, so this branch renumbered.

### The queue, and why it is ordered this way

| REQ | Status | Gate | Note |
|---|---|---|---|
| REQ-380 | pending | none | Cross-Reference Convention in `capture-reference.md`. `tdd: false`, one file, closes UR-074. |
| REQ-383 | pending | none | Resolve ticket mentions in Go; delete the client Markdown scanner. **The foundation — start here.** |
| REQ-381 | pending | `depends_on: [REQ-383]` | Citation index + filter. Consumes REQ-383's AST walk. |
| REQ-382 | pending | `depends_on: [REQ-381]` | Markdown-link ids, drawer and clipboard. F6 folded in. |

**REQ-380 and REQ-383 are safe to run concurrently.** REQ-380 writes only
`skills/do-work/actions/capture-reference.md`; REQ-383 writes `citations.go`, `citations_test.go`,
`generate.go`, `web/board-clipboard.js`, `generate_test.go` and the board prime. No overlap.

**REQ-381 must not run beside REQ-383** — it consumes REQ-383's AST walk, and building first would
mean writing a mention scanner that REQ-383 then deletes. **REQ-382 must not run beside REQ-381**;
both write `generate_test.go`. Both edges are `depends_on` fields, not sentences — `write_set` gates
nothing, a lesson this branch learned twice from the reviewer.

**Critical path is REQ-383 → REQ-381 → REQ-382.** Start at REQ-383, not on REQ-380, if only one
thing runs.

### Why REQ-383 comes first, and what it replaces

REQ-383 was originally captured to *harden* the hand-rolled Markdown fence scanner REQ-379 shipped.
It has been rewritten to **delete** it, on the user's direction after the design was put to them with
probe output.

Every external finding against that scanner was one shape — the scanner disagreeing with CommonMark:
blockquoted fences, backtick info strings, list-item fences, code spans crossing a newline, indented
code blocks, link reference definitions. Six symptoms, one cause. The parser was always available and
simply sat on the wrong side: `render.go:25` already builds a goldmark renderer and already parses
every one of these bodies to make the drawer's HTML.

The mechanism is probe-verified, not assumed. One AST walk returned exact byte ranges for every
failing case, and `` ```lang`invalid `` produced no node at all because goldmark treats it as prose,
exactly as CommonMark says. Two API constraints came out of that probe and are recorded in the REQ:
`Lines()` **panics** on inline nodes, so a code span's extent comes from its child `Text` segments;
and offsets are **body-relative**, so every one shifts by the `bodyStartOffset` `splitFrontmatter`
already computes.

The acceptance signal is unusual and worth knowing: **REQ-379's existing clipboard assertions must
pass unmodified.** An assertion that needs rewriting means a behaviour changed that should not have.

### Heads-up list — things that will bite in the first ten minutes

- **`_dev/tests/maintainer-verify.sh` exits 1 before running any check** unless ShellCheck 0.11.0,
  `just`, and Go 1.26.1 are on `PATH`. The container ships ShellCheck 0.9.0, no `just`, and Go
  1.24.7. All three are fetchable; this session put them in its scratchpad. Prefix every invocation
  with `GOTOOLCHAIN=go1.26.1`, or the Go gate alone fails first. **Whoever acts:** the next session,
  before trusting any "the gate cannot run here" claim — that is a property of the container, not the
  repository, and a builder already filed it as a Discovered Task once in error.
- **`TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` fails on this container's
  Chromium** (Playwright chromium-1194) and is pre-existing — confirmed by stashing, in a file
  neither this branch nor the main merge touches. **REQ-375, which arrived from main, already owns
  it.** Do not attribute it to new work.
- **A real browser is at `/opt/pw-browsers/chromium`.** Set `QUEUE_KANBAN_BROWSER` to it or the whole
  browser lane silently skips, and REQ-381's acceptance depends on that lane.
- **Three vacuous mutations shipped in this batch** — a mutation whose anchor lands on a line no
  fixture exercises passes, and reports coverage that does not exist. Two were written by this
  session. REQ-383's requirements name this explicitly. **Whoever acts:** any builder adding an
  assertion — run the mutation *and* confirm the fixture reaches the mutated line.
- **Codex reviews this PR on `@codex review`** and has found a real defect on every pass so far
  (five across two REQs). Budget for a verify-fix-pin-push cycle per REQ rather than treating a green
  suite as done.

### Calibration, for whoever plans the next block

The estimator's p50 is *active agent minutes* and both completed REQs overran it: REQ-378 estimated
35, took 222 wall (6.3×); REQ-379 estimated 35, took 119 (3.4×). The gap is review round-trips, not
building. Remaining critical path reads as 75 estimator-minutes and realistically 4–6 hours at this
session's observed rate.

### Worktrees

None. The survey found one checkout, `/home/user/skill-do-work`, clean, with no `worktree-agent-*`
branches merged or unmerged. Nothing to remove and no foreign claim to preserve.
