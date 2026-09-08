---
id: REQ-624
title: 'Addendum: Trim the shipped changelog to 50 entries'
status: pending
created_at: 2026-09-08T21:49:05Z
user_request: UR-130
domain: general
prime_files: ["_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
depends_on: ["REQ-623"]
related: ["REQ-622", "REQ-623"]
batch: instruction-history-cleanup
write_set: ["CHANGELOG.md", "skills/do-work/CHANGELOG.md", "CHANGELOG-20*.md"]
required_lessons: ["_dev/primes/lessons-releases.md#canonical-link-outlives-its-target", "skills/do-work/tools/lessons-do-work-update.md"]
addendum_to: REQ-022
---

# Trim the Shipped Changelog to 50 Entries

## What
Keep the newest 50 release entries in CHANGELOG.md and its identical installed mirror. Preserve all older history in dated export-ignored archives with working links. Verify history preservation and mirror equality; retain existing release cadence and version semantics.

## Prior Implementation
REQ-022 (CHANGELOG: keep newest 20 entries live, archive the rest with a no-git pointer), completed at commit `4bd7342`, split older entries into `CHANGELOG-archive.md`, added export exclusion and a tarball-safe history pointer, and preserved the recent-version view. It changed the root changelog, archive, `.gitattributes`, and the then-current `actions/version.md`. This new bounded trim uses the approved 50-entry window, current dated archive convention, and current installed mirror; the historical request stays unchanged.

## Detailed Requirements
- Keep exactly the newest 50 release entries in both live files, including the release entry for this shipped change under the current release rules.
- Preserve older entries and their order in dated archives covered by the existing `CHANGELOG-20*.md` export-ignore convention, retaining existing history and working archive links.

## Constraints
Retain existing release cadence and version semantics. Fifty is the chosen window, not a measured optimum; history packaging is the scope.

## Dependencies
Depends on REQ-623 (Reduce orchestration instruction duplication), which depends on REQ-622 (Correct verified prose drift). This request's edge preserves the user-approved serial order rather than a technical prerequisite.

## Builder Guidance
Follow the existing archive conventions. No pending or pending-answers candidate in any UR shares this request's root cause; the seven existing queued requests address different release, commit-integrity, and diagnostics defects.

## Completion Proof
Reconcile every release entry exactly once across the logical live history and archives, preserving order and existing archives; the installed mirror is a required duplicate of the live file. Check byte equality of the two live files, archive exclusion, tarball-safe history links, and existing recent version/history output. Preserve or adapt a demonstrated history consumer before trimming if it requires older entries.

## Full Context
See `do-work/user-requests/UR-130/input.md` for the user instruction, adopted report, and batch decisions. The approved prompt excerpts are in `do-work/user-requests/UR-130/assets/rev-v2-approved-scope.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the listed primes and agent rules; record a brief approach.
- [ ] **[APPLY]:** Make the scoped changes.
- [ ] **[UNIFY]:** Review every changed file and record the focused checks and outcomes.
