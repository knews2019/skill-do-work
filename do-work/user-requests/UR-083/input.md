---
id: UR-083
title: 'Capture all accepted validation findings'
created_at: 2026-08-31T14:19:37Z
requests: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
word_count: 6
---

# Capture All Accepted Validation Findings

## Summary

Capture every issue that received an **Accept** verdict in the preceding `do-work validate-feedback` report. The report adjudicated 22 accepted comments representing 17 underlying root causes. Nine root causes already have exact queue or in-flight homes, and eight uncaptured or independently gated root causes become new REQs in this UR.

The one **Discuss** item—the Go version floor as written for Go 1.23—is intentionally excluded. Exact-toolchain evidence disproved that target, and REQ-427 (Confirm Go Version Floor) separately implemented the confirmed Go 1.25 floor.

## Extracted Requests

| Accepted root cause | Destination |
|---|---|
| Unsupported timestamps can anchor doctor repairs | REQ-434 (Refuse Unsupported Timestamp Ordering Anchors) |
| UR closure can outrun required member archival | REQ-430 (Couple UR Closure to Terminal Member Archival) |
| Documentation rewrites can describe refused moves | REQ-431 (Couple Documentation Rewrites to Their Owning Moves) |
| Consumed scratch can bypass the nonempty-index `--commit` guard | REQ-432 (Enforce the Commit Guard for Consumed Scratch Cleanup) |
| Scratch-only `cleanup --commit` can delete outside the requested commit | REQ-444 (Refuse Untracked Consumed Scratch in Cleanup Commit Mode) |
| Shared atomic publication strips special mode bits | REQ-436 (Audit Special-Mode Preservation in Remaining File Publication) |
| Doctor reports active and terminal working REQs as stuck | REQ-437 (Stop False Stuck-Work Findings for Active and Terminal REQs) |
| Forensics requires data absent from its declared authority | REQ-435 (Complete the Doctor-Forensics Delegation Contract) |
| Filename-only dependency collisions are mislabeled missing | REQ-428 (Preserve Filename-Only Collision Evidence in Dependency Graphs) |
| `caveman` is absent from typed schema projection | REQ-429 (Complete Normalized Schema-Field Projection) |
| Misplaced UR files share an excessive conflict domain | REQ-433 (Split Misplaced UR Partial-Merge Conflicts by Item) |
| Nested repository roots can bypass Git transaction guards | REQ-438 (Refuse Mismatched Git Transaction Roots) |
| Trailing timeline windows anchor after cosmetic padding | REQ-439 (Anchor Trailing Timeline Windows Before Display Padding) |
| Static board publication can delete non-file targets | REQ-440 (Refuse Non-File Static Board Output Targets) |
| HTTP archive validation occurs after public-target replacement | REQ-441 (Validate HTTP Archives Before Publication) |
| Untimed claimed work reserves no forecast time | REQ-442 (Reserve Forecast Time for Claimed Work Without a Parseable Stamp) |
| Slash-containing branches create multi-component archive prefixes | REQ-443 (Keep Git Fallback Archive Prefixes to One Component) |

## Batch Constraints

- Preserve each finding's original severity, concrete replay, evidence, origin, and Surface-cost judgment through implementation and review.
- Apply the Finding-Closure Ratchet: each fix must land with the named regression or exact deletion/simplification proof.
- Prefer the direct remedies recorded by validation. Do not replace them with general workflow DAGs, timestamp recovery machinery, or other broader defensive apparatus.
- Keep exact duplicate comments as provenance, but do not create duplicate implementation work.

## Folded Requests

- REQ-434 (Refuse Unsupported Timestamp Ordering Anchors) — Finding 1: `[P1] Doctor repair can manufacture and commit a future timestamp` because diagnosis-only timestamps still become repair anchors.
- REQ-430 (Couple UR Closure to Terminal Member Archival) — Finding 2: `[P1] Cleanup can archive a UR while leaving one of its active REQs behind` when a member group is refused independently.
- REQ-431 (Couple Documentation Rewrites to Their Owning Moves) — Finding 3: `[P1] Documentation rewrites can point at moves that were refused` because the global rewrite is owned only by the first move group.
- REQ-432 (Enforce the Commit Guard for Consumed Scratch Cleanup) — Finding 4: `[P2] cleanup --commit can delete consumed scratch despite a nonempty index` because the scratch exception swallows the global empty-index failure.
- REQ-436 (Audit Special-Mode Preservation in Remaining File Publication) — Finding 5: `[P2] The special-mode-bit fix currently misses the shared atomic writer` and cleanup move publication.
- REQ-435 (Complete the Doctor-Forensics Delegation Contract) — Finding 7: `[P2] The forensics action requests data its sole authority cannot provide`; prefer deleting unused count requirements unless a consumer earns a typed summary.
- REQ-428 (Preserve Filename-Only Collision Evidence in Dependency Graphs) — Finding 8: filename-only collision targets are labeled missing because nil lookup wins before ambiguity evidence.
- REQ-429 (Complete Normalized Schema-Field Projection) — Finding 9: `caveman` is defined by schema normalization but omitted from the typed request record.
- REQ-433 (Split Misplaced UR Partial-Merge Conflicts by Item) — Finding 10: one occupied misplaced-UR destination blocks safe siblings because the whole directory is one group.
- REQ-434 (Refuse Unsupported Timestamp Ordering Anchors) — Finding 13 repeats the diagnosis-only timestamp-anchor defect with the same mixed offset/supported successor replay.
- REQ-430 (Couple UR Closure to Terminal Member Archival) — Finding 14 repeats the requirement that closure depend on every required member move.
- REQ-431 (Couple Documentation Rewrites to Their Owning Moves) — Finding 15 repeats the requirement that rewrites remain tied to moves that actually apply.
- REQ-433 (Split Misplaced UR Partial-Merge Conflicts by Item) — Finding 17 repeats the per-item conflict-domain requirement.
- REQ-439 (Anchor Trailing Timeline Windows Before Display Padding) — Finding 21 repeats the 95-day drained-board case where 1.9 days of padding make “Last day” entirely empty.

## Full Verbatim Input

> ````text
> `do-work capture-request:` all the accepted issues
> ````

---
*Captured: 2026-08-31T14:19:37Z*
