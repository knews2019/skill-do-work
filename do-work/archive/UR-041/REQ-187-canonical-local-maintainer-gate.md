---
id: REQ-187
title: No single local maintainer command proves shell plus both Go modules
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T11:37:24Z
completed_at: 2026-08-15T12:11:43Z
commit: c20110d
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
write_set: [CLAUDE.md, justfile, _dev/tests/maintainer-verify.sh, _dev/tests/contract-regressions.sh]
route: C
kb_status: pending
kb_entry:
---

# No Single Local Maintainer Command Proves Shell Plus Both Go Modules

## What

Add one export-ignored local maintainer verification command as the source of truth for strict shell checks, the aggregate once, and vet/test in both Go modules; make documentation and any root recipe delegate to it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Create one export-ignored maintainer script with exact Go/ShellCheck version checks, warning-level tracked-shell lint, the aggregate once, both modules' vet/tests, and the REQ-185 strict JavaScript lane when Node is present. Give that script a recursion-safe `--self-test` mode with PATH shims and per-family failure replays; invoke only that mode from the aggregate, and make CLAUDE/Just thin delegates.
- [x] **[APPLY]:** Added the self-test first and captured its missing-dispatch RED, then implemented the exact-version, tracked-shell, aggregate, two-module, and conditional strict-JavaScript gate plus the CLAUDE/Just delegates within the four-file Scope.
- [x] **[UNIFY]:** Reviewed all four files, tightened the shims against wrong strict regex, missing warning severity, altered tracked-shell discovery, and unexpected Go arguments, then verified the complete canonical gate, direct self-test, Just parsing/delegation, Bash syntax, direct ShellCheck, export-ignore, and `git diff --check`; no debug artifacts remain.

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

## Triage

**Route C** — the required command families are firm, but the export-ignored script location, tool-version policy, recursion-safe contract seam, optional Node behavior, and repository-only Just placement needed an explicit plan before implementation.

## Plan

1. Add `_dev/tests/maintainer-verify.sh` as the sole production command inventory. Resolve the repository from `BASH_SOURCE`, fail on mismatched pinned Go/ShellCheck versions, lint every tracked shell file at warning severity, run the aggregate once, then vet/test queue-kanban and audit-metrics.
2. When Node is available, invoke REQ-185's exact strict JavaScript maintainer lane after the ordinary board tests; when it is absent, print an explicit optional-skip note and continue so this gate remains independently useful.
3. Implement `--self-test` in the same script. Use a private fixture and argument/cwd-aware PATH shims to assert every success stage exactly once, prove deliberate failures in version, ShellCheck, aggregate, board, strict-JavaScript, and audit-metrics families propagate nonzero, and prove the no-Node path succeeds. The self-test must never invoke the real aggregate.
4. Invoke only `maintainer-verify.sh --self-test` from `contract-regressions.sh`, so a normal maintainer run reaches the focused contract through its one aggregate call without recursion.
5. Point `CLAUDE.md` to the canonical script only and add a thin `maintainer-verify` Just delegate after the managed consumer marker.

## Exploration

- Both Go modules declare language version 1.26 and pass current vet/tests; the audited/current toolchain is Go 1.26.1. ShellCheck 0.11.0 passes all tracked shell files at warning severity.
- The aggregate already owns shipped-reference once and reaches prescribed-shell through staged-skills after REQ-186. The canonical gate must invoke only the aggregate, never those children separately.
- `/_dev` is export-ignored. A script under `_dev/tests/` and the maintainer-only CLAUDE pointer cannot affect installed consumers.
- The root `justfile` managed span ends at `# <<< do-work:recipes <<<`; a delegate after that marker is repository-owned and does not change the consumer template or bare-`just` default.
- A same-script `--self-test` mode avoids a second long-lived contract file. The real aggregate may invoke that mode safely because its fixture's shimmed child `bash` records the aggregate stage instead of executing it.
- REQ-185's strict lane is distinct from ordinary board `go test ./...`: the parent runs only when selected by its exact regex. Running it conditionally on Node adds behavior evidence without making this REQ depend on Node's presence.

## Scope

**Files I will touch:**
- `CLAUDE.md` (modify) — replace the hand-back command list with one canonical maintainer-script pointer
- `justfile` (modify) — add a thin repository-only delegate outside the managed consumer markers
- `_dev/tests/maintainer-verify.sh` (new) — canonical command inventory plus recursion-safe focused self-test
- `_dev/tests/contract-regressions.sh` (modify) — invoke only the canonical script's focused `--self-test` contract

**Files I will NOT touch:** consumer templates/runtime, Go module source or manifests, hosted CI, tool configuration files, or unrelated existing contract assertions.

**Acceptance criteria (restated from REQ):**
- [ ] One export-ignored local script owns tool versions, warning-level ShellCheck, one aggregate run, and vet/test for both Go modules.
- [ ] A deliberate failure in every child family makes the canonical command nonzero.
- [ ] CLAUDE and Just delegate without repeating the child command list; the Just recipe remains outside managed consumer markers.
- [ ] REQ-185's strict JavaScript lane runs when Node is available, while Node absence remains an explicit successful skip.
- [ ] Consumer optional-tool degradation is unchanged and no hosted CI or generic test registry is added.

## Discovered Tasks

- Four pre-existing late `assert_contains` calls in `_dev/tests/contract-regressions.sh` still open `Justfile` instead of the tracked lowercase `justfile`. They are outside REQ-187's canonical-gate scope and must not be repaired inline; a focused case-sensitive addendum to REQ-180 should own them.

## Implementation Summary

**Files changed:**
- `CLAUDE.md` (modified) — points hand-back verification at the canonical script without restating child commands
- `justfile` (modified) — adds a thin repository-only `maintainer-verify` delegate after the managed consumer marker
- `_dev/tests/maintainer-verify.sh` (new) — adds an executable gate that owns exact Go/ShellCheck versions, tracked warning-level shell lint, one aggregate run, both modules' vet/tests, conditional strict JavaScript behavior, and its focused recursion-safe self-test
- `_dev/tests/contract-regressions.sh` (modified) — invokes only `maintainer-verify.sh --self-test` so the contract stays automatic without recursively running the aggregate

**Behavior:** One repository-local command now proves the required shell and native families. Missing/mismatched required tools or any failing child aborts the gate; Node-backed JavaScript behavior runs strictly when Node exists and produces an explicit successful skip when it does not. Installed consumer behavior is unchanged because the gate lives under export-ignored `_dev`.

## Testing

**RED:** With the self-test contract in place but normal dispatch deliberately absent, `bash _dev/tests/maintainer-verify.sh --self-test` exited 1 because the all-success fixture could not observe the required stages.

**GREEN:**
- `bash _dev/tests/maintainer-verify.sh --self-test` — PASS; exact-once success, version mismatch, ShellCheck lint, aggregate, board vet/ordinary/strict, audit vet/test, wrong strict regex, and no-Node fixtures all behave as required
- `bash _dev/tests/maintainer-verify.sh` — PASS; aggregate final marker, queue-kanban vet/tests/strict JavaScript, and audit-metrics vet/tests all reached
- `just --dry-run maintainer-verify` and `just --list` — PASS
- `bash -n _dev/tests/maintainer-verify.sh _dev/tests/contract-regressions.sh` — PASS
- `shellcheck --severity=warning _dev/tests/maintainer-verify.sh` — PASS
- `git check-attr export-ignore -- _dev` — PASS (`set`)
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports the four-file Implementation Summary exactly matches Scope; the pre-existing capital-J assertions and foreign queue edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found the new executable and three modified files, complete P-A-U evidence, and no debug artifacts.
- **Substance and traceability:** PASS — the script is the only production command inventory, pins the audited Go/ShellCheck versions, and the self-test rejects removal, duplication, argument drift, or successful failure in every required family.
- **Wiring/data flow:** PASS — CLAUDE and Just each delegate once; normal mode invokes the aggregate once, whose only reverse edge is `--self-test`; that mode uses fixture shims and cannot invoke the real aggregate, so recursion terminates.

## Review

**Result:** Approve — Acceptance: Pass
**Overall score:** 98%

- **Requirements (100%):** One export-ignored command owns every required shell/native family, exact tool pins, failure propagation, delegates, and consumer boundary.
- **Code quality (96%):** The Bash 3.2-compatible implementation is cwd-independent, fail-fast, NUL-safe for tracked shell paths, and trap-cleans private fixtures.
- **Test adequacy (98%):** Exact-once, wrong-version, per-family failure, argument-shape, wrong-regex, Node-present/absent, recursion, and full live execution are covered.
- **Scope (100%):** Exactly the four declared files changed.

**Important findings:** None.
**Minor findings:** An invocation whose first argument is empty but which has later arguments enters normal verification instead of returning usage status 2. The documented zero-argument and `--self-test` entrypoints are unaffected.
**Explicit remediation:** None required for acceptance; a future touch may tighten the no-argument branch to require `$# -eq 0` and add the invalid invocation to self-test.

## Lessons Learned

- A canonical verification wrapper prevents command-list drift only when documentation and recipes are thin delegates and its behavioral contract counts/fails semantic stages rather than copying a second production command inventory.
- When a canonical gate must invoke an aggregate that also tests the gate, give the gate a fixture-only self-test mode that cannot reach the real aggregate. That closes the recursive ownership loop without duplicating a normal execution edge.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Orientation

**[MAP CHANGED]** Repository hand-back health now has one local entrypoint: `_dev/tests/maintainer-verify.sh`. It owns strict shell verification, the aggregate, both Go modules, and available Node-backed board behavior; CLAUDE and Just only delegate to it. The prescribed-shell prime remains current and gains this ownership/recursion lesson after archive.
