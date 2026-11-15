# Pre-dispatch guide — run work-2026-11-15-012821, ten claimed REQs (REQ-723 to REQ-733 without REQ-727)

ROOT = /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 (the main checkout; you work here, never in a worktree)
RUN = ROOT/do-work/runs/work-2026-11-15-012821
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
9. Write `RUN/REQ-NNN-brief.md`, the builder brief. Model it on `ROOT/do-work/runs/2026-11-09-layout-queue/REQ-716-brief.md` (read it in full), with these changes for this wave:
   - Operative name / branch: `worktree-agent-REQ-NNN-<3-to-6-word-kebab-slug>`. Worktree path: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/<operative name>`. The coordinator creates it from main HEAD right before dispatch; the builder must not create it.
   - Hand-back file: `RUN/REQ-NNN-handback.md` by its absolute main-tree path (the one main-tree file a builder may write; never staged or committed by the builder).
   - Up to ten builders run at the same time in this run. The builder runs focused suites for the files it touched and the contract scripts the REQ names, and **never** the repository gate `_dev/tests/maintainer-verify.sh` (the integrator runs it). A 30-second per-file budget that fails once under load is rerun once and both runs are recorded.
   - Commit trailer: `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`. Commit subjects start with `[REQ-NNN]`.
   - Hand-back headings: `## Commits`, `## Files`, `## Scope changes`, `## Red-green evidence`, `## Verification`, `## P-A-U`, `## Decisions`, `## Discovered Tasks`, `## Integration seams`, `## Proposed CHANGELOG entry` (descriptive title, the "what changes for an installed copy" line), `## Proposed lesson` (one bullet or "none"), `## Lessons read`.
   - Integration seam section: from the sibling table below.
10. Record your timing spans with `record-timing-event --request-path P --category <planning|exploration-preflight> --operation "<short>" --started-at <the UTC instant you took when that span began>` (it times start to now; record each span right as it ends). Use `--run work-2026-11-15-012821` if the command asks for a run identity.

