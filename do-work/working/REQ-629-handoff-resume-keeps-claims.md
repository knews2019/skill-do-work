---
id: REQ-629
title: '[impact-rule-change] Resuming a handoff must not reset its own claimed REQs'
status: claimed
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  basis:
  - Route B
  - 4-file write set
  - 2 subsystems involved
  - 4 acceptance criteria
  calculated_at: 2026-10-02T20:55:30Z
created_at: 2026-10-02T19:49:03Z
user_request: UR-133
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: ["REQ-628"]
batch: ur-132-follow-ups
required_lessons: ["skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#rule-direction-checked-against-callers"]
write_set: ["skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go", "skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go", "skills/do-work/actions/restart-with-parallel-handoff.md", "skills/do-work/actions/work-reference.md"]
dispatch_at: 2026-10-02T20:57:11Z
claimed_at: 2026-10-02T20:53:37Z
---
# Resuming a Handoff Must Not Reset Its Own Claimed REQs

## What
Make the handoff (`phandoff`, `actions/restart-with-parallel-handoff.md`) and the crash-recovery takeover agree on what "resume" means, so a session that resumes a handoff continues the handed-off claimed REQs from their written sections instead of resetting them.

## Why
On 2026-10-02 the UR-132 run was handed off with REQ-626 claimed, merged at 505ec74a and gate-green, and REQ-627 claimed with Triage, estimate, Exploration and Scope written. The handoff prompt said "run canonical recover, then continue each claimed REQ". After a machine restart, `recover` reported `RECOVERY-TAKEOVER-AVAILABLE` with the argv `recover --take-over REQ-NNN` ("working claim requires explicit authority"). Running it, as the finding suggested, returned both REQs to `do-work/queue/` as `pending` and stripped every orchestrator-generated section, including route, write_set, Scope, Implementation Summary, Qualification, Testing and Decisions of an already-merged REQ. The run recovered only by re-claiming and restoring the sections byte-for-byte from the handoff commit 1794c268 (decisions REQ-626 D-07 and REQ-627 D-05 in `do-work/archive/UR-132/`). The takeover is documented as a "canonical reset" (`actions/work-reference.md` → Crash Recovery), so the command did what it says; the handoff and the finding's suggested next step led the resuming session into it.

