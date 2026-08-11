---
id: REQ-165
title: Shell-block lint harness for shipped action files
status: completed
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
write_set: [_dev/tests/action-shell-blocks.sh, _dev/tests/contract-regressions.sh, skills/do-work-board/actions/board.md]
claimed_at: 2026-08-11T11:59:35Z
route: C
completed_at: 2026-08-11T12:20:15Z
kb_status: pending
kb_entry:
---

# Shell-Block Lint Harness for Shipped Action Files

## What

Add a `_dev/tests/` check that extracts every fenced shell block (` ```bash ` / ` ```sh `) from the shipped `skills/` tree — action files, crew files, and shipped hook scripts — and lints them: `bash -n` (syntax) always, `shellcheck` when available. This makes the suite's largest defect generator (prose-prescribed shell that nothing executes until an agent runs it in a consumer repo) testable in CI.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Build one Bash probe that preserves Markdown source locations, narrowly neutralizes angle-bracket placeholders, lints every fence and shipped shell file, and proves its own failure path; wire it through the existing explicit child-probe seam.
- [x] **[APPLY]:** Added the scoped standalone lint/self-test probe and explicit aggregate-runner invocation; no shipped skill or runtime file was changed.
- [x] **[UNIFY]:** Reviewed `git diff --stat` and the complete contents/diff of `_dev/tests/action-shell-blocks.sh` and `_dev/tests/contract-regressions.sh`; ran `bash -n`, ShellCheck at warning severity, `git diff --check`, debug-artifact grep, targeted positive/negative/degraded checks, and the aggregate contract suite.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** This is a new repository-wide test harness with extraction, placeholder normalization, optional ShellCheck integration, negative self-proof, suite wiring, and likely pre-existing findings to disposition.

**Planning:** Required

## Plan

1. Add `_dev/tests/action-shell-blocks.sh` as a standalone behavioral probe. It will enumerate shipped Markdown and hook scripts under `skills/`, extract every `bash`/`sh` fence while preserving source path and starting line, neutralize only angle-bracket prose placeholders without changing line count, run `bash -n` for every block and shipped `.sh` hook, and run ShellCheck when available with snippet-context-only exclusions.
2. Give the probe an explicit `--self-test` mode that feeds a deliberately malformed fenced block through the same checker and passes only when the checker fails with the fixture path, source line, and diagnostic. Its normal mode will aggregate findings, name each source location, note when ShellCheck is unavailable, and return non-zero on real findings.
3. Modify `_dev/tests/contract-regressions.sh` to invoke both the negative self-test and the clean-tree scan beside the other standalone behavioral probes, with a missing-probe guard and actionable failure summary.
4. Verify requirement coverage by running the probe self-test, the normal scan with ShellCheck present, the normal scan with ShellCheck hidden from `PATH`, `bash -n` on both edited scripts, ShellCheck on both edited scripts, and the full contract regression suite.

**Architectural decisions:** Keep the checker in Bash/POSIX utilities to match `_dev/tests/` conventions; use temporary per-block files plus a manifest so diagnostics can be remapped to original Markdown lines; treat SC2034/SC2154 as block-isolation noise while retaining warning-or-higher ShellCheck findings; lint shipped hook `.sh` files directly because they are part of the named shipped-shell surface.

**Requirement mapping:** Extraction and attribution are covered by task 1; Bash/ShellCheck linting by tasks 1 and 4; fail-capable proof by task 2; suite wiring by task 3; filename/line/diagnostic output by tasks 1 and 2; placeholder handling by task 1.

*Generated by Plan phase after repository inspection*

**Plan validation:** All five Detailed Requirements map to planned tasks, and both planned files trace directly to the harness and its required suite seam. Scope is four ordered tasks because verification is separated from the two implementation files; no product or shipped-skill edits are planned unless the new checker exposes a genuine blocking defect.

## Exploration

