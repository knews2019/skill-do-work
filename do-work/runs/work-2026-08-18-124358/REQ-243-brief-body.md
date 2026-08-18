---
id: REQ-243
title: Check that shipped markdown pointers actually resolve
status: claimed
status_changed_at: 2026-08-18T12:43:06Z
created_at: 2026-08-18T12:20:30Z
user_request: UR-042
addendum_to: REQ-238
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- _dev/tests/prescribed-shell-canonicalization.sh
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-18T12:43:06Z
  basis:
    - Route B
    - 1-file write set
    - 5 acceptance criteria
    - full-suite verification
claimed_at: 2026-08-18T12:43:06Z
route: B
---

# Check That Shipped Markdown Pointers Actually Resolve

## What

`_dev/tests/prescribed-shell-canonicalization.sh` proves a restatement is **absent** from shipped markdown. Nothing proves that the pointer which replaced it **resolves** — not the relative path, not the `#anchor`. Add that check.

## Context

Raised as a suggested-testing item by REQ-230's review, and confirmed as a real gap by REQ-238, which is the second fix of the same class. Both REQs did the identical work by hand:

- normalize the relative path from the citing file's directory and confirm the target exists
- confirm the anchor text exists as a heading in that target

Every fix of this class replaces prose with a link. It trades a staleness risk that **has** a detector for a broken-link risk that has **none**. Two instances in, the trade is systematic rather than incidental — and the manual check is mechanical enough that a human doing it a third time is the actual defect.

Neither existing check covers it. The canonicalization suite greps for phrases that must be absent; the shipped-package reference contract checks that shipped files do not cite maintainer-only paths. A pointer to a real file at a wrong relative depth passes both.

## Requirements

- Every relative `.md` link in shipped markdown resolves, from the **citing file's own directory**, to a file that exists. The citing directory is the whole point — a link correct from the repo root and wrong from its own file is exactly the failure mode.
- Where a link carries a `#anchor`, the anchor matches a heading in the target file, compared the way the anchor was generated from it (lowercased, spaces to hyphens, punctuation dropped).
- A broken link names the citing file, the line, the raw link, and what it resolved to — enough to fix without re-deriving.
- Links to things the check cannot resolve (external `http(s)://`, `mailto:`, bare anchors into the same file, paths outside the repo) are skipped explicitly rather than silently, and the skip rule is stated in the file.
- `bash _dev/tests/maintainer-verify.sh` still exits 0.

## Constraints

- **Shipped markdown only** — the same tree the canonicalization scan already walks. Archived REQs under `do-work/` are history and routinely reference paths that have since moved; scanning them would produce noise that trains readers to ignore the check.
- No new test file if an existing suite is the natural home. This is a property of shipped markdown, and `prescribed-shell-canonicalization.sh` already walks exactly that set — check whether it belongs there before adding a file.
- `maintenance: true`: this is a pass on the skill's own instructions, so ask whether anything can be **removed** before adding. In particular, check whether the shipped-package reference contract already walks these links for a different reason and could carry this one instead of a second walker.

## Red-Green Proof

**RED prompt/case:** introduce a pointer with a wrong relative depth — e.g. change one existing `../../do-work/docs/prescribed-shell-primitives.md` to `../do-work/docs/prescribed-shell-primitives.md` — and run the check.
**Why RED now:** no check reads markdown links at all, so the wrong depth passes the entire suite today.
**GREEN when:** that mutation fails the suite naming the citing file and line, and reverting it passes. A second mutation on the anchor (`#portfolio-summary-publication` → `#portfolio-summary-publications`) fails the same way.
**Validation:** Review finding on REQ-238; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] Twice today a fix has replaced a repeated paragraph with a link to the one place that paragraph lives — REQ-230 and REQ-238, and there will be more, because that is the standard fix for this class. Each time, the test suite confirms the repeated text is gone, and nothing at all confirms the link works. Both times I checked the link by hand: that the path points at a file that exists, and that the heading it names is really in that file. A wrong link would pass every check we have. The fix is a check that walks the links in shipped instruction files and resolves them. I am asking rather than doing it because it is new machinery rather than a repair — a whole new class of check, with its own decisions about what to skip (external URLs, links into archived history) and its own maintenance cost — and you may prefer to keep checking these by hand while the class is only a few instances old. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — keep checking pointers by hand until the class is big enough to justify a checker.
  *(Answered by user via do-work clarify, 2026-08-18T12:24:42Z)*

---

# Builder Guardrails (orchestrator-issued — binding)

## Your tree

