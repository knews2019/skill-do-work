---
id: REQ-075
title: Five files still explain write_set's display-only status with a reason fan-out made false
status: completed
claimed_at: 2026-08-03T15:58:05Z
completed_at: 2026-08-03T16:07:20Z
commit: 738e9fe
kb_status: promoted
kb_entry: REQ-075-five-files-still-explain-write-set-s-di.md
route: B
created_at: 2026-08-03T15:20:00Z
user_request: UR-013
addendum_to: REQ-073
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
maintenance: true
review_generated: true
discovered_during: REQ-073
write_set: [actions/board.md, actions/capture-reference.md, actions/work-reference.md, docs/board-guide.md, docs/work-guide.md, tools/queue-kanban/prime-do-kanban.md, tools/queue-kanban/model.go, tools/queue-kanban/web/board.js, _dev/tests/contract-regressions.sh]
---

# Five Files Still Explain `write_set`'s Display-Only Status With a Reason Fan-Out Made False

## What

The conclusion is still right — nothing schedules, gates, or dispatches on `write_set` — but five files
give a reason for it that REQ-073 falsified. They say some version of:

> Under the exclusive-session model one REQ runs at a time, so the badge schedules nothing.

Since REQ-073, several builders **can** run at once under a single queue owner. So a reader following
that reasoning concludes the opposite of the contract: if the premise ("one REQ at a time") no longer
holds, the conclusion ("nothing schedules on it") looks like it should no longer hold either — and
`write_set` becomes a scheduling input. That is explicitly forbidden.

Replace the reason, keep the conclusion. The correct reason after REQ-073: `write_set` is **advisory
input to the human's pick, never a gate**, and the **merge** — `git merge --no-ff --no-commit`
refusing — is the non-interference proof. That holds at any builder count.

## Why This Is Worth a REQ Rather Than a Note

The stale text does not merely read oddly — it argues for the wrong behavior, in the files an agent is
most likely to read while touching the board or the scheduler. `actions/board.md` is loaded whenever
the board runs; `tools/queue-kanban/prime-do-kanban.md` is the prime an agent reads before editing the
tool. Both currently hand that agent a premise that has become false and a conclusion that depends on
it.

## Context

Each site, with the phrasing to replace:

- `actions/board.md:92` — "Under the exclusive-session model one REQ runs at a time, so the badge schedules nothing"
- `actions/board.md:117` — "display-only under the exclusive-session model, since one REQ runs at a time"
- `docs/board-guide.md:39` — "Under the exclusive-session model `do-work run` builds one REQ at a time, so the badge schedules nothing"
- `tools/queue-kanban/prime-do-kanban.md:57` — "Nothing schedules on `write_set` under the exclusive-session model (the work pipeline runs one REQ at a time)"
- `actions/capture-reference.md:113` — "nothing schedules on it under the exclusive-session model" (weakest of the five — it names the model without asserting one-REQ-at-a-time, so it may only need the pointer updated)

The corrected wording already exists in three places to copy from, all landed by REQ-073:
`actions/work.md` § Rules, `actions/work.md` Step 5.5, and `CLAUDE.md` § Shipped Tooling. The canonical
statement is `actions/work-reference.md` → Worktree Dispatch Mode → **Fan-Out Dispatch**.

Nothing about the board's *behavior* changes: `annotateWriteSetOverlap` still runs after bucketing,
still feeds only the badge and drawer row, and `tools/queue-kanban/` column logic stays untouched. This
is a prose-only correction.

## Detailed Requirements

1. **Replace the reason at all five sites**, keeping each file's existing voice and length. The
   conclusion ("nothing schedules / gates / dispatches on it") does not change.
2. **Do not restate the Fan-Out Dispatch contract** — point at it. `actions/work-reference.md` is its
   canonical home, and a sixth copy of the reasoning is what created this REQ.
3. **Keep "absence reads as unknown, not safe"** wherever it already appears. That property is
   independent of builder count and is load-bearing for anyone using the badge to pick a fan-out set.
