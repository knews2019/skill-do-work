---
id: UR-092
title: 'Canonicalize REQ reservation marker filenames across allocation flows'
created_at: 2026-09-01T12:11:03Z
requests: [REQ-485]
---

# Canonicalize REQ Reservation Marker Filenames Across Allocation Flows

The REQ-number collision guard is defeated by a marker filename mismatch. queue-kanban's `next-req` creates zero-padded markers (`do-work/.req-reservations/REQ-000482`), while capture-files manifests in the wild carry unpadded marker paths (`do-work/.req-reservations/REQ-482`), and the exclusive-create semantics that are supposed to make concurrent reservations collide never fire across the two spellings. Both flows scan the same max and therefore compute the same candidate number, and both "win".

This happened for real on 2026-09-01: this session's `next-req` returned 482 and created `REQ-000482` at 11:50Z; a concurrent capture created `REQ-482` and committed `do-work/queue/REQ-482-stack-verify-findings-full-width.md` (commit 78847fe4) at 11:56Z. Only the orchestrator's manual git-status inspection caught the duplicate before a second REQ-482 file was committed; the fix was a hand renumber to REQ-483.

Fix: one canonical marker filename convention shared by every allocation flow — queue-kanban `allocate.go`, the capture-files manifest guidance in `skills/do-work/actions/capture.md`, and `skills/do-work/scripts/cleanup-req-reservations.sh`'s reaping match — with read-side acceptance of both legacy spellings so existing markers still count in the max-scan and still reap. Add a lock-in test that reproduces today's collision shape: two flows reserving the same number through different spellings must collide, not both succeed. The board parser lock-step rule (_dev/primes/prime-kanban-board.md) governs if the board reads marker names anywhere.
