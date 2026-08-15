---
id: REQ-194
title: Forward stray REQs through verify and forensics
status: pending
created_at: 2026-08-15T09:12:04Z
user_request: UR-043
domain: backend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-193]
maintenance: false
related: [REQ-193]
batch: closed-ur-documentation-hardening
write_set: [skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go, skills/do-work/actions/forensics.md, skills/do-work/actions/work-reference.md]
---

# Forward Stray REQs Through Verify and Forensics

## What

Make forensics surface every REQ file that the board already detects outside `queue/`, `working/`, and `archive/`, reusing the board's detector rather than describing another scan. Align the existing archived-UR/live-member probe with REQ-193 so legitimate review-generated follow-ups do not teach users to reopen a closed UR.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Forensics currently scans only queue, working, and archive when looking for stranded finished REQs. A completed file inside an open UR folder can therefore receive a clean forensics report even though the board correctly warns that it is invisible as a card.

## Context

- Original priority: Secondary; no P-level severity was supplied.
- Verbatim claim: “It has no check for finished ticket files physically parked inside an open user request folder.”
- Evidence: `skills/do-work/actions/forensics.md:67` and `:117` omit `user-requests/`; `skills/do-work-board/tools/queue-kanban/model.go:353` already turns `StrayRequestFiles` into warnings; `board_synthetic_test.go:243` pins the exact misplaced-REQ replay.
- Surface-cost: Earned. Commit `1323982` supplies the incident and regression fixture; forwarding the canonical detector is cheaper than maintaining a second scan.
- The board's derivative dependency warnings are correct consequences of the malformed tree, not false warnings.

## Detailed Requirements

- Give the board's existing stray-request warning a stable structured handoff that `queue-kanban verify` can lift without re-walking the filesystem.
- Emit one read-only, non-fixable verify finding for every REQ file found outside `queue/`, `working/`, and `archive/`, regardless of status.
- Preserve board behavior: stray files remain warnings and must not be parsed into `AllRequests` or rendered as cards.
- Add the minimum parser/schema support needed for verify to recognize the existing `review_generated: true` marker.
- Exempt a queued or working `review_generated: true` member from `ur-archived-with-live-member` when its original UR is already archived; this is REQ-193's legitimate same-UR follow-up shape.
- Keep the archived-UR/live-member finding for ordinary live members, but remove any remedy that recommends moving the archived UR back to `user-requests/`.
- Document the canonical stray-file coverage and review-generated exception in Forensics Check 14 without presenting its probe table as authoritative.
- Add focused verify fixtures for both the stray-file forwarding and the closed-UR review exception.

## Constraints

- Reuse `StrayRequestFiles`/the board warning path; do not implement another directory scan in `forensics.md` or `verify.go`.
- Keep `verify` read-only. Do not mark stray relocation or closed-UR repair mechanically fixable.
- Apply the all-misplaced condition chosen by the user; do not limit detection to terminal statuses.
- Preserve existing duplicate, stranded-finished, and genuine premature-archive findings without double-reporting.
- Do not change the board's three write surfaces.

## Dependencies

Depends on REQ-193, which defines the canonical same-UR/stays-closed lifecycle behavior this diagnostic must recognize.

## Builder Guidance

Firm intent. Reuse the board's canonical detector and add only the smallest marker parse and verify routing needed to distinguish the legitimate review-generated case.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Add `TestVerifyFlagsStrayRequestFiles` with a REQ under `do-work/user-requests/UR-NNN/`, and `TestVerifyAllowsReviewGeneratedMemberUnderClosedUserRequest` with an archived UR plus a queued `review_generated: true` member.
**Why RED now:** The board warns about the stray but `runVerifyProbes` does not forward it, while the archived-UR probe incorrectly flags the legitimate review-generated follow-up and recommends reopening the folder.
**GREEN when:** Both named tests pass: verify reports every stray through the existing detector, stays silent for the legitimate review-generated member, still reports an ordinary live member without recommending unarchive, and `go test ./...` remains green.
**Validation:** User confirmed all-misplaced coverage and the same-UR/stays-closed exception on 2026-08-15; apply `actions/work-reference.md` → Finding-Closure Ratchet.

## Assets

None.

## Full Context

See `do-work/user-requests/UR-043/input.md` for the complete verbatim request, validated evidence, and batch constraints.

---
*Source: documentation-hardening finding validated through `do-work-toolbox validate-feedback`; original priority Secondary.*
