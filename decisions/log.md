# Decision Log — Timeline

Append-only timeline of every decision entered into the ADR log. Newest entries go at the bottom (chronological order, oldest first), matching the BKB `wiki/log.md` pattern.

Each line: `YYYY-MM-DD | ADR-NNNN | action | title`

Actions: `new` (first-time entry), `amended` (frontmatter or body change), `superseded` (marked inactive because a newer ADR replaced it).

---

## 2026

- 2026-04-15 | ADR-0001 | new | Capture ≠ Execute Boundary
- 2026-04-15 | ADR-0002 | new | UR + REQ Pairing on Every Capture
- 2026-04-15 | ADR-0003 | new | Immutable Working and Archive Folders
- 2026-04-15 | ADR-0004 | new | Platform-Agnostic Action Files
- 2026-04-15 | ADR-0005 | new | Subagent Dispatch with Direct-Read Fallback
- 2026-04-15 | ADR-0006 | new | Priority-Ordered Routing Table
- 2026-04-15 | ADR-0007 | new | Crew Members — JIT-Loaded Domain Rules
- 2026-04-15 | ADR-0008 | new | Queue Canonical Path — `do-work/queue/`
- 2026-04-15 | ADR-0009 | new | Companion Reference File Pattern
- 2026-04-15 | ADR-0010 | new | REQs as Validated Intent

_This log was seeded in one pass on 2026-04-15 from a review of CHANGELOG.md v0.1.0 through v0.64.1. Dates in frontmatter reflect when the original decision was made; the `new` entries here reflect when each was codified as an ADR._
