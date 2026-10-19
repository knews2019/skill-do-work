## Review: REQ-701

**Approve** — a finished second-opinion page now survives a failed model check, and a subagent's model beside the requested one is recorded instead of treated as ambiguity. Ship-ready.
Route A | merge 5b1536ac (range 6256fd58..5b1536ac)

### What's built
- `ai-report second-opinion` keeps the staged page on every refusal after the backend wrote it; the failure reason names the kept file. Only a zero-byte or already-published staging file is removed.
- A model check that fails after the page is written publishes the pair as `index-<cli>-<model>-unverified[-N].html` / `<cli>-<model>-unverified[-N]-reply.md`, state `unverified`, no backlink check, no judge, no banner, with finding `AI-REPORT-SECOND-OPINION-UNVERIFIED` naming the requested model and every model seen.
- For `claude` and `gemini`, the requested model wins when the native usage map lists it; extras land in the new `seen_models` row field (always a JSON array). Lives in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go`.

### Decisions / risks for you
- None. The one behavior change (a model switch now publishes an unverified pair instead of refusing) is what the REQ asked for; the REQ-698 test that locked the old rule was renamed and flipped, traced in the Testing section.

### Findings

**Important:**
- None

**Minor:**
- M1 `report_second_opinion.go` ~435-466: the unverified branch never checks that a finished page exists. A backend that wrote nothing (zero-byte staging) and reported a mismatched or multi-model map now publishes an empty `index-<cli>-<model>-unverified.html` and reply with state `unverified` and a reason that blames the model; before this diff the same run was a plain `failed`. The verified branch catches this indirectly through the backlink check, which the unverified branch skips on purpose. A `len(page) == 0` guard before `secondOpinionNames` would keep that run on the `failed` path with the staging file removed by the size rule. — impact-user-visible → report only
- M2 `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:5` (the `native-backend-identity` bullet from REQ-698): "refuse ambiguous or mismatched identity before naming artifacts" now restates the old rule; after this diff mismatched identity publishes under an explicit `-unverified` name and only ambiguity between the requested model and a subagent's is no longer a refusal at all. The hand-back's proposed REQ-701 bullet beside it reconciles the two once the lessons phase lands it. — impact-negligible → report only

**Nit:**
- N1 D-01 put requested-model-wins in the shared gemini/claude branch, but only `claude` fakes exercise it; a one-case gemini fixture with `--model` plus a second `stats.models` key would pin the shared path. — impact-negligible → report only
- N2 The run summary line prints `model=` empty for an unverified row; the requested model is only in the reason text. Readable, but a reader scanning `model=` sees nothing. — impact-negligible → report only

### Requirements Checklist

- [x] No code path in `secondOpinionRun` deletes a finished staged page; cleanup limited to the never-written case — delivered (deferred closure keyed on `stagePublished` or `Size() == 0`; stat failure removes nothing and appends nothing; the missing-backlink test asserts the kept non-empty file and the "staged page kept at" reason)
- [x] Verification failure after the page exists publishes it under an unverified name, no banner, typed finding with requested and seen models — delivered (`TestSecondOpinionUnverifiedModelKeepsFinishedPage` asserts page bytes, reply bytes, `banner_state: not-inserted`, `judge_verdict: not-run`, unchanged `index.html`, finding code and all three model names)
- [x] `claude`: requested `--model` wins when in `modelUsage`, else the single entry; extras recorded, not ambiguity; only empty or multi-without-requested fails — delivered (`TestSecondOpinionSubagentModelIsNotAmbiguous`; the empty map falls to the sentinel because `len(seen) == 0`)
- [x] Both behaviors locked in `report_second_opinion_test.go` — delivered (two new tests plus the renamed codex mismatch test)
- [x] UR: a 12-minute generation is never discarded on a probe verdict — delivered (M1 is the one gap: a page that was never written is still published, which is the opposite failure and harmless to finished work)

Builder claims checked against the diff: every file and symbol the hand-back manifest names is in the merge range; the Implementation Summary's description of the defer rule, the sentinel, `seen_models` and the two doc sentences matches the code. Red-green: the three tests the Testing section lists trace to `## Red-Green Proof` (RED case = `TestSecondOpinionSubagentModelIsNotAmbiguous`, GREEN second case = `TestSecondOpinionUnverifiedModelKeepsFinishedPage`); RED ran on a signature-only scaffold so the failures were assertion failures (D-02), which is valid RED-before-GREEN. Cross-REQ update to the REQ-698 test is traced. P-A-U: all three boxes ticked. Decisions D-01 to D-05 are in the REQ (copied from the hand-back) and match the diff.

