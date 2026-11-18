---
id: UR-184
title: 'Report gallery feedback and guided design rounds'
created_at: 2026-11-18T15:50:47Z
requests: [REQ-745, REQ-746, REQ-747, REQ-748]
word_count: 173
---
## Summary
Use reader feedback on reports and figures to guide later designs, retain their design history, and make any design reusable on demand.

## Extracted Requests

| Request | Outcome |
| --- | --- |
| REQ-745 | Save likes, dislikes and notes on reports and figures |
| REQ-746 | Copy a starter prompt that selects a specific report design |
| REQ-747 | Generate traceable report variations with GPT-6.2-astra in T3 Code |
| REQ-748 | Run two rounds of five reports on one real evidence bundle |

## Batch Constraints
Depends on REQ-744 (gallery matching and second-opinion suggestions). The user specified GPT-6.2-astra through T3 Code for generating the rounds. Actual user marks separate the two rounds. Likes express design preference; a factual or unit error disqualifies a report.

## User clarification
> Let the builder choose a suitable existing bundle and record it

Use that same real bundle for all ten reports.

## Full Verbatim Input
> ```
> Add a like/dislike feedback loop to the report gallery. Depends on REQ-744.
> 1) `just run-report-gallery` opens the gallery page, where I can mark a report or one figure as liked or disliked, with an optional short note. The marks are saved and still there after I close the page.
> 2) An agent command generates the next round of reports with GPT-6.2-astra through T3 Code. It reads my marks, and each new report records which report it came from, which marks it used, and what it changed.
> 3) First run: ten reports in two rounds of five, from one real evidence bundle. I mark the first five before the second five are generated. A factual or unit error disqualifies a report even if I like its design.
> 4) Every design carries a proper `do-work-toolbox ai-report` starter prompt: each shipped sample, each sample I add, and each generated report. The gallery page shows it with a copy control. When I paste it, the report writer uses the named design instead of ranking the samples itself.
> ```
