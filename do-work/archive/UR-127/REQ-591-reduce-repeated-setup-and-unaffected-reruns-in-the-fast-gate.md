---
id: REQ-591
title: 'Reduce repeated setup and unaffected reruns in the fast gate'
status: completed
created_at: 2026-09-05T19:43:25Z
user_request: UR-127
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: []
related: [REQ-574]
claimed_at: 2026-09-05T20:04:13Z
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-05T20:06:28Z
  basis:
    - Route C
    - 6-file write set
    - 1 new files
    - 2 subsystems involved
    - 4 acceptance criteria
    - performance instrumentation
    - cross-route regression gates
    - full-suite verification
write_set:
  - _dev/tests/session-start-hook-behavior.sh
  - skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go
  - _dev/tests/fast-stages.json
  - _dev/tests/maintainer-verify.sh
  - _dev/tests/fast-stage-reuse-behavior.sh
  - _dev/tests/contracts/probe-lanes.sh
route: C
dispatch_at: 2026-09-05T20:35:46Z
builder_handback_at: 2026-09-05T21:44:19Z
remediation_at: 2026-09-05T22:20:13Z
review_at: 2026-09-05T22:45:16Z
integration_at: 2026-09-05T21:44:19Z
planning_at: 2026-09-05T20:34:44Z
commit: 63f5288d5c1cf8a65842280a49248f65c11fe5d9
completed_at: 2026-09-06T03:36:21Z
release_at: 2026-09-06T03:36:21Z
---

# Reduce Repeated Setup and Unaffected Reruns in the Fast Gate

## What

Make this repository's routine verification faster and cheaper without losing the failure coverage that protects queue state and completed work. Start by removing repeated fixture/build setup, then reuse verification results only when the complete inputs of the relevant checks are unchanged.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The Plan agent produced the approach from a measured exploration: the fourteen serial gate stages with their costs, the root cause of the slowest probe, and the real invalidation contracts of the two existing evidence mechanisms. The build followed its task order and corrected three of its assumptions from evidence rather than from preference (D-12, D-13, D-14).
- [x] **[APPLY]:** Five commits, eight files, every one inside the declared scope. The one declared file that was not touched is struck through above with its reason.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file over the cumulative merge range; the full gate green with the reuse switch unset and with it forced off; module tests, `go vet`, `gofmt -l` and ShellCheck clean on every shell file touched or added; and eleven mutations run against the new code — seven on the Go engine and four on the gate wrapper — to prove the negative cases fail rather than decorate. The orchestrator additionally re-read the merged range, took the whole-gate four-condition comparison itself, and verified the three new commands reach the shipped CLI's own help output.

*Authored across the builder's and the remediation builder's hand-back in `do-work/runs/work-2026-09-05-170806/REQ-591-handback.md`; transcribed and checked here by the orchestrator, which is the only writer of this file in worktree dispatch mode.*

## Why

The user wants the queue to drain faster while maintaining good enough quality. The CPU investigation found several independent test pools competing for one machine. Avoiding repeated work lowers the verification cost per request; increasing concurrency alone does not establish an improvement.

## Context and Prior Work

- REQ-574 (reducing CLI fixture setup costs) is already completed, merged at `50569e88c8f1f5234cbdfaf0efaede671d72b13c`. It reused initialized Git repositories in finalization/publication tests and avoided repeated launcher/toolchain work in inventory cases. Its recorded whole-module comparison was 65s to 61s with 772 tests and every assertion retained. Build on that work instead of repeating it. Source: `do-work/archive/UR-115/REQ-574-bring-do-work-cli-test-files-under-the-30s-budget.md`.
- `_dev/tests/session-start-hook-behavior.sh` still copies the complete CLI module separately for each scenario. This is a concrete setup-reuse candidate, not a requirement to bypass the launcher path under test.
- `_dev/tests/maintainer-verify.sh` selects the ordinary board tests and the complete CLI module with `-short`; `_dev/tests/run-go-tests-with-budget.sh` forces `-count=1`. Existing fast/heavy separation is already implemented.
- Existing heavy-lane fingerprints and green-gate evidence provide starting points. Inspect their actual invalidation and caller contracts before extending them; a prior green at HEAD is not proof that current dirty or external inputs match.
- `do-work/test-durations.tsv` provides per-file observations. Summed file times are not complete-gate wall time or CPU time, and loaded-window observations cannot establish an uncontended speedup.
- Fold-first review found no eligible pending request in any UR with this root cause. The queued helper deduplication and shell-guide requests do not close the repeated fast-gate setup/rerun problem.

## Detailed Requirements

1. Reduce repeated setup in measured hot paths. Reuse immutable fixture material and the current build where that preserves the boundary being tested. Keep writable test state isolated. Retain explicit coverage of the real launcher, build failures, missing tools and installed layout where those are the behavior under test.
2. After the setup improvement, avoid repeating verification whose complete relevant inputs have not changed. Use the existing verification/evidence mechanisms where practical. Relevant inputs include source, transitive dependencies, fixtures, test/gate scripts, configuration and effective toolchain/runtime inputs, including uncommitted changes. Unknown impact, incomplete evidence or an unverifiable input must select the broader verification rather than produce a false green.
3. Make selective results reviewable: record which checks executed, which reused evidence and why. Preserve failure and interruption exit statuses. A skipped, failed or incomplete run cannot supply successful evidence. A board-only change may reuse unrelated CLI evidence only after proving that the relevant CLI inputs are unchanged.
4. Keep useful failure coverage. Preserve real checks for rollback, interrupted recovery, locking, concurrent writers and process cleanup. Do not silently remove assertions, replace a real boundary with a mock, or move essential coverage behind the heavy tier to improve the number. Any proposed consolidation must name the failure it still catches and its retained test.
5. Measure before and after on comparable revisions with a fixed worker limit, recorded toolchain/cache state and no competing expensive gate or synthetic load. Record complete-gate wall time and process-tree CPU cost (or explicitly identify unavailable measurements), alongside the existing per-file durations. Use the smallest bounded comparison that establishes a repeatable improvement; do not benchmark by repeatedly saturating the machine.

