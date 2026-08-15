---
id: REQ-192
title: Migrate completed-work presentation routing documentation and contracts
status: completed
claimed_at: 2026-08-15T17:48:42Z
route: C
completed_at: 2026-08-15T18:46:00Z
commit:
kb_status: pending
kb_entry:
created_at: 2026-08-15T09:10:53Z
user_request: UR-042
domain: testing
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-189, REQ-190, REQ-191]
maintenance: true
effort_estimate: normal
related: [REQ-189, REQ-190, REQ-191]
batch: completed-work-presentation-consolidation
write_set: [README.md, skills/do-work-toolbox/SKILL.md, skills/do-work-toolbox/actions/help.md, skills/do-work-toolbox/actions/tutorial.md, skills/do-work-toolbox/actions/completed-work-presentation-reference.md, skills/do-work-toolbox/docs/present-work-guide.md, skills/do-work/actions/help.md, skills/do-work/actions/capture.md, skills/do-work/actions/review-work.md, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/abandon.md, skills/do-work/crew-members/prompt-injection.md, skills/do-work/crew-members/anti-slop.md, skills/do-work-knowledge/crew-members/prompt-injection.md, skills/do-work-knowledge/crew-members/anti-slop.md, skills/do-work-toolbox/crew-members/prompt-injection.md, skills/do-work-toolbox/crew-members/anti-slop.md, _dev/tests/contract-regressions.sh, _dev/tests/staged-skills-contract.sh]
---

# Migrate Completed-Work Presentation Routing Documentation and Contracts

## What

Update toolbox routing, discovery, completion-flow recommendations, cross-references, tests, and release notes so the suite presents one unambiguous completed-work choice: detailed report through `ai-report`, portfolio through `present-work`, and animated walkthrough through `present-video`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add failing ownership contracts first, migrate toolbox routing/discovery plus live lifecycle guidance, generalize the three mirrored guardrail caller conditions, align only necessary guides/inventories, and verify the complete three-command split while preserving the focused REQ-197/198/199/201 defects for their owners.
- [x] **[APPLY]:** Added semantic ownership assertions first, confirmed the intended RED failures, then migrated exactly the 20 scoped discovery, guidance, guardrail, inventory, and test files without touching presentation mechanics or focused follow-up behavior.
- [x] **[UNIFY]:** Reviewed the complete scoped diff (20 files), verified no debug artifacts, ran both updated suites to GREEN, passed shell syntax and `git diff --check`, and confirmed excluded action/guide/follow-up/foreign paths have no diff.

## Why

The live router, help, tutorials, guides, next-step recommendations, caller lists, and historical successor cross-references currently point at overlapping command ownership. Without a coordinated migration, users can still enter retired detail and video workflows even after the underlying actions are separated.

## Context

REQ-189 establishes canonical detailed reporting, REQ-190 narrows `present-work`, and REQ-191 adds `present-video`. This dependent REQ makes those boundaries discoverable and adds regression coverage that rejects their regrowth.

## Detailed Requirements

### Routing and discovery

- Update the toolbox router and argument hint for the three command families.
- Route `showcase`, `visual report`, and `proof of work` to `ai-report`.
- Route `present-work`, `portfolio`, and `work portfolio` to portfolio-only `present-work`.
- Route `present-video`, `remotion`, and `video walkthrough` to `present-video`.
- Preserve explicit `ai-report` routing and do not rename it.
- Ensure bare `present-work` is usage-only and item-specific `present-work` is migration-guidance-only; do not silently delegate those invocations in the router.

### Documentation and cross-references

- Update toolbox help, tutorials, action guides, next-step recommendations, caller lists, and current cross-references for the new ownership.
- Replace completion-flow recommendations that currently suggest `present-work UR-NNN` with `ai-report UR-NNN`.
- Keep cross-project portfolio examples on `present-work all|portfolio`.
- Explain that every valid portfolio run guides the user through an optional timestamped snapshot after disclosing that the canonical summary will be refreshed. Do not introduce a snapshot flag or require a secondary command the user must remember.
- Document the safe fallback: if no snapshot answer can be collected, preserve a timestamped snapshot by default.
- Document that snapshots are never deleted automatically and may be cited by future REQs or Lessons Learned without back-editing archived records.
- Add explicit animated walkthrough examples using `present-video <ID>` only.
- Present one unambiguous choice everywhere current behavior is documented:
  - detailed report → `ai-report`;
  - cross-project portfolio → `present-work`;
  - animated walkthrough → `present-video`.
