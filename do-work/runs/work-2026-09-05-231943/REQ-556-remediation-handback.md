# REQ-556 remediation hand-back

REQ-556 is the request that cut the debug-artifact rule prose from three action
files because `do-work-cli qualify` already enforces that rule in code. The
review scored it 63% (Partial): the prose cuts were right, the lock-in that
protects them was wrong in both directions. This hand-back covers only the
remediation.

**Branch:** `worktree-agent-REQ-556-cut-debug-artifact-prose`
**Head:** `3fa05e223c5c27162dc4c34e1b5890df8741340b` (one commit on top of `25e1b82`)
**Files changed:** `_dev/tests/audit-lockins.sh`, `skills/do-work/actions/work.md`
Both are in the request's declared write set. Nothing else was touched. Not
pushed. No release performed.

## What each finding needed and what closed it

### F1 (Important) — the lock-in fired on a pure reflow

The old block used `grep -c`, which counts lines, and the ceiling was exactly
the current count with no headroom. `review-work.md:106` carries both matched
strings on one physical line, so splitting that bullet after "no debug
artifacts —" with no word changed took the count from 2 to 3.

Fixed by counting matches instead of lines, using `rg -n -o` with the exit
status read from `$?` on the very next line. The reviewer's suggested
`grep -o … | wc -l` was not used: a pipeline whose writer can die is the trap
REQ-593 (the request that made the heavy test tier's verdict honest) spent a
whole request removing from this repository, and `wc -l` reports 0 both when
nothing matched and when the scan never ran. The chosen form has no pipeline at
all.

RED evidence: with the identical reflow applied, the old block printed
`FAIL: 3 debug-artifact rule mentions … ceiling is 2` and exited 1.
GREEN evidence: with the same reflow applied, the new block printed
`Audit lock-in regressions passed.` and exited 0.

### F2 (Important) — the comment claimed more than the code did

I chose to widen the pattern set rather than write a smaller comment, because
widening was measurable and cost nothing. Measured against the three files as
they stand today, `\b(debugger|TODO|FIXME)\b` matches zero times, so adding the
marker vocabulary that `checks.go:24` actually uses produces zero false
positives. The two old literals were also relaxed to case-insensitive and
singular-or-plural (`debug artifacts?`), which is what catches the capitalised
and singular rewrites the reviewers demonstrated.

The three marker words are split across adjacent quoted strings in the script,
exactly as `checks.go` splits them. Without that, the literal words would sit in
this file's own diff and `qualify` would flag the lock-in as a debug artifact.

Honesty was closed on the other side too: the comment now names what the block
still cannot catch — a restatement written in neither vocabulary, for example a
rationalization row about leftover print statements. Finding those is the
maintainability audit's job, not this block's.

RED evidence, one reworded restatement pasted into each of the three files in
turn, each run separately:
- `work.md` — a bullet naming only the three marker words: old block GREEN, new
  block FAIL listing 6 sites and exit 1.
- `review-work.md` — "Debug artifacts left behind in the diff": old block GREEN,
  new block FAIL listing 4 sites including `review-work.md:514 (Debug artifacts)`
  and exit 1.
- `work-reference.md` — a rationalization row saying "a debug artifact"
  (singular): old block GREEN, new block FAIL listing 4 sites including
  `work-reference.md:871 (debug artifact)` and exit 1.

### F3 (Minor) — the ratchet was one-sided

A floor was added. The pin is now exact at 3: above it a restatement came back,
below it one of the two mentions that are not restatements was cut. The two
protected reads are `review-work.md`'s standalone-review hygiene bullet (a read
that `qualify` never makes, because standalone review sees a diff `qualify`
never saw) and the emitted P-A-U template payload, which is byte-identical in
four shipped files.

An exact pin means a deliberate change to how often these files name debug
markers has to move the pin in the same commit. That is stated in the block's
comment, and it is the same contract the sibling REQ-552 block carries for its
phrase list.

RED evidence:
- `review-work.md:106` deleted outright: old block GREEN, new block
  `FAIL: the debug-artifact rule mention count fell to 1 … the pin is 3` and
  exit 1.
