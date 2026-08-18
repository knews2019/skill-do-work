---
id: REQ-238
title: "Review fix: point present-work at the canonical independent-bytes rationale"
status: pending
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

- [ ] `skills/do-work-toolbox/actions/present-work.md:136`: "A snapshot that shared storage with the canonical file would silently follow every later in-place edit of it, which is the opposite of what preserving a snapshot is for" — restating `prescribed-shell-primitives.md:88`'s "a snapshot linked to the canonical file would follow every later in-place edit of it".

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
