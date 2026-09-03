---
id: REQ-528
title: '[impact-critical] Anchor resolved-question marker matching to its position'
status: pending
created_at: 2026-09-03T03:10:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-460, REQ-413]
sweep: true
sweep_key: answer-line-marker-position-spoofing
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
---

# Anchor Resolved-Question Marker Matching to Its Position

## What

`allResolvedQuestionsMatch` decides whether every resolved question on a REQ carries the same disposition marker, and that verdict drives a terminal status write. It tests `bytes.Contains(line, marker)` against the whole `- [x] ` line, so **answer text** that happens to contain the marker counts as the marker. A plain one-line answer summary containing `→ Discarded:` therefore makes an *answered* question read as *discarded*, and the REQ is silently cancelled and archived.

## Instances

- `skills/do-work/tools/do-work-cli/internal/publication/answer.go:410-421` — `allResolvedQuestionsMatch` matches the marker anywhere in the line rather than at the position the answer writer put it.
- Demonstrated: answer Q1 with the summary `keep it → Discarded: not really`, then discard Q2. The REQ lands `status: cancelled` with `completed_at` set and the terminal archive path selected, despite Q1 having been answered.
- The `→ Confirmed:` variant, with `builder_decided: true`, reaches `status: completed` by the same route.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding F1** — `impact-critical` — from REQ-460's independent review (Approve, 89%), reproduced in a scratch fixture rather than reasoned about.
- **Pre-existing, not a regression.** The old ten-prefix predicate inlined the same text, so REQ-460 neither introduced nor widened this. It surfaced because REQ-460 redefined exactly this contract: `skills/do-work/actions/clarify.md:103` names `- [x]` as one of "this file's own delimiters" precisely so that answer text cannot be read as one, and REQ-460's new doc comment claims completeness over three *Markdown* ingredients while the real contract is broader than Markdown.

## Detailed Requirements

- The disposition marker must be recognized only at the position the answer writer places it, never anywhere in the line. Answer text containing the marker's characters must not satisfy it.
- A resolved question whose answer text contains `→ Discarded:` or `→ Confirmed:` must not contribute to a terminal-status verdict as though it carried that disposition.
- No terminal status (`cancelled`, `completed`) may be written on evidence a user's own answer text can forge.
- Keep the existing verdict for genuinely uniform dispositions; this is a matching fix, not a policy change.

## Constraints

- Do not solve this in `summaryRequiresContainment`. The reviewer was explicit that the predicate is the wrong layer: containment classifies a summary's shape, while this is about where a marker is anchored on a line the writer itself composed.
- Preserve the existing refusal codes and typed results.

## Dependencies

No request prerequisite. Independent of REQ-460, which is already archived.

## Red-Green Proof

**RED prompt/case:** Answer one open question with the one-line summary `keep it → Discarded: not really`, then discard a second open question, and inspect the resulting `status`.
**Why RED now:** `bytes.Contains` on the whole line lets the answer's own text supply the marker, so the answered question is counted as discarded and the REQ is cancelled.
**GREEN when:** That REQ stays non-terminal with Q1 recorded as answered; a genuinely all-discarded REQ still reaches `cancelled`; and the `→ Confirmed:` plus `builder_decided: true` route still reaches `completed` only when every resolved question really carries it.

---
*Source: REQ-460 independent review finding F1.*
