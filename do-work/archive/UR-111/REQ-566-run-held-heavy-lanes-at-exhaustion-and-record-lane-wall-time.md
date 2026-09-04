---
id: REQ-566
title: '[impact-rule-change] Run held heavy lanes at queue exhaustion without asking, and record per-lane wall time'
status: completed
created_at: 2026-09-04T13:19:11Z
user_request: UR-111
domain: backend
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
priority: now
effort_estimate: effort-substantive
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/publication_types.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - _dev/tests/maintainer-verify.sh
  - skills/do-work/actions/work.md
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work-board/tools/queue-kanban/timeline.go
route: C
planning_at: 2026-09-04T13:36:09Z
exploration_at: 2026-09-04T13:39:25Z
dispatch_at: 2026-09-04T13:40:45Z
builder_handback_at: 2026-09-04T14:08:09Z
integration_at: 2026-09-04T14:11:30Z
testing_at: 2026-09-04T14:12:48Z
status_changed_at: 2026-09-04T14:18:07Z
commit: 4a28946e782cf42397a329ae65b51ad3d5694a1b
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-04T13:24:40Z
  basis:
    - Route C
    - 8-file write set
    - 3 subsystems involved
    - 5 acceptance criteria
    - persistence changes
    - cross-route regression gates
heavy_verified_at: 2026-09-04T14:18:07Z
heavy_verified_revision: 58e1c9c948bb68f3805e704b9c7db39fff38f504
claimed_at: 2026-09-04T14:18:08Z
review_at: 2026-09-04T14:23:10Z
kb_status: pending
completed_at: 2026-09-04T14:27:48Z
release_at: 2026-09-04T14:27:48Z
---

# Run Held Heavy Lanes At Queue Exhaustion Without Asking, And Record Per-Lane Wall Time

## What

Two changes, in this order.

1. **Record per-lane wall time.** The generated `## Heavy Verification Result` section records wall seconds for every lane next to its exit status, and the typed `heavy_testing` evidence (`HeavyLaneResult` in `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go`) carries that field so the result is machine-readable.
2. **Lift the human permission prompt.** At queue exhaustion, the work loop runs the batched heavy plan itself at HEAD for every `pending-heavy-testing` REQ, then records green or red through the same `answer --manifest` transaction that `do-work clarify` uses today. The `pending-heavy-testing` status, the non-blocking hold, and the source-ready rule for dependents all stay. Only the permission question goes. `do-work clarify` keeps working for a maintainer who runs the lanes by hand.

The fold-first scan found no pending or pending-answers REQ in any UR that shares this root cause. REQ-564 (reuse matching per-lane verification evidence for four hours) avoids rerunning unaffected lanes and REQ-559 (retry a red repository gate once before a repair REQ) handles a red gate; neither removes the permission prompt or records lane durations.
Target revision: `4a28946e782cf42397a329ae65b51ad3d5694a1b`
Execution revision: `58e1c9c948bb68f3805e704b9c7db39fff38f504`

- queue-kanban-javascript: exit 0, 6s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: exit 0, 82s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: exit 0, 67s — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0, 25s — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0, 52s — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0, 24s — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Builder read the crew rules, four primes and three lessons satellites, then the plan and exploration notes; approach settled before code: one `run-heavy-verification` command beside the planner, `wall_seconds` carried identically through resultmodel and publication, the engine skip made honest in the shell, and the permission-gate sweep. (Mirrored from the builder hand-back by the orchestrator.)
- [x] **[APPLY]:** Every changed file is inside `## Scope`; nothing under `do-work/`, no release file, `heavy-lanes.json` and `main.go` untouched. Ten deviations recorded as D-02 through D-10 in the hand-back, none outside the scope's file set.
- [x] **[UNIFY]:** `git diff --stat a37ea0d8..HEAD` reviewed on the branch (15 files, +952/-63); `gofmt -l` and `go vet` clean in both Go modules; debug-artifact sweep over the range found nothing; per-file checks listed in the hand-back.
## Why

The permission gate was a cost control from when the whole maintainer gate cost about seven minutes per REQ. It now protects roughly three and a half minutes of CPU per batch, and it is the only thing that stops an unattended queue drain. The user's goal is to drain the queue faster; the human wait is the bottleneck, not the tests.

## Context

Measured from `do-work/test-durations.tsv` on 2026-09-04:

| run | wall time |
|---|---|
| typical full gate, 126 files | 180 to 220 s |
| worst of the last ten runs | 326 s |
| fast tier alone, 14 files | 43 to 63 s |

