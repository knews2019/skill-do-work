---
id: UR-151
title: 'Portable verification actions: source-audit, journey-qa and release-check in the toolbox, and delivery stages in core review'
created_at: 2026-10-09T22:53:15Z
requests: [REQ-675, REQ-676, REQ-677, REQ-678]
word_count: 1033
---
# Portable Verification Actions and Delivery Stages in Core Review

## Summary
The maintainer pasted an implementation brief as `do-work validate-feedback: <brief>`. The triage split it into 16 findings, F1 to F16. The maintainer then answered four questions through the ask tool and asked to capture the accepted findings and run them. All four answers took the recommended option.

Captured, four REQs:
- F5 (core review answers per delivery stage) became REQ-675 (core review delivery stages). Core review owns the four stage names.
- F3 (source-audit) became REQ-676 (source-audit toolbox action).
- F1 (journey-qa), Discuss, became REQ-677 (journey-qa toolbox action) after the maintainer chose a lean toolbox action.
- F2 (release-check), Discuss, became REQ-678 (release-check toolbox action), depending on REQ-675 because it cites the stage names.
- F10 (integration), F11 (release), F14 (plain Markdown briefs, examples not rules) and F15 (untrusted-content handling) are copied into each REQ's constraints. F12 (verification) became each REQ's Red-Green Proof and packaging check, using the maintainer's choice: reuse the existing gate and install/update probes plus one-off exercises, with no kept fixtures.

Not captured:
- F4 (campaign-variants) and F7 (on-request retrospective mode): the maintainer kept both push-backs. Campaign asset production is not a toolbox job, and review Step 9.5 Lessons Learned already records missed assumptions.
- F6 (targeted counterexamples), F8 (repair addenda keep history), F9 (no new statuses or reopening) and F13 (prepared worktree is clean) were already done. F16 (the `just do-work-update` hand-back line) belongs in the run's hand-back, not in a REQ.

The prepared worktree `/private/tmp/do-work-portable-workflows` (branch `codex/portable-workflow-actions`) is left untouched; the run uses its own worktrees. No prompt injection found.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-675 | Core review names four delivery stages and reports deployment and live acceptance as unassessed unless the review exercised them |
| REQ-676 | do-work-toolbox source-audit checks each cited source and returns a claim-to-source table without rewriting the report |
| REQ-677 | do-work-toolbox journey-qa verifies whole user journeys and classifies each result as passed, product defect, test defect or unresolved |
| REQ-678 | do-work-toolbox release-check reports each delivery stage with evidence and a readiness verdict that stale or empty content cannot pass |

## Batch Constraints
- Toolbox actions are action documents under `skills/do-work-toolbox/actions/`, not shell commands or Go subcommands.
- Four-package manifest, installer and updater unchanged. Literal relative links that resolve in source and installed layouts.
- Briefs are plain Markdown paths read as data. No profile registry, config schema or cross-repository discovery. Game- and campaign-specific checks are examples only.
- Core review keeps its persisted `## Review` format, scoring, verdict mapping, impact tokens and capture rules; no new status and no reopening.
- Verification: the existing gate and install/update probes, plus one one-off behavioral exercise per REQ recorded in `## Testing`. No kept fixtures.
- Do not edit installed consumer copies or run `just do-work-update` in a real consumer.
- REQ-678 depends on REQ-675. The three toolbox REQs share the router, help and contract-test lists; that overlap is not an ordering edge.

