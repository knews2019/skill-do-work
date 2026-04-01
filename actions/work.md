# Work Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to process the queue. Processes requests from the `do-work/` folder in your project.

An orchestrated build system that processes request files created by the capture requests action. Uses complexity triage to route simple requests straight to implementation and complex ones through planning and exploration first.

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary (mandatory file manifest), Testing, Review. This ensures full traceability — what was planned vs done, what files were touched, and whether triage was accurate.

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
├── REQ-018-pending-task.md       # Pending (root = queue)
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

- **Root**: The queue — only pending `REQ-*.md` files
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
prime_files: []  # list paths to relevant prime-*.md files, or leave empty
created_at: 2025-01-26T10:00:00Z
user_request: UR-001          # May be absent on legacy REQs

# Set by work action when claimed
claimed_at: 2025-01-26T10:30:00Z
route: A | B | C

# Set by work action when finished
completed_at: 2025-01-26T10:45:00Z
status: completed | completed-with-issues | failed
commit: abc1234               # If git repo
error: "Description"          # Only if failed
---
```

**Status flow (frontmatter values):** `pending` → `claimed` → `completed` / `completed-with-issues` / `failed`

The intermediate phases (planning, exploring, implementing, testing, reviewing) are tracked by which `##` sections exist in the REQ file, not by frontmatter status changes. Only three status transitions are written to frontmatter: `pending` → `claimed` (Step 2), then `claimed` → final status (Step 8).

**Special status:** `pending-answers` — a follow-up REQ whose Open Questions need user input before it can be worked. These accumulate in the queue and get batch-reviewed when the user runs `do work clarify`.

## Workflow

**The work action is an orchestrator.** You handle ALL file management (moving files, updating frontmatter, appending sections, archiving). Spawned agents handle implementation work only.

### Step 1: Find Next Request

**Crash Recovery:** Before checking the queue, look inside `do-work/working/` for any `REQ-*.md` files. If any exist, a previous run was interrupted. For each recovered REQ:
1. Reset frontmatter: set `status` to `pending` (but if the REQ was `pending-answers` before being claimed — check for `## Open Questions` with unresolved `- [ ]` items — restore to `pending-answers` instead). Remove `claimed_at` and `route`.
2. Strip sections generated during the interrupted run: remove `## Triage`, `## Exploration`, `## Plan`, `## Scope`, `## Pre-Flight`, `## Implementation Summary`, `## Qualification`, `## Testing`, `## Review`, `## Lessons Learned`, `## Decisions`, and `## Discovered Tasks` sections (and their content) if present — these may be incomplete or stale from the crash. Leave `## Open Questions` and user-authored content intact.
3. Move the REQ back to the `do-work/` root

Once `working/` is empty, proceed with finding the next request.

Glob for `do-work/REQ-*.md` (root of `do-work/`, **not** a subdirectory — there is no `queue/` folder). Sort by number. Read the frontmatter of each (in number order) to check `status` — pick the first with `status: pending` (skip `pending-answers` — those wait for user input). Don't read the full body at this stage. If no `pending` REQs found, report completion and exit. If only `pending-answers` REQs remain, report them to the user so they can batch-review the questions.

**REQ validation:** When reading each REQ's frontmatter, verify it has the required fields (`id`, `status`, `title`). If a REQ file has missing or unparseable frontmatter, skip it and report: `⚠ Skipping [filename]: missing required frontmatter ([field]).` Do not let a single malformed REQ block the entire work loop — skip it and continue to the next.

**Exact glob pattern:** `do-work/REQ-*.md` — if this returns no results, do NOT conclude the queue is empty. Verify by listing `do-work/` contents to rule out a bad pattern.

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
3. Number decisions sequentially per REQ (D-01, D-02, D-03...). Open Questions decisions and Implementation Decisions (Step 6) share the same D-XX ID space — if Open Questions uses D-01 through D-03, the first implementation decision is D-04. These IDs can be referenced by future REQs.
4. Proceed with implementation using those decisions.

The follow-up REQs for builder-decided questions are created during **Step 8 (Archive)** — not here. Step 3.5 just records the decisions; the archive step handles the paperwork after the REQ is fully complete.

**Why not block?** Human time is the bottleneck. The optimal windows for user interaction are: (1) capture time, when the user is actively fleshing out requests, and (2) batch-review time, when the user returns to answer accumulated questions. Blocking mid-build wastes builder capacity on idle waiting.

