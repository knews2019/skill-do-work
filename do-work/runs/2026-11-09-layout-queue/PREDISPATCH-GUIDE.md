# Pre-dispatch guide — run 2026-11-09-layout-queue, wave REQ-716 + REQ-717

ROOT = /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 (the main checkout; you work here, never in a worktree)
RUN = ROOT/do-work/runs/2026-11-09-layout-queue
CLI = `bash ROOT/skills/do-work/tools/do-work-cli.sh --repo-root ROOT --format json <command> ...` (always `--format json`; read `findings[].code`, `advance.phase`, `advance.missing_evidence`; a refusal names the missing input: fix it and call again).

You are the pre-dispatch agent for ONE REQ (named in your prompt). The coordinator (the main session) claimed it already. You take it from Triage to a ready builder brief. You do not build, you do not dispatch, you do not integrate.

## What you do, in order

Follow `ROOT/skills/do-work/actions/work.md` from Step 3 (Triage) through **Pre-Build Evidence Judgment**, for your REQ only, and read only the `work-reference.md` / `fan-out-reference.md` sections those steps name. Drive every step with the read-only `advance REQ-NNN --request-path do-work/working/<file>.md` classifier: it names the phase and the missing evidence.

1. Read the REQ fully, its UR `input.md` (`do-work/user-requests/UR-NNN/input.md`), and the audit report it came from if the UR names one. REQ-713 (folder-structure-check sees recipes with default parameter values) is archived under `do-work/archive/UR-165/`; REQ-714 and REQ-715 under `do-work/archive/UR-166/` and `UR-167/`. Read their `## Lessons Learned` if your change touches the same files.
2. Step 3 Triage: write `route:` in frontmatter and `## Triage`. Step 3.5: open questions get best judgment, never a question to the user; record each as a decision (D-01, ...) in the REQ.
3. Estimate: `advance REQ-NNN --request-path P -- --route B --write-set 18 --subsystems 2 --acceptance 7` (integers: counts, not lists). Persist the returned `estimate:` block in frontmatter BEFORE writing `## Plan`. `effort_estimate:` contains `estimate:`: match `\nestimate:` at line start.
4. `## Plan` (Route C: Step 4 planning; Routes A/B: the short Plan or skip note work.md asks for).
5. Step 5: Required-Lessons consult (all routes); `## Exploration` for Routes B/C, with concrete file:line pointers.
6. Step 5.5 Scope (Routes B/C). Traps: every backticked token in a "Files I will touch" bullet is read as a path, so keep commands and flags out of backticks there. The scope-drift gate compares Scope bullets with Implementation Summary bullets: describe a folder move (many files) in a plain sentence, not as one bullet per file. `frontmatter set REQ-NNN write_set "[...]"` writes a quoted string; edit the `write_set:` line by hand instead.
7. Pre-Build Evidence Judgment (Routes B/C): write `RUN/REQ-NNN-preflight-probe.sh` (the focused pre-build baseline argv) and run it once yourself at the current main HEAD. **Do NOT record pre-flight with advance and do NOT run `_dev/tests/maintainer-verify.sh`.** The coordinator runs the repository gate once and records pre-flight for both REQs serially, because pre-flight rewrites the shared `do-work/working/baseline.json`. Stop there: advance should report phase `preflight` (Route B/C).
8. Write `RUN/REQ-NNN-probe.sh`: the GREEN focused probe the integrator's test gate runs (`-- --probe-file`). Rules: `#!/usr/bin/env bash`, `set -euo pipefail`, paths relative to `git rev-parse --show-toplevel` so it checks whichever tree it runs in, under 30 s, names only this REQ's tests (never a whole Go package: `go test -run '<names>'`), asserts PASS lines (a skip is not a pass), here-strings (`grep -q ... <<<"$output"`) never `| grep -q` (the quiet-grep audit rejects it), no process substitution (advance runs probes via `sh -c`). Check it with `bash ROOT/_dev/tests/quiet-grep-pipeline-audit.sh` (read the script first to see whether it scans tracked files only; comply by hand if it cannot see an untracked probe). At base the GREEN probe is expected to fail; say so in the brief.
9. Write `RUN/REQ-NNN-brief.md`, the builder brief. Model it on `RUN/REQ-714-brief.md` (read it in full), with these changes for this wave:
   - Operative name / branch: `worktree-agent-REQ-NNN-<3-to-6-word-kebab-slug>`. Worktree path: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/<operative name>`. The coordinator creates it from main HEAD right before dispatch; the builder must not create it.
   - Hand-back file: `RUN/REQ-NNN-handback.md` by its absolute main-tree path (the one main-tree file a builder may write; never staged or committed by the builder).
   - Two builders run at the same time in this wave. The builder runs focused suites for the files it touched and the contract scripts the REQ names, and **never** the repository gate `_dev/tests/maintainer-verify.sh` (the integrator runs it). A 30-second per-file budget that fails once under load is rerun once and both runs are recorded.
   - Commit trailer: `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`. Commit subjects start with `[REQ-NNN]`.
   - Hand-back headings: `## Commits`, `## Files`, `## Scope changes`, `## Red-green evidence`, `## Verification`, `## P-A-U`, `## Decisions`, `## Discovered Tasks`, `## Integration seams`, `## Proposed CHANGELOG entry` (descriptive title, the "what changes for an installed copy" line), `## Proposed lesson` (one bullet or "none"), `## Lessons read`.
   - Integration seam section: from the sibling table below.
