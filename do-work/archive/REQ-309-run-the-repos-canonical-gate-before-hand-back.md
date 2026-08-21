---
id: REQ-309
title: "[impact-rule-change] Run the repo's canonical gate before hand-back, not only the changed area's tests"
status: completed
created_at: 2026-08-20T23:18:07Z
status_changed_at: 2026-08-21T00:00:00Z
claimed_at: 2026-08-21T15:41:01Z
completed_at: 2026-08-21T16:19:08Z
commit: d9bf150
user_request: UR-055
addendum_to: REQ-262
domain: general
route: C
impact: impact-rule-change
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-21T15:42:21Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
depends_on: []
maintenance: true
kb_status: pending
write_set:
- skills/do-work/actions/work.md
- _dev/tests/contract-regressions.sh
---

# Run the Repo's Canonical Gate Before Hand-Back

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the request, project rules, domain guidance, and declared primes; recorded a
  Route C plan, exploration, scope, acceptance criteria, and policy decisions before implementation.
- [x] **[APPLY]:** Added the semantic regression contract first, captured RED, then changed only the
  declared action and contract-test files until the same test passed.
- [x] **[UNIFY]:** Reviewed the two-file diff, ran syntax and whitespace checks, reproduced the
  historical defect, confirmed scope alignment, and verified no debug artifacts were introduced.

## What

REQ-283 archived with a Testing section listing four green checks — `go test ./...` and a fresh build in `skills/do-work-board/tools/queue-kanban`, plus `queue-kanban verify` returning `OK: no findings`. Every one of them was true. None of them was `bash _dev/tests/maintainer-verify.sh`, which the change had just turned red by adding a second `./actions/board.md` routing row that `_dev/tests/staged-skills-contract.sh` counts.

The gate stayed red across REQ-279, REQ-295 and REQ-283's own metadata commits, and REQ-262's run was the first to notice — because Step 5.75's pre-flight ran the command and reported the baseline failing.

`actions/work.md` Step 6.5 resolves test commands from the prime file's testing section, then falls back to generic detection per changed file. Both are **area-scoped by construction**: they answer "what tests cover the files I touched", which is the right question for a regression and the wrong question for a repo that also has one whole-repo gate its own `CLAUDE.md` calls "the canonical baseline pass/fail check before any hand-back."

## Context

Discovered by REQ-262's orchestrator while repairing the gate in order to be able to verify REQ-262 at all (REQ-262 `## Decisions` D-01). The repair itself shipped separately at version 0.222.5, commit `8e9cc46` — this REQ is about the process gap that let the breakage land and persist, not about the four files it broke.

Worth noting what did *not* fail: REQ-283's review scored it and passed it. The Restatement Sweep (`actions/review-work.md` Step 6) should in principle have caught a routing row whose count another file asserts on, but a routing table is not obviously "something other text restates" until you know the contract test counts its rows.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-262: the work pipeline's test resolution is area-scoped, so a change can pass Step 6.5 and Step 7 while leaving the repository's own canonical whole-repo gate red — which is what happened to `maintainer-verify.sh` at REQ-283 and stayed that way for three REQs. Should I process this as a new task? → Confirmed: Yes, add to queue
  *(2026-08-21)* User confirmed via `do-work clarify`: process as a new task. Nothing put out of scope — the two follow-on decisions (where the knowledge lives, gate-vs-report) are still open and belong to whoever builds this REQ.

**Two things to decide if you say yes**, because they pull in different directions:

1. **Where the knowledge lives.** A prime file's testing section is the sanctioned home for project-specific test mapping, and this repo's `_dev/primes/` has no such section for the gate. Adding one is the smallest change and needs no action-file edit — but it only helps REQs whose `prime_files` list it. The alternative is a line in `actions/work.md` Step 6.5 that says a whole-repo gate, where the project declares one, always runs in addition to the area-scoped commands. That reaches every REQ and costs a rule.
2. **Whether it is a gate or a report.** Running the gate on every REQ is slower (it is a multi-minute suite here) and would have blocked REQ-279 and REQ-295 on a failure neither caused. Reporting it without blocking keeps those REQs moving but is exactly the posture that let three of them ship past a red gate without anyone noticing.

