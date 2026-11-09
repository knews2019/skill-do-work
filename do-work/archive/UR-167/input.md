---
id: UR-167
title: 'Move the dated changelog archives to docs/changelog/'
created_at: 2026-11-09T21:47:36Z
requests: [REQ-715]
word_count: 122
---
## Full Verbatim Input
> ```
> do-work capture-request: Move the five dated changelog archives CHANGELOG-2026-*-up-to-v*.md from the repository root to docs/changelog/ (layout rule 1: the root CHANGELOG.md covers the whole repository; old changelog archives go to docs/changelog/). Depends on the REQ that moves kb/ and decisions/ under docs/. Use git mv. Rewrite the 15 live references: the GitHub blob links in skills/do-work/CHANGELOG.md and the root CHANGELOG.md, the /CHANGELOG-20*.md export-ignore glob and its comment in .gitattributes (it becomes /docs/changelog), and any release or archive-cut code or prime that names the root location (check _dev/primes/prime-releases.md), so the next archive is written to docs/changelog/. Verify: folder-structure-check no longer reports the archives, git archive HEAD does not contain docs/changelog/, the release --dry-run still parses the live CHANGELOG.md, and the contract suites pass.
> ```
