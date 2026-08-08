---
id: REQ-145
title: "Remove the Stateful Pipeline"
status: completed
claimed_at: 2026-08-08T17:03:56Z
completed_at: 2026-08-08T17:48:23Z
route: C
kb_status: pending
kb_entry:
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: []
write_set:
  - .gitignore
  - README.md
  - _dev/tests/staged-skills-contract.sh
  - _dev/tests/install-suite-behavior.sh
  - _dev/tests/contract-regressions.sh
  - tools/install-do-work-suite.sh
  - skills/do-work/tools/install-do-work-suite.sh
  - skills/do-work/SKILL.md
  - skills/do-work/actions/help.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/abandon.md
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/kb-lessons-handoff.md
  - skills/do-work/actions/pipeline.md
  - skills/do-work/actions/pipeline-reference.md
  - skills/do-work/next-steps.md
  - skills/do-work/hooks/hooks.json
  - skills/do-work/hooks/session-start.sh
  - skills/do-work/hooks/pipeline-guard.sh
  - skills/do-work/docs/standing-preferences.md
  - skills/do-work/crew-members/anti-slop.md
  - skills/do-work/crew-members/approach-directives.md
  - skills/do-work/crew-members/background-agents.md
  - skills/do-work-knowledge/actions/interview.md
  - skills/do-work-knowledge/hooks/memory-stop-capture.sh
  - skills/do-work-knowledge/crew-members/anti-slop.md
  - skills/do-work-knowledge/crew-members/background-agents.md
  - skills/do-work-toolbox/actions/ai-report.md
  - skills/do-work-toolbox/actions/note.md
  - skills/do-work-toolbox/actions/present-work.md
  - skills/do-work-toolbox/actions/tutorial.md
  - skills/do-work-toolbox/docs/ai-report-guide.md
  - skills/do-work-toolbox/actions/slop-check.md
  - skills/do-work-toolbox/actions/stray-check.md
  - skills/do-work-toolbox/docs/stray-check-guide.md
  - skills/do-work-toolbox/crew-members/anti-slop.md
  - skills/do-work-toolbox/crew-members/approach-directives.md
  - skills/do-work-toolbox/crew-members/background-agents.md
  - decisions/records/adr-005-pipeline-is-stateful-and-resumable.md
  - decisions/records/adr-001-modular-action-prompts-and-companion-references.md
  - decisions/records/adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles.md
  - decisions/records/adr-007-close-the-pipeline-with-present-and-a-technical-debrief.md
  - decisions/records/adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats.md
  - decisions/records/adr-019-four-skill-suite-contract.md
  - decisions/topics/_index_workflow-orchestration.md
  - decisions/topics/_index_pipeline-deliverables.md
  - decisions/topics/_index_skill-architecture.md
  - decisions/_master_index.md
  - decisions/log.md
tdd: true
suggested_spec: refactor
depends_on: [REQ-144]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-146]
batch: do-work-four-skill-suite
---

# Remove the Stateful Pipeline

