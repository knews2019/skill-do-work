---
id: REQ-453
title: 'Keep targeted UR dependency closures in the run'
status: completed
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-450, REQ-451, REQ-452, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
write_set:
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - _dev/tests/contract-regressions.sh
claimed_at: 2026-09-02T12:10:38Z
route: C
planning_at: 2026-09-02T12:19:08Z
dispatch_at: 2026-09-02T12:33:08Z
builder_handback_at: 2026-09-02T12:50:22Z
integration_at: 2026-09-02T12:50:22Z
review_at: 2026-09-02T12:53:31Z
remediation_at: 2026-09-02T13:04:07Z
re_review_at: 2026-09-02T13:14:02Z
estimate:
  p50_active_minutes: 40
  confidence: medium
  basis:
    - Route C
    - 5-file write set after exploration added an action-contract regression
    - 2 subsystems involved
    - 7 acceptance criteria
    - cross-route regression gates
    - full-suite verification
  calculated_at: 2026-09-02T12:10:44Z
completed_at: 2026-09-02T13:28:36Z
release_at: 2026-09-02T13:31:23Z
commit: 62ef510d
kb_status: promoted
kb_entry: REQ-453-keep-targeted-ur-dependency-closures-in.md
---

# Keep Targeted UR Dependency Closures in the Run

## What

Keep every pending member of a targeted user-request dependency closure in the authoritative run set, then re-evaluate downstream members after their prerequisites integrate. Do not falsely declare a dependent concurrently runnable during fan-out, and do not silently leave it behind when the targeted workflow stops after the returned set.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this targeted-run dependency-closure root cause.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add selector regressions first, then reclassify only UR-expanded `DEPENDENCIES-UNMET` records whose prerequisites reach a fixed point of runnable or fan-out-limited scoped candidates; preserve all existing evidence and update the targeted action contract plus its semantic regression.
- [x] **[APPLY]:** Added the planned fixed-point selector classification and targeted ledger, then remediated F1 by freezing the effective numeric fan-out bound, replaying without `--fan-out`, projecting frozen IDs, and bounding only projected dispatch. The A/B/C RED proved the selector already returns the needed authoritative records, so remediation did not change selector runtime code.
- [x] **[UNIFY]:** Re-ran the exact starvation and dependency-closure tests, the targeted semantic/mutation contract (including the mutation that restores bounded replay), scoped vet, gofmt, and diff checks; all passed. A full `contract-regressions.sh` retry reached the targeted contract cleanly but was contaminated later by unrelated concurrent syntax errors under `internal/finalization`, outside this REQ's five-file scope.

## Finding Provenance

- **Finding #7 — P1 — source:** `internal/nextselection/next_selection.go:210-212`

> ````text
> [P1] Keep UR dependencies in the targeted run — [prj].claude/skills/do-work/tools/do-work-cli/
> internal/nextselection/next_selection.go:210-212
> For a targeted UR-NNN containing pending B that depends on pending A, B is emitted as DEPENDENCIES-UNMET even though A is selected
> earlier in the same dependency-ordered target set. The targeted workflow stops after its returned selected set, so do-work run UR-
> NNN completes A and silently leaves B behind instead of draining the requested UR as specified in actions/work.md targeted mode.
> ````

- **Evidence:** UR members are expanded in dependency-depth order at `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go:57-86`, but every member is evaluated against the unchanged repository graph. `next_selection.go:187-213` excludes a pending dependent before an earlier member can integrate. Targeted mode stops on the returned set at `skills/do-work/actions/work.md:185-194`, while the estimate and run contract at lines `172-183` requires the loop to drain dependents.
- **Surface-cost result:** N/A — this is a direct execution-contract correction, not added defensive apparatus.

## Detailed Requirements

- Preserve the complete scoped closure or equivalent deferred run set returned for targeted UR execution.
- Re-evaluate a dependent after a prerequisite from the same targeted run integrates successfully.
- Continue until the targeted UR closure is drained or a genuine terminal blocker occurs.
- Preserve dependency-depth ordering.
- Do not classify a dependent as concurrently runnable while its prerequisite is still pending or in flight.
- Report genuine failed, cancelled, external, or unresolved dependency blockers truthfully.

## Constraints

- The authoritative result must support targeted execution without a prohibited queue rescan.
- Fan-out must retain dependency safety.

## Dependencies

No request prerequisite. This REQ changes how targeted runs represent their own internal prerequisites; it does not depend on another captured REQ.

## Builder Guidance

Certainty level: Firm. Represent deferred-in-this-run separately from permanently excluded work if needed.

## Red-Green Proof

**RED prompt/case:** Target a UR containing pending B depending on pending A, with no external blocker, and execute only from the returned authoritative result.
**Why RED now:** B is excluded against the pre-run graph and targeted mode stops after A, leaving B queued.
**GREEN when:** A runs first, B is retained and re-evaluated after A integrates, and the targeted run drains both without making B concurrently runnable too early.
**Validation:** User confirmed after validate-feedback accepted Finding #7.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---

