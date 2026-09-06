---
id: REQ-556
title: '[impact-negligible] Cut the debug-artifact rule prose that do-work-cli qualify already enforces'
status: completed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-06T02:27:50Z
  basis:
    - Route B
    - 4-file write set
    - 4 acceptance criteria
    - cross-route regression gates
route: B
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, skills/do-work/actions/work-reference.md, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T02:27:19Z
completed_at: 2026-09-06T03:15:36Z
commit: 2acc465a5fbb7e371a45b63caaa95198c6ecb3f0
release_at: 2026-09-06T03:15:36Z
---

# Cut the debug-artifact rule prose that do-work-cli qualify already enforces

## What
The debug-artifact and P-A-U-honesty rule that `do-work-cli qualify` enforces at Step 6.3 (finding codes `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED`, `QUALIFY-UNIFY-DISARMED` in `internal/corehelpers/checks.go`) is written a second time as an agent instruction at five prose sites across `work.md`, `review-work.md` and `work-reference.md`. Keep one sentence in `work.md` Step 6.3 naming the three finding codes; cut the other four sites to a pointer. Keep the judgment prose the Go check explicitly defers ("judge entry-point or dynamic-wiring exceptions"), which has no duplicate.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Both `prime_files` read before either shell or action-file edit, plus the crew rules
  and the exploration. Approach: follow the exploration's site-by-site list rather than the request's
  stale count, locate every edit by text because REQ-486 had moved one anchor, and keep the two
  mentions that are not restatements.
- [x] **[APPLY]:** Four files, exactly the declared `write_set`, across the original change and its
  remediation.
- [x] **[UNIFY]:** `git diff --stat` reports four files for the original and two for the remediation,
  all declared. Linters: `bash _dev/tests/action-shell-blocks.sh` — `Shell-block lint passed`, exit 0;
  `shellcheck --severity=warning` on the modified script — exit 0; `bash -n` — exit 0;
  `contract-regressions.sh`, `prescribed-shell-canonicalization.sh` and the heavy-tier
  `staged-skills-contract.sh` all exit 0. No debug artifacts: no added line carries `TODO`, `FIXME` or
  `debugger` as a whole word, which matters unusually much here because this request's own subject is
  that check.

