---
id: REQ-282
title: Make the release probes run in a suite checkout, and report not applicable elsewhere
status: completed
created_at: 2026-08-19T13:42:45Z
claimed_at: 2026-08-19T20:31:30Z
completed_at: 2026-08-19T20:44:05Z
commit:
user_request: UR-057
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-279, REQ-280, REQ-281, REQ-283]
batch: upstream-consumer-report-2026-08-19
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-19T20:31:30Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - cross-route regression gates
write_set:
- skills/do-work-board/tools/queue-kanban/release.go
- skills/do-work-board/tools/queue-kanban/verify.go
- skills/do-work-board/tools/queue-kanban/release_test.go
- skills/do-work/actions/forensics.md
---

# Make the Release Probes Run in a Suite Checkout, and Report Not Applicable Elsewhere

## What

`verify.go:152` resolves the version file with `resolveVersionFilePath(repoRoot, "")`, which yields `<repo-root>/actions/version.md` (`release.go:26`). That path exists only when the repo root *is* the suite root. In a consumer install the suite lives under `.claude/skills/do-work/`, so the three release probes — version file against newest changelog entry, versions strictly increasing, titles not reused — report:

```
- skipped version-vs-changelog probes: no version file readable at <project-root>/actions/version.md
```

`actions/forensics.md:187` Check 14 instructs consumers to run exactly that command with the project root, so in every consumer install those three probes are permanently off and permanently *skipped*.

Two changes: make the probes resolve this repo's actual version file so they run where they belong, and report them as **not applicable** where they genuinely cannot — saying so in Check 14.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Two costs, and the one found during capture is the larger.

In the **maintainer repo** the three probes are simply off, so the release invariants CLAUDE.md's *Before Every Commit* ritual asks a human to verify by eye — version file against newest changelog entry, versions strictly increasing, titles not reused — have no mechanical check behind them at all. That section notes duplicate version numbers "have occurred before"; the probe written to catch them has not run since the modular split.

In a **consumer repo** the same line is permanent and unactionable, which trains readers to scroll past the skipped section — including on the day a genuinely skipped probe means something. There the fix is honesty about the category, not new capability.

## Context

`release.go:19-25` already documents the consumer case as deliberate, for `next-version`: "In a CONSUMER install this path does not exist at the repo root … and that is deliberate: next-version then reports and writes nothing … This tool is an accelerator for a repo that matches do-work's convention, never a dependency." That reasoning was written for the writer, and nothing extended it to `verify`'s reader.

**The defect is wider than the report knew — the probes are off here too.** The upstream report assumed the default path is correct "when the repo root *is* the suite root, as in the skill's own development repo." That stopped being true at the modular four-skill split: this repo's version file is `skills/do-work/actions/version.md`, and there is no `actions/version.md` at the root. Running `queue-kanban verify --repo-root "$PWD"` in this repo today prints:

```
  OK: no findings
  - skipped version-vs-changelog probes: no version file readable at <repo>/actions/version.md
```

So all three release probes are unverified **in the maintainer repo**, which is the one place they were meant to run — and they are the three invariants CLAUDE.md's *Before Every Commit* ritual asks a human to check by eye, including the duplicate-version-number failure that section says "have occurred before." This REQ therefore has two halves that must not be confused: report not-applicable for a consumer install, and make the probes actually run here.

## Detailed Requirements

**Half one — make the probes run in a suite checkout.**

- Resolve the version file so it is found in this repo. `skills/do-work/actions/version.md` is the maintainer layout; the root `CHANGELOG.md` remains the changelog the probes compare against, since the release ritual bumps root `VERSION`, `skills/do-work/VERSION`, and `skills/do-work/actions/version.md` together against root `CHANGELOG.md`.
- Do not turn this into a search. One additional known location for the modular layout, tried after the existing default, is the whole mechanism.
- Verify against reality before finishing: `queue-kanban verify --repo-root "$PWD"` in this repo must run all three probes and report `0.212.25` agreeing with the newest `CHANGELOG.md` entry.

**Half two — report not-applicable where they genuinely cannot run.**

