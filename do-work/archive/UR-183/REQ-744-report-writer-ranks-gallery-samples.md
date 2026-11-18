---
id: REQ-744
title: 'Report writer ranks gallery samples, borrows from the best, and offers the runners-up as second opinions'
status: completed
commit: bbc1e5d0db0edfb0462b46b51af43b83cd809ac1
heavy_verified_at: 2026-11-18T16:40:41Z
heavy_verified_revision: bbc1e5d0db0edfb0462b46b51af43b83cd809ac1
kb_status: pending
created_at: 2026-11-18T15:29:51Z
user_request: UR-183
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
route: B
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-11-18T16:18:05Z
  basis:
    - Route B
    - 4-file write set
    - 7 acceptance criteria
    - dependency depth 1
    - browser evidence
    - cross-route regression gates
    - full-suite verification
tdd: false
maintenance: false
batch: report-sample-gallery
related: [REQ-743]
depends_on: [REQ-743]
addendum_to: REQ-712
claimed_at: 2026-11-18T16:15:41Z
write_set: [skills/do-work-toolbox/actions/ai-report.md, skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work-toolbox/docs/ai-report-kit/skeleton.html, skills/do-work-toolbox/docs/ai-report-guide.md]
dispatch_at: 2026-11-18T16:22:42Z
builder_handback_at: 2026-11-18T16:30:55Z
integration_at: 2026-11-18T16:32:53Z
review_at: 2026-11-18T16:40:15Z
completed_at: 2026-11-18T16:40:42Z
release_at: 2026-11-18T16:40:42Z
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
- `_dev/primes/lessons-action-files.md` — 8761 tokens. Matches shared report routing and alternate artifact writers.
- `_dev/primes/lessons-shell-commands.md` — 9341 tokens. Matches shell argv and quoting for the displayed second-opinion calls.

Both index rows have partial family coverage and exceed the 2000-token budget. No required lesson entry fits.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Put matching at the shared Step 3 seam and carry the choice and prompts through the design rules, skeleton, guide, and Step 8.
- [x] **[APPLY]:** The merge range changes the four declared files.
- [x] **[UNIFY]:** Reviewed the four-file range and ran `git diff --check`. The action routes ordinary, kind, and revise through matching, the reference defines quoting and preset lookup, the skeleton keeps intentional author-fill markers, and the guide reflects the writer. The builder reports a passing contract lint.

## Full Context
See [UR-183 input](input.md) for complete verbatim input.

## Triage

**Route: B** - Medium

**Reasoning:** The acceptance criteria and four implementation files are specified. Exploration is needed to trace the shared writer step through each creation route and verify the existing gallery and preset contracts.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- Ordinary creation enters Steps 1–8. Both `--kind` forms replace the Step 1 target resolution and Step 5 narrative and share Steps 3–4 and 6–8. Revise step 5 also enters Steps 2–8 after rechecking current evidence. Put the shared ranking at the end of Step 3. Second opinion dispatches before those steps and uses the existing clarity-editor prepare path.
- Step 5 item 9, Step 8, the verification checklist, Report Design Rules, the kit’s `#another-layout` band and the guide all describe the existing layout and teaching lines. Keep the band anchor and teaching choice. In the kit, preserve provenance receipts, the three depths, role colors and judge-visible FILL markers for author input.
- REQ-743 finalized at 0.325.0 with the initial historical upgrade limitation recorded in its handback. The installed gallery list exposes the sample identity, fit text and record path. Use that projection to rank before opening selected Borrow records. Preset resolution checks the bundle first and the shipped directory next. Session prepare accepts preset plus direction and returns the combined prompt.
- Claim-time lessons were matched against the action and the emitted command. The two relevant partial satellites exceed the budget. The action prime's alternate-writer trap directly applies to the kind and revise routes.

*Explored by work action (inline)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — shared Step 3 ranking, Step 5 authoring, Step 8 summary and verification checklist
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — gallery use and choice in Report Design Rules, including the closing-band contract
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modify) — report-head marker and closing band
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — writer selection, recording and follow-up prompts

**Files I will NOT touch:** Gallery and second-opinion CLI, shipped gallery records and presets, presentation CSS and JS, existing published reports.

**Acceptance criteria (restated from REQ):**
- [x] Ordinary creation, proposal, root-cause and revise list and rank from the reader's question and collected evidence before visuals. Second opinion still follows its own editor direction.
- [x] The writer consults only top-ranked Borrow records, adapts moves using its own facts and preserves role colors, receipts and three depths. Simplifying or combining moves is allowed.
- [x] The qualification fixture shows success count rising 90 → 160 and success rate falling 90% → 80%, with two unit-labelled charts spanning the same windows.
- [x] The head records `two-measures-side-by-side` for that fixture, and Step 8 identifies the pick and its reason.
- [x] At most two fitting runner-ups receive rationale and shell-valid commands naming this bundle. The fixture yields at least one such command with a valid preset and separate single-argument direction. Use direction alone when a preset does not fit.
- [x] The label correction records `none`, and an absent runner-up retains one content-specific alternative layout. The existing teaching method and its command remain available in every case.
- [x] The handback prints every closing line. A runner-up line resolves through the existing prepare route with its selected preset and direction.

## Pre-Flight

**Git:** Main at `82a0eefe05f5cefee388151560ec56aeb378d6de`. Clean on entry. Subsequent changes are the preparation record, focused baseline and run evidence.

**Tests baseline:** The focused probe required five named Go tests for gallery listing and preset/prepare compatibility. It launched and exited 0. Advance wrote `do-work/working/baseline.json`.

