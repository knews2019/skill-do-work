# REQ-557 hand-back — deduplicating six Go helper names in do-work-cli

Branch `worktree-agent-REQ-557-deduplicate-go-helper-names`, head `a8ba69cb4c518d842a0f1990c9fcfb8b7e4baa40`.
One commit. Nothing staged or committed in the main checkout, which is clean at `adfbaee`.

**Bottom line first:** fifteen definitions of six helper names became seven. Every plan step landed. Two
things differ from the plan and both are named below: the write set grew by one test file, and the
lock-in pins seven, not six, because the plan's own number was only reachable by counting a name that
decision D-04 chose never to create.

## Files changed

New:
- `skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go`
- `skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go` — **the one widening**

Modified:
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go`
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
- `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go`
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go`
- `skills/do-work/tools/do-work-cli/internal/publication/release.go`
- `skills/do-work/tools/do-work-cli/internal/publication/release_mirrors.go`
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go`
- `_dev/tests/audit-lockins.sh`
- `do-work/working/REQ-557-deduplicate-six-go-helper-names-defined-fourteen-times-across-do-work-cli.md`
  — write_set frontmatter, `## Scope` list, the P-A-U boxes, and the `## Implementation Summary` and
  `## Decisions — implementation` sections the request requires ("a silent pick is a review refusal")

## The write-set widening

**One path added: `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go`.**

D-02 requires a test that fails when the restored guard is disabled again. The plan's Scope named
`finalization_recovery_test.go` as the only finalization test file it would touch, but that file is a
recovery-fixture file and the guard test is not a recovery fixture. This package already has a per-REQ
test-file convention — `finalization_req499_test.go`, `_req512_`, `_req547_`, `_req560_`, `_req565_` —
so the new test went in `finalization_req557_test.go`.

Declared in three places: the request's `write_set` frontmatter, the request's `## Scope` "Files I will
touch" list, and the commit message. `do-work-cli scope-drift --request-path …` prints
`OK: Implementation Summary matches the Scope declaration` and exits 0.

Nothing else outside the plan's declared set was touched.

## Decisions, each with its evidence

### D-01 — `UniqueSortedStrings` keeps every value, including the empty string

Two of the four deleted copies filtered blanks (`knowledgecommands`, `repairvalidation`), two did not
(`corehelpers`, `finalization`). Evidence that the filter cannot fire today, read at every producer:

- `corehelpers`: `firstBacktickedPaths` and `allBacktickedPaths` both append only when the backtick span
  is non-empty (`if rest[:second] != ""`, `if parts[index] != ""`); `scopeDeclaredPaths` gets its paths
  from `backtickedPathsOnLine`, which does the same check.
- `repairvalidation`: `gitPaths` passes `strings.FieldsFunc` output, which by definition never yields an
  empty field; `gitStatusPaths` passes `record[3:]` of a record already checked `len >= 4`.
- `knowledgecommands`: every site passes map keys of a write set, or `filepath.Join` / `filepath.Dir`
  output, neither of which returns "".

Dropping the filter is the safe direction: a blank that ever does appear becomes visible in an evidence
list instead of vanishing from it. The parameter is `values`, not `paths`, because
`already_green.go:283` passes reason codes. The result is always a non-nil empty slice — every deleted
copy also guaranteed that, so JSON evidence still marshals `[]` and never `null`; pinned by
`TestUniqueSortedStringsReturnsAnEmptySliceRatherThanNil`.

### D-02 — the finalization wrapper is deleted and the guard it disabled is restored (behaviour change)

The wrapper's whole body was `result, _ := normalizeRepositoryPaths(paths); return result`. At
`finalization_prepare.go` a required commit path that was empty, absolute, or escaping made
`normalizeRepositoryPaths` return `(nil, error)`; the error was thrown away, so the whole required set
became empty, the subtraction found nothing, `len(missing)` was 0, and the
`commit_paths omits planned lifecycle or release targets` error on the very next line could never fire.

