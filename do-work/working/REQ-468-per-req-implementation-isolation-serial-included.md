---
id: REQ-468
title: 'Per-REQ branch/worktree isolation for all implementation, serial included'
status: claimed
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
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/crew-members/background-agents.md, _dev/tests/contract-regressions.sh]
claimed_at: 2026-09-03T16:19:49Z
---

# Per-REQ Branch/Worktree Isolation for All Implementation, Serial Included

## What

Make every REQ's implementation run on its own per-REQ branch/worktree — serial runs included, not just `--fan-out` — so that a REQ set aside after implementation edits can never contaminate the next REQ's diff, tests, qualification, staging, or commit. Integration stays serial with one releaser per queue.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
