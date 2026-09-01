---
source_type: req_lesson
req_id: REQ-303
req_path: do-work/archive/UR-062/REQ-303-run-the-pinned-live-archive-assertions-only-in-a-suite-checkout.md
date: 2026-08-21
domain: testing
module: _dev/primes
tags: [testing, pinned, live, archive, assertions]
---

# Lessons from REQ-303: Run the pinned live-archive assertions only in a suite checkout

## What the REQ was about

`TestLiveArchiveDurationsMatchTheCalibratedFigures` (`durations_test.go:209-232`) pins exact medians
and counts from this repo's archive — 2026-07-31 median 2.5 with 2 completed / 1 kept, 2026-08-15
count 25 with median 19.6 — but `liveBoard` (`board_live_test.go:17-32`) resolves the repo root by
walking up from the test's working directory. The `_test.go` files ship (only `/do-work` is
export-ignored, not the tests), so in a consumer install the same test loads that project's queue and
fails on data it was never calibrated against. Apply REQ-282's Route B shape: run where the
assertions apply, report not-applicable elsewhere.

## Solution summary

`suiteCheckoutSkipReason` asks REQ-282's already-shipped
`resolveReleaseProbeVersionFilePath` whether the resolved root is a suite checkout, and returns a
reason naming that condition — and no path — when it is not.
`TestLiveArchiveDurationsMatchTheCalibratedFigures` consults it before building the board. The
pinned figures moved unchanged into `calibratedLiveArchiveFindings`, which returns one line per
disagreement instead of fataling, so a second test can feed it wrong figures and prove it still
bites. Production code is untouched.

## What worked

- **An inline `t.Fatalf` chain cannot be proven to still bite.** A pinned check that silently
  stopped biting is indistinguishable from a passing one, and editing the check to test it proves
  nothing about the shipped version. Returning findings from a function makes the bite testable with
  a wrong input, which is what `release_test.go` already did for the release probes.
- **A test for repo-independence must itself be repo-independent.** The first instinct was to assert
  `suiteCheckoutSkipReason(liveRepoRoot(t)) == ""` — which would pass here and fail in exactly the
  install the REQ exists to protect. Building both roots as fixtures is what makes the guard
  portable.
- **Prove a skip by running it, not by reading it.** `go test -c` plus a working directory inside
  a consumer fixture turned "it should skip there" into an observed SKIP line. The gate function's
  unit test would have passed either way.

## Back-reference

See `do-work/archive/UR-062/REQ-303-run-the-pinned-live-archive-assertions-only-in-a-suite-checkout.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `69f3319`.
