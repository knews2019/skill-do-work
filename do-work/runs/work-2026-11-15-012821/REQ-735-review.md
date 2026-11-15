## Review: REQ-735

**Approve** — the three rule decisions, both precedence clauses and the new commit rule shipped as decided, with one wording fix recommended before release.
Route B | merge `e3f36c6c` (builder commit `82db9eb2`), range `08de436c..e3f36c6c`

### What's built
- R1 (read numbers as approximate) and R2 (keep existing UI during a redesign) now live once, in the always-on `communication-style.md:28-29`. Both `coding-guardrails.md` copies lost the same two lines.
- R3 (pick the model per agent) stays only in the three `background-agents.md` copies, which are identical. Each now yields to a model the action names.
- The ask-tool rule yields to an action's own way of asking, and `commit.md:223` says never to revert, stash or discard a change you did not make.

### Decisions / risks for you
- DR1: after `~/.claude/CLAUDE.md` is deleted, two rules no longer reach an ordinary session. R3 (decision D-03): a single lookup agent runs on the default model, which costs more tokens but gives the same result. F5's rule (decision D-05): it loads only when the commit action runs, so an ordinary session can revert another session's file. The second loss can destroy work, not only money. See F2.

### Findings

**Minor:**
- F1 `skills/do-work/crew-members/communication-style.md:25` — the clause "unless the action in use says how to ask" sits at the end of the second sentence, after "load `crew-members/clear-questions.md` first". A reader can attach it to the load step only. The bold lead "A decision for the user goes through the ask-user tool." stays unconditional, so it still contradicts `skills/do-work-toolbox/actions/tutorial.md:355` (plain-text menus, no ask-user tool). Recommend fix before release (moves the clause, adds no words). Replacement line: `- **A decision for the user goes through the ask-user tool, unless the action in use says how to ask.** Use the agent's ask-user tool when it has one, never a prose list of questions, and load \`crew-members/clear-questions.md\` first. A reversible choice whose answer is clear from context is the DECIDE & STATE tier of \`crew-members/coding-guardrails.md\`: make it yourself and report it in one line.` `clear-questions.md` still loads on its own trigger (JIT_CONTEXT names "an option menu"). The REQ-735 probe phrase does not change, so the probe still passes. — impact-rule-change → report only
- F2 `skills/do-work/actions/commit.md:223` — the only home of "Never revert, stash or discard a change you did not make" is the commit action. D-05 rejected an always-on copy because REQ-724 (moving the global standing rules into the skill) scoped exact-path staging to `commit.md`. That reason covers staging, which happens only when committing. Reverting or stashing can happen in any session. The global file applies the rule in every session today. No fix before release, because the REQ allowed this choice and D-05 records it. The maintainer should decide it before deleting the global file. If moved, the bullet goes after `communication-style.md:29` as `- Never revert, stash or discard a change you did not make. Another session or the user may own it.` and is deleted from `commit.md`. — impact-rule-change → report only

**Nit:**
- F3 `skills/do-work/crew-members/communication-style.md:7` — the scope sentence says "These rules govern conversational output". R1 is about reading a request and R2 about changing UI, so two bullets sit outside that scope. D-02 names this risk. No fix. — impact-negligible → report only
- F4 `skills/do-work/docs/commit-guide.md:26-32` — the "Key rules" summary for humans does not list the new revert rule. It does not claim to be complete. No fix needed. — impact-negligible → report only
- F5 `CHANGELOG.md:25-26` and `skills/do-work/CHANGELOG.md:25-26` — the 0.322.0 history entry still names `coding-guardrails.md` as the home of R1 and R2, and gives the model sentence without the new clause. History is not rewritten. The coming REQ-735 release entry has to say that the rules moved. — impact-negligible → report only
- F6 `skills/do-work-toolbox/actions/architecture-report.md:81` — R2 now loads in every session. This action asks for a full redesign each run with "no requirement to reproduce the prior design". Invoking the action counts as the user asking, and each run writes a new HTML file, so the conflict is weak. No fix. — impact-negligible → report only

### Requirements Checklist