## Acceptance Criteria

- Repeated setup is measurably reduced in at least one identified hot path, with its behavior assertions retained and writable fixtures independent.
- A matching-input verification case reuses a recorded success; changing a relevant source/dependency, fixture, gate script or runtime input invalidates it. Uncertain impact and dirty relevant inputs cannot reuse stale evidence. Tests cover these positive and negative cases.
- The routine gate's end-to-end duration improves under the recorded comparable conditions; the evidence separates setup savings, avoided execution and contention. A one-off noisy sample, a renamed/split test file, a raised timeout or greater concurrency alone is not acceptance.
- The existing fast-tier per-file budget and heavy-tier policy remain satisfied. Required correctness checks pass; retained rollback/recovery/locking/cleanup tests remain capable of detecting their named failures. Reused results are visibly distinguishable from freshly executed ones.

## Constraints

- Repository scope is `skill-do-work2`. The other projects seen in process metadata are diagnostic context, not authorization to inspect or change their test configurations.
- No machine-wide scheduler, background reaper or watchdog in this request. Shared machine concurrency policy is adjacent work; use a controlled worker limit for the comparison here.
- Preserve existing heavy-test permission/defer behavior. This capture is not authorization to run an otherwise permission-gated heavy suite.
- Synthetic load is unnecessary for this optimization. Any deliberately backgrounded process still needs its own lifetime bound under the existing testing rule.
- Do not optimize to a fabricated percentage or treat the 30-second per-file ceiling as a whole-gate speed target.

## Builder Guidance

Implement and verify setup reuse before selective rerun behavior. Choose the smallest existing code path that supports conservative invalidation; do not introduce a second verification platform. Keep one coherent request, with small verified increments. Exact files and test cases are resolved during planning rather than guessed into `write_set` at capture.

## Red-Green Proof

**RED prompt/case:** Exercise the normal fast gate's current repetition on unchanged relevant inputs, then a fixture where only unrelated board inputs change. Record which expensive checks rerun. In the existing harness, add a test of the intended reuse decision plus counter-cases for relevant source, fixture, gate script, runtime and uncommitted-input changes.
**Why RED now:** Routine verification still executes broad uncached selections and repeats fixture setup; the proposed fast-gate input-aware reuse is not established by the existing heavy-lane or revision-only evidence rules.
**GREEN when:** The matching/unrelated-input cases reuse only valid success, every relevant or uncertain change executes the needed checks, retained failure scenarios still fail when broken, and comparable before/after evidence shows lower routine gate cost.
**Validation:** Inferred during capture from the recommendation the user asked to capture. The user approved the objective; exact regression fixtures and a defensible performance baseline are builder-resolved details.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — index cost 10543 tokens; exceeds the 2000-token budget and is `slugged: partial`, so targeted selection is not eligible. Matches `fixture-cost-is-subprocess-spawning`, `background-worker-self-bound`, `smoke-vs-characterization` and `shipped-module-test-self-containment`.
- `_dev/primes/lessons-shell-commands.md` — index cost 3385 tokens; exceeds the budget and is `slugged: partial`. Matches the shell/launcher fixture and migration-parity surface.

## Open Questions

None about user intent. The builder establishes the exact measured scope and comparison conditions.

## Full Context

See `do-work/user-requests/UR-127/input.md` for the capture instruction and prior conversation context.

---

## Triage

**Route: C** - Complex

**Reasoning:** Two separable pieces of work, and the second one is a new decision boundary rather than a change to an existing one. Reducing repeated fixture setup is bounded and measurable. Deciding when a verification stage may reuse a recorded success is a correctness question about invalidation: the request itself says an unknown impact, incomplete evidence or an unverifiable input must select the broader verification rather than produce a false green, and that is exactly the rule that has to be designed before it is written. It also has to be built on the evidence mechanisms that already exist rather than a second platform, which means reading their real invalidation contracts first. A written plan decides what is in scope before any of it is typed.

**Planning:** Required.

**A claim-time note on the measurement condition.** This request's acceptance asks for a before/after comparison with no competing expensive gate and no synthetic load, and forbids a one-off noisy sample. This checkout is shared with several concurrent sessions and the load average has ranged from 5 to 59 over the last two hours; five canonical gate runs during this work run failed on per-file wall-clock budgets purely from that contention, on four different files, while the same files pass in a quiet window. The comparison this request needs is therefore obtainable but only inside a measured quiet window, and every recorded number has to name the load it was taken under. That constraint is stated here so a later reader can see it was priced rather than discovered.

## Plan

*Generated by Plan agent. The full working plan, with every code block and the complete measurement protocol, is in the run directory as `do-work/runs/work-2026-09-05-170806/REQ-591-plan.md`; the exploration it rests on is beside it as `REQ-591-exploration.md`. This section is the durable record.*

**Scope judgment: both parts, sequenced, five tasks.** Part two is the whole of acceptance criterion 2 and most of criteria 3 and 4, so splitting it out would leave this request unable to meet its own acceptance — a capture-level change rather than a builder decision. Five tasks is the signal line, not past it, and the clean seam if smaller increments are ever wanted is inside T2, between the evidence engine and the commands that expose it.

**Tasks.** T1 share one skill root in `_dev/tests/session-start-hook-behavior.sh`; T2 the fast-stage evidence engine and its three commands, tests first; T3 the `_dev/tests/fast-stages.json` manifest; T4 the gate wiring, the reuse summary line, the self-test bypass guard and the maintainer-tree pinning; T5 the end-to-end reuse probe, its registration, and the measurement.

**Decisions.**

