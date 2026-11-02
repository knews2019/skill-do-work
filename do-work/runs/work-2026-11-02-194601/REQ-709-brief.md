# REQ-709 builder brief

You are the builder for REQ-709 (recover ignores loose do-work files that no request record can own). Route A: triaged as simple, aim for a focused minimal change. `tdd: true`.

## Where you work

- Your worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-709-recover-ignores-unowned-do-work-notes` on branch `worktree-agent-REQ-709-recover-ignores-unowned-do-work-notes` (already created from `main`). Do every edit, test and commit there.
- Every commit subject starts with `[REQ-709]`. Commit by exact path (never `git add -A`, never `git commit -a`). Never push.
- Never write the main tree (`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`) except your hand-back file, written by absolute path: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-194601/REQ-709-handback.md`. Never stage or commit it.
- Never write any `do-work/` path inside your worktree. The REQ file and the lessons satellite are the orchestrator's (the lesson bullet is written at finalization).
- Scratch files go under a private folder such as `/tmp/req709-builder/`, never the shared scratchpad root.

## Read first

1. The REQ (read-only, main tree): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-709-recover-ignores-unowned-do-work-notes.md`. Its `## Red-Green Proof` is user-confirmed: anchor RED on it.
2. Crew rules (in your worktree): `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `backend.md`, `testing.md`.
3. Prime: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`.
4. Lessons (the whole satellite was dropped for budget, but the touch-conditional rule applies because you change recovery evidence): read only the bullets carrying `[family: closed-enumeration-for-a-condition]` and `[family: observed-subset-is-not-semantic-completeness]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (use fixed-string grep). Name what you read in the hand-back.

## The bug

`recover` (and `recover-finalization --discover`) refuse `FINALIZATION-DISCOVERY-AMBIGUOUS` when the only uncommitted file is a note under `do-work/` that no request record can own, such as an untracked `do-work/inbox/note.md`. Cause: `sharedFinalizationPath` in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` returns true for every path under `do-work/`, so the note lands in `unresolvedShared` (when legacy candidates exist, as in this repository) or in `ambiguousSharedRemainder` (when none exist). `finalization_apply.go:525` uses the same predicate for a discovered journal's commit-safety check. Reproduced at claim time in the main repository.

## Orchestrator design calls (log them in your hand-back `## Decisions` as D-01 to D-03, DECIDE & STATE)

- **D-01** The condition is "where request records live", stated once in code: the queue, working, archive and user-request folders (`do-work/queue/`, `do-work/working/`, `do-work/archive/`, `do-work/user-requests/`), REQ reservation markers (`do-work/.req-reservations/`), run records (`do-work/runs/`), and the two lifecycle record files every finalization writes, `do-work/CHECKPOINT.md` (claim records) and `do-work/calibration-log.tsv` (one row per finished request). Those two files stay evidence because the REQ requires real unfinished finalization work to refuse or settle exactly as today. Do NOT list note folders such as `inbox/`, `audits/` or `HANDDOWN-*`: every other file under `do-work/` is covered by the condition. Doc comment states the condition and why (a note has no request owner, and `do-work commit` must be able to commit it).
- **D-02** The rule lives in the finalization package, not in the inventory chokepoint `untrackedDoWorkMetadata` (`internal/corehelpers/inventory.go`): `do-work commit` reads that inventory and must still list the note so it can commit it.
- **D-03** Discovery neither refuses on such a file nor associates it with a request: drop it from discovery's `dirty` set (path and rename origin judged one at a time, so a rename out of a record folder still leaves its record-side path as evidence), and the narrowed predicate also narrows the apply-time check at `finalization_apply.go:525`. Consequence to state in the hand-back: other non-record files directly under `do-work/` (for example `lessons-index.md`) also stop being evidence, so a legacy tail now settles around them and leaves them dirty for `do-work commit`, where today it refused.

Keep the change minimal. No new package, no new exported API unless a test truly needs it. Keep the existing refusal codes and messages unchanged.

## Tests (RED first, then GREEN)

Write the failing tests first, run them, and record the exact failure text before you change production code. Keep each test file under 30 s. Name each new test after the real failure it pins and put a one-line comment above it saying which failure (the 2026-11-02 UR-161 run refusal on `do-work/inbox/2026-10-28_...md`).

Cover at least:
1. Finalization discovery (`internal/finalization`, same helpers as `TestRecoverFinalizationIgnoresFinderMetadataUnderDoWork` in `finalization_recovery_test.go`): an untracked `do-work/inbox/note.md` with no terminal request → success, no finalization; a legacy tail (`seedSemanticLegacyTail`) recovers to `cleanup_complete` around the note and the note stays untracked and uncommitted; a tracked note outside the record folders (for example `do-work/HANDDOWN-UR-1.md`) that is modified is ignored the same way.
2. The negative (the REQ's "exactly as today"): real unowned shared dirt next to the note still refuses, for example a foreign claim line in `do-work/CHECKPOINT.md` (see `TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest`), and the refusal's paths name the checkpoint but not the note. Also show a dirty file inside a record folder still counts (for example an untracked file under `do-work/runs/<dir>/` or `do-work/working/`).
3. The public `recover` composition (`internal/lifecycleadvance`, `recovery_commands_test.go` style): the REQ's RED case, a repository whose only uncommitted file is an untracked `do-work/inbox/note.md`, returns success; and the protected inventory (`corehelpers.ReadProtectedInventory`) still lists the note as an `A` row, which is what `do-work commit` groups as an ordinary change with no owning request. If importing corehelpers there is awkward, put this assertion where it fits and say where.

## Verify before hand-back

From `skills/do-work/tools/do-work-cli` in your worktree: `gofmt -l ./internal/finalization ./internal/lifecycleadvance` (empty), `go vet ./internal/finalization/ ./internal/lifecycleadvance/`, `go test ./internal/finalization/ ./internal/lifecycleadvance/ -count=1`. Do not run `_dev/tests/maintainer-verify.sh` or any heavy lane.

## Hand-back (write to the absolute path above, then end)

Sections: `## Summary` (2-4 sentences), `## Files` (each path with created/modified/deleted, tests included), `## Commits` (hashes and subjects), `## Red-Green Evidence` (test name, exact RED failure text, GREEN pass), `## P-A-U` (your [PLAN] approach, [APPLY] scope note, [UNIFY]: `git diff --stat main...HEAD`, gofmt/vet results, each file checked and what you checked), `## Lessons Read` (entries, family-targeted, anything missing), `## Decisions` (D-01 to D-03 as given, then D-04 onward for your own choices, each with reasoning; ESCALATE ones add Value and Risk), `## Discovered Tasks` (out-of-scope finds with an impact guess, or "none"), `## Test Commands` (exact commands and results).
