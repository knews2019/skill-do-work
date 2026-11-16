# Pre-dispatch guide — run work-2026-11-16-202422, one claimed REQ (REQ-739, UR-181)

ROOT = /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 (the main checkout; you work here, never in a worktree)
RUN = ROOT/do-work/runs/work-2026-11-16-202422
CLI = `bash ROOT/skills/do-work/tools/do-work-cli.sh --repo-root ROOT --format json <command> ...` (always `--format json`; read `findings[].code`, `advance.phase`, `advance.missing_evidence`; a refusal names the missing input: fix it and call again).

You are the pre-dispatch agent for ONE REQ (named in your prompt). The coordinator (the main session) claimed it already. You take it from Triage to a ready builder brief. You do not build, you do not dispatch, you do not integrate.

## What you do, in order

Follow `ROOT/skills/do-work/actions/work.md` from Step 3 (Triage) through **Pre-Build Evidence Judgment**, for your REQ only, and read only the `work-reference.md` / `fan-out-reference.md` sections those steps name. Drive every step with the read-only `advance REQ-NNN --request-path do-work/working/<file>.md` classifier: it names the phase and the missing evidence.

1. Read the REQ fully, its UR `input.md` (`do-work/user-requests/UR-NNN/input.md`), and the audit report it came from if the UR names one. Related archived REQs are named in your REQ body; read their `## Lessons Learned` if your change touches the same files.
2. Step 3 Triage: write `route:` in frontmatter and `## Triage`. Step 3.5: open questions get best judgment, never a question to the user; record each as a decision (D-01, ...) in the REQ.
3. Estimate: `advance REQ-NNN --request-path P -- --route B --write-set 18 --subsystems 2 --acceptance 7` (integers: counts, not lists). Persist the returned `estimate:` block in frontmatter BEFORE writing `## Plan`. `effort_estimate:` contains `estimate:`: match `\nestimate:` at line start.
4. `## Plan` (Route C: Step 4 planning; Routes A/B: the short Plan or skip note work.md asks for).
5. Step 5: Required-Lessons consult (all routes); `## Exploration` for Routes B/C, with concrete file:line pointers.
6. Step 5.5 Scope (Routes B/C). Traps: every backticked token in a "Files I will touch" bullet is read as a path, so keep commands and flags out of backticks there. The scope-drift gate compares Scope bullets with Implementation Summary bullets: describe a folder move (many files) in a plain sentence, not as one bullet per file. `frontmatter set REQ-NNN write_set "[...]"` writes a quoted string; edit the `write_set:` line by hand instead.
7. Pre-Build Evidence Judgment (Routes B/C): write `RUN/REQ-NNN-preflight-probe.sh` (the focused pre-build baseline argv) and run it once yourself at the current main HEAD. **Do NOT record pre-flight with advance and do NOT run `_dev/tests/maintainer-verify.sh`.** The coordinator runs the repository gate once and records pre-flight for every Route B/C REQ serially, because pre-flight rewrites the shared `do-work/working/baseline.json`. Stop there: advance should report phase `preflight` (Route B/C).
8. Write `RUN/REQ-NNN-probe.sh`: the GREEN focused probe the integrator's test gate runs (`-- --probe-file`). Rules: `#!/usr/bin/env bash`, `set -euo pipefail`, paths relative to `git rev-parse --show-toplevel` so it checks whichever tree it runs in, under 30 s, names only this REQ's tests (never a whole Go package: `go test -run '<names>'`), asserts PASS lines (a skip is not a pass), here-strings (`grep -q ... <<<"$output"`) never `| grep -q` (the quiet-grep audit rejects it), no process substitution (advance runs probes via `sh -c`). Check it with `bash ROOT/_dev/tests/quiet-grep-pipeline-audit.sh` (read the script first to see whether it scans tracked files only; comply by hand if it cannot see an untracked probe). At base the GREEN probe is expected to fail; say so in the brief.
9. Write `RUN/REQ-NNN-brief.md`, the builder brief. Model it on `ROOT/do-work/runs/work-2026-11-16-171725/REQ-738-brief.md` (read it in full; REQ-738 shipped the header this REQ fixes), with these changes for this wave:
   - Operative name / branch: `worktree-agent-REQ-NNN-<3-to-6-word-kebab-slug>`. Worktree path: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/<operative name>`. The coordinator creates it from main HEAD right before dispatch; the builder must not create it.
   - Hand-back file: `RUN/REQ-NNN-handback.md` by its absolute main-tree path (the one main-tree file a builder may write; never staged or committed by the builder).
   - One builder runs in this run. The builder runs focused suites for the files it touched and the contract scripts the REQ names, and **never** the repository gate `_dev/tests/maintainer-verify.sh` (the integrator runs it). A 30-second per-file budget that fails once under load is rerun once and both runs are recorded.
   - Commit trailer: `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`. Commit subjects start with `[REQ-NNN]`.
   - Hand-back headings: `## Commits`, `## Files`, `## Scope changes`, `## Red-green evidence`, `## Verification`, `## P-A-U`, `## Decisions`, `## Discovered Tasks`, `## Integration seams`, `## Proposed CHANGELOG entry` (descriptive title, the "what changes for an installed copy" line), `## Proposed lesson` (one bullet or "none"), `## Lessons read`.
   - Integration seam section: from the sibling table below.
