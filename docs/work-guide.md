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

Builders never block on ambiguities. They mark questions as `- [~]` with best-judgment reasoning and create `pending-answers` follow-ups. Run `do-work clarify` later to review these decisions as a batch.

## Checkpoints

At session end, a `do-work/CHECKPOINT.md` is written with the last completed REQ, queue state, and where any in-progress work stopped — so the next session can resume cleanly.

## What happens when you run it

A typical `do-work run` session:

1. **Queue scan** — finds the next `pending` REQ file in `do-work/queue/`
2. **Claim** — moves it to `working/` and sets `status: claimed` so no other agent grabs it
3. **Triage** — reads the REQ, assesses complexity, picks Route A/B/C
4. **Build** — implements the request (planning and exploration for B/C routes)
5. **Test** — runs the project's test suite, validates red-green if TDD targets exist
6. **Review** — scores the work against requirements, code quality, and acceptance criteria
7. **Archive** — moves the REQ to `archive/`, creates follow-up REQs if the review flagged issues
8. **Commit** — one atomic commit per REQ with explicit file staging
9. **Loop** — wipes context and picks the next REQ (or exits if the queue is empty)

Each REQ is fully processed before the next one starts. If context limits are hit mid-REQ, a checkpoint is written so the next session can resume.

## What `run` does (and does not) do

A bulk `do-work run` has three properties worth knowing before firing 20 REQs at once.

- **Dependency-aware ordering (opt-in via frontmatter).** If REQs declare `depends_on: [REQ-IDs]` in their frontmatter, the work loop honors it — a REQ is only picked up once every member of its `depends_on` has reached `completed` or `completed-with-issues`. REQs without `depends_on` fall back to numeric ID order. Cycles in `depends_on` are detected and the affected REQs are held under `status: blocked-dependency-cycle` for the user to resolve. Run `do-work roadmap` before a bulk run to see what's classified as Ready vs Blocked. To force a scoped run that ignores dependency gating for a specific REQ, use `do-work run REQ-NNN`. For wave-by-wave execution one dependency depth at a time, use `do-work run --wave N` (roots are depth 0).
- **No mid-run pause for clarification.** Open Questions are answered by the builder with logged reasoning and a `pending-answers` follow-up REQ is queued for batch review. You'll see the questions when you next run `do-work clarify` — the loop itself never blocks on a prompt.
- **Halt on failure is opt-in.** By default, a failed REQ is classified, archived as `failed` with a follow-up REQ created when appropriate, and the loop continues to the next pending REQ. Failures that trace back to a failed upstream REQ (via `addendum_to` or `depends_on`) are auto-classified as `spec` with an upstream pointer in the error message — so cascading failures aren't misdiagnosed as fresh code bugs. Add `--halt-on-failure` to stop the loop after the first failed or completed-with-issues REQ; useful for high-stakes foundation work where you want a checkpoint to review before continuing.

## Trigger aliases

All of these do the same thing — process the queue:

```
do-work run
do-work go
do-work start
do-work begin
do-work process
do-work execute
do-work build
do-work continue
do-work resume
```

Use whichever feels natural. `continue` and `resume` read well after a break; `run` and `go` are good for fresh starts.

## Tips

- **`continue` vs fresh `run`** — No functional difference. Both scan the queue and pick the next pending REQ. Use `continue` when you're resuming a session; use `run` when you're starting fresh. The checkpoint system handles the actual resume logic.
- **Failed items** — If a REQ fails review, the system tries one remediation pass. If it still fails, it archives with issues noted and optionally creates a follow-up REQ. You don't need to intervene manually.
- **Context limits** — Long-running queues may hit context limits. The system writes `do-work/CHECKPOINT.md` before stopping. Just run `do-work run` again in a new session — it picks up where it left off.
- **One at a time** — The work action processes one REQ per loop iteration. This keeps commits atomic and reviews focused. Don't try to batch multiple REQs into one pass.

## Clarify mode

```
do-work clarify
do-work questions
do-work pending
```

Reviews all `pending-answers` REQs. You can confirm the builder's choice, override it, skip, or discard. Answered REQs flip back to `pending` and re-enter the queue.
