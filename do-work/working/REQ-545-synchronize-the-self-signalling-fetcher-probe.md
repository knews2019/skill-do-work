---
id: REQ-545
title: 'Synchronize the self-signalling upstream-fetcher probe'
status: claimed
priority: now
created_at: 2026-09-03T15:25:00Z
user_request: UR-104
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-525]
write_set:
  - _dev/tests/update-script-behavior.sh
claimed_at: 2026-09-04T13:17:16Z
---

# Synchronize the Self-Signalling Upstream-Fetcher Probe

## What

`_dev/tests/update-script-behavior.sh`'s upstream-fetcher signal probe races its own signal delivery. It replaces `bash` with a function that signals the current shell and then returns 1; when the signal lands after the `kill` but before the shell acts on it, the function's `return 1` wins and the probe observes exit 1 instead of the conventional 128+signo.

## Instances

- `_dev/tests/update-script-behavior.sh:661-679` — the `for signal_case in HUP:129 INT:130 TERM:143` loop. Observed failing once as `upstream fetcher: HUP exits with its conventional status — expected exit 129, got 1`, which took down the whole gate lane with `update-script or suite installer behavior probes failed`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- Observed in a canonical gate run on the merged tree at `aeff306`. **Attribution was established rather than assumed:** the same probe passes on a clean `origin/main` worktree, and passes three consecutive times on the merged tree. Four passes against one failure, so it is a flake in the probe and not a regression in the code under test.
- Same family as REQ-525: a signal test whose own delivery is unsynchronized. That one was mis-captured as a test race and turned out to be a product defect; this one is the genuine article — the code under test is not involved, because the probe never reaches it.

## Detailed Requirements

- The probe must observe the shell's signal disposition, not a race between the signal and the overriding function's return value.
- Keep asserting the conventional 128+signo status for HUP, INT and TERM, and keep asserting that no archive is published and no success is reported.
- Do not lengthen a sleep, do not retry, and do not widen the accepted status set to include 1.
- The fix must hold under gate load, where the failure was observed; a pass in isolation proves nothing here, exactly as REQ-525 established.

## Constraints

- Test-side only. The upstream fetcher's own signal handling is not implicated: `got 1` is the override function's return value, so the fetcher was never reached.

## Dependencies

No request prerequisite.

## Red-Green Proof

**RED prompt/case:** run the probe repeatedly under parallel gate load; the HUP case intermittently reports `expected exit 129, got 1`.
**Why RED now:** the overriding `bash` function returns 1 after signalling the current shell, so the return value can beat the signal.
**GREEN when:** the three signal cases report their conventional statuses across repeated runs under load, with the non-publication and no-success assertions unchanged.

---
*Source: canonical gate run on the merged tree; attribution confirmed against clean origin/main.*
