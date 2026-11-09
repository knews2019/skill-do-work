---
id: REQ-714
title: 'Suite writes its knowledge base and ADR log under docs/'
status: claimed
created_at: 2026-11-09T21:47:36Z
user_request: UR-166
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-11-09T22:05:41Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: Make the do-work suite write its knowledge base and ADR log under docs/ (layout rule 7: all docs live under docs/). Change the BKB default from kb/ to docs/kb/ (the bkb-init/bkb-status/bkb-lint-structure defaults in the Go CLI and in skills/do-work-board/justfile.template, plus every action and prompt in skills/do-work-knowledge that names kb/), and change the ADR-log prompt (skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md and the actions that cite decisions/) to write docs/decisions/. Existing consumers keep working: a repo that already has kb/ or decisions/ and no docs/ copy must still be found (or told the one-line --kb override), not silently get a second empty scaffold. In this repo, git mv kb/ to docs/kb/ and decisions/ to docs/decisions/, update the .gitattributes export-ignore lines (/kb becomes /docs/kb, /decisions becomes /docs/decisions) and the contract-regressions.sh assertion that checks them, and add docs/kb/ and docs/decisions/ to the Target Layout tree in skills/do-work-toolbox/actions/folder-structure.md. Rewrite the 118 live references (not counting do-work/ history). Verify: folder-structure-check no longer reports kb/ or decisions/, a fresh bkb-init in a temp repo creates docs/kb/, git archive HEAD contains neither folder, and the contract suites pass.*
