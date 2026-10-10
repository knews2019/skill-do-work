## Review: REQ-693

**Approve** — a missing builder worktree folder no longer breaks `worktree status` or `worktree cleanup`, and all three dirt readers ignore the owned `do-work/worktree-links` symlinks. One stale line in the by-hand cleanup recipe should be fixed now.
Route B | range 54ecd108..03f9beda (builder 22d7773e, addendum df527501, merge 03f9beda)

### What's built
- `worktree status` lists a builder worktree whose folder is gone as `missing` and exits 0. `worktree cleanup REQ-N` for that worktree runs `git worktree prune` before `git branch -d`.
- Cleanup Pass 5, the handoff survey and the board's verify probe ignore owned links. Pass 5 deletes those links (the links only) before its non-forced `git worktree remove`. An unreadable or invalid link config filters nothing, so the readers fail toward consent.
- Coordinator addendum D-13: the tutorial and the roadmap guide send "stuck" to `do-work status` and "broken or failed" to `do-work forensics`.

### Decisions / risks for you
- None. All of D-01 to D-13 are recorded and match the diff.

### Findings

**Important:**
- None.

**Minor:**
- F1: `skills/do-work/actions/fan-out-reference.md:101` (the by-hand "Cleanup — happy path") still runs `git worktree remove <path>` directly and says "`worktree remove` refuses on uncommitted builder work". When `do-work/worktree-links` exists, the non-forced remove also refuses on the owned links. The M1 mutation below shows git's refusal ("contains modified or untracked files"). An integrator who follows the hand path would then report "uncommitted builder work" that is not there and stop. Fix now (same file is already in this REQ's write set): before `git worktree remove <path>`, add "first delete the owned `do-work/worktree-links` symlinks in the worktree (links only, never their targets)", or say that those links do not count as builder work. — impact-user-visible → fix now (report only; no queue entry made)

