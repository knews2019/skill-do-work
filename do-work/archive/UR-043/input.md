---
id: UR-043
title: Harden closed-UR review and forensics contracts
created_at: 2026-08-15T09:12:04Z
requests: [REQ-193, REQ-194]
word_count: 87
---

# Harden Closed-UR Review and Forensics Contracts

## Summary

File two validated findings as one documentation-hardening report rather than a tool bug: prevent standalone review from teaching an agent to reopen an archived User Request, and make forensics reuse the board's existing misplaced-REQ detector instead of duplicating a filesystem scan in prose.

## Extracted Requests

| REQ | Priority | Request |
|---|---|---|
| REQ-193 | Primary | State and pin the closed-UR invariant at standalone review's archived-input read and archive path |
| REQ-194 | Secondary | Forward the existing board stray detector through verify/forensics and align its closed-UR diagnostic |

## Batch Constraints

- Classify the work as documentation hardening, not a queue-kanban defect. The board's existing stray warning is correct and is the detector that exposed the malformed tree.
- Keep the two findings in one UR but as two REQs. REQ-194 depends on REQ-193 because its diagnostic behavior must follow the lifecycle contract established there.
- A standalone-review follow-up keeps the original `user_request`, goes into `do-work/queue/`, and does not move or reopen the original archived UR folder. A completed follow-up is archived into that already-closed folder in place.
- Forward all REQ files found outside `queue/`, `working/`, and `archive/`, regardless of status; every such file is invisible as a board card.
- Reuse the board's existing detector through `verify`; do not introduce an independent prose-described filesystem scan.
- Preserve the finding provenance and surface-cost evidence. Commit `1323982` and `TestSyntheticStrayRequestFlagged` are the verified replay case for the misplaced-file class.
- The downstream `sa2-sentence-aligner2` installation will receive the fix through the normal suite update path; do not edit that repository directly.
- Capture only in this invocation. Do not implement these REQs until a later explicit `do-work run`.

## Full Verbatim Input

do-work capture-request: Documentation-hardening report, not a tool bug. Primary: “Missing: any statement that the archived UR folder stays put.” Add the immutable-archive rule beside review-work.md Step 3’s archived-input read; evidence: review-work.md:60/412, abandon.md:129, capture.md:156; surface-cost earned by commit 1323982’s replay case; test the co-located contract. Secondary: “It has no check for finished ticket files physically parked inside an open user request folder.” Reuse the board’s existing stray detector through verify/forensics rather than adding a duplicate prose scan; evidence: forensics.md:67/117, model.go:353, board_synthetic_test.go:243; surface-cost earned; add a focused verify fixture.

---
*Captured: 2026-08-15T09:12:04Z*