That computation now lives in `missingCommitPaths(requiredCommitPaths, effectiveCommitPaths)`, which
propagates the refusal — exactly as the `normalizeRepositoryPaths` call nine lines above it in the same
function already does. This is deliberate: carrying a silently disabled guard forward under a new name
would ship the drift this request exists to remove.

The other 21 production call sites plus 1 test site are plain dedupe and now call
`sharedprimitives.UniqueSortedStrings`. All of them pass repository-relative slash paths produced by git
output, journal images (`imagePaths`), `requestRepositoryPath`, or manifest lock paths, so the
normalization the wrapper also performed was a no-op there. Where a path WOULD have been rejected, the
old behaviour emptied the entire refusal list; the new behaviour shows the offending path in it.

Test: `TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet` (red/green evidence
below) and `TestMissingCommitPathsNamesTheRequiredTargetsCommitPathsOmits`.

### D-03 — one orientation, one strict parser, an explicit `parsed` flag

The two copies were inverted, which is why picking one body and leaving both predicates alone would have
reversed the release guard with every test still green:
- `knowledgecommands/interview_commands.go` returned -1 when the FIRST argument was older.
- `publication/release.go` returned +1 when the SECOND argument was newer.

`CompareSemanticVersions(left, right) (int, bool)` uses the standard Go orientation (negative when left
is older) and `publication`'s strict rules verbatim: exactly three dot-separated parts, no empty part,
no leading zero on a multi-character part, digits only, non-negative.

- `publication/release.go` now reads
  `versionOrdering, versionsParsed := …; if !versionsParsed || versionOrdering >= 0 { refuse }`. Same
  input set it refused before: previously `compareSemver` returned 0 for unparseable input and the guard
  was `!= 1`, so unparseable already refused.
- `knowledgecommands` keeps its `< 0` predicate and now returns the section's own
  `"template and session versions must be bare semver"` error for a version its lenient parser used to
  score as zero. Before, such a version silently skipped the template-version stamp. This is the surfaced
  error the plan called for.
- `ParseSemanticVersion` is exported so `publication/release_mirrors.go`'s admission check keeps working
  when the private `parseSemver` goes.

Pinned by `TestCompareSemanticVersionsUsesTheStandardOrientation` and
`TestCompareSemanticVersionsReportsUnparseableInputInsteadOfScoringItEqual`.

### D-04 — `physicalPath` is not merged; two contracts keep two names

`suiteinstall/update_transaction.go`'s copy is renamed `existingPhysicalPath` (EvalSymlinks + Abs,
absence is an error, result always absolute). `knowledgecommands/commands.go` keeps `physicalPath`
(walks missing ancestors, absence succeeds). No behaviour change: three call sites renamed, one
definition renamed, one doc comment added saying which contract this one is and why the other exists.

Merging would cost an explicit existence check re-added in `resolveUpdateRoots`, where the missing-path
error IS the existence check for the installed skill root before the `strings.HasPrefix` containment
test decides whether the skill lives inside the project.

### D-05 — `RequestIDLess` uses the permissive parser

`repositorymodel.requestNumberFromText` stays private and is the parser the exported comparator uses. It
is also the parser that assigns request identity elsewhere in the same file, so ordering and naming
cannot disagree.

One observable difference, stated rather than hidden: `nextselection`'s deleted `requestIDLess` used the
strict `numericID`, which rejects any non-digit after the prefix. An id like `REQ-12x` previously fell
through to the plain string comparison in `nextselection` and now sorts by 12. Such an id cannot reach
either of the two `nextselection` sorts — both sort ids drawn from a repository snapshot, whose
filenames and `id:` fields are produced by `formatRequestID`. `numericID` itself is untouched and keeps
its nine other callers.

### D-06 — DEVIATION: the lock-in pins seven, not the six the plan states

This is arithmetic, not judgment, and here is the measurement.

