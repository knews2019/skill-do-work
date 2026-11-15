# REQ-725 Hand-back — sort the local memory notes by destination

Branch `worktree-agent-REQ-725-memory-note-destination-report`, tip `f90c21fa`.

## Commits
- `f90c21fa` [REQ-725] sort the local memory notes by destination (proposal report)

## Files
- `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/index.html` — added (generated from a scratch row table, then hand-checked)
- `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/presentation.css` — added (kit copy, unchanged)
- `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/presentation.js` — added (kit copy, unchanged)

## Scope changes
None. `ai-reports/catalog.json` and `ai-reports/index.html` were not regenerated (integration seam).

## Red-green evidence
- Base (`df1bf423`, before the bundle existed): `REQ-725 probe FAIL: expected exactly one ai-reports/*_proposal-REQ-725-*/index.html, found 0` (exit 1).
- Tip `f90c21fa`: `REQ-725 probe PASS: 47 notes, one row each, in ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/index.html` (exit 0).

## Verification
1. GREEN probe from the worktree: PASS, 47 notes (above).
2. `git diff --stat main...HEAD`: 3 files, all inside the new bundle (677 insertions).
3. D-03 snapshots: `shasum -a 256 <memory>/*` before the first read and after the last; `diff` shows no difference (48 files: `MEMORY.md` + 47 notes). The memory folder was never written.
4. Step 7 render check: `ai-report-judge` verdict `pass`, no findings. I inspected wide light, wide dark and phone dark captures (header, counts, options, table). The judge output folder was deleted afterwards. No `FILL:` marker remains.

## P-A-U
- PLAN: read all 47 notes as data, check each claim with fixed-string searches at `df1bf423`, apply the REQ-724 and REQ-732 seams, write one proposal bundle with one `data-memory-note` row per note.
- APPLY: the bundle only; rows sorted by destination, then file name; counts 5 + 13 + 29 = 47 shown at the top; options O1 (whole list), O2 (smallest change: delete only), O3 (do nothing), each priced; capture lines printed only.
- UNIFY: `git diff --stat main...HEAD` lists only `index.html`, `presentation.css`, `presentation.js` in the bundle; checked that the kit files are byte copies, the probe passes, the judge passes, and the memory snapshots match.

## Decisions
- **D-04** (DECIDE & STATE): `ask-tool-when-action-says-ask.md` is **delete** ("already recorded at communication-style.md and clear-questions.md (after REQ-724)"), not a move. Its rule is REQ-724's G2, and G2 has no autonomy exception, so the incident it records is covered.
- **D-05** (DECIDE & STATE): `one-orchestrator-per-queue.md` is **move to the skill**, but only its path-limited commit rule (`git commit … -- <paths>`, target `skills/do-work/actions/commit.md`). Its incident log and "never probe for sessions" stance are not moved. Reason: `commit.md:125-127` and `144-146` still show `git add` then a bare `git commit`, which commits a sibling's staged files in a shared checkout.
- **D-06** (DECIDE & STATE): five run-recipe notes (`single-req-run-recipe`, `coordinator-parallel-builders-serial-integrators`, `gate-runs-need-a-quiet-machine`, `shared-tree-dirt-detached-checkout`, plus the live traps only) target one new prime, `_dev/primes/prime-work-runs.md`, instead of `CLAUDE.md`, which every session loads. The report asks this as Q2.
- **D-07** (ESCALATE): `maintainer-verify-slated-for-deletion.md` is marked **delete** on a stale premise (two months later the gate is the standard integrator check). Value: one fewer note that contradicts current practice. Risk: if the maintainer still means to delete the gate, the note's "never re-add a CLAUDE.md pointer" warning is lost. Asked as Q1 in the report.
- **D-08** (DECIDE & STATE): `run-prompts-not-in-root-docs.md` is **delete**, judged against REQ-732 (which allows `docs/run-prompts/` in this run) as well as main, where the note is still true. The reason cell says both.
- **D-09** (DECIDE & STATE): the bundle timestamp is local time (`2026-11-15_0137`), as the other reports in `ai-reports/` use it.
- **D-10** (DECIDE & STATE): the 13 "move to the repository" notes include `maintainability-audit-brief-shape.md`, targeted as a run prompt in `docs/run-prompts/` after REQ-732 (today `_dev/run-prompts/`), because it is a reusable brief, not a rule.

## Discovered Tasks
- 11 of 18 crew-member files in `skills/do-work-toolbox/crew-members/` differ from their `skills/do-work/crew-members/` originals (`cmp`): anti-slop, background-agents, clear-questions, coding-guardrails, debugging, general, maintenance, prompt-injection, security, supply-chain, testing. The gap was noted on 2026-10-23 and nothing checks it. REQ-724 edits four of these core files, so the drift grows. → report only
- `skills/do-work/actions/commit.md` commit examples stage with `git add` and then run `git commit` without a pathspec; in a shared checkout that commits another session's staged files. Already in the report as the D-05 move. → report only
- `release --manifest` without `--dry-run` or `--commit` reportedly writes the tree (memory note, 2026-09-05). Not re-run here. If still true, a bare call that mutates is a CLI design trap worth a fix. → report only

## Integration seams
- No sibling edits this bundle.
- REQ-724: two notes (`ask-tool-when-action-says-ask.md`, `subagent-model-split.md`) are "delete, already recorded at … (after REQ-724)". If REQ-724 does not merge as scoped, those two rows need a recheck.
- REQ-732: `run-prompts-not-in-root-docs.md` and `maintainability-audit-brief-shape.md` assume `docs/run-prompts/` after REQ-732.
- The catalog (`ai-reports/catalog.json`, `ai-reports/index.html`) is not regenerated; the integrator rebuilds it once.
- The report cites `file:line` locations at `df1bf423`; a sibling that deletes lines above a cited line (for example REQ-724 in `communication-style.md`, `coding-guardrails.md`) can shift them. The cited lines in REQ-724's files are named without line numbers except `coding-guardrails.md:130`.

## Proposed CHANGELOG entry
none (no shipped file changes): `/ai-reports` is export-ignore in `.gitattributes`.

## Proposed lesson
- A memory note that names a REQ "until it ships" goes stale silently: 3 of the 29 deletes (`fast-gate-caches-since-0-304`, `gate-hold-route-to-repair-req`, `dowork-followup-batch-resume`) were notes waiting on a REQ that had since shipped or been dropped.

## Lessons read
- None required (`do-work/lessons-index.md` matched none for this report-only REQ). Read for context: `skills/do-work-toolbox/actions/ai-report.md`, `ai-report-reference.md` (Proposal and Root-Cause Kinds, Report Design Rules), `completed-work-presentation-reference.md` (Safety Load Order, Collision-Safe Publication), `communication-style.md`, and the REQ-724 and REQ-732 working files.
