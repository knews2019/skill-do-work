# Changelog

What's new, what's better, what's different. Most recent stuff on top.

---

## 0.69.12 — The Dark Code Kit (2026-04-23)

Captures a three-prompt kit for fighting "dark code" — code that was never understood by anyone at any point in its lifecycle. Shared `dark-code-kit_` prefix groups them as sibling tools in the library.

- `prompts/dark-code-kit_audit.md`: four-group interview (architecture, AI tool usage, team/ownership, deployment) that produces a hotspot map across structural and velocity dimensions, with severity ratings, ownership gaps, and a prioritized action plan.
- `prompts/dark-code-kit_context-layer-generator.md`: per-module interview walking through structural → semantic → philosophical context, emitting a module manifest, behavioral contracts, and a decision log that make the module self-describing.
- `prompts/dark-code-kit_comprehension-gate.md`: senior-engineer-style PR review across seven dimensions (credentials, cross-service side effects, blast radius, state, tokens, assumptions, comprehension) with CLEAR / REVIEW REQUIRED / HOLD verdicts.
- `prompts/README.md`: three entries added to the Available prompts table.

## 0.69.11 — The Ingest Correction (2026-04-23)

Fixes the kb-lessons handoff's user-facing messages: both the no-KB fallback and the promoted confirmation told users to run `bkb triage` alone, but triage only sorts inbox files — compilation into the wiki happens in `bkb ingest`. Following the old messages left lessons stuck in `capture/notes/` and invisible in the wiki.

- `actions/kb-lessons-handoff.md`: no-KB fallback now documents the full re-promotion path — `bkb init` → re-run handoff (e.g. `do-work review REQ-NNN`) → `bkb triage` → `bkb ingest`. Previously it stopped at triage and also glossed over the fact that the handoff set `kb_status: pending` without dropping the file, so even a correct triage+ingest pair would have found an empty inbox.
- `actions/kb-lessons-handoff.md`: "Promoted to …" confirmation now instructs `bkb triage` then `bkb ingest`. Previously users on the happy path were told `bkb triage` was the last step, leaving the lesson sorted but uncompiled.

## 0.69.10 — The Gap Patrol (2026-04-23)

Audit-driven cleanup of the three recent handoff commits (0.69.7–0.69.9). Fills in the spots where the new `kb_status`/`kb_entry` fields and the handoff flow weren't yet mentioned in sibling docs. Nothing behavioral — just the cross-references finally catching up with the feature.

- `actions/work.md`: `## Request File Schema` now documents the two optional `kb_status` and `kb_entry` frontmatter fields alongside the existing ones. Previously only `sample-archived-req.md` mentioned them, so agents reading the schema block thought they were non-standard.
- `next-steps.md`: "After work" and "After review work" blocks now suggest `do-work bkb triage` as a follow-up when lessons were promoted, and `do-work bkb init` when the handoff deferred because no `kb/` existed.
- `actions/build-knowledge-base.md`: `triage` classification table now recognizes `.md` files with `source_type: req_lesson` frontmatter (written by the kb-lessons handoff). They route to `capture/notes/` — no new capture subdir needed — with a note that the `domain` field is a reliable topic hint and `req_path` is a back-reference to the originating REQ.

## 0.69.9 — The Handoff Cleanup (2026-04-23)

Two bot-reviewer findings against the kb-lessons handoff, both legitimate and both fixed. Metadata now populates correctly in pipeline mode, and the `declined` vs `skipped` statuses are actually reachable as designed.

