# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a four-skill do-work suite installed inside the invoking project; the already-installed full-suite installer owns module/configuration reconciliation, and `actions/version.md` remains the canonical agent-driven update contract.

## Read first

- `tools/do-work-update.sh` — project-root gate, upstream fetch, confirmation, extraction, and audit.
- `tools/install-do-work-suite.sh` — trusted installed module, managed-configuration (Just section, agent-instructions section, hook composition), verification, and exact recovery transaction. The installer is authoritative for the current surface set; any list here is illustrative and must not be read as closed.
- `tools/replace-text-section.sh` — byte-preserving managed-section replacement for whichever marker pair the caller names (`--begin-marker`/`--end-marker`, defaulting to the Just recipe markers), plus no-Just reserved-recipe collision validation.
- `actions/version.md` — canonical safeguards the script must keep aligned with.

## Do not edit

- `<project-root>/do-work/` — runtime queue, archives, and deliverables; the updater's explicit managed-path plan must never include it.
- Any consumer-owned bytes outside the installer’s consent-gated managed marker spans — the condition is the rule, not the file. That covers a justfile outside the recipe markers and a `CLAUDE.md` outside the communication-style markers alike, and it covers whatever surface the installer manages next.

## Stakes

- `do-work-update.sh` — project-local overwrite boundary
  Req: reject a skill root outside the invoking project, require `--project-root` to name the Git worktree root, refuse an upstream version that is not newer, delegate the reviewed diff and the single confirmation to the installed full-suite installer (honouring its cancel status as a no-change exit), and verify the installed version afterward.
  Value: users can update without an agent turn while retaining the protection against clobbering a shared install or local customization.
  Risk: weakening any guard can overwrite user work or runtime queue data. The updater requires the project Git root and validates the suite archive with its already-installed manifest validator before delegating that same archive to the installed full-suite installer. The installer reviews modules plus every owned configuration change, snapshots exact managed originals, and restores them on failure — a new managed surface earns a diff section, a backup, a recovery branch, and a post-write byte check in the same commit that adds it. Runtime, KB, application paths, unrelated settings, and any bytes outside a managed marker span must never enter that plan. Dirty module changes are named before the one confirmation; accepting discards them from both index and worktree before installation. `_dev/tests/update-script-behavior.sh` holds current-suite, hostile-manifest, dirty-consent, and forced-recovery behavior.

## Lessons

See [`lessons-do-work-update.md`](lessons-do-work-update.md) — read it before changing what **Read first** or **Traps** name above.