- Update prompt-injection and anti-slop caller guidance so every archive-reading or human-artifact action is covered by the condition and current examples are not stale.
- Correct current successor statements and cross-references without deleting or rewriting accurate historical records of previously generated artifacts.
- Preserve all existing generated artifacts; do not migrate or delete old briefs, `.single.html` files, reports, or video directories.

### Contract and acceptance coverage

- Lock in that `ai-report` is the only action capable of producing stakeholder-facing detailed HTML.
- Lock in that UI work retains screenshots, SVG callout annotations, responsive layout, authentic before/after evidence, and full-page light/dark render verification.
- Lock in that backend and refactor work succeeds in non-visual evidence mode without fabricated screenshots and states that UI captures were not expected.
- Prove that `present-work all` and `present-work portfolio` always refresh `do-work/deliverables/portfolio-summary.md`, produce no per-item artifact, and create no conditional output other than the guided timestamped portfolio snapshot.
- Prove all snapshot prompt branches: no creates only the canonical summary; yes creates the canonical summary plus one timestamped snapshot; no available answer uses the yes branch as the safe default.
- Prove that snapshot publication never overwrites an existing snapshot and that no presentation path automatically deletes one.
- Prove that bare `present-work` writes nothing and prints compact usage.
- Prove that item-specific `present-work` writes nothing and prints the exact `ai-report <ID>` and `present-video <ID>` replacement commands.
- Prove that `present-video` creates a valid Remotion source tree with `registerRoot` and no MP4.
- Prove that `ai-report` and `present-work` never create video artifacts and never invoke `present-video` automatically.
- Prove that all archive-reading actions load prompt-injection guidance before archived user content.
- Prove that presentation actions accept the terminal-success states `completed` and `completed-with-issues` and reject cancelled or unfinished work.
- Make contract tests reject stale Detail Mode, Interactive Explainer, detail Client Brief, sibling-link/detail-depth, automatic video, and unsafe Remotion preview workflows.
- Update shipped-package inventories, caller lists, action/guide counts, routing fixtures, and other contract-owned enumerations required by the new files and routes.

### Release record

- Add a concise changelog entry explaining the command migration and the three resulting ownership boundaries.
- Apply the repository’s required shared version bump and installed changelog mirror synchronization at the integrating commit.

## Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- Do not add `--with-video` or automatic video behavior to `ai-report`.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
- Preserve all existing generated artifacts.
- Treat the canonical portfolio summary as the sole intentional in-place refresh; preserve all timestamped snapshots until an explicit user-approved cleanup.
- Keep tests focused on the real ownership and archive-safety failures; do not add decorative snapshots of prose.

## Dependencies

Depends on REQ-189, REQ-190, and REQ-191 so the router, docs, inventory, and tests describe actions that already have their final ownership.

## Builder Guidance

Firm intent. Treat current-behavior references as a condition-based sweep, not a hand-maintained filename checklist; keep accurate history intact while removing stale live guidance. Start with failing contract cases for the current overlapping ownership and unsafe preview rules.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Search current routers, argument hints, help, tutorials, completion-flow examples, guides, caller lists, shipped inventories, and contract fixtures for the three command families and the retired `present-work` detail/video workflows, then run the baseline contract suites.
**Why RED now:** `showcase` still routes to `present-work`; completion guidance still recommends `present-work UR-NNN`; no `present-video` action is discoverable; live docs and contracts still endorse Client Brief, Interactive Explainer, detail-depth, and embedded Remotion behavior while baseline checks pass.
**GREEN when:** Every live surface presents the same three-way command ownership and guides portfolio users through snapshot retention without a memorized flag; exact routing aliases dispatch correctly; tests cover no/yes/unanswered snapshot prompts, canonical refresh, snapshot no-clobber, and explicit-only cleanup; focused contract tests fail if detailed HTML or video regrows under `present-work`, if non-visual work redirects away from `ai-report`, if unsafe Remotion preview commands or MP4 rendering return, or if terminal-success/prompt-injection rules drift; inventories and baseline suites pass; the changelog records the migration without modifying old generated artifacts.
**Validation:** User-adjusted during verification for guided snapshot behavior and its safe default; the remaining GREEN acceptance criteria were user-specified and the concrete RED case was inferred during capture.

