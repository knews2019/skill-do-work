---
id: REQ-314
title: "[impact-rule-change] Judge effort_estimate on review-minted follow-ups too"
status: completed
created_at: 2026-08-21T08:56:44Z
status_changed_at: 2026-08-21T18:51:21Z
claimed_at: 2026-08-21T18:20:11Z
completed_at: 2026-08-21T18:51:21Z
kb_status: pending
route: B
user_request: UR-064
addendum_to: REQ-308
domain: general
review_generated: true
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work-board/tools/queue-kanban/model.go
  - _dev/tests/contract-regressions.sh
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-21T18:20:33Z
  basis:
    - trivial short-circuit
---

# Judge effort_estimate on Review-Minted Follow-Ups Too

## What

REQ-308 made capture judge `effort_estimate` on every REQ it mints, by the same three-way contract
`impact:` already carried. The other writer of new REQs kept the weaker rule.
`actions/work-reference.md` → **Discovered Tasks Classification (Step 8)** and
`actions/review-work.md` Step 10 both tell that path to "write `effort-mechanical` only when you
have actually judged the fix small, and otherwise leave it absent to read as `effort-substantive`".

Half of that is right and stays: never invent `effort-mechanical`. The other half is permission not
to judge, and it is the rule capture just lost.

## Why

Review-minted follow-ups are a large share of the queue — this `do-work run` alone created three
(REQ-312, REQ-313, and this one) — so they are a large share of what `do-work run-simple-reqs`
cannot see. The measurement REQ-308 was built on was 14 of 22 pending REQs carrying the field; the
missing eight are exactly this population.

Leaving the two writers on different standards also reintroduces the asymmetry REQ-308 removed, one
level down: a reader who lands on Step 8's rule learns that leaving the field absent is fine.

## Requirements

- The review/discovered-task follow-up path judges `effort_estimate` by the same three-way contract
  capture now uses: judge it, or put the judgment to the user, or leave it absent because neither was
  possible — never a copied default, in either direction.
- Every site stating the weaker rule is updated, not just the one this REQ names. REQ-308's sweep
  found a fifth site its own REQ had not listed; run the same sweep.
- The two axes stay independent: `effort_estimate` is never derived from the finding's `impact:`
  token. That sentence already exists at both sites and is being enforced, not restated.
- The lock-in check pins the property, the way REQ-308's does — ideally by extending that check to
  cover this writer, rather than adding a second one that can drift from it.
- No backfill of existing REQs, and no enum growth. Both exclusions carry over from REQ-308 unchanged.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** Read `actions/work-reference.md` → **Discovered Tasks Classification (Step 8)**
and `actions/review-work.md` Step 10. Both still say to leave `effort_estimate` absent when it is
not obviously mechanical, which is permission not to judge.
**GREEN when:** A check fails on a follow-up-minting rule that permits an unjudged
`effort_estimate`, and passes once both sites carry capture's contract.
**Validation:** Discovered task from REQ-308; apply `actions/work-reference.md` →
**Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-308: that REQ made request capture
  judge the size field on every new request instead of leaving it blank, because blank quietly reads
  as "big" and hides the request from the cheaper-model queue verb. But requests are also created a
  second way — automatically, when a code review finds something worth following up — and that path
  still says leaving the field blank is fine. Those follow-ups are a large share of the queue, so
  the problem REQ-308 fixed is mostly still there. Fixing it means applying the same rule one writer
  over. Should I process this as a new task? → Confirmed: Yes, add to queue

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  **Answered 2026-08-21** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the builder's recommendation via
  `do-work clarify`: apply capture's explicit effort judgment to review-created follow-ups while
  retaining the REQ's existing exclusions. No additional scope was requested.

  Why this is yours rather than mine: REQ-308 deliberately scoped itself to capture, and you were
  the one who filed it that way rather than folding it into the pull request that discovered it. So
  the question of whether the same rule should reach the automatic path is a scope call you have
  already made once, in the other direction, and it should be yours to make again rather than mine
  to assume.

---

## Triage

**Route: B** — Medium

**Reasoning:** The prose edits are mechanically small, but the requirement is a semantic sweep over
every follow-up-minting restatement plus a property-level regression extension. Route B provides an
independent exploration pass to prove the population and avoid shipping a second half-landed rule.

