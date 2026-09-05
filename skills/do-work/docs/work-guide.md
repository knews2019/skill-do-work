# Work (Process the Queue)

The central orchestrator. It owns judgment while canonical commands advance queue claims, recover, checkpoint, finalize, archive, release, and commit queue state.

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
1. Advance selection and claim (one typed, committed queue transaction)
2. Triage (route A/B/C)
3. Estimate (P50 active minutes — printed before planning starts)
4. Plan and explore when the route requires them
5. Implement and test
6. Review and prepare finalization
7. Finalize (one resumable lifecycle/archive/release/commit transaction)
8. Continue from `advance` or record the exit checkpoint
```

### The estimate line

Right after triage, before any planning, the run prints something like:

```
Starting REQ-1459 — Add SD-39 review links and QA gates
Estimated active duration: approximately 125 minutes (P50, medium confidence)
Dominant factors: Route C, browser evidence matrix, performance gate, storage auditing
```

**P50 means roughly a 50% chance of completing within the estimated active minutes** — half of similar REQs finish faster, half slower. "Active" counts only time the agents are actually working (planning, building, testing, review, remediation); waiting on your input, paused sessions, and queue wait are excluded. It is an informational forecast, never a deadline: nothing in the pipeline gates on it, and estimation failures never stop a run. Multi-REQ runs also print per-REQ estimates plus two labeled totals — total effort (everything summed) and critical path (the longest dependency chain, which is what wall-clock time actually follows when work runs in parallel).

## What accumulates in the REQ file

As the request moves through the pipeline, sections are appended:

- `## Triage` — route decision and reasoning
- `estimate:` frontmatter block — P50 active-minutes forecast, frozen once execution begins
- `## Plan` — implementation plan (Route C)
- `## Exploration` — key files, patterns, gotchas (Routes B & C)
- `## Scope` — declared files to touch and acceptance criteria (Routes B & C)
- `## Implementation Summary` — manifest of files changed (mandatory)
- `## Testing` — test results, red-green validation
- `## Review` — scores, findings, acceptance result
- `## Decisions` — numbered implementation choices (D-01, D-02...)
- `## Discovered Tasks` — out-of-scope issues found during work; only `impact-critical` items auto-queue, and every noncritical line ends `→ report only` in the archived REQ
- `## Lessons Learned` — what worked, what didn't

## Review gate

After testing, a multi-dimensional review runs:

| Result | Action |
|--------|--------|
| Pass (75%+) | Archive as completed |
| Partial (50-74%) | Archive; auto-queue only `impact-critical` findings and keep the rest in the report |
| Fail (<50%) | One remediation attempt, then archive with issues |

## Open Questions

Builders never block on ambiguities. They mark questions as `- [~]` with best-judgment reasoning and create `pending-answers` follow-ups. Run `do-work clarify` later to review these decisions as a batch.

If you answer a question yourself mid-run — a long run stopping to ask you something between REQs — that answer is written into the REQ as `- [x]` before any building starts. It survives the session that heard it, so no builder re-decides it in a fresh context.

## Checkpoints

Claims write structural entries immediately. At session exit, `advance --checkpoint` is the sole checkpoint writer: it refreshes queue state while preserving every live foreign or unlabelled record. On the next run, canonical `recover` settles finalization first and returns typed claim decisions. A takeover happens only through its exact command or an explicit sole-authority run; writer labels remain data and never become shell source.

## What happens when you run it

A typical `do-work run` session:

1. **Advance** — selects claimable queue work, records holds, and commits each claim
2. **Triage** — reads the REQ, assesses complexity, picks Route A/B/C
3. **Estimate** — ensures a P50 forecast exists and prints the estimate line before planning
4. **Build** — implements the request (planning and exploration for B/C routes)
5. **Test** — runs the project's test suite, validates red-green if TDD targets exist
6. **Review** — scores the work against requirements, code quality, and acceptance criteria
7. **Prepare** — chooses the terminal result and routes critical findings
8. **Finalize** — journals and applies archive, release, commit, and provenance once
9. **Loop** — follows `advance`'s frozen continuation, starts a fresh advance, or writes the checkpoint through `advance --checkpoint`

