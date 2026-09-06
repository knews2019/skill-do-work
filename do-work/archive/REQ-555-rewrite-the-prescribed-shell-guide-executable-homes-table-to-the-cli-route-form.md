---
id: REQ-555
title: '[impact-negligible] Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form'
status: completed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-06T03:17:52Z
  basis:
    - trivial short-circuit
    - Route A
    - 3-file write set
depends_on: [REQ-554]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: A
write_set: [skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T03:16:18Z
completed_at: 2026-09-06T03:57:19Z
commit: 123243554b9a1c220c27604748b8bccc0cfc08a4
release_at: 2026-09-06T03:57:19Z
---

# Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form

## What
The "Shipped executable homes" table in `skills/do-work/docs/prescribed-shell-primitives.md` assigns owned mechanics to nine `*.sh` paths that are each a 6-to-11-line `exec` shim over `do-work-cli.sh` (the mechanics moved to Go at 0.260.1), and one sentence below it says `scripts/protected-inventory.sh` "orchestrates" two check scripts, which a six-line shim cannot do. Reword the route column to the `tools/do-work-cli.sh … <subcommand>` form the toolbox rows already use and delete the orchestration sentence.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` (§ Unchecked Exit Status Reads as
  Content, § Closed Enumerations Go Stale) and re-verified every claim the request makes against
  HEAD before choosing Route A. Approach: take each route from the shim's own `exec` line, delete
  the false orchestration clause, and pin both shapes with one assertion block that fails loudly
  when its target is renamed.
- [x] **[APPLY]:** Two files changed, both in the declared write set:
  `skills/do-work/docs/prescribed-shell-primitives.md` and `_dev/tests/audit-lockins.sh`. Nothing
  else was touched. `_dev/tests/prescribed-shell-canonicalization.sh` was declared and turned out
  not to need a change; that is stated in the summary rather than papered over with an edit.
- [x] **[UNIFY]:** `git diff --stat` — `_dev/tests/audit-lockins.sh` +53, the guide +12/-10. Both
  files read in full. The guide: fourteen table rows now read the same way, each subcommand
  checked against its shim's `exec` line, and the added sentence states a condition rather than a
  count. The lock-in: every command's exit status is read, the awk pass exits 3 when the heading
  is gone and that status is checked, and `rg`'s 0/1/>1 are told apart. No debug artifacts, no
  commented-out code, no `TODO`. `bash _dev/tests/audit-lockins.sh` exits 0 and
  `bash _dev/tests/prescribed-shell-canonicalization.sh` exits 0; four ablations each print the
  intended FAIL line and exit 1.

## Why
The guide is the pointer target from 16 shipped files (ratchet-enforced) and it currently misroutes readers to shims and describes an orchestration that no longer exists.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 7, sweep_key `stale-shell-ownership-prose`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -5. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- Seven rows naming `scripts/show-commit-diff.sh`, `add-local-git-exclude.sh`, `atomic-download.sh`, `capture-screenshot.sh`, `run-blocked-check.sh`, `protected-inventory.sh`, `stage-exact-deletion.sh` (9-11 lines each, all `exec … do-work-cli.sh`).
- Two rows naming `../../do-work-knowledge/scripts/lexical-memory-recall.sh` and `install-memory-hooks.sh` (6 lines each).
- The sentence `which orchestrates \`tools/checks/uncommitted-inventory.sh\` and \`tools/checks/associate-files.sh\`` is false at HEAD: delete it.
- The shims themselves are not touched (retained by the 0.260.1 decision); only the guide's description of them changes.
- Reproduce at dc8a64e3: `awk 'NR>=9 && NR<=22 && /\.sh`/' skills/do-work/docs/prescribed-shell-primitives.md | grep -oE '[^`]*\.sh' | while read -r p; do case "$p" in ../../*) fp="skills/${p#../../}";; *) fp="skills/do-work/$p";; esac; echo "$(wc -l < "$fp" | tr -d ' ') lines $fp $(grep -c do-work-cli "$fp") cli-exec"; done; rg -n 'which orchestrates' skills/do-work/docs/prescribed-shell-primitives.md`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Same guide and ratchet as REQ-554: land after it so the heading and pointer counts are re-baselined once.
- Lock-in limit: shim rows in the executable-homes table: 0 after this REQ (today 9).

## Dependencies
Depends on REQ-554, which already edits this guide and re-baselines the ratchet, so this REQ is a table rewrite only.

## Builder Guidance
Firm: the route column names the CLI subcommand, not a shim.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine shim rows (6-11 lines each, all `do-work-cli` execs) and the orchestration sentence.
**GREEN when:** The table has zero `.sh` rows for mechanics owned by Go, the sentence is gone; the lock-in pins shim rows in that table at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for stale-shell-ownership-prose.*

## Triage

**Route: A** — Build directly.

**Reasoning:** Every claim the request makes was re-checked against HEAD and all of them hold, so there
is nothing left to discover. The nine paths in the "Shipped executable homes" table are 6-to-11-line
`exec` shims over `tools/do-work-cli.sh`, each naming its own subcommand on that exec line, so the
route column can be rewritten mechanically from the shims themselves. The orchestration sentence is
still present, and it is still false: `scripts/protected-inventory.sh` is six lines, and the two files
it is said to orchestrate — `tools/checks/uncommitted-inventory.sh` and
`tools/checks/associate-files.sh` — are themselves 9-line and 19-line compatibility launchers over the
same Go command, which no code in `internal/` ever launches. The work is one table rewrite, one
sentence deletion and one counting assertion.

**Planning:** Skipped.

**REQ-554 landed first, as this request requires.** It rewrote the paragraph the orchestration sentence
sits in and left that sentence untouched by an explicit decision, so this request edits settled text
rather than racing it.

## Plan

**Planning not required** — Route A: the write set, the target table and the sentence to delete are all
named by the request, and every claim was re-checked against HEAD before the route was chosen.

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `_dev/tests/audit-lockins.sh` (modified)

**What was done:** The nine `.sh` rows in the "Shipped executable homes" table now name the
`tools/do-work-cli.sh … <subcommand>` route, taken from each shim's own `exec` line rather than
guessed, so all fourteen rows read the same way. The false orchestration clause is deleted. One
assertion block in `_dev/tests/audit-lockins.sh` pins both shapes.

**One sentence was added, and it is earned.** Rewriting the route column alone leaves the guide
contradicting itself: the table would name only CLI subcommands while the paragraph below it still
says the inventory ships behind `scripts/protected-inventory.sh`, which is true — `actions/commit.md`
and `../../do-work-toolbox/actions/inspect.md` invoke that path today. A reader following the table
would conclude the launchers are gone and could delete one. The added sentence states the condition
(where a route also ships a retained launcher of the same name) rather than listing which nine, so it
cannot go stale as rows move, and it says plainly that a behaviour change is made in the command and
never in the launcher.

**The expected net line delta was -5; the actual guide delta is +2** (12 insertions, 10 deletions).
The nine rewritten rows are net zero, the deleted clause was part of a longer line, and the added
sentence is +2 with its blank line. The estimate assumed a pure deletion; the correctness of the guide
was preferred to the number.

`_dev/tests/prescribed-shell-canonicalization.sh` is in the declared write set and was not changed: it
checks only that the nine shims exist and are executable, which this request does not touch. Declaring
a file and then not needing it is recorded rather than papered over.

## Decisions — implementation

- **D-01 — the route column names the subcommand each shim actually execs. DECIDE & STATE.** Every one
  of the nine was read: `show-commit-diff`, `add-local-git-exclude`, `atomic-download`,
  `capture-screenshot`, `run-blocked-check`, `protected-inventory`, `stage-exact-deletion`,
  `lexical-memory-recall`, `install-memory-hooks`. None was inferred from the file name.
- **D-02 — the retained launchers are stated once, keyed on a condition, not counted. DECIDE &
  STATE.** "Where a route above also ships a retained `scripts/*.sh` launcher of the same name" holds
  as rows are added or removed. A sentence saying "the first nine rows" would be a hand-maintained
  list, which `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale bans.
- **D-03 — the shims themselves are untouched, and so is the sentence that routes readers to one.**
  The 0.260.1 decision retained them and shipped actions still invoke them by path. Rewriting those
  invocations is a different change with a different blast radius, and this request does not name it.
- **D-04 — the lock-in pins two shapes in one block because they are one defect. DECIDE & STATE.** A
  shim in the route column and the claim that a six-line launcher orchestrates two scripts are both
  "the prose describes a shell ownership that no longer exists". The request asks for one assertion,
  and splitting them would have produced two.
  **Corrected after review.** The first version of this decision claimed both halves were keyed on the
  claim rather than the wording. Neither was. The route half tested for a backticked `.sh` in one cell
  and the orchestration half was a 19-character fixed string; the review put the same defect back seven
  ways past the first and three ways past the second. Both are now keyed on the condition — see the
  remediation qualification.
- **D-05 — a missing guide or a renamed heading fails rather than reads clean. DECIDE & STATE.** The
  awk pass exits 3 when it never sees the heading, and that status is checked. A ratchet that goes
  quiet when its target is renamed is the failure mode REQ-552 and REQ-554 both hit in this same file.

## Discovered Tasks

- **Row 13's Mechanics cell describes shell the Go command does not run.** It credits `run-blocked-check`
  with "GNU timeout selection and isolated stock-Bash process-group timeout/cleanup"; the Go
  implementation does neither. Equally stale before this change, and the request's constraint is "scope
  is exactly this finding class", so it is captured as **REQ-595** rather than folded in. That request
  also asks for the other thirteen Mechanics cells to be checked in the same pass, because the one row
  that was checked says nothing about the rows that were not.

## Qualification

**Passed.** Read from the range `a49a542f..ad8a8050`, two files, 65 insertions and 10 deletions.
Canonical `qualify` is satisfied.

- **Every claim in the request holds at HEAD, and each was checked rather than assumed.** The nine
  paths are 9, 10, 9, 9, 11, 6, 9, 6 and 6 lines, and every one of them `exec`s
  `do-work-cli.sh --format text <subcommand> "$@"`. The orchestration clause was still present. This is
  the opposite of the sibling REQ-556, whose stated baseline did not survive contact with HEAD.
- **The orchestration clause is false in two independent ways, both verified.**
  `scripts/protected-inventory.sh` is six lines, and `tools/checks/uncommitted-inventory.sh` and
  `tools/checks/associate-files.sh` are themselves 9-line and 19-line compatibility launchers over the
  same command. Nothing under `internal/` launches either of them: the only matches in the Go tree are
  in `inventory_test.go`, which drives the launcher chain deliberately to hold the launcher contract.
- **One sentence was added, and the addition is declared rather than smuggled.** `maintenance.md` asks
  for a concrete case that fails without it, and that case is the paragraph immediately below the table,
  which still routes readers to `scripts/protected-inventory.sh` because that is genuinely what
  `commit.md` and `inspect.md` invoke. Without the sentence the table and the paragraph contradict each
  other. It is stated in the Implementation Summary as an addition and is the reason the net line delta
  is +2 rather than the expected -5.
- **The declared write set is a ceiling that was not filled.**
  `_dev/tests/prescribed-shell-canonicalization.sh` was declared and not changed. It still exits 0.
  **Corrected after review**: the first version of this bullet said it "only asserts the nine shims
  exist and are executable". It is 169 lines and also pins twelve required headings in this same guide,
  sixteen pointer sites across shipped files, a no-direct-curl rule over `tools/`, and stale-prose
  scans. None of those is touched by a route-column rewrite, which is why it needed no change — but the
  reason is "nothing this request edits is in its coverage", not "it covers almost nothing".

### Remediation qualification (after review)

**Passed.** Remediation merge range `5b0f7b8f..12324355`, two files, both already in the declared write
set. The review scored 62% and asked for changes; all four of its confirmed non-cosmetic findings are
closed in code, and three claims in this record are corrected rather than defended.

- **The sentence this request added was false, and it was false in the exact way the request exists to
  remove.** It said the retained launchers "do nothing but `exec` the subcommand" and that behaviour
  changes are "never in the launcher". Seven of the nine translate a legacy positional call into the
  flags the subcommand requires — `scripts/show-commit-diff.sh <commit>` becomes `--commit <commit>` —
  so the same arguments passed straight to the command are rejected with `unknown option`. And
  `scripts/protected-inventory.sh` sets `DO_WORK_COMPATIBILITY_SHIM=1`, which eighteen non-test sites
  under `internal/` read, and which selects the `<tag>\t<path>` output that `actions/commit.md` and
  `../../do-work-toolbox/actions/inspect.md` parse one row per file from. Adding false prose to a
  request whose purpose is deleting false prose is the worst available outcome, and the reviewer was
  right to score it high. The sentence now states what is true and says which of the two is invoked and
  which owns the mechanics, which also closes the reviewer's separate finding that the guide had begun
  giving two canonical answers.
- **The route check is now positive and reads the file it names.** It was a negative test — "no
  backticked `.sh` in cell two of a row under this heading" — which the review evaded seven ways: an
  argument inside the backtick span, the table's own `…` house style, no backticks at all, the shim
  moved into the Mechanics cell, one leading space, an inserted sub-heading, and an emptied table. It
  now finds each table by its own header row, scans every cell, and flags a named `.sh` path when that
  file is itself a do-work-cli launcher — read from the file, so a genuinely shell-owned route would
  pass and a renamed shim would not. A row naming a file that does not exist fails, and a table with no
  rows fails.
- **The orchestration check is keyed on the co-occurrence, not on a verb.** It matched the literal
  `which orchestrates`; "that orchestrates", "which coordinates" and "a wrapper that drives" each
  restored the identical false claim green. It now fails any line naming `scripts/protected-inventory.sh`
  beside either check launcher it was said to drive. That is deliberately broad: the two names appearing
  in one sentence is the defect shape, and a future sentence that legitimately needs both would be a
  conscious edit to this check.
- **Three record claims are corrected in place** rather than left standing: D-04's claim that both
  halves were keyed on the claim, the Testing bullet's claim that the check matched the claim rather
  than the wording, and the Qualification's understatement of what
  `_dev/tests/prescribed-shell-canonicalization.sh` covers.
- **One finding is recorded and not fixed.** Row 13's Mechanics cell attributes GNU-timeout selection
  and stock-Bash process groups to a Go command that does neither. It was equally stale before this
  change, and the request's constraint is "scope is exactly this finding class", so it is a discovered
  task rather than a widening.

## Testing

**The lock-in was proven red four ways before it was accepted**, and the guide was restored to the
green state after each.

- One shim row returned: `prescribed-shell-primitives.md:9 routes owned mechanics to a shim
  (\`scripts/show-commit-diff.sh\`)`, exit 1.
- All nine returned: nine `FAIL:` lines, one per row, exit 1. The count is what the request pins at 0.
- The orchestration clause returned, reworded to a shorter sentence than the deleted one: caught, exit
  1. ~~The check matches the claim, not the original wording.~~ **False, and corrected after review:**
  the shorter sentence still contained the literal `which orchestrates`, which is all the check
  matched. Three rewordings that keep the same false claim — "that orchestrates", "which coordinates",
  "a wrapper that drives" — each passed. The ablation could not have found this because it kept the
  string the check was looking for.
- The table heading renamed from `## Shipped executable homes` to `## Executable homes`: `the "##
  Shipped executable homes" heading is gone ... (awk exit 3); the route ratchet cannot run`, exit 1.
  This is the case that matters most — a rename is how a ratchet goes quiet instead of red.

Green at HEAD: `bash _dev/tests/audit-lockins.sh` prints `Audit lock-in regressions passed.` and exits
0; `bash _dev/tests/prescribed-shell-canonicalization.sh` prints its passed line and exits 0. The full
fast gate is run once for the batch rather than per request.

### Remediation testing (after review)

**Thirteen ablations, each run against the real file with it restored from a green copy between runs,
and the last line of the sequence is the restored green.**

Seven route-row evasions the review demonstrated, all of which used to pass, all now `EXIT=1` with the
same message naming the path and its line: a bare backticked shim path (the control), an argument
inside the backtick span, the `…` house style, no backticks, the shim moved into the Mechanics cell,
one leading space, and a sub-heading inserted between the heading and the table.

Three structural cases: the table emptied of rows gives `the executable-homes table ending at line 8
has no rows; an empty table names no home`; the table's header row renamed gives `no "| Canonical
executable route |" table header remains ... (awk exit 3); the route ratchet cannot run`; a route
naming a script that does not exist gives `names scripts/does-not-exist.sh as a canonical route and no
such file exists`.

Three orchestration rewordings that each used to pass — "that orchestrates", "which coordinates", "a
wrapper that drives" — now fail with the same message, as does the original wording.

Green after restore: `bash _dev/tests/audit-lockins.sh` prints `Audit lock-in regressions passed.` and
exits 0. `bash _dev/tests/prescribed-shell-canonicalization.sh` exits 0. `bash -n` and
`shellcheck --severity=warning` on the lock-in both exit 0.

**Fast gate at the remediation revision:** `Maintainer verification passed.`, exit 0, wall 85s, with
`do-work-cli` at 784 tests. One `SKIP` line, the heavy-only one every fast run prints.

## Review

**Overall: 62%** | 2026-09-06T04:15:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 65% |
| Code Quality | 60% |
| Test Adequacy | 35% |
| Scope | 58% |
| Risk | Medium |
| Acceptance | Request changes |

**Verdict: Request changes** — the table rewrite itself is correct and the nine subcommands each match
their launcher's own `exec` line, but the change added one sentence of shipped prose that is false for
seven of the nine launchers it describes, and neither half of the lock-in held the property it claimed.
Three reviewers worked independently; the synthesis reproduced every finding before accepting it.

Where the reviewers disagreed, and what was picked:

- The added sentence. Reviewer 1 scored it a positive for stating a condition instead of listing nine
  rows. Reviewers 2 and 3 called it a high-severity false claim. Settled by reading all nine launcher
  bodies and running `protected-inventory` through both routes against the same dirty tree: seven of
  nine do more than `exec`, and the two routes share no output line. Picked reviewers 2 and 3 —
  compliance with one rule does not make a sentence true.
- How many ways the route ratchet could be evaded. Reviewer 2 named one, reviewer 3 named three,
  reviewer 1 named eleven across both halves. Settled by running each candidate in a sandbox copy with
  the guide restored from git between runs: seven route-row evasions reproduce.
- Whether "the route the table now names is not runnable" is a finding. Rejected as framed: the `…` in
  a route cell stands for global options and is the same shape the five pre-existing toolbox rows
  already used, so a route cell was never a paste-ready command in this table.

**Important findings (each with its recorded impact token):**

- The added sentence is false for seven of the nine launchers, and materially wrong for one. Six
  launchers rewrite positional argv into the flags the subcommand requires; `protected-inventory.sh`
  sets `DO_WORK_COMPATIBILITY_SHIM=1`, which eighteen non-test sites under `internal/` read and which
  selects the `<tag>\t<path>` output `commit.md:59` and `inspect.md:67` parse one row per file from.
  Reproduction: `bash scripts/show-commit-diff.sh ad8a8050` exits 0 while
  `bash tools/do-work-cli.sh --format text show-commit-diff ad8a8050` exits 2 with
  `unknown option ad8a8050`. — impact-user-visible → fixed in remediation
- The shim-row half of the lock-in is evadable seven verified ways, including the table's own `…` house
  style. It required a backtick immediately after `.sh`, looked only at cell two, only at lines starting
  with `|`, and only inside one heading's window. — impact-user-visible → fixed in remediation
- The orchestration half pins a 19-character literal, so one changed word restores the identical false
  claim green: "that orchestrates", "which coordinates", "a wrapper that orchestrates" all passed. The
  builder's own ablation kept the literal, so it could not have found this. — impact-user-visible →
  fixed in remediation
- The guide gave two different canonical answers for the same three mechanics: the table named the
  command while the prose sections at lines 43, 47, 80, 99 and 115 still prescribed the launcher with
  positional arguments. — impact-negligible → fixed in remediation

**Minor findings:**

- The negative check would have failed spuriously on a legitimate shell-owned route row and on a true
  sentence containing "which orchestrates". Neither exists in the guide today, so it was scored low. —
  impact-negligible → fixed in remediation, which reads the named file rather than matching its name
- Row 13's Mechanics cell attributes GNU-timeout selection and stock-Bash process groups to a Go
  command that does neither. Equally stale before this change. — impact-negligible → discovered task
- This record understated what `_dev/tests/prescribed-shell-canonicalization.sh` covers. —
  impact-negligible → corrected in place
- Two consecutive blank lines at the end of the new block where the file uses one. — impact-negligible
  → fixed in remediation

**Requirements checklist:**

- [x] The nine `.sh` rows read `tools/do-work-cli.sh … <subcommand>`, each subcommand matching its
  launcher's own `exec` line — delivered
- [x] The false orchestration clause is gone and line 43 still parses — delivered
- [x] The shims themselves are untouched; the diff is exactly two files — delivered
- [x] No test file changed beyond the one lock-in — delivered
- [x] The lock-in is green today, and the loud structural guards fire — delivered
- [ ] → [x] "Red the moment the number regrows" — **not delivered at review, delivered in remediation.**
  It was red only for a byte-exact repeat of the deleted shapes; thirteen ablations now cover the seven
  route evasions, three structural cases and three orchestration rewordings.
- [ ] → [x] "Scope is exactly this finding class" — **not delivered at review, delivered in
  remediation.** The added sentence was prose the request did not ask for and was false; it now states
  what is true and resolves the contradiction the review found.

**Acceptance testing**

**Result: Partial at review, Pass after remediation.** Three reviewers each ran the lock-in against
sandbox copies with the guide restored between runs. Every evasion they demonstrated was re-run against
the remediated check and now fails; the restored tree is green, and the fast gate exits 0.

**Follow-ups created:** 1 — the stale `run-blocked-check` Mechanics cell, recorded as a discovered task
in the same finding class rather than widening this request.

*Reviewed by review-work action*

## Lessons Learned

- **Reading a script's `exec` line is not reading the script.** Nine launchers were checked and all
  nine execed the subcommand the table now names, which is what the route column needed. But six of
  them rewrite positional arguments into flags first and one sets an environment variable that changes
  the command's output shape, and a sentence was written claiming they "do nothing but exec". The
  question answered was "which subcommand does it call"; the sentence made a claim about "what else
  does it do", and that was never checked. When a sentence generalizes over a set, check the sentence
  against the whole set, not against the fact that put the set on screen.
- **An addition made to prevent one misreading can create a worse one.** The sentence existed because
  the rewritten table would otherwise contradict the paragraph below it. That reasoning was sound and
  the sentence still shipped a false claim into the file this request exists to make true. An earned
  addition still has to be verified like any other claim.
- **A negative ratchet forbids the spelling you deleted.** "No backticked `.sh` in cell two under this
  heading" was red for exactly the shape that was removed and green for seven others, including the
  table's own `…` house style. The positive form — every route row must name the command, and a named
  `.sh` file that is itself a launcher is an offender — cannot be evaded by punctuation, indentation,
  a moved cell, or a second table, because it asks what the row means rather than how it is spelled.
- **A check that reads the file it names beats one that matches the name.** Deciding "is this route a
  shim?" by opening the script and looking for `do-work-cli` removed the whole false-positive class in
  the same change: a genuinely shell-owned route passes, and a shim renamed to something innocuous
  still fails.
- **An ablation that keeps the string the check looks for proves nothing.** The orchestration clause
  was reworded for the ablation but the reworded version still contained `which orchestrates`, so the
  test passed for the wrong reason and the record claimed it "matches the claim, not the original
  wording". Ablate the property, not the sentence: change the words the check is not allowed to depend
  on.
- **Three independent reviewers disagreed, and the disagreement was the signal.** One scored the added
  sentence as a positive. The two who called it false had read the launcher bodies. The one who had
  read more of the code was right, which is the argument for running three lenses rather than one.

## Orientation

`skills/do-work/docs/prescribed-shell-primitives.md` has one "Shipped executable homes" table and every
row now names `tools/do-work-cli.sh … <subcommand>`. The table says where a mechanic is **owned**.
Where a row also has a retained `scripts/*.sh` launcher of the same name, that launcher is what the
guide's own prose and the shipped actions **invoke**, and the two are not interchangeable: six
launchers translate a legacy positional call into the subcommand's flags, and
`scripts/protected-inventory.sh` sets `DO_WORK_COMPATIBILITY_SHIM=1` to select the `<tag>\t<path>`
output that `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` parse. Changing the
command's flags or output shape changes the launcher's contract; fix both together.

The guard is Finding 7 in `_dev/tests/audit-lockins.sh`, beside REQ-552's, REQ-554's and REQ-556's. Two
conditions, neither keyed on a spelling:

- **Route rows.** Each table is found by its own `| Canonical executable route |` header row, so a
  sub-heading or a second table cannot hide one. Every cell of every row is scanned for a `.sh` path;
  each is resolved the way the guide writes them (`../../<pkg>/…` is a sibling package under `skills/`,
  anything else is relative to the do-work package) and the file is read. A path whose file contains
  `do-work-cli` is a shim row and fails. A path with no file fails. A table with no rows fails.
- **The orchestration claim.** Any line naming `scripts/protected-inventory.sh` beside
  `tools/checks/uncommitted-inventory.sh` or `tools/checks/associate-files.sh` fails. That is
  deliberately broad — the two names in one sentence is the defect shape — so a future sentence that
  legitimately needs both is a conscious edit to this check.

Recorded and unfixed: row 13's Mechanics cell attributes GNU-timeout selection and stock-Bash process
groups to a Go command that does neither. It was stale before this request and is captured as a
discovered task in the same finding class.
