---
id: REQ-429
title: '[impact-rule-change] Review fix: Complete normalized schema-field projection'
status: pending
domain: backend
created_at: 2026-08-30T19:21:33Z
user_request: UR-081
addendum_to: REQ-408
review_generated: true
impact: impact-rule-change
effort_estimate: effort-mechanical
tdd: true
sweep: true
sweep_key: normalized-schema-projection-completeness
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Complete Normalized Schema-Field Projection

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Make the shared typed request record project every field governed by the Schema Read Contract. Done means the class cannot recur: a table-driven completeness test must fail whenever a contracted field lacks normalized value and `FieldResult` evidence.

## Context
Found during review of REQ-408. The remediation projected the specifically listed downstream fields but omitted `caveman`, leaving consumers to normalize it ad hoc despite the one-model contract.

## Requirements
- Add normalized `caveman` value and explicit normalization evidence to `RequestRecord` and `TypedRecord`.
- Replace the manually incomplete coverage assertion with a table-driven completeness ratchet over every Schema Read Contract field.
- Preserve generic parser evidence alongside normalized typed evidence.

## Instances
- [ ] `caveman` is defined by `schemanormalization` but absent from the typed request record.
- [ ] `TestTypedRecordProjectsEveryNormalizedSchemaField` does not enumerate every contracted field, so its name overstates its protection.

## Red-Green Proof
**RED prompt/case:** Add `caveman: light` to the every-normalized-field fixture and require typed value `lite` plus recognized `FieldResult` evidence; also make the coverage table enumerate all schema contracts.
**Why RED now:** `TypedRecord` never normalizes or projects `caveman`, and the existing test omits it.
**GREEN when:** Every contracted schema field has typed normalized evidence and the completeness ratchet prevents future omissions.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
