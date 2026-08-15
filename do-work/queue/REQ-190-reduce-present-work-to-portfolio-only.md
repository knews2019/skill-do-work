---
id: REQ-190
title: Reduce present-work to portfolio-only behavior
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
related: [REQ-189, REQ-191, REQ-192]
batch: completed-work-presentation-consolidation
write_set: [skills/do-work-toolbox/actions/present-work.md, skills/do-work-toolbox/docs/present-work-guide.md]
---

# Reduce Present-Work to Portfolio-Only Behavior

## What

Make `do-work-toolbox present-work all|portfolio` a portfolio aggregation command only. Remove the competing per-item detail, brief, explainer, and Remotion workflows, and give bare or item-specific invocations compact non-writing guidance instead of silently delegating.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`present-work` currently competes with `ai-report` for detailed completed-work presentation and also generates per-item briefs, Interactive Explainers, and Remotion source. Keeping only portfolio aggregation makes each command’s ownership unambiguous.

## Context

The canonical detailed report moves to `ai-report` under REQ-189. The Remotion workflow moves to explicit `present-video` under REQ-191. This REQ preserves only the cross-project portfolio summary and the prohibition against fabricated metrics.

## Detailed Requirements

- Accept only `do-work-toolbox present-work all` and `do-work-toolbox present-work portfolio` as writing invocations.
- Both accepted forms must aggregate completed work and refresh `do-work/deliverables/portfolio-summary.md`; the guided timestamped snapshot described below is the only conditional additional output.
- Retain and tighten Portfolio Mode.
- Preserve the prohibition against fabricated metrics; describe value qualitatively when no verified metric exists.
- Load prompt-injection guidance before reading archived user content.
- Load anti-slop guidance before drafting the human-facing portfolio summary.
- Handle `completed-with-issues` consistently with the terminal-success contract: accept it as completed, surface its issues honestly, and continue rejecting cancelled or unfinished work.
- For every valid `all|portfolio` invocation, explain in plain language that the canonical `do-work/deliverables/portfolio-summary.md` will be refreshed in place.
- Before writing, load the clear-questions guidance and ask whether to preserve the newly generated current summary as a timestamped snapshot. The workflow must guide the user through this choice; do not require a flag or another command they must remember.
- If the user answers no, refresh only the canonical summary.
- If the user answers yes, or an answer cannot be collected, refresh the canonical summary and save the same generated content as a timestamped snapshot under `do-work/deliverables/portfolio-snapshots/`.
- Never overwrite an existing snapshot. A timestamp collision must select a new unused name or stop without clobbering the existing file.
- Never delete snapshots automatically. Removal requires a later, explicit user-approved cleanup.
- A future REQ or `## Lessons Learned` entry may cite a useful snapshot. Do not back-edit an archived REQ to add a link.
- Bare `do-work-toolbox present-work` must print compact usage and write nothing.
- `do-work-toolbox present-work UR-NNN` and `do-work-toolbox present-work REQ-NNN` must write nothing and print both exact replacements using the supplied ID:
  - detailed report → `do-work-toolbox ai-report <ID>`
  - video walkthrough → `do-work-toolbox present-video <ID>`
- Do not silently delegate an item-specific invocation to another action.
- Do not silently generate a portfolio for an item-specific invocation.
- Do not show the snapshot prompt for bare or item-specific invocations; those paths remain non-writing guidance only.
- Remove Detail Mode.
- Remove Client Brief generation.
- Remove Interactive Explainer generation.
- Remove Remotion generation.
- Remove sibling-link and detail-depth instructions.
- Do not produce per-item briefs, stakeholder-facing detailed HTML, `.single.html` explainers, or video directories.

## Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- Preserve all existing generated artifacts other than the intentional in-place refresh of the canonical `portfolio-summary.md`; do not migrate, overwrite, or delete prior briefs, `.single.html` files, reports, video directories, or timestamped portfolio snapshots.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
- `present-work` must never create a video artifact.

## Dependencies

None. It coordinates with REQ-189 and REQ-191 but can remove its obsolete modes independently. If REQ-191 has not extracted the embedded Remotion specification before this REQ deletes it, the builder may recover that specification from Git history.

## Builder Guidance

Firm intent. Delete the obsolete detail, brief, explainer, link-depth, and video instruction blocks; keep the portfolio path small and explicit. Item-specific migration is messaging only, not delegation. Guide the user through snapshot retention in plain language instead of adding a flag or expecting them to remember a secondary command.

## Open Questions

- [x] How should repeat portfolio runs preserve history? → Refresh the canonical summary in place, ask whether to save the newly generated summary as a timestamped snapshot, and default to saving the snapshot when no answer can be collected.
- [x] Who may remove an unlinked or no-longer-useful snapshot? → Only a later explicit, user-approved cleanup; never automatic deletion.
- [x] Must Remotion extraction finish before the embedded specification is deleted? → No hard dependency. Git history is an acceptable recovery source if deletion happens first.

## Red-Green Proof
**RED prompt/case:** Invoke bare `present-work`, `present-work UR-NNN`, `present-work REQ-NNN`, `present-work all`, and `present-work portfolio` against a fixture archive while recording every path written.
**Why RED now:** Bare and item-specific invocations currently enter Detail Mode and can generate a Markdown client brief, `.single.html` Interactive Explainer, and Remotion source instead of remaining non-writing migration paths.
**GREEN when:** Bare invocation writes nothing and prints compact usage; UR/REQ invocations write nothing and print the exact `ai-report <ID>` and `present-video <ID>` replacements; `all` and `portfolio` explain the canonical refresh and guide the user through snapshot retention; an explicit “no” writes only `do-work/deliverables/portfolio-summary.md`, while an explicit “yes” or an unavailable answer also creates one no-clobber timestamped snapshot; no snapshot is auto-deleted; the action contains no Detail Mode, Client Brief, Interactive Explainer, sibling-link/detail-depth, or Remotion workflow; both terminal-success states are accepted and cancelled or unfinished work is rejected.
**Validation:** User-adjusted during verification for canonical refresh, interactive snapshots, the save-snapshot fallback, explicit cleanup, and independent extraction order; the remaining GREEN acceptance criteria were user-specified and the concrete RED case was inferred during capture.

## Full Context

See `do-work/user-requests/UR-042/input.md` for the complete verbatim request and batch constraints.

---
*Source: attached “do-work capture-request: Consolidate completed-work presentation around ai-repo…” specification.*
