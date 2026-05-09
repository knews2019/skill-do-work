---
name: work-operating-model
description: |
  Elicits the operating model of a person at work. Produces agent-ready
  artifacts (USER.md, SOUL.md, HEARTBEAT.md) plus machine-readable exports.
  Based on the five-layer Work Operating Model by Nate B. Jones and
  Jonathan Edwards.
version: 2.0.0
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

## Migration from v1.x

Two breaking shape changes since v1.0.0:

- `details.interruptions` — was `list[string]`, now `list[{source: string, priority: "low" | "medium" | "high"}]`. Required because `HEARTBEAT.md`'s "What to ignore" section draws from `low`-priority entries.
- `details.time_windows` — gained a required `days` field (`list` of weekday abbreviations from `Mon … Sun`). Required so `schedule-recommendations.json` can emit `days` without inventing data.

In-flight v1.x sessions will fail validation against v2.0. The `actions/interview.md` Step 2 migration check runs these steps automatically when it sees a session whose `template_version` is missing or older than `2.0.0`:

1. For each `operating_rhythms` entry's `details.interruptions`, replace each string `"<source>"` with `{"source": "<source>", "priority": "medium"}` (use `"medium"` as a safe default; the user revisits during the next `review` pass).
2. For each `details.time_windows` entry, add `"days": ["Mon", "Tue", "Wed", "Thu", "Fri"]` (or the actual days the window applies). Without `days`, the export will refuse.
3. Set `template_version` in `session.json` to `2.0.0`.

Sessions stamped `2.0.0` or higher load without migration. If a hand migration is preferred, the same three steps applied directly to `./do-work/interview/work-operating-model/session.json` produce an equivalent state — the action treats hand-migrated sessions and auto-migrated sessions identically once `template_version: 2.0.0` is recorded. To skip migration entirely and re-elicit from scratch, run `do-work interview work-operating-model reset`.

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

When the `export` sub-command runs against an approved session of this template, render each artifact below using the session's canonical entries.

**Dialect.** Templates are mechanical handlebars with three constructs and nothing else:

1. **Field substitution** — `{{path.to.field}}` resolves against the canonical entry contract plus layer-specific `details`.
2. **Iteration** — `{{#each <expr>}} … {{/each}}` over a list. `<expr>` may be a path (`recurring_decisions.entries`) or a `where` / `sorted by` extension (`recurring_decisions.entries where status != "stale" and details.reversible == true`, `friction.entries where status != "stale" sorted by details.priority desc, details.time_cost desc`). The `where` predicate uses field comparisons (`==`, `!=`, `exists`, `contains`, `mentions`) and boolean `and`/`or`. `sorted by` takes one or more `<field> asc|desc`.
3. **Synthesized fields** — `{{derived.<name>}}` resolves to a value pre-computed by the renderer per the **Synthesized fields** block listed immediately under each template. These are the only places prose generation enters the pipeline; everywhere else is mechanical substitution.

Omit sections whose source layer has no qualifying entries. Anything that looks like a natural-language directive embedded inside `{{ … }}` is a bug — promote it to a `{{derived.<name>}}` slot and document it in the template's Synthesized fields block.

### `USER.md` — narrative profile

```markdown
# Work Operating Model — {{session.role_or_name_or_repo}}

_Generated {{session.last_exported_at}}. Based on the work-operating-model template, version {{template.version}}._

## How the week actually runs

{{derived.rhythm_synthesis}}

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

{{#each friction.entries where status != "stale" sorted by details.priority desc, details.time_cost desc}}
- [{{details.priority}}] **{{title}}** — {{details.frequency}}, ~{{details.time_cost}} per occurrence.
  Workaround: {{details.current_workaround}}. Systems: {{details.systems_involved}}.
  Automation candidate: {{details.automation_candidate}}.
{{/each}}

## Stale or deprecated

_Entries flagged as stale during an `update` run. Kept for context but not applied as active rules._

{{#each all_layers.entries where status == "stale"}}
- **[{{layer_id}}] {{title}}** — {{summary}}. (last validated {{last_validated_at}})
{{/each}}
```

**Synthesized fields for `USER.md`:**

