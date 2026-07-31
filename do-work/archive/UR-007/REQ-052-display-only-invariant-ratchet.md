---
id: REQ-052
title: Contract-regression ratchet pinning the overlap badge's display-only invariant
status: completed
route: A
claimed_at: 2026-07-29T13:55:51Z
commit: 7861204
domain: testing
tdd: false
maintenance: false
prime_files: []
created_at: 2026-07-29T09:32:07Z
user_request: UR-007
addendum_to: REQ-034
depends_on: []
write_set: ["_dev/tests/contract-regressions.sh"]
---

# Contract-regression ratchet pinning the overlap badge's display-only invariant

## What

Add a ratchet to `_dev/tests/contract-regressions.sh` pinning the display-only invariant of the board's write-set overlap annotation: `annotateWriteSetOverlap` is called *after* `bucketColumns` in `tools/queue-kanban/`, and `actions/board.md` retains its display-only wording.

## Why

The invariant (overlap annotation never influences column placement) is currently protected by one Go test (`TestWriteSetOverlapNeverAffectsColumnPlacement`) plus prose. A ratchet also guards the instruction-side wording against a future edit quietly turning the display annotation into column logic or deleting the doc claim. Approved by the user via `do-work clarify` on 2026-07-29 (follow-up from REQ-034, surfaced in REQ-041).

## Constraints

- Anchor on the call-site line and the heading — REQ-033's review showed weak string anchors can be gutted while the check still passes. Make the anchors specific enough that moving the call before bucketing, or rewording the doc claim away, fails the suite with a message naming the file and the fix.

## Acceptance

- `_dev/tests/contract-regressions.sh` fails if `annotateWriteSetOverlap` no longer runs after bucketing, or if `actions/board.md` drops the display-only claim.
- The suite passes on the current tree.

## Triage

