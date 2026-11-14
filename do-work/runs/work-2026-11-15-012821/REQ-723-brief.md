# Builder Brief — REQ-723 (the archive-collision helper skips assets folders)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-723-archive-collision-skips-assets`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-723-archive-collision-skips-assets`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its Red-Green Proof and decisions D-01 and D-02: read `do-work/working/REQ-723-archive-collision-helper-skips-assets-folders.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-723-handback.md`
- Route A (simple): no exploration, no scope section, no pre-flight. The REQ names both files and the rule to copy.
- Up to ten builders run at the same time in this run. Machine load is high; see the time-budget rule below.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. The prime: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. Required lessons: none fit the budget. Dropped only for budget, read this one bullet (line 5 of `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `[family: guard-on-one-entry-path]`, from REQ-722): "a skip rule that lives in one walker is not a rule for the tree … A walker that picks do-work records by file name must apply the same folder skips as `repositorymodel` … Pin the skip with a fixture that also holds a real blanked record beside the skipped folder, so a skip that grows too wide fails too."

## Decided (do not reopen)
- **D-01:** copy the skip rule into the walk; do not switch the command to the `repositorymodel` snapshot.
- **D-02:** skip only folders named `assets`. Leave dot-folders and the rest of the walk as they are.

## What to build
Two commits.

**Commit 1 (RED, test only).** Add `TestArchiveCollisionSkipsAssetsFolders` as a new function at the end of `skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go`. Use the existing helpers `writeMatrixFile` and `testContext`. The test pins both directions:
1. A fixture whose archive holds only `do-work/archive/UR-001/assets/REQ-001-evidence.md`: `handleArchiveCollision(testContext(root), []string{"REQ-001"})` returns `OutcomeSuccess` with empty `ExactTextOutput`.
2. Then add the real record `do-work/archive/UR-001/REQ-001-done.md` beside the assets folder: the result is `OutcomeFindings` with output exactly `do-work/archive/UR-001/REQ-001-done.md\n` (the evidence file is still not listed). This half is what fails if the skip grows too wide.
Run it and record the failure (step 1 lists the evidence path).

**Commit 2 (the fix).** In `handleArchiveCollision` (`skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:158-193`), inside the `filepath.WalkDir` callback after the archive-root checks, return `fs.SkipDir` when `entry.IsDir() && entry.Name() == "assets"`. Add a one-line comment in the same words as `cleanup/cleanup_git.go:115` ("assets/ folders hold evidence files, never records; repositorymodel skips them the same way."). Keep the missing-archive, not-a-directory and symlink behavior unchanged; both existing tests must still pass unmodified.

**Not yours:** `VERSION`, both `CHANGELOG.md` files, the version line in `skills/do-work/actions/version.md` (the integrator writes the release); anything under `do-work/`; `repositorymodel` and `cleanup` packages; `skills/do-work/tools/checks/archive-collision.sh`.

## Rules
- Write only the two files above. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-723]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/corehelpers` prints nothing; `go vet ./internal/corehelpers/`; `go test -count=1 ./internal/corehelpers/`.
2. End to end: in a `mktemp -d` Git repository holding only `do-work/archive/UR-001/assets/REQ-001-evidence.md`, run `bash <worktree>/skills/do-work/tools/do-work-cli.sh --repo-root <that repo> --format text archive-collision REQ-001` and expect exit 0 with no output. Record the exit status.
3. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-723-probe.sh` run with your worktree as the current directory prints `REQ-723 probe PASS`. It fails at base by design (the new test does not exist; the pre-dispatch agent confirmed this at main HEAD).

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created, modified or deleted, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the new test failing after commit 1 (the output line that lists the evidence path) and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-03 onward (D-01 and D-02 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — see below; name any shared file you touched anyway.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Archive-Collision Check Ignores Evidence in assets/ Folders") and the line "What changes for an installed copy": `do-work-cli archive-collision REQ-NNN` (and `tools/checks/archive-collision.sh`) no longer reports evidence files inside `do-work/archive/**/assets/` as archived twins; real archived records are still reported.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
No sibling REQ in this run edits `internal/corehelpers/commands.go` or `archive_collision_test.go`. The other nine REQs touch report-judge Go files, `ai-report` docs, crew-member files, decision records, the staged-skills contract and one `ai-reports/` bundle. Expect a clean merge; if `commands.go` conflicts anyway, the conflict is someone else's edit to a different handler, so keep both.

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
