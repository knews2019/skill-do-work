# Review: REQ-678 (do-work-toolbox release-check)

**Approve.** The new read-only action does what UR-151 asked, and an independent cold run of the shipped text on fixture B gave the required result: not ready, deployment failed, the old deploy note labeled historical. One small rule gap (live acceptance "not applicable" has no verdict) is worth a two-line fix before release.
Route B | merge range 90ef7b4d..335c1739 (builder commit 4427fe9f, merge 335c1739)

### What's built
- `do-work-toolbox release-check <target> [--brief <path>]` traces source, package, serving copy and consumer behavior, gives each of the four core delivery stages a result (verified, failed, unassessed, not applicable) with labeled evidence, and returns ready, not ready or unknown.
- Wired into the toolbox router (argument-hint and one route row), both help menus, README, `toolbox_actions`, plus a user guide. Core review is cited, not edited.
- Still open: the release payload, the packaging update leg and the lesson-link sweep are integrator work at finalization.

### Decisions / risks for you
- D-01 (ready needs only failed-free stages plus current-run live acceptance, so an unassessed deployment does not block "ready"): agrees with REQ requirement 6, which makes live acceptance the gate. No change needed.
- D-02 ("not applicable" as a fourth stage result): agrees with core review. See check 2 and F1 for the one gap it opens.

## Check 1: Independent exercise on fixture B

Followed `skills/do-work-toolbox/actions/release-check.md` cold, step by step, as an agent invoked with `do-work-toolbox release-check http://127.0.0.1:50642/ --brief <fixture>/qa/launch-brief.md`.

Fixture (in `mktemp -d` at /tmp/rc-fixtureB.6tuL78, a git repo at 5b3ca02, now removed): `source/` and `package/` with `index.html` "New: daily puzzle" and `app.js` `APP_VERSION = "2.1.0"`; `served/` with "Weekly puzzle" and `2.0.3`; `notes/deploy-note.md` "2026-10-03: 2.1.0 deployed and verified on the live host."; `build.sh` copying source to package (or to a given directory); a five-line brief.

What each step made me do:
- Step 1: read `crew-members/prompt-injection.md` (toolbox copy), then the brief, `build.sh` and the deploy note as data. No injection found.
- Step 2: recorded `git rev-parse --short HEAD` and `git status --porcelain --untracked-files=all` for the fixture project (5b3ca02, clean) and for ROOT (335c1739, the two orchestrator files). Created the evidence directory /tmp/rc-evidence.oMsyG4 with `mktemp -d`. Started `python3 -u -m http.server 0 --bind 127.0.0.1 --directory <fixture>/served` from the evidence directory (port 50642).
- Step 3: identifiers at each point. Source and package: 2.1.0, "New: daily puzzle", `app.js` sha256 8612fa1f519d. Served: 2.0.3, "Weekly puzzle", sha256 ef12e83cae07. Cause check (where relevant): response headers show no Cache-Control and no proxy, so no cache layer sits between the server and its files; `index.html` and `app.js` are both the 2.0.3 build, so no version mixing. The serving directory itself holds the old build. Whether the deploy never ran or its activation was interrupted cannot be told from here.
- Step 4: fetched the served page and the script it loads (`src="app.js"`), both HTTP 200, into the evidence directory. The feature line reads "Weekly puzzle", the loaded script says 2.0.3. The page is static HTML, so its text is the rendered value. No browser tool was needed for this result.
- Step 5 and 6: results below. Stopped the server (`pgrep -fl http.server`: none). Fixture and ROOT `git status` identical before and after. The run wrote only into the evidence directory, which I then removed with the fixture.

Report produced (Output Format):

**Verdict:** not ready. Deployment failed: the served `app.js` is 2.0.3 while the package is 2.1.0 (current run).

