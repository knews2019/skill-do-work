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
claimed_at: 2026-11-15T07:30:04Z
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
