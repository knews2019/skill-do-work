---
id: REQ-594
status: claimed
domain: testing
created_at: 2026-09-06T03:08:29Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
route: C
planning_at: 2026-09-06T04:34:44Z
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: false
depends_on: [REQ-593]
related: [REQ-593]
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-09-06T03:37:11Z
  basis:
    - Route C
    - 24-file write set
    - 3 subsystems involved
    - 4 acceptance criteria
write_set: [_dev/tests/quiet-grep-pipeline-audit.sh, _dev/tests/contracts/probe-lanes.sh, _dev/tests/update-script-behavior.sh, _dev/tests/contracts/core-checks.sh, _dev/tests/contracts/queue-kanban.sh, _dev/tests/select-simple-reqs-behavior.sh, _dev/tests/p50-estimator-determinism.sh, _dev/tests/audit-lockins.sh, _dev/tests/staged-skills-contract.sh, _dev/tests/install-suite-behavior.sh, _dev/tests/prescribed-shell-cases/qualify.sh, _dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh, _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh, _dev/tests/prescribed-shell-cases/generate-report-image.sh, _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh, _dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh, _dev/tests/prescribed-shell-cases/atomic-download.sh, _dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh, _dev/tests/prescribed-shell-cases/protected-inventory.sh, _dev/tests/prescribed-shell-cases/capture-screenshot.sh, _dev/tests/prescribed-shell-cases/show-commit-diff.sh, _dev/tests/prescribed-shell-cases/lexical-memory-recall.sh, _dev/tests/prescribed-shell-cases/install-memory-hooks.sh, _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh, skills/do-work/tools/select-simple-reqs.sh]
title: 'Generalize the SIGPIPE fix: about 130 quiet-grep pipelines remain across the maintainer test tree'
claimed_at: 2026-09-06T03:36:45Z
---

# Generalize the SIGPIPE Fix Across the Maintainer Test Tree

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Three read-only censuses plus a judge that re-measured every load-bearing claim by
  running commands. `_dev/primes/prime-shell-commands.md` read before any shell was written. Approach:
  convert the sites whose producer status is already asserted mechanically, give the sites whose
  producer can fail an explicit status assertion, then move the guard out of one probe file and give it
  the whole repository.
- [x] **[APPLY]:** 25 files changed, all inside the declared write set: 24 under `_dev/tests/` plus the
  one shipped line the guard's scope would otherwise have to be drawn around. Three builders in
  sequence in one worktree, twelve commits, so a failure is attributable to one step.
- [x] **[UNIFY]:** `git diff --stat` over the merge — 25 files, 488 insertions, 249 deletions. Every
  converted file re-scanned individually with REQ-593's own scanner and reporting zero; `bash -n` and
  the ShellCheck lane form (every tracked `*.sh` in one invocation) both exit 0. The whole-tree scanner
  reports **0** offending logical lines, down from 129. No debug artifacts: the one site the first pass
  rewrote outside the defect class was reverted and the conversion tool tightened to require an
  early-leaving option.

## What

REQ-593 removed a defect where a writer piped into `grep -q` under `set -o pipefail` reported the
writer's SIGPIPE death as grep's verdict. It fixed eleven sites and added a source scanner — but the
scanner reads `${BASH_SOURCE[0]}`, its own file and nothing else. The class is closed in **one file
of twenty-three**.

`grep -rc -E '\|[[:space:]]*grep[[:space:]]+-[A-Za-z]*q' _dev/tests/` returns 23 non-empty files and
roughly 130 sites, including nineteen in the gate's own `_dev/tests/contracts/core-checks.sh` and
thirty-one in `_dev/tests/prescribed-shell-cases/qualify.sh`. `core-checks.sh` already disagrees with
itself: one site uses the herestring form, another still pipes.

The consequence is measured, not theoretical. The defect is silent below roughly 36 KB of producer
output and certain above about 200 KB, and it is wrong in **both** directions — a positive matcher
misses a pattern that is present, and a negative matcher fails to flag one it should. Two scripts have
already produced false failures in the gate during a single work run.

## Requirements

- No check anywhere under `_dev/tests/` decides on a quiet grep fed by a pipeline whose writer can die.
- The guard is repository-wide rather than per-file, and it is not defeatable by ordinary shell:
  REQ-593's scanner was evaded by five spellings a reviewer found in minutes — the pipe at end of line
  with the grep on the next, `grep --quiet`, `grep --silent`, `| LC_ALL=C grep -q`, and
  `| command grep -q`.
- Each conversion preserves what the assertion measured. REQ-593 found that capturing a producer's
  output and discarding its exit status can silently narrow an assertion — a truncated archive whose
  partial listing still contains the marker passed where the pipeline form failed.
- A probe fixture carries every evasion spelling, so the widened guard is itself pinned.

## Context

Found during the independent review of REQ-593, which fixed the two named matchers and their nine
siblings in `_dev/tests/update-script-behavior.sh`. This request is the generalization that request's
own Scope deferred, promoted from a note to a request because the review showed the deferral had no
destination.

## Full Context

See the REQ-593 review in `do-work/archive/` (or `do-work/working/` while it is in flight) for the
five evasion spellings with their reproductions, and the site census by file.

## Triage

**Route: C** — Explore, plan, then build.

**Reasoning:** The technique is settled — REQ-593 shipped a widened scanner that catches all five
evasion spellings — but the census is not. `grep -rc -E '\|[[:space:]]*grep[[:space:]]+-[A-Za-z]*q'
_dev/tests/` returns 23 files and about 139 matches, and a spot check of the two files REQ-593 already
fixed shows two of their four matches are comments and one is a herestring, so the real site count is
unknown and the raw number over-reports. Each conversion is also a judgment: REQ-593 proved that
capturing a producer's output and discarding its exit status can silently narrow the assertion it
replaced, so a mechanical rewrite of 130 sites is exactly the shape that ships 130 quiet narrowings.
Every site must be read for what its assertion measures before it is converted.

