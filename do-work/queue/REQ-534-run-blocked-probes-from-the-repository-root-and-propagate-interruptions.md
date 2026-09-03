---
id: REQ-534
title: 'Review fix: run blocked probes from the repository root and propagate interruptions'
status: pending
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
write_set: [skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go, skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go, skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go]
---

# Run Blocked Probes From the Repository Root and Propagate Interruptions

## What

The `next` command's blocked-probe runner has three defects in one file. It runs `blocked_check` in the caller's working directory instead of the selected repository root, so a relative probe gives the wrong answer under `--repo-root`. It reports a SIGINT/SIGHUP/SIGTERM received during a probe as an ordinary probe failure, so `next` swallows Ctrl-C, keeps evaluating, and can select another REQ with exit code 0. And it installs its signal handler only after the child has started, leaving a window where the parent dies by default action while the isolated child process group keeps running.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares these root causes. REQ-505 (moving selection and claim behind `advance`) re-homes the caller of `next`; it does not touch the probe runner.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
