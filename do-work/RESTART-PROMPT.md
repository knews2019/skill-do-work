```
do-work run --fan-out 3

This command is sufficient; everything below it is context.

You are picking up the do-work queue in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2.
Nothing is in flight. The tree is clean, `git worktree list` shows only the main tree, there are
no `worktree-agent-*` branches, `bash _dev/tests/maintainer-verify.sh` exits 0, and
`queue-kanban verify` reports no findings. Version is 0.212.7.

Nine REQs are queued and every ordering constraint is already encoded as `depends_on`, so the
command above schedules correctly on its own. It will start REQ-246, REQ-248 and REQ-249 — the
three chain-heads — and pick up REQ-247, REQ-252, REQ-253, REQ-250, REQ-251 and REQ-254 on
later waves.

Read do-work/CHECKPOINT.md before you build. Its Session Notes carry one finding that recurred
five times in the previous session and will apply to most of this queue: every REQ shipped a
mechanism that looked like it closed a class and closed only the instance, and in three of five
the remaining hole was exactly where the real data lives. Assume your first fix has that shape
and go looking for it before a reviewer does.

Start with REQ-248 if you build anything by hand. Panel B's leftmost bar sits in the axis gutter
on this repository's own board, and a one- or two-day board renders Panel B entirely off-canvas
at x=-3330.

Two standing rules the previous session learned the hard way:
  - Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never
    carry a timestamp forward from earlier in the session and never compute one.
  - For anything that changes what appears on screen, generate a board and look at it. A passing
    assertion is not evidence about two glyphs sharing a coordinate.

The disk is at 100% (about 260 MiB free of 123 GiB). `maintainer-verify.sh` failed mid-session
on that alone, with 36 `No space left on device` errors inside the installer probe. Free space
before a long run, and delete scratch when you finish with it.
```

---

# Reference

Everything below is for humans and debugging. The paste block above is sufficient to start.

## State at handoff

| Fact | Value |
|---|---|
| HEAD | `0a1ecf9` |
| Version | 0.212.7 (`VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` all agree) |
| `maintainer-verify.sh` | exit 0 |
| `queue-kanban verify` | OK: no findings |
| `git status --short --untracked-files=all` | empty |
| `git worktree list` | main tree only |
| `git branch --list 'worktree-agent-*'` | none |
| `do-work/working/` | `baseline.json` only — no claimed REQs |
| In-flight REQs | none |
| Foreign claims | none — `do-work/CHECKPOINT.md`'s In Progress section is empty |

## Worktree verdicts

`git worktree list` shows exactly one entry, the main tree. There is nothing to classify as
ACTIVE, FOREIGN or REMOVABLE, and no `git worktree remove` command is owed. The worktrees parent
directory `../skill-do-work2-worktrees/` was removed once empty.

All five of the previous session's worktrees and branches were removed with `git worktree remove`
and `git branch -d` — never `-D`. Every `-d` succeeded, which is the only mechanical proof those
merges landed.

## Shipped in the previous session

All five REQs merged, archived, and released. Every hash confirmed with
`record-commit-hash.sh --verify`. `maintainer-verify.sh` exited 0 at every commit boundary.

| Version | REQ | Merge commit | Merge range |
|---|---|---|---|
| 0.212.3 | REQ-241 reconcile the Durations label metrics | `90c74b7` | `2432f45..90c74b7` |
| 0.212.4 | REQ-243 resolve heading anchors on shipped pointers | `37d7729` | `2432f45..37d7729` |
| 0.212.5 | REQ-245 name fabricated stamps in future-stamp warnings | `23bad9d` | `2ad71eb..23bad9d` |
| 0.212.6 | REQ-242 move Panel B's slowest-day annotation | `48263dd` | `51beffc..48263dd` |
| 0.212.7 | REQ-244 cite the Timestamp rule at every stamp site | `f733365` | `37d7729..f733365` |

Ranges are cumulative where a remediation re-merged, per the remediation rule: first pre-merge
tip to latest merge.

