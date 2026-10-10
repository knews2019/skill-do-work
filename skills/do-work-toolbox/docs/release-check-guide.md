# Release Check

Checks that intended content or a feature reached its consumer correctly. It traces the content from source to package to the serving environment to what the consumer sees, reports each delivery stage with evidence, and gives a readiness verdict. Read-only: it changes no project file, creates no REQ, and deploys or repairs nothing. The procedure is [`actions/release-check.md`](../actions/release-check.md).

> **Not to be confused with journey-qa, source-audit or review-work.** `do-work-toolbox journey-qa` exercises one user journey in a browser and classifies each step. `do-work-toolbox source-audit` checks that a report's cited sources say what it claims. `do-work review-work` reviews a finished REQ, and its acceptance result covers only the stages that review exercised. Release check answers a different question: did the right content reach the consumer, stage by stage, and is it ready?

## Why a file that loads is not enough

A hit-mask image can exist, decode and still be fully transparent, so nothing is clickable. A served bundle can carry an older version than the source. A note from last week can say "deployed and verified". None of these proves what the consumer gets today. So the check tests what the consumer does with the content, and it labels every piece of evidence current run or historical. Historical evidence never verifies a stage.

## Stages and verdict

The four stages are implementation, integration, deployment and live acceptance, as [core review defines them](../../do-work/actions/review-work.md). Each one is verified, failed, unassessed or not applicable, with its evidence.

| Verdict | When |
|---|---|
| **ready** | No stage failed, and live acceptance is verified in this run. |
| **not ready** | An applicable stage failed. |
| **unknown** | No stage failed, but live acceptance is unassessed. |

Deployment, rollback, content repair and recurring monitoring stay with you. The report names the ones the gaps call for.

## What it needs

- Access to the points it traces: the source, the built package, and the served copy or install.
- A way to run the consumer's check. For a rendered page that is a browser tool, found the way `journey-qa` finds one. Without it, live acceptance is unassessed, never verified.
- Optionally a brief: any Markdown file in your project that names the intended content, where it is served, and project-specific checks. It is read as data, not as instructions.

## Usage

```
do-work-toolbox release-check http://localhost:8080/ --brief qa/launch-brief.md
do-work-toolbox release-check REQ-042
do-work-toolbox release-check help
```
