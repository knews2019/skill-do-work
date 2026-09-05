---
id: REQ-562
title: 'Addendum: record lightweight per-REQ lifecycle timings'
status: completed
created_at: 2026-09-03T21:28:31Z
user_request: UR-108
addendum_to: REQ-448
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-448, REQ-531, REQ-539, REQ-542, REQ-559]
claimed_at: 2026-09-05T00:33:24Z
route: C
estimate:
  p50_active_minutes: 45
  confidence: low
  calculated_at: 2026-09-05T08:13:09Z
  basis:
    - Route C
    - 12-file write set
    - 3 subsystems involved
    - 5 acceptance criteria
planning_at: 2026-09-05T01:35:32Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing.go
  - skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands_test.go
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/docs/prescribed-shell-primitives.md
completed_at: 2026-09-05T08:15:45Z
commit: 4c2cd380
release_at: 2026-09-05T08:15:45Z
---

# Addendum: Record Lightweight Per-REQ Lifecycle Timings

## What

Extend REQ-448's phase milestones with a lightweight, flat timing stream so a completed `do-work run` attributes one REQ's elapsed time to its major workflow stages and material external commands. Keep the raw stream in Git-private local state, create no instrumentation commits, and fold one concise timing summary into the REQ during finalization.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the request, both required primes, the action-file lessons, the guardrails and testing crew files, then read REQ-448's implementation commit `01a1676`, the current 823-line `work-reference.md`, the CLI prime, the finalization and gate-evidence packages, and the existing `do-work/test-durations.tsv` writer before writing a line. Chose one new package plus three argv verbs, and rejected a nesting model, a critical path and a second test-duration system in writing (D-01, D-02, decisions above).
- [x] **[APPLY]:** Twelve tests written against real stub signatures failed on assertions first (11 red, output quoted above), then the writer, the command layer and the prose landed. Scope stayed inside the derived boundary: one new package, one result-model extension, one command registration, one recovery-symmetry line, and four documentation surfaces.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file before staging. `gofmt -l .` printed nothing, `go vet ./...` exited 0, `go test -count=1 ./...` is green everywhere except the pre-existing `heavyverification` runtime probe proven red on a stashed tree, `core-checks.sh` and `shipped-package-reference-contract.sh` passed, and the fast contract aggregate reported only a load-induced duration breach. No debug print, no stray file, nothing under `do-work/` staged.

## Why

REQ-531 took 1 hour 27 minutes, but its durable evidence could attribute that time only to broad milestone intervals. The maintainer needs command- and wait-level evidence to decide whether gates, agent latency, merge conflicts, retries, or lifecycle bookkeeping are the highest-value optimization targets.

## Context

REQ-448 already records optional flat timestamps for planning, dispatch, builder handback, integration, review, remediation, re-review, and release, renders a phase breakdown, preserves `claimed_at` to `completed_at` calibration, and keeps historical REQs compatible. This addendum deepens that observation model only enough to answer how long a serial REQ took and why; it does not discard those fields or introduce a general tracing system.

Related pending work owns adjacent optimizations: REQ-539 records per-test durations, REQ-542 moves gate execution into a background runner, and REQ-559 retries a transiently red repository gate before creating repair work. This REQ measures those paths and consumes their evidence where available; it does not reimplement them.

## Prior Implementation

REQ-448 added eight optional successful-event milestone fields across the canonical lifecycle schema, work action, board parser, duration model, generated payload, and card/detail UI. Recovery strips stale phase observations, absent fields remain backwards-compatible, and claim-to-completion stays the calibration span. Its implementation commit is `01a167697743c505e0b77e21a83db23fdc255cc0`.

