---
id: REQ-187
title: No single local maintainer command proves shell plus both Go modules
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: refactor
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-183, REQ-184, REQ-185, REQ-186, REQ-188]
batch: audit-findings-2026-08-14
write_set: [CLAUDE.md, justfile]
---

# No Single Local Maintainer Command Proves Shell Plus Both Go Modules

## What

Add one export-ignored local maintainer verification command as the source of truth for strict shell checks, the aggregate once, and vet/test in both Go modules; make documentation and any root recipe delegate to it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The required hand-back commands run shell contracts but neither Go module's vet/tests. The root `justfile` has no maintainer verification entry, so establishing native repository health requires manually coordinating five command families.

## Context

- Audit priority: P3; impact 2; effort normal.
- Root-cause key: `canonical-local-maintainer-gate`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 7.
- Reproduce: `rg -n 'contract-regressions|shipped-package-reference' CLAUDE.md && rg -n 'go (test|vet)' _dev/tests/contract-regressions.sh _dev/tests/shipped-package-reference-contract.sh || true && find skills -name go.mod -print && sed -n '1,40p' justfile`.

## Detailed Requirements

- Add one export-ignored maintainer verification script as the command source of truth.
- Have it own exact tool-version checks, ShellCheck at warning severity, the aggregate contract suite once, and `go vet ./...` plus `go test ./...` in both Go modules.
- Point `CLAUDE.md`'s hand-back instruction to the canonical script rather than maintaining a second command list.
- If a root Just recipe is added, keep it as a thin repository-only delegate outside the managed consumer markers.
- Add focused contract coverage proving that a deliberate failure in each child family propagates nonzero.
- Preserve consumer runtime's optional-tool degradation; strictness is for repository maintainers.

## Constraints

- Keep one command list; do not duplicate the command families in documentation or YAML.
- Do not add hosted CI in this REQ. The audit preserves that as a separate Discuss decision.
- Lock-in limit: zero required maintainer check families outside the canonical script.

## Dependencies

None. REQ-185's future strict JavaScript lane may be invoked when available, but this gate must remain useful independently. Coordinate `CLAUDE.md` with REQ-186.

## Builder Guidance

Firm intent with implementation latitude for the export-ignored script's exact repository-local path. The PLAN must update the capture-seeded `write_set` with that chosen path before building.

## Open Questions

None. Hosted Linux CI is explicitly out of scope pending a separate maintainer decision.

## Red-Green Proof
**RED prompt/case:** Follow the current hand-back instructions and inspect the root recipes; neither path proves vet/test in both Go modules, and there is no one local command for the complete set.
**Why RED now:** Native health requires manually coordinating strict ShellCheck, the aggregate, and four Go invocations.
**GREEN when:** One local script runs every required family exactly once, one documented command invokes it, any thin root recipe delegates to it, and a deliberate failure in each family makes the canonical command exit nonzero.
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 07, labeled P3, impact 2, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 7 in the canonical audit.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
