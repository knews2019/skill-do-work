---
id: REQ-238
title: "Review fix: point present-work at the canonical independent-bytes rationale"
status: completed
completed_at: 2026-08-18T12:20:30Z
commit:
claimed_at: 2026-08-18T11:57:10Z
route: B
estimate:
  p50_active_minutes: 15
  confidence: high
  calculated_at: 2026-08-18T11:57:10Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
status_changed_at: 2026-08-18T11:56:00Z
created_at: 2026-08-18T11:00:00Z
user_request: UR-042
addendum_to: REQ-230
domain: general
review_generated: true
effort_estimate: trivial
sweep: true
sweep_key: caller-doc-restates-canonical-publication-rationale
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
depends_on: []
write_set:
- skills/do-work-toolbox/actions/present-work.md
- _dev/tests/prescribed-shell-canonicalization.sh
---

# Review Fix: Point present-work at the Canonical Independent-Bytes Rationale

## What

The bullet immediately above the one REQ-230 just fixed is the same class: `skills/do-work-toolbox/actions/present-work.md:136` restates the independent-bytes rationale that `skills/do-work/docs/prescribed-shell-primitives.md:88` already carries. Replace it with local intent plus a pointer, and add its phrase to the stale-pattern list so a future copy fails the suite.

Done means the same thing it meant for REQ-230: after this REQ, a shipped file other than the canonical guide that restates *this* rationale fails `_dev/tests/prescribed-shell-canonicalization.sh`.

## Context

Found by REQ-230's mandatory restatement sweep — the fix landed one bullet away from a second copy of the same shape, which is why it surfaced at all. REQ-230 was scoped to one named instance with the maintainer's consent recorded on it, so widening it in place would have spent that consent on something else; it was recorded as a Discovered Task instead.

Sized by grep, not assumed: two files carry the phrase `follow every later in-place edit` — the canonical guide (line 88) and this one caller (line 136). One instance.

## Instances

- [x] `skills/do-work-toolbox/actions/present-work.md:136`: "A snapshot that shared storage with the canonical file would silently follow every later in-place edit of it, which is the opposite of what preserving a snapshot is for" — restating `prescribed-shell-primitives.md:88`'s "a snapshot linked to the canonical file would follow every later in-place edit of it".

## Requirements

- Replace the restatement with local intent plus a pointer, matching the two pointers the same file now carries one and three lines away (`#portfolio-summary-publication`, `#verified-exact-publication`).
- Keep the caller's own policy — "the helper verifies each output against the source separately and neither output can rewrite the other afterwards" is this action's policy, not the shared primitive. REQ-230 kept its equivalent clause for the same reason.
- Add the rationale's phrase to the stale-pattern list in `_dev/tests/prescribed-shell-canonicalization.sh`.
- Do not touch `CHANGELOG.md` or `skills/do-work/CHANGELOG.md` — that test already skips changelogs because their matching sentences are dated release history.

## Red-Green Proof

**RED prompt/case:** `_dev/tests/prescribed-shell-canonicalization.sh` with the new phrase added to its stale-pattern list, before the prose is changed.
**Why RED now:** the added pattern matches `skills/do-work-toolbox/actions/present-work.md:136`, so the test fails naming that file — and naming only that file, which is what proves the pattern is neither inert nor over-broad.
**GREEN when:** the same test passes with the pattern still in the list, because the restatement has become a pointer.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-230 I found a second copy of the same kind, one bullet above the one you just approved fixing. The same instruction file repeats a second paragraph of shell reasoning — about why a saved snapshot must not share storage with the file it was copied from — that also already has a single shared home. Nothing is broken and no behaviour changes either way: both copies say the right thing today. The cost of leaving it is the same one you accepted the argument for last time — the next correction to the shared paragraph leaves this copy teaching the old version. The fix is the same shape and about the same size: a small edit to one shipped file plus one line in a test. I am asking rather than just doing it because REQ-230 was deliberately scoped to one named instance with your consent recorded on it, and treating that consent as covering whatever the sweep finds next is exactly how a scoped approval quietly becomes an open one. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — one pointer per file is enough and the second copy is not worth the round trip.

---

## Triage

**Route: B** - Medium

**Reasoning:** REQ-230 established the shape an hour earlier, so the "what" was settled; the work was choosing a pattern phrase narrow enough not to fire on legitimate helper policy, which required reading the suite's exclusion mechanism first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/present-work.md` (modify) — restatement becomes a pointer
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — add the stale pattern

**Files I will NOT touch:** either `CHANGELOG.md` (dated release history the test already skips), `skills/do-work/docs/prescribed-shell-primitives.md` (the canonical home is correct), and every file the three sibling builders hold.

**Acceptance criteria (restated from REQ):**
- [ ] The named instance carries local intent plus a pointer, matching the two pointers the same file already has
- [ ] The caller's own policy clause survives
- [ ] A future restatement of this rationale in shipped markdown fails the canonicalization suite
- [ ] Neither changelog is touched

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified)

**What was done:** Step 6's independent-bytes bullet no longer explains *why* a snapshot must not share storage with the canonical file. It states the requirement locally — identical bytes as independent files — points at the canonical [Portfolio summary publication](../../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication) contract for the reason, and keeps the caller's own policy clause verbatim. `'follow every later in-place edit'` joined the stale-pattern list beside `'container rather than a collision'`, the other publication-side rationale REQ-230 added.

## Qualification

Passed — 2 files verified in the merge range `ac1fa97..d783ec9`, 4 acceptance criteria traced.

