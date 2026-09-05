# Independent Review — REQ-583

## Review: REQ-583

**Approve** — every one of the three behaviours is now held by a named test that fails on its own assertion, and three bug-preserving mutations I built myself were all caught.
Route B | merge range `a22ddfcf..722f5ada` (one file, 207 insertions, 0 deletions)

### What's built

REQ-583 (pinning the evidence-gate remedy redirection, the layered focused-test guard, and the interrupted-probe finding code) adds three tests and one helper to `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go`. Before this change, deleting the remedy redirection, either half of the guard, or the interrupted-probe classification arm left `go test ./internal/lifecycleadvance` fully green. Now each deletion fails a named test whose message says which behaviour is gone. No product code changed.

### Decisions / risks for you

- **D-01 (using a fake gate argv to reach the second remedy-redirection call site) is sound, and it buys more than the acceptance criterion.** I audited the builder's structural claim against `internal/gateevidence/gate_commands.go` myself rather than taking the hand-back. It holds. `composeGreenGate` always builds a well-formed argv for both helpers, so `GATE-EVIDENCE-USAGE` — the only finding whose remedy re-enters the helper naturally — is unreachable from `advance`. The other findings carry either a `git` remedy or the caller's own gate argv. There is no producer the builder missed. More importantly, that subtest is the **only** assertion in the package that pins the verb discrimination in `advanceArgvCommandVerb`. I proved this: changing `redirectHelperRemedies` to match on `argv[0] == "do-work-cli"` alone, dropping the verb check, is a real regression that hijacks a remedy naming a different helper, and it leaves the entire package green except for that one subtest's second assertion. So the escalated construction is not a test pinning its own fixture.
- **The `focusedGateState` guard is inert, not dead.** It runs on every focused-gate invocation, but no input its single caller can produce makes it change the answer: `handleBlockedCheck` calls `compareFocusedBaseline` only when the execution finished on its own, so every failed, unlaunched, or timed-out execution arrives carrying `FocusedBaselineNotCompared`, where the switch below the guard has no case and returns the same state the guard would have. Deleting the guard is behaviour-preserving today. Keeping it is right — REQ-506's remediation plan asked for the layering — but the pin now means a future simplifier who correctly identifies it as inert will redden a named test rather than notice it silently. The test's own comment says all of this in plain terms, so the trade is documented at the point a reader meets it. Worth knowing, not worth changing.
- **The interrupted-probe row signals the CLI from inside the probe.** It is safe. See Acceptance Testing below for the measurements.

### Findings

**Important:**
- `skills/do-work/actions/work.md:339` states a closed list of what the focused-test gate distinguishes — "green, matching recorded red, new red, timeout, launch failure, and an unusable `launched: false` baseline" — and the interrupted class is missing from it. That class exists in the code (`internal/corehelpers/commands.go`, the `case runError != nil` arm, `BLOCKED-PROBE-FAILED` at error severity with gate state `failed`), was introduced by REQ-506 (running the evidence gates from advance), and is exactly what this REQ has now pinned with `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure`. An orchestrator reading that sentence has no listed class for an interrupted focused test and would map it onto "launch failure", which is a different code with a different cause. The mechanical outcome is unchanged because the gate state is `failed` either way, so this is not a correctness break. It is the project's own "Closed Enumerations Go Stale" trap standing in a shipped action file. The REQ never declared this file, and it is reported here because the review's Restatement Sweep is scoped to meaning, not to the write set. — `impact-rule-change` → report only

