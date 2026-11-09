---
source_type: req_lesson
req_id: REQ-084
req_path: do-work/archive/UR-016/REQ-084-verify-misses-committed-owner-impersonation.md
date: 2026-08-04
domain: general
module: tools/queue-kanban
tags: [queue-kanban, queue-state, do-work, verify, probe]
---

# Lessons from REQ-084: verify's queue-state probe misses a builder that committed its do-work edits

## What the REQ was about

The owner-impersonation probe runs only `git status --porcelain … -- do-work/` inside the builder's
worktree, so it sees uncommitted changes and nothing else. A builder that edits `do-work/` and
**commits** on its own branch leaves a clean worktree and passes verification — which is the more
likely shape, since a builder commits its work by design.

## Solution summary

`appendWorktreeFindings` now resolves the integration ref once per run (`resolveIntegrationBranchRef`, from the repo-root checkout, falling back to the commit id on a detached checkout) and runs `worktreeCommittedQueueState` — `git diff --name-only <ref>...<branch> -- do-work/` — for every `worktree-agent-*` name, emitting the new `worktree-committed-queue-state` finding when the builder's branch carries queue edits. The existing porcelain check and `worktreeDirtyQueueState` are untouched, so both states are now covered and reported distinguishably. Neither is `Fixable`. A ref that cannot be resolved becomes a `SkippedProbes` entry rather than silence.

## What worked

- **Scoping the probe to the rule, not to a detection method.** The root cause here was not a missed case — it was a probe named after `git status` instead of after "a builder must not write queue state." Once the question was re-asked as "what does the *branch* change", three-dot diff semantics answered it in one command. Worth asking of any probe: does its name describe the rule or the tool?
- **The two passing tests carry the real weight, but only because a failing one exists.** The stale-snapshot and porcelain-survives guards would both pass against a probe that never fires; they become meaningful the moment the RED test proves it does. Writing all three together is what makes the widening safe rather than merely asserted.
- **Requirement 6 asked for confirmation and the fixture beat the argument.** The structural reasoning from REQ-082 was sound and would have been accepted — but planting a real hand-back file in the main tree and asserting silence costs four lines and cannot rot the way an argument can.

## What didn't work

- **The first live reproduction ran against a stale binary and printed nothing.** `queue-kanban` had been built during REQ-083, so the manual GREEN check silently exercised the old code and looked like the fix had failed. Any manual reproduction in this module must rebuild first — the binary is gitignored and there is nothing to remind you it is old.

## Worth knowing

- **`git diff A...B` (three dots) is the only comparison that is safe here.** It diffs merge-base(A,B) to B, so it reports what the builder's branch changed and is blind to how far the integration branch has moved. Two dots or `diff -r` would fire on every run in a repo that commits `do-work/`, because the orchestrator claims and archives constantly — which is exactly why the original probe narrowed to porcelain instead. The narrowing was a reasonable response to a real problem; it just traded away the main case rather than picking a sharper comparison.
- **The integration ref must be named, not passed as `HEAD`.** From the repo root they are the same thing, which is what makes the shortcut tempting — but `HEAD` means "whatever checkout this command runs in", so inside a worktree the comparison silently becomes branch-against-itself. Same class as `git branch -d`'s trap.
- **`appendWorktreeFindings` now runs three sub-probes and is the file's longest function.** Two REQs in a row were forbidden from restructuring `verify.go`, correctly — but the next change in this neighbourhood should extract rather than append.

## Back-reference

See `do-work/archive/UR-016/REQ-084-verify-misses-committed-owner-impersonation.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0d61054`.
