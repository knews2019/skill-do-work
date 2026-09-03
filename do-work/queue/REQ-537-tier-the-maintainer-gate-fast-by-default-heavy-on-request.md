---
id: REQ-537
title: 'Tier the maintainer gate: fast by default, --heavy on request'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - _dev/tests/maintainer-verify.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/probe-batch.sh
---

# Tier the Maintainer Gate: Fast by Default, --heavy on Request

## What

`bash _dev/tests/maintainer-verify.sh` with no flags becomes the fast tier and the only canonical gate: toolchain floors, ShellCheck, gofmt, `go vet` for both modules, the fast aggregate, queue-kanban tests without JavaScript or browser probes, and do-work-cli tests. `--heavy` adds the board's JavaScript probes (strict marker on, zero-probe guard on), the browser lane when an engine is present, and the heavy aggregate probes. `--heavy-surfaces` prints the repo-relative path globs whose changes warrant a heavy run, so REQ-541 can consult it mechanically.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"then we could have a heavy weight test, but that should be executed only with user permission". Today's 301 s gate runs on every step; the maintainer decided the fast tier is the gate and the heavy tier is theirs alone (D1, D2).

## Context

- Today: lint and vet about 10 s, aggregate 112 s, queue-kanban 149 s, do-work-cli 35 s wall. The JavaScript probes are 48 s of the board run and the two re-exec meta-tests another 58 s (REQ-538 owns those).
- `_dev/tests/gate-runner.sh` and every `skills/` caller invoke the bare script and never pass `--heavy` (D2). The gate-evidence record is keyed by exact argv, so fast and heavy evidence are already distinct.
- The shimmed `--self-test` (`run_self_test`, `assert_success_stages`) enumerates stages per run; it must cover both tiers, and the strict-marker mutation check stays.

## Detailed Requirements

- Default run: floors, ShellCheck, gofmt, both `go vet`, fast aggregate, queue-kanban tests with JavaScript probes skipped through REQ-538's knob (until REQ-538 lands, pass the marker off and let probes run), do-work-cli tests.
- `--heavy`: everything the default runs plus JavaScript probes with `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1`, the browser lane when present, and the heavy aggregate probes (REQ-539 classifies them; until then the aggregate stays whole).
- `--heavy-surfaces`: prints globs, one per line: the board's `web/**` and its JavaScript probe test files, installer and updater sources under `skills/do-work/tools/` and `tools/`, and anything else you classify heavy. Keep the list in the script, not in prose.
- Self-test: expected stage lists per tier; a `--heavy` fixture run; the marker-mutation check unchanged.
- No caller under `skills/` or the gate runner passes `--heavy`.
- Supersedes REQ-519's path-scoped `--changed` lane: the fast tier is the per-REQ check for every path, so no path scoping is built. If the fast tier cannot reach 120 s without it, record that in the commit body and capture path scoping separately; do not add it here silently.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** `bash _dev/tests/maintainer-verify.sh --heavy` and `bash _dev/tests/maintainer-verify.sh --heavy-surfaces`.
**Why RED now:** both exit 2 with the usage line; the default run executes the JavaScript probes and the full aggregate.
**GREEN when:** the default run exits 0 in under 120 s once REQ-538 to REQ-540 land (under 200 s on its own), `--heavy` exits 0 and executes the JavaScript probes, `--heavy-surfaces` prints globs, and `--self-test` exits 0.
**Validation:** User confirmed (D1, D2)

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a shipped-style shell harness and its self-test.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.

