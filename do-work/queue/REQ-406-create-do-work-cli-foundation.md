---
id: REQ-406
title: 'Create the shared do-work-cli runtime and Git transaction foundation'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Create the Shared Do-Work-CLI Runtime and Git Transaction Foundation

## What
Create the suite-wide `do-work-cli` Go module and its shared execution contracts.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Put one `do-work-cli` module under the installed core package and require Go 1.26.1+.
- Support global `--repo-root` and `--format text|json` options.
- Build on demand when the binary is missing or older than its Go sources, using the standard library unless an existing dependency is demonstrably necessary.
- Define one typed result model that renders both text and stable JSON with schema version, command, outcome, repository root, findings, changes, skipped work, and rollback result.
- Give every finding a stable code, severity, affected IDs/paths, observed evidence, fixability, automation-stop reason, exact next argv/Just recipe, and verification command.
- Enforce exit codes 0–4 exactly as specified by the UR.
- Implement Git target preflight, dry-run support, optional exact-path commit, rollback, committed-path verification, and committed-state risk reporting.

## Constraints
- Mutations require Git; read-only commands remain usable outside Git.
- Dirty target paths are refused while unrelated dirty paths remain allowed; `--commit` requires an empty pre-existing index.
- Pre-commit rollback restores tracked targets and removes only invocation-created paths; post-commit failure reports `git revert <sha>` and never rewrites history.

## Dependencies
Foundation REQ; all later batch REQs depend on it directly or transitively.

## Builder Guidance
Certainty level: Firm. Establish narrow packages and behavioral contracts that later migrations can consume without parallel implementations.

## Red-Green Proof
**RED prompt/case:** Invoke the absent `do-work-cli` in text and JSON modes and exercise a mutation against dirty, clean, rollback, and commit fixtures.
**Why RED now:** No shared CLI, typed result schema, build-cache launcher, or common Git transaction layer exists.
**GREEN when:** The same typed result drives stable text/JSON, documented exit codes are observed, and Git transaction fixtures prove exact-path refusal, rollback, and commit behavior.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
