---
id: REQ-079
title: Two guards pin the weaker fingerprint of the premise they exist to retire
status: completed
claimed_at: 2026-08-03T21:58:03Z
completed_at: 2026-08-03T22:04:14Z
commit: 8fdce3c
kb_status: pending
route: B
created_at: 2026-08-03T16:53:42Z
user_request: UR-015
domain: testing
prime_files: []
tdd: true
depends_on: []
maintenance: true
addendum_to: REQ-075
---

# Two guards pin the weaker fingerprint of the premise they exist to retire

## What

REQ-075 (v0.166.2) established that a retired premise leaves two fingerprints — the thing it *said*
("one REQ at a time") and the thing it was *called* ("under the exclusive-session model") — and that
the second is the more dangerous of the two. Its own regression check pins only the first. Its regex is
also narrower than the class it names. And `actions/cleanup.md:31` still argues the safety of the
skill's one destructive pass from a premise REQ-071 spent an entire REQ falsifying.

No live defect today: the shipped tree is clean of both fingerprints. This is a durability gap — the
guards will not catch the recurrence they were written for.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty. Loaded `crew-members/general.md`, `coding-guardrails.md`,
  `maintenance.md` (`maintenance: true`), `testing.md` (`tdd: true`, and this REQ is almost entirely
  test work). Approach: hoist the strong-form regex into one named variable so the three consumers
  stop carrying three copies (subtraction — the triplication in F4b is itself the bug), widen it to
  the class, add the weak-form twin as a line sweep plus two file-level negatives, then rule on the
  three remaining exclusive-session sites individually.
- [x] **[APPLY]:** Two files, both declared. Requirement 5's rulings are recorded below; two of four
  sites were deliberately left alone.
- [x] **[UNIFY]:** `git diff --stat` → `_dev/tests/contract-regressions.sh` and `actions/cleanup.md`.
  Verified: `shellcheck` clean, `bash -n` parses, suite exits 0 on the clean tree. Both patterns are
  defined once and referenced by every consumer — no third copy. Read the `actions/cleanup.md` hunks
  against Pass 0's actual steps to confirm the new safety argument matches what the pass does (step 4
  is what makes "terminal only" true). No debug artifacts.

## Why

A guard that matches only the wording that was already removed is worse than no guard: the suite goes
green, and the green reads as coverage. REQ-075 wrote down the exact insight needed to avoid this —
in `tools/queue-kanban/prime-do-kanban.md:60`, where it calls the weak form "the more dangerous,
because the model it names is still true and only its relevance died" — and then did not encode it.
The lesson is filed; the assertion does not implement it.

## Context

**Finding F4a — no assertion covers the weak fingerprint.** `_dev/tests/contract-regressions.sh:137`
greps `one REQ( [a-z]+)? at a time` and filters to lines mentioning `write_set|overlaps`. Nothing
asserts against the weak form near `write_set`. Verified: `grep -n "exclusive-session" _dev/tests/contract-regressions.sh`
returns lines 184, 190, 191, 205, 213, 222 — all about REQ-069's removed machinery, none about
`write_set`.

**Finding F4b — the regex is narrower than the class.** The optional word slot sits *after* `REQ`, so
the pattern matches "one running REQ at a time" but not:

- `one active REQ at a time`
- `a single REQ at a time`
- `one builder at a time`
- `only one REQ is ever building`

Each of those reintroduces the premise and passes the suite. The same regex is reused verbatim in both
file-level negatives (lines ~153 and ~158), so the gap is triplicated.