**Minor:**
- The two `wantContinuation` literals in `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation` assert the whole advance argv byte-for-byte. The same file's existing continuation test, `TestAdvanceMissingInputContinuationsPreserveArgumentChannels`, deliberately asserts channel placement by index instead, precisely so that adding a token to `advanceGateContinuation` does not redden it. Any future addition to that builder now reddens both new subtests for a reason unrelated to redirection, and the fix is a two-place edit. I judge the trade acceptable and worth stating: the exact slice is what caught my BP-3 mutation, where the redirection pointed at a continuation that had silently dropped the probe argv, and a channel-shaped assertion would have passed that. The cost is real but the coverage it buys is real too. — `impact-negligible` → report only
- `evidence_gates_test.go` carries no build constraint, and the interrupted row now needs POSIX signal delivery, a `sh` that supports `$PPID`, and the probe's parent being the `do-work-cli` process. This repository puts `//go:build unix` on every other test file with those needs: `internal/nextselection/blocked_probe_test.go`, `internal/ownedprocess/owned_process_group_unix_test.go`, `internal/gittransaction/git_transaction_cancellation_test.go`, and the four `internal/heavyverification/*` files. On Windows, `runOwnedProbe` in `blocked_probe_windows.go` returns `Launched: false` with an error, so the row would fail on its own execution-facts assertion. This is largely pre-existing — the file's older tests already write `focused.sh` and run it through `sh -c` — and the new row deepens it rather than creating it. Two things keep the practical risk at zero: there is no CI workflow in this repository, and the prime's declared Windows checks are compile-only. I ran `GOOS=windows GOARCH=amd64 go vet ./internal/lifecycleadvance` and it exits 0, so the added test compiles cross-platform and nothing the prime mandates is broken. Adding the tag would change the whole file's build scope and is larger than this REQ. — `impact-negligible` → report only
- The `composeGreenGate` remedy-redirection call site remains unreachable through any realistic canonical gate argv, and it now costs a test to remove. The builder recorded this as a discovered task and D-01 states the reversal is one subtest deletion. The maintainer decision — earned defence in depth, or delete the site — is still open, and the pin makes the deletion slightly more visible rather than harder. Noted so it does not get lost behind a green test. — `impact-negligible` → report only

**Nit:**
- The `focusedGateState` table row named `eligible execution reddened by a new red baseline` uses `subordinateState: AdvanceGateSatisfied`, a combination the one caller cannot produce: `FOCUSED-NEW-RED` is emitted at error severity (`internal/corehelpers/commands.go:625`), so `handleBlockedCheck` has already raised the outcome to `findings` before `focusedGateState` sees it. The row is still correct and still load-bearing as a unit assertion over a pure function — it pins the `FocusedBaselineNewRed` arm of the switch, so deleting that arm reddens it. Only the row's name implies a scenario that does not occur in practice. — `impact-negligible` → report only
- `advanceGateFinding` returns the first matching finding, so a change that emitted both `BLOCKED-PROBE-LAUNCH-FAILED` and `BLOCKED-PROBE-FAILED` for one execution would pass. No path can do that today; the classification is a single `switch` with one winner. — `impact-negligible` → report only

**Positive (informational):**
- The Coding-Guardrails naming check passes cleanly. Every introduced name with reach is two words or more and survives a plain-text grep: `advanceGateFinding`, `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation`, `TestFocusedGateStateKeepsSubordinateAuthority`, `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure`.
- Every new test's comment names the deletion it catches, which is the explicit fourth acceptance criterion and the guard against the `[family: smoke-vs-characterization]` trap the prime names for this area.
- All three P-A-U boxes are marked `[x]` in the REQ, and the UNIFY narration matches what the merge range actually contains.

### Requirements Checklist

- [x] Assert on a finding's rewritten remedy, not only record-level `NextArgv`, so removing either `redirectHelperRemedies` call site fails — delivered. Both call sites reddened under mutation, and I confirmed independently that the green-gate subtest is the load-bearing one for its site.
- [x] Cover both halves of the `focusedGateState` guard independently, so removing either one fails — delivered, and a third half the REQ did not name (`!focusedTest.Launched`) is covered too under D-02.
- [x] A public case for an interrupted focused test asserts the finding code is `BLOCKED-PROBE-FAILED` — delivered, paired with an ordinary red control at the same exit status 143.
- [x] Each new test states what it pins and which deletion it catches — delivered. The comments name `composeCoreGate`, `composeGreenGate`, each guard clause, and the old finding code by name.
- [x] Constraint: tests only, none of the three behaviours changes — verified from the merge range. One file, 207 insertions, 0 deletions. I restored `evidence_gates.go` after my own mutations and its SHA-256 came back as `9d53f5e26a8a82a76c5fcccf26c238a8d3120dd8f36be389efcec6b0c8d0bb75`, which is the exact pre-work digest the hand-back recorded.
- [x] UR-119 batch constraint: the RED case is the mutation the reviewer performed, not a substitute — honored, and extended with per-call-site and per-guard-half differentials.