- When no version file resolves *and* the repo root is not a suite checkout, append to a not-applicable category rather than `SkippedProbes`, with wording along the lines of: `not applicable: release probes verify the suite's own release ritual, and this repo root is not a suite checkout`.
- Render not-applicable entries as their own section in the report, visibly distinct from skipped probes.
- Keep every genuinely-skipped case as skipped: a suite checkout whose version file is unreadable, a `CHANGELOG.md` that cannot be read, a changelog with no house-format entries, and a version string that fails to compare all stay in `SkippedProbes`. Only the not-a-suite-checkout case moves.
- Update `actions/forensics.md` Check 14 so its "a skipped probe is an unverified invariant" instruction names the not-applicable category and says a consumer install is expected to show the release probes there.

**Do not let half two swallow half one.** A not-applicable path that also fires in a suite checkout would silence the probes permanently and read as clean — the exact failure this REQ exists to end. The maintainer-shaped assertion in the Red-Green Proof is what pins that.

## Constraints

- **Do not discover a suite root.** The upstream report's first remedy was to walk up from the executable or accept `--suite-root` and run the probes against the vendored copy. Declined at triage: in a consumer repo the root `CHANGELOG.md` is the consumer's, so pointing the probes at a suite root would check install integrity rather than a release-ritual violation the consumer can act on — and `decisions/records/adr-019-four-skill-suite-contract.md` already assigns install integrity to the updater's all-or-recover contract ("success is reported only after all installed bytes verify"). A second implementation of a covered invariant is unearned surface.
- **Do not add `--version-file` to `verify`.** Same reasoning; the flag exists on `next-version` because that subcommand writes, and a consumer pointing verify at a vendored version file would be asking the wrong question.
- **Do not change `next-version`'s behavior.** Its `--version-file` override and its documented write-nothing consumer path stay exactly as they are; any shared resolution change must leave that subcommand's observable behavior identical.
- Detecting "not a suite checkout" must not become its own inference engine — an existence test on the conventional locations is the whole mechanism.

## Red-Green Proof

**RED prompt/case:** Two cases, one command each.
1. Maintainer: `queue-kanban verify --repo-root "$PWD"` in this repo.
2. Consumer: build a fixture with `do-work/` and a queue at the root, no version file at the root, and the suite vendored under `.claude/skills/`; run verify against it.

