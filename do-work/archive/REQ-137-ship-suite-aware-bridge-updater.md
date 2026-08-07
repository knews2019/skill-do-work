---
id: REQ-137
title: "Ship the Suite-Aware Bridge Updater"
status: completed
completed_at: 2026-08-07T21:03:11Z
commit: bef7334
claimed_at: 2026-08-07T20:46:54Z
route: C
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md]
tdd: true
kb_status: pending
kb_entry:
suggested_spec: refactor
depends_on: [REQ-136]
maintenance: false
related: [REQ-135, REQ-136, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
write_set: [actions/version.md, tools/do-work-update.sh, _dev/tests/update-script-behavior.sh, _dev/tests/contract-regressions.sh]
---

# Ship the Suite-Aware Bridge Updater

## What
Ship a bridge release whose updater understands both the current all-in-one archive and the future four-skill suite, while leaving the live layout unchanged.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Replace the updater's duplicated monolithic mutation flow with one tested engine that detects legacy versus suite archives, validates all managed paths before writing, records the exact recoverable write set, and restores only that set after failure; make the agent-facing action delegate to the same engine.
- [x] **[APPLY]:** Replaced the updater with the planned dual-layout transaction engine, delegated the agent update action, and replaced the obsolete manual-recovery probes with strict RED→GREEN coverage in exactly the four declared implementation files.
- [x] **[UNIFY]:** Reviewed the full diff for `actions/version.md`, `tools/do-work-update.sh`, `_dev/tests/update-script-behavior.sh`, and this REQ's isolated hunks in `_dev/tests/contract-regressions.sh`; warning ShellCheck, Bash syntax, all standalone/contract shell suites, queue-kanban test/vet/build, queue verification, and diff checks pass with no debug artifacts.

## Why
An updater already executing cannot adopt replacement migration logic from the archive it is currently extracting. Existing clients therefore need a suite-aware bridge before the source layout changes.

## Detailed Requirements
- Make `actions/version.md` delegate update mutation to the same script used by `just run-do-work-update`.
- Support legacy archives and future `suite/modules.tsv` archives.
- Add `--capabilities`, which prints `suite-layout-v2` without modifying anything.
- Download and review one archive, validate every declared module before writing, show one suite-wide diff, ask once, and verify installed bytes against the reviewed archive.
- Implement the suite contract as one all-or-recover transaction: validate all declared modules before the first managed write, succeed only after all installed bytes verify, and recover every changed managed module path on failure. Do not claim a cross-directory filesystem-atomic rename.
- If managed skill paths are dirty, list them and warn that continuing discards those changes before accepting confirmation.
- On failure after writing, restore every changed tracked managed skill path from Git and clean only newly created files within validated module destinations, leaving no partially updated suite.
- Never restore, reset, clean, or delete application paths, `do-work/` runtime data, `kb/`, Justfiles, or settings files.
- Preserve both current update entry points through the bridge release.
- Add hermetic coverage for legacy updates, future manifests, malformed manifests, traversal, dirty confirmation, partial failure, and runtime-data preservation.

## Constraints
- The live repository remains an all-in-one skill in this REQ.
- The bridge must be released and installed in every client before REQ-144 can start.

## Dependencies
Requires REQ-136's manifest and version contract.

## Builder Guidance
Certainty level: Firm. The bridge is a compatibility release, not the modular cutover.

## Red-Green Proof
**RED prompt/case:** Run the current updater against a fixture archive with a valid four-module manifest.
**Why RED now:** The updater assumes one skill root and cannot install or recover multiple sibling skill destinations.
**GREEN when:** The bridge updater handles both fixture shapes, reports `suite-layout-v2`, and both public update paths exercise the same tested engine.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for rollout and Git-recovery decisions.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This changes the update transaction and recovery model across legacy and future archive layouts, with security-sensitive path validation and two public entry points.

**Planning:** Required

## Plan

1. Rewrite the hermetic updater probes first so legacy success, suite success, capability discovery, malformed/traversing manifests, dirty confirmation, runtime preservation, and mid-transaction recovery all fail against the current engine.
2. Refactor `tools/do-work-update.sh` around one validated managed-module plan, one archive/diff/prompt cycle, and a narrowly scoped Git recovery ledger for both archive shapes.
3. Replace the duplicated agent update procedure in `actions/version.md` with delegation to that script, then run focused, contract, lint, syntax, and Go validation before review.

**Plan validation:** Every detailed requirement maps to one of the three tasks, with the destructive transaction tested before implementation and both entry points sharing the same code path.

*Generated by the primary agent because subagent dispatch is unavailable for this run*

## Exploration

The current updater accepts only `--project-root`, assumes a single skill root, and merely prints manual Git recovery commands after a partial extraction. The agent-facing `actions/version.md` separately prescribes its own fetch/diff/extract procedure. REQ-136's validator and temporary export guards allow the future archive shape to be tested without changing the live monolithic distribution.

*Generated by the primary agent*

## Scope

**Files I will touch:**
- `tools/do-work-update.sh` — implement capability discovery, dual-layout planning, validation, one confirmation, byte verification, and scoped automatic recovery
- `_dev/tests/update-script-behavior.sh` — replace the old no-rollback probes with hermetic legacy/suite/security/recovery coverage
- `_dev/tests/contract-regressions.sh` — adjust only the updater contract invocation/assertions if required
- `actions/version.md` — delegate update mutation to the updater engine used by the Just recipe

**Files I will NOT touch:** The live skill layout, suite export guards, application files, runtime `do-work/`/`kb/`, Justfiles, settings, or unrelated REQ-134 changes.

**Acceptance criteria (restated from REQ):**
- [x] `--capabilities` prints only `suite-layout-v2` and writes nothing.
- [x] Legacy and valid four-module archives use one reviewed archive, one suite-wide diff, one prompt, and verified installed bytes.
- [x] Every suite module and destination validates before the first managed write; malformed and escaping manifests fail without mutation.
- [x] Any failure after writing restores all changed tracked managed paths and removes only files newly created within validated destinations.
- [x] Dirty managed files are listed with an explicit discard warning before confirmation.
- [x] Application paths, runtime data, Justfiles, and settings remain untouched.
- [x] `do-work update` and `just run-do-work-update` delegate to the same tested script.

## Pre-Flight

**Git:** ⚠ Unrelated REQ-134 implementation changes and later UR-031 capture corrections are present; preserve them, stage shared contract-test hunks interactively, and exclude all unrelated paths.
**Tests baseline:** ✓ REQ-136 passed warning-level ShellCheck, Bash syntax, contract regressions, suite-manifest probes, and queue-kanban Go tests.
**Dependencies:** ✓ REQ-136 is archived; Bash, Git, tar, diff, ShellCheck, and Go are available.

*Checked by work action*

## Decisions

- **D-01 (DECIDE & STATE):** Require the consuming project root to be the Git worktree root before any managed write. The approved batch contract says clients use Git, and automatic restoration of committed managed bytes cannot be truthfully guaranteed for a non-Git install.
- **D-02 (DECIDE & STATE):** Resolve every nearest existing destination parent physically before confirmation. A manifest path can be textually confined to `.claude/skills` yet still escape through a client-side symlink; the updater rejects that state without writing.
- **D-03 (DECIDE & STATE):** Replace each managed module/path completely after confirmation rather than merge-extracting it. This makes post-update byte verification exact and removes stale shipped content; dirty tracked and untracked managed content is listed before the single discard confirmation.
- **D-04 (DECIDE & STATE):** Validate future archives with the already-installed bridge validator, never the archive's replacement copy. The bridge is the trusted migration boundary; an adversarial archive must not redefine the manifest rules used to authorize its own writes.

## Implementation Summary

**Files changed:**
- `tools/do-work-update.sh` (modified) — dual-layout capability, planning, validation, review, confirmation, exact copy verification, and scoped Git recovery engine
- `_dev/tests/update-script-behavior.sh` (modified) — hermetic capability, legacy/root-fallback, suite, malformed/traversal, symlink-escape, dirty, partial-failure, and runtime-preservation probes
- `_dev/tests/contract-regressions.sh` (modified) — shared-engine and automatic-recovery contract assertions
- `actions/version.md` (modified) — agent update mode delegates to the updater used by the Just recipe

**What was done:** Shipped the bridge updater without activating the modular archive. Both public update paths now share one engine that recognizes legacy and suite artifacts, validates a complete suite before mutation, warns and confirms once, verifies reviewed bytes, and automatically recovers only its explicit managed paths.

## Qualification

Passed — four implementation files match the declared scope; the updater and test rewrite are substantive; every detailed requirement is traced by a hermetic probe or shared-engine assertion; P-A-U is complete; data flows through real archive extraction, manifest validation, managed copy, forced failure, Git restore, and byte verification rather than static-only mocks.

## Testing

**Tests run:** warning-level ShellCheck and Bash syntax for updater/contract scripts; `bash _dev/tests/record-commit-hash-guards.sh`; `bash _dev/tests/suite-manifest-contract.sh`; `bash _dev/tests/update-script-behavior.sh`; `/bin/bash _dev/tests/contract-regressions.sh`; `go test ./...`; `go vet ./...`; `go build ./...`; `go run . verify --repo-root ...`; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Bridge capability: ✗ old CLI rejected `--capabilities` → ✓ exact side-effect-free `suite-layout-v2`
- Future suite: ✗ old updater tried to read missing root `actions/version.md` and installed no module → ✓ all four validated modules install and byte-verify from one archive
- Unsafe manifests/destinations: ✗ old updater never reached suite validation → ✓ malformed, traversal, and physical symlink escapes fail before confirmation or managed writes
- Validator trust boundary: ✗ an archive-bundled permissive validator could authorize its own traversal manifest → ✓ the installed bridge validator rejects it before diff or confirmation
- Dirty/recovery contract: ✗ old updater only printed manual recovery and the action duplicated mutation logic → ✓ confirmed dirty state is discarded only inside managed paths, forced module-two failure restores the old suite automatically, and both entry points share the engine

**New tests added:** Reworked `_dev/tests/update-script-behavior.sh` as the bridge transaction suite and updated shared contract assertions.

*Verified by work action*

## Review

**Overall: 98%** | 2026-08-07T21:02:20Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 1 — `actions/install.md` still narrates the pre-bridge manual-recovery behavior in its `just-kanban` background text. Runtime behavior and the installed recipe invocation are correct; REQ-143 already owns the full-suite installer and managed recipe model, so that narrative should be replaced there rather than widening this bridge engine's explicit write set.
**Acceptance:** Pass — both archive layouts, both installed locations, the shared public entry point, dirty consent, hostile validation inputs, exact byte verification, forced partial failure, and protected project/runtime surfaces are exercised hermetically.
**Suggested testing:** At REQ-144 cutover, run the bridge-to-modular fixture against the actual staged packages in addition to these synthetic module payloads.
**Follow-ups created:** None; the non-blocking installer narrative is covered by existing dependent REQ-143.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Turning every managed write into an explicit source/destination plan made review, dirty reporting, recovery, and verification share one boundary.
- A deterministic failing `cp` wrapper tested the real trap and Git recovery path without permission-dependent fixtures.

**What didn't:**
- The first suite validator call trusted the downloaded archive's validator. An adversarial RED fixture exposed that the already-installed bridge must remain the authority for migration validation.

**Worth knowing:**
- The bridge requires the exact project Git worktree root and clears confirmed dirty managed content from both index and worktree before installing reviewed bytes.
- Future suite destinations are checked both textually by the manifest validator and physically against existing client-side symlinks.

## Orientation

[MAP CHANGED] `tools/do-work-update.sh` is now the single mutation engine for agent and Just update paths, advertises `suite-layout-v2`, and can migrate a validated four-module archive while the live upstream tarball remains monolithic behind REQ-136's export guards.
