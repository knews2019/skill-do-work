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

## Tone

Direct, practical, specific. No generic productivity advice. No fake certainty. Keep momentum moving without bulldozing confirmation.
