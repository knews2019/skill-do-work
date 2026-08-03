---
id: REQ-088
title: "Confirm: fix memory-reference.md's citation of the export-ignored CLAUDE.md"
status: pending-answers
created_at: 2026-08-03T21:54:55Z
user_request: UR-015
domain: general
addendum_to: REQ-078
discovered_during: REQ-078
maintenance: true
---

# Confirm: fix memory-reference.md's citation of the export-ignored CLAUDE.md

## What

While working REQ-078 the builder noticed, in a file it was already editing, that
`actions/memory-reference.md:88` cites this repo's `CLAUDE.md` inline:

> Steps 3–4 are scoring/formatting the agent performs on the grep output — they need no additional
> shell state, so nothing carries between command blocks (CLAUDE.md: shell state does not survive
> between prescribed blocks).

`CLAUDE.md` is the maintainer document and is `export-ignore`d from the distribution tarball, so in
every consumer install that citation points at a file that is not there. The skill's own rule for
shipped files is to restate the rule inline or point at a shipped home instead. The suite greps for
the common citation idioms; this phrasing is not one of them, which is why it survived.

It was **not** fixed inline: it is unrelated to REQ-078's premise (the timestamp command's single
home), and the implementation-time rule is to record an out-of-scope find rather than sweep it in.

## Why this is your call, not the builder's

The fix is trivial, but *which* fix is a small editorial choice — delete the parenthetical, restate
the rule in six words, or point at a shipped home — and it is a one-line change to a shipped file
that nothing is currently broken by. A dangling citation costs a consumer one confused lookup; it
corrupts nothing. So it is worth queuing and not worth a builder deciding unilaterally while working
on something else.

## What Would Change

Whichever option you pick, one line of `actions/memory-reference.md` changes and nothing else does.
If you also want the suite's citation grep widened so the next phrasing is caught mechanically, say
so — that is a second, larger change and would touch `_dev/tests/contract-regressions.sh`.

## Open Questions

- [ ] How should `actions/memory-reference.md:88`'s dangling `CLAUDE.md` citation be fixed?
  Recommended: **restate the rule inline** — replace "(CLAUDE.md: shell state does not survive
  between prescribed blocks)" with "(shell state does not survive between prescribed command
  blocks)". The sentence keeps its point and needs no external file.
  Value: the consumer reading it gets the whole rule where they are standing; nothing to look up.
  Risk: none material — it is a parenthetical in one sentence, and reverting is a one-line edit.
  Also: **(a)** delete the parenthetical entirely, since the sentence already explains itself
  ("they need no additional shell state, so nothing carries between command blocks"); **(b)** point
  at a shipped home instead of restating — but there is no single shipped file that owns this rule
  today, so this option means first choosing one; **(c)** leave it, and instead widen the suite's
  citation grep so future occurrences are caught — this does not fix the existing line.
