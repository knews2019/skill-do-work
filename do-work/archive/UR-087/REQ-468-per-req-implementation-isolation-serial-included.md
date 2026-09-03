---
id: REQ-468
title: 'Per-REQ branch/worktree isolation for all implementation, serial included'
status: cancelled
created_at: 2026-09-01T04:29:16Z
user_request: UR-087
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-469, REQ-470, REQ-471, REQ-472]
batch: non-blocking-orchestration
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/crew-members/background-agents.md, skills/do-work-knowledge/crew-members/background-agents.md, skills/do-work-toolbox/crew-members/background-agents.md, _dev/tests/contract-regressions.sh]
claimed_at: 2026-09-03T16:19:49Z
completed_at: 2026-09-03T20:48:01Z
---

# Per-REQ Branch/Worktree Isolation for All Implementation, Serial Included

## What

Make every REQ's implementation run on its own per-REQ branch/worktree — serial runs included, not just `--fan-out` — so that a REQ set aside after implementation edits can never contaminate the next REQ's diff, tests, qualification, staging, or commit. Integration stays serial with one releaser per queue.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`, `crew-members/coding-guardrails.md`, `crew-members/anti-slop.md` and `crew-members/communication-style.md`. Approach in `## Plan` below.
- [x] **[APPLY]:** Six files changed, all in the Scope list below (`write_set` extended for the two sibling `background-agents.md` copies — D-06).
- [x] **[UNIFY]:** `git diff --stat` reviewed line by line. `skills/do-work/actions/work.md` — every remaining "serial" mention re-read; the survivors are all about the serial *scan/integration* (concurrency), not about where implementation is written. `skills/do-work/actions/work-reference.md` — same sweep; the Isolation ladder reads coherently with the Naming, hand-back, cleanup and Fan-Out sections that follow it. The three `background-agents.md` copies — `diff` against each other shows only the three pre-existing citation-depth regions, so no copy inherited another package's paths. `_dev/tests/contract-regressions.sh` — `bash -n` clean and the full suite exits 0. No debug artifacts, no new shell fence, no version or changelog edit.

## Detailed Requirements

- "Make setting aside safe after implementation edits. A blocked REQ's code, tests, decisions, and evidence must remain durable without contaminating the next REQ's diff, tests, qualification, staging, or commit. Use per-REQ branch/worktree isolation for implementation, including serial runs, or an equally safe durable mechanism." — the user chose always-on per-REQ branch/worktree isolation interactively during capture (over a set-aside-time-only holding branch).
- Acceptance test: "A blocked REQ with implementation edits cannot affect another REQ's diff, qualification, tests, staging, or commit."
- "Keep one releaser per queue and preserve explicit per-REQ commits, archive rules, changelog/version behavior, scope evidence, and UR closure semantics."

## Constraints

- The per-REQ machinery already exists and is written per REQ, not per wave: `skills/do-work/actions/work-reference.md` → Worktree Dispatch Mode ("Everything in this section is written per REQ and therefore already holds for any number of concurrent builders") and is documented as hand-drivable outside `--fan-out`. Serial mode becomes single-builder dispatch of that machinery rather than a new mechanism.
- Integration stays serial: merge → qualify → test → review → changelog → archive, one REQ at a time (`actions/work.md` "integration stays serial"). This REQ changes where implementation happens, not who integrates.
- Respect the existing degrade path: worktree support already degrades silently to serial dispatch when unavailable (`work-reference.md` → Worktree Dispatch Mode precondition). Define what "isolated" means on the floor agent without worktree support (a plain per-REQ branch is acceptable; define the degrade explicitly rather than silently losing isolation).
- Serial-only files stay serial-only: `actions/version.md` and `CHANGELOG.md` remain integrator-owned; builders never bump them.
- `_dev/tests/contract-regressions.sh` predicates that pin any edited prose (fan-out block, worktree dispatch block, hand-back block, State-stays-home assertions) must change in the same commit as that prose.
- Agent compatibility floor: action files must remain followable by the simplest agent that can read/write files and run shell commands.

## Dependencies

None — this is the batch's root. REQ-469 (blocked set-aside) depends on this isolation being in place; REQ-470/471/472 follow.

## Builder Guidance

Certainty: Firm on the decision (always-on per-REQ isolation, serial included — user-confirmed); latitude on the mechanics of expressing serial runs as single-builder dispatch and on the exact degrade path wording. Prefer reusing Worktree Dispatch Mode's existing per-REQ contract over writing a parallel serial-isolation mechanism.

