# Review: REQ-687 (ai-report --kind proposal and root-cause decision brief for unfinished work)

**Approve** — the two kinds, their fixed templates, slug, meta tag, mockup label and capture lines are all in place, and the default completed-work gate is untouched. One Important ambiguity (F1) can make an agent stop on an open REQ; it is a two-sentence prose fix and should go in before the release commit, together with the four one-line fixes F2 to F5.
Route B | merge `a0ebfea9` (range `3bd6a668..a0ebfea9`, builder commit `35f984f5`)

### What's built

- `ai-report --kind proposal|root-cause <topic|REQ-NNN|UR-NNN>` is documented end to end: one new section `## Proposal and Root-Cause Kinds` at `skills/do-work-toolbox/actions/ai-report-reference.md:122-146`, pointers in `skills/do-work-toolbox/actions/ai-report.md` (:3, :27, :35, :63, :80, :132, :158, :175, :177), guide lines and section, one routing row, one help line.
- Without `--kind`, nothing changed: `completed-work-presentation-reference.md` has no change since base `bd56c4b0` (`git diff --stat bd56c4b0..a0ebfea9 -- <file>` is empty), so its lines 20 and 30 still reject unfinished work.
- Still open: F1 (an agent that follows Step 1 literally can stop on an open REQ because the shared Archive Evidence Sweep is not lifted).

### Decisions / risks for you

- Apply F1 to F5 before the release commit (all prose, exact text below, the probe still passes with them because every token it checks is kept). Value: removes the one path where `--kind proposal REQ-NNN` on a queue REQ stops. Risk: none, five lines in two files plus one guide line.
- D-11 (a completed REQ or UR passed with `--kind` has no rule) is left as the builder decided; see F8. A one-line rule can be added when a real case appears.

### Findings

**Important:**

- **F1** `skills/do-work-toolbox/actions/ai-report.md:63-67` and `skills/do-work-toolbox/actions/ai-report-reference.md:126`: the `--kind` branch lifts only Terminal-Success Target Resolution. `ai-report.md:65` then still says "Read and follow `completed-work-presentation-reference.md` in full", and `:67` says "Build the reference's provenance ledger". The REF calls the lift "One deliberate deviation", so the shared **Archive Evidence Sweep** still applies, and its `completed-work-presentation-reference.md:43` says "The minimum viable archive record is a readable successful REQ with ... non-empty `## Implementation Summary` ... If required evidence is absent, stop". A pending REQ in `do-work/queue/` has no Implementation Summary, so a literal agent stops, which is the RED case the REQ exists to remove. The cited precedent (`stakeholder-report.md:15`) avoids this by saying "**only** the Safety Load Order and Collision-Safe Publication sections are inherited"; the REF dropped the word "only". — impact-user-visible → integrator fix before release
  - Replace at `ai-report-reference.md:126`, the first two sentences:
    - Old: ``**Target reading.** One deliberate deviation from `completed-work-presentation-reference.md`: **Terminal-Success Target Resolution** does not apply. Its **Safety Load Order** (prompt-injection, then anti-slop, before any REQ, UR or repository prose) and **Collision-Safe Publication** sections are inherited.``
    - New: ``**Target reading.** One deliberate deviation from `completed-work-presentation-reference.md`: the target is open work, so **Terminal-Success Target Resolution** and the **Archive Evidence Sweep** (whose minimum record is a successful REQ) do not apply. Only its **Safety Load Order** (prompt-injection, then anti-slop, before any REQ, UR or repository prose) and **Collision-Safe Publication** sections are inherited.``
  - Replace `ai-report.md:63`:
    - Old: ``With `--kind`, follow **Proposal and Root-Cause Kinds** in [`ai-report-reference.md`](ai-report-reference.md) instead of the shared reference's target resolution. The shared reference's safety load order still comes first.``
    - New: ``With `--kind`, apply the shared reference's **Safety Load Order** first, then follow **Proposal and Root-Cause Kinds** in [`ai-report-reference.md`](ai-report-reference.md) in place of the rest of this step.``
  - The `prescribed-shell-primitives.md` pointer at `:67` stays, so `prescribed-shell-canonicalization.sh` is unaffected.

