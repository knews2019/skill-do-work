---
id: REQ-282
title: Make the release probes run in a suite checkout, and report not applicable elsewhere
status: pending
created_at: 2026-08-19T13:42:45Z
user_request: UR-057
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-279, REQ-280, REQ-281, REQ-283]
batch: upstream-consumer-report-2026-08-19
write_set:
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