The plan's literal regex names `physicalPath|ResolvePhysicalPath`. `ResolvePhysicalPath` was the merged
resolver from the competing plan; **D-04 decided not to create it**, and D-04 instead creates
`existingPhysicalPath`, which the plan's regex does not name. Run against the finished tree:

```
$ rg -n … '^func (…|physicalPath|ResolvePhysicalPath)\(' internal/   →  6
$ rg -n … '^func (…|physicalPath|existingPhysicalPath)\(' internal/  →  7
```

So the plan's six is reachable only by counting a phantom identifier and omitting the real one. A
lock-in in that shape names something that can never exist and leaves the one name the change introduces
unguarded — the "Closed Enumerations Go Stale" trap, applied to the very edit that creates the name.

The shipped block therefore pins **7**: five canonical helpers in `sharedprimitives` and
`repositorymodel`, plus the two deliberately separate path resolvers. The block says in place why the
total is seven and not six, so the next reader does not have to rediscover it.

Related, and the same lineage: the request's baseline is also off by one. Its Reproduce line claims
fourteen; at the revision this was built from it prints **fifteen**, because the audit missed
`internal/repairvalidation/already_green.go`'s fourth `uniqueSorted`. The plan's Pre-Flight repeated the
fourteen.

## Every deletion, traced

Nine definitions deleted, one exported in place, one renamed.

| Deleted | Canonical it now uses | Why every call site behaves the same |
|---|---|---|
| `corehelpers.uniqueSorted` (`checks.go`) | `sharedprimitives.UniqueSortedStrings` | Identical body (`stringSet` + sort, no empty filter), identical non-nil empty result. 4 sites: `checks.go` `scopeDeclaredPaths`, `firstBacktickedPaths`, `allBacktickedPaths`, `porcelainPaths`; plus 2 in `inventory.go`. |
| `corehelpers.subtractPaths` | `sharedprimitives.SubtractStringValues` | Identical body. 2 sites in `handleScopeDrift` (`missing`, `extra`). Left order and non-nil empty result preserved. |
| `corehelpers.firstError` | `sharedprimitives.FirstNonNilError` | Byte-identical body. 1 site, the `SCOPE-PATH-LIST-MALFORMED` finding. |
| `knowledgecommands.uniqueSorted` (`interview_commands.go`) | `sharedprimitives.UniqueSortedStrings` | Only difference is the empty-string filter, which no producer can trigger — D-01. 6 sites in `interview_commands.go`, 2 in `memory_commands.go`. |
| `knowledgecommands.compareSemver` | `sharedprimitives.CompareSemanticVersions` | 1 site, `migrateInterviewSession`. Same orientation, same `< 0` predicate. Strictness change is D-03 and it surfaces an error where the old code silently did nothing. |
| `repairvalidation.uniqueSorted` (the copy the audit missed) | `sharedprimitives.UniqueSortedStrings` | Only difference is the empty-string filter — D-01, producers checked. 5 sites. |
| `publication.firstError` (`capture_files.go`) | `sharedprimitives.FirstNonNilError` | Byte-identical body. 3 sites across `capture_files.go`, `release.go`, `answer.go`. |
| `publication.compareSemver` (`release.go`) | `sharedprimitives.CompareSemanticVersions` | 1 site, the release guard. Orientation flips, and the predicate flips with it in the same edit — D-03. Same inputs refused. |
| `publication.parseSemver` (`release.go`) | `sharedprimitives.ParseSemanticVersion` | Body copied verbatim into the shared package. 1 remaining site, `release_mirrors.go` `plainVersionFileValue`. |
| `dependencygraph.requestIDLess` | `repositorymodel.RequestIDLess` | 4 sites. Parsers agree (survey: 8,016 inputs, 160,000 ordered pairs). |
| `dependencygraph.requestNumber` | — deleted, no caller left | Its only caller was `dependencygraph.requestIDLess`. Confirmed by build: `strconv` became an unused import and was removed. |
| `nextselection.requestIDLess` | `repositorymodel.RequestIDLess` | 2 sites. Parser change is D-05; `numericID` stays for its nine other callers. |
| `finalization.subtractPaths` | `sharedprimitives.SubtractStringValues` | Identical body. 1 site, inside the guard rewritten by D-02. |
| `finalization.uniqueSorted` | `sharedprimitives.UniqueSortedStrings` at 21 sites; `missingCommitPaths` at 1 | This one is NOT behaviour-preserving at the guard site, and that is D-02. |

