---
id: REQ-715
title: '[impact-negligible] Move the dated changelog archives to docs/changelog/'
status: claimed
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
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-09T23:10:49Z
  basis:
    - Route B
    - 9-file write set
    - 4 acceptance criteria
write_set: [.gitattributes, CHANGELOG.md, skills/do-work/CHANGELOG.md, _dev/tests/staged-skills-contract.sh, docs/changelog/CHANGELOG-2026-04-07-up-to-v0.49.0.md, docs/changelog/CHANGELOG-2026-04-13-up-to-v0.64.1.md, docs/changelog/CHANGELOG-2026-07-07-up-to-v0.109.0.md, docs/changelog/CHANGELOG-2026-07-13-up-to-v0.121.0.md, docs/changelog/CHANGELOG-2026-09-05-up-to-v0.303.6.md]
claimed_at: 2026-11-09T23:09:15Z
route: B
dispatch_at: 2026-11-09T23:13:02Z
builder_handback_at: 2026-11-09T23:14:43Z
integration_at: 2026-11-09T23:14:44Z
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
- [x] **[PLAN]:** git mv the five archives into docs/changelog in one commit with the export-ignore rule, the header links of both live changelogs, each archive's own link, and the staged-skills exclusion and assertion.
- [x] **[APPLY]:** One commit, 742256f3: five renames and four modified files, all in Scope.
- [x] **[UNIFY]:** git status before the commit: 5 R and 4 M, nothing else. The two live changelogs are still byte-identical (cmp). shellcheck clean on the changed script. The reference contract, contract-regressions and the staged-skills contract exit 0 on the branch.
*Source: do-work capture-request: Move the five dated changelog archives CHANGELOG-2026-*-up-to-v*.md from the repository root to docs/changelog/ (layout rule 1: the root CHANGELOG.md covers the whole repository; old changelog archives go to docs/changelog/). Depends on the REQ that moves kb/ and decisions/ under docs/. Use git mv. Rewrite the 15 live references: the GitHub blob links in skills/do-work/CHANGELOG.md and the root CHANGELOG.md, the /CHANGELOG-20*.md export-ignore glob and its comment in .gitattributes (it becomes /docs/changelog), and any release or archive-cut code or prime that names the root location (check _dev/primes/prime-releases.md), so the next archive is written to docs/changelog/. Verify: folder-structure-check no longer reports the archives, git archive HEAD does not contain docs/changelog/, the release --dry-run still parses the live CHANGELOG.md, and the contract suites pass.*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (five files move to one folder), but the REQ asks to find any release or archive-cut code and every reference that names the root location, which needs a short exploration before the move.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**No archive-cut code exists.** Cutting an archive is a manual rename plus a pointer edit, as the `.gitattributes` comment says. No Go code, shell test or prime names the root location of the dated archives: `git grep` for the archive name pattern outside `do-work/`, `docs/kb/` and `ai-reports/` finds only `.gitattributes`, the two live changelogs and the archives themselves. The release tool's insertion anchor is the last archive link line of the live changelog; it ends with the file name, so it still matches after the path in the link changes.

**References to rewrite (9 files, not the 15 the REQ estimated).**
- `.gitattributes` line 46 (`/CHANGELOG-20*.md export-ignore`) and its comment at lines 27-33.
- The five GitHub links in the header of `CHANGELOG.md` and of its byte-identical mirror `skills/do-work/CHANGELOG.md`.
- Line 3 of each archive: its own GitHub link ("a tarball install browses it on GitHub").
- `_dev/tests/staged-skills-contract.sh`: `docs` is a legacy runtime path, so line 108 must also exclude the new folder, and the export-ignore assertion loop at line 74 should cover it.

**Left as they are.** The Older and Newer links between archives are relative (`./CHANGELOG-...`) and keep working inside one folder. Links inside old release entries (for example line 2433 and the ADR link at line 2221 of the 0.303.6 archive) are verbatim history, like the records under `do-work/`.

**Baseline.** `folder-structure-check` reports five `FOLDER-STRUCTURE-ROOT-ENTRY CHANGELOG-2026-*` gaps. `git archive HEAD` lists no archive file.

*Explored by work action (orchestrator, inline)*

## Scope

**Files I will touch:**
- `.gitattributes` (modify)
- `CHANGELOG.md` (modify)
- `skills/do-work/CHANGELOG.md` (modify)
- `_dev/tests/staged-skills-contract.sh` (modify)
- `docs/changelog/CHANGELOG-2026-04-07-up-to-v0.49.0.md` (modify)
- `docs/changelog/CHANGELOG-2026-04-13-up-to-v0.64.1.md` (modify)
- `docs/changelog/CHANGELOG-2026-07-07-up-to-v0.109.0.md` (modify)
- `docs/changelog/CHANGELOG-2026-07-13-up-to-v0.121.0.md` (modify)
- `docs/changelog/CHANGELOG-2026-09-05-up-to-v0.303.6.md` (modify)

