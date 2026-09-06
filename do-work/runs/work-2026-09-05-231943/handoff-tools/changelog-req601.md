## 0.305.13 — Correct Stale Mechanism and Script Claims Across Shipped Callers (2026-09-06)

Documentation and action files across toolbox, board, and core now accurately describe underlying Go command implementations, staging paths, error conditions, and conflict resolution semantics.

- `ai-report-reference.md` documents that single-image and batch image generation stage in the system temporary directory rather than adjacent to the target, and removes nonexistent retained script fallback references.
- `install.md` documents temporary skill downloads landing in `SKILL.md.download.<random>`, specifies that local Git ignores are managed inline following the canonical contract rather than calling an external helper, and removes nonexistent compatibility script references.
- `present-work.md` clarifies that portfolio summary publication reads the source once and writes both outputs from that buffer, and removes nonexistent compatibility script references.
- `board.md` clarifies that adding a local git exclude on non-git projects exits zero with a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning finding.
- `prescribed-shell-primitives.md` accurately describes inventory conflict resolution: a parseable `completed_at` beats an unparseable or missing timestamp regardless of walk order, and ties fall back to `working/` before `archive/`.
- `architecture-report.md` removes nonexistent compatibility script references.
- `work-reference.md` clarifies that `run-timed-command` connects both child streams directly to the CLI's stderr handle.
