# Pen-and-Paper Exercise to Prepare a Prompt

> A pre-flight thinking exercise the user does **away from the screen** before opening any significant AI session. The agent's job is to hand off the seven questions, get out of the way, then help structure the answers once the user returns.

**Aliases:** `human-prompt`, `pre-flight`, `pen-and-paper`

**When to use:**
- Before any significant AI task — drafting a spec, delegating to an agent, producing a deliverable that matters
- Before running any other `prompt-kit-*` prompt (diagnostic, context doc, spec engineer, constraints, etc.)
- Whenever the user notices they're reacting to AI's framing instead of driving the conversation

**Inputs / flags:**
- `--skip-handoff` — the user has already done the exercise offline and just wants help structuring their notes; skip straight to the debrief
- `--debrief` — same as `--skip-handoff`

---

## Instructions for the executing agent

This is **not a prompt you answer for the user.** It is a handoff. Your job is to give the user the seven questions, get them to step away from the screen, and then help them structure what they bring back. If you start answering the questions yourself, you have defeated the entire purpose of the exercise.

### Why this exists (share with the user in your own words if they ask)

AI is too fluent. If the user shows up with half-formed thinking, the AI's polished, confident framing will overwrite their own before they notice. The times AI's version beats the user's are the times the user already did the thinking and could evaluate output against their own criteria. The times it leads them astray are the times they showed up empty-handed and let a language model fill in the blanks with statistical plausibility.

Pen and paper don't have opinions. A whiteboard doesn't autocomplete a strategy. A voice memo doesn't rephrase intent. Those are features, not limitations.

### Step 1 — Confirm the user is ready to step away

Unless `--skip-handoff` or `--debrief` was passed, open with something like:

> This exercise takes 5–15 minutes and needs to happen **away from this chat window.** Grab a pen and paper, a whiteboard, a napkin, or a voice memo app — anything that doesn't autocomplete your thinking. I'll give you seven questions. Work through them in order, offline. Come back when you have answers (messy is fine — the mess is yours).
>
> Ready? Say "go" and I'll hand over the questions. Or say "why?" if you want the short version of why this matters before you commit ten minutes to it.

Wait for the user to say "go" (or equivalent). If they ask "why?", give the short rationale above in 2–3 sentences and re-offer the handoff. Do **not** start asking the seven questions conversationally — that turns a pen-and-paper exercise into yet another chat.

### Step 2 — Hand over the seven questions

Print the questions as a single block the user can screenshot or copy. Format them so they're easy to read on paper:

---

**The Seven Questions**

Work through these in order. Pen, voice, or whiteboard — not here.

1. **What am I actually trying to accomplish?**
   Not the task — the outcome. One sentence. If you can't get it to one sentence, keep talking it through until you can.

2. **Why does this matter?**
   What happens if it goes well? What happens if you skip it? Separates the high-stakes from the must-merely-exist.

3. **What does "done" look like?**
   Describe the finished thing — format, length, tone, audience, what the audience should feel/do/know afterward. If you can't describe done, you're not ready to delegate.

4. **What does "wrong" look like?**
   The subtle failure mode. What would be polished, technically correct, and still make you say "no, that's not what I meant"? This is the one people skip — and the most important.

5. **What do I already know about this that I haven't written down?**
   Institutional knowledge. Unwritten rules. The things that are obvious to you but wouldn't be obvious to a newcomer.

6. **What are the pieces?**
   Components, subtasks, dependencies. What comes first? What can run independently?

7. **What's the hard part?**
   The one piece that's genuinely difficult versus the pieces that are just effort. Where are the judgment calls? Where are you least certain?

Come back when you have answers. Say "done" and I'll help you structure what you have.

---

### Step 3 — Wait

Do nothing until the user returns with notes. If they send a short message like "back" or "done", proceed to Step 4. If they come back with partial answers and ask for help continuing, nudge them to finish offline rather than finishing in-chat — the whole point is uncontaminated thinking.

**Exception:** if the user explicitly refuses to do it offline ("just run it conversationally" / "I'm driving, ask me the questions"), honor that, but note at the end that the output will be softer than an offline pass would produce.

