# Review: REQ-559

**Approve with follow-ups** — the one-shot retry is real, correctly bounded, and records only the second run; one sentence in the already-green repair section still says that branch gets a single gate run, which the retry can now make false.
Route B | merge range `904b4d3..0e07aac` (merge `0e07aac`, work commit `bce5fb4`)

### What's built

- A red baseline test command in `tools/checks/preflight.sh` is now rerun once, immediately, with the same argv. Only the second run reaches `do-work/working/baseline.json` and `baseline-failures.txt`, and one WARN line names both exit statuses. Verified end-to-end, not just in unit tests.
- One condition-keyed prose rule, `actions/work-reference.md` → **One retry before classification**, is cited from `work.md` Step 5 (pre-flight), Step 5.75 (repository-gate baseline) and Step 6.5 item 4 (post-merge attribution). Every downstream branch reads the second run's status, output and fingerprint.
- Still missing: the already-green repair section was not swept for the sentence that now conflicts with the rule (finding I1 below).

### Decisions / risks for you

- **D-01 (retry in the Go handler, not `preflight.sh`) is correct.** `skills/do-work/tools/checks/preflight.sh` is a six-line launcher that `exec`s `do-work-cli.sh --format text preflight`; it contains no command execution to change. The behavior contract the REQ names at that entry point is genuinely satisfied — running `bash skills/do-work/tools/checks/preflight.sh bash <flaky>` in a scratch repository produced the retry, the single WARN line, and `exit_status: 0` in `baseline.json`. Editing the launcher would have meant duplicating the runner in shell.
- **D-02 (rule keyed on the condition, not on the Step 6.5 site alone) is faithful, not scope creep.** The REQ's `## What` says "at the baseline (Step 5 pre-flight) or after integration (Step 6.5)", and the incident in `## Why` (REQ-531 deferred, REQ-548 minted) happened in the Step 5.75 repository-gate baseline lane, which is the only lane that calls `defer-gate` before source edits. A rule keyed only at Step 6.5 would not have prevented the incident the REQ exists to prevent. The condition form also matches CLAUDE.md § "State conditions, not lists".
- **Wall-clock cost in this repository.** `_dev/tests/maintainer-verify.sh` is both the Step 5 pre-flight test command and the Step 5.75 canonical gate here, so a genuinely red suite can now cost four full gate runs before a deferral instead of two. The builder disclosed this in `## Discovered Tasks`. Bounded and correct, but it is the price of a genuinely red gate on a slow machine.
- **The retry does not weaken failure attribution.** Recording only the second run means a failure seen only on the first run is *not* written into `baseline-failures.txt`, so if it reappears at Step 6.5 it is treated as a candidate regression rather than silently excluded. That is the safe direction.

### Findings

**Important:**
- `skills/do-work/actions/work-reference.md:442` — "**The pre-build run is this branch's only gate run.**" and, in the same paragraph, "Done means one gate run plus bookkeeping — under ten minutes of wall clock". A repository-gate repair whose pre-build gate exits non-zero and then green on its one retry reaches the already-green no-op path after two gate runs. That is exactly the REQ-548 shape this REQ targets, so the sentence is stale in the case that matters most, and an agent reading it as a hard budget could conclude the retry does not apply to a repair REQ's pre-build gate. Detailed Requirement 5 ("Delete any sentence that would now contradict the rule") is therefore only partly met — `impact-rule-change` → report only

**Minor:**
- The rule's reach over the Step 6.5 *diagnostic worktree* base run at `<pre>` is unstated. `work-reference.md:372` says the retry applies "in every lane that launches the gate directly — the baseline lane and late attribution alike", and the base run at `<pre>` is a direct run of the canonical gate argv inside late attribution; `work.md:410` reads narrower ("everything after this sentence ... reads the second run's status"), implying only the current-tree run is retried. Either reading is safe, but they cost different wall clock and the REQ's own wording ("A red rerun enters the existing path unchanged: fingerprint, diagnostic worktree, defer-gate") suggests the base run was not meant to be a retry site — `impact-rule-change` → report only
- Because the same script serves Step 5 pre-flight and Step 5.75 in this repository, a genuinely red suite now pays two extra full runs, not the "one extra run" the new prose at `work-reference.md:374` promises — `impact-user-visible` → report only
- Test gap: no probe pins that a *green* first run is invoked exactly once (the "no retry on zero exit" half of the bound), and none covers a retry whose second run cannot launch. I exercised both by hand: a green command runs once, and `/definitely/not/here` produces `exit_status: 127`, `launched: false`, plus both the retry WARN and the not-launched WARN — `impact-negligible` → report only
- `skills/do-work/docs/command-line-guide.md:46` still describes the `defer-gate` manifest as binding the "direct non-zero gate result". Still true, but that result is now the second run's; the builder already recorded this in `## Discovered Tasks` — `impact-negligible` → report only

