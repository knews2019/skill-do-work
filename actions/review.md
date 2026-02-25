# Review Action

> **Part of the do-work skill.** Invoked when routing determines the user wants a post-work code review. Evaluates implementation quality after coding + testing are done.

A focused code-review pass that is distinct from `verify`:
- **Verify** checks whether REQs accurately capture the original user request.
- **Review** checks whether completed code changes are high quality and production-ready.

## When to Use

- After `work` has completed one or more REQs
- Before shipping or merging implementation commits
- When the user says "review", "code review", "review code", or "audit code"

## Workflow

### Step 1: Resolve Target Scope

1. **If user specifies REQ** (e.g., `review REQ-021`): review that request only
2. **If user specifies UR** (e.g., `review UR-010`): review all completed REQs in that UR
3. **If no target specified**: use the most recently completed REQ (prefer `do-work/archive/`)

If the target has not been completed yet (`status != completed`), report that and ask whether to review current in-progress changes anyway.

### Step 2: Gather Implementation Context

For each target REQ:

1. Read the REQ file (triage notes, implementation summary, testing section)
2. If `commit` exists, inspect that commit and files changed
3. If no `commit` exists, inspect current working tree diff relevant to the REQ
4. Identify and run relevant tests/lint/type-check commands when available

### Step 3: Evaluate Quality Dimensions

Score and comment on:

1. **Correctness**
   - Does implementation satisfy REQ requirements?
   - Any obvious logic bugs, edge-case failures, or regressions?

2. **Test Quality**
   - Are existing tests relevant and passing?
   - Were new tests added for new behavior and regressions?
   - Are assertions meaningful (not just snapshot/no-op checks)?

3. **Code Quality**
   - Readability, cohesion, naming, and maintainability
   - Duplication or dead code
   - Error handling and input validation

4. **Safety & Compatibility**
   - Backwards compatibility concerns
   - Security/privacy issues
   - Performance risks for hot paths

5. **Traceability**
   - Do REQ notes match what was actually shipped?
   - Are testing notes reproducible and specific?

### Step 4: Severity Classification

Classify findings as:
- **Blocker**: must fix before merge/release
- **Major**: high-value fix strongly recommended now
- **Minor**: quality improvement, can be deferred
- **Nit**: optional stylistic suggestion

### Step 5: Produce Review Report

Use this format:

```markdown
## Code Review Report: <target>

**Recommendation:** Approve | Approve with follow-ups | Changes requested

### Summary
- [1-3 bullets]

### Findings
- **Blocker:** ...
- **Major:** ...
- **Minor:** ...
- **Nit:** ...

### Test Assessment
- Commands run: ...
- Result: ...
- Gaps: ...

### Next Actions
1. ...
2. ...
```

If there are no findings, explicitly say that and still provide test evidence.

### Step 6: Offer Auto-Fix Pass

After presenting the report, ask whether to apply fixes now.

- If yes: make changes, run tests again, and post a follow-up report with updated recommendation.
- If no: leave code untouched and record the review outcome only.

## What NOT To Do

- Don't re-run REQ extraction checks here (that's `verify`)
- Don't invent new product requirements beyond REQ scope
- Don't mark approval without at least one concrete test/check result (or a clear environment limitation)
