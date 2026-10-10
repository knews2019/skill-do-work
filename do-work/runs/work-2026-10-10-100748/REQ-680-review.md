## Review: REQ-680 (archive fetch test keeps an absent target absent)

**Approve.** The new `absent target` row pins the named failure with the same failing request as the old case, and the old case keeps every check it had.
Route A | merge 2414535d on 1dbe058e (builder commit 9b8dee8b)

### What's built
- `TestTotalFailurePreservesTheTargetAndLeavesNoScratch` is now a two-row table: `pre-existing target` and `absent target`.
- The absent row seeds no file, sends the same failing fetch (HTTP server returns "not an archive", Git route points at a missing repository), checks the same four failure-report fragments, then checks that the target path does not exist and that no `.fetching.` scratch file is left.
- Production code is unchanged.

### Decisions / risks for you
- None. The builder's D-03 (the REQ's "only the new case fails" mutation sentence cannot hold, because removing private staging also breaks the pre-existing row and other staging tests) is a correct reading of the code, recorded in the REQ `## Testing` section.

### Findings

**Important:**
- None.

**Minor:**
- None.

**Nit:**
- N1. `archive_fetch_test.go:423`: when the target exists, the message prints `a failed fetch created the absent target: <nil>`, because a nil error is the value formatted. Saying "target exists" would read more clearly. — impact-negligible → report only
- N2. `archive_fetch_test.go:425` and `:442`: both branches call `assertNoArchiveScratch`. One call before the branch would remove the early `return` and the duplicate. — impact-negligible → report only

### Requirements Checklist

- [x] Absent-target case added to the existing failure test as a table row (REQ requirement 1, builder latitude on row versus sibling test) — delivered
- [x] Same failing request as the pre-existing row — delivered (the server, request and fragment checks are shared code inside the loop body)
- [x] "A failed fetch never creates a target that did not exist" — delivered: `os.Lstat(targetPath)` must report not-exist. Any other result, including an existing file or symlink, fails the row
- [x] "No scratch left" — delivered: `assertNoArchiveScratch` scans for `.fetching.`, which matches both staging names in `archive_fetch.go:205` and `:486`
- [x] Pre-existing row keeps full coverage — delivered: failure fragments, byte equality, mode 0644 and the scratch scan are all still asserted
- [x] No consumer commit IDs or REQ numbers in added comments (REQ requirement 2) — delivered: grep of added lines for `REQ-`, `UR-` and 7+ hex characters prints nothing. The comment says in plain words what each starting state protects
- [x] Test-only, no production change — delivered: the diff touches one `_test.go` file
- [x] P-A-U boxes checked — all three `[x]`

### Anti-bloat count
New helpers 0, new options 0, new files 0, new constants 0, new test functions 0, decorative tests 0. One anonymous struct field (`existingBytes`) belongs to the table and replaces a repeated literal. Nothing to name as a finding.

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `go test -count=1 -run TestTotalFailurePreservesTheTargetAndLeavesNoScratch -v ./internal/archivefetch/` at HEAD: both rows PASS, package `ok`.
- `gofmt -l internal/archivefetch/` prints nothing.
- Mutation proof not re-run by the reviewer (read-only review). The builder's recorded result (new row fails on its own assertion when private staging is removed) matches what the assertion checks.
- The integrator's full gate at 2414535d passed (recorded in REQ `## Testing`).

### Suggested Additional Testing
- None. No deployment or live-acceptance stage applies to a test-only change.

### Restatement sweep
Nothing redefined. A test-table change does not change a contract that other text restates.

### Scores (on the record, not the headline)

**Overall: 100%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Both requirements and all constraints met |
| Code Quality | 100% | Two nits only, zero weight |
| Test Adequacy | 100% | Row pins the named failure; mutation evidence recorded |
| Scope | 100% | Exactly the one declared file |
| Risk | None | Test-only |
| Acceptance | Pass | Targeted test green at the merge |

### Follow-ups created
- None (2 findings report only)

## Review

**Overall: 100%** | TIMESTAMP-PLACEHOLDER

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None. Nit N1: `archive_fetch_test.go:423` prints `<nil>` when the absent target was created; "target exists" would be clearer — impact-negligible → report only. Nit N2: `archive_fetch_test.go:425` and `:442` call `assertNoArchiveScratch` in both branches; one call before the branch would do — impact-negligible → report only
**Acceptance:** Pass — implementation and integration stages: targeted `go test -run TestTotalFailurePreservesTheTargetAndLeavesNoScratch` green at 2414535d (both rows), gofmt clean
**Restatement sweep:** nothing redefined
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*
