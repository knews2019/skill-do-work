---
id: REQ-557
title: '[impact-negligible] Deduplicate six Go helper names defined fourteen times across do-work-cli'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-09-06T03:18:43Z
  basis:
    - Route C
    - 10-file write set
    - 8 subsystems involved
    - 4 acceptance criteria
depends_on: [REQ-550, REQ-552]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
route: C
planning_at: 2026-09-06T04:01:02Z
write_set: [skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go, skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives_test.go, skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go, skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go, skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go, skills/do-work/tools/do-work-cli/internal/publication/capture_files.go, skills/do-work/tools/do-work-cli/internal/publication/release.go, skills/do-work/tools/do-work-cli/internal/publication/release_mirrors.go, skills/do-work/tools/do-work-cli/internal/publication/answer.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T03:16:29Z
---

# Deduplicate six Go helper names defined fourteen times across do-work-cli

## What
`uniqueSorted`, `subtractPaths`, `requestIDLess`, `firstError`, `compareSemver` and `physicalPath` are defined fourteen times across `internal/`; in every case the duplicating package already imports the package holding an earlier copy, and three names have copies that disagree. Export one canonical definition per name in the lowest already-imported package (`corehelpers` for the path and error helpers, `repositorymodel` for `requestIDLess`), delete the other eight, and record the three semantic reconciliations as named decisions in this REQ.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The validated plan is in `## Plan` below: one new leaf package `internal/sharedprimitives`
  holding four exported helpers plus an exported strict semver parser, `RequestIDLess` exported in place
  in `repositorymodel`, and two path resolvers deliberately kept apart under two names. Ten ordered
  steps, the tree compiling after each.
- [x] **[APPLY]:** All ten steps applied. One path was added to the write set — the D-02 guard test —
  and it is declared in the frontmatter and in `## Scope` above. One number in the plan was wrong and
  the deviation is D-06 below.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file: 19 modified, 3 added. `go build ./...`,
  `go vet ./...` clean; `gofmt -l` over the module empty; `go test ./...` green at 808 top-level tests
  (794 pass, 14 skip, 0 fail) against 784 before. `shellcheck -S warning` clean on the lock-in file, and
  the lock-in was proved red in three directions before being accepted. No debug artifacts, no
  commented-out code, no leftover probes in the diff.

## Why
Per-REQ helper files that duplicate an existing helper are the agent-creep class; the three semantic splits (`uniqueSorted` drops empty strings in one copy, `compareSemver` accepts unparseable input in one copy, `physicalPath` has two contracts) are latent correctness drift.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 4, sweep_key `per-req-duplicate-go-helpers`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -70. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` — `subtractPaths` and `uniqueSorted` duplicate `internal/corehelpers/checks.go` (introduced 761d8e6a, REQ-498; `finalization` already imports `corehelpers`).
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` — third `uniqueSorted`, silently drops empty strings (01d920dd).
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` and `internal/repositorymodel/repository_model.go` — two `requestIDLess` bodies in one commit (ac2e3acd, REQ-408) with two different number parsers; `internal/nextselection/next_types.go` — third `requestIDLess` (625d49aa, REQ-411; `nextselection` already imports `repositorymodel`).
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` — `firstError` byte-identical to `corehelpers/checks.go` (cf111a50, REQ-413).
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` — `compareSemver` returns 0 for unparseable input while `internal/publication/release.go` rejects it.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` — `physicalPath` is `EvalSymlinks`+`Abs`; `internal/knowledgecommands/commands.go` walks missing ancestors; same name, different contract.
- Each reconciliation (empty-string handling, unparseable semver, missing-ancestor paths) is written into this REQ's Implementation Summary as a decision with the behaviour each caller keeps; a silent pick is a review refusal.
- Reproduce at dc8a64e3 (prints 14 lines): `rg -n --glob '*.go' --glob '!*_test.go' '^func (uniqueSorted|subtractPaths|requestIDLess|firstError|compareSemver|physicalPath)\(' skills/do-work/tools/do-work-cli/internal/ | sort`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- No import cycles introduced: the canonical home is always a package the duplicator already imports.
- Tests unchanged except where a test named a deleted private helper directly; then it points at the canonical one.
- Lock-in limit: definitions of the six helper names: 6 after this REQ (today 14).

## Dependencies
Depends on REQ-550 and REQ-552 so `corehelpers` is settled before helpers move in. REQ-558 depends on this REQ.

## Builder Guidance
Firm on one definition per name and on recording the three reconciliations; latitude on the exported names.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints fourteen definitions of six names.
**GREEN when:** It prints six lines, one per name; `go test ./...` green; the lock-in pins definitions of these six names at 6.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for per-req-duplicate-go-helpers.*

## Triage

**Route: C** — Explore, plan, then build.

**Reasoning:** Three of the six names have copies that disagree with each other, and the request says a
silent pick is a review refusal. Deciding what each caller keeps for empty-string handling in
`uniqueSorted`, for unparseable input in `compareSemver`, and for the two different `physicalPath`
contracts is a behaviour decision per call site, not a rename. That needs every caller read before any
definition is deleted, so exploration is real work and the plan it produces is what the builder
executes. The write set is nine Go files plus one lock-in, and `corehelpers` has been edited twice
already in this run, so the copies the request names must be re-checked against HEAD rather than
against the audited commit it cites.

**Planning:** Required. The exploration produces the per-name canonical home and the per-caller
behaviour decision; the plan turns that into an ordered edit sequence that never leaves the tree
uncompilable between steps.

## Plan

Two independent plans were produced from one three-slice survey; the judge agent returned an empty
synthesis, so the orchestrator judged them and re-verified the load-bearing claims directly.

**The request's proposed home is wrong, and this was measured, not argued.**
`go list -deps` at HEAD: `internal/corehelpers` is **not** in the dependency closure of
`knowledgecommands`, `repairvalidation`, `publication` or `suiteinstall` — four of the five duplicating
packages. Routing the helpers there would violate the request's own "a package the duplicator already
imports" rule and add 6-9 transitive packages to each. `internal/resultmodel` **is** in all their
closures and is itself a leaf.

**D-A — the canonical home is a new leaf package, `internal/sharedprimitives`, not `resultmodel`.
DECIDE & STATE.** `resultmodel` would cost zero new packages and zero new import lines, and it is the
result-and-finding model: a 1333-line file that would gain a semantic-version comparator and a
string-set helper. The reason these six names were duplicated fourteen times is that there was nowhere
obvious to put them; putting them in the result model leaves that still true for the seventh. A leaf
package with zero internal imports satisfies the no-cycle purpose of the request's rule absolutely,
even though it does not satisfy its letter. That deviation is the reason this is a stated decision.

**The survey corrected the request's premise in four places, each verified:**

- There is a **fourth** `uniqueSorted`, in `internal/repairvalidation/already_green.go:442`. The request
  names three.
- `finalization`'s `uniqueSorted` is **not a dedupe helper**. Its body is
  `result, _ := normalizeRepositoryPaths(paths); return result` — a validator whose refusal is thrown
  away, returning `nil` for the whole slice when any one element is rejected.
- The two `compareSemver` copies are **inverted**, not merely differently strict. The request describes
  only the strictness. `knowledgecommands` returns -1 when the first argument is older; `publication`
  returns +1 when the second is newer.
- `physicalPath` differs on **two** axes, not one. Beside the named absence-is-an-error axis, the
  `knowledgecommands` copy ends at `filepath.Clean` with no `filepath.Abs`, so it returns a relative
  path for a relative existing input.

### Canonical homes

| Name | Canonical | Home |
|---|---|---|
| `uniqueSorted` | `UniqueSortedStrings(values []string) []string` | `internal/sharedprimitives` |
| `subtractPaths` | `SubtractStringValues(leftValues, rightValues []string) []string` | `internal/sharedprimitives` |
| `firstError` | `FirstNonNilError(firstCandidate, secondCandidate error) error` | `internal/sharedprimitives` |
| `compareSemver` | `CompareSemanticVersions(leftVersion, rightVersion string) (int, bool)` | `internal/sharedprimitives` |
| `requestIDLess` | `RequestIDLess` exported in place | `internal/repositorymodel` |
| `physicalPath` | not merged — see D-04 | stays private in two packages |

### Decisions — the reconciliations the request requires

- **D-01 — `UniqueSortedStrings` keeps every value it is given, including the empty string.** Two of
  the four copies filter empties. Every producer at every one of the 42 call sites was read and none can
  emit one — the corehelpers callers pre-filter at `checks.go:779`, `:799`, `:897` and `:900`, and
  `repairvalidation` passes `strings.FieldsFunc` output and `record[3:]` of a `len >= 4` record. So the
  filter cannot fire today, and dropping it is the safe direction: if a producer ever does emit a blank,
  it becomes visible in an evidence list instead of vanishing. The parameter is `values`, not `paths`,
  because `already_green.go:283` passes reason codes.
- **D-02 — `finalization`'s wrapper is deleted, and the one call site whose result feeds a guard gets
  its error back.** 22 of its 23 non-test call sites want plain dedupe and move to
  `SubtractStringValues`/`UniqueSortedStrings`. `finalization_prepare.go:162` is the exception: today a
  required commit path that is empty, absolute or escaping makes the wrapper return `nil` for the whole
  slice, so `subtractPaths(nil, …)` is empty, `len(missing)` is 0, and the "commit_paths omits planned
  lifecycle or release targets" error at line 164 **never fires**. That site calls
  `normalizeRepositoryPaths` directly and propagates its error, exactly as line 153 already does nine
  lines above it in the same function. **This is a behaviour change and it is deliberate**: carrying a
  silently disabled guard forward under a new name would be shipping the drift this request exists to
  remove. It needs a test that fails when the guard is disabled again.
- **D-03 — one orientation, one strict parser, and an explicit `parsed` flag.**
  `CompareSemanticVersions` uses the standard Go orientation (negative when the first argument is
  older) and `publication`'s strict rules verbatim: exactly three dot-separated parts, no empty part,
  no leading zero on a multi-character part, integer only, non-negative. It returns `(ordering, parsed)`
  so no caller can read "unparseable" as "equal" — which is exactly the latent drift the request names.
  `publication/release.go:26` becomes `ordering, parsed := …; if !parsed || ordering >= 0 { refuse }`,
  the same input set it refuses today. `knowledgecommands` keeps its `< 0` predicate and now surfaces an
  error for a version its lenient parser used to score as zero. Export `ParseSemanticVersion` too, so
  `publication/release_mirrors.go:109` keeps working when the private `parseSemver` is deleted.
- **D-04 — `physicalPath` is not merged. The two contracts keep their behaviour under two names.**
  `suiteinstall/update_transaction.go` renames its copy to `existingPhysicalPath` (EvalSymlinks + Abs,
  absence is an error, result always absolute); `knowledgecommands/commands.go` keeps `physicalPath`
  (walks missing ancestors, absence succeeds). Merging is possible but costs an explicit existence check
  re-added at `update_transaction.go:179`, where the missing-path error **is** the existence check for
  the installed skill and feeds a `strings.HasPrefix` containment test. The competing plan called that
  check "the whole risk" of its own approach. The defect the request names is one name carrying two
  contracts; two names resolves it at zero behaviour change and four line touches.
- **D-05 — `RequestIDLess` uses the permissive parser.** `repositorymodel`'s
  `requestNumberFromText` (`(?i)^REQ-0*([0-9]+)`, anchored at the start only) is the parser that
  **assigns** identity elsewhere in the same file, so a comparator that disagreed with it would order
  ids by one rule and name them by another. `dependencygraph`'s hand-rolled byte-scan parser agrees with
  it on 8,016 inputs and 160,000 ordered pairs, so deleting it is behaviour-preserving.
  `nextselection`'s strict `numericID` is **kept** — it has nine other callers and is the right
  validator for CLI target tokens and the `UR-` prefix — but its `requestIDLess` goes.

### Ordered steps, each leaving the tree compiling

1. Create `internal/sharedprimitives` with the four exported helpers and `ParseSemanticVersion`, plus
   its own unit tests. Nothing calls it yet.
2. Export `repositorymodel.RequestIDLess` in place; keep `requestNumberFromText` private.
3. Point `dependencygraph` and `nextselection` at it; delete their copies and `dependencygraph`'s
   now-callerless `requestNumber`.
4. Move `corehelpers`, `knowledgecommands` and `repairvalidation` off their `uniqueSorted` copies;
   delete all three. Keep `corehelpers.stringSet` — it has two other callers.
5. Move `corehelpers` and `finalization` off `subtractPaths`; delete both copies.
6. Move `corehelpers` and `publication` off `firstError`; delete both copies.
7. Rewrite both `compareSemver` call sites per D-03; delete both copies and `publication`'s private
   `parseSemver`.
8. Delete `finalization`'s `uniqueSorted` wrapper; route its 22 pure-dedupe sites to the shared helper
   and give `finalization_prepare.go:162` the error back, with a test that fails when the guard is
   disabled again.
9. Rename `suiteinstall`'s `physicalPath` to `existingPhysicalPath`.
10. Add the lock-in.

### The lock-in

One block in `_dev/tests/audit-lockins.sh`, in the file's existing per-Finding shape. It counts
definitions of the union of the old and the new names and pins the total at ~~**exactly 6**~~
**exactly 7** — a floor as well as a ceiling, because a one-sided ratchet guards half the property
(REQ-556's lesson). **Corrected by the builder (D-06):** the 6 below counts `ResolvePhysicalPath`,
which is the merged resolver D-04 decided *not* to create, and omits `existingPhysicalPath`, the name
D-04 actually creates. Pinning 6 would have left the new name unguarded and shipped a stale enumeration
into the very edit that creates it. The shipped regex names `existingPhysicalPath` and pins 7:

```
rg -n --glob '*.go' --glob '!*_test.go' \
  '^func (uniqueSorted|UniqueSortedStrings|subtractPaths|SubtractStringValues|requestIDLess|RequestIDLess|firstError|FirstNonNilError|compareSemver|CompareSemanticVersions|physicalPath|existingPhysicalPath)\(' \
  skills/do-work/tools/do-work-cli/internal/
```

The union is used, not the six original names, because four of them cease to exist when the canonical
definitions are exported: the request's own reproduce command would print 1 rather than 6 after a
correct change. `rg`'s exit status is read rather than a piped total, and a status above 1 fails.

*Validated by orchestrator against the survey; the two competing plans and the survey are in the run
directory.*

## Exploration

Three read-only agents, one per name group, re-verified against HEAD rather than against the audited
commit the request cites (`dc8a64e3`, which is not an object in this clone). Full reports in the run
directory as `REQ-557-survey.json`; the two competing plans built from them are
`REQ-557-competing-plans.json`.

**Fourteen definitions is the wrong count and the wrong shape.** There are four `uniqueSorted` bodies,
not three — the request misses `internal/repairvalidation/already_green.go:442` — and one of the four,
`finalization`'s, is not a dedupe helper at all. Its body discards a validator's error and returns
`nil` for the whole slice when any single element is rejected, which at `finalization_prepare.go:162`
disables the guard on line 164 entirely.

**Both `compareSemver` copies were read side by side and they are inverted**, which the request never
mentions: `knowledgecommands/interview_commands.go:1535` returns -1 when the first argument is older,
`publication/release.go:210` returns +1 when the second is newer. A merge that picked one body and left
both predicates alone would have reversed `publication`'s release guard while every test stayed green,
because no test pins the orientation.

**`physicalPath` differs on two axes.** Beside the absence-is-an-error axis the request names, the
`knowledgecommands` copy ends at `filepath.Clean` with no `filepath.Abs` and so returns a relative path
for a relative existing input. That is recorded as a discovered task rather than fixed here.

**The import closures were measured, not assumed.** `go list -deps` at HEAD puts `corehelpers` outside
the closure of `knowledgecommands`, `repairvalidation`, `publication` and `suiteinstall`, and puts
`resultmodel` inside all of them with a closure of only itself. `dependencygraph` and `nextselection`
both already import `repositorymodel`, and `repositorymodel`'s own closure is `atomicfile`,
`requestmodel` and `schemanormalization`, so it is a leaf and the request's home for `requestIDLess`
is correct.

**Two parsers were compared by execution, not by reading.** `dependencygraph`'s hand-rolled byte-scan
request-number parser and `repositorymodel`'s regexp agree on 8,016 inputs and produce identical
orderings across 160,000 ordered pairs, which is what makes deleting the former behaviour-preserving.

*Generated by three Explore agents, judged by the orchestrator*

## Scope

Every path below is repository-relative, matching the `## Implementation Summary` list, because the
scope-drift comparator reads the backticked paths in both sections and compares them directly.

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go` (new) — the four exported helpers plus the exported semantic-version parser
- `skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives_test.go` (new) — unit tests for each, including the parsed-flag contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` (modify) — delete three private helpers, point their call sites at the shared package, keep the local string-set helper which has two other callers
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` (modify) — two dedupe call sites
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modify) — export the request-id comparator in place, keep its parser private
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` (modify) — delete the duplicate comparator and its now-callerless number parser, use the exported one
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go` (modify) — delete the third comparator; the strict validator beside it stays, it has nine other callers
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (modify) — call site
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modify) — delete the wrapper whose validator error is discarded; the one site whose result feeds a guard calls the validator directly and propagates its error
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modify) — eighteen call sites
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modify) — three call sites
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modify) — one reference to a deleted private helper
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go` (new, added by the builder) — the D-02 guard test, in the per-REQ test-file form this package already uses
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go` (modify) — delete two private helpers, rewrite the version predicate
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go` (modify) — two call sites
- `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go` (modify) — delete the fourth dedupe helper the request does not name
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` (modify) — delete the byte-identical error helper
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (modify) — delete the inverted comparator and the private strict parser, rewrite the release guard
- `skills/do-work/tools/do-work-cli/internal/publication/release_mirrors.go` (modify) — use the exported parser
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — one call site
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` (modify) — rename the path resolver, no behaviour change
- `_dev/tests/audit-lockins.sh` (modify) — one assertion block

**The declared set is larger than the write_set the request carried, and the widening is stated here
rather than discovered by a reviewer.** The request listed ten paths; it named the files holding the
definitions and not the files holding the call sites, and it missed the fourth dedupe helper entirely.
Most added paths are a call site of a deleted definition or a file the survey showed holds one.
**Corrected after review:** three are neither — the new canonical home
`internal/sharedprimitives/shared_primitives.go`, its tests, and the guard test D-02 requires. Those
exist because the plan chose a new leaf package and because one decision changed behaviour, not because
a definition was found there.

**Files I will NOT touch:** the string-set helper in corehelpers, which has two callers unrelated to
this class; the strict numeric validator in nextselection, which has nine other callers and is the
right validator for target tokens; any action file, prime or document. No shipped behaviour outside the
CLI module changes.

**Acceptance criteria:**
- [x] Exactly ~~six~~ **seven** definitions across the union of the old and the new names, and the
  lock-in fails in both directions. **Corrected after review**, to the number D-06 established and the
  shipped lock-in pins; the six was written before D-04 settled that the two path resolvers stay
  separate
- [x] Each of the reconciliations is a named decision stating the behaviour each caller keeps —
  delivered, though two of them stated their own extent wrongly and are corrected after review
- [x] No import cycle, and the shared package imports nothing else inside the module — delivered,
  `go list -deps` returns only itself and every package's module-internal closure changed by exactly
  `+sharedprimitives`
- [ ] **Not met.** Every current call site keeps its observable behaviour, except the one guard whose
  restoration is D-02 and which carries its own test. Two further call sites change observable
  behaviour: `nextselection`'s ordering when a request carries a non-canonical frontmatter `id:` (D-05),
  and seven version pairs in `knowledgecommands` that turn a stamp into a refusal (D-03). Both are
  accepted on their merits and both are now named where they belong; neither was named here when the
  criterion was written, which is why it stays unticked rather than being reworded to fit
- [x] Tests unchanged except where a test named a deleted private helper — delivered for the existing
  suite; `finalization_recovery_test.go` changes exactly one call and one import line

**One request constraint is superseded, and naming it is part of not quietly outgrowing it.** The
request says "no test files beyond the lock-in". The validated plan's step 1 requires unit tests for the
new shared package, D-02 requires a test for the restored guard, and the review requires two more. The
constraint was written for a pure move; this is a move plus three behaviour decisions, and each one that
carries a behaviour change carries a test.

## Pre-Flight

**Green gate at `03901f6d`**, the revision the builder branches from.
`bash _dev/tests/maintainer-verify.sh` printed `Maintainer verification passed.` and exited 0, gate
wall 76s, with the CLI module at 784 tests and the board module inside its 30s budget. One `SKIP` line,
the heavy-only one every fast run prints.

**The baseline this request's Red-Green Proof names is re-verified at HEAD, not replayed from the
audited commit.** ~~The reproduce command prints fourteen definitions of the six names, which is what
the request claims.~~ **Corrected by the builder:** it prints **fifteen** at the branch point. The
request's count of fourteen predates the fourth `uniqueSorted` the survey found in
`internal/repairvalidation/already_green.go`, which the audit missed. The shape is wrong as well as the
count: one of the fifteen is a validator wearing a deduper's name.

**The environment the gate needs is recorded once for this run.** `NODE_OPTIONS` and the
`GIT_CONFIG_COUNT` / `GIT_CONFIG_KEY_*` / `GIT_CONFIG_VALUE_*` triples unset, `GIT_CONFIG_GLOBAL`
pointed at a config carrying `commit.gpgsign = false`, and `QUEUE_KANBAN_BROWSER` at the container's
Chromium. A gate run refuses on an opaque runtime extension or an opaque Git configuration override,
and an unusable global signing key makes a fixture's own commit fail inside a lane.

**The builder works in an isolated worktree** at
`../skill-do-work-worktrees/worktree-agent-REQ-557-deduplicate-go-helper-names`, branched from
`03901f6d`, and hands back one file to the main checkout without staging or committing anything there.

## Implementation Summary

Every path below is repository-relative, because the qualification gate resolves the backticked paths
in this section against the repository root; the `## Scope` declaration writes the same set in its own
module-relative shorthand.

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go`
- `skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives_test.go`
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
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
- `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go`
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go`
- `skills/do-work/tools/do-work-cli/internal/publication/release.go`
- `skills/do-work/tools/do-work-cli/internal/publication/release_mirrors.go`
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go`
- `_dev/tests/audit-lockins.sh`

**What was done:** Fifteen definitions of six helper names became seven. ~~Nine~~ **Twelve** private copies were
deleted outright (corrected after review — the first count omitted three); one comparator was exported in place; one path resolver was renamed so that one name
no longer carries two contracts. The new `internal/sharedprimitives` imports nothing else in the module
(`go list -deps` shows only itself), so it can be called from any package without a cycle — which is
the reason the copies existed in the first place.

**The request's home was wrong and the plan measured why**, so the canonical helpers went to a new leaf
package rather than to `corehelpers`, which is outside the import closure of four of the five
duplicating packages. That deviation is D-A in `## Plan`.

**Two of the request's own numbers were off by one and both were re-measured.** The reproduce command
prints fifteen at the revision this was built on, not fourteen: the audit missed
`internal/repairvalidation/already_green.go`'s fourth `uniqueSorted`. The lock-in's pinned total is
seven, not six, for the reason in D-06.

## Decisions — implementation

- **D-01 — `UniqueSortedStrings` keeps the empty string.** Two of the four deleted copies filtered
  blanks (`knowledgecommands`, `repairvalidation`), two did not (`corehelpers`, `finalization`). No
  producer at any call site can emit one: the `corehelpers` sites pre-filter empty backtick spans, the
  `repairvalidation` sites pass `strings.FieldsFunc` output (which never yields an empty field) and
  `record[3:]` of a `len >= 4` record, and every `knowledgecommands` site passes map keys of a write set
  or `filepath.Join`/`filepath.Dir` output. The filter therefore cannot fire today, and dropping it is
  the safe direction: a blank that ever does appear becomes visible in an evidence list instead of
  vanishing from it. The result is always a non-nil empty slice, so JSON evidence still marshals `[]`
  and never `null`. ~~which every deleted copy also guaranteed~~ **False, corrected after review:**
  `finalization`'s copy returned `nil` whenever `normalizeRepositoryPaths` refused, and that `nil` is
  the exact defect D-02 fixes. A third `repairvalidation` producer shape the first version of this
  bullet missed: `requeststate` plan change paths, appended only under a non-empty guard.
- **D-02 — `finalization`'s wrapper is deleted and the guard it disabled is restored. This is a
  deliberate behaviour change.** The wrapper's whole body was `result, _ := normalizeRepositoryPaths(paths);
  return result`. Twenty-one production call sites and one test site wanted plain dedupe and now call
  `sharedprimitives.UniqueSortedStrings`; every one of them passes repository-relative slash paths that
  come from git output, journal images, request relative paths or manifest lock paths, so the discarded
  normalization was a no-op there. The twenty-second site is
  `finalization_prepare.go` — a required commit path that was empty, absolute or escaping turned the
  whole required set into `nil`, the subtraction then found nothing, and the "commit_paths omits planned
  lifecycle or release targets" error on the next line could never fire. That computation now lives in
  `missingCommitPaths`, which propagates the refusal exactly as the `normalizeRepositoryPaths` call nine
  lines above it already does. `internal/finalization/finalization_req557_test.go` fails if the error is
  discarded again.
- **D-03 — one orientation, one strict parser, an explicit `parsed` flag.** The two `compareSemver`
  bodies were inverted, not merely differently strict: `knowledgecommands` returned -1 when the first
  argument was older, `publication` returned +1 when the second was newer. `CompareSemanticVersions`
  uses the standard Go orientation (negative when the first argument is older) and `publication`'s
  strict rules verbatim. It returns `(ordering, parsed)` so no caller can read "unparseable" as "equal".
  `publication/release.go` refuses on `!parsed || ordering >= 0`, which is the same input set it refused
  before. `knowledgecommands` keeps its `< 0` predicate and now returns the section's own
  "template and session versions must be bare semver" error for a version its lenient parser used to
  score as zero — previously that version silently skipped the template-version stamp.
  **Corrected after review: that is one direction of two.** Seven version pairs also turn
  *stamp* into *error* — a pair the lenient comparator scored -1, meaning migrate and stamp, now aborts
  the interview command. `semverMajor` is unchanged and gates the branch on three dot-separated parts
  with an `Atoi`-able first part, so a leading zero, an empty middle part, a trailing dot or a
  non-numeric third part gets past it and used to be scored numerically. Session `1.0.0` against
  template `1.09.0` is one such pair. `template.Version` is read verbatim from a project-authored
  template's frontmatter, so this is project data rather than tool-written data. The remediation adds
  the test that pins it.
  `ParseSemanticVersion` is exported so `publication/release_mirrors.go` keeps its admission check when
  the private `parseSemver` goes.
- **D-04 — `physicalPath` is not merged; two contracts keep two names.** `suiteinstall`'s copy becomes
  `existingPhysicalPath` (EvalSymlinks + Abs, absence is an error, result always absolute);
  `knowledgecommands` keeps `physicalPath` (walks missing ancestors, absence succeeds). Merging would
  cost an explicit existence check re-added in `resolveUpdateRoots`, where the missing-path error IS the
  existence check for the installed skill root and feeds the `strings.HasPrefix` containment test. Two
  names resolve the defect the request names — one name carrying two contracts — at zero behaviour
  change.
- **D-05 — `RequestIDLess` uses the permissive parser.** `repositorymodel`'s `requestNumberFromText`
  is the parser that assigns request identity elsewhere in the same file, so a comparator that
  disagreed with it would order ids by one rule and name them by another. `dependencygraph`'s
  hand-rolled byte-scan parser agrees with it on every input the survey tried, so deleting it preserves
  behaviour. `nextselection`'s strict `numericID` stays — nine other callers, and it is the right
  validator for CLI target tokens — but its `requestIDLess` goes, so an id with trailing non-digits
  (`REQ-12x`) now sorts by 12 in `nextselection` where it previously fell to the string comparison.
  ~~Such an id cannot reach these two sorts: both sort ids drawn from a repository snapshot, whose
  filenames and `id:` fields are produced by `formatRequestID`.~~ **False, and the reviewer reached it
  end to end.** Only `FilenameID` goes through `requestIDFromText`. The frontmatter `id:` is raw
  `document.scalarValue("id")`, never normalized, and `nextselection.requestID()` *prefers* it over
  `FilenameID`. A filename/frontmatter mismatch only appends to `snapshot.WarningMessages`, which
  `nextselection` never reads and `next` never prints. Built from both revisions against a fixture
  carrying `id: REQ-12x`, the base binary selects `REQ-100` and the head binary selects `REQ-12x`. So
  in a repository holding a non-canonical `id:`, the queue order and the request `next` selects both
  change. **The change is accepted anyway**, because the permissive parser is the one that assigns
  identity elsewhere in the same file, and a comparator that disagreed with it would order ids by one
  rule and name them by another. What was wrong was the claim of unreachability, not the pick.
- **D-06 — the lock-in pins the union at seven, not the six the plan states. This is a deviation and it
  is arithmetic, not judgment.** The plan's own regex reaches six only by naming `ResolvePhysicalPath`,
  a merged resolver that D-04 decided NOT to create, and by omitting `existingPhysicalPath`, the name
  D-04 actually creates. A union of the old and new names that covers what the change really produced
  counts seven: five canonical helpers, plus the two deliberately separate path resolvers. Keeping the
  six would have shipped a ratchet that names an identifier which can never exist and leaves the one new
  name unguarded. The block says in place why the total is seven.

## Qualification

**Passed.** Read from the merge range `adfbaeeb..1b119056`, 23 files, 530 insertions and 231 deletions.
Canonical `qualify` and `scope-drift` both satisfied.

- **The builder corrected the plan three times rather than following it, and each correction is right.**
  The lock-in pins the union at **7**, not the 6 the plan wrote: the plan's regex named
  `ResolvePhysicalPath`, which is the merged resolver D-04 decided not to create, and omitted
  `existingPhysicalPath`, the name D-04 actually creates. Pinning 6 would have counted a phantom, left
  the new name unguarded, and shipped a stale enumeration into the edit that creates it — the exact trap
  `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale names. The request's own baseline
  is off by one for the same reason the survey found: the reproduce command prints **15** at the branch
  point, not 14, because the audit missed the fourth `uniqueSorted`. And the finalization wrapper's call
  sites are 21 pure-dedupe plus 1 guard plus 1 test, not 22 plus 1. All three are corrected in the plan
  above rather than left standing.
- **The one deliberate behaviour change is D-02 and it carries its own test.** The deleted wrapper's
  body discarded a validator's error, so a required commit path that was empty, absolute or escaping
  made it return `nil` for the whole slice, `subtractPaths(nil, …)` came back empty, and the
  `commit_paths omits planned lifecycle or release targets` error never fired. Disabling the restored
  propagation reds
  `TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet` with
  `missingCommitPaths accepted the unusable required path "" and returned missing=[]string{}; the
  commit_paths guard is disabled again`.
- **The inversion the request never mentions was the real risk, and it is closed.** The two
  `compareSemver` copies returned opposite signs for every unequal pair. A merge that picked one body
  and left both predicates alone would have reversed `publication`'s release guard with every test
  staying green, because nothing pinned the orientation. Both call sites were rewritten and both
  orientations are now pinned by name.
- **The new package is a leaf, measured.** `go list -deps ./internal/sharedprimitives` reports exactly
  one package inside the module — itself. No import cycle is possible, which is the purpose the
  request's "a package the duplicator already imports" rule exists to serve.
- **One declared widening.** `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req557_test.go`
  is new. D-02 requires a test that fails when the restored guard is disabled again, and the plan's
  Scope named only the recovery-fixture file, which this is not. The package already carries a per-REQ
  test-file convention (`finalization_req499_test.go`, `_req512_`, `_req547_`, `_req560_`, `_req565_`).
  Declared in the frontmatter, the Scope list, the commit message and the hand-back.
- **A path-form defect in this record was corrected, not worked around.** `qualify` resolves the
  backticked paths in `## Implementation Summary` against the repository root, and `scope-drift`
  compares them against `## Scope`; the builder wrote both in module-relative shorthand, which made
  `qualify` report seven missing paths. Both sections now carry the repository-relative form and say why.
- **Gate green on the merged tree.** `Maintainer verification passed.`, exit 0, wall 82s, CLI module at
  794 tests — ten more than the 784 at the branch point. The builder's own worktree run reported the
  same 794 with `EXECUTING (fingerprint_mismatch)`, so the stage really ran rather than reusing evidence.

### Remediation qualification (after review)

**Passed.** Remediation merge range `2acb34fb..234eda9`, two test files, 176 insertions and 1 deletion.
The review scored 86% and called the engineering sound; both of its code findings are closed, six record
claims are corrected in place, and one thing the review asked for turned out to be impossible and is
recorded as such rather than faked.

- **The restored guard is now pinned at its call site.**
  `TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget` drives
  `prepareBoundJournal` with a manifest whose `commit_paths` omits the planned archive target and
  asserts the error names it, that no journal is returned, and that none is left on disk. Deleting the
  `if len(missing) > 0` block reds it: `prepare accepted commit_paths that omit
  do-work/archive/REQ-721.md ... the omitted-target guard no longer fires`.
- **The other mutation the review asked for cannot be pinned, and the builder said so instead of
  claiming otherwise.** Rewriting the call site as `missing, _ := missingCommitPaths(...)` leaves the
  package green — and no behaviour test can change that, because no accepted manifest can reach the
  helper with a path it would refuse: `state_plan.go:111-112` refuses a target the snapshot did not
  resolve with `REQUEST-SNAPSHOT-STALE`, and every release target passes `publication`'s `containedPath`
  check. The propagation is covered inside the helper by the first test; the call site is covered for
  the guard it feeds. The test file's header now says exactly that, per test, replacing a comment that
  claimed one test failed "the moment the error is discarded again, whatever the helper is called".
- **The semver change is measured, not described.** The predicate was run beside a reconstruction of the
  pre-change function over an ordered cross product of 21 versions: **441 pairs, 242 changed, every one
  of them toward a refusal.** 159 skip-to-error and 83 stamp-to-error. Nothing that used to error now
  succeeds. The exact rule predicts all 242: `semverMajor` accepts both sides, the majors are equal, and
  at least one side fails the strict parser. `semverMajor` wants only three dot-separated parts with an
  `Atoi`-able first part, so a leading zero, an empty part, a trailing dot, a negative or non-numeric
  third part, a `-rc1` or `+1` suffix and a trailing space all get past it and used to be scored
  numerically. Five of those pairs are now a table test, including `1.0.0` against `1.09.0`, which is a
  stamp-to-error case.
- **The count is a property of the grid; the rule is the durable fact**, and the record says so. A
  reviewer reading "83 pairs" should not think 83 is a property of the code.
- **Six claims in this record are corrected rather than defended**: D-05's unreachability sentence,
  D-03's one-directional description, D-01's claim that every deleted copy returned non-nil, the
  deletion count, acceptance criterion 1's stale six, and the Scope rule that did not describe three of
  the paths it added. The picks themselves all stand.

## Testing

### Remediation testing (after review)

**Both new tests were shown red by ablation and green after restore.**

- `TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget` — green at HEAD; red with
  the `if len(missing) > 0` refusal deleted, naming the omitted path and the effective set it accepted.
  Green, and provably unable to be otherwise, with the call site's error discarded: see the
  qualification above.
- `TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder` — five rows, three of
  them refusals; green at HEAD, and the three refusing rows red when the `if !versionsParsed` guard is
  removed.

**The measured grid is the evidence for D-03, not a claim about it.** 441 ordered pairs run through both
the shipped predicate and a verbatim reconstruction of its predecessor; 242 changed; predicted-by-rule
242 equals actual-changed 242, which is what makes the rule a description rather than a summary.

**Gate on the merged tree:** `Maintainer verification passed.`, exit 0, gate wall 79s, CLI module at
**796** tests — up from 784 at the branch point, 794 after the first merge. Two `heavyverification` tests
fail identically at the branch point for environmental reasons and are not this request's; the builder
checked that rather than reporting them as its own.

## Review

**Overall: 86%** | 2026-09-06T07:10:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 88% |
| Code Quality | 95% |
| Test Adequacy | 72% |
| Scope | 90% |
| Risk | Low-moderate |
| Acceptance | Partial at review, Pass after remediation |

**Verdict: Accept with record corrections** — the engineering is clean and the hard parts are right.
Fifteen definitions became seven, the new package is a true leaf whose `go list -deps` returns only
itself, every touched package's module-internal closure changed by exactly `+sharedprimitives`, the
release-guard inversion risk the plan named as the biggest hazard is closed and pinned by five existing
tests, and the lock-in is red in both directions. What did not hold up was the record: two of the seven
decisions stated their own extent wrongly, one count was wrong, one acceptance criterion still named the
number D-06 had corrected, and the one deliberate behaviour change was pinned inside its helper but not
at its call site.

Where the reviewers disagreed, and what was picked:

- Whether a non-canonical frontmatter `id:` can reach the two sorts. One reviewer accepted D-05's
  unreachability sentence; another disputed it. Settled by building the CLI from both revisions and
  running `next` against a fixture carrying `id: REQ-12x`: the base binary selects `REQ-100`, the head
  binary selects `REQ-12x`. Picked the disputing reviewer.
- The hunk-level facts of D-03. One reviewer's table had two wrong rows and another described the
  stamp cases as passing through silently, which they did not. Settled by running the real
  `migrateInterviewSession` in both trees; the substance survived both readings.
- Whether deleting the whole `if len(missing) > 0` block is a REQ-557 regression. Reproduced green — but
  the identical deletion at the branch point is also green, so it is a coverage gap the change inherited
  rather than caused.
- Whether the record's "808 top-level tests green" holds. Two reported failures reproduce identically at
  the branch point and are environmental.

**Important findings:**

- D-05's unreachability sentence is false, demonstrated end to end with two built binaries. The pick
  stands; the claim did not. — impact-user-visible → corrected in the record
- The one deliberate behaviour change (D-02) was pinned inside its helper but not at its call site, and
  the test file's own comment overclaimed. — impact-user-visible → fixed in remediation

**Minor findings:**

- D-03 described the skip-to-error direction and omitted 83 stamp-to-error pairs; nothing tested the new
  refusal. — impact-negligible → measured and tested in remediation
- The Implementation Summary said nine deletions where there are twelve. — impact-negligible → corrected
- Acceptance criterion 1 still pinned six, contradicting D-06 and the shipped lock-in; all five criteria
  were unticked. — impact-negligible → corrected
- D-02's precondition is filtered upstream, so the restored guard changes no observable behaviour today;
  the record read as a fixed live bug. — impact-negligible → stated in the remediation qualification
- D-01's producer enumeration missed one `repairvalidation` call site, and its "every deleted copy also
  guaranteed non-nil" is false for the very copy D-02 fixes. — impact-negligible → corrected
- The lock-in guards top-level `func` declarations in one `internal/` tree only. — impact-negligible →
  report only, inherent to a name-list ratchet whose comment says the name list is the pinned set
- Scope's stated rule for the widening did not describe three of the paths it added, and one request
  constraint is superseded without being named. — impact-negligible → corrected

**Requirements checklist:**

- [x] One definition per canonical name — delivered; the request's reproduce regex prints 15 at the
  branch point and 1 now, and the union regex prints exactly 7
- [x] No import cycle, and the shared package imports nothing else in the module — delivered, measured
  package by package
- [x] Each reconciliation is a named decision — delivered, though two stated their extent wrongly and
  are corrected
- [x] The lock-in is one assertion in the existing file, registration untouched, red in both directions
  — delivered
- [ ] Every call site keeps its observable behaviour except D-02 — **not delivered**, and now named:
  `nextselection`'s ordering under a non-canonical `id:`, and 83 version pairs in `knowledgecommands`
- [x] → D-02 carries a test that fails when the guard is disabled again — partially delivered at review,
  delivered in remediation for the reachable mutation, with the unreachable one explained rather than
  faked
- [x] Tests for the other two behaviour changes — delivered in remediation for the semver refusal; the
  ordering change is recorded, not tested, because it needs a repository carrying a non-canonical `id:`
- [x] D-01 is safe — delivered, confirmed dynamically by instrumenting the helper to panic on an empty
  input and running the whole suite
- [x] D-03's publication half is input-set-identical — delivered
- [x] D-05's parser-equivalence claim — delivered, over 1,304 hostile inputs including unicode digits,
  400-digit runs and int64 overflow
- [x] D-04 is zero behaviour change — delivered, the renamed body is byte-identical
- [x] Scope holds: 23 changed files are the 22 declared paths plus the record — delivered
- [x] Module hygiene: build, vet, gofmt clean — delivered

**Acceptance testing**

**Result: Partial at review, Pass after remediation.** Three reviewers built the CLI from both
revisions and ran it rather than reading it. The gate is green at 796 tests.

**Follow-ups created:** None — every finding is either fixed, corrected in the record, or a
report-only observation about a ratchet whose limits its own comment states.

*Reviewed by review-work action*

## Lessons Learned

- **"A package the duplicator already imports" was the right rule and the wrong package.** The request
  named `corehelpers`, and `go list -deps` puts it outside the closure of four of the five duplicating
  packages. The rule exists to prevent a cycle; a new leaf package with zero module-internal imports
  serves that purpose absolutely while failing the rule's letter. When a stated constraint and its
  purpose come apart, say which one you followed and why — that is what makes it a decision instead of a
  deviation.
- **A duplicate can be a different function wearing the same name.** `finalization`'s `uniqueSorted`
  was `result, _ := normalizeRepositoryPaths(paths); return result` — a validator with its refusal
  thrown away, which returned `nil` for the whole slice and silently disabled the guard on the next
  line. Counting definitions found it; reading them is what identified it. The audit that produced this
  request counted.
- **Two implementations that differ in strictness can also differ in sign.** The request described the
  two `compareSemver` copies as differing on unparseable input. They also returned opposite signs for
  every unequal pair, which it never mentioned. Merging by picking one body and leaving both predicates
  alone would have reversed a release guard with every test green, because nothing pinned the
  orientation. Read both bodies before you merge two functions, even when a request has already told you
  how they differ.
- **"Cannot reach" is a claim, and this one was false.** D-05 said an id with a trailing non-digit could
  not reach the two sorts, reasoning from where ids are *generated*. The frontmatter `id:` is raw text
  the selector *prefers*, and a mismatch only appends a warning nothing prints. A reviewer built both
  binaries and watched them select different requests. Reachability arguments have to trace the value,
  not the convention.
- **A test that calls the helper directly does not pin the call site.** The guard test reddened on a
  discard written inside `missingCommitPaths` and stayed green on the same discard written one line
  above the call. The fix was a second test through `prepareBoundJournal` — and it turned out the
  helper-level discard is unreachable from any accepted manifest, so the honest record says which
  mutation each test kills and why one of them cannot be killed at all.
- **A number in a record should say whether it is a property of the code or of the sample.** "83 pairs
  turn a stamp into a refusal" is true of a 21-version grid. The durable fact is the rule that predicts
  all 242 changed pairs. The record carries both, labelled.

## Orientation

`internal/sharedprimitives` is the canonical home for helpers more than one package needs. It imports
nothing else inside the module — `go list -deps` returns only itself — so any package can call it
without a cycle, which is the property that stops these helpers being copied again. It holds
`UniqueSortedStrings`, `SubtractStringValues`, `FirstNonNilError`, `CompareSemanticVersions` and
`ParseSemanticVersion`.

Four things about those helpers are load-bearing:

- `UniqueSortedStrings` keeps every value it is given, including the empty string, and always returns a
  non-nil slice. Two of the four copies it replaced dropped blanks; no producer can emit one today, and
  a blank that ever appears should be visible in an evidence list rather than vanish from it.
- `CompareSemanticVersions` returns `(ordering, parsed)` and uses the standard Go orientation — negative
  when the left argument is older. The two copies it replaced were **inverted**. The `parsed` flag exists
  so no caller can read "unparseable" as "equal", which is what one of them did.
- `repositorymodel.RequestIDLess` uses the permissive parser, the one that also assigns request identity
  in the same file. That is deliberate: a comparator disagreeing with the identity parser would order ids
  by one rule and name them by another. Its observable cost is recorded in D-05 — a repository carrying
  a non-canonical frontmatter `id:` gets a different queue order, and `next` selects a different request.
- `physicalPath` was **not** merged. `knowledgecommands` keeps it (walks missing ancestors, absence
  succeeds); `suiteinstall` has `existingPhysicalPath` (EvalSymlinks + Abs, absence is an error). The
  second one's error *is* the existence check before a containment test, so merging them would mean
  re-adding that check by hand.

`finalization` no longer has a dedupe helper of its own. Its old one discarded a validator's error;
`missingCommitPaths` propagates it, and `finalization_req557_test.go` says per test what each one pins —
one covers the propagation inside the helper, one covers the call site's guard, and its header records
why no test can cover a discard written at the call site.

The ratchet is Finding 4 in `_dev/tests/audit-lockins.sh`: the union of the old and new names must be
exactly seven definitions, a floor as well as a ceiling. It guards top-level `func` declarations in
`internal/` only, and its comment says the name list is the pinned set — a differently-named copy of the
same body is outside it, which is inherent to any name-list ratchet.

Recorded and unfixed: `knowledgecommands.semverMajor` is a third, looser version parser nine lines from
the deleted one; `suiteinstall.compareSemanticVersions` is a fourth comparator with its own orientation;
and `gittransaction.stringSet` and `corehelpers.stringSet` are two definitions of one helper with
different value types. All three are the same class as this request and outside its six names.
