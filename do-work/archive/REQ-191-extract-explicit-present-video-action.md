---
id: REQ-191
title: Extract an explicit standalone present-video action
status: completed
claimed_at: 2026-08-15T17:13:15Z
route: C
completed_at: 2026-08-15T17:44:24Z
commit:
kb_status: pending
kb_entry:
created_at: 2026-08-15T09:10:53Z
user_request: UR-042
domain: frontend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-189]
maintenance: false
effort_estimate: normal
related: [REQ-189, REQ-190, REQ-192]
batch: completed-work-presentation-consolidation
write_set: [skills/do-work-toolbox/actions/present-video.md, skills/do-work-toolbox/docs/present-video-guide.md]
---

# Extract an Explicit Standalone Present-Video Action

## What

Move the existing Remotion specification into `do-work-toolbox present-video [UR|REQ|most recent]` and a dedicated guide. Generate a valid, source-only animated walkthrough only after an explicit video request, with no automatic invocation and no MP4 render path.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Extract the useful historical Remotion contract into two standalone docs, delegate completed-target safety/evidence/collision mechanics to the shared reference, and validate an explicit-only, proportional, source-only package contract with safe foreground preview guidance.
- [x] **[APPLY]:** Added only the scoped standalone action and guide, with shared-reference delegation, the pre-write skip gate, the complete four-scene source contract, and manual foreground preview guidance.
- [x] **[UNIFY]:** Reviewed both new files in full; scoped status showed only those additions; focused positive/negative assertions, action shell-block lint, shipped-reference checks, contract regressions, and `git diff --check` passed with no debug artifacts.

## Why

Video source generation is currently embedded in `present-work` and can run as part of a broader presentation workflow. A standalone action makes the cost and intent explicit and allows `ai-report` and `present-work` to stay video-free.

## Context

REQ-189 creates `actions/completed-work-presentation-reference.md`; this action must consume that reference rather than restating its archive and evidence rules. REQ-192 owns public routing and command discovery after the action and guide exist. REQ-190 may remove the embedded Remotion instructions first; Git history remains an acceptable source for extracting the existing specification.

## Detailed Requirements

- Add `do-work-toolbox present-video [UR|REQ|most recent]` as a standalone action with its own guide.
- Run it only through an explicit `present-video`, `remotion`, or `video walkthrough` request.
- Resolve a completed UR, completed REQ, or most-recent completed target through the shared completed-work presentation reference.
- Accept `completed` and `completed-with-issues`, surface known issues honestly, and reject cancelled or unfinished work.
- Load anti-slop and prompt-injection guidance before reading archived content.
- Produce `do-work/deliverables/<ID>-video/`.
- Generate valid Remotion source only and never render an MP4.
- Move and preserve the useful existing Remotion contract:
  - Problem → Solution → Architecture → Value structure;
  - proportional generation rules so output depth matches the work;
  - real archived content rather than invented feature claims;
  - a module-level `registerRoot` requirement;
  - a no-external-assets rule;
  - qualitative value rules with no fabricated metrics.
- Use `remotion studio src/Root.tsx` as the package preview command.
- Keep the preview in the foreground.
- Do not background the preview process.
- Do not sleep and assume the server is ready.
- Do not assume port 3000.
- Do not invoke macOS `open`.
- Do not add any MP4 render command or render an MP4.
- Skip trivial or genuinely non-visual work with a concise explanation instead of generating filler.
- Honor the shared no-overwrite behavior for an existing output directory.
- Do not duplicate terminal-success resolution, archive-reading safety, required archive fields, merge-aware commit inspection, current-code inspection, crew loading, or no-overwrite rules from the shared reference.

## Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- `present-video` is never called automatically by `ai-report`, `present-work`, or a completion flow.
- Preserve all existing generated artifacts; do not migrate, overwrite, or delete prior briefs, `.single.html` files, reports, or video directories.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
- Use only local/generated source and content allowed by the no-external-assets rule.

## Dependencies

Depends on REQ-189, which creates the shared completed-work presentation reference this action must use. It has no hard dependency on REQ-190; if the old embedded specification has already been removed, recover it from Git history.

