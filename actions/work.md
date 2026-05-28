# Work Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to process the queue. Processes pending requests from the `do-work/queue/` folder in your project. User-facing walkthrough: [`docs/work-guide.md`](../docs/work-guide.md).

An orchestrated build system that processes request files created by the capture requests action. Uses complexity triage to route simple requests straight to implementation and complex ones through planning and exploration first.

## When to Use

**Use when:**
- The queue has `pending` REQs and the user wants them built (`do-work run`, `start`, `go`, etc.).
- The pipeline dispatches the work action as its build step.
- A specific REQ id was named (`do-work run REQ-042`) — the action scopes to it.

**Do NOT use when:**
- The queue is empty — tell the user and stop; suggest `do-work capture request: [describe]` instead.
- The only REQs left are `pending-answers` — route to `do-work clarify` so the user can resolve them first.
- See `SKILL.md` routing table for sibling action selection (inspect, verify requests, review work, etc.).

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary (mandatory file manifest), Testing, Review. This ensures full traceability — what was planned vs done, what files were touched, and whether triage was accurate.

This living log is also the **trail of intent**. The REQ starts as a validated statement of what the user wants (written by capture). As the work action processes it, each appended section documents how intent was interpreted and realized: builder decisions (## Decisions) record where the builder exercised judgment beyond stated intent, scope declarations (## Scope) record what the builder committed to, and implementation summaries record what was actually built. The gap between captured intent and realized implementation is visible in a single file.

## Architecture

```
work action (orchestrator - lightweight, stays in loop)
  │
  ├── Read CHECKPOINT.md if exists (resume context from previous session)
  │
  ├── For each pending request (skip pending-answers):
  │     │
  │     ├── TRIAGE: Assess complexity (no agent, just read & categorize)
  │     │
  │     ├── OPEN QUESTIONS? ── - [ ] items exist ──► Mark - [~], builder decides
  │     │                      (none / all resolved) ──► continue
  │     │     │
  │     │     ├── Route A (Simple) ──────────────────┐
  │     │     │   Skip plan/explore, direct to build │
  │     │     │                                      │
  │     │     ├── Route B (Medium) ───────┐          │
  │     │     │   Explore, scope declare  │          │
  │     │     │                           ▼          │
  │     │     └── Route C (Complex) ──► Plan ──► Explore ──► Scope declare
  │     │                                            │
  │     │                                            ▼
  │     │                                     Implementation agent
  │     │                                            │
  │     │                                            ▼
  │     │                                  Implementation Summary
  │     │                                            │
  │     │                                            ▼
  │     │                              Qualify (orchestrator verifies)
  │     │                                            │
  │     │                                            ▼
  │     │                                        Testing
  │     │                                            │
  │     │                                            ▼
  │     │                                  Review ◄─── Fail? ──► Remediate ──► Re-review
  │     │                                            │
  │     │                                            ▼
  │     ├── Archive ──► classify discovered tasks ──► queue follow-ups
  │     │                                            │
  │     │                                            ▼
  │     └── Commit (git repos only)
  │
  └── Context wipe → Loop | Write CHECKPOINT.md → cleanup → report
```

> **Remember:** Every completed request gets a git commit (Step 9) before looping to the next request.

**Sub-agent note:** This document uses "spawn agent" language. Use your platform's subagent mechanism when available. If your tool doesn't support subagents, run phases sequentially in the same session and label outputs clearly.

## Complexity Triage

Before spawning any agents, assess the request to determine the right route.

### Route A: Direct to Builder (Simple)

Skip planning and exploration entirely.

**Indicators:** Bug fix with clear steps, value/config change, single UI element add/remove, styling tweak, request names specific files, well-specified with obvious scope (<~50 words), copy changes, feature flag toggle.

**Examples:** "Change button color from blue to green", "Fix crash when clicking submit with empty form", "Update API timeout from 30s to 60s"

### Route B: Explore then Build (Medium)

Skip planning, run exploration. The "what" is clear, the "where" or existing patterns need discovery.

**Indicators:** Clear outcome but unspecified location, "like the existing X", need to find similar implementations, modifying something at an unknown location.

**Examples:** "Add form validation like we have on the login page", "Create a new API endpoint following our existing patterns"

### Route C: Full Pipeline (Complex)

Plan, then explore, then implement.

**Indicators:** New feature requiring multiple components, architectural changes, ambiguous scope ("improve", "refactor"), touches multiple systems, external service integration, long request (100+ words) with many requirements, user explicitly asked for a plan.

**Examples:** "Add user authentication with OAuth", "Implement dark mode across the app", "Refactor state management to use Zustand"

### Decision Flow

```
Read the request
  ├── Names specific files AND has clear changes? → Route A
  ├── Bug fix with clear reproduction? → Route A
  ├── Simple value/config/copy change? → Route A
  ├── Clear outcome but location/pattern unknown? → Route B
  ├── Ambiguous, multi-system, or architectural? → Route C
  └── Default: Route B (builder can request planning if needed)
```

**When uncertain, prefer Route B.** Under-planning is recoverable; over-planning is wasted time.

## Folder Structure

```
do-work/
├── queue/                         # Pending REQ files (the work queue)
│   └── REQ-018-pending-task.md
├── user-requests/                 # UR folders (verbatim input + assets)
│   └── UR-003/
│       ├── input.md
│       └── assets/
├── working/                       # Currently being processed
│   └── REQ-020-in-progress.md
└── archive/                       # Completed work
    ├── UR-001/                    # Archived as self-contained unit
    │   ├── input.md
    │   ├── REQ-013-feature.md
    │   └── assets/
    ├── REQ-010-legacy-task.md     # Legacy REQs (no UR) archive directly
    └── legacy/                    # Consolidated legacy items
```

- **`queue/`**: The queue — only pending `REQ-*.md` files
- **`working/`**: Claimed requests. Immutable to all actions except the work pipeline.
- **`archive/`**: Completed UR folders (self-contained) and legacy REQs/CONTEXT docs
- **`user-requests/`**: Active UR folders. Moved to `archive/` when all REQs complete.

## Request File Schema

Request files use YAML frontmatter added progressively:

```yaml
---
# Set by capture action
id: REQ-001
title: Short descriptive title
status: pending
domain: frontend  # choose one: frontend, backend, ui-design, or general
tdd: false       # optional — set true when test-first applies (per capture's TDD heuristic); drives Step 6 testing-crew loading and RED/GREEN mode
caveman: false   # optional — `true` or intensity `lite` | `full` | `ultra`; loads crew-members/caveman.md to compress agent prose
prime_files: []  # list paths to relevant prime-*.md files, or leave empty
created_at: 2025-01-26T10:00:00Z
user_request: UR-001          # May be absent on legacy REQs
addendum_to: REQ-NNN          # optional — present only when this REQ amends an in-flight or completed REQ; set by capture, or by review when creating follow-ups. **Legacy alias:** every read site (Step 8 upstream walk, Step 8 cycle detection, Step 8 follow-up generation, and roadmap Blocked classification) also recognizes `amends:`, `parent:`, and `amendment_to:` as synonyms when `addendum_to` is absent so natural-English glosses don't silently drop the parent linkage; `addendum_to` wins when multiple are present. Capture and follow-up REQs always emit `addendum_to:` — never propagate the alias.
depends_on: []                # optional list of REQ IDs that must reach `completed` or `completed-with-issues` before this REQ runs. Semantically distinct from `addendum_to` ("amends that REQ"): depends_on is "requires that REQ to be done first." A REQ can have both. Honored by Step 1's selection scan and by Step 8's upstream-failure classification. **Legacy alias:** every read site (Step 1 selection, Step 1 cycle detection, Step 1 `--wave` depth, Step 8 upstream walk, roadmap classification) also recognizes a `dependencies:` key as a synonym so muscle-memory typos from Python/Node/Cargo conventions don't silently bypass gating; `depends_on` wins when both are present. Capture and follow-up REQs always emit `depends_on:` — never propagate the alias.

# Set by work action when claimed
claimed_at: 2025-01-26T10:30:00Z
route: A | B | C

# Set by work action when finished
completed_at: 2025-01-26T10:45:00Z
status: completed | completed-with-issues | failed
commit: abc1234               # If git repo
error: "Description"          # Only if failed

# Set by kb-lessons handoff (Step 7.5 pipeline / Step 9.5 standalone). Optional; absent on REQs that predate the handoff.
kb_status: promoted | pending | declined | skipped
kb_entry: REQ-042-lesson-slug.md   # filename only (survives bkb moves from inbox/ to capture/ to processed/); present only when kb_status: promoted
---
```

### Schema Read Contract

Seven fields above are enum-or-boolean-valued, and an audit of `0.76.2`'s `dependencies:` → `depends_on` patch surfaced that several silently swallow natural typo variants from sister conventions (snake_case vs kebab-case YAML, `done`/`finished`/`closed` as English glosses of `completed`, lowercase route letters, etc.). Pure silent-alias is risky for enum values because an unknown value should not be quietly remapped — it should leave a footprint. Every read site in this file (and in `actions/roadmap.md`) honors a uniform **normalize-and-warn contract** for these fields:

1. **Normalize first.** Apply the per-field alias map below. If a canonical match results, use it silently.
2. **Warn-on-fallback.** If after normalization the value still doesn't match the canonical enum, emit:

   ```
   ⚠ {field}: '{value}' not recognized — expected one of [{enum}]. Treating as '{default}'.
   ```

   and proceed with the documented default.
3. **Never silently drop.** The warning is the missing feedback channel that allowed `dependencies:` to go unnoticed pre-0.76.2. Warnings render in the queue-status summary block (Step 1) or, for fields read outside Step 1, alongside the operation that triggered the read.

| Field (read sites) | Canonical enum | Normalization | Default on unknown |
|---|---|---|---|
| `domain` (Step 4 Route C plan-agent spawn, Step 6 crew load, Step 9 review-work spawn) | `frontend`, `backend`, `ui-design`, `general` | `back-end`/`back_end` → `backend`; `front-end`/`front_end` → `frontend`; `ui_design` → `ui-design` | `general` |
| `status` (Step 1 scan + categorization, Step 10 archive trigger) | `pending`, `claimed`, `completed`, `completed-with-issues`, `failed`, `pending-answers`, `blocked-archive-collision`, `blocked-dependency-cycle` | `done`/`finished`/`closed` → `completed` | skip REQ at Step 1 with the warning text — never claim or archive an unrecognized status silently |
| `route` (Step 3 dispatch, Step 5 scope-drift comparison) | `A`, `B`, `C` | lowercase `a`/`b`/`c` → uppercase | treat as needing re-triage in Step 3 |
| `caveman` (Step 6 crew load) | `false`, `true`, `lite`, `full`, `ultra` | truthy strings (`yes`/`on`) → `true`; `light` → `lite` | `false` |
| `tdd` (Step 6 testing-crew load, Step 6.5 TDD-evidence gate; emission validated in `actions/capture.md`) | `true`, `false` (YAML boolean) | `test_first`/`yes`/`on`/`t` → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 testing crew not loaded; Step 6.5 gate not enforced) |
| `error_type` (Step 8 failure classification, Step 8 upstream-failure short-circuit, forensics) | `intent`, `spec`, `code`, `environment` | (no common typo aliases identified) | `code` |
| `kb_status` (Step 7.5 / Step 9.5 kb-lessons handoff, roadmap lessons rollup) | `promoted`, `pending`, `declined`, `skipped` | `skip` → `skipped`; `rejected` → `declined` | `pending` |

