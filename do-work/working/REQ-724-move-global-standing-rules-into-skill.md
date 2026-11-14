---
id: REQ-724
title: '[impact-rule-change] Move the maintainer''s global standing rules into the skill'
status: claimed
created_at: 2026-11-14T22:44:21Z
user_request: UR-175
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
related: [REQ-725]
batch: machine-local-rules
write_set: [skills/do-work/crew-members/communication-style.md, skills/do-work/crew-members/clear-questions.md, skills/do-work/crew-members/coding-guardrails.md, skills/do-work/crew-members/background-agents.md, skills/do-work/CHANGELOG.md]
claimed_at: 2026-11-14T23:28:11Z
---
# Move the Maintainer's Global Standing Rules Into the Skill
## What
Move the standing rules that live only in the maintainer's `~/.claude/CLAUDE.md` into the shipped do-work skill, so every machine and every repository gets them. Add only the rules the skill does not carry yet, in existing crew-member files, with the smallest addition.
## Why
The global file exists on one machine. The maintainer's words: "it does not go to other machines, and I get weird behaviour. This is why I have the common do-work skill that is checked into every repo".
## Context
- The complete source text of the global file (35 lines, copied 2026-11-15) is in `do-work/user-requests/UR-175/input.md`. Use that copy. The build session may not have the file.
- A text search on 2026-11-15 found this state in `skills/`. Re-check every line at HEAD.

| Code | Rule group | State in `skills/` |
| --- | --- | --- |
| G1 | Explaining state | Mostly present in `skills/do-work/crew-members/communication-style.md` ("Answer first, next action last"). Missing: big picture before the exact picture, and a plain definition of every status word on first use. |
| G2 | Decisions through the ask tool | Not found. `skills/do-work/crew-members/clear-questions.md` covers how to word a question, not "use the ask tool, recommended option first, value and risk per option, decide reversible trivia yourself". |
| G3 | Commands for the user | Not found. |
| G4 | Scope | The operator-work rule is shipped (0.305.85). The other three lines (a bug report means fix the code, numeric targets are approximate, never remove UI in a redesign) are not found. |
| G5 | Prompts for another session | Only the folder name, in `skills/do-work-toolbox/actions/folder-structure.md`. The unattended shape, the stall check and the short paste block are not found. |
| G6 | Several sessions in one checkout | Exact-path commits are shipped. The `checkout-lease` hook stays local (maintainer's answer during capture). Nothing to move. |
| G7 | Model per subagent | Not found. |

## Detailed Requirements
- For each source line, first look for shipped text that already says it. Add a line only when none does. Never restate a shipped rule in a second file.
- G1, G2 and G3 must apply outside a work loop. `communication-style.md` is the only file the installer links from a consumer project's agent instructions, so these rules live there, or in a file it tells the agent to load on a named condition.
- Write conditions, not this maintainer's machine. "zsh" becomes "write commands for the user's login shell, and wrap a payload for a remote host in that host's shell". "AskUserQuestion" becomes "the agent's ask-user tool, when it has one". "haiku" becomes "the cheapest model tier that can do the lookup".
- G3 also carries: fill in every value so a pasted command has no placeholder, and logic longer than about 5 lines goes into a committed script with a one-line launcher.
- G5 uses the project's run-prompts folder as `folder-structure.md` defines it. Do not write a second path.
- Suggested homes, which the builder may change with a stated reason: G1 to G3 in `communication-style.md`; the option-shape half of G2 in `clear-questions.md`; G4 in `coding-guardrails.md`; G7 in `background-agents.md`; G5 wherever the skill already describes a handoff to another session.
- Release the change as `_dev/primes/prime-releases.md` requires for shipped files (changelog and version).
## Constraints
- Never read, edit or delete `~/.claude/CLAUDE.md` or `~/.claude/settings.json`. This REQ touches no file outside the repository.
- No new crew-member file and no new always-on file.
- No vendor, tool or host name in the shipped text.
- No shipped rule is reworded.
- The Aliases section and the Examples section of `communication-style.md` do not change.
- About 25 shipped lines in total is the target. It is approximate.
- No new test file. These are prose rules, and a test of their wording would pin nothing real (maintainer confirmed during capture).
## Dependencies
None. REQ-725 (sorting the local memory notes) is related but independent.
## Builder Guidance
- Certainty: the maintainer chose this direction firmly ("move first, delete after"). The homes of the rules are suggestions.
- In this repository the run-prompts folder question is still open (`_dev/run-prompts/2026-11-15-run-prompts-folder-rule.md`). That does not block G5, because G5 points at `folder-structure.md` and names no path of its own.
## Red-Green Proof
**RED prompt/case:** At the starting commit, run one text search in `skills/` for each of G1 (status-word definition), G2, G3, G4, G5 and G7, and record the six commands and their empty results in the REQ.
**Why RED now:** The six rule groups exist only in `~/.claude/CLAUDE.md` on one machine.
**GREEN when:** The same six searches each find exactly one shipped file. A search for `zsh`, `AskUserQuestion`, `haiku` and `EnterWorktree` in the four changed crew-member files finds nothing. The fast gate passes.
**Validation:** User confirmed
## Open Questions
- [x] Should the skill's installer ship the `checkout-lease` session hook, which exists only in `~/.claude/settings.json`? → No, leave it local. Nothing is captured for it.
- [x] Should the notes in the local memory folder get the same treatment? → Yes, as a separate review request: REQ-725.
- [~] Where does G5 live if the skill has no natural home for "prompts for another session"? → Builder decides. If no existing file fits, record that in the REQ and leave G5 out instead of creating a new file.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the budget, and its row is `slugged: partial` so no targeted entry is legal. Matched on "budgeted context routing": this REQ adds text to an always-on file.
## Full Context
See `do-work/user-requests/UR-175/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-move-global-rules-into-skill.md*
