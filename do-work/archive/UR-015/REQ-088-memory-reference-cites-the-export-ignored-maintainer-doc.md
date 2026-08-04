---
id: REQ-088
title: "Confirm: fix memory-reference.md's citation of the export-ignored CLAUDE.md"
status: completed
status_changed_at: 2026-08-03T22:49:47Z
claimed_at: 2026-08-04T04:21:50Z
completed_at: 2026-08-04T04:28:03Z
commit: bb8cf3b
created_at: 2026-08-03T21:54:55Z
kb_status: pending
user_request: UR-015
domain: general
addendum_to: REQ-078
discovered_during: REQ-078
maintenance: true
---

# Confirm: fix memory-reference.md's citation of the export-ignored CLAUDE.md

## What

While working REQ-078 the builder noticed, in a file it was already editing, that
`actions/memory-reference.md:88` cites this repo's `CLAUDE.md` inline:

> Steps 3–4 are scoring/formatting the agent performs on the grep output — they need no additional
> shell state, so nothing carries between command blocks (CLAUDE.md: shell state does not survive
> between prescribed blocks).

`CLAUDE.md` is the maintainer document and is `export-ignore`d from the distribution tarball, so in
every consumer install that citation points at a file that is not there. The skill's own rule for
shipped files is to restate the rule inline or point at a shipped home instead. The suite greps for
the common citation idioms; this phrasing is not one of them, which is why it survived.

It was **not** fixed inline: it is unrelated to REQ-078's premise (the timestamp command's single
home), and the implementation-time rule is to record an out-of-scope find rather than sweep it in.

## Why this is your call, not the builder's

The fix is trivial, but *which* fix is a small editorial choice — delete the parenthetical, restate
the rule in six words, or point at a shipped home — and it is a one-line change to a shipped file
that nothing is currently broken by. A dangling citation costs a consumer one confused lookup; it
corrupts nothing. So it is worth queuing and not worth a builder deciding unilaterally while working
on something else.

## What Would Change

Whichever option you pick, one line of `actions/memory-reference.md` changes and nothing else does.
If you also want the suite's citation grep widened so the next phrasing is caught mechanically, say
so — that is a second, larger change and would touch `_dev/tests/contract-regressions.sh`.

## Open Questions

- [x] How should `actions/memory-reference.md:88`'s dangling `CLAUDE.md` citation be fixed? → Confirmed: restate the rule inline
  Recommended: **restate the rule inline** — replace "(CLAUDE.md: shell state does not survive
  between prescribed blocks)" with "(shell state does not survive between prescribed command
  blocks)". The sentence keeps its point and needs no external file.
  Value: the consumer reading it gets the whole rule where they are standing; nothing to look up.
  Risk: none material — it is a parenthetical in one sentence, and reverting is a one-line edit.
  Also: **(a)** delete the parenthetical entirely, since the sentence already explains itself
  ("they need no additional shell state, so nothing carries between command blocks"); **(b)** point
  at a shipped home instead of restating — but there is no single shipped file that owns this rule
  today, so this option means first choosing one; **(c)** leave it, and instead widen the suite's
  citation grep so future occurrences are caught — this does not fix the existing line.

## Answer

- **[2026-08-03] Decided by:** user, via `do-work clarify`
- **Decision:** restate the rule inline — the builder's recommendation. Replace
  `(CLAUDE.md: shell state does not survive between prescribed blocks)` with
  `(shell state does not survive between prescribed command blocks)` at
  `actions/memory-reference.md:88`. Nothing else in the file changes.
- **Explicitly out of scope:** widening `_dev/tests/contract-regressions.sh`'s citation grep
  (offered as option (c)) was not chosen and is not part of this REQ. Choosing a shipped
  owner-file for the shell-state rule (option (b)) is likewise out of scope.