- **D-01 — share one physical skill root rather than symlink nine of them. DECIDE & STATE.** The symlink shape's win rests on Go canonicalizing a path component the shell never canonicalized; a shared root reaches one physical module by construction and needs no such assumption.
- **D-02 — the banner input becomes a required parameter of every case helper. DECIDE & STATE.** With a shared root, "forgot to set the banner input" would be an ordering bug that reads as a pass; a required argument makes it a syntax error instead.
- **D-03 — extend `internal/heavyverification` rather than add a package. DECIDE & STATE.** The coverage rules, manifest decoding, evidence-store mechanics, disposition constants, four-hour ceiling and fingerprint probe runner are all unexported there; a new package would duplicate them or force a wide export surface, and the request forbids a second verification platform. The package name then under-describes its contents, which is a separate mechanical rename and not done here.
- **D-04 — the fast-stage fingerprint seals working-tree bytes, not committed object ids. DECIDE & STATE.** The heavy lane can seal committed objects only because it first refuses a dirty tree. The fast gate exists to run on a dirty tree, so committed objects would be a false green the moment anything is uncommitted, which requirement 2 names explicitly.
- **D-05 — fast-stage records live in their own key space with their own schema version. DECIDE & STATE.** A fast green is computed on a possibly-dirty tree and is not attributable to a revision; if the two key spaces could cross, a dirty-tree fast green could authorize a heavy lane, which is a weaker guarantee standing in for a stronger one.
- **D-06 — exactly the two Go test stages are reusable. DECIDE & STATE.** They are 73s of a ~102s gate and about 85% of its CPU. The rest are excluded with reasons: `go vet` is already content-addressed and 0.2s warm, so caching it caches a cache; gofmt is 0.07s and the fingerprint would cost more than the check; the version floors *are* the toolchain check; the wide reference probes resolve link targets against arbitrary repository paths, so anything narrower than a whole-tree seal is a false green, and the narrow probes are already 0.02-1.4s. ShellCheck lint at 3.7s over a closed 92-file input set is a legitimate future candidate, named as a follow-up rather than taken here.
- **D-07 — three named environment variables are removed from the decision child's environment, in the shell, with the reason at the call site. DECIDE & STATE.** The gate exports a run label, a `ps` sample and an output path before any tier runs; all three change between runs, and the fingerprint seals the whole environment, so with them present reuse could never fire. Excluding them in the shell keeps the Go rule "seal every variable, no exceptions" absolute and makes the exclusion one greppable line. The two variables that change the verdict are deliberately not excluded.
- **D-08 — a reused stage inherits its per-file budget verdict along with its pass/fail verdict, and the reuse line says so. ESCALATE.** A reused stage writes no duration rows and enforces no per-file budget for that run. The budget exists to catch tests getting slower as they are added, which is an input-determined property the fingerprint covers; it does not cover a breach caused purely by machine contention — and five gate runs during this very work run failed on exactly that. **Value:** the gate stops failing on other sessions' load for a stage whose inputs are provably unchanged. **Risk:** a real contention problem goes unreported for up to four hours on a matching tree. Reversible by disabling reuse for that stage in the manifest. The recommendation is to accept and to print the inherited verdict, naming the run whose budget verdict is standing in.
- **D-09 — the record command recomputes the fingerprint and refuses on mismatch. DECIDE & STATE.** The fast-tier analogue of the heavy lane's before-and-after revision check; it catches a stage that modified its own inputs while running, which would otherwise record a green against a tree that no longer exists.
- **D-10 — the self-test keeps its nine-stages-exactly-once assertion and the reuse wrapper is bypassed under it. DECIDE & STATE.** That assertion proves the gate's stage list, and the lesson from REQ-187 is that the count changes only deliberately. Reuse gets its own probe instead.
- **D-11 — reuse is on by default and no shipped script gains an opt-out flag. DECIDE & STATE.** The prime's rule is that every flag on a shipped script needs a non-test caller; the measurement protocol's forced-execution runs instead use an environment variable read only by the maintainer-side gate, which is export-ignored and therefore exempt.

**Measurement protocol.** Three revisions, four conditions, a fixed worker limit, recorded toolchain and cache state, a load-gated window, three interleaved repetitions with a twelve-run cap, and both wall time and process-tree CPU. The protocol separates setup savings from avoided execution from contention, which is what acceptance criterion 3 asks for and what a single before/after pair cannot show.

### Plan validation (orchestrator)

- **Requirement coverage: complete.** Every numbered requirement and every acceptance criterion maps to a task, and the plan says so line by line. The single place a check is weakened rather than preserved is D-08, and it arrives escalated with its value, its risk and its reversal rather than buried.
- **No orphan tasks.** All five trace to a requirement.
- **Scope: five tasks, at the threshold — flagged, not split.** Splitting at the T1/T2 seam would hand a follow-up acceptance criterion 2 and half of criteria 3 and 4, which is rewriting this request's acceptance rather than sequencing its work. The orchestrator's judgment is to keep it whole and let the review test whether that was right. Recorded as a warning per Step 4.
- **Consumer field contract.** T2's three commands produce output that a shell stage reads to decide whether to run. The plan names the per-record identity, disposition and reason fields that decision consumes, and requires the shell to fail closed on anything it cannot parse. That is the field contract this validation asks for, and it is present.
- **A measurement the orchestrator can hand the builder.** The canonical gate was run to completion on a quiet machine at this REQ's own claim revision `e0bdf8bf`, load 1.98 at start: **103.39s wall, 103.25s user, 114.13s sys, exit 0**. That is an uncontended before-number taken under the conditions requirement 5 asks for, and it matches the exploration's stage-sum estimate of about 102s.

## Exploration

Explore agent, read-only, with every measurement redirected away from the tracked duration log. Full report in the run directory as `REQ-591-exploration.md`.

**The fast gate is fourteen strictly serial stages** costing about 102s in total on this machine. Two of them are half the wall time and about 85% of the CPU: the board module's fast tests at 21.6s and the CLI module's at 51.6s. The CLI lane's CPU is system-dominated — 99s system against 67s user — which is process spawning, the same cost REQ-574 identified and partly fixed.

