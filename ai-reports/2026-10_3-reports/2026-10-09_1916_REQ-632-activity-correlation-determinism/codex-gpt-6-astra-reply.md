Created [the published report](index-codex-gpt-6-astra.html).

- Added a clear overview, clickable architecture diagram and component cards.
- Added an interactive example drawn from an existing test.
- Separated current behavior from historical decisions and measurements.
- Added corrections, evidence limits and a commit-pinned source map.

Corrections include:

- Merge expansion was already removed; automatic prefix stamping was never added.
- Current file counts and test coverage replace outdated figures.
- “Largest gap” measures recorded events, not proven inactivity.
- Clarified history windows, command-count limits, clock sharing and the prefix survey’s limits.

Open questions remain about current test results, performance, consumer prefix compliance and whether historical attribution was fully preserved.

Model-only preview note: Static checks passed; original inputs are unchanged. Browser security policy rejected the local preview because `file:` URLs are blocked, so visual rendering remains unverified.

Launcher verification: the generated page passed the shared browser render judge at wide/phone widths in light/dark. The verifier then removed source URLs naming an unpublished commit and made the report link relative; post-edit verification is recorded separately in the builder evidence.
