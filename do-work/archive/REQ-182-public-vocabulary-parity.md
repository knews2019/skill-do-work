---
id: REQ-182
title: Public work and schema vocabularies drift while suites stay green
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T09:52:42Z
completed_at: 2026-08-15T10:12:17Z
user_request: UR-041
domain: testing
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-183, REQ-184, REQ-185, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [README.md, skills/do-work/SKILL.md, skills/do-work/docs/work-guide.md, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh]
route: C
kb_status: pending
kb_entry:
---

# Public Work and Schema Vocabularies Drift While Suites Stay Green

## What

Restore parity at the public work-guide/router and testing-schema/normalizer seams, and correct the two short workflow summaries that omit canonical states while the baseline suites remain green.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Keep the guide as the single public work-alias inventory; compare it bidirectionally with the core router in the existing contract suite. Compare the testing-status schema aliases bidirectionally with the Go normalizer in the same suite. Prove both helpers reject additions and removals on either side, then repair the six live mismatches and the two stale workflow summaries.
- [x] **[APPLY]:** Added the two parity ratchets first, recorded RED for both live mismatches, then synchronized the guide/router aliases, schema/normalizer aliases, README ownership pointer, schema gloss, and dependency-cycle summary vocabulary within the six-file scope.
- [x] **[UNIFY]:** Reviewed the full six-file diff and `git diff --stat`; verified the README/guide/router ownership chain, action/reference vocabulary, and both parser/mutation helpers. `bash -n`, the full contract suite, the board Go suite, the shipped-reference contract, and `git diff --check` all pass with no debug artifacts.

## Why

The public work guide advertises three aliases the first-match router cannot dispatch, workflow glosses omit canonical fields/statuses, and the board accepts two testing-status aliases absent from the schema table. These are user-visible or maintainer-facing drifts that existing suites do not detect.

## Context

- Audit priority: P2; impact 3; effort normal.
- Root-cause key: `public-vocabulary-parity`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 2.
- Reproduce: `rg -n 'do-work (begin|execute|build)|Other trigger words|enum/boolean fields|Queue: N pending|blocked-dependency-cycle|selected for testing|returned with feedback' README.md skills/do-work/SKILL.md skills/do-work/docs/work-guide.md skills/do-work/actions/work.md skills/do-work/actions/work-reference.md skills/do-work-board/tools/queue-kanban/testing.go`.

## Detailed Requirements

- Restore the pre-modular `do-work begin`, `do-work execute`, and `do-work build` aliases advertised by `skills/do-work/docs/work-guide.md`.
- Keep one public work-alias list and compare it with the router so additions or removals cannot occur on one side only.
- Delete or generalize the stale `enum/boolean fields` gloss in `skills/do-work/actions/work.md`.
- Make queue summaries include dependency-cycle holds rather than implying only `N pending`.
- Add the already-supported `selected for testing` and `returned with feedback` aliases from `testing.go` to the canonical `testing_status` schema table.
- Add compact bidirectional parity checks with mutation cases at both declared seams.

## Constraints

- Do not create a repository-wide prose/schema generator.
- The lock-in limit is zero one-sided aliases across the work-guide/router and testing-schema/normalizer seams.
- Keep runtime alias support and public documentation synchronized without widening unrelated vocabulary.

## Dependencies

None. This REQ is semantically independent, though its documented `write_set` overlaps REQ-181.

## Builder Guidance

Firm intent. The audit attributes the two parity checks to incidents `9ba534e` and `ea0fd94`; keep the assertions compact and seam-local.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Mutate either the work guide or router by one alias, and mutate either the testing schema table or normalizer by one alias; current baseline checks remain green and the two sides disagree.
**Why RED now:** Three documented work aliases are absent from the router, two runtime testing aliases are absent from the schema table, and two workflow summaries omit canonical states.
**GREEN when:** Seam-local tests fail for either one-sided addition or removal; all six current alias mismatches and both stale summary glosses are corrected.
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 02, labeled P2, impact 3, normal effort, among the eight validated audit findings.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 2 in the canonical audit.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*

## Triage

**Route C** — the requested behavior spans the public router/guide seam, the action/reference workflow contract, and the board normalizer/schema seam. The locations are known, but the parity rules and mutation-sensitive regression design require a coordinated multi-file plan.

## Plan

1. Add two compact source-contract checks to the existing regression suite: guide aliases versus the work router, and documented `testing_status` aliases versus `normalizeTestingStatus`.
2. Exercise each comparison helper with in-memory one-sided additions and removals so both directions of drift are proven to fail.
3. Make the guide's trigger block the sole public work-alias list, restore its missing router aliases, and replace README's duplicate list with a pointer.
4. Generalize the stale schema gloss, include dependency-cycle holds in both queue summaries, and document the two already-supported spaced testing aliases.
5. Run focused RED/GREEN evidence, the board Go suite, both repository handoff suites, and diff hygiene checks.

## Exploration