**Write paths are unaffected.** Step 2 claim, Step 8 archive, Step 8 follow-up generation, the kb-lessons handoff, and capture emission always write the canonical key and canonical enum value — never an alias, never the typo'd input. The normalize-and-warn contract is read-only.

**Status flow (frontmatter values):** `pending` → `claimed` → `completed` / `completed-with-issues` / `failed`

The intermediate phases (planning, exploring, implementing, testing, reviewing) are tracked by which `##` sections exist in the REQ file, not by frontmatter status changes. Only three status transitions are written to frontmatter: `pending` → `claimed` (Step 2), then `claimed` → final status (Step 8).

**Special statuses — these REQs stay in the queue but Step 1 won't pick them up (they're not `pending`, so the "find next pending REQ" scan walks right past them):**
- `pending-answers` — a follow-up REQ whose Open Questions need user input before it can be worked. These accumulate in the queue and get batch-reviewed when the user runs `do-work clarify`.
- `blocked-archive-collision` — set by Step 2.0 when a queue file's REQ id is already archived. Non-destructive holding state; the user flips it back to `pending` (or removes/renames the duplicate) after deciding what to do.
- `blocked-dependency-cycle` — set by Step 1 when a REQ's `depends_on` graph contains a cycle (e.g., REQ-A depends on REQ-B which depends on REQ-A). Non-destructive holding state; the user edits the `depends_on` chain to break the cycle, then flips the status back to `pending`.

## Input

`$ARGUMENTS` may contain:

- **Specific REQ IDs** (e.g., `REQ-042`, `REQ-042 REQ-043`) — process only those REQs and stop (do not process the full queue). This is how the pipeline action scopes work to a specific batch. Targeted mode bypasses `depends_on` gating — the user explicitly named the REQs.
- **`--wave N`** (integer flag, default mode only) — run only REQs at dependency depth N. Roots (no `depends_on`, or all `depends_on` resolve to archived REQs) are depth 0; depth grows by one per dependency layer. Mutually exclusive with targeted REQ IDs — reject the combination with an error.

When no REQ IDs and no flags are provided, process all pending REQs in dependency-aware order (default behavior).

## Steps

**The work action is an orchestrator.** You handle ALL file management (moving files, updating frontmatter, appending sections, archiving). Spawned agents handle implementation work only.

### Step 1: Find Next Request

**Crash Recovery:** Before checking the queue, look inside `do-work/working/` for any `REQ-*.md` files. If any exist, a previous run was interrupted. For each recovered REQ:
1. Reset frontmatter: set `status` to `pending`, **unless** the REQ file contains a `## Open Questions` section with at least one unresolved `- [ ]` item — in that case, restore `status` to `pending-answers`. (If the `## Open Questions` section exists but all items are already `[x]` or `[~]`, or if no `## Open Questions` section exists at all, set `status` to `pending`.) Remove `claimed_at` and `route`.
2. Strip sections generated during the interrupted run: remove `## Triage`, `## Exploration`, `## Plan`, `## Scope`, `## Pre-Flight`, `## Implementation Summary`, `## Qualification`, `## Testing`, `## Review`, `## Lessons Learned`, `## Decisions`, and `## Discovered Tasks` sections (and their content) if present — these may be incomplete or stale from the crash. Leave `## Open Questions` and user-authored content intact.
3. Move the REQ back to `do-work/queue/`

Once `working/` is empty, proceed with finding the next request.

Glob for `do-work/queue/REQ-*.md`. Sort by number. Read the frontmatter of each (in number order) to check `status`. Don't read the full body at this stage.

**Dependency-aware selection.** For each `pending` REQ, evaluate its `depends_on` field (or its legacy alias `dependencies:` — recognized for back-compat; `depends_on` wins when both present). A REQ is **dependency-ready** when every ID in the resolved dependency list reaches a REQ with `status: completed` or `status: completed-with-issues`. Resolve each dependency ID by globbing `do-work/archive/**/REQ-NNN-*.md`, `do-work/archive/**/REQ-NNN.md`, `do-work/queue/REQ-NNN-*.md`, and `do-work/working/REQ-NNN-*.md`. Cache resolution within a single Step 1 invocation — a 20-REQ queue with 3 deps each is 60 globs; cache hits keep the cost flat. A REQ with unmet dependencies is **dependency-blocked** and is skipped by the scan; it surfaces in the composed exit summary if no other pending REQ is dependency-ready. Process dependency-ready REQs in numeric ID order.

**REQs with neither `depends_on` nor `dependencies:` are roots** and are always dependency-ready. Existing REQs (captured before the field existed) behave exactly as before.

**Cycle detection for `depends_on`.** Before evaluating a REQ's dependencies, walk its `depends_on` graph (or `dependencies:` if `depends_on` is absent — same alias rule) collecting visited IDs into a seen set. If you encounter the current REQ's ID during the walk, the graph contains a cycle — set the REQ's `status` to `blocked-dependency-cycle`, report it, and skip. Mirrors the `addendum_to` cycle-detection approach already used in Step 8 substep 5. Non-destructive — the user breaks the cycle by editing the dependency list and flips status back to `pending`.

**Wave execution (`--wave N`).** If the `--wave N` flag is set, compute each pending REQ's dependency depth before the dependency-ready filter:

- Depth 0: REQs with no dependency list (neither `depends_on` nor the legacy `dependencies:` alias), or whose dependency members are all already archived (completed/completed-with-issues).
- Depth K (K > 0): `max(depth of each dependency member in the current pending set) + 1`.
- A dependency member that is neither archived (completed/completed-with-issues) nor in the current pending set — i.e. it sits in `pending-answers`, `blocked-archive-collision`, `blocked-dependency-cycle`, `claimed`, or `failed` — contributes depth 0 to this computation. Depth is only about ordering waves; the member's own gating is handled separately by the dependency-ready filter below, which holds the dependent REQ until every member reaches `completed`/`completed-with-issues`.

Filter the pending list to REQs whose depth equals N, then apply the dependency-ready filter normally. If no REQ at depth N is dependency-ready (or none exists at that depth), render the composed exit summary with a leading `No REQs at wave N (depth-N set is empty or fully gated).` line and exit. `--wave` and targeted REQ IDs are mutually exclusive — reject the combination at parse time with a clear error.

**Targeted mode bypasses dependency gating.** When `$ARGUMENTS` contains explicit REQ IDs, process them in the given order regardless of `depends_on` (or its `dependencies:` alias). The user named them explicitly.

**Queue status summary:** After reading all REQ frontmatter, categorize every REQ by status and print a summary before proceeding:

```
Queue: N pending | N completed/done (awaiting archive) | N pending-answers | N blocked-archive-collision
```

Count `completed`, `completed-with-issues`, and `done` statuses together as "completed/done (awaiting archive)." Count `blocked-archive-collision` separately so held duplicates don't disappear into the silence between "no pending" and "no REQs at all." If any completed/done REQs exist in `do-work/queue/`, add:

```
⚠ N completed REQs across M URs awaiting archive. Run `do-work cleanup` after this session.
```

**Targeted mode:** If `$ARGUMENTS` contains specific REQ IDs, find only those REQs in `do-work/queue/`. Verify each exists and has `status: pending`. If a targeted REQ is missing or not pending, report the issue and skip it. Process only the targeted REQs, then stop after the last one completes (skip the loop-or-exit logic in Step 10).

