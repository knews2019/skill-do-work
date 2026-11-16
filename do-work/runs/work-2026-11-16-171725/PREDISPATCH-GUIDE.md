# Pre-dispatch guide — run work-2026-11-16-171725, two claimed REQs (REQ-737 and REQ-738, both of UR-180)

ROOT = /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 (the main checkout; you work here, never in a worktree)
RUN = ROOT/do-work/runs/work-2026-11-16-171725
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
9. Write `RUN/REQ-NNN-brief.md`, the builder brief. Model it on `ROOT/do-work/runs/work-2026-11-15-012821/REQ-727-brief.md` (read it in full; REQ-727 is the parent of this run's work), with these changes for this wave:
   - Operative name / branch: `worktree-agent-REQ-NNN-<3-to-6-word-kebab-slug>`. Worktree path: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/<operative name>`. The coordinator creates it from main HEAD right before dispatch; the builder must not create it.
   - Hand-back file: `RUN/REQ-NNN-handback.md` by its absolute main-tree path (the one main-tree file a builder may write; never staged or committed by the builder).
   - Two builders run at the same time in this run. The builder runs focused suites for the files it touched and the contract scripts the REQ names, and **never** the repository gate `_dev/tests/maintainer-verify.sh` (the integrator runs it). A 30-second per-file budget that fails once under load is rerun once and both runs are recorded.
   - Commit trailer: `Co-Authored-By: Claude Fable 5.2 <noreply@anthropic.com>`. Commit subjects start with `[REQ-NNN]`.
   - Hand-back headings: `## Commits`, `## Files`, `## Scope changes`, `## Red-green evidence`, `## Verification`, `## P-A-U`, `## Decisions`, `## Discovered Tasks`, `## Integration seams`, `## Proposed CHANGELOG entry` (descriptive title, the "what changes for an installed copy" line), `## Proposed lesson` (one bullet or "none"), `## Lessons read`.
   - Integration seam section: from the sibling table below.
10. Record your timing spans with `record-timing-event --request-path P --category <planning|exploration-preflight> --operation "<short>" --started-at <the UTC instant you took when that span began>` (it times start to now; record each span right as it ends). Use `--run work-2026-11-16-171725` if the command asks for a run identity.

## Write boundary (strict)
- You may write ONLY: your REQ file `ROOT/do-work/working/REQ-NNN-*.md`, and files in RUN whose names start with `REQ-NNN-`. One other pre-dispatch agent writes the sibling REQ's files at the same time; never touch them.
- Scratch files go in a private folder: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/predispatch-REQ-NNN/`. Never delete anything outside it.
- Never `git add`, `git commit`, `git stash`, `git checkout`, `git worktree`, or anything that changes the index or HEAD. The coordinator commits your files by exact path.
- Never edit `do-work/CHECKPOINT.md`, `do-work/working/baseline.json`, `RUN/manifest.md`, any other REQ, or any file outside `do-work/`.
- Never run `recover` (any form), `advance` without your REQ id, `advance --fan-out`, `advance --checkpoint`, or the repository gate.

## The run's members (for the brief's Integration seam section)

Both are claimed and build in parallel worktrees from the same main HEAD (VERSION 0.324.0). Integration is serial: **REQ-738 first, REQ-737 last**. The integrator writes `CHANGELOG.md` and `VERSION`; **a builder never edits `CHANGELOG.md`, `VERSION`, or any `do-work/` path**, even when the REQ's `write_set` lists them (the builder proposes the changelog entry in its hand-back). Nothing under `/Users/t2/.claude/` is ever edited (read only).

| REQ | What | Shared files with the sibling |
| --- | --- | --- |
| REQ-737 | Preset text only: the shipped `relation-map` preset keeps box titles readable after a click on a phone; one real proof page built on the tstack bundle and measured | none in source. Seam in behaviour: its proof page is decorated by the launcher that REQ-738 changes |
| REQ-738 | Go, tdd: the `ai-report second-opinion` launcher shows a short header and folds the full direction, so the page content starts on the first phone screen; contract text in `ai-report.md` and `ai-report-guide.md` | none in source |

Run-specific rules for the briefs:

- **The maintainer must see phone captures before either REQ is called done** (both REQs say so). The builder saves its final captures (390 x 844, and 1440 x 1000 where the REQ measures it) as PNG files under `ROOT/do-work/runs/work-2026-11-16-171725/REQ-NNN-captures/` by absolute path (a second main-tree location the builder may write, never staged by the builder), lists each file with one line on what it shows in a `## Captures for the maintainer` hand-back section, and names the layout form it chose and the forms it rejected. The coordinator shows them to the maintainer before integration starts. Any other Playwright capture is deleted before hand-back (captures written by the Playwright tools land in the main tree's gitignored `.playwright-mcp/`; delete those and record it as a D-XX).
- **Browser check matrix (put it in the brief).** Page REQs took two or three review rounds in the last run (label overlap, phone width, focus order). The builder measures with ONE script that clicks every shown box from each starting view (first view, and after "Show everything") at 390 x 844 and 1440 x 1000, and reports per click: clicked title px, each shown related title px, boxes out of frame, sideways page scroll, plus Tab focus in view, the no-JavaScript view, and the list under the map. The script path and its raw output go in the hand-back.
- **REQ-737 seam.** REQ-737's builder branches from a main that still has the 0.324.0 header block. It builds and measures its proof page there, and says in `## Integration seams` exactly how the page was produced (the full `ai-report second-opinion` argv, the bundle, the page id) and how the measurement script is run, so the REQ-737 integrator can re-decorate the page with the REQ-738 launcher after REQ-738 lands and re-run the same script (REQ-737 Builder Guidance: "Re-measure after finalize, because the launcher's direction block moves the map"). REQ-737's GREEN criteria are about title size after a click and do not depend on the block; the builder still reports the map top y it saw.
- **REQ-737 preset size.** About 3,500 bytes at most; no check removed to save bytes. The GREEN probe asserts the byte bound, that all 20 checks (or more) are still numbered, and the existing shipped-presets Go test by name.
- **REQ-738 seam.** Published pages under `ai-reports/` are not rewritten; the proof is a scratch copy, kept out of the commit. The brief names the kit files to check against (`presentation.css` lines 196-199, `presentation.js` lines 44-50) and the old-block strip case (a page decorated by 0.324.0, decorated again by the new code).
- **Session second opinion.** REQ-737's proof page is made through a Session second opinion, where the builder agent itself is the backend. Read how REQ-727's builder did it (`ROOT/do-work/runs/work-2026-11-15-012821/REQ-727-handback.md` and `REQ-727-brief.md`) and the `#session-is-the-backend` lesson in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (grep for the slug; do not read the whole file), and put the exact steps in the brief.

Unattended run: never ask the user a question. An open choice becomes a D-XX decision in the REQ and one line in your final message.

## Your final message to the coordinator (under 1500 characters)
Route; operative name; the files you wrote; the last `advance` phase and code; the exact pre-flight argv the coordinator should record (`advance REQ-NNN --request-path P --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-11-16-171725/REQ-NNN-preflight-probe.sh`); any open choice you made that the coordinator should know; any seam you found that the table above missed.
