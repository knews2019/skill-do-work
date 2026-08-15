---
id: REQ-194
title: Forward stray REQs through verify and forensics
status: completed
created_at: 2026-08-15T09:12:04Z
claimed_at: 2026-08-15T13:10:43Z
completed_at: 2026-08-15T13:42:43Z
commit: ca34ef2
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
route: C
kb_status: pending
kb_entry:
---

# Forward Stray REQs Through Verify and Forensics

## What

Make forensics surface every REQ file that the board already detects outside `queue/`, `working/`, and `archive/`, reusing the board's detector rather than describing another scan. Align the existing archived-UR/live-member probe with REQ-193 so legitimate review-generated follow-ups do not teach users to reopen a closed UR.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Retain the walker-owned structured stray list on `Board` while preserving warnings/no-card behavior; parse only the exact review-generated marker; forward one non-fixable verify finding per retained stray; narrow the archived-UR anomaly to ordinary live members with a stays-closed remedy; document the marker and Forensics Check 14; and lock the caller-level paths with focused RED/GREEN fixtures.
- [x] **[APPLY]:** Added both named verify tests first and captured their distinct REDs, then retained the canonical structured stray evidence, parsed the exact review marker, routed non-fixable findings, narrowed the closed-UR anomaly/remedy, and updated schema/forensics guidance within the exact five-file Scope.
- [x] **[UNIFY]:** Reviewed all five files and corrected the final no-Go forensics restatement. Verified `model.go` retains but never parses strays into cards and reads only exact marker truth; `verify.go` performs no second walk or warning parse and preserves stranded/ordinary-live ownership; `verify_test.go` covers pending/completed strays plus mixed true/false/noncanonical/terminal members; `forensics.md` reports stray coverage unverified without Go; `work-reference.md` documents the marker-only contract. Focused/full uncached tests, vet, gofmt, aggregate contracts, and `git diff --check` pass with no debug artifacts.

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

## Triage

**Route C** — the outcome is firm, but the change crosses the board model, structured verify findings, frontmatter marker parsing, closed-UR exception logic, focused fixtures, and forensics documentation; the exact handoff shape and double-reporting boundaries require a written plan and independent code trace.

## Plan

1. Extend the existing board model with the walker-owned stray records and an exact boolean `ReviewGenerated` marker, without parsing strays into request/card collections or changing warning text.
2. Add a dedicated read-only verify category that consumes only retained structured strays, and narrow the archived-UR live-member probe to skip only non-terminal queue/working members whose marker is exactly true.
3. Replace the obsolete reopen remedy with closed-UR-safe guidance, document the marker contract and Forensics Check 14 routing, and keep failure/stranded/ordinary-live ownership intact.
4. Write the two named caller-level tests first, include mixed ordinary/generated and pending/completed stray shapes, then run focused mutation-sensitive cases, the full Go module, aggregate contracts, and canonical maintainer gate.

## Exploration

- `enumerateDoWorkTree` already produces structured absolute/relative stray records for every detector-supported off-section `REQ-*.md`; `buildBoard` turns them into warnings and then drops the list. Retaining that same list on `Board` is the smallest handoff and avoids both warning-prose parsing and a second filesystem walk.
- Strays must stay outside `AllRequests`, `RequestsById`, cards, calendar, UR membership, dependencies, and duplicate/stranded probes. The existing synthetic board test already pins warning plus no-card behavior.
- `RequestTicket` can parse `review_generated` through the existing scalar coercion path with an exact `== "true"` check. The marker has no aliases, normalization warnings, display, or scheduling role.
- `appendArchivedUserRequestLiveMemberFindings` already excludes terminally resolved queue/working members so the stranded-finished probe owns them. The new exception belongs after that carve-out and skips only marked non-terminal members; mixed URs must still report ordinary siblings.
- The current archived-UR remedy incorrectly recommends moving the UR back to `user-requests/`. The replacement keeps the UR closed and asks the human to resolve, abandon, or correct the ordinary REQ's association.
- Forensics Check 14 is a forwarding consumer, not an alternate detector. Its table should name structured stray coverage and the exact marker exception while remaining an illustrative snapshot.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — retain structured stray records on `Board` and parse the exact review-generated marker
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — emit one non-fixable stray finding per retained path and narrow the archived-UR live-member probe/remedy
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — add the named stray-forwarding and closed-UR exception RED/GREEN fixtures plus adjacent ownership assertions
- `skills/do-work/actions/forensics.md` (modify) — document Check 14's detector forwarding, marker exception, and stays-closed remedy
- `skills/do-work/actions/work-reference.md` (modify) — define `review_generated: true` as the exact marker parsed for the diagnostic exception

**Files I will NOT touch:** `walk.go`, frontmatter infrastructure, board warning/UI/generator behavior, board write surfaces, cleanup, generic scanners, consumer repositories, or contract-regression scripts.

