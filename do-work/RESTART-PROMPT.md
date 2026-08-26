```
do-work clarify
do-work run --fan-out 2

This command is sufficient; everything below it is context.

Two REQs are ready and their write sets are disjoint, so a fan-out of two is safe.
One REQ is held at pending-answers with a design question that clarify will surface;
answer it before that REQ is built, because the answer decides whether it is written
at all.
```

---

## Reference

### Where the work stands

Branch `claude/request-id-autocomplete-xnhkow`, head `a90e19a`, pushed, clean tree, no worktrees,
no builder branches, nothing claimed. PR [#169](https://github.com/knews2019/skill-do-work/pull/169)
is open and mergeable with no CI configured on this repo.

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
| REQ-381 | pending | none | Citation index + filter — UR-076's "find tickets that cite REQ-1679". Route C. |
| REQ-382 | pending | `depends_on: [REQ-381]` | Markdown-link ids, drawer and clipboard. F6 folded in. |
| REQ-383 | **pending-answers** | `depends_on: [REQ-382]` | Held. See the design question below. |

**REQ-380 and REQ-381 are safe to run concurrently.** REQ-380 writes only
`skills/do-work/actions/capture-reference.md`; REQ-381 writes `citations.go`, `generate.go`,
`web/board-filters.js`, `citations_test.go` and `generate_test.go`. No overlap.

**REQ-382 must not run beside REQ-381.** Both write `generate_test.go`. That is why it carries a
`depends_on` edge and not a sentence — `write_set` is display-only and gates nothing. This branch
learned that twice from the Codex reviewer, on REQ-381 and again on REQ-382 an hour later.

**Critical path is REQ-381.** Start there, not on REQ-380, if only one thing runs.

### The held question, and why it matters more than it looks

REQ-383 was captured to *harden* the hand-rolled Markdown fence scanner REQ-379 introduced. Before
building it, read its `## Open Questions` — the recommendation is to **delete the scanner instead**.

Every external finding on REQ-379 was the same shape: the hand-rolled scanner disagreeing with
CommonMark. Blockquoted fences, backtick info strings, list-item fences, multi-line code spans,
indented code blocks, link reference definitions — six symptoms, one cause. The Go side already runs
goldmark and already ships a build-time index for exactly this class (`repoFileMentions`). Computing
annotatable byte ranges there would delete `codeFenceRunFor`, `codeFenceRunCloses`,
`findMatchingBacktickRun`, `stripContainerPrefix` and the paragraph lookahead outright.

**REQ-381 needs a Go-side body scan anyway** for its citation index, so if the answer is "delete it",
fold that work into REQ-381 rather than running REQ-383 fourth. That changes the order, which is why
the question is gated ahead of the build rather than left for the builder.

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
