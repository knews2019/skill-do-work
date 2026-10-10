## Review: REQ-656

**Approve with follow-ups**: the revise form is delivered and correct for its main path, but Step 8 still asks for a `pass` verdict that a revise can never get, and the "Steps 2 to 8 apply" rule forces a non-completed-work report (the deploy guide in the UR) into the shipped-work narrative.
Route B | merge range `c44d9f29..6c837fab`

### What's built
- `ai-report revise <dir|latest> [what changed]` is a seven-step subsection of `skills/do-work-toolbox/actions/ai-report.md`. It resolves `latest` from the regenerated catalog and follows `superseded_by` forward, reads the prior bundle as untrusted data, re-checks its claims, writes a `-rev<N>` sibling through Collision-Safe Publication with a rev block and an `ai-report-supersedes` meta, regenerates the catalog, and does not commit.
- A length-keyed table-of-contents rule sits in Report Design Rules. Every bundle-writing run ends with the path plus a `file://` link. The guide, help, routing row and `$ARGUMENTS` line name the form.
- Missing: the Step 8 and checklist verdict rule for a revise (F1). A defined section shape for revising a report that is not a completed-work report (F2).

### Decisions / risks for you
- F2 needs a call: should a revise keep the prior report's own section order when that report is not about completed work? Value: the UR's real cases (a deploy guide, "findings so far") revise cleanly. Risk: one more sentence of rule in the revise form. The proposed fix below says yes.

### Findings

**Important:**
- F1 `ai-report.md:150` (Step 8: "the last `judge.json` verdict is `pass` (or `skipped` or `error` ...)"), `:170` (render checklist, same verdict list) and `ai-report-guide.md:55` ("reruns until the verdict is `pass`") contradict revise step 5 (D-21), which says the prior-bundle link finding is expected. Reproduced in a scratch repo: `ai-report-judge` on a bundle whose only defect is `href="../<prior>/index.html"` returns `verdict=fail`, exit 1, one `AI-REPORT-JUDGE-BROKEN-LINK`. A literal agent cannot complete Step 8 for any revise, or reports a false pass. Fix now, yes. Replacement text:
  - `:150`: "Confirm `index.html` exists, the last `judge.json` verdict is `pass` (or `skipped` or `error` and disclosed in the footer, or for a revise `fail` whose only finding is the expected link to the prior bundle), screenshots open at full resolution, ..." (rest unchanged)
  - `:170`: "... (verdict `pass`, or `skipped` or `error` and disclosed in the footer, or for a revise `fail` with only the expected prior-bundle link finding)."
  - guide `:55`, append: "For a revised report, the link back to the old report is the one expected finding, because the check serves only the new folder."
  — impact-user-visible → report only
- F2 `ai-report.md:53` revise step 5: "Steps 2 to 8 apply to the new bundle as for its kind". No kind is defined in `ai-report.md` at this merge (REQ-654, the kinds request, was cancelled; REQ-687, `--kind proposal|root-cause`, is not merged, and neither covers a deploy guide). So Step 5's required nine-section narrative (Verdict, What Shipped, ...), Step 3's archive-evidence search, Step 6's "verify each claim against the provenance ledger" (built only in Step 1, which revise replaces) and checklist `:169` all apply to the UR's own example, a deploy guide. A literal agent writes "What Shipped" sections into a deploy guide. Fix now, yes. Replacement for the first clause of step 5: "Steps 2 to 8 apply to the new bundle as for its kind (step 4 above names the folder); when the prior bundle is not a completed-work report, keep its own section order in place of Step 5's narrative, and the step 3 re-check serves as the provenance ledger," (rest unchanged). — impact-user-visible → report only

**Minor:**
- F3 `ai-report.md:27` Do-NOT-use bullet "The target is unfinished or unsuccessful; report its status instead of presenting it as shipped." The UR's "update the ai-report with the findings so far" is a revise of a report on unfinished work. This bullet can stop that route before the revise form is read. Fix now, optional. Replacement: "The target is unfinished or unsuccessful; report its status instead of presenting it as shipped (a revise targets an existing report, and its step 2 covers unfinished linked work)." REQ-687 also edits this line, so expect a merge seam. — impact-user-visible → report only
- F4 `ai-report.md:49` `latest` is "the newest bundle of any kind", so it can pick an `architecture-report` or stakeholder-questions bundle. Scratch check: a same-minute `2026-10-10_1600_architecture-report` sorts ahead of `..._deploy-guide-rev1`. A revise of it writes `..._architecture-report-rev1`, which the architecture scanner (`architecture.go:25`, regex `_architecture-report(?:-(\d+))?$`) never sees, and it skips that action's own re-check contract. This follows the REQ's capture assumption ("of any kind"), so no fix now. — impact-negligible → report only

