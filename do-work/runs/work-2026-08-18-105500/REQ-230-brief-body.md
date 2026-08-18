---
id: REQ-230
title: "Review fix: point caller docs at the canonical publication rationale"
status: claimed
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

- [ ] `skills/do-work-toolbox/actions/present-work.md:137`: restates "`ln` and `mv` treat a directory in the destination's place as a container rather than a collision" in full, together with the portfolio helper's suffix-advance and fail-closed policy.

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