## Why
Prose that restates a rule code already enforces is read on every REQ by the three highest-churn action files (272, 93 and 162 commits in twelve months) and can drift from the code. The audit labelled this class INFERRED: the code enforces the rule, and whether reviewer prose is meant as a second independent read is the builder's judgment to confirm before cutting.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 1, sweep_key `qualify-debug-artifact-prose-restated`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -15. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `work.md` — Red Flags row `diff contains \`console.log\`, \`debugger\`, or \`TODO\`` ⇐ `QUALIFY-DEBUG-ARTIFACT`.
- `work.md` — Common Rationalizations row `A checked \`[UNIFY]\` over a diff containing \`console.log\`` ⇐ `QUALIFY-UNIFY-DISARMED`.
- `work-reference.md` — `Read the actual diff for debug artifacts` ⇐ `QUALIFY-PAU-UNCHECKED`. If queued REQ-510 (sweep work-reference sections owned by CLI tests) has already removed this site when this REQ is claimed, skip it; shared files are not a dependency.
- `review-work.md` — `Builder checked all P-A-U boxes but the diff contains` and `Diff hygiene — no debug artifacts — console.log/print lines` ⇐ the same `checks.go` pair.
- Contrast (house style, not an instance): `work.md` says of the blocked probe "this supersedes scripts/run-blocked-check.sh ... prose must not execute the probe a second time"; `capture.md` names the SessionStart hook as owner instead of restating it.
- Reproduce at dc8a64e3: `rg -n -e 'console\.log' -e 'debug artifacts' skills/do-work/actions/work.md skills/do-work/actions/work-reference.md skills/do-work/actions/review-work.md && rg -n 'QUALIFY-DEBUG-ARTIFACT|QUALIFY-PAU-UNCHECKED|QUALIFY-UNIFY-DISARMED' skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- `_dev/tests/contract-regressions.sh` may pin some of these sentences; delete the matching predicates in the same commit rather than keeping a sentence to satisfy a pin.
- Prime `_dev/primes/prime-action-files.md` first.
- Lock-in limit: debug-artifact rule mentions across work.md, review-work.md, work-reference.md: ≤ 3 after this REQ (today 9).

## Dependencies
No dependency. Overlaps queued REQ-510 on `work-reference.md`; overlap is not a dependency, the builder checks the site at claim time.

Verify repair (2026-09-03, `do-work verify-requests` on UR-105): the audit's plan line said "fold the work-reference.md site into REQ-510"; this REQ deliberately owns all five sites instead, because REQ-510 is last in the nine-deep REQ-502 chain and a fold would park one sentence behind eight REQs. The maintainer accepted that shape.

## Builder Guidance
Mixed: firm that one sentence naming the finding codes stays in `work.md` Step 6.3; latitude to keep a reviewer-side sentence if reading the code shows the review pass runs on a diff `qualify` never saw.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine mentions across the three files (work.md 5, review-work.md 3, work-reference.md 1) and the three finding codes in checks.go.
**GREEN when:** At most three mentions remain across the three files; the lock-in pins the mention count at the post-fix value.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for qualify-debug-artifact-prose-restated.*

## Triage

**Route: B** — Explore then build.

**Reasoning:** The work itself is four text deletions, one four-word trim, one added sentence and one
shell assertion, in files the request names — Route A shaped. What made exploration necessary is that
the request's own baseline does not survive contact with HEAD: its site count, the commit its Reproduce
line names, and the claim that the sentence it says to keep already exists are all wrong. Deciding
which mentions are restatements and which two are independent reads is discovery, and it is done.

**Planning:** Skipped.

**The request's own baseline is stale and is corrected here rather than followed.** It claims nine
prose sites across three action files. Several of its claims did not survive contact with HEAD, and one
of them would have sent the builder to a heading that no longer exists. The exploration lists each.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Explore agent, read-only, re-verified against HEAD rather than against the audited commit. Full report
in the run directory as `do-work/runs/work-2026-09-05-231943/REQ-556-exploration.md`.

**Ten of the request's own statements were checked and several do not hold**, including its headline
count, the commit its Reproduce line names, and the claim that the sentence it says to *keep* already
exists. The exploration replaces the count with a site-by-site list and settles which two mentions must
survive: `review-work.md`'s standalone-review hygiene bullet, which is a read the canonical `qualify`
never makes, and the emitted P-A-U template payload, which is byte-identical across four shipped files.
Neither is a restatement, and cutting either to chase a smaller number would remove a real check.

**The rule itself lives in code, not prose.** `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and
`QUALIFY-UNIFY-DISARMED` in `internal/corehelpers/checks.go` are what enforce it, which is what makes
the action-file copies restatements rather than the rule.

**The lock-in counts rather than name-lists**, because a new restatement is the regression whatever
words it uses, and it fails loudly on a missing target file so that a rename cannot retire it silently.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — four edits: one four-word trim and three deletions of prose the canonical qualify already enforces
- `skills/do-work/actions/review-work.md` (modify) — delete the Red Flags bullet that restates the rule; the Step 6 hygiene bullet stays, because it is a read qualify never makes
- `skills/do-work/actions/work-reference.md` (modify) — delete one anti-rationalization table row, leaving the other five and the heading
- `_dev/tests/audit-lockins.sh` (modify) — one counting assertion in the file's existing per-Finding shape, beside REQ-552's and REQ-554's

**Files I will NOT touch:** `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` — the
rule lives there and is not changing. The three other shipped files carrying the byte-identical P-A-U
template payload. Any action file outside the three named.

**Acceptance criteria (restated from REQ):**
- [ ] The debug-artifact and P-A-U-honesty rule is stated once, in code, with the action files pointing at it rather than restating it
- [ ] The two mentions that are not restatements survive
- [ ] A lock-in assertion fails if a restatement returns, counted rather than name-listed
- [ ] A renamed or missing target file fails the assertion loudly instead of counting zero
- [ ] The gate is green