## Full Context

See `do-work/user-requests/UR-042/input.md` for the complete verbatim request and batch constraints.

---
*Source: attached “do-work capture-request: Consolidate completed-work presentation around ai-repo…” specification.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This is a coordinated migration across routing, help, tutorials, completion guidance, guardrail callers, shipped inventories, and several contract suites after three dependent action contracts changed. It requires a condition-based surface sweep, test-first ownership assertions, and careful separation from the focused review follow-ups.

**Planning:** Required

## Plan

1. **Write RED-first presentation ownership contracts.** Add complete source-seam assertions for routing aliases, hint/help and shipped inventories, detailed visual/non-visual report evidence, portfolio-only branches and snapshot retention, explicit source-only video validity, terminal-success and guardrail load conditions, and deletion of retired detail/video/unsafe-preview language. Run the focused cases before prose changes and confirm assertion failures rather than syntax failures.
2. **Migrate toolbox routing and discovery.** Route showcase/report/proof language to `ai-report`, portfolio language to portfolio-only `present-work`, and explicit video/Remotion language to `present-video`; update the argument hint, help inventory, and action/guide enumeration without preserving broad retired detail aliases.
3. **Migrate live lifecycle and tutorial guidance.** Make item-level completion examples use `ai-report <ID>`, retain cross-project portfolio examples on `present-work all|portfolio`, add explicit `present-video <ID>` examples, and correct live caller/cross-reference prose without rewriting accurate archived history or generated artifacts.
4. **Generalize guardrail caller contracts.** Update all three shipped prompt-injection and anti-slop mirrors so the archive-ingestion and human-artifact conditions are canonical and action names remain illustrative, not exhaustive caller lists.
5. **Align only guides that remain stale.** Preserve the completed report/action implementations; touch user guides only where the current three-way ownership, guided portfolio snapshot disclosure/default/retention, evidence modes, or source-only video boundary is missing.
6. **Unify and release.** Run focused ownership assertions, contract/shell/staged-package suites, scope drift, full-file review, `git diff --check`, and canonical maintainer verification. The integrating work phase owns the shared version bump and mirrored changelog entry after product contracts pass.

**Requirement coverage:** Tasks 1–2 cover routing, discovery, inventories, and durable ownership rejection. Tasks 2–5 cover help/tutorial/caller/cross-reference migration, snapshot guidance, evidence modes, and condition-based guardrails. Task 1 covers all acceptance branches and terminal states. Task 6 covers release synchronization and final proof.

**Candidate implementation set:** toolbox `SKILL.md`, help/tutorial, the three presentation guides where proven stale, core help/capture/work/abandon live guidance, all three prompt-injection and anti-slop mirrors, and presentation-owned portions of `contract-regressions.sh`, `staged-skills-contract.sh`, and `prescribed-shell-canonicalization.sh`. Exploration will reduce this to files with a concrete required edit before Scope is declared.

**Risks:** Broad term replacement can corrupt accurate history; vague prose assertions can pass without locking behavior; inventories can become competing sources of truth. Search live shipped instructions and tests conditionally, assert complete predicates at canonical seams, and preserve all generated/archive artifacts.

**Non-goals:** No presentation implementation mechanics, shared resolver/publication algorithms, schemas, archive formats, `review-work`, generated artifacts, publishing/hosting/search, MP4 rendering, automatic video, snapshot flags, queue-kanban code, or foreign REQ-200 changes. REQ-197, REQ-198, REQ-199, and REQ-201 retain their focused defects.

**Plan validation:** Every Detailed Requirements group maps to a task, RED precedes migration, historical surfaces remain excluded, and no focused follow-up is duplicated.

*Generated by Plan agent*

## Exploration

