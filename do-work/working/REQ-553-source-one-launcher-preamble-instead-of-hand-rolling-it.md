---
id: REQ-553
title: '[impact-negligible] Source one do-work-cli launcher preamble instead of hand-rolling it in every launcher'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-551]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
write_set: [tools/do-work-cli-preamble.sh, skills/do-work/tools/do-work-cli-preamble.sh, tools/install-do-work-suite.sh, tools/fetch-upstream-archive.sh, tools/replace-text-section.sh, tools/validate-suite-manifest.sh, skills/do-work/tools/install-do-work-suite.sh, skills/do-work/tools/fetch-upstream-archive.sh, skills/do-work/tools/replace-text-section.sh, skills/do-work/tools/validate-suite-manifest.sh, skills/do-work/tools/checks/associate-files.sh, _dev/tests/staged-skills-contract.sh, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-05T00:41:15Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route B
    - 13-file write set
    - 3 subsystems involved
    - 4 acceptance criteria
---

# Source one do-work-cli launcher preamble instead of hand-rolling it in every launcher

## What
After REQ-551 deletes the four toolbox copies, nine shell files still hand-roll the do-work-cli launcher preamble in two spellings (`for cli_candidate in` in the four root tools and their four byte-locked mirrors; `launcher_arguments=(--format text)` in `tools/checks/associate-files.sh`). Promote one sourceable preamble beside the root tools, with its byte-locked mirror under `skills/do-work/tools/`, that resolves the launcher path and the `--format text` argument array, and source it from every remaining copy.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Two mutually incompatible spellings of one primitive, copied instead of sourced, is the 0.202 canonicalization class; the ratchet `_dev/tests/prescribed-shell-canonicalization.sh` does not look at this pattern.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 2, sweep_key `cli-launcher-preamble-copied`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -40. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- New `tools/do-work-cli-preamble.sh` (sourced, not executed) and its mirror `skills/do-work/tools/do-work-cli-preamble.sh`; add the pair to the byte-identical mirror list in `_dev/tests/staged-skills-contract.sh`.
- The bootstrap path stays self-contained: `tools/install-do-work-suite.sh --print-bootstrap-command` prints a literal heredoc and must keep working before any skill is installed, so the preamble lives beside the root tools, never under `skills/do-work/scripts/`.
- `skills/do-work/tools/checks/associate-files.sh` switches from the `launcher_arguments=(--format text)` spelling to the sourced preamble.
- Root tools and their mirrors change in lock-step (the staged-skills contract requires byte identity).
- Reproduce at dc8a64e3 (prints 13 files; 9 after REQ-551): `rg -n --no-heading 'for cli_candidate in|^launcher_arguments=\(--format text\)$' --glob '*.sh' skills tools`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (create it on first use, executable, invoked from `_dev/tests/contract-regressions.sh` in the fast tier the way `_dev/tests/defensive-surface-audit.sh` is, with the same missing-or-not-executable FAIL line), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- No launcher behaviour change: the differential fixtures for the installer and updater must pass unchanged.
- Prime `_dev/primes/prime-shell-commands.md` first; the sourcing must survive `set -euo pipefail` and a missing Go toolchain the same way the current preambles do.
- Lock-in limit: hand-rolled launcher preambles outside the preamble pair: 0 after this REQ (today 13); the Reproduce command prints at most 2 paths, the preamble file and its byte-identical mirror (verify repair 2026-09-03: the plan line's "target ≤ 1" counted the helper once, the mirror makes it two).

## Dependencies
Depends on REQ-551, which deletes four of the thirteen copies so this REQ touches nine files instead of thirteen.

## Builder Guidance
Firm on one sourced preamble and on the bootstrap constraint; latitude on the helper's name and on whether the argument array is a function or a variable.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements after REQ-551 has landed.
**Why RED now:** It prints nine files hand-rolling the preamble.
**GREEN when:** It prints at most the preamble file itself and its mirror; installer and updater fixtures green; the lock-in pins hand-rolled preamble copies at 0 outside the preamble pair.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for cli-launcher-preamble-copied.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear — one sourceable preamble replacing nine hand-rolled copies — but how each of the nine callers must locate the shared file, and whether the byte-locked mirror needs a modules.tsv declaration, has to be established by reading the callers. Exploration required; the shape of the fix is not in doubt.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*