Batching already exists: on 2026-09-04 one heavy run at revision `c0d8ce1c` stamped five held REQs (REQ-475, REQ-483, REQ-485, REQ-502, REQ-561) with the same `heavy_verified_at`. The lane selector (REQ-563, archived) already limits each plan to the lanes whose coverage paths the diff touches.

Risks the builder must handle:

- **R1, a skipped lane must stop the run, not loop.** The browser lane skips when no Chrome binary is found, and a skip counts as not confirmed. An automatic run that meets a skip records the skip, leaves the REQ at `pending-heavy-testing`, reports a typed finding naming the lane, and ends the run. It never routes the REQ through remediation for a lane that did not execute.
- **R2, lane commands must survive the deletion of the maintainer gate script.** Every lane argv in `_dev/tests/heavy-lanes.json` still calls `_dev/tests/maintainer-verify.sh`, which is slated for deletion. Do not add a new dependency on that script; the lane runner reads argv from the manifest so the commands can move without touching this feature.

Out of scope: a duration budget guard ("run automatically under N minutes, ask above"). Add it only when recorded lane durations show a lane that needs it.

## Red-Green Proof
**RED prompt/case:** One REQ sits at `pending-heavy-testing` with a green fast tier and a recorded heavy plan, the ready set is empty, and the user runs `do-work run` unattended.
**Why RED now:** The run ends by asking for permission to run the lanes. The REQ stays held until a human answers, and the existing `## Heavy Verification Result` shape records exit status only, so nothing measures how long a lane took.
**GREEN when:** The same run executes the planned lanes at HEAD, records exit status and wall seconds per lane in `## Heavy Verification Result` and in the typed `heavy_testing` evidence, returns the REQ to `pending` with `heavy_verified_at` and `heavy_verified_revision` set, and the next selection returns it with `resume_phase: review`. A red or skipped lane leaves the REQ held with the lane named in a typed finding and the run ends cleanly.
**Validation:** User confirmed

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, `slugged: partial` so only selectable bare; matches because this REQ changes a status contract and downstream readers in `actions/work.md` and `actions/clarify.md`. Over the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5924 tokens, `slugged: partial`; matches because this REQ changes typed evidence the answer command writes and reads (family `opaque-evidence-projection`). Over the 2000-token budget.

## Assets

None.

---
*Source: "stopping on pending-heavy-testing because that needs your permission." <- I asked for this, but it turns out it's not very good, because the work just stops, what I need is an inteligent way to run it (either faster by 80/20 rule), or groupped with multiple tasks (but that can get complicated). Any suggestions? I would start by monitoring the duration and lifting the human blocking, after all the goal is to do the queue faster.*

---

## Triage

**Route: C** - Complex

**Reasoning:** The change spans three subsystems (work and clarify action prose, the Go heavy-verification and answer evidence types, and the contract tests), replaces a permission gate with an autonomous run path, and adds a schema field to typed evidence. Multi-system with an architectural choice about where the lane runner lives.

**Planning:** Required

## Plan

**Architecture decision:** one small Go command, `run-heavy-verification`, registered from `internal/heavyverification` beside `plan-heavy-verification`. It executes named lanes at HEAD, times them, and returns typed per-lane results. It does not plan (that is `plan-heavy-verification`), does not compare against the stored `## Heavy Verification Plan` prose (drift refusal stays action judgment), and does not build the `answer` manifest (clarify.md Step 2.5 already assembles it; work.md cites it). Batching falls out: the action unions lane ids across every held REQ, runs once, then answers each REQ from the subset it selected. Prose-only was rejected: the exhaustion paragraph would describe an exact command sequence, which is the maintainer's "programs beat prose" condition.

### Task 1 — A1: `wall_seconds` on lane evidence (Go, TDD)
1. `internal/publication/answer_test.go` — RED first: in `TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure` (~line 262) give the lane `WallSeconds: 26`, tighten the substring `"browser: exit 0"` to `"- browser: exit 0, 26s — "` (red case `"- browser: exit 1, 26s — "`). New `TestHeavyTestingAnswerRejectsNegativeWallSeconds` (expects `ANSWER-HEAVY-EVIDENCE-INVALID`). New `TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane` (`Skipped: true, WallSeconds: 1`: `confirmed` → `ANSWER-HEAVY-EVIDENCE-NOT-GREEN`; `answered` renders `"- browser: skipped, 1s — "`).
2. `internal/publication/publication_types.go` — `WallSeconds int \`json:"wall_seconds"\`` on `HeavyLaneResult` (plain int, always present).
3. `internal/publication/answer.go` — `validateHeavyTestingEvidence` refuses `WallSeconds < 0` under `ANSWER-HEAVY-EVIDENCE-INVALID` and copies the field; `appendHeavyVerificationResult` renders `- <lane>: exit N, <s>s — \`argv\`` and `- <lane>: skipped, <s>s — \`argv\``. Header lines unchanged. No exact-byte fixtures exist elsewhere.