**Planning:** Required. The plan must settle where the repository-wide guard lives, how it avoids
flagging its own documentation and the fixture that pins it, and the per-file conversion order so a
failure is attributable.

**One requirement is already partly met and must not be re-done.**
`_dev/tests/update-script-behavior.sh` carries `quiet_grep_pipeline_offenders`, which already matches
on the defect's ingredients rather than one spelling and already has a five-shape fixture. The work is
to give it a repository-wide scope and a home that is not one probe file — not to write a second
scanner beside it.

## Plan

Three read-only censuses plus a judge that re-measured everything itself rather than trusting them.
The judge corrected the censuses and the competing plan in six places; those corrections are carried
below. One of the two planning agents failed on its structured output and returned nothing, so the
judge worked from one plan plus its own measurements — stated here rather than left to look like two.

### The measured census

**129 offending logical lines carrying 135 grep invocations, in 22 of the repository's 93 tracked shell
files.** Measured by running REQ-593's shipped scanner verbatim over `git ls-files -z -- '*.sh'`, not by
the request's regex, which both over-reports (comments, herestrings) and under-reports (`--quiet`,
`--silent`, continuations, prefixed greps).

- Fast tier, 32 lines: `contracts/core-checks.sh` 13 (carrying 19 greps — the only file where lines and
  invocations differ), `select-simple-reqs-behavior.sh` 15, `contracts/queue-kanban.sh` 2,
  `p50-estimator-determinism.sh` 1, `audit-lockins.sh` 1.
- Heavy tier, 90 lines in `prescribed-shell-cases/` across 14 of that directory's 18 files, led by
  `qualify.sh` at 31, `audit-archive-timestamps.sh` at 15 and `repair-req-timestamps.sh` at 13; plus
  `staged-skills-contract.sh` 4 and `install-suite-behavior.sh` 2.
- **One shipped file**: `skills/do-work/tools/select-simple-reqs.sh:47`. That single line is what makes
  this a release.

**45 of the 129 are the silent-pass class** — 35 carrying `&& fail` and 10 more in
`if <matching grep>; then fail` form — where a dead writer makes the assertion pass rather than fail.
The competing plan said 28; the judge measured 45.

**Two sites are already broken, and the judge reproduced both.**
`staged-skills-contract.sh:82-86` and `install-suite-behavior.sh:142` are
`find … -print -quit | grep -q .` guards: with a stubbed `find` that prints a leftover path and then
exits 1, `pipefail` makes the guard stay silent and report the legacy runtime as retired. The
capture-with-status form fires correctly.

**The shipped scanner has a real hole.** A comment line between the pipe and the reader — `producer |`,
then `# note`, then `grep -q marker` — parses and runs as one pipeline, and the scanner misses it. Both
the bare-pipe and the backslash-continuation forms were confirmed with `bash -n` and by execution. A
two-line awk fix catches both, and its output is byte-identical to the shipped scanner across all 93
files (129 findings each).

### Decisions

- **D-01 — the guard moves to its own file and runs in the fast tier.** New executable
  `_dev/tests/quiet-grep-pipeline-audit.sh` holds the scanner (moved, not copied), the repository walk
  and the fixture. Two lines register it in `_dev/tests/contracts/probe-lanes.sh` above the heavy block.
  Measured cost: 0.24s over 93 files against a 30s per-file fast budget. A violation now blocks an
  ordinary commit rather than only a `--heavy` run, which is the point.
- **D-02 — the scope is a condition, not a list.** `git ls-files -z -- '*.sh'` is the same enumeration
  the ShellCheck lane already uses at `_dev/tests/maintainer-verify.sh:670`. Verified that the suffix is
  not a partial enumeration: `git grep -lI '^#!.*\(ba\)\?sh'` over tracked files, minus `*.sh`, returns
  nothing.
- **D-03 — the shipped line is converted and the release is paid.** Drawing the guard's scope around
  `_dev/tests/` to avoid a version bump would be an exemption list with one entry, aimed at the one
  known instance of the defect the guard exists to forbid. The conversion is provably lossless: that
  producer's status is already captured and asserted four lines earlier.
- **D-04 — the guard's self-exemption is three mechanisms and none of them is a list.** Prose is `.md`
  and outside the walk by construction. The fixture heredocs interpolate the pipe and hash characters,
  so the audit file's own bytes never carry the shape and it scans itself clean. And a `#` suppresses
  exactly what bash suppresses, so moving a real offender into a comment deletes the check rather than
  hiding it.
- **D-05 — the fixture asserts by name, not by count.** 14 must-flag shapes and 7 must-not-flag shapes,
  each paired with its name so a failure says which shape was lost. Mutation-tested: dropping `--quiet`
  loses shape 3, `--silent` loses 4, `--max-count` loses 12, the `[qm]` class loses 11, the `(e|f)?`
  alternation loses 13, the comment fix loses both 14 and 17. The must-not-flag half is not decoration —
  the request's own census regex flags 2 of those 7.
- **D-06 — `verify_output_matchers_read_greps_verdict` stays in `update-script-behavior.sh`.** It drives
  the two shipped matchers over a 256 KiB variable and is the only thing proving the mechanism on real
  data. It is not the scanner and must not travel with it.

### Ordered steps

