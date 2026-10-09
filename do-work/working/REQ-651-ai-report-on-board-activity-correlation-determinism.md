---
id: REQ-651
title: 'AI report on the board activity correlation separates deterministic git structure from scaffold around agent behaviour'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer review "Activity Correlation … hundreds of lines of code to analyze Git history, merge branches, and track ancestry just to put a date on a dashboard card", triaged by `do-work-toolbox validate-feedback` on 2026-10-09 as F7; maintainer answer "create an ai-report and explain this to me, my concerns are that I don't want to build scaffold around not-deterministic behaviour".*
