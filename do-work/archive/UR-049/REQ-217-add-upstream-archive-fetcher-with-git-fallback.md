---
id: REQ-217
title: Add the upstream archive fetcher with a git fallback
status: completed
completed_at: 2026-08-17T19:13:20Z
commit: 0e8cf0d
claimed_at: 2026-08-17T19:05:40Z
created_at: 2026-08-17T17:16:28Z
route: B
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-17T19:06:20Z
  basis:
    - Route B
    - 9-file write set
    - 1 new files
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
    - full-suite verification
user_request: UR-049
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-216]
maintenance: false
related: [REQ-216, REQ-218]
batch: resilient-upstream-fetch
write_set:
  - tools/fetch-upstream-archive.sh
  - _dev/tests/staged-skills-contract.sh
  - skills/do-work/tools/fetch-upstream-archive.sh
  - tools/install-do-work-suite.sh
  - skills/do-work/tools/install-do-work-suite.sh
  - skills/do-work/tools/do-work-update.sh
  - skills/do-work/actions/version.md
  - README.md
  - _dev/tests/update-script-behavior.sh
---

# Add the Upstream Archive Fetcher With a Git Fallback

## What

Add `tools/fetch-upstream-archive.sh` (plus its `skills/do-work/tools/` mirror) that fetches this repo's
tarball with an anonymous HTTP route and a **`git archive` fallback**, then route `do-work-update.sh` and the
installer's no-`--archive` branch through it. Add `DO_WORK_UPSTREAM_URL` as a supported override so a blocked
consumer has an escape hatch that is not "edit a vendored file".

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One fetcher at `tools/`, mirrored, taking `<archive-target-path> <tarball-url> [repo-url]`. Route 1 delegates to REQ-216's `atomic-download.sh`, located by probing both mirror-relative depths since this file lives in two layouts. Route 2 shallow-clones and repacks with `git archive --prefix`, deriving the repo URL and the prefix from the tarball URL's `/archive/refs/heads/<branch>.tar.gz` shape. Verify each produced archive with `tar tzf` before publishing, so an unreadable 200 falls through to git rather than being accepted. Report the winning route on stdout; name both outcomes plus the `DO_WORK_UPSTREAM_URL` escape hatch on stderr when both fail.
- [x] **[APPLY]:** Nine files. One scope extension recorded as D-01 below.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file — `tools/fetch-upstream-archive.sh` and its mirror (new, byte-identical); `tools/install-do-work-suite.sh` and its mirror (env default, fetcher sibling, fetching branch, bootstrap heredoc retry flags); `skills/do-work/tools/do-work-update.sh` (env default, fetcher sibling and presence check, fetch call); `README.md` (bootstrap twin, byte-identical to the heredoc); `skills/do-work/actions/version.md` (escape-hatch guidance); `_dev/tests/update-script-behavior.sh` (fixture helpers sourced, fetcher + primitive installed into both fixture trees, three new cases plus two caller-delegation checks); `_dev/tests/staged-skills-contract.sh` (mirror enumeration generalized). ShellCheck via `maintainer-verify.sh` clean on the new script and every modified one. No debug artifacts, no leftover fixtures.

## Why (if provided)

Feedback-brief finding **F4**: the reported 429 persisted across three attempts spanning ~2 minutes, so
retry alone would not have unblocked the consumer. `git clone` was the only probe that succeeded. The git
transport sits behind a different limiter than `codeload.github.com`, which makes the git fallback the
load-bearing fix — retry (REQ-216) is correct for genuine blips but does not close the class.

## Detailed Requirements

**Two routes, in order:**

1. The anonymous `https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz`, fetched through
   REQ-216's hardened `skills/do-work/scripts/atomic-download.sh` so it inherits the retry.
2. Shallow clone and repack:
   ```bash
   git clone --depth 1 --quiet "$upstream_repo_url" "$clone_directory"
   git -C "$clone_directory" archive --format=tar.gz \
     --prefix=skill-do-work-main/ HEAD > "$archive_target_path"
   ```
   `--prefix` preserves the single-top-level-directory layout that both callers' existing
   `tar --strip-components=1` depends on.