- `skills/do-work/docs/work-guide.md` advertises `begin`, `execute`, and `build`; the `./actions/work.md` router row omits all three and additionally supports `work`, which the guide omits.
- `README.md` maintains a second enumerated alias list, so it can drift independently from both the guide and router.
- `skills/do-work/actions/work-reference.md` omits the spaced aliases `selected for testing` and `returned with feedback` that `normalizeTestingStatus` already accepts in `skills/do-work-board/tools/queue-kanban/testing.go`.
- `skills/do-work/actions/work.md` names a stale closed set of schema fields and its queue summary omits `blocked-dependency-cycle`; the checkpoint `queue_state` example in `work-reference.md` has the same omission.
- Existing behavior tests cover individual router and normalizer behavior, but no test compares either declared seam. `_dev/tests/contract-regressions.sh` is therefore an earned scope addition required by this REQ's TDD contract; `testing.go` is comparison input and needs no runtime edit.

## Scope

**Files I will touch:**
- `README.md` (modify) — replace the duplicate public alias inventory with a pointer to the canonical guide list
- `skills/do-work/SKILL.md` (modify) — restore all guide-advertised work aliases in the first-match router
- `skills/do-work/docs/work-guide.md` (modify) — keep the complete canonical public work-alias list
- `skills/do-work/actions/work.md` (modify) — generalize the schema gloss and include dependency-cycle holds in the live queue summary
- `skills/do-work/actions/work-reference.md` (modify) — synchronize testing aliases and checkpoint queue-state vocabulary
- `_dev/tests/contract-regressions.sh` (modify) — add bidirectional seam-parity checks with mutation cases

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/testing.go` (its runtime aliases are already correct and serve as comparison input), unrelated queue files, or any repository-wide prose/schema generator.

**Acceptance criteria (restated from REQ):**
- [ ] The guide and router expose exactly one identical work-alias set, with `begin`, `execute`, and `build` restored.
- [ ] The testing schema table and runtime normalizer expose exactly one identical alias set, including both supported spaced aliases.
- [ ] One-sided additions or removals at either seam fail compact mutation-sensitive regression checks.
- [ ] The schema gloss no longer claims a stale closed field list, and both queue summaries include dependency-cycle holds.
- [ ] No unrelated vocabulary or files are widened.

## Pre-Flight

**Git:** ⚠ Four pre-existing edits under `do-work/queue/` (REQ-189–192) belong to other work and must remain unstaged; the declared source/test files are clean.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh` passed before implementation, confirming the drift was not previously detected.
**Dependencies:** ✓ Repository shell/Python tooling and the board Go module are available.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `README.md` (modified) — replaced the duplicate alias inventory with a canonical guide anchor
- `skills/do-work/SKILL.md` (modified) — synchronized the work router with all ten public aliases
- `skills/do-work/docs/work-guide.md` (modified) — added the router-supported `work` alias to the canonical public list
- `skills/do-work/actions/work.md` (modified) — generalized the schema gloss and surfaced dependency-cycle holds in the queue summary
- `skills/do-work/actions/work-reference.md` (modified) — documented both spaced testing aliases and dependency-cycle checkpoint state
- `_dev/tests/contract-regressions.sh` (modified) — added exact seam comparisons and bidirectional add/remove mutation probes

**Behavior:** Public aliases and testing-status aliases now have one documented inventory each plus executable parity mirrors; any one-sided addition, removal, or testing-alias remap fails the existing contract suite. Queue summaries no longer hide dependency-cycle holds.

## Testing

**RED:** With only the new contract block added, `bash _dev/tests/contract-regressions.sh` exited 1 and reported guide-only `begin`, `build`, `execute`, router-only `work`, and normalizer-only `selected for testing`, `returned with feedback`.

**GREEN:**
- `bash _dev/tests/contract-regressions.sh` — PASS
- `cd skills/do-work-board/tools/queue-kanban && go test ./...` — PASS
- `bash _dev/tests/shipped-package-reference-contract.sh` — PASS
- `bash -n _dev/tests/contract-regressions.sh` — PASS
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports the Implementation Summary exactly matches all six declared files; unrelated REQ-189–192 edits remain outside the implementation inventory.
- **Mechanical checks:** PASS — `qualify.sh` found all files present in the diff, all P-A-U phases complete, and no debug artifacts.
- **Substance and traceability:** PASS — each detailed requirement maps to a visible source change or one of the two exact seam comparisons; both comparison helpers reject one-sided additions and removals, and the testing map also rejects remapped aliases.
- **Wiring/data flow:** PASS — the existing contract-regression entrypoint executes the new checks; the router reads the guide-matched alias set and the board keeps using the unchanged runtime normalizer whose source is compared to the schema row.

## Review

**Result:** Approve — Acceptance: Pass  
**Overall score:** 98%

- **Requirements (100%):** All six drift instances, both stale summaries, and both seam-local parity ratchets are delivered.
- **Code quality (92%):** The embedded parity block is longer than the word “compact” suggests, but remains seam-local, readable, and justified by exact parsing, clear diagnostics, README ownership checks, and bidirectional mutation proof.
- **Test adequacy (100%):** Independent review confirmed additions and removals fail on either side of both seams, testing-alias remaps fail, and all repository checks pass.
- **Scope (100%):** The six changed files exactly match the declared Scope; unrelated queue edits remain excluded.

**Important findings:** None.  
**Minor findings:** Test block length only; no remediation required.  
**Explicit remediation:** None.

## Lessons Learned

- When prose is intentionally authoritative but runtime must remain independently readable, a seam-local exact comparison with bilateral mutation probes is enough to prevent drift without introducing a generator.
- A duplicate public inventory is itself a third drift surface; replace it with an anchored pointer before asserting parity between the remaining owner and mirror.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.