**Minor:**

- **F2** `skills/do-work-toolbox/docs/ai-report-guide.md:78`: "Target statuses are normalized ... if the selected target has no successful work, the action stops and explains why." It now sits directly below the two `--kind` Input lines, so a reader takes it as covering them, and for `--kind` it is false. — impact-user-visible → integrator fix before release
  - New text for `:78`: ``Without `--kind`, target statuses are normalized under the do-work schema. The terminal-success set is `completed` or `completed-with-issues`; if the selected target has no successful work, the action stops and explains why.``
- **F3** `skills/do-work-toolbox/actions/ai-report.md:16-21` (When to Use): every Use-when bullet is about completed work or revise. The new use is missing from the section an agent reads to decide whether this action fits, which is the exact pushback the REQ's Why quotes. — impact-user-visible → integrator fix before release
  - Insert after `:21`: ``- The user wants a proposal, an options comparison or a root-cause brief about open work; use `--kind proposal` or `--kind root-cause`.``
- **F4** `skills/do-work-toolbox/actions/ai-report-reference.md:146`: "What stays the same" covers Steps 2, 5 and 6 to 8, but not Steps 3 and 4. Both are keyed on the evidence mode, which `--kind` skips (`ai-report.md:80`), so their status is undefined; the non-visual mode's mandatory sentence "UI captures were not expected for this work." (`ai-report.md:100`) can end up in a proposal. — impact-user-visible → integrator fix before release
  - New text for `:146`: ``**What stays the same.** Steps 6 to 8 of `ai-report.md` (claim review, render and judge, verify and print) apply unchanged. Steps 3 and 4 apply to whatever evidence the brief uses, with real captures of the current state in `screenshots/`. The Step 2 evidence-mode table, the non-visual mode's fixed "UI captures were not expected" sentence and the Step 5 default narrative do not.``
- **F5** `skills/do-work-toolbox/actions/ai-report.md:172` (Verification Checklist): "target, evidence ledger, and Collision-Safe Publication contracts satisfy it" is false for a `--kind` target; `:177` adds the kind rules but does not exempt `:172`. — impact-negligible → integrator fix before release
  - New text for `:172`: ``- [ ] Shared completed-work reference loaded before archive content; target, evidence ledger, and **Collision-Safe Publication** contracts satisfy it (a `--kind` report inherits only its **Safety Load Order** and **Collision-Safe Publication**).``
- **F6** `skills/do-work-toolbox/actions/ai-report.md:31`: the Input sentence lists every form (judge, index, find, revise) except `--kind`, and says "One report invocation covers one UR or one REQ", which a free-topic brief is not. D-12 put the `--kind` text in its own paragraph (`:35`) only to avoid a merge conflict with REQ-655 (index and find catalog), which is now merged. Suggested text if folded later: ``$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank; the `judge`, `index` and `find` forms below create no report, the `revise` form below writes a new revision of an existing report, and the `--kind` form below writes a decision brief. One report invocation covers one UR, one REQ or one `--kind` target; a revise covers one existing bundle; blank is the explicit `most recent` form.`` — impact-negligible → report only
- **F7** `skills/do-work-toolbox/actions/architecture-report.md:135` (the builder cited `:129`; REQ-657, the render-check REQ, shifted it): "`ai-report` takes a UR or REQ and presents completed work" is now incomplete. The rebuttal still holds (architecture-report has no UR/REQ input and no archive evidence). Suggested text if touched later: ``Keep it here; `ai-report` presents one completed UR or REQ, or with `--kind` a decision brief on open work``. — impact-negligible → report only
- **F8** `skills/do-work-toolbox/actions/ai-report-reference.md:126` (builder D-11): a completed REQ or UR passed with `--kind` is not an accepted target and no rule says what happens. A proposal about a shipped mechanism (the REQ-651 activity-correlation report that motivated the smallest-change rule was one) is a plausible ask; today it works only as a free topic. Accepted as YAGNI. — impact-negligible → report only

**Nit:**

