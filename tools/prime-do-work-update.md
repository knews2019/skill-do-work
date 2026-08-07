# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a do-work skill installed inside the invoking project; `actions/install.md` installs the just recipe, and `actions/version.md` remains the canonical agent-driven update contract.

## Read first

- `tools/do-work-update.sh` — project-root gate, upstream fetch, confirmation, extraction, and audit.
- `actions/install.md` — shipped recipe block, drift upgrade, and installer verification.
- `actions/version.md` — canonical safeguards the script must keep aligned with.

## Do not edit

- `<project-root>/do-work/` — runtime queue, archives, and deliverables; the updater must always exclude it from extraction.
- A consumer's justfile outside the installer’s consent-gated recipe spans.

## Stakes

- `do-work-update.sh` — project-local overwrite boundary
  Req: reject skill roots outside the invoking project, show the reviewed diff, require confirmation, and verify the installed version afterward.
  Value: users can update without an agent turn while retaining the protection against clobbering a shared install or local customization.
  Risk: weakening any guard can overwrite user work or runtime queue data. The script keeps **no** rollback copy — version control is the undo — so a failure inside the destructive region must report the partial install with runnable recovery commands (`print_recovery_instructions`), and runtime state must never need recovery because it is never touched. Do not reintroduce a `cp -R` snapshot; `_dev/tests/contract-regressions.sh` fails the build if you do.
  Because git only restores what was **committed**, the one thing the snapshot did cover is now a warning, not a mechanism: uncommitted edits to shipped files die at the extraction, so the script must name them before the confirmation prompt and repeat in the recovery path that the printed `git checkout` will not bring them back. `_dev/tests/update-script-behavior.sh` Probe 4 holds both messages.

## Lessons

- REQ-061: semantic-version comparison must execute on the platform’s `awk` implementation; avoid names such as `index`, which some implementations reserve as built-ins.
- [REQ-136: staged suite artifacts must remain export-ignored through the bridge release](../do-work/archive/REQ-136-define-four-skill-suite-contract.md#lessons-learned)
