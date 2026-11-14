---
id: UR-173
title: 'Cleanup Pass 6 skips assets folders when scanning for blanked records'
created_at: 2026-11-14T21:47:51Z
requests: [REQ-722]
word_count: 61
---
## Summary
One accepted finding from a validate-feedback triage of pasted consumer feedback (2026-11-14). Verdict: Accept, no Discuss items. Captured as REQ-722 (Cleanup Pass 6 skips assets folders when scanning for blanked records).
## Full Verbatim Input
> ```
> do-work capture-request: cleanup Pass 6 (AddBlankedRepairs, cleanup_git.go:117) flags assets/REQ-*.md evidence files as BLANKED-RECORD-UNRECOVERABLE and exits 1. Source: pasted consumer feedback, 49 false warnings after archiving UR-610 and UR-639. Evidence: the walk does not skip assets/ folders, while repository_model.go:192 does; reproduced in a scratch repo. Surface-cost: N/A. Remedy: skip assets/ folders in the Pass 6 walk and add one test for it.
> ```