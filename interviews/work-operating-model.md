---
name: work-operating-model
description: |
  Elicits the operating model of a person at work. Produces agent-ready
  artifacts (USER.md, SOUL.md, HEARTBEAT.md) plus machine-readable exports.
  Based on the five-layer Work Operating Model by Nate B. Jones and
  Jonathan Edwards.
version: 1.1.0
topic_cluster: operating-model
layers:
  - id: operating_rhythms
    title: Operating Rhythms
    order: 1
  - id: recurring_decisions
    title: Recurring Decisions
    order: 2
  - id: dependencies
    title: Dependencies
    order: 3
  - id: institutional_knowledge
    title: Institutional Knowledge
    order: 4
  - id: friction
    title: Friction
    order: 5
exports:
  - path: USER.md
    kind: narrative
  - path: SOUL.md
    kind: decision-framework
  - path: HEARTBEAT.md
    kind: checklist
  - path: operating-model.json
    kind: machine-readable
  - path: schedule-recommendations.json
    kind: machine-readable
---

# Work Operating Model Template

The first job is not to automate the user. It is to help them see and describe how their work actually runs.

## Layer 1: Operating Rhythms

Map how the user's days, weeks, and months actually unfold — not the calendar version, the real one.

### Prompt patterns
- "Walk me through a real Monday from the last two weeks."
- "Where does your calendar lie to you?"
- "When are you actually good for deep work versus admin or reactive work?"
- "What repeats weekly or monthly even when it isn't formally scheduled?"

### Details shape
Every entry's `details` field must include:
- `time_windows` — list of `{start, end, label, days}` objects describing recurring time blocks. `days` is a list of weekday abbreviations drawn from `Mon`, `Tue`, `Wed`, `Thu`, `Fri`, `Sat`, `Sun`. Required so `schedule-recommendations.json` can emit the `days` field without inventing data.
- `energy_pattern` — string describing when the user has energy for which kind of work
- `interruptions` — list of `{source, priority}` objects. `source` names the recurring interruption and who/what it comes from; `priority` is `low`, `medium`, or `high`. Required because `HEARTBEAT.md`'s "What to ignore" section draws from the `low`-priority entries — without the marker the export either violates its schema or invents a ranking.
- `non_calendar_reality` — string describing what actually happens that isn't on the calendar

## Layer 2: Recurring Decisions

Capture the judgment calls the user makes over and over, especially ones where the answer depends on context rather than a checklist.

### Prompt patterns
- "What decisions do you make over and over where the answer depends on context, not a checklist?"
- "What do you look at before you decide?"
- "When do you escalate versus handle it yourself?"
- "Which decisions are reversible if you get them wrong?"

### Details shape
- `decision_name` — short name for the decision
- `decision_inputs` — list of data sources or signals checked
- `thresholds` — list of `{metric, value, direction}` — the numbers that matter
- `escalation_rule` — when the user passes this up or brings someone else in
- `reversible` — boolean: can this decision be undone cheaply

## Layer 3: Dependencies

Map who and what the user waits on, and what breaks when those inputs are late or wrong.

### Prompt patterns
- "What part of your week depends on someone else sending, approving, or clarifying something?"
- "What breaks when that doesn't happen on time?"
- "What's your fallback when you're blocked?"

### Details shape
- `dependency_owner` — person or system the user waits on
- `deliverable` — what they send / approve / provide
- `needed_by` — timing window
- `failure_impact` — what breaks if it's late or wrong
- `fallback` — what the user does when blocked

## Layer 4: Institutional Knowledge

Surface what the user knows that isn't written down anywhere — the context only they carry.

### Prompt patterns
- "What do you know that your team relies on but nobody has really documented?"
- "What mistakes would a smart new hire make because the real context is still in your head?"
- "What would break if you disappeared for two weeks?"

### Details shape
- `knowledge_area` — short name for the domain
- `why_it_matters` — why this context is load-bearing
- `where_it_lives` — "in my head," or a specific partial source
- `who_else_knows` — list of people who partially share this
- `risk_if_missing` — what goes wrong without this knowledge

## Layer 5: Friction

Name the recurring annoyances that eat time — the tooling gaps, the duplicate work, the waits.

### Prompt patterns
- "What keeps eating 10-20 minutes at a time?"
- "Where do you keep doing work the hard way because the systems never quite line up?"
- "What's the same broken handoff you've been complaining about for months?"

