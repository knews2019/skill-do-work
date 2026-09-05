---
id: REQ-555
title: '[impact-negligible] Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form'
status: pending
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-554]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/audit-lockins.sh]
---

# Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form

## What
The "Shipped executable homes" table in `skills/do-work/docs/prescribed-shell-primitives.md` assigns owned mechanics to nine `*.sh` paths that are each a 6-to-11-line `exec` shim over `do-work-cli.sh` (the mechanics moved to Go at 0.260.1), and one sentence below it says `scripts/protected-inventory.sh` "orchestrates" two check scripts, which a six-line shim cannot do. Reword the route column to the `tools/do-work-cli.sh … <subcommand>` form the toolbox rows already use and delete the orchestration sentence.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The guide is the pointer target from 16 shipped files (ratchet-enforced) and it currently misroutes readers to shims and describes an orchestration that no longer exists.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 7, sweep_key `stale-shell-ownership-prose`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -5. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- Seven rows naming `scripts/show-commit-diff.sh`, `add-local-git-exclude.sh`, `atomic-download.sh`, `capture-screenshot.sh`, `run-blocked-check.sh`, `protected-inventory.sh`, `stage-exact-deletion.sh` (9-11 lines each, all `exec … do-work-cli.sh`).
- Two rows naming `../../do-work-knowledge/scripts/lexical-memory-recall.sh` and `install-memory-hooks.sh` (6 lines each).
- The sentence `which orchestrates \`tools/checks/uncommitted-inventory.sh\` and \`tools/checks/associate-files.sh\`` is false at HEAD: delete it.
- The shims themselves are not touched (retained by the 0.260.1 decision); only the guide's description of them changes.
- Reproduce at dc8a64e3: `awk 'NR>=9 && NR<=22 && /\.sh`/' skills/do-work/docs/prescribed-shell-primitives.md | grep -oE '[^`]*\.sh' | while read -r p; do case "$p" in ../../*) fp="skills/${p#../../}";; *) fp="skills/do-work/$p";; esac; echo "$(wc -l < "$fp" | tr -d ' ') lines $fp $(grep -c do-work-cli "$fp") cli-exec"; done; rg -n 'which orchestrates' skills/do-work/docs/prescribed-shell-primitives.md`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Same guide and ratchet as REQ-554: land after it so the heading and pointer counts are re-baselined once.
- Lock-in limit: shim rows in the executable-homes table: 0 after this REQ (today 9).

## Dependencies
Depends on REQ-554, which already edits this guide and re-baselines the ratchet, so this REQ is a table rewrite only.

## Builder Guidance
Firm: the route column names the CLI subcommand, not a shim.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine shim rows (6-11 lines each, all `do-work-cli` execs) and the orchestration sentence.
**GREEN when:** The table has zero `.sh` rows for mechanics owned by Go, the sentence is gone; the lock-in pins shim rows in that table at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for stale-shell-ownership-prose.*
