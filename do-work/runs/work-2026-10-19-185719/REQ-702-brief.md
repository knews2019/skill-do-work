# Builder brief — REQ-702 (run work-2026-10-19-185719)

You are the **builder** for REQ-702: the foreign-backend second-opinion launcher shows live progress (heartbeats) and mirrors the backend's own output on its stderr while the backend subprocess runs. Route A, `tdd: true`.

## Where you work
- Worktree (your only write location): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-702-second-opinion-launcher-live-progress-foreign-backends`, branch `worktree-agent-REQ-702-second-opinion-launcher-live-progress-foreign-backends`. Commit there. Never write the main checkout `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` (call it ROOT). Never push. Never touch `do-work/`.
- Hand-back file: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/49f49f76-f600-48fb-908a-3a36236378f7/scratchpad/REQ-702/handback.md`. Use only that `REQ-702/` subfolder of the scratchpad and delete nothing outside it.

## Read first
1. The REQ: `ROOT/do-work/working/REQ-702-second-opinion-launcher-live-progress-foreign-backends.md` (read-only for you). Its `## Requirements`, `## Builder Guidance`, and `## Red-Green Proof` are the contract.
2. Crew rules: `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `testing.md`, `communication-style.md` (paths relative to the worktree).
3. Primes: `_dev/primes/prime-shell-commands.md` (before going near `secondOpinionSupervisor`; prefer not touching it) and `_dev/primes/prime-releases.md` (you do NOT release; the orchestrator does — read only so you know a changelog entry follows).
4. Lessons (both satellites are over the budget, so read these families only): in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` every bullet tagged `[family: reaped-by-its-own-parent]`, `[family: background-worker-self-bound]`, `[family: interruptible-blocking-io]`, `[family: native-backend-identity]`, `[family: session-is-the-backend]`, `[family: fixture-cost-is-subprocess-spawning]`; in `_dev/primes/lessons-shell-commands.md` the `[family: assertion-passes-when-feature-is-dead]` bullets. List what you read in the hand-back.

## The code
`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go`:
- `secondOpinionCapture` (about line 820) runs the backend under the `/bin/sh` supervisor with `command.Stdout`/`command.Stderr` set to two `secondOpinionBuffer` values (16 MB cap, `truncated` flag) read only after exit.
- `secondOpinionRun` (about line 513) calls it for the agy `models` probe, the read-only model probe, and the generation. The staging page name (`stageName`) exists only from the generation onward; before that the heartbeat says "not written yet".
- The render judge (`secondOpinionJudge` = `secondOpinionRender`) runs in-process, not as a backend subprocess. The REQ lists "judge" as a heartbeat phase: decide whether any judge step actually runs a subprocess long enough to need one; if not, record a D-XX decision (DECIDE & STATE) and skip it.
- Existing package pattern for test seams: `var secondOpinionJudge = secondOpinionRender`. A package-level `io.Writer` progress sink (default `os.Stderr`) plus an unexported heartbeat interval variable the test shortens fits the REQ's YAGNI guidance. `commandruntime.ExecutionContext` has no stderr field; do not add one.

## What to build (outcome fixed, mechanism yours)
- Heartbeat line on the sink at least once per minute while a backend subprocess runs: elapsed, phase, backend name, staging page size in bytes or "not written yet". It stops when the capture returns (no goroutine outlives the capture).
- Every backend stdout/stderr line mirrored to the sink as it arrives, with a prefix that tells it apart from heartbeats. The two copy goroutines and the ticker share the sink: serialize writes so lines never interleave mid-line. A partial last line without a newline is still flushed at the end. The in-memory buffers still receive the complete bytes, so the 16 MB cap and "backend output exceeded capture limit" behave exactly as today (decide whether mirroring continues past the cap; a mirror that keeps printing a runaway backend is acceptable only if it cannot grow memory).
- stdout of the CLI stays only the typed JSON. `--format text` output unchanged apart from stderr progress.
- Supervisor contract untouched: closed stdin, fd 3 status, `TerminateGroup` before `Wait`, 16-minute bound. Note: `exec.Cmd` copies to a non-`*os.File` writer through goroutines that `Wait` joins; the existing code already relies on that, keep the same shape (writers stay non-file values).
- Docs: `skills/do-work-toolbox/actions/ai-report.md` § Second-opinion form, one or two sentences (the launcher path shows progress and the backend's output on stderr; read the typed summary from stdout as before). `skills/do-work-toolbox/docs/ai-report-guide.md` § A Second Opinion, one plain-words sentence.
- Files you may touch: exactly the four in the REQ's `write_set`. Anything else: stop and report in the hand-back.

## TDD
RED first: add `TestSecondOpinionCaptureStreamsProgressWhileBackendRuns` in `report_second_opinion_test.go` (this exact name: the orchestrator's probe greps for it). A fake backend script prints a line, sleeps past the test-shortened heartbeat interval, writes the page and prints its envelope; assert the sink received the backend line (with its prefix) and at least one heartbeat BEFORE the backend exited (for example: the backend waits for a marker the test writes only after it has seen both on the sink, or compare ordering against a line the backend prints last). Reuse the existing fake-backend fixtures in that test file. Keep the test under a few seconds. Run it, see it fail for the right reason, record the failure text. Then GREEN, then run `go test -C skills/do-work/tools/do-work-cli -count=1 -run SecondOpinion ./internal/toolboxcommands` and `go vet` on the package, plus `gofmt -l`.
The orchestrator's probe script lives only in ROOT; the go test command above is its equivalent in your worktree.

## Hand-back (write the file, then reply with ≤ 1500 chars summarizing it)
Sections: `## Files` (each path with created/modified, plus tests), `## P-A-U` (three short paragraphs for [PLAN], [APPLY], [UNIFY] including `git diff --stat` and lint/vet/gofmt results), `## Red-Green` (test name, the failing output before, passing after), `## Decisions` (D-01… each with reasoning; ESCALATE ones add Value/Risk), `## Discovered Tasks` (out-of-scope finds, or "None"), `## Lessons Read`, `## Proposed Lesson` (one bullet in the satellite's style if something non-obvious bit you, else "None"), `## Commits` (hashes on your branch). Commit with message `[REQ-702] <what changed>`.
