---
id: REQ-687
title: 'ai-report --kind proposal and root-cause writes a decision-first brief for unfinished work'
status: completed
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-655, REQ-656, REQ-657]
batch: ai-report-modes
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/ai-report-reference.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md"]
claimed_at: 2026-10-10T12:52:18Z
status_changed_at: 2026-10-10T13:14:50Z
required_lessons: ["_dev/primes/lessons-releases.md"]
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-10-10T13:19:15Z
  basis:
    - Route B
    - 6-file write set
    - 9 acceptance criteria
builder_handback_at: 2026-10-10T13:27:53Z
integration_at: 2026-10-10T17:26:47Z
review_at: 2026-10-10T17:38:45Z
kb_status: pending
commit: 8c0cf7ed44e66f763b8efbbb1bc182dc49bd4b9e
heavy_verified_at: 2026-10-10T17:39:16Z
heavy_verified_revision: 8c0cf7ed44e66f763b8efbbb1bc182dc49bd4b9e
completed_at: 2026-10-10T17:39:37Z
release_at: 2026-10-10T17:39:37Z
---
# ai-report --kind Proposal and Root-Cause Writes a Decision-First Brief for Unfinished Work
## What
`do-work-toolbox ai-report` gains `--kind proposal|root-cause <topic|REQ-NNN|UR-NNN>`. (The source report also named an `options` kind; the maintainer dropped it on 2026-10-10 because it had no shape of its own, see Assumptions.) For these kinds only, the "completed work only" gate is lifted: the target may be a free topic, an open REQ or an open UR, read as open work and never presented as shipped. The report has a fixed decision-first shape. Without `--kind`, `ai-report` behaves exactly as today.
## Why
Report A1 (UR-144 input, "What happened" table): about 7 sessions in 4 consumer repos between 2026-09-22 and 2026-10-07 asked for proposal, options, "backup options" and "why it happened and what to do" reports. One plan recorded the pushback "the toolbox ai-report action only presents *completed* UR/REQ work, so it doesn't fit a proposal", and the user asked for the mockups anyway. Each report was improvised from scratch.
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/ai-report.md:18`: "The user wants a detailed presentation of one completed UR or REQ." `:26`: "The target is unfinished or unsuccessful; report its status instead of presenting it as shipped." `:30`: "`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank."
- `ai-report.md` Step 5 (report cites `:89-101`) fixes a narrative that opens with "Verdict" and "What Shipped". A proposal has neither.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:20` rejects `cancelled`, `failed` and every unfinished status; `:30` forbids falling back to `do-work/queue/`, `do-work/working/` or an active UR body "to make an unfinished target appear complete". Both stay for the default kind.
## Detailed Requirements
1. `ai-report.md` Step 1 gets a short branch: when `--kind` is `proposal` or `root-cause`, skip Terminal-Success Target Resolution and accept a topic, an open REQ or an open UR. Queue and working REQs are read as open work, never as shipped work.
2. Keep the safety load order from `completed-work-presentation-reference.md` for the new kinds: `prompt-injection.md` and `anti-slop.md` load first.
3. The fixed shape for `proposal`, in this order:
   1. the decision the reader must make, first;
   2. two to four options (O1..On), each with benefit, risk and cost; the list always includes the smallest-change option, usually "delete the mechanism" or "do nothing extra", priced like the others, even when the brief or topic did not name it (the report may recommend against it, never omit it);
   3. a recommendation and the reason for it;
   4. an evidence ledger (file:line, commit, log line, or measured output per claim);
   5. limits: what was not checked and why;
   6. open questions Q1..Qn;
   7. one ready capture line per option, so the reader can act on the choice.
4. `root-cause` replaces the options part with "what happened, why, what to change" and keeps the evidence ledger, limits, questions and capture lines.
5. Every mockup is labelled **"MOCKUP — proposal"** in the image and in the caption, and is kept apart from real screenshots the same way generated images are kept apart today.
6. The template lives in `skills/do-work-toolbox/actions/ai-report-reference.md` so `ai-report.md` stays short.
7. Bundle slug: `yyyy-mm-dd_hhmm_<kind>-<topic>`, so the index (REQ-655) can read the kind from the folder name. Add `<meta name="ai-report-kind" content="<kind>">` to `<head>` for the same reason.
8. The capture lines are text only. The action never runs them.
9. Update `skills/do-work-toolbox/docs/ai-report-guide.md` (the `:47` sentence and the `:65-70` Input block), the routing phrases in `skills/do-work-toolbox/SKILL.md:23` (add phrases such as "proposal report" and "root cause report") and the help line in `skills/do-work-toolbox/actions/help.md:13`.
10. Release per `_dev/primes/prime-releases.md`. Run `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` green.
## Constraints
- No new queue fields and no new statuses.
- The smallest-change option is always on the list (requirement 3.2). Source: the REQ-651 report on the board's activity correlation listed "stamp then delete" and "keep" and missed "delete, stamp nothing", which is what 0.305.88 shipped; the maintainer ruled on 2026-10-10 that decision reports always price it (UR-150).
- The default kind keeps the completed-work gate exactly as today, including `completed-work-presentation-reference.md:20` and `:30`.
- Out of scope: a free-form brainstorm (`deep-explore` owns it), running the capture lines, committing or publishing the report.
- Existing bundles are immutable (`completed-work-presentation-reference.md:71`, `:73`).
- This REQ is its own release; do not fold REQ-655, REQ-656 or REQ-657 into it.
## Assumptions (recorded at capture, no questions asked)
- The `options` kind named in the source report is dropped: its own capture assumption said it would share the `proposal` template and differ only in label, so a report that compares options is `--kind proposal`. If a phrase such as "options report" is routed, it routes to `proposal`.
- The report writes the capture line as `/do-work capture-request ...`, a consumer slash-command spelling. The builder writes the invocation form the toolbox already uses in its other printed next-step lines (`do-work capture-request: <task>` per `skills/do-work/SKILL.md:46`).
- A free topic has no REQ or UR id; the slug topic part is a short kebab form of the topic text.
- "Open REQ" means any non-terminal status (`pending`, `pending-answers`, `blocked`, claimed). A `failed` or `cancelled` REQ is accepted for `root-cause`, where explaining a failure is the point, and is never described as shipped.
- The "kept apart like generated images" rule reuses the existing generated-image separation in `ai-report-reference.md`; no new image pipeline.
## Dependencies
None upstream. REQ-655 (index and find) reads the kind meta and slug this REQ defines, but works without it (it falls back to unknown kind), so there is no ordering edge.
## Builder Guidance
Certainty is high on the shape (the report lists it item by item). Latitude: exact wording, where the Step 1 branch sits.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox ai-report --kind proposal "backup options"` in a repo with no completed work.
**Why RED now:** `ai-report.md:30` accepts only `UR-NNN`, `REQ-NNN`, `most recent` or blank, and `:26` stops on unfinished targets, so there is no path to a proposal report.
**GREEN when:** The action produces a bundle named `yyyy-mm-dd_hhmm_proposal-<topic>` whose first section is the decision, with 2 to 4 options each showing benefit, risk and cost, a recommendation, an evidence ledger, a limits section, Q1..Qn, and one capture line per option; its `<head>` carries `ai-report-kind`; every mockup carries "MOCKUP — proposal". And `ai-report REQ-NNN` on an unfinished REQ, without `--kind`, still stops exactly as today.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing action routing and status contracts, which this REQ does; family `canonical-authoring-vs-tolerant-reading` fits a kind that reads open work without presenting it as shipped.
## Full Context
See `do-work/user-requests/UR-153/input.md` for the 2026-10-10 decision record. The cancelled original, `do-work/archive/UR-144/REQ-654-ai-report-kind-proposal-root-cause-options.md`, carries the complete body and the UR-150 addendum; the source report is `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A1, with verbatim input in `do-work/archive/UR-144/input.md`.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Prose-only change in five toolbox files, per brief and D-01..D-10. All rules of the two kinds go into one section appended at the end of `ai-report-reference.md` (clear of REQ-656's Report Design Rules edits). `ai-report.md` gets pointers only and stays clear of `:109` (REQ-655) and Step 7 (REQ-657). `completed-work-presentation-reference.md` is not edited (D-03); the deviation is declared in the new section the way `stakeholder-report.md:15` does. New guide/help lines are placed directly below the existing ones; the guide section goes after `## Evidence Safety`. No lesson in `lessons-releases.md` contradicts the plan (no history links or mirrors touched).
- [x] **[APPLY]:** Edits made as planned in the five write-set paths only. Example slug uses `REQ-NNN` rather than a real-looking number so no consumer REQ number lands in shipped prose. The existing `../../do-work/docs/prescribed-shell-primitives.md` pointers at Step 1 and Step 3 are kept.
- [x] **[UNIFY]:**
  - `git diff bd56c4b0 --stat`:
    ```
     skills/do-work-toolbox/SKILL.md                    |  2 +-
     .../do-work-toolbox/actions/ai-report-reference.md | 26 ++++++++++++++++++++++
     skills/do-work-toolbox/actions/ai-report.md        | 15 +++++++++----
     skills/do-work-toolbox/actions/help.md             |  1 +
     skills/do-work-toolbox/docs/ai-report-guide.md     |  8 ++++++-
     5 files changed, 46 insertions(+), 6 deletions(-)
    ```
  - `git diff --check`: clean (exit 0).
  - `bash .../REQ-687-probe.sh`: exit 0, `REQ-687 probe passed.`, wall 2 s.
  - `bash _dev/tests/contract-regressions.sh`: exit 0, `Contract regression checks passed.`, wall 36 s (single run, no budget failure).
  - `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, wall 2 s (also run inside the probe).
  - Files checked: all five diffs read in full. Checked: only the write set changed; `completed-work-presentation-reference.md:20` and `:30` byte-identical (probe lines pass); `ai-report.md:26` still carries "unfinished or unsuccessful"; no citation of `CLAUDE.md` or `_dev/`; no debug artifacts; no `do-work/` path staged.
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A1: "`ai-report --kind proposal|root-cause|options <topic|REQ-NNN|UR-NNN>` for work that is not completed. Lift the "completed work only" gate for these kinds only."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (a fixed report shape, a slug, a meta tag, routing and help lines), but it is prose spread across six toolbox files whose exact insertion points and neighbouring contracts (completed-work gate, generated-image separation, collision-safe publication) had to be located first. No new code, no new subsystem, one package.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Read at main HEAD bd56c4b0 (2026-10-10) by the pre-dispatch agent.

**Required-lessons consult.** `do-work/lessons-index.md` matches two rows. `_dev/primes/lessons-releases.md` (666 tokens, `slugged: full`) matches "the shipped-package reference contract" and the release this REQ carries; added whole to `required_lessons`. `_dev/primes/lessons-action-files.md` (7756 tokens, `slugged: partial`) stays dropped for budget (see `## Required Lessons — Dropped for Budget`). Its index row (action routing, status contracts, downstream readers) still matches; the part that applies is the `[family: alternate-writer-contract-drift]` trap, which `_dev/primes/prime-action-files.md` § Traps carries and `prime_files` already loads. Its one `canonical-authoring-vs-tolerant-reading` bullet is about schema parsers and does not apply to prose.