**`pending-answers` REQs:** These accumulate in the queue. When the user returns, they run `do work clarify` to review all `pending-answers` REQs at once, answer the questions, and flip the status to `pending` so the next work run picks them up. The work loop skips `pending-answers` REQs — it only processes `pending` ones.

If all `- [ ]` items are already `[x]` or `[~]`, or no Open Questions section exists, skip this step entirely.

### Step 4: Planning (Route C only)

**Route C:** Spawn a **Plan agent** with the request content, project context, the `rules-[domain].md` file (if domain is missing or the file doesn't exist, skip loading it), and any files listed in the `prime_files` array. Instruct it to use the prime files as the strict index for discovering the source of truth. Do not load global architecture. Ask it to produce a specific implementation plan (files to modify, order of changes, architectural decisions, testing approach). If a `## Plan` section does not already exist, append the output:

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

1. **Git clean:** Run `git status --porcelain`. If there are uncommitted changes unrelated to `do-work/`, warn: "Uncommitted changes detected — the commit step may stage unrelated files." List the files.
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

Spawn a **general-purpose agent** with the `agent-rules/rules-general.md` file (always loaded — contains PRIME Files Philosophy and cross-domain rules) plus the `rules-[domain].md` file (if domain is set and the file exists), any files listed in the `prime_files` array, and context appropriate to the route:

- **Route A**: Request content only — "triaged as simple, aim for a focused minimal change"
- **Route B**: Request + exploration output — "follow existing patterns identified above"
- **Route C**: Request + plan + exploration output — "implement according to the plan"

All routes include these instructions to the agent:

```
- **Prime Files:** If `prime_files` are attached to this REQ, READ THEM FIRST. They are your map to the codebase. If NO prime file exists for the primary utility you are modifying, you MUST investigate the utility and create one (`prime-[name].md`). Prime files must be low-noise, high-value, point to code as the source of truth, and avoid volatile metrics (like test counts). If you create one, update the REQ's frontmatter to include it.
- **Lessons:** If any prime file you loaded has a `## Lessons` section, read the linked REQ lessons before implementing. These are previous mistakes and discoveries from this exact area of the codebase. Pay particular attention to "What didn't work" entries — these prevent repeating failed approaches. If a lesson directly contradicts your planned approach, note the conflict in your PLAN phase and explain why you're proceeding differently (or adjust your plan).
- You have full access to edit files and run shell commands
- If you find the request is more complex than expected, you can explore or plan as needed
- Document any blockers clearly
- Identify existing tests related to your changes
- **Check the prime file for a testing section** — if the prime maps code areas to specific test commands (e.g., "changes to lib/inpainting.js → run `npm run test:api`"), follow that mapping. This takes precedence over generic test detection.
- **Write pragmatic tests:** For bug fixes and new features, prefer red-green validation — write or identify tests that validate the request's requirements, run them before implementing (they should fail), then verify they pass after. For refactors, config changes, documentation, and cleanup, red-green may not apply — targeted regression tests, lint/build validation, or non-regression evidence is sufficient. The goal is proof that the change works, not ceremony.
- Write new tests for new functionality / regression tests for bug fixes
- Update existing tests if behavior intentionally changed
- **If existing tests break:** When your changes cause tests from a prior request to fail, determine if the behavior change is intentional. If yes: update the failing tests to match the new behavior and document which REQ's tests changed and why in the Testing section — this creates traceability for which request altered which other request's behavior. If no: fix your implementation to preserve the existing behavior.
- When complete, report back: list every source file you created, modified, or deleted (with the action — new/modified/deleted), and summarize what tests exist and what new tests were written. The orchestrator uses this to write the formal `## Implementation Summary`.
- **State Machine Updates:** As you progress, you MUST physically edit this REQ file to change the `[ ]` checkboxes in the "AI Execution State (P-A-U Loop)" section to `[x]`.
- **TDD Mode:** If the REQ has `tdd: true` in frontmatter, follow the red-green-refactor cycle:
  1. **RED:** Write a failing test that validates the requirement. Run it — confirm it fails.
  2. **GREEN:** Write the minimum code to make the test pass. No more.
  3. **REFACTOR:** Clean up while tests stay green.
  Report the red-green evidence in your output: test name, failure message before, pass after. The Testing section (Step 6.5) will verify this evidence is present.
- **[PLAN] Phase:** Before writing any code, write your brief technical approach next to the `[PLAN]` checkbox in the REQ file.
- **[APPLY] Phase:** Stay strictly focused on the planned scope. Resist the urge to refactor unrelated code or fix adjacent issues. (Note: You are required to edit this REQ file to update your state checkboxes).
- **[UNIFY] Phase:** Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked — the orchestrator will audit this in Step 6.3.
- **Decisions:** When you make a significant implementation decision not covered by the plan or requirements (e.g., choosing between two valid approaches, deciding on an API shape, selecting a library), log it as a numbered decision. Append a `## Decisions` section to the REQ file:
  ```
  ## Decisions
  - **D-01**: [from Open Questions] Sidebar supports dark mode — consistent UX
  - **D-02**: Used zustand over jotai for state — matches existing project pattern
  - **D-03**: API returns paginated results (20/page) — no explicit requirement, follows existing endpoints
  ```
  Continue numbering from any D-XX decisions already recorded in Open Questions (Step 3.5). Future REQs can reference these: "per D-02 in REQ-003, we use zustand."
- **Out-of-Scope Discoveries:** If you discover unrelated bugs, technical debt, or missing prerequisites, do not fix them inline. Instead, append a `## Discovered Tasks` section to your summary and list them as bullet points so the orchestrator can queue them for later.
```

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
5. **Wired:** For each `(new)` source file, verify it is imported or referenced by at least one other file in the project (grep for the filename or an exported symbol). A new component/module that nothing imports is dead code — flag it. **Exceptions:** Entry points (e.g., `main.ts`, `index.html`), config files, test files, and standalone scripts don't need to be imported by other files.
6. **Flowing:** For files that handle data (API endpoints, data stores, handlers, services), verify the data path isn't hardcoded or stubbed. Check for: hardcoded empty arrays `return []`, placeholder strings like `"TODO"` or `"placeholder"`, `return null` in data-fetching functions, commented-out database calls. If found, flag as hollow implementation — the file exists and is wired but doesn't actually do anything.

**Anti-rationalization rules** (apply when evaluating the above):

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The summary says files changed" | Check the file system | The summary is a claim, not evidence |
| "Tests pass so requirements are met" | Compare requirements to diff, word by word | Tests can be incomplete |
| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |
| "This is probably fine" | Verify it specifically | Probably ≠ verified |

**If qualification fails on any check:**
1. Append a `## Qualification` section to the REQ noting what failed and why.
2. Return to Step 6 — spawn the builder again with the specific failures as context.
3. Maximum **2 re-qualification attempts**. After that, note remaining issues and proceed to Testing (Step 6.5). The review step will catch what remains.

**If qualification passes:**
- Append a brief `## Qualification` section: "Passed — [N] files verified, [N] requirements traced, P-A-U confirmed."
- Proceed to Step 6.5.

### Step 6.5: Testing

Before marking complete, verify tests pass:

1. **Check the prime file for test guidance** — if the REQ's `prime_files` reference a prime with a testing section (test commands, code-area-to-test mappings), use that as the primary source for what to run. Prime test maps are project-specific knowledge that generic detection can't replicate (e.g., "changes to `lib/inpainting.js` require `npm run test:api`" or "`npm test` is always safe but `npm run test:e2e` costs money").
2. **Fall back to generic detection for unmapped files** — if the prime has no testing section, or if you changed files the prime's test map doesn't cover, fall back to generic detection for those files: look for `package.json` test scripts, `jest.config.*`, `pytest.ini`, `Cargo.toml`, `*_test.go`, etc. A partial prime map is not an excuse to skip tests — matched files use the prime's commands, unmatched files use generic detection. If neither source yields test commands for a file, skip testing for it and note it.
3. **Run relevant tests** — target tests related to changed code, not the full suite (unless it's fast). If the prime specifies different commands for different code areas, run only the commands relevant to the files you changed. For unmapped files, run whatever generic detection found.
4. **If tests fail** — return to implementation to fix. On attempt 2+, load `agent-rules/rules-debugging.md` for the builder to follow the structured debugging methodology. Loop until passing or mark as failed after 3 attempts.
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

**TDD verification:** If the REQ has `tdd: true`, the `Red-green validation` section is mandatory — the builder must show test-first evidence (test written before implementation, failed, then passed after). If this evidence is missing, flag it as a qualification issue.

### Step 7: Review

Run the [review work action](./review-work.md) in **pipeline mode** against this REQ.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?).

**How to run it:** Spawn an agent with the review work action file, the REQ path, and the `rules-[domain].md` file (if the domain has one and the file exists). Or read `actions/review-work.md` and follow its pipeline mode instructions in the current session.

**What happens next depends on the review result:**

- **Acceptance = Pass AND overall ≥ 75%**: Append the Review section to the REQ and continue to archive as `completed`. Minor findings go in the report only.
- **Acceptance = Partial OR overall 50-74%**: Append Review, continue to archive as `completed`, but the review **MUST** create follow-up REQs for every Important finding. These are not optional — they enter the queue and block the UR from being considered "done" until addressed.
- **Acceptance = Fail OR overall < 50%**: **Do NOT archive as completed.** Instead:
  1. Append the Review section to the REQ.
  2. Return to Step 6 (Implementation) with the review findings as context for the builder. Load `agent-rules/rules-debugging.md` for the remediation attempt — the builder needs structured debugging methodology, not just "try again."
  3. The builder gets **ONE remediation attempt**.
  4. Re-run Steps 6.25 through 7 (Summary → Qualification → Testing → Review) on the remediated code.
  5. If still failing after remediation: update frontmatter to `status: completed-with-issues`, `completed_at: <timestamp>`, append a `## Remediation` section documenting both attempts, and create follow-up REQs for all remaining Important findings. Then proceed to archive (Step 8) — the frontmatter is already set, so Step 8 should not overwrite it.

The status `completed-with-issues` means the REQ was archived but has known unresolved problems. It counts toward UR completion for archiving purposes, but the follow-up REQs must be processed before the work is considered ship-ready. This status is visible in the recap and present-work actions.

**Follow-up REQs are created based on finding severity, not score.** The review creates follow-up REQs for each **Important** finding (regardless of overall score). Minor and Nit findings go in the report only. The follow-up REQs enter the queue and get processed in a future loop iteration. Follow-up REQs created by the review step must include: `status: pending`, `user_request: [same UR as the reviewed REQ]`, `addendum_to: [reviewed REQ id]`, `domain: [same domain]`, and `review_generated: true`. Place them in `do-work/` root. Cycle detection (Step 8, substep 5) applies to these follow-ups — check the `addendum_to` chain before creating.

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

**Update prime files:** After writing the Lessons Learned section, check the REQ's `prime_files` frontmatter. For each listed prime file, append a link to the lesson under a `## Lessons` section in that prime file (create the section if it doesn't exist):

