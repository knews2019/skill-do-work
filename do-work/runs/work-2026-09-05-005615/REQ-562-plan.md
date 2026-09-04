# REQ-562 — Lightweight serial timing plan

**Route C proposed.** Add one small timing owner, connect it to the existing action's real stage/command boundaries, and seal its compact summary through the existing finalization transaction. The queued REQ's September 4 addendum supersedes UR-108's earlier critical-path proposal. No hierarchy, board changes, retry policy, background scheduling, model telemetry, or instrumentation-only commit belongs here.

Read-only prep based on the current queued REQ, UR-108, REQ-448 and REQ-539 records, action/shell/CLI primes and relevant shared/general/testing/coding/security guidance. No source or queue changes and no tests were run. Re-resolve source owners, request identity, and scope at claim; ongoing REQ-504/506 work and unrelated publication edits may move these files.

## Existing boundaries and decisions

- REQ-448's eight optional successful-event milestones are still observations in `work.md`/`work-reference.md`; leave their schema, recovery behavior, board readers, and `claimed_at → completed_at` calibration alone. They are not precise stage starts and cannot be retroactively converted into new timing events.
- REQ-539 already owns `_dev/tests/test-duration-log.sh` and the TSV shape `run_id`, `file`, `seconds`, `other_gate_processes`. Gate callers can supply `DO_WORK_TEST_RUN_ID` and `DO_WORK_TEST_DURATION_LOG`. Ingest matching existing rows as bounded diagnostic detail; never create another persistent per-test timer or add per-test seconds to stage totals.
- `ownedprocess` already owns process groups and teardown. A timing command wrapper should reuse it and inherit stdio directly; it must not buffer/log output, change the child argv, add retries, detach a gate, or introduce a resident timer. Keep the canonical gate's direct, unpiped execution and raw exit status.
- Git-private shared state must resolve `git rev-parse --git-common-dir` relative to the selected checkout, as `gateevidence` already does. A worktree's `.git` is a file and its private git-dir is not the shared repository store.
- **One flat run stream, keyed records:** a bounded JSONL stream under `<git-common-dir>/do-work-timing/<opaque-run-id>/events.jsonl`, with run, REQ, event identity, operation/category, UTC start/end, elapsed seconds, outcome, optional revision/responsible identity, and optional command status. Begin/end state is private and owned by the same package. Protect concurrent appends and completion/cleanup with one narrow store serialization primitive; test independent processes/worktrees rather than assuming multiple buffered writes are atomic. Do not introduce a general lock manager or silently truncate another writer's tail.
- **Safe identity:** persist a static operation key such as `repository-gate`, `handback-merge`, or `heavy-verification`, and a validated executable family where useful. Never persist raw argv, shell text, command output, arbitrary labels, environment values, URL credentials, or a hash of raw secret-bearing argv. Pass the real command separately for execution; only the allowlisted metadata reaches JSONL or the archive.
- **Flat accounting:** serial stage events supply category totals; command and delegated-wait events supply slowest-detail attribution and are not added again to total stage time. There is no parent link or nesting calculation. Total observed wall window is first observed boundary to the summary cutoff; unattributed time is its residual after the serial stage ledger, with malformed/inconsistent timing explicitly incomplete rather than fabricated. Missing categories are absent, not zero. Concurrent runs/agents are separable by identity; this delivery makes no parallelism claim.
- **Honest clocks:** an in-process wrapper uses `time.Now`/monotonic subtraction through an injectable clock. A stage beginning in one CLI process and ending in another has UTC wall elapsed with an explicit wall-clock basis, not invented monotonic precision. Clock reversal produces incomplete evidence. Agent waits mean observed dispatch-to-handback or explicit wait intervals, not model-internal compute time.
- **Necessary observation cutoff:** the existing finalization commit cannot contain the measured duration of that same commit or the cleanup that follows it. Seal a summary immediately before its primary commit, identifying `observed_through` and a finalization-to-cutoff interval. Later commit/verification/raw deletion/worktree cleanup is explicitly outside that archived observation window; do not call it zero or create an additional timing commit to measure it. Cleanup operations already completed earlier can appear normally. This limitation should be recorded as a D-XX decision at claim and stated in `## Timing`; if full post-commit-tail accounting is required, the captured no-extra-commit constraint must change rather than being silently violated.

## Proposed exact scope

New timing owner (four files):

- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_stream.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_stream_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands_test.go`

Production registration and consumers (five files):

- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/docs/command-line-guide.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`

Existing finalization integration and tests (six files):

- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_journal.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_timing.go` (new)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_timing_test.go` (new)

This 15-path declaration is a proposal, not authority to write. Inspect the exact finalization preimage/application seams before dispatch; a required path absent from this set needs an explicit scope adjustment. No `_dev` duration writer, board file, request milestone schema, lifecycle selector, process-group implementation, or commandruntime change is proposed.

## Task 1 — Build the single timing owner with synthetic-clock RED proof

