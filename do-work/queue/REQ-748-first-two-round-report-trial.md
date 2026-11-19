---
id: REQ-748
title: 'Run two rounds of five reports on one real evidence bundle'
status: blocked
blocked_by: 'User saves at least one Like or Dislike on each of the five round-one reports in the live report gallery'
blocked_at: 2026-11-19T11:23:19Z
estimate:
  calculated_at: 2026-11-18T23:55:11Z
  p50_active_minutes: 85
  confidence: medium
  basis:
    - Route B
    - 25-file write set
    - 22 new files
    - 2 subsystems involved
    - 5 acceptance criteria
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-gallery-feedback
related: [REQ-745, REQ-746, REQ-747]
depends_on: [REQ-747]
claimed_at: 2026-11-18T23:53:11Z
---
# Run two rounds of five reports on one real evidence bundle

## What
Carry out the first feedback trial with ten reports, five per round. The builder chooses a suitable existing real evidence bundle and records which one was used.

## Detailed Requirements
Keep the same bundle across both rounds. Deliver the first five for the user to mark, then pause. Generate the second five only after the user has marked the first five. Use the saved feedback and record lineage and changes. Check facts and units against the bundle, and disqualify an erroneous report even if the user liked it.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** Observe five first-round reports and no second-round generation before the user’s marks. After genuine user feedback, observe five second-round reports tied to the same evidence, with traceable changes and correctness decisions.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers gallery routing and report writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.

## Interrupted round-one publication

The coordinator stopped the trial worker on the explicit handoff request. Continue from the five existing page bundles on main. The evidence manifest is `do-work/runs/work-2026-11-18-191126/trial-source/gallery-input.json`. Each bundle contains a design record, the T3 execution receipt and a factual/render review. The original candidate for Warm/cold labels had a correction. The review and diagnostic artifacts are in the local Git directory at `do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-one.lkQZp8/`, including the preserved failed version.

| First-round page | Recorded factual/render verdict | Current digest equals review |
| --- | --- | --- |
| `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-latency-ranking/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-replicate-variability/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html` | pass/pass | True |
| `ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html` | pass/pass | True |

Next integrator: verify the five main-checkout digests and any changed bytes, complete the live gallery check without writing marks, then publish the five links to the user. Preserve the recorded T3 authorship. Do not generate the first five again. Wait until the user has marked the first five before round two. Then generate five children with the same evidence manifest, verify provenance, changed figures, factual and unit verdicts, and finish canonical finalization. Use ordinary direct repository test commands.

## Round-one publication and feedback handoff

Verified all five reports from `160a96b5` in the live gallery. The [round-one handoff](../runs/work-2026-11-18-191126/REQ-748-round-one/README.md) links the five pages and explains how to save marks and resume. The [manifest](../runs/work-2026-11-18-191126/REQ-748-round-one/manifest.json) records the verified source pins, Astra terminal runs, reviewed bundle digests, live links, exact starters, and the original warm-label rejection.

The next wave waits for the user to save at least one current report or figure mark on each of the five cards. Verify the saved entries with the manifest and read-only preparation command in the handoff before clearing the hold. Round two remains pending.


## Reader feedback — gallery comparison

The user gave this feedback while reviewing the gallery on 2026-11-19.

> request-lifecycle, 
> cli-command-map is kind of nice, but samples are nicer.
>
> Also why don't we have preview image for all of them?

Use the samples as the stronger visual reference and retain the favorable comparison to request-lifecycle and cli-command-map. This comment refers to gallery examples. The five first-round report preferences remain to be collected. The coordinator is fixing missing actual page screenshots on generated and older report cards before the next round.