Key implementation surfaces were `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, and the queue-kanban model, duration, generator, test, and card/detail files under `skills/do-work-board/tools/queue-kanban/`.

## Detailed Requirements

- Persist one flat, append-safe timing stream per run in Git-private state shared by the repository's worktrees, associated with the REQ and run identity. JSONL is an acceptable initial representation.
- Record the existing major lifecycle boundaries under stable categories: recovery/selection, planning, exploration/preflight, builder work, hand-back/merge, verification/gate, review, remediation, finalization, and cleanup. Record delegated-agent waits and canonical external commands that materially contribute wall time; do not trace file reads, individual shell primitives, or nested implementation details.
- Every event records its operation/category, wall-clock start and end, elapsed duration measured in-process when one process observes both ends, outcome, relevant revision when known, and responsible agent or process when known. Command events also record exit status and a safely redacted command identity. No parent span or nesting model is required.
- At finalization, calculate total observed elapsed time, time by category, the slowest stage and command events, and unattributed wall time. Do not calculate exclusive time, a critical path, overlap, or parallelism savings.
- Fold one compact, structured `## Timing` summary into the archived REQ in the existing finalization commit, then remove the Git-private raw stream. Instrumentation must not create an agent turn, tracked run artifact, or Git commit of its own.
- Reuse or ingest REQ-539's per-test duration evidence when available rather than introducing a second persistent test-duration system.
- Preserve REQ-448's milestone fields and `claimed_at` to `completed_at` calibration for compatibility. Historical REQs and runs with no structured timing stream continue to parse and display normally.
- Keep timing evidence bounded by recording metadata only, never command output or raw arguments that could contain secrets or user-controlled content.
- Use deterministic synthetic-clock fixtures for duration, category aggregation, command attribution, final-summary generation, raw-stream cleanup, and unattributed-time assertions.

## Constraints

- One canonical timing writer owns timestamp and duration mechanics; do not scatter `date` calls or independent timing formats through action prose and scripts.
- Use a monotonic clock where one process observes both ends of an event and UTC wall timestamps for human correlation. Do not add a resident process merely to manufacture cross-process monotonic spans.
- Concurrent runs and agents must remain separable by run, REQ, and span identity.
- A hierarchy, parent/child spans, exclusive-time accounting, critical-path analysis, parallelism metrics, Chrome Trace export, board visualization, and unavailable model-internal token/latency telemetry are explicitly out of scope.
- Test-duration ownership remains with REQ-539; retry policy remains with REQ-559; background gate scheduling remains with REQ-542.

## Builder Guidance

Certainty: Firm on the flat evidence and compact final summary, exploratory on the exact Git-private storage schema and canonical writer. Prefer one small executable primitive with deterministic tests over procedural timing prose. Optimize for zero additional agent turns and zero instrumentation-only commits.

## Red-Green Proof