### Task 2 — A2 runner: `run-heavy-verification` (Go, TDD)
4. `internal/heavyverification/heavy_run_test.go` (new) — RED first, reusing `newHeavyTestRepository`, `writeHeavyTestFile`, `commitHeavyTestChanges` from heavy_verification_test.go. Fixture manifest lanes with argv `sh lanes/<name>.sh` committed in the fixture: green.sh (`sleep 1; exit 0`), red.sh (`exit 3`), skip.sh (`printf 'SKIP: no browser\n'; exit 0`). Tests: `TestRunLanesRecordsExitStatusAndWallSeconds` (manifest order, exit 0 and 3, `WallSeconds >= 1` for green, `ExecutionRevision == HEAD`, command outcome success, warning `HEAVY-RUN-LANE-RED`); `TestRunLanesMarksSkipFromExplicitSkipLine` (`Skipped == true`, exit 0, warning `HEAVY-RUN-LANE-SKIPPED` with the SKIP line as evidence); `TestRunLanesRefusesUnknownLaneBeforeExecuting` (`HEAVY-RUN-LANE-UNKNOWN`, nothing executed); `TestRunLanesRefusesDirtyTrackedTree` (`HEAVY-RUN-DIRTY-TREE`); `TestRunLanesTerminatesTimedOutLaneGroup` (unix tag, shape of nextselection/blocked_probe_test.go: `sleep 30` with 1s bound → exit 124, group gone). `heavy_commands_test.go`: `TestRunHeavyVerificationHandlerRejectsMissingLane` (`HEAVY-RUN-USAGE`), `TestHandlersRegisterRunHeavyVerification`.
5. `internal/resultmodel/result_model.go` — `HeavyVerificationRun` and `HeavyLaneExecution` (tags `lane_id`, `command_argv`, `exit_status`, `skipped`, `wall_seconds`, identical to publication's so the action pastes `lanes` into `heavy_testing.lanes` unchanged), `heavy_verification_run` field on `CommandResult` beside `HeavyVerification`, nil-slice normalization, text renderer.
6. `internal/heavyverification/heavy_run.go` (new) — `RunLanes(...)`: resolve HEAD; refuse `HEAVY-RUN-DIRTY-TREE` on non-empty `git status --porcelain --untracked-files=no`; read manifest at HEAD via existing `readManifestAtRevision`; refuse `HEAVY-RUN-LANE-UNKNOWN` before running anything; run lanes sequentially in manifest order as `exec.Command(argv[0], argv[1:]...)`, Dir = root, `ownedprocess.ConfigureGroup` fail-closed (shape of toolboxcommands/report_image_process.go); stdout+stderr teed to CLI stderr and a line scanner; `Skipped` = any output line starts with `SKIP:`; timeout → `TerminateGroup`, exit 124; SIGINT/SIGTERM → 128+signal, remaining lanes not run. Result fields per lane: `lane_id`, `manifest_path`, `execution_revision`, `command_argv`, `skipped`, `exit_status`, `wall_seconds`; top-level `outcome`; findings `HEAVY-RUN-LANE-RED` / `HEAVY-RUN-LANE-SKIPPED` (warnings; a red lane is evidence, not a command failure).
7. `internal/heavyverification/heavy_commands.go` — `CommandRunHeavyVerification = "run-heavy-verification"`, `handleRunHeavyVerification`, `parseRunArguments` (`[--manifest _dev/tests/heavy-lanes.json] --lane <id> [--lane <id>...] [--lane-timeout-seconds 1800]`), added to `Handlers()`. main.go already ranges over `heavyverification.Handlers()`.
8. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — the `internal/heavyverification/` line: the runner owns lane execution and wall time; publication still owns result evidence.

### Task 3 — R1 at the source: isolated lanes print the same `SKIP:` line as `--heavy`
9. `_dev/tests/maintainer-verify.sh` — extract browser discovery (~584-605) and Node discovery (~554-567) into two functions that `run_heavy_lane` (~662) reuses for `queue-kanban-browser` and `queue-kanban-javascript`. Today `--heavy-lane queue-kanban-browser` with no browser runs the strict Go lane, which skips inside and trips the zero-probe guard: red exit, no `SKIP:` line. After: prints `SKIP: no browser is available; ... Set QUEUE_KANBAN_BROWSER to name one.` and exits 0, as `--heavy` does. Add a `--self-test` case mirroring the `heavy-no-node` assertion (~388). Not a new caller of the script (R2): it is the lane's own "did not run" contract.

### Task 4 — prose sweep and board label
10. `skills/do-work/actions/work.md` — L16: "…route answers to `do-work clarify`; held heavy lanes run themselves at exhaustion (Step 6.5)". L91: "permission-gated heavy-test hold" → "non-blocking heavy-test hold". L97: drop "permission-gated"; last sentence → "At queue exhaustion the loop runs the batched lanes at HEAD itself (Step 6.5, heavy-lane drain); no permission is asked." L189: "awaiting permission" → "awaiting the exhaustion drain". L539: "otherwise the separate clarify permission flow receives …" → "otherwise the hold below records each selected lane's stable id, exact argv, and reasons for the exhaustion drain"; delete "Never infer permission from a prior run or from the user's request to implement this policy."; keep the force-all sentence. L544 heading → "**Heavy-test hold (non-blocking):**", steps 1-4 unchanged. Replace L551-553 with "**Heavy-lane drain (at queue exhaustion):**" — when no claimable pending REQ remains but at least one undrained `pending-heavy-testing` REQ exists: per REQ recompute `plan-heavy-verification` at the stored base and `commit:`; refuse drift against the stored `## Heavy Verification Plan` (drifted or `historical-revalidation` plans stay held and route to clarify); union the lane ids; run once with the fenced runner command; one `answer` manifest per REQ per clarify.md Step 2.5 (all exit 0 and none skipped → `confirmed`; any skipped → no answer, REQ stays held, finding reported, excluded from further drains this run; else → `answered` listing every lane; `target_revision` = REQ `commit:`, `execution_revision` = runner's); recompute selection and continue (green REQs return claimable with `resume_phase: review`). A `commit:` drained once this run is not drained again; a dirty tree is a typed refusal, never a hand edit. L794 row → "Run the heavy-lane drain (Step 6.5). A REQ still held afterwards (skipped lane, plan drift, dirty tree) is listed in the composed exit summary with its typed finding; not a failed run."
11. `skills/do-work/actions/clarify.md` — L11: "whose lanes did not run at exhaustion (skipped lane, plan drift, or historical plan) and a maintainer wants to run them by hand". Step 2.5 title → "Run held heavy lanes by hand"; delete the batch permission question, keep malformed-entry refusals; replace "Run the deduplicated commands there and retain …" with the fenced runner command; delete "The work loop never runs these commands; clarify … owns the permission."; keep the historical `plan-heavy-revalidation` path and the isolated-worktree sentence; add wall seconds to the evidence list. L203 "left awaiting permission" → "left held".
12. `skills/do-work/actions/work-reference.md` — L230 comment "permission-gated heavy lanes" → "heavy lanes for the exhaustion drain"; L318 "while that permission hold is open" → "while that hold is open"; Composed Exit Summary gains section "Held-for-heavy-testing" (condition: any `pending-heavy-testing` REQ after the drain), line `REQ-NNN — [title] (held: <finding code> <lane>)`, remedy: a skipped browser lane needs `QUEUE_KANBAN_BROWSER=<path>` and a re-run of `do-work run`; plan drift or a historical plan goes through `do-work clarify`. Testing template unchanged.
13. `skills/do-work-board/tools/queue-kanban/timeline.go` (~525) — label `waiting for permission to run heavy tests` → `waiting for the heavy lanes to run at queue exhaustion`. Not test-pinned. Board versioning folds into the skill bump (prime-kanban-board).

### Testing approach (tdd: true)
RED tests named per task above, written and observed failing before each implementation. No contract test under `_dev/tests` pins the heavy prose (sentence pins were retired in 0.273.1; do not add one); the lock-in is the Go pair: runner emits `HEAVY-RUN-LANE-SKIPPED`, answer refuses `confirmed` with a skip. Shell: `bash _dev/tests/maintainer-verify.sh --self-test` gains the browser-skip case. Gates: `go test ./internal/publication ./internal/heavyverification ./internal/resultmodel`, `go vet ./...`, `GOOS=windows GOARCH=amd64 go build ./...`, `bash _dev/tests/action-shell-blocks.sh`, `bash _dev/tests/contract-regressions.sh`, then the unpiped `bash _dev/tests/maintainer-verify.sh`.

### Release inputs (Step 9)
Shipped paths change, so this is a maintainer release. Version mirrors: `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (`**Current version**:` line), root `CHANGELOG.md`, and `skills/do-work/CHANGELOG.md` byte-identical to root. `suite/modules.tsv` untouched. Drift found: `CHANGELOG.md` head and `version.md` say 0.274.1 while both `VERSION` files say 0.274.0 (commit 69305f98 bumped only the changelog and version.md). Bump from the higher: release 0.275.0 (new command plus status-contract change) and repair both VERSION files in the same release transaction. One release entry after Task 4, title "Heavy Lanes Run Themselves at Queue Exhaustion and Record Wall Time".

### Risks
- R1: Task 3 makes an isolated skip honest at the source; the runner keys `skipped` on the condition "a line starts with `SKIP:`" (no lane-name enumeration); the answer side already refuses `confirmed` with a skip and the action never submits `answered` for a skip; drain per REQ per `commit:` at most once per run; other held REQs whose lanes all executed still get answered (one broken pipe).
- R2: the runner takes `--lane <id>` and reads argv from the manifest at HEAD; nothing in Go names the maintainer gate script; only `_dev/tests/heavy-lanes.json` changes when the commands move.
- Out of scope, recorded: no evidence cache (REQ-564), no duration budget guard, no retries (REQ-559), historical-revalidation plans stay clarify-only.

*Generated by Plan agent*

### Plan validation (orchestrator)
- Requirement coverage: A1 → Task 1; A2 → Tasks 2 and 4; "clarify keeps working by hand" → Task 4 item 11; "status, hold, dependents stay" → work.md steps 1-4 unchanged; R1 → Task 3 plus runner and answer guards; R2 → runner reads argv from the manifest. No uncovered requirement.
- No orphan tasks: the board label (item 13) closes a restatement of the permission gate, which the Restatement Sweep would otherwise flag.
- Scope sanity: 4 tasks across 13 files against an 8-file estimate. Warning only; the estimate is frozen. Task 3 is the only cut candidate and cutting it violates R1, so it stays.
- Consumer field contract: per-lane identity, provenance, state, and outcome fields are named in Task 2 item 6 and match publication's tags.
- Orchestrator guidance for the builder, D-01: the dirty-tree refusal must ignore paths under `do-work/`. Queue bookkeeping is never lane input, and this checkout routinely carries uncommitted `do-work/` edits from the orchestrator at the moment the drain runs. Dirt outside `do-work/` still refuses.
- The builder must read `_dev/primes/prime-kanban-board.md` before touching timeline.go (item 13) and `_dev/primes/prime-shell-commands.md` before editing the fenced runner command or maintainer-verify.sh.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (modify) — `WallSeconds` on `HeavyLaneResult`
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — validate and render wall seconds
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modify) — RED tests for A1 and the skip refusal
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` (new) — lane runner
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go` (new) — RED tests for the runner
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (modify) — `run-heavy-verification` handler and registration
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go` (modify) — handler usage and registration tests
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — `heavy_verification_run` result type, normalization, renderer
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — renderer and normalization coverage, only if the package's existing tests cover the sibling field
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — heavyverification ownership sentence
- `_dev/tests/maintainer-verify.sh` (modify) — isolated browser and JavaScript lanes print `SKIP:` and exit 0; self-test case
- `skills/do-work/actions/work.md` (modify) — hold wording and the heavy-lane drain at exhaustion
- `skills/do-work/actions/clarify.md` (modify) — by-hand path calls the runner; permission question removed
- `skills/do-work/actions/work-reference.md` (modify) — schema comments, source-ready wording, exit-summary section
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modify) — held-status label

**Files I will NOT touch:** `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (release tail is the orchestrator's at Step 9); `_dev/tests/heavy-lanes.json` (argv unchanged); anything under `do-work/`; `cmd/do-work-cli/main.go` (already ranges over the package handlers).

**Acceptance criteria (restated from REQ):**
- [ ] `## Heavy Verification Result` lists exit status and wall seconds for every lane, the typed `heavy_testing` evidence carries `wall_seconds`, and a negative value is refused
- [ ] `run-heavy-verification` executes manifest lanes at HEAD in manifest order, records exit status, skip, and wall seconds per lane, emits typed `HEAVY-RUN-LANE-RED` and `HEAVY-RUN-LANE-SKIPPED` findings, refuses an unknown lane and a dirty tree outside `do-work/` before running anything, and bounds each lane with a timeout that terminates the process group
- [ ] At queue exhaustion the work loop runs the drain and answers each held REQ through the answer transaction without asking; a skipped lane leaves the REQ held with a typed finding and the answer side refuses `confirmed` with a skip
- [ ] `do-work clarify` remains the by-hand path and calls the same runner; no restatement of the permission gate remains in work.md, clarify.md, work-reference.md, or the board label
- [ ] Isolated `queue-kanban-browser` and `queue-kanban-javascript` lanes print the `SKIP:` line and exit 0 when the browser or Node is absent, pinned by a `--self-test` case


## Pre-Flight

**Git:** ⚠ 2 pre-existing uncommitted files from a concurrent `do-work run-simple-reqs` session on this checkout — preserve and exclude from this REQ (`_dev/tests/update-script-behavior.sh`, `do-work/working/REQ-545-synchronize-the-self-signalling-fetcher-probe.md`). The builder runs in an isolated worktree, so neither can enter this REQ's build or merge range.
**Tests baseline:** ✓ `check-green-gate` matched HEAD `a37ea0d8` exactly (persisted green run by the background gate runner), so the pre-build gate was skipped; baseline revision a37ea0d8. No `baseline.json` was written to avoid overwriting the other session's record.
**Dependencies:** ✓ Go toolchain present

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `_dev/tests/maintainer-verify.sh` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified)

