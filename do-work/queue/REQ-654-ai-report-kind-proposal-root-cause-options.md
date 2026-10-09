---
id: REQ-654
title: 'ai-report --kind proposal, root-cause and options writes a decision-first brief for unfinished work'
status: pending
created_at: 2026-10-09T21:10:00Z
user_request: UR-144
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-655, REQ-656, REQ-657]
batch: ai-report-modes
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/ai-report-reference.md", "skills/do-work-toolbox/actions/completed-work-presentation-reference.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md"]
---
# ai-report --kind Proposal, Root-Cause and Options Writes a Decision-First Brief for Unfinished Work
## What
`do-work-toolbox ai-report` gains `--kind proposal|root-cause|options <topic|REQ-NNN|UR-NNN>`. For these kinds only, the "completed work only" gate is lifted: the target may be a free topic, an open REQ or an open UR, read as open work and never presented as shipped. The report has a fixed decision-first shape. Without `--kind`, `ai-report` behaves exactly as today.
## Why
Report A1 (UR-144 input, "What happened" table): about 7 sessions in 4 consumer repos between 2026-09-22 and 2026-10-07 asked for proposal, options, "backup options" and "why it happened and what to do" reports. One plan recorded the pushback "the toolbox ai-report action only presents *completed* UR/REQ work, so it doesn't fit a proposal", and the user asked for the mockups anyway. Each report was improvised from scratch.
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/ai-report.md:18`: "The user wants a detailed presentation of one completed UR or REQ." `:26`: "The target is unfinished or unsuccessful; report its status instead of presenting it as shipped." `:30`: "`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank."
- `ai-report.md` Step 5 (report cites `:89-101`) fixes a narrative that opens with "Verdict" and "What Shipped". A proposal has neither.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:20` rejects `cancelled`, `failed` and every unfinished status; `:30` forbids falling back to `do-work/queue/`, `do-work/working/` or an active UR body "to make an unfinished target appear complete". Both stay for the default kind.
## Detailed Requirements
1. `ai-report.md` Step 1 gets a short branch: when `--kind` is `proposal`, `root-cause` or `options`, skip Terminal-Success Target Resolution and accept a topic, an open REQ or an open UR. Queue and working REQs are read as open work, never as shipped work.
2. Keep the safety load order from `completed-work-presentation-reference.md` for the new kinds: `prompt-injection.md` and `anti-slop.md` load first.
3. The fixed shape for `proposal` and `options`, in this order:
   1. the decision the reader must make, first;
   2. two to four options (O1..On), each with benefit, risk and cost;
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
- The default kind keeps the completed-work gate exactly as today, including `completed-work-presentation-reference.md:20` and `:30`.
- Out of scope: a free-form brainstorm (`deep-explore` owns it), running the capture lines, committing or publishing the report.
- Existing bundles are immutable (`completed-work-presentation-reference.md:71`, `:73`).
- This REQ is its own release; do not fold REQ-655, REQ-656 or REQ-657 into it.
## Assumptions (recorded at capture, no questions asked)
- `options` uses the same decision-first template as `proposal`; the report gives no separate shape, so the two differ only in the kind label and the meta value. If the builder finds a real difference worth keeping, record it in Decisions.
- The report writes the capture line as `/do-work capture-request ...`, a consumer slash-command spelling. The builder writes the invocation form the toolbox already uses in its other printed next-step lines (`do-work capture-request: <task>` per `skills/do-work/SKILL.md:46`).
- A free topic has no REQ or UR id; the slug topic part is a short kebab form of the topic text.
- "Open REQ" means any non-terminal status (`pending`, `pending-answers`, `blocked`, claimed). A `failed` or `cancelled` REQ is accepted for `root-cause`, where explaining a failure is the point, and is never described as shipped.
- The "kept apart like generated images" rule reuses the existing generated-image separation in `ai-report-reference.md`; no new image pipeline.
## Dependencies
None upstream. REQ-655 (index and find) reads the kind meta and slug this REQ defines, but works without it (it falls back to unknown kind), so there is no ordering edge.
## Builder Guidance
Certainty is high on the shape (the report lists it item by item). Latitude: exact wording, where the Step 1 branch sits, and whether `options` and `proposal` share one template block.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox ai-report --kind proposal "backup options"` in a repo with no completed work.
**Why RED now:** `ai-report.md:30` accepts only `UR-NNN`, `REQ-NNN`, `most recent` or blank, and `:26` stops on unfinished targets, so there is no path to a proposal report.
**GREEN when:** The action produces a bundle named `yyyy-mm-dd_hhmm_proposal-<topic>` whose first section is the decision, with 2 to 4 options each showing benefit, risk and cost, a recommendation, an evidence ledger, a limits section, Q1..Qn, and one capture line per option; its `<head>` carries `ai-report-kind`; every mockup carries "MOCKUP — proposal". And `ai-report REQ-NNN` on an unfinished REQ, without `--kind`, still stops exactly as today.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing action routing and status contracts, which this REQ does; family `canonical-authoring-vs-tolerant-reading` fits a kind that reads open work without presenting it as shipped.
## Full Context
See `do-work/user-requests/UR-144/input.md` for complete verbatim input (sections Request A1, What happened, Where the behaviour lives today, Proposed direction A1, Acceptance check). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A1: "`ai-report --kind proposal|root-cause|options <topic|REQ-NNN|UR-NNN>` for work that is not completed. Lift the "completed work only" gate for these kinds only."*

## Addendum (2026-10-10)

User added (through the ask tool, recorded in UR-150):

> ```
> Addendum to REQ-654 (ai-report --kind proposal, root-cause and options): every decision report must include the smallest-change option.
> 
> Context: the REQ-651 report on the board's activity correlation (ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism) listed two options from its brief, "stamp the [REQ-NNN] prefix then delete the scaffold" and "keep the code". The option that was chosen and shipped in 0.305.88, "delete the scaffold and stamp nothing", was not in the report. The maintainer was asked through the ask tool: "Decision reports today list only the options the brief named. Should that become a rule for the queued proposal-report feature (REQ-654)?" and answered "Yes, add it to REQ-654".
> 
> Constraint to add: for the proposal and options kinds, the option list always includes the smallest-change option, usually "delete the mechanism" or "do nothing extra", priced like every other option (benefit, risk, cost), even when the brief or the topic did not name it. The report may recommend against it, but it must be on the list.
> ```

- Requirement 3's option list (and the `options` kind that shares it) always includes the smallest-change option, usually "delete the mechanism" or "do nothing extra", with benefit, risk and cost like every other option, even when the brief or topic did not name it. The report may recommend against it; it may not omit it.
- Source: the REQ-651 report listed "stamp then delete" and "keep" and missed "delete, stamp nothing", which is what 0.305.88 shipped.
