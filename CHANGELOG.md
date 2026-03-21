# Changelog

What's new, what's better, what's different. Most recent stuff on top.

---

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

## 0.27.0 — The Cartographer (2026-03-12)

The work orchestrator now speaks prime files. Plan and implementation agents receive prime files as first-class context, and the builder is instructed to create missing ones on the fly. The archived REQ example also carries the new field for reference.

- Added `prime_files: []` to Request File Schema YAML example
- Updated Step 4 (Planning): Plan agent now receives prime files and uses them as the strict index
- Updated Step 6 (Implementation): general-purpose agent receives prime files alongside domain rules
- Added Prime Files bullet to agent instructions — read first, create if missing, keep low-noise
- Added `prime_files: []` to the Archived Request File Example frontmatter

## 0.26.0 — The Compass (2026-03-12)

Capture now knows about prime files — semantic index files that point agents to the right source code. REQs carry a `prime_files` array in frontmatter, the PLAN phase reads them alongside agent rules, and Step 1 routes to matching prime files automatically.

- Added `prime_files: []` field to Simple REQ YAML frontmatter
- Updated PLAN checkbox to read listed `prime_files` and agent rules
- Added prime file routing bullet to Step 1: Parse and Assess
- Updated Step 5 item 2 to populate `prime_files` with discovered paths

## 0.25.1 — The Billboard (2026-03-12)

README and SKILL.md now advertise the two new features. Users browsing the docs will see Human UAT under Review Work and Interactive Explainer under Present Work without digging into the action files.

- README: Added Human UAT bullet to Review Work section
- README: Added Interactive Explainer bullet to Present Work section
- SKILL.md: Updated help menu description for `do work present work`

## 0.25.0 — The Show Floor (2026-03-12)

Present work now generates an interactive HTML explainer alongside the client brief and video. It's a single `.html` file — no build steps, no npm — that stakeholders can double-click to open in any browser.

- Added section 4c: Interactive Explainer (Single-File HTML) to `present-work.md`
- Zero dependencies: HTML5 + Tailwind CDN + Vanilla JS in one file
- Includes Before/After toggle, step-by-step architecture walkthrough, and value summary
- Updated Step 5 summary to list the HTML file with double-click-to-open instructions
- Renumbered Portfolio artifacts from 4c to 4d

## 0.24.0 — The Feedback Loop (2026-03-12)

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

Commits weren't happening during `do work run` because the architecture diagram — the visual flow agents follow — never mentioned Step 9 (Commit). The detailed instructions existed but agents never reached them. Now the diagram shows Commit as an explicit step after Archive.

- **work.md**: Added "Commit (git repos only)" node to the architecture diagram between Archive and Loop
- **work.md**: Added bold reminder callout below the diagram reinforcing that every completed request gets a commit before looping

## 0.22.6 — The Safety Net (2026-03-04)

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

Two fixes from PR review. CHANGELOG.md moved back to root so `do work changelog` works for installed users (it was accidentally excluded with the `_dev/` move). Standalone reviews now find UR input regardless of whether the UR has been archived yet.

- Moved `CHANGELOG.md` back to root — the changelog command needs it at install time
- Fixed `review-work.md` Step 3: UR input lookup now checks `user-requests/` first, falls back to `archive/` — works in both pipeline and standalone modes regardless of UR completion state

## 0.22.3 — The Lighter Carry (2026-03-04)

Sample deliverables no longer ship with the skill. The three do-work-specific outputs (client brief, video script, portfolio summary) moved to `_dev/deliverables/` so they're excluded from installation. Users generate their own via `do work present`.

- Moved 3 sample deliverable files from `do-work/deliverables/` to `_dev/deliverables/`

## 0.22.2 — The Tidy Install (2026-03-04)

Dev-only files no longer tag along when someone installs the skill. CLAUDE.md, CHANGELOG.md, AGENTS.md, and GEMINI.md now live in `_dev/` — excluded by the skills CLI's underscore convention. A root symlink keeps CLAUDE.md discoverable for repo development.

- Moved 4 dev-only files to `_dev/` directory (excluded during `npx skills add` installation)
- Added `CLAUDE.md` symlink at repo root for Claude Code auto-discovery
- Updated CLAUDE.md changelog path reference to `_dev/CHANGELOG.md`

## 0.22.1 — The Signpost (2026-03-04)

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

## 0.21.0 — The Consistency Pass (2026-03-04)

Ten cross-file inconsistencies and instruction gaps cleaned up. Agents following these docs literally should now get consistent behavior across all action files.

