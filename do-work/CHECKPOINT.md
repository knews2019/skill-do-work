---
session_ended: 2026-08-04T21:38:00Z
last_completed: REQ-101
queue_state: 1 pending, 1 pending-answers, 0 blocked, 0 in-progress
reqs_processed_this_session: 7
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

## Completed This Session

- REQ-095: Two-clone acceptance run — poisoning repro and claim-conflict evidence (Route B, 93%) — v0.170.2, commit `0526e44`
- REQ-096: Execution-model re-grain — claim anywhere, one releaser (Route B, 92%) — v0.171.0, commit `7024c4a`
- REQ-097: `assigned_to` advisory field — schema, scan skip, board parse (Route B, 94%) — v0.172.0, commit `13328a8`
- REQ-098: Verify probes — assignment drift and half-closed URs (Route B, 95%) — v0.173.0, commit `47cd408`
- REQ-099: Automatic wave dispatch — the loop computes the ready set (Route B, 91%) — v0.174.0, commit `0cf9420`
- REQ-100: Live auto-wave acceptance run — wall-clock concurrency proven (Route B, 93%) — v0.174.1, commit `7ab69e3`
- REQ-101: Multi-checkout guide and ADR-018 session ownership (Route B, 94%) — v0.174.2, commit `e452989`

## Still Queued

- REQ-104: Label-less checkpoint entries — "locally modified" is not evidence of authorship (**pending, dependency-ready**) — generated from REQ-095's finding F-07; `maintenance: true`. **This is a known safety hole in what this batch shipped** and was deliberately left unbuilt: it fell outside the scope named for this run, and its fix carries a real either/or for the user to settle — narrow the heuristic to "modified *and* no merge in progress", or drop it so a label-less entry is always report-only, losing auto-recovery for pre-0.170.0 crashes. The REQ recommends dropping it.
- REQ-103: Checkpoint frontmatter writer identity (pending-answers — waits for `do-work clarify`)

## Context Summary

**UR-018 is functionally complete and stays open.** All eight originally-captured members (REQ-094 through REQ-101) are archived; the UR remains in `do-work/user-requests/` because REQ-103 and REQ-104 are unresolved, so consolidation into `do-work/archive/UR-018/` waits. The nine archived members sit in `do-work/archive/` root awaiting that consolidation.

**What the batch shipped, one line each.** 0.170.2 proved the checkpoint writer label works and found two defects. 0.171.0 re-grained ownership to claim-anywhere/one-releaser and renamed the Execution Model section. 0.172.0 added the advisory `assigned_to` field end to end. 0.173.0 added two `verify` probes. 0.174.0 added `do-work run --fan-out` auto-wave. 0.174.1 measured 4.109 seconds of real builder overlap. 0.174.2 wrote the guide and ADR-018.

**Two contract renames the next session must not re-derive.** The section is now `## Execution Model — Claim Anywhere, One Releaser` (cited by that name from four files), and the counted invariant is now `one releaser per queue` — stated exactly once in `actions/` and asserted by the suite, with a second assertion that fails if the retired `one queue owner per checkout` reappears anywhere in `actions/`, `docs/`, or `SKILL.md`.

**The suite is red for an environment reason, not a code one.** Eight FAIL lines, all from `_dev/tests/update-script-behavior.sh`, whose failure injection is `chmod 500` — which root ignores, so the updater succeeds where the probe expects failure. Every REQ this session was gated on "still exactly those eight, name-for-name". Do not read a green suite as achievable in this container, and do not read the eight as a regression.

**Re-read rather than trust carried assumptions.** Three things in this batch were settled by *observation* and contradict what reasoning predicted: a cross-checkout double claim is a content conflict and never a rename conflict; `do-work/CHECKPOINT.md` conflicts on every concurrent claim even when the two REQs overlap in nothing; and byte-identical claim writes leave the `writer:` label as the only detector. All three are recorded in `do-work/archive/REQ-095-two-clone-acceptance-run.md`'s `## Testing` section. Read it before writing any new prose about cross-checkout behavior.

## Session Notes

- **Rebuild `tools/queue-kanban/queue-kanban` before trusting any of its output.** A stale binary cost two false negatives in REQ-097 (the payload looked like it was missing the new field). Also: `generate --out` writes the payload to **`board-data.js`**, never `index.html` — asserting against the wrong file cost a third.
- **Wrap tool rebuilds in a subshell:** `(cd tools/queue-kanban && go build -o queue-kanban .)`. A bare `cd` moves the outer shell and the next repo-relative invocation fails in a way that reads like a build error.
- **Write multi-paragraph prose edits to a script file and run it by path.** Inlining them in a heredoc broke twice, on `→` and on quoting.
- **Adding a `do-work run` flag is five edits**, not one: `actions/work.md`'s Input list, its unrecognized-argument strip list, and its usage string, plus `SKILL.md`'s routing row and dispatch table. The strip list is the one that bites — omit it and the guard rejects the flag the same file just documented.
- **A new ADR is five file changes**: the record, the topic index twice (frontmatter `related:` and body bullet), the master index row, and `decisions/log.md`. `decisions/` is `export-ignore`d, so it may cite the maintainer doc; `docs/` may not.
- **The pipeline ran with in-session phases rather than spawned subagents** — the floor path `actions/work.md` explicitly permits it, and this harness withheld subagent dispatch. Every builder and review phase is labelled as such in the REQ records.
