---
id: UR-150
title: 'Addendum to REQ-654: a decision report always lists the smallest-change option'
created_at: 2026-10-09T22:20:10Z
requests: []
word_count: 158
---
# Addendum to REQ-654: A Decision Report Always Lists the Smallest-Change Option

## Summary
The maintainer ruled through the ask tool on 2026-10-10 that the proposal and options kinds of `ai-report` (REQ-654, pending in `do-work/queue/`) must always list the smallest-change option, priced like the others, even when the brief did not name it. The REQ-651 report left out the option that was later chosen and shipped (0.305.88). REQ-654 is pending and unclaimed, so this rides as an Addendum section on that file and mints no REQ.

## Folded Requests

- REQ-654 (ai-report-kind-proposal-root-cause-options) — the whole input: the smallest-change-option constraint

## Full Verbatim Input
> ```
> Addendum to REQ-654 (ai-report --kind proposal, root-cause and options): every decision report must include the smallest-change option.
> 
> Context: the REQ-651 report on the board's activity correlation (ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism) listed two options from its brief, "stamp the [REQ-NNN] prefix then delete the scaffold" and "keep the code". The option that was chosen and shipped in 0.305.88, "delete the scaffold and stamp nothing", was not in the report. The maintainer was asked through the ask tool: "Decision reports today list only the options the brief named. Should that become a rule for the queued proposal-report feature (REQ-654)?" and answered "Yes, add it to REQ-654".
> 
> Constraint to add: for the proposal and options kinds, the option list always includes the smallest-change option, usually "delete the mechanism" or "do nothing extra", priced like every other option (benefit, risk, cost), even when the brief or the topic did not name it. The report may recommend against it, but it must be on the list.
> ```
