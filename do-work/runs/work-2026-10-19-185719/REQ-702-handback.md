# REQ-702 hand-back (builder)

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-702-second-opinion-launcher-live-progress-foreign-backends`, branch `worktree-agent-REQ-702-second-opinion-launcher-live-progress-foreign-backends`.

## Files

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` — modified. Adds the progress sink (`secondOpinionProgressSink`, default `os.Stderr`), `secondOpinionHeartbeatInterval` (30s), a sink mutex, `secondOpinionProgress {phase, backend, pagePath}`, `secondOpinionHeartbeat` (ticker goroutine, stop joins it). `secondOpinionBuffer` now holds a named `captured bytes.Buffer` (no embedding), mirrors each line with a prefix, bounds the unfinished-line buffer at 64 KiB, and `flushMirror` prints a newline-less last line after `Wait`. `secondOpinionCapture` takes a `secondOpinionProgress` and runs the heartbeat from `Start` until return. `secondOpinionRun` passes phase `probe` (agy `models` and the read-only probe) and `generation` (with the staging page path). `secondOpinionPublication` wraps `secondOpinionJudge` in a `judge` heartbeat.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` — modified. New `TestSecondOpinionCaptureStreamsProgressWhileBackendRuns` (the RED/GREEN test) with a `progressRecorder` sink. New `TestSecondOpinionCaptureCapsBackendOutput`. Three existing direct `secondOpinionCapture` calls get `secondOpinionProgress{}` for the new parameter.
- `skills/do-work-toolbox/actions/ai-report.md` — modified, § Second-opinion form: one sentence (heartbeat and prefixed backend lines on stderr; read the typed summary from stdout as before).
- `skills/do-work-toolbox/docs/ai-report-guide.md` — modified, § A Second Opinion: one plain sentence.

No file outside the `write_set` was touched.

## P-A-U

[PLAN] Tee inside the existing `secondOpinionBuffer.Write`: append to the 16 MB capture as before and split complete lines to a shared, mutex-guarded sink with a `[<cli> stdout] ` / `[<cli> stderr] ` prefix. A ticker goroutine per capture writes `[second-opinion] <cli> <phase> running <elapsed>, page <N bytes|not written yet>` and is joined before the capture returns. Supervisor untouched; writers stay non-file values so `Wait` still joins exec's copy goroutines before the flush. Seams: package-level sink and interval variables, following `secondOpinionJudge`.

[APPLY] As planned, plus one change the RED-to-GREEN step forced: the first GREEN attempt produced heartbeats but zero mirrored lines, because `secondOpinionBuffer` embedded `bytes.Buffer`, whose promoted `ReadFrom` `io.Copy` prefers, so `Write` never ran. Making the buffer a named field fixed the mirror and revived the 16 MB cap (dead since it was written; see D-03). Judge heartbeat added (D-01).

[UNIFY] `git diff --stat` (commit aa61c81a): ai-report.md 2 +-, ai-report-guide.md 2 +-, report_second_opinion.go 108 +++--, report_second_opinion_test.go 75 +++-; 4 files, 172 insertions, 15 deletions. `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands`: clean. `gofmt -l` on the package: empty. `go test -C skills/do-work/tools/do-work-cli -count=1 -run SecondOpinion ./internal/toolboxcommands`: ok 13.6s. The whole package (`./internal/toolboxcommands`, all tests): ok 15.2s. I reviewed the diff for debug artifacts: none. stdout path untouched (`commandruntime` and the result model are unchanged, and the sink never touches stdout).

## Red-Green

Test: `TestSecondOpinionCaptureStreamsProgressWhileBackendRuns`. It goes through the real handler with the fake Codex fixture. After the page is written, the fake prints `drafting` to stderr and waits (bounded at 4s) for a marker file. The recorder sink writes that marker only when it has seen `[codex stderr] drafting\n`, a `codex generation` heartbeat and ` bytes` (page size). So the marker existing after the run proves that all three reached the sink before the backend exited. After the run the test also asserts `[codex stdout] Clarity changes\n` (stdout mirror) and `[codex stderr] tail without newline\n` (partial-line flush).

RED (only the sink and interval variables existed):
```
--- FAIL: TestSecondOpinionCaptureStreamsProgressWhileBackendRuns (6.08s)
    report_second_opinion_test.go:791: backend exited before its line and a heartbeat reached the sink: ""