- `derived.rhythm_synthesis` — 2–3 sentence prose paragraph that names the time windows, energy pattern, and non-calendar reality across all `operating_rhythms.entries where status != "stale"`. Specific, not generic; cites the user's own phrasing where it appeared in the session.

### `SOUL.md` — agent decision framework

```markdown
# Agent Operating Instructions

_Use this file to decide how to act on behalf of the user described in `USER.md`. Do not override these rules with defaults inferred from general context._

## When to act autonomously

{{#each recurring_decisions.entries where status != "stale" and details.reversible == true}}
- **{{details.decision_name}}**: apply the thresholds in `USER.md` and act. Do not escalate for this decision class.
{{/each}}

## When to escalate

{{#each recurring_decisions.entries where status != "stale" and details.escalation_rule exists and details.escalation_rule != "never"}}
- **{{details.decision_name}}**: {{details.escalation_rule}}
{{/each}}

Additionally, always escalate when:
- A dependency from `USER.md` is late and its fallback is not defined
- An institutional_knowledge item marked `risk_if_missing: high` is needed but the owner is unreachable
- Any threshold in `USER.md` is within 10% of being crossed and the decision is irreversible

## Data sources — trust hierarchy

**Authoritative** (cite these directly, do not second-guess):
{{#each derived.authoritative_inputs}}
- {{this}}
{{/each}}

**Advisory** (consider, but cross-check before acting):
{{#each derived.advisory_inputs}}
- {{this}}
{{/each}}

**Tacit** (do not assume present; ask the user if needed):
{{#each derived.tacit_knowledge}}
- {{title}} — {{details.where_it_lives}}
{{/each}}

## Tone rules by audience

{{#each derived.stakeholder_tones}}
- **{{stakeholder}}**: {{tone}}
{{/each}}

## "Good enough" thresholds

{{#each recurring_decisions.entries where status != "stale"}}
- For **{{details.decision_name}}**: proceed when {{details.thresholds}} are met. Do not hold for perfection.
{{/each}}

## What never to do

- Do not act on behalf of the user in a domain not covered by `USER.md`.
- Do not fabricate information for a decision whose `decision_inputs` are unavailable.
- Do not smooth over contradictions between `USER.md` sections — surface them.
```

**Synthesized fields for `SOUL.md`:**