```markdown
## Lessons

- [REQ-NNN: 1-line summary of the lesson](<relative-path-to-req>#lessons-learned)
```

**Path must be relative to the prime file's location**, not the repo root. Compute the correct relative path from the prime file's directory to the archived REQ file. For example, if the prime file is at `src/utils/prime-auth.md` and the REQ is at `do-work/archive/UR-005/REQ-042-auth-fix.md`, the link should use `../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`.

Only add a link when the lesson is relevant to that prime file's scope — don't spray every lesson into every prime file. If the REQ has no `prime_files` or the lessons aren't relevant to any prime file, skip this.

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
   These go in `do-work/` with `status: pending-answers`. The user reviews them via `do work clarify`.
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
   This ensures non-critical discoveries require the user's explicit permission via `do work clarify` before execution.
5. **Cycle detection:** Before creating any follow-up REQ, check the `addendum_to` chain. If the proposed follow-up would reference a REQ that itself has an `addendum_to` pointing back to the current REQ (or any ancestor in the chain), this is a circular reference. Do not create the follow-up — instead, report the cycle to the user: `⚠ Cycle detected: REQ-NNN → REQ-MMM → REQ-NNN. Skipping follow-up — manual resolution needed.`
6. Archive based on REQ type:

| REQ has... | Archive behavior |
|------------|-----------------|
| `user_request: UR-NNN` | Check if ALL REQs in the UR are finished (status: `completed`, `completed-with-issues`, or `failed`). Check `do-work/` root, `do-work/working/`, `do-work/archive/` root, and `do-work/archive/UR-NNN/` for REQs belonging to this UR. If all finished: move completed/completed-with-issues REQs into UR folder (failed REQs stay at archive root), move entire UR folder to `archive/`. If any REQ is still `pending`, `pending-answers`, or `claimed`: move this REQ to `archive/` root; UR stays in `user-requests/` until last REQ finishes. |
| `context_ref` (legacy) | Move REQ to `archive/`. If all related REQs are now archived, move the CONTEXT doc too. |
| Neither (standalone legacy) | Move directly to `archive/`. |

