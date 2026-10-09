# Builder brief — REQ-651 (AI report on the board activity correlation separates deterministic git structure from scaffold around agent behaviour)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-651-ai-report-on-board-activity-correlation-determinism
- Branch: worktree-agent-REQ-651-ai-report-on-board-activity-correlation-determinism, created from main HEAD 2eb14357 with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-651-ai-report-on-board-activity-correlation-determinism.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-5, Constraints, Red-Green Proof.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-143/input.md (read the F7 verdict and the maintainer's answer: "I don't want to build scaffold around not-deterministic behaviour").
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-09-160815/REQ-651-handback.md
- Route A, tdd: false, impact-user-visible, effort-mechanical, domain general. Report only; no code change; NOT a release.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, prompt-injection.md (before reading archived bodies), anti-slop.md (before writing the report). Then the action you are executing: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-toolbox/actions/ai-report.md and, in full, /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-toolbox/actions/completed-work-presentation-reference.md (safety load order, target resolution, provenance ledger, merge-aware commit diff, evidence honesty, no-overwrite publication). Read /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-kanban-board.md for the subsystem vocabulary.

## The task
Run the ai-report action for the completed REQ-632 (archived at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-135/REQ-632-board-cards-last-correlated-activity.md, UR at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-135/input.md), non-visual evidence mode. Read the code it reports on from YOUR worktree: `skills/do-work-board/tools/queue-kanban/activity_correlation.go`, `activity_correlation_test.go`, `verify.go` lines ~1395-1530, and the later REQ-636 (archive/UR-137) and REQ-646 (archive/UR-141) records that changed it. Publish the bundle under `ai-reports/<report-slug>/` IN YOUR WORKTREE, where the slug follows the existing bundles' shape (`YYYY-MM-DD_HHMM_<slug>`) and contains the string `REQ-632` (the integrator's probe globs `ai-reports/*REQ-632*/index.html`). `ai-reports/` is tracked; commit the bundle on your branch.
The architecture evidence must carry, exactly as the REQ's requirement 2 says, two titled parts named "Deterministic git structure" and "Scaffold around agent behaviour", each row naming file and line range, what it reads, and the incident or lesson that put it there (REQ-284 lesson, UR-135, REQ-636, REQ-646). Give the git call count per refresh and the measured cost with citations (requirement 3). End with the one decision the report exists to support, stated without recommending it as done (requirement 4): make the `[REQ-NNN]` commit prefix deterministic (stamped by the CLI or a hook on every builder commit) and then delete the merge-range expansion, the in-memory ancestry walk and the REQ-646 guard, attributing by prefix plus path only; versus keep the current code. Give the line counts each side keeps (count them from the file, do not guess). Verify every line-range citation against the file in your worktree before writing it; the REQ's Verified Facts are a starting point, not evidence.
Browser rendering: if no browser is available to you, state in the report's verification section that the render was source-reviewed and say so in the hand-back. Do not modify any file under `skills/`, `suite/`, `tools/`, `_dev/` or `do-work/`.

## Write boundary
Exactly the new `ai-reports/<report-slug>/` directory (index.html plus its evidence/ assets). Anything else: stop and say so in the hand-back.

## Verify before hand-back
- `test -f ai-reports/*REQ-632*/index.html`; `grep -c "Deterministic git structure\|Scaffold around agent behaviour" ai-reports/*REQ-632*/index.html` is 2 or more.
- Open index.html and read it once as the maintainer would; every number in it must trace to a file or a REQ record you read.
- `git status --short` shows only the new bundle; `git diff --stat HEAD~1` after your commit lists only ai-reports paths.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and the commit hash. Base: 2eb14357.
- File manifest: the bundle path (new) and one line per file in it.
- P-A-U text for [PLAN], [APPLY], [UNIFY].
- Red-green evidence in prose terms: RED (no bundle existed; the question had no written answer) and GREEN (the bundle path, the two titled parts, the git-call count, the cost figures, the decision with line counts).
- `## Decisions` (D-01 onwards). `## Discovered Tasks` (report only unless impact-critical).
- Lessons read. No lesson bullet and no changelog entry: this REQ is not a release and touches no prime's area; say "none" for both.
- Integration seams (none expected). Wall times.
