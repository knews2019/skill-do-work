---
id: REQ-527
title: 'Teach cleanup Pass 5 that merged is not finished'
status: completed
priority: now
created_at: 2026-09-03T02:00:00Z
user_request: UR-086
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-458]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-458]
review_generated: true
addendum_to: REQ-458
write_set:
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go
  - skills/do-work/actions/cleanup.md
  - skills/do-work/docs/cleanup-guide.md
  - skills/do-work/actions/work-reference.md
claimed_at: 2026-09-03T23:34:31Z
route: C
planning_at: 2026-09-03T23:39:35Z
exploration_at: 2026-09-03T23:39:35Z
dispatch_at: 2026-09-03T23:42:48Z
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-03T23:39:35Z
  basis:
    - Route C
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
completed_at: 2026-09-04T00:22:00Z
commit: 9ccbc9f2928a
---

# Teach Cleanup Pass 5 That Merged Is Not Finished

## What

REQ-458 stopped `verify` from advertising an active builder's worktree as a fixable leftover. `cleanup` Pass 5 still decides on merged-ness alone, so the two now contradict each other in the opposite direction: verify protects the worktree, cleanup removes it.

## Instances

- `skills/do-work/actions/cleanup.md:140` — "Both commands succeeding means the work was merged and the leftover was pure residue: report `Removed merged worktree <name>`." For a clean worktree belonging to a run that has not reached review, `git worktree remove` and `git branch -d` both succeed, so Pass 5 removes an in-flight builder's worktree with no consent gate.
- `skills/do-work/docs/cleanup-guide.md:29` — "Already-merged leftovers are removed automatically." Same premise, user-facing.
- `skills/do-work/docs/cleanup-guide.md:56` — same contract restated.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Require a fresh, exact request-tree observation alongside cleanliness and mergedness; fail closed to the existing consent path for every absent, ambiguous, malformed, unreadable, or working state.
- [x] **[APPLY]:** Added settled-request evidence to Pass 5, regression fixtures for every request state, and aligned the three shipped cleanup restatements.
- [x] **[UNIFY]:** Reviewed the exact five-file diff; focused, race, full-module, vet, shipped-reference, sequential contract, and whitespace checks passed.

## Finding Provenance

- **Finding F2** — `impact-user-visible` — from REQ-458's independent review (Approve with follow-ups, 73%), found by its restatement sweep. The builder's own sweep missed it; the reviewer swept the whole repository for consumers of the redefined contract and found three.
- REQ-458's Scope declared `cleanup.md` as "will NOT touch" on the premise that "Pass 5's consent-gated behavior is already correct". That premise holds for **dirty** worktrees, where non-forced `git worktree remove` refuses. It is false for **clean in-flight** ones, where both commands succeed.

## Detailed Requirements

- Pass 5 must not mechanically remove a `worktree-agent-*` worktree whose REQ is still in `do-work/working/`. Both commands succeeding is evidence the work is merged, never evidence the run is over.
- Use the same evidence REQ-458 established — the leftover name's `REQ-NNN` id and that REQ's tree section — so verify and cleanup agree on what "finished" means rather than each deciding separately.
- Introduce no liveness signal: no heartbeat, lock, PID probe, mtime heuristic, claim registry, or time threshold. REQ-073's ban applies here exactly as it does to verify.
- Keep removing genuinely finished, clean, merged residue mechanically; this must not turn Pass 5 into a prompt for the ordinary case.
- Keep `cleanup.md` and `cleanup-guide.md` consistent with each other and with `verify`'s corrected categories.

## Constraints

- Never `-D`, never `--force`. The existing non-forced commands stay non-forced.
- Do not relocate the consent gate; extend what routes to it.

## Dependencies

Depends on REQ-458, which establishes the evidence and the corrected verify categories this must agree with.

## Red-Green Proof

**RED prompt/case:** Run cleanup Pass 5 against a tree holding a `worktree-agent-REQ-NNN-*` worktree that is clean and merged while `do-work/working/REQ-NNN-*.md` still exists. Today both commands succeed and the worktree is reported as removed residue.
**Why RED now:** Pass 5 routes on the two commands' exit status alone, which answers "is it merged" and not "is the run over".
**GREEN when:** That worktree is left in place and reported as belonging to an unfinished run; a clean merged worktree whose REQ is archived is still removed mechanically; and `cleanup-guide.md` no longer says already-merged leftovers are removed automatically without qualification.

---
*Source: REQ-458 independent review finding F2.*

---

## Triage

**Route: C** - Complex

**Reasoning:** The captured two-document scope cannot fix the executable cleanup path. The behavioral correction spans Go cleanup planning/tests plus every shipped Pass 5 and crash-recovery restatement.

**Planning:** Required

## Plan

1. Add RED fixtures for clean merged worktrees and branch-only leftovers whose exact REQ is still in `do-work/working/`, while preserving automatic removal for positively settled archived REQs.
2. Discover current request-tree evidence at Pass 5 execution time, parse the owner only from the anchored `worktree-agent-REQ-NNN-...` name, and classify working/settled/unknown without a liveness signal.
3. Admit automatic removal only when the candidate is clean, merged, and positively settled; route unfinished or unknown candidates through the existing consent-required finding without mutation.
4. Update cleanup action, guide, and crash-recovery restatements to the same three-fact rule, then run focused cleanup tests, contract checks, vet, and the canonical gate.

