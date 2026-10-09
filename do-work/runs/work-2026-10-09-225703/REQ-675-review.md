## Review: REQ-675

**Approve**: the edited Steps 7 and 8 make a reviewer scope Acceptance to the stages it checked and name deployment and live acceptance as unassessed. Four small report-only findings.
Route A | merge range 623f00d5..f7dad3b1 (merge commit f7dad3b1, builder commit f12beed1)

### What's built
- `review-work.md` Step 7 defines four delivery stages (implementation, integration, deployment, live acceptance), says which apply is a judgment, and adds one paragraph: the Acceptance result scores only the stages this review exercised, and its one-line summary names them.
- Step 8 has a new first category, "Unassessed delivery stages", required whenever an applicable stage was not exercised.
- `docs/review-work-guide.md` Phase 3 has one matching paragraph in user words.

### Decisions / risks for you
- None blocking. The main residual risk is F2: the persisted `**Acceptance:**` line keeps the stage names only if the reviewer applies the Step 7 sentence while filling a template placeholder that gives no cue.

### Findings

**Important:**
- None

**Minor:**
- F1. `skills/do-work/crew-members/shared-principles.md:15` ("Acceptance cannot be exercised | Record Untested and the exact check that could not run") does not cover the new partial case, where implementation and integration were exercised but deployment or live acceptance could not be. Review Step 2 makes the reviewer read this row before Step 7, and the stage name "Live acceptance" shares the word "acceptance" with the score, so a reviewer can record Untested instead of a scoped Pass. `actions/work.md:366-368` has no Untested branch, so that result has no routing. Smallest fix: add "when only some stages can be exercised, score those and list the rest as unassessed (review-work.md Step 8)" to the row. — impact-rule-change → report only
- F2. `skills/do-work/actions/review-work.md:202` says "Name the stages you covered in the result's one-line summary" but does not say which line. The human report has three candidates: the `**Result:**` line (:273), the Scores table Notes cell (:294), and the persisted `**Acceptance:** [Pass/Partial/Fail/Untested] — [1-line summary]` line (:413), whose placeholder gives no cue. The persisted block keeps only `**Suggested testing:** [count] items`, so the stage names survive in the durable record, which `do-work-toolbox/actions/completed-work-presentation-reference.md:39` reads as acceptance evidence, only through that summary. Smallest fix inside the write set: name the `**Acceptance:**` line in the Step 7 sentence. The template stays unchanged. — impact-user-visible → report only
- F3. `skills/do-work/actions/sample-archived-req.md:99` (`**Acceptance:** Pass — component renders correctly with all avatar states`) names no stage. `actions/work.md:534` presents this file as a complete example of what the pipeline generates. The sample was already behind the current Append template (old `**Findings:**` line, no `**Restatement sweep:**` line, bare `None` follow-ups), so this adds one more drift to a known stale example. — impact-negligible → report only

**Nit:**
- F4. `skills/do-work/actions/review-work.md:460`, the Route A row of Calibrating Review Depth, says "Suggested testing is usually empty or 1 item". A Route A REQ with a deployment stage now needs at least two required lines (deployment and live acceptance). "Usually" keeps the two texts compatible. — impact-negligible → report only

### Requirements Checklist

- [x] Step 7 defines the four stages once, each in one plain sentence, worded as the REQ asks (`review-work.md:188-192`). Delivered.
- [x] Applicability is a judgment, and many REQs have no deployment stage (`:188`). Delivered.
- [x] Acceptance scores only the exercised stages (`:202`). Delivered.
- [x] An applicable unexercised stage is listed in Step 8 as "unassessed", by stage name (`:210`). Delivered.
- [x] One matching guide line in Phase 3 (`review-work-guide.md:36`). Delivered. It is two sentences, which fits the intent.
- [x] Constraint: no change to the Append to REQ File template, score table, Scoring Guidelines and caps, Verdict mapping, impact tokens, Step 10 or any status. Verified: the diff has three hunks in review-work.md (two inside Step 7, one inside Step 8) plus the guide hunk, with 11 insertions and 0 deletions.
- [x] Constraint: no retrospective mode, no restated counterexample rule, no repair-addenda change. Verified in the diff.
- [x] Constraint: release per `prime-releases.md`. N/A at this point. The release happens at run finalization, not in the merge range.

### Code review notes (quick scan, Route A)

