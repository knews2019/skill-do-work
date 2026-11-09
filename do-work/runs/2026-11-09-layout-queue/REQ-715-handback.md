# REQ-715 Hand-back

Branch `worktree-agent-REQ-715-changelog-archives-to-docs-changelog`, one commit: 742256f3.

## Files
- `docs/changelog/CHANGELOG-2026-04-07-up-to-v0.49.0.md`, `docs/changelog/CHANGELOG-2026-04-13-up-to-v0.64.1.md`, `docs/changelog/CHANGELOG-2026-07-07-up-to-v0.109.0.md`, `docs/changelog/CHANGELOG-2026-07-13-up-to-v0.121.0.md`, `docs/changelog/CHANGELOG-2026-09-05-up-to-v0.303.6.md` (moved from the root with git mv; line 3 of each modified)
- `.gitattributes` (modified)
- `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` (modified, still byte-identical)
- `_dev/tests/staged-skills-contract.sh` (modified)

## Verification (on the branch)
- `shipped-package-reference-contract.sh` exit 0 (changelog mirror matches, canonical links resolve).
- `contract-regressions.sh` exit 0.
- `staged-skills-contract.sh` in the heavy tier exit 0.
- `git archive HEAD | tar -t` lists no `docs/` path and no dated archive.
- `folder-structure-check` reports no archive (7 gaps left, down from 12).

## Decisions
- D-01 (DECIDE & STATE): the export-ignore rule is the folder `/docs/changelog`, not a file glob. Everything cut into that folder is covered without a pattern to keep in sync.
- D-02 (DECIDE & STATE): links inside old release entries are left as written (the ADR link at line 2221 and the archive list at line 2433 of the 0.303.6 archive). The archives say they hold the release notes verbatim, like the records under `do-work/`. Only each file's own navigation link on line 3 follows the move.
- D-03 (DECIDE & STATE): no release or archive-cut code names the root location, so nothing in Go changes. The `.gitattributes` comment now documents that a new archive is cut into `/docs/changelog`.

## Discovered Tasks
- The 0.303.6 archive links `blob/main/decisions/records/adr-020-...` at line 2221, which breaks on GitHub now that REQ-714 moved the ADR log. Verbatim history, left as written. impact-negligible → report only

## Lessons read
- `_dev/primes/lessons-releases.md` (whole satellite, the REQ's required lesson): applied `canonical-link-outlives-its-target` (the header links of both live changelogs and each archive's own link moved in the same commit as the files).

## Addendum (review F1, F2)

- F1: each archive's `**Current:**` link pointed at `./CHANGELOG.md`, which does not exist in `docs/changelog/`. It now points at `../../CHANGELOG.md`. The first sweep searched only for the dated archive names and missed it.
- F2: the `.gitattributes` comment names all three lines the staged-skills contract asserts.
- Second commit on the branch; the reference contract exits 0.
