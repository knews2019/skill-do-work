---
id: REQ-569
title: 'Refuse a release that leaves a version or changelog mirror undeclared'
status: claimed
created_at: 2026-09-04T22:35:42Z
user_request: UR-113
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-releases.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-565]
claimed_at: 2026-09-04T22:56:08Z
---

# Refuse a Release That Leaves a Version or Changelog Mirror Undeclared

## What

Make the release planner in `do-work-cli` discover the version and changelog mirrors itself and refuse, with a typed finding, any release manifest whose declared targets leave a mirror behind on the old version. The set of mirrors is a condition, not a list: every tracked file the planner can read as a plain version file that currently carries the old version, and every tracked changelog whose bytes match the declared changelog preimage.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The 0.282.0 release commit (cefe971d) wrote root `VERSION` and root `CHANGELOG.md` but not `skills/do-work/actions/version.md` or `skills/do-work/CHANGELOG.md`. The board's version-changelog probe reported the drift and a repair commit (2a5a698b) closed it by hand. The planner already refuses mirrors that disagree, stale preimages, and payloads without the version string, but only across the targets the manifest names. Nothing asks which tracked files still carry the old version, so a forgotten mirror passes. The user's words: "isn't this a mechanical sync? since it's able to check it, why not update it as well?"

## Context

- `internal/publication/release.go` → `BuildReleasePlan` is the single chokepoint for the direct `release` command, finalization prepare, and journal replay, and it already receives the repository root.
- `internal/finalization/finalization_discovery.go` → `configuredReleaseMetadataPaths` already computes the same mirror set for recovery-time discovery (VERSION basenames, the version action file's `**Current version**:` line, changelogs byte-identical to the changed one). The plan-time check must agree with it, so the two never disagree about what a mirror is.
- Workspace package manifests (npm, Cargo, uv) are owned by REQ-512 and REQ-565 (closing residual workspace release identity gaps); this REQ does not widen into them.
- The board probe stays read-only. It cannot know whether the version file or the entry heading is wrong, and the reverse direction is a normal mid-release state.

## Requirements

- Enumerate tracked files at plan time; an enumeration failure refuses rather than skipping the check.
- Refuse with one typed code that names every undeclared mirror path and the remedy (declare it as a target, or explain why it is not a mirror).
- Keep every existing release refusal and acceptance unchanged, including maintainer mirrors, project-owned partitions, and bootstrap changelogs.
- Add one sentence to the Changelog Entry Procedure in `actions/work-reference.md` so the manifest author knows the planner enforces the mirror set.
- Record the lesson in the do-work-cli lessons satellite and refresh its token count in `do-work/lessons-index.md`.

## Red-Green Proof

**RED prompt/case:** Build a release plan in a git-tracked fixture holding `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, and `skills/do-work/CHANGELOG.md` all on 1.0.0, with a manifest that declares only root `VERSION` and root `CHANGELOG.md`.

**Why RED now:** `BuildReleasePlan` returns a runnable two-mutation plan; the three undeclared mirrors stay on 1.0.0 after publication, which is exactly the 0.282.0 drift.

**GREEN when:** The plan is refused with a typed code whose paths list the three undeclared mirrors, and the same fixture with all five declared still produces a five-mutation plan. Existing release tests stay green.

**Validation:** Inferred during capture from the 0.282.0 incident; user approved the proposed fix ("capture it, and build it").

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6625 tokens, over the 2000-token budget and `slugged: partial`, so no targeted entry is legal. Matches families `closed-enumeration-for-a-condition` and `publication-target-topology-classification`.

*Source: "isn't this a mechanical sync? since it's able to check it, why not update it as well?" / "capture it, and build it"*