## Pre-Flight

**Git:** ✓ Clean. Canonical `recover` reports `FINALIZATION-NONE`.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at this REQ's claim revision
`d3ceca3` — **88s wall**, exit status read directly from `$?`.

**Tests baseline:** ✓ `bash _dev/tests/audit-lockins.sh` exited 0, launched true. That file already
carries this run's blocks from REQ-552 and REQ-554, so a later red in it is attributable.

**A stale-anchor hazard specific to this request.** Its exploration was taken before REQ-486 landed,
and REQ-486 edited one of the three files this request cuts from. Every edit is located by text rather
than by line number, and the builder reports which anchors had moved.

**Dependencies:** ✓ `depends_on` empty. Toolchain at or above every floor the gate requires.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/work-reference.md`
- `_dev/tests/audit-lockins.sh`

**What was done:** Seven mentions of the debug-artifact and P-A-U-honesty rule across three shipped
action files become two. The rule itself lives in `internal/corehelpers/checks.go` as
`QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and `QUALIFY-UNIFY-DISARMED`, so the prose copies
were restatements; one added sentence points at the owner instead. The two surviving mentions are kept
deliberately — `review-work.md`'s standalone-review hygiene bullet is a read the canonical `qualify`
never makes, and the emitted P-A-U template payload is byte-identical across four shipped files.

**The request claimed nine sites; HEAD has seven**, and the audited commit its Reproduce line names is
not an object in this clone, so its captured red could not be replayed and no lines were manufactured
to reach nine. One anchor had moved: REQ-486 pushed the anti-rationalization row eight lines down
earlier in this run. Every edit was located by text.

**After the review, the guard was rebuilt.** It counted lines rather than matches against a ceiling
with no headroom, so a pure reflow failed the gate; and its comment claimed a reworded restatement
would be caught while it matched two case-sensitive literals. It now counts matches, reads its
scanner's exit status without a pipeline, matches the marker vocabulary the code itself uses, carries a
floor as well as a ceiling, and names path, line and matched text per site. A companion assertion pins
every `QUALIFY-*` code named in the action files against `checks.go` — which closes the drift class
this request exists to remove, rather than moving it one file over.

Merge ranges `c40c6d1b..32ecdba0` (36 insertions, 6 deletions) and `e05c2371..2acc4657` (101
insertions, 15 deletions).

## Decisions — implementation

- **D-01 — the ceiling is 2, not the exploration's fallback 3. DECIDE & STATE.** Reaching 3 meant
  keeping a four-word phrase in `work.md` whose own preamble calls that block pointers, while the real
  instruction ships in the template payload inside every REQ file. After the review the pin became
  exact — a floor as well as a ceiling — because losing one of the two protected mentions was silently
  green.
- **D-02 — the added sentence names the finding codes rather than a count. DECIDE & STATE.** The
  exploration's draft said "a separate pair"; there are three, plus a fourth beside the main code.
  Shipping a count would have put a fresh wrong number into the prose this request exists to
  de-duplicate. After the review the enumeration is pinned against `checks.go` and introduced as
  illustrative, so it can go stale loudly rather than quietly.
- **D-03 — the sentence is an addition, not a retention, and it is earned. DECIDE & STATE.** The
  request says to keep a sentence naming the codes; no such sentence existed. `work.md` invokes the CLI
  with `--format json`, and before this change no shipped file routed a `QUALIFY-*` token to its owner.
  `maintenance.md` asks for a concrete case that fails without an addition, and that is it.
- **D-04 — widen the pattern set rather than narrow the claim. DECIDE & STATE, measured.** Faced with a
  comment that overclaimed, the choice was to make the code match the comment or the comment match the
  code. Widening to the marker vocabulary `checks.go` uses measured **zero** false positives against
  the three files as they stand, so the honest option was also the stronger one. What remains uncaught
  is named in the comment instead of being papered over.

