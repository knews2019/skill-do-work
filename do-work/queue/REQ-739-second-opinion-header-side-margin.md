---
id: REQ-739
title: 'Second-opinion header keeps a side margin on a phone'
status: pending
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
