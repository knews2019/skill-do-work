---
id: REQ-724
title: '[impact-rule-change] Move the maintainer''s global standing rules into the skill'
status: completed
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
write_set: [skills/do-work/crew-members/communication-style.md, skills/do-work/crew-members/clear-questions.md, skills/do-work/crew-members/coding-guardrails.md, skills/do-work/crew-members/background-agents.md]
claimed_at: 2026-11-14T23:28:11Z
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-14T23:31:57Z
  basis:
    - Route B
    - 4-file write set
    - 7 acceptance criteria
dispatch_at: 2026-11-15T00:01:58Z
builder_handback_at: 2026-11-15T02:09:01Z
integration_at: 2026-11-15T02:09:32Z
kb_status: pending
review_at: 2026-11-15T02:26:21Z
heavy_verified_at: 2026-11-15T02:31:30Z
heavy_verified_revision: 0074c2744688f020c37b352325a96a1cbb01b605
commit: 0074c2744688f020c37b352325a96a1cbb01b605
completed_at: 2026-11-15T02:32:03Z
release_at: 2026-11-15T02:32:03Z
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
- [~] Where does G5 live if the skill has no natural home for "prompts for another session"? → Builder decides. If no existing file fits, record that in the REQ and leave G5 out instead of creating a new file. → **D-01**: Builder chose: `skills/do-work/crew-members/communication-style.md`, one Core Patterns bullet that points at the run-prompts folder in `skills/do-work-toolbox/actions/folder-structure.md` and names no path of its own. Reasoning: the only shipped handoff, `actions/restart-with-parallel-handoff.md`, writes the queue-specific `do-work/RESTART-PROMPT.md` and runs only on `phandoff` or a context limit; a request for a prompt to run in another session or tool arrives in a plain session, so only the always-on file sees it. Value: G5 applies on every machine with no new file. Risk: one more always-on line; fully reversible by deleting it.
<!-- D-XX counter: last used D-15. Next decision: D-16. -->
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the budget, and its row is `slugged: partial` so no targeted entry is legal. Matched on "budgeted context routing": this REQ adds text to an always-on file.
## Full Context
See `do-work/user-requests/UR-175/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Add only the missing rules, each once, in the four planned files, at the end of the named lists or sections, in the house style of `communication-style.md`.
- [x] **[APPLY]:** (from the builder hand-back) One commit (`ab80e3a2`), 13 added lines and 1 joined line, prose only.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat` shows only the four crew-member files. No existing sentence is reworded, Aliases and Examples are unchanged, the cross-package citation resolves, and `zsh|AskUserQuestion|haiku|EnterWorktree` is absent. Integrator re-check after the merge: see Qualification.
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-move-global-rules-into-skill.md*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and the four target files are named, but every source line must first be checked against shipped text at HEAD, and the exact home of G4, G5 and G7 needs discovery. Prose only, one subsystem (crew-member files), no Go or shell change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Pre-dispatch agent, base `1389e2be`. Every source line was re-checked against `skills/` at HEAD (fixed-string `grep -rliF`, `CHANGELOG.md` excluded).

**Required-lessons consult.** `required_lessons` stays `[_dev/primes/lessons-releases.md]` (856 tokens, read). Its families (`canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`) do not fire: nothing moves and no manifest changes. `lessons-action-files.md` stays dropped for budget (8761 tokens, `slugged: partial`). Its `[family: budgeted-context-routing]` bullets (REQ-478, REQ-479) are about the required-lessons pointer, not about always-on prose, so the drop costs nothing here. No other index row matches.

