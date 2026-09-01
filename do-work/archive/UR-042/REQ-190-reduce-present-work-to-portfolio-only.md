---
id: REQ-190
title: Reduce present-work to portfolio-only behavior
status: completed
claimed_at: 2026-08-15T16:38:32Z
route: C
completed_at: 2026-08-15T17:10:21Z
commit: c66d11c
kb_status: promoted
kb_entry: REQ-190-reduce-present-work-to-portfolio-only-be.md
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
- [x] **[PLAN]:** Delete the obsolete per-item presentation surface before adding a small argument dispatcher. Bare, UR, REQ, and invalid inputs will stop after non-writing guidance; `all|portfolio` will share one ordered workflow: announce the canonical refresh, load prompt-injection before archive bodies, normalize terminal-success records, load anti-slop before one staged draft, load clear-questions before one retention choice, resolve any snapshot collision before writes, then refresh the canonical file and conditionally publish byte-identical staged bytes without clobbering or cleanup. Rewrite the guide to the same reduced contract. Preserve only a legitimate pointer to the canonical prescribed-shell guide; publication remains prose because a split shell block around an interactive prompt would create unsafe state and clobber risks. RED is the focused source seam: obsolete per-item terms and routes are present, while snapshot-retention and clear-question terms are absent.
- [x] **[APPLY]:** Replaced the two-mode action with a portfolio-only dispatcher and one ordered aggregation/publication workflow; rewrote the guide around the same writing and non-writing paths. No files outside the two declared product files and this P-A-U/decision record were edited by the builder.
- [x] **[UNIFY]:** Reviewed the complete scoped diff and final file contents. Verified `actions/present-work.md` for dispatcher side-effect boundaries, safety order, terminal-success issue handling, one-draft publication, snapshot fallback/collisions, artifact preservation, and the required prescribed-shell pointer; verified `docs/present-work-guide.md` mirrors every public invocation and preservation rule; verified this REQ contains only P-A-U and decision-log additions from the builder. Focused obsolete-term/snapshot/source assertions, `contract-regressions.sh`, `prescribed-shell-canonicalization.sh`, `shipped-package-reference-contract.sh`, `action-shell-blocks.sh`, trailing-whitespace checks, and `git diff --check` all passed. The diff contains no debug artifacts. Maintainer verification was intentionally left to the orchestrator.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** This is a destructive instruction refactor with several distinct invocation paths, an interactive snapshot-retention branch, no-clobber output behavior, archive-safety requirements, and extensive obsolete-mode deletion. Planning and exploration are needed to preserve only the intended portfolio contract and prove every other path is non-writing.

**Planning:** Required

## Plan

1. **RED first, then delete obsolete behavior from `present-work.md`.** Capture the five-invocation source/fixture contract, remove Detail Mode, Client Brief, Interactive Explainer, Remotion, sibling-link, depth-calibration, and per-item presentation blocks, then rebuild a small dispatcher: bare and item-specific paths write nothing; `all|portfolio` run one archive-safe portfolio workflow. The writing path announces the canonical refresh, loads prompt-injection before archive bodies and anti-slop before drafting, accepts both terminal-success states while surfacing issues, asks one clear snapshot-retention question, defaults unavailable answers to saving, resolves a collision-safe timestamped path before writes, refreshes the canonical summary, and optionally writes the same bytes as a snapshot without deleting prior artifacts.
2. **Rewrite `present-work-guide.md` around the reduced public contract.** Document portfolio-only ownership, both valid writing forms, exact bare/UR/REQ guidance, canonical refresh, guided snapshot choice and fallback, collision safety, terminal-success handling, qualitative value, and exclusions for briefs, detailed HTML, `.single.html`, and videos.
3. **Prove the five paths and preservation rules.** GREEN assertions must show bare/UR/REQ have zero writes and exact replacements; `all|portfolio` differ only by snapshot answer; yes/unavailable writes one byte-identical collision-safe snapshot; no writes only the canonical file; pre-existing snapshots/artifacts stay unchanged; terminal-success and issue honesty are preserved; and every obsolete workflow term is absent. Then run focused source assertions, shell/reference suites, contract regressions, `git diff --check`, and canonical maintainer verification.

**Invocation contract:** bare → compact usage/no writes/no prompt; `UR-NNN` and `REQ-NNN` → exact `ai-report` + `present-video` replacements/no writes/no prompt; `all` and `portfolio` → identical portfolio aggregation/canonical refresh/one guided snapshot decision.

