# REQ-704 integrator brief — run work-2026-10-20-200259

- REQ: REQ-704 (second opinions build a linked iteration tree on any base page)
- REQ path: do-work/working/REQ-704-second-opinion-linked-iteration-tree.md
- UR: UR-160 (do-work/user-requests/UR-160/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-20-200259/
- Shared guide (read it fully first): do-work/runs/work-2026-10-20-200259/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-10-20-200259/REQ-704-handback.md
- Operative name / branch: worktree-agent-REQ-704-second-opinion-linked-iteration-tree
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-704-second-opinion-linked-iteration-tree
- Dispatch instant (for the builder-work timing event): 2026-10-20T22:00:40Z
- Route: C (thorough review against the Plan)
- Wave membership: REQ-703, REQ-704 (serial). Set aside: none.
- Closes UR-160: YES (UR_CLOSES=1; archive do-work/archive/UR-160/REQ-704-second-opinion-linked-iteration-tree.md; REQ-703 and REQ-705 move from their flat archive paths into do-work/archive/UR-160/ in the same transaction: put both flat paths, both UR-160 destinations, and the UR input.md (do-work/user-requests/UR-160/input.md and its archive/UR-160 destination) in commit_paths)
- Last successful integration of the wave (run the Restatement Sweep): YES: run the Restatement Sweep; the other members are REQ-703 and REQ-705
- Lesson satellite for a Go lesson: skills/do-work/tools/do-work-cli/lessons-do-work-cli.md (action-file lessons: _dev/primes/lessons-action-files.md)
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- Pre-flight was recorded by the coordinator after REQ-705's release; `## Pre-Flight` is in the REQ. Never re-run pre-flight.
- tdd: true. Testing carries RED/GREEN evidence for the REQ-704 tests named in the hand-back.
- RUN/REQ-703-probe.sh and RUN/REQ-705-probe.sh are earlier members' probes; your probe is RUN/REQ-704-probe.sh.
- Link sweep before make-manifests (canonical-link-outlives-its-target): every shipped link to do-work/archive/REQ-703-*.md or REQ-705-*.md (for example skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, the REQ-703 bullet) must be re-pointed to do-work/archive/UR-160/. Run `git grep -n "archive/REQ-70[35]" -- ':!do-work/'` and fix every hit; then add those files and do-work/lessons-index.md (token refresh) to commit_paths.
- Release: re-read VERSION right before prep-payloads (0.306.7 expected after REQ-705); this REQ ships the next patch.
- Reviewer anti-bloat: ask it to count helpers, options, files or tests the brief did not name.