## Open Questions

None — the isolation mechanism was resolved with the user during capture.

## Red-Green Proof
**RED prompt/case:** `actions/work.md` currently states the serial loop builds in the main tree with concurrency opt-in ("This action processes one REQ at a time unless you ask it not to… concurrency is opt-in rather than resident"); a contract-regressions assertion that serial implementation runs on a per-REQ branch/worktree fails today.
**Why RED now:** Serial implementation edits land directly in the shared working tree, so a set-aside REQ's edits would contaminate the next REQ's diff, tests, and commit.
**GREEN when:** `actions/work.md`/`work-reference.md` instruct per-REQ branch/worktree implementation for every run mode including serial; a new contract-regressions lane pins it and `bash _dev/tests/contract-regressions.sh` exits zero.
**Validation:** User confirmed (isolation choice selected interactively during capture)

## Full Context
See `do-work/user-requests/UR-087/input.md` for complete verbatim input.

---
*Source: UR-087 — "Use per-REQ branch/worktree isolation for implementation, including serial runs, or an equally safe durable mechanism." (user selected always-on branch/worktree isolation)*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fixed by the REQ (always-on per-REQ isolation, serial included) and the machinery already exists; the work was finding every place that exempts a serial run and deciding what "isolated" means where no worktree exists.

**Planning:** Not required beyond the approach below.

## Plan