**RED prompt/case:** Run a synthetic serial Route C lifecycle containing planning, exploration, a delegated builder wait, a timed gate command, review, and finalization. Today the completed REQ has milestone timestamps but no durable command timings, category totals, slowest-event report, or measure of unattributed wall time.
**Why RED now:** REQ-448 records broad phase boundaries only, so a long interval cannot distinguish agent work, material commands, waits, and lifecycle bookkeeping.
**GREEN when:** The deterministic fixture emits flat Git-private events and one final `## Timing` summary whose total observed time, category totals, slowest stage/command events, and unattributed time match the synthetic clock; finalization removes the raw stream without an instrumentation-only commit, and a historical no-trace fixture remains unchanged.
**Validation:** User narrowed the captured request after reviewing the cost evidence; serial per-REQ attribution is the goal, not general distributed tracing.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`; matched because this addendum extends lifecycle fields and action-owned evidence readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget and `slugged: partial`; matched because external command timing and a canonical execution primitive may touch shipped shell boundaries.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5924 tokens, over budget and `slugged: partial`; matched because the request adds structured evidence projection and lifecycle identity.

## Addendum (2026-09-04)

User added:

> ````text
> do it
> ````

- Narrow the requested instrumentation to the minimum needed to explain one REQ's serial wall time.
- Keep flat major-stage and material-command timings; write raw events to Git-private state and fold one summary into the REQ at finalization.
- Drop nested spans, parent relationships, exclusive-time and critical-path calculation, overlap and parallelism metrics, Chrome Trace output, and board visualization.
- Require zero extra agent turns, tracked run artifacts, or instrumentation-only commits.

Resolved conflict: the original request required a general nested tracing model with critical-path and parallelism analysis → the user chose lightweight serial per-REQ attribution with a flat timing stream and one final summary.

## Full Context

See `do-work/user-requests/UR-108/input.md` for the complete verbatim invocation and its capture-time summary. This addendum extends REQ-448 rather than replacing it.

---
*Source: "ok, capture it using do-work capture-request" — referring to the agreed structured lifecycle timing and critical-path proposal.*
## Triage

**Route: C** — the request spans a new package, three commands, prose in two action files and a shipped doc.

**Reasoning:** The change adds a new Go package to the shipped CLI, extends the shared result model, registers new command handlers and edits lifecycle prose that other steps read, so it is multi-component rather than one located edit. It also had to be reconciled against REQ-448 (the eight optional phase-milestone timestamps it extends) and against three neighbouring pending requests that own per-test durations, background gate scheduling and gate retries. That mix of new surface plus contested boundaries is what Route C is for.

**Planning:** Required. The plan is the request's own Detailed Requirements, which are unusually precise: they name the storage shape (one flat, append-safe, Git-private stream, JSONL accepted), the ten categories, the per-event fields, the four summary numbers, the fold point inside the existing finalization commit, and the explicit exclusions. Planning confirmed that boundary and chose the writer shape inside it instead of inventing scope.

## Plan

One new package holds the canonical timing writer: stream location, event record, validation, redaction, summary, section fold and stream removal. One run's events for one request live in a single append-only JSONL file under the Git common directory, so the stream is Git-private, shared by every linked worktree, and never reachable by the index. Three argv verbs sit on top of it: one closes an interval the caller already lived through, one runs and measures a command itself, one folds the summary and deletes the raw stream. A boundary event with no explicit start chains from the previous event's end, which is what lets a serial lifecycle be recorded without any caller holding its own timestamp; a wrapped command derives its start from the elapsed time the process measured, so waiting before the launch stays honestly unattributed. The fold runs at Step 8, before the finalization manifest is authored, so the section lands in the preimage the existing finalization transaction archives and commits. No new commit, no new agent turn, no tracked artifact. Around that: one result-model extension, one command registration, one recovery-symmetry line, and four documentation surfaces.

The plan deliberately excluded a whole class of analysis. No parent or child spans, no nesting, no exclusive time, no critical path, no overlap or parallelism metrics, no Chrome Trace export and no board visualization. The request's Detailed Requirements and Constraints forbid all of them, and a dated user decision in the request's own Addendum resolved the original nested-tracing idea down to lightweight serial per-request attribution. The filename still ends with the words about reporting the critical path, and the capture agent's restatement still promises one; the requirements win over both. Nesting also buys nothing here, because no event refers to another: totals are a sum over independent rows, and the gaps between rows are exactly the unattributed wall time the request asked for. Two further exclusions were deliberate: no second test-duration system, because per-test duration evidence stays with the maintainer-only harness that owns it and a shipped CLI cannot read an artifact that does not exist in a consumer install; and no changelog entry or version bump, because those belong to the orchestrator at finalization.

## Exploration

REQ-448's eight milestone fields were read first, together with the board's duration model and generator, to establish what already existed: flat optional stamps for planning, dispatch, hand-back, integration, review, remediation and release, plus the claim-to-completion calibration span. This addendum deepens that model only enough to answer how long one request took and why, so nothing in it may displace those fields or that span.

Two findings shaped the design. The Git common directory is the only location that is shared by every linked worktree of one repository and outside the index in all of them, which is what the request's "Git-private state shared by the repository's worktrees" requires — a path under the repository's own `.git` would not be shared by a linked worktree, and anything in the working tree would be stageable.

REQ-539's per-test duration evidence was examined for reuse and deliberately not ingested. Its log is git-ignored and its only writer lives under `_dev/`, which `.gitattributes` marks export-ignore, so a shipped command reading it in a consumer install would read a file that never exists there. The boundary is stated in one sentence of shipped documentation instead.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands_test.go`
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/docs/prescribed-shell-primitives.md`

`work.md` and `work-reference.md` take pointer-level edits only. The request declares no `write_set`, so this boundary is derived from its own Detailed Requirements and its `## Prior Implementation` section.

