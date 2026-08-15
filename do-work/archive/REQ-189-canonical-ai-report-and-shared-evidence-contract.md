---
id: REQ-189
title: Canonicalize ai-report and the shared completed-work evidence contract
status: completed
claimed_at: 2026-08-15T16:01:10Z
route: C
completed_at: 2026-08-15T16:35:04Z
kb_status: pending
kb_entry:
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
- [x] **[PLAN]:** Read the action-file prime, both linked Lessons, required crew rules, and refactor spec. Replace ai-report's duplicated target/archive/safety/commit/output rules with one shared completed-work reference; rewrite the action and guide around canonical detailed HTML with visual and non-visual evidence modes; retain the existing image, SVG, responsive-layout, prescribed-shell, and render-verification machinery; then prove stale clauses are gone and focused/full contract checks pass.
- [x] **[APPLY]:** Created the shared completed-work presentation contract; rewrote ai-report around canonical detailed HTML with visual and non-visual evidence modes; retained the established visual backend, SVG, responsive layout, and render-verification references; and aligned the user guide. Implementation stayed within the four scoped product files plus this permitted execution-state/decision record.
- [x] **[UNIFY]:** Ran `git diff --stat`, reviewed all four scoped product files plus this permitted REQ record, and separated the owner's pre-existing CHECKPOINT/queue move from the implementation diff. Verified `completed-work-presentation-reference.md` owns resolution/safety/evidence/current-code/no-overwrite; `ai-report.md` contains one shared pointer and both report modes while retaining prescribed-shell and visual QA; `ai-report-reference.md` retains image-backend/SVG/layout machinery with aligned step names; and `ai-report-guide.md` documents the canonical two-mode report and exclusions. Focused semantic assertions, `contract-regressions.sh`, prescribed-shell canonicalization, shipped-reference checks, action shell-block lint/ShellCheck, and `git diff --check` pass; no debug artifacts were introduced. The integration owner requested that this builder leave `maintainer-verify.sh` to orchestration.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** This REQ establishes a new shared completed-work evidence contract and substantially consolidates the detailed-report behavior of multiple action and guide surfaces. Its visual and non-visual modes, safety rules, and TDD contract coverage need coordinated planning and exploration before implementation.

**Planning:** Required

## Plan

1. **Extract the shared completed-work contract.** Create `skills/do-work-toolbox/actions/completed-work-presentation-reference.md`; first delete duplicated safety, target-resolution, archive-field, commit-reading, and output-collision prose from `ai-report.md`, then replace it with one named pointer. The shared reference will own terminal-success resolution for UR/REQ/most-recent targets, honest `completed-with-issues` treatment, safe archive ingestion, required versus optional evidence, merge-aware commit and current-code inspection, anti-slop/prompt-injection loading, and collision-safe no-overwrite publication.
2. **Make `ai-report` the single detailed-report contract.** Update `ai-report.md`, its reference, and guide to remove pixel-only framing and redirects to `present-work`; define visual and non-visual evidence modes; absorb the stakeholder narrative; preserve the timestamped self-contained bundle and stronger visual machinery; and explicitly exclude Markdown briefs, `.single.html`, video/Remotion/MP4 output, `--with-video`, and automatic video behavior.
3. **Prove the contract change.** Before editing, capture RED assertions for the missing shared reference, stale pixel-only/detail-routing language, and duplicated mechanics. After editing, rerun those assertions for GREEN, verify the positive visual/non-visual/shared-contract clauses, run `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/maintainer-verify.sh`, and review the complete scoped diff.

**Requirement coverage:** Task 1 covers shared resolution, archive safety, evidence, and no-overwrite requirements. Task 2 covers canonical detailed HTML, stakeholder narrative, both evidence modes, artifact preservation, and video/brief exclusions. Task 3 covers the captured RED/GREEN proof and regression verification.

**Non-goals:** No `present-work`, future `present-video`, router/help/tutorial/next-step, schema/archive, review-work, publishing/hosting/search, dependency, or prior-artifact mutations; REQ-190 through REQ-192 own the batch migration surfaces.

**Plan validation:** All detailed requirements map to one of the three tasks; no orphan task or over-sized task list found.

