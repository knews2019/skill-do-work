---
source_type: req_lesson
req_id: REQ-298
req_path: do-work/archive/UR-056/REQ-298-sweep-unchecked-exit-status-across-the-shipped-scripts.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, review, sweep, unchecked, exit]
---

# Lessons from REQ-298: Review fix: sweep the unchecked-exit-status primitive across every shipped script

## What the REQ was about

REQ-268 closed five instances of one condition — a command or process substitution whose
exit status is discarded while only its content is judged, so a tool that never ran reads
as a tool that found nothing — and stated that condition in the headers of the two scripts
it touched. Its independent review then found the same primitive **inherited verbatim from
a third script that neither REQ-268 nor its Requirements named**:
`skills/do-work/tools/checks/record-commit-hash.sh` is the file REQ-268's own header cites
as its mandatory guard-style template, and it carries the pattern itself (for example
`head_blob_bytes="$(git cat-file -s … 2>/dev/null || true)"`). The copy direction was
template → copies, so fixing the copies and leaving the template means the next script
written from it re-imports the defect.

**Done means the class cannot recur:** the condition is stated once where shipped shell is
governed, and every shipped script that takes a substitution and judges only its content
either checks the status or says in place why the content alone is sufficient. Patching N
sites one at a time is what this REQ exists to avoid.

## Solution summary

Reproduced the fail-open in the repo's own incident guard, fixed it by asking whether the blob exists before asking its size, judged the three fail-safe sites in place, swept every shipped script and measured the result, stated the condition where shipped shell is governed, and added the narrow mechanical check the corpus can actually support.

## What worked

Building the broad check *first*, running it, and counting. The REQ recorded a dissent that predicted false positives and a user decision to build it anyway with an evidence-based escape hatch — and the only honest way to use that hatch is to have built the thing and measured it. Fifteen flagged, zero defects is a number; "I think it would flag too much" is an opinion, and the REQ was explicit that only the former would do.

## What didn't work

The first attempt to run `record-commit-hash-guards.sh` directly failed with "must exist and be executable" against a file that plainly was. The suite tests the *canonical* `tools/checks/` path, and `contract-regressions.sh` invokes it through a `sed` that repoints it at the shipped copy. Worth knowing before debugging a phantom permissions problem: in this repo a `_dev/tests/*.sh` file is not necessarily runnable on its own terms, because nothing auto-discovers them and some are rewritten at their call site.

## Worth knowing

The sharpest form of this defect is that **the fallback value looked like an answer**. `|| true` on a size query produces `""`, and `""` is what "there is no blob" also produces — so the guard could not tell a missing file from a broken tool, and the safe-looking branch was the wrong one for one of them. The general shape to watch: whenever a fallback value is in the same domain as a legitimate result, the failure has been laundered into data. `'?'` is safe precisely because it is not a number.

## Back-reference

See `do-work/archive/UR-056/REQ-298-sweep-unchecked-exit-status-across-the-shipped-scripts.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `01abc28`.
