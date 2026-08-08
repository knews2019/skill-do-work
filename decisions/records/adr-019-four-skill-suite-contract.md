---
title: "ADR-019: Distribute do-work as One Four-Skill Suite"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-08-07
sources:
  - do-work/user-requests/UR-031/input.md
  - VERSION
  - suite/modules.tsv
  - tools/validate-suite-manifest.sh
related:
  - page: adr-013-harden-the-vendored-skill-distribution-model
    rel: extends
  - page: adr-016-vendor-queue-kanban-into-the-skill
    rel: extends
created: 2026-08-07
updated: 2026-08-07
confidence: high
---

# ADR-019: Distribute do-work as One Four-Skill Suite

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-013-harden-the-vendored-skill-distribution-model]] (extends), [[adr-016-vendor-queue-kanban-into-the-skill]] (extends)

## Context

The current distribution installs one large `do-work` skill. Queue orchestration, the Kanban application, knowledge tooling, and general repository utilities therefore share one routing and context boundary even though most commands use only one of those groups. The repository already ships and versions them together, and users rely on the current feature-rich `do-work run` behavior plus both update entry points.

The split must reduce per-command context without creating independently versioned modules, partial installs, or a cutover that old updaters cannot understand. Existing clients first need a bridge updater that can read both the monolithic archive and the new manifest.

## Decision

Distribute one required suite containing four sibling skills:

- `do-work` owns queue capture, verification, orchestration, review, release, and lifecycle behavior.
- `do-work-board` owns the queue-kanban application and board-facing commands.
- `do-work-knowledge` owns knowledge-base, memory, exploration, interview, and prompt-library commands.
- `do-work-toolbox` owns reporting, feedback validation, repository audits, design review, and companion utilities.

The canonical source and client destinations live in `suite/modules.tsv`. Its grammar is exactly two tab-separated columns, `source` and `destination`, followed by the four mappings under `skills/<name>/` to `.claude/skills/<name>`. `tools/validate-suite-manifest.sh` is the single validator used by update and installation paths. It rejects malformed or incomplete manifests, unsafe or duplicate paths, unexpected modules, symlink escapes, and missing or empty `SKILL.md` files before any managed write.

Root `VERSION` is the suite version. The four modules never carry independent versions and are always installed or updated together. Until the modular cutover, `actions/version.md` remains the active monolithic display and mirrors the same value.

An update is **all-or-recover**, not a filesystem-atomic four-directory rename: every module validates in staging before the first managed write; success is reported only after all installed bytes verify; any later failure restores every changed managed module path and removes only newly created files inside validated destinations. A client must therefore observe either the previous complete suite or the new complete suite after recovery, never a reported partial success.

The live archive stays monolithic for the bridge release. Temporary root-anchored `export-ignore` entries keep `VERSION`, `suite/`, and `skills/` out of that archive while the manifest, validator, staged packages, and installer land in the source repository. REQ-144 removes the three guards only after every known client reports the bridge capability `suite-layout-v2`.

## Alternatives

1. **Keep one all-in-one skill.** Rejected because unrelated command families continue paying the same routing and context cost.
2. **Version or install modules independently.** Rejected because cross-skill references would permit incompatible combinations and multiply update choices for one product.
3. **Switch the archive layout immediately.** Rejected because an updater already running cannot acquire migration logic from the archive it is in the middle of replacing.
4. **Claim filesystem atomicity across four destinations.** Rejected because sibling directory replacements cannot be one portable rename. Staged validation plus complete recovery provides the observable guarantee users need.

## Consequences

Every client receives four skill directories at one version and one confirmation boundary. Both public update commands and the fresh installer can share one manifest parser and one transaction contract. The cost is a deliberate two-release migration: a bridge release must reach every existing client before the live layout changes, and one modular release retains explanatory command shims before their removal.

The validator becomes security-sensitive distribution code. New manifest fields, module names, or destinations require an ADR update and matching fixtures rather than an ad hoc caller-specific parser.

## References

- [VERSION](../../VERSION) — canonical suite version
- [suite/modules.tsv](../../suite/modules.tsv) — exact four-module mapping
- [tools/validate-suite-manifest.sh](../../tools/validate-suite-manifest.sh) — shared validator
- [do-work/user-requests/UR-031/input.md](../../do-work/user-requests/UR-031/input.md) — confirmed rollout and ownership decisions
- [[adr-013-harden-the-vendored-skill-distribution-model]] — existing whole-tree archive and update safeguards
- [[adr-016-vendor-queue-kanban-into-the-skill]] — current board ownership and distribution history

## Stateful Pipeline Retirement

After modular cutover, REQ-145 removed the separate stateful pipeline action, companion reference, `do-work/pipeline.json` lifecycle, and Stop guard. The core router exposes no `pipeline` or `full` alias. The approved successor is the copyable capture → verify → `do-work run` → `do-work-toolbox present-work` prompt in the root README and core help; testing and review remain built into `do-work run`, with no replacement state or separate testing stage.

This retirement supersedes ADR-005 through ADR-008 while preserving them as historical records.