**Default mode (no REQ IDs in arguments):** Scan for the first REQ with `status: pending` (skip `pending-answers` — those wait for user input).

**Exit paths when no `pending` REQs found:**

The exit report is **composed**, not picked from disjoint branches. Whenever no `pending` REQs are found, lead with `No pending REQs in queue.` and then append every section that has at least one REQ. Three sections may apply, in this order:

1. **Completed/done section** — applies if any REQ in `do-work/queue/` has status `completed`, `completed-with-issues`, or `done`. Read the `user_request` frontmatter field from each to group by UR. Render:

   ```
   ⚠ N completed REQs awaiting archive (UR-137: 3 REQs, UR-138: 1 REQ, ...):
     REQ-351 — [title] (done)
     REQ-352 — [title] (completed)
     ...

   Run `do-work cleanup` to archive completed work, then `do-work recap` to see full history.
   ```

2. **Pending-answers section** — applies if any REQ has status `pending-answers`. Render from frontmatter only — do not open the REQ body to count `## Open Questions` items at this stage (Step 1 reads frontmatter per the queue scan). The count is deferred to `do-work clarify`, which is the action that reads Open Questions sections:

   ```
   ⚠ N REQs awaiting clarification:
     REQ-NNN — [title] (pending-answers)
     ...

   Run `do-work clarify` to batch-review the open questions; resolved REQs flip to `pending` and re-enter the queue.
   ```

3. **Blocked-archive-collision section** — applies if any REQ has status `blocked-archive-collision`. Read the matching archive path from each blocked REQ's frontmatter if recorded; otherwise re-run the Step 2.0 glob (`do-work/archive/**/REQ-NNN-*.md` and `do-work/archive/**/REQ-NNN.md`) to find it. Render:

   ```
   ⚠ N REQs held by archive-collision guard:
     REQ-NNN — [title] (queue file: do-work/queue/REQ-NNN-slug.md)
       already archived at <archive-path>
       recover: rename the queue file (if this is an intentional re-do) or delete it (if it's a stale duplicate), then flip status back to `pending`
     ...
   ```

4. **Blocked-by-dependencies section** — applies if any `pending` REQ has an unmet `depends_on` reference (dependency-blocked) or any REQ has status `blocked-dependency-cycle`. Pending REQs stay `pending` (the gating is dynamic — they become ready as upstream REQs complete); only cycle-detected REQs are flipped to a held status. Render both groups under one heading:

   ```
   ⚠ N REQs blocked by unmet dependencies:
     REQ-NNN — [title] (pending; depends on REQ-MMM, status: <pending|claimed|pending-answers>)
     REQ-PPP — [title] (blocked-dependency-cycle; chain: REQ-PPP → REQ-QQQ → REQ-PPP)
     ...

   Resolve the blocking REQs first, then re-run. To force a scoped run that ignores dependency gating for a specific REQ, use `do-work run REQ-NNN`. To break a dependency cycle, edit the REQ's `depends_on` and flip its status back to `pending`.
   ```

**After rendering all applicable sections, exit the work loop** — do not proceed to Step 2.0 or beyond. There is no `pending` REQ to claim. Step 1's contract on the no-pending path is "render the composed summary, then stop"; the only path that continues is the one where Step 1 finds at least one dependency-ready `pending` REQ.

If **no section applies** (no REQs at all in `do-work/queue/`), report completion and exit. Never silently exit when any of the four sections applies — every non-pending or non-ready REQ in the queue is something the user needs to see.

**Composition is deliberate.** A queue with both `pending-answers` and `blocked-archive-collision` REQs (and no completed/done) renders both sections back-to-back. A queue with all four categories renders all four. The user sees the full picture in one report instead of a single branch's slice.

**REQ validation:** When reading each REQ's frontmatter, verify it has the required fields (`id`, `status`, `title`). If a REQ file has missing or unparseable frontmatter, skip it and report: `⚠ Skipping [filename]: missing required frontmatter ([field]).` Do not let a single malformed REQ block the entire work loop — skip it and continue to the next.

**Exact glob pattern:** `do-work/queue/REQ-*.md` — if this returns no results, do NOT conclude the queue is empty. Verify by listing `do-work/queue/` contents to rule out a bad pattern.

### Step 2.0: Pre-Claim Archive Collision Check

Before claiming the queue file, verify it isn't a duplicate of an already-archived REQ. The most common footgun is rerunning `do-work work` against a queue file whose twin was already archived in a prior run, which silently re-processes and re-commits the duplicate.

Extract `REQ-NNN` from the queue filename or frontmatter. Glob `do-work/archive/**/REQ-NNN-*.md` AND `do-work/archive/**/REQ-NNN.md` (both forms — the second catches REQs archived without a slug suffix). If any match exists, **bail without moving or claiming**:

1. Update the duplicate queue file's frontmatter to `status: blocked-archive-collision`. This is non-destructive — the user can flip it back to `pending` after deciding what to do. It also prevents the next Step 10 → Step 1 re-glob from picking up the same REQ and bailing again (livelock).
2. Report:

   ```
   REQ-NNN already archived at <archiveCollisionMatchPath>; remove the duplicate from do-work/queue/ or rename if this is a re-do. Status set to `blocked-archive-collision` to skip future loop iterations.
   ```

3. Skip the rest of the work loop for this REQ and continue to the next pending REQ. Do not delete the queue file — the user decides whether it's a stale duplicate or an intentional re-do that needs renaming.

**Scope (minimal):** This guard only checks the archive. It does NOT add post-move verification or pre-commit collision guards — those are parallel-orchestrator concerns out of scope for single-orchestrator usage.

If no archive match is found, proceed to Step 2.

### Step 2: Claim the Request

1. `mkdir -p do-work/working` and move the REQ file there
2. Update frontmatter: `status: claimed`, `claimed_at: <timestamp>`

### Step 3: Triage

Read the request, apply the decision flow, update frontmatter with `route`. If a `## Triage` section does not already exist, append to the request file:

```markdown
---

## Triage

**Route: [A/B/C]** - [Simple/Medium/Complex]

**Reasoning:** [1-2 sentences]

**Planning:** [Required/Not required]
```

Report the triage decision briefly to the user.

**Addendum REQs:** If the REQ has `addendum_to` in frontmatter, read the original REQ before building. If the original includes a `## Prior Implementation` section, use it. If it doesn't (e.g., the original was in-flight when the addendum was captured but has since completed), find the original in `do-work/archive/` and read it to understand what was already built — key files, patterns, and approach. This prevents duplicating or conflicting with existing work.

### Step 3.5: Open Questions — Best Judgment, Not a Gate

After triage, scan the REQ for a `## Open Questions` section with `- [ ]` items. Open Questions are **not a blocker** — the builder proceeds with its best judgment and completes the REQ.

Open Questions use checkbox syntax:
- `- [ ]` — **Unresolved**: has `Recommended:` and `Also:` choices from capture
- `- [x]` — **Resolved**: user answered (answer follows `→`)
- `- [~]` — **Deferred**: builder used its best judgment (reasoning follows `→`)

**If unresolved `- [ ]` items exist:**

1. Note them. Read the `Recommended:` default and `Also:` alternatives for each.
2. Mark each as `- [~]` with a numbered decision and the builder's reasoning: `- [~] [question] → **D-01**: Builder chose: [choice]. Reasoning: [why]`
3. Number decisions sequentially per REQ (D-01, D-02, D-03...). Open Questions decisions and Implementation Decisions (Step 6) share the same D-XX ID space — if Open Questions uses D-01 through D-03, the first implementation decision is D-04. After resolving all `- [ ]` items, append a counter comment immediately after the `## Open Questions` section so Step 6 knows the next available ID: `<!-- D-XX counter: last used D-03. Next decision: D-04. -->` If no decisions were made in this step, write `<!-- D-XX counter: none used. Next decision: D-01. -->` These IDs can be referenced by future REQs.
4. Proceed with implementation using those decisions.

The follow-up REQs for builder-decided questions are created during **Step 8 (Archive)** — not here. Step 3.5 just records the decisions; the archive step handles the paperwork after the REQ is fully complete.

**Why not block?** Human time is the bottleneck. The optimal windows for user interaction are: (1) capture time, when the user is actively fleshing out requests, and (2) batch-review time, when the user returns to answer accumulated questions. Blocking mid-build wastes builder capacity on idle waiting.

**`pending-answers` REQs:** These accumulate in the queue. When the user returns, they run `do-work clarify` to review all `pending-answers` REQs at once, answer the questions, and flip the status to `pending` so the next work run picks them up. The work loop skips `pending-answers` REQs — it only processes `pending` ones.

If all `- [ ]` items are already `[x]` or `[~]`, or no Open Questions section exists, skip this step entirely.

### Step 3.7: Spec Loading (optional)

After triage, check if a specification template matches this REQ's domain or task type.

