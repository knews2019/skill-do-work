---
title: "Lessons from REQ-258: Split the prescribed shell behavior suite per script"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-04/REQ-258-split-the-prescribed-shell-behavior-suit.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-258: Split the prescribed shell behavior suite per script

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

`_dev/tests/prescribed-shell-scripts-behavior.sh` now carries 47 named cases, and the ten reservation-cleanup + timestamp-repair cases dominate its tail. If it keeps growing, per-script files may read and fail more legibly. Organizational only — no case changes.

## Solution summary

**What was done:** Lifted the suite's shared fixture preamble into a new harness file, split all 76 named cases into one file per script under test, and reduced `prescribed-shell-scripts-behavior.sh` to a runner that executes each case file as its own process, aggregates failures, and derives the reported case count from the case files at run time. No case was changed: the assertions, fixtures, and failure strings all moved verbatim, verified by a line-multiset comparison (1756 non-blank case lines in, 1756 out, zero missing, zero added). Two groups that were interleaved in the monolith (`generate-report-image` around `generate-report-image-batch`, `repair-req-timestamps` around `audit-archive-timestamps`) are now each contiguous in their own file, and the `publish-portfolio-summary` fixture setup that had been sitting under the `install-memory-hooks` header moved to the file it belongs to. Each case file is runnable on its own.

## What worked

- Splitting by a one-shot extraction script instead of by hand, then proving the result with a **line-multiset comparison** against the original. 1756 in, 1756 out, zero drift — a claim of "no case changes" that is checkable rather than asserted. Hand-splitting 1882 lines across 17 files would have made the same claim unverifiable.
- Auditing cross-group variable and function coupling *before* choosing the architecture. The answer (zero coupling) is what made per-file processes safe, and it took one script to learn rather than a failed run to discover.

## What didn't work

- The first pass at `prescribed_shell_finish` used `$([ "$n" -eq 1 ] && printf case || printf cases)` for pluralization. It works, but the `&&`/`||` ternary breaks silently the moment the first branch can fail, and it was cleverness bought for a cosmetic. Replaced with plain `if` blocks (D-05).
- Hoisting `core_checks` into the shared harness looked tidier than leaving it inside the `qualify` block. It would have been a case change and an unused-variable warning, for nothing. Reverted before it shipped: in a "no case changes" REQ, the tidying reflex is the failure mode.

## Worth knowing

- **`tools/checks/qualify.sh` cannot tell a moved line from an added one.** It runs `git diff` with no `-M`/`-C` and greps `^+`, so a relocation REQ gets FAILed on every pre-existing `TODO`/debug marker inside the moved text. This REQ hit it on four fixture strings that are byte-identical in HEAD. The correct response is to prove the lines pre-exist (`git show HEAD:<file> | grep`) and record the override with that evidence — not to un-check `[UNIFY]`, which is the reflex the audit exists to prevent. Queued as a discovered task.
- **A restructure's real blast radius is the text that plans around the old structure, not the code.** Every consumer of this file was fine (one caller, exit status only). What went stale was `RESTART-PROMPT.md` and three queued `write_set` fields describing a scheduling bottleneck that no longer exists. `write_set` self-heals at Step 5.5, so the damage is to human wave-planning, not to the build.
- Groups were interleaved in the monolith (`generate-report-image` around its `-batch` sibling, `repair-req-timestamps` around `audit-archive-timestamps`) and one fixture-setup block sat under the wrong case header entirely. Nobody could see either while it was one file. That is the concrete argument for the split, and it is worth stating more strongly than the REQ's "may read more legibly".

## Back-reference

See `do-work/archive/UR-056/REQ-258-split-the-shell-behavior-suite-per-script.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1cc1836`.
