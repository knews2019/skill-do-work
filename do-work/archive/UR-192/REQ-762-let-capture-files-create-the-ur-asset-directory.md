---
id: REQ-762
title: 'Let capture-files create the UR asset directory'
status: completed
route: A
heavy_verified_at: 2026-11-26T11:27:34Z
heavy_verified_revision: b0f2f38112842cc30ed7b9f3d8cc362d5abac8b1
review_at: 2026-11-26T11:26:22Z
commit: 4287819b305852988677c6d9bba3871231f046b1
kb_status: pending
integration_at: 2026-11-26T11:23:00Z
estimate:
  calculated_at: 2026-11-26T11:21:34Z
  p50_active_minutes: 10
  confidence: high
  basis:
    - Route A
    - 1-file write set
    - 2 acceptance criteria
created_at: 2026-11-26T11:18:47Z
user_request: UR-192
domain: general
prime_files: [_dev/primes/prime-action-files.md]
write_set: [skills/do-work/actions/capture.md]
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: false
maintenance: true
claimed_at: 2026-11-26T11:20:41Z
completed_at: 2026-11-26T11:27:34Z
release_at: 2026-11-26T11:27:34Z
---
# Let capture-files create the UR asset directory
## What
Remove the complex-capture step that creates the destination assets directory. Keep asset payloads in private staging and list their destinations in the manifest so capture-files creates the parent directories during publication.

## Original review finding

**[P2] Let `capture-files` create the UR asset directory — [capture.md:204](/Users/t2/Desktop/sa2-sentence-aligner2/.claude/skills/do-work/actions/capture.md:204)**

The complex-capture instructions now create `UR-NNN/assets/` before the transaction runs. If the manifest is refused, that directory survives rollback. A subsequent `verify-requests` selects the highest-numbered UR and then fails to find its `input.md`. Stage the asset payloads privately and let the transaction create their destination directories.

## Verification evidence

In skill-do-work2 at `f2ee470f`, commit `efcd3618` introduced the directive at `skills/do-work/actions/capture.md:204` to create the destination before publication. The same action already requires private staging at line 191 and invokes the transaction at line 211.

`skills/do-work/tools/do-work-cli/internal/publication/capture_files.go:390` records only missing directories for transaction-owned creation. A pre-created directory therefore survives refusal. An untargeted verification selects the highest-numbered UR at `skills/do-work/actions/verify-requests.md:47` and reads its `input.md` at line 55.

The review accepted F2 by inspecting the local source. The original link identifies the reporting checkout. Surface-cost is N/A because the remedy deletes a conflicting instruction.

The queue had no candidate in any UR sharing this root cause.

## Red-Green Proof

**RED prompt/case:** Follow complex capture for an asset-bearing request with a manifest that will be refused.
**Why RED now:** The instruction creates the next UR directory before validation. Refusal leaves that directory without input.md.
**GREEN when:** The explicit directory-creation bullet at skills/do-work/actions/capture.md:204 is gone. Following the instructions through refusal leaves the new UR absent, so untargeted verification still selects the previous UR. An accepted manifest publishes the asset at its declared path.
**Validation:** Inferred during capture from the accepted review evidence.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` (8,851 tokens). Its alternate-writer guidance applies to this action. The partially indexed file exceeds the 2,000-token budget.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Delete the destination directory instruction. Exercise refusal and acceptance with an asset manifest.
- [x] **[APPLY]:** Removed the complex-mode asset directory creation step.
- [x] **[UNIFY]:** Inspected the capture.md merge diff, checked whitespace and exercised both transaction outcomes.
*Source: Finding F2, pasted review (P2), accepted in skill-do-work2. Capture F2 from this review with its original P2 claim and verification evidence.*

## Triage

**Route: A** - Simple

**Reasoning:** The review identifies one conflicting bullet and the existing staging instruction that covers it.

**Planning:** Not required.

## Plan

**Planning not required** - Route A: Direct implementation.

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified)

**What was done:** Deleted the pre-publication directory creation bullet. The private staging and manifest instructions already cover asset-bearing requests.

## Qualification

Typed qualification was satisfied for 7f323216..4287819b. The diff deletes the one directive that creates a permanent path before validation. Step 4 and the common Step 5 still stage the asset and declare its manifest destination. Both publication cases were exercised against the merged action.

## Testing

**Red-green validation:** An isolated Git repository had one readable UR. With the old instruction, a missing staged asset caused `capture-files --commit` to refuse and left UR-002 as the newest directory without input.md. With the corrected instructions, the identical refusal kept UR-001 newest and readable. Supplying the asset then published UR-002/input.md and the exact 256 asset bytes.

**Repository checks:** The focused replay and the canonical `bash _dev/tests/maintainer-verify.sh` gate passed on the merged revision. Independent review also ran three existing capture publication tests covering bootstrap, asset ordering, and target collisions.

**Heavy lane selected:** staged-skills.

## Review

**Overall: 99%. Acceptance: Pass.** Independent review checked the request, the original input and the merge range. Requirements 100%, quality 96%, test adequacy 100%, scope 100%. Risk is low. The reviewer confirmed the named line was deleted and three targeted Go publication tests passed.

**Restatement sweep:** The screenshot dispatch paragraph in `skills/do-work/SKILL.md:59` still attributes permanent file installation to Step 4. The capture action stages screenshots there and publishes them in Step 5. This discrepancy applies to delegated screenshot capture. Impact-user-visible → report only.

## Lessons Learned

The manifest publisher can roll back directories it creates. A destination that an instruction created earlier survives refusal. The existing private staging instruction covers this case.

## Orientation

Complex capture keeps assets private until the request and its files can be published together.

## Heavy Verification Plan

- Base revision `7f32321658cd0bb818fb50119f3cb9b82987edeb`
- Target revision `4287819b305852988677c6d9bba3871231f046b1`
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. skills/do-work/actions/capture.md matched subtree skills.

## Heavy Verification Result

Target revision `4287819b305852988677c6d9bba3871231f046b1`, execution revision `b0f2f38112842cc30ed7b9f3d8cc362d5abac8b1`. The staged-skills lane executed in 46 seconds and exited 0.

## Timing

Observed 2026-11-26T11:23:16Z to 2026-11-26T11:26:05Z: 2m 49s total, 2m 48s attributed across 1 events, 1s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 2m 48s | 1 |

Slowest command: verification-gate / Maintainer verification, 2m 48s, exit 0, bash (2 argv tokens).
