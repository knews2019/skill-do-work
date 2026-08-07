---
id: REQ-144
title: "Activate the Four-Skill Distribution"
status: claimed
claimed_at: 2026-08-07T22:21:03Z
status_changed_at: 2026-08-07T21:40:04Z
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md, tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-137, REQ-143]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Activate the Four-Skill Distribution

## What
Switch the live distribution to the four staged skills and migrate bridge-enabled client installations through one all-or-recover suite transaction.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** First close the newly discovered bridge seam by making the bridge updater delegate a future suite archive to the tested full-suite installer, with configuration diffs and one confirmation. Ship that while export guards remain, then require renewed client rollout evidence before deleting the monolith and activating the archive. After the gate, add exact moved-command shims, publish the canonical bootstrap, switch maintainer/runtime paths, remove shipped legacy duplication, and prove both update entry points converge.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Staging alone does not reduce the installed router or context; the manifest and live archive must cut over only after clients can understand the new layout.

## Detailed Requirements
- Activate the live four-module manifest.
- Require user attestation that every known client repository's installed updater reports exactly `suite-layout-v2` from `--capabilities` before changing the live layout; no stored client-by-client inventory is required.
- Make `skills/do-work`, `skills/do-work-board`, `skills/do-work-knowledge`, and `skills/do-work-toolbox` the only shipped runtime sources.
- Remove duplicated legacy root actions/tools in the same change.
- Update export rules, README installation, root maintainer Just fallback, changelog/version handling, and every runtime cross-skill reference. Publish the exact canonical, tested, copy-paste bootstrap command from REQ-143 as the single fresh-install path in the README.
- Add one-release core shims for moved commands. Each shim prints the exact new invocation and stops; it must not forward permanently.
- Ensure `do-work update` and `just run-do-work-update` produce byte-equivalent installations from a bridge client.
- Enforce the all-or-recover suite contract: validate all four modules before the first managed write, report success only after all installed bytes verify, and restore every changed managed module path on failure so clients never observe a partially successful update. Do not claim a cross-directory filesystem-atomic rename.
- Refresh the managed Just section and migrate known memory-hook paths during update.
- Preserve project queue data, KB data, application files, and unrelated configuration.
- Run hermetic bridge-to-modular and fresh-install tests plus ShellCheck, formatting, Go tests/vet, contract tests, and `git diff --check`.

## Constraints
- Do not remove the live export guards until the user attests that every existing client has configuration-aware bridge v0.183.20 or later. The earlier `suite-layout-v2` capability alone predates managed Just/settings migration and is not sufficient evidence.
- Preserve the current feature-rich work orchestrator.
- Keep compatibility shims for exactly one modular release.

## Dependencies
Requires REQ-137 and REQ-143. The original capability gate was cleared by user attestation via `do-work clarify` on 2026-08-08; implementation then proved that a final configuration-aware bridge release must reach those clients before live cutover.

## Builder Guidance
Certainty level: Firm. Treat this as the release cutover; no partial module activation is acceptable.

## Red-Green Proof
**RED prompt/case:** Update a hermetic legacy client that has the bridge updater, legacy Just recipes, queue data, KB data, and an enabled legacy memory hook.
**Why RED now:** The live manifest and installed layout remain monolithic.
**GREEN when:** The user has attested that every known client reports the bridge capability; the README's exact bootstrap command succeeds from a clean Git client; one update installs exactly four current skill directories or restores the prior suite on failure; managed configuration migrates, project data survives, and both update entry points yield identical bytes.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the bridge rollout gate and cutover contract.

## Blocked

- [2026-08-07] blocked on "user supplies an inventory of every known client repository confirming that each installed updater's --capabilities output is exactly suite-layout-v2" — cleared by user via clarify
- [2026-08-08] live cutover is blocked until v0.183.20 is published and the user confirms every known client has updated to that version or later through `do-work update` or `just run-do-work-update`

## Decisions

- [2026-08-08] The user confirmed every known client updater reports exactly `suite-layout-v2` and chose their attestation as sufficient evidence. A stored client-by-client inventory is not required.
- [2026-08-08] A strict bridge-to-suite fixture showed the original capability did not migrate managed Just recipes or legacy memory-hook paths. Ship one final monolithic bridge before cutover; use its installed validator and installer as the trusted engine rather than executing an archive-provided replacement.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

## Triage

**Route: C** — live archive cutover, bridge migration, runtime-source removal, command compatibility, updater/install transactions, configuration reconciliation, and release publication all meet at this boundary.

## Scope

**Files I will touch before the renewed rollout gate:** bridge/full-suite installer engines, their hermetic tests, this REQ/checkpoint, and bridge release metadata. Live export guards and legacy runtime sources remain unchanged in this preparatory release.

**Files I will touch after the gate:** staged core shims/router/help, export rules, README, root maintainer Justfile, update/install contracts, and legacy runtime export/removal surfaces required by cutover.

**Files I will not touch:** project queue/KB/application data, unrelated client configuration, or unrelated dirty REQ-134/REQ-147 implementation work.

## Pre-Flight

The existing bridge updater passes its original suite tests but deliberately never mutates Justfiles or settings. A bridge-to-suite fixture that requires managed recipes and legacy memory-hook migration is therefore RED. Capability output alone does not prove this later configuration-aware bridge behavior; the preparatory release must reach clients before activation.

## Bridge Preparation

Strict RED first proved that the prior bridge left Testing-era board paths and legacy memory hooks unchanged during a suite update. The bridge now delegates the already-downloaded archive to its already-installed full-suite transaction engine, which reviews modules and owned configuration together, asks once, verifies the result, and recovers exact managed file bytes on failure. A hostile-archive fixture also proves that a downloaded installer cannot replace the installed trust boundary. The staged core now carries the installer so the first modular update does not remove the engine needed by later updates.

Version v0.183.20 is intentionally still a monolithic bridge release: `.gitattributes` continues to export-ignore `VERSION`, `suite/`, and `skills/`. Removing those guards, publishing the modular README bootstrap, and deleting legacy shipped duplication remain pending until renewed client rollout evidence is supplied.

**Bridge qualification:** `update-script-behavior`, `install-suite-behavior`, `staged-skills-contract`, the full contract regressions, warning-level ShellCheck, Bash syntax, Just parsing, `gofmt` cleanliness, queue-kanban `go test`/`go vet`/`go build`, synchronized runtime copies, and `git diff --check` all pass for v0.183.20.
