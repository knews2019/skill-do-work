## Review: REQ-684

**Approve**: `qualify` now says which Implementation Summary problem it found. The change is the smallest one that fixes the problem.
Route A | range d342db9b..5cc6f7c7 (merge 5cc6f7c7)

### What's built
- `handleQualify` (`skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go:290-296`) reports "Implementation Summary section not found" when the heading is absent and "Implementation Summary lists no backticked file paths" when the heading exists but no `- \`path\`` line is under it. The finding code, severity, fixability and next-action text are unchanged.

### Decisions / risks for you
- None.

### Findings

**Important:** None.

**Minor:** None.

**Nit:**
- `checks.go:291`: the parser (`allBacktickedPaths`, `checks.go:794`) only reads list lines that start with `` - ` ``. A summary that names a backticked path in a prose sentence gets "lists no backticked file paths", even though a path is visible in backticks. The word "lists" covers this case, so no change is needed. impact-negligible → report only

### Checks the caller asked for
- (a) The default string is reachable only when `found` is true, `paths` is empty and `parseError` is nil. The outer condition is `parseError != nil || !found || len(paths) == 0`. The first branch catches `parseError != nil` and the `else if` catches `!found`, so only `found && len(paths) == 0` is left. Met.
- (b) `parseError` wins. It is tested first, and `allBacktickedPaths` returns `found=true` with its only error (unmatched backtick, `checks.go:799`), so the order cannot hide it. Met. No test covers `parseError` precedence inside `handleQualify`, but that branch is unchanged and the REQ does not ask for a test.
- (c) Both tests assert the finding count, the code and the exact evidence string. I ran the merged `checks_test.go` against the base `checks.go` (exported with git archive into the scratchpad). Both tests failed on the old string. On the merged tree they pass (`go test -count=1 -run TestQualify ./internal/corehelpers/`, ok).

### Anti-bloat
- New helpers: 0. New options or flags: 0. New files: 0. New finding codes: 0. New branches: 1 (`else if !found`), which is the split the REQ asked for. New tests: 2, one per message as required, none decorative. The inline fixture setup repeats about 6 lines per test instead of extending `writeQualificationRequest`. Hand-back decision D-02 records this, and extending the helper would have added a parameter. No verification-only route, marker or validator was added (F14a in the upstream report stays rejected).

### Requirements Checklist
- [x] `!found` → "Implementation Summary section not found": delivered
- [x] found and no paths → "Implementation Summary lists no backticked file paths": delivered
- [x] Code `QUALIFY-SUMMARY-MISSING`, severity, fixability, other fields and the `parseError` branch unchanged: delivered (the `helperFinding` line is untouched)
- [x] One test per message with the exact string: delivered
- [x] Constraints (no new code, helper, file or route; message change only): kept

### Acceptance Testing
**Result: Pass** (implementation and integration stages)
- Focused qualify tests pass on the merged tree. RED was confirmed on the base code.
- The REQ records a green repository gate at 5cc6f7c7 and a green probe through `advance`. The qualify gate ran this REQ's own summary through the new code.

### Suggested Additional Testing
- None. No deployment or live stage applies beyond the normal release.

### Restatement sweep (own elements; not a wave end)
Redefined element: the evidence text of `QUALIFY-SUMMARY-MISSING` (old text "Implementation Summary is missing or empty", one condition).
- `grep -rn "missing or empty"` over `skills/`, `_dev/` and `README.md` (CHANGELOG and `do-work/archive` excluded) found only `_dev/tests/contracts/core-checks.sh:833` ("shared principles file missing or empty"), which is a different check. No match for the old qualify text.
- `grep -rn QUALIFY-SUMMARY-MISSING` outside the archive found code only: `checks.go:297`, the two new tests, and `commands.go:398`. That line maps the code to a `sed` inspection command and does not depend on the evidence text. No action, crew-member, lessons satellite or doc describes when the code fires.
- Doc phrasings with "no/missing Implementation Summary" (`skills/do-work/actions/forensics.md:102`, `skills/do-work/docs/forensics-guide.md:14`, `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:43`) belong to other checks (hollow completions, archive records), not qualify.
- The old text is still in historical and third-party records only: the UR-152 input, the upstream report in the inbox, this run's brief and hand-back, and a source snapshot inside `ai-reports/2026-09-06_2113_do-work-cli-guide/index.html` pinned to an old commit. None of these restates the current contract.
- Result: no stale restatement.

### Scores (on the record, not the headline)

## Review

**Overall: 97%** | <STAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None. Nit: `checks.go:291` "lists no backticked file paths" also fires when a backticked path appears only in prose, because the parser reads `` - ` `` list lines only. The wording still fits. impact-negligible → report only
**Acceptance:** Pass. Implementation and integration stages: both new tests RED on the base code and GREEN at 5cc6f7c7, the qualify tests pass, and the REQ records a green gate.
**Anti-bloat:** 0 helpers, 0 options, 0 files, 0 finding codes, 1 branch (the required `else if !found`), 2 tests (one per message, as required). No verification-only route.
**Restatement sweep:** redefined `QUALIFY-SUMMARY-MISSING` evidence text. Grepped `skills/`, `_dev/` and `README.md` for "missing or empty" and for the finding code. No doc, action, crew file, lessons satellite or test restates the old text or its single condition. `commands.go:398` keys on the code only. The old text survives only in historical and third-party records (UR-152 input, inbox report, run brief and hand-back, old ai-report source snapshot).
**Suggested testing:** 0 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*