- **Report which route succeeded** so a consumer can tell a clean fetch from a degraded one. On total
  failure, name both route outcomes rather than only the last.
- Both callers default `upstream_url` from **`DO_WORK_UPSTREAM_URL`**, matching the existing
  `DO_WORK_INSTALL_CANCEL_EXIT_STATUS` idiom at `install-do-work-suite.sh:7`.
- Replace the bare `curl` at `skills/do-work/tools/do-work-update.sh:72` and at
  `install-do-work-suite.sh:173` (the `--archive`-absent branch) with a call to the new fetcher.
- The `BOOTSTRAP` heredoc at `install-do-work-suite.sh:18-32` and its `README.md:16` twin get the
  **`--retry` flags only** — nothing is installed yet at that point, so the inline `curl` there is correct
  layering, not duplication. No fallback chain, no fetcher call. The two blocks are already forced
  byte-identical by `_dev/tests/install-suite-behavior.sh:185-195`, so edit both.
- Update `skills/do-work/actions/version.md`'s user-facing update text to name `DO_WORK_UPSTREAM_URL` as the
  escape hatch, so a blocked consumer learns the way out from the error itself. Keep the existing rule at
  `version.md:48` intact — the engine remains the authoritative safety boundary and the action still must not
  hand-roll a second download.

## Constraints

- **`git archive`, never a worktree copy.** This is the whole of finding F3 and the single most important
  constraint in the batch. `cp -R`, `rsync`, and `tar` of a clone do **not** honor `export-ignore`;
  `git archive` does. Measured on this tree: 734 tracked files vs **294** archive entries. A worktree-copy
  implementation would ship `do-work/`, `kb/`, `_dev/`, `decisions/`, `ai-reports/`, and `CLAUDE.md` into
  consumer installs — and `CLAUDE.md` in particular is called out in `.gitattributes:35-43` as actively
  harmful downstream. **Nothing catches this:** `tools/validate-suite-manifest.sh --root <leaky tree>` was
  verified to print `suite manifest valid: v0.199.2 (4 modules)` and exit `0`. RED test 3 is the only guard.
- **The fetcher lives in `tools/`, not `scripts/`** — see Context below.
- No change to the installer's transaction, confirmation prompt, managed-section handling, or `--archive`
  contract. The reported workaround proved that path is sound; the fetch is the whole bug.
- No new runtime dependency. Git is already a hard prerequisite — `do-work-update.sh:57-58` refuses to run
  outside a Git worktree.
- Both `tools/` mirrors must move together; `_dev/tests/staged-skills-contract.sh:809` enforces byte-identity.

## Context

**Correction to the brief's C2 placement.** C2 proposed `skills/do-work/scripts/fetch-upstream-archive.sh`.
That does not work: `install-do-work-suite.sh` resolves its siblings via `$script_dir/`
(`tools/install-do-work-suite.sh:74-76`, alongside `validate-suite-manifest.sh` and
`replace-text-section.sh`) and is byte-mirrored into two locations with different relative layouts — a helper
under `skills/do-work/scripts/` is `../scripts/` from the installed copy and `../skills/do-work/scripts/`
from the archive-root copy. Putting the fetcher at `tools/fetch-upstream-archive.sh`, mirrored the same way
as its siblings, makes `$script_dir/` correct from both.

**Correction to the brief's C2 route list.** The brief's route 1
(`GH_TOKEN` → `api.github.com/repos/knews2019/skill-do-work/tarball/main`) is **dropped**. Its pre-signed
redirect still terminates at codeload, so the claim that it escapes the limiter is asserted rather than
demonstrated, and it is the one route that cannot be exercised offline in a fixture test. Under the
earned-defense rubric it does not clear the bar. The generic `GH_TOKEN`/`GITHUB_TOKEN` header still lands in
`atomic-download.sh` via REQ-216, where it is cheap and benefits every caller — so a consumer with a token
configured still gets authenticated requests on route 1.

## Dependencies

`depends_on: [REQ-216]` — route 1 fetches through the hardened primitive, so the retry behavior must land
first. REQ-218's ratchet only turns GREEN once this REQ removes the bare `curl` from both `tools/` scripts.

## Builder Guidance

