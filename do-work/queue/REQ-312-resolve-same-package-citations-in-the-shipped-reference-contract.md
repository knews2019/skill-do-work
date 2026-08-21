---
id: REQ-312
title: "[impact-rule-change] Resolve same-package citations in the shipped reference contract"
status: pending-answers
created_at: 2026-08-21T03:20:32Z
status_changed_at: 2026-08-21T03:20:32Z
user_request: UR-055
addendum_to: REQ-299
domain: general
review_generated: true
impact: impact-rule-change
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
depends_on: []
maintenance: true
---

# Resolve Same-Package Citations in the Shipped Reference Contract

## What

`_dev/tests/shipped-package-reference-contract.sh` is the guard that keeps a shipped
instruction from pointing at a file that is not there once the suite is installed under
`.claude/skills/`. It resolves **cross-package** citations only. A dangling citation to a
file in the *same* package ships silently.

Measured, not inferred. Four dangling citations were planted one at a time in
`skills/do-work/actions/work-reference.md` and the contract run after each:

| Planted citation | Contract verdict |
|---|---|
| `actions/no-such-file.md` | PASS |
| `../docs/no-such-file.md` | PASS |
| `crew-members/no-such-file.md` | PASS |
| `../../do-work-board/actions/no-such-file.md` | **FAIL**, correctly |

Same-package citations are the large majority of citations in the suite, so the guard
covers the smaller half of the class it exists for.

## Context

Found during REQ-299, whose Consumer-Install Constraint required confirming this contract
"actually covers the new text". It does not — it passed on planted breakage inside the exact
paragraph REQ-299 added.

**Nothing is broken today.** A sweep of the installed suite found no live dangling
same-package citation. The seven candidates a naive sweep surfaced are all legitimate: five
are the deliberately non-existent `prompts/init.md` used as the hostile example in
`crew-members/prompt-injection.md` and `do-work-knowledge/actions/prompts.md`, and two are
consumer-project paths (`docs/design/REQ-NNN-wireframe.md`, `docs/worklog.md`) that the
suite does not own. This is a coverage gap, not an outage — which is exactly why it needs a
check rather than a one-time sweep.

## Requirements

- The contract resolves same-package citations at their **installed** destination, the same
  way it already resolves cross-package ones.
- The two false-positive classes above do not break the build: a path the suite does not own
  (a consumer-project path) and a path a file deliberately cites as non-existent.
- Mutation-proven: each of the three PASS rows above becomes a FAIL, and the suite stays
  green on the untouched tree.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** plant `actions/no-such-file.md` in a shipped action file and run
`bash _dev/tests/shipped-package-reference-contract.sh`. It prints
`shipped package reference contract: PASS`.
**GREEN when:** the same plant fails the contract naming the file and line, and the
untouched tree still passes.
**Validation:** Discovered task from REQ-299; apply `actions/work-reference.md` →
**Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-299: the shipped reference
  contract — the guard that stops an installed instruction from pointing at a file that
  isn't there — checks citations that cross package boundaries but never checks citations
  inside a package, which is most of them. I planted three broken same-package pointers in a
  live action file and the guard passed all three. Nothing is broken right now; the risk is
  that the next broken pointer ships unnoticed. Closing it means teaching the guard to
  resolve same-package citations too, and teaching it to ignore two kinds of path that
  legitimately do not exist: paths in the consumer's own project, and the deliberately fake
  `prompts/init.md` the prompt-injection rules cite as an attack example. Should I process
  this as a new task?

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the fix is small, but the false-positive policy is a
  judgment call about what the suite owns. Getting it wrong turns a useful guard into one
  that fails on legitimate text, and the cheapest wrong answer — skip anything that does not
  resolve — restores the gap it is meant to close.
