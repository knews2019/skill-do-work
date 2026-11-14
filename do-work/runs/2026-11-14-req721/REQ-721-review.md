# Review: REQ-721 (every ai-report ends with a Tinka canvas suggestion, with pros and cons)

**Approve with follow-ups** — the closing band, the design rule, the action sentences and the guide section are all shipped and consistent, but the hand recipe has two defects an agent would hit when it follows the text literally.
Route B | merge range `192d50fa..2285c121` (builder commit `cd48924d`)

Reviewer: independent, orchestrated mode, standard depth. Nothing was edited except this file.

## What's built

- Every report built from the kit now ends with a third block in the band "Try another layout or teaching method": a short description of a Tinka canvas, a Pros line, a Cons line and one line an agent can follow.
- The guide has a new section "Review on a Tinka Canvas" after "A Second Opinion": six rules and a 12-step hand recipe.
- Still missing: the recipe text does not match what the builder actually ran for tile ids (F1), and it downloads Tinka where the request asked it to stop (F2).

## Decisions / risks for you

- **F2 needs a ruling.** The request says the recipe's first step checks for the pinned checkout and stops with one line. The shipped step 3 clones Tinka from GitHub when it is missing. Value of the shipped form: the line from the band works in a fresh project with no manual setup. Risk: one pasted line makes an agent download and run third-party code, and the choice is not written in the Decisions section. The commit pin limits the risk, because Git checks content by hash.
- **Risk level: Low.** Tinka listens on `127.0.0.1` only, the manual layout never writes `AGENTS.md`, and all Tinka state stays inside `.tinka/`. I confirmed these three points in the Tinka source at the pinned commit, read as data.

## Findings

**Important:**

- **F1** — `skills/do-work-toolbox/docs/ai-report-guide.md:200` (commands at `:203-204`): the tile-id rule makes a page tile and its reply tile share one id. The text says "Build each tile id from the bundle folder name and the tree id", and both commands use `<tile-id>`. For `v1c01` the page and the reply get the same id. Tinka answers "Item already exists." for the second one (`lib/workspace.mjs`, `makeItem`), and the next sentence tells the agent this answer is fine. The reply tile is then silently missing. The builder's own run did not follow the text: its ids on disk are `<bundle>_<tree-id>_page` and `<bundle>_<tree-id>_reply` (`.tinka/boards/.tinka/items/`). Fix: write "Build each tile id as `<bundle>_<tree-id>_page` or `<bundle>_<tree-id>_reply`" and put those two forms in the two commands. — impact-user-visible → report only
- **F2** — `skills/do-work-toolbox/docs/ai-report-guide.md:161-168` against REQ line 65 ("The recipe's first step checks for the pinned checkout and stops with one line that says what is missing"): the check is step 3, after step 2 has already created `.tinka/` and edited the local exclude file, and a missing checkout is cloned, not reported. The recipe stops only when Node is missing or the clone fails. Decision D-03 records the folder and the pin, not the automatic download. Fix: move the Node and checkout check before any write; print one line before cloning ("Tinka is not in `.tinka/src`; cloning it at commit 77ba29c…"); add "When `rev-parse HEAD` prints another commit, stop"; record the choice as a decision. The literal alternative is to stop and print the two clone commands. — impact-user-visible → report only

**Minor:**