**Where each change lands (line numbers at bd56c4b0; two drifted from the capture):**
- `skills/do-work-toolbox/actions/ai-report.md` (140 lines): `:3` description blockquote ("for one completed UR or REQ"); `:18` When to Use; `:26` Do-NOT-use line "The target is unfinished or unsuccessful" (must gain "without `--kind`"); `:30` Input line; `:34-38` Step 1 (pointer to the completed-work reference); `:51` slug rule; `:67` screenshots/generated separation; `:89-101` Step 5 narrative (Verdict, What Shipped, ...); `:109` Step 6 "does not publish, host, or search" (REQ-655 edits this line); `:113` Step 7 render check (REQ-657 rewrites it); `:121` Step 8 summary; `:125` Output Format; `:133-140` Verification Checklist. `:38` and `:57` carry the pointer `../../do-work/docs/prescribed-shell-primitives.md` that `_dev/tests/prescribed-shell-canonicalization.sh:107` requires; keep at least one.
- `skills/do-work-toolbox/actions/ai-report-reference.md` (119 lines): `:3` companion blockquote; `:49-59` generated-image rules (`:51` own `generated/` folder, `:53` visible "AI-generated" badge); `:70-90` Report Design Rules (`:90` disclosure); `:117-119` Output Format Template. `_dev/tests/staged-skills-contract.sh:215` requires this file to keep the text `<skill-root>/../do-work/tools/do-work-cli.sh`.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md`: `:20` and `:30` are the default gate and stay byte-identical. Precedent for a consumer that is not completed work: `skills/do-work-toolbox/actions/stakeholder-report.md:15` says "Terminal-Success Target Resolution does not apply — only the Safety Load Order and Collision-Safe Publication sections are inherited", declared in the consumer, with no edit to the shared reference. The kinds follow that precedent, so this file needs no edit (D-03).
- `skills/do-work-toolbox/docs/ai-report-guide.md` (78 lines): `:3` intro, `:47` "Cancelled, failed, and unfinished work is rejected", `:63-72` Input block.
- `skills/do-work-toolbox/SKILL.md:24` ai-report routing row (capture said `:23`).
- `skills/do-work-toolbox/actions/help.md:14` ai-report help line (capture said `:13`).
- Capture-line form: `do-work capture-request: <task>`, as at `skills/do-work-toolbox/actions/release-check.md:127` and `journey-qa.md:107`.

**Tests.** No test pins the wording of these six files beyond the two text anchors above and the file-existence list at `_dev/tests/staged-skills-contract.sh:173-176`. `bash _dev/tests/shipped-package-reference-contract.sh` checks every relative link and path a shipped file names (green at base, 1.5 s). `bash _dev/tests/contract-regressions.sh` is green at base (32 s). No Go code or router-budget check reads the toolbox routing table.

**Restatements outside the write set (alternate-writer-contract-drift sweep).** `architecture-report.md:27` and `:129`, `present-work.md:5` and `:23`, `present-video.md:23`, `tutorial.md:245-252` and `:318`, `skills/do-work/actions/help.md:37` and `:55`, and the prompt-injection/anti-slop JIT comments all describe the default kind and stay true. Only `architecture-report.md:129` ("`ai-report` takes a UR or REQ and presents completed work") becomes incomplete; it is a rationalization row whose argument still holds, and REQ-657 edits that file, so it is left alone and reported as a discovered task (D-08).

**Pre-dispatch decisions** (no `## Open Questions` section existed; these are the builder's starting decisions, ID space shared with Step 6):
- **D-01** One label for both kinds: every mockup in a `proposal` or `root-cause` report carries the exact text "MOCKUP — proposal" inside the image and in its caption. Value: one string to check. Risk: low; a later kind-specific label is a one-line change.
- **D-02** Placement: `ai-report.md` gets a short `--kind` branch at the top of Step 1, a Do-NOT-use qualifier, the Input line, a Step 5 pointer and a Step 8 summary note. Everything else (target reading, the two templates, slug, meta, mockup rule, capture lines) is one new section at the END of `ai-report-reference.md`, so it does not touch lines REQ-655, REQ-656 or REQ-657 edit. Value: `ai-report.md` stays short (requirement 6) and merges stay clean. Risk: none.
- **D-03** `completed-work-presentation-reference.md` is not edited; the new section declares the deviation the way `stakeholder-report.md:15` does (Safety Load Order and Collision-Safe Publication inherited, Terminal-Success Target Resolution not applied). Value: `:20` and `:30` cannot drift; smallest change. Risk: none; the capture's write_set listed the file only as a candidate.
- **D-04** Topic part of the slug: for an id target, the id plus a short kebab summary (`proposal-REQ-651-activity-correlation`); for a free topic, a short kebab form of the topic text (`proposal-backup-options`). Value: matches the default slug style at `ai-report.md:51`. Risk: low.
- **D-05** `--kind` takes exactly `proposal` or `root-cause`; any other value stops with one line naming the two. A plain-language ask for an "options report" or "backup options" is a `proposal`. `--kind` with no target stops with a one-line usage. Value: no hidden alias. Risk: low.
- **D-06** Mockups are synthetic, so they live in `generated/` beside generated images and inherit their separation (`ai-report-reference.md:51`, `:53`), never in `screenshots/`. Hand-authored mockups carry the label as SVG or HTML text inside the mockup; raster mockups get the label in the generation prompt and as a badge overlaid inside the image frame. No new image pipeline.
- **D-07** `root-cause` shape: what happened, why, what to change, evidence ledger, limits, Q1..Qn, one capture line per change. When "what to change" offers more than one change, the smallest-change option is on the list and priced, as for `proposal` (the UR-150 ruling covers decision reports in general). Value: consistent rule. Risk: low.
- **D-08** `architecture-report.md:129` is not edited (outside the write set; REQ-657 owns edits there); reported as a discovered task.
- **D-09** Routing phrases added to the ai-report row: `proposal report`, `root cause report`, `options report`. Help gets its own new line under the existing ai-report line instead of a rewrite of `:14`.
- **D-10** The default kind writes no `ai-report-kind` meta and is otherwise unchanged ("Without `--kind`, `ai-report` behaves exactly as today"). REQ-655 reads an absent meta as unknown kind.
<!-- D-XX counter: last used D-10. Next decision: D-11. -->