## Write boundary (strict)
- You may write ONLY: your REQ file `ROOT/do-work/working/REQ-NNN-*.md`, and files in RUN whose names start with `REQ-NNN-`. Nine other pre-dispatch agents write the other REQs' files at the same time; never touch them.
- Scratch files go in a private folder: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/predispatch-REQ-NNN/`. Never delete anything outside it.
- Never `git add`, `git commit`, `git stash`, `git checkout`, `git worktree`, or anything that changes the index or HEAD. The coordinator commits your files by exact path.
- Never edit `do-work/CHECKPOINT.md`, `do-work/working/baseline.json`, `RUN/manifest.md`, any other REQ, or any file outside `do-work/`.
- Never run `recover` (any form), `advance` without your REQ id, `advance --fan-out`, `advance --checkpoint`, or the repository gate.

## The run's members (for the brief's Integration seam section)

All ten are claimed and build in parallel worktrees from the same main HEAD. Integration is serial, one integrator at a time, in hand-back order. The integrator writes `CHANGELOG.md` and `VERSION`; **a builder never edits `CHANGELOG.md`, `VERSION`, or any `do-work/` path**, even when the REQ's `write_set` lists them (the builder proposes the changelog entry in its hand-back).

| REQ | What | Shared files with siblings |
| --- | --- | --- |
| REQ-723 | Go: the archive-collision helper skips assets folders | none expected |
| REQ-724 | Move the maintainer's global standing rules into the skill's crew-member files | `skills/do-work/crew-members/*.md` (no sibling edits them) |
| REQ-725 | Report only: sort the local memory notes by destination. No release | none; the report lands under `ai-reports/` |
| REQ-726 | Decision record: no external review canvas. No release | `docs/decisions/*` only |
| REQ-728 | Go: `ai-report-judge` renders one named page of a bundle | `report_judge.go`, `report_judge_test.go` (REQ-729), `skills/do-work-toolbox/actions/ai-report.md` and `docs/ai-report-guide.md` (REQ-729, REQ-730, REQ-731) |
| REQ-729 | Go or docs: `ai-report-judge` resolves or documents a relative `--out` | same Go files as REQ-728; `ai-report.md` |
| REQ-730 | Go: the report catalog lists each bundle's second-opinion pages | `report_index.go`; `ai-report.md`, `ai-report-guide.md` |
| REQ-731 | Prose: `ai-report` prints the catalog rebuild command when a new report is missing from the catalog | `ai-report.md`, `ai-report-guide.md` |
| REQ-732 | Allow `docs/run-prompts/` in the staged-skills contract and move the run prompts there | `_dev/tests/staged-skills-contract.sh`, `.gitattributes` |
| REQ-733 | Rebuild the relation map page inside one `ai-reports/` bundle. No release | the bundle's `index.html` (REQ-734 builds later, after REQ-733 lands) |

Rules for a shared file: each builder edits only the lines its REQ needs, adds new text as its own paragraph or bullet (never rewrites a neighbour's paragraph for style), adds new Go tests as new functions at the end of the test file, and names in `## Integration seams` every shared file it touched with the function or heading it changed. REQ-728 and REQ-729 builders: say in the seam section which functions of `report_judge.go` you changed, so the second integrator can resolve the overlap.

A no-release REQ (only `ai-reports/`, `docs/decisions/` or other unshipped paths change) says so in its brief: `## Proposed CHANGELOG entry` is then "none (no shipped file changes)". Check "shipped" against `.gitattributes` export-ignore and `suite/modules.tsv`, not by guess.

Route A REQs skip Exploration, Scope and pre-flight: write Triage, estimate, the Plan skip note, the lessons consult, the GREEN probe and the brief, then stop. Say "Route A, no pre-flight" in your final message.

Unattended run: never ask the user a question. An open choice becomes a D-XX decision in the REQ and one line in your final message.

## Your final message to the coordinator (under 1500 characters)
Route; operative name; the files you wrote; the last `advance` phase and code; the exact pre-flight argv the coordinator should record (`advance REQ-NNN --request-path P --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-11-15-012821/REQ-NNN-preflight-probe.sh`); any open choice you made that the coordinator should know; any seam you found that the table above missed.


## Second pass (2026-11-15): REQ-735 and REQ-736 from UR-179

The eleven members in the table above are finalized; ignore that table. This pass has two members, claimed together, built in parallel, integrated one at a time (REQ-736 first, REQ-735 last).

| REQ | What | Shared files with the sibling |
| --- | --- | --- |
| REQ-735 | Decide, rule by rule, which of three moved standing rules load in ordinary sessions; instruction-file change in `skills/do-work/crew-members/` and the package copies under `skills/do-work-toolbox/crew-members/` and `skills/do-work-knowledge/crew-members/` | none |
| REQ-736 | One line: the judge line in `skills/do-work-toolbox/actions/help.md` shows the `--page` option | none |

Every rule of this guide still applies, including: a builder never edits `CHANGELOG.md`, `VERSION` or a `do-work/` path; never ask the user; nothing under `/Users/t2/.claude/` is ever edited (read only). REQ-724 (archived at `do-work/archive/UR-175/`) is the parent of REQ-735: read its Decisions, Review and Lessons sections and `do-work/runs/work-2026-11-15-012821/REQ-724-integration-report.md`. Its review required that each new line also lands in every package copy of the same crew file; put that in the REQ-735 brief and probe.


## Third pass (2026-11-15): REQ-727 alone

Every member in the tables above is finalized; ignore those tables. This pass has one member, REQ-727 (ship a relation-map preset for `ai-report second-opinion`), unblocked by the maintainer through clarify: the maintainer confirmed the relation map as a standard option after both example pages existed. No sibling builds at the same time. Every rule of this guide still applies: a builder never edits `CHANGELOG.md`, `VERSION` or a `do-work/` path; never ask the user; nothing under `/Users/t2/.claude/` is ever edited.

Read before planning: the two dependency sources, `do-work/archive/UR-177/REQ-733-rebuild-relation-map-page.md` (the map page; its Review took three rounds) and `do-work/archive/UR-178/REQ-734-proposal-table-and-map-comparison.md` (the table page and the reviewer's map-versus-table comparison), plus their integration reports in this run directory. The example pages are `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-*.html` (map) and `index-v1c04-*.html` (table). Also read how the shipped presets are stored and tested (REQ-703 and REQ-705 in `do-work/archive/UR-160/`).

The preset must make a generated map avoid the defects the example needed three review rounds to remove, and the open findings that review left: relation words visible only after a click; labels covering box titles; words broken inside narrow labels; a forced zoom that pushes related boxes out of view; Tab focus on off-screen boxes; a phone view that needs panning; no relation words without JavaScript. Turn each into a checkable instruction in the preset text where a preset can carry it, and say in a decision which ones a preset cannot guarantee.

This REQ closes UR-176 (its other six members are archived flat). Say so in the brief's notes for the integrator; the builder does nothing about it.
