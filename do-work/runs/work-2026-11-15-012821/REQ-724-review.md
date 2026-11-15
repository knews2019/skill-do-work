# Review: REQ-724 (move the maintainer's global standing rules into the skill)

**Approve** — the rules now ship in plain, machine-free words. Two small edits should land before the release: F1 (sibling-package copies) and F2 (one misstated pointer).
Route B | merge range 60c8a8f8..47565ced (builder commit ab80e3a2)

## What's built

- Rule groups G1, G2, G3, G5 and the bug-report half of G4 are new Core Patterns bullets in the always-on `skills/do-work/crew-members/communication-style.md` (lines 24-28).
- The G2 option shape is in `clear-questions.md` principle 5. The numeric-target and redesign lines of G4 are in `coding-guardrails.md` § 1 and § 3. G7 is one paragraph in `background-agents.md`.
- Missing: the toolbox and knowledge packages keep their own copies of three of these files, and those copies did not get the new lines (F1).

## Decisions / risks for you

- F1 asks for the three additions to be copied into the sibling-package copies. Value: toolbox `code-review` and `deep-explore`, the main lookup fan-outs, then see the model rule. Risk: the REQ's GREEN check "exactly one shipped file" then reads one file per package, so the integrator must note that reading. Reversible.
- F6: before the global file is deleted, decide whether the two G4 lines and G7 may stay out of plain sessions. Today they load only during implementation or a fan-out.

## Findings

**Important:**

- F1 `skills/do-work-toolbox/crew-members/clear-questions.md:29`, `skills/do-work-knowledge/crew-members/clear-questions.md:29`, `skills/do-work-toolbox/crew-members/coding-guardrails.md:21` and `:85`, `skills/do-work-toolbox/crew-members/background-agents.md:189`, `skills/do-work-knowledge/crew-members/background-agents.md:189` — the sibling-package copies of the changed crew-member files did not get the new lines. These copies are kept in lock-step by hand: commits bf9d9b8a, f06c6082 and d1bf4855 changed every copy together. Toolbox actions load their own copy: `code-review.md`, `deep-explore.md`, `note.md`, `tidy-repo.md` load `background-agents.md`, and `folder-structure.md`, `ai-report.md`, `prime.md` and others load `clear-questions.md`. So the G7 model rule never reaches the toolbox fan-outs that run the most lookup agents. `_dev/primes/prime-action-files.md` § Traps `[family: alternate-writer-contract-drift]` names sibling-package restatements as the thing to grep. Suggested fix: append the same sentence to principle 5 of both `clear-questions.md` copies, the two bullets to the toolbox `coding-guardrails.md` (after `:21` and after `:85`), and the model paragraph after the `flow.` line (`:189`) of both `background-agents.md` copies. Keep each copy's existing `../../do-work/` path spellings. Record in the REQ that the GREEN searches for `cheapest model tier` and `unless the user says strict` now read one file per package. (Pre-existing drift in the toolbox `coding-guardrails.md`, the missing "Earned defense" paragraph, is out of scope.) — impact-rule-change → required fix in this run, not queued.
- F2 `skills/do-work/crew-members/communication-style.md:25` — "A reversible choice is the DECIDE & STATE tier" drops the source word "trivia" and misstates the canonical gate. `coding-guardrails.md:33-41` puts a choice in DECIDE & STATE only when it is reversible, its intent is inferable from context, and no reasonable person would disagree. A reversible choice that depends on taste or that people would dispute is ESCALATE (b) or (c). Read as written, the always-on rule tells every consumer session to decide those choices silently. Suggested fix text: `A reversible choice whose answer is clear from context is the DECIDE & STATE tier of \`crew-members/coding-guardrails.md\`: make it yourself and report it in one line.` — impact-rule-change → required fix in this run, not queued.

**Minor:**

