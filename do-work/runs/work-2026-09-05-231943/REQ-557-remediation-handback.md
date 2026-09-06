# REQ-557 remediation hand-back — two behaviour changes now pinned

REQ-557 is the request that deduplicated six Go helper names across the CLI's internal
packages. The three-lens review accepted the engineering and asked for two tests. Both
are added. No production code changed. Nothing outside `internal/finalization` and
`internal/knowledgecommands` was touched.

- Worktree: `/home/user/skill-do-work-worktrees/worktree-agent-REQ-557-deduplicate-go-helper-names`
- Branch: `worktree-agent-REQ-557-deduplicate-go-helper-names`
- Commit: `764c1c210df64b4a8c31b85849b3bebdf83c986c` — `[REQ-557] Pin the two REQ-557 behaviour changes nothing tested`

## Files changed

- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-557-deduplicate-go-helper-names/skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go`
  — one new test, plus the corrected header comment.
- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-557-deduplicate-go-helper-names/skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_version_migration_test.go`
  — new file, one table test.

## T1 — the commit_paths guard now has a call-site test

`TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget` seeds a
repository with `newFinalizationRepository`, commits a claimed REQ-721 plus a checkpoint
entry, and writes a manifest whose `commit_paths` lists the request, the checkpoint and
`implementation.txt` but NOT the planned `do-work/archive/REQ-721.md`. It calls
`prepareBoundJournal` and asserts three things: the error text contains
`commit_paths omits planned lifecycle or release targets`, the error names the omitted
archive path, and no journal file is left behind.

### Green at HEAD

```
=== RUN   TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet
--- PASS: TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet (0.00s)
=== RUN   TestMissingCommitPathsNamesTheRequiredTargetsCommitPathsOmits
--- PASS: TestMissingCommitPathsNamesTheRequiredTargetsCommitPathsOmits (0.00s)
=== RUN   TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget
--- PASS: TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget (0.09s)
PASS
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	0.102s
```

### RED under mutation (b) — the `if len(missing) > 0 { ... }` block deleted

```
=== MUTATION (b) rerun ===
--- FAIL: TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget (0.23s)
    finalization_req557_test.go:110: prepare accepted commit_paths that omit do-work/archive/REQ-721.md and wrote a journal for do-work/archive/REQ-721.md with effective paths [do-work/CHECKPOINT.md do-work/working/REQ-721.md implementation.txt] (resumed=false); the omitted-target guard no longer fires
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	0.251s
FAIL
```

The mutation was reverted; `git diff --stat` after the revert showed only the test file.

### NOT red under mutation (a) — and it cannot be, from any accepted manifest

Mutation (a) is the call site rewritten as
`missing, _ := missingCommitPaths(requiredCommitPaths, effectiveCommitPaths)`, with the
`if missingError != nil` block deleted. Measured with the new test in place:

```
=== MUTATION (a): the call site discards the helper's error ===
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	0.121s
=== whole package under mutation (a) ===
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	24.352s
```

The reason is structural, not a gap in the test. The error the call site would discard
can only come from `normalizeRepositoryPaths` refusing an empty, absolute, or escaping
REQUIRED path, and the required set is built from tool-derived paths only:

- `statePlan.TargetPaths` — every entry is a path `requeststate` resolved out of the
  repository snapshot. `internal/requeststate/state_plan.go:111-112` refuses anything
  else with `REQUEST-SNAPSHOT-STALE` before a plan is ever runnable, so a manifest
  `request_path` that escapes the repository never reaches the guard.
- release postimage paths — every one passed `containedPath` in
  `internal/publication/release.go:36` and `:61`, refusing with `RELEASE-PATH-UNSAFE`.

So mutation (a) changes no observable behaviour for any manifest the validators accept.
It is invisible to every behaviour test, including this one. What still covers it is
`TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet`, and only
for a discard written INSIDE `missingCommitPaths`. See "Found and not fixed" below.

### Comment correction at `finalization_req557_test.go`