**Acceptance criteria:** the stream is invisible to git from every worktree and shared by all of them; the summary folds into the request inside the existing finalization commit and the raw stream is then removed; instrumentation creates no agent turn, tracked artifact or commit of its own; REQ-448's fields and calibration span still work; a request with no stream parses and displays unchanged.

## Implementation Summary

**Files changed:**

- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/timing_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**What was done:** A new CLI package records flat per-request lifecycle events into a Git-private JSONL stream and folds them into one compact Timing section in the request file at finalization, then deletes the stream. Three new commands drive it: record-timing-event, run-timed-command and fold-timing-summary.

Each stream line is a complete, independent event: schema version, event id, run id, request id, category, operation, UTC start, UTC end, elapsed seconds, elapsed source, outcome, and optionally revision, responsible agent, exit status and a redacted command identity. Ten categories are defined: recovery-selection, planning, exploration-preflight, builder-work, handback-merge, verification-gate, review, remediation, finalization and cleanup. Commands are redacted to a base name plus an argv token count, so no argument value and no full path ever reaches the stream. An undecodable line is skipped and counted rather than treated as fatal, so a crash mid-append cannot block finalization.

The summary reports total observed elapsed time, per-category totals, the slowest stage event, the slowest command event and unattributed wall time. Nothing else. Folding on an absent stream returns success with the stream marked absent, writes no bytes and leaves the request byte for byte unchanged, so an uninstrumented or historical request finalizes exactly as before. No reader requires the section: it is not in any section-ordering list, the board does not parse it, and the read contract in the reference action now states that tolerance in the same paragraph that already grants it to REQ-448's eight optional stamps.

REQ-448 is untouched. Its eight milestone fields are still parsed, rendered and phase-differenced by the board, and calibration is still claim to completion only; the summary derives its window from events alone and never mentions either field. The one point of contact is recovery, which now strips a stale Timing section for the same reason it already resets the eight stamps: a fresh attempt must not inherit observations from an interrupted run.

The continuation fixed a defect in the section replacement that could silently destroy tracked content. The heading test ran against every line, so a fenced example of the section inside a request body matched, and everything from that match to the next heading was deleted. The heading scan and the section-end scan now share a fence mask, with a fence defined as a run of three or more of one ASCII punctuation mark, excluding the four marks CommonMark defines as thematic breaks and setext underlines. A second defect in the same function dropped the trailing newline when the replaced section was last; the rebuilt document now always ends with exactly one newline. The continuation also confined the fold's request path to the repository root, made the wrapper report shell-faithful exit statuses for signals and missing executables, and removed a three-site contradiction about how the gate is launched.

**Implementation range:** `87226175..4c2cd380`. Builder commits `d3c14c6a` and `eb5a5ed3`.

## Decisions

