---
source_type: req_lesson
req_id: REQ-375
req_path: do-work/archive/UR-074/REQ-375-restore-the-strict-browser-lane-on-current-chromium.md
date: 2026-08-27
domain: testing
module: _dev/primes
tags: [testing, restore, strict, browser, lane]
---

# Lessons from REQ-375: Restore the strict browser lane on current Chromium

## What the REQ was about

`TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` fails on Chromium 141.0.7390.37, tripping its own vacuity guard: *"the capture-swallowed outside-release trial never crossed the host pointerleave boundary; the isolator was not exercised and the mutation pair is vacuous."* Because it fails, `TestMaintainerStrictBrowserBehaviorLane` fails, so the whole strict browser lane is unavailable as a gate.

## Solution summary

- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modified). Reuses the existing DevTools-pipe session for measurement probes, waits for a populated result node on the page's real clock, reads literal textContent, and retains object-shape/caller JSON validation and strict probe coun

## Worth knowing

Browser process exit is not a reliable result-readiness signal: this Chrome build emitted complete dump-DOM output but did not exit, even for a tiny local page. A bounded protocol read of the completed result node restores observability without relaxing product assertions. Read textContent when the contract is JSON text; serializing HTML can silently change literal clipboard content. A failed mutation-isolator trial identifies that trial's missing event, not all event behavior in that engine or every older release.

## Back-reference

See `do-work/archive/UR-074/REQ-375-restore-the-strict-browser-lane-on-current-chromium.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `54de194b`.