*Generated by Plan agent*

## Exploration

The new shared reference does not exist. `ai-report.md` currently owns duplicated safety, target-resolution, archive-extraction, evidence, commit-inspection, and collision language; it also retains pixel-only framing, redirects non-visual work to `present-work`, permits a synthetic prose “before,” and contains stale `.single.html`/Interactive Explainer comparisons. Its report-folder creation happens before an adequate same-minute collision gate.

`ai-report-reference.md` is the established visual-machinery source: preserve its image backend and provenance rules, SVG/layout rules, optional comparison toggle, and bundle format. The guide preserves useful bundle, visual-evidence, render verification, invocation, and safety content but needs the same two-mode/canonical-report rewrite as the action.

Existing sources of truth are the terminal-success/schema contract in `skills/do-work/actions/work-reference.md`, `scripts/show-commit-diff.sh` plus prescribed-shell guidance for merge-aware evidence, and neighboring archive readers for the prompt-injection-before-ingestion ordering. Existing tests pin image-backend isolation, terminal-success guide language, prescribed-shell references, and shipped Markdown-link integrity.

RED probes: the shared file is absent; stale pixel-only/detail-routing phrases and duplicated mechanics remain. GREEN probes: one shared-contract pointer, explicit visual and non-visual modes, an explicit no-UI-captures statement, current-code plus merge-aware evidence, honest `completed-with-issues`, collision-safe no-overwrite output, preserved visual clauses, and explicit no-brief/no-single-HTML/no-video exclusions.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (new) — canonical target, archive-ingestion, evidence, safety, commit/current-code, and no-overwrite contract shared by detailed completed-work presentations
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — canonical detailed-report workflow with explicit visual and non-visual evidence modes
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — keep the heavy report machinery aligned with the action and shared-contract boundary
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — user-facing canonical detailed-report guidance for both evidence modes

**Files I will NOT touch:** `present-work`, the future `present-video`, toolbox routing/help/tutorial/next-step files, `_dev/tests/`, UR/REQ schemas, archive formats, `review-work`, previously generated artifacts, changelog/version files until the serial integration commit.

**Acceptance criteria (restated from REQ):**
- [ ] `ai-report [UR|REQ|most recent]` is the only detailed stakeholder HTML workflow and retains its timestamped self-contained bundle
- [ ] Stakeholder verdict, shipped change, problem/change, operation, qualitative value, files/commits, verification, and optional lessons/questions are covered without fabricated evidence
- [ ] Visual mode preserves authentic screenshots, SVG annotations, real before/after evidence, generated provenance, responsive layout, and light/dark rendered QA
- [ ] Non-visual mode uses architecture/data-flow, merge-aware commit, current-code, test, and operational evidence and states that UI captures were not expected
- [ ] One shared reference owns terminal-success resolution, safe archive reading, required/optional evidence, merge/current-code inspection, anti-slop/prompt-injection loading, and no-overwrite behavior
- [ ] No Markdown brief, `.single.html`, video, Remotion/MP4, `--with-video`, automatic video, publishing, hosting, search, schema, or prior-artifact mutation is introduced

## Decisions

- **D-01 — DECIDE & STATE — Most-recent compatibility:** Keep the existing ID-based meaning: highest-numbered successful archived UR, then highest-numbered successful legacy REQ as fallback. This centralizes current behavior without silently changing target selection semantics.
- **D-02 — DECIDE & STATE — Non-visual bundle shape:** Keep the timestamped self-contained report folder, but create `screenshots/` only when authentic captures exist. An empty screenshot directory would preserve ceremony while contradicting the requirement not to force non-visual work into a screenshot-shaped report.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (new) — defines the shared safety load order, terminal-success target resolution, archive evidence ledger, merge/current-code inspection, evidence honesty, and collision-safe publication contract
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — replaces duplicated archive mechanics and pixel-only routing with one canonical detailed-report workflow spanning visual and non-visual evidence modes
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — aligns visual-generation, SVG, layout, comparison, and bundle-shape machinery with the consolidated action and optional image directories
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — documents the canonical report, stakeholder narrative, two evidence modes, artifact exclusions, and safety/evidence behavior

