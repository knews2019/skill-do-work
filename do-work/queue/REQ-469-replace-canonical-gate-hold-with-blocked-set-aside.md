---
id: REQ-469
title: 'Replace the unrelated canonical-gate hold with a blocked set-aside'
status: pending
created_at: 2026-09-01T04:29:16Z
user_request: UR-087
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-468]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-468, REQ-470, REQ-471, REQ-472]
batch: non-blocking-orchestration
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh]
---

# Replace the Unrelated Canonical-Gate Hold With a Blocked Set-Aside

## What

When the declared canonical repository gate fails for reasons unrelated to the current REQ, stop preserving the claim and halting: mark the active REQ `blocked` with durable evidence, release it from `working/` back to `do-work/queue/`, and continue with the next runnable REQ. Preserve the implementation on the REQ's isolated branch/worktree (REQ-468) and resume it after the gate turns green. Never waive the gate.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- "Replace the unrelated/pre-existing canonical-gate hold that currently leaves a REQ claimed and stops the session:
  - Mark the active REQ `blocked`.
  - Record the gate command, failing tests, evidence that the failures are unrelated to the REQ, and the phase where it stopped.
  - Remove `claimed_at`, remove the current checkout's checkpoint entry, move the REQ back to `do-work/queue/`, and continue with the next runnable REQ."
- "Never waive the canonical gate. A blocked implementation cannot be archived, versioned, or released until the gate passes."
- "Record enough durable state to resume the same implementation after its condition clears. Do not discard completed work or restart from scratch unnecessarily."
- "Re-probe blocked conditions when possible. When the condition clears, return the REQ to a runnable state and resume its qualification/release pipeline."
- Acceptance tests: "REQ-A encounters an unrelated canonical-gate failure while REQ-B is pending: A becomes blocked and leaves `working/`; B is subsequently claimed and processed." / "A blocked REQ's checkpoint entry is removed when it leaves `working/`." / "When the gate becomes green, the blocked REQ resumes from its preserved implementation and can complete normally." / "A gate failure caused by the current REQ still follows remediation and code-failure handling rather than being misclassified as unrelated."

## Constraints

- Current hold to replace: `actions/work.md` Step 6.5 item 4 ("preserve the claimed REQ and its checkpoint and stop before successful archive, commit, or hand-back"), the Error Handling row that excludes the unrelated case from code-failure handling ("Step 6.5 owns that hold"), and the Step 6.5 orchestrator-checklist line. The current-diff branch of Step 6.5 item 4 (remediation loop, 3 attempts, `error_type: code` follow-up) stays as is.
- Reuse the existing mid-run blocked flip (`actions/work.md` "During any step, if progress is blocked…") rather than inventing a parallel release path: sets `status: blocked`/`blocked_by`/`blocked_at`, removes `claimed_at` and `route`, appends `## Blocked`, moves back to `do-work/queue/`, drops the checkpoint entry per the In-Progress Record rule (own entry removed whenever the REQ leaves `working/`), reports "released, continuing". Its "no substantive implementation edits" clause must be relaxed for this case — safe only because REQ-468's isolation holds the edits on the REQ's own branch/worktree.
- The `## Blocked` evidence payload for this case records: the gate command, the failing tests, the evidence the failures are unrelated to this REQ's diff, and the phase where the pipeline stopped — durable in the REQ file, not only in `CHECKPOINT.md`.
- Re-probe constraints: the selector's `blocked_check` probe is bounded at 30 seconds (`actions/work.md` Step 1 probe; `scripts/run-blocked-check.sh`), while canonical gates commonly run longer — so either a dedicated gate recheck (e.g. at claim time or a distinct recheck step), or an explicit carve-out; also `blocked_check` is never pipeline-invented (`actions/capture-reference.md` external-condition fields) — the gate command is project-declared, so if it is recorded as a probe the never-invent contract needs an explicit, narrow carve-out for it. Builder decides between dedicated recheck and carve-out; record the decision in the REQ.
- Crash-recovery interplay: a set-aside REQ back in `queue/` with `status: blocked` is outside crash recovery's input (recovery preserves `blocked`), which also fixes today's defect where a gate-held claim in `working/` is classified "own crash" next session and destructively reset (13 generated sections stripped, moved to queue) — the current "preserve for resumption" does not survive a restart.
- Resume semantics: on re-claim after the gate is green, restore/continue from the preserved branch/worktree implementation and rejoin the pipeline at qualification/testing — never restart implementation from scratch.
- Never waive: no archive, no version bump, no changelog entry, no release for a gate-blocked implementation until the gate passes. Only an orchestration-wide integrity failure may stop the whole run.
- `_dev/tests/contract-regressions.sh` canonical-gate lane predicates (including "preserve the claimed req and checkpoint" and "stop before successful archive, commit, or hand-back") pin the old behavior and change in the same commit as the prose.