**What was done:** Added `wall_seconds` to per-lane heavy verification evidence (validated non-negative, rendered as `exit N, Ss` or `skipped, Ss`), added the `run-heavy-verification` command that executes named manifest lanes at HEAD inside owned process groups with a timeout and returns typed per-lane results plus `HEAVY-RUN-LANE-RED` / `HEAVY-RUN-LANE-SKIPPED` findings, made the isolated browser and JavaScript lanes announce `SKIP:` and exit 0 when their engine is absent, and replaced the heavy-lane permission prompt in work.md, clarify.md, work-reference.md, and the board timeline label with the exhaustion-time drain. Merge range `8e4089ab..4a28946e` (builder commits a673878e, 3f10d548, 19916bee, 5d33c582).

## Qualification

Passed — 15 files verified in the merge range `8e4089ab..4a28946e` (2 new, 13 modified; `qualify.sh` and `scope-drift.sh` run with `DO_WORK_DIFF_RANGE`), 5 requirements traced (wall seconds in evidence and result; runner with exit/skip/wall time, typed findings, unknown-lane and dirty-tree refusals, timeout; drain replaces the prompt in work.md with clarify as the by-hand path; isolated lanes announce `SKIP:`; skipped lane keeps the REQ held), P-A-U confirmed from the hand-back and mirrored above. Script notes judged: the two `(new)` Go files show no filename reference because Go references packages by symbol — `RunLanes` is called from `heavy_commands.go`, so neither is dead code; the scope-drift "declared but never touched" tokens are backticked identifiers inside the Scope descriptions, not files, and every declared path was touched. Orchestrator debug-artifact grep over added lines: none. Contamination check against the previous REQ in this session: not applicable (first REQ this run).