Certainty: **Firm on the constraints, latitude on the shell.** The two routes, the `git archive` requirement,
the file placement, and the reporting behavior are fixed. How the fetcher structures its argument handling and
route reporting is yours, subject to `_dev/primes/prime-shell-commands.md`.

## Red-Green Proof

**RED prompt/case:** Three fixture cases in `_dev/tests/update-script-behavior.sh`, all against a local
fixture repo standing in for the upstream remote:
1. *Sustained 429 falls back to git* — a fake `curl` on `PATH` that always returns 429; assert a valid
   archive is still produced.
2. *The fallback honors `export-ignore`* — give the fixture repo a `.gitattributes` containing
   `/private-path export-ignore`, force the git route, and assert `private-path` is **absent** from the
   produced tarball.
3. *Target preserved on total failure* — every route fails; assert the pre-existing target file is unchanged
   and no `.download.*` scratch leaked.

**Why RED now:** There is no fetcher and no fallback — `do-work-update.sh:72` is a single `curl` that exits
non-zero on the first 429, which is exactly the reported incident. Case 2 is the one that catches a
`cp -R`-based implementation: such an implementation passes cases 1 and 3 and fails only this one.
**GREEN when:** All three cases pass, and `do-work update` completes via the git route when the HTTP route is
unreachable, reporting which route it used.
**Validation:** User confirmed (capture approved the three-REQ split, the dropped route 1, and the `tools/`
placement).

## Finding-Closure

Origin: external feedback brief, triaged via `do-work-toolbox validate-feedback` — findings F3 and F4 and
proposals C2 (adjusted) and C3, verdict **Accept**. Surface-cost: **earned** — the incident is named and
reproducible, the added surface is one fallback branch plus one env override, all three RED cases replay it,
and the one unearned piece (route 1) was cut at capture. Named regression tests: the three cases above in
`_dev/tests/update-script-behavior.sh`, with case 2 as the F3 guard.

## Full Context

See `do-work/user-requests/UR-049/input.md` for the complete verbatim brief.

## Triage

**Route:** B — Explore then Build

**Reasoning:** The routes, the placement, and the constraints are fixed by the REQ, but both callers' sibling-resolution shape, the fixture harness's script-installation points, and the mirror-contract enumeration all had to be read before the fetcher could be wired without breaking them.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route B. The REQ carries the design; exploration supplied the wiring.

## Exploration

**Key files:**
- `skills/do-work/tools/do-work-update.sh` — resolves siblings from `$script_dir`, already refuses to run outside a Git worktree (so git is a settled prerequisite).
- `tools/install-do-work-suite.sh` (+ mirror) — same `$script_dir` idiom for `validate-suite-manifest.sh` and `replace-text-section.sh`; its `--archive`-absent branch is the only fetch site.
- `_dev/tests/update-script-behavior.sh` — builds two fixture trees (`build_installed_project`, `build_suite_tree`), each copying an explicit list of scripts. **It did not source `_dev/tests/fixture-repo.sh`**, so the repo-fixture helpers were unavailable to new cases.
- `_dev/tests/staged-skills-contract.sh` — enforced mirror byte-identity against a hand-written list of three tool names.

**Concerns found:**
- The fetcher is mirrored into two layouts with different depths, so it cannot reach `atomic-download.sh` by one fixed relative path — it has to probe.
- Requiring the fetcher beside the installer unconditionally would break `--archive` installs, including the canonical bootstrap and any install from an archive that predates the fetcher.

## Decisions

- **D-01: extend the scope by one file — `_dev/tests/staged-skills-contract.sh`.** The mirror byte-identity contract enumerated three tool names by hand, so a fourth mirrored tool would have shipped unguarded. Rather than adding a fourth line to the list, the loop now derives the mirrored set from `tools/*.sh` present in both trees. This is the repo's own **Closed Enumerations Go Stale** rule applied to the enumeration that would otherwise have gone stale on this very commit. Recorded here rather than written silently, per the write-boundary rule; the file is added to `write_set`. ESCALATE-grade by reach but not by risk. **Value:** every future mirrored tool is covered automatically. **Risk:** the loop now depends on `tools/*.sh` naming, so a mirrored tool with a different extension would be missed (low, reversible).
- **D-02: require the fetcher only on the branch that fetches.** An unconditional sibling check next to the validator and replacer looked symmetrical but broke the `--archive` path, which needs no fetcher at all. Reasoning: a hard requirement should sit where the dependency is real; putting it at the top would make every supplied-archive install depend on a file that install never uses. DECIDE & STATE — caught by `_dev/tests/install-suite-behavior.sh` before it could ship.

