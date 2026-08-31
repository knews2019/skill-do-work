# REQ-436 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-436-audit-special-mode-preservation-in-remaining-file-publication`
Branch: `worktree-agent-REQ-436-audit-special-mode-preservation-in-remaining-file-publication`
REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-436-audit-special-mode-preservation-in-remaining-file-publication.md`
Exploration: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-436-exploration.md`

Implement only the frozen four-file scope:

- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go`
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go`

Requirements:

- Add RED-first fixtures for replacement, exclusive creation, and the real cleanup move using setuid, setgid, and sticky cases; assertions must project Go special flags back to Unix low-twelve-bit notation.
- Fix both cleanup masks. `CreateExclusiveAt` must apply the sanitized complete mode after writing and before sync; do not carry file-type bits.
- Preserve no-overwrite, containment, source revalidation, rollback, and source-deletion sequencing.
- Audit every scoped `Mode().Perm()` publication result and record the classification in the handback.
- Run focused tests, full CLI tests/vet, exact Go 1.25 compatibility, Windows atomic compile, gofmt, and diff/scope hygiene.

State stays home: do not edit or commit anything under `do-work/`, version/changelog files, primes, or unrelated files in the builder branch. Commit the four implementation files on the branch, keep the worktree clean, and write the handback only to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-436-handback.md`. Include commit hash, RED/GREEN evidence, commands/results, exact files, P-A-U evidence, decisions, discovered tasks, and readiness for integration.
