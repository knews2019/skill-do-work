# Builder brief: REQ-694 (req append-section refuses a body that hides later sections; frontmatter set refuses status and id)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-694-append-section-and-set-guards
- Branch: worktree-agent-REQ-694-append-section-and-set-guards, created by the coordinator from main HEAD with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch; never create the worktree yourself).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-694-append-section-hidden-sections-and-set-status-guard.md. Read it fully: What, Prior Implementation, Detailed Requirements, Red-Green Proof, Constraints, Builder Guidance, Required Lessons Dropped for Budget, and the orchestrator's `## Triage` (decisions D-01 to D-05; your first decision is D-06). The release belongs to the integrator; you PROPOSE the changelog entry and a lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-154/input.md (review follow-ups R1 to R5; this REQ is R2 and R3).
- Original work (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-145/REQ-659-frontmatter-set-and-req-append-section.md and the findings it left: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-131527/REQ-659-rereview.md (N1) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-131527/REQ-659-review.md (F3).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-192201/REQ-694-handback.md
- Route A, tdd: true, impact-user-visible, effort-mechanical, domain backend. Two small guards in one Go file plus their tests.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, backend.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md, /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release itself is the integrator's). Lessons: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file, small). The capture dropped `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` for budget (20802 tokens); read only its `[family: rule-direction-checked-against-callers]` bullets (lines 13 and 34) and `[family: lifecycle-section-evidence]` bullets (lines 40 and 93), because both guards change what a shared writer accepts and how a section is counted.

## The change (decided, do not reopen)
All edits are in `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` and `request_writers_test.go` of your worktree.

1. N1, `handleRequest` (`request_writers.go:70-176`). In the existing loop at `:121-130` (it already skips `HeadingIndent != 0`), also count the column-0 visible sections of the body before the insert. After `document.ReplaceBodySpan` (`:158`), count the column-0 visible sections of `document.BodyBytes()` and refuse `SECTION-WRITE-FAILED` unless the new count equals the old count plus one. This REPLACES the single-copy re-check at `:161-171` (D-02): a hidden new heading does not raise the count, so the end-of-file case it guarded stays covered. Rewrite the comment at `:161-162` to say what the count check protects (an open fence or comment in the body, or one already open at the end of the file, hides the new heading or every later section). Keep the refusal before `replaceRequestFile`, so the file stays byte-identical. The pre-insert heading check at `:102-106` stays as it is.
2. F3, `handleFrontmatterSet` (`request_writers.go:25-65`). Right after the argument parsing at `:30`, before the stamp checks and before `resolveActiveRequest`, refuse `field == "status"` and `field == "id"` with a new code `FRONTMATTER-FIELD-OWNED` (D-03) through `writerRefusal(code, target, evidence)`. The evidence names the owner. For status: the lifecycle commands (`advance`, `unblock`, `finalize`) or the hand write at the transition's defining site change it, never `frontmatter set`. For id: an id is never rewritten. Every other scalar field keeps today's behaviour. Update the doc comment at `:21-24` with one sentence on the refusal.
3. Tests in `request_writers_test.go` (fixture `writerFixtureRequest` at `:13`, helpers `writeMatrixFile`, `runRegisteredWriter`, `readFixture`, `assertRefusedUnchanged` at `:44`). Use exactly these names (D-05), because the GREEN probe names them:
   - `TestRequestAppendSectionRefusesBodyThatHidesLaterSections`: a table of three rows, each `--section Testing` into a REQ that is refused `SECTION-WRITE-FAILED` with the file byte-identical. Row 1: body `## What` then `## Review`, `--from` file holding a line of three backticks then `## x`. Row 2: same REQ, `--from` file holding `text <!-- open`. Row 3 (D-04): a REQ whose body ends inside an unclosed fence and has no later canonical section, with an ordinary `--from` body (the end-of-file case the removed check guarded).
   - `TestFrontmatterSetRefusesLifecycleOwnedFields`: `frontmatter set REQ-701 status bogus-status` and `frontmatter set REQ-701 id REQ-999` are each refused `FRONTMATTER-FIELD-OWNED` with the file byte-identical, and the finding's evidence names the owner (for status it contains `advance`; for id it says the id is never rewritten).
   - Each test opens with a one-line comment naming the failure it pins, like the existing tests (for example "Pins review N1: an open fence in the body hid the later ## Review from advance.").

## Anti-bloat (YAGNI), from the REQ's Constraints; the maintainer asked for this to be watched
- Smallest change that fixes the named failure. Prefer deleting or reusing existing code over adding.
- No new allowed-values table for other fields: only `status` and `id` are refused.
- No new package-level helpers, flags, options, config or files. A one-line local count closure inside `handleRequest` is fine; a new package-level function is not (REQ-693 edits `handoff.go` in the same Go package at the same time, so a new package-level name can collide).
- Tests only pin the named failure (one test per failure); no decorative or "while we're here" tests.
- Do not fix adjacent things you notice; record them as discovered tasks instead.
In the hand-back, paste `git diff --stat` and list every function, constant, option, file or test you added that the REQ did not name (expected: the new refusal code only). Each one gets a one-line reason, or you remove it.

## Write boundary
Exactly `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` and `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go`. Anything else: stop and say so in the hand-back. No prose file changes: no shipped action tells an agent to run `frontmatter set` on `status` or `id` (checked at pre-dispatch), so no caller needs an edit.

## Hard rules
- Every commit subject on your branch starts with `[REQ-694]` (for example `[REQ-694] refuse hidden-section bodies and owned frontmatter fields`). The board credits commits by that prefix only.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (which you write but never stage or commit). The integrator's queue guard refuses a branch that touches `do-work/`.
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, `just do-work-update`, or the repository gate (`_dev/tests/maintainer-verify.sh`); the integrator owns the gate.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or any version mirror: the release belongs to the integrator.
- Go: run focused `go test -run` and `go vet` on `internal/corehelpers` only, never the whole module and never the repository gate.
- Scratch files and fixtures go in `mktemp -d` directories outside both trees and are never committed. Clean up anything you start.
- Load is high while four builders run at once: if a wall-time budget fails once under load, rerun it once and record both runs.

## Integration seam
No shared file with the other three builders. REQ-693 (worktree status and cleanup survive a missing worktree folder) edits `internal/corehelpers/handoff.go` in the SAME Go package, so the two branches compile together: add no package-level identifier (see Anti-bloat). REQ-695 (`run-status` fixes) edits `internal/runstatus/`; REQ-696 (stale-wording sweep) edits prose and `_dev/tests/staged-skills-contract.sh`. Neither touches your files.

## Proof to run and record (from the REQ's Red-Green Proof)
Module dir below is `skills/do-work/tools/do-work-cli` of your worktree.
1. RED: write both tests first, before any production edit. Run `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestRequestAppendSectionRefusesBodyThatHidesLaterSections|TestFrontmatterSetRefusesLifecycleOwnedFields)$' ./internal/corehelpers/`. Expect rows 1 and 2 of the N1 test to fail (exit 0 path: the outcome is success and `## Review` is hidden) and both F3 cases to fail (`status: bogus-status` written). Row 3 passes at base (the old single-copy check refuses it); say so. Record the failing assertion text.
2. GREEN: implement, rerun the same command: every row passes, no `--- SKIP`.
3. Mutation (Constraint "each must fail with the new guard removed"): delete the count check and run the N1 test (all three rows must fail, including row 3, because the old check is gone too); restore it. Delete the F3 guard and run the F3 test (both cases fail); restore it. Record each failure line. `git diff` must show the guards back before you commit.

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-192201/REQ-694-probe.sh` run from your worktree root (it reads paths relative to `git rev-parse --show-toplevel`, so it checks your tree; it runs the two new tests plus `TestRequestAppendSectionLandsInCanonicalOrder`, gofmt and go vet on `internal/corehelpers`): exit 0. At base it exits 1 (the new tests' `--- PASS` lines are missing); that is expected.
- `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^(TestFrontmatterSet|TestRequestAppendSection|TestRequestIDResolves)' ./internal/corehelpers/` passes (the existing writer tests still hold).
- `gofmt -l skills/do-work/tools/do-work-cli/internal/corehelpers` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/corehelpers/` is clean.
- `git diff <base> --stat` shows the two files only; `git diff --check` is clean.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path, base commit (`git -C <worktree> rev-parse --short HEAD` before your first edit) and every commit hash (one commit preferred).
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: the RED run (failing assertion text, and row 3 passing at base), the GREEN run (test names, no skips), and the two mutation runs (failure lines).
- `## Decisions` (D-06 onwards; each DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellites and families).
- Anti-bloat check (see above): the `git diff --stat` and the list of anything added that the REQ did not name.
- A proposed CHANGELOG entry: a descriptive title that says what shipped, in plain words, one or two sentences on why it matters, then specific bullets (name the new `FRONTMATTER-FIELD-OWNED` code and the widened `SECTION-WRITE-FAILED` refusal). The integrator writes it with the version.
- A proposed lesson bullet, naming the satellite you think it belongs in, in that file's bullet shape (`- [family: <slug>] [REQ-694: <one-line lesson>](<relative archive link; the integrator fixes the path>)`). Say "none" if nothing was learned that a later builder would need.
- Integration seams (see the note above) and test wall times.