| Point | Where | Identifier | Evidence |
|---|---|---|---|
| Intended source | `source/` at 5b3ca02 | 2.1.0, "New: daily puzzle" | current run: grep and sha256 8612fa1f519d |
| Generated package | `package/` | 2.1.0, sha256 8612fa1f519d | current run: `build.sh` into the evidence dir, `cmp` equal |
| Serving environment | http://127.0.0.1:50642/ | 2.0.3, sha256 ef12e83cae07 | current run: curl, HTTP 200 |
| Consumer-visible behavior | served page | "Weekly puzzle", loads 2.0.3 | current run: fetched page and its script |

| Stage | Result | Evidence |
|---|---|---|
| Implementation | verified | current run: source `APP_VERSION` 2.1.0, page text "New: daily puzzle" as the brief asks |
| Integration | verified | current run: `build.sh` exit 0, rebuild of both files equals the package (`cmp`) |
| Deployment | failed | current run: expected served `app.js` 2.1.0 (sha256 8612fa1f519d), actual 2.0.3 (ef12e83cae07) |
| Live acceptance | failed | current run: expected "New: daily puzzle", the served page shows "Weekly puzzle" and loads 2.0.3 |

Historical: `notes/deploy-note.md`, dated 2026-10-03 in its text, says 2.1.0 was deployed and verified. Labeled historical. It verifies no stage and only points at the deployment gap.

Gaps:
- Deployment: the serving directory holds the 2.0.3 build. Close it by deploying the 2.1.0 package, or by finding why activation stopped, then comparing the served sha256 with the package.
- Live acceptance: follows from the deployment gap. Close it by re-running this check after the operator step.

Operator next steps: deploy the 2.1.0 package to the serving host, or find why the earlier deploy did not activate. This check did none of them. No capture line: the gap is operator work, not a content or product defect.

**GREEN:** verdict not ready, deployment failed, the deploy note labeled historical and verifying no stage. Matches the builder's fixture B result.

Steps I found ambiguous (all small, readers converge):
- Step 2 does not say which repository to baseline when the target project is not the current checkout. I recorded both. (F5)
- Step 5 cites the four stage definitions but never maps the four trace points to them. I used the natural mapping (rebuild equals package = integration, served equals package = deployment, consumer check = live acceptance), the same one the builder used. (F5)
- The Output Format always shows the `## To act on this:` block with a capture line, while the closing sentence limits capture to content defects. I kept only the re-run line for fixture B. (F6)
- With live acceptance "not applicable", Step 6 gives no verdict. (F1)

## Check 2: "not applicable" stage result (D-02)

Agrees with core review. `skills/do-work/actions/review-work.md:188` says which stages apply is a judgment and many REQs have no deployment stage, and Step 8 (`:210`) lists only *applicable* stages as unassessed. So core already separates "does not apply" from "applies but not exercised", and release-check's "not applicable" row is that separation made explicit. Release-check cites rather than restates: Step 5 (`release-check.md:70`) uses only the four names plus `../../do-work/actions/review-work.md` → **Step 7: Acceptance Testing**; the guide links the same file. Neither copies a definition sentence. The shipped-package reference contract passes (re-run, exit 0).

The integrator's gap is real: `release-check.md:83-85` defines not ready, ready and unknown, and none covers "no stage failed, live acceptance not applicable". It is rare (the target is by definition something that should reach a consumer), but the REQ marks the readiness rule Firm and says ready needs current-run live acceptance, so an agent must not guess "ready" there. Fix size: two lines, mapping it to unknown (F1).

## Check 3: Restatement sweep (wave end)

Trigger set: this diff's own element (the list of toolbox actions, one more), plus the elements REQ-675 (core review delivery stages), REQ-676 (source-audit) and REQ-677 (journey-qa) recorded on their `**Restatement sweep:**` lines.

