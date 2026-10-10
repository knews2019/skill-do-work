# Integrator brief: REQ-685 (git transaction rollback reuses the transaction's open repository root; the no-root rollback path is deleted)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-685.

- REQ id: REQ-685. Working REQ path: `do-work/working/REQ-685-rollback-uses-open-root.md`. UR: UR-152 (`do-work/user-requests/UR-152/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-100748/` (run id `work-2026-10-10-100748`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-100748/REQ-685-handback.md`. Builder summary: one commit `62306f91` on base `b629e5cd`; 4 files, +45/-252 (net -207). `rollbackFailure` takes the open root and calls `rollbackWithRoot` only; `rollbackWithoutRoot`, `openRollbackRoot`, the fallback branch and the two no-root tests are deleted; one test with a deleted subtest was collapsed (de-indented); three comments that named the deleted code fixed; unused `syscall` import removed; two-line "literal HEAD is deliberate" comments at `exact_commit.go` and `git_transaction.go`; `_dev/tests/audit-lockins.sh` ratchet comment updated, ratchet code untouched. gittransaction 0 (10.9 s), finalization 0 (64.2 s), probe 0, audit-lockins 0, vet/gofmt clean, no remaining references.
- Review focus: (1) the REQ-598 nil-handle panic must stay impossible: confirm no `rollbackFailure` call site can pass a nil root (every call is after the open at ~line 520 succeeds). (2) Confirm the collapsed test still asserts what its kept subtest asserted. (3) Anti-bloat: count added lines that are not comments.
- Guide section 12 "The merged build runs your own tools" applies: your finalization runs the merged transaction code.
- Operative name (branch = worktree basename): `worktree-agent-REQ-685-rollback-uses-open-root`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-685-rollback-uses-open-root`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T10:29:20Z`.
- Route B, tdd: false, impact-negligible, effort-substantive, domain backend. Pre-flight green gate already recorded by the coordinator. Test-gate probe: `do-work/runs/work-2026-10-10-100748/REQ-685-probe.sh`.
- Integration order (coordinator ruling): hand-back order. Earlier integrators of this run bump VERSION first; re-read VERSION right before the payloads. REQ-686 stays last.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer records its own **Restatement sweep:** line: grep `lessons-do-work-cli.md`, `_dev/primes/`, actions and docs for "rollbackWithoutRoot", "no-root rollback", "reopens the root" or REQ-598 wording that describes the deleted path, and report what still restates it (the shipped CHANGELOG history stays as is).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, cleanup); the coordinator does after you report.
- Lessons: none proposed by the builder.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-685-rollback-uses-open-root.md`.
- Changelog title direction: plain words, what shipped, for example "Transaction Rollback Reuses the Repository Folder It Already Opened, and the Second Rollback Path Is Gone"; the builder's proposed entry is in the hand-back.
