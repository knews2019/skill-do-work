---
id: REQ-735
title: '[impact-rule-change] Decide, rule by rule, which moved standing rules load in ordinary sessions'
status: completed
created_at: 2026-11-15T07:25:23Z
user_request: UR-179
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
related: [REQ-736]
batch: run-012821-follow-ups
write_set: [skills/do-work/crew-members/communication-style.md, skills/do-work/crew-members/coding-guardrails.md, skills/do-work-toolbox/crew-members/coding-guardrails.md, skills/do-work/crew-members/background-agents.md, skills/do-work-toolbox/crew-members/background-agents.md, skills/do-work-knowledge/crew-members/background-agents.md, skills/do-work/actions/commit.md]
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-11-15T07:34:32Z
  basis:
    - Route B
    - 7-file write set
    - 2 subsystems involved
    - 8 acceptance criteria
claimed_at: 2026-11-15T07:30:04Z
dispatch_at: 2026-11-15T07:40:29Z
builder_handback_at: 2026-11-15T07:54:20Z
integration_at: 2026-11-15T07:54:43Z
review_at: 2026-11-15T08:07:34Z
commit: 3b2d0ff6797cc684e879b920492a951b0ae77085
heavy_verified_at: 2026-11-15T08:08:38Z
heavy_verified_revision: 3b2d0ff6797cc684e879b920492a951b0ae77085
kb_status: pending
completed_at: 2026-11-15T08:09:31Z
release_at: 2026-11-15T08:09:31Z
---
# Decide, Rule by Rule, Which Moved Standing Rules Load in Ordinary Sessions
## What
REQ-724 (moving the maintainer's global standing rules into the skill, released 0.322.0) put three rules in files that load only during implementation or an agent fan-out. Decide, for each of the three rules, whether it moves to a file that loads in an ordinary session, and ship that decision. The goal is that the maintainer can then delete `~/.claude/CLAUDE.md` without losing a rule.
## Why
The maintainer's words: "capture a follow-up REQ that decides, rule by rule, which of the three move to an always-on file, so that ~/.claude/CLAUDE.md can then be deleted". In an ordinary session the global file applies these rules today. The skill copy does not, so the delete would silently drop them. The REQ-724 review raised this as F6 and asked for a decision before the delete step.
## Context
"Always-on file" here means `skills/do-work/crew-members/communication-style.md`, the one file the installer links from a consumer project's agent instructions, or a file that it tells the agent to load on a named condition.

The three rules and where they live now (verified at capture, commit `db458e9e`):

| Rule | Current home | Loads when |
| --- | --- | --- |
| R1 approximate numeric targets | `skills/do-work/crew-members/coding-guardrails.md:22` ("Read a number in the request ... as approximate unless the user says strict.") | implementation only (Step 6, per the file's JIT_CONTEXT) |
| R2 never remove existing UI while redesigning | `skills/do-work/crew-members/coding-guardrails.md:89` | implementation only |
| R3 pick the model per agent | `skills/do-work/crew-members/background-agents.md:190-193` (paragraph "Pick each agent's model by its work.") | agent fan-out only |

Lock-step copies of the same lines, which REQ-724 decision D-14 kept identical: `skills/do-work-toolbox/crew-members/coding-guardrails.md:22` and `:87`, `skills/do-work-toolbox/crew-members/background-agents.md:192`, and `skills/do-work-knowledge/crew-members/background-agents.md:192`.

The maintainer's earlier answer on R3: REQ-724 decision D-05, "confirmed by the maintainer, 2026-11-15", says the model rule "lives only in `crew-members/background-agents.md`, with no copy or pointer line in `communication-style.md`" (`do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md:175`).

Open review findings from REQ-724 (`do-work/runs/work-2026-11-15-012821/REQ-724-review.md:26` and `:28`, re-stated in `REQ-724-review-delta.md:45` and `:47`):
- F3: the ask-tool rule (`skills/do-work/crew-members/communication-style.md:25`) and the model rule (`background-agents.md:192` and its copies) do not say that an action's own rule wins. Two shipped texts already say otherwise: `skills/do-work-toolbox/actions/tutorial.md:355` prints plain-text menus and does not use the ask-user tool, and `skills/do-work/actions/ai-evaluation-loop-reference.md:19-20` names the reviewer models for `sol-loop` and `council`. The review suggested "unless the action in use says how to ask" and "keep the session model unless the action names another model".
- F5: the rule group about several sessions in one checkout is only partly shipped. Exact-path staging ships (`skills/do-work/actions/commit.md:222`, under `## Rules` at `:218`). "Never revert files you did not change" has no home in the skill, and `git commit -a` is not named. The review suggested one bullet in `commit.md` § Rules: "Never revert, stash or discard a change you did not make."

Sources: `do-work/runs/work-2026-11-15-012821/REQ-724-integration-report.md` and `do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md`.
## Detailed Requirements
- Decide R1, R2 and R3 one by one. For each, either give it a home that loads in an ordinary session, or keep it where it is with a stated reason why an ordinary session does not need it. Record each decision in this REQ with its value and its risk.
- R3: respect D-05 (no copy or pointer line in `communication-style.md`). If the builder decides a pointer is now needed, the REQ must say plainly why D-05 no longer holds.
- When a rule moves, it is stated once. Do not leave a second copy in its old home; a short pointer is allowed only where the old home's readers need it.
- Keep the lock-step copies in the toolbox and knowledge packages identical to the do-work copy after the change.
- Fold in F3 if it fits: add the "an action's own rule wins" condition to the ask-tool rule and the model rule, so they no longer contradict `tutorial.md:355` and `ai-evaluation-loop-reference.md:19-20`.
- Fold in F5 if it fits: give "never revert, stash or discard a change you did not make" one shipped home.
- If F3 or F5 does not fit this REQ, record why in the REQ and leave it report-only. Do not mint another REQ for it.
## Constraints
- Never read, edit or delete `~/.claude/CLAUDE.md` or `~/.claude/settings.json`. This REQ touches no file outside the repository. Deleting the global file is the maintainer's own step after this REQ ships.
- Write each rule as a condition, with no machine, vendor, model or tool name.
- Nothing under `skills/` cites the maintainer doc `CLAUDE.md`.
- The Aliases and Examples sections of `communication-style.md` do not change. Keep that file short: it loads in every session.
- Do not reword a shipped rule except to add the precedence condition from F3.
- The change touches shipped files under `skills/`, so it is a release. Follow `_dev/primes/prime-releases.md`.
- Independent of REQ-736 (the help.md judge line). No `depends_on`.
## Dependencies
None. REQ-724 is already completed and archived.
## Builder Guidance
- Certainty: the maintainer firmly chose to decide rule by rule. The outcome for each rule is open; the builder decides with stated reasons, using the DECIDE & STATE and ESCALATE tiers of `coding-guardrails.md`.
- Delete before you add: prefer a move or one added condition over new text.
## Red-Green Proof
**RED prompt/case:** Load only the files an ordinary session loads (`communication-style.md` and what it names on a condition). Search them for each of R1 ("unless the user says strict"), R2 ("while redesigning") and R3 ("cheapest model tier"). All three searches find nothing, and nothing in this REQ records why that is fine.
**Why RED now:** REQ-724 placed R1 and R2 in `coding-guardrails.md` (implementation only) and R3 in `background-agents.md` (fan-out only). With `~/.claude/CLAUDE.md` deleted, no ordinary session applies them.
**GREEN when:** For each of R1, R2 and R3, either the search finds the rule in an ordinary-session file, or this REQ records a decision with a reason that it stays where it is. R3's outcome matches D-05, or the REQ says plainly why a pointer is needed. F3 and F5 are each folded in or recorded as not fitting. The lock-step copies match.
**Validation:** Inferred during capture
## Full Context
See `do-work/user-requests/UR-179/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes crew-member files with lock-step copies, the `alternate-writer-contract-drift` family that REQ-724's review F1 hit.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** apply D-01 to D-05 as written: two byte-identical line moves, two precedence clauses (one in three lock-step copies), one new commit rule. No other wording. (from the builder hand-back)
- [x] **[APPLY]:** one Python pass with an exact-once assertion per substitution over the seven files, then one commit by exact path (`82db9eb2`). (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat 547d4cc5..HEAD` lists the seven Scope files only; checked each hunk: R1/R2 gone from both guardrails copies, present once in communication-style.md after the bug-report bullet; the three model paragraphs identical, with "cheapest model tier" still on one line; commit.md bullet sits right after the `git add -A` bullet, which is unchanged; no debug artifacts. (from the builder hand-back)
*Source: item 1 of `do-work/user-requests/UR-179/input.md` (rule scope after REQ-724)*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is a per-rule decision plus a few one-line prose moves, but each rule's home depends on which files load when, so the load sites must be checked first. Prose only, two subsystems (crew-member files with lock-step package copies, and the commit action), no Go or shell change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Pre-dispatch agent, base `547d4cc5`. Every pointer below was re-read at that commit.

**Required-lessons consult.** `required_lessons` stays `[_dev/primes/lessons-releases.md]` (856 tokens, read). Its families (`canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`) do not fire: no file moves and no manifest changes, only lines inside files. `_dev/primes/lessons-action-files.md` stays dropped for budget (8761 tokens, `slugged: partial`, so it cannot be narrowed). Two of its bullets match and were read anyway: `[family: alternate-writer-contract-drift]` (the lock-step copies, REQ-724 review F1) and `[family: condition-preserving-prose-extraction]` (REQ-509: a moved rule must keep its condition, and every reader of the old home must load the new home). The second one is answered below under "Who reads the old homes". No other index row matches (no shell, Go or board change).

**Who loads which file (the evidence for D-01 to D-03):**
- `skills/do-work/agent-instructions.template.md:4` is the always-on link: an installed project's agent instructions say to read `crew-members/communication-style.md` before the first reply and follow it in every response. In this repository `CLAUDE.md` imports it. `communication-style.md` exists only in the do-work package, so it has no lock-step copy.
- `skills/do-work/crew-members/coding-guardrails.md:3` (JIT_CONTEXT): loaded during implementation (Step 6) only. `skills/do-work/actions/work-implementation-reference.md:9` loads it, and `:12` (item 2b) also always loads `communication-style.md` in the same step.
- `skills/do-work/crew-members/background-agents.md:3` (JIT_CONTEXT): loaded by any action that fans work out to sub-agents. Not loaded for a single agent in an ordinary session.
- `skills/do-work/actions/commit.md` § Rules (`:218-224`) is the canonical commit-staging guard. Every other action that stages or commits points there (`:220`).

**Who reads the old homes (condition-preserving check for R1 and R2).** The do-work `coding-guardrails.md` is read in Step 6, and Step 6 also loads `communication-style.md` (2b), so a line moved there still reaches every implementation. The toolbox copy `skills/do-work-toolbox/crew-members/coding-guardrails.md` is cited by toolbox actions for review criteria (`code-review.md:140` § Simplicity First, `quick-wins.md:181` surgical changes, `maintainability-audit.md:106` and `validate-feedback.md:84` point at the do-work copy). None of them relies on the numeric or redesign line. Those sessions get the moved lines through the always-on link.

**The three rules at base:**
- R1: `skills/do-work/crew-members/coding-guardrails.md:22` and `skills/do-work-toolbox/crew-members/coding-guardrails.md:22`, byte-identical: `- Read a number in the request (a size, a count, a time) as approximate unless the user says strict.`
- R2: `skills/do-work/crew-members/coding-guardrails.md:89` and `skills/do-work-toolbox/crew-members/coding-guardrails.md:87`, byte-identical: `- Never remove existing UI while redesigning unless the user asked for it.` Nearby but different: `crew-members/ui-design.md:58` ("Don't redesign what isn't broken") is about review scope, not about removing UI. It stays.
- R3: the paragraph "Pick each agent's model by its work." at `skills/do-work/crew-members/background-agents.md:190-193`, `skills/do-work-toolbox/crew-members/background-agents.md:191-194`, `skills/do-work-knowledge/crew-members/background-agents.md:191-194`. The three paragraphs hash the same (checked with `awk` + `shasum` at base).

**F3 targets:** the ask-tool bullet `skills/do-work/crew-members/communication-style.md:25`, first sentence ("Use the agent's ask-user tool when it has one, never a prose list of questions, and load `crew-members/clear-questions.md` first."). The model sentence "Planning, review and build agents keep the session model." in the three R3 paragraphs. The two shipped texts it contradicts today: `skills/do-work-toolbox/actions/tutorial.md:355` (plain-text menus, no ask-user tool) and `skills/do-work/actions/ai-evaluation-loop-reference.md:19-21` (named reviewer models for `sol-loop`, `council`, `--jester`).

**F5 target:** `skills/do-work/actions/commit.md` § Rules, after the exact-path staging bullet at `:222`. Partial coverage elsewhere, not changed: `skills/do-work/actions/work-reference.md:31` (during a run, preserve a diff you did not make) and `skills/do-work/actions/restart-with-parallel-handoff.md:68` (files you did not change stay untouched). `commit.md` has no package copy.

**Size:** `communication-style.md` is 99 lines at base. The change adds two moved bullets and two clauses there. `_dev/tests/instruction-reading-volume.py:43-44` counts both `coding-guardrails` and `communication-style` in the build crew set, so the move is volume-neutral for implementation.

**Negative check at base:** `zsh|AskUserQuestion|haiku|EnterWorktree|CLAUDE\.md` has no hit in the seven Scope files.

**RED searches (at `547d4cc5`), the probe's fixed phrases:**
1. `unless the user says strict` is in the two `coding-guardrails.md` copies, not in `communication-style.md`.
2. `while redesigning`: same.
3. `unless the action in use says how to ask`: no hit in `skills/`.
4. `unless the action names another model`: no hit in `skills/`.
5. `change you did not make`: no hit in `skills/`.

**Pre-dispatch decisions.** The maintainer delegated this judgment (REQ § Builder Guidance), so each ESCALATE-tier choice was made here with value and risk instead of a question, per the unattended-run rule.

- **D-01 (ESCALATE, decided): R1 (approximate numeric targets) moves to `communication-style.md`.** The line moves byte-identical as a plain bullet right after the bug-report bullet (`communication-style.md:27`). It is deleted from § 1 of both `coding-guardrails.md` copies, and no pointer stays there. Reason: the rule is about reading the user's request, and numbers arrive in ordinary sessions before any REQ exists, often for prose (a word count, a line count). That is the same reason REQ-724 D-02 gave for putting the bug-report line in the always-on file. Implementation loses nothing, because Step 6 always loads `communication-style.md` too (`work-implementation-reference.md:12`). Value: the rule applies in every session, stated once. Risk: a consumer project that declined the installer's always-on link no longer gets R1 inside toolbox actions that read only the toolbox `coding-guardrails.md`; do-work Step 6 still loads it explicitly. Reversible: move the line back.
- **D-02 (ESCALATE, decided): R2 (never remove existing UI while redesigning) moves to `communication-style.md`.** Same form as D-01: byte-identical, right after the R1 bullet, deleted from § 3 of both `coding-guardrails.md` copies, no pointer. Rejected alternative: one pointer line "before you change code, load `coding-guardrails.md`". It is one line, but it would load the whole guardrails file (about 130 lines: the four behavioral guardrails, the decide-vs-escalate gate, naming for reach) into every ordinary coding session. That is a much wider behavior change than this REQ asks for, and it contradicts that file's JIT_CONTEXT. Value: a redesign request in an ordinary session keeps the rule, at the cost of one line. Risk: `communication-style.md` holds one more scope rule next to the bug-report line, so its "conversational output" scope stretches a little further. Reversible.
- **D-03 (ESCALATE, decided): R3 (pick the model per agent) stays only in `background-agents.md`. No copy and no pointer in `communication-style.md`. REQ-724 D-05 holds.** Reason: the maintainer confirmed D-05 on 2026-11-15, and D-05 already named this exact risk ("a single lookup agent outside a fan-out never sees the rule"). No new fact has appeared since then. The rule saves cost and does not change results. Overturning a confirmed answer from the same day is not a call to make silently. Value: the maintainer's answer stands, and the always-on file stays about talking. Risk: after `~/.claude/CLAUDE.md` is deleted, an ordinary session that starts one lookup agent runs it on the harness's default model, which costs more tokens but gives the same result. So the delete does lose R3 for single agents in ordinary sessions. Reversible: one pointer line in `communication-style.md` if the maintainer wants it.
- **D-04 (DECIDE & STATE): F3 folds in.** Two precedence clauses, worded as the REQ-724 review suggested. (a) The first sentence of the ask-tool bullet (`communication-style.md:25`) ends with "unless the action in use says how to ask". (b) In all three `background-agents.md` copies, "Planning, review and build agents keep the session model" gains "unless the action names another model". The REQ's constraints allow exactly this rewording. It removes the contradiction with `tutorial.md:355` and `ai-evaluation-loop-reference.md:19-21`.
- **D-05 (DECIDE & STATE): F5 folds in as one new bullet in `commit.md` § Rules,** after the exact-path staging bullet: "Never revert, stash or discard a change you did not make." plus at most one short reason sentence. `commit.md` is the canonical commit-staging guard, so the concurrent-session rules now live together there. `git commit -a` is not added: naming it would reword the shipped `git add -A` bullet, which the constraints forbid, and that bullet's "Stage only the specific files" already rules it out. Not done: an always-on copy of F5. An ordinary session that never commits through the commit action does not load this rule. The concurrent-session rule group was already scoped this way by REQ-724 (G6, exact-path staging lives only in `commit.md`).
- **D-06 (DECIDE & STATE): the builder never edits `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or a version mirror.** It proposes the entry, and the integrator releases. Proposed bump: minor, because R1 and R2 now apply in every session of every installed project (impact-rule-change, same reasoning as REQ-724 D-12).
- **D-07 (DECIDE & STATE): fixed proof phrases** for the RED searches, the GREEN probe and review: the five phrases under "RED searches" above, plus `cheapest model tier` (must stay in the three `background-agents.md` copies only, which proves D-03).

<!-- D-XX counter: last used D-07. Next decision: D-08. -->

*Generated by pre-dispatch agent (in-context exploration)*

## Scope

**Files I will touch:**
- `skills/do-work/crew-members/communication-style.md` (modify) — two moved bullets after the bug-report bullet: R1 approximate numbers (D-01) and R2 no UI removal while redesigning (D-02), byte-identical to their old lines; the ask-tool bullet's first sentence gains the precedence clause (D-04 a)
- `skills/do-work/crew-members/coding-guardrails.md` (modify) — delete the R1 line from § 1 and the R2 line from § 3
- `skills/do-work-toolbox/crew-members/coding-guardrails.md` (modify) — the same two deletions (lock-step copy)
- `skills/do-work/crew-members/background-agents.md` (modify) — the model paragraph gains the precedence clause (D-04 b), rewrapped
- `skills/do-work-toolbox/crew-members/background-agents.md` (modify) — the same paragraph, byte-identical to the do-work copy
- `skills/do-work-knowledge/crew-members/background-agents.md` (modify) — the same paragraph, byte-identical to the do-work copy
- `skills/do-work/actions/commit.md` (modify) — one new bullet in the Rules section after the exact-path staging bullet (D-05)

**Files I will NOT touch:** `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` and the version mirrors (the integrator releases, D-06); the Aliases and Examples sections of `communication-style.md`; `skills/do-work-toolbox/actions/tutorial.md` and `skills/do-work/actions/ai-evaluation-loop-reference.md` (cited, not changed); `skills/do-work/crew-members/ui-design.md`; any file outside the repository; no new file and no test file.

**Acceptance criteria (restated from REQ):**
- [ ] R1, R2 and R3 each have a recorded decision with value and risk (D-01, D-02, D-03).
- [ ] R1 and R2 are each stated once, in `communication-style.md`, with no copy and no pointer left in either `coding-guardrails.md` copy.
- [ ] R3 stays only in the three `background-agents.md` copies, with no copy or pointer in `communication-style.md` (REQ-724 D-05 respected).
- [ ] F3 is folded in: the ask-tool rule and the model rule each say that an action's own rule wins.
- [ ] F5 is folded in: one shipped home for never reverting, stashing or discarding a change you did not make.
- [ ] Lock-step copies stay identical: the model paragraph matches in all three packages, and both `coding-guardrails.md` copies lose the same two lines.
- [ ] Aliases and Examples of `communication-style.md` are unchanged; no machine, vendor, model or tool name; nothing under `skills/` cites the maintainer doc.
- [ ] The release follows the releases prime (proposed by the builder, written by the integrator).

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-735-preflight-probe.sh) and the repository gate green at 547d4cc5 (147 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*

## Decisions

(from the builder hand-back)

- **D-08 (DECIDE & STATE):** The F5 reason sentence is "Another session or the user may own it." It is the example the brief gave, one short sentence, and it names the condition (someone else owns the change) without naming a tool.
- **D-09 (DECIDE & STATE):** The model paragraph rewrap breaks after "unless the action", so lines 1 and 2 stay unchanged and the new line 3 is 81 characters, within the paragraph's existing width (line 1 is 86).

## Implementation Summary

**Files changed:**
- `skills/do-work/crew-members/communication-style.md` (modified: R1 and R2 bullets added after the bug-report bullet; the ask-tool bullet gains "unless the action in use says how to ask")
- `skills/do-work/crew-members/coding-guardrails.md` (modified: R1 line deleted from § 1, R2 line deleted from § 3)
- `skills/do-work-toolbox/crew-members/coding-guardrails.md` (modified: the same two deletions)
- `skills/do-work/crew-members/background-agents.md` (modified: the model sentence gains "unless the action names another model", rewrapped)
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified: byte-identical paragraph)
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified: byte-identical paragraph)
- `skills/do-work/actions/commit.md` (modified: one new bullet in § Rules after the exact-path staging bullet)

**What was done:** R1 (read numbers as approximate) and R2 (never remove existing UI while redesigning) moved byte-identical from both `coding-guardrails.md` copies into the always-on `communication-style.md`, with no pointer left behind (D-01, D-02). R3 (pick the model per agent) stays only in the three `background-agents.md` copies, so REQ-724 D-05 holds (D-03). F3 folded in: the ask-tool rule and the model rule now yield to an action's own rule (D-04). F5 folded in: `commit.md` § Rules says never to revert, stash or discard a change you did not make (D-05). The builder changed no CHANGELOG, VERSION or mirror (D-06). Diff: seven files, +10/-11, builder commit `82db9eb2`. Review F1 fix (integrator, on the builder branch, commit `0b10c310`, re-merge `3b2d0ff6`): the ask-tool precedence clause moved into the rule's bold lead, so it covers the whole rule; no word added or removed. Cumulative range `08de436c..3b2d0ff6`.

## Discovered Tasks

(from the builder hand-back)

None.

## Qualification

**Gate records:** `advance --diff-range 08de436c..e3f36c6c` returned phase `qualify`, gate records `qualify` and `scope-drift` satisfied. The first call refused `QUALIFY-CLAIMED-PATH-MISSING` and `SCOPE-UNDECLARED-TOUCH` because my Implementation Summary backticked the command text "git add -A", which the qualifier reads as a path; I reworded that phrase and changed nothing else.

**Requirement trace against `git diff 08de436c..e3f36c6c` (7 files, +10/-11):**
- R1, R2 and R3 each have a recorded decision with value and risk: D-01, D-02 and D-03 in § Exploration.
- R1 and R2 are stated once, in the always-on file: `skills/do-work/crew-members/communication-style.md:28` and `:29` hold the old lines byte for byte. Both `coding-guardrails.md` copies lose exactly those two lines, with no pointer left (diff hunks at § 1 and § 3).
- R3 stays only in the three `background-agents.md` copies; `communication-style.md` has no "model tier", "background-agents" or "model by its work" text (probe check). REQ-724 D-05 holds.
- F3: `communication-style.md:25` ends its first sentence with "unless the action in use says how to ask"; the model paragraph in all three `background-agents.md` copies reads "keep the session model unless the action names another model".
- F5: `skills/do-work/actions/commit.md:223` says "Never revert, stash or discard a change you did not make." It sits right after the unchanged exact-path staging bullet.
- Lock-step: the three model paragraphs are byte-identical (probe), and the two guardrails copies drop the same two lines.
- Aliases and Examples of `communication-style.md` are unchanged against `547d4cc5`; no machine, vendor, model or tool name in the seven files; nothing under `skills/` cites the maintainer doc (probe and `shipped-package-reference-contract.sh`).
- Release: proposed by the builder (minor, impact-rule-change), written by the integrator at finalization.

**Scope:** Route B. The declared Scope list names seven files and the diff touches exactly those seven. No CHANGELOG, VERSION or mirror hunk. No debug artifacts. `communication-style.md` grows from 99 to 101 lines.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge `e3f36c6cc20db8dbf8820d59f9dee040fdbc1011` (detached drain-head checkout, `DO_WORK_FAST_STAGE_REUSE=off`, 1-minute load 4.10 at start, no other gate running); GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-735-probe.sh` through `advance --probe-file`.
**Result:** ✓ Gate exit 0, 123 s, "Maintainer verification passed." Probe exit 0 ("REQ-735 probe PASS", 15 PASS lines); gate records `run-blocked-check` and `green-gate` satisfied. After the review F1 fix, the gate ran again at re-merge `3b2d0ff6797cc684e879b920492a951b0ae77085` (1-minute load 4.62): exit 0, 127 s; the green record was written with `record-green-gate` at that revision, and the probe passed again.

**Pre-flight baseline (PREFLIGHT-BASELINE-RETRIED):** the pre-build probe `REQ-735-preflight-probe.sh` (shipped-package reference contract plus contract regressions) at base `547d4cc5` failed on its first attempt (non-zero exit, under load from a parallel agent; the coordinator's finding does not keep the exact status) and passed on the retry (exit 0). The repository gate at `547d4cc5` was green (147 s).

**Red-green validation:** (from the builder hand-back)
- `REQ-735-probe.sh`: ✗ at base `547d4cc5` (exit 1, 7 FAIL lines: R1, R2, F3-ask, F5, the two moved lines, F3-model) → ✓ at tip `82db9eb2` (exit 0, "REQ-735 probe PASS").

**New tests added:**
- None (prose moves and clauses; the run probe is run evidence, not a shipped test).

**Heavy verification plan:** `plan-heavy-verification` over 08de436c86f851633660d785fceb52c29d81dc0e..e3f36c6cc20db8dbf8820d59f9dee040fdbc1011 selected one lane (the re-merge to `3b2d0ff6` touches only `communication-style.md`, already covered by the same lane):
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: each of the seven changed files matched subtree `skills`.

## Review

**Overall: 94%** | 2026-11-15T08:07:34Z (delta re-review after the F1 fix: 96%, Approve)

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 `skills/do-work/crew-members/communication-style.md:25`: "unless the action in use says how to ask" reads as scoping only the clear-questions load, and the bold lead stays unconditional against `tutorial.md:355`. Fixed before release: the integrator applied the suggested line on the builder branch (`0b10c310`), re-merged as `3b2d0ff6`, re-ran the gate (green, 127 s) and had the delta re-reviewed (`REQ-735-review-delta.md`, 96%, Approve). — impact-rule-change → report only
- F2 `skills/do-work/actions/commit.md:223`: the never-revert rule loads only with the commit action, so an ordinary session loses it once the global file is deleted (D-05). This is a maintainer decision before that delete. — impact-rule-change → report only
- F3 `skills/do-work/crew-members/communication-style.md:7`: the "conversational output" scope sentence does not cover R1 and R2 (D-02 risk). — impact-negligible → report only
- F4 `skills/do-work/docs/commit-guide.md:26-32`: the Key rules summary omits the new revert rule. — impact-negligible → report only
- F5 `CHANGELOG.md:25-26`, `skills/do-work/CHANGELOG.md:25-26`: the 0.322.0 history still names `coding-guardrails.md` as the home of R1 and R2. History is not rewritten, and the release entry must state the move. — impact-negligible → report only
- F6 `skills/do-work-toolbox/actions/architecture-report.md:81`: a full redesign each run, now co-loaded with R2. The conflict is weak. — impact-negligible → report only
- F7 `CHANGELOG.md:23`, `skills/do-work/CHANGELOG.md:23` (delta review): the 0.322.0 history states the ask rule without the new condition. History is not rewritten; the REQ-735 release entry names the new precedence condition. — impact-negligible → report only

**Acceptance:** Pass — the implementation and integration stages were exercised: the REQ-735 probe (15 PASS, exit 0) and `shipped-package-reference-contract.sh` (exit 0) at merge `e3f36c6c`, and again at re-merge `3b2d0ff6` (delta review). Deployment and live acceptance are unassessed.
**Restatement sweep:** redefined the home of R1 and R2 (now `communication-style.md:28-29`; no live restatement in `skills/`, `docs/`, `_dev/`, `CLAUDE.md`, `AGENTS.md` or `README.md`; history only at `CHANGELOG.md:25` and `skills/do-work/CHANGELOG.md:25`, F5, plus the input quotes in `docs/run-prompts/2026-11-15-move-global-rules-into-skill.md:24,35-36,48`), the ask-tool precedence (now in the bold lead at `communication-style.md:25` after the F1 fix, no longer stale against `tutorial.md:355`; `trace.md:39` and the three `clear-questions.md:3` triggers consistent; history at `CHANGELOG.md:23`, F7), the model-rule precedence (three copies identical, consistent with `ai-evaluation-loop-reference.md:19-21`; history at `CHANGELOG.md:26`), and the new revert rule at `commit.md:223` (consistent with `work-reference.md:31` and `restart-with-parallel-handoff.md:68`; `commit-guide.md:26-32` omits it, F4). Inherited elements from REQ-724 were rechecked with no stale text: the DECIDE & STATE pointer (now `coding-guardrails.md:32`, cited without a line number), the clear-questions option shape (three copies identical apart from package paths), the per-agent model rule, the numeric-target and redesign lines, and the stall-check wording (`communication-style.md:30`). REQ-736 recorded nothing redefined.
**Suggested testing:** 3 items
**Follow-ups created:** None (7 findings report only; F1 fixed)
**Anti-bloat:** 0 unrequested helpers, options, files or tests
**Full report:** `do-work/runs/work-2026-11-15-012821/REQ-735-review.md` and the delta `do-work/runs/work-2026-11-15-012821/REQ-735-review-delta.md`

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: `08de436c86f851633660d785fceb52c29d81dc0e`
- Target: `3b2d0ff6797cc684e879b920492a951b0ae77085` (re-merge after the review F1 fix; the first merge `e3f36c6c` selected the same lane)
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: each of the seven changed files matched subtree `skills`

## Heavy Verification Result

- Target: `3b2d0ff6797cc684e879b920492a951b0ae77085`; execution revision `3b2d0ff6797cc684e879b920492a951b0ae77085` (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set).
- `staged-skills` — executed, exit 0, 31 s.

## Lessons Learned

**What worked:** A fixed-phrase probe that names each rule's only allowed home (and checks the three lock-step copies byte for byte) made the move, the two precedence clauses and the D-03 "no pointer" decision mechanical to verify, before and after the review fix.
**What didn't:** The builder placed the ask-tool precedence clause at the end of the rule's first body sentence, where it read as covering only the clear-questions load; the bold lead stayed unconditional. The review caught it (F1) and the clause moved into the bold lead.
**Worth knowing:** A precedence condition on a rule with a bold lead belongs in the lead, because readers and sweeps treat the lead as the rule. The qualifier reads any backticked text in an Implementation Summary as a path, so a command such as "git add -A" must stay unbackticked there.

## Orientation

`skills/do-work/crew-members/communication-style.md` is the only always-on rules file (linked from consumer agent instructions and imported by this repository's maintainer doc); `coding-guardrails.md` loads only in implementation and `background-agents.md` only in fan-outs, so a rule that must apply in every session lives in the first.

## Timing

Observed 2026-11-15T07:34:32Z to 2026-11-15T08:08:38Z: 34m 06s total, 29m 00s attributed across 7 events, 5m 06s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 13m 51s | 1 |
| review | 5m 30s | 1 |
| verification-gate | 5m 09s | 3 |
| exploration-preflight | 4m 17s | 1 |
| handback-merge | 13s | 1 |

Slowest stage: builder-work / builder worktree build, 13m 51s, outcome success.
