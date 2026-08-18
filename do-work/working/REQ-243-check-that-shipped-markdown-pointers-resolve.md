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
- _dev/tests/shipped-package-reference-contract.sh
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

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — add heading-anchor resolution to the markdown-target walk that already lives here, plus the anchor-slug generator and its fixtures

**Files I will NOT touch:** `_dev/tests/prescribed-shell-canonicalization.sh` — the declared write-set file. The sweep required by `maintenance: true` found that the reference contract already walks this exact corpus and already resolves relative targets from the citing file's own directory; adding a second walker there would have duplicated a ~500-line markdown parser. See D-01 and the deviation note below.

**Acceptance criteria (restated from REQ):**
- [x] Every relative `.md` link in shipped markdown resolves from the citing file's own directory — **already satisfied before this REQ**; proven by running the REQ's own wrong-depth mutation against the unmodified tree and getting FAIL
- [x] Where a link carries a `#anchor`, the anchor matches a heading in the target, compared the way the anchor was generated (lowercased, spaces to hyphens, punctuation dropped) — **this is the work this REQ actually did**
- [x] A broken link names the citing file, the line, the raw link, and what it resolved to — already satisfied for paths; the new anchor message matches that shape
- [x] Unresolvable link classes (external scheme, bare fragment, root-absolute, templated) are skipped explicitly, with the skip rule stated in the file as a condition rather than a list
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

**Scope deviation from the declared `write_set`, logged rather than silent:** the REQ declared `_dev/tests/prescribed-shell-canonicalization.sh`; the implementation touched `_dev/tests/shipped-package-reference-contract.sh` instead. `write_set` above is mirrored from this section, per Step 5.5.

## Pre-Flight

**Git:** ✓ clean at claim (`2432f45`); builder worked in an isolated worktree on its own branch
**Tests baseline:** ✓ `maintainer-verify.sh` exit 0 before dispatch; both affected suites green (`shipped package reference contract: PASS`, `Prescribed shell primitive canonicalization checks passed.`)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0 present

*Checked by work action*

## AI Execution State (P-A-U Loop)

<!-- Transcribed verbatim from the builder's hand-back. In worktree dispatch mode a builder may not
     write do-work/, so the orchestrator carries the filled loop into the REQ. Source of record:
     do-work/runs/work-2026-08-18-124358/REQ-243-handback.md -->


- [x] **[PLAN]**

Read `_dev/primes/prime-shell-commands.md`, CLAUDE.md, and the crew members
(`general.md`, `coding-guardrails.md`, `communication-style.md`, `maintenance.md`).

`maintenance: true` plus CLAUDE.md's **Delete before you add** made the first move a sweep, not a
build. The brief named the sweep explicitly: does `shipped-package-reference-contract.sh` already
walk these links? I read it in full (1009 lines) before writing anything.

It does. Lines 957–991 of the pre-change file walk every markdown target in every module source
tree from `suite/modules.tsv`, resolve each relative target from the citing file's own directory
(`markdown_path.parent / relative_target`), *and* resolve it a second time in installed topology
(`.claude/skills/...`), reporting `{markdown_path}:{line_number}: relative target is missing in
{source and/or installed} topology: {target}`. That is requirement 1 and requirement 3 of this REQ,
already shipped, with a ~500-line CommonMark-grade parser behind it that masks fenced code,
indented code, inline code, HTML comments, and escaped link syntax so that links inside examples
are not walked.

What it does not do is read the fragment. `urllib.parse.urlsplit(target).path` drops `#anchor`, so
`prescribed-shell-primitives.md#does-not-exist` resolves to a real file and passes.

So the plan became: **do not add a walker; add the fragment to the existing one.**

- Generate heading anchors the way the anchors were generated from the headings — GitHub's rule:
  rendered heading text, lowercased, every non-word/non-hyphen/non-space character dropped, spaces
  to hyphens, repeated slugs numbered `-1`, `-2` in document order.
- Take heading text from the **raw** line but gate the line on the **code-masked** line. Masking
  blanks inline code, which would turn ``## What `run` does`` into `## What does` and produce the
  wrong slug; but masking is exactly what makes a `#` inside a fence not a heading. Using masked to
  decide and raw to read gets both. `strip_markdown_code` preserves byte and line offsets — there
  are existing fixtures asserting precisely that — so the two line arrays index in lock step.
- Key the skip rule on a **condition**, per Closed Enumerations Go Stale, and write it into the file.
- Check the anchor in each *distinct* resolved target across the two topologies.

Planned write: `_dev/tests/shipped-package-reference-contract.sh` (not the declared write-set file).

- [x] **[APPLY]**

Written exactly as planned, one file, `+153 / -0`.

