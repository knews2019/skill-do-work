---
id: REQ-497
title: 'Review fix: Strictly normalize frontmatter collision identities'
status: pending
created_at: 2026-09-02T12:03:38Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-452]
maintenance: false
addendum_to: REQ-452
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-452]
---

# Review Fix: Strictly Normalize Frontmatter Collision Identities

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Make repository collision evidence use strict whole-field normalization for frontmatter REQ IDs while retaining suffix-tolerant parsing for filenames. A malformed frontmatter value such as `REQ-452junk` must not alias valid `REQ-452` or veto its documented explicit-target overrides.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this strict-frontmatter-versus-filename identity root cause.

## Context

Found during the post-remediation review of REQ-452. Its first remediation correctly unified numeric-equivalent IDs such as `REQ-452` and `REQ-0452`, but reused `requestIDFromText`, whose suffix tolerance is appropriate for filenames and too permissive for frontmatter identity.

## Requirements

- Parse frontmatter IDs with the selector's strict whole-value numeric grammar.
- Keep filename-derived claims suffix-tolerant so existing repository discovery behavior remains intact.
- Do not report a collision between valid `REQ-452` and malformed frontmatter `REQ-452junk`.
- Preserve collision evidence for genuine numeric equivalents such as `REQ-452` and `REQ-0452`, including every contributing path.
- Preserve dependency, assignment, and impact overrides for a unique explicit target.

## Red-Green Proof

**RED prompt/case:** Create unrelated queue filenames whose frontmatter values are `REQ-452` and `REQ-452junk`; discover the repository and explicitly select `REQ-452` in both discovery orders.
**Why RED now:** Collision normalization uses the suffix-tolerant filename parser for frontmatter claims, so the malformed value aliases the valid numeric identity and returns `DEPENDENCY-AMBIGUOUS`.
**GREEN when:** Repository collision evidence does not merge the malformed value with `REQ-452`, explicit selection chooses the unique valid record in both orders with its documented overrides, and the genuine `REQ-452`/`REQ-0452` collision regressions remain green.
**Validation:** Important post-remediation review finding from REQ-452; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/working/REQ-452-refuse-ambiguous-explicit-request-ids.md` until archival, then the corresponding REQ-452 archive record.

---
*Source: review of REQ-452.*