- Clean prose and no em-dash in the added lines (`git diff 623f00d5..f7dad3b1 | grep '^+' | grep -c '—'` returned 0). `git diff --check` exit 0.
- D-01 (the unassessed rule lives in Step 8, Step 7 points to it) and D-02 (Pass/Partial/Fail/Untested bullets unchanged, scoping paragraph after them) are sound. They keep one home for each rule and leave the score values that `work.md` routes on byte-identical.
- D-03 (the one-line summary names the stages) is the right use of the latitude. F2 is the gap it leaves.
- D-04 (colon instead of the siblings' em-dash in the new Step 8 bullet) is a small style mismatch. It is accepted and not a finding.
- D-05 judgment: the overlap between the "Integration testing" category (third-party services, APIs, migrations, auth) and the new "Integration" stage is not confusing enough to matter. The stage line is required and keyed by stage name. The category is optional and scoped by its own examples. Nothing routes on Step 8 categories. The worst case is a reviewer listing a third-party check under "Integration: unassessed", which still surfaces the check. No finding.
- Decisions section: read from the REQ (copied from the hand-back) and from `REQ-675-handback.md`. Both agree, and every behavior-shaping choice in the diff is documented.
- P-A-U boxes: all three ticked.
- Coding-guardrails spot check: surgical, declared write set equals touched files. Nothing extra.

### Acceptance Testing

**Result: Pass** (implementation and integration)
- Implementation: walked the diff against every requirement above.
- Integration: re-ran the GREEN probe `do-work/runs/work-2026-10-09-225703/REQ-675-probe.sh` at f7dad3b1, exit 0 (including `shipped-package-reference-contract.sh: PASS`). The repository gate result at f7dad3b1 is the work action's (exit 0, recorded in `## Testing`). It was not re-run here.
- Behavioral exercise (below): GREEN.

#### Behavioral exercise

Scenario: a REQ titled "publish the updated rules page". Its diff correctly updates the rules page content. The local test suite passes. The reviewer has no access to the served site.

Applying the edited Steps 7 and 8 at f7dad3b1:
- Step 7: implementation is exercised (the diff does what the REQ asks). Integration is exercised (the merged tree's test suite passes). Deployment applies (the REQ says "publish", so the built page must reach the serving environment) and was not exercised. Live acceptance applies and was not exercised. The result covers only the exercised stages and they both work, so the result is Pass, with the covered stages named in the summary.
- Step 8: the "Unassessed delivery stages" category is required because two applicable stages were not exercised.

Exact `**Acceptance:**` line:

```
**Acceptance:** Pass — implementation and integration: the diff updates the rules page content as the REQ asks and the local test suite passes; deployment and live acceptance unassessed (no access to the served site).
```

Exact Suggested Additional Testing lines:

```
- Deployment: unassessed. Confirm the built rules page reached the serving environment, for example that the served file or the deployed build revision matches the merged commit.
- Live acceptance: unassessed. Open the served rules page and confirm the new rules show.
```

Result: **GREEN**. The Acceptance result is scoped to implementation and integration and names them in its one-line summary. Deployment and live acceptance are each listed as unassessed by stage name.

Contrast with the base text (`623f00d5`, Steps 7-8): nothing scopes the result or names stages. The code can run, so the "If you can't run the code" escape does not fire. Pass ("Feature works end-to-end as specified") follows from the passing suite. Step 8 has no required entry, and the Route A row allows an empty list. The base text yields `**Acceptance:** Pass — rules page content updated and local tests pass.` with, at most, an optional "Manual verification: check the rules page renders" line and no deployment or live line. That is the RED shape.

Can a reviewer still produce the RED shape with the edited text? Only by skipping two explicit sentences, one of them marked "Required". The remaining ambiguities are:
- F2: "the result's one-line summary" does not name the persisted `**Acceptance:**` line. A reviewer filling the Append template 200 lines later can write a bare summary there, which is the RED shape on the durable record even if the human report is correct. This is the most likely path back to RED.
- F1: the shared-principles row can pull a reviewer to Untested. That is wrong, but it is not the RED shape.
- "Which ones apply is a judgment" lets a reviewer decide deployment does not apply. For "publish" this is implausible, because Step 8's own example is the served rules page. For a less obvious REQ (for example, a docs change that a site build picks up) the definitions ("reached its serving environment or consumer install") are the only guide.
- "this review exercised" versus integration evidence from the work action's Step 6.5 run: Step 7.1 lets the reviewer rely on that run, but the scoping sentence counts only what "this review exercised". A strict reviewer may mark integration unassessed. That gives a narrower result, not the RED shape.

### Suggested Additional Testing

- Deployment: unassessed. After run finalization releases this change, confirm a consumer's installed `do-work/actions/review-work.md` contains the "Delivery stages" paragraph and the new Step 8 category.
- Live acceptance: unassessed. In a consumer project, run one real orchestrated review of a REQ with a served or installed output, and confirm its `**Acceptance:**` line names the covered stages and Suggested Additional Testing lists the rest as unassessed.

### Scores (on the record, not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All three requirements and every constraint met |
| Code Quality | 92% | Clean, minimal prose. F2's unnamed summary line is the one gap |
| Test Adequacy | 85% | Word-presence probe plus one-off behavioral exercise, as the REQ chose. No kept test, by design |
| Scope | 100% | Declared write set equals touched files |
| Risk | Low | No code. F1 can mislead a reviewer to Untested |
| Acceptance | Pass | Implementation and integration. Deployment and live acceptance unassessed |

Overall: (100 + 92 + 85 + 100) / 4 = 94.25, rounded down to 94%. No cap or penalty applies.

### Follow-ups created
- None (4 findings report only)

### Restatement sweep detail

Redefined elements: the meaning of the Acceptance result (it now covers only the delivery stages the review exercised, and its one-line summary names them) and the content of Step 8's Suggested Additional Testing (a new required category for unassessed stages).

Readers checked:
- `crew-members/shared-principles.md:14` ("Acceptance is inferred from narrower tests"): agrees. It asks for end-to-end exercise where applicable, and the new text scopes the result to what was exercised.
- `crew-members/shared-principles.md:15` ("Acceptance cannot be exercised"): stale for the partial case (F1).
- `actions/work.md:358-368` (Step 7 Pass/Partial/Fail routing) and `:518` (Error Handling, Acceptance = Fail): agree. The values and their names are unchanged. Untested has no branch there, which predates this REQ.
- `actions/review-work.md` Append to REQ File template (`:390-419`): unchanged in the diff. Its `[1-line summary]` placeholder is where F2 bites.
- `actions/review-work.md` Scoring Guidelines (`:221-230`), Verdict mapping (`:300`), Step 10 (`:334`): unchanged, and they agree.
- `actions/review-work.md` Calibrating Review Depth (`:460`): mild tension (F4). Verification Checklist (`:504`, `:506`): agrees.
- `docs/review-work-guide.md:34-36`, `:49`, `:61`: agree.
- `actions/sample-archived-req.md:96`, `:99`: the example line names no stage (F3).
- `tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go:327`: a fixture body only. Nothing parses the summary.
- `do-work-toolbox/actions/completed-work-presentation-reference.md:39`: reads `## Review` as acceptance evidence. It does not restate semantics, and it is why F2 matters.
- Broad grep of `skills/` for "Suggested Additional Testing", "Suggested testing", "Untested", "unassessed", "delivery stage", "live acceptance", "**Acceptance:**" and "Acceptance =": no other restatement of Acceptance semantics or Step 8 categories. The other "untested" hits are unrelated (test-coverage wording).

Not a wave end (REQ-675 is not its wave's last successful integration), so no inherited elements.

## Review

**Overall: 94%** | <REVIEW-TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1. `crew-members/shared-principles.md:15` "Acceptance cannot be exercised → Record Untested" does not cover the partial case (implementation and integration exercised, deployment or live acceptance not), so a reviewer can score Untested instead of a scoped Pass. — impact-rule-change → report only
- F2. `actions/review-work.md:202` "the result's one-line summary" does not name the persisted `**Acceptance:**` line (:413), whose placeholder gives no cue, so the durable record can still carry a bare Pass. Fix by naming that line in Step 7. — impact-user-visible → report only
- F3. `actions/sample-archived-req.md:99` example `**Acceptance:**` line names no stage. The sample was already behind the current template. — impact-negligible → report only
- F4 (Nit). `actions/review-work.md:460` Route A "Suggested testing is usually empty or 1 item" versus the now-required unassessed-stage lines. "Usually" keeps them compatible. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration: every requirement is in the diff, the GREEN probe re-run at f7dad3b1 exits 0, and the behavioral exercise on the "publish the updated rules page" scenario is GREEN; deployment and live acceptance unassessed (release happens at finalization, no consumer install checked).
**Restatement sweep:** redefined the Acceptance result's scope (only exercised delivery stages, named in the one-line summary) and Step 8's required unassessed-stage category; swept shared-principles.md:14-15, work.md:358-368 and :518, the review-work.md Append template, Scoring Guidelines, Verdict mapping, Step 10, Calibrating Review Depth and Verification Checklist, review-work-guide.md, sample-archived-req.md, advance_commands_test.go:327 and completed-work-presentation-reference.md:39; stale: shared-principles.md:15 (F1), sample-archived-req.md:99 (F3), review-work.md:460 (F4)
**Suggested testing:** 2 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Delta re-review (f7dad3b1..fec8d220)

**Verdict: Approve, 95%.** The one new sentence closes F2. No new finding.

Delta: commit 67e0ed86, merged as fec8d220. One line changed in `skills/do-work/actions/review-work.md:202` (Step 7). `git diff --stat f7dad3b1..fec8d220` shows 1 insertion and 1 deletion in that file only.

1. **F2 closed, nothing else touched.** The sentence now says to name the covered stages "both in the report and on the `**Acceptance:**` line of the Append to REQ File block, because that line is the durable record". That names the exact persisted line (`:413`) that F2 said was missing. The cumulative `git diff 623f00d5..fec8d220` still has only the two Step 7 hunks, the one Step 8 hunk and the guide hunk (11 insertions, 0 deletions). The Append to REQ File template (`:390-419`), Scoring Guidelines, Verdict mapping, Step 10, impact tokens and all statuses are unchanged.
2. **Wording.** Plain and clear. No em-dash in the added text (`git diff f7dad3b1..fec8d220 | grep '^+' | grep -c '—'` returned 0). `git diff --check` exit 0. The phrase "line of the Append to REQ File block" matches how the Verification Checklist (`:500`) already points at the **Restatement sweep:** line, so a reader meets a familiar pattern. "In the report" does not say which human-report line (the `**Result:**` line or the Scores Notes cell). That is acceptable: the human report is read once, and the durable line is the one that mattered.
3. **Behavioral exercise, re-applied.** Scenario: "publish the updated rules page". The diff is correct, local tests pass, no access to the served site. Edited Step 7: implementation and integration are exercised, deployment and live acceptance apply and were not. Result is Pass, scoped to the two exercised stages. The sentence now says those stage names go on the persisted `**Acceptance:**` line, so the bare `**Acceptance:** Pass — rules page content updated and local tests pass.` (the RED shape) now breaks an explicit instruction that names that line. The expected line is the same one the first review wrote. Step 8 still requires the two unassessed lines. Result: **GREEN**. The GREEN probe `REQ-675-probe.sh` re-run at fec8d220 exits 0.
4. **New restatements staled.** None. Grep of `skills/` for "one-line summary" finds no other reader of this rule. The template placeholder `[1-line summary]` stays generic, which the REQ's constraint requires. `docs/review-work-guide.md:36` says the result covers only the checked stages and does not restate where they are named, so it still agrees. F3 (the sample's `**Acceptance:**` line names no stage) is now one step further from the rule but is the same finding.

## Review

**Overall: 95%** | <REVIEW-TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1. `crew-members/shared-principles.md:15` "Acceptance cannot be exercised → Record Untested" does not cover the partial case (implementation and integration exercised, deployment or live acceptance not), so a reviewer can score Untested instead of a scoped Pass. — impact-rule-change → report only
- F2. Resolved in 67e0ed86 (merged fec8d220): `actions/review-work.md:202` now names the `**Acceptance:**` line of the Append to REQ File block as where the covered stages go. — resolved
- F3. `actions/sample-archived-req.md:99` example `**Acceptance:**` line names no stage. The sample was already behind the current template. — impact-negligible → report only
- F4 (Nit). `actions/review-work.md:460` Route A "Suggested testing is usually empty or 1 item" versus the now-required unassessed-stage lines. "Usually" keeps them compatible. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration: every requirement is in the diff, the F2 fix names the persisted `**Acceptance:**` line, the GREEN probe re-run at fec8d220 exits 0, and the behavioral exercise on the "publish the updated rules page" scenario is GREEN; deployment and live acceptance unassessed (release happens at finalization, no consumer install checked).
**Restatement sweep:** redefined the Acceptance result's scope (only exercised delivery stages, named in the one-line summary and on the persisted `**Acceptance:**` line) and Step 8's required unassessed-stage category; swept shared-principles.md:14-15, work.md:358-368 and :518, the review-work.md Append template, Scoring Guidelines, Verdict mapping, Step 10, Calibrating Review Depth and Verification Checklist, review-work-guide.md, sample-archived-req.md, advance_commands_test.go:327, completed-work-presentation-reference.md:39 and every "one-line summary" hit in skills/; stale: shared-principles.md:15 (F1), sample-archived-req.md:99 (F3), review-work.md:460 (F4)
**Suggested testing:** 2 items
**Follow-ups created:** None (3 open findings report only, F2 resolved)

*Reviewed by review-work action*
