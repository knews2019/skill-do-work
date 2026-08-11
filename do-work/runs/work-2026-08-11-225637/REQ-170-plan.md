# REQ-170 Implementation Plan

## Discovery Boundary

- `prime_files` is empty, so no architecture/prime material is loaded and no prime updates are planned.
- Use the five paths already declared in `write_set` as the complete implementation surface. Do not expand into `actions/work.md`, capture templates, board parsing, metrics/trend storage, version files, changelogs, or test-source edits.
- This is instruction-contract work with `tdd: false`; prove it with the captured RED/GREEN case, focused textual inspection, and the existing contract suites rather than inventing runtime machinery.

## Single-Home Decisions

1. **Finding-closure ratchet home:** `skills/do-work/actions/work-reference.md`, next to the Testing Section Template and before later pipeline templates. Define once what counts as finding-origin work and the only two valid closure proofs: a named regression test with before-fix failure and after-fix success, or deletion of the surface where the finding lived. A bare patch is explicitly insufficient.
2. **Finding-origin detection without a new schema field:** use the existing durable provenance already available to the pipeline: `review_generated: true` for review-created follow-ups, or an explicit review/triage origin recorded in the REQ/UR (including capture's Red-Green Proof). Do not overload `review_generated`, add a parallel boolean, or add parser/board work.
3. **Earned-defense rubric home:** one paragraph within `skills/do-work/crew-members/coding-guardrails.md` § 2, “Simplicity First.” It preserves the user's exact question — “what earned this, and is the fix still cheaper than the surface it added?” — and requires the earning incident plus a cheaper-than-risk surface judgment. Keep it subordinate to the existing simplicity principle so the five-principle taxonomy, JIT comment, and caller glosses do not need expansion.
4. **Caller posture:** `capture.md`, `review-work.md`, and `validate-feedback.md` cite the canonical homes and contain only their local action/enforcement consequence. They must not create alternate definitions.
5. **No metric machinery:** no finding counters, trend log, schema additions, or other state. Convergence comes from the closure gate itself.

## Ordered Changes

### 1. Add the canonical earned-defense paragraph

**File:** `skills/do-work/crew-members/coding-guardrails.md`

- Insert one paragraph in § 2 after the existing “Simplify ≠ strip” qualification.
- Scope it to fixes that add long-lived defensive surface, name the incident/replay case that earned the defense, preserve the user's exact rubric wording, and require the remedy's added surface to be cheaper than the risk covered.
- Keep it to one paragraph and do not add a new numbered guardrail, checklist, rationalization table, success-indicator bullet, or duplicate examples.

### 2. Add the canonical finding-closure ratchet

**File:** `skills/do-work/actions/work-reference.md`

- Add a short, findable `## Finding-Closure Ratchet (Steps 6.5–7)` section immediately after `## Testing Section Template (Step 6.5)`.
- Canonically define finding-origin provenance as either the existing `review_generated: true` marker or explicit REQ/UR review/triage provenance.
- State the closure invariant once: the REQ closes only when evidence names either (a) a regression test that failed before the fix and passes after it, or (b) the deleted finding surface. Passing unrelated tests, `tdd: false`, or an untested bare patch does not satisfy the invariant.
- Keep the statement compact enough for callers to cite by section rather than paraphrase.

### 3. Sharpen capture-time proof and preserve triage provenance

**File:** `skills/do-work/actions/capture.md`

- In Step 1's Red-green proof inference bullet, add one local hook pointing to `actions/work-reference.md` → **Finding-Closure Ratchet**: when input explicitly originates in review or finding triage, the captured GREEN criterion must name the intended regression test or the exact surface to delete.
- Ensure this requirement applies independently of `tdd`; a deletion closure may legitimately use `tdd: false`, while a test closure still records the named test.
- Do not restate the ratchet's definition, add a frontmatter marker, or broaden ordinary capture behavior.

**File:** `skills/do-work-toolbox/actions/validate-feedback.md` (handoff line)

- Preserve review/triage origin in the suggested capture command by labeling the handed-off text as an accepted feedback/triage finding, so capture and later review can recognize the provenance from the durable UR/REQ text rather than infer it from a vague bug description.
- Keep Capture ≠ Execute and the existing user gate unchanged.

### 4. Enforce both rules in review

**File:** `skills/do-work/actions/review-work.md`

- In the requirements/testing review path, add a concise closure gate that points to `actions/work-reference.md` → **Finding-Closure Ratchet**. For a finding-origin REQ, verify the Testing section names matching fail-before/pass-after regression evidence or verify the diff deletes the named finding surface.
- Make missing evidence a hard bounce: record an Important finding and set Acceptance to `Fail`, which uses the existing Step 9 verdict mapping and `actions/work.md` orchestration behavior to return the REQ for remediation. Do not waive the gate because the REQ has `tdd: false`, a high score, or a generally green suite.
- Add one Verification Checklist item confirming that every finding-origin REQ passed that canonical gate. This makes skipped enforcement auditable without duplicating its definition.
- In the existing Coding-Guardrails Principle Check's “Simplicity First” item, add a one-line citation to the canonical earned-defense paragraph and require the review to flag a surface-adding defense that cannot clear it. Keep severity/report routing under the existing review machinery.

### 5. Reduce validate-feedback to triage-specific application

**File:** `skills/do-work-toolbox/actions/validate-feedback.md`

- Replace Step 4 item 5's standalone rubric explanation with a compact citation to `../do-work/crew-members/coding-guardrails.md` § 2 and retain only the triage-specific work: identify the surface-adding remedy boundary, produce the `Surface-cost` judgment, and classify direct fixes/deletions/simplifications as `N/A`.
- Preserve the user's exact question inside that one-line citation so the existing REQ-169 contract assertion remains valid while the normative explanation has one home.
- Keep Step 5's Accept bar, the `Surface-cost: N/A / Earned / Flagged` output field, and the Verification Checklist contract intact. These are action-specific enforcement/output, not an alternate rubric.
- Keep the current rule that an unearned or net-costly surface-adding remedy cannot receive plain Accept; only remove duplicated explanatory prose.

## Requirement Mapping

| REQ-170 requirement | Planned implementation | Verification |
|---|---|---|
| Ratchet stated once in work-reference | New canonical section after the Testing template | Inspect canonical wording and grep other touched files for citation-only mentions |
| Review/triage finding cannot close with a bare patch | Canonical invariant plus review hard gate that forces Acceptance Fail | Exercise the RED thought-case against the review instructions; confirm checklist and Fail route are explicit |
| Capture GREEN names test or deletion | Step 1 finding-origin proof hook; validate-feedback handoff preserves provenance | Inspect capture hook and handoff text together |
| Review sends missing evidence back, no waiver | Important finding + Acceptance Fail via existing verdict/remediation path | Verify `tdd: false`, score, and green-suite non-waivers are explicit |
| Rubric is always available during implementation | One paragraph under always-loaded coding guardrails § 2 | Confirm one paragraph and exact user wording |
| Review gate cites rubric | Extend existing Simplicity First principle check with canonical citation | Inspect review-work caller text for citation, not restatement |
| validate-feedback keeps Surface-cost and Accept contracts | Condense Step 4; retain Step 5, output field, rules/checklist | Existing REQ-169 assertions in `contract-regressions.sh` remain green |
| Minimal surface / no metric machinery | No schema, board, trend-log, action-workflow, or test-source changes | `git diff --stat`, `git diff --numstat`, and changed-path audit |
| Shipped references remain valid | Only shipped relative references; no new maintainer-doc citations | `shipped-package-reference-contract.sh` |

## Testing and Unification

1. **Record RED before editing:** use focused searches to confirm that `work-reference.md` has no finding-closure contract and `coding-guardrails.md` has no earned-defense paragraph; note that review currently allows a bare-patch finding closure.
2. **Focused content checks after editing:** verify the canonical ratchet is normative only in `work-reference.md`; verify the defensive rubric is normative only in `coding-guardrails.md`; verify each caller names its canonical section and states only its local consequence. Confirm validate-feedback still contains the surface-adding boundary, direct-fix/deletion/simplification `N/A`, non-Accept routing, and exact `Surface-cost` output token expected by REQ-169's assertions.
3. **Run contract suites from repository root:**
   - `bash _dev/tests/contract-regressions.sh`
   - `bash _dev/tests/shipped-package-reference-contract.sh`
4. **Check patch hygiene:** run `git diff --check`, `git diff --stat`, and `git diff --numstat`; inspect the complete diff of every changed file and confirm only the five declared paths changed.
5. **Account for surface:** report additions/deletions per touched file and the approximate net number of instruction lines. Keep net growth to a handful of lines by replacing validate-feedback's duplicate prose rather than layering new sections around it.
6. **Confirm explicit exclusions:** no metrics/trend log, no new schema field, no `actions/work.md` edit, no maintainer-doc citation in shipped files, and no builder-side `VERSION`/`CHANGELOG.md` edit (serial version/changelog work remains the integrator's Step 9 responsibility).

## Planned File Set

- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work-toolbox/actions/validate-feedback.md`

