## Review: REQ-702

**Approve** — the launcher now shows a heartbeat and the backend's own lines on stderr while a foreign backend runs, stdout stays the typed JSON, and the 16 MB cap really works for the first time.
Route A | merge range e28c3945..26a20f10 (builder commit aa61c81a, merge 26a20f10)

### What's built
- `secondOpinionCapture` mirrors every backend stdout/stderr line to stderr with a `[<cli> stdout] ` / `[<cli> stderr] ` prefix and prints a `[second-opinion]` heartbeat every 30 s (backend, phase, elapsed, page size). The render judge gets the same heartbeat.
- The stream buffer no longer embeds `bytes.Buffer`, so `io.Copy` goes through `Write`. That makes the mirror work and makes the 16 MB cap and its truncation error real (D-03).
- Still missing (not builder-owned): the release and changelog entry the REQ asks for. Finalization owns it.

### Decisions / risks for you
- D-03 (the 16 MB cap now fires). Value: memory is bounded and a runaway backend gets a typed error. Risk: a backend that legitimately prints more than 16 MB on one stream now fails instead of succeeding. No such backend is known. The changelog entry should say this in one sentence, because it is a behavior change, not only a progress feature.

### Findings

**Important:**
- None.

**Minor:**
- F1. The mirror writes to the real stderr synchronously, on exec's copy goroutines, under the global `secondOpinionProgressMutex` (`report_second_opinion.go:55-58`, called from `secondOpinionBuffer.Write`). Two new failure modes follow. (a) Closed reader: with `... 2>&1 | head -n 20`, Go exits on SIGPIPE at the next mirrored line. The run dies mid-generation with no typed JSON, a staging page left in the bundle, and no launcher teardown of the backend group. A scratch Go program that does the same stderr writes exits 141 under `| head -n 2` and never prints its stdout line. Before this REQ the launcher wrote nothing until the end, so `| head` still got the JSON. (b) Stalled reader: if a caller pipes stderr and reads stdout first (for example Popen with `stderr=PIPE`), the copier blocks in `Write` while holding the mutex. The backend then blocks on its own stdout. After the 16-minute bound, `command.Wait()` still waits for the blocked copier, and the heartbeat's `stop` waits for the mutex, so the 16-minute bound no longer bounds the run. No caller in this repository does either; agent shells drain stderr. Smallest fixes: ignore SIGPIPE (`signal.Ignore(syscall.SIGPIPE)`, write errors are already discarded), and drop mirror lines when the sink is slow instead of blocking the copier. — impact-user-visible → report only

**Nit:**
- F2. `TestSecondOpinionCaptureCapsBackendOutput` checks `captured.err != nil` but not the text `backend output exceeded capture limit` (`report_second_opinion_test.go`, the cap test). Another capture error with a 16 MiB stdout would also pass. The red result (16777300 bytes, err nil) still proves it pins the dead-cap defect. — impact-negligible → report only
- F3. The `secondOpinionBuffer` comment says `pending` holds "at most one unfinished line of up to 64 KiB". The cut runs only after a whole `Write` is appended, so `pending` can reach about 96 KiB (64 KiB − 1 plus one 32 KiB io.Copy chunk), and a cut line can be that long. Memory is still bounded; only the number in the comment is off. — impact-negligible → report only

