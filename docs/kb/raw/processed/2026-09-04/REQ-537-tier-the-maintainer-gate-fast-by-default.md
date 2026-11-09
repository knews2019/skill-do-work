---
source_type: req_lesson
req_id: REQ-537
req_path: do-work/archive/REQ-537-tier-the-maintainer-gate-fast-by-default-heavy-on-request.md
date: 2026-09-03
domain: testing
module: _dev/primes
tags: [testing, tier, maintainer, gate]
---

# Lessons from REQ-537: Tier the maintainer gate: fast by default, --heavy on request

## What the REQ was about

`bash _dev/tests/maintainer-verify.sh` with no flags becomes the fast tier and the only canonical gate: toolchain floors, ShellCheck, gofmt, `go vet` for both modules, the fast aggregate, queue-kanban tests without JavaScript or browser probes, and do-work-cli tests. `--heavy` adds the board's JavaScript probes (strict marker on, zero-probe guard on), the browser lane when an engine is present, and the heavy aggregate probes. `--heavy-surfaces` prints the repo-relative path globs whose changes warrant a heavy run, so REQ-541 can consult it mechanically.

## Solution summary

**Files changed** (all in commit `8d9d1bb`):
- `_dev/tests/maintainer-verify.sh` (modified, +182/-117) — the tier mechanism, `--heavy`, `--heavy-surfaces`, and the two-tier self-test.
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (modified) — 0.269.0 to 0.270.0.
- `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` (modified) — the "0.270.0 — Two-Tier Maintainer Gate" entry, byte-identical mirror.
- `_dev/tests/contract-regressions.sh`, `_dev/tests/probe-batch.sh` — in `write_set`, deliberately unchanged (D-02).

## What worked

- Deriving the heavy board-test list instead of writing it down. The nine files share no filename token, so a hand list would have looked right and been wrong the first time a probe moved. Making the derivation fail loudly on zero matches is what keeps it honest, because the grep sits in a process substitution whose status nothing reads.
- Recording the miss in the commit body instead of reaching for path scoping. The REQ pre-authorized exactly one of those two responses, and the other would have shipped a `--changed` lane that REQ-538 through REQ-540 make unnecessary.

## What didn't work

- The REQ's own sequencing. REQ-537 delivers the tier seam, but every second the seam could save belongs to REQ-538 and REQ-540, so it lands as a 2 s difference and a mechanism. Shipping the knob before the things it switches is a defensible order for review size, and it means this REQ cannot be judged on its own numbers — the changelog entry says so, and so does this record.

## Worth knowing

- **A green strict-marker run is the probe-execution proof; a green plain run is not.** The board package's `TestMain` refuses exit 0 when the marker is set and the probe counter is zero. That is why `--heavy` can assert "the probes ran" while the default tier can only assert "the package passed".
- One package can be three quarters of a test gate. 467 s of a 641 s run came from `queue-kanban` alone, which is the whole reason the two-tier split cannot pay off before REQ-538.
- The gate runner's own log under `$TMPDIR/do-work-gate-runs/<revision>.log` retains per-package times long after the run. It is the cheapest place to re-derive a gate breakdown without paying for the gate again.

## Back-reference

See `do-work/archive/REQ-537-tier-the-maintainer-gate-fast-by-default-heavy-on-request.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8d9d1bb`.
