---
id: REQ-184
title: Live board origin checks have no trusted Host anchor
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: security
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-183, REQ-185, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [skills/do-work-board/tools/queue-kanban/serve.go, skills/do-work-board/tools/queue-kanban/testing_api.go, skills/do-work-board/tools/queue-kanban/testing_test.go, skills/do-work-board/tools/queue-kanban/filementions_test.go]
---

# Live Board Origin Checks Have No Trusted Host Anchor

## What

Anchor live-board request authority to the actual configured listener so matching request-controlled `Origin.Host` and request `Host` values cannot reach the handler unless that authority is accepted by the listener policy.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The live board compares `Origin.Host` only with request `Host` and validates neither against the actual listener. A request with Host and Origin set to `example.com` reaches all six routes under a loopback peer despite the documented hostile-page protection.

## Context

- Audit priority: P2; impact 3; effort normal.
- Root-cause key: `live-board-configured-host-anchor`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 4.
- Reproduce: `sed -n '75,107p' skills/do-work-board/tools/queue-kanban/testing_api.go && sed -n '66,90p;492,545p' skills/do-work-board/tools/queue-kanban/serve.go && sed -n '625,633p' skills/do-work-board/tools/queue-kanban/testing_test.go && sed -n '163,169p' skills/do-work-board/tools/queue-kanban/filementions_test.go`.

## Detailed Requirements

- After bind, wrap the production handler with one normalized authority validator derived from the actual listener and port.
- Define accepted behavior for concrete bind addresses, loopback names/addresses, IPv6 authorities, wildcard binds, and intentional LAN access.
- Reject unconfigured matching Host/Origin values before they reach any of the six inner read or testing-write routes.
- Retain the existing loopback-peer, content-type, Origin, path, and write-scope controls.
- Add focused listener/route table tests for accepted and rejected authorities.

## Constraints

- Do not add authentication or a general middleware framework.
- Preserve supported loopback and explicitly configured LAN behavior.
- Lock-in limit: zero requests with an unaccepted normalized Host reach the inner handler.

## Dependencies

None.

## Builder Guidance

Firm security intent. One outer validator plus focused tests is the earned surface; use the actual post-bind listener authority rather than a second request-controlled comparison.

## Open Questions

None. Listener cases must be made explicit during the PLAN phase using current CLI behavior and existing tests.

## Red-Green Proof
**RED prompt/case:** Under a safe loopback peer, send a request with both `Host` and `Origin` set to `example.com` to every live-board route.
**Why RED now:** Matching request-controlled values pass the current check, and the actual listener authority is never consulted.
**GREEN when:** Every read and testing-write route rejects the unconfigured `example.com` replay before the inner handler; configured concrete, loopback, IPv6, wildcard/LAN cases behave according to explicit tests.
**Validation:** Inferred during capture from the independently validated audit reproduction and lock-in proposal.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 04, labeled P2, impact 3, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 4 in the canonical audit.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
