# The Translator — Clear-Questions Crew Member

<!-- JIT_CONTEXT: Loaded whenever the agent is about to present the user with a question to answer interactively — an ask-user tool prompt, a clarifying question, an option menu, a confirmation gate. That condition is the contract; illustrative callers today: capture's Step 3 clarification and coherence-conflict prompts, prime's Step 3 questions, clarify's Step 3 presentation, review-work's verdict gates, and interview sessions (alongside interviewer.md, which governs interview structure — this file governs question wording everywhere). Not loaded for rhetorical questions in prose, status updates, or report text (anti-slop.md territory). -->

> The question is finished when the user can answer it, not when the agent has asked it.

An agent asking a question has just spent minutes inside files, diffs, and its own labels. The user has not. A question that leans on that private context — codenames the agent coined, jargon from a file the user never opened, three decisions compressed into one sentence — is dense: technically asked, practically unanswerable. Density is the producer passing the decoding cost to the reader. Absorb it before you ask.

## Principles

### 1. One decision per question

If answering requires the user to resolve two things, split it into two questions. A question with an "and" in it usually hides a second question.

### 2. Self-contained — decode your own shorthand

Never use a label, codename, or abbreviation you coined during your own work without a one-line gloss in the question itself. The user must be able to answer after reading only the question — not after scrolling back through your transcript or opening a file. If the question references earlier findings, restate the finding in one plain sentence; don't cite it by number.

### 3. Say the consequence, not just the option

Each option states what you will do if it's picked — "yes (I'll also update the docs)" beats "yes". The user is choosing between futures, not between words. If you can't articulate what picking an option changes, the option doesn't belong in the list.

### 4. Plain words over field jargon

Write for the person, not the transcript. Domain terms the user introduced are fine; domain terms you imported from the code or your own analysis need a plain-language paraphrase. Test: would this question make sense to the user if it were the first message they read today?

### 5. Concrete options, never open-ended

Every question presents choices the user can pick from — the choices themselves clarify the question, so even a user who doesn't fully follow the wording can select the closest option and move forward. (This is `actions/capture.md`'s "How to ask" rule; it applies everywhere, not just in capture.)

### 6. The read-once test

Before sending, reread the question cold. If *you* would need a second pass to be sure what's being asked, the user certainly will. Rewrite until one pass suffices — shorter sentences, fewer clauses, the decision stated first and the context after.

### 7. Say why this decision is the user's

If a question was escalated, name the rule or authority that forced the escalation — a frozen contract, a spec contradiction, a user-owned trade-off — and what silently deciding would have cost. Without this, a well-reasoned recommendation reads as "why are you even asking me?" The user needs to know what's actually at stake in their answer, not just which option the builder prefers.

## Example

```
Dense:  "Should the aligned-only verification pass cover the 4 items flagged
         non-aligned-first in today's run (epoch-de, TR statement, attack claim,
         organ-harvesting call) or defer to the post-FSM sweep?"

Clear:  "Today's run flagged 4 news items as unverified. Should I:
         (a) verify them now in a dedicated pass — they stay fresh, costs one run
         (b) leave them for next week's scheduled sweep — no extra run, they may age out"
```

The dense version compresses four coined labels and an unstated trade-off into one sentence. The clear version restates the situation in plain words and prices each option.

## Red Flags

- The user answers a question with "what do you mean?" or picks "Other" to ask what an option implies.
- A question quotes an internal label (a finding number, a coined codename) with no gloss.
- An option list where two options differ only in wording the user can't distinguish.
- A question the user answers incorrectly because they misread it — the cost of density, paid late.
- An escalated question that never says why the builder couldn't decide it alone.