1. **Match by task type:** If the REQ's title or What section clearly indicates a task type (API endpoint, UI component, refactor, bug fix), check `specs/` for a matching template (`specs/api-endpoint.md`, `specs/ui-component.md`, `specs/refactor.md`, `specs/bug-fix.md`).
2. **Match by suggested spec:** If the REQ's frontmatter contains a `suggested_spec` field (set during capture), check `specs/` for that template.
3. **If a matching spec exists**, read it and use it to inform:
   - The implementation checklist order (pass to the planning or implementation agent)
   - Quality standards to verify against (pass to the review step)
   - Common pitfalls to watch for (include in the builder's context)
4. **The spec is guidance, not override** — the REQ's specific requirements always take priority. If the REQ's requirements conflict with a spec's recommendations, follow the REQ.
5. **If no matching spec exists**, proceed normally. Specs are optional — their absence never blocks work.

### Step 4: Planning (Route C only)

**Route C:** Spawn a **Plan agent** with the request content, project context, the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain is missing, falls back to `general` for an unknown value, or the file doesn't exist, skip loading it), and any files listed in the `prime_files` array. Instruct it to use the prime files as the strict index for discovering the source of truth. Do not load global architecture. Ask it to produce a specific implementation plan (files to modify, order of changes, architectural decisions, testing approach). If a `## Plan` section does not already exist, append the output:

```markdown
## Plan

[Plan agent output]

*Generated by Plan agent*
```

**Plan validation (Route C only):** After the Plan agent returns, run a quick quality check before proceeding:

1. **Requirement coverage:** Re-read the REQ's What/Detailed Requirements. Every requirement should map to at least one planned task. Flag uncovered requirements.
2. **No orphan tasks:** Every planned task should trace back to at least one requirement. Tasks that don't address any requirement suggest scope creep.
3. **Scope sanity:** Count the planned tasks. If 5+, flag: "Plan has [N] tasks — quality degrades past 3. Consider splitting this REQ into multiple smaller REQs."
4. **File conflicts:** If the plan mentions modifying files that are currently claimed by another REQ in `do-work/working/`, flag the conflict.

Append validation findings to the `## Plan` section (if any issues found). These are **warnings, not blockers** — the builder can adapt. But flag them visibly so the orchestrator and review step are aware.

**Routes A and B:** Append a skip note (if not already present):

```markdown
## Plan

**Planning not required** - [Route A: Direct implementation / Route B: Exploration-guided implementation]

*Skipped by work action*
```

### Step 5: Exploration (Routes B and C)

Spawn an **Explore agent** to find relevant files, existing patterns, types/interfaces, and testing conventions.

- **Route C**: Give it the plan and ask it to find files mentioned in the plan plus similar implementations
- **Route B**: Give it the request and ask it to find where the change should go and what patterns to follow
- **Both routes**: If the REQ's `prime_files` reference primes with a `## Lessons` section, include them in the explore context. Previous failed approaches and gotchas from this codebase area save the explorer from repeating dead ends.

If an `## Exploration` section does not already exist, append the output:

```markdown
## Exploration

[Explore agent findings — key files, patterns, concerns]

*Generated by Explore agent*
```

### Step 5.5: Scope Declaration (Routes B and C)

Before the builder starts coding, declare intent. This prevents scope drift from being discovered only at review time, after the code is already written.

**Route A:** Skip — scope is inherently constrained (single file, single change).

**Routes B and C:** Based on the plan (Route C) or exploration output (Route B), write a `## Scope` section into the REQ file:

```markdown
## Scope

**Files I will touch:**
- `src/stores/theme-store.ts` (new) — theme state management
- `src/components/settings/SettingsPanel.tsx` (modify) — add toggle
- `tests/theme-store.test.js` (new) — unit tests

**Files I will NOT touch:** [any files that seem related but are out of scope]

**Acceptance criteria (restated from REQ):**
- [ ] Dark mode toggle visible in settings
- [ ] Theme persists across page reload
- [ ] OS preference respected on first visit
```

The Scope section serves two purposes:
1. The builder commits to a file list before writing code — drift becomes measurable.
2. The acceptance criteria, restated from the REQ, become the word-by-word comparison target for review.

The review step (Step 7) **MUST** compare the Implementation Summary's file list against the Scope declaration (Routes B and C only). Any file touched that was not declared, or any declared file not touched, is flagged as scope drift (Important finding if significant, Minor if trivial like a forgotten import update). **Route A** has no Scope declaration — skip the scope-drift comparison for Route A REQs.

### Step 5.75: Pre-Flight Check (Routes B and C)

Quick environment sanity check before the builder starts coding. All checks are **warnings, not blockers.** Append findings to REQ as `## Pre-Flight` section only if issues are found — skip the section entirely if clean.

**Route A:** Skip pre-flight — too lightweight to justify the overhead.

**Routes B and C:**

1. **Git clean:** Run `git status --porcelain --untracked-files=all` (the `-uall` flag lists files inside untracked dirs individually instead of collapsing them to a single `?? dir/` row). If there are uncommitted changes unrelated to `do-work/`, warn: "Uncommitted changes detected — the commit step may stage unrelated files." List the files.
2. **Tests baseline:** If the project has a test command (check the prime file's testing section, or look for `package.json` test scripts, `pytest.ini`, etc.), run it. If tests already fail on HEAD before any changes, note this: "Baseline tests failing — builder should not be blamed for pre-existing failures." Record which tests fail.
3. **Dependencies:** If `package.json` exists but `node_modules/` doesn't, or `requirements.txt` exists without an active venv, warn: "Dependencies may not be installed."

```markdown
## Pre-Flight

**Git:** ⚠ 3 uncommitted files (src/temp.ts, .env.local, notes.md)
**Tests baseline:** ✓ All passing (47 tests)
**Dependencies:** ✓ Installed

*Checked by work action*
```

### Step 6: Implementation

**Agent rules loading:** Before spawning the implementation agent, load domain-specific rules:

1. **Always load** `crew-members/general.md` — cross-domain rules and PRIME Files Philosophy
2. **Always load** `crew-members/karpathy.md` — behavioral guardrails (think before coding, simplicity, surgical changes, goal-driven execution)
3. **Conditionally load** `crew-members/[domain].md` — normalize the REQ's `domain` frontmatter per the Schema Read Contract first (e.g., `back-end` → `backend`, `ui_design` → `ui-design`), then load if the resolved domain is set AND the file exists (e.g., `domain: ui-design` → `ui-design.md`). An unknown value after normalization emits the contract's warning and falls back to `general` — no additional domain-specific crew loads (the always-loaded `general.md` from step 1 is the base).
4. **Conditionally load** `crew-members/testing.md` — if the REQ's `tdd` frontmatter normalizes to `true` per the Schema Read Contract (accepts `test_first`/`yes`/`on`/`t` as truthy aliases), or `domain: testing`
5. **Conditionally load** `crew-members/caveman.md` — if the REQ's `caveman` frontmatter normalizes to a non-`false` value per the Schema Read Contract (any of `true`, `lite`, `full`, `ultra`, plus `yes`/`on` → `true`, `light` → `lite`). Compresses agent prose ~65-75% while keeping code and technical terms exact.
6. **If a rules file is missing**, proceed without it — never block on a missing rules file

**Approach directive assignment (multi-REQ only):** If multiple REQs are being processed in parallel, read `crew-members/approach-directives.md` and assign each sub-agent a distinct directive from the pool. Include the directive in the sub-agent's context block. Record the assigned directive in the REQ's Implementation Summary section. For single-REQ processing, no directive is needed — skip this.

Spawn a **general-purpose agent** with the loaded rules, any files listed in the `prime_files` array, and context appropriate to the route:

- **Route A**: Request content only — "triaged as simple, aim for a focused minimal change"
- **Route B**: Request + exploration output — "follow existing patterns identified above"
- **Route C**: Request + plan + exploration output — "implement according to the plan"

All routes include these instructions to the agent (pointers — the underlying rules live in the loaded crew-members files and in the REQ frontmatter the orchestrator already wrote):

- **Crew rules govern behavior:** `crew-members/general.md` (always loaded) carries the Prime Files philosophy, Lessons-discipline, test-writing posture, cross-REQ test-break rules, and Discovered-Tasks contract. `crew-members/karpathy.md` (always loaded) enforces think-before-code, surgical scope, and goal-driven execution. Domain/testing/caveman crews layer on top per Step 6's loading order. The builder reads these — do not re-state their contents inline.
- **Prime files come first:** Read every path in `prime_files` before touching code. If the primary utility you are modifying has no prime, investigate and create one (`prime-[name].md`), then update REQ frontmatter. Lessons sections in those primes encode prior mistakes — heed them.
- **P-A-U phasing is mandatory:** Edit the REQ's "AI Execution State (P-A-U Loop)" checkboxes in real time. [PLAN] writes a brief technical approach. [APPLY] stays in declared scope. [UNIFY] runs `git diff --stat`, runs native linters, verifies no debug artifacts, and lists each file checked (the orchestrator audits this in Step 6.3).
- **TDD mode when `tdd: true`:** Follow RED → GREEN → REFACTOR. Anchor RED on the REQ's `## Red-Green Proof` section if present. Report the red-green evidence (test name, failure-before, pass-after) — Step 6.5 verifies it.
- **Captured proof first:** If `## Red-Green Proof` is present, its RED prompt/case and GREEN outcome are the primary behavior tests must prove. Only adapt with documented reason.
- **Log Decisions as D-XX:** Significant implementation choices not dictated by plan/requirements become numbered entries in a `## Decisions` section. Continue numbering from the `<!-- D-XX counter: ... -->` comment Step 3.5 left behind; if none, start at D-01. Each decision needs reasoning — without it, the intent trail breaks.
- **Out-of-scope finds go to `## Discovered Tasks`** (a separate section, not nested inside Implementation Summary) — do not fix inline. Step 8 classifies and queues them.
- **Report back the file manifest:** list every source file created/modified/deleted with the action verb, plus tests touched. The orchestrator writes the formal `## Implementation Summary` from your report.
- **Standard freedoms and obligations:** Full file/shell access. Escalate to explore or plan if the work proves harder than triaged. Document blockers explicitly. Identify and run related existing tests; honor any test-command map in the prime file (takes precedence over generic detection).

### Step 6.25: Implementation Summary

After implementation completes, write a manifest of what changed to the request file. This is the primary auditability artifact — without it, there's no way to verify the REQ was implemented without digging through git history.

**If a `## Implementation Summary` section already exists** (e.g., from a re-qualification or remediation loop), **replace it entirely** with the new content. Do not append a second copy. The most recent implementation is the one that matters.

Append (or replace) in the request file:

```markdown
## Implementation Summary

**Files changed:**
- `src/stores/theme-store.ts` (new)
- `src/components/settings/SettingsPanel.tsx` (modified)
- `tests/theme-store.test.js` (new)

**What was done:** [1-2 sentences — what the implementation actually did]
```

**Rules:**
- **Mandatory for all routes.** Route A gets a short list. Route C gets a detailed list.
- List all project files that changed — source code, config (`package.json`, `Dockerfile`, CI YAML), documentation, etc. Exclude only `do-work/` metadata files.
- Mark files as `(new)`, `(modified)`, or `(deleted)`.
- The "What was done" summary should be factual, not aspirational — describe what you built, not what the REQ asked for.
- This section is the primary auditability artifact. If `Files changed` only lists `do-work/` paths or is empty, the REQ was not implemented.
- **Design-artifact exception:** For `domain: ui-design` requests that produce design deliverables rather than code (wireframes, IA specs, visual specs, interaction specs), the artifact files themselves count as project files. Place them in the project's design docs directory (e.g., `docs/design/`) — not inside `do-work/`. The Implementation Summary lists these files normally.

### Step 6.3: Qualify Implementation

After the builder returns and the Implementation Summary is written, the **orchestrator** (not the builder) independently verifies the builder's claims before proceeding. This is not self-reporting — the orchestrator reads actual output, not the builder's description of it.

**Qualification checklist:**

1. **Files exist:** For every file listed in the Implementation Summary, verify on disk. `(new)` files must exist. `(modified)` files must show in `git diff` or `git diff --staged`. `(deleted)` files must be gone. Run the commands — don't trust the summary.
2. **Changes are substantive:** For each `(new)` file, verify it is not a placeholder (more than boilerplate/empty exports/TODO comments — minimum 10 meaningful lines for source files, 3 for config). For `(modified)` files, verify the diff contains changes related to the REQ's requirements, not just whitespace or import shuffling.
3. **Requirements traced:** Re-read the REQ's What/Detailed Requirements section. For each stated requirement, confirm at least one file in the Implementation Summary plausibly addresses it (by filename and diff content). Flag any requirement with no corresponding file change.
4. **P-A-U box audit:** Read the REQ's AI Execution State section. If any box is still `[ ]`, the builder did not complete that phase — flag it. If `[UNIFY]` is checked but the diff contains debug artifacts (`console.log`, `print()`, `debugger`, TODO/FIXME added by this change), un-check it and flag.
5. **Wired:** For each `(new)` source file, verify it is imported or referenced by at least one other file in the project (grep for the filename or an exported symbol). A new component/module that nothing imports is dead code — flag it. **Exceptions** (do not flag): Entry points (e.g., `main.ts`, `index.html`), config files, test files, standalone scripts, framework-convention files discovered by file-system routing (e.g., Next.js `pages/`/`app/` routes, SvelteKit `routes/`, Remix `routes/`, Nuxt `pages/`, Astro `pages/`), files re-exported through a barrel index (`index.ts`/`index.js` that re-exports them), files that are side-effect-only imports (CSS modules, polyfills, global stylesheets imported for their side effects), and files used exclusively via dynamic import (`import()` or `require()` with a variable path) where static grep won't find a reference.
6. **Flowing:** For files that handle data (API endpoints, data stores, handlers, services), verify the data path isn't hardcoded or stubbed. Check for: hardcoded empty arrays `return []`, placeholder strings like `"TODO"` or `"placeholder"`, `return null` in data-fetching functions, commented-out database calls. If found, flag as hollow implementation — the file exists and is wired but doesn't actually do anything.

**Anti-rationalization rules** (apply when evaluating the above):

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The summary says files changed" | Check the file system | The summary is a claim, not evidence |
| "Tests pass so requirements are met" | Compare requirements to diff, word by word | Tests can be incomplete |
| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |
| "This works on my test case" | Test at least 2 additional cases including an edge case | One test case proves nothing about generality |
| "The existing code was already like this" | Flag it in Discovered Tasks | Pre-existing problems are still problems |
| "It's just a small deviation from the plan" | Log it as a Decision (D-XX) | Unlogged deviations break traceability |

## Red Flags

- REQ in `do-work/working/` for >1 hour with no new git commits (builder may be stuck)
- Implementation Summary lists files but `git diff` shows no changes in those files (hollow implementation)
- All P-A-U checkboxes marked complete but diff contains `console.log`, `debugger`, or `TODO` (debug artifacts)
- No Triage section appended to the REQ after processing begins
- Scope section declares 3 files but Implementation Summary lists 12 (scope creep)
- Builder created files only inside `do-work/` and no source files changed (no real work done)

## Verification Checklist

- [ ] All pending REQs processed or explicitly skipped with documented reason
- [ ] Every completed REQ has an Implementation Summary section with file manifest
- [ ] No REQ files remain in `do-work/working/` after the work loop ends
- [ ] CHECKPOINT.md written if ending mid-session (for resume)
- [ ] Git commit created for each completed REQ
- [ ] Cleanup pass triggered at end of work loop

**If qualification fails on any check:**
1. Append a `## Qualification` section to the REQ noting what failed and why.
2. Return to Step 6 — spawn the builder again with the specific failures as context.
3. Maximum **2 re-qualification attempts**. After that, note remaining issues and proceed to Testing (Step 6.5). The review step will catch what remains.

**If qualification passes:**
- Append a brief `## Qualification` section: "Passed — [N] files verified, [N] requirements traced, P-A-U confirmed."
- Proceed to Step 6.5.

### Step 6.5: Testing

Before marking complete, verify tests pass:

1. **Check the prime file for test guidance** — if the REQ's `prime_files` reference a prime with a testing section (test commands, code-area-to-test mappings), use that as the primary source for what to run. **Before running, verify each listed command still exists**: for npm scripts check it's present in `package.json`; for other tools verify the config file exists (`jest.config.*`, `pytest.ini`, `Cargo.toml`, etc.). If a prime test command is no longer valid, fall back to generic detection for that command and note: `Prime test command '[cmd]' not found — falling back to generic detection.` Prime test maps are project-specific knowledge that generic detection can't replicate (e.g., "changes to `lib/inpainting.js` require `npm run test:api`" or "`npm test` is always safe but `npm run test:e2e` costs money").
2. **Fall back to generic detection for unmapped files** — if the prime has no testing section, or if you changed files the prime's test map doesn't cover, fall back to generic detection for those files: look for `package.json` test scripts, `jest.config.*`, `pytest.ini`, `Cargo.toml`, `*_test.go`, etc. A partial prime map is not an excuse to skip tests — matched files use the prime's commands, unmatched files use generic detection. If neither source yields test commands for a file, skip testing for it and note it.
3. **Run relevant tests** — target tests related to changed code, not the full suite (unless it's fast). If the prime specifies different commands for different code areas, run only the commands relevant to the files you changed. For unmapped files, run whatever generic detection found.
4. **If tests fail** — check whether the failures were already recorded as baseline failures in Step 5.75 (Pre-Flight). If a failing test matches a pre-existing baseline failure (same test name/file, same failure mode), exclude it from the pass/fail gate — the builder should not be blamed for pre-existing failures. Only **new regressions** (tests that passed at baseline but fail after implementation) require fixing. Return to implementation to fix new regressions. On attempt 2+, load `crew-members/debugging.md` and `crew-members/testing.md` for the builder to follow the structured debugging methodology and review test quality. Loop until passing or mark as failed after 3 attempts.
5. **If new tests are needed** — spawn a general-purpose agent to write them following existing patterns, then run them.

Append to the request file:

```markdown
## Testing

**Tests run:** [command]
**Result:** ✓ All passing (X tests)

**Red-green validation:** *(for bug fixes and new features)*
- [test name/file]: ✗ before implementation → ✓ after
- [test name/file]: ✗ before implementation → ✓ after

**New tests added:**
- [list]

**Existing tests updated (cross-REQ impact):**
- [test file] (from REQ-NNN): [what changed and why — intentional behavior change]

*Verified by work action*
```

Omit `Red-green validation` if no request-specific tests were written or identified, or if the change is non-behavioral (refactor, config, docs, cleanup) — use regression evidence instead. Omit `Existing tests updated` if no prior tests were modified.

When the REQ includes `## Red-Green Proof`, the `Red-green validation` entries should trace back to that captured RED/GREEN pair. If the implemented test uses a nearby equivalent instead of the exact captured prompt/case, explain why.

**TDD verification:** If the REQ has `tdd: true`, the `Red-green validation` section is mandatory — the builder must show test-first evidence that they used RED/GREEN TDD (test written before implementation, failed, then passed after). If this evidence is missing, treat it as a test failure: return to implementation (same path as step 4 above) with explicit instructions to provide red/green evidence — write the failing test first, confirm it fails, then make it pass.

### Step 7: Review

Run the review work action in **pipeline mode** against this REQ.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?).

**How to run it:** Spawn an agent with the review work action file, the REQ path, and the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain has a matching file, load it; otherwise skip). Or read the review work action file and follow its pipeline mode instructions in the current session.

**What happens next depends on the review result:**

- **Acceptance = Pass AND overall ≥ 75%**: Append the Review section to the REQ and continue to archive as `completed`. Minor findings go in the report only.
- **Acceptance = Partial OR overall 50-74%**: Append Review, continue to archive as `completed`, but the review **MUST** create follow-up REQs for every Important finding. These are not optional — they enter the queue and block the UR from being considered "done" until addressed.
- **Acceptance = Fail OR overall < 50%**: **Do NOT archive as completed.** Instead:
  1. Append the Review section to the REQ.
  2. Return to Step 6 (Implementation) with the review findings as context for the builder. Load `crew-members/debugging.md` for the remediation attempt — the builder needs structured debugging methodology, not just "try again."
  3. The builder gets **ONE remediation attempt**.
  4. Re-run Steps 6.25 through 7 (Summary → Qualification → Testing → Review) on the remediated code.
  5. If still failing after remediation: update frontmatter to `status: completed-with-issues`, `completed_at: <timestamp>`, append a `## Remediation` section documenting both attempts, and create follow-up REQs for all remaining Important findings. Then proceed to archive (Step 8) — the frontmatter is already set, so Step 8 should not overwrite it.

The status `completed-with-issues` means the REQ was archived but has known unresolved problems. It counts toward UR completion for archiving purposes, but the follow-up REQs must be processed before the work is considered ship-ready. This status is visible in the recap and present-work actions.

**Follow-up REQs are created based on finding severity, not score.** The review creates follow-up REQs for each **Important** finding (regardless of overall score). Minor and Nit findings go in the report only. The follow-up REQs enter the queue and get processed in a future loop iteration. Follow-up REQs created by the review step must include: `status: pending`, `user_request: [same UR as the reviewed REQ]`, `addendum_to: [reviewed REQ id]`, `domain: [same domain]`, and `review_generated: true`. Place them in `do-work/queue/`. Cycle detection (Step 8, substep 5) applies to these follow-ups — check the `addendum_to` chain before creating.

**Calibrate depth to route:** Route A gets a quick scan (skip dimensions that don't apply). Route B gets a standard review. Route C gets a thorough review comparing against the plan.

Append to the request file:

```markdown
## Review

**Overall: [X]%** | [timestamp]

| Dimension | Score |
|-----------|-------|
| Requirements | X% |
| Code Quality | X% |
| Test Adequacy | X% |
| Scope | X% |
| Risk | [level] |
| Acceptance | [result] |

**Findings:** [count] important, [count] minor
**Acceptance:** [Pass/Partial/Fail/Untested] — [1-line summary]
**Suggested testing:** [count] items
**Follow-ups created:** [REQ-NNN, REQ-NNN] or "None"

*Reviewed by review work action*
```

### Step 7.5: Lessons Learned

Before archiving, capture what's worth remembering. This section is the institutional memory — when someone revisits this code in six months, the REQ file tells them what happened, what was tried, and why things ended up the way they did.

Append to the request file:

```markdown
## Lessons Learned

**What worked:** [1-2 bullets — approaches, patterns, or tools that paid off]
**What didn't:** [1-2 bullets — dead ends, failed approaches, and *why* they failed]
**Worth knowing:** [Anything the next person touching this code should know — gotchas, edge cases, non-obvious dependencies]
```

**Rules:**
- Keep it concise — pointers to code, not walls of text. The code is the source of truth.
- **Required for Routes B and C** — there's always something worth recording when exploration or planning was involved. **Optional for Route A** — skip if the change was straightforward with no unexpected discoveries, no failed approaches, and no gotchas worth noting. If anything surprised you (undocumented behavior, unexpected test failures, a file that wasn't where you expected), record it.
- "What didn't work" is the most valuable part — it prevents repeating mistakes.
- File lists are no longer needed here — they're covered by the mandatory Implementation Summary (Step 6.25).

**Update prime files (deferred to Step 8):** After writing the Lessons Learned section, check the REQ's `prime_files` frontmatter. For each listed prime file relevant to this lesson, **collect a pending prime-link write** — do NOT execute the write here. The REQ is still in `do-work/working/`, so any link pointing to its eventual archive location would either be broken or tempt a link to the transient working path.

Record each pending write as a tuple: `{ primeFilePath, relativeLinkText, lessonSummary }`. Hold them in memory (or a small scratch file under `do-work/working/`) until Step 8.

**Path computation rule (for use in Step 8):** the link path must be relative to the prime file's location, not the repo root. Count how many directories deep the prime file sits (i.e., the number of path components before the filename). Prepend that many `../` steps to the REQ's repo-root-relative archive path. Examples:
- Prime at `prime-auth.md` (0 dirs deep) → `do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Prime at `src/utils/prime-auth.md` (2 dirs deep: `src/` and `utils/`) → `../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Prime at `web/src/auth/prime-auth.md` (3 dirs deep) → `../../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`

The existence-verify check on the resolved path runs in Step 8 (post-move) — that's the whole reason for deferring.

Only add a link when the lesson is relevant to that prime file's scope — don't spray every lesson into every prime file. If the REQ has no `prime_files` or the lessons aren't relevant to any prime file, skip this and clear the pending list.

**Knowledge-base handoff.** After the Lessons Learned section is written and prime-file links are in place, follow `actions/kb-lessons-handoff.md` to offer dropping a structured source document into `kb/raw/inbox/` so the next `bkb triage` + `bkb ingest` cycle compiles the lessons into the wiki. The handoff asks the user before writing and records `kb_status` (plus `kb_entry` on success) back onto the REQ. In unattended pipeline runs with no human in the loop, the handoff defaults to `kb_status: pending` — it never writes to the KB without consent. If the project has no `kb/` directory, the handoff points the user at `do-work bkb init` and defers; it never blocks archival.

### Step 8: Archive

**On success:**

1. Update frontmatter: `status: completed`, `completed_at: <timestamp>`
2. Verify `## Implementation Summary` is present (written in Step 6.25). If missing, append it now — this should not happen in normal flow, but crash recovery may skip it.
3. **Create follow-ups for builder-decided questions:** If the REQ has any `- [~]` items in Open Questions where the builder's choice affects what the user sees or interacts with, create a follow-up REQ for each. **Create follow-ups for:** UX decisions (interaction behavior, visibility, layout), scope boundaries (what's included/excluded), data representation choices. **Skip follow-ups for:** purely technical decisions (caching strategy, algorithm choice, internal naming, DB indexes) that don't change user-facing behavior. Template:
   ```markdown
   ---
   id: REQ-NNN
   title: "Confirm: [brief description of the choice]"
   status: pending-answers
   created_at: [timestamp]
   user_request: [same UR as the original REQ]
   addendum_to: [original REQ id]
   builder_decided: true
   ---

   # Confirm: [Brief Description]

   ## What
   During [REQ-id], the builder chose [choice] for [question]. This follow-up
   confirms whether that choice matches your intent or if you'd prefer a different approach.

   ## What the Builder Chose
   [Description of the choice and its impact on the implementation]

   ## What Would Change
   [If the user picks a different option, what would need to change]

   ## Open Questions
   - [ ] [Original question]
     Recommended: [builder's choice — already implemented]
     Also: [other alternatives]
   ```
   These go in `do-work/queue/` with `status: pending-answers`. The user reviews them via `do-work clarify`.
