---
id: REQ-624
title: 'Addendum: Trim the shipped changelog to 50 entries'
status: completed
route: A
integration_at: 2026-09-08T22:57:17Z
builder_handback_at: 2026-09-08T22:56:10Z
dispatch_at: 2026-09-08T22:50:19Z
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
review_at: 2026-09-09T09:00:39Z
commit: 959cb107d4b95395dcd54d5c44f7509197b97114
heavy_verified_at: 2026-09-09T09:00:39Z
heavy_verified_revision: 707c3368aab57a7a496f05d40b763fb24cf1029a
kb_status: pending
completed_at: 2026-09-09T09:00:39Z
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
- [x] **[APPLY]:** Integrated the 610-entry history split and archive navigation in four changelog files.
- [x] **[UNIFY]:** Reviewed all changed headers and proved every release block byte-identical; merged-tree reconciliation, mirror and tarball checks passed. Native gate follows.

## Triage

**Route: A** — Simple.

**Reasoning:** One mechanical history split with a required mirror and archive navigation; exact byte/entry reconciliation supplies the proof. REQ-022 and current archive conventions are known.

**Planning:** Not required; direct builder.

## Plan

**Planning not required** — Route A: direct builder.

## Decisions

### D-01 — Retain version 0.305.37 for metadata-only history packaging

DECIDE & STATE: Current release ownership classifies CHANGELOG files as release metadata, and the finalization release guard excludes metadata-only changes from releases. The captured condition “including the release entry … under the current release rules” therefore does not add an entry for this trim. Preserve version semantics and keep the latest 50 existing entries, beginning at 0.305.37. No release-manifest or version bump applies.

### D-02 — Preserve complete-history readers through the live header

DECIDE & STATE: Direct links to every archive and an explicit reading order preserve the ADR fallback's history access within the requested packaging scope. The archive date follows the newest archived release, matching existing conventions.

### D-03 — Retain historical separator bytes at the new live EOF

DECIDE & STATE: Plain git diff --check reports one blank-at-EOF warning per live copy because the preserved 0.303.7 separator now ends the file. Retain these harmless bytes to satisfy exact historical preservation; the check passes with only blank-at-EOF disabled. No other whitespace issue exists.

## Implementation Summary

**Files changed:**
- `CHANGELOG.md` (modified) — newest 50 entries and canonical archive navigation
- `skills/do-work/CHANGELOG.md` (modified) — identical live mirror
- `CHANGELOG-2026-09-05-up-to-v0.303.6.md` (new) — 610 verbatim older entries
- `CHANGELOG-2026-07-13-up-to-v0.121.0.md` (modified) — successor navigation; release bodies unchanged

**What was done:** Retained releases 0.305.37 through 0.303.7 live and archived 0.303.6 through 0.121.1. The logical history still contains 1,012 unique release entries. Existing release semantics leave version 0.305.37 unchanged (D-01).

**Merge evidence:** `49c61ba76932abd5d6a40f7509ac3f56193dfea8..959cb107d4b95395dcd54d5c44f7509197b97114`. The pre-merge builder queue-state guard passed; exactly four captured changelog paths were integrated.

**Preservation evidence:** Builder and parent independently compared all complete release blocks against the prior tree. Every entry and separator is byte-identical, version uniqueness and descending order hold, and the three untouched archives remain whole-file identical. The root/mirror are each 60,178 bytes. Ordered concatenated release blocks before and after have SHA-256 `b33f120eba66aafc9701a172f91e68f758b18030c0484defcc271133601094e3`.

| Logical history | Entries | Newest → oldest |
|---|---:|---|
| Live | 50 | 0.305.37 → 0.303.7 |
| New September archive | 610 | 0.303.6 → 0.121.1 |
| July 13 archive | 16 | 0.121.0 → 0.110.0 |
| July 7 archive | 144 | 0.109.0 → 0.65.0 |
| April 13 archive | 56 | 0.64.1 → 0.50.0 |
| April 7 archive | 136 | 0.49.0 → 0.1.0 |

**Packaging and consumers:** Parent verified all 23 archive/header links against source paths and an actual merged-HEAD `git archive`: both live files ship identically, all five dated archives are excluded. Installed history links use canonical GitHub URLs. The version action's first-80-lines/five-releases convention still yields 0.305.33 through 0.305.37 with newest last; the builder compared complete output bytes. Full-history ADR readers reach every archive through explicit live-header guidance. Board release/parser fixtures passed (0.711 seconds). No demonstrated consumer needed a source change. New archive URLs name committed local targets; publication follows a later push.

