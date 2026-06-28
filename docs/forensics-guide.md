# Forensics

Pipeline diagnostics — detects stuck work, hollow completions, orphaned URs, scope contamination, and other health issues. Read-only and safe to run anytime.

> **Sister action:** `do-work roadmap` is the read-only survey for *intended* state — what's queued, in-progress, and feasible to pick up next. Forensics looks for *broken* state. If nothing is broken but you want to know "where are we and what's next," see `docs/roadmap-guide.md`.

## What it checks

| Check | What it detects |
|-------|----------------|
| **Stuck work** | REQs in `working/` claimed >1hr (warning) or >24hr (critical) |
| **Hollow completions** | Completed REQs with no Implementation Summary or file changes |
| **Missing qualifications** | REQs lacking `## Qualification` section (post-v0.38.0) |
| **Orphaned URs** | UR folders in `user-requests/` where all REQs are archived but UR wasn't moved |
| **Scope contamination** | Files modified by 3+ unrelated REQs, or overlapping files within same UR |
| **Failed without follow-up** | Failed REQs missing error classification or follow-up REQ |
| **Stale pending-answers** | REQs waiting for user input for >7 days |
| **Git divergence** | Files from completed REQs later modified or deleted without tracking |
| **Stranded finished REQs** | Terminal-status REQs left in `do-work/queue/` or `working/` instead of archived |
| **Recurring corrections** | The same lesson/correction theme surfacing across 2+ archived REQs (2 = watch, 3+ = strong signal) — a sign to fix the harness, not the next run |

## Output

Markdown report organized by severity:

- **Critical Findings** — needs immediate attention
- **Warnings** — should be addressed soon
- **Info** — awareness items

Each finding includes a suggested fix. Sections with no findings are omitted.

## Key rules

- Read-only — reports findings, never auto-fixes
- User decides what to act on

## Usage

```
do-work forensics
do-work diagnose
do-work health check
do-work health
```