- `actions/kb-lessons-handoff.md`: `date` now falls back to today's date when `completed_at` isn't set yet — the handoff runs at Step 7.5 (pipeline mode), before Step 8 writes `completed_at`, so the old "source from `completed_at`" rule produced empty dates on every pipeline run.
- `actions/kb-lessons-handoff.md`: user's explicit "Skip" choice in Step 3/4 now records `kb_status: declined` instead of `skipped`, matching Step 5's semantics (`declined` = active refusal, `skipped` = silent auto-skip when trigger conditions aren't met). Previously `declined` was effectively unreachable.

## 0.69.8 — The Homegrown Handoff (2026-04-23)

Replaces the compound-engineering integration from 0.69.7 with a zero-dependency version that uses do-work's own knowledge base (`kb/`). After a REQ's review passes and Lessons Learned are captured, do-work drops a structured source document into `kb/raw/inbox/` so the existing `bkb triage` → `bkb ingest` pipeline compiles it into the wiki. Same consent-driven shape as before, just no external plugin required.

- `actions/kb-lessons-handoff.md`: New handoff reference. Writes to `<kb>/raw/inbox/REQ-NNN-<slug>.md`, defers to `kb_status: pending` if no `kb/` exists (never auto-inits), and stops at the drop — triage and ingest stay in the bkb action's lane.
- `actions/review-work.md`, `actions/work.md`: Step 9.5 / Step 7.5 now call the kb-lessons handoff instead of the CE one. Unattended runs default to `kb_status: pending`.
- `actions/sample-archived-req.md`: Frontmatter fields renamed — `ce_compound_status` → `kb_status`, `ce_solution_path` → `kb_entry` (filename only, survives bkb's moves through `capture/` and `processed/`).
- `CLAUDE.md`: "Compound-engineering Integration" section replaced with a shorter "Lessons → Knowledge Base Handoff" section that documents the in-skill flow only.
- Removed: `actions/ce-compound-handoff.md`, `docs/ce-integration-guide.md` — both were CE-specific and no longer apply.

## 0.69.7 — The Compound Handoff (2026-04-23)

First integration point with the [compound-engineering plugin](https://github.com/EveryInc/compound-engineering-plugin). After a REQ's review passes and Lessons Learned are captured, do-work now offers to promote those lessons into CE's `docs/solutions/` knowledge base via the `ce-compound` skill. The handoff asks before dispatching, degrades to a saved prompt if CE isn't installed, and never blocks archival.

- `actions/ce-compound-handoff.md`: New reference file describing the handoff payload shape, user consent flow, and REQ frontmatter updates. Both review-work Step 9.5 (standalone) and work Step 7.5 (pipeline) dispatch into this single reference.
- `actions/review-work.md`: Step 9.5 now runs the compound handoff after lesson capture in standalone mode.
- `actions/work.md`: Step 7.5 now runs the compound handoff after lesson capture in pipeline mode. Unattended runs default to `ce_compound_status: pending` — no auto-promotion.
- `actions/sample-archived-req.md`: Sample frontmatter now shows the two new optional fields (`ce_compound_status`, `ce_solution_path`) so REQ authors know the schema.
- `CLAUDE.md`: New "Compound-engineering Integration" section documents the augmentation model, the three CE artifact paths, and the current integration point.
- `docs/ce-integration-guide.md`: New user-facing guide covering install, the handoff flow with sample payload, troubleshooting, roadmap for future integration points (reviewer agents, ce-plan, ce-brainstorm), and design principles for contributors wiring up the next seam.

## 0.69.6 — The Audit Ratchet (2026-04-22)

Close the contradictions and gaps found in a self-audit of the skill: a broken link, a missing `next-steps.md` entry, an out-of-date README, a missing docs guide, two action files that didn't follow the template, and a wave of missing `When to Use` / `Red Flags` / `Verification Checklist` sections across core actions. Nothing behavioral — just the docs finally matching the conventions CLAUDE.md claims.

- `decisions/records/adr-012-interview-v2-gap-closure.md`: Fixed broken link to v1 spec — the filename is date-prefixed (`2026-04-16_expand-skill-do-work-interview.md`).
- `next-steps.md`: New `After interview` blocks covering session-in-progress, all-layers-complete, export, and list — previously absent despite `interview` being a first-class action.
- `README.md`: `bkb` usage list now includes `defrag`, `garden`, and `crew [action]` — all three were already in the action file and `next-steps.md`, just missing from the README overview.
- `CLAUDE.md`: `prompts/` tree entry now points at `prompts/README.md` as the authoritative index instead of listing one outdated prompt.
- `docs/prompts-guide.md`: New guide for the prompts dispatcher — sub-commands, name resolution, safety model, and how to add a new prompt.
- `actions/install-ui-design.md`, `actions/install-bowser.md`: Restructured to follow the CLAUDE.md action template (When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist).
- `actions/capture.md`, `actions/clarify.md`, `actions/work.md`, `actions/pipeline.md`, `actions/ui-review.md`, `actions/prompts.md`, `actions/present-work.md`, `actions/prime.md`, `actions/version.md`, `actions/tutorial.md`, `actions/build-knowledge-base.md`, `actions/forensics.md`, `actions/deep-explore.md`, `actions/scan-ideas.md`: Added missing `When to Use`, `Red Flags`, and/or `Verification Checklist` sections per CLAUDE.md's action-template spec. All 14 core actions now carry the full template.

## 0.69.5 — The Hyphen Hustle (2026-04-22)

Every `do work` command invocation is now written `do-work` across docs, actions, crew rules, and the session-start hook. Matches the skill's actual name and makes it unambiguous to agents that it's a real command, not a verb phrase.

- All `*.md` files and `hooks/session-start.sh`: `do work <action>` → `do-work <action>`, including README examples, SKILL.md routing tables, action files, docs, CHANGELOG prose, crew rules, prompts, and decision records.
- No behavior change — natural-language triggering still works; the skill's name has always been `do-work`, so hyphenated references stay consistent with the skill manifest.

## 0.69.4 — The Review Ratchet (2026-04-17)

Follow-up to 0.68.2: fixes three defects from code review on the interview v2 gap-closure patch. One was a JSON rendering bug, one was a reference to a session field that doesn't exist, and one was a stale-entry leak into agent rules that violated ADR-012's own promise. ADR-012 gets a "Post-merge corrections" section documenting each.

- `interviews/work-operating-model.md`: `operating-model.json` template now uses `{{json_entries <layer>}}` instead of `[ "{{canonical_entries}}" ]` — emits a proper JSON array of entry objects instead of a single-element array of strings.
- `interviews/work-operating-model.md`: All `{{session.completed_at}}` references changed to `{{session.last_exported_at}}` (the field that actually exists on `session.json`).
- `actions/interview.md`: `export` sub-command reordered to stamp `last_exported_at` in-memory **before** rendering (step 2), then persist `session.json` after artifacts are on disk (step 4). Prevents templates from substituting a null timestamp on first export.
- `interviews/work-operating-model.md`: SOUL.md and HEARTBEAT.md templates now filter `where status != "stale"` on every entry-iterating block. USER.md's active sections do the same, plus a new "Stale or deprecated" section labels stale entries at the bottom (narrative context preserved, but they no longer appear as active rules).
- `actions/interview-reference.md`: Ingest frontmatter `created:` fields follow the template fix (`last_exported_at` in place of `completed_at`).
- `decisions/records/adr-012-interview-v2-gap-closure.md`: "Post-merge corrections" section added under Consequences.

## 0.69.3 — The Honored Flag (2026-04-17)

Fixes an inconsistency in the eval-harness prompt flagged in code review: `--tasks <n>` was documented but the interview and output flow were hard-coded to exactly three test cases. The prompt now resolves N from the flag up front (default 3, clamped to 1–7) and uses N everywhere — task inventory, priority selection, case count, verification.

- `prompts/prompt-kit-step5-eval-harness.md`: new Step 0 resolves and clamps N; Steps 1, 2, 3, 5 reference N instead of literal 3; Rules and Verification Checklist enforce the contract; Red Flags call out suite-size drift; template placeholder for the per-case index changed from `[N]` to `[#]` to avoid visual collision with the count variable.

## 0.69.2 — The Topical Shelving (2026-04-17)

Regroups the five AI-industry analytical prompts by the discipline they're drawn from — business, economics, or tech — dropping the redundant `ai-` umbrella (the whole library is AI-oriented). One of the tech prompts gains an `infrastructure` sub-prefix to mark it as an infra decision rather than an architecture one.

- `prompts/ai-vendor-strategic-sort.md` → `prompts/business-vendor-strategic-sort.md`
- `prompts/inference-economics-stress-test.md` → `prompts/economics-inference-stress-test.md`
- `prompts/saas-repricing-exposure.md` → `prompts/economics-saas-repricing-exposure.md`
- `prompts/compute-geography-risk.md` → `prompts/tech-infrastructure-compute-geography-risk.md`
- `prompts/inference-architecture-decision.md` → `prompts/tech-inference-architecture-decision.md`
- `prompts/README.md` index rows updated to match. Historical references in `CHANGELOG.md` left as-is.

## 0.69.1 — The Spelled-Out Name (2026-04-17)

Renames the ADR-log prompt so its filename actually says what it does. Establishes a `[noun]_[action]` convention (underscore between the subject and the verb phrase) that leaves room for sibling actions on the same noun later.

- `prompts/adr-log.md` → `prompts/architecture-decisions-log_create-or-expand.md`: renamed; H1 and aliases updated inside the file (`adr`, `adr-log`, `decisions`, `architecture-decisions` all still work as documentation hints — the dispatcher resolves via prefix match against the new filename).
- Cross-references updated in `SKILL.md`, `CLAUDE.md`, `README.md`, `actions/prompts.md`, and `prompts/README.md`. Historical references in `CHANGELOG.md` left as-is.

## 0.69.0 — The Seven Steps (2026-04-17)

Extracts the Prompt Kit article's progression into the library as seven numbered prompts. One pre-flight pen-and-paper exercise plus six runnable disciplines — diagnostic, context doc, spec engineer, intent framework, eval harness, constraints — all `step[n]`-prefixed so they sort in workflow order.

- `prompts/prompt-kit-step0-pen-and-paper-exercises-to-prepare-prompt.md`: handoff prompt that tells the user to step away from the screen and work the seven questions offline, then structures the returning notes into a PRE-FLIGHT BRIEF.
- `prompts/prompt-kit-step1-four-discipline-diagnostic.md`: scored audit across Prompt Craft, Context, Intent, Specification — with a 4-month personalized roadmap.
- `prompts/prompt-kit-step2-personal-context-doc.md`: seven-domain interview producing the user's "CLAUDE.md for everything."
- `prompts/prompt-kit-step3-spec-engineer.md`: collaborative spec builder for real projects — acceptance criteria, constraint architecture, task decomposition, definition of done.
- `prompts/prompt-kit-step4-intent-and-delegation-framework.md`: extracts implicit decision rules into a deployable framework, with a Klarna Test self-check.
- `prompts/prompt-kit-step5-eval-harness.md`: Lütke-pattern test suite over the user's actual recurring tasks.
- `prompts/prompt-kit-step6-constraint-architecture.md`: pre-delegation Must Do / Must Not / Prefer / Escalate document tied to the user's stated failure modes.
- `prompts/README.md`: index updated with all seven new entries.

## 0.68.2 — The Paved Cowpath (2026-04-17)

Closes five v1 gaps in the `interview` action per the v2 imported spec — export templates move into the template file as mechanical render templates, `update` goes entry-level, mid-layer quits become recoverable, and `ingest` lands 10 files in `kb/raw/inbox/` instead of inventing its own frontmatter shape. Surgical patches, not a rewrite. Recorded as ADR-012.

- `interviews/work-operating-model.md`: New `## Export Templates` section with verbatim handlebars-style render templates for `USER.md`, `SOUL.md`, `HEARTBEAT.md`, `operating-model.json`, and `schedule-recommendations.json`. An implementation can now render exports mechanically against the approved session — different runs produce the same file shape.
- `actions/interview-reference.md`: `## Export Schemas` trimmed to framework-level invariants only (narrative tone, source-confidence filtering, cadence, traceability). Template-specific rendering now lives in the template.
- `actions/interview-reference.md`: `update` re-run mode rewritten to walk entries individually — `[confirm / edit / mark-stale / delete / skip]` per entry. Explicitly overrides v1's "do not invent a per-entry patch path." CHANGELOG format for update runs is now `N confirmed, N edited, N marked stale, N deleted, N added`.
- `actions/interview-reference.md`: New `### Mid-layer recovery` section. On resume, the action checks for `.draft-<layer-id>.md` written opportunistically during the interview and offers pick-up vs. start-over.
- `actions/interview-reference.md`: `## Ingest Frontmatter` rewritten as `## Ingest File Mapping`. Specifies 5 export files + 5 layer summaries = 10 files per run for `work-operating-model`, plus a manifest row per file in `kb/raw/_inbox_queue.md`. Frontmatter aligns with BKB's canonical schema (`sources:` list, `related:` with `rel`, `type: source-summary` for exports, `type: concept` for layer summaries).
- `actions/interview.md`: New draft-checkpoint step in the layer interview workflow. Subsequent steps renumbered. `ingest` sub-command body rewritten to reference the new File Mapping section in the reference.
- `decisions/records/adr-012-interview-v2-gap-closure.md`: New ADR documenting the five patches. Extends ADR-011. Crew placement audit confirmed `crew-members/interviewer.md` stays put — the directory is a generic persona pool, not `work`-scoped.
- `decisions/_master_index.md` + `decisions/topics/_index_skill-architecture.md`: Bumped to list ADR-012.

## 0.68.1 — The Rename Tag (2026-04-16)

Renames the Weekly Structural Diff prompt so "original" is explicit in the filename — clears the way for variant versions of the same framework to coexist in the library.

- `prompts/weekly-structural-diff.md` → `prompts/weekly-structural-diff-original.md`: renamed; index entry in `prompts/README.md` updated to match. Invoke with `do-work prompts run weekly-structural-diff-original` (prefix match `weekly-structural-diff` still resolves unambiguously while it's the only variant).

## 0.68.0 — The Promptkit Drop (2026-04-16)

Six new reusable prompts ingested from the Prompt Kit article on the 2026 capability-phase → economics-phase transition. They turn the article's analytical framework into runnable tools for tracking AI news, stress-testing product economics, mapping infrastructure risk, pricing SaaS seat compression, sorting vendors, and designing inference architectures.

- `prompts/weekly-structural-diff.md`: Signal/noise sort for AI news across five altitudes (physics, monetization, geography, business models, geopolitics), with a "what didn't change" calibration and prioritized takeaways.
- `prompts/inference-economics-stress-test.md`: Sora-style economics stress test — sustainability ratio, three-scenario pressure test, emoji verdict (🟢/🟡/🟠/🔴), and a concrete "what would fix it" plan. Benefits from a thinking-capable model.
- `prompts/compute-geography-risk.md`: Location-by-location risk matrix (power/grid, permitting/politics, geopolitics, data residency) with a deployment strategy and contingency playbook per location.
- `prompts/saas-repricing-exposure.md`: Seat compression estimate, "The Clock" (months until compression shows in reported numbers), transition readiness score, and an Atlassian benchmark.
- `prompts/ai-vendor-strategic-sort.md`: Vendor assessment matrix across five structural-sustainability dimensions, one tripwire event per vendor, and a portfolio concentration score.
- `prompts/inference-architecture-decision.md`: API vs. self-hosted vs. hybrid comparison, model selection matrix, Sora test, and a Now / 3× / 10× migration path with triggers.
- `prompts/README.md`: Index table extended with the six new prompts.

## 0.67.5 — The Weekly Witness (2026-04-17)

New prompt in the library: `weekly-signal-diff` — a weekly structural diff of AI-industry news, personalized via BKB. Ships with a 10-lane core starter universe and auto-loads a personal sidecar at `prompts/weekly-signal-diff-personal.md` when present for user-specific lanes. Every loaded lane gets full coverage every week — no lane is ever compressed or dropped.

- `prompts/weekly-signal-diff.md`: New prompt. Produces both an inline digest and a durable deliverable at `do-work/deliverables/weekly-signal-diff/<week-ending>.md` staged for BKB ingest. Idempotent per week-ending date (appends timestamped revisions rather than overwriting). Supports `--week-ending`, `--source-packet`, `--topic`, `--dry-run`, `--no-ingest`. Aliases: `wsd`, `signal-diff`.
- `prompts/weekly-signal-diff-personal.md`: New placeholder template. Ships with no real lanes — users copy it anywhere in their project (project root, `.claude/`, `do-work/`, etc.) and fill in real lanes. At Phase 3 the main prompt searches the user's project and loads whatever project-local copy it finds; the shipped placeholder is only a template, never treated as a source of real lanes. Library prompt and shipped placeholder stay generic; personal content lives exclusively in the user's project.
- `prompts/README.md`: New rows for `weekly-signal-diff` and `weekly-signal-diff-personal` in the Available prompts table.
- `decisions/imported-specs/2026-04-16_weekly-signal-diff-authoring-prompt.md`, `decisions/imported-specs/2026-04-17_starter-universe.md`: Spec updates — demotion language removed, 3–7 shift cap removed, forbidden-memory-layer name-drops stripped, personal sidecar pattern documented.

## 0.67.4 — The Gap Sealer (2026-04-16)

Folds in the legitimate improvements from a parallel branch that landed alongside 0.67.2/0.67.3. The earlier "Unified Trunk" merge tried to combine both lines but truncated `CHANGELOG.md` and rewrote `actions/version.md` losing the global-install guard and the recap section — that merge was reverted and only the load-bearing changes were re-applied here.

- `actions/version.md`: Widened the auto-update dirty check to scope every shipped editable path (`prompts/`, `interviews/`, `specs/`, `docs/`, `decisions/`, `hooks/`, `CLAUDE.md`, `AGENTS.md`, `next-steps.md`) — anything tar would clobber. Anything dirty in those paths now blocks the update.
- `actions/version.md`: New pre-clean step (4) for `prompts/` and `interviews/` — top-level `.md` files are deleted before extraction so upstream-removed entries don't linger as ghost workflows in `do-work prompts list` / `do-work interview list`. Subsequent steps renumbered 4→5, 5→6, 6→7.
- `actions/interview-reference.md`: `update` re-run mode now tracks an in-memory `any_edits` flag. If any layer's approval committed a non-zero diff, the export gate state (`review_completed_at`, `review_runs`) is cleared on completion — the user must re-run `review` before the next `export`. Pure re-confirms leave the gate untouched.
- `actions/interview-reference.md`: `fresh` and `version` empty session shapes now include `last_activity_at: <now>` so the freshness preflight has something to compare against on the very first export.
- `actions/interview.md`: Exports gate rule documents that `update` clears the review state when edits are committed.
- `interviews/work-operating-model.md`: Layer 1 schema fix — `time_windows` entries gain a required `days` field (weekday abbreviations) so `schedule-recommendations.json` can emit `days` without inventing data; `interruptions` is now a list of `{source, priority}` objects (priority drawn from `low`/`medium`/`high`) so `HEARTBEAT.md`'s "What to ignore" section has a real signal to filter on. Template version bumped 1.0.0 → 1.1.0.

## 0.67.3 — The Right Shelf (2026-04-16)

Moves the 0.67.2 export freshness stamp out of `exports/` and into `session.json.last_exported_at`. The sidecar-file approach would have been picked up by `ingest`'s "for each file" loop and polluted `kb/raw/inbox/` with bogus timestamp documents. Caught in review; the field-on-session.json approach was always the right one.

- `actions/interview.md`: `export` preflight and stamp-write now read/write `session.json.last_exported_at` instead of a sidecar file. Empty session shape gains the new field.
- `actions/interview-reference.md`: `session.json` schema gains `last_exported_at`. Status Vocabulary row updated with a note explaining why the stamp lives on the session, not in `exports/`. `fresh` re-run mode writes `last_exported_at: null` in the new empty session.

## 0.67.2 — The Status Ledger (2026-04-16)

Interview recipe gains a stale-export warning and a consolidated status vocabulary — small operational patches for when an operating model gets re-run in anger. Addresses gaps surfaced by a recent design review of the `work-operating-model` activation path.

- `actions/interview.md`: `export` sub-command now stamps `exports/.exported_at` after each run and does a freshness preflight on the next run — if `session.json.last_activity_at` is newer than the stamp, the user hears about it before exports are regenerated.
- `actions/interview-reference.md`: New Status Vocabulary table consolidates the four independent status fields (session `status`, layer `approved`, entry `status`, export freshness stamp) into a single reference. Explicitly notes that prior runs are archived directories, not `superseded` flags.
- `actions/interview-reference.md`: `update` re-run mode now documents the "empty a layer" path (user can nuke a layer; same approval gate applies, empty layer still counts as approved) and calls out that per-entry edit friction is intentional — the approval gate is the whole point.

## 0.67.1 — The Settled Tenant (2026-04-16)

Interview action now works the moment the skill is installed into a project, and session state lives in `do-work/` alongside the rest of the per-repo workspace — tracked in git like URs and REQs.

- Templates resolve from `<skill-root>/interviews/` (the `interviews/` directory inside the skill bundle), not the user's project root. Fixes `do-work interview list` and `do-work interview <template>` finding nothing when the skill ships from `~/.claude/skills/do-work/`.
- Session state moved from `./interview/<template>/` to `./do-work/interview/<template>/`. It joins `queue/`, `user-requests/`, `archive/`, and `working/` under the canonical workspace and is tracked in git — the elicited operating model is durable per-repo knowledge, not transient orchestration state.
- Removed the stale `interview/` entry from the skill repo's own `.gitignore` so the skill no longer models the wrong behaviour.

## 0.67.0 — The Open Ear (2026-04-16)

New `interview` action — a generalized elicitation framework that runs prescriptive templates to turn tacit work knowledge into agent-ready operating artifacts. First template `work-operating-model` walks the five-layer Work Operating Model (Nate B. Jones and Jonathan Edwards) across ~45 focused minutes and produces `USER.md` / `SOUL.md` / `HEARTBEAT.md` plus machine-readable exports. Session state is resumable, cross-layer contradictions get surfaced explicitly, and exports flow into BKB via `ingest` for querying.

- `actions/interview.md`: New sub-command dispatcher — `list`, `<template>`, `<template> status`, `<template> review`, `<template> export`, `<template> ingest`, `<template> reset`, `<template> versions`. Session state lives at `./interview/<template>/session.json` and persists across sessions per ADR-005. Export gates on all layers approved + at least one review pass complete. Re-run modes (`fresh`, `update`, `version`) archive prior runs as immutable `versions/v<N>-<date>/` directories.
- `actions/interview-reference.md`: Companion per ADR-001 holding the heavy content — template file format, canonical 11-field entry contract, `session.json` schema (including `review_completed_at` + `review_runs` gate fields), checkpoint format, per-export schemas for the five `work-operating-model` artifacts, re-run mode specifications, versioning scheme, and ingest frontmatter shape.
- `interviews/work-operating-model.md`: First template. Five layers — operating rhythms, recurring decisions, dependencies, institutional knowledge, friction — each with concrete prompt patterns and layer-specific `details` shape. Declares four named cross-layer contradiction checks the `review` sub-command surfaces.
- `crew-members/interviewer.md`: New persona loaded during every interview sub-command. Concrete-before-abstract, one-question-at-a-time, checkpoint-gated, honest-confidence standards. Never invents fields the user didn't provide.
- `docs/interview-guide.md`: Onboarding guide — when to run (45-minute focused session), the five export files, re-run cadence (quarterly), BKB integration flow, multi-repo context separation, and troubleshooting.
- `decisions/records/adr-011-interview-framework-with-prescriptive-templates.md`: New ADR documenting the prescriptive-not-minimal template shape, single-instance-per-repo design, and local-files-only constraint. Depends on ADR-001, ADR-002, ADR-005; complements ADR-010.
- `SKILL.md`: Registered in action list, routing table (priority 19), Verb Reference, Action Dispatch table, bare-invocation help menu (new Interviews block), and foreground-dispatch list. Frontmatter `argument-hint` updated.
- `README.md`: New numbered scenario "19. Run a structured interview"; renumbered later scenarios 20→21, 21→22, 22→23.
- `decisions/_master_index.md`, `decisions/_progress.md`, `decisions/topics/_index_skill-architecture.md`: ADR-011 added to the index and topic cluster; progress tracker bumped to `Next ADR number: ADR-012`.
- `.gitignore`: New `interview/` line so per-repo session state isn't accidentally committed. Templates under `interviews/` remain tracked.

## 0.66.1 — The Local Landlord (2026-04-16)

The `do-work update` flow now refuses to overwrite a global/shared install. If `SKILL.md` lives under `~/.claude/skills/`, `~/.gemini/skills/`, or anywhere else outside the current project's git root, the update stops and tells the user to either `cd` into the owning project or install the skill locally — no more silent updates to a user-wide copy.

- `actions/version.md`: Added an explicit preflight location check as step 2 of the update flow that resolves the skill root, compares it to `git rev-parse --show-toplevel`, and refuses to proceed if the skill sits under a user-wide skills directory. Renumbered the dirty-tree / run / verify / report steps accordingly. The curl command is now prefixed with `cd <skill-root> &&` so extraction can't land in a global directory by mistake. The fetch-failed fallback message was rewritten to call out the global paths by name.

## 0.66.0 — The Four Corners (2026-04-16)

Deliverables now follow an unambiguous naming convention: `.marp.md` for LLM-authored Marp source, `.marp.html` for the marp-cli export of that source, and `.single.html` for LLM-authored standalone HTML (explainer or debrief). The pipeline now ships four files per completion — three LLM renderings plus the mechanical Marp HTML export — so a stakeholder without marp-cli can still view the deck.

- `actions/pipeline.md`: Step 5 table expanded to four rows — `.md`, `.marp.md`, `.marp.html`, `.single.html` — with the Marp HTML export marked as mechanically produced by `npx @marp-team/marp-cli ... --html`. Narrative, rationalizations, red flags, and the verification checklist updated to distinguish the three LLM renderings from the fourth mechanical export, and to scope the Tailwind/Mermaid CDN constraint specifically to `.single.html`.
- `actions/pipeline-reference.md`: Section 3 heading + filename renamed to `.single.html`. Section 2 (Marp Slide Deck) now calls out the `.marp.html` export with the exact command. Sibling-link lists, preview commands, and the HTML Related-deliverables card grid updated to link both `.marp.html` and `.single.html` where relevant.
- `actions/present-work.md`: Interactive explainer renamed to `{UR-NNN}-interactive-explainer.single.html` with a note explaining the `.single.` vs `.marp.` distinction. Client-brief "Related Reading" footer and the terminal summary updated to the new filenames, and the Keep-exploring footer now links both pipeline summary formats.

## 0.65.2 — The Dry-Run Reprieve (2026-04-15)

Fixes two review findings on the `adr-log` prompt. Phase 0 no longer hard-blocks every run on `main`/`master` — `--dry-run` now skips the tree/branch blockers entirely (they're zero-risk in a read-only run), and non-dry-run invocations on `main` pause and ask for authorization instead of refusing outright. README's description of the prompt's source model was stale; it now accurately reflects the layered spine (`implementation-history.md` primary, `lessons-learned/` secondary, code verification, `CHANGELOG.md` fallback).

- `prompts/adr-log.md`: Rewrote Phase 0 to parse flags first, skip dirty-tree / branch-name blockers under `--dry-run`, and prompt for authorization on `main`/`master` (with three accepted responses: yes / feature-branch-name / no). Authorization persists across resume via `authorized_main_branch: true` in `_progress.md`. Updated the "Never push to main/master" and "`--dry-run` means read-only" rules to match. Added two new Common Rationalization rows (--no-push on main is still a write; dry-run can't skip source verification).
- `README.md`: Replaced the stale "mines CHANGELOG.md for load-bearing decisions" description in scenario 19 with the current layered source model (implementation-history primary, lessons-learned secondary, code verification, CHANGELOG fallback).

## 0.65.1 — The Layered Spine (2026-04-15)

Rewrote `prompts/adr-log.md` to merge the better ideas from the user's own ADR-extraction prompt with the safety envelope from the first draft. Same prompt, much sharper — layered source mining with `implementation-history.md` as the primary spine, REQ/UR-keyed idempotency instead of fuzzy CHANGELOG-version matching, proper YAML `related: [{page, rel}]` relationships, per-cluster `topics/_index_*.md` wiki pages, and a completion report that forecasts remaining work sized S/M/L.

- `prompts/adr-log.md`: Replaced the mining spine (`CHANGELOG.md` → `implementation-history.md` primary, `lessons-learned/` secondary, current code for verification, `CHANGELOG.md` as portable fallback). Replaced the frontmatter schema (now `req:`, `ur:`, `sources:`, `related: [{page, rel}]`, `confidence`). Moved ADR files into `decisions/records/` and clusters into `decisions/topics/_index_<cluster>.md` as first-class wiki pages. Added explicit supersession workflow that flips the old ADR's `status` and adds the inverse `rel: superseded-by` to its `related` list in the same commit. Commit messages now follow `docs(adr): …` conventional shape. Added a completion-report section with a remaining-candidates forecast (sized S/M/L per UR). Kept the pre-flight safety checks, `--dry-run` / `--no-push` / `--batch-size` / `--from` flags, "infer alternatives if absent and mark `(inferred)`" guidance, and the Common Rationalizations / Red Flags / Verification Checklist guardrails.
- `prompts/README.md`: Updated the `adr-log` description to reflect the layered source model and REQ/UR-based idempotency.

## 0.65.0 — The Prompt Shelf (2026-04-15)

New `prompts` action — a dispatcher over a growing library of reusable, battle-tested prompts for recurring jobs the skill doesn't have a first-class action for. Seeded with `adr-log`, a create-or-update prompt that builds a project-wide Architecture Decision Record log at `decisions/` (BKB wiki pattern) by mining `CHANGELOG.md` for load-bearing decisions. Idempotent, resumable, supersession-aware.

- `actions/prompts.md`: New sub-command dispatcher (`list`, `show <name>`, `run <name>`, shorthand `<name>`) that resolves prompt names against `prompts/*.md` by exact match or unambiguous prefix. `show` is strictly read-only; `run` adopts the body below the `---` separator as operational instructions.
- `prompts/README.md`: Library index explaining the prompt file shape (title + blockquote + metadata + `---` + body) and how to add new entries.
- `prompts/adr-log.md`: First library entry. Detects create-vs-update mode via `decisions/_master_index.md`, resumes from `_progress.md` mid-run, allocates sequential `ADR-NNNN` numbers without reuse, handles supersession (sets `status: superseded` + `superseded_by` on the old ADR, never deletes), de-duplicates on re-run via a `source:` frontmatter field, and commits+pushes in batches (scaffolding → mining → ADRs in groups of 3 → final reconciliation).
- `SKILL.md`: New priority-19 routing row for `prompts` / `prompt`, new Verb Reference entry, new Action Dispatch entry, new "Prompt library:" block in the bare-invocation help menu, and `prompts` added to the foreground-dispatch list.
- `next-steps.md`: Three new post-action sections (`prompts list`, `prompts show`, `prompts run`).
- `README.md`: New numbered scenario "19. Run a saved prompt"; renumbered later scenarios 19→20, 20→21, 21→22.
- `CLAUDE.md`: Registered `actions/prompts.md` and the `prompts/` directory in the Project Structure tree.

## 0.64.1 — The Companion Split (2026-04-13)

`actions/pipeline.md` had grown past the 10k-token read limit, which meant agents couldn't load it in one pass. Extracted the three Pipeline Completion Report rendering templates (markdown / Marp / HTML) plus their composition rules into a new `pipeline-reference.md` — same pattern as `work.md` + `work-reference.md` and `deep-explore.md` + `deep-explore-reference.md`. Pipeline.md drops from 549 lines to 377; the templates live in a companion file loaded at Step 5 Completion.

- `actions/pipeline-reference.md`: New companion file holding the three renderings (plain markdown template, Marp 11-slide sequence + frontmatter skeleton, HTML 12-section sequence + CDN stack + design requirements) and the seven composition rules that apply across formats.
- `actions/pipeline.md`: The former three-renderings subsection is now a short pointer paragraph listing what the reference contains — rules, markdown skeleton, Marp sequence, HTML sequence. No content lost; every rule, template, and constraint moved verbatim.
- `CLAUDE.md`: Registered `pipeline-reference.md` in the project-structure listing.

## 0.64.0 — The Cross-Linked Set (2026-04-13)

Pipeline summaries and present-work deliverables now serve both audiences in every file and link to each other. A stakeholder landing on any summary opens straight into a plain-language "What got built" section before the audit data; a developer landing on the interactive explainer finds commit SHAs and `git show` commands alongside the Before/After demo. Each artifact's footer lists its siblings as clickable relative links grouped by audience — "Start here if you want to understand what was built" vs. "Audit the run" — so readers can drill in regardless of which file a teammate sent them.

- `actions/pipeline.md`: All three summary formats (`.md`, `.marp.md`, `.html`) now open with a "What got built" narrative copied verbatim from the client brief, followed by an optional architecture diagram, then the existing audit sections. Added the rendering to the markdown template, a new slide 2–3 pair to the Marp required sequence (with renumbering), and hero-adjacent sections 2–3 to the HTML required sequence.
- `actions/pipeline.md`: Deliverables section now groups sibling artifacts by audience (understand-what-was-built vs. audit-the-run) and renders them as clickable relative links in markdown, real `<a>` tiles in HTML, and a two-column next-steps slide in Marp. New "Serve both audiences in every file" and "Reuse client-brief content verbatim" composition rules, two rationalization rows covering the duplication and dev/stakeholder split, red flags for summaries that skip the narrative or ship unlinked paths, and checklist items enforcing word-for-word parity with the brief.
- `actions/present-work.md`: Client brief template grows a "Related Reading" footer linking the interactive explainer, video, and pipeline-summary siblings when they exist. Interactive explainer guidelines now require a "For the developer" section with commit SHAs and `git show` blocks, plus a "Keep exploring" navigation card grid to sibling deliverables. Terminal summary notes that artifacts link to each other.

## 0.63.3 — The Retro (2026-04-13)

Agents working in this repo now close multi-turn conversations with a short "how you could have one-shotted this" retrospective when it helps. Not for every reply — only when three-plus clarification turns landed somewhere materially different from a naïve reading of the first ask, with specific phrases the user could have used up front.

- `CLAUDE.md`: New "One-Shot Suggestions (Prompt Retrospectives)" section describing when to offer the retrospective, when to skip it (iterative-by-design work, unfolding user thinking, small tasks), its shape (diagnosis → concrete one-shot prompt in the user's voice → disambiguating phrases with reasons → optional meta-lesson), and framing rules (feedback not self-flagellation, be concrete, surface the receiving-agent-vs-embedded-content split explicitly).

## 0.63.2 — The Triple Render (2026-04-13)

The pipeline debrief now ships in three formats — plain markdown, a Marp slide deck, and a standalone HTML page — all rendered from the same extracted dataset. A 12-REQ pipeline deserves more than one surface: a developer scans the `.md` in a PR, a stakeholder sits through the deck, a non-technical reader browses the HTML. One pass over the data, three files on disk, zero drift.

- `actions/pipeline.md`: Step 5 Completion now renders `{UR-NNN}-pipeline-summary.md`, `.marp.md`, and `.html` from one composition pass. Added format-specific templates and design constraints for each (Marp frontmatter skeleton + required slide sequence, HTML stack limited to Tailwind + Mermaid CDN with light/dark theming).
- `actions/pipeline.md`: New composition rule — the three renderings must carry identical facts, no format-specific editorializing. Added rules, rationalization rows, red flags, and checklist items enforcing parity across formats and flagging the common skip-the-HTML shortcut.

## 0.63.1 — The Debrief (2026-04-13)

Pipeline completion now educates instead of just checking a box. After the six steps finish, the pipeline assembles a technical debrief — Final summary table (REQ/commit/scope/one-line), Test state before→after, Cross-REQ coherence highlights from the review, Carry-forward candidates, Deliverables, and a copy-pasteable How-to-verify recipe — and persists it to `do-work/deliverables/{UR-NNN}-pipeline-summary.md`. Long pipelines deserve a digest, not a checkmark.

- `actions/pipeline.md`: Rewrote Step 5 Completion to assemble and save the Pipeline Completion Report. Added the report format to Output Format with composition rules (cite commits, pull from primary sources, flag missing baselines, never auto-capture carry-forward). Added completion-status block with Duration/Branch/Verdict metadata.
- `actions/pipeline.md`: New Rule on completion-as-education, two new Common Rationalization rows (hollow completion + invented baselines), and Red Flag / Verification Checklist additions covering missing report sections and fabricated metrics.
- `actions/present-work.md`: Added a "How to Verify" section to the client brief template so non-technical readers also get a concrete validation recipe.

## 0.63.0 — The Closing Act (2026-04-13)

The pipeline now closes the loop. Added `present` as the sixth step so a full pipeline run ends with client-facing deliverables (brief, architecture diagram, video, HTML explainer) — no more remembering to run `do-work present` manually after every pipeline.

- `actions/pipeline.md`: Added `present` step after `review` — dispatches to the present work action with the UR ID from the capture step's artifacts. Skips gracefully if capture produced no artifacts.
- `actions/pipeline.md`: Updated state schema, status block example, help menu, dispatch table, completion check (5 → 6 steps), Rules, and Common Rationalizations to include the new step.
- `SKILL.md`: Updated pipeline description and help menu to reflect the six-step sequence.
- `README.md`: Updated the pipeline section to mention `present` in the full cycle.
- `next-steps.md`: Post-pipeline suggestions now point at `present all` (portfolio mode) instead of the per-UR brief that's already been generated.

## 0.62.5 — The Few Words (2026-04-12)

New crew member: caveman mode. Tag a REQ with `caveman: true` (or `caveman: lite|full|ultra`) and the builder compresses prose output ~65-75% while keeping code and technical terms exact. Adapted from JuliusBrussee/caveman.

- Added `crew-members/caveman.md` — token-efficient communication rules with three intensity levels (lite, full, ultra) and auto-clarity escape hatch for security warnings
- `actions/work.md`: Step 6 agent rules loading now includes caveman.md conditional on `caveman` frontmatter
- `CLAUDE.md`: Documented caveman crew member loading behavior

## 0.62.4 — The Dangling Pointer (2026-04-12)

Fixed a dangling cross-reference in `bkb init` Step 5 — it pointed to a "Schema File" section that had been extracted out into `bkb-reference.md`. Init now correctly points to the schema content's real home.

- `actions/build-knowledge-base.md`: Step 5 now references the "Schema File Content" section of the bkb-reference action, matching the pattern already used in Steps 3 and 4

## 0.62.3 — The Same Rake (2026-04-12)

Gap-closure pass after 0.62.2. The earlier release caught a Rules-at-end ordering bug in `quick-wins.md` and fixed it there, but never grepped for the same pattern elsewhere — turns out two other action files had the identical structure. Also formalized a capture.md deviation that was noticed during the first review but never fixed.

- `actions/review-work.md`, `actions/verify-requests.md`: Moved `## What NOT to Do` (functionally a Rules section) from the very end back up before `## Common Rationalizations`, matching the CLAUDE.md template order
- `actions/capture.md`: Renamed `## Core Rules` → `## Philosophy` so the opening invariants block lives in the template's pre-When-to-Use slot instead of pretending to be Rules before Steps

## 0.62.2 — The Own Medicine (2026-04-12)

Ran `do-work code-review` on the skill itself and actually fixed the findings. Four first-class actions (`pipeline`, `scan-ideas`, `deep-explore`, `tutorial`) were missing from the README usage scenarios despite being prominent in the SKILL.md help menu — now every listed action has a README section. Also cleared the remaining template-ordering drift and filled in missing Red Flags blocks.

- `README.md`: Added four new usage-scenario sections (`pipeline`, `scan-ideas`, `deep-explore`, `tutorial`) and renumbered to 21 scenarios — closes the discoverability gap between the help menu and the README
- `actions/scan-ideas.md`, `actions/review-work.md`, `actions/deep-explore.md`: Swapped `## Philosophy` and `## When to Use` to match the CLAUDE.md template order
- `actions/quick-wins.md`, `actions/ui-review.md`: Moved `## Rules` back before `## Common Rationalizations` / `## Verification Checklist`
- `actions/inspect.md`: Renamed `## Core Rules` → `## Rules` and repositioned it after Output so it matches the template's post-Output placement
- `actions/code-review.md`, `actions/inspect.md`, `actions/ui-review.md`: Added the missing `## Red Flags` section each had been skipping — brings the encouraged-elements coverage up to parity with peers
- `actions/work.md`: Replaced three `[text](./file.md)` cross-references with short-name prose (CLAUDE.md rule: "SKILL.md owns the file-path mappings")
- `hooks/session-start.sh`: Warn on stderr when the version parse falls back to "unknown" instead of silently hiding format drift
- `hooks/pipeline-guard.sh`: Documented that the jq-absent fallback is best-effort and depends on well-formed JSON
- `CLAUDE.md`: Softened the `docs/` description — not every action has a per-action guide (install-*, tutorial, scan-ideas, deep-explore, pipeline, clarify intentionally rely on the action file + README)

## 0.62.1 — The Senior Engineer Test (2026-04-12)

Refined the Karpathy crew-member and wired its principles into review-work, so the four guardrails aren't just applied during the build — they're audited during review. Also added an oversimplification hedge, because "simplify" is not "strip."

- `crew-members/karpathy.md`: Sharpened Success Indicators with four concrete observable behaviors (clarifying questions first, small diffs, untouched neighbors, verification language)
- `crew-members/karpathy.md`: Added "Simplify ≠ strip" clarification under Simplicity First — foundation should not be removed just because it could be
- `actions/review-work.md`: Added Karpathy Principle Check as an informational pass in Step 6 Code Review — a mnemonic audit against the four principles, without double-penalizing issues already caught by existing dimensions

## 0.62.0 — The Karpathy Nod (2026-04-12)

Adopted Andrej Karpathy's four coding guardrails as an always-loaded crew-member, so every REQ — not just multi-agent waves — benefits from them. Complements do-work's workflow machinery: the queue decides *what* to build; these principles shape *how*.

- `crew-members/karpathy.md`: New file — four behavioral principles (Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution) adapted from [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills)
- `actions/work.md`: Step 6 now always-loads `karpathy.md` alongside `general.md`
- `CLAUDE.md`: Agent Rules section documents the new always-loaded file

## 0.61.3 — The Finer Edges (2026-04-12)

Round 2 of self-run quick-wins — structural template nits and documentation accuracy. Every action file now either matches the CLAUDE.md template or is documented as an accepted variant.

- `actions/deep-explore.md`: Wrapped 9 steps under a `## Steps` parent and demoted them from H2 to H3 (largest template deviation remaining after 0.61.2)
- `actions/cleanup.md`: Renamed `## What It Does` → `## Steps` (Pass 0/1/2/3 sub-sections keep their domain-appropriate "Pass" terminology)
- `CLAUDE.md`: Clarified crew-members description (not all files are domain-gated); added `approach-directives.md` loading rule to the Agent Rules list
- `CLAUDE.md`: Added `tutorial.md` and `forensics.md` to the Accepted Variants section (multi-mode with dispatcher + checklist-based diagnostic patterns)
- `actions/work.md`: Completed the `specs/` example list in Step 3.7 — now names all 4 shipped spec templates instead of just 2

## 0.61.2 — The Quick Sweep (2026-04-12)

Fixes from a self-run `do-work quick-wins` on the skill itself. Consistency nits the team would notice before users ever would.

- `next-steps.md`: Added missing `**After <action>:**` blocks for `cleanup`, `install-ui-design`, `install-bowser` — SKILL.md's "suggest next steps after every action" rule now holds for every action
- `actions/deep-explore.md`: Renamed second `## When to Use` (a comparison table, not a use-case section) to `## Scan-Ideas vs Deep-Explore` — no more duplicate headers
- `crew-members/general.md`: Added the `JIT_CONTEXT` comment convention the other 8 crew files already follow (always-loaded during Step 6)
- 9 action files: Renamed `## Workflow` → `## Steps` to match the CLAUDE.md template's "Required elements: Steps (numbered)" — `capture`, `commit`, `inspect`, `install-bowser`, `install-ui-design`, `review-work`, `ui-review`, `verify-requests`, `work`
- `_dev/code-review-20-commits.md`: Marked resolved (both findings already addressed in later versions)

## 0.61.1 — The Lean Cut (2026-04-11)

Trimmed low-value additions from 0.61.0 and split the largest action file. Guardrails stay where they earn their token cost; template bloat moves to a companion file.

- Removed guardrail sections (rationalizations, checklists) from 5 low-stakes actions: forensics, scan-ideas, prime, present-work, clarify — these are read-only reporting actions where the guardrails restated what the Steps already say
- Removed Role Identity sections from 3 crew-member files — a motivational paragraph doesn't change behavior when 150+ lines of specific rules follow
- Removed `CONTRIBUTING.md` (CLAUDE.md already serves as the contributor guide) and `docs/skill-anatomy.md` (same)
- Split `build-knowledge-base.md` (1077→687 lines) — extracted seed file templates, agent crew definitions, and KB schema into new `bkb-reference.md` companion, following the work.md/work-reference.md pattern

## 0.61.0 — The Bright Standard (2026-04-10)

Quality guardrails, routing clarity, and a session-start hook across the skill.

- 7 action files: Added Common Rationalizations tables, Red Flags sections, and Verification Checklists (capture, cleanup, commit, inspect, verify-requests, pipeline, quick-wins)
- 10 action files: Added "When to Use / When NOT to Use" sections to commonly confused routes (review-work, code-review, verify-requests, inspect, commit, cleanup, forensics, scan-ideas, quick-wins, deep-explore)
- `hooks/session-start.sh` + `hooks/hooks.json`: SessionStart hook injects version, pending REQ count, and pipeline status
- `CLAUDE.md`: Updated action file template with When to Use, Common Rationalizations, Red Flags, and Verification Checklist conventions
- `README.md`: Added token efficiency guidance and hooks installation section

## 0.60.5 — The Honest Mirror (2026-04-10)

Cross-file contradiction audit — fixes 13 inconsistencies spanning stale paths, duplicate codenames, missing scoping mechanisms, and documentation gaps.

- `actions/scan-ideas.md`, `actions/deep-explore.md`: Fixed stale `do-work/` queue path → `do-work/queue/` (missed by v0.60.3)
- `CHANGELOG.md`: Renamed 15 duplicate codenames (The Safety Net ×3, The Consistency Pass ×4, The Signpost ×2, The Compass ×2, The Cartographer ×2, The Feedback Loop ×2, The Gap Closer ×2, The Full Picture ×2, The Clarity Pass ×2) — each entry now has a unique codename
- `CHANGELOG.md`: Corrected v0.60.2 entry that claimed `do-work/` root was canonical (subsequently reversed by v0.60.3)
- `actions/work.md`: Added Input section with `$ARGUMENTS` support for targeted REQ IDs (e.g., `do-work run REQ-042`) — fixes pipeline scoping gap where pipeline.md told work to process specific REQs but work had no mechanism to accept that constraint
- `SKILL.md`: Updated work verb reference to document optional REQ ID arguments; updated priority 4 routing to accept trailing REQ IDs
- `actions/commit.md`: Documented commit message format distinction from work.md (`Traced-to:` vs `Implements:`) and added commit pathway deconfliction note
- `actions/verify-requests.md`: Fixed Step 3 to include `do-work/queue/` while keeping `do-work/` as legacy fallback
- `actions/review-work.md`: Added P-A-U checkbox verification to the Verification Checklist
- `specs/README.md`: Fixed `--spec` hint claim → `suggested_spec` frontmatter field (matches actual capture.md implementation)
- `actions/deep-explore.md`: Added `surviving_directions` and `total_directions_explored` to Step 8 state.json update
- `CLAUDE.md`: Added missing `docs/` directory and `AGENTS.md` to project structure listing

## 0.60.4 — The Vivid Voice (2026-04-10)

Enriched all four deep-explore subagent persona prompts — Free Thinker, Grounder, Writer, Explorer — from dry task specs into rich, conversational creative briefings with distinct voices, example phrases, and clear "what to avoid" guidance.

- `actions/deep-explore-reference.md`: Rewrote Free Thinker persona with divergent identity, "do NOT self-censor" directive, exploration dimensions, and example phrases
- `actions/deep-explore-reference.md`: Rewrote Grounder persona as brainstorm partner (not analyst), with taste-driven instincts, enthusiasm for good ideas, and direct example reactions
- `actions/deep-explore-reference.md`: Rewrote Writer persona with invisible-observer identity and philosophy about preserving agents' actual language
- `actions/deep-explore-reference.md`: Rewrote Explorer persona with tenacious-researcher identity, focused report structure, and "facts only" boundary

## 0.60.3 — The Paved Path (2026-04-10)

Pending REQ files now live in `do-work/queue/` instead of `do-work/` root. The `queue/` path is what people kept writing instinctively — paving the cow path prevents a recurring class of stale-path bugs.

- `actions/work.md`: All queue glob patterns, directory diagrams, REQ placement paths, crash recovery, and git staging updated to `do-work/queue/`
- `actions/capture.md`: REQ output paths, duplicate-check scans, addendum destinations, and all example outputs updated to `do-work/queue/`
- `actions/cleanup.md`: Sweep globs, report messages, git staging, and relocation paths updated to `do-work/queue/`
- `actions/pipeline.md`, `actions/clarify.md`, `actions/forensics.md`, `actions/review-work.md`, `actions/version.md`, `actions/code-review.md`: Queue scan paths and REQ placement references updated
- `CLAUDE.md`: Queue Path Convention section updated to document `do-work/queue/` as canonical location
- `README.md`, `docs/capture-guide.md`, `docs/work-guide.md`, `docs/cleanup-guide.md`, `docs/forensics-guide.md`: Directory diagrams and path references updated

## 0.60.2 — The Clean Ledger (2026-04-10)

Changelog and action file hygiene — fixes that prevent recurring errors.

- `CHANGELOG.md`: Fixed duplicate version numbers (two entries for 0.52.0, two for 0.51.8) by renumbering displaced entries to 0.51.8–0.51.11 in correct monotonic order
- `CHANGELOG.md`: Renamed 4 duplicate codenames (The Tight Scope → The Narrow Pipe, The Crew → The Agent Crew, The Safety Net → The Guard Dog, The Second Brain → The Knowledge Forge)
- `actions/scan-ideas.md`: Fixed header from "Ideate Action" to "Scan-Ideas Action" (missed in v0.57.0 rename)
- `actions/scan-ideas.md`, `actions/deep-explore.md`: Updated queue path references (subsequently moved to `do-work/queue/` in v0.60.3)
- `CLAUDE.md`: Added version dedup guard (verify new version > existing before committing) and codename uniqueness check
- `CLAUDE.md`: Added Queue Path Convention section (subsequently updated to `do-work/queue/` in v0.60.3)

## 0.60.1 — The Clear Head (2026-04-10)

Reverted wave-based pipeline processing — it duplicated what `do-work run` already handles natively (sequential queue draining with fresh agents per REQ). Pipeline Step 5a is back to the original simple continuation loop.

- `actions/pipeline.md`: Removed wave-based processing (Step 5a.1), wave output formats, wave rules. Restored original Step 5a with 3-cycle cap.

## 0.60.0 — The Many Lenses (2026-04-10)

Per-agent approach directives for multi-REQ processing. When sub-agents work on parallel or sequential REQs, each gets a distinct implementation lens (Correctness-First, Simplicity-First, etc.) to improve solution diversity and reduce convergent thinking.

- `crew-members/approach-directives.md`: New file — 8 implementation lenses, assignment rules, and sub-agent context template
- `actions/work.md`: Added approach directive assignment before sub-agent dispatch in Step 6
- `actions/review-work.md`: Added Directive Alignment Check in Step 6 — evaluates whether the assigned lens was applied and flags blind spots

## 0.59.0 — The Quality Blueprint (2026-04-10)

New `specs/` directory with reusable specification templates for common task types. Specs define output structure, quality standards, implementation checklists, and common pitfalls — loaded automatically during work when the REQ matches a template.

- `specs/`: New directory with README and four templates: `api-endpoint.md`, `ui-component.md`, `refactor.md`, `bug-fix.md`
- `actions/work.md`: Added Step 3.7 (Spec Loading) — checks `specs/` for matching templates after triage, passes guidance to builder and reviewer
- `actions/capture.md`: Added optional `suggested_spec` frontmatter field and spec hint inference during parsing
- `CLAUDE.md`: Updated project structure to include `specs/` directory

## 0.57.1 — The Tidy Sweep (2026-04-10)

Quick-wins cleanup: shell script hardening, broken link fix, and next-steps consolidation.

- `hooks/pipeline-guard.sh`: Quoted command substitution (line 27), replaced `2>/dev/null` error suppression with numeric validation on PENDING comparison (line 53)
- `actions/capture.md`: Fixed broken relative link on line 149 — replaced with inline code path
- `next-steps.md`: Expanded generic bkb entry into 11 per-sub-command next-step blocks (moved from `build-knowledge-base.md`)
- `actions/build-knowledge-base.md`: Removed embedded next-steps section (88 lines) — canonical source is now `next-steps.md`

## 0.57.0 — The Deep Dive (2026-04-10)

New `do-work deep-explore` action for multi-round structured exploration of concepts. Spawns divergent/convergent subagent dialogue (Free Thinker, Grounder, Writer, optional Explorer) to develop seed ideas into vision documents and idea briefs. Also renames `ideate` to `scan-ideas` for clarity — `ideate` still works as a trigger keyword.

- `actions/deep-explore.md`: New action — multi-round exploration with session directories, continue mode, convergence rubric, and 4 subagent roles
- `actions/deep-explore-reference.md`: Companion file — persona prompts, document templates, state schema, error handling
- `actions/ideate.md` → `actions/scan-ideas.md`: Renamed for clarity (quick scan vs deep exploration)
- `SKILL.md`: Add deep-explore routing (priority 21), rename ideate → scan-ideas (priority 20), update verb reference, help menu, action dispatch, subagent config
- `CLAUDE.md`: Update project structure for scan-ideas, deep-explore, deep-explore-reference
- `next-steps.md`: Add post-deep-explore suggestions, update post-scan-ideas suggestions

## 0.56.2 — The Tight Scope (2026-04-10)

Two fixes to pipeline queue continuation (Step 5a) from PR review feedback.

- Continuation reviews now always target individual REQ IDs — removed UR shortcut that would re-review all completed REQs under a UR, not just the current batch
- Error recovery guidance is now context-aware: suggests `do-work review REQ-NNN` when review fails (since processed REQs are already completed and `do-work run` would no-op), and `do-work run` only when the run step itself failed

## 0.56.1 — The Safety Net (2026-04-10)

Three gaps closed in the pipeline queue continuation (Step 5a).

- Error handling for continuation: if run or review fails mid-continuation, report the error, print progress, and stop — don't retry or update `pipeline.json`
- Max iteration cap: continuation loop limited to 3 cycles to prevent runaway loops from review-generated follow-ups
- Explicit review targeting: continuation now records pending REQ IDs before dispatching run, then passes them to the review action by ID (or by shared UR)

## 0.56.0 — The Clean Sweep (2026-04-10)

Pipeline now drains the full queue after completing its primary request. If pending REQs remain (from prior captures, follow-ups, or review-generated work), the pipeline automatically continues with run + review cycles until the queue is empty.

- `actions/pipeline.md`: Added Step 5a (Queue Continuation) — scans for remaining `status: pending` REQs after pipeline completion and processes them in a loop
- `actions/pipeline.md`: Added continuation notice to Output Format section and drain rule to Rules section
- `next-steps.md`: Updated pipeline completion label to reflect queue-drained state

## 0.55.0 — The Outside Eye (2026-04-10)

Enriched security, accessibility, and testing guidance after reviewing the claude-skills-collection catalog and cross-referencing with Trail of Bits skills, claude-a11y-skill, and testing-anti-patterns approaches.

- `crew-members/security.md`: New "Static Analysis Tooling" section — tool detection table (CodeQL, Semgrep, Bandit, Brakeman, gosec), what SAST catches vs misses, variant analysis concept, guidance to use project's existing tools
- `actions/ui-review.md`: New "Automated Accessibility Tooling" subsection in Step 7 — tool detection for eslint-plugin-jsx-a11y, axe-core, Pa11y with run commands and integration guidance
- `crew-members/testing.md`: Three new anti-patterns — test-per-method symmetry, catch-all assertions, ignoring test output

## 0.54.1 — The Sharp Eye (2026-04-09)

Fix three bugs in v0.54.0 crew-member additions caught by PR review.

- `crew-members/testing.md`: Rust detection no longer requires `[dev-dependencies]` — any `Cargo.toml` is sufficient. RSpec pattern fixed from `*.test.rb` to `spec/*_spec.rb, .rspec`.
- `crew-members/performance.md`: Reverted JIT_CONTEXT to match actual work.md loader rules — removed aspirational "backend API" loading claim that was never wired up.

## 0.54.0 — The Test Bench (2026-04-09)

New testing crew member and enhanced domain knowledge for performance/observability and async/concurrency. Inspired by patterns from the wshobson/agents plugin marketplace — distilled into do-work's platform-agnostic crew-member format.

- `crew-members/testing.md`: New "Verifier" crew member — test framework detection, testing pyramid guidance, mocking boundaries, fixture patterns, flaky test prevention, TDD workflow, and anti-patterns. Loads on `tdd: true`, `domain: testing`, or after 2+ test failures
- `crew-members/performance.md`: Added observability basics section (structured logging, health checks, metric naming, trace context). Broadened loading to include backend API and data-intensive work
- `crew-members/backend.md`: Added async/concurrency section (blocking I/O in async paths, shared state protection, parallel I/O, cancellation) and dependency awareness section (vulnerability checks, lockfile hygiene, pinned versions)
- `actions/work.md`: Updated crew-member loading rules in Step 6 and Step 6.5 to include testing.md
- `CLAUDE.md`: Documented testing.md loading behavior in Agent Rules

## 0.53.2 — The Short Circuit (2026-04-09)

Bare "code review" (no hyphen, no scope) now routes to `code-review` instead of falling through to `review-work`. No more surprise routing.

- `SKILL.md`: Move bare "code review" from priority 9 (review-work) to priority 7 (code-review), update verb reference, remove help menu warning

## 0.53.1 — The Mirror Check (2026-04-09)

Fixes two documentation gaps from 20-commit audit: adds Performance dimension to code-review guide, surfaces routing distinction in help menu.

- `docs/code-review-guide.md`: Add Performance Anti-Pattern Scan section
- `SKILL.md`: Add UX note about "code review" vs "code-review" routing
- `_dev/code-review-20-commits.md`: Updated review — 2 valid findings, 2 false positives dismissed

## 0.53.0 — The Spark (2026-04-09)

New `do-work ideate` action — generates grounded ideas for what to build, improve, or explore next. Scans prime files, project history, TODOs, coverage gaps, and codebase patterns to produce ranked suggestions with effort estimates. Every idea references something concrete in the code.

- `actions/ideate.md`: New action with 7 idea categories (features, improvements, performance, DX, reliability, integrations, docs), size tags (S/M/L), and confidence levels
- `SKILL.md`: Add ideate routing (priority 20), verb reference, help menu entry, action dispatch, subagent config
- `CLAUDE.md`: Add ideate.md to project structure
- `next-steps.md`: Add post-ideate suggestions

## 0.52.3 — The Full Map (2026-04-09)

Tutorial's "File structure" topic now covers the knowledge base layout (raw/, wiki/, agents/) alongside the do-work/ directory.

- `actions/tutorial.md`: Expand Topic 8 guidance to include KB directory structure

## 0.52.2 — The Plain Prompt (2026-04-09)

Tutorial now uses plain text menus instead of the ask-user tool. The ask tool caps at 4 options, which truncated the 8-topic interactive tour. Menus are printed as text and the agent waits for the user to reply naturally.

- `actions/tutorial.md`: Replace ask-tool requirement with plain text print-and-wait pattern in mode selection, tour topic selection, and rules

## 0.52.1 — The Tidy Menu (2026-04-09)

Moved tutorial to a single line in the "Maintenance & info" section, right before `help`. Keeps the help menu compact.

- `SKILL.md`: Consolidate tutorial from separate "Learn" section into one line before `help`

## 0.52.0 — The Onboarding (2026-04-09)

New `do-work tutorial` command with four modes: quick-start (hands-on walkthrough), concepts (mental model explainer), recipes (scenario → command cheat sheet), and interactive tour (menu-driven deep dives). Bare invocation asks which mode to run.

- `actions/tutorial.md`: New multi-mode tutorial action with Quick Start, Concepts, Recipes, and Interactive Tour
- `SKILL.md`: Add tutorial routing (priority 21), verb reference, help menu entry, action dispatch, subagent config
- `CLAUDE.md`: Add tutorial.md to project structure
- `next-steps.md`: Add post-tutorial suggestions

## 0.51.11 — The Guard Rails (2026-04-09)

Strengthens anti-rationalization guards, adds verification checklists, and deepens crew member guidance — inspired by patterns from addyosmani/agent-skills.

- `actions/work.md`: Expanded anti-rationalization table from 4 to 9 rows in Step 6.3
- `actions/code-review.md`: Added Common Rationalizations table and Verification Checklist
- `actions/review-work.md`: Added Common Rationalizations, Red Flags, and Verification Checklist
- `actions/ui-review.md`: Added Common Rationalizations and Verification Checklist
- `actions/quick-wins.md`: Added performance/security smell scanning (Steps 3 + 3.5) and Common Rationalizations
- `crew-members/frontend.md`: Expanded with animation perf, error handling depth, and frontend security
- `crew-members/backend.md`: Expanded with API resilience and performance awareness
- `crew-members/performance.md`: New crew member covering Core Web Vitals, backend optimization, and bundle analysis

## 0.51.10 — The Help Desk (2026-04-09)

Per-command help — any action now supports `do-work <command> help` to show a brief usage summary. Actions with sub-commands (pipeline, prime, bkb) already handled this; all other actions now generate a compact summary from their action file. Footer line added to the main help menu to advertise the feature.

- `SKILL.md`: Add "Per-Command Help" section with rendering template and dispatch rules
- `SKILL.md`: Add tip footer to help menu

## 0.51.9 — The Trim Down (2026-04-09)

Condensed the help menu from ~80 lines to ~35. Removed duplicate entries, collapsed BKB sub-commands into a single line, reduced per-action examples, and merged related sections.

- `SKILL.md`: Help menu compressed — grouped related actions, eliminated duplicates (clarify listed twice), collapsed 12 BKB examples into inline sub-command list

## 0.51.8 — The Safe Exit (2026-04-09)

Fix pipeline-guard stop hook crashing when jq is unavailable and no pipeline is active.

- `hooks/pipeline-guard.sh`: Add `|| true` to grep fallback for `active` field so a no-match doesn't trigger `set -e` exit

## 0.51.7 — The Cross-Check (2026-04-08)

Fixes stale references and underspecified instructions found during a 20-commit code review.

- `actions/code-review.md`: Fix stale "see Step 9" → "see Step 10" after step renumbering from perf-audit fold
- `actions/code-review.md`: Security severity mapping now explicitly includes Critical → Critical
- `actions/pipeline.md`: Session ID increment logic simplified — single-file state can't track prior IDs
- `actions/pipeline.md`: Clarify `investigate` step completes immediately when no uncommitted changes exist
- `actions/pipeline.md`: Fix contradictory rule about `run` step queue scope
- `actions/verify-requests.md`: Add scoring formula for per-REQ Overall and Overall Confidence
- `actions/capture.md`: Unify addendum coherence resolution protocol for both queued and in-flight paths

## 0.51.6 — The Narrow Pipe (2026-04-08)

Pipeline hardening — request isolation, synchronous dispatch, and robust gitignore handling.

- `actions/pipeline.md`: `run` step now scoped to captured REQs only — no longer drains the full work queue, preventing unrelated backlog from executing during a pipeline
- `actions/pipeline.md`: All pipeline-dispatched actions run foreground (blocking), overriding SKILL.md's background default for `work` — prevents race between `run` and `review`
- `actions/pipeline.md`: `.gitignore` is now created if absent (previously only appended to existing), ensuring `pipeline.json` is always excluded from commits
- `SKILL.md`: Added pipeline foreground dispatch exception to subagent config

## 0.51.5 — The Full Send (2026-04-08)

End-to-end pipeline orchestration — chain investigate, capture, verify, run, and review in one command with resumable state tracking.

- `actions/pipeline.md` (NEW): Stateful multi-action pipeline with `do-work/pipeline.json` state tracking, resume across sessions, status display, and error recovery
- `actions/pipeline.md`: Explicit sub-agent context passing — each step documents what artifacts and IDs to forward so sub-agents can target the correct UR/REQs
- `actions/pipeline.md`: Pipeline initialization auto-adds state file to `.gitignore` (transient session state, not for version control)
- `hooks/pipeline-guard.sh` (NEW): Optional Claude Code stop hook to prevent agent from stopping mid-pipeline; uses `$CLAUDE_PROJECT_DIR` for robust path resolution
- SKILL.md: Added pipeline routing (priority 3), dispatch entry, help menu section, verb reference, subagent config
- next-steps.md: Added pipeline next-step suggestions
- CLAUDE.md: Added pipeline.md and hooks/ directory to project structure

## 0.51.4 — The Deeper Cuts (2026-04-08)

Cherry-picked five improvements from a Graph-of-Thought analysis of the bkb action — better cross-source awareness, smarter queries, and fewer deferred problems. Also fixed a bug where clustered resolve left contradictions permanently open.

- `build-knowledge-base.md`: Triage now enriches queue entries with `topic_hint` and `priority` fields
- `build-knowledge-base.md`: Ingest detects confidence transitions at ingest time (medium→high on corroboration, high→low on contradiction) instead of deferring to lint
- `build-knowledge-base.md`: Batch ingest cross-references claims across sources — catches agreements and contradictions at merge time
- `build-knowledge-base.md`: Query follows typed relationships up to 2 hops deep for richer multi-source answers
- `build-knowledge-base.md`: Resolve groups related contradictions into clusters and resolves them as a unit to prevent cascading inconsistencies
- `build-knowledge-base.md`: Resolve emits one `[RESOLVED]` marker per original contradiction in a cluster (not one per cluster), preventing ghost re-detection
- `build-knowledge-base.md`: Lint adds a confidence-audit check (flags mismatches between source evidence and confidence level)

## 0.51.3 — The Intent Trail (2026-04-08)

Elevates intent tracking to a first-class concept. REQs are now explicitly framed as validated statements of user intent, not just task descriptions.

- `SKILL.md`: New "Trail of Intent" blockquote — the skill produces a trail of intent, not just code
- `capture.md`: "Validated artifacts" principle — captured REQs are user-validated, not drafts
- `capture.md`: Coherence Rule — addenda must not contradict existing REQ content; conflicts trigger user resolution
- `capture.md`: Coherence across addendum chains — cross-file contradictions flagged before writing
- `capture.md`: "Capture produces validated intent" closing — names the output of capture-phase clarification
- `work.md`: Living log connected to intent trail — builder decisions and scope declarations are intent documentation
- `work.md`: Decisions linked to intent trail — decisions without reasoning are not traceable
- `verify-requests.md`: "REQs are validated intent" philosophy bullet — verify checks validation actually happened
- `verify-requests.md`: Internal Coherence evaluation dimension (0-100%) — catches self-contradictory REQs
- `verify-requests.md`: Coherence column added to verification report table

## 0.51.2 — The One Scale (2026-04-08)

Security findings in code-review now use the same severity scale as the rest of the report (Critical / Important / Minor / Nit) instead of a separate High / Medium / Low scale that had no mapping to follow-up REQ creation.

- Aligned Step 5 security classification to the file's existing 4-level scale with explicit mapping from security.md levels

## 0.51.1 — The Lean Cut (2026-04-08)

Removed standalone `test-strategy` and `perf-audit` actions — their best ideas now live inside `code-review` instead of duplicating scope across three actions.

- Deleted `actions/test-strategy.md` and `actions/perf-audit.md`
- Enhanced **code-review** with new Step 6 (Performance Anti-Pattern Scan) covering N+1 queries, unbounded queries, sequential I/O, bundle bloat, and more
- Enhanced **code-review** Step 7 (Test Coverage Assessment) with risk-driven gap prioritization — flags critical-risk + untested combinations
- Enhanced **code-review** Step 5 (Security) to load `crew-members/security.md` when present
- Cleaned up SKILL.md routing, help menu, action dispatch, and next-steps.md

## 0.51.0 — The Sentinel Suite (2026-04-07)

New audit actions and a security crew member, inspired by techniques from [awesome-prompts](https://github.com/ai-boost/awesome-prompts). Fills gaps in performance diagnosis, test planning, and security review.

- New **crew-members/security.md** — OWASP Top 10 checklist with framework-specific patterns (Node/Express, Python/Django, Java/Spring, React, Go) and severity classification. Loads JIT when working on auth, crypto, or input handling code
- New **actions/test-strategy.md** — risk-driven test strategy designer. Identifies what tests should exist based on risk assessment, gap analysis, and test pyramid health. Includes flaky test prevention and CI quality gate checks
- New **actions/perf-audit.md** — evidence-based performance diagnosis. Scans for backend, frontend, and database anti-patterns (N+1 queries, bundle bloat, missing indexes), quantifies impact, and ranks fixes by effort vs improvement
- Enhanced **crew-members/debugging.md** — added tool-selection-by-failure-class table, Heisenbug identification heuristics, and confidence-level labeling for diagnostic claims
- Enhanced **actions/quick-wins.md** — added objective complexity metrics for tie-breaking (cyclomatic complexity, nesting depth, import count, change frequency), false positive checks, and behavior-preservation rule
- Updated SKILL.md routing table, help menu, action dispatch, and next-steps.md for new actions

## 0.50.5 — The Second Pass (2026-04-07)

Self-review of 0.50.4 patch kit — fixed 5 issues found in our own changes.

- `work.md`: Relative path algorithm reworded to count directory components, not `/` separators (less ambiguous)
- `work.md`: Path verification now specifies failure behavior (report broken link, don't silently write it)
- `work.md`: Step 6 builder instructions now explicitly say to read the D-XX counter before numbering decisions
- `work.md`: Cycle detection rewritten to check the current REQ's existing chain for loops (clearer logic)
- `cleanup.md`: Pass 1 "skip" behavior made explicit (leave UR in `user-requests/` untouched)
- `cleanup.md`: Pass 3a explains why the canonical-location check exists despite Pass 0

## 0.50.4 — The Patch Kit (2026-04-07)

Code review fixes: addressed bugs and ambiguities found across a 20-commit audit. Improves reliability of the core work pipeline and diagnostic actions.

- `work.md`: Clarified crash recovery logic for `pending-answers` restoration (explicit condition check)
- `work.md`: D-XX counter annotation required after Step 3.5 to prevent ID collision with Step 6 decisions
- `work.md`: Expanded "Wired" check exceptions to cover barrel re-exports, dynamic imports, CSS side-effect modules
- `work.md`: Prime test commands are now validated against `package.json`/config before use; stale commands fall back to generic detection
- `work.md`: Added explicit relative-path algorithm (depth-counting) for prime file lesson links
- `work.md`: Cycle detection rewritten with clear traversal algorithm (handles chains of any length)
- `cleanup.md`: Pass 1 now flags duplicate REQ-IDs found in multiple archive locations
- `cleanup.md`: Pass 3a defers to Pass 0 before overwriting a canonical REQ
- `forensics.md`: Stuck work check now includes explicit remediation steps
- `inspect.md`: Renamed "Committed" verdict to "Already Committed" to avoid confusion with quality judgments
- `prime.md`: Area index prime is now defined (filename pattern + required section)
- `SKILL.md`: "scan" routing clarified — bare path vs. path + descriptive text

## 0.50.3 — The Lint Brush (2026-04-07)

Fix 9 bugs and consistency issues found during code review of last 20 commits.

- Fix doubled path in ui-review (`crew-members/crew-members/` → `crew-members/`)
- Fix master index line limit contradiction (50 → 80) in BKB directory tree comment
- Align topic index split threshold to 40 articles everywhere (was 80 in some places, 40 in defrag)
- Add missing defrag/garden staleness warnings to BKB `status` sub-command
- Fix misleading cleanup Pass 2 comment about Pass 1 handling loose REQs
- Add `build-knowledge-base.md` to CLAUDE.md project structure
- Add missing BKB sub-commands (`defrag`, `garden`, `rollup`, `crew`) to SKILL.md help menu
- Clarify SKILL.md routing table to distinguish scoped code-review (priority 6) from unscoped review (priority 8)
- Disambiguate `CLAUDE.md` reference in BKB architect agent to mean KB schema file

## 0.50.2 — The Typo (2026-04-07)

Fixed incorrect queue path in work guide.

- `docs/work-guide.md`: `do-work/requests/` → `do-work/` (REQ files live at the do-work root, not a requests subdirectory)

## 0.50.1 — The Roll Call (2026-04-07)

Named the crew members — each now has a title that reflects their role.

- The Compass (general) — cross-domain orientation, PRIME philosophy
- The Renderer (frontend) — components, state, performance, accessibility
- The Engineer (backend) — APIs, data layer, security, error boundaries
- The Artisan (ui-design) — 6-phase design pipeline, visual craft, interaction specs
- The Detective (debugging) — scientific method, investigation techniques, bias guards

## 0.50.0 — The Crew (2026-04-07)

Renamed `agent-rules/` to `crew-members/` and dropped the `rules-` prefix from all files inside.

- `agent-rules/` → `crew-members/` directory rename
- `rules-general.md` → `general.md`, `rules-frontend.md` → `frontend.md`, etc.
- Updated all references across work.md, ui-review.md, code-review.md, prime.md, review-work.md, version.md, sample-archived-req.md, CLAUDE.md, README.md, and docs/
- Historical CHANGELOG entries preserved as-is (they describe the state at time of release)

## 0.49.0 — The Architect (2026-04-07)

Major clarity and modularization pass — smaller files, consolidated routing, documented conventions.

- SKILL.md verb sections (18 H3s, ~130 lines) collapsed into a single scannable Verb Reference table
- work.md split: orchestrator checklist, error handling, progress template, and common mistakes extracted to `work-reference.md`
- CLAUDE.md now documents accepted action file variants (sub-command dispatchers, multi-mode, state-based) alongside the standard template
- Project structure updated to reflect all new files (`work-reference.md`, `clarify.md`, `next-steps.md`)

## 0.48.0 — The Splitter (2026-04-07)

Modularized the skill for clarity — smaller files, explicit conventions, cleaner separation of concerns.

- Extracted `next-steps.md` from SKILL.md (130 lines of per-action suggestions now live in their own file)
- Extracted `actions/clarify.md` from work.md (clarify questions is now a standalone action file, not a mode of work.md)
- Added Action File Conventions and Agent Rules loading docs to CLAUDE.md
- Added explicit agent-rules loading steps in work.md Step 6 (always load general, conditionally load domain, never block on missing)
- Updated SKILL.md dispatch table: clarify now points to `./actions/clarify.md` instead of `./actions/work.md` with mode flag

## 0.47.5 — The Signpost (2026-04-07)

Gave the most important command in the skill a proper help entry and expanded its guide.

- Help menu now shows description text and trigger aliases for `do-work run`
- Added `do-work continue` and `do-work clarify` to the help menu's "Process the queue" section
- Expanded `docs/work-guide.md` with practical session walkthrough, full alias list, and tips

## 0.47.4 — The Compass (2026-04-06)

Prime action no longer assumes project-specific CLAUDE.md sections or directory layouts.

- Conventions section and Step 5 registry check are now conditional — skip gracefully if CLAUDE.md has no prime registry
- Step 3 adds explicit skip list and CLAUDE.md-guided scan targets
- Cross-cutting prime registration references made generic
- Step 1 skip rules use portable patterns instead of project-specific directory names

## 0.47.2 — The Quoter (2026-04-06)

Fixed invalid YAML frontmatter in SKILL.md that caused strict parsers to fail.

- Quoted `argument-hint` value to escape colons
- Removed non-standard `upstream` frontmatter key

---

## 0.47.1 — The Badge Guard (2026-04-06)

Agent dispatch now gracefully handles missing `agents/` directory instead of failing.

- Main and custom agent dispatch sections skip loading when `agents/` is absent (fresh `init` or legacy KBs)
- Legacy KBs get a clear message pointing to `bkb init --fill-gaps` for migration
- Init `--fill-gaps` documentation clarified as the migration path for pre-v0.46.0 KBs
- Schema crew section notes the guard

## 0.47.0 — The Full Crew (2026-04-06)

BKB imports the best ideas from My-Brain-Is-Full-Crew's 13 skills — adapted for zero-dependency knowledge base operations.

- `bkb defrag` — weekly structural maintenance: re-evaluates cluster boundaries, proposes merges/splits, promotes growing concepts, demotes stale clusters
- `bkb garden` — metadata hygiene: topic cluster balance, relationship type distribution, orphaned indexes, reciprocity checks, reclassification suggestions
- `bkb crew create/list/edit/remove` — custom agent lifecycle: extend the 8 built-in agents with domain-specific roles via guided interview
- Enhanced transcript handling during ingest: multi-speaker detection, key decisions/action items/open questions extraction, topic segmentation, speaker entity pages, structured source summary format
- Status sub-command now reports last defrag/garden dates and custom agent count
- All 8 built-in agents updated with new "When active" entries for defrag, garden, and crew

## 0.46.0 — The Agent Crew (2026-04-06)

BKB gets a crew of 8 specialized agents that define the roles the LLM adopts during each operation. Inspired by My-Brain-Is-Full-Crew's multi-agent architecture, adapted for knowledge base operations.

- 8 agent definition files created during `bkb init` in `kb/agents/`: Architect, Sorter, Compiler, Seeker, Connector, Librarian, Reviewer (QA), Editor (readability)
- Agent dispatch table maps each sub-command to its crew: sequential handoff (→) for ingest pipeline, concurrent standards (+) for lint
- Reviewer agent adds QA gate: confidence auditing, source verification, untested-claim detection
- Editor agent adds readability checks: scannable articles, clear titles, navigation quality, stub detection
- Schema file (CLAUDE.md) updated with crew dispatch section

## 0.45.0 — The Second Brain (2026-04-06)

BKB gets smarter queries, richer connections, and zero-friction capture. Inspired by the best ideas from karpathy-llm-wiki, open-brain-server, claude-knowledge-vault, and My-Brain-Is-Full-Crew.

- Self-improving retrieval agent (`wiki/agent.md`) — learns from past queries to prioritize future lookups, regenerates hot topics every 5 queries
- Three-tier query routing (Synthesize/Record/Skip) — only files answers as wiki pages when they connect 2+ sources, preventing wiki bloat
- Typed relationships in `related:` frontmatter — six relationship types (extends, contradicts, evidence-for, complements, supersedes, depends-on) with 8-per-page cap
- Lint expanded with relationship density, relationship validity, and agent staleness checks

## 0.44.0 — The Cartographer (2026-04-06)

Prime file operations now live inside do-work. `do-work prime create` generates prime files via interactive Q&A; `do-work prime audit` runs a full health check on all primes.

- New `prime` action with `create` and `audit` sub-commands
- Routing at priority 16 (between forensics and BKB)
- BKB, quick-wins, install, capture renumbered to 17-20

## 0.43.5 — The Final Polish (2026-04-05)

Four cleanup fixes from final review pass — two bugs and two documentation gaps.

- Fixed duplicate step number in `status` sub-command
- Fixed stale "today" target reference in sub-commands table
- Companion files (image+md, audio+transcript) explicitly move to `processed/` together as a unit
- Schema file (CLAUDE.md) now includes confidence rules, non-text handling, and `[RESOLVED]` convention

---

## 0.43.4 — The Guard Dog (2026-04-05)

Four defensive gaps closed — init can't clobber existing KBs, queue prunes itself, lint runs are trackable, and confidence levels have clear rules.

- Init: pre-flight check stops if KB already exists, suggests `--fill-gaps` for repair
- Rollup: archives done queue rows older than 30 days to `_inbox_queue_archive.md`
- Lint: appends timestamped entry to `log.md` so `status` can report last-lint date
- Page conventions: confidence heuristic — high (primary/corroborated), medium (single secondary, default), low (no source or contradiction flagged)

---

## 0.43.3 — The Loose Ends (2026-04-05)

Three remaining BKB edge cases tightened up after full review pass.

- Removed stale `ingest today` from help menu (target was already removed)
- Resolve: added `[RESOLVED]` log convention so open vs closed contradictions are distinguishable
- Ingest by path: files outside `capture/` now handled correctly — moved to `processed/` from any source location, queue entry added for traceability

---

## 0.43.2 — The Gap Closer (2026-04-05)

Nine process gaps fixed in the BKB command — removed ghost folders, added collision handling, non-text source support, per-file fault tolerance, contradiction resolution, and more.

- Removed `raw/daily/` and `raw/monthly/` ghost folders (nothing was ever written there)
- Removed dead `ingest today` target (referenced the ghost folder)
- Triage: HHMMSS- prefix on filename collisions in `capture/`
- Ingest: non-text file handling — LLM vision for images, companion transcript required for audio/video
- Ingest: per-file move-and-mark-done (partial failure no longer loses completed work)
- Ingest: `sources:` frontmatter explicitly uses `raw/processed/` paths (stable location)
- New sub-command: `bkb resolve` — interactive contradiction resolution workflow
- Close: now refreshes `wiki/overview.md` and suggests (but doesn't auto-run) git commit
- Schema file updated to match all lifecycle changes

---

## 0.43.1 — The Clean Handoff (2026-04-05)

BKB file lifecycle fixed — ingest now moves sources out of capture/ so triage never re-queues already-ingested files.

- Ingest step 6: **move** (not copy) sources from `capture/` to `processed/{today}/`; filename collisions get an `HHMMSS-` prefix
- Triage step 4: queue updates scoped to files moved in the current pass only (append-only ledger)
- Ingest step 3: duplicate detection before wiki page creation — exact duplicates merge into existing pages, near-duplicates get bidirectional `related:` links
- Manifest table gains "Processed Path" column tracking the final location
- Schema file updated to match new lifecycle semantics

---

## 0.43.0 — The Knowledge Forge (2026-04-05)

New `bkb` command — build and maintain a persistent LLM Knowledge Base wiki compiled from raw source documents. Based on Karpathy's methodology: raw sources go in, structured interlinked Markdown wiki comes out.

- New action: `build-knowledge-base.md` with sub-commands: init, triage, ingest, query, lint, close, rollup, status
- Routing: `do-work bkb [subcommand]`, `do-work build knowledge base`, `do-work knowledge base`, `do-work kb`
- Three-layer architecture: raw pipeline (inbox → capture → daily → processed), wiki (hierarchical indexes → articles), schema file
- Two-hop index navigation for scaling to thousands of articles

---

## 0.42.7 — The Open Sesame (2026-04-05)

Remotion preview script now works around a macOS LaunchServices bug that prevents the browser from opening.

- Present-work Remotion package.json: preview script uses `--no-open` flag + native macOS `open` command to reliably launch the browser

---

## 0.42.6 — The Missing Call (2026-04-05)

Remotion video template now calls `registerRoot()` — fixes the "Waiting for registerRoot()" hang in Remotion Studio.

- Present-work Remotion Root.tsx template: added `registerRoot(RemotionRoot)` call at module level
- Added explicit warning note so agents don't revert to the broken `export const` pattern

---

## 0.42.5 — The Light Switch (2026-04-05)

Interactive explainer now defaults to light theme with automatic OS-level dark mode support via `prefers-color-scheme`.

- Present-work interactive explainer: light theme by default, CSS custom properties for theming, `@media (prefers-color-scheme: dark)` override for dark mode users

---

## 0.42.4 — The Tightened Bolts (2026-04-04)

Five logic gaps fixed from code review — baseline test handling, deleted file inspection, domain propagation, framework routing, and TDD enforcement.

- Work Step 6.5: baseline test failures from Pre-Flight are now excluded from the pass/fail gate — only new regressions require fixing
- Work Step 6.3: "wired" qualification check now exempts framework-convention files (Next.js pages/, SvelteKit routes/, etc.)
- Work Step 6.5: TDD missing-evidence is now a test failure (returns to implementation), not an orphaned qualification flag
- Inspect: committed `(deleted)` files are noted by path without attempting `git show` (which would error)
- Review-work: follow-up REQ template now includes `domain` field to preserve rules-[domain].md loading

---

## 0.42.3 — The Sweeper (2026-04-04)

Cleanup now sweeps finished REQs stranded in queue root or working/ — completed work that fell through the cracks gets archived automatically.

- Cleanup Pass 0: sweep REQs with terminal statuses (`completed`, `done`, `failed`, etc.) from queue root and `working/` to archive
- Cleanup Pass 0: normalize non-standard statuses (`done` → `completed`, `finished` → `completed`, `closed` → `completed`) before archiving
- Forensics Check 9 expanded: "Stranded Finished REQs" now scans both queue root and `working/` for terminal-status REQs

---

## 0.42.2 — The Clear Eyes (2026-04-04)

Queue state is now visible everywhere — no more blind spots where completed work silently disappears.

- Work action prints queue status summary at start (pending/completed/pending-answers counts)
- Work action lists completed-but-unarchived REQs instead of silently saying "no pending REQs"
- Recap reads from both archive and active queue, labels non-archived URs
- Forensics Check 9 detects completed REQs sitting in queue awaiting archive

---

## 0.42.1 — The Full Picture (2026-04-04)

Scoped inspect (`do-work inspect REQ-NNN` / `UR-NNN`) now shows all files from the Implementation Summary — not just uncommitted ones. Committed files are read via `git show` and reported alongside pending changes so you get the complete picture in one report.

- Inspect: committed files from Implementation Summary included when scoped to REQ/UR
- Inspect: committed files get "Committed" verdict (informational, not actionable)
- Inspect: Readiness Summary table gains a Status column (committed / uncommitted)
- Inspect: REQ groups split into Uncommitted and Committed Files subsections
- Inspect: no longer exits early when scoped REQ has no uncommitted files — continues with committed-only inspection
- Inspect: pre-associates all Implementation Summary files in scoped mode, skipping the matching step

---

## 0.42.0 — The Careful Eye (2026-04-04)

Every captured REQ now ends with "Think carefully before answering." — a prompt-level nudge for downstream builders to slow down and reason before acting.

- REQ file template in `capture.md` appends the phrase after the source line

---

## 0.41.0 — The Looking Glass (2026-04-02)

New inspect action explains uncommitted changes — what they are, why they were made, and whether they're ready to commit. Read-only companion to the commit action, with a hybrid narrative+table report format.

- New action: `inspect` — analyzes uncommitted files, traces to REQs, assesses commit readiness
- Three scoping modes: all changes (default), per-REQ, or per-UR
- Six readiness signals: completeness, test coverage, REQ traceability, coherence, safety, improvement hints
- Hybrid report format: narrative What/Why per group + compact readiness summary table
- SKILL.md: routing at priority 12, verbs, examples, dispatch, next-steps
- CLAUDE.md: added inspect.md to project structure

## 0.40.0 — The Audit Pass (2026-04-01)

Code review across last 20 commits uncovered logic bugs, stale docs, and routing gaps. This release fixes them all.

- work.md: Fix commit hash invalidation — use separate metadata commit instead of `--amend` (which changed the hash after writing it)
- work.md: Fix crash recovery overwriting `pending-answers` status — restore original status when REQ had unresolved Open Questions
- work.md: Fix crash recovery strip list — now removes all generated sections (Scope, Pre-Flight, Qualification, Review, Lessons Learned, Decisions, Discovered Tasks), not just the original five
- work.md: Fix UR completion check — now scans `do-work/` root, `working/`, and `archive/` root for REQs belonging to the UR, handles `completed-with-issues` and `failed` statuses explicitly
- work.md: Fix duplicate Implementation Summary — replace existing section on re-qualification instead of appending a second copy
- work.md: Fix checkpoint deletion timing — keep CHECKPOINT.md until crash recovery succeeds, not immediately after reading
- work.md: Fix Route A scope-drift review — skip scope comparison when no Scope declaration exists
- work.md: Specify follow-up REQ requirements for review step — must include `status`, `user_request`, `addendum_to`, `domain`, `review_generated`; cycle detection applies
- SKILL.md: Fix `AskUserQuestion` tool-specific reference → generalized "your environment's ask-user prompt/tool"
- SKILL.md: Add missing verbs to routing table — `audit`, `review requests`, `ui audit`, `design audit`, `clean up`, `commit files`, `save changes`, `scan`, `setup bowser`, `setup playwright`
- SKILL.md: Add `commit`, `forensics` to argument-hint frontmatter
- SKILL.md: Add `commit`, `install-ui-design`, `install-bowser` to subagent foreground list
- capture.md: Replace `AskUserQuestion` with generalized language (2 occurrences)
- rules-debugging.md: Fix hardcoded `error_type: code` — now instructs proper classification per work.md failure table
- rules-general.md: Fix typo "consolididation" → "consolidation"
- CLAUDE.md: Add missing `forensics.md` to project structure
- CHANGELOG.md: Fix wrong year on v0.37.1 (2025 → 2026)

## 0.39.3 — The Clear Intent (2026-03-30)

Use `capture request:` prefix for capture commands so intent is unambiguous. Add `do-work help` as an explicit route.

- SKILL.md: Add `capture request:` as preferred capture prefix in routing, help menu, content signals, and all next-steps suggestions
- SKILL.md: Route `do-work help` to the help menu (priority 1)
- SKILL.md: Add `do-work help` reminder rule to every next-steps block
- README.md: Update capture examples and help section to match

## 0.39.2 — The Clear Menu (2026-03-30)

Rewrote README to organize all features as numbered usage scenarios instead of API-style action docs.

- README: Replace "Actions" reference section with 14 scenario-driven sections (capture, process, verify, review, clarify, code-review, ui-review, quick-wins, present, commit, cleanup, diagnostics, install, version)
- Each scenario leads with the why, then shows exact commands
- Added forensics, ui-review, and install-bowser which were missing or underrepresented

## 0.39.1 — The Bug Sweep (2026-03-30)

Fixes bugs and inconsistencies found during code review of the last 20 commits.

- SKILL.md: Renumber routing priorities to clean integers (no more 2.5/5.5/11.5), fix "check for updates" matching before action verbs
- SKILL.md: Remove duplicate "check for updates" from priority 10 (already handled at priority 2)
- SKILL.md: Add forensics to Action Dispatch table, help menu, and subagent dispatch list
- SKILL.md: Update all priority number references in verb sections to match new numbering
- work.md: Fix architecture diagram — Explore now correctly shown before Scope declare (matches actual step ordering)
- work.md: Fix status flow — only documents actual frontmatter values (pending/claimed/completed/failed), intermediate phases tracked by section presence
- work.md: Add REQ validation in Step 1 — malformed REQs are skipped instead of blocking the work loop
- work.md: Clarify decision numbering — Open Questions and Implementation Decisions share the same D-XX ID space per REQ
- work.md: Add cycle detection for follow-up REQs in Step 8 — prevents infinite addendum chains
- work.md: Clarify Route A lessons learned skip condition with concrete criteria

## 0.39.0 — The Deep Verify (2026-03-30)

Round 2 of quality improvements. The pipeline now self-diagnoses, debugs methodically, validates plans before execution, and tracks decisions for audit trails.

- work.md: Qualify step (6.3) now checks wired (is the file imported?) and flowing (is data real, not stubbed?) — catches dead code and hollow implementations
- work.md: Plan validation (Step 4, Route C) — checks requirement coverage, orphan tasks, scope sanity (5+ tasks = split warning), and file conflicts
- work.md: Pre-flight check (Step 5.75, Routes B/C) — verifies git is clean, tests pass on HEAD, dependencies installed before builder starts
- work.md: Decision traceability — Open Questions and builder decisions numbered D-01/D-02/D-03 in a `## Decisions` section for audit trail
- work.md: TDD task type — `tdd: true` frontmatter flag triggers red-green-refactor cycle with mandatory evidence
- work.md: Test failure loop loads `rules-debugging.md` on attempt 2+; remediation loads it too
- work.md: Session checkpoint enriched with `session_depth` (light/moderate/heavy) — heavy sessions get a Context Summary for next session
- agent-rules/rules-debugging.md (NEW): Full debugging methodology — scientific method, investigation techniques (binary search, minimal repro, working backwards, differential), cognitive bias guards, knowledge capture
- actions/forensics.md (NEW): Read-only pipeline diagnostic — detects stuck work, hollow completions, orphaned URs, scope contamination, failed-without-follow-up, stale pending-answers, git divergence
- SKILL.md: Added forensics routing (forensics, diagnose, health check)
- review-work.md: Scope Discipline now checks for undocumented decisions that changed behavior
- capture.md: Added `tdd` frontmatter field with test-first heuristic in capture workflow

## 0.38.0 — The Quality Gate (2026-03-30)

Borrowed the best quality mechanisms from PAUL and GSD to make the work pipeline self-verifying. Builder claims are now independently checked, failing acceptance tests block archiving, and failures get classified for smarter recovery.

- work.md: New Step 6.3 (Qualify Implementation) — orchestrator independently verifies files exist, changes are substantive, requirements trace, and P-A-U boxes are honest. Includes anti-rationalization rules.
- work.md: New Step 5.5 (Scope Declaration) — Routes B/C declare target files and acceptance criteria before coding. Review compares against this.
- work.md: Step 7 review gating — Acceptance=Fail now triggers a remediation loop instead of silent archiving. New `completed-with-issues` status.
- work.md: Diagnostic failure routing — failures classified as Intent/Spec/Code/Environment with appropriate follow-up REQs.
- work.md: Discovered task severity triage — `[critical]` auto-queues, `[normal]`/`[low]` go through pending-answers.
- work.md: Session checkpoint (CHECKPOINT.md) — written on exit, consumed on resume for session continuity.
- work.md: Context wipe is now structural — fresh agents per REQ, contamination check via file list comparison.
- work.md: Builder and explorer agents now read `## Lessons` from prime files before implementing.
- capture.md: UNIFY checkbox now requires listing files verified and checks performed.
- rules-frontend.md: Populated with component patterns, state management, performance, error handling, quality checks, scope discipline.
- rules-backend.md: Populated with API design, data layer, security baseline, error handling, quality checks, scope discipline.
- review-work.md: Scope Discipline dimension now compares against `## Scope` declaration when present.

## 0.37.4 — The Self-Contained Shell (2026-03-27)

Each code block in install-bowser is now independently runnable — no stale variable assumptions across blocks.

- install-bowser: Recompute `PROJECT_ROOT` in fallback and verify blocks so they work in a fresh shell

## 0.37.3 — The Root Anchor (2026-03-27)

install-bowser now resolves the project root via `git rev-parse --show-toplevel` so the skill installs correctly regardless of cwd.

- install-bowser: All `.claude/skills/` paths anchored to `$PROJECT_ROOT` instead of bare relative paths
- install-bowser: Notes updated to reference `<project-root>/.claude/skills/`
- ui-review: Bowser detection checks from project root for consistency
- SKILL.md: "After ui-review" next steps text now mentions both Playwright CLI and Bowser skill

## 0.37.2 — The Clarity Patch (2026-03-27)

Minor wording fixes for install-bowser and ui-review to make instructions more precise and actionable.

- install-bowser: Clarified fallback curl text — "alternative path in the same repository" instead of misleading "repository's root"
- ui-review: Made Bowser skill detection concrete — check for `.claude/skills/playwright-bowser/SKILL.md` instead of vague "loaded in your environment"

## 0.37.1 — The Route Fix (2026-03-25)

Fixed routing conflicts where commands could dispatch to the wrong action.

- "check for updates" no longer misroutes to verify-requests — added exact-phrase rule at priority 2.5 above Verify
- Added `code review <scope>` and `codebase review` to the code-review priority row so the table matches the verb docs and examples
- Clarified `code review` (no scope) falls through to review-work at priority 6 — annotated in both the table and verb section
- Updated code-review "Read-only" label to "Source-code read-only" with explicit note about optional REQ file creation

## 0.37.0 — The Bowser Install (2026-03-25)

One command to get browser automation. `do-work install-bowser` installs Playwright CLI globally and downloads the Bowser skill from github.com/disler/bowser into the project.

- Added `actions/install-bowser.md` — installs `playwright-cli` (global), Chromium browsers, and Bowser skill (project-scoped from upstream repo)
- Added `install-bowser` routing, verbs, help menu entry, dispatch table, and routing examples in SKILL.md
- `ui-review` now recommends `do-work install-bowser` (instead of raw npm command) when no browser tools detected
- Updated CLAUDE.md, README.md with install-bowser documentation

## 0.36.2 — The Bowser Eye (2026-03-25)

Visual verification now uses Playwright CLI (`playwright-cli`) and the Bowser skill instead of deprecated MCP browser tools. Concrete session-based workflow with viewport screenshots at 320/768/1280px.

- ui-review Step 2.4 detects `playwright-cli` or Bowser skill; recommends `npm install -g @anthropic-ai/playwright-cli@latest` when missing
- ui-review Step 8.5 rewritten with `playwright-cli` session commands: named sessions, `--headless`, viewport via env var, snapshot for a11y, `close-all` cleanup
- Removed all MCP/Puppeteer references — Bowser skill is the browser automation layer

## 0.36.1 — The Rendered Eye (2026-03-25)

Visual verification layer for ui-review. Playwright CLI or browser tools (when available) now screenshot at 320/768/1280px, run axe accessibility audits, and catch rendered-page issues that static code analysis misses.

- ui-review Step 2.4 detects Playwright CLI, browser MCP/skill, or recommends `npm init playwright@latest`
- ui-review Step 8.5 runs visual verification: viewport screenshots, axe audit, rendered layout checks
- Report template gains `Visual verification` status line and `Visual Verification` findings category
- Fixed: missing `design audit` in routing examples
- Fixed: `review work` next steps now suggests `ui-review` for ui-design domain REQs
- Fixed: restored `quick-wins` alongside `ui-review` in code-review next steps
- Added `UI Review` section to README.md
- Removed conflicting `check ui`/`ui check` verbs (consumed by verify-requests at priority 3)

## 0.36.0 — The Design Audit (2026-03-25)

Read-only UI validation that combines the structured 6-phase design checklist with the `frontend-design` skill's aesthetic eye. Points out what needs fixing without touching the code.

- Added `actions/ui-review.md` — validates UI quality across structure/IA, visual aesthetics, component consistency, UX copy, interaction/accessibility, and implementation patterns
- Produces a severity-rated findings report with file:line references and concrete fix suggestions
- Leverages both `rules-ui-design.md` (6-phase design workflow) and `frontend-design` skill (if installed) for comprehensive coverage
- Optional follow-up: capture high/medium findings as `domain: ui-design` REQs in the queue
- SKILL.md routing: `ui-review`, `review ui`, `design review`, `validate ui`, `ui audit` keywords (priority 5.5)
- Updated help menu, dispatch table, and next steps suggestions

## 0.35.0 — The Code Lens (2026-03-25)

Standalone codebase review, scoped by prime files and/or directories. Review consistency, patterns, security, and architecture without needing the REQ/UR queue.

- Added `actions/code-review.md` — full codebase review action with prime file and directory scoping
- SKILL.md routing: `code-review`, `audit codebase`, `review codebase` keywords (priority 5, before review work)
- Supports combined scoping: `do-work code-review prime-auth src/utils/` reviews the union of both
- Interactive mode when no scope given — lists available prime files and asks
- Optionally creates follow-up REQs for Critical/Important findings
- Updated help menu, README, and CLAUDE.md with code-review documentation

## 0.34.1 — The Consistency Pass (2026-03-24)

Documentation alignment across the entire skill.

- CLAUDE.md project structure now lists all action files (added quick-wins.md, install-ui-design.md, README.md)
- README.md now documents all actions (added Quick-Wins, Commit, Install UI Design, Version/Recap sections)
- SKILL.md action list now includes version and recap entries
- install-ui-design.md uses real installation mechanism (mkdir + curl) instead of nonexistent `claude skill add` command
- Fixed changelog dates for v0.32.0 and v0.33.0 (2026-03-23 → 2026-03-24)

## 0.34.0 — The Design Pipeline (2026-03-24)

Design-only deliverables now flow through the full pipeline without inventing code changes. Domain-specific review criteria are actually applied during review.

- Design-artifact exception in work.md — wireframe specs, IA docs, and visual specs placed outside `do-work/` (e.g., `docs/design/`) satisfy Implementation Summary and commit validation
- Work.md Step 7 now passes `rules-[domain].md` to the review agent (matching how Plan and Implementation agents already receive it)
- review-work.md gains a generic Domain-Specific Review hook — any domain can define review criteria in its rules file
- rules-ui-design.md gains a Design Artifacts section explaining where to place non-code deliverables

## 0.33.0 — The Design Install (2026-03-24)

One command to get production-grade UI design capabilities. `do-work install-ui-design` installs Anthropic's `frontend-design` skill into the current project.

- Added `actions/install-ui-design.md` — installs `frontend-design` Claude skill with automatic fallback to manual curl
- Added `install-ui-design` routing, verbs, help menu entry, and dispatch table in SKILL.md
- Works alongside `domain: ui-design` rules for a complete design workflow

## 0.32.0 — The Design Eye (2026-03-24)

UI design gets first-class domain support. Requests tagged `domain: ui-design` now load dedicated design rules covering IA, wireframing, visual aesthetics, component systems, UX copy, interaction specs, and heuristic reviews.

- Added `agent-rules/rules-ui-design.md` — phased design workflow (IA → wireframe → visuals → components → copy → interaction) with quality checks and implementation patterns
- Added `ui-design` as a domain option in capture and work action file schemas
- Includes accessibility baseline, heuristic review criteria, and handoff note guidance

## 0.31.2 — The Review Fix (2026-03-22)

Addresses PR review feedback on the Implementation Summary feature.

- Clarified Step 8 Discovered Tasks wording — explicitly notes it's a separate section from `## Implementation Summary`, not nested inside it
- Broadened "source files" to "project files" in Implementation Summary rules — config, CI, docs, and Dockerfiles now included in the manifest
- Scoped Step 9 validation check to successful REQs only — failed REQs may have no summary or staged project files, and that's expected

## 0.31.1 — The Alignment Check (2026-03-22)

Follow-up fixes to ensure Implementation Summary is consistently referenced across all action files.

- Added Implementation Summary to architecture diagram in work.md
- Added `## Implementation Summary` to crash recovery strip list (Step 1)
- Clarified Step 6 agent instructions: agent reports file list, orchestrator writes the formal summary
- Fixed present-work.md: removed stale "key files" reference from Lessons Learned extraction, pointed Key Files section to Implementation Summary
- Added `Summary...` line to progress reporting example

## 0.31.0 — The Archive Proof (2026-03-22)

Every completed REQ now carries a mandatory implementation manifest — no more guessing whether a REQ was actually built or just filed away.

- Added Step 6.25 (Implementation Summary) to work.md — mandatory file-change manifest for all routes (A, B, C)
- Removed `Key files:` from Lessons Learned (Step 7.5) — now covered by Implementation Summary
- Added commit validation check to Step 9 — flags mismatches between Implementation Summary and staged files
- Updated archived REQ example to reflect new format
- Added two new common-mistakes entries for Implementation Summary gaps

## 0.30.6 — The Alias Restore (2026-03-22)

Restored version/release aliases and anchored changelog lookup to skill root — fixes from PR review feedback.

- Restored `what's new`, `release notes`, `updates`, `history`, `what's changed` as version keywords (were dropped in 0.30.4)
- Anchored CHANGELOG.md lookup to skill root directory (same level as SKILL.md) so it doesn't pick up the project's own changelog
- Added route examples for all restored aliases

## 0.30.5 — The Missing Piece (2026-03-22)

Filled in gaps from the recap/version split — missing route examples, dispatch table entry, next steps, and help menu clarity.

- Added "Routes to Version" example section (was the only action without one)
- Added recap row to Action Dispatch table with `mode: recap` context
- Added recap to foreground subagent list
- Split help menu: `do-work version` and `do-work update` now on separate lines
- Added "After version / recap" to Suggest Next Steps
- Updated version.md header to mention recap handling

## 0.30.4 — The Recap (2026-03-22)

Split changelog into two focused commands: `do-work version` now shows last 5 skill releases alongside the version number, and `do-work recap` shows last 5 completed URs with their REQs. No more wall-of-text changelog dumps.

- `do-work version` now includes last 5 releases from CHANGELOG.md (read first ~80 lines only)
- Added `do-work recap` command — shows last 5 archived URs with REQ titles
- Removed `do-work changelog` / `changelog all` routes (version and recap cover both use cases)
- Updated help menu, routing table, and argument-hint

## 0.30.2 — The Scoped Check (2026-03-22)

Dirty-tree check now only looks at shipped skill files, so captured requests in `do-work/` no longer trigger false "dirty" warnings during updates. Also replaced the dangerous "delete directory contents" clean-update advice with a safe list of skill-owned paths. Prime-file lesson links now compute proper relative paths.

- Scoped `git status --porcelain` to specific skill paths (`SKILL.md`, `actions/`, `agent-rules/`, etc.), explicitly excluding `do-work/`
- Clean-update guidance in version.md and README.md now lists only safe-to-delete skill paths instead of suggesting a blanket directory wipe
- Lesson links in prime files now use paths relative to the prime file's location, so links work regardless of nesting depth
- Added stale-file warning to version.md update flow and README.md install section
- Added dirty-tree guard to the update flow with scoped file checking

## 0.30.1 — The Guard Rail (2026-03-22)

Fixes three review findings in quick-wins: ambiguous verb routing, glob-vs-directory mismatch, and hardcoded language list.

- Ambiguous verbs (`scan`, `opportunities`, `what can we improve`) now route to quick-wins only when used alone or with a directory path — followed by descriptive content they fall through to capture requests, matching the existing `check`/`review` pattern
- Dropped "or glob pattern" from the `$TARGET` interface — quick-wins accepts a directory path, not a file-level glob
- Replaced hardcoded extension list with language-adaptive detection — scans whatever languages are actually in the repo instead of only JS/TS/Python/PHP

## 0.30.0 — The Low Hanging Fruit (2026-03-22)

New quick-wins action scans a target directory for refactoring opportunities and low-hanging tests without modifying any files.

- Added `actions/quick-wins.md` — read-only codebase scanner that identifies long functions, copy-paste, god files, dead code, deep nesting, mixed concerns, and untested pure functions
- Outputs a structured markdown report ranked by effort vs impact (trivial+high first)
- Full SKILL.md integration: routing, dispatch, help menu, verb section, examples, and next-steps suggestions

## 0.29.5 — The Clean Slate (2026-03-22)

Crash recovery now strips stale phase sections so interrupted runs don't poison the next attempt. Pass 3a relocation spec tightened with explicit path preservation and conflict handling rules.

- Crash recovery clears `## Triage`, `## Exploration`, `## Plan`, and `## Testing` from interrupted REQs — prevents stale partial data from being trusted on retry
- Pass 3a now specifies exactly how each subtree type is relocated, what counts as a conflict vs. a duplicate, and when to leave misplaced copies for manual resolution

## 0.29.4 — The Pragmatist (2026-03-22)

Addressed three code review findings: cleanup detection gaps, internal contradictions, and overly rigid test requirements.

- Pass 3a now detects misplaced `do-work/` directories directly instead of relying on narrow file patterns that miss partial trees
- Resolved contradiction between Pass 3a relocation and "does not touch" exclusions — misplaced trees get fully relocated (error recovery), canonical root stays untouched (work action's domain)
- Red-green testing is now the default for bug fixes and new features; refactors, config, docs, and cleanup use pragmatic evidence instead

## 0.29.3 — The Self Updater (2026-03-22)

Version update now actually runs the update instead of telling you to do it. You're the user, not the agent's assistant.

- Changed "update available" flow to execute the curl command directly instead of printing it
- Agent verifies the update succeeded by re-reading the version file

## 0.29.2 — The Drift Catcher (2026-03-22)

Cleanup now detects `do-work/` directories created in the wrong location when an agent's CWD drifts into a subdirectory. Pass 3a scans the repo for misplaced `do-work/` trees and relocates their contents to the canonical root queue.

- Added Pass 3a: scan for misplaced `do-work/` directories outside project root
- Renamed existing Pass 3 to Pass 3b (misplaced folders within the archive)
- Conflict-safe: duplicates are reported for manual resolution, not auto-deleted
- Updated commit staging to include misplaced directory paths

## 0.29.1 — The Gap Plug (2026-03-22)

Fixed a gap in Step 6.5 where a partial prime test map caused zero tests to run for unmapped files. The fallback is now per-file: matched files use the prime's commands, unmatched files fall back to generic detection.

- Step 6.5: fallback to generic detection applies per changed file, not per prime section
- Explicit: a partial prime map is not an excuse to skip tests for unmapped files

## 0.29.0 — The Red-Green Trace (2026-03-21)

Each request now proves itself with red-green test validation and cross-REQ traceability. Tests must fail before implementation and pass after — proving the request is delivered, not just that tests exist. When a request intentionally changes behavior tested by a prior request, the builder documents which REQ's tests changed and why.

- Builder instructions: write/identify tests before code, confirm they fail, then implement
- Builder instructions: when existing tests break, document cross-REQ impact with originating REQ reference
- Step 6.5 testing template: added red-green validation and cross-REQ impact sections
- review-work.md Step 6: reviewers check for red-green evidence and cross-REQ test traceability
- review-work.md Step 7: acceptance testing verifies cross-REQ test updates are intentional and documented

## 0.28.2 — The Test Map (2026-03-21)

Agents now check prime files for project-specific test commands before falling back to generic detection. If your prime maps code areas to test commands, builders and reviewers will follow that mapping instead of just running `npm test`.

- work.md builder instructions: added bullet to check prime file testing sections
- work.md Step 6.5: prime test guidance comes first, generic detection is the fallback
- review-work.md Step 6: Test Adequacy now checks whether the *right* tests were run per the prime
- review-work.md Step 7: Acceptance testing checks prime for test command mappings

## 0.28.1 — The Light Install (2026-03-18)

Update command no longer pulls in the `skills` npm package. Now it's a single curl+tar one-liner that downloads files directly from GitHub — no npm, no intermediary tools.

- Replaced `npx skills add` with `curl | tar` in update commands and install docs
- `_dev/` folder excluded automatically during extraction
- Added directory guidance ("run from the skill's root directory") to prevent extracting into the wrong location

## 0.28.0 — The Feedback Loop (2026-03-18)

Lessons learned now flow back into prime files. When a REQ captures lessons, the relevant prime files get a link under a `## Lessons` section — so future agents working on that area of the codebase benefit from past experience without re-reading archived REQs.

- Added prime file update step to work.md Step 7.5 (pipeline mode)
- Added prime file update step to review-work.md Step 9.5 (standalone mode)
- Links are scoped: only lessons relevant to a prime file's domain get added

## 0.27.9 — The Right File (2026-03-18)

Fixed Step 9.5 targeting the archived REQ in pipeline mode — the file hasn't been archived yet at that point. Lesson capture is now standalone-only (work.md Step 7.5 handles it in pipeline mode), self-validation still runs in both modes.

- Changed "archived REQ file" → "the REQ file" in Step 9.5
- Made lesson capture standalone-only to avoid duplication with work.md Step 7.5
- Reordered steps: self-validation first, then lesson capture

## 0.27.8 — The Self-Check (2026-03-18)

Replaced human validation gate in review-work with automated self-validation. The review now re-examines its own findings, captures lessons learned, and creates follow-up REQs for anything it missed — no human prompt blocking the flow.

- Removed Step 9.5 human validation prompt (standalone mode)
- Added self-validation pass that runs in both pipeline and standalone modes
- Lessons learned are now captured automatically by the review itself

## 0.27.7 — The Trim (2026-03-14)

Rewrote `CLAUDE.md` as a proper prime file — project structure map, concise commit rules, and agent compatibility guidance. Cut the noise, kept the signal.

- Added project structure overview with file-level descriptions
- Condensed changelog formatting rules from 10 bullets to a template + one-liner
- Tightened agent compatibility section from 5 verbose bullets to 3 clear ones

## 0.27.6 — The Unboxing (2026-03-14)

Moved `agent-rules/` out of the `do-work/` subdirectory to the repo root. When the skill is installed into a project that already uses `do-work/` as its working directory, the old layout would create a nested `do-work/do-work/` path. Now the rules live at `agent-rules/` — no nesting, no confusion.

- Moved `do-work/agent-rules/` → `agent-rules/` at repo root
- Updated path references in `work.md`

## 0.27.5 — The Spring Clean (2026-03-14)

Trimmed the `_dev/` folder and fixed a fragile symlink. Root `CLAUDE.md` was a symlink to `_dev/CLAUDE.md` — replaced it with a real file so it can't break when `_dev/` gets cleaned up. Removed the stub agent config files too.

- Replaced `CLAUDE.md` symlink with a standalone file (was `CLAUDE.md -> _dev/CLAUDE.md`)
- Deleted `_dev/CLAUDE.md` (now lives at root as a real file)
- Deleted `_dev/AGENTS.md` and `_dev/GEMINI.md` (one-line stubs with no real value)

## 0.27.4 — The Stage Call (2026-03-13)

The video preview actually works now. The Remotion project was missing `registerRoot()`, so `npm run preview` would launch Studio with nothing to show. Added a proper entry file and pointed the preview script at it.

- Added `src/index.ts` with `registerRoot(RemotionRoot)` call
- Updated `package.json` preview script to point at `src/index.ts` instead of `src/Root.tsx`

## 0.27.3 — The Right Lane (2026-03-13)

Discovered-task approvals now route to `pending` instead of `completed`. Previously, confirming "Yes, add to queue" on a discovered task hit the "Builder Was Right" fast-path, which archived it immediately — the task never actually ran.

- Added "Approved Discovered Task" section to clarify workflow with correct `pending` routing
- Updated "Confirm builder's choice" logic to distinguish discovered tasks from builder-decision follow-ups
- Discovered tasks confirmed for processing stay in `do-work/` and enter the normal work queue

## 0.27.2 — The Safety Catch (2026-03-13)

Restored the missing-domain fallback guard for loading `rules-[domain].md` in the work pipeline. Steps 4 and 6 now gracefully skip the rules file when `domain` is absent from frontmatter or the file doesn't exist, instead of assuming it's always resolvable.

- Restored `(if domain is missing or the file doesn't exist, skip loading it)` guard to Step 4 (Planning, Route C)
- Added the same guard to Step 6 (Implementation) for consistency

## 0.27.1 — The Field Guide (2026-03-12)

General agent rules now include the Prime Files Philosophy. Agents know what prime files are, how to write them, and what to avoid — before they ever encounter one in a REQ.

- Added PRIME Files Philosophy section to `rules-general.md`
- Covers purpose, conciseness, pointer-not-copy pattern, volatile metric avoidance, and multi-aspect support

## 0.27.0 — The Trailblazer (2026-03-12)

The work orchestrator now speaks prime files. Plan and implementation agents receive prime files as first-class context, and the builder is instructed to create missing ones on the fly. The archived REQ example also carries the new field for reference.

- Added `prime_files: []` to Request File Schema YAML example
- Updated Step 4 (Planning): Plan agent now receives prime files and uses them as the strict index
- Updated Step 6 (Implementation): general-purpose agent receives prime files alongside domain rules
- Added Prime Files bullet to agent instructions — read first, create if missing, keep low-noise
- Added `prime_files: []` to the Archived Request File Example frontmatter

## 0.26.0 — The Wayfinder (2026-03-12)

Capture now knows about prime files — semantic index files that point agents to the right source code. REQs carry a `prime_files` array in frontmatter, the PLAN phase reads them alongside agent rules, and Step 1 routes to matching prime files automatically.

- Added `prime_files: []` field to Simple REQ YAML frontmatter
- Updated PLAN checkbox to read listed `prime_files` and agent rules
- Added prime file routing bullet to Step 1: Parse and Assess
- Updated Step 5 item 2 to populate `prime_files` with discovered paths

## 0.25.1 — The Billboard (2026-03-12)

README and SKILL.md now advertise the two new features. Users browsing the docs will see Human UAT under Review Work and Interactive Explainer under Present Work without digging into the action files.

- README: Added Human UAT bullet to Review Work section
- README: Added Interactive Explainer bullet to Present Work section
- SKILL.md: Updated help menu description for `do-work present work`

## 0.25.0 — The Show Floor (2026-03-12)

Present work now generates an interactive HTML explainer alongside the client brief and video. It's a single `.html` file — no build steps, no npm — that stakeholders can double-click to open in any browser.

- Added section 4c: Interactive Explainer (Single-File HTML) to `present-work.md`
- Zero dependencies: HTML5 + Tailwind CDN + Vanilla JS in one file
- Includes Before/After toggle, step-by-step architecture walkthrough, and value summary
- Updated Step 5 summary to list the HTML file with double-click-to-open instructions
- Renumbered Portfolio artifacts from 4c to 4d

## 0.24.0 — The Human Touch (2026-03-12)

Reviews in standalone mode now pause for human validation before closing out. The reviewer presents its report, then asks the user to test manually and share feedback. Lessons learned go straight into the archived REQ; bugs become follow-up REQs automatically.

- Added Step 9.5: Human Validation (Standalone Mode Only) to `review-work.md`
- Lessons learned / architectural feedback appended to the archived REQ's `## Lessons Learned` section
- Bugs and fix requests treated as Important findings, routed to Step 10 for follow-up REQ generation
- Pipeline mode skips the step entirely — no blocking the automated loop

## 0.23.7 — The Softer Touch (2026-03-12)

Toned down the APPLY and Out-of-Scope agent instructions. Same constraints, less adversarial language — agents follow guidance better when it reads like coaching, not a legal contract.

- Rewrote APPLY phase: "stay focused" instead of "you are forbidden"
- Rewrote Out-of-Scope: "do not fix them inline" instead of "DO NOT fix them. You must strictly adhere to..."

## 0.23.6 — The Reference Card (2026-03-12)

The archived REQ example now shows what a completed P-A-U loop looks like. Agents have a concrete reference for how the execution state checkboxes should read when a request is done.

- Added completed `## AI Execution State (P-A-U Loop)` section to the archived request file example

## 0.23.5 — The Fine Tuning (2026-03-12)

Two small fixes in work.md: the domain field in the REQ schema no longer looks like a pipe-delimited value (it's a single choice), and the APPLY phase now explicitly permits editing the REQ file to update state checkboxes.

- Fixed `domain` field in Request File Schema — shows example value with comment instead of ambiguous pipe syntax
- Added REQ-file exception to APPLY phase scope restriction — agents can update their own state checkboxes

## 0.23.4 — The Crash Guard (2026-03-12)

Appending steps in the work loop are now idempotent. If a crash or re-entry happens mid-REQ, Steps 3, 4, and 5 skip sections that already exist instead of writing duplicates.

- Step 3 (Triage): guards `## Triage` append with existence check
- Step 4 (Planning): guards `## Plan` append and skip-note with existence checks
- Step 5 (Exploration): guards `## Exploration` append with existence check

## 0.23.3 — The Tidy Tenant (2026-03-12)

Agent rules now live inside `do-work/` where they belong. No more root-level pollution — the skill's runtime directory holds everything it creates.

- Moved `agent-rules/` to `do-work/agent-rules/`
- Updated all references in `capture.md` and `work.md` to use the new path

## 0.23.2 — The Atomic Ledger (2026-03-09)

Uncommitted files no longer pile up without a home. The new commit action analyzes your working tree, traces files back to archived REQs when possible, semantically groups the rest, and commits everything in small atomic batches — each one traceable.

- **commit.md**: New action — analyzes uncommitted files, associates with archived REQs for traceability, groups semantically into atomic commits (1-5 files each), and reports a summary
- **SKILL.md**: Added routing (priority 8), commit verbs, action list, dispatch table, help menu, next steps, and examples for the commit action

## 0.23.1 — The Paper Trail (2026-03-09)

Every action now commits its own work. Capture, cleanup, review-work, and work all have explicit git commit steps so changes are never left unstaged. The work action also writes the real commit hash back into the archived REQ for traceability.

- **capture.md**: Added Step 7 — commits the UR folder and new REQ files after capture, with addendum-aware message format
- **cleanup.md**: Added Commit section — commits structural moves (archive consolidation, legacy, misplaced folders) after all three cleanup passes
- **review-work.md**: Added Commit section for standalone mode — commits the appended Review section and any follow-up REQs (pipeline mode defers to work Step 9)
- **work.md**: Step 1 now uses explicit glob pattern `do-work/REQ-*.md` with a fallback verification to prevent false "queue empty" results
- **work.md**: Step 9 now writes the real commit hash back to the archived REQ's `commit:` frontmatter field via `--amend`, giving review-work and present-work reliable traceability

## 0.23.0 — The Director's Cut (2026-03-07)

Present work now generates real Remotion video projects instead of markdown video scripts. The video deliverable is a full React/TypeScript project with animated scenes you can preview in the browser via `npx remotion studio` — no mp4 rendering needed.

- **present-work.md**: Replaced section 4b markdown video script template with Remotion project structure (Root, Video, scene components, styles)
- **present-work.md**: Added scene content guidelines, animation patterns, and project scaffolding instructions
- **_dev/deliverables**: Replaced `do-work-video-script.md` with a complete `do-work-video/` Remotion project as the reference example
- **SKILL.md, README.md**: Updated video deliverable descriptions to reflect Remotion video format

## 0.22.7 — The Missing Step (2026-03-07)

Commits weren't happening during `do-work run` because the architecture diagram — the visual flow agents follow — never mentioned Step 9 (Commit). The detailed instructions existed but agents never reached them. Now the diagram shows Commit as an explicit step after Archive.

- **work.md**: Added "Commit (git repos only)" node to the architecture diagram between Archive and Loop
- **work.md**: Added bold reminder callout below the diagram reinforcing that every completed request gets a commit before looping

## 0.22.6 — The Guardrail (2026-03-04)

Six cross-file fixes addressing safety gaps, missing guardrails, and inconsistent instructions.

- **work.md**: Expanded commit failure guidance — explicit prohibition of `--no-verify` and `--no-gpg-sign`, instructions to investigate and fix hook errors instead of bypassing them
- **work.md**: Restored Orchestrator Checklist (per-request step verification) and Common Mistakes to Avoid section — prevents file management errors, premature archiving, and unsafe git operations
- **work.md**: Clarified follow-up creation filter for `- [~]` items — create follow-ups for UX/scope/data-representation decisions, skip purely technical decisions (caching, algorithms, internal naming)
- **work.md**: Clarified Lessons Learned scope — required for Routes B/C, optional for Route A (consistent with present-work.md)
- **review-work.md**: Clarified Test Adequacy N/A handling — explicitly excluded from overall score average (not counted as 0%)

## 0.22.5 — The Regression Fix (2026-03-04)

Three regressions restored from content that was lost during prior simplification passes.

- **SKILL.md**: Restored "Human time has two optimal windows" section — explains the two-phase interaction model (capture phase for real-time clarification, batch review for accumulated questions) that underpins the entire system
- **capture.md**: Restored full "Step 3: Capture-Phase Clarification" — was reduced to a single paragraph ("Clarify Only If Needed"), losing the AskUserQuestion guidance, good/bad examples, what NOT to ask about, and after-capture open questions flow
- **work.md**: Restored safe git staging — `git add -A` replaced with specific file staging, plus safety instructions warning against `git add -A` / `git add .` (risk of staging secrets, `.env` files, or unrelated changes)

## 0.22.4 — The Course Correct (2026-03-04)

Two fixes from PR review. CHANGELOG.md moved back to root so `do-work changelog` works for installed users (it was accidentally excluded with the `_dev/` move). Standalone reviews now find UR input regardless of whether the UR has been archived yet.

- Moved `CHANGELOG.md` back to root — the changelog command needs it at install time
- Fixed `review-work.md` Step 3: UR input lookup now checks `user-requests/` first, falls back to `archive/` — works in both pipeline and standalone modes regardless of UR completion state

## 0.22.3 — The Lighter Carry (2026-03-04)

Sample deliverables no longer ship with the skill. The three do-work-specific outputs (client brief, video script, portfolio summary) moved to `_dev/deliverables/` so they're excluded from installation. Users generate their own via `do-work present`.

- Moved 3 sample deliverable files from `do-work/deliverables/` to `_dev/deliverables/`

## 0.22.2 — The Tidy Install (2026-03-04)

Dev-only files no longer tag along when someone installs the skill. CLAUDE.md, CHANGELOG.md, AGENTS.md, and GEMINI.md now live in `_dev/` — excluded by the skills CLI's underscore convention. A root symlink keeps CLAUDE.md discoverable for repo development.

- Moved 4 dev-only files to `_dev/` directory (excluded during `npx skills add` installation)
- Added `CLAUDE.md` symlink at repo root for Claude Code auto-discovery
- Updated CLAUDE.md changelog path reference to `_dev/CHANGELOG.md`

## 0.22.1 — The Badge (2026-03-04)

Verify now clearly identifies itself as capture QA, so agents (and users) don't confuse it with implementation review.

- Added "capture QA" clarification to `verify-requests.md` opening description
- Fixed stale "Critical" reference in verify's "What NOT To Do" section (should be "Important" per severity alignment)

## 0.22.0 — The Alignment (2026-03-04)

Cross-file severity levels and extraction lists are now consistent. Agents following one action file won't contradict another.

- Aligned `verify-requests.md` severity levels to match `review-work.md` — replaced Critical/Important/Minor with Important/Minor/Nit (Ambiguous stays as-is since it's verify-specific)
- Added Builder Guidance to `review-work.md` Step 2 extraction list — reviewers now calibrate expectations based on certainty level (Firm vs Exploratory)
- Marked Lessons Learned as optional in `present-work.md` Step 2 — Route A REQs skip this section per `work.md`

## 0.21.1 — The Addendum Fix (2026-03-04)

Addendum REQs now work reliably in non-git environments and the builder knows what to do with them.

- Made commit hash conditional ("if available") in `capture.md`'s Prior Implementation section — non-git projects legitimately have no hash
- Added addendum_to handling to `work.md` Step 3 (Triage) — builder now reads the original REQ for context, closing the timing gap where capture skips Prior Implementation for in-flight originals that complete before the addendum is built

## 0.21.0 — The Level Set (2026-03-04)

Ten cross-file inconsistencies and instruction gaps cleaned up. Agents following these docs literally should now get consistent behavior across all action files.

- Fixed `cleanup.md` self-contradiction — Pass 1 no longer references `do-work/working/` (work action handles its own files before cleanup runs)
- Fixed stale "Critical" severity in `README.md` — the defined levels are Important/Minor/Nit
- Aligned follow-up REQ creation rules — `work.md` now matches `review-work.md`: follow-ups are per-Important-finding, not score-gated
- Added overall review score formula to `review-work.md` — average of percentage dimensions with Risk/Acceptance modifiers
- Closed 200–500 word classification gap in `capture.md` — removed the >500 word floor from Complex (features/constraints matter more than word count)
- Fixed `review-work.md` Step 2: "Plan (if Route B/C)" → "Plan (if Route C)" — planning is Route C only
- Fixed `capture.md` leading-slash references (`/do-work verify requests` → `do-work verify requests`)
- Fixed `version.md` agent compatibility — replaced tool-specific "WebFetch" with generalized language
- Documented `hold/` directory in `cleanup.md` archive structure
- Added review annotation exception to `capture.md` immutability rule (cross-ref to `review-work.md`)

## 0.20.7 — The Chunker (2026-03-04)

Clarify workflow now chunks questions by count (max 4 per prompt), not by REQ. A single REQ with 6 questions gets 2 prompts instead of blowing the limit.

- Fixed question batching in clarify mode to respect per-prompt limits

## 0.20.6 — The Context Bridge (2026-03-04)

Addendum REQs for archived work no longer leave the builder guessing. When creating a follow-up to a completed request, capture now reads the original archived REQ and includes a `## Prior Implementation` section — key files, patterns used, commit hash — so the builder has full context without re-discovering what already exists.

- Added `## Prior Implementation` section to the addendum REQ template in `capture.md`
- Added "Context is critical" guidance — instructs capture to read the original archived REQ before writing the addendum
- Updated "Addendum to Archived Request" example to show the prior-implementation flow

## 0.20.5 — The Tight Seal (2026-03-03)

Addendum rules in `capture.md` are now airtight. When an original request is archived, creating an addendum always produces a new UR + REQ in `do-work/` root — so the work loop can pick it up. The archive stays immutable.

- Added explicit "New REQ lands in" column to the duplicate-handling table
- Strengthened the Immutability Rule to state that new addendum REQs always go to `do-work/` root
- Clarified that archived URs are immutable — addendums always get a fresh UR
- Added "Addendum to Archived Request" example to make the pattern unambiguous

## 0.20.4 — The Right Folders (2026-03-03)

Two instruction bugs that would cause literal agent implementations to fail. Standalone review mode now searches UR subfolders for recent work, and Step 1 of the work loop now explicitly reads frontmatter before selecting the next request.

- `review-work.md` Step 1: "no target specified" now searches `do-work/archive/UR-NNN/` subdirectories in addition to the archive root — completed REQs live in UR folders after cleanup, not the root
- `work.md` Step 1: replaced "List (don't read) ... pick first with `status: pending`" with an explicit frontmatter-read step — status is in YAML frontmatter, not the filename, so listing alone can't filter by status

## 0.20.3 — The Bug Hunt (2026-03-03)

Three pre-existing bugs squashed. Nothing new, just things that were quietly wrong.

- `verify-requests.md`: Removed dangling "per Step 3.5" reference — that step doesn't exist
- `review-work.md`: Removed phantom "Critical" severity from Step 10 — the defined levels are Important/Minor/Nit, not Critical
- `SKILL.md`: Added `audit code` to routing table row 5 so the table matches the verb list updated in 0.20.2

## 0.20.2 — The Fine Print (2026-03-02)

Three small clarity improvements borrowed from sibling branches. Nothing dramatic — just sharper routing, a missing severity level, and a better signpost for confused users.

- `verify-requests.md`: Added a redirect note under "When to Use" — if you want code review, use `review work` instead
- `review-work.md`: Added **Nit** as a fourth finding severity (below Minor; carries zero score weight — stylistic suggestions only)
- `SKILL.md`: Disambiguated `audit` routing — `audit` alone stays in verify, `audit code` and `audit implementation` now correctly route to review work

## 0.20.1 — The Self-Portrait (2026-03-02)

The do-work skill presented itself. Generated the first set of client-facing deliverables for the skill as a product — a client brief with full architecture diagrams and data flow, a 3-minute video script (7 scenes, capture through portfolio), and a portfolio summary covering all 20 releases.

- Generated `do-work/deliverables/do-work-client-brief.md` — architecture, data flow, value proposition, competitive advantage, and roadmap
- Generated `do-work/deliverables/do-work-video-script.md` — 7-scene walkthrough from problem to install command
- Generated `do-work/deliverables/do-work-portfolio-summary.md` — all 20 versions catalogued with cumulative value prop and cross-project lessons

## 0.20.0 — The Pitch Deck (2026-03-02)

Completed work can now speak for itself. The new "present work" action reads your archive — requests, implementation history, code diffs, and lessons learned — and generates client-facing deliverables: briefs that explain what was built and how it works, value propositions that sell the impact, and video scripts for demo walkthroughs. Run it on a single UR or across the full portfolio. Also added a "Lessons Learned" section to archived REQs so institutional knowledge survives between sessions, and refined diff hygiene to protect those lessons from cleanup.

- New `present work` action (`actions/present-work.md`) — two modes:
  - **Detail mode** — deep dive on a specific UR or REQ: client brief with architecture diagrams, data flow, value proposition, and optional video script (Remotion/Loom-ready)
  - **Portfolio mode** (`do-work present all`) — summary of all completed work with cumulative value proposition and cross-project lessons learned
- Artifacts saved to `do-work/deliverables/` for reuse and sharing
- New `## Lessons Learned` section in work.md — archived REQs now capture what worked, what didn't, key files, and gotchas (Step 7.5, between Review and Archive)
- Refined diff hygiene in review-work.md — explicitly protects comments that document reasoning, failed approaches, or architectural decisions
- SKILL.md: new routing (priority 6), dispatch table, help menu, verb section, examples, and next-step suggestions for present work
- README.md: new Present Work section

## 0.19.1 — The Neighbor Check (2026-03-02)

Review work now checks whether your change broke something nearby. Regression risk analysis reads the diff to identify callers and dependents, acceptance testing exercises adjacent features, and suggested testing flags regression scenarios. Also catches leftover debug artifacts and commented-out experiments.

- Added regression risk to Risk Assessment — identifies callers/dependents of changed code, flags changed interfaces, notes shared utilities
- Replaced "Check integration" with broader "Check for regressions" in acceptance testing — run adjacent tests, exercise other consumers of shared code, verify bug fixes don't break related behaviors
- Added "Regression scenarios" category to Suggest Additional Testing
- Added diff hygiene to Code Quality — catches debug artifacts, console.log/print statements, commented-out experiments, temp files

## 0.19.0 — The Wide Lens (2026-03-02)

Review and verify got proper names and bigger jobs. "Verify" becomes "verify requests" — it checks capture quality. "Review" becomes "review work" — and now it does requirements checking (did we build what was asked?), code review, acceptance testing (actually run the thing), and suggests additional testing the user should do. Every action now ends with suggested next prompts so you always know what to do next.

- Renamed `verify.md` → `verify-requests.md`; action name is now "verify requests" across routing, dispatch, help menu, and examples
- Renamed `review.md` → `review-work.md`; action name is now "review work" — enhanced with three new phases:
  - **Requirements check** — walks through every REQ requirement line-by-line to confirm it was delivered
  - **Acceptance testing** — runs the app/tests and verifies the feature works end-to-end, not just in the diff
  - **Suggested additional testing** — recommends manual verification, integration, edge cases, and environment-specific checks
- Review report now includes a requirements checklist, acceptance result (Pass/Partial/Fail/Untested), and suggested testing section
- Added "Suggest Next Steps" section to SKILL.md — every action now ends with 2-3 fully qualified prompt suggestions (`do-work verify requests`, not just `verify`)
- Updated capture.md, work.md, README.md with new action names and references

## 0.18.0 — The Plain English (2026-03-02)

Actions now say what they mean. "Capture" becomes "capture requests," the confusing "answers mode" becomes "clarify questions" with `do-work clarify`, and bare `do-work` shows a help menu instead of jumping straight to the work loop.

- Renamed action: "capture" → "capture requests" across SKILL.md, README, dispatch table, capture.md, work.md
- Renamed "answers mode" → "clarify questions" — new primary verb is `do-work clarify` (old verbs still work)
- Bare invocation (`do-work` with no arguments) now shows a help menu with sample prompts instead of asking "Start the work loop?"
- Added `clarify questions` row to action dispatch table (routes to work.md with `mode: clarify`)
- README now documents "Clarify Questions" as a standalone section alongside the other actions

## 0.17.0 — The Name Tag (2026-03-02)

The "do action" is now the "capture action." No more confusion between the skill name (`do-work`) and the action that captures requests. `do.md` becomes `capture.md`, all references updated across the codebase. Also fixes three workflow consistency issues found during a full trace.

- Renamed `actions/do.md` → `actions/capture.md` and updated all references in SKILL.md, work.md, README.md, CLAUDE.md
- SKILL.md routing: `→ do` becomes `→ capture` everywhere — routing table, content signals, examples, dispatch table
- Architecture diagram in work.md now shows Open Questions and the pending-answers follow-up flow
- Step 10 (Loop or Exit) now runs cleanup even when only `pending-answers` REQs remain, then reports them
- Answers Mode step 5 now explicitly skips REQs already completed by the Builder Was Right path

## 0.16.0 — The Full Loop (2026-03-02)

The Open Questions system now has a complete lifecycle — from capture to drain. Five improvements tighten the feedback loop: verify resolves ambiguous questions on the spot (user is present, why not ask?), `do-work answers` gives users a dedicated command to batch-review accumulated questions, follow-up REQ creation moves to the archive step so timing is unambiguous, confirmed builder choices skip the work loop entirely, and verify's question handling is explicitly documented as different from review's.

- `verify.md`: Ambiguous gaps now get presented to the user immediately during verify — resolve on the spot, defer, or leave for the builder
- `SKILL.md` + `work.md`: New "answers mode" — `do-work answers`/`questions`/`pending` presents all `pending-answers` REQs for batch review
- `work.md`: "Builder Was Right" path — when user confirms builder's choice, follow-up archives directly with no work cycle
- `work.md`: Follow-up REQ creation moved from Step 3.5 to Step 8 (Archive) with full template — timing is now explicit
- `verify.md`: Clarified that verify never sets `pending-answers` status — it already asked the user; remaining questions stay on `pending` REQs

## 0.15.0 — The No-Block Build (2026-02-26)

Open Questions no longer block the build phase. The builder uses its best judgment, completes the REQ, and creates `pending-answers` follow-up REQs for decisions that need user validation. Human interaction is optimized for two windows: capture time (ask freely) and batch-review time (user returns to answer accumulated questions). Questions now include recommended defaults and alternatives so they're answerable at a glance.

- `do.md`: Open Questions now include `Recommended:` and `Also:` choices; capture time is the primary ask window — use the ask tool immediately instead of deferring
- `work.md`: Step 3.5 is no longer a blocking gate — builder marks `- [~]` with reasoning, completes the REQ, then queues `pending-answers` follow-ups for user review
- `work.md`: New `pending-answers` status — work loop skips these; user batch-reviews them between runs
- `work.md`: Step 1 now skips `pending-answers` REQs and reports them when queue is otherwise empty
- `verify.md`: Ambiguous gaps use the choice format (`Recommended:` / `Also:`) when adding Open Questions
- `review.md`: Ambiguous-requirement follow-ups use `status: pending-answers` with choice format

## 0.14.0 — The Clarification Gate (2026-02-26)

Ambiguous requirements now get caught before code gets written. Open Questions in REQs use a structured checkbox format, the work action pauses at a new Step 3.5 checkpoint to resolve them with the user, verify flags genuinely ambiguous gaps for clarification instead of just failing them, and review creates follow-up REQs with Open Questions when the root cause is unclear intent rather than a code bug.

- `do.md`: Open Questions now use `- [ ] question text` checkbox format with `(context: ...)` annotations
- `work.md`: New Step 3.5 — Resolve Open Questions checkpoint that pauses for user input before implementation
- `verify.md`: New "Ambiguous" gap classification that generates Open Questions on the REQ instead of just reporting a gap
- `review.md`: Follow-up REQs for ambiguous-requirement findings now include `## Open Questions` to trigger the clarification checkpoint

## 0.13.0 — The Second Opinion (2026-02-25)

Every completed request now gets a code review before it's archived. The work pipeline gained a new step between testing and archive that reads the actual diff, compares it against the original requirements and UR, scores the implementation across five dimensions, and creates follow-up REQs when it finds real issues. You can also invoke it manually on anything already shipped.

- New `review` action (`actions/review.md`) — post-work code review with requirements tracing
- Two modes: **pipeline** (auto-triggered in the work loop after tests pass) and **standalone** (manual via `do-work review`)
- Scores on Requirements Compliance, Code Quality, Test Adequacy, Scope Discipline, and Risk Assessment
- Creates follow-up REQs (using `addendum_to` pattern) for Critical/Important findings — they re-enter the queue automatically
- Review depth scales with route: quick scan for Route A, standard for B, thorough for C
- `work.md` updated: new Step 7 (Review), renumbered Archive→8, Commit→9, Loop→10
- `SKILL.md` routing updated: "review"/"review code"/"code review" → review action (priority 4); "review requests"/"review reqs" still → verify
- REQ living documents now include a `## Review` section with per-dimension scores and follow-up links

## 0.12.7 — The Cold Start (2026-02-25)

The do action now knows what to do the very first time it runs. Previously, agents following the instructions would try to scan `do-work/` for duplicates and numbering before the directory existed — a guaranteed stumble on first use. Now there's explicit guidance for bootstrapping the folder structure, starting numbering at 1, skipping duplicate checks on an empty project, and ensuring directories exist before writing files.

- Added "First-Run Bootstrap" subsection under Core Rules — create `do-work/` and `user-requests/`, don't pre-create `working/`/`archive/`
- Added fallback to File Naming: start at 1 when no existing files found
- Added fresh-bootstrap skip to Step 2 (Duplicate Check): no files means no duplicates to scan
- Added directory-creation guard to Step 5 (Write Files): ensure paths exist before writing

## 0.12.6 — The Missed Spot (2026-02-25)

Fixed the last bare `archive/UR-*/` path in the duplicate-check instructions. The v0.12.4 path fix caught the numbering section but missed the same issue in the Step 2 duplicate scan — agents following it literally would skip archived UR subfolders and let duplicates through.

- Fixed `do.md` line 153: `archive/UR-*/` → `do-work/archive/UR-*/`

## 0.12.5 — The Deep Check (2026-02-25)

Duplicate detection now actually reads queued request files instead of just glancing at filenames. A `REQ-042-ui-cleanup.md` whose `## What` says "fix spacing on the settings page" will now correctly match a new submission of "fix the spacing and layout on the settings page" — no more phantom duplicates slipping through because the slug didn't match the phrasing.

- Queued requests (`do-work/`): agent now inspects `title`, heading, and `## What` for semantic intent matching
- In-flight and archived requests (`working/`, `archive/`): still filename-scan only (fast, and files are immutable anyway)
- Decision table and addendum formats unchanged — this is a detection improvement, not a workflow change

## 0.12.4 — The Right Address (2026-02-25)

Fixed ambiguous paths in the REQ/UR numbering instructions. The do action told agents to scan `working/` and `archive/` for existing IDs — bare paths that miss the `do-work/` prefix every other reference uses. Agents following the instructions literally would scan nonexistent directories and risk creating duplicate request IDs.

- Fixed `do.md` line 47: `working/` → `do-work/working/`, `archive/` → `do-work/archive/`
- Explicitly listed UR scan locations (`do-work/user-requests/UR-*/` and `do-work/archive/UR-*/`)
- Added file pattern hints (`REQ-*.md`, `UR-*`) so agents know what to look for

## 0.12.3 — The Time Traveler (2026-02-25)

Fixed three changelog entries (0.12.0, 0.12.1, 0.12.2) that were dated 2025 instead of 2026. The release chronology is now consistent across all versions.

- Corrected year in 0.12.0, 0.12.1, and 0.12.2 headings from 2025-02-25 to 2026-02-25

## 0.12.2 — The New Address (2026-02-25)

Upstream references updated to the forked repository. README install command, SKILL.md upstream URL, and version action URLs all now point to `knews2019/skill-do-work` instead of the original `bladnman/do-work`.

- Updated README.md install command to `npx skills add knews2019/skill-do-work`
- Updated SKILL.md upstream URL to `knews2019/skill-do-work`
- Updated version.md upstream URL, install commands, and GitHub link to `knews2019/skill-do-work`
- CHANGELOG.md historical entries left unchanged (they reference the original repo accurately)

## 0.12.1 — The Passport Check (2026-02-25)

Removed a hardcoded `Co-Authored-By: Claude <noreply@anthropic.com>` trailer from the commit template in work.md. Agents on other platforms would stamp Claude-specific metadata onto their commits just by following the template verbatim — violating the agent compatibility rules. The trailer is now a documented option with a generic example, not a baked-in default.

- Removed tool-specific co-author line from the commit template example
- Added guidance: use your platform's co-author convention if it has one, otherwise omit

## 0.12.0 — The Diet (2026-02-25)

The skill shed two-thirds of its weight. `do.md` dropped from 883 to 288 lines, `work.md` from 1,277 to 383. Same behavior, dramatically less noise. Redundancy across files (folder structure repeated 4 times, schemas defined twice, checklists restating the workflow) was consolidated or cut. Agent prompt templates in work.md were merged into one. The 158-line retrospective section, 7 overlapping examples, and standalone "What NOT to do" sections — all trimmed to their essentials.

- `do.md`: 883 → 288 lines (67% reduction) — consolidated formats, trimmed examples from 7 to 4, folded checklists into workflow, cut platform-specific screenshot bloat
- `work.md`: 1,277 → 383 lines (70% reduction) — unified agent prompt template, cut duplicate retrospective section, merged error handling into a table, removed redundant orchestrator checklist
- All behavioral rules preserved — UR+REQ pairing, immutability, complexity triage, living logs, capture≠execute boundary
- Zero behavior changes — this is a documentation refactor, not a feature change

## 0.11.1 — The Soft Landing (2026-02-24)

Subagent dispatch no longer assumes subagents exist. Environments without Task subagents can now fall back to reading the action file directly in the current session — no more broken routing in simpler tools. The dispatch section is restructured as "if available / if not" so the skill stays portable.

- Added fallback path: read action file directly when subagents are unavailable
- Removed Claude Code-specific language (Task tool, `run_in_background`)
- Dispatch table simplified — background column moved into subagent-specific guidance

## 0.11.0 — The Delegate (2026-02-24)

Actions now run in subagents instead of the main context window. The 170-220KB of action file content that used to flood the conversation stays out of sight — the main thread only handles routing and receives a summary. `work` and `cleanup` run in the background so you get your conversation back immediately.

- Replaced "Action References" with "Action Dispatch" in SKILL.md
- Actions dispatched to `general-purpose` Task subagents via prompt pattern
- `work` and `cleanup` run in background (non-blocking)
- `do`, `verify`, `version` run in foreground (blocking)
- Screenshots bridged to `do` subagent via temp files + text descriptions

## 0.10.0 — The Hard Stop (2026-02-16)

Capture no longer slides into execution. The do action now has an explicit boundary: after writing files and reporting back, it stops. No helpful "let me go ahead and start building that for you." The user decides when to run the queue — always. Both SKILL.md (routing level) and do.md (action level) enforce this, so even eager agents get the message.

- Added "Capture ≠ Execute" guardrail to SKILL.md core concepts
- Added "STOP After Capture" section to do.md workflow, before the checklist
- Only exception: user explicitly asks for capture + execution in the same invocation

## 0.9.5 — The Reinstall (2026-02-04)

`npx skills update` silently fails to update files despite reporting success. Switched the update command to `npx skills add bladnman/do-work -g -y` which does a full reinstall and actually works. Also fixed the upstream URL — version checks now hit `version.md` where the version number actually lives.

- Update command changed from `npx skills update` to `npx skills add -g -y` (full reinstall)
- Upstream URL fixed: `SKILL.md` → `actions/version.md`

## 0.9.4 — The Passport (2026-02-04)

Install and update commands are no longer tied to a single CLI tool. Switched from `npx install-skill` / `npx add-skill` to the portable `npx skills` CLI, which works across multiple agentic coding tools. Update checks now point to `npx skills update` instead of a reinstall command.

- README install command updated to `npx skills add bladnman/do-work`
- Version action "update available" message now suggests `npx skills update`
- Fallback/manual update uses `npx skills add` instead of `npx install-skill`

## 0.9.3 — The Timestamp (2026-02-04)

Every changelog entry now carries a date. Backfilled all existing entries from git history so nothing's undated. Future entries get dates automatically — the CLAUDE.md format template and rules were updated to enforce it.

- Added `(YYYY-MM-DD)` dates to all 12 existing changelog entries via git history
- Updated CLAUDE.md changelog format template to include date
- Added "Date every entry" rule to changelog guidelines

## 0.9.2 — The Front Door (2026-02-04)

The SKILL.md frontmatter was broken — missing closing delimiters and markdown syntax mixed into the YAML. The `add-skill` CLI couldn't parse the skill metadata properly. Now it's valid YAML frontmatter that tools can actually read.

- Fixed SKILL.md frontmatter: removed `##` from name field, added closing `---`
- Cleaned up upstream URL (was wrapped in a markdown link inside YAML)

## 0.9.1 — The Gatekeeper (2026-02-04)

Keywords like "version" and "changelog" were sneaking past the routing table and getting treated as task content. Fixed by reordering the routing table so keyword patterns are checked before the descriptive-content catch-all, and added explicit priority language so agents match keywords first.

- Routing table now has numbered priority — first match wins, top to bottom
- "Descriptive content" catch-all moved to last position (priority 7)
- Step 2 clarifies that single keywords matching the table are routed actions, not content
- Fixes: `do-work version` no longer asks "Add this as a request?"

## 0.9.0 — The Rewind (2026-02-04)

You can now ask "what's new" and actually see what's new — right at the bottom of your terminal where you're already looking. The version action gained changelog display with a twist: it reverses the entries so the latest changes land at the bottom of the output, no scrolling required. Portable across skills — any project with a CHANGELOG.md gets this for free.

- Changelog display added to the version action: `do-work changelog`, `release notes`, `what's new`, `updates`, `history`
- Entries print oldest-to-newest so the most recent version appears at the bottom of terminal output
- Routing table updated with changelog keyword detection
- Works with any skill that has a CHANGELOG.md in its root

## 0.8.0 — The Bright Light (2026-02-03)

The UR system was hiding in plain sight — documented everywhere but easy to miss if you weren't reading carefully. This release restructures the do action and skill definition so the UR + REQ pairing is unmissable, even for agents that skim. Also added agent compatibility guidance to CLAUDE.md so future edits keep the skill portable across platforms.

- Added "Required Outputs" section to top of do.md — UR + REQ pairing stated upfront as mandatory
- Restructured Step 5 Simple Mode — UR creation now has equal weight with REQ creation
- Added Do Action Checklist at end of workflow — mirrors the work action's orchestrator checklist
- Moved UR anti-patterns to general "What NOT To Do" section (was under complex-only)
- Updated SKILL.md with core concept callout about UR + REQ pairing
- Added Agent Compatibility section to CLAUDE.md — generalized language, standalone-prompt design, floor-not-ceiling

## 0.7.0 — The Nudge (2026-02-01)

Complex requests now get a gentle suggestion to run `/do-work verify` after capture. If your input had lots of features, nuanced constraints, or multiple REQs, the system lets you know verification is available — so you can catch dropped details before building starts. Simple requests stay clean and quiet.

- Verify hint added to do action's report step for meaningfully complex requests
- Triggers on: complex mode, 3+ REQ files, or notably long/nuanced input
- Two complex examples updated to show the hint in action
- No change for simple requests — no hint, no noise

## 0.6.0 — The Bouncer (2026-02-01)

Working and archive folders are now off-limits. Once a request is claimed by a builder or archived, nobody can reach in and modify it — not even to add "one more thing." If you forgot something, it goes in as a new addendum request that references the original. Clean boundaries, no mid-flight surprises.

- Files in `working/` and `archive/` are now explicitly immutable
- New `addendum_to` frontmatter field for follow-up requests
- Do action checks request location before deciding how to handle duplicates
- Work action reinforces immutability in its folder docs

## 0.5.0 — The Record Keeper (2026-02-01)

Now you can see what changed and when. Added this very changelog so the project has a memory. CLAUDE.md got updated with rules to keep it honest — every version bump gets a changelog entry, no exceptions.

- Added `CHANGELOG.md` with full retroactive history
- Updated commit workflow: version bump → changelog entry → commit

## 0.4.0 — The Organizer (2026-02-01)

The archive got a brain. New **cleanup action** automatically tidies your archive at the end of every work loop — closing completed URs, sweeping loose REQs into their folders, and herding legacy files where they belong. Also introduced the **User Request (UR) system** that groups related REQs under a single umbrella, so your work has structure from capture to completion.

- Cleanup action: `do-work cleanup` (or automatic after every work loop)
- UR system: related REQs now live under UR folders with shared context
- Routing expanded: cleanup/tidy/consolidate keywords recognized
- Work loop exit now triggers automatic archive consolidation

## 0.3.0 — Self-Aware (2026-01-28)

The skill learned its own version number. New **version action** lets you check what you're running and whether there's an update upstream. Documentation got a glow-up too.

- Version check: `do-work version`
- Update check: `do-work check for updates`
- Improved docs across the board

## 0.2.0 — Trust but Verify (2026-01-27)

Added a **testing phase** to the work loop and clarified what the orchestrator is (and isn't) responsible for. REQs now get validated before they're marked done.

- Testing phase baked into the work loop
- Clearer orchestrator responsibilities
- Better separation of concerns

## 0.1.1 — Typo Patrol (2026-01-27)

Fixed a username typo in the installation command. Small but important — can't install a skill if the command is wrong.

- Fixed: incorrect username in `npx install-skill` command

## 0.1.0 — Hello, World (2026-01-27)

The beginning. Core task capture and processing system with do/work routing, REQ file management, and archive workflow.

- Task capture via `do-work <description>`
- Work loop processing with `do-work run`
- REQ file lifecycle: pending → working → archived
- Git-aware: auto-commits after each completed request

