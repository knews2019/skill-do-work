# Roadmap

Read-only queue survey — what's done, what's in progress, what's pending, and a feasibility read on what's actionable next. Never modifies REQs, frontmatter, or files.

> **Sister action:** `do-work forensics` is the read-only diagnostic for *broken* state (stuck, hollow, orphaned). Roadmap looks at *intended* state. If you suspect something is wrong rather than "where are we," see `docs/forensics-guide.md`.

## What it surveys

| Section | What it shows |
|---------|---------------|
| **Notes** | Lightweight next-step hints from `do-work note` (`do-work/notes.md`), verbatim in append order — rendered only when notes exist |
| **Ready to Pick Up** | Queue REQs with clear scope, no `pending-answers`, no unresolved blockers |
| **Needs Clarification** | `status: pending-answers` or open questions in the body |
| **Blocked** | Depends on a REQ still pending/in-progress, or `status: blocked` waiting on an external condition (named in `blocked_by`) |
| **Stale** | Created >30 days ago and never claimed — re-confirm before working |
| **TDD Eligible** | `tdd: false` but the behavior is testable (Red-Green Proof, I/O example, backend domain) |
| **In Progress** | REQs in `working/` — id, route, claimed-for, current phase, TDD posture |
| **Recently Completed** | Grouped by UR or week; flags URs whose REQs are all done but the UR isn't archived |
| **Lessons Promoted / Pending** | REQs whose Lessons Learned are staged in `<kb>/raw/inbox/` (promoted) vs. captured but not staged (pending) |

## Output

Markdown report. When `do-work/notes.md` is non-empty, a **Notes** block renders first — even when the queue is empty — then **Ready to Pick Up** so the actionable section is the next thing the reader sees. Empty sections are omitted — if the queue is empty, the report says so explicitly. Caps each section at 20 entries by default.

## Key rules

- Read-only — never modifies REQs, moves files, or creates commits
- Feasibility is a read, not a verdict — flags concerns, never reclassifies a REQ as blocked
- Cites evidence for every classification (frontmatter field, section, missing artifact)
- No ideation — surveys what exists; for new ideas use `do-work scan-ideas`

## Usage

```
do-work roadmap
do-work roadmap pending
do-work roadmap UR-014
do-work roadmap since 2026-04-01
do-work queue-status
do-work where are we
do-work what's left
do-work what should I work on next
```

## When NOT to use

- Suspect something is *broken* or *stuck* → `do-work forensics`
- Want *new ideas* for what to build → `do-work scan-ideas`
- Want to *review specific completed code* → `do-work review-work` or `do-work code-review`
- Want to *explain uncommitted local changes* → `do-work inspect`
