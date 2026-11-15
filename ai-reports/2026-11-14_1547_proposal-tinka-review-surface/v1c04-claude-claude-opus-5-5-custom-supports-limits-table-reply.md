Page: [index-v1c04-claude-claude-opus-5-5-custom-supports-limits-table.html](index-v1c04-claude-claude-opus-5-5-custom-supports-limits-table.html) (the launcher publishes it under its final name).

What changed for clarity:

- The decision comes first: "Pilot O2 by hand. Build nothing yet.", the pilot rule in two lines, and the order of the options (Decision → O2, Decision → O1, O2 → O3, O3 → O4).
- One table replaces the prose and the map: one row per option O1 to O4, a Supports column and a Limits column. Each entry reads "From → To: relation words", then names the other box in one line with its source.
- "Limits or blocks" and "gap in the evidence" both sit in the Limits column; a gap is marked "Gap in the evidence".
- The eight relations that touch no option (L1, L6, L9 → live check; L4, L5 → S5; L2, L7 → decision; L3 → S9) sit below the table, grouped by the box they touch. All 37 relations of the v1c03 map appear on the page.
- The page has no script and no zoom. At phone width each option row stacks its name, then Supports, then Limits.

Corrections to the original report (checked in this repository):

- The footer says "The five images in screenshots/". The folder holds four files.
- O1 says "only the direction text survives, as a meta tag". A directed page also shows the direction in a visible block at its top (ai-report.md:87).
- O1 says "The catalog page lists every bundle". ai-reports/catalog.json has no Tinka bundle (grep -c -i tinka: 0).

Open questions:

- The Tinka line references and the live-check outputs could not be checked: the Tinka checkout is not in this repository.
- The folded-text shares of 38 to 46% could not be checked: the measuring script is not in the bundle.
- O1 has no relation in either column. The report relates no kill attempt, seed or limit to O1, so the row shows "None".

Variation slug: supports limits table
