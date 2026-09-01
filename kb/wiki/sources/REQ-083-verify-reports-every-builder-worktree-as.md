---
title: "Lessons from REQ-083: verify reports every builder worktree as a fixable orphan, including active and unmerged ones"
type: source-summary
topic_cluster: worktree-and-parallel-dispatch
sources: [raw/processed/2026-09-01/REQ-083-verify-reports-every-builder-worktree-as.md]
related:
  - page: concept-worktree-isolation-and-parallelism
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-083: verify reports every builder worktree as a fixable orphan, including active and unmerged ones

Part of the [[concept-worktree-isolation-and-parallelism]] cluster.

## What the REQ was about

`appendWorktreeFindings` emits an `orphan-worktree` finding with `Fixable: true` for every
`worktree-agent-*` worktree and branch it finds, with no merged / unmerged / still-building
distinction. So an in-flight builder is reported as a leftover, and unmerged builder work — which
`actions/cleanup.md` Pass 5 will not delete without asking — is advertised in the report footer as
mechanically fixable.

## Solution summary

`appendWorktreeFindings` now asks `git merge-base --is-ancestor <name> HEAD` (pinned to the repo root with `-C`) for each `worktree-agent-*` leftover and routes the answer through `routeWorktreeLeftover`, which sets `Fixable: true` for merged residue only. `verifyCategoryOrphanWorktree` is replaced by three categories — `merged-worktree-leftover`, `unmerged-worktree-leftover`, `worktree-merge-state-undetermined` — so the report states the case it found instead of leaving it in remedy prose. Tests gained a real `git init` + `git worktree add` fixture covering all five states from requirement 6 plus a `FixableCount() == 0` contract assertion. Check 14's probe row in `actions/forensics.md` was updated to describe the split.

## What worked

- **Landing the constants before the logic bought a real RED.** Adding the three category names first, with the old unconditional `Fixable: true` still in place, turned a compile error into an assertion failure — `FixableCount = 1, want 0` against a rendered report showing `[fixable]` and `1 fixable: run do-work cleanup`. A build error proves the test is new; only an assertion failure proves the bug is real. This is the session's fourth consecutive case of observing the RED paying for itself in under a minute.
- **Reading `actions/cleanup.md` Pass 5 first turned a design question into a lookup.** "What should fixable mean here" reads open-ended until you notice Pass 5 has exactly three branches — mechanical, consent-gated, and refuse-to-act-without-a-merge-target. The three categories fell out of that table rather than being invented; the REQ's Builder Guidance pointing at Pass 5 was the highest-value line in the file.

## What didn't work

- **The requirement list was a floor, exactly as this session's checkpoint warned.** Requirement 6 enumerates five fixture cases, none of which exercises the third state the split introduces — so a fixture that satisfied the requirement in full still left a new code path untested. The review caught it; the requirement never would have. A REQ that says "split into N states" and separately lists fixtures is two lists that can disagree, and here they did.
- **The first live reproduction misled for a minute.** Creating the worktree at `../wt-active` while naming its branch `worktree-agent-REQ-001-active` printed `branch only (no worktree)` and looked like a detection bug. It was the fixture: worktree dispatch mode keeps the directory basename and the branch name *the same string*, and `listWorktreeAgentWorktrees` keys on the basename. Any manual worktree reproduction has to honor that or it tests a shape the skill never produces.

## Worth knowing

- **`git merge-base --is-ancestor` has three outcomes, not two, and the third one matters.** Exit 0 is merged, exit 1 is git's answer "no", and anything else (usually 128) is git *declining to answer* — most often a worktree whose branch is gone. Folding 128 into "unmerged" would report an unanswered question as a definite answer, which is the same defect class this REQ fixed, one level down.
- **Merged-ness is HEAD-relative and the question must be pinned to the repo root.** `git -C repoRoot` is doing real work here, not tidiness: asked from a builder's worktree, a merged branch reads unmerged and vice versa. Same trap `git branch -d` carries in `actions/work-reference.md` → Worktree Dispatch Mode, *Cleanup — happy path*.
- **`verify_test.go` now has a real git-repo fixture** (`newWorktreeFixtureRepo` and friends). Before this REQ every test used a plain temp directory, which silently *skipped* all worktree probes — the coverage hole REQ-072's review predicted. Any future probe that shells out to git should build on these helpers rather than adding a second fixture style.
- **`_dev/tests/contract-regressions.sh` is red on the base branch** — 7 update-script probes. Confirm pre-existing by stashing before treating any of it as your own regression. Filed as a Discovered Task.
  - **[Superseded 2026-08-04, not part of REQ-083's own record]** This observation does not reproduce. The follow-up it was filed as (REQ-090) was cancelled as not-reproducible: the suite exits 0 with zero FAIL lines, the seven probes pass, and `tools/do-work-update.sh` contains all five strings they assert. Lingering worktrees, cwd sensitivity, a vacuous skip path, and an intervening fix were all ruled out; the original observation remains unexplained. The transferable lesson is the *habit* — stash and re-run before blaming your own change — not the red state, which was never confirmed to exist outside that one session.

## Back-reference

See `do-work/archive/UR-016/REQ-083-verify-advertises-unmerged-work-as-fixable.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f6c1514`.
