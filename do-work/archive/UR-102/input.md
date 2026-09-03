---
id: UR-102
title: 'Close the review tap and delete dead gate work before optimizing'
created_at: 2026-09-03T11:42:36Z
requests: [REQ-531, REQ-532]
word_count: 402
---

# Close the Review Tap and Delete Dead Gate Work Before Optimizing

## Summary

The roadmap survey of 2026-09-03 found the queue at equilibrium: intake tracks output because reviews and builds mint most of the queue (29 of 45 pending REQs), and the maintainer gate is slow because the test corpus pins prose sentence by sentence and re-runs its own harness. The maintainer chose three dispositions. D1 and D2 become the two REQs below. D4 (hand-triage of the remaining queue) is a disposition applied in-session through the canonical cancel and fold transactions, so it mints no REQ.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-531 | D1: a review finding below impact-critical is recorded in the report and never mints a REQ; the maintainer captures one by hand when wanted |
| REQ-532 | D2: delete the test scripts the gate never executes and the harness self-tests nested inside the aggregate suite, before UR-100 optimizes what remains |

## Batch Constraints

- REQ-531 and REQ-532 are independent and may run in parallel.
- No new sentence pins in `_dev/tests/contract-regressions.sh`; that file does not grow past its current line count (UR-100 constraint). Pins that assert the old minting behavior are deleted, not rewritten.
- The full gate (`bash _dev/tests/maintainer-verify.sh`) exits 0 at every integrating commit.
- The duplicate strict JavaScript lane and the line-count ratchet named in D2 are already owned by REQ-520 and REQ-519; REQ-532 carries only the deletions.

## Full Verbatim Input

> ```
> [maintainer, message 1]
> do-work roadmap <- and tell me if it is on the right path.
> 
> 
> The entire development cycle seems to be a never ending story, now we have 46 reqs and they get more faster then being implemented.
> 
> One of the problems that I can see is this test suite:
> 
> """
> time bash ./_dev/tests/maintainer-verify.sh
> maintainer-verify: checking Go go1.26.1 or newer
> maintainer-verify: checking ShellCheck 0.11.0 or newer
> maintainer-verify: ShellCheck warning-level lint (83 tracked files)
> maintainer-verify: gofmt formatting check (223 tracked files)
> maintainer-verify: aggregate contract suite
> [... full gate transcript: contract regression checks passed; queue-kanban tests ok 180.100s; strict JavaScript lane PASS 54.07s; browser lane SKIP; do-work-cli packages ok, finalization 51.600s, publication 29.216s, corehelpers 23.770s ...]
> Maintainer verification passed.
> bash ./_dev/tests/maintainer-verify.sh  304.60s user 121.46s system 95% cpu 7:24.53 total
> """
> 
> [assistant verdict, 2026-09-03, plan file do-work-roadmap-and-sorted-eich.md]
> Not on the right path, and the repo already knows it. ai-reports/2026-08-20_1131_remaining-work-and-the-req-treadmill/ concluded two weeks ago: "The queue is not a backlog. It is an equilibrium. The exit is a rule change, not a sprint." Since 2026-08-20: 214 REQs created, 195 completed, 459 releases in 30 days. 39 of 45 queued REQs are maintenance: true; 29 of 45 were spawned by a review or an addendum; 0 of 45 carry impact-negligible, so the skip filter removes nothing. The gate is slow because the instruction prose has a test per sentence: 1,009 Go tests with zero t.Parallel, 19,177 lines of shell tests, contract-regressions.sh alone 8,468 lines with about 341 sentence-pin assertions, 55 JavaScript probes run twice per gate, the gate's own self-test re-run 15 times from inside the aggregate suite, 6 dead test scripts.
> D1 Close the tap before draining: review findings below impact-critical land in the review report only; the maintainer mints a REQ by hand if wanted.
> D2 Delete before optimizing the gate: remove the dead test scripts, the self-test call from inside the aggregate suite, and the duplicate strict JavaScript lane; add a line-count ratchet on contract-regressions.sh. Then UR-100 runs against a smaller target.
> D3 Run UR-098 as the only substantive stream.
> D4 Hand-triage the rest in one sitting, following the 2026-09-01 disposition model.
> D5 Freeze maintenance capture for a fixed window and use do-work on a repository that is not do-work.
> 
> [assistant question] Which parts of the disposition should I capture as REQs now?
> 
> [maintainer selection, 2026-09-03]
> D1 Close the tap (Recommended), D2 Delete before optimizing the gate, D4 Hand-triage the other 31
> ```

---
*Captured: 2026-09-03T11:42:36Z*