### Step 4 — Debrief and structure (only after the user returns)

Once the user shares their notes — or if `--skip-handoff`/`--debrief` was passed upfront — do the following. **Do not rewrite or "improve" their thinking.** Transcribe faithfully, then flag gaps.

Produce this artifact:

```
=== PRE-FLIGHT BRIEF ===
Date: [today]
Task: [one-line descriptor the user gave]

1. OUTCOME
[Q1 — one sentence, in the user's voice]

2. STAKES
[Q2 — what changes if this goes well vs. is skipped]

3. DONE LOOKS LIKE
[Q3 — concrete, observable, verifiable]

4. WRONG LOOKS LIKE
[Q4 — the subtle failure mode; quote the user if they had a vivid phrase]

5. INSTITUTIONAL CONTEXT
[Q5 — unwritten rules, background, things the user takes for granted]

6. DECOMPOSITION
[Q6 — pieces and dependencies, in the order the user described them]

7. THE HARD PART
[Q7 — the judgment calls and areas of uncertainty]
```

Then add two short sections:

- **GAPS I NOTICED** — places where their answers were thin, contradictory, or where Q4 (wrong looks like) didn't pair with a matching constraint in Q5/Q6. Phrased as questions, not corrections.
- **WHERE TO USE THIS** — one sentence naming the next action (e.g., "Paste this at the top of your next AI session before asking for X" or "Feed this into `prompt-kit-spec-engineer`"). If the user mentioned another `prompt-kit-*` prompt or a do-work action by name, route them there.

### Step 5 — Close

End with a single line reminding the user this is **their** thinking, not yours — they should evaluate any future AI output against this brief, not against the AI's later paraphrase of it.

## Rules

- **Do not answer the seven questions on the user's behalf.** Ever. Even if they ask you to "just take a first pass." The exercise only works if the thinking is theirs first.
- **Do not polish the user's notes into something that sounds better than what they wrote.** Fluency is what this exercise is inoculating them against. Preserve their voice, including the rough bits.
- **Do not turn the handoff into a conversational interview.** The whole point is to get the user off the screen. If they refuse, honor that but flag the tradeoff.
- **Do not assume missing answers.** If they skipped Q4, leave Q4 empty and flag it in GAPS I NOTICED — don't invent a failure mode for them.
- **Do not load other context documents or prior session notes during Step 2.** The user's offline pass should not be primed by your reading of their files.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll just ask the questions conversationally — it'll be faster" | Hand over the written block and wait | Fluent conversation is exactly the contamination the exercise prevents |
| "Their Q3 answer is vague, I'll fill in what they probably meant" | Transcribe vague, flag in GAPS | If their own definition of done is vague, papering over it with yours is how the 80%-problem starts |
| "They said they're too busy — I'll just synthesize the brief from what they told me earlier" | Offer the abridged 3-question variant (Q1, Q3, Q4) as a fallback, still offline | A shorter offline pass beats a longer in-chat pass; the medium is the point |
| "They came back with half the questions answered — I'll fill in the rest by inferring from Q1" | Give them back the unanswered questions and ask them to finish offline | Inference from Q1 reproduces the AI's framing they're trying to escape |

## Red Flags

- You produced the PRE-FLIGHT BRIEF without the user ever leaving the chat — the exercise did not happen, regardless of what the artifact looks like
- The brief reads more polished than the user's normal writing — you rewrote rather than transcribed
- GAPS I NOTICED is empty — real first-pass thinking always has gaps; an empty gap list means you smoothed them over
- You answered Q4 ("wrong looks like") for the user because they left it blank — this is the single most important question; flag blanks, don't fill them

## Verification Checklist

- [ ] The user was told to step away from the screen before the seven questions were delivered
- [ ] The seven questions were handed over as a single copy/screenshot-able block
- [ ] The agent did not answer any of the seven questions on the user's behalf
- [ ] The PRE-FLIGHT BRIEF preserves the user's voice (no upgrading vague phrasing into polished prose)
- [ ] GAPS I NOTICED is populated with questions, not corrections
- [ ] The final line reminds the user the brief is their thinking, to be used as the yardstick for future AI output
