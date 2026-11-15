# REQ-724 Hand-back (move the maintainer's global standing rules into the skill's crew-member files)

Branch `worktree-agent-REQ-724-move-standing-rules-into-skill`, base `d2f8cefc`, tip `ab80e3a2`.

## Commits
- `ab80e3a2` [REQ-724] move the maintainer's standing rules into the crew-member files

## Files
- `skills/do-work/crew-members/communication-style.md` — modified: five bullets appended at the end of `## Core Patterns` (G1, G2, G3, G4 bug-report line, G5).
- `skills/do-work/crew-members/clear-questions.md` — modified: one sentence appended to principle `### 5. Concrete options, never open-ended` (G2 option shape).
- `skills/do-work/crew-members/coding-guardrails.md` — modified: one bullet in `## 1. Think Before Coding` ("Before implementing" list) and one in `## 3. Surgical Changes` ("When editing existing code" list).
- `skills/do-work/crew-members/background-agents.md` — modified: one paragraph at the end of `## Match the Pattern to the Harness`, before `## Manifest Format` (G7).

## Scope changes
None. Only the four planned files changed.

## Red-green evidence
RED at base `d2f8cefc` (later than `1389e2be`), six searches, all empty:
```
every status word -> []
prose list of questions -> []
login shell -> []
unless the user says strict -> []
20-minute stall check -> []
cheapest model tier -> []
```
GREEN probe at base: exit 1, `REQ-724 probe FAIL (12)` (12 FAIL lines, as the brief expected).

GREEN probe at tip `ab80e3a2`:
```
PASS: G1 found once in skills/do-work/crew-members/communication-style.md
PASS: G1 found once in skills/do-work/crew-members/communication-style.md
PASS: G2 found once in skills/do-work/crew-members/communication-style.md
PASS: G3 found once in skills/do-work/crew-members/communication-style.md
PASS: G3 found once in skills/do-work/crew-members/communication-style.md
PASS: G4 found once in skills/do-work/crew-members/communication-style.md
PASS: G4 found once in skills/do-work/crew-members/coding-guardrails.md
PASS: G4 found once in skills/do-work/crew-members/coding-guardrails.md
PASS: G5 found once in skills/do-work/crew-members/communication-style.md
PASS: G7 found once in skills/do-work/crew-members/background-agents.md
PASS: G2 option shape in clear-questions.md
PASS: G5 cites folder-structure.md
PASS: no machine or tool names in the four crew-member files
PASS: shipped-package-reference-contract
REQ-724 probe PASS
```

## Verification
- `bash do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` (worktree as cwd): `REQ-724 probe PASS`, rerun after the final wording change, PASS again.
- `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, `shipped package reference contract: PASS` (1.1 s), PASS again after the final wording change.
- `bash _dev/tests/contract-regressions.sh`: exit 0, `Contract regression checks passed.` (24 s, one run, no rerun needed).
- `git diff --stat d2f8cefc..HEAD`: 4 files, 13 insertions, 1 deletion. The one deletion is principle 5 of `clear-questions.md`, where the new sentence joins the existing paragraph line.
- Patterns to Avoid re-read: no stock emphasis phrase, no analogy, no semicolon or em-dash chain in added text, no vendor, tool, host or machine name.

## P-A-U
- PLAN: add only the missing rules, each once, in the four planned files, at the end of the named lists or sections, in the house style of `communication-style.md`.
- APPLY: one commit, 13 added lines and 1 joined line, prose only.
- UNIFY: `git diff --stat` shows only the four crew-member files. I checked that no existing sentence is reworded, Aliases and Examples are unchanged, the cross-package citation resolves, and `zsh|AskUserQuestion|haiku|EnterWorktree` is absent.

## Decisions
- **D-05 (confirmed by the maintainer):** the model-per-agent rule (G7) lives only in `crew-members/background-agents.md`. No copy or pointer line in `communication-style.md`, per the Coordinator note of 2026-11-15.
- **D-07 (DECIDE & STATE):** none of the eight rules was already shipped. I searched `skills/` for the ideas (big picture, recommended option, placeholder, paste-ready, model tier, stall check, pending answer). The only near hit is the `ttsr` alias appendix ("paste-ready commands"), which covers one answer format, not commands in general. So every rule was added.
- **D-08 (DECIDE & STATE):** "parks questions as pending answers" is written as "records open questions as pending answers instead of waiting", and "operator steps" as "steps only a person can do". "Stall check" gets a plain gloss in brackets. Reason: the house style asks for plain words, and "parks" and "operator" are jargon for a first-time reader.
- **D-09 (DECIDE & STATE):** the G5 bullet's long first sentence was split in two, so each sentence stays near 20 words.
- **D-10 (DECIDE & STATE):** the numeric-target rule says "a number in the request (a size, a count, a time)". The bracket is marked illustrative by "a", it is not a closed list.
- **D-11 (DECIDE & STATE):** the G4 bug-report line is a plain bullet without a bold lead, because it is one short sentence, like the other short bullets in the list.

## Discovered Tasks
None.

## Integration seams
- `communication-style.md` → `## Core Patterns` (five bullets after the last existing bullet). This file is also imported by this repository's `CLAUDE.md`, so maintainer sessions get the rules at once.
- `clear-questions.md` → `### 5. Concrete options, never open-ended`.
- `coding-guardrails.md` → `## 1. Think Before Coding` and `## 3. Surgical Changes`.
- `background-agents.md` → `## Match the Pattern to the Harness`.
- REQ-732 (allowing `docs/run-prompts/` in the staged-skills contract) does not touch these files. The G5 bullet names no folder path.

## Proposed CHANGELOG entry
Title: **Standing Rules for Status Explanations, User Decisions, Pasted Commands and Subagent Models Ship in the Crew Files**

What changes for an installed copy: every session in an installed project now follows the new always-on rules in `communication-style.md`. They are: big picture first with status words defined, decisions through the ask-user tool, paste-ready commands, a bug report means fix the code, and a prompt for another session is a committed file with a short paste block. Question wording, implementation guardrails and fan-out model choice gain one rule each.

Bump: **minor**. Reason: an always-on file that every installed session loads gains new required behaviour, so agent behaviour changes in every consumer project, not only a wording fix. If the integrator reads `work-reference.md` → Changelog Entry Procedure as "prose-only rule additions are patch", patch is acceptable. No code, schema or command changes.

## Proposed lesson
None.

## Lessons read
- `_dev/primes/lessons-releases.md`: its families (`canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`) do not fire. No file moved and no manifest changed. The one new relative link was checked by the reference contract in both topologies.
