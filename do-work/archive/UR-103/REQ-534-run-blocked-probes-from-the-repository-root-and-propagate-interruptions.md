---
id: REQ-534
title: 'Review fix: run blocked probes from the repository root and propagate interruptions'
status: completed
priority: now
created_at: 2026-09-03T12:20:21Z
user_request: UR-103
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-535, REQ-536]
batch: validate-feedback-2026-09-03
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: [skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go, skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go, skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go, skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go]
claimed_at: 2026-09-03T23:34:31Z
route: C
planning_at: 2026-09-03T23:38:54Z
exploration_at: 2026-09-03T23:38:54Z
dispatch_at: 2026-09-03T23:42:48Z
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-09-03T23:38:54Z
  basis:
    - Route C
    - 7-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
completed_at: 2026-09-04T00:24:00Z
commit: 8b488b5c7e89
---

# Run Blocked Probes From the Repository Root and Propagate Interruptions

## What

The `next` command's blocked-probe runner has three defects in one file. It runs `blocked_check` in the caller's working directory instead of the selected repository root, so a relative probe gives the wrong answer under `--repo-root`. It reports a SIGINT/SIGHUP/SIGTERM received during a probe as an ordinary probe failure, so `next` swallows Ctrl-C, keeps evaluating, and can select another REQ with exit code 0. And it installs its signal handler only after the child has started, leaving a window where the parent dies by default action while the isolated child process group keeps running.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares these root causes. REQ-505 (moving selection and claim behind `advance`) re-homes the caller of `next`; it does not touch the probe runner.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Capture the selected root without widening the stable runner interface, register Unix signals before launch, propagate a typed interruption through selection, and retain the established exit-code owner.
- [x] **[APPLY]:** Implemented the root-aware execution path, pre-launch signal registration, typed interruption/short-circuit behavior, platform parity, and three acceptance fixtures in the exact seven-file scope.
- [x] **[UNIFY]:** Reviewed all seven files; focused, race, full-module, vet, Windows cross-compile, whitespace, and debug-artifact checks passed.

## Context

Finding provenance, carried per the Finding-Closure Ratchet. Three accepted findings from the same external review, adjudicated by `do-work validate-feedback` on 2026-09-03; the full blocks are preserved in UR-103 input.md.

**Finding 1 (P1, Accept) — probe working directory.**
- Verbatim claim: "[P1] Run blocked probes from the selected repository — blocked_probe_unix.go:15-15. When next is invoked from outside the target repository using --repo-root, a relative blocked_check inherits the caller's working directory because this command has no Dir. Checks such as test -e config/ready can therefore produce the opposite result and incorrectly unblock or exclude work from the selected repository."
- Evidence: `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go:15` builds `sh -c` with no `Dir`; `internal/commandruntime/command_runtime.go:95-145` records `--repo-root` and never chdirs. Both shipped callers pass `--repo-root` (`skills/do-work/actions/work.md:143`, `skills/do-work-board/justfile.template:32`). Reproduced on a scratch build: `blocked_check: 'test -e config/ready'` probes exit 0 from the repo root and exit 1 (`BLOCKED-PROBE-FAILED`) from another directory with `--repo-root`.
- Surface-cost: N/A, direct fix.

**Finding 2 (P1, Accept) — interruption swallowed.**
- Verbatim claim: "[P1] Propagate probe interruptions out of next — blocked_probe_unix.go:49-50. When SIGINT, SIGHUP, or SIGTERM arrives during a blocked probe, this returns 128+signal as an ordinary probe failure with no error. evaluateCandidate consequently records BLOCKED-PROBE-FAILED and Select continues with other candidates and returns success, so Ctrl-C can be ignored and another request may still be selected."
- Evidence: `blocked_probe_unix.go:44-50` returns `128+signal, nil`; `next_selection.go:184-194` treats any non-zero as `BLOCKED-PROBE-FAILED` and continues. Reproduced: SIGINT during a `sleep 20` probe gave exit code 0, outcome `success`, the blocked REQ excluded with probe exit 130, and the next pending REQ selected. Child cleanup was correct (no leftover process); only propagation is wrong. `internal/toolboxcommands/report_image.go:159,229` already show the intended shape (`mutationSignalContext` plus `ExitCodeOverride`).
- Surface-cost: N/A, direct fix.

