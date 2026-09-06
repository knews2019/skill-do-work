---
id: REQ-553
title: '[impact-negligible] Source one do-work-cli launcher preamble instead of hand-rolling it in every launcher'
status: completed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-551]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
write_set: [tools/do-work-cli-preamble.sh, skills/do-work/tools/do-work-cli-preamble.sh, tools/install-do-work-suite.sh, tools/fetch-upstream-archive.sh, tools/replace-text-section.sh, tools/validate-suite-manifest.sh, skills/do-work/tools/install-do-work-suite.sh, skills/do-work/tools/fetch-upstream-archive.sh, skills/do-work/tools/replace-text-section.sh, skills/do-work/tools/validate-suite-manifest.sh, skills/do-work/tools/checks/associate-files.sh, _dev/tests/staged-skills-contract.sh, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-05T00:41:15Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route B
    - 13-file write set
    - 3 subsystems involved
    - 4 acceptance criteria
completed_at: 2026-09-05T08:05:05Z
commit: eabee330
release_at: 2026-09-05T08:05:05Z
---

# Source one do-work-cli launcher preamble instead of hand-rolling it in every launcher

## What
After REQ-551 deletes the four toolbox copies, nine shell files still hand-roll the do-work-cli launcher preamble in two spellings (`for cli_candidate in` in the four root tools and their four byte-locked mirrors; `launcher_arguments=(--format text)` in `tools/checks/associate-files.sh`). Promote one sourceable preamble beside the root tools, with its byte-locked mirror under `skills/do-work/tools/`, that resolves the launcher path and the `--format text` argument array, and source it from every remaining copy.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** — Read the REQ in full, `_dev/primes/prime-shell-commands.md`, all of `_dev/primes/lessons-shell-commands.md`, `skills/do-work/docs/prescribed-shell-primitives.md`, `suite/modules.tsv`, `skills/do-work/crew-members/coding-guardrails.md`, `CLAUDE.md`, then all nine callers, `_dev/tests/staged-skills-contract.sh`, `_dev/tests/audit-lockins.sh`, `_dev/tests/contracts/probe-lanes.sh` and both differential fixtures before writing anything. Evidence: the plan named the `${BASH_SOURCE[0]}` anchor and the byte-copy mirror rule before any file was created, and I took green baselines of …
- [x] **[APPLY]:** — Scope stayed inside the declared write set. Evidence: `git show --stat HEAD` lists exactly 12 files — the two new preamble files, the eight launchers, `associate-files.sh` and `_dev/tests/staged-skills-contract.sh`; `git diff --cached --name-only | grep '^do-work/'` printed nothing before the commit; `_dev/tests/audit-lockins.sh` is untouched and handed back above as a seam.
- [x] **[UNIFY]:** — `git show --stat`: 12 files changed, 112 insertions(+), 83 deletions(-). Reviewed every changed file's full diff: each launcher lost the same nine-line block and gained four lines plus two ShellCheck comments, no debug output, no leftover scaffolding, no unrelated edit; the two preamble files are `cmp`-identical. Linters and suites, all after the final byte of the commit: `shellcheck --severity=warning` over all twelve files exit 0; the repo linter's own `--format=gcc --shell=bash --severity=warning` single-file argv clean on all eleven shell files; `bash _dev/tests/action-shell-blocks.sh` …

## Why
Two mutually incompatible spellings of one primitive, copied instead of sourced, is the 0.202 canonicalization class; the ratchet `_dev/tests/prescribed-shell-canonicalization.sh` does not look at this pattern.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 2, sweep_key `cli-launcher-preamble-copied`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -40. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- New `tools/do-work-cli-preamble.sh` (sourced, not executed) and its mirror `skills/do-work/tools/do-work-cli-preamble.sh`; add the pair to the byte-identical mirror list in `_dev/tests/staged-skills-contract.sh`.
- The bootstrap path stays self-contained: `tools/install-do-work-suite.sh --print-bootstrap-command` prints a literal heredoc and must keep working before any skill is installed, so the preamble lives beside the root tools, never under `skills/do-work/scripts/`.
- `skills/do-work/tools/checks/associate-files.sh` switches from the `launcher_arguments=(--format text)` spelling to the sourced preamble.
- Root tools and their mirrors change in lock-step (the staged-skills contract requires byte identity).
- Reproduce at dc8a64e3 (prints 13 files; 9 after REQ-551): `rg -n --no-heading 'for cli_candidate in|^launcher_arguments=\(--format text\)$' --glob '*.sh' skills tools`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- No launcher behaviour change: the differential fixtures for the installer and updater must pass unchanged.
- Prime `_dev/primes/prime-shell-commands.md` first; the sourcing must survive `set -euo pipefail` and a missing Go toolchain the same way the current preambles do.
- Lock-in limit: hand-rolled launcher preambles outside the preamble pair: 0 after this REQ (today 13); the Reproduce command prints at most 2 paths, the preamble file and its byte-identical mirror (verify repair 2026-09-03: the plan line's "target ≤ 1" counted the helper once, the mirror makes it two).

## Dependencies
Depends on REQ-551, which deletes four of the thirteen copies so this REQ touches nine files instead of thirteen.

## Builder Guidance
Firm on one sourced preamble and on the bootstrap constraint; latitude on the helper's name and on whether the argument array is a function or a variable.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements after REQ-551 has landed.
**Why RED now:** It prints nine files hand-rolling the preamble.
**GREEN when:** It prints at most the preamble file itself and its mirror; installer and updater fixtures green; the lock-in pins hand-rolled preamble copies at 0 outside the preamble pair.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for cli-launcher-preamble-copied.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear — one sourceable preamble replacing nine hand-rolled copies — but how each of the nine callers must locate the shared file, and whether the byte-locked mirror needs a modules.tsv declaration, has to be established by reading the callers. Exploration required; the shape of the fix is not in doubt.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*
## Exploration

The nine hand-rolled copies were located by their two spellings — `for cli_candidate in` in the four root tools and their four byte-locked mirrors, and `launcher_arguments=(--format text)` in `tools/checks/associate-files.sh` — and each caller was read to find where the resolver ends and its own work begins.

The finding that shaped the design came from running the fixtures rather than reading them. A single unconditional `source` line fails `_dev/tests/install-suite-behavior.sh`, because that fixture builds an archive root whose `tools/` receives only the launchers while `skills/` is copied whole. The eight mirrored launchers are the same bytes in two directories, so one source line has to resolve from `tools/` and from `skills/do-work/tools/`. `_dev/tests/update-script-behavior.sh` builds the same shape, so two fixtures depend on the two-depth probe.

`suite/modules.tsv` was checked and needs no change: it maps four module directories rather than files, and the new staged file lives inside the already-declared `skills/do-work` module.

## Scope

**Files I will touch:**
- `tools/do-work-cli-preamble.sh`
- `skills/do-work/tools/do-work-cli-preamble.sh`
- `tools/install-do-work-suite.sh`
- `tools/fetch-upstream-archive.sh`
- `tools/replace-text-section.sh`
- `tools/validate-suite-manifest.sh`
- `skills/do-work/tools/install-do-work-suite.sh`
- `skills/do-work/tools/fetch-upstream-archive.sh`
- `skills/do-work/tools/replace-text-section.sh`
- `skills/do-work/tools/validate-suite-manifest.sh`
- `skills/do-work/tools/checks/associate-files.sh`
- `_dev/tests/staged-skills-contract.sh`
- `_dev/tests/heavy-lanes.json`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go`
- `_dev/tests/audit-lockins.sh`

`_dev/tests/audit-lockins.sh` was an **integration seam**, not a builder write: another builder held that file concurrently, so the exact block and its anchor were handed back and the orchestrator applied them inside the merge commit.

**Acceptance criteria:** every caller still reaches do-work-cli with unchanged argv and exit status; the preamble pair stays byte-identical under the staged-skills contract; the lock-in pins hand-rolled copies at zero and goes red on a reintroduced one.

## Implementation Summary

**Files changed:**
- `tools/do-work-cli-preamble.sh` (new)
- `skills/do-work/tools/do-work-cli-preamble.sh` (new, byte-identical mirror)
- `tools/validate-suite-manifest.sh` (modified)
- `tools/replace-text-section.sh` (modified)
- `tools/fetch-upstream-archive.sh` (modified)
- `tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/validate-suite-manifest.sh` (modified)
- `skills/do-work/tools/replace-text-section.sh` (modified)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modified)
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/checks/associate-files.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `_dev/tests/audit-lockins.sh` (modified)
- `_dev/tests/heavy-lanes.json` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go` (modified)

**What was done:** One sourceable preamble file now resolves the do-work-cli path and the launcher argument array, and the nine shell files that each hand-rolled that preamble source it instead. A lock-in assertion pins hand-rolled launcher preambles at zero outside the preamble pair.

The preamble anchors on BASH_SOURCE zero, which under source is the preamble's own file, not the caller's directory. That is what makes one file correct both for the root copy, which must reach up into the staged tools directory, and for the staged copy, which sits beside do-work-cli itself. It also removes any dependence on what the caller names its directory variable: eight launchers call it script_dir and the associate-files check calls it script_directory.

Eight of the nine callers are one byte-identical file kept in two places, so their single source line has to work from both locations. They probe two depths for the preamble. That probe is not defensive padding: the installer behavior fixture builds a synthetic archive root whose tools directory receives only the three launchers while the skills tree is copied whole, and an unconditional single-path source failed there with a concrete no-such-file error before the fallback was added.

Review found the original no-guard decision wrong, and the continuation commit gave all five launchers a typed guard for a missing preamble. Each guard keeps that launcher's own message prefix and its original exit status, and names the preamble file rather than do-work-cli so two guards four lines apart cannot print the same string for two different missing files.

The continuation also closed two coverage gaps. The new root preamble was added to the updater and installer heavy lanes and to the Go coverage test table, because an uncovered path marks the whole plan uncertain and selects every lane. The lock-in exclusion was narrowed from a trailing-path match to two exact repository paths using grep with the whole-line and fixed-string flags.

The merged range carries nine other requests besides this one, so its 52-file stat is wider than this manifest. All fifteen files listed above appear in that stat, and no file this request touched is missing from the list.

**Implementation range:** `027cffc3..eabee330`. Builder commit `7f51a7c56c3e38a069a23459830acc130c4a59ae`, continuation commit `f44d8cf5eedac624fadc779ea92e6e7f51575254`.

## Decisions

- **D-01 — The preamble anchors on its own path, not the caller's:** it reads BASH_SOURCE zero, which under source is the preamble's own file. One file is therefore correct for the root copy, which must reach up into the staged tools directory, and for the staged copy, which sits beside do-work-cli. It also removes any dependence on the caller's variable name. Referencing a caller-supplied script_dir would have silently resolved wrong for the associate-files check, which sits one directory deeper.
- **D-02 — Launchers probe two depths for the preamble:** required, not defensive. Callers one through eight are the same bytes in two locations, so one source line has to work from both, and the installer behavior fixture builds an archive root whose tools directory holds launchers with no preamble beside them. The builder measured this: the unconditional single-path version failed that fixture with a concrete error before the probe was added. It is the same two-depth shape the launchers already used for do-work-cli.
- **D-03 — No guard for a missing preamble. Wrong as written, superseded by D-10:** the original claim was that a failed source stops the script and bash names the missing file. That holds only for the four launchers running set -euo pipefail. It is false for the upstream-archive fetcher and its staged twin, which run set -u with no set -e: the failed source does not stop the script, execution falls through to the empty-path check, and the script dies with an unbound-variable error naming neither the missing file nor the launcher. Two further exit statuses had also silently changed from 2 to 1. This deviated from the request's requirement that the sourcing survive set -euo pipefail the same way the current preambles do, and the continuation corrected it.
- **D-04 — Two ShellCheck directive lines per caller:** ShellCheck follows a source directive only with external sources enabled, which neither the shell-block lint nor the maintainer verify script passes. Without following, every caller draws SC2154 for an unassigned launcher_arguments, which is a warning and so fails the gate. The first attempt looked clean only because all files were linted in one command, which makes ShellCheck read the preamble as an input; single-file linting, which is what the gate does, failed. Adding external sources to the two test files was rejected because they are outside the write set and the verify script's self-test shim asserts an exact three-argument argv. A repository-root shellcheck config was rejected because it changes lint behavior for every shell file in the repository. The chosen fix is a narrow per-line disable plus a source-path pointer, so a genuine typo elsewhere in the file is still caught.
- **D-05 — launcher_arguments stays an array variable, not a function:** the request gave latitude here. A wrapper function would have removed the nine SC2154 disables, but it would have reordered argv for the installer, where the repo-root flag currently precedes the format flag, and it does not fit the associate-files check, which needs a compatibility-shim environment variable and an exec. No launcher behavior change outranks nine comment lines.
- **D-06 — The loop variable kept its name:** it lives in a sourced file, so it would leak into the caller's namespace; the preamble unsets it and its directory variable after the loop. The existing spelling was kept because the request's Reproduce command greps for it and expects to find it in the preamble pair and nowhere else.
- **D-07 — The associate-files check loses one word from one error message:** it used to exec bash against a hard-coded relative path and now execs the resolved variable. If do-work-cli were missing, bash previously named the path and now names an empty string. The exit status of 127 and every success path are identical. The script has no missing-launcher guard today and adding one is out of scope, so the difference was recorded rather than papered over. No test covers that case.
- **D-08 — The staged-skills contract gains existence checks, not a mirror list:** the request said to add the pair to the byte-identical mirror list, but no such list exists any more. The contract derives the mirrored set from disk, so every root tool that also exists under the staged tools directory must compare equal, and the new pair is covered automatically. What that loop cannot see is a pair where one side is missing, so the staged copy was added to the core files list and the root copy to the retained bootstrap tool list. Deleting either half now fails the contract. This is a plain deviation from the request's literal wording, driven by the list no longer existing.
- **D-09 — Net line count is plus 29, not the plan's minus 40:** the nine hand-rolled blocks came out at minus 83 lines and the preamble pair went in at plus 48 for two 24-line files, which is the shape the plan predicted. The overrun is the 18 ShellCheck directive lines that D-04 forced across the nine callers, plus the preamble's header comment. The builder trimmed both once, taking the preamble header from 11 comment lines to 7 and the per-launcher probe comment from two lines to one, and stopped rather than deleting the explanation of why the two-depth probe exists. The finding this request closes is one primitive with two spellings copied instead of sourced, and that is now zero outside the pair.
- **D-10 — The preamble guard is now earned, which D-03 wrongly denied:** the incident is concrete and was reproduced. A launcher running set -u without set -e continues past a failed source and dies on an unbound variable. The coding guardrails ask for the incident, the surface that remains and what keeps it live: the incident is that transcript, the surface is one line per launcher, and the lock-in plus the five typed messages keep it visible. The earlier claim that no incident earned it was an assertion untested against the two launchers that behave differently.
- **D-11 — The guards name the preamble file, not do-work-cli:** the suggested shape said to reuse each launcher's existing typed message verbatim. Each launcher's prefix, wording and exit status were kept, but the filename inside the message was changed, because the two guards now sit four lines apart and a verbatim reuse would print the identical string for two different missing files. The builder checked first that no test, script or document asserts those strings: the only matches outside the changelog are the launchers themselves and the update script. Statuses are unchanged, which is the contractual part.
- **D-12 — The associate-files check gets exit 2, not its accidental 127:** it never had a typed behavior here. Before this request a missing do-work-cli produced an exec failure with 127; after the first commit it was 1 through set -e. Neither was designed. It was given 2, matching both its sibling launchers' cannot-run status and its own existing silent exit 2 for a bad repo-root argument, with a message in its own voice. This is the one status in the set that does not restore a prior value, which is why it is called out.
- **D-13 — The lock-in exempts by exact path, not by a regex:** the exclusion uses grep with the whole-line and fixed-string flags against the two literal repository paths. Those flags cannot be widened by accident the way an anchored regex can. Three comment lines above the filter record why, naming the toolbox path the reviewer used, so the narrowing survives the next person who wants to shorten two long lines.
- **D-14 — One long guard line in the upstream-archive fetcher:** its preamble guard is a single printf of roughly 300 characters carrying both output lines through one embedded newline, rather than the five-line block that would mirror the file's existing style. The file already carries a printf of roughly 230 characters, so the precedent exists, and keeping all five guards the same one-line shape was judged worth more than local style symmetry.

## Qualification

Passed the request-bound advance qualify and scope-drift gates for the cumulative range `027cffc3..eabee330`, both satisfied. Fifteen files across the request's two merge commits, including the `_dev/tests/audit-lockins.sh` block handed back as an integration seam and applied by the orchestrator inside the merge commit.

Independent review proved argv and exit-status parity rather than assuming it: both trees rebuilt side by side with a fake CLI printing its argv, all nine callers run with space-containing arguments at two exit statuses, transcripts identical. It broke path resolution six ways and all held.

Verified here after merge: all five mirror pairs byte-identical by `cmp`, and the lock-in red on a reintroduced hand-rolled copy then green once removed.

The P-A-U boxes were reconciled from the builder hand-back, which is where worktree dispatch puts them.
## Testing

**Red-green validation:** before the change the request's Reproduce command printed nine shell files hand-rolling the launcher preamble in two spellings. The hand-back states the lock-in value plainly: "hand-rolled launcher preambles outside the preamble pair **= 0** (13 before REQ-551, 9 at the start of this REQ)." After the change the same command prints exactly two paths, the preamble and its mirror:

```
$ rg -n --no-heading 'for cli_candidate in|^launcher_arguments=\(--format text\)$' --glob '*.sh' skills tools
tools/do-work-cli-preamble.sh:16:for cli_candidate in \
tools/do-work-cli-preamble.sh:27:launcher_arguments=(--format text)
skills/do-work/tools/do-work-cli-preamble.sh:16:for cli_candidate in \
skills/do-work/tools/do-work-cli-preamble.sh:27:launcher_arguments=(--format text)
```

The hand-back notes those line numbers are from the pre-trim revision and are 13 and 24 in the committed file. Two paths either way.

The lock-in in `_dev/tests/audit-lockins.sh` was driven red and green in both directions in the continuation, in this order:

```
=== GREEN (current tree) ===
Audit lock-in regressions passed.
exit=0

=== RED direction 1: the reviewer's third copy under do-work-toolbox ===
FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: <worktree>/skills/do-work-toolbox/tools/do-work-cli-preamble.sh
exit=1

=== RED direction 2: a hand-rolled copy back in a launcher ===
FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: <worktree>/tools/zz-red-probe.sh
exit=1

=== GREEN again ===
Audit lock-in regressions passed.
exit=0
```

The missing-preamble failure was also reproduced red before the guards were written. The worst case was the launcher running `set -u` without `-e`:

```
-- fetch-upstream-archive.sh a b
<path>/fetch-upstream-archive.sh: line 27: ...do-work-cli-preamble.sh: No such file or directory
<path>/fetch-upstream-archive.sh: line 28: do_work_cli: unbound variable
  -> exit 1
```

After the guards, on the same trees with the same invocations, that becomes a typed message naming the missing file, and the five measured exit statuses are 2, 2, 1, 1, 2 — each launcher's original status for "my runtime is not here", except the associate-files check, which is covered by D-12.

Named test functions: `TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint` passes after the heavy-lane coverage rows were added (`--- PASS ... (0.00s)`, package `ok ... 0.003s`). Two other functions in that package fail, and the hand-back states explicitly that they fail identically without the change: `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes` and `TestShippedGitIsolationPreservesGenericLaneInheritance`, at `heavy_runtime_contract_test.go:107` and `:181`. The builder stashed the working tree and re-ran both at the merged tip, getting the same two failures at the same lines with none of the edits present. Under the isolation environment the lanes actually use (`GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null`) the fingerprint test passes and only the Git-isolation one remains, which points at this container's Git configuration. Diagnosed as not belonging to this change.

**Controls preserved:**

- `_dev/tests/install-suite-behavior.sh` — the installer differential fixture the request required to pass unchanged. It is also the fixture that proves the two-depth preamble probe is load-bearing, because it builds an archive root whose `tools/` holds launchers with no preamble beside them. Green baseline was taken on the untouched worktree before any edit, so "passed" means "passed the same way it passed before".
- `_dev/tests/update-script-behavior.sh` — the updater differential fixture, also baselined green before the change.
- `_dev/tests/staged-skills-contract.sh` — protects byte identity of every root tool that also exists under the staged tools directory, and now also the existence of both halves of the preamble pair.
- `_dev/tests/action-shell-blocks.sh` — the repository's shell-block lint over fenced blocks and shipped shell files.
- `_dev/tests/contracts/core-checks.sh` — pins the associate-files check's early guard: `--repo-root` with no value still exits 2, confirmed after the guards were added.
- `_dev/tests/contract-regressions.sh` — the fast aggregate. One FAIL, diagnosed below as not this change.

**Module verification:**

| Command | Result |
|---|---|
| `shellcheck --severity=warning` over all twelve touched shell files | exit 0 |
| Same files, the repo linter's own `--format=gcc --shell=bash --severity=warning`, one file at a time | exit 0 |
| `cd / && shellcheck -x --severity=warning` over the launchers from an unrelated working directory | exit 0, proving the source-path pointers resolve |
| `bash _dev/tests/audit-lockins.sh` | green, red on a toolbox-path copy, red on a launcher copy, green again |
| `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` | `staged skills contract: PASS`, exit 0 |
| `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/install-suite-behavior.sh` | `suite installer behavior probes passed.`, exit 0 |
| `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/update-script-behavior.sh` | `update-script behavior probes passed.`, exit 0 |
| `bash _dev/tests/action-shell-blocks.sh` | `Shell-block lint passed: 73 fenced blocks and 33 shipped shell files; ShellCheck enabled.`, exit 0 |
| `python3 -c "import json; json.load(open('_dev/tests/heavy-lanes.json'))"` | parses, no output |
| `gofmt -l skills/do-work/tools/do-work-cli/internal/heavyverification/` | no output |
| `go test -count=1 -run 'TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint' ./internal/heavyverification/` | PASS in 0.00s, package ok in 0.003s |
| `bash _dev/tests/contract-regressions.sh` | one FAIL, see below |

Argv parity was re-measured after the guards, against a fake CLI that prints its argv, on two trees built from `7bbbc325` and from the current worktree. All nine callers, every argument containing a space, at `FAKE_CLI_EXIT=0` and `FAKE_CLI_EXIT=7`: "IDENTICAL transcripts (argv + exit status)", with an empty `diff -u` in both runs. The bootstrap heredoc's md5 is in the transcript and matches.

Mirror equality after the guards, by `cmp`: all five pairs identical, with md5 sums `b000412e739d76ae8cb7f239d018e042` for the preamble, `9650b5a92637852792527d0da4f07d0e` for the installer, `3a08e1c10bf4c13a4ecfed72078541e2` for the manifest validator, `79279f4f672ae432ce047bcd1bec4d03` for the section replacer and `7474bc1a7de9765ebf725a228185a80d` for the archive fetcher.

The one `contract-regressions.sh` FAIL is diagnosed as not belonging to this change:

```
FAIL: .../_dev/tests/session-start-hook-behavior.sh took 35s; each test file must finish under 30s
```

Run alone that probe takes `real 0m26.354s` and prints `SessionStart hook behavior probes passed.` It breaches 30s only when the aggregate runs probes in parallel on this container, and this change touches no hook file and no file that probe reads.

Durations beyond the Go package timings and the hook probe timings above are not in the hand-back, so none are claimed for the shell suites.

## Discovered Tasks

- `_dev/tests/session-start-hook-behavior.sh` breaches the 30s per-file limit enforced by `_dev/tests/contract-regressions.sh` when the aggregate runs probes in parallel (35s and 37s in two runs), but takes 26s run alone. Pre-existing capacity mismatch on this container, not a hook regression. → report only
- Two tests fail in the heavy-verification package independently of this change: `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes` and `TestShippedGitIsolationPreservesGenericLaneInheritance`, at `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_runtime_contract_test.go:107` and `:181`. Confirmed failing at the merged tip with the change stashed. The Git-isolation one still fails under the lanes' own isolation environment. → queue as follow-up
- The `--format text` prefix is still not defined once in the repository. 24 shipped shell files write it inline, and seven scripts under `skills/do-work/tools/checks/` do not source the preamble at all: `archive-collision.sh`, `blanked-req-scan.sh`, `preflight.sh`, `qualify.sh`, `record-commit-hash.sh`, `scope-drift.sh`, `uncommitted-inventory.sh`. Outside this request's declared scope and deliberately not widened into. → queue as follow-up
- `skills/do-work/tools/checks/associate-files.sh` has no guard for a missing `do-work-cli.sh`. Per D-07, a missing launcher now makes bash name an empty string instead of a path; exit status 127 and every success path are unchanged, and no test covers that case. → queue as follow-up
- Release paperwork is owed for this commit. The releases prime says a commit changing shipped files under `skills/`, `tools/` or `suite/` is a release, so a VERSION bump, a `CHANGELOG.md` entry and its byte-identical mirror under the skill are due. Those files are outside the builder's write set and shared by all ten builder branches in this run, so they were left to finalize. → report only

## Review

**Overall: 94%**
**Acceptance: Pass.** The reviewer proved argv equality rather than assuming it: both trees rebuilt side by side with a fake CLI printing its argv, all nine callers run with space-containing arguments, transcripts identical in argv and exit status. It broke path resolution six ways — relative, absolute from `/`, four `..` segments, via PATH, through a symlinked parent — and all held.

The two-depth probe was independently confirmed load-bearing: removing the fallback fails the install fixture with exactly the reported error, and a second fixture builds the same topology.

Three findings were closed in the continuation. A decision was factually wrong for the two callers running `set -u` without `-e`, where a failed source did not stop the script and execution died on an unbound variable naming neither the missing file nor the launcher. The heavy-lane coverage gap was recorded backwards — an uncovered path marks the run uncertain and selects every lane, so it forced the full heavy set rather than skipping one. And the lock-in exclusion matched by trailing path shape, so a third copy under a different module stayed green.

The reviewer judged the +29 lines against a planned −40 and concluded the cure does not cost more than the disease: 81 lines were nine copies of one resolver that could drift, replaced by one definition held in two files that `cmp` equal under an enforced contract.

## Lessons Learned

Verify a lint gate with the gate's own argv, not with a convenient batch. ShellCheck reads a sourced file as an input when it is passed in the same command, so linting a whole directory at once hides the SC2154 that single-file linting reports. Any check whose result depends on how many files are handed to it must be re-run in the exact shape the gate uses before it is called clean.

A failed `source` only stops the script when `set -e` is on. In a file running `set -u` alone, execution continues past the failure and dies later on an unbound variable that names neither the missing file nor the caller. When a shared preamble becomes a new required input for several scripts, check each one's `set` line individually and give each a typed guard, rather than reasoning about the group from the strictest member.

## Orientation

There is now one definition of the do-work-cli launcher preamble, held in two files that must compare byte-equal, and a lock-in that turns red the moment a tenth copy appears anywhere outside that pair. Every launcher that loses its preamble now says so by name and exits with the status it always used for a missing runtime, instead of dying on an unbound variable.