## Qualification

**Passed.** Read from the remediation range `e05c2371..2acc465a`; canonical `qualify` and `scope-drift`
both satisfied. The original range is `c40c6d1b..32ecdba0`.

- **One warning is judged rather than obeyed.** `QUALIFY-PATH-NOT-IN-DIFF` names
  `work-reference.md` and `review-work.md`, which the Implementation Summary claims and the remediation
  range does not contain. They belong to the original range. The summary lists the union deliberately,
  because that is what this request built.
- **The request's own baseline was corrected rather than followed.** It claimed nine sites; seven is
  what HEAD has, and the audited commit its Reproduce line names is not an object in this clone, so no
  lines were manufactured to reach nine. One anchor had moved eight lines under REQ-486 earlier in this
  run, and every edit was located by text rather than by line number.
- **Nothing load-bearing was cut, and three reviewers each traced it independently.** Every deleted
  sentence has a live successor: the builder-side instruction ships verbatim in the P-A-U template
  payload inside every REQ file, and the orchestrator-side reads are the `QUALIFY-*` findings in
  `checks.go`. The two mentions that are not restatements survive.
- **The guard was wrong in both directions and both were shown, not argued.** It counted lines against
  a zero-headroom ceiling, so splitting one bullet across two lines without changing a word failed the
  gate with a message claiming a restatement had returned — reproduced by three reviewers with three
  different reflows. And its comment claimed a reworded restatement would be caught while it matched
  two case-sensitive literals; three rewordings each left it green.
- **The fix widened the code to match the claim rather than shrinking the claim, and that was a
  measurement, not a preference.** The marker vocabulary `checks.go` actually uses matches **zero**
  times across the three files as they stand, so widening added no false positives. What is still
  uncaught is now named in the comment instead of being papered over.
- **The remediation also closed the drift class this request is about, rather than moving it.** A
  companion assertion pins every `QUALIFY-*` code named in the three action files against `checks.go`,
  and it was proven from both sides: renaming the code in the prose fails, and renaming it in
  `checks.go` fails.

Requirements traced: the rule is stated once, in code, with the action files pointing at it; the two
non-restatements survive; the assertion fails when a restatement returns, in any of the three files and
in either vocabulary; a renamed or unreadable target file fails loudly instead of counting zero; and
the gate is green.

*Checked by work action*

## Testing

**Tests run:** `audit-lockins.sh`, `action-shell-blocks.sh`, `contract-regressions.sh`,
`prescribed-shell-canonicalization.sh`, the heavy-tier `staged-skills-contract.sh`, plus `shellcheck
--severity=warning` and `bash -n` on the modified script.

**Result:** ✓ Green, each with its own line: `Audit lock-in regressions passed.`;
`Shell-block lint passed: 74 fenced blocks and 33 shipped shell files; ShellCheck enabled.`;
`Contract regression checks passed.`; `Prescribed shell primitive canonicalization checks passed.`;
`Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).`
and `staged skills contract: PASS`. The canonical gate exited 0 at this REQ's claim revision
`d3ceca3` — 88s wall.

**Eleven ablations, each reverted alone and restored with its sha256 compared.** Three restatements
pasted back, one per file. Three rewordings, one per vocabulary. One reflow. One deletion of a
protected mention. One reworded protected mention. One broken scanner regex. One moved target file.
Two code-name renames, one in the prose and one in `checks.go`. Every one red; every one green after
restore.

*Verified by work action*

## Review

**Overall: 63%** | 2026-09-06 02:56 UTC

| Dimension | Score |
|-----------|-------|
| Requirements | 70% |
| Code Quality | 65% |
| Test Adequacy | 55% |
| Scope | 100% |
| Risk | Medium |
| Acceptance | Partial |