## Detailed Requirements
- Establish, by a real check, whether a resuming session on the same checkout can continue a handed-off claim with per-request `advance REQ-NNN --request-path <working path>` without any takeover, and what `recover`'s plain result then does on the next queue-mode `advance`.
- Fix the side that is wrong (see Open Questions) so that following the handoff's paste block plus the commands' own suggested next steps never resets a claim the handoff described as in flight.
- The handoff paste block (or the handoff action's rules for writing it) states what to run for each claimed REQ and that `recover --take-over` resets a claim, when that remains true after the fix.
- A REQ that is merged (`commit:` or `integration_at` present) must not lose its evidence sections through the documented resume path.

## Constraints
State is binding; prose is advisory (the handoff action's own rule): prefer a fix that the commands enforce over one that only adds prose, when the command change is small. No change to crash recovery for a genuinely crashed claim from another checkout. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-628 (board guide update).

## Builder Guidance
Medium certainty on the fix direction; high certainty on the failure. Reproduce it on a scratch repository first: claim a REQ, write sections, write a handoff, then follow the paste block in a fresh process.

## Open Questions
- [~] Which side changes: the handoff instructions, or the takeover command? → **D-01**: Builder chose: Both, minimally — the handoff action tells the resuming session to continue its own claims with per-request advance and never to run takeover on them, and `recover`'s takeover finding stops presenting `--take-over` as a plain next step for that case or states plainly that it resets the claim. Reasoning: the trap was the suggested argv in the command output; a prose-only fix leaves it there, and a new resume mode is only needed if per-request advance cannot continue a claim, which exploration checks first. Value: the documented resume path and the command's own next step agree, so a resumed run keeps merged work. Risk: a reworded or narrowed finding could hide the takeover from a genuinely crashed claim; reversible by restoring the argv, and the REQ's constraint keeps crash recovery for another checkout unchanged.
  Recommended: Both, minimally. Also: prose only, or a new resume mode in `recover`.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Red-Green Proof
**RED prompt/case:** In a scratch repo, claim a REQ, append Triage and Scope, write and commit a handoff, then in a new session follow the paste block and the `recover` finding's next argv: the REQ ends up in `do-work/queue/` as `pending` with Triage and Scope gone.
**Why RED now:** The handoff says "continue each claimed REQ", and the only next step `recover` offers for that claim is `--take-over`, which is a canonical reset.
**GREEN when:** Following the paste block and the commands' suggested next steps leaves the REQ claimed in `do-work/working/` with its sections intact, and the classifier names its next phase.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (17879 tokens, over the 2000 budget; `slugged: partial`). Matched: semantic recovery completeness, structured evidence projection. Narrowed at claim time to its `rule-direction-checked-against-callers` family (changing a finding's next step must be checked against every shipped caller that relies on it), which is now in `required_lessons`.
- `_dev/primes/lessons-action-files.md` (5879 tokens, over budget; `slugged: partial`). Matched: action routing, downstream readers. No family names the handoff-versus-command contract, so no narrower entry exists.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-133/input.md` for complete verbatim input.

*Source: capture both as REQs (the handoff-versus-takeover mismatch reported at the end of the UR-132 run)*

---

## Triage

**Route: B** - Medium

**Reasoning:** The failure and the desired outcome are clear and the Open Question has a recommended direction, but where `recover` decides a claim needs takeover, whether per-request advance can continue a claim without it, and what the handoff paste block says need discovery before dispatch. Two surfaces (one action file, one Go finding), not an architectural change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors and a scratch-repo reproduction: `do-work/runs/work-2026-10-02-204757/REQ-629-exploration.md`. The findings that change what the builder writes (paths under `skills/do-work/tools/do-work-cli/internal/`):

- **The authority check has no writer comparison.** `lifecycleadvance/recovery_commands.go:90` authorizes only `--assume-sole-authority` or `--take-over <this REQ>`; otherwise `:91-102` emits `RECOVERY-TAKEOVER-AVAILABLE` whose next argv is `recover --take-over REQ`. The writer label (`hostname:repoRoot`) is evidence text only, so after a restart on the same machine the takeover is still the only offered step. `recovery_commands_test.go:115` pins `"--take-over"` in the output with a foreign-writer fixture.
- **Bare `advance REQ-NNN` continues a claim; `advance REQ --request-path P` is refused at judgment phases** (`ADVANCE-GATE-INPUT-IRRELEVANT`). Reproduced in a scratch repo: plain `recover` changed nothing, bare `advance REQ-001` named the next phase, and `recover --take-over REQ-001` requeued the REQ and stripped route, Triage, Plan, Exploration, Scope and its checkpoint entry. The UR-132 handoff prompt used the refused `--request-path` form.
- **Queue-mode `advance` ignores working claims** (`nextselection/next_selection.go:371`), so a plain `do-work run` sees an empty queue and never resumes a claim; the resuming session must name each claim.
- **Prose:** `skills/do-work/actions/restart-with-parallel-handoff.md:49` tells the handoff writer to record each claim's "exact takeover command"; its paste-block rules say nothing about claimed REQs. `skills/do-work/actions/work-reference.md:310` (Crash Recovery) says "Follow the result". Other mentions (`docs/work-guide.md`, `run-with-recovery.md`, `forensics.md`) describe genuine crash recovery.
- **Smallest command fix:** in the not-authorized branch, make the next argv the read-only `advance REQ` (the heavy-lanes hold at `:115` already uses that shape) and make the stop reason say that `recover --take-over REQ` resets the claim and strips its sections. No writer comparison, so it cannot misfire on a foreign claim; `--take-over` stays available and documented for a genuinely crashed claim.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (modify) — the takeover finding's next argv becomes read-only `advance REQ`; the stop reason names `--take-over` as a reset
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modify) — update the pinned argv; add a test that a same-writer claim survives plain `recover` then `advance REQ` byte-identical in `working/`
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modify) — record each claim's continue command (bare `advance REQ-NNN`), not a takeover; the paste block says what to run per claim and that `--take-over` resets a claim
- `skills/do-work/actions/work-reference.md` (modify) — Crash Recovery says the takeover finding's next step is read-only and that `--take-over` is a reset for a genuinely abandoned claim

**Files I will NOT touch:** the reset itself (`requeststate/state_apply.go`), queue selection, `docs/work-guide.md`, `run-with-recovery.md`, `forensics.md` (they describe genuine crash recovery; the builder reports any sentence that tells a resuming session to take over), release paths, anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] A real check establishes whether per-request `advance` continues a handed-off claim without takeover, and what plain `recover` then does
- [ ] Following the handoff paste block plus the commands' suggested next steps never resets a claim the handoff described as in flight
- [ ] The handoff states what to run for each claimed REQ and that `recover --take-over` resets a claim
- [ ] A merged REQ (`commit:` or `integration_at`) keeps its evidence sections through the documented resume path

## Pre-Flight

**Git:** ✓ Integration tip 5f3038aa on `main`; the only dirt is this run's untracked `do-work/runs/` trail — no third-party paths
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 5f3038aa, gate wall 66s (do-work-cli fast stage reused by fingerprint)
**Baseline:** ✓ `bash do-work/runs/work-2026-10-02-204757/helpers/probe-629.sh` (`go test -count=1 -run Recover ./internal/lifecycleadvance/`) green. A first attempt passed the probe path as the test argv and could not launch (status 127, rerun 127); re-run with `bash <probe>` and satisfied.
**Dependencies:** ✓ Go toolchain present; no new dependency planned

*Checked by work action*
