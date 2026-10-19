---
id: REQ-702
title: 'Foreign-backend second-opinion launcher shows live progress and streams the backend output while it runs'
status: completed
created_at: 2026-10-19T18:20:06Z
user_request: UR-159
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-700]
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
claimed_at: 2026-10-19T18:57:19Z
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-10-19T18:57:48Z
  basis:
    - Route A
    - 4-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
dispatch_at: 2026-10-19T18:59:08Z
builder_handback_at: 2026-10-19T19:05:37Z
integration_at: 2026-10-19T19:05:42Z
review_at: 2026-10-19T19:10:55Z
kb_status: pending
commit: 26a20f1050ee9a9861fd4f0d2cd2382ebc57ca96
heavy_verified_at: 2026-10-19T19:12:53Z
heavy_verified_revision: 26a20f1050ee9a9861fd4f0d2cd2382ebc57ca96
completed_at: 2026-10-19T19:13:16Z
release_at: 2026-10-19T19:13:16Z
---
# Foreign-Backend Second-Opinion Launcher Shows Live Progress And Streams The Backend Output While It Runs
## What
When `ai-report-second-opinion` runs a foreign backend (codex, gemini, agy, or `claude` targeted from another agent's session) as a subprocess, the launcher must show progress while the backend runs: a heartbeat line at least once a minute, and the backend's own stdout and stderr lines mirrored live, both on the launcher's stderr. The typed JSON on stdout, the in-memory capture used for model identity and the reply, the 16 MB capture cap, and the supervisor's process-group guarantees stay as they are.
## Why
Today the launcher buffers everything the backend prints until it exits, up to the 16-minute bound, so the agent's shell pane shows "No output available" for the whole generation and nobody can tell a working backend from a stalled one. REQ-700 (second opinion from the session model runs in-session, released in 0.306.4) fixed this for the session's own model by moving the editorial pass into the session, and named the foreign-backend path as a separate request. This is that request.
## Context
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go`: `secondOpinionCapture` runs the backend under the `/bin/sh` supervisor (`secondOpinionSupervisor`) with `command.Stdout` and `command.Stderr` set to two `secondOpinionBuffer` values that are read only after the process ends; `secondOpinionRun` calls it twice per backend (a short read-only model probe, then the generation) and `secondOpinionMetadata` parses the captured stdout for the reply and the model identity.
- The launcher is invoked with `--format json`, so stdout is the typed result and must stay parseable; progress belongs on stderr. `commandruntime` already prints findings to stderr in text mode, so stderr is the established human channel.
- The session path (`--prepare` / `--finalize`, REQ-700) has no backend subprocess and is out of scope.
## Requirements
- While a backend subprocess runs (probe or generation), the launcher writes a heartbeat line to its stderr at least once per minute: elapsed time, phase (probe, generation, judge), backend name, and the staging page's current size in bytes (or "not written yet"). Lines stop changing when the backend stalls, which is the visible stall signal.
- Every line the backend writes to its stdout or stderr is mirrored to the launcher's stderr as it arrives, prefixed so the reader can tell backend output from launcher heartbeats (the builder picks the prefix). The in-memory capture still receives the complete bytes, so `secondOpinionMetadata`, the reply extraction, the 16 MB cap, and the truncation error behave exactly as today.
- stdout carries only the typed JSON, byte-for-byte the same shape as today; `--format text` keeps its current output plus the same stderr progress.
- The supervisor contract is unchanged: closed stdin, fd 3 completion status, group termination bottom-up, the 16-minute bound. Streaming must not reintroduce the reap-before-cleanup trap recorded under `reaped-by-its-own-parent` in the do-work-cli lessons satellite.
- Update `skills/do-work-toolbox/actions/ai-report.md` § Second-opinion form (one or two sentences: the launcher path shows progress and the backend's output on stderr; read the typed summary from stdout as before) and `skills/do-work-toolbox/docs/ai-report-guide.md` § A Second Opinion (one plain-words sentence). Release per `_dev/primes/prime-releases.md` with a changelog entry.
## Dependencies
None. REQ-700 is archived; this REQ builds on the code it left (`secondOpinionCapture` and `secondOpinionRun` as of 0.306.4) and does not amend it.
## Builder Guidance
Certainty: high on the outcome (heartbeat plus live mirror on stderr, JSON stdout untouched), open on the mechanism. The natural shape is a tee: the backend's pipes feed both the existing buffers and a line-oriented stderr writer, with a ticker goroutine for the heartbeat that stops when the capture returns. Make the stderr sink injectable (for example through the execution context or a package-level writer the test replaces) so a test can assert ordering without reading the process's real stderr. Keep YAGNI: no new flags unless a test needs an interval override (an environment variable or an unexported variable the test sets is enough); no progress file; no change to the prompt template or the session path. Load `_dev/primes/prime-shell-commands.md` before touching `secondOpinionSupervisor`, and prefer not touching it at all.
## Red-Green Proof
**RED prompt/case:** In an agent session, run `.claude/skills/do-work/tools/do-work-cli.sh --repo-root . --format json ai-report-second-opinion --backend codex ai-reports/<bundle>`; for the whole generation the shell pane shows "No output available". A Go test in `report_second_opinion_test.go` with a fake backend that prints a line, sleeps past the (test-shortened) heartbeat interval, then writes the page and prints its envelope asserts that the launcher's stderr sink received the backend's line and at least one heartbeat before the backend exited.
**Why RED now:** `secondOpinionCapture` connects the backend to two in-memory buffers and nothing else; no heartbeat exists; stderr stays empty until the typed result is printed.
**GREEN when:** The same invocation shows, while the backend runs, heartbeat lines at least once a minute (elapsed, phase, backend, staging page size) and the backend's own output lines as they arrive, all on stderr; stdout still holds only the typed JSON and parses as before; the fake-backend test sees the backend line and a heartbeat before exit, and the existing second-opinion tests (model identity, reply extraction, truncation, supervisor safety) stay green.
**Validation:** User adjusted (the user added live streaming of the backend's own output to the heartbeat proposal).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22518 tokens, over the 2000 budget; `slugged: partial` so no targeted entry is legal. Matched because the REQ changes the backend subprocess capture in do-work-cli internals (families reaped-by-its-own-parent, background-worker-self-bound, interruptible-blocking-io, native-backend-identity).
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, over budget; `slugged: partial`. Matched because the supervisor is a prescribed shell block shipped inside the Go source.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Tee inside the existing `secondOpinionBuffer.Write`: append to the 16 MB capture as before and split complete lines to a shared, mutex-guarded sink with a `[<cli> stdout] ` / `[<cli> stderr] ` prefix. A ticker goroutine per capture writes `[second-opinion] <cli> <phase> running <elapsed>, page <N bytes|not written yet>` and is joined before the capture returns. Supervisor untouched; writers stay non-file values so `Wait` still joins exec's copy goroutines before the flush. Seams: package-level sink and interval variables, following `secondOpinionJudge`.
- [x] **[APPLY]:** As planned, plus one change the RED-to-GREEN step forced: the first GREEN attempt produced heartbeats but zero mirrored lines, because `secondOpinionBuffer` embedded `bytes.Buffer`, whose promoted `ReadFrom` `io.Copy` prefers, so `Write` never ran. Making the buffer a named field fixed the mirror and revived the 16 MB cap (dead since it was written; see D-03). Judge heartbeat added (D-01).
- [x] **[UNIFY]:** `git diff --stat` (commit aa61c81a): ai-report.md 2 +-, ai-report-guide.md 2 +-, report_second_opinion.go 108 +++--, report_second_opinion_test.go 75 +++-; 4 files, 172 insertions, 15 deletions. `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands`: clean. `gofmt -l` on the package: empty. `go test -C skills/do-work/tools/do-work-cli -count=1 -run SecondOpinion ./internal/toolboxcommands`: ok 13.6s. The whole package (`./internal/toolboxcommands`, all tests): ok 15.2s. I reviewed the diff for debug artifacts: none. stdout path untouched (`commandruntime` and the result model are unchanged, and the sink never touches stdout).
*Source: progress feedback for the foreign-backend second-opinion launcher path*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact function (`secondOpinionCapture`), file, and test file, and its Context section already holds the exploration: wire the backend pipes into a tee with a stderr sink and add a heartbeat ticker. The mechanism choice is small and local to one Go file plus two one-sentence doc edits.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

(from the builder hand-back)

- D-01 (DECIDE & STATE) Judge heartbeat: included. `secondOpinionRender` runs a headless-browser subprocess through `judgeRender`: browser startup plus 2 viewports × 2 color schemes, each wait bounded by `judgeProtocolDeadline` = 30s. So the worst case is over a minute of silence. Wrapping the `secondOpinionJudge` call costs 2 lines and reuses the same heartbeat. Its page size is the published page. Side effect: `secondOpinionPublication` is shared with the session `--finalize` path, so a finalize whose judge runs longer than 30s also prints heartbeats on stderr. No change to prepare/finalize logic or output. The browser's own output is not mirrored because it lives in `report_judge.go`, outside the write set.
- D-02 (DECIDE & STATE) Interval 30s, so "at least once per minute" holds with margin. Prefixes: `[second-opinion] ` for heartbeats, `[<cli> stdout] ` / `[<cli> stderr] ` for backend lines. The page shows as "not written yet" while the staging file is absent or still empty (CreateTemp leaves it at 0 bytes).
- D-03 (DECIDE & STATE, flagged for the orchestrator) The 16 MB cap now actually runs. The brief said the cap must "behave exactly as today", but today it never fires: `io.Copy` used the promoted `bytes.Buffer.ReadFrom`, so a runaway backend was buffered without limit and with no error (shown above: 16777300 bytes captured, err nil). The mirror needs `Write` to run, so the dead cap could not stay dead without a separate bypass. The result now matches the documented REQ-698 intent. Value: memory bound and the typed "backend output exceeded capture limit" error are real. Risk: a backend that legitimately emits more than 16 MB per run (no evidence of one; agy stream-json is the largest producer) now fails instead of succeeding. To undo, raise the constant in `Write`.
- D-04 (DECIDE & STATE) The mirror keeps printing past the 16 MB cap. Memory stays bounded because the only retained mirror state is one unfinished line, cut at 64 KiB.
- D-05 (DECIDE & STATE) `secondOpinionCapture` takes the progress descriptor as a new parameter rather than reading it from context or a global. The function is unexported and has three callers. The three direct test calls pass `secondOpinionProgress{}`, which produces prefixes like `[ stdout] ` that only tests see.

## Discovered Tasks

(from the builder hand-back)

- The `ownedprocess`/report-image code may have the same embedded-`bytes.Buffer` writer shape. I did not check it, because it is out of scope. Worth one grep for structs that embed `bytes.Buffer` and override `Write` to cap, across do-work-cli: impact-noncritical → report only.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)

**What was done:** `secondOpinionCapture` now takes a progress descriptor (phase, backend, page path), starts a 30-second heartbeat ticker after `Start` that is joined before the capture returns, and its two stream buffers mirror each backend line to a mutex-guarded package-level stderr sink with a `[<cli> stdout] ` / `[<cli> stderr] ` prefix while still capturing the full bytes; the render judge call is wrapped in the same heartbeat. The buffer holds its `bytes.Buffer` in a named field, which makes the 16 MB cap run for the first time (D-03). Two new tests pin the live mirror plus heartbeat before backend exit and the cap; both docs gained one sentence each.

## Qualification

**Gate record:** `advance REQ-702 --diff-range e28c3945..26a20f10` → qualify `satisfied`, no findings (no debug artifacts, P-A-U boxes ticked from the hand-back). Route A, so no scope comparison.

**Diff:** 4 files, +172 −15, exactly the REQ's `write_set` (`git diff --stat e28c3945..26a20f10`).

**Requirement trace (read against the merged code):**
- Heartbeat at least once a minute with elapsed, phase, backend, page size or "not written yet": `secondOpinionHeartbeat` ticks every 30 s; phases `probe` (agy models + read-only probe), `generation` (staging page path), `judge` (D-01). The ticker goroutine is joined by `stop`, which `secondOpinionCapture` defers right after `Start`.
- Backend lines mirrored live with a prefix, full bytes still captured: `secondOpinionBuffer.Write` splits complete lines to the mutex-guarded sink (`[<cli> stdout] ` / `[<cli> stderr] `), keeps at most one 64 KiB unfinished line, flushes a newline-less tail after `Wait`, and still appends to `captured`. The 16 MB cap now runs (D-03): before this REQ, `io.Copy` used the embedded buffer's promoted `ReadFrom` and bypassed `Write`. This differs from "exactly as today", but "today" was a defect: the cap and its truncation error had never fired. The new test `TestSecondOpinionCaptureCapsBackendOutput` pins the cap.
- stdout only the typed JSON: the sink defaults to `os.Stderr`; `commandruntime` and the result model are unchanged.
- Supervisor contract: `secondOpinionSupervisor` is byte-identical; `TerminateGroup` still runs before `Wait`; the writers stay non-file values, so exec's copiers are joined by `Wait` before `flushMirror`.
- Docs: one sentence in `ai-report.md` § Second-opinion form, one in `ai-report-guide.md` § A Second Opinion.

**Judgment notes:** the mirror writes from exec's copy goroutines. If the caller's stderr were a full, unread pipe, those copies would block and the backend would stall on its own stdout. Agent shells drain stderr, so this stays report-only (see Review).

## Testing

**Tests run:** `bash do-work/runs/work-2026-10-19-185719/REQ-702-probe.sh` (runs `go test -count=1 -v -run SecondOpinion ./internal/toolboxcommands` in the do-work-cli module and requires `--- PASS: TestSecondOpinionCaptureStreamsProgressWhileBackendRuns`), 12.7 s; repository gate `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge 26a20f10, load 1.8 at start, wall 122 s.
**Result:** ✓ All passing. The probe passed, the gate exited 0 on its first run ("Maintainer verification passed."), and `advance` recorded `run-blocked-check` and `green-gate` as satisfied.

**Red-green validation:**
- `TestSecondOpinionCaptureStreamsProgressWhileBackendRuns` (`report_second_opinion_test.go`): ✗ before implementation (`backend exited before its line and a heartbeat reached the sink: ""`) → ✓ after. This traces the `## Red-Green Proof` RED/GREEN pair: a fake Codex backend prints a line, waits for a marker that the recording sink writes only after it has seen the prefixed backend line, a `codex generation` heartbeat and a page size, then exits.
- `TestSecondOpinionCaptureCapsBackendOutput`: ✗ against the old embedded-`bytes.Buffer` shape (`captured 16777300 bytes, err=<nil>`) → ✓ after (D-03).

**New tests added:**
- `TestSecondOpinionCaptureStreamsProgressWhileBackendRuns`
- `TestSecondOpinionCaptureCapsBackendOutput`

**Existing tests updated (cross-REQ impact):**
- `report_second_opinion_test.go` (from REQ-698/700/701): three direct `secondOpinionCapture` calls pass `secondOpinionProgress{}` for the new parameter. No behavior assertion changed.

**Heavy verification plan:**
- Range: e28c3945b14e574740b245a47c9fef5bb2defe69..26a20f1050ee9a9861fd4f0d2cd2382ebc57ca96
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`: report_second_opinion.go and its test matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`: all four changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`: the two Go files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`: the two Go files matched subtree skills/do-work/tools/do-work-cli

*Verified by work action*

## Review

**Overall: 93%** | 2026-10-19T19:10:55Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 the stderr mirror writes synchronously on exec's copy goroutines under the global mutex (`report_second_opinion.go:55-58`), so a closed stderr reader (`2>&1 | head`) kills the launcher by SIGPIPE mid-run with no JSON and a leftover staging page, and a stalled stderr reader blocks the copier so `Wait` outlives the 16-minute bound; no in-repo caller does either; fix by ignoring SIGPIPE and dropping lines when the sink is slow — impact-user-visible → report only. Nits: F2 the cap test checks `err != nil`, not the `backend output exceeded capture limit` text — impact-negligible → report only; F3 the buffer comment's "up to 64 KiB" pending bound is really about 96 KiB — impact-negligible → report only
**Acceptance:** Pass — implementation + integration: vet clean; `-race -run SecondOpinionCapture` ok (streaming before exit, cap, cancellation, drain, status); live acceptance unassessed
**Restatement sweep:** redefined second-opinion launcher stderr output (heartbeat + prefixed backend lines) and the 16 MB capture cap (now enforced); swept skills/ and _dev/primes, no stale restatements
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: e28c3945b14e574740b245a47c9fef5bb2defe69
- Target revision: 26a20f1050ee9a9861fd4f0d2cd2382ebc57ca96
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`: report_second_opinion.go and its test matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`: all four changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`: the two Go files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`: the two Go files matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target revision: 26a20f1050ee9a9861fd4f0d2cd2382ebc57ca96
- Execution revision: 26a20f1050ee9a9861fd4f0d2cd2382ebc57ca96 (detached checkout `.git/work-run-work-2026-10-19-185719/drain-head`)
- do-work-cli-integrations: exit 0, executed, 62 s
- staged-skills: exit 0, executed, 36 s
- updater: exit 0, executed, 65 s
- installer: exit 0, executed, 27 s

## Lessons Learned

**What worked:** Writing the RED test first, against the real handler with a fake backend and a marker-file handshake, showed at once that heartbeats arrived but mirrored lines never did. That led straight to the dead `Write`. A package-level sink plus an interval variable, following the existing `secondOpinionJudge` seam, kept the test deterministic without new flags.
**What didn't:** The first GREEN attempt put the tee in `secondOpinionBuffer.Write` while the type still embedded `bytes.Buffer`. `io.Copy` (and `exec.Cmd`'s pipe copier) prefers the promoted `ReadFrom`, so `Write`, and with it the 16 MB cap, had never run since REQ-698.
**Worth knowing:** The mirror writes synchronously to the process's stderr. A caller that closes stderr early (`| head`) or stops reading it can kill the run by SIGPIPE or stall it past the 16-minute bound (review F1, report only). If a caller like that appears, drop mirror lines when the sink is slow instead of blocking the copier. agy `stream-json` mirrors one line per NDJSON event, which can be loud in an agent's tool output.

## Orientation

Now a second opinion from codex, gemini, agy, or a cross-agent `claude` shows a 30-second heartbeat and the backend's own output on stderr while it runs, and its 16 MB capture cap is enforced. This lives in the do-work-cli second-opinion launcher (`toolboxcommands/report_second_opinion.go`), described in the toolbox `ai-report` second-opinion form. Prime spot-check: `prime-shell-commands.md` and `prime-releases.md` reference no path this change touched; neither is stale.

## Timing

Observed 2026-10-19T18:59:08Z to 2026-10-19T19:12:54Z: 13m 46s total, 12m 11s attributed across 4 events, 1m 35s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 6m 29s | 1 |
| verification-gate | 5m 36s | 2 |
| handback-merge | 6s | 1 |

Slowest stage: builder-work / builder worktree build, 6m 29s, outcome success.
