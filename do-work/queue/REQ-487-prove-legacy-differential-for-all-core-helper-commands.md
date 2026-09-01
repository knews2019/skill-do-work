---
id: REQ-487
title: 'Review fix: Prove legacy differential behavior for every core-helper command'
status: pending
domain: testing
created_at: 2026-09-01T18:13:23Z
user_request: UR-081
addendum_to: REQ-420
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
---

# Review Fix: Prove Legacy Differential Behavior for Every Core-Helper Command

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Replace the detached synthetic comparator and same-binary text/JSON checks with immutable, per-command legacy observations for all 17 retained core-helper commands. The matrix must fail when any command drifts in exact status, ordered facts, affected paths, recovery or verification argv, filesystem bytes, Git index/worktree state, or private state.

No pending REQ, including existing sweep REQs in any UR, shares this legacy-oracle root cause; the fold-first scan therefore routes this non-trivial user-visible finding to its own follow-up.

## Context
Found during the terminal re-review of REQ-420 after its single permitted remediation. The current broad parity suite is green, but the focused 17-command lane derives both renderings from the current binary and therefore cannot reject a behavioral regression shared by text and JSON.

## Requirements
- Capture immutable legacy expectations for each of the 17 public core-helper commands before changing the current implementation.
- Compare each real command execution against exact expected status, ordered facts, affected paths, recovery and verification argv, filesystem bytes, Git index/worktree state, and private state.
- Route every required mutation dimension through the actual 17-command matrix comparator, not a detached synthetic observation.
- Prove the comparator rejects one mutation in every required dimension and identifies the affected command.
- Preserve the six untracked REQ-418 run artifacts and keep this follow-up's changes isolated from them.

## Red-Green Proof
**RED prompt/case:** Run the 17-command legacy differential after independently mutating status, ordered facts, affected paths, recovery argv, verification argv, filesystem bytes, Git index/worktree state, and private state for a named command.
**Why RED now:** The current focused lane compares text and JSON from the same implementation, and its mutation comparator is connected only to one synthetic inventory observation, so shared regressions and per-command omissions pass.
**GREEN when:** Immutable expectations cover all 17 commands and every named mutation is rejected through the real per-command comparator with the drifting command and dimension identified.
**Validation:** Terminal REQ-420 review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
