---
id: REQ-115
title: status/testing_status normalization suppressed its own contract warning
status: completed
created_at: 2026-08-05T19:39:11Z
claimed_at: 2026-08-05T19:39:11Z
completed_at: 2026-08-05T19:41:30Z
route: A
user_request: UR-023
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: [REQ-112]
addendum_to: REQ-112
write_set: [tools/queue-kanban/frontmatter_cli.go, tools/queue-kanban/model.go, tools/queue-kanban/frontmatter_cli_test.go]
maintenance: false
review_generated: true
---

# status/testing_status Normalization Suppressed Its Own Contract Warning

## What

`queue-kanban frontmatter get <file> status --normalize` on a typo like `completedd` printed the typo to stdout, exited 0, and emitted no Schema Read Contract warning. Same for `testing_status`. The command shipped in REQ-112 to replace hand-rolled status reads, and this left exactly the no-feedback path it exists to close.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Root cause is a forced `recognized = true` for the two fields that are not alias-map rows in the contract table. Give them canonical-vocabulary rows (no alias maps — those stay in `normalizeStatus`/`normalizeTestingStatus`, which `normalizeSchemaField` dispatches to), then let the CLI's normal `isKnownSchemaFieldValue` check run for every field. Verify: RED on 9 new cases, then GREEN, then exercise the binary.
- [x] **[APPLY]:** Two table rows in `model.go`, the special-case block removed from `frontmatter_cli.go`, and an accurate warning for fields with no documented default.
- [x] **[UNIFY]:** `gofmt -l` clean, `go vet` clean, `go test ./...` passes. Exercised the binary on four fixtures: typo'd status warns and prints the truth; aliased status stays silent; typo'd status fails `--in-set` with a warning; `domain` still substitutes its real default. `_dev/tests/contract-regressions.sh` at its 7 pre-existing failures, unchanged.

## Why

Found by Codex's review of PR #130 (P2). Verified by reproduction before accepting.

## Implementation Summary

**What was done:** Restored the contract warning for `status` and `testing_status`, and made the warning text truthful for fields the contract gives no default.

**Files changed:**
- `tools/queue-kanban/model.go` (modified) — added canonical-vocabulary rows for `status` and `testing_status` (no alias maps, by design); `schemaFieldWarningText` now says "No default is defined; reporting it unchanged" instead of the false `Treating as ''`
- `tools/queue-kanban/frontmatter_cli.go` (modified) — removed the forced `recognized = true`; every contract-row field now goes through `isKnownSchemaFieldValue`
- `tools/queue-kanban/frontmatter_cli_test.go` (modified) — 2 new tests, 9 cases

## Decisions

- **D-01: vocabulary rows without alias maps.** DECIDE & STATE. Copying `normalizeStatus`'s aliases into the table would create the second definition the table exists to prevent. The rows carry only `canonicalValues`, so `isKnownSchemaFieldValue` and `schemaFieldWarningText` can answer while the aliases stay in one place.
- **D-02: an unrecognized value with no documented default prints what was found.** DECIDE & STATE. The contract explicitly declines to define a `status` default ("never claim or archive an unrecognized status silently"), so substituting one would be inventing data. The caller sees `completedd` plus a warning — the truth and the feedback.

## Qualification

Passed — 3 files in the diff, root cause traced to a specific branch, fix verified end-to-end rather than only in unit tests.

## Testing

**Red-green validation:** RED — `TestRunFrontmatterCommandWarnsOnUnrecognizedStatus` failed with `warning present = false, want true`, and the `--in-set` variant failed too. GREEN — both pass; 9 cases cover aliases staying silent, typos warning, and hand-edited values warning.

**Regression check:** aliases must not start warning. `status: done` → prints `completed`, stderr empty. Verified against the binary.

## Review

**Approve** — closes a functional hole in freshly shipped code, with the alias path proven still silent.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor. **Follow-ups created:** None.

## Lessons Learned

**What worked:** Reproducing the finding against the built binary before touching code. It confirmed both the symptom and the exit code, so the fix targeted the real branch rather than the one I assumed.

**What didn't:** I had already flagged this exact code in REQ-112's own review as a Minor finding — "it puts a second piece of dispatch knowledge next to `normalizeSchemaField`'s own" — and filed it as a code-organization smell. It was a functional hole. Noticing that a branch is structurally awkward and *not* asking what it does behaviourally is how a self-review misses a bug it has already looked straight at.

**Worth knowing:** A special case added to make a type check pass (`recognized = true` so the two dispatch-only fields wouldn't fall through) silently became a policy decision about warnings. When a branch exists to satisfy a mechanism, state what it means for behaviour, or it will mean something nobody chose.

## Orientation

`queue-kanban frontmatter --normalize` now warns for all nine contract fields instead of seven, so a typo'd `status` is reported rather than passed through. Lives in `tools/queue-kanban/`. Leaf fix; no contract or interface change.
