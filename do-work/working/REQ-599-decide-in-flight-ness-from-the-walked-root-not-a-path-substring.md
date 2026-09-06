---
id: REQ-599
status: claimed
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
  loses the path to any `archive/` REQ that has one, whatever the walk order. Only when **both** are
  missing does the first found stand. Folded into REQ-601, whose write set now includes the guide.

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
