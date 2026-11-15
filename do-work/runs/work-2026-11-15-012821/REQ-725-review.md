# REQ-725 Review — sort the local memory notes by destination

Reviewer, orchestrated mode, Route A. Range `e5972dc5..27e77d53` on `main` (tip `27e77d53`). The Restatement Sweep was not run: this REQ is not the last successful integration of the run.

## Verdict first

The report does what the REQ asks. It gives every one of the 46 notes in the memory folder exactly one row, with a destination, a reason and a target or an "already recorded at" path. The counts are right, the four rows that depended on REQ-724 (moving the global standing rules into the skill, release 0.322.0) and REQ-732 (allowing `docs/run-prompts/`) match main, and the change set touches no shipped file. I found two wrong citations in 16 sampled rows. Neither changes a destination. No fix is required before approval.

## Checks

### 1. Sampled rows (16 rows: 3 move to the skill, 5 move to the repository, 8 delete)

For each row I read the memory note, then checked every named path and cited `file:line` on main at `27e77d53`.

| # | Note | Destination | What I checked on main | Verdict |
|---|------|-------------|------------------------|---------|
| S1 | `one-orchestrator-per-queue.md` | move to the skill | `skills/do-work/actions/commit.md:125` and `:144` are `git add …`, `:127` and `:146` are a bare `git commit -m`. The note's lasting fix is the path-limited commit. Moving the rule only, not the incident log, is fair. | Correct |
| S2 | `warnings-over-exemptions.md` | move to the skill | No hit for `carve-out` or `exemption` in `skills/do-work/crew-members/`. The note is a general rule for health checks. Target `coding-guardrails.md` exists. | Correct |
| S3 | `make-the-proposal-means-recommend-only.md` | move to the skill | No hit for `make the proposal` under `skills/`. Target `skills/do-work-toolbox/actions/ai-report-reference.md` exists. | Correct |
| R1 | `bare-cli-subcommand-executes.md` | move to the repository | `skills/do-work/tools/do-work-cli/internal/publication/release.go:50` returns `RELEASE-PREIMAGE-STALE`, as the row says. Target `_dev/primes/lessons-releases.md` is a repository file. | Correct |
| R2 | `browser-heavy-lane-skips-silently.md` | move to the repository | `_dev/tests/maintainer-verify.sh:42` holds the skip line naming `QUEUE_KANBAN_BROWSER`. The Chrome path in the note is this machine's fact, so a repository prime fits. | Correct |
| R3 | `maintainability-audit-brief-shape.md` (REQ-732 seam) | move to the repository | `docs/run-prompts/` exists and `_dev/run-prompts/` is gone. `maintainability-audit.md` and its reference have no §Score category-reuse rule. `do-work/audits/audit-2026-09-03.md` exists. A reusable brief as a run prompt is fair. | Correct |
| R4 | `record-lessons-after-fixes.md` | move to the repository | `skills/do-work/actions/work.md:80` is the Required-lessons **consult** (reading lessons before exploration). The phase that **writes** a lesson is Step 7.5 Lessons-Capture Phase at `work.md:108`. The destination is still fair. | Wrong line (F1) |
| R5 | `finalization-manifest-recipe.md` | move to the repository | `do-work/runs/2026-11-09-layout-queue/helpers/make-manifests.py` exists. The note's `maintainer_release: true` field is specific to this repository. | Correct |
| D1 | `ask-tool-when-action-says-ask.md` (REQ-724 seam) | delete | `skills/do-work/crew-members/communication-style.md:25` says "A decision for the user goes through the ask-user tool", with no autonomy exception. `clear-questions.md` principle 5 (line 29) has "recommended option first … value and its risk". Both landed in 0.322.0 (`7c15a2dd`). | Correct |
| D2 | `subagent-model-split.md` (REQ-724 seam) | delete | "Pick each agent's model by its work" is at `skills/do-work/crew-members/background-agents.md:190`, and the toolbox and knowledge copies carry it at line 191. The rest of the note is one machine's launch command. | Correct |
| D3 | `fast-gate-caches-since-0-304.md` | delete | `do-work/archive/UR-127/REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages.md` exists. `_dev/tests/fast-stages.json:10` covers `do-work/archive/UR-003/input.md`, and the board stage covers the `do-work` subtree. The note's false green is closed. | Correct |
| D4 | `velocity-lessons-report.md` | delete | The cited folder `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/` does not exist. Commit `ffc2728d` (Organize ai-reports by month, 2026-11-04, before the builder's base) moved it to `ai-reports/2026-09_18-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`. The catalog lists it there, so the destination is still fair. | Wrong path (F2) |
| D5 | `capture-files-containment-recipe.md` | delete | `skills/do-work/actions/capture-reference.md:196` states the fence rule, "no info string", and that `capture-files` refuses a UR without the block. | Correct |
| D6 | `takeover-resets-claimed-reqs.md` | delete | `skills/do-work/actions/work-reference.md:316` says plain recovery preserves claims and points at `advance REQ-NNN`. | Correct |
| D7 | `gate-hold-route-to-repair-req.md` | delete | `work-reference.md:115` defines the `repository_gate_repair` marker, and `do-work/archive/UR-095/` exists. The note's "until REQ-491/492 ship" is stale. | Correct |
| D8 | `one-broken-pipe-principle.md`, `compound-greppable-trigger-words.md`, `descriptive-changelog-titles.md`, `stack-preference-go-sqlite.md` | delete | `CLAUDE.md:14` and `work-reference.md:29` (broken pipe). `coding-guardrails.md:132` and `CLAUDE.md:53` (naming). `_dev/primes/prime-releases.md:7` (changelog title). `skills/do-work-toolbox/actions/folder-structure.md:45` (Go with SQLite). Each line says what the row claims. | Correct |

I also confirmed `do-work/HANDDOWN-UR-031.md` is the only handoff file, as the `dowork-followup-batch-resume.md` row says.

### 2. Counts

- The live memory folder holds 47 files: `MEMORY.md` (the index) and 46 notes.
- The page has 46 `data-memory-note` rows, and the file names match the folder exactly (`diff` of the two sorted lists is empty).
- Destinations: 5 `move-to-skill`, 13 `move-to-repository`, 28 `delete`. The page's tiles and the line "5 + 13 + 28 = 46" agree.
- `run-prompts-not-in-root-docs.md` is absent from the folder, and the page says it was deleted on 2026-11-15 when REQ-732 landed. `MEMORY.md` no longer lists it.
- No page text still says 47 notes or 29 deletes as the current state. The one "47 notes" is in the evidence ledger, which says what was read before the deletion and then says the folder now holds 46.

### 3. Seam rows re-checked by the integrator

- `ask-tool-when-action-says-ask.md` and `subagent-model-split.md`: correct on main (D1, D2 above).
- `maintainability-audit-brief-shape.md`: correct on main (R3 above).
- `run-prompts-not-in-root-docs.md`: the row is removed, and the page explains why (check 2). The ledger cites `do-work/archive/REQ-732-allow-docs-run-prompts-folder.md`, which exists.
- The ledger cites `do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md`. Today the file is at `do-work/archive/REQ-724-move-global-standing-rules-into-skill.md`. D-12 says the UR-175 finalization creates the cited path. That is plausible, because UR-175 holds only REQ-724 and REQ-725, but it is not true on main yet (F3).

### 4. Probe

`bash do-work/runs/work-2026-11-15-012821/REQ-725-probe.sh` from ROOT:

```
REQ-725 probe PASS: 46 notes, one row each, in ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/index.html
```

Exit 0.

### 5. Shipped files and catalog

- `git diff --name-only e5972dc5..27e77d53` lists five files, all under `ai-reports/`. Nothing under `skills/`, `suite/`, `_dev/`, `do-work/` or `CLAUDE.md` changed.
- `presentation.css` and `presentation.js` are byte copies of the kit files (`cmp` exit 0).
- `ai-reports/catalog.json` lists the new bundle as kind `proposal` with `linked_ids: ["REQ-725"]`, and `ai-reports/index.html` has its row under "proposal".

### 6. Anti-bloat

Zero items the REQ did not name. The two kit files are part of the normal `ai-report` bundle. The options, open questions, capture lines and "another layout" section are the proposal kind's standard parts. No helper, test or extra file was added.

## Review

**Overall: 92%** | 2026-11-15T03:02:45Z

| Dimension | Score |
|-----------|-------|
| Requirements | 96% |
| Code Quality | 88% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Code Quality here means citation accuracy: 2 wrong citations in 16 sampled rows.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1: the `record-lessons-after-fixes.md` row cites `skills/do-work/actions/work.md:80` as the "Lessons phase", but line 80 is the Required-lessons consult (reading). The lesson-writing phase is Step 7.5 at `work.md:108`. The destination is unchanged. impact-negligible → report only
- F2: the `velocity-lessons-report.md` row cites `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`, which does not exist. The folder moved to `ai-reports/2026-09_18-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/` in `ffc2728d` (2026-11-04). The destination is unchanged. impact-negligible → report only
- F3: the evidence ledger cites `do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md`, which exists only after the UR-175 finalization. Today the file is `do-work/archive/REQ-724-move-global-standing-rules-into-skill.md`. If finalization does not create the UR-175 folder, the citation is wrong. impact-negligible → report only

**Sampled rows:** S1-S3 (move to the skill), R1-R5 (move to the repository), D1-D8 (delete; D8 covers four notes). 14 correct, R4 has a wrong line (F1), D4 has a wrong path (F2). Every sampled destination is a fair judgment of its note.
**Anti-bloat:** 0 unrequested helpers, options, files or tests.
**Acceptance:** Pass — 46 rows for 46 notes, counts 5 + 13 + 28, probe PASS, no shipped file in the range, catalog lists the bundle.
**Restatement sweep:** not run (not the last successful integration of the run).
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
