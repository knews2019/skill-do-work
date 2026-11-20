---
id: UR-187
title: 'Run second-opinion backends in a throwaway Git worktree'
created_at: 2026-11-19T23:37:12Z
requests: [REQ-754]
word_count: 46
---
## Summary

Conversation context (pipeline-authored, not the user's words): `ai-report second-opinion` refuses to start when any file outside the report bundle has uncommitted changes, because it compares `git status` outside the bundle before and after the backend run. The user proposed running the backend in a new worktree, checking it, merging the result back, and always removing the worktree. Capture-time decisions, made through the ask tool:

- Isolation: Git worktree (chosen over a temporary copy of only the bundle, and over keeping the current check).
- Session mode (prepare/finalize, where the calling session is the editor): drop the outside-bundle Git status check; keep the bundle preservation check.
- Proof: both RED/GREEN cases recorded in REQ-754 were confirmed.
- Delivery: capture as a do-work REQ, not build in the capturing session.

## Full Verbatim Input
> ```
> If we are talking about this level of untrustworthiness, wouldn't then it make more sense to create a new work tree, work on that and then check it before merging it in? Also, make sure to clean up any work tree after the work is done.
> ```
