---
id: REQ-291
title: Browser behavior probe lane beside the Node behavior lane
status: completed
created_at: 2026-08-19T14:36:44Z
status_changed_at: 2026-08-19T14:36:44Z
claimed_at: 2026-08-21T01:01:12Z
completed_at: 2026-08-21T01:12:05Z
kb_status: pending
commit:
route: B
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
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The REQ names the shape to copy line by line (`runJavaScriptBehaviorProbe`, `lookupNodeForJavaScriptProbe`, the strict lane and its zero-probe guard) and states every requirement. What needed discovery was the exact wiring — where the zero-probe counter is gated, how the verify self-test insulates its fixtures, and whether `--dump-dom` is sufficient to get measurements back as data.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- `generate_test.go:22-27` holds the Node lane's three constants and its `atomic.Int64` counter; `:194-202` is `TestMain`, where the zero-probe guard actually lives. **`TestMain` is per-package, not per-file**, so the browser lane's guard had to be added there rather than in the new file.
- `:224` / `:240` are the zero-probe guard and the strict lane; `:204` `testEnvironmentWithOverrides` is the helper both use to scrub the child environment — reusable as is.
- `maintainer-verify.sh:477-485` is the Node lane's guarded invocation, and **`:286-330` is a self-test that runs the script against a fixture repo with stub binaries on a controlled PATH.** That self-test is what made the naive wiring fail: the fixture inherits the parent environment, so an exported `QUEUE_KANBAN_BROWSER` pointed the fixture's browser lane at a real engine while its `go` was a shim.
- `--headless --dump-dom --virtual-time-budget` on the available Chromium returns the post-script DOM, confirmed by smoke test before writing any Go. That settled the driver question without needing Playwright or npm.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (new) — the lane: discovery, probe runner, result extraction, the real probe, the strict lane and its zero-probe guard
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — extend `TestMain`'s zero-probe gate to the browser counter
- `_dev/tests/maintainer-verify.sh` (modify) — the guarded lane invocation, plus insulating the self-test fixtures from the override

**Files I will NOT touch:** `durations.go`, `durations_test.go`, or any measured-face constant. The REQ builds the capability; REQ-292 uses it.

**Acceptance criteria (restated from the REQ):**
1. A probe helper renders in a real engine and returns what the page measured, as data.
2. Prefer a driver needing no package manager; say in the code which was chosen and why.
3. Discovery: env override, then well-known binaries, then skip.
4. A strict lane that fails instead of skipping, plus a zero-probe guard.
5. No wall-clock waits.
6. `maintainer-verify.sh` runs it when a browser is present, prints an explicit SKIP when not, exits 0 either way.
7. One real probe: a `.durations-mark-label` `<text>` at 11px via `getBBox()`, asserted positive and finite.
8. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Decisions

- **D-01** (DECIDE & STATE): Driver is headless Chrome/Chromium invoked directly with `--dump-dom` and `--virtual-time-budget`; no Playwright, no npm. Reasoning: the REQ asks for the package-manager-free option unless it proves insufficient, and a smoke test confirmed `--dump-dom` returns the post-script DOM, which is all the "results cross as data" requirement needs. The choice and its rationale are written into the file's header comment, as the REQ requires — including the condition under which a future probe should reach for a real driver.
- **D-02** (DECIDE & STATE): An override naming a browser that does not exist is a **fatal error**, not a fall-through to PATH discovery. Reasoning: the caller asked for a specific engine; silently measuring in a different one would make the result a lie about which engine produced it, which is the exact failure the stale hand-transcribed constants came from.
- **D-03** (ESCALATE): Extended scope inside `maintainer-verify.sh` to insulate its three self-test fixture invocations with `QUEUE_KANBAN_BROWSER=''`. Reasoning: the fixtures run against stub binaries on a controlled PATH so their stages are deterministic; inheriting a real browser path made the fixture's browser lane run `go test` against a `go` shim and fail for a reason unrelated to what the fixture tests. Value: the self-test stays deterministic on a maintainer machine that has a browser — which, after this REQ, is the machine the lane is *for*. Risk: three more lines the self-test must keep in step if a fourth fixture invocation is added; each sits directly under its `PATH=` line, where the same reasoning already applies.
- **D-04** (DECIDE & STATE): The result contract is one DOM node holding JSON, and the page writes it **last and once**. Reasoning: it makes presence-of-node the readiness sentinel with no timing assumption, so a script that throws leaves the node empty and fails loudly rather than reporting a zero measurement — which is what `getBBox()` returns for an unrendered element and exactly how this lane could have passed forever.

## Implementation Summary

