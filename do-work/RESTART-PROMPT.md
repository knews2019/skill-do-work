```
do-work run --fan-out 3

This command is sufficient; everything below it is context.
```

---

# Reference

## State at handoff

| Fact | Value |
|---|---|
| Version | 0.212.19 (`VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` agree) |
| `maintainer-verify.sh` | exit 0 |
| Working tree | clean; `git worktree list` shows only the main tree; no `worktree-agent-*` branches |
| `do-work/working/` | empty of REQs (`baseline.json` only, locally excluded) |
| In-flight REQs | none — `CHECKPOINT.md`'s In Progress section is empty |
| Queue | 12 `pending`, 0 `pending-answers`, 0 `blocked` |
| Branch / PR | `claude/restart-prompt-handoff-nprlgd` → PR #145, open and mergeable, all review threads answered |

## Environment — read this first

This checkout is a Linux container, not the maintainer's machine. Three things had to be repaired before any work could be judged, and a fresh container will need them again:

- **Go 1.26.1, ShellCheck 0.11.0 and `just` were installed by hand.** The canonical gate pins exact versions of the first two and fails closed without them. `apt-get install shellcheck` gives 0.9.0 — too old; fetch the release binary. Go 1.26.1 was placed at `/usr/local/go` (the stock 1.24.7 moved aside).
- **One pre-existing test failure was Linux-only** and is fixed on this branch (0.212.8): the board's JavaScript behaviour probes now reach Node on stdin, because a probe embedding the assembled client exceeds Linux's 128 KiB per-argument limit. macOS has no such cap, which is why it had never been seen.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB binary into the source tree that is gitignored (so nothing warns you) and multiplies through the installer probe's copies. Build to scratch.

Disk here is fine (~30 GB free). The previous session's 100%-full warning does not apply to this container.

## The queue, and how it schedules

Twelve REQs, all `pending`, **every one already answered by the user in a `do-work clarify` session** — there are no open questions and nothing is gated on a human. Only one dependency edge exists (`REQ-257 depends_on REQ-255`, and REQ-255 shipped), so **all twelve are dependency-ready**; `--fan-out 3` takes the first three in numeric id order and recomputes each wave.

Write sets cluster into four families. Anything within a family should not run beside its siblings:

- **The timestamp repairer / archive auditor** — REQ-257, REQ-267, REQ-268 all touch `scripts/repair-req-timestamps.sh` or the file that sources it, plus the same behaviour suite. Run them one at a time. **REQ-267 carries the only known live defect worth prioritising:** a file whose frontmatter fence never closes can wedge the SessionStart hook into exiting nonzero every session, with no self-heal.
- **The qualify gate** — REQ-263 and REQ-264 both edit `tools/checks/qualify.sh` and its lock-ins. Serialise.
- **Board measurements** — REQ-265 and REQ-266 touch the Durations test constants and the JS renderer's comments. Adjacent but separable.
- **Docs and conventions** — REQ-258, REQ-259, REQ-260, REQ-261, REQ-262 are disjoint from everything above and from each other.

A safe first wave: **REQ-267, REQ-263, REQ-259** — one per family, disjoint files, and the highest-value item first.

## What the last session learned, in one paragraph

Nine of twelve REQs shipped a mechanism that looked like it closed a class and closed only the instance. The two that beat it did the same thing: they grepped or fuzzed the **primitive** before declaring the class closed, and both found the real hole where no instance list pointed. Assume your first fix has that shape. An instance list — even one assembled from two independent reviews — is a sample.

Three review channels ran this session and none was redundant: the pipeline's own adversarial reviews (what execution reveals), an external automated reviewer on the PR (what fresh eyes reveal), and the orchestrator re-reading its own bookkeeping (which caught a commit message claiming a transcription that had silently no-op'd). Expect all three to fire on work the others passed.

## Standing rules that bit last time

- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward, never compute one.
- For anything that changes what appears on screen, generate a board and look at it. A passing assertion is not evidence about two glyphs sharing a coordinate.
- Exit code zero is the only proof a check passed. Never pipe the gate through `tail` — the pipeline's status hides the failure.
- Builders were repeatedly killed mid-run by server-side 529s. Nothing was lost, because each committed in small increments; instruct builders to keep doing that, and resume rather than restart.
- `qualify.sh`'s box audit only fires when the REQ file carries a checked P-A-U block. Review-generated REQs created before this session lack one — add it at claim time, or the gate is silently half-off (that is REQ-264).

## Where the evidence lives

- `do-work/CHECKPOINT.md` — the full session record: twelve shipped REQs with hashes and review scores, decisions with reach, calibration caveats.
- `do-work/runs/work-2026-08-18-*/` — four run directories: per-builder briefs, hand-backs with red-green evidence, orchestrator manifests.
- `do-work/archive/REQ-24[6-9]-*.md`, `REQ-25[0-6]-*.md` — every archived REQ carries its P-A-U trail, decisions, independent review, and lessons.
- `do-work/calibration-log.tsv` — twelve rows appended. **Do not recalibrate from them**: the wall spans measure API failures and serial integration queuing, not the work.