Correctness points asked for, each traced in the merged source:
- Deferred cleanup: cannot delete a finished page (removal needs `stagePublished` or size 0; `stagePublished` flips only after the page bytes were written to the sibling). Cannot leave an empty staging file (size 0 is removed; stat failure means it is already gone). Defer order is correct: the cleanup closure runs before `bundleRoot.Close()`.
- The unverified branch skipping backlink and empty-reply checks is deliberate and documented in the comment; the cost is M1.
- `row.ActualModel = ""` while `requested_model` keeps `options.model`: correct, and the reason string carries the probe-resolved model even when `--model` was not passed.
- `secondOpinionMetadata` requested-wins rule: `expected != ""` guards the map lookup; a nil map (`"modelUsage": null` or absent) has length 0 and returns the sentinel with `seen == nil`; `errors.Is` matches through the `%w` wrap.
- `seen` for codex/agy derived after the switch (`if seen == nil && model != ""`): one entry when a model was read, nil on identity errors; the slug-pattern failure still returns nil seen as a non-identity error.
- Probe call passes `model` (possibly empty) as `expected`: empty falls to the single-entry rule; a sentinel on the probe is a plain `fail` with no page yet, which is right.
- JSON `seen_models` never null: row literal initializes `[]string{}` and the generation path uses `append([]string{}, seen...)`.

Anti-bloat count: the hand-back's claim is accurate. Named items present: `errSecondOpinionModelUnverified`, `SeenModels`/`seen_models`, the `expected` parameter and `seen` return, `stagePublished`, `secondOpinionFakeClaude`, two new tests, one renamed test, the kept-staging assertion. Unnamed additions: the `unverified` and `publishedModel` locals, the `fmt` test import, and a `names` slice local inside the renamed test. Helpers, options, flags, files or decorative tests the REQ did not name: none.

### Acceptance Testing

**Result: Pass** (stages exercised: Implementation and Integration)
- `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/` at HEAD 5b1536ac: 19 of 19 `--- PASS`, `ok` in 15.6 s. This runs the merged tree, so it covers Integration for this package; the integrator's full gate at 5b1536ac is on the record in the REQ's Testing section and was not rerun.
- Original bug: the RED case (two-entry `modelUsage` holding the requested model) now ends with state `ran`, the page under the requested model's name and no staging leftover; the envelope lacking the requested model ends with the page kept under `-unverified`.
- Regression: the pre-existing identity, preservation, judge, batch and cancellation tests in the same family pass unchanged.

### Suggested Additional Testing

- Live acceptance: unassessed. Run `ai-report second-opinion --backend claude --model <session-model>` on a real bundle where the editor dispatches an Explore subagent, and confirm the page lands as `index-claude-<session-model>.html` with both models in `seen_models`.
- Live acceptance for gemini: unassessed. Run `--backend gemini --model <slug>` once and confirm the `stats.models` key equals the `--model` slug; the shared requested-wins rule depends on that key equality (N1).
- Deployment: unassessed. The change reaches consumers only through the next do-work-cli release; confirm the launcher's fetched binary reports state `unverified` on a forced mismatch.

