---
id: REQ-182
title: Public work and schema vocabularies drift while suites stay green
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: testing
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-183, REQ-184, REQ-185, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [README.md, skills/do-work/SKILL.md, skills/do-work/docs/work-guide.md, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work-board/tools/queue-kanban/testing.go]
---

# Public Work and Schema Vocabularies Drift While Suites Stay Green

## What

Restore parity at the public work-guide/router and testing-schema/normalizer seams, and correct the two short workflow summaries that omit canonical states while the baseline suites remain green.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The public work guide advertises three aliases the first-match router cannot dispatch, workflow glosses omit canonical fields/statuses, and the board accepts two testing-status aliases absent from the schema table. These are user-visible or maintainer-facing drifts that existing suites do not detect.

## Context

- Audit priority: P2; impact 3; effort normal.
- Root-cause key: `public-vocabulary-parity`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 2.
- Reproduce: `rg -n 'do-work (begin|execute|build)|Other trigger words|enum/boolean fields|Queue: N pending|blocked-dependency-cycle|selected for testing|returned with feedback' README.md skills/do-work/SKILL.md skills/do-work/docs/work-guide.md skills/do-work/actions/work.md skills/do-work/actions/work-reference.md skills/do-work-board/tools/queue-kanban/testing.go`.

## Detailed Requirements

- Restore the pre-modular `do-work begin`, `do-work execute`, and `do-work build` aliases advertised by `skills/do-work/docs/work-guide.md`, or otherwise make the advertised list and first-match router exactly agree.
- Keep one public work-alias list and compare it with the router so additions or removals cannot occur on one side only.
- Delete or generalize the stale `enum/boolean fields` gloss in `skills/do-work/actions/work.md`.
- Make queue summaries include dependency-cycle holds rather than implying only `N pending`.
- Add the already-supported `selected for testing` and `returned with feedback` aliases from `testing.go` to the canonical `testing_status` schema table.
- Add compact bidirectional parity checks with mutation cases at both declared seams.

## Constraints

- Do not create a repository-wide prose/schema generator.
- The lock-in limit is zero one-sided aliases across the work-guide/router and testing-schema/normalizer seams.
- Keep runtime alias support and public documentation synchronized without widening unrelated vocabulary.

## Dependencies

None. This REQ is semantically independent, though its documented `write_set` overlaps REQ-181.

## Builder Guidance

Firm intent. The audit attributes the two parity checks to incidents `9ba534e` and `ea0fd94`; keep the assertions compact and seam-local.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Mutate either the work guide or router by one alias, and mutate either the testing schema table or normalizer by one alias; current baseline checks remain green and the two sides disagree.
**Why RED now:** Three documented work aliases are absent from the router, two runtime testing aliases are absent from the schema table, and two workflow summaries omit canonical states.
**GREEN when:** Seam-local tests fail for either one-sided addition or removal; all six current alias mismatches and both stale summary glosses are corrected.
**Validation:** Inferred during capture from the audit's exact instances and lock-in proposal.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 02, labeled P2, impact 3, normal effort, among the eight validated audit findings.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 2 in the canonical audit.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