- F3 `skills/do-work/crew-members/communication-style.md:25` and `skills/do-work/crew-members/background-agents.md:192` — the two new general rules do not say that an action's own rule wins. Two shipped texts already say otherwise. First, `skills/do-work-toolbox/actions/tutorial.md:355` says to print plain-text menus and not use the ask-user tool. Second, `skills/do-work/actions/ai-evaluation-loop-reference.md:19-20` names the reviewer models for `sol-loop` and `council`, while the new text says review agents keep the session model. Suggested fix: end the G2 sentence with "unless the action in use says how to ask", and write "keep the session model unless the action names another model". — impact-rule-change → report only
- F4 `skills/do-work/crew-members/communication-style.md:28` — the last sentence of the G5 bullet is about 38 words, against the 20-word target. The stall-check gloss also drops half of the rule: the UR's own prompt defines it as "If nothing moves for 20 minutes, stop and report what blocked". The source's `<date>-<slug>.md` file name is also dropped. Suggested fix text: `The prompt must run unattended. It records open questions as pending answers instead of waiting, and it leaves out steps only a person can do. It includes a 20-minute stall check: if nothing moves for 20 minutes, the run stops and reports what blocked it.` — impact-negligible → report only
- F5 REQ § Context row G6 — "nothing to move" is only partly verified. Exact-path staging ships (`skills/do-work/actions/commit.md:222`). The REQ/UR id re-read is covered by the `capture-files` reservations, which are created with exclusive-create (`skills/do-work/actions/capture.md:82`). But "never revert files you did not change" has no shipped home, and `git commit -a` is not named. The maintainer's capture answer covered only the session hook. Suggested fix, if wanted: one bullet in `commit.md` § Rules, "Never revert, stash or discard a change you did not make." — impact-rule-change → report only
- F6 `skills/do-work/crew-members/coding-guardrails.md:22` and `:89` — the numeric-target and redesign lines now load only during implementation (the file's JIT_CONTEXT says Step 6). G7 loads only during a fan-out (D-05, confirmed by the maintainer). The REQ allowed these homes, because only G1-G3 had to be always on. But in a plain session, which is where the global file applied them, these three rules stop applying once `~/.claude/CLAUDE.md` is deleted. This is not drift. It is a decision to take before the "delete after" step. — impact-rule-change → report only

## Requirements Checklist

- [x] Each source line checked against shipped text first. G1 "close with the next command" is covered by `communication-style.md:12` ("make the last line the next action: the command to run"). Verified.
- [x] G1 big picture first and status words defined — delivered (`communication-style.md:24`).
- [~] G2 ask-user tool, no prose list, decide reversible trivia — delivered, but the trivia pointer misstates the gate (F2).
- [x] G2 option shape in `clear-questions.md` — delivered in `do-work`. The sibling copies are missing it (F1).
- [x] G3 login shell, remote host shell, no placeholder, script over about 5 lines — delivered (`:26`).
- [x] G4 bug report (`:27`), numeric targets (`coding-guardrails.md:22`), redesign UI (`:89`) — delivered. The operator-work rule is verified as shipped (0.305.85, `actions/work.md:142`).
- [x] G5 points at `folder-structure.md`, which lists `run-prompts/` at line 38, and writes no path of its own — delivered. Wording gaps are in F4.
- [x] G6 — hook stays local by the maintainer's answer. Exact-path staging ships. One line has no home (F5).
- [x] G7 in `background-agents.md` — delivered in `do-work`. The sibling copies are missing it (F1).
- [x] No machine, vendor, tool or host name. `zsh|AskUserQuestion|haiku|EnterWorktree|Explore` are absent from the four files.
- [x] No shipped sentence reworded. The only deleted line, principle 5, keeps its two sentences byte-identical and only gains one appended sentence.
- [x] Aliases and Examples unchanged. The diff touches only `communication-style.md` lines 24-28.
- [x] About 25 lines (approximate): 13 added and 1 joined. No new file and no test file.
- [n/a] Release per `prime-releases.md` — the integrator writes it at finalization (D-12, minor).

## Acceptance Testing

**Result: Pass** (implementation and integration stages)

- `bash _dev/tests/shipped-package-reference-contract.sh` → `shipped package reference contract: PASS` (1.2 s). The `../../do-work-toolbox/actions/folder-structure.md` citation resolves in the repo (`skills/do-work-toolbox/actions/folder-structure.md`) and in the installed topology (`.claude/skills/do-work-toolbox/...` per `suite/modules.tsv`). The `crew-members/...` pointers are same-package local paths, as `prime-action-files.md` § Cross-Referencing requires.
- `bash do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` → 14 PASS lines, `REQ-724 probe PASS`, exit 0.
- I did not run the full gate. The integrator ran it green (135 s) and drains the staged-skills and installer lanes.

## Suggested Additional Testing

- Deployment: unassessed. The integrator's `installer` heavy lane (`_dev/tests/install-suite-behavior.sh`) covers the link from a consumer project's agent instructions.
- Live acceptance: unassessed. In a consumer project, open a plain session, ask "what is the state of the queue?", and confirm the reply gives a one-paragraph big picture and defines each status word.
- Live acceptance: unassessed. Ask a consumer session for a command to run on a remote host, and confirm it is wrapped for that host's shell and has no placeholder.
- After F1: grep `cheapest model tier` in `skills/` and confirm one hit per package that ships `background-agents.md`.

## Scores (on the record — not the headline)

**Overall: 88%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 85% | All groups in `do-work`. Sibling copies missing (F1). G2 trivia pointer misstated (F2). |
| Code Quality | 85% | Plain words and conditions. F2 misstatement, F3 precedence gap, F4 long sentence. |
| Test Adequacy | 90% | Prose only, and the REQ forbids a wording test. Search-based red-green proof plus the reference contract. |
| Scope | 95% | Exactly the four declared files. Lock-step copies not touched (F1). |
| Risk | Low | Prose rules. F2 could cause silent decisions on taste choices. |
| Acceptance | Pass | Probe and reference contract pass. |

## Follow-ups created

None (6 findings report only, F1 and F2 marked as required fixes in this run)

- Open item, already parked in `_dev/run-prompts/2026-11-15-run-prompts-folder-rule.md`: in this repo, the always-on G5 bullet points through `folder-structure.md` to `docs/run-prompts/`. Maintainer practice here is `_dev/run-prompts/`, because root `docs/` breaks the staged-skills lane. This is not a regression, because the global file named `docs/run-prompts/` too.

## Review

**Overall: 88%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 85% |
| Code Quality | 85% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 sibling-package copies of `clear-questions.md` (toolbox, knowledge), `coding-guardrails.md` (toolbox) and `background-agents.md` (toolbox, knowledge) lack the new lines, which breaks the hand-kept lock-step; toolbox fan-outs never see G7 — impact-rule-change → required fix in this run, not queued
- F2 `communication-style.md:25` "A reversible choice is the DECIDE & STATE tier" misstates the canonical three-condition gate (`coding-guardrails.md:33-41`); fix: "A reversible choice whose answer is clear from context is …" — impact-rule-change → required fix in this run, not queued

**Minor findings:**
- F3 the new G2 rule (`communication-style.md:25`) and G7 rule (`background-agents.md:192`) state no precedence for an action's own rule (`do-work-toolbox/actions/tutorial.md:355`, `actions/ai-evaluation-loop-reference.md:19-20`) — impact-rule-change → report only
- F4 `communication-style.md:28` last sentence about 38 words; stall-check gloss drops "stop and report what blocked"; `<date>-<slug>.md` file name dropped — impact-negligible → report only
- F5 G6 "nothing to move" partly unverified: "never revert files you did not change" has no shipped home — impact-rule-change → report only
- F6 G4 numeric/redesign lines (`coding-guardrails.md:22`, `:89`) and G7 stop applying in plain sessions once the global file is deleted; decide before the delete step — impact-rule-change → report only

**Acceptance:** Pass — implementation and integration stages: `shipped-package-reference-contract.sh` PASS, `REQ-724-probe.sh` 14 PASS exit 0; deployment and live acceptance unassessed
**Restatement sweep:** redefined the DECIDE & STATE pointer (stale against `coding-guardrails.md:33-41`, F2), the clear-questions option shape (sibling copies stale, F1; consistent with `capture-reference.md:137` and `work-reference.md:605-613`), the clear-questions load trigger (consistent with `SKILL.md:68` and the `clear-questions.md` JIT_CONTEXT), the run-prompts folder (single home `folder-structure.md:38`, consistent), and the per-agent model rule (sibling copies stale, F1; precedence gap with `ai-evaluation-loop-reference.md`, F3)
**Anti-bloat count:** 0 (no helper, option, file or test the REQ did not name)
**Suggested testing:** 4 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*