**Finding 3 (P2, Accept, folded into this REQ) — handler installed after launch.**
- Verbatim claim: "[P2] Install probe signal handling before launch — blocked_probe_unix.go:33-35. If an interrupt arrives after command.Start() but before this signal.Notify, the CLI takes the default signal action without terminating or reaping the newly isolated process group. Because the child is in a different process group, the probe and its descendants can continue after the parent exits."
- Evidence: `blocked_probe_unix.go:17` starts the child; `signal.Notify` is at line 34; the child sits in its own process group (line 16).
- Surface-cost: N/A, a reorder.

## Requirements

- The probe command runs with its working directory set to the selected repository root. `RunBlockedProbe` (or the `ProbeRunner` it implements) receives the root; `Select` callers and the fake runner in tests are updated to the new signature.
- A SIGINT, SIGHUP, or SIGTERM received during a probe is not an ordinary probe failure. The runner returns a typed interruption; `Select` stops evaluating further candidates, the result is not `success`, and the process exit status is `128+signal` (via `ExitCodeOverride` or by re-raising the signal after cleanup), matching the toolbox commands' pattern. Existing child-tree termination and reaping behaviour is preserved.
- `signal.Notify` and its `defer signal.Stop` are installed before `command.Start()`.
- Timeout (`124`) and launch-failure (`125`) statuses keep their current meaning and rendering.
- Closing checks: a lock-in test in `internal/nextselection` for the working directory (relative probe passes only when run against the repository root from a different cwd) and one for interruption (a signal during a probe yields a non-success result, no further candidate selected, exit override `128+signal`). Both name the failure they pin.

## Red-Green Proof
**RED prompt/case:** In a temporary repository with a `blocked` REQ carrying `blocked_check: 'test -e config/ready'` and an existing `config/ready`, run `do-work-cli --repo-root <that repo> --format json next` from a different directory. Separately, with a REQ carrying `blocked_check: 'sleep 20'` and a second `pending` REQ, run `next` and send SIGINT while the probe is running.
**Why RED now:** The first run excludes the REQ with `BLOCKED-PROBE-FAILED` (probe exit 1) although the same command from the repo root selects it. The second run exits 0 with outcome `success`, records probe exit 130 as an ordinary failure, and selects the second REQ.
**GREEN when:** The first run selects the REQ with `probe_status: succeeded` regardless of the caller's cwd. The second run stops selection, reports the interruption, selects nothing further, and exits with status 130. `go test ./internal/nextselection/...` carries the two lock-in tests.
**Validation:** Inferred during capture (both RED cases reproduced during the triage)

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3124 tokens, over the 2000-token budget; `slugged: partial`, so the targeted `#interruptible-blocking-io` form is not eligible. Matched on the interruptible-blocking-io family (signal handling around a blocking wait). The owning prime is listed in `prime_files` instead.

## Assets
None.

---
*Source: `do-work validate-feedback` triage of 2026-09-03, Findings 1, 2, and 3 (Accept); full blocks preserved in UR-103 input.md.*

---

## Triage

**Route: C** - Complex

**Reasoning:** The fix is localized to `nextselection`, but it changes signal ordering, process-group cleanup, typed interruption propagation, process exit status, and cross-platform build parity.

**Planning:** Required

## Plan

1. Introduce a typed probe-interruption result and make the platform runner accept the selected repository root.
2. Set the Unix command directory and install signal notification before process launch while preserving owned-group termination and reaping.
3. Adapt `handleNext` with a root-capturing closure, then make selection short-circuit on typed interruption with a non-success result and `128+signal` exit override.
4. Add RED/GREEN coverage for root-relative probes, real interruption cleanup, selection short-circuiting, timeout/launch compatibility, and Windows signature parity.

