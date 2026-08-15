---
id: UR-041
title: Capture validated maintainability audit findings
created_at: 2026-08-15T07:13:20Z
requests: [REQ-181, REQ-182, REQ-183, REQ-184, REQ-185, REQ-186, REQ-187, REQ-188]
word_count: 4
---

# Capture Validated Maintainability Audit Findings

## Summary

Capture all eight validated improvement classes from the 2026-08-14 maintainability audit as independent queue requests. Preserve their severity, exact evidence, bounded remedies, surface-cost limits, and finding-closure proofs without capturing the audit's Discuss, Pre-empted, WATCH, or NOT-MEASURED items.

## Extracted Requests

| REQ | Audit priority | Request |
|---|---|---|
| REQ-181 | P1 | README understates install and implementation write boundaries |
| REQ-182 | P2 | Public work and schema vocabularies drift while suites stay green |
| REQ-183 | P2 | Static board generation can publish a mixed three-file bundle |
| REQ-184 | P2 | Live board origin checks have no trusted Host anchor |
| REQ-185 | P2 | JavaScript behavior probes can all skip while the board suite passes |
| REQ-186 | P3 | Required baseline verification executes two child suites twice |
| REQ-187 | P3 | No single local maintainer command proves shell plus both Go modules |
| REQ-188 | P3 | Hotspot output silently drops unavailable tracked paths |

## Batch Constraints

- The canonical evidence source is `do-work/audits/audit-2026-08-14.md`, audited at commit `58eb4f84f408dce1ec9828a07aef0b174930ce34` and independently validated before capture.
- Keep one REQ per validated root-cause class. Do not merge distinct classes because each has a separate reproduction, remedy boundary, and closure ratchet.
- Do not capture the two Discuss choices, WATCH items, Pre-empted items, or NOT-MEASURED gaps as implementation work. In particular, REQ-180 already owns the `Justfile`/`justfile` case mismatch.
- Preserve each finding's surface-cost limit: use the smallest earned lock-in and do not introduce repository-wide generators, generic transaction frameworks, middleware frameworks, browser frameworks, or test-graph registries.
- Treat overlapping `write_set` entries as visibility for serial coordination, not as new dependencies or a safety guarantee.
- Capture only in this invocation. Do not implement the findings until a later explicit `do-work run`.

## Source Asset

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows the audit report's validated-findings list with filters for all eight classes: one P1, four P2, and three P3. Each row shows its ordinal, severity, title, impact, effort estimate, and collapsed evidence control. The first row is keyboard-focused; no finding body is expanded.

## Full Verbatim Input

do-work capture-request for these

---
*Captured: 2026-08-15T07:13:20Z*