**What was done:** Added a browser behavior probe lane beside the existing Node lane, built to the same shape. A page renders in a real engine, writes its measurements as JSON into a single sentinel node, and the Go side parses that node and asserts on the numbers. Discovery is override-then-PATH-then-skip; a strict lane refuses to skip and a zero-probe guard stops an empty lane reporting green. `maintainer-verify.sh` runs the lane when an engine is present and prints an explicit SKIP naming what did not run when it is not.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (new) — 260 lines: lane constants, `lookupBrowserForBehaviorProbe`, `runBrowserBehaviorProbe`, `extractBrowserProbeResult`, `TestBrowserBehaviorMarkLabelTextExtent`, `TestMaintainerStrictBrowserBehaviorLane`, `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes`.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — `TestMain`'s zero-probe gate extended to the browser counter, with a comment saying why both lanes gate in one place.
- `_dev/tests/maintainer-verify.sh` (modified) — guarded lane invocation mirroring the Node lane's; three self-test fixture invocations insulated from the override (D-03).

**Tests touched:** the file is tests. Nothing existing changed meaning; `TestMain` gained a second guard beside the first.

## Qualification

Passed — 3 files verified, 8 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `shellcheck` clean on the edited script (it is in `maintainer-verify`'s 73-file lint set). No debug artifacts; the one `t.Logf` prints the measurement, which is the probe's product. Both `maintainer-verify.sh` paths exit 0.
- **Substantive:** the probe measured a real glyph — `"REQ-042 · 3h 20m"` at 11px came back 101.20 × 13.00. Not a stub.
- **Requirements traced:** AC1 → `runBrowserBehaviorProbe` returns JSON the test unmarshals; AC2 → D-01, stated in the file header; AC3 → `lookupBrowserForBehaviorProbe`; AC4 → both lane tests, each exercised; AC5 → `--virtual-time-budget` plus the sentinel node, no `sleep` anywhere in the diff; AC6 → both paths run below; AC7 → the probe and its five assertions; AC8 → verify exits 0.
- **Flowing:** the measurement crosses the boundary as parsed data and is asserted on numerically — not an exit code, not a stub value.

## Testing

- `go test -run '^TestBrowserBehavior' -v .` with a browser — PASS, logging `measured "REQ-042 · 3h 20m" at 11px in ui-sans-serif, …: 101.20 x 13.00`.
- `gofmt -l .` clean; `go vet ./...` clean; `shellcheck _dev/tests/maintainer-verify.sh` clean.

**Red-green validation** — the REQ's RED is "no test in this repo can observe a rendered glyph", so GREEN is a test that does:

| | Before | After |
|---|---|---|
| A test asks a real engine for a text extent and asserts on the answer | impossible — nothing could observe a rendered glyph | `TestBrowserBehaviorMarkLabelTextExtent` measures 101.20 × 13.00 and asserts positive, finite, 11px, and wider-than-tall |

**All three environment paths exercised, because the skip is only safe if the strict lane is not:**

| Path | Command | Result |
|---|---|---|
| Browser present | `go test -run '^TestBrowserBehavior'` | PASS, real measurement |
| No browser, ordinary run | same with an empty PATH | **SKIP**, suite still `PASS` |
| No browser, strict lane | `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` | PASS — i.e. it correctly observed the strict child **failing** with `strict browser behavior lane executed zero probes` |
| Browser present, strict lane | `TestMaintainerStrictBrowserBehaviorLane` | PASS |

**Both `maintainer-verify.sh` paths, end to end:**

| Environment | Exit | Output |
|---|---|---|
| `QUEUE_KANBAN_BROWSER` set | **0** | `maintainer-verify: queue-kanban strict browser behavior lane` |
| no browser anywhere | **0** | `SKIP: no browser is available; strict browser behavior lane was not run. Set QUEUE_KANBAN_BROWSER to name one.` |

**The zero-probe guard is not decorative:** `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` passes only by observing the strict child exit non-zero with the diagnostic, so it fails if the guard is removed from `TestMain`.

## Review

**Overall: 94%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Probe helper renders in a real engine; results cross as data | ✅ JSON in one node, unmarshalled and asserted numerically |
| Driver needing no package manager; choice and reason stated in code | ✅ headless Chrome/Chromium `--dump-dom`; rationale in the file header, including when to reach for a driver instead |
| Discovery: override → well-known binaries → skip | ✅ |
| Strict lane fails instead of skipping, plus zero-probe guard | ✅ both exercised |
| No wall-clock waits | ✅ virtual time + sentinel node; no `sleep` in the diff |
| `maintainer-verify.sh` runs it or prints an explicit SKIP; exits 0 either way | ✅ both paths run, both exit 0 |
| One real probe: mark-label `<text>` at 11px via `getBBox()`, positive and finite | ✅ plus three further assertions |
| `maintainer-verify.sh` exits 0 | ✅ |

### Findings

**Important — none.**

**Minor:**

- **M1:** The lane is Chromium-family only. Firefox and WebKit have no `--dump-dom`, so a machine with only those skips. That is the REQ's own trade (prefer the no-package-manager driver) and the override cannot rescue it — a future cross-engine need is the trigger to reach for a real driver, which the file header says explicitly.
- **M2:** D-03 touched the verify self-test's fixture invocations, which is adjacent to what this REQ is about. It was forced: without it the canonical gate fails on any machine that has a browser, which after this REQ is the machine the lane exists for.

**Nit:**

- **N1:** `browserProbeWellKnownBinaries` is a hand-maintained list, the shape `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale warns about. It is annotated as a convenience list rather than a closed set, and the override is the escape hatch that keeps it from being load-bearing — the same arrangement the Node lane has with a single name.

### Restatement Sweep

Redefined element: none — this REQ adds a lane and adds a guard beside an existing one. It changes no contract, field, or output shape.

Swept for statements that would now be stale: `maintainer-verify.sh`'s own stage list and its self-test assertions (`assert_success_stages` checks named stages; the browser stage is guarded and skipped in every fixture, so no assertion needed updating — verified by the self-test passing). `_dev/primes/prime-kanban-board.md` — grepped for statements about what the board's test lanes cover; it defers to the suite rather than enumerating lanes, so nothing there went stale. The Node lane's own comments still describe the Node lane correctly.

No stale restatement.

### Acceptance Testing

Every environment path was exercised rather than reasoned about, because the whole design rests on a skip being safe: browser-present (real measurement), no-browser-ordinary (skip, suite green), no-browser-strict (fails loudly with the diagnostic), browser-present-strict (passes). Then both `maintainer-verify.sh` paths end to end, both exiting 0 with the correct line. The measurement itself — 101.20 × 13.00 for a 16-character label at 11px — is plausible on inspection, which is the point: it is the first number in this repo about rendered text that nobody typed by hand.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope Discipline | 90% |
| Risk | Low |
| Acceptance | Pass |

Scope Discipline 90% for D-03. Risk Low rather than None: this adds a stage to the canonical gate, and a stage that can hang or flake would be felt on every commit — mitigated by there being no wall-clock wait anywhere in it and by the virtual-time budget bounding the browser's own run.

### Follow-up REQs Created

None. REQ-292, already queued and `depends_on` this REQ, is what the lane was built for and is now unblocked.

## Lessons Learned

**What worked:** Smoke-testing `--headless --dump-dom` from a shell before writing a line of Go. It settled the driver question — the REQ's "prefer no package manager, reach for Playwright only if that proves insufficient" needed an answer to "insufficient for what", and thirty seconds of shell gave it. Building the whole lane against a driver that then turned out to need npm would have been the expensive version of the same discovery.

**What didn't:** The first `maintainer-verify.sh` wiring failed the script's own self-test, and the reason is worth keeping: the self-test runs the script against a fixture repo with **stub binaries on a controlled PATH**, but it inherits the rest of the environment. An exported `QUEUE_KANBAN_BROWSER` therefore pointed the fixture's browser lane at a real engine while its `go` was a shim. Any environment variable a new stage reads has to be neutralized at those fixture invocations, right beside the `PATH=` line that is already doing exactly that job for the same reason.

Also: `shellcheck` rejects `VAR= \` as a probable typo (SC1007) and wants `VAR='' \`. Worth knowing before writing three of them.

**Worth knowing:** `getBBox()` returns zeros for an unrendered or detached element, so a browser probe's default failure mode is a *successful-looking measurement of nothing*. That is why the result node is written last and only once, and why the assertions check positive-and-finite plus a known font size rather than merely "no error". A browser lane that renders nothing measurable passes forever, and that is the specific failure this probe was built to be incapable of.

## Orientation

The repo can now ask a real rendering engine for a measurement inside a test, instead of a person running a browser by hand and pasting the number into a comment. That ritual is why the Durations view's measured-face constants are stale and admittedly not a supremum over the face space. Lives beside the existing Node behavior lane in the board tool's test suite (`skills/do-work-board/tools/queue-kanban/browser_probe_test.go`), wired into the canonical gate the same guarded way, and indexed by `_dev/primes/prime-kanban-board.md`.

**[MAP CHANGED]** — a second behavior lane exists, with its own marker, run pattern, probe counter, and zero-probe guard, and `maintainer-verify.sh` gained a stage. Anyone adding a probe now has two lanes to choose between, and the choice is whether the assertion needs a rendering engine or only a JavaScript runtime. Nothing about the Durations view changed: this REQ built the capability and proved it works.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-shell-commands.md` — referenced paths still resolve; neither enumerates test lanes, so neither was made stale.
