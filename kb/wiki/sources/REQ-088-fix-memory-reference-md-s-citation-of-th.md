---
title: "Lessons from REQ-088: Confirm: fix memory-reference.md's citation of the export-ignored CLAUDE.md"
type: source-summary
topic_cluster: knowledge-and-memory
sources: [raw/processed/2026-09-01/REQ-088-fix-memory-reference-md-s-citation-of-th.md]
related:
  - page: concept-knowledge-and-memory-systems
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-088: Confirm: fix memory-reference.md's citation of the export-ignored CLAUDE.md

Part of the [[concept-knowledge-and-memory-systems]] cluster.

## What the REQ was about

While working REQ-078 the builder noticed, in a file it was already editing, that
`actions/memory-reference.md:88` cites this repo's `CLAUDE.md` inline:

> Steps 3–4 are scoring/formatting the agent performs on the grep output — they need no additional
> shell state, so nothing carries between command blocks (CLAUDE.md: shell state does not survive
> between prescribed blocks).

## Solution summary

Replaced the parenthetical `(CLAUDE.md: shell state does not survive between prescribed blocks)` at line 88 with `(shell state does not survive between prescribed command blocks)`, exactly as authorized in `## Answer`. `git diff --numstat` reads `1 1` — one line changed, nothing else. The file now contains no reference to `CLAUDE.md`.

## What worked

- Running three grep shapes before trusting the REQ's one-site inventory. The REQ named one site; the shipped-paths sweep found six more of the same defect class in `tools/queue-kanban/*.go`, plus proof that the suite's own guard pattern matches none of them — including the defect this REQ was filed about. Running the suite's `self_citation_pattern` verbatim and getting **0 hits** turned "the guard is probably too narrow" into a measured fact.

## What didn't work

- Trusting the inherited test baseline. The checkpoint recorded "8 FAIL" as this repo's pre-existing suite state, confirmed twice last session by stash-and-compare. It does not reproduce: the suite exits 0 with zero FAIL lines, the seven update-script probes pass, and `tools/do-work-update.sh` demonstrably contains all five strings REQ-090 says are absent (lines 166, 194, 202, 204, 218). Checked for the obvious environmental causes — lingering worktrees (none), cwd sensitivity (green from a subdirectory), and a vacuous skip path (the probe has exactly one, `git` unavailable, and `git` is present). The cause of last session's observation is still unexplained, which is itself the finding.

## Worth knowing

- A stale `baseline.json` in `do-work/working/` outlives the session that wrote it and is silently available to the next REQ's Step 6.5 comparison. Route A skips pre-flight, so a Route A REQ that reaches for a baseline is always reading someone else's. Measure your own.

## Back-reference

See `do-work/archive/UR-015/REQ-088-memory-reference-cites-the-export-ignored-maintainer-doc.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bb8cf3b`.