4. **Queue Discovered Tasks:** Check the REQ file for a `## Discovered Tasks` section (appended by the implementation agent as a separate section — not inside `## Implementation Summary`). For every item listed, classify by severity and create follow-up REQs accordingly.

   The builder should classify each discovered task when appending them:
   ```
   ## Discovered Tasks
   - **[critical]** SQL injection vulnerability in user search endpoint
   - **[normal]** Dead code in utils/legacy-parser.ts can be removed
   - **[low]** Variable naming inconsistency in auth module
   ```

   If the builder did not classify them, the orchestrator classifies based on:
   - **critical**: Security vulnerability, data loss risk, broken functionality in production paths
   - **normal**: Technical debt, missing tests, minor bugs in non-critical paths
   - **low**: Style issues, naming, dead code, documentation gaps

   **For `[critical]` discoveries:** Create follow-up REQ with `status: pending` (not `pending-answers`) — these skip user confirmation and go straight into the work queue. Add a note in Open Questions: `- [x] Auto-approved: critical severity (security/data/production risk). → Added to queue immediately.` Report prominently: `⚠ CRITICAL discovered: [description] — auto-queued as REQ-NNN`

   **For `[normal]` and `[low]` discoveries:** Use the existing `pending-answers` flow:
   - Set frontmatter: `status: pending-answers`, `user_request: [same UR]`, `addendum_to: [current REQ id]`, `domain: [same domain as current REQ]`.
   - Add an `## Open Questions` section with this checkbox format:
     `- [ ] I discovered this out-of-scope task while working on [current REQ]: [Task Description]. Should I process this as a new task?`
     `  Recommended: Yes, add to queue (will flip to 'pending').`
     `  Also: No, discard it.`
   This ensures non-critical discoveries require the user's explicit permission via `do-work clarify` before execution.
