## 0.305.22 — Batch Shell Audits while Preserving Diagnostics (2026-09-06)

Batched compatible ShellCheck invocations and Markdown pattern audits across the shell verification suite, eliminating repeated process startup while preserving fragment syntax isolation and exact line attribution.

- `_dev/tests/action-shell-blocks.sh` batches ShellCheck across 74 extracted Markdown fences and 33 shipped shell files, writing `.meta` companion files to map GCC-formatted diagnostics back to exact Markdown paths and line numbers, eliminating 106 ShellCheck Haskell process spawns (-33.5% wall time).
- `_dev/tests/prescribed-shell-canonicalization.sh` scans 165 Markdown files using multi-pattern `grep -F -f` passes, eliminating over 2,600 child grep process launches (-84.1% wall time).
- Verifies defect detection across fragment syntax errors, quiet-grep pipeline violations, ShellCheck warnings, and canonicalization guide restatements.
