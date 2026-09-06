---
id: REQ-603
status: completed
domain: general
created_at: 2026-09-06T08:19:05Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-shell-commands.md]
tdd: true
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-06T16:05:00Z
  basis:
    - Route B
    - 6-file write set
    - 2 subsystems (launcher shell + corehelpers Go CLI)
depends_on: [REQ-597, REQ-601]
related: [REQ-597, REQ-601]
write_set: [skills/do-work/scripts/protected-inventory.sh, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Let the protected-inventory launcher pass global flags, and stop its shim discarding the text it prepared'
claimed_at: 2026-09-06T13:04:17Z
completed_at: 2026-09-06T13:14:46Z
commit: 4afd5e9f768426e1a741bca33e245dda7b806dae
release_at: 2026-09-06T13:14:46Z
---

# Let the Protected-Inventory Launcher Pass Global Flags, and Stop Its Shim Discarding the Text It Prepared

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Implement global flag pass-through (`--repo-root`, `--format`) in `skills/do-work/scripts/protected-inventory.sh` ahead of the subcommand token. Update `internal/corehelpers/inventory.go` compatibility shim to preserve prepared diagnostic strings and findings on failure. Add unit tests in `inventory_test.go` covering swallowed outcomes (`NO-DO-WORK-DIR`, `PARSE-FAILED`, walk error `HELPER-USAGE`). Update `commit.md`, `inspect.md`, and `prescribed-shell-primitives.md` to reflect new `--repo-root` behavior and correct exit 2 diagnostic vs skip handling and quarantine lifecycle.)
- [x] **[APPLY]:** (Agent: Implemented flag forwarding in `protected-inventory.sh`, shim fix in `inventory.go`, tests in `inventory_test.go`, and documentation in `commit.md`, `inspect.md`, and `prescribed-shell-primitives.md`.)
- [x] **[UNIFY]:** (Agent: Ran `git diff --stat`, `gofmt -w`, native project tests, and all maintainer guard checks (`audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, `quiet-grep-pipeline-audit.sh`, `gate.sh`).)

## Triage

**Route: B** — Substantive change with TDD and multiple caller updates.

**Reasoning:**
- Modifies both shell launcher wrapper (`scripts/protected-inventory.sh`) and Go implementation (`internal/corehelpers/inventory.go`).
- TDD required: replicated red test before fixing compatibility shim error discarding.
- Modifies documentation in core and toolbox actions (`commit.md`, `inspect.md`) and the prescribed shell guide.

## Plan

1. **Launcher global flag forwarding:** Sift arguments in `skills/do-work/scripts/protected-inventory.sh` to extract global flags (`--repo-root`, `--format`) and place them before the `protected-inventory` command token, forwarding remaining arguments after.
2. **Shim error preservation (TDD):**
   - Write red unit test in `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` verifying that `DO_WORK_COMPATIBILITY_SHIM=1` does not discard `NO-DO-WORK-DIR`, `PARSE-FAILED`, or walk error findings.
   - Modify `handleProtectedInventory` in `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` so `ExactTextOutput` is only overwritten when `result.Outcome == resultmodel.OutcomeSuccess`.
   - Run tests to verify green.
3. **Action and documentation updates:**
   - Update `skills/do-work/actions/commit.md`: note `--repo-root` pass-through from subdirectories, exit 2 coverage for status/quarantine failures, re-run using `start` with quarantine replacement vs union, and exit 2 error reporting vs skip.
   - Update `skills/do-work-toolbox/actions/inspect.md`: note `--repo-root` pass-through, `associate` exit 2 after `start --dry-run`, and exit 2 error reporting.
   - Update `skills/do-work/docs/prescribed-shell-primitives.md`: document `--repo-root` flag forwarding in launcher sentence.
4. **Verification:**
   - Run unit tests, repository guards, and full maintainer gate.

## Exploration

Explored `skills/do-work/scripts/protected-inventory.sh`, `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`, and callers.
- Found that `commandruntime` expects global flags before the command verb; passing them after resulted in `unknown option` errors.
- Found that `handleProtectedInventory` unconditionally rebuilt `ExactTextOutput` when `DO_WORK_COMPATIBILITY_SHIM=1`, which overwrote any error text prepared by `handleAssociate` on non-success outcomes.
- Confirmed that `start` replaces the quarantine file whereas `associate` unions with it.
- Confirmed that running `start --dry-run` writes no quarantine file, leading to exit 2 on subsequent `associate`.

## Scope

Files in scope for this change:
- `skills/do-work/scripts/protected-inventory.sh`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go`
- `skills/do-work/actions/commit.md`
- `skills/do-work-toolbox/actions/inspect.md`
- `skills/do-work/docs/prescribed-shell-primitives.md`

## What

Two defects in one wrapper, both measured by REQ-597's builders while correcting the prose that
describes it (evidence and fixtures in `do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`).

**The launcher cannot pass a global flag.** `skills/do-work/scripts/protected-inventory.sh:6` is
`exec bash .../do-work-cli.sh --format text protected-inventory "$@"`: everything a caller passes lands
after the command token, so `--repo-root` is rejected as `unknown option` and the runtime takes the
current directory as the root (`commandruntime/command_runtime.go:103-107`). Run from a subdirectory,
`start` prints the same rows but `associate` exits 2 with nothing on either stream. Neither
`commit.md` nor `inspect.md` can run the wrapper from anywhere but the project root, and both now say
so because REQ-597 had no better answer in prose.

**The compatibility shim discards the text it prepared.** `internal/corehelpers/inventory.go:445-456`
replaces the result text unconditionally, so the `NO-DO-WORK-DIR` line prepared at `:205`, the
`PARSE-FAILED` line at `:212` and a walk error's `HELPER-USAGE` finding at `:215` never reach a
caller. Measured silent exit 2 for: a missing `do-work/` directory, an unmatched backtick in an
Implementation Summary line, and an unreadable REQ file. Both callers read every exit 2 as "skip REQ
tracing", which is the reading that let `inspect.md`'s broken blocks ship unnoticed.

## Why

A silent exit 2 that means "skip" is indistinguishable from a silent exit 2 that means "the tool could
not run", and the callers have been treating the second as the first. The wrapper prepares the text that
would tell them apart and then throws it away.

## Detailed Requirements

- The launcher passes global flags through. `tools/checks/associate-files.sh:10-17` shows the
  translating shape; after the change `protected-inventory.sh --repo-root <root> start` (or the flag
  after the mode, whichever the CLI's own convention is; say which) works from any directory. State the
  chosen convention in the guide's launcher sentence.
- The shim keeps what the handlers prepared: `NO-DO-WORK-DIR`, `PARSE-FAILED` and a walk error's
  finding reach stdout or stderr as the handler intended. One test per swallowed outcome, red before
  the change, asserting the text is present.
- Then the callers: `commit.md` and `inspect.md` stop reading a printed finding as the skip condition.
  Exit 2 with a finding is an error to report; a genuine "nothing to associate against" is whatever the
  handler prints for it. Rewrite the exit-2 sentences in both files from the new behaviour, measured.
- Three neighbouring prose defects in the same sentences, found on the same pass: `commit.md:67` tells
  a re-run to append to the retained quarantine, but `start` replaces it (`inventory.go:393`) and only
  `associate` unions (`:420-425`), so say which command a re-run uses; `commit.md:61`'s "exit 2 means
  this is not a git repo" also covers a `git status` failure (`:368`) and a quarantine write failure
  (`:394`); `associate` after `start --dry-run` exits 2 as not-started (`:412-415`, `:392-396`),
  which neither action states.
- The guide's sentences REQ-597 wrote about the current-directory dependence are updated to the new
  behaviour in the same commit; a sentence describing the old launcher must not outlive it.

## Constraints

- Every prose sentence is written from a measured run, not from the code alone (REQ-596's rule for this
  file class). Keep the fixtures REQ-597 left under the scratchpad if they still exist; rebuild them if not.
- Shipped Go and shipped shell: a release.

## Open Questions

None.

## Implementation Summary

- `skills/do-work/scripts/protected-inventory.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work-toolbox/actions/inspect.md` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**Passed global flags through launcher, preserved prepared failure text in compatibility shim, and corrected caller action guidance:**
- `protected-inventory.sh`: Parsed and forwarded `--repo-root` and `--format` ahead of the `protected-inventory` command token so callers can pass `--repo-root <path>` either before or after the mode (`start` or `associate`).
- `inventory.go`: Fixed compatibility shim in `handleProtectedInventory` to overwrite `ExactTextOutput` only when `result.Outcome == resultmodel.OutcomeSuccess`, preserving prepared error strings (`NO-DO-WORK-DIR: nothing to associate against`, `PARSE-FAILED: ...`) and findings on failure.
- `inventory_test.go`: Added `TestProtectedInventoryCompatibilityShimPreservesErrors` covering all three swallowed error conditions (`NO-DO-WORK-DIR`, `PARSE-FAILED`, and directory walk error `HELPER-USAGE`).
- `commit.md`: Documented that `--repo-root` can be passed when invoked from a subdirectory, clarified that exit 2 covers git status or quarantine write failures as well as non-git repositories, clarified that re-running inventory uses `start` (which replaces quarantine rather than appending, whereas `associate` unions), and clarified that exit 2 with a finding or error is an error to report rather than a skip condition.
- `inspect.md`: Documented `--repo-root` pass-through from subdirectories, noted that `associate` exits 2 with `HELPER-USAGE` if `start` was run with `--dry-run`, and clarified that exit 2 with a finding or error is an error to report.
- `prescribed-shell-primitives.md`: Documented launcher global flag forwarding and `--repo-root` pass-through in the executable homes table.

## Decisions

- **D1 Flag sifting in launcher:** Sifting `--repo-root` and `--format` into `global_arguments` allows callers to place `--repo-root` before or after the mode token without violating `do-work-cli`'s requirement that global flags precede the command verb.
- **D2 Condition for shim text formatting:** Only format association tab-separated output when `result.Outcome == resultmodel.OutcomeSuccess`, ensuring non-zero outcomes retain their diagnostic messages.

## Qualification

**Passed.** Read from the range `5a7b213608433a944d4ef34b6f6cf3f1d2ad4b5e..4afd5e9f768426e1a741bca33e245dda7b806dae`, 6 files, 130 insertions, 18 deletions.
Canonical `qualify` and `scope-drift` both satisfied.

- `skills/do-work/tools/do-work-cli` unit tests passed (including new tests in `inventory_test.go`).
- Repository guards `audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, and `quiet-grep-pipeline-audit.sh` all passed with exit 0.
- Full maintainer gate `gate.sh` passed with exit 0 (126s wall time, 799 tests).

## Testing

**Commands executed:**
- `go test -count=1 ./internal/corehelpers/...` — passed, exit 0.
- `go test -count=1 ./...` — all packages passed, exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` — passed, exit 0.
- `bash _dev/tests/audit-lockins.sh` — passed, exit 0.
- `bash _dev/tests/action-shell-blocks.sh` — passed, exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — passed, exit 0.
- `DO_WORK_GATE_ROOT="$(pwd)" bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh` — `Maintainer verification passed.`, exit 0.

## Review

**Overall: 98%** | 2026-09-06T16:13:00Z | Synthesis of review lenses (code correctness, test adequacy, documentation fidelity, gate verification)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Pass.** Launcher now transparently forwards global flags, diagnostic text is preserved under the compatibility shim, all swallowed error conditions are pinned with unit tests, and caller action documentation accurately reflects execution mechanics and quarantine lifecycles.

## Remediation

None needed.

## Lessons Learned

**What worked:**
- Sifting options in the bash wrapper cleanly bridges positional invocation patterns with strict global-flag CLI structures.
- Unit testing each swallowed failure mode with `t.Setenv("DO_WORK_COMPATIBILITY_SHIM", "1")` directly locked in diagnostic visibility.

**What didn't:**
- An unconditional string builder loop in a compatibility shim can silently discard errors prepared by underlying handlers if outcome status isn't checked first.

**Worth knowing:**
- `start` replaces the Git-private quarantine file, while `associate` unions with it. Re-runs must account for this difference when preserving previous exclusions.

## Orientation

Fixes `scripts/protected-inventory.sh` to forward `--repo-root` and `--format` ahead of the command token, and fixes `inventory.go` compatibility shim so prepared failure diagnostics (`NO-DO-WORK-DIR`, `PARSE-FAILED`, and walk error findings) reach callers instead of being replaced with empty output. Updates `commit.md`, `inspect.md`, and `prescribed-shell-primitives.md` accordingly.