The invariant is already written at `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** ("a builder writes only its own tree and owns no queue state"). What contradicted it was one label plus a set of serial carve-outs, so the change is mostly deletion:

1. Delete the `**Optional, advanced harnesses only.**` label that made per-REQ isolation opt-in, and the serial exemptions in `actions/work.md` and `actions/work-reference.md` (no merge range in serial, no `DO_WORK_DIFF_RANGE` in serial, the serial working-tree branch of the blocked-flip probe, "serial mode has no isolated committed implementation range", the serial-vs-worktree provenance fork, and the serial clauses in the two fail-safe-stop rows).
2. Replace the single silent degrade with an explicit three-rung **Isolation ladder**: worktree rung, plain per-REQ branch rung under the same `worktree-agent-REQ-NNN-<suffix>` name, and a no-rung case (not a git repository) that reports one progress line.
3. Re-key provenance and the two late-attribution fail-safes on facts that survive the change — "does this REQ have an implementation commit" and "can the saved base be isolated" — instead of on the run mode.
4. Reuse the existing branch name rather than adding a mechanism, so the crash sweep and `actions/cleanup.md` Pass 5 keep finding leftovers.
5. Add one grep-only contract-regressions lane (two positives, two deletion-ratchet negatives) and update the two predicates that quoted deleted clauses.

Not done deliberately: no rename of "Worktree Dispatch Mode" (34 shipped citations, three test section boundaries), no new shell fence, and no script for the branch/merge sequence — it needs judgment at three points (categorising `do-work/` paths, conflict resolution, applying seams), and the REQ asks for reuse of the existing per-REQ contract over a parallel serial mechanism.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — the two axes, the branch-rung write boundary, delete the serial evidence exemptions, re-key provenance and the fail-safe rows
- `skills/do-work/actions/work-reference.md` (modify) — the Isolation ladder, the sweep's applicability, delete the remaining serial carve-outs
- `skills/do-work/crew-members/background-agents.md` (modify) — `do-work run` isolates every builder in every run mode
- `skills/do-work-knowledge/crew-members/background-agents.md` (modify) — same sentence, sibling citation depth preserved
- `skills/do-work-toolbox/crew-members/background-agents.md` (modify) — same sentence, sibling citation depth preserved
- `_dev/tests/contract-regressions.sh` (modify) — one new REQ-468 lane, two predicate updates

**Files I will NOT touch:** `VERSION`, `CHANGELOG.md` and their mirrors, `actions/version.md`, `_dev/tests/maintainer-verify.sh`, `do-work/lessons-index.md` (a concurrent release owns them); `actions/review-work.md` and `docs/work-guide.md` (routed to REQ-471, see Discovered Tasks); `actions/cleanup.md` (see Discovered Tasks); no Go changes.

**Acceptance criteria (restated from REQ):**
- [x] `actions/work.md` and `actions/work-reference.md` instruct per-REQ branch/worktree implementation for every run mode including serial
- [x] The degrade is defined explicitly, with the one rung that loses isolation named and reported
- [x] One releaser per queue, per-REQ commits, archive rules, changelog/version behaviour, scope evidence and UR closure semantics unchanged
- [x] A new contract-regressions lane pins it and `bash _dev/tests/contract-regressions.sh` exits zero
- [ ] Acceptance test "a blocked REQ with implementation edits cannot affect another REQ's diff, qualification, tests, staging, or commit" — **not ticked**: see Testing for why no test in this repo can prove it

## Decisions

- **D-01 — "isolated" on the plain-branch rung means committed work (DECIDE & STATE).** A per-REQ branch in a shared working tree isolates only what is committed; an uncommitted edit is still in the next REQ's `git status`. Taken the REQ's own wording at face value ("a plain per-REQ branch is acceptable") and added one explicit sentence to the branch rung: the rung isolates committed work only, so setting a REQ aside there requires the builder's edits committed on its branch first, and an uncommitted set-aside must never be reported as isolated. The alternative reading (isolation must also cover uncommitted edits) would need a commit-or-refuse rule at set-aside time, which is REQ-469's subject, not this one. That sentence is what REQ-469 (making a blocked set-aside safe) needs to relax its "no substantive implementation edits landed" clause safely.
- **D-02 — builder and orchestrator are roles, not processes (DECIDE & STATE).** On the floor agent there is no second process to be a builder. Stated in one clause in both files: one agent may play both roles and the boundary between them is the branch, not the process. Without it, Step 6's "In worktree dispatch mode, never write the main tree" and the reference's "State stays home" read as process-level rules and the REQ becomes unimplementable on exactly the agent the compatibility floor protects. What keeps the boundary mechanical on that rung is that no `do-work/` path is ever committed on a builder branch — the hand-back sequence's step 2 queue guard already checks precisely that.
- **D-03 — the degrade is resolved by axis (DECIDE & STATE).** The concurrency degrade stays silent, as today. Dropping from the worktree rung to the branch rung also stays silent, because the REQ is still isolated. Only the no-git-repository rung, where no isolation is possible at all, gets one progress line. Reversing this would either re-introduce the silent loss the REQ forbids or turn an ordinary harness limitation into noise on every run.
- **D-04 — provenance re-keyed on having an implementation commit (DECIDE & STATE).** The prose used to fork on serial vs worktree; every isolated run now has a merge hash, so `supplied_commit` is the ordinary path and `primary_commit` is left for a finalization with no implementation commit to record (the already-green repair no-op and a terminal `fail` archive today, stated as a condition rather than a closed list). No Go change: the CLI accepts both modes from the manifest and never inspects the run mode.
- **D-05 — the two fail-safe stops re-keyed, not deleted (DECIDE & STATE).** "A serial dirty implementation without isolated base/merge authority" and "serial late failure" were the wording, but the real condition is that the saved base cannot be isolated — which still happens on the branch rung, where there is no `git worktree` for a detached diagnostic checkout. Both rows now say "an inability to isolate the saved base" / "unverifiable attribution", so the safety property survives the deletion of the serial framing.
- **D-06 — scope extended to the two sibling `background-agents.md` copies (DECIDE & STATE).** The REQ's `write_set` named only the core copy, but `_dev/tests/shipped-package-reference-contract.sh` requires all three packages to carry the same paragraph, so the completion proof requires the file class the declaration named. Flagged here, proceeded, and extended both the Scope list and `write_set`. The three copies differ only in citation depth elsewhere in the file; the edited sentence carries no citation, so it is identical in all three.
- **D-07 — the crash sweep's applicability had to change with the ladder (DECIDE & STATE).** The sweep said "skip it entirely when the repo has no `git worktree` support", which would have left every branch-rung leftover unswept and made the ladder's reuse-the-name argument false. It now applies to every run in a git repository and only enumerates branches where worktree support is missing. The equivalent precondition in `actions/cleanup.md` Pass 5 is outside this REQ's write_set and is filed as a discovered task.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/crew-members/background-agents.md` (modified)
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified)
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Deleted the opt-in label and thirteen serial carve-outs across the two pipeline files, added the three-rung Isolation ladder plus the two-axes and roles-not-processes statements, re-keyed provenance and both fail-safe-stop rows on conditions that survive the change, and added one grep-only contract-regressions lane with four asserts. Net 74 insertions, 38 deletions; no new shell fence, no new file, no Go change.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ exit 0