**On failure:**

Before archiving, classify the failure to determine the correct recovery path. Read the error description and any test/build output, then classify:

| Type | Symptoms | Recovery |
|------|----------|----------|
| **Intent** | Requirements are ambiguous or contradictory; builder couldn't determine what to build | Create a follow-up REQ with `status: pending-answers` containing the specific ambiguities as Open Questions. Archive original as `failed` with `error_type: intent`. |
| **Spec** | Requirements are clear but the technical approach was wrong (wrong files, wrong pattern, wrong architecture) | Create a follow-up REQ with a `## Prior Attempt` section summarizing what was tried and why it failed. Set `status: pending`. Archive original with `error_type: spec`. |
| **Code** | Approach was right but implementation has bugs (tests fail, runtime errors, logic errors) | Create a follow-up REQ targeting the specific code issue. Set `status: pending`. Archive original with `error_type: code`. |
| **Environment** | External dependency unavailable, permissions issue, tooling broken | No follow-up REQ — user must fix the environment. Archive with `error_type: environment` and a clear description of what's needed. |

**Procedure:**
1. Classify using the table above.
2. Update frontmatter: `status: failed`, `error: "description"`, `error_type: [intent|spec|code|environment]`
3. For Intent/Spec/Code failures: create the appropriate follow-up REQ (details above). Set `addendum_to` to the failed REQ's ID so context chains.
4. Move to `archive/` (failed REQs always go to archive root, not into UR folders).
5. Report to user: `[REQ-NNN] failed ([type]): [description]. Follow-up: [REQ-NNN] / None.`

