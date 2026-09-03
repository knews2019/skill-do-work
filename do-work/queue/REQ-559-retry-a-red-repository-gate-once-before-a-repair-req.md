---
id: REQ-559
title: '[impact-rule-change] Retry a red repository gate once before deferring or minting a repair REQ'
status: pending
priority: now
created_at: 2026-09-03T20:05:46Z
user_request: UR-106
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-548, REQ-531, REQ-560]
batch: lifecycle-overhead
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/checks/preflight.sh
---

# Retry a Red Repository Gate Once Before Deferring or Minting a Repair REQ

## What

When the repository gate exits non-zero, at the baseline (Step 5 pre-flight) or after integration (Step 6.5), rerun the exact same argv once, immediately, before any classification. A green rerun is recorded as green and the run continues; the retry is written to the run's progress output as one line. A red rerun enters the existing path unchanged: fingerprint, diagnostic worktree, `defer-gate`, repair REQ. The retry happens in the program that launches the gate where one exists (`preflight.sh` for the baseline), and as one rule in `work-reference.md` cited from Step 6.5 for the direct post-merge run.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"that is the lifecycle around it that we need to improve a lot": REQ-548 was an already-green repair. The gate passed when rerun, and the only thing the 28-minute detour proved was that the first run had flaked.

## Context

- 2026-09-03: REQ-531 claimed 21:19; its baseline gate failed with `update-script-behavior.sh: printf write error: Broken pipe` (REQ-548's recorded diagnostic evidence); deferred 21:37; REQ-548 claimed 21:37, found the gate green, archived 21:47 as an already-green no-op; REQ-531 re-claimed 21:47. Two full gate runs and one REQ lifecycle to learn that nothing was broken.
- The already-green repair path (0.266.8) made that no-op cheaper, one gate run instead of three, but it still costs a claim, a checkpoint entry, a finalization and an archive. A retry before classification removes the whole detour for every transient failure and changes nothing for a real one: a red retry carries the same fingerprint into the same defer path.
- Exactly one retry. A second failure is a real failure by definition here; bounding it keeps a genuinely red gate from doubling its cost.
- The retry counts as a gate run for the one-gate-per-machine budget rule in `_dev/tests/`; it is the same argv, run alone.

## Detailed Requirements

- Baseline: `preflight.sh` reruns the gate command once when the first run exits non-zero, records only the second result in `baseline.json`, and prints one WARN line naming the retry and both exit codes.
- Post-merge (Step 6.5 item 4): one rule in `work-reference.md`, cited from the step: on a non-zero direct exit, run the exact argv once more, directly and unpiped; zero exit records green and continues; non-zero continues with the existing diagnostic and defer procedure using the second run's output as the fingerprint source.
- The run's progress output shows the retry as one line; the REQ's Testing section records both exit codes when a retry happened.
- No new status, no new flag, no new REQ type. `defer-gate`, `repository_gate_repair`, and the already-green path are untouched; they simply fire less often.
- Delete any sentence that would now contradict the rule (for example, wording that treats the first non-zero exit as final).

## Constraints

- Mechanics in the script, judgment in prose; no new prose walking a shell sequence (CLAUDE.md, prime-shell-commands.md).
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.
- The repository gate itself, `_dev/tests/maintainer-verify.sh`, is not edited by this REQ.

## Red-Green Proof
**RED prompt/case:** Run a REQ through `do-work run` while the baseline gate fails once transiently (the recorded shape: a broken pipe from a probe's `printf` while its reader has exited) and passes on the next run.
**Why RED now:** the first non-zero exit is final: the REQ is deferred, a `repository_gate_repair` REQ is minted and claimed, and its already-green no-op completion is the only thing that lets the parent resume. REQ-548 on 2026-09-03 is the record: 28 minutes, two gate runs, one archived REQ, zero code changed.
**GREEN when:** the same transient failure produces one immediate rerun; the green rerun is recorded, the REQ continues without deferral, no repair REQ exists, and the progress output carries one retry line. A gate that fails twice with the same fingerprint still defers exactly as today.
**Validation:** Inferred during capture; the maintainer approved the capture ("do 1, 2 and 3").

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a pipeline step contract and its downstream readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget and `slugged: partial`. Matched because this REQ changes a shipped check script.

## Full Context
See `do-work/user-requests/UR-106/input.md` for complete verbatim input.