## Triage

**Route: C** - Complex

**Reasoning:** The fix changes the authoritative selector result contract and the targeted work-loop lifecycle across dependency integration and fan-out safety. It needs an explicit representation plan, caller-field audit, and end-to-end regression proof.

**Planning:** Required

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2299 tokens; matches selector result and evidence projection work but exceeds the 2000-token required-lessons budget and the partially slugged satellite cannot be narrowed safely. It was still read under the touch-conditional Lessons Discipline rule.

## Plan

1. **Represent targeted in-run holds without widening the result schema.** In `internal/nextselection/next_selection.go`, evaluate the complete targeted candidate set and reclassify only UR-expanded dependency exclusions whose unmet prerequisites are themselves progressable within the scoped invocation as `TARGET-DEPENDENCY-DEFERRED`. Keep them in `excluded`; keep ready-but-bound records as `FAN-OUT-LIMIT`; keep only runnable-now records in `selected`.
2. **Prove closure, ordering, and truthful blockers.** Extend `internal/nextselection/next_selection_test.go` with RED/GREEN chains, forks, fan-out bounds, mixed explicit/UR provenance, and genuine blocker cases. Preserve stable order and exact identity/provenance/state/outcome evidence, including replay argv.
3. **Drain the frozen targeted ledger through canonical re-evaluation.** Update `actions/work.md` and `actions/work-reference.md` so targeted runs retain initial `selected`, `FAN-OUT-LIMIT`, and `TARGET-DEPENDENCY-DEFERRED` members; after each serial integration or fan-out wave, call canonical `next` again with the original target tokens/flags, consume only frozen members, and stop only when drained or a genuine exclusion replaces the temporary disposition. No action-side queue rescan or reason parsing.

**Consumer field contract:** deferred records retain `request_id`, `request_path`, `title`, `provenance`, `selection_priority`, `original_status`, all probe/unblock fields, collection=`excluded`, code=`TARGET-DEPENDENCY-DEFERRED`, an actionable reason, and exact next/verification argv. The work action branches only on collection/code and re-evaluates canonically.

**Verification:** focused nextselection tests; full do-work-cli tests and vet; exact Go 1.25 compatibility; action contract regressions; canonical maintainer verification.

*Generated by Plan agent*

**Plan validation:** All seven detailed requirements map to the three tasks; no orphan task or uncovered requirement was found. The three-task plan stays below the complexity warning threshold, and every action-owned mutation consumer has explicit identity, provenance, state, outcome, and replay fields.

## Exploration

- `next_targets.go` already expands UR membership in dependency-depth order and preserves explicit-versus-UR provenance.
- `next_selection.go` evaluates every candidate against one unchanged graph, so `evaluateCandidate` emits `DEPENDENCIES-UNMET` before the selector can distinguish an in-run prerequisite from a genuine blocker. A post-evaluation fixed point is the narrow seam that can retain deeper chains without making dependents runnable early.
- Existing `SelectionExclusion` fields and invocation-wide `verification_argv` carry the required identity, provenance, probe, and replay evidence; no result-model schema change is needed.
- `work.md` currently treats an empty targeted `selected` set as terminal and stops after the initially returned set. Its Step 10 fresh-queue path cannot drain a targeted closure without violating the selector-only queue-read contract.
- `work-reference.md` already owns canonical recomputation and `FAN-OUT-LIMIT`; it is the correct home for the frozen targeted-ledger contract.
- Exploration expanded the planned scope by one regression file: `_dev/tests/contract-regressions.sh` must pin the new shipped action directives because its existing checks do not cover targeted closure draining.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modify) — classify targeted in-run dependency holds after candidate evaluation
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modify) — prove chains, forks, fan-out safety, provenance, evidence, and genuine blockers
- `skills/do-work/actions/work.md` (modify) — drain a frozen targeted ledger through canonical re-evaluation
- `skills/do-work/actions/work-reference.md` (modify) — define selector dispositions and targeted replay semantics
- `_dev/tests/contract-regressions.sh` (modify) — pin the shipped targeted-drain action contract

**Files I will NOT touch:** result-model schema/renderers, target expansion and dependency-graph packages, other actions, or queue records outside this REQ lifecycle.

**Acceptance criteria (restated from REQ):**
- [ ] The authoritative targeted result retains the complete pending UR dependency closure as runnable, fan-out-limited, or in-run deferred members.
- [ ] Dependents are re-evaluated after same-run prerequisites integrate, using the exact original target tokens and flags.
- [ ] A targeted run continues until its frozen closure is drained or a genuine terminal blocker replaces a temporary disposition.
- [ ] Dependency-depth order is preserved across chains and forks.
- [ ] A dependent is never selected concurrently while its prerequisite remains pending or in flight.
- [ ] Failed, cancelled, external, missing, ambiguous, cyclic, assigned, negligible-filtered, and otherwise unresolved prerequisites remain truthful blockers rather than in-run deferrals.
- [ ] The action consumes only canonical selector records and never rescans, parses, or sorts the queue itself.

