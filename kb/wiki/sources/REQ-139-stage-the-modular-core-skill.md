---
title: "Lessons from REQ-139: Stage the Modular Core Skill"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-139-stage-the-modular-core-skill.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-139: Stage the Modular Core Skill

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Create a self-contained staged `skills/do-work` package while the repository-root all-in-one distribution remains active.

## Solution summary

[MAP CHANGED] The modular suite now has a staged core context boundary at `skills/do-work`: queue orchestration remains feature-complete, while board, knowledge, and toolbox capabilities are named sibling dependencies instead of hidden local files. The repository-root all-in-one skill remains the live distribution until the migration gate.

## What worked

- Copying the feature-rich actions first and changing only their ownership-boundary references preserved behavior while making the split auditable.
- A staged-suite contract can validate not-yet-created siblings safely by requiring their names in the exact suite manifest, then automatically tighten to file existence when each sibling directory appears.

## What didn't work

- A naive ban on every textual `CLAUDE.md` mention incorrectly treated historical changelog text and consumer-project discovery guidance as live citations; the contract needed to target actual links and directive citations.

## Worth knowing

- The temporary pipeline necessarily crosses into toolbox for inspect and present; REQ-145 removes that stateful compatibility path only after activation.
- Core still carries all domain crew files because `work` selects them dynamically; moving those guardrails would weaken the orchestrator even though several toolbox actions also consume them.

## Back-reference

See `do-work/archive/UR-031/REQ-139-stage-modular-core-skill.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9ba534e`.