### Acceptance Testing

**Result: Pass**

Work was done in a detached worktree cut at `722f5ada` under the session scratchpad, which has been removed. The main tree was never touched.

- Baseline: all three new tests pass, 2.18s for the three.
- **The mutation table was not re-run.** The orchestrator's five mutations are already recorded. Instead I built three *bug-preserving* changes to answer the question a mutation table cannot: could a change that keeps the defect still make these tests pass?
  - **BP-1, drop the verb discrimination.** Rewrote `redirectHelperRemedies` to match on `argv[0] == "do-work-cli"` alone. This hijacks a remedy that names a different helper, which is a real regression. Caught, by exactly one assertion: `evidence_gates_test.go:544` in the `green gate call site` subtest. Nothing else in the package noticed.
  - **BP-2, rewrite `NextArgv` but leave `VerificationArgv` stale.** Caught at `evidence_gates_test.go:497`, showing the stale `do-work-cli --format json run-blocked-check --probe-file focused.sh`.
  - **BP-3, redirect to a continuation that drops the phase argv.** The remedy still points at this advance invocation but has lost `-- --probe-file focused.sh --timeout-seconds 2`, so an action following it would land straight back in `needs_input`. Caught in both subtests. This is the mutation the exact-slice assertion exists for.
- **Interrupted-probe flakiness and process hygiene.** 25 consecutive runs of `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure`, all green, 13.9s total (about 0.56s per iteration). No `sleep 5` process survived the run. The design is sound: the probe runs as `sh -c` under `Setpgid`, so `$PPID` is the `do-work-cli` process whose `signal.Notify` is already armed before the probe launches; `kill -TERM $PPID` has no leading minus so it targets one PID and can never reach `go test`. The 5-second sleep against a 10-second timeout gives signal delivery a wide margin, and if that margin were ever lost the failure mode is a clean assertion message about the execution facts, never a hang or a false pass.
- **Portability.** `GOOS=windows GOARCH=amd64 go vet ./internal/lifecycleadvance` exits 0, so the added test compiles cross-platform and the prime's declared Windows compile checks are unaffected. Runtime on Windows is a separate matter, covered in the Minor finding above.
- **Regression.** `go test -count=1 ./internal/lifecycleadvance ./internal/corehelpers ./internal/gateevidence` all green (27.6s, 16.1s, 4.1s). The two packages whose product code the mutations touched are included deliberately. `gofmt -l` on the package prints nothing.
- **Finding-Closure Ratchet.** REQ-583 originates from an independent review finding, so the ratchet applies. Its `## Red-Green Proof` names three mutations as the RED case and requires the package to be green with all of them reverted. Both halves verified: the named tests exist, they fail on their own assertions under the mutations, and the tree is green when restored.
- **Domain review.** `crew-members/backend.md` exists but its Quality Checks are HTTP and API shaped and do not apply to a Go test-only change. The one row that transfers, "Existing tests pass", is satisfied by the regression run above.

### Restatement Sweep

The diff is tests only and redefines nothing, so the sweep was run on the contracts these tests now pin rather than on anything the diff changed. Three elements were swept across the shipped packages. **The remedy-redirection contract**: `redirectHelperRemedies` and `advanceArgvCommandVerb` appear in no shipped Markdown at all, only in `do-work/` request and run files, so there is nothing to go stale. **The focused-test gate's outcome classes**: one restatement exists, at `skills/do-work/actions/work.md:339`, and it is stale — it enumerates six classes and omits the interrupted one, which is the Important finding above. `skills/do-work/actions/work-reference.md` has no competing enumeration. **The blocked-probe finding codes**: `BLOCKED-PROBE-FAILED`, `BLOCKED-PROBE-LAUNCH-FAILED`, `BLOCKED-PROBE-TIMED-OUT`, and `BLOCKED-PROBE-SUCCEEDED` appear in no shipped file, only in archived requests and run records, which are historical evidence and correctly left alone. I also checked whether the earlier REQ-534 review finding about `actions/work.md` claiming a probe "never halts the work loop and never raises an error" still stands. It does not — the phrase is gone from the file, so it was already repaired and is not re-reported here. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and `lessons-do-work-cli.md` were both read for restatements of the three pinned behaviours and carry none; the prime's `internal/lifecycleadvance/` and `internal/gateevidence/` descriptions stay at an altitude these tests do not touch.