## Scope

**Files I will touch:**
- `tools/fetch-upstream-archive.sh` (new)
- `skills/do-work/tools/fetch-upstream-archive.sh` (new)
- `tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/do-work-update.sh` (modified)
- `skills/do-work/actions/version.md` (modified)
- `README.md` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified — added by D-01)

**Acceptance criteria (restated from the REQ):**
1. Two routes in order: hardened HTTP through `atomic-download.sh`, then shallow clone repacked with `git archive`.
2. `git archive`, never a worktree copy — `export-ignore` must be honored.
3. The succeeding route is reported; total failure names both route outcomes.
4. Both callers default `upstream_url` from `DO_WORK_UPSTREAM_URL` and fetch through the shared fetcher.
5. The bootstrap heredoc and its README twin get retry flags only, and stay byte-identical.
6. `version.md` names `DO_WORK_UPSTREAM_URL` as the escape hatch without hand-rolling a second download.

## Pre-Flight

- Working tree clean outside `do-work/`.
- Baselines passing before the change: `update-script-behavior.sh`, `install-suite-behavior.sh`, `staged-skills-contract.sh`.

## Implementation Summary

**Files changed:**
- `tools/fetch-upstream-archive.sh` (new)
- `skills/do-work/tools/fetch-upstream-archive.sh` (new)
- `tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/do-work-update.sh` (modified)
- `skills/do-work/actions/version.md` (modified)
- `README.md` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)

**What was done:** Added the two-route fetcher and routed both no-archive fetch sites through it.

*The fetcher.* `fetch-upstream-archive.sh <archive-target-path> <tarball-url> [repo-url]`. Route 1 delegates to `atomic-download.sh`, located by probing `../scripts/` then `../skills/do-work/scripts/` — the two depths this file's mirrors sit at. Route 2 shallow-clones with `GIT_TERMINAL_PROMPT=0` and repacks with `git archive --format=tar.gz --prefix=<repo>-<branch>/`, which preserves the single-top-level-directory layout both callers' `tar --strip-components=1` depends on. When no repo URL is supplied it is derived from the tarball URL's `/archive/refs/heads/<branch>.tar.gz` shape, which is also where the prefix comes from. **Every produced archive is verified with `tar tzf` before it is accepted**, so an unreadable 200 falls through to git rather than being published as success. The winning route is printed on stdout; on total failure both route outcomes and the `DO_WORK_UPSTREAM_URL` escape hatch go to stderr and the exit is non-zero. The pre-existing target is never touched until a route has produced a verified archive.

*The callers.* Both now default `upstream_url` from `DO_WORK_UPSTREAM_URL` and call the fetcher instead of a bare `curl`. `do-work-update.sh` checks for the fetcher beside it like its other siblings; the installer checks only inside its fetching branch (D-02). The `BOOTSTRAP` heredoc and its `README.md` twin get the retry flags alone — nothing is installed at that point, so the inline `curl` there is correct layering — and remain byte-identical.

*The guidance.* `version.md` now explains that the engine tries two routes, that its failure message names both, and that `DO_WORK_UPSTREAM_URL` is the supported way out. It still forbids a second download; the engine remains the authoritative safety boundary.

**Tests touched:** `_dev/tests/update-script-behavior.sh` now sources `_dev/tests/fixture-repo.sh` (it never did), installs the fetcher and `atomic-download.sh` into both fixture trees, and carries the REQ's three cases plus two caller-delegation checks. `_dev/tests/staged-skills-contract.sh`'s mirror enumeration is derived rather than hand-listed.

## Qualification

Passed — 9 files verified, 6 requirements traced, no debug artifacts.

