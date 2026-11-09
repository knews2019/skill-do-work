---
id: REQ-717
title: 'Suite-managed justfile recipes call a script in one line'
status: claimed
created_at: 2026-11-09T21:47:36Z
user_request: UR-169
domain: general
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-shell-commands.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-713]
claimed_at: 2026-11-09T23:39:43Z
---
# Suite-managed justfile recipes call a script in one line
## What
Move the logic out of the suite-managed justfile recipes in `skills/do-work-board/justfile.template` so each managed recipe is one line that calls a script shipped with do-work-board. Today `run-kanban` has 3 body lines and `kanban-static` has 2, so every consumer repository gets this layout gap.
## Detailed Requirements
- `run-kanban`: move the port validation, the stale queue-kanban listener stop, and the build-and-serve into a shipped script (for example under `skills/do-work-board/tools/`). Keep the exact behavior and messages.
- `kanban-static`: move the build-and-generate and the `.git/info/exclude` write into a shipped script. Keep the exact behavior and messages.
- The installer rewrites the managed block in consumer justfiles, so also update this repository's root justfile.
- Depends on REQ-713 (folder-structure-check sees recipes with default parameter values); without it, the check cannot prove `run-kanban` is fixed.
- Layout rule this applies: a justfile recipe holds no logic and calls a script in one line (rule 5).
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository after REQ-713 lands.
**Why RED now:** It reports `FOLDER-STRUCTURE-JUST-LOGIC` for `justfile: kanban-static` and `justfile: run-kanban`.
**GREEN when:** No `FOLDER-STRUCTURE-JUST-LOGIC` gap is reported, `just run-kanban` still replaces a stale instance on its port, `just kanban-static` still writes the `build/` snapshot and the exclude line, and the board contract tests pass.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` — 6122 tokens, `slugged: partial`; matched because the run-kanban and kanban-static recipes drive the board.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, `slugged: partial`; matched because the recipe logic becomes shipped shell.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: Move the logic out of the suite-managed justfile recipes so each managed recipe is one line (layout rule 5: a justfile recipe holds no logic and calls a script in one line). In skills/do-work-board/justfile.template, run-kanban has 3 body lines (port validation, stop a stale queue-kanban listener, build and serve) and kanban-static has 2 (build and generate, then the .git/info/exclude write). Move that logic into scripts shipped with do-work-board (for example under skills/do-work-board/tools/), keep the exact behavior and messages, and make each recipe a single call. The installer rewrites the managed block in consumer justfiles, so also update this repo's root justfile. Depends on the REQ that fixes folder-structure-check for recipes with default parameter values. Verify: folder-structure-check reports no FOLDER-STRUCTURE-JUST-LOGIC on this repo, just run-kanban still replaces a stale instance on its port, just kanban-static still writes the build/ snapshot and the exclude line, and the board contract tests pass.*