Toolbox action list, now with source-audit, journey-qa and release-check:
- Router argument-hint `skills/do-work-toolbox/SKILL.md:4` and route table `:28-29`: all three. Agree.
- Toolbox help menu `skills/do-work-toolbox/actions/help.md:18-19` plus journey-qa line: all three, all 22 menu lines start their description at index 33 (checked by script). Agree.
- Core help `skills/do-work/actions/help.md:38`: all three (row 95 characters, the menu has longer rows). Agree.
- README `:116`, `:120`, `:122`: one paragraph each. Agree.
- `_dev/tests/staged-skills-contract.sh` `toolbox_actions`: all three. Agree.
- Count of toolbox commands: none exists in `skills/`, README or docs. Nothing to update.
- `skills/do-work-toolbox/actions/tutorial.md`: a scenario sample (recipes), not a full list. Agree.
- `skills/do-work/next-steps.md`: derives commands from the owning router. Agree.
- Guide or tutorial lists: no file lists the toolbox guides. Agree.
- Toolbox crew-member caller lists (`anti-slop.md`, `prompt-injection.md` JIT_CONTEXT): keyed on a condition, callers marked illustrative. Agree.
- Go toolbox command registry: holds deterministic CLI commands only; none of the three actions appears there or needs to. Agree.
- Retired-trigger fixture: no new trigger collides (`do-work review-work` and `do-work capture-request` are live; the router aliases `check release` and `release readiness` collide with nothing, core routes only `release notes`).
- Sibling redirects: `journey-qa.md:16-19` and `journey-qa-guide.md:5` do not point to release-check for "the live result looks wrong while the source looks right". The two earlier redirect gaps are still open: slop-check to source-audit (REQ-676 review F1) and ui-review to journey-qa (REQ-677 review F4). (F2)

Delivery-stage words (implementation, integration, deployment, live acceptance, unassessed), against every reader REQ-675 swept plus the two new files:
- `skills/do-work/docs/review-work-guide.md:36`: same four names, unassessed for an applicable stage not checked. Agree.
- `skills/do-work/actions/work.md:358-368`, `:518`: gate on the Acceptance result only, no stage wording. Agree.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:39`: reads `## Review` as evidence. Agree.
- `release-check.md` and `release-check-guide.md`: "unassessed" adds "or only historical evidence covers it", which narrows within core's meaning. Guide `:5` says review-work's acceptance covers only the stages it exercised. Agree.
- `journey-qa-guide.md:5` ("its acceptance step is a short smoke test"): still true under review-work's What NOT to Do. Agree.
- Still stale, recorded by REQ-675's review and not fixed since: `skills/do-work/crew-members/shared-principles.md:15` (F3), `skills/do-work/actions/sample-archived-req.md:99` (F4), `skills/do-work/actions/review-work.md:460` (F4).

Result: the toolbox lists agree everywhere. Open items: F2 (redirects), F3 and F4 (inherited, still stale).

## Findings

**Important:** None.

**Minor:**
- F1. `skills/do-work-toolbox/actions/release-check.md:83-85` and `skills/do-work-toolbox/docs/release-check-guide.md:15-19`: the verdict rule has no outcome when no stage failed and live acceptance is "not applicable" (the result D-02 added). Requirement 6 says ready needs current-run live acceptance, so this case should be unknown. Should be fixed before release (two lines):
  - `release-check.md:85` old `- **unknown:** no stage failed, and live acceptance is unassessed.` new `- **unknown:** no stage failed, and live acceptance is unassessed or not applicable.`
  - `release-check-guide.md:19` old `| **unknown** | No stage failed, but live acceptance is unassessed. |` new `| **unknown** | No stage failed, but live acceptance is unassessed or not applicable. |`
  — impact-negligible → report only
