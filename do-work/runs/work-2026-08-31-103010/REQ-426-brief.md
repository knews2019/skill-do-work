# Builder Brief — REQ-426

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-426-preserve-special-mode-bits`
Branch: `worktree-agent-REQ-426-preserve-special-mode-bits`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-103010/REQ-426-handback.md`
Route: A (focused direct implementation)
TDD: true

## Request

Preserve setuid, setgid and sticky bits on managed files instead of stripping them.

REQ-407's Go port reads file modes with `info.Mode().Perm()` (mask `0o777`) where the Python it replaced read `stat.S_IMODE(...)` (mask `0o7777`). The three special bits are silently dropped from `Justfile`, `CLAUDE.md` and `.claude/settings.json` on every install and update.

Requirements:

- `permissionsOf` in `managed_section.go` and the settings-mode read in `install_transaction.go` must carry setuid, setgid and sticky bits, not only the low nine.
- Add RED-first regressions: at minimum one `managedsection` case on a `0o2644` target and one install case asserting setgid survives a real install.
- Seed the broader installer fixture with `Justfile` 2644, `CLAUDE.md` 4644, and `.claude/settings.json` 1644 if the existing test structure supports the captured proof without widening scope.
- Keep ordinary permission behavior unchanged. The expected implementation is the mask `info.Mode().Perm() | (info.Mode() & (os.ModeSetuid|os.ModeSetgid|os.ModeSticky))`; do not add a syscall path.
- Record a concise root-cause note in the hand-back.

Allowed project files:

- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go`
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go`

Never touch any path under `do-work/` in the builder worktree. The only main-tree path you may write is the absolute hand-back path above. Do not modify release/version/changelog files.

## Required rules and context

Before editing, read these files completely from the main checkout:

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-shell-commands.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-shell-commands.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/specs/bug-fix.md`

Follow RED → GREEN → REFACTOR and record the failing and passing test commands/output. Run focused Go tests and any directly relevant installer tests. Commit the complete implementation on your branch; do not bump versions or edit changelogs.

## Hand-back format

Write the absolute hand-back file before returning. Include:

- branch name and final commit hash;
- P-A-U evidence (brief plan, applied scope, UNIFY checks and files inspected);
- RED/GREEN evidence;
- exact file manifest with `(modified)`, `(new)`, or `(deleted)`;
- tests run with direct results;
- root cause;
- integration seams, if any (exact line and destination; do not edit shared files);
- `## Decisions` and `## Discovered Tasks` sections when applicable, or state that there were none.

Return only a one-line status after the file is written.
