---
id: REQ-138
title: "Add Managed Text-Section Replacement"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-136]
maintenance: false
related: [REQ-135, REQ-136, REQ-137, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
write_set: [tools/replace-text-section.sh, justfile, actions/install.md, _dev/tests/contract-regressions.sh]
---

# Add Managed Text-Section Replacement

## What
Add a deterministic managed text-section utility and use it to own only do-work's generated Just recipes.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The modular installer must update generated recipe paths without reformatting or overwriting the client's unrelated Justfile content.

## Detailed Requirements
- Use exact sentinels `# >>> do-work:recipes >>>` and `# <<< do-work:recipes <<<`.
- Replace exactly one balanced managed section atomically.
- Preserve all bytes outside the section and preserve the target file mode.
- Create an absent Justfile from a supplied complete template.
- Detect and migrate the legacy five do-work recipes without duplication.
- Reject duplicate, nested, reversed, or unbalanced markers without changing the target.
- Support `justfile`, `Justfile`, and `.justfile`, including paths containing spaces.
- Repeated execution must be byte-idempotent.
- Run ShellCheck, Just parsing, and repository contract tests.

## Constraints
- Do not introduce a fictional Just attribute; this is a repository-owned sentinel convention and tested replacement utility.
- Do not touch unrelated recipes or variables.

## Dependencies
Requires REQ-136's distribution contract.

## Builder Guidance
Certainty level: Firm. Treat malformed ownership markers as a hard stop, not an invitation to guess.

## Red-Green Proof
**RED prompt/case:** Reconcile the current legacy do-work recipe block inside a Justfile that also contains a custom user recipe.
**Why RED now:** No managed-section mechanism exists, so upgrades require fragile per-recipe extraction and replacement.
**GREEN when:** The legacy block becomes one managed section, the custom recipe is byte-identical, `just --list` passes, and a second reconciliation produces no diff.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the installation decisions and exact sentinels.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
