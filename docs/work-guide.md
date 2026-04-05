# Work (Process the Queue)

The central orchestrator. Picks up pending requests and works through them one by one — triage, plan, build, test, review, archive.

## Complexity triage

Each request is assessed and routed:

| Route | When | Pipeline |
|-------|------|----------|
| **A (Simple)** | Bug fixes, config changes, copy updates | Build → Test → Review |
| **B (Medium)** | Clear goal, unknown location | Explore → Build → Test → Review |
| **C (Complex)** | New features, architectural changes | Plan → Explore → Build → Test → Review |

When uncertain, the system defaults to Route B (under-planning is recoverable; over-planning wastes time).

## Pipeline steps

```
1. Find next pending REQ
2. Claim it (move to working/, status: claimed)
3. Triage (route A/B/C)
4. Plan (Route C only — architecture, file list, testing approach)
5. Explore (Routes B & C — find relevant files and patterns)
6. Implement (all routes — build the thing)
7. Test (run tests, validate red-green if TDD)
8. Review (requirements check, code quality, acceptance testing)
9. Archive (move to archive/, create follow-ups if needed)
10. Commit (one commit per REQ, explicit file staging)
11. Loop or exit (context wipe, pick next REQ)
```

## What accumulates in the REQ file

As the request moves through the pipeline, sections are appended:

- `## Triage` — route decision and reasoning
- `## Plan` — implementation plan (Route C)
- `## Exploration` — key files, patterns, gotchas (Routes B & C)
- `## Scope` — declared files to touch and acceptance criteria (Routes B & C)
- `## Implementation Summary` — manifest of files changed (mandatory)
- `## Testing` — test results, red-green validation
- `## Review` — scores, findings, acceptance result
- `## Decisions` — numbered implementation choices (D-01, D-02...)
- `## Discovered Tasks` — out-of-scope issues found during work
- `## Lessons Learned` — what worked, what didn't

## Review gate

After testing, a multi-dimensional review runs:

| Result | Action |
|--------|--------|
| Pass (75%+) | Archive as completed |
| Partial (50-74%) | Archive, create follow-up REQs for important findings |
| Fail (<50%) | One remediation attempt, then archive with issues |

## Open Questions

Builders never block on ambiguities. They mark questions as `- [~]` with best-judgment reasoning and create `pending-answers` follow-ups. Run `do work clarify` later to review these decisions as a batch.

## Checkpoints

At session end, a `do-work/CHECKPOINT.md` is written with the last completed REQ, queue state, and where any in-progress work stopped — so the next session can resume cleanly.

## Usage

```
do work run
do work go
do work start
do work continue
```

## Clarify mode

```
do work clarify
do work questions
do work pending
```

Reviews all `pending-answers` REQs. You can confirm the builder's choice, override it, skip, or discard. Answered REQs flip back to `pending` and re-enter the queue.