- **Why this stayed queued rather than closing:** confirming the recommendation authorizes the
  edit, it does not stand in for it — no implementation exists yet, so this flips to `pending`
  for the next `do-work run` instead of taking the builder-was-right fast path.

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact file and line, and the `## Answer` section already
authorizes the exact replacement text. One line of one file changes; there is no "where" or
"how" left to discover.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `actions/memory-reference.md` (modified)

**What was done:** Replaced the parenthetical `(CLAUDE.md: shell state does not survive between
prescribed blocks)` at line 88 with `(shell state does not survive between prescribed command
blocks)`, exactly as authorized in `## Answer`. `git diff --numstat` reads `1 1` — one line
changed, nothing else. The file now contains no reference to `CLAUDE.md`.

## Qualification

**Passed** — 1 file verified, 1 requirement traced, single-line diff confirmed by
`git diff --numstat` (`1	1`).

- **Files exist / show in diff:** `actions/memory-reference.md` appears in `git diff` with the
  expected single-line hunk at line 88.
- **Change is substantive:** the authorized edit landed verbatim; the removed idiom is gone
  (`grep -n 'CLAUDE\.md' actions/memory-reference.md` returns no hits).
- **Requirement traced:** the `## Answer` decision named one replacement in one file. That is the
  whole delivered change. The two explicitly out-of-scope items (widening the suite's citation
  grep; choosing a shipped owner file for the shell-state rule) were **not** touched — confirmed
  by the diff containing only `actions/memory-reference.md`.
- **No scope drift:** Route A skips Step 5.5, so there is no `## Scope` section to compare
  against and `write_set` stays as captured. The diff's single file is its own proof.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Passing (exit 0, zero FAIL lines)

**Baseline comparison (established for this REQ, not inherited):**
- Pre-change: edit stashed with `git stash push -- actions/memory-reference.md` → suite exit **0**
- Post-change: edit restored → suite exit **0**

No regression. A stale `do-work/working/baseline.json` from a prior session was deliberately
**not** used: Route A skips Step 5.75, so that file describes another REQ's run, and its
`"exit_status": 0` would have been an inherited claim rather than this REQ's measurement.

**Note on the recorded baseline:** the previous session's checkpoint recorded an "8 FAIL" suite
baseline as this repo's pre-existing state. That state does **not** reproduce here — the suite
exits 0 with zero FAIL lines. That observation is evidence bearing on **REQ-090**, the live
`pending-answers` REQ whose entire premise is those seven failures; it is reported into REQ-090 via
`do-work clarify` rather than duplicated as a new REQ from here. Details in `## Lessons Learned`.

*Verified by work action*

## Decisions

- **D-01**: Kept the user's authorized replacement text verbatim even though
  `crew-members/maintenance.md` § 1 ("delete before you add") points at the declined option (a),
  deleting the parenthetical outright. Reasoning: the `## Answer` section is an explicit user
  decision made from a menu that *included* option (a); a builder re-litigating it would convert a
  settled choice back into an open one. The redundancy this creates is reported as a review
  finding for the user to judge rather than silently corrected. **DECIDE & STATE** — the edit is
  one line and trivially reversible either way.

## Discovered Tasks

- **[normal]** Six shipped sites across two files in `tools/queue-kanban/` cite this repo's
  export-ignored `CLAUDE.md` as the source of the rule they implement — `verify.go:123`, `verify.go:156`,
  `verify.go:186`, `verify_test.go:89`, `verify_test.go:141`, `verify_test.go:171`, all of the
  form `CLAUDE.md § Before Every Commit`. `tools/queue-kanban/` is **not** export-ignored, so it
  ships to every consumer while `CLAUDE.md` does not — the same dangling-citation defect this REQ
  fixed in `actions/memory-reference.md`, six more times. A seventh, lower-confidence site:
  `prompts/prompt-kit-step6-constraint-architecture.md:78` attributes a rule to a "CLAUDE.md
  standard" that is not in this repo's `CLAUDE.md` at all. Additionally,
  `_dev/tests/contract-regressions.sh`'s `self_citation_pattern` matches **zero** of these — and
  matched zero of the defect this REQ fixed. Out of scope here: the `## Answer` explicitly
  excluded both widening that grep and touching anything beyond the one line.

