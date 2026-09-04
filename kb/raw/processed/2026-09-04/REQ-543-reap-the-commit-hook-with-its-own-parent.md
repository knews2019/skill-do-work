---
source_type: req_lesson
req_id: REQ-543
req_path: do-work/archive/REQ-543-reap-the-commit-hook-with-its-own-parent.md
date: 2026-09-03
domain: backend
module: skills/do-work/tools/do-work-cli
tags: [backend, reap, commit, hook]
---

# Lessons from REQ-543: Reap the commit hook with its own parent

## What the REQ was about

Make cancellation of a media Git transaction terminate the whole owned process group, not just the direct `git` child. A `pre-commit` hook that ignores `SIGTERM` currently outlives the cancelled transaction and keeps running after the command has returned.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group.go` (new)
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix.go` (new)
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unsupported.go` (new)
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_cancellation_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_unix.go` (deleted)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_windows.go` (deleted)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

## What worked

- Measuring before diagnosing. Process-table sampling at 20 ms during a live failing run refuted three plausible causes (`Setpgid` is applied, the negative-PID group signal is delivered, `trap '' TERM` is inherited by `sleep`) and produced the descendants-first shape directly. A guessed fix here would have been "make the escalation synchronous", which does not work.
- Testing every guard by breaking it. The neuter table is what turned "the tests pass" into "these four tests red when the fix is removed" — and it is the same instrument that found three guards nothing reds on (F3, F12). A green suite would have hidden both.
- Cross-building for plan9. The builder's first attempt named the non-Unix fallback `*_windows.go`, whose implicit constraint left every other non-Unix target with no implementation. One cross-build caught it before merge; `_unsupported.go` is the fix.

## What didn't work

- The builder reported four always-load crew-member files as absent from the tree. All 18 are present and tracked. The guardrails were never read; the orchestrator read the diff against them afterwards instead. Filed as a discovered task, because if the check was a relative `ls` from the Go module directory then every builder dispatched into a subdirectory silently skips the always-on guardrails.
- The Scope declaration. Prose paths ("a new shared package under `.../internal/`") instead of literal ones, and backticked identifiers in the trailing descriptions, which `scope-drift.sh` reads as declared paths. **Second occurrence** — REQ-457 hit it first. The fix is that a Scope bullet's description carries no backticks; recording it a third time would not be a fix.
- `write_set` was never amended when the file list grew from 2 to 10. `## Scope` drift was reconciled in prose and the frontmatter guard was not, and under fan-out it is the guard, not the prose, that prevents collisions (F10). Corrected in this record's frontmatter at closure.
- The first sweep implementation re-derived terminable leaves each round and **hung the `generate-report-image` shell probes forever**: a backend that respawns a helper in a loop always has a live child, so the sweep never climbed to the parent. D-04's single-snapshot levelling is the replacement.

## Worth knowing

- **A zombie satisfies `kill(pid, 0)`.** That one fact is why making the escalation synchronous is not sufficient and why the ordering had to change: only letting `git` survive to `waitpid()` its own hook removes the window, because a wait on our own child ends in milliseconds while a wait on init is measured in seconds.
- `ESRCH` from `Getpgid(leader)` does **not** mean "the group is not isolated" — the leader can be reaped while the group is still alive and running. `internal/nextselection`'s `cleanupReapedProcessGroup` (`blocked_probe_unix.go:54`) already handles exactly this case and is the reference behaviour F1 lacks.
- Without `PR_SET_CHILD_SUBREAPER` — barred here, since it lives in `x/sys/unix` and this module is standard-library-only — an orphan can be proved not-running but never reaped, and a `setsid` escapee cannot be reached at all. Any future wording of this contract has to carry all three outcomes.
- The prime's `GOOS=windows` Verify lines are run by no gate: `grep -n GOOS _dev/tests/*.sh` is empty. The portability contract is written down and unenforced.
- The REQ's `## Traps` names `_dev/tests/do-work-cli-go125-compatibility.sh` as the Go 1.25 floor enforcer. Main deleted that script in `5e0e166` and the merge removed its prime Verify line, so the floor now rests on `go.mod`'s `go 1.25.0` alone. Not this REQ's doing, but the Traps text is stale for the next reader.
- New lesson family **`reaped-by-its-own-parent`**, promoted into the prime's `## Traps` (`prime-do-work-cli.md:56`) by the implementation and now carried by its satellite bullet in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as well. The gap the review recorded as finding F9 — the index row at `do-work/lessons-index.md:11` listing a family the satellite did not carry, which contradicts the index's own header rule ("`families` is the exact sorted set of `[family: <slug>]` markers present in lesson bullets") — is closed in the release commit below, together with the `tokens` recount the header rule prescribes (`ceil(bytes / 4)`: 5127 to 5660). The disagreement recorded at hand-back stands as written: the two halves belonged in one commit and were split only because that closure could not write under `skills/`.

## Back-reference

See `do-work/archive/REQ-543-reap-the-commit-hook-with-its-own-parent.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1cc3beb`.
