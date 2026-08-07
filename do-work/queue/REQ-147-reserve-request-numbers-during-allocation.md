---
id: REQ-147
title: "Addendum: reserve request numbers during allocation"
status: pending
created_at: 2026-08-07T19:15:15Z
user_request: UR-033
addendum_to: REQ-072
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-134]
maintenance: false
---

# Addendum: Reserve Request Numbers During Allocation

## What

Change `queue-kanban next-req` from a read-only maximum scan into an atomic reservation operation so every successful invocation owns a distinct REQ number before it prints that number.

## Context

REQ-134 was briefly observed as a duplicate while UR-031 and UR-032 were captured concurrently. The committed metadata now has a unique mapping—UR-032 owns REQ-134 and UR-031 owns REQ-135 through REQ-146—but the collision exposed a real allocator defect: calculating `max(existing REQ)+1` does not reserve the result, so two callers can receive the same identifier before either creates its request file.

## Prior Implementation

REQ-072 introduced `next-req` in commit `5db22ea` as a read-only allocator that scans existing request filenames and frontmatter. That implementation deliberately avoided writes, which means sequential or concurrent calls made before request creation can return the same number.

## Requirements

- Atomically reserve a REQ number before returning it from `queue-kanban next-req`.
- Count both existing request records and prior reservations when selecting the next number.
- Ensure sequential calls and concurrent processes receive distinct numbers without relying on a long-lived global lock.
- Use an exclusive filesystem operation so competing allocators retry safely instead of overwriting another reservation.
- Keep abandoned reservations as accepted gaps; never recycle a number merely because its request file has not appeared.
- Keep command output semantics unchanged: on success, print exactly one decimal request number.
- Keep reservation paths contained under the repository metadata tree and reject unsafe symlink or path-escape conditions.
- Update the capture workflow contract so a successful capture stages its reservation marker with the UR and REQ records.
- Add regression tests for sequential allocation, concurrent allocation, existing reservations, and unsafe reservation paths.
- Run `go test ./...`, `go vet ./...`, and the relevant contract regression checks.

## Constraints

- Limit implementation changes to the queue-kanban allocator and directly coupled capture contract.
- Do not modify `contact_processor`, generated assets, unrelated skill documentation, cards, warnings, or pipeline fields.
- Do not introduce a changelog write or change any other queue-kanban command semantics.

## Red-Green Proof

**RED prompt/case:** Invoke `queue-kanban next-req` twice before writing a request, and run multiple allocator processes concurrently against the same repository.

**Why RED now:** A maximum-only scan has no durable claim step, so callers can observe identical state and all print the same next number.

**GREEN when:** Every successful call leaves a durable, exclusively created reservation; sequential and concurrent calls return unique numbers; existing REQs and markers are both honored; unsafe marker paths fail closed; and the Go plus capture-contract verification suites pass.

**Validation:** User explicitly required the Go allocator to reserve numbers so the next call receives a different ID.

## Assets

None.

---
*Source: the go app should reserve the numbers, so the next call gets a different id*