**Plan validation:** Each requirement maps to one task. The output consumer keeps per-record probe status for ordinary failures while typed interruption bypasses candidate evaluation and controls the command result directly. Seven files are required because the build-tagged Windows runner must match the Unix signature.

*Generated from the Plan-agent findings*

## Exploration

- `commandruntime` resolves `--repo-root` without changing the caller's current directory; `handleNext` currently passes bare `RunBlockedProbe`.
- The Unix runner starts the child before installing `signal.Notify`, and returns `128+signal, nil`; `evaluateCandidate` therefore records an ordinary failed probe and continues.
- `CommandResult.ExitCodeOverride` already carries process-style interruption status, so the result model and runtime need no change.
- The existing `ProbeRunner` type can remain stable by passing a closure that captures `ExecutionContext.RepositoryRoot`.
- `blocked_probe_windows.go` is a required signature-parity touch omitted from the captured write set.

*Generated from the Explore-agent findings*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go` (modify) — typed interruption and root-aware runner contract
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go` (modify) — cwd, pre-launch signal registration, cleanup, and typed return
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go` (modify) — platform signature parity
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modify) — interruption short-circuit and selection outcome
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (modify) — bind repository root and exit override
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` (modify) — root and real-signal regression coverage
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modify) — no-further-candidate interruption coverage

**Files I will NOT touch:** commandruntime, resultmodel, or unrelated selector packages.

**Acceptance criteria (restated from REQ):**
- [ ] Relative blocked probes execute from the selected repository root even when the caller is elsewhere.
- [ ] SIGINT, SIGHUP, and SIGTERM stop selection, select no later REQ, and exit with `128+signal` after cleanup.
- [ ] Signal notification is installed before child launch.
- [ ] Existing child-tree termination and reaping remain intact.
- [ ] Timeout 124 and launch-failure 125 retain their current result shapes.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)

**What was done:** The next command now runs blocked probes from the selected repository root. Unix signal ownership is installed before the child starts; an interrupt is forwarded, the owned group is terminated and reaped, and a typed interruption stops candidate evaluation with no selection and the shell-compatible exit override. Timeout, launch failure, ordinary probe behavior, and Windows fail-closed parity remain intact. Integrated at 8b488b5c7e89d8ecf2761f1367a7daf8fffbe020 from implementation range 5c06a7d7bd2c2511d0458b4724ffa1e14e651408..8b488b5c7e89d8ecf2761f1367a7daf8fffbe020.

## Decisions

- **D-01:** Keep the public probe-runner shape stable by binding the authoritative repository root in the next-command adapter.
- **D-02:** Return a typed interruption to selection and let the existing command runtime remain the sole process-exit owner through ExitCodeOverride.
- **D-03:** Preserve the current-cwd compatibility wrapper for direct callers while the command path uses the root-aware entry point.
- **D-04:** Preserve the existing Windows launch-status refusal because standard-library process-tree ownership is unavailable there; verify signature parity by cross-compilation.

## Testing

- RED then PASS: root-relative probe, typed real-signal/group-reaping, and interruption short-circuit acceptance fixtures.
- PASS: complete nextselection package and its race suite; `go vet ./...`; `go test -count=1 ./...`.
- PASS: Windows amd64 nextselection cross-compile; `git diff --check`; no debug-artifact matches in the seven files.

## Review

**Overall: 83%** | 2026-09-04T00:12:33Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 90% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The canonical run action still says a probe never halts or raises an error, contradicting the new required typed-interruption stop at `actions/work.md:153` — `impact-rule-change` → report only

**Minor findings:**
- `blocked_probe.go:29-31` still promises a public 0–4 envelope despite the new 129/130/143 interruption override — `impact-negligible` → report only
- The handback's required-read list names six paths that did not exist at the integration base or merge end, so that durable evidence is not verifiable (`REQ-534-handback.md:91-97`) — `impact-negligible` → report only

**Acceptance:** Partial — implementation, race, full module, vet, process cleanup, and Linux/Windows compile checks pass; live restatements remain contradictory.
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
