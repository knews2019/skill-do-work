# Pre-dispatch guide — run work-2026-10-19-152525

ROOT = /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 (the main checkout; you work here, never in a worktree)
RUN = ROOT/do-work/runs/work-2026-10-19-152525
CLI = `bash ROOT/skills/do-work/tools/do-work-cli.sh --repo-root ROOT --format json <command> ...` (always `--format json`; read `findings[].code`, `advance.phase`, `advance.missing_evidence`; a refusal names the missing input: fix it and call again).

You are the pre-dispatch agent for ONE REQ (named in your prompt). The coordinator (the main session) claimed it already. You take it from Triage to a ready builder brief. You do not build, you do not dispatch, you do not integrate.

## What you do, in order

Follow `ROOT/skills/do-work/actions/work.md` from Step 3 (Triage) through **Pre-Build Evidence Judgment**, for your REQ only, and read only the `work-reference.md` / `fan-out-reference.md` sections those steps name. Drive every step with the read-only `advance REQ-NNN --request-path do-work/working/<file>.md` classifier: it names the phase and the missing evidence. Concretely:

1. Read the REQ fully, its UR `input.md` (path from `user_request:`; `do-work/user-requests/UR-NNN/input.md`), and every `related:` REQ it cites (archived ones live under `do-work/archive/`).
2. Step 3 Triage: write `route:` in frontmatter and `## Triage`. Step 3.5: open questions get best judgment, never a question to the user; record each as a decision in the REQ.
3. Estimate: `advance REQ-NNN --request-path P -- --route B --write-set 18 --subsystems 2 --acceptance 7` (integers: counts, not lists; `-- --trivial` is accepted for effort-mechanical). Persist the returned `estimate:` block in frontmatter BEFORE writing `## Plan` (a Plan written first refuses with "later lifecycle evidence exists before estimate.p50_active_minutes"). `effort_estimate:` contains `estimate:`: match `\nestimate:` at line start.
4. `## Plan` (Route C: Step 4 planning; Routes A/B: the short Plan or skip note work.md asks for).
5. Step 5: Required-Lessons consult (all routes); `## Exploration` for Routes B/C.
6. Step 5.5 Scope (Routes B/C). Trap: every backticked token in a "Files I will touch" bullet is read as a path. Keep commands and flags out of backticks there.
7. Pre-Build Evidence Judgment (Routes B/C): write `RUN/REQ-NNN-preflight-probe.sh` (the focused pre-build baseline argv) and run it once yourself to see it behave as expected at the current main HEAD. **Do NOT record pre-flight with advance and do NOT run `_dev/tests/maintainer-verify.sh`.** The coordinator runs the repository gate once, then records pre-flight for every Route B/C REQ serially, because pre-flight rewrites the shared `do-work/working/baseline.json`. Stop there: advance should report phase `preflight` (Route B/C) or the implementation phase (Route A).
8. Write `RUN/REQ-NNN-probe.sh`: the GREEN focused probe the integrator's test gate runs (`-- --probe-file`). Rules: `#!/usr/bin/env bash`, `set -euo pipefail`, reads paths relative to `git rev-parse --show-toplevel` so it checks whichever tree it runs in, finishes under 30 s, names only this REQ's tests (never a whole Go package: `go test -run '<names>'`), asserts `--- PASS:` lines (a skip is not a pass), uses here-strings (`grep -q ... <<<"$output"`) never `| grep -q` (the contract quiet-grep audit rejects that), no process substitution (advance runs probes via `sh -c`). Run `bash ROOT/_dev/tests/quiet-grep-pipeline-audit.sh` against your probe files if it accepts paths; otherwise read it and comply by hand. At base the GREEN probe is expected to fail (the code does not exist yet); say so in the brief.
9. Write `RUN/REQ-NNN-brief.md`, the builder brief. Model it on `ROOT/do-work/runs/work-2026-10-10-100748/REQ-679-brief.md` (read it in full): header bullets with absolute paths (worktree, branch, REQ, UR, hand-back file), Rules to load first, The change (decided, do not reopen; concrete file:line pointers from your Exploration), Anti-bloat (YAGNI) block, Write boundary, Hard rules (copy them, adjusting the REQ id; keep "never touch do-work/ except the hand-back file", "commit subject starts with `[REQ-NNN]`", "never run the gate / advance / recover", "no CHANGELOG/VERSION/mirror edits"), Integration seam (from the sibling table below), Proof to run and record (RED then GREEN from the REQ's Red-Green Proof, if it has one), Verify before hand-back (run `bash RUN/REQ-NNN-probe.sh` from the worktree root; gofmt/go vet for Go), Hand-back format (same list as the example, including `## Decisions`, `## Discovered Tasks`, a proposed CHANGELOG entry with a descriptive title, and a proposed lesson bullet or "none").
   - Operative name / branch: `worktree-agent-REQ-NNN-<3-to-6-word-kebab-slug>`. Worktree path: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/<operative name>`. The coordinator creates it with `git worktree add -b` from main HEAD right before dispatch; the builder must not create it.
   - Load note for the brief: one builder runs at a time in this run (the two members are serialized because they edit the same Go file); a wall-time budget or browser probe that fails once under load is rerun once and both runs are recorded.
   - Go builds: the brief tells the builder to run focused `go test -run` and `go vet` on the touched package only, never the repository gate.
10. Record your timing spans with `record-timing-event --request-path P --category <planning|exploration-preflight> --operation "<short>" --started-at <the UTC instant you took when that span began>` (it times start to now; record each span right as it ends).

## Write boundary (strict)
- You may write ONLY: your REQ file `ROOT/do-work/working/REQ-NNN-*.md`, and files in RUN whose names start with `REQ-NNN-`. One other pre-dispatch agent writes the other REQ files and run-dir files at the same time.
- Never `git add`, `git commit`, `git stash`, `git checkout`, `git worktree`, or anything that changes the index or HEAD. The coordinator commits your files by exact path.
- Never edit `do-work/CHECKPOINT.md`, `do-work/working/baseline.json`, `RUN/manifest.md`, any other REQ, or any file outside `do-work/`.
- Never run `recover` (any form), `advance` without your REQ id, `advance --fan-out`, `advance --checkpoint`, or the repository gate.

## The run's members (for the brief's Integration seam section)

Two REQs of UR-158 (both from the same user feedback about `ai-report second-opinion`). Both change `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` and `report_second_opinion_test.go`, so the coordinator **serialized** them: REQ-701 is built and integrated first; REQ-700's builder is dispatched from main only after REQ-701 has been released.

| REQ | What | Order | Likely shared files |
| --- | --- | --- | --- |
| REQ-701 | post-run model check never deletes a finished page; `claude` model resolution = requested model when present in `modelUsage`, subagent entries recorded not ambiguous; unverified-name publication + typed finding | first | `report_second_opinion.go` (`secondOpinionRun` cleanup defer, `secondOpinionMetadata`), its test, `lessons-do-work-cli.md` |
| REQ-700 | Session form never execs the agent's own CLI; CLI split into a prepare phase (bundle, pre-checks, reserved names, filled prompt) and a finalize phase (preservation, outside-bundle Git, backlink, judge, banner, reply, typed summary) the agent calls around its own in-session editorial pass; action prose + user guide updated | second | same Go file + test, `skills/do-work-toolbox/actions/ai-report.md`, `skills/do-work-toolbox/docs/ai-report-guide.md`, `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md`, result model types |

**Pre-dispatch for REQ-700 (second member):** your Exploration and brief are written against today's main, but the builder will branch from main AFTER REQ-701 lands, so (a) cite functions by name and not only by line number, (b) say in the brief that REQ-701's change (no deletion of a finished page, the `claude` model-resolution rule) is already on main and the finalize phase must keep it, (c) the coordinator will append a "Coordinator note" with REQ-701's landed diff summary before dispatch. Still write `REQ-700-preflight-probe.sh` and `REQ-700-probe.sh` now; the coordinator records pre-flight after REQ-701's release.

## Your final message to the coordinator (under 1500 characters)
Route; operative name; the files you wrote; the last `advance` phase and code; the exact pre-flight argv the coordinator should record (`advance REQ-NNN --request-path P --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-10-19-152525/REQ-NNN-preflight-probe.sh`, or "none, Route A"); any open choice you made that the coordinator should know; any seam you found that the table above missed.
