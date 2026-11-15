---
id: REQ-735
title: '[impact-rule-change] Decide, rule by rule, which moved standing rules load in ordinary sessions'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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
