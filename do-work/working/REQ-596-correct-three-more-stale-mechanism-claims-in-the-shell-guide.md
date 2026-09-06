---
id: REQ-596
status: claimed
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
- **The publication sentence no longer generalizes across a difference.** That shape — one sentence true
  of some of the things it covers — has now shipped twice in this file, in REQ-555 and again here, and
  both times a review caught it rather than the author. The replacement states each command separately
  and keeps only the two clauses that are true of both.
- **All three guards over this file exit 0**: `audit-lockins.sh` (route column, orchestration claim,
  Mechanics-column shell vocabulary), `prescribed-shell-canonicalization.sh` (twelve headings, sixteen
  pointer sites), and the new `quiet-grep-pipeline-audit.sh`.
- **The ten findings outside these two sections are captured, not fixed.** Widening a third time in one
  file is how a review ends up finding prose nobody asked for; the request drew its own line at "the
  sections these three sit in", and REQ-597 carries the rest with the sweep's report as its input.

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
