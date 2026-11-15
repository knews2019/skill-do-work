# Re-review (delta): REQ-724 (move the maintainer's global standing rules into the skill)

**Approve** — the required fixes F1 and F2 landed as asked, and the applied F4 wording is accurate. The delta adds no new problem that needs a fix.
Delta range 47565ced..0074c274 (builder fix commit e8bbcd7c, probe commit 0bd17956) | whole REQ range 60c8a8f8..0074c274

## What the delta changed

- F1: the new lines are in the five sibling-package copies (toolbox and knowledge `clear-questions.md`, toolbox `coding-guardrails.md`, toolbox and knowledge `background-agents.md`).
- F2: the DECIDE & STATE pointer in `skills/do-work/crew-members/communication-style.md:25` now carries the "clear from context" condition.
- F4: the G5 unattended sentence in `communication-style.md:28` is split in three, and the stall check now says what it does.
- The probe `do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` now accepts one home per package for the three phrases that live in copied files.

## Checks

1. **F1 — verified.** `diff` of each sibling copy against the `do-work` copy shows only differences that existed before: the copies' own `../../do-work/` path spellings, two re-wrapped paragraphs in `background-agents.md`, and the toolbox `coding-guardrails.md` gap (missing "Earned defense" paragraph, out of scope per the first review). The added lines are byte-identical and sit in the same position: principle 5 end (`clear-questions.md:29`), after the "If something is unclear" bullet and after the "unrelated dead code" bullet (`coding-guardrails.md`), and after the `flow.` line (`background-agents.md`). The delta changes no existing `../../do-work/` spelling: the only changed old line in each `clear-questions.md` copy keeps `../../do-work/actions/capture.md` and gains one appended sentence. `bash _dev/tests/shipped-package-reference-contract.sh` → `shipped package reference contract: PASS`, exit 0.
2. **F2 — verified.** `coding-guardrails.md:32-35` defines DECIDE & STATE as reversible, intent inferable from context, and no reasonable person would disagree. "A reversible choice whose answer is clear from context" names the first two conditions, and "answer is clear" covers the third closely enough for a pointer. The pointer sends the reader to the canonical tier, so the full gate stays one click away. This is the exact text the first review suggested.
3. **F4 — verified.** The source rule says "include a 20-minute stall check". The UR capture prompt (`do-work/user-requests/UR-175/input.md:45`) says "If nothing moves for 20 minutes, stop and report what blocked." The new sentence "if nothing moves for 20 minutes, the run stops and reports what blocked it" matches both. Sentence lengths: 5, 20 and 20 words. The added lines hold no semicolon and no em-dash. The `<date>-<slug>.md` file name stays out on purpose (D-15: `folder-structure.md` owns the folder and the naming), which is a reasonable call.
4. **Probe — verified.** `bash do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` → 14 PASS lines, `REQ-724 probe PASS`, exit 0. Negative test in a throwaway scratch copy of `skills/` (git-initialised, removed after): a stray `cheapest model tier` in `skills/do-work-toolbox/actions/stray-note.md` and a stray `unless the user says strict` in a knowledge crew file each produced a FAIL line naming the extra file, exit 1. Removing the phrase from the knowledge `background-agents.md` copy also produced a FAIL. So the probe still fails on an unplanned file and on a missing copy. The `sort | tr | sed` pipeline under `set -euo pipefail` keeps the `|| true` guard for a phrase with no match.
5. **Breakage and bloat — none found.** The shipped diff for the whole REQ is 9 files, +27/-3, which matches the REQ's re-qualification paragraph. `zsh|AskUserQuestion|haiku|EnterWorktree|Explore` is absent from all nine files. No new file, helper or test in `skills/`.

## Findings

**Minor:**

- N1 `do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh:41` and `:60` — the G2 option-shape check and the machine-name check read only the four `do-work` copies, so the five sibling copies are checked only by the three single-home phrases. The PASS message also still says "found once in" when it lists three files. This review confirmed the copies by `diff`, and the probe is a run artifact, not a shipped test. — impact-negligible → report only

## Review

**Overall: 92%** | <TIMESTAMP>

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
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*