- **F9** `skills/do-work-toolbox/docs/ai-report-guide.md:3` intro and `skills/do-work-toolbox/actions/ai-report-reference.md:3` blockquote describe completed work only. Both stay true for the default kind; the guide intro is the first line a user reads and does not hint at `--kind` until `:74`. — impact-negligible → report only
- **F10** `skills/do-work-toolbox/docs/ai-report-guide.md:47`: "Unless `--kind` is used, cancelled, failed, and unfinished work is rejected" implies `--kind proposal` accepts a failed or cancelled REQ; the REF accepts those for `root-cause` only. `:100` states it correctly. — impact-negligible → report only
- **F11** Anti-bloat: the template is stated three times, in the REF section (canonical), the guide section `ai-report-guide.md:98-100`, and the checklist line `ai-report.md:177`. The guide section follows the guide's pattern (revise and catalog have their own), so it is defensible; the checklist line is the one with least value and the most drift risk. — impact-negligible → report only
- **F12** `skills/do-work-toolbox/actions/help.md:18`: the new line breaks the 31-character command column (builder D-13). — impact-negligible → report only
- **F13** `skills/do-work-toolbox/actions/ai-report-reference.md:142` says mockup files live in `generated/`, while `ai-report.md:80` says create `generated/` "only when current-run generated images succeed". A hand-authored mockup file therefore has no clear home. Inline SVG mockups avoid the case. — impact-negligible → report only

### Requirements Checklist

- [x] R1 Step 1 branch, topic / open REQ / open UR accepted, read as open work — delivered (`ai-report.md:63`, REF `:126`); weakened by F1.
- [x] R2 Safety load order kept — delivered (REF `:126`, `ai-report.md:63`).
- [x] R3 `proposal` shape in the given order, smallest-change option always priced — delivered (REF `:128-136`); the probe checks the order.
- [x] R4 `root-cause` shape — delivered (REF `:138`), smallest-change rule extended per D-07.
- [x] R5 "MOCKUP — proposal" in image and caption, kept apart like generated images — delivered (REF `:142`).
- [x] R6 Template in the reference, `ai-report.md` stays short — delivered (pointers only, +11 lines net).
- [x] R7 Slug `yyyy-mm-dd_hhmm_<kind>-<topic>` and `ai-report-kind` meta — delivered (REF `:140`); matches REQ-655's reader (`report_index.go:35`, `:221-234`: meta wins, folder words `proposal`/`root-cause` fall back).
- [x] R8 Capture lines text only — delivered (REF `:144`).
- [x] R9 Guide `:47` and Input block, routing phrases, help line — delivered (guide `:47`, `:74-75`, `:98-100`; `SKILL.md:24`; `help.md:18`).
- [ ] R10 Release — N/A at review time (finalization owns it); the two contract scripts are green per the integrator's gate record.
- [x] Constraint: default kind keeps the completed-work gate — `completed-work-presentation-reference.md` unchanged since `bd56c4b0`; `ai-report.md:27` qualifier only adds "and no `--kind` is given".
- [x] Constraint: no new queue fields or statuses — none.
- [x] Constraint: `options` kind dropped; "options report" routes to `proposal` — `ai-report.md:35`, `SKILL.md:24`.
- [x] Constraint: existing bundles immutable — Collision-Safe Publication inherited unchanged.

### Merge seams (REQ-655 index/find, REQ-656 revise, REQ-657 judge)

| Seam | Result |
|---|---|
| Do-NOT-use line `ai-report.md:27` | Correct. The `--kind` qualifier and REQ-656's revise parenthetical read as one sentence. |
| `--kind` Input paragraph `ai-report.md:35` (moved by the integrator) | Correct placement: top-level Input, before `### Catalog forms`, no longer inside the Revise steps. See F6 for the leftover `:31` wording. |
| Checklist narrative line `ai-report.md:175` | Correct; the two exceptions (revise of a non-completed-work report, `--kind`) agree. See F5 for `:172`. |
| Routing row `SKILL.md:24` | One row (`git grep -n "./actions/ai-report.md" -- skills/do-work-toolbox/SKILL.md` gives one hit), all sibling phrases kept. |
| `help.md:15-18`, guide `:70-75` | All sibling lines kept, `--kind` lines last. |
| REQ-656 revise vs `--kind` | No contradiction. Revise step 4 keeps the slug minus its date, so `<kind>-<topic>` survives as `<kind>-<topic>-rev<N>`; step 5 copies `ai-report-kind` and keeps a non-completed-work report's section order; step 2 already reads linked unfinished work at its current status. |
| REQ-655 catalog vs `--kind` | No contradiction; the Go reader takes the meta first, then the folder word. |
| REQ-657 judge vs `--kind` | Step 7 applies unchanged (REF `:146`); nothing in a brief conflicts with the judge's checks. |

