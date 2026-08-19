---
id: UR-042
title: Consolidate completed-work presentation commands
created_at: 2026-08-15T09:10:53Z
requests: [REQ-189, REQ-190, REQ-191, REQ-192]
word_count: 972
---

# Consolidate Completed-Work Presentation Commands

## Summary

Make `ai-report` the only detailed completed-work HTML report, narrow `present-work` to cross-project portfolio aggregation, and extract Remotion source generation into an explicit `present-video` action. Unify archive-reading and evidence rules, migrate routing and documentation, and lock the command boundaries into contract tests.

## Extracted Requests

| Request | Title | Scope |
| --- | --- | --- |
| REQ-189 | Canonicalize ai-report and shared completed-work evidence contract | Shared archive/evidence rules and detailed visual/non-visual reporting |
| REQ-190 | Reduce present-work to portfolio-only behavior | Portfolio aggregation, safe archive reading, and migration-only handling for item targets |
| REQ-191 | Extract explicit present-video Remotion action | Standalone source-only animated walkthrough generation |
| REQ-192 | Migrate presentation routing, documentation, and contract coverage | Router/help/tutorial/next-step migration, cross-references, tests, and changelog |

## Batch Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- `ai-report` and `present-work` must never create video artifacts, and `present-video` must never render MP4 output.
- Treat `completed` and `completed-with-issues` as terminal-success states; reject cancelled or unfinished work.
- Load prompt-injection guidance before reading archived user content and apply anti-slop rules to human-facing artifacts.
- Preserve all existing generated artifacts; do not migrate, overwrite, or delete prior briefs, `.single.html` files, reports, or video directories.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, automatic video generation, or an `ai-report --with-video` path.

## Full Verbatim Input

do-work capture-request: Consolidate completed-work presentation around ai-report, reduce present-work to portfolio aggregation, and extract Remotion into an explicit standalone present-video action.

Context

ai-report and present-work currently duplicate much of the same completed-work presentation workflow:

- resolving completed UR/REQ targets;
- reading requirements, implementation summaries, reviews, tests, and lessons;
- explaining what changed and how it works;
- producing responsive HTML with light/dark support;
- presenting before/after states, architecture, files, commits, and verification;
- applying anti-slop rules.

ai-report is the newer and stronger detailed-report system. It additionally owns screenshot provenance, SVG annotations, generated-visual provenance, timestamped report bundles, and render-and-judge visual QA.

present-work still owns portfolio aggregation, Markdown client briefs, an overlapping Interactive Explainer, and Remotion generation. This creates two competing HTML contracts and contradictory routing for non-visual work.

Goal

Make ai-report the single canonical detailed report for completed URs and REQs.

Reduce present-work to portfolio-only behavior.

Move Remotion into a separate command that runs only when explicitly requested.

Command ownership

1. do-work-toolbox ai-report [UR|REQ|most recent]

- Canonical detailed presentation for one completed UR or REQ.
- The only action that produces stakeholder-facing detailed HTML.
- Preserve the timestamped `ai-reports/<timestamp>_<slug>/` bundle.
- Support visual and non-visual work.

2. do-work-toolbox present-work all|portfolio

- Portfolio aggregation only.
- Produce `do-work/deliverables/portfolio-summary.md`.
- Do not produce per-item briefs, HTML explainers, or videos.
- Bare `present-work` must print compact usage and write nothing.
- `present-work UR-NNN` or `present-work REQ-NNN` must write nothing and print migration guidance:
  - detailed report → `do-work-toolbox ai-report <ID>`
  - video walkthrough → `do-work-toolbox present-video <ID>`

Do not silently delegate or generate a portfolio for an item-specific invocation.

3. do-work-toolbox present-video [UR|REQ|most recent]

- New standalone Remotion action.
- Run only through an explicit `present-video`, `remotion`, or `video walkthrough` request.
- Produce `do-work/deliverables/<ID>-video/`.
- Generate valid Remotion source only; never render an MP4.
- ai-report and present-work must never generate video.
- Do not add `--with-video` or automatic video behavior to ai-report.

ai-report changes

Absorb the useful detail-level narrative from present-work:

