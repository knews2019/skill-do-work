---
source_type: req_lesson
req_id: REQ-374
req_path: do-work/archive/UR-074/REQ-374-show-how-long-each-done-card-took.md
date: 2026-08-26
domain: frontend
module: _dev/primes
tags: [frontend, show, long, each]
---

# Lessons from REQ-374: Show how long each done card took

## What the REQ was about

A card in the Kanban board's Recently Done column states when the work finished (`done Aug 26, 12:47 UTC · 9min ago`) but never states how long it took. Add the implementation span to that card: the time from when the builder started the REQ to when it landed in Done with a completed status.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

## What worked

- Extracting the span and its verdict into one Go helper before adding the second reader. The card and the Durations chart cannot disagree about the four-hour ceiling because there is nothing for them to disagree with — the ceiling has one definition and both are readers of it. Shipping the applied verdict rather than the threshold is the same move the board already made for `excludedReason`.
- Checking what the card's completion instant actually *is* before measuring against it. `resolveCompletionTime` falls back to the commit's git date, so the obvious implementation — span from `claimed_at` to the instant already on the line — would have printed durations for exactly the REQs the Durations view excludes, and nobody would have noticed until the two numbers were compared.

## What didn't work

- Pointing the builder at `formatDurationMinutes` because it was "the existing span formatter." It is the *chart's* formatter: `34.0 min`, where the card wants `34m 00s`. Two formatters existed and the vocabulary question — which surface is this? — was the one that decided it, not "is there already a function." A PR reviewer caught the mismatch against the REQ's own acceptance strings while the builder was still running.
- `omitempty` on a `float64` that a presence flag already guards. D-06 added `hasImplementationSpan` precisely because a zero-minute span is real, then left `omitempty` on the value — so the flag said "present" and the value vanished, and the client multiplied `undefined` into `took NaNs`. A presence flag and `omitempty` on the same field are contradictory by construction: one exists to preserve the zero the other deletes.
- Telling the builder "do not edit the REQ file" while the pipeline requires the builder to tick the P-A-U boxes in that file. Qualification failed on unticked boxes that could not have been ticked. The dispatch has to leave the builder the one write the pipeline demands of it, or the orchestrator has to own the audit — which is what happened, and is arguably better, but it should be chosen rather than discovered.

## Worth knowing

- **A test that spans a threshold widely does not test the threshold.** The agreement test used 40 min, 18 h and −3 h against a four-hour ceiling and looked thorough; a second ceiling of six hours passed it silently, because any threshold in the gap classifies all three samples identically. Only a pair straddling the real boundary — derived from the constant, so it moves with it — can catch a second definition. This is the sharper form of the prime's REQ-322 lesson: reading the constant is necessary but not sufficient, the fixture has to land where the constant decides something.
- **A property argued at length in a comment is not a property under test.** Two of this REQ's three review findings were exactly that shape: the code explained why a finished span must not tick and why a zero span must survive, and nothing asserted either. The comment is where a reader looks for intent; it is also where an unpinned invariant hides most comfortably.
- Go's `parseTimestamp` accepts `2006-01-02 15:04:05` and reads it as UTC; V8's `Date.parse` reads the same string as local time — a nine-hour divergence under `TZ=Asia/Tokyo`. Any board feature that re-parses a frontmatter stamp client-side inherits that bug. Format from Go-measured values instead.

## Back-reference

See `do-work/archive/UR-074/REQ-374-show-how-long-each-done-card-took.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5ad1d3d`.