- `inline_link_label_pattern`, `atx_heading_pattern`, `heading_slug_cache` — module-level.
- `heading_rendered_text()` — collapses `[label](dest)` / `![alt](dest)` / `[label][ref]` to the
  label, to fixpoint, so a link's destination never pollutes the slug.
- `heading_anchor_slug()` — the GitHub rule above.
- `heading_anchor_slugs_from_text()` — the masked-gate/raw-read heading walk with duplicate numbering.
- `heading_anchor_slugs()` — path-keyed cache over it, `None` on unreadable.
- `run_anchor_slug_fixtures()` — six fixtures; registered next to `run_parser_fixtures()`.
- The check itself, inside the existing relative-target loop, after a new `continue` on the
  missing-target branch (so a link at a missing file is not also reported as a missing anchor).

No other file touched. No debug prints left; the instrumentation used in UNIFY was a copy under
`/tmp` and a temporary `_dev/tests/zz-req243-probe.sh` that was deleted in the same command that ran it.

- [x] **[UNIFY]**

```
 _dev/tests/shipped-package-reference-contract.sh | 153 +++++++++++++++++++++++
 1 file changed, 153 insertions(+)
```

One file changed. What I verified on it:

- `_dev/tests/shipped-package-reference-contract.sh` — read the complete diff. Additions only; no
  existing line altered except the `if missing_locations:` branch gaining a `continue`. Confirmed
  that `continue` cannot skip anything, because that branch was previously the last statement in
  the loop body. Confirmed the new code runs *after* every existing skip (scheme'd URL, dynamic
  target, root-absolute, empty path) so it can never widen what the contract walks.
- ShellCheck warning-level over 50 tracked shell files: clean (the Python lives in a quoted
  heredoc, so the added lines are data to ShellCheck; the surrounding shell is unchanged).
- Working tree clean apart from the committed file — `git status --porcelain --untracked-files=all`
  shows nothing after commit. No stray scratch files in the worktree or the repo root.
- Queue guard `git diff --name-only <pre>...HEAD -- do-work/` → empty.
- Serial-file guard over `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`,
  `CHANGELOG.md` → empty.
- Full `maintainer-verify.sh` → exit 0 (pasted below).

## Implementation Summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified)

**What was done:** Added `#anchor` resolution to the markdown-target walk the shipped-package reference contract already performs — a GitHub-rule heading-slug generator, a masked-gate/raw-read heading walk, a path-keyed cache, six slug fixtures, and the check itself inside the existing relative-target loop. The relative-path half of the REQ was found already implemented and was not rebuilt.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 (unpiped; `echo $?` printed `0` on its own line)

**Red-green validation:**
- wrong relative depth (`../../do-work/...` → `../do-work/...`): ✗ **already ✗ before the change** — this mutation was RED on the unmodified tree, which is what disproved the REQ's stated premise
- wrong anchor (`#portfolio-summary-publication` → `#portfolio-summary-publications`): ✓ passed before implementation → ✗ after → ✓ on revert. **This is the REQ's real coverage.**
- slug rule mutation (`.lower()` removed from `heading_anchor_slug`): ✓ before → ✗ after, all six fixtures naming expected vs got → ✓ on restore. Proves the fixtures are load-bearing, not decorative.

**Non-vacuity:** an instrumented copy run from inside the repo reports **27 anchors checked across 4 packages** on the live corpus, so the green run is not a green-because-everything-was-skipped.

**New tests added:**
- `run_anchor_slug_fixtures()` — six fixtures pinning the slug rule: case folding + punctuation removal, inline code contributing content not backticks, link labels replacing destinations, code blocks/comments holding no headings, repeated headings taking numbered suffixes, closing hash sequences excluded

**Existing tests updated (cross-REQ impact):** none

*Verified by work action*

## Decisions


**D-01 — The check lives in `shipped-package-reference-contract.sh`, not in the declared write-set
file.** That contract already walks the same corpus, already resolves relative targets from the
citing file's directory in both source and installed topology, and already owns the code-masking
markdown parser needed to distinguish a live link from one inside a fence. Adding an anchor walker
to `prescribed-shell-canonicalization.sh` would have meant a second parser over the same files.
Reach: any future "check something about shipped markdown links" REQ should extend this same walk.

**D-02 — Anchor slugs follow GitHub's generation rule, because that is the rule the anchors in this
repo were generated under.** The REQ says "compared the way the anchor was generated from it".
Concretely: rendered text, lowercased, characters outside `[\w\- ]` dropped, spaces to hyphens,
duplicates numbered `-1`/`-2` in document order. Reach: if the repo's docs are ever published by a
renderer with a different slug rule, this function is the single place to change.