- Every file is in `write_set` (one added by the recorded D-01); nothing undeclared was touched.
- Substantive: a 100-line new script, two rewired callers, and three fixture cases — no reformatting.
- Requirements traced: two ordered routes (1); `git archive` with the RED case that catches a copy-based implementation (2); route reporting and dual-outcome failure text (3); env default and delegation asserted for both callers (4); heredoc/README retry flags with byte-identity still enforced by `install-suite-behavior.sh` (5); `version.md` escape hatch (6).
- Flowing: both routes are reachable and both were exercised — the HTTP route by the pre-existing updater probes, the git route by the new cases.

## Testing

**Tests run:** `bash _dev/tests/update-script-behavior.sh`, `bash _dev/tests/install-suite-behavior.sh`, `bash _dev/tests/staged-skills-contract.sh`, `bash _dev/tests/maintainer-verify.sh`

**Result:** ✓ all pass; maintainer-verify exit 0 with zero FAIL lines.

**Red-green validation:** two independent RED stages.

1. ✗ **No fetcher, callers on bare `curl`.** With the new cases in place and the fetcher removed and both callers reverted, the suite exits 1: `tools/fetch-upstream-archive.sh must exist and be executable`, plus `does not delegate its fetch to the shared fetcher` and `does not honor DO_WORK_UPSTREAM_URL` for each caller.
2. ✗ **A worktree-copy implementation (the F3 guard).** Replacing only the `git archive` call with `tar czf --exclude .git` of the clone leaves cases 1 and 3 passing and fails exactly one assertion: `the git route shipped an export-ignore path into the archive`. This is the constraint the REQ called the single most important in the batch, and it is now the only thing standing between a plausible implementation and shipping the maintainer tree into consumer installs.

→ ✓ **GREEN** — with the real implementation, a host that answers 429 forever still produces a valid archive via git, that archive contains `VERSION` and excludes the `export-ignore`d path under a single top-level directory, the winning route is named in the report, and a total failure leaves the pre-existing target byte-for-byte with no `.download.*` or `.fetching.*` scratch and a message naming both routes and the escape hatch.

**Existing tests updated (cross-REQ impact):** `_dev/tests/update-script-behavior.sh`'s two fixture builders now also install `fetch-upstream-archive.sh` and `atomic-download.sh`. This is required, not incidental: the updater's fetch now runs through both, so a fixture without them is not a faithful install. All pre-existing probes — including `expected exactly one archive download` — pass unchanged, which is the evidence that the HTTP route still wins when it can. `_dev/tests/staged-skills-contract.sh`'s mirror loop is generalized under D-01; it still fails on any divergent mirror, and now covers the new one too.

*Verified by work action*

## Review

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All six delivered; the dropped route 1 stayed dropped |
| Code Quality | 92% | The two-depth probe for the primitive is the price of the mirror; it is commented at the point it matters |
| Test Adequacy | 98% | The F3 guard was proven by building the wrong implementation and watching exactly one assertion fire |
| Scope | 95% | One file added mid-build, recorded as D-01 before it was written |
| Risk | Low | A new network route and a `git clone` of a derived URL — both fail closed, and the target is untouched until an archive verifies |
| Acceptance | Pass | Two RED stages, full baseline green |

**Verdict: Approve** — the fallback is the load-bearing half of this batch, and the one constraint that could have quietly ruined consumer installs is now the only test that a copy-based implementation fails.

### Findings

**Important:**
- `_dev/tests/update-script-behavior.sh` never sourced `_dev/tests/fixture-repo.sh`, so any fixture repository it tried to build silently failed with "command not found" and the resulting probe failures pointed at the code under test instead of the harness. It cost a debugging cycle here and would cost the next author the same. Fixed in this REQ (one `source` line) rather than deferred — the file was already open and in scope, and leaving a harness that misattributes its own failure would have been worse than the one-line change. — gate: trivial

**Minor:**
- The repo URL is derived by string surgery on the tarball URL. It handles the canonical GitHub branch-tarball shape and degrades honestly (`no repository URL supplied and none derivable`) on anything else, but a consumer pointing `DO_WORK_UPSTREAM_URL` at a non-GitHub archive gets one route, not two. That is the right trade — the alternative is a second env var nobody asked for — but it is worth stating in the guide if a second host ever appears.
- The git route clones the default branch rather than the branch named in the tarball URL. For the canonical `main` URL these agree. A `DO_WORK_UPSTREAM_URL` naming a non-default branch would fall back to the wrong ref, which is a silent difference rather than a failure.

