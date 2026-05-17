# Changelog (archive: spring 2026, 0.50.0 – 0.64.x)

This archive contains older entries from the do-work changelog. The live changelog is in [CHANGELOG.md](./CHANGELOG.md). Even older entries are in [CHANGELOG-pre-0.50.md](./CHANGELOG-pre-0.50.md).

---

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

