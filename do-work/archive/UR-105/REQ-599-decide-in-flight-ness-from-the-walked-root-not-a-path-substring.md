---
id: REQ-599
status: completed
domain: backend
created_at: 2026-09-06T06:31:11Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
route: A
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-06T07:26:22Z
  basis:
    - Route A
    - 1-file write set plus its test
    - 3 acceptance criteria
maintenance: false
depends_on: []
related: [REQ-596]
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go]
title: 'Decide a REQ in-flight from the root being walked, not from a substring of its absolute path'
claimed_at: 2026-09-06T07:26:22Z
completed_at: 2026-09-06T08:06:01Z
commit: 75da24a1fdafca824a4a8dba51aebcc0ed08db46
release_at: 2026-09-06T08:06:01Z
---

# Decide In-Flight-ness From the Walked Root, Not a Path Substring

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`. Approach: write the
  failing test first against a checkout whose absolute path contains a `working` component, then carry
  an active flag per walked root into the callback and delete the substring test.
- [x] **[APPLY]:** Two files, both in the one package: `inventory.go` (+15/-6) and `inventory_test.go`
  (+28). Nothing else.
- [x] **[UNIFY]:** `git diff --stat` over the merge — 2 files, 43 insertions, 6 deletions. `go build`,
  `go vet` and `gofmt -l` clean. The package suite green and unchanged; the new test red on the old code
  and green on the new. No debug artifacts.

## What

`AssociateProjectPaths` decides whether a REQ is in-flight with
`active := strings.Contains(filepath.ToSlash(path), "/working/")` on the **absolute** walked path
(`internal/corehelpers/inventory.go:281`). The loop already knows which of the two roots it is walking;
it should ask that instead.

## Why

A repository checked out anywhere beneath a directory named `working` — a common name for a scratch or
workspace directory — makes every archived REQ satisfy the test. The terminal-success status filter on
the next line is then skipped, so a blocked, cancelled or abandoned archived REQ claims paths it should
not, and the wrong REQ is named as a file's owner during a commit.

## Context

Found during the independent three-lens review of REQ-596, which corrected the guide's description of
this rule. The prose is right and the tool is wrong, which is why it is its own request rather than part
of that one. One reviewer raised it; the synthesizer reproduced it.

## Detailed Requirements

- Decide in-flight-ness from the root the walk is currently in, not from a substring of the path.
- The observable rule must not change for a repository checked out anywhere else: a `working/` REQ
  counts whatever its status says, an `archive/` REQ only on a terminal-success alias.
- Add a test that fails on the current code: a fixture repository whose absolute path contains a
  `working` component, holding an archived REQ with a non-terminal status that claims a path. Today it
  wins the path; after the change it does not.

## Constraints

- One file plus its test. No prose change: the guide already describes the intended rule correctly.
- The package's existing inventory and association tests must stay green unchanged.

## Red-Green Proof

**RED case:** a checkout under a directory named `working`, an archived REQ with `status: blocked`
naming a path in its Implementation Summary, and that path uncommitted.
**Why RED now:** the substring test marks the archived REQ active, the status filter is skipped, and the
blocked REQ is reported as the path's owner.
**GREEN when:** the same fixture reports the path unassociated, and the existing tests are unchanged.

## Open Questions

None.

## Triage

**Route: A** — Build directly.

**Reasoning:** The defect is one expression at a known line, the intended rule is already stated
correctly in the guide, and the reproduction is a fixture whose absolute path contains a `working`
component. There is nothing to discover: the walk loop already knows which of its two roots it is in,
and the fix is to ask that instead of the path. The test is the point as much as the fix — no test in the
package reaches the archived-with-non-terminal-status case today, which is why the substring could sit
there.

**Planning:** Skipped.

## Plan

**Planning not required** — Route A: one expression, one test that fails on the current code.

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go`

