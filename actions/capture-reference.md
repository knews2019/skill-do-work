# Capture Requests — Reference

> **Companion file to `actions/capture.md`.** Holds the REQ/UR templates, the schema-alias table, and the addendum-REQ template — the hard-contract field specifications that `actions/work.md`, `actions/roadmap.md`, and `tools/queue-kanban/model.go` all read. Each section below is pointed to by name from the matching step in `actions/capture.md`. Loading this file is only necessary when you reach the step that references it — and read only the named section. If this file is already in context from an earlier step this session, reuse it; don't re-read it at every reference site.

---

## Request File Formats

### Simple REQ

```markdown
---
id: REQ-001
title: Brief descriptive title
status: pending
created_at: 2025-01-26T10:00:00Z  # current UTC instant — date -u +%Y-%m-%dT%H:%M:%SZ, never local time with a Z suffix (Timestamp rule, actions/work-reference.md)
user_request: UR-001
domain: frontend  # choose one: frontend, backend, ui-design, general, security, or testing
prime_files: []  # list paths to relevant prime-*.md files, or leave empty
tdd: true  # default true when a runnable RED test can be written in this project's harness; false otherwise (see heuristic below)
suggested_spec:  # optional — spec template name if one clearly matches (e.g., "api-endpoint", "bug-fix")
depends_on: []  # optional — list of REQ IDs that must complete before this REQ runs; honored by actions/work.md's selection scan
maintenance: false  # set true ONLY for a deliberate removal/narrowing of the skill's OWN operating instructions — see Step 1's maintenance assessment; loads crew-members/maintenance.md in actions/work.md Step 6
# External-condition fields — set ONLY when the task waits on something outside the queue (see Step 1's external-condition assessment). Omit all three for normal REQs.
# status: blocked            # use INSTEAD OF `status: pending` when the REQ cannot start until an external condition is met — distinct from depends_on (another REQ) and Open Questions (a question for the user)
# blocked_by: "..."          # free text naming the condition (always YAML-quoted), e.g. "LM Studio running locally", "designer answered on mockups"
# blocked_at: 2026-01-26T10:00:00Z   # stamp the moment it was captured blocked — current UTC instant, same Timestamp rule as created_at
# blocked_check: "..."       # OPTIONAL shell probe (YAML-quoted) — emit ONLY when the user supplies or explicitly confirms the command; never invent one
---

# [Brief Title]

## What
[1-3 sentences describing what is being requested]

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)
[User's stated reasoning — omit if not provided]

## Context
[Additional context, constraints, or details mentioned]

## Red-Green Proof
**RED prompt/case:** [Minimal prompt, repro, or example that should fail or be missing today]
**Why RED now:** [What is currently broken or absent]
**GREEN when:** [Observable result that proves the request is done]
**Validation:** [User confirmed / User adjusted / Inferred during capture]

## Assets
[Description of screenshots or links to saved files]

---
*Source: [original verbatim request]*

Think carefully before answering.
```

Include `## Red-Green Proof` when the request is behavior-changing and can be proven with a prompt, repro, or example. If `tdd: true`, this section is mandatory. The goal is proof of behavior, not implementation detail.

Treat defining the RED state as essential, high-value capture work. It is one of the most helpful things you can do for the downstream builder because it turns vague intent into a concrete failing proof target. Do not treat this as paperwork. Lean into it. Be eager to find the best RED case: the smallest, clearest prompt/repro/example that proves the behavior is missing now and will clearly turn GREEN later.

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

**Capture time is the optimal window for resolving these.** During capture (this action), use the ask tool if your environment provides one; otherwise use your environment's normal ask-user prompt/tool. Present Open Questions immediately. The user is here, engaged, and fleshing out the request — don't defer what you can clarify now. Only leave questions as `- [ ]` if you genuinely can't ask (e.g., batch processing, async capture).

Only add questions where the user's intent is genuinely unclear — don't add questions the builder can answer by reading the codebase.

