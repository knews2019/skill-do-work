---
id: REQ-716
title: 'Move the root installer scripts from tools/ to cmd/'
status: claimed
created_at: 2026-11-09T21:47:36Z
user_request: UR-168
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [skills/do-work/tools/lessons-do-work-update.md]
claimed_at: 2026-11-09T23:39:42Z
---
# Move the root installer scripts from tools/ to cmd/
## What
Move the five root installer scripts in `tools/` (`install-do-work-suite.sh`, `fetch-upstream-archive.sh`, `do-work-cli-preamble.sh`, `replace-text-section.sh`, `validate-suite-manifest.sh`) to `cmd/` with `git mv`, without breaking a fresh install or an update from an older installed version.
## Detailed Requirements
- Map every caller before moving: the README bootstrap (`source/tools/install-do-work-suite.sh`), the Go update path, the shipped skills (about 34 files name these five scripts), and the contract tests.
- Rule out a break for already-installed consumers: if an older installed updater runs `tools/<script>` from inside the newly downloaded archive, keep a thin `tools/` forwarder, or prove no installed version does. Record which one in this REQ.
- Layout rule this applies: scripts live in `cmd/` (rule 5); the maintainer chose `cmd/` over `tools/` in the audit.
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository.
**Why RED now:** It reports `FOLDER-STRUCTURE-ROOT-ENTRY tools/`.
**GREEN when:** `tools/` is no longer reported (or only a recorded forwarder remains), the README bootstrap installs the suite in a temp repository, an update from the previous release works, and the contract suites pass.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, `slugged: partial`; matched because the moved files are shipped shell.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: Move the five root installer scripts in tools/ (install-do-work-suite.sh, fetch-upstream-archive.sh, do-work-cli-preamble.sh, replace-text-section.sh, validate-suite-manifest.sh) to cmd/ (layout rule 5: scripts live in cmd/; the maintainer chose cmd over tools). Use git mv. First map every caller: the README bootstrap (source/tools/install-do-work-suite.sh), the Go update path, the shipped skills (80 files mention tools/; only the 34 that name these five scripts matter), and the contract tests. Rule out a break for already-installed consumers: if an older installed updater runs tools/<script> from inside the newly downloaded archive, keep a thin tools/ forwarder (or prove no installed version does) and say which in the REQ. Verify: folder-structure-check no longer reports tools/, the README bootstrap runs in a temp repo, an update from the previous release works, and the contract suites pass.*