1-10. The mechanical sweep, ~110 sites where the producer replays an already-captured variable and its
status is already asserted on the capture line: `contracts/core-checks.sh`, `contracts/queue-kanban.sh`,
`select-simple-reqs-behavior.sh`, `p50-estimator-determinism.sh`, `audit-lockins.sh`, then the
`prescribed-shell-cases/` files including the 21 `find … -print -quit | grep -q .` sites.
11-12. The judgment sites: `staged-skills-contract.sh` (4) and `install-suite-behavior.sh` (2), which
include the two already-broken guards and two `git check-attr` positive tests for a forbidden state,
where a dead producer yields no output, no match and a silent pass.
13. The one release: convert `skills/do-work/tools/select-simple-reqs.sh:47`, bump, changelog, mirror.
14. Move and widen the guard, register it, and delete the moved block from `update-script-behavior.sh`.

**Step order is load-bearing**: 13 must land before 14, or the repository-wide guard is red between them.

### Risks carried forward

- Two thirds of the sweep (90 of 128 `_dev/tests` sites) is heavy-only, reached through
  `staged-skills-contract.sh:189`, so those conversions are not exercised by a fast run.
- The scanner splits on a bare pipe, so a `|` inside a quoted grep pattern would create a phantom stage.
  Zero false positives across all 93 files today.
- The reader set is grep/egrep/fgrep only. `rg -q`, `head`, `sed -n '/x/{p;q;}'`, `awk '/x/{exit}'` and
  `read` are the same condition with a different reader and are not covered.
- Herestrings append a newline that `printf '%s'` did not. Identical for every `grep -q`/`-qx`/`-qxF`
  here; it would not be for a byte-exact reader.

*Validated by the judge against its own measurements; the three censuses and the surviving plan are in
the run directory.*

## Exploration

Three read-only census agents split the tree by directory, and a judge re-measured every load-bearing
claim by running commands rather than reading reports. Full census in the run directory as
`REQ-594-census.json`; the judged plan is `REQ-594-judged-plan.json`.

**The request's own numbers are close but its method is not.** "About 130 across 23 files" is right at
129 lines across 22 files, but the regex it prescribes is wrong in both directions: it counts comments
and already-converted herestrings, and it cannot see `--quiet`, `--silent`, a continuation, or a
prefixed grep. The measurement that matters came from running REQ-593's own scanner over
`git ls-files -z -- '*.sh'`.

**The polarity split is the number that says what this buys.** 45 of the 129 are negative assertions
where a dead writer makes the check pass silently — 35 written `&& fail` and 10 as
`if <matching grep>; then fail`. The rest are positive assertions, where the same failure is loud. A
plan that treats all 129 as equally urgent is not reading the defect.

**Two of the 129 are not latent.** `staged-skills-contract.sh:82-86` and `install-suite-behavior.sh:142`
were reproduced broken with a stubbed `find` that prints a path and exits 1: `pipefail` makes the guard
report the legacy runtime as retired. Those two are bug fixes, not conversions.

**The scanner REQ-593 shipped can be walked past with a comment.** `producer |` on one line, a `#` note
on the next, the reader on the third: bash builds and runs the pipeline, the scanner does not see it.
Confirmed both by `bash -n` and by execution (`seq 1 5 | / # c / grep -c .` prints 5). The
backslash-continuation form has the same hole. A two-line awk fix closes both and produces byte-identical
output to the shipped scanner across all 93 tracked files.

**Six claims the censuses and the surviving plan made did not hold**, and each was corrected by
measurement: the negative-assertion count (28, actually 45), a stale byte figure, a wrong line number
for the ShellCheck lane's enumeration, a claim that the working tree was dirty (it is clean), a
directory file count, and a right conclusion reached by wrong reasoning about manifest subtree coverage.

*Generated by three Explore agents and one judge; one of the two planning agents failed on structured
output and contributed nothing.*

## Scope

