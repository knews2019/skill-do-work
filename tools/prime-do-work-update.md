# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a do-work skill installed inside the invoking project; `actions/install.md` installs the just recipe, and `actions/version.md` remains the canonical agent-driven update contract.

## Read first

- `tools/do-work-update.sh` — project-root gate, upstream fetch, confirmation, extraction, and audit.
- `actions/install.md` — shipped recipe block, drift upgrade, and installer verification.
- `actions/version.md` — canonical safeguards the script must keep aligned with.

## Do not edit

- `<project-root>/do-work/` — runtime queue, archives, and deliverables; the updater's explicit managed-path plan must never include it.
- A consumer's justfile outside the installer’s consent-gated recipe spans.

## Stakes

- `do-work-update.sh` — project-local overwrite boundary
  Req: reject skill roots outside the invoking project, show the reviewed diff, require confirmation, and verify the installed version afterward.
  Value: users can update without an agent turn while retaining the protection against clobbering a shared install or local customization.
  Risk: weakening any guard can overwrite user work or runtime queue data. The bridge requires the project Git root, validates the full archive with its already-installed manifest validator, inventories only explicit managed destinations, and automatically restores their committed bytes plus removes only newly created managed paths after failure. Runtime and application paths must never enter that plan. Dirty managed changes are named before the one confirmation; accepting discards them from both index and worktree before installation. `_dev/tests/update-script-behavior.sh` holds legacy, suite, hostile-manifest, dirty-consent, and forced-recovery behavior.

## Lessons

- REQ-061: semantic-version comparison must execute on the platform’s `awk` implementation; avoid names such as `index`, which some implementations reserve as built-ins.
- [REQ-136: staged suite artifacts must remain export-ignored through the bridge release](../do-work/archive/REQ-136-define-four-skill-suite-contract.md#lessons-learned)
- [REQ-137: the installed bridge validator must remain authoritative over a future archive](../do-work/archive/REQ-137-ship-suite-aware-bridge-updater.md#lessons-learned)
- [REQ-138: ambiguous unmarked legacy recipe spans must fail instead of absorbing client content](../do-work/archive/REQ-138-add-managed-text-section-replacement.md#lessons-learned)
- [REQ-139: staged sibling references must be manifest-declared before their packages exist](../do-work/archive/REQ-139-stage-modular-core-skill.md#lessons-learned)
- [REQ-140: copy staged modules from committed Git state when the active tree has unrelated edits](../do-work/archive/REQ-140-stage-modular-board-skill.md#lessons-learned)