Judgment checks, run against the merged tree:
- **RED and GREEN reproduced independently**, using the builder's own suggested method rather than the one that failed on REQ-230: hold the pattern fixed and check out the *pre-merge blob* of the caller file by hash. `git checkout ac1fa97 -- …present-work.md` → `FAIL: … restates canonical prescribed-shell rationale <follow every later in-place edit>` (exit 1, one file named). `git checkout d783ec9 -- …` → passes (exit 0). REQ-230's lesson was that `git stash push` on a clean file stashes nothing and yields a green run that reads as proof; the builder flagged that trap unprompted in its hand-back.
- **The pointer was resolved, not trusted.** `skills/do-work-toolbox/actions/` + `../../do-work/docs/prescribed-shell-primitives.md` normalizes to a file that exists, and `## Portfolio summary publication` is present in it as a heading.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-canonicalization.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on both, run unpiped

**Red-green validation:** as above — one failure naming exactly one file, which is simultaneously the proof that the pattern fires on the real instance and that it is not over-broad across the shipped markdown tree the scan walks.

**New tests added:** none — the mechanism is an entry in a list the suite already maintains, which is the shape this class of finding is pinned with.

*Verified by work action*

## Decisions

- **D-01**: The pattern is the rationale phrase `follow every later in-place edit`, not the requirement vocabulary (`identical bytes`, `independent files`, `share storage`). Both existing phrasings contain it verbatim, so a paraphrased re-copy is still caught, while a caller stating *what* its outputs must be stays legal. DECIDE & STATE.
- **D-02**: The canonical guide's exemption was **verified before** choosing a phrase drawn from the canonical sentence — the scan loop skips it by full-path comparison, so the guide is exempt from every pattern rather than a curated subset. Had it been per-pattern, the phrase would have had to be one the guide does not use, which would have made it weaker. DECIDE & STATE.
- **D-03**: Kept the caller's policy clause verbatim, for REQ-230's D-03 reason: deleting a restatement is not deleting policy. DECIDE & STATE.
- **D-04**: Pointed at `#portfolio-summary-publication` — the section that owns the sentence — rather than the neighbouring bullet's `#verified-exact-publication`. The repeated pointer to the same section one paragraph above is deliberate: a reader landing on this bullet directly should not depend on having read the intro line. DECIDE & STATE.
- **D-05**: Changed "carry" to "must carry". With the rationale gone, the clause has to stand as the requirement it always was. Only wording change beyond the removal. DECIDE & STATE.

## Discovered Tasks

- [normal] **Nothing in the suite verifies that a markdown pointer resolves.** Raised as a suggested-testing item by REQ-230's review and still open; this is the second fix of the class, and it again traded a staleness risk for a broken-link risk with no detector. Both fixes needed the same by-hand check — relative path normalizes to an existing file, anchor text exists as a heading in it — and a mechanical check would get cheaper the more of this class ships. **Queued as REQ-243.**
- [low] `## Portfolio summary publication` is now referenced by name from a caller, which raises the cost of rewording it. No action: the suite's `required_heading` list already pins the heading and the new stale-pattern entry pins the rationale. Recorded so a future edit knows a pointer depends on it.

## Review

**Overall: 98%** | 2026-08-18T12:19:59Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 1 (report only)
- Second instance of a class whose *detector gap* is now confirmed rather than suspected: the canonicalization suite proves a restatement is absent and says nothing about whether the pointer replacing it resolves. Both REQ-230 and this REQ needed the same manual check. Routed to REQ-243 rather than bolted on here.

**Restatement sweep:** this REQ redefines what shipped markdown may say about the independent-bytes rationale. `grep` across `skills/` finds the phrase in exactly two files — the canonical guide (exempt by full-path skip) and the caller just fixed — and the suite passing is that check's mechanical form. Hits under `do-work/` are the queue REQ and REQ-230's archive record, which the scan does not walk and which are history. No stale restatement remains.

**Acceptance:** Pass — RED and GREEN both reproduced against the merged tree, pointer and anchor both resolved by hand, `maintainer-verify.sh` exits 0.

**Suggested testing:** none beyond REQ-243, which is that suggestion made durable.

**Follow-ups created:** REQ-243; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reading the suite's exclusion mechanism *before* choosing the pattern phrase. The canonical guide is skipped by full-path comparison, not per-pattern, which is what made the exact rationale sentence available as the pattern. Had that not been checked, the natural defensive move would have been picking a phrase the guide does not use — a weaker pattern chosen to work around a constraint that does not exist.

**What didn't:** Nothing in this build — but the trap REQ-230 recorded was avoided rather than re-hit, because the builder read that record first and named it in its own hand-back. That is the lessons chain working: the mistake cost one confusing green run an hour ago and cost nothing here.

**Worth knowing:** This class of fix has a systematic blind spot. Every instance replaces prose with a link, trading a staleness risk that *has* a detector for a broken-link risk that does not. Two instances in, the pattern is clear enough to be worth a detector of its own rather than a third round of manual checking.

## Orientation

`skills/do-work-toolbox/actions/present-work.md` now points at the shipped shell guide's **Portfolio summary publication** section for why a snapshot must not share storage with the canonical file, keeping only its own requirement and policy; the canonicalization suite will fail any future copy of that rationale in shipped markdown. Lives in the prescribed-shell-primitives subsystem (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — the second instance of a class whose canonical home and enforcement mechanism both already existed. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves and § *Closed Enumerations Go Stale* still describes the file accurately. The prime is not stale.