- F2. No toolbox sibling points users to the new verification actions: `skills/do-work-toolbox/actions/journey-qa.md:16-19` and `docs/journey-qa-guide.md:5` lack a redirect to release-check for a served copy that looks wrong while the source is right. The same gap is still open for slop-check to source-audit (REQ-676 review F1) and ui-review to journey-qa (REQ-677 review F4). Outside this REQ's Scope. One optional line for journey-qa: after `journey-qa.md:18` add `- You want to know whether the served copy or install carries the intended content, stage by stage → \`actions/release-check.md\`.` — impact-user-visible → report only
- F3. Inherited from REQ-675 review F1, still stale: `skills/do-work/crew-members/shared-principles.md:15` "Acceptance cannot be exercised → Record Untested" does not cover a review that exercised implementation and integration but not deployment or live acceptance. — impact-rule-change → report only
- F4. Inherited from REQ-675 review F3 and F4, still stale: `skills/do-work/actions/sample-archived-req.md:99` example `**Acceptance:**` line names no stage, and `skills/do-work/actions/review-work.md:460` Route A "Suggested testing is usually empty or 1 item" sits beside the required unassessed-stage lines. — impact-negligible → report only

**Nit:**
- F5. `release-check.md:45` (Step 2) does not say which repository to baseline when the target project is not the current checkout, and Step 5 (`:70`) does not map trace points to stages. Two independent readers (builder and reviewer) made the same choices, so no change is needed now. — impact-negligible → report only
- F6. `release-check.md:126-128`: the Output Format always shows the capture line, while `:131` limits it to content or product defects. Saying "omit the capture line when every gap is operator work" would remove the guess. — impact-negligible → report only

### Requirements Checklist

- [x] R1 trace source, package, serving environment, consumer behavior: Step 3 and the Trace table. Delivered.
- [x] R2 consumer behavior, never presence, decode or identifier alone: `:7`, Step 4. Delivered.
- [x] R3 version mixing, stale caches, missing dependencies, interrupted activation as examples: Step 3, marked illustrative. Delivered.
- [x] R4 four stages reported separately, definitions cited, "unassessed": Step 5, Stages table. Delivered (F1 is a gap in the verdict, not in R4).
- [x] R5 current run or historical with date or revision, historical never verifies: `:7`, Step 5. Delivered (exercised on fixture B).
- [x] R6 ready, not ready or unknown with gaps, ready needs current-run live acceptance, operator acts named only: Step 6, `:5`. Delivered.
- [x] R7 toolbox ownership, no queue machinery, no REQ for operator work: blockquote `:3`, `:5`, `:131`. Delivered.
- [x] Integration list: route row and argument-hint, both help menus, `toolbox_actions`, guide, README, ownership blockquote, literal cross-package links (contract passes), per-command help from When to Use and Input, modules.tsv/installer/updater unchanged, brief as data with illustrative examples, prompt-injection before reading and anti-slop before writing. Delivered. Release per prime: integrator's, at finalization.
- [x] Do NOT use when names journey-qa, source-audit and review-work, each with the right reason. Delivered.
- [x] Traps: example brief path starts with `qa/`; no `do-work <retired word>`; help menu column holds; no em-dashes in added lines (`git diff | grep '^+' | grep -c '—'` = 0).

### Acceptance Testing

