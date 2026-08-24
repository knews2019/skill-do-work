---
id: REQ-339
title: "[impact-rule-change] Addendum: count every prescribed-shell case the runner reports"
status: completed
completed_at: 2026-08-24T09:22:00Z
claimed_at: 2026-08-23T23:08:00Z
status_changed_at: 2026-08-23T22:32:23Z
created_at: 2026-08-23T19:30:00Z
user_request: UR-065
addendum_to: REQ-325
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
route: B
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - _dev/tests/prescribed-shell-harness.sh
---

# Count Every Prescribed-Shell Case the Runner Reports

## What

`prescribed_shell_finish` counts a case file's cases with `grep -cE '^# [a-z0-9][a-z0-9-]*: '`, so a
header with a space or a comma before its colon is invisible. `generate-report-image.sh` reports 7
cases and contains 9: `# generate-report-image caller contract: …` and
`# generate-report-image, interrupted directly: …` are both uncounted. The aggregate figure the
runner prints ("96 named script cases across 17 per-script files") inherits the same undercount.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-shell-commands.md` (§ Closed Enumerations Go Stale decided the approach). Rejected renaming the two odd headers; chose to widen the rule and anchor it on the case file's own basename, with the rule stated in one sourced file so the harness and the aggregate cannot disagree.
- [x] **[APPLY]:** Landed the refactor first with the OLD expression to prove it changed nothing, then swapped the expression. Three files beyond the declared write set, each declared as a scope extension with its reason (see `## Decisions`).
- [x] **[UNIFY]:** Audited by the orchestrator against the merged range `47c61cf..97a5ca8`, not from the builder's report: read all four files' diffs, grepped the range for debug artifacts (one hit, a `mktemp` XXXXXX template — not an artifact), re-derived both changed counts independently (9 and 4) and confirmed the other fifteen files are unmoved.

## Why

REQ-234 replaced a hand-maintained literal with a derived count precisely so the figure would stop
being a remembered number reported as a measured one. A regex that silently skips real cases is the
same untruth with a different cause, and it is the more durable kind: nothing fails, the number just
reads low forever. A reader using the count to judge coverage is misled in the direction that
matters.

## Detailed Requirements