The old header ended: "This test fails the moment the error is discarded again, whatever
the helper is called." That is true only for a discard inside the helper. The header now
states what each of the three tests pins — helper propagation, helper answer, call-site
refusal — says the third is the only one that fails when the `if len(missing) > 0`
refusal is deleted, and records why a call-site test cannot pin propagation, naming
`REQUEST-SNAPSHOT-STALE` and `RELEASE-PATH-UNSAFE` as the validators that make an
unusable required path unreachable.

## T2 — the D-03 semver refusal now has a table test

`TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder` in
`internal/knowledgecommands/interview_version_migration_test.go`, five rows:

| session version | template version | expected |
| --- | --- | --- |
| `1.0.x` | `1.0.1` | refusal `template and session versions must be bare semver`, version not stamped |
| `1.0.0` | `1.09.0` | same refusal (a stamp -> error pair), version not stamped |
| `1.0.0` | `1.0.` | same refusal, version not stamped |
| `1.0.0` | `1.1.0` | stamped `1.1.0`, unchanged behaviour |
| `1.0.0` | `1.0.0` | left at `1.0.0`, unchanged behaviour |

### Green at HEAD

```
--- PASS: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder (0.00s)
    --- PASS: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/session_patch_part_is_not_a_number (0.00s)
    --- PASS: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/template_minor_part_carries_a_leading_zero (0.00s)
    --- PASS: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/template_patch_part_is_empty (0.00s)
    --- PASS: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/both_versions_are_bare_semver_and_the_template_is_newer (0.00s)
    --- PASS: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/both_versions_are_bare_semver_and_equal (0.00s)
PASS
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/knowledgecommands	0.003s
```

### RED with the `if !versionsParsed` refusal removed

Ablation: `versionOrdering, versionsParsed := ...` plus its refusal replaced by
`versionOrdering, _ := sharedprimitives.CompareSemanticVersions(old, template.Version)`.

```
=== ABLATION: versionsParsed refusal removed ===
--- FAIL: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder (0.00s)
    --- FAIL: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/session_patch_part_is_not_a_number (0.00s)
        interview_version_migration_test.go:78: migrateInterviewSession("1.0.x" -> "1.0.1") accepted a version the strict parser cannot read; the lenient comparator's numeric scoring is back
    --- FAIL: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/template_minor_part_carries_a_leading_zero (0.00s)
        interview_version_migration_test.go:78: migrateInterviewSession("1.0.0" -> "1.09.0") accepted a version the strict parser cannot read; the lenient comparator's numeric scoring is back
    --- FAIL: TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder/template_patch_part_is_empty (0.00s)
        interview_version_migration_test.go:78: migrateInterviewSession("1.0.0" -> "1.0.") accepted a version the strict parser cannot read; the lenient comparator's numeric scoring is back
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/knowledgecommands	0.003s
FAIL
```

The ablation was reverted from a byte copy taken before it was applied.

## Measured version-pair grid — the record's numbers should come from here

Method: a throwaway test in `internal/knowledgecommands` ran the SHIPPED
`migrateInterviewSession` beside a reconstruction of the pre-REQ-557 function (identical
except that the same-major branch uses the deleted lenient `compareSemver`, recovered
verbatim from commit `a8ba69c`). Both ran over the full ordered cross product of one
version list. Outcome per pair: ERROR / MIGRATE / STAMP (`template_version` rewritten) /
SKIP (left alone). The scratch file was deleted before the commit.

Version list (21 values):
`1.0.0`, `1.0.1`, `1.1.0`, `1.0.10`, `1.10.0`, `2.0.0`, `1.0.x`, `1.09.0`, `01.0.0`,
`1.01.0`, `1.0.01`, `1.0.`, `1..0`, `1.0.0-rc1`, `1.0.0+1`, `1.0.-1`, `1.0.0 ` (trailing
space), `1.0`, `1.0.0.0`, `x.0.0`, `` (empty).

Totals:

```
pairs=441 changed=242
transition SKIP->ERROR = 159
transition STAMP->ERROR = 83
predicted-by-rule=242 actual-changed=242
```

Two facts the record should carry:

1. Every changed pair moves to ERROR. Nothing that used to error now succeeds, no stamp
   became a skip, no skip became a stamp. The change is one-directional.
2. The changed set has an exact rule, verified against the measurement (242 predicted =
   242 measured): `semverMajor` accepts both sides, the two majors are equal, and at
   least one side fails `sharedprimitives.ParseSemanticVersion`. `semverMajor` only
   requires three dot-separated parts with an `Atoi`-able first part, so a leading zero,
   an empty part, a trailing dot, a non-numeric or negative third part, a `-rc1` or `+1`
   suffix and a trailing space all get through it and used to be scored numerically.

The 83 stamp -> error pairs, grouped by session version (this is the direction the
record does not currently describe — the session was silently re-stamped before, the
interview command aborts now):

```
session "1.0.0"      -> 1.09.0, 1.01.0, 1.0.01
session "1.0.1"      -> 1.09.0, 1.01.0
session "1.1.0"      -> 1.09.0
session "1.0.10"     -> 1.09.0, 1.01.0
session "1.0.x"      -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session "1.09.0"     -> 1.10.0
session "01.0.0"     -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session "1.01.0"     -> 1.10.0, 1.09.0
session "1.0.01"     -> 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0
session "1.0."       -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session "1..0"       -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session "1.0.0-rc1"  -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session "1.0.0+1"    -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session "1.0.-1"     -> 1.0.0, 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.0.x, 1.09.0, 01.0.0,
                        1.01.0, 1.0.01, 1.0., 1..0, 1.0.0-rc1, 1.0.0+1, "1.0.0 "
session "1.0.0 "     -> 1.0.1, 1.1.0, 1.0.10, 1.10.0, 1.09.0, 1.01.0, 1.0.01
session ""           -> 1.09.0, 1.01.0, 1.0.01   (empty session version defaults to 1.0.0)
```

The count depends on which versions are in the grid, so quote the rule, not the 83.
`template.Version` is read verbatim from a project-authored template's YAML
frontmatter, so the template side of every one of these pairs is project data. The
session side is tool-written, but a session file carried forward from an older template
keeps whatever that template stamped.

## Gate

Run from the worktree, pointed at the worktree:

```
DO_WORK_GATE_ROOT=.../worktree-agent-REQ-557-deduplicate-go-helper-names bash .../gate.sh
EXIT=0
```

Last line: `Maintainer verification passed.` Gate wall 98s. `do-work-cli` module
reported `wall=31s tests=796`. The two heavy verification tests that fail at the branch
point for environmental reasons did NOT appear in this run — nothing failed.

Also clean before the commit: `go build ./...`, `go vet ./...`, `gofmt -l .` empty, and
`go test ./internal/finalization/ ./internal/knowledgecommands/` green.

## Found and not fixed

1. **A call-site discard of the helper's error is pinned by nothing, and no behaviour
   test can pin it.** Detail and evidence above. It is behaviour-neutral today; it
   becomes a live hole the day a manifest-derived path joins the required set. Two ways
   to close it, neither taken here because both go past "add a test":
   - Collapse the helper and the guard into one function returning a single `error`
     (nil when covered, the normalize refusal, or the "commit_paths omits..." error).
     A one-result signature makes `missing, _ :=` unwritable, and the helper test then
     pins both halves. This is a production-code change to a change the review said not
     to rework, so it needs a decision, not a builder.
   - A source-text lock-in asserting the call site does not discard. Cheap, but it
     breaks on any rename and pins text rather than behaviour.
2. **`_dev/tests/audit-lockins.sh` was left alone as instructed** — another builder is
   editing it in a different worktree. The optional improvement the review suggested
   there is untouched.
3. **`semverMajor` is now the loosest gate in the path.** It accepts anything with three
   dot-separated parts and an `Atoi`-able first part, which is why the refusal set is as
   wide as the grid shows. Tightening it to the strict parser would move most of those
   159 skip -> error pairs into the same refusal earlier and with a clearer message, but
   that is a further behaviour change and belongs in its own request.