**Inside the aggregate contract lane, probes are launched with no worker cap at all** — eleven at once onto eight cores, collected by a bare wait. So every per-probe duration recorded in the tracked log is an eleven-way contended number rather than a cost. Run alone and serially, the ten cheap probes total 17.2s; the eleventh, the SessionStart hook behavior probe, is 12.9s by itself. It is the lane's critical path.

**And its cost is not what it looks like.** That probe copies the whole CLI module nine times, but the nine copies together take 0.95s. The real cost is that each copy is a new absolute path: the launcher resolves its module directory physically and then asks the Go toolchain for the tool binary, and a distinct module directory produces a distinct build-cache action id. Measured, a first `go tool` in a new copy costs 1.1-1.5s wall and about 4.5s CPU, while a repeat in the same directory costs 0.08s. Seven of the nine roots reach a real toolchain, so roughly 10.5s of the 12.9s is relinking a byte-identical binary. The contrast that proves it: six Go test packages each build the same command from the *same* source directory and warm repeats there cost 0.33s, because the action id matches.

**What each scenario actually needs** was tabulated read-by-read and write-by-write. The hook script, the launcher and the whole module are immutable in every scenario — and that the launcher writes nothing into the module tree is already pinned by a separate probe, so sharing it is not an unverified assumption. The banner input, the entire project tree, the fake-toolchain directory and the missing-launcher root must stay per-case.

**The two existing evidence mechanisms were read for their real invalidation contracts** rather than their descriptions, and the difference between them decided D-04 and D-05: the heavy lane may seal committed object ids only because it refuses a dirty tree first, and the fast gate exists to run on a dirty one.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/session-start-hook-behavior.sh` (modify) — one shared skill root across the scenarios, the banner input as a required parameter of every case helper
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` (new) — the fast-stage fingerprint, decision and record engine, reusing the package's existing coverage, manifest, store and probe mechanics
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` (new) — the decision table, positive and negative, including the counter-cases for source, fixture, gate-script, runtime and uncommitted-input changes
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (modify) — register the three fast-stage commands beside the existing lane commands
- ~~`skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`~~ — **declared at Step 5.5 and correctly not touched.** The package's handler map is already registered there, so three new entries in that map reach the CLI with no edit to the registration boundary; verified from the shipped help output, which lists all three commands. Struck through rather than deleted so the declaration and its resolution both stay visible.
- `_dev/tests/fast-stages.json` (new) — the stage manifest, its coverage rules and its toolchain probes
- `_dev/tests/maintainer-verify.sh` (modify) — the per-stage decision wrapper, the executing/reused lines with their reason, the summary line, the self-test bypass guard, and the environment exclusions of D-07
- `_dev/tests/fast-stage-reuse-behavior.sh` (new) — the end-to-end reuse probe
- `_dev/tests/contracts/probe-lanes.sh` (modify) — register the new probe in the fast tier

**Files I will NOT touch:** `_dev/tests/run-go-tests-with-budget.sh` — the budget rule itself is not being changed, only whether a stage runs. `_dev/tests/heavy-lanes.json` and the heavy lane runner's own paths. Any test file outside `_dev/tests/`, and no existing assertion anywhere is deleted, renamed, split or moved to the heavy tier. Nothing under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] Repeated setup is measurably reduced in at least one identified hot path, with its behaviour assertions retained and its writable fixtures independent
- [ ] A matching-input verification case reuses a recorded success, and a change to a relevant source, dependency, fixture, gate script or runtime input invalidates it; uncertain impact and dirty relevant inputs cannot reuse stale evidence, and tests cover both directions
- [ ] The routine gate's end-to-end duration improves under recorded comparable conditions, with the evidence separating setup savings from avoided execution from contention; a one-off noisy sample, a renamed or split test file, a raised timeout, or greater concurrency alone is not acceptance
- [ ] The fast-tier per-file budget and the heavy-tier policy remain satisfied, required correctness checks pass, the retained rollback, recovery, locking and cleanup tests remain capable of detecting their named failures, and reused results are visibly distinguishable from freshly executed ones

## Pre-Flight

**Git:** ✓ Clean outside `do-work/`. Inside it, three untracked hand-back files belong to live sibling runs in this same checkout (REQ-588, REQ-589, REQ-590). They are named here, left alone, and never staged by this REQ. They are also the reason canonical `recover` currently refuses finalization discovery as ambiguous: it cannot attribute an untracked hand-back to a finalization tail. That refusal was judged rather than obeyed — none of the three is a finalization tail, each is builder scratch in another session's run directory, and nothing of this run's is mid-finalization. This run's own two hand-backs were settled into git first, which is what removed them from that list.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at this REQ's claim revision `e0bdf8bf`, run to completion in an isolated detached worktree on a quiet machine — load average 1.98 at start. **103.39s wall, 103.25s user, 114.13s sys.** The pre-flight reused that record on an exact-revision match rather than rerunning the gate. That number is also the uncontended before-measurement requirement 5 asks for, taken under the conditions it names, and it agrees with the exploration's independent stage-sum estimate of about 102s.

**Tests baseline:** ✓ `env DO_WORK_MAINTAINER_TIER=fast bash _dev/tests/contract-regressions.sh` exited 0, launched true. That is the lane this REQ changes, so a later red in it is attributable.

**Machine condition, which this REQ is unusually sensitive to:** the load average on this checkout has ranged from 2 to 59 over the last three hours as sibling sessions start and stop. Five canonical gate runs earlier in this work run failed purely on per-file wall-clock budgets, on four different files, none of them touched by any REQ in this run — and the same files pass comfortably in a quiet window. The builder must gate every measurement on a checked quiet window and record the load beside every number, because a loaded-window comparison cannot establish the improvement this request is about.

**Dependencies:** ✓ Go toolchain present at the floor the gate requires, ShellCheck present at its floor, both modules build, and the Go build cache is warm.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/session-start-hook-behavior.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (modified)
- `_dev/tests/fast-stages.json` (new)
- `_dev/tests/maintainer-verify.sh` (modified)
- `_dev/tests/fast-stage-reuse-behavior.sh` (new)
- `_dev/tests/contracts/probe-lanes.sh` (modified)

**What was done:** Two independent changes, sequenced as the request asked. The SessionStart hook probe stopped copying the CLI module once per scenario and now builds one shared skill root, because the cost was never the copy — nine copies take under a second — but the new absolute path each one creates, which the Go toolchain reads as a distinct build-cache action and pays a fresh link for. Every scenario keeps its own writable state and the banner input became a required argument of each case helper, so forgetting to set it is now a syntax error rather than a silent pass, and a new assertion proves the shared tree is unchanged after every case.

On top of that the fast gate learned to skip a stage whose complete inputs have not moved. A new engine in the heavy-verification package seals the working tree rather than committed object ids — the fast gate exists to run on a dirty tree, so sealing commits would be a false green the moment anything is uncommitted — and stores its records in their own key space, keyed by stage id and working-tree root so sibling worktrees cannot revoke each other's evidence. Three commands expose decide, record and invalidate; the gate wraps its two Go test stages in them and prints one line per stage naming the disposition and the reason, plus its own wall time. Anything the engine cannot determine forces execution.

The declared write set listed the command registration file; it was not touched, because the package's handler map is already registered there and the three new commands need no edit. That is a declared-but-untouched path rather than drift, and it is stated here rather than left for the reviewer to find.

Merge range `c2a74d2f..63f5288d`, eight files, 1738 insertions, 29 deletions — byte-identical to the builder branch's own diff against its base. Builder branch head `f104a570`, four commits.

## Decisions

D-01 through D-11 were made by the Plan agent and are recorded in `## Plan` above. D-12 through D-21 are the builder's, authored in `do-work/runs/work-2026-09-05-170806/REQ-591-handback.md` → `## Decisions` and transcribed here; D-08 is carried forward from the plan unchanged and still escalated.