- `derived.authoritative_inputs` — list of input names that appear in `details.decision_inputs` of **2 or more** `recurring_decisions.entries where status != "stale"`. Deduplicated, sorted alphabetically.
- `derived.advisory_inputs` — list of input names that appear in `details.decision_inputs` of **exactly one** `recurring_decisions.entries where status != "stale"`. Deduplicated, sorted alphabetically.
- `derived.tacit_knowledge` — `institutional_knowledge.entries where status != "stale"` whose `details.where_it_lives` contains the substring `"head"` or `"undocumented"` (case-insensitive). Each item exposes the entry's full canonical fields plus `details`.
- `derived.stakeholder_tones` — list of `{stakeholder, tone}` objects, one per unique stakeholder name appearing in any layer's entries. `tone` is one of `terse`, `formal`, `informal`, derived from which layers the stakeholder appears in: appears only in `dependencies` → `formal`; appears in `friction` or `recurring_decisions` and not `dependencies` → `terse`; appears only in `institutional_knowledge` → `informal`. Stakeholders appearing across mixed categories use the strongest signal in that order (formal > terse > informal).

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
- Scan these sources for overnight changes:
{{#each derived.overnight_scan_sources}}
  - {{this}}
{{/each}}

## First heartbeat Monday after 08:00

- Review last week's friction log from `USER.md` (entries where `status != "stale"`). Any high-priority items unchanged? Flag for user.
- For each `institutional_knowledge` entry with `status != "stale"` and `risk_if_missing: high`: was this knowledge used last week? By whom? Log.

## First heartbeat on the 1st of the month

- Produce a one-page delta: what in `USER.md` no longer matches reality? Flag for the user's next quarterly interview re-run.
```

**Synthesized fields for `HEARTBEAT.md`:**

- `derived.overnight_scan_sources` — deduplicated list of strings combining (a) `details.non_calendar_reality` from `operating_rhythms.entries where status != "stale"` and (b) every entry in `details.decision_inputs` across `recurring_decisions.entries where status != "stale"`. Sorted alphabetically.

### `operating-model.json` — machine-readable dump

```json
{
  "template": "work-operating-model",
  "template_version": "{{template.version}}",
  "session_id": "{{session.session_id}}",
  "generated_at": "{{session.last_exported_at}}",
  "previous_version": "{{session.previous_version}}",
  "layers": {
    "operating_rhythms": { "entries": {{json_entries operating_rhythms}} },
    "recurring_decisions": { "entries": {{json_entries recurring_decisions}} },
    "dependencies": { "entries": {{json_entries dependencies}} },
    "institutional_knowledge": { "entries": {{json_entries institutional_knowledge}} },
    "friction": { "entries": {{json_entries friction}} }
  }
}
```

`{{json_entries <layer>}}` emits a **JSON array** (not a quoted string) whose elements are the layer's canonical entry objects serialized verbatim with all 11 required fields from the entry contract — including `status`, so stale entries are preserved in the machine-readable dump. Consumers that want only active entries should filter by `status != "stale"` themselves.

### `schedule-recommendations.json` — derived scheduling data

```json
{
  "generated_at": "{{session.last_exported_at}}",
  "source_template": "work-operating-model",
  "source_session": "{{session.session_id}}",
  "time_blocks": [
    {{#each derived.time_blocks}}
    {
      "label": "{{label}}",
      "days": {{json days}},
      "start": "{{start}}",
      "end": "{{end}}",
      "type": "{{type}}",
      "source_entries": {{json source_entries}}
    }{{#unless @last}},{{/unless}}
    {{/each}}
  ],
  "avoid_windows": [
    {{#each derived.avoid_windows}}
    {
      "label": "{{label}}",
      "days": {{json days}},
      "start": "{{start}}",
      "end": "{{end}}",
      "reason": "{{reason}}",
      "source_entries": {{json source_entries}}
    }{{#unless @last}},{{/unless}}
    {{/each}}
  ],
  "standing_slots": [
    {{#each derived.standing_slots}}
    {
      "label": "{{label}}",
      "cadence": "{{cadence}}",
      "day": "{{day}}",
      "time": "{{time}}",
      "counterparty": "{{counterparty}}",
      "source_entries": {{json source_entries}}
    }{{#unless @last}},{{/unless}}
    {{/each}}
  ]
}
```

**Synthesized fields for `schedule-recommendations.json`:**

- `derived.time_blocks` — list of `{label, days, start, end, type, source_entries}` from `operating_rhythms.entries[*].details.time_windows` (where `status != "stale"`). `label` echoes the time window's `label`. `days` is the window's `days` array. `start` / `end` are `HH:MM`. `type` is `deep_work | admin | reactive`, classified by joining the entry's `details.energy_pattern` against the window's label (focus/deep → `deep_work`; admin/email/reactive → `admin` or `reactive`). `source_entries` is `["operating_rhythms.<entry_id>"]`.
- `derived.avoid_windows` — list of `{label, days, start, end, reason, source_entries}` drawn from two sources: (a) `operating_rhythms.entries where status != "stale"` whose `details.non_calendar_reality` describes a consistent recurring interruption pattern (label = a short summary, reason = the `non_calendar_reality` string), and (b) `friction.entries where status != "stale" and details.priority == "high"` whose `details.systems_involved` implies a recurring time loss (label = entry `title`, reason = entry `summary`). `days` / `start` / `end` come from the source entry's nearest time window when one exists; otherwise omit the entry rather than invent times.
- `derived.standing_slots` — list of `{label, cadence, day, time, counterparty, source_entries}` from `dependencies.entries where status != "stale"` whose `details.needed_by` parses to a regular cadence (daily, weekly, monthly, or named day). `label` = `details.deliverable`; `cadence` = the parsed cadence string; `day` = weekday abbreviation or `"rolling"`; `time` = `HH:MM` or `""` if not specified; `counterparty` = `details.dependency_owner`; `source_entries` = `["dependencies.<entry_id>"]`.

The `{{json X}}` helper emits `X` as a JSON value (array or object) rather than a quoted string, identical in semantics to `{{json_entries X}}` used in `operating-model.json`.

## Tone

Direct, practical, specific. No generic productivity advice. No fake certainty. Keep momentum moving without bulldozing confirmation.