- The count matches the number of case blocks the file actually contains.
- Decide and state what a case header *is* — the counting rule lives in a comment beside the count
  (REQ-234's contract), so whatever shape is accepted has to be written down there too.
- Sweep every file under `_dev/tests/prescribed-shell-cases/` for headers in the uncounted shapes
  and confirm the new count agrees with a hand count on at least the two files that change.
- The aggregate figure in `prescribed-shell-scripts-behavior.sh` moves in step (it greps the same
  pattern).

## Constraints

- `_dev/primes/prime-shell-commands.md` governs any shell that ships or gates. Read it first — in
  particular § *Closed Enumerations Go Stale*: prefer widening the condition over enumerating the
  two header spellings that happen to exist today.
- Renaming the offending headers to fit the existing regex is the cheaper fix and is on the table —
  but say why, because it leaves the next author free to write a header the count skips.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-325: the prescribed-shell case count
  skips any header with a space or comma before its colon, so `generate-report-image.sh` reports 7
  of its 9 cases and the runner's aggregate is low by the same amount. Should I process this as a
  new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - [2026-08-23] User approved via clarify: the reported case count must match what the file
    holds — a figure that reads low forever is the same untruth REQ-234 removed when it
    replaced the hand-maintained literal. Nothing put out of scope; renaming the two odd
    headers instead of widening the rule stays on the table as the REQ's Constraints say,
    provided the reason is stated.

## Red-Green Proof

**RED prompt/case:** `grep -cE '^# [a-z0-9][a-z0-9-]*: ' _dev/tests/prescribed-shell-cases/generate-report-image.sh`
prints 7 while `grep -c '^# generate-report-image'` on the same file prints 9.

**GREEN when:** the count reported by that file equals its hand-counted case blocks, and a case
header written in either previously-uncounted shape is counted.

**Validation:** Inferred during REQ-325's implementation — a Discovered Task, not a user request.

---
*Source: Discovered Task, REQ-325 (UR-065).*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is stated exactly (the reported count must match what the file holds) but the counting rule's shape was undecided, and the sweep across `_dev/tests/prescribed-shell-cases/` had to find where else the pattern is applied. Clear what, unknown where.

**Planning:** Not required.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `_dev/tests/prescribed-shell-harness.sh` (modify) — replace the counting regex with a call to the shared rule
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — the aggregate greps the same pattern and moves in step
- `_dev/tests/prescribed-shell-case-count.sh` (new) — one home for the counting rule
- `_dev/tests/contract-regressions.sh` (modify) — lock-in for the undercount

**Files I will NOT touch:** the 17 case files under `_dev/tests/prescribed-shell-cases/` — renaming their headers was the rejected alternative (D-01), so their content stays as written.

**Acceptance criteria (restated from REQ):**
- [x] The count matches the number of case blocks the file actually contains
- [x] What a case header *is* is decided and stated in a comment beside the count
- [x] Every file under `_dev/tests/prescribed-shell-cases/` swept for uncounted header shapes
- [x] New count agrees with a hand count on at least the two files that change
- [x] The aggregate in `prescribed-shell-scripts-behavior.sh` moves in step

## Implementation Summary

**Files changed:**
- `_dev/tests/prescribed-shell-case-count.sh` (new)
- `_dev/tests/prescribed-shell-harness.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** The counting rule moved out of two hand-synchronized literal regexes into one sourced function, `count_named_case_headers`, and widened from `^# [a-z0-9][a-z0-9-]*: ` to `^# ${script_under_test}[^.:]*: ` — a header is a column-zero comment opening with the case file's own basename and reaching a colon before any period. Two files' counts corrected (`generate-report-image` 7→9, `generate-report-image-batch` 2→4), aggregate 96→100, with a lock-in in `contract-regressions.sh` that asserts 3 against a fixture carrying all three header spellings plus three near-misses.

## Decisions

<!-- D-XX counter: last used D-05. Next decision: D-06. -->

- **D-01 — Widen the rule rather than rename the two headers. DECIDE.** Renaming the two odd headers to fit the old regex is a two-line diff, but it fixes the two spellings that exist today and leaves the count silently skipping whatever the next author writes — the same shape REQ-234 removed. Both qualifiers also carry real information: `caller contract` distinguishes a multi-job contract case from the single-invocation cases around it. The REQ's Constraints put renaming on the table provided the reason was stated; this is the reason it was not taken.
- **D-02 — Anchor on the script under test rather than a generic token. DECIDE.** The obvious widening `^# [a-z0-9][a-z0-9-]*[^:]*: ` over-counts badly: 23 lines across the suite are wrapped prose continuations carrying a colon, 12 of them starting with a lowercase word, and all would read as cases. Anchoring on the case file's own basename is what makes an open-ended qualifier safe — the rule becomes "names its script, then a colon" rather than an enumeration of the separators used so far, which is what `prime-shell-commands.md` § Closed Enumerations Go Stale asks for.
- **D-03 — Stop the qualifier at the first period. DECIDE.** Without it, a wrapped continuation line (`# qualify.sh should be updated to match (REQ-250: ...`) poses as a header and pushes qualify.sh from 21 to 22. One character class, stating a real property: a header's qualifier is a phrase, not a sentence.
- **D-04 — One shared definition in a new sourced file rather than the same expression in two places. DECIDE.** The harness and the runner previously carried identical literal regexes kept in step by hand. Under the new rule they are no longer the same code shape (single file vs loop), so duplicating would be worse than before — and the point of this REQ is that a reported number must not be able to disagree with the files. Rejected alternative: having the runner parse each case file's printed summary line, which deletes a copy of the rule but adds a parser coupled to a print format and forces stdout buffering that reorders output against the live `FAIL:` lines on stderr.
- **D-05 — Lock-in placed inline in `contract-regressions.sh`. DECIDE.** This bug class is invisible without one: reverting the fix drops the aggregate to 96 and nothing else fails. The fixture deliberately carries three near-misses (capitalised prose, a lowercase continuation line, a sentence naming a script before a colon) so the test also fails against an over-broad widening, not only against a narrowing.

## Scope Extensions

Three files beyond the declared write set of `_dev/tests/prescribed-shell-harness.sh`, each with its reason:

- `_dev/tests/prescribed-shell-scripts-behavior.sh` — pre-authorised by the REQ's own requirements; the aggregate greps the same pattern and had to move in step. Now 100.
- `_dev/tests/prescribed-shell-case-count.sh` (new) — the consequence of D-04. One home for the rule is what makes the two counters unable to disagree.
- `_dev/tests/contract-regressions.sh` — the consequence of D-05. A silent-undercount fix with no lock-in regresses silently by construction.

## Discovered Tasks

- `_dev/tests/prescribed-shell-cases/qualify.sh` lines 127-130 and 400-404 are wrapped prose continuations opening with `qualify`. They are excluded today only because the qualifier stops at the first period. A future continuation line that opens with the script name and reaches a colon with no period in between would be counted as a case. A stronger discriminator (requiring the header to open a comment block — its preceding line not being a comment) would close that, at the cost of moving the counter from grep to awk. Not worth it for the current corpus; worth revisiting if a false positive appears.
- Nothing validates that a case file's basename matches the script it actually covers, so a case file containing headers for a different script would be counted under its own name. No instance exists — all 17 files use their own basename exclusively.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh`, `bash _dev/tests/contract-regressions.sh`, `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing — gate exit 0 (verified twice: after this REQ merged, and again after REQ-340 merged on top)

**Red-green validation:**
- `contract-regressions.sh` case-count lock-in: ✗ before implementation (reports 1 of the fixture's 3 headers) → ✓ after
- The REQ's own stated RED — `generate-report-image.sh` reporting 7 while holding 9 — resolved: the file now reports 9

**Mutation evidence (orchestrator, not the builder):** reverting `count_named_case_headers` to the old
regex makes the lock-in report 1 of 3 and FAIL with exit 1. The test is load-bearing, not decorative.

**New tests added:**
- `_dev/tests/contract-regressions.sh` — one lock-in block asserting exactly 3 against a fixture carrying three header spellings plus three near-misses

**Cross-REQ evidence:** REQ-340 merged after this and added a fifth case to
`generate-report-image-batch.sh`. The per-file line and the aggregate moved from 4/100 to 5/101 with
no edit to either counter — the derivation this REQ installs, demonstrated live rather than argued.

*Verified by work action*

## Review

**Overall: 89%** | 2026-08-24T09:12:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:**
- The stated rule and the implemented rule differ at the script-name boundary: the comment says a header "opens with the name of the script", while `^# ${script_under_test}[^.:]*: ` accepts any token *beginning* with that name, so a line like `# qualifying the untracked scan below is prose: it must not count.` counts as a case (reproduced: 2 against a 1-case fixture). No corpus line exercises it today and the lock-in carries no near-miss of this shape. — `impact-negligible` → **REQ-355** created.

**Minor findings:** 4 (report only) — the prose omits the space the regex requires after the colon; the basename is interpolated into an ERE with no stated regex-safe assumption (a case file named `a+b.sh` would silently count 0); `write_set` frontmatter still lists one file while four were delivered, so the board's overlap badge understated the surface during the fan-out wave (no collision occurred — siblings verified disjoint); `rm -rf` at one fixture cleanup omits the `--` its eight neighbours use.

**Acceptance:** Pass — runner exit 0 with per-file lines summing exactly to the printed aggregate; `contract-regressions.sh` exit 0; `shellcheck --severity=warning` exit 0 on all three shell files.

**Restatement sweep:** done, nothing stale. Every live restatement of the old rule moved with the diff; remaining hits are dated history (CHANGELOG figures, archived hand-backs) that record what the number was at the time.

**Counting audit:** the reviewer sourced the shipped rule, listed all 100 matched lines and read every one — each a genuine case header opening a fixture block — then listed the 19 column-zero comments it rejects, all plainly wrapped prose. No over-count and no undercount in the current corpus.

**Follow-ups created:** REQ-355

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Landing the refactor first with the OLD expression, proving it changed nothing, then
swapping the expression — the two-step made it impossible to confuse a refactor regression with a
rule change. Anchoring the widened rule on the case file's own basename is what let the qualifier
stay open-ended without over-counting the 23 wrapped prose lines that carry a colon.

**What didn't:** The obvious widening `^# [a-z0-9][a-z0-9-]*[^:]*: ` over-counts badly — 12 wrapped
continuation lines start with a lowercase word and would all have read as cases. Enumerating the two
header spellings that exist today was the cheaper fix and was explicitly on the table; it was
rejected because it leaves the next author free to write a header the count skips, which is the
defect this REQ removes.

**Worth knowing:** The rule as landed is still one boundary short — it accepts a token *starting*
with the script name, which is the over-count mirror of the bug fixed here (REQ-355). Also: a case
file whose basename carries regex metacharacters would silently count 0, because the basename is
interpolated straight into an ERE. Neither is reachable with the current 17 files.

## Orientation

The prescribed-shell test suite now derives its reported case count from one rule in one file rather
than two hand-synchronized regexes. Lives in `_dev/tests/`, the repo's own gate tooling — no shipped
skill code changed. A reader who trusts the runner's "N named script cases across 17 per-script
files" line is now trusting a measured figure rather than an approximation of one.
