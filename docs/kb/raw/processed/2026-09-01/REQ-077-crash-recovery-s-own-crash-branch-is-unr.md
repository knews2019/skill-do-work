---
source_type: req_lesson
req_id: REQ-077
req_path: do-work/archive/UR-015/REQ-077-crash-recovery-own-crash-branch-unreachable.md
date: 2026-08-03
domain: general
module: actions
tags: [actions, own-crash, recovery, branch, unreachable]
---

# Lessons from REQ-077: Crash recovery's own-crash branch is unreachable, and its retired premise survives in the same file

## What the REQ was about

REQ-071 (v0.164.0) gated crash recovery on the checkpoint's `## In Progress (interrupted)` record.
That record has exactly one write site — `actions/work.md:627`, Step 10, **session end** — so on a
hard crash it does not exist, every abandoned REQ classifies as a foreign claim, and the REQ leaves
the pipeline permanently. Separately, the premise REQ-071 retired still sits ~110 lines below the
paragraph REQ-071 rewrote, in the same file, and the guard meant to pin it absent matches a different
wording.

## Solution summary

Moved the `## In Progress (interrupted)` record's write site from session end (Step 10) to claim time (Step 2), which is what makes Crash Recovery's own-crash branch reachable by the event it handles. The mechanics live in a new named procedure, `actions/work-reference.md` → **In-Progress Record (Step 2)**, placed in step order between the Composed Exit Summary and the Triage template — deliberately outside the `sed`-delimited `crash_recovery_block` the suite reads, so no existing assertion's haystack changed. The procedure specifies the record as a list (one entry per claimed REQ), append-never-rewrite, create-with-only- that-section when `do-work/CHECKPOINT.md` is absent, removal on every departure from `working/`, and a deterministic literal path; a dedicated paragraph states what it is *not* (no exclusivity, no coordination, no second reader) so it cannot accrete into the lock REQ-069 deleted.

## What worked

- **Running the old suite against the old tree** was the cheapest possible proof of the regression: `OLD-SUITE-EXIT=0` with the stale sentence sitting at line 224 is one line of evidence that no amount of reading the assertion would have produced as convincingly.
- **Placing the new reference section outside the `sed`-delimited block an existing assertion reads.** Checked before writing, not after: `crash_recovery_block` runs `/^## Crash Recovery (Step 1)/,/^## Worktree Dispatch Mode/`, so a section inserted between those two headings would have silently changed six existing assertions' haystack.

## What didn't work

- **The first draft of the removal rule enumerated three movers and was wrong within the hour** — the Restatement Sweep found two more in files the REQ never mentioned. Writing a closed list *while fixing a REQ about closed lists going stale* is how strong the pull toward enumeration is.
- **A single-file `git archive` reproduction of HEAD failed twice** before working: `_dev/`, `CLAUDE.md`, and `.gitattributes` are all `export-ignore`d, so a tarball of HEAD has no test suite and the suite's own export-ignore assertions fail against it. Reproducing a suite regression needs `git show HEAD:<path>` into the live repo, not `git archive`.

## Worth knowing

- `assert_file_not_contains` is `grep -Eiq` — **case-insensitive**. One pattern therefore covers a capital-E and lowercase fingerprint for free, which is why the merged sweep needed three patterns and not six.
- The premise had **three** fingerprints in the tree, not the two the REQ diagnosed: the third (`actions/work.md:549`, Step 8 substep 6, "no lock or claim record is updated") was in neither the REQ's F3 trace nor the external audit's. Both audits found the sentence that *justified* the missing claim record; neither found the one that *asserted* it again 325 lines later. When sweeping a premise, grep the claim it makes, not the wording you found it in.
- The record now has a live machine consumer worth remembering: `tools/queue-kanban/verify.go`'s `appendCheckpointGhostFindings` flags any id the checkpoint names that exists in none of `queue/`/`working/`/`archive/`. A claim-time entry always names a `working/` REQ, so the write is ghost-safe by construction — but a future change that writes an id before its file exists would trip it.

## Back-reference

See `do-work/archive/UR-015/REQ-077-crash-recovery-own-crash-branch-unreachable.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `598ef35`.