10. Record your timing spans with `record-timing-event --request-path P --category <planning|exploration-preflight> --operation "<short>" --started-at <the UTC instant you took when that span began>` (it times start to now; record each span right as it ends). Use `--run 2026-11-09-layout-queue` if the command asks for a run identity.

## Write boundary (strict)
- You may write ONLY: your REQ file `ROOT/do-work/working/REQ-NNN-*.md`, and files in RUN whose names start with `REQ-NNN-`. One other pre-dispatch agent writes the other REQ's files at the same time; never touch them.
- Scratch files go in a private folder: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/008bccaf-aa66-4644-8fb6-cac0974a4198/scratchpad/predispatch-REQ-NNN/`. Never delete anything outside it.
- Never `git add`, `git commit`, `git stash`, `git checkout`, `git worktree`, or anything that changes the index or HEAD. The coordinator commits your files by exact path.
- Never edit `do-work/CHECKPOINT.md`, `do-work/working/baseline.json`, `RUN/manifest.md`, any other REQ, or any file outside `do-work/`.
- Never run `recover` (any form), `advance` without your REQ id, `advance --fan-out`, `advance --checkpoint`, or the repository gate.

## The wave's members (for the brief's Integration seam section)

Source: the folder-layout audit of this repository on 2026-11-09. Recorded exceptions (do not move): `VERSION`, `skills/` and `suite/` stay at the root; `_dev/` stays as the maintainer-only folder. REQ-713 to REQ-715 are released or finalized. Both members below are claimed and build in parallel; integration is serial, REQ-716 first.

| REQ | What | Likely shared files |
| --- | --- | --- |
| REQ-716 | Move the five root installer scripts from `tools/` to `cmd/` with `git mv`. Riskier: it must prove an update from an already-installed older version still works (an older updater may run `tools/<script>` from inside the newly downloaded archive), or keep a thin `tools/` forwarder; record which. The README bootstrap, the Go update path, about 34 shipped skill files and the contract tests name these scripts. `.gitattributes` and `suite/modules.tsv` may name `tools/`. | `_dev/tests/staged-skills-contract.sh` (its `legacy_runtime_paths` entries under `tools/`), the root `justfile` only if a recipe names `tools/`, `README.md` |
| REQ-717 | The suite-managed recipes in `skills/do-work-board/justfile.template` (`run-kanban` 3 body lines, `kanban-static` 2) become one-line calls to scripts shipped with do-work-board; this repository's root `justfile` managed block follows. Behavior-preserving. | the root `justfile`, `_dev/tests/staged-skills-contract.sh` if it lists shipped scripts |

Each builder edits only its own lines of a shared file; the integrator resolves overlap at the merge. A REQ-716 builder must not rename anything REQ-717 adds, and the reverse. If your REQ's new script would live under a `tools/` folder, say in the brief whether the layout rule ("scripts live in `cmd/`") applies to a skill package's own `tools/` folder, using what `folder-structure-check` actually reports.

## Your final message to the coordinator (under 1500 characters)
Route; operative name; the files you wrote; the last `advance` phase and code; the exact pre-flight argv the coordinator should record (`advance REQ-NNN --request-path P --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/2026-11-09-layout-queue/REQ-NNN-preflight-probe.sh`); any open choice you made that the coordinator should know; any seam you found that the table above missed.
