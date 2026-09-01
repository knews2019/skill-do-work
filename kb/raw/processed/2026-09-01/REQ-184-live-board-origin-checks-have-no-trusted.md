---
source_type: req_lesson
req_id: REQ-184
req_path: do-work/archive/UR-041/REQ-184-live-board-host-anchor.md
date: 2026-08-15
domain: security
module: _dev/primes
tags: [security, live, board, origin, checks]
---

# Lessons from REQ-184: Live board origin checks have no trusted Host anchor

## What the REQ was about

Anchor live-board request authority to the actual configured listener so matching request-controlled `Origin.Host` and request `Host` values cannot reach the handler unless that authority is accepted by the listener policy.

## Solution summary

**Behavior:** Production requests are now authorized against the configured bind and actual post-bind listener port before the inner router runs. Wildcard binds accept only the accepted connection's concrete numeric local address, with loopback aliases limited to loopback connections; arbitrary DNS Host values never inherit wildcard authority.

## What worked

- When HTTP Host is an authority boundary, matching it to Origin is not enough because both are request-controlled; accepted authorities must be derived after bind from the listener and, for wildcard sockets, from the accepted connection's concrete local address.
- Wildcard network reachability does not justify wildcard DNS authority. Numeric local destinations preserve intentional LAN access without reopening DNS-rebinding-style Host trust.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Back-reference

See `do-work/archive/UR-041/REQ-184-live-board-host-anchor.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6711636`.
