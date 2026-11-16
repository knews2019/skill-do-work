---
id: REQ-739
title: 'Second-opinion header keeps a side margin on a phone'
status: claimed
created_at: 2026-11-16T20:15:49Z
user_request: UR-181
domain: frontend
prime_files: ["_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-738]
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION"]
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-16T20:26:12Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
claimed_at: 2026-11-16T20:24:22Z
dispatch_at: 2026-11-16T20:30:25Z
builder_handback_at: 2026-11-16T20:39:28Z
---
# Second-Opinion Header Keeps a Side Margin on a Phone
## What
The short header that `ai-report second-opinion` writes above a directed page has no side margin. On a 390 px phone, "Iteration v1cNN" and the fold line "▸ Preset: …" start at x = 0 and touch the screen edge. Give the header a side margin on pages that do not style it themselves, without overriding pages that do.
## Context
- The header was shipped in 0.324.1 by REQ-738 (second-opinion direction block keeps the page content on the first phone screen). The 0.324.0 header had the same gap.
- The block is written by `secondOpinionDecorate` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (the `block :=` line builds `<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">` between the `<!-- second-opinion-header -->` markers). Its header test is `TestSecondOpinionHeaderFoldsTheFullDirection` in `report_second_opinion_test.go`.
- Evidence:
  - `do-work/runs/work-2026-11-16-171725/REQ-737-captures/final-390-load.png`: the final REQ-737 proof page at 390 px; the header text touches the left edge.
  - REQ-738's review finding F3 (`do-work/archive/UR-180/REQ-738-second-opinion-direction-block-first-screen.md`, `## Review`): "no side gutter for the header on kit and plain pages (`asideLeft` 0 px on the kit page at 390)". The same file's `## Discovered Tasks` and `do-work/runs/work-2026-11-16-171725/REQ-738-handback.md` § Discovered Tasks say an inline padding would override the tstack page's own rule, so the fix needs its own decision.
  - `do-work/runs/work-2026-11-16-171725/REQ-737-handback.md` § Discovered Tasks (pass 3): "Iteration v1c04" and the summary triangle start at x = 0 at 390 px. The page may not style the block (REQ-737 D-07), so this is a launcher change.
- A page that styles the header itself: `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html` has `body>aside[aria-label="Second-opinion direction"]{…margin:12px max(16px,…);padding:8px 14px;…}` (REQ-727 D-13). An inline `style` padding or margin on the `aside` would win over that rule and change the page's look.
## Detailed Requirements
- On a page with no rule of its own for the header (a plain page, a report-kit page), the header's text starts about 16 px from the left screen edge at 390 px, and the right side keeps a matching margin.
- On a page that styles `body>aside[aria-label="Second-opinion direction"]` itself, the page's own margin and padding still apply unchanged.
- Decorating stays idempotent (decorating twice gives the same bytes), and a second decorate still strips the whole old block, including a 0.324.0 or 0.324.1 block.
- The page never scrolls sideways at 390 px, closed or open, and `ai-report-judge` still passes.
## Constraints
- Only the launcher's header block changes. Head metas, markers, the fold, the presets and the iteration list stay as they are.
- Do not rewrite published pages under `ai-reports/`. Measure on scratch copies.
- Release rules from `_dev/primes/prime-releases.md` apply: a CHANGELOG entry and a VERSION bump, written by the integrator.
## Dependencies
None. REQ-738 is archived; this REQ changes the block it shipped.
## Red-Green Proof
**RED prompt/case:** A Go test (beside `TestSecondOpinionHeaderFoldsTheFullDirection`) decorates a page and checks that the block carries the side-margin rule in the agreed low-priority form, and carries no inline `padding` or `margin` on the `aside` itself.
**Why RED now:** `secondOpinionDecorate` writes the `aside` with only `overflow-wrap:anywhere`, so on a page without its own rule the header text starts at x = 0.
**GREEN when:** (1) The Go test passes. (2) A scratch page decorated by the new code with no header style of its own, opened at 390 x 844 px, has the header text's left edge at about 16 px (> 0), measured by the `aside`'s first text box. (3) A scratch copy of the tstack v1c03 page decorated by the new code shows the same computed margin (12 px top, at least 16 px sides) and padding (8 px 14 px) on the `aside` as before. (4) No sideways scroll at 390 px, closed and open, and `ai-report-judge` passes.
**Validation:** Inferred during capture. The 390 px measurement and "the v1c03 page keeps its own styling" are the coordinator's stated proof; the 16 px figure is approximate.
## Builder Guidance
- Decided at capture (unattended, no question asked): put the margin in a rule the page's own rule always beats, for example a `<style>` element inside the header markers with `:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}` (zero specificity, so any page rule for the same element wins). Keeping it inside the markers keeps the strip and idempotence working. The builder may choose another form if it meets GREEN (2) and (3), and records the choice in Decisions.
- Check the report kit page too: the kit must not restyle the new rule, and the fold line must still line up with the iteration line.
- Run the whole `toolboxcommands` package tests, since the block is on every directed page.
- Show a 390 px capture of a plain page's header and of the v1c03 header to the maintainer before the REQ is called done.
## Full Context
See `do-work/user-requests/UR-181/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (bare; the row is `slugged: partial`, so a targeted form is not legal) — cost 25,074 tokens against the 2,000-token budget — matches "a session that is its own second-opinion backend" and "writer-anchored-marker-position": the change edits what the launcher writes between the header markers.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Integrate + capture fix (Recommended): … Capture one new small request for a side margin on the launcher's header, queued for a later run. … Risk: until that request runs, the header touches the screen edge on a phone.*

## Triage

**Route: B** - Medium

**Reasoning:** The file and the function are named, but the fix must not override a page's own rule for the header, and it must survive the report kit's CSS, so the builder needs exploration of the kit, the tstack page rule and the judge before choosing the rule's form, plus a browser proof on scratch copies.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Where the change goes**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:1603`, function `secondOpinionDecorate`. Line 1643 builds the block: `headerStart + metas + <aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere"> + paragraphs + </aside> + headerEnd`. Lines 1597-1602 are the function comment; it must name the new margin rule. Line 1636 (the `details` fold), the metas (1625-1631), the strip calls (1607-1618) and the insert position (1644-1647) stay as they are.
- Callers: line 841 (strip only, before the backlink check) and line 864 (the real decorate at publish). Neither changes.
- `secondOpinionStripBlocks` (1565-1589) removes everything between the two marker comments on the raw-text-masked surface. `secondOpinionRawHTML` (1473) masks `<style>` bodies, so a `<style>` element placed between the markers is removed by the same strip, and its text can never forge a marker. A 0.324.0 or 0.324.1 block is stripped with no extra work.

**Tests that pin the header today**
- `report_second_opinion_test.go:1008` `TestSecondOpinionHeaderFoldsTheFullDirection` (REQ-738). Its `header` helper (1022-1036) treats everything outside the `details` body as visible; a CSS rule with no direction text keeps checks 1-6 true. Check 7 (line 1074) asserts the exact `aside` open tag with `style="overflow-wrap:anywhere"`, so the D-01 form keeps that tag byte-identical. Check 5 (1057) holds a literal 0.324.0 block.
- `report_second_opinion_test.go:943` `TestSecondOpinionDirectedRunNamesAndMarksPage`: markers, metas and free text (969); idempotence, metas before the real head end, and a plain decorate strips back to the input bytes (990-1001).
- `report_second_opinion_tree_test.go:43` `TestSecondOpinionTreeTargetsNamesAndParentHeader`: the parent link and the id and base metas (79); a plain run has no marks (84).
- No other Go code, script or shipped doc reads the block (grep for the label and the marker over `skills`, `_dev`, `suite`).

**Pages and kit the rule meets**
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html` styles `body>aside[aria-label="Second-opinion direction"]{max-width:1368px;margin:12px max(16px,calc((100% - 1368px)/2));padding:8px 14px;...}` and `... p{margin:.25em 0}` (REQ-727 D-13). Its selector has specificity (0,1,2). A rule wrapped in `:where()` has specificity 0, so the page's padding wins whatever the order; an inline `style` padding would beat it.
- The same bundle's v1c01 and v1c04 pages set `body{margin:0}` and have no rule for the header, so their header text starts at x = 0 today. v1c04 is REQ-737's final page.
- Report kit: `skills/do-work-toolbox/docs/ai-report-kit/presentation.css:62` sets `* { box-sizing: border-box; }` and lines 66-68 set `body { margin: 0; ... }`; no kit rule names `aside`. The kit `p { margin: 0 0 1em; }` (line 79) and the `details`/`summary` rules (196-199) are already cancelled or accepted by REQ-738. `presentation.js:44-50` ("Open all details") toggles only `details` inside its own section; the header is outside every section.
- `skeleton.html:8` asks page CSS to sit in one `<style>` in the head; a launcher `<style>` inside the header block is not page CSS and does not restyle the kit, but it is a `<style>` element in the body. Browsers apply it; the judge (`report_judge.go`) does not inspect `<style>` placement.
- No published page sets a Content-Security-Policy meta (the grep hits are prose), so an inline `<style>` element is allowed wherever the existing inline `style` attribute is.

**Measured baseline (from REQ-738's integration)**: header 87.4 px and map frame top y = 830 on the v1c03 page at 390 x 844; `asideLeft` 0 px on the kit page at 390 (REQ-738 review F3). The v1c03 page keeps its own rule, so its numbers must not move.

**Lessons read**
- `_dev/primes/lessons-releases.md` (required): no archive move and no manifest edit here; the release is the integrator's.
- Dropped for budget, read by fixed-string grep: `[family: writer-anchored-marker-position]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:105-107` (no whitespace outside the markers; strip on the masked surface; decorating twice gives the same bytes) and `[family: session-is-the-backend]` lines 17-20 (measure a header with the page's own CSS at the real width; look at a phone capture, not only the numbers).
- `do-work/archive/UR-180/REQ-738-second-opinion-direction-block-first-screen.md` Decisions, Review (F1-F4) and Lessons Learned; `do-work/runs/work-2026-11-16-171725/REQ-738-integration-report.md`.

*Generated by the pre-dispatch agent (no Explore subagent; the guide forbids subagents)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — the header block string in secondOpinionDecorate and its function comment
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — the new RED test for the side-margin rule (D-04)
- `CHANGELOG.md` (modify, integrator only) — release entry
- `skills/do-work/CHANGELOG.md` (modify, integrator only) — byte-identical mirror
- `VERSION` (modify, integrator only) — version bump (D-06)

**Files I will NOT touch:** every other function of the launcher file (strip, metas, fold, presets, iteration list); the tree test file (its asserts hold); the two contract docs (D-05); the report kit; the judge; every published file under the reports folder (proofs use scratch copies); anything under the queue folder.

**Acceptance criteria (restated from REQ):**
- [ ] A Go test, RED first, decorates a page and shows the block carries the side-margin rule in a low-priority form and no inline padding or margin on the `aside`.
- [ ] On a page with no rule of its own (plain page, report-kit page), the header text starts about 16 px from the left edge at 390 px, with a matching right margin.
- [ ] On a page that styles `body>aside[aria-label="Second-opinion direction"]` itself (a scratch copy of the tstack v1c03 page), the computed margin (12 px top, at least 16 px sides) and padding (8 px 14 px) of the `aside` are unchanged.
- [ ] Decorating twice gives the same bytes, and a second decorate strips a whole 0.324.0 or 0.324.1 block.
- [ ] Head metas, markers, the fold, the presets and the iteration list are unchanged.
- [ ] No sideways scroll at 390 px with the fold closed or open, and the render judge passes.
- [ ] The maintainer is shown the 390 px captures before the REQ is called done (D-07).

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-16-202422/REQ-739-preflight-probe.sh) and the repository gate green at 49e54593 (156 s, detached checkout, stage reuse off, load 4.6)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-16-202422)*

## Decisions

(pre-dispatch)

- **D-01** (DECIDE & STATE): Recommended form, which the builder may change if GREEN (2) or (3) fails (it then records its own decision and names the forms it rejected). Inside the markers, right before the `aside`, the block gets one `<style>` element: `:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}`. `:where()` has zero specificity, so any page rule for the same element (the tstack v1c03 `padding:8px 14px`) wins whatever its order. The `aside` open tag stays byte-identical (`style="overflow-wrap:anywhere"` only), so REQ-738's check 7 and every page selector still match. Rejected: an inline padding or margin on the `aside` (it beats the page's own rule, REQ-738 Discovered Tasks); a `<style>` in the head beside the metas (a second thing to strip with its own matcher, and a page with no head gets the block anyway); a class with a rule (a new name for the same effect); padding on the `p` and `summary` elements inline (beats the page's `p{margin:.25em 0}` neighbour rules and still needs the fold to line up).
- **D-02** (DECIDE & STATE): The 16 px is padding inside the `aside`, so the text's left edge is the body margin plus 16 px: 16 px on the report kit and on pages with `body{margin:0}` (v1c01, v1c04), 24 px on a page that keeps the browser's default 8 px body margin. The REQ's 16 px is approximate, so both count as met; the hand-back reports the exact numbers. The right side gets the same padding.
- **D-03** (DECIDE & STATE): The rule applies at every width, with no media query. A header touching the screen edge is wrong on a desktop too, and one rule is simpler than two.
- **D-04** (DECIDE & STATE): The RED test is one new function named exactly `TestSecondOpinionHeaderKeepsAPhoneSideMargin` in `report_second_opinion_test.go`, beside `TestSecondOpinionHeaderFoldsTheFullDirection`. It calls `secondOpinionDecorate` directly (milliseconds). It pins: the block between the markers carries the side-margin rule in a zero-specificity form; the `aside` carries no inline `padding` or `margin`; decorating twice gives the same bytes; a page holding a literal 0.324.1 block (copied from today's line 1636 and 1643 shape) decorates to the same bytes as the clean page. The GREEN probe names this test.
- **D-05** (DECIDE & STATE): No contract text changes. `skills/do-work-toolbox/actions/ai-report.md:89` describes what the header shows (the preset at the top, the full direction folded below it), which stays true; spacing is not part of that contract.
- **D-06** (DECIDE & STATE): The integrator picks the version bump. Suggested: a patch (0.324.2 to 0.324.3), because an existing header changes its look and no option is added. The builder only proposes the changelog entry.
- **D-07** (ESCALATE): The REQ asks to show a 390 px capture of a plain page's header and of the v1c03 header to the maintainer before the REQ is called done. The run is unattended, so the builder saves captures under `do-work/runs/work-2026-11-16-202422/REQ-739-captures/` and the coordinator shows them before integration starts. Value: the maintainer sees the margin on real pages before it ships to every directed page. Risk: if the maintainer wants another width or form, the change is one block string and one test, fully reversible in a follow-up; pages published in between keep the header they got until a new pass rebuilds them.

<!-- D-XX counter: last used D-07. Next decision: D-08. -->
