# Work Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to process the queue. Processes requests from the `do-work/` folder in your project.

An orchestrated build system that processes request files created by the capture requests action. Uses complexity triage to route simple requests straight to implementation and complex ones through planning and exploration first.

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary, Testing, Review. This ensures full traceability — what was planned vs done, where failures happened, and whether triage was accurate.

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
  │     │                                        Testing
  │     │                                            │
  │     │                                            ▼
  │     │                                        Review
  │     │                                            │
  │     │                                            ▼
  │     └── Archive ──► create pending-answers follow-ups for - [~] items
  │
  └── Loop until queue empty → cleanup → report (tip: `do work clarify` for pending-answers)
```

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

List (don't read) `REQ-*.md` filenames in `do-work/`. Sort by number, pick the first with `status: pending` (skip `pending-answers` — those wait for user input). If no `pending` REQs found, report completion and exit. If only `pending-answers` REQs remain, report them to the user so they can batch-review the questions.

### Step 2: Claim the Request

1. `mkdir -p do-work/working` and move the REQ file there
2. Update frontmatter: `status: claimed`, `claimed_at: <timestamp>`

### Step 3: Triage

**If the REQ has an `addendum_to` field:** Before triaging, locate and read the referenced REQ from its current location (check `do-work/working/`, `do-work/archive/`, and `do-work/archive/UR-*/` in that order). Use its context (What, Requirements, Implementation Summary if completed) alongside the addendum REQ when triaging and when passing context to the builder in later steps. This ensures the builder understands what already exists and what the addendum is changing.

Read the request, apply the decision flow, update frontmatter with `route`. Append to the request file:

```markdown
---

## Triage

**Route: [A/B/C]** - [Simple/Medium/Complex]

**Reasoning:** [1-2 sentences]

**Planning:** [Required/Not required]
```

Report the triage decision briefly to the user.

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

**Route C:** Spawn a **Plan agent** with the request content and project context. Ask it to produce a specific implementation plan (files to modify, order of changes, architectural decisions, testing approach). Append the output:

```markdown
## Plan

[Plan agent output]

*Generated by Plan agent*
```

**Routes A and B:** Append a skip note:

```markdown
## Plan

**Planning not required** - [Route A: Direct implementation / Route B: Exploration-guided implementation]

*Skipped by work action*
```

### Step 5: Exploration (Routes B and C)

Spawn an **Explore agent** to find relevant files, existing patterns, types/interfaces, and testing conventions.

- **Route C**: Give it the plan and ask it to find files mentioned in the plan plus similar implementations
- **Route B**: Give it the request and ask it to find where the change should go and what patterns to follow

Append the output:

```markdown
## Exploration

[Explore agent findings — key files, patterns, concerns]

*Generated by Explore agent*
```

### Step 6: Implementation

Spawn a **general-purpose agent** with context appropriate to the route:

- **Route A**: Request content only — "triaged as simple, aim for a focused minimal change"
- **Route B**: Request + exploration output — "follow existing patterns identified above"
- **Route C**: Request + plan + exploration output — "implement according to the plan"

All routes include these instructions to the agent:

```
- You have full access to edit files and run shell commands
- If you find the request is more complex than expected, you can explore or plan as needed
- Document any blockers clearly
- Identify existing tests related to your changes
- Write new tests for new functionality / regression tests for bug fixes
- Update existing tests if behavior intentionally changed
- When complete, summarize: what changed, what tests exist, what new tests were written
```

### Step 6.5: Testing

Before marking complete, verify tests pass:

1. **Detect testing infrastructure** — look for `package.json` test scripts, `jest.config.*`, `pytest.ini`, `Cargo.toml`, `*_test.go`, etc. If none found, skip testing and note it.
2. **Run relevant tests** — target tests related to changed code, not the full suite (unless it's fast)
3. **If tests fail** — return to implementation to fix. Loop until passing or mark as failed after 3 attempts.
4. **If new tests are needed** — spawn a general-purpose agent to write them following existing patterns, then run them.

Append to the request file:

```markdown
## Testing

**Tests run:** [command]
**Result:** ✓ All passing (X tests)

**New tests added:**
- [list]

*Verified by work action*
```

### Step 7: Review

Run the [review work action](./review-work.md) in **pipeline mode** against this REQ.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?).

**How to run it:** Spawn an agent with the review work action file and the REQ path, or read `actions/review-work.md` and follow its pipeline mode instructions in the current session.

**What happens next depends on the review score:**

- **75%+ overall**: Append the Review section to the REQ and continue to archive. Minor findings go in the report only.
- **Below 75%**: Review creates follow-up REQ files in `do-work/` (using the `addendum_to` pattern). Append the Review section to the REQ and continue to archive — the current REQ is still marked completed. The follow-up REQs enter the queue and get processed in a future loop iteration.

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
**Key files:** [Pointers to the most important files — these are the source of truth, not this summary]
**Worth knowing:** [Anything the next person touching this code should know — gotchas, edge cases, non-obvious dependencies]
```

