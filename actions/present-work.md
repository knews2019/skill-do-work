# Present Work Action

> **Part of the do-work skill.** Generates client-facing deliverables from completed work — briefs, architecture diagrams, value propositions, and video scripts. Turns the technical archive into something that educates and sells.

The code is done. Now communicate its value. This action reads the full history of completed work — what was requested, what was built, how it works, and what the Lessons Learned say — and produces artifacts that explain it to someone who doesn't read diffs.

## Philosophy

- **Educate first, sell second.** The client should understand what they got before hearing why it's valuable.
- **Technical accuracy in plain language.** Architecture and data flow matter, but describe them at the level of components and interactions, not code.
- **Honest value, not hype.** Say "enables X" not "will increase revenue by Y%." Don't fabricate metrics.
- **Pointers over prose.** Reference key files and code — the implementation is the source of truth. Don't rewrite it in paragraphs.
- **Proportional effort.** A config change gets a 2-paragraph brief. A multi-feature system gets architecture diagrams and a video script.

## Two Modes

| Mode | Trigger | What it does |
|------|---------|-------------|
| **Detail** | `do work present UR-003`, `do work present REQ-005`, or `do work present` (most recent) | Deep dive on specific completed work |
| **Portfolio** | `do work present all` or `do work present portfolio` | Summary of all completed work across the archive |

## Detail Mode Workflow

### Step 1: Find the Target

Same pattern as review work standalone mode:

1. **If user specifies a REQ** (e.g., "present REQ-005"): Find it in `do-work/archive/` or `do-work/archive/UR-NNN/`
2. **If user specifies a UR** (e.g., "present UR-003"): Find all completed REQs under that UR — present them as one deliverable
3. **If no target specified**: Find the most recently completed UR (or REQ if no UR). Check `do-work/archive/` for the highest UR/REQ number with `status: completed`

If the target has no completed REQs, report that there's nothing to present and exit.

### Step 2: Read the Full History

For each completed REQ, read the full file and extract:

- **What was requested** — the What/Detailed Requirements sections
- **Original input** — read the UR's `input.md` for the user's own words
- **Triage** — what route was chosen and why
- **Plan** — what was planned (Route C)
- **Exploration** — what was discovered about the codebase
- **Implementation Summary** — what the builder says it did
- **Testing** — what tests exist and pass
- **Review** — scores, findings, acceptance result
- **Lessons Learned** (if present) — what worked, what didn't, key files, gotchas. Route A REQs may skip this section.

### Step 3: Read the Code

