```
do-work run --fan-out 3

This command is sufficient; everything below it is context.
```

---

# Reference

## State at handoff

| Fact | Value |
|---|---|
| Version | 0.212.25 (`VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` agree) |
| `maintainer-verify.sh` | exit 0 |
| Working tree | clean; `git worktree list` shows only the main tree; no `worktree-agent-*` branches |
| `do-work/working/` | empty of REQs (`baseline.json` only, locally excluded) |
| In-flight REQs | none — `CHECKPOINT.md`'s In Progress section is empty |
| Queue | **16**: 12 `pending`, 4 `pending-answers`, 0 `blocked` |
| Branch / PR | `claude/restart-prompt-handoff-nprlgd` → PR #145, open, everything pushed |

## Environment — read this first

This checkout is a Linux container. A fresh one needs the same three repairs:

- **Go 1.26.1, ShellCheck 0.11.0 and `just` installed by hand.** The gate pins exact versions of the first two and fails closed without them. `apt-get install shellcheck` gives 0.9.0 — too old; fetch the release binary. Go 1.26.1 goes at `/usr/local/go`.
- **The gate now also runs `gofmt`** (added by REQ-260, 0.212.21), resolved from `$(go env GOROOT)/bin/gofmt` rather than PATH. A Go install whose GOROOT lacks `bin/gofmt` fails the gate with a named path.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB gitignored binary into the source tree. Build to scratch with `-o`.

## The queue, and how it schedules

**One file is the whole bottleneck.** `_dev/tests/prescribed-shell-scripts-behavior.sh` is written by **REQ-258, 263, 264, 268 and 271** — at most one per wave. That forces about four more waves whatever the builder count. **REQ-258 restructures that file wholesale**, so run it alone or first; anything else sharing the file will land its cases in a file REQ-258 has dissolved.

Other write-set families, all mutually disjoint:
- **qualify.sh** — REQ-263, REQ-264 (also suite writers; serialise against the above)
- **Board / Durations** — REQ-266, REQ-273, REQ-277, REQ-278
- **Docs, actions, citations** — REQ-262, REQ-269, REQ-270, REQ-272, REQ-274

A safe next wave: **REQ-258** (alone in the suite family), **REQ-266**, **REQ-262**.

## Four REQs need a human first

Run `do-work clarify` before or during the next run — all four are `pending-answers`:

- **REQ-274** — "the SessionStart hook exits nonzero" is **false** (`hooks/session-start.sh:59` runs the repairer under `|| true`). Real consequence: a FAILED line in every session banner, forever, no self-heal. The false mechanism is still stated where it carries REQ-255 D-04's rationale.
- **REQ-275** — the repairer repairs any `_at` field by suffix; the board's `detectFutureTimestampFields` checks six hand-kept names. Latent until the schema grows an `_at` field.
- **REQ-276** — `record-commit-hash.sh` guards its *writer* against an unterminated fence but not its *readers*, on the last check every REQ passes through.
- **REQ-278** — nothing bounds the Durations label face off Linux. Scoped to measuring, not geometry surgery.

## What the last session learned, in one paragraph

Six REQs shipped and **ten** new ones were created — every one from a review finding a real defect. The queue does not converge; stop when the findings stop being worth fixing, not when the list empties. Three separate corrections landed on claims that were *right in conclusion and false in reasoning*: two builders argued correct decisions from provably false premises, and both were caught only because the reviewer checked the argument rather than the verdict. Twice the limiting factor was **measurement, not the fix** — a parity fuzz that compares only whether a file was mutated cannot see a shape that refuses, and a fuzz's blind spots are exactly the axes it holds constant.

## Standing rules that bit last time

- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry one forward.
- Exit code zero is the only proof a check passed. Never pipe the gate through `tail`.
- **Check a handed-back integration seam before applying it.** A seam is the one part of a worktree REQ the builder cannot test, and the integrator is its only reader. One shipped a false safety claim last session and was caught by review, not by the integrator who applied it.
- **Add the P-A-U block at claim time if the REQ lacks one.** Review-generated REQs often have none, and `qualify.sh`'s box audit then passes vacuously with the gate silently half-off (REQ-264 exists for this).
- **Transcribe the hand-back's Discovered Tasks into the REQ before Step 8.** A worktree builder cannot write the REQ file, so Step 8 finds nothing and drops them silently (REQ-270 exists for this).
- For anything that changes what appears on screen, generate a board and look at it.

## Where the evidence lives

- `do-work/CHECKPOINT.md` — the full session record: six shipped REQs with hashes and scores, corrections made, decisions with reach.
- `do-work/runs/work-2026-08-18-{211613,230100}/` — two run directories: briefs, hand-backs with red-green evidence, orchestrator manifests.
- `do-work/archive/REQ-2{57,59,60,61,65,67}-*.md` — each archived REQ carries its P-A-U trail, decisions, independent review, and lessons.
- `do-work/calibration-log.tsv` — six rows appended. **Do not recalibrate from them**: the spans measure serial integration queuing, not work.
