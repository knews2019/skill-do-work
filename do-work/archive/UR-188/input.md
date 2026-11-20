---
id: UR-188
title: 'Repair second-opinion cleanup and outside-change checks'
created_at: 2026-11-20T01:30:31Z
requests: [REQ-755, REQ-756]
word_count: 9
---
## Summary

The user explicitly authorized capture and execution of both reproduced review findings from REQ-754 (second-opinion worktree isolation), plus recording the retrospective. REQ-755 repairs the new read-only-directory cleanup regression. REQ-756 repairs inherited outside-change detection gaps for committed edits and ignored new files. Implement them serially, cleanup first. This ordering is the session implementation choice, not a new functional prerequisite.

The later review ran against 4e48c2ef (version 0.331.2). Existing second-opinion tests passed in 38.532 seconds. Three added probes failed. Comparison with 02f3a5e1 confirmed the cleanup regression is new and both detection gaps predate the worktree change. The prior archived UR remains unchanged.

## Evidence

- `assets/REQ-755-756-review-probes.go.txt` preserves the original runnable Go overlay probes.
- `assets/REQ-755-756-review-evidence.txt` preserves observed current/baseline results.
- The retrospective is appended as post-work review metadata to `do-work/archive/UR-187/REQ-754-second-opinion-throwaway-worktree.md`.

## Full Verbatim Input
> ```
> Capture both fixes, record the retrospective, then run them.
> ```
