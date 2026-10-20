---
id: REQ-705
title: 'Confirm: a second opinion refuses to run while a project preset file is uncommitted'
status: claimed
created_at: 2026-10-20T20:47:28Z
user_request: UR-160
addendum_to: REQ-703
builder_decided: true
status_changed_at: 2026-10-20T21:19:07Z
claimed_at: 2026-10-20T21:19:46Z
---

# Confirm: A Second Opinion Refuses To Run While A Project Preset File Is Uncommitted

## What
During REQ-703 (second opinions accept a named preset and a free-text editor direction), the plan kept one existing safety check unchanged: before and after each second-opinion run, the launcher refuses to continue when Git shows any change outside the report bundle. A project preset file, `ai-reports/second-opinion-presets/<name>.md`, sits outside the bundle. So the first run after you write or edit a project preset is refused until you commit that file. This was escalated because it trades convenience against the one check that catches a backend writing files outside the bundle, and only you can judge which matters more.

## What the Builder Chose
Keep the check whole. A run with an uncommitted project preset fails with `outside-bundle changes: ai-reports/second-opinion-presets/<name>.md`, and `skills/do-work-toolbox/docs/ai-report-guide.md` tells you to commit a new preset before using it. Shipped presets (`bold`, `visual`) are not affected, because they live in the installed toolbox, not in your project.

## What Would Change
If you prefer the exemption: the outside-bundle check in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` would ignore the one preset file the run actually read, and only when its content is unchanged between the check before the run and the check after it. One test would pin that a backend editing the preset is still refused.

## Open Questions
- [x] Should a second opinion run while the project preset file it uses has uncommitted changes? → Neither: presets live only inside the report bundle (bundle-local presets folder) plus the shipped ones; drop the shared ai-reports/second-opinion-presets/ folder
  Recommended: No, keep refusing until the preset is committed (already implemented).
  Value: the check that catches any backend write outside the report bundle stays whole, with no special case.
  Risk: the first run after you write or edit a preset is refused, and you must commit the file and run again; reversible later by adding the exemption.
  Also: Yes, exempt the preset file the run read, as long as its bytes did not change during the run.


## Answer Notes

- 2026-10-20 - [ ] Should a second opinion run while the project preset file it uses has uncommitted changes?: Neither: presets live only inside the report bundle (bundle-local presets folder) plus the shipped ones; drop the shared ai-reports/second-opinion-presets/ folder
> ```
> User, asked whether a run should refuse while a project preset in ai-reports/second-opinion-presets/ is uncommitted: "why would you put the second opinion presets outside the report folder?"
> User, asked where their own presets should live (options: shared folder + commit first / also inside the bundle / only inside the bundle): "Only inside the bundle".
> Reasoning: a preset kept inside the report bundle is never outside-bundle dirt, so the question disappears instead of needing a carve-out. Out of scope now: the shared project folder ai-reports/second-opinion-presets/ (drop its lookup and the catalog skip that only existed for it); shipped presets beside the prompt stay.
> ```
