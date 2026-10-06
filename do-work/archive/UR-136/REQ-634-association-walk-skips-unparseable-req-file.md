---
id: REQ-634
title: 'A REQ file the association walk cannot parse claims no paths, and do-work commit continues'
status: completed
route: A
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
  - Route A
  - 9-file write set
  - 3 subsystems involved
  - 7 acceptance criteria
  calculated_at: 2026-10-06T15:49:54Z
created_at: 2026-10-06T15:27:35Z
user_request: UR-136
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go, _dev/tests/contracts/core-checks.sh, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work/CHANGELOG.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, do-work/lessons-index.md]
review_at: 2026-10-06T16:04:59Z
heavy_verified_at: 2026-10-06T16:07:11Z
heavy_verified_revision: b9e04c02cb6de9616bdb7cfbd6cdc1519d21149b
commit: b9e04c02cb6de9616bdb7cfbd6cdc1519d21149b
kb_status: pending
integration_at: 2026-10-06T15:55:42Z
builder_handback_at: 2026-10-06T15:55:05Z
dispatch_at: 2026-10-06T15:49:54Z
claimed_at: 2026-10-06T15:47:41Z
completed_at: 2026-10-06T16:07:32Z
release_at: 2026-10-06T16:07:32Z
---
# A REQ File the Association Walk Cannot Parse Claims No Paths, and do-work Commit Continues

## What
When `AssociateProjectPaths` (`skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`) hits a REQ file whose Implementation Summary fails to parse, that one file claims no paths, the walk continues, and `protected-inventory associate` exits 0 with its owner rows on stdout and one warning on stderr naming the skipped file. The `PARSE-FAILED` exit-2 path goes away. Nothing about how a path is written changes.

## Why
A consumer reported (do-work 0.305.68) that one archived REQ with an odd number of backticks on one Summary line makes every `do-work commit` in a ~2,000-REQ repository exit 2 at Step 3, permanently, because archive records are immutable. The walk at `inventory.go:301-304` returns the parse error from the WalkDir callback, so a single-file problem becomes a walk-level failure. The maintainer's position, recorded at capture: the commit must not die on a formatting issue in one record, and the parser must not grow a grammar for every path representation a non-deterministic model can produce. So the fix is tolerance at the walk, not a smarter parser.

## Detailed Requirements
1. In `AssociateProjectPaths`, a parse error from `allBacktickedPaths` on one REQ file means that file claims no paths. Return nil from the callback for that file and keep walking. The function returns the skipped files alongside the association map so the caller can report them.
2. `handleAssociate` (`inventory.go:213-218`) drops the `PARSE-FAILED` branch and emits one warning finding per skipped file, code `ASSOCIATION-SUMMARY-UNPARSED`, with the REQ file's repository-relative path as the affected path and the parser's message as evidence. Outcome stays success, exit 0. In shim text mode the stdout rows (`<owner>\t<path>`) stay exactly as they are; the warning reaches stderr through the runtime's existing finding printer (`internal/commandruntime/command_runtime.go:95-104`), not through a direct `os.Stderr` write in the handler.
3. Flip the two lock-ins that pin the old behaviour: the `associate_unmatched` probe in `_dev/tests/contracts/core-checks.sh` (around lines 365-389) and the `preserves PARSE-FAILED` subtest in `inventory_test.go` (around lines 636-660). Both now assert: exit 0, the good file's owner row present on stdout, `PARSE-FAILED` absent, and the stderr finding naming the unparseable REQ file. Add a unit test on `AssociateProjectPaths` with one unparseable archived REQ and one claiming archived REQ: the claiming REQ's path is associated and the unparseable file is reported once.
4. Rewrite the sentence at `skills/do-work/actions/commit.md:85` and `skills/do-work-toolbox/actions/inspect.md:117` so `PARSE-FAILED` is no longer listed as an exit-2 case; say instead that a REQ file whose Summary cannot be parsed claims nothing and is named on stderr, and the agent reports it. Add one bullet under "What `associate` settles" in `skills/do-work/docs/prescribed-shell-primitives.md` stating the same rule.
5. Release: changelog entry and version bump per `_dev/primes/prime-releases.md` (shipped Go, action prose, and the docs guide all change).
6. Lessons: one entry in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (a per-record parse error inside a lookup walk must not become a walk-level failure) and the matching token refresh in `do-work/lessons-index.md`, in the same commit.

