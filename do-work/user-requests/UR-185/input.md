---
id: UR-185
title: 'Improve the AI Report Gallery by hill climbing'
created_at: 2026-11-18T17:43:03Z
requests: [REQ-749, REQ-750, REQ-751, REQ-752]
word_count: 909
---
## Summary
Improve selection among ten real reports, with a ten-second reader goal. Establish a reproducible baseline, then keep changes only when a fresh blind comparison and seven browser checks support them.

## Extracted Requests

| Request | Outcome |
| --- | --- |
| REQ-749 | Addendum: Build the gallery from ten real report pages |
| REQ-750 | Establish the seven gallery checks and freeze the scored baseline |
| REQ-751 | Run blind Opus comparisons and keep qualifying gallery changes |
| REQ-752 | Compare the frozen endpoints and publish the gallery climb evidence |

## Batch Constraints
Never edit a report's content to improve its preview. Every preview must be an
actual page screenshot stored as an image file. Preserve the presets list and
the “Add a sample” section. Put open decisions in the final report. Support
each result with command output or evaluator output.

Round 0 comes first. Freeze the report identities, scripts, reader needs, private answer key and rubric after establishing the baseline. Use fresh Opus 5.5 judges through T3 Code and record the stop condition. The ten-second target requires timed testing with new readers before it can be reported as measured.

## Existing work to credit
REQ-743 built the gallery generator, sample discovery, presets and Add a sample. REQ-744 added evidence-based sample selection. Commits e143415e and 137d723d made the gallery wide and responsive and added ten named cards, real screenshot assets, clipboard controls and report links. Reuse those pieces where they meet the new checks.

At 137d723d, the curated cards open shipped `skills/do-work-toolbox/docs/report-gallery/*/example.html` through a separate report link. Six adapt repository evidence and four use labeled fictional data. The preview opens an image viewer, the names copy in title case and the card fragments begin with `#sample-`. The earlier Opus reviews and browser checks are useful design history. They did not use the new frozen corpus, answer key or blind rubric, so they provide no starting score for this run.

## Full Verbatim Input
> ```
> Improve the AI Report Gallery by hill climbing
> 
> Goal
> A new reader can pick the right report in ten seconds. Improve the gallery at
> docs/report-gallery/index.html. Count work already completed toward this request.
> 
> Required result
> Show exactly ten worked reports. Each card displays a screenshot of the actual
> report page. Each report has a short, memorable name. Clicking the name copies
> it. Clicking the preview opens the report. The gallery URL with #name scrolls
> directly to that card.
> 
> Round 0 — Establish the baseline
> 1. Verify that playwright-cli can capture a screenshot. If it cannot and no
>    equivalent browser capture tool works, stop and report the failed command.
> 2. Save screenshots of the original gallery before changing it.
> 3. Select ten real reports from ai-reports/. Include at least two examples of
>    each kind: completed work, proposal, root cause, and architecture. If examples
>    are missing, use do-work-toolbox ai-report with real archived REQs or URs.
> 4. Write a capture script that screenshots the top of each report's index.html
>    at a fixed viewport and saves each image as a file next to the gallery.
>    Record which report produced each image.
> 5. Write a gate script that exits 0 only when all seven checks pass:
>    G1  Exactly ten example cards appear.
>    G2  Each card displays the nonempty screenshot captured from the report it opens.
>    G3  Each name is unique, lowercase, two or three hyphenated words, at most
>        24 characters in total.
>    G4  Clicking each name copies that exact name. Verify the clipboard in a browser.
>    G5  Opening the gallery at #name scrolls to that card. Its preview opens the
>        correct report.
>    G6  Every link resolves and every image loads.
>    G7  At 1440 px and 390 px, in both light and dark mode, names and controls
>        remain visible and usable, with no overlap or horizontal overflow.
> 6. Write four concrete reader needs, one per report kind. Record acceptable
>    report choices for each before evaluation. Keep that answer key from the judge.
> 7. Bring the gallery through all seven gates. Save its screenshots and get a
>    baseline score from a fresh Opus 5.5 evaluator through T3 Code. This is the
>    first incumbent. Freeze the ten report identities, capture script, gate
>    script, reader needs, answer key, and rubric for the rest of the run.
> 
> Evaluator
> Use a fresh Opus 5.5 evaluator through T3 Code for each round. Randomly assign
> the incumbent and candidate to A and B. Show matching screenshots at the same
> viewports and themes. Include all ten cards. Give the evaluator the four reader
> needs and have it record its first choice for each version before showing the
> HTML. Then supply both versions' HTML and the rubric. Keep the version order,
> change description, prior scores, and builder's notes hidden from the evaluator.
> 
> Score each line from 0 to 20, for a total of 100:
> R1  Previews — Can a reader distinguish the reports from the images alone?
> R2  Names — Does each name convey the subject and seem easy to recall?
> R3  Choice — How easily can a reader select a report for each of the four needs?
> R4  Range — Do the examples demonstrate distinctly different report kinds and shapes?
> R5  Attention — Does each visible element earn its place? Deduct for filler and
>     repeated text.
> 
> Return the five scores, total, and single weakest point with its page location
> for each version. Record the four first choices against the frozen answer key.
> Claim a measured ten-second result only if a timed test with new readers was run.
> 
> Hill-climbing loop — Repeat up to 12 rounds
> 1. Start from the current incumbent. Take its evaluator's single weakest point.
> 2. Make one candidate change aimed at that weakness. One idea, one commit.
> 3. Run the capture script, save the candidate screenshots, and run the gates.
>    If a gate fails, record which one, revert only the candidate change, and
>    count the round as reverted.
> 4. If the gates pass, run the blind A/B evaluation above.
> 5. Keep the candidate only if this round's judge scores it at least 3 points
>    above the incumbent AND no rubric line drops by more than 2 points. A kept
>    candidate becomes the incumbent. Otherwise revert only the candidate change.
> 6. Append one line to gallery-climb-log.md with the round, change, gate result,
>    both totals, and keep/revert decision. Link the evidence. For a gate failure,
>    record “unscored” rather than inventing totals.
> 7. After two consecutive reverts, use the next round for one larger idea, such
>    as a different card layout. Apply the same gates and keep rule.
> 
> Stop when the incumbent reaches 90/100 with all gates passing and all four
> reader needs matched, when three consecutive rounds are reverted, or when
> 12 rounds are complete. Record which condition ended the run.
> 
> Constraints
> Never edit a report's content to improve its preview. Every preview must be an
> actual page screenshot stored as an image file. Preserve the presets list and
> the “Add a sample” section. Put open decisions in the final report. Support
> each result with command output or evaluator output.
> 
> Final report
> Give one verdict line with start score, end score, rounds run, changes kept,
> and changes reverted. For comparable start and end scores, have one fresh
> Opus 5.5 judge score the frozen baseline and final incumbent together, blind.
> Show the original, baseline, and final screenshots. Include a table of all ten
> names, kinds, and report paths. Link the climb log and gate evidence. List
> open decisions as D1, D2, and so on.
> ```
