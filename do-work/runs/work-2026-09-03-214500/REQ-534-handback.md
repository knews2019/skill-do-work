# REQ-534 builder handback

## Builder result

- Status: complete
- Branch: `worktree-agent-REQ-534-run-blocked-probes-from-the-repository-root-and-propagate-interruptions`
- Commit: `2ad658b37c8dfbe1d1930b49744ba1b2abd4b26a`
- Commit subject: `[REQ-534] run blocked probes from repository root and propagate interruptions`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-534-run-blocked-probes-from-the-repository-root-and-propagate-interruptions`

## Exact changed-file manifest

1. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go`
2. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go`
3. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`
4. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go`
5. `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go`
6. `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
7. `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`

No lifecycle, request, run-manifest, version, release, REQ-539, REQ-563, or REQ-564 file was changed in the builder worktree.

## Requirement mapping and implementation

- Repository-root execution: `handleNext` now closes over `ExecutionContext.RepositoryRoot` and supplies a root-aware runner to `Select`; the Unix command sets `exec.Cmd.Dir` to that selected root. The public `ProbeRunner` shape remains unchanged.
- Pre-start signal ownership: the Unix runner installs `signal.Notify` for `SIGHUP`, `SIGINT`, and `SIGTERM` before constructing and starting the child and defers `signal.Stop`.
- Interruption propagation: after forwarding the signal, terminating the owned process group, and reaping the child, the runner returns a typed `BlockedProbeInterruption` with `128 + signal` (129, 130, or 143).
- Stop semantics: selection recognizes the typed interruption, stops evaluating candidates, returns no selected REQs, produces one `BLOCKED-PROBE-INTERRUPTED` error finding for the interrupted request, and sets `ExitCodeOverride` to the shell-compatible interruption status.
- Existing outcome preservation: timeout remains 124, launch/validation failure remains 125, and ordinary probe success/non-zero behavior remains unchanged.
- Platform parity: the Windows implementation accepts the root-aware internal signature while preserving its existing fail-closed process-tree ownership response.

## RED evidence

Tests were added before production changes, then run from `skills/do-work/tools/do-work-cli`:

```text
go test -count=1 ./internal/nextselection -run 'TestBlockedProbe(RunsFromSelectedRepositoryRoot|InterruptionIsTypedAndReapsDescendants|InterruptionStopsSelection)' -v
```

Exit: 1.

Literal failing assertions:

```text
blocked_probe_test.go: ... status=130 err=<nil> <nil>, want typed interruption 130
next_selection_test.go: ... selected = [], want root-relative blocked probe success; ... BLOCKED-PROBE-FAILED ... probe exit 1
next_selection_test.go: ... interrupted selection outcome=success override=0, want non-success/130
```

The RED failures demonstrated all three missing behaviors: the runner did not return typed interruption evidence, the probe inherited the caller cwd instead of the selected repository root, and selection treated an interruption as ordinary probe failure and continued.

## GREEN and verification evidence

Final focused run from `skills/do-work/tools/do-work-cli`:

```text
go test -count=1 ./internal/nextselection -run 'TestBlockedProbe(RunsFromSelectedRepositoryRoot|InterruptionIsTypedAndReapsDescendants|InterruptionStopsSelection)' -v
=== RUN   TestBlockedProbeInterruptionIsTypedAndReapsDescendants
--- PASS: TestBlockedProbeInterruptionIsTypedAndReapsDescendants (0.64s)
=== RUN   TestBlockedProbeRunsFromSelectedRepositoryRoot
--- PASS: TestBlockedProbeRunsFromSelectedRepositoryRoot (0.01s)
=== RUN   TestBlockedProbeInterruptionStopsSelection
--- PASS: TestBlockedProbeInterruptionStopsSelection (0.00s)
PASS
ok github.com/knews2019/skill-do-work/do-work-cli/internal/nextselection 1.113s
```

Additional verification, all successful:

```text
go test -count=1 ./internal/nextselection -v
PASS (package elapsed 2.461s; guarded heavy go-run cases skipped as designed)

go test -race -count=1 ./internal/nextselection
ok (3.626s)

GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <temporary-directory>/nextselection.test.exe
exit 0; temporary artifact removed

go vet ./...
exit 0

go test -count=1 ./...
exit 0; all module packages passed

git diff --check
exit 0

rg -n 'TODO|FIXME|fmt\\.Print|log\\.Print' <all-seven-changed-files>
no matches
```

## Test coverage added

- `TestBlockedProbeRunsFromSelectedRepositoryRoot` deliberately changes the test process cwd and proves a root-relative `blocked_check` succeeds only from the selected repository.
- `TestBlockedProbeInterruptionIsTypedAndReapsDescendants` runs a real background descendant, delivers SIGINT, asserts status 130 plus typed interruption evidence, and proves the descendant is gone before return.
- `TestBlockedProbeInterruptionStopsSelection` proves the first interrupted probe ends the pass, later probes are not invoked, no request is selected, and the result carries failure plus exit override 130.

## P-A-U record

- Plan: introduce a root-aware internal execution path, keep the stable `ProbeRunner` interface through a closure, register Unix signals before child start, and propagate a typed interruption from runner through selector to command result.
- Apply: wrote the three acceptance tests first, captured the literal RED failures, implemented the seven-file change, and retained timeout, launch-failure, and ordinary probe behavior.
- Unify: reviewed every changed file and the complete diff; reran focused tests after formatting; passed package, race, full-module, vet, Windows compile, diff, and debug-artifact checks; committed one coherent source/test increment.

## Guidance read

- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/lessons/lessons-do-work-cli.md` (complete file, including single-exit-owner, reaped-by-own-parent, and interruptible-blocking-I/O guidance)
- `skills/do-work/crew/general.md`
- `skills/do-work/crew/backend.md`
- `skills/do-work/crew/testing.md`
- `skills/do-work/crew/coding-guardrails.md`
- `skills/do-work/crew/communication-style.md`

## Decisions, seams, and risks

- Compatibility decision: an existing direct test outside the authorized write set calls `RunBlockedProbe`. That function remains as a compatibility wrapper using the current cwd, while command execution uses the new `RunBlockedProbeAtRoot` path with the authoritative selected root. No out-of-scope edit was needed.
- Integration seam: `handleNext` owns root selection; the captured runner supplies it without widening `ProbeRunner` or changing external call sites.
- Result seam: the typed interruption crosses `evaluateCandidate` into `Select`, which sets the existing `ExitCodeOverride`; the established command runtime remains the sole process-exit owner.
- Process seam: signal delivery, owned process-group cleanup, and child reaping remain inside the process that started the child.
- Residual platform risk: Windows still deliberately returns launch status 125 because the standard-library-only process-tree ownership contract is unavailable there; this is existing behavior, and cross-compilation confirms signature parity.
- Scope deviation: none.
- Discovered follow-up tasks: none.

## Merge guidance

Cherry-pick `2ad658b37c8dfbe1d1930b49744ba1b2abd4b26a`. After integration, rerun at minimum the focused `internal/nextselection` tests, `go test -race -count=1 ./internal/nextselection`, `go vet ./...`, and `go test -count=1 ./...` from `skills/do-work/tools/do-work-cli`. Release and lifecycle finalization remain with the orchestrator.