## Testing

**Tests run:** `go test ./internal/publication ./internal/heavyverification ./internal/resultmodel` (30 s, orchestrator rerun on the merged tree; builder run 23 s); `go test ./...` in the board module (41 s); builder gates on its branch: `go vet ./...`, `gofmt -l .`, `GOOS=windows GOARCH=amd64 go build ./...`, `go test -count=1 ./...` (49 s, 28 packages), `bash _dev/tests/maintainer-verify.sh --self-test` (2 s), `bash _dev/tests/action-shell-blocks.sh` (4 s), `bash _dev/tests/contract-regressions.sh` (18 s); canonical repository gate `bash _dev/tests/maintainer-verify.sh` unpiped from the project root on the merged tree: exit 0, recorded through `record-green-gate` at `4a28946e`.
**Result:** ✓ All passing. Slowest test file repository-wide is the pre-existing `internal/publication/defer_gate_test.go` at 21 to 26 s; every file this REQ wrote or changed is under 4 s.

**Red-green validation:** (traced to `## Red-Green Proof`; RED observed on assertions before each change, builder hand-back holds the failure output)
- `TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure` (tightened to `exit 0, 26s`): ✗ before implementation → ✓ after
- `TestHeavyTestingAnswerRejectsNegativeWallSeconds`: ✗ before → ✓ after
- `TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane` (`answered` half renders `skipped, 1s`; `confirmed` half is a lock-in that already refused): ✗ before → ✓ after
- `TestRunLanesRecordsExitStatusAndWallSeconds`, `TestRunLanesMarksSkipFromExplicitSkipLine`, `TestRunLanesRefusesUnknownLaneBeforeExecuting`, `TestRunLanesRefusesDirtyTrackedTree` (both halves), `TestRunLanesTerminatesTimedOutLaneGroup`: ✗ against stubs → ✓ after
- `TestRunHeavyVerificationHandlerRejectsMissingLane`, `TestHandlersRegisterRunHeavyVerification`: ✗ before → ✓ after
- `TestHeavyVerificationRunTextAndJSONCarryTheSameTypedLanes`: ✗ before → ✓ after
- `maintainer-verify.sh --self-test` isolated-lane skip case: ✗ with the guards removed → ✓ restored; confirmed end to end on this machine (no Chrome): `--heavy-lane queue-kanban-browser` prints the SKIP line and exits 0
- The GREEN condition's remaining half (the drain answers a held REQ without a prompt and returns it with `resume_phase: review`) is exercised on this REQ itself at queue exhaustion; see `## Heavy Verification Result` when recorded.

