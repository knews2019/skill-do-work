# Roadmap Action

> **Part of the do-work skill.** Invoked when the user wants a survey of the queue: what's done, what's in progress, what's pending, and a feasibility read on what's actionable next. Read-only — never modifies REQs, frontmatter, or files.

A planning aid, not a diagnostic. Where `forensics` looks for *broken* state (stuck, hollow, orphaned), `roadmap` looks at *intended* state — the shape of remaining work and whether it's ready to be picked up. Run it when you want to know "where are we, and what's worth doing next?"

## When to Use

**Use when:**
- User asks "what's left?", "what's done?", "where are we?", "status", "queue-status", "queue status", "roadmap", "what's feasible?", "what should I work on next?"
- Planning a session and wanting to scope what's actionable
- Onboarding to a project with existing `do-work/` history and wanting a lay of the land

**Do NOT use when:**
- User suspects something is *broken* or *stuck* — route to the forensics action instead
- User wants to *generate new ideas* for what to build — route to the scan-ideas action instead
- User wants to *review specific completed code* — route to the review-work or code-review action instead
- User wants to *explain uncommitted local changes* — route to the inspect action instead

## Core Rules

- **Read-only.** Never modifies REQs, moves files, updates frontmatter, or creates commits.
- **Feasibility is a read, not a verdict.** Flag concerns; don't unilaterally reclassify a REQ as blocked.
- **Cite evidence.** Every feasibility judgment should point to a specific frontmatter field, section, or missing artifact.
- **No ideation.** This action surveys what already exists in the queue. If nothing's queued, say so — don't propose new work.

## Input

Optional argument may scope the report:

- *(no argument)* — full survey across queue, working, and archive
- `pending` — only `do-work/queue/` REQs (what's actionable now)
- `in-progress` — only `do-work/working/` REQs
- `done` — only archived REQs
- `UR-NNN` — scope to a single user request and its REQs
- `since <date>` — filter archive entries to those completed on/after the date

If the argument is unrecognized, default to the full survey and note the unrecognized argument in the report.

## Steps

### Step 1: Inventory

Walk the do-work tree and collect:

- `do-work/queue/REQ-*.md` — pending and `pending-answers`
- `do-work/working/REQ-*.md` — actively claimed
- `do-work/archive/**/REQ-*.md` — terminal status (completed, completed-with-issues, failed)
- `do-work/user-requests/UR-*/` — open URs and their referenced REQs

For each REQ, capture: id, title, status, route (if set), `user_request`, `created_at`, `claimed_at`, `completed_at`, `domain`, `addendum_to`, `kb_status`, `tdd` (frontmatter, default false if absent), and which `##` sections exist (note especially the presence/absence of `## Red-Green Proof`).

### Step 2: Classify Pending Work by Feasibility

For each REQ in `do-work/queue/`, assign a feasibility bucket using only what's visible in the file:

- **Ready** — has a clear `## What`, no `pending-answers` status, no unresolved `addendum_to` chain, dependencies (if listed) point to archived/completed REQs.
- **Needs clarification** — `status: pending-answers`, OR the request body contains explicit open questions, OR scope is too vague to triage (one-line title with no `## What` body).
- **Blocked** — references a REQ in `addendum_to` or a dependencies list that is still pending or in-progress; or names an external dependency (waiting on an API, a decision, a third-party).
- **Stale** — `created_at` more than 30 days old AND not yet claimed. Flag for re-confirmation; the user may no longer want it.

Each classification must cite the specific evidence that drove it (e.g., "status: pending-answers", "addendum_to: REQ-031 (still pending)", "no `## What` section").

### Step 2.5: Assess TDD Posture for Pending Work

For each REQ in `do-work/queue/`, classify TDD posture using only frontmatter and section evidence:

- **TDD on** — `tdd: true` in frontmatter. Note whether `## Red-Green Proof` exists (mandatory for TDD-on per the capture contract; flag as a gap if missing).
- **TDD eligible** — `tdd: false` or absent, but the REQ describes testable behavior the heuristic would flag (pure logic, data transformations, API handlers, utility functions, behavior-changing bug fixes). Strong signals: a `## Red-Green Proof` section, an explicit input/output example in `## What`, or `domain: backend | testing`. Surface as "could turn TDD on."
- **TDD not applicable** — `tdd: false` or absent, and the REQ is UI layout, copy/content, config tweak, glue code, or pure refactor. Don't surface a recommendation; just record the posture.

Cite the specific evidence that drove the classification (e.g., "tdd: true + Red-Green Proof present", "tdd: false but Red-Green Proof present → eligible", "domain: ui-design, copy change → not applicable"). Never reclassify the frontmatter — this is a read.

### Step 3: Roll Up Completed Work

For each REQ in `do-work/archive/`:

- Group by `user_request` (UR-NNN) and by completion week.
- Note any UR with all REQs completed (candidate for UR archival — surface, don't act).
- Note `kb_status: pending` lessons that haven't been promoted (candidate for `do-work bkb triage`).
- Record `tdd` posture per REQ so completed work shows whether tests went in test-first.

### Step 4: Highlight In-Progress Work

For each REQ in `do-work/working/`:

- Report id, title, route, current phase (most recent `##` section), how long claimed, and `tdd` posture (on/off).
- Do **not** flag stuck work here — that's forensics' job. Just report state.

### Step 5: Compose the Report

Render the report per the Output Format below. Lead with the actionable section (what's Ready) so the reader can act on it without scrolling.

## Output Format

```markdown
# Roadmap

**Scan date:** [timestamp]
**Scope:** [full | pending | in-progress | done | UR-NNN | since <date>]
**Totals:** [N ready] · [N needs clarification] · [N blocked] · [N in-progress] · [N completed] · [N failed]
**TDD posture (pending):** [N on] · [N eligible] · [N not applicable]

## Ready to Pick Up

- **REQ-NNN — <title>** (UR-NNN, route: <route or "untriaged">, tdd: on | eligible | n/a)
  Brief one-line scope summary. Evidence: <why it's ready>. TDD: <on with proof | eligible — suggest enabling | n/a>.

## Needs Clarification

- **REQ-NNN — <title>** (status: pending-answers, age: 4d, tdd: on | eligible | n/a)
  Open questions: <count or summary>. Suggested next step: `do-work clarify`.

## Blocked

- **REQ-NNN — <title>** (depends on REQ-MMM, still pending)
  Unblock when REQ-MMM lands.

## Stale

- **REQ-NNN — <title>** (created 47d ago, never claimed)
  Re-confirm relevance with the user before working.

## TDD Eligible (Could Turn On)

REQs where `tdd: false` but the behavior is testable and a test-first approach would apply.

- **REQ-NNN — <title>** (UR-NNN)
  Signal: <Red-Green Proof present | input/output example in What | domain: backend>. To enable, set `tdd: true` and confirm `## Red-Green Proof`.

## In Progress

- **REQ-NNN — <title>** (route: <route>, claimed: 2h ago, phase: Implementation, tdd: on | off)

## Recently Completed

Grouped by UR or by week:

- **UR-NNN — <ur title>** — 4/4 REQs complete (candidate for UR archival)
  - REQ-NNN <title> (commit: abc1234, tdd: on | off)
  - REQ-NNN <title> (commit: def5678, tdd: on | off)

## Lessons Awaiting Promotion

- REQ-NNN — kb_status: pending, kb_entry: <filename or none>
  Suggested next step: `do-work bkb triage`.

## Suggested Next Steps

1. Pick up REQ-NNN (top of Ready) — clearest scope, no blockers.
2. Run `do-work clarify` to drain the N pending-answers REQs.
3. Consider enabling `tdd: true` on the N TDD-eligible REQs before they're picked up.
4. Confirm or discard the N stale REQs with the user.
```

Omit sections with no entries. If the queue is empty and nothing is in-progress, report:

```markdown
# Roadmap

**Scan date:** [timestamp]
**Scope:** [scope]

Queue is empty — no pending or in-progress work.

[If archive non-empty: brief summary of recent completions.]
[If archive empty: "No archived work yet — run `do-work capture` to add a request."]
```

## Rules

- **Cap each section at 20 entries** by default; if more exist, list the top 20 and note "(N more — narrow scope with an argument)".
- **Don't editorialize on REQ quality** — the verify-requests action handles that. Roadmap reports state, not content health.
- **Don't recommend code changes.** Suggested next steps must be do-work commands or human decisions, not implementation work.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll quickly fix this REQ's stale `pending-answers` status while I'm here." | Report it under Needs Clarification and stop. | Roadmap is read-only; mutations belong to clarify/cleanup. |
| "Let me re-triage this untriaged REQ since I'm reading it anyway." | Note it as Ready (untriaged) and let `do-work work` handle triage. | Triage is part of the work action's contract; doing it here splits responsibility. |
| "These three REQs look like duplicates — I'll consolidate them." | Surface the overlap under Suggested Next Steps and let the user decide. | Consolidation requires user judgment; roadmap is a survey, not an editor. |
| "Nothing is Ready, so I'll suggest new ideas." | Say "no ready work" and stop. | Ideation is `scan-ideas`. Roadmap surveys what exists. |
| "This REQ looks testable — I'll set `tdd: true` myself." | Surface it under TDD Eligible and let the user decide. | Roadmap is read-only; `tdd` is an authoring decision. |

## Red Flags

- Report classified a REQ as Blocked but cited no specific dependency — feasibility judgment without evidence.
- Roadmap modified any file under `do-work/` — read-only contract violated.
- Suggested Next Steps recommends writing code rather than running a do-work action or making a human decision.
- Report duplicated forensics findings (stuck work, hollow completions) — wrong action; redirect to forensics.
- Every pending REQ landed in the same bucket — classifier is degenerate; review the rubric.
- A REQ shipped with `tdd: true` but no `## Red-Green Proof` section, and the report missed it.
- TDD posture reported but no evidence cited (frontmatter value, proof section, domain, or What-section example).

## Verification Checklist

- [ ] Zero changes to `do-work/` — read-only contract held.
- [ ] Every feasibility classification cites concrete evidence (frontmatter field, section, or absence thereof).
- [ ] Ready section appears first; user can act on it without scrolling.
- [ ] Empty sections were omitted, not rendered with "(none)".
- [ ] Suggested Next Steps lists do-work commands or human decisions, not code work.
- [ ] If scope argument was unrecognized, the report notes it and falls back to full survey.
- [ ] Every pending REQ has a TDD posture (on / eligible / n/a) with cited evidence; the totals line in the header reflects the same counts.
