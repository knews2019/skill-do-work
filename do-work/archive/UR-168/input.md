---
id: UR-168
title: 'Move the root installer scripts from tools/ to cmd/'
created_at: 2026-11-09T21:47:36Z
requests: [REQ-716]
word_count: 128
---
## Full Verbatim Input
> ```
> do-work capture-request: Move the five root installer scripts in tools/ (install-do-work-suite.sh, fetch-upstream-archive.sh, do-work-cli-preamble.sh, replace-text-section.sh, validate-suite-manifest.sh) to cmd/ (layout rule 5: scripts live in cmd/; the maintainer chose cmd over tools). Use git mv. First map every caller: the README bootstrap (source/tools/install-do-work-suite.sh), the Go update path, the shipped skills (80 files mention tools/; only the 34 that name these five scripts matter), and the contract tests. Rule out a break for already-installed consumers: if an older installed updater runs tools/<script> from inside the newly downloaded archive, keep a thin tools/ forwarder (or prove no installed version does) and say which in the REQ. Verify: folder-structure-check no longer reports tools/, the README bootstrap runs in a temp repo, an update from the previous release works, and the contract suites pass.
> ```