**Files I will touch:**
- `_dev/tests/quiet-grep-pipeline-audit.sh` (new, executable) — the scanner moved out of the updater probe with the comment fix, the repository walk, and the two-sided named-shape fixture
- `_dev/tests/contracts/probe-lanes.sh` (modify) — two lines registering the audit as a fast-tier probe
- `_dev/tests/update-script-behavior.sh` (modify) — delete the moved block; keep the matcher self-check, which drives the shipped matchers over a 256 KiB variable and is not the scanner
- `_dev/tests/contracts/core-checks.sh` (modify) — 13 logical lines carrying 19 grep invocations
- `_dev/tests/contracts/queue-kanban.sh` (modify) — 2 sites
- `_dev/tests/select-simple-reqs-behavior.sh` (modify) — 15 sites
- `_dev/tests/p50-estimator-determinism.sh` (modify) — 1 site
- `_dev/tests/audit-lockins.sh` (modify) — 1 site
- `_dev/tests/staged-skills-contract.sh` (modify) — 4 sites, all needing judgment, one of them already broken
- `_dev/tests/install-suite-behavior.sh` (modify) — 2 sites, both needing judgment, one of them already broken
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modify) — 31 sites
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (modify) — 15 sites, two of them three-stage
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` (modify) — 13 sites
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (modify) — 6 sites
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (modify) — 6 sites
- `_dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh` (modify) — 4 sites
- `_dev/tests/prescribed-shell-cases/atomic-download.sh` (modify) — 3 sites
- `_dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh` (modify) — 3 sites
- `_dev/tests/prescribed-shell-cases/protected-inventory.sh` (modify) — 3 sites, one of them the security-relevant negative assertion in this sweep
- `_dev/tests/prescribed-shell-cases/capture-screenshot.sh` (modify) — 2 sites
- `_dev/tests/prescribed-shell-cases/show-commit-diff.sh` (modify) — 1 site
- `_dev/tests/prescribed-shell-cases/lexical-memory-recall.sh` (modify) — 1 site
- `_dev/tests/prescribed-shell-cases/install-memory-hooks.sh` (modify) — 1 site
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modify) — 1 site
- `skills/do-work/tools/select-simple-reqs.sh` (modify) — one line, and the reason this request is a release

**Files I will NOT touch:** `_dev/tests/contract-regressions.sh`, which is a 77-line file sitting at a
77-line ratchet — the registration goes to `probe-lanes.sh` instead. `_dev/tests/fast-stages.json` and
`_dev/tests/heavy-lanes.json`, which declare Go test stages and never enumerate shell test files.
`_dev/primes/prime-shell-commands.md` and `skills/do-work/docs/prescribed-shell-primitives.md`, which
are prose and outside the walk by file type rather than by exemption.

**Acceptance criteria:**
- [ ] The repository-wide audit reports zero offenders over every tracked shell file, and it fails when
  any one of the fourteen must-flag shapes is reintroduced
- [ ] The fixture fails if any single evasion is removed from the scanner, and says which shape was lost
- [ ] The audit file scans itself clean, without a path exemption
- [ ] Every conversion preserves what its assertion measured; the sites whose producer can fail get an
  explicit status assertion rather than a bare herestring
- [ ] The two already-broken guards fire on the input that used to pass silently
- [ ] The release lands before the repository-wide guard, so the gate is never red between steps

## Pre-Flight

**Green gate at `ceeea69`**, the revision the builder branches from.
`bash _dev/tests/maintainer-verify.sh` printed `Maintainer verification passed.` and exited 0, gate
wall 61s. One `SKIP` line, the heavy-only one every fast run prints.

**A fast gate is not sufficient evidence for this request, and the builder is told so.** 96 of the 128
`_dev/tests` sites live under `prescribed-shell-cases/`, reached only through
`staged-skills-contract.sh:189`, which `probe-lanes.sh` registers inside the heavy block. Two thirds of
the conversions are invisible to a fast run, so the `staged-skills`, `updater` and `installer` heavy
lanes are required evidence, each reporting its own exit line.

**The baseline is re-measured at the branch point rather than replayed from the request.** REQ-593's
shipped scanner over `git ls-files -z -- '*.sh'` reports 129 offending logical lines in 22 of 93 files;
the request's own regex reports 140 matches over 23 files. The scanner's number is the one this request
drives to zero.

**The environment the gate needs**, recorded once: `NODE_OPTIONS` and the `GIT_CONFIG_COUNT` /
`GIT_CONFIG_KEY_*` / `GIT_CONFIG_VALUE_*` triples unset, `GIT_CONFIG_GLOBAL` pointed at a config
carrying `commit.gpgsign = false`, and `QUEUE_KANBAN_BROWSER` at the container's Chromium.

**The builder works in an isolated worktree** at
`../skill-do-work-worktrees/worktree-agent-REQ-594-quiet-grep-pipelines`, branched from `ceeea69`, and
hands back one file to the main checkout without staging or committing anything there.

## Implementation Summary

**Files changed:** 25, every path repository-relative and every one inside the declared write set.

- `_dev/tests/quiet-grep-pipeline-audit.sh` — new, executable — the guard's new home
- `_dev/tests/contracts/probe-lanes.sh` — two lines, fast-tier registration
- `_dev/tests/update-script-behavior.sh` — 94 lines deleted, deletion only
- `_dev/tests/contracts/core-checks.sh` — 13 logical lines carrying 19 grep invocations, one of them split into if/elif
- `_dev/tests/contracts/queue-kanban.sh` — 2 sites
- `_dev/tests/select-simple-reqs-behavior.sh` — 15 sites
- `_dev/tests/p50-estimator-determinism.sh` — 1 site
- `_dev/tests/audit-lockins.sh` — 1 site
- `_dev/tests/staged-skills-contract.sh` — 4 sites, all judgment, one of them already broken
- `_dev/tests/install-suite-behavior.sh` — 2 sites, both judgment, one of them already broken
- `_dev/tests/prescribed-shell-cases/qualify.sh` — 31 sites
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` — 15 sites, two of them three-stage
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` — 13 sites
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` — 6 sites
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` — 6 sites
- `_dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh` — 4 sites
- `_dev/tests/prescribed-shell-cases/atomic-download.sh` — 3 sites
- `_dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh` — 3 sites
- `_dev/tests/prescribed-shell-cases/protected-inventory.sh` — 3 sites, one the security-relevant negative assertion
- `_dev/tests/prescribed-shell-cases/capture-screenshot.sh` — 2 sites
- `_dev/tests/prescribed-shell-cases/show-commit-diff.sh` — 1 site
- `_dev/tests/prescribed-shell-cases/lexical-memory-recall.sh` — 1 site
- `_dev/tests/prescribed-shell-cases/install-memory-hooks.sh` — 1 site
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` — 1 site
- `skills/do-work/tools/select-simple-reqs.sh` — one line, and the reason this is a release

**What was done: 129 offending logical lines became 0**, measured with REQ-593's own scanner over
`git ls-files -z -- '*.sh'` at both ends. ~~122 sites were mechanical … Seven needed a decision.~~
**Corrected after review, and the correction is on the axis that matters most here.** The measured split
is **99** variable-replay producers and **30** real-command producers: 21 `find … -print -quit` leak
checks, 2 `find` guards, 2 `sed -n` routing extractions, 2 `git check-attr`, 1 `git diff --cached`, 1
`git add -u` and 1 `associate-files.sh`. Every one of those 30 carries a live producer status, which is
exactly what the request says must be judged rather than swept — and each of them did get an explicit
status assertion. The delivered work was right; the summary described seven decisions where there were
thirty, and its own D-05 already contradicted it by saying find's status had to be asserted at 21 sites.

**Two of the 129 were not latent.** `staged-skills-contract.sh:82` and `install-suite-behavior.sh:142`
are `find … -print -quit | grep -q .` guards. Reproduced broken with a stub that prints a leftover file
and then exits non-zero — which is what real `find` does when a sibling directory is unreadable: under
`pipefail` the guard stayed silent and reported the legacy runtime as retired. Both fire now.