**Requirement coverage:** Task 1 covers dispatch, aggregation, safety, terminal success, metric honesty, snapshots, collisions, and artifact preservation. Task 2 covers user-facing migration and command ownership. Task 3 covers the captured RED/GREEN proof and regressions.

**Non-goals:** No shared schema/archive, router/help/tutorial/test, `ai-report`, future `present-video`, `review-work`, prior-artifact, archived-REQ, publishing/hosting/search, MP4, automatic-video, or snapshot-cleanup edits; REQ-191 and REQ-192 own the adjacent migration surfaces.

**Plan validation:** All detailed requirements map to one of three tasks; no orphan task or over-sized task list found.

*Generated by Plan agent*

## Exploration

`present-work.md` currently mixes a large Detail Mode (target resolution, briefs, `.single.html`, sibling links, and embedded Remotion) with a smaller portfolio workflow. The retained nucleus already scans archived URs/legacy REQs, drafts a canonical portfolio summary, and prohibits fabricated metrics, but it lacks explicit archive-safety loading, terminal-success normalization and issue disclosure, canonical-refresh explanation, guided snapshot retention, collision handling, and no-deletion rules.

Delete the entire two-mode/detail surface and depth calibration, then rebuild around a small input dispatcher plus the portfolio nucleus. The guide needs the same removal and a compact public usage matrix. Do not reuse the one-item completed-work reference wholesale: portfolio aggregation intentionally refreshes its canonical file in place, while that reference makes every consumer output immutable.

Safety order is prompt-injection before archived bodies, anti-slop before drafting, and clear-questions before the single snapshot prompt. Snapshot publication should stay a prose contract in this two-file REQ: a shell block would introduce state-across-prompt and check-then-clobber hazards without a helper inside scope. Resolve the complete timestamped name before either write, publish snapshots exclusively/no-clobber, and write the same generated bytes to canonical and snapshot targets.

Existing tests still pin several obsolete full-cycle routing/caller assertions; REQ-192 owns those migrations. This REQ must preserve the legitimate prescribed-shell guide pointer until that test moves, even after deleting the old merge-diff block.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/present-work.md` (modify) — replace detail/brief/explainer/video behavior with one portfolio-only dispatcher and guided snapshot workflow
- `skills/do-work-toolbox/docs/present-work-guide.md` (modify) — document the reduced command contract, terminal-success honesty, canonical refresh, and snapshot decision

**Files I will NOT touch:** `ai-report`, the future `present-video`, the shared completed-work reference, toolbox router/help/tutorial/next-step files, root README, `_dev/tests/`, crew caller inventories, UR/REQ schemas, archive formats, `review-work`, prior generated artifacts, changelog/version files until the serial integration commit.

**Acceptance criteria (restated from REQ):**
- [ ] Only `present-work all|portfolio` write; both run one archive-safe portfolio workflow and refresh the canonical summary
- [ ] Bare input prints compact usage; UR/REQ input prints both exact ID-specific replacements; all three paths write nothing and never prompt
- [ ] Portfolio scans normalize terminal success, surface `completed-with-issues`, reject cancelled/failed/unfinished work, and load prompt-injection before archive bodies
- [ ] Drafting loads anti-slop and uses verified metrics only, falling back to qualitative value
- [ ] Before writing, clear-questions guides one snapshot choice: no → canonical only; yes/unavailable → canonical plus one byte-identical timestamped snapshot
- [ ] Snapshot publication resolves collisions without clobbering and never deletes snapshots or mutates prior artifacts
- [ ] Detail Mode, Client Brief, Interactive Explainer, Remotion, sibling/depth behavior, per-item briefs/HTML/`.single.html`, and video output are removed

## Decisions

- **D-01 — Keep snapshot publication prose-only.** The scoped action has no executable helper and the snapshot choice splits execution across an interactive boundary. The action therefore resolves the final name before writes and requires exclusive no-clobber creation in prose, while retaining only the canonical merge-aware commit guide pointer required by the shipped shell contract.
- **D-02 — Treat unsupported arguments like bare input.** Since only `all|portfolio` may write, an invalid token prints compact usage and stops with the same no-read, no-prompt, no-write boundary; item-shaped canonical IDs receive the two exact migration commands instead.
- **D-03 — Derive one summary byte sequence.** The canonical file and optional snapshot are written from the same retained content rather than drafted or serialized independently, making byte identity an explicit publication invariant.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/present-work.md` (modified) — removes every per-item/detail/brief/explainer/video workflow and defines one side-effect-free dispatcher plus archive-safe portfolio generation and guided snapshot publication
- `skills/do-work-toolbox/docs/present-work-guide.md` (modified) — documents the two writing forms, bare/item migration guidance, terminal-success portfolio evidence, canonical refresh, snapshot fallback, and preservation boundary

