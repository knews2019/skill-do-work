---
id: REQ-143
title: "Build the Full-Suite Installer and Reconciler"
status: completed
claimed_at: 2026-08-07T22:04:11Z
completed_at: 2026-08-07T22:14:02Z
commit: 8f22cbe
route: C
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md, tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-138, REQ-139, REQ-140, REQ-141, REQ-142]
maintenance: false
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
write_set: [tools/install-do-work-suite.sh, suite/modules.tsv, README.md, _dev/tests/contract-regressions.sh, _dev/tests/install-suite-behavior.sh]
kb_status: promoted
kb_entry: REQ-143-build-the-full-suite-installer-and-recon.md
---

# Build the Full-Suite Installer and Reconciler

## What
Build the canonical fresh-install bootstrap and client-configuration reconciler for the complete four-skill suite, without activating the live modular distribution yet.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add strict RED hermetic installer/bootstrap fixtures first; implement one archive-backed full-suite transaction with manifest/source/destination validation, exact managed-path backups, one consent boundary, board-template reconciliation through the managed-section utility, and jq/Python/manual hook composition; verify installed bytes and configuration before declaring success.
- [x] **[APPLY]:** Added the canonical one-archive bootstrap and full-suite installer; implemented four-module validation/install/byte verification, managed Just reconciliation, core-hook composition with exact memory-path migration, one confirmation, and signal/failure recovery from temporary originals.
- [x] **[UNIFY]:** Reviewed the installer, its 456-line hermetic behavior suite, the contract invocation, and release/request metadata; verified destructive targets, path containment, exact backups, signal traps, jq/Python/manual branches, source/destination validation, bootstrap artifact count, mode preservation, ShellCheck, Bash syntax, full contracts, Just parsing, Go checks, and diff cleanliness.

## Why
Clients need one reliable install command and automated migration of the two project-owned configuration surfaces affected by the split.

## Detailed Requirements
- Add `tools/install-do-work-suite.sh` requiring `--project-root` and a Git repository.
- Define one canonical, copy-paste bootstrap command that requires no preinstalled do-work skill, downloads one upstream artifact, and invokes the installer for the target project. Test the exact command against a clean hermetic Git project; REQ-144 publishes it as the live README installation path at cutover.
- Download one upstream archive and validate `VERSION`, `suite/modules.tsv`, and all four non-empty `SKILL.md` files in staging before writing.
- Always install the full suite; never install a subset.
- Create a complete Justfile from the board template when none exists, or migrate/replace only the managed `do-work:recipes` section in an existing file.
- Enable core Claude hooks by composing with valid existing settings.
- Do not enable memory hooks on fresh installs.
- If known legacy memory-hook command strings already exist, rewrite only those strings to `do-work-knowledge`.
- Use `jq` when available, Python 3 as fallback, and otherwise leave settings unchanged with an exact manual instruction.
- Restore exact temporary originals when Just/settings validation fails.
- Cover fresh install, reinstall, legacy recipes, custom recipes, invalid markers/settings, hook composition, memory-hook migration, spaces, interruption, and exact four-module verification.
- Keep the live repository in bridge mode until REQ-144.

## Constraints
- All four modules share one version and one confirmation boundary.
- Installer-owned reconciliation must not overwrite unrelated client configuration.

## Dependencies
Requires the managed-section utility and all four staged skill packages.

## Builder Guidance
Certainty level: Firm. Reuse the manifest validator and managed-section utility instead of adding parallel implementations.

## Red-Green Proof
**RED prompt/case:** Run the exact canonical bootstrap command against an empty Git project, then run the installer against a project with custom Just recipes and existing core/memory hooks.
**Why RED now:** There is no full-suite installer, no four-module verification, and no deterministic configuration reconciliation.
**GREEN when:** The documented copy-paste command installs the complete suite without a preinstalled skill; both fixtures receive the same-version four-module suite; custom content survives, the managed block is current, core hooks are enabled, and memory hooks are only migrated when already present.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for bootstrap, Justfile, hook, and Git policies.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

## Triage

**Route: C** — this is a multi-surface install transaction spanning archive trust, four managed module trees, Just ownership, JSON hook composition, confirmation, rollback, and a bootstrap that must work without an installed skill.

## Scope

**Files I will touch:** `tools/install-do-work-suite.sh`, `_dev/tests/contract-regressions.sh`, `do-work/working/REQ-143-build-full-suite-installer-reconciler.md`, and release metadata. The already-verified bootstrap requirement in this REQ remains part of the committed request record.

