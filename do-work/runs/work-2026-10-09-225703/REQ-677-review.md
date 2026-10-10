## Review: REQ-677

**Approve**. The journey-qa action does what the REQ asks, it is wired into every toolbox surface beside source-audit, and the combined-journey exercise reproduces on the merged tree. Four small wording findings, all report only.
Route B | merge range 68bd8bd7..1cb75f9f (builder commit 4b2c7a5c, merge 1cb75f9f)

### What's built
- `do-work-toolbox journey-qa <target> [--brief <path>]` (in `skills/do-work-toolbox/actions/journey-qa.md`, 111 lines) reproduces the reported sequence first, then runs isolated and combined transitions. It waits on observable state and records revision, environment, steps, expected, actual and evidence. It classifies each result as passed, product defect, test defect or unresolved, and names the smallest repair without applying it.
- User guide `skills/do-work-toolbox/docs/journey-qa-guide.md`, route row, argument-hint, both help menus, README paragraph and `toolbox_actions` entry are all present.
- Missing: no reverse pointer from `ui-review` to journey-qa (F4).

### Decisions / risks for you
- Builder decision D-02 (a browser tool the session already provides counts as a browser tool). This complies with the Agent Compatibility rule. It is generalized language and names no tool API, so it follows that rule more closely than the ui-review step it cites. The cost is that journey-qa and ui-review now disagree on detection (F2). Option O1 (recommended): keep the clause and later widen ui-review Step 2 item 4 the same way. Value: honest results in sessions that only have an MCP or built-in browser. Risk: two detection rules until ui-review catches up. Option O2: drop the clause. Value: one rule. Risk: a false "unresolved" in browser-capable sessions.

### Findings

**Important:**
- None.

**Minor:**
- F1: The guide disagrees with the action on browser detection. `skills/do-work-toolbox/docs/journey-qa-guide.md:24` says "the same detection as `ui-review`: Playwright CLI or the Bowser skill", but `actions/journey-qa.md:44` also counts a browser tool the session already provides. A user who only has a session browser reads the guide and installs Bowser for nothing. Fix: replace line 24 of the guide with `- A browser tool. It looks for Playwright CLI or the Bowser skill the way \`ui-review\` does (\`do-work-toolbox install bowser\` installs both), and a browser automation tool your agent session already provides also counts. Without one, rendered journeys are unresolved, never passed.` — impact-user-visible → report only
- F2: Builder decision D-02 creates two detection rules. `actions/journey-qa.md:44` counts a session-provided browser tool. `actions/ui-review.md:60-67` (Step 2 item 4) and `:150-152` (Step 8.5, "If Playwright CLI or the Bowser skill was detected") do not. The citation is honest because it says "also counts", but one capability check now has two meanings in one package. Fix (O1): in `skills/do-work-toolbox/actions/ui-review.md`, insert a bullet before the "**Neither available**" bullet in Step 2 item 4: `   - **A browser automation tool the session already provides**: counts as available. Use it for Step 8.5 the same way.` Then change the Step 8.5 first sentence to `If a browser tool was detected in Step 2.4, use it to validate the rendered UI.` After that, journey-qa line 44 can drop its clause. Fix (O2): delete `; a browser automation tool the session already provides also counts` from journey-qa.md:44. — impact-rule-change → report only
- F3: The Step 6 cleanup is broader than the read-only claim supports. Line 5 says "It changes no project file", and line 77 says "If the browser tool wrote anything inside the project, remove it". Removing the action's own leftovers is consistent with read-only in intent, because the net change is zero. The wording has two problems. "anything" also covers a modified tracked file, and undoing that change is a destructive revert. And plain `git status --porcelain` collapses an untracked directory into one line, so the before/after diff can miss a new file or blame the tool for a file the user created in the meantime. Fix: line 46, replace `` and the output of `git status --porcelain`, so Step 6 can prove nothing changed. `` with `` and the output of `git status --porcelain --untracked-files=all`, so Step 6 can prove nothing changed. ``. Line 77, replace the whole line with `` Compare `git status --porcelain --untracked-files=all` with the Step 2 record. Remove only new untracked paths the browser tool wrote (session files, screenshots, traces) and name each in the report. Never revert a tracked file or delete a path you cannot tie to the browser tool. Report it instead. `` — impact-user-visible → report only
- F4: Next-step guidance only points one way. UR-151 INTEGRATION asks for "relevant next-step guidance". journey-qa sends users to ui-review (`journey-qa.md:17`), but `skills/do-work-toolbox/actions/ui-review.md:23-26` (Do NOT use when) never sends them to journey-qa. This is the same gap that the REQ-676 review (source-audit) recorded as its F1 for slop-check. The REQ's write set left ui-review out, so this is not scope drift by the builder. Fix: insert after `skills/do-work-toolbox/actions/ui-review.md:25`: `- The user wants a whole user journey checked in a browser (combined steps, product defect versus test defect) → \`actions/journey-qa.md\`.` — impact-user-visible → report only

