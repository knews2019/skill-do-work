---
id: REQ-557
title: '[impact-negligible] Deduplicate six Go helper names defined fourteen times across do-work-cli'
status: pending
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-550, REQ-552]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go, skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go, skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go, skills/do-work/tools/do-work-cli/internal/publication/capture_files.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go, _dev/tests/audit-lockins.sh]
---

# Deduplicate six Go helper names defined fourteen times across do-work-cli

## What
`uniqueSorted`, `subtractPaths`, `requestIDLess`, `firstError`, `compareSemver` and `physicalPath` are defined fourteen times across `internal/`; in every case the duplicating package already imports the package holding an earlier copy, and three names have copies that disagree. Export one canonical definition per name in the lowest already-imported package (`corehelpers` for the path and error helpers, `repositorymodel` for `requestIDLess`), delete the other eight, and record the three semantic reconciliations as named decisions in this REQ.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Per-REQ helper files that duplicate an existing helper are the agent-creep class; the three semantic splits (`uniqueSorted` drops empty strings in one copy, `compareSemver` accepts unparseable input in one copy, `physicalPath` has two contracts) are latent correctness drift.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 4, sweep_key `per-req-duplicate-go-helpers`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -70. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `internal/finalization/finalization_prepare.go` — `subtractPaths` and `uniqueSorted` duplicate `internal/corehelpers/checks.go` (introduced 761d8e6a, REQ-498; `finalization` already imports `corehelpers`).
- `internal/knowledgecommands/interview_commands.go` — third `uniqueSorted`, silently drops empty strings (01d920dd).
- `internal/dependencygraph/dependency_graph.go` and `internal/repositorymodel/repository_model.go` — two `requestIDLess` bodies in one commit (ac2e3acd, REQ-408) with two different number parsers; `internal/nextselection/next_types.go` — third `requestIDLess` (625d49aa, REQ-411; `nextselection` already imports `repositorymodel`).
- `internal/publication/capture_files.go` — `firstError` byte-identical to `corehelpers/checks.go` (cf111a50, REQ-413).
- `internal/knowledgecommands/interview_commands.go` — `compareSemver` returns 0 for unparseable input while `internal/publication/release.go` rejects it.
- `internal/suiteinstall/update_transaction.go` — `physicalPath` is `EvalSymlinks`+`Abs`; `internal/knowledgecommands/commands.go` walks missing ancestors; same name, different contract.
- Each reconciliation (empty-string handling, unparseable semver, missing-ancestor paths) is written into this REQ's Implementation Summary as a decision with the behaviour each caller keeps; a silent pick is a review refusal.
- Reproduce at dc8a64e3 (prints 14 lines): `rg -n --glob '*.go' --glob '!*_test.go' '^func (uniqueSorted|subtractPaths|requestIDLess|firstError|compareSemver|physicalPath)\(' skills/do-work/tools/do-work-cli/internal/ | sort`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- No import cycles introduced: the canonical home is always a package the duplicator already imports.
- Tests unchanged except where a test named a deleted private helper directly; then it points at the canonical one.
- Lock-in limit: definitions of the six helper names: 6 after this REQ (today 14).

## Dependencies
Depends on REQ-550 and REQ-552 so `corehelpers` is settled before helpers move in. REQ-558 depends on this REQ.

## Builder Guidance
Firm on one definition per name and on recording the three reconciliations; latitude on the exported names.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints fourteen definitions of six names.
**GREEN when:** It prints six lines, one per name; `go test ./...` green; the lock-in pins definitions of these six names at 6.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for per-req-duplicate-go-helpers.*
