# Independent review — REQ-591 (reduce repeated setup and unaffected reruns in the fast gate)

## Review: REQ-591

**Approve with follow-ups** — the setup saving and the reuse engine both work end to end and fail closed on nine of ten input classes, but the tenth is a demonstrated false green: the gate prints `Maintainer verification passed.` and exits 0 while one of the reused stages' own tests fails on the current tree.
Route C | merge range `c2a74d2f..fcf07ea4` (5 commits, 8 files, 1751 insertions, 29 deletions)

### What's built

- The SessionStart hook probe stopped rebuilding the same Go tool binary nine times. One shared physical skill root, per-case writable state, the banner input now a required argument, and a new assertion that the shared tree is unchanged after every case. Measured 12.44s → 3.31s standalone.
- The fast gate learned to skip a Go test stage whose complete inputs have not moved. A new engine in `internal/heavyverification` seals working-tree bytes (not committed objects), stores records in their own key space keyed by stage id plus working-tree root, and exposes decide / record / invalidate as three CLI commands the gate calls around each of its two Go test stages. I confirmed end to end in my own detached worktree: cold gate 119s, warm gate 28s, both stages printing `REUSED (fingerprint_match, recorded …; per-file budget verdict inherited from that run)`, exit 0 both times.
- What is still missing: the manifest declares `do-work/` as a tree no stage reads. Both stages read it. See F1.

### Decisions / risks for you

- **D-08 (a reused stage enforces no per-file budget for that run) is correctly escalated, and its blast radius is wider than the hand-back states.** The four-hour ceiling *is* enforced for fast records — verified in code at `fast_stage_evidence.go:508` and pinned by `TestFastStageEvidenceExpiresIndependentlyOfTheFingerprint`. The stated risk is "a contention problem goes unreported for up to four hours". There is a second, input-caused case: the board module's live tests walk the whole `do-work/` tree (`board_live_test.go:12-44`), so a growing archive makes those test files genuinely slower — a real budget regression, caused by an input, that the seal cannot see because `do-work/` is unsealed. Accept D-08 if you want, but on the current manifest it does not only suppress contention noise.
- **The `do-work/` exemption is pinned as correct by two new tests.** `fast_stage_evidence_test.go:183-185` comments it as "Queue bookkeeping is not stage input", and `fast-stage-reuse-behavior.sh` case 4 is named `queue state alone still reuses`. Fixing F1 means changing those two assertions, not only the manifest.
- **D-12's asymmetry is right and I verified both directions independently.** A change confined to the `do-work-cli` module leaves the board stage reusable; a board change forces both. Confirmed by reading and by execution (see Acceptance).

### Findings

**Important:**