**Nit:**
- None.

### Requirements Checklist

- [x] R1: Reads the requirements, the testing guidance and the brief as data, with prompt-injection loaded first. Delivered (Step 1, lines 38-40).
- [x] R2: Reuses existing checks and reproduces the reported sequence first. Delivered (Step 2 line 46, Step 3 line 52).
- [x] R3: Combined transitions. An isolated pass never stands in for a combined one, and examples are marked illustrative. Delivered (line 7, Step 4 line 56).
- [x] R4: Waits on observable state and keeps instrumentation outside timed spans. Delivered (lines 52, 58).
- [x] R5: Records revision, environment, steps, expected, actual and evidence per journey. Delivered (Step 2, Step 4, Output Format lines 92-98).
- [x] R6: Four classes. An environment limit is unresolved with its reason, and emulation is labeled. Delivered (Step 5 table, lines 44, 60).
- [x] R7: Smallest justified repair, remaining checks, diagnostic only, user-run capture line. Delivered (Step 6, lines 100-111).
- [x] R8: Detection cites ui-review Step 2 item 4 instead of copying it. Delivered, with one added clause (D-02, see F2).
- [x] R9: Toolbox ownership is justified in the blockquote. Delivered (line 3).
- [x] Integration: argument-hint, route row, both help menus, README, `toolbox_actions` and the guide. Delivered. Per-command help is served by the router (`skills/do-work-toolbox/SKILL.md:38`).
- [x] Constraints: read-only on project source, evidence in a temporary or ignored directory, no browser means unresolved. Delivered (the Step 6 wording is in F3).
- [x] `suite/modules.tsv`, the installer and the updater are unchanged. Delivered (not in the range).
- [~] UR "relevant next-step guidance". Partially delivered (F4).

### Merge seam with REQ-676 (source-audit)

Both names are present in every place at 1cb75f9f:
- `skills/do-work-toolbox/SKILL.md:4` (argument-hint: `ui-review | journey-qa | ...` and `slop-check | source-audit | ...`), and the route rows at `:23` and `:28`.
- `skills/do-work-toolbox/actions/help.md:13` and `:18`. All 21 menu lines put the description at the same column (34, 1-based).
- `skills/do-work/actions/help.md:38` (`present-video · slop-check · source-audit · journey-qa`). That row is 75 characters, and the widest row in the toolbox block is 79.
- `README.md:116` (journey-qa) and `:120` (source-audit).
- `_dev/tests/staged-skills-contract.sh:153` and `:158`.

Both run probes (`REQ-676-probe.sh`, `REQ-677-probe.sh`) exit 0 on the merged tree.

### Restatement sweep

