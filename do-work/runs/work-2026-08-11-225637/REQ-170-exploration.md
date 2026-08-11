# REQ-170 Exploration

## Boundary

- Read the REQ, the generated implementation plan, the five declared files, and only the focused portions of the two named contract suites needed to identify current assertions and shipped-reference behavior.
- `prime_files` is empty; no global architecture or prime material was loaded.
- No implementation or test-source edits were made.

## Current RED State

- There is no finding-closure/finding-origin/bare-patch contract in any of the five planned files. The only current defensive-surface wording is in `skills/do-work-toolbox/actions/validate-feedback.md` Step 4.
- `review-work.md` currently asks for red-green evidence for behavioral changes, but this remains part of Test Adequacy and is conditional on the kind of change. It does not identify finding-origin REQs, does not accept deletion as the alternate closure proof, and does not force Acceptance to `Fail`. A finding-origin REQ can therefore close with a bare patch if its general tests and acceptance check pass.
- `capture.md` currently infers a RED/GREEN proof for `tdd: true` work and clearly behavioral fixes/features only. It does not preserve or act on review/triage provenance, and it does not require the GREEN criterion to name a regression test or deletion.
- `validate-feedback.md` currently hands off the vague text `[paste an accepted finding]`; the resulting UR/REQ can preserve the finding body but does not necessarily preserve the fact that this was accepted review/triage feedback.
- `coding-guardrails.md` § 2 currently owns general YAGNI/simplicity guidance but has no earned-defense paragraph.

## Exact File Locations and Existing Patterns

### `skills/do-work/crew-members/coding-guardrails.md`

- § 2, **Simplicity First**, begins at line 49.
- Its existing canonical-home declaration is lines 51–53: the general YAGNI statement ends with “This is the canonical statement of the principle; other files point here.”
- The planned insertion point is immediately after the single **Simplify ≠ strip** paragraph at lines 65–68 and before `## 3. Surgical Changes` at line 70.
- The file's top JIT comment says it is always loaded during Step 6 and that guardrails 1–4 plus Naming for Reach are the five-principle taxonomy. The new rubric should therefore remain an unnumbered paragraph subordinate to § 2; a sixth principle, checklist, table, or Success-indicator bullet would expand the taxonomy.
- Existing caller style is citation rather than redefinition. For example, `review-work.md` line 150 says `See crew-members/coding-guardrails.md § Simplicity First for the canonical statement.`
- The exact requested canonical question from REQ-170 is: “what earned this, and is the fix still cheaper than the surface it added?” The paragraph also needs the earning incident/replay case and the cheaper-than-covered-risk judgment, without expanding into the triage-specific `Surface-cost` verdict mechanics.

### `skills/do-work/actions/work-reference.md`

- The canonical insertion seam is directly after `## Testing Section Template (Step 6.5)`, lines 564–583, and before `## Deferred Prime-Link Path Computation (Step 7.5)`, line 585.
- The existing Testing template already has the evidence shape needed by the test branch:
  - `**Tests run:**` and `**Result:**` at lines 569–570.
  - Named fail-before/pass-after rows under `**Red-green validation:**` at lines 572–574.
  - New and cross-REQ test lists at lines 576–580.
- No template field currently names a deleted surface. The canonical section can define deletion as the alternate proof without growing the template: capture names the intended surface, and review checks that named surface against the diff.
- Existing reference convention within core is `` `actions/work-reference.md` → **Exact Section Heading** ``. The planned heading should therefore be stable and findable as `## Finding-Closure Ratchet (Steps 6.5–7)`, while callers cite **Finding-Closure Ratchet** rather than restating it.
- The only existing durable review-follow-up marker is `review_generated: true`; it is already described and consumed in `review-work.md` Step 10. Explicit review/triage origin in REQ/UR prose (including the Red-Green Proof) is the non-schema counterpart. No parser, alias row, or new boolean is needed.
- Canonical content needs to distinguish the invariant from ordinary TDD: finding origin is `review_generated: true` or explicit REQ/UR review/triage provenance, and closure evidence is either a named regression test with before-fix failure plus after-fix success, or deletion of the finding surface. Unrelated green tests, `tdd: false`, and a bare patch are explicit non-evidence.

### `skills/do-work/actions/capture.md`