Each REQ is fully processed before the next one starts. If context limits are hit mid-REQ, a checkpoint is written so the next session can resume.

## What `run` does (and does not) do

A bulk `do-work run` has a few properties worth knowing before firing 20 REQs at once.

- **Dependency-aware ordering (opt-in via frontmatter).** If REQs declare `depends_on: [REQ-IDs]` in their frontmatter, the work loop honors it — a REQ is only picked up once every member of its `depends_on` has reached `completed` or `completed-with-issues`. REQs without `depends_on` fall back to numeric ID order. Cycles in `depends_on` are detected and the affected REQs are held under `status: blocked-dependency-cycle` for the user to resolve. Run `do-work roadmap` before a bulk run to see what's classified as Ready vs Blocked. To force a scoped run that ignores dependency gating for a specific REQ, use `do-work run REQ-NNN`; to scope a run to a whole capture, use `do-work run UR-NNN` — a UR expands to its member REQs and runs them in dependency order (unlike an explicitly named REQ, a UR-expanded member does *not* bypass gating). For wave-by-wave execution one dependency depth at a time, use `do-work run --wave N` (roots are depth 0) — and add `--fan-out [N]` to build that depth's REQs concurrently instead of one at a time. REQs that use `dependencies:` instead of `depends_on:` are honored as a legacy alias so muscle-memory typos don't silently bypass gating; `depends_on:` is canonical and wins when both are present.
- **Claim from anywhere; release from one place.** Any checkout you point at a queue — a worktree, a second workspace, a clone, a cloud session — may capture REQs, claim them, and build them. A successful claim atomically commits its queue move and checkpoint entry, so it is available to other checkouts on their next ordinary git sync. Two unsynced checkouts can still claim the same REQ; nothing locks or arbitrates them, and their committed footprints conflict when the branches meet. What stays single is the **release tail**: one checkout merges, bumps the version, writes the `CHANGELOG.md` entry, moves files into `archive/`, and closes URs (`actions/work-reference.md` → Execution Model — Claim Anywhere, One Releaser). `do-work verify`'s duplicate-id probe catches colliding captures. Within one checkout the loop still finishes a REQ before starting the next **unless you pass `--fan-out`**: `do-work run --fan-out [N]` computes the ready set itself — pending, dependencies met, unclaimed, not earmarked for someone else, and not dropped by a filter you passed — and dispatches that many builders at once, each in its own tree, with no confirmation prompt. You can still pick the set by hand instead, and either way the releaser merges, verifies, and archives them one at a time (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). The saving is in the build phase; everything after the merge stays serial. On a harness without worktree support the flag quietly does nothing and you get the serial loop.
- **`write_set` helps you pick, but never picks for you.** A REQ's `write_set` frontmatter (the repo-relative paths it expects to write) is a display-only hint that feeds the board's *overlaps* badge — useful when choosing which REQs to run together, but nothing schedules on it, and no badge means *unknown*, not safe.
- **No mid-run pause for clarification.** Open Questions are answered by the builder with logged reasoning and a `pending-answers` follow-up REQ is queued for batch review — except a question whose real answerer is a named outside stakeholder, which lands on that person's stakeholder REQ with an HTML report to share (`do-work stakeholder-answers` ingests their reply later). You'll see your own questions when you next run `do-work clarify` — the loop itself never blocks on a prompt, in either direction.
- **Waiting on an external condition uses `status: blocked`.** When a REQ can't start until something outside the queue is true — LM Studio running, a designer answering, credentials provisioned — it carries `status: blocked` and a free-text `blocked_by` naming the condition (set at capture or when the builder hits the missing precondition mid-run before any edits land). This is distinct from `pending-answers` (a question for you) and `depends_on` (a wait on another REQ). Blocked REQs sit out of the run. If the REQ has an optional `blocked_check` shell probe, queue-mode `advance` re-runs it and atomically unblocks a success before claim; otherwise confirm the condition via `do-work clarify` or edit the status back to `pending`. The one blocked shape with its own exit is a stakeholder-questions REQ (`stakeholder:` in frontmatter) — it clears through `do-work stakeholder-answers`, never a probe or a clarify confirm. They surface on the board's *Needs input · Blocked* column with a "blocked by" badge.
- **Failures classify, archive, and queue follow-ups; the loop always continues.** A failed REQ is classified, archived as `failed`, and a follow-up REQ is queued when appropriate; the loop then proceeds to the next pending REQ. Failures that trace back to a failed upstream REQ (via `addendum_to` or `depends_on`) are auto-classified as `spec` with an upstream pointer in the error message — so cascading failures aren't misdiagnosed as fresh code bugs. To triage what landed (including any `pending-answers` follow-ups for completed-with-issues outcomes), run `do-work clarify` after the queue drains.