4. **Add a contract-suite assertion** that no shipped file justifies `write_set`'s display-only status
   with a one-REQ-at-a-time premise — e.g. no file under `actions/`, `docs/`, or `tools/queue-kanban/`
   matches both a `write_set`/`overlaps` mention and a "one REQ at a time"-shaped clause in the same
   line. Without it, the sixth copy arrives with the next edit.
5. **Leave `tools/queue-kanban/*.go` alone.** No behavior change, no schema change.

## Discovered During

REQ-073 (`do-work/archive/` — fan-out dispatch), by that REQ's restatement sweep. The three sites
inside REQ-073's declared Scope were fixed there; these five sat outside it, and the sweep rule routes
an out-of-scope stale restatement to a follow-up rather than widening the original REQ's diff.

## Red-Green Proof

**RED prompt/case:** The assertion from requirement 4, run against the current tree — it must fail,
naming the five files.

**Why RED now:** All five still carry the stale premise; only the three REQ-073 touched were corrected.

**GREEN when:** The assertion passes, `grep -rn "one REQ at a time" actions/ docs/ tools/queue-kanban/`
returns nothing that is offered as a reason for `write_set`'s status, and the full contract suite stays
green (including REQ-073's own exactly-once invariant check).

## Full Context

See `do-work/user-requests/UR-013/input.md` for the batch input, and REQ-073's `## Review` for the
sweep that surfaced this.

---

## Triage

**Route: B** - Medium

**Reasoning:** The five sites come with line numbers, but the replacement reason has to be sourced from
the three landed sites rather than re-derived (requirement 2 forbids a sixth copy of the contract), and
requirement 4's assertion has to be designed against the real tree. Locations known, wording and
mechanism to discover.

**Planning:** Not required

## Plan

**Skipped** — Route B. Exploration below carries the discovery; no plan agent needed for a prose
correction plus one grep assertion.

## Exploration

**The REQ's site list is incomplete — the sweep finds eleven, not five.** Grepping the premise shape
(`one REQ at a time` / `one REQ runs at a time` / `runs one REQ at a time` / `builds one REQ at a time`)
across `actions/`, `docs/`, `tools/`, `crew-members/`, `_dev/`, `SKILL.md`, `CLAUDE.md`:

Listed by the REQ (all confirmed present, all stale):
1. `actions/board.md:92` — "so the badge schedules nothing"
2. `actions/board.md:117` — "display-only under the exclusive-session model, since one REQ runs at a time"
3. `docs/board-guide.md:39` — "builds one REQ at a time, so the badge schedules nothing"
4. `tools/queue-kanban/prime-do-kanban.md:57` — "Nothing schedules on `write_set` under the exclusive-session model (the work pipeline runs one REQ at a time)"
5. `actions/capture-reference.md:113` — "nothing schedules on it under the exclusive-session model"

**Not listed, same defect:**
6. `actions/work-reference.md:108` — the `write_set` **schema line itself**: "**Display only** — under the
   exclusive-session model one REQ runs at a time (`actions/work.md` Step 1), so nothing schedules on
   this field." This is the field's canonical home; leaving it stale while fixing five downstream
   restatements inverts the point of the REQ.
7. `actions/work-reference.md:441` — Scope Declaration Template: "(`write_set` is display, not
   scheduling — one REQ runs at a time under the exclusive-session model)".
8. `tools/queue-kanban/model.go:103` — `WriteSet` struct-field comment, same premise.
9. `tools/queue-kanban/model.go:1147` — `annotateWriteSetOverlap` doc comment, same premise.
10. `tools/queue-kanban/web/board.js:528` — the `overlaps` badge **tooltip a user actually reads**:
    "nothing schedules on write_set (do-work runs one REQ at a time)".
11. `_dev/tests/contract-regressions.sh:121` — the comment above the parser lock-step assertion states
    the premise as settled fact ("Under the exclusive-session model … one REQ runs at a time").

**Correctly-worded lines that must NOT be touched** (they say something else, or were already fixed by
REQ-073): `docs/work-guide.md:91` (fan-out overview — already correct), `actions/work-reference.md:302`
("Integration is serial" — a true statement about integration, not about `write_set`),
`actions/help.md:23` and `actions/capture.md:284` (describe queue processing, no `write_set` claim).