**State per rule group at HEAD:**
- G1. `skills/do-work/crew-members/communication-style.md:12` already says "Answer first, next action last" (the closing next command). `:17` already says "When you must use a technical term, explain it in plain words the first time" — status words are the case agents skip because they look like plain English. Missing: big picture first, and the status-word clause. Search `big picture in one short paragraph` and `every status word`: no hit.
- G2. `skills/do-work/crew-members/clear-questions.md:3` loads before any interactive question; principles 3 (`:19-21`, consequence per option) and 5 (`:27-29`, concrete options) exist. Missing: ask-user tool over a prose list, recommended option first, value and risk per option. "Decide reversible trivia yourself" is already the DECIDE & STATE tier, `skills/do-work/crew-members/coding-guardrails.md:32-35`, but that file loads only in Step 6. Search `prose list of questions`: no hit.
- G3. No hit for `login shell`, `one-line launcher`, or a no-placeholder rule.
- G4. Operator-work rule shipped (0.305.85). No hit for `unless the user says strict`, `A bug report means`, `while redesigning`. `coding-guardrails.md` § 3 Surgical Changes (`:78-94`) is the natural home for the UI line; § 1 Think Before Coding (`:12-24`) for the numeric-target line.
- G5. Only the folder: `skills/do-work-toolbox/actions/folder-structure.md:38` (`docs/` holds `run-prompts/`). The only shipped handoff, `skills/do-work/actions/restart-with-parallel-handoff.md`, writes `do-work/RESTART-PROMPT.md` on `phandoff` or a context limit. That is a different artifact. `stall check` appears only in `actions/work-finalization-reference.md` (a different subject). No hit for `20-minute stall check`.
- G6. Nothing to move (exact-path commits shipped; hook stays local).
- G7. `skills/do-work/crew-members/background-agents.md` § Match the Pattern to the Harness (`:135-188`, ends before `## Manifest Format` at `:190`). No hit for `cheapest model tier`.
- Negative check: `zsh|AskUserQuestion|haiku|EnterWorktree` in the four files: no hit at base.

**Cross-package citation.** From `crew-members/` the toolbox file is `../../do-work-toolbox/actions/folder-structure.md` (`_dev/primes/prime-action-files.md` § Cross-Referencing); `_dev/tests/shipped-package-reference-contract.sh` checks it. Existing precedent: `crew-members/general.md:15`.

**RED searches (recorded at `1389e2be`, all empty):**
1. G1: grep -rliF --exclude=CHANGELOG.md 'every status word' skills
2. G2: grep -rliF --exclude=CHANGELOG.md 'prose list of questions' skills
3. G3: grep -rliF --exclude=CHANGELOG.md 'login shell' skills
4. G4: grep -rliF --exclude=CHANGELOG.md 'unless the user says strict' skills
5. G5: grep -rliF --exclude=CHANGELOG.md '20-minute stall check' skills
6. G7: grep -rliF --exclude=CHANGELOG.md 'cheapest model tier' skills

**Pre-dispatch decisions** (D-01 is in Open Questions):
- **D-02** (DECIDE & STATE): the G4 bug-report line goes to `communication-style.md`, not `coding-guardrails.md`. A bug report arrives in a plain session before any REQ exists, and `coding-guardrails.md` loads only in Step 6. The numeric-target and redesign lines stay in `coding-guardrails.md` as suggested.
- **D-03** (DECIDE & STATE): the status-word rule is a clause in the new G1 bullet. The shipped technical-term line (`communication-style.md:17`) is not reworded.
- **D-04** (DECIDE & STATE): the "decide reversible trivia yourself" half of G2 is a pointer to the DECIDE & STATE tier in `coding-guardrails.md`, plus "report it in one line". It does not restate the tier.
- **D-05** (ESCALATE): G7 goes to `background-agents.md` as suggested. Value: the rule sits where agents are dispatched, and the always-on file stays about talking. Risk: that file loads only for fan-out, so a single lookup agent outside a fan-out never sees the rule. Reversible: move the paragraph to `communication-style.md`.
- **D-06** (DECIDE & STATE): the proof phrases are fixed (the six RED searches above) and exclude `CHANGELOG.md`, because the changelog mirror under `skills/do-work/` may quote a phrase. Per the run rules the builder never edits `CHANGELOG.md` or `VERSION`; the integrator releases, so they leave Scope even though capture's `write_set` named `skills/do-work/CHANGELOG.md`.

*Generated by pre-dispatch agent (in-context exploration)*

## Scope

**Files I will touch:**
- `skills/do-work/crew-members/communication-style.md` (modify) — new Core Patterns bullets for G1 (big picture first, every status word defined), G2 (ask-user tool, no prose list of questions, decide reversible choices per the coding-guardrails gate), G3 (commands for the login shell, no placeholder, longer logic in a committed script with a one-line launcher), the G4 bug-report line (D-02) and G5 (D-01)
- `skills/do-work/crew-members/clear-questions.md` (modify) — principle 5 gains the option shape: recommended option first, value and risk per option
- `skills/do-work/crew-members/coding-guardrails.md` (modify) — numeric targets approximate unless the user says strict (§ 1), never remove existing UI while redesigning unless asked (§ 3)
- `skills/do-work/crew-members/background-agents.md` (modify) — one paragraph: model per agent, cheapest model tier for lookups (D-05)

