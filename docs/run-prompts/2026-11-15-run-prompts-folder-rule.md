# Capture prompt: decide where run prompts live in this repository

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.

## The request

Settle one conflict between a maintainer rule and a repository check, with the smallest change.

### Context

- The maintainer's global rule says: write a prompt for another session to `docs/run-prompts/<date>-<slug>.md` and commit it.
- In this repository that file breaks a check. `_dev/tests/staged-skills-contract.sh` (the loops near lines 75 and 99 to 114, message "legacy root runtime must be retired at modular cutover") allows only `docs/kb/`, `docs/decisions/` and `docs/changelog/` under the root `docs` folder. On 2026-11-14 the file `docs/run-prompts/2026-11-14-tinka-canvas-suggestion.md` turned the heavy lane `staged-skills` red five commits after it was added, and was removed in commit 00568685. The fast gate did not catch it.
- Since then run prompts for this repository are kept in `_dev/run-prompts/`, which is maintainer-only.

### The choice (park it as a pending answer; recommended first)

1. **Allow `docs/run-prompts/` in the contract** (recommended). It is a maintainer folder like the three allowed ones, and the global rule then works here as in every other project. Add it to both lists in the contract, check that it is export-ignored like the other maintainer docs, and move the files from `_dev/run-prompts/` there.
2. **Keep `_dev/run-prompts/`.** Change nothing in the contract. Write one line in the maintainer `CLAUDE.md` that names this folder as the home of run prompts in this repository.

### Files to change

- Choice 1: `_dev/tests/staged-skills-contract.sh`, the export-ignore rules if needed, and the moved prompt files.
- Choice 2: `CLAUDE.md` only.

### What must not change

- Nothing under `skills/` or `suite/`, so this is no release.
- The contract keeps refusing every other new file under the root `docs` folder.

### Verification

RED first (choice 1): with a file `docs/run-prompts/probe.md` present, `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` fails with the legacy-root message.

GREEN (choice 1): the same lane passes with that file present and still fails with a file `docs/other/probe.md`; the probe files are removed afterwards. GREEN (choice 2): `CLAUDE.md` names the folder and the lane passes.
