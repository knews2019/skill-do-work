---
id: REQ-372
title: "[impact-rule-change] State the builder's out-of-scope response once"
status: completed
claimed_at: 2026-08-25T08:29:59Z
completed_at: 2026-08-25T08:36:51Z
commit: a2e28d3
route: A
created_at: 2026-08-24T23:37:06Z
user_request: UR-073
addendum_to: REQ-365
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
impact: impact-rule-change
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-25T08:30:49Z
  basis:
    - trivial short-circuit
write_set:
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/actions/work.md
---

# State the Builder's Out-of-Scope Response Once

## What

REQ-365's shipped invariant tells a builder to flag the contradiction and **proceed**;
`actions/work.md` tells a builder discovering an out-of-scope file to **stop and report to the
orchestrator**. Both are shipped, both describe the same moment, and they prescribe different
actions. State the response once, in the file the builder actually loads, and cite it from the other.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the action-file prime, prior implementation, and crew rules. Keep the
  conditional-completeness invariant at capture, move the complete builder response to `work.md`,
  and distinguish a contradictory declaration from a genuinely new scope expansion.
- [x] **[APPLY]:** Centralized the two-path builder response in `work.md` and replaced the
  capture-side duplicate with a section pointer. No files outside the declared set were changed.
- [x] **[UNIFY]:** Reviewed the two-file implementation diff and `git diff --check`; confirmed the
  capture-side response was removed, the work-side two-path rule is canonical, and no debug
  artifacts landed. The full maintainer verification gate passed.

## Why

The two texts diverge on what the builder does next:

- `actions/capture-reference.md` → **Populating `write_set`**, conditional completeness invariant:
  "If a builder meets this contradiction in an already-queued REQ, flag it before editing, **proceed
  with the required file class despite the declared set**, and report the contradiction plus the
  actual files in the handback."
- `actions/work.md` § **Write only inside the declared scope**: "Discovering mid-build that you need
  a file outside it is a **stop-and-report to the orchestrator, never a silent write.** The
  orchestrator records the request and its resolution in the REQ trail as a `## Decisions` D-XX entry
  (it is a scope judgment) and extends both the Scope list and `write_set`."

