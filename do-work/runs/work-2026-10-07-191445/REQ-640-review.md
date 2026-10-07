## Review

**Overall: 91%** | 2026-10-07T20:05:00Z

Mode: pipeline (worktree dispatch, delegated integration). Diff: `git diff 2ab092c0..b36f7c65` (builder commits d7e2e186, 92f662d7; merge b36f7c65). Decisions read from `do-work/runs/work-2026-10-07-191445/REQ-640-handback.md` (D-01 to D-09).

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

### Requirements Checklist

| # | Requirement (REQ Detailed Requirements / UR-138 "REQ B") | Status | Where |
|---|---|---|---|
| 1 | `### Mid-Run Messages (any step)` directly after Step 3.5, condition-keyed | Delivered | work.md:165-173 (after Step 3.5, before Mechanical Evidence-Gate Loop); branches keyed on what the message changes, no message-type list |
| 1a | Progress/state question answered from disk, run does not stop | Delivered | work.md:167, 169 |
| 1b | Not-yet-handed-back change: `## Addendum (mid-run)`, framing line, Outside-text containment; serial writes now; delegated: brief then integrator writes before merge; forward or say builder will not see it; takeover keeps it | Delivered | work.md:170 (plus one-section-with-entries rule, D-03) |
| 1c | Handed-back change: capture in-flight path (`addendum_to`), reported as queued | Delivered | work.md:171 (includes running integrator, D-02) |
| 1d | Queued REQ or new work: capture paths in the writing gap | Delivered | work.md:172 |
| 1e | Hand-back emphasis: run manifest, Decision Brief reads it | Delivered (narrowed) | work.md:173; work-reference.md:480 (manifest row), 848 (Decision Brief); no-run-directory case held in session (D-01); see F3 |
| 1f | User's own words, never paraphrase; route to every affected REQ | Delivered | work.md:167 |
| 2 | review-work Step 5 item 1 reads any `## Addendum` | Delivered | review-work.md:87; Step 2 line 55 swept too (D-06) |
| 3 | capture.md `working/` row keeps NEVER modify + pointer | Delivered | capture.md:124 (arrow-form citation resolves) |
| 4 | Orchestrator Checklist line | Delivered | work.md:494, exact wording |
| C | Every work-reference.md heading unchanged; contract tests pass | Delivered | no heading lines in diff; tests below |
| C | Prose only, no Go/flags/new actions/SKILL.md rows | Delivered | diff --stat: 5 .md files |

### Mechanism claims checked against Go

- Takeover reset keeps the section: TRUE. `generatedRecoveryHeading` (tools/do-work-cli/internal/requeststate/state_apply.go:981) lists only Triage … Timing; line 992 keeps every other section.
- `advance` refuses a duplicated section name: TRUE. `advanceSections` (internal/lifecycleadvance/advance_commands.go:344-360) returns "duplicate lifecycle section <name>" for any visible `## ` name seen twice; line 133-135 turns it into ADVANCE-EVIDENCE-MISSING. `## Addendum (mid-run)` and capture's `## Addendum (<date>)` are distinct names, so they coexist.
- Classifier ignores unknown sections: TRUE (lines 138-259 test named lifecycle sections only).
- Outside-text containment keeps user words from opening a section: TRUE, though the reason is the `> ` prefix, not fence-skipping. In `VisibleSections` (internal/requestmodel/visible_sections.go) a `> ` line has a leading run of one `>` (no fence) and its content does not start with `## `, so no quoted line becomes a section. Shipped text cites clarify.md Step 4, which states the reason correctly; only the Exploration note was imprecise.

### Findings

- **F1 (Important)** work.md:173 — The emphasis-note branch says "write it as the hand-back emphasis note in the run manifest, not in memory" with no writing-gap clause, while the three sibling branches (170-172) each say the coordinator waits for its next writing gap. Under delegated integration a run directory always exists, so a literal reader edits `manifest.md` while an integrator runs, which `work-reference.md:470` forbids ("no manifest edits"). Fix: add "Under delegated integration the coordinator writes it in its next writing gap and keeps it in its own session scratch until then." — impact-rule-change → report only
- **F2 (Minor)** work.md:170 — The not-handed-back branch names only "the serial loop" (writes now) and "delegated integration" (brief, integrator writes). Fan-out where the main session integrates itself is in neither case. Fix: key it on the condition, e.g. "When the session that received the message is the only writer under the project root, it writes the section now; under delegated integration …". — impact-rule-change → report only
- **F3 (Minor)** work-reference.md:848 — The Decision Brief reads the emphasis note only to "order what each section shows; the sections … stay the same". The REQ's own RED example ("put decisions first in the update") asks for section order, which this cannot honor, and the narrowing is not recorded as a D-XX. Fix: either state that a note may also move DECISIONS FOR YOU first (it is still "what needs the user", not a self-grade), or record the narrowing as a decision. — impact-rule-change → report only
- **F4 (Minor, restatement)** work-reference.md:310 and docs/work-guide.md:130 — Both say takeover strips the claim's "orchestrator sections". `## Addendum (mid-run)` is written by the orchestrator, yet the reset keeps it (work.md:170 says so correctly). A reader of the recovery text can conclude the addendum is lost. Fix: say "generated sections" (the Go term) in both places, or add "user-intent sections such as `## Addendum (mid-run)` survive". — impact-rule-change → report only

