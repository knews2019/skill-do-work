# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a do-work installation inside the invoking project; `actions/install.md` installs the bridge-era recipe, the already-installed `tools/install-do-work-suite.sh` owns modular module/configuration reconciliation, and `actions/version.md` remains the canonical agent-driven update contract.

## Read first

- `tools/do-work-update.sh` — project-root gate, upstream fetch, confirmation, extraction, and audit.
- `tools/install-do-work-suite.sh` — suite module, Just section, hook composition, verification, and exact recovery transaction.
- `actions/install.md` — shipped recipe block, drift upgrade, and installer verification.
- `actions/version.md` — canonical safeguards the script must keep aligned with.

## Do not edit

- `<project-root>/do-work/` — runtime queue, archives, and deliverables; the updater's explicit managed-path plan must never include it.
- A consumer's justfile outside the installer’s consent-gated recipe spans.

## Stakes

- `do-work-update.sh` — project-local overwrite boundary
  Req: reject skill roots outside the invoking project, show the reviewed diff, require confirmation, and verify the installed version afterward.
  Value: users can update without an agent turn while retaining the protection against clobbering a shared install or local customization.
  Risk: weakening any guard can overwrite user work or runtime queue data. The bridge requires the project Git root and validates a suite archive with its already-installed manifest validator before delegating that same archive to the full-suite installer. The installer reviews modules plus the owned Just/settings changes, snapshots exact managed originals, and restores them on failure. Runtime, KB, application paths, exterior Just bytes, and unrelated settings must never enter that plan. Dirty module changes are named before the one confirmation; accepting discards them from both index and worktree before installation. `_dev/tests/update-script-behavior.sh` holds legacy, suite, configuration migration, hostile-manifest, dirty-consent, and forced-recovery behavior.

## Lessons

- REQ-061: semantic-version comparison must execute on the platform’s `awk` implementation; avoid names such as `index`, which some implementations reserve as built-ins.
- [REQ-136: staged suite artifacts must remain export-ignored through the bridge release](../do-work/archive/REQ-136-define-four-skill-suite-contract.md#lessons-learned)
- [REQ-137: the installed bridge validator must remain authoritative over a future archive](../do-work/archive/REQ-137-ship-suite-aware-bridge-updater.md#lessons-learned)
- [REQ-138: ambiguous unmarked legacy recipe spans must fail instead of absorbing client content](../do-work/archive/REQ-138-add-managed-text-section-replacement.md#lessons-learned)
- [REQ-139: staged sibling references must be manifest-declared before their packages exist](../do-work/archive/REQ-139-stage-modular-core-skill.md#lessons-learned)
- [REQ-140: copy staged modules from committed Git state when the active tree has unrelated edits](../do-work/archive/REQ-140-stage-modular-board-skill.md#lessons-learned)