**Dependencies:** REQ-743 is terminal at suite 0.325.0. Fresh installed consumer listed its four shipped samples and prepared a second opinion with a shipped preset and an explicit direction.

**Canonical gate:** `bash _dev/tests/maintainer-verify.sh` exited 0 on attempt 1 at 16:20:33 UTC (121 s). Request-bound `preflight` and `green-gate` both returned `satisfied`, with the canonical record for the exact argv and main revision above. [Run evidence](../../runs/work-2026-11-18-153357/REQ-744-preparation.md)

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — consult the installed gallery between evidence and visuals and repeat every closing prompt in the handback.
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — define Borrow, report-head metadata, runner-up and preset selection, and command quoting.
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modified) — add a gallery choice marker and a repeatable layout block.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — explain sample selection, recording, and follow-up commands.

**What was done:** The shared writer ranks gallery samples after it collects evidence, applies fitting presentation moves, and records its choice in each new report. The report and summary offer up to two fitting runner-up layouts or one content-specific fallback, plus the teaching prompt.

## Qualification

**Result:** Pass. The request-bound `advance` records returned `satisfied` for qualification and scope drift over `84f371e5..bbc1e5d0`.

The main merge contains exactly the four scoped files. Its shipped tree matches the independently tested candidate `5bdc38b2`. Ordinary creation, proposal, root-cause, and revise enter the same evidence-to-visuals seam. Second opinion uses the existing dispatcher. The kit placeholders are required author input. `git diff --check` found no errors.

## Testing

**Tests run:** `bash do-work/runs/work-2026-11-18-153357/REQ-744-focused-probe.sh` through request-bound `advance`. `bash _dev/tests/maintainer-verify.sh` through the timing wrapper at `bbc1e5d0`.

**Result:** Five named gallery and second-opinion compatibility tests passed. The canonical gate passed on its first execution in 121 seconds, including both uncached Go stages (428 board and 995 CLI tests), shipped package reference checks and shell lint. The request-bound scope, focused baseline comparison, and canonical gate records all returned `satisfied`.

**Behavioral baseline:** The installed 0.325.0 writer had no gallery ranking seam or choice marker. The reduced count/rate fixture had zero preset-plus-direction runner-ups. Independent reports supply the end-to-end acceptance evidence for this change.

**Heavy verification plan:**
- Range: `84f371e577923ed1e156581bf1f8886cf01d5432..bbc1e5d0db0edfb0462b46b51af43b83cd809ac1`.
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four paths matched the shipped skills subtree.

*Verified by work action*

## Review

**Overall: 97%** | 2026-11-18T16:39:56Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:**
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css:17` still describes one suggested layout. The comment is outside the implementation scope and has no effect on the generated page or the authoritative instructions — impact-negligible → report only

**Acceptance:** Pass — the four-file merge, repo gates, and three independent disposable installs verified the ordinary, proposal, root-cause, revise and second-opinion prepare routes. The comparison had 90 → 160 successes and 90% → 80% success rate in aligned, unit-labelled charts. Its runner-up resolved with the shipped `80-20` preset and exact direction. The text-only revision reranked to `none`, retained the prior bytes, and had only the expected prior-bundle judge link finding. The prepared comparison-evidence revision scenario was not separately run.

**Restatement sweep:** redefined the shared evidence-to-visuals match, the `ai-report-gallery-sample` head marker, and closing layout cardinality and emitted commands. Wave 2 has one member, REQ-744. The writer, kind and revise paths, guide, template, kit and existing CLI consumers were checked.

**Suggested testing:** 0 items.

**Follow-ups created:** None (1 finding report only).

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Put gallery ranking in the shared evidence-to-visuals step. Independent disposable consumers then exercised the actual installed instructions and checked shell arguments, evidence and rendering.

**What didn't:** Reading the changed action and guide alone missed the old one-layout rule repeated in the kit CSS header. Reduced HTML fixtures established shape and launcher compatibility but could not show how a writer chose from real evidence.

**Worth knowing:** Revisions recheck the evidence and enter the shared writer. The judge serves the revised bundle alone and can report the expected link to its sibling even when the link resolves on disk.

## Orientation

[MAP CHANGED] The AI report writer now uses the gallery between evidence collection and visual design. The report records the selected sample and offers other fitting layouts as second-opinion prompts. The existing editor still handles those prompts. The action prime still routes to the existing action-file lessons.

## Discovered Tasks

- The copied kit CSS header still describes one alternate layout. It is a source comment with no rendering or command effect. The action, reference, skeleton and guide define the current limit — impact-negligible → report only

## Heavy Verification Plan

Target range `84f371e577923ed1e156581bf1f8886cf01d5432..bbc1e5d0db0edfb0462b46b51af43b83cd809ac1`.

- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. The four merged paths match the shipped skills subtree.

## Heavy Verification Result

Target and execution revision `bbc1e5d0db0edfb0462b46b51af43b83cd809ac1`.

- `staged-skills` — executed in 31 seconds, exit 0. The runtime fingerprint could not confirm reuse.

## Timing

Observed 2026-11-18T16:22:42Z to 2026-11-18T16:35:59Z: 13m 17s total, 10m 20s attributed across 2 events, 2m 57s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 8m 19s | 1 |
| verification-gate | 2m 01s | 1 |

Slowest stage: builder-work / Implement gallery sample selection and runner-up prompts, 8m 19s, outcome success.
Slowest command: verification-gate / merged maintainer gate attempt 1, 2m 01s, exit 0, bash (2 argv tokens).