**Plan validation:** All detailed requirements map to these four tasks. The executable change and its exact tests are necessary to satisfy the user-visible behavior; prose changes only align the alternate consumers. No new force path or liveness mechanism is introduced.

*Generated from the Plan-agent findings*

## Exploration

- `internal/cleanup/cleanup_git.go` currently decides from cleanliness and mergedness alone, then mechanically removes the candidate without consulting request-tree state.
- Existing cleanup fixtures encode the old premise by removing clean merged REQ worktrees with no request evidence.
- REQ-458 established the three facts that make residue fixable: merged branch, clean worktree, and a positively identified REQ outside `do-work/working/`; absent or ambiguous identity is unknown, not finished.
- `actions/cleanup.md`, `docs/cleanup-guide.md`, and `actions/work-reference.md` all restate the merged-only rule and must move together under the alternate-writer contract.
- The existing explicit `--discard-worktree` consent path already owns forced deletion; this REQ does not add force or widen automatic removal.

*Generated from the Explore-agent findings*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modify) — require exact settled-REQ evidence before automatic removal
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (modify) — active, settled, unknown, and branch-only RED/GREEN fixtures
- `skills/do-work/actions/cleanup.md` (modify) — Pass 5 mechanics and consent wording
- `skills/do-work/docs/cleanup-guide.md` (modify) — user-facing three-fact rule
- `skills/do-work/actions/work-reference.md` (modify) — align crash-recovery cleanup semantics

**Files I will NOT touch:** queue-kanban verify code, cleanup command registration, or the explicit-consent force path.

**Acceptance criteria (restated from REQ):**
- [ ] A clean merged `worktree-agent-*` candidate whose REQ remains in `do-work/working/` is preserved and reported as unfinished.
- [ ] Cleanup and verify derive finishedness from the same anchored REQ id and request-tree evidence.
- [ ] No heartbeat, lock, PID, mtime, claim registry, or time threshold is introduced.
- [ ] Clean merged residue with a positively settled REQ is still removed mechanically.
- [ ] Cleanup action, guide, and crash-recovery reference state the same rule and retain the existing consent gate.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/docs/cleanup-guide.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)

**What was done:** Cleanup Pass 5 now automatically removes a builder residue only when a fresh observation proves the worktree clean, its branch or detached head merged into the integration head, and exactly one readable matching request settled outside the working tree. Every missing, negative, malformed, unreadable, or ambiguous fact uses the established consent-required finding and exact discard command without mutation. Shipped action, guide, and recovery wording now state the same three-fact rule. Integrated at 9ccbc9f2928a50ad155e7e3101880feff7393df0 from implementation range f57d9176..9ccbc9f2928a50ad155e7e3101880feff7393df0.

## Decisions

- **D-01:** Refresh repository discovery when Pass 5 begins so earlier cleanup passes cannot leave the deletion decision using stale lifecycle evidence.
- **D-02:** Treat queue and archive as settled only when exactly one parseable request owns the anchored worktree REQ id; all collision or read evidence fails closed.
- **D-03:** Reuse the existing consent finding and exact discard path rather than adding liveness, age, or new force mechanisms.

## Testing

- PASS: focused RED/GREEN cleanup fixture group covering active, absent, ambiguous, malformed, unreadable, settled, unmerged, detached, and explicit-discard behavior.
- PASS: `go test -count=1 ./internal/cleanup`; `go test -race -count=1 ./internal/cleanup`; `go vet ./...`; `go test -count=1 ./...`.
- PASS: `bash _dev/tests/shipped-package-reference-contract.sh`; sequential `bash _dev/tests/contract-regressions.sh`; `git diff --check`.
- One parallel contract run crossed a per-file timing budget; its standalone fixture and the sequential complete rerun passed.

## Review

**Overall: 73%** | 2026-09-04T00:00:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 75% |
| Code Quality | 80% |
| Test Adequacy | 78% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Dry-run can report consent-required for a terminal REQ still in `working/` while the real run first moves that REQ and then automatically deletes its worktree (`cleanup_commands.go:40-43`; `cleanup_git.go:215,261,268`) — `impact-user-visible` → report only
- Board verify deduplicates ambiguous REQ identity into one settled winner and advertises the worktree `[fixable]`, while cleanup refuses the same candidate as `request_state=ambiguous` (`queue-kanban/model.go:667-701`; `queue-kanban/verify.go:912-920`; `cleanup_git.go:373-383`) — `impact-user-visible` → report only
- `actions/cleanup.md:128` still says only unmerged work needs explicit consent, omitting every newly consent-gated state — `impact-rule-change` → report only
- Existing REQ-409 cleanup tests were behaviorally updated without the required cross-REQ provenance in Testing (`cleanup_git_test.go:141-177,313-333`) — `impact-negligible` → report only

**Minor findings:** `actions/restart-with-parallel-handoff.md:51-55` still equates removable with an archived claim rather than the new exact outside-`working/` settled state — `impact-rule-change` → report only
**Acceptance:** Partial — active-run preservation and settled cleanup pass, but built fixtures reproduce dry-run/apply and board/cleanup contradictions.
**Suggested testing:** 2 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

## Remediation

One bounded remediation pass must align dry-run with post-Pass-0 cleanup decisions, make board verification fail closed on duplicate identity, correct the two stale operative restatements, and record the intentional REQ-409 test updates.
