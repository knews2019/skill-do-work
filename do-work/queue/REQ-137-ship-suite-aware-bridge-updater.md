---
id: REQ-137
title: "Ship the Suite-Aware Bridge Updater"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md]
tdd: true
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
An updater already executing cannot adopt replacement migration logic from the archive it is currently extracting. Existing clients therefore need a suite-aware bridge before the source layout changes.

## Detailed Requirements
- Make `actions/version.md` delegate update mutation to the same script used by `just run-do-work-update`.
- Support legacy archives and future `suite/modules.tsv` archives.
- Add `--capabilities`, which prints `suite-layout-v2` without modifying anything.
- Download and review one archive, validate every declared module before writing, show one suite-wide diff, ask once, and verify installed bytes against the reviewed archive.
- If managed skill paths are dirty, list them and warn that continuing discards those changes before accepting confirmation.
- On failure after writing, restore tracked managed skill paths from Git and clean only newly created files within validated module destinations.
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