**Planning:** Not required for Route B; REQ-308 supplies the canonical three-way contract and the
lock-in pattern to extend.

**Estimate:** 5 active minutes (P50, high confidence; effort-mechanical short-circuit).

## Exploration

The semantic population is five weak or stale statements across three shipped files, plus the
existing REQ-308 regression block. `review-work.md` permits mechanical-or-omit and its emitted
frontmatter template omits the field; `work-reference.md` repeats that permission in Discovered
Tasks and still describes capture as the only judging writer in both schema surfaces; `model.go`
mirrors the capture-only schema claim and explicitly requires lock-step correction.

The canonical replacement keeps impact independence and requires the same three-way result at both
automatic follow-up writers: judge size and emit either canonical value; if genuinely unclear put
the size judgment to the user; omit only when neither judging nor asking was possible; never copy a
default. `docs/review-work-guide.md` already says effort is a separate judgment and stays unchanged.
`work.md` names required review-created follow-up fields without delegating the effort judgment, so
its Step 7 restatement joins the write set with a concise pointer to Step 10.

Extend REQ-308's property-level Python check rather than adding a second contract. The extension
compares both follow-up writer sections to capture's canonical alternatives and requires the Step 10
frontmatter template to carry an effort placeholder. Mutations must catch either weak writer alone,
both weakened together, a missing template field, and judged non-small being replaced by omission.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — delegate review-minted follow-up effort judgment to
  `review-work.md` Step 10 without duplicating the three-way contract
- `skills/do-work/actions/review-work.md` (modify) — require explicit effort judgment at Step 10 and
  carry the result in the emitted follow-up template
- `skills/do-work/actions/work-reference.md` (modify) — apply the same contract to Discovered Tasks
  and generalize the schema/read-contract writer descriptions
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — keep its schema-mirror comment in
  lock-step; parser behavior stays byte-for-byte unchanged
- `_dev/tests/contract-regressions.sh` (modify) — extend the existing REQ-308 property check across
  both automatic writers and the review template, with adversarial mutations

**Files I will NOT touch:** capture files (their contract is already canonical), guides,
selectors, enum definitions, existing REQs, or any runtime parser behavior.

**Acceptance criteria (restated from REQ):**
- [x] Review and Discovered-Task follow-up writers judge and emit either effort value, ask when
  genuinely unclear, and omit only when neither judgment nor a question was possible.
- [x] The two axes remain independent and no default is copied in either direction.
- [x] Every stale capture-only/mechanical-or-omit restatement in the measured class is reconciled.
- [x] REQ-308's existing property-level check covers both writers and the Step 10 template and is
  mutation-proven against independent and symmetric weakening.
- [x] No backfill, enum growth, selector change, or parser behavior change occurs.
- [x] The direct canonical repository gate exits 0.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/`; this REQ's own claim/answer bookkeeping is isolated.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh`
**Dependencies:** ✓ Python, Bash, Go 1.26.1, and ShellCheck are available.

*Checked by work action*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the action prime and REQ-308 lessons; extend its canonical three-way property
  check across the two automatic writers, template, schema surfaces, and parser mirror.
- [x] **[APPLY]:** Extended REQ-308's property-level contract first, captured RED against both
  weak follow-up writers, then applied the three-way size judgment to review Step 10, its emitted
  template, Discovered Tasks, the schema/read-contract surfaces, and the Go schema mirror. Review
  remediation added `work.md`'s delegation and anchored every semantic check at its own effort
  directive.
- [x] **[UNIFY]:** Reviewed the exact five-file diff and mutation-proved independent and symmetric
  writer weakening, template-field deletion, non-small omission, all three schema restatements,
  both ask-when-unclear clauses, and review's never-copy-default clause. Focused contract, `bash -n`,
  ShellCheck (existing informational findings only), `gofmt -l`, `go vet ./...`, scope checks, and
  `git diff --check` are clean.

## Decisions

- **D-01 — DECIDE & STATE: test writer instructions separately from emitted templates.** The first
  mutation replay exposed that Step 10's template comment could satisfy a section-wide prose check
  after the actual writer instruction was weakened. The property checker now strips fenced payloads,
  isolates the paragraph carrying the writer's `effort_estimate` directive, and starts at that field
  before judging its semantic clauses; the frontmatter template is checked independently at the
  emission seam.