**Red-green validation:**
- contract-regressions REQ-468 lane, assert 1 (`actions/work.md` states isolation for every run mode, serial included): ✗ before implementation → ✓ after
- contract-regressions REQ-468 lane, assert 2 (`actions/work-reference.md` names the plain per-REQ branch rung): ✗ before implementation → ✓ after
- contract-regressions REQ-468 lane, assert 3 (`actions/work.md` carries no `Serial mode (has no merge|omits the variable|ignores this bullet)` exemption): ✗ before implementation → ✓ after
- contract-regressions REQ-468 lane, assert 4 (`actions/work-reference.md` carries neither `Optional, advanced harnesses only` nor `no isolated committed implementation range`): ✗ before implementation → ✓ after

All four failed on the unmodified tree in one run of the real suite before any prose changed, and pass after.

**New tests added:**
- `_dev/tests/contract-regressions.sh` — one REQ-468 lane beside the fan-out lane: two positives (isolation resident for every run mode; the branch rung is named) and two deletion-ratchet negatives (no serial evidence exemption; no opt-in label and no serial no-isolated-range exception).

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` `repository_gate_defects` predicate "late fail-safe branches" (from REQ-309/REQ-492): now expects "an inability to isolate the saved base" instead of "a serial dirty implementation", matching the re-keyed clause in `actions/work.md` Step 6.5. The property it pins — a late red gate that cannot be attributed stops safely — is unchanged.
- `_dev/tests/contract-regressions.sh` `repeated_failure_defects` predicate "fail-safe stop branch" (same origin): now expects "unverifiable attribution" instead of "serial late failure", for the same reason.

**What this test lane cannot prove — read this before ticking the acceptance box.** The REQ's acceptance test is "a blocked REQ with implementation edits cannot affect another REQ's diff, qualification, tests, staging, or commit." No test in this repository can demonstrate that. The deliverable is prose that an agent follows, so the property holds only when an agent actually follows it; a behavioural fixture would have to stand in for the agent, and it would then be proving that the fixture branches correctly, not that the instructions do. The four asserts are the honest mechanical form of the weaker claim "no shipped file exempts a run mode from per-REQ isolation" — two positives that the rule and its degrade ladder are stated, and two negatives that the specific exemptions which existed are gone. The negatives matter more than the positives here, because "serial has no exemption" is only checkable as the absence of the exemptions. The acceptance criterion is therefore left unticked in Scope above.

**Also not proved mechanically:** that the branch rung actually works end to end on a harness without `git worktree`. That would need a run on such a harness, which this repo does not have. What is checkable is that the rung reuses the existing name, so the existing leftover sweeps enumerate it, and that its hand-back sequence is the same one the worktree rung uses.

*Verified by work action*

## Discovered Tasks

1. **`skills/do-work/docs/work-guide.md` frames `worktree-agent-*` branches as `--fan-out`-only** (around its worktree and cleanup passages). After this REQ every serial run creates one, so the user-facing guide understates what a plain `do-work run` does. Outside this REQ's `write_set`; route to REQ-471, which owns reader and doc consistency for this batch.
2. **`skills/do-work/actions/review-work.md` restates the serial-versus-worktree diff fork in three places.** With the working-diff branch no longer occurring in a git repository, those three restatements are stale. Outside this REQ's `write_set`; route to REQ-471.
3. **`skills/do-work/actions/cleanup.md` Pass 5 step 1 skips the whole pass when there is no `worktree` subcommand.** On the branch rung a leftover is a branch with no worktree, so that precondition drops exactly the leftovers this REQ creates on a floor harness. The equivalent precondition in Crash Recovery was fixed here (D-07); Pass 5 is outside this REQ's `write_set` and needs the same one-line change. Behavioural, not doc consistency — worth its own REQ rather than folding into REQ-471.
4. **`_dev/tests/contract-regressions.sh` REQ-459 calibration lane holds four dead partitions** (`work_serial`, `work_worktree`, `reference_serial`, `reference_worktree`): only the whole blocks feed `surfaces`, and the reference anchor never matches. Untouched here because this REQ did not edit their anchors. Cleanup, no behaviour change.

## Cancelled

- **When:** 2026-09-03T20:48:01Z
- **Why:** landed in place by 7a099340 (per-REQ implementation isolation resident for every run mode). The claim was held by the cloud VM writer that went idle; the maintainer released it on 2026-09-03 23:35 local.
- **Decided by:** user, via `do-work abandon`
