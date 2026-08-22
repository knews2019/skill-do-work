---
session_ended: 2026-08-21T19:41:36Z
last_completed: REQ-317
queue_state: [0 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress]
reqs_processed_this_session: 7
session_depth: heavy
---

# Session Checkpoint

## Completed This Session

- REQ-309: Run the repo's canonical gate before hand-back (Route C, 79%)
- REQ-310: Check template payload citations at their destination (Route B, 98%)
- REQ-312: Resolve same-package citations in the shipped reference contract (Route C, 96%)
- REQ-313: Count the breaks the timeline actually draws (Route B, 99%)
- REQ-314: Judge effort_estimate on review-minted follow-ups too (Route B, 99%)
- REQ-316: Audit the calibration-log write step for the REQ-274 bug class (Route B, 100%)
- REQ-317: Preserve canonical-gate holds in error handling (Route B, 99%)

## In Progress (interrupted)

- REQ-319: List only the REQs the selected window covers — claimed 2026-08-22T23:23:42Z — writer: vm:/home/user/skill-do-work

## Still Queued

None.

## Session Notes

- Every implementation and staged release tree passed the direct canonical maintainer gate,
  including uncached Go, strict JavaScript, strict browser, and audit-metrics lanes.
- Releases advanced from 0.236.3 through 0.236.10. Completed UR-055, UR-057, UR-062, and UR-064
  were consolidated and their durable prime links repointed.
- Cleanup found no terminal queue/working files, no active UR closures, no loose archive REQs, no
  consumed run manifests, no orphaned worktree-agent branches/worktrees, and no blanked REQ/UR
  files. Eight historical run directories remain because their manifests are absent or do not carry
  `Status: consumed`; they were reported and left untouched.

## Context Summary

- Canonical verification is additive to focused tests and is judged only by its direct zero exit
  status. Current-diff failures enter remediation; unrelated or pre-existing canonical-gate failures
  preserve the claimed REQ and checkpoint and stop before archive/commit/hand-back. Re-read
  `_dev/primes/prime-action-files.md` and `_dev/primes/prime-shell-commands.md` before changing this
  policy.
- Shipped-reference validation now covers same-package citations across both source and installed
  topology. Consumer-owned examples declare `<project-root>` explicitly.
- Timeline break counts are a row-deduplicated union over filtered completion anomalies, reversed
  waits, and drawable reversed work spans.
- Every current REQ writer judges effort explicitly. Semantic instruction tests must isolate the
  directive itself; neighboring prose and fenced templates can otherwise mask a weakened rule.
- Calibration arithmetic reads `claimed_at` and `completed_at` from the just-archived REQ
  frontmatter at calculation time, never from values carried earlier in context.