- The three presentation actions already carry the intended behavior. This REQ should change discovery, live successor guidance, caller descriptions, and durable assertions; it must not repair the implementation defects owned by REQ-197, REQ-198, REQ-199, or REQ-201.
- Live stale surfaces are root `README.md`; toolbox `SKILL.md`, help, tutorial, shared-reference description, and portfolio guide; core help/capture/review/work/work-reference/abandon guidance; and all six prompt-injection/anti-slop mirrors. The portfolio guide alone needs new user-facing content: snapshots may be cited by future REQs or Lessons without back-editing archives.
- Router ownership must be exact: `ai-report|showcase|visual report|proof of work` → `ai-report`; `present-work|portfolio|work portfolio` → portfolio; `present-video|remotion|video walkthrough` → video. Retire broad `present` and `client brief` aliases rather than silently mapping detail intent to a portfolio.
- `review-work.md` still claims `present-work` parses persisted scores even though portfolio-only `present-work` no longer does. Core terminal-success and commit-traceability examples should point to the shared completed-work presentation reader instead of a closed action enumeration.
- Contract RED belongs in `_dev/tests/contract-regressions.sh` and `_dev/tests/staged-skills-contract.sh`: router/hint/help/inventory, README/tutorial completion examples, condition-based guardrail callers, detailed visual/non-visual report evidence, all portfolio branches and retention, explicit source-only video validity, shared terminal success/safety, and precise retired/unsafe workflow negatives.
- Correctly excluded: the three action implementations, `ai-report-reference.md`, `ai-report-guide.md`, `present-video-guide.md`, `next-steps.md`, prescribed-shell tests, and the retired-core trigger fixture. `present-video` was never a moved core command, and its manual preview does not restate a promoted shell primitive.
- Foreign REQ-200 currently owns queue-kanban product edits and has temporarily touched the five shared release files. Preserve and exclude all of it; defer this REQ's version/changelog integration until that release finishes.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `README.md` (modify) — migrate the full-cycle presentation step and expose the three distinct presentation commands
- `skills/do-work-toolbox/SKILL.md` (modify) — update the argument hint and exact routing aliases
- `skills/do-work-toolbox/actions/help.md` (modify) — describe detailed report, portfolio, and source walkthrough separately
- `skills/do-work-toolbox/actions/tutorial.md` (modify) — migrate full-cycle and post-work recipes to the explicit three-way ownership
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (modify) — correct the stale future-action description only
- `skills/do-work-toolbox/docs/present-work-guide.md` (modify) — document durable snapshot citation without archive back-editing
- `skills/do-work/actions/help.md` (modify) — add video discovery and make the full-cycle prompt use `ai-report`
- `skills/do-work/actions/capture.md` (modify) — generalize preserved-input reader guidance
- `skills/do-work/actions/review-work.md` (modify) — remove the retired portfolio score-parser successor claim
- `skills/do-work/actions/work.md` (modify) — generalize completed-with-issues presentation visibility
- `skills/do-work/actions/work-reference.md` (modify) — refresh illustrative terminal-success and traceability consumers without implementing ID normalization
- `skills/do-work/actions/abandon.md` (modify) — generalize archived body-reader guidance
- `skills/do-work/crew-members/prompt-injection.md` (modify) — refresh the condition-based archive-ingestion caller example
- `skills/do-work/crew-members/anti-slop.md` (modify) — refresh illustrative human-artifact callers
- `skills/do-work-knowledge/crew-members/prompt-injection.md` (modify) — mirror the condition-based archive-ingestion caller example
- `skills/do-work-knowledge/crew-members/anti-slop.md` (modify) — mirror illustrative human-artifact callers
- `skills/do-work-toolbox/crew-members/prompt-injection.md` (modify) — mirror the condition-based archive-ingestion caller example
- `skills/do-work-toolbox/crew-members/anti-slop.md` (modify) — mirror illustrative human-artifact callers
- `_dev/tests/contract-regressions.sh` (modify) — add focused three-command ownership and safety/evidence regressions
- `_dev/tests/staged-skills-contract.sh` (modify) — update shipped action/guide inventories and the approved full-cycle prompt

**Files I will NOT touch:** presentation action implementation mechanics; `ai-report-reference.md`; `ai-report-guide.md`; `present-video-guide.md`; prescribed-shell tests; retired-core route fixtures; schemas/archive formats; `review-work` behavior beyond its stale reader description; generated artifacts; queue-kanban files; or foreign REQ-200 work. Shared version/changelog files remain integrator-only lifecycle edits after product verification.