## Builder Guidance

Firm intent. Extract the proven Remotion content proportionally and delete its unsafe preview assumptions. A skipped trivial or non-visual target is a successful concise outcome, not a reason to manufacture scenes.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Request a video walkthrough for a completed visual REQ and inspect the generated package, then request one for trivial or genuinely non-visual work; also search the command for backgrounding, sleeps, fixed port 3000, macOS `open`, and MP4 rendering.
**Why RED now:** Remotion generation exists only inside `present-work`, its preview script backgrounds Studio, sleeps, assumes port 3000, and invokes `open`, and there is no explicit standalone routing boundary.
**GREEN when:** Explicit `present-video`, `remotion`, or `video walkthrough` routing can generate a valid `<ID>-video` Remotion source tree with module-level `registerRoot`, real archived content, no external assets, and a foreground `remotion studio src/Root.tsx` preview; no MP4 exists or can be produced by the action; trivial/non-visual work is skipped concisely; no other action invokes it automatically.
**Validation:** The GREEN acceptance criteria were user-specified; the concrete RED case was inferred during capture. The independent REQ-190/REQ-191 ordering was user-confirmed during verification.

## Full Context

See `do-work/user-requests/UR-042/input.md` for the complete verbatim request and batch constraints.

---
*Source: attached “do-work capture-request: Consolidate completed-work presentation around ai-repo…” specification.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This adds a new standalone action and guide, recovers a removed embedded specification from Git history, consumes the new shared evidence contract without duplication, and must define a safe source-only Remotion package plus multiple skip and preview boundaries. Planning and exploration are required before authoring the new public workflow.

**Planning:** Required

## Plan

1. **Author `actions/present-video.md` as an explicit source-only action.** Justify toolbox ownership in the description, gate explicit video intent before archive reads/writes, consume `completed-work-presentation-reference.md` as the sole target/evidence/safety/no-overwrite contract, skip trivial or genuinely non-visual work concisely, and generate a valid private Remotion/React/TypeScript tree with Problem → Solution → Architecture → Value scenes, module-level `registerRoot`, proportional timing, real archived claims, qualitative value, no external assets, and exactly the foreground preview script `remotion studio src/Root.tsx`. Never install, launch, render, or create media automatically.
2. **Write `docs/present-video-guide.md`.** Document explicit accepted forms, never-automatic ownership, eligibility versus successful skip, the four-scene proportional narrative, honest issues/value, source-only output and preservation, and manual preview from the reported directory with Studio kept in the foreground. Exclude backgrounding, sleeps, readiness guesses, fixed ports, browser launching, macOS `open`, and every MP4/render path.
3. **Prove extraction and boundaries.** RED is the absence of both standalone files plus the pre-REQ-190 embedded unsafe preview script. GREEN assertions require the shared-reference pointer, explicit-only and skip gates, exact source tree/scenes, module-level registration, exact preview command, proportional/evidence/no-assets/value clauses, and source-only/no-render boundaries; reject render scripts, backgrounding, sleeps, fixed port 3000, and shell `open`. Run focused assertions, contract/shell/reference suites, canonical maintainer verification, `git diff --check`, and review only the two-file scope.

**Requirement coverage:** Task 1 covers explicit invocation, shared evidence resolution, terminal status/issues, skip gates, output safety, source validity, scene structure, proportionality, no assets, value honesty, and no MP4. Task 2 covers public usage, safe foreground preview, and artifact preservation. Task 3 covers TDD and regressions.

**Non-goals:** No router/SKILL/help/tutorial/next-step/inventory/test, shared-reference, `ai-report`, `present-work`, schema/archive/review-work, publishing/hosting/search, prior-artifact, dependency-install, changelog/version, or automatic invocation edits; REQ-192 owns discovery and test integration, while REQ-197 owns shared ID normalization.

**Plan validation:** All detailed requirements map to one of three tasks; no orphan task or over-sized task list found.

*Generated by Plan agent*

## Exploration