- Work **only** inside your worktree (path below). It is a full checkout on your own branch.
- **Never write anywhere in the main tree** except the single hand-back file named below. That is the only main-tree path you may touch.
- **Never touch `do-work/`** — not the queue, not `working/`, not `CHECKPOINT.md`, not `archive/`. Queue state is the orchestrator's alone. Your branch must contain **zero** commits touching `do-work/`; the orchestrator runs `git diff --name-only <pre>...<your-branch> -- do-work/` and a single path there stops your hand-back.
- **Never touch `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, or `CHANGELOG.md`.** Those are serial-only integrator-owned files. A bump on your branch races every sibling.
- **Scratch files go in `/tmp` or inside your worktree — never the main tree root.** A previous builder left a PNG in the repo root; that is a write-set violation. Screenshots, fixtures, generated boards: `/tmp`.

## Commit on your branch

Commit your implementation on your own branch before handing back. Message body only — no version bump, no changelog entry. End the message with:

```
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
```

## The P-A-U loop is yours to fill

The REQ body contains an `## AI Execution State (P-A-U Loop)` section with three checkboxes, or the orchestrator will add one. **You must tick all three and write the required content into each**, in your worktree's copy of nothing — instead, put the filled P-A-U block **into your hand-back file** verbatim under a `## P-A-U` heading, since you may not write `do-work/`. `qualify.sh` FAILs on unticked boxes and the orchestrator will otherwise have to fill them from your evidence.

- **[PLAN]** — read the listed `prime_files` and agent rules, then write the technical approach. No code yet.
- **[APPLY]** — code written exactly as planned, scope strictly limited to planned files.
- **[UNIFY]** — run `git diff --stat`, review every changed file, run the project's linters/tests, confirm no debug artifacts. List each file you verified and what you checked.

## Evidence rules — every one of these was learned by getting it wrong

1. **Two REDs when the first is a reference error.** A test that fails because a constant or function does not exist yet proves nothing. Put the code in place, break exactly one rule, and let the assertion fail *for the reason it exists*. Report both RED outputs.
2. **`git stash push` on a clean file stashes nothing** — and the resulting green run reads as proof when it is vacuous. To reproduce RED against pre-change code, check out the pre-change blob by hash (`git show <hash>:<path>`) instead.
3. **Assert page identity inside the same call that reads the DOM.** If you drive a browser, return `location.href` (and, where relevant, the page's own rule text) from the *same* `evaluate` call as every measurement. A shared browser instance can be navigated by a sibling between your navigate and your evaluate, and the numbers come back confident, well-formed, and about somebody else's page. A URL checked *before* navigating is not the same claim. Prefer an isolated browser context.
4. **A programmatic `.focus()` does not trigger `:focus-visible` in Chrome.** Use a real `Tab` keypress if focus styling is in question.
5. **Generate the artifact and look at it.** For anything that changes what appears on screen, a passing assertion is not evidence about two glyphs sharing a coordinate. Measure `getBoundingClientRect()` intersections in the live DOM when the question is "do two things overlap"; read the rendered text when the question is "what does this say".
6. **Push back if the brief is wrong.** If a requirement contradicts an existing test, or a piece of code you wrote turns out unneeded, say so in the hand-back rather than quietly editing the test or keeping dead code. Two builders pushed back last session and both were right.

## Verification bar

`bash _dev/tests/maintainer-verify.sh` from your worktree root. **Exit code 0 is the only proof.** Never pipe it through `tail`/`head` — the pipeline's exit status hides the failure. Run it, then `echo $?` on its own line, and paste that.

## Hand-back

Write **one** file, at the absolute path given below, containing:

1. `## Branch` — your branch name.
2. `## P-A-U` — the three filled, ticked checkboxes with their content.
3. `## Files Changed` — `git diff --stat` against your branch's merge base, plus one line per file saying what changed and why.
4. `## Red-Green Evidence` — the RED output(s) and the GREEN output, quoted.
5. `## Verification` — the `maintainer-verify.sh` tail and its `echo $?` line.
6. `## Integration Seams` — anything the orchestrator must apply by hand in the merge commit (shared registries, cross-REQ text). Say "none" if none.
7. `## Decisions` — numbered D-01, D-02… for choices with reach beyond this REQ.
8. `## Lessons Learned` — what a future session should know. Omit if genuinely nothing.
9. `## Pushback` — anything in this brief you think is wrong. Omit if none.

Your final message back should be a short summary; the hand-back file is the real deliverable.

## Your Assignment

- **Worktree path (your working directory):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-243-check-that-shipped-markdown-pointers-resolve`
- **Branch name:** `worktree-agent-REQ-243-check-that-shipped-markdown-pointers-resolve`
- **Hand-back file (absolute, main tree — the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-18-124358/REQ-243-handback.md`
- **Repo root of the MAIN tree (read-only for you):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Orchestrator Notes for This REQ

- `maintenance: true`. Read `_dev/primes/prime-shell-commands.md` first, and honour CLAUDE.md's **Delete before you add**: before adding a walker, check whether `_dev/tests/shipped-package-reference-contract.sh` already walks these links for a different reason and could carry this one. Report what you checked and why you chose where the check lives.
- CLAUDE.md's **Closed Enumerations Go Stale** rule applies: key the skip rule on a condition, not on a hand-maintained list of files or link forms. State the skip rule in the file itself.
- Both mutations named in the Red-Green Proof must be demonstrated: the wrong relative depth **and** the wrong anchor. Revert each and show green.
- Your declared `write_set` is one file. If the sweep convinces you the check belongs somewhere else, **say so in `## Pushback` and proceed** — but declare the actual files you touched in `## Files Changed`.
