# REQ-718 replay repair integration

Merged into main with `git merge --no-ff`, hooks enabled, subject `[REQ-718] integrate archived experiment replay repair`.

- Merge: `b02547a7bdcd1053357067a6076a6b633b7341e2`.
- Parents: `93c1971d7649c900f5886b0234586212920bc84d` and worker `6c3d6d856e79c5bcd89e5fbe1e3aa6f401d58957`.
- Verified worker branch, HEAD, clean checkout, and exact four-file experiment-only diff before merging. Shipped source under `skills/` and `suite/`, and all `do-work/` tracked state, remain unchanged.
- Merged-tree checks passed (exit 0): `test-replay-progress.py` (one offline test), `verify-evidence.py` (eight source reads, fingerprints, boundary, continuation and process checks), AST syntax for all three Claude Python files, merge diff whitespace and working-tree diff checks. Evidence regeneration left no tracked changes.
- Removed clean merged worktree `/Users/t2/.t3/worktrees/skill-do-work2/t3-782e38ad` with `git worktree remove` and branch `t3/repair-req718-replay-progress` with `git branch -d`; neither used force. Verified both absent.

The archived request remains failed/environment because actual T3 compaction was unavailable. No archived status edits, recovery, finalization, live probe, full-suite repeat, push, PR, or Step10 ran. Scratch remains uncommitted and the run directory is preserved for coordinator cleanup.