- concise stakeholder verdict;
- what shipped;
- problem and change;
- how it works;
- qualitative value delivered without fabricated metrics;
- key files and commits;
- verification commands;
- lessons and open questions when available.

Retain ai-report’s stronger visual contract:

- real screenshots outrank synthetic visuals;
- SVG callout annotations;
- side-by-side before/after evidence;
- generated-image provenance;
- self-contained timestamped report directory;
- full-page light/dark render-and-judge verification.

Define two explicit evidence modes:

Visual mode:
- screenshots, annotations, and authentic before/after evidence;
- never fabricate a visual “before” state.

Non-visual mode:
- architecture or data-flow diagrams;
- commit and current-code evidence;
- tests and operational verification;
- an explicit statement that UI captures were not expected for this work.

Remove claims that ai-report is only pixel-anchored. Remove redirects that send educational or non-visual detailed reports to present-work.

Do not carry over arbitrary mandatory interactivity from the Interactive Explainer. Add interaction only when it improves comparison or comprehension.

Do not generate a separate Markdown client brief or `.single.html` explainer.

present-work changes

- Remove Detail Mode.
- Remove Client Brief generation.
- Remove Interactive Explainer generation.
- Remove Remotion generation.
- Remove sibling-link and detail-depth instructions.
- Retain and tighten Portfolio Mode.
- Preserve the prohibition against fabricated metrics.
- Load prompt-injection guidance before reading archived user content.
- Handle `completed-with-issues` consistently with the canonical terminal-success contract.

present-video changes

Move the existing Remotion specification into this action and its guide.

Preserve:

- Problem → Solution → Architecture → Value structure;
- proportional generation rules;
- real archived content;
- `registerRoot` requirement;
- no-external-assets rule;
- qualitative value rules.

Use `remotion studio src/Root.tsx` as the package preview command.

Do not:

- background the preview process;
- sleep and assume the server is ready;
- assume port 3000;
- invoke macOS `open`;
- render an MP4.

Skip trivial or genuinely non-visual work with a concise explanation instead of generating filler.

Load anti-slop and prompt-injection guidance before reading archived content.

Shared contract

Create `actions/completed-work-presentation-reference.md`, shared by ai-report and present-video.

It must canonically define:

- terminal-success target resolution;
- handling of `completed-with-issues`;
- safe reading of UR input and REQ bodies;
- required archive fields;
- merge-aware commit inspection;
- current-code inspection;
- anti-slop and prompt-injection loading;
- no-overwrite behavior for existing outputs.

Do not duplicate these rules across the two action files.

Routing and migration

Update toolbox routing, argument hints, help, tutorials, guides, next-step recommendations, caller lists, and cross-references.

Route:

- `showcase`, `visual report`, `proof of work` → ai-report
- `present-work`, `portfolio`, `work portfolio` → portfolio-only present-work
- `present-video`, `remotion`, `video walkthrough` → present-video

Replace completion-flow recommendations that currently suggest:

`present-work UR-NNN`

with:

`ai-report UR-NNN`

Preserve all existing generated artifacts. Do not migrate or delete old briefs, `.single.html` files, reports, or video directories.

Add a changelog entry explaining the command migration.

Acceptance tests

- ai-report is the only action capable of producing detailed HTML.
- UI work retains screenshots, annotations, responsive layout, and light/dark render verification.
- Backend and refactor work succeeds in non-visual evidence mode without fabricated screenshots.
- `present-work all` and `present-work portfolio` produce only the portfolio summary.
- Bare `present-work` writes nothing and prints usage.
- Item-specific present-work invocations write nothing and print exact replacement commands.
- present-video creates a valid Remotion source tree and no MP4.
- ai-report and present-work never create video artifacts.
- All actions accept terminal-success states, including `completed-with-issues`, and reject cancelled or unfinished work.
- All archive-reading actions load prompt-injection guidance.
- Contract tests reject stale Interactive Explainer and detail Client Brief workflows.
- Documentation presents one unambiguous choice:
  - detailed report → ai-report
  - cross-project portfolio → present-work
  - animated walkthrough → present-video

Non-goals

- Do not rename ai-report.
- Do not change UR/REQ schemas or archive formats.
- Do not modify review-work or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
---
*Captured: 2026-08-15T09:10:53Z*