*Generated by the pre-dispatch agent (Explore role)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — the kind branch at Step 1, the Do-NOT-use qualifier, the Input line, pointers from Step 5 and Step 8, one checklist line
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — one new end-of-file section with target reading, both templates, slug, meta tag, mockup label and capture lines
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — the sentence at line 47, two Input lines, one short section on the two kinds
- `skills/do-work-toolbox/SKILL.md` (modify) — three routing phrases on the ai-report row
- `skills/do-work-toolbox/actions/help.md` (modify) — one new help line under the ai-report line

**Files I will NOT touch:** completed-work-presentation-reference.md (D-03: the default gate at its lines 20 and 30 stays byte-identical and the kinds declare their deviation in ai-report-reference.md); architecture-report.md, stakeholder-report.md, present-work.md, present-video.md, tutorial.md (default-kind restatements stay true; D-08); any Go code, test script, CHANGELOG, VERSION or mirror (the release is the integrator's).

**Acceptance criteria (restated from REQ):**
- [ ] `ai-report --kind proposal|root-cause <topic|REQ-NNN|UR-NNN>` skips Terminal-Success Target Resolution and accepts a free topic, an open REQ or an open UR (and a `failed` or `cancelled` REQ for `root-cause`), read as open work and never described as shipped.
- [ ] The kinds keep the safety load order: `prompt-injection.md`, then `anti-slop.md`, before any REQ, UR or repository prose is read.
- [ ] `proposal` shape in order: decision first; two to four options O1..On, each with benefit, risk and cost, always including the smallest-change option priced like the others; recommendation with reason; evidence ledger (file:line, commit, log line or measured output per claim); limits; Q1..Qn; one capture line per option.
- [ ] `root-cause` replaces the options part with what happened, why, what to change, and keeps the evidence ledger, limits, questions and capture lines.
- [ ] Every mockup carries "MOCKUP — proposal" in the image and in the caption, and stays apart from real screenshots the way generated images do.
- [ ] The template lives in `ai-report-reference.md`; `ai-report.md` stays short.
- [ ] Bundle slug is `yyyy-mm-dd_hhmm_<kind>-<topic>` and `<head>` carries `<meta name="ai-report-kind" content="<kind>">`.
- [ ] Capture lines are text only (`do-work capture-request: <task>`); the action never runs them.
- [ ] The guide (line 47 sentence, Input block), the SKILL.md routing phrases and the help line are updated; without `--kind`, `ai-report REQ-NNN` on an unfinished REQ still stops exactly as today.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)

