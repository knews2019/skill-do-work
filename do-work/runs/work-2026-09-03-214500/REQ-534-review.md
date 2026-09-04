# REQ-534 Independent Review

## Decision Brief

## Review: REQ-534

**Approve with follow-ups** — the exact implementation closes the repository-root, interruption-propagation, and pre-launch signal-registration defects without weakening timeout, launch-failure, process-group cleanup, or Windows fail-closed behavior. Two live restatements still describe the retired exit contract, and the builder handback's required-read evidence names paths that did not exist at the integration base.

Route C | merge range `5c06a7d7bd2c2511d0458b4724ffa1e14e651408..8b488b5c7e89d8ecf2761f1367a7daf8fffbe020`

### What's built

- `handleNext` binds the authoritative `ExecutionContext.RepositoryRoot`; the Unix runner assigns it to `exec.Cmd.Dir`, while the direct `RunBlockedProbe` compatibility wrapper retains current-working-directory behavior.
- Unix signal notification for SIGHUP, SIGINT, and SIGTERM is registered before `command.Start`. A received signal is forwarded to the isolated process group, the leader is waited, remaining group members are escalated, and a typed interruption reaches `Select`.
- `Select` stops immediately on the typed interruption, returns no selected records, emits `BLOCKED-PROBE-INTERRUPTED`, marks the result failed, and sets `ExitCodeOverride` to 129, 130, or 143.
- Timeout 124, launch failure 125, ordinary non-zero probe exclusions, successful-probe evidence, and the Windows standard-library ownership refusal retain their previous shapes.

### Decisions / risks for you

- No product decision is needed. The source change is safe to retain; one bounded documentation/comment correction can close the live contract drift.
- D-01 through D-04 are supported by the code: preserving `ProbeRunner` with a root-capturing closure avoids interface churn; `ExitCodeOverride` keeps one process-exit owner; the direct-call wrapper preserves compatibility; Windows remains deliberately fail-closed at 125.
- Residual platform risk is unchanged: Windows cannot run these process-tree-owned probes with the standard-library implementation. Cross-compilation proves parity, not Windows execution.

### Findings

**Important:**

- F1. The canonical run action still states that a probe “never halts the work loop and never raises an error” (`skills/do-work/actions/work.md:153`). That directly contradicts the new typed interruption path in `next_selection.go:32-43` and `blocked_probe_unix.go:45-52`, where Ctrl-C/SIGHUP/SIGTERM must halt selection and return a non-success process status. The preceding action rule at `work.md:149` already says a failed selector command stops, so the two adjacent instructions now disagree about the exact event this REQ changed. An orchestrator following the narrower restatement could misclassify `BLOCKED-PROBE-INTERRUPTED` as an ordinary failed probe and continue the run. — `impact-rule-change` → report only

**Minor:**

