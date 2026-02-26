# Do Action

> **Part of the do-work skill.** Invoked when routing determines the user is adding a request. Creates a `do-work/` folder in your project for request tracking.

A fast-capture system for turning ideas into structured request files. Speed over perfection — minimal interaction when intent is clear.

## Core Rules

Every invocation produces exactly two things, always paired:

1. **A UR folder** at `do-work/user-requests/UR-NNN/` with `input.md` containing the full verbatim input
2. **One or more REQ files** at `do-work/REQ-NNN-slug.md`, each linked via `user_request: UR-NNN` in frontmatter

Never create one without the other. A REQ without `user_request` is orphaned. A UR without REQs is pointless. The verify action depends on this linkage.

**Principles:**
- Represent, don't expand — if the user says 5 words, write a 5-word request (with structure)
- The building agent solves technical questions — you're capturing intent, not making architectural decisions
- Never be lossy — for complex inputs, preserve ALL detail in the UR's verbatim section
- After capture, **STOP** — do not start processing the queue or transition into the work action unless the user explicitly asked for both (e.g., "add this and start working")

### First-Run Bootstrap

If `do-work/` doesn't exist yet (first invocation in a project):

1. Create `do-work/` and `do-work/user-requests/`
2. Do **not** pre-create `working/` or `archive/` — those are created by the work action on demand
3. Start numbering at REQ-001 and UR-001

## Simple vs Complex

| Mode | When | Approach |
|------|------|----------|
| **Simple** | Short input (<200 words), 1-2 features, no detailed constraints | Lean format, minimal UR |
| **Complex** | >500 words, 3+ features, detailed requirements/constraints/edge cases, dependencies between features, or user says "spec"/"PRD"/"requirements" | Full preservation with detailed REQ sections |

**When uncertain, treat as complex.** Over-preserving is better than losing requirements.

## File Locations

- `do-work/` root — ONLY for pending `REQ-*.md` files (the queue)
- `do-work/user-requests/UR-NNN/` — verbatim input (`input.md`) and assets (`assets/`)
- **NEVER write to** `do-work/working/` or `do-work/archive/` — those belong to the work action

### Immutability Rule

Files in `working/` and `archive/` are **immutable**. If someone wants to add to an in-flight or completed request, create a new addendum REQ that references the original via `addendum_to` in frontmatter. The addendum goes through the queue like any other request.

## File Naming

- **REQ files:** `REQ-[number]-[slug].md` in `do-work/` root
- **UR folders:** `do-work/user-requests/UR-[number]/` containing `input.md` and optional `assets/`
- **Assets:** `do-work/user-requests/UR-NNN/assets/REQ-[num]-[descriptive-name].png`

To get the next REQ number, check existing `REQ-*.md` files across `do-work/`, `do-work/working/`, and `do-work/archive/` (including inside `do-work/archive/UR-*/`), then increment from the highest. For the next UR number, check `do-work/user-requests/UR-*/` and `do-work/archive/UR-*/`. REQ and UR use separate numbering sequences. If no existing files are found anywhere, start at 1.

### Backward Compatibility

Legacy REQ files (pre-UR system) may lack `user_request` and reference `CONTEXT-*.md` files or `do-work/assets/` instead. This is fine — the work action handles both patterns. New REQs always get `user_request`.

## Request File Formats

### Simple REQ

```markdown
---
id: REQ-001
title: Brief descriptive title
status: pending
created_at: 2025-01-26T10:00:00Z
user_request: UR-001
---

# [Brief Title]

## What
[1-3 sentences describing what is being requested]

## Why (if provided)
[User's stated reasoning — omit if not provided]

## Context
[Additional context, constraints, or details mentioned]

## Assets
[Description of screenshots or links to saved files]

---
*Source: [original verbatim request]*
```

### Complex REQ (additional sections)

Complex requests use the same base format plus these sections:

```markdown
## Detailed Requirements
[Extract EVERY requirement from the original input that applies to THIS feature.
DO NOT SUMMARIZE — use the user's words. Include specific values, constraints,
conditions, edge cases, "must"/"should"/"never" statements.]

## Constraints
[Limitations, restrictions, batch-level concerns that apply to this REQ]

## Dependencies
[What this feature needs or what needs it — reference other REQ IDs]

## Builder Guidance
[Certainty level: Exploratory / Firm / Mixed. Scope cues like "keep it simple."
Any latitude given to the builder.]

## Open Questions
- [ ] [Question about ambiguity the user needs to clarify]
  Recommended: [best default based on context]
  Also: [alternative A], [alternative B]

Open Questions use checkbox syntax with recommended choices. Each question includes a `Recommended:` line (the best default if the user doesn't answer) and an `Also:` line with alternatives. The choices make questions answerable even when the question itself isn't fully understood — the user can just pick one.

`- [ ]` = unresolved, `- [x]` = answered (answer follows `→`), `- [~]` = deferred to builder (note follows `→`).

**Capture time is the optimal window for resolving these.** During capture (this action), use your environment's ask-user prompt/tool to present Open Questions immediately. The user is here, engaged, and fleshing out the request — don't defer what you can clarify now. Only leave questions as `- [ ]` if you genuinely can't ask (e.g., batch processing, async capture).

Only add questions where the user's intent is genuinely unclear — don't add questions the builder can answer by reading the codebase.

## Full Context
See [user-requests/UR-NNN/input.md](./user-requests/UR-NNN/input.md) for complete verbatim input.
```

**Additional frontmatter for complex requests:**
- `related: [REQ-006, REQ-007]` — other REQs in this batch
- `batch: auth-system` — batch name grouping related requests
- `addendum_to: REQ-005` — if this amends an in-flight/completed request

### UR input.md

Created for every invocation. For simple requests, it's minimal:

```markdown
---
id: UR-005
title: Add keyboard shortcuts
created_at: 2025-01-26T10:00:00Z
requests: [REQ-020]
word_count: 4
---

# Add keyboard shortcuts

## Full Verbatim Input

add keyboard shortcuts

---
*Captured: 2025-01-26T10:00:00Z*
```

For complex requests, add a Summary, an Extracted Requests table, and a Batch Constraints section before the Full Verbatim Input. The verbatim section must contain the COMPLETE, UNEDITED input — never summarize or clean it up.

## Workflow

### Step 1: Parse and Assess

Read the user's input. Determine:
- **Single vs multiple requests** — look for "and also", comma-separated lists, numbered items, distinct topics
- **Simple vs complex** — apply the detection criteria above

### Step 2: Check for Duplicates

**Queued requests** — read each `REQ-*.md` in `do-work/` and compare the new request's intent against the existing file's `title`, heading, and `## What` section. Slugs are lossy — a file named `REQ-042-ui-cleanup.md` may contain the exact requirement being re-submitted under different phrasing. Match on intent, not just keywords.

**In-flight and archived requests** — list filenames in `do-work/working/` and `do-work/archive/` (including inside `do-work/archive/UR-*/`). A filename scan is sufficient here since these files are immutable regardless.

If `do-work/` is freshly bootstrapped (no existing REQ files anywhere), skip duplicate checking entirely.

For each parsed request, check for similar existing ones across both tiers.

| Existing request is in... | Action |
|---------------------------|--------|
| `do-work/` (queue) | If same: tell user, skip. If similar: ask. If enhancement: append an Addendum section to the pending file |
| `do-work/working/` | **NEVER modify.** Create a new addendum REQ with `addendum_to` field |
| `do-work/archive/` | **NEVER modify.** Create a new addendum REQ with `addendum_to` field |

**Addendum to a queued request** — don't rewrite, append:

```markdown
## Addendum (2025-01-27)

User added: "dark mode should also affect the sidebar"

- Sidebar must also respect dark mode theme
```

**Addendum for in-flight/completed requests** — create a new REQ:

```markdown
---
id: REQ-021
title: "Addendum: dark mode sidebar support"
status: pending
created_at: 2025-01-27T09:00:00Z
user_request: UR-006
addendum_to: REQ-005
---

# Addendum: Dark Mode Sidebar Support

## What
Add sidebar support to the existing dark mode implementation (REQ-005).

## Context
Addendum to REQ-005, which is currently [in progress / completed].
The user wants the sidebar to also support dark mode.

## Requirements
- Sidebar must respect the dark mode theme
```