- **F1 — `_dev/tests/fast-stages.json` declares `do-work/` as `non_stage_coverage`, but both fast stages read the live `do-work/` tree, so a `do-work/`-only change reuses stale evidence and the gate reports a false green.** Demonstrated at the gate level, not inferred. In a detached worktree at `fcf07ea4` with a warm store I appended one newline to `do-work/archive/UR-003/input.md`. `TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass` (`skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go:397-416`, not behind `testing.Short()`) then fails with `production legacy fixture changed size: got 5609 bytes`; the whole gate run immediately afterwards printed `stage do-work-cli-fast-tests: REUSED (fingerprint_match…)`, then `Maintainer verification passed.`, exit 0. The board stage has the same exposure through `board_live_test.go`'s `liveBoard`, which builds the board from the real `do-work/` tree, and through `durations_test.go:270` and `citations_test.go`. The exclusion itself is inherited from the heavy lane (`queueStatePrefix`, skipped at `fast_stage_evidence.go:195` and `:223`), where it is safe because the heavy lane refuses a dirty tree and attributes its result to a revision; the fast gate has neither protection. This is exactly what Detailed Requirement 2 forbids ("Unknown impact, incomplete evidence or an unverifiable input must select the broader verification rather than produce a false green") and what Acceptance Criterion 2 requires ("changing a relevant … fixture … invalidates it") — the failing test calls that file its "production legacy fixture" in its own message. — **impact-critical → follow-up REQ specified below under Follow-ups created (this review was directed to write no file under `do-work/`; the orchestrator must land it)**
- **F2 — `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:28` still describes `internal/heavyverification` as owning only the heavy-lane contract.** It says the package seals "covered and globally unclassified **committed** bytes" and "runs named manifest lanes at HEAD". The package now also owns a second evidence mechanism with the opposite seal (working-tree bytes), a second key space (`do-work-fast-stages`, keyed by stage id plus working-tree root), no revision on the record, and three new commands. A reader or agent acting on that paragraph would believe a fast-stage record is revision-attributable, which is the one thing D-05 and D-16 say it must never be read as. The `## Verify` list at `:88-101` also gained no entry for the new engine. D-03 anticipated the package name under-describing its contents and deferred the rename; the prime's *description* is a separate, non-mechanical fix and was not done. — **impact-rule-change → report only**
- **F3 — the D-07/D-17 environment exclusions are correct, but the gate's own duration log is an input the exclusion set makes unobservable in a second way.** All five excluded names behave as documented — I forced each into the decision child and reuse still fired for `DO_WORK_TEST_RUN_ID`, `DO_WORK_TEST_OTHER_GATE_PROCESSES`, `DO_WORK_TEST_DURATION_LOG`, `SHLVL` and `OLDPWD`, while an unrelated new variable, `DO_WORK_TEST_ENFORCE_BUDGET=no` and `GOFLAGS` all forced execution (the last as `fingerprint_uncertain`, the toolchain probe refusing an opaque flag set). Grepping the readers confirms none of the five decides an assertion: `DO_WORK_TEST_OTHER_GATE_PROCESSES` and `DO_WORK_TEST_DURATION_LOG` are only written into `test-duration-log.sh`'s rows, and the one variable that does decide (`DO_WORK_TEST_ENFORCE_BUDGET`, read at `run-go-tests-with-budget.sh:78`) stays sealed. The finding is narrower than the exclusion: the duration log itself lives under `do-work/` and is therefore doubly unsealed, so no run of a reused stage can contribute the per-file rows D-08 says the next run inherits a verdict from. The inherited verdict can be arbitrarily old and the reuse line does not say how old in budget terms, only when the record was written. — **impact-user-visible → report only**

**Minor:**

