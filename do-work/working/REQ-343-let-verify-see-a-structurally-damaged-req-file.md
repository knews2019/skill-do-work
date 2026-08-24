---
id: REQ-343
title: "Let verify see a structurally damaged REQ file"
status: claimed
claimed_at: 2026-08-24T08:55:00Z
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-342, REQ-344]
maintenance: false
route: B
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
---

# Let Verify See a Structurally Damaged REQ File

## What

`queue-kanban verify` reports `OK: no findings` and exits 0 on REQ files whose structure is broken.
Give it a structural-anomaly probe, and lift the unrecognized-status warnings the board already
produces into findings the same way the duplicate-id probe does.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

On the user's fixture, six of seven REQs carried delimiter damage — including one whose opening
frontmatter fence was broken so its `status`, `title` and `user_request` all parsed empty — and verify
printed `OK: no findings` and exited 0. Verify is the mechanical half of the pre-commit ritual, so a
clean exit there is what an operator trusts before committing. It is also the safety net for the
damage REQ-342 and REQ-344 exist to prevent: without this probe, that damage is silent twice over.

The mechanism is the parser's leniency, which is correct on its own terms: recovery means the damage
surfaces as *empty fields* rather than a parse error, so nothing throws. `buildBoard`'s
unrecognized-status warning is the only trace it leaves.

## Context

Verified against the source. `collectVerifyFindings` (`verify.go:141-158`) runs thirteen `append*`
probes. Exactly one of them reads `board.Warnings`: `appendDuplicateRequestIdFindings`
(`verify.go:282-292`), and it filters by `duplicateRequestIdWarningPrefix`, so every other warning
class — unrecognized status included — is never lifted. `ExitCode` (`verify.go:104`) keys on
`len(report.Findings)`, so an unlifted warning cannot affect the exit status.

`appendCompletionAnomalyFindings` is the pattern to follow rather than invent: its own comment
records that "verify was blind to every anomaly class until then" (REQ-214), and it forwards the
board's structured evidence instead of re-walking the tree or parsing warning prose.

**Capture decision — the missing-`user_request` class was narrowed, and the user did not ask for
that.** The request lists "a missing `user_request` pointer" as an anomaly the probe should fail on.
Taken literally that fires on every stakeholder-questions REQ, which carries no `user_request` **by
design** (`actions/work-reference.md` → Stakeholder REQ Template), and on every REQ in
`archive/legacy/`, which predates the field. A probe that flags correct files is a probe someone
turns off, so the requirement below carves those two out. **Value:** the probe stays trustworthy, so
its findings keep meaning something. **Risk:** a genuinely damaged REQ that happens to look like a
stakeholder REQ would slip through — narrow, because the carve-out keys on the documented shape and
not on the field's mere absence. Reversible: delete the carve-out and the literal reading is back.

## Detailed Requirements

- A structural-anomaly probe fails the mechanical check on a REQ file with any of: no leading
  frontmatter fence, an empty or unrecognized `status`, an empty `id`, or a missing `user_request`
  pointer.
- Existing unrecognized-status warnings are lifted into findings the same way
  `appendDuplicateRequestIdFindings` lifts duplicate-id warnings.
- Forward the board's structured evidence; do not parse warning prose and do not re-walk the tree —
  the same rule `appendCompletionAnomalyFindings` states for itself.
- Each finding names the broken field and its remedy, so the operator can act without opening the
  tool's source.
- The probe distinguishes damage from legitimate absence: a stakeholder-questions REQ and a REQ in
  `archive/legacy/` both legitimately carry no `user_request`, and neither is a finding. This narrows
  the request's literal wording — see the capture decision in `## Context` for why, and overturn it
  there if the narrowing is unwanted.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first, including the parser
  lock-step convention.
- **Keep the parser's leniency.** The point is to report the damage, not to start rejecting files:
  a REQ with one bad line must still parse and still appear on the board.
- Do not weaken `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` or any existing verify
  probe, and do not change what the board renders — this REQ adds detection, not display.

## Builder Guidance

**Certainty: firm — the gap, the mechanism and the pattern to copy are all confirmed in the source.**
The one real judgment is the legitimate-absence carve-out above; get that wrong and verify cries wolf
on every stakeholder REQ, which is how a probe gets disabled.

Write the fixture before the probe. The user's shape — several REQs damaged, at least one with a
broken opening fence — is the right RED, and it is also the fixture that proves the carve-out, so
include a stakeholder REQ and a legacy REQ in it that must NOT be flagged.

## Open Questions

None — the user named the four anomaly classes and the lifting pattern.

## Red-Green Proof

**RED prompt/case:** Build a fixture repo whose `do-work/queue/` holds REQ files with each of the four
damage shapes, plus a stakeholder REQ and a legacy REQ that are structurally fine. Run
`queue-kanban verify --repo-root <fixture>`: it prints `OK: no findings` and exits 0.

**Why RED now:** Only duplicate-id warnings are lifted into findings, and `ExitCode` keys on the
finding count, so an unrecognized or empty status cannot fail the check.

**GREEN when:** The same fixture produces one finding per damaged REQ, each naming the broken field,
verify exits nonzero, and the stakeholder and legacy REQs produce no finding. A healthy fixture still
exits 0 with `OK: no findings`.

**Validation:** User confirmed — the fixture, the counts, the mechanism and the lifting pattern are
stated verbatim in the input, and each was re-verified against the source during capture.

## Assets

None. The fixture is the deliverable's own RED.

---
*Source: UR-068 — see `do-work/user-requests/UR-068/input.md` for complete verbatim input.*