**What was done:** Added a `## Proposal and Root-Cause Kinds` section at the end of `ai-report-reference.md` that holds every rule of the two kinds: target reading (Terminal-Success Target Resolution lifted, Safety Load Order and Collision-Safe Publication inherited, open REQ/UR or free topic accepted, failed or cancelled REQ for `root-cause` only), the seven-part `proposal` template with the always-priced smallest-change option, the `root-cause` template, the `yyyy-mm-dd_hhmm_<kind>-<topic>` slug and `ai-report-kind` meta tag, the "MOCKUP — proposal" label with mockups kept in `generated/`, and text-only `do-work capture-request:` lines. `ai-report.md` gained pointers only: a `--kind` Input paragraph, a Step 1 branch, Step 2 slug and Step 5 template pointers, a Step 8 summary note, a Do-NOT-use qualifier and one checklist line. The guide, the toolbox routing row (`proposal report`, `root cause report`, `options report`) and the help list gained the new form. `completed-work-presentation-reference.md` is unchanged, so the default kind keeps its completed-work gate. Integrator merge (conflicts with REQ-655, REQ-656, REQ-657 resolved as a union): the `--kind` Input paragraph was moved out of the Revise form, where it had auto-merged, to the top-level Input text, and the narrative checklist line says a `--kind` report uses its own template. Review fixes (integrator, on the builder branch, `84d23733`): a `--kind` report inherits only the Safety Load Order and Collision-Safe Publication, so the Archive Evidence Sweep no longer stops on an open REQ (F1); the guide's terminal-success stop sentence is scoped to the default kind (F2); When to Use names the brief (F3); the reference says how Steps 3 and 4 apply to a brief (F4); the checklist's shared-reference line exempts `--kind` (F5).

