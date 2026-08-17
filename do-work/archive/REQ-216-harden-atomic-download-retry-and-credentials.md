---
id: REQ-216
title: Teach atomic-download retry and optional credentials
status: completed
completed_at: 2026-08-17T19:04:40Z
commit:
claimed_at: 2026-08-17T18:57:52Z
created_at: 2026-08-17T17:16:28Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-17T18:58:30Z
  basis:
    - trivial short-circuit
user_request: UR-049
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-217, REQ-218]
batch: resilient-upstream-fetch
write_set:
  - skills/do-work/scripts/atomic-download.sh
  - skills/do-work/docs/prescribed-shell-primitives.md
  - _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Teach Atomic-Download Retry and Optional Credentials

## What

Add transient-failure retry and opt-in GitHub credentials to the shipped download primitive at
`skills/do-work/scripts/atomic-download.sh`, so every caller inherits the resilience instead of it being
patched into one download site.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add the four flags to the existing single `curl`, and build the optional `-H` argument with `set --` rather than a bash array — macOS ships bash 3.2, where expanding an empty array under `set -u` is an unbound-variable error, while `"$@"` is exempt. Read `GH_TOKEN` first, then `GITHUB_TOKEN`. Leave the mktemp/rename/status shape untouched. Add a fake `curl` that models curl's own internal retry loop so the retry is observable, plus header-log assertions for both token variables.
- [x] **[APPLY]:** Three files, all declared in `write_set`. No change to the publication contract, no new dependency, no shell retry loop wrapped around curl.
- [x] **[UNIFY]:** `git diff --stat` → `_dev/tests/prescribed-shell-scripts-behavior.sh` (+72), `skills/do-work/docs/prescribed-shell-primitives.md` (+4/−1), `skills/do-work/scripts/atomic-download.sh` (+12/−1). ShellCheck (via `maintainer-verify.sh`, 46 tracked files, warning level) flagged `SC1007` on two `GH_TOKEN=` clearings in the new test; corrected to `GH_TOKEN=''` and re-run clean. No debug artifacts, no leftover fixtures, no `console`/`TODO` markers in the diff.

## Why (if provided)

Feedback-brief finding **F1** — the four download sites in this repo are each a single `curl` attempt with
no retry and no credential support. Finding **F2** — `atomic-download.sh` is the canonical primitive for
exactly this and is bypassed by both `tools/` scripts, which is why the fix would otherwise have to be
written three or four times.

## Detailed Requirements

- Add `--retry 3 --retry-delay 2 --retry-max-time 60` to the existing `curl` at
  `skills/do-work/scripts/atomic-download.sh:26`.
- Do **not** use `--retry-all-errors` — it needs curl ≥ 7.71, and plain `--retry` already treats HTTP 429 as
  transient from curl 7.51.0 (2016). The 7.51 floor is consistent with how this repo has reasoned about curl
  versions before: `skills/do-work/CHANGELOG.md:2127` rejected `--remove-on-error` for requiring ≥ 7.83.
- When `GH_TOKEN` or `GITHUB_TOKEN` is non-empty, send `-H "Authorization: Bearer <token>"`. Opt-in only —
  absent or empty means the request goes out exactly as it does today.
- Preserve the publication contract **exactly**: private adjacent `mktemp` sibling, rename only on success,
  non-zero exit preserved, no `.download.*` scratch left behind on any path. The script runs under `set -u`
  (not `-e`) and captures `$?` explicitly into `download_status` — keep that shape.
- Update the **Atomic download publication** section of `skills/do-work/docs/prescribed-shell-primitives.md`
  (heading at line 63) so the canonical guide describes the retry and credential behavior. The heading text
  itself must not change — `_dev/tests/prescribed-shell-canonicalization.sh` greps for it with `grep -Fqx`.

## Constraints

- **Additive flags only.** The brief explicitly does not ask for a change to the publication contract, and
  neither does this REQ.
- The existing behavior case at `_dev/tests/prescribed-shell-scripts-behavior.sh:54-62` must keep passing.
  Its fake `curl` parses `-o` and ignores unrecognized flags, so the new flags do not disturb it — verify
  rather than assume.
- No new runtime dependency.

## Dependencies

None. REQ-217 depends on this.

## Builder Guidance

Certainty: **Firm.** The change is small and the shape is specified. Resist expanding it — the git fallback
belongs to REQ-217 and the ratchet to REQ-218.

## Red-Green Proof

