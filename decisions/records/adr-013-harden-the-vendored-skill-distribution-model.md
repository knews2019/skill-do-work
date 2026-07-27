---
title: "ADR-013: Harden the Vendored-Skill Distribution Model"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-06-15
sources:
  - hooks/hooks.json
  - hooks/session-start.sh
  - hooks/pipeline-guard.sh
  - actions/version.md
  - .gitattributes
  - README.md
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: complements
created: 2026-06-15
updated: 2026-06-15
confidence: high
---

# ADR-013: Harden the Vendored-Skill Distribution Model

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-001-modular-action-prompts-and-companion-references]] (complements)

## Context

do-work is distributed as a _vendored_ skill: consumers `curl | tar` the `main` tarball into `.claude/skills/do-work/` and commit it, and `do-work version update` re-extracts a newer tarball over that copy. That model has three sharp edges, each of which caused a real defect:

1. **Hook sample paths were project-relative.** `hooks/hooks.json` (and the two script header comments) shipped `bash hooks/session-start.sh`. Claude Code runs hook handlers from the _project root_, not the skill directory, so the path resolved to `<project-root>/hooks/...` — which doesn't exist — and the SessionStart status hook and Stop pipeline guard failed with "No such file or directory."

2. **The update silently clobbered committed customizations.** The pre-update dirty check ran `git status --porcelain`, which only sees _uncommitted_ edits. A consumer who _committed_ a local fix (e.g. re-anchoring the hook paths above) had a clean tree, passed the check, and had the customization silently overwritten by the tar extraction. This is exactly how the hook-path fix kept getting reverted: downstream consumers re-applied and committed it, every `version update` re-reverted it, and the loop repeated across multiple releases.

3. **The tarball shipped maintainer-internal files.** `decisions/` (this ADR set), the dev dotfiles (`.gitignore`/`.gitattributes`), and `_dev/` (flag-excluded only, never `export-ignore`d — so a GitHub "Download ZIP" would leak it) all landed in consumers' installs as clutter.

## Decision

Harden the distribution model on all three fronts.

1. **Anchor hook sample paths** to `${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/...` in `hooks/hooks.json` and both script header comments. `$CLAUDE_PROJECT_DIR` is the project root Claude Code runs hooks from; `.claude/skills/do-work/` is the canonical install path. Each script header carries a "do NOT simplify back to a relative path — it has regressed before" guard, and `README.md` documents the install-location assumption (adjust the path for non-canonical installs; `${CLAUDE_PLUGIN_ROOT}` is the more robust anchor if do-work is ever distributed as a Claude Code plugin).

2. **Make the update non-clobbering.** `actions/version.md` now (a) also detects _committed_ customizations by diffing the shipped paths from the last version-bump commit to `HEAD`, (b) snapshots non-git installs before extraction so an overwrite is recoverable, and (c) runs a post-update audit of the diff so a reverted customization is caught after the fact. The git-clean precondition makes that diff a reliable audit surface and `git restore` the rollback.

3. **Ship only consumer-relevant files.** `.gitattributes` `export-ignore`s `decisions/`, `_dev/`, `.gitignore`, and `.gitattributes`. `decisions/` is the maintainer's own design history and is never loaded at runtime by any action, so consumers don't need it.

## Consequences

- The bundled hooks work on a default install with no manual edit, and the recurring downstream re-clobber loop is closed at its source — consumers no longer need a local divergence to fix the paths.
- A committed local customization can no longer be silently lost on update: it is detected before, recoverable during, and audited after. The detection is best-effort (a sync that hand-merges a customization into its own commit can still hide it), which is why the post-update audit — not the pre-check — is the real safety net.
- Consumer installs are leaner (107 vs 134 files) and free of maintainer plumbing. Existing installs that already vendored `decisions/` keep their copy (tar never deletes) but stop receiving updates to it.
- The hooks path hardcodes `.claude/skills/do-work/`; non-canonical installs must adjust it. Accepted as the pragmatic cost of a vendored (non-plugin) distribution; revisit if do-work ships as a plugin.
- This ADR lives in `decisions/`, which no longer ships — it is a maintainer-facing record. The enforceable guards for decision 1 live where an editor will actually see them: the script headers and the README.