- **D-01 — No critical path, no spans, no nesting, no exclusive time, no parallelism metrics:** the request's title says critical path and its Detailed Requirements and Constraints both forbid computing one. The requirements were followed. The summary reports total observed elapsed, per-category totals, the slowest stage event, the slowest command event and unattributed wall time, and nothing else.
- **D-02 — Three verbs, not one:** `record-timing-event` closes an interval the caller already lived through, `run-timed-command` runs and measures a command itself, `fold-timing-summary` summarizes and cleans up. The wrapper earns its place because it is the only path where one process observes both ends, so the only path that can honestly claim a monotonic measurement, and the only way to capture an exit status without the caller reporting it. The alternative, making callers pass their own timestamps, is forbidden by the Constraints rule against scattering date calls.
- **D-03 — A missing start chains from the previous event's end:** this makes a serial lifecycle recordable with one flag-light invocation per boundary. A wrapped command is the exception; its start is derived from the measured elapsed, so the wait before the launch stays unattributed instead of being absorbed into the command's number. That is what makes unattributed wall time a real, non-zero quantity rather than an always-zero field.
- **D-04 — The fold runs before the finalization manifest is authored, not inside the finalization transaction:** a hook point inside the transaction exists, but that path binds journal preimages and is high risk with other builders active. Folding into the working request before the manifest is authored puts the section into the preimage the transaction archives and commits, which satisfies the in-the-existing-commit requirement with no coupling to the finalization package.
- **D-05 — `--request-path` is required for the fold:** no repository discovery, no guessing which copy of a request is current. The qualify and scope-drift commands already take the same option, so the orchestrator is already holding the value.
- **D-06 — Commands are redacted to a base name plus an argv token count:** no argument value and no full path ever reaches the stream. A test asserts the stream bytes contain neither a planted secret nor the script's own file name.
- **D-07 — An undecodable stream line is skipped and counted, not fatal:** a crash mid-append must never be the reason a finished request cannot be finalized. The count surfaces as a `TIMING-STREAM-LINE-UNREADABLE` skipped-work entry rather than a new result field.
- **D-08 — Recovery strips a stale Timing section:** REQ-448 already removes its eight stamps on recovery so a fresh attempt cannot inherit observations from an interrupted run. The section is the same class of observation, so it gets the same treatment. That is why the recovery writer and its test are in the manifest; the same file was extended by REQ-448 for the same reason, so it counts as inside the boundary rather than a widening.
- **D-09 — Single-word category values are acceptable:** the naming rule governs identifiers with reach, and these are values in a closed, validated vocabulary, like a status or route value elsewhere in the same schema. They match the pipeline's own step names, which is what makes them findable. Every Go identifier introduced is two words or more.
- **D-10 — No changelog entry, no version bump:** both belong to the orchestrator at finalization.
- **D-11 — The fence rule excludes the four marks `-`, `_`, `*` and `=`:** REQ-544's blanket punctuation class is correct where a region ends at the first block construct and never resumes. Here it would treat the unpaired thematic break above every request's source line as an unclosed fence, hide the section the fold itself appended, and make the next fold stack a duplicate. The four excluded marks are exactly the ones CommonMark defines as thematic breaks and setext underlines, which enclose nothing, so the exclusion is a condition drawn from the spec and not a spelling list. Every other punctuation mark, known to the package or not, still fences.
- **D-12 — An unclosed fence fences everything after it, and ends a replaced span at its opener:** both directions fail toward leaving bytes alone. A heading after an unclosed fence is not matched, so the fold appends rather than deletes; a fence opened inside the section being replaced stops the span at the opener rather than running to end of file.
- **D-13 — Only two of the ten categories are wired:** the builder wait and the repository gate are the two intervals that dominate a serial request, and they were enough to prove the mechanism end to end. Merge, retries, review, remediation, finalization and cleanup currently land in the unattributed number. This is a wiring gap in the action prose, not a gap in the mechanism, and it is recorded as a discovered task rather than fixed here because each remaining call site is a separate prose edit in a file this request was told to touch only at cross-reference level.
- **D-14 — The child's exit status is read through an anonymous interface, not a build-tagged file:** the process state type has no portable signal accessor. A two-file build-tagged split for a status number is more machinery than the problem deserves; the interface assertion matches on Unix, fails harmlessly elsewhere, and the fallback is documented at the function.

## Qualification

Passed the request-bound advance qualify and scope-drift gates for the cumulative range `87226175..4c2cd380`, both satisfied. Twelve files across the request's two commits, matching the hand-back's own manifest; the range stat is larger only because integration was serial and interleaved with other requests.

Independent review reproduced the RED in a scratch copy, with eleven of twelve tests failing on real assertions and eight failure lines matching the hand-back character for character. It verified the mechanism end to end against the built binary rather than the harness.

Its one critical finding — a fold that silently deleted request content when the section heading appeared inside a fenced example — was closed in the continuation and is now covered by regressions in both directions.

Verified after merge: build, vet and gofmt all clean over the CLI module, which mattered because this merged into two files another request had just rewritten.

The P-A-U boxes were reconciled from the builder hand-back, which is where worktree dispatch puts them.
## Testing

Original work, RED first. Stubs with real signatures were written before the tests, so every failure is an assertion and not a compile or import error. `go test -count=1 ./internal/lifecycletiming/` from `skills/do-work/tools/do-work-cli` failed with 11 of 11 tests red:

```
--- FAIL: TestAppendTimingEventWritesFlatGitPrivateStreamAndChainsDefaultStart (0.17s)
    lifecycle_timing_test.go:35: first = resultmodel.LifecycleTimingResult{RequestID:"", RunID:"", StreamPath:"", StreamState:"", ... SectionWritten:false}
--- FAIL: TestAppendTimingEventRedactsCommandArgumentsAndKeepsExitStatus (0.14s)
    lifecycle_timing_test.go:103: recorded = (*resultmodel.TimingEventRecord)(nil)
--- FAIL: TestRunTimedCommandMeasuresInProcessAndReportsChildExitStatus (0.10s)
    lifecycle_timing_test.go:133: exit status = 0, want 3
--- FAIL: TestFoldTimingSummaryReplacesAnExistingTimingSection (0.10s)
    lifecycle_timing_test.go:262: stale summary survived:
--- FAIL: TestAppendTimingEventRejectsUnknownCategoryAndUnsafeIdentity (0.11s)
    lifecycle_timing_test.go:279: expected an unknown category to be rejected
--- FAIL: TestConcurrentRunsAndRequestsKeepSeparateStreams (0.11s)
    lifecycle_timing_test.go:319: streams collided: "" "" ""
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/lifecycletiming	1.282s
```

The other five red tests were `TestFoldTimingSummaryWritesCompactSectionAndRemovesRawStream`, `TestFoldTimingSummaryOnAbsentStreamLeavesTheRequestUnchanged`, `TestTimingCommandsRecordThenFoldThroughArgv`, `TestRunTimedCommandArgvKeepsChildOutputOffTheResultEnvelope` and `TestTimingCommandsRefuseIncompleteArgv`. A twelfth test, `TestRunTimedCommandRealClockCoversTheChildsOwnDuration`, was written after the wrapper existed, because it pins the one property a synthetic clock cannot show.

GREEN for the original work: `go test -count=1 -v ./internal/lifecycletiming/` — PASS, 2.121s, all 12 tests passing, ending `ok github.com/knews2019/skill-do-work/do-work-cli/internal/lifecycletiming 2.121s`. An end-to-end smoke through a real built binary in a throwaway repository recorded one boundary event and one wrapped command, then folded and deleted the stream; the request file gained a `## Timing` section reading `Observed 2026-09-05T00:00:00Z to 2026-09-05T01:33:44Z: 1h 33m 44s total, 1h 33m 44s attributed across 2 events, 0s unattributed.`, and `find .git/do-work-timing -type f` after the fold returned nothing.

Continuation, five defects, RED first. `go test -count=1 -run 'Fenced|Condition|TrailingNewline|OutsideTheRepository|ShellStatuses' ./internal/lifecycletiming/` — five of five failing on assertions:

```
--- FAIL: TestFoldTimingSummaryLeavesAFencedTimingExampleIntact (0.10s)
    lifecycle_timing_test.go:493: the fenced example was damaged:
--- FAIL: TestFoldTimingSummaryReadsFencesAsAConditionNotASpelling (0.09s)
    lifecycle_timing_test.go:523: an unknown fence spelling did not hide its example:
--- FAIL: TestFoldTimingSummaryKeepsATrailingNewlineWhenTheSectionIsLast (0.10s)
    lifecycle_timing_test.go:560: folded file must end with exactly one newline, got ", 20m 00s, exit 0, bash (2 argv tokens)."
--- FAIL: TestFoldTimingSummaryRefusesARequestPathOutsideTheRepository (0.09s)
    lifecycle_timing_test.go:579: fold accepted a request path outside the repository: "../outside-target.md"
--- FAIL: TestRunTimedCommandReportsShellStatusesForSignalsAndMissingExecutables (0.09s)
    lifecycle_timing_test.go:604: signalled status = 128, want 143 (128 + SIGTERM)
```

The first reproduces the reviewer's destruction exactly: the closing fence and the paragraph under it are gone, and the real summary was written inside the example.

