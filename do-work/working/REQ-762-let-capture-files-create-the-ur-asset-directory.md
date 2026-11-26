---
id: REQ-762
title: 'Let capture-files create the UR asset directory'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Finding F2, pasted review (P2), accepted in skill-do-work2. Capture F2 from this review with its original P2 claim and verification evidence.*
