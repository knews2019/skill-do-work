---
id: REQ-090
title: "Confirm: seven update-script behavior probes fail on the base branch"
status: cancelled
completed_at: 2026-08-04T09:00:52Z
created_at: 2026-08-03T23:52:10Z
user_request: UR-016
addendum_to: REQ-083
discovered_during: REQ-083
domain: general
---

# Confirm: Seven Update-Script Behavior Probes Fail on the Base Branch

## What

While running the full `_dev/tests/contract-regressions.sh` suite as a side-check during REQ-083, seven
probes failed. They are unrelated to REQ-083 (which changed `tools/queue-kanban/verify.go` and its
tests) and were confirmed **pre-existing** — the identical failures appear with REQ-083's changes
stashed away.

All seven live in `_dev/tests/update-script-behavior.sh` and cover `tools/do-work-update.sh`, the
script a consumer runs to pull a newer version of this skill into their project.

## What Is Failing

Two scenarios, both about how the updater behaves when something goes wrong:

**1. A failure partway through an update** (5 probes). The probe sets up an update that dies mid-way
and expects the script to (a) exit with a non-zero status, and (b) print four specific recovery lines:
the words `Update did not complete`, the words `may be partially updated`, a runnable
`git ... checkout -- ...` command to restore the overwritten files, and a `git ... clean -nd -- ...`
command to list files the partial update added. Observed instead: the script **exits 0** and prints
none of the four.

**2. An update run against an install with uncommitted local edits** (2 probes). The probe expects a
non-zero exit and a warning containing the phrase `restores the COMMITTED content` — the point being
that `git checkout` brings back the last committed state, so a user's uncommitted edits are *not*
recoverable that way. Observed instead: exits 0, warning absent.

## Why This Needs Your Decision Rather Than Mine

There are two very different explanations and the correct fix is opposite in each case, so guessing
would be expensive:

- **The updater regressed** — the behavior existed and was lost. Then the fix is in
  `tools/do-work-update.sh`, and until it lands, a consumer whose update dies partway through gets a
  success exit code and no warning that their install is half-updated.
- **The probes were written ahead of the implementation** — the behavior was specified and tested
  first, and the script side was never built (or was deliberately deferred). Then the "fix" is either
  to build it or to retire the probes, which is a scope call, not a bug fix.

Telling these apart means reading the history of both files, and the answer changes whether this is
urgent or merely untidy. Note that a related design decision is already recorded nearby:
`_dev/tests/contract-regressions.sh` line ~730 asserts that `do-work-update.sh` must **not**
reintroduce a pre-update rollback copy, because "git is the undo" and "a mid-update failure reports the
partial install instead" — which reads as though the reporting behavior was intended to exist.

## What Would Change

If accepted as a bug: `tools/do-work-update.sh` gains a failure path that exits non-zero and prints the
four recovery lines, plus the dirty-install warning. No consumer-visible behavior changes on the happy
path.

If accepted as unbuilt-by-design: either the probes are removed (and the suite goes green), or a REQ is
captured to build the behavior deliberately. Either way the suite stops being red, which matters
because a permanently-red suite trains everyone to ignore it — REQ-083 had to stash its own changes to
establish that these seven failures were not its doing.

## Open Questions

- [x] Should the seven failing update-script probes be investigated and resolved as a work item? → No — the premise did not reproduce; cancelled as not-reproducible (see `## Cancelled`)
  Recommended: Yes, add to queue (will flip to 'pending'). The first step is diagnosis — read the
  history of `tools/do-work-update.sh` and `_dev/tests/update-script-behavior.sh` to establish which of
  the two explanations above is true, then fix accordingly.
  Value: a green suite is trustworthy; a suite with seven permanent failures makes every future REQ
  either ignore it or repeat REQ-083's stash-and-compare to prove innocence.
  Risk: if it turns out the behavior was deliberately deferred, this spends effort on something already
  decided — the diagnosis step is cheap and bounds that risk. Low and fully reversible either way.
  Also: No, discard it — you already know why these are red and it is intentional.

## Cancelled

- **When:** 2026-08-04T09:00:52Z
- **Why:** Not reproducible. The premise — seven update-script behavior probes failing on the base
  branch — does not hold. Cancelled as not-reproducible rather than fixed or built, because there is
  nothing failing to fix.
- **Decided by:** user, via `do-work abandon`

### Evidence gathered before cancelling (2026-08-04, during REQ-088)

The diagnosis step this REQ's Open Question asked for was performed. Findings:

- `bash _dev/tests/contract-regressions.sh` exits **0** with **zero** FAIL lines.
- `bash _dev/tests/update-script-behavior.sh` standalone reports `update-script behavior probes
  passed.`
- `tools/do-work-update.sh` **contains** all five strings this REQ reports as absent:
  `may be partially updated` (line 166), the `git … checkout --` restore command (194),
  `restores the COMMITTED content` (202), the `git … clean -nd --` command (204), and
  `Update did not complete.` (218).
- The probe asserts exactly those strings at lines 218, 219, 221 and 266, and they pass.

So neither of the two explanations this REQ offered applies: the updater was **not** regressed (the
behavior is present), and the probes were **not** written ahead of the implementation (the
implementation exists and the probes pass against it).

**Ruled out as causes of the earlier observation:**

- **Lingering worktrees** — `git worktree list` shows only the main checkout; no
  `worktree-agent-*` branches remain. (The prior session ran two worktree builders, so this was the
  leading hypothesis.)
- **Working-directory sensitivity** — the suite and the standalone probe both pass when invoked
  from a subdirectory, consistent with `repo_root` deriving from `BASH_SOURCE`.
- **A vacuous pass** — the probe has exactly one skip path (`git` unavailable, which prints an
  explicit `SKIP:` line and exits 0). `git` is present, so the pass is a real pass, not a skip.
- **An intervening fix** — neither `tools/do-work-update.sh` nor
  `_dev/tests/update-script-behavior.sh` has been committed to since `b583c78`, well before the
  observation.

**What remains unexplained:** why the prior session observed 8 FAIL lines twice, including once
with its own changes stashed. That is recorded rather than resolved. The failure state is not
present now, and no cause was found that would make it recur.

**If it ever returns:** the four checks above are the ones already run — start past them. The most
likely remaining candidate is something environmental and transient in that session (a shell,
toolchain, or filesystem condition not captured here), which is precisely why this was cancelled
rather than left open: an unreproducible red suite cannot be fixed, and a permanently-queued REQ
about it trains readers to ignore the queue.
