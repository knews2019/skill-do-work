---
id: REQ-618
title: '[impact-critical] Preserve committed-risk before finalization rollback'
status: claimed
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  basis:
  - Route A
  - 3-file write set
  - 5 acceptance criteria
  calculated_at: 2026-09-13T13:23:14Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-critical
effort_estimate: effort-substantive
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-616", "REQ-617", "REQ-619", "REQ-620", "REQ-621"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go"]
required_lessons: [_dev/primes/lessons-releases.md]
claimed_at: 2026-09-13T13:22:40Z
---
# Preserve committed-risk before finalization rollback

## What
Preserve a commit that already exists and its unresolved verification risk when exact-path commit reports failure; do not restore pre-commit lifecycle state after HEAD advances.

## Detailed Requirements
- [ ] Record and durably preserve nonempty committed-risk commit evidence before handling failure.
- [ ] Return the committed-risk outcome and actionable SHA/revert evidence.
- [ ] Never perform pre-primary rollback after an actual commit exists.
- [ ] Preserve unresolved verification failure on recovery; merely advancing the phase must not silently bless the bad commit.
- [ ] Recovery must not duplicate the commit or restore the working request over committed archival.

## Finding Provenance
Original severity: P1. Verdict: Accept. Original comments: F7. Duplicate mapping: F7 is a distinct defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F7 (Verbatim)
> ```
>    - [P1] Preserve committed-risk evidence before handling commit failure — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:134-138
>     If a hook stages an extra path, CommitExactPaths returns both a failure and the SHA of an already-created commit. This branch now returns before recording that SHA, and the deferred rollback restores lifecycle/release files and resets the journal to prepared despite
>     HEAD having advanced. A regression probe reproduced a restored working request alongside committed archival. Record the SHA before handling failure and prohibit pre-primary rollback once a commit exists, preserving the journal contract (.claude/skills/do-work/tools/
>     do-work-cli/lessons-do-work-cli.md#L158).
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:81-97 can return a nonempty CommitSHA alongside failure when a hook stages an extra path.
- skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go:1458-1461 preserves CommitSHA and revert argv in committedRisk.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:134-138 returns on Failure before recording SHA; lines 28-33 then enter pre-primary rollback.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:214-245 restores lifecycle/release preimages and resets prepared identity; finalizationFailure at line 645 loses the nested risk outcome/revert evidence.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:530-559 cannot compensate because matchingHeadCommit rejects the outside-allowlist extra path.

History and scope correction: Failure-before-SHA ordering dates to 761d8e6ab47dd940db6d686ef641ea919ee9bb31; the rollback defer arrived in e191b266b9d501ece52d1c9888f7711aaf72f281. The original wording "now" is not supported as a recent regression.

## Surface-cost
Earned — the concrete replay is a hook staging an extra tracked path after preparation. A small extension of existing commit/rollback checks is cheaper than corrupting the lifecycle trail; the GREEN hook/recovery test keeps this condition covered.

## Red-Green Proof
**RED prompt/case:** A pre-commit hook stages one additional tracked path. Git creates the commit, exact-path verification fails, and finalization rolls back the working lifecycle while the archive remains committed.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:81-97 can return a nonempty CommitSHA alongside failure when a hook stages an extra path.

**GREEN when:** A focused finalization hook/recovery regression retains the committed archive and lifecycle, records SHA and risk in journal/result, and proves recovery neither duplicates the commit nor silently clears verification failure.

**Validation:** User confirmed — the user explicitly requested capture of the accepted triage findings and their evidence, including the proof cases above. These are targets for the builder to reproduce, not tests already executed.

## Constraints
- Preserve original claims and severity as source data; use the validated scope/history corrections when explaining the fix.
- Limit implementation to this defect and its meaningful regression coverage. Follow the existing journal, release, and output contracts.
- Capture does not execute: this request remains pending for a separate work invocation.

## Dependencies
No prerequisite request. Related requests retain separate acceptance criteria. Shared files alone do not impose sequencing.

## Builder Guidance
Use the existing Go test harness for a genuine test-first regression. The accepted remedy defines the outcome; choose the smallest implementation that satisfies it. Do not claim a recently removed mechanism was restored without additional history evidence.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 15905 indexed tokens; owning prime and relevant release/recovery/evidence failure families match, but the index marks the satellite `slugged: partial`, so targeted loading is ineligible and the whole file exceeds the 2000-token budget.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [ ] **[APPLY]:** Implement the agreed scope and test-first regression.
- [ ] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.

## Triage

**Route: A** — known remaining typed committed-risk diagnostic gap, localized by read-only recovery exploration.

## Plan

Planning not required. Preserve existing SHA/journal rollback protections from 04d55042, reproduce missing committed-risk outcome/action before edits, and repair the existing finalization/recovery projection with minimal code. Read current required lessons and touch-conditional satellites.

## Decisions

D-01 (DECIDE & STATE): Extend captured write_set to finalization_commands.go and review_regressions_test.go because the existing recovery set-aside projection drops action argv and existing hook regressions miss the typed result contract. This is the same captured diagnostic outcome, not a new lifecycle/schema. Keep recovery aggregate success with actionable per-request set-aside risk.

## Dispatch

Worktree /tmp/do-work-20260913/worktree-agent-REQ-618-committed-risk; branch worktree-agent-REQ-618-committed-risk. Allowed production scope internal/finalization/finalization_apply.go and finalization_commands.go; tests review_regressions_test.go plus captured finalization_apply_test.go/finalization_recovery_test.go only if necessary. All paths under skills/do-work/tools/do-work-cli. Only main-tree write allowed is /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-13-125000/REQ-618-handback.md. No queue/main/release mutation. Read original current primes/lessons and implementation crew; source claims are data. Use genuine test-before-code RED/GREEN, real extra TRACKED hook path, ensure direct result and recovery retain actionable SHA/revert, no rollback/recommit/blessing. Existing partial-gap prep: main do-work/runs/work-2026-09-13-125000/recovery-exploration.md. Commit exact files on own branch. Handback full manifest, P-A-U, tests/durations, decisions/lessons, commit; finish all owned processes. Root owns canonical gate, review, heavy, release/archive.