## Constraints
- No bullet-joining, no continuation-line grammar, no "first span is the path" rule. `allBacktickedPaths`, `firstBacktickedPaths` and `qualificationSummaryEntries` are untouched; qualify (`checks.go:321`) and scope-drift (`checks.go:152`) stay strict because there the REQ being parsed is the one being finalized and the author can fix it.
- The multi-path contract stays: `TestAssociationParserRetainsEveryClosedToken` and the core-checks.sh multi-path probe are unchanged.
- A skipped REQ's own files become unassociated and go to commit Step 4 grouping. Accepted.
- The reporter's optional in-flight marker (`ASSOCIATION-FOUND-INFLIGHT`) is dropped. In-flight `working/` claims stay included as documented at `inventory.go:257-263`; no marker, no confirm step.
- The reporter asked for no JSON code. The finding is the CLI's only stderr channel, so it also appears under `--format json`. The maintainer accepted this at plan approval.

## Builder Guidance
Certainty is high. The code change is small: the callback swallows the error into a skipped list, the handler turns the list into warning findings, one branch is deleted. Most of the work is the two test inversions, three prose sites, the changelog, and the lessons entry. Builder latitude: the exact evidence wording of the stderr line, and how `AssociateProjectPaths` returns the skipped list (second return value or a small result struct).

## Red-Green Proof
**RED prompt/case:** A fixture repository with an uncommitted `good-file.txt`, an archived completed REQ-502 whose Summary bullet is the existing core-checks.sh unmatched probe body (one closed `legacy-file.txt` span, then a second span opened with a backtick and never closed), and an archived completed REQ-503 claiming `good-file.txt`. Run `scripts/protected-inventory.sh start` then `scripts/protected-inventory.sh associate`.
**Why RED now:** `associate` exits 2, prints `PARSE-FAILED: unmatched backtick in Implementation Summary`, and prints no owner row for `good-file.txt`.
**GREEN when:** `associate` exits 0, stdout contains the row `REQ-503<TAB>good-file.txt`, stderr contains one finding line naming the REQ-502 file, and `PARSE-FAILED` appears nowhere.
**Validation:** User confirmed (2026-10-06, decisions D1 to D4 at capture)

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18027 tokens, over the 2000 budget; `slugged: partial`). Matched: this REQ changes a corehelpers classifier walk and the finding a command emits.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over the 2000 budget; `slugged: partial`). Matched: a prescribed command block's documented exit contract changes in commit.md and inspect.md.
- `_dev/primes/lessons-action-files.md` as a whole satellite (5879 tokens, over the 2000 budget; `slugged: partial`). Matched: action prose in commit.md and inspect.md changes.