This REQ redefines the list of toolbox actions (one more action). I checked every place that lists or counts them:
- All the surfaces above agree.
- `skills/do-work-toolbox/actions/tutorial.md:327` is a "key commands" sample for reviews, not a full list (it also omits slop-check and source-audit).
- `skills/do-work/next-steps.md:8` derives commands from the owning router, so it has no list.
- The toolbox `crew-members/prompt-injection.md`, `anti-slop.md` and `testing.md` caller lists are marked "illustrative, not exhaustive" or have no caller list.
- The Go `toolboxcommands.Handlers()` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go:21-30`) registers only deterministic CLI phases. journey-qa has none.
- `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` covers moved commands only. `journey qa` and `user journey` do not appear there, and the heavy contract's retired-trigger counts still pass.
- No numeric count of toolbox actions exists anywhere in `skills/`, `README.md` or `_dev/primes/`.

Result: no stale restatement. The only gap is the reverse pointer in F4, which is next-step guidance, not a stale list. This is not the wave's last integration, so no inherited sweep was run.

### Acceptance Testing

**Result: Pass.** Implementation and integration were assessed. Deployment was assessed partly (fresh install only). Live acceptance was not assessed.

Re-run by this review, on the merged tree:
- `REQ-677-probe.sh` exit 0 and `REQ-676-probe.sh` exit 0. Both include `shipped-package-reference-contract.sh`.
- `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` exit 0 in 32 s ("staged skills contract: PASS").
- `bash _dev/tests/contracts/core-checks.sh` exit 0 in 8 s (near-identical cross-file pairs: 0).
- `git diff --check 68bd8bd7..1cb75f9f` is clean. There are no em-dashes in either new file. The working tree matched 1cb75f9f for `skills/`, `README.md` and `_dev/` before these runs, and `git status` was unchanged after them.
- Route resolves: `SKILL.md:23` sends `journey-qa` to `./actions/journey-qa.md`, which exists.
- Per-command help: the router rule at `SKILL.md:38` says per-command help "reads the selected action without executing it", and the format comes from `skills/do-work/actions/help.md:62` (Input and When to Use, at most 15 lines). journey-qa.md has a `## When to Use` section and an `## Input` section with the exact usage line and two examples, so help prints usage and runs no journey. I took the rendered help text from the builder's hand-back and did not re-render it.
- Real behavioral re-run: I copied the builder's `zoom-app` fixture (revision e3fddf2) into a fresh `mktemp -d` directory and served it with a server limited to 4 minutes. I ran the builder's `journeys.sh` through `playwright-cli` 0.1.14 (headless Chrome 155, 1280x720), started from the evidence directory. Results: zoom alone gave scale 2, x 0. Pan alone gave scale 1, x 50. Reset alone gave scale 1, x 0. Zoom → pan → reset gave scale 1, x 50 (product defect reproduced). Pan → zoom → reset gave scale 1, x 0. `playwright-cli` wrote `.playwright-cli/` into the evidence directory only, which confirms builder decision D-03. The fixture tree was unchanged afterwards (`git status --porcelain --ignored` showed no diff), and the browser and server were closed.
- Packaging at the merge revision: `git archive 1cb75f9f` was installed with the canonical installer into a fresh Git consumer. Results: "install-suite: success", "rollback: not_needed", version 0.305.91. `actions/journey-qa.md`, `docs/journey-qa-guide.md` and `actions/source-audit.md` are present under `.claude/skills/do-work-toolbox/`. All seven relative links from the installed action resolve. The sha256 of `qa/brief.md`, `do-work/queue/REQ-001-sample.md` and `src/app.js` is the same before and after install.

Taken from the builder hand-back, not re-run:
- The wrong-selector variant was classified test defect.
- The phone-emulation run was labeled emulation.
- The simulated no-browser run was unresolved.
- The updater pass returned already-current at 0.305.89 and did not re-stage files.

### Suggested Additional Testing

- Deployment: partly unassessed. Re-run the updater check after this REQ's release bumps VERSION, so the updater actually re-stages files instead of stopping as already current.
- Live acceptance: unassessed. After release, run `just do-work-update` in a real consumer and then `do-work-toolbox journey-qa help` there. Confirm that it prints usage only.
- Real no-browser agent: run the action in an agent with no browser tool at all (the builder simulated this with a narrowed PATH). Confirm the report says unresolved, not passed.
- Physical device: no reviewer or builder had one. One run on a real phone would confirm that the report keeps emulation and device results apart.

### Scores (on the record, not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | Every REQ requirement delivered. UR next-step guidance is partial (F4). |
| Code Quality | 90% | Short and on template. Guide and action disagree on detection (F1), and the Step 6 cleanup wording is too broad (F3). |
| Test Adequacy | 90% | Real behavioral exercise per the maintainer's one-off choice. The probe is keyword-level. The no-browser case was simulated. |
| Scope | 100% | 7 files declared, the same 7 touched. Decisions D-01 to D-07 are recorded. |
| Risk | Low | Step 6 cleanup could revert or delete a file the tool did not write (F3). |
| Acceptance | Pass | Implementation and integration. Deployment partly (fresh install at the merge revision). Live acceptance unassessed. |

