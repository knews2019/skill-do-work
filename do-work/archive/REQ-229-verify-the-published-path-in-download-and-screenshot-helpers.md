---
id: REQ-229
title: Verify the published path in the download and screenshot helpers
status: completed
completed_at: 2026-08-18T01:43:47Z
claimed_at: 2026-08-18T01:38:34Z
domain: general
created_at: 2026-08-18T00:18:44Z
user_request: UR-042
addendum_to: REQ-225
effort_estimate: normal
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-18T01:40:00Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
maintenance: false
write_set:
- skills/do-work/scripts/atomic-download.sh
- skills/do-work/scripts/capture-screenshot.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
- skills/do-work/docs/prescribed-shell-primitives.md
---

# Discovered Task: Verify the Published Path in the Download and Screenshot Helpers

## What

`skills/do-work/scripts/atomic-download.sh` and `skills/do-work/scripts/capture-screenshot.sh` publish their payload and return without checking the path they actually wrote. When the destination path is occupied by a directory, `mv` and `ln` nest the payload inside it and exit zero, so both helpers report a success that did not happen. Close the same defect class already closed in four sibling helpers.

## Context

Found while implementing REQ-225, which states the rule once in the shipped guide as `## Verified exact publication`. Writing the condition down is what surfaced these two — the same defect had been found four times before, every time by a review sweep and never by reading the guide, which is the argument REQ-225 was captured on. Both instances below were reproduced, not inferred.

**`atomic-download.sh:44`** runs `mv "$download_path" "$target_path"`. With the target occupied by a directory the download nests as `<target>/<target-basename>.download.XXXXXX`, line 48 then clears `download_path` so the EXIT trap spares it, and the script exits 0. Reproduced against a `file://` source: exit 0, target still a directory, private file abandoned inside it. Callers include the suite installer's `SKILL.md` downloads and `tools/fetch-upstream-archive.sh`, both of which read the exit status as proof the file landed.

**`capture-screenshot.sh:33`** runs `ln "$copy_path" "$destination_path"`. `ln` refuses an occupied *file* destination, which is where the no-clobber guarantee comes from, but nests on an occupied *directory* and exits zero. Under `--staged` this compounds into data loss: the success path continues to `rm "$source_path"` at line 45 and destroys the staged screenshot — the only copy the capture dispatch holds — while the destination never receives it. Reproduced: exit 0, staged source deleted, destination still a directory holding an orphaned `.copying.XXXXXX` file.

## Requirements

- Both helpers verify the path they actually wrote after publishing, and fail closed when the write nested instead of publishing.
- Each helper removes only its own nested artifact and leaves the occupying directory exactly as it was, matching the four helpers that already do this (`publish-portfolio-summary.sh:102-163`, `generate-report-image.sh:112-117`, `generate-report-image-batch.sh:169-173`, `install-last30days.sh:98-103`).
- `capture-screenshot.sh` must not delete the staged source when publication did not happen — the staged source is preserved on every other failure path in that script and this one is the exception.
- Do not weaken either helper's existing guarantees: `atomic-download.sh` keeps rename-on-success and failure preservation, `capture-screenshot.sh` keeps byte verification and the no-clobber refusal on an occupied file.
- Delete the temporary sentence in `skills/do-work/docs/prescribed-shell-primitives.md` § Atomic download publication that records these two helpers as not yet making the check.

## Red-Green Proof

