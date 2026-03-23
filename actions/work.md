# Work Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to process the queue. Processes requests from the `do-work/` folder in your project.

An orchestrated build system that processes request files created by the capture requests action. Uses complexity triage to route simple requests straight to implementation and complex ones through planning and exploration first.

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary (mandatory file manifest), Testing, Review. This ensures full traceability — what was planned vs done, what files were touched, and whether triage was accurate.

## Architecture

```
work action (orchestrator - lightweight, stays in loop)
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
  │     │     │   Explore only, then build│          │
  │     │     │                           ▼          │
  │     │     └── Route C (Complex) ──► Plan ──► Explore
  │     │                                            │
  │     │                                            ▼
  │     │                                     Implementation agent
  │     │                                            │
  │     │                                            ▼
  │     │                                  Implementation Summary
  │     │                                            │
  │     │                                            ▼
  │     │                                        Testing
  │     │                                            │
  │     │                                            ▼
  │     │                                        Review
  │     │                                            │
  │     │                                            ▼
  │     ├── Archive ──► create pending-answers follow-ups for - [~] items
  │     │                                            │
  │     │                                            ▼
  │     └── Commit (git repos only)
  │
  └── Loop until queue empty → cleanup → report (tip: `do work clarify` for pending-answers)
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
status: completed | failed
commit: abc1234               # If git repo
error: "Description"          # Only if failed
---
```

**Status flow:** `pending` → `claimed` → `[planning]` → `[exploring]` → `implementing` → `testing` → `reviewing` → `completed` / `failed`

**Special status:** `pending-answers` — a follow-up REQ whose Open Questions need user input before it can be worked. These accumulate in the queue and get batch-reviewed when the user runs `do work clarify`.

## Workflow

**The work action is an orchestrator.** You handle ALL file management (moving files, updating frontmatter, appending sections, archiving). Spawned agents handle implementation work only.

### Step 1: Find Next Request

**Crash Recovery:** Before checking the queue, look inside `do-work/working/` for any `REQ-*.md` files. If any exist, a previous run was interrupted. For each recovered REQ:
1. Reset frontmatter: set `status` to `pending`, remove `claimed_at` and `route`
2. Strip sections generated during the interrupted run: remove `## Triage`, `## Exploration`, `## Plan`, `## Implementation Summary`, and `## Testing` sections (and their content) if present — these may be incomplete or stale from the crash. Leave `## Open Questions` and user-authored content intact.
3. Move the REQ back to the `do-work/` root

Once `working/` is empty, proceed with finding the next request.

Glob for `do-work/REQ-*.md` (root of `do-work/`, **not** a subdirectory — there is no `queue/` folder). Sort by number. Read the frontmatter of each (in number order) to check `status` — pick the first with `status: pending` (skip `pending-answers` — those wait for user input). Don't read the full body at this stage. If no `pending` REQs found, report completion and exit. If only `pending-answers` REQs remain, report them to the user so they can batch-review the questions.

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
2. Mark each as `- [~]` with the builder's chosen approach and reasoning: `- [~] [question] → Builder chose: [choice]. Reasoning: [why]`
3. Proceed with implementation using those decisions.

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

If an `## Exploration` section does not already exist, append the output:

```markdown
## Exploration

[Explore agent findings — key files, patterns, concerns]

*Generated by Explore agent*
```

### Step 6: Implementation

