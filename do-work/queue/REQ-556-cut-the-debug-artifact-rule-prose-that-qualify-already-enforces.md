---
id: REQ-556
title: '[impact-negligible] Cut the debug-artifact rule prose that do-work-cli qualify already enforces'
status: pending
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, skills/do-work/actions/work-reference.md, _dev/tests/audit-lockins.sh]
---

# Cut the debug-artifact rule prose that do-work-cli qualify already enforces

## What
The debug-artifact and P-A-U-honesty rule that `do-work-cli qualify` enforces at Step 6.3 (finding codes `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED`, `QUALIFY-UNIFY-DISARMED` in `internal/corehelpers/checks.go`) is written a second time as an agent instruction at five prose sites across `work.md`, `review-work.md` and `work-reference.md`. Keep one sentence in `work.md` Step 6.3 naming the three finding codes; cut the other four sites to a pointer. Keep the judgment prose the Go check explicitly defers ("judge entry-point or dynamic-wiring exceptions"), which has no duplicate.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Prose that restates a rule code already enforces is read on every REQ by the three highest-churn action files (272, 93 and 162 commits in twelve months) and can drift from the code. The audit labelled this class INFERRED: the code enforces the rule, and whether reviewer prose is meant as a second independent read is the builder's judgment to confirm before cutting.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 1, sweep_key `qualify-debug-artifact-prose-restated`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -15. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `work.md` — Red Flags row `diff contains \`console.log\`, \`debugger\`, or \`TODO\`` ⇐ `QUALIFY-DEBUG-ARTIFACT`.
- `work.md` — Common Rationalizations row `A checked \`[UNIFY]\` over a diff containing \`console.log\`` ⇐ `QUALIFY-UNIFY-DISARMED`.
- `work-reference.md` — `Read the actual diff for debug artifacts` ⇐ `QUALIFY-PAU-UNCHECKED`. If queued REQ-510 (sweep work-reference sections owned by CLI tests) has already removed this site when this REQ is claimed, skip it; shared files are not a dependency.
- `review-work.md` — `Builder checked all P-A-U boxes but the diff contains` and `Diff hygiene — no debug artifacts — console.log/print lines` ⇐ the same `checks.go` pair.
- Contrast (house style, not an instance): `work.md` says of the blocked probe "this supersedes scripts/run-blocked-check.sh ... prose must not execute the probe a second time"; `capture.md` names the SessionStart hook as owner instead of restating it.
- Reproduce at dc8a64e3: `rg -n -e 'console\.log' -e 'debug artifacts' skills/do-work/actions/work.md skills/do-work/actions/work-reference.md skills/do-work/actions/review-work.md && rg -n 'QUALIFY-DEBUG-ARTIFACT|QUALIFY-PAU-UNCHECKED|QUALIFY-UNIFY-DISARMED' skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- `_dev/tests/contract-regressions.sh` may pin some of these sentences; delete the matching predicates in the same commit rather than keeping a sentence to satisfy a pin.
- Prime `_dev/primes/prime-action-files.md` first.
- Lock-in limit: debug-artifact rule mentions across work.md, review-work.md, work-reference.md: ≤ 3 after this REQ (today 9).

## Dependencies
No dependency. Overlaps queued REQ-510 on `work-reference.md`; overlap is not a dependency, the builder checks the site at claim time.

Verify repair (2026-09-03, `do-work verify-requests` on UR-105): the audit's plan line said "fold the work-reference.md site into REQ-510"; this REQ deliberately owns all five sites instead, because REQ-510 is last in the nine-deep REQ-502 chain and a fold would park one sentence behind eight REQs. The maintainer accepted that shape.

## Builder Guidance
Mixed: firm that one sentence naming the finding codes stays in `work.md` Step 6.3; latitude to keep a reviewer-side sentence if reading the code shows the review pass runs on a diff `qualify` never saw.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine mentions across the three files (work.md 5, review-work.md 3, work-reference.md 1) and the three finding codes in checks.go.
**GREEN when:** At most three mentions remain across the three files; the lock-in pins the mention count at the post-fix value.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for qualify-debug-artifact-prose-restated.*
