---
id: REQ-744
title: 'Report writer ranks gallery samples, borrows from the best, and offers the runners-up as second opinions'
status: pending
created_at: 2026-11-18T15:29:51Z
user_request: UR-183
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-sample-gallery
related: [REQ-743]
depends_on: [REQ-743]
addendum_to: REQ-712
---
# Report writer ranks gallery samples, borrows from the best, and offers the runners-up as second opinions

## What
Before it drafts a detailed AI report, the writer ranks the gallery samples against the reader's main question and the evidence in hand. It borrows from the best sample, records which one it used, and closes the report with ready-to-run second-opinion lines for the next-best samples, each with the best-fitting preset.

## Why
The reader should be able to understand the main finding in most reports. One ranking serves both wishes of the user: the best match shapes the report, and the runners-up are the variations worth a second opinion.

## Detailed Requirements

### Rank and borrow
- The ranking happens after evidence is collected (`ai-report.md` Step 3) and before visuals are built (Step 4). The writer cannot judge a fit before it knows the evidence.
- The writer runs `ai-report-gallery --list` (REQ-743), ranks by the reader's question and the evidence, and opens only the samples it ranked first. It does not read the whole gallery.
- The writer borrows the moves in the sample's Borrow list. It never borrows the sample's facts or its look. The kit's design rules still apply in full: role colors, a receipt on every claim, three depths.
- The writer may combine two samples or simplify one when that helps the reader.
- "No sample fits" is a valid result. The writer then uses the kit as it does today. A report with no numbers to compare is one example where this is likely.
- This applies to every form that writes a new `index.html`: ordinary creation, `--kind`, and revise. The second-opinion form is not changed.

### Record the choice
- The report's `<head>` carries `<meta name="ai-report-gallery-sample" content="<sample name>">`, or `none`.
- The Step 8 summary names the sample and gives one line on why it fits, or says that none fit.

### Offer the runners-up
- The closing band "Try another layout or teaching method" changes its layout half. Each of the next-best samples that also fits (two at most) becomes one line: `ai-report second-opinion --preset <name> --editor-direction '<the sample's moves as one instruction>' <bundle path>`.
- The preset is the one that best serves what that variation helps the reader see. The writer picks it from the shipped presets and the bundle's own presets. When no preset fits, the line carries `--editor-direction` alone.
- When no runner-up fits, the band keeps today's rule: one layout unlike the page, chosen for this content.
- The teaching-method suggestion stays as REQ-712 built it.
- Every line names this bundle, holds no single quote inside the direction, and runs as pasted. One or two sentences beside each line say what the variation makes clearer.
- Step 8 prints every line. The action never runs them.

### Files
The change belongs in `skills/do-work-toolbox/actions/ai-report.md`, the **Report Design Rules** in `ai-report-reference.md`, the closing band of `docs/ai-report-kit/skeleton.html`, the verification checklist, and the user guide `docs/ai-report-guide.md`.

## Constraints
- A sample supplies a way to present evidence. Each report supplies its own facts.
- Keep the exploration driven by the report's content. Do not add a random style, a random preset or a table of style axes; the user declined those on 2026-11-02.
- No new second-opinion mechanism and no CLI change in this request. `--preset` and `--editor-direction` already combine.

## Red-Green Proof
**RED prompt/case:** With the shipped gallery in place, ask for a detailed report on work where two measures could be confused (the confirmed example: a rate that falls while a level still rises).
**Why RED now:** `ai-report.md` starts from the kit and a fixed narrative. It consults no samples, records no choice, and its closing band offers one free layout idea and one teaching method.
**GREEN when:** The report shows the two measures side by side over the same period, each labelled by what it measures, with the report's own evidence. Its head names `two-measures-side-by-side`. The closing band holds at least one runner-up line with a preset that exists, and the line runs as pasted. A report on work with nothing to compare records `none` and keeps today's closing band.
**Validation:** User confirmed the example during the UR-182 capture; the `none` case is inferred from the review on 2026-11-18.

## Prior Implementation
REQ-712 (choose the best teaching method for the report, commit 0c06eb10) added the teaching-method suggestion beside the layout suggestion. It lives in the **End with another layout and a teaching method** rule of `ai-report-reference.md`, in `ai-report.md` Step 5 item 9 and Step 8, and in the `#another-layout` band of `skeleton.html`. This request extends the layout half of that band and leaves the teaching-method half as it is.

## Dependencies
Needs REQ-743 (the gallery, its four samples and `ai-report-gallery --list`) finished first.

## Builder Guidance
What to rank on, when to rank, what to record and the shape of the lines are fixed. The wording of the rules and where each sentence sits in the two action files are the builder's choice. Add to the existing steps and rules; do not add a new step number unless the text cannot fit an existing one.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens; matches report routing and alternate writers, but the partially slugged satellite cannot be narrowed below the 2000-token budget.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-183/input.md` for complete verbatim input.