5. **Cycle detection:** Before creating any follow-up REQ, verify the current REQ's own `addendum_to` chain is not already circular. Algorithm: walk `addendum_to` links (honoring the `amends`/`parent`/`amendment_to` alias per the Schema Read Contract when the canonical key is absent) starting from the current REQ, collecting each visited ID into a seen set. If you encounter the current REQ's ID again during the walk, the chain is already circular — do not create any follow-ups. Report: `⚠ Cycle detected in addendum_to chain: REQ-NNN → REQ-MMM → ... → REQ-NNN. Skipping follow-up — manual resolution needed.` This handles chains of any length.
6. Archive based on REQ type:

| REQ has... | Archive behavior |
|------------|-----------------|
| `user_request: UR-NNN` | Check if ALL REQs in the UR are finished (status: `completed`, `completed-with-issues`, or `failed`). Check `do-work/queue/`, `do-work/working/`, `do-work/archive/` root, and `do-work/archive/UR-NNN/` for REQs belonging to this UR. If all finished: move completed/completed-with-issues REQs into UR folder (failed REQs stay at archive root), move entire UR folder to `archive/`. If any REQ is still `pending`, `pending-answers`, or `claimed`: move this REQ to `archive/` root; UR stays in `user-requests/` until last REQ finishes. |
| `context_ref` (legacy) | Move REQ to `archive/`. If all related REQs are now archived, move the CONTEXT doc too. |
| Neither (standalone legacy) | Move directly to `archive/`. |