**Acceptance criteria (restated from REQ):**
- [x] Every live router, hint, help, tutorial, README, and completion-flow surface exposes the same exact detailed-report, portfolio, and animated-walkthrough ownership.
- [x] Bare and item-specific `present-work` remain non-writing guidance; `all|portfolio` refresh the canonical summary and guide No, Yes, and unanswered snapshot outcomes without flags, overwrite, deletion, or archive back-editing.
- [x] Prompt-injection and anti-slop caller guidance is condition-based across all three shipped mirrors and covers current presentation readers/artifacts without a closed enumeration.
- [x] Durable contracts pin detailed visual and non-visual `ai-report` evidence, portfolio-only output and snapshot branches, explicit valid source-only `present-video`, terminal-success/archive-safety rules, and absence of retired or unsafe workflows.
- [x] Shipped action/guide inventories and approved full-cycle fixtures include the new action and use `ai-report` for item-level completion.
- [x] Accurate historical records and existing generated artifacts remain byte-for-byte untouched, while REQ-197/198/199/201 retain their focused defects.

## Decisions

### D-01: Retire ambiguous presentation aliases instead of guessing ownership

**Decision:** Remove the broad `present` and `client brief` routes. Keep only the explicit alias families assigned by the request.

**Why:** Mapping ambiguous item-detail language to portfolio-only `present-work` would preserve the overlap this batch removes; unknown descriptive requests can fall back to help/core capture.

### D-02: Make guardrail applicability conditional and examples illustrative

**Decision:** All three prompt-injection and anti-slop mirrors name the ingestion/artifact condition first, then cite the current presentation family only as examples.

**Why:** A closed caller list was the source of the stale `present-work` ownership and would fail again when another archive reader or human artifact is added.

## Implementation Summary

**Files changed:**
- `README.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work-toolbox/actions/tutorial.md` (modified)
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (modified)
- `skills/do-work-toolbox/docs/present-work-guide.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `skills/do-work/actions/capture.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/abandon.md` (modified)
- `skills/do-work/crew-members/prompt-injection.md` (modified)
- `skills/do-work/crew-members/anti-slop.md` (modified)
- `skills/do-work-knowledge/crew-members/prompt-injection.md` (modified)
- `skills/do-work-knowledge/crew-members/anti-slop.md` (modified)
- `skills/do-work-toolbox/crew-members/prompt-injection.md` (modified)
- `skills/do-work-toolbox/crew-members/anti-slop.md` (modified)

**What was done:** Migrated live routing, discovery, full-cycle guidance, caller descriptions, and portfolio retention documentation to one explicit three-command ownership model. Added RED-first durable contracts for routing, inventories, report evidence modes, portfolio branches, source-only video, archive safety/status, and retired-workflow rejection while preserving the four focused follow-up defects for their owning REQs.

## Qualification

**Attempt 1:** Mechanical qualification and scope drift passed, and all 20 files are substantive and wired. Requirements tracing found three incomplete negative/order ratchets in `_dev/tests/contract-regressions.sh`: it did not prove prompt-injection-before-archive-content order for every presentation reader, did not reject retired detail/brief/explainer/depth regrowth directly inside `present-work.md`, and did not scan both the video action and guide for executable background/sleep/fixed-port/browser-open/render forms. Returned to the builder for one focused in-scope test correction.

**Attempt 2:** Passed — 20 files verified, 6 acceptance groups traced, P-A-U confirmed, and scope drift clean. The focused correction now pins ordered safety loads, direct retired-workflow rejection in `present-work`, unsafe executable forms across the video action and guide without false-positive prohibition prose, and explicit no-video ownership in `ai-report`.

## Testing

**Tests run:** preflight baseline with `bash _dev/tests/contract-regressions.sh`; test-only RED runs of `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/staged-skills-contract.sh`; GREEN reruns of both; `bash -n _dev/tests/contract-regressions.sh`; `bash -n _dev/tests/staged-skills-contract.sh`; `skills/do-work/tools/checks/qualify.sh do-work/working/REQ-192-migrate-presentation-routing-docs-and-contracts.md`; `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-192-migrate-presentation-routing-docs-and-contracts.md`; `bash _dev/tests/maintainer-verify.sh`; `git diff --check`

**Result:** ✓ All passing. Canonical maintainer verification completed with exit 0.

**Red-green validation:**
- Routing/discovery ownership: ✗ new tests failed because `showcase` still routed to `present-work`, `present-video` was absent from hint/router/help/inventory, and public action count was 17 rather than 18 → ✓ exact three-family aliases, destinations, hint, help, and shipped inventory pass.
- Completion guidance: ✗ README, core help, tutorial recipes, and Topic 4 still recommended item-specific `present-work` → ✓ all item-level full-cycle examples use `ai-report`, with portfolio and video kept explicit.
- Safety/caller contracts: ✗ guardrail mirrors retained obsolete artifact ownership, the shared reference called video future work, and core caller descriptions named stale individual readers → ✓ condition-based mirrored callers, current shared-reference wording, ordered safety loads, and illustrative presentation-family readers pass.
- Artifact ownership: ✗ no durable block rejected detailed/brief/explainer regrowth under portfolio work or unsafe/automatic video workflows → ✓ detailed visual/non-visual report evidence, all portfolio branches/retention, valid source-only video, retired-workflow negatives, and executable preview negatives pass.

**New tests added:**
- `_dev/tests/contract-regressions.sh` — one focused presentation ownership block with semantic router, evidence, branch, safety/order, inventory, and regrowth assertions.

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/staged-skills-contract.sh` — adds the REQ-191 action/reference/guide to shipped inventory and migrates the approved full-cycle fixture from retired item-specific portfolio guidance to `ai-report`.