**What was done:** Extracted shared completed-work resolution and evidence rules into one toolbox reference, narrowed `ai-report` to its report-specific workflow, preserved its stronger visual verification machinery, and added a first-class non-visual evidence path without introducing alternate briefs, explainers, or video behavior.

**Requirements addressed:** Canonical detailed HTML; stakeholder narrative; authentic visual evidence and rendered QA; non-visual architecture/commit/current-code/test evidence; terminal-success and archive-ingestion safety; merge-aware inspection; missing-evidence honesty; no-overwrite publication; and explicit exclusion of brief, `.single.html`, video, publishing, hosting, and search artifacts.

## Qualification

Passed — 4 project files verified, 6 acceptance groups traced to substantive diff content, P-A-U confirmed, shared-contract wiring resolved, and no hollow data path or debug artifact found. The only other dirty paths are this REQ's orchestrator-owned claim/checkpoint lifecycle files.

## Testing

**Tests run:** focused GREEN semantic assertions; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/prescribed-shell-canonicalization.sh`; `bash _dev/tests/shipped-package-reference-contract.sh`; `bash _dev/tests/action-shell-blocks.sh`; `bash _dev/tests/maintainer-verify.sh`; `git diff --check`
**Result:** ✓ All passing (all invoked suites and shell assertions exited 0)

**Red-green validation:**
- Shared reference presence: ✗ `test -s skills/do-work-toolbox/actions/completed-work-presentation-reference.md` before implementation → ✓ non-empty canonical reference after implementation
- Canonical report routing: ✗ action/guide contained pixel-only framing, Interactive Explainer/`.single.html` comparisons, and non-visual redirects to `present-work` → ✓ stale-routing scan clean while visual and non-visual evidence clauses are present
- Shared evidence contract: ✗ `ai-report.md` directly owned target/safety/evidence mechanics → ✓ action has one shared-reference pointer and the new reference owns terminal success, archive safety, merge/current-code evidence, and no-overwrite publication

**New tests added:** None — REQ-192 owns durable batch contract assertions; this REQ used captured before/after semantic probes plus the existing contract suites.

*Verified by work action*

## Review

**Overall: 82%** | 2026-08-15T16:34:08Z

| Dimension | Score |
|-----------|-------|
| Requirements | 96% |
| Code Quality | 84% |
| Test Adequacy | 86% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- Shared target resolution does not yet inherit case-insensitive, numeric-value ID canonicalization, so equivalent inputs such as `req-42` can miss stored `REQ-042` — gate: user-visible → REQ-197 created
- The generated-image shell block creates `generated/` before a backend succeeds, contradicting the new conditional bundle-shape contract on an all-failed run — gate: user-visible → REQ-198 created

**Minor findings:** 1 (focused GREEN commands are summarized but not preserved as a replayable block; report only)
**Acceptance:** Partial — canonical visual/non-visual reporting and shared evidence rules are delivered, with two executable edge seams queued for focused correction.
**Suggested testing:** 2 items
**Follow-ups created:** REQ-197, REQ-198; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Subtracting duplicated target/evidence prose before extracting one shared reference kept the report action focused on presentation-specific judgment.
- Naming visual and non-visual evidence modes explicitly preserved the mature screenshot machinery without forcing backend and refactor work into screenshot-shaped output.

**What didn't:**
- The first shared resolver described exact archive lookup without inheriting the suite-wide ID-token normalization contract; a canonical reference must also inherit every upstream input grammar it consumes.
- The retained image-generation block still created its public directory before success was known, contradicting the newly conditional bundle shape.

**Worth knowing:** When an instruction says an output directory exists only on success, audit the prescribed shell block itself—not just surrounding prose—for eager `mkdir`. Completed-work ID readers inherit `work-reference.md`'s Target ID Resolution contract even when their search locations differ.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

[MAP CHANGED] Detailed completed-work reporting now has one shared archive/evidence contract and one `ai-report` workflow that selects authentic visual evidence or code/commit/test evidence by work type. The command boundary is clearer, while REQ-197 and REQ-198 carry the two review-discovered edge corrections before the UR closes.