## Decisions

*(from the builder hand-back)*

- **D-11** (DECIDE & STATE) A completed REQ or UR passed with `--kind` is not in the accepted-target list, and no extra rule says what happens. Adding one was not asked for (anti-bloat). If it comes up, a one-line rule can be added later.
- **D-12** (DECIDE & STATE) The `--kind` Input sentence went into a new paragraph under `ai-report.md:30` instead of onto that line, because REQ-655 edits `:30`. The merge stays clean.
- **D-13** (DECIDE & STATE) The help line cannot fit the 31-character command column, so it uses two spaces before its description. Same for the two guide Input lines.
- **D-14** (DECIDE & STATE) The slug example uses `REQ-NNN` (`proposal-REQ-NNN-retry-policy`) so no consumer REQ number appears in shipped prose.

## Discovered Tasks

*(from the builder hand-back)*

- impact-minor: `skills/do-work-toolbox/actions/architecture-report.md:129` still says "`ai-report` takes a UR or REQ and presents completed work", which is now incomplete (D-08) → report only
- impact-minor: `skills/do-work-toolbox/docs/ai-report-guide.md:3` intro and `skills/do-work-toolbox/actions/ai-report-reference.md:3` companion blockquote still describe completed work only. Both remain true for the default kind, and the brief did not name them → report only