### Step 3: Clarify Only If Needed

Ask questions ONLY when the request is genuinely ambiguous (could mean two very different things), or when a duplicate/similar request makes intent unclear. Don't ask about implementation details — that's for the building agent.

### Step 4: Handle Screenshots

If the user provides a screenshot:
1. Find the image using your platform's attachment mechanism or cache
2. Copy to `do-work/user-requests/UR-NNN/assets/REQ-[num]-[slug].png`
3. Reference in the REQ's Assets section
4. Write a thorough text description (what it shows, visible text, layout, problems visible) — this is the primary record for searchability

### Step 5: Write Files

Before writing, ensure `do-work/` and `do-work/user-requests/UR-NNN/` exist (create if needed).

**For all requests (simple and complex):**
1. Create `do-work/user-requests/UR-NNN/input.md` with verbatim input (leave `requests` array empty initially)
2. Create `REQ-NNN-slug.md` files using the appropriate format, with `user_request: UR-NNN`
3. Update the UR's `requests` array with all created REQ IDs

**Complex mode additionally:**
- Create `assets/` subfolder in the UR folder
- Extract EVERY requirement into the appropriate REQ — do not summarize
- Set `related` and `batch` fields across the batch
- Add Batch Constraints to the UR (cross-cutting concerns, scope cues, sequencing)
- Duplicate batch-level constraints into each relevant REQ's Constraints section
- Re-read the original input to verify nothing was dropped — especially UX/interaction details and intent signals (certainty level, scope cues)

### Step 6: Report Back

Brief summary of created files. If the request was meaningfully complex (complex mode, 3+ REQs, or notably long/nuanced input), add:

> That was a pretty detailed request — it's possible the capture missed some nuances. You can run `/do-work verify` to check coverage against your original input.

## Examples

### Simple Capture

```
User: do work add keyboard shortcuts

Created:
- do-work/user-requests/UR-003/input.md
- do-work/REQ-004-keyboard-shortcuts.md
```

### Multiple Requests

```
User: do work add dark mode, also the search feels slow, and we need an export button

Created:
- do-work/user-requests/UR-004/input.md
- do-work/REQ-005-dark-mode.md
- do-work/REQ-006-search-performance.md
- do-work/REQ-007-export-button.md
```

### Addendum to In-Flight Request

```
User: do work dark mode should also affect the sidebar

[Checks existing — REQ-005-dark-mode.md is in do-work/working/]

REQ-005 is currently being worked on — creating a follow-up request instead.

Created:
- do-work/user-requests/UR-006/input.md
- do-work/REQ-021-addendum-dark-mode-sidebar.md (addendum_to: REQ-005)
```

### Complex Multi-Feature Request

```
User: do work [detailed auth system requirements — OAuth, profiles, sessions, password reset...]

Created:
- do-work/user-requests/UR-001/input.md (full verbatim input, 1847 words)
- do-work/REQ-010-oauth-login.md (user_request: UR-001)
- do-work/REQ-011-user-profiles.md (user_request: UR-001)
- do-work/REQ-012-session-management.md (user_request: UR-001)
- do-work/REQ-013-password-reset.md (user_request: UR-001)

That was a pretty detailed request — it's possible the capture missed some
nuances. You can run `/do-work verify` to check coverage against your original input.
```

## Edge Cases

- **Vague request ("fix the search")**: Capture what was said. The builder can clarify.
- **References earlier conversation**: Include that context in the request file.
- **Seems impossible or contradictory**: Capture it. Add contradictions as `- [ ]` Open Questions with recommended resolutions — and ask the user right now if they're available.
- **Requirement applies to multiple features**: Include in ALL relevant REQ files. Duplication beats losing it.
- **User changes mind mid-request**: Capture the final decision, note the evolution in the UR.
- **Mentioned once in passing**: Still a requirement. Capture it.
