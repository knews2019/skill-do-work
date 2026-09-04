# REQ-503 independent review

**Approve with follow-ups** — the read-only foundation works across the tested lifecycle phases, but two additional acceptance cases expose incorrect phase decisions. Both are report-only findings; neither authorizes a new request.

Reviewed Route C implementation range `d1e87a31d358d27479816f9181555887e71af3f6..f38a78b0ad80b34b3f5cd332b31e21ae63a7602d`, the complete claimed REQ, UR-098, the analysis step table, and the relevant action, CLI prime, shell, backend, general, and coding-guardrail contracts. Later queue/gate/finalization mutations are intentional chain work and are not attributed to this foundation.

## Review

**Overall: 78.75%** | 2026-09-04T22:01:46Z

**Verdict: Approve with follow-ups.** Typed identity, phase, evidence coordinates, next/replay argv, deterministic rendering, command registration, and the read-only boundary are delivered. Acceptance is partial because legitimate Route A completion and fenced Markdown examples are misclassified.

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

Overall calculation: `(90 + 85 + 80 + 100) / 4 - 10 = 78.75`; the ten-point deduction is the documented Partial modifier.

**Important findings:**

- At the reviewed revision, `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:243` requires `Lessons Learned` for every route and refuses an existing `Orientation` without it. `skills/do-work/actions/work.md:627` expressly permits straightforward Route A requests to skip lessons. A Route A fixture with Triage, Plan, Implementation Summary, Qualification, Testing, Review, and Orientation returned exit 1 / `ADVANCE-EVIDENCE-MISSING` instead of finalization-manifest judgment. Preserve that optional Route A path and cover it with a focused regression. This is an incorrect advisory refusal, with the existing prose still providing the legitimate completion path — impact-user-visible → report only
- At the reviewed revision, `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:334` scans every matching `##` line without excluding fenced code. A Route A request awaiting implementation returned `qualify` after adding only a fenced Markdown example containing `## Implementation Summary` beneath What; the control without that example correctly returned implementation judgment. Treat fenced examples as content rather than completed phase evidence and cover both backtick and tilde fences. This misdirects the advisory next step; downstream qualification remains an independent authority — impact-user-visible → report only

**Minor findings:** None.

**Requirements checklist:**

- Route A/B/C normal phase progression and mechanical/judgment separation: delivered, except the optional Route A lessons branch above.
- Typed phase, exact request identity/path, missing-evidence coordinates, next argv and verification argv in text/JSON: delivered and exercised through the real binary.
- Typed malformed/ambiguous/impossible-state refusals: delivered for the covered cases; fenced content can be mistaken for real lifecycle evidence as recorded above.
- Read-only repository/REQ/checkpoint/Git behavior: delivered; the implementation discovers and projects without calling mutation handlers, and the byte-digest test passed in both output formats.
- Prime ownership and route table: delivered. The six production/test/doc files match the expanded Scope; the additional range entry is the REQ's own lifecycle evidence. No action prose or sentence predicates were deleted, as this explicitly scoped foundation requires; subsequent REQs own those moves.

**Acceptance: Partial.** In an isolated detached worktree at the exact target, `go test -count=1 ./internal/lifecycleadvance ./internal/resultmodel ./cmd/do-work-cli` passed (1.486s, 0.700s, 0.502s). These tests invoke the real CLI across queue, Route A/B/C, archive provenance, malformed identity, contradictory evidence, text/JSON, and byte-for-byte immutability cases. An independently built exact-target binary reproduced both findings above; the ordinary Route A implementation control passed. Four saved exact-target heavy lanes already passed without skips and were not rerun. The detached worktree was clean and removed without force; temporary binary and fixture directories were removed. No background task remains.

**Restatement sweep:** Compared the new phase table and classifier with the exact-target work action and reference templates, and searched shipped consumers for `advance`, its typed result, and refusal tokens. Registration, renderer and prime agree with the new public projection. The optional Route A lessons mismatch is recorded above. The phase classifier's omission of fenced-content boundaries was confirmed with the real binary rather than inferred from tests. Current later revisions still contain both underlying classifier paths; intentional mutation extensions are outside this review's attribution.

**Self-validation:** All three P-A-U boxes are checked. The scope expansion and canonical blocked-probe choice are documented decisions. The stored RED failure exercised the missing command, and GREEN exercises registration rather than merely compiling helpers. Extra edge cases were tested because the existing matrix did not exercise optional lessons or Markdown examples. No source edits, queue capture, or commits were made by this reviewer.

**Suggested testing:** 2 items — Route A optional-lessons completion and lifecycle-looking headings inside fenced examples, as specified in the findings.

**Follow-ups created:** None (2 findings report only).

*Reviewed by review-work action; orchestrated artifact for the owning work action to persist.*
