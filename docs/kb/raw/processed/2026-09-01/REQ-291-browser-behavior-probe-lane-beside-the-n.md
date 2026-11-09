---
source_type: req_lesson
req_id: REQ-291
req_path: do-work/archive/UR-061/REQ-291-browser-behavior-probe-lane.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, browser, behavior, probe, lane]
---

# Lessons from REQ-291: Browser behavior probe lane beside the Node behavior lane

## What the REQ was about

Every font measurement recorded in `durations_test.go` got there by a person running
Playwright by hand and pasting the number into a comment with the browser build written
beside it. That ritual is the reason the constants are stale — the repo's own comment admits
the recorded box height "is NOT a supremum over the face space" — and it is the reason
REQ-278 proposed surveying operating systems by hand. The user rejected the survey and chose
to move placement into the browser instead (UR-061), which needs a way to assert against a
real rendering engine automatically.

This REQ builds only the lane. Nothing about the Durations view changes here.

## Solution summary

Added a browser behavior probe lane beside the existing Node lane, built to the same shape. A page renders in a real engine, writes its measurements as JSON into a single sentinel node, and the Go side parses that node and asserts on the numbers. Discovery is override-then-PATH-then-skip; a strict lane refuses to skip and a zero-probe guard stops an empty lane reporting green. `maintainer-verify.sh` runs the lane when an engine is present and prints an explicit SKIP naming what did not run when it is not.

## What worked

Smoke-testing `--headless --dump-dom` from a shell before writing a line of Go. It settled the driver question — the REQ's "prefer no package manager, reach for Playwright only if that proves insufficient" needed an answer to "insufficient for what", and thirty seconds of shell gave it. Building the whole lane against a driver that then turned out to need npm would have been the expensive version of the same discovery.

## What didn't work

The first `maintainer-verify.sh` wiring failed the script's own self-test, and the reason is worth keeping: the self-test runs the script against a fixture repo with **stub binaries on a controlled PATH**, but it inherits the rest of the environment. An exported `QUEUE_KANBAN_BROWSER` therefore pointed the fixture's browser lane at a real engine while its `go` was a shim. Any environment variable a new stage reads has to be neutralized at those fixture invocations, right beside the `PATH=` line that is already doing exactly that job for the same reason.

Also: `shellcheck` rejects `VAR= \` as a probable typo (SC1007) and wants `VAR='' \`. Worth knowing before writing three of them.

## Worth knowing

`getBBox()` returns zeros for an unrendered or detached element, so a browser probe's default failure mode is a *successful-looking measurement of nothing*. That is why the result node is written last and only once, and why the assertions check positive-and-finite plus a known font size rather than merely "no error". A browser lane that renders nothing measurable passes forever, and that is the specific failure this probe was built to be incapable of.

## Back-reference

See `do-work/archive/UR-061/REQ-291-browser-behavior-probe-lane.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6fa130d`.