**What was done:** `AssociateProjectPaths` no longer decides whether a REQ is in flight by searching its
absolute path for `/working/`. The two roots it walks are now a slice of directory-plus-active pairs —
`do-work/working` active, `do-work/archive` not — and the walk callback reads the pair it is walking.
A slice rather than a map keeps the original walk order, so the tie-break between roots is unchanged.
The substring line is deleted and a four-line comment states the rule.

**The test came first and failed on the old code.**
`TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest` builds a repository at
`<tmp>/working/project`, an archived REQ with `status: blocked` whose Implementation Summary names an
uncommitted file, and asserts the file is not associated. On the old code:
`blocked archived request claimed project.txt through the checkout path: owner="REQ-905"`. It also
guards that the fixture path really contains `/working/`, and pins the other half of the rule in the
same checkout — a blocked `working/` REQ still claims its path.

Merge range `dc5d8180..890ed0c8`, two files, 43 insertions and 6 deletions.

## Decisions — implementation

- **D-01 — the roots carry the flag; the path does not. DECIDE & STATE.** The loop already iterated two
  known directories. Reading the flag off the pair is the only design in which a path's spelling cannot
  change the answer.
- **D-02 — a slice, not a map, so walk order is preserved.** The tie-break between two REQs claiming one
  path is `completed.After(current.completed)`, and which claim is seen first decides equal cases. A map
  would have made that order unspecified.
- **D-03 — the test pins both halves in one checkout.** Asserting only that the archived blocked REQ
  loses would pass under a fix that made *everything* inactive. The same test asserts a blocked
  `working/` REQ in the same nested checkout still claims its path.

## Discovered Tasks

- **The guide's tie-break sentence is wrong for a missing timestamp, and it is a sentence this run
  wrote in REQ-596.** It says "when the timestamps are equal or missing the first claim found stands".
  `ParseTimestamp` yields the zero time for a missing `completed_at` and the comparison is
  `completed.After(current.completed)`, so a `working/` REQ with no timestamp — every in-flight REQ —
  loses the path to any `archive/` REQ that has one, whatever the walk order. ~~Only when **both** are
  missing does the first found stand.~~ **Corrected after review:** the first claim found stands whenever
  the two timestamps compare equal, present or missing, because the comparison is a strict `After`; and
  `ParseTimestamp`'s error is discarded at the call, so an unparseable `completed_at` counts as missing.
  The sentence REQ-601 should write: a REQ with a parseable `completed_at` beats one without, whichever
  root is read first; when both parse equal, or neither parses, the first claim found stands (`working/`
  before `archive/`, name order within a root). Folded into REQ-601, whose write set now includes the guide.

## Qualification

**Passed.** Read from the merge range `dc5d8180..890ed0c8`, two files, 43 insertions and 6 deletions.
Canonical `qualify` satisfied.

- **The test was red before the fix existed, and the failure names the defect.** On the unchanged code:
  `blocked archived request claimed project.txt through the checkout path: owner="REQ-905"`. That is the
  exact behaviour REQ-596's review reproduced by hand.
- **The fix is the smallest one that makes the path's spelling irrelevant.** The roots were already two
  known directories iterated in order; carrying an active flag on each and reading it in the callback
  removes the only place the path was consulted. Nothing else in the walk changed.
- **Walk order is preserved deliberately**, because the tie-break between two claims on one path
  depends on which is seen first. A slice keeps it; a map would not.
- **The normal-checkout rule is unchanged, and three existing tests already pin it**: an archived
  `completed` REQ claims (`commands_test.go`), a `claimed` `working/` REQ claims
  (`inventory_test.go` via `writeInventoryOwner`), and the alias list rejects `cancelled`
  (`TestTerminalSuccessAliases`). The new test adds the case none of them reached: a `blocked` `working/`
  REQ inside a nested checkout still claims.
- **One thing found on the way is not fixed here, and it is prose this run wrote.** The builder noticed,
  while choosing a slice over a map, that a `working/` REQ with no `completed_at` loses a contested path
  to any archived REQ with one, whatever the order — because a missing timestamp parses to zero. REQ-596
  wrote the opposite into the guide. It is folded into REQ-601 rather than fixed here, because this
  request's write set is one Go file and its test.

