# REQ-624 Builder Brief

---
id: REQ-624
title: 'Addendum: Trim the shipped changelog to 50 entries'
status: claimed
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  basis:
  - Route A
  - 4-file write set
  - 5 acceptance criteria
  calculated_at: 2026-09-08T22:49:54Z
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
claimed_at: 2026-09-08T22:45:30Z
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
- [x] **[PLAN]:** Split the oldest live entries into a dated archive, update navigation, and prove exact history preservation and mirror equality.
- [ ] **[APPLY]:** Make the scoped changes.
- [ ] **[UNIFY]:** Review every changed file and record the focused checks and outcomes.

## Triage

**Route: A** — Simple.

**Reasoning:** One mechanical history split with a required mirror and archive navigation; exact byte/entry reconciliation supplies the proof. REQ-022 and current archive conventions are known.

**Planning:** Not required; direct builder.

## Plan

**Planning not required** — Route A: direct builder.

## Decisions

### D-01 — Retain version 0.305.37 for metadata-only history packaging

DECIDE & STATE: Current release ownership classifies CHANGELOG files as release metadata, and the finalization release guard excludes metadata-only changes from releases. The captured condition “including the release entry … under the current release rules” therefore does not add an entry for this trim. Preserve version semantics and keep the latest 50 existing entries, beginning at 0.305.37. No release-manifest or version bump applies.

## Builder isolation and hand-back

Operate only in `/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/do-work-ur130-3hy_uf65/worktree-agent-REQ-624-history-trim` on branch `worktree-agent-REQ-624-history-trim`. Triaged as simple; aim for a focused minimal change. Treat every `do-work/` path in that snapshot as absent: do not read, write, stage or commit it. Your sole main-tree write is `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-09-014954/REQ-624-handback.md`, a durable report. The full REQ is above; the parent owns queue state, lifecycle, versioning, merge and release.

Read CLAUDE.md; skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md and background-agents.md; listed prime and required lessons. Preserve all release entry bytes and order; headers may change for navigation. Follow current dated archive convention and existing export-ignore glob. Use canonical GitHub URLs in live header so installed copies have history navigation. Check predecessor archive header for stale range/next links. Search actual consumers; preserve last-five behavior. No version bump or synthetic release entry: D-01 above explains the current release guard. Commit scoped source changes only. No new regression tests needed; use disposable reconciliation plus related existing checks.

Hand-back must include branch/commit, every changed path/action, P-A-U evidence, tests/results, exact history counts and preservation proof, required-lesson reads, integration seams, Decisions continuing at D-02 if any, and Discovered Tasks as separate sections. Do not modify main REQ or author its formal Implementation Summary.