## Qualification

**Gate records (`advance --diff-range 3bd6a668..a0ebfea9`):** `qualify` satisfied, `scope-drift` satisfied, on the first call. No debug artifacts.

**Requirement trace against the diff (`git diff 3bd6a668..a0ebfea9 --stat`: 5 files, +47/-7) and the merged files** (AR = `skills/do-work-toolbox/actions/ai-report.md`, REF = `skills/do-work-toolbox/actions/ai-report-reference.md`):
1. Step 1 branch: AR:63 sends `--kind` to REF **Proposal and Root-Cause Kinds**; REF:126 lifts Terminal-Success Target Resolution and accepts a free topic, an open REQ (queue, working or archive) or an open UR, plus a failed or cancelled REQ for `root-cause` only; everything is read as open work. Input: AR:35. Do-NOT-use qualifier: AR:27.
2. Safety load order kept: REF:126 inherits Safety Load Order and Collision-Safe Publication; AR:63 "safety load order still comes first".
3. `proposal` shape in order: REF:130-136 (Decision, Options with benefit/risk/cost and the always-priced smallest-change option, Recommendation, Evidence ledger, Limits, Open questions, Capture lines). The probe checks the order.
4. `root-cause` shape: REF:138, smallest-change option also priced when there is more than one change (D-07).
5. Mockup label and separation: REF:142 ("MOCKUP — proposal" in the image and the caption, files in `generated/`, never `screenshots/`).
6. Template in the reference; AR has pointers only (AR:35, :63, :80, :132, :158, :177).
7. Slug and meta: REF:140; Step 2 pointer AR:80.
8. Capture lines text only: REF:144.
9. Guide `:47` sentence and Input block (guide:47, :74-75), end section (guide:98-100), routing row (`skills/do-work-toolbox/SKILL.md:24`), help line (`skills/do-work-toolbox/actions/help.md:18`).
10. Release and the two contract scripts: finalization and the test gate.
- Default kind unchanged: `completed-work-presentation-reference.md` is not in the diff (its lines 20 and 30 are byte-identical; the probe checks both lines).
- REQ-656 revise rule: revise keeps the prior bundle's kind and copies its `ai-report-kind` meta (AR revise step 5), and keeps a non-completed-work report's section order. That matches the `--kind` rules; no contradiction.

**Merge seams (integrator):** four files conflicted with REQ-655, REQ-656 and REQ-657, resolved as a union in the merge commit. The `--kind` Input paragraph had auto-merged after the revise steps, where it read as part of the Revise form; I moved it to the top-level Input text (AR:35) and wrote "instead of a completed-work report" for "instead". The narrative checklist line (AR:175) now says a `--kind` report uses its own template, since a brief has no verdict or shipped change. Both seams go to the reviewer.

**Scope:** declared `write_set` = 5 files; touched = the same 5. No drift. Anti-bloat: no new files, fields, statuses, scripts or tests (hand-back Anti-Bloat Check).