### Restatement Sweep

Redefined elements: (a) who may write a working REQ mid-run, (b) the requirement-source list, (c) the run manifest's contents and Decision Brief inputs.
- (a) capture.md:59, 63 (NEVER write to working/, Immutability Rule): capture-scoped, still true; pointer added only in Step 2 table (D-07). work-reference.md:51 exempts the work pipeline. work.md:155 (escalated answer) consistent. work-reference.md:470 already carries "any mid-run addendum text" in the brief and the one-writer rule; kept by work.md:170-172, broken only by F1. work-reference.md:310 and work-guide.md:130 imprecise (F4).
- (b) work.md:202, review-work.md:55, 87 all updated; no other "What/Detailed Requirements" reader in skills/ (verify-requests.md and toolbox presentation reference are unrelated).
- (c) work-reference.md:480 manifest row updated; work.md:536 Progress Reporting restates section order and stays consistent with the within-section narrowing (see F3); crew-members/background-agents.md manifest text is generic and consistent. Go reads `manifest.md` only in cleanup (`manifestIsConsumed`, first `status:` line); a note below the header does not interfere.
- docs/work-guide.md:72, 101 and the new 136 paragraph consistent; clarify.md:99, 259 consistent; docs/capture-guide.md:76 still true for capture.

### Acceptance Testing

Prose read of each changed section against the REQ, mechanism claims checked in Go (above). `bash _dev/tests/shipped-package-reference-contract.sh` PASS; `bash _dev/tests/contract-regressions.sh` PASS; `git diff --check 2ab092c0..b36f7c65` clean. RED/GREEN: the new subsection names the in-flight path and review reads `## Addendum` — GREEN.

### Scope

Five files, all in `write_set`; work-guide.md extension recorded as D-09, the Step 2 and plan-validation sweeps as D-05/D-06. No heading changed, no out-of-scope edits.

**Important findings:**
- F1: emphasis-note branch lacks the writing-gap clause, so it directs a manifest edit while an integrator runs — impact-rule-change → report only

**Minor findings:** F2 (fan-out without delegation not covered) — impact-rule-change → report only; F3 (emphasis cannot reorder sections; undocumented narrowing) — impact-rule-change → report only; F4 (recovery text says "orchestrator sections stripped") — impact-rule-change → report only
**Acceptance:** Pass — every requirement delivered; contract tests green; Go claims verified
**Suggested testing:** 0 items (prose-only; no behavior to pin)
**Follow-ups created:** None (4 findings report only)

**Verdict:** Approve

*Reviewed by review-work action*

## Re-check

Range: `git diff b36f7c65..33c41c50` (fix 9079ab68, merged as 33c41c50). Contract tests re-run: shipped-package-reference-contract.sh PASS, contract-regressions.sh PASS; `git diff --check` clean.

- **F1: resolved.** work.md:173 now says the coordinator writes the emphasis note in its next writing gap and keeps it in session scratch until then. The one-writer rule now holds on all four branches.
- **F2: resolved.** work.md:170 is keyed on "the session that received the message is the only writer under the project root", and names the serial loop and a self-integrating fan-out run as examples.
- **F3: resolved.** work-reference.md:848 lets the note order the sections as well as their contents, so the RED example ("decisions first") can be honored.
- **F4: still open, report only.** The integrator's reason is accepted: work.md:170 states the rule correctly.
- **F5 (Minor, new, restatement)** work.md:536 (Progress Reporting) still gives the Decision Brief section order as fixed ("lead with WHAT'S BEING BUILT … then DECISIONS FOR YOU …") and does not mention the emphasis note that F3's fix now lets reorder sections. Fix: add "unless the run manifest's hand-back emphasis note orders them otherwise". — impact-rule-change → report only

Overall after re-check: **93%** (Requirements 97, Code Quality 92, Test Adequacy 85, Scope 95; Risk Low; Acceptance Pass). Verdict: **Approve**. Follow-ups created: None (2 findings report only).

## Re-check 2

Range: `git diff 33c41c50..0982069e` (fix 5c68b501, merged as 0982069e). One line, work.md:536. shipped-package-reference-contract.sh PASS, contract-regressions.sh PASS, `git diff --check` clean.

- **F5: resolved.** work.md:536 now says the section order holds "unless the run manifest's hand-back emphasis note orders them otherwise". The `(**Mid-Run Messages (any step)**, above)` citation points to work.md:165, which is above this line.
- **Nothing new.** F4 stays the only open finding (Minor, report only).

Final: **94%** (Requirements 97, Code Quality 94, Test Adequacy 85, Scope 95; Risk Low; Acceptance Pass). Verdict: **Approve**. Follow-ups created: None (1 finding report only).
