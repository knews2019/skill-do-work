---
id: REQ-757
title: 'Addendum: Show repository reports before gallery samples'
status: claimed
dispatch_at: 2026-11-25T23:12:29Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-25T23:12:22Z
  basis:
  - trivial short-circuit
route: A
created_at: 2026-11-25T22:17:40Z
user_request: UR-189
domain: frontend
prime_files: []
tdd: true
impact: impact-user-visible
effort_estimate: effort-mechanical
addendum_to: REQ-749
maintenance: false
claimed_at: 2026-11-25T23:11:32Z
---
# Show repository reports before gallery samples
## What
Open the gallery on reports from the current repository, newest first. Label the sample designs clearly and place them after all repository report listings. When the repository has no reports, show the samples immediately.

## Prior Implementation

REQ-749 (build the gallery from ten real report pages) added a curated worked-report section when a manifest is present. The shared template in `report_gallery_page.go` currently emits the sample section before generated and older reports, including in repositories without a worked-report manifest.

## Red-Green Proof

**RED prompt/case:** Generate and open a gallery for a repository with two dated reports, shipped samples, and no worked-report manifest.
**Why RED now:** The first grid shows shipped sample cards. The repository reports appear further down under Generated report designs.
**GREEN when:** The newer repository report appears first, the older one follows, and the sample designs follow the repository listings. With no repository reports, samples are the first cards.
**Validation:** Inferred during capture from the screenshot and discussion.

## Assets

[Gallery on opening](../user-requests/UR-189/assets/REQ-757-screenshot-1-gallery-opening.png) shows the gallery in sa2-sentence-aligner2. The first grid is headed “10 report designs” and starts with Before and After, Bottleneck Hunt, Cause Chain, and Changing Mix. Generated designs is a navigation link at the top.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: [Original request](../user-requests/UR-189/input.md)*

## Triage

Route A. The gallery template and before-and-after ordering are identified, including the empty repository case.

## Plan

**Planning not required** — Route A, direct implementation.
