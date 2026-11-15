## Delta Review: REQ-735 (F1 fix)

**Approve.** The delta moves the ask-tool precedence clause into the bold lead of the rule, byte for byte as F1 suggested, and F1 is resolved.
Delta `e3f36c6c..3b2d0ff6` | builder fix `0b10c310` | 1 file, +1/-1 (`skills/do-work/crew-members/communication-style.md:25`)

### Checks

1. **Delta matches F1 and resolves it.** The new line 25 equals F1's replacement line exactly (`diff` of the suggested text, with the escaped backticks unescaped, against line 25 at `3b2d0ff6` prints nothing). The bold lead now reads "A decision for the user goes through the ask-user tool, unless the action in use says how to ask." The lead is no longer unconditional. `skills/do-work-toolbox/actions/tutorial.md:355` ("Plain text menus ... Do not use the ask-user tool") is an action that says how to ask, so it now falls under the exception and no longer contradicts the rule.
2. **No other change in meaning.** `git diff --word-diff` shows one change only: the clause "unless the action in use says how to ask" moves from the end of the second sentence to the end of the bold lead, and the periods move with it. No word is added or removed. The DECIDE & STATE sentence and the rest of the file are untouched. The REQ constraint (do not reword a shipped rule except to add the precedence condition) holds. The clause now scopes the whole rule instead of only the load step, which is the meaning F3 of the REQ asked for. `clear-questions.md` still loads on its own trigger (its JIT_CONTEXT names "an option menu"), so a plain-text menu still gets the question shape.
3. **Style holds.** The line names no tool, vendor or model. The bold lead is 17 words. The second sentence drops from about 27 to about 21 words, near the file's 20-word target. Plain words only.
4. **Restatement sweep.** `git grep` for "ask-user tool", "decision for the user", "prose list of questions", "says how to ask" and "AskUserQuestion" across `skills`, `docs`, `_dev`, `CLAUDE.md`, `AGENTS.md` and `README.md` finds no live restatement of the rule's lead outside `communication-style.md:25`. The other hits are consistent or are history:
   - `skills/do-work/actions/trace.md:39` tells the agent to ask with the ask-user tool. That is an action naming its own way of asking, which is consistent.
   - The three `clear-questions.md:3` JIT_CONTEXT lines name the ask-user tool as one trigger. That is consistent.
   - `CHANGELOG.md:23` and `skills/do-work/CHANGELOG.md:23` (the 0.322.0 entry) state the rule without the condition. History is not rewritten. See F1 below.
   - `docs/run-prompts/2026-11-15-move-global-rules-into-skill.md:24` quotes the global rule as input. It is a prompt record, not a rule.
   - `docs/changelog/CHANGELOG-2026-04-13-up-to-v0.64.1.md:317` records why the tutorial uses plain-text menus. History, consistent.

### Findings

- F1 `CHANGELOG.md:23`, `skills/do-work/CHANGELOG.md:23`: the 0.322.0 entry says to send a decision through the ask-user tool with no condition. History stays as written. The coming REQ-735 release entry should say that the ask rule now yields to an action's own way of asking. This extends first-review F5. — impact-negligible → report only

### Acceptance

`bash do-work/runs/work-2026-11-15-012821/REQ-735-probe.sh` passed at `3b2d0ff6`: exit 0, 15 PASS lines, "REQ-735 probe PASS", run once by the delta reviewer. The gate run (green, 127 s at `3b2d0ff6`) is the integrator's record. The reviewer did not rerun it.

---

## Review

**Overall: 96%** | delta `e3f36c6c..3b2d0ff6` | 2026-11-15

**Verdict:** Approve. First-review F1 is resolved by the exact suggested line. No new issue in the shipped text.

**Findings:**
- F1 `CHANGELOG.md:23`, `skills/do-work/CHANGELOG.md:23`: the 0.322.0 history states the ask rule without the condition; the REQ-735 release entry should name the new precedence condition (extends first-review F5). — impact-negligible → report only

**Acceptance:** Pass — REQ-735 probe exit 0 (15 PASS) at `3b2d0ff6`; gate green per integrator record.
**Restatement sweep:** the ask-tool precedence now lives in the bold lead at `communication-style.md:25` and is no longer stale against `tutorial.md:355`. No other live text restates the rule's lead. `trace.md:39` and the three `clear-questions.md:3` triggers are consistent. History only at `CHANGELOG.md:23` and `skills/do-work/CHANGELOG.md:23` (F1), plus the input quote at `docs/run-prompts/2026-11-15-move-global-rules-into-skill.md:24`.
**Follow-ups created:** None (1 finding report only)

*Reviewed by review-work action (delta)*
