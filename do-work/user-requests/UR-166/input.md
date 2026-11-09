---
id: UR-166
title: 'Suite writes its knowledge base and ADR log under docs/'
created_at: 2026-11-09T21:47:36Z
requests: [REQ-714]
word_count: 178
---
## Full Verbatim Input
> ```
> do-work capture-request: Make the do-work suite write its knowledge base and ADR log under docs/ (layout rule 7: all docs live under docs/). Change the BKB default from kb/ to docs/kb/ (the bkb-init/bkb-status/bkb-lint-structure defaults in the Go CLI and in skills/do-work-board/justfile.template, plus every action and prompt in skills/do-work-knowledge that names kb/), and change the ADR-log prompt (skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md and the actions that cite decisions/) to write docs/decisions/. Existing consumers keep working: a repo that already has kb/ or decisions/ and no docs/ copy must still be found (or told the one-line --kb override), not silently get a second empty scaffold. In this repo, git mv kb/ to docs/kb/ and decisions/ to docs/decisions/, update the .gitattributes export-ignore lines (/kb becomes /docs/kb, /decisions becomes /docs/decisions) and the contract-regressions.sh assertion that checks them, and add docs/kb/ and docs/decisions/ to the Target Layout tree in skills/do-work-toolbox/actions/folder-structure.md. Rewrite the 118 live references (not counting do-work/ history). Verify: folder-structure-check no longer reports kb/ or decisions/, a fresh bkb-init in a temp repo creates docs/kb/, git archive HEAD contains neither folder, and the contract suites pass.
> ```