**The guard left one probe file and covers the repository.** `quiet_grep_pipeline_offenders` moved into
`_dev/tests/quiet-grep-pipeline-audit.sh`, which walks `git ls-files -z -- '*.sh'` — the same
enumeration the ShellCheck lane uses — and runs in the fast tier at 0.20s over 94 files. Its fixture
carries 14 must-flag and 7 must-not-flag shapes, each asserted by name, and 16 mutations of the scanner
each produce a failure naming the shape that was lost. The scanner also gained a one-condition fix: a
comment line between the pipe and the reader is invisible to bash and was invisible to the scanner too.
That fix is additive — both scanner bodies produce byte-identical output over every tracked file at the
branch point, 129 findings each.

**`verify_output_matchers_read_greps_verdict` stayed behind.** It drives the two shipped matchers over a
256 KiB variable and is the only thing proving the mechanism on real data. It is not the scanner.

**One shipped line changed**, `skills/do-work/tools/select-simple-reqs.sh:47`, so this request is a
release. Drawing the guard's scope around `_dev/tests/` to avoid that would have been an exemption list
with one entry, aimed at the only known instance of the defect the guard forbids. The conversion is
provably lossless: that producer's status is captured and asserted four lines earlier, so the pipeline
had no producer status left to lose, and the two forms were checked byte- and verdict-identical over
five output shapes.

## Decisions — implementation

- **D-01 — the guard's scope is a condition, not a directory. DECIDE & STATE.**
  `git ls-files -z -- '*.sh'` covers every tracked shell file. Verified that the suffix is not a partial
  enumeration: no tracked shell file lacks it, and one that did would be invisible to the ShellCheck
  lane too — one problem, one place to fix.
- **D-02 — the guard runs in the fast tier. DECIDE & STATE.** A violation now blocks an ordinary commit
  rather than only a `--heavy` run. Measured cost 0.20s over 94 files against a 30s per-file budget.
- **D-03 — the shipped line is converted and the release is paid. DECIDE & STATE.** See above.
- **D-04 — the fixture asserts by name, not by count. DECIDE & STATE.** 16 mutations, each naming the
  shape it lost. A count would have said only that something broke.
- **D-05 — `find`'s exit status is asserted at 20 of the 21 leak checks. DECIDE & STATE**, with its own
  caveat: it is right at those 20 because every search root is created by the fixture before the check. A
  future case searching a root that may legitimately be absent must not copy the form blindly.
  **Corrected after review:** the first version said all 21, while the Qualification below said twenty
  and was right. The 21st, `generate-report-image.sh:414`, is a readiness poll inside a bounded wait
  loop rather than an assertion; it reads only the captured path's emptiness, deliberately and with a
  comment saying why it differs from its twenty neighbours. **Also corrected:** "every search root is
  created by the fixture" is proven for 18. Three of the 21 sit inside `if false; then` blocks that
  predate this change and never execute, so their only evidence is `bash -n`, ShellCheck and the
  scanner. That was in the builders' hand-back and should have been in this record.
- **D-06 — one site was reverted rather than converted.** `repair-req-timestamps.sh:262` pipes into
  `grep -c`, which reads to EOF, so there is no early-leaving reader and no defect. The conversion tool
  was tightened to require an early-leaving option, matching the scanner's own definition, so the diff
  carries only the class this request is about.
- **D-07 — one statement had to be split rather than converted.** `core-checks.sh:349` stayed flagged
  after every reader on it was converted: the scanner splits on a bare pipe and cannot see a
  command-substitution boundary, so a `printf | associate-files.sh` capture inside `$( )` counted as a
  stage feeding three herestring greps that shared the line through continuations. Reaching zero needed
  the statement split into an `if` capture plus `elif` reads. That is the plan's stated R3 risk showing
  up as a real edit.

## Qualification

**Passed.** Read from the merge range `234eda9c..4edde87`, 25 files, 488 insertions and 249 deletions.
Canonical `qualify` and `scope-drift` both satisfied.

- **The number this request exists to drive is measured at both ends with the same instrument.**
  REQ-593's shipped scanner, run verbatim over `git ls-files -z -- '*.sh'`, reports **129** offending
  logical lines across 22 of 93 tracked files at the branch point and **0** across 94 files now. The
  request's own regex reports 140 matches over 23 files; it was not used, because it counts comments and
  already-converted herestrings and cannot see `--quiet`, `--silent`, a continuation or a prefixed grep.
- **Two of the 129 were live defects, reproduced before they were fixed.** The stub is a `find` that
  reaches a leftover file, prints it, and then exits non-zero because a sibling directory is unreadable
  — which is what real `find` does. Under `pipefail` the shipped guard at
  `staged-skills-contract.sh:82` and its twin at `install-suite-behavior.sh:142` stayed silent with the
  leftover file sitting there, and reported the legacy runtime as retired. Both fire now.
- **Seven sites got an explicit producer-status assertion rather than a bare herestring**, and each was
  mutation-tested in both directions. The one that matters most is the twenty `find … -print -quit` leak
  checks: a planted leftover reports the leak, and a missing search root reports "could not search",
  where the old form passed silently.
- **The guard's self-exemption is three mechanisms and none of them is a path list.** Prose is `.md` and
  outside `git ls-files -- '*.sh'` by construction. Both fixture heredocs interpolate the pipe and hash
  characters, so the audit's own bytes never carry the shape — verified by running the whole-tree scan
  with the new file in the index and getting 0 across all 94 files. And a `#` suppresses exactly what
  bash suppresses, so moving a real offender into a comment deletes the check rather than hiding it.
- **The scanner's own fix is additive, and that was proven rather than argued.** Both scanner bodies —
  the shipped one and the comment-fixed one — run over every tracked shell file at the branch point
  produce byte-identical output, 129 findings each. So the fix adds two shapes and removes nothing.
