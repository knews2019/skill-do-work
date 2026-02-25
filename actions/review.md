# Review Action

> **Part of the do-work skill.** Invoked when routing determines the user wants a post-work code review. Reviews implementation quality after work is completed.

A structured post-implementation review system that is separate from capture verification.

- **Verify** asks: "Did we capture the request correctly?"
- **Work testing** asks: "Do the implemented changes pass tests?"
- **Review** asks: "Is the implementation itself high quality and production-ready?"

## When to Use

- After one or more REQs are marked completed
- Before merging or releasing
- When the user says "review", "code review", "audit implementation", "post-work review", or "quality review"
- When the user asks for a deeper quality pass beyond test results

## Scope

Review can target:

1. A specific REQ (e.g., `do work review REQ-022`)
2. A specific UR batch (e.g., `do work review UR-005`)
3. The most recent completed REQ/UR when no target is provided

If both REQ and UR are supplied, prefer the REQ and note the choice.

## Workflow

### Step 1: Resolve Target

1. Parse explicit target from user input (REQ or UR)
2. If REQ: locate that file in `do-work/archive/` (or legacy location)
3. If UR: locate the folder in `do-work/archive/UR-NNN/`
4. If no target: choose the most recently completed REQ from archive

### Step 2: Gather Review Context

For each target REQ:

1. Read frontmatter and execution log sections (Triage, Plan, Exploration, Implementation Summary, Testing)
2. Identify files changed for that request from notes and git history
3. Read the changed code and relevant neighboring code
4. Read existing tests and any tests added by the request

### Step 3: Evaluate Quality

Assess with clear, implementation-focused criteria:

- **Correctness:** Is the implementation aligned with the request and edge cases?
- **Test quality:** Are tests meaningful, maintainable, and covering key behavior (not only happy path)?
- **Design fit:** Does the change match existing architecture and patterns?
- **Readability:** Is the code understandable and maintainable?
- **Risk:** Any regressions, perf, security, or migration risks?

Severity levels:

- **Blocker:** Must fix before merge/release
- **Major:** Should fix soon; likely to cause defects or maintenance pain
- **Minor:** Improvement suggested; optional for immediate merge

### Step 4: Produce Review Report

Output this format:

```markdown
## Code Review Report: [REQ/UR target]

**Decision:** [Approve | Approve with follow-ups | Changes requested]

### Findings

| Severity | Area | Finding | Recommendation |
|----------|------|---------|----------------|
| Blocker  | Tests | Missing regression test for null input path | Add focused regression test in ... |

### Test Assessment
- Existing tests reviewed: [...]
- New tests reviewed: [...]
- Coverage confidence: [High/Medium/Low] with 1 sentence rationale

### Risk Notes
- [Any release or maintenance risks]
```

### Step 5: Offer Follow-up

After presenting the report:

1. Ask whether to create follow-up REQs for review findings
2. If user says yes, create one REQ per actionable fix
3. Link each follow-up REQ back to the reviewed REQ/UR in the body

## Guardrails

- Do not re-run capture verification logic here; that belongs to `verify`
- Do not invent new product requirements; review implementation quality only
- Prefer concrete, file-specific comments over generic style opinions
- If evidence is insufficient, say so and downgrade confidence rather than guessing
