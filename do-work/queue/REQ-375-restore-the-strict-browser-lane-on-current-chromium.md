---
id: REQ-375
title: '[impact-rule-change] Restore the strict browser lane on current Chromium'
status: pending-answers
created_at: 2026-08-26T14:40:00Z
user_request: UR-074
addendum_to: REQ-374
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
---

# Restore the Strict Browser Lane on Current Chromium

## What

`TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` fails on Chromium 141.0.7390.37, tripping its own vacuity guard: *"the capture-swallowed outside-release trial never crossed the host pointerleave boundary; the isolator was not exercised and the mutation pair is vacuous."* Because it fails, `TestMaintainerStrictBrowserBehaviorLane` fails, so the whole strict browser lane is unavailable as a gate.

## Context

Discovered while working on REQ-374. Reproduced at HEAD in a separate worktree with REQ-374's diff absent, so it predates that REQ.

The root cause is already recorded in `do-work/prose-backlog.md` as a stale stated reason: `web/board-timeline.js:2571-2574` and `timeline_browser_probe_test.go:2196-2197` both justify the Timeline's pointer capture with "Chromium suppresses the boundary events while a button is held", which did not reproduce on Chromium 1194 — a buttoned exit delivered `pointerleave` to the host four times. The backlog line covers the prose; this REQ covers the failing test, which is not prose-only.

The prime's own note applies: a measured browser behaviour is per-browser, and this probe was calibrated against an engine that behaved differently.

## Red-Green Proof

**RED prompt/case:** `QUEUE_KANBAN_BROWSER=<chromium-141> go test -count=1 -run '^TestMaintainerStrictBrowserBehaviorLane$' .` in `skills/do-work-board/tools/queue-kanban` fails, naming the timeline pointer-capture probe.
**Why RED now:** the probe assumes Chromium suppresses boundary events while a button is held; Chromium 141 delivers `pointerleave` to the host, so the trial the mutation pair depends on never runs and the guard correctly refuses to pass vacuously.
**GREEN when:** that command exits 0 on Chromium 141, with the probe still exercising the isolator — the guard must be satisfied by driving the gesture the engine actually produces, never by weakening or deleting the guard.
**Validation:** Inferred during capture — the failure is reproduced and understood; the fix is not.

## Open Questions
- [ ] I discovered this out-of-scope task while working on REQ-374: the strict browser behavior lane fails at HEAD on Chromium 141 because the timeline pointer-capture probe's vacuity guard fires. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

## Full Context
See `do-work/user-requests/UR-074/input.md` and REQ-374's `## Discovered Tasks`.
