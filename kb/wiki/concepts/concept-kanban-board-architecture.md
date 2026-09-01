---
title: "Kanban Board and UI"
type: concept
topic_cluster: kanban-board-and-ui
sources:
  - raw/processed/2026-09-01/REQ-015-sync-the-deferred-status-between-the-que.md
  - raw/processed/2026-09-01/REQ-016-remove-the-producer-less-severity-frontm.md
  - raw/processed/2026-09-01/REQ-017-just-run-kanban-replace-a-stale-board-se.md
  - raw/processed/2026-09-01/REQ-040-board-overlap-badge-use-path-match-and-d.md
  - raw/processed/2026-09-01/REQ-087-the-board-and-verify-hand-the-user-the-p.md
  - raw/processed/2026-09-01/REQ-089-the-board-drawer-s-copy-button-omits-the.md
  - raw/processed/2026-09-01/REQ-097-assigned-to-advisory-field-schema-line-s.md
  - raw/processed/2026-09-01/REQ-116-normalize-route-at-the-board-s-read-site.md
  - raw/processed/2026-09-01/REQ-117-an-unrecognized-domain-must-leave-a-foot.md
  - raw/processed/2026-09-01/REQ-119-an-off-vocabulary-route-warns-on-the-boa.md
  - raw/processed/2026-09-01/REQ-122-the-by-ur-lens-counts-recently-done-work.md
  - raw/processed/2026-09-01/REQ-134-addendum-make-queue-kanban-atomic-replac.md
  - raw/processed/2026-09-01/REQ-140-stage-the-modular-board-skill.md
  - raw/processed/2026-09-01/REQ-175-align-board-question-preprocessing-with-.md
  - raw/processed/2026-09-01/REQ-184-live-board-origin-checks-have-no-trusted.md
  - raw/processed/2026-09-01/REQ-185-javascript-behavior-probes-can-all-skip-.md
  - raw/processed/2026-09-01/REQ-195-modularize-the-framework-free-queue-boar.md
  - raw/processed/2026-09-01/REQ-200-render-png-file-mentions-as-images.md
  - raw/processed/2026-09-01/REQ-207-render-html-file-mentions-as-folder-awar.md
  - raw/processed/2026-09-01/REQ-266-name-builds-beside-the-js-renderer-s-mea.md
  - raw/processed/2026-09-01/REQ-277-state-the-mark-label-face-constant-s-rea.md
  - raw/processed/2026-09-01/REQ-284-emit-every-verify-finding-from-the-board.md
  - raw/processed/2026-09-01/REQ-285-render-a-verify-findings-strip-on-the-bo.md
  - raw/processed/2026-09-01/REQ-376-raise-the-done-line-s-faint-text-to-read.md
  - raw/processed/2026-09-01/REQ-381-index-cited-ticket-ids-and-let-the-filte.md
  - raw/processed/2026-09-01/REQ-382-expand-ticket-ids-written-as-markdown-li.md
  - raw/processed/2026-09-01/REQ-385-treat-an-underscore-as-a-ticket-id-bound.md
  - raw/processed/2026-09-01/REQ-386-make-the-drawer-and-the-paste-agree-abou.md
  - raw/processed/2026-09-01/REQ-387-keep-a-spliced-title-from-changing-how-t.md
  - raw/processed/2026-09-01/REQ-388-settle-the-last-two-drawer-clipboard-div.md
related:
  - page: concept-duration-estimation-and-breaks
    rel: complements
  - page: entity-queue-kanban
    rel: extends
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Kanban Board and UI

Architectural overview and synthesis for the Kanban Board and UI subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 30 source documents:

- [[REQ-015-sync-the-deferred-status-between-the-que]] — Sync the `deferred` status between the queue-kanban parser and the Schema Read Contract
- [[REQ-016-remove-the-producer-less-severity-frontm]] — Remove the producer-less `severity` frontmatter field from queue-kanban
- [[REQ-017-just-run-kanban-replace-a-stale-board-se]] — `just run-kanban`: replace a stale board server on the port and open the default browser
- [[REQ-040-board-overlap-badge-use-path-match-and-d]] — Board overlap badge: use path.Match and document the glob dialect
- [[REQ-087-the-board-and-verify-hand-the-user-the-p]] — The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing
- [[REQ-089-the-board-drawer-s-copy-button-omits-the]] — The board drawer's Copy button omits the ticket's frontmatter, so the paste carries no status, domain or timestamps
- [[REQ-097-assigned-to-advisory-field-schema-line-s]] — assigned_to advisory field — schema line, scan skip-and-report, board parse (lock-step)
- [[REQ-116-normalize-route-at-the-board-s-read-site]] — Normalize route at the board's read site and correct 0.174.15's board-wide claim
- [[REQ-117-an-unrecognized-domain-must-leave-a-foot]] — An unrecognized domain must leave a footprint on the board, not become general in silence
- [[REQ-119-an-off-vocabulary-route-warns-on-the-boa]] — An off-vocabulary route warns on the board like domain does
- [[REQ-122-the-by-ur-lens-counts-recently-done-work]] — The By UR lens counts recently-done work as active and honors the window buttons
- [[REQ-134-addendum-make-queue-kanban-atomic-replac]] — Addendum: make queue-kanban atomic replacement cross-platform and symlink-safe
- [[REQ-140-stage-the-modular-board-skill]] — Stage the Modular Board Skill
- [[REQ-175-align-board-question-preprocessing-with-]] — Align board question preprocessing with valid Markdown fences
- [[REQ-184-live-board-origin-checks-have-no-trusted]] — Live board origin checks have no trusted Host anchor
- [[REQ-185-javascript-behavior-probes-can-all-skip-]] — JavaScript behavior probes can all skip while the board suite passes
- [[REQ-195-modularize-the-framework-free-queue-boar]] — Modularize the framework-free queue board client
- [[REQ-200-render-png-file-mentions-as-images]] — Render PNG file mentions as images
- [[REQ-207-render-html-file-mentions-as-folder-awar]] — Render HTML file mentions as folder-aware previews
- [[REQ-266-name-builds-beside-the-js-renderer-s-mea]] — Name builds beside the JS renderer's measured face numbers
- [[REQ-277-state-the-mark-label-face-constant-s-rea]] — State the mark-label face constant's real scope at its canonical home
- [[REQ-284-emit-every-verify-finding-from-the-board]] — Emit every verify finding from the board's Go producer
- [[REQ-285-render-a-verify-findings-strip-on-the-bo]] — Render a verify-findings strip on the board
- [[REQ-376-raise-the-done-line-s-faint-text-to-read]] — Raise the done line''s faint text to readable contrast
- [[REQ-381-index-cited-ticket-ids-and-let-the-filte]] — Index cited ticket ids and let the filter box match them
- [[REQ-382-expand-ticket-ids-written-as-markdown-li]] — Expand ticket ids written as Markdown links
- [[REQ-385-treat-an-underscore-as-a-ticket-id-bound]] — Treat an underscore as a ticket-id boundary on both surfaces
- [[REQ-386-make-the-drawer-and-the-paste-agree-abou]] — Make the drawer and the paste agree about a body H1 that restates the title
- [[REQ-387-keep-a-spliced-title-from-changing-how-t]] — Keep a spliced title from changing how the pasted Markdown parses
- [[REQ-388-settle-the-last-two-drawer-clipboard-div]] — Settle the last two drawer/clipboard divergences: fence info strings and ids inside paths

## Cross-References

See related system components and verification gates.

## Related Entities

- [[entity-queue-kanban]]