**Approve with follow-ups** — the prose cuts are correct and nothing load-bearing lost a home, but the new lock-in that is supposed to protect them fires on a formatting change and stays silent on a reworded restatement.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Zero headroom plus line-counting makes the lock-in fail on a pure reflow: `review-work.md:106` carries both matched strings on one physical line, `grep -c` counts lines, so splitting that bullet after "no debug artifacts —" with no word changed takes the count 2 → 3 and prints "FAIL: 3 debug-artifact rule mentions ... ceiling is 2". I reproduced this in a scratch copy (baseline had no such FAIL, reflow produced it); reviewers 2 and 3 reproduced it independently by different reflows. The message claims a restatement returned, which is false, and the fix it implies (delete a mention) is the wrong action. `audit-lockins.sh` runs in the fast probe lane (`_dev/tests/contracts/probe-lanes.sh:29`) over the three highest-churn action files, so every future editor of those files can trip it. Fix: count matches, not lines (`grep -o ... | wc -l`, ceiling 3) — `_dev/tests/audit-lockins.sh:290,304` — impact-rule-change → report only
- The lock-in's comment claims "a new restatement is the regression whatever words it uses"; it greps two case-sensitive literals, so most rewordings pass green. I appended `- Builder checked all P-A-U boxes but the diff still adds \`debugger\`, \`TODO\` or \`FIXME\` markers.` to `work.md` in a scratch copy and the suite reported zero debug-artifact failures. Reviewers 2 and 3 each showed further green cases: the deleted review bullet with "console.log," dropped, capitalised "Debug artifacts", singular "debug artifact", and a leftover-print-statements rationalization row. The vocabulary that passes is exactly the vocabulary `checks.go:24` matches (`\b(debugger|TODO|FIXME)\b`). Acceptance criterion 3 ("counted rather than name-listed") is met in form, not substance — counting fixes the list of sites, not the list of spellings. Fix the comment and the criterion claim to say it pins two literal spellings, or widen the pattern set to the marker vocabulary `checks.go` uses — `_dev/tests/audit-lockins.sh:279-311` — impact-rule-change → report only

**Minor findings:**
- The added sentence names 7 of the 15 `QUALIFY-*` codes `checks.go` emits, with nothing pinning the list to the code, which is the drift class this REQ removed — the unnamed eight are `QUALIFY-SUMMARY-MISSING`, `-DIFF-RANGE-INVALID`, `-SUMMARY-MALFORMED`, `-DELETED-PATH-PRESENT`, `-PATH-NOT-IN-DIFF`, `-CLAIMED-PATH-MISSING`, `-NEW-FILE-UNWIRED`, `-NO-PROJECT-FILES` (`work.md:335`) — impact-rule-change → report only
- `work.md:335` says `--format json` hands you "raw codes"; it does not — `resultmodel.CommandFinding` ships `automation_stop_reason` (`result_model.go:54`) with a plain-English meaning per code, and `renderText` prints the same string. Reviewer 1 says this also weakens the hand-back's replay case; reviewer 3 says the case still holds. I side with reviewer 3 on keeping the sentence (the stop reason names the meaning, not the owner, and no shipped prose routed the token to `advance`) and with reviewer 1 on the wording: cut or reword "as raw codes" to "as finding codes with a stop reason" — impact-negligible → report only
- The ratchet is one-sided — only `-gt ceiling` is checked, so rewording or deleting `review-work.md:106` drops the count to 1 and stays green, losing a read `qualify` never makes; the block's comment warns against exactly this and enforces nothing (`_dev/tests/audit-lockins.sh:307`) — impact-negligible → report only
- The failure message reports one aggregate count and a hardcoded three-name string, naming no file or line, unlike the sibling REQ-552 and REQ-554 blocks in the same script which print `path:line` per offending site (`audit-lockins.sh:308` vs `:249-251`, `:271-274`) — impact-negligible → report only
- The hand-back's deferred-release note says VERSION is "at 0.303.10 on this branch"; it is 0.304.0 at the base, the branch head, and HEAD, so the follow-on release is 0.304.1 in `skills/do-work/VERSION`, root `VERSION`, `skills/do-work/actions/version.md:5`, plus a new top entry in root `CHANGELOG.md` copied byte-identically to `skills/do-work/CHANGELOG.md` — impact-negligible → report only
- The missing-file guard covers absence but not unreadability: `debug_rule_file_hits=$(grep -c ...)` discards grep's status, so exit 2 yields an empty string that the arithmetic folds to 0 and the ceiling stays green. Reviewers 1 and 3 disagree on weight — reviewer 3 calls the `[ ! -f ]` guard sufficient for the case that matters. I keep it as Minor: reachability is low, but both sibling blocks capture `$?` and fail on exit > 1, so this is an unexplained departure in the same file (`audit-lockins.sh:304` vs `:263-268`) — impact-negligible → report only
- `debugger`, `TODO` and `FIXME` now appear in none of the three action files. Orchestrated mode is covered (`QUALIFY-DEBUG-ARTIFACT` is an Error), standalone review is not — `review-work.md:106` names only console.log/print and temporary files. Reviewer 1 rated this a Nit because the surviving bullet states the condition and the house rule treats trailing examples as illustrative; I keep it Minor because standalone review is the mode used to justify keeping that bullet — impact-negligible → report only
- "`advance`'s qualification gate owns the ... P-A-U-honesty mechanics" overstates `checks.go`: the three P-A-U codes fire independently and nothing correlates a checked box with a dirty diff, which is what all four deleted sentences asserted. The dirty diff is still caught (Error, regardless of box state), so the damage is bounded to the lost inference — impact-negligible → report only