### Remediation qualification (after review)

**Passed by inspection of the range `98deb1c..75da24a`**, two files, 7 insertions and 8 deletions,
both already named in the Implementation Summary; the canonical gate cannot be re-invoked at the review
phase, and the earlier `qualify` and `scope-drift` over `dc5d818..890ed0c` stand.

- **The test now pins the rule, not the bug.** The fixture sits under `<tmp>/do-work/working/project`,
  so a narrower substring test on the absolute path (`"/do-work/working/"`) fails it where the walked-root
  flag passes. Before, that mutation passed: the test defeated the exact substring the bug used and no
  other. The fixture-path guard, true by construction, and the `git init` the walk never runs are gone.
- **The dead field is deleted, not re-plumbed.** `claims` stored an `active` value nothing read; the
  original change rewrote the literal to carry `walkedRoot.active` and the record called that the
  smallest fix. It was not. The value is `{id, completed}` now.
- **The comment says why a slice.** Walk order decides equal-timestamp ties, so `working/` before
  `archive/` is deliberate and a map would make the tie-break unspecified.
- **One output change the record did not name, now named.** A blocked REQ stored under
  `do-work/archive/working/` in an ordinary checkout was treated as in flight by the old substring test and
  is filtered now. That is the same defect reached by a second route and the stated rule's answer; it is
  not a regression, and the Detailed Requirements' "no observable change for a repository checked out
  anywhere else" was one route too narrow.

## Testing

**Red, then green, on the same test.**

- Red, `inventory.go` untouched:
  `--- FAIL: TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest` —
  `blocked archived request claimed project.txt through the checkout path: owner="REQ-905"`, exit 1.
- Green after the fix: `--- PASS`, `ok .../internal/corehelpers 0.007s`.
- Whole package: `ok .../internal/corehelpers 4.169s`, exit 0, against a 6.482s baseline before any
  change. `go build ./...` and `go vet ./...` exit 0; `gofmt -l` on both files prints nothing.

**Gate at the builder's head:** `Maintainer verification passed.`, exit 0, wall 86s, CLI module at
**797** tests — one more than the 796 at the branch point, which is this test.

### Remediation testing (after review)

- `go test -count=1 ./internal/corehelpers/` at `75da24a` — `ok`, exit 0; `gofmt -l` prints nothing;
  `go build ./...` and `go vet ./internal/corehelpers/` exit 0.
- **The narrower-substring mutation is red against the deepened fixture**: with the flag test replaced by
  `!strings.Contains(filepath.ToSlash(path), "/do-work/working/")`, the test fails at
  `inventory_test.go:394` with `owner="REQ-905"`. Against the old fixture the same mutation passed.
- **The two flag mutations stay red**: every root inactive fails at `:397` with `owner=""`; every root
  active fails at `:394` with `owner="REQ-905"`. Each restored from a copy and the restore diffed.
- Fast gate at `75da24a`: `Maintainer verification passed.`, exit 0, gate wall 101s.

## Review

**Overall: 88%** | 2026-09-06T09:10:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 85% |
| Test Adequacy | 78% |
| Scope | 98% |
| Risk | Low |
| Acceptance | Accept; tidy-ups applied before archive |

**Verdict: Accept.** Three independent reviewers (fix-and-test, behaviour preservation by differential
testing of the parent and merge binaries, quality-and-record) and a synthesizer who reproduced every
finding. The substring test is gone, the roots carry the flag, the new test is red on the parent and green
at the merge, and every mechanical claim in the record reproduced: build, vet, gofmt, package green, diff
stat, the test count by the gate's counter, canonical qualify. The seven findings that survived are all
low or informational, and all seven are applied in the remediation above rather than deferred.

**Findings that survived verification:**

- **F1 (low, applied)** The test did not pin "never from the path": replacing the flag with a narrower
  `"/do-work/working/"` substring kept it and the whole package green. Fixture deepened; the mutation is
  red now.
- **F2 (low, applied)** `claims.active` was written and never read, and the change re-plumbed it rather
  than deleting it. Deleted.
