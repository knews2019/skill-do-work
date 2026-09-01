---
id: REQ-310
title: "[impact-rule-change] Check a template payload's citations against where the payload lands"
status: completed
created_at: 2026-08-20T23:36:10Z
status_changed_at: 2026-08-21T00:00:00Z
claimed_at: 2026-08-21T16:25:24Z
completed_at: 2026-08-21T16:51:29Z
commit: 61cfc28
user_request: UR-055
addendum_to: REQ-269
domain: general
route: B
impact: impact-rule-change
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-21T16:34:12Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
depends_on: []
maintenance: false
kb_status: promoted
kb_entry: REQ-310-check-a-template-payload-s-citations-aga.md
write_set:
- skills/do-work-toolbox/actions/code-review.md
- skills/do-work-toolbox/actions/validate-feedback.md
- _dev/tests/contract-regressions.sh
---

# Check a Template Payload's Citations Against Where the Payload Lands

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the request, parent, project rules, prime, and live producer/test seams;
  measured the affected payload population and fixed a three-file scope before implementation.
- [x] **[APPLY]:** Captured the four-token RED surface, repaired the two payloads, observed the
  existing tests reject their stale assumptions, then narrowed the existing regression assertion.
- [x] **[UNIFY]:** Reviewed all three diffs; verified the four old tokens are absent and four new
  destinations present; ran both citation contracts, shell syntax/lint, and whitespace checks.

## What

REQ-269 split the fence exemption by its reason: a fenced payload is exempt because its text lands in *another* file, so it is not a citation from the file it is written in. That reasoning is sound and the exemption survives on it — but it only establishes that the payload is not checkable **from here**. It says nothing about whether the payload resolves **there**, and nothing checks that today.

A live example, found while confirming REQ-269's scope-drift was intended:

`skills/do-work-toolbox/actions/code-review.md:328` sits inside a ```markdown block that is a REQ template a reviewer copies into a new REQ file. It reads:

```
**Validation:** Review finding; apply `../do-work/actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
```

That payload lands in `do-work/queue/REQ-NNN-....md` in a consumer project. From there, `../do-work/actions/work-reference.md` resolves to `do-work/actions/work-reference.md`, which is not where the installed action lives — the installed path is under `.claude/skills/do-work/actions/`. So the citation is wrong at its destination, and both the pre-REQ-269 checker and the post-REQ-269 checker are structurally blind to it, for the same good reason.

`skills/do-work-toolbox/actions/validate-feedback.md:114` carries the same shape.

## Context

Discovered by REQ-269's orchestrator (REQ-269 `## Decisions` D-05) while verifying that its two declared-but-untouched files were correctly untouched. They were — this is the *next* question, and it is genuinely a different one: REQ-269 closed "is this token a citation from here", and this is "does a payload citation resolve where the payload goes".

Naming the shape honestly: this is the fourth marker. REQ-269 retired the backtick, the leading `../`, and the fence. Each retirement revealed the next boundary. This one is a boundary of a different kind — not a spelling standing in for a meaning, but a genuine gap in coverage that the correct exemption creates.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-269: a fenced template's payload is exempt from citation checking because it lands in another file, and nothing verifies it resolves at that destination. Two live templates carry a citation that is wrong where it lands. Should I process this as a new task? → Confirmed: Yes, add to queue
  *(2026-08-21)* User confirmed via `do-work clarify`: process as a new task. Nothing put out of scope — which of the three sizing options (hand-fix only, destination convention, installed-topology-only rule) to take is still open and belongs to whoever builds this REQ.

**Worth knowing before you decide**, because it changes the size a lot:

Checking a payload at its destination needs the destination to be **declared**, and today it is not. A fenced block does not say "this lands in a consumer REQ file" — a human infers it. So the honest options are:

1. **Fix the two sites, add nothing.** Smallest. Leaves the class open, which is exactly the posture REQ-269 spent its whole scope arguing against.
2. **Declare the destination where a template block is written** — an info-string convention (` ```markdown lands=do-work/queue/ `) or a marker comment — and have the checker resolve payload citations from there. Closes the class, costs a convention every template author must learn, and the checker can only check blocks that opted in.
3. **Treat installed-topology paths as the only correct form inside a template payload** and check that instead. No new convention, but it asserts a rule about template content that may not be true for every template.

Per `crew-members/maintenance.md` the deletion questions were asked: option 1 is the removal-shaped answer and it is genuinely on the table, because two hand-fixed sites may be the whole population. Counting the live template blocks that contain package-shaped paths is the first thing to do if you say yes — if it is two, option 1 is right and this REQ should close as a two-line fix.

---

## Triage

**Route: B** - Medium

**Reasoning:** The population measurement found exactly two destination-bound Markdown template
blocks with package-shaped citations, both in the toolbox actions already named by the request.
That selects the request's hand-fix option and avoids adding a destination convention or checker,
but the two blocks land in different contexts and still require destination-specific judgment.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The shipped corpus contains nine core-action-bearing Markdown payload blocks. Seven are core
templates that already use destination-safe skill-local `actions/...` citations. Exactly two are
affected sibling-package blocks carrying source-relative core-action citations:

- `skills/do-work-toolbox/actions/code-review.md` — a REQ template written under
  `do-work/queue/`, with two source-relative do-work action citations inside its payload.
- `skills/do-work-toolbox/actions/validate-feedback.md` — a user-facing capture handoff with two
  do-work action citations in one line.

All other fenced package paths are executable examples or schema annotations, not affected
destination-bound Markdown payload citations. The affected count selects option 1 from the approved
request: repair the two blocks and add no new declaration syntax or checker surface. An existing
contract-regression assertion must also be corrected because it currently requires the stale
source-relative spelling in toolbox-generated REQs.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/code-review.md` (modify) — make the generated REQ payload's
  citations resolve from its destination rather than from the toolbox source file
- `skills/do-work-toolbox/actions/validate-feedback.md` (modify) — make the generated capture
  handoff's citations resolve in the consumer project rather than from the toolbox source file
- `_dev/tests/contract-regressions.sh` (modify) — replace the stale source-package citation
  assertion with the destination-safe REQ-payload convention it should enforce

**Files I will NOT touch:** the shipped-package reference checker and prime files. The existing
regression assertion is adjusted rather than adding a new destination checker or opt-in convention.

**Acceptance criteria (restated from REQ):**
- [x] Every do-work action citation in the two destination-bound Markdown payloads resolves in the
  context where that payload is consumed.
- [x] No destination marker convention or new checker surface is added for a two-block population.
- [x] The prior source-file citation contract and the full canonical repository gate remain green.

## Decisions

- **D-01 (DECIDE & STATE):** Kept the fix to the two measured destination-bound payloads. The
  generated REQ and user-facing handoff use installed-root `.claude/skills/...` citations. This
  remains valid when a generated REQ moves from queue to working or archive, and satisfies the
  existing staged-runtime reference gate. No destination marker or general checker was added.
- **D-02 (DECIDE & STATE):** Replaced the review-generated producer test's source-package branch
  with one source-independent destination assertion. It accepts the two destination-safe forms
  already used by producers and rejects any source-relative `../do-work/actions/` payload citation.
- **D-03 (DECIDE & STATE):** The producer-only RED exposed the stale review assertion and the
  staged-runtime gate exposed that queue-relative `actions/...` is not stable after a REQ moves.
  Used the installed-root spelling instead of weakening or bypassing either existing gate.

## Implementation Summary

- `skills/do-work-toolbox/actions/code-review.md` (modified) — changed the generated REQ payload's
  two core-action citations to installed-root paths.
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified) — changed the user-facing
  capture handoff's two core-action citations to installed-root paths.
- `_dev/tests/contract-regressions.sh` (modified) — replaced the producer-source branch with a
  destination-safe Finding-Closure assertion that also rejects the old source-relative form.

The implementation adjusted an existing regression rather than adding a destination declaration or
checker. It stayed within the declared three-file write set.

## Testing

**Tests run:**
- Fence-scoped finding surface — RED before implementation with exactly four source-relative
  `(?:../)+do-work/actions/` occurrences; GREEN afterward with zero old occurrences and all four
  installed-root replacements present.
- `bash _dev/tests/contract-regressions.sh` — expected producer-only RED on the stale package-safe
  assertion, then exit 0 after the assertion was corrected.
- `bash _dev/tests/shipped-package-reference-contract.sh` — exit 0.
- `bash _dev/tests/staged-skills-contract.sh` — exit 0.
- `bash -n` and ShellCheck on the changed shell test — pass.
- `git diff --check` on the three implementation files — pass.
- `QUEUE_KANBAN_BROWSER=<headless-chromium> bash _dev/tests/maintainer-verify.sh` — exit 0,
  run directly and unpiped from the project root on the completed implementation tree.

**New tests added:** None; an existing review-generated-payload assertion was corrected and
narrowed.

**Existing tests updated (cross-REQ impact):** The REQ-170 Finding-Closure producer assertion now
checks the emitted payload's destination-safe form rather than the source package's relative path.

## Review

**Overall: 98%** | 2026-08-21T16:48:59Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — both affected payloads resolve in consumer topology, the four-token finding
surface is deleted, and the two-producer regression is non-vacuous and rejects the old form.
**Suggested testing:** 1 item
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Counting the affected payloads before designing selected the direct repair, while
running both citation gates exposed the existing assertion that encoded the wrong source-relative
form.

**What didn't:** The first replacement used `actions/...`, matching core-generated REQs but failing
the staged toolbox reference gate. A consumer-root installed path is stable across queue, working,
and archive moves and satisfies both readers.

**Worth knowing:** Template citations have two locations: the producer and the emitted artifact.
Source-package correctness says nothing about the destination; tests must judge the emitted form.

## Orientation

The two toolbox payloads that point into core actions now emit consumer-root installed citations,
and the existing review-generated producer contract rejects their retired source-relative form.
This is a leaf repair to the action/template reference contract; no new destination syntax was added.