- Fixed `cleanup.md` self-contradiction — Pass 1 no longer references `do-work/working/` (work action handles its own files before cleanup runs)
- Fixed stale "Critical" severity in `README.md` — the defined levels are Important/Minor/Nit
- Aligned follow-up REQ creation rules — `work.md` now matches `review-work.md`: follow-ups are per-Important-finding, not score-gated
- Added overall review score formula to `review-work.md` — average of percentage dimensions with Risk/Acceptance modifiers
- Closed 200–500 word classification gap in `capture.md` — removed the >500 word floor from Complex (features/constraints matter more than word count)
- Fixed `review-work.md` Step 2: "Plan (if Route B/C)" → "Plan (if Route C)" — planning is Route C only
- Fixed `capture.md` leading-slash references (`/do-work verify requests` → `do work verify requests`)
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

## 0.20.5 — The Gap Closer (2026-03-03)

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
  - **Portfolio mode** (`do work present all`) — summary of all completed work with cumulative value proposition and cross-project lessons learned
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

## 0.19.0 — The Full Picture (2026-03-02)

Review and verify got proper names and bigger jobs. "Verify" becomes "verify requests" — it checks capture quality. "Review" becomes "review work" — and now it does requirements checking (did we build what was asked?), code review, acceptance testing (actually run the thing), and suggests additional testing the user should do. Every action now ends with suggested next prompts so you always know what to do next.

- Renamed `verify.md` → `verify-requests.md`; action name is now "verify requests" across routing, dispatch, help menu, and examples
- Renamed `review.md` → `review-work.md`; action name is now "review work" — enhanced with three new phases:
  - **Requirements check** — walks through every REQ requirement line-by-line to confirm it was delivered
  - **Acceptance testing** — runs the app/tests and verifies the feature works end-to-end, not just in the diff
  - **Suggested additional testing** — recommends manual verification, integration, edge cases, and environment-specific checks
- Review report now includes a requirements checklist, acceptance result (Pass/Partial/Fail/Untested), and suggested testing section
- Added "Suggest Next Steps" section to SKILL.md — every action now ends with 2-3 fully qualified prompt suggestions (`do work verify requests`, not just `verify`)
- Updated capture.md, work.md, README.md with new action names and references

## 0.18.0 — The Clarity Pass (2026-03-02)

Actions now say what they mean. "Capture" becomes "capture requests," the confusing "answers mode" becomes "clarify questions" with `do work clarify`, and bare `do work` shows a help menu instead of jumping straight to the work loop.

- Renamed action: "capture" → "capture requests" across SKILL.md, README, dispatch table, capture.md, work.md
- Renamed "answers mode" → "clarify questions" — new primary verb is `do work clarify` (old verbs still work)
- Bare invocation (`do work` with no arguments) now shows a help menu with sample prompts instead of asking "Start the work loop?"
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

The Open Questions system now has a complete lifecycle — from capture to drain. Five improvements tighten the feedback loop: verify resolves ambiguous questions on the spot (user is present, why not ask?), `do work answers` gives users a dedicated command to batch-review accumulated questions, follow-up REQ creation moves to the archive step so timing is unambiguous, confirmed builder choices skip the work loop entirely, and verify's question handling is explicitly documented as different from review's.

- `verify.md`: Ambiguous gaps now get presented to the user immediately during verify — resolve on the spot, defer, or leave for the builder
- `SKILL.md` + `work.md`: New "answers mode" — `do work answers`/`questions`/`pending` presents all `pending-answers` REQs for batch review
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
- Two modes: **pipeline** (auto-triggered in the work loop after tests pass) and **standalone** (manual via `do work review`)
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

## 0.11.1 — The Safety Net (2026-02-24)

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
- Fixes: `do work version` no longer asks "Add this as a request?"

## 0.9.0 — The Rewind (2026-02-04)

You can now ask "what's new" and actually see what's new — right at the bottom of your terminal where you're already looking. The version action gained changelog display with a twist: it reverses the entries so the latest changes land at the bottom of the output, no scrolling required. Portable across skills — any project with a CHANGELOG.md gets this for free.

- Changelog display added to the version action: `do work changelog`, `release notes`, `what's new`, `updates`, `history`
- Entries print oldest-to-newest so the most recent version appears at the bottom of terminal output
- Routing table updated with changelog keyword detection
- Works with any skill that has a CHANGELOG.md in its root

## 0.8.0 — The Clarity Pass (2026-02-03)

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

- Cleanup action: `do work cleanup` (or automatic after every work loop)
- UR system: related REQs now live under UR folders with shared context
- Routing expanded: cleanup/tidy/consolidate keywords recognized
- Work loop exit now triggers automatic archive consolidation

## 0.3.0 — Self-Aware (2026-01-28)

The skill learned its own version number. New **version action** lets you check what you're running and whether there's an update upstream. Documentation got a glow-up too.

- Version check: `do work version`
- Update check: `do work check for updates`
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

- Task capture via `do work <description>`
- Work loop processing with `do work run`
- REQ file lifecycle: pending → working → archived
- Git-aware: auto-commits after each completed request

