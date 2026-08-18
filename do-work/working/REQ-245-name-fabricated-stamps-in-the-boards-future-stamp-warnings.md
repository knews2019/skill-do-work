---
id: REQ-245
title: Name fabricated stamps in the board's future-stamp warnings
status: claimed
created_at: 2026-08-18T12:28:33Z
user_request: UR-055
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-244]
batch: timestamp-stamping-integrity
effort_estimate: trivial
write_set: ["skills/do-work-board/tools/queue-kanban/model.go", "skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/*_test.go"]
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T12:43:06Z
  basis:
    - trivial short-circuit
claimed_at: 2026-08-18T12:43:06Z
route: A
---

# Name Fabricated Stamps in the Board's Future-Stamp Warnings

## What

The board's future-stamp diagnosis messages name exactly one cause — "likely local wall-clock time stamped with a Z suffix" — but a fully fabricated value is a second, now-observed cause, and the current wording sends that reader to the wrong fix. Reword the diagnosis clauses to name both causes; keep the fix instruction (rewrite with the current UTC instant per the Timestamp rule) unchanged, since it is correct for both.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

Sibling messages found at capture — update together so they don't drift:

- `skills/do-work-board/tools/queue-kanban/model.go:379` — generate-time data warning: "…likely local wall-clock time stamped with a Z suffix; fix: rewrite with the current UTC instant…"
- `skills/do-work-board/tools/queue-kanban/model.go:1232` — reversed-span message: "…one stamp is usually local wall-clock time written with a Z suffix…"
- `skills/do-work-board/tools/queue-kanban/verify.go:371` — verify-time future `claimed_at`: "…usually local wall-clock time written with a Z suffix"

Comments asserting the single-cause story (e.g. `timestamp_test.go:42`, `completion_anomaly_test.go:227`, `model.go:1338`) should be brought in line where they would otherwise contradict the new wording.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this change — versioning, parser lock-step, build outputs. Read it before touching the tool.
- Message-text change only: no new checks, no threshold changes, the 2-minute skew allowance stays as is.
- Finding provenance (validate-feedback triage, this session): verdict Accept; Surface-cost N/A — text accuracy fix to an existing warning, no new surface.

## Red-Green Proof

**RED prompt/case:** A Go test asserting the future-stamp warning message names fabrication as a possible cause (alongside the wall-clock/Z-suffix cause) fails against the current strings.
**Why RED now:** All three diagnosis messages assert the timezone cause alone; a fabricated stamp — the observed incident — is misdiagnosed by the rendered warning.
**GREEN when:** The three messages name both causes with the fix instruction unchanged, the new assertion passes, and `go test ./...` in the tool directory exits 0.
**Validation:** Inferred during capture

## Full Context

See `do-work/user-requests/UR-055/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding 3 — "Broaden the board's future-stamp warning text: 'local wall-clock time with a Z suffix' is one cause; a fully fabricated value is a second, now-observed one, and the current message sends the reader to the wrong fix."*
