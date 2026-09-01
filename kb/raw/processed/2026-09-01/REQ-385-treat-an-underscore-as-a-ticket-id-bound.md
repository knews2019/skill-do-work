---
source_type: req_lesson
req_id: REQ-385
req_path: do-work/archive/UR-075/REQ-385-treat-underscore-as-a-ticket-id-boundary.md
date: 2026-08-27
domain: frontend
module: _dev/primes
tags: [frontend, treat, underscore, ticket, boundary]
---

# Lessons from REQ-385: Treat an underscore as a ticket-id boundary on both surfaces

## What the REQ was about

`\b` counts `_` as a word character, so the mention pattern behaves differently around an underscore
than a reader expects. Change the boundary on both the Go and the client side together, in one
commit, so the agreement test stays green.

## Solution summary

Underscores now delimit ticket IDs. Boundary checks run after candidate matching, preventing the regular expression from falling back to a shorter UR alternative. Ticket candidates remain consumed even when resolution intentionally suppresses them; only non-ticket runs keep the previous retry behavior.

## What worked

Compound-first alternation does not guarantee compound-first behavior when a failing boundary permits fallback. Consume before checking boundaries, and keep intentionally suppressed ticket candidates consumed too: a caller's retry can recreate a fallback the regular expression no longer performs.

## Back-reference

See `do-work/archive/UR-075/REQ-385-treat-underscore-as-a-ticket-id-boundary.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `259b1479`.
