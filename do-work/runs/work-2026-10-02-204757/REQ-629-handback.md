# REQ-629 Builder Hand-Back

**Builder commit:** `42c37e7e` on `worktree-agent-REQ-629-handoff-resume-keeps-claims` (worktree `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/REQ-629`), parent `89145e03`.

## File Manifest

- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (modified): the `RECOVERY-TAKEOVER-AVAILABLE` finding's `next_argv` is now `do-work-cli --format json advance REQ-NNN` (same shape as the heavy-lanes hold). The `automation_stop_reason` says: "working claim preserved; continue it with advance REQ-NNN, or, only for a claim no live session owns, recover --take-over REQ-NNN resets it (requeues it as pending and strips its orchestrator sections)". Code, severity, decision, verification argv and `--take-over` itself are unchanged. There is no writer comparison.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modified): `TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim` now decodes the finding and pins the read-only `advance` next argv, plus a stop reason that names `recover --take-over REQ-713` and "resets". Its take-over half still runs. It adds helper `recoveryTakeoverFinding` and the new test `TestRecoverNextStepContinuesAClaimWithoutResettingIt`.
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modified): Step 3 records each claim's continue command (bare `advance REQ-NNN`) instead of the "exact takeover command". The Step 4 paste block adds one `advance REQ-NNN` line per working claim and says `recover --take-over` resets a claim. The sufficiency sentence becomes plural. The checklist line was updated to match.
- `skills/do-work/actions/work-reference.md` (modified): Crash Recovery now says the takeover finding's next step is the read-only `advance REQ-NNN`. It also says the reset requeues the claim and strips its sections, and is only for a claim no live session owns, never for resuming a handoff.

## P-A-U

- **PLAN:** Use Option A from the exploration. Change only the not-authorized branch's next argv and stop reason. Pin it with TDD in the existing test plus one UR-132 test. Then align the two prose files and grep every caller first (lesson family rule-direction-checked-against-callers).
- **APPLY:** Made exactly the four Scope files. The tests were written and seen to fail first, then the Go change was made, then the prose.
- **UNIFY:** `git diff --stat`: 4 files, 93 insertions, 10 deletions (handoff md +8/-2 lines changed, work-reference md 1 line, recovery_commands.go 11, recovery_commands_test.go 82). Checks per file:
  - recovery_commands.go: gofmt clean, `go vet ./...` clean, no debug output.
  - recovery_commands_test.go: `go test -count=1 ./internal/lifecycleadvance/` ok in 24.2s (under 30s). `./internal/resultmodel/` and `./internal/doctor/` are also ok.
  - restart-with-parallel-handoff.md: read the full diff. No test pins its wording (grep of `_dev/` and `skills/`).
  - work-reference.md: `_dev/tests/contracts/recovery-set-aside.sh` and `_dev/tests/contracts/core-checks.sh` both pass.

## Red-Green Evidence

- `TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim`. RED: `takeover finding next argv = []string{"do-work-cli", "recover", "--take-over", "REQ-713"}, want read-only advance`. GREEN: pass.
- `TestRecoverNextStepContinuesAClaimWithoutResettingIt`. RED: following the offered next argv ran `recover --take-over REQ-714`. The output showed `RECOVERY-CLAIM-RECOVERED`, the file moved `working/ -> queue/`, and a commit was made. That is the UR-132 failure reproduced. The test then failed on decoding the text-mode output. GREEN: pass. The claim stays byte-identical in `working/`, the tree digest is unchanged, and the classifier names `preflight`.

## Scratch-Repo GREEN Run

Directory: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f3fafe8a-e2c7-4452-928a-c2f42556a5bd/scratchpad/req629-green`. The CLI was built from the worktree with `go build -o <scratch>/do-work-cli ./cmd/do-work-cli`.

1. Fixture setup: queue REQ-001, then `claim REQ-001 --writer "$(hostname):$PWD"`. Then added route B, write_set, `commit: abc1234`, estimate, and Triage/Plan/Exploration/Scope. Committed it with a RESTART-PROMPT paste block written to the new rules (`do-work run --fan-out 1`, `advance REQ-001`).
2. `do-work-cli --format json recover` gave: outcome success, `FINALIZATION-NONE`, `RECOVERY-TAKEOVER-AVAILABLE` with next_argv `['do-work-cli','--format','json','advance','REQ-001']`, and the stop reason naming the reset.
3. Following that next_argv / the paste line `advance REQ-001` gave: success, tree_section `working`, phase `preflight`.
4. Queue-mode `advance` (the resume command's selection) gave: success, no findings, no mutation.
5. `shasum -c` on the REQ file reported OK. `git status` was clean. `do-work/queue/` was empty and `do-work/working/REQ-001-sample.md` was present.

## Decisions

- **D-02:** The paste-block sufficiency sentence changes from "This command is sufficient; everything below it is context." to "These commands are sufficient; everything below them is context." Reasoning: with claims, the block carries more than one command. Selection never continues a working claim, so the old singular sentence was not honestly true whenever a claim existed. Only an old `ai-reports/` evidence file quotes the old sentence, and no test pins it.
- **D-03:** The new test's fixture carries `commit: abc1234` and a same-host writer label. Reasoning: this models the UR-132 case (a merged claim resumed on the same machine). The finding is writer-agnostic, so the label is realism only. The expected phase is `preflight`, the real next phase for a Route B claim with Scope and write_set written.

No ESCALATE decisions.

## Discovered Tasks

- `skills/do-work/docs/work-guide.md:165` (Context limits): "Start a new session with `do-work run`; its first canonical recovery result says whether anything needs takeover authority." This steers a resuming session toward takeover. It also misses that `do-work run` never continues a working claim (it should say `advance REQ-NNN` per claim). Not edited because it is out of scope.
- `skills/do-work/docs/work-guide.md:130`: "Plain `recover` preserves working claims and returns typed takeover options". This is now slightly stale, because the offered next step is `advance`. Suggested fix: name `advance REQ-NNN` as the next step and `--take-over` as a reset. Not edited.
- `run-with-recovery.md:41` and `forensics.md:54/100/168` describe genuine crash recovery or delegate to Crash Recovery. They do not tell a resuming session to take over. No action needed.

## Lessons

- Read: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `[family: rule-direction-checked-against-callers]` (one bullet, REQ-604 / 0.305.30). Applied: grepped `RECOVERY-TAKEOVER-AVAILABLE`, `take-over` and `takeover` across `skills/` and `_dev/tests/`. The only code caller pinning the argv was the updated test. The prose callers are listed above. The resultmodel text-render test covers the take-over mode only and still passes.
- Missing: none.
