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
write_set: [skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives.go, skills/do-work/tools/do-work-cli/internal/sharedprimitives/shared_primitives_test.go, skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go, skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go, skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go, skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go, skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go, skills/do-work/tools/do-work-cli/internal/publication/capture_files.go, skills/do-work/tools/do-work-cli/internal/publication/release.go, skills/do-work/tools/do-work-cli/internal/publication/release_mirrors.go, skills/do-work/tools/do-work-cli/internal/publication/answer.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T03:16:29Z
---

# Deduplicate six Go helper names defined fourteen times across do-work-cli

## What
`uniqueSorted`, `subtractPaths`, `requestIDLess`, `firstError`, `compareSemver` and `physicalPath` are defined fourteen times across `internal/`; in every case the duplicating package already imports the package holding an earlier copy, and three names have copies that disagree. Export one canonical definition per name in the lowest already-imported package (`corehelpers` for the path and error helpers, `repositorymodel` for `requestIDLess`), delete the other eight, and record the three semantic reconciliations as named decisions in this REQ.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Per-REQ helper files that duplicate an existing helper are the agent-creep class; the three semantic splits (`uniqueSorted` drops empty strings in one copy, `compareSemver` accepts unparseable input in one copy, `physicalPath` has two contracts) are latent correctness drift.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 4, sweep_key `per-req-duplicate-go-helpers`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -70. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `internal/finalization/finalization_prepare.go` — `subtractPaths` and `uniqueSorted` duplicate `internal/corehelpers/checks.go` (introduced 761d8e6a, REQ-498; `finalization` already imports `corehelpers`).
- `internal/knowledgecommands/interview_commands.go` — third `uniqueSorted`, silently drops empty strings (01d920dd).
- `internal/dependencygraph/dependency_graph.go` and `internal/repositorymodel/repository_model.go` — two `requestIDLess` bodies in one commit (ac2e3acd, REQ-408) with two different number parsers; `internal/nextselection/next_types.go` — third `requestIDLess` (625d49aa, REQ-411; `nextselection` already imports `repositorymodel`).
- `internal/publication/capture_files.go` — `firstError` byte-identical to `corehelpers/checks.go` (cf111a50, REQ-413).
- `internal/knowledgecommands/interview_commands.go` — `compareSemver` returns 0 for unparseable input while `internal/publication/release.go` rejects it.
- `internal/suiteinstall/update_transaction.go` — `physicalPath` is `EvalSymlinks`+`Abs`; `internal/knowledgecommands/commands.go` walks missing ancestors; same name, different contract.
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
definitions of the union of the old and the new names and pins the total at **exactly 6** — a floor as
well as a ceiling, because a one-sided ratchet guards half the property (REQ-556's lesson):

```
rg -n --glob '*.go' --glob '!*_test.go' \
  '^func (uniqueSorted|UniqueSortedStrings|subtractPaths|SubtractStringValues|requestIDLess|RequestIDLess|firstError|FirstNonNilError|compareSemver|CompareSemanticVersions|physicalPath|ResolvePhysicalPath)\(' \
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

All Go paths below are under `skills/do-work/tools/do-work-cli/`.

**Files I will touch:**
- `internal/sharedprimitives/shared_primitives.go` (new) — the four exported helpers plus the exported semantic-version parser
- `internal/sharedprimitives/shared_primitives_test.go` (new) — unit tests for each, including the parsed-flag contract
- `internal/corehelpers/checks.go` (modify) — delete three private helpers, point their call sites at the shared package, keep the local string-set helper which has two other callers
- `internal/corehelpers/inventory.go` (modify) — two dedupe call sites
- `internal/repositorymodel/repository_model.go` (modify) — export the request-id comparator in place, keep its parser private
- `internal/dependencygraph/dependency_graph.go` (modify) — delete the duplicate comparator and its now-callerless number parser, use the exported one
- `internal/nextselection/next_types.go` (modify) — delete the third comparator; the strict validator beside it stays, it has nine other callers
- `internal/nextselection/next_targets.go` (modify) — call site
- `internal/finalization/finalization_prepare.go` (modify) — delete the wrapper whose validator error is discarded; the one site whose result feeds a guard calls the validator directly and propagates its error
- `internal/finalization/finalization_discovery.go` (modify) — eighteen call sites
- `internal/finalization/finalization_apply.go` (modify) — three call sites
- `internal/finalization/finalization_recovery_test.go` (modify) — one reference to a deleted private helper
- `internal/knowledgecommands/interview_commands.go` (modify) — delete two private helpers, rewrite the version predicate
- `internal/knowledgecommands/memory_commands.go` (modify) — two call sites
- `internal/repairvalidation/already_green.go` (modify) — delete the fourth dedupe helper the request does not name
- `internal/publication/capture_files.go` (modify) — delete the byte-identical error helper
- `internal/publication/release.go` (modify) — delete the inverted comparator and the private strict parser, rewrite the release guard
- `internal/publication/release_mirrors.go` (modify) — use the exported parser
- `internal/publication/answer.go` (modify) — one call site
- `internal/suiteinstall/update_transaction.go` (modify) — rename the path resolver, no behaviour change
- `_dev/tests/audit-lockins.sh` (modify) — one assertion block

**The declared set is larger than the write_set the request carried, and the widening is stated here
rather than discovered by a reviewer.** The request listed ten paths; it named the files holding the
definitions and not the files holding the call sites, and it missed the fourth dedupe helper entirely.
Every added path is either a call site of a deleted definition or a file the survey showed holds one.

**Files I will NOT touch:** the string-set helper in corehelpers, which has two callers unrelated to
this class; the strict numeric validator in nextselection, which has nine other callers and is the
right validator for target tokens; any action file, prime or document. No shipped behaviour outside the
CLI module changes.

**Acceptance criteria:**
- [ ] Exactly six definitions across the union of the old and the new names, and the lock-in fails in
  both directions
- [ ] Each of the reconciliations is a named decision stating the behaviour each caller keeps
- [ ] No import cycle, and the shared package imports nothing else inside the module
- [ ] Every current call site keeps its observable behaviour, except the one guard whose restoration is
  D-02 and which carries its own test
- [ ] Tests unchanged except where a test named a deleted private helper