**Files I will NOT touch:** `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` and the version mirrors (the integrator releases, D-06); the Aliases and Examples sections of `communication-style.md`; any shipped rule's existing wording; `skills/do-work/actions/restart-with-parallel-handoff.md` and `skills/do-work-toolbox/actions/folder-structure.md` (cited, not changed); no new crew-member file, no test file, nothing outside the repository.

**Acceptance criteria (restated from REQ):**
- [ ] Each rule group G1, G2, G3, G4, G5 and G7 is shipped once; a line is added only where no shipped text already says it, and no shipped rule is reworded.
- [ ] G1, G2 and G3 live in `communication-style.md` (always on), with the option shape of G2 in `clear-questions.md`.
- [ ] Shipped text names no machine, vendor, tool or host: the login shell, the agent's ask-user tool when it has one, the cheapest model tier that can do the lookup.
- [ ] G5 points at the run-prompts folder of `folder-structure.md` and writes no path of its own.
- [ ] The six RED searches each find exactly one shipped file (CHANGELOG excluded), and `zsh`, `AskUserQuestion`, `haiku`, `EnterWorktree` are absent from the four changed files.
- [ ] About 25 shipped lines in total (approximate); Aliases and Examples unchanged; no new file.
- [ ] The release follows `_dev/primes/prime-releases.md` (proposed by the builder, written by the integrator).

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-724-preflight-probe.sh) and the repository gate green at the REQ-729 merge 572e8b94 (141 s)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*