10. Record your timing spans with `record-timing-event --request-path P --category <planning|exploration-preflight> --operation "<short>" --started-at <the UTC instant you took when that span began>` (it times start to now; record each span right as it ends). Use `--run work-2026-11-16-202422` if the command asks for a run identity.

## Write boundary (strict)
- You may write ONLY: your REQ file `ROOT/do-work/working/REQ-NNN-*.md`, and files in RUN whose names start with `REQ-NNN-`. No other agent writes in the checkout while you run.
- Scratch files go in a private folder: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/predispatch-REQ-NNN/`. Never delete anything outside it.
- Never `git add`, `git commit`, `git stash`, `git checkout`, `git worktree`, or anything that changes the index or HEAD. The coordinator commits your files by exact path.
- Never edit `do-work/CHECKPOINT.md`, `do-work/working/baseline.json`, `RUN/manifest.md`, any other REQ, or any file outside `do-work/`.
- Never run `recover` (any form), `advance` without your REQ id, `advance --fan-out`, `advance --checkpoint`, or the repository gate.

## This run's member

One member: REQ-739 (the second-opinion header gets a side margin on a phone; UR-181), a follow-up of REQ-738 (archived at `do-work/archive/UR-180/REQ-738-second-opinion-direction-block-first-screen.md`; read its Decisions, Review and Lessons, and `do-work/runs/work-2026-11-16-171725/REQ-738-integration-report.md`). No sibling. A builder never edits `CHANGELOG.md`, `VERSION` or any `do-work/` path; nothing under `/Users/t2/.claude/` is ever edited.

Run-specific rules for the brief:
- **Captures for the maintainer.** The builder saves 390 x 844 captures of the header on (1) a page without its own header style (a scratch kit page and a scratch custom page) and (2) a scratch copy of the tstack v1c03 page that styles the header itself, plus the final REQ-737 page `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html` re-decorated in scratch, under `ROOT/do-work/runs/work-2026-11-16-202422/REQ-739-captures/` by absolute path, and lists them in a `## Captures for the maintainer` hand-back section. The coordinator shows them to the maintainer before integration.
- **Browser check matrix (put it in the brief):** per page at 390 x 844 and 1440 x 1000: header left text edge x, header height, map frame top y where the page has a map, sideways page scroll, and the computed margin/padding of the `aside` (to prove the page's own rule wins where it exists). One script, path and raw output in the hand-back.
- **Published pages are not rewritten.** Proofs use scratch copies; nothing under `ai-reports/` changes.
- Reuse the REQ-738 kit checks (`presentation.css` details/summary rules, `presentation.js` "Open all details").

Unattended run: never ask the user a question. An open choice becomes a D-XX decision in the REQ and one line in your final message.

## Your final message to the coordinator (under 1500 characters)
Route; operative name; the files you wrote; the last `advance` phase and code; the exact pre-flight argv the coordinator should record (`advance REQ-739 --request-path P --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-11-16-202422/REQ-739-preflight-probe.sh`), or "Route A, no pre-flight"; any open choice you made that the coordinator should know.
