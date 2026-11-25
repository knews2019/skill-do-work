# REQ-757 builder brief

Implement the request below in your T3-assigned checkout and branch. Record their absolute path and name in the handback. Read CLAUDE.md, the applicable primes and the implementation crew rules in skills/do-work/actions/work-implementation-reference.md, including required lesson entries. Consult do-work/lessons-index.md and the linked user request. Preserve sibling changes. Only this request is yours. Use the REQ's own validation, capture the baseline before edits, follow red/green where specified, run focused checks, review the complete diff, and commit the implementation with REQ-757 in the subject. Do not publish, merge or modify queue records, shared version files or CHANGELOG.

Write the completed handback to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-757-handback.md in the main checkout. Provide operative branch, SHA, base, file manifest, changed behavior, focused commands and results, baseline/red evidence, decisions and discovered tasks (explicitly say when none), integration seams and any limitations. Append UTC start/end lines to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-757-progress.log for long operations. These two main-tree paths are the only shared writes. No other agent owns your implementation. Full gate and release are serial integration work. No project run-policy file.

---
id: REQ-757
title: 'Addendum: Show repository reports before gallery samples'
status: claimed
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