**RED prompt/case:** A fake `curl` placed on `PATH` that returns 429 on its first invocation and succeeds on
its second — the same idiom as the existing partial-download case at
`_dev/tests/prescribed-shell-scripts-behavior.sh:54-62`. Invoke `atomic-download.sh` against it and assert
the target file is published complete.
**Why RED now:** `atomic-download.sh:26` is a single `curl -fsSL -o` attempt. The first 429 fails the script
outright, the target is never published, and the fake `curl`'s second invocation never happens.
**GREEN when:** The fetch succeeds on the retry and the target contains the complete second-attempt payload,
with no `.download.*` sibling left behind.
**Validation:** User confirmed (capture approved the three-REQ split and the C1 scope; route 1 dropped).

## Finding-Closure

Origin: external feedback brief, triaged via `do-work-toolbox validate-feedback` — findings F1, F2, and
proposal C1, all verdict **Accept**. Surface-cost: **earned** — the incident is named and reproducible
(sustained codeload 429 on `do-work update`), the added surface is one flag set plus one conditional header,
and the RED test above replays it. The named regression test is RED test 1, landing in
`_dev/tests/prescribed-shell-scripts-behavior.sh`.

## Full Context

See `do-work/user-requests/UR-049/input.md` for the complete verbatim brief.

## Triage

**Route:** A — Direct to Builder

**Reasoning:** The REQ names the file, the line, the exact flag set, the flags to avoid and why, and the test idiom to copy. Nothing about the "where" or the "how" needed discovery.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route A.

## Implementation Summary

**Files changed:**
- `skills/do-work/scripts/atomic-download.sh` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** Added `--retry 3 --retry-delay 2 --retry-max-time 60` to the single existing `curl`, and an opt-in `Authorization: Bearer` header built from `GH_TOKEN`, falling back to `GITHUB_TOKEN`. Both are additive; the mktemp-beside-target, rename-on-success, preserved-exit-status shape is untouched, and the script still runs under `set -u` with `$?` captured explicitly into `download_status`.

The optional header is assembled with `set --` rather than a bash array. macOS ships bash 3.2, where expanding an empty array under `set -u` is an unbound-variable error; `"$@"` is exempt from that rule, so the empty case stays safe without duplicating the flag list across two `curl` invocations. The two real arguments were captured into `source_url` and `target_path` before the positional list is rebuilt.

`--retry-all-errors` was deliberately not used: plain `--retry` has treated 429 as transient since curl 7.51.0, and `--retry-all-errors` would raise the required version to 7.71 for nothing — the same reasoning that rejected `--remove-on-error` for requiring 7.83.

The canonical guide's **Atomic download publication** section now describes both behaviors. Its heading is unchanged, as `_dev/tests/prescribed-shell-canonicalization.sh` matches it with `grep -Fqx`.

**Tests touched:** three new named cases in `_dev/tests/prescribed-shell-scripts-behavior.sh` — transient-429 retry, `GH_TOKEN` credential, `GITHUB_TOKEN` fallback. Named-case count 31 → 34.

## Qualification

Passed — 3 files verified, all requirements traced, no debug artifacts.

