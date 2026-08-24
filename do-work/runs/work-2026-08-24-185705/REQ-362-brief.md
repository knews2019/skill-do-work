# REQ-362 exploration brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-362-stop-a-multi-path-bullet-disabling-the-scope-drift-check`

Read `_dev/primes/prime-shell-commands.md` completely and inspect only repository files needed to
answer the exploration questions below. Do not edit code yet. Do not read or write any `do-work/`
path in the worktree; report findings to the orchestrator in conversation.

Explore:

- the current Scope and Implementation Summary extraction functions in `scope-drift.sh`;
- a portable extraction path for every backticked repo-relative path on one bullet;
- a loud, non-vacuous treatment for filename-only backticks without breaking prose-only bullets;
- REQ-344's archived multi-path example and the existing scope-drift case corpus;
- exact mutation and before/after compatibility evidence.

Return a concise proposed implementation, exact declared file list, acceptance-to-test mapping,
risks, and any reason the captured one-file write set must expand. Wait for an explicit follow-up
before coding.
