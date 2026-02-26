# Verify Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to verify the quality of captured requests. Evaluates REQ files against their originating User Request (UR) to find gaps.

A confidence evaluation system that compares extracted REQ files against the original user input to identify lost requirements, dropped UX details, missing intent signals, and incomplete coverage.

## Philosophy

- **The original input is the source of truth** — the UR's input.md contains everything the user said
- **REQs should be lossless extractions** — every requirement in the input should appear in at least one REQ
- **Intent signals matter** — not just WHAT was requested, but HOW firmly and with what scope guidance
- **Actionable output** — don't just report problems, offer to fix them

## When to Use

- After creating REQs from a complex request (to validate the extraction)
- When the user says "verify", "check", "evaluate", "review requests"
- Before starting `do work` processing, as a quality gate

## Workflow

### Step 1: Find the Target UR

1. **If user specifies a UR** (e.g., "verify UR-003"): Use that UR directly
2. **If user specifies a REQ** (e.g., "verify REQ-018"): Read the REQ's `user_request` field to find the UR
3. **If no target specified**: Find the most recent UR folder in `do-work/user-requests/` (highest UR number)

**Legacy support:** If the user points to a REQ with `context_ref` instead of `user_request`, read the referenced CONTEXT file from `do-work/assets/` and use its verbatim input as the source of truth.

### Step 2: Read the Original Input

1. Read `do-work/user-requests/UR-NNN/input.md`
2. Extract the full verbatim input section
3. Note the `requests` array to know which REQs to evaluate
4. Note any Batch Constraints section

### Step 3: Read All Related REQs

1. Find all REQ files listed in the UR's `requests` array
2. Check `do-work/`, `do-work/working/`, and `do-work/archive/` for each
3. Read the full content of each REQ file

### Step 4: Evaluate Each REQ

For each REQ, score it on these dimensions:

**Requirements Coverage (0-100%)**
- Does the REQ capture all requirements from the original input that apply to this feature?
- Are specific values, constraints, and conditions preserved?
- Are edge cases and error handling requirements included?

**UX/Interaction Details (0-100%)**
- Are interaction behaviors captured? (e.g., "auto-scroll to current file," "collapse on click")
- Are visual/layout requirements noted?
- Are state transitions described?

**Intent Signals (0-100%)**
- Does the Builder Guidance section (if applicable) accurately reflect the user's tone?
- Is the certainty level correct (exploratory vs firm)?
- Are scope cues preserved ("keep it simple," "don't over-build")?

**Batch Context (0-100%)** — only for multi-REQ batches
- Do cross-cutting constraints from the UR appear in this REQ's Constraints section?
- Are sequencing requirements noted?
- Are shared design principles captured?

### Step 5: Identify Gaps

For each gap found:
1. Quote the relevant section from the original input
2. Identify which REQ should contain it (or if a new REQ is needed)
3. Classify the severity:
   - **Critical**: A firm requirement that was completely dropped
   - **Important**: A clear requirement that was partially captured or summarized too aggressively
   - **Minor**: A passing mention or soft preference that was missed
   - **Ambiguous**: The original input doesn't contain enough information to resolve this — neither the REQ nor the UR has a clear answer. This isn't a gap in the REQ; it's a gap in the original request that only the user can fill.

### Step 6: Generate Report

Output a confidence report in this format:

```
## Verification Report: UR-NNN

**Overall Confidence: [X]%**

### Per-REQ Scores

| REQ | Title | Coverage | UX Detail | Intent | Batch | Overall |
|-----|-------|----------|-----------|--------|-------|---------|
| REQ-018 | TOC Panel | 85% | 70% | 90% | 80% | 81% |
| REQ-019 | File Tree | 90% | 60% | 90% | 80% | 80% |

### Gaps Found

**Critical:**
- [None / list of dropped firm requirements with source quotes]

**Important:**
- [List of partially captured or over-summarized requirements]

**Minor:**
- [List of missed passing mentions]

**Ambiguous (needs client input):**
- [List of requirements where the original input is unclear — these become Open Questions on the REQ]

### Recommendations

1. [Specific fix: "Add 'auto-scroll to current file' to REQ-018 Detailed Requirements"]
2. [Specific fix: "Add batch constraint about stability-first sequencing to REQ-019"]
```

### Step 7: Offer Fixes

After presenting the report:

1. Ask the user if they want to apply the recommended fixes
2. If yes, update the REQ files directly:
   - **Critical/Important/Minor gaps**: Add missing requirements to the appropriate sections, add or update Builder Guidance sections, add batch constraints to Constraints sections
   - **Ambiguous gaps**: Don't fix the REQ content — instead add an Open Question to the REQ's `## Open Questions` section (create the section if it doesn't exist). Use the choice format:
     ```
     - [ ] [Question]
       Recommended: [best default based on context]
       Also: [alternative A], [alternative B]
     ```
     The recommended choice lets the builder proceed with best judgment if the user doesn't answer before the work action picks it up.
3. Re-score after fixes to confirm improvement (Ambiguous items don't affect the re-score — they're resolved by the user during work, not by verify)

## Scoring Guidelines

**90-100%**: Excellent — all requirements captured with full detail. Ready to build.
**75-89%**: Good — minor gaps that probably won't affect implementation. Fix if convenient.
**50-74%**: Needs attention — important requirements or interaction details missing. Fix before building.
**Below 50%**: Significant gaps — major requirements dropped. REQ needs substantial rework.

## Legacy REQ Handling

For REQs created before the UR system:
- They won't have `user_request` in frontmatter
- They may reference `assets/CONTEXT-*.md` via `context_ref`
- They won't have a Builder Guidance section
- Score them the same way, but note that missing Builder Guidance is expected (not a gap) for legacy REQs
- If the user wants to verify legacy REQs and has the original CONTEXT file, use its verbatim input

## What NOT To Do

- Don't expand requirements beyond what the user said — you're checking coverage, not inventing new features
- Don't penalize REQs for missing details the user never mentioned
- Don't treat implementation details as gaps — those are for the builder to decide
- Don't classify something as Ambiguous when the answer is in the original input — that's a Critical or Important gap. Ambiguous means the *user's input itself* doesn't contain the answer.
- Don't block on verification — it's advisory, not a gate (unless the user wants it as a gate)
