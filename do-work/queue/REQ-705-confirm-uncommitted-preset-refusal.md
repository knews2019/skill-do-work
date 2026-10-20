---
id: REQ-705
title: 'Confirm: a second opinion refuses to run while a project preset file is uncommitted'
status: pending-answers
created_at: 2026-10-20T20:47:28Z
user_request: UR-160
addendum_to: REQ-703
builder_decided: true
---

# Confirm: A Second Opinion Refuses To Run While A Project Preset File Is Uncommitted

## What
During REQ-703 (second opinions accept a named preset and a free-text editor direction), the plan kept one existing safety check unchanged: before and after each second-opinion run, the launcher refuses to continue when Git shows any change outside the report bundle. A project preset file, `ai-reports/second-opinion-presets/<name>.md`, sits outside the bundle. So the first run after you write or edit a project preset is refused until you commit that file. This was escalated because it trades convenience against the one check that catches a backend writing files outside the bundle, and only you can judge which matters more.

## What the Builder Chose
Keep the check whole. A run with an uncommitted project preset fails with `outside-bundle changes: ai-reports/second-opinion-presets/<name>.md`, and `skills/do-work-toolbox/docs/ai-report-guide.md` tells you to commit a new preset before using it. Shipped presets (`bold`, `visual`) are not affected, because they live in the installed toolbox, not in your project.

## What Would Change
If you prefer the exemption: the outside-bundle check in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` would ignore the one preset file the run actually read, and only when its content is unchanged between the check before the run and the check after it. One test would pin that a backend editing the preset is still refused.

## Open Questions
- [ ] Should a second opinion run while the project preset file it uses has uncommitted changes?
  Recommended: No, keep refusing until the preset is committed (already implemented).
  Value: the check that catches any backend write outside the report bundle stays whole, with no special case.
  Risk: the first run after you write or edit a preset is refused, and you must commit the file and run again; reversible later by adding the exemption.
  Also: Yes, exempt the preset file the run read, as long as its bytes did not change during the run.
