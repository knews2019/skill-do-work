---
id: REQ-525
title: 'Block the install confirmation before signalling it'
status: claimed
created_at: 2026-09-03T00:35:00Z
user_request: UR-085
domain: testing
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-457, REQ-524]
claimed_at: 2026-09-03T11:44:30Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-03T11:45:48Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 4 acceptance criteria
    - async lifecycle behavior
---

# Block the Install Confirmation Before Signalling It

## What

Make `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation` wait until the installer is actually blocked on its confirmation prompt before delivering the signal. Today it signals on a timer, so under parallel load the installer can finish rendering its managed-install diff and exit before the signal is handled, and the assertion reads `exit = <nil>, want 130`.

## Instances

- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/install-suite/INT` and `/install-suite/HUP` (`suite_commands_test.go:292`). Observed failing twice under `go test ./...` parallel load during the REQ-457 run, and passing 5/5 with `-count=5` in isolation immediately afterwards. It did not fail in the pre-REQ-457 baseline, and REQ-457's diff touches no file in this package.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Discovered by:** REQ-457's Step 6.5 gate attribution, as the one gate failure that was neither in the recorded red baseline nor in a package this REQ's diff touches.
- **Evidence:** the captured stderr shows the installer still emitting `--- managed destination: .claude/skills/do-work ---` diff output when the signal arrives, i.e. it had not yet reached the blocked confirmation read the test name asserts against.

## Detailed Requirements

- Deliver the signal only after the installer has demonstrably reached its blocked confirmation read, rather than after a fixed delay or an unsynchronized amount of output.
- Keep asserting exit 130 and the existing managed-path non-effects; this is a synchronization fix, not a weaker assertion.
- Do not lengthen a sleep to paper over the race, and do not mark the test as skipped or flaky-tolerant.
- The fix must hold under `go test ./...` parallel load, not only under `-run` in isolation.

## Constraints

- Test-side only unless the product genuinely cannot signal readiness; if a product change is required, say so rather than widening silently.

## Dependencies

No request prerequisite.

## Red-Green Proof

**RED prompt/case:** `go test ./... ` for the module under parallel load, repeated; the subtest fails intermittently with `signal interrupt exit = <nil>, want 130`.
**Why RED now:** the test races the installer's own diff rendering instead of synchronizing on the confirmation prompt.
**GREEN when:** the subtest passes across repeated full-module parallel runs with the exit-130 and non-effect assertions unchanged.

---
*Source: REQ-457 gate attribution, captured during the work run.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Captured as mechanical on the assumption that the test signalled on a timer. It does not — it polls for the prompt text — so the actual race is still unlocated and needs discovery before any edit.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Two things the capture assumed that are wrong, established before dispatch:

1. **The test does not signal on a fixed delay.** `runBuiltCLIAtBlockedConfirmation` (`suite_commands_test.go:322`) starts the built CLI, then calls `waitForPromptInFile` (`:374`), which polls the stderr file every 10 ms for up to 10 s and returns as soon as the prompt text appears. Only then does it signal. So "deliver the signal after a fixed delay" is not the defect.
2. **The prompt text is not install-only.** `grep` for `[y/N]` across `internal/suiteinstall` finds exactly one prompt, `install_transaction.go:747` — `update-suite` reaches the same one, because updating runs a full-suite install. So the failing `update-suite/*` subtests are not waiting on a string their command never prints.

What remains is the seam at `install_transaction.go:747-748`: the prompt is written with `narrateWithoutNewline`, and only then does `readConfirmation(ctx)` run. The test's poll observes the *written bytes*, which is strictly earlier than the process being blocked in a cancellable read. Anything the process does between those two points is a window in which the signal can land somewhere that does not produce exit 130 — and the observed failure is exactly `exit = <nil>`, a clean exit rather than a missed signal.

That window is the hypothesis to confirm or refute, not a conclusion. Four sightings this session, on three distinct subtests (`install-suite/INT`, `update-suite/INT`, `update-suite/TERM`, `update-suite/HUP`), always green in isolation and only failing under full-suite parallel load — which fits a window that widens when the machine is contended.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modify) — synchronize on the process actually being blocked in its confirmation read, rather than on the prompt bytes having been written

**Files I will NOT touch:** `install_transaction.go` and the rest of the production path, unless the root cause proves to be a product defect rather than a test race — in which case the REQ's own constraint requires saying so explicitly instead of widening silently.

**Acceptance criteria (restated from REQ):**
- [ ] The signal is delivered only after the installer has demonstrably reached its blocked confirmation read
- [ ] Exit 130 and the existing managed-path non-effect assertions are unchanged — this is a synchronization fix, not a weaker assertion
- [ ] No sleep is lengthened to paper over the race, and the test is neither skipped nor made flaky-tolerant
- [ ] The fix holds under `go test ./...` parallel load, not only under `-run` in isolation
