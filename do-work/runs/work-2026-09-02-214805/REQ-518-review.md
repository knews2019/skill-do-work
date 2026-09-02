## Review: REQ-518

**Approve with follow-up** — the green-gate evidence commands are sound and well tested, but REQ-518 alone does not yet make the next REQ skip its baseline; the existing REQ-519 must close that integration-boundary gap.
Route C | uncommitted

### What's built

- `check-green-gate` and `record-green-gate` provide typed, exact-argv, repository-bound evidence in Git-private storage.
- Exact `HEAD` and `_dev/gate-runs/`-only descendants match; different argv/repositories, divergent or missing revisions, project commits, malformed records, and unsafe targets fail closed or return a valid non-match.
- The work action checks evidence before the baseline, records every direct green baseline/final run, and keeps Step 6.5's full gate unconditional.
- The per-REQ time saving is not operational until the full gate moves after finalization/release, because those later commits currently invalidate the just-written record.

### Decisions / risks for you

- Keep REQ-519 as the owner of the remaining integration-boundary work. Its acceptance test should run two consecutive REQs and count full-gate invocations; REQ-518 should not be treated as independently delivering “once per REQ” until that test passes.
- No new follow-up REQ was created, per the review instruction. The Important finding should fold into pending REQ-519.

### Findings

**Important:**

- The final record is written at Step 6.5 before the successful release tail (`skills/do-work/actions/work.md:531`, then finalization at `skills/do-work/actions/work.md:679-705`). The reference requires every successful REQ to change and commit its lifecycle/release files after that gate (`skills/do-work/actions/work-reference.md:1000-1054`), while the new matcher deliberately invalidates every non-`_dev/gate-runs/` commit (`skills/do-work/actions/work-reference.md:392,416`). The next REQ therefore observes `matches: false` and reruns the full baseline, so the headline requirement and the explicit “next REQ in the same run skips its baseline” requirement are only partially delivered. The command-level smoke test reproduced this exact sequence: record → exact match; add a lifecycle/release-shaped project commit → `invalidated_by_non_gate_log_commit`, `matches: false`; add only a gate-log commit after a new record → `gate_log_descendant_match`, `matches: true`. — impact-user-visible → route by folding the closure test and integration-boundary fix into existing pending REQ-519; do not create a new REQ.

**Minor:**

- None.

**Nit:**

- None.

### Requirements Checklist

- [x] Durable evidence records the canonical repository identity, exact argv and digest, direct-zero provenance, and full current revision outside hand-edited pipeline state.
- [x] Evidence survives separate CLI invocations and is shared through the Git common directory with linked worktrees.
- [x] Exact `HEAD` and descendants containing only `_dev/gate-runs/` changes match.
- [x] Different argv/repository, missing or divergent revisions, project-changing history, malformed content, unsafe targets, and non-private records never match; unverifiable evidence fails closed.
- [x] Step 5.75 consumes typed evidence, skips only a match, falls back to the direct baseline on a valid non-match, and stops on check/record failure.
- [x] Direct green baseline and Step 6.5 runs record evidence; Step 6.5 remains direct and mandatory.
- [x] Fingerprinting, detached-base attribution, gate deferral, and repository-gate-repair branches remain unchanged.
- [x] The public CLI seam proves record, exact match, project invalidation, and gate-log-only descendant behavior.
- [x] The non-green recorder refuses without mutation and names the gate argv rather than itself as the next action.
- [x] Existing contract predicates were strengthened line-neutrally; `_dev/tests/contract-regressions.sh` remains 8,440 lines with four additions and four deletions.
- [ ] A completed REQ leaves evidence reusable by the next REQ, so a normal run performs only one full gate per REQ — not delivered until REQ-519 moves the gate after the release/finalization commit or otherwise closes the same boundary without broadening transparent history.

### Scope and Restatement Sweep

- `scope-drift.sh` passed: the Implementation Summary and declared 12-file Scope match exactly. The run artifacts and previously identified AI report are outside builder scope and were not attributed to REQ-518.
- The repository sweep covered both command names, `gate_evidence`, `baseline_revision`, recorded-green phrasing, exact-argv semantics, and gate-log transparency. `work.md`, `work-reference.md`, the CLI/result model, tests, prime, and existing contract predicates agree on the shipped command contract.
- Pending REQ-519 is the only additional semantic consumer found. It is not a stale restatement: it explicitly owns moving the full gate to the integration commit and deriving its changed range from the recorded revision. It is also the required route for the Important finding above.
- No new debug artifacts, unrelated refactor, sentence predicate, or undeclared source path was found.

