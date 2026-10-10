# Hand-back: REQ-687 (ai-report --kind proposal and root-cause writes a decision-first brief for unfinished work)

- Branch: `worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds`
- Base commit: `bd56c4b0`
- Commits: `35f984f5` `[REQ-687] add ai-report proposal and root-cause kinds` (one commit)

## File Manifest

- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified): new end-of-file section `## Proposal and Root-Cause Kinds` with target reading, the `proposal` and `root-cause` templates, slug and meta tag, mockup label, capture lines, and what stays the same.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): `:3` clause, `:26` default-kind qualifier, `--kind` Input paragraph, Step 1 branch, Step 2 slug pointer, Step 5 template pointer, Step 8 summary note, one checklist line.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): `:47` sentence gains "Unless `--kind` is used", two Input lines, new end section `## Proposal and Root-Cause Reports`.
- `skills/do-work-toolbox/SKILL.md` (modified): `proposal report`, `root cause report`, `options report` appended to the ai-report routing row.
- `skills/do-work-toolbox/actions/help.md` (modified): new line under the unchanged ai-report line: `ai-report --kind proposal|root-cause <target>  Decision brief for open work`.

## P-A-U

- [x] **[PLAN]:** Prose-only change in five toolbox files, per brief and D-01..D-10. All rules of the two kinds go into one section appended at the end of `ai-report-reference.md` (clear of REQ-656's Report Design Rules edits). `ai-report.md` gets pointers only and stays clear of `:109` (REQ-655) and Step 7 (REQ-657). `completed-work-presentation-reference.md` is not edited (D-03); the deviation is declared in the new section the way `stakeholder-report.md:15` does. New guide/help lines are placed directly below the existing ones; the guide section goes after `## Evidence Safety`. No lesson in `lessons-releases.md` contradicts the plan (no history links or mirrors touched).
- [x] **[APPLY]:** Edits made as planned in the five write-set paths only. Example slug uses `REQ-NNN` rather than a real-looking number so no consumer REQ number lands in shipped prose. The existing `../../do-work/docs/prescribed-shell-primitives.md` pointers at Step 1 and Step 3 are kept.
- [x] **[UNIFY]:**
  - `git diff bd56c4b0 --stat`:
    ```
     skills/do-work-toolbox/SKILL.md                    |  2 +-
     .../do-work-toolbox/actions/ai-report-reference.md | 26 ++++++++++++++++++++++
     skills/do-work-toolbox/actions/ai-report.md        | 15 +++++++++----
     skills/do-work-toolbox/actions/help.md             |  1 +
     skills/do-work-toolbox/docs/ai-report-guide.md     |  8 ++++++-
     5 files changed, 46 insertions(+), 6 deletions(-)
    ```
  - `git diff --check`: clean (exit 0).
  - `bash .../REQ-687-probe.sh`: exit 0, `REQ-687 probe passed.`, wall 2 s.
  - `bash _dev/tests/contract-regressions.sh`: exit 0, `Contract regression checks passed.`, wall 36 s (single run, no budget failure).
  - `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, wall 2 s (also run inside the probe).
  - Files checked: all five diffs read in full. Checked: only the write set changed; `completed-work-presentation-reference.md:20` and `:30` byte-identical (probe lines pass); `ai-report.md:26` still carries "unfinished or unsuccessful"; no citation of `CLAUDE.md` or `_dev/`; no debug artifacts; no `do-work/` path staged.

## Proof Record

**RED (base bd56c4b0, before any edit):** probe exit 1, `REQ-687 probe: 21 failure(s)`; the two completed-work-gate checks passed. Samples:
- `FAIL: skills/do-work-toolbox/actions/ai-report.md lacks: --kind`
- `FAIL: skills/do-work-toolbox/actions/ai-report-reference.md lacks: ## Proposal and Root-Cause Kinds`

At base, `ai-report.md:26`: "- The target is unfinished or unsuccessful; report its status instead of presenting it as shipped." and `:30`: "`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank. One invocation covers one UR or one REQ; blank is the explicit `most recent` form." No path to a proposal report.

**GREEN (after 35f984f5):** probe exit 0, `REQ-687 probe passed.`

**GREEN-when trace** (line numbers at 35f984f5):

| Clause | Where |
|---|---|
| Bundle `yyyy-mm-dd_hhmm_proposal-<topic>` | `actions/ai-report-reference.md:139`; pointer `actions/ai-report.md:55` |
| Decision first | `actions/ai-report-reference.md:129` |
| 2 to 4 options, each with benefit, risk, cost | `actions/ai-report-reference.md:130` |
| Smallest-change option always listed | `actions/ai-report-reference.md:130` (proposal), `:137` (root-cause) |
| Recommendation with reason | `actions/ai-report-reference.md:131` |
| Evidence ledger | `actions/ai-report-reference.md:132` |
| Limits | `actions/ai-report-reference.md:133` |
| Q1..Qn | `actions/ai-report-reference.md:134` |
| One capture line per option, never run | `actions/ai-report-reference.md:135`, `:143` |
| `<head>` carries `ai-report-kind` meta | `actions/ai-report-reference.md:139` |
| Every mockup carries "MOCKUP — proposal" | `actions/ai-report-reference.md:141` |
| root-cause shape | `actions/ai-report-reference.md:137` |
| Safety load order kept, target resolution lifted | `actions/ai-report-reference.md:125`; `actions/ai-report.md:38` |
| Default kind still stops on an unfinished REQ | `actions/ai-report.md:26` (qualified "and no `--kind` is given"), `completed-work-presentation-reference.md:20`, `:30` unchanged |
| Guide, routing, help | `docs/ai-report-guide.md:47`, `:70-71`, `:82-84`; `SKILL.md:24`; `actions/help.md:15` |

## Decisions

- **D-11** (DECIDE & STATE) A completed REQ or UR passed with `--kind` is not in the accepted-target list, and no extra rule says what happens. Adding one was not asked for (anti-bloat). If it comes up, a one-line rule can be added later.
- **D-12** (DECIDE & STATE) The `--kind` Input sentence went into a new paragraph under `ai-report.md:30` instead of onto that line, because REQ-655 edits `:30`. The merge stays clean.
- **D-13** (DECIDE & STATE) The help line cannot fit the 31-character command column, so it uses two spaces before its description. Same for the two guide Input lines.
- **D-14** (DECIDE & STATE) The slug example uses `REQ-NNN` (`proposal-REQ-NNN-retry-policy`) so no consumer REQ number appears in shipped prose.

## Discovered Tasks

- impact-minor: `skills/do-work-toolbox/actions/architecture-report.md:129` still says "`ai-report` takes a UR or REQ and presents completed work", which is now incomplete (D-08) → report only
- impact-minor: `skills/do-work-toolbox/docs/ai-report-guide.md:3` intro and `skills/do-work-toolbox/actions/ai-report-reference.md:3` companion blockquote still describe completed work only. Both remain true for the default kind, and the brief did not name them → report only

## Lessons Read

- `_dev/primes/lessons-releases.md` (whole; required_lessons). Families: `canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`. Neither applies: no history links or manifests changed.
- `_dev/primes/prime-action-files.md` § Traps `alternate-writer-contract-drift`: the restatement sweep in Exploration still holds. Only `architecture-report.md:129` became incomplete.
- `_dev/primes/lessons-action-files.md`: not read (dropped for budget, per brief).

## Anti-Bloat Check

`git diff --stat` above: five files, 46 insertions, 6 deletions. Items added that the REQ or brief did not name:
- The guide section title `## Proposal and Root-Cause Reports`. Reason: the brief asked for "one short section" without a name.
- In the mockups paragraph, how the label is placed (text inside an SVG/HTML mockup, a badge inside a raster frame). Reason: this is D-06, so it is named.
- Nothing else. No new files, kinds, frontmatter fields, statuses, image pipeline, `options` alias, default-kind meta, or tests.

## Proposed CHANGELOG Entry

**AI Report Proposal and Root-Cause Briefs**

`ai-report` can now write a decision brief for work that is not finished, so a proposal or a "why did this fail" report has a fixed shape instead of being made up each time. Without `--kind`, the report works as before.

- `do-work-toolbox ai-report --kind proposal <topic|REQ-NNN|UR-NNN>` writes a brief that opens with the decision. It then gives two to four options with benefit, risk and cost, a recommendation, an evidence ledger, limits, open questions Q1..Qn and one `do-work capture-request:` line per option.
- `--kind root-cause` uses "what happened, why, what to change" in place of the options, and also accepts a failed or cancelled REQ.
- The options always include the smallest-change option, such as deleting the mechanism or doing nothing extra.
- Bundles are named `yyyy-mm-dd_hhmm_<kind>-<topic>` and carry `<meta name="ai-report-kind">`. Mockups are labelled "MOCKUP — proposal" and kept in `generated/`.
- Routing adds `proposal report`, `root cause report` and `options report`. The toolbox help lists the new form.

## Proposed Lesson Bullet

None. The change followed existing precedent (`stakeholder-report.md:15` deviation pattern) and nothing went wrong that a later builder would need to know.

## Integration Seams

- `SKILL.md:24`: one table row that REQ-655 and REQ-657 also extend. A conflict is expected. Keep every sibling's phrases and add `proposal report`, `root cause report`, `options report` at the end.
- `actions/help.md:15`: new line directly under the ai-report line. Siblings may add lines in the same place, so keep all of them.
- `docs/ai-report-guide.md`: two Input lines directly under `most recent`, and a new section after `## Evidence Safety` at the end of the file. Sibling lines in the same block or end-of-file sections may conflict as adjacent additions, so keep both.
- `actions/ai-report.md`: `:30` is unchanged (my paragraph is below it). `:109` and Step 7 are not touched.
- Check wall times: probe 2 s, contract-regressions 36 s, shipped-package-reference-contract 2 s.
