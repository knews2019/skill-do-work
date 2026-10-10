## Delta re-review: REQ-656 (ai-report revise writes a superseding sibling bundle)

**Approve**: the delta closes F1, F2 (in the revise step), F3, F5 and F6 and adds no contradiction. One checklist line still restates the completed-work narrative for every report (N1), so F2's checklist part is still open.
Delta `6c837fab..85395d8e` (2 files, +7/-7) | cumulative `c44d9f29..85395d8e` | read at HEAD `85395d8e` | quick-scan depth

### Closure of the first review's findings
- F1 closed. `ai-report.md:150` (Step 8) and `:170` (checklist) now accept, for a revise, a `fail` whose only finding is the expected prior-bundle link. Guide `:55` appends the matching sentence. This matches revise step 5 (`:53`) and the reproduced judge behaviour (verdict=fail, exit 1, one `AI-REPORT-JUDGE-BROKEN-LINK`).
- F2 closed in revise step 5 (`:53`): a prior bundle that is not a completed-work report keeps its own section order in place of Step 5's narrative, and the step 3 re-check serves as the provenance ledger. The checklist part is still open (N1).
- F3 closed. `:27` now has the revise exception in parentheses.
- F4 left as report only, as agreed. No change.
- F5 closed. `:53` now reads "check on disk that the prior bundle the link names exists".
- F6 closed. Guide `:9` says "Every report-writing invocation", and `:31` says "One report invocation covers one UR or one REQ; a revise covers one existing bundle".

### New findings

**Minor:**
- N1 `ai-report.md:169` checklist: "Stakeholder narrative includes verdict, shipped change, problem/change, operation, ..." still applies to every report. A revise of a report that is not a completed-work report (the deploy guide in the UR) now keeps its own section order (`:53`), so this box cannot be ticked honestly. The first review's sweep listed `:169` under F2, and the delta did not touch it. Proposed fix: append " (a revise of a report that is not a completed-work report keeps the prior report's section order instead)." — impact-user-visible → report only

**Nit:**
- N2 `ai-report.md:53`: in "when the prior bundle is not a completed-work report, keep its own section order ..., and the step 3 re-check serves as the provenance ledger", the ledger clause reads as part of the condition. A revise of a completed-work report also has no ledger, because Step 1 builds it and "Steps 2 to 8" leaves Step 1 out, so Step 6 (`:132`) "Verify each claim against the provenance ledger" has no source. Proposed fix: "...keep its own section order in place of Step 5's narrative. For any revise, the step 3 re-check serves as the provenance ledger." — impact-negligible → report only
- N3 `ai-report.md:144` Step 7 still says "rerun until the verdict is `pass`". Revise step 5 overrides it and Step 8 is now widened, so a careful agent follows the revise rule. Optional fix: "rerun until the verdict is `pass` (for a revise, until only the expected prior-bundle link finding remains)". — impact-negligible → report only

### Checks
- Step 5 narrative (`:116`), Step 6 ledger and boundary (`:132`, `:134`), Step 7 (`:144`), Step 8 (`:150`), checklist (`:166-172`) and the guide (`:9`, `:21`, `:36-47`, `:55`, `:92-94`): read at HEAD. No new contradiction apart from N1-N3.
- Wording follows communication-style.md: plain words, no stock phrases, no em-dash chains. The `:150` parenthetical has three alternatives and is dense, but it is still clear.
- `git diff --check 6c837fab..85395d8e`: clean.
- Not run: probe, `maintainer-verify.sh` (forbidden by the brief), and an agent-driven revise. The delta is prose only.

## Review

**Overall: 92%** | 2026-10-10

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

(100+88+85+95)/4 = 92%.

**Important findings:** none
**Minor findings:** N1 `ai-report.md:169` checklist still requires the completed-work narrative for a revise of a report that is not a completed-work report; append the revise exception — impact-user-visible → report only; N2 (nit) `:53` ledger clause reads as conditional, so a completed-work revise has no ledger source for Step 6; make it apply to any revise — impact-negligible → report only; N3 (nit) `:144` Step 7 "rerun until the verdict is `pass`" not widened for a revise (step 5 overrides) — impact-negligible → report only
**Closed:** F1, F2 (revise step; checklist part open as N1), F3, F5, F6. F4 unchanged (report only).
**Acceptance:** Pass. Implementation stage: closures verified by reading HEAD. The judge's expected-link behaviour was reproduced in the first review. Diff check clean.
**Restatement sweep:** the delta redefined the end verdict a revise may accept, the revise section shape and ledger, the Do-NOT-use unfinished-target bullet, and the invocation scope lines. Grepped `skills/` and `README.md` for "verdict is `pass`", "until the verdict", "Every invocation", "One invocation covers", "unfinished or unsuccessful", "BROKEN-LINK", "provenance ledger", "Stakeholder narrative includes". Stale: `ai-report.md:169` (N1) and `:144` (N3). Checked, not stale: `architecture-report.md:103` (its own draft-directory rule), `completed-work-presentation-reference.md:34`, `:63` (ledger owner), `present-video.md` ledger lines, guide `:47` (normal target resolution; the revise section at `:94` covers revise), `capture.md:11` and `capture-guide.md:7` (unrelated "Every invocation").
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action (delta re-review)*

## Delta 2 (85395d8e..65547570)

The proposed texts were applied to `skills/do-work-toolbox/actions/ai-report.md` (+3/-3). `git diff --check` is clean.
- N1 closed. Checklist `:169` now has the revise exception for a report that is not a completed-work report, which matches revise step 5 (`:53`).
- N2 closed. `:53` now says "For any revise, the step 3 re-check serves as the provenance ledger" as its own sentence, so Step 6 (`:132`) has a ledger source for every revise.
- N3 closed. Step 7 `:144` now says to rerun "until only the expected prior-bundle link finding remains" for a revise, which matches step 5, Step 8 (`:150`) and checklist `:170`.
- No new contradiction. Re-read Steps 5-8, checklist `:166-172` and guide `:55` at `65547570`. The checklist `:166` "evidence ledger" is satisfied by the step 3 re-check for a revise.

The score stays 92%, Approve. Open findings: none (F4 is report only, as before).