- **[low]** `do-work/working/baseline.json` is a stale pre-flight artifact left behind by a
  previous session (mtime 2026-08-04T00:48, `"exit_status": 0`). `do-work/working/` is supposed to
  be empty at a clean queue boundary, and a leftover baseline is exactly the kind of inherited
  claim Step 6.5 warns about — a later REQ could compare its failures against another REQ's
  record. It is untracked, so no commit carries it. Sweeping it belongs to `do-work cleanup`, not
  to this REQ.

## Review

**Overall: 96%** | 2026-08-04T04:21:50Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — the authorized line is in place, the file no longer cites `CLAUDE.md`, and
the suite is green before and after.
**Suggested testing:** 0 items
**Follow-ups created:** REQ-093

**Restatement Sweep (run — mandatory):** The change inlines a rule whose canonical home is
export-ignored, so the sweep asked whether any *shipped* file already owns that rule (which would
make the inline copy a second, drift-prone statement). It does not. The rule is already restated
inline at **nine** shipped sites — `actions/work.md:394`, `actions/work-reference.md:302`,
`actions/work-reference.md:421`, `actions/work-reference.md:775`, `actions/install.md:394`,
`actions/board.md:80`, `actions/ai-report.md:256`, `tools/checks/qualify.sh:54`, and now
`actions/memory-reference.md:88` — and **none** of the other eight cites `CLAUDE.md`. So
inline-restatement is the established house pattern, this change makes the tenth site consistent
with the other nine, and option (b)'s premise ("there is no single shipped file that owns this
rule today") is confirmed correct rather than assumed. No stale restatement anywhere depended on
the removed citation.

**Minor finding — the authorized text is redundant with the clause it follows.** The line now
reads: *"…they need no additional shell state, so nothing carries between command blocks (shell
state does not survive between prescribed command blocks)."* The parenthetical restates the clause
immediately before it almost word for word. The other nine sites each state the rule **once**;
this is the only one that says it twice. Declined option (a) — delete the parenthetical, since the
sentence already carries the point — would have produced the cleaner line, and the option menu did
not surface that the surrounding clause already said it. Left as authorized per D-01; a one-line
follow-up if the user wants it.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Running three grep shapes before trusting the REQ's one-site inventory. The REQ
named one site; the shipped-paths sweep found six more of the same defect class in
`tools/queue-kanban/*.go`, plus proof that the suite's own guard pattern matches none of them —
including the defect this REQ was filed about. Running the suite's `self_citation_pattern`
verbatim and getting **0 hits** turned "the guard is probably too narrow" into a measured fact.

**What didn't:** Trusting the inherited test baseline. The checkpoint recorded "8 FAIL" as this
repo's pre-existing suite state, confirmed twice last session by stash-and-compare. It does not
reproduce: the suite exits 0 with zero FAIL lines, the seven update-script probes pass, and
`tools/do-work-update.sh` demonstrably contains all five strings REQ-090 says are absent
(lines 166, 194, 202, 204, 218). Checked for the obvious environmental causes — lingering
worktrees (none), cwd sensitivity (green from a subdirectory), and a vacuous skip path (the probe
has exactly one, `git` unavailable, and `git` is present). The cause of last session's observation
is still unexplained, which is itself the finding.

**Worth knowing:** A stale `baseline.json` in `do-work/working/` outlives the session that wrote
it and is silently available to the next REQ's Step 6.5 comparison. Route A skips pre-flight, so a
Route A REQ that reaches for a baseline is always reading someone else's. Measure your own.

## Orientation

The `do-work memory` recall procedure's shell-state note now states its rule inline instead of
citing the maintainer doc, so a consumer reading `actions/memory-reference.md` no longer follows a
pointer to a file their install does not contain. Leaf change — one parenthetical in one action's
companion reference; no contract, data flow, or concept renamed. No prime files are listed on this
REQ, and no prime was made stale.
