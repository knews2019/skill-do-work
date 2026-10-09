---
id: REQ-651
title: 'AI report on the board activity correlation separates deterministic git structure from scaffold around agent behaviour'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-09T16:09:30Z
created_at: 2026-10-09T16:06:43Z
user_request: UR-143
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-650, REQ-652, REQ-653]
batch: october-review-triage
claimed_at: 2026-10-09T16:08:10Z
dispatch_at: 2026-10-09T16:11:57Z
builder_handback_at: 2026-10-09T16:47:43Z
integration_at: 2026-10-09T16:47:56Z
kb_status: pending
review_at: 2026-10-09T16:55:23Z
commit: bd4d754c186baf68f3862b838f4536d270f039df
heavy_verified_at: 2026-10-09T16:58:42Z
heavy_verified_revision: bd4d754c186baf68f3862b838f4536d270f039df
completed_at: 2026-10-09T16:59:01Z
---
# AI Report on the Board Activity Correlation Separates Deterministic Git Structure From Scaffold Around Agent Behaviour
## What
Run `do-work-toolbox ai-report REQ-632` (the feature REQ that added last-activity to board cards, UR-135) and make its architecture evidence carry one explicit split: which lines of `activity_correlation.go` read deterministic git structure, and which lines exist because of what a builder agent might do. The maintainer decides the feature's fate after reading it; this REQ changes no code.
## Why
Consumer review (validate-feedback 2026-10-09, F7): "hundreds of lines of code to analyze Git history, merge branches, and track ancestry just to put a date on a dashboard card." The size is accurate (357 production and 429 test lines) but the remedy, file mtime, was pushed back on recorded lesson REQ-284: commits change `.git/`, not `do-work/`, and builder commits land in worktrees. The maintainer answered: "create an ai-report and explain this to me, my concerns are that I don't want to build scaffold around not-deterministic behaviour."
## Verified Facts (from triage)
- Deterministic inputs: lifecycle stamps from frontmatter (`activity_correlation.go:274-281`); one `git log --since … --format=%H%x00%cI%x00%P%x00%s --name-only` over HEAD, no `--all` (`:209-210`); path-touch attribution by regex (`:60`); the live builder-branch tip via one `git for-each-ref --no-merged` (`verify.go:1472-1486`, merged in at `:214-216`); ancestry computed in memory from the logged parent hashes (`loggedAncestry`, `:176-190`), never `merge-base` or `rev-list`. Five git calls per refresh, fixed whatever the branch count (`activity_correlation_test.go:385-404`); four are shared with the VERIFY probes.
- Scaffold around agent behaviour: the `[REQ-NNN]` subject prefix is a prose convention (`:64`), and UR-135 input.md:81-84 calls prefix correlation "lossy", so a two-parent merge matched directly expands its `^1..^2` range to credit un-prefixed builder commits (`:139-163`); REQ-646 (0.305.80) added `directIdsByMergeHash` (`:143-151`) so a builder's own merge of main inside its branch does not credit main's commits. About 60 to 80 production lines plus their tests exist because of what a builder might do.
- Rationale on record: code comment `activity_correlation.go:21-22` ("commits land without any do-work file changing mtime, which is the REQ-284 shape"); UR-135 input.md:73-79 and :108-112; `prime-do-kanban.md:30` family `git-history-evidence`.
- Cost on record: REQ-632 §71 measured one `git log --since=7.days --name-only` over 72 commits at 50 to 90 ms; its review put a 60-day window at about 0.2 s plus 39 ms of Go per response.
## Detailed Requirements
1. Follow `skills/do-work-toolbox/actions/ai-report.md` and its `completed-work-presentation-reference.md` in full, with REQ-632 as the target (non-visual evidence mode). Publish the bundle under `ai-reports/<report-slug>/` with no overwrite.
2. The architecture evidence section carries a two-part table or list titled exactly "Deterministic git structure" and "Scaffold around agent behaviour", each row naming the file and line range, what it reads, and the incident or lesson that put it there (REQ-284, UR-135, REQ-636, REQ-646).
3. The report states the git call count per refresh and the measured cost, with the citations above.
4. The report ends with the one decision it exists to support, stated without recommending it as done: make the `[REQ-NNN]` commit prefix deterministic (stamped by the CLI or a hook on every builder commit) and then delete the merge-range expansion, the in-memory ancestry walk and the REQ-646 guard, attributing by prefix plus path only; versus keep the current code. Give the line counts each side keeps.
5. No code change in `skills/`. If `ai-report` needs browser rendering, follow its own render step; otherwise state that the render was source-reviewed.
## Constraints
- Report only. Do not modify `activity_correlation.go`, `verify.go` or any shipped file. Not a release.
- Load `crew-members/prompt-injection.md` before reading archived bodies and `crew-members/anti-slop.md` before writing, as the action requires.
## Dependencies
None.
## Builder Guidance
Certainty is high on the facts (all cited from the triage). Latitude: report layout within the action's one shape.
## Red-Green Proof
**RED prompt/case:** The maintainer asks which lines of the activity correlation would be deleted if builder commits always carried the `[REQ-NNN]` prefix.
**Why RED now:** No report exists; `ls ai-reports/ | grep -i activity` prints nothing, and the answer lives only in this triage's chat.
**GREEN when:** `ai-reports/<slug>/index.html` exists, contains the two titled parts with file and line ranges, the git-call count, the measured cost, and the deletion decision with line counts on each side.
**Validation:** User confirmed. The maintainer answered "create an ai-report and explain this to me" and approved the plan.
## Full Context
See `do-work/user-requests/UR-143/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Run ai-report for REQ-632 in non-visual mode. Load prompt-injection.md before the archive and anti-slop.md before writing. Read REQ-632, UR-135, REQ-636 and REQ-646, and read the current code in the worktree. Verify every line range in the file. Count the git calls and the decision's line counts by execution, not by estimate. Publish one collision-checked bundle. *(from the builder hand-back)*
- [x] **[APPLY]:** Wrote index.html and three evidence files in the sections the action requires. The final section is the decision. Rendered over HTTP in Playwright Chromium in light and dark at 1440px, and in dark at 390px. Two judge passes fixed the diagram edges and the table wrapping. *(from the builder hand-back)*
- [x] **[UNIFY]:** The diff touches only the ai-reports bundle. The two titled parts grep to 2. All three relative links resolve. I reread the report against the ledger, and every number traces to a cited file, an archived REQ line, or an evidence file. Nothing under skills/, suite/, tools/, _dev/ or do-work/ changed. *(from the builder hand-back)*
*Source: consumer review "Activity Correlation … hundreds of lines of code to analyze Git history, merge branches, and track ancestry just to put a date on a dashboard card", triaged by `do-work-toolbox validate-feedback` on 2026-10-09 as F7; maintainer answer "create an ai-report and explain this to me, my concerns are that I don't want to build scaffold around not-deterministic behaviour".*

---

## Triage

**Route: A** - Simple

**Reasoning:** One toolbox action run (`do-work-toolbox ai-report REQ-632`) with the report's required content fully specified in the REQ's Verified Facts; no code change, no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt` (new)