All five were independently reviewed and **all five returned PASS-WITH-FINDINGS** (80, 88, 88,
88, 66 percent). None shipped on its first pass. Every one was remediated on its builder's own
branch and re-merged.

## The queue, and why it is ordered the way it is

Nine REQs. Three dependency chains of depth two, plus three leaves.

```
REQ-246 ──▶ REQ-247          mechanical timestamp repair (hook, then archive audit)
REQ-248 ──▶ REQ-252          Durations geometry, then measurement provenance
REQ-249 ──▶ REQ-253          citation-form sweep, then Timestamp rule paragraphs

REQ-250   REQ-251   REQ-254  leaves, no gates
```

**Ready set** (status `pending`, no unmet `depends_on`): REQ-246, REQ-248, REQ-249, REQ-250,
REQ-251, REQ-254.

**Auto-wave takes the first N in numeric id order**, so `--fan-out 3` starts REQ-246, REQ-248 and
REQ-249 — which are exactly the three chain-heads. That is a consequence of the numbering, not a
coincidence worth relying on blindly: if you change N, re-check what the first N actually are.

### Why 3 and not more

The three chain-heads are the critical path — each unblocks a second REQ, and the leaves unblock
nothing. Starting them first maximises what a second wave can pick up. Three also keeps the wave
inside domains that do not touch: core hook script, board JavaScript plus its test, and shipped
action markdown.

Raising N to 4 or more adds REQ-251, which edits `timestamp.go` and `verify_test.go` — the same
Go package as REQ-248's `generate_test.go`. Different files, and REQ-251 is comment and fixture
text only, so the risk is small. But this exact shape is what broke the previous session's build:
REQ-241 and REQ-242 each declared `durationsMeasuredAxisTitleAscentUnits` in different files of
one package, never touched adjacent lines, and git merged them cleanly into code that would not
compile. If you raise N, expect to resolve that class by hand and **run both sides' tests by
name** after resolving, not merely the compile.

### Must not run concurrently, and why

Both of these are encoded as `depends_on`, so the one-line resume cannot violate them.
`write_set` is display-only and auto-wave deliberately does not read it, so a collision expressed
only as overlapping `write_set` paths would not have gated anything.

- **REQ-252 must follow REQ-248.** They share
  `skills/do-work-board/tools/queue-kanban/generate_test.go`. The binding reason is semantic
  rather than textual: REQ-252 records the provenance of measurements taken against Panel B's
  geometry, and REQ-248 moves that geometry. Recording provenance for a layout about to change
  produces numbers that are wrong on arrival.
- **REQ-253 must follow REQ-249.** Both edit `skills/do-work/actions/work-reference.md`,
  `ui-review.md`, `memory.md` and `memory-reference.md` — a guaranteed conflict in at least four
  files. REQ-249 runs first so REQ-253's new citations conform to a settled convention instead of
  one about to be swept.
- **REQ-247 must follow REQ-246** — pre-existing gate, unchanged. The archive auditor builds on
  the repair path the hook script establishes.

### Safe to run concurrently

REQ-246 (core script plus SessionStart hook), REQ-248 (`web/board-durations.js`,
`generate_test.go`), REQ-249 (shipped action markdown across four skills), REQ-250
(`_dev/tests/shipped-package-reference-contract.sh`) and REQ-254
(`skills/do-work/tools/checks/qualify.sh`) touch disjoint files.

One caveat worth stating rather than discovering: **REQ-246 has no declared `write_set`** and may
add prose to an action file when it documents the hook wiring. REQ-249 sweeps all action
markdown. Any overlap there is an ordinary textual conflict that git will show you — not the
invisible kind — so it is a merge to resolve, not a reason to serialise.

### Critical path

REQ-248 is the highest-value item and the one to start on if you build anything by hand. It is a
live visual defect on the maintainer's own board: the leftmost Panel B bar spans x=37.1 to 49.1,
entirely left of `DURATIONS_MARGIN_LEFT` (54) and struck through by the "0" axis tick; a one-day
board puts the annotation at x=-3330 and the bar at x=-3342, both off-canvas, so Panel B renders
empty. Root cause is `xOfEpoch` mapping each day bucket's midnight while `timeStart` is the first
completion instant. The suggested fix is to floor `timeStart` to its UTC midnight and ceil
`timeEnd` to the following midnight before computing `timeSpan` — check it against Panels A and
C, which read the same domain.

