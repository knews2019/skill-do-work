---
id: REQ-715
title: '[impact-negligible] Move the dated changelog archives to docs/changelog/'
status: pending
created_at: 2026-11-09T21:47:36Z
user_request: UR-167
domain: general
prime_files: [_dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
depends_on: [REQ-714]
required_lessons: [_dev/primes/lessons-releases.md]
---
# [impact-negligible] Move the dated changelog archives to docs/changelog/
## What
Move the five dated changelog archives `CHANGELOG-2026-*-up-to-v*.md` from the repository root to `docs/changelog/` with `git mv`, and make the next archive cut land there too.
## Detailed Requirements
- Rewrite the 15 live references: the GitHub blob links in `skills/do-work/CHANGELOG.md` and the root `CHANGELOG.md`, and the `/CHANGELOG-20*.md` export-ignore glob plus its comment in `.gitattributes` (it becomes `/docs/changelog`).
- Update any release or archive-cut code or prime that names the root location (check `_dev/primes/prime-releases.md`), so the next archive is written to `docs/changelog/`.
- Depends on REQ-714 (suite writes kb/ and decisions/ under docs/), which creates `docs/` first.
- Layout rule this applies: the root `CHANGELOG.md` covers the whole repository; old changelog archives go to `docs/changelog/` (rule 1).
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository.
**Why RED now:** It reports five `FOLDER-STRUCTURE-ROOT-ENTRY CHANGELOG-2026-*` gaps.
**GREEN when:** None of the archives is reported, `git archive HEAD` does not contain `docs/changelog/`, a `release --dry-run` still parses the live `CHANGELOG.md`, and the contract suites pass.
**Validation:** User confirmed
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: Move the five dated changelog archives CHANGELOG-2026-*-up-to-v*.md from the repository root to docs/changelog/ (layout rule 1: the root CHANGELOG.md covers the whole repository; old changelog archives go to docs/changelog/). Depends on the REQ that moves kb/ and decisions/ under docs/. Use git mv. Rewrite the 15 live references: the GitHub blob links in skills/do-work/CHANGELOG.md and the root CHANGELOG.md, the /CHANGELOG-20*.md export-ignore glob and its comment in .gitattributes (it becomes /docs/changelog), and any release or archive-cut code or prime that names the root location (check _dev/primes/prime-releases.md), so the next archive is written to docs/changelog/. Verify: folder-structure-check no longer reports the archives, git archive HEAD does not contain docs/changelog/, the release --dry-run still parses the live CHANGELOG.md, and the contract suites pass.*
