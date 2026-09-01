# REQ-418 Sole Remediation Handback

## Range and scope

- Reviewed implementation: `a7c975c57822a7b050370142e4c0e817cce075c7`.
- Branch: `worktree-agent-REQ-418-toolbox-migration`.
- Remediation scope: 21 changed paths, all within the documented 32-path ceiling (the frozen 30 plus the conditionally activated `internal/gittransaction/git_transaction.go` and `_test.go`).
- The conditional expansion was necessary because tracked rollback had no caller API capable of distinguishing this invocation's publication from a later writer's replacement. The scratch plan records that decision before product edits.
- No lifecycle, retained-script, standalone-oracle, REQ-419/420, or REQ-462 inventory path changed.

## Review closure matrix

1. **Root confinement:** one `os.Root` publisher now validates linked ancestors, pins parent handles, stages privately, claims creations exclusively, and binds owned directories. Linked-ancestor fixtures prove note/architecture/portfolio/media paths cannot write outside the declared repository.
2. **last30days eligibility/rollback:** tracked targets refuse before clone or mutation; project containment and local-exclude eligibility precede publication; exclude undo and target restore compare publication identity/bytes and preserve replacements; publication uses an exclusive directory claim and verifies a complete payload before discarding backup evidence.
3. **Concurrent tracked replacement:** `gittransaction` captures each completed mutation's tracked publication before Git, restores only when inode/content still match, otherwise unstages and reports incomplete rollback while preserving the replacement. A failing-hook fixture proves the second writer's bytes survive.
4. **Media cancellation:** one context now spans backend generation and Git transaction. Backend groups use TERM/grace/KILL with live-group inspection after leader exit; Git and hooks run in a cancellable owned process group with escalation. Fixtures cover pre-cancel no-launch, TERM-deaf descendants, and blocking commit-hook cancellation plus rollback.
5. **Batch diagnostics/publication:** every failed worker produces typed `REPORT-IMAGE-MISSING` evidence and retained fallback text. `generated/` is claimed exclusively, its identity is pinned for child publication, and cleanup removes only the owned directory.
6. **Portfolio snapshot-first:** non-commit snapshot mode publishes the immutable snapshot in its own completed transaction before canonical publication. Canonical-directory and later canonical failures retain and report the snapshot rather than rolling it back.
7. **Architecture evidence/claims:** relative and absolute scan arguments use separate display/filesystem paths; `ReadDir` errors are findings. Bundle directories are exclusively claimed before hidden file publication and remain durably occupied after a later publication failure.
8. **Literal flag-shaped data:** mutation flags are recognized only in the leading option region; `--` ends it and preserves literal `--dry-run`/`--commit` arguments.
9. **Audit characterization:** option presence is distinct from sentinel value, closing all four explicit-default malformed cases. Committed tests build the standalone oracle, compare exact Markdown/status for all four modes on binary/unterminated/rename history fixtures, assert typed ordering, and mutation-check status/text/row-order comparison. The Go 1.25 lane skips only the standalone build because that retained module declares Go 1.26; canonical CLI tests still run there.
10. **Minor ReadDir:** unreadable/missing reports directories now return incomplete evidence.
11. **Minor P-A-U evidence:** named remediation tests now persist the previously unsupported adversarial claims. Per the no-lifecycle instruction, the orchestrator—not this branch—owns marking the REQ's lifecycle checkboxes from this handback and gate record.

## Verification

Green, unpiped where applicable:

- focused `gittransaction`, `toolboxcommands`, `resultmodel`, and `commandruntime` tests;
- race tests for `gittransaction` and `toolboxcommands`;
- full do-work-cli tests and vet;
- exact Go 1.25 compatibility;
- Windows toolbox cross-compile/fail-closed boundary;
- standalone audit tests/vet;
- valid four-mode exact Markdown differential and four malformed-option status differential;
- five focused retained toolbox fixture files and the complete 118-case prescribed-shell suite;
- contract regressions, staged-skills contract, installer behavior, and updater behavior;
- exact scope comparison and `git diff --check`;
- canonical `_dev/tests/maintainer-verify.sh` (passed; its optional browser lane declared the standard no-browser skip).

No gate failure was waived.
