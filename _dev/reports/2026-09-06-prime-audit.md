# Prime audit — 2026-09-06

Seven source primes audited. Three oversized primes were reduced without deleting their guidance; remaining routing defects are recorded below. This audit changed documentation only.

## Changes made

- **Stakes:** added 3 (actions, shell, releases), refreshed 1 (updater), current 3 (CLI, board, maintainer-board delegation).
- **Shrink:** 3 primes; 137 nonblank original lines relocated into their paired satellites. Every original line of each shrunk prime remains in its prime/satellite pair; prior satellite histories are intact.
- **Lesson promotion:** added the 8 missing recurring family markers as generalized traps, retaining the incident evidence in the satellites.
- **Lesson index:** recomputed byte-based estimates, family sets and coverage for the 3 changed satellites. All 7 satellites have matching index rows and paired primes; no orphan, missing/dead row, estimate drift, family drift or false-full coverage remains.
- **CLI reading cost:** 25,431 → 9,154 bytes (approximately two-thirds less). Four existing Read first entries remain; the package catalog, phase table, scoped contracts and verification recipes are linked in the satellite.
- **Updater Stakes correction:** an equal upstream version is successful skipped work; only an older upstream version fails. Confirmed in `RunUpdate` rather than inferred from the old prime sentence.

| Prime | Lines before → after | Last commit date before audit |
| --- | ---: | --- |
| [prime-action-files.md](../../_dev/primes/prime-action-files.md) | 120 → 58 | 2026-09-04 |
| [prime-kanban-board.md](../../_dev/primes/prime-kanban-board.md) | 45 → 45 | 2026-09-04 |
| [prime-releases.md](../../_dev/primes/prime-releases.md) | 9 → 24 | 2026-09-06 |
| [prime-shell-commands.md](../../_dev/primes/prime-shell-commands.md) | 72 → 42 | 2026-09-06 |
| [prime-do-work-cli.md](../../skills/do-work/tools/do-work-cli/prime-do-work-cli.md) | 112 → 59 | 2026-09-06 |
| [prime-do-work-update.md](../../skills/do-work/tools/prime-do-work-update.md) | 36 → 36 | 2026-09-01 |
| [prime-do-kanban.md](../../skills/do-work-board/tools/queue-kanban/prime-do-kanban.md) | 53 → 55 | 2026-09-05 |

The whole-file ceiling is now met by all seven primes. Some retained trap paragraphs are still dense: audit preserves existing routing claims except for relocation and promotion from recorded lessons. Meeting a line limit alone does not establish that a prime is optimally concise.

## Remaining health findings — report only

The audit action makes routing content read-only except for its defined shrink/promotion operation. These are follow-ups, not silent repairs performed by this audit.

- [ ] **Stale updater traps:** [prime-do-work-update.md](../../skills/do-work/tools/prime-do-work-update.md) says the running binary is inside the replaced install and that staleness is mtime-based. [The launcher](../../skills/do-work/tools/do-work-cli.sh) now resolves a content-cached Go tool outside that source tree. Refresh these historical claims in a routing-content revision.
- [ ] **Mis-rooted updater pointers:** its Read first entries use `tools/...` and `actions/version.md` from a file already inside `tools/`. Resolve the CLI paths as `do-work-cli/...`, the launcher as `do-work-cli.sh`, and the action as `../actions/version.md`.
- [ ] **Other mis-rooted pointers:** [the board prime](../../skills/do-work-board/tools/queue-kanban/prime-do-kanban.md) cites `actions/board.md` from `tools/queue-kanban/` (the local target is `../../actions/board.md`). The CLI's relocated package-routing entry retains `tools/do-work-cli.sh`; from the CLI directory its target is `../do-work-cli.sh`.
- [ ] **Dead canonical history link:** [lessons-action-files.md](../../_dev/primes/lessons-action-files.md) links REQ-509 at the former flat archive path. The existing record is [UR-098/REQ-509](../../do-work/archive/UR-098/REQ-509-merge-common-rationalizations-into-one-crew-member.md). Refresh the canonical URL to include `UR-098/`.
- [ ] **Missing area index / cross-links:** the four cross-cutting primes in `_dev/primes/` have no common index and incomplete mutual routing. The core updater and CLI primes also lack reciprocal links. Do not solve this by making a shipped prime cite an export-ignored maintainer file.
- [ ] **Volatile detail in a compact-looking file:** [prime-kanban-board.md](../../_dev/primes/prime-kanban-board.md) is only 45 lines but over 10 KB, including historical ticket-count and browser-measurement detail. It would benefit from a separately scoped content revision that keeps durable traps and routes measurements to their evidence.

No additional utility prime was required by the module scan: both independent Go modules have primes. Cross-cutting entries registered in `CLAUDE.md` exist; utility primes are discovered by convention. No `file:///` portability violation was found (the board's `file://` protocol requirement is not a machine-local path).

## Verification

- Original-text conservation and exact relocated-block checks passed.
- All newly introduced relative Markdown links and section anchors resolve.
- Every repeated lesson-family marker is present in its paired prime's Traps section.
- `bash _dev/tests/contract-regressions.sh` passed, including core, board, suite manifest, shell, launcher and reference contracts. Heavy-only staged/install probes were explicitly skipped by that runner.
- `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/audit-lockins.sh` passed.

The dead REQ-509 canonical link above is a maintainer-satellite gap the current shipped-reference and relative-satellite checks do not cover. A passing existing check is not evidence that this finding disappeared.