- The same bullet reworded so it names neither vocabulary ("no leftover prints
  — output lines …"): old block GREEN, new block the same FAIL and exit 1.

### F4 (Minor) — the failure message named no file or line

Every failure now lists `path:line (matched text)` per matching site, one per
line, matching the sibling REQ-552 and REQ-554 blocks in the same script. The
two directions get different first lines, so the reader is told whether to
delete a mention or restore one. Visible in every RED transcript above.

### F5 (Minor) — the guard covered absence but not unreadability

`rg`'s status is now captured into `debug_rule_scan_status` on the line after
the scan, and an exit greater than 1 FAILs with that status, matching both
sibling blocks. The missing-file guard is unchanged and still fires first.

RED evidence, forcing the scan to fail by breaking the regex:
- old block: grep printed `Unmatched ( or \(` three times to stderr and the
  script still exited 0 — the defect exactly as described.
- new block: three `FAIL: could not scan <path> for debug-artifact rule prose
  (rg exit 2).` lines and exit 1.
- Missing target file (`work-reference.md` moved aside): `FAIL: debug-artifact
  prose lock-in cannot read skills/do-work/actions/work-reference.md; the file
  moved and the lock-in is dead`, exit 1.

### F6 (Minor, prose) — the added sentence in work.md

Three separate fixes in the Step 6.3 sentence at `work.md:335`:

1. **The enumeration is now pinned, not just marked.** The list is introduced as
   "The codes this step judges most often, illustrative and not the full set",
   which is the house rule for lists that can go stale. On top of that, a
   companion assertion in the same lock-in block requires every `QUALIFY-*` code
   named in the three action files to still exist in `checks.go`. Completeness
   is deliberately not pinned — the prose does not claim completeness — but a
   dangling name is now impossible. The three shorthand suffixes were written
   out as full code names (`QUALIFY-DEBUG-ARTIFACT-RELOCATED`,
   `QUALIFY-OUTPUT-RELOCATED`, `QUALIFY-REPORTER-OUTPUT`) so the pin covers them
   too.
2. **"raw codes" corrected.** `--format json` ships an `automation_stop_reason`
   per finding (`resultmodel.CommandFinding`, `result_model.go:54`) that states
   in plain English why the run stopped. The sentence now says so.
3. **`debugger` relabelled.** `QUALIFY-DEBUG-ARTIFACT` now reads "a debug
   statement or unfinished-work marker the diff added"; the old wording called
   `debugger` an unfinished-work marker, which it is not.

RED evidence for the companion pin, both directions:
- Prose side: renaming `QUALIFY-UNIFY-DISARMED` to `QUALIFY-UNIFY-MISSING` in
  `work.md` gave `FAIL: skills/do-work/actions/work.md:335 names
  QUALIFY-UNIFY-MISSING, which no longer exists in … checks.go; the pointer is
  stale prose.` and exit 1.
- Code side: renaming `QUALIFY-PAU-UNCHECKED` inside `checks.go` gave the same
  shape of FAIL naming that code, and exit 1.

## Verification

All run in the worktree at the committed head:

- `bash _dev/tests/audit-lockins.sh` — `Audit lock-in regressions passed.`, exit 0
- `bash _dev/tests/action-shell-blocks.sh` — exit 0, "Shell-block lint passed:
  74 fenced blocks and 33 shipped shell files; ShellCheck enabled."
- `bash _dev/tests/contract-regressions.sh` — exit 0, "Contract regression checks
  passed."
- `shellcheck --severity=warning _dev/tests/audit-lockins.sh` — exit 0, no output
- `bash -n _dev/tests/audit-lockins.sh` — exit 0

The full canonical gate was not run, on instruction. Every ablation above was
applied to the working tree, measured, and reverted with `git checkout --`
before the next one; the tree was confirmed clean between ablations and the
committed diff is 2 files, 101 insertions, 15 deletions.

## Things a later reader should know

- **The split marker words in the pattern are load-bearing and must not be
  "tidied".** `debug_rule_marker_pattern` contains `'…(debug''ger|TO''DO|FIX''ME)…'`.
  Joining those adjacent quoted strings would put the literal marker words in the
  file and make `qualify` flag any future diff of `audit-lockins.sh` as a debug
  artifact. The comment above the pattern says this, and `checks.go:24` does the
  same thing for the same reason.
- **The code lookup is a substring match.** `QUALIFY-DEBUG-ARTIFACT` is found
  inside `QUALIFY-DEBUG-ARTIFACT-RELOCATED`, so deleting the base code while
  keeping its downgrade variant would not be caught. That pair is emitted by one
  check and only ever moves together, so the looser match was accepted rather
  than adding a boundary regex. This is stated in the block's comment.
- **The review's other Minor and Nit findings were left alone**, as they are
  outside the six findings this remediation was scoped to: the standalone-review
  coverage gap for `debugger`/`TODO`/`FIXME` in `review-work.md:106`, the
  "`advance`'s qualification gate owns the P-A-U-honesty mechanics" overstatement,
  the `work.md:292` paraphrase Nit, and `QUALIFY-NEW-FILE-UNWIRED` not being
  named. The stale VERSION note in the earlier hand-back also stands uncorrected
  here — the review already recorded that the base is 0.304.0, so the follow-on
  release is 0.304.1 across the four mirrors, and the release was not performed.

## Not done, by instruction

No release, no push, no full canonical gate, no staging or committing of this
hand-back file.