**Required lessons:** Builder and parent read the full releases satellite (both canonical-link-outlives-its-target family bullets) and the full updater lesson satellite. No required entry was missing. Release and archive-link ownership remain unchanged.

## Discovered Tasks

None.

## Qualification

Passed by orchestrator judgment at `49c61ba76932abd5d6a40f7509ac3f56193dfea8..959cb107d4b95395dcd54d5c44f7509197b97114`. Canonical advance returned qualify state `findings`: six `QUALIFY-DEBUG-ARTIFACT-RELOCATED` warnings and four `QUALIFY-LIBRARY-OUTPUT` errors, all in the new history archive. These are historical prose quoting the checker itself, not runtime output or unfinished implementation. Parent checked every observed line against pre-change release bytes, and independent whole-history reconciliation proves exact preservation. No source code, output statements or new work markers were introduced. P-A-U and the four-file manifest agree with the merged diff. Static findings remain recorded here rather than editing history or changing the checker.

| Typed finding | Observed historical text | Disposition |
|---|---|---|
| QUALIFY-DEBUG-ARTIFACT-RELOCATED | - The untracked artifact scans skip binary files before reading them. A binary asset carrying the bytes `TODO` or `print(` produced a grep diagnostic naming a file nobody can inspect, and on BSD grep (and GNU before 3.5) that diagnostic lands on stdout, where it read as a matched source line and could fail a checked `[UNIFY]`. | Historical release prose; exact prior bytes retained. |
| QUALIFY-DEBUG-ARTIFACT-RELOCATED | - The relocation test counted substring matches, so replacing `# TODO remove deprecated parser` with a bare `# TODO` left the count unchanged and a brand-new marker read as relocated — while the warning claimed its exact text already existed. It now compares whole lines. `git grep` has no whole-line flag, so its fixed-string prefilter feeds a real `grep -c -x -F`, which needs no regex escaping on either side. | Historical release prose; exact prior bytes retained. |
| QUALIFY-DEBUG-ARTIFACT-RELOCATED | The check that catches leftover `console.log`s and stray `TODO`s could be talked out of a finding three | Historical release prose; exact prior bytes retained. |
| QUALIFY-DEBUG-ARTIFACT-RELOCATED | - `qualify.sh`'s output-primitive tokens (`print(`, `console.log`) are now judged by process-exit ownership: a file that ends its own process has a terminal audience, so its added output surfaces as a legible WARN for judgment; a file that never ends its process is library code, so the same line still FAILs, naming the file and reason. Unfinished-work markers (`TODO`, `FIXME`, `debugger`) FAIL anywhere, unchanged. | Historical release prose; exact prior bytes retained. |
| QUALIFY-DEBUG-ARTIFACT-RELOCATED | - Three lock-ins pin the boundary, including that the reporter exemption never pardons a TODO. | Historical release prose; exact prior bytes retained. |
| QUALIFY-DEBUG-ARTIFACT-RELOCATED | - `tools/checks/qualify.sh`: `grep -q` on a piped `git diff` could SIGPIPE the pipeline and mark genuinely-changed files as absent (false WARNs); the diff file list is now computed once. The debug-artifact grep now excludes `do-work/` at the pathspec level, so REQ prose merely *mentioning* console.log/TODO no longer FAILs clean implementations. | Historical release prose; exact prior bytes retained. |
| QUALIFY-LIBRARY-OUTPUT | - The untracked artifact scans skip binary files before reading them. A binary asset carrying the bytes `TODO` or `print(` produced a grep diagnostic naming a file nobody can inspect, and on BSD grep (and GNU before 3.5) that diagnostic lands on stdout, where it read as a matched source line and could fail a checked `[UNIFY]`. | Historical release prose; exact prior bytes retained. |
| QUALIFY-LIBRARY-OUTPUT | The check that catches leftover `console.log`s and stray `TODO`s could be talked out of a finding three | Historical release prose; exact prior bytes retained. |
| QUALIFY-LIBRARY-OUTPUT | - `qualify.sh`'s output-primitive tokens (`print(`, `console.log`) are now judged by process-exit ownership: a file that ends its own process has a terminal audience, so its added output surfaces as a legible WARN for judgment; a file that never ends its process is library code, so the same line still FAILs, naming the file and reason. Unfinished-work markers (`TODO`, `FIXME`, `debugger`) FAIL anywhere, unchanged. | Historical release prose; exact prior bytes retained. |
| QUALIFY-LIBRARY-OUTPUT | - `tools/checks/qualify.sh`: `grep -q` on a piped `git diff` could SIGPIPE the pipeline and mark genuinely-changed files as absent (false WARNs); the diff file list is now computed once. The debug-artifact grep now excludes `do-work/` at the pathspec level, so REQ prose merely *mentioning* console.log/TODO no longer FAILs clean implementations. | Historical release prose; exact prior bytes retained. |