- **F3** — `skills/do-work-toolbox/docs/ai-report-guide.md:179`: "Its captures are written whatever the verdict is" is wrong for `skipped` and `error`. With no browser engine the judge returns before it renders anything (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go:140-144`), so step 6 finds no `wide-light.png` and the recipe has no rule for that. Fix: "The judge writes its captures on `pass` and on `fail`. On `skipped` or `error` there is no picture: stop with one line that gives the judge's finding." — impact-user-visible → report only
- **F4** — `skills/do-work-toolbox/docs/ai-report-guide.md:182,185,191,196`: three placeholders have no rule. `<lineage>` has no naming rule, so a later session can pick another id and create a second canvas (the builder invented `tinka-review-surface-1547`). `<tree-id>` is not defined in this section, and a plain second-opinion page (`index-<cli>-<model>.html`, guide `:122`) has none. `ai-reports/<bundle>` is wrong for a bundle inside a month folder (guide `:106`). Fix: add one paragraph before the steps: `<bundle>` is the bundle's folder name and its path may include a month folder; `<lineage>` is the oldest bundle's folder name; `<tree-id>` is the page's iteration id (`v1` for `index.html`), or for a page without one its file name without `index-` and `.html`. — impact-user-visible → report only
- **F5** — `skills/do-work-toolbox/docs/ai-report-guide.md:175,182,188`: every command uses a relative `.tinka/…` path or `$PWD`, but no rule says to run from the project root. A shell in a subfolder creates a second `.tinka/` there. Fix: add the rule "Run every command from the project root." — impact-user-visible → report only
- **F6** — `skills/do-work-toolbox/docs/ai-report-guide.md:163-168`: one block holds the check and the clone. On a later run a literal agent runs `git clone` again, it fails because the folder exists, and the prose says a failed clone means stop. Fix: split into a check block (`node --version`, `git -C .tinka/src rev-parse HEAD`) and a clone block marked "only when the checkout is missing". — impact-user-visible → report only
- **F7** — `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:200`: the `FILL:` placeholder sits in the middle of the line and its end is not marked. The fixed tail ("and its revisions on a Tinka canvas, following Review on a Tinka Canvas in the ai-report guide") can be replaced by mistake, and then the line no longer names the guide section. The two existing lines (`:189`, `:192`) end with the placeholder. Fix: `ai-report: follow Review on a Tinka Canvas in the ai-report guide for this report and its revisions: ai-reports/FILL: this bundle folder`. — impact-negligible → report only
- **F8** — `do-work/runs/2026-11-14-req721/REQ-721-handback.md:13-19`: the proof holds no newly created report. The GREEN condition "A newly created report ends with the canvas suggestion: its pros, its cons and the line" rests on two greps of the skeleton. The recipe run covered one bundle with tree pages only: no revision (the "one row per bundle" rule) and no plain second-opinion page. — impact-negligible → report only
- **F9** — `skills/do-work-toolbox/docs/ai-report-guide.md:134` against `skills/do-work-toolbox/SKILL.md:25`: the guide says you can "ask for a canvas in your own words", but the router row for `ai-report` has no canvas trigger, and decision D-02 chose no router change. A request such as "put this on a Tinka canvas" reaches the action only when it also says `ai-report`. Fix: add `tinka canvas` to the trigger list in `SKILL.md:25`, or delete the clause in the guide. — impact-user-visible → report only

**Nit:**

- **F10** — `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:194`: "stays on its page" can be read as "is written into the report page". The guide says "stays on its tile" (`ai-report-guide.md:132`). Fix: "stays on the canvas beside its page". — impact-negligible → report only
- **F11** — `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:193` with `presentation.css:79`: the new `<h3>` has no top margin and touches the code box above it, in both widths of my render. Fix: `#another-layout h3 { margin-top: 1.2em; }` next to `presentation.css:367`. — impact-negligible → report only
- **F12** — three small wording mismatches: `skills/do-work-toolbox/actions/ai-report.md:25` says "skips Steps 1–8" where lines `:49` and `:53` say "Steps 1 to 8 do not apply"; `ai-report-guide.md:144` says every file goes into `.tinka/` while step 2 (`:155`) also edits the local exclude file; `ai-report-reference.md:92` says Step 8 prints "the offer's line" while `ai-report.md:195` says "the Tinka canvas offer with its line". Fix: use the second wording in each pair. — impact-negligible → report only

## Requirements checklist

- [x] Third block in the band, after the layout and the teaching method — delivered (`skeleton.html:193-200`)
- [x] Every report, no `.tinka/` condition — delivered (no condition anywhere in the shipped text)
- [x] Two or three pros and cons from facts counted without Tinka — delivered (`skeleton.html:196-197`, `ai-report-reference.md:92` names all three facts)
- [x] One line that names the bundle and the guide section — delivered (`skeleton.html:200`; see F7)
- [x] Guide section "Review on a Tinka Canvas" after "A Second Opinion" — delivered (`ai-report-guide.md:130`)
- [x] Pinned source checkout — delivered (`:166`)
- [x] `.tinka/` created, hidden through the local exclude file, never `.gitignore` — delivered (`:155-159`)
- [x] Manual layout with `init --config` and `--boards` inside `.tinka/`; update policy off — delivered (`:175-176`)
- [x] One PNG per page with `ai-report-judge --out`, scratch copy for a second-opinion page — delivered (`:179-189`; see F3)
- [x] Copies in one folder, registered as the only resource root — delivered (`:191-198`)
- [ ] One image tile per capture and one Markdown tile per reply — partially delivered (F1)
- [x] Feedback read with `board.get` and `feedback.list`, handed over as quoted data — delivered (`:213-221`)
- [x] Step 8 summary names the canvas offer with its line — delivered (`ai-report.md:195`)
- [x] Declared files changed; no Go code, no new command — delivered (5 files, all in the write set)
- [x] "Check first" test of a resource root inside `.tinka/` — delivered (recorded in the REQ's Exploration; `workspace.json` on disk holds `../../captures`)
- [x] No bundle changed; `git status --porcelain` clean of `.tinka/`; `git diff --quiet HEAD -- AGENTS.md` passes — delivered (both checked now)
- [x] No `project init`; no HTML or bundle folder as a resource root — delivered (`:142-143`; `.tinka/captures` holds 3 `.png` and 2 `.md`)
- [x] Creating a report never starts Tinka; unattended runs print text only — delivered (`ai-report-reference.md:92`, `ai-report-guide.md:140`)
- [ ] Absent Tinka: the first step checks for the checkout and stops with one line — partially delivered (F2)
- [x] Recipe followed once on a real lineage — delivered (artifacts on disk match the hand-back; see F8 for its limits)
- N/A `CHANGELOG.md` and `VERSION` — written by finalization, not in this range

## Acceptance testing

**Result: Partial** — stages covered: Implementation and Integration.

- Ran `bash _dev/tests/action-shell-blocks.sh` (97 blocks, pass), `bash _dev/tests/shipped-package-reference-contract.sh` (pass) and `bash _dev/tests/contract-regressions.sh` (pass, exit 0).
- Ran the two proof greps: 6 and 2.
- Rendered the kit skeleton and a band-only copy with `ai-report-judge` in the scratchpad: verdict `pass`, no overflow at 1440 and 390 pixels, light and dark. The new block reads correctly (see F11).
- Read the builder's leftover `.tinka/` folder and the Tinka source at 77ba29c as data. Confirmed: checkout at the pinned commit, `init` refuses an existing config, `start` runs in the background and listens on `127.0.0.1`, the manual layout never touches `AGENTS.md`, both "already exists" answers, a relative `--out` resolves against `--repo-root`.
- Partial, not Pass, because the recipe text and the builder's real run differ on tile ids (F1).

## Restatement sweep

The diff redefines what the closing band holds and what the Step 8 summary prints. I searched `skills/` and `suite/` for "two suggestions", "both suggestions", "teaching method", "another layout", "closing band", `another-layout`, "suggestion" and "Tinka". Every restatement agrees with the new meaning: `ai-report-reference.md:91-92`, `ai-report.md:25,165,167,195,221`, `ai-report-guide.md:38,134`, `presentation.css:17-19`, `skeleton.html:29,177-184`. The words "two suggestions" and "both suggestions" stay true because the canvas is called an "offer" everywhere. No Go code and no test reads the band. No added line cites `CLAUDE.md`, `AGENTS.md` of this repository or `_dev/`. The two relative links and the anchor `#review-on-a-tinka-canvas` resolve. One cosmetic leftover: the comment at `presentation.css:366` says "closing suggestions" and its rule now also styles the offer's line.

## Suggested additional testing

- Deployment: unassessed. After finalization, confirm `CHANGELOG.md`, the mirror and `VERSION` carry this change.
- Live acceptance: unassessed. Create one ordinary report in a consumer project and read its band: counted pros and cons, a filled line, and the offer in the Step 8 summary.
- Give the guide section alone to a fresh agent on a lineage with one revision and one plain second-opinion page, and count the reply tiles afterwards.
- Run the recipe on a machine with no browser engine (verdict `skipped`).
- Run the recipe from a subfolder and from a linked worktree.
- Manual: open the canvas, leave a note and a star, and confirm step 10 reads both back.

## What I verified and what I did not

- Verified by reading: the whole diff, the four shipped files at the merge commit, the REQ, the UR, the hand-back, the judge's Go code, and the parts of Tinka's source the recipe depends on.
- Verified by running: the three scripts above, the two greps, two renders, `git status --porcelain` and the `AGENTS.md` diff check.
- Not verified: I did not run the recipe myself, did not start Tinka, and did not clone from the network. The hand-back's claims of HTTP 200 for five files, and of a note and a star written and read back, rest on the builder's word. The "identical before and after" Git status is confirmed only for the state now.

## Scores (on the record — not the headline)

**Overall: 74%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 90% | 18 of 20 fully delivered; F1 and F2 partial |
| Code Quality | 80% | Follows the kit and action patterns; recipe has one false sentence and several unmarked choices |
| Test Adequacy | 75% | Lint and contract checks pass, recipe run once; no new report created, narrow lineage |
| Scope | 90% | Five declared files only; one behavior choice missing from Decisions |
| Risk | Low | Pinned third-party code, local-only server |
| Acceptance | Partial | Band and rules work; recipe text differs from the run that proved it |

Average of the four percentages is 84%. Acceptance Partial takes 10 points.

## Follow-ups created

None (12 findings report only)

---

## Ready-to-paste section for the REQ file

```markdown
## Review

**Overall: 74%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 80% |
| Test Adequacy | 75% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1: `skills/do-work-toolbox/docs/ai-report-guide.md:200` builds a page tile and its reply tile from the same two parts, so both get one id; Tinka answers "Item already exists." and the text calls that fine, so the reply tile is silently missing. The builder's own run used `<bundle>_<tree-id>_page` and `_reply`. Fix: prescribe those two id forms. — impact-user-visible → report only
- F2: `skills/do-work-toolbox/docs/ai-report-guide.md:161-168` clones Tinka when the checkout is missing, and checks only in step 3 after two writes; the REQ constraint asks the first step to check and stop with one line. The choice is not in Decisions. Fix: check before any write, print one line before cloning, stop when `rev-parse HEAD` is not the pin, record the decision. — impact-user-visible → report only

**Minor findings:**
- F3: `skills/do-work-toolbox/docs/ai-report-guide.md:179` says captures are written whatever the verdict is; on `skipped` or `error` none exist (`report_judge.go:140-144`) and the recipe has no stop rule for it. — impact-user-visible → report only
- F4: `skills/do-work-toolbox/docs/ai-report-guide.md:182,185,191,196` leave `<lineage>` unnamed, `<tree-id>` undefined for a plain second-opinion page, and `ai-reports/<bundle>` wrong for a month folder. — impact-user-visible → report only
- F5: `skills/do-work-toolbox/docs/ai-report-guide.md:175,182,188` use relative paths and `$PWD` with no rule to run from the project root. — impact-user-visible → report only
- F6: `skills/do-work-toolbox/docs/ai-report-guide.md:163-168` puts the check and the clone in one block, so a later literal run re-clones, fails and stops. — impact-user-visible → report only
- F7: `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:200` has its `FILL:` placeholder mid-line with no marked end, so the tail that names the guide section can be lost. — impact-negligible → report only
- F8: `do-work/runs/2026-11-14-req721/REQ-721-handback.md:13-19` holds no newly created report, and the recipe run had no revision and no plain second-opinion page. — impact-negligible → report only
- F9: `skills/do-work-toolbox/docs/ai-report-guide.md:134` promises an own-words request, but `skills/do-work-toolbox/SKILL.md:25` has no canvas trigger. — impact-user-visible → report only
- F10 (Nit): `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:194` says "stays on its page" where the guide says "stays on its tile". — impact-negligible → report only
- F11 (Nit): `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:193` heading touches the code box above it; add a top margin beside `presentation.css:367`. — impact-negligible → report only
- F12 (Nit): wording pairs disagree at `ai-report.md:25` against `:49`, `ai-report-guide.md:144` against `:155`, and `ai-report-reference.md:92` against `ai-report.md:195`. — impact-negligible → report only
**Acceptance:** Partial — Implementation and Integration covered: lint, contract checks and a skeleton render pass, and the band and rules are right; the recipe text differs from the builder's real run on tile ids (F1). Deployment and live acceptance unassessed.
**Restatement sweep:** redefined what the closing band holds (two suggestions plus one Tinka canvas offer) and what the Step 8 summary prints; every restatement in `skills/` and `suite/` agrees, one cosmetic comment at `presentation.css:366`
**Suggested testing:** 6 items
**Follow-ups created:** None (12 findings report only)

*Reviewed by review-work action*
```