## Full Context
See `do-work/user-requests/UR-136/input.md` for complete verbatim input.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Signature stub first so the tests compile, RED at three layers (unit test, shim subtest, shell probe) plus the end-to-end fixture, then: the walk callback appends an unparsed-summary record and returns nil; the handler deletes the message-string branch and appends one warning finding per skipped file; a fix-hint entry; three prose sites. Keyed on the condition "the parser returned an error", not on a message string. (From the builder hand-back.)
- [x] **[APPLY]:** Done as planned, within the write boundary; commands.go and commands_test.go were the only files outside the brief's list, both for the fix-hint contract and the changed signature. (From the builder hand-back.)
- [x] **[UNIFY]:** `git diff --stat` 8 files, 108 insertions, 28 deletions, each file reviewed; gofmt, go vet, shellcheck and `git diff --check` clean; no debug artifacts. (From the builder hand-back.)
*Source: consumer bug report pasted into `do-work validate-feedback` on 2026-10-06, triaged with the maintainer; the full text is the UR's verbatim input.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names every file and the exact change at each (one swallowed error in the walk callback, one finding in the handler, two test inversions, three prose sentences, release and lessons entries), and the triage that produced it already read the code paths. Exploration would re-discover what capture recorded.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modified)
- `_dev/tests/contracts/core-checks.sh` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work-toolbox/actions/inspect.md` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**What was done:** `AssociateProjectPaths` records a REQ file whose Implementation Summary fails to parse as claiming no paths and keeps walking, returning the skipped files as a third value; `handleAssociate` drops the `PARSE-FAILED` exit-2 branch and emits one `ASSOCIATION-SUMMARY-UNPARSED` warning per skipped file, which the runtime prints to stderr in shim text mode while the stdout owner rows stay unchanged. The two lock-ins that pinned `PARSE-FAILED` were flipped, a unit test pins the skip-and-keep-walking walk, and the commit and inspect actions plus the shell primitives guide state the new rule. Parsers, qualify and scope-drift are untouched.

## Qualification

**Range:** `09c5ddf7..b9e04c02` (builder commit e86b314d merged as b9e04c02)
**Gate record:** `advance` qualification satisfied with no findings: no debug artifacts, no leftover output primitives, every P-A-U box ticked from the hand-back, no `do-work/` path inside the range.
**Orchestrator judgment:** The diff is substantive and matches every requirement. R1: the WalkDir callback appends an `UnparsedSummaryRecord` and returns nil (`inventory.go`), so the file claims nothing and the walk continues. R2: the `PARSE-FAILED` branch is deleted and one `ASSOCIATION-SUMMARY-UNPARSED` warning per skipped file is appended with the repository-relative REQ path; the handler writes nothing to stderr itself, the runtime's finding printer does. R3: `TestProtectedInventoryCompatibilityShimPreservesErrors` and the `associate_unmatched` probe in `core-checks.sh` now assert exit 0, the good owner row, no `PARSE-FAILED`, and the skipped file named on stderr; `TestAssociationSkipsUnparseableSummaryAndKeepsWalking` pins the walk. R4: `commit.md`, `inspect.md`, and the "What `associate` settles" list carry the new sentence. Constraints hold: `checks.go` is untouched, the multi-path contract test is unchanged, the in-flight rule is unchanged. Two files outside the brief's list, `commands.go` and `commands_test.go`, are the fix-hint entry the house finding contract requires and a signature-only caller update, both judged in scope. Live flow verified by the builder's end-to-end transcript on a scratch repository (RED exit 2 with `PARSE-FAILED`; GREEN exit 0 with `REQ-503<TAB>good-file.txt` on stdout and the warning naming `REQ-502` on stderr).

## Testing

**Tests run:** `bash do-work/runs/work-2026-10-06-154820/REQ-634-probe.sh` (go test, corehelpers, `-run 'Associat|ProtectedInventory|CompatibilityShim|TerminalSuccess'`) through the `advance` test gate: probe status 0, gate record satisfied. Builder-side: `go test -count=1 ./internal/corehelpers/ ./internal/commandruntime/` (7.3s and 0.8s package wall, every file under the 30s budget), `bash _dev/tests/contracts/core-checks.sh` (exit 0, 4s), gofmt, go vet, shellcheck, `git diff --check` all clean.
**Result:** ✓ All passing
**Repository gate:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at b9e04c02, exit 0 on the first run, gate wall 171s; both Go stages EXECUTING (reuse_disabled); slowest files strict_behavior_regression_test.go 22.33s and finalization_recovery_test.go 28.64s, both under the 30s budget.

**Red-green validation:** (traces to `## Red-Green Proof`)
- `TestAssociationSkipsUnparseableSummaryAndKeepsWalking` (inventory_test.go): ✗ before, "one unparseable REQ file must not fail the walk: unmatched backtick in Implementation Summary" → ✓ after
- `TestProtectedInventoryCompatibilityShimPreservesErrors/unparseable summary claims no paths and warns`: ✗ before, outcome failure with ExitCodeOverride 2 → ✓ after
- `_dev/tests/contracts/core-checks.sh` associate_unmatched probe: ✗ before, "got exit 2, stdout: PARSE-FAILED: unmatched backtick" → ✓ after (exit 0, `REQ-503<TAB>legacy-file.txt` on stdout, REQ-502 file named on stderr)
- End-to-end on a scratch repository with `scripts/protected-inventory.sh`: before, `associate` exit 2 `PARSE-FAILED`; after, exit 0, stdout `REQ-503<TAB>good-file.txt`, stderr `finding ASSOCIATION-SUMMARY-UNPARSED [warning]: ... paths: do-work/archive/UR-301/REQ-502-unmatched-summary.md`. This is the captured RED/GREEN pair exactly.

**New tests added:**
- `TestAssociationSkipsUnparseableSummaryAndKeepsWalking` (inventory_test.go)

