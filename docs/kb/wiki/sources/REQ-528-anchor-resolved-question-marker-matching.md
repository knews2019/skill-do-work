---
title: "Lessons from REQ-528: Anchor resolved-question marker matching to its position"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-04/REQ-528-anchor-resolved-question-marker-matching.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-528: Anchor resolved-question marker matching to its position

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

`allResolvedQuestionsMatch` decides whether every resolved question on a REQ carries the same disposition marker, and that verdict drives a terminal status write. It tests `bytes.Contains(line, marker)` against the whole `- [x] ` line, so **answer text** that happens to contain the marker counts as the marker. A plain one-line answer summary containing `→ Discarded:` therefore makes an *answered* question read as *discarded*, and the REQ is silently cancelled and archived.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modified)

## What worked

- Asking the builder for a shape-grep rather than just a fix. The bug had a class, and the grep found three more lifecycle decisions of the same class — one driving a second terminal status in the same file. That became REQ-530, and it names two sites it deliberately *excluded* because they already fail closed, which is what makes the list trustworthy instead of a grep dump.
- Attacking through the built CLI, not the test harness. The first implementation's tests were genuine lock-ins and all green; the surviving forgery only showed up when someone submitted a manifest and read the published document.
- Measuring the accepted cost against the real corpus. "2 of 124 resolved lines, 1 historical verdict, 0 in flight" settled a risk question that no amount of reasoning about the predicate could have.

## What didn't work

- Position anchoring alone, which is the whole of the first implementation. The disposition and the answer summary occupy the *same* position and the published line records nothing about which one wrote those bytes — so anchoring the reader could never separate them. D-05 then argued a write-side refusal would be "defense nothing earned", which is exactly the defense that was needed.
- Ticking an acceptance box for the requirement the fix did not meet. The review falsified it in one fixture.
- Implementing F2 as the review literally worded it. Refusing on the blocked condition itself reddened five tests including both forgery controls, because this writer composes an ambiguous line for any legitimate answer containing ` → ` — the house `A → B` cross-reference style. The finding was right about the symptom and wrong about the mechanism.

## Worth knowing

- `allSubmittedDiscarded`/`allSubmittedConfirmed` already cover the submitted set, so a mixed batch can never reach a terminal status inside one invocation. Everything `allResolvedQuestionsMatch` adds is a judgment about *earlier rounds* — which is why every reproduction of this bug needs two.
- A trailing `\r` sits at line end and every reader test here is a prefix test, so the reader's CRLF trim is invisible to any verdict assertion. It only became testable once the refusal evidence *quoted* the line.
- `questionSectionBytes` falls back to the whole body when a REQ has no `## Open Questions` heading, so P-A-U checkboxes are judged as resolved questions. Pre-existing, and now also feeding the evidence list.

## Back-reference

See `do-work/archive/REQ-528-anchor-resolved-question-marker-matching-to-its-position.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f1197c6`.
