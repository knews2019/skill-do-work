---
id: UR-183
title: 'Report gallery that the report writer matches against, in two requests'
created_at: 2026-11-18T15:29:51Z
requests: [REQ-743, REQ-744]
word_count: 299
---
## Summary
Make reports easier to understand by matching each new report against a growing gallery of samples. This capture replaces UR-182 (the first report-gallery capture, three requests, all cancelled on 2026-11-18). The user reviewed that capture, accepted eight changes and chose two requests instead of three. Lines in square brackets inside the verbatim input are capture labels, not the user's words. The four screenshots from the first message stay in `do-work/archive/UR-182/assets/`.

## Extracted Requests

| Request | Outcome |
| --- | --- |
| REQ-743 | A report gallery: a sample format, four shipped samples, one command that lists samples and builds the gallery page, and a form that adds a sample the user found |
| REQ-744 | The report writer ranks the gallery samples before drafting, borrows from the best one, records its choice, and offers the runners-up as second-opinion prompts |

## Batch Constraints
- The project gallery lives in `docs/report-gallery/`. The user accepted that folder on one condition: it is deployed into a project by `just run-do-work-update`.
- Four rebuilt patterns ship in the toolbox on the first day. The third-party screenshots are references only and never ship.
- A sample supplies a way to present evidence. Each report supplies its own facts.
- The writer uses judgment: it may combine samples, simplify, or decide that no sample fits.
- Second-opinion prompts come after the report is built. Running them is a later act by the user.
- Order: REQ-743 first, then REQ-744.

## Confirmed acceptance example
Save the inflation example, create a report where two measures could be confused, and check that the report adapts the side-by-side comparison and offers preset-based prompts for other views. The user confirmed this example during the UR-182 capture.

## Full Verbatim Input
> ```
> [First message, 2026-11-18, with four screenshots]
> what ideas can you take from here for a better data visualization?
> 
> the goal is to have an ai-report that can be more easily understood in most cases.
> 
> [Request captured as UR-182, now cancelled and replaced by this capture]
> I need a gallery of sample reports that lists all types of reports that we can kind of do. And it should be also a bit more interactive, meaning that if I see a nice report somewhere, I can then add it in. And the important part is that when the deep report is created. First, you check which report from the gallery matches most and then use it on judgment to build up the best one. And after the report is built, also gives me prompts for second opinions on the variations of the current report with the best preset possible. Thank you.
> 
> [Acceptance example proposed during the UR-182 capture; the user answered "Yes, use this example"]
> Save the inflation example, create a report where two measures could be confused, and check that the report adapts the side-by-side comparison and offers preset-based prompts for other views.
> 
> [Follow-up messages]
> before we continue and I want to know what will be the file shape of this feature? where will the sample gallery live? where will the instruction live?
> 
> changed the model, see if you can improve on these REQs/suggestions
> 
> [Answers to three questions about the rewrite]
> Should I rewrite the three queued report-gallery requests with the eight changes I listed? = Two requests (Recommended)
> Where should a project keep the report samples you add? = docs/report-gallery is fine, if it gets deployed into the repositories when using 
> just run-do-work-update
> Which samples should ship inside the toolbox on the first day? = Four rebuilt patterns (Recommended)
> ```