### Details shape
- `frequency` — how often this friction hits
- `time_cost` — rough minutes or hours lost per occurrence
- `current_workaround` — what the user does today
- `systems_involved` — tools, services, or people in the friction loop
- `automation_candidate` — boolean: could this reasonably be automated
- `priority` — `low`, `medium`, or `high` when the user is willing to rank

## Cross-layer contradiction checks

During the `review` sub-command, surface these specific tensions:

- **Rhythm vs Dependencies** — A claimed deep-work window that collides with a standing dependency handoff.
- **Decisions vs Knowledge** — A decision rule that cites data the user also said isn't written down anywhere.
- **Friction vs Rhythm** — A recurring friction pattern that implies the stated rhythm isn't real.
- **Dependencies vs Knowledge** — A dependency owner who is the same person the user said carries undocumented context (single point of failure).

## Export Templates

When the `export` sub-command runs against an approved session of this template, render each artifact below using the session's canonical entries. Field references use handlebars-style `{{field}}` syntax against the canonical entry contract plus layer-specific `details`. Iteration uses `{{#each layers.<layer_id>.entries}} … {{/each}}`. Omit sections whose source layer has no qualifying entries.

### `USER.md` — narrative profile

```markdown
# Work Operating Model — {{session.role_or_name_or_repo}}

_Generated {{session.completed_at}}. Based on the work-operating-model template, version {{template.version}}._

## How the week actually runs

{{synthesis_paragraph from operating_rhythms — describe time_windows, energy_pattern, and non_calendar_reality in 2–3 sentences}}

### Deep work windows
{{#each operating_rhythms.entries where details.energy_pattern mentions "deep" or "focus"}}
- {{details.time_windows}} — {{summary}}
{{/each}}

### What the calendar hides
{{#each operating_rhythms.entries}}
- {{details.non_calendar_reality}}
{{/each}}

## Recurring decisions

{{#each recurring_decisions.entries}}
### {{details.decision_name}}
- **Cadence:** {{cadence}}. **Trigger:** {{trigger}}.
- **Inputs:** {{details.decision_inputs}}
- **Thresholds:** {{details.thresholds}}
- **Escalate when:** {{details.escalation_rule}}
- **Reversible:** {{details.reversible}}
{{/each}}

## Dependencies

{{#each dependencies.entries}}
- **{{details.dependency_owner}}** — {{details.deliverable}}, needed {{details.needed_by}}.
  Failure impact: {{details.failure_impact}}. Fallback: {{details.fallback}}.
{{/each}}

## Institutional knowledge I carry

{{#each institutional_knowledge.entries}}
- **{{details.knowledge_area}}** — {{details.why_it_matters}}.
  Currently lives: {{details.where_it_lives}}. Partial sharers: {{details.who_else_knows}}.
  Risk if missing: {{details.risk_if_missing}}.
{{/each}}

## Active friction

{{#each friction.entries sorted by details.priority desc, details.time_cost desc}}
- [{{details.priority}}] **{{title}}** — {{details.frequency}}, ~{{details.time_cost}} per occurrence.
  Workaround: {{details.current_workaround}}. Systems: {{details.systems_involved}}.
  Automation candidate: {{details.automation_candidate}}.
{{/each}}
```

### `SOUL.md` — agent decision framework

```markdown
# Agent Operating Instructions

_Use this file to decide how to act on behalf of the user described in `USER.md`. Do not override these rules with defaults inferred from general context._

## When to act autonomously

{{#each recurring_decisions.entries where details.reversible == true}}
- **{{details.decision_name}}**: apply the thresholds in `USER.md` and act. Do not escalate for this decision class.
{{/each}}

## When to escalate

{{#each recurring_decisions.entries where details.escalation_rule exists and details.escalation_rule != "never"}}
- **{{details.decision_name}}**: {{details.escalation_rule}}
{{/each}}

Additionally, always escalate when:
- A dependency from `USER.md` is late and its fallback is not defined
- An institutional_knowledge item marked `risk_if_missing: high` is needed but the owner is unreachable
- Any threshold in `USER.md` is within 10% of being crossed and the decision is irreversible

## Data sources — trust hierarchy

**Authoritative** (cite these directly, do not second-guess):
{{items appearing in 2+ recurring_decisions.details.decision_inputs}}

**Advisory** (consider, but cross-check before acting):
{{items appearing in only 1 recurring_decisions.details.decision_inputs}}

**Tacit** (do not assume present; ask the user if needed):
{{institutional_knowledge.entries where details.where_it_lives contains "head" or "undocumented"}}

## Tone rules by audience

{{#for each unique stakeholder across all layers}}
- **{{stakeholder}}**: {{derived tone — terse/formal/informal based on which layers they appear in}}
{{/for}}

## "Good enough" thresholds

{{#each recurring_decisions.entries}}
- For **{{details.decision_name}}**: proceed when {{details.thresholds}} are met. Do not hold for perfection.
{{/each}}

## What never to do

- Do not act on behalf of the user in a domain not covered by `USER.md`.
- Do not fabricate information for a decision whose `decision_inputs` are unavailable.
- Do not smooth over contradictions between `USER.md` sections — surface them.
```