**Nit findings:**
- "An unfinished-work marker the diff added" mislabels `debugger`, which is a debug statement, not an unfinished-work marker; one word fixes it (`work.md:335`) — impact-negligible → report only
- The four-word trim leaves `work.md:292` paraphrasing three of the template's four UNIFY clauses; if the bullet is a pointer, point at the template instead of paraphrasing most of it — impact-negligible → report only
- `QUALIFY-NEW-FILE-UNWIRED` is the fourth code that belongs by the sentence's own logic — its judgment prose sits one sentence earlier at `work.md:335` without the token attached — impact-negligible → report only

**Requirements checklist:**
- [x] The debug-artifact and P-A-U-honesty rule is stated once in code, with the action files pointing at it — delivered; four restatements removed, `work.md:335` names the codes and ends "do not restate the rule they enforce"
- [x] Nothing load-bearing was cut — delivered; all three reviewers traced every deleted sentence to a live successor (the builder-side instruction ships verbatim in the P-A-U template payload inside every REQ file, `work.md:292` routes the builder there; the orchestrator-side reads are the `QUALIFY-*` findings in `checks.go`)
- [x] The two mentions that are not restatements survive — delivered; `review-work.md:106` (standalone review reads a diff `qualify` never sees) and `:374` (emitted template payload)
- [ ] A lock-in assertion fails if a restatement returns, counted rather than name-listed — not delivered in substance; it catches a pasted-back copy but not a reworded one, and it false-fires on a reflow (Important findings above). Reviewer 1 marked this met with caveats, reviewers 2 and 3 marked it not met; I follow the two who showed the failing runs
- [x] A renamed or missing target file fails loudly instead of counting zero — delivered; reviewer 2 renamed each file singly and all three together and got a per-file FAIL naming the correct path
- [x] Nothing changed outside the four declared files — delivered; `git diff --stat c40c6d1..32ecdba` is 4 files, 36 insertions, 6 deletions, all in the REQ's `write_set`
- [x] All three edited action files still satisfy `prime-action-files.md` — delivered; both Red Flags sections keep earned content, the Common Rationalizations table keeps do-work-specific nouns, and the contract-pinned heading at `work-reference.md:606` is untouched
- [ ] The release is correctly identified with the right number — partially; deferring to finalization is right, the named base version is stale (Minor above)

**Acceptance:** Partial — `bash _dev/tests/audit-lockins.sh` at HEAD prints "Audit lock-in regressions passed." exit 0 and the working tree is clean, but the new assertion fails its two advertised properties on demonstration: green on a reworded restatement, red on a whitespace-only reflow.

