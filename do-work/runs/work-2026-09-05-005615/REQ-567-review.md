# Review: REQ-567 — Repair shipped lesson links

**Approve** — the three shipped lesson links now point to their existing UR-095 archives and Lessons Learned sections. The focused reference contract passes. One pre-existing, local knowledge-base pointer issue remains report-only.

Route C | Implementation target `2dda18a1f816e148258295b2b351f12695a049b4`

## What's built

The implementation changes only three URL destinations in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, lines 25–27. It preserves their labels, family markers, and lesson prose. The REQ's Decisions section explicitly records the choice of canonical UR-095 paths with `#lessons-learned` anchors.

## Requirements Checklist

- [x] All three shipped URLs identify the canonical UR-095 archive files and their Lessons Learned anchors. Independently verified in the exact target Git objects.
- [x] No obsolete root-level archive URL for REQ-491, REQ-492, or REQ-493 remains under `skills/`.
- [x] The recorded canonical repository-gate failure is repaired without changing the gate or unrelated files. The request records the same gate passing after integration; exact-target heavy verification subsequently passed all four selected lanes without skips.
- [x] The repair preserves the parent UR-098 intent by unblocking its lifecycle-command work. This repair performs no lifecycle migration itself, so the parent's four-part migration write set does not apply to this dependency repair.
- [x] All three P-A-U checkboxes are checked; the declared one-file Scope matches the complete implementation diff.

## Acceptance Testing

**Result: Pass.**

- Reviewed `git diff ffc6cf03c6ab77c8ead85a8925c193436b62c686..2dda18a1f816e148258295b2b351f12695a049b4`: one file, three insertions and three deletions, all URL replacements.
- Independently resolved the three old URLs against base Git objects: all three root-level archive targets are absent. Resolved the three replacement URLs against target Git objects: all three UR-095 files exist and contain an exact `## Lessons Learned` heading; every URL ends in `#lessons-learned`. This directly verifies deletion of the captured finding surface and its intended replacement.
- Mechanically masked only the three URL values and compared the full before/after file bodies: identical. Compared the three current bullets with the exact target: identical.
- Ran `bash _dev/tests/shipped-package-reference-contract.sh` on the current checkout: exit 0, `shipped package reference contract: PASS`, 0.78 seconds. This is a fresh current-checkout check, distinct from the saved exact-revision evidence.
- Reused recorded exact-target heavy evidence at `2dda18a1f816e148258295b2b351f12695a049b4`: `do-work-cli-integrations` exit 0 / 50s; `staged-skills` exit 0 / 22s; `updater` exit 0 / 53s; `installer` exit 0 / 24s. No skips. The heavy suite was not rerun during review.

The historical Testing section records green-gate evidence at `e738544ff7e6de5349aee910664e774759f0495c`; that older entry is not represented as fresh exact-target execution. The separate Heavy Verification Result supplies the exact-target evidence above.

## Restatement Sweep

No API or lifecycle semantics changed. Swept the corrected archive-path references across shipped files and other repository text. Shipped references agree with the canonical destinations. Local knowledge-base source pages still contain obsolete root-level pointers, recorded below; corresponding `kb/raw/processed/` documents retain historical provenance and were not altered. Queue/history references were read as historical records, not instructions.

## Review

**Overall: 100%** | 2026-09-04T21:57:19Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

Overall is the arithmetic mean of the four percentage dimensions; no qualitative modifier applies. The report-only finding below predates this narrowly scoped shipped-reference repair and does not indicate an implementation defect or scope drift.

**Important findings:**

- F1 — Existing local knowledge-base source pointers still send readers to absent root-level archive paths: `kb/wiki/sources/REQ-491-add-canonical-repository-gate-deferral-l.md:69`, `kb/wiki/sources/REQ-492-integrate-repository-gate-deferral-and-r.md:45`, and `kb/wiki/sources/REQ-493-review-fix-complete-repository-gate-defe.md:34`; their actual REQs live under `do-work/archive/UR-095/`. These pages are outside the shipped-link repair's declared scope. — impact-user-visible → report only

**Minor findings:** None.

**Acceptance:** Pass — all three corrected destinations and anchors exist at the reviewed target, the named stale surface is absent, and the focused reference contract exits 0.

**Suggested testing:** 0 items. The existing reference contract and exact-target heavy results cover this change; remote publication availability was not tested or claimed.

**Follow-ups created:** None (1 findings report only).

**Self-validation:** Completed. Distinguished current-checkout execution from saved exact-target evidence, checked actual archive headings independently of the reference contract, checked unchanged prose and scope, and recorded the stale local pointers without mutating their records. No source, queue, request, archive, or Git state was changed by this review.

*Reviewed by review-work action*
