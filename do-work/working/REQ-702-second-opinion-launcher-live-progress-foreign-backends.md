---
id: REQ-702
title: 'Foreign-backend second-opinion launcher shows live progress and streams the backend output while it runs'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: progress feedback for the foreign-backend second-opinion launcher path*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact function (`secondOpinionCapture`), file, and test file, and its Context section already holds the exploration: wire the backend pipes into a tee with a stderr sink and add a heartbeat ticker. The mechanism choice is small and local to one Go file plus two one-sentence doc edits.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