**Suggested testing:** 5 items
- Re-run `bash _dev/tests/audit-lockins.sh` after switching `grep -c` to match-counting with ceiling 3, and confirm the reflow of `review-work.md:106` no longer trips it
- Paste the deleted Red Flags bullet back into each of the three files in turn and confirm each goes red independently, then repeat with the three reworded variants (dropped "console.log,", capitalised "Debug artifacts", singular "debug artifact")
- Delete `review-work.md:106` outright and confirm whether the intended behavior is silence or a failure; the block's comment says the mention is protected but nothing enforces a floor
- Run `shellcheck --severity=warning` and `bash -n` on `_dev/tests/audit-lockins.sh` after any change to the block (currently exit 0; `-S style` reports only 5 pre-existing findings at lines 28/48/118/229/230)
- At finalization, confirm all four version mirrors move together to 0.304.1 — omitting one refuses with RELEASE-MIRROR-UNDECLARED, and `skills/do-work-board/tools/queue-kanban/VERSION` is independently versioned and must not be dragged in

**Follow-ups created:** None (13 findings report only)

*Reviewed by review-work action*

## Lessons Learned

- **A ratchet with no headroom fails on formatting, and `grep -c` counts lines.** Two independently
  reasonable choices — pin the ceiling at the current measurement, count with `grep -c` — combine into
  a gate that fails when someone reflows a paragraph. The lesson is not "leave headroom": headroom
  absorbs real regressions too. It is that the unit you count must be the unit the thing you are
  guarding is made of. Restatements are matches, not lines.
- **A comment that claims more than the code does is worse than no comment.** This one said "a new
  restatement is the regression whatever words it uses" over a grep for two case-sensitive literals.
  The next maintainer reads the claim, believes the class is guarded, and stops looking. When the code
  and the comment disagree, either widen the code or narrow the claim — and prefer widening only when
  you have measured the false-positive rate, which here was zero.
- **A one-sided ratchet guards half the property.** The count could only go up; deleting one of the two
  mentions the block's own comment says must survive was silently green. If a number is protected in
  one direction, ask out loud why the other direction is not a regression.
- **Deleting prose is only safe once each sentence's successor is named.** Three reviewers each traced
  every deleted sentence to a live successor before accepting the cut. The one that could have gone
  wrong is the builder-side instruction: it survives because it ships verbatim in the template payload
  inside every REQ file, not because the finding codes cover it — those are an orchestrator-side read.
- **A request's stated baseline is a claim about the world, and this one had three wrong parts.** Nine
  sites where there were seven; an audited commit not present in the clone; and a sentence it said to
  keep that did not exist, making the "retention" an addition that had to earn itself separately.

## Orientation

The rule lives in `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` as
`QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and `QUALIFY-UNIFY-DISARMED`. Three shipped action
files now point at it rather than restating it, and one sentence in `work.md` Step 6.3 names the codes.

The guard is in `_dev/tests/audit-lockins.sh`, beside REQ-552's Finding 9 and REQ-554's Finding 6. It
is an exact pin, not a ceiling: a restatement returning fails it, and losing one of the two protected
mentions fails it too. The two protected mentions are `review-work.md`'s standalone-review hygiene
bullet — a read the canonical `qualify` never makes, because a standalone review sees a diff `qualify`
does not — and the emitted P-A-U template payload, which must stay byte-identical across four shipped
files. A companion assertion in the same block pins every `QUALIFY-*` code named in the action files
against `checks.go`, from both directions.

Two things are recorded and unfixed. `review-work.md`'s surviving bullet names `console.log` and print
statements but not `debugger`, `TODO` or `FIXME`, so standalone review covers less than orchestrated
review does — changing it would move the pin, so it is a follow-up. And the sentence saying the
qualification gate owns the P-A-U-honesty mechanics slightly overstates `checks.go`, where the three
codes fire independently and nothing correlates a checked box with a dirty diff; the dirty diff is
still caught regardless of box state, so what is lost is the inference, not the check.