- The pre-REQ-190 `present-work` history contains the reusable source shape: `package.json`, `tsconfig.json`, `src/Root.tsx`, `src/Video.tsx`, `src/styles.ts`, and Problem, Solution, Architecture, and Value scene components. Preserve module-scope `registerRoot`, internally consistent sequence offsets/durations, proportional depth, progressive React/CSS/inline-SVG animation, real archive claims, and honest qualitative value.
- Do not carry forward the old backgrounded `npx remotion studio ... --no-open & sleep 3 && open http://localhost:3000` script, fixed-port/readiness assumptions, external-asset exceptions, fabricated before states, a client-brief dependency, or any render/media path. Use only the local package binary through the exact foreground preview script `remotion studio src/Root.tsx`.
- `actions/completed-work-presentation-reference.md` is the sole owner of target resolution, safety order, terminal-success/issues handling, archive fields, provenance, merge/current-code inspection, and collision-safe publication. The action names the preferred `<canonical-ID>-video` directory and applies the reference before any output creation; REQ-197 owns the shared ID-normalization defect.
- Direct dispatch is the explicit authority: the action receives only a target argument, so it must define never-automatic ownership rather than reparsing trigger words. Trivial and genuinely non-visual work is a successful concise skip, but evidence-backed architecture or data flow can remain eligible even without UI screenshots.
- The package contract must require valid complete JSON/TSX, compatible local React/Remotion/TypeScript dependencies, a private package, consistent frame totals, system fonts and inline assets only, no lockfile or install, and no automatic Studio launch or rendered media creation. The guide may give manual `npm install` and `npm run preview` steps, with Studio remaining in the foreground.
- RED is the absence of both standalone files plus the historical unsafe script. GREEN probes should check exact positive contract tokens and narrowly reject executable render/background/sleep/fixed-port/open forms; REQ-192 owns durable routing and regression-test integration.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/present-video.md` (new) — define the explicit source-only walkthrough workflow and valid Remotion package contract
- `skills/do-work-toolbox/docs/present-video-guide.md` (new) — document invocation, eligibility, safe manual preview, and output boundaries

**Files I will NOT touch:** `SKILL.md`, routing/help/tutorial/next-step surfaces, test files, shared completed-work reference, `ai-report`, `present-work`, schemas/archive/review-work behavior, dependencies, generated artifacts, publishing/hosting/search, or any existing deliverable; REQ-192 owns caller and durable-test migration, and REQ-197 owns shared ID normalization.

**Acceptance criteria (restated from REQ):**
- [x] An explicit `present-video`, Remotion, or video-walkthrough request can select a terminally successful completed target through the shared reference, with `completed-with-issues` qualifications preserved and unfinished/cancelled targets rejected.
- [x] Trivial or genuinely non-visual work skips successfully before output creation; eligible work produces a collision-safe `do-work/deliverables/<ID>-video/` source tree and never overwrites prior artifacts.
- [x] The complete valid Remotion/React/TypeScript source contract uses Problem → Solution → Architecture → Value, module-scope `registerRoot`, consistent proportional frame totals, verified claims, and qualitative value without fabricated metrics or before states.
- [x] Source uses only React/CSS/system fonts/inline SVG with no network, CDN, or imported image/video/audio/font/CSS assets; the action never installs, launches Studio, renders, or creates rendered media.
- [x] The only package script is `"preview": "remotion studio src/Root.tsx"`; the guide's manual preview stays in the foreground and prescribes no backgrounding, sleep, readiness guess, fixed port, browser launch, macOS `open`, or render command.
- [x] `ai-report`, `present-work`, and completion flows never call this action automatically, and routing/test/discovery migration remains deferred to REQ-192.

## Decisions

### D-01: Treat direct dispatch as explicit video authority

**Decision:** The action accepts only target arguments and treats dispatch to `present-video` as the explicit opt-in; it does not attempt to reparse the consumed trigger phrase.

**Why:** This matches the toolbox action boundary and avoids rejecting valid routed aliases whose command word is no longer present in `$ARGUMENTS`.

### D-02: Keep evidence-backed non-UI flows eligible

**Decision:** Backend or infrastructure work can pass the animation gate when verified architecture or data flow materially benefits from animation.

**Why:** “Genuinely non-visual” is an evidence judgment, not a synonym for “no screenshots.”

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/present-video.md` (new)
- `skills/do-work-toolbox/docs/present-video-guide.md` (new)

