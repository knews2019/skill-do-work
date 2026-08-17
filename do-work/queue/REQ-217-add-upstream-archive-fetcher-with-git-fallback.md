---
id: REQ-217
title: Add the upstream archive fetcher with a git fallback
status: pending
created_at: 2026-08-17T17:16:28Z
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
- [ ] **[PLAN]:** (Agent: Read `_dev/primes/prime-shell-commands.md`, both caller scripts, and REQ-216's landed primitive. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
