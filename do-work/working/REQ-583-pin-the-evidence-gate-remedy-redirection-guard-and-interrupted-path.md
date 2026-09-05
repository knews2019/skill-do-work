---
id: REQ-583
title: 'Addendum: pin the evidence-gate remedy redirection, layered guard and interrupted focused-test code'
status: claimed
created_at: 2026-09-05T01:30:57Z
user_request: UR-119
addendum_to: REQ-506
depends_on: [REQ-506]
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
write_set: [skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-09-05T17:05:17Z
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-05T17:07:24Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - async lifecycle behavior
route: B
dispatch_at: 2026-09-05T17:15:54Z
builder_handback_at: 2026-09-05T17:33:12Z
review_at: 2026-09-05T17:52:43Z
integration_at: 2026-09-05T17:33:12Z
commit: 722f5ada02df491c9b44e22d6815dc328dfa63ec
---

# Addendum: Pin the Evidence-Gate Remedy Redirection, Layered Guard and Interrupted Focused-Test Code

## What
Three pieces of behaviour in `internal/lifecycleadvance/evidence_gates.go` have no test that fails when they are removed. Delete the remedy-redirection call sites, or either half of the layered guard in `focusedGateState`, and `go test ./internal/lifecycleadvance` stays fully green. Add the tests that hold each one in place, plus a public case for the interrupted focused test whose finding code changed.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the prime and, because its Read-first entries name both `internal/lifecycleadvance/` and `internal/nextselection/`, the whole `lessons-do-work-cli.md` satellite under the touch-conditional rule. Three pins in one file: drive the real CLI at the test-gate phase for the remedy redirection and assert on a finding's remedy; an in-process table over the focusedGateState guard keyed on launched, timed-out and baseline-state; a public interrupted-probe case paired with an ordinary red run at the same exit status.
- [x] **[APPLY]:** Applied in the REQ's own order, each mutation kept applied while the test that catches it was written, then reverted: M1 both sites, M1a, M1b, revert; M2a, M2b, M2c, revert; M3, M3b, revert. Product files restored by copying the pristine originals back and comparing SHA-256 rather than by re-editing.
- [x] **[UNIFY]:** `git diff --stat` shows one file, 207 insertions, 0 deletions. The whole added block was read line by line: one import (`slices`), one helper, three test functions with their comments; no print or log calls, no skipped tests, no commented-out code, no scratch paths, no leftover mutation, and no existing test, helper or constant edited or reordered. `internal/lifecycleadvance/evidence_gates.go` and `internal/corehelpers/commands.go` verified byte-identical to their pre-work state by SHA-256 and absent from `git status --short`. `go vet ./internal/lifecycleadvance` clean; `gofmt -l internal/lifecycleadvance/` prints nothing.

*Authored by the builder in `do-work/runs/work-2026-09-05-170806/REQ-583-handback.md` → `## P-A-U`; transcribed and checked here by the orchestrator, which is the only writer of this file in worktree dispatch mode. The orchestrator re-read the merged diff itself before checking the UNIFY box.*

## Why
One root cause: new behaviour landed without a lock. Each of the three was delivered on purpose — the remedy redirection to answer a prior review finding, the layered guard because the remediation plan explicitly asked for it — and nothing stops a future change deleting any of them silently.

## Context
All three were verified by mutation. The reviewer deleted the code and the package stayed green. Line numbers below are from the main tree at commit `a55f24ce`; the reviewer's own numbers, taken in the REQ-506 worktree, are given where they differ.

**M1, the remedy redirection.** `redirectHelperRemedies` (line 343; reviewer: around 334-368) and its helper `advanceArgvCommandVerb` (line 357) rewrite a subordinate finding's remedy so it points at the continuation rather than at the subordinate command. Deleting both call sites — line 171 (reviewer: 168) and line 211 — leaves `go test ./internal/lifecycleadvance` fully green, because the existing tests assert record-level `NextArgv` and never a finding's rewritten remedy.

**M2, the layered guard.** `focusedGateState` (line 183) opens with `if subordinateState == resultmodel.AdvanceGateFailed || !focusedTest.Launched || focusedTest.TimedOut` (line 184; reviewer: 180). Deleting either half keeps the package green, because the eligibility guard `finishedOnItsOwn` in `internal/corehelpers/commands.go` (line 545, consumed at line 564; reviewer cites 565) already leaves those executions at `FocusedBaselineNotCompared`. The layering is correct and was explicitly asked for by the remediation plan — the finding is that nothing holds it there.

**The interrupted focused test.** That path's finding code changed from `BLOCKED-PROBE-LAUNCH-FAILED` to `BLOCKED-PROBE-FAILED` (`internal/corehelpers/commands.go` lines 550 and 554). Both are `failure` severity, so no existing assertion moved, and nothing pins the change. The one live assertion on either code is `evidence_gates_test.go:347`, which checks `BLOCKED-PROBE-LAUNCH-FAILED`.

This REQ amends REQ-506 (running the evidence gates from advance), which delivered all three and is in flight. It is not a defect report against that work — the behaviour is right, it is just unpinned.

## Red-Green Proof
**RED prompt/case:** Three mutations, run one at a time with `go test ./internal/lifecycleadvance`.
1. Delete the `redirectHelperRemedies(...)` call at `evidence_gates.go:171` and at `evidence_gates.go:211`.
2. Delete `subordinateState == resultmodel.AdvanceGateFailed ||` from the guard at `evidence_gates.go:184`, then restore it and instead delete `|| focusedTest.TimedOut`.
3. Change the interrupted focused-test path to emit `BLOCKED-PROBE-LAUNCH-FAILED` again in place of `BLOCKED-PROBE-FAILED`.

**Why RED now:** Every one of those mutations leaves the package green. There is no test that reads a finding's rewritten remedy, no test that reaches `focusedGateState` with an execution the eligibility guard has already excluded, and no public case for the interrupted focused test at all.
**GREEN when:** Each mutation, applied alone, fails at least one named test in `internal/lifecycleadvance`, and the failure message says which behaviour is gone. With all mutations reverted the package is green.
**Validation:** User confirmed — the mutations are the ones the reviewer performed, recorded as the RED case at the maintainer's instruction.

## Detailed Requirements
- Assert on a finding's rewritten remedy, not only on record-level `NextArgv`, so removing either `redirectHelperRemedies` call site fails.
- Cover both halves of the `focusedGateState` guard independently, so removing either one fails. Building the state needs an execution the eligibility guard in `internal/corehelpers/commands.go` already leaves at `FocusedBaselineNotCompared`; that is what makes the layering invisible to today's tests.
- Add a public case for an interrupted focused test that asserts the finding code is `BLOCKED-PROBE-FAILED`.
- State in each new test what it pins and which deletion it catches, so a later reader can tell it apart from a smoke test.

## Constraints
- Tests only. None of the three behaviours changes.
- `depends_on: [REQ-506]` because these tests pin code that request delivers; do not start before it completes.

## Builder Guidance
Apply one mutation at a time and keep it applied while writing the test that catches it. A test that passes under the mutation is the exact defect this REQ exists to close, and REQ-581 (making the descendant-cleanup tests fail on a real process-group leak) is the same failure in the same tool — worth reading before writing assertions here.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 8700 tokens at claim time (`do-work/lessons-index.md`), over the 2000-token budget in `actions/capture-reference.md` → Required Lessons Budget Contract, and `slugged: partial`, so the targeted `path#family-slug` form is not legal and the entry can only be taken whole. Matches on families `lifecycle-section-evidence`, `opaque-evidence-projection`, `smoke-vs-characterization` and `silent-skip-reads-as-red`. Dropped as a stamped entry; the touch-conditional Lessons Discipline rule still applies through `prime-do-work-cli.md`, and the builder is told to read the satellite if the prime's Read-first or Traps entries name the code it touches.

## Full Context
See `do-work/user-requests/UR-119/input.md` for complete verbatim input.

*Source: independent review of REQ-506 (running the evidence gates from advance), findings M1 and M2, work run `do-work/runs/work-2026-09-05-003420/`.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The file to change and the three behaviours to pin are named exactly, which looks like Route A. The discovery that makes it Route B is how to *reach* each behaviour from a test: the layered guard in `focusedGateState` is unreachable through the ordinary path because the eligibility guard in `internal/corehelpers/commands.go` already excludes those executions, so the test has to construct a state today's fixtures never produce. Finding that construction is exploration.

**Planning:** Not required — the Detailed Requirements are already the design, and the Red-Green Proof already names the three mutations that must fail.

## Plan

**Planning not required** - Route B: exploration-guided implementation

*Skipped by work action*

## Exploration

Explore agent, read-only, verified by building the CLI into a scratch directory and running it against throwaway fixture repos. No repository file was modified during exploration.

**Test file shape.** `internal/lifecycleadvance/evidence_gates_test.go` is 537 lines and declares `package lifecycleadvance` — the same package as the product code, so every unexported identifier in `evidence_gates.go` is directly callable from a test. Every existing test in the file drives a real built CLI binary as a subprocess and decodes JSON; no test calls product-code unexported functions in process. Reusable helpers: `runAdvanceGateJSON` (line 260), `findAdvanceGate` (281), `gateHasFinding` (293), `initAdvanceGitFixture` (442), `recordAdvanceGreenGate` (451), `gitOnlyPathDirectory` (461), plus `writeAdvanceRequest` / `writeAdvanceFile` / `advanceCLIBinary` in the sibling `advance_commands_test.go`. Constants `focusedGateRouteABody` (318) and `canonicalGateFixtureBinary` (322) are the minimal claimed-Route-A fixture.

**M1 — what the redirection actually rewrites.** `redirectHelperRemedies` (`evidence_gates.go:339-353`) changes exactly two fields **on each finding**: `NextArgv` and `VerificationArgv`. It never touches `NextJustRecipe`, and the record-level `AdvanceGateRecord.NextArgv` is written elsewhere (`missingAdvanceGateInput:262`, green-gate direct-run branch `:216`) — which is precisely why every existing assertion on `gate.NextArgv` proves nothing about this function. The reliable producer is `run-blocked-check` (`corehelpers/commands.go:577`), whose findings carry `do-work-cli run-blocked-check --probe-file <p>` before the rewrite and, after it, `do-work-cli --format json advance REQ-NNN --request-path <p> -- --probe-file <p> --timeout-seconds N`. Verified live on a Route A fixture with `focused.sh` = `exit 0`: finding `BLOCKED-PROBE-SUCCEEDED` carries the rewritten advance argv, while the sibling `FOCUSED-BASELINE-MISSING` finding keeps `next_argv: null` and the `gateevidence` git remedies (`git status --short`) are left alone because their argv[0] is not `do-work-cli`. Those two are the negative controls.

**M2 — why the layered guard is unreachable from the public seam.** `BaselineState` moves off `not_compared` only inside `compareFocusedBaseline`, which `corehelpers/commands.go:545-561` calls only when `finishedOnItsOwn` — launched, not timed out, no runner error. Every outcome that maps to `AdvanceGateFailed` arrives from `!Launched` or `runError != nil`, both of which are `!finishedOnItsOwn`, and `compareFocusedBaseline` can only raise Success to Findings, never to Failure. So `Failed`, `!Launched` and `TimedOut` all reach `focusedGateState` carrying `not_compared`, where the switch has no case and falls through to the same `return subordinateState` the guard would have produced. Deleting the guard is therefore invisible at the CLI seam, which is why the existing REQ-506 pin `TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline` passes either way. `focusedGateState` is package-private and the test file is in the same package, so a direct in-process table over `{Failed, Findings} x Launched x TimedOut x {Green, MatchingRed, NewRed}` is the only construction that reaches both halves, and it runs in microseconds with no subprocess.

**The interrupted focused test.** The dispatch in `corehelpers/commands.go:549-557` is ordered: `!Launched` gives `BLOCKED-PROBE-LAUNCH-FAILED` (failure/error); `TimedOut` gives `BLOCKED-PROBE-TIMED-OUT` (findings/warning); launched-with-runner-error gives `BLOCKED-PROBE-FAILED` at **failure/error**; a plain non-zero status gives `BLOCKED-PROBE-FAILED` at **findings/warning**. Case three is the interrupted run and has no assertion anywhere in the package — the only existing check, at `evidence_gates_test.go:347`, covers case one. A probe file containing `kill -TERM $PPID; sleep 5` reproduces it: the probe's parent is the `do-work-cli` process, whose `signal.Notify` (`nextselection/blocked_probe_unix.go:21`) is already armed, so the CLI catches the signal, tears the probe group down, still renders JSON, and exits non-zero. Observed: `exit_status 143`, `launched true`, `timed_out false`, gate state `failed`, finding `BLOCKED-PROBE-FAILED` at severity error, 0.5s wall. Pairing that row with an ordinary red row (exit 143 without the kill, which lands `launched` + no runner error and therefore gate state `findings`) is what makes "collapse case three into case four" red. Do not copy the in-process `syscall.Kill(os.Getpid(), ...)` technique from `nextselection/blocked_probe_test.go`; here the CLI is a subprocess and killing through `$PPID` needs no goroutine, no sleep and no signal guard in the test process.

**Two traps carried from the prime.** Without a `--gate-arg` the test-gate phase always adds a second `check-green-gate` record in `needs_input`, so the aggregate outcome is `findings` even when the probe passes — assert on the `run-blocked-check` record, not the aggregate, unless a green gate is recorded. And per `[family: closed-enumeration-for-a-condition]`, the M2 rows must key on `Launched` / `TimedOut` / `BaselineState`, never on the exit-status values that happen to produce them today.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go` (modify) — three new pins: a finding-level remedy assertion for M1, an in-process table over the focusedGateState guard for M2, and a public interrupted-focused-test case with its ordinary-red negative control

**Files I will NOT touch:** `internal/lifecycleadvance/evidence_gates.go` and `internal/corehelpers/commands.go` — the REQ's own constraint is "tests only, none of the three behaviours changes", so the mutations are applied and reverted, never landed. No other test file, including `internal/corehelpers/commands_test.go` (REQ-506's D-18 deliberately left it alone) and `internal/nextselection/blocked_probe_test.go`.

**Acceptance criteria (restated from REQ):**
- [ ] A test asserts on a finding's rewritten remedy, not only on record-level `NextArgv`, so deleting either `redirectHelperRemedies` call site fails a named test
- [ ] Both halves of the `focusedGateState` guard are covered independently, so deleting either one fails a named test
- [ ] A public case for an interrupted focused test asserts the finding code is `BLOCKED-PROBE-FAILED`
- [ ] Each new test states what it pins and which deletion it catches, so a later reader can tell it apart from a smoke test

## Pre-Flight

**Git:** ⚠ Clean outside `do-work/`. Inside it, a sibling session in this same checkout holds a live claim on REQ-588 and is writing `do-work/working/REQ-588-…md` and its own run directory `do-work/runs/work-2026-09-05-170800/`. Those bytes belong to that session: named here, left alone, never staged by this REQ. One shared file was collided on — `do-work/working/baseline.json` is a single path both sessions write, so the record below replaced whatever stood there; it is a stale-by-design file, not queue state, and this is recorded rather than repaired.
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at revision `9d58ae17`, the exact revision this REQ dispatches from, and the green record was written for that exact argv. Wall 110.8s. An earlier run of the same argv at `93f856ee` also exited 0 (122.9s); the only commit between the two is this REQ's own claim, which touches `do-work/` alone.
**Tests baseline:** ✓ `go -C skills/do-work/tools/do-work-cli test ./internal/lifecycleadvance` exited 0, launched true — a usable green baseline, so any later red in this package is attributable to the builder.
**Dependencies:** ✓ Go toolchain present and the module builds; `go build ./...` and `go vet ./internal/lifecycleadvance` both clean in the builder worktree before dispatch.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go` (modified)

**What was done:** Three behaviours that could be deleted without any test noticing are now each held by a named test, and every test says in its own comment which deletion it catches. The remedy redirection is pinned by reading a finding's rewritten remedy rather than the record-level continuation the old assertions read, at both call sites, with three negative controls proving the rewrite is selective: a sibling finding whose remedy stays null, a git remedy left untouched because its first argv token is not the CLI, and one finding where one field is rewritten and the other is not. The layered guard in focusedGateState is pinned by a nine-row in-process table keyed on the launched, timed-out and baseline-state ingredients rather than on the exit statuses that happen to produce them today; four of those rows are controls showing the guard admits an eligible execution instead of refusing everything. The interrupted focused test gets a public case that signals the CLI from inside the probe, paired with an ordinary red run at the same exit status so the only variable between them is whether the runner observed an interruption.

Behaviour is unchanged: the three product behaviours were mutated temporarily to prove each test red, then restored and verified byte-identical by digest. The merge range `a22ddfcf..722f5ada` contains one file and 207 insertions with no deletions, which matches the builder's manifest exactly. Builder branch head `0e04630a`.

## Qualification

**Passed.** Read from the merge range `a22ddfcf..722f5ada`, and the central claim re-verified independently rather than taken from the hand-back.

- **The claim this REQ makes is exactly the claim a hand-back cannot settle by itself.** The REQ exists because three behaviours could be deleted with the package staying green, so "I mutated it and the test failed" is the one assertion that has to be reproduced. The orchestrator cut a detached worktree at the merge commit, applied all five mutations one at a time, and ran the focused package for each. Every one reddened, each on the named test the builder claimed, and the tree was green before the first mutation and green again after the last. Nothing in the qualification below rests on the builder's own report of that.
- **Substantive, and the diff matches the manifest.** One file, 207 insertions, 0 deletions — the same file the Scope declared and the same numbers the hand-back reported. `internal/lifecycleadvance/evidence_gates.go` and `internal/corehelpers/commands.go` are byte-identical to their pre-work state, which is what the REQ's tests-only constraint requires; the mutations were applied and reverted, never landed.
- **The tests are characterization, not smoke.** Each of the three carries a comment naming what it pins and which deletion it catches — an explicit acceptance criterion, and the guard against the `[family: smoke-vs-characterization]` trap the prime names for this area. The rows key on the execution facts the code reads (launched, timed-out, baseline state) rather than on the exit statuses that produce them today, which is the `[family: closed-enumeration-for-a-condition]` rule applied where it belongs.
- **The negative controls are real controls.** Three of them fail if the redirection becomes a blanket rewrite: a sibling finding whose remedy must stay null, a git remedy whose owner's argv must survive, and one finding where one field is rewritten and the other is not. The focusedGateState table adds four rows proving the guard still admits an eligible execution instead of refusing everything. And both interrupted-probe rows exit 143, so the classification cannot pass by keying on the integer.
- **The builder found the REQ's own premise incomplete, and said so.** The REQ treats M1's two call sites as symmetric. They are not: with only the composeGreenGate site deleted the package stayed green even with the natural-producer test in place, because no realistic canonical gate argv produces a finding whose remedy re-enters that helper. That is recorded as D-01 with its value and risk, and as two report-only discovered tasks. Reproduced here: M1b reddens only through the constructed subtest.
- **One warning was corrected rather than argued with.** Scope-drift first reported `focusedGateState` as a declared-but-untouched path. It is a Go identifier that the orchestrator had written in backticks inside the Scope bullet, and the checker reads backticked tokens there as paths. The prose was fixed; no code or declaration changed.

Requirements traced, one by one: a finding's rewritten remedy is asserted at both call sites; both halves of the guard are covered independently, and a third half the REQ did not name is covered too; a public case asserts `BLOCKED-PROBE-FAILED` for an interrupted focused test; every new test states what it pins and which deletion it catches.

*Checked by work action*

## Testing

**Tests run:** `go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/lifecycleadvance`
**Result:** ✓ All passing, 20.99s, run through the focused-test gate against the pre-build baseline, which was green — so this is a green-on-green comparison, not an unusable baseline.

**Canonical repository gate:** `bash _dev/tests/maintainer-verify.sh` exited 0 at revision `71eb49f3`, run directly and unpiped from the project root, wall 110.8s. No retry was needed; the gate has not exited non-zero once during this REQ. Three runs of the same argv over this REQ's life all exited 0: `93f856ee` before the claim, `9d58ae17` for the pre-build baseline, `71eb49f3` after the merge.

**A note on the revision the final gate ran at.** A sibling session in this same checkout completed REQ-588 and released while this REQ was between merge and gate, so `HEAD` moved from this REQ's merge commit `722f5ada` to `71eb49f3`. The gate was started and recorded at `71eb49f3` and matches it exactly. This REQ's merge commit is an ancestor of that revision, its merge range is unchanged, and nothing the sibling landed touches the package under test.

**Red-green validation** (traced to `## Red-Green Proof`, whose RED case is three named mutations), reproduced independently by the orchestrator in a detached worktree cut at the merge commit — not read from the hand-back:

- Unmutated tree: `go test ./internal/lifecycleadvance` exit 0.
- M1a, deleting the `redirectHelperRemedies` call in `composeCoreGate`: ✗ `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation/core_gate_call_site`.
- M1b, deleting the call in `composeGreenGate`: ✗ `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation/green_gate_call_site`. Before this REQ, and even with the natural-producer test alone, this deletion left the package green — see D-01.
- M2a, deleting `subordinateState == resultmodel.AdvanceGateFailed ||`: ✗ two rows of `TestFocusedGateStateKeepsSubordinateAuthority`, both failed-subordinate rows.
- M2b, deleting `|| focusedTest.TimedOut`: ✗ both timed-out rows of the same test.
- M2c, deleting `!focusedTest.Launched ||`: ✗ the never-launched row. This half is not named in the REQ; the builder added it because the guard has three parts, not two.
- M3, restoring `BLOCKED-PROBE-LAUNCH-FAILED` on the interrupted path: ✗ `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure/interrupted_run`.
- Tree restored: exit 0.

Every mutation was applied alone, and the ordinary-red control row stayed green under all of them, which is what proves it is a control and not a duplicate of the interrupted row.

**New tests added:**
- `TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation` — two subtests, one per call site, plus three negative controls.
- `TestFocusedGateStateKeepsSubordinateAuthority` — nine rows: five that catch a deleted guard half, four that prove the guard still admits an eligible execution.
- `TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure` — an interrupted run and an ordinary red run at the same exit status.
- Helper `advanceGateFinding`, which returns the finding under a code so a test can read its remedy. `gateHasFinding` was left alone rather than widened, because four existing tests use it (D-05).

**Existing tests updated (cross-REQ impact):** None. No existing test, helper or constant was edited or reordered.

**Timing against the per-file budget:** `evidence_gates_test.go` rose from 6.27s to 6.99s against the 30s per-file budget the gate enforces. `-count=8` over the three new tests found no flake.

**Heavy verification plan:** *(planned, held for the drain at queue exhaustion)*
- Range: `a22ddfcf`..`722f5ada`
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the changed file matched subtree `skills/do-work/tools/do-work-cli`
- `staged-skills`: same argv with `--heavy-lane staged-skills` — matched subtree `skills`
- `updater`: same argv with `--heavy-lane updater` — matched subtree `skills/do-work/tools/do-work-cli`
- `installer`: same argv with `--heavy-lane installer` — matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Decisions

Authored by the builder in `do-work/runs/work-2026-09-05-170806/REQ-583-handback.md` → `## Decisions`; transcribed here by the orchestrator, which is the only writer of this file in worktree dispatch mode.

- **D-01 — ESCALATE. The composeGreenGate call site has no natural producer, so its test supplies the redirection condition through the gate argv.** Deleting only that call leaves the package green, and it still does with M1's natural-producer test in place. The reason is structural: of the findings the two green-gate helpers can return, two carry a git remedy, one carries the caller's own gate argv, and the only one whose remedy re-enters the helper cannot be produced from advance because composeGreenGate always builds a well-formed argv for it. The choice was between leaving that site unpinned — which fails the REQ's "deleting **either** call site fails a named test" — and supplying the condition the redirection keys on through the one channel that can carry it: the gate argv, which advance copies verbatim into the not-green remedy and never executes. The builder took the second. **Value:** the acceptance criterion is met for both call sites, and the pin is on the stated condition rather than on the spellings today's callers happen to produce. **Risk:** the fixture's gate argv is not a realistic repository gate, so a reader could mistake the subtest for a contrived assertion; a comment above it states why no natural producer exists. Fully reversible — deleting the subtest returns the tree to one site pinned and one not, which is where the evidence says it stands.
- **D-02 — DECIDE & STATE. A third M2 row was added for the `!focusedTest.Launched` half.** The REQ names two deletions; the guard has three parts. The never-launched half is masked today by the failed half, so deleting it alone was green. The row costs microseconds and completes the condition.
- **D-03 — DECIDE & STATE. The green-gate negative control asserts the first argv token is git, not the exact git argv.** Asserting the argv byte-for-byte would redden this test whenever `internal/gateevidence` reworded its own remedy, which has nothing to do with the redirection. The first token states the contract the redirection actually has and still catches a blanket rewrite.
- **D-04 — DECIDE & STATE. Both interrupted-probe rows exit 143.** A different status on the control row would let the test pass while the classification keyed on the integer.
- **D-05 — DECIDE & STATE. `advanceGateFinding` was added rather than widening `gateHasFinding`.** The existing helper answers a different question and four tests use it; changing its signature would touch tests outside the write set.
- **D-06 — DECIDE & STATE. The focusedGateState table calls product code in process; no other new test in this package does.** It is stated in the test's own comment: the guard is defence in depth, its only current caller cannot produce the input, and a public test would need a second producer of the focused-test result type.

## Review

**Overall: 96%** | 2026-09-05T17:52:43Z

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

## Lessons Learned

**What worked:** Keeping each mutation applied while writing the test that catches it, one mutation at a time rather than all three at once. That is what turned "the package is green under this deletion" from a claim into a measurement, and it is what surfaced the thing the REQ itself had wrong — that M1's two call sites are not symmetric. It also produced honest per-half differentials for the guard, which a single combined revert could never have given.

**What didn't:** The REQ's premise that pinning the remedy redirection means writing one test. Deleting only the `composeGreenGate` call site left the whole package green even with the natural-producer test in place, because nothing a realistic canonical gate argv produces reaches that rewrite: two of the green-gate findings carry a git remedy, one carries the caller's own gate argv, and the only one whose remedy re-enters the helper cannot be produced because that function always builds a well-formed argv for it. The second pin had to supply the condition through the gate argv itself. The reviewer audited that from the source rather than the hand-back and found no producer the builder had missed, and then showed the subtest earns more than the acceptance criterion: it is the only assertion in the package that pins the verb discrimination, so dropping the verb check and matching on the tool name alone — a real regression that hijacks another helper's remedy — is caught by that one line and nothing else.

**Worth knowing:** Two of the three behaviours here are *inert*, not dead. The `focusedGateState` guard runs on every focused-gate invocation, and no input its single caller can produce makes it change the answer, because an upstream eligibility check has already left those executions at a baseline state the switch below the guard has no case for. Deleting it is behaviour-preserving today. That is a genuinely awkward thing to pin, and the honest form is what was done: a direct in-process table with a comment saying the guard is defence in depth, its only caller cannot produce the input, and a public test would need a second producer. The cost is real and should be stated rather than hidden — a future simplifier who correctly identifies the guard as inert will now redden a named test instead of noticing it silently.

## Orientation

The evidence-gate layer of do-work-cli now has tests that fail when its three quietest behaviours are removed: the redirection that sends a subordinate command's remedy back through the same request-bound advance invocation, the guard that stops a saved baseline clearing an execution it was never eligible to be compared against, and the classification that separates an interrupted focused test from one that merely exited non-zero. Lives in the do-work-cli lifecycle-advance subsystem, governed by `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`. No behaviour changed and no map changed — this closes an independent-review finding that three pieces of REQ-506's work had no lock. The prime and its lesson satellite were spot-checked for staleness against this change and both are current; the one stale restatement found anywhere is in `skills/do-work/actions/work.md`, whose list of focused-test gate outcomes omits the interrupted class, reported by the review and left as report-only.

## Heavy Verification Plan

- **Base revision:** `a22ddfcf7f16207433e4d399a37b9442ebdbefdc`
- **Target revision:** `722f5ada02df491c9b44e22d6815dc328dfa63ec` (the recorded `commit:`)
- **Changed paths in range:** `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go` — no uncovered paths, planner not forced, not uncertain.

| Lane | Argv | Selection reason |
| --- | --- | --- |
| `do-work-cli-integrations` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` | matched subtree `skills/do-work/tools/do-work-cli` |
| `staged-skills` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` | matched subtree `skills` |
| `updater` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` | matched subtree `skills/do-work/tools/do-work-cli` |
| `installer` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` | matched subtree `skills/do-work/tools/do-work-cli` |

Held at Step 7.7: the lanes are not run now and the queue loop is not held open. The recorded `commit:` above is what makes this REQ's source ready for anything that depends on it while it waits for the drain at queue exhaustion.