**RED prompt/case:** Two new cases in `_dev/tests/prescribed-shell-scripts-behavior.sh`: (1) invoke `atomic-download.sh` with a `file://` source and a target path occupied by a directory; (2) invoke `capture-screenshot.sh --staged` with a destination path occupied by a directory.
**Why RED now:** Both exit 0 today. Case 1 leaves the target a directory with a `.download.XXXXXX` file inside it. Case 2 additionally deletes the staged source, so the screenshot is gone and was never published.
**GREEN when:** Both cases exit nonzero, the occupying directory is unchanged and holds no nested private artifact, and case 2's staged source is still present.
**Validation:** Reproduced by hand during REQ-225; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] Auto-approved: critical severity (data-loss risk on `capture-screenshot.sh --staged`, silent false success in both). → Added to queue immediately.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect, both call sites, and the fix are all named exactly in the REQ, and the pattern is already written four times in this repo. What needed discovery was the shape of the behavior-test harness and how the four references structure their nested-path checks, so the fix reads as the fifth instance rather than a fifth invention.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The four references, read before writing anything.** `generate-report-image.sh:112-121` and `install-last30days.sh:98-112` share one shape: probe `"$destination/${staged##*/}"` after the rename, remove only that nested artifact, clear the variable the cleanup trap reads, print a `REFUSING:` line naming the directory, and exit nonzero. `publish-portfolio-summary.sh:102-114` and `:150-163` do the same twice, once per output. The fix here is that shape applied a fifth and sixth time.

**The harness.** `_dev/tests/prescribed-shell-scripts-behavior.sh` builds one fixture tree under `mktemp -d`, fakes `curl` on `PATH` for the download cases, and records failures through `fail_case` rather than exiting at the first one — so a case can assert several consequences and report each independently. Two existing cases already pin the guarantees this REQ must not weaken: the partial-publication case (a failed transfer never changes the target and leaks no scratch) and the coordinated-race case (exactly one of two concurrent writers installs).

**The cleanup trap is why the ordering matters.** `capture-screenshot.sh`'s `--staged` branch removes the staged source at line 45, *after* the `ln`. That is the data-loss path: the removal is gated on the `ln` having succeeded, and `ln` succeeds on a directory. The verification therefore has to sit between them.

*Explored by work action (inline, serial mode)*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/atomic-download.sh` (modify) — verify the published path after the rename.
- `skills/do-work/scripts/capture-screenshot.sh` (modify) — verify the published path after the link, before the staged source can be removed.
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — the two RED cases, and its case count.
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — delete the temporary caveat REQ-225 left recording this gap (D-01).

**Files I will NOT touch:** the four helpers that already implement the rule, `_dev/tests/prescribed-shell-canonicalization.sh` (the heading ratchet is unaffected).

**Acceptance criteria (restated from REQ):**
- [ ] Both helpers verify the path they actually wrote and fail closed when the write nested.
- [ ] Each removes only its own nested artifact and leaves the occupying directory exactly as it was.
- [ ] `capture-screenshot.sh` does not delete the staged source when publication did not happen.
- [ ] Neither helper's existing guarantees are weakened.
- [ ] The temporary caveat in the shipped guide is deleted.

## Implementation Summary

**Files changed:**
- `skills/do-work/scripts/atomic-download.sh` (modified)
- `skills/do-work/scripts/capture-screenshot.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**What was done:** Both helpers now probe the nested path after publishing, remove only their own artifact, print a `REFUSING:` line naming the occupying directory, and exit nonzero — the same shape the four sibling helpers already use. In `capture-screenshot.sh` the check sits between the `ln` and the `--staged` branch's `rm` of the source, which is the ordering that turns a silent data loss into a preserved staged file. Two cases were added to the shipped-script behavior suite, and the caveat REQ-225 wrote into the shipped guide to record this gap was deleted.

## Decisions

- **D-01**: Extended the write set to `skills/do-work/docs/prescribed-shell-primitives.md`. Requirement 5 asks for exactly this deletion, and the sentence is only true while the defect is open — leaving it would ship a guide that reports a gap it no longer has. Replaced rather than merely deleted, so the section keeps the pointer REQ-225 gave it and now states what both helpers actually do. DECIDE & STATE.
- **D-02**: Bumped the behavior suite's closing `(45 named script cases)` literal to 47 rather than deriving it. The number does not match any obvious count in the file — there are 40 case headers, not 45 — so it counts something whose definition is not recoverable from the file, and inventing a derivation would risk silently changing what the line reports. Adding two cases in the existing style makes the bump correct under whatever convention produced 45. The undefined, hand-maintained count is recorded in `## Discovered Tasks` rather than repaired blind. DECIDE & STATE.

## Discovered Tasks