### Suggested Additional Testing

- Run the interrupted-probe row once on Linux before relying on it there. It passed 25 of 25 on darwin, and `$PPID` plus SIGTERM through `sh -c` is POSIX, but a different `/bin/sh` (dash rather than bash) has never been exercised for this row.
- If a second producer of `resultmodel.FocusedTestResult` is ever added, revisit `TestFocusedGateStateKeepsSubordinateAuthority`. The guard stops being inert at that moment and the in-process table should be joined by a public case at the CLI seam.
- Before deleting the `composeGreenGate` remedy-redirection call site, if that decision is ever taken, keep the verb-discrimination assertion. BP-1 showed it is the only thing in the package that pins it, and deleting the subtest wholesale would take that pin with it.

### Scores (on the record — not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All four acceptance criteria plus both constraints delivered and independently verified |
| Code Quality | 90% | Clear characterization comments, real negative controls, greppable names; two style points deducted for exact-slice assertions diverging from the file's own convention and for the deepened Unix runtime dependency without the repo's build-tag convention |
| Test Adequacy | 96% | Three self-built bug-preserving mutations all caught; assertions pin code, severity, gate state and verb discrimination; 25 repetitions with no flake |
| Scope | 100% | One declared file, one touched; product files digest-identical to their pre-work state |
| Risk | None | Test-only, no product change, stdlib-only import, +0.72s against a 30s per-file budget |
| Acceptance | Pass | Package and both adjacent packages green; interrupted probe measured, not assumed |

### Follow-ups created

None (6 findings report only)

---

## Review

**Overall: 96%** | 2026-09-05T18:12:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `skills/do-work/actions/work.md:339` enumerates what the focused-test gate distinguishes and omits the interrupted class that REQ-506 introduced and this REQ pins, so an orchestrator has no listed class for an interrupted focused test and would map it onto "launch failure" — `impact-rule-change` → report only

**Minor findings:**
- The two `wantContinuation` literals assert the full advance argv byte-for-byte while the same file's `TestAdvanceMissingInputContinuationsPreserveArgumentChannels` deliberately asserts channel placement, so a future token added to `advanceGateContinuation` reddens both new subtests for an unrelated reason; the trade is acceptable because the exact slice is what caught a dropped-phase-argv mutation — `impact-negligible` → report only
- `evidence_gates_test.go` has no `//go:build unix` tag while every other test file in this repository needing POSIX signals, `sh`, or process ownership carries one, and the interrupted row needs all three; largely pre-existing, zero practical risk (no CI workflow, and `GOOS=windows GOARCH=amd64 go vet ./internal/lifecycleadvance` exits 0) — `impact-negligible` → report only
- The `composeGreenGate` remedy-redirection call site is still unreachable through any realistic canonical gate argv, verified independently from `internal/gateevidence/gate_commands.go`; the maintainer decision to keep it as earned defence or delete it is open and now costs one subtest to reverse — `impact-negligible` → report only

**Nit findings:**
- The `eligible execution reddened by a new red baseline` table row uses `subordinateState: AdvanceGateSatisfied`, which the one caller cannot produce because `FOCUSED-NEW-RED` is error severity; the row still correctly pins the `FocusedBaselineNewRed` switch arm, only its name implies a scenario that does not occur — `impact-negligible` → report only
- `advanceGateFinding` returns the first matching finding, so an implementation emitting two blocked-probe codes at once would pass; no path can do that today — `impact-negligible` → report only

**Acceptance:** Pass — three self-built bug-preserving mutations were all caught, the interrupted probe ran 25 of 25 with no flake and no orphaned process, and `internal/lifecycleadvance`, `internal/corehelpers` and `internal/gateevidence` are all green.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*
