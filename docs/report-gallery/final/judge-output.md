# Gallery assessment: A vs B

Both versions show the same 10 reports with the same previews, names and one-line descriptions. They differ in how the reports are organized. A puts all 10 cards in one grid, gives each card its own kind label, and adds a "Report kinds" table. B splits the cards into four headed groups (Completed work, Proposals, Diagnose a problem, Architecture) and drops both the per-card kind labels and the table.

First choices (saved before the HTML was released, identical for both): `second-opinion`, `memory-note-routing`, `overflow-labels`, `cli-command-map`.

## Scores

| Rubric line | A | B |
|---|---|---|
| R1 Previews | 15 | 15 |
| R2 Names | 16 | 16 |
| R3 Choice | 15 | 18 |
| R4 Range | 15 | 15 |
| R5 Attention | 12 | 16 |
| **Total** | **73** | **80** |

## Notes per rubric line

- **R1 Previews (tie, 15).** The two galleries use the same images. Most previews show a large readable headline, for example "Pilot Tinka by hand…", "Decision: approve this sorted list…" and "Inside do-work-cli". Three previews are harder to tell apart: `second-opinion` is dense small text, `request-lifecycle` shows a dark banner and a table of contents, and `queue-refills` shows "Less work around each request.", which does not hint at a queue. Five of the 10 previews start with a similar pale UI header, so they look alike at thumbnail size.
- **R2 Names (tie, 16).** Most names state their subject (`memory-note-routing`, `cli-command-map`, `overflow-labels`). `tinka-review` means nothing to a reader who does not know what Tinka is. `release-dossier` describes the report's format, not its subject. `overflow-labels` does not mention durations, although the description does.
- **R3 Choice (B 18, A 15).** In B the group headings follow the reader's goal, so each need maps to one group: need 1 to Completed work, need 2 to Proposals, need 3 to Diagnose a problem, need 4 to Architecture. In A the reader has to scan all 10 cards and read the small kind label under each name. A's "Report kinds" table adds a second, different set of links: its Proposal row and Architecture row point to other reports than the cards do. That can pull a reader away from the right card.
- **R4 Range (tie, 15).** The reports cover four kinds, and the shapes really differ: a one-line deck, a dossier with numbers, a file tree, a long serif essay, an engineering guide. Both versions have the same weakness. The `overflow-labels` preview itself says "PROPOSAL — DESIGN ALTERNATIVES, NOT A COMPLETED-WORK REPORT", but A labels it "Root cause" and B places it under "Diagnose a problem". So the kind shown on the page and the kind shown in the preview disagree.
- **R5 Attention (B 16, A 12).** A repeats the kind label on every card ("Completed work" three times in a row, and so on). The "Report kinds" table then lists the same four kinds again, and the header has an extra "Report kinds" nav link. B removes all three. B loses points for rows with only two cards in a three-column grid, which leaves empty space in the "Diagnose a problem" and "Architecture" rows. B's HTML also has a stray `</div>` after the presets section; it does not affect the page.

## Weakest point

- **A:** The "Report kinds" table, on the bottom left below the grid (`#report-kinds`, the "Root cause" row). The row says "Create one with `do-work-toolbox ai-report --kind root-cause …`" as if no example existed, but two Root cause cards (`overflow-labels`, `queue-refills`) sit directly above it. The table contradicts the grid it summarizes. It also repeats the kind labels, and its Architecture link (`…_1740_architecture-report`) points to a different report than the `request-lifecycle` card (`…_1706_architecture-report`).
- **B:** The `overflow-labels` card, first card under "Diagnose a problem" (third group, left column). The group heading promises a diagnosis, but the preview banner reads "PROPOSAL — DESIGN ALTERNATIVES", and the preview title says "audit and alternatives". A reader with need 3 (why long-running labels bunch together) has to trust the description over the preview. A reader with need 2 (alternatives and trade-offs) could reasonably open this card instead of one in Proposals. In B the grouping is the main way readers choose, so this mismatch costs B more than the same mismatch costs A.