**Rules:**
- Keep it concise — pointers to code, not walls of text. The code is the source of truth.
- Only include entries that have value. If everything went smoothly (Route A, no surprises), skip this section entirely.
- "What didn't work" is the most valuable part — it prevents repeating mistakes.
- Always reference specific files rather than describing their contents.

### Step 8: Archive

**On success:**

1. Update frontmatter: `status: completed`, `completed_at: <timestamp>`
2. Append implementation summary if not already present
3. **Create follow-ups for builder-decided questions:** If the REQ has any `- [~]` items in Open Questions where the builder's choice meaningfully affects the user experience, create a follow-up REQ for each:
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
4. Archive based on REQ type:

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
git add -A
git commit -m "$(cat <<'EOF'
[REQ-003] Dark Mode (Route C)

Implements: do-work/archive/REQ-003-dark-mode.md

- Created src/stores/theme-store.ts
- Modified src/components/settings/SettingsPanel.tsx

EOF
)"
```

**Format:** `[{id}] {title} (Route {route})` + `Implements:` line + summary bullets. Add a co-author trailer if your platform convention calls for one (e.g., `Co-Authored-By: Agent <agent@example.com>`), otherwise omit.

One commit per request. Stage everything with `git add -A`. Don't bypass pre-commit hooks — fix issues and retry. Failed requests get committed too.

### Step 10: Loop or Exit

Re-check `do-work/` for `REQ-*.md` files (fresh check, not cached).

- **`pending` REQs found**: Loop to Step 1.
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
4. **Collect answers**: For each question, the user can:
   - **Answer it** → update to `- [x] [question] → [user's answer]`
   - **Confirm builder's choice** → update to `- [x] [question] → Confirmed: [builder's choice]` and mark the REQ `status: completed` (no implementation needed — see "Builder Was Right" below)
   - **Pick a different option** → update to `- [x] [question] → [user's chosen option]`
   - **Skip for now** → leave as `- [ ]`, REQ stays `pending-answers`
5. **Activate answered REQs**: For each REQ that wasn't already completed by the Builder Was Right path: if all questions are now `[x]` or `[~]`, flip `status` from `pending-answers` to `pending`. These enter the queue for the next `do work run`.
6. **Report**: Summary of what was resolved and what's still pending

### Builder Was Right

When the user reviews a `pending-answers` follow-up and confirms that the builder's original choice was correct (i.e., no implementation change needed):

1. Update the question to `- [x] [question] → Confirmed: [builder's choice]`
2. Update frontmatter: `status: completed`, `completed_at: <timestamp>`
3. Archive the follow-up REQ directly (skip the work loop — there's nothing to build)
4. Append a brief note: `## Implementation\n\n**No changes needed.** User confirmed builder's choice from [original REQ].\n\n*Resolved via clarify questions*`

This avoids wasting a work cycle on a REQ that just needs sign-off.

## Progress Reporting

Keep the user informed:

```
Processing REQ-003-dark-mode.md...
  Triage: Complex (Route C)
  Open Questions: 2 found → builder decided (follow-ups queued)
  Planning...     [done]
  Exploring...    [done]
  Implementing... [done]
  Testing...      [done] ✓ 12 tests passing
  Reviewing...    [done] 92% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → abc1234

Processing REQ-004-fix-typo.md...
  Triage: Simple (Route A)
  Implementing... [done]
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
| Commit fails | Report error, continue to next request — changes remain uncommitted but archived |
| Unrecoverable error | Stop loop, report clearly, leave queue intact for manual recovery |

## What This Action Does NOT Do

- Create new request files (use the capture requests action)
- Make architectural decisions beyond what's in the request
- Run without user present (this is supervised automation)
- Modify already-completed requests
- Allow external modification of files in `working/` or `archive/`

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
---

# Add User Avatar Component

## What
[Original request content]

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
- Created src/components/UserAvatar.tsx
- Added tests in tests/user-avatar.spec.ts
*Completed by work action (Route B)*

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
**Key files:** `src/components/UserAvatar.tsx`, `tests/user-avatar.spec.ts`
**Worth knowing:** Avatar sizes are constrained by the grid layout in `AppShell.tsx` — don't go above 48px without checking the sidebar.
```

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = implementation time. Route + timestamps let you calibrate triage accuracy over time.