FAIL
```
GREEN:
```
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	0.822s   (-run TestSecondOpinionCaptureStreamsProgressWhileBackendRuns)
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	13.643s  (-run SecondOpinion)
```
Both new tests against a temporary re-embedded `bytes.Buffer` variant (restored afterwards) fail, which proves they pin the dead-Write failure:
```
report_second_opinion_test.go:792: backend exited before its line and a heartbeat reached the sink: "[second-opinion] codex probe running 0s, ..."   (heartbeats only, no mirrored lines)
--- FAIL: TestSecondOpinionCaptureCapsBackendOutput (0.08s)
    report_second_opinion_test.go:810: captured 16777300 bytes, err=<nil>
```

## Decisions

- D-01 (DECIDE & STATE) Judge heartbeat: included. `secondOpinionRender` runs a headless-browser subprocess through `judgeRender`: browser startup plus 2 viewports × 2 color schemes, each wait bounded by `judgeProtocolDeadline` = 30s. So the worst case is over a minute of silence. Wrapping the `secondOpinionJudge` call costs 2 lines and reuses the same heartbeat. Its page size is the published page. Side effect: `secondOpinionPublication` is shared with the session `--finalize` path, so a finalize whose judge runs longer than 30s also prints heartbeats on stderr. No change to prepare/finalize logic or output. The browser's own output is not mirrored because it lives in `report_judge.go`, outside the write set.
- D-02 (DECIDE & STATE) Interval 30s, so "at least once per minute" holds with margin. Prefixes: `[second-opinion] ` for heartbeats, `[<cli> stdout] ` / `[<cli> stderr] ` for backend lines. The page shows as "not written yet" while the staging file is absent or still empty (CreateTemp leaves it at 0 bytes).
- D-03 (DECIDE & STATE, flagged for the orchestrator) The 16 MB cap now actually runs. The brief said the cap must "behave exactly as today", but today it never fires: `io.Copy` used the promoted `bytes.Buffer.ReadFrom`, so a runaway backend was buffered without limit and with no error (shown above: 16777300 bytes captured, err nil). The mirror needs `Write` to run, so the dead cap could not stay dead without a separate bypass. The result now matches the documented REQ-698 intent. Value: memory bound and the typed "backend output exceeded capture limit" error are real. Risk: a backend that legitimately emits more than 16 MB per run (no evidence of one; agy stream-json is the largest producer) now fails instead of succeeding. To undo, raise the constant in `Write`.
- D-04 (DECIDE & STATE) The mirror keeps printing past the 16 MB cap. Memory stays bounded because the only retained mirror state is one unfinished line, cut at 64 KiB.
- D-05 (DECIDE & STATE) `secondOpinionCapture` takes the progress descriptor as a new parameter rather than reading it from context or a global. The function is unexported and has three callers. The three direct test calls pass `secondOpinionProgress{}`, which produces prefixes like `[ stdout] ` that only tests see.

## Discovered Tasks

- The `ownedprocess`/report-image code may have the same embedded-`bytes.Buffer` writer shape. I did not check it, because it is out of scope. Worth one grep for structs that embed `bytes.Buffer` and override `Write` to cap, across do-work-cli: impact-noncritical → report only.

## Lessons Read

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: every `[family: native-backend-identity]` (0.306.0, 0.306.3), `[family: session-is-the-backend]` (0.306.4), `[family: reaped-by-its-own-parent]` (0.306.0, REQ-543), `[family: background-worker-self-bound]` (2026-09-05), `[family: interruptible-blocking-io]` (REQ-451), `[family: fixture-cost-is-subprocess-spawning]` (REQ-574).
- `_dev/primes/lessons-shell-commands.md`: both `[family: assertion-passes-when-feature-is-dead]` bullets (REQ-591, REQ-592). This family applied directly: the cap's `Write` was dead and no test read it.
- Crew: general, coding-guardrails, shared-principles, testing, communication-style. Primes: prime-shell-commands, prime-releases (read only; no release made).
- The fake backend's wait loop carries its own 4s bound (background-worker-self-bound). The supervisor is untouched (reaped-by-its-own-parent).

## Proposed Lesson

- [family: assertion-passes-when-feature-is-dead] 0.306.x (REQ-702, 2026-10-19): **a type that embeds `bytes.Buffer` and overrides `Write` never sees `io.Copy`'s bytes.** The embedded buffer promotes `ReadFrom`, and `io.Copy` (including `exec.Cmd`'s pipe copier) prefers it, so the second-opinion 16 MB cap and its truncation error never ran from REQ-698 onward. The live mirror added on top of `Write` showed the defect at once: heartbeats arrived, mirrored lines never did. Hold the buffer in a named field so `Write` is the only path, and pin any cap with a test that pushes past it.

## Commits

- `aa61c81a` [REQ-702] stream second-opinion backend output and heartbeats to stderr