### Step 9: Commit (Git repos only)

Check for git with `git rev-parse --git-dir 2>/dev/null`. If not a git repo, skip.

```bash
# Stage implementation files + archived REQ
git add src/stores/theme-store.ts src/components/settings/SettingsPanel.tsx \
  do-work/archive/UR-002/REQ-003-dark-mode.md

# Stage follow-up REQs created in Step 8 (if any)
git add do-work/REQ-025-confirm-sidebar-palette.md

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

One commit per request. Stage all files created, modified, moved, or deleted during this request's lifecycle: implementation files (listed in the Implementation Summary), the archived REQ file, any follow-up REQs created in Step 8 (`pending-answers` files in `do-work/`), and any UR-folder moves to `archive/`. Do not use `git add -A` or `git add .` — these risk staging secrets, `.env` files, or unrelated changes. Don't bypass pre-commit hooks — fix issues and retry. Failed requests get committed too.

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

Re-check `do-work/` for `REQ-*.md` files (fresh check, not cached).

- **`pending` REQs found**: **CONTEXT WIPE** (see below). Then loop to Step 1.
- **Only `pending-answers` REQs remain**: Write a **Session Checkpoint** (see below), run the [cleanup action](./cleanup.md), then report final summary including a list of the `pending-answers` REQs and their unresolved questions so the user can run `do work clarify` when ready.
- **No REQs at all**: Write a Session Checkpoint, run cleanup, report final summary and exit.

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
queue_state: [N pending, N pending-answers, N in-progress]
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

When invoked with `do work clarify` (or `answers`, `questions`, `pending`, `what's blocked`), the work action enters **clarify mode** instead of the normal work loop. This is the batch-review drain for `pending-answers` REQs.

### Clarify Workflow

1. **Scan the queue**: Find all `REQ-*.md` files in `do-work/` with `status: pending-answers`
2. **If none found**: Report "No pending questions — queue is clear" and exit
3. **Present questions**: For each `pending-answers` REQ, show:
   ```
   REQ-025 — Review fix: dark mode sidebar
   (follow-up to REQ-003, from review)

   1. [ ] Should the sidebar use the same dark palette as the main content?
      Recommended: Yes, match main content palette
      Also: Separate sidebar palette, User-configurable

   2. [ ] Should dark mode persist across sessions?
      Recommended: Yes, save to localStorage
      Also: Reset on refresh, Follow OS preference
   ```
