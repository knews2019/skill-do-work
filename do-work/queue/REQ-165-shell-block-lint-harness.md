---
id: REQ-165
title: Shell-block lint harness for shipped action files
status: pending
created_at: 2026-08-11T11:46:50Z
user_request: UR-036
domain: testing
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-166, REQ-167, REQ-168]
batch: stabilization-audit
write_set: [_dev/tests/action-shell-blocks.sh]
---

# Shell-Block Lint Harness for Shipped Action Files

## What

Add a `_dev/tests/` check that extracts every fenced shell block (` ```bash ` / ` ```sh `) from the shipped `skills/` tree — action files, crew files, and shipped hook scripts — and lints them: `bash -n` (syntax) always, `shellcheck` when available. This makes the suite's largest defect generator (prose-prescribed shell that nothing executes until an agent runs it in a consumer repo) testable in CI.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

User goal: stabilize the skill so reviews stop returning 3–5 findings per pass. Every entry in CLAUDE.md's "Prescribed Shell Commands" trap list originated on this untested surface; closing the class beats fixing instances.

## Context

- Follow the existing harness conventions in `_dev/tests/` (e.g. `contract-regressions.sh`, `shipped-package-reference-contract.sh`) — same invocation style, same failure-reporting style (name the offending file and the fix).
- Blocks containing prose placeholders (`<suite-root>`, `REQ-NNN`, `[slug]`, `{n}`) are expected — the harness must handle them (substitute dummies, or lint with placeholders neutralized) rather than skip those blocks wholesale; a blanket skip would neuter the check exactly where it matters (see CLAUDE.md's skip-list trap).
- `shellcheck` may be absent — degrade to `bash -n` with a note, don't fail the suite over a missing linter.
- Repo-only: lives in `_dev/tests/`, does not ship with any package.

## Detailed Requirements

- Extract every fenced `bash`/`sh` block across the shipped `skills/` tree, tracking source file and line so failures are attributable.
- `bash -n` each block; run `shellcheck` on each block when the binary is present.
- A deliberately broken fixture block (or self-test mode) proves the harness actually fails on bad input — a checker that cannot fail is decoration.
- Wire it into however `_dev/tests/` is normally run, alongside the existing contract tests.
- Findings style: name the file, the line, and the diagnostic.

## Builder Guidance

Certainty: Firm on intent, exploratory on mechanics. The extraction mechanics (awk/sed vs. a small script) are the builder's choice. Expect the first run to surface real pre-existing findings in shipped blocks — fix trivial ones in-scope, report substantive ones as follow-up candidates rather than ballooning this REQ.

## Red-Green Proof

**RED prompt/case:** Today no `_dev/tests/` check reads the fenced shell blocks in shipped action files — a block with a syntax error (or the `set -euo pipefail` dead-fallback pattern in `session-start.sh`) sails through the whole suite.
**Why RED now:** The prescribed-shell surface is prose; nothing executes or parses it before an agent hits it in a real repo. Nine documented traps in CLAUDE.md all shipped this way.
**GREEN when:** The new check runs with the suite, fails naming file + line when a seeded bad block is introduced, and passes on the clean tree (after any pre-existing findings are fixed or dispositioned).
**Validation:** Inferred during capture (plan discussed and endorsed in-session).

## Full Context

See `do-work/user-requests/UR-036/input.md` for complete verbatim input.

---
*Source: "do-work capture-request for audit and fix to simplify and make it robust" (UR-036)*