**Existing tests updated (cross-REQ impact):**
- `TestProtectedInventoryCompatibilityShimPreservesErrors` subtest `preserves PARSE-FAILED` (from REQ-603, preserved shim diagnostics) renamed and inverted: the shim no longer emits PARSE-FAILED by design; the subtest now pins the warning finding and the untouched owner rows.
- `_dev/tests/contracts/core-checks.sh` associate_unmatched probe (from REQ-539's contract split) inverted for the same reason.
- `TestGenericAssociationNeverOwnsSharedDoWorkMetadata`, `TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest`: signature-only caller updates.

**Heavy verification plan:**
- Range: 09c5ddf7..b9e04c02 (6 of 6 lanes: `_dev/tests/contracts/core-checks.sh` matches the `_dev/tests` subtree every lane covers)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests (+4 more matches)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests (+7 more matches)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests (+4 more matches)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests (+4 more matches)

*Verified by work action*

## Review

**Overall: 97%** | 2026-10-06T16:03:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 98% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** None
**Acceptance:** Pass — reviewer reproduced the captured RED (pre-change code: exit 2, `PARSE-FAILED`, no owner row) and GREEN (exit 0, `REQ-503<TAB>good-file.txt` on stdout, `ASSOCIATION-SUMMARY-UNPARSED` naming REQ-502 on stderr, no `PARSE-FAILED`) on a scratch repository; the finding also appears under `--format json`; corehelpers Go tests and core-checks.sh pass.
**Suggested testing:** 3 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

**Base revision:** 09c5ddf7f4ad5842f47b25e61446641ac5ee907e
**Target revision:** b9e04c02cb6de9616bdb7cfbd6cdc1519d21149b
**Selected lanes** (6 of 6; every lane covers the `_dev/tests` subtree that `_dev/tests/contracts/core-checks.sh` sits in):
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `_dev/tests/contracts/core-checks.sh` matched subtree `_dev/tests`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same match
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — same match, plus the four `internal/corehelpers` files matched `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — same match, plus every changed path under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — same match, plus the `skills/do-work/tools/do-work-cli` files
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — same match, plus the `skills/do-work/tools/do-work-cli` files

## Heavy Verification Result

**Target revision:** b9e04c02cb6de9616bdb7cfbd6cdc1519d21149b
**Execution revision:** b9e04c02cb6de9616bdb7cfbd6cdc1519d21149b (detached worktree `.git/work-run-2026-10-06-154820/drain-head`, `QUEUE_KANBAN_BROWSER` set, 2026-10-06T16:00:37Z to 16:06:44Z)
- queue-kanban-javascript: executed, exit 0, 11s
- queue-kanban-browser: executed, exit 0, 104s
- do-work-cli-integrations: executed, exit 0, 87s
- staged-skills: executed, exit 0, 40s
- updater: executed, exit 0, 69s
- installer: executed, exit 0, 29s

## Lessons Learned

**What worked:** Triage before capture: the validate-feedback pass had already read every code path, so a Route A dispatch with the exact design in the brief produced a first-try green build (8 files, 108 insertions) with RED at three layers. Running the six heavy lanes from a detached worktree at the merge revision while the review ran, instead of after it.
**What didn't:** The first `advance` estimate call and the first REQ edit both tripped on string matching (`effort_estimate:` contains `estimate:`; the plan JSON keys lanes as `command_argv`, not `argv`). Assert on the exact line start, and print a record's keys before composing prose from it.
**Worth knowing:** The runtime, not the handler, is the only stderr writer: a warning finding under `ExactTextOutput` is printed to stderr by `command_runtime.go` after the exact text, so "say it on stderr" means "emit a warning finding", and the finding is also visible under `--format json`. A lookup walk over immutable archive records must never inherit a validator's refusal: strictness stays where the author can fix the file (qualify, scope-drift). The shell contract probe captures stderr to its own file because the warning and the owner rows travel on different streams.

## Orientation

Now `do-work commit` survives a consumer archive holding one request whose Implementation Summary does not parse: association completes, the skipped record is named on stderr, and the owner rows are unchanged. Lives in the do-work-cli corehelpers subsystem (`prime-do-work-cli.md`) with its prose contract in the commit and inspect actions and the shell primitives guide. No map change: one exit code case became a warning; no new module, data flow or renamed concept. Prime spot-check: `prime-do-work-cli.md` names no path this change touched, so it is not stale.

## Timing

Observed 2026-10-06T15:49:54Z to 2026-10-06T16:07:04Z: 17m 10s total, 19m 18s attributed across 4 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 9m 38s | 2 |
| builder-work | 5m 11s | 1 |
| review | 4m 29s | 1 |

Slowest stage: verification-gate / heavy drain, six lanes at b9e04c02 from the detached worktree, 6m 27s, outcome success.
