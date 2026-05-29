# Slop Check

Validates a human-facing artifact against the seven anti-slop principles before it ships. Read-only by default — flags findings, offers a rewrite, never auto-applies.

> **Not to be confused with code-review or ui-review.** `do-work code-review` reviews source code for consistency, security, performance. `do-work ui-review` validates UI quality against design best practices. `do-work review work` is a REQ-scoped acceptance check that runs as part of `do-work run`. Slop-check is the anti-slop guardrail for *prose* — briefs, reports, summaries, drafts.

## The seven principles checked

The source of truth is `crew-members/anti-slop.md`; slop-check loads it and runs each principle as an explicit check:

| # | Principle | What it means |
|---|-----------|---------------|
| 1 | **Don't send what you wouldn't read** | If you wouldn't want this in your own inbox, don't send it. Length isn't a substitute for thought. |
| 2 | **Do the verification yourself** | Every claim, every citation, every code snippet — checked before sending. AI drafts pass the cost on; you absorb it. |
| 3 | **Compress before sending** | Cut bullets, headers, throat-clearing, hedging. Taking a 1,000-word draft to 200 forces clarity. |
| 4 | **Lead with the conclusion** | First sentence carries the verdict. Reasoning follows for those who want it. |
| 5 | **Be honest about what's AI-generated and unchecked** | "First-pass draft, numbers not fact-checked" lets the reader calibrate. Pretending unreviewed output is your own erodes trust. |
| 6 | **Ask whether the artifact needs to exist at all** | Most "let me write this up" instincts could be a two-line answer. Default to less. |
| 7 | **Match the medium to the stakes** | Quick question → quick answer. Real decision → real thinking, which usually means *less* AI scaffolding, not more. |

Each check produces PASS or FLAG with one-line evidence. Principles 2 and 5 can come back N/A — slop-check can't verify upstream claims for you, and "be honest about what's AI-generated" only applies if the artifact was AI-drafted. N/A is a documented outcome, not a silent skip.

## Output

A markdown report keyed to each principle:

```
# Slop-check: <artifact source>

| # | Principle | Result | Evidence |
|---|-----------|--------|----------|
| 1 | Don't send what you wouldn't read | FLAG | 1,400-word draft for a status update — paragraphs 3–6 are throat-clearing |
| 2 | Do the verification yourself     | N/A  | Slop-check can't verify the cited numbers; flag for self-review |
| 3 | Compress before sending          | FLAG | Section 4 repeats Section 1's point in three different phrasings |
| 4 | Lead with the conclusion         | PASS | Verdict in first sentence |
| 5 | Honest about AI-drafted          | FLAG | No disclosure block; draft reads like considered work |
| 6 | Does this need to exist?          | FLAG | A two-line message would replace the entire artifact |
| 7 | Match medium to stakes            | PASS | Low-stakes status update → short prose is right |
```

The report says what to cut and why. It does not paraphrase the artifact.

## Rewrite mode

After the findings, slop-check offers a rewrite — but never auto-applies it. The flow:

1. Report findings.
2. Ask the user: `Want a rewrite addressing these flags? [yes / specific principles / no]`.
3. If yes, generate the compressed version and **show it for review**.
4. The user copies or edits before adopting. Slop-check does not overwrite the original file unless the user explicitly approves.

This is deliberate: an AI rewriting AI-flagged prose without human review just shifts the slop one layer down. The point is to surface the cuts, not perform them silently.

## Input

```
do-work slop-check do-work/deliverables/UR-003-client-brief.md   File path
do-work slop-check REQ-042                                       Resolve to the REQ's deliverable / summary / review
do-work slop-check UR-003                                        Resolve to the UR's most relevant artifact
do-work slop-check most recent                                   Newest authored artifact under do-work/deliverables/
do-work slop-check                                               Same as "most recent"
do-work slop-check <pasted multi-paragraph prose>                Treat the input as the artifact directly
```

For "most recent", the action skips `.marp.html` (mechanical Marp-CLI exports) and `*-video/` (Remotion source). It only looks at *authored* prose — `.md` and `.single.html` — to find the newest by mtime.

## Key rules

- **Read-only by default.** The original artifact is never modified without explicit user consent on the rewrite.
- **The seven principles come from `crew-members/anti-slop.md`** — don't paraphrase, don't skip any. The crew-member is the source of truth.
- **N/A is a real outcome.** Principles 2 (verify-yourself) and 5 (disclose-AI) can be N/A when the action can't verify or the artifact wasn't AI-drafted. Document the N/A, don't silently skip.
- **Don't run slop-check on the slop-check report itself.** Recursive self-check produces no signal.

## When NOT to use

- Reviewing code → `do-work code-review`.
- Reviewing UI quality → `do-work ui-review`.
- REQ-scoped acceptance check → `do-work review work` (auto-runs in `do-work run`).
- Internal agent status updates — `caveman.md` governs that. Slop-check is for *human-facing* output.
- Commit messages and PR titles — already short by convention. The principles apply implicitly, no check needed.
