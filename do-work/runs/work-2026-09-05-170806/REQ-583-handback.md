# Hand-back — REQ-583 (pinning the evidence-gate remedy redirection, the layered guard, and the interrupted focused-test code)

Branch head: `0e04630a077cf4122a04f72895a92c2da9357d94` on `worktree-agent-REQ-583-pin-the-evidence-gate-remedy-redirection-guard-and-interrupted-path`.

## File manifest

| Verb | Path |
|---|---|
| modified | `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go` |

`git diff --stat` against the branch base shows that one file, 207 insertions, 0 deletions. Nothing else is changed and nothing else is staged. `internal/lifecycleadvance/evidence_gates.go` and `internal/corehelpers/commands.go` were mutated temporarily and restored; both were verified back to their pre-work SHA-256 after every mutation (`9d53f5e26a8a82a76c5fcccf26c238a8d3120dd8f36be389efcec6b0c8d0bb75` and `56da07c445df56c57818814b7c2f4ab6c660b3e9e8a9705fc8c2a3f96bfd56b6`).

Three tests plus one helper were added:

- `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation` — two subtests, one per `redirectHelperRemedies` call site.
- `TestFocusedGateStateKeepsSubordinateAuthority` — a nine-row in-process table over `focusedGateState`.
- `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure` — two rows, an interrupted run and an ordinary red run.
- `advanceGateFinding` — returns the finding a gate record carries under a code, so a test can read its remedy instead of only its presence. `gateHasFinding` stays where it is; nothing was rewritten to use the new helper.

## Red-Green Evidence

Every test below was written while its mutation was applied, and only then checked against the reverted tree. That order is the REQ's own requirement.

### M1 — the remedy redirection

The REQ's mutation is both call sites at once. Each site was also mutated on its own, because the test's own comment claims either deletion is caught.

**M1 (both sites).** Deleted `record.Findings = redirectHelperRemedies(record.Findings, gateID, advanceGatePhaseContinuation(advance, inputs))` from `composeCoreGate` and `record.Findings = redirectHelperRemedies(record.Findings, handlerName, advanceGatePhaseContinuation(advance, inputs))` from `composeGreenGate`.

Failed test: `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation`.

```
--- FAIL: TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation (0.52s)
    evidence_gates_test.go:493: the probe finding's remedy was not redirected to this advance invocation:
         got []string{"do-work-cli", "run-blocked-check", "--probe-file", "focused.sh"}
        want []string{"do-work-cli", "--format", "json", "advance", "REQ-726", "--request-path", "do-work/working/REQ-726-fixture.md", "--gate-arg", "/usr/bin/true", "--", "--probe-file", "focused.sh", "--timeout-seconds", "2"}
```

**M1a (`composeCoreGate` only).** Same failure, same test, subtest `core gate call site`:

```
--- FAIL: TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation (0.73s)
    evidence_gates_test.go:493: the probe finding's remedy was not redirected to this advance invocation:
         got []string{"do-work-cli", "run-blocked-check", "--probe-file", "focused.sh"}
        want []string{"do-work-cli", "--format", "json", "advance", "REQ-726", ...}
```

**M1b (`composeGreenGate` only).** This is the one that needed a new construction. With only that site deleted, `go test ./internal/lifecycleadvance` was **fully green**, including the `core gate call site` subtest — the green-gate call site has no natural producer (see `## Exploration` below). The `green gate call site` subtest supplies the condition directly and reddens:

```
--- FAIL: TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation/green_gate_call_site (0.01s)
    evidence_gates_test.go:539: the green gate's remedy was not redirected to this advance invocation:
         got []string{"do-work-cli", "record-green-gate"}
        want []string{"do-work-cli", "--format", "json", "advance", "REQ-727", "--request-path", "do-work/working/REQ-727-fixture.md", "--gate-exit-status", "3", "--gate-arg", "do-work-cli", "--gate-arg", "record-green-gate", "--", "--probe-file", "focused.sh", "--timeout-seconds", "2"}
```

**After revert:** `--- PASS: TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation (0.65s)`, both subtests PASS.

Negative controls in the same test, all passing: the sibling `FOCUSED-BASELINE-MISSING` finding keeps `next_argv` and `verification_argv` null; the green gate's `GATE-EVIDENCE-CHECK-FAILED` remedy (`git status --short`) keeps `git` as its argv[0]; and the `GATE-EVIDENCE-NOT-GREEN` finding's `verification_argv`, which names `check-green-gate` rather than the `record-green-gate` this gate ran, is left alone while its `next_argv` is rewritten. That last one is one finding with one field rewritten and one field not, which is the strongest available evidence that the rewrite keys on the verb.

### M2 — the layered guard

Failed test for all three: `TestFocusedGateStateKeepsSubordinateAuthority`.

