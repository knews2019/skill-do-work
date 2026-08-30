---
id: REQ-414
title: 'Migrate remaining core checks, publication helpers, Git helpers, and surveys'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-413]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Migrate Remaining Core Checks, Publication Helpers, Git Helpers, and Surveys

## What
Move all remaining core utility domain logic into `do-work-cli` subcommands.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Migrate preflight, qualification, scope drift, protected/uncommitted inventories, file association, and commit-hash recording.
- Migrate screenshot and download publication, local Git exclude, blocked checks, commit-diff display, exact deletion staging, timestamp/reservation helpers, and handoff surveys.
- Preserve existing exit status, output, filesystem effects, byte verification, permissions, redaction, and error behavior through characterization fixtures.
- Make findings actionable in both text and JSON without forcing an LLM to rescan.

## Constraints
- Existing `.sh` paths remain for compatibility but ultimately contain no domain logic.
- Target-specific Python project checks remain valid target probes rather than CLI implementation dependencies.

## Dependencies
Depends on REQ-413 (publication and common transaction primitives).

## Builder Guidance
Certainty level: Firm. Inventory every shipped core utility before migration and map each to a named CLI subcommand.

## Red-Green Proof
**RED prompt/case:** Run the current core utility fixture matrix against corresponding absent CLI subcommands and compare status, output, and filesystem effects.
**Why RED now:** Core behavior still lives across numerous shell implementations.
**GREEN when:** Every mapped Go subcommand passes parity fixtures and its compatibility path only builds/executes the canonical binary.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Folded From REQ-390 (2026-08-30)

- **`tools/checks/scope-drift.sh` reads every backticked token in a "Files I will
  touch" bullet as a declared path.** `emit_backticked_paths` splits the whole line
  on backticks and prints each odd-indexed part, so a bullet whose rationale contains
  an ordinary code span — `` - `path/to/file` (modify) — adds one `flex-wrap`
  declaration `` — declares a phantom second path and reports it as "declared in
  ## Scope but never touched". Observed on REQ-390, where it produced a false DRIFT
  line that had to be worked around by rewording the REQ rather than fixing the
  check. The path is the first backticked token on the bullet; the rest is prose.
  Worth closing when this REQ ports the checks to Go, where the extraction can be
  written once and tested.

