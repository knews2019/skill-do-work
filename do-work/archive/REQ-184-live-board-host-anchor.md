---
id: REQ-184
title: Live board origin checks have no trusted Host anchor
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T10:29:10Z
completed_at: 2026-08-15T10:51:51Z
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
route: C
kb_status: pending
kb_entry:
---

# Live Board Origin Checks Have No Trusted Host Anchor

## What

Anchor live-board request authority to the actual configured listener so matching request-controlled `Origin.Host` and request `Host` values cannot reach the handler unless that authority is accepted by the listener policy.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Wrap the production handler only after a successful bind, deriving the actual port and concrete listener address from `listener.Addr()`. Normalize request authorities, accept only configured concrete names/addresses or the wildcard connection's actual local numeric address (plus loopback aliases on loopback), reject unconfigured authorities before routing, and normalize the testing-write Origin check against the same accepted HTTP authority.
- [x] **[APPLY]:** Added the hostile six-route RED matrix first, then implemented the post-bind authority wrapper, concrete/wildcard normalization policy, production wiring, normalized HTTP Origin comparison, and focused `/file` integration within the four-file scope.
- [x] **[UNIFY]:** Reviewed the complete four-file diff and normalization tables. Verified actual `:0` port anchoring, DNS/IP/IPv6 forms, concrete and wildcard/LAN policy, missing/wrong/malformed authorities, all six routes, normalized Origin, existing peer/path/write guards, `go test -count=1 ./...`, `go vet ./...`, and `git diff --check`; no debug artifacts remain.

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
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 04, labeled P2, impact 3, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 4 in the canonical audit.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*

## Triage

**Route C** — the security failure is reproduced, but the accepted-authority policy must be derived from post-bind listener behavior across concrete, loopback, IPv6, wildcard, and intentional LAN configurations before implementation.

## Plan

1. Add a production-only outer handler constructor in `serve.go` that receives the configured bind spelling, actual post-bind listener address, and inner live-board handler.
2. Normalize Host authorities as DNS/IP plus numeric port; use the actual assigned port, canonical DNS/IP forms, bracketed IPv6, and no request-time DNS lookup.
3. For concrete binds, accept the explicitly configured hostname/address and the actual listener address; apply loopback aliases only to loopback binds. For wildcard binds, accept the connection-local numeric IP, and loopback aliases only when the accepted connection is loopback; reject arbitrary DNS and wildcard literals.
4. Install the wrapper before all route dispatch. Retain all existing peer, method, content-type, path-containment, and write-scope guards.
5. Normalize write Origin authority and require the live board's `http` scheme while preserving absent/`null` behavior.
6. Prove hostile matching Host/Origin values reach zero inner routes, and table-test actual port 0 resolution, concrete/loopback/LAN/IPv6/wildcard accepted and rejected cases.

## Exploration

