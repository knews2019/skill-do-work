---
id: REQ-532
title: 'Delete dead test scripts and nested self-tests from the maintainer gate'
status: pending
created_at: 2026-09-03T11:42:36Z
user_request: UR-102
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: [REQ-533]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
    - trivial short-circuit
  calculated_at: 2026-09-03T12:02:36Z
related: [REQ-531, REQ-519, REQ-520, REQ-521]
batch: close-the-tap
write_set:
  - _dev/tests/contract-regressions.sh
  - _dev/tests/maintainer-verify.sh
  - _dev/tests/fixtures/shipped-shell-command-map.tsv
  - _dev/tests/do-work-cli-go125-compatibility.sh
  - _dev/tests/flat-just-recipes-behavior.sh
  - _dev/tests/memory-hook-behavior.sh
  - _dev/tests/record-commit-hash-guards.sh
  - _dev/tests/shipped-shell-parity.sh
  - _dev/tests/shipped-shell-thinness.sh
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
gate_deferred: 'true'
---

# Delete Dead Test Scripts and Nested Self-Tests From the Maintainer Gate

## What

Remove from `_dev/tests/` every script that no gate path executes, and take the harness self-tests out of the per-run aggregate suite so they run only when someone edits the harness. This is the deletion pass that precedes UR-100's optimization of what remains.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The gate takes 7.5 minutes and every second of it is lint, test, or self-test over 19,177 lines of shell tests plus 1,009 Go tests. Dead scripts still cost ShellCheck time on every run (the gate lints all tracked `*.sh`) and one of them broke the gate on 2026-09-01 (commit 2d140f63 fixed SC2034 in a script nothing runs). Delete before you add.

## Context

Measured 2026-09-03 against `_dev/tests/contract-regressions.sh` (8,468 lines) and `_dev/tests/maintainer-verify.sh` (556 lines):

- Scripts no gate path executes (the condition is the rule; these are today's instances): `do-work-cli-go125-compatibility.sh`, `flat-just-recipes-behavior.sh`, `memory-hook-behavior.sh`, `record-commit-hash-guards.sh`, `shipped-shell-parity.sh`, and `shipped-shell-thinness.sh` (reachable only from `shipped-shell-parity.sh`). Their remaining references are `_dev/tests/fixtures/shipped-shell-command-map.tsv`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `decisions/audits/2026-08-11-defensive-surface.md`, and archived REQ prose.
- Harness self-tests nested in the aggregate: `contract-regressions.sh:7319` runs `maintainer-verify.sh --self-test`, which re-executes the gate script 15 times through shims; `contract-regressions.sh:7356` runs the action shell-block probe's `--self-test`.

Decisions taken at capture:

- D1 A self-test stays runnable standalone (`bash _dev/tests/maintainer-verify.sh --self-test`) and is documented as the check to run after editing that harness. It leaves the per-run aggregate.
- D2 Deleting a script also deletes or updates every reference that would otherwise dangle: the fixture map row, the prime file sentence, the audit record gets a one-line note. Archived REQ and changelog prose is history and is left alone.
- D3 The duplicate strict JavaScript lane is REQ-520's; the line-count ratchet is REQ-519's; parallel sub-suites are REQ-521's. This REQ does not touch them.

## Detailed Requirements

- After the change, `git ls-files _dev/tests/*.sh` lists no script that neither `maintainer-verify.sh`, `contract-regressions.sh`, nor a script they invoke executes. Prove it with a one-off listing in the REQ's Testing section, not a new permanent check.
- `grep -n -- '--self-test' _dev/tests/contract-regressions.sh` returns nothing.
- `_dev/tests/contract-regressions.sh` shrinks; it never grows.
- `bash _dev/tests/maintainer-verify.sh` exits 0 and `bash _dev/tests/maintainer-verify.sh --self-test` exits 0.
- Record the before and after `time bash _dev/tests/maintainer-verify.sh` wall clock in the REQ's Testing section.

## Constraints

- No new prose that walks a shell sequence; mechanics stay in the scripts.
- Nothing under `skills/` changes except the one prime-file sentence that names a deleted script.

## Red-Green Proof
**RED prompt/case:** `grep -c -- '--self-test' _dev/tests/contract-regressions.sh` and `ls _dev/tests/shipped-shell-parity.sh`.
**Why RED now:** The count is 2 and the file exists; the gate lints and, through the aggregate, re-runs itself 15 times on every run.
**GREEN when:** The count is 0, the six named scripts are absent from `git ls-files`, no remaining tracked file outside `do-work/archive/`, `kb/`, and the changelogs names them, and both gate invocations exit 0.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes shipped-shell parity fixtures and prescribed command blocks.

## Full Context
See `do-work/user-requests/UR-102/input.md` for complete verbatim input.

---
*Source: D2 of the 2026-09-03 roadmap disposition, selected by the maintainer.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the exact scripts and references to remove, constrains the only retained self-test behavior, and defines concrete verification commands.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Repository Gate Deferral

- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** contract-regressions:claude-write-surface-count-stale
- **Repair dependency:** REQ-533
- **Diagnostic evidence:** "CLAUDE.md must state the tool has exactly three write surfaces once next-req reserves ids; testing fields, next-version, and reservation markers are the complete set."