**What was done:** Reduced `present-work` from a mixed item/portfolio artifact generator to one cross-project portfolio command. Only `all|portfolio` may read archives, prompt, or write; bare and item-specific invocations now stop after compact guidance, while portfolio output uses one evidence-backed draft and optional byte-identical no-clobber snapshot.

**Requirements addressed:** Portfolio-only command ownership; exact non-writing guidance; archive ingestion safety; terminal-success and issue honesty; anti-slop and metric provenance; guided snapshot choice with unavailable-answer preservation; collision-safe additive snapshots; canonical-only overwrite boundary; and deletion of detail, brief, explainer, sibling/depth, Remotion, HTML, and video behavior.

## Qualification

Passed — 2 project files verified, all 7 acceptance groups traced to substantive diff content, P-A-U confirmed, dispatcher/safety/publication ordering checked, and no hollow workflow, debug artifact, or scope drift found.

## Testing

**Tests run:** focused GREEN source assertions; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/prescribed-shell-canonicalization.sh`; `bash _dev/tests/shipped-package-reference-contract.sh`; `bash _dev/tests/action-shell-blocks.sh`; `bash _dev/tests/maintainer-verify.sh`; `git diff --check`
**Result:** ✓ All passing (all invoked suites and semantic assertions exited 0)

**Red-green validation:**
- Obsolete workflows: ✗ Detail Mode, Client Brief, Interactive Explainer, sibling/depth, and Remotion instruction blocks existed before implementation → ✓ exact obsolete headings/workflows absent after subtraction
- Non-writing dispatch: ✗ bare and item-specific calls entered Detail Mode → ✓ bare/invalid stop at compact usage and UR/REQ stop after both exact replacement commands, with explicit no-read/no-prompt/no-write boundaries
- Portfolio snapshots: ✗ no clear-question, fallback, collision, or preservation contract existed → ✓ `all|portfolio` share one workflow where No publishes canonical-only and Yes/unavailable publishes the same retained bytes to a numeric-suffix, exclusive no-clobber snapshot

**New tests added:** None — REQ-192 owns durable batch contract tests; this REQ used captured source-seam RED/GREEN assertions plus the existing full contract suites.

*Verified by work action*

## Review

**Overall: 82%** | 2026-08-15T17:09:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 88% |
| Test Adequacy | 86% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- Item-specific migration dispatch does not explicitly inherit case-insensitive, numeric-value ID canonicalization, so equivalent tokens can fall through to generic usage — gate: user-visible → REQ-197 expanded with this instance
- Canonical publication is ordered before exclusive snapshot publication, so a late snapshot failure can leave the Yes/unavailable branch with only one promised output — gate: user-visible → REQ-199 created

**Minor findings:** 0
**Acceptance:** Partial — the portfolio-only command, safety gates, and snapshot contract are delivered, with two executable edge corrections queued.
**Suggested testing:** 2 items
**Follow-ups created:** REQ-199; **existing follow-ups expanded:** REQ-197; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Deleting the entire detail-mode span before rebuilding exposed a small, auditable dispatcher and one portfolio workflow instead of leaving hidden per-item fallthroughs.
- Keeping snapshot publication as a prose contract avoided unsafe shell state across an interactive prompt while still making byte identity and no-clobber behavior explicit.

**What didn't:**
- The first dispatcher described canonical-looking UR/REQ tokens without inheriting the suite-wide case-insensitive, numeric-value token grammar.
- Resolving a collision-safe snapshot name was not enough; the first publication wording refreshed the canonical file before exclusive snapshot success, permitting partial completion.

**Worth knowing:** A non-writing migration path is still an ID-taking action and inherits Target ID Resolution. For a branch promising an immutable snapshot plus a mutable canonical file, publish the no-clobber artifact first and atomically replace the mutable target only after success.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

[MAP CHANGED] `present-work` is now a cross-project portfolio command only: bare and item-specific calls are non-writing guidance, while `all|portfolio` refresh one evidence-backed canonical summary and guide optional snapshot preservation. REQ-197 and REQ-199 carry the two review-discovered edge corrections before the batch closes.
