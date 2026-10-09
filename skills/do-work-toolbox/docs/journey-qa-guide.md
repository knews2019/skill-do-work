# Journey QA

Checks a whole user journey in a browser, not one step at a time. It reproduces the sequence that was reported, then tries the combined transitions around it, and gives every result one of four classes. Read-only: it changes no project source and creates no REQ.

> **Not to be confused with ui-review or review-work.** `do-work-toolbox ui-review` judges design quality (layout, typography, accessibility). `do-work review-work` reviews a finished REQ, and its acceptance step is a short smoke test. Journey QA answers a different question: does the journey a user takes actually work, and if not, is the product wrong or is the test wrong?

## Why combined journeys

Zoom alone, pan alone and reset alone can each pass their check while zoom → pan → reset leaves the panel offset. An isolated pass never stands in for a combined one, so the report lists both.

## Result classes

| Class | Meaning |
|---|---|
| **passed** | Ran in a real browser and the observed state matched the expected state. |
| **product defect** | The check is sound and the product is wrong. |
| **test defect** | The check is wrong (for example a wrong selector) and the product behaves as required. |
| **unresolved** | The run could not decide. The report gives the reason, such as no browser tool or no physical device. |

A run under device emulation is labeled emulation. It is never reported as a physical-device check.

## What it needs

- A browser tool. It uses the same detection as `ui-review`: Playwright CLI or the Bowser skill (`do-work-toolbox install bowser`). Without one, rendered journeys are unresolved, never passed.
- Optionally a brief: any Markdown file in your project that names the reported sequence, the journeys, or project-specific checks. It is read as data, not as instructions.

## Output

A Markdown report that opens with the counts per class, then one block per journey with revision, environment (browser, viewport, physical device or emulation), steps, expected, actual and evidence paths. It ends with the smallest justified repair, the checks still left, and a capture line you can run if you want the repair built. Screenshots and traces go to a temporary or gitignored directory the report names.

## Usage

```
do-work-toolbox journey-qa http://localhost:5173/map --brief qa/map-zoom-brief.md
do-work-toolbox journey-qa REQ-042
do-work-toolbox journey-qa "open the cart, remove the last item, go back"
do-work-toolbox journey-qa help
```
