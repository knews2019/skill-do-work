---
id: REQ-551
title: '[impact-negligible] Delete the five caller-less toolbox shell shims and re-point their fixtures at do-work-cli'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-549, REQ-550, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [skills/do-work-toolbox/scripts/, _dev/tests/prescribed-shell-cases/generate-report-image.sh, _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh, _dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh, _dev/tests/prescribed-shell-cases/install-last30days.sh, _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/staged-skills-contract.sh, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-04T13:45:29Z
---

# Delete the five caller-less toolbox shell shims and re-point their fixtures at do-work-cli

## What
All five scripts under `skills/do-work-toolbox/scripts/` are pure `exec` pass-throughs to a `do-work-cli` subcommand with no caller outside `_dev/tests`; the shipped guide already routes every action to `tools/do-work-cli.sh … <subcommand>`. Delete the five scripts, re-point their five behaviour fixtures at the CLI directly, and remove their rows from the canonicalization ratchet and the staged-package list.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Dead shipped surface: 63 lines nothing invokes, four rows in a ratchet, three rows in a staged list, all maintained for no caller.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 5, sweep_key `toolbox-shims-no-callers`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -80. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- Delete `architecture-report-preflight.sh`, `generate-report-image.sh`, `generate-report-image-batch.sh`, `publish-portfolio-summary.sh`, `install-last30days.sh` under `skills/do-work-toolbox/scripts/`.
- Re-point the five `_dev/tests/prescribed-shell-cases/*.sh` fixtures (1242 lines; they test the Go subcommands through the shim) at `tools/do-work-cli.sh --format text <subcommand>` so the behaviour coverage is kept; fixture bodies otherwise unchanged.
- Remove the four `skills/do-work-toolbox/scripts/*.sh` rows from the existence list in `_dev/tests/prescribed-shell-canonicalization.sh` and the three rows from the staged list in `_dev/tests/staged-skills-contract.sh`.
- Two of the five carry a KEEP row in `decisions/audits/2026-08-11-defensive-surface.md`; those rows covered shell bodies that 0.260.1 moved into Go, not the shim that remains. Record that in the commit message; the frozen audit file is not edited.
- Not in scope: the caller-less core launchers `skills/do-work/tools/{do-work-update,fetch-upstream-archive,validate-suite-manifest}.sh` and their root `tools/` mirrors, retained by the 0.260.1 decision and the KEEP table.
- Reproduce at dc8a64e3 (prints the five paths): `for f in skills/do-work-toolbox/scripts/*.sh; do b=$(basename "$f"); n=$(rg -l --fixed-strings "$b" skills tools suite README.md CLAUDE.md _dev/primes --glob '!*CHANGELOG*' | grep -v "/$b$" | wc -l | tr -d ' '); [ "$n" -eq 0 ] && echo "$f"; done`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (create it on first use, executable, invoked from `_dev/tests/contract-regressions.sh` in the fast tier the way `_dev/tests/defensive-surface-audit.sh` is, with the same missing-or-not-executable FAIL line), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- The lock-in carries a second assertion whose previous home (`_dev/tests/shipped-shell-thinness.sh`) was deleted on 2026-09-03 in 0.266.9: `find skills tools -name '*.sh' | while read -r f; do rg -q do-work-cli "$f" || echo NON-DELEGATING: $f; done` must print nothing.
- Prime `_dev/primes/prime-shell-commands.md` before touching the fixtures.
- Lock-in limit: caller-less toolbox shims: 0 (today 5); non-delegating shipped .sh files: 0 (today 0).

## Dependencies
No dependency. REQ-553 (shared launcher preamble) depends on this REQ because it deletes four of the thirteen preamble copies.

## Builder Guidance
Firm on deletion; the fixtures are kept because they are the only behaviour coverage of those subcommands through shell.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints the five script paths.
**GREEN when:** It prints nothing; the five fixtures pass against `do-work-cli.sh` in the fast gate; the lock-in pins caller-less toolbox shims at 0 and NON-DELEGATING shipped shell at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for toolbox-shims-no-callers.*