- Step 1 starts at line 96. The smallest local hook belongs in the existing **Red-green proof inference** bullet at line 103, immediately after the long **TDD assessment** bullet at line 102.
- Current semantics matter:
  - `tdd: true` is limited to a realistically runnable test-first case.
  - `tdd: false` still requires proof for behavior-changing work.
  - The current proof-inference bullet is grammatically scoped to `tdd: true` requests and clearly behavioral fixes/features.
- The finding-origin hook must be written as an additional independent case, not nested under that existing scope. Otherwise doc/config/rule findings and deletion closures with `tdd: false` would escape the ratchet.
- Caller wording should point to `` `actions/work-reference.md` → **Finding-Closure Ratchet** `` and state only capture's local action: when the input explicitly comes from review/finding triage, GREEN names the intended regression test or exact surface to delete. It should not redefine provenance or the closure invariant.
- Step 5 currently requires a Red-Green Proof for `tdd: true` and adds it for behavior-changing work when meaningful (lines 254–264). The Step 1 hook can make the finding-origin proof meaningful without adding a template/frontmatter field.
- Integration warning: this file already has an unrelated uncommitted diff at lines 207–246 (REQ-168 screenshot-copy hardening: unique `mktemp` copy). Any REQ-170 integration must preserve that change and stage/merge only its own Step 1 hunk.

### `skills/do-work/actions/review-work.md`

- The requirements/testing enforcement path is split across:
  - Step 5 Requirements Check, lines 75–91.
  - Step 6 Test Adequacy, lines 102–109.
  - Step 7 Acceptance Testing and the Acceptance enum, lines 169–202.
  - Step 9 verdict mapping, line 299.
  - Verification Checklist, lines 520–532.
- The closest existing evidence check is Test Adequacy line 107. It asks for red-green validation for bug fixes/features but explicitly allows alternate proof for non-behavioral work. It is advisory/scored, not a closure gate.
- The hard-bounce consequence cannot be expressed only as an Important finding. Current Step 9 says a high-enough score plus Acceptance Pass can still produce **Approve**, and Important findings ordinarily flow to future follow-ups in Step 10. The new local hook must explicitly record an Important finding **and set Acceptance to `Fail`**; line 299 then deterministically maps that result to **Request changes**, independent of score.
- The gate should cite `` `actions/work-reference.md` → **Finding-Closure Ratchet** `` and inspect:
  - finding origin from `review_generated: true` or explicit REQ/UR provenance;
  - a matching named fail-before/pass-after test in the Testing section, or deletion of the named finding surface in the diff.
- The non-waiver needs to remain explicit because three current paths would otherwise weaken it: `tdd: false`, Test Adequacy `N/A`, and generally green/unrelated tests. The gate is independent of all three.
- The best existing location for earned-defense review is Step 6's **Coding-Guardrails Principle Check**, specifically Simplicity First at line 150. That bullet already cites the canonical § 2 and flags speculative/defensive handling, so it only needs the local review consequence: a surface-adding defense that cannot clear the canonical paragraph is a finding. Do not reproduce incident/cost/test criteria here.
- The principle pass is currently described as informational and normally routes anything it uniquely catches as Minor (line 155). If the intended earned-defense failure must be stronger than Minor, the implementation must say so explicitly; otherwise the current paragraph controls severity. REQ-170 only requires the gate to cite/flag the rubric, not a separate hard-fail for every rubric miss.
- Add the audit item to the Verification Checklist near the Acceptance item at line 529: every finding-origin REQ passed the canonical gate. Citation-only wording avoids a second definition.
- Existing review provenance pattern is Step 10's follow-up frontmatter (`review_generated: true`, line 365) and the generation-2 marker-only rule (around lines 398–403). The new ratchet may read this marker but must not alter its existing cascade-depth meaning.

### `skills/do-work-toolbox/actions/validate-feedback.md`

- Step 4 item 5, line 69, is the duplicated normative rubric that should be condensed.
- Triage-specific behavior that must remain:
  - Explicit surface-adding boundary: `guard, fallback, retry, validation layer, rule, or warning apparatus`.
  - `Surface-cost: N/A` for direct bug fixes, deletions, and simplifications.
  - Step 5's Accept bar at line 80: an unearned/net-costly added defense cannot receive plain Accept and routes to Push back/Discuss.
  - Output field at line 95: `**Surface-cost:** N/A / Earned / Flagged`.
  - Rules hook at line 126 and Verification Checklist item at line 152.