- **M1 — an ignored untracked file that no stage covers is never sealed, with no guard on what may enter that class.** `fast_stage_evidence.go:226-229` skips it deliberately, with a stated reason (refusing every ignored file would disable reuse on any normal checkout). Today `.gitignore` holds only `.playwright-mcp`, `/build/`, `/do-work/test-durations.tsv` and `.DS_Store`, none of which any gate stage reads, so the hole is currently empty — I confirmed creating `build/review-probe.txt` leaves both stages reusable. Nothing stops a future ignored generated artifact from becoming a stage input inside that blind spot. — **impact-user-visible → report only**
- **M2 — the shell wrapper's `decision_unparseable` branch has no test.** `maintainer-verify.sh` invents two reasons the Go engine never emits, `decision_unavailable` and `decision_unparseable`. The end-to-end probe covers the first (case 9, unreadable manifest) and not the second, so the field-count and empty-field checks at the parse site are unexercised. — **impact-negligible → report only**
- **M3 — the fast-stage evidence store has no reaper.** Records are keyed by working-tree root and a deleted worktree's record is never removed. I found six records in `.git/do-work-fast-stages` from four different worktrees, two of which no longer exist. Small files, no correctness effect (a worktree recreated at the same path only reuses when the sealed bytes still match), but the directory grows without bound on a machine that uses worktree dispatch. — **impact-negligible → report only**
- **M4 — two near-duplicate blocks, disclosed by the builder rather than hidden.** `sealedEnvironmentVariables` (`fast_stage_evidence.go:279-305`) repeats `laneFingerprint`'s environment loop, and `ReadFastStageEvidence` repeats `ReadLaneEvidence`'s Lstat/permission/SameFile/read dance almost line for line. Extracting either needs `heavy_evidence.go`, which was outside `## Scope`. Recorded as a discovered task in the hand-back. — **impact-negligible → report only**
- **M5 — my own view on the intermittent, which differs from the hand-back only in emphasis.** I agree it is pre-existing and not caused by this change, and I checked the load-bearing claims myself rather than accepting them: `internal/heavyverification` contains zero `t.Parallel()` calls and zero `os.Setenv`/`os.Chdir` calls before or after this diff; the new engine makes exactly three Git calls, all read-only (`ls-files --cached`, `ls-files --others`, `rev-parse --git-common-dir`) at `fast_stage_evidence.go:186,218,351`; the two evidence stores are different directories with different digest inputs, and I confirmed by planting a real heavy-lane record at a fast stage's key that it reads as `evidence_unusable` rather than as evidence. The new shell probe runs in the aggregate contract stage, which is a different serial gate stage that completes before either Go stage starts, so it cannot overlap them. What remains is the hand-back's own honest connection: one more test file in the same package means more short-lived Git subprocesses inside one binary on a host whose `kern.maxprocperuid` is 3333. Pre-existing, exposed by process pressure. The named next step (pass a `LaneOutputWriter` so `laneSkipWatcher` stops discarding the lane's own error text at `heavy_run.go:425`) is the right first move. — **impact-negligible → report only**

**Nit:**

- **N1 — D-07's stated reason for keeping two variables sealed is inaccurate.** It says `DO_WORK_TEST_ENFORCE_BUDGET` and `DO_WORK_TEST_FILE_BUDGET_SECONDS` are not excluded because "they change the verdict". In this gate neither can: `run_budgeted_go_tests` sets both as command-prefix assignments on the stage invocation (`maintainer-verify.sh:84-86`), which overrides anything inherited, and `test_file_budget_seconds` is a literal 30. Sealing them is harmless and defensible as defence in depth; the reason written beside it is not the reason it is right. — **impact-negligible → report only**
- **N2 — the manifest's `argv` is an identity token list, not the executed command.** `["run-go-tests-with-budget.sh", "skills/do-work/tools/do-work-cli", "-short", "./..."]` omits the `bash`, the path, and the budget variables the gate actually runs. That is fine — it is a drift key, and D-14 makes the gate supply the same tokens — but the field reads like a command and one day someone will try to execute it. — **impact-negligible → report only**

### Requirements Checklist

- [x] **R1** reduce repeated setup in a measured hot path; immutable material shared, writable state isolated, real launcher / build failures / missing tools / installed layout retained — **delivered.** I read the T1 diff line by line: all four banner cases, both tail cases, the authority matrix, both launcher-failure cases and the missing-launcher case survive, each with its original assertion. `missing_root` keeps its own root, `$fake_bin` stays per case, the static grep is untouched, and `assert_shared_skill_root_unchanged` is added after every case. No assertion deleted, renamed, split, mocked, or moved to the heavy tier.
- [ ] **R2** avoid repeating verification whose complete relevant inputs are unchanged; unknown or unverifiable input must select the broader verification — **partially delivered.** Nine input classes invalidate correctly (below). The `do-work/` class does not, and both stages read it. See F1.
- [x] **R3** record which checks executed, which reused, and why; preserve failure and interruption statuses; a skipped/failed/incomplete run supplies no evidence — **delivered.** Per-stage `EXECUTING (<reason>)` / `REUSED (…)` lines plus a `gate wall <n>s` line, all observed on real runs. The probe pins exact status 9 for a failing stage and 143 for a signalled one, with no evidence left behind either way.
- [x] **R4** keep useful failure coverage; no silent assertion removal, no mock for a real boundary, no move to the heavy tier — **delivered.** Confirmed from the diff: the only test-file changes are additions. `--self-test` still asserts nine stages exactly once and passes (`Maintainer verification self-test passed.`, exit 0, run by me).
- [x] **R5** before/after on comparable revisions, fixed worker limit, recorded state, no competing gate, wall and process-tree CPU, smallest bounded comparison — **delivered.** Four conditions, three interleaved repetitions, `GOMAXPROCS=4`, load recorded on both sides, and one measurement explicitly declared unavailable rather than guessed. I did not re-take it.
- [x] **AC1** repeated setup measurably reduced, assertions retained, writable fixtures independent — **delivered** (12.44s → 3.31s, 73%).
- [ ] **AC2** a matching-input case reuses; a relevant source / dependency / fixture / gate-script / runtime change invalidates; uncertain impact and dirty relevant inputs cannot reuse stale evidence — **partially delivered.** See F1: a relevant fixture changed and the stage reused.
- [x] **AC3** end-to-end duration improves under recorded comparable conditions, separating setup savings from avoided execution from contention — **delivered.** Independently reproduced in kind: 119s cold, 28s warm, same worktree, same window.
- [x] **AC4** fast-tier per-file budget and heavy-tier policy still satisfied; retained rollback/recovery/locking/cleanup tests still detect their failures; reused results visibly distinguishable — **delivered**, with D-08 escalated and F3 noted.

### Acceptance Testing

**Result: Partial**

Everything below was run in my own detached worktree at `fcf07ea4`, created with `git worktree add --detach`, load 2.7–3.2, removed afterwards along with the two evidence records it wrote. The main tree was never touched and nothing was committed.

What passed:

- `bash _dev/tests/maintainer-verify.sh --self-test` → `Maintainer verification self-test passed.`, exit 0.
- `bash _dev/tests/fast-stage-reuse-behavior.sh` → `Fast-stage evidence reuse probes passed.`, exit 0.
- Full gate, cold store → exit 0, both stages `EXECUTING (fingerprint_mismatch)`, 119s.
- Full gate immediately after, warm store → exit 0, both stages `REUSED`, 28s. The feature fires end to end and the disposition lines are what the run prints.

Invalidation matrix, driven through the real `decide-fast-stage` against the shipped manifest. `board` is `queue-kanban-fast-tests`, `cli` is `do-work-cli-fast-tests`:

| mutation | board | cli | correct? |
|---|---|---|---|
| nothing changed | reused | reused | yes |
| board source edited, uncommitted | executed `fingerprint_mismatch` | executed | yes |
| CLI source edited, uncommitted | **reused** | executed | yes — D-12's asymmetry, verified in both directions |
| `_dev/tests/maintainer-verify.sh` edited | executed | executed | yes |
| `_dev/tests/run-go-tests-with-budget.sh` edited | executed | executed | yes |
| new untracked `.go` inside the board module | executed | executed | yes |
| `skills/do-work/actions/work.md` edited (unclassified) | executed | executed | yes |
| tracked covered file deleted from the worktree | executed `fingerprint_uncertain` | executed | yes |
| tracked covered file replaced by a symlink | executed `fingerprint_uncertain` | executed | yes |
| `chmod +x` only, contents unchanged | executed `fingerprint_mismatch` | executed | yes |
| new env var in the decision child | — | executed | yes |
| `GOFLAGS=-mod=mod` (toolchain probe refuses) | — | executed `fingerprint_uncertain` | yes |
| a heavy-lane record planted at the fast stage's key | executed `evidence_unusable` | — | yes |
| **new `do-work/queue/REQ-…` file** | **reused** | **reused** | **no — F1** |
| **tracked `do-work/` file edited** | **reused** | **reused** | **no — F1** |
| gitignored file at repo root (`build/…`) | reused | reused | documented trade-off, M1 |
| gitignored file inside the board module | executed | executed | yes |

Also verified: `.git/do-work-fast-stages` holds records from four separate worktrees with no collision, each carrying `working_tree_root`, so D-15's key space works on a machine that actually runs sibling worktrees.

Why Partial rather than Pass: the feature works end to end and 14 of the 17 mutation classes behave exactly as specified (one more is a documented trade-off, M1), but the two `do-work/` rows reproduce a gate-level false green, which is the specific outcome the request was written to prevent.

### Restatement Sweep

This diff redefines three things that are stated in more than one place: what the fast gate does (a stage may now report a stored result instead of running), what `internal/heavyverification` owns (a second evidence mechanism with the opposite seal and its own key space), and what a green gate attests (a run in which the two most expensive stages may not have executed).

Swept: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (the `internal/gateevidence` and `internal/heavyverification` descriptions at `:27-28`, the import rules at `:38`, the Advance phase table at `:40-51`, the `## Verify` map at `:88-101`), `_dev/primes/prime-shell-commands.md` (all five sections; § Unchecked Exit Status Reads as Content and § Every Flag on a Shipped Script Needs a Non-Test Caller are both honoured by the new wrapper, and neither restates the gate's behaviour), `skills/do-work/actions/work.md` and `actions/work-reference.md` (the canonical-gate lanes, the one-retry rule, the green-gate record contract, the heavy-lane hold and drain — none of which asserts that the gate executes every stage), `skills/do-work-board/docs/board-guide.md` and the board and toolbox action files (heavy-lane and source-ready language only, unaffected), `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and `_dev/primes/lessons-shell-commands.md`, and every consumer of the gate's stdout (`grep -rn "Maintainer verification passed"` outside `do-work/` returns only archived REQ records, so the new `gate wall <n>s` line and the per-stage lines break no parser).

One stale restatement found: `prime-do-work-cli.md:28`, recorded as F2. Its `## Verify` list also has no entry for the new engine, folded into the same finding. Everything else agrees with the new meaning.

### Suggested Additional Testing

1. After F1 is fixed, re-run the two `do-work/` mutation rows above and require both stages to execute; the fix is not proved by a green gate, because a green gate is what the defect produces.
2. Exercise a stage whose covered subtree contains a gitignored generated artifact the stage reads (M1) — there is none today, which is why nothing catches the class.
3. Run the gate twice from two different callers (an interactive shell and a `nohup`'d script) on one unchanged tree and require the second to reuse. That is the shape D-17 found by measurement, and only a cross-caller run can regress it.
4. Let a fast record age past four hours on an unchanged tree and confirm `evidence_expired`. The unit test pins the arithmetic; nothing exercises the real clock.
5. Run the gate concurrently in two sibling worktrees of the same repository and confirm neither revokes the other's record. D-15 is the reason this is safe and I verified the keys differ, but not under simultaneous writes to one shared store directory.
6. Kill a gate mid-stage (SIGINT at the Go test) and confirm the next run reports `no_prior_evidence` rather than reusing the pre-kill record. The probe covers a signalled stage command; it does not cover a signalled gate.

### Scores (on the record — not the headline)

**Overall: 60%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 83% | R2 and AC2 partially delivered; everything else fully delivered and independently verified |
| Code Quality | 88% | Follows the package's existing shapes exactly, comments carry the reasoning, names are two-word and greppable; two disclosed duplications, one factually wrong comment (`fast_stage_evidence_test.go:183-185`), one inaccurate rationale (N1) |
| Test Adequacy | 90% | 26-case decision table asserting exact reason codes, expiry independence, recording refusals, manifest strictness, a 9-case end-to-end probe, and 11 mutations. Deducted because two of those tests pin the F1 defect as correct behaviour, and `decision_unparseable` is untested |
| Scope | 100% | Exactly the 8 declared files; the one declared-but-untouched path is struck through with its reason and I confirmed the CLI registers all three commands without it |
| Risk | Critical | A demonstrated false green in the gate that protects every other change |
| Acceptance | Partial | Feature works end to end (119s → 28s, exit 0, correct disposition lines); two of seventeen mutation classes reuse when they must execute |

Arithmetic: (83 + 88 + 90 + 100) / 4 = 90.25; Acceptance Partial applies a 10-point penalty → 80.25; Risk Critical caps the overall at 60%.

### Follow-ups created

One critical finding, F1. Per Step 10 it routes to its own REQ in `do-work/queue/`, but this review was instructed to write nothing under `do-work/` and to produce only this file, so **the orchestrator must land the entry below**. Fold-First check: I found no pending REQ in any UR whose root cause is the fast gate's coverage declaration — the nearest neighbours are REQ-574 (reducing CLI fixture setup costs, already completed) and this REQ itself, whose `## Scope` explicitly excludes changing the manifest's meaning after the fact.

```markdown
---
id: REQ-NNN
status: pending
domain: testing
created_at: <current UTC instant>
user_request: UR-127
addendum_to: REQ-591
review_generated: true
impact: impact-critical
effort_estimate: effort-mechanical
title: '[impact-critical] Review fix: seal the do-work tree into both fast gate stages'
---

# Review Fix: Seal the do-work Tree Into Both Fast Gate Stages

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

`_dev/tests/fast-stages.json` declares `do-work/` as `non_stage_coverage`, which asserts that no
fast gate stage reads it. Both stages do. The do-work-cli stage's
`TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass` reads and byte-checks
`do-work/archive/UR-003/input.md`; the queue-kanban stage's `board_live_test.go`, `durations_test.go`
and `citations_test.go` build the board from the real `do-work/` tree. Because
`fast_stage_evidence.go` skips `queueStatePrefix` in both its tracked and untracked seal loops, a
`do-work/`-only change reuses stale evidence and the gate reports a false green.

Either seal `do-work/` into the stages that read it, or stop those stages from reading it. The
narrowest honest fix is the first: give the fast-stage seal its own exclusion set instead of
inheriting the heavy lane's `queueStatePrefix`, and declare in the manifest exactly which subtrees
under `do-work/` no stage reads (if any). The heavy lane may keep its own exclusion, which is safe
there because it refuses a dirty tree and attributes its result to a revision.

Two existing assertions pin the current behaviour as correct and must move with the fix:
`fast_stage_evidence_test.go` case `queue state changed` (expects `reused`), and
`_dev/tests/fast-stage-reuse-behavior.sh` case `queue state alone still reuses`.

## Context

Found during independent review of REQ-591. The review reproduced a gate-level false green: with a
warm evidence store, appending one newline to `do-work/archive/UR-003/input.md` makes the
do-work-cli stage's own test fail, while the whole gate prints `Maintainer verification passed.`
and exits 0 with that stage `REUSED`.

## Requirements

- A change to any `do-work/` path a fast gate stage reads must force that stage to execute.
- The manifest's `non_stage_coverage` must state only trees no stage reads, verified rather than
  assumed.
- The two tests that currently assert the opposite are updated in the same change, each naming the
  failure it now catches.
- The gate's own `do-work/test-durations.tsv` must keep not invalidating its own stage; it is
  gitignored and written by the stage itself, so it needs an explicit narrow exclusion rather than
  the whole-tree one.

## Red-Green Proof

**RED prompt/case:** In a detached worktree at the merge revision, run the gate once to record
evidence, append one newline to `do-work/archive/UR-003/input.md`, and run the gate again.
**Why RED now:** The second run prints `stage do-work-cli-fast-tests: REUSED (fingerprint_match…)`
and `Maintainer verification passed.` with exit 0, while
`go test -short -run TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass
./internal/repositorymodel/` fails on the same tree with `production legacy fixture changed size`.
**GREEN when:** The second run prints `EXECUTING (fingerprint_mismatch)` for that stage and the gate
exits non-zero, and the same sequence with only `do-work/test-durations.tsv` changed still reuses.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
```

---

## Review

**Overall: 60%** | 2026-09-05T22:44:10Z

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
