---
id: UR-170
title: 'Reproduce instruction retention across compaction'
created_at: 2026-11-10T23:19:33Z
requests: [REQ-718]
word_count: 3
---
## Summary
Capture the sole accepted finding, F9 (compaction reproduction), from the preceding validate-feedback report. F1 and F3 were Already done; F2 and F4 were Push back; F5–F8 were Discuss and are not promoted by this invocation.

## Referenced Accepted Finding

> **F9 — “Investigating the break points … is absolutely worth it” → Accept**
> 
> The useful next step is a small behavioral experiment: load a real work phase and its crew rules, compact mid-phase, then inspect subsequent file reads and compliance with a harmless rule. Record the harness version and configuration.
> 
> **Remedy:** Run a bounded compaction reproduction before proposing new hooks or instruction rules.
> **Surface-cost:** N/A for the experiment; it adds no permanent runtime layer.

## Provenance and Evidence
- Finding: F9, the only Accept verdict in the preceding validate-feedback report.
- Original source: pasted text, Matt Maher, “AI Agent Secrets Nobody Tells You”, dated November 10, 2026; transcript near 19:24–19:32.
- Verbatim source claim: “Investigating the break points that aren't written down anywhere is absolutely worth it.”
- Original severity: not supplied.
- Evidence: `skills/do-work/actions/work.md:21–25` requires procedure loads once per context, introduced by commit `22f82fed`; `skills/do-work/tools/do-work-cli/internal/hookcommands/session_start.go:40–87` performs housekeeping but does not restore active instruction files. Static inspection did not establish whether required reads recur after compaction.
- Supporting documentation: https://code.claude.com/docs/en/context-window documents bounded restoration of file and skill contents; https://learn.chatgpt.com/docs/build-skills documents progressive disclosure, not the transcript's absolute Codex restoration claim.
- Surface-cost: N/A for the bounded experiment, which adds no permanent runtime layer. New restoration hooks or enforcement need separate evidence and authorization.

## Full Verbatim Input
> ```
> capture-request all accepted
> ```