## What
Remove the separate resumable pipeline state machine after modular cutover and replace it with a copyable full-cycle prompt.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Ratchet the three declared shell contract suites first and capture assertion failures against the live pipeline runtime; then delete the three runtime files, remove only pipeline-specific routes/state readers/docs, reconcile both byte-identical installers so they preserve custom Stop hooks while pruning the retired guard, add the approved UR-031 prompt byte-for-byte, update living decisions, and run the focused plus distribution regression checks.
- [x] **[APPLY]:** Retirement contracts were written and observed RED before runtime edits; the three stateful runtime files, routes, state readers, guard hook, and stale live guidance are removed; README/help carry the exact approved prompt; both installer paths preserve custom Stop hooks while pruning only the retired guard; ADR-005 through ADR-008 are superseded. Scope stayed within the declared list plus the D-01 expansion recorded below.
- [x] **[UNIFY]:** `git diff --stat` reports 50 project files (419 insertions, 865 deletions). `git diff --check`, Bash syntax, warning-threshold ShellCheck, installer byte identity, the retired-surface sweep, protected-file hashes, and added-debug-marker audit all passed. Per-file review:
  - Contract tests — `_dev/tests/staged-skills-contract.sh`, `_dev/tests/install-suite-behavior.sh`, `_dev/tests/contract-regressions.sh`: retirement absence, exact prompt, jq/Python mixed-entry custom Stop preservation, and RED→GREEN outputs checked.
  - Installer/runtime — `tools/install-do-work-suite.sh`, `skills/do-work/tools/install-do-work-suite.sh`, `.gitignore`, `skills/do-work/hooks/hooks.json`, `skills/do-work/hooks/session-start.sh`: byte identity, only-retired-guard pruning, JSON/shell validity, and no state reader checked; `skills/do-work/actions/pipeline.md`, `skills/do-work/actions/pipeline-reference.md`, and `skills/do-work/hooks/pipeline-guard.sh` verified deleted.
  - Core routing/guidance — `README.md`, `skills/do-work/SKILL.md`, `skills/do-work/actions/help.md`, `skills/do-work/actions/work.md`, `skills/do-work/actions/review-work.md`, `skills/do-work/actions/abandon.md`, `skills/do-work/actions/capture.md`, `skills/do-work/actions/kb-lessons-handoff.md`, `skills/do-work/next-steps.md`, `skills/do-work/docs/standing-preferences.md`: no alias/state lifecycle, exact prompt, built-in tests/review, and surviving work lifecycle checked.
  - Core crew copies — `skills/do-work/crew-members/anti-slop.md`, `skills/do-work/crew-members/approach-directives.md`, `skills/do-work/crew-members/background-agents.md`: removed retired caller/state claims while preserving generic workflow guidance.
  - Knowledge surfaces — `skills/do-work-knowledge/actions/interview.md`, `skills/do-work-knowledge/hooks/memory-stop-capture.sh`, `skills/do-work-knowledge/crew-members/anti-slop.md`, `skills/do-work-knowledge/crew-members/background-agents.md`: stale state/guard/caller references removed; Stop capture still non-blocking.
  - Toolbox surfaces — `skills/do-work-toolbox/actions/ai-report.md`, `skills/do-work-toolbox/actions/note.md`, `skills/do-work-toolbox/actions/present-work.md`, `skills/do-work-toolbox/actions/tutorial.md`, `skills/do-work-toolbox/actions/slop-check.md`, `skills/do-work-toolbox/actions/stray-check.md`, `skills/do-work-toolbox/docs/ai-report-guide.md`, `skills/do-work-toolbox/docs/stray-check-guide.md`, `skills/do-work-toolbox/crew-members/anti-slop.md`, `skills/do-work-toolbox/crew-members/approach-directives.md`, `skills/do-work-toolbox/crew-members/background-agents.md`: retired completion/state guidance removed and full-cycle teaching points at capture → verify → run → present.
  - Decisions — `decisions/records/adr-001-modular-action-prompts-and-companion-references.md`, `decisions/records/adr-005-pipeline-is-stateful-and-resumable.md`, `decisions/records/adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles.md`, `decisions/records/adr-007-close-the-pipeline-with-present-and-a-technical-debrief.md`, `decisions/records/adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats.md`, `decisions/records/adr-019-four-skill-suite-contract.md`, `decisions/topics/_index_workflow-orchestration.md`, `decisions/topics/_index_pipeline-deliverables.md`, `decisions/topics/_index_skill-architecture.md`, `decisions/_master_index.md`, `decisions/log.md`: supersession/history consistency checked; ADR-019's pre-existing attestation hunk preserved separately from the new retirement hunk.

## Why
Pipeline duplicates composition and persistent state while the existing `do-work run` orchestrator already owns testing and review.

## Detailed Requirements
- Delete pipeline action/reference, routing, `pipeline.json` lifecycle, pipeline guard hook, session-start reporting, tests, and stale docs.
- Do not weaken `do-work run`; retain triage, implementation, testing, review, lessons, archival, and commit behavior.
- Add the approved copyable prompt to core help and documentation.
- The prompt must capture the request, record its UR, verify it, run its REQs with built-in tests/review, invoke `do-work-toolbox present-work`, and stop/report on failure.
- Explain that testing is already inside `do-work run`; do not create another stateful testing stage.
- Do not retain a pipeline alias or recreate pipeline state under another name.

## Constraints
- This happens only after REQ-144 activates the suite.
- Preserve the exact approved full-cycle prompt in the UR context.

## Dependencies
Requires REQ-144.

