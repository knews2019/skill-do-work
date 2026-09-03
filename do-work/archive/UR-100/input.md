---
id: UR-100
title: 'Make the maintainer gate cheap to run per REQ'
created_at: 2026-09-02T21:27:16Z
requests: [REQ-518, REQ-519, REQ-520, REQ-521, REQ-522, REQ-523]
word_count: 500
---

# Make the Maintainer Gate Cheap to Run per REQ

## Summary

`bash _dev/tests/maintainer-verify.sh` takes about 6.5 minutes and the work pipeline runs it twice per REQ, so 167 REQ commits in 14 days cost on the order of 36 hours of gate time. The analysis report at `ai-reports/2026-09-03_0010_maintainer-verify-gate-cost-analysis/index.html` measured every stage: about 40% of one run is the same work repeated (the board's JavaScript probes run three times, the 14 aggregate sub-suites run one after another, no board test uses `t.Parallel`). These six REQs run the full gate once per REQ, add a path-scoped fast lane, remove the repeated work, and log every run.

Message 3 (update the report with plain-terms contents, timings, and technology) was done directly in the same session and minted no REQ.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-518 | A1: run the full gate once per REQ by skipping the baseline when HEAD is the last recorded green revision |
| REQ-519 | A2 + A6: a path-scoped fast lane the pipeline runs per REQ, the full gate at integration, and a line-count ratchet on the contract file |
| REQ-520 | A3: run the board's JavaScript probes once per gate run |
| REQ-521 | A4: run the aggregate contract sub-suites in parallel |
| REQ-522 | A5: opt queue-kanban tests into `t.Parallel` |
| REQ-523 | Log every gate run (start, duration, exit status, full output) and commit the log by exact path |

## Batch Constraints

- Order: REQ-518 first, then REQ-519. REQ-520, REQ-521, REQ-522 and REQ-523 are independent of each other and may run in parallel.
- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417); REQ-519 adds the ratchet.
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Full Verbatim Input

> ```
> [maintainer, message 1]
> This test is crazy, because of it my development time got really slow. 
> What is this doing? Is it an overkill?
> 
> Show me a breakdown as an ai-report use do-work ai-report.
> 
> Also include improvement suggestions. 
> [screenshot: `time bash ./_dev/tests/maintainer-verify.sh` ending in `6:27.79 total`]
> 
> [assistant report the maintainer pointed at, section "Improvements, in the order to do them"]
> A1 Run the gate once per REQ. Record the revision hash of the last green run. Before dispatch, if HEAD equals that hash, take the baseline as green without running. Keep the Step 6.5 run. Saves 6.5 min per REQ, half of all gate time.
> A2 Path-scoped fast lane. maintainer-verify.sh --changed: lint and vet always; Go tests only for a module whose files changed, with the test cache on; the aggregate only when skills/**/*.md, scripts, or _dev/tests changed. The full uncached gate stays the integration-commit check. Typical REQ from 388 s to 30-60 s.
> A3 Run the JS probes once. Set the strict marker on the ordinary run whenever Node is present, so TestMain's zero-probe guard applies to that run, and delete the separate lane. Narrow the two zero-probe meta-tests to re-execute one cheap test instead of the whole set. Saves about 110 s per full run.
> A4 Run aggregate sub-suites in parallel. Launch the 14 sub-suite processes together and collect statuses; they already own private fixture roots. Saves about 80 s per full run (149 s to about 65 s, bounded by update-script at 60 s).
> A5 t.Parallel() in queue-kanban. 0 of 451 tests opt in; the package runs strictly serial on 8 cores. Saves about 100 s per full run (146 s of test time to 40-50 s wall).
> A6 Freeze the contract file. A new sentence-predicate lane must delete one or land as a Go behavior test. Split what exists by owning action file. Little wall-clock; stops the maintenance curve.
> 
> [maintainer, message 2]
> ok, do-work capture-request to optimize it
> 
> [maintainer, message 3]
> also update ai-reports/2026-09-03_0010_maintainer-verify-gate-cost-analysis/index.html and show me what the test contains, how much time it takes, also tell me in plain terms what it does test and why it exists.
> 
> I also want to know the underlying technology as well.
> 
> [maintainer, message 4]
> I also want to have a log file of, when the suite ran, how much time it took, and what was the output of it, and also commit this log
> 
> [capture-time answers, 2026-09-03]
> Q1 Which improvements become REQs? -> A1 to A5 as REQs, A6 as policy folded into the fast-lane REQ.
> Q2 Which gate does the pipeline run per REQ once a fast lane exists? -> Fast lane per REQ, full gate at integration.
> Q3 How does the pipeline stop running the full gate twice per REQ? -> Skip the baseline when HEAD is the last green revision.
> Q4 What numbers make this done? -> Full gate under 3 min, typical REQ gate under 60 s.
> Q5 Who commits the gate run log? -> The gate script commits its own log entry.
> ```

---
*Captured: 2026-09-02T21:27:16Z*
