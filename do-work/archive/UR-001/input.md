---
id: UR-001
title: Add folder-cohesion and cyclomatic complexity dimensions to code-review
created_at: 2026-05-25T11:06:37Z
requests: [REQ-001]
word_count: 67
---

# Add folder-cohesion and cyclomatic complexity dimensions to code-review

## Full Verbatim Input

do-work capture request: Add a folder-cohesion / orphan-files check to the code-review action (Step 4 Architecture) — detect files whose imports, naming, or concerns don't match the folder's apparent domain, and folders that have accumulated unrelated files. Also promote cyclomatic complexity from a tie-breaker (quick-wins Step 5) to an explicit named dimension in code-review's Architecture section, distinct from cyclic dependencies.

## Conversation Context

This request emerged from a comparison of `do-work quick-wins` vs `do-work code-review` against a user-described audit checklist (folders with too many files, files/functions with too many lines, cyclomatic dependencies, orphan files in modular folders).

Findings of the comparison:
- `quick-wins` covers file/function size (god files >300 lines, long functions >50 lines) but not folder-level cohesion or cyclic/cyclomatic concerns.
- `code-review` already covers cyclic dependencies (Step 3, "Circular dependencies?") and folder *structure* (Step 3, "Inconsistent file/folder structure"), but does NOT explicitly check for orphan files within folders, and does NOT name cyclomatic complexity (McCabe number) as a first-class signal. Cyclomatic complexity is only mentioned once in `quick-wins` Step 5 as a tie-breaker for risk-impact scoring.
- The user-described "cyclomatic dependencies" is ambiguous — it could mean cyclic deps (already covered) or McCabe complexity (not first-class anywhere). Making cyclomatic complexity an explicit named dimension resolves the ambiguity.
- File/function size thresholds were deliberately NOT added to code-review — that remains quick-wins's job; the two actions stay complementary.

---
*Captured: 2026-05-25T11:06:37Z*