**M2a.** Deleted `subordinateState == resultmodel.AdvanceGateFailed ||` from the guard.

```
--- FAIL: TestFocusedGateStateKeepsSubordinateAuthority/failed_subordinate_against_a_green_baseline (0.00s)
    evidence_gates_test.go:593: a saved baseline decided the gate for subordinate=failed launched=true timed_out=false baseline=green: got satisfied, want failed
--- FAIL: TestFocusedGateStateKeepsSubordinateAuthority/failed_subordinate_against_a_matching_red_baseline (0.00s)
    evidence_gates_test.go:593: a saved baseline decided the gate for subordinate=failed launched=true timed_out=false baseline=matching_red: got satisfied, want failed
```

**M2b.** Restored M2a, then deleted `|| focusedTest.TimedOut`.

```
--- FAIL: TestFocusedGateStateKeepsSubordinateAuthority/timed_out_against_a_green_baseline (0.00s)
    evidence_gates_test.go:593: a saved baseline decided the gate for subordinate=findings launched=true timed_out=true baseline=green: got satisfied, want findings
--- FAIL: TestFocusedGateStateKeepsSubordinateAuthority/timed_out_against_a_matching_red_baseline (0.00s)
    evidence_gates_test.go:593: a saved baseline decided the gate for subordinate=findings launched=true timed_out=true baseline=matching_red: got satisfied, want findings
```

**M2c (not in the REQ, added for completeness).** The guard has three halves, not two. Deleted `!focusedTest.Launched ||`:

```
--- FAIL: TestFocusedGateStateKeepsSubordinateAuthority/never_launched_against_a_green_baseline (0.00s)
    evidence_gates_test.go:593: a saved baseline decided the gate for subordinate=findings launched=false timed_out=false baseline=green: got satisfied, want findings
```

**After revert:** `ok ... 0.275s`, all nine rows pass. Four of the rows are negative controls: an execution the guard admits is still classified by its baseline (green and matching-red clear it, new-red reddens it, not-compared falls through), so the guard is a boundary and not a blanket refusal.

### The interrupted focused test

Failed test for both: `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure`.

