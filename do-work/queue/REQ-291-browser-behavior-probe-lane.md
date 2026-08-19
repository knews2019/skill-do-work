---
id: REQ-291
title: Browser behavior probe lane beside the Node behavior lane
status: pending
created_at: 2026-08-19T14:36:44Z
status_changed_at: 2026-08-19T14:36:44Z
user_request: UR-061
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/browser_probe_test.go
- _dev/tests/maintainer-verify.sh
---

# Browser Behavior Probe Lane Beside the Node Behavior Lane

## What

Every font measurement recorded in `durations_test.go` got there by a person running
Playwright by hand and pasting the number into a comment with the browser build written
beside it. That ritual is the reason the constants are stale — the repo's own comment admits
the recorded box height "is NOT a supremum over the face space" — and it is the reason
REQ-278 proposed surveying operating systems by hand. The user rejected the survey and chose
to move placement into the browser instead (UR-061), which needs a way to assert against a
real rendering engine automatically.

This REQ builds only the lane. Nothing about the Durations view changes here.

The shape already exists one layer over: `runJavaScriptBehaviorProbe` pipes a probe to
`node -` and fails the test on a non-zero exit, `lookupNodeForJavaScriptProbe` skips when
Node is absent, `TestMaintainerStrictJavaScriptBehaviorLane` refuses to skip when selected
directly, and `TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes` stops an empty
lane passing as a green one. Follow that shape rather than inventing a second convention.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- **A probe helper that renders a page in a real engine and returns what the page measured.**
  The page writes its results into one DOM node; the Go side reads that node and parses it.
  Results cross the boundary as data the test can assert on, never as an exit code alone.
- **Prefer a driver that needs no package manager.** A headless browser invoked directly
  against a temporary HTML file — Chrome and Chromium both take `--headless --dump-dom` with
  `--virtual-time-budget` — needs only a binary that is already on most machines. Reach for
  Playwright or npm only if that proves insufficient, and **say in the code which was chosen
  and why**; a silent npm dependency is a build input nobody agreed to.
- **Discovery, then skip.** An environment-variable override is consulted first, then a short
  list of well-known binaries. When none is found the probe skips exactly as
  `lookupNodeForJavaScriptProbe` does. A machine with no browser still runs everything else.
- **A strict lane that fails instead of skipping**, mirroring
  `TestMaintainerStrictJavaScriptBehaviorLane`, plus its zero-probe guard so a lane that ran
  nothing cannot report green. Both are what make the skip safe.
- **No wall-clock waits anywhere.** Readiness is a sentinel the page writes into the dumped
  DOM, or the driver's own virtual-time budget. A `sleep` in a test lane is the flake this
  lane will be blamed for.
- **`_dev/tests/maintainer-verify.sh` gains the lane guarded the way the Node lane already is:**
  run it when a browser is present, print an explicit SKIP line naming what did not run when
  it is not. The script exits 0 in both cases.
- **One real probe, not a decorative one.** Measure a `.durations-mark-label` `<text>` at the
  board's own 11px through `getBBox()` and assert the returned box is positive and finite.
  It names the failure it pins: a lane that renders nothing measurable is a lane that will
  pass forever.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED:** no test in this repo can observe a rendered glyph. The measured-face constants are
hand-transcribed, and `TestDurationsMeasuredConstantsNameTheirChromiumBuild` can only check
that a human wrote a build name in a comment — it cannot check the number.

**GREEN:** a test asks a real engine for a text extent and asserts on the answer. On a
machine without a browser the suite still exits 0, and the maintainer-strict selection fails
loudly rather than skipping quietly.

## Context

Prerequisite for REQ-292, which cannot be verified without it.

The lane's value is not this one probe. It is that the next person who needs a rendered
measurement gets it from a test instead of from a Playwright session and a comment.

Deliberately **not** in scope: changing any Durations geometry, deleting any existing
constant, or touching `durations.go`. This REQ adds a capability and proves it works.