The five archive files are moved from the repository root with git mv; the only content edit in each is its own GitHub link on line 3.

**Files I will NOT touch:** the text of old release entries inside the archives; anything under do-work/, docs/kb/ and ai-reports/; the release version files.

**Acceptance criteria (restated from REQ):**
- [ ] The five dated archives live under docs/changelog and were moved with git mv.
- [ ] The links in both live changelogs and the export-ignore rule point at the new folder, and the next archive cut is documented to land there.
- [ ] folder-structure-check no longer reports the archives, and git archive HEAD contains no docs/changelog path.
- [ ] A release dry run still parses the live changelog, and the contract suites pass.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/` at 2ed0be64; only this REQ's working file, `baseline.json` and the run directory are changed
**Tests baseline:** ✓ Focused probe green (`bash do-work/runs/2026-11-09-layout-queue/REQ-715-probe.sh`, the shipped reference contract); repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 146s at 2ed0be64, run from a detached checkout
**Dependencies:** ✓ None added

*Checked by work action*

## Decisions

- D-01 (DECIDE & STATE): the export-ignore rule is the folder `/docs/changelog`, not a file glob. Everything cut into that folder is covered without a pattern to keep in sync.
- D-02 (DECIDE & STATE): links inside old release entries are left as written (the ADR link at line 2221 and the archive list at line 2433 of the 0.303.6 archive). The archives say they hold the release notes verbatim, like the records under `do-work/`. Only each file's own navigation link on line 3 follows the move.
- D-03 (DECIDE & STATE): no release or archive-cut code names the root location, so nothing in Go changes. The `.gitattributes` comment now documents that a new archive is cut into `/docs/changelog`.
- D-04 (DECIDE & STATE, review F1): each archive's `Current` link now reads `../../CHANGELOG.md`. The first sweep searched only for the dated archive names and missed this link to the live file.

## Discovered Tasks

- The 0.303.6 archive links `blob/main/decisions/records/adr-020-...` at line 2221, which breaks on GitHub now that REQ-714 (knowledge base and ADR log under docs/) moved the ADR log. Verbatim history, left as written. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `.gitattributes` (modified)
- `CHANGELOG.md` (modified)
- `skills/do-work/CHANGELOG.md` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `docs/changelog/CHANGELOG-2026-04-07-up-to-v0.49.0.md` (modified)
- `docs/changelog/CHANGELOG-2026-04-13-up-to-v0.64.1.md` (modified)
- `docs/changelog/CHANGELOG-2026-07-07-up-to-v0.109.0.md` (modified)
- `docs/changelog/CHANGELOG-2026-07-13-up-to-v0.121.0.md` (modified)
- `docs/changelog/CHANGELOG-2026-09-05-up-to-v0.303.6.md` (modified)

**What was done:** The five dated changelog archives moved from the repository root to `docs/changelog/` with git mv, in one commit. Each archive's own GitHub link and its link back to the live changelog, the five header links of both live changelogs, and the export-ignore rule (now the folder `/docs/changelog`) follow the move. The staged-skills contract excludes the new folder from its legacy `docs` check and asserts its export-ignore line. The `.gitattributes` comment says a new archive is cut into that folder.

## Qualification

**Range:** `2b48f7db..1a000d08` (two merges of the REQ branch: 742256f3, then the review fix ab39e61c; first merge f2ccb2a7)
**Mechanical gates:** qualify ✓ satisfied, scope-drift ✓ satisfied (9 declared paths, 9 touched).

**Requirement trace (orchestrator, read against the diff):**
- Move: five `R` entries from the root to `docs/changelog/`. Each archive has two one-line edits: its own GitHub link and, after review finding F1, its `Current` link to the live changelog (`../../CHANGELOG.md`).
- References: the header of `CHANGELOG.md` and its mirror link `blob/main/docs/changelog/...` five times each, and the two files are byte-identical. `.gitattributes` carries `/docs/changelog  export-ignore` in place of the root glob, and its comment says where a new archive is cut.
- Release and archive-cut code: none names the root location (Exploration), so no Go change; the release tool's insertion anchor still occurs exactly once in the live changelog.
- The staged-skills contract excludes the new folder from the legacy `docs` scan and asserts its export-ignore attribute.
- The REQ counted 15 live references; 9 files hold them. Links inside old release entries stay verbatim (D-02).

**Live flow:** at the merge, `folder-structure-check` reports no archive, and `git archive` of the merge lists no `docs/` path and no dated archive file.

*Qualified by work action*