- `runServeCommand` currently installs `liveBoardServer` directly in `http.Server.Handler` after `bindServeListenerAndAnnounce`; this is the single production integration point.
- `liveBoardServer.ServeHTTP` dispatches six routes before any Host validation: `/`, `/board-data.js`, `/board-markdown.js`, `/file`, `/api/testing/profile`, and `/api/testing/status`.
- Testing writes compare only request-controlled `Origin.Host` and `Request.Host`; `/file` and writes have independent loopback-peer guards that must remain behind the new outer validator.
- `resolveServeListenAddress` preserves explicit host forms and turns bare ports into loopback binds. `listener.Addr()` is the only authority for the assigned port when configured with `:0`.
- Go supplies the accepted connection's local address through `http.LocalAddrContextKey`, allowing wildcard listeners to preserve numeric LAN access without trusting arbitrary DNS labels.
- Existing direct `httptest.NewServer(newLiveBoardServer(...))` unit tests intentionally bypass production composition; focused wrapper fixtures should exercise the new boundary without rewriting unrelated route tests.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/serve.go` (modify) — authority normalizer/policy, outer wrapper, and post-bind production wiring
- `skills/do-work-board/tools/queue-kanban/testing_api.go` (modify) — normalized HTTP same-origin comparison for testing writes
- `skills/do-work-board/tools/queue-kanban/testing_test.go` (modify) — listener/authority/route/write security tables and RED/GREEN proof
- `skills/do-work-board/tools/queue-kanban/filementions_test.go` (modify) — accepted and rejected wrapped `/file` integration coverage

**Files I will NOT touch:** CLI flags, authentication/session design, unrelated board routes, or files outside the declared write set.

**Acceptance criteria (restated from REQ):**
- [ ] Every unaccepted normalized Host is rejected before any of the six inner routes.
- [ ] The trusted port comes from the actual post-bind listener, including `:0`.
- [ ] Concrete, loopback, IPv6, wildcard, and intentional LAN authorities follow explicit tested rules.
- [ ] Existing loopback-peer, content-type, Origin, path, and write-scope controls remain effective.
- [ ] No authentication or general middleware framework is introduced.

## Pre-Flight

**Git:** ⚠ Four pre-existing edits under `do-work/queue/` (REQ-189–192) belong to other work and must remain unstaged; all four scoped implementation files are clean.
**Tests baseline:** ✓ `go test -count=1 ./...` passed before implementation (5.316s).
**Dependencies:** ✓ Go toolchain and board module dependencies are available.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified) — adds normalized authority parsing, concrete/wildcard listener policy, an outer Host gate, and post-bind production wiring
- `skills/do-work-board/tools/queue-kanban/testing_api.go` (modified) — compares normalized HTTP Origin and request authorities while preserving existing write guards
- `skills/do-work-board/tools/queue-kanban/testing_test.go` (modified) — covers hostile Host/Origin replay on all six routes plus concrete, loopback, LAN, IPv6, wildcard, port, malformed, and Origin cases
- `skills/do-work-board/tools/queue-kanban/filementions_test.go` (modified) — proves accepted `/file` access and hostile authority rejection through production composition

**Behavior:** Production requests are now authorized against the configured bind and actual post-bind listener port before the inner router runs. Wildcard binds accept only the accepted connection's concrete numeric local address, with loopback aliases limited to loopback connections; arbitrary DNS Host values never inherit wildcard authority.

## Testing

**RED:** Before the wrapper existed, a loopback client using matching hostile `Host` and `Origin` values for `example.com` reached all six live-board routes and received HTTP 200.

**GREEN:**
- Focused authority, Origin, six-route, and `/file` tests — PASS
- `cd skills/do-work-board/tools/queue-kanban && go test -count=1 ./...` — PASS (5.297s)
- `cd skills/do-work-board/tools/queue-kanban && go vet ./...` — PASS
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports the four-file Implementation Summary exactly matches Scope; unrelated queue edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found all files in the diff, all P-A-U phases complete, and no debug artifacts.
- **Substance and traceability:** PASS — the validator uses the actual listener-assigned port, normalizes authorities without DNS lookup, applies explicit concrete/loopback/LAN/IPv6/wildcard rules, and rejects before dispatch.
- **Wiring/data flow:** PASS — `runServeCommand` constructs the outer handler only after bind and installs it in `http.Server.Handler`; all six routes and existing inner security guards remain behind that boundary.

## Review

**Result:** Approve — Acceptance: Pass  
**Overall score:** 98%

- **Requirements (100%):** Every firm authority, route, policy, and preservation requirement is delivered.
- **Code quality (98%):** The policy is fail-closed, avoids request-time DNS, canonicalizes IP forms, and stays one outer boundary rather than a framework.
- **Test adequacy (96%):** Strong post-bind, six-route, concrete/wildcard, IPv6, malformed, Origin, and `/file` coverage; a real scoped/link-local IPv6 fixture remains optional because it is platform-specific.
- **Scope (100%):** Exactly the four declared files changed.

**Important findings:** None.  
**Minor findings:** None.  
**Explicit remediation:** None.

**Reviewer scope judgment:** The pre-existing `--port ...:0` banner/browser URL still displaying port 0 is adjacent but not a REQ-184 failure; the firm requirement is that the security gate trust the listener-assigned port, which it does.

## Lessons Learned

- When HTTP Host is an authority boundary, matching it to Origin is not enough because both are request-controlled; accepted authorities must be derived after bind from the listener and, for wildcard sockets, from the accepted connection's concrete local address.
- Wildcard network reachability does not justify wildcard DNS authority. Numeric local destinations preserve intentional LAN access without reopening DNS-rebinding-style Host trust.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.