- **F3 (low, recorded)** "No observable change elsewhere" missed the `do-work/archive/working/` route,
  where an archived blocked REQ flips from owner to unowned. Rule-consistent; named in the remediation
  qualification.
- **F4 (low, applied)** The Discovered Task's tie-break sentence said first-found stands only when both
  timestamps are missing; it stands on any equal comparison, and an unparseable timestamp counts as
  missing. Corrected in place with the sentence REQ-601 should write.
- **F5 (low, applied)** The test ran `git init`, which the walk never needs. Deleted.
- **F6 (info, applied)** The comment stated the in-flight rule but not why the roots must stay a slice.
  One sentence added.
- **F7 (info, applied)** The fixture-path guard was true by construction. Deleted with the fixture change.

**Findings that did not survive:** "the record does not say this is a release" (that judgment belongs to
finalization, which had not run); "gate exit 0 could not be reproduced" (two `heavyverification` tests
fail in the reviewers' sandbox at the parent as well, so environmental, and the builder's hand-back
records exit 0 with the heavy probes skipped); "797 tests not reproduced" (`go test -list` gives 811 and
the gate's counter 797; the record uses the gate's number and both show the same +1); one reviewer's
`QUALIFY-SUMMARY-MISSING` came from a clone whose record predated the Implementation Summary.

**Disagreements and how settled:** whether the dead field was a finding (settled by removing it: builds,
vets and passes, 1 insertion and 3 deletions); whether the test pinned the design claim (settled by
running the narrower mutation: it passed); whether the `archive/working/` flip was a finding (settled by a
scratch test on parent and merge sources: owner on the parent, unowned at the merge).

## Lessons Learned

- **A test that defeats the bug you found does not pin the rule you stated.** The fixture mirrored the
  exact substring the bug used, so every narrower substring passed it. Pin the rule by making the fixture
  hostile to the whole family: put the forbidden token where the most specific wrong fix would still see it.
- **"Smallest fix" is checked field by field.** The change touched a struct and rewrote a literal to keep
  carrying a value nothing reads. When a change touches a struct, ask of each field whether anything reads
  it; a field that survives a change only because it was there before is the drift the delete-before-add
  rule exists for.
- **Say what a fix changes, not only what it preserves.** The old bug also fired for a REQ stored under
  `do-work/archive/working/`, so the fix changes output there. The record said "no observable change
  elsewhere" because the author had one route in mind. A differential run of the two binaries found the
  other one in minutes.
- **A correction of a rule needs the code's whole condition, not the case that was noticed.** The
  tie-break sentence written here to correct REQ-596's was itself incomplete: it named the missing
  timestamp case and missed equal present timestamps and unparseable ones, which the same strict `After`
  decides the same way. Read the comparison, then write the sentence.

## Orientation

`AssociateProjectPaths` in `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` decides
which REQ owns a changed project path. It walks two roots in a fixed order, `do-work/working` then
`do-work/archive`, each carried as a `{directory, active}` pair in a slice, and the walk callback reads
the pair. **In-flight-ness comes from the root being walked, never from the absolute path.** A `working/`
REQ counts whatever its status says; an `archive/` REQ counts only on a terminal-success alias
(`terminalSuccessStatus`). A REQ file stored under `do-work/archive/working/` is an archive REQ.

**Ties.** The winner is `!exists || completed.After(current.completed)`, a strict comparison, so a REQ
with a later parseable `completed_at` wins whichever root is read first, and on an equal comparison the
claim seen first stands. `ParseTimestamp`'s error is discarded at the call: a missing or unparseable
`completed_at` is the zero time. That is why the roots must stay a slice and why `working/` is read
first; the comment above `walkedRoots` says so.

**The pin.** `TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest` builds its
checkout under `<tmp>/do-work/working/project`, so a blocked archived REQ must be skipped and a blocked
working REQ must claim. That fixture is hostile to every path-substring fix, not only the one that shipped:
reverting to any substring test on the absolute path makes it red.
