# Interview Guide

Sit down for 45 focused minutes and produce agent-ready operating artifacts. The interview action runs a structured, multi-layer conversation that turns tacit work knowledge into files an agent (or a new hire) can actually act on.

## When to run it

- You want to hand work off — to a delegate, a collaborator, an agent.
- You're onboarding an AI assistant and need it to decide on your behalf.
- You want your judgment patterns explicit for yourself — to catch your own inconsistencies.
- You're moving roles or restructuring responsibilities and the old operating model no longer fits.

This is not a productivity coaching session. The interview describes reality; it does not suggest improvements.

## Expected time

- First run of the `work-operating-model` template: **~45 minutes** of focused attention. Block the time on your calendar and treat it like a meeting with yourself.
- Possibly more on the first run if the user surfaces more than 5 entries per layer.
- Subsequent `update` runs (quarterly recheck): **~15 minutes** if nothing major has changed.

## Output files

After `do-work interview work-operating-model export`, you get five files in `./do-work/interview/work-operating-model/exports/`:

| File | What it is |
|---|---|
| `USER.md` | Narrative profile of who you are at work — role, rhythms, decisions, dependencies, knowledge, frictions. |
| `SOUL.md` | Decision framework an agent follows acting for you — when to escalate, when to decide, tone rules, trusted data sources, "good enough" thresholds. |
| `HEARTBEAT.md` | Checklist the agent reviews on a cadence — what to check, what signals mean act, what to ignore. |
| `operating-model.json` | Full machine-readable dump of approved entries, grouped by layer. |
| `schedule-recommendations.json` | Derived suggested time blocks, standing slots, and avoid-windows. |

Hand these to an AI agent (as system prompt context, or via a tool-use loop) and the agent can decide on your behalf within the rules you wrote.

## Re-run cadence

**Quarterly** — operating models drift. Role shifts, team changes, new tooling, new dependencies. Pick a re-run mode:

- `update` — revalidate each layer's stored entries in place. Fastest. Use when most of the operating model is still correct.
- `version` — archive the current run and start fresh, with a reference back. Use when the operating model has shifted materially but you want the old one for comparison.
- `fresh` — archive and start completely over. Use after a major role change where the prior run doesn't help.

## Integration with bkb

After export, feed the operating model into the knowledge base:

```bash
do-work interview work-operating-model ingest
do-work bkb triage
do-work bkb ingest
```

After ingest, the operating model is queryable alongside everything else in the KB:

```bash
do-work bkb query "when do I escalate dependency delays?"
do-work bkb query "what's my deep-work window on Tuesdays?"
```

The `ingest` sub-command writes files to `kb/raw/inbox/` with the `topic_cluster` from the template's frontmatter (default `operating-model`), so BKB groups them in the right cluster on ingest.

## Context separation — multiple operating models

One operating model per repo. If you want distinct operating models for different contexts (personal vs. work, or two different roles), install the skill in two repos. The interview action reads and writes relative to the current working directory — `./do-work/interview/<template>/` — so each repo has its own independent session.

## Troubleshooting

**The agent started asking abstract questions — "how would you describe your work?"**
Push back. Say "ask about last week" or "give me a concrete question about a real day." The Interviewer persona is supposed to start concrete; if it drifted, steer it back.

**A checkpoint feels wrong and I already said "save."**
The next layer's checkpoint can still reference or revise the prior layer. Alternatively, run `do-work interview work-operating-model review` — the contradiction pass often catches "wait, that's not quite right" after the fact, and the review UX lets you revise.

**I want to edit an entry in the middle of a layer.**
Just say so. The Interviewer overwrites the draft checkpoint on each revision. Don't approve the checkpoint until it matches what you meant.

**`session.json` is corrupt (invalid JSON).**
Do not try to repair it by hand — the action refuses to load a corrupt file and does not auto-repair. Options: inspect and fix manually (the schema is in `actions/interview-reference.md`), or archive the run with `do-work interview <template> reset --confirm` and start over. If reset itself fails because `session.json` can't be parsed, back up the whole `./do-work/interview/<template>/` directory manually, delete it, and start fresh.

**Review flagged a contradiction but I want to leave both versions.**
Pick `both-are-true` during the review pass. The contradiction gets noted on both entries' `constraints` list so the downstream exports acknowledge the tension rather than hiding it.

**I approved a layer but realized I want to change an entry.**
Two options: run `do-work interview work-operating-model` again and pick `update` mode to revalidate in place, or run `review` — any revision inside a review regenerates the checkpoint and requires re-approval.