**Result: Pass** (implementation and integration exercised; deployment from the builder's install leg; live acceptance unassessed)
- Independent cold run of the shipped action on fixture B: GREEN (check 1).
- `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, PASS.
- Routing traced: `release-check` row at `skills/do-work-toolbox/SKILL.md:29` to `./actions/release-check.md`; per-command help served by the router rule at `:40` from When to Use and Input, runs nothing.
- Repository gate (exit 0 at 335c1739) and green probe: integrator's, not re-run. Heavy lanes: being drained elsewhere, not run.

### Suggested Additional Testing

- Deployment: partly assessed. The builder's fresh install of 4427fe9f put both files under `.claude/skills/do-work-toolbox/` with links resolving. Still to run: the update leg after the version bump.
- Live acceptance: unassessed. In a real consumer project, invoke `do-work-toolbox release-check` against a served page and `do-work-toolbox release-check help`, and confirm the agent finds the action and the help prints usage only.
- Edge case: a script-rendered page, to exercise the Step 4 browser-tool path (no one has run it).
- Edge case: a target with no deployment stage (for example a library REQ), to see "not applicable" used and, after F1, the unknown verdict.

### Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All 7 requirements and the Integration list delivered |
| Code Quality | 92% | Clear, follows the journey-qa and source-audit shape; F1 verdict gap, F5 and F6 nits |
| Test Adequacy | 90% | One-off exercises by maintainer choice, independent fixture B GREEN; update leg pending |
| Scope | 100% | 7 declared files, 7 touched |
| Risk | None | Instruction-only, read-only action |
| Acceptance | Pass | Implementation and integration exercised; live acceptance unassessed |

Average 95.5%, recorded as 95% (not rounded up).

### Follow-ups created
None (6 findings report only)

## Review

**Overall: 95%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1. `skills/do-work-toolbox/actions/release-check.md:83-85` and `docs/release-check-guide.md:15-19`: no verdict when no stage failed and live acceptance is "not applicable"; map it to unknown (two-line edit) — impact-negligible → report only. F2. `journey-qa.md:16-19` and `journey-qa-guide.md:5` give no redirect to release-check; the slop-check to source-audit (REQ-676 F1) and ui-review to journey-qa (REQ-677 F4) redirects are still missing too — impact-user-visible → report only. F3. Inherited from REQ-675 F1, still stale: `skills/do-work/crew-members/shared-principles.md:15` "Record Untested" ignores a stage-scoped Pass — impact-rule-change → report only. F4. Inherited from REQ-675 F3 and F4, still stale: `skills/do-work/actions/sample-archived-req.md:99` names no stage; `skills/do-work/actions/review-work.md:460` "usually empty or 1 item" — impact-negligible → report only. F5 (Nit). Step 2 does not say which repository to baseline, and Step 5 does not map trace points to stages; two independent readers converged — impact-negligible → report only. F6 (Nit). Output Format always shows the capture line, while `:131` limits it to content defects — impact-negligible → report only.
**Acceptance:** Pass — implementation and integration: the reviewer followed the shipped action cold on fixture B (served 2.0.3 against package 2.1.0, a 2026-10-03 "deployed and verified" note) and got not ready, deployment failed, live acceptance failed, the note labeled historical and verifying no stage; shipped-package reference contract re-run green; routing and per-command help traced. Deployment only as the builder's fresh install (update leg pending); live acceptance in a real consumer unassessed.
**Restatement sweep:** redefined the list of toolbox actions (one more: release-check). With the inherited REQ-675 stage words and the REQ-676 and REQ-677 action lists, checked the router argument-hint and route table, both help menus, README, `toolbox_actions`, the absent toolbox count, `tutorial.md`, `next-steps.md`, guide lists, crew-member caller lists, the Go registry, the retired-trigger fixture, review-work-guide.md, work.md Step 7, completed-work-presentation-reference.md and the two new files. Lists agree; stale: shared-principles.md:15 (F3), sample-archived-req.md:99 and review-work.md:460 (F4); missing redirects (F2).
**Suggested testing:** 4 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Re-check of 90ef7b4d..44548b59

**Approve.** Overall 96% (Requirements 100%, Code Quality 94%, Test Adequacy 90%, Scope 100%, Risk None). Acceptance: Pass, the same stages as above (implementation and integration exercised, deployment only as the builder's install, live acceptance unassessed).

- Delta `git diff 335c1739..44548b59`: 2 files, 2 lines, exactly the F1 edits (`release-check.md:85`, `release-check-guide.md:19`). No other change.
- F1: resolved. The verdict rule now covers every case: any applicable stage failed gives not ready. Otherwise live acceptance verified gives ready, and unassessed or not applicable gives unknown. That keeps requirement 6 (ready needs current-run live acceptance). The action and the guide use the same wording.
- Checks: `bash _dev/tests/shipped-package-reference-contract.sh` PASS at 44548b59, and the added lines have no em-dashes. The gate and heavy lanes were not run (integrator's).
- New findings: none. F2 to F6 are unchanged and stay report only.
