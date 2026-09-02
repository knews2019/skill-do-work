---
id: REQ-519
title: '[impact-rule-change] Path-scoped fast lane for the maintainer gate'
status: pending
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-518]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-518, REQ-520, REQ-521, REQ-522, REQ-523]
batch: cheap-maintainer-gate
write_set: [_dev/tests/maintainer-verify.sh, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, CLAUDE.md, justfile, _dev/tests/contract-regressions.sh]
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-02T21:34:07Z
  basis:
    - Route B
    - 6-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
---

# Path-Scoped Fast Lane for the Maintainer Gate

## What

Add `bash _dev/tests/maintainer-verify.sh --changed`: a fast lane that scopes the gate to what the working tree and the commits since the last green revision actually touched. Lint and vet always run; a Go module's tests run only when its files changed, with the Go test cache on; the aggregate contract suite runs only when action Markdown, shipped scripts, or `_dev/tests` changed. The pipeline runs the fast lane per REQ at Step 6.5 and the full uncached gate once at the integrating commit. Also freeze `_dev/tests/contract-regressions.sh` at its current size with a ratchet (item A6, captured as policy).

The fold-first scan found no pending REQ that owns a fast lane; REQ-510 (Sweep work-reference sections whose contract is now a CLI behavior test) shrinks the contract file but adds no ratchet.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

There is one tier for every purpose. A one-line prose edit to an action file pays for 211 seconds of board tests and 39 seconds of CLI tests it cannot affect, and `-count=1` disables Go's test cache everywhere. Capture-time answers Q2 and Q4: fast lane per REQ, full gate at integration, typical REQ gate under 60 seconds.

## Context

The gate's self-test already models its stage set with shimmed tools and asserts each stage runs exactly once, so the fast lane's stage selection can be proven the same way. `CLAUDE.md` § Verify names the full script as the canonical baseline check before hand-back; that sentence changes to name the fast lane for a builder's hand-back and the full gate for the integrating commit, mirroring how version bumps are owner-only. `justfile` gets a matching recipe.

## Detailed Requirements

- `maintainer-verify.sh --changed` derives the changed path set from `git diff --name-only` against the last recorded green revision (REQ-518) plus the working tree, and selects the expensive stages by path class: `*.go` under a module selects that module's `go test` (cached, no `-count=1`); `skills/**/*.md`, `skills/*/scripts/**`, `tools/**`, or `_dev/tests/**` select the aggregate contract suite; the JavaScript probes run only when `skills/do-work-board/tools/queue-kanban/web/**` or its Go tests changed.
- Version floors, ShellCheck, gofmt and `go vet` for both modules run in every lane; together they cost about 6 seconds, and the input says "lint and vet always". (verify-requests, 2026-09-03: the earlier path-scoped lint and vet contradicted the What section.)
- The lane prints, before running, which stages it selected and which it skipped and why, so a skip is never silent.
- The self-test gains fixtures for the fast lane: a Markdown-only change runs no Go tests; a queue-kanban-only change runs only that module; an empty change set runs floors, lint and gofmt only.
- `actions/work.md` Step 6.5 runs the fast lane for a builder's hand-back; Step 9's integrating commit runs the full gate. `CLAUDE.md` § Verify and `work-reference.md` say the same thing once.
- A6 policy: `CLAUDE.md` gains one sentence under Verify: a new sentence-predicate lane in `contract-regressions.sh` must delete one or land as a Go behavior test. The self-test (or the aggregate itself) fails when `contract-regressions.sh` exceeds 8,417 lines; lowering the ratchet is a one-line edit whenever the file shrinks.
- Out of scope, by the maintainer's capture-time choice (Q1): splitting the existing `contract-regressions.sh` by owning action file. The verbatim input's A6 mentions it; do not do it under this REQ. (verify-requests, 2026-09-03)

## Constraints

- The full gate's behavior, argv, and output are unchanged; `--changed` is additive.
- The fast lane never counts as the release check and never records a green revision for REQ-518.
- No new prose that walks the shell sequence; the script owns the selection rules.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

Depends on REQ-518 (Run the full gate once per REQ): the changed-path set is computed against the green revision that REQ records, and both edit the same gate lanes in `work.md`. Related to REQ-510 by subject.

## Builder Guidance

Certainty level: Firm on the two-tier rule and the ratchet; latitude on path classes beyond the ones listed. Keep the selection table in the script beside the stage it selects.

## Red-Green Proof

**RED prompt/case:** `bash _dev/tests/maintainer-verify.sh --changed` after editing only `skills/do-work/actions/help.md`.
**Why RED now:** The script exits 2 with a usage message; the only lane is the full 6.5-minute run.
**GREEN when:** The command runs floors, lint, gofmt and the aggregate suite, skips both Go modules with a printed reason, and finishes under 60 seconds on the maintainer's machine; the self-test fixtures for the three path classes pass; `work.md` Step 6.5 names the fast lane and Step 9 names the full gate; adding a line to `contract-regressions.sh` past 8,417 turns the ratchet red.
**Validation:** User confirmed (capture-time answers Q1, Q2, Q4, 2026-09-03).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing shipped shell and prescribed command blocks.
- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing pipeline fields and the run action's gate lanes.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A2 and A6 of the analysis report's improvements, captured by UR-100.*