- **The fixture is asserted by name and survives 16 mutations.** Every one of the 21 shapes is covered
  by at least one mutation, and each failure names the shape that was lost rather than a count.
- **The 94-line deletion left the surviving probe intact**, which the updater heavy lane is what proves.
  `verify_output_matchers_read_greps_verdict` stayed behind deliberately: it drives the two shipped
  matchers over a 256 KiB variable and is the only thing proving the mechanism on real data.
- **One site was reverted rather than converted, and one statement had to be split.** Neither is a
  concession — see D-06 and D-07. The split at `core-checks.sh:349` is the plan's own R3 risk arriving
  as a real edit: the scanner splits on a bare pipe and cannot see a command-substitution boundary.
- **The release is paid rather than avoided.** One shipped line is in the diff, so this request carries a
  version bump. Drawing the guard's scope around `_dev/tests/` would have been an exemption list with one
  entry, pointed at the only known instance of the defect the guard forbids.

### Remediation qualification (after review)

**Passed.** Remediation range `ea71f09e..90f1b6e` plus one follow-up commit, one file. The review scored
86% and answered its own central question cleanly — **no conversion narrowed an assertion** — and then
failed the guard on its second requirement.

- **Five ordinary spellings walked past the guard, in two families, and the request's own Requirement
  named exactly this risk.** The parser skipped comment lines but not blank ones, so a blank line after
  a pipe or after a continuation flushed the joined command and the reader was then scanned as a fresh
  pipe-free command. A note after the pipe on the same line left the line not ending in `|`, closing the
  command before the reader arrived. And the option class was anchored so a digit immediately after `q`
  or `m` blocked every backtrack — which is why `grep -m 1` was pinned and `grep -m1` was not. All five
  are one running pipeline in bash and all five reproduce the defect at exit 141.
- **The three fixes are additive, proven rather than argued.** Both scanner bodies — the one shipped at
  the branch point and the widened one — report **exactly 129** offending lines over the 93 tracked shell
  files at that revision, and the widened one reports **0** over all 94 at HEAD. So the widening adds
  five shapes and changes nothing about what it already caught.
- **The five join the fixture by name**, taking it to 19 must-flag and 7 must-not-flag shapes, and three
  mutations each name the shapes they lose: dropping the blank-line rule loses two, dropping the
  note-after-pipe strip loses one, narrowing the option class loses two.
- **The header claimed coverage the parser did not have.** It said "there is no state where the offender
  is both hidden and live" — and a note after an open pipe was exactly that state. That claim is gone,
  replaced by what the parser does plus an explicit list of what no source scan can see: a non-grep
  reader, a runtime-assembled pipeline, and three textual false-positive shapes that are all loud.
- **The comment skip's own false-positive class is now disclosed too.** Making the skip unconditional
  also joins a backslash-continued command to the next one when a comment terminates it. The additive
  proof showed nothing was removed; it says nothing about what was added on that side, and a reviewer
  said so.
- **Four counts in this record are corrected rather than defended**: the mechanical/judgment split is
  99/30 and not 122/7, find's status is asserted at 20 of 21 leak checks and not all 21, three of those
  21 sit in `if false` blocks that never execute, and 96 of 128 sites are heavy-only rather than 90. The
  delivered work was right in every case; the record described it wrongly.