## Full Context
See `do-work/user-requests/UR-NNN/input.md` for complete verbatim input.
```

**Additional frontmatter for complex requests:**
- `related: [REQ-006, REQ-007]` — other REQs in this batch
- `batch: auth-system` — batch name grouping related requests
- `addendum_to: REQ-005` — if this amends an in-flight/completed request
- `write_set: [src/auth/session.ts, src/auth/*.test.ts]` — repo-relative paths/globs this REQ expects to write

(`maintenance` is **not** complex-only — it lives in the base schema above. Step 1's *Maintenance assessment* is its authoritative home.)

**Populating `depends_on`.** When the request body mentions prior REQs that must complete first (e.g., "after REQ-486 lands", "depends on the auth refactor"), populate `depends_on` in the frontmatter with the REQ IDs. Don't rely on numeric ID ordering — actions/work.md honors `depends_on`, not ID-based heuristics. The optional prose `## Dependencies` section in REQ bodies remains for human readers; the frontmatter field is the source of truth for tooling (work-action selection, roadmap classification, upstream-failure detection).

**Populating `write_set`.** Seed it only when the request already names the files, or when the slice is inherently per-file ("rewrite each adapter in `src/adapters/`"). Otherwise **omit it** — the field is display-only (it feeds the board's *overlaps* badge; nothing schedules on it at any builder count — `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), so an invented set is strictly worse than absence: absence simply gets no overlaps badge (it reads as unknown, not conflict), while a wrong set makes the board badge misleadingly. Capture's value is a hint, never a commitment: the work pipeline's Scope step (`actions/work-reference.md` → Scope Declaration Template) firms it up and overwrites it — for Routes **B and C** only. A **Route A** REQ never runs that step (`actions/work.md` Step 5.5), so it keeps the capture-seeded value for the whole run. The field is optional on every REQ.

**Slicing convention.** When a single user request slices into multiple REQs with internal dependencies, the slicer should populate `depends_on` per the dependency graph it produced. actions/work.md then runs roots first, gates downstream REQs on their prerequisites, and supports `--wave N` for checkpointed execution one dependency depth at a time. A clean DAG in `depends_on` makes foundation-phase batches predictable; sloppy or missing `depends_on` returns to numeric-ID order and risks cascade misclassification. Prefer slice boundaries that also give each REQ its own files — split a surface per concern rather than leaving several REQs editing one shared block — because per-concern files keep each REQ's diff and review self-contained and its board *overlaps* badge quiet. That is a nudge, not a gate: when a file-clean boundary would distort the request, slice for intent and declare the unavoidable overlap in each REQ's `write_set` so the board's overlap badge can surface it.

`depends_on` is semantically distinct from `addendum_to`: `addendum_to: REQ-N` says "this REQ amends REQ-N" (used for follow-ups and review-generated remediation); `depends_on: [REQ-N, REQ-M]` says "this REQ requires REQ-N and REQ-M to be completed first." A REQ can carry both.

### Schema Aliases

Several fields above accept legacy aliases at read time so muscle-memory typos from sister tools don't silently drop information. The canonical key wins when multiple are present; capture always emits the canonical — aliases are read-only, never propagated on write.

| Canonical field | Aliases recognized | Read sites |
|---|---|---|
| `addendum_to` | `amends`, `parent`, `amendment_to` | capture's duplicate check, `actions/work.md`'s upstream-failure walk + cycle detection, `actions/roadmap.md` Blocked classification |
| `depends_on` | `dependencies` | capture's slicing convention, `actions/work.md`'s dependency-aware selection / cycle detection / `--wave` depth / upstream walk, `actions/roadmap.md` Ready/Blocked rubrics |
| `batch` | `batch_name` | `actions/roadmap.md` batch grouping; verify-requests cross-REQ summarization |
| `related` | `related_reqs` | `actions/roadmap.md` cross-REQ surfacing; verify-requests batch coverage |
| `suggested_spec` | `spec_hint`, `suggested-spec` | `actions/work.md`'s spec pre-load hint |

For enum-valued and boolean fields shared with `actions/work.md` (`status`, `domain`, `route`, `caveman`, `tdd`, `maintenance`, `error_type`, `kb_status`), capture honors the **normalize-and-warn contract** defined in the Schema Read Contract (in the companion `actions/work-reference.md`): invalid values trigger a warning and a documented default rather than silent acceptance. When writing the REQ files, if the captured value for any normalize-and-warn field doesn't match the canonical enum (after applying the contract's normalization), prompt the user to confirm the intended value before emitting the REQ — capture is the human-attention window for catching typos at the source. Never write a non-canonical value silently.

### UR input.md

Created for every invocation. For simple requests, it's minimal:

```markdown
---
id: UR-005
title: Add keyboard shortcuts
created_at: 2025-01-26T10:00:00Z  # current UTC instant — date -u +%Y-%m-%dT%H:%M:%SZ, never local time with a Z suffix (Timestamp rule, actions/work-reference.md)
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

## Addendum REQ Template

Written for **Addendum for in-flight/completed requests** (`actions/capture.md` Step 2) — a new UR + REQ pair where the REQ's `addendum_to` links back to the original. The fenced block shows the exact frontmatter and body shape, including the `## Prior Implementation` section used when the original is archived:

```markdown
---
id: REQ-021
title: "Addendum: dark mode sidebar support"
status: pending
created_at: 2025-01-27T09:00:00Z  # current UTC instant — date -u +%Y-%m-%dT%H:%M:%SZ (Timestamp rule)
user_request: UR-006        ← new UR created for this addendum
addendum_to: REQ-005        ← links back to the original request
---

# Addendum: Dark Mode Sidebar Support

## What
Add sidebar support to the existing dark mode implementation (REQ-005).

## Context
Addendum to REQ-005, which is currently [in progress / completed].
The user wants the sidebar to also support dark mode.

## Prior Implementation
[For archived/completed originals: read the original REQ from the archive and
summarize what was built, key files modified, patterns used, and commit hash
(if available). Skip this section for in-flight originals — the builder will
encounter the work in progress naturally.]

## Requirements
- Sidebar must respect the dark mode theme
```