GREEN for the continuation: `go test -count=1 -v ./internal/lifecycletiming/` — 19 of 19 pass, 2.843s. The seven new tests all pass, including `TestFoldTimingSummaryLeavesAFencedTimingExampleIntact`, `TestFoldTimingSummaryReadsFencesAsAConditionNotASpelling`, `TestFoldTimingSummaryKeepsATrailingNewlineWhenTheSectionIsLast`, `TestFoldTimingSummaryRefusesARequestPathOutsideTheRepository`, `TestRunTimedCommandReportsShellStatusesForSignalsAndMissingExecutables`, `TestRunTimedCommandArgvCarriesSignalAndLaunchStatusesToTheProcessStatus` and `TestFoldTimingSummaryArgvRefusesAnEscapingRequestPath`. A second end-to-end smoke, on a request body carrying a fenced example of the section, a following paragraph and the unpaired thematic break above a source line, folded twice: the closing fence and the paragraph survived, exactly one unfenced section existed after the re-fold with its content replaced rather than stacked, the last byte was a newline, and an escaping request path was refused with `request path "../escape.md" escapes or names the repository root`.

Module verification, exact commands and durations:

- `gofmt -l .` — printed nothing, both times. `go vet ./...` — exit 0, both times.
- `go test -count=1 ./...` from `skills/do-work/tools/do-work-cli` — 47.7s wall, every package passing except `internal/heavyverification`.
- `go test -count=1 ./internal/lifecycletiming ./internal/requeststate ./internal/resultmodel` — ok in 2.838s, 2.441s and 0.004s.
- `go test -count=1 -run 'Duration|Phase|Milestone|Release' ./...` from `skills/do-work-board/tools/queue-kanban` — ok, 4.001s. This is the REQ-448 guard (the eight optional phase-milestone timestamps and their board phase rows).
- `bash _dev/tests/contracts/core-checks.sh` from the worktree root — `core-checks contract probes passed.`
- `bash _dev/tests/shipped-package-reference-contract.sh` — `shipped package reference contract: PASS`, 7s.
- `bash _dev/tests/contract-regressions.sh` — first run: every content probe passed, one duration breach, `_dev/tests/session-start-hook-behavior.sh took 43s; each test file must finish under 30s` under concurrent-builder load, with that file's own probes reporting `SessionStart hook behavior probes passed.` Continuation run: `Contract regression checks passed.`, exit 0, every duration budget met including that file.
- `GOOS=windows GOARCH=amd64 go build ./...` and `go vet ./internal/lifecycletiming` — both clean. The signal reading goes through an anonymous interface rather than a build-tagged file, so a platform without signals compiles and degrades to the bare 128.

The `internal/heavyverification` failure was diagnosed as a container-environment effect, not a repository defect. It fails at `heavy_runtime_contract_test.go:107` with `default runtime must have a determinable fingerprint` and at `:181` with `shipped runtime probe did not isolate host Git configuration`. Two facts confirmed it. First, it fails identically at the integration tip with the builder's changes stashed, so it is not caused by this request. Second, `env -i PATH=... HOME=... go test -count=1 ./internal/heavyverification/` returns `ok` in 20.563s, while the same command with the session environment fails. The first hand-back filed it as a repository defect; the continuation corrected that entry, so there is nothing to file.

## Discovered Tasks