- **[low]** `_dev/tests/prescribed-shell-scripts-behavior.sh:1099` closes with a hand-maintained literal — `Prescribed shell script behavior probes passed (N named script cases)` — whose N matches no derivable count in the file: there are 40 case-header comments against a reported 45. It is the closed-enumeration pattern `_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* names, in a summary line that reports the suite's own size, so a reader takes it as measured when it is remembered. Either derive it (count the case headers, and make the header shape the definition) or delete the count and report only pass/fail. Not repaired here because guessing the intended convention could silently change what the line claims.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing — both exit 0

**Red-green validation:** the REQ's captured RED, run as written.

- `atomic-download occupied-target case`: ✗ three assertions — *reported success for a publication that nested*, *abandoned its private file inside the occupant* → ✓. The case fakes a succeeding `curl` and points the helper at a target path occupied by a directory holding a pre-existing file.
- `capture-screenshot occupied-destination case`: ✗ three assertions — *reported success*, ***destroyed the staged source it never published***, *abandoned its private copy inside the occupant* → ✓. The middle one is the data-loss path, and it failed before the fix exactly as the REQ predicted.
- Both cases also assert the occupying directory survives with its contents intact, so the fix cannot pass by deleting the obstacle.

**Guarantees explicitly not weakened (REQ requirement 4):** the existing `atomic-download partial-publication` case (a failed transfer never changes the target and leaks no scratch), the `retry`, `credential` and `fallback-credential` cases, and `capture-screenshot coordinated-race` (exactly one of two concurrent writers installs) all pass unchanged.

**Reproductions re-run by hand,** against the same commands that produced the original finding:
- `atomic-download.sh` with a `file://` source and a directory-occupied target: exit 1, `REFUSING: out/SKILL.md is a directory — download discarded, existing directory left unchanged`, target still a directory, no nested `.download.*`.
- `capture-screenshot.sh --staged` with a directory-occupied destination: exit 1, `REFUSING: … screenshot not installed, staged source preserved: stage/screenshot-1.png`, **staged source still present**, no nested `.copying.*`.

**New tests added:** two cases in `_dev/tests/prescribed-shell-scripts-behavior.sh`.

**Existing tests updated (cross-REQ impact):** none — only the suite's own closing case count (D-02).

## Lessons Learned

**What worked:** Reading all four existing implementations before writing the fifth. The nested-path probe is only half the pattern; the other half is *where* it sits relative to whatever the success status authorizes. In `generate-report-image.sh` it guards a fallback re-stage, in `install-last30days.sh` it guards a rollback that opens with `rm -rf`, and here it guards the removal of the staged source. Copying the probe without copying that placement question would have produced a check that fires after the damage.

**What didn't:** Nothing failed, but one thing was nearly missed. The obvious reading of this defect is "the helper reports success wrongly" — a status bug. The severity is entirely in what reads that status: `ln` nesting is harmless on its own, and becomes data loss only because line 45 removes the staged source on the strength of it. A fix that verified after the removal would have passed a naive test and preserved the bug.

**Worth knowing:** This is the fifth and sixth time this defect class has been closed here, and the first time it was found by reading a written-down rule rather than by a review sweep — REQ-225 wrote the condition down and it surfaced these two immediately. The remaining shipped publication helpers were checked while the rule was in hand; these were the last two.

## Review

**Overall: 96%** | 2026-08-18T01:43:47Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 96% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None. The restatement sweep found no stale restatement: the shipped guide's caveat was the only text describing this gap and it was deleted as Requirement 5 asks, and the four sibling helpers' own sections are unaffected.

**Minor findings:** 1 (report only)
- The behavior suite's closing case count is a hand-maintained literal that matches no derivable count in the file. Bumped by two rather than repaired, and recorded in `## Discovered Tasks` — guessing the intended convention could silently change what the line claims.

**Acceptance:** Pass — all five restated criteria verified: four by the two new mutation-proof cases (each failed before the fix and passes after), and the guide deletion by the shipped reference contract staying green.
**Suggested testing:** 0 items
**Follow-ups created:** see Step 8 (Discovered Tasks); **sweeps appended to:** None