### `HEARTBEAT.md` — recurring checklist

```markdown
# Heartbeat Checklist

_Review on a 30-minute cadence. For each item: act, defer, or ignore. Log the decision._

## Every heartbeat

- Scan `USER.md` dependencies. Any expected deliverable past its `needed_by` window?
  - If yes and `fallback` is defined → execute fallback per `SOUL.md`
  - If yes and no fallback → escalate
- Scan `USER.md` recurring decisions. Any whose `cadence` or `trigger` fires now?
  - If yes → pull `decision_inputs`, apply `thresholds`, act or escalate per `SOUL.md`

## First heartbeat after 08:00 local

- Load today's calendar. Compare to deep work windows in `USER.md`. Flag conflicts.
- Scan these sources for overnight changes: {{list derived from operating_rhythms.details.non_calendar_reality + recurring_decisions.details.decision_inputs}}

## First heartbeat Monday after 08:00

- Review last week's friction log from `USER.md`. Any high-priority items unchanged? Flag for user.
- For each `institutional_knowledge` entry with `risk_if_missing: high`: was this knowledge used last week? By whom? Log.

## First heartbeat on the 1st of the month

- Produce a one-page delta: what in `USER.md` no longer matches reality? Flag for the user's next quarterly interview re-run.
```

### `operating-model.json` — machine-readable dump

```json
{
  "template": "work-operating-model",
  "template_version": "{{template.version}}",
  "session_id": "{{session.session_id}}",
  "generated_at": "{{session.completed_at}}",
  "previous_version": "{{session.previous_version}}",
  "layers": {
    "operating_rhythms": { "entries": [ "{{canonical_entries}}" ] },
    "recurring_decisions": { "entries": [ "{{canonical_entries}}" ] },
    "dependencies": { "entries": [ "{{canonical_entries}}" ] },
    "institutional_knowledge": { "entries": [ "{{canonical_entries}}" ] },
    "friction": { "entries": [ "{{canonical_entries}}" ] }
  }
}
```

Canonical entries are serialized verbatim with all 11 required fields from the entry contract.

### `schedule-recommendations.json` — derived scheduling data

```json
{
  "generated_at": "{{session.completed_at}}",
  "source_template": "work-operating-model",
  "source_session": "{{session.session_id}}",
  "time_blocks": [
    {
      "label": "{{derived from operating_rhythms.details.energy_pattern}}",
      "days": ["{{day}}"],
      "start": "HH:MM",
      "end": "HH:MM",
      "type": "deep_work | admin | reactive",
      "source_entries": ["operating_rhythms.<entry_id>"]
    }
  ],
  "avoid_windows": [
    {
      "label": "{{why this window should be protected}}",
      "days": ["{{day}}"],
      "start": "HH:MM",
      "end": "HH:MM",
      "reason": "{{non_calendar_reality or friction source}}",
      "source_entries": ["<layer>.<entry_id>"]
    }
  ],
  "standing_slots": [
    {
      "label": "{{deliverable or handoff}}",
      "cadence": "{{from dependency entry}}",
      "day": "{{day or 'rolling'}}",
      "time": "HH:MM",
      "counterparty": "{{dependency_owner}}",
      "source_entries": ["dependencies.<entry_id>"]
    }
  ]
}
```

Derivation rules:
- `time_blocks` come from `operating_rhythms.entries[*].details.time_windows` joined against `details.energy_pattern` to classify type.
- `avoid_windows` come from (a) `operating_rhythms.details.non_calendar_reality` when it describes a consistent interruption pattern, and (b) high-priority `friction.entries` whose `systems_involved` implies a recurring time loss.
- `standing_slots` come from `dependencies.entries[*]` where `needed_by` has a regular cadence.

## Tone

Direct, practical, specific. No generic productivity advice. No fake certainty. Keep momentum moving without bulldozing confirmation.