**Files I will not touch:** the live export layout, active root router/actions, client runtime `do-work/` or `kb/` data, unrelated application/configuration bytes, or unrelated dirty REQ-134/REQ-147 work. README publication waits for REQ-144.

## Pre-Flight

The full contract baseline passes before installer work. The strict installer fixture is RED because `tools/install-do-work-suite.sh` does not exist; no canonical no-preinstall bootstrap or four-module configuration transaction is currently executable.

## Implementation Summary

**Files changed:**
- `tools/install-do-work-suite.sh` (new) — canonical bootstrap printer plus recoverable four-module installer/configuration reconciler
- `_dev/tests/install-suite-behavior.sh` (new) — hermetic fresh, reinstall, legacy, custom, malformed, fallback, rollback, interruption, spaces, cancellation, and Git-root fixtures
- `_dev/tests/contract-regressions.sh` (modified) — runs the suite-installer behavior contract from the repository acceptance suite

The installer downloads or reuses exactly one archive, validates the shared version and exact manifest before client writes, rejects unsafe sources/destinations, installs all four packages behind one confirmation, reconciles the first existing Justfile variant through `replace-text-section.sh`, composes core hooks with jq or Python, migrates only known memory command paths, leaves settings exact with a precise manual step when neither JSON tool exists, verifies every installed byte, and restores exact managed originals on failure or interruption. README publication and live export activation remain deferred to REQ-144.

## Qualification

Passed — all requirement paths are substantive and connected: the printed bootstrap executes the same installer with its single downloaded archive; the validated manifest drives exactly four destination writes; the board template and managed-section utility drive Just ownership; hook fragments flow through jq/Python composition or the explicit unchanged/manual branch; and every write participates in the snapshot, post-write verification, and recovery transaction.

## Testing

**Tests run:** strict RED/GREEN installer behavior suite; exact canonical bootstrap in a clean Git repo with spaces; fresh/reinstall/idempotence; custom and legacy Justfiles; invalid markers/settings; corrupt module validation; jq/Python/manual settings modes; exact Just/settings rollback; TERM interruption recovery; cancellation; non-Git/subdirectory refusal; full contract regressions; warning-level ShellCheck; Bash syntax; Just parse; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Installer/bootstrap: ✗ executable absent → ✓ one copy-paste command fetches one archive and installs four verified modules without a preinstalled skill
- Configuration: ✗ no suite-owned reconciliation → ✓ complete/managed/legacy Just cases and core/memory hook policies are deterministic
- Recovery: ✗ no cross-surface transaction → ✓ forced Just failure, forced settings failure, and TERM restore exact prior modules/configuration

**New tests added:** `_dev/tests/install-suite-behavior.sh`, invoked by the main contract suite.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-07T22:11:47Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 0
**Acceptance:** Pass — the staged suite now has one hermetically proven fresh-install path and a configuration-safe, all-or-recover client transaction, without activating the modular export early.
**Suggested testing:** REQ-144 must run this bootstrap against the cutover archive after removing export-ignore guards and compare both updater entry points byte-for-byte after bridge migration.
**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Passing the already-downloaded bootstrap archive into the installer makes “one artifact” mechanically testable and avoids a second network trust boundary.
- Preparing Just and settings candidates before confirmation moves malformed-input failures ahead of all managed writes while keeping post-write recovery independently testable.

**What didn't:**
- An initial bootstrap assertion accidentally matched only the generic `--` token; tightening it to the literal archive argument prevented a false-positive contract.
- Resetting signal traps to defaults before recovery let a second signal interrupt rollback; recovery must ignore termination signals until originals are restored.

**Worth knowing:**
- The managed-section utility requires Python for existing Justfiles. A Python-free fresh project can still copy the complete board template and take the explicit manual settings step; an existing Justfile fails unchanged rather than receiving an unsafe shell rewrite.
- Hook migration is limited to `command` string values containing the two known legacy memory paths. Other settings strings and every unrelated hook entry survive.

## Orientation

[MAP CHANGED] `tools/install-do-work-suite.sh` is the staged suite’s canonical bootstrap/reconciliation boundary. It owns archive validation, full four-module installation, managed Just recipes, core-hook composition, known memory-path migration, confirmation, byte verification, and recovery; live distribution still remains in bridge mode until REQ-144.