### Scores (on the record — not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All four REQ requirements and the UR's intent delivered |
| Code Quality | 90% | Clean, comments explain the rule; M1 empty-page gap on the unverified branch |
| Test Adequacy | 90% | Red-green traced to the proof; gemini shared-path and empty-page cases not pinned |
| Scope | 100% | Five files, all named by the REQ or holding the contract it changed; no unnamed surface |
| Risk | Low | New `state` value `unverified` for row consumers; `secondOpinionMetadata` signature changed, both callers updated; no security or data-loss path |
| Acceptance | Pass | Implementation and Integration stages via the merged-tree test family |

### Follow-ups created
- None (4 findings report only)

## Review

**Overall: 95%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** M1 `report_second_opinion.go` unverified branch publishes a zero-byte page and reply when the backend wrote nothing and the model check failed, where the old path was a plain `failed` — impact-user-visible → report only; M2 `lessons-do-work-cli.md:5` (`native-backend-identity`, REQ-698) still says mismatched identity is refused before naming, now reconciled only by the REQ-701 bullet the hand-back proposes — impact-negligible → report only; N1 no gemini case pins the shared requested-model-wins branch (D-01) — impact-negligible → report only; N2 summary line prints `model=` empty for unverified rows — impact-negligible → report only. Anti-bloat: no helpers, options, flags, files or decorative tests beyond the named items (unnamed: locals `unverified`, `publishedModel`, test local `names`, `fmt` test import).
**Acceptance:** Pass — Implementation and Integration stages: `go test -run '^TestSecondOpinion'` at 5b1536ac, 19 of 19 pass; full gate per the Testing section not rerun
**Restatement sweep:** redefined the published filename contract (`index-<cli>-<actual-model>[-N].html` gains the `-unverified` form), the run row schema (`seen_models`, state `unverified`), the `AI-REPORT-SECOND-OPINION-<STATE>` code set (`-UNVERIFIED`), and the identity rule (mismatch publishes unverified; a subagent model is not ambiguity); swept `skills/` and `_dev/`: `actions/ai-report.md` and `docs/ai-report-guide.md` updated in this diff, `ai-report-reference.md`, `SKILL.md`, `do-work-cli.sh`, `_dev/tests/do-work-cli-launcher-behavior.sh` and the 0.306.0 changelog entry hold no stale restatement ("never substitutes" still holds); one stale restatement in `lessons-do-work-cli.md:5` recorded as M2 | unread for the sweep: REQ-700
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

### Delta re-review (5b1536ac..e8b063c1)

M1 is resolved. Read-only delta review (no tests run, per the integrator; family tests were green in the worktree). The guard `if unverified && len(page) == 0` sits right after `os.ReadFile` of the staging page and before the backlink check, so an unverifiable model with no page written is now a plain `failed` row whose reason names the requested and seen models; `stagePublished` stays false and the zero-byte staging file is removed by the existing size rule. The verified path is untouched (an empty page there still fails on the backlink check). No new defect: `row.SeenModels` is already set when the guard runs, `fail` overwrites the reason as elsewhere, and no other branch can reach `secondOpinionNames` with an empty page. No bloat: one 5-line guard with a comment and one pinning test `TestSecondOpinionUnverifiedModelWithoutPageFails` that reuses `secondOpinionFakeClaude`, strips its page write with an exact-string replace that fails loudly if the fake's line changes, and asserts state `failed`, the reason text, and a bundle holding only `index.html`. No new helpers, options, flags or files.

Scores change: Code Quality 90% → 95% (M1 closed), Test Adequacy 90% → 95% (the empty-page case is pinned). Requirements, Scope, Risk Low and Acceptance Pass unchanged. **Overall: 95% → 97%** (mean 97.5, not rounded up). Findings now: M2, N1, N2 (3 findings report only). The `## Review` block above records the 5b1536ac state; the integrator carries the delta values into the REQ.
