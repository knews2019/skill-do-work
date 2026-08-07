---
id: REQ-132
title: Kanban restart refuses a surviving port listener
status: completed
claimed_at: 2026-08-07T11:05:02Z
completed_at: 2026-08-07T11:09:38Z
route: B
created_at: 2026-08-07T08:45:11Z
user_request: UR-029
addendum_to: REQ-017
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-129, REQ-130, REQ-131]
batch: audited-safety-fixes
effort_estimate: normal
write_set: [justfile, actions/install.md, _dev/tests/contract-regressions.sh, actions/version.md, CHANGELOG.md]
---

# Addendum: Kanban Restart Refuses a Surviving Port Listener

## What

Harden every shipped `run-kanban` recipe so it waits specifically for the old listener to release the requested port, then refuses to start whenever any listener still owns that port. Keep the root recipe and installation template synchronized.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Both current recipe copies send SIGTERM and poll `kill -0` for twenty iterations of one hundred milliseconds. They do not verify during the wait that the old PID remains the listener on the requested port, do not inspect the port after waiting, and proceed to build and serve even when a listener survives. Existing contract coverage checks executable identification but not timeout, post-wait refusal, or exact recipe synchronization.

## Prior Implementation

REQ-017 introduced stale-listener replacement and browser opening in commit `0d5c1f5`. It deliberately kills only a verified queue-kanban executable, refuses to kill an unrelated listener, degrades gracefully when `lsof` is unavailable, and waits briefly for graceful shutdown. Preserve those safety boundaries and the existing `--open` behavior.

## Detailed Requirements

- Update every shipped `run-kanban` recipe, including the repository root recipe and the installer template.
- After terminating a verified queue-kanban listener, wait for at most 320 iterations of 100 milliseconds.
- During each iteration, use `lsof` to determine whether that specific old PID still listens on the requested port; stop waiting as soon as it no longer does.
- After the wait, query the requested port again for any listener, including a replacement process or a listener with a different PID.
- If any listener remains, print its PID and command, exit non-zero, and refuse to execute the build-and-serve line.
- Preserve the refusal to kill unrelated non-queue-kanban processes and preserve graceful behavior when `lsof` is unavailable.
- Keep numeric port validation, executable-path identity checks, browser opening, and other recipe behavior unchanged.
- Add regression coverage that prevents the root recipe and installation template from drifting from one hardened canonical shutdown behavior.

## Constraints

- Use the existing shell contract-regression framework, not a consumer-specific Jest test.
- Keep the recipe portable across supported shells and `just`'s separate-shell recipe-line behavior.
- Validate both recipe parsing and synchronization. Run the shell regressions, `bash -n`/ShellCheck where applicable, queue-kanban `go test ./...`, and `go vet ./...`.
- Preserve Claude-specific skill metadata such as `argument-hint` and validate update/install behavior remains synchronized.
- Do not stage, commit, or push without separate user authorization.

## Builder Guidance

Firm on listener identity, timeout, and refusal. A process can remain alive after releasing the port, and another process can acquire the port during shutdown, so process existence is not a substitute for listener-specific checks plus a final port query.

## Red-Green Proof

**RED prompt/case:** Exercise or stub a verified queue-kanban listener that remains bound after SIGTERM. The current recipe stops checking after twenty process-existence polls and reaches the build-and-serve line without a final port refusal. Also alter one recipe copy and observe that the current contract suite does not detect semantic drift.
**Why RED now:** Shutdown tracks PID existence rather than PID-on-port state, and the tests assert only selected tokens independently in both files.
**GREEN when:** The recipe performs up to 320 listener-specific polls, exits early when the old PID releases the port, performs a final any-listener query, names and refuses every surviving listener, never kills a foreign process, and the regression fails whenever the two shipped recipe copies diverge from the hardened canonical behavior.
**Validation:** User confirmed by requesting capture of the audited findings after reviewing the stale-wait analysis.

## Dependencies

No queued dependency. This addendum corrects completed REQ-017.

## Full Context

See `do-work/user-requests/UR-029/input.md` for the complete capture context.

---
*Source: UR-029 - "run do-work capture-request on these issues"*

---

## Triage

**Route: B** - Medium

**Reasoning:** The shell behavior is compact but safety-critical, duplicated in an installer recipe, and requires both executable behavior coverage and drift protection.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Both recipe copies contain the same shutdown line: executable-path identity protects foreign processes, but the bounded wait polls process existence for only twenty iterations and has no final port query. The contract suite already ratchets executable identity, so it can extract the shutdown line, execute it with shell-function seams for lsof, ps, kill, and sleep, and compare the two canonical copies exactly.

## Scope

**Files I will touch:**
- `justfile` (modify) - wait on listener ownership and refuse if the port remains occupied
- `actions/install.md` (modify) - keep the installer recipe byte-synchronized
- `_dev/tests/contract-regressions.sh` (modify) - add stuck-listener behavior and recipe-drift regressions
- `actions/version.md` (modify) - bump the integration version
- `CHANGELOG.md` (modify) - record the port-shutdown safety fix

**Files I will NOT touch:** The Go server, browser opener, board action, or archived antecedent.

**Acceptance criteria (restated from REQ):**
- [ ] The old PID is polled specifically as a listener on the requested port for at most 320 iterations.
- [ ] A final port query refuses startup whenever any listener remains.
- [ ] Refusal names the remaining PID and full command.
- [ ] Foreign-process protection, missing-lsof degradation, port validation, build, and browser opening remain unchanged.
- [ ] Root and installer shutdown recipes remain synchronized by regression coverage.

## Implementation Summary

**Files changed:**
- `justfile` (modified)
- `actions/install.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `actions/version.md` (modified)
- `CHANGELOG.md` (modified)

**What was done:** Replaced process-existence polling with a 320-iteration listener-specific wait, added a final all-listener port refusal naming PID and command, and added executable behavior plus exact recipe-drift regressions.

## Qualification

Passed - 5 files verified, 5 acceptance requirements traced, P-A-U confirmed.

## Testing

**Tests run:**
- /bin/bash _dev/tests/contract-regressions.sh - passed
- just --justfile justfile --list - passed

**Red-green validation:**
- RED: Both recipe copies lacked listener-specific and final-query contracts, retained kill -0, and the stuck-listener probe exited zero.
- GREEN: The contract suite passed with a 320-iteration listener wait, final PID/command refusal, foreign-process non-kill proof, and exact recipe-copy comparison.

**Existing tests updated:** _dev/tests/contract-regressions.sh executes the canonical shutdown line with lsof, ps, kill, and sleep seams and ratchets installer synchronization.

## Root Cause

The recipe waited for process termination rather than port-listener release, stopped after two seconds, and treated timeout as permission to proceed without checking the port again.

## Review

**Overall: 100%** | 2026-08-07

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition - this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass - a persistent listener blocks startup with PID and command, foreign listeners remain protected, both recipes match, and the justfile parses.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Extracting and executing the shipped shell line made the regression behavioral while an exact copy comparison prevents installer drift.

**What didn't:** Process existence was an indirect shutdown signal and timeout had no refusal semantics, so a live socket could survive both checks.

**Worth knowing:** Port safety needs two layers: poll the old PID as a listener, then query the port without assuming the old PID is still the occupant.
