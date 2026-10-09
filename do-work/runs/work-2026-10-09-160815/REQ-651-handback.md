# REQ-651 hand-back (AI report on the board activity correlation)

- Branch: worktree-agent-REQ-651-ai-report-on-board-activity-correlation-determinism
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-651-ai-report-on-board-activity-correlation-determinism
- Commit: 5ac41434bbf884ff16e32808cc5d3077147858f4 (one commit). Base: 2eb14357.

## File manifest
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/` (new)
  - `index.html`: the report.
  - `evidence/git-call-count.txt`: output and source of a throwaway test that counted git commands per served response, run in a scratch copy of the package.
  - `evidence/git-log-timing.txt`: windowed `git log` timings on this repo for 7 and 60 days.
  - `evidence/prefix-coverage.txt`: how many commits inside prefixed merge ranges lack both a bracketed token and a REQ path.

`git diff --stat HEAD~1` lists only these four files. `git status --short` is clean.

## P-A-U
- [PLAN]: Run ai-report for REQ-632 in non-visual mode. Load prompt-injection.md before the archive and anti-slop.md before writing. Read REQ-632, UR-135, REQ-636 and REQ-646, and read the current code in the worktree. Verify every line range in the file. Count the git calls and the decision's line counts by execution, not by estimate. Publish one collision-checked bundle.
- [APPLY]: Wrote index.html and three evidence files in the sections the action requires. The final section is the decision. Rendered over HTTP in Playwright Chromium in light and dark at 1440px, and in dark at 390px. Two judge passes fixed the diagram edges and the table wrapping.
- [UNIFY]: The diff touches only the ai-reports bundle. The two titled parts grep to 2. All three relative links resolve. I reread the report against the ledger, and every number traces to a cited file, an archived REQ line, or an evidence file. Nothing under skills/, suite/, tools/, _dev/ or do-work/ changed.

## Red-green evidence
- RED: no bundle existed. `ls ai-reports/ | grep -i REQ-632` printed nothing before this commit, and the question had no written answer.
- GREEN: the bundle is `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html`. It has two titled parts, "Deterministic git structure" with 5 rows and "Scaffold around agent behaviour" with 6 rows. Each row names the file and line range, what it reads, and its origin: the REQ-284 lesson, UR-135, REQ-632, REQ-636 or REQ-646.
  - Git calls: five per refresh, the same with 0, 1 or 4 builder branches, counted by execution. A detached checkout adds a sixth.
  - Cost: the report cites REQ-632's archived lines 71, 116, 224 and 256, and adds today's timings. A 7-day window took 44–47 ms warm and a 60-day window took 152–161 ms warm.
  - Decision: option 1 (stamp the prefix, then delete) keeps 303 production lines and 337 test lines, plus stamp code that is not written yet. Option 2 (keep) keeps 357 and 429. The report recommends neither.

## Decisions
- D-01 (DECIDE & STATE): the decision's line counts come from pruning a scratch copy. The deleted ranges are activity_correlation.go 70, 96, 111-113, 116-120, 139-164 and 174-191 (54 lines), and test 141-232 (92 lines). The REQ's estimate of "about 60 to 80" is corrected in the report.
- D-02 (DECIDE & STATE): the prefix rule (:62-64, :134-136) is listed under scaffold, because the prefix is prose-only (commit.md:126). Both options keep it.
- D-03 (DECIDE & STATE): the owned-tip read (verify.go:1472-1486) is listed as deterministic, and both options keep it. Its REQ key comes from the branch name, which follows the prose rule in work-reference.md:416. The report says so in that row.
- D-04 (DECIDE & STATE): three triage facts are corrected in the report. The scaffold is 54 lines, not 60-80. The "lossy" lesson is at UR-135 input.md:82-84. Verify consumes three of the four shared commands, and `for-each-ref` feeds only activity.
- D-05 (DECIDE & STATE): I added evidence the REQ did not ask for. Over the last 7 days, all 44 commits inside 37 prefixed merge ranges carry their own prefix. Over 60 days, 34 of 250 do not, and the newest is dated 2026-09-06. This bears on the decision, and the report presents it without a recommendation.
- D-06 (DECIDE & STATE): the Playwright MCP can only write to the main tree's gitignored `.playwright-mcp/`. My judge captures and page snapshots went there and were deleted afterwards. Nothing was staged, and no other file in the main tree was touched.

## Discovered Tasks
- The test that pins the git call count checks only that the count stays the same across branch counts (activity_correlation_test.go:379-406), not that it equals five. — impact-negligible → report only
- A detached repo-root checkout costs a sixth git command (verify.go:1404). — impact-negligible → report only

## Lessons read
prompt-injection.md, anti-slop.md, general.md, shared-principles.md, prime-kanban-board.md, and the REQ-284 entry in lessons-kanban-board.md:49. Lesson bullet: none. Changelog entry: none. This is not a release.

## Integration seams
None. The commit adds new files under ai-reports/ only.

## Wall times
Started 2026-10-09 19:12 EEST, handed back about 19:25 EEST, roughly 13 minutes. Go tests took 0.44 s. Each Playwright render pass took a few seconds.