Start with the captured synthetic serial Route C lifecycle: recovery/selection, planning, exploration/preflight, builder work and an observed delegated wait, hand-back/merge, verification/gate, review, remediation when present, finalization-to-cutoff, and any completed cleanup stage. The failing fixture must assert exact flat events, stage/category totals, slowest stage and command, and unattributed seconds, not merely that a helper exists.

Implement bounded metadata validation, begin/end correlation for cross-process stages, same-process elapsed measurement for commands, append-safe shared-worktree storage, and a pure summary function with deterministic ordering/tie handling. Use stable opaque IDs and known operation/category values; do not accept a user-authored command label as safe metadata. Tests inject separate wall/elapsed clocks, cover gaps, clock reversal, missing endpoints, duplicate/replayed closure, multiple REQs/runs, concurrent append integrity, and secret sentinels in actual child arguments that must never appear in stored/returned timing. A missing timing stream remains a no-op.

Expose only the small primitives the action needs: begin/end a named stage or observed wait, and execute one material external command while timing it. A command wrapper inherits stdin/stdout/stderr and returns the child's raw status; in text passthrough mode its own result must not contaminate the gate's output. Use the existing handler registration pattern and ownedprocess lifecycle. Do not add a worker daemon, scheduler, general tracing API, or a timer around every shell primitive.

## Task 2 — Connect real execution, not just a standalone helper

Update `work.md` at the existing transitions to invoke the writer in the **same tool invocation** already doing the work or persisting its boundary. Use the existing run identity, pass it explicitly with REQ identity to delegated builders, and carry responsible identity only where known. Metadata calls must not require a new agent turn, a tracked timing artifact, or a timing-only Git commit. Keep principles in the action and the compact data/command contract in the reference and CLI guide.

The wiring acceptance ledger must name actual callers:

- recovery/selection around the canonical queue-boundary calls;
- planning/exploration and review around dispatch/handback, with observed waits identified distinctly;
- preflight and focused evidence around their existing advance invocation;
- builder work from accepted dispatch through completed handback, not from a milestone guessed afterward;
- hand-back/merge, direct repository-gate attempts, selected heavy runs, and material cleanup commands through the timing wrapper;
- remediation uses another flat event in its stable category, keeping the existing retry owner unchanged;
- finalization consumes the same context inside its canonical owner, not through an action-authored `## Timing` body.

The wrapper around a real direct gate can pass a matching test run ID through the existing REQ-539 interface and ingest only that run's rows for slow-test detail. It does not time every test itself and does not count test totals again. Existing background-gate and retry behavior is observed where available; no attempt is made to replace REQ-542/559. Add public CLI fixture coverage that follows this caller sequence and checks exit/stdio parity, real event emission, and unchanged failure outcomes. A helper-only unit test does not complete this task.

## Task 3 — Seal Timing through existing finalization and clean only consumed evidence

The finalizer currently freezes lifecycle/release preimages and postimages, applies them, commits exact paths, records provenance if needed, verifies, and deletes its journal. Integrate Timing into that same image/journal authority before the primary commit identity is frozen. Freeze the run/REQ identity, observation cutoff, exact bounded summary bytes, and consumed event identities in optional journal data; recovery must replay those same bytes rather than recompute a different duration. Archive summary changes must participate in the existing conflict checks, rollback images, exact commit allowlist, and final-state verification. Do not patch the working REQ behind its manifest digest or append an unjournaled archive mutation after the commit.

Emit one compact `## Timing` block with observed wall seconds, serial stage/category totals, the slowest stage/command, unattributed seconds, measurement basis/coverage, and the cutoff limitation. Preserve milestone bytes and calibration. Optional slow-test detail comes from matching REQ-539 rows; no copied raw log is committed. A historical/no-stream request follows its existing byte and commit path without an empty Timing block.

Only after the finalizer proves its successful committed state may the canonical timing owner remove that REQ's consumed events and temporary endpoints. Since one run may still contain other REQs, prune only the exact consumed identities and delete the raw stream/directory when empty; never delete another REQ's or concurrently appended event. On refusal/interruption, retain raw evidence and journal state for retry. Recovery and already-completed replay must not append a second summary or recreate a cleaned stream.

Synthetic-clock finalization tests must assert exact summary values, unchanged milestone/calibration fields, no-trace byte parity, raw cleanup, another REQ/run preserved, interrupted replay at relevant existing phase seams, and **the same Git commit count** as the corresponding uninstrumented transaction. Include the public registered command path plus an isolated linked-worktree storage test. Then run the affected package tests, applicable vet/race checks, and the project's canonical merged gate; the orchestrator owns renewed heavy planning. Do not broaden to board testing solely because REQ-448 once changed the board.

## Claim-time acceptance decision

The compact pre-commit cutoff is the only implementable default that respects both real measured durations and zero instrumentation commits. Review must not describe it as full end-to-end finalization/cleanup timing. If the maintainer expects those later durations inside the same immutable commit, report the conflicting requirements before source work; a new tracing system does not solve that ordering problem.