**Re-merge after review (`3bd6a668..8c0cf7ed`, 5 files, +50/-9):** the review-fix delta (`84d23733`) changes `skills/do-work-toolbox/actions/ai-report-reference.md` (Target reading, What stays the same), `skills/do-work-toolbox/actions/ai-report.md` (When to Use bullet, Step 1 branch, checklist line) and `skills/do-work-toolbox/docs/ai-report-guide.md` (terminal-success sentence), all inside the declared write set, with the reviewer's exact text. One conflict in the When to Use list (REQ-656's revise bullet and the new `--kind` bullet) kept both. The classifier had moved past qualify, so the re-merge gate was recorded with `record-green-gate`.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `a0ebfea9` (1-minute load 1.92 before the run, no other gate running), then `advance REQ-687 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-687-probe.sh`
**Result:** ✓ Gate exit 0, wall 113 s (do-work-cli 903 Go tests, queue-kanban 420 Go tests, contract and shell suites, including `contract-regressions.sh` and `shipped-package-reference-contract.sh`). GREEN probe exit 0; `test-gate`, `scope-drift`, `run-blocked-check` and `green-gate` satisfied.

**Repository gate after the review fix:** same argv at re-merge `8c0cf7ed` (load 2.19 before, no other gate running): exit 0, wall 117 s; GREEN probe exit 0 (1 s); green record written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` (recorded revision `8c0cf7ed`).

**Red-green validation:** *(from the builder hand-back proof record; tdd: false, prose REQ)*
- `REQ-687-probe.sh` at base `bd56c4b0`: ✗ exit 1 with 21 failures (for example `ai-report.md lacks: --kind`, `ai-report-reference.md lacks: ## Proposal and Root-Cause Kinds`); the two completed-work-gate lines passed → ✓ exit 0 in the worktree at `35f984f5` and again at merge `a0ebfea9`.
- The probe also checks the part order of both templates and that `completed-work-presentation-reference.md` lines 20 and 30 are byte-identical, so the default kind still stops on an unfinished target.

**New tests added:** none (prose change; the run probe is the per-REQ check).

**Heavy verification plan:**
- Range: 3bd6a668638d72c93f0e84f446779255ebaf7224..8c0cf7ed44e66f763b8efbbb1bc182dc49bd4b9e (the plan at the first merge `a0ebfea9` selected the same single lane)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed files matched subtree `skills` (this lane also checks that ai-report stays routed in one toolbox SKILL.md row)

*Verified by work action*

## Review

**Overall: 89%** | 2026-10-10T17:38:45Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. Independent reviewer (general-purpose agent, standard Route B review of `3bd6a668..a0ebfea9`, full report `do-work/runs/work-2026-10-10-131527/REQ-687-review.md`). The reviewer recommended F1 to F5 before release; I applied all five with its exact text on the builder branch (`84d23733`) and re-merged (`8c0cf7ed`). The delta is the reviewer's text only, so I did not spawn a second reviewer; gate, probe and heavy lane re-ran green at `8c0cf7ed`.

**Important findings:**
- F1 `ai-report.md` Step 1 + `ai-report-reference.md` Target reading: `--kind` lifted only Terminal-Success Target Resolution, so the shared Archive Evidence Sweep (`completed-work-presentation-reference.md:43`, minimum record is a successful REQ) would stop a literal agent on an open REQ — impact-user-visible → fixed in `84d23733` (only Safety Load Order and Collision-Safe Publication are inherited)

**Minor findings:**
- F2 `ai-report-guide.md` terminal-success stop sentence read as covering `--kind` — impact-user-visible → fixed in `84d23733` ("Without `--kind`,")
- F3 `ai-report.md` When to Use lacked the brief — impact-user-visible → fixed in `84d23733`
- F4 `ai-report-reference.md` What stays the same: Steps 3 and 4 undefined for `--kind` — impact-user-visible → fixed in `84d23733`
- F5 `ai-report.md` checklist shared-reference line false for `--kind` — impact-negligible → fixed in `84d23733`
- F6 `ai-report.md:31` Input sentence omits the `--kind` form and says "one UR or one REQ" — impact-negligible → report only
- F7 `architecture-report.md:135` "`ai-report` ... presents completed work" now incomplete; its argument still holds — impact-negligible → report only
- F8 `ai-report-reference.md` (D-11): no rule for a completed REQ or UR passed with `--kind` — impact-negligible → report only