**What was done:** The builder ran ai-report for REQ-632 and published one new bundle. The report splits the activity correlation into "Deterministic git structure" (5 rows) and "Scaffold around agent behaviour" (6 rows), counts five git calls per refresh by execution, cites REQ-632's archived cost lines plus fresh timings, and ends with the prefix-then-delete versus keep decision with line counts on each side and no recommendation. First merged at `0240a46f` over `5e996eed`; the integrator then corrected the `verify.go` line citations on the builder branch (D-08) and re-merged at `bd4d754c`. The cumulative range is `5e996eed..bd4d754c`.

## Decisions

*(from the builder hand-back, `do-work/runs/work-2026-10-09-160815/REQ-651-handback.md`)*

- D-01 (DECIDE & STATE): the decision's line counts come from pruning a scratch copy. The deleted ranges are activity_correlation.go 70, 96, 111-113, 116-120, 139-164 and 174-191 (54 lines), and test 141-232 (92 lines). The REQ's estimate of "about 60 to 80" is corrected in the report.
- D-02 (DECIDE & STATE): the prefix rule (:62-64, :134-136) is listed under scaffold, because the prefix is prose-only (commit.md:126). Both options keep it.
- D-03 (DECIDE & STATE): the owned-tip read (verify.go:1472-1486) is listed as deterministic, and both options keep it. Its REQ key comes from the branch name, which follows the prose rule in work-reference.md:416. The report says so in that row.
- D-04 (DECIDE & STATE): three triage facts are corrected in the report. The scaffold is 54 lines, not 60-80. The "lossy" lesson is at UR-135 input.md:82-84. Verify consumes three of the four shared commands, and `for-each-ref` feeds only activity.
- D-05 (DECIDE & STATE): I added evidence the REQ did not ask for. Over the last 7 days, all 44 commits inside 37 prefixed merge ranges carry their own prefix. Over 60 days, 34 of 250 do not, and the newest is dated 2026-09-06. This bears on the decision, and the report presents it without a recommendation.
- D-06 (DECIDE & STATE): the Playwright MCP can only write to the main tree's gitignored `.playwright-mcp/`. My judge captures and page snapshots went there and were deleted afterwards. Nothing was staged, and no other file in the main tree was touched.