**D-03 — Headings are read from the raw line and gated on the code-masked line.** Reading masked
text alone would corrupt any heading containing inline code (``## What `run` does`` → `what-does`);
reading raw text alone would treat `#` comments inside fenced shell blocks as headings, and this
repo's markdown is full of those. The gate works only because `strip_markdown_code` preserves line
offsets — an existing invariant with its own fixtures. Reach: anything else that needs "the real
lines of this markdown file" should use the same two-array technique rather than re-deriving fences.

**D-04 — Anchors resolve against ATX headings only, and that limitation is stated in the file.**
The repo has no setext headings. Supporting them is speculative machinery per YAGNI, and the failure
mode if someone writes one and links to it is a loud FAIL naming the file and anchor — not a silent
miss. Recorded as a decision because the honest framing matters: this is a stated limitation with a
visible failure, not an undetected hole.

**D-05 — The skip rule is stated as a condition in the file, per Closed Enumerations Go Stale.**
A link is anchor-checked exactly when it survives the pre-existing skips and carries both a fragment
and a Markdown target. Everything skipped is skipped for a reason the comment gives — a scheme'd URL
is not fetched, a bare fragment names no file, a root-absolute or templated target has no single
on-disk meaning, and a heading only exists in Markdown. No list of files or link forms to go stale.

## Lessons Learned

- **The gap a REQ names and the gap a REQ has are not always the same gap.** REQ-243's Red-Green
  Proof asserted "no check reads markdown links at all, so the wrong depth passes the entire suite
  today." Running the stated mutation against the unmodified tree took one command and disproved it:
  `shipped-package-reference-contract.sh` has resolved relative markdown targets from the citing
  file's own directory, in two topologies, for some time. Run the RED against pre-change code
  *before* building — not to satisfy a ritual, but because a "RED" that was already red is the
  cheapest possible signal that the work is half-done already.
- **Two hand-fixes of the same class is good evidence a checker is needed; it is not evidence that
  no checker exists.** REQ-230 and REQ-238 both verified the path by hand and the anchor by hand.
  The path half was machine-checked the whole time and nobody knew. Before writing the third
  instance of a manual check, grep for the check as well as for the defect.
- **"0 checks ran" from a relocated script is usually the harness.** The instrumented probe reported
  zero anchors checked when copied to `/tmp`, which reads exactly like a vacuous pass. The cause was
  `repo_root` derived from `BASH_SOURCE`. Any suite that locates itself relative to its own file
  cannot be probed from outside the tree — copy it *into* the tree, run, delete.
- **Where masking helps and where it hurts, in one file.** `strip_markdown_code` is what makes link
  extraction trustworthy and what makes heading extraction wrong, for the same reason. The
  offset-preserving property — already fixture-locked — is what lets both be right at once.

## Builder Pushback (accepted by the orchestrator)

**1. The REQ's core premise is half wrong, and I proceeded on the corrected version.**

> **Why RED now:** no check reads markdown links at all, so the wrong depth passes the entire suite
> today.

This is false. `_dev/tests/shipped-package-reference-contract.sh` reads markdown links, resolves
them from the citing file's own directory, and additionally resolves them in installed topology. The
brief's own suggested mutation fails on the unmodified tree — evidence quoted above under
*Pre-change probe*. Requirements 1 and 3 of this REQ ("every relative `.md` link resolves from the
citing file's own directory", "a broken link names the citing file, the line, the raw link, and what
it resolved to") were satisfied before I started.

The genuinely uncovered requirement was requirement 2, the `#anchor`. I built that and nothing else.
This is the **Delete before you add** sweep returning a real deletion: roughly half the specified
work was already present, and the right output was ~150 lines inside an existing walk rather than a
new walker.

**2. The declared write set was wrong, and I deviated as the brief permitted.**

Declared: `_dev/tests/prescribed-shell-canonicalization.sh`. Touched:
`_dev/tests/shipped-package-reference-contract.sh`. The REQ's own Constraints anticipated this
("check whether the shipped-package reference contract already walks these links for a different
reason and could carry this one instead of a second walker") — it does, and it now does. The
canonicalization suite proves phrases are *absent*; it has no parser and no notion of a link. Putting
a link checker there would have meant a second CommonMark-grade parser over the same file set.

No collision risk from the deviation: no sibling REQ in this wave declares
`shipped-package-reference-contract.sh`, and the file is maintainer-side (`_dev/` is export-ignored).

**3. One scope question I did *not* decide unilaterally.**

The REQ lists "bare anchors into the same file" among the things to skip, and I skipped them as
specified. Worth flagging for a future pass: now that heading slugs are computed and cached,
validating a same-file `#anchor` is about three lines and closes a hole of exactly the same class —
a `[see below](#some-section)` pointing at a heading that was renamed is silently broken today.
I left it out because the REQ named it a skip deliberately rather than by omission, and narrowing or
widening a scope decision the author made on purpose is the author's call, not mine. If you want it,
it is a very small follow-up.