7. **Execute deferred prime-link writes (from Step 7.5):** Now that the REQ is at its final archive path, walk the `pendingPrimeLinkWrites` collected during Step 7.5. For each pending entry:
   - Compute the relative path from the prime file to the REQ's actual archived location (UR folder if the UR was just consolidated, or `archive/` root if the UR is incomplete).
   - Verify the resolved path points to an existing file. If it doesn't, report the broken link and skip — do NOT silently write a broken link.
   - Append the link to a `## Lessons` section in the prime file (create the section if it doesn't exist):
     ```markdown
     ## Lessons

     - [REQ-NNN: 1-line summary of the lesson](<relative-path-to-archived-req>#lessons-learned)
     ```
   - Stage the prime file along with the implementation files in Step 9.

   This is the post-move execution that makes the existence-verify meaningful — Step 7.5 only collected; the writes happen here.

**On failure:**

Before classifying via the symptom table below, **check for upstream failure**. Cascades from a failed prerequisite often present as plausible-looking `code` or `spec` symptoms in the downstream REQ; misclassifying them sends the builder chasing phantom bugs in the wrong domain.

**Upstream-failure short-circuit:**

Read the frontmatter of every REQ this one depends on:
- `addendum_to` (single parent, if set; or `amends`/`parent`/`amendment_to` as the legacy alias if `addendum_to` is absent — same back-compat shape as the `depends_on`/`dependencies:` pair; `addendum_to` wins when multiple are present)
- every entry in `depends_on` (if set, or every entry in the legacy `dependencies:` alias if `depends_on` is absent — same back-compat rule as Step 1; `depends_on` wins when both present)

Resolve each ID by globbing `do-work/archive/**/REQ-NNN-*.md`, `do-work/archive/**/REQ-NNN.md`, `do-work/queue/REQ-NNN-*.md`, and `do-work/working/REQ-NNN-*.md`. If any referenced REQ has `status: failed`, skip the symptom table and short-circuit classification:

- `status: failed`
- `error_type: spec` (the local approach is downstream-correct only if the upstream is correct; with the upstream broken, the local spec is implicitly unsound)
- `error: "Upstream REQ-NNN failed (error_type: <ancestor.error_type>); downstream blocked. Original error: <original error message>"`

Create the follow-up REQ per the Spec row below. It inherits `addendum_to: <this failed REQ>`; the cascade is now visible in the addendum chain and the follow-up's error description names the upstream root cause. The follow-up should also carry the original dependency list so it re-blocks on the same upstream until the upstream's own follow-up lands — and it always emits the canonical `depends_on:` key, even if the failed REQ used the legacy `dependencies:` alias. Don't propagate the alias on follow-ups.

If no upstream REQ is `failed`, fall through to the symptom-based classification table:

| Type | Symptoms | Recovery |
|------|----------|----------|
| **Intent** | Requirements are ambiguous or contradictory; builder couldn't determine what to build | Create a follow-up REQ with `status: pending-answers` containing the specific ambiguities as Open Questions. Archive original as `failed` with `error_type: intent`. |
| **Spec** | Requirements are clear but the technical approach was wrong (wrong files, wrong pattern, wrong architecture) | Create a follow-up REQ with a `## Prior Attempt` section summarizing what was tried and why it failed. Set `status: pending`. Archive original with `error_type: spec`. |
| **Code** | Approach was right but implementation has bugs (tests fail, runtime errors, logic errors) | Create a follow-up REQ targeting the specific code issue. Set `status: pending`. Archive original with `error_type: code`. |
| **Environment** | External dependency unavailable, permissions issue, tooling broken | No follow-up REQ — user must fix the environment. Archive with `error_type: environment` and a clear description of what's needed. |

**Anti-rationalization addition.** When checking the symptom table:

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This REQ failed on a code bug" | Check whether any `addendum_to` or `depends_on` ancestor (or `dependencies:` alias) is also `failed` first | Downstream failures often inherit upstream rot; misclassifying as `code` chases phantom bugs in the wrong domain |

**Procedure:**
1. Run the upstream-failure short-circuit. If it fires, jump to step 3.
2. Otherwise classify using the symptom table above.
3. Update frontmatter: `status: failed`, `error: "description"`, `error_type: [intent|spec|code|environment]`
4. For Intent/Spec/Code failures: create the appropriate follow-up REQ (details above). Set `addendum_to` to the failed REQ's ID so context chains. Preserve the original dependency list on the follow-up when the failure was upstream-driven — always emit it under the canonical `depends_on:` key, even if the failed REQ used the legacy `dependencies:` alias.
5. Move to `archive/` (failed REQs always go to archive root, not into UR folders).
6. Report to user: `[REQ-NNN] failed ([type]): [description]. Follow-up: [REQ-NNN] / None.` When the short-circuit fired, prefix the report with `(upstream cascade — original failure at REQ-NNN)`.

### Step 9: Commit (Git repos only)

Check for git with `git rev-parse --git-dir 2>/dev/null`. If not a git repo, skip.

```bash
# Stage implementation files + archived REQ
git add src/stores/theme-store.ts src/components/settings/SettingsPanel.tsx \
  do-work/archive/UR-002/REQ-003-dark-mode.md

# Stage follow-up REQs created in Step 8 (if any)
git add do-work/queue/REQ-025-confirm-sidebar-palette.md

# Stage UR-folder move (if this was the last REQ and the UR moved to archive/)
# Both the old path (deletion) and new path (addition) must be staged.
git add do-work/user-requests/UR-002/ do-work/archive/UR-002/

git commit -m "$(cat <<'EOF'
[REQ-003] Dark Mode (Route C)

Implements: do-work/archive/UR-002/REQ-003-dark-mode.md

- Created src/stores/theme-store.ts
- Modified src/components/settings/SettingsPanel.tsx

EOF
)"
```

**Format:** `[{id}] {title} (Route {route})` + `Implements:` line + summary bullets. Add a co-author trailer if your platform convention calls for one (e.g., `Co-Authored-By: Agent <agent@example.com>`), otherwise omit.

One commit per request. Stage all files created, modified, moved, or deleted during this request's lifecycle: implementation files (listed in the Implementation Summary), the archived REQ file, any follow-up REQs created in Step 8 (`pending-answers` files in `do-work/queue/`), and any UR-folder moves to `archive/`. If Step 8 substep 7 wrote prime-file lessons links, the modified prime files must also be staged — they are part of the REQ's lifecycle changes even though they aren't listed in the Implementation Summary's `Files changed`. Do not use `git add -A` or `git add .` — these risk staging secrets, `.env` files, or unrelated changes. Don't bypass pre-commit hooks — fix issues and retry. Failed requests get committed too.

**Validation check (successful REQs only):** Before committing, compare the `## Implementation Summary` file list against the staged files (excluding `do-work/` paths). If the Implementation Summary lists files that aren't staged, or if the only staged files are `do-work/` metadata, flag the mismatch — the commit may not contain the actual implementation. Fix the staging or update the Implementation Summary before proceeding. Design-artifact files placed outside `do-work/` satisfy this check — they are project deliverables. **Skip this check for failed REQs** — they may have no Implementation Summary or no project files staged, and that's expected.

**Write commit hash back to the archived REQ.** After the commit succeeds, retrieve the hash with `git rev-parse --short HEAD` and update the archived REQ's frontmatter `commit:` field with the actual value. Then create a **separate metadata commit** (do not amend — amending changes the hash and invalidates what you just wrote):

```bash
# After the initial commit succeeds:
HASH=$(git rev-parse --short HEAD)
# Update the commit: field in the archived REQ frontmatter
# (replace "commit:" line or add it if missing)
git add do-work/archive/UR-NNN/REQ-NNN-slug.md
git commit -m "[REQ-NNN] record commit hash ${HASH}"
```

This ensures the `commit:` field in the archived REQ contains the real implementation commit hash, which the review-work and present-work actions depend on for traceability. The metadata commit is a lightweight bookkeeping entry — it does not contain implementation changes.

### Step 10: Loop or Exit

Re-check `do-work/queue/` for `REQ-*.md` files (fresh check, not cached).

- **Dependency-ready `pending` REQs found**: **CONTEXT WIPE** (see below). Then loop to Step 1.
- **No dependency-ready `pending` REQs remain** (queue may still have dependency-blocked or held REQs): Write a **Session Checkpoint** (see below), run the cleanup action, then report the final summary using the **same composed structure** as Step 1's "Exit paths when no `pending` REQs found" — render the completed/done section, the pending-answers section, the blocked-archive-collision section, and the blocked-by-dependencies section in that order, including only those that have at least one REQ. If none of the four sections applies (queue is fully empty), report completion and exit. Mixed cases render all applicable sections in one summary.

#### Context Wipe — Verified

Before looping to Step 1 for the next REQ:

1. **Fresh agents:** Spawn a NEW agent for the next REQ. Do not reuse the previous builder/explorer/planner agent. Each REQ gets clean agents with no carried-over context.
2. **Explicit declaration:** State in your progress message: `Context wipe: previous REQ was [REQ-NNN] working on [files]. Now starting fresh for next REQ.`
3. **Contamination check:** When the next REQ's builder returns its Implementation Summary (Step 6.25), compare the file list against the *previous* REQ's Implementation Summary. Unexpected overlap — files from the previous REQ appearing without an explicit `addendum_to` or `related` link — is a scope contamination signal. Flag it in the Qualification step (Step 6.3).

#### Session Checkpoint

At the end of every work session (whether all REQs completed, user stops, or session is ending), write `do-work/CHECKPOINT.md`. Scale the checkpoint to how much happened:

```markdown
---
session_ended: [timestamp]
last_completed: REQ-NNN
queue_state: [N pending, N pending-answers, N blocked-archive-collision, N in-progress]
reqs_processed_this_session: N
session_depth: light | moderate | heavy
---

# Session Checkpoint

## Completed This Session
- REQ-NNN: [title] (Route [X], [score]%)
- REQ-NNN: [title] (Route [X], [score]%)

## In Progress (interrupted)
- REQ-NNN: [title] — stopped at [phase: triage/planning/exploring/implementing/testing/reviewing]
  Last known state: [1-2 sentences]
  Key files being modified: [list]
  Known issues: [any blockers or concerns]

## Still Queued
- REQ-NNN: [title] (pending)
- REQ-NNN: [title] (pending-answers — [N] questions)

## Session Notes
[Environment issues, user preferences expressed, patterns discovered, decisions made outside REQ files]

## Context Summary (heavy sessions only)
[Recap of key decisions (D-XX references), architectural patterns encountered, and prime files
that the next session should re-read before starting. Include this section when 6+ REQs were
processed — at that volume, carried-over assumptions are unreliable.]
```

**Session depth guide:**
- **light** (1-2 REQs): Minimal checkpoint — Completed + Still Queued sections are sufficient
- **moderate** (3-5 REQs): Add Session Notes with patterns observed and environment quirks
- **heavy** (6+ REQs): Add Context Summary recapping key decisions and recommending the next session re-read prime files fresh rather than trusting carried-over assumptions

**On session start (Step 1 addition):** Before crash recovery, check for `do-work/CHECKPOINT.md`. If it exists:
1. Read it and report a brief summary: `Resuming from previous session. Last completed: REQ-NNN. [N] REQs still queued.`
2. Use the "In Progress" section to inform crash recovery context.
3. **Do not delete yet.** Keep the checkpoint until crash recovery completes successfully (all files moved out of `working/`). Then delete it. This prevents losing resume context if the session crashes again during crash recovery.

This is NOT a blocking gate. If no checkpoint exists, the session starts normally with existing crash-recovery logic.

## Clarify Questions

The clarify workflow has its own action. Run `do-work clarify` — it handles batch-review of `pending-answers` REQs, where the user confirms, overrides, or discards builder decisions. Resolved REQs flip back to `pending` and re-enter the work queue.

## Orchestrator Checklist (per request)

```
□ Step 1: Find next request (read CHECKPOINT.md if exists, crash recovery, validate frontmatter, pick first pending)
□ Step 2: Claim request (mkdir -p working/, move REQ, update status & claimed_at)
□ Step 3: Triage (decide route, append ## Triage, read original if addendum)
□ Step 3.5: Handle Open Questions (mark - [~] with D-XX numbered decisions)
□ Step 4: Plan (Route C: spawn Plan agent + validate plan / Routes A & B: note skipped)
□ Step 5: Explore (Routes B & C: spawn Explore agent, include prime file lessons)
□ Step 5.5: Scope Declaration (Routes B & C: declare files + acceptance criteria in REQ)
□ Step 5.75: Pre-Flight Check (Routes B & C: git clean, test baseline, dependencies)
□ Step 6: Implement (spawn agent with lessons + TDD mode if set, log decisions as D-XX)
□ Step 6.25: Implementation Summary (append file manifest — mandatory for all routes)
□ Step 6.3: Qualify (orchestrator verifies: files exist, substantive, wired, flowing, requirements traced, P-A-U audit)
□ Step 6.5: Test (run relevant tests, load debug rules on attempt 2+, verify TDD evidence if tdd:true)
□ Step 7: Review (spawn review action — gate on acceptance: Pass→archive, Fail→remediate with debug rules)
□ Step 7.5: Lessons Learned (append section, update prime files, skip for Route A if no surprises)
□ Step 8: Archive (update status, classify failures, triage discovered tasks, cycle-check follow-ups, queue follow-ups, move to archive/)
□ Step 9: Commit (stage explicit files, commit if git repo, write hash to REQ in separate metadata commit)
□ Step 10: Loop or Exit (context wipe + contamination check if looping, else write CHECKPOINT.md with depth + cleanup)
```

## Common Mistakes

- Spawning implementation agent without first moving file to `working/`
- Letting spawned agents handle file management (only the orchestrator moves/archives files)
- Forgetting to update status in frontmatter (only two transitions: `claimed` at Step 2, final status at Step 8)
- Archiving a UR folder before all its REQs are complete
- Forgetting Planning status note for Routes A/B ("Planning not required")
- Using `git add -A` instead of staging specific files
- Using `--no-verify` to bypass a failing pre-commit hook instead of fixing the issue
- Committing without validating Implementation Summary file list against staged files
- Implementation Summary that only lists `do-work/` paths (means the REQ wasn't actually implemented — exception: `domain: ui-design` design artifacts placed in project directories like `docs/design/`)
- Creating follow-ups for every `- [~]` item instead of only UX-affecting decisions

## Error Handling

| Phase | Action |
|-------|--------|
| `pending-answers` REQs remain after queue is empty | Report them to the user: list each REQ and its unresolved questions. Suggest `do-work clarify` to batch-review. |
| Plan agent fails (Route C) | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Explore agent fails (B/C) | Proceed to implementation with reduced context — builder can explore on its own |
| Implementation fails | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Tests fail repeatedly | After 3 fix attempts, classify as Code failure, create follow-up REQ with test failure details, archive as failed |
| Review: Acceptance = Fail | Return to Step 6 for ONE remediation attempt, then re-review. If still failing: archive as `completed-with-issues` with follow-up REQs |
| Review work agent fails | Skip review, note it in the REQ file, continue to archive — review failure is not a gate |
| Commit fails | Investigate the error (usually a pre-commit hook failure). Fix the underlying issue, re-stage, and retry as a **new** commit. Do NOT use `--no-verify` to skip hooks or `--no-gpg-sign` to bypass signing — fix the root cause. If unfixable, report the error to the user and continue to next request — changes remain uncommitted but archived. |
| Unrecoverable error | Stop loop, report clearly, leave queue intact for manual recovery |

## Progress Reporting

Keep the user informed with this format:

```
Processing REQ-003-dark-mode.md...
  Triage: Complex (Route C)
  Open Questions: 2 found → builder decided (follow-ups queued)
  Planning...     [done]
  Scope...        [done] 4 files declared
  Exploring...    [done]
  Implementing... [done]
  Summary...      [done] 3 files changed
  Qualifying...   [done] ✓ files verified, requirements traced
  Testing...      [done] ✓ 12 tests passing
  Reviewing...    [done] 92% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → abc1234

Processing REQ-004-fix-typo.md...
  Triage: Simple (Route A)
  Implementing... [done]
  Summary...      [done] 1 file changed
  Qualifying...   [done] ✓ verified
  Testing...      [done] ✓ 3 tests passing
  Reviewing...    [done] 88% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → def5678

All 2 requests completed:
  - REQ-003 (Route C) → abc1234 [review: 92%]
  - REQ-004 (Route A) → def5678 [review: 88%]
```

## What This Action Does NOT Do

- Create new request files (use the capture requests action)
- Make architectural decisions beyond what's in the request
- Run without user present (this is supervised automation)
- Modify already-completed requests
- Allow external modification of files in `working/` or `archive/`

## Archived Request File Example

See [sample-archived-req.md](./sample-archived-req.md) for a complete example of what an archived REQ looks like after processing through the full pipeline (Route B). Every section shown there is generated by the steps above.

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = implementation time. Route + timestamps let you calibrate triage accuracy over time.