- [x] R1, R2 and R3 each have a decision with value and risk — delivered (D-01, D-02, D-03)
- [x] R3 respects REQ-724 D-05: no copy and no pointer in `communication-style.md` — delivered (grep for "model" finds nothing there)
- [x] A moved rule is stated once, with no copy left in its old home — delivered (R1 and R2 appear only in `communication-style.md:28-29`, CHANGELOG excluded)
- [x] Lock-step copies stay identical — delivered (the three model paragraphs hash `55740a80`. The two guardrails copies differ only by the package-path rewrites and the "Earned defense" paragraph, and both differences were already there at `08de436c`)
- [x] F3 folded in — delivered (the model clause scopes cleanly; the ask clause has the placement problem in F1)
- [x] F5 folded in, one shipped home — delivered (`commit.md:223`)
- [x] No machine, vendor, model or tool name, and nothing under `skills/` cites `CLAUDE.md` — delivered (no added line matches the name regex; `shipped-package-reference-contract.sh` exit 0)
- [x] Aliases and Examples unchanged; file kept short — delivered (both sections hash the same at `547d4cc5` and now; file grew from 99 to 101 lines)
- [x] Release follows the releases prime — delivered as a proposal (minor bump, impact-rule-change). The integrator writes it, and `VERSION` is still 0.322.1.

### Instruction-file checks

1. **Moved rules against `/Users/t2/.claude/CLAUDE.md`. Meaning is preserved for every rule.**
   - Numbers. The global text is "Numeric targets are approximate unless the user says strict." The shipped text is "Read a number in the request (a size, a count, a time) as approximate unless the user says strict." The shipped text is a little broader: "a time" can also be read as a deadline. That wording came from REQ-724 and moved byte for byte, and the constraints forbid rewording, so it is not a finding.
   - UI. The global text is "Never remove existing UI while redesigning unless asked." The shipped text says "...unless the user asked for it." Same meaning.
   - Ask tool. The global rule is "goes through the AskUserQuestion tool: recommended option first, value and risk for each option. Never a prose list of questions." The shipped rule keeps the tool and the no-prose-list parts, and adds the precedence clause. The option shape lives in `clear-questions.md:29` in all three copies.
   - Model. The global rule gives lookups, searches, file sweeps and Explore agents "model: haiku", and keeps "inherit the session model" for planning, review and build agents. The shipped rule says "cheapest model tier that can do them" and "keep the session model unless the action names another model". This is the same rule without names, plus the precedence clause. The cheap-model sentence has no clause, and no shipped action names a model for a lookup (grep for Haiku/Sonnet/Opus in `skills/` finds nothing).
   - Concurrent sessions. The global rule is "never stage or revert files you did not change" and "never `git add -A` or `git commit -a`". The shipped rule says "never revert, stash or discard a change you did not make", next to "Stage only the specific files", which also rules out `git commit -a` (D-05). This is broader than the global rule and does not conflict with it. `commit.md:69` (unstage a pre-staged secret path) and `fan-out-reference.md:76` (`git restore --staged`) only unstage, so they do not conflict either.
2. **No rule ships twice.** R1 and R2 are each in exactly one file under `skills/`. F5 is only in `commit.md`. The model rule is in the three lock-step copies only. This was checked with the probe and an independent `git grep` across `skills`, `docs`, `_dev`, `CLAUDE.md`, `AGENTS.md` and `README.md`.
3. **The style of `communication-style.md` holds.** The two moved bullets are short and plain, and they name no tool, vendor, model or machine. The ask bullet's second sentence is now about 27 words with four instructions in it, over the file's 20-word target. That is part of F1, and the replacement line fixes it.
4. **D-03 holds.** `communication-style.md` has no model text and no pointer to `background-agents.md`.
5. **Ask-clause scope.** Not clear. See F1.
6. **Anti-bloat.** The diff adds 0 helpers, 0 options, 0 files and 0 tests. It is 7 prose files, +10/-11, all declared in Scope. The probe is a run artifact under `do-work/runs/` and does not ship.

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `bash do-work/runs/work-2026-11-15-012821/REQ-735-probe.sh` passed with exit 0, 15 PASS lines and "REQ-735 probe PASS", run by the reviewer at merge `e3f36c6c`.
- `bash _dev/tests/shipped-package-reference-contract.sh` passed with exit 0, run by the reviewer.
- The integrator's gate run (`maintainer-verify.sh`, exit 0, 123 s) is recorded in the REQ. The reviewer did not rerun it.
- The builder's RED evidence at base `547d4cc5` has 7 FAIL lines, which matches the rule state there.