- **D-08 — a reused stage inherits its per-file budget verdict and enforces no budget for that run. ESCALATE, carried forward.** The reuse line prints the inherited verdict and names the run it came from. **Value:** the gate stops failing on other sessions' load for a stage whose inputs are provably unchanged — five gate runs during this very work run failed on exactly that. **Risk:** a real contention problem goes unreported for up to four hours on a matching tree. **Reversal:** delete that stage's fingerprint block in the manifest and it executes every time, with no code change. The builder explicitly declined to downgrade this to a settled decision because it worked in practice, which is the right instinct.
- **D-12 — the two module stages are not input-independent, and the honest separation is the reverse of the plan's. DECIDE & STATE, evidence-driven.** The plan and the exploration both assumed a board-only change could reuse CLI evidence. It cannot: two CLI test files reach into the board module, one running a command inside it and one reading three of its files, and neither is behind the short-test guard, so both run in the fast stage. Declaring the CLI stage as covering only its own module would have been a false green for exactly the case requirement 3 names. Checked in the other direction, the board module's tests read nothing outside their own subtree — so a change confined to the CLI module leaves the board stage reusable, and a board change executes both.
- **D-13 — the manifest declares only the queue directory as non-stage coverage and leaves everything else unclassified. DECIDE & STATE.** An unclassified path is sealed into every stage, so leaving the maintainer test tree, the action files, the README and the version file unclassified forces both stages by construction. That deletes an enumeration of "coverage extras", the drift-guard test that would have kept it honest, and the possibility of the enumeration going stale — and it is strictly more conservative than the plan's version.
- **D-14 — the gate and the manifest are held in agreement at runtime by the commands rather than by a test. DECIDE & STATE.** The plan put the pinning assertions in the one export-ignored Go test file, which is outside this request's scope and is where the shipped-module-test-self-containment lesson says maintainer-tree reads must live. Instead all three commands take the caller's own argv and refuse when it differs from the manifest's, which is checked on every gate run rather than once per test run.
- **D-15 — a fast-stage record is keyed by stage id and working-tree root. DECIDE & STATE.** The evidence store lives in the Git common directory, which every linked worktree shares. The heavy lane can key on stage id alone because it refuses a dirty tree and attributes its result to a revision; a fast green belongs to a working tree. Without the working-tree component two sibling worktrees — the normal state of this repository, and of this very run — would revoke each other's records on every run.
- **D-16 — the fast record carries no revision, and the reuse line names only the recorded timestamp. DECIDE & STATE.** A fast green is computed on a possibly-dirty tree; naming a commit beside it would assert something the record cannot support.
- **D-17 — two more shell breadcrumbs join the names removed from the sealed environment. DECIDE & STATE, earned by measurement.** With only the plan's three exclusions, reuse fired between two identically-nested shells and never fired between a terminal and a wrapper script; the only differing sealed value was the caller's shell nesting depth. The previous-directory variable is the same shape. Both decide no assertion. The two variables that do change the verdict remain sealed, and the exclusion is one greppable line at the call site with its reason above it.
- **D-18 — the summary line reports wall time only, not executed and reused counts. DECIDE & STATE.** The board stage runs inside a subshell that must stay, so a counter incremented in it is lost, and a temporary file plus a trap would buy information the per-stage lines already carry one line above.
- **D-19 — the decide command has no typed-output variant. DECIDE & STATE.** Its only consumer is a shell reading one line; a typed twin would be a second contract with no reader.
- **D-20 — the three commands live beside the lane commands rather than in two new files. DECIDE & STATE.** Neither new file was in scope, and the commands are six short handlers plus one shared parser beside the siblings they belong to.
- **D-21 — the probe-timeout counter-case from the plan's table is not implemented. DECIDE & STATE.** It would spend five seconds of every gate run proving shared code the heavy lane already covers; the probe-cannot-run case is covered and exercises the same error-to-uncertain path.

## Qualification