## Testing

**Tests run:** Direct merged-tree `bash _dev/tests/maintainer-verify.sh` through `run-timed-command`; focused `bash _dev/tests/shipped-package-reference-contract.sh` through canonical advance; independent history reconciliation and actual git-archive inspection. Builder also ran the existing board release/parser fixture subset.

**Result:** All passed. The merged maintainer gate exited 0 in 108 seconds: 402 board tests and 815 CLI tests, every per-file budget below 30 seconds. Canonical advance returned `run-blocked-check: satisfied` and `green-gate: satisfied`. The focused probe passed using the existing recorded green baseline for that same command; Route A required no separate preflight. No new tests were added. The plain diff-check EOF warnings and qualification findings are preserved and judged above.

**Heavy verification plan:** Stored at `do-work/runs/work-2026-09-09-014954/req624-heavy-plan.json`. Range `49c61ba76932abd5d6a40f7509ac3f56193dfea8..959cb107d4b95395dcd54d5c44f7509197b97114`; selects only `staged-skills` because `skills/do-work/CHANGELOG.md` matches subtree `skills`. Exact argv: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. This heavy lane has not run for REQ-624; run it after review, then record its actual result.

## Review

**Overall: 100%** | 2026-09-09T09:00:39Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None
**Acceptance:** Pass — Exact history preservation of all 1,012 releases across 50 live entries (0.305.37 through 0.303.7) and 5 dated archives, mirror byte equality, tarball exclusion, and 23 valid header links verified on the merged tree.
**Suggested testing:** 0 items
**Follow-ups created:** None (0 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Automated whole-history block reconciliation (`verify-req624-history.py`) against git baseline proved exact byte preservation of all 1,012 releases, order, mirror equality, and tarball exclusion deterministically.

**What didn't:** Raw diff-hygiene tools flag historical release prose quoting output primitives or debug keywords. Verifying line-for-line identity against pre-change bytes provides the durable qualification proof without modifying historical notes or disarming checks.

**Worth knowing:** CHANGELOG-only modifications are classified as release metadata under current finalization release guards and do not bump the version or create synthetic release entries (D-01). Preserving trailing separator bytes at live EOF maintains exact release block hashes while causing only benign blank-at-EOF git diff check notices (D-03).

## Orientation

The shipped changelog now retains the newest 50 releases while older history is preserved across dated export-ignored archives with full GitHub navigation. Shipped release and documentation prime files remain valid; no subsystem map changed.

## Heavy Verification Plan

Base revision: `49c61ba76932abd5d6a40f7509ac3f56193dfea8`
Target revision: `959cb107d4b95395dcd54d5c44f7509197b97114`

- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/CHANGELOG.md matched subtree skills

## Heavy Verification Result

Target revision: `959cb107d4b95395dcd54d5c44f7509197b97114`. Execution revision: `707c3368aab57a7a496f05d40b763fb24cf1029a`. Stored plan reproduced exactly.
- staged-skills: exit 0, 29 seconds; executed (fingerprint_mismatch), no skip.

## Resume State

Independent review completed with 100% acceptance. Recomputed heavy verification plan reproduced the stored plan, and the staged-skills heavy lane passed (exit 0, 29s) at execution revision `707c3368aab57a7a496f05d40b763fb24cf1029a`. Proceeding to canonical finalization with supplied commit `959cb107d4b95395dcd54d5c44f7509197b97114` (no release manifest per D-01) to close UR-130, consolidate REQ-622, REQ-623, and REQ-624, remove retained worktree, refresh checkpoint, and run canonical cleanup.

## Timing

Observed 2026-09-08T22:50:19Z to 2026-09-08T22:59:57Z: 9m 38s total, 7m 40s attributed across 2 events, 1m 58s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 5m 51s | 1 |
| verification-gate | 1m 49s | 1 |

Slowest stage: builder-work / implementation, 5m 51s, outcome success.
Slowest command: verification-gate / merged-tree, 1m 49s, exit 0, bash (2 argv tokens).