### Concurrency check (asked by the orchestrator)
- Data races: none found. Each buffer has one writer (exec's copier). `flushMirror`, `truncated` and `captured` are read only after `Wait` joins the copiers. The heartbeat goroutine reads a by-value `progress` copy and the interval global. The sink global is read under the mutex. Tests swap the globals before the invoke and restore them in `t.Cleanup`, after every heartbeat is joined. The file has no `t.Parallel`. `go test -race -run SecondOpinionCapture` is clean.
- Goroutine leaks: none. The heartbeat starts only after a successful `Start`, `stop` is deferred right after, and `stop` closes the channel once and waits for the ticker goroutine to exit. The judge heartbeat is stopped on the next line (it leaks only if the judge panics, and a panic ends the process anyway).
- Deadlock if the sink blocks: yes, see F1(b).
- Line splitting at 64 KiB: correct. Lines that have a newline are emitted whole at any length. A run with no newline is emitted once `pending` is ≥ 64 KiB, with a newline added. The newline-less tail is flushed after `Wait`. A cut may split a multi-byte UTF-8 character. That is cosmetic, and it is on stderr only.
- 16 MB cap (D-03): correct. With `bytes.Buffer` embedded, `io.Copy` uses the promoted `ReadFrom`, so `Write` never ran. A named field removes `ReadFrom` from the method set. The new cap test pins this: the builder's red run against the re-embedded type captured 16777300 bytes with err nil. The mirror continues past the cap (D-04) without growing memory.
- Supervisor contract: `secondOpinionSupervisor` is byte-identical. `TerminateGroup` still runs before `Wait`, so the reap-before-cleanup trap is not reintroduced.

### Anti-bloat count
- Unrequested flags, options, env vars or files: 0.
- New helpers: `secondOpinionProgress` (3 fields), `secondOpinionProgressLine` (2 callers), `secondOpinionHeartbeat`, `mirrorLine` (2 callers), `flushMirror`. Each one is used by the required behavior.
- The judge heartbeat (D-01) is in scope, because the REQ lists `judge` among the phases.
- Tests: 2 new tests plus a `progressRecorder` sink. The recorder's marker-file handshake is what makes "before the backend exited" deterministic. The cap test pins the D-03 behavior change. Neither test is decorative.
- Verdict: no bloat found.

### Builder's discovered task (checked)
- The builder asked whether other do-work-cli structs embed `bytes.Buffer` the same way. A grep for embedded `bytes.Buffer` fields across `skills/do-work/tools/do-work-cli` finds none. Nothing to capture.

### Requirements Checklist

- [x] Heartbeat at least once a minute with elapsed time, phase (probe, generation, judge), backend and staging page size or "not written yet" — delivered (30 s ticker)
- [x] Every backend stdout/stderr line mirrored live with a prefix; in-memory capture still gets the full bytes — delivered
- [x] Metadata, reply extraction, cap and truncation error behave as before — delivered with one documented deviation: the cap was dead and now runs (D-03)
- [x] stdout carries only the typed JSON; `--format text` keeps its output plus the stderr progress — delivered (the sink defaults to `os.Stderr`; `commandruntime` is unchanged)
- [x] Supervisor contract unchanged, no reap-before-cleanup — delivered
- [x] One or two sentences in `ai-report.md` § Second-opinion form and one in `ai-report-guide.md` § A Second Opinion — delivered
- [ ] Release with changelog entry — not in the merge range; owned by finalization, not scored here

### Acceptance Testing

**Result: Pass** (implementation and integration stages: targeted tests through the real handler with a fake codex backend, under the race detector)
- `go vet ./internal/toolboxcommands`: clean.
- `go test -count=1 -race -run 'SecondOpinionCapture' ./internal/toolboxcommands`: ok, 5.1 s. This covers streaming before backend exit, the cap, the cancellation and drain cases, and status/channel preservation.
- The builder's red/green evidence shows both new tests fail against the old embedded-buffer type and pass after the change.
- The full package and the repository gate were not run, by instruction (the orchestrator runs the gate).

### Suggested Additional Testing

- Live acceptance: unassessed. Run `ai-report-second-opinion --backend codex` (or agy) on a real bundle from an agent shell. Confirm that heartbeats and prefixed backend lines appear during generation, and that stdout still parses as JSON.
- Edge case: run a real agy generation with `--output-format stream-json` and check how much mirrored stderr it produces. Every NDJSON event becomes one line, which may flood an agent's tool output.

### Restatement Sweep

The diff redefines two things: the launcher's stderr output during a foreign-backend run (silent before, heartbeat plus mirrored lines now), and the 16 MB capture cap (dead before, enforced now). I grepped `skills/` and `_dev/primes` for `16 MB`/`16 MiB`, `capture limit`/`capture cap`, `No output available`, and launcher stdout/stderr/progress wording. I also read every `.md` file that mentions second opinion: `ai-report.md`, `ai-report-guide.md`, `ai-report-reference.md`, `SKILL.md`, the `lessons-do-work-cli.md` second-opinion bullets, and `prime-agy-usage.md`. The two updated sentences are the only text that describes launcher output. `prime-agy-usage.md` covers direct AGY runs, not the launcher. No restatement describes the cap. `_dev/tests/do-work-cli-launcher-behavior.sh` makes no stderr assertion for this command. Stale restatements: none.

### Scores (on the record — not the headline)

**Overall: 93%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All builder-owned requirements delivered; release left to finalization |
| Code Quality | 88% | Clean and well commented; sink is coupled to the copiers (F1) |
| Test Adequacy | 85% | Red-green on both new tests, race-clean; cap test does not check the error text (F2) |
| Scope | 100% | Exactly the 4-file write_set |
| Risk | Low | F1 only bites callers that close or stop reading stderr |
| Acceptance | Pass | Implementation + integration via targeted race-enabled tests |

### Follow-ups created
None (3 findings report only)

---

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