- F2. `RunBlockedProbe`'s live API comment says public command rendering remains in the standard 0–4 envelope (`skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go:29-31`), while the new `next` path deliberately exposes `ExitCodeOverride` 129/130/143 (`next_selection.go:32-43`). The comment should distinguish ordinary probe results from invocation interruption. — `impact-negligible` → report only
- F3. The builder handback records one lessons path and five crew paths under `skills/do-work/lessons/` and `skills/do-work/crew/` as read, but none exists at either integration base or merge end; the live sources are `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and `skills/do-work/crew-members/*.md` (`do-work/runs/work-2026-09-03-214500/REQ-534-handback.md:91-97`). This does not reveal a code defect—the implementation follows the prime's signal/process rules—but the durable handback does not truthfully substantiate the brief's required reads. — `impact-negligible` → report only

**Nit:** None.

### Requirements Checklist

- [x] Relative probes execute from the selected repository root even when the caller's cwd differs.
- [x] SIGINT, SIGHUP, and SIGTERM use a typed interruption, stop candidate evaluation, select nothing, return non-success, and surface `128+signal`.
- [x] Signal notification and deferred stop are installed before child launch.
- [x] Owned process-group termination and leader reaping are preserved on success, timeout, and interruption.
- [x] Timeout 124 and launch-failure 125 retain their meanings and result handling.
- [x] The required cwd and interruption lock-in tests were added and name the defects they pin.
- [~] Live restatements agree with the new interruption contract — the prime is correct, but `actions/work.md:153` and `blocked_probe.go:29-31` remain stale.
- [x] Scope reconciles exactly to the seven declared files; no lifecycle, request, queue, release, or unrelated source path is present in the merge range.

### Acceptance Testing

**Result: Partial**

- Focused root, interruption, timeout, raw-status, background-descendant, and selection-short-circuit tests — PASS (`2.461s`).
- `go test -race -count=1 ./internal/nextselection` — PASS (`3.651s`).
- `go test -count=1 ./...` from `skills/do-work/tools/do-work-cli` — PASS (all packages, including `internal/finalization` at `60.239s`).
- `go vet ./...` — PASS.
- `GOOS=linux GOARCH=amd64 go test -c ./internal/nextselection -o <temporary-path>` — PASS.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <temporary-path>` — PASS.
- `git diff --check 5c06a7d7bd2c2511d0458b4724ffa1e14e651408..8b488b5c7e89d8ecf2761f1367a7daf8fffbe020` — PASS.
- Exact range/scope check — PASS: seven changed files, all and only the declared `internal/nextselection` paths.
- Signal/process inspection — PASS: notification precedes `Start`, group isolation is verified before group signalling, `Wait` remains owned by the spawning process, and interruption shares the established termination path.
- Restatement/process-evidence sweep — PARTIAL on F1–F3.

### Restatement Sweep

- Confirmed aligned: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:31` explicitly assigns `nextselection` its signal-forwarding and `128+signal` contract; source and tests match it.
- Found the contradictory current consumer at `skills/do-work/actions/work.md:153` and the stale API comment at `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go:29-31`.
- Historical changelog wording was treated as historical evidence, not a live instruction.
- No second executable blocked-probe runner or alternate queue-selection consumer was found.

### Suggested Additional Testing

- Convert the real-signal fixture into a Unix table for SIGHUP, SIGINT, and SIGTERM so all three public exit statuses (129, 130, 143) are exercised dynamically; the current test dynamically covers SIGINT while the other two share the inspected code branch.
- Add a command-runtime acceptance fixture that invokes `next`, interrupts a running probe, and asserts the final executable status as well as the typed `CommandResult`; current coverage composes a real runner test, a selector test with a typed fake, and the generic runtime override test.

### Scores (on the record — not the headline)

**Overall: 83%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 92% | All executable requirements are delivered; one canonical action restatement remains contradictory. |
| Code Quality | 90% | The change is small, fail-closed, and keeps one exit owner; one API comment is stale. |
| Test Adequacy | 88% | Strong RED/GREEN, race, cleanup, and platform compile coverage; only SIGINT is dynamically exercised and the final executable exit is compositional. |
| Scope | 100% | Exact seven-file declared range with no lifecycle or unrelated changes. |
| Risk | Low | No child leak, continued selection, cwd error, timeout regression, or platform regression reproduced. |
| Acceptance | Partial | Runtime behavior passes; current consumer/comment restatements and handback read evidence need correction. |

Raw percentage average: 92.5%. The documented 10-point Acceptance Partial penalty yields 82.5%, recorded as 83% after rounding.

### Directive and Guardrail Checks

- No approach directive was assigned beyond the explicit root-capturing, typed-interruption, and platform-parity constraints; the patch follows all three.
- Think Before Coding: D-01 through D-04 and the root/result/process seams are recorded in both the REQ and handback.
- Simplicity First: the stable `ProbeRunner` closure and existing `ExitCodeOverride` avoid new public infrastructure.
- Surgical Changes: exact seven-file Scope match.
- Goal-Driven Execution: literal RED evidence and focused/race/full/vet/platform GREEN evidence are present.
- Naming for Reach: the new exported `BlockedProbeInterruption` and `RunBlockedProbeAtRoot` names are specific to their domain and do not introduce ambiguous single-word API surface.

### Self-validation

- Re-read the exact merge range, current source, REQ, UR-103 finding provenance, builder handback, prime, lessons, action consumer, runtime exit owner, and all new tests.
- Rechecked signal registration order, command cwd, process-group isolation, termination/reaping paths, typed-error recognition, selection early return, status override bounds, timeout/launch behavior, Windows signature parity, and exact scope.
- A second restatement sweep found only the two current contradictions recorded above; no additional security, provenance, destructive-action, or compatibility defect emerged.

### Follow-ups created

None (3 findings report only). This review modified only its requested artifact.

## Append-ready durable Review block

```markdown
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
```