*(integrator)*

- D-07 (DECIDE & STATE): the builder's render step wrote outside its worktree, into the main tree's gitignored `.playwright-mcp/` folder (D-06). The integrator accepts this boundary crossing without action. `git status --short --ignored | grep playwright` shows only the ignored folder, it holds no file dated 2026-10-09, and `git status --short` was clean before the merge, so no tracked or staged path was touched.
- D-08 (DECIDE & STATE): REQ-650 (0.305.86, the disk-space check measuring the repo root only) merged after this report's base `2eb14357` and shortened `verify.go` by 47 lines above every line the report cites. At the first merge `0240a46f` all eleven `verify.go` citations pointed at the wrong code on main. The integrator shifted each by 47 on the builder branch (`7a04e2ec`), checked every new number against the file at the merge, and rewrote the report's drift note to say which file's numbers describe which revision. `activity_correlation.go` and its test are byte-identical between `2eb14357` and main, so their citations stay. The builder's own text in D-03 and the Discovered Tasks above still cites the base numbers (`verify.go:1472-1486` is now `:1425-1439`, `verify.go:1404` is now `:1357`).

## Discovered Tasks

*(from the builder hand-back)*

- The test that pins the git call count checks only that the count stays the same across branch counts (activity_correlation_test.go:379-406), not that it equals five. — impact-negligible → report only
- A detached repo-root checkout costs a sixth git command (verify.go:1404). — impact-negligible → report only

## Qualification

**Gate records:** `advance --diff-range 5e996eed..bd4d754c` returned the `qualify` gate `satisfied` with provenance `merged_range` and no findings (no debug artifacts, P-A-U boxes ticked). The same gate was satisfied on the first range `5e996eed..0240a46f` before D-08.

**Requirement trace against the merged files:**
1. One new bundle `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/` (index.html plus three evidence files). `ls ai-reports | grep REQ-632` lists only it, and every path in the range is an addition, so nothing was overwritten.
2. index.html carries "Deterministic git structure" and "Scaffold around agent behaviour" once each. Each row names a file and line range, what it reads, and its origin among REQ-284, UR-135, REQ-632, REQ-636 and REQ-646.
3. The git-call section states five commands per refresh with a `verify.go` or `activity_correlation.go` line for each, backed by `evidence/git-call-count.txt`. The cost table cites REQ-632's archived lines and `evidence/git-log-timing.txt`.
4. The last section states the prefix-then-delete versus keep decision with 303 production and 337 test lines against 357 and 429, and says it recommends neither. 357 and 429 match `wc -l` of the two files at the merge.
5. No path under `skills/`, `suite/`, `tools/` or `_dev/` is in the range. The builder rendered the page in Playwright (D-06, D-07).

