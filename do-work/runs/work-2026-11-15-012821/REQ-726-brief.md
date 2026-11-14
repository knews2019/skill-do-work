# Builder Brief — REQ-726 (record the decision to adopt no external review canvas)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and finalizes. This is a **docs-only, no-release** REQ: every file you write is under `docs/decisions/`, which `.gitattributes` marks `export-ignore` and `suite/modules.tsv` does not ship.

- Operative name and branch: `worktree-agent-REQ-726-no-review-canvas-decision`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-726-no-review-canvas-decision`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with Triage (Route A) and decisions D-01 to D-04: read `do-work/working/REQ-726-record-no-review-canvas-decision.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-726-handback.md`
- Up to ten builders run at the same time in this run. Route A: no exploration, no scope section, no pre-flight.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`, `anti-slop.md` (the record is a human-facing artifact).
2. `CLAUDE.md` (maintainer preferences: simple, delete before you add).
3. The format model: `docs/decisions/records/adr-022-one-broken-pipe-does-not-stop-the-factory.md` (frontmatter fields, `Topic cluster:` line, `## Context` / `## Decision` / `## Consequences`-style body). Copy its shape, not its length.
4. The facts: `do-work/user-requests/UR-176/input.md` (the verbatim capture prompt, section "Contents of _dev/run-prompts/2026-11-15-record-no-review-canvas-decision.md"), `CHANGELOG.md` entries 0.316.0, 0.316.1 and 0.317.0 (lines 27-45 at base), `do-work/archive/UR-172/REQ-721-ai-report-closing-suggests-tinka-canvas.md`, and skim the two bundles `ai-reports/2026-11-14_1452_proposal-tinka-review-surface/` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/` for the measured numbers.
5. Required lessons: none. The pre-dispatch consult of `do-work/lessons-index.md` found no satellite whose condition covers a decision record under `docs/decisions/`.

## Decided (do not reopen)
- **D-01:** link from `docs/decisions/topics/_index_pipeline-deliverables.md`. Widen that index's title and opening sentence so it covers the live ai-report deliverable as well as the retired completion reports; list the new record as accepted (no "superseded" mark); update `updated:`. In `_master_index.md`, update that cluster's line (name, count wording, related pages) and its `updated:`.
- **D-02:** `decided: 2026-11-15`. `created:`/`updated:` = your build date.
- **D-03:** add only this record to `_master_index.md`; do not backfill ADR-017 to ADR-022. Report the gap under Discovered Tasks.
- **D-04:** do not touch `docs/decisions/_progress.md` (its "Next ADR number: ADR-017" is stale). Report it under Discovered Tasks.

## What to build
One commit, subject `[REQ-726] Record the decision to adopt no external review canvas`.

1. New `docs/decisions/records/adr-NNN-<slug>.md`, NNN = the next free number at build time (023 at base: `ls docs/decisions/records`). The slug is two words minimum and contains `review-canvas` (suggested: `adr-023-adopt-no-external-review-canvas.md`). Frontmatter as adr-022: `status: accepted`, `topic_cluster: pipeline-deliverables`, `sources:` naming both proposal bundles, `do-work/archive/UR-172/REQ-721-ai-report-closing-suggests-tinka-canvas.md`, `CHANGELOG.md` (0.316.0, 0.316.1, 0.317.0) and `do-work/user-requests/UR-176/input.md`. Body must state every fact in the REQ's Detailed Requirements: the decision; the ~5,000-pixel page as an unreadable strip tile with folded sections closed (38 to 46% of the text); 127.0.0.1-only so no phone; a live tile from Tinka's own address can write feedback; a tile from a second local port cannot (browser blocked the read and one write, Tinka refused two writes with HTTP 403) but is live only after two clicks and one page at a time; the star/vote/note marks move feedback out of chat and add steps; the direction taken instead (`index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` in the second bundle); and the two reopen conditions. Name Tinka's repository (github.com/bladnman/tinka) once.
2. `docs/decisions/log.md`: a new dated section at the end, `## [2026-11-15] <plain title>`, one `- Accepted [[<slug>]]: ...` bullet in the house style of the adr-022 entry (line 122 at base).
3. `docs/decisions/_master_index.md` and `docs/decisions/topics/_index_pipeline-deliverables.md` per D-01, linking `[[<slug>]]` the way existing records are linked (the topic index uses `[[slug]] — [ADR-NNN](../records/<slug>.md): one-line summary`, plus a `related:` frontmatter entry).

**Trap:** `_dev/tests/audit-lockins.sh` scans every topic index's `sources:` list for dead paths: an entry with no space and no `*` must be a substring of some tracked path. If you add a source line to the topic index, make it a tracked path or a `CHANGELOG.md (0.317.0 ...)` style entry with spaces. Record files and `log.md` are not scanned.

**Not yours:** anything under `skills/`, `suite/`, `ai-reports/` or `do-work/`; `CHANGELOG.md`; `VERSION`; `docs/decisions/_progress.md`; other ADRs.

## Rules
- Write only the four files above. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-726]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-726-probe.sh` with your worktree as the current directory prints `REQ-726 probe PASS (<slug>)`. It fails at base by design ("no docs/decisions/records/adr-*review-canvas*.md record"); the pre-dispatch agent confirmed both the failure at base and a pass on a simulated record.
2. `bash _dev/tests/contract-regressions.sh` exits 0 (about 30 s; it passed at base).
3. `bash _dev/tests/audit-lockins.sh` exits 0 (about 1 s; it guards the topic-index `sources:` paths).
4. `git diff --stat main...HEAD -- skills/ suite/ ai-reports/ do-work/` prints nothing, and `git diff --stat main...HEAD` lists only the four `docs/decisions/` files.
5. Read the record once as the maintainer would: plain words, every Detailed Requirements fact present, no fact the REQ does not support.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject.
- `## Files` — each file with created/modified.
- `## Scope changes` — "none" if none.
- `## Red-green evidence` — the GREEN probe failing at base and passing at your tip, each with its output line.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-05 onward (D-01 to D-04 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks` — at least D-03's index backfill gap and D-04's stale `_progress.md` line.
- `## Integration seams` — see below.
- `## Proposed CHANGELOG entry` — "none (no shipped file changes)". `docs/decisions` is `export-ignore` in `.gitattributes` and not a source in `suite/modules.tsv`.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read` — "none required (no matching satellite)", plus any you chose to read.

## Integration seams
No sibling in this run edits `docs/decisions/`. REQ-732 (allow `docs/run-prompts/` in the staged-skills contract) edits `_dev/tests/staged-skills-contract.sh` and `.gitattributes`, which you do not touch. If the integrator meets a conflict, it can only be the ADR number: if another branch landed an `adr-023-*` first, renumber yours to the next free number and update the four files' links. Say in this section which number you used.

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