- **Two findings outside this request are captured as REQ-600**: the prime `CLAUDE.md` calls the
  hard-won trap list still does not mention this class after two requests and 140 fixed sites, and one
  shipped ```bash block that agents copy carries the forbidden shape.

## Testing

**Three heavy lanes and the fast gate, each with its own exit line, because a lane that skips reports
success.** Two thirds of this change (96 of the 128 `_dev/tests` sites, **corrected after review from
90**) is heavy-only: 90 under `prescribed-shell-cases/`, reached only through
`staged-skills-contract.sh:189`, plus `staged-skills-contract.sh`'s own 4 and
`install-suite-behavior.sh`'s 2, both of which refuse to run outside the heavy tier. A fast run alone
would have been no evidence at all. The three lanes below were required and run, so no evidence was
missing — only the number was wrong.

- Fast gate — `Maintainer verification passed.`, exit 0, gate wall 99s at the builder's final head and
  79s on the merged tree, CLI module at 796 tests. The audit's own line:
  `quiet-grep pipeline audit passed (94 tracked shell files, 14 must-flag and 7 must-not-flag shapes).`
- `--heavy-lane staged-skills` — EXIT 0, wall 24s,
  `Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).`
  This is the lane that reaches all 90 prescribed-shell-case conversions.
- `--heavy-lane updater` — EXIT 0, wall 72s, `update-script behavior probes passed.` This is the lane
  that proves the 94-line deletion left the surviving probe intact.
- `--heavy-lane installer` — EXIT 0, wall 16s, `suite installer behavior probes passed.`

**Beyond the lanes:** every touched probe and case file was also run standalone and reported zero
failures, and the ShellCheck lane form — every tracked `*.sh` in one invocation from the repository root
— exits 0.

**The guard's own cost was measured, not estimated.** `time bash _dev/tests/quiet-grep-pipeline-audit.sh`
reports 0.196s at 93 files and 0.201s at 94, against a 30s per-file fast budget. The gate's launcher
reports it as 2s because it measures in whole seconds.

**The two already-broken guards were mutation-tested on the real files**, not only in the harness: a
planted leftover fires the replacement, and the shipped form on the same fixture reports zero failures.

### Remediation testing (after review)

**The five evasions were shown to be real defects before they were fixed.** Each was run as a shell
fragment with a 400,000-line producer whose first line carries the marker, under `pipefail`: a blank
line after the pipe, a blank line after a continuation, `grep -m1`, `grep -qm1`, and a note after the
pipe on the same line all exit **141** — the writer's SIGPIPE death reported as the reader's verdict,
where a correct verdict is 0. The control, a plain `grep -q`, exits 141 too.

**And they were shown invisible to the shipped guard**, in a scratch repository holding a verbatim copy
of the audit plus one file carrying all five shapes and two controls: the audit reported only the two
controls and passed on all five evasions.

**Additive proof, run at both ends.** Both scanner bodies over the 93 tracked shell files at the branch
point: `TOTAL 129` each. The widened body over all 94 files at HEAD: `TOTAL 0`. The audit's own line:
`quiet-grep pipeline audit passed (94 tracked shell files, 19 must-flag and 7 must-not-flag shapes).`

**Three mutations, each naming what it lost:**

- the blank-line rule removed → `a blank line between the pipe and the reader`, `a blank line between
  the backslash continuation and the reader`
- the note-after-pipe strip removed → `a note after the pipe on the same line, with the reader on the
  next`
- the option class narrowed back → `grep -m1, with no space before the count`, `grep -qm1, quiet bundled
  with a count`

**Lint and gate.** `bash -n` and `shellcheck --severity=warning` on the audit both exit 0. The fast gate
at the remediation revision prints `Maintainer verification passed.`, exit 0, gate wall 75s, with the
audit's own passing line inside it.

**The reviewers re-ran the original evidence rather than reading it**, including the three heavy lanes,
the two already-broken guards, and a direction check on the security-relevant site: with 200,000 lines
after a leaked `.env.local`, the old form of `protected-inventory.sh:19` reported zero failures — the
secret leaked and the check passed — and the new form fires.

## Review

**Overall: 86%** | 2026-09-06T08:30:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 82% |
| Code Quality | 85% |
| Test Adequacy | 84% |
| Scope | 92% |
| Risk | Low |
| Acceptance | Accept with remediation |

**Verdict: Accept with remediation.** The central question answers cleanly: **no conversion turned a
real check into a weaker one.** The synthesizer verified that four independent ways, all executed — a
mechanical producer/reader pairing across all 129 sites (106 pairs matched with zero variable
mismatches, identical grep flags and identical pattern bytes; the 5 unpaired only because their patterns
contain a literal `|`, each read by hand); a flag-and-pattern multiset diff per file, where the only
matchers that disappear anywhere are the leak checks, each replaced by an emptiness test plus a new
`find`-status assertion; an empty-input equivalence check running all 112 head matchers against a single
empty line, since `printf '%s' ""` emits zero bytes where `<<<""` emits one; and a direction check on
the security-relevant site.

What failed was the guard's own second requirement — "not defeatable by ordinary shell". Five more
spellings survived it. Both findings are closed in the remediation.

Where the reviewers disagreed, and what was picked:

- The runtime. Two reviewers and the record said 0.20s; one measured 0.32-0.40s. Settled by seven
  consecutive runs with timestamps taken outside the measured region: 0.192-0.203s. Picked the 0.20s.
- Which evasions defeat the guard. One reviewer named three, another named three different ones with one
  overlap. Settled by running each candidate as a real pipeline and then against the shipped audit; five
  survive in total, the union of both lists minus one that does not reproduce.
- The mechanical/judgment split. One reviewer said the record's justification is false for at least 33 of
  the 122; another said the real split is 99/30. Both describe the same defect; the measured numbers are
  99 and 30.
- Whether the reader-set limitation is disclosed. One reviewer said nowhere; it is in this record's
  Risks section verbatim and in the hand-back. Rejected.

**Important findings:**

- Five ordinary grep spellings walk past the repository-wide guard — a blank line after the pipe, a blank
  line after a continuation, a note after the pipe on the same line, `grep -m1` and `grep -qm1`. All five
  reproduce the defect at exit 141 and all five are silent to the audit. None exists in the tree today,
  so 129-to-0 stands; what was at risk is the guard's only job. — impact-rule-change → fixed in
  remediation
- The audit's own header stated two things about its coverage that are false, including "there is no
  state where the offender is both hidden and live" — which is exactly the state a note after an open
  pipe occupies. — impact-rule-change → fixed in remediation

**Minor findings:**

- The record's mechanical/judgment split is wrong by a factor of four on the exact axis the review tests:
  99 variable-replay producers and 30 real-command producers, not 122 and 7. Every one of the 30 did get
  its status assertion, so the work was right and the summary was not. — impact-negligible → corrected
- D-05 said find's status is asserted at all 21 leak checks; it is 20, and the Qualification in the same
  record already said twenty. — impact-negligible → corrected
- Three of the 21 sit in pre-existing `if false` blocks and never execute, so "every search root is
  created by the fixture" is proven for 18. In the hand-back, not in this record. — impact-negligible →
  corrected
- "90 of the 128 sites are heavy-only" is 96: the two contract files refuse to run outside the heavy tier
  too. Both lanes were required and run, so no evidence was missing. — impact-negligible → corrected
- The comment fix adds an undisclosed false-positive class: an unconditional comment skip joins a
  backslash-continued command to the next one when a comment terminates it. Loud, never silent, zero
  instances. — impact-negligible → disclosed in the audit's header
- One shipped ```bash block that agents copy carries the forbidden shape, so "prose is .md" is not the
  whole reason the walk's boundary is where it is. — impact-rule-change → REQ-600
- The prime `CLAUDE.md` names as the hard-won trap list still carries no mention of this class, after two
  requests and about 140 fixed sites. — impact-rule-change → REQ-600

**Findings raised and rejected:**

- "The must-not-flag fixture pins a helper-function reader as safe" — the fixture line is a
  file-reading grep, not a pipeline.
- "The reader-set limitation is admitted nowhere" — it is in this record's Risks section, naming the
  same five readers.