## Pre-Flight

**Git:** ✓ Working tree clean outside `do-work/`; lifecycle claim, scope, and baseline records are the only current changes.
**Tests baseline:** ⚠ `go test -count=1 ./...` launched and exited 1 because pre-existing `internal/corehelpers` command-fixture tests observed transient fixture files/status from sibling cases. `internal/nextselection` passed. Exact output is stored in `do-work/working/baseline-failures.txt` with `launched: true` in `baseline.json`.
**Dependencies:** ✓ Go test tooling launched successfully.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Added a post-evaluation fixed point that retains UR-expanded dependents as `TARGET-DEPENDENCY-DEFERRED` only when every unmet prerequisite is progressable inside the targeted invocation. Added selector coverage for chained, forked, fan-out-limited, mixed-provenance, evidence-preservation, genuine-blocker, and newly captured same-UR starvation cases. Defined a frozen targeted ledger that replays canonical `next` with the exact original targeting tokens and non-fan-out filters, projects authoritative records onto frozen IDs, and applies the saved numeric bound only to dispatch until the ledger drains or a genuine blocker appears. Mutation tests pin the shipped action contract and reject bounded replay observation.

## Qualification

Passed after remediation — 5 files verified, 7 acceptance criteria traced, P-A-U confirmed. The selector change is wired before fan-out bounding; the action readers branch on structured collection/code, project before scheduling, and preserve canonical replay authority. Focused `nextselection` tests and vet pass, including the A/B/C starvation case. Concurrent finalization-related files are unrelated, excluded from this REQ, and currently prevent an uncontaminated full-module rerun; no placeholder, debug artifact, undeclared REQ touch, or hollow data path was found.

## Review

**Acceptance:** Fail — one remediation required.

**Important finding (F1, P1):** A newly captured same-UR root can consume the original fan-out slot before projection onto the frozen ledger. In an initial `A → B` run with `--fan-out 1`, integrating A and then adding a new root C makes bounded replay select out-of-ledger C and leave frozen B as `FAN-OUT-LIMIT`; the action then stops without a genuine blocker. Remediation must freeze the effective numeric fan-out bound, omit `--fan-out` from canonical replay observation, project onto frozen IDs, and apply the saved bound only when dispatching projected `selected` records. Add the A/B/C replay-starvation regression and mutation-pin the distinction between replay filters and the scheduling bound.

**Score:** 50% — acceptance-failure cap. Scope and fixed-point classification otherwise passed review.

*Generated by Review agent*

## Decisions

<!-- D-XX counter: next D-02 -->

### D-01: Separate targeted replay observation from fan-out scheduling

**Decision:** Freeze the initial invocation's effective numeric fan-out bound, but omit `--fan-out` when replaying canonical `next`. Preserve the exact original targeting tokens and every non-fan-out filter, project the unbounded authoritative result onto frozen ledger IDs, and apply the saved bound only to dispatch of projected `selected` records.

**Reasoning:** The original plan replayed every flag exactly, which lets a newly captured out-of-ledger UR member consume a selector slot and starve frozen work. Separating observation from scheduling keeps the selector authoritative and concurrency bounded without growing the frozen ledger or inventing a blocker.

## Re-Review

**Acceptance:** Pass — F1 closed with no remaining findings.

The A/B/C regression reproduces bounded-replay starvation, then proves that canonical replay without `--fan-out`, frozen-ID projection, and the saved numeric dispatch bound drain the original closure without admitting the new same-UR member. Focused selector tests, scoped vet, semantic mutation checks, formatting, and diff checks pass.

**Score:** 100% — all requirements satisfied; risk low.

*Generated by Re-Review agent*

## Lessons Learned

**What worked:** A post-evaluation fixed point retained dependency chains without widening the result schema, and mutation-based action tests pinned the selector/action ownership boundary.

**What didn't:** Replaying the initial fan-out bound before projecting frozen membership let a newly captured out-of-ledger root consume the only slot and falsely strand retained work.

**Worth knowing:** In targeted replay, observe the complete canonical ready set first, project it onto frozen membership second, and apply the saved scheduling bound last. Bounding before projection changes semantics, not just throughput.

## Orientation

[MAP CHANGED] Targeted UR runs now retain their initial dependency closure in the do-work CLI selector contract and drain it through a frozen action-owned ledger. Dependents remain non-runnable until prerequisites integrate, while newly captured UR members cannot enter or starve the active run.

---
*Source: validate-feedback Finding #7, captured by UR-085.*
