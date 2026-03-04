# Changelog

What's new, what's better, what's different. Most recent stuff on top.

---

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

