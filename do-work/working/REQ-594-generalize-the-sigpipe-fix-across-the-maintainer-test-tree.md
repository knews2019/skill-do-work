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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
