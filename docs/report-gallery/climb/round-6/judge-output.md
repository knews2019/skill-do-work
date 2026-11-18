Both versions get the same sealed choices. B scores one point higher (79 against 78) because it gives one card a better name. The HTML files differ only in that card: A calls it `tinka-review` and B calls it `report-review-canvas`.

**Sealed choices** (written to `/tmp/gallery-mcq3i3iy/choices.json` before I opened any HTML):

| Need | A | B |
|---|---|---|
| 1. Second opinion on a report, original protected | second-opinion | second-opinion |
| 2. Agent guidance on one machine, alternatives | memory-note-routing | memory-note-routing |
| 3. Durations chart labels still collide | overflow-labels | overflow-labels |
| 4. Add a command to the Go tool | cli-command-map | cli-command-map |

**Rubric scores** (0–20 each):

| | R1 Previews | R2 Names | R3 Choice | R4 Range | R5 Attention | Total |
|---|---|---|---|---|---|---|
| A | 16 | 14 | 15 | 17 | 16 | **78** |
| B | 16 | 15 | 15 | 17 | 16 | **79** |

What drove the scores, the same in both versions unless noted:
- **R1 Previews:** the ten previews look clearly different from each other. But at 390 px the text inside them is too small to read, so the reader relies on the name and the one-line description.
- **R2 Names:** `release-dossier` does not tell you its subject, which is "Board client modules and deployment evidence". A's `tinka-review` only makes sense if you already know what Tinka is. B's `report-review-canvas` says what the report covers, but its description wraps to two lines on a phone.
- **R3 Choice:** needs 3 and 4 are easy. For need 1, the `second-opinion` description ("Guarded report rebuilds across models") does not say that the original is protected. You only learn that from the preview text ("fences writes to a reserved staging file"). For need 2, the description of `memory-note-routing` ("Where each memory note should go") says nothing about alternatives or trade-offs. `checkout-control-repo` sits next to it and competes for that need.
- **R4 Range:** the set has many shapes: a deck, a dossier, a proposal, a diagnosis, an architecture report and a guide. One mismatch is that the `overflow-labels` preview is labeled "PROPOSAL — DESIGN ALTERNATIVES" but sits under the "Diagnose a problem" heading.
- **R5 Attention:** every card repeats a `#` permalink and a copy icon. The HTML also has a stray closing `</div>` at line 205 in A, and the same line exists in B.

**Weakest point for each version:**
- **A:** the `tinka-review` card, the first card in the "Proposals" row (`A/gallery.html:130-136`). Its name and description only make sense if you know what Tinka is. The top of its preview shows a list of second-opinion versions, so the card looks like a second-opinion report and competes with `second-opinion` for need 1.
- **B:** the same card, the first one in "Proposals" (`B/gallery.html:130-136`). The new name fixes the Tinka problem. But the preview still opens with a list of second-opinion versions, and the new description says "compare feedback". Together they pull the card even closer to need 1 and make it harder to tell apart from `second-opinion`.