### Restatement Sweep

**Triggered** — the diff changes how upstream archives are obtained. Swept `upstream_url`, `archive/refs/heads`, and `curl -fsSL` across `skills/`, `tools/`, `README.md`, and `_dev/`. Results: the two `tools/` scripts flagged by finding F2 still carry their own `curl` and are exactly REQ-218's subject, left untouched by design; the bootstrap heredoc and its README twin were updated together and their byte-identity is still enforced; `skills/do-work/actions/version.md` is the one prose consumer of the update flow and was updated in the same commit; `docs/prescribed-shell-primitives.md` describes the download primitive rather than the fetch policy and needed no change.

### Requirements Checklist

- [x] Two ordered routes, HTTP through the hardened primitive then git — delivered
- [x] `git archive` with `--prefix`, never a worktree copy — delivered and adversarially proven
- [x] Winning route reported; total failure names both outcomes — delivered
- [x] Both callers honor `DO_WORK_UPSTREAM_URL` and delegate the fetch — delivered
- [x] Bootstrap heredoc and README twin get retry flags only, byte-identical — delivered
- [x] `version.md` names the escape hatch, engine still the safety boundary — delivered

### Acceptance Testing

**Result: Pass**
- `bash _dev/tests/update-script-behavior.sh` — exit 0, including every pre-existing updater probe.
- `bash _dev/tests/install-suite-behavior.sh` — exit 0; this is what caught D-02 before it shipped.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines.
- Finding-Closure Ratchet: the named regression tests are the three cases in `update-script-behavior.sh`, with case 2 as the F3 guard, and each was measured failing before passing.

### Suggested Additional Testing

- A real `do-work update` against a genuinely rate-limited codeload would confirm the route reporting reads well in the situation it was built for; the fixture proves the mechanism, not the message's usefulness under stress.
- A `DO_WORK_UPSTREAM_URL` pointing at a non-default branch is untested and would currently clone the default branch — worth a case if that usage is ever supported rather than merely possible.

### Follow-up REQs Created

None — REQ-218 already owns the remaining `tools/` bypass, and the two Minor findings are trade-offs rather than defects.

## Lessons Learned

**What worked:** Building the *wrong* implementation on purpose. Swapping `git archive` for a `tar` of the clone and watching exactly one assertion fire is the only way to know that the export-ignore guard is load-bearing rather than incidentally satisfied — and it took about a minute.

**What didn't:** Requiring the fetcher beside the installer unconditionally. It read as symmetry with the validator and replacer, and it broke the canonical bootstrap, which supplies `--archive` and never fetches. A dependency check belongs on the branch that has the dependency. Separately, the first fixture run failed inside the *harness*: `update-script-behavior.sh` had never sourced `fixture-repo.sh`, so `fixture_repo_init` was simply not a command and every downstream assertion blamed the fetcher.

**Worth knowing:** `git archive` is the only one of `cp -R`, `rsync`, `tar`-the-clone, and `git archive` that honors `export-ignore` — measured on this tree as 734 tracked files versus 294 archive entries, and the manifest validator happily approves the leaky tree. Also: a script mirrored into two layouts cannot reach a sibling subtree by one relative path, so either the sibling is mirrored too or the script probes both depths; probing keeps the two copies byte-identical, which is what the mirror contract requires.

## Orientation

`do-work update` and the suite installer now survive a host that stays rate-limited: they fall back to a shallow clone repacked with `git archive`, report which route they used, and name `DO_WORK_UPSTREAM_URL` when neither works. Lives in the suite update/install subsystem — a new `tools/fetch-upstream-archive.sh` (mirrored into `skills/do-work/tools/`) now sits between both callers and the network. **[MAP CHANGED]** — this adds a component where there was previously a direct `curl` in each caller, and it is the seam REQ-218's ratchet will hold shut. The `git archive` requirement is load-bearing rather than stylistic: it is the only fetch mechanism that honors `export-ignore`, so it is what keeps the maintainer-only tree out of consumer installs.

Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` — all referenced paths resolve; its `curl -o` lesson still reads true and now has a fetcher above it rather than a correction.
