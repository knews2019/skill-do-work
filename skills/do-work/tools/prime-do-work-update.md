# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a four-skill do-work suite installed inside the invoking project; the already-installed full-suite installer owns module/configuration reconciliation, and `actions/version.md` remains the canonical agent-driven update contract.

## Read first

- `tools/do-work-update.sh` — project-root gate, upstream fetch, confirmation, extraction, and audit.
- `tools/install-do-work-suite.sh` — trusted installed module, Just section, hook composition, verification, and exact recovery transaction.
- `tools/replace-text-section.sh` — byte-preserving managed-section replacement plus no-Just reserved-recipe collision validation.
- `actions/version.md` — canonical safeguards the script must keep aligned with.

## Do not edit

- `<project-root>/do-work/` — runtime queue, archives, and deliverables; the updater's explicit managed-path plan must never include it.
- A consumer's justfile outside the installer’s consent-gated recipe spans.

## Stakes

- `do-work-update.sh` — project-local overwrite boundary
  Req: reject skill roots outside the invoking project, show the reviewed diff, require confirmation, and verify the installed version afterward.
  Value: users can update without an agent turn while retaining the protection against clobbering a shared install or local customization.
  Risk: weakening any guard can overwrite user work or runtime queue data. The updater requires the project Git root and validates the suite archive with its already-installed manifest validator before delegating that same archive to the installed full-suite installer. The installer reviews modules plus the owned Just/settings changes, snapshots exact managed originals, and restores them on failure. Runtime, KB, application paths, exterior Just bytes, and unrelated settings must never enter that plan. Dirty module changes are named before the one confirmation; accepting discards them from both index and worktree before installation. `_dev/tests/update-script-behavior.sh` holds current-suite, hostile-manifest, dirty-consent, and forced-recovery behavior.

## Lessons

- REQ-061: semantic-version comparison must execute on the platform’s `awk` implementation; avoid names such as `index`, which some implementations reserve as built-ins.
- [REQ-136: the suite manifest is the sole module source/destination contract](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-136-define-four-skill-suite-contract.md#lessons-learned)
- [REQ-137: the installed manifest validator must authorize every candidate archive](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-137-ship-suite-aware-bridge-updater.md#lessons-learned)
- [REQ-138: managed recipe markers preserve exterior client bytes and reject malformed ownership](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-138-add-managed-text-section-replacement.md#lessons-learned)
- [REQ-144: every install must use one verified four-skill recovery transaction](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-144-activate-four-skill-distribution.md#lessons-learned)
- [REQ-146: modular updates retain one installed all-or-recover transaction](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-146-remove-modular-migration-shims.md#lessons-learned)
- [REQ-162: multiline literal state must carry recipe-header state too](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-162-handle-ordinary-multiline-backtick-commands.md#lessons-learned)
- [REQ-163: Markdown reference classification must honor syntax adjacency and container scope](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-031/REQ-163-complete-remaining-inline-link-and-list-fence-classification.md#lessons-learned)