Spawn a **general-purpose agent** with the `rules-[domain].md` file (if domain is missing or the file doesn't exist, skip loading it), any files listed in the `prime_files` array, and context appropriate to the route:

- **Route A**: Request content only — "triaged as simple, aim for a focused minimal change"
- **Route B**: Request + exploration output — "follow existing patterns identified above"
- **Route C**: Request + plan + exploration output — "implement according to the plan"

All routes include these instructions to the agent:

```
- **Prime Files:** If `prime_files` are attached to this REQ, READ THEM FIRST. They are your map to the codebase. If NO prime file exists for the primary utility you are modifying, you MUST investigate the utility and create one (`prime-[name].md`). Prime files must be low-noise, high-value, point to code as the source of truth, and avoid volatile metrics (like test counts). If you create one, update the REQ's frontmatter to include it.
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
- **[PLAN] Phase:** Before writing any code, write your brief technical approach next to the `[PLAN]` checkbox in the REQ file.
- **[APPLY] Phase:** Stay strictly focused on the planned scope. Resist the urge to refactor unrelated code or fix adjacent issues. (Note: You are required to edit this REQ file to update your state checkboxes).
- **[UNIFY] Phase:** Run native project linters and manually review your own diff to ensure no debug artifacts (e.g., console.log, TODOs) are left behind before checking the `[UNIFY]` box. Do not rely on external bash scripts for this.
- **Out-of-Scope Discoveries:** If you discover unrelated bugs, technical debt, or missing prerequisites, do not fix them inline. Instead, append a `## Discovered Tasks` section to your summary and list them as bullet points so the orchestrator can queue them for later.
```

### Step 6.25: Implementation Summary

After implementation completes, append a manifest of what changed to the request file. This is the primary auditability artifact — without it, there's no way to verify the REQ was implemented without digging through git history.

Append to the request file:

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

### Step 6.5: Testing

Before marking complete, verify tests pass:

1. **Check the prime file for test guidance** — if the REQ's `prime_files` reference a prime with a testing section (test commands, code-area-to-test mappings), use that as the primary source for what to run. Prime test maps are project-specific knowledge that generic detection can't replicate (e.g., "changes to `lib/inpainting.js` require `npm run test:api`" or "`npm test` is always safe but `npm run test:e2e` costs money").
2. **Fall back to generic detection for unmapped files** — if the prime has no testing section, or if you changed files the prime's test map doesn't cover, fall back to generic detection for those files: look for `package.json` test scripts, `jest.config.*`, `pytest.ini`, `Cargo.toml`, `*_test.go`, etc. A partial prime map is not an excuse to skip tests — matched files use the prime's commands, unmatched files use generic detection. If neither source yields test commands for a file, skip testing for it and note it.
3. **Run relevant tests** — target tests related to changed code, not the full suite (unless it's fast). If the prime specifies different commands for different code areas, run only the commands relevant to the files you changed. For unmapped files, run whatever generic detection found.
4. **If tests fail** — return to implementation to fix. Loop until passing or mark as failed after 3 attempts.
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

### Step 7: Review

Run the [review work action](./review-work.md) in **pipeline mode** against this REQ.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?).

**How to run it:** Spawn an agent with the review work action file and the REQ path, or read `actions/review-work.md` and follow its pipeline mode instructions in the current session.

**What happens next depends on the review score:**

- **75%+ overall**: Append the Review section to the REQ and continue to archive. Minor findings go in the report only.
- **Below 75%**: Same as above — append Review, continue to archive. The current REQ is still marked completed.

**Follow-up REQs are created based on finding severity, not score.** The review creates follow-up REQs for each **Important** finding (regardless of overall score). Minor and Nit findings go in the report only. The follow-up REQs enter the queue and get processed in a future loop iteration.

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
- **Required for Routes B and C** — there's always something worth recording when exploration or planning was involved. **Optional for Route A** — skip if everything went smoothly with no surprises.
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
4. **Queue Discovered Tasks:** Check the REQ file for a `## Discovered Tasks` section (appended by the implementation agent as a separate section — not inside `## Implementation Summary`). For every item listed, create a new follow-up REQ file in the `do-work/` root.
   - Set frontmatter: `status: pending-answers`, `user_request: [same UR]`, `addendum_to: [current REQ id]`, `domain: [same domain as current REQ]`.
   - Add an `## Open Questions` section with this checkbox format:
     `- [ ] I discovered this out-of-scope task while working on [current REQ]: [Task Description]. Should I process this as a new task?`
     `  Recommended: Yes, add to queue (will flip to 'pending').`
     `  Also: No, discard it.`
   This ensures out-of-scope discoveries are safely captured but require the user's explicit permission via `do work clarify` before execution.
5. Archive based on REQ type:

| REQ has... | Archive behavior |
|------------|-----------------|
| `user_request: UR-NNN` | Check if ALL REQs in the UR are complete. If yes: move completed REQs into UR folder, move entire UR folder to `archive/`. If no: move REQ to `archive/` root; UR stays in `user-requests/` until last REQ completes. |
| `context_ref` (legacy) | Move REQ to `archive/`. If all related REQs are now archived, move the CONTEXT doc too. |
| Neither (standalone legacy) | Move directly to `archive/`. |

**On failure:**

1. Update frontmatter: `status: failed`, `error: "description"`
2. Move to `archive/` (failed REQs always go to archive root, not into UR folders)
3. Report failure to user

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

**Validation check (successful REQs only):** Before committing, compare the `## Implementation Summary` file list against the staged files (excluding `do-work/` paths). If the Implementation Summary lists files that aren't staged, or if the only staged files are `do-work/` metadata, flag the mismatch — the commit may not contain the actual implementation. Fix the staging or update the Implementation Summary before proceeding. **Skip this check for failed REQs** — they may have no Implementation Summary or no project files staged, and that's expected.

**Write commit hash back to the archived REQ.** After the commit succeeds, retrieve the hash with `git rev-parse --short HEAD` and update the archived REQ's frontmatter `commit:` field with the actual value. Then amend the commit to include this update:

```bash
# After the initial commit succeeds:
HASH=$(git rev-parse --short HEAD)
# Update the commit: field in the archived REQ frontmatter
# (replace "commit:" line or add it if missing)
git add do-work/archive/UR-NNN/REQ-NNN-slug.md
git commit --amend --no-edit
```

This ensures the `commit:` field in the archived REQ contains the real hash, which the review-work and present-work actions depend on for traceability. Without this step, the field would be empty or a placeholder.

### Step 10: Loop or Exit

Re-check `do-work/` for `REQ-*.md` files (fresh check, not cached).

- **`pending` REQs found**: **CONTEXT WIPE**. Before looping back to Step 1 to grab the next pending REQ, you MUST clear your sub-agents' memory and close all open editor tabs. Drop all architectural assumptions from the previous REQ. Treat the next REQ as an entirely new, isolated project to prevent context drift. Then loop to Step 1.
- **Only `pending-answers` REQs remain**: Run the [cleanup action](./cleanup.md), then report final summary including a list of the `pending-answers` REQs and their unresolved questions so the user can run `do work clarify` when ready.
- **No REQs at all**: Run cleanup, report final summary and exit.

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
  Exploring...    [done]
  Implementing... [done]
  Summary...      [done] 3 files changed
  Testing...      [done] ✓ 12 tests passing
  Reviewing...    [done] 92% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → abc1234

Processing REQ-004-fix-typo.md...
  Triage: Simple (Route A)
  Implementing... [done]
  Summary...      [done] 1 file changed
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
| Plan agent fails (Route C) | Mark failed, continue to next request |
| Explore agent fails (B/C) | Proceed to implementation with reduced context — builder can explore on its own |
| Implementation fails | Mark failed, preserve plan/exploration outputs for retry |
| Tests fail repeatedly | After 3 fix attempts, mark failed with test failure details |
| Review work agent fails | Skip review, note it in the REQ file, continue to archive — review is advisory, not a gate |
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
□ Step 1: Find next request (pick first pending, skip pending-answers)
□ Step 2: Claim request (mkdir -p working/, move REQ, update status & claimed_at)
□ Step 3: Triage (decide route, append ## Triage, read original if addendum)
□ Step 3.5: Handle Open Questions (mark - [~] with builder's choice)
□ Step 4: Plan (Route C: spawn Plan agent / Routes A & B: note skipped)
□ Step 5: Explore (Routes B & C: spawn Explore agent, append ## Exploration)
□ Step 6: Implement (spawn agent, execute P-A-U state loop, agent updates checkboxes)
□ Step 6.25: Implementation Summary (append file manifest — mandatory for all routes)
□ Step 6.5: Test (run relevant tests, append ## Testing)
□ Step 7: Review (spawn review action in pipeline mode)
□ Step 7.5: Lessons Learned (append section, skip for Route A if no surprises)
□ Step 8: Archive (update status, append summary, queue follow-ups & discovered tasks, move to archive/)
□ Step 9: Commit (stage explicit files, commit if git repo, amend hash to REQ)
□ Step 10: Loop or Exit (CONTEXT WIPE if looping, else cleanup and exit)
```

**Common mistakes to avoid:**
- Spawning implementation agent without first moving file to `working/`
- Letting spawned agents handle file management (only the orchestrator moves/archives files)
- Forgetting to update status in frontmatter at each phase transition
- Archiving a UR folder before all its REQs are complete
- Forgetting Planning status note for Routes A/B ("Planning not required")
- Using `git add -A` instead of staging specific files
- Using `--no-verify` to bypass a failing pre-commit hook instead of fixing the issue
- Committing without validating Implementation Summary file list against staged files
- Implementation Summary that only lists `do-work/` paths (means the REQ wasn't actually implemented)
- Creating follow-ups for every `- [~]` item instead of only UX-affecting decisions

## Archived Request File Example

After processing, each archived REQ contains its full history:

```markdown
---
id: REQ-007
title: Add user avatar component
status: completed
created_at: 2025-01-26T09:30:00Z
claimed_at: 2025-01-26T11:00:00Z
route: B
completed_at: 2025-01-26T11:08:00Z
commit: a1b2c3d
prime_files: []
---

# Add User Avatar Component

## What
[Original request content]

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read `agent-rules/rules-[domain].md`. Write brief technical approach here. Do not write code yet.) -> Analyzed Avatar.tsx and determined we need to build UserAvatar.tsx wrapping it.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Native project linters run, tests passed, and diff hygiene manually verified to remove debug logs.)

---

## Triage
**Route: B** - Medium
**Reasoning:** Clear feature but need to find existing component patterns.
**Planning:** Not required

## Plan
**Planning not required** - Route B: Exploration-guided implementation
*Skipped by work action*

## Exploration
- Found similar component at src/components/Avatar.tsx
- Uses pattern X for state management
*Generated by Explore agent*

## Implementation Summary

**Files changed:**
- `src/components/UserAvatar.tsx` (new)
- `tests/user-avatar.spec.ts` (new)

**What was done:** Created a UserAvatar component wrapping the existing Avatar.tsx with user-specific props and default state handling.

## Testing
**Tests run:** npm test -- --testPathPattern="user-avatar"
**Result:** ✓ All passing (4 tests)
*Verified by work action*

## Review

**Overall: 90%** | 2025-01-26T11:06:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — component renders correctly with all avatar states
**Suggested testing:** 1 item
**Follow-ups created:** None

*Reviewed by review work action*

## Lessons Learned

**What worked:** Reused existing Avatar.tsx patterns — saved time and kept consistency.
**What didn't:** Initially tried CSS modules for scoping, but project uses styled-components everywhere — switched after exploration.
**Worth knowing:** Avatar sizes are constrained by the grid layout in `AppShell.tsx` — don't go above 48px without checking the sidebar.
```

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = implementation time. Route + timestamps let you calibrate triage accuracy over time.