## Builder Guidance
Certainty level: Firm. This is a deletion/narrowing pass; remove pipeline-specific machinery rather than adding another orchestration layer.

## Red-Green Proof
**RED prompt/case:** Inspect core routing/help and the runtime tree for pipeline routes, actions, hooks, or state-file handling.
**Why RED now:** Stateful pipeline machinery is present and advertised as a core workflow.
**GREEN when:** No pipeline runtime surface remains, the full-cycle prompt is copyable and accurate, and `do-work run` still passes its full orchestrator tests.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the exact replacement prompt.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

---

## Triage

**Route: C** - Complex

**Reasoning:** Removing the stateful pipeline changes routing, actions, hooks, tests, and public documentation while preserving the full `do-work run` lifecycle and introducing an exact replacement prompt.

**Planning:** Required

## Plan

1. **Turn the retirement contract RED.** Update the staged-skill, installer-behavior, and core contract regression suites so they require the pipeline action/reference/guard and public aliases to be absent, require custom Stop hooks to survive legacy-guard removal, and require the approved replacement prompt byte-for-byte in core help and the README.
2. **Delete the stateful runtime.** Remove `skills/do-work/actions/pipeline.md`, `skills/do-work/actions/pipeline-reference.md`, and `skills/do-work/hooks/pipeline-guard.sh`; remove the `do-work/pipeline.json` ignore entry.
3. **Remove live routes and state readers.** Delete the `pipeline`/`full` routes and help exceptions, the hooks Stop event, session-start `pipeline.json` reporting, pipeline-specific next-step/help/work/abandon/review wording, and standing-preference claims. Keep `do-work run` behavior unchanged.
4. **Publish the exact successor.** Add the UR-031 full-cycle prompt verbatim to core help and the root README, explaining that `do-work run` already owns testing and review and that no stateful testing stage remains.
5. **Reconcile installed hooks safely.** Update both installer copies so jq and Python composition delete only the retired pipeline-guard Stop entry, preserve custom Stop hooks, prune an empty Stop list when appropriate, and validate absence rather than presence.
6. **Update shipped sibling docs and historical decisions.** Remove live references to the deleted runtime from toolbox/knowledge docs and duplicated crew files. Mark ADR-005 through ADR-008 superseded and update workflow indexes/log/ADR-019 without rewriting changelogs, audits, KB evidence, or archived REQs.
7. **Verify and release.** Run focused RED/GREEN tests, suite/update/hash guards, shell syntax/lint, board Go tests, a shipped-surface negative sweep, and diff review. The owner then performs the shared version/changelog commit.