**New tests added:**
- `internal/heavyverification/heavy_run_test.go` (five runner tests, unix build tag), two handler tests in `heavy_commands_test.go`, one parity test in `resultmodel/result_model_test.go`, two answer tests in `publication/answer_test.go`, one self-test case in `_dev/tests/maintainer-verify.sh`

**Heavy verification plan:** *(all six lanes selected: the gate script is in every lane's coverage)*
- Range: 8e4089ab..4a28946e
- queue-kanban-javascript: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; `skills/do-work-board/tools/queue-kanban/timeline.go` matched the board subtree
- queue-kanban-browser: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; `skills/do-work-board/tools/queue-kanban/timeline.go` matched the board subtree
- do-work-cli-integrations: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; ten files matched the `skills/do-work/tools/do-work-cli` subtree
- staged-skills: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; fifteen files matched the `skills` subtree
- updater: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; ten files matched the `skills/do-work/tools/do-work-cli` subtree
- installer: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; ten files matched the `skills/do-work/tools/do-work-cli` subtree

*Verified by work action*

## Heavy Verification Plan

Mode: `exact-revision`
Base revision: `8e4089ab`
Target revision: `4a28946e782cf42397a329ae65b51ad3d5694a1b`
Manifest: `_dev/tests/heavy-lanes.json`

- queue-kanban-javascript: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; `skills/do-work-board/tools/queue-kanban/timeline.go` matched the board subtree
- queue-kanban-browser: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; `skills/do-work-board/tools/queue-kanban/timeline.go` matched the board subtree
- do-work-cli-integrations: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; ten files matched the `skills/do-work/tools/do-work-cli` subtree
- staged-skills: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; fifteen files matched the `skills` subtree
- updater: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; ten files matched the `skills/do-work/tools/do-work-cli` subtree
- installer: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — `_dev/tests/maintainer-verify.sh` matched its exact coverage path; ten files matched the `skills/do-work/tools/do-work-cli` subtree

## Open Questions

- [x] Run the selected heavy lane commands at `4a28946e782cf42397a329ae65b51ad3d5694a1b`; did every command exit 0? → Confirmed: every selected heavy lane exited 0 at 58e1c9c948bb68f3805e704b9c7db39fff38f504 during the queue-exhaustion drain (queue-kanban-javascript 6s, queue-kanban-browser 82s, do-work-cli-integrations 67s, staged-skills 25s, updater 52s, installer 24s)
  Recommended: Yes
  Also: No — <failing lane>


## Answer Notes

- 2026-09-04 - [ ] Run the selected heavy lane commands at `4a28946e782cf42397a329ae65b51ad3d5694a1b`; did every command exit 0?: Confirmed: every selected heavy lane exited 0 at 58e1c9c948bb68f3805e704b9c7db39fff38f504 during the queue-exhaustion drain (queue-kanban-javascript 6s, queue-kanban-browser 82s, do-work-cli-integrations 67s, staged-skills 25s, updater 52s, installer 24s)

## Review

**Overall: 95%** | 2026-09-04T14:23:10Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `skills/do-work/actions/roadmap.md:68` still tells the roadmap action to say that `do-work clarify` owns "the one permission prompt" for a `pending-heavy-testing` REQ; that prompt no longer exists — impact-user-visible → report only

**Minor findings:**
- `skills/do-work/actions/restart-with-parallel-handoff.md:64` describes clarify as the way "to authorize held heavy tests"; the routing is right, the verb is stale — impact-user-visible → report only
- `skills/do-work-board/actions/board.md:92` lists "heavy-test permission" among what makes a Needs-input card operator-actionable — impact-negligible → report only
- Five maintainer probe scripts print "run `_dev/tests/maintainer-verify.sh --heavy` after user permission" (`_dev/tests/update-script-behavior.sh:6`, `prescribed-shell-harness.sh:12`, `install-suite-behavior.sh:5`, `prescribed-shell-scripts-behavior.sh:9`, `staged-skills-contract.sh:5`) — impact-negligible → report only
- An interrupted run breaks the lane loop and still returns `outcome: success` with no typed finding that the run was cut short (`heavy_run.go:117`); the drain's "present in the run" condition fails closed, but a `HEAVY-RUN-INTERRUPTED` finding would make that mechanical — impact-rule-change → report only
- The rename and copy branch of `refuseDirtyTrackedTree` (`heavy_run.go:141-148`) has no test case; both halves of `TestRunLanesRefusesDirtyTrackedTree` use plain modified paths — impact-negligible → report only
- Nit: `WallSeconds` truncates, so a sub-second lane records `0s`, indistinguishable from an unset field — impact-negligible → report only
- Nit: `defaultLaneTimeoutSeconds` sits apart from the const block its siblings share (`heavy_run.go:314`) — impact-negligible → report only

**Acceptance:** Pass — focused Go tests and `go vet` green; self-test with the two new isolated-lane skip cases green; the isolated browser lane with no `QUEUE_KANBAN_BROWSER` printed its `SKIP:` line and exited 0 on this Chrome-less PATH; `--lane nope` returned `HEAVY-RUN-LANE-UNKNOWN` and a duplicated `--lane` returned `HEAVY-RUN-USAGE`; the recorded drain at `58e1c9c9` covers the full heavy tier through the new command.
**Suggested testing:** 4 items — a rename case for the dirty-tree test; a multi-REQ drain with differing lane sets; an interrupted mid-lane drain; a drain on a checkout dirty outside `do-work/` to see the refusal reach the exit summary.
**Follow-ups created:** None (8 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Splitting the runner from planning and answering: `run-heavy-verification` only executes and times named lanes, `plan-heavy-verification` stays the one planner, and clarify's existing `answer` manifest records the result, so batching across held REQs fell out of the split instead of needing a super-command.
- Dogfooding the drain on this REQ at its own queue exhaustion: the hold, the typed `HEAVY-RUN-LANE-SKIPPED` finding, the browser remedy, the green answer, and the `resume_phase: review` resume all ran unattended and matched the prose.

**What didn't:**
- The first drain skipped the browser lane because no Chrome is on PATH; the morning batch only passed that lane because a person had set `QUEUE_KANBAN_BROWSER`. A skip that surfaces as a typed finding with a named remedy is the right outcome, but the gate runner and any headless drain should export that variable, or the browser lane skips on every unattended run.
- The hold leaves `claimed_at` in place until the answer clears it, so the selector reports a held REQ as `ALREADY-CLAIMED` rather than as heavy-held. Harmless for the drain, misleading in the replay output.

**Worth knowing:**
- The runner keys `skipped` on a `SKIP:` output line, never on a lane name, and the isolated browser and JavaScript lanes now print that line at their own source (`_dev/tests/maintainer-verify.sh`). A lane that skips silently inside its Go test still reads as red, so any new engine-gated lane must announce its skip the same way.
- `wall_seconds` travels under one JSON tag through `resultmodel.HeavyLaneExecution` and `publication.HeavyLaneResult`, so the runner's `lanes` array pastes straight into `heavy_testing.lanes`.
- The dirty-tree refusal ignores `do-work/` on purpose (D-01): the orchestrator's living-log edits are always uncommitted when a drain runs.

## Orientation

Now the work loop finishes held heavy verification on its own: at queue exhaustion it runs the selected lanes through `run-heavy-verification`, records exit status and wall seconds per lane, and returns green REQs for review without a human answering a prompt. Lives in the do-work-cli `heavyverification` and `publication` packages plus the work and clarify actions. [MAP CHANGED] The heavy-test hold is now a non-blocking pause the loop resolves itself; the permission gate is gone and `do-work clarify` is the by-hand fallback only. Prime spot-check: `prime-do-work-cli.md`'s `heavyverification`, `ownedprocess`, and package-direction sentences were updated in this REQ; `prime-action-files.md` references still resolve.
