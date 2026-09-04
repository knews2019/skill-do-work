---
title: "Lessons from REQ-245: Name fabricated stamps in the board's future-stamp warnings"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-04/REQ-245-name-fabricated-stamps-in-the-board-s-fu.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-245: Name fabricated stamps in the board's future-stamp warnings

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

The board's future-stamp diagnosis messages name exactly one cause — "likely local wall-clock time stamped with a Z suffix" — but a fully fabricated value is a second, now-observed cause, and the current wording sends that reader to the wrong fix. Reword the diagnosis clauses to name both causes; keep the fix instruction (rewrite with the current UTC instant per the Timestamp rule) unchanged, since it is correct for both.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timestamp_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

## What worked

- **A message that grows a second cause makes every nearby comment a half-truth, and the comments outnumber the strings.** The REQ named three strings and three comments; the real count in `model.go` alone was two strings and four comments, two of which (`:42`, `:202`) the REQ had not listed. The tell is any comment containing "the usual cause", "the signature of", or "the specific" — those phrasings are single-cause assertions wearing a description's clothes. Grep for the *claim shape*, not just the string being edited. (This is REQ-231's lesson — "a change that adds a second cause makes every sentence naming the first one a half-truth" — reappearing in a different file; it is worth treating as a standing check rather than a per-REQ discovery.)
- **The Go string and the JS string are two renderings of one message, and only one of them is in the write set.** `web/board-cards.js` hand-duplicates `model.go`'s warning for the badge tooltip, and nothing in the suite compares them — I only found the drift by generating a board and grepping the output `index.html`, which showed the new clause and the old clause in the same file. A `write_set` scoped to the Go files silently scopes out the copy the user actually reads. Any future REQ touching a board *message* should have `web/` in its write set by default, or the pairing should be made mechanical (the way `futureInstantSkewAllowanceMs` mirrors `futureTimestampSkewAllowance` with a "keep the two in lock-step" comment — the comment exists for the constant but not for the message).
- **A comment can contradict itself in four lines and still read fine.** My first draft of the constant's doc said the old wording sent readers "to a fix that cannot help them" and then, two lines later, that "the remedy is the same either way" — both cannot be true. The old message's failure is the *diagnosis*, not the remedy. Worth re-reading a freshly written doc comment as a whole rather than as sentences.
- **The lane is a name pattern, not a registry — which is why a "you may not touch that file" constraint did not block adding a probe to it.** `TestMaintainerStrictJavaScriptBehaviorLane` re-execs the binary with `-test.run=^TestJavaScriptBehavior`, so membership is a naming convention. Worth knowing before anyone concludes that JS behavior coverage has to live in `generate_test.go`.
- **A shared browser writes into whichever repo root it considers its working directory, and that was the main tree.** The Playwright MCP dropped a console log and a page snapshot into `skill-do-work2/.playwright-mcp/` — the main tree, which builders may not write. It is gitignored and already held 36 files from sibling sessions, so this is not a git-hygiene failure, but it does mean **a builder can write outside its worktree without ever issuing a write**. I removed only my own two files (matched by timestamp) and left the siblings' alone; deleting the directory would have destroyed other agents' evidence. Any brief that says "never write in the main tree" should treat browser tooling as an exception to state explicitly rather than a rule to assume.
- **Prefer the pre-change artifact over memory when reporting a before/after count.** I had the BEFORE numbers from earlier in the session, but rebuilt the tool from `HEAD`'s blobs into `/tmp` and regenerated the board rather than quoting them — the earlier figures came from a differently-seeded fixture, and a table comparing two different fixtures would have been quietly wrong while looking authoritative.

## Back-reference

See `do-work/archive/UR-055/REQ-245-name-fabricated-stamps-in-the-boards-future-stamp-warnings.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `23bad9d`.