*Verified by work action*

## Review

**Overall: 83% — Approve with follow-up (Partial)**

**Route:** C
**Reviewed:** 2026-08-15T18:45:11Z

### Summary

The three-command migration is coherent and the current source contracts are safe. Routing, hints, help, tutorials, lifecycle guidance, guardrail conditions, shared-reference wording, guides, and shipped inventories agree on detailed HTML through `ai-report`, portfolio output through `present-work`, and source-only animation through `present-video`.

### Findings

**Important — rule-change:** `_dev/tests/contract-regressions.sh` does not reject every unsafe preview form that the new contract claims to pin. Executable `remotion studio src/Root.tsx --port 3000` / `--port=3000` and nonliteral platform-opener commands such as `open "$REMOTION_PREVIEW_URL"` can escape the matcher even though literal fixed-URL, background/sleep, and render forms are caught. A focused follow-up must complete the mutation detector while preserving safe foreground-preview examples and negative prose.

**Minor:** None.

### Acceptance

**Result: Partial.** Scope drift, both changed shell files' syntax, focused contract regressions, staged-skill contracts, `git diff --check`, P-A-U, and canonical maintainer verification pass. Focused mutation probes demonstrated the remaining detector gap.

### Follow-ups

- `REQ-202` — Complete unsafe Remotion preview mutation detection (pending; review-generated sweep).

*Reviewed by independent Review agent*

## Lessons Learned

**What worked:**
- A condition-based surface sweep reduced a broad migration to 20 live files while leaving accurate history, generated artifacts, presentation mechanics, and focused review follow-ups untouched.
- Adding the routing and inventory tests before product edits produced a useful RED signal for every stale command family and then verified the exact three-owner model.

**What didn't:**
- The first GREEN contract block under-specified load order, retired portfolio workflows, and unsafe guide commands. Qualification caught those gaps and sent the test seam back for focused correction.
- The widened unsafe-preview matcher still encoded literal examples too narrowly; review mutation probes exposed fixed-port flag and nonliteral platform-opener escapes, now owned by REQ-202.

**Worth knowing:** Presentation routing aliases must remain exact and mutually exclusive, while guardrail applicability should be condition-based with action names treated only as examples. Test executable command segments separately from negative explanatory prose so safety assertions can be both broad and precise.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

[MAP CHANGED] Completed-work presentation discovery now has three explicit owners everywhere: detailed visual or non-visual HTML through `ai-report`, cross-project portfolio through `present-work`, and explicit source-only animation through `present-video`. Lifecycle guidance, guardrail conditions, shipped inventories, and contract tests share that map.