**What was done:** Added the explicit source-only walkthrough action and guide with shared completed-work evidence delegation, proportional verified scenes, successful skip behavior, immutable collision-safe output, and manual foreground-only preview guidance. No tests were touched because REQ-192 owns durable routing and contract-test integration.

## Qualification

Passed — 2 substantive files verified, 6 acceptance criteria traced, P-A-U confirmed, and scope drift passed. The mechanical checker warned that the new guide had no tracked static reference; direct inspection confirmed the new action links it, and the warning exists because both files remain untracked until this REQ's commit. Data flow is complete from explicit dispatch through the shared evidence ledger, eligibility gate, fresh source publication, verification, and manual preview handoff.

## Testing

**Tests run:** focused present-video positive/negative source assertions; `bash _dev/tests/action-shell-blocks.sh`; `bash _dev/tests/shipped-package-reference-contract.sh`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/prescribed-shell-canonicalization.sh`; `skills/do-work/tools/checks/qualify.sh do-work/working/REQ-191-extract-explicit-present-video-action.md`; `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-191-extract-explicit-present-video-action.md`; `bash _dev/tests/maintainer-verify.sh`; `git diff --check`

**Result:** ✓ All passing. Canonical maintainer verification completed with exit 0.

**Red-green validation:**
- Standalone action and guide presence: ✗ both files absent before implementation → ✓ both substantive files present after implementation.
- Unsafe preview provenance: ✗ historical embedded script backgrounded Studio, slept, assumed port 3000, and invoked macOS `open` → ✓ the standalone contract contains only the exact foreground preview script and targeted executable-form scans find none of those patterns.
- Source contract: ✗ no independently invocable source package specification existed → ✓ shared-reference delegation, four exact scenes, module-scope `registerRoot`, proportional frame math, no-external-assets rules, and source-only/no-render boundaries all pass focused assertions.

**New tests added:** None; REQ-192 owns durable caller and contract-test integration for this batch.

**Existing tests updated (cross-REQ impact):** None.

*Verified by work action*

## Review

**Overall: 83%** | 2026-08-15T17:43:17Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 92% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- `present-video.md` partially restates the shared reference's collision/no-overwrite algorithm after declaring that reference canonical, creating wording that can drift — gate: trivial → REQ-201 created as a root-cause sweep.

**Minor findings:** 1 (the shared reference still calls `present-video` a future action; report only, with current cross-reference migration already owned by REQ-192)
**Acceptance:** Partial — the standalone source-only action and guide pass static and contract verification; one canonical-publication restatement remains and generated-package execution is deferred to REQ-192.
**Suggested testing:** 4 items
**Follow-ups created:** REQ-201; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Recovering the removed Remotion section from Git history separated its useful source/package contract from the unsafe preview wrapper without restoring obsolete `present-work` modes.
- Treating direct action dispatch as explicit authority kept `$ARGUMENTS` target-only and made the never-automatic boundary auditable.

**What didn't:**
- The first action draft delegated publication canonically and then restated the collision branch in later steps and checks; REQ-201 now sweeps that duplicated-rule class.
- Qualification's tracked-reference heuristic warned on the new guide because both linked files were untracked; full-file inspection and the shipped-reference suite were needed to judge the actual link.

**Worth knowing:** Source-only video work still needs a complete, internally consistent package contract. Keep consumer-specific preferred naming and scene validity local, but leave archive, evidence, and publication algorithms in the shared completed-work presentation reference.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

[MAP CHANGED] The toolbox presentation family now has an explicit animated-walkthrough action alongside detailed `ai-report` and portfolio-only `present-work`. `present-video` owns only animation eligibility, verified four-scene storytelling, and valid source shape; shared completed-work evidence and publication stay centralized, with public routing and durable generated-package coverage following in REQ-192.