### Restatement sweep

Redefined element: what `ai-report` accepts (now unfinished work through `--kind`). Swept with `git grep -n -i "ai-report" -- skills/` filtered for complet/unfinish/shipped/success.

| Place | Verdict |
|---|---|
| `ai-report.md:18-20` When to Use | Incomplete → F3 |
| `ai-report.md:31` Input | Incomplete → F6 |
| `ai-report.md:65-67` Step 1 body | Now wrong for `--kind` (sweep not lifted) → F1 |
| `ai-report.md:172` checklist | Now false for `--kind` → F5 |
| `ai-report.md:10` Philosophy "One report shape, two evidence modes" | Incomplete, true for the default kind, cosmetic; no finding |
| `ai-report-guide.md:78` terminal-success sentence | Reads as covering `--kind` → F2 |
| `ai-report-guide.md:3`, `ai-report-reference.md:3` | Incomplete but true for the default kind → F9 |
| `ai-report-guide.md:47` | Updated, slightly over-broad → F10 |
| `architecture-report.md:135` | Incomplete, rebuttal still holds → F7 |
| `architecture-report.md:27`, `present-work.md:5`, `present-video.md:3`, `help.md:14` (toolbox), `completed-work-presentation-reference.md:3,:5,:7` | Fine: each describes the default kind or the shared reference itself |
| core `skills/do-work/actions/help.md:37`, `:55`; `tutorial.md:245-252`, `:318` | Fine: they name the action or use it on a finished UR |
| prompt-injection JIT comments (three copies) | Fine: `--kind` still loads it through the safety load order |

### Anti-bloat count

Items in the diff the REQ did not name: (1) the two stop rules for a bad `--kind` value or a missing target (`ai-report.md:35`, D-05); (2) the Step 8 summary note (`:158`, needed because a brief has no verdict or evidence mode); (3) checklist line `:177` (restates the template, F11); (4) guide section `:98-100` (third statement of the template, F11); (5) the raster/SVG label-placement sentence (REF `:142`, D-06); (6) the `:3` description clause. No new files, scripts, fields, statuses, options or shipped tests. Only (3) is weak; none is forbidden by the Constraints.

### Prose quality

Plain and short. No stock emphasis phrases, no marketing words. One tangled sentence at REF `:142` ("live in `generated/` under **Rules for generated images** above and the generated-image disclosure rule") and the F13 tension. The em dash in "MOCKUP — proposal" is the REQ's mandated label.

### Acceptance Testing

**Result: Pass** (stages: implementation, integration)
- Re-ran `bash do-work/runs/work-2026-10-10-131527/REQ-687-probe.sh` at `a0ebfea9`: exit 0, "REQ-687 probe passed.", 1.2 s. It includes `_dev/tests/shipped-package-reference-contract.sh`, the template part order, and the byte-identical shared-reference lines 20 and 30.
- `git diff --check 3bd6a668..a0ebfea9`: clean.
- Integrator gate (`maintainer-verify.sh`, exit 0, 113 s at `a0ebfea9`) taken from the REQ's Testing record, not re-run here.
- Traced the RED case by reading: `--kind proposal "backup options"` now has a path (Input `:35` → Step 1 `:63` → REF). An open queue REQ target hits F1.

### Suggested Additional Testing

