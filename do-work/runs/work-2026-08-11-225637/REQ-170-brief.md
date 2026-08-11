# Builder Brief — REQ-170

## Dispatch

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-170-finding-closure-ratchet-and-rubric-home`
- Branch / operative name: `worktree-agent-REQ-170-finding-closure-ratchet-and-rubric-home`
- Durable hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-11-225637/REQ-170-handback.md`
- Route: C
- Domain: general
- TDD: false, but the captured finding-origin RED/GREEN proof is mandatory.

## Request and plan

Implement two single-home rules with minimal instruction growth:

1. Canonical Finding-Closure Ratchet in `skills/do-work/actions/work-reference.md`: a review/triage finding-origin REQ closes only with a named regression check/test that fails before and passes after, or deletion of the named finding surface. Bare patches, unrelated green tests, `tdd: false`, and high scores are not evidence.
2. Canonical one-paragraph earned-defense rubric in `skills/do-work/crew-members/coding-guardrails.md` § 2, preserving exactly: “what earned this, and is the fix still cheaper than the surface it added?”
3. `capture.md` states only its local hook: explicitly finding-origin capture GREEN names the intended regression test/check or exact deletion surface, independent of ordinary TDD inference.
4. `review-work.md` cites the ratchet, verifies matching evidence/deletion, and makes a miss both Important and Acceptance `Fail`; it also cites the earned-defense rubric in the existing Simplicity First pass.
5. `validate-feedback.md` keeps only triage-specific Surface-cost/Accept behavior plus citations and preserves accepted-feedback provenance in the capture handoff.

No schema marker, metrics, trend log, parser, or board work.

## Required context

Read completely before editing:

- `CLAUDE.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- the five scoped files below
- focused REQ-169 assertions in `_dev/tests/contract-regressions.sh` (read-only)

Important explored constraints:

- Keep the five-principle guardrail taxonomy unchanged; the rubric is one paragraph under Simplicity First.
- The work-reference insertion seam is immediately after the Step 6.5 Testing template.
- Preserve validate-feedback's exact existing tokens, including `what incident earned this, and is the fix still cheaper than the surface it added?`, the surface-adding boundary, direct-fix/deletion/simplification `N/A`, non-Accept routing, and `**Surface-cost:** N/A / Earned / Flagged`.
- Core callers use backticked `actions/...`/`crew-members/...` citations. Toolbox→core semantic prose uses backticked `../do-work/...` paths; do not turn that into a broken Markdown link.
- Main `capture.md` has an unrelated uncommitted screenshot-publication hunk. Your branch owns only the Step 1 finding-proof hook. The orchestrator will preserve/reconcile the main hunk.
- The pre-existing main-tree full-suite failure is only the uncommitted REQ-173 prime link; branch tests must not absorb or fix it.

## Scope

- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work-toolbox/actions/validate-feedback.md`

Do not touch tests, `do-work/`, `actions/work.md`, release metadata, or any main-tree path other than the exact hand-back file. If a durable test-source edit proves unavoidable, stop and request a scope extension in the hand-back; do not write it silently.

## Proof and verification

- Record the current RED state before editing.
- Because this REQ itself originates from a finding, give the focused closure check a stable descriptive name in the hand-back and show its exact command/output failing before and passing after. A committed test-source edit is not required if the named deterministic check exercises the shipped canonical/caller contracts directly.
- Run `bash _dev/tests/contract-regressions.sh`, `bash _dev/tests/shipped-package-reference-contract.sh`, focused canonical-home/restatement greps, `git diff --check`, `git diff --stat`, `git diff --numstat`, and a changed-path/debug-artifact audit.
- Report instruction lines added/removed and confirm net growth is a handful of lines.

## Builder rules

- Use `apply_patch`; commit on the operative branch; never rebase.
- Do not read or write a worktree `do-work/` snapshot.
- P-A-U evidence goes in the hand-back: PLAN, APPLY, UNIFY with exact files checked.
- Significant choices become D-XX entries with reasoning; out-of-scope findings are reported, not fixed.

## Hand-back format

Branch and commit; P-A-U evidence; named RED/GREEN closure proof; file manifest; tests/results; net-surface accounting; decisions; integration seams; discovered tasks/blockers. Write it to the exact absolute hand-back path, then return only a one-line status.
