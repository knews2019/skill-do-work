---
id: REQ-269
title: Draw the cross-package citation class by what a citation is, not by the punctuation around it
status: pending-answers
created_at: 2026-08-18T21:44:57Z
status_changed_at: 2026-08-18T21:44:57Z
user_request: UR-055
addendum_to: REQ-259
domain: general
review_generated: true
sweep: true
sweep_key: citation-class-drawn-at-the-marker
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- _dev/tests/shipped-package-reference-contract.sh
- _dev/primes/prime-action-files.md
- skills/do-work/crew-members/prompt-injection.md
- skills/do-work-toolbox/crew-members/prompt-injection.md
- skills/do-work-knowledge/crew-members/prompt-injection.md
- skills/do-work/actions/work-reference.md
---

# Draw the Cross-Package Citation Class by What a Citation Is, Not by the Punctuation Around It

## What

Eight consecutive REQs in this area have bounded the cross-package citation rule by a **marker** — a leading `../`, a backtick, a fence character — instead of by the **thing** the rule governs: a path in shipped text that a reader is expected to follow from the citing file's own directory. Each fix closes one spelling and leaves the class open, so the next review finds the same defect wearing different punctuation. This REQ ends that by making the checker's condition *be* the rule's condition.

Concretely, three markers are currently doing the bounding, and each has live escapees:

1. **The backtick.** `_dev/tests/shipped-package-reference-contract.sh` checks backticked spans only. Two of REQ-259's three sites existed solely because they were bare text; REQ-259 fixed them by backticking, which moved those two sites into coverage without closing the hole.
2. **The leading `../`.** REQ-259's own sweep found the same skills-folder base spelled with no `../` at all and reasoned it out of the class on the grounds that "it never claimed to be a relative path" — a spelling test standing in for a semantic one.
3. **The fence.** REQ-249 exempted fenced blocks with a stated rationale: their text lands in *some other file*, so it is content rather than a citation from here. The exemption is keyed on the fence character instead of on that rationale, so it also shields in-file annotations that are citations from the file in every sense that matters.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-259's independent review, Important finding 1 (gate: user-visible), consolidating the builder's own D-T1, D-T2 and D-T3 — one root cause, so one sweep REQ rather than three follow-ups. Created `pending-answers` per the generation-≥2 cascade depth stop, since REQ-259 is itself `review_generated: true`.

The review verified the hole by mutation rather than by argument: restoring the exact pre-REQ-259 state (bare text at the wrong depth) makes `shipped-package-reference-contract.sh` **PASS**. That is the class, still open, demonstrated.

## Instances

- [ ] **`skills/do-work/crew-members/prompt-injection.md:3`** — `through do-work-toolbox/actions/completed-work-presentation-reference.md`, bare, no `../`; resolves from nowhere but `skills/`.
- [ ] **`skills/do-work-toolbox/crew-members/prompt-injection.md:3`** — same text; wrong twice over, since from this file the correct citation is the same-package `../actions/completed-work-presentation-reference.md`.
- [ ] **`skills/do-work-knowledge/crew-members/prompt-injection.md:3`** — same text.
- [ ] **`skills/do-work/actions/work-reference.md` lines 130, 132, 137, 204** — four `../do-work-board/...` citations at the retired one-`../` depth (correct is `../../` from `actions/`), shielded by the fence around the Schema Read Contract yaml block whose `#` annotations are documentation for the agent reading that very file.
- [ ] **The fence exemption itself** — split it by its stated rationale (fenced text that lands in another file) rather than by the fence character, in `_dev/primes/prime-action-files.md` § Cross-Referencing and in the checker.
- [ ] **The checker's span condition** — widen `_dev/tests/shipped-package-reference-contract.sh` past backticked spans so a bare-text citation is checked like any other.

## Requirements

- The checker's condition is the rule's condition: a cross-package citation in shipped text is checked whether or not it is backticked, and whether or not it leads with `../`.
- The fence exemption is decided by whether the text lands in another file, not by the presence of a fence.
- **Done means the class cannot recur, not that six more spots were patched.** The Red-Green Proof is the mutation the review ran: restoring a bare, wrong-depth citation anywhere in shipped text must make the checker FAIL, naming the file and line. Today it passes.
- Every instance above is either fixed or reported with a reason.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED:** in the current tree, revert `skills/do-work/actions/commit.md:17` to the bare, wrong-depth form `../do-work-toolbox/actions/inspect.md` (no backticks) and run `bash _dev/tests/shipped-package-reference-contract.sh`. Observed today: `shipped package reference contract: PASS`, exit 0 — the blind spot.

**GREEN:** the same mutation FAILs, naming `skills/do-work/actions/commit.md:17`, and the unmutated tree still passes.

## Open Questions

- [ ] REQ-259's review found that the cross-package citation class is bounded by punctuation rather than by what a citation is, with six live escapees across three packages, and demonstrated by mutation that the pre-REQ-259 defect still passes the checker. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — accept that bare-text citations are checked by review rather than mechanically.
