I picked the same four reports in both versions. The two pages are nearly identical. A full diff of the HTML found one difference: the overflow-labels card (row 3, first card). A labels it "Proposal / Why duration labels collide and which layouts could fix them". B labels it "Root cause / Why labels collide above the duration ceiling". B scores 72 and A scores 71.

**Sealed choices** (I saw all four screenshots per version before choosing, then wrote `/tmp/gallery-rlu2pnt_/choices.json`; the release script then sealed them)

| Need | A | B |
|---|---|---|
| 1. Second-model take on a report, original protected | second-opinion | second-opinion |
| 2. Agent guidance on one machine, alternatives before approving | memory-note-routing | memory-note-routing |
| 3. Duration labels still collide | overflow-labels | overflow-labels |
| 4. Go tool dispatch and implementation packages | cli-command-map | cli-command-map |

**Scores** (0 to 20 each)

| Rubric line | A | B | Why |
|---|---|---|---|
| R1 Previews | 15 | 15 | Most previews are distinct. phase-by-phase, tinka-review and memory-note-routing share the same "Depth 1 of 3" template, so they look alike at thumbnail size. The second-opinion preview is a dense block of text. |
| R2 Names | 14 | 14 | `release-dossier` names the report's form, not its subject (board client modules). `second-opinion` and `overflow-labels` are vague without the subtitle. |
| R3 Choice | 16 | 17 | All four needs map to one card. B's "above the duration ceiling" matches need 3 (long-running items bunching together) more directly than A's subtitle. |
| R4 Range | 14 | 14 | A has 4 proposals and only 1 root cause, so the mix is lopsided. B is more balanced (3 proposals, 2 root causes), but its new label contradicts the preview. |
| R5 Attention | 12 | 12 | "Report kinds" repeats grid entries. It also says "Create one" for root cause although the grid already has a root-cause example. The 13 preset chips and a "#" link on every card add noise. |
| **Total** | **71** | **72** | |

**Single weakest point**
- **A:** The "Report kinds" section below the grid, on the "Root cause" row. It tells the reader to create a root-cause report, but queue-refills, a root-cause example, sits in the grid above it. The other rows link reports the grid already shows, so the whole section repeats the grid and contradicts it.
- **B:** The overflow-labels card (row 3, first card). The card's kind says "Root cause", but its own preview's banner reads "Proposal — design alternatives, not a completed-work report". The same "Report kinds" contradiction is also present, now with two root-cause examples in the grid.