**Nit:**
- A command that cannot launch at all (exit 126/127) is still retried and still prints the retry WARN line beside "could not run the test command". D-03 accepts this deliberately and the second attempt costs nothing measurable, but the extra line is noise on a path where nothing could ever have been transient — `impact-negligible` → report only

### Requirements Checklist

- [x] Baseline: `preflight.sh` reruns the gate command once when the first run exits non-zero — delivered (`checks.go:88-93`, verified end-to-end through the launcher)
- [x] Records only the second result in `baseline.json` — delivered; `baselineRecord` is written after the retry, and `baseline-failures.txt` holds the second run's output (`checks.go:95-110`), pinned by both new probes
- [x] Prints one WARN line naming the retry and both exit codes — delivered (`PREFLIGHT-BASELINE-RETRIED` plus its `preflightCompatibilityText` case); the probe asserts the line appears exactly once
- [x] Post-merge: one rule in `work-reference.md`, cited from Step 6.5 item 4 — delivered (`work-reference.md:370-374`, cited at `work.md:410`, and additionally at `work.md:273` and `work.md:279`)
- [x] Zero exit on the rerun records green and continues; non-zero continues with the existing diagnostic and defer procedure using the second run's output as the fingerprint source — delivered (`work-reference.md:380`, `:387`, `:404`)
- [x] Progress output shows the retry as one line; REQ Testing section records both exit codes — delivered (Testing Section Template gained an optional `**Repository gate retry:**` line at `work-reference.md:797`; this REQ's own `## Testing` uses it)
- [x] No new status, no new flag, no new REQ type; `defer-gate`, `repository_gate_repair` and the already-green path untouched — delivered; `baseline.json` keeps its three fields (D-04), no CLI flag added, no code outside `handlePreflight` changed
- [~] Delete any sentence that would now contradict the rule — partially delivered; the branch-table row, the baseline paragraph and the late-attribution paragraph were all rewritten, but `work-reference.md:442` still asserts a single gate run for the already-green repair branch (finding I1)
- [x] Constraint: mechanics in the program, judgment in prose; no new prose walking a shell sequence — met; the new rule states a condition and an outcome, not a command sequence
- [x] Constraint: no `do-work/` file touched, `_dev/tests/maintainer-verify.sh` not edited — met; the merge range touches four files, none of them under `do-work/` or `_dev/`

Requirements Compliance: **90%** (nine of ten fully delivered, one partial).

### Restatement Sweep

The diff redefines what "the direct non-zero gate result" means (it is now the second run's) and adds one finding code plus one Testing-template line. Swept:

- Every shipped statement of the gate's non-zero branch — `work.md:273`, `:279`, `:410`; `work-reference.md:372`, `:380`, `:387`, `:396`, `:404`, `:408` — all agree with the new meaning.
- `work-reference.md:442` (already-green repair, "only gate run") — **stale**, finding I1.
- `docs/command-line-guide.md:46` — technically still true, cosmetically stale (Minor above).
- Finding-code consumers: `commands.go:361-402` `findingSpecificCommands` needs no case because `PREFLIGHT-BASELINE-RETRIED` supplies its own `next`/`verify` argv; `preflightCompatibilityText` has the new case; no other file enumerates preflight codes.
- Testing-section template: `work-reference.md:794-800` is the only canonical template; `work.md:434-442` cites it and needs no change; `actions/sample-archived-req.md` is an example of a non-retried run and correctly omits the optional line.
- Baseline artifact consumers: `work.md:415` (Step 6.5 item 5) reads `launched` and the failure list — both semantics unchanged, only the values now come from the second run.

### Acceptance Testing

**Result: Pass**

- `go vet ./...` clean, `gofmt -l .` empty.
- `go test ./internal/corehelpers/ -run Preflight -count=1` — 3 tests pass.
- RED reproduced independently: copied the module to a scratch directory, restored `checks.go` from `904b4d3`, kept the new probes. Both fail with `baseline command ran 1 times, want exactly 2`. The builder's red-green claim is real.
- End-to-end through the shipped launcher, in a scratch git repository (the real working tree was never used as the target):
  - transient red (exit 7 then 0) → `WARN: baseline command exited 7 on its first run and was rerun once; the rerun exited 0 and is the only recorded result`, then `OK: test baseline passing`; `baseline.json` holds `exit_status: 0`; no `baseline-failures.txt` left behind.
  - persistent red (always exit 5) → the same retry WARN naming `5` and `5`, followed by the existing `WARN: baseline tests failing BEFORE any changes`; `baseline.json` holds `exit_status: 5`; `baseline-failures.txt` holds `red run 2`, the second run's output. The existing red path is reached unchanged.
  - command that cannot launch → `exit_status: 127`, `launched: false`, both the retry WARN and `WARN: could not run the test command`.
  - The retry is bounded: the invocation counter reads exactly 2 in every red case, never 3.
- Gate checks relevant to the changed areas: `_dev/tests/action-shell-blocks.sh` pass, `_dev/tests/shipped-package-reference-contract.sh` pass, `_dev/tests/contract-regressions.sh` fails only on the pre-existing `session-start-hook-behavior.sh took 32s; each test file must finish under 30s` budget miss — every assertion inside that file passes, and it is the failure already recorded in `do-work/working/baseline-failures.txt` before this REQ. Not attributable to this change.
- Not run: the full `_dev/tests/maintainer-verify.sh` (already run by the orchestrator; excluded by the review brief).

### Suggested Additional Testing

- Add a probe that a zero-exit baseline command is invoked exactly once, so the "no retry when green" half of the bound is pinned mechanically rather than by reading the `if status != 0` guard.
- Add a probe for a retry whose second run cannot launch (first run exits non-zero, second run's binary is missing), so `launched: false` after a retry stays pinned.
- Watch the first real transient gate failure in the wild and confirm the progress output shows exactly one retry line and the REQ's `## Testing` carries both statuses through the new template slot — the agent-driven Step 5.75 and Step 6.5 lanes have no mechanical enforcement, only prose.
- Time one genuinely red run end to end in this repository to see whether four full `maintainer-verify.sh` runs before a deferral is acceptable, or whether the pre-flight result should be reused for the Step 5.75 baseline lane.

### Scores (on the record — not the headline)

**Overall: 90%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 90% | Nine of ten delivered; the contradicting-sentence sweep missed `work-reference.md:442` |
| Code Quality | 92% | `runBaselineCommand` is a faithful extraction, correctly named, commented with the reason not the mechanism; vet and gofmt clean; no debug artifacts |
| Test Adequacy | 85% | Two meaningful behavior probes with independently reproduced RED; missing the green-runs-once and retry-cannot-launch cases |
| Scope | 95% | Four files, all declared in `## Scope`, all six decisions recorded; the `write_set`/`## Scope` divergence is explained by D-01 rather than silently taken |
| Risk | Low | Bounded at one extra run; only the pre-flight lane is mechanical, the rest is prose an agent must follow |
| Acceptance | Pass | Retry, boundedness, second-result-only recording, and the unchanged red path all verified by running the shipped launcher |

### Follow-ups created

None (6 findings report only)

---

**Verdict: Pass**

Acceptance is Pass and the overall score is 90%, above the 75% Approve band. The retry is genuinely bounded at one — the invocation counter reads exactly 2 in every red case in both the unit probes and the end-to-end runs — only the second result is recorded, and a second failure reaches the existing fingerprint and `defer-gate` path with no code or contract change. The two builder decisions under scrutiny both hold: D-01 put the mechanics at the only site that can execute the command while leaving the named entry point's behavior contract satisfied, and D-02's condition-keyed rule covers the lane that actually caused the incident this REQ exists to remove. The one requirement not fully met — deleting every now-contradicting sentence — leaves a single stale sentence in the already-green repair section; it misstates a run count, does not change any mechanical behavior, and is recorded above as `impact-rule-change` report-only rather than a blocker.
