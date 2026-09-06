---
id: REQ-596
status: completed
domain: general
created_at: 2026-09-06T05:18:15Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: B
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T05:21:33Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
maintenance: true
depends_on: [REQ-595]
related: [REQ-555, REQ-595]
write_set: [skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Correct three more stale mechanism claims in the prescribed-shell guide, in sections REQ-595 never opened'
claimed_at: 2026-09-06T05:20:24Z
completed_at: 2026-09-06T06:58:37Z
commit: 0bbd10d299efef80e78cbe2b1d7b61f00e55fbb2
release_at: 2026-09-06T06:58:37Z
---

# Correct Three More Stale Mechanism Claims in the Prescribed-Shell Guide

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` and ran a 69-claim sweep of the guide
  against the Go implementations before touching the file. Approach: correct every claim the sweep
  proved false inside the two sections the request's three sit in, using fixtures rather than reading,
  and capture the ten in other sections as their own request.
- [x] **[APPLY]:** One file changed, `skills/do-work/docs/prescribed-shell-primitives.md`: eleven
  corrections in two sections. No Go source, no route column, no Mechanics column, no heading, nothing
  under `_dev/tests/`.
- [x] **[UNIFY]:** `git diff --numstat` — the guide at +8/-7, one file. Each replacement was checked
  against the fixture output the sweep produced and against the code it cites. All three guards over
  this file pass: `audit-lockins.sh`, `prescribed-shell-canonicalization.sh` and the new
  `quiet-grep-pipeline-audit.sh` each exit 0. No debug artifacts, no code change.

## What

Three sentences in `skills/do-work/docs/prescribed-shell-primitives.md` describe behaviour the Go code
does not have. All three are in sections REQ-595 did not open, so that request left them alone rather
than widening a second time. Each was verified against the code by an independent reviewer.

## Why

Same finding class as REQ-555 and REQ-595: the guide is the pointer target from sixteen shipped files
and currently describes implementations that do not exist. One of the three loosens a secret-quarantine
instruction, which is the only one with a safety edge.

## Context

Found during the independent three-lens review of REQ-595. That review scored 86% and confirmed all
three against the code; they are outside REQ-595's declared class, so they are captured here.

## Detailed Requirements

- **Line 124, atomic download.** The guide says the download verifies what it wrote and removes its own
  nested artifact. The code pre-checks instead, and a rename cannot nest. Describe what
  `atomic-download` really does.
- **Line 72, inventory ordering.** "An archived REQ outranks an in-flight one" is not what the code
  does: the active flag is stored at `internal/corehelpers/inventory.go:307` and never read, and the
  only comparison is on the completion instant.
- **Line 57, secret quarantine.** The guide's pattern is `*credentials*`, narrower than the code's
  `strings.Contains(base, "credential")`. A reader following the guide's by-hand fallback would miss a
  file the tool quarantines. Correct the pattern to match the code.
- Check the sections these three sit in for the same class in the same pass, the way REQ-595 checked all
  fourteen Mechanics cells. The sentence that was checked says nothing about its neighbours.

## Constraints

- Guide prose only: no behaviour change, no code change, no route-column change, no heading change.
- `_dev/tests/audit-lockins.sh` Finding 7 pins the route column, the orchestration claim and the
  Mechanics column's shell vocabulary. Do not weaken any of the three.
- If a claim cannot be verified against code, say so in the request rather than guessing.

## Dependencies

Depends on REQ-595, which rewrote three Mechanics cells and two prose claims in the same file.

## Open Questions

None.

## Triage

**Route: B** — Explore then build.

**Reasoning:** The three named claims are already settled — each was verified against the code by the
REQ-595 reviewer and re-verified here before the route was chosen. What makes exploration real is the
request's fourth requirement: check the sections these three sit in for the same class in the same pass.
That is the requirement that turned REQ-595 from a one-cell edit into a three-cell one, and the same
argument applies here — the sentence that was reported says nothing about its neighbours.

**Planning:** Skipped. One file, and the work is whatever the sweep finds.

**One of the three has a safety edge and is the reason this is not "later".** The guide's by-hand
secret-quarantine fallback lists `*credentials*`; the code matches `credential`. A reader following the
guide when the tool is unavailable would not quarantine `credential.json`. The fallback is what someone
executes by hand, so a pattern narrower than the code's is a miss, not a nit.

## Plan

**Planning not required** — Route B: one file, and the edits are whatever the section sweep confirms.

*Skipped by work action*

## Exploration

Three read-only agents checked **69 factual claims** across the guide against the Go implementations,
each running fixtures rather than reading. Full report in the run directory as
`REQ-596-section-sweep.json`. **21 of the 69 do not hold.**

**The three claims the request names are all confirmed, and two of them are worse than reported.**

- Line 57's `*credentials*` is narrower than the code's `strings.Contains(base, "credential")`. Fixture:
  `credential.txt` and `MyCredential.YAML` both come back as excluded rows, and neither matches the
  guide's glob. Every other pattern in that sentence matches the code exactly.
- Line 72's "An archived REQ outranks an in-flight one" is false in the opposite direction as well. The
  `active` flag is computed and stored and never read; the winner is `!exists || completed.After(...)`,
  and because `do-work/working` is walked before `do-work/archive` with a strict `After`, an equal or
  missing timestamp leaves the **in-flight** claim in place.
- Line 124 generalizes three clauses over two commands that behave differently. "Verifies the path it
  actually wrote" is true of the screenshot install, which byte-compares and then confirms the published
  inode, and false of the download, which stats only to report a size and discards that error. This is
  the same shape REQ-555 shipped and its review caught: a sentence true of some of the things it covers.

**The same pass found eight more in the two sections the three sit in**, none of them reported before:

- The `XD` tag is documented as "deleted secret-shaped path" but is also assigned to a deletion that is
  only secret-*derived* — a rename whose destination is ordinary and whose origin is secret-shaped.
  Fixture: `git mv api-secret.pem plain-name.txt` then delete gives `XD` rows for both.
- "Read each tag class this way, and no other way" then gives rules for four of the five tags. `X` is
  the one left without a rule, and `X` is produced for ordinary additions promoted by the ambiguity rule
  as well as for real secrets — so a reader applies the closest listed rule and reads it as a new file,
  which is exactly what `actions/commit.md` forbids. Fixture: an ordinary `brandnew.txt` came back `X`.
- The by-hand inventory fallback is wrong in three ways: it says to use the *complete* porcelain output
  where the tool drops untracked hidden files under `do-work/`; it says to classify each rename origin as
  a path of its own where the tool gives the origin a row only when it is secret-shaped; and it says to
  *add* X paths to the quarantine and overlay them, where `start` *replaces* the file and performs no
  overlay at all.
- The by-hand association fallback filters in-flight REQs out. The tool applies the terminal-success
  status test only to non-working REQs; a `working/` REQ counts whatever its status says. Following the
  guide, every `claimed` or `blocked` REQ is skipped and its files come back unassociated — the exact
  failure the bullet above it exists to prevent.
- Its glob is wrong twice: without `shopt -s globstar` bash expands `**` as a single `*`, so
  `do-work/archive/**/REQ-*.md` requires exactly one intervening directory and misses every REQ sitting
  directly in `do-work/archive/`, and `do-work/working/REQ-*.md` never descends.

**Ten more, in sections these three do not sit in**, are left for a follow-up rather than widening a
third time: claims about lifecycle timing, the merge-aware commit diff, the commit file listing, the
verified-exact-publication rule and the portfolio summary.

*Generated by three Explore agents*

## Scope

**Files I will touch:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — eleven corrections in the two sections the request's three claims sit in

**Eleven, not three, and the extra eight are inside the request's own fourth requirement.** It says to
check the sections these three sit in for the same class in the same pass, because the sentence that was
reported says nothing about its neighbours. The pass found eight more in those two sections. Fixing the
three and leaving the eight would satisfy the request's list and defeat its reason.

**Ten more found outside those two sections are NOT fixed here.** They are in sections this request
never opened, they are the same class, and they go to their own request. That line is drawn where the
request drew it: "the sections these three sit in".

**Files I will NOT touch:** any Go source — prose only, no behaviour change, as the request requires.
The route column and the Mechanics column, which REQ-555 and REQ-595 settled and which
`_dev/tests/audit-lockins.sh` pins in three ways. `_dev/tests/` at all: the corrections here are prose,
and the one guard that could pin a prose claim already exists.

**Acceptance criteria:**
- [ ] The three claims the request names are corrected, each against the code that contradicts them
- [ ] The eight more in the same two sections are corrected in the same pass
- [ ] The by-hand fallbacks — the parts a reader executes when the tool is unavailable — match what the
  tool does, for the secret patterns, the metadata drop, the rename origin, the quarantine write, the
  status test and the file glob
- [ ] No sentence generalizes over two commands that behave differently
- [ ] `_dev/tests/audit-lockins.sh` and `_dev/tests/prescribed-shell-canonicalization.sh` both still exit 0
- [ ] The ten findings outside these two sections are captured rather than fixed

## Pre-Flight

**Green gate at `06be3df`.** `bash _dev/tests/maintainer-verify.sh` printed
`Maintainer verification passed.` and exited 0, gate wall 88s.

**Three guards over this file are green and must stay green.** `_dev/tests/audit-lockins.sh` Finding 7
now pins three things: the route column, the orchestration claim, and the Mechanics column's shell
vocabulary. `_dev/tests/prescribed-shell-canonicalization.sh` pins twelve headings in this document plus
sixteen pointer sites. None of the three columns and no heading changes here.

**This request writes prose only, so no existing lane reads the sentences that change.** The correctness
evidence is the sweep's fixtures — each of the 21 findings was produced by running a fixture through the
real command, not by reading the code — re-checked by the orchestrator for the three the request names.
Saying "no lane reads this" is a statement about the lanes, not about what is possible: REQ-595 added a
guard over the Mechanics column for exactly that reason. Nothing here is guardable the same way, because
these are by-hand procedures whose correctness is a comparison between prose and behaviour.

## Implementation Summary

**Files changed:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**What was done:** Eleven corrections, all inside the two sections the request's three claims sit in.
Ten in "Protected inventory fallbacks" and one in the atomic download section.

**The three the request names.** The secret pattern is `*credential*`, matching the code's
`strings.Contains(base, "credential")`, and `credential.txt` joins the examples. The conflict-resolution
rule now says the archived/in-flight distinction plays no part, and says what actually decides a tie:
the first claim found, with `working/` read first. The publication sentence no longer generalizes across
two commands that differ — the screenshot install verifies and the download does not, and the sentence
now says which does what.

**Eight more in the same two sections, each found by running a fixture.** The `XD` tag covers a
secret-*derived* deletion as well as a secret-shaped one. `X` gained the reading rule it never had,
which matters because `X` is produced for ordinary additions promoted by the ambiguity rule, and a
reader who applies the closest listed rule reads one as a new file and opens it. The by-hand inventory
now drops the metadata the tool drops, classifies a rename origin without giving it a row, and writes
the quarantine the way `start` writes it — replacing, not adding, and with no overlay of its own. The
by-hand association now counts a `working/` REQ whatever its status says, which is what stops every
claimed REQ being filtered out, and uses `find` instead of a `**` glob that needs `globstar` and still
misses the files sitting directly in `do-work/archive/`.

**The by-hand fallbacks are the reason this is not cosmetic.** They are what a person executes when the
tool will not run. A pattern narrower than the code's leaves a secret unquarantined; a status test the
tool does not apply leaves every in-flight REQ's files unassociated. Both were wrong.

**Ten more findings are not fixed here.** They are in sections this request never opened — lifecycle
timing, the merge-aware commit diff, the commit file listing, the verified-exact-publication rule and
the portfolio summary. The request's own line is "the sections these three sit in", and widening past it
a third time in one file is how a review finds prose nobody asked for.

## Discovered Tasks

- **Ten more stale claims in the guide's other sections**, at lines 30, 32, 80, 92 (two), 106 (two) and
  128 (two) — lifecycle timing, the merge-aware commit diff, the commit file listing, the
  verified-exact-publication rule and the portfolio summary. Two are outright false rather than
  incomplete. Captured as **REQ-597**, with the sweep's report as its input.

## Qualification

**Passed.** Read from the range `06be3df6..7c296fc0`, one file, 8 insertions and 7 deletions. Canonical
`qualify` and `scope-drift` both satisfied.

- **Each of the eleven replacements was written from a fixture result, not from a reading.** The sweep
  produced its findings by running the real command against a purpose-built repository: `credential.txt`
  and `MyCredential.YAML` came back as excluded rows and neither matches the old glob; an ordinary
  `brandnew.txt` came back `X` because a secret was present in the same inventory; `git mv
  api-secret.pem plain-name.txt` followed by a delete produced `XD` rows for both paths; a candidate no
  REQ claimed produced no row at all. The three the request names were re-verified by the orchestrator
  against `internal/corehelpers/inventory.go` and `internal/corehelpers/commands.go` before the edits.
- **The two claims with a real cost are the by-hand fallbacks.** They are what a person runs when the
  tool will not. The secret pattern being narrower than the code's means a by-hand inventory leaves a
  file unquarantined that the tool would have caught, and the status test the guide prescribed but the
  tool does not apply means every `claimed` or `blocked` REQ is filtered out and its files come back
  unassociated — which is the exact failure the bullet two lines above it exists to prevent.
- **The publication sentence no longer generalizes across a difference.** ~~That shape has now shipped
  twice in this file, in REQ-555 and again here.~~ **Corrected after review:** the sentence replaced here
  predates the visible per-REQ history — `git log -S` puts it in the squashed import — and the visible
  history holds three instances, not two: REQ-555, REQ-595 one commit earlier, and this one. Naming two
  and skipping the one directly above understates the pattern the paragraph exists to argue. Every time,
  a review caught it rather than the author. The replacement states each command separately
  and keeps only the two clauses that are true of both.
- **All three guards over this file exit 0**: `audit-lockins.sh` (route column, orchestration claim,
  Mechanics-column shell vocabulary), `prescribed-shell-canonicalization.sh` (twelve headings, sixteen
  pointer sites), and the new `quiet-grep-pipeline-audit.sh`.
- **The ten findings outside these two sections are captured, not fixed.** Widening a third time in one
  file is how a review ends up finding prose nobody asked for; the request drew its own line at "the
  sections these three sit in", and REQ-597 carries the rest with the sweep's report as its input.

### Remediation qualification (after review)

**Passed.** Remediation range `41313790..0bbd10d`, two files. The review scored 72% and asked for
remediation, and it earned that by **executing both by-hand fallbacks** rather than reading them: a
reviewer transcribed each procedure literally into a script and diffed its output against the tool.
Neither reproduced it. That is this request's own headline acceptance criterion.

- **The classification clause was mine and it was wrong.** "A rename or copy destination is M" was
  written with no condition; `classifyInventory` tests `strings.Contains(status, "D")` **first**, so
  porcelain `RD` and `CD` classify the destination as D, and as XD when the origin is secret-shaped. Two
  divergences followed: an ordinary deleted rename destination came out M by hand and D from the tool,
  sending the reader to the M rule — `git diff` on a file that is gone — and a secret-derived one came
  out X by hand and XD from the tool, so the by-hand run quarantined a path the tool permits staging.
  Worse, the same commit widened the XD legend four lines above to describe exactly that path, so the
  shipped file gave two different tags for one file. The clause now states the code's own precedence.
- **The globstar parenthetical was false in its own right.** "A `**` glob needs `shopt -s globstar` and
  **still** misses the files sitting directly in `do-work/archive/`" — with globstar on, `**/` matches
  zero or more directories, so those files are matched. Only the first half was true. The record's
  Exploration had it right and the shipped sentence compressed it into a falsehood, which means that
  replacement was never re-checked against the evidence cited for it.
- **The prescribed command is hardened three ways**, each closing a divergence a reviewer measured:
  `-type f` skips the symlinks the walk skips (a symlinked REQ took a path by hand that the tool gave to
  another); `2>/dev/null` survives a project that has never archived, where `find` writes to stderr and
  exits 1 and `set -euo pipefail` aborts the fallback — in exactly the situation a fallback is most
  needed; and `LC_ALL=C sort` makes the tie-break reproducible, because `find` returns raw directory
  order where the walk reads each directory name-sorted, so a tie answered the same way on two machines
  by the tool was answered differently by hand.
- **One sentence taught a two-way state as one-way.** "A candidate missing from that output is
  unassociated" is also true when the quarantine still holds the path from an earlier inventory, and
  `actions/commit.md` states that as the design. The unsafe direction is the one the sentence taught.
- **"The download inspects nothing afterwards" was checkable and false**: it stats the published target
  to report a byte count and discards the error. The point stands and the wording did not.
- **A guard now exists for the class this request was created for.** It derives the expected patterns
  from `secretPath` itself — each `HasPrefix` must appear as `X*`, each `HasSuffix` as `*X`, each
  `Contains` as `*X*` — and requires the guide to name every one. It fails on exactly the `*credentials*`
  drift, on a new code pattern the guide lacks, on the guide dropping one, and on `secretPath` being
  renamed away. The record's "nothing here is guardable by a text scan" was wrong, and a reviewer built
  the prototype that proved it.
- **Two findings outside this request are captured**: the caller drift in `actions/commit.md`, added to
  REQ-597's write set, and a tool bug where in-flight-ness is decided by a substring of the absolute
  path, so a checkout beneath any directory named `working` makes every archived REQ active — REQ-599.

## Testing

**No test reads the sentences that changed, and the correctness evidence is the sweep's fixtures.**
Each of the 21 findings came from running the real command against a purpose-built repository under
`scratchpad/audit596/`, not from reading the code; the three the request names were re-verified
independently before the edits.

- `bash _dev/tests/audit-lockins.sh` — `Audit lock-in regressions passed.`, exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` — passed, exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — `quiet-grep pipeline audit passed (94 tracked shell
  files, 14 must-flag and 7 must-not-flag shapes).`, exit 0.
- `bash _dev/tests/maintainer-verify.sh` at the pre-change revision `06be3df` —
  `Maintainer verification passed.`, exit 0, gate wall 88s.

Unlike REQ-595, nothing here is guardable by a text scan: these are by-hand procedures whose correctness
is a comparison between prose and behaviour, and the only instrument for that is a fixture plus a reader.

### Remediation testing (after review)

**The new guard was proven red four ways**, each against the real files with them restored from a green
copy between runs:

- The pre-REQ-596 plural restored: `secretPath tests `*credential*` and prescribed-shell-primitives.md
  does not name it; a by-hand inventory built from the guide would miss what the tool excludes.` — exit 1.
  This is the drift the request was created for, now caught.
- A new pattern added to `secretPath` that the guide does not name (`.jks`): caught, exit 1.
- The guide dropping a pattern the code tests (`*.pem`): caught, exit 1.
- `secretPath` renamed away: `secretPath is gone from inventory.go (awk exit 3); the secret-pattern
  drift check cannot run.` — exit 1, rather than a clean read.

Restored: `Audit lock-in regressions passed.`, exit 0.

**The expected patterns are derived, not listed.** Adding a pattern to the code without adding it to the
guide fails; so does narrowing one in either place. The block's comment states what it does not catch —
that the guide may name extras, and that the surrounding prose can still describe the classification
wrongly, which is exactly what this remediation had to fix by hand.

**Guards and gate.** `bash -n` and `shellcheck --severity=warning` on the lock-in exit 0. All three
guards over the guide exit 0: `audit-lockins.sh`, `prescribed-shell-canonicalization.sh` and
`quiet-grep-pipeline-audit.sh`. The fast gate at the remediation revision:
`Maintainer verification passed.`, exit 0, gate wall 83s.

## Review

**Overall: 72%** | 2026-09-06T08:15:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 62% |
| Code Quality | 68% |
| Test Adequacy | 55% |
| Scope | 88% |
| Risk | Moderate |
| Acceptance | Remediate before archive |

**Verdict: Remediate.** Eight of the eleven replacements are true and verified against the code or a
fixture, including all three the request names; the `*credential*` correction and the conflict-resolution
rewrite are real safety and correctness fixes. But the change's own headline acceptance criterion failed:
**a person following either by-hand fallback still got a different answer from the tool**, and one newly
written sentence was false in its own right.

The reviewers did not read the fallbacks. They transcribed each one literally into a script and diffed
it against `protected-inventory.sh`. On a broad fixture — staged rename with a secret origin, a copy from
a secret origin, an ordinary copy, a deletion, a modification, an untracked addition, hidden `do-work/`
metadata — the two outputs were byte-identical. They diverged on porcelain `RD`, and on three separate
things in the association fallback.

Where the reviewers disagreed, and what was picked:

- All three found the `RD` divergence and the false globstar parenthetical independently; no dispute.
- Only one reviewer found the tie-break ordering divergence. Reproduced by the synthesizer with eight
  archived REQs where two claim one path at the same `completed_at`: the tool answers one, `find`
  presents the other first.
- One reviewer scored three pre-existing defects in unopened sections into this change at 78. The
  synthesizer verified all three as genuinely false but kept them out of this request's score.

**Important findings:**

- The by-hand inventory tags a deleted rename destination `X` where the tool tags `XD`, and `M` where the
  tool tags `D` — and the `XD` legend widened in the same commit says `XD`. The clause "a rename or copy
  destination is M" was written with no condition; the code tests for `D` first. No secret is exposed in
  either direction, but the by-hand run quarantines a path the tool permits staging, and the ordinary case
  sends a reader to `git diff` on a file that is gone. — impact-user-visible → fixed in remediation

**Minor findings:**

- "A `**` glob needs `shopt -s globstar` and **still** misses the files sitting directly in
  `do-work/archive/`" is false: with globstar on, `**/` matches zero directories and those files are
  matched. — impact-user-visible → fixed in remediation
- The by-hand association names a different owner than the tool on a `completed_at` tie, because `find`
  returns raw directory order and the walk reads name-sorted. Not exotic: two in-flight REQs both lack
  `completed_at`. Also not reproducible across machines, since the directory hash seed is per-filesystem.
  — impact-user-visible → fixed in remediation
- "A candidate missing from that output is unassociated" is false whenever the quarantine retains the
  path from an earlier inventory, which `actions/commit.md` states as the design. —
  impact-user-visible → fixed in remediation
- One sentence gave contradictory instructions on whether a rename origin gets a row of its own. —
  impact-negligible → fixed in remediation
- The prescribed `find` fails on a project that has never archived, where the tool tolerates a missing
  root, and lists symlinked REQ files the walk skips. — impact-negligible → fixed in remediation
- "The download inspects nothing afterwards" is contradicted by the `os.Stat` it runs. The discarded
  error is also a latent nil-FileInfo dereference. — impact-negligible → fixed in remediation
- The `XD` reading bullet still said "Deleted secret-shaped files" after the same commit widened the
  legend. — impact-negligible → fixed in remediation
- The record's "nothing here is guardable by a text scan" is false; a reviewer built the fifteen-line
  prototype and ran it at both ends of the range. — impact-rule-change → the guard is now shipped
- `actions/commit.md` now contradicts the guide on both the glob and the associate output. —
  impact-rule-change → added to REQ-597's write set
- The tool decides in-flight-ness with a substring of the absolute path, so a checkout beneath any
  directory named `working` makes every archived REQ active. The prose is right and the code is wrong. —
  impact-user-visible → REQ-599
- The guide now contradicts itself about whether a rename over an occupying directory nests: line 107
  describes shell `mv`, the new line 125 describes `rename(2)`. The `ln`/`mv` half is already in
  REQ-597. — impact-negligible → report only
- The record's "shipped twice in this file" history is wrong: the sentence replaced here predates the
  visible history, and the visible history holds a third instance — REQ-595, one commit earlier. —
  impact-negligible → corrected below

**Requirements checklist:**

- [x] The three claims the request names are corrected — delivered, each against the code
- [x] The eight more in the same two sections are corrected — delivered
- [ ] → [x] The by-hand fallbacks match what the tool does — **not delivered at review**, five
  divergences measured by execution; **delivered in remediation**
- [ ] → [x] No sentence generalizes over two commands that behave differently — **not delivered at
  review**: the shape moved from the two-command axis to the `R`/`RD` axis; **delivered in remediation**
- [x] Both guards over the file still exit 0 — delivered, and a third now guards the pattern set
- [x] The ten findings outside these two sections are captured rather than fixed — delivered as REQ-597

**Acceptance testing**

**Result: Remediate at review, Pass after remediation.** Three reviewers built fixture repositories and
ran both procedures against the tool. Every divergence they measured is closed, and the pattern half of
the class is now pinned by a guard rather than by a reader.

**Follow-ups created:** 2 — REQ-599 for the tool's path-substring bug, and the caller drift added to
REQ-597.

*Reviewed by review-work action*

## Lessons Learned

- **A procedure is not checked until it is executed.** Eleven sentences were corrected against the code
  and eight of them were right. The three that were wrong only showed up when a reviewer transcribed the
  whole paragraph into a script and diffed it against the tool. Reading a fallback tells you whether each
  sentence is true; running it tells you whether the sequence produces the tool's answer, which is the
  only thing a fallback is for.
- **The generalizing sentence moved axis rather than going away.** This request existed partly to fix a
  sentence that generalized across two commands, and it shipped one that generalized across `R` and `RD`
  — "a rename or copy destination is M", with no condition, where the code tests for `D` first. Fixing an
  instance of a shape is not the same as learning the shape.
- **A correction can contradict another correction in the same commit.** The `XD` legend was widened four
  lines above the clause that then tagged the same file `X`. Two edits, each defensible alone, disagreeing
  in the shipped file. When one pass changes several statements about one mechanism, re-read them
  together at the end, not one at a time as they are written.
- **Compressing a verified sentence is a rewrite.** The Exploration had the globstar behaviour right; the
  shipped one-line version added the word "still" and became false. A shorter version of a checked claim
  is an unchecked claim.
- **"Nothing here is guardable" was wrong twice in two requests.** REQ-595's record said no lane could
  distinguish a prose change, and a reviewer showed a guard was feasible. This record said the same and a
  reviewer built the prototype. Both times the guard turned out to be fifteen lines. The reflex to reach
  for is "what is the smallest thing that would have caught this", not "this class is unguardable".
- **Order can be part of a procedure's answer.** The tie-break sentence said "the first claim found
  stands" and the prescribed command was `find`, which returns raw directory order where the tool reads
  name-sorted. The sentence and the command were each defensible; together they made the answer depend on
  a filesystem hash seed.

## Orientation

`skills/do-work/docs/prescribed-shell-primitives.md` § **Protected inventory fallbacks** owns what an
inventory tag means, how each tag class is read, what association settles, and the by-hand procedure for
each mode. Two of those are **procedures a person executes when the tool will not run**, and they are
held to a stronger standard than the rest of the guide: the output must match
`scripts/protected-inventory.sh`.

Three things about them are easy to get wrong and were:

- **Classification order.** Any status containing `D` is D; otherwise `??` or an index-column `A` is A;
  otherwise M. So a rename or copy destination is M only when its record reports no deletion — porcelain
  `RD` is a deletion.
- **The rename origin** gets a row of its own only when the record is a rename *and* the origin is
  secret-shaped, in which case the row is `XD`. A copy origin never gets one.
- **The association `find`** needs `-type f` (the walk skips symlinks), `2>/dev/null` (the walk tolerates
  a missing `do-work/archive/`, which a project that has never archived does not have), and
  `LC_ALL=C sort` (the walk reads each directory name-sorted; `find` does not, and the tie-break depends
  on order).

`X` is not only "a secret". An ordinary addition is promoted to `X` whenever any `X` or `XD` is present in
the same inventory, because Git cannot identify a copy when both source and destination are untracked.
That is why `X` has its own reading rule and why the rule says not to relax it because the name looks
ordinary.

**The pattern set is guarded.** Finding 3 in `_dev/tests/audit-lockins.sh` derives the expected globs
from `secretPath` itself and requires the guide to name each one. It catches a pattern added to the code
and not the guide, a pattern dropped from either, and `secretPath` being renamed away. It does **not**
catch the guide naming extras, and it cannot check that the surrounding procedure is right — that is what
the fixtures are for.

Recorded and unfixed: the guide's line 107 describes shell `mv` nesting into a directory while line 125
describes `rename(2)` refusing, and the two read as one contradiction — REQ-597 owns line 107.
`actions/commit.md` still names the glob this request replaced and a `-` row the command never prints,
added to REQ-597. And the tool decides in-flight-ness from a substring of the absolute path, so a
checkout beneath a directory named `working` makes every archived REQ active — REQ-599.
