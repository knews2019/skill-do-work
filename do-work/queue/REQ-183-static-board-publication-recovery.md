---
id: REQ-183
title: Static board generation can publish a mixed three-file bundle
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-184, REQ-185, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/generate_test.go]
---

# Static Board Generation Can Publish a Mixed Three-File Bundle

## What

Make static board publication all-or-recover across `board-data.js`, `board-markdown.js`, and `index.html`, so a handled write failure cannot expose files from different generations.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`generateStaticSite` writes the public files directly in sequence. A forced second-write failure exits nonzero after publishing new rendered data beside the old copy payload and shell, leaving a plausible but internally inconsistent shareable artifact.

## Context

- Audit priority: P2; impact 3; effort normal.
- Root-cause key: `static-board-all-or-recover-publication`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 3.
- The canonical audit contains the complete temporary-fixture reproduction command and the observed new-card/old-copy result.

## Detailed Requirements

- Build all three payloads privately before publishing any public target.
- Publish only `board-data.js`, `board-markdown.js`, and `index.html` as one bounded all-or-recover operation.
- Use unique backups/restoration so a handled publication failure restores all pre-invocation target bytes.
- Preserve unrelated entries already present under `--out`.
- Cover seeded failure, first generation, successful replacement, rollback, and private-residue cleanup.
- Preserve the existing static, zero-network board output contract.

## Constraints

- Define no stronger promise than all-or-recover for these three files; do not claim one filesystem-atomic rename across them.
- Do not introduce a general transaction framework.
- Lock-in limit: after a handled failure, zero public target bytes differ from their pre-invocation state and no private residue remains.

## Dependencies

None. The `generate_test.go` overlap with REQ-185 is declared for coordination, not a semantic prerequisite.

## Builder Guidance

Firm intent. The audit says one bounded three-file staging/backup helper is earned by the reproduced failure and `_dev/lessons/validated-runtime-boundaries.md`.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Generate an initial fixture, change the ticket, make the second target unwritable, and generate again.
**Why RED now:** The command exits 1 after `board-data.js` contains the new title while `board-markdown.js` still contains the old body.
**GREEN when:** The seeded second-write failure returns nonzero while all three public targets remain byte-identical to their pre-invocation state, unrelated output entries remain, and no staging or backup residue survives; first-generation and success paths also pass.
**Validation:** Inferred during capture from the audit's executed failure replay and closure ratchet.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 03, labeled P2, impact 3, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 3 in the canonical audit for the exact shell reproduction.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