**Route:** A
**Reasoning:** Single file (`_dev/tests/contract-regressions.sh`); add a ratchet with specific anchors (the `annotateWriteSetOverlap`-after-`bucketColumns` call-site + board.md display-only heading). Route A; anchor specificity is the risk.
**Rigor:** Standard main-context review (part of the parallel disjoint-write_set batch 051/052/054/057/058; single-file, no spec-cluster overlap).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the REQ, `crew-members/general.md`, `crew-members/coding-guardrails.md`, `crew-members/testing.md`, the existing ratchets in `_dev/tests/contract-regressions.sh` (idiom: `assert_contains` / `assert_block_contains` / bespoke `if` blocks for anything grep can't express), the real call sites via `grep -rn 'annotateWriteSetOverlap\|bucketColumns' tools/queue-kanban/`, and `actions/board.md`'s Rules section. Approach: the ordering condition is not expressible as a string match, so extract the `buildBoard` block (`sed -n '/^func buildBoard(/,/^}/p'`) and compare the *line numbers* of the two call sites — `bucketColumns(board.AllRequests` must precede `annotateWriteSetOverlap(board.AllRequests)` — with a distinct failure branch for "one call site vanished" so a rename can't make the check silently vacuous. Anchor deliberately on the calls, **not** the `// Deliberately AFTER bucketing` comment: the comment survives the move (proven in RED 1). The doc condition gets two `assert_block_contains` scoped to the `## Rules` block (not file-wide — the REQ-044 masking failure mode where a match anywhere in the file satisfies the assertion), pinning the after-bucketing claim with its `model.go` pointer and the `never column logic` display-only claim.
- [x] **[APPLY]:** Implemented exactly as planned. Scope strictly limited to `_dev/tests/contract-regressions.sh` (+33 lines at the end of the suite, before the `fail_count` exit). No other file touched — red-green was run against a `tar`-cloned sandbox of the tree in the session scratchpad precisely so the temporary mutations to `tools/queue-kanban/model.go` and `actions/board.md` never landed in the live tree that four concurrent builders were writing to; sandbox deleted afterward.
- [x] **[UNIFY]:** Red-green per condition, all five in the sandbox clone (baseline green there first). **(a) Ordering — RED 1:** hoisted `annotateWriteSetOverlap(board.AllRequests)` above the `bucketColumns` call inside `buildBoard`, *leaving the "Deliberately AFTER bucketing" comment untouched* → `FAIL: tools/queue-kanban/model.go calls annotateWriteSetOverlap BEFORE bucketColumns in buildBoard …`, exit 1; restored → passed. **RED 4:** renamed the call site to `annotateOverlapAndPlace(board)` → `FAIL: … must call both bucketColumns(board.AllRequests ...) and annotateWriteSetOverlap(board.AllRequests) — one call site was renamed or removed …`, exit 1 (proves the anchor fails loud instead of going vacuous); restored → passed. **(b) Doc claim — RED 2:** deleted `runs *after* column bucketing` from the board.md Rules sentence → `FAIL: actions/board.md Rules must keep the claim that annotateWriteSetOverlap runs after column bucketing …`, exit 1. **RED 3:** deleted `— never column logic, never blocking` → `FAIL: actions/board.md Rules must keep the overlap annotation display-only claim …`, exit 1. **RED 5:** moved the whole parser lock-step bullet out of `## Rules` (text still present in the file) → both board.md assertions fired, exit 1 — confirms the heading scoping, not just a file-wide grep. Every restore returned the suite to green. Live tree: `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` (exit 0); `shellcheck _dev/tests/contract-regressions.sh` clean; `bash -n` clean. `git diff --stat -- _dev/tests/contract-regressions.sh` → `1 file changed, 33 insertions(+)`; the other three modified files in `git status --porcelain` (`actions/capture-reference.md`, `actions/forensics.md`, `tools/queue-kanban/generate_test.go`) are the concurrent builders' work, not this REQ's. No git, frontmatter, or commit operations performed.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — added ratchet block pinning the board write-set overlap annotation's display-only invariant (call-order assertion + two Rules-block doc assertions)

**What was done:**

Added one ratchet block at the end of the suite pinning the board write-set overlap annotation's display-only invariant on both sides — code and instruction.

1. **Call-site ordering assertion (code side).** Extracts the `buildBoard` body from `tools/queue-kanban/model.go` and resolves the block-relative line numbers of `bucketColumns(board.AllRequests` and `annotateWriteSetOverlap(board.AllRequests)` with `grep -nF … | head -1 | cut -d: -f1 || true`. Two failure branches: either call site missing → FAIL naming the file and telling the maintainer to restore the calls (annotation last) or update the anchor in the same commit; annotation line before the bucketing line → FAIL naming the file, the fix (move it back below `bucketColumns`), and where the co-dispatch decision actually belongs (`actions/work.md` Step 1's gate). Anchored on the calls rather than the adjacent comment, which a hoist leaves intact. Anchored on the call argument (`bucketColumns(board.AllRequests`) rather than the assignment LHS so a variable rename doesn't trip it.
2. **Two `assert_block_contains` on `actions/board.md`'s `## Rules` block (instruction side).** One pins ``​`annotateWriteSetOverlap` in `tools/queue-kanban/model.go` runs *after* column bucketing`` — the claim plus its file pointer; the other pins `never column logic`. Both are scoped to the `## Rules` → `## Common Rationalizations` block rather than the whole file, so relocating the bullet out of Rules fails too (the REQ-044 masking pattern, where a file-wide match satisfies the assertion while the load-bearing restatement rots).
3. **Comment header** explaining why the invariant is one movable line, why the assertion is on the call order and not a comment, and that the doc claim is the second way the invariant quietly stops being a promise.

Idiom matches the surrounding suite: the two-word-minimum variable names (`build_board_block`, `bucket_columns_call_line`, `overlap_annotation_call_line`, `board_rules_block`), `sed`-extracted blocks fed to `assert_block_contains`, and a bespoke `if`/`elif` for the one condition grep can't express — the same shape as the router word-budget and hardened-check-script blocks.