**Acceptance criteria (restated from REQ):**
- [x] Every stray retained by the board's existing detector produces exactly one read-only, non-fixable verify finding regardless of status, without being parsed into cards or other REQ probes.
- [x] Verify consumes structured board evidence and performs no second filesystem scan or warning-prose parse.
- [x] The exact `review_generated: true` marker is parsed with absent/non-true values false and no alias or display behavior.
- [x] A non-terminal queued/working review-generated member under its archived UR is exempt from `ur-archived-with-live-member`; ordinary live members still report and terminal strays remain owned by stranded-finished.
- [x] The archived-UR remedy keeps the UR closed and never recommends moving it back to `user-requests/`.
- [x] Forensics Check 14 documents the canonical stray coverage and narrow exception without making its table authoritative.
- [x] The named focused tests close both captured findings and the full board/canonical suites remain green.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — retains canonical structured stray paths on the board and parses exact `review_generated: true` into a verify-only marker
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — forwards one non-fixable finding per retained stray and exempts only marked non-terminal members from the archived-UR anomaly with a stays-closed remedy
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified) — adds the two named caller-level fixtures and locks status independence, exact marker truth, mixed ordinary members, terminal ownership, and remedy behavior
- `skills/do-work/actions/forensics.md` (modified) — documents structured stray forwarding, the review-generated exception, and explicit unverified coverage when Go is unavailable
- `skills/do-work/actions/work-reference.md` (modified) — defines the exact marker's absent/noncanonical semantics and sole diagnostic consumer

**Behavior:** `queue-kanban verify` now surfaces every request file retained by the board's existing off-section detector as one read-only, non-fixable finding without turning it into a card. An archived UR may keep legitimate queued/working review-generated follow-ups, while ordinary live members still report and the remedy keeps the UR closed.

## Testing

**Finding-closure RED:** Before implementation, `TestVerifyFlagsStrayRequestFiles` found zero structured verify findings, while `TestVerifyAllowsReviewGeneratedMemberUnderClosedUserRequest` reported exact-true review follow-ups and exposed the obsolete remedy that moved the UR back to `user-requests/`.

**GREEN:**
- focused five-test verify/board lane — PASS
- formal review-remediation seam test — PASS: warning-parser and filesystem-rewalk mutants each fail when only `Board.StrayRequestFiles` carries evidence
- `go test -count=1 ./...` — PASS
- `go vet ./...` — PASS
- `gofmt -l model.go verify.go verify_test.go` — PASS (no output)
- `bash _dev/tests/contract-regressions.sh` — PASS after implementation and after the final Forensics no-Go clarification
- `bash _dev/tests/maintainer-verify.sh` — PASS on the final diff, including ShellCheck, aggregate contracts, both Go modules, and the strict JavaScript lane
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports the exact five-file Implementation Summary matches Scope; foreign queue edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found complete P-A-U evidence, the exact modified file set, and no debug artifacts.
- **Substance and traceability:** PASS — every detailed requirement maps to the retained structured evidence, dedicated verify category, exact marker parse, narrow diagnostic carve-out, caller-level tests, or the two canonical documentation seams.
- **Wiring/data flow:** PASS — one filesystem enumeration populates `Board.StrayRequestFiles`; both warning and verify consumers read it while ordinary request/card collections remain separate. Marker data flows only from frontmatter into the archived-UR probe after terminal ownership is resolved.

## Review

**Overall: 99%** | 2026-08-15T13:42:14Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 99% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None. The initial review's sole Important finding—warning parsing could replace the required structured handoff while tests stayed green—was closed in the allowed remediation attempt with a direct structured-only seam test.
**Minor findings:** None.
**Acceptance:** Pass — canonical stray evidence now reaches verify without a second scan or prose parse, and the closed-UR exception recognizes only the exact review-generated shape.
**Suggested testing:** None beyond the canonical maintainer gate and recorded warning-parser/rewalk mutations.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Retaining the walker's structured stray records on the board lets warnings and verify share one detector without turning misplaced files into cards.
- A direct seam test with no warnings or filesystem evidence proves the intended data source more strongly than an end-to-end fixture alone.

**What didn't:**
- The first integration tests allowed verify to reconstruct identical output from warning prose, so the forbidden coupling survived mutation.
- Forensics initially claimed every Go-backed probe had a manual equivalent, contradicting the deliberate no-second-scan boundary for strays.

**Worth knowing:** Apply the archived-UR ownership checks in order: terminal queue/working members remain stranded-finished evidence; among non-terminal members, only exact `review_generated: true` is legitimate under a closed UR. Ordinary siblings must still report.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Orientation

**[MAP CHANGED]** Queue verification now shares the board walker's structured stray-request evidence and understands the closed-UR review-follow-up marker. Misplaced REQs stay invisible as cards but are no longer invisible to forensics, while legitimate same-UR review work does not trigger advice to reopen a closed folder.
