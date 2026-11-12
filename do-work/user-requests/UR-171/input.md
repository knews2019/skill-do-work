---
id: UR-171
title: 'Phase-loading triage follow-ups: prime byte budget and wider reading measurement'
created_at: 2026-11-12T09:05:50Z
requests: [REQ-719, REQ-720]
word_count: 282
---
## Summary
The user promoted two accepted findings from the validate-feedback triage of 2026-11-11 into queue work. That triage checked a proposal to load do-work instructions by phase against commit 22f82fed (work procedures load by phase, version 0.315.3). Five of its seven findings were already done; these two were still open.

## Extracted Requests

| REQ | Source finding | Request |
| --- | --- | --- |
| REQ-719 | Finding 2 | Add a byte budget of about 6,000 bytes beside the ~60-line prime budget, and shrink the primes it flags where a shrink is possible. |
| REQ-720 | Finding 5 | Extend the reading-volume script and its report to count `work-reference.md` by cited section and the always-loaded crew files. |

## Batch Constraints
- Byte or token volume is used to spot unnecessary reading in unrelated phases. It is never proof of instruction quality and never a reason to delete a required rule.
- No order was requested between the two requests. Neither depends on the other.
- REQ-720 is measurement only: it changes no shipped file.

## Capture Clarifications
Asked during capture and answered by the user on 2026-11-12. The answers are recorded on each REQ.

1. Should the byte-budget task also shrink this repository's kanban-board prime file, which the new budget marks as too large (45 lines, 10,551 bytes)? The user's answer:

> ```
> rule and audit phase when we can if shrink is possible
> ```

2. Proof for the byte-budget task: the user chose "Audit check only": the prime audit's size check, run by hand before and after the change, with no new automated test.
3. Proof for the measurement task: the user chose "Script and report": the script output and the report table must both show the wider totals.

## Full Verbatim Input
> ```
> do-work capture-request: Two follow-ups from the phase-loading triage of 2026-11-11.
> 
> Finding 2, source: pasted text, no severity given.
> Verbatim claim: "Extend prime and instruction size checks beyond raw line counts (such as the ~60-line prime threshold that mainly relocates lessons) to measure bytes or estimated tokens as well. Use byte or token volume to spot unnecessary reading in unrelated phases, never as proof of instruction quality or a reason to delete required rules."
> Evidence: skills/do-work-toolbox/actions/prime.md:216 and :291 count lines only; skills/do-work/crew-members/general.md:12 and skills/do-work-toolbox/crew-members/general.md:10 say "roughly 60 lines"; _dev/primes/prime-kanban-board.md passes at 45 lines but is 10,551 bytes.
> Surface-cost: Earned. Incident: that 10,551-byte prime passes the line budget. Surface: one byte figure in an existing rule and one wc -c in the audit step. Test: the prime audit on this repo flags prime-kanban-board.md.
> Remedy: add a byte budget of about 6,000 bytes beside the 60-line budget, using the ceil(bytes / 4) formula from skills/do-work/actions/capture-reference.md:78, in prime.md and both general.md copies. Relocate overflow, never delete required rules.
> 
> Finding 5, source: pasted text, no severity given.
> Verbatim claim: "RED baseline: Trace the current `work` and `review-work` / `code-review` instruction paths across three representative scenarios: (a) a simple mechanical fix, (b) a complex multi-file change, and (c) a review-only run. Record the exact instruction files loaded, total bytes read, and every required check and recovery rule reached."
> Evidence: _dev/reports/work-phase-loading.md:18 excludes work-reference.md (119,671 bytes); work-selection-reference.md:53 sends a simple fix to its Request File Schema (35,527 bytes) and Schema Read Contract (17,233 bytes) sections; the split saved 12,693 bytes.
> Surface-cost: N/A.
> Remedy: extend _dev/tests/instruction-reading-volume.py and the report to count work-reference.md by cited section and the always-loaded crew files for the same three walkthroughs. Measurement only, no shipped file changes.
> ```
