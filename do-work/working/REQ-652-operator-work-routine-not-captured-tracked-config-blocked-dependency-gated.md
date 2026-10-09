---
id: REQ-652
title: '[impact-rule-change] Routine operator work is never a REQ; a tracked operator configuration is a blocked, dependency-gated REQ that clarify completes'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-09T16:09:30Z
created_at: 2026-10-09T16:06:43Z
user_request: UR-143
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-650, REQ-651, REQ-653]
batch: october-review-triage
write_set: ["skills/do-work/actions/capture.md", "skills/do-work/actions/clarify.md", "skills/do-work/actions/work.md", "skills/do-work/actions/work-reference.md", "skills/do-work/docs/work-guide.md"]
claimed_at: 2026-10-09T16:08:10Z
dispatch_at: 2026-10-09T16:11:57Z
builder_handback_at: 2026-10-09T16:16:57Z
---
# Routine Operator Work Is Never a REQ; a Tracked Operator Configuration Is a Blocked, Dependency-Gated REQ That Clarify Completes
## What
Capture, clarify and the run orchestrator get one small rule for work the operator does with their own hands. A routine operator act after the code ships (deploy, publish, approve, verify on live hosts) is never captured; the user does it on their own time, and only the AI-buildable parts become REQs. When the request itself asks to track a special operator configuration, capture that one REQ as `status: blocked` with `blocked_by` naming the operator, `blocked_at`, and `depends_on` naming the AI-buildable REQs from the same request, with the checklist in its body. The board already shows such a REQ under Pending → Waiting while a dependency is unmet and under Needs input · Blocked once all are met. Clarify's blocked prompt gains a fourth option, "Done — I did it myself", which records the receipts and flips the REQ to `completed`, so it shows under Done and `cleanup` archives it. A run that discovers mid-flight that only a routine operator act remains completes the REQ with what the AI built and names the operator step in the hand-back; it does not flip to `blocked`, does not fail, and does not recommend abandon.
## Why
Upstream proposal "operator work is not a REQ" (pasted 2026-10-09, verified against 0.305.84 as F11): one consumer REQ, "install, enable and verify CMS publishing in production", touched 53 commits over six days with three block/clarify/earmark round trips and two multi-hour sittings, and the user cancelled it: "I don't need a REQ for it, I'll publish on my own, these REQs are for things that still need to implemented by ai-coders." The proposal asked for "no REQ, a report-only line, and an abandon recommendation". The maintainer rejected that shape and said: "operator action needs to show up in the Needs input · Blocked column, when the code is ready to be released (but not before, before it should have dependencies that might block it, in which case it should be in pending column)"; "when the operator section is done, the REQ should be moved to the DONE column"; and "don't open operator tasks everytime, I'll deploy on my own time, that is nothing special, unless special configuration needs to be done (again this is prose that mostly this orchestrator skill should not concern itself too much)".
## Verified Facts (from triage)
- `skills/do-work/actions/capture.md:107` (Earmark assessment) routes "I need to be at the keyboard for this one" to a `blocked` capture; `:108` (External-condition assessment) defines `blocked` as a condition the work cannot start without, with examples that are all preconditions to AI work. Neither says what to do when the operator act is the whole remaining work.
- `skills/do-work/actions/work.md:523` (Error Handling row added by REQ-648) flips a queued REQ the orchestrator wants to park on the operator to `blocked`. `skills/do-work/actions/work-reference.md:788` (Failure Classification, Environment row) holds the blocked-flip test with no branch for "the missing thing is the whole remaining work".
- `skills/do-work-board/tools/queue-kanban/model.go:1723` places `status: blocked` with unmet `depends_on` under Pending → Waiting; `:1735` places `blocked` with none under Needs input · Blocked; `completed` goes to Done. `model_test.go:638-647` already pins blocked-with-unmet-deps → Waiting, so no board change is needed.
- `skills/do-work/actions/clarify.md:179-190`: the plain-blocked prompt offers "1. Yes — unblock it  2. Not yet — leave it  3. Abandon this REQ"; "Yes" runs the `unblock` transaction, which returns the REQ to `pending` and the next run claims it. `clarify.md:267` already flips a confirmed `builder_decided: true` follow-up to `completed` in place, and `actions/cleanup.md` Pass 0 archives terminal REQs left in the queue.
- The core contract in `skills/do-work/SKILL.md` ("Capture does not execute": a UR plus one or more REQs) and `capture.md` Philosophy are untouched by this rule: a request whose only content is a routine operator act still gets its UR, with zero REQs only when nothing AI-buildable is in it, which `capture.md:16` already allows as the fold-only shape; say so in the capture summary.
## Detailed Requirements
1. `actions/capture.md:108`, External-condition assessment, two sentences keyed on the condition, not a list: a routine operator act after the code ships (deploy, publish, approve, verify on live hosts are examples, not a list to match) is never captured as a REQ; the user does it on their own time, and only the AI-buildable parts become REQs, with the capture summary naming the operator step left to the user. When the request itself asks to track a special operator configuration, capture that one REQ as `status: blocked`, `blocked_by: '<the operator, in the user's words>'`, `blocked_at`, and `depends_on` listing every AI-buildable REQ minted from the same request, with the checklist or runbook in the REQ body; it waits under Pending → Waiting until those complete, then surfaces under Needs input · Blocked.
2. `actions/capture.md:107`, Earmark assessment: the "at the keyboard" sentence points at the new rule instead of flatly capturing `blocked`.
3. `actions/clarify.md` Step 5.5, plain-blocked prompt: add `4. Done — I did it myself`. On that answer clarify appends a `## Operator receipts` section holding the user's words verbatim (apply Step 4's Outside-text containment), sets `status: completed` and `status_changed_at: <now>` in place, and does not run `unblock`. The summary in Step 6 lists each REQ completed this way. Precedent to cite: the confirmed `builder_decided: true` path at `clarify.md:267`, and `actions/cleanup.md` Pass 0 for the archive move. No new CLI verb.
4. `actions/work.md:523` row and `actions/work-reference.md:788` Environment row, one sentence each: when a run finds that the remaining work is a routine operator act, the REQ completes with what the AI built and the hand-back names the operator step; it is not flipped to `blocked`, not failed, and abandon is not recommended. The blocked flip stays for a precondition to AI work and for a tracked special configuration.
5. `docs/work-guide.md`, the Earmarking or blocked paragraph: one sentence on the shape and where it shows on the board (Pending → Waiting, then Needs input · Blocked, then Done).
6. Acceptance greps: `grep -n "own time" skills/do-work/actions/capture.md` finds the routine-act sentence; `grep -n "Done — I did it myself" skills/do-work/actions/clarify.md` finds the option; `grep -n "operator" skills/do-work/actions/work.md skills/do-work/actions/work-reference.md` finds the two inverse sentences.
7. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` green. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Prose only. No new field, status, flag, transaction, probe or board change. `blocked`, `depends_on`, `unblock` and the default scan keep their meaning.
- Keep it short. The maintainer said the orchestrator skill "should not concern itself too much" with operator work: two sentences in capture, one option in clarify, one sentence each in the two run rows, one in the guide. Delete before you add where an existing sentence already says part of it.
- Record the challenge in the Decisions of this REQ: the proposal's "no REQ, report-only line, abandon" shape was rejected by the maintainer in favour of the blocked + `depends_on` placement the board already has and a clarify completion path.
- Read `_dev/primes/prime-action-files.md` before editing; cross-references use that prime's spelling.
- Batch constraint: this REQ is its own release and its own commit; REQ-653 depends on it and must not be folded into it.
## Dependencies
None upstream. REQ-653 (fan-out prose to a companion reference) depends on this REQ because both edit `actions/work.md` and `actions/work-reference.md`.
## Builder Guidance
Certainty is high on the shape (the maintainer stated each column and the Done move). Latitude: exact wording, and whether the clarify Done path writes the receipts section before or after the status line.
## Red-Green Proof
**RED prompt/case:** A capture whose request ends "then deploy it to production and verify on the live host" and a clarify session on a `blocked` REQ whose operator work is finished.
**Why RED now:** Capture mints a `blocked` operator REQ (capture.md:107-108) with no `depends_on`; clarify offers only unblock, leave, or abandon, so the finished REQ goes back to `pending`, the run claims it, finds nothing to build, and flips it to `blocked` again (work.md:523).
**GREEN when:** Capture writes no REQ for the routine deploy and says so in its summary; a request that asks to track a special configuration yields one `blocked` REQ with `depends_on` on the AI-buildable REQs; clarify's prompt has the Done option and leaves a `completed` REQ with `## Operator receipts`, which the board shows under Done and `cleanup` archives; the acceptance greps in requirement 6 all hit; both contract tests are green.
**Validation:** User adjusted. The maintainer replaced the proposal's shape with the blocked, dependency-gated, Done-column shape in two messages, quoted verbatim in the UR.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (6855 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` (a rule added in one action stales sibling restatements) is the risk when capture, clarify, work and the guide each carry part of this rule.
## Full Context
See `do-work/user-requests/UR-143/input.md` for complete verbatim input (the full proposal, the maintainer's answers and refinements, and the triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream proposal "operator work is not a REQ" (verified by `do-work-toolbox validate-feedback` on 2026-10-09 as F11), reshaped by the maintainer: "operator action needs to show up in the Needs input · Blocked column, when the code is ready to be released …", "when the operator section is done, the REQ should be moved to the DONE column", "don't open operator tasks everytime, I'll deploy on my own time … unless special configuration needs to be done".*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the five files and the exact lines, and states each sentence to add; prose only, no Go, no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