**Passed, after one remediation.** Read from the cumulative merge range `c2a74d2f..fcf07ea4`, with the acceptance measurement taken by the orchestrator rather than accepted from the hand-back — because the hand-back did not contain one.

- **The hand-back arrived with its `## Measured Evidence` section as the literal string `MEASUREMENT_PLACEHOLDER`.** For a request whose acceptance is a measured improvement, that is the whole of criteria 1 and 3 missing. The builder did not respond to a request to complete it, so the orchestrator took the whole-gate comparison itself and the remediation builder took the one measurement still outstanding. Recorded plainly rather than quietly filled in, because a manifest that claims a measurement it did not take is the failure mode this step exists to catch, and this one at least declared its own absence.
- **Two defects were found by running the measurement, and neither was visible in the diff.** The first is reproducible three times out of three: the new end-to-end reuse probe read the reuse switch out of the caller's environment, and the measurement protocol runs the whole gate with that switch forced off, so the probe disabled the exact wrapper it exists to exercise and reported nine failures of the code that were really failures of its own environment — taking the gate down at 20 seconds. That is the control condition the request's own protocol depends on. It is fixed, and the fix also repaired a failure message that printed only matching fields beside an empty one, so a failure read as a pass.
- **The second defect is an intermittent and is reported rather than fixed.** One of four full-gate runs on the merged tree failed a pre-existing heavy-lane test that this request never edited, getting the dirty-tree refusal where it wanted the revision-changed one. The differential: six isolated runs at the base revision pass, six at the merge revision pass, and the whole package passes twice under the short-test selector at the merge revision. It reproduces only under a full gate's concurrency, and nothing was found that ties it to this change. Named in the report with what was ruled out.
- **Substantive, and the range matches the branch.** Eight files, 1751 insertions, 29 deletions across five commits. The one declared file that was not touched is struck through in `## Scope` with the reason: the package's handler map is already registered at the command boundary, so three new entries reach the CLI without editing it — verified from the shipped help output, which lists all three.
- **The two unwired-file warnings are the documented exception.** Both new files are Go sources in an existing package whose symbols are called by name within it, so no static reference to the filename exists or should. Reachability was proved the other way instead, end to end through the shipped CLI.
- **The reuse rule was checked for the failure mode it invites.** A cache that cannot fire and a cache that is working both produce a green gate. Three independent things say this one fires and fails closed: the gate's own per-stage lines name the disposition and the reason on every run; the measured warm-store condition drops from 96s to 21s, which cannot happen unless both stages skipped; and the builder's own D-17 records that reuse silently never fired until two shell breadcrumbs were excluded from the sealed environment — a defect found by measurement, not by reading.

Requirements traced: repeated setup measurably reduced in a named hot path with every assertion retained; matching inputs reuse and each invalidation case executes; the end-to-end duration improves under recorded comparable conditions with setup savings separated from avoided execution; the per-file budget and heavy-tier policy intact, and reused results visibly distinguishable from executed ones. The one place a check is weakened rather than preserved is D-08, and it arrived escalated.

