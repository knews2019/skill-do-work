---
id: REQ-714
title: 'Suite writes its knowledge base and ADR log under docs/'
status: completed
created_at: 2026-11-09T21:47:36Z
user_request: UR-166
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 90
  confidence: low
  calculated_at: 2026-11-09T22:06:14Z
  basis:
    - Route C
    - 60-file write set
    - 3 subsystems involved
    - 6 acceptance criteria
claimed_at: 2026-11-09T22:05:41Z
route: C
planning_at: 2026-11-09T22:17:54Z
write_set: [skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/dream_scan.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init_test.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan_test.go, skills/do-work-board/justfile.template, justfile, skills/do-work-knowledge/actions/bkb.md, skills/do-work-knowledge/actions/dream.md, skills/do-work-knowledge/actions/interview.md, skills/do-work-knowledge/actions/interview-reference.md, skills/do-work-knowledge/actions/memory-value.md, skills/do-work-knowledge/actions/memory-reference.md, skills/do-work-knowledge/actions/prompts.md, skills/do-work-knowledge/docs/bkb-guide.md, skills/do-work-knowledge/docs/dream-guide.md, skills/do-work-knowledge/docs/interview-guide.md, skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md, skills/do-work-knowledge/prompts/README.md, skills/do-work/actions/kb-lessons-handoff.md, skills/do-work/actions/review-work.md, skills/do-work/actions/work.md, skills/do-work/actions/roadmap.md, skills/do-work/actions/version.md, skills/do-work/docs/standing-preferences.md, skills/do-work/docs/verify-requests-guide.md, skills/do-work-toolbox/actions/folder-structure.md, skills/do-work-toolbox/actions/maintainability-audit.md, skills/do-work-toolbox/actions/maintainability-audit-reference.md, skills/do-work-toolbox/actions/architecture-report.md, skills/do-work-toolbox/actions/validate-feedback.md, skills/do-work-toolbox/actions/tutorial.md, skills/do-work-board/tools/queue-kanban/release.go, .gitattributes, README.md, _dev/tests/staged-skills-contract.sh, _dev/tests/audit-lockins.sh, _dev/primes/prime-releases.md, _dev/primes/prime-kanban-board.md, docs/decisions/_master_index.md, docs/decisions/_progress.md, docs/decisions/topics/_index_skill-architecture.md]
dispatch_at: 2026-11-09T22:20:14Z
builder_handback_at: 2026-11-09T22:44:01Z
integration_at: 2026-11-09T22:44:19Z
review_at: 2026-11-09T23:03:41Z
kb_status: pending
commit: 409a3e585f7ae0643f489a7886107d4ea74f0def
heavy_verified_at: 2026-11-09T23:08:30Z
heavy_verified_revision: 409a3e585f7ae0643f489a7886107d4ea74f0def
completed_at: 2026-11-09T23:08:42Z
release_at: 2026-11-09T23:08:42Z
---
# Suite writes its knowledge base and ADR log under docs/
## What
Make the do-work suite write its knowledge base (BKB) to `docs/kb/` and its ADR log to `docs/decisions/` instead of the root `kb/` and `decisions/`, and move this repository's own `kb/` and `decisions/` there. The audit found that the suite itself creates these two root folders, so every repository that uses them gets two false layout gaps.
## Detailed Requirements
- Change the BKB default from `kb/` to `docs/kb/`: the bkb-init, bkb-status and bkb-lint-structure defaults in the Go CLI and in `skills/do-work-board/justfile.template`, plus every action and prompt in `skills/do-work-knowledge` that names `kb/`.
- Change the ADR-log prompt (`skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md`) and the actions that cite `decisions/` to write `docs/decisions/`.
- Existing consumers keep working: a repository that already has `kb/` or `decisions/` and no `docs/` copy is still found (or told the one-line `--kb` override). It must never silently get a second, empty scaffold.
- In this repository: `git mv kb docs/kb` and `git mv decisions docs/decisions`; change the `.gitattributes` export-ignore lines (`/kb` to `/docs/kb`, `/decisions` to `/docs/decisions`) and the `_dev/tests/contract-regressions.sh` assertion that checks them.
- Add `docs/kb/` and `docs/decisions/` to the Target Layout tree in `skills/do-work-toolbox/actions/folder-structure.md`.
- Rewrite the live references (about 118 files, not counting `do-work/` history records).
- Layout rule this applies: all docs live under `docs/` (rule 7).
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** In a temp Git repository, run `bkb-init` with no `--kb` argument.
**Why RED now:** It creates `kb/` at the root, which `folder-structure-check` reports as `FOLDER-STRUCTURE-ROOT-ENTRY kb/`.
**GREEN when:** It creates `docs/kb/`; a repository with an existing root `kb/` still gets that one, not a second scaffold; `folder-structure-check` on this repository no longer reports `kb/` or `decisions/`; `git archive HEAD` contains neither folder; the contract suites pass.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8581 tokens, `slugged: partial`; matched because the knowledge actions and prompts change.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24067 tokens, `slugged: partial`; matched because the bkb defaults live in the Go CLI.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One resolver for the default order, three pass-through recipes, one canonical prose section, one atomic move commit (the REQ's Plan section).
- [x] **[APPLY]:** Built as planned in four commits (11afb731, f328ce6d, 4934f0bd, aaa4f4a5) plus one wording commit (00eac450). Two additions the requirement forced: bkb-init creates the parent folder (D-02) and the move hint starts with mkdir (D-01).
- [x] **[UNIFY]:** git diff -M --name-status over the range: 42 modified files, all in Scope; 520 pure renames and 24 renames with edits; 0 added, 0 deleted. Read the Go diff and the prose word-diff; no unfinished-work markers or print statements added. gofmt, go vet and shellcheck clean.
*Source: do-work capture-request: Make the do-work suite write its knowledge base and ADR log under docs/ (layout rule 7: all docs live under docs/). Change the BKB default from kb/ to docs/kb/ (the bkb-init/bkb-status/bkb-lint-structure defaults in the Go CLI and in skills/do-work-board/justfile.template, plus every action and prompt in skills/do-work-knowledge that names kb/), and change the ADR-log prompt (skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md and the actions that cite decisions/) to write docs/decisions/. Existing consumers keep working: a repo that already has kb/ or decisions/ and no docs/ copy must still be found (or told the one-line --kb override), not silently get a second empty scaffold. In this repo, git mv kb/ to docs/kb/ and decisions/ to docs/decisions/, update the .gitattributes export-ignore lines (/kb becomes /docs/kb, /decisions becomes /docs/decisions) and the contract-regressions.sh assertion that checks them, and add docs/kb/ and docs/decisions/ to the Target Layout tree in skills/do-work-toolbox/actions/folder-structure.md. Rewrite the 118 live references (not counting do-work/ history). Verify: folder-structure-check no longer reports kb/ or decisions/, a fresh bkb-init in a temp repo creates docs/kb/, git archive HEAD contains neither folder, and the contract suites pass.*

## Triage

**Route: C** - Complex

**Reasoning:** The change crosses three systems: the Go CLI defaults for the knowledge base, the prose of the knowledge and toolbox actions, and this repository's own tracked folders with their export-ignore lines, contract tests and heavy-lane path maps. About 118 live files reference the two folders, and existing consumers must keep working, which needs a designed fallback rule.

**Planning:** Required

## Plan

**Compatibility rule.** With no explicit `--kb`, the CLI uses `docs/kb` when that directory exists. Otherwise it uses a root `kb/` (then `knowledge-base/`) that already is a knowledge base, and adds one info finding with the `git mv kb docs/kb` fix. Otherwise it creates or expects `docs/kb`. An explicit `--kb` is always taken literally. When both `docs/kb` and a root `kb/` exist, `docs/kb` wins and one warning finding names both. The ADR log follows the same order in prose only (`<project-root>/docs/decisions/`, then a legacy root `decisions/` plus the hint, else create under docs), because no Go code reads that folder.

1. **RED tests (commit 1).** In `bkb_init_test.go` and `bkb_scan_test.go` of `internal/knowledgecommands`: a dry-run `bkb-init` with no `--kb` in a temp Git repository plans only paths under `docs/kb/`; `bkb-status` with only a `docs/kb` fixture reports `location=docs/kb`; a lock-in that a repository with only a root `kb/` is still found by init (`--fill-gaps --dry-run` lists only `kb/` paths, no `docs` folder appears) and by status (`location=kb` plus the hint). Reuse `runGitFixture`, `writeFixture` and `treeDigest`.
2. **GREEN Go (commit 2).** One resolver in `bkb_scan.go` replaces the three `"kb"` literals (`commands.go:69`, `interview_commands.go:52`, `memory_commands.go:181`). The parsers leave the target empty unless supplied; `locateKnowledgeBase` treats empty as the ordered default; `handleBKBInit` and `dream_scan.go` call the same resolver; `interview_commands.go:1377` compares against "explicit", not the literal.
3. **Recipes and prose (commit 3).** The three `bkb-*` recipes in `skills/do-work-board/justfile.template` and the root `justfile` become pass-through (`*args`), so they no longer hard-code the old path and defeat the fallback. `bkb.md` § Locating the Knowledge Base becomes the one canonical order and the other knowledge, core and toolbox files point at it or name the new default with the `<project-root>/` prefix. The ADR prompt resolves its root in Phase 1 with the same order. The Target Layout tree gains the two folders. The staged-skills exemption for the two new path spellings lands in this commit.
4. **Repository move (commit 4, atomic).** `git mv kb docs/kb` and `git mv decisions docs/decisions`; `.gitattributes` export-ignore lines and their comment; the `legacy_runtime_paths` exclusions and the export-ignore assertion in `_dev/tests/staged-skills-contract.sh`; `_dev/tests/audit-lockins.sh`; the link-depth rewrite in the 21 ADR records and the self-path tokens in the three index files; the two primes, `README.md`, the `standing-preferences.md` link and the `release.go` comment.
5. **Release (orchestrator, Step 9).** Changelog entry, minor bump and mirrors through the finalizer.

**Testing approach.** `go test ./internal/knowledgecommands/` after commits 1 and 2. After commit 4: `bkb-lint-structure` and `bkb-status` with no `--kb` match the baseline with `location=docs/kb`; `folder-structure-check` no longer reports the two folders; `git archive HEAD | tar -t` lists no `docs/`; the repository gate, then the heavy lanes.

**Plan validation (orchestrator):**
- Requirement coverage: every Detailed Requirement maps to a task. One requirement cannot be met as written: `_dev/tests/contract-regressions.sh` holds no export-ignore assertion (it was removed by REQ-539 and the file sits at its 77-line ceiling). Task 4 puts the assertion in `staged-skills-contract.sh` beside the existing export-ignore checks and corrects the stale sentence in `.gitattributes` that names the old home.
- No orphan tasks. The pass-through recipes are not named by the REQ but are required by it: a recipe that always passes `--kb kb` defeats the fallback.
- Scope sanity: the plan has 5 tasks (4 for the builder). Quality degrades past 3. The tasks are kept in one REQ because the suite default and this repository's move must release together; the builder commits them as four separate commits.
- Consumer field contract: `bkb-status` prints `location=<path>`; the tests assert that field and the new finding codes.
- The "about 118 live references" in the REQ counted history-like files. The plan finds 42 files outside the two folders and 24 inside the ADR folder.

*Generated by Plan agent*

## Exploration

**KB location in Go.** Three literals set the default: `commands.go:69` (`bkbOptions{target:"kb"}`, shared by bkb-init, bkb-status, bkb-lint-structure and dream-scan), `interview_commands.go:52` and `memory_commands.go:181`, all in `skills/do-work/tools/do-work-cli/internal/knowledgecommands`. `locateKnowledgeBase` (`bkb_scan.go:436-458`) adds `knowledge-base` only when the request equals the literal, then walks three ancestors. bkb-init bypasses it (`bkb_init.go:567`). Every existing test passes `--kb kb`, so none pins the default.

**Recipes.** `skills/do-work-board/justfile.template:73-82` and the root `justfile:73-82` always pass `--kb` with the old default. No test pins the bkb recipes: `install-suite-behavior.sh:256` compares the installed justfile to the template, and `replace-text-section.sh` reads recipe names only.

**Prose that states the location.** In do-work-knowledge: `bkb.md`, `dream.md`, `interview.md`, `interview-reference.md`, `memory-value.md`, `memory-reference.md`, `prompts.md`, the three guides, the ADR prompt and the prompts README. In core: `kb-lessons-handoff.md`, `review-work.md:334`, `work.md:431`, `roadmap.md`, `version.md:60`, `standing-preferences.md:12`, `verify-requests-guide.md:63`. In the toolbox: `maintainability-audit.md`, its reference, `architecture-report.md:63`, `validate-feedback.md:72,184`, `tutorial.md:333`. Relative paths inside the knowledge base (`<kb>/raw/...`) and the scaffold bytes need no change.

**Repository machinery.** `.gitattributes:44,50` carry the two export-ignore lines (comment lines 12, 17, 21). `_dev/tests/staged-skills-contract.sh:75` lists `docs` in `legacy_runtime_paths`, so any tracked file under a root `docs/` fails lines 90-102 until the two folders are excluded. `_dev/tests/audit-lockins.sh:106` globs `decisions/topics/*.md` and would skip silently after the move. `heavy-lanes.json` and `fast-stages.json` name neither folder. No Go test reads the live folders. `update-script-behavior.sh` uses a `kb/sentinel.txt` consumer fixture, still valid as the legacy case.

**Links.** No Markdown link leaves `kb/` (wikilinks only). `decisions/records/` holds 102 outbound relative links in 21 files: 41 resolve today and would break one level deeper, 61 are already dead. No Markdown link from `skills/`, `_dev/` or the root docs points into either folder. `bkb-lint-structure` checks none of these.

**Baseline.** `bkb-lint-structure --kb kb`: success, `BKB-STRUCTURE-CLEAN`. `bkb-status --kb kb`: findings (defrag and garden overdue, 68 days), `location=kb articles=245 topic_clusters=11 inbox=0 ready=0`.

**Gate risks.** A bare `docs/kb/...` in shipped Markdown fails `staged-skills-contract.sh:1041-1062` and `shipped-package-reference-contract.sh`, which read it as a same-package citation. Prose uses `<project-root>/docs/kb/`. The move and the `.gitattributes` change must be one commit, or `git archive` ships 544 maintainer files. `kb/usage-ledger.jsonl`, `decisions/audits/2026-07-15-phase1-hashes.txt` and `.gitattributes` are paths no heavy lane covers, so every heavy lane is selected.

*Generated by Explore agent (combined Plan and Explore dispatch)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/dream_scan.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init_test.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan_test.go` (modify)
- `skills/do-work-board/justfile.template` (modify)
- `justfile` (modify)
- `skills/do-work-knowledge/actions/bkb.md` (modify)
- `skills/do-work-knowledge/actions/dream.md` (modify)
- `skills/do-work-knowledge/actions/interview.md` (modify)
- `skills/do-work-knowledge/actions/interview-reference.md` (modify)
- `skills/do-work-knowledge/actions/memory-value.md` (modify)
- `skills/do-work-knowledge/actions/memory-reference.md` (modify)
- `skills/do-work-knowledge/actions/prompts.md` (modify)
- `skills/do-work-knowledge/docs/bkb-guide.md` (modify)
- `skills/do-work-knowledge/docs/dream-guide.md` (modify)
- `skills/do-work-knowledge/docs/interview-guide.md` (modify)
- `skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md` (modify)
- `skills/do-work-knowledge/prompts/README.md` (modify)
- `skills/do-work/actions/kb-lessons-handoff.md` (modify)
- `skills/do-work/actions/review-work.md` (modify)
- `skills/do-work/actions/work.md` (modify)
- `skills/do-work/actions/roadmap.md` (modify)
- `skills/do-work/actions/version.md` (modify)
- `skills/do-work/docs/standing-preferences.md` (modify)
- `skills/do-work/docs/verify-requests-guide.md` (modify)
- `skills/do-work-toolbox/actions/folder-structure.md` (modify)
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modify)
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (modify)
- `skills/do-work-toolbox/actions/architecture-report.md` (modify)
- `skills/do-work-toolbox/actions/validate-feedback.md` (modify)
- `skills/do-work-toolbox/actions/tutorial.md` (modify)
- `skills/do-work-board/tools/queue-kanban/release.go` (modify)
- `.gitattributes` (modify)
- `README.md` (modify)
- `_dev/tests/staged-skills-contract.sh` (modify)
- `_dev/tests/audit-lockins.sh` (modify)
- `_dev/primes/prime-releases.md` (modify)
- `_dev/primes/prime-kanban-board.md` (modify)
- `docs/decisions/_master_index.md` (modify)
- `docs/decisions/_progress.md` (modify)
- `docs/decisions/topics/_index_skill-architecture.md` (modify)

Plus two folder moves with no content edits beyond the link fixes named in the Plan: the 506 tracked files of the root kb folder move under docs/kb, and the 38 tracked files of the root decisions folder move under docs/decisions (the three index files above and the link-depth rewrite in the 21 files under docs/decisions/records are the only content edits).

**Files I will NOT touch:** anything under do-work/, ai-reports/ and the dated root changelog archives; content under docs/kb/raw/processed and docs/kb/wiki/sources; _dev/tests/contract-regressions.sh (at its line ceiling, nothing to update); _dev/tests/update-script-behavior.sh (its kb fixture stays as the legacy case); _dev/tests/heavy-lanes.json and fast-stages.json; folder_structure.go; the board's walk.go, board_live_test.go and verify_test.go; skills/do-work-knowledge/SKILL.md, bkb-reference.md and roadmap-guide.md. The release metadata (VERSION, both changelogs, version.md) belongs to the finalizer.

**Acceptance criteria (restated from REQ):**
- [ ] With no explicit path, bkb-init, bkb-status, bkb-lint-structure and the other knowledge commands default to docs/kb, in the Go CLI and in the justfile recipes.
- [ ] The ADR-log prompt and the actions that cite the ADR folder write and read docs/decisions.
- [ ] A repository that already has a root kb or decisions folder and no docs copy is still found and is told the one-line move; it never gets a second, empty scaffold.
- [ ] This repository's two folders are moved with git mv, the export-ignore lines and their assertion follow, and git archive HEAD contains neither folder.
- [ ] The Target Layout tree in the folder-structure action lists both folders under docs.
- [ ] Live references are rewritten; history records under do-work/ are not.
- [ ] folder-structure-check no longer reports the two folders, a fresh bkb-init in a temp repository creates docs/kb, and the contract suites pass.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/` at a79e94bb; only this REQ's working file, `baseline.json` and the run directory are changed
**Tests baseline:** ✓ Focused probe green (`bash do-work/runs/2026-11-09-layout-queue/REQ-714-probe.sh`, knowledge-base tests, 2s); repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 119s at a79e94bb, run from a detached checkout
**Dependencies:** ✓ None added

*Checked by work action*

## Decisions

- D-01 DECIDE & STATE: the move hint is `mkdir -p docs && git mv kb docs/kb`. The Plan's `git mv kb docs/kb` fails when `docs/` does not exist ("renaming 'kb' failed: No such file or directory"; probed).
- D-02 DECIDE & STATE: `bkb-init` now creates the missing parent folders of a nested target. Without it the new default failed in a Git repository. Rollback removes a parent this run created.
- D-03 DECIDE & STATE: codes are `BKB-LEGACY-LOCATION` (info) and `BKB-DUPLICATE-LOCATION` (warning). `bkb-init`, `bkb-status` and `bkb-lint-structure` emit them. `dream-scan`, `interview-ingest` and `memory-audit` use the same order but emit no hint.
- D-04 DECIDE & STATE: a notice never changes the `bkb-init` outcome. On status and lint the warning gives outcome `findings`; the info hint keeps `success`.
- D-05 DECIDE & STATE (superseded by D-15): "both exist" means the `docs/kb` directory exists and a root `kb/` or `knowledge-base/` holds `raw/` and `wiki/`. If `docs/kb` is an empty folder beside a real root `kb/`, `docs/kb` still wins (the Plan's literal rule) and the warning names both.
- D-06 DECIDE & STATE: an explicit `--kb kb` no longer falls back to `knowledge-base/`. Before, the code could not tell an explicit `kb` from the default.
- D-07 DECIDE & STATE: `bkb.md` tells the agent to pass `--kb` only when the user gave a path. The command then applies the order and the hint reaches the user.
- D-08 DECIDE & STATE: the staged-skills exemption is the prefixes `docs/kb/` and `docs/decisions/` plus the exact token `docs/decisions`. The contract cannot see how a token is spelled in context. Bare spellings remain only in literal command arguments: six `--exclude-path docs/kb/`, the `verify-requests --against` example and the fenced `git mv` in the ADR prompt. `shipped-package-reference-contract.sh` needed no change.
- D-09 DECIDE & STATE: the ADR-014 pointer in `standing-preferences.md` is now a canonical GitHub link. A bare `docs/decisions/...` path fails the shipped reference contract. The link target exists on GitHub only after this work reaches `main`.
- D-10 DECIDE & STATE: `audit-lockins.sh` fails when the topic-index glob matches nothing. This is the silent skip the Exploration named.
- D-11 DECIDE & STATE: the builder ran `git add -A docs` once, limited to the `docs` path, to stage the 24 edited files inside the moved folders. This breaks the letter of "never `git add -A`". It checked the result before committing: 544 renames, 0 added files.
- D-12 (DECIDE & STATE, orchestrator): `skills/do-work/actions/version.md` line 60 named the old folder. The builder left it because the brief reserved that file for the finalizer; the ban was about the version line, not this sentence. The orchestrator changed it on the branch (commit 00eac450) to say the update does not mutate "the knowledge base", so the sentence is right for both locations.
- D-13 (DECIDE & STATE, orchestrator, plan validation): the export-ignore assertion lives in `_dev/tests/staged-skills-contract.sh`, because `_dev/tests/contract-regressions.sh` no longer holds one (removed by REQ-539) and sits at its line ceiling. The stale sentence in `.gitattributes` that named the old home is corrected.
- D-14 (DECIDE & STATE, orchestrator): the relative links in the dated ADR records were rewritten one level deeper, although `audit-lockins.sh` says records cite the tree as it was. The rewrite changes only path depth, and without it 41 links that resolve today would break.
- D-15 (DECIDE & STATE, review F1 and F4): with no explicit path, `docs/kb` wins only when it is a knowledge base or when no legacy root folder exists. A partial root folder (only `raw/` or only `wiki/`) counts as the legacy location for `bkb-init`, which fills it in place as the old CLI did. When a non-knowledge-base `docs/kb` is in the way, the hint says to remove it before the move.
- D-16 (DECIDE & STATE, review F3): the six `audit-metrics` command blocks in `maintainability-audit.md` exclude both `docs/kb/` and `kb/`, because `--exclude-path` repeats and a prefix that matches nothing is harmless.

## Discovered Tasks

- `CHANGELOG-2026-09-05-up-to-v0.303.6.md:2221` links to `blob/main/decisions/records/adr-020-...`; the link breaks on GitHub once the move reaches `main`. The file is in the "will NOT touch" list. impact-negligible → report only
- 61 relative links in `docs/decisions/records/*.md` were dead before the move and are still dead. impact-negligible → report only
- Files inside `docs/kb/` and the dated ADR records still name `kb/` and `decisions/` in their text (no content edits were in scope). impact-negligible → report only
- A knowledge base folder whose name starts with `-` cannot be passed to the pass-through recipes as a bare word; `just bkb-status --kb <name>` works. impact-negligible → report only
- `interview-ingest` and `memory-audit` now read `--kb ""` as "not supplied". impact-negligible → report only

## Implementation Summary

**Files changed:**
- `.gitattributes` (modified)
- `README.md` (modified)
- `_dev/primes/prime-kanban-board.md` (modified)
- `_dev/primes/prime-releases.md` (modified)
- `_dev/tests/audit-lockins.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `justfile` (modified)
- `skills/do-work-board/justfile.template` (modified)
- `skills/do-work-board/tools/queue-kanban/release.go` (modified)
- `skills/do-work-knowledge/actions/bkb.md` (modified)
- `skills/do-work-knowledge/actions/dream.md` (modified)
- `skills/do-work-knowledge/actions/interview-reference.md` (modified)
- `skills/do-work-knowledge/actions/interview.md` (modified)
- `skills/do-work-knowledge/actions/memory-reference.md` (modified)
- `skills/do-work-knowledge/actions/memory-value.md` (modified)
- `skills/do-work-knowledge/actions/prompts.md` (modified)
- `skills/do-work-knowledge/docs/bkb-guide.md` (modified)
- `skills/do-work-knowledge/docs/dream-guide.md` (modified)
- `skills/do-work-knowledge/docs/interview-guide.md` (modified)
- `skills/do-work-knowledge/prompts/README.md` (modified)
- `skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md` (modified)
- `skills/do-work-toolbox/actions/architecture-report.md` (modified)
- `skills/do-work-toolbox/actions/folder-structure.md` (modified)
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (modified)
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modified)
- `skills/do-work-toolbox/actions/tutorial.md` (modified)
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified)
- `skills/do-work/actions/kb-lessons-handoff.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/roadmap.md` (modified)
- `skills/do-work/actions/version.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/docs/standing-preferences.md` (modified)
- `skills/do-work/docs/verify-requests-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/dream_scan.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go` (modified)
- `docs/decisions/_master_index.md` (modified)
- `docs/decisions/_progress.md` (modified)
- `docs/decisions/topics/_index_skill-architecture.md` (modified)

Folder moves: the 506 tracked files of the root kb folder now live under docs/kb and the 38 tracked files of the root decisions folder under docs/decisions (544 renames; 520 are byte-identical, 24 carry edits: the three index files listed above and the link-depth rewrite in the 21 files under docs/decisions/records).

**What was done:** The knowledge commands of the CLI (`bkb-init`, `bkb-status`, `bkb-lint-structure`, `dream-scan`, the interview and memory commands) share one resolver: with no explicit `--kb` they use `docs/kb`, then an existing root `kb` or `knowledge-base` folder with an info finding that gives the move command, else they create or expect `docs/kb`. When both hold a knowledge base, `docs/kb` wins and a warning names both. `bkb-init` creates the missing parent folder. The three `bkb-*` justfile recipes pass their arguments through, so they no longer force the old path. The knowledge, core and toolbox prose names the new default, with `bkb.md` holding the one canonical order, and the ADR-log prompt resolves its folder the same way (docs first, a legacy root folder with a move hint, else create under docs). This repository's own knowledge base and ADR log moved under `docs/` in one commit together with the export-ignore lines, the staged-skills exclusions and assertion, and the audit lock-in glob. The Target Layout tree lists both folders. After the first review, the resolver was tightened so a `docs/kb` folder that is not a knowledge base never hides a root one, the ADR prompt got the same rule (the docs folder is the log only when it holds the index file), the maintainability-audit commands exclude both locations, and the README names both.

## Qualification

**Range:** `3ec1215b..409a3e58` (two merges of the REQ branch: RED 11afb731, GREEN f328ce6d, recipes and prose 4934f0bd, atomic move aaa4f4a5, wording 00eac450; then the review fixes 06841482, 828144e7, a61ef3ef; first merge fae09f99)
**Mechanical gates (on the first merge):** qualify ✓ satisfied, scope-drift ✓ satisfied (45 declared paths, 45 touched; the 544 renames are described in prose in both sections, as planned).

**Requirement trace (orchestrator, read against the diff):**
- Default location: `resolveKnowledgeBaseLocation` in `bkb_scan.go` is the one resolver. An explicit path is literal; with none it returns `docs/kb` when that folder is a knowledge base (holds `raw/` and `wiki/`) or when no legacy root folder exists, else a root `kb` or `knowledge-base`, else `docs/kb`. The first merge let any existing `docs/kb` directory win; review finding F1 showed that an empty one then hid a real root `kb/`, and the second merge fixed it. The three old literals are gone, and `bkb-init`, status, lint, dream-scan, interview and memory commands call the resolver.
- Existing consumers: a legacy root folder is used and reported with `BKB-LEGACY-LOCATION` (info, with the move command `mkdir -p docs && git mv kb docs/kb`); two real knowledge bases give `BKB-DUPLICATE-LOCATION` (warning). A root folder that holds only `raw/` or only `wiki/` is filled in place by a bare `bkb-init` (review F4). The lock-in test pins that a root-only repository gets no `docs` folder from `bkb-init --fill-gaps`.
- Recipes: the three `bkb-*` recipes in the template and the root justfile pass their arguments through; none names the old path.
- ADR log: the prompt resolves its folder in its first phase with the same order, and the README row, `prompts.md` and the toolbox actions name the new default with the `<project-root>/` prefix.
- This repository: 544 renames in one commit with `.gitattributes`, the staged-skills exclusions and assertion, and `audit-lockins.sh`; `git archive HEAD` lists no `docs/` path (builder probe, repeated by the reviewer).
- Target Layout: `folder-structure.md` lists `kb/` and `decisions/` under `docs/`.
- History under `do-work/` and `ai-reports/` is untouched (0 files in the range).

**Live flow:** the built CLI at the merge, with no `--kb`, reports `location=docs/kb articles=245 topic_clusters=11` on this repository, the same counts as the baseline at the old path.

**Judgment on the qualification warnings:** none reported.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/2026-11-09-layout-queue/REQ-714-probe.sh` (the knowledge-base location tests in `internal/knowledgecommands`), then `bash _dev/tests/maintainer-verify.sh` from a detached checkout at each merge revision. The builder also ran the gate once on its branch before the hand-back (exit 0, 157s).
**Result:** ✓ Focused probe green. Repository gate exit 0 at the first merge fae09f99 (gate wall 142s) and exit 0 at the final merge 409a3e58 (gate wall 146s). Both Go stages reported EXECUTING, not REUSED. On the branch: `shipped-package-reference-contract.sh` PASS, `contract-regressions.sh` exit 0, `staged-skills-contract.sh` (heavy tier) PASS, `audit-lockins.sh` passed.

**Red-green validation:**
- Captured RED case (a bare `bkb-init` in a temp Git repository creates a root `kb/`): `TestBKBInitWithoutTargetCreatesDocsKnowledgeBase` and the `docs/kb` cases of `TestBKBScanDefaultLocationOrder` ✗ at commit 11afb731 (planned paths under `kb/`; status `BKB-NOT-FOUND`) → ✓ at commit f328ce6d.
- Lock-in for existing consumers (`TestBKBInitWithoutTargetKeepsLegacyRootKnowledgeBase` and the root-only table case): a repository with only a root `kb/` is found by status (`location=kb` plus `BKB-LEGACY-LOCATION`) and `bkb-init --fill-gaps --dry-run` plans no path under `docs/`. Green before and after, by design; it pins the second-empty-scaffold failure.
- Review F1 and F4: `TestBKBInitWithoutTargetIgnoresDocsFolderThatIsNotAKnowledgeBase`, `TestBKBInitWithoutTargetFillsPartialLegacyRootFolderInPlace` and the table case "docs folder that is not a knowledge base beside legacy kb" ✗ at commit 06841482 (init planned `docs/kb/raw/_inbox_queue.md`; status `BKB-NOT-FOUND`) → ✓ at commit 828144e7.
- Captured GREEN on this repository: with no `--kb`, the built CLI at 409a3e58 reports `location=docs/kb articles=245 topic_clusters=11`; `folder-structure-check` no longer lists `kb/` or `decisions/`; `git archive` of the merge lists no `docs/` path.

**New tests added:**
- `TestBKBInitWithoutTargetCreatesDocsKnowledgeBase`, `TestBKBInitWithoutTargetKeepsLegacyRootKnowledgeBase`, `TestBKBInitWithoutTargetIgnoresDocsFolderThatIsNotAKnowledgeBase`, `TestBKBInitWithoutTargetFillsPartialLegacyRootFolderInPlace` and `TestBKBScanDefaultLocationOrder` (six cases, each also run through `memory-audit`)
- The export-ignore assertion and its mutation check in `_dev/tests/staged-skills-contract.sh`; the non-empty-glob check in `_dev/tests/audit-lockins.sh`

**Heavy verification plan:**
- Range: 3ec1215b234454d74284d155a43acf842551afe4..409a3e585f7ae0643f489a7886107d4ea74f0def
- All six lanes are selected because coverage is uncertain for `.gitattributes`, the root `justfile`, and the moved `usage-ledger.jsonl` and `phase1-hashes.txt` files.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

*Verified by work action*

## Review

**Overall: 95%** | 2026-11-09T23:03:41Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Reviewed range: `3ec1215b..409a3e58` (first review covered `3ec1215b..fae09f99` and scored 89%; the re-check covered `fae09f99..409a3e58`, 10 files, all inside the declared Scope). Decision D-05 (a `docs/kb` folder wins as soon as it exists) is superseded by the F1 fix and is recorded as such in the hand-back addendum.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 A `docs/kb` folder that is not a knowledge base (empty, or holding other files) hid a real root `kb/`: status reported not found, and a bare `bkb-init` wrote a second scaffold. FIXED in the re-check (commits 06841482 RED, 828144e7 GREEN): `resolveKnowledgeBaseLocation` in `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go` now lets `docs/kb` win only when it holds `raw/` and `wiki/`, and `BKB-DUPLICATE-LOCATION` is reported only when both folders are knowledge bases. Reproduction rerun with the built CLI at 409a3e58, once with an empty `docs/kb` and once with `docs/kb/faq.md`, beside a real root `kb/`: `bkb-status` reports `location=kb`, `bkb-lint-structure` is `success`, `memory-audit --engine bkb` reports `Active` at `kb`, `dream-scan` targets `kb`, a bare `bkb-init` is `refused` with `BKB-INIT-EXISTS`, `--fill-gaps` works in `kb/`, and nothing is written under `docs/` in either case. The hint for this case reads "docs/kb/ exists and is not a knowledge base. Remove that folder first (move its files elsewhere), then: mkdir -p docs && git mv kb docs/kb"; it is correct, because the plain move command would nest the knowledge base inside the existing folder — impact-user-visible → report only

**Minor findings:**
- F2 The ADR prompt chose `docs/decisions` as soon as the directory existed. FIXED: `skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md:47` now asks for `_master_index.md`, line 48 keeps a root log in use beside a `docs/decisions` folder without the index and tells the agent to announce a manual merge instead of the move command, and line 53 says existing files in the folder stay as they are. An agent with a root `decisions/` log continues there in every case — impact-user-visible → report only
- F3 Toolbox actions named only the new location. FIXED: the six command blocks in `skills/do-work-toolbox/actions/maintainability-audit.md` pass `--exclude-path docs/kb/ --exclude-path kb/` (the tool matches a repository-relative prefix, so `kb/` covers only the root folder and is harmless when absent), `maintainability-audit-reference.md:50` says so, `architecture-report.md:63`, `validate-feedback.md:72` and `maintainability-audit.md:64` name a root `decisions/` as the second place, and the two red-flag lines no longer carry a path — impact-user-visible → report only
- F4 A root `kb/` holding only `raw/` or only `wiki/` got a second folder under `docs/`. FIXED: the resolver accepts such a folder when no complete legacy folder exists. Reproduced both forms: a bare `bkb-init` fills `kb/` in place, keeps the existing file bytes, creates no `docs` folder and carries the move hint. A complete `knowledge-base/` is preferred over a partial `kb/`. A root `kb/` with neither subfolder is ignored and `docs/kb` is created — impact-negligible → report only
- F5 Test gaps. PARTLY FIXED: the F1 case is in `TestBKBScanDefaultLocationOrder`, two new init tests pin F1 and F4, and every table case also runs `memory-audit --engine bkb`, so a second command is pinned on both the `docs/kb` and the legacy branch. STILL OPEN: no default-location test for `interview-ingest` (the command that writes into the knowledge base) or `dream-scan`, and no test that a failed init removes a `docs` parent it created — impact-negligible → report only
- F6 `README.md` named only `docs/kb/`. FIXED: lines 33, 37 and 49 name both `docs/kb/` and a root `kb/` — impact-negligible → report only
- F7 (nit) Commit subject and names. FIXED for the subject (`docs(adr): scaffold the decisions log`, line 131). STILL OPEN by instruction: the single-word names `notices`, `relative` and `legacy` in `bkb_scan.go`; the one new field, `defaultOccupied`, has two words — impact-negligible → report only
- F8 (nit, new) The hint for a `docs/kb` folder that is in the way gives only one way forward (remove that folder, then move). It does not say that staying on the root `kb/` is also fine, and the info finding repeats on every status, lint and init run while both folders exist — impact-negligible → report only

**Acceptance:** Pass — Implementation and Integration stages on the merged tree at 409a3e58: `go test -count=1 ./internal/knowledgecommands/` ok, `gofmt -l` and `go vet` clean, `shipped-package-reference-contract.sh` PASS, `git archive HEAD | tar -t | grep -c '^docs/'` is 0. Built CLI in temp Git repositories: the F1 and F4 reproductions above; no regression in the first-pass scenarios (fresh repository creates `docs/kb`; root `kb/` only and `knowledge-base/` only are found with `BKB-LEGACY-LOCATION`, bare init refused, and the printed move command works; both real gives `docs/kb` plus the warning on status, lint and init; `--kb kb` and `--kb nowhere/x` are literal; `--kb ""` is a usage error). Extra edges tried: `docs/kb` as a plain file beside a real `kb/` (root folder used, hint shown), a partial `docs/kb` alone (filled in place), a complete `docs/kb` beside a partial root `kb/` (`docs/kb` used, no notice). The `just` recipes were exercised in the first pass and did not change. The staged-skills contract and the repository gate were not run here (orchestrator). Deployment and Live acceptance: not exercised.
**Restatement sweep:** redefined the default knowledge base location (root `kb/` to `docs/kb/`, with the legacy order: `docs/kb` when it is a knowledge base, then a root `kb/` or `knowledge-base/`, also when only partly built, else `docs/kb`), the default ADR log location (root `decisions/` to `docs/decisions/`, keyed on `_master_index.md`), the meaning of an explicit `--kb kb` (no `knowledge-base/` fallback any more), the three `bkb-*` recipe signatures (pass-through arguments), the stated parent search depth (root plus 2 levels), the wording of the two location notices, and the home of the export-ignore assertion; swept `skills/`, `README.md`, `_dev/`, `CLAUDE.md`, `suite/`, `tools/`. `skills/do-work-knowledge/actions/bkb.md` § Locating the Knowledge Base (steps 2, 3, 5 and the closing sentence) agrees with the Go code after the fix. No other file quotes the old notice texts or the old "when that directory exists" rule. `kb-lessons-handoff.md:31`, `memory-value.md` and `interview.md` point at the canonical section and stay correct. The only stale link left is `CHANGELOG-2026-09-05-up-to-v0.303.6.md:2221` (already a Discovered Task).
**Suggested testing:** 5 items — (1) Integration: take the staged-skills contract, the repository gate and the heavy lanes from the orchestrator's run at 409a3e58; (2) Deployment: unassessed. Update an installed consumer that has a root `kb/` and the old recipes, then confirm the Just section holds the pass-through recipes and `just bkb-status` reports `location=kb` with the move hint; (3) Live acceptance: unassessed. After the change reaches `main` on GitHub, open the ADR-014 link in `skills/do-work/docs/standing-preferences.md:12`; (4) Agent run of the ADR prompt in a repository with a root `decisions/` log and a `docs/decisions` folder of hand-written ADRs: confirm it reports the root folder, UPDATE mode and the manual-merge note before writing; (5) Edge: a repository whose `docs` is a symbolic link to a folder outside the repository. A bare `bkb-init` builds a standalone knowledge base with its own `.git` at the link target (reproduced in the first pass; the old CLI did the same for an explicit `--kb docs/kb`).
**Follow-ups created:** None (8 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: 3ec1215b234454d74284d155a43acf842551afe4
- Target revision: 409a3e585f7ae0643f489a7886107d4ea74f0def
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for `.gitattributes`, the root `justfile` and the moved ledger and hash files
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same uncertain coverage
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — same uncertain coverage
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — same uncertain coverage
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — same uncertain coverage
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — same uncertain coverage

## Heavy Verification Result

- Target revision: 409a3e585f7ae0643f489a7886107d4ea74f0def
- Execution revision: 409a3e585f7ae0643f489a7886107d4ea74f0def (detached checkout)
- queue-kanban-javascript: exit 0, 7s, executed
- queue-kanban-browser: exit 0, 96s, executed (real browser, not skipped)
- do-work-cli-integrations: exit 0, 87s, executed
- staged-skills: exit 0, 39s, executed
- updater: exit 0, 70s, executed
- installer: exit 0, 33s, executed

## Lessons Learned

**What worked:** One Plan-and-Explore dispatch before the build. It found that the justfile recipes always passed `--kb kb`, which would have defeated any fallback in the Go code, and that a root `docs` folder is a "legacy runtime path" in the staged-skills contract. Four separate commits (RED, GREEN, prose, atomic move) kept the 544 renames out of every code diff. The reviewer's scenario runs in temp repositories found the one real defect.
**What didn't:** The first resolver let `docs/kb` win as soon as the directory existed. An empty or unrelated `docs/kb` then hid a real root `kb/`, and a bare `bkb-init` wrote a second scaffold, which is the exact failure the REQ forbids. The rule had to test what the folder is (it holds `raw/` and `wiki/`), not that it exists. The REQ also named an assertion in `contract-regressions.sh` that REQ-539 had removed; the captured "118 live references" counted history files, and the real number was 42.
**Worth knowing:** The location order is stated once, in `bkb.md` § Locating the Knowledge Base, and implemented once, in `resolveKnowledgeBaseLocation`. Change both together. Shipped prose must write `<project-root>/docs/kb/`, because the reference contracts read a bare `docs/...` as a path inside the skill package; the staged-skills contract carries an exemption only for literal command arguments. `git mv kb docs/kb` fails when `docs/` does not exist, so every hint starts with `mkdir -p docs`. The move and the `.gitattributes` lines must stay in one commit, or `git archive` ships the maintainer's knowledge base.

## Orientation

Now the suite keeps a project's knowledge base in `docs/kb/` and its decision log in `docs/decisions/`, and every knowledge command finds an older root `kb/` or `decisions/` and tells the user the one-line move. This lives in the knowledge commands of the do-work CLI (one shared resolver), the do-work-knowledge actions and ADR prompt, and the managed justfile recipes. `[MAP CHANGED]` The default data location of the knowledge package moved, and this repository's own `kb/` and `decisions/` now sit under `docs/`, still excluded from the distributed archive.

Why it matters: an installed repository's root now holds only `do-work/` and `ai-reports/` from the suite, which is what the folder-structure layout expects.

Prime staleness check: `prime-releases.md` and `prime-kanban-board.md` were updated in this change to name `docs/decisions/`; `prime-action-files.md` names neither folder.

## Timing

Observed 2026-11-09T22:20:14Z to 2026-11-09T23:08:29Z: 48m 15s total, 32m 22s attributed across 3 events, 15m 53s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 23m 47s | 1 |
| verification-gate | 8m 35s | 2 |

Slowest stage: builder-work / builder agent, four commits, 23m 47s, outcome success.