4. **Collect answers**: If your environment has a structured question prompt (multi-question UI), batch questions in groups of **at most 4 per prompt** — chunk by question count, not by REQ. A REQ with 6 questions needs 2 prompts. For each question, the user can:
   - **Answer it** → update to `- [x] [question] → [user's answer]`
   - **Confirm builder's choice** → update to `- [x] [question] → Confirmed: [builder's choice]`. Then check the REQ type:
     - *Discovered-task REQ* (has a "Should I process this as a new task?" question with recommended "Yes, add to queue"): flip `status` to `pending` so the task enters the work queue — see "Approved Discovered Task" below
     - *All other REQs* (builder-decision follow-ups): mark `status: completed` (no implementation needed — see "Builder Was Right" below)
   - **Pick a different option** → update to `- [x] [question] → [user's chosen option]`
   - **Skip for now** → leave as `- [ ]`, REQ stays `pending-answers`
   - **Discard it** → update to `- [x] [question] → Discarded`, then mark the REQ `status: completed`, `completed_at: <timestamp>`, and archive it directly (same pattern as "Builder Was Right" — no implementation work)
5. **Activate answered REQs**: For each REQ that wasn't already completed or discarded: if all questions are now `[x]` or `[~]`, flip `status` from `pending-answers` to `pending`. These enter the queue for the next `do work run`.
6. **Report**: Summary of what was resolved and what's still pending

### Builder Was Right / Discarded

When the user reviews a `pending-answers` follow-up and confirms that the builder's original choice was correct (i.e., no implementation change needed):

1. Update the question to `- [x] [question] → Confirmed: [builder's choice]`
2. Update frontmatter: `status: completed`, `completed_at: <timestamp>`
3. Archive the follow-up REQ directly (skip the work loop — there's nothing to build)
4. Append a brief note: `## Implementation\n\n**No changes needed.** User confirmed builder's choice from [original REQ].\n\n*Resolved via clarify questions*`

**Discarded discovered tasks:** When the user reviews a discovered-task follow-up and chooses "No, discard it", the same fast-path applies. Mark `status: completed`, archive directly, and append: `## Implementation\n\n**Discarded.** User chose not to process this discovered task from [original REQ].\n\n*Resolved via clarify questions*`

### Approved Discovered Task

When the user reviews a discovered-task follow-up (one whose question is "Should I process this as a new task?" with recommended "Yes, add to queue") and confirms the recommendation:

1. Update the question to `- [x] [question] → Confirmed: Yes, add to queue`
2. Update frontmatter: `status: pending` (NOT `completed` — this task needs to be built)
3. **Do not archive.** The REQ stays in `do-work/` and enters the normal work queue for the next `do work run`

This is distinct from "Builder Was Right" because confirming a discovered task means the user wants it *executed*, not signed off. The task has no prior implementation to confirm — it's a new piece of work that needs a full work cycle.

This avoids wasting a work cycle on a REQ that just needs sign-off or rejection, while correctly routing approved discovered tasks into the build queue.

## Progress Reporting

Keep the user informed:

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

## Error Handling

| Phase | Action |
|-------|--------|
| `pending-answers` REQs remain after queue is empty | Report them to the user: list each REQ and its unresolved questions. Suggest `do work clarify` to batch-review. |
| Plan agent fails (Route C) | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Explore agent fails (B/C) | Proceed to implementation with reduced context — builder can explore on its own |
| Implementation fails | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Tests fail repeatedly | After 3 fix attempts, classify as Code failure, create follow-up REQ with test failure details, archive as failed |
| Review: Acceptance = Fail | Return to Step 6 for ONE remediation attempt, then re-review. If still failing: archive as `completed-with-issues` with follow-up REQs |
| Review work agent fails | Skip review, note it in the REQ file, continue to archive — review failure is not a gate |
| Commit fails | Investigate the error (usually a pre-commit hook failure). Fix the underlying issue, re-stage, and retry as a **new** commit. Do NOT use `--no-verify` to skip hooks or `--no-gpg-sign` to bypass signing — fix the root cause. If unfixable, report the error to the user and continue to next request — changes remain uncommitted but archived. |
| Unrecoverable error | Stop loop, report clearly, leave queue intact for manual recovery |

## What This Action Does NOT Do

- Create new request files (use the capture requests action)
- Make architectural decisions beyond what's in the request
- Run without user present (this is supervised automation)
- Modify already-completed requests
- Allow external modification of files in `working/` or `archive/`

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

**Common mistakes to avoid:**
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

## Archived Request File Example

See [sample-archived-req.md](./sample-archived-req.md) for a complete example of what an archived REQ looks like after processing through the full pipeline (Route B). Every section shown there is generated by the steps above.

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = implementation time. Route + timestamps let you calibrate triage accuracy over time.