Kept on purpose, both verified as having callers outside this class:
- `corehelpers.stringSet` — still used by `checks.go` and `inventory.go`.
- `nextselection.numericID` — nine other callers, and the right validator for CLI target tokens and the
  `UR-` prefix.

## Red/green evidence, verbatim

### Build, vet, format — CLI module

```
$ cd skills/do-work/tools/do-work-cli && gofmt -l ./ ; echo "GOFMT_EXIT=$?" ; go build ./... && echo BUILD_OK && go vet ./... && echo VET_OK
GOFMT_EXIT=0
BUILD_OK
VET_OK
```

`gofmt -l` printed nothing (empty), exit 0.

### No import cycle: the shared package is a leaf

```
$ go list -deps ./internal/sharedprimitives | grep 'do-work-cli'
github.com/knews2019/skill-do-work/do-work-cli/internal/sharedprimitives
$ go list -deps ./internal/sharedprimitives | grep 'do-work-cli/internal' | grep -v 'internal/sharedprimitives$' || echo "(none)"
(none)
```

The only module package in its dependency closure is itself.

### Full test run — CLI module

```
$ go test ./... -count=1 -v   (with NODE_OPTIONS and the GIT_CONFIG_* triples unset)
GO_TEST_EXIT=0
top-level tests: 808
PASS: 794  SKIP: 14  FAIL: 0
--- last line ---
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	1.171s
```

794 passing against 784 at pre-flight; the ten added are eight in `sharedprimitives` and two in
`finalization`. The 14 skips are pre-existing.

### D-02 guard test, RED then GREEN

RED — the propagation replaced by the discard the wrapper used to do
(`normalizedRequired, _ := normalizeRepositoryPaths(requiredCommitPaths)`):

```
=== RED RUN (guard disabled: error discarded again) ===
--- FAIL: TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet (0.00s)
    finalization_req557_test.go:23: missingCommitPaths accepted the unusable required path "" and returned missing=[]string{}; the commit_paths guard is disabled again
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	0.016s
```

GREEN — propagation restored:

```
=== GREEN RUN (guard restored) ===
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	0.012s
```

### The lock-in, RED in three directions then GREEN

Direction 1 — an eighth definition reappears (a `uniqueSorted` wrapper appended to
`repairvalidation/already_green.go`):

```
FAIL: 8 definitions of the shared helper names under skills/do-work/tools/do-work-cli/internal; REQ-557 pinned exactly 7, one per canonical name plus the two deliberately separate path resolvers:
  …/internal/suiteinstall/update_transaction.go:215:func existingPhysicalPath(path string) (string, error) {
  …/internal/repositorymodel/repository_model.go:658:func RequestIDLess(leftID string, rightID string) bool {
  …/internal/sharedprimitives/shared_primitives.go:18:func UniqueSortedStrings(values []string) []string {
  …/internal/sharedprimitives/shared_primitives.go:34:func SubtractStringValues(leftValues, rightValues []string) []string {
  …/internal/sharedprimitives/shared_primitives.go:51:func FirstNonNilError(firstCandidate, secondCandidate error) error {
  …/internal/sharedprimitives/shared_primitives.go:65:func CompareSemanticVersions(leftVersion, rightVersion string) (int, bool) {
  …/internal/repairvalidation/already_green.go:459:func uniqueSorted(values []string) []string {
  …/internal/knowledgecommands/commands.go:147:func physicalPath(path string) (string, error) {
exit=1
```

Direction 2 — one canonical definition renamed away (`FirstNonNilError` → `firstNonNilErrorMovedAway`),
which is the floor half a one-sided ratchet would miss:

```
FAIL: 6 definitions of the shared helper names under skills/do-work/tools/do-work-cli/internal; REQ-557 pinned exactly 7, one per canonical name plus the two deliberately separate path resolvers:
  …/internal/knowledgecommands/commands.go:147:func physicalPath(path string) (string, error) {
  …/internal/suiteinstall/update_transaction.go:215:func existingPhysicalPath(path string) (string, error) {
true exit=1
```

Direction 3 — the scan itself cannot run (search root pointed at a path that does not exist), which is
the "Unchecked Exit Status Reads as Content" case: rg exits 2 with empty output, and a check that judged
only the text would read that as "no duplicates found":

```
FAIL: could not scan skills/do-work/tools/do-work-cli/internal-moved-away for shared-helper definitions (rg exit 2); the duplicate-helper ratchet did not run.
exit=1
```

GREEN, and `shellcheck -S warning` clean on the file:

```
$ bash _dev/tests/audit-lockins.sh ; echo "exit=$?"
Audit lock-in regressions passed.
exit=0
$ shellcheck -S warning _dev/tests/audit-lockins.sh && echo SHELLCHECK_OK
SHELLCHECK_OK
```

### The union command's output at the end — 7 lines

```
$ rg -n --glob '*.go' --glob '!*_test.go' \
  '^func (uniqueSorted|UniqueSortedStrings|subtractPaths|SubtractStringValues|requestIDLess|RequestIDLess|firstError|FirstNonNilError|compareSemver|CompareSemanticVersions|physicalPath|existingPhysicalPath)\(' \
  skills/do-work/tools/do-work-cli/internal/ | sort
skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go:147:func physicalPath(path string) (string, error) {
skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go:658:func RequestIDLess(leftID string, rightID string) bool {
skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go:18:func UniqueSortedStrings(values []string) []string {
skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go:34:func SubtractStringValues(leftValues, rightValues []string) []string {
skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go:51:func FirstNonNilError(firstCandidate, secondCandidate error) error {
skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go:65:func CompareSemanticVersions(leftVersion, rightVersion string) (int, bool) {
skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go:216:func existingPhysicalPath(path string) (string, error) {
```

Seven lines, not the six the task named. D-06 above is the reason, with the measurement.

The request's own Reproduce command — the six ORIGINAL names only — now prints one line, the
`knowledgecommands` path resolver that D-04 deliberately keeps:

```
skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go:147:func physicalPath(path string) (string, error) {
```

That is why the lock-in uses the union and not the request's Red-Green Proof line, which would read
"1 of 6" on a correct change.

### Gate

The wrapper `bash /tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gate.sh`
contains `cd /home/user/skill-do-work` on line 17, so running it FROM the worktree still gates the main
checkout. Both were run.

Exact wrapper, invoked from the worktree as instructed:

```
EXACT_WRAPPER_EXIT=0
last line: Maintainer verification passed.
maintainer-verify: gate wall 31s
```

That result is about the main checkout, which is clean at `adfbaee` and carries none of this work. The
same wrapper with only its `cd` retargeted at the worktree — identical environment scrub — is the run
that actually gates this branch:

```
WORKTREE_GATE_EXIT=0
last line: Maintainer verification passed.
maintainer-verify: stage do-work-cli-fast-tests: EXECUTING (fingerprint_mismatch)
go-test budget: module=…/worktree-agent-REQ-557-…/skills/do-work/tools/do-work-cli wall=21s tests=794 slowest-file=internal/nextselection/blocked_probe_test.go:6.96s limit=<30s
maintainer-verify: gate wall 53s
```

`fingerprint_mismatch` means the CLI stage really executed rather than reusing evidence, and 794 is the
gate's own count of passing tests, matching the direct run.

### Scope-drift check on the request document

```
$ DO_WORK_COMPATIBILITY_SHIM=1 do-work-cli scope-drift --request-path do-work/working/REQ-557-….md
OK: Implementation Summary matches the Scope declaration
exit=0
```

