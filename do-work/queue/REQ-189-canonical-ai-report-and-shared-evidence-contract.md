---
id: REQ-189
title: Canonicalize ai-report and the shared completed-work evidence contract
status: pending
created_at: 2026-08-15T09:10:53Z
user_request: UR-042
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: refactor
depends_on: []
maintenance: true
effort_estimate: normal
related: [REQ-190, REQ-191, REQ-192]
batch: completed-work-presentation-consolidation
write_set: [skills/do-work-toolbox/actions/completed-work-presentation-reference.md, skills/do-work-toolbox/actions/ai-report.md, skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work-toolbox/docs/ai-report-guide.md]
---

# Canonicalize Ai-Report and the Shared Completed-Work Evidence Contract

## What

Make `do-work-toolbox ai-report [UR|REQ|most recent]` the single canonical detailed presentation for one completed UR or REQ, for both visual and non-visual work. Extract the archive-reading, evidence, and output-safety rules shared with `present-video` into one completed-work presentation reference instead of duplicating them across action files.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`ai-report` and `present-work` duplicate completed-work resolution, archive reading, narrative explanation, responsive HTML, evidence, and anti-slop rules. `ai-report` is the newer, stronger report system and already owns screenshot provenance, SVG annotation, generated-visual provenance, timestamped bundles, and rendered visual QA; one detailed-report contract removes the competing HTML workflow and the contradictory redirect for non-visual work.

## Context

This REQ owns the shared foundation and the detailed report. REQ-191 consumes the new shared reference from the standalone video action, while REQ-190 removes detail output from `present-work` and REQ-192 migrates the public routing and cross-cutting documentation.

## Detailed Requirements

### Canonical detailed report

- `do-work-toolbox ai-report [UR|REQ|most recent]` is the canonical detailed presentation for one completed UR or REQ.
- It is the only action that produces stakeholder-facing detailed HTML.
- Preserve the timestamped self-contained `ai-reports/<timestamp>_<slug>/` report bundle.
- Support visual UI work, backend work, refactors, infrastructure work, and other non-visual completed work.
- Remove claims that `ai-report` is only pixel-anchored.
- Remove redirects that send educational, stakeholder-facing, backend, refactor, infrastructure, or other non-visual detailed reports to `present-work`.
- Do not generate a separate Markdown client brief or `.single.html` explainer.
- Do not carry over arbitrary mandatory interactivity from the former Interactive Explainer. Add interaction only when it improves comparison or comprehension.
- Do not generate video, add `--with-video`, or introduce any automatic video behavior.

### Stakeholder narrative

Absorb the useful detail-level narrative formerly owned by `present-work`:

- a concise stakeholder verdict;
- what shipped;
- the problem and the change;
- how the result works;
- qualitative value delivered without fabricated metrics;
- key files and commits;
- verification commands;
- lessons and open questions when available.

### Visual evidence mode

- Preserve the stronger visual contract already owned by `ai-report`: real screenshots outrank synthetic visuals, SVG callout annotations, side-by-side before/after evidence, generated-image provenance, responsive layout, self-contained timestamped output, and full-page light/dark render-and-judge verification.
- Use screenshots, annotations, and authentic before/after evidence for visual work.
- Never fabricate a visual “before” state.
- Keep real screenshots physically and narratively distinct from synthetic or generated images.

### Non-visual evidence mode

- Use architecture or data-flow diagrams, merge-aware commit and current-code evidence, tests, and operational verification when UI captures are not part of the work.
- State explicitly that UI captures were not expected for the work.
- Do not fabricate screenshots or force a screenshot-shaped report for non-visual work.

### Shared completed-work contract

- Create `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` and have both `ai-report` and `present-video` point to it.
- Canonically define terminal-success target resolution for UR, REQ, and most-recent targets.
- Accept both `completed` and `completed-with-issues`, surface the latter honestly, and reject cancelled or unfinished work.
- Canonically define safe reading of UR input and REQ bodies, including prompt-injection handling before archived user content is read.
- Define the archive fields required to produce a presentation and how missing optional evidence is handled without invention.
- Preserve the completed-work evidence sweep: read requirements, implementation summaries, reviews, tests, and lessons when present, and state an optional field's absence rather than inventing evidence.
- Canonically define merge-aware commit inspection and current-code inspection.
- Canonically define anti-slop and prompt-injection loading.
- Canonically define no-overwrite behavior for existing output paths.
- Do not duplicate these shared rules across the `ai-report` and `present-video` action files.

## Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- Preserve all existing generated artifacts; do not migrate, overwrite, or delete prior briefs, `.single.html` files, reports, or video directories.
- Do not rename `ai-report`.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
- Load prompt-injection guidance before reading archived user content and apply anti-slop rules to the human-facing report.
- `ai-report` must never create a video artifact.

## Dependencies

None. This REQ creates the shared reference consumed by REQ-191.

## Builder Guidance

Firm intent. Consolidate by deleting duplicated and contradictory report instructions before adding new prose. Keep the existing proven visual machinery and make the evidence mode explicit rather than creating a second report implementation.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Run or inspect `ai-report` for a completed backend/refactor REQ with no UI captures, then compare its result and routing guidance with the current per-item `present-work` detailed explainer path.
**Why RED now:** `ai-report` still describes itself as pixel-anchored and redirects some non-visual or educational detail work to `present-work`, while the two actions duplicate target resolution, archive inspection, stakeholder narrative, and detailed HTML contracts.
**GREEN when:** One `ai-report` action produces the only detailed stakeholder HTML in either visual evidence mode or non-visual evidence mode; its visual path retains screenshots, SVG annotations, responsive light/dark render verification, and authentic before/after rules; its non-visual path succeeds with code, commit, architecture, test, and operational evidence plus an explicit no-UI-captures statement; shared archive and evidence rules live only in the new reference; no brief, `.single.html`, or video artifact is created.
**Validation:** The GREEN acceptance criteria were user-specified; the concrete RED case was inferred during capture.

## Full Context

See `do-work/user-requests/UR-042/input.md` for the complete verbatim request and batch constraints.

---
*Source: attached “do-work capture-request: Consolidate completed-work presentation around ai-repo…” specification.*