## Decisions
- **D-05 (confirmed by the maintainer, 2026-11-15):** the model-per-agent rule (G7) lives only in `crew-members/background-agents.md`, with no copy or pointer line in `communication-style.md`. The maintainer confirmed this through a question on 2026-11-15, so the pre-dispatch ESCALATE is closed and no pending-answers follow-up is written. (from the builder hand-back and the integrator brief)
- **D-07 (DECIDE & STATE, from the builder hand-back):** none of the eight rules was already shipped. The builder searched `skills/` for the ideas (big picture, recommended option, placeholder, paste-ready, model tier, stall check, pending answer). The only near hit is the `ttsr` alias appendix ("paste-ready commands"), which covers one answer format, not commands in general. So every rule was added.
- **D-08 (DECIDE & STATE, from the builder hand-back):** "parks questions as pending answers" is written as "records open questions as pending answers instead of waiting", and "operator steps" as "steps only a person can do". "Stall check" gets a plain gloss in brackets. Reason: the house style asks for plain words, and "parks" and "operator" are jargon for a first-time reader.
- **D-09 (DECIDE & STATE, from the builder hand-back):** the G5 bullet's long first sentence was split in two, so each sentence stays near 20 words.
- **D-10 (DECIDE & STATE, from the builder hand-back):** the numeric-target rule says "a number in the request (a size, a count, a time)". The bracket is marked illustrative by "a", it is not a closed list.
- **D-11 (DECIDE & STATE, from the builder hand-back):** the G4 bug-report line is a plain bullet without a bold lead, because it is one short sentence, like the other short bullets in the list.
- **D-12 (DECIDE & STATE, integrator):** release as a minor bump. Reason: `communication-style.md` is always on in every installed project, so the new rules change agent behavior in every consumer project (impact-rule-change). The builder proposed minor too. Value: the version tells consumers that behavior changed. Risk: none, a version number only.
- **D-13 (DECIDE & STATE, integrator):** drain the `installer` heavy lane in addition to the one selected lane (`staged-skills`). Reason: `_dev/tests/install-suite-behavior.sh` (installer lane) is the contract test that names `communication-style.md` (it checks that the installer links the file from a consumer project's agent instructions), and the brief asks that the tests naming that file still pass. Value: direct evidence for the always-on link. Risk: about 30 s more wall time.
- **D-14 (DECIDE & STATE, integrator, review F1):** copy the added lines into the lock-step copies of `clear-questions.md` (toolbox, knowledge), `coding-guardrails.md` (toolbox) and `background-agents.md` (toolbox, knowledge). Reason: each package ships its own copy, earlier rule changes (bf9d9b8a, f06c6082, d1bf4855) changed every copy together, and toolbox fan-outs (`code-review`, `deep-explore`) load their own `background-agents.md`, so G7 would never reach them. A lock-step copy of the same file is the same rule in each package, not a restatement in a second file. The GREEN probe now accepts one home per package for `unless the user says strict`, `while redesigning` and `cheapest model tier`. Value: the rules reach every installed package. Risk: five more files to keep in step by hand, as today; reversible.
- **D-15 (DECIDE & STATE, integrator, review F2 and F4):** the DECIDE & STATE pointer now carries the tier's "clear from context" condition (F2, required). The G5 unattended sentence is split in three and the stall check says what it does: "if nothing moves for 20 minutes, the run stops and reports what blocked it", the wording of the UR's own capture prompt (F4, report-only finding, applied because the re-merge was happening anyway and the gloss had lost half of the rule). The `<date>-<slug>.md` file name stays out: `folder-structure.md` owns the folder and a second naming rule would be a second path rule.

## Implementation Summary

**Files changed:**
- `skills/do-work/crew-members/communication-style.md` (modified) — five bullets appended at the end of the Core Patterns section: G1 big picture first and status words defined, G2 decisions through the ask-user tool, G3 commands that run as pasted, the G4 bug-report line, G5 a prompt for another session is a committed file.
- `skills/do-work/crew-members/clear-questions.md` (modified) — principle 5 gains one sentence: recommended option first, value and risk per option.
- `skills/do-work/crew-members/coding-guardrails.md` (modified) — one bullet in § 1 Think Before Coding (numbers are approximate unless the user says strict) and one in § 3 Surgical Changes (never remove existing UI while redesigning unless asked).
- `skills/do-work/crew-members/background-agents.md` (modified) — one paragraph at the end of § Match the Pattern to the Harness: pick each agent's model by its work, cheapest model tier for lookups.

- `skills/do-work-toolbox/crew-members/clear-questions.md`, `skills/do-work-knowledge/crew-members/clear-questions.md` (modified, at review F1) — the same principle 5 sentence.
- `skills/do-work-toolbox/crew-members/coding-guardrails.md` (modified, at review F1) — the same two bullets.
- `skills/do-work-toolbox/crew-members/background-agents.md`, `skills/do-work-knowledge/crew-members/background-agents.md` (modified, at review F1) — the same model paragraph.

**What was done:** The rule groups G1, G2, G3, G4 (two of the three missing lines in `coding-guardrails.md`, the bug-report line in `communication-style.md`), G5 and G7 from the maintainer's machine-local instruction file now ship in the crew-member files, each once, written as conditions with no machine, vendor or tool name. Hand-back merge range 60c8a8f8..47565ced, no conflicts, no integration seams. Review fixes (F1, F2, F4) were committed on the builder branch as `e8bbcd7c` and re-merged with the same base: cumulative range 60c8a8f8..0074c274 (it also holds the run-artifact commit `0bd17956` for the probe, under `do-work/` only).

## Discovered Tasks

(from the builder hand-back) None.

## Qualification

**Gate records:** `advance --diff-range 60c8a8f8..47565ced` (qualify) first refused `QUALIFY-CLAIMED-PATH-MISSING` and `SCOPE-UNDECLARED-TOUCH` on a backticked section heading in the Implementation Summary, which the parser read as a path. I reworded that one mention to plain text; the rerun returned `success` with no findings: no scope drift, no debug artifacts (the diff adds prose only, no TODO, FIXME or stray marker).

**Requirement trace (against `git diff 60c8a8f8..47565ced`, 4 files, +13/-1):**
- Each rule group G1, G2, G3, G4, G5 and G7 is shipped once, added only where nothing said it: the six RED searches (CHANGELOG excluded) each print exactly one file at the merge: `every status word`, `prose list of questions`, `login shell`, `20-minute stall check` → `communication-style.md`; `unless the user says strict` → `coding-guardrails.md`; `cheapest model tier` → `background-agents.md`. G6 has nothing to move. ✓
- G1, G2 and G3 are always on: three bold-lead bullets in the Core Patterns section of `communication-style.md`; the G2 option shape (recommended option first, value and risk per option) is one sentence in principle 5 of `clear-questions.md`; "decide reversible trivia yourself" points at the DECIDE & STATE tier (`coding-guardrails.md:33`) and adds only "report it in one line". ✓
- No machine, vendor, tool or host name: the text says "the user's login shell", "the agent's ask-user tool when it has one", "the cheapest model tier that can do them"; `zsh|AskUserQuestion|haiku|EnterWorktree` is absent from the four files (probe). ✓
- G5 names no path of its own: it cites `../../do-work-toolbox/actions/folder-structure.md`, whose tree at line 38 lists `run-prompts/` under `docs/`; `docs/run-prompts` does not appear in `communication-style.md`; the shipped-package reference contract passes in both topologies. ✓
- About 25 lines target (approximate): 13 added lines and 1 joined line. Aliases and Examples are unchanged (the diff touches only the end of Core Patterns). No new file, no test file. ✓
- No shipped rule reworded: the only deleted line is principle 5 of `clear-questions.md`, whose existing two sentences stay byte-identical with one sentence appended. ✓
- Release per `_dev/primes/prime-releases.md`: written at finalization (minor, D-12). Pending until Step 9.

**Scope comparison (Route B):** Scope declares the four crew-member files; the merge touched exactly those four. No undeclared file. `CHANGELOG.md` and `VERSION` were left for the integrator, as Scope states (D-06).

**Re-qualification after the review fixes (range 60c8a8f8..0074c274):** the shipped diff is now 9 files, +27/-3. The five added files are the lock-step sibling-package copies of three Scope files (D-14, review F1); they carry byte-identical added lines, and each copy keeps its own `../../do-work/` path spellings. The six RED searches now print one file per package that ships the file: `cheapest model tier` → the do-work, toolbox and knowledge `background-agents.md`; `unless the user says strict` → the do-work and toolbox `coding-guardrails.md`; the other four → `communication-style.md` only. `zsh|AskUserQuestion|haiku|EnterWorktree` is absent from all nine files.

## Testing

**Tests run:**
- Repository gate `bash _dev/tests/maintainer-verify.sh` via `helpers/gate.sh` at merge `47565ced69018a3adb24ad2d7ce3298b367cab04` (detached drain-head checkout, stage reuse off), started at 1-minute load 4.61 → `GATE-EXIT=0 WALL=135s`, `Maintainer verification passed.`
- GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` through `advance --probe-file` → probe status 0 (`REQ-724 probe PASS`, 14 PASS lines: ten single-home phrase checks, the G2 option shape, the G5 citation and no own path, no machine or tool names, the shipped-package reference contract); green gate recorded for the merge revision.
- The six RED searches re-run on the merged main checkout: each prints exactly one shipped file (see Qualification).
- After the review fixes (re-merge `0074c2744688f020c37b352325a96a1cbb01b605`): the repository gate again via `helpers/gate.sh`, started at 1-minute load 4.29 → `GATE-EXIT=0 WALL=127s`; the amended GREEN probe (one home per package, D-14) → `REQ-724 probe PASS`; the green gate recorded for that revision with `record-green-gate` (the review phase refuses gate input to `advance`).

**Result:** ✓ All passing

**Red-green validation:** `tdd: false` (prose rules; the REQ forbids a wording test). RED at builder base `d2f8cefc`: the six searches were empty and the GREEN probe exited 1 with 12 FAIL lines (from the hand-back). GREEN at the builder tip `ab80e3a2` and at the merge `47565ced` (re-run here).

**New tests added:** None (Constraints: no new test file).

**Heavy verification plan:** the first plan (range 60c8a8f8..47565ced) selected only `staged-skills`. The cumulative range after the review fixes holds the run-artifact commit `0bd17956`, whose `do-work/runs/` paths no lane covers, so the planner marks coverage uncertain and selects all six lanes. The final plan is in `## Heavy Verification Plan`.

*Verified by work action (integrator, run work-2026-11-15-012821)*

## Review

**Overall: 92%** | 2026-11-15T02:26:21Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 sibling-package copies lacked the new lines; fixed in this run (e8bbcd7c, D-14), byte-identical lines in the same position in all five copies, own `../../do-work/` spellings untouched — impact-rule-change → report only
- F2 `communication-style.md:25` DECIDE & STATE pointer misstated the gate; fixed in this run (e8bbcd7c, D-15), now "A reversible choice whose answer is clear from context" — impact-rule-change → report only

**Minor findings:**
- F3 the new G2 rule (`communication-style.md:25`) and G7 rule (`background-agents.md:192`, and now its two sibling copies) state no precedence for an action's own rule (`do-work-toolbox/actions/tutorial.md:355`, `actions/ai-evaluation-loop-reference.md:19-20`) — impact-rule-change → report only
- F4 `communication-style.md:28` long G5 sentence and half-dropped stall-check gloss; fixed in this run (e8bbcd7c, D-15), three sentences of 5, 20 and 20 words, stall check matches `UR-175/input.md:45` — impact-negligible → report only
- F5 G6 "nothing to move" partly unverified: "never revert files you did not change" has no shipped home — impact-rule-change → report only
- F6 G4 numeric/redesign lines (`coding-guardrails.md:22`, `:89`) and G7 stop applying in plain sessions once the global file is deleted; decide before the delete step — impact-rule-change → report only
- N1 probe G2 option-shape and machine-name checks read only the `do-work` copies; PASS message says "found once" for three files — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: `shipped-package-reference-contract.sh` PASS, `REQ-724-probe.sh` 14 PASS exit 0, probe negative test fails on a stray file and on a missing copy; the integrator's gate at 0074c274 `GATE-EXIT=0`; deployment and live acceptance unassessed
**Restatement sweep:** redefined the DECIDE & STATE pointer (now consistent with `coding-guardrails.md:32-35`), the clear-questions option shape (all three copies identical), the per-agent model rule (all three `background-agents.md` copies identical; F3 precedence gap remains), the numeric-target and redesign lines (both `coding-guardrails.md` copies identical), and the stall-check wording (consistent with `UR-175/input.md:45`)
**Anti-bloat count:** 0 (the five added files are lock-step copies of Scope files, not new surface)
**Suggested testing:** 4 items (carried from the first review; the F1 grep item is now done)
**Review history:** first review 88% Approve (`do-work/runs/work-2026-11-15-012821/REQ-724-review.md`, F1 and F2 marked required); fixes F1, F2, F4 re-merged at `0074c274`; delta re-review 92% Approve (`do-work/runs/work-2026-11-15-012821/REQ-724-review-delta.md`).
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: 60c8a8f8b88b5c30f6923fb33984cfa3f58bd551
- Target: 0074c2744688f020c37b352325a96a1cbb01b605
- Selection: the shipped paths (nine crew-member files) match only the `staged-skills` subtree; the run-artifact commit `0bd17956` adds `do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` and `REQ-724-progress.log`, which no lane covers, so the planner reports coverage uncertain and selects every lane.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh, do-work/runs/work-2026-11-15-012821/REQ-724-progress.log
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh, do-work/runs/work-2026-11-15-012821/REQ-724-progress.log
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh, do-work/runs/work-2026-11-15-012821/REQ-724-progress.log
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh, do-work/runs/work-2026-11-15-012821/REQ-724-progress.log
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh, do-work/runs/work-2026-11-15-012821/REQ-724-progress.log
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh, do-work/runs/work-2026-11-15-012821/REQ-724-progress.log

## Heavy Verification Result

- Target: 0074c2744688f020c37b352325a96a1cbb01b605
- Execution revision: 0074c2744688f020c37b352325a96a1cbb01b605 (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set, start load 2.79)
- queue-kanban-javascript: executed, exit 0, 7 s
- queue-kanban-browser: executed, exit 0, 78 s
- do-work-cli-integrations: executed, exit 0, 64 s
- staged-skills: executed, exit 0, 36 s
- updater: executed, exit 0, 66 s
- installer: executed, exit 0, 27 s (holds `install-suite-behavior.sh`, which checks the link to `communication-style.md`)

## Lessons Learned

**What worked:** Fixed proof phrases chosen at pre-dispatch (six searches, CHANGELOG excluded) made RED, GREEN and the review all check the same thing.
**What didn't:** The builder brief and the probe treated `skills/do-work/crew-members/` as the only home of each file, so the toolbox and knowledge lock-step copies were missed until review (F1). The probe's "exactly one shipped file" check then had to learn that a home is one file per package.
**Worth knowing:** `clear-questions.md`, `coding-guardrails.md` and `background-agents.md` have hand-kept copies in `do-work-toolbox` and `do-work-knowledge`; a rule added to one must be added to every copy (`_dev/primes/prime-action-files.md` § Traps, `[family: alternate-writer-contract-drift]`, already names this). The toolbox `coding-guardrails.md` already lacks the "Earned defense" paragraph, a pre-existing drift left out of scope. No new lesson bullet: the builder proposed none, and the existing family covers it.

## Orientation

Crew members: the always-on `communication-style.md` now carries the status, decision, pasted-command, bug-report and other-session-prompt rules; `clear-questions.md`, `coding-guardrails.md` and `background-agents.md` gained one rule each, in all their package copies.

## Timing

Observed 2026-11-15T00:01:58Z to 2026-11-15T02:31:30Z: 2h 29m 32s total, 2h 30m 11s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 2h 07m 04s | 1 |
| review | 12m 05s | 1 |
| verification-gate | 10m 14s | 3 |
| handback-merge | 48s | 1 |

Slowest stage: builder-work / builder worktree build, 2h 07m 04s, outcome success.