**Finding F7 — `actions/cleanup.md:31` argues from the falsified premise.** It reads: "**Safe under the
exclusive-session model.** This pipeline assumes no other `do-work` session is running against this
checkout … so there is no live coexisting claim to protect — Pass 0 needs no lock and consults none."
REQ-071 exists *because* a live coexisting claim can be there. Pass 0's behavior is probably still
safe — it sweeps only terminal statuses, and integration is serial — but this is the skill's one
destructive pass reasoning from the premise the safety REQ removed, and its next sentence ("that
session's own REQ has already been moved out of `working/`") is singular where fan-out is plural.

## Detailed Requirements

1. **Add a weak-form assertion** for the `write_set` rule: no shipped file may justify `write_set`'s
   display-only status by naming the exclusive-session model. Keep REQ-075's per-class shape — a line
   sweep for prose, file-level negatives for the comment-carrying sources (`tools/queue-kanban/model.go`,
   `tools/queue-kanban/web/board.js`) whose comments wrap past a line-granularity grep.
2. **Widen the strong-form regex** to cover the class rather than one phrasing. Move the optional word
   slot so it can precede `REQ`, and cover the `single`/`builder` variants. Per the Closed Enumerations
   rule, the assertion's comment must state the trigger *condition* — "no shipped file may argue
   `write_set`'s display-only status from any builder count" — and mark the matched wordings as
   illustrative, not exhaustive.
3. **Do not let the widened pattern false-positive on the canonical section.** REQ-075's original
   comment records why granularity is one line: `actions/work-reference.md`'s Fan-Out Dispatch section
   says "integration … runs one REQ at a time" two lines below the advisory-`write_set` bullet, and
   that statement is **true**. Re-verify after widening that the suite still passes on a clean tree —
   a widened pattern is exactly the change that could start matching it.
4. **Correct `actions/cleanup.md:31`.** Replace the exclusive-session justification with the durable
   one: Pass 0 is safe because it sweeps only *terminal* statuses and integration is serial under one
   queue owner, not because no other claim can exist. Fix the singular in the second sentence — under
   fan-out several REQs sit in `working/` at once. Point at
   `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch rather than restating it.
5. **Sweep the remaining weak-form sites and rule on each.** `actions/work.md:549` and
   `actions/cleanup.md:40` also invoke the exclusive-session model. Some uses are legitimate — the
   model *is* still true, only its relevance to a given conclusion died. For each site, state whether
   the conclusion still follows from the premise, and correct only the ones where it does not. Do not
   sweep the phrase mechanically; that would delete true statements.
6. **Leave the badge's behavior alone.** No Go logic, no schema field, no board column. Same as
   REQ-075: the explanation changes, the code does not.

## Constraints

- **This REQ can only be proven by making the suite fail.** A guard change that is never observed
  failing is not evidence of anything. Every new or widened assertion must be demonstrated red by
  reintroducing the wording it targets, then green after reverting — record the observed failure text.
- Requirement 5 is where this REQ can do damage. `crew-members/maintenance.md`'s delete-before-you-add
  rule cuts both ways: removing a *true* premise because it pattern-matches a retired one is the
  failure mode here. Err toward leaving a site alone and saying why.
- `_dev/tests/contract-regressions.sh` is also touched by REQ-077 (the REQ-071 guard at ~line 364).
  Different block; if both are ever built concurrently the merge is the non-interference proof, not the
  overlaps badge.

## Dependencies

`addendum_to: REQ-075` — completes the guard REQ-075 specified in its lesson and under-implemented in
its assertion. Overlaps REQ-077 in one file (see Constraints). No `depends_on`: buildable immediately.

## Builder Guidance

**Certainty: Firm on requirements 1–4, genuinely open on 5.**

The regex and assertion work is mechanical and the findings were verified by grep; don't re-derive
them, but do re-check line numbers since REQ-077 may move them.

Requirement 5 needs judgment, and the right answer may well be "all three remaining sites are fine."
That is a valid outcome — record the per-site ruling rather than forcing a change to justify the REQ.
The value here is the two assertions, not a body count.

Scope note: this REQ is the guard, not another sweep. If requirement 5 turns up a *new* class of stale
reasoning beyond the exclusive-session premise, capture it as a discovered task rather than expanding
this REQ — REQ-075's scope doubling mid-build is the precedent to avoid repeating.

## Red-Green Proof

**RED case:** Add the line
`` `write_set` is display-only because only one active REQ builds at a time under the exclusive-session model. ``
to `actions/board.md` and run `bash _dev/tests/contract-regressions.sh`. It passes today — the strong
regex misses "one active REQ" (the optional slot is on the wrong side of `REQ`), and no assertion
covers the exclusive-session clause at all.

**Why RED now:** Both guards match one historical phrasing rather than the class, so the premise
REQ-075 retired can walk back in under any of four wordings, including the one REQ-075's own lesson
names as the more dangerous.

**GREEN when:** (1) The RED line above fails the suite and the failure message names the file and the
fix. (2) Each of the four wordings in Context fails the suite. (3) The suite still passes on an
otherwise clean tree, with Fan-Out Dispatch's true "integration runs one REQ at a time" sentence
untouched. (4) `actions/cleanup.md:31` no longer argues Pass 0's safety from the exclusive-session
premise. (5) The per-site rulings from requirement 5 are recorded, including the leave-alone ones.

**Validation:** Inferred during an adversarial audit; remediation plan reviewed and approved by the
user before capture.

## Full Context

See `do-work/user-requests/UR-015/input.md` for the audit's provenance and the findings it cleared.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (two assertions, one widened regex, one corrected
paragraph) and every site was named and grep-verified at capture. What needed discovery was only
*where* the false-positive risk actually sits after widening, and how the four remaining
exclusive-session sites each reason — exploration, not planning.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Line numbers re-checked after REQ-077 and REQ-078 moved things.** The REQ-075 guard block is now at
`_dev/tests/contract-regressions.sh:150-181` (was ~137), because REQ-078 inserted the timestamp
single-source block above it. `actions/cleanup.md:31` and `:40` are unmoved.

**The false-positive risk is real and the existing design already handles it.** `actions/work-reference.md:316`
— Fan-Out Dispatch's "**Integration is serial.** … merge → qualify → test → review → changelog →
archive runs one REQ at a time." — matches the widened pattern, as it matched the narrow one. What
saves it is the `write_set|overlaps` line filter, not the pattern: that sentence names neither. So
requirement 3 is satisfied by *keeping* REQ-075's line-granularity design, not by tuning the regex,
and the widened pattern was dry-run against the whole tree with the filter applied before anything
was committed to.

**A false-positive risk the REQ did not name:** `tools/queue-kanban/prime-do-kanban.md:60` is the
REQ-075 **lesson entry itself**, and it quotes both fingerprints verbatim ("one REQ at a time" and
"under the exclusive-session model") in one line. It is inside `tools/queue-kanban`, which is one of
the three swept directories. It survives for the same reason the Fan-Out line does — the line names
neither `write_set` nor `overlaps` — and it is the reason no file-level negative may ever be added
for that file. Recorded in the assertion's comment so a later "let's be thorough" edit doesn't add one.

**Requirement 5 inventory — four exclusive-session sites in shipped files, ruled individually:**

| Site | Conclusion drawn | Does it follow? | Ruling |
|---|---|---|---|
| `actions/cleanup.md:31` | Pass 0 is safe / "no live coexisting claim to protect" | **No.** `actions/work-reference.md` → Crash Recovery exists because a `working/` claim may not be this session's. | **Corrected** (requirement 4) |
| `actions/cleanup.md:40` | a terminal REQ in `working/` came from a crashed prior run | Yes, but the clause carries no inferential load — a terminal REQ is finished whoever finished it | **Narrowed** — clause deleted, sentence unchanged |
| `actions/work.md:550` | "Nothing else is released … the exclusive-session model keeps none" | **Yes.** The claim is that the model keeps no lock, which is exactly what the model says. | **Left alone** |
| `actions/work-reference.md:407` | "the skill keeps no lock, heartbeat, or liveness check" | **Yes.** Same shape — citing the model for its own content. | **Left alone** |

Two of four left alone, which is the outcome Builder Guidance said was valid. The distinguishing test:
the premise is being cited for *what it asserts* (no lock) versus for a *consequence it does not
support* (no coexisting claim, or any statement about builder counts).

**No new class of stale reasoning turned up**, so nothing was captured as a discovered task on that
account — the scope note's warning about REQ-075's mid-build doubling did not need to bite.

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modify) — hoist + widen the strong-form pattern, add the
  weak-form line sweep and two file-level negatives, restate the block's trigger condition
- `actions/cleanup.md` (modify) — Pass 0's safety argument (requirement 4) and the `:40` clause

**Files I will NOT touch:** `actions/work.md:550` and `actions/work-reference.md:407` (ruled
legitimate above), `tools/queue-kanban/prime-do-kanban.md` (the lesson entry must keep quoting both
fingerprints), and all Go/JS logic — requirement 6 is explicit that the explanation changes and the
code does not.

**Acceptance criteria (restated from REQ):**
- [ ] The REQ's stated RED line fails the suite, naming the file and the fix
- [ ] Each of the four Context wordings fails the suite
- [ ] The suite still passes on a clean tree, Fan-Out Dispatch's true sentence untouched
- [ ] `actions/cleanup.md:31` no longer argues Pass 0's safety from the exclusive-session premise
- [ ] Every requirement-5 site has a recorded ruling, including the leave-alone ones

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `actions/cleanup.md` (modified)

**What was done:** The strong-form pattern is now defined **once** as `builder_count_premise_pattern`
and referenced by all three of its consumers, instead of being spelled out three times — F4b's gap was
triplicated precisely because the pattern was. It was widened from `one REQ( [a-z]+)? at a time` to
cover the class: the count word may be `one`/`a single`/`only one`, up to two adjectives may precede
the noun, the noun may be `REQ`/`builder`/`coder`/`agent`, and the tail may be `at a time`/`at
once`/`concurrently` or the `is ever building|running|in flight` shape.

A second pattern, `exclusive_session_premise_pattern`, pins the **weak** fingerprint REQ-075's own
lesson named as the more dangerous one and then failed to encode. It gets the same two-shape treatment
as the strong form: a line sweep over `actions/`, `docs/` and `tools/queue-kanban` filtered to lines
mentioning `write_set`/`overlaps`, plus file-level negatives on `tools/queue-kanban/model.go` and
`web/board.js`, whose comments wrap past a line grep. Four assertions where there were two.

The block's comment now states the **trigger condition** — no shipped file may argue `write_set`'s
display-only status from any builder count, in either fingerprint — and marks the enumerated wordings
illustrative. It also records why `tools/queue-kanban/prime-do-kanban.md` must never get a file-level
negative: its REQ-075 lesson entry quotes both fingerprints on one line, legitimately.

`actions/cleanup.md`'s Pass 0 preamble no longer argues safety from the exclusive-session model. The
durable reason replaces it — Pass 0 sweeps only *terminal* statuses, and a terminal REQ is finished
whoever finished it, which holds at any builder count — plus an explicit "do not argue it from the
exclusive-session model" and a pointer to Crash Recovery for why a `working/` claim may not be this
session's. The singular "that session's own REQ" became plural with a Fan-Out Dispatch pointer, and
`:40`'s load-free `, under the exclusive-session model` clause was deleted.

## Qualification

**Passed** — 2 files verified; 6 requirements traced.

- r1 → the weak-form line sweep + two file-level negatives, in REQ-075's per-class shape.
- r2 → the widened pattern, hoisted to one definition, with the trigger condition in the comment and
  the wordings marked illustrative.
- r3 → dry-run against the whole tree before adoption; clean tree exits 0 and Fan-Out Dispatch's
  sentence is untouched. Demonstrated below, including as an explicit negative control.
- r4 → `actions/cleanup.md:31` rewritten.
- r5 → four sites ruled in the Exploration table; two corrected, two left alone with reasons.
- r6 → zero Go/JS logic changes; `git diff --stat` shows only the suite and one action file.
- **P-A-U audit:** three boxes with evidence; diff matches. **Contamination check:** REQ-078 also
  touched `_dev/tests/contract-regressions.sh` — expected and predicted by this REQ's own Constraints
  ("Different block"); verified by reading the hunks, they are ~140 lines apart and independent.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Passing on the clean tree (exit 0)

**Red-green validation** — the REQ's Constraints require every assertion be *observed* failing, so
each was, by appending the offending line to `actions/board.md`, running the suite, and reverting:

| Injected line | Result |
|---|---|
| `write_set` is display-only because only one active REQ builds at a time under the exclusive-session model. *(the REQ's stated RED case)* | **RED** — trips **both** guards, 2 lines naming `board.md` |
| `write_set` is display-only because only one REQ runs at a time. | **RED** |
| The `overlaps` badge is advisory since a single REQ at a time is ever claimed. | **RED** |
| The `overlaps` badge is safe because one builder at a time touches the tree. | **RED** |
| `write_set` never gates because only one REQ is ever building. | **RED** |
| `write_set` needs no lock under the exclusive-session model. *(weak form alone)* | **RED** |
| **Negative control:** Integration is serial: merge through archive runs one REQ at a time. | **stays GREEN** — correct: a true statement naming neither `write_set` nor `overlaps` |

- **The gap is real, and the old suite proves it.** Running `HEAD`'s copy of the suite against three of
  those injected lines — including the REQ's stated RED case and the weak form — exits **0** every
  time. That is F4a and F4b as a single observation: the guards were green on the recurrence they
  exist to catch.
- **File-level weak-form negative fires:** appending
  `// write_set is display-only under the exclusive-session model.` to `tools/queue-kanban/model.go`
  produces exactly the new `model.go must not invoke the exclusive-session model` failure. Reverted.
- **Clean tree, before and after every case:** exit 0, `git status` clean of the probe files.

**New tests added:** 2 assertions (weak-form line sweep; two file-level weak-form negatives — counted
as one rule in two files), plus the widening of the existing three.

**Existing tests updated (cross-REQ impact):** the REQ-075 strong-form pattern (all three consumers)
was widened and de-duplicated. Intentional and strictly broader — every string the old pattern caught
is still caught.

*Verified by work action*

## Decisions

- **D-01 — Hoisted the strong-form pattern into one variable before widening it.** DECIDE & STATE.
  F4b reports the gap as "triplicated," which is the tell: the pattern was spelled out three times, so
  a fix had to be applied three times and could be applied twice. Naming it once makes the triplication
  structurally impossible rather than caught-by-review. This is subtraction, and it is why the widening
  itself is a one-line change. Reversible.
- **D-02 — Kept the `write_set|overlaps` line filter rather than tightening the regex to spare the
  Fan-Out sentence.** DECIDE & STATE. Requirement 3 warns that a widened pattern is exactly what could
  start matching `actions/work-reference.md:316`, and it does match it — the pattern cannot tell a true
  statement about integration from a false one about `write_set`, and no regex can. Only the *filter*
  can, because the premise is wrong only where it explains the field. So the answer to "don't
  false-positive on the canonical section" is REQ-075's existing design, unchanged. Reversible.
- **D-03 — Left `actions/work.md:550` and `actions/work-reference.md:407` alone.** DECIDE & STATE, and
  the outcome Builder Guidance predicted. Both cite the exclusive-session model for *what it asserts*
  (the pipeline keeps no lock), which is still true and still the point. Sweeping them would delete
  true statements to make a phrase-count go to zero — the specific damage requirement 5 and
  `maintenance.md`'s "subtraction is not vandalism" both warn about. The per-site test that separated
  them from `cleanup.md:31` is recorded in the Exploration table so the next sweeper doesn't re-derive
  it.

## Lessons Learned

**What worked:**
- **Dry-running the widened pattern against the whole tree, with the filter applied, before adopting
  it.** Requirement 3 asked for a re-check after widening; doing it as the *first* step instead turned
  a possible late surprise into a design input — it is what showed that the filter, not the regex, is
  what protects the canonical sentence.
- **A ruling table with a stated test.** Four sites, two corrected, two left alone, and one sentence
  saying how they were told apart (cited for what the premise asserts vs. for a consequence it does not
  support). That is reusable; a list of four verdicts is not.

**What didn't:**
- **The obvious way to "be thorough" would have broken the lesson.** `tools/queue-kanban/prime-do-kanban.md:60`
  quotes both fingerprints verbatim and sits inside a swept directory. Adding a file-level negative
  there — the same treatment `model.go` and `board.js` get, and the natural next step if you are
  pattern-matching on "cover every file in `tools/`" — would make the suite fail on the very lesson
  entry that documents the rule. It survives only because the line filter finds no `write_set` on it.
  Written into the assertion's comment as a do-not.

**Worth knowing:**
- **A guard's blast radius is the pattern *and* the filter, and they do different jobs.** The pattern
  decides what the premise looks like; the filter decides where the premise is wrong. Widening the
  first without understanding the second is how you either miss the class or flag true statements —
  and this REQ needed both to stay separate to satisfy requirements 2 and 3 at once.
- **The weak fingerprint is not always wrong**, which is what makes it survive review: two of the four
  sites cite the exclusive-session model correctly. A mechanical sweep of the phrase would have been a
  regression. The line filter is again what makes the assertion safe — it only fires where the premise
  is being used to explain `write_set`.
- **"Suite passes" was the false signal this REQ existed to remove.** Running `HEAD`'s suite against
  the injected wordings and watching it exit 0 three times is worth more than any argument about
  regex coverage — and it took about a minute.

## Orientation

The two guards protecting `write_set`'s display-only rule now catch the class of premise they were
written for, instead of the one wording that had already been removed. Both of the retired premise's
fingerprints are pinned — the count form and the exclusive-session form — across prose by line and
across comment-carrying source by file. Lives in the contract-regression suite; no runtime behaviour
changed anywhere.

Separately, `do-work cleanup`'s Pass 0 now argues its own safety from something true at any builder
count (it sweeps only terminal statuses) rather than from the premise a prior REQ spent itself
falsifying.

No `[MAP CHANGED]`: no new module, data flow, contract, or renamed concept — the explanations changed
and the code did not, exactly as requirement 6 required. `prime_files` is empty; the REQ-075 lesson in
`tools/queue-kanban/prime-do-kanban.md` was re-read and is still accurate (its "sweep prose by line,
guard comment-carrying source by file" prescription is what this REQ implements for the second
fingerprint).

## Review

**Overall: 97%** | 2026-08-03T21:58:03Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Acceptance:** Pass — all five GREEN conditions demonstrated, each by an observed suite failure and
a matching revert.
**Suggested testing:** 2 items
**Follow-ups created:** None

### Findings

- **[Minor] The widened pattern is still an enumeration, and the comment says so.** `builder`,
  `coder`, `agent`, `REQ`; `at a time`, `at once`, `concurrently`. A phrasing like "no two REQs are
  ever claimed together" states the premise and matches nothing. This is inherent — the alternative is
  a semantic check, which a grep is not — and the mitigation is the stated trigger condition plus the
  illustrative-not-exhaustive marker, which is the repo's standing answer for exactly this. Recorded
  rather than fixed.
- **[Minor] Two guards now cover four files, and `prime-do-kanban.md`'s exemption is comment-only.**
  Nothing mechanically stops a later edit from adding the file-level negative that would break the
  lesson entry; the protection is a paragraph a maintainer has to read. Making it mechanical would mean
  asserting the *absence* of an assertion, which is worse. Left as a comment deliberately.

### Requirements Checklist

| # | Requirement | Status |
|---|---|---|
| 1 | Weak-form assertion, per-class shape | Delivered — line sweep + 2 file-level negatives |
| 2 | Widen the strong form; trigger condition; illustrative wordings | Delivered — and de-triplicated (D-01) |
| 3 | No false positive on the canonical section | Delivered — verified by dry run *and* as an explicit negative-control test case |
| 4 | Correct `actions/cleanup.md:31` | Delivered — durable reason, plus the plural fix |
| 5 | Rule on every remaining site | Delivered — 4 sites, 2 corrected, 2 left alone with the test recorded |
| 6 | Leave the badge's behavior alone | Delivered — zero Go/JS changes in the diff |

### Acceptance Testing

Seven injected-line cases run against the live suite, each reverted: six RED (including the REQ's
stated RED case, which trips both guards at once), one deliberate negative control that correctly
stays GREEN. Three of the six re-run against `HEAD`'s suite exit 0, which is the durability gap
observed rather than argued. One file-level negative fired on `tools/queue-kanban/model.go` and was
reverted. Clean tree exits 0 before and after; `git status` clean of every probe.

`shellcheck` clean, `bash -n` parses, `qualify.sh` and `scope-drift.sh` both OK.

### Suggested Additional Testing

- **Re-run the seven cases after the next edit to this block.** They are recorded as a table above but
  not as a fixture — nothing replays them automatically, and a future widening could quietly lose one.
  Worth considering whether the suite should grow a self-test harness; deliberately not built here,
  since that is a larger change than this REQ's scope.
- **A phrasing not in the enumeration** — e.g. "no two REQs are ever claimed together" — confirms the
  Minor finding above is the real boundary and not a hypothetical.

*Reviewed by review-work action (pipeline mode, in-session)*