Use `git show <commit>` (from the REQ's `commit` frontmatter) to get the diff. Then read the actual created/modified files to understand:

- **Architecture** — what components exist, their roles, how they connect
- **Data flow** — how information enters, moves through, and exits the system
- **Patterns used** — frameworks, libraries, conventions followed
- **Scale** — how many files, how much new code, how much modified

Don't just summarize the diff — understand the system that was built.

### Step 4: Generate Artifacts

Based on the work scope, generate the appropriate deliverables:

#### 4a: Client Brief (always generated)

Write to `do-work/deliverables/UR-NNN-client-brief.md` (or `REQ-NNN-client-brief.md` for standalone REQs):

```markdown
# [Feature/Project Name]

## What We Built

[1-2 paragraphs. Plain language. What it does from the user's perspective.
No code, no jargon — what problem it solves and how the user interacts with it.]

## How It Works

### Architecture

[ASCII diagram: components, their roles, how they connect.
Use boxes, arrows, labels. Keep it to one screen.]

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  User Input │────►│  Processing  │────►│   Output    │
│  (Form/API) │     │  (Service)   │     │  (UI/File)  │
└─────────────┘     └──────┬───────┘     └─────────────┘
                           │
                    ┌──────▼───────┐
                    │   Storage    │
                    │  (Database)  │
                    └──────────────┘
```

### Data Flow

[Numbered walkthrough showing how data moves through the system:]
1. User submits [input] via [interface]
2. [Component] validates and transforms the data
3. [Service] processes it by [doing what]
4. Result is stored in [where] and returned to [whom]

### Key Design Decisions

- [Why this approach over alternatives — what made it the right call]
- [Trade-offs made and why they're acceptable]
- [What the Lessons Learned revealed about the approach]

## Why This Works

[1 paragraph. Business-level justification. Why this solution solves the
problem well, not just technically but for the business.]

## Value Delivered

### Immediate Impact
- [What the client gets right now — specific capabilities]
- [Problems solved — pain points eliminated]

### Revenue & Growth Opportunities
- [How this enables the client to make more money]
- [New capabilities unlocked — what's now possible that wasn't before]
- [Efficiency gains — time saved, errors prevented, processes streamlined]

### Competitive Advantage
- [What this gives the client that their competitors likely don't have]

## Key Files

[Pointers to the most important files in the implementation.
These are the source of truth — read them for the full picture.]

- `src/components/Feature.tsx` — main component
- `src/services/feature-service.ts` — business logic
- `tests/feature.spec.ts` — test coverage

## What's Next

- [Follow-up opportunities — natural extensions of this work]
- [Phase 2 ideas — features that build on what was just shipped]
- [Quick wins — small additions that compound the value]
```

#### 4b: Video Script (when the feature is user-facing/demo-able)

Write to `do-work/deliverables/UR-NNN-video-script.md`:

```markdown
# Video Script: [Feature Name]

**Duration:** ~[X] minutes
**Format:** Screen recording with narration (Remotion / Loom / manual)
**Audience:** [Client stakeholders / end users / both]

---

## Scene 1: The Problem (~10-15s)

**Visual:** [What to show on screen — the existing pain point, gap, or workflow without this feature]
**Narration:** "[Setup — what was missing, broken, or manual. Keep it relatable.]"
**Transition:** [How to move to the next scene]

## Scene 2: The Solution (~20-30s)

**Visual:** [Demo the feature — show the happy path, step by step]
**Narration:** "[Walk through what's happening on screen. Match narration to actions.]"
**Key moments:**
- [Moment 1: The "aha" — where the value becomes obvious]
- [Moment 2: The detail — a specific capability that impresses]

## Scene 3: Under the Hood (~15-20s)

**Visual:** [Architecture diagram or data flow from the client brief]
**Narration:** "[Brief technical explanation — how it works at a high level. No code. Confidence-building, not education.]"

## Scene 4: The Value (~10-15s)

**Visual:** [Key capability comparison: before vs. after, or a list of what's now possible]
**Narration:** "[Why this matters to the business. End on the revenue/growth angle.]"

---

**Production notes:**
- [Screen resolution / browser / app state needed for recording]
- [Test data to use for the demo]
- [Anything to set up before recording]
```

**When to generate a video script:**
- The feature has visible UI or user-facing output
- There's a clear before/after to demonstrate
- The work is substantial enough to warrant a walkthrough (Route B or C)

**When to skip:**
- Backend-only changes, config tweaks, refactors, infrastructure
- Bug fixes that aren't visually interesting
- Route A changes (too small)

#### 4c: Portfolio artifacts (portfolio mode only — see below)

### Step 5: Save and Present

1. Create `do-work/deliverables/` directory if it doesn't exist
2. Write all generated artifacts to `do-work/deliverables/`
3. Present a summary to the user:

```
Generated deliverables for UR-003:

  do-work/deliverables/UR-003-client-brief.md      Client brief with architecture + value prop
  do-work/deliverables/UR-003-video-script.md       Video script (4 scenes, ~90s)

Key value points:
  - [Top 1-2 value propositions from the brief]
```

## Portfolio Mode Workflow

### Step 1: Scan the Archive

List all UR folders and completed REQs in `do-work/archive/`:
- Read each UR's `input.md` for the title and request list
- Read each completed REQ's frontmatter (title, route, commit, completed_at) and Review section (overall score)
- Check `do-work/archive/legacy/` for standalone completed REQs

### Step 2: Build the Overview

For each UR/REQ, extract a one-line summary from the What section and the review score.

### Step 3: Generate Portfolio Summary

Write to `do-work/deliverables/portfolio-summary.md`:

```markdown
# Work Portfolio

**Total completed:** [N] requests across [M] user requests
**Period:** [earliest created_at] — [latest completed_at]

---

## Completed Work

### UR-001: [Title from input.md]

[1-2 sentences on what was delivered — synthesized from the REQs' What sections]

| REQ | Title | Route | Review | Commit |
|-----|-------|-------|--------|--------|
| REQ-010 | [title] | C | 92% | abc1234 |
| REQ-011 | [title] | B | 88% | def5678 |

---

### UR-003: [Title]

[1-2 sentences]

| REQ | Title | Route | Review | Commit |
|-----|-------|-------|--------|--------|
| REQ-020 | [title] | A | 95% | ghi9012 |

---

## Cumulative Value Proposition

### What Was Built
[High-level summary across all completed work — the big picture.
What capabilities does the system now have?]

### Total Value Delivered
- [Aggregated business impact across all work]
- [Capabilities added — the full list]
- [Technical quality — average review scores, test coverage]

### Growth Opportunities
- [Cross-cutting opportunities that only become visible when viewing all work together]
- [Natural next phases based on what's been built]
- [Quick wins that compound existing value]

### Lessons Learned (Cross-Project)
- [Patterns that worked well across multiple REQs]
- [Common pitfalls encountered]
- [Architectural decisions that should inform future work]
```

### Step 4: Save and Present

Same as detail mode — save to `do-work/deliverables/` and summarize.

## Calibrating Depth

| Work scope | Artifact depth |
|-----------|---------------|
| Single Route A REQ (bug fix, config) | Minimal brief — 2-3 paragraphs, skip architecture diagram, skip video script, minimal value prop. |
| Single Route B REQ (clear feature) | Standard brief with architecture section. Video script if user-facing. |
| Route C or multi-REQ UR | Full brief with detailed architecture, data flow, value prop, design decisions. Video script for demo-able features. |
| Portfolio (all work) | Portfolio summary with per-UR breakdowns and cumulative value proposition. |

## What NOT to Do

- Don't include code snippets — this is client-facing, not developer-facing. Point to files instead.
- Don't oversell — be honest about what was built and its limitations
- Don't fabricate metrics — quantify where data exists, qualify where it doesn't
- Don't include internal details (review scores, triage routes, test commands) in client briefs — those are for the portfolio summary
- Don't generate a video script for non-visual changes
- Don't write walls of text when pointers to code would be more accurate and durable
- Don't regenerate deliverables that already exist without the user asking — check `do-work/deliverables/` first and offer to update