### Acceptance Testing

**Result: Partial**

- `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-518-run-the-full-gate-once-per-req.md` — PASS.
- `go test -count=1 ./cmd/do-work-cli ./internal/gateevidence ./internal/resultmodel` — PASS.
- `go test -race -count=1 ./internal/gateevidence` — PASS.
- `go vet ./...` — PASS.
- `go test -count=1 ./...` for the full CLI module — PASS.
- `bash _dev/tests/do-work-cli-go125-compatibility.sh` — PASS.
- `bash _dev/tests/contract-regressions.sh` — PASS, including the strengthened baseline/final-gate mutations.
- Direct public-CLI smoke test in a temporary Git repository — PASS for record, exact match, project-commit invalidation, replacement, and gate-log-only descendant matching. It also reproduced the next-REQ miss described in the Important finding.
- `git diff --check` — PASS; the contract file stayed line-neutral.
- The work action records the canonical unpiped repository gate as green for this integrated state. Review did not rerun the 6.5-minute gate because Step 7 permits using the just-completed Step 6.5 evidence.
- RED evidence was handled as reproduced evidence, not contemporaneous builder telemetry: the builder retained no original failing transcript, but the work action reran the public test against the detached pre-implementation revision and observed the expected assertion-level `UNKNOWN-COMMAND` failure. Together with the current GREEN run, this establishes the behavioral transition; it does not independently prove the builder's moment-by-moment TDD ordering.

### Suggested Additional Testing

- After REQ-519 lands, run two ordinary REQs consecutively and instrument the canonical full-gate entry point. Require exactly one full-gate invocation per REQ and prove the second REQ's baseline uses a recorded match after the first REQ's lifecycle/release commit.
- Repeat that end-to-end check with REQ-523's gate-log commit enabled to verify the recorded revision is the post-log integration revision and remains reusable after a session restart.

### Coding-Guardrails and Domain Review

- **Think Before Coding:** the handback records the Git-common-dir, transparent-history, action-owned execution, and shipped-launcher decisions. The known release-boundary limitation was surfaced rather than hidden.
- **Simplicity First:** the implementation is a small standard-library package behind the existing handler/result/atomic-file seams; no new dependency or speculative configurability was added.
- **Surgical Changes:** all implementation paths are declared and every changed line traces to command registration, evidence mechanics/results/tests, the two action consumers, the prime index, or existing contract predicates.
- **Goal-Driven Execution:** public-seam, focused, race, full-module, Go 1.25, contract, and command-level checks are substantive. The detached RED is valid regression evidence with the historical limitation stated above.
- **Naming for Reach:** introduced exported commands/types and new files use descriptive, greppable names; the two CLI subcommands fall under the single-word-by-design invocation exemption.
- **Backend/testing:** invalid inputs and unsafe persistence states fail closed; writes are atomic and private; the check path is read-only; tests exercise the public CLI plus storage/history boundaries without adding a second framework.

### Scores (on the record — not the headline)

**Overall: 82%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 80% | The evidence mechanism and action seams are delivered, but the normal next-REQ skip is deferred to REQ-519. |
| Code Quality | 95% | Clear typed model, exact argv handling, fail-closed Git validation, private atomic persistence, and focused action integration. |
| Test Adequacy | 93% | Strong public, unit, race, full-module, compatibility, and contract coverage; detached RED proves the pre-change failure but not original builder chronology. |
| Scope | 100% | Mechanical drift check passed; all 12 declared paths and no implementation extras. |
| Risk | Low | No integrity/security defect found; current risk is that the intended 6.5-minute saving is not realized until REQ-519. |
| Acceptance | Partial | Commands work end to end, but the next normal REQ still reruns its baseline after finalization/release. |

Raw percentage average: 92%; Acceptance Partial applies the documented 10-point penalty.

### Follow-up Routing

- New REQs created: None (prohibited by review instruction).
- Recommended fold: append the Important finding's two-REQ invocation-count acceptance case to existing pending REQ-519, whose root cause is the same final-gate integration boundary.

### Self-Validation

- Rechecked the negative paths, action ordering, release/finalization ordering, public command outputs, changed-path set, and repository-wide restatements after the initial review.
- The self-validation confirmed one material gap only: the known Step 9 invalidation prevents REQ-518 from independently delivering its headline outcome. No second Important finding emerged.

*Reviewed by review-work action*
