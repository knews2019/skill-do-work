## 0.305.9 — A Maintainer-Only Test No Longer Ships, and a Skewed Claim Keeps Its Progress Figures Ticking (2026-09-06)

Two fixes from an external review of the open pull request, applied directly.

- One test in the shipped `do-work-cli` module read this repository's own maintainer tree and skipped only when a directory six levels above the package was missing. In an installed copy that directory is `<project>/.claude/_dev/tests`; a project that happens to have it ran the test and failed on a manifest that never ships. The test now lives in the export-ignored file with the other maintainer-tree tests, so an installed module's `go test ./...` cannot meet it.
- On the board, a request claimed with a timestamp ahead of the viewer's clock is shown as clock skew rather than counted. Its user-request summary was never subscribed to the one-second refresh, so when the clock caught up the card's stopwatch moved while the header and drawer figures stayed on clock skew until something else redrew them. The summary now subscribes while a claim is skewed and recovers on the next tick.