Proceed-then-report and stop-then-wait are not the same instruction. A contradiction between two
shipped instructions changes what an agent does, which is why it earns a REQ rather than a prose
note (`actions/capture-reference.md` → Fold-First Rule → the prose-only test's second exemption).

The divergence is also in the wrong place: builders load `actions/work.md`, not the capture
reference. A builder-behaviour rule that only capture reads is a rule with no reader.

## Prior Implementation

REQ-365 (`addendum_to`) is completed and archived: `do-work/archive/REQ-365-a-tdd-req-must-name-a-test-file-in-its-write-set.md`, Route A, `commit: 6265f1c406ccfe72f4cb12f1e0376c483eac9b65` (merge "Merge REQ-365 write-set completeness guidance"; the implementation itself is `b7d2362`, "docs: require complete declared write sets").

**What shipped — two files, three lines:**

- `skills/do-work/actions/capture-reference.md` gained the **Conditional completeness invariant**
  paragraph immediately after **Populating `write_set`**: a declared set must name at least one path
  from any file class the REQ's own requirements or completion proof require writing, or omit
  `write_set` entirely; `tdd: true` is named as today's instance rather than a special case. Its last
  two sentences are the builder-response clause this REQ relocates, plus a restatement that the field
  stays display-only.
- `skills/do-work/actions/capture.md` gained **one trailing sentence** on the Step 1 TDD-assessment
  bullet, pointing at that invariant. It is a pointer, not a second copy of the rule.

**Pattern to preserve:** the rule lives once at its canonical home and the decision site cites it.
That is why this REQ moves a clause rather than adding a third statement, and why `capture.md` needs
no edit here — its pointer already reaches whatever the invariant says.

**Not shipped, deliberately:** no gate, no new frontmatter field, no lock-in test. Nothing under
`_dev/tests/` asserts the wording of either paragraph, so this REQ's edit needs no test update.

## Detailed Requirements

- Decide which behaviour is intended and state it **once**, in `actions/work.md` § "Write only inside
  the declared scope". The proceed-anyway case is a real one — REQ-346's builder proceeded, the
  orchestrator upheld it, and an unattended fan-out builder has no orchestrator to wait on — so if
  that is the intended response, it belongs in that bullet as a stated condition, not as a second
  rule elsewhere.
- Reduce the invariant's sentence in `actions/capture-reference.md` to a citation of that bullet.
  The invariant itself stays where it is: it is a capture-time rule about what a declared set must
  contain, and only its builder-response clause is being relocated.
- Key the surviving statement on the condition (a needed file outside the declared set), not on
  tests or on `tdd: true` — the invariant already establishes that tests are one instance
  (`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale).
- Keep both files' existing guarantee intact: `write_set` stays display-only, and an honest absent
  set stays better than an invented one.

## Builder Guidance

Certainty: **Firm** on the problem, **open** on which behaviour wins. The contradiction is verified in
both files; which of the two responses is correct is a judgment the builder makes and states.

Scope cue: keep it small. One statement moves, one pointer replaces it. Do not reconcile the two by
writing a third rule, and do not reword the invariant's capture-time half while you are in there.

Batch (UR-073): this REQ and REQ-373 overlap on `skills/do-work/actions/work.md`, and REQ-373 may
also need `skills/do-work/actions/capture-reference.md` depending on how its Open Question resolves.
Run this one first — it edits the same invariant paragraph REQ-373 might have to amend.

## Constraints

- `_dev/primes/prime-action-files.md` governs. Read it first.
- Delete before you add: the fix is one statement plus one pointer, not a third statement
  reconciling the two.
- Do not touch the invariant's capture-time half, and do not make `write_set` a gate
  (`actions/work-reference.md` → Request File Schema).

## Red-Green Proof

**RED prompt/case:**

```bash
grep -n "proceed with the required file class" skills/do-work/actions/capture-reference.md
grep -n "stop-and-report to the orchestrator" skills/do-work/actions/work.md
```

Both return a line today, each prescribing a different next action for the same moment.

**Why RED now:** REQ-365 added the builder-response clause to `capture-reference.md` without
reconciling it against the declared-scope bullet that already governed builders.

**GREEN when:** exactly one of the two files states what the builder does, that file is
`actions/work.md`, the other cites it by section name, and the stated response covers the unattended
case that made REQ-346's builder proceed.

**Validation:** Inferred during capture.

---
*Source: UR-073 finding F1 — reconciling this branch's capture against REQ-365 as shipped on main (`6265f1c`).*

---

## Triage

**Route: A** - Simple

**Reasoning:** The contradiction and both target files are already identified. The change is a
focused instruction move plus a canonical pointer, with no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — Split the response by whether the REQ already requires the file class.** A required but
  omitted class is contradictory metadata, so the builder may flag, proceed, and report it. A file
  that the REQ did not already require is a real scope expansion and remains stop-and-report. This
  preserves unattended builders without turning `write_set` into permission for unrelated edits.

## Implementation Summary

- `skills/do-work/actions/work.md` (modified): states the canonical two-path response for a file
  outside declared scope, including the unattended/worktree handback case.
- `skills/do-work/actions/capture-reference.md` (modified): replaces the duplicate response with a
  pointer to the canonical work-action rule while retaining the capture-time invariant.

**What was done:** Builder behavior now has one owner. A requirement-backed omission can proceed
with a recorded contradiction; a genuinely new file remains a stop-and-report scope expansion.

## Discovered Tasks

None.

## Testing

- **Focused checks:** The captured RED grep no longer finds the proceed instruction in
  `capture-reference.md`; semantic searches find the flag/proceed branch and stop/report branch in
  `work.md`, plus the single canonical pointer at capture.
- **Regression:** `git diff --check` passed.
- **Canonical gate:** `bash _dev/tests/maintainer-verify.sh` passed all contract, Go, strict
  JavaScript, and audit-metrics lanes. The optional strict browser lane skipped because no browser
  is configured.

## Qualification

- Mechanical qualification passed with the two implementation files present in the diff and all
  P-A-U phases completed.
- The changes are substantive and trace every detailed requirement: one canonical response,
  condition-based branching, an explicit unattended-builder path, and preserved display-only
  `write_set` semantics.
- Route A correctly has no `## Scope`; its two changed files match the capture-authored
  `write_set`, and the implementation contains no hollow or unwired surface.

## Review

**Overall: 100%** | 2026-08-25T08:36:28Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None

**Minor findings:** 0
**Acceptance:** Pass — the two-path rule covers serial and unattended builders, the capture copy is
gone, and the display-only contract is unchanged.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

The restatement sweep checked live action, reference, crew, guide, and schema surfaces. No live text
still prescribes the retired unconditional response; historical changelog entries remain history.

*Reviewed by review-work action*

## Orientation

Builders now have one canonical out-of-scope rule: requirement-backed declaration mistakes may
proceed with an explicit report, while genuinely new scope still stops for orchestrator approval.
