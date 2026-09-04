# Builder Brief — REQ-566

You are the implementation builder for do-work REQ-566. You work ONLY inside this worktree and commit ONLY on its branch:

- Worktree (your working directory): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-566-run-held-heavy-lanes-at-exhaustion-and-record-lane-wall-time`
- Branch: `worktree-agent-REQ-566-run-held-heavy-lanes-at-exhaustion-and-record-lane-wall-time` (already checked out there, based on main `a37ea0d8`)
- Go module: `<worktree>/skills/do-work/tools/do-work-cli` (Go 1.26). Node v22 is present. No Chrome on PATH.

## Never touch
- Anything under `do-work/` in the worktree. It is a stale snapshot; never read status from it, never write it, never commit it. The orchestrator owns queue state in the main tree.
- The main tree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` — with exactly one exception: your hand-back file (below).
- Release files: `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`. The orchestrator writes the release at finalization.
- `_dev/tests/heavy-lanes.json` (argv unchanged), `cmd/do-work-cli/main.go` (already ranges over the package handlers).

## Hand-back (mandatory, the only main-tree path you may write)
Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-132415/REQ-566-handback.md` (absolute path, never staged or committed) with these sections:
1. `## Branch` — the branch name and every commit hash on it.
2. `## File Manifest` — every source file created/modified/deleted with the verb, tests included.
3. `## AI Execution State (P-A-U Loop)` — the three checkboxes with your notes ([PLAN] brief approach, [APPLY] scope kept, [UNIFY] `git diff --stat` reviewed, `gofmt`/`go vet` run, no debug artifacts, list each file checked).
4. `## Testing` — every command run with elapsed seconds per test file/command, RED evidence (test name, failure output before the change) and GREEN evidence (pass after) for every RED test in the plan, plus gate results.
5. `## Decisions` — numbered `D-NN` entries. Record the orchestrator's directive as D-01 ("dirty-tree refusal ignores paths under `do-work/`; dirt outside `do-work/` still refuses") and continue from D-02. Each entry: reasoning; ESCALATE-tier entries add `Value:` and `Risk:`.
6. `## Discovered Tasks` — out-of-scope finds, never fixed inline (or "None").
7. `## Integration Seams` — any line the orchestrator must apply in a shared main-tree file (or "None").
8. `## Lesson Evidence` — the `required_lessons` entries read (none are stamped on this REQ; both index matches were over budget — see the REQ's Dropped section), plus the prime `## Lessons` satellites you read because the touched code is named by a prime's Read-first or Traps entries.

## Rules that govern you (read before coding, in this order)
- `<worktree>/skills/do-work/crew-members/general.md`
- `<worktree>/skills/do-work/crew-members/coding-guardrails.md`
- `<worktree>/skills/do-work/crew-members/communication-style.md`
- `<worktree>/skills/do-work/crew-members/backend.md`
- `<worktree>/skills/do-work/crew-members/testing.md` (tdd: true)
- Prime files, read in full before touching code: `<worktree>/_dev/primes/prime-action-files.md`, `<worktree>/skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, and their `## Lessons` satellites `_dev/primes/lessons-action-files.md`, `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`. Also `<worktree>/_dev/primes/prime-shell-commands.md` before editing `_dev/tests/maintainer-verify.sh` or any fenced command block in an action file, and `<worktree>/_dev/primes/prime-kanban-board.md` before touching `skills/do-work-board/tools/queue-kanban/timeline.go`.
- TDD is mandatory: for every RED test named in the plan, write it first, run it, capture the failing output, then implement, then capture the pass. Anchor on the REQ's `## Red-Green Proof`.
- Write only inside `## Scope` → "Files I will touch". Needing another file: if the REQ's own requirements already require that file class, flag it and proceed; otherwise stop and report in the hand-back, never write it silently.
- Implement according to `## Plan` below. Deviations are `D-NN` decisions with reasoning.
- Naming for reach: two-word minimum for exported identifiers, struct fields, files, CLI flags, finding codes.
- Commit on your branch as you go (small complete slices are fine, one per task is ideal). Never `git add -A` across `do-work/`. Never push.
- Gates before hand-back (all from the worktree root unless noted): `go test ./internal/publication ./internal/heavyverification ./internal/resultmodel` (from the Go module), `go vet ./...`, `gofmt -l .` (module), `GOOS=windows GOARCH=amd64 go build ./...` (module), `go test -count=1 ./...` (module), `bash _dev/tests/maintainer-verify.sh --self-test`, `bash _dev/tests/action-shell-blocks.sh`, `bash _dev/tests/contract-regressions.sh`, then the unpiped `bash _dev/tests/maintainer-verify.sh` from the worktree root. Record elapsed seconds; any single test file at or above 30 s is a test-maintenance failure to fix.
- Exploration notes with exact patterns to copy: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-132415/REQ-566-exploration.md` (absolute main-tree path; read-only).

---

# The REQ (verbatim from the main tree at dispatch time)

---
id: REQ-566
title: '[impact-rule-change] Run held heavy lanes at queue exhaustion without asking, and record per-lane wall time'
status: claimed
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
claimed_at: 2026-09-04T13:24:15Z
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
---

# Run Held Heavy Lanes At Queue Exhaustion Without Asking, And Record Per-Lane Wall Time

## What

Two changes, in this order.

1. **Record per-lane wall time.** The generated `## Heavy Verification Result` section records wall seconds for every lane next to its exit status, and the typed `heavy_testing` evidence (`HeavyLaneResult` in `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go`) carries that field so the result is machine-readable.
2. **Lift the human permission prompt.** At queue exhaustion, the work loop runs the batched heavy plan itself at HEAD for every `pending-heavy-testing` REQ, then records green or red through the same `answer --manifest` transaction that `do-work clarify` uses today. The `pending-heavy-testing` status, the non-blocking hold, and the source-ready rule for dependents all stay. Only the permission question goes. `do-work clarify` keeps working for a maintainer who runs the lanes by hand.

The fold-first scan found no pending or pending-answers REQ in any UR that shares this root cause. REQ-564 (reuse matching per-lane verification evidence for four hours) avoids rerunning unaffected lanes and REQ-559 (retry a red repository gate once before a repair REQ) handles a red gate; neither removes the permission prompt or records lane durations.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
