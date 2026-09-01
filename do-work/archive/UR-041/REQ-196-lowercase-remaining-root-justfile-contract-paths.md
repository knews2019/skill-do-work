---
id: REQ-196
title: Remaining late contract assertions use capitalized root Justfile
status: completed
created_at: 2026-08-15T12:11:13Z
claimed_at: 2026-08-15T14:30:05Z
completed_at: 2026-08-15T14:40:23Z
commit: 5f15929
user_request: UR-041
addendum_to: REQ-180
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-180, REQ-187]
write_set: [_dev/tests/contract-regressions.sh]
route: A
kb_status: promoted
kb_entry: REQ-196-remaining-late-contract-assertions-use-c.md
---

# Remaining Late Contract Assertions Use Capitalized Root Justfile

## What

Replace the four remaining late `assert_contains` inputs that open `Justfile` with the tracked lowercase `justfile`, and prove those assertions execute on a case-sensitive filesystem.

## Why

REQ-180 repaired the two earlier casing mismatches that aborted the aggregate, but four later assertions still open a nonexistent capital-J root path. Default macOS filesystems resolve the wrong spelling, so the aggregate can pass locally while those assertions fail on Linux or another case-sensitive filesystem.

## Context

- Discovered during REQ-187 while tracing the root Just delegate and intentionally left outside that canonical-gate Scope.
- The remaining inputs are the four `"Justfile"` arguments near the final managed-marker, board-source, and updater-fallback assertions in `_dev/tests/contract-regressions.sh`.
- This is an addendum to REQ-180's exact tracked-casing repair.

## Detailed Requirements

- Change all four remaining root-file inputs from `Justfile` to `justfile` without changing their assertion meanings.
- Keep intentional prose and fixture references to the general `Justfile` filename variant unchanged.
- Add or reuse case-sensitive evidence that the four late assertions read the tracked root file and remain reachable.
- Preserve the canonical maintainer gate and aggregate final pass marker.

## Constraints

- Do not add case-insensitive path resolution or filesystem-specific production logic.
- Do not rewrite unrelated `Justfile` prose or multi-variant fixtures.
- Lock-in limit: zero live contract inputs open a root filename whose case differs from the tracked `justfile`.

## Red-Green Proof

**RED prompt/case:** Run the late root-justfile assertions against a case-sensitive checkout; all four capital-J inputs name a nonexistent file.
**Why RED now:** The current checkout tracks only lowercase `justfile`, while the four inputs survive because the development filesystem is case-insensitive.
**GREEN when:** All live root inputs use exact tracked casing, the focused case-sensitive replay reaches the four assertions, and the full aggregate passes through its final marker.

## Open Questions

- [x] Auto-approved: test-only mechanical hygiene (normal). → Added to queue.

---
*Source: discovered during REQ-187 canonical maintainer gate implementation.*

---

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Confirmed the request names the exact four late assertion inputs, one implementation file, and a case-sensitive proof requirement. Kept Route A direct implementation and used the tracked lowercase `justfile` spelling as the sole repair.
- [x] **[APPLY]:** Added the non-vacuous source ratchet first, captured all four capital-J inputs in RED, then changed only those four live assertion paths to lowercase while preserving assertion patterns, prose, and filename-variant fixtures.
- [x] **[UNIFY]:** Reviewed the sole source diff and exact four-input inventory. Focused mutations, the full aggregate final marker, canonical maintainer gate, `bash -n`, warning-level ShellCheck, and `git diff --check` pass with no unrelated or debug changes.

## Triage

**Route: A** - Simple

**Reasoning:** The defect, all four literal edits, the sole source file, and the required case-sensitive replay are already identified. This is a focused test-harness casing repair with no architecture or source-discovery uncertainty.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — changes the four remaining live root assertion inputs from `Justfile` to tracked `justfile` and adds a case-sensitive, deletion-resistant source ratchet for their exact patterns

**Behavior:** The late managed-marker, board-source, and updater-fallback assertions now open the tracked root file on every filesystem. Their meanings and execution order are unchanged, and the aggregate still reaches its final pass marker.

## Testing

**RED:** With the ratchet added before the literal repair, `bash _dev/tests/contract-regressions.sh` exited 1 after all earlier child suites. It reported each of the four expected patterns using `['Justfile']`, reported the wrong-case live-root inventory, and did not print the final aggregate pass marker.

**GREEN:**
- focused source replay — PASS with exactly four lowercase live inputs
- wrong-case, assertion-deletion, assertion-duplication, and source-shape-drift mutations — correctly rejected
- `bash _dev/tests/contract-regressions.sh` — PASS through `Contract regression checks passed.`
- `bash _dev/tests/maintainer-verify.sh` — PASS, including ShellCheck, aggregate contracts, both Go modules, and the strict JavaScript lane
- `bash -n _dev/tests/contract-regressions.sh` — PASS
- direct warning-level ShellCheck — PASS
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — only `_dev/tests/contract-regressions.sh` changed; foreign REQ-189–192 edits and orchestration/release state are excluded.
- **Mechanical checks:** PASS — P-A-U is complete, the modified file is present in the diff, and no debug artifacts were added.
- **Substance and traceability:** PASS — the four captured path defects map one-for-one to the four literal changes, while the exact-pattern ratchet supplies the requested filesystem-independent casing/reachability evidence.
- **Wiring/data flow:** PASS — the repaired paths feed the existing top-level `assert_contains` calls before the shared `fail_count` gate and final pass marker.

## Review

**Overall: 99%** | 2026-08-15T14:40:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 99% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — all four live inputs use exact tracked casing, the ratchet rejects wrong-case/deleted/duplicated/source-shape mutations, and the aggregate remains reachable through its final marker.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Parsing the exact live assertion source shape provides case-sensitive evidence even when the host filesystem aliases `Justfile` and `justfile`. Pinning each expected pattern separately also proves all four late assertions remain present and reachable.

**What didn't:** Fixing only the first two occurrences in REQ-180 left four later inputs hidden by macOS case-insensitive lookup. A local green aggregate was therefore not sufficient evidence that every live path used tracked casing.

**Worth knowing:** Keep intentional filename variants in prose and fixture loops; the enforceable boundary is the path argument consumed by a live root-file assertion, not every textual occurrence of “Justfile.”

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Orientation

The canonical maintainer aggregate now uses the tracked lowercase root `justfile` in every live contract input, so its late checks behave consistently on case-sensitive and case-insensitive filesystems.