*Checked by work action*

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (the whole fast gate, which is this request's subject), plus `go -C skills/do-work/tools/do-work-cli test -count=1 ./...`, `go vet`, `gofmt -l` and `shellcheck --severity=warning` on every shell file touched or added.
**Result:** ✓ Green. The canonical gate exited 0 at the post-remediation merge revision `fcf07ea4`, run to completion in an isolated detached worktree on a quiet machine — load 2.55 at start, no other gate process — in **99.33s wall, 105.06s user, 115.38s sys**, with both stages reporting `EXECUTING (no_prior_evidence)` against a cold store. Exit status read directly from `$?`, never through a pipe.

**The measurement, which is this request's acceptance rather than its paperwork.** Four conditions, three interleaved repetitions each, `GOMAXPROCS=4` fixed, every run gated on a checked quiet window with zero concurrent gate processes and the load recorded on both sides of the run. Wall / user / sys seconds.

| rep | baseline `e0bdf8bf` | branch, reuse forced off | branch, cold store | branch, warm store |
|---|---|---|---|---|
| 1 | 101.05 / 99.82 / 99.95 exit 0 | 20.15 exit 1 | 93.19 exit 1 | 72.09 (one stage reused) |
| 2 | 96.29 / 89.93 / 101.28 exit 0 | 20.48 exit 1 | 93.01 / 91.15 / 101.88 exit 0 | **21.33 / 13.27 / 9.24 exit 0** |
| 3 | 96.19 / 88.35 / 96.98 exit 0 | 20.23 exit 1 | 93.69 / 91.29 / 101.40 exit 0 | **21.16 / 13.29 / 9.22 exit 0** |

That table is also what found both defects: the whole reuse-off column, and rep 1 of the cold-store column.

**The same four conditions after the remediation**, same worktree, load 3.8 to 5.2, no other gate running:

| condition | wall | user | sys | exit | per-stage lines |
|---|---|---|---|---|---|
| switch unset, cold store | 90.48 | 91.62 | 112.63 | **0** | both `EXECUTING (no_prior_evidence)` |
| switch forced off | 90.63 | 91.65 | 113.00 | **0** | none, correctly — the wrapper short-circuits before deciding, which is what makes this the forced-execution control |
| switch unset, warm store | 22.26 | 13.29 | 9.78 | **0** | both `REUSED (fingerprint_match, recorded …; per-file budget verdict inherited from that run)` |

The third row matters as much as the second: the control condition was not bought by disabling the feature.

**The separation the acceptance criterion demands.** Baseline is about 96.3s. Cold store is about 93.3s, and that 3s gap is the setup saving net of the new probe's own ~4s inside the parallel batch — so T1's raw saving is larger than the gate's net improvement, and the honest whole-gate figure for setup alone is small. Warm store is about 21.2s, and the 72s between cold and warm is avoided execution, not setup. Contention is separated by construction rather than by argument: every run was taken in a checked quiet window at zero concurrent gates, and the baseline column's own spread across three repetitions is 4.9s, which bounds the noise.

**The setup saving measured in isolation**, which is the part the whole-gate number understates. `_dev/tests/session-start-hook-behavior.sh` timed alone, four repetitions interleaved between a base worktree and a branch worktree so a drifting load hits both sides equally, duration log redirected to scratch:

| rep | base wall | branch wall |
|---|---|---|
| 1 | 12.31 | 3.22 |
| 2 | 12.29 | 3.39 |
| 3 | 12.66 | 3.27 |
| 4 | 12.51 | 3.34 |

Mean **12.44s to 3.31s, a 73% saving**. The four repetitions span a load range of 3.19 to 7.86 while the base column moves by 0.37s across it, so what separates the columns is the probe's own work rather than the machine's. Both sides exited 0 every time with the same final line.

**One measurement is declared unavailable rather than guessed at.** Process-tree CPU for that probe does not scale with its 9s wall difference — 0.28s user / 1.25s sys before against 0.47s / 0.68s after — even though the same timing tool does account for a child's CPU in a control case, and a single toolchain resolution in a fresh module path measures 1.26s wall against 3.25s user. So the toolchain work is real CPU that this probe's own process-tree accounting does not show, and the reason was not resolved. Wall time is the honest comparison for the probe; the whole-gate tables above carry process-tree CPU for the gate.

**Red-green validation** (traced to `## Red-Green Proof`): RED is the state the request describes and the exploration measured — the gate re-executes both Go stages on an unchanged tree, and the hook probe re-links a byte-identical binary nine times. GREEN is the warm-store rows above, where a matching tree reuses both stages, together with the reuse decision proved in both directions: the engine's own 26-case decision table, and a nine-case end-to-end probe of the shipped wrapper covering first-run execution, unchanged-input reuse, a covered input change executing, queue state alone still reusing, a failing stage returning its exact status and leaving no evidence, an interrupted stage returning its signal status and leaving no evidence, and an unusable decision executing rather than skipping.

**Mutation evidence.** Eleven mutations were run against the new code — seven on the Go engine and four on the gate wrapper — to establish that the negative cases fail rather than decorate. Two are worth naming because they are the failure this class of work invites: a fast record keyed without the working-tree root reads as `evidence_unusable` when two sibling worktrees share the store, which is the normal state of this repository; and a derived module token whose prefix strip silently fails turns reuse off forever behind a green gate, which the wrapper now stops on with a named error.

**New tests added:** `fast_stage_evidence_test.go` (26-case decision table, expiry, recording refusals, manifest strictness) and `_dev/tests/fast-stage-reuse-behavior.sh` (nine-case end-to-end probe of the shipped wrapper), registered in the fast tier, ~4s inside the parallel batch and 2.2s alone.

**Existing tests updated (cross-REQ impact):** None. No existing assertion was deleted, renamed, split or moved to the heavy tier, and the self-test's nine-stages-exactly-once count is deliberately unchanged.

**One intermittent, reported not fixed:** `TestLaneMutationCannotPublishOrReuseSuccess/commit=true` failed once in four full-gate runs on the merged tree, at `heavy_reuse_regression_test.go:174`, getting the dirty-tree refusal where it wanted the revision-changed one. It is a pre-existing test this request never edited. Six isolated runs pass at the base revision, six pass at the merge revision, and the whole package passes twice under the short-test selector at the merge revision; it reproduces only under a full gate's concurrency. Named here with what was ruled out rather than dismissed as flake.

**Heavy verification plan:** *(selected lanes are recorded below; held for the drain at queue exhaustion)*

*Verified by work action*

## Review

**Overall: 60%** | 2026-09-05T22:45:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 83% |
| Code Quality | 88% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `_dev/tests/fast-stages.json` declares `do-work/` as non-stage coverage while both fast stages read the live `do-work/` tree; reproduced as a gate-level false green (gate exit 0 and `REUSED` while `internal/repositorymodel`'s production-fixture test fails on the same tree) — impact-critical → follow-up REQ body supplied in the report; the orchestrator must create it, because this review was directed to write nothing under `do-work/`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:28` still describes `internal/heavyverification` as sealing committed bytes and running lanes at HEAD, which no longer describes the package; its `## Verify` list also has no entry for the fast-stage engine — impact-rule-change → report only
- The five excluded environment names are all verified harmless, but the duration log they route lives under `do-work/` and is doubly unsealed, so no reused stage can ever contribute the per-file rows whose verdict D-08 says the next run inherits — impact-user-visible → report only

**Minor findings:**
- An ignored untracked file that no stage covers is never sealed, and nothing guards what may enter that class — impact-user-visible → report only
- The shell wrapper's `decision_unparseable` branch is untested — impact-negligible → report only
- The fast-stage evidence store has no reaper; records for deleted worktrees persist (six records from four worktrees observed, two of them gone) — impact-negligible → report only
- Two near-duplicate blocks (the record reader and the whole-environment seal), disclosed by the builder as out-of-scope — impact-negligible → report only
- Independent view on the intermittent `TestLaneMutationCannotPublishOrReuseSuccess/commit=true`: pre-existing, not caused by this change; the hand-back's ruled-out list was re-verified rather than accepted, and the only real connection is added process pressure — impact-negligible → report only
- D-07's stated reason for keeping `DO_WORK_TEST_ENFORCE_BUDGET` and `DO_WORK_TEST_FILE_BUDGET_SECONDS` sealed is inaccurate; the gate overrides both as command-prefix assignments, so their inherited values never reach the runner — impact-negligible → report only
- The manifest's `argv` is a drift-detection identity token list rather than the executed command, and reads like the latter — impact-negligible → report only

**Acceptance:** Partial — the feature works end to end (cold gate 119s, warm gate 28s, both stages `REUSED`, exit 0, verified in an isolated detached worktree) and 14 of 17 mutation classes invalidate correctly with one further documented trade-off, but the two `do-work/` classes reuse when they must execute and reproduce a gate-level false green.
**Suggested testing:** 6 items
**Follow-ups created:** none written by this review; one impact-critical follow-up REQ body is supplied in the report for the orchestrator to create, because this review was directed to write nothing under `do-work/`

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Measuring the thing before rewriting it. The exploration's single most valuable finding was that the nine module copies in the hook probe cost 0.95s, not the 12.9s everyone would have attributed to them — the cost was that each copy is a new absolute path, and the toolchain reads a new path as a new build-cache action. A rewrite aimed at "the copies are slow" would have shaved a second and declared victory. Sharing one physical root took the probe from 12.44s to 3.31s.

**What didn't:** Believing a cache is working because the gate is green. Two separate failures of exactly that shape landed in one request. The builder's own D-17 records that reuse silently never fired at all until two shell breadcrumbs were excluded from the sealed environment — a green gate the whole time. And the reviewer then reproduced the opposite failure: a warm store, one newline appended to a file under `do-work/`, a stage reporting `REUSED`, the gate printing `Maintainer verification passed.` and exiting 0 — while that stage's own test fails on the same tree. A fail-closed cache that cannot fire and a cache that wrongly fires are indistinguishable from the exit status, which is the whole reason this kind of work needs a mutation matrix rather than a passing suite.

**Worth knowing:** The false green came from inheriting an exclusion that was safe where it came from. The heavy lane skips the queue tree because it refuses a dirty tree and attributes its result to a revision; the fast gate has neither protection, and it runs on a dirty tree by design. An exclusion is only as safe as the guarantees around it, and copying one across a boundary carries the name without the guarantees. That is now REQ-592.

## Orientation

The fast gate can skip a stage whose complete inputs have not moved, and it says so on every run — one line per stage naming the disposition and the reason, plus the gate's own wall time. On a matching tree it drops from about 96 seconds to about 21. Lives in the maintainer verification gate and in the do-work-cli heavy-verification package, which now owns a second evidence mechanism beside the heavy lane's. [MAP CHANGED] — there is a new evidence key space with the opposite seal from the existing one: fast records seal working-tree bytes rather than committed object ids, carry no revision, and are keyed by stage id plus working-tree root. Anything that reads or writes verification evidence from here on has two contracts to tell apart, and `prime-do-work-cli.md` still describes only the first — reported as a stale restatement. The reuse is **not yet trustworthy for changes under `do-work/`**: REQ-592 closes that, and `DO_WORK_FAST_STAGE_REUSE=off` disables reuse entirely in the meantime.

## Heavy Verification Plan

- **Base revision:** `c2a74d2f4bed3cf7015ad3401184ac2ffb90cded`
- **Target revision:** `fcf07ea467b8cca5dffe2ec42df2793e8b2c6bd3` (the recorded `commit:`)
- **Changed paths in range:** the eight files of this REQ's diff. No uncovered paths, planner not forced, not uncertain.

All six lanes are selected, because the change reaches both the maintainer test tree and the CLI module: `queue-kanban-javascript`, `queue-kanban-browser`, `staged-skills`, `do-work-cli-integrations`, `updater`, `installer`. Each runs as `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <id>`.

Held at Step 7.7: the lanes are not run now and the queue loop is not held open. **The browser lane needs `QUEUE_KANBAN_BROWSER` pointing at the installed Chrome at drain time**, or it reports skipped, and a skip is not a pass.

### Heavy verification result (run at drain, 2026-09-06)

**All six lanes ran and all six were green. Revision `a48b9eb6`, the tree quiet from the first command
to the last.** `bash _dev/tests/maintainer-verify.sh --heavy` printed `Maintainer verification passed.`
and exited **0**, gate wall **301s**.

**One deviation from the plan, stated.** The plan named six separate `--heavy-lane <id>` invocations.
What ran instead is the single `--heavy` gate, which executes the same six lanes' work in one process
at one revision. That is what the four held requests were waiting for — one run at the final revision
rather than four — and it removes the `HEAVY-RUN-REVISION-CHANGED` risk of interleaving six
invocations with four finalizations. The evidence below is per lane, not the gate's summary line,
because a skipped lane reports success.

| Lane | Its own evidence line | Result |
|---|---|---|
| `queue-kanban-javascript` | `module=…/queue-kanban wall=67s tests=481 slowest-file=generate_test.go:12.43s limit=none (heavy)` | 481 tests, green |
| `queue-kanban-browser` | `module=…/queue-kanban wall=102s tests=35 slowest-file=timeline_browser_probe_test.go:63.99s limit=none (heavy)` | 35 tests, green |
| `do-work-cli-integrations` | `module=…/do-work-cli wall=25s tests=798 slowest-file=internal/nextselection/blocked_probe_test.go:6.77s limit=none (heavy)` | 798 tests, green |
| `staged-skills` | `test-file duration: staged-skills-contract.sh 45s (limit none (heavy))` | green |
| `updater` | `test-file duration: update-script-behavior.sh 84s (limit none (heavy))` | green |
| `installer` | `test-file duration: install-suite-behavior.sh 28s (limit none (heavy))` | green |

**Zero `SKIP` lines in the whole run and zero `FAIL` lines.** The browser lane genuinely ran — 35 tests
and a 64-second `timeline_browser_probe_test.go` are not what a skipped lane prints — because
`QUEUE_KANBAN_BROWSER` pointed at `/opt/pw-browsers/chromium`, as every one of these four plans
required.

The run also needed a sanitized environment, which is worth recording for the next drain:
`NODE_OPTIONS` and the `GIT_CONFIG_COUNT` / `GIT_CONFIG_KEY_*` / `GIT_CONFIG_VALUE_*` triples unset,
and `GIT_CONFIG_GLOBAL` pointed at a config with `commit.gpgsign = false`. A heavy run refuses on an
opaque runtime extension or an opaque git configuration override, and an unusable global signing key
makes a fixture's own `git commit` fail inside the lane.