- "The release is not in the reviewed range" — the shipped line is in the range; the version bump and
  changelog belong to finalization, which is where every request in this batch puts them.
- "Three leak checks in dead code" as a standalone medium — correct and reproduced, but the `if false`
  blocks predate this change.

**Requirements checklist:**

- [x] No check under `_dev/tests/` decides on a quiet grep fed by a pipeline whose writer can die —
  delivered, 129 to 0 measured with the same instrument at both ends
- [x] The guard is repository-wide rather than per-file — delivered, walking the same enumeration the
  ShellCheck lane uses
- [ ] → [x] The guard is not defeatable by ordinary shell — **not delivered at review**, five spellings
  survived; **delivered in remediation**, with the widening proven additive
- [x] Each conversion preserves what its assertion measured — delivered, verified four ways
- [x] → A probe fixture carries every evasion spelling — delivered for the fourteen named shapes at
  review, nineteen after the remediation
- [x] The two already-broken guards fire on the input that used to pass silently — delivered
- [x] The release lands before the repository-wide guard, so the gate is never red between steps —
  delivered
- [x] Scope is exactly the declared write set — delivered, set difference empty in both directions
- [x] The guard fires on a real reintroduced offender, not only on its fixture — delivered

**Acceptance testing**

**Result: Partial at review, Pass after remediation.** Three reviewers re-ran the gate, the three heavy
lanes, the two already-broken guards and the additive proof themselves, and left the working tree clean.

**Follow-ups created:** 1 — REQ-600, for the prime and the one shipped block.

*Reviewed by review-work action*

## Lessons Learned

- **A guard written to catch five known evasions catches five known evasions.** REQ-593's scanner was
  walked past by five spellings a reviewer found in minutes; this request widened it and it was walked
  past by five more. The second five were not exotic — a blank line, a note after a pipe, a missing
  space before a count. The lesson is not "widen again": it is that a textual scanner's coverage is
  bounded by the author's imagination, and the only honest response is to say in the file what it cannot
  see. That list is now longer than the list of what it can.
- **`grep -m 1` was pinned and `grep -m1` was not, and the fixture said nothing.** The option class
  ended in `[A-Za-z]*`, so a digit right after `q` or `m` blocked every backtrack. A fixture that pins
  the spelling you wrote proves the scanner matches the spelling you wrote. Vary the whitespace, the
  bundling and the punctuation of every shape you pin, because that is exactly what the next author will
  vary without thinking about it.
- **"Additive" needs proving in one direction and disclosing in the other.** Running both scanner bodies
  over the branch-point tree and getting 129 each proves nothing was removed. It says nothing about
  false positives added, and a reviewer had to point that out. Two different claims, two different
  experiments.
- **A comment that overclaims survives longer than code that does.** "There is no state where the
  offender is both hidden and live" was written as a reassurance and was false for a shape three lines
  of awk away. Code gets exercised; a comment is believed until someone tests it, and the person most
  likely to believe it is the next author deciding not to look further.
- **Counting the easy category first makes the hard one look small.** The record said 122 mechanical and
  7 judgment; the measured split is 99 and 30. The work was right — every one of the 30 got its status
  assertion — but the summary described the change as four times more mechanical than it was, on the
  exact axis a reviewer would use to decide how hard to look.
- **The guard is not where authors read.** After two requests and about 140 fixed sites, the prime
  `CLAUDE.md` calls the hard-won trap list still does not mention this class. A guard catches the defect
  after it is written; a prime stops it being written. Both are needed, and only one of them existed.

## Orientation

The defect: a writer piped into an early-leaving reader under `set -o pipefail` reports the writer's
SIGPIPE death as the reader's verdict. Silent below roughly 36 KB of producer output, certain above
about 200 KB, and **wrong in both directions** — a positive matcher misses a pattern that is present, a
negative one fails to flag one it should. The negative ones are the dangerous half: a dead writer makes
`&& fail` pass.

**The guard is `_dev/tests/quiet-grep-pipeline-audit.sh`**, registered as a fast-tier probe in
`_dev/tests/contracts/probe-lanes.sh`. It walks `git ls-files -z -- '*.sh'` — the same enumeration the
ShellCheck lane uses at `maintainer-verify.sh:670` — and runs in 0.20s over 94 files. Its fixture pins
**19 must-flag and 7 must-not-flag shapes, each by name**, so a mutation says which shape it lost rather
than that a count moved.

**Read the header before changing the parser.** It carries the three things that make it work — logical
lines rather than physical, any pipe stage after the first, any early-leaving option — and, more
importantly, what it cannot see: a reader that is not grep (`rg -q`, `head`, `sed -n '1p;q'`,
`awk '/x/{exit}'`, `read`), a pipeline assembled at runtime, and three textual false-positive shapes
that are all loud rather than silent.

**Converting a site is not mechanical when the producer can fail.** 99 of the 129 replayed a variable
whose status was already asserted; 30 had a live producer status, and each of those got an explicit
assertion, because capturing output into a variable discards the status the pipeline used to carry. That
is the trap REQ-593 proved: a truncated archive still produced a listing whose first line carried the
marker, so a naive herestring passed where the pipeline had failed.

Two guards in this tree were **already broken** and are fixed: `staged-skills-contract.sh` and
`install-suite-behavior.sh` each had a `find … -print -quit | grep -q .` check that reported the legacy
runtime as retired whenever `find` exited non-zero. If you write another `find`-based leak check, assert
`find`'s status — 20 of the 21 here do, and the 21st is a readiness poll that says in place why it
differs.

Recorded and unfixed: the class is still missing from `_dev/primes/prime-shell-commands.md`, and one
shipped ```bash block at `skills/do-work-knowledge/actions/memory-reference.md:88` carries the shape.
Both are REQ-600.