## Several checkouts against one queue

You can point more than one checkout at the same `do-work/` queue — a second local workspace, a clone, a cloud session, a spawned worktree — and have them cooperate. The rule is short: **anyone may claim and build; exactly one checkout releases.**

**Claiming from anywhere.** Run `do-work run` in whichever checkout you're sitting in. It claims a REQ with one committed transaction containing the queue-to-working move, frontmatter flip, and checkpoint entry. Another checkout sees that claim after ordinary git sync; two unsynced checkouts can still claim and build the same REQ, and the duplicate surfaces when their committed footprints merge. Nothing locks, nothing waits, and no checkout asks another for permission.

**One releaser.** Pick one checkout to run the release tail: merging integration, bumping the version, writing the `CHANGELOG.md` entry, moving REQs into `archive/`, and closing URs. Everything else can happen anywhere; that part should not happen twice. Two changelog prepends against one queue is the failure to watch for, and unique version numbers do not make it safe. There's no mechanism enforcing this — it's a decision you make once and keep.

**Earmarking with `assigned_to`.** To say "leave this one for me", add one field to a pending REQ:

```yaml
assigned_to: 'cloud-alpha'
```

Use the single-quoted **Frontmatter Quoting** form from `../actions/work-reference.md`; double an apostrophe inside a session name. Any other checkout's `do-work run` then skips it and tells you why (`REQ-042 — assigned to cloud-alpha`), and the board shows an `assigned` badge. It's a courtesy, not a lock: nothing checks whether that session exists or is still running. Two ways to take an earmarked REQ anyway — name it explicitly (`do-work run REQ-042`), which claims it and clears the field, or delete the field by hand. Reaching it through `do-work run UR-011` does *not* override the earmark: naming a whole capture is a weaker signal than naming the REQ.

**When two checkouts claim the same REQ.** Nothing stops them, and the fix happens where the branches meet. Expect two things from that merge:

- **The REQ file** conflicts on content — the two `claimed_at` values and whatever sections each side wrote. Keep one claim: whichever checkout actually has the work. (It's never a rename conflict, even though the file moved — both sides moved it to the same place, so git resolves that part silently.)
- **`do-work/CHECKPOINT.md`** conflicts too, and it conflicts on *every* concurrent claim, including two REQs that overlap in nothing. **Keep every entry from both sides.** This one matters: taking only yours deletes another checkout's record of live work, and taking only theirs means your own crash can't be recovered. One REQ listed under two checkouts is not a contradiction to tidy up — it's the honest record of two claims.

If both checkouts made byte-identical claims (same second, nothing written yet), the REQ file won't conflict at all. The checkpoint still will, because each entry is stamped with the checkout that wrote it — which is the only thing that catches that case.

**Captures collide the same way, and the same fix applies.** Two checkouts capturing at once can both mint `REQ-042`. Merge them, renumber one, and run `do-work verify` — its duplicate-id probe is there for exactly this. `do-work verify` also flags two related drifts: a REQ you're building here that's still earmarked for somewhere else, and a UR archived while one of its members is still live.

**Recovery is one public command.** Plain `recover` preserves working claims and returns typed takeover options after finalization is clean. `recover --take-over REQ-NNN` authorizes one claim; `run-with-recovery` uses `recover --assume-sole-authority` to reset working claims and every same-request checkpoint entry without interpreting writer text as code — except a REQ whose finalization tail recovery set aside, which keeps its claim so the run does not select it again.

**Building several REQs at once.** Within one checkout, `do-work run --fan-out [N]` computes the ready set itself — pending, dependencies met, unclaimed, not earmarked elsewhere — and dispatches that many builders concurrently, each in its own git worktree, with no confirmation prompt. `--fan-out 3` sets the count; bare `--fan-out` uses your harness's limit, or two. It composes with `--wave N` (`--wave` picks *which* REQs, `--fan-out` picks *how many at once*), and on a harness without worktree support it quietly does nothing and you get the ordinary serial loop.

The wave doesn't check whether the REQs touch the same files. A computed set means "these are all runnable", not "these won't collide" — collisions surface when the branches merge, and that's deliberate: a REQ's `write_set` is a display hint whose absence means *unknown*, so scheduling on it would produce sets that look safe without being safe. Integration stays serial regardless, so the time you save is in the build phase and nowhere else.

**Skipping work nobody would notice.** A REQ can carry an `impact:` verdict — whether anyone would ever notice the work — and capture writes one on every new REQ. `do-work run --skip-impact-negligible` leaves the `impact-negligible` ones where they are and builds the rest, telling you how many it passed over so a narrowed run never reads as an empty queue. It picks *which* REQs run, so it stacks with `--wave N` and composes with `--fan-out [N]` the same way. It errs toward doing the work: a REQ with no `impact:`, or one whose value is misspelled, reads as user-visible and is never skipped — only a REQ somebody actually judged negligible is passed over. Naming a REQ outright — `do-work run REQ-042` — still builds it: an explicit name beats the filter, exactly as it beats dependency gating and another session's earmark. Capture also writes the token into the REQ's title, so typing `impact-negligible` into the board's search box lists them without the flag.

**What isn't specified.** Two checkouts both running the release tail, and two `do-work` sessions in the same working tree. Nothing prevents either, nothing detects them, and if one happens, `do-work forensics` shows you the damage and `do-work cleanup` helps you fix it. Don't do them.

**One caveat that changes everything above.** All of it depends on `do-work/` being committed to git. That's the arrangement this skill's own repo uses, and it's what lets a claim, a checkpoint entry, or a captured REQ travel between checkouts at all. If your `do-work/` is untracked (the default for most projects), nothing syncs: the poisoning can't happen, and neither can any of the merge-time detection. `do-work verify`'s duplicate-id probe is then the only cross-checkout check you have, run by hand in each one.

## Trigger aliases

All of these do the same thing — process the queue:

```
do-work run
do-work go
do-work start
do-work work
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
- **Failed items** — If a REQ fails review, the system tries one remediation pass. If it still fails, it archives with issues noted and auto-queues only critical findings. Noncritical findings stay in that report; promote one by running `do-work capture` with its complete finding line quoted as the source.
- **Context limits** — Claims are recorded immediately and `advance --checkpoint` refreshes the exit state. Start a new session with `do-work run`; its first canonical recovery result says whether anything needs takeover authority.
- **One at a time** — The work action processes one REQ per loop iteration. This keeps commits atomic and reviews focused. Don't try to batch multiple REQs into one pass.

## Clarify mode

```
do-work clarify
do-work questions
do-work pending
```

Reviews all `pending-answers` REQs. You can confirm the builder's choice, override it, skip, or discard. Answered REQs flip back to `pending` and re-enter the queue.
