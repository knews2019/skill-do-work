---
id: REQ-230
title: "Review fix: point caller docs at the canonical publication rationale"
status: completed
completed_at: 2026-08-18T11:12:00Z
commit: 19669fc
claimed_at: 2026-08-18T10:56:00Z
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-18T10:56:00Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - cross-route regression gates
route: B
status_changed_at: 2026-08-18T10:26:34Z
domain: general
created_at: 2026-08-18T00:19:09Z
user_request: UR-042
addendum_to: REQ-225
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: caller-doc-restates-canonical-publication-rationale
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
write_set:
- skills/do-work-toolbox/actions/present-work.md
- _dev/tests/prescribed-shell-canonicalization.sh
---

# Review Fix: Point Caller Docs at the Canonical Publication Rationale

## What

A shipped action file carries its own full copy of the publication rationale that REQ-225 just made canonical in `skills/do-work/docs/prescribed-shell-primitives.md` § **Verified exact publication**. Replace the copy with a pointer, and extend the canonicalization test's stale-restatement patterns so a future copy fails the suite instead of surviving in it.

Done means the class cannot recur: after this REQ, a shipped file other than the canonical guide that restates the container-not-a-collision rationale fails `_dev/tests/prescribed-shell-canonicalization.sh`, the same way that test already catches restatements of the other seven canonical primitives. Patching the one instance without adding the pattern leaves the next copy free to land.

## Context

Found during review of REQ-225 by its mandatory restatement sweep. REQ-225 moved the rationale into a shared section; this is the one caller-side copy the move left stale. The canonicalization test already has the mechanism — a list of stale phrases that no shipped markdown outside the guide may contain — and this rationale was never added to it, which is why the copy has survived since it was written.

## Requirements

- Replace the restated rationale in each instance below with local intent plus a pointer to the canonical section, matching how the same file already points at `#portfolio-summary-publication` one paragraph earlier. Keep each caller's own policy statements — they are the action's policy, not the shared primitive.
- Add the rationale's phrase to the stale-pattern list in `_dev/tests/prescribed-shell-canonicalization.sh` so a future restatement in any shipped markdown fails the suite.
- Do not touch `CHANGELOG.md` or `skills/do-work/CHANGELOG.md`: their matching sentences are dated release history, and that test already skips changelogs for exactly this reason.

## Instances

- [x] `skills/do-work-toolbox/actions/present-work.md:137`: restates "`ln` and `mv` treat a directory in the destination's place as a container rather than a collision" in full, together with the portfolio helper's suffix-advance and fail-closed policy.

## Red-Green Proof

**RED prompt/case:** `_dev/tests/prescribed-shell-canonicalization.sh` with "container rather than a collision" added to its stale-pattern list.
**Why RED now:** The added pattern matches `skills/do-work-toolbox/actions/present-work.md:137`, so the test fails naming that file.
**GREEN when:** The same test passes with the pattern still in the list, because the restatement has become a pointer.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-225 I found one shipped instruction file that repeats, word for word, a paragraph of shell rationale that REQ-225 had just moved into a single shared home. Nothing is broken and no behavior changes either way — the copy is accurate today. The cost of leaving it is that the next time the shared paragraph is corrected, this copy keeps teaching the old version, which is the failure this repeated rule has already produced elsewhere. The cost of fixing it is a small edit to a shipped file plus one line in a test, and a slightly less self-contained read for anyone who was relying on that paragraph being right there. This is your call rather than mine because REQ-225 was explicitly scoped to one file with your consent recorded on it, and quietly widening that scope to a second package would spend the consent you gave on something you did not agree to. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  *Answered 2026-08-18 via `do-work clarify` — user approved queueing the pointer-plus-test-pattern fix.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" was fully specified — one named instance and one named test — but the pointer had to match the file's own existing convention and the stale pattern had to be narrow enough not to fire on legitimate helper-policy prose, both of which needed reading the surrounding code first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/present-work.md` (modify) — restatement becomes a pointer
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — add the stale pattern

**Files I will NOT touch:** `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` (dated release history; the test already skips changelogs for this reason), `skills/do-work/docs/prescribed-shell-primitives.md` (the canonical home is already correct).

**Acceptance criteria (restated from REQ):**
- [ ] The named instance carries local intent plus a pointer, matching the `#portfolio-summary-publication` pointer one paragraph earlier
- [ ] The caller's own policy statements survive — they are the action's policy, not the shared primitive
- [ ] A future restatement of this rationale in any shipped markdown fails `_dev/tests/prescribed-shell-canonicalization.sh`
- [ ] Neither changelog is touched

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified)