**Planned project files:**
- `.gitignore`
- `_dev/tests/staged-skills-contract.sh`
- `_dev/tests/install-suite-behavior.sh`
- `_dev/tests/contract-regressions.sh`
- `skills/do-work/actions/pipeline.md` (delete)
- `skills/do-work/actions/pipeline-reference.md` (delete)
- `skills/do-work/hooks/pipeline-guard.sh` (delete)
- `skills/do-work/SKILL.md`
- `skills/do-work/actions/help.md`
- `README.md`
- `skills/do-work/next-steps.md`
- `skills/do-work/hooks/hooks.json`
- `skills/do-work/hooks/session-start.sh`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/abandon.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/docs/standing-preferences.md`
- `tools/install-do-work-suite.sh`
- `skills/do-work/tools/install-do-work-suite.sh`
- shipped toolbox/knowledge references and duplicated crew files identified during exploration
- ADR-005 through ADR-008, ADR-019, workflow/pipeline topic indexes, master index, and decision log
- owner release files: `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`

**Requirement coverage:** runtime/state/hook/routing deletion is covered by tasks 1–3; exact capture→verify→run→present prompt and built-in tests/review language by task 4; no alias/recreated state/testing stage by the negative ratchets; stale documentation and safe installed-hook migration by tasks 5–6.

**Plan validation:** All Detailed Requirements map to planned tasks, and every planned task traces to the retirement or its verification. ⚠ The plan has 7 tasks — quality degrades past 3; the breadth is inherent to deleting a cross-cutting public workflow, so implementation must stay within the enumerated retirement surfaces and avoid unrelated cleanup.

*Generated by Plan agent*

## Exploration

The live retirement surface spans three deleted core files, core routing/help/hooks, the two byte-identical suite installers, three shell contract suites, sibling documentation, duplicated crew guidance, and living decision indexes. The root and staged-core installer copies must remain byte-identical; removing `Stop` from `hooks.json` alone is insufficient because existing clients may retain a legacy guard entry. Both jq and forced-Python composition paths must remove only the retired guard while preserving custom Stop hooks.

Historical changelogs, audits, archived REQs, KB raw evidence, and migration progress remain immutable. Generic uses of “pipeline” for `do-work run`, CI, data flows, or orchestration are not the retired state machine and must not be over-deleted. The approved prompt must be copied byte-for-byte from UR-031. `README.md`'s bootstrap block must remain identical to installer output. The existing unrelated ADR-019 hunk is preserved separately; only the new retirement hunk may enter this REQ's commit.

Relevant verification follows the repository's shell fixture convention: staged-skill, installer-behavior, and contract-regression suites for RED/GREEN; suite-manifest, update-script, and commit-hash guard suites for distribution regressions; syntax/ShellCheck, installer byte identity, board Go tests, and a narrow negative sweep for the removed runtime tokens.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `.gitignore`
- `README.md`
- `_dev/tests/staged-skills-contract.sh`
- `_dev/tests/install-suite-behavior.sh`
- `_dev/tests/contract-regressions.sh`
- `tools/install-do-work-suite.sh`
- `skills/do-work/tools/install-do-work-suite.sh`
- `skills/do-work/SKILL.md`
- `skills/do-work/actions/help.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/abandon.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/actions/kb-lessons-handoff.md`
- `skills/do-work/actions/pipeline.md` (delete)
- `skills/do-work/actions/pipeline-reference.md` (delete)
- `skills/do-work/next-steps.md`
- `skills/do-work/hooks/hooks.json`
- `skills/do-work/hooks/session-start.sh`
- `skills/do-work/hooks/pipeline-guard.sh` (delete)
- `skills/do-work/docs/standing-preferences.md`
- `skills/do-work/crew-members/anti-slop.md`
- `skills/do-work/crew-members/approach-directives.md`
- `skills/do-work/crew-members/background-agents.md`
- `skills/do-work-knowledge/actions/interview.md`
- `skills/do-work-knowledge/hooks/memory-stop-capture.sh`
- `skills/do-work-knowledge/crew-members/anti-slop.md`
- `skills/do-work-knowledge/crew-members/background-agents.md`
- `skills/do-work-toolbox/actions/ai-report.md`
- `skills/do-work-toolbox/actions/note.md`
- `skills/do-work-toolbox/actions/present-work.md`
- `skills/do-work-toolbox/actions/tutorial.md`
- `skills/do-work-toolbox/docs/ai-report-guide.md`
- `skills/do-work-toolbox/actions/slop-check.md`
- `skills/do-work-toolbox/actions/stray-check.md`
- `skills/do-work-toolbox/docs/stray-check-guide.md`
- `skills/do-work-toolbox/crew-members/anti-slop.md`
- `skills/do-work-toolbox/crew-members/approach-directives.md`
- `skills/do-work-toolbox/crew-members/background-agents.md`
- `decisions/records/adr-005-pipeline-is-stateful-and-resumable.md`
- `decisions/records/adr-001-modular-action-prompts-and-companion-references.md`
- `decisions/records/adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles.md`
- `decisions/records/adr-007-close-the-pipeline-with-present-and-a-technical-debrief.md`
- `decisions/records/adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats.md`
- `decisions/records/adr-019-four-skill-suite-contract.md` (new retirement hunk only; preserve the pre-existing attestation hunk unstaged)
- `decisions/topics/_index_workflow-orchestration.md`
- `decisions/topics/_index_pipeline-deliverables.md`
- `decisions/topics/_index_skill-architecture.md`
- `decisions/_master_index.md`
- `decisions/log.md`

**Files I will NOT touch:** `do-work/user-requests/UR-031/input.md`, `do-work/queue/REQ-146-remove-modular-migration-shims.md`, historical changelog entries, audits, archived REQs, KB raw evidence, migration progress, unrelated generic pipeline language, and owner-only release files before Step 9.

**Acceptance criteria (restated from REQ):**
- [ ] Stateful pipeline action/reference, routes/aliases, `pipeline.json` lifecycle, guard hook, session-start reporting, tests, and stale live docs are removed.
- [ ] `do-work run` retains triage, implementation, built-in testing, review, lessons, archival, and commit behavior.
- [ ] Core help and README contain the exact approved UR-031 full-cycle prompt.
- [ ] The prompt captures the request, records its UR, verifies that UR, runs its REQs with built-in tests/review, invokes `do-work-toolbox present-work`, and stops/reports on failure.
- [ ] Documentation explains that testing already lives inside `do-work run`; no replacement testing state or stage is introduced.
- [ ] Existing installs lose only the retired guard and preserve unrelated/custom Stop hooks.
- [ ] No pipeline alias or renamed replacement state remains.

## Decisions

- **D-01 — DECIDE & STATE:** Extend the declared scope to `slop-check.md`, `stray-check.md`, and `stray-check-guide.md`. Exploration initially treated their “pipeline completion/health” language as possibly generic, but the builder's final retirement sweep proved those passages specifically instruct against the deleted public workflow. Updating them is a reversible wording correction directly required by “stale docs”; leaving them would preserve broken live guidance.

## Pre-Flight

**Git:** ⚠ One pre-existing project-file modification (`decisions/records/adr-019-four-skill-suite-contract.md`). Preserve its existing attestation hunk and exclude that hunk from this REQ's commit; a separate retirement hunk in the same file is permitted by the declared scope.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh` passed before implementation (`launched: true`).
**Dependencies:** ✓ Repository shell/Go tooling needed by the planned checks is present or has a documented optional fallback.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `.gitignore` (modified)
- `README.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `decisions/_master_index.md` (modified)
- `decisions/log.md` (modified)
- `decisions/records/adr-001-modular-action-prompts-and-companion-references.md` (modified)
- `decisions/records/adr-005-pipeline-is-stateful-and-resumable.md` (modified)
- `decisions/records/adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles.md` (modified)
- `decisions/records/adr-007-close-the-pipeline-with-present-and-a-technical-debrief.md` (modified)
- `decisions/records/adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats.md` (modified)
- `decisions/records/adr-019-four-skill-suite-contract.md` (modified)
- `decisions/topics/_index_pipeline-deliverables.md` (modified)
- `decisions/topics/_index_skill-architecture.md` (modified)
- `decisions/topics/_index_workflow-orchestration.md` (modified)
- `skills/do-work-knowledge/actions/interview.md` (modified)
- `skills/do-work-knowledge/crew-members/anti-slop.md` (modified)
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified)
- `skills/do-work-knowledge/hooks/memory-stop-capture.sh` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/note.md` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `skills/do-work-toolbox/actions/slop-check.md` (modified)
- `skills/do-work-toolbox/actions/stray-check.md` (modified)
- `skills/do-work-toolbox/actions/tutorial.md` (modified)
- `skills/do-work-toolbox/crew-members/anti-slop.md` (modified)
- `skills/do-work-toolbox/crew-members/approach-directives.md` (modified)
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/docs/stray-check-guide.md` (modified)
- `skills/do-work/SKILL.md` (modified)
- `skills/do-work/actions/abandon.md` (modified)
- `skills/do-work/actions/capture.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `skills/do-work/actions/kb-lessons-handoff.md` (modified)
- `skills/do-work/actions/pipeline-reference.md` (deleted)
- `skills/do-work/actions/pipeline.md` (deleted)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/crew-members/anti-slop.md` (modified)
- `skills/do-work/crew-members/approach-directives.md` (modified)
- `skills/do-work/crew-members/background-agents.md` (modified)
- `skills/do-work/docs/standing-preferences.md` (modified)
- `skills/do-work/hooks/hooks.json` (modified)
- `skills/do-work/hooks/pipeline-guard.sh` (deleted)
- `skills/do-work/hooks/session-start.sh` (modified)
- `skills/do-work/next-steps.md` (modified)
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `tools/install-do-work-suite.sh` (modified)

**What was done:** Removed the separate stateful pipeline runtime, aliases, state reporting, and installed Stop guard; replaced the public workflow with the exact approved capture → verify → `do-work run` → toolbox presentation prompt. Installer reconciliation now removes only the retired guard while preserving custom Stop hooks in jq and Python paths, and live docs/decisions/tests reflect the prompt-based successor without weakening `do-work run`.

## Qualification

Passed — 50 project files verified (47 modified, 3 deleted), all seven acceptance criteria traced, and P-A-U confirmed. The changes are substantive and net subtractive; installer hook data remains dynamically reconciled rather than stubbed. Overlap with REQ-144 is expected and explicit through `depends_on`/`related` in the same modular-cutover batch, so the contamination check found no unexplained prior-REQ carryover. Mechanical qualification passed after normalizing the ADR-019 manifest verb; the unrelated pre-existing ADR-019 attestation hunk remains separately identified for partial staging.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; `go test ./...` in `skills/do-work-board/tools/queue-kanban`; Bash syntax on both installers, changed hooks, and three contract suites; installer byte comparison; hooks JSON parse; `git diff --check`.
**Result:** ✓ All passing. The contract suite reported suite-manifest, commit-hash, updater, staged-skill, installer-behavior, and full regression success; Go tests passed; all syntax/identity/JSON/diff guards exited 0.

**Red-green validation:**
- `_dev/tests/staged-skills-contract.sh`: ✗ while the three runtime files, aliases/state lifecycle, fresh Stop guard, and missing prompt remained → ✓ after retirement and exact-prompt publication.
- `_dev/tests/install-suite-behavior.sh`: ✗ while fresh installs and jq/Python reconciliation retained `pipeline-guard.sh` → ✓ with targeted legacy-guard removal and mixed custom Stop-hook preservation.
- `_dev/tests/contract-regressions.sh`: ✗ on the new retirement ratchets and nested focused suites → ✓ after the full implementation; failures were assertion-level evidence against the old live behavior, not launch/syntax errors.

**New tests added:**
- Retirement-absence and exact-prompt contracts in the staged-skill and core regression suites.
- Fresh-install and jq/forced-Python legacy Stop-hook reconciliation cases, including mixed custom entries.

**Existing tests updated (cross-REQ impact):**
- The REQ-144-era staged distribution/installer contracts now assert the scheduled post-cutover deletion rather than requiring the one-release pipeline compatibility surface.

*Verified by work action*

## Review

**Overall: 77%** | 2026-08-08T17:45:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 71% |
| Code Quality | 90% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The no-JSON-tool installer path preserves the retired pipeline-guard Stop entry while deleting its script, contrary to the existing-install migration contract (`tools/install-do-work-suite.sh:6`, `skills/do-work/tools/install-do-work-suite.sh:6`, `_dev/tests/install-suite-behavior.sh:378`) — gate: user-visible → REQ-151 created

**Minor findings:** 1 (report only: `skills/do-work/actions/work-reference.md` still calls the work-driven lessons handoff “pipeline mode” rather than “orchestrated mode”)
**Acceptance:** Partial — runtime retirement, exact prompt, jq/Python reconciliation, and regressions pass; manual settings reconciliation remains incomplete.
**Suggested testing:** 2 items (exercise the no-tool mixed custom Stop-hook fallback after REQ-151; repeat the retired-workflow restatement sweep)
**Follow-ups created:** REQ-151; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Ratcheting the deleted paths, aliases, state token, exact successor prompt, and installer behavior before implementation made the intended removal observable as RED and protected `do-work run` from accidental weakening.
- Treating the two installer copies and the duplicated crew guidance as lock-step surfaces kept the modular distribution internally consistent.

**What didn't:**
- The first installer pass covered executable jq and Python reconciliation but left the human no-JSON-tool fallback saying to preserve every hook. That dead path only surfaced during the independent review, not the GREEN automation.

**Worth knowing:**
- Retiring a hook requires three aligned paths: fresh configuration, automated migration, and manual fallback. Tests must assert the retired entry disappears while unrelated/custom hooks survive in each path.
- A narrow stale-surface sweep is safer than banning the word “pipeline”; the repository still legitimately uses it for the `do-work run` orchestration, CI, and data flows.

## Orientation

[MAP CHANGED] End-to-end suite work is now a copyable capture → verify → `do-work run` → `do-work-toolbox present-work` sequence instead of a resumable pipeline state machine. Core `do-work run` remains the sole implementation orchestrator and still owns triage, tests, review, lessons, archival, and commits; `REQ-151` carries the remaining manual installer-fallback cleanup.