- Eight of the ten lifecycle categories have no caller in the shipped instructions. `skills/do-work/actions/work.md` records only builder-work at its dispatch paragraph and verification-gate at its testing paragraph, so recovery and selection, exploration and preflight, hand-back and merge, review, remediation, finalization and cleanup all fall into unattributed wall time. Each needs one sentence at its own step. Until then a reader of a Timing section should treat the unattributed number as everything not yet wired, not as lifecycle bookkeeping. → queue as follow-up
- Run identity has no defined source outside worktree dispatch mode, and a mismatch is silent. Recording under one run name and folding under another, or under the standalone default, succeeds with the stream reported as absent, writes no section, and leaves the stream orphaned under the Git common directory. The fix is to name one source for the value, with the run directory name the obvious candidate, and to say so where the fold is invoked. → queue as follow-up
- An interrupted run leaves an orphaned stream under the Git common directory, because only a successful fold deletes it and nothing sweeps that directory. `skills/do-work/actions/cleanup.md` is the natural owner. Not urgent: the residue is Git-private, per-request kilobytes, and a later run uses a new run identity so it never reads stale events. → queue as follow-up
- The write-time chaining defended in decision D-03 has no caller outside the tests, because both wired call sites pass an explicit start. Five command options also have no non-test caller: `--outcome`, `--revision`, `--exit-status`, `--command` and the `--run` default. None is worth deleting while the eight unwired categories remain open, since wiring them is what would use these, but nothing today depends on them. → report only
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_runtime_contract_test.go` lines 107 and 181 fail in this container and pass under a cleaned environment. Environment effect, corrected in the continuation, no repository defect. → report only
- `_dev/tests/session-start-hook-behavior.sh` took 43s against its 30s per-file budget on the first run while several builders were running; its own probes passed and the continuation run met the budget. Worth splitting only if it breaches on an idle machine. → report only

## Review

**Overall: 50% at first review; the critical finding was then closed.**
**Acceptance: Partial at review, remediated after.** The reviewer confirmed the mechanism works end to end against the built binary — the stream is invisible to git from both the main tree and a linked worktree, both share one stream root, redaction holds against a planted token and a password in a URL, and exit statuses 0, 1, 2, 124, 125, 126, 127 and 255 pass through unchanged. It reproduced the RED independently, with eleven of twelve tests failing on real assertions and eight failure lines matching the hand-back character for character.

It settled the title-versus-requirements conflict by reading the original user input at `do-work/user-requests/UR-108/input.md`, where the user's verbatim words are six: "ok, capture it using do-work capture-request". The critical-path promise exists only in the capture agent's restatement and the stale filename, and a dated user decision in the request's own Addendum supersedes it. Following the requirements was correct.

The critical finding: `replaceTimingSection` matched `## Timing` anywhere including inside a fenced code block, then deleted everything to the next heading. The reviewer reproduced a request losing its closing fence and the paragraph under it — silent content destruction in a tracked file. A second defect in the same function dropped the trailing newline when the replaced section was last, which the builder's own re-fold test missed because its fixture had a later section.

Both were closed in the continuation, along with repository-root confinement for `--request-path`, correct signal exit statuses, and a three-site contradiction where `work-reference.md` said to run the gate directly while `work.md` said to launch it through the wrapper.

Eleven further findings are report-only, the largest being that eight of the ten lifecycle categories the request names have no caller in the shipped instructions, so a real run attributes them to unattributed time.

## Lessons Learned

A lifecycle classifier that reads request prose must exclude fenced regions before it matches anything. Any scan that locates a section by heading, and then deletes or rewrites the span to the next heading, will eventually match an example of that section quoted inside a request body and destroy the text around it. The exclusion belongs in the same pass as the heading match, sharing one fence mask, so the heading scan and the section-end scan can never disagree about where prose ends.

State the fence rule as a condition, not as a list of fence spellings. A fence is a run of three or more of one ASCII punctuation mark; testing for backticks alone means a dialect fence the code has never heard of stops hiding what it wraps, and the list goes stale the first time someone writes an example a different way.

The repository-specific trap is that every request file ends with an unpaired thematic break above its source line. A blanket punctuation-run rule reads that as a fence opener that never closes, which hides the section the fold itself appended and makes the next fold stack a duplicate instead of replacing. The four marks CommonMark defines as thematic breaks and setext underlines enclose nothing, so excluding exactly those keeps the rule a condition drawn from the spec rather than a hand-maintained list.

REQ-544 (stating the evidence-position rule as a condition rather than a spelling list) closed the same defect class in the same run on a different surface. Two independent surfaces hit it in one run, which is the signal that the fence condition belongs at subsystem altitude and not in one function's comments.

## Orientation

A finished `do-work run` can now say where one request's wall time went: the shipped CLI records flat, Git-private lifecycle events and folds one compact Timing section into the request inside the existing finalization commit, with no extra agent turn, tracked artifact or instrumentation commit. Two of the ten categories are wired today, so the unattributed number currently means everything not yet instrumented rather than true lifecycle overhead.