- The canonical reference from this toolbox action should follow the established semantic cross-package style: `` `../do-work/crew-members/coding-guardrails.md` § 2, Simplicity First ``. Keep only triage's boundary, `Surface-cost` judgment, N/A classification, and Accept consequence here.
- The capture handoff is at lines 113–116. Label the pasted payload as accepted feedback/triage finding so the durable UR/REQ text carries provenance, while leaving `Capture ≠ Execute`, the read-only rule, and the explicit user-triggered `do-work capture-request:` handoff intact.

## Current Contract-Test Assertions

`_dev/tests/contract-regressions.sh` resolves `actions/validate-feedback.md` to the toolbox package, then REQ-169 asserts all five of these regex contracts (lines 1950–1976):

1. `what incident earned this, and is the fix still cheaper than the surface it added\?`
2. `guard, fallback, retry, validation layer, rule, or warning apparatus`
3. `Direct bug fixes, deletions, and simplifications.*N/A`
4. `must not receive.*Accept.*Push back.*Discuss`
5. `\*\*Surface-cost:\*\* N/A / Earned / Flagged`

Consequences:

- Condensing line 69 is safe only if every regex still appears somewhere in `validate-feedback.md`; Step 5/output/checklist already satisfy parts 4–5, but parts 1–3 currently live only in Step 4 item 5.
- There is a wording conflict to reconcile: REQ-170 requires the canonical exact question “what earned this, and …,” while the existing REQ-169 assertion requires `what incident earned this, and …` in the toolbox caller. Because test-source edits are out of scope, the toolbox file must retain the older asserted phrase in citation/application text while the canonical paragraph preserves the REQ-170 wording. Keep the extra occurrence citation-sized so it is not a competing normative explanation.
- There are no current contract-regression assertions for the new ratchet heading, capture hook, review hard gate, or core earned-defense paragraph. The planned focused searches are therefore required in addition to running the suites.
- The maintainer-document scan checks shipped files for `CLAUDE.md`/`AGENTS.md`; `actions/capture.md` and `actions/validate-feedback.md` are already per-file allowlisted because they discuss consumer-project decision stores. REQ-170 should add no new maintainer-doc citation.

`_dev/tests/shipped-package-reference-contract.sh` reads `suite/modules.tsv`, maps all four sibling skill sources to separate installed destinations, and validates Markdown links in both source and installed topology. It ignores backticked prose citations as link targets. Relevant conventions:

- Core-to-core semantic references use skill-root-relative backticked paths such as `` `actions/work-reference.md` `` and `` `crew-members/coding-guardrails.md` ``.
- Toolbox-to-core semantic references already use `` `../do-work/actions/...` `` in backticks throughout `validate-feedback.md` and other toolbox actions.
- If the new citation is made a real Markdown link instead of backticked prose, its path must be document-relative: from `skills/do-work-toolbox/actions/validate-feedback.md` the valid source/installed target is `../../do-work/crew-members/coding-guardrails.md`; `../do-work/...` is only the established skill-root-relative prose convention and would fail as an actual link.
- The allowed reference direction is toolbox → core; `suite/modules.tsv` installs both as sibling skill directories, so that reference resolves in the shipped topology.

## Integration Concerns and Reconciliation Checklist

1. Preserve the unrelated existing `capture.md` working-tree hunk; REQ-170 owns only the Step 1 proof bullet in that file.
2. Preserve all five REQ-169 regex tokens in `validate-feedback.md`, especially the old “what incident earned this” phrase, despite moving the canonical REQ-170 wording to core.
3. Keep exactly one normative ratchet definition (work-reference) and one normative defense rubric paragraph (core coding-guardrails). Capture/review/validate-feedback should cite and state only their local consequence.
4. Make the capture finding-origin clause independent of `tdd` and of “clearly behavioral” scope.
5. Make review's failure consequence explicit as both Important and Acceptance `Fail`; Important alone creates a follow-up and does not bounce a high-scoring change.
6. Preserve deletion as a first-class closure proof. It does not need a new Testing-template field, but the named deletion from capture/REQ must match the diff at review.
7. Preserve the five-principle coding-guardrails taxonomy and one-paragraph limit.
8. Keep validate-feedback's surface boundary, N/A classification, Accept routing, output token, read-only posture, and capture user gate byte-for-meaning even when condensing prose.
9. Use stable exact section names in citations: **Finding-Closure Ratchet** and **Simplicity First**.
10. Run both existing suites, focused canonical-home/restatement greps, `git diff --check`, and a five-path/net-line audit. No test-source, schema, metrics, trend-log, `actions/work.md`, version, or changelog edit belongs in the builder patch.