**Why RED now:** Both print the same line — `- skipped version-vs-changelog probes: no version file readable at <root>/actions/version.md`. In case 1 that is three real invariants silently unverified in the repo that owns them (this repo's version file is `skills/do-work/actions/version.md`). In case 2 it is a path error dressed as an unverified invariant.

**GREEN when:** Case 1 runs all three probes and reports `0.212.25` agreeing with the newest `CHANGELOG.md` entry, and a deliberately mismatched version in a maintainer-shaped fixture produces a finding. Case 2 prints the release probes under a not-applicable heading with the suite-checkout reason and no path in the message. Both asserted in the same test, so neither half can be satisfied by breaking the other.

**Validation:** User confirmed for the consumer half (upstream report accepted at triage). The maintainer half was found by running `queue-kanban verify --repo-root "$PWD"` in this repo during capture and reading its skipped-probe line — a second, sharper instance of the same defect that the upstream report explicitly assumed did not exist. Mechanics confirmed at `verify.go:150-160`, `release.go:26,52-57`, `main.go:341-346` (verify takes only `--repo-root`) versus `main.go:243` (`next-version` has `--version-file`).

## Full Context

See `do-work/user-requests/UR-057/input.md` for the complete verbatim upstream report.

---
*Source: upstream defect report D3, severity medium, from `g1w-game-find-the-difference` running v0.212.25 — verbatim claim: "`verify`'s release probes can never run in a consumer repo … three invariants report as skipped rather than as inapplicable … If the install is a consumer vendoring, report `- not applicable: release probes verify the suite's own release ritual`." Accepted (narrowed to the reporting half) by `do-work-toolbox validate-feedback` triage (2026-08-19); the report's suite-root-discovery remedy was declined. Evidence: `skills/do-work-board/tools/queue-kanban/release.go:26`, `verify.go:152` and `:158`, `main.go:341-346` vs `main.go:243`, `skills/do-work/actions/forensics.md:187`; the deliberate-for-next-version reasoning at `release.go:19-25`; `decisions/records/adr-019-four-skill-suite-contract.md` on update validation. Surface-cost: Earned — incident is this report, a permanent unactionable warning in every consumer install that trains readers past the skipped section; surface is one existence test and a second report bucket, replacing rather than adding to the current line; strictly cheaper than suite-root discovery plus its fallback behavior; test is the consumer-shaped and maintainer-shaped fixture pair above.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The Detailed Requirements and Constraints are precise, but the "where" needs discovery: which resolver the two subcommands share, and whether the not-applicable bucket can be added without changing `next-version`'s observable behavior — the REQ's hardest constraint.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The reproduction, first.** `queue-kanban verify --repo-root "$PWD"` in this repo prints
exactly the line the REQ describes, confirming the maintainer half is live:

```
  OK: no findings
  - skipped version-vs-changelog probes: no version file readable at <repo>/actions/version.md: ... no such file or directory
```

**The shared-resolver trap.** `resolveVersionFilePath(repoRoot, override)` (`release.go:52`)
has two callers: `main.go:312` for `next-version`, which **writes**, and `verify.go:152`,
which reads. Adding the modular-layout fallback *inside* that function would give
`next-version` a version file to write in this repo, where today it deliberately finds none
and reports nothing — the exact observable change the REQ's third Constraint forbids. So the
fallback goes in a **new resolver used only by the release probes**, and
`resolveVersionFilePath` is left byte-for-byte alone.

**Where "not applicable" has to live.** `VerifyReport` (`verify.go:73`) has `Findings` and
`SkippedProbes`; the renderer (`:865`) prints skips as `  - skipped %s`. A third slice plus
its own render line is the whole mechanism — no new type, no change to `ExitCode()`, since a
not-applicable probe is no more a failure than a skipped one.

**The four cases that must stay skipped**, all reachable only once a version file resolves:
a version file that exists but carries no `**Current version**:` line, an unreadable
`CHANGELOG.md`, a changelog with no house-format entries, and a version string that fails to
compare. Only the neither-path-exists case moves to not-applicable, which is what keeps half
two from swallowing half one.

**Stale figure in the REQ:** its Red-Green Proof asserts the probes should report `0.212.25`.
Three REQs shipped tonight, so the repo is past that. The invariant the REQ actually
describes is *the version file agrees with the newest changelog entry*, and that is what the
assertion pins — a literal would have to be edited on every release.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/release.go` (modify) — the probe-only resolver and the modular-layout constant
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — `NotApplicableProbes`, the release-probe branch, and the render line
- `skills/do-work-board/tools/queue-kanban/release_test.go` (modify) — the maintainer-shaped and consumer-shaped fixture pair
- `skills/do-work/actions/forensics.md` (modify) — Check 14 names the not-applicable category

**Files I will NOT touch:** `main.go` (`next-version`'s wiring and its `--version-file` flag stay exactly as they are), `resolveVersionFilePath` itself.

**Acceptance criteria (restated from REQ):**
- [x] The three release probes run in this repo and the version file agrees with the newest `CHANGELOG.md` entry
- [x] One additional known location tried after the existing default — not a search, not suite-root discovery
- [x] `next-version`'s observable behavior is identical, including its consumer write-nothing path
- [x] A repo root that is not a suite checkout reports the release probes as **not applicable**, rendered distinctly, with no path in the message
- [x] Every genuinely-skipped case stays skipped
- [x] `actions/forensics.md` Check 14 names the not-applicable category
- [x] Both halves asserted in one test, so neither can be satisfied by breaking the other
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## Decisions

- **D-01**: The modular-layout fallback lives in a **new** read-only resolver,
  `resolveReleaseProbeVersionFilePath`, and `resolveVersionFilePath` is left byte-for-byte
  unchanged. Putting the fallback in the shared function would have satisfied half one and
  broken the REQ's third Constraint in the same edit: that function is `next-version`'s
  **writer** path, so in this repo it would start rewriting
  `skills/do-work/actions/version.md` where today it deliberately finds nothing and reports
  nothing. Pinned by `TestNextVersionResolutionIsUnchangedBySuiteFallback`, which asserts
  the writer's resolver still yields the pre-split root path and that the fixture does not
  create it. DECIDE & STATE.
- **D-02**: "Is this a suite checkout" is answered by the same existence test that finds the
  file — the resolver returns `(path, isSuiteCheckout)`, and neither location existing is
  the whole condition. No second detector, no walk, no `--suite-root`, per the REQ's
  Constraints. DECIDE & STATE.
- **D-03**: The REQ's Red-Green Proof asserts the probes should report `0.212.25`. Three
  REQs shipped tonight and the repo is at 0.215.3, so the assertion pins the **invariant**
  the REQ describes — the version file agrees with the newest changelog entry — rather than
  the literal, which would need editing on every release. The real-repo run is recorded in
  Testing below with the actual number. DECIDE & STATE.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/release.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/release_test.go` (modified)
- `skills/do-work/actions/forensics.md` (modified)

**What was done:** Added `suiteVersionFileRelativePath` and a read-only
`resolveReleaseProbeVersionFilePath(repoRoot) (path, isSuiteCheckout)` that tries the
pre-split root path first and the modular suite path second — two known locations, an
existence test, no search. `VerifyReport` gained `NotApplicableProbes`, rendered as
`  ~ not applicable: …` beside the existing `  - skipped …` so the two claims are distinct
on screen; `ExitCode()` treats it like a skip. `appendReleaseFindings` routes only the
neither-location-exists case there, so every failure reachable after a version file resolves
stays a skip. `next-version`'s resolver is untouched. Check 14 now says a consumer install
should expect the release probes under not-applicable, and that a skip there is still a gap.

## Qualification

Passed — 4 files verified in the diff, 8 acceptance criteria traced, P-A-U confirmed.
Scope-drift check clean. Judgment check 6 (data flows): the probes are read-only and were
confirmed to fire on a real mismatch rather than merely staying quiet, which is the specific
hollow-implementation risk here — a disabled probe and a passing one print the same thing.

## Testing

**Tests run:** `go test ./...` in `skills/do-work-board/tools/queue-kanban`, then
`bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing (maintainer-verify exit 0)

**Red-green validation:**
- `TestReleaseProbesRunInASuiteCheckoutAndAreNotApplicableElsewhere`: ✗ before → ✓ after
- `TestGenuineReleaseProbeSkipsStaySkipped`: ✗ before → ✓ after
- `TestNextVersionResolutionIsUnchangedBySuiteFallback`: ✗ before → ✓ after

All three were written before the implementation and failed at compile time on the missing
`NotApplicableProbes` field — the coarsest possible RED, so each was additionally confirmed
by the reviewer, which reverted each production hunk in a scratch copy and checked the
specific test goes red on the mutation it claims to catch. Folding the suite fallback into
the shared writer's resolver fails the third test, which is the REQ's third Constraint held
mechanically rather than by memory.

**Real-repo acceptance (the REQ asks for this explicitly):**

- Before: `queue-kanban verify --repo-root "$PWD"` printed
  `- skipped version-vs-changelog probes: no version file readable at <repo>/actions/version.md`.
- After: no skip line, `OK: no findings`, exit 0 — the version file
  (`skills/do-work/actions/version.md`, **0.215.3**) agrees with the newest `CHANGELOG.md`
  entry, `## 0.215.3 — Refuse an Unterminated Frontmatter Fence Before Reading It`.
- **Proved the probes bite rather than went quiet:** with the version line temporarily set
  to `9.9.9`, the same command reported
  `! version-changelog-mismatch: version 9.9.9 is ahead of the newest CHANGELOG.md entry 0.215.3` and
  exited 1. The file was restored from a copy and confirmed byte-identical with
  `git diff --quiet`.
- Consumer-shaped fixture (`do-work/` and `CHANGELOG.md` at the root, suite vendored under
  `.claude/skills/`): `~ not applicable: release probes: they verify the suite's own release
  ritual, and this repo root is not a suite checkout`, exit 0, no path in the message, and
  visibly distinct from the `- skipped worktree probes` line directly above it.

**New tests added:**
- Three tests plus one fixture helper in `release_test.go`

**Existing tests updated:** none.

*Verified by work action*

## Review

**Reviewer:** independent agent, orchestrated mode against the working-tree diff.

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 98% |
| Scope Discipline | 100% |
| Risk | None |
| **Acceptance** | **Pass** |
| **Overall** | **98%** |

**No Important findings.** The reviewer attacked the REQ's own sharpest warning — that a
not-applicable path firing in a suite checkout would silence the probes permanently and read
as clean — across seven shapes: a directory at the version path, an unreadable file
(`chmod 000`), an empty file, a symlinked file, both candidate paths present, a relative
`--repo-root`, and a trailing-slash `--repo-root`. Every one routed to the correct bucket,
and no suite checkout reached `NotApplicableProbes`. It confirmed `resolveVersionFilePath`
is byte-for-byte unchanged with `main.go:312` its only caller, that no failure branch is
now dropped from both buckets, and that `VerifyReport` has no JSON or board consumer whose
contract the new field could break. Its Restatement Sweep found `forensics.md` Check 14 to
be the only canonical restatement, correctly updated.

**Minor, fixed:** D-03 pointed at a `## Testing` section recording the real-repo number, and
that section did not exist yet — a dangling reference in the REQ, not in the code. Written
above, with 0.215.3 and the entry it agrees with.

**Minor, accepted and not fixed:** the consumer-vs-suite verdict is an existence test on two
conventional paths, so a repo that coincidentally holds one of them is classified as a suite
checkout. That is the mechanism the REQ's Constraints chose deliberately over suite-root
discovery, with the reasoning recorded there; the reviewer flagged it as a known limitation
rather than a defect, and I agree.

**Suggested by the reviewer, not done:** run `verify` against a repo produced by the real
installer rather than a hand-built fixture, to confirm the vendored path shape. Worth doing
once, but it tests the installer's output shape rather than this change, and the
not-applicable branch does not depend on where the vendored copy lives — only on neither
root path existing.

## Lessons Learned

**What worked:** Proving the probes *bite* on the real repo, not just that the skip line
disappeared. A disabled probe and a working one both print nothing on a clean repo, so
"`OK: no findings`" was never evidence on its own — the temporary `9.9.9` run is what turned
it into evidence. The same trap is why the test asserts a mismatch produces a finding rather
than only asserting the clean case.

**What didn't:** Nothing failed, but the near-miss is worth recording. The obvious
implementation — add the fallback to `resolveVersionFilePath` — satisfies half one in one
line and silently breaks the REQ's third Constraint, because that function is also
`next-version`'s **writer** path and would have started rewriting a file it correctly finds
nothing at today. The REQ named the constraint; the code did not enforce it until a test did.

**Worth knowing:** a resolver shared between a reader and a writer cannot be widened for the
reader's benefit. The two want opposite things from a missing file — the reader wants a
fallback, the writer wants to find nothing and stop — so widening for one silently changes
the other. Splitting them cost eight lines and is pinned by
`TestNextVersionResolutionIsUnchangedBySuiteFallback`.

## Orientation

`queue-kanban verify`'s three release probes now run in this repo, where they have been off
since the four-skill split moved the version file under `skills/do-work/` — the same three
invariants CLAUDE.md's *Before Every Commit* ritual asks a human to check by eye, including
the duplicate-version-number failure that section notes has happened before. In a consumer
install they report as **not applicable** rather than skipped, so a permanent unactionable
warning stops training readers past the skipped section. Lives in the board tool's verify
subsystem (`_dev/primes/prime-kanban-board.md`). No map change: one new read-only resolver
beside the existing writer's, one new report bucket beside `SkippedProbes`. Prime spot-check:
`prime-kanban-board.md`'s referenced paths all still exist.