**Nit:**
- F5 `ai-report.md:53` D-21 sentence: "check that the link opens once the bundle is written". In `ai-report` the bundle is already written into its final folder at Step 5, before Step 7, unlike `architecture-report`'s draft directory. So the timing clause gives no instruction. Optional replacement: "..., so rerun until no other finding remains and check on disk that the prior bundle the link names exists." — impact-negligible → report only
- F6 Restatements left stale by REQ-655 and REQ-657, not by this REQ: `ai-report-guide.md:9` "Every invocation creates a fresh timestamped bundle" is false for `index`, `find` and `judge`. `ai-report.md:31` "One invocation covers one UR or one REQ" does not fit a revise, which covers one bundle. Optional replacements: guide `:9` "Every report-writing invocation creates a fresh timestamped bundle under `ai-reports/`:" and `:31` "One report invocation covers one UR or one REQ; a revise covers one existing bundle; blank is the explicit `most recent` form." — impact-negligible → report only

### Anti-bloat count (maintainer request)
Items in the diff that the REQ did not name. Each is small, and none breaks a Constraint:
- B1 Use-when bullet `ai-report.md:21`.
- B2 New checklist line `:172`.
- B3 D-16 parenthetical: N read from metas when the catalog refused.
- B4 "stop at a repeated path" cycle guard in step 4.
- B5 "(step 4 above names the folder)" in step 5.
- B6 Copy the prior `ai-report-kind` meta (D-05).
- B7 "Earlier rev blocks are not carried forward".
- B8 Guide `## Revising a Report` (one paragraph). It restates the TOC and path-last rules, so these are two more places that can drift.
- B9 Second routing phrase `update the report`. The REQ named one phrase "such as revise the report". "update the report" is broad and could catch "update the architecture report".

There is no new file, verb, flag, helper, script, test, queue field or status. Net change: +27/-6 lines in five files.

### Conflict resolutions with REQ-657 (judged)
- `skills/do-work-toolbox/SKILL.md:24`: one `ai-report` row holds the index, find, judge and revise phrases. Correct.
- `ai-report.md` checklist end: REQ-657's `ai-report-judge` render line is kept, the builder's widened exception line is kept, and the revise line is appended. Correct, apart from F1, which this merge seam created.
- `help.md` and the guide Input block: both siblings' lines are kept and columns align (help description at column 34). Correct.
- D-21 (integrator): the reason is right and confirmed by the reproduction above. The fix is incomplete because Step 8 and the checklist were not widened (F1).
- `$ARGUMENTS` line (`ai-report.md:31`): keeps REQ-657's judge, index and find clause and adds the revise clause. Requirement 10 is met.

### Requirements Checklist
- [x] R1 Resolve `<dir|latest>`, read the prior bundle as untrusted data: delivered (steps 1 and 2; chain-forward covers the same-minute tie, confirmed in scratch: rev1 sorts before rev2 and points forward to rev2)
- [x] R2 `yyyy-mm-dd_hhmm_<same-slug>-rev<N>` through Collision-Safe Publication: delivered (step 4; the slug rule strips the date, the time, `-rev<K>` and the collision suffix; the example path starts with `ai-reports/`)
- [x] R3 rev-N block with date, Changed, Still to do: delivered (step 5)
- [x] R4 `ai-report-supersedes` meta in `<head>`: delivered (bare folder name; the catalog matches by base name)
- [x] R5 Catalog regenerated, `superseded_by` forward, prior bytes unchanged: delivered (step 6; confirmed in scratch)
- [x] R6 Readable-style baseline: delivered (full Report Design Rules, no prior CSS that breaks them)
- [x] R7 TOC rule keyed on length, in Report Design Rules: delivered (`ai-report-reference.md:91`)
- [x] R8 Path plus `file://` link as the last line: delivered (Output Format, every bundle-writing run)
- [x] R9 Committing stays out: delivered (step 7)
- [x] R10 Arguments line, guide, routing, help: delivered (the arguments line was fixed by the integrator in `8a56cd5b`). Probe and reference contract pass.
- [x] Constraints: no in-place edit, no queue field or status, no commit or publish, its own release. Met.
- [ ] Revise reaches Step 8 cleanly: partially delivered (F1).

### Acceptance Testing

