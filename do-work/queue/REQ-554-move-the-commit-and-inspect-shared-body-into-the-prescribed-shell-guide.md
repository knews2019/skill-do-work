---
id: REQ-554
title: '[impact-negligible] Move the 46 lines commit.md and inspect.md share into the prescribed-shell guide'
status: pending
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/audit-lockins.sh]
---

# Move the 46 lines commit.md and inspect.md share into the prescribed-shell guide

## What
`skills/do-work/actions/commit.md` and `skills/do-work-toolbox/actions/inspect.md` share 46 byte-identical non-blank lines in runs of three or more: the M/A/D/X/XD classification legend, four file-reading bullets, and two complete "If the script is missing or will not run, do it by hand" fallbacks that restate the algorithms `internal/corehelpers/inventory.go` and `associate-files` implement. Move the legend, the two fallbacks, and the bullets into one section of `skills/do-work/docs/prescribed-shell-primitives.md` (which already owns the protected-inventory heading), cite it from both actions, and keep only the read-only-versus-staging deltas local.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Two prose files stating one rule drift; this is the only non-crew-member prose pair with duplicated three-line windows in the audited surface, and no INLINE RESIDUE row in `decisions/audits/2026-08-11-prescribed-shell-primitives.md` covers it.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 6, sweep_key `commit-inspect-shared-body`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -60. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `commit.md` and `inspect.md` — `If the script is missing or will not run, do it by hand from the complete NUL-delimited output` (207 words each, two relative-path fixups differ).
- `commit.md` and `inspect.md` — `If the script is missing or will not run, do it by hand: glob both directories` (byte-identical).
- `commit.md` and `inspect.md` — `- **Modified files**: Read the \`git diff\` for each file.` and the three bullets after it.
- `commit.md` and `inspect.md` — the `M`/`A`/`D`/`X`/`XD` legend and `Secret-shaped matching is case-insensitive` (X/XD suffixes reworded for read-only mode: keep that delta local).
- The canonicalization ratchet counts headings in the guide and pointers to it from named files: re-baseline those counts in `_dev/tests/prescribed-shell-canonicalization.sh` in the same commit.
- Reproduce at dc8a64e3 (prints 46, then the four fallback sites): `python3 -c "import difflib;a=[l.rstrip() for l in open('skills/do-work/actions/commit.md')];b=[l.rstrip() for l in open('skills/do-work-toolbox/actions/inspect.md')];print(sum(1 for i,j,s in difflib.SequenceMatcher(None,a,b).get_matching_blocks() if s>=3 for k in range(s) if a[i+k].strip()))" && rg -n 'If the script is missing or will not run' skills --glob '*.md'`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (create it on first use, executable, invoked from `_dev/tests/contract-regressions.sh` in the fast tier the way `_dev/tests/defensive-surface-audit.sh` is, with the same missing-or-not-executable FAIL line), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Prime `_dev/primes/prime-action-files.md` before editing either action; the guide edit follows `_dev/primes/prime-shell-commands.md`.
- Lock-in limit: commit.md/inspect.md shared lines ≤ 10 after this REQ (today 46); fallback sentences in actions: 0 (today 4).

## Dependencies
No dependency. REQ-555 (rewrite the guide's executable-homes table) depends on this REQ because it re-baselines the same guide and ratchet counts.

## Builder Guidance
Firm on one home for the shared text; latitude on the section title and on how the read-only delta is phrased.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints 46 shared lines and four fallback sites.
**GREEN when:** Shared lines are at most 10 and the fallback sentence appears 0 times in actions; the lock-in pins the difflib count at the post-fix value and the fallback sentence count at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for commit-inspect-shared-body.*