### Suggested Additional Testing

- Deployment: unassessed. After release, reinstall into a consumer project and confirm that the linked `communication-style.md` carries lines 28-29.
- Live acceptance: unassessed. In an ordinary session of an installed project (no do-work action), ask for "about 300 words" and for a redesign of a page. Confirm that the agent treats the number as approximate and keeps the existing controls.
- Live acceptance: unassessed. Run `do-work tutorial` and confirm that it prints a plain-text menu and does not open the ask-user tool. This check is worth more after the F1 fix.

### Scores (on the record — not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | Every item delivered; the F3 ask clause is ambiguous in placement (F1) |
| Code Quality | 90% | Byte-identical moves; F1 placement |
| Test Adequacy | 90% | Prose change; red-green run probe plus contract test, no decorative shipped test |
| Scope | 100% | Exactly the seven declared files; D-08 and D-09 cover the visible choices |
| Risk | Low | Two rules stay out of ordinary sessions by recorded decision (DR1, F2) |
| Acceptance | Pass | Probe and contract test at the merge |

### Follow-ups created
None (6 findings report only)

---

## Review

**Overall: 94%** | 2026-11-15T08:02:13Z

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
- F1 `skills/do-work/crew-members/communication-style.md:25`: "unless the action in use says how to ask" reads as scoping only the clear-questions load, and the bold lead stays unconditional against `tutorial.md:355`. Fix before release: move the clause into the bold lead (exact line in the full report). — impact-rule-change → report only
- F2 `skills/do-work/actions/commit.md:223`: the never-revert rule loads only with the commit action, so an ordinary session loses it once the global file is deleted (D-05). This is a maintainer decision before that delete. — impact-rule-change → report only
- F3 `skills/do-work/crew-members/communication-style.md:7`: the "conversational output" scope sentence does not cover R1 and R2 (D-02 risk). — impact-negligible → report only
- F4 `skills/do-work/docs/commit-guide.md:26-32`: the Key rules summary omits the new revert rule. — impact-negligible → report only
- F5 `CHANGELOG.md:25-26`, `skills/do-work/CHANGELOG.md:25-26`: the 0.322.0 history still names `coding-guardrails.md` as the home of R1 and R2. History is not rewritten, and the release entry must state the move. — impact-negligible → report only
- F6 `skills/do-work-toolbox/actions/architecture-report.md:81`: a full redesign each run, now co-loaded with R2. The conflict is weak. — impact-negligible → report only

**Acceptance:** Pass — the implementation and integration stages were exercised: the REQ-735 probe (15 PASS, exit 0) and `shipped-package-reference-contract.sh` (exit 0) at merge `e3f36c6c`. Deployment and live acceptance are unassessed.
**Restatement sweep:** redefined the home of R1 and R2 (now `communication-style.md:28-29`; no live restatement in `skills/`, `docs/`, `_dev/`, `CLAUDE.md`, `AGENTS.md` or `README.md`; history only at `CHANGELOG.md:25` and `skills/do-work/CHANGELOG.md:25`, F5, plus the input quotes in `docs/run-prompts/2026-11-15-move-global-rules-into-skill.md:24,35-36,48`), the ask-tool precedence (stale against `tutorial.md:355` until F1 is fixed), the model-rule precedence (three copies identical, consistent with `ai-evaluation-loop-reference.md:19-21`; history at `CHANGELOG.md:26`), and the new revert rule at `commit.md:223` (consistent with `work-reference.md:31` and `restart-with-parallel-handoff.md:68`; `commit-guide.md:26-32` omits it, F4). Inherited elements from REQ-724 were rechecked with no stale text: the DECIDE & STATE pointer (now `coding-guardrails.md:32`, cited without a line number), the clear-questions option shape (three copies identical apart from package paths), the per-agent model rule, the numeric-target and redesign lines, and the stall-check wording (`communication-style.md:30`). REQ-736 recorded nothing redefined.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)
**Anti-bloat:** 0 unrequested helpers, options, files or tests
**Full report:** `do-work/runs/work-2026-11-15-012821/REQ-735-review.md`

*Reviewed by review-work action*