- `_dev/tests/contract-regressions.sh` is the repository's aggregate runner. Its standalone behavioral probes are explicitly guarded and invoked near the existing suite-manifest, shipped-reference, commit-hash, updater, staging, and installer probes because no glob auto-discovers `_dev/tests/*.sh`.
- Existing probes use Bash, `fail_count`, `mktemp -d`, cleanup traps, and `FAIL:` diagnostics; the aggregate runner converts a child failure into one additional named summary without suppressing the child's details.
- The shipped `skills/` tree currently contains 49 fenced `bash`/`sh` blocks plus three shipped hook scripts. Replacing angle-bracket prose placeholders such as `<skill-root>` or `<pre>` with a same-line dummy token is sufficient for `bash -n` to parse all current fences.
- Per-block ShellCheck otherwise reports isolation-only SC2034 (assigned in this excerpt, consumed later) and SC2154 (assigned in an earlier prescribed block). Excluding only those two codes and setting severity to `warning` retains actionable warnings/errors while avoiding false failures from fenced-excerpt boundaries; shipped hook scripts are clean under the same threshold.
- Source attribution should be owned by the harness rather than inferred from temporary filenames: each extracted block needs its original relative path and first body line, with tool line numbers offset back to the Markdown source.

*Generated by Exploration phase*

## Scope

**Files I will touch:**
- `_dev/tests/action-shell-blocks.sh` (new) — extract, normalize, lint, self-test, and report shipped shell blocks and hook scripts
- `_dev/tests/contract-regressions.sh` (modified) — invoke the new behavioral probe and its negative self-test from the aggregate suite
- `skills/do-work-board/actions/board.md` (modified) — fail immediately when the prescribed repository-root `cd` cannot succeed, closing the one real warning exposed by the complete scan

**Files I will NOT touch:** other shipped `skills/` content, unrelated test probes, runtime hooks, product code, or release metadata during implementation

**Acceptance criteria (restated from REQ):**
- [ ] Every fenced `bash`/`sh` block across shipped `skills/` content is extracted with source file and line attribution; shipped hook `.sh` files are linted directly.
- [ ] Every extracted block runs through `bash -n`; ShellCheck runs when available and absence is a non-failing noted degradation.
- [ ] Prose placeholders are neutralized narrowly rather than causing whole-block skips.
- [ ] A deliberately invalid fixture proves the checker returns failure and reports file, line, and diagnostic.
- [ ] The aggregate `_dev/tests/contract-regressions.sh` runner invokes the new probe beside existing behavioral probes.
- [ ] Clean-tree findings name the source file, source line, and diagnostic.

## Decisions

- **D-01 — DECIDE & STATE:** ShellCheck excludes SC2034 and SC2154 for extracted fences only. Those diagnostics require context outside a standalone Markdown block (a variable may be assigned or consumed in a different prescribed block), while syntax and all other warning-or-higher diagnostics remain enforced. Complete shipped `.sh` sources receive the full warning set with no snippet-context exclusions.
- **D-02 — DECIDE & STATE:** The direct-source scan covers every shipped `skills/**/*.sh` file, not only `hooks/*.sh`. This is a small superset of the named hook surface and prevents another shipped shell entry point from sitting outside the same behavioral gate.
- **D-03 — DECIDE & STATE (remediation):** Treat zero-to-three leading blanks as part of a valid Markdown fence boundary, and make the deliberately broken self-test use that indented form. Review independently counted 59 valid shipped fences versus 49 recognized by the first implementation; changing both discovery and the fixture closes the blind spot rather than patching the current ten instances.
- **D-04 — DECIDE & STATE (scope extension):** Extend the declared scope to `skills/do-work-board/actions/board.md` for the one SC2164 warning revealed once indented fences were included. Adding `|| exit 1` to the prescribed root-directory change is the direct fail-fast correction and keeps the harness strict; suppressing the warning would preserve a real execution defect merely to keep the original two-file boundary.

## Implementation Summary

**Files changed:**
- `_dev/tests/action-shell-blocks.sh` (new)
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work-board/actions/board.md` (modified)

**What was done:** Added an attributable Bash/ShellCheck harness for all shipped Markdown shell fences and shipped shell source files, including valid fence indentation, narrow placeholder normalization, a deliberately broken indented self-test fixture, an explicit no-ShellCheck degradation path, and aggregate contract-suite invocation. Fixed the single real warning exposed by the complete scan by making the board action's prescribed repository-root directory change fail fast.

## Qualification

Passed after remediation — 3 implementation files verified, 5 Detailed Requirements traced, P-A-U confirmed. The new probe is substantive, the runner wiring invokes the real public test seam, source discovery is dynamic rather than hardcoded, and targeted checks covered malformed indented input, all 59 current fences, complete shipped shell sources, and the optional-linter fallback.

## Testing

**Tests run:**
- `bash -n _dev/tests/action-shell-blocks.sh`
- `bash -n _dev/tests/contract-regressions.sh`
- `shellcheck --severity=warning _dev/tests/action-shell-blocks.sh _dev/tests/contract-regressions.sh`
- `_dev/tests/action-shell-blocks.sh --self-test`
- `_dev/tests/action-shell-blocks.sh`
- `PATH=/usr/bin:/bin _dev/tests/action-shell-blocks.sh`
- `bash _dev/tests/contract-regressions.sh`
- `git diff --check`

**Result:** ✓ All passing after remediation. Normal mode linted 59 fenced blocks and 15 shipped shell files with ShellCheck enabled; degraded mode linted the same surface with `bash -n` and a non-failing absence note. The full contract suite completed with exit 0.

**Red-green validation:**
- Captured missing-checker case: baseline contract suite passed while containing no shell-fence probe → the new `--self-test` now proves a deliberately malformed fenced block is rejected with its fixture path, mapped source line, and Bash diagnostic.
- Aggregate seam: standalone probe was not invoked before implementation → the full contract suite now runs both its negative self-test and clean shipped-tree scan.

**New tests added:**
- `_dev/tests/action-shell-blocks.sh` — behavioral lint probe plus fail-capable self-test mode

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` — explicitly invokes the REQ-165 behavioral probe because `_dev/tests/*.sh` is not auto-discovered

*Verified by work action*

## Remediation

The first review failed acceptance because an exact column-zero fence matcher scanned 49 of 59 valid shipped fences. The remediation widened discovery to Markdown-valid zero-to-three-space indentation, changed the deliberately broken fixture to the indented form, and fixed the one genuine SC2164 warning the complete scan exposed. Independent enumeration, targeted modes, scope drift, qualification, and the aggregate contract suite all pass after the fix.

## Review

**Overall: 99%** | 2026-08-11T12:19:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — all 59 shipped fences and 15 shipped shell files are linted; the indented malformed fixture fails attributably, optional ShellCheck degradation is clean, and the aggregate suite passes.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action after one successful remediation attempt*

## Lessons Learned

**What worked:**
- Pairing the checker with a deliberately broken fixture proved the failure path, while an independent count of Markdown-valid fences tested the discovery boundary the checker could not validate about itself.
- Keeping the new probe behind `_dev/tests/contract-regressions.sh` matched the repository's explicit child-probe convention and made the ratchet part of the normal suite.

**What didn't:**
- The first extractor regex matched only column-zero fences. Its own clean run looked convincing while silently skipping 10 indented fences; comparing 49 reported blocks with an independent 59-block enumeration exposed the blind spot.

**Worth knowing:**
- SC2034 and SC2154 are excluded only for isolated Markdown fences because assignments and uses may live in separate prescribed blocks. Complete shipped `.sh` files receive the full warning set.
- A newly complete scan can reveal genuine pre-existing findings. Fix a narrow blocking defect, as with the board action's unguarded `cd`, instead of weakening the checker to preserve a green baseline.

## Orientation

[MAP CHANGED] The repository's contract-test subsystem now treats shipped shell guidance as executable input: the aggregate runner invokes a fail-capable probe that discovers every Markdown `bash`/`sh` fence plus every shipped shell source, remaps diagnostics to source locations, and degrades cleanly when ShellCheck is absent.