**Replacement wording to copy from**, all landed by REQ-073: `actions/work.md` § Rules (the fullest
form), `actions/work.md` Step 5.5, `CLAUDE.md` § Shipped Tooling. Canonical home to point at:
`actions/work-reference.md` → Worktree Dispatch Mode → **Fan-Out Dispatch**.

**Requirement 4 vs. requirement 5 conflict, and requirement 4 vs. a correct file.** See D-01 and D-02.

## Scope

**Files I will touch:**
- `actions/board.md` — replace the reason at both sites
- `actions/capture-reference.md` — replace the reason
- `docs/board-guide.md` — replace the reason
- `tools/queue-kanban/prime-do-kanban.md` — replace the reason
- `actions/work-reference.md` — the schema line and the Scope template (**added mid-build, D-01**)
- `tools/queue-kanban/model.go` — two comments only, no code (**added mid-build, D-01**)
- `tools/queue-kanban/web/board.js` — the badge tooltip string (**added mid-build, D-01**)
- `docs/work-guide.md` — split the correct-but-overlong bullet so the assertion needs no exclusion (**added mid-build, D-02**)
- `_dev/tests/contract-regressions.sh` — the new assertion (requirement 4) and its own stale comment

**Acceptance criteria (restated from the REQ):**
- The reason is replaced at every site; the conclusion ("nothing schedules / gates / dispatches on it") is unchanged everywhere.
- No site restates the Fan-Out Dispatch contract — each points at its canonical home.
- "Absence reads as unknown, not safe" survives wherever it already appeared.
- A contract assertion fails on a shipped file that justifies `write_set`'s display-only status with a one-REQ-at-a-time premise.
- `tools/queue-kanban/*.go` behavior and schema unchanged.

## Decisions