## Dependencies

- REQ-468 (Per-REQ branch/worktree isolation) — releasing a REQ after implementation edits is only safe once isolation exists.
- REQ-470 (fold-first routing of the failures) and REQ-471 (readers/docs) build on this REQ.

## Builder Guidance

Certainty: Firm on the behavior (the spec enumerates it and maps onto existing mechanisms); latitude on the re-probe mechanics (dedicated recheck vs `blocked_check` carve-out) and on exactly how resume-from-preserved-branch is phrased into the claim step. Prefer extending the existing blocked flip over adding a new status or a parallel release path.

## Open Questions

- [~] Re-probe mechanism for gate-blocked REQs: dedicated gate recheck at claim/run time, or a narrow carve-out letting the pipeline record the project-declared gate command as `blocked_check` despite the never-invent rule? → deferred to builder; the 30-second probe bound and the never-invent contract are the constraints either choice must satisfy.

## Red-Green Proof
**RED prompt/case:** With REQ-A claimed and an unrelated pre-existing canonical-gate failure, `actions/work.md` instructs "preserve the claimed REQ and its checkpoint and stop" — pending REQ-B is never claimed; `_dev/tests/contract-regressions.sh` gate-lane predicates assert exactly that hold language.
**Why RED now:** One unrelated red test anywhere in the repository parks the entire queue behind a single claimed REQ, and the "preserved" claim is destroyed by crash recovery on the next session.
**GREEN when:** The instructions route the unrelated case through the blocked set-aside (REQ blocked with gate evidence, out of `working/`, checkpoint entry removed, run continues to REQ-B) while the caused-by-current-diff case still remediates; the updated contract-regressions lane pins the new language and `bash _dev/tests/contract-regressions.sh` exits zero.
**Validation:** Inferred during capture (from the spec's acceptance tests)

## Folded From REQ-472

Hand triage 2026-09-03, maintainer approved: REQ-472 (End-to-end regression scenarios for non-blocking orchestration) is cancelled. Its scenario list becomes the acceptance criteria for REQ-469, REQ-470, and REQ-471; each of those REQs proves the scenarios it owns in its own Testing section, and no separate test REQ exists. Every lock-in names the real failure it pins; no decorative smoke tests, and no new sentence pins in `_dev/tests/contract-regressions.sh`.

- REQ-A hits an unrelated canonical-gate failure while REQ-B is pending: A becomes blocked and leaves `working/`; B is subsequently claimed and processed, in serial and fan-out modes. (REQ-469)
- The unrelated failures create or fold into `pending-answers` REQs instead of being stored only in `CHECKPOINT.md`. (REQ-470)
- A blocked REQ with implementation edits cannot affect another REQ's diff, qualification, tests, staging, or commit. (REQ-469, with REQ-468)
- A blocked REQ's checkpoint entry is removed when it leaves `working/`. (REQ-469)
- When the gate becomes green, the blocked REQ resumes from its preserved implementation and can complete normally. (REQ-469)
- A gate failure caused by the current REQ still follows remediation and code-failure handling rather than being misclassified as unrelated. (REQ-469)
- Repeated runs over a queue containing gate-blocked REQs stay stable: no re-hold, no duplicate fold targets, no checkpoint residue. (REQ-470, REQ-471)
- A UR with gate-blocked members is not closed by UR-closure readers until those members resolve. (REQ-471)

## Full Context
See `do-work/user-requests/UR-087/input.md` for complete verbatim input.

---
*Source: UR-087 — "Replace the unrelated/pre-existing canonical-gate hold that currently leaves a REQ claimed and stops the session"*