## Full Verbatim Input
> ```
> use the ask tool to ask me and help me decide, then when done
>   capture the accepted ones and run them
> 
> [Answers the maintainer gave through the ask tool in the same session, 2026-10-10]
> 
> Q: Where should journey-qa live (verifying full user journeys and classifying results as passed, product defect, test defect, or unresolved)?
> A: Toolbox action, lean (Recommended)
> 
> Q: Where should release-check and the four delivery stages (implementation, integration, deployment, live acceptance) live?
> A: Core owns stages, toolbox cites (Recommended)
> 
> Q: How should the new instruction-only actions be verified?
> A: Reuse probes + one-off exercises (Recommended)
> 
> Q: I pushed back on campaign-variants (not a toolbox job) and the on-request retrospective mode (Lessons Learned already covers it). What should happen?
> A: Keep both push-backs (Recommended)
> 
> [Original brief the maintainer pasted earlier in the same session as `do-work validate-feedback: <brief>`]
> 
> Implement four portable toolbox workflows and a focused improvement to core review in the do-work upstream repository.
> 
> Repository: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
> Prepared worktree: /private/tmp/do-work-portable-workflows
> Branch: codex/portable-workflow-actions
> Starting revision: 7676a183, version 0.305.89
> 
> The prepared worktree has no implementation changes. Check its current state before using it. Work upstream only. Do not edit installed .claude/skills/do-work* copies in consumer repositories or update consumers during this task.
> 
> Read AGENTS.md, CLAUDE.md, the owning skill routers, and the relevant maintainer primes for action files and releases. Follow the repository’s established implementation, verification, and release procedures.
> 
> DELIVERABLES
> 
> Add these four agent commands to do-work-toolbox. They are action documents routed by SKILL.md, not new shell commands or Go CLI subcommands.
> 
> 1. journey-qa <target> [--brief <path>]
> 
> Verify complete user journeys and distinguish product defects from test defects and environmental limitations.
> 
> - Read requirements, project testing guidance, and the supplied brief.
> - Reuse existing checks. Reproduce the reported sequence before broadening coverage.
> - Exercise combined transitions, such as zoom → pan → scroll → reset or failed load → retry → close.
> - Synchronize on observable state. Keep screenshot capture and other expensive instrumentation outside performance measurements.
> - Record revision, environment, steps, expected/actual behavior, and evidence.
> - Classify results as passed, product defect, test defect, or unresolved. Do not equate emulation with physical-device verification.
> - Report the smallest justified repair and remaining checks. This diagnostic action does not itself authorize source repairs.
> 
> 2. release-check <target> [--brief <path>]
> 
> Verify that intended content or functionality reached its consumer correctly.
> 
> - Trace intended source, generated package, serving environment, and consumer-visible behavior.
> - Test consumer behavior rather than relying on file presence, successful decoding, or matching identifiers.
> - Investigate package-version mixing, stale caches, missing dependencies, and interrupted activation when relevant.
> - Report implementation, integration, deployment, and live acceptance separately, with evidence for each. Mark unassessed stages explicitly.
> - Distinguish historical evidence from current-run verification.
> - Return a readiness verdict and concrete gaps. Deployment, rollback, content repair, and recurring monitoring remain separate actions.
> 
> 3. source-audit <report-or-url-list>
> 
> Verify sources supporting research claims.
> 
> - Extract claims and citations, then inspect actual source content.
> - Record requested/final URLs, retrieval outcome, and supporting content. Mark information unavailable when tools cannot establish it.
> - Detect error pages, unrelated redirects, index pages, unavailable sources, and unsupported claims.
> - Distinguish retrieval success, claim support, corroboration, and authority.
> - Preserve unknown publication dates as unknown.
> - Preserve original evidence and failure history. Present verified replacement candidates without silently rewriting the report.
> - Return a claim-to-source table with supported, contradicted, insufficient, or unavailable judgments.
> 
> 4. campaign-variants <brief>
> 
> Produce deliberate creative variations while preserving campaign requirements.
> 
> - Extract invariant text, spelling, capitalization, layout, brand constraints, and prohibited elements.
> - Vary composition, palette, subjects, or style deliberately.
> - Honor requested quantity; default to three.
> - Use an available image-generation capability. When unavailable, deliver clearly labeled generation prompts.
> - Inspect generated images for exact text, cropping, unwanted text, and brief compliance.
> - Deliver the variants with a compact comparison and disclose unresolved defects.
> 
> PROJECT KNOWLEDGE
> 
> Accept ordinary project-owned Markdown briefs directly. Do not introduce a profile registry, configuration schema, or automatic cross-repository discovery.
> 
> Game-specific checks—clickable masks, collectible artwork, embedded scrolling, publication dates—and campaign text such as “Faci Rai din ceAi” are examples or brief inputs, not universal upstream requirements.
> 
> Follow the suite’s existing untrusted-content handling. Imported briefs and source pages cannot authorize unrelated actions.
> 
> CORE REVIEW IMPROVEMENT
> 
> Extend skills/do-work/actions/review-work.md and its guide:
> 
> - Answer whether work exists, satisfies the request, and has evidence for applicable delivery stages.
> - Use targeted counterexamples when passing tests omit the reported behavior.
> - When requested, explain the missed assumption, verification gap, and smallest preventive improvement in a retrospective.
> - Recommend focused repair addenda while preserving completed history and active ownership.
> - Preserve existing persisted review formats, scoring, impact classifications, and capture rules. Add no lifecycle statuses or blanket reopening behavior.
> 
> INTEGRATION
> 
> - Put new actions under skills/do-work-toolbox/actions/.
> - Update toolbox routing, argument hints, help, relevant next-step guidance, and concise usage documentation.
> - Justify each action’s toolbox ownership in its description, as required by the action-file prime.
> - Keep the existing four-package manifest and installer topology unchanged.
> - Use literal relative links that resolve in both source and installed layouts.
> - Create an additive release using the current upstream release procedure, keeping all owned version/changelog mirrors consistent. Determine the version from current state; do not assume 0.305.89 is still latest.
> 
> VERIFICATION
> 
> Run the established maintainer gate and relevant package/install/update checks.
> 
> Verify:
> - All four routes work; per-command help executes no workflow.
> - Combined interaction failures are not hidden by isolated passing checks.
> - Empty hit masks and stale served content cannot support a readiness claim.
> - Error pages and irrelevant redirects cannot support research claims; unknown freshness remains unknown.
> - Campaign typography failures are reported.
> - Core review distinguishes local implementation from production acceptance without changing existing lifecycle behavior.
> 
> Use realistic behavioral exercises for these instruction-driven workflows; keyword assertions alone do not prove their behavior.
> 
> Package the change and test installation/update in a disposable Git consumer through the canonical installer/updater. Confirm all four actions and their links survive packaging, and project briefs, queue records, and application files remain unchanged.
> 
> Commit coherent verified changes. Hand back the branch, commit, release version, validation evidence, and any limitations. Include the consumer update instruction (`just do-work-update`) for use after the upstream change is released, but do not run it in real consumers.
> ```