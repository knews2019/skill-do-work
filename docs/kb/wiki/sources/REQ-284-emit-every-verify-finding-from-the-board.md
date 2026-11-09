---
title: "Lessons from REQ-284: Emit every verify finding from the board's Go producer"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-284-emit-every-verify-finding-from-the-board.md]
related:
  - page: REQ-280-probe-timestamp-ordering-and-point-check
    rel: complements
  - page: REQ-281-reconcile-the-calibration-log-against-th
    rel: complements
  - page: REQ-285-render-a-verify-findings-strip-on-the-bo
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-284: Emit every verify finding from the board's Go producer

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Split `runVerifyProbes` into a board-taking `collectVerifyFindings(repoRoot, board, now)` plus a thin
wrapper that builds the board first, then carry the resulting findings into `generatedBoardData` as
`verifyFindings` and `verifySkipped`. Suppress the three categories the board already renders, in the
producer, so the client can render the list blindly. Wire both callers: `generate` for the static
snapshot, and `serve` per `/board-data.js` request outside the mtime cache.

## Solution summary

Split `runVerifyProbes` into a board-taking `collectVerifyFindings` plus a thin wrapper, then carried the findings into the board payload as `verifyFindings`/`verifySkipped` for both producers. Suppression of the three already-rendered categories and reduction of absolute paths both happen once, in the Go producer. Serve computes findings fresh on every `/board-data.js` request, outside the mtime cache and with no cache of their own.

## What worked

Checking the call-site count before changing a signature. `buildGeneratedBoardData` looked like the obvious place to thread the new data through until `grep` showed twelve callers, ten of them tests with no interest in verify — at which point attaching at the two real callers became obviously right. The second thing: proving "the CLI report is byte-identical" by actually building both binaries and diffing, rather than reasoning that a pure refactor could not have changed it. It cost two minutes and turned an assumption into evidence.

## What didn't work

The first serve wiring referenced `liveServer.cachedBoard`, a field that did not exist — the server cached the *projected payload* but not the board behind it. The REQ's own text said "calling `collectVerifyFindings(repoRoot, cachedBoard, ...)`", which read as describing something already there. Worth noting as a small instance of the same pattern this session keeps hitting: a capture-time description of the code is a hypothesis about the code.

## Worth knowing

The probes must never be cached, and the reason is not performance — it is that two of their inputs are not files. Wall-clock claim age and `git worktree list` both change while every mtime stays identical, which is exactly the blind spot that let two REQs sit claimed for 13 hours against a 3-hour threshold with the age printed on their cards. `TestCollectVerifyFindingsSeesAClaimGoStaleWithoutAnyFileChanging` is the guard: it reuses one board across two calls and only advances `now`. If someone later "optimizes" the per-request probes behind the mtime cache, that test is what should stop them.

## Back-reference

See `do-work/archive/UR-058/REQ-284-emit-verify-findings-from-the-board-producer.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8f61f69`.
