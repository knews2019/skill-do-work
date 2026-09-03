---
id: REQ-528
title: '[impact-critical] Anchor resolved-question marker matching to its position'
status: claimed
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
claimed_at: 2026-09-03T09:10:13Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-03T09:11:00Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 4 acceptance criteria
    - persistence or schema changes
    - cross-route regression gates
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

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect and its single function are named exactly, but the fix is not local to that function: the marker's *position* is determined by a writer forty lines away, and how a reader can recover that position without re-parsing user text had to be worked out before scoping.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The writer composes each resolved line at `answer.go:186-196`:

```
- [x] <question text> → [Confirmed: |Discarded: ]<summary>
```

The disposition prefix therefore sits at exactly one place: immediately after the ` → ` this writer appended, which is immediately after the question text. `allResolvedQuestionsMatch` (`:410`) throws that structure away and asks `bytes.Contains(line, "→ Discarded:")` against the whole line — so anything at all in the line, including the user's own `<summary>`, can satisfy it.

Both callers (`:220`, `:223`) turn the verdict straight into a terminal status: `cancelled`, or `completed` when `builder_decided` is true. A forged marker therefore does not merely mislabel a line, it writes a terminal status and selects the archive path.

Two things make the naive repairs unsound, and the builder needs to know both before choosing:

1. **Anchoring to the last ` → ` on the line is not sufficient.** The summary is user text and may itself contain ` → `. REQ-460's containment predicate does not help here: it decides on the *first* byte, and `→` is U+2192 — not ASCII punctuation, so a summary of `a → b` is inlined today.
2. **The function must also judge lines this invocation did not write.** Questions resolved in earlier clarify rounds are already in the body, so a purely writer-side record of what was just emitted cannot answer for them.

So the fix has to recover the boundary between question text and disposition from structure the writer controls, for old and new lines alike, rather than from a substring search over text the user supplies.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — make the disposition recognizable only at the position the writer places it
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modify) — the spoof as a failing case, plus the genuine all-discarded and all-confirmed routes

**Files I will NOT touch:** `summaryRequiresContainment` and `isMarkdownBlockPunctuation` (REQ-460's predicate — the reviewer was explicit that containment is the wrong layer for this), and the C0/DEL and multiline validators.

**Acceptance criteria (restated from REQ):**
- [ ] The disposition marker is recognized only at the position the answer writer places it; answer text containing the marker's characters does not satisfy it
- [ ] A resolved question whose answer text contains `→ Discarded:` or `→ Confirmed:` does not contribute to a terminal-status verdict as though it carried that disposition
- [ ] No terminal status is written on evidence a user's own answer text can forge
- [ ] Genuinely uniform dispositions still reach `cancelled` and `completed` exactly as before