### Follow-ups created
- None (4 findings report only)

## Review

**Overall: 94%** | <timestamp>

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
- F1: `skills/do-work-toolbox/docs/journey-qa-guide.md:24` says detection is only Playwright CLI or Bowser, but `actions/journey-qa.md:44` also counts a session-provided browser tool. The guide misleads users who only have a session browser. — impact-user-visible → report only
- F2: Builder decision D-02 gives journey-qa a wider browser-detection rule than `actions/ui-review.md:60-67` and Step 8.5. This follows the Agent Compatibility rule, but one capability check now has two meanings. Widen ui-review (O1) or drop the clause (O2). — impact-rule-change → report only
- F3: `actions/journey-qa.md:77` "remove anything the browser tool wrote" could revert a tracked file or delete a file the user created, and plain `git status --porcelain` at `:46` collapses untracked directories. Narrow it to new untracked tool output and use `--untracked-files=all`. — impact-user-visible → report only
- F4: No reverse next-step pointer. `actions/ui-review.md:23-26` never sends users to journey-qa, although UR-151 asks for relevant next-step guidance (the same gap as the REQ-676 review's F1). — impact-user-visible → report only

**Acceptance:** Pass. Implementation and integration assessed: probes, heavy staged-skills, core-checks, a real zoom → pan → reset re-run that reproduces the product defect, and per-command help confirmed from the router rule. Deployment partly assessed (fresh install at 1cb75f9f). Live acceptance unassessed.
**Restatement sweep:** redefined the list of toolbox actions (one more). Checked: toolbox `SKILL.md` argument-hint and route table, both help menus, `README.md:116,120`, `toolbox_actions`, `tutorial.md:327` (a sample, not a list), `next-steps.md` (derived from the router), the toolbox crew-member caller lists (illustrative), the Go `toolboxcommands.Handlers()` (deterministic CLI phases only), and the retired-trigger fixture (does not apply). No count of toolbox actions exists. All agree.
**Suggested testing:** 4 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Re-check of 68bd8bd7..ec1746ee

**Overall: 95%** | Acceptance: Pass. Implementation and integration were re-checked. Deployment and live acceptance are unchanged from the first review.

Delta read: `git diff 1cb75f9f..ec1746ee`. It contains builder commit 84738e35 and merge ec1746ee, and changes 2 files (3 insertions, 3 deletions). Both files sit inside the declared Scope.

- **F1 closed.** `skills/do-work-toolbox/docs/journey-qa-guide.md:24` now matches the proposed text word for word. It agrees with `actions/journey-qa.md:44`: Playwright CLI or Bowser as ui-review detects them, and a browser tool the session already provides also counts.
- **F3 closed.** `actions/journey-qa.md:46` and `:77` now match the proposed text word for word. Both the Step 2 record and the Step 6 comparison use `git status --porcelain --untracked-files=all`. The cleanup removes only new untracked paths the browser tool wrote, never reverts a tracked file, and reports anything it cannot tie to the tool.
- **No new finding.** Read-only wording: line 5 ("changes no project file") and Step 6 now agree, because the only removal is of the run's own new untracked output. Agent Compatibility: the new text uses generalized language and names no tool API. There are no em-dashes in either file, and `git diff --check` is clean. Citations: no path changed, and the heavy staged-skills contract passes.
- **Checks re-run on the merged tree.** The working tree matches ec1746ee for `skills/`, `README.md` and `_dev/`. `REQ-677-probe.sh` exit 0 (includes `shipped-package-reference-contract.sh`). `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` exit 0 ("staged skills contract: PASS"). `bash _dev/tests/contracts/core-checks.sh` exit 0. `git status` was unchanged after the runs.
- **Still open, report only (outside Scope):** F2 (two browser-detection rules, journey-qa versus ui-review) — impact-rule-change → report only. F4 (no reverse pointer from ui-review to journey-qa) — impact-user-visible → report only.

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Follow-ups created:** None (2 findings report only)
