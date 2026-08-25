---
id: REQ-363
title: "[impact-user-visible] Exempt supported UR-less REQ shapes from structural damage"
status: cancelled
completed_at: 2026-08-24T13:42:26Z
created_at: 2026-08-24T11:35:00Z
user_request: UR-068
addendum_to: REQ-343
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
---

# Exempt Supported UR-Less REQ Shapes From Structural Damage

## What

REQ-343's structural probe flags any REQ without `user_request` as damaged, carving out only a
`stakeholder:` marker and an `archive/legacy/` path. Two other **documented, supported** shapes carry
no `user_request` by design, and the probe reports both as damage with a remedy that names a UR which
does not exist.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-343's own Detailed Requirements say the probe must distinguish damage from legitimate absence,
because "a probe that flags correct files is a probe someone turns off". Two legitimate shapes were
missed:

**1. `code-review.md`'s REQ template mints no `user_request`.** Verified verbatim in the shipped
template:

```
id: REQ-NNN
title: '[<impact token>] Code review: [brief description]'
status: pending
created_at: <timestamp>
review_generated: true
source: code-review
scope: [prime file or directory that surfaced this]
---
```

It carries `source: code-review` and `scope:` instead. Every REQ a user confirms out of a code review
would be reported damaged.

**2. `actions/work.md`'s archive table routes UR-less REQs to `archive/` root, not `archive/legacy/`.**
Lines 588-589:

| `context_ref` (legacy) | Move REQ to `archive/`. … |
| Neither (standalone legacy, or a stakeholder-questions REQ …) | Move directly to `archive/`. |

REQ-343's `isLegacyArchiveRequestPath` looks for the `archive/legacy/` segment, so a legacy REQ
archived by the documented path lands outside the carve-out.

Neither shape exists in the tree today — that is why REQ-343's real-tree false-positive check came
back byte-identical, and why this shipped. The writers that produce them are live.

## Detailed Requirements

- A REQ minted by the `code-review.md` template produces no `user_request` finding. Key it on the
  documented shape (`source: code-review`, or `review_generated: true` with a `scope:`), not on the
  field's mere absence — the same discipline REQ-343 applied to the stakeholder marker.
- A legacy REQ archived to `archive/` root by the documented path produces no `user_request` finding.
  `context_ref` is the discriminator the schema already defines for that shape.
- The fence, `id` and status checks still apply to all of these. Only the `user_request` finding is
  exempted, exactly as REQ-343 scoped its existing carve-outs.
- **A genuinely damaged REQ must not be able to wear either exemption cheaply.** REQ-343's review
  already flagged that the stakeholder carve-out keys on the marker's presence alone; do not repeat
  that shape. State what each new discriminator requires and why it cannot be produced by damage.
- `verify_test.go` gains a case per shape, plus the mutation that proves each carve-out is
  load-bearing — remove it and the fixture must be flagged.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs. Read it first.
- Do not widen the exemption to "any REQ without `user_request`". That deletes the check.
- Keep the parser's leniency, as REQ-343 required.

## Red-Green Proof

**RED prompt/case:** Build a REQ from `code-review.md`'s template verbatim (no `user_request`,
`review_generated: true`, `source: code-review`) and run `queue-kanban verify`. It reports
`structurally-damaged-req: … carries no user_request: pointer` and exits 1, on a file the shipped
template told an agent to write. Same for a `context_ref` REQ placed in `do-work/archive/` root.

**GREEN when:** both produce no `user_request` finding; a REQ that is genuinely missing the field and
matches neither documented shape is still flagged; and removing either carve-out makes its fixture
fail.

**Validation:** Reported by an external reviewer on PR #166; both halves independently verified
against the shipped template and `actions/work.md`'s archive table before capture.

---
*Source: PR #166 review finding (UR-068) — a false-positive class in REQ-343's structural probe.*

## Cancelled

- **When:** 2026-08-24T13:42:26Z
- **Why:** folded into REQ-357 (sweep `req-343-structural-probe-remediation`) at the maintainer's instruction — both supported-shape exemptions survive as an instance there, reconciled with the maintainer's clarify reversal of D-08 (recognition by schema discriminator, never by path; backfill wherever a real UR exists). Also the decision-revalidation scan's one candidate for that reversal; the reconciliation resolves it
- **Decided by:** user, via `do-work clarify` fold confirmation