**Result: Partial** (implementation and integration stages)
- `REQ-656-probe.sh`: exit 0, 1.3 s. `shipped-package-reference-contract.sh`: PASS, 1.2 s. `git diff --check` on the range: clean.
- Scratch repo, four bundles: `ai-report-index` produced the chain original → rev1 → rev2 with correct `supersedes` and `superseded_by`. Same-minute ties sort by path ascending (architecture-report, then rev1, then rev2), so step 1's chain-forward rule is needed and works.
- `ai-report-judge` on a rev bundle whose only defect is the parent link: `verdict=fail`, exit 1, one `AI-REPORT-JUDGE-BROKEN-LINK` → F1.
- Not run: a full agent-driven revise, and `maintainer-verify.sh` (forbidden by the brief; the integrator's run passed at `6c837fab`).

### Suggested Additional Testing
- Deployment: unassessed. Check the staged-skills heavy lane on the range (it is in the REQ's heavy plan) before release.
- Live acceptance: unassessed. In a consumer repo, run `ai-report revise latest "..."` on a non-completed-work report (a deploy guide) and confirm the section shape after the F2 fix.
- Edge case: run revise after the catalog command refuses on a hand-made `ai-reports/index.html`, with an explicit `<dir>`. Confirm N comes from the metas and the path line still prints.

### Scores (on the record, not the headline)

**Overall: 79%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All 10 delivered; Step 8 path for a revise is incomplete (F1) |
| Code Quality | 80% | F1 contradiction, F2 undefined "kind" for non-completed-work reports |
| Test Adequacy | 85% | Probe RED→GREEN plus a literal scratch walkthrough; Step 7 was not run in the walkthrough, which hid F1 |
| Scope | 95% | Five declared files exactly; small undeclared additions B1-B9 |
| Risk | Low | Prose only; the prior-bundle immutability holds |
| Acceptance | Partial | Core flow verified; a revise cannot satisfy Step 8 as written |

(95+80+85+95)/4 = 88.75, minus 10 for Partial = 78.75, shown as 79%.

### Follow-ups created
None (6 findings report only)

## Review

**Overall: 79%** | 2026-10-10T16:37:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 80% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `ai-report.md:150`, `:170`, `ai-report-guide.md:55` still require a final `pass` verdict, but a revise always ends `fail` on the expected prior-bundle link (reproduced: verdict=fail, exit 1); widen Step 8, the checklist and the guide to accept that one finding — impact-user-visible → report only
- F2 `ai-report.md:53` "Steps 2 to 8 apply ... as for its kind": no kind is defined at this merge, so a non-completed-work report (the UR's deploy guide) is forced into Step 5's shipped-work narrative and Step 6's provenance-ledger check; keep the prior section order and use the step 3 re-check as the ledger — impact-user-visible → report only

**Minor findings:** F3 `ai-report.md:27` Do-NOT-use "target is unfinished" bullet can block a revise of a findings-so-far report — impact-user-visible → report only; F4 `latest` of any kind can revise an architecture-report into a `-rev1` folder its scanner never sees (follows the capture assumption) — impact-negligible → report only; F5 (nit) D-21 "once the bundle is written" adds no timing in ai-report — impact-negligible → report only; F6 (nit) guide `:9` "Every invocation creates a fresh bundle" and `ai-report.md:31` "One invocation covers one UR or one REQ" are stale, from REQ-655/657 — impact-negligible → report only
**Acceptance:** Partial — implementation and integration stages: probe and reference contract pass, scratch catalog chain and tie-forward verified; judge returns fail on the expected link, so Step 8 is unreachable for a revise
**Restatement sweep:** redefined the revise input form, `latest`, the path-last Output Format rule, the length-keyed TOC rule, the "only the bundle / only exceptions to searching" boundary (now also revise), and (via D-21) the end verdict of Step 7 for a revise. Stale: F1 (`ai-report.md:150`, `:170`, guide `:55`), F2 (`:53`, `:169`), F3 (`:27`), F6 (guide `:9`, `:31`). Inherited from REQ-655 (input forms and boundary): `:31` now names revise, Output Format still true, guide `:21` widened. Inherited from REQ-657 (judge verdicts): `:144`, `:150`, `:170`, guide `:55` → F1. Checked, not stale: README.md:112, toolbox help.md:14, command-line-guide.md:26 (no new verb), architecture-report.md, stakeholder-report.md (TOC rule applies by length, as intended), completed-work-presentation-reference.md:7 (architecture-report named as an example only)
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*