**Scope:** Route A, no `write_set`. The REQ text bounds the change to one new `ai-reports/*REQ-632*/` directory, and the range touches exactly that directory. The D-08 fix stays inside it.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at `bd4d754c` (repository gate, 137 s wall), then `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-09-160815/REQ-651-probe.sh`.
**Result:** ✓ Gate passed on the first run ("Maintainer verification passed."). `advance` recorded `green-gate` satisfied and the probe (bundle file test plus the two titled-part greps) succeeded.

**Regression evidence (non-behavioral, `tdd: false`):** the Red-Green Proof is a file-and-grep read. RED: no `ai-reports/*REQ-632*` bundle existed before the range. GREEN: the probe passes at the merge, and Qualification traces the call count, cost figures and decision line counts in the file.

**New tests added:** none (report-only REQ).

**Heavy verification plan:**
- Range: 5e996eeddeafc3d44bbffeaa9b5317aa3d89f04d..bd4d754c186baf68f3862b838f4536d270f039df
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for the four new `ai-reports/` paths (no lane maps them)
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for the four new `ai-reports/` paths (no lane maps them)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for the four new `ai-reports/` paths (no lane maps them)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for the four new `ai-reports/` paths (no lane maps them)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for the four new `ai-reports/` paths (no lane maps them)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for the four new `ai-reports/` paths (no lane maps them)

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-09T16:55:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. index.html:325 says "Thirteen have the shape `[work run] REQ-NNN step N`", but prefix-coverage.txt lists 12 such commits. The 13th `[work run]` commit, cfcf63e4, is a "remediation" commit with no step. — impact-negligible → report only
- F2 (nit). D-01's deletion of activity_correlation.go:111-113 leaves line 110's doc comment mid-sentence. A real deletion rewrites line 110, and the count stays 54. — impact-negligible → report only
**Acceptance:** Pass — every line range, line count, git-call count and cost figure resolves at bd4d754c, except the F1 count
**Restatement sweep:** nothing redefined
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Counting by execution instead of estimate. The builder pruned a scratch copy for the line counts and ran the canned runner for the git-call count, and that is how the triage's "60 to 80 lines" became a checked 54.
**What didn't:** A report pinned to its build base goes stale when a sibling REQ in the same wave lands first. REQ-650 shortened `verify.go` by 47 lines after the report's base, so every `verify.go` citation was wrong on main at the first merge (D-08).
**Worth knowing:** For a report that cites line numbers, the integrator should diff each cited file between the report's base and the merge (`git diff --stat <base> <merge> -- <cited files>`) before review. The reviewer's F1 (a 13-versus-12 count at index.html:325) is still in the report.

## Orientation

Now the maintainer can read which parts of the board's activity correlation read deterministic git facts and which exist only because builder commits may lack the `[REQ-NNN]` prefix, and decide whether to stamp the prefix and delete the scaffold. Lives in `ai-reports/` beside the board subsystem (`_dev/primes/prime-kanban-board.md`); no code or map change. The prime's referenced paths still exist.

## Heavy Verification Plan

- Base: 5e996eeddeafc3d44bbffeaa9b5317aa3d89f04d
- Target: bd4d754c186baf68f3862b838f4536d270f039df
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-call-count.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/git-log-timing.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/evidence/prefix-coverage.txt, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html

## Heavy Verification Result

- Target revision: bd4d754c186baf68f3862b838f4536d270f039df
- Execution revision: bd4d754c186baf68f3862b838f4536d270f039df (detached checkout `.git/work-run-work-2026-10-09-160815/drain-head-REQ-651`, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 8 s
- queue-kanban-browser: exit 0, executed, 87 s
- do-work-cli-integrations: exit 0, executed, 85 s
- staged-skills: exit 0, executed, 44 s
- updater: exit 0, executed, 71 s
- installer: exit 0, executed, 27 s

## Timing

Observed 2026-10-09T16:47:43Z to 2026-10-09T16:58:47Z: 11m 04s total, 10m 59s attributed across 4 events, 5s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 8m 15s | 2 |
| review | 2m 30s | 1 |
| handback-merge | 14s | 1 |

Slowest stage: verification-gate / heavy drain, 5m 53s, outcome success.