**D-01 — Extended scope to the six unlisted stale sites, including two Go comments and the board
tooltip. ESCALATE (surfaced here, decided by the evidence).**
Requirement 5 says "Leave `tools/queue-kanban/*.go` alone. No behavior change, no schema change." But
requirement 4 asks for an assertion covering `tools/queue-kanban/`, and `model.go` carries the stale
premise in two comments — so as literally specified, requirement 4's assertion cannot pass while
requirement 5 is honored. Resolution: read requirement 5 by its stated reason ("no behavior change, no
schema change") and change **comments only**. Not one line of Go logic, no struct field, no test
fixture. The same argument covers `board.js:528`, which is a user-visible tooltip rather than code, and
`actions/work-reference.md:108`, which is the field's canonical home — fixing five restatements while
the definition they restate stays stale would leave the REQ's own conclusion unsupported.
**Value:** the sweep actually converges; an agent editing the board or the schema reads the true reason.
**Risk:** low, reversible — every hunk is prose. The declared file list grew from five files to nine,
which is real scope growth a reviewer should weigh; the alternative is a REQ that fixes 5 of 11 sites
and a green assertion that certifies the other 6 do not exist.

**D-03 — The assertion is a line sweep plus two file-level negatives, not a proximity window.
DECIDE & STATE.**
Requirement 4 suggested "both in the same line." Built that way first, it passed over `model.go` and
`board.js` entirely — their comments wrap, so the premise and the field land on different lines, and the
check was blind inside the very tool it guards. Widening to a 3-line window fixed that and immediately
false-positived on `actions/work-reference.md:302`, two lines under the advisory-`write_set` bullet in
the canonical Fan-Out section, where "integration … runs one REQ at a time" is a true statement about
integration. So: line granularity for the prose class (long lines, no false positives), plus a
file-level "must not contain a builder-count claim" negative for `model.go` and `board.js`, which have
no legitimate reason to state one. Both halves are in the suite with the reasoning as comments.

**D-04 — Also fixed the weaker "under the exclusive-session model" phrasing inside the same files.
DECIDE & STATE.**
After the eight one-REQ-at-a-time sites were clean, `model.go` (2×) and `board.js` (1×) still justified
the field with "under the exclusive-session model nothing schedules on write_set" — the REQ's own
site-5 shape. The exclusive-session model is *not* falsified (one queue owner still holds), but it
stopped being the reason: fan-out runs several builders under exactly that model. No new files — all
three are inside D-01's declared set.

**D-02 — Split `docs/work-guide.md:91` rather than exclude it from the assertion. DECIDE & STATE.**
That bullet is correct (REQ-073 updated it) but is five sentences on one physical line, so it carries
both "by default one REQ at a time" and a `write_set` mention — a line-granularity assertion
false-positives on it. The options were an exclusion list in the check or splitting the bullet.
Excluded-file lists are the staleness pattern this repo already ratchets against, and the split is a
readability win on its own (the `write_set` hint becomes findable as its own bullet). No wording
changes — the sentences move, they do not change.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — the `write_set` schema line and the Scope-template mirror
  note now give the durable reason and point at Fan-Out Dispatch
- `actions/board.md` (modified) — both sites (badge explanation, parser lock-step note)
- `actions/capture-reference.md` (modified) — Populating `write_set`, the weaker pointer-only shape
- `docs/board-guide.md` (modified) — "Reading the `overlaps` badge", rewritten for a user audience
- `docs/work-guide.md` (modified) — the `write_set` sentence split into its own bullet (D-02); wording
  unchanged
- `tools/queue-kanban/prime-do-kanban.md` (modified) — the REQ-032/034 lesson entry
- `tools/queue-kanban/model.go` (modified) — **four comments only**, zero code: the `WriteSet` field
  comment, `annotateWriteSetOverlap`'s doc, `writeSetPatternsIntersect`'s doc, `writeSetsIntersect`'s doc
- `tools/queue-kanban/web/board.js` (modified) — the badge tooltip string a user reads, plus the comment
  above it
- `_dev/tests/contract-regressions.sh` (modified) — the new two-part assertion (requirement 4) and the
  stale premise in the existing parser lock-step comment

**What was done:** Replaced the reason at every site that argued `write_set`'s display-only status from
a builder count, keeping every conclusion and every "absence reads as unknown" clause intact. Each site
now points at `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch rather than
restating it (requirement 2). Eleven sites carried the strong premise or its weaker
exclusive-session-model variant — six more than the REQ listed; the extras are recorded in D-01 and D-04.
No Go logic, no schema field, no test fixture changed: `go vet` and `go test ./...` pass unchanged.

## Qualification

Passed — 9 files verified, 5 acceptance criteria traced, scope-drift comparison clean (`OK:
Implementation Summary matches the Scope declaration` after the D-01/D-02 scope extension was written
into `## Scope` and `write_set` *before* the extra files were touched, not after). Mechanical checks
green via `tools/checks/qualify.sh`.

## Testing

**Tests run:**
- `bash _dev/tests/contract-regressions.sh` → ✓ passed
- `cd tools/queue-kanban && go vet ./... && go test ./...` → ✓ `ok` (no behavior change expected or seen)

**Red-green validation:** *(REQ carries `tdd: true`; the assertion was written before any fix)*
- `contract-regressions.sh` REQ-075 line sweep: ✗ FAILed listing 8 offending lines across
  `actions/work-reference.md`, `actions/board.md`, `docs/board-guide.md`, `docs/work-guide.md`,
  `tools/queue-kanban/prime-do-kanban.md`, `tools/queue-kanban/web/board.js` → ✓ passes after the fixes
- `contract-regressions.sh` model.go file-level negative: ✗ FAIL → ✓ pass
- `contract-regressions.sh` board.js file-level negative: ✗ FAIL → ✓ pass
- All three re-verified by stashing only the source fixes and keeping the assertions, then restoring —
  not inferred from the pre-fix run

**GREEN criteria from `## Red-Green Proof`, checked individually:**
- The assertion passes ✓
- `grep -rIEn 'one REQ( [a-z]+)? at a time' actions/ docs/ tools/queue-kanban/ | grep -E 'write_set|overlaps'`
  → no matches ✓
- Second shape swept too: `grep -rn "exclusive-session" … | grep -iE "write_set|overlap"` → no matches ✓
- Full contract suite green, REQ-073's exactly-once invariant included ✓

**Surviving `one REQ` mentions, checked individually and deliberately left:**
`actions/work-reference.md:302` (integration is serial — true, and about integration),
`docs/work-guide.md:91` and `:120` (queue processing, no `write_set` claim), `actions/help.md:23`
(command summary), `actions/capture.md:284` (a splitting rationalization row).

## Review

**Overall: 91%** | 2026-08-03T16:06:37Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 80% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor

- **Minor — the declared scope nearly doubled mid-build** (5 files → 9). Handled correctly by the
  contract (`## Scope` and `write_set` extended before the extra files were touched, reasoning in D-01),
  and the added sites were the same defect, not new work. Still the honest read: the REQ's own site
  inventory was 6 short, which is worth noting because that inventory came from a restatement sweep —
  the same instrument this REQ exists to serve. A sweep that reports "five files" should be read as a
  floor.
- **Minor — requirement 5 was honored by intent, not by letter.** "Leave `tools/queue-kanban/*.go`
  alone" was read as "no behavior or schema change" and four comments were edited. A reviewer who reads
  requirement 5 literally will see a violation; D-01 states the conflict and why the literal reading
  makes requirement 4 unsatisfiable.

**Restatement sweep:** This REQ *is* a restatement sweep, so the check turned on its own output. Two
shapes exist for this premise, not one — the strong form ("one REQ at a time") and the weak form
("under the exclusive-session model") — and finding the second only after the first was clean (D-04) is
the reusable lesson.

**Acceptance:** Pass — every GREEN criterion checked individually rather than as a bundle; the Go build
and test suite confirm the comment-only claim.
**Suggested testing:** 1 item — render the board (`do-work board`) and read the `overlaps` badge tooltip
in a browser, the one changed string with a human audience. Text-only change, no layout risk, but it is
the only site the suite cannot judge for readability.
**Follow-up REQs created:** None — both findings are Minor and self-documenting in the trail.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Writing the assertion first and letting it enumerate the sites, instead of trusting the
REQ's list. The list said five; the check found eight, and a second grep shape found three more. On a
sweep REQ the check is the inventory — the prose inventory is a starting hypothesis.

**What didn't:** Two dead ends, both instructive. (1) Line-granularity matching — the REQ's own
suggestion — is blind to wrapped Go/JS comments, so the first green run was green for the wrong reason.
(2) Widening to a 3-line proximity window then false-positived on the canonical Fan-Out section itself,
where "integration runs one REQ at a time" is true. The shape that works is per-class: a line sweep for
prose, a file-level negative for files that have no business naming a builder count at all.

**Worth knowing:** A falsified premise leaves two fingerprints, and grepping for one of them reads as
clean. Here the strong form named the count ("one REQ at a time") and the weak form named the model
("under the exclusive-session model"). The weak form is the more dangerous of the two, because the model
it names is *still true* — only its relevance died — so the sentence reads correct on inspection. When a
premise is retired, sweep for the thing it was *called* as well as the thing it *said*.

## Orientation

The board's file-overlap badge now explains itself in terms that survive fan-out: it is advisory input to
a human's pick, and the merge is what proves two builders didn't collide. Same behavior, same badge,
same column logic — what changed is that an agent reading the board action, the schema, the Go source,
or the tooltip is no longer told the badge is harmless *because only one REQ runs at a time*, which
stopped being true at 0.166.0. Lives in the queue-kanban board subsystem plus the pipeline's `write_set`
schema. Not `[MAP CHANGED]` — no new module, contract, or data flow.

**Prime staleness spot-check** (`tools/queue-kanban/prime-do-kanban.md`, this REQ's only prime): its
referenced paths (`tools/queue-kanban/model.go`, `web/board.js`, `_dev/tests/contract-regressions.sh`)
all still exist, and the REQ-032/034 lesson entry was updated in place. Not stale.
