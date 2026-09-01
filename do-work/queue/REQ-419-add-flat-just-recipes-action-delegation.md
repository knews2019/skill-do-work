---
id: REQ-419
title: 'Add flat Just recipes, collision validation, action delegation, and compatibility aliases'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-418]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-420]
batch: go-no-llm-command-platform
---

# Add Flat Just Recipes, Collision Validation, Action Delegation, and Compatibility Aliases

## What
Expose every public command through flat Just recipes and make the existing skill aliases delegate to the canonical CLI.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Add flat recipes for all core, knowledge, toolbox, interview, memory, Dream, BKB, media, last30days, and absorbed audit-metrics public commands named by UR-081.
- Preserve existing board recipes and `run-do-work-update` as compatibility aliases.
- Implement dynamic reserved-recipe collision validation for the expanded generated section.
- Update natural-language actions, help, guides, install/update behavior, and upgrade contracts so deterministic phases delegate to `do-work-cli`.
- Missing or failed canonical tooling must stop with actionable output and never fall back to free-form mutation.

## Constraints
- The public aliases remain unchanged; the deterministic implementation becomes singular.
- Just recipes must remain flat and directly runnable without an LLM.

## Dependencies
Depends on REQ-418 (all command families must exist before publishing the complete interface).

## Builder Guidance
Certainty level: Firm. Derive the reserved recipe set from the actual managed template rather than duplicating a hard-coded list.

## Red-Green Proof
**RED prompt/case:** List and invoke every advertised recipe in a fresh install, test a custom collision for each reserved recipe source, and invoke an action with the CLI missing.
**Why RED now:** Only a handful of board/update recipes exist and natural-language actions still own deterministic mutations.
**GREEN when:** Every advertised command runs mechanically via Just, compatibility aliases work, collisions are refused dynamically, and action delegation stops actionably on canonical-tool failure.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Review Addendum — Shell-Safe Publication Recipe Arguments

REQ-413's fresh re-review found that publication results render manifest paths with Go double-quoted strings. When pasted into a shell-backed Just recipe, valid paths containing `$HOME`, `$(...)`, or backticks can expand or execute even though machine `next_argv` remains exact.

### Additional Requirements

- Render every generated publication recipe argument with shell-safe, byte-preserving quoting suitable for direct execution by the emitted flat Just recipe.
- Keep `next_argv` exact and ensure the human-runnable recipe represents the same argument bytes without expansion.
- Cover spaces, single and double quotes, dollar signs, command substitutions, backticks, tabs, and newlines where the interface permits them.

### Additional Red-Green Proof

**RED prompt/case:** Generate a publication recovery recipe for manifest paths containing shell expansions and quoting metacharacters, execute it through the advertised Just/shell boundary, and compare received argv bytes.
**Why RED now:** `strconv.Quote` emits Go double quotes, which are not shell-literal for `$`, command substitution, or backticks.
**GREEN when:** The generated recipe executes without expansion and passes byte-identical arguments for every supported hostile path shape.