**Nit:**
- F2: `cleanup_git.go:234-236` (`linkPaths = nil` on a config error) and `worktree_lifecycle.go:511-513` (the early `return nil` in `OwnedWorktreeLinks`) are redundant. `readWorktreeLinkConfig` already returns nil `linkPaths` on both failure branches. They still document fail-toward-consent, so keeping them is fine. — impact-negligible → report only
- F3: the two copies of the porcelain parse differ. do-work-cli `worktreeDirtyPaths` reads through `splitZero`, which sorts the NUL fields (REQ-660 review F7, left untouched on purpose), so the `index++` rename skip can pass over the wrong field. The board copy keeps git's order. The "change both together" comments cover the link config and the ownership check, not the parse. The dirty/clean answer agrees: a rename entry is always counted dirty before any skip, and the scratch test with an owned link plus a staged rename read dirty in both copies. Only the listed paths in `WORKTREE-DIRTY` can be wrong. — impact-negligible → report only
- F4 (the builder's discovered task, rated): Pass 5 cannot settle a `worktree-agent-REQ-41-*` worktree. Rating: not impact-critical. Report only, at impact-negligible rather than the builder's impact-user-visible. do-work only creates zero-padded ids (`formatRequestID` → `REQ-%03d`), so a real worktree is named `worktree-agent-REQ-041-*` and `requestIDFromWorktree` returns `REQ-041`, which matches `FilenameID`. `worktree new REQ-41` cannot even resolve a canonical `REQ-041` file (`RequestsByID` is keyed `REQ-041`). Only a REQ whose own id is the non-canonical `REQ-41` reaches this path. The model already treats that file as malformed, and consent is the safe answer for it. — impact-negligible → report only

### Requirements Checklist

- [x] Missing folder: `worktree status` exits 0 and lists the row marked missing — delivered (`reportStatus` stat, `WorktreeStatusRow.Missing`, text renders `missing`)
- [x] Missing folder: `worktree cleanup` skips the dirt check and `worktree remove`, and runs prune then `branch -d` — delivered
- [x] Pass 5 `worktreeClean` reads clean for a worktree whose only untracked paths are owned links, and Pass 5 then removes it — delivered (link removal before the non-forced remove, D-06)
- [x] Handoff survey ignores owned links — delivered (`cleanup.OwnedWorktreeLinks`, live fixture below)
- [x] queue-kanban `worktreeHasUncommittedWork` uses the identical rule (same file, same validation expression, same symlink-target check) — delivered. I compared the two copies line by line and found no drift.
- [x] Fail toward consent: an unreadable or invalid config filters nothing, and nothing is deleted — delivered (scratch tests below)
- [x] One missing-folder test and one owned-link test, no per-reader tests — delivered
- [x] Crash Recovery prose is unchanged. Correct: `work-reference.md:308-313` does not restate the dirt rule.
- [x] Addendum D-13: the tutorial and roadmap-guide routing matches `forensics.md:10`, `README.md:176` and `SKILL.md:40` — delivered
- [x] P-A-U boxes all checked. `## Scope` declares ten files and the range touches exactly those ten.

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `go test -count=1` on `internal/cleanup`, `internal/corehelpers`, `internal/resultmodel` at HEAD `03f9beda`: ok. queue-kanban `-run 'Worktree|Verify'`: ok.
- Mutations on a scratch copy, run against the two new tests. All four were caught:
  - M1: no Pass 5 link removal → `WORKTREE-REMOVE-FAILED` "contains modified or untracked files".
  - M2: `branch -d` before prune for a missing folder → "cannot delete branch … used by worktree".
  - M3: `worktreeClean` given no owned links → `WORKTREE-REQUIRES-CONSENT dirty=true`.
  - M4: status reads dirt in a missing folder → `WORKTREE-GIT-FAILED`.
- Scratch tests, all pass:
  - An invalid config (`../escape`) and an unreadable config (a directory) leave Pass 5 at `WORKTREE-REQUIRES-CONSENT` with the link still in place.
  - A link the user re-pointed elsewhere is dirt, and its target is untouched.
  - A dry run plans the removal and deletes no link.
  - An owned link plus a staged rename reads dirty.
  - A missing folder whose branch is not merged still refuses with `WORKTREE-NOT-MERGED`, and its status row is `missing` with `ahead 1`.
- Live CLI fixture (binary built in scratch). The link config named `shared-deps` and the non-ASCII `nöde`. `git status` quotes the second name (`"n\303\266de"`), so this checks the `-z` parse. `worktree new REQ-141`, then `handoff-state-survey`, gave `HANDOFF-WORKTREE-CLEAN` for both checkouts, and `worktree status` showed `clean`. With an invalid config line, the survey gave `HANDOFF-WORKTREE-DIRTY` for `nöde` and `shared-deps`.
- Board scratch test on the same fixture:
  - Valid config → clean.
  - Unknown main tree, invalid config and absent config → dirty.
  - Staged rename → dirty.
  - Rename undone → clean.

### Suggested Additional Testing

- Deployment: unassessed. After release, run `do-work-cli worktree status` in a consumer repo with a `do-work/worktree-links` file and one builder folder deleted by hand.
- Live acceptance: unassessed. Run a real `do-work cleanup` Pass 5 on a merged, settled builder worktree that has a linked `node_modules` (ignore line `node_modules/`), and confirm it is removed without asking.

### Scores (on the record — not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All seven criteria and the addendum delivered |
| Code Quality | 95% | Small, follows `ownedLinks`/`dirtyPaths` lineage; two redundant guards (F2) |
| Test Adequacy | 90% | Both RED cases pinned and every core mutation caught. Fail-toward-consent, handoff and board have no lock-in test (REQ constraint) |
| Scope | 100% | Ten declared files, ten touched |
| Risk | Low | Link deletion is limited to verified owned symlinks; readers fail toward dirt |
| Acceptance | Pass | Implementation and integration stages |

### Follow-ups created
- None (4 findings report only)

## Review

**Overall: 96%** | 2026-10-10T21:09:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `skills/do-work/actions/fan-out-reference.md:101`: the by-hand cleanup runs `git worktree remove <path>` without first deleting the owned `do-work/worktree-links` symlinks and says a refusal means uncommitted builder work, which is stale now that owned links are not dirt. Fix now: add "first delete the owned `do-work/worktree-links` symlinks (links only)" before that step — impact-user-visible → report only; F2 (Nit) redundant nil guards at `cleanup_git.go:234-236` and `worktree_lifecycle.go:511-513` — impact-negligible → report only; F3 (Nit) the do-work-cli parse (sorted `splitZero`, REQ-660 F7) and the board parse (git order) differ, but the clean/dirty answer agrees (tested). Only the listed dirty paths can differ — impact-negligible → report only; F4 (Nit, the builder's discovered task) Pass 5 cannot settle a `worktree-agent-REQ-41-*` worktree. Real ids are zero-padded (`REQ-%03d`), so only a non-canonical REQ id reaches this, and consent is safe there. Not critical — impact-negligible → report only
**Acceptance:** Pass. Implementation and integration stages: package tests green; 4 of 4 mutations caught; scratch tests for fail-toward-consent, dry run, foreign link, rename and missing-plus-unmerged; live handoff and board fixture with a non-ASCII link.
**Restatement sweep:** redefined the worktree dirt rule (owned links are not dirt; Pass 5 removes them before the non-forced remove), the `missing` worktree status, and Pass 5 link removal. Stale: `fan-out-reference.md:101` (F1). Consistent: `cleanup.md:140`, `fan-out-reference.md:53,105`, `cleanup-guide.md:29`, `restart-with-parallel-handoff.md:40`, `prime-do-kanban.md:34`, `work-reference.md:308-313`. Wave-end inherited elements, all consistent: REQ-694 `SECTION-WRITE-FAILED`/`FRONTMATTER-FIELD-OWNED` (`work.md:21,202,289,310`, usage strings in `request_writers.go:30,79`); REQ-695 C3 and the `run:` line (`status.md:48,67`, `status-guide.md:26`, `fan-out-reference.md:163`); REQ-696 fence rule (`clarify.md:106`, `capture-reference.md:196`) and "stuck" routing (`SKILL.md:40`, `forensics.md:10`, `README.md:176`, `tutorial.md:260-264`, `roadmap-guide.md:5,48`, `roadmap.md:5,124,278`, `stray-check.md:16`, `stray-check-guide.md:5,59`). REQ-696 F1 and F2 are now closed by D-13.
**Anti-bloat:** 0 unnamed additions. The four package-level helpers (`readWorktreeLinkConfig`, `ownedWorktreeLinks`, exported `OwnedWorktreeLinks`, `worktreeDirtyPaths`) are D-01's named lift, and two of them replace deleted methods. The board's `ownedWorktreeLinkPaths` is D-03. `Missing` is D-05. The new `mainWorktreePath` parameter and field are the plumbing D-02 needs. 2 tests, both named by the REQ. No new options, files or finding codes. Only the 2 redundant guards (F2) could go.
**Suggested testing:** 2 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

### Delta re-review (03f9beda..8211de3c)

**Approve** — F1 is closed. | 2026-10-10T21:11:03Z

- Range: `dd0f2c31`, merged as `8211de3c`. It changes one sentence in `skills/do-work/actions/fan-out-reference.md:101` and nothing else.
- F1: closed. The by-hand sequence now deletes the owned `do-work/worktree-links` symlinks (links only, never their targets) before `git worktree remove <path>`. So "`worktree remove` refuses on uncommitted builder work" is true again for linked worktrees. It now matches `cleanup.md:140` (Pass 5), `fan-out-reference.md:53` and the closing CLI sentence on the same line.
- Rest of the paragraph and the file: consistent. Line 53 defines the link file, line 105 (crash path) does not restate the dirt rule, and the "Never `--force`" and "Report the refusal and stop" rules are unchanged.
- N1 (Nit, new): the by-hand order deletes the links before any check. `worktree cleanup` first refuses an unmerged branch or a dirty worktree, and leaves the links in place when it refuses. So on a refusal, the by-hand path has already removed the links. This is harmless at Step 9 because finalization has finished and the links can be recreated from the main tree. — impact-negligible → report only
- Updated Overall: **96%**. The dimension scores do not change: Requirements 100%, Code Quality 95%, Test Adequacy 90%, Scope 100%, Risk Low, Acceptance Pass. Remaining findings: F2, F3, F4, N1, all report only.