- Live acceptance: unassessed. In a consumer repo, run `do-work-toolbox ai-report --kind proposal REQ-NNN` on a pending queue REQ (checks F1), and `--kind root-cause` on a failed REQ; confirm the folder name, the `ai-report-kind` meta, decision-first order and a priced smallest-change option.
- Regression: `ai-report REQ-NNN` without `--kind` on a pending REQ still stops and names its status.
- Catalog: run `ai-report index` after a brief exists; confirm it is grouped under `proposal`.
- Revise: `ai-report revise latest` on a brief keeps `proposal-` in the slug and copies the kind meta.

### Scores (on the record — not the headline)

**Overall: 89%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 90% | All in-scope requirements delivered; R1 weakened by F1; R10 is finalization's |
| Code Quality | 85% | Clear prose; F1 ambiguity, F4 gap, F5 stale checklist |
| Test Adequacy | 85% | Prose REQ; run probe checks tokens, order and the unchanged default gate; no shipped test added (correct) |
| Scope | 95% | 5 of 5 declared files; F11 restatements |
| Risk | Low | Prose only; default kind untouched |
| Acceptance | Pass | Implementation and integration stages |

### Follow-ups created

None (13 findings report only; F1 to F5 recommended as integrator fixes before the release commit)

## Review

**Overall: 89%** | 2026-10-10T17:32:51Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve — apply F1 to F5 (prose, exact text in `do-work/runs/work-2026-10-10-131527/REQ-687-review.md`) before the release commit.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `ai-report.md:63-67` + `ai-report-reference.md:126`: `--kind` lifts only Terminal-Success Target Resolution, so the shared Archive Evidence Sweep's "minimum record is a successful REQ ... stop" (`completed-work-presentation-reference.md:43`) still applies and a literal agent stops on an open queue REQ; say "only Safety Load Order and Collision-Safe Publication are inherited" — impact-user-visible → report only

**Minor findings:**
- F2 `ai-report-guide.md:78`: terminal-success stop sentence sits under the `--kind` lines and reads as covering them; prefix "Without `--kind`," — impact-user-visible → report only
- F3 `ai-report.md:16-21`: When to Use lacks the proposal / root-cause use — impact-user-visible → report only
- F4 `ai-report-reference.md:146`: Steps 3 and 4 status undefined for `--kind`; the non-visual fixed sentence can land in a brief — impact-user-visible → report only
- F5 `ai-report.md:172`: checklist says the target satisfies the shared reference, false for `--kind` — impact-negligible → report only
- F6 `ai-report.md:31`: Input sentence omits the `--kind` form and says "one UR or one REQ" — impact-negligible → report only
- F7 `architecture-report.md:135`: "`ai-report` ... presents completed work" now incomplete; rebuttal still holds — impact-negligible → report only
- F8 `ai-report-reference.md:126` (D-11): completed REQ/UR with `--kind` has no rule — impact-negligible → report only

**Nit findings:** F9 guide `:3` and REF `:3` completed-only intros — impact-negligible → report only; F10 guide `:47` implies `proposal` accepts failed/cancelled — impact-negligible → report only; F11 template restated in guide `:98-100` and checklist `:177` — impact-negligible → report only; F12 `help.md:18` breaks the command column — impact-negligible → report only; F13 REF `:142` mockup files in `generated/` vs `ai-report.md:80` "generated/ only when generated images succeed" — impact-negligible → report only
**Acceptance:** Pass — implementation and integration stages: probe re-run exit 0 at `a0ebfea9`, `git diff --check` clean, integrator gate exit 0 on record; live acceptance unassessed.
**Restatement sweep:** redefined what `ai-report` accepts (unfinished work through `--kind`). Stale: F1 (`ai-report.md:65-67`), F2 (guide `:78`), F3 (`ai-report.md:18-20`), F5 (`:172`), F6 (`:31`), F7 (`architecture-report.md:135`), F9 (guide `:3`, REF `:3`). Checked, not stale: `architecture-report.md:27`, `present-work.md:5`, `present-video.md:3`, toolbox `help.md:14`, core `help.md:37`, `:55`, `tutorial.md:245-252`, `:318`, `completed-work-presentation-reference.md:3,:5,:7`, the three prompt-injection JIT comments.
**Suggested testing:** 4 items
**Follow-ups created:** None (13 findings report only)

*Reviewed by review-work action*