**Nit findings:** F9 guide `:3` and reference `:3` completed-only intros — impact-negligible → report only; F10 guide `:47` implies `proposal` accepts failed/cancelled — impact-negligible → report only; F11 template restated in the guide section and checklist line — impact-negligible → report only; F12 `help.md` new line breaks the command column — impact-negligible → report only; F13 mockup files in `generated/` vs "create `generated/` only when generated images succeed" — impact-negligible → report only
**Acceptance:** Pass — implementation and integration stages: probe exit 0 at `a0ebfea9` and `8c0cf7ed`, `git diff --check` clean, gate exit 0 at both merges; live acceptance unassessed.
**Restatement sweep:** redefined what `ai-report` accepts (unfinished work through `--kind`). Stale and fixed: `ai-report.md` Step 1 body, When to Use, checklist line; guide terminal-success sentence. Stale, report only: `ai-report.md:31`, `architecture-report.md:135`, guide `:3`, reference `:3`. Checked, not stale: `architecture-report.md:27`, `present-work.md:5`, `present-video.md:3`, toolbox `help.md:14`, core `help.md:37`, `:55`, `tutorial.md:245-252`, `:318`, `completed-work-presentation-reference.md:3,:5,:7`, the three prompt-injection JIT comments.
**Anti-bloat:** six small additions the REQ did not name (bad-value and missing-target stop rules, Step 8 summary note, checklist line, guide section, label-placement sentence, `:3` clause); no new files, fields, statuses, scripts or tests. Only the checklist line adds little (F11).
**Suggested testing:** 4 items (live `--kind proposal` on a pending REQ, `--kind root-cause` on a failed REQ, default-kind regression, catalog and revise on a brief)
**Follow-ups created:** None (8 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Putting every rule of the two kinds in one end-of-file section of `ai-report-reference.md` kept the change to pointers in `ai-report.md` and let three sibling REQs' edits merge as a plain union. The run probe checked tokens, template order and the byte-identical default gate, so each re-merge was verified in 1 to 2 s.

**What didn't:** The first wording declared the deviation from the shared reference by naming the one lifted section ("Terminal-Success Target Resolution does not apply"). That left the shared reference's other sections in force, including the Archive Evidence Sweep, which stops on any REQ without an Implementation Summary. The precedent it copied (`stakeholder-report.md:15`) says "**only** Safety Load Order and Collision-Safe Publication are inherited"; the word "only" carried the rule. Review caught it (F1).

**Worth knowing:** A three-way merge put the new `--kind` Input paragraph after the sibling's newly added Revise form steps without any conflict, where it read as part of the Revise form. Conflict hunks are not the whole merge review: read the merged section around each auto-merged insertion.

## Orientation

do-work-toolbox ai-report now has two input families: the default completed-work report (gated by `completed-work-presentation-reference.md`) and the `--kind proposal|root-cause` decision brief for open work, whose whole contract is the `## Proposal and Root-Cause Kinds` section at the end of `ai-report-reference.md`.

## Heavy Verification Plan

- Base: 3bd6a668638d72c93f0e84f446779255ebaf7224
- Target: 8c0cf7ed44e66f763b8efbbb1bc182dc49bd4b9e
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: all five changed files under `skills/do-work-toolbox/` matched subtree `skills`

## Heavy Verification Result

- Target: 8c0cf7ed44e66f763b8efbbb1bc182dc49bd4b9e; execution revision 8c0cf7ed44e66f763b8efbbb1bc182dc49bd4b9e (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-687`, `QUEUE_KANBAN_BROWSER` set, removed afterwards)
- staged-skills: executed, exit 0, 33 s (first merge `a0ebfea9`: executed, exit 0, 44 s)

## Timing

Observed 2026-10-10T17:24:58Z to 2026-10-10T17:38:45Z: 13m 47s total, 16m 33s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 8m 55s | 1 |
| verification-gate | 5m 28s | 4 |
| handback-merge | 2m 10s | 2 |

Slowest stage: review / independent review, 8m 55s, outcome success.

Notes: no builder-work event; the hand-back had landed hours before integration began (fan-out-reference Landed hand-back). The review event runs to the moment `## Review` was written, so it overlaps the review-fix re-merge, gate and drain; that overlap is why attributed time exceeds total.