- **D-02 — DECIDE & STATE: emit both canonical effort values explicitly.** Review and Discovered
  Tasks now name both `effort-mechanical` and `effort-substantive`; omission is reserved for the
  same neither-judgment-nor-question escape hatch as capture. This prevents the old default from
  standing in for a substantive judgment while preserving impact independence.
- **D-03 — DECIDE & STATE: extend scope to `actions/work.md`.** Review found that Step 7 restated
  required review-created fields but did not delegate `effort_estimate` to the newly canonical
  writer. A one-sentence pointer closes that reader without copying the three-way contract; the
  file was added to `write_set` and Scope before modification.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Both automatic follow-up writers now judge size and emit the matching canonical
effort token, ask the user when size is genuinely unclear, and omit only when neither judging nor
asking was possible. The schema and board mirror name every writer, `work.md` delegates its
review-follow-up field list to the canonical Step 10 judgment, and REQ-308's existing semantic
contract independently verifies both field-anchored instructions and the emitted template.

## Qualification

Passed — 5 implementation files verified, 6 acceptance criteria traced, P-A-U confirmed. Mechanical
qualification and Scope/Implementation Summary parity pass; the Go diff is comment-only, all three
action edits are substantive, and the regression extension reaches each writer and the emitted
template.

## Testing

**Tests run:** `bash -n _dev/tests/contract-regressions.sh`; ShellCheck 0.11.0; `gofmt -l model.go`;
`go vet ./...`; `bash _dev/tests/contract-regressions.sh`; scope-drift and qualification checks.
**Result:** ✓ Remediated focused contract exited 0. ShellCheck reported only existing informational
findings outside the changed block. The orchestrator owns the post-review canonical repository gate.

**Red-green validation:**
- Extended REQ-308 property contract: ✗ before prose changes — named the weak review writer and
  missing omission boundary, with Discovered Tasks weak too → ✓ after both writers, template, and
  schema mirrors adopted the canonical contract.
- Eleven restored mutations failed: weak review alone; weak Discovered Tasks alone; both weakened;
  missing template field; judged non-small replaced by omission; and removal of each of the schema,
  read-contract, and Go-mirror writer populations; plus deletion of review's ask-when-unclear clause,
  deletion of Discovered Tasks' ask-when-unclear clause, and deletion of review's
  never-copy-default clause.
- Fenced template payloads are stripped from instruction checks and tested separately, so template
  wording cannot mask a weakened writer directive; unrelated wording earlier in the directive's
  paragraph cannot satisfy the field-anchored clauses either.

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` (REQ-308): expanded its existing property-level effort
  judgment contract from capture to every current new-REQ writer.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-21T18:50:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None. The first review found that loose section-wide tokens let unrelated prose mask deletion of
required effort clauses; remediation anchored the check to each writer's own directive and added
the three missing deletion mutations.

**Minor findings:** None after remediation.
**Acceptance:** Pass — every current REQ-minting writer now carries the same effort judgment, and
the regression contract independently proves the writer clauses and emitted template.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Comparing every REQ writer with one canonical semantic property keeps capture, review follow-ups,
  and Discovered Tasks aligned without adding another competing contract.
- Mutation-testing each writer independently caught both one-sided and symmetric weakening, while a
  separate template assertion kept emission coverage at the actual output seam.

**What didn't:**
- A section-wide keyword scan initially looked semantic but allowed unrelated Step 10 prose to
  satisfy the ask-user and never-copy-default legs. Fenced payload isolation alone was insufficient;
  the check also had to isolate the paragraph containing the actual `effort_estimate` directive.

**Worth knowing:** Schema prose and its Go mirror are active contract surfaces even when parser
behavior does not change. When a writer population changes, those read-side descriptions must move
in the same commit or they will keep teaching the old boundary.

## Orientation

**[MAP CHANGED]** Effort judgment is now a property of every path that mints a REQ, not a capture-only
behavior. Capture, review-generated follow-ups, and Discovered Tasks all judge and emit one of the
two canonical effort tokens, ask when genuinely unclear, and reserve omission for cases where
neither judging nor asking was possible. The review action owns its emitted-template details; the
work orchestrator delegates there rather than restating the rule.