Per `crew-members/maintenance.md`, the deletion questions were asked first: the drift is not caused by a stale source, a bad example, or a too-broad tool. It is a genuinely absent instruction, which is why this is written as an addition candidate rather than a removal — and why it wants a replay case (REQ-283's diff, which must fail the new check) before anything is added.

---

## Triage

**Route: C** - Complex

**Reasoning:** This changes a repository-wide verification contract and must resolve both where a
project's canonical gate is declared and whether that gate blocks hand-back. The change spans work
orchestration, project-specific guidance, and regression coverage.

**Planning:** Required

## Plan

1. Add a condition-keyed rule to the work action's testing phase: when project guidance explicitly
   declares a canonical repository-wide pass/fail gate, run it from the project root on the final
   implementation or merged state in addition to focused tests. Require its direct exit status to
   be zero before successful archive, commit, or hand-back.
2. Keep focused tests and pre-flight attribution intact. A gate failure caused by the current diff
   uses the existing remediation loop; a pre-existing or unrelated gate failure stops successful
   hand-back with the claimed request and checkpoint preserved rather than being waived or fixed
   inside unrelated scope.
3. Add a focused contract regression that pins the trigger, additive behavior, repository-root and
   final-state execution, direct exit-zero requirement, no baseline exemption, and stop-before-
   hand-back outcome. Prove the contract RED before the action edit, mutation-test its load-bearing
   clauses, replay REQ-283's escaped failure from commit `2308afd`, then run the current canonical
   gate unpiped.

**Policy decisions:** Put the invariant in the shipped work action rather than a prime, because
prime test maps are selected per request and can repeat the blind spot. Treat the repository's own
declared gate as a hard gate rather than a report, because this repository explicitly defines exit
zero as the only hand-back proof.

**Plan validation:** Every captured concern maps to a task above; no task is orphaned, and the plan
has three implementation units.

*Generated by Plan agent*

## Exploration

The executable repository gate already exists and propagates failures directly. Project guidance
already declares it canonical, so the shipped work action must respond to that declaration by
condition rather than naming this repository's command or maintainer-only instruction files.

The minimal product seam is the testing phase in `actions/work.md`, after focused-test resolution
and before its failure handling. The existing baseline exemption remains valid for focused tests;
the new repository gate needs its own explicit no-exemption rule and a stop-before-hand-back path
for pre-existing or unrelated failures. The contract suite already contains semantic Python checks
and mutation matrices suitable for pinning those meanings without copying the production wording.

The historical replay is available at commit `2308afd`: its board router has two rows for one action,
so that commit's staged-skills contract and canonical gate fail even though its archived REQ records
green board-area checks. No change is needed in the gate itself, the project instructions, a prime
test map, or the testing-section template.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — require a project-declared canonical repository gate
  in addition to focused tests before successful hand-back
- `_dev/tests/contract-regressions.sh` (modify) — pin the rule semantically and mutation-test every
  load-bearing clause

**Files I will NOT touch:** `_dev/tests/maintainer-verify.sh`, `CLAUDE.md`, the prime files, and
`actions/work-reference.md`; they already provide the executable gate, project declaration,
maintainer index, and sufficient Testing record shape.

**Acceptance criteria (restated from REQ):**
- [x] Every REQ inherits a canonical repository gate declared by project guidance, regardless of
  its `prime_files` list.
- [x] Focused tests still run and remain the attribution signal; the repository gate is additive.
- [x] The declared gate runs from the project root on the final implementation or merged state and
  its direct exit status must be zero.
- [x] A red gate cannot be waived as a baseline failure or pass through successful archive, commit,
  or hand-back; unrelated pre-existing failures preserve the claim for resumption.
- [x] REQ-283's escaped regression reproduces as a nonzero canonical-gate run, while the current
  repository passes the new contract and the full canonical gate.

## Decisions

- **D-01 — Put the invariant in the universal work action, not a prime.** Prime selection is
  request-specific and was the blind spot this change closes; a condition keyed to explicit
  project guidance reaches every request without naming this repository in the shipped action.
- **D-02 — Treat the declared repository check as a hard gate.** Focused tests remain additive,
  but successful archive, commit, and hand-back require the canonical command's direct exit status
  to be zero.
- **D-03 — Preserve focused-test attribution without weakening the repository gate.** Existing
  baseline exclusions still apply to focused tests. They cannot waive the canonical gate: current-
  diff failures use remediation, while unrelated or pre-existing failures preserve the claim and
  checkpoint for resumption.

## Implementation Summary

- `skills/do-work/actions/work.md` (modified) — added the project-declared canonical repository gate lane to
  Step 6.5 and reflected it in the orchestration checklist.
- `_dev/tests/contract-regressions.sh` (modified) — added a semantic detector plus eleven in-memory mutations
  covering the trigger, additive behavior, execution context, direct verdict, remediation, and
  blocking failure paths.

The implementation stayed within the declared two-file write set. No gate script, project
guidance, prime, or testing template needed modification.

## Testing

**Tests run:**
- `bash _dev/tests/contract-regressions.sh` — RED before the action edit (exit 1: missing canonical-
  gate lane), then GREEN after the edit (exit 0: contract regression checks passed).
- `bash -n _dev/tests/contract-regressions.sh` — exit 0.
- `git diff --check -- skills/do-work/actions/work.md _dev/tests/contract-regressions.sh` — exit 0.
- Historical replay at `2308afdd582fd65f878a431fd7cc3c92a52b078d` — the staged-skills
  contract exited 1 on duplicate board routing (`found 2`); the historical canonical gate also
  exited 1, and after neutralizing an unrelated historical ShellCheck incompatibility in the
  isolated replay only, it reached and failed on the same staged-skills contract.
- `QUEUE_KANBAN_BROWSER=<headless-chromium> bash _dev/tests/maintainer-verify.sh` — exit 0, run
  directly and unpiped from the project root on the completed release tree.

**New tests added:** REQ-309 semantic Step 6.5 contract and eleven mutation cases.

**Existing tests updated (cross-REQ impact):** None.

## Review

**Overall: 79%** | 2026-08-21T16:18:11Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The later Error Handling row can still archive an unrelated or pre-existing canonical-gate
  failure that Step 6.5 says must preserve the claim and checkpoint. — impact-user-visible →
  REQ-317 created

**Minor findings:** 0 (report only)
**Acceptance:** Partial — the new lane works and its finding-closure evidence is adequate, but the
downstream generic failure row needs a traceable reconciliation.
**Suggested testing:** 2 items
**Follow-ups created:** REQ-317; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** A semantic detector plus adversarial mutations turned a prose policy into an
executable contract, and replaying REQ-283 proved the new gate catches the exact escaped defect.

**What didn't:** Testing only the newly inserted Step 6.5 lane missed a later generic Error Handling
row that can oppose it. REQ-317 carries that downstream-reader reconciliation.

**Worth knowing:** Focused tests establish attribution; a project-declared repository gate establishes
whether the final tree is hand-backable. They are complementary verdicts, not substitutes.

## Orientation

[MAP CHANGED] The work action now inherits any canonical repository-wide gate explicitly declared
by project guidance and requires its direct zero exit status on the final tree, independently of
per-request prime selection. The contract-regression suite mutation-tests that policy; REQ-317 will
align the remaining downstream error-handling reader.