**M3 (the REQ's mutation).** In `internal/corehelpers/commands.go`, changed `case runError != nil:` back to emitting `BLOCKED-PROBE-LAUNCH-FAILED`.

```
--- FAIL: TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure/interrupted_run (0.98s)
    evidence_gates_test.go:640: gate "run-blocked-check" carries no BLOCKED-PROBE-FAILED finding: []resultmodel.CommandFinding{resultmodel.CommandFinding{Code:"BLOCKED-PROBE-LAUNCH-FAILED", Severity:"error", ...
```

**M3b (the collapse the pairing exists for).** Deleted the whole `case runError != nil:` arm so an interrupted run falls through to the ordinary non-zero arm:

```
--- FAIL: TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure/interrupted_run (1.21s)
    evidence_gates_test.go:642: an interrupted probe and a probe that chose 143 for itself were classified alike: severity=warning want error, gate state=findings want failed
```

**After revert:** both rows PASS. The `ordinary red run` row stayed green under both mutations, which is what makes it a control rather than a duplicate.

Observed facts for the interrupted row, matching the orchestrator's exploration exactly: `exit_status 143`, `launched true`, `timed_out false`, gate state `failed`, CLI exit 2, JSON still rendered, wall about 0.9 s per invocation. Both rows exit 143, so the only difference between them is whether the runner observed an interruption — never the integer the child picked.

### Stability and timing

- `go test -count=8` over the three new tests: `ok ... 6.304s`. No flake.
- `go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/lifecycleadvance`: **green**, package wall **24.5 s** after (20.9 s before). The package number moved mostly with machine load — other test files in the same package rose too on the same run.
- Per-file budget (the number `_dev/tests/run-go-tests-with-budget.sh` enforces, which sums top-level test Elapsed): `evidence_gates_test.go` **6.27 s → 6.99 s**, against a 30 s budget. The three new tests cost 0.57 s (interrupted probe), 0.04 s (remedy redirection) and 0.00 s (the `focusedGateState` table).
- `go -C skills/do-work/tools/do-work-cli test -count=1 ./...`: exit 0.
- `go vet ./internal/lifecycleadvance`: clean. `gofmt -l internal/lifecycleadvance/`: no output.

## P-A-U

### [PLAN]

Read `prime-do-work-cli.md`, then `lessons-do-work-cli.md` in full (the touch-conditional rule fires: the prime's Read-first entries name both `internal/lifecycleadvance/` and `internal/nextselection/`). Read the REQ's Exploration and confirmed each claim against the source before writing anything.

Three pins, one file:

1. **M1** — drive the real CLI at the test-gate phase on a fixture that is deliberately not a Git repository and has no saved baseline, so one invocation yields a redirected probe remedy, a remedy-less baseline finding, and a git remedy from the green gate. Assert on the finding's `next_argv`/`verification_argv` against the exact continuation. Verify each call site separately.
2. **M2** — an in-process table over `focusedGateState`, keyed on `Launched` / `TimedOut` / `BaselineState` per `[family: closed-enumeration-for-a-condition]`, never on exit statuses. Cover all three halves of the guard plus four admitted-execution controls.
3. **The interrupted focused test** — a public CLI case with a `kill -TERM $PPID; sleep 5` probe, paired with `echo failing; exit 143`. Both rows exit 143 so the condition, not the integer, is what separates them.

Conflicts with the loaded lessons: none. Two lessons shaped the work rather than opposing it — `[family: closed-enumeration-for-a-condition]` (key rows on the condition's ingredients, and run one differential per half, because a single revert cannot isolate a condition when a rename reddens everything) and `[family: smoke-vs-characterization]` (a registration-shaped assertion stays green while semantics diverge, which is exactly the defect being closed).

### [APPLY]

Applied exactly as planned, in the REQ's own order, with each mutation kept applied while the test that catches it was written. Sequence: M1 both sites → M1a → M1b → revert/verify green → M2a → M2b → M2c → revert/verify green → M3 → M3b → revert/verify green. Product files restored by copying the pristine originals back and comparing SHA-256, not by re-editing.

One thing the plan did not anticipate: M1b left the package fully green, so the green-gate subtest had to be constructed rather than observed. See `D-01`.

### [UNIFY]

- `git diff --stat`: `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go | 207 +++++`, 1 file changed, 207 insertions, 0 deletions.
- Reviewed the whole diff line by line. It contains: one added import (`slices`), one helper (`advanceGateFinding`), three test functions with their comments. No `fmt.Println`, no `t.Log`, no `t.Skip`, no commented-out code, no scratch paths, no leftover mutation. No existing test, helper or constant was edited or reordered.
- Verified `internal/lifecycleadvance/evidence_gates.go` and `internal/corehelpers/commands.go` are byte-identical to their pre-work state by SHA-256, and that `git status --short` lists neither.
- `go vet ./internal/lifecycleadvance` clean; `gofmt -l internal/lifecycleadvance/` prints nothing.
- `go test -count=1 ./internal/lifecycleadvance` green; `go test -count=1 ./...` exit 0.
- Files checked, and what was checked in each: `evidence_gates_test.go` — the entire added block for scope, debug artifacts, naming, and that every new test's comment states what it pins and which deletion catches it. `evidence_gates.go` and `commands.go` — restoration by digest. No other file was opened for writing at any point.

## Decisions

**D-01 — ESCALATE. The `composeGreenGate` call site has no natural producer, so its test supplies the redirection condition through the gate argv.**

Deleting only the `redirectHelperRemedies` call in `composeGreenGate` leaves the whole package green, and it still would after M1's natural-producer test. The reason is structural: of the findings the two green-gate helpers can return, `GATE-EVIDENCE-CHECK-FAILED` and `GATE-EVIDENCE-RECORD-FAILED` carry a `git` remedy, `GATE-EVIDENCE-NOT-GREEN` carries the caller's own gate argv, and `GATE-EVIDENCE-USAGE` — the only one whose remedy re-enters the helper — cannot be produced from advance, because `composeGreenGate` always builds a well-formed argv for it. So nothing an ordinary canonical gate argv produces reaches that rewrite.

The choice was between leaving the site unpinned (and failing the REQ's "deleting **either** call site fails a named test") or supplying the condition the redirection keys on — argv[0] is `do-work-cli` and its verb is the subordinate command — through the one channel that can carry it: the gate argv, which advance copies verbatim into the not-green remedy and never executes. I took the second. The subtest passes `--gate-exit-status 3 --gate-arg do-work-cli --gate-arg record-green-gate`.

**Value:** the REQ's acceptance criterion is met for both call sites, and the pin is on the stated condition rather than on the spellings today's callers produce, which is what `[family: closed-enumeration-for-a-condition]` asks for. **Risk:** the fixture's gate argv is not a realistic repository gate — no real project's canonical gate is a `do-work-cli record-green-gate` invocation — so a reader could mistake it for a contrived assertion. Mitigated by a comment above the subtest that states exactly why no natural producer exists. Fully reversible: deleting the subtest returns the tree to "one call site pinned, one not", which is where M1's own evidence says it stands.

**D-02 — DECIDE & STATE. Added a third M2 row for the `!focusedTest.Launched` half.**

The REQ names two deletions; the guard has three halves. The `!Launched` half is masked today by the first half (a launch failure already arrives as `AdvanceGateFailed`), so deleting it alone was green. The row costs microseconds and completes the condition. Verified red under M2c.

**D-03 — DECIDE & STATE. The green-gate negative control asserts argv[0] is `git`, not the exact git argv.**

Asserting `{"git","status","--short"}` byte-for-byte would redden this test whenever `internal/gateevidence` reworded its own remedy, which has nothing to do with the redirection. Asserting argv[0] states the contract the redirection actually has — a remedy that does not run `do-work-cli` is left alone — and still catches a blanket rewrite.

**D-04 — DECIDE & STATE. Both interrupted-probe rows exit 143.**

Using a different exit status for the control row would have let the test pass while the classification keyed on the integer. Same status, different condition, per `[family: closed-enumeration-for-a-condition]`.

**D-05 — DECIDE & STATE. Added `advanceGateFinding` rather than widening `gateHasFinding`.**

`gateHasFinding` is used by four existing tests and answers a different question. Changing its signature would touch tests outside the write set.

**D-06 — DECIDE & STATE. The `focusedGateState` table calls product code in process; no other new test does.**

The brief authorised this and required it be stated in the test's own comment, which it is: the guard is defence in depth, its only current caller cannot produce the input, and a public test would need a second producer of `FocusedTestResult`.

## Discovered Tasks

- The `redirectHelperRemedies` call site in `composeGreenGate` (`internal/lifecycleadvance/evidence_gates.go`) is unreachable through any realistic canonical gate argv: every finding the two green-gate helpers can return carries either a `git` remedy or the caller's own gate argv, and the one whose remedy re-enters the helper (`GATE-EVIDENCE-USAGE`) cannot be produced because `composeGreenGate` always builds a well-formed argv. Worth a maintainer decision on whether the site is earned defence or should be deleted — `impact-negligible` → report only.
- `GATE-EVIDENCE-NOT-GREEN`'s `verification_argv` names `check-green-gate` while the subordinate advance ran is `record-green-gate`, so the redirection leaves it pointing at a helper the action may not be able to call, even though the next advance invocation is exactly what it should run. This may be deliberate; it is now pinned as current behaviour by the `green gate call site` subtest's second assertion, so changing it will redden a named test rather than pass silently — `impact-negligible` → report only.

## Lesson evidence

- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — read whole (the REQ's declared prime).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — read **whole**, not family-targeted. The touch-conditional Lessons Discipline rule fires because the prime's `## Read first` names both `internal/lifecycleadvance/` and `internal/nextselection/`, which is the code under test. The four families the brief called out were all present: `[family: smoke-vs-characterization]`, `[family: closed-enumeration-for-a-condition]`, `[family: reaped-by-its-own-parent]`, `[family: silent-skip-reads-as-red]`.
- Crew members read: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`. `debugging.md` was not needed — no test reached a second repair attempt.
- Nothing listed was missing.

Lessons that changed the work, rather than merely being read:

- `[family: closed-enumeration-for-a-condition]` — drove three choices: the M2 rows key on `Launched`/`TimedOut`/`BaselineState` instead of the exit statuses that produce them; both interrupted-probe rows exit 143 so the condition is the only variable; and each guard half and each call site got its own differential, because one combined revert cannot isolate a condition.
- `[family: smoke-vs-characterization]` — the reason each new test's comment names the deletion it catches, so a later reader can tell it apart from a registration check.
- `[family: reaped-by-its-own-parent]` — read before choosing the interrupt technique, and it is why the test signals the CLI through `$PPID` from inside the probe instead of reaching into the process tree from the test process.

## Integration seams

None. No line of this work belongs in a file outside the write set.

One thing the orchestrator owns, stated so it is not forgotten: `_dev/primes/prime-releases.md` says a commit that changes shipped files under `skills/` is a release, and this change is under `skills/`. The version bump, changelog entry and installed changelog mirror are finalization's business and are outside the REQ's one-file write set, so none of them was touched here.

## Exploration

The orchestrator's `## Exploration` is accurate on every claim I checked, including the exact `BLOCKED-PROBE-SUCCEEDED` rewrite, the `FOCUSED-BASELINE-MISSING` null remedy, the `git status --short` remedy being left alone, the `focusedGateState` unreachability argument, the interrupted-probe facts (`exit_status 143`, `launched true`, `timed_out false`, gate state `failed`, about 0.5 s), and the second `check-green-gate` record appearing in `needs_input` when no `--gate-arg` is supplied.

One thing it does not say, and which changed the work: **the M1 mutation is not symmetric across its two call sites.** The exploration names `run-blocked-check` as "the reliable producer" and lists the two negative controls, all of which is right for the `composeCoreGate` site at line 171. The `composeGreenGate` site at line 211 has *no* producer at all — verified by mutation, with only that line deleted the whole package including the new natural-producer test stayed green (`ok ... 19.719s`). Reaching it needs a gate argv that is itself a `do-work-cli record-green-gate` invocation, so that the not-green remedy advance copies verbatim carries the subordinate command's own verb. See `D-01`.
