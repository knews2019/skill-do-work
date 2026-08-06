---
session_ended: 2026-08-06T11:16:00Z
last_completed: REQ-118
queue_state: 1 pending (REQ-114), 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 3
session_depth: moderate
---

# Session Checkpoint

## In Progress (interrupted)

## Completed This Session

- REQ-116: Normalize route at the board's read site and correct 0.174.15's board-wide claim (Route A) — v0.176.2, commit `2a2cd59` (review: Pass, 97%)
- REQ-117: An unrecognized domain must leave a footprint on the board, not become general in silence (Route A) — v0.176.3, commit `42f71e2` (review: Pass, 96%)
- REQ-118: The normalize flag must stop calling vocabulary-less field values unrecognized (Route A) — v0.176.4, commit `8d1a9f2` (review: Pass, 96%; closed UR-024)

Renumbered from 0.175.3/4/5 to 0.176.2/3/4 when main's own 0.176.0/0.176.1 (`just run-kanban-cli`) landed first and PR #133 hit a CHANGELOG/version conflict — the same collision REQ-111 hit at 0.174.14. All three came from an external review of the 0.174.15–0.175.2 series, triaged in this session via `do-work validate-feedback`. The review's fourth item (one-commit-per-REQ) was verdicted **Already done** and deliberately never captured.

## Still Queued

- REQ-114: `pending` — the three residual shell-logic extraction candidates, restated as greps rather than line numbers. **Still not approved work** (unchanged from last session): each needs its own floor decision, and each candidate's grep must be re-run because the as-of-census site counts are explicitly untrusted. UR-022 stays open in `user-requests/` — REQ-114 is its only member.

## Session Notes

- **The run was scoped to `UR-024`, not the whole queue.** The user's `do-work run` arrived with the trailing menu gloss "Process the captured fixes" attached. Read literally, work.md's unrecognized-argument guard would have halted the run; read as a full-queue default, it would have picked up REQ-114, which its own body marks as not-approved work. Scoping to the UR just captured was the only reading that matched intent. Worth remembering that the next-steps menu lines in this skill's own reports are copy-paste-shaped, so a description can arrive looking like an argument.
- **A `git add` whose pathspec list contains one non-existent path stages nothing at all** and exits non-zero — it is not partial. REQ-116's first commit therefore captured only the REQ file's rename with 0 content changes, because the list still named the pre-`git mv` queue path. Recovered with `git reset --soft HEAD~2` (both commits local and unpushed), the stale `commit:` line removed, then re-staged and re-committed. `record-commit-hash.sh --verify` is what caught it — the FAIL names the extra lines in the metadata commit, which is exactly the symptom of an implementation commit that never happened. **Check column 1 of `git status --porcelain` after every `git add`**, and note that `git mv` stages the *old* blob at the new path, so the content edits still need an explicit add afterwards.
- **`_dev/tests/contract-regressions.sh` still fails exactly 7 probes**, all in the update-script behavior section, byte-identical to the documented baseline and caused by running as root (the suite says up front it needs a non-root runner). Unrelated to this session's work. Run it as `bash _dev/tests/contract-regressions.sh` — the file is not executable in this checkout, and `./` returns permission-denied with exit 0 through a pipe, which reads as "0 failures" if you count FAIL lines without checking the verdict.
- **Recurring shape across all three REQs:** each was a *feedback leg* the contract already specified and the code silently skipped — route uppercasing, the domain warning, the vocabulary-less no-op. REQ-111 shipped the table and wired one field; three follow-ups were needed for the rest. When a contract is implemented as a table plus per-field wiring, the wiring is where it goes wrong, and a table test passes right over it.
- All three carry `kb_status: pending` — the lessons handoff was not run, consistent with REQ-110/111/112/113 and the REQ-104/108/109 backlog. Worth one `do-work bkb` pass over the whole set rather than per-REQ.