- All three files are in the REQ's captured `write_set`; nothing else was touched.
- Substantive: 12 added lines in the script are the whole behavior change; the tests add a fake `curl` that models curl's own retry loop rather than asserting on flags.
- Requirements traced: the four flags added and `--retry-all-errors` avoided; both token variables honored with empty treated as absent; publication contract byte-for-byte unchanged; guide section updated with its heading intact; the pre-existing failure case at lines 54–62 verified still passing rather than assumed.
- Flowing: `download_token` is read, tested, and consumed; nothing is set and unused.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh` (baseline, RED, GREEN); `bash _dev/tests/maintainer-verify.sh`

**Result:** ✓ prescribed-shell suite exit 0, 34 named cases; ✓ maintainer-verify exit 0, zero FAIL lines.

**Red-green validation:** ✗ RED — with the three new cases in place and the script reverted, the suite exits 1 with seven failures: `retry case did not survive a transient 429`, `retry case did not publish the successful attempt`, `retry case did not let curl retry the rate-limited transfer`, `credential case returned nonzero`, `credential case did not send GH_TOKEN as a bearer credential`, `fallback-credential case returned nonzero`, `fallback-credential case did not fall back to GITHUB_TOKEN`. → ✓ GREEN — all seven pass; the transfer log records exactly 2 attempts, the published target holds the second attempt's payload, no `.download.*` sibling remains, and the no-token case sends no header at all.

The fake `curl` serves 429 on its first attempt and 200 on the next, looping internally the way real curl does under `--retry`. It therefore survives only if the caller actually passed a retry allowance — the assertion is on behavior under rate limiting, not on the presence of a flag string.

**Existing tests updated:** none. The pre-existing partial-publication case (lines 54–62) passes unchanged, confirming the REQ's constraint that its flag-ignoring fake `curl` is undisturbed by the new arguments.

*Verified by work action*

## Review

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every bullet delivered, including the two explicit "do not" constraints |
| Code Quality | 95% | `set --` needs its comment, and has it; the alternative was duplicating the flag list |
| Test Adequacy | 95% | Behavioral fake rather than a flag-string assertion; both token variables covered |
| Scope | 100% | Exactly the three captured `write_set` files |
| Risk | Low | A bearer token is now sent when the environment carries one — opt-in, and `-fsSL` keeps it out of output |
| Acceptance | Pass | Seven assertions fail on the old script, pass on the new one |

**Verdict: Approve** — resilience now lives in the primitive, which is what makes REQ-217 and REQ-218 worth doing.

### Findings

**Minor:**
- A token is sent to whatever host the caller names, not only GitHub. The brief scoped this to GitHub downloads and every current caller is a GitHub URL, so narrowing by hostname would be speculative today — but it is the obvious next question if a non-GitHub caller ever appears.

**Nit:**
- `--retry-delay 2` with `--retry 3` means a rate-limited fetch can now take ~6 seconds longer before failing. Intended, and bounded by `--retry-max-time 60`.

### Restatement Sweep

**Triggered** — the diff changes what `atomic-download.sh` guarantees its callers. Swept `atomic-download` across `skills/`, `tools/`, and `_dev/`. The canonical **Atomic download publication** section is the only prose that describes the helper's behavior and was updated in the same commit. The two `tools/` scripts that bypass the helper are precisely REQ-218's subject and were left alone; `skills/do-work/actions/capture.md` cites the section by anchor for its no-clobber mechanics, which this REQ did not touch.

### Requirements Checklist

- [x] `--retry 3 --retry-delay 2 --retry-max-time 60` added — delivered
- [x] `--retry-all-errors` avoided, with the version reasoning recorded — delivered
- [x] `GH_TOKEN` / `GITHUB_TOKEN` sent as a bearer credential when non-empty, opt-in — delivered
- [x] Publication contract, `set -u` shape, and explicit `$?` capture preserved — delivered
- [x] Guide's **Atomic download publication** section updated, heading unchanged — delivered
- [x] Pre-existing behavior case still passing — delivered (verified, not assumed)

### Acceptance Testing

**Result: Pass**
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh` — exit 0, 34 named cases.
- `bash _dev/tests/maintainer-verify.sh` — exit 0. Its ShellCheck pass caught `SC1007` on the new test's `GH_TOKEN=` clearings, which was fixed before the final run.
- Finding-Closure Ratchet: the named regression test is the retry case, and it fails before / passes after.

### Suggested Additional Testing

- A real rate-limited fetch against codeload is the incident this came from and cannot be replayed in the suite; worth one manual `do-work update` under a live 429 if the opportunity arises.
- The bearer header is asserted through a fake `curl`'s argument log. A real request against a private repository would confirm the header is accepted as well as sent.

### Follow-up REQs Created

None — the remaining findings from this brief are already REQ-217 and REQ-218.

## Lessons Learned

**What worked:** Modelling curl's *internal* retry loop inside the fake rather than asserting that the flag string was passed. The test now fails for the reason the incident failed — a 429 that nobody retried — instead of for a missing substring, and it would keep passing if the flags were ever expressed differently.

**What didn't:** The first version of the no-token case used `GH_TOKEN= GITHUB_TOKEN=` to clear both variables, which ShellCheck flags as `SC1007` (a space after `=` usually means a mistyped assignment). `GH_TOKEN=''` says the same thing unambiguously.

**Worth knowing:** macOS ships bash 3.2, where `"${array[@]}"` on an empty array under `set -u` is an unbound-variable error — the standard trick for optional command arguments does not work in any script that has to run there. `set --` plus `"$@"` does, because `$@` is explicitly exempt from `set -u`.

## Orientation

Every caller of the shipped download primitive now survives a rate-limited host and can authenticate when the environment offers a token, without any caller changing. Lives in the prescribed-shell primitives subsystem (`skills/do-work/scripts/atomic-download.sh`, documented in `docs/prescribed-shell-primitives.md`). The publication contract is unchanged, so the system's shape is unchanged — this is the primitive gaining the resilience its callers were each about to grow separately, which is what REQ-217 and REQ-218 then build on.

Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` — all referenced paths resolve; its `curl -o` lesson still reads true and gains a companion rather than a correction.
