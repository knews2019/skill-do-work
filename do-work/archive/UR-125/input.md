---
id: UR-125
title: 'Verify findings strip as a slim collapsible band (M4)'
created_at: 2026-09-05T18:19:11Z
requests: [REQ-589]
word_count: 152
---

# Verify Findings Strip as a Slim Collapsible Band (M4)

## Summary

After REQ-588 (the M1 rows, release 0.303.2) shipped, the user looked at the round-one report and said neither the shipped rows nor M1 was visually nice and both were huge, asked for better mock-ups, and said collapsible was fine. A second round (three compact options) was answered with: show all the options in the mock-up, nothing left to imagine, beautiful and professional. A third round, the slim-band gallery, rendered every state of M4, M5 and M6 at board width in light and dark with the real top bar. The user picked M4: closed, the strip is one line naming every subject with a Show button; open, one row per finding with a chevron that opens that row's remedy. This UR captures one addendum REQ against REQ-588.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-589 | Replace the finding rows with the M4 slim band: a one-line closed state naming every subject, an open state with one row per finding and a per-row remedy toggle, the open state remembered per browser, the slim-band styling shared by the gallery's three options |

## Batch Constraints

- Producer and payload are untouched: subject, category, detail, remedy, fixable, skipped probes stay as REQ-579 and REQ-588 left them. The change is the board client only.
- REQ-578's hide-on-Activity rule keeps reading the two host ids; the strip keeps hiding when there is nothing to say.

## Full Verbatim Input

> ```
> [Screenshot 2: the round-one mock-up report, the M0 and M1 frames side by side: the shipped paragraph rows and the M1 rows with the remedy on its own line, four items each taking about 300 px of strip.]
> 
> <- neither of these is visually nice and they are huge, please provide better mocks, it's fine to be colapsible as well
> 
> I want all of the options in the mockup, don't make me imagine what would be, also make it beautiful and professional
> 
> I asked for all the mockups, you are asking me again, why?
> 
> [Screenshot 3: the round-three gallery, M4 section: the live frame closed (one line: warning glyph, VERIFY, "3 findings · 1 probe not checked", the three subjects and "1 probe not checked" each with a weight dot, a Show button at the right), then State 1 closed and State 2 open in light and dark.]
> 
> ok, M4 is good
> ```
