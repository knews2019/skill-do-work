---
id: UR-174
title: 'Archive-collision helper skips assets folders'
created_at: 2026-11-14T22:31:56Z
requests: [REQ-723]
word_count: 26
---
## Summary
The user promoted review finding F1 of REQ-722 (Cleanup Pass 6 skips assets folders when scanning for blanked records), which was report only, into queue work. Captured as REQ-723 (Archive-collision helper skips assets folders).
## Full Verbatim Input
> ```
> do-work capture-request: the standalone archive-collision helper (internal/corehelpers/commands.go:164-176) does not skip assets/ folders, the same gap REQ-722 fixed in cleanup Pass 6. Source: REQ-722 review finding F1.
> ```