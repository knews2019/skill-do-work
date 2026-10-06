# REQ-634 hand-back (association walk skips an unparseable REQ file)

- Branch: worktree-agent-REQ-634-association-walk-skips-unparseable-req-file
- Base: 5b7c3f9a ([REQ-634] claim request lifecycle)
- Commit: e86b314d ([REQ-634] association walk skips an unparseable REQ file and keeps going)
- Not touched, per brief: VERSION, CHANGELOGs, version.md, lessons satellite, lessons index, checks.go, anything under do-work/. The orchestrator owns the release and the lessons entry.

## File manifest

All modified, none new.

- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`: new `UnparsedSummaryRecord` type; `AssociateProjectPaths` now returns `(map[string]string, []UnparsedSummaryRecord, error)`; walk callback records the parse error and returns nil; `handleAssociate` drops the PARSE-FAILED branch (which keyed on the message string "unmatched backtick") and emits one `ASSOCIATION-SUMMARY-UNPARSED` warning per skipped file. No os.Stderr write in the handler. `handleProtectedInventory` unchanged.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`: fix-hint case for `ASSOCIATION-SUMMARY-UNPARSED` (see D-02).
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go`: new `TestAssociationSkipsUnparseableSummaryAndKeepsWalking`; subtest `preserves PARSE-FAILED` flipped and renamed to `unparseable summary claims no paths and warns`; existing caller updated for the new signature.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`: caller updated for the new signature only.
- `_dev/tests/contracts/core-checks.sh`: `associate_unmatched` probe flipped (adds REQ-503 claiming legacy-file.txt, stderr captured to a separate file).
- `skills/do-work/actions/commit.md`, `skills/do-work-toolbox/actions/inspect.md`: exit-2 sentence no longer lists PARSE-FAILED; new sentence states the skip rule. HELPER-USAGE and its cases kept verbatim.
- `skills/do-work/docs/prescribed-shell-primitives.md`: one new bullet under "What `associate` settles". By-hand fallback paragraph left alone (it does not contradict the rule).

## Red-green evidence

Step 0 (no behaviour change): signature-only stub so the new tests compile; behaviour identical to HEAD, so the RED below is the dispatch-HEAD behaviour.

1. **TestAssociationSkipsUnparseableSummaryAndKeepsWalking**
   - RED: `inventory_test.go:415: one unparseable REQ file must not fail the walk: unmatched backtick in Implementation Summary`
   - GREEN: `--- PASS: TestAssociationSkipsUnparseableSummaryAndKeepsWalking (0.00s)`
2. **TestProtectedInventoryCompatibilityShimPreservesErrors/unparseable_summary_claims_no_paths_and_warns**
   - RED: `inventory_test.go:672: associate must succeed past an unparseable REQ file: ... Outcome:"failure" ... ExitCodeOverride:2`
   - GREEN: `--- PASS: TestProtectedInventoryCompatibilityShimPreservesErrors (0.34s)` (all three subtests pass; walk-error HELPER-USAGE subtest still green, so an unreadable REQ file still fails loudly)
3. **core-checks.sh associate_unmatched probe**
   - RED (exit 1): `FAIL: tools/checks/associate-files.sh must skip a REQ file whose Implementation Summary has an unmatched backtick: exit 0, the other REQ claim on stdout, the skipped REQ file named on stderr, no PARSE-FAILED; got exit 2, stdout: PARSE-FAILED: unmatched backtick in Implementation Summary, stderr: `
   - GREEN: exit 0, `core-checks contract probes passed.`

The runtime writes findings straight to the process stderr (command_runtime.go writeResult uses os.Stderr, no writer seam), so the Go subtest asserts the finding; the stderr line itself is asserted by the core-checks probe and the end-to-end transcript.

## Verification

| Check | Result |
|---|---|
| gofmt -l . | empty |
| go vet ./... | clean |
| go test -count=1 ./internal/corehelpers/ | pass, 7.3s package wall time (every file under the 30s budget) |
| go test -count=1 ./internal/commandruntime/ | pass, 0.8s |
| bash _dev/tests/contracts/core-checks.sh | exit 0, 4s wall |
| shellcheck -S warning core-checks.sh | 0 findings before and after |
| git diff --check | clean |

PARSE-FAILED grep (skills, _dev; go/sh/md), every remaining hit accounted for:
- `skills/do-work/CHANGELOG.md:390`: history line, fine.
- `inventory_test.go:678-679`: negative assertion (PARSE-FAILED must be absent).
- `core-checks.sh:403,405`: negative assertion and its FAIL message.

"unmatched backtick" sweep: `_dev/tests/prescribed-shell-cases/qualify.sh` and the core-checks scope-drift probe still pin the strict errors of qualify and scope-drift, as the REQ requires. No other shipped prose restates the old associate exit contract.

## P-A-U

- [x] **[PLAN]:** Signature stub first so tests compile, RED at three layers (unit, shim subtest, shell probe) plus end-to-end, then: callback appends `UnparsedSummaryRecord{repo-relative path, parser message}` and returns nil; handler deletes the message-string branch and appends warnings; fix hint added; three prose sites. Keyed on the condition "allBacktickedPaths returned an error", not on a message string (closed-enumeration-for-a-condition). Skip is announced through the runtime's stderr printer, never silent (silent-skip-reads-as-red); handler returns a typed result and never writes stderr (single-exit-owner).
- [x] **[APPLY]:** Done as planned, within the write boundary.
- [x] **[UNIFY]:**
```
 _dev/tests/contracts/core-checks.sh                | 30 +++++++++--
 skills/do-work-toolbox/actions/inspect.md          |  2 +-
 skills/do-work/actions/commit.md                   |  2 +-
 skills/do-work/docs/prescribed-shell-primitives.md |  1 +
 .../do-work-cli/internal/corehelpers/commands.go   |  2 +
 .../internal/corehelpers/commands_test.go          |  2 +-
 .../do-work-cli/internal/corehelpers/inventory.go  | 36 +++++++++----
 .../internal/corehelpers/inventory_test.go         | 61 ++++++++++++++++++----
 8 files changed, 108 insertions(+), 28 deletions(-)
```
  Reviewed each file in the diff. gofmt, go vet, shellcheck and git diff --check clean. No debug artifacts.

## Decisions

- **D-01 DECIDE & STATE: return shape.** Third return value `[]UnparsedSummaryRecord` rather than a result struct. Smallest change to the three callers. The record carries a repository-relative, slash-separated path computed with filepath.Rel from the repository root.
- **D-02 DECIDE & STATE: fix-hint entry added.** The commands matrix test (commands_test.go, "lacks exact actions") pins that every non-info finding carries exact next and verification argv that are not family-wide. The associate matrix case does not hit this code today, but the new warning would otherwise ship with empty argv against that house contract. Next is `sed -n '/^## Implementation Summary/,$p' -- <req>`, verification is `grep -n '`' -- <req>` (shows the backtick lines). The new subtest asserts both are non-empty.
- **D-03 DECIDE & STATE: evidence wording.** Evidence is the parser message plus "; this REQ file claims no paths", so the stderr line tells the reader the consequence. The runtime prints the REQ file on the following `  paths:` line. If the orchestrator wants the path on the same line as the code, prefix it in the evidence string in handleAssociate; one-line change.
- **D-04 DECIDE & STATE: JSON mode.** The finding also appears under `--format json`, as the maintainer accepted at plan approval. Not changed.

## Discovered Tasks

- None found that need action. → report only

## Lessons read

- Prime traps in prime-do-work-cli.md (silent-skip-reads-as-red, closed-enumeration-for-a-condition read in full), and the satellite bullets for families single-exit-owner (REQ-525), silent-skip-reads-as-red (REQ-566), closed-enumeration-for-a-condition (0.303.7, REQ-460, REQ-461, 0.283.0) by grep in lessons-do-work-cli.md. Primes prime-shell-commands.md and prime-action-files.md read in full. Whole satellites were dropped for budget per the REQ; none missing.
- Lesson candidate for the orchestrator (family suggestion: per-record-error-in-lookup-walk): a parse error in one record inside a lookup walk must not end the walk; the old branch also keyed recovery on a message substring, which is the closed-enumeration shape.

## End-to-end transcripts

Scratch repo outside both trees: committed archive with REQ-502 (unmatched bullet) and REQ-503 (claims good-file.txt), uncommitted good-file.txt. stdout shown with `cat -t`, so ^I is a TAB.

### RED (dispatch behaviour)
```
$ protected-inventory.sh --repo-root <scratch> start
exit=0
--- stdout
A	good-file.txt
--- stderr
$ protected-inventory.sh --repo-root <scratch> associate
exit=2
--- stdout
PARSE-FAILED: unmatched backtick in Implementation Summary
--- stderr
```

### GREEN (commit e86b314d)
```
$ protected-inventory.sh --repo-root <scratch> start
exit=0
--- stdout
A	good-file.txt
--- stderr
$ protected-inventory.sh --repo-root <scratch> associate
exit=0
--- stdout
REQ-503^Igood-file.txt
--- stderr
finding ASSOCIATION-SUMMARY-UNPARSED [warning]: unmatched backtick in Implementation Summary; this REQ file claims no paths
  paths: do-work/archive/UR-301/REQ-502-unmatched-summary.md
```

## Integration seams

None. `AssociateProjectPaths` has no callers outside corehelpers. `handleProtectedInventory` still rebuilds exact text from ASSOCIATION-FOUND findings only, and the warning stays in result.Findings, so the runtime prints it.
