# REQ-718 archived experiment replay repair

Ready for coordinator integration. Claude progress now goes to timestamped, flushed stderr, and replay instructions no longer require the deleted original scratch paths. Historical evidence remains unchanged.

## Actual identities

- Branch: `t3/repair-req718-replay-progress`
- Worktree: `/Users/t2/.t3/worktrees/skill-do-work2/t3-782e38ad`
- Head: `6c3d6d856e79c5bcd89e5fbe1e3aa6f401d58957`
- Commit: `[REQ-718] make archived experiment replay independent of scratch paths`
- Working tree after commit: clean.

## Committed files

- `_dev/experiments/REQ-718/claude/run-probe.py`: remove absolute progress-file dependency; print UTC progress to stderr with flushing.
- `_dev/experiments/REQ-718/claude/test-replay-progress.py`: add focused offline regression for progress with filesystem opens failing.
- `_dev/experiments/REQ-718/claude/experiment-report.md`: describe portable progress and offline regression command.
- `_dev/experiments/REQ-718/codex-t3/probe-prompts.md`: require replay templates to substitute current checkout, archived request, and new local output paths. Preserve the quoted original brief verbatim.

## Checks and outcomes

1. Ran `python3 _dev/experiments/REQ-718/claude/test-replay-progress.py` before repair: exit 1, reproducing `FileNotFoundError` from `PROGRESS_PATH.open` when directory access is unavailable.
2. Ran the same test after repair: exit 0, one test passed. It loads the runner without executing its main entry point, disables `Path.open` to simulate absent scratch storage, calls the real progress function, and checks the message and UTC timestamp on stderr. No coordinator directory was deleted or moved.
3. Ran `python3 _dev/experiments/REQ-718/claude/verify-evidence.py`: exit 0. All eight complete source reads match, fingerprints match, manual boundary and function observations pass, continuation reads remain zero, sanitization and reaped-process assertions pass, and runner syntax passes. Regenerated `focused-checks.json` was byte-identical and is not part of the change.
4. Parsed every Claude experiment Python file with `ast.parse`: all three passed, without generating bytecode artifacts.
5. Ran `git diff --check` before commit and `git diff HEAD^ HEAD --check` afterward: exit 0.
6. Committed with ordinary `git commit`: exit 0. No hook bypass flags, environment overrides, or configuration changes were used. Git emitted no hook diagnostics.

## Replay-input review and limits

Reviewed both experiment directories, the three Claude prompt files, runner, verifier, and Codex replay instructions. Claude prompts use repository-relative paths. The verifier translates recorded absolute read paths relative to the recorded cwd before accessing this checkout, so historical worktree paths do not require the old worktree to exist.

The genuine additional defect was Codex replay step 1 instructing reuse of the original dispatch/brief, which referenced a now-archived working request and original scratch destinations. The replay procedure now explicitly substitutes current paths. Original transcripts, metadata, fingerprints, historical brief, observations, and captured CLI help were not edited.

No costly live model probe ran. No global settings, hooks, queue records, or archived request were modified. No delegates, pushes, merges, or worktree cleanup were performed. No child processes remain from verification. The coordinator owns integration and any broader gate.

## Rules and decisions

Read `CLAUDE.md` and crew rules `communication-style.md`, `general.md`, `coding-guardrails.md`, `anti-slop.md`, and `testing.md`. No applicable experiment-local primes or required lessons were found. Chose stderr rather than a new progress artifact or configurable path because progress is operational output and needs no filesystem dependency.

## Discovered Tasks

None beyond the replay defects repaired here.

Next action: integrate commit `6c3d6d856e79c5bcd89e5fbe1e3aa6f401d58957` through the main session's integration delegate.
