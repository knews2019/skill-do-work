# Decisions

Architecture Decision Records (ADRs) for the do-work skill — a wiki-style log of the load-bearing choices that shape how the skill works, modeled on the `build-knowledge-base` pattern (`_master_index` + topic clusters + typed `related:` links + `log.md` timeline).

An ADR captures a decision that is **hard to reverse without also reversing the behaviors built on top of it**. If you could swap the choice out silently and nothing downstream would notice, it probably doesn't belong here.

## When to Add an ADR

Add a new ADR when:

- A change defines a **new invariant** that other actions will rely on (e.g., "REQs are immutable once claimed").
- A change **closes off an alternative** that was seriously considered (e.g., "bare `do-work/` vs. `do-work/queue/`").
- A change **shapes the skill's surface area** — a new directory, a new file-naming convention, a new routing priority.
- Someone reading the codebase in a year would ask "why is it this way?" and the CHANGELOG entry alone wouldn't be enough.

Do NOT add an ADR for:

- Bug fixes, typos, or doc polish (those live in CHANGELOG only).
- Additive features that reuse an existing pattern (they extend a prior ADR — cite it, don't repeat it).
- Reversible experiments (come back once they've stuck).

## Directory Layout

```
decisions/
  README.md                         # This file
  _master_index.md                  # Nav — all ADRs by topic + by date
  log.md                            # Append-only timeline of decisions
  topics/
    _index_queue-model.md           # UR/REQ structure, queue path, immutability
    _index_platform-portability.md  # Agent-agnostic, subagent dispatch
    _index_routing-dispatch.md      # SKILL.md routing table
    _index_content-structure.md     # Crew members, companion refs, templates
    _index_philosophy.md            # Capture ≠ Execute, validated intent
  adr/
    NNNN-short-slug.md              # One file per decision
```

## ADR Page Schema

Every ADR page has YAML frontmatter followed by five fixed sections. No section is optional.

```yaml
---
id: NNNN                            # 4-digit, monotonically increasing
title: Short descriptive title
status: accepted                    # accepted | proposed | superseded | deprecated
decided: YYYY-MM-DD                 # Date of the release that introduced it
version: X.Y.Z                      # Changelog version where it landed
topic: queue-model                  # One of the topic clusters
supersedes: []                      # List of ADR ids this replaces
superseded_by: null                 # ADR id that replaces this one, or null
related:
  - adr: NNNN
    rel: complements                # extends | supersedes | depends-on | complements | contradicts
---

# ADR-NNNN: Title

## Context
Why the decision needed to be made. What problem was visible. What would break or stay broken if we did nothing.

## Decision
The choice, stated in one to three sentences. No hedging.

## Alternatives Considered
What else was on the table and why each was rejected. At least two alternatives — if there was only one option, this isn't really a decision.

## Consequences
What this enables. What it closes off. What new obligations it creates. Be honest about the costs.

## References
- CHANGELOG entry (version, codename)
- Primary action files the decision lives in
- Related ADRs (cite by id)
```

## Typed Relationships

Match the BKB vocabulary — relationships between ADRs use these verbs:

| Relationship  | Meaning                                                           |
|---------------|-------------------------------------------------------------------|
| `extends`     | Builds on an earlier decision without replacing it                |
| `supersedes`  | Replaces an earlier decision (target marked `status: superseded`) |
| `depends-on`  | Only makes sense because an earlier decision holds                |
| `complements` | Covers adjacent ground without extending or depending             |
| `contradicts` | Explicit tension — both ADRs should name the conflict             |

Relationships are **bidirectional**. If A `extends` B, B's frontmatter lists A as a relation (typically as `complements`, or as `superseded_by` if A replaces B).

## How to Add a Decision

1. Pick the next free id by scanning `adr/` (currently `0001`–`0010`).
2. Create `adr/NNNN-short-slug.md` using the schema above.
3. Add a line to `_master_index.md` under the correct topic cluster **and** the chronological section.
4. Append a one-line entry to `log.md` with the date, id, title, and action (`new` / `superseded` / `amended`).
5. If the new ADR supersedes an older one, update the older ADR's frontmatter (`status: superseded`, `superseded_by: NNNN`) and link it explicitly.
6. Update the correct `topics/_index_*.md` file with the new ADR entry.

## How to Read This Log

- **Looking for the big picture?** Start at `_master_index.md`.
- **Wondering why a specific pattern exists?** Find the topic cluster index (e.g., `topics/_index_queue-model.md`) and follow the links.
- **Want a timeline?** `log.md` lists every decision in the order it was made.
- **Tracing how a choice evolved?** Follow `supersedes` / `superseded_by` chains in the frontmatter.
