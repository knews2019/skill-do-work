---
id: REQ-342
title: "[impact-critical] Neutralize user text written into a REQ body"
status: pending
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: security
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-343, REQ-344]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/clarify.md
  - _dev/tests/contract-regressions.sh
---

# Neutralize User Text Written Into a REQ Body

## What

The Canonical answered-question format writes a user's typed answer verbatim into a REQ body, and
nothing states how that text is neutralized first. Arbitrary typed text can therefore forge the
body's own delimiters. State the neutralization rule at the named entry point — `actions/clarify.md`
Step 4 — so every caller that records a user answer obeys it, not just clarify.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Three consequences, all proven on the user's fixture, and none of them announce themselves:

- A typed answer containing a line-leading `- [ ]` becomes a real unchecked open question. Clarify
  Step 5's "any remaining unchecked item wins" then pins that REQ in `pending-answers` forever and
  re-presents the pasted fragment at every session.
- An unclosed code fence swallows every following section into a code block on the board, while the
  actions' prose greps still see those sections — so the board and the pipeline disagree about what
  the REQ says.
- A paste landing above the opening frontmatter fence loses the entire frontmatter. Status, title and
  `user_request` all read empty, so the REQ silently drops out of its UR.

Pasting a code snippet as an answer is ordinary use, not an attack, and an unbalanced fence is the
common accident. That is what makes this `impact-critical` rather than a hardening nicety: the queue
loses a REQ's identity and its UR membership without erroring.

## Context

`actions/clarify.md` Step 4 is the **Canonical answered-question format** — a named entry point that
`actions/work.md` Step 3.5 already cites for a mid-run answer, so a rule stated there is inherited
rather than copied. Verified: the file today contains no escaping, neutralization or fence-balancing
language at all. Its only "verbatim" mentions are about rewriting stored *question* wording for a
cold reader (`:38`, `:216`, `:231`), which is a different concern in the opposite direction.

This REQ was captured immediately after a live instance of the same class in this repo: a clarify
write keyed its regex on the file's first line-leading `- [ ]`, which is the P-A-U `[PLAN]` checkbox
rather than the question, and falsely ticked `[PLAN]` on three REQs. That was an agent's own text,
not the user's, and it was caught and reverted — but it is the same failure mode from the other side:
this format's delimiters are ambiguous, and neither the writer nor the reader is told so.

## Detailed Requirements

- The Canonical answered-question format states how user text is neutralized before it is written.
- The rule covers at minimum: a leading `- [ ]` or `- [x]`, a bare `---` line, a `## ` heading, and an
  unbalanced code fence.
- **State the condition, not a closed list of characters.** The four above are the shapes the fixture
  proved; the rule must be keyed on "text that could be read as one of this file's own delimiters"
  so a fifth shape is covered without an edit (`_dev/primes/prime-shell-commands.md` § Closed
  Enumerations Go Stale).
- Stated **once**, at the named entry point, and cited by callers rather than restated. Grep every
  current citation of the format and confirm each inherits rather than duplicating.
- The neutralization must preserve what the user actually said — a reader must still be able to see
  the answer's real content. A rule that silently drops characters is not acceptable.

## Constraints

- `_dev/primes/prime-action-files.md` governs any action-file change. Read it first.
- Do not weaken clarify Step 5's "any remaining unchecked item wins" rule — it is correct, and the
  defect is that forged text can create an item it then acts on.
- Scope is the writing rule and its lock-in. Repairing REQ files already damaged by this is REQ-343's
  detection half plus a separate repair decision, not this REQ.

## Builder Guidance

**Certainty: firm on the defect and the entry point, open on the neutralization mechanism.** The user
proved all three consequences and named the minimum shapes; how to neutralize them (indent the
pasted block, fence it with a longer fence, prefix a zero-width guard, quote it as a blockquote) is
the builder's call. Prefer whichever keeps the answer legible to a human reading the archived REQ,
because that trail is the reason the answer is written down at all.

Note the tension this sits next to: `actions/capture-reference.md` § REQ Title Convention says the
parser's line-based recovery "is a salvage path with a narrower contract than the parser proper",
while `frontmatter.go:104` says "recovery is the contract here". REQ-344 covers the frontmatter half;
do not resolve that tension here, but do not lean on recovery either.

## Open Questions

None — the user stated the defect, the entry point, and the minimum shapes.

## Red-Green Proof

**RED prompt/case:** Record an answer whose text contains a line-leading `- [ ]` in a REQ that has no
other open questions, then run `do-work clarify`: the REQ is still `pending-answers` and the pasted
fragment is presented as a question. Separately, record an answer containing one unbalanced code
fence and open the board: every section below it renders inside a code block.

**Why RED now:** Nothing in the Canonical answered-question format says the text is neutralized, so
it is written through as-is.

**GREEN when:** Both pastes are recorded with their content intact and neither forges a delimiter —
the REQ resolves out of `pending-answers`, the board renders the following sections normally, and a
semantic instruction test fails if the neutralization rule is removed or narrowed from the named
entry point.

**Validation:** User confirmed — the defect, the entry point and the minimum shapes are stated
verbatim in the input, from a fixture the user ran.

## Assets

None. Reproduced on the user's fixture; the three consequences are recorded in the UR input.

---
*Source: UR-068 — see `do-work/user-requests/UR-068/input.md` for complete verbatim input.*
