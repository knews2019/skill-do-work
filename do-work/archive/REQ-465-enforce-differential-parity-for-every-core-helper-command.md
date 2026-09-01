---
id: REQ-465
title: '[impact-rule-change] Review fix: Enforce differential parity for every core helper command'
status: cancelled
created_at: 2026-09-01T02:23:45Z
user_request: UR-081
domain: testing
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-414]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-414, REQ-420, REQ-462, REQ-463, REQ-464]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-414
sweep: true
sweep_key: core-helper-command-differential-parity
completed_at: 2026-09-01T09:18:35Z
---

# Enforce Differential Parity for Every Core Helper Command

## What

Replace the registration/smoke matrix with retained-versus-Go characterization for all 17 core helper commands. Done means every semantic status, ordered fact, affected path, actionable projection, and filesystem/Git effect is compared mechanically, and mutation tests prove the adapter cannot silently accept a divergence.

The fold-first scan found REQ-420 with the same eventual whole-suite parity intent, but it is not dependency-ready because REQ-419 is pending and therefore is not conversion-eligible. No eligible pending or pending-answers REQ or sweep shares this command-matrix root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `internal/corehelpers/commands_test.go`: the all-17 matrix accepts any exit status up to one and checks renderer agreement/generic collection presence rather than retained semantic facts and effects. (found by REQ-414 / UR-081)
- [ ] Parity adapter: no mutation proof demonstrates that changing a status, fact, path, action, or side effect makes the matrix fail. (found by REQ-414 / UR-081)

## Requirements

- Run every retained helper beside its Go command on the same characterized fixtures in text and JSON modes.
- Compare exact semantic status, ordered facts, affected paths, recovery/verification argv, and filesystem/index/private-state effects for all 17 surfaces.
- Cover happy, non-clean, hostile-path, combined-state, dry-run, refusal, and concurrent-change cases earned by each helper's retained contract.
- Mutation-test the parity adapter across status, fact, action, path, and effect dimensions.
- Coordinate with REQ-420 so the eventual shim conversion reuses this authority rather than building a second parity framework.

## Red-Green Proof

**RED prompt/case:** Run the current matrix after introducing a controlled `AD` inventory misclassification, unborn-repository reservation deletion, collapsed dirty-path evidence, or wrong recovery argv.
**Why RED now:** All of those observable divergences pass the smoke-only matrix.
**GREEN when:** Each controlled divergence fails the named differential fixture, and the unmodified retained and Go implementations agree across all characterized statuses, facts, actions, paths, and effects.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-414-rereview.md`.

---
*Source: REQ-414 fresh re-review finding 4.*

## Cancelled

- **When:** 2026-09-01T09:18:35Z
- **Why:** folded into REQ-420 as parity acceptance criteria (maintainer decision, 2026-09-01 queue analysis); no separate parity framework before the shim conversion
- **Decided by:** user, via `do-work abandon`