### Nothing is held back

No REQ is at `pending-answers` and none is `blocked`. The three open questions from the previous
session were answered by the user through another session in commit `4b15b0e`, and both REQs
flipped to `pending`:

- **REQ-249** — literal `../../` citation paths become the rule, so a citation a reader pastes
  actually resolves and a future checker can verify it mechanically.
- **REQ-253** — `ui-review.md`'s report header date is UTC, matching every other stamp; and the
  `## HH:MM UTC` time-of-day headings are declared out of the Timestamp rule's scope and marked
  so a future sweep walks past them.

Those answers are what the two REQs now build. Do not re-ask them.

## Heads-up — things that will bite in the first ten minutes

- **The disk is at 100%**, roughly 260 MiB free of 123 GiB. `maintainer-verify.sh` failed
  mid-session with 36 `No space left on device` errors inside the installer probe, which needs
  room to build install and rollback snapshots. About 350 MB of scratch was reclaimed (including
  a 168 MB Chromium a reviewer downloaded and an 11 MB compiled `queue-kanban` binary a builder
  had left in the source tree). *Who should act:* the maintainer, by freeing real space. Any
  agent: delete `/tmp` scratch when finished, and remove Playwright browser downloads.
- **`queue-kanban verify` labels a live builder's worktree `merged-worktree-leftover [fixable]`
  and advises `do-work cleanup`.** It is true — the orchestrator merged the branch, so its tip is
  in HEAD — and acting on it mid-remediation deletes work in progress. *Who should act:* the
  orchestrator, by not running cleanup while a builder is still committing.
- **Browser tooling writes outside a builder's worktree.** The Playwright MCP drops a
  `.playwright-mcp/` directory into whatever it considers its working root, which can be the main
  tree, without the builder issuing a write. It is gitignored so nothing breaks, and it holds
  sibling agents' evidence — remove only your own files from it. *Who should act:* every builder
  brief should state this as an exception rather than leaving it to be discovered.
- **`qualify.sh` cannot distinguish a check's own success output from a leftover debug print.**
  It FAILed REQ-244's remediation on the success line that REQ-244's own review had just required
  it to add. The override is recorded in REQ-244's archived Review Remediation section with the
  diff quoted. REQ-254 fixes this. *Who should act:* whoever hits a `+print(` FAIL — read the
  diff, and if it is a contract's output, record the override rather than silencing it.
- **A builder left a compiled `queue-kanban` binary in the source tree** during the previous
  session. It is gitignored so git stays clean, but it inflates every install-probe copy. If the
  installer suite starts failing on space, check
  `skills/do-work-board/tools/queue-kanban/queue-kanban` first.

## Where the evidence lives

- `do-work/CHECKPOINT.md` — session notes, decisions with reach, and the context summary. Read
  this before building.
- `do-work/runs/work-2026-08-18-124358/` — the previous run directory: per-builder briefs,
  hand-backs with full red-green evidence and addenda, and the orchestrator manifest.
- `do-work/archive/REQ-24{1,2,3,4,5}-*.md` — the archived REQs, each carrying its transcribed
  P-A-U loop, scope, testing evidence, decisions, lessons, builder pushback and independent
  review.
- `do-work/calibration-log.tsv` — five rows appended: 35/28, 20/20, 5/48, 25/34, 55/32
  (estimate/actual active minutes). The Route C estimate overshot again. REQ-245's 5/48 is not an
  estimator failure — it was correctly estimated trivial and then had its scope widened twice by
  the orchestrator on the builder's own findings; exclude it from any recalibration or record why.
- `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-shell-commands.md` — updated with
  this session's lessons, including one new board convention: a measured face is per-browser, so
  record the build beside the number and take the larger where two disagree.