## Deviations from the plan

1. **D-06 — the lock-in pins seven, not six.** Full reasoning and measurement above. The task said
   "EXACTLY 6"; the union of old and new names is 7, and the two halves of that instruction cannot both
   be satisfied because D-04 chose not to create `ResolvePhysicalPath`. I took the definition half as
   carrying the intent and the number as a computed consequence of a stale draft.
2. **The write set grew by one file** — `finalization_req557_test.go`, declared in the frontmatter, the
   Scope list, and the commit message.
3. **`missingCommitPaths` is a named helper the plan did not ask for.** The plan says
   "`finalization_prepare.go:162` calls `normalizeRepositoryPaths` directly and propagates its error".
   Inlining it that way would have left the guard testable only through a full repository fixture
   driving `prepareBoundJournal`. Extracting the two lines into `missingCommitPaths` in the same file
   gives D-02's required test a direct target and changes no behaviour: the helper does exactly the
   normalize-then-subtract the plan describes, and its caller propagates the error. If you would rather
   have the fixture-level test as well, say so — it is a different test, not a replacement.
4. **The plan's counts were wrong twice and I re-measured rather than repeating them.** The Reproduce
   command prints 15 at HEAD, not the 14 the request and the Pre-Flight both state; the missing one is
   `repairvalidation/already_green.go`'s `uniqueSorted`, which the plan's own Exploration does name.
   Also, the plan says finalization's wrapper has 23 non-test call sites of which 22 are pure dedupe;
   the true shape is 22 production call sites (21 pure dedupe + 1 guard) plus 1 test call site — the
   plan counted the definition line. Nothing turns on either number; both are stated so a reviewer
   re-running the commands is not surprised.

## Found and NOT fixed

1. **`knowledgecommands.physicalPath` returns a relative path for a relative existing input.** It ends
   at `filepath.Clean` with no `filepath.Abs`, unlike `suiteinstall`'s resolver. The plan's Exploration
   records this as a discovered task rather than a fix, and D-04's scope is the name collision only. All
   three of its call sites pass `filepath.Clean(...)` of an already-absolute path today, so it does not
   bite. Left alone. Worth a follow-up REQ.
2. **`knowledgecommands.semverMajor` is a third, even looser version parser** sitting nine lines from
   the `compareSemver` that was just deleted: `strings.Split` on ".", `len(parts) != 3`, `Atoi(parts[0])`.
   It accepts `1.x.y`. It is outside this request's six names, so it stays. Folding it onto
   `ParseSemanticVersion` would be a small, clean follow-up.
3. **`suiteinstall.compareSemanticVersions`** (private, lowercase) is a fourth semver comparator, with
   its own orientation ("1 when remote is newer"). Not one of the six names and not in the write set. It
   is a genuine candidate for `sharedprimitives.CompareSemanticVersions`, but merging it means reading
   its callers and its `SplitN(…, 3)` leniency, which is more than this request authorises.
4. **`gittransaction.stringSet` and `corehelpers.stringSet` are two definitions of the same helper**
   with different value types (`map[string]struct{}` vs `map[string]bool`). Same duplication class as
   this request's finding, different name, so not covered by the lock-in as written. Noted, not touched.
5. **`nextselection.numericID` and `repositorymodel.requestNumberFromText` remain two request-number
   parsers with different strictness.** D-05 deliberately keeps both — the strict one validates CLI
   target tokens, the permissive one assigns identity. That is a decision, not an omission, but a
   reviewer scanning for "one parser per concept" will find two.

## What a reviewer should look at hardest

The release guard rewrite in `publication/release.go`. The two `compareSemver` bodies were inverted and
no test pinned the orientation, so the predicate and the comparator had to flip together in one edit.
The refused input set is unchanged — `!versionsParsed || versionOrdering >= 0` refuses exactly what
`compareSemver(old, new) != 1` refused, including unparseable input, which the old code scored as 0 —
but that is the one place where a reading error would have been silent.
