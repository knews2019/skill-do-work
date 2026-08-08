# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a do-work suite installed inside the invoking project; the bridge's already-installed full-suite installer owns modular module/configuration reconciliation, and `actions/version.md` remains the canonical agent-driven update contract.

## Read first

- `tools/do-work-update.sh` — project-root gate, upstream fetch, confirmation, extraction, and audit.
- `tools/install-do-work-suite.sh` — trusted installed module, Just section, hook composition, verification, and exact recovery transaction.
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
- [REQ-144: bridge and fresh installs must converge on one verified four-skill recovery transaction](../../../do-work/archive/REQ-144-activate-four-skill-distribution.md#lessons-learned)
