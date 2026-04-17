I want to improve the weekly-signal-diff skill prompt so its output more directly serves this repo's Goal section ("bring value when I can help my clients proactively and shine light on blindspots"). The prompt currently produces a rigorous structural diff centered on the operator; I want to add client-facing actionability without diluting the structural-shift philosophy.

Four concrete edits. Before touching anything, read in this order: /CLAUDE.md (operator context + Goal + Philosophy), .claude/skills/do-work/prompts/weekly-signal-diff.md (the base prompt, ~433 lines — this is where the edits land), and input-personal/weekly-signal-diff-personal.md (personal sidecar, ~117 lines).

Create a feature branch like claude/digest-actionability-<short-slug>. Commit each of the four edits as a separate commit. Push at the end.

**Edit 1 — Add "Top of mind this week" as the first subsection of the inline digest.**

In weekly-signal-diff.md Phase 7, add a new subsection BEFORE "Coverage note":

### Top of mind this week
Hard cap: 5 bullets, 150 words. Name the 3–5 things the operator should hold in working memory this week — the synthesis, not the detail. Everything else in the digest is support material for mid-week re-reading. If the week is thin, give fewer bullets rather than padding.

Add matching Rule ("Top of mind is mandatory; hard cap enforced; thin weeks produce fewer bullets, not padded ones") and a Verification checklist entry.

**Edit 2 — Move Actions from optional tail to mandatory top, split operator vs. client.**

In the same Phase 7:
- Remove the existing "### Actions (optional)" subsection from its current position at the bottom.
- Add this new subsection RIGHT AFTER "Top of mind this week" and BEFORE "Coverage note":

### Actions this week
Two mandatory groups. Be concrete; if a group is empty, say so explicitly — that is a finding, not a hole.

**For the operator** — 1–3 things the operator could act on this week, formatted as `do work capture request: <short description>` so they can capture any they want to pursue. Do not auto-capture.

**For clients** — 1–3 proactive client-outreach angles: which client archetype, what finding to raise, one-line draft of the outreach. If no shift this week has a client angle, state "No client-facing actions this week — purely structural."

Add matching Rule ("Actions section is mandatory, split operator vs. client; empty groups are stated explicitly, never omitted") and a Verification checklist entry.

**Edit 3 — Add a "For client archetypes" bullet to each headline structural shift.**

In the same Phase 7, in the "Headline structural shifts" template, add a new bullet AFTER "Why it matters to this user":

- **For client archetypes** — optional per-shift. If this shift is useful to a specific client type the operator serves, name the archetype and a one-line outreach angle. If nothing client-facing, write "No direct client angle." Never collapse this into the "why it matters to this user" paragraph — keep them visually distinct.

Add a new Common Rationalizations row: "'The client angle is obvious from context' → write it out anyway → obvious to you ≠ obvious at a glance mid-week."

**Edit 4 — Ritualize calibration.**

In input-personal/weekly-signal-diff-personal.md, find lane 42 (Calibration — weekly-signal-diff forecast review and accuracy ledger). In its "Why this lane matters" cell, replace the promote-trigger sentence ("Promote quarterly or when a predicted shift lands") with:

"Promote every weekly run once 4+ prior runs exist. Each week begins with a short forecast scoreboard — prior weeks' 'Watch next' items scored happened / partial / not yet / stale. Accumulates into a calibration ledger over time."

**Anti-drift guardrails — do NOT do these:**

- Do not rewrite the structural-questions filter in Phase 5. The constraint-movement gate is intentional and load-bearing.
- Do not collapse "why it matters in general" and "why it matters to this user" and "for client archetypes" into a single framing. The repo's operator context requires them visually separated.
- Do not add a cap on the headline shift count; thin weeks are by design.
- Do not auto-run the pipeline to test. Validate by re-reading the prompt file end-to-end and checking internal consistency.
- Do not touch the pipeline prompt at .claude/skills/do-work/prompts/weekly-signal-diff-pipeline.md — it's orchestration plumbing, not content.

**Verification before pushing:**

- Re-read Phase 7 top-to-bottom. Section order should be: Top of mind → Actions this week → Coverage note → Headline structural shifts → Per-lane scan notes → Triggered-scan sweep → What didn't change → What changed from last week → Likely big news horizon → Watch next.
- Check the Verification checklist section at the bottom of weekly-signal-diff.md — new sections (Top of mind, Actions) need checklist entries.
- Verify the lane 42 row in the sidecar still parses as a valid markdown table (pipe count consistent with other rows).

Report: branch name, one-line summary per commit, file paths touched.

