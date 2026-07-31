---
id: REQ-042
title: "Confirm: three worktree-mode evidence-path hardening follow-ups from REQ-037's review"
status: completed
created_at: 2026-07-29T08:34:00Z
completed_at: 2026-07-29T09:32:07Z
user_request: UR-007
addendum_to: REQ-037
depends_on: []
domain: general
prime_files: []
tdd: false
review_generated: true
write_set: []
maintenance: false
---

# Confirm: three worktree-mode evidence-path hardening follow-ups from REQ-037's review

## What

While building and reviewing REQ-037 (placing the worktree merge in the pipeline and re-pointing the evidence-consuming steps at the merge range), three hardening items surfaced that were outside REQ-037's declared scope. None is a defect in what REQ-037 shipped — its feature works and is tested. Each is a "should we also do this?" question. Approve any subset; approved items become a small REQ each (or one combined REQ).

## What the Builder / Review Found

Two came from the builder's out-of-scope discoveries and one from the adversarial review's contradiction-hunter (logged as out-of-scope-but-noted). All three concern how worktree dispatch mode interacts with evidence/diff-reading paths REQ-037 did *not* touch.

## What Would Change

Each approved item is a small, independent change to the work pipeline or its checks; declining any leaves today's behavior untouched. Worktree dispatch mode is an optional advanced-harness feature, so none of these affects a normal serial run.

## Open Questions

- [x] **Standalone re-review of a worktree-merged REQ sees an empty diff — fix the standalone diff source?** → Confirmed: Yes — queued as REQ-055. In worktree dispatch mode a finished REQ's `commit:` frontmatter now holds the `--no-ff` *merge* commit hash (REQ-037 established this). When someone later runs a standalone review (`do-work review REQ-NNN`), the review's "Get the Diff" step runs `git show <commit>` on that hash — but `git show` on a merge commit prints a *combined* diff that is usually empty, so the standalone re-review would find nothing to review. (This only affects standalone re-review; the pipeline-mode review that runs during the build reads the merge range correctly and is unaffected.) The decision is yours because it trades a small `actions/review-work.md` change against accepting that worktree-merged REQs simply aren't re-reviewable standalone.
  Recommended: Yes — detect a merge commit in standalone mode and diff it as `git show --first-parent -m <commit>` (or `git diff <commit>^1..<commit>`).
  Value: Standalone re-reviews of worktree-merged REQs actually show the change instead of nothing.
  Risk: Low and reversible — a scoped `git show` flag change that only affects the merge-commit case; serial-mode `commit:` hashes are ordinary commits and are unaffected.
  Also: No — accept that standalone review isn't meaningful for worktree-merged REQs and rely on the pipeline-mode review that already ran during the build.

- [x] **qualify.sh false-FAILs when a diff edits its own debug-artifact tokens — make the scan self-aware?** → Confirmed: No / low priority — revisit only if the false FAIL recurs; no REQ queued `tools/checks/qualify.sh` flags leftover debug artifacts by grepping a REQ's diff for the literal tokens `console.log`, `debugger`, `print(`, `TODO`, `FIXME`. When a REQ's diff contains those tokens *as data* — qualify.sh's own detection pattern, a linter, or a doc about debugging — the scan matches its own pattern text and emits a false `FAIL` (this happened on REQ-037, where the new ranged branch duplicated the grep chain; the orchestrator verified zero real artifacts and overrode the FAIL by hand). The decision is yours because every clean fix is imperfect and the trigger is narrow.
  Recommended: No / low priority — the trigger is rare (only diffs that add these tokens as data) and the orchestrator's judgment already catches it.
  Value: A maintainer editing qualify.sh's pattern (or a token-carrying file) would not hit a spurious FAIL that needs a manual override.
  Risk: Low, but the fix candidates each have edge cases — excluding `qualify.sh`'s own path from its scan misses other token-carrying files; requiring the token in *added code* rather than a quoted grep pattern is hard to detect reliably.
  Also: Yes — fix it (start by excluding `tools/checks/qualify.sh`'s own path from its debug-artifact scan) if the false FAIL turns out to recur.

- [x] **Mid-run "blocked" flip can misclassify a worktree builder that failed after real work — consult the branch instead of the main tree?** → Confirmed: Yes — queued as REQ-056. The Step 8 "blocked" flip (which parks a REQ as a non-terminal hold when a builder hits a missing external precondition without having done real work) decides "did edits land this attempt?" by running `git diff` on the *main tree*. In worktree dispatch mode a builder commits on its own branch, so the main tree reads clean even when the builder did substantial work — so a worktree builder that hits a missing precondition *after* real work could be wrongly flipped to `blocked` instead of archived as `failed` (with its follow-up). This is the rare failure/blocked path, and it pre-dates REQ-037. The decision is yours because it's a narrow-path correctness nicety against leaving one more worktree-mode guard reading the wrong tree.
  Recommended: Yes, eventually — in worktree mode the "did edits land" test should consult the builder's branch (or the hand-back manifest), not the main-tree diff.
  Value: A worktree builder that fails after doing real work is classified correctly (`failed`, with a follow-up REQ) instead of parked as `blocked`.
  Risk: Low and reversible — a worktree-mode branch in one guard; the serial path is unchanged.
  Also: No — accept the mis-flip on this rare path; a `blocked` REQ is non-destructive and a human catches it at the next `do-work clarify`.

## Implementation

**No changes needed in this REQ.** All three questions resolved by the user via `do-work clarify` on 2026-07-29, each confirming the builder's recommendation. Approved items queued as REQ-055 (standalone review merge-commit diff) and REQ-056 (worktree-mode blocked-flip branch check); the qualify.sh self-scan fix was declined as low priority, to be revisited only if the false FAIL recurs.

*Resolved via clarify questions*