**What was done:** Step 6's fourth bullet no longer restates why `ln` and `mv` treat a directory in the destination's place as a container rather than a collision. It now names the canonical [Verified exact publication](../../do-work/docs/prescribed-shell-primitives.md#verified-exact-publication) check and states only this helper's own answers to it — suffix-advance for an occupied snapshot candidate, fail-closed for an occupied canonical path, occupying directory untouched and nothing left nested inside it. The construction deliberately mirrors the `#portfolio-summary-publication` pointer one paragraph above, so the two read as one convention. `'container rather than a collision'` joined the stale-pattern list in `_dev/tests/prescribed-shell-canonicalization.sh`, next to the other publication-side rationale, so the next copy of this paragraph fails the suite instead of surviving in it.

## Qualification

Passed — 2 files verified in the merge range `b2d5840..19669fc`, 4 acceptance criteria traced, no P-A-U section on this REQ (predates the block).

Judgment checks: both are substantive edits, not whitespace. The pointer's relative path was resolved independently rather than taken on the builder's word — `skills/do-work-toolbox/actions/` + `../../do-work/docs/prescribed-shell-primitives.md` normalizes to `skills/do-work/docs/prescribed-shell-primitives.md`, which exists, and `## Verified exact publication` is present in it, so the anchor is real. Nothing hollow: the added pattern is the mechanism, and its RED run proves it fires.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-canonicalization.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on both, run unpiped

**Red-green validation:** reproduced independently on the **merged** tree, not taken from the builder's report — the pattern was held fixed and only the caller file was moved between its pre-merge and post-merge versions:
- `git checkout b2d5840 -- skills/do-work-toolbox/actions/present-work.md` → `FAIL: skills/do-work-toolbox/actions/present-work.md restates canonical prescribed-shell rationale <container rather than a collision>; keep local intent and point at the guide.` (exit 1)
- `git checkout 19669fc -- …` → `Prescribed shell primitive canonicalization checks passed.` (exit 0)

Exactly one failure, naming exactly the target file, which is also the evidence the pattern is neither inert nor over-broad across the rest of the shipped markdown tree.

**New tests added:** none — the mechanism is a new entry in an existing suite's stale-pattern list, which is the shape this class of finding is pinned with.

*Verified by work action*

## Discovered Tasks

- [normal] **The bullet immediately above this REQ's instance is the same class.** `skills/do-work-toolbox/actions/present-work.md:136` restates the independent-bytes rationale ("A snapshot that shared storage with the canonical file would silently follow every later in-place edit of it"), which `skills/do-work/docs/prescribed-shell-primitives.md:88` already carries ("a snapshot linked to the canonical file would follow every later in-place edit of it"). No stale pattern covers it, so a correction to the canonical version would leave this copy teaching the old one. Two files, one instance — sized by grep, not assumed.

## Review

**Overall: 98%** | 2026-08-18T11:10:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 1 (report only)
- The link target was checked rather than trusted: `skills/do-work-toolbox/actions/` + `../../do-work/docs/prescribed-shell-primitives.md` normalizes to a file that exists and contains the `## Verified exact publication` heading the anchor names. Nothing in the suite would have caught a wrong relative depth here, which is worth knowing rather than fixing — the canonicalization test proves a restatement is *absent*, never that a pointer *resolves*.

**Restatement sweep:** this REQ redefines what shipped markdown may say about the container-not-a-collision rationale, so the sweep asked which other text states it. `grep` across `skills/` found no remaining copy — the suite passing is that check's mechanical form. The sweep did surface the adjacent independent-bytes paragraph as the same root cause; it is recorded as a Discovered Task and routed to REQ-238 rather than fixed here, because it is outside this REQ's named instance and wants its own RED/GREEN.

**Acceptance:** Pass — RED and GREEN both reproduced independently against the merged tree, and `maintainer-verify.sh` exits 0 there.

**Suggested testing:** 1 item
- Nothing in the suite verifies that a markdown pointer *resolves*. Every fix of this class replaces prose with a link, so the class trades a staleness risk for a broken-link risk that currently has no detector. Worth considering as its own decision, not as part of this REQ.

**Follow-ups created:** REQ-238; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Adding the stale pattern *first* and running the suite before touching the prose. That single ordering turns the pattern from a claim into a measurement: one failure naming one file proves simultaneously that the pattern fires on the real instance and that it does not fire anywhere else in the shipped tree. A pattern added after the fix would have proven only that it was quiet.

**What didn't:** The first attempt to reproduce RED on the merged tree used `git stash push` on a file that was clean, so it stashed nothing and the suite passed — a green result that looked like evidence and was nothing at all. Checking out the pre-merge blob by hash is the version that actually moves the file. A verification step that cannot fail is worse than none, because it reads as proof.

**Worth knowing:** The canonicalization suite proves a restatement is *absent*; it says nothing about whether the pointer that replaced it resolves. Every fix of this class swaps a staleness risk for a broken-link risk, and only the first has a detector.

## Orientation

`skills/do-work-toolbox/actions/present-work.md` now points at the shipped shell guide's **Verified exact publication** section instead of carrying its own copy of the rationale, and `_dev/tests/prescribed-shell-canonicalization.sh` will fail any future copy of that paragraph in shipped markdown. Lives in the prescribed-shell-primitives subsystem (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — this closes an instance of a class the canonical home already owned, using the mechanism that class already had. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves and its § *Closed Enumerations Go Stale* still describes the file accurately. The prime is not stale.
