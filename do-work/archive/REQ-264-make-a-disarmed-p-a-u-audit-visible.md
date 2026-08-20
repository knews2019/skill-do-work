---
id: REQ-264
title: Make a disarmed P-A-U audit visible in qualify
status: completed
created_at: 2026-08-18T19:52:15Z
status_changed_at: 2026-08-18T20:55:14Z
claimed_at: 2026-08-20T12:16:31Z
completed_at: 2026-08-20T12:25:10Z
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-20T12:16:31Z
  basis:
    - trivial short-circuit
user_request: UR-055
addendum_to: REQ-254
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-300, REQ-263]
maintenance: true
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-cases/qualify.sh
- skills/do-work/actions/review-work.md
- skills/do-work-toolbox/actions/code-review.md
- _dev/tests/contract-regressions.sh
---

# Make a Disarmed P-A-U Audit Visible in Qualify

## What

Both of qualify's UNIFY-gated FAIL branches key on a checked `[UNIFY]` box in the REQ file — so a REQ with **no** P-A-U section at all sails through Check 4 with the FAIL half silently disarmed. Every review-generated REQ from the previous session (REQ-250 through REQ-254) lacks the section, and REQ-254's own qualification "Passed" that way: its review re-ran the range armed and got FAILs (fixture TODO lines, qualify's own regex, a doc seam line — false positives of the protected class) with no override on the record. qualify should WARN when the REQ file carries no P-A-U section, so a disarmed audit is visible instead of silent.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduce the silent pass first. Then WARN when the section is absent (counting boxes in any state, since the question is whether the audit exists, not its verdict) — and answer the REQ's build-time question by fixing the class at its source: the two templates that mint buildable REQs get the block, pinned by an assertion in the enumeration that already walks them.
- [x] **[APPLY]:** `qualify.sh` — `pau_box_total` plus the disarmed-audit WARN. `review-work.md` and `code-review.md` — the P-A-U block in each follow-up template. `contract-regressions.sh` — a fifth assertion in the existing `review_generated` template loop. `prescribed-shell-cases/qualify.sh` — the WARN case plus its vacuity guard.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `shellcheck --severity=warning` clean on all three shell files (only the two pre-existing structural `SC2154`s from the sourced harness); `maintainer-verify.sh` exits 0; suite 83 → 85 cases. Per-file: **`qualify.sh`** — the new grep counts `[ ]`, `[x]` and `[~]` so a partially-worked REQ is not misread as sectionless, and the WARN sits above the unchecked-box FAIL so a reader meets "disarmed" before "incomplete". **`review-work.md` / `code-review.md`** — block inserted directly under the `#` title in the fenced template only; no prose outside the fences touched, verified by reading the whole diff hunk for each. **`contract-regressions.sh`** — the assertion rides the existing loop and its count-parity vacuity guard, so a new producer package cannot escape it by not being listed. **`prescribed-shell-cases/qualify.sh`** — the vacuity guard reuses `$qualify_reporting_output` from the first case rather than building a second fixture, so it cannot pass by testing nothing.

## Context

REQ-254 review, Important finding 2 (gate: rule-change). Created `pending-answers` per the generation-≥2 depth stop. The companion record-keeping (the armed-run override for REQ-254's own range) was written into REQ-254's archived trail by the orchestrator at integration; this REQ is the rule so the next disarmed audit cannot be silent. Worth deciding at build time whether review-created follow-up REQ templates should simply include the P-A-U block (the capture template already does), which removes the class at the source.

## Open Questions

- [ ] REQ-254's review found qualify's box audit silently disarmed for REQs without a P-A-U section. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — orchestrators should notice a missing section themselves.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

---

## Triage

**Route: B** - Medium

**Reasoning:** The rule is one line ("WARN when the REQ carries no P-A-U section"), but the REQ also asks for a build-time decision about whether the follow-up templates should carry the block — which needed finding those templates, checking whether anything already pins them, and counting how live the class is.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The silent pass, reproduced before any edit.** A REQ file containing only `## Implementation Summary`, with a real leftover `print(raw_text)` added to library code, qualified clean: `OK: mechanical qualification passed`, exit 0. Nothing mentioned the missing section, and the debug print — genuine instrumentation in a file that never ends its own process — went unreported.

**Why absence satisfies every branch.** Check 4 has one FAIL keyed on an **unchecked** `[UNIFY]` box and (as of REQ-263) four keyed on a **checked** one. A REQ with no box at all satisfies all five by absence. The audit is disarmed, not passed, and the distinction was invisible.

**The class is live, and it has a source.** Six of the 26 queued REQs carry no P-A-U section — REQ-293, REQ-294, REQ-295, REQ-297, REQ-298, REQ-299 — every one `review_generated`. `grep -rn '^review_generated: true$' skills/*/actions/*.md` finds exactly two producer templates, and **neither emitted the block**: `skills/do-work/actions/review-work.md:368` (Step 10's Review Fix template) and `skills/do-work-toolbox/actions/code-review.md:302`. Both mint `status: pending` REQs, so both mint work that gets built with a disarmed audit.

**Something already walks those templates.** `_dev/tests/contract-regressions.sh` enumerates every fenced `review_generated: true` producer block across all shipped actions, asserts four properties per template, and guards against vacuity by requiring the template count to equal the count of exact `review_generated: true` fields. That is the right home for a fifth assertion: it inherits the vacuity guard, and a future producer package cannot escape it by not being on a list — the trap `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale names.

**Two nearby things that are correctly *not* in scope.** `work-reference.md`'s Builder-Decided Follow-up Template mints `status: pending-answers` "Confirm: X" REQs that clarify archives without ever building — a P-A-U block there would be three checkboxes nothing ever ticks. And `review-work.md:537`'s review checklist reads "**if** the REQ has an 'AI Execution State (P-A-U Loop)' section", tolerating the absence on the review side too; reported as a finding rather than fixed, since qualify's WARN now surfaces the same condition earlier.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/checks/qualify.sh` (modify) — the disarmed-audit WARN
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modify) — the WARN case and its vacuity guard
- `skills/do-work/actions/review-work.md` (modify) — P-A-U block in the Review Fix template
- `skills/do-work-toolbox/actions/code-review.md` (modify) — P-A-U block in the Code Review Fix template
- `_dev/tests/contract-regressions.sh` (modify) — fifth assertion in the existing producer-template loop

**Files I will NOT touch:**
- `skills/do-work/actions/work-reference.md` — its Builder-Decided Follow-up Template mints non-building `pending-answers` REQs; a P-A-U block there would be ceremony
- `skills/do-work/actions/work.md` — Step 6 already states P-A-U phasing is mandatory; this REQ enforces that statement rather than changing it
- the six queued REQs that lack the section — see the review's Minor finding for why they are left alone

**Acceptance criteria (restated from REQ):**
- [ ] qualify WARNs when the REQ file carries no P-A-U section, so a disarmed audit is visible instead of silent
- [ ] The WARN does not fire on a REQ that carries the section
- [ ] The build-time question about follow-up templates is decided and recorded
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ✓ passing (`bash _dev/tests/maintainer-verify.sh`, exit 0)
**Dependencies:** ✓ go1.26.1, ShellCheck 0.11.0, `just` present

*Checked by work action*

## Decisions

- **D-01 — WARN, not FAIL.** DECIDE & STATE. The REQ asks for a WARN and that is right: the missing section is the orchestrator's paperwork, not evidence about the code, and REQs written before the section existed are still legitimately qualifiable. A FAIL would block work on a defect in the REQ rather than in the change. The WARN text says "DISARMED, not passed" and spells out the consequence — every `[UNIFY]`-gated FAIL is unreachable — because "no P-A-U section" alone does not tell a reader that the run's OK means less than it looks.

- **D-02 — Yes: both follow-up templates get the block, and the assertion pins them.** DECIDE & STATE, answering the REQ's explicit build-time question. Six live REQs and two producer templates make this a source that is actively reproducing the class, not a hypothetical. Fixing the rule without fixing the source would leave qualify warning forever about REQs the skill itself keeps minting — `_dev/primes/prime-shell-commands.md` § Lessons, REQ-255: fix the class at the primitive. The boundary is *buildable*: a template that mints `status: pending` work gets the block; the `pending-answers` confirm template does not, because it has no implementation phase to phase.

- **D-03 — Scope extended mid-build, from two files to five.** The REQ's `write_set` named only `qualify.sh` and its case file, and D-02 needs three more. Recorded here and reflected in `## Scope` and `write_set` per work.md Step 6's rule that a file outside the declared boundary is a stop-and-report the orchestrator resolves in the trail, never a silent write. Reasoning: the extension is the REQ's own stated question answered yes; declining it to stay inside a capture-time field would have been the field governing the work instead of the reverse.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/checks/qualify.sh` (modified)
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work-toolbox/actions/code-review.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** qualify counts P-A-U boxes in any state and WARNs when there are none, naming the audit as disarmed rather than passed and stating that every `[UNIFY]`-gated FAIL is therefore unreachable. Both shipped templates that mint buildable review-generated REQs now emit the P-A-U block under their title, and the existing `review_generated` producer-template enumeration in `contract-regressions.sh` gained a fifth per-template assertion requiring it — inheriting that loop's count-parity vacuity guard. Two cases added: the WARN itself, and a vacuity guard proving it stays quiet on a REQ that has the section.

## Qualification

**Passed** — `tools/checks/qualify.sh` exit 0. Five files verified in the diff, all four acceptance criteria traced, P-A-U confirmed. One WARN, correct: the case file gained `print(`/`console.log` strings and owns its process exit, so they read as its own reporting.

`tools/checks/scope-drift.sh`: `OK: Implementation Summary matches the Scope declaration` (exit 0) — the mid-build scope extension in D-03 is reflected in both `## Scope` and `write_set`, so the widened boundary is declared rather than drifted into.

**Judgment checks:** *(2) Substantive* — a new counted grep plus a WARN branch in `qualify.sh`, a real block in each of two templates, a fifth assertion wired into an existing enumeration, two new cases. *(3) Requirements traced* — each criterion maps to a named case or run in `## Testing`. *(6) Flowing* — the WARN reads the actual REQ file and the assertion reads actual fenced template text; both are mutation-proven, which a stubbed check could not be.

Worth noting for the next builder of this file: the fixture-payload FAIL that REQ-263 hit did **not** recur here, because this REQ's two new cases assert on `grep -q 'DISARMED'` and reuse an existing fixture rather than writing new `TODO`-bearing payloads. The class from REQ-263's hand-back is real but it is triggered by adding marker payloads, not by adding cases as such.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-cases/qualify.sh`, `bash _dev/tests/contract-regressions.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ `qualify: 12 cases, 0 failures`; ✓ contract regressions exit 0; ✓ `maintainer-verify.sh` exit 0; suite total 83 → **85 named script cases across 17 per-script files**

**Red-green validation:**

- `REQ-960` no-P-A-U section: ✗ before, `OK: mechanical qualification passed`, exit 0, silent about both the missing section and a real leftover `print(raw_text)` in library code → ✓ after, the run names the audit `DISARMED` and states that every `[UNIFY]`-gated FAIL is unreachable
- vacuity guard: ✓ the same WARN does not appear on the first case's REQ, which carries the section — so the check keys on absence, not on every run
- template assertion: ✗ before, both producer templates lacked the block and nothing noticed → ✓ after, `contract-regressions.sh` asserts it per template inside the loop that already enumerates them

**Mutation testing — both new guards proven able to fail:**

| Mutation | Result |
|---|---|
| the disarmed-audit WARN removed from `qualify.sh` | `REQ-960` fails: "a REQ with no P-A-U section still passes Check 4 silently" |
| the P-A-U block dropped from `review-work.md`'s template | `contract-regressions.sh` FAILs, naming `review-work.md:359` and the consequence |

**New tests added:** two cases in `_dev/tests/prescribed-shell-cases/qualify.sh` (`REQ-960` plus its vacuity guard), and one per-template assertion in `_dev/tests/contract-regressions.sh`'s `review_generated` producer loop.

**Existing tests updated:** none. The ten cases from REQ-254 and REQ-263 pass unchanged.

*Verified by work action*

## Review

**Overall: 92%** | 2026-08-20T12:23:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 85% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 3 (report only)
- `skills/do-work/actions/review-work.md:537`'s review checklist still reads "**if** the REQ has an 'AI Execution State (P-A-U Loop)' section", which tolerates on the review side the same absence this REQ makes visible on the qualify side. Left unfixed deliberately: qualify now WARNs earlier in the pipeline, so the condition is surfaced before review reaches it, and tightening a review-contract line would start flagging the six existing sectionless REQs as findings — a separate decision with its own consequences.
- **Six queued REQs still carry no P-A-U section** (REQ-293, REQ-294, REQ-295, REQ-297, REQ-298, REQ-299) and were not backfilled. Each will now WARN when claimed, which is the intended behavior rather than a gap — its builder adds the section as part of its own `[PLAN]` phase, which is where the section belongs. Backfilling six queue files inside this commit would have been undeclared scope creep on top of an already-extended scope.
- The `## Scope` "Files I will NOT touch" list names `work-reference.md`'s confirm template as out of scope on the grounds that it never builds. That reasoning is sound but lives only here; nothing pins the *boundary* (buildable templates carry the block, confirm templates do not), so a future reader could add the block there for consistency and nothing would object.

**Restatement sweep:** the diff changes what a gate reports and what two templates emit, so the sweep asked who else states either. Ten sites mention P-A-U; the substantive ones are `work.md:391` ("P-A-U phasing is mandatory" — unchanged and now actually enforced, so it got *less* stale), `work.md:446` (the "in the diff" gloss already reported as a Minor under REQ-263), `review-work.md:527` (a Red Flag about checked boxes over a dirty diff — still accurate), and `review-work.md:537` (the Minor above). `capture-reference.md:54` and `sample-archived-req.md:30` carry the canonical block, and the two templates this REQ edited now match them verbatim.

**Acceptance:** Pass — all four criteria met with named evidence, and the REQ's build-time question is answered in D-02 rather than deferred.

**Suggested testing:** 1 item
- Nothing asserts that the block the two templates emit is *identical* to the canonical one in `capture-reference.md`. The new assertion checks the heading and all three box labels, which is the load-bearing part, but the parenthetical guidance text could drift between the three copies unnoticed.

**Scope 85%:** the declared `write_set` grew from two files to five mid-build. Recorded as D-03 with reasoning and reflected in both `## Scope` and `write_set` before the review ran, so it is a declared extension rather than drift — but it is still a 2.5× widening of a capture-time boundary and the score should say so.

**Follow-ups created:** None; **sweeps appended to:** None

## Lessons Learned

**What worked:** Asking who already walks the thing I wanted to pin, before writing a new check. `contract-regressions.sh` already enumerated every fenced `review_generated: true` producer template and already had a count-parity vacuity guard, so the fifth assertion cost six lines and inherited the guarantee that a future producer package cannot escape it by not being on a list. A standalone test for the same property would have been longer and weaker.

**What didn't:** Nothing failed outright, but the first framing was too narrow. Read literally, the REQ is one WARN in one script; its closing sentence asks a build-time question whose honest answer changes three more files. Treating the `write_set` as the boundary would have shipped a gate that warns forever about REQs the skill itself keeps minting — the rule without the source.

**Worth knowing:** "Disarmed" and "passed" print identically unless something says otherwise, and this was the second instance of that shape in two REQs: qualify's Check 4 here, and (REQ-263) its artifact scans reading no untracked files. Both were guards whose *input* was empty rather than clean. When a check is keyed on the state of a marker, the case worth writing first is the one where the marker is absent entirely.

## Orientation

A qualification run can no longer look clean because its audit was missing. When a REQ carries no P-A-U section, qualify says the box audit is disarmed and that every `[UNIFY]`-gated failure is unreachable, so an OK on such a REQ reads as what it is. The two shipped templates that mint buildable review-generated REQs now emit the block themselves, so the skill stops producing the REQs that trip it, and the enumeration that already walks those templates enforces it. Lives in the Step 6.3 qualification gate and the review/code-review follow-up producers (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — no new checklist item, field, or caller; one existing gate stopped being silently skippable and two templates gained a section the canonical capture template already had. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves, and its § Closed Enumerations Go Stale is the reason the assertion rides an enumerating loop rather than a filename list. The prime is not stale.
