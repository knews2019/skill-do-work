# Hand-back: REQ-656 (ai-report revise writes a new sibling bundle that supersedes the prior one)

- Branch: `worktree-agent-REQ-656-ai-report-revise-sibling-bundle`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-656-ai-report-revise-sibling-bundle`
- Base: `7fbeff30` (recorded before the first edit)
- Commits: `d281018f` `[REQ-656] add ai-report revise form: new -rev<N> sibling bundle that supersedes the prior one` (one commit)

## File Manifest

- `skills/do-work-toolbox/actions/ai-report.md` (modified): Use-when bullet for updating a report; new `### Revise form` subsection (7 numbered steps, no bash fence) after the catalog forms; path-last `file://` sentence in `## Output Format`; D-15 widening on the Step 6 paragraph (two sentences) and the last checklist line; one new checklist line at the end.
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified): one length-keyed table-of-contents bullet at the end of `## Report Design Rules (Step 5)`.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): one Input fence line for `revise`; `:21` widened (D-15); new `## Revising a Report` section at the end (6 sentences).
- `skills/do-work-toolbox/SKILL.md` (modified): `revise the report`, `update the report` appended to the ai-report routing row.
- `skills/do-work-toolbox/actions/help.md` (modified): one line `ai-report revise <dir|latest>` after the index/find line, description at column 34.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Prose-only change in five toolbox files, following the brief's decided shape (D-01 to D-15). Reuse existing contracts by name (catalog command fence, Safety Load Order, Terminal-Success Target Resolution, Evidence Honesty, Collision-Safe Publication, `architecture-report.md` Step 3, Report Design Rules) instead of restating them. No Go, no tests, no new bash fence. Proof: RED probe at base, GREEN walkthrough in a `mktemp -d` scratch repo for rev1 and rev2, GREEN probe, reference contract, contract-regressions.
- [x] **[APPLY]:** Edits exactly as planned, inside the five-file write set; each edit kept local to its anchor so the REQ-657 and REQ-687 merges stay small.
- [x] **[UNIFY]:** `git diff 7fbeff30 --stat`:
  ```
   skills/do-work-toolbox/SKILL.md                      |  2 +-
   .../do-work-toolbox/actions/ai-report-reference.md   |  1 +
   skills/do-work-toolbox/actions/ai-report.md          | 20 +++++++++++++++++---
   skills/do-work-toolbox/actions/help.md               |  1 +
   skills/do-work-toolbox/docs/ai-report-guide.md       |  7 ++++++-
   5 files changed, 26 insertions(+), 5 deletions(-)
  ```
  Checks (run from the worktree root):
  - `bash .../REQ-656-probe.sh`: exit 0, 1.2 s.
  - `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, 1.1 s.
  - `bash _dev/tests/contract-regressions.sh`: exit 0, 24 s.
  - `git diff 7fbeff30 --check`: exit 0 (clean).
  - `bash _dev/tests/staged-skills-contract.sh`: exit 2, refuses with "heavy-only; run maintainer-verify.sh --heavy". Not run further: the brief forbids the gate. The example paths I added all start with `ai-reports/`, which is the brief's safe form. Integrator's gate covers it.
  Files checked: all five, by reading the full diff. Each probe token is inside `### Revise form`; no `#` line inside the subsection; no citation of `CLAUDE.md` or `_dev/`; `completed-work-presentation-reference.md` and `architecture-report.md` untouched (probe confirms the immutability and sibling-suffix lines are byte-identical).

## Proof Record

- **RED at base 7fbeff30:** probe exit 1, **23 failures** (no revise form or its tokens, no `file://` in Output Format, no TOC rule, no guide/routing/help lines). Matches the pre-dispatch count.
- **GREEN walkthrough rev1** (scratch repo from `mktemp -d`, one committed bundle `ai-reports/2026-09-11_1430_deploy-guide/index.html` with a `<title>` and 7 `<h2>` sections; prose followed literally for `ai-report revise latest "sitemaps re-checked"`; catalog command `bash <worktree>/skills/do-work/tools/do-work-cli.sh --repo-root <scratch> --format text ai-report-index`):
  - Step 1 catalog: exit 0, `1 bundles`; `latest` = `ai-reports/2026-09-11_1430_deploy-guide`, `superseded_by` empty.
  - Step 4: hops 0, so N = 1; slug `deploy-guide`; new folder `2026-10-10_1903_deploy-guide-rev1` (no collision).
  - Step 5: rev block is the first content after the `<h1>` title: `<h2>rev-1 (2026-10-10)</h2>`, **Changed** list (`sitemaps re-checked`), **Still to do** list, relative link `../2026-09-11_1430_deploy-guide/index.html`. `<head>` holds `<meta name="ai-report-supersedes" content="2026-09-11_1430_deploy-guide">`. The prior bundle has no `ai-report-kind` meta, so none was copied. 7 sections > 6, so a `<nav>` TOC: 7 `href="#..."` links, 0 missing ids (checked with a short Python script).
  - Step 6 catalog: exit 0, `2 bundles`; prior bundle `superseded_by: ai-reports/2026-10-10_1903_deploy-guide-rev1`.
  - `git -C <scratch> diff --exit-code -- ai-reports/2026-09-11_1430_deploy-guide`: exit 0.
  - `git -C <scratch> status --short`: only `?? ai-reports/2026-10-10_1903_deploy-guide-rev1/`, `?? ai-reports/catalog.json`, `?? ai-reports/index.html`.
  - Last output line: `ai-reports/2026-10-10_1903_deploy-guide-rev1/index.html file:///var/folders/.../tmp.yn7IhVRjFQ/ai-reports/2026-10-10_1903_deploy-guide-rev1/index.html`.
- **GREEN walkthrough rev2** (after committing rev1, `ai-report revise latest "second pass"`):
  - `latest` = `ai-reports/2026-10-10_1903_deploy-guide-rev1`, 1 hop, N = 2, slug `deploy-guide` (date, time and `-rev1` stripped), new folder `2026-10-10_1903_deploy-guide-rev2`.
  - TOC: 8 links, 0 missing ids. Prior bundles unchanged (`git diff --exit-code` exit 0 on the original and rev1 `index.html`). Status: only the new folder plus modified `catalog.json` and `index.html`.
  - Catalog chain: original `superseded_by` rev1, rev1 `supersedes` original and `superseded_by` rev2, rev2 `supersedes` rev1. Chain original → rev1 → rev2 confirmed.
  - Finding: rev1 and rev2 were written in the same minute, and the catalog listed rev1 **before** rev2 (same date, tie broken by path ascending). A third `latest` would resolve to rev1. Step 1's chain-forward rule fixes this, so I made it apply explicitly to `latest` as well (D-17).
- Step 7 render check: **not run** (optional for this fixture; no browser automation used).
- Scratch repo and helper scripts removed after the run.
- **GREEN probe:** exit 0, "REQ-656 probe passed.", 1.2 s.

## Decisions

- **D-16 (DECIDE & STATE)** N counting when the catalog command refused and the user gave an explicit `<dir>`: step 4 says to read the hop count from each bundle's `ai-report-supersedes` meta instead of the catalog. The brief only covered `latest` under refusal. Without this fallback an explicit-folder revise would have no defined N. Reversible: one parenthetical.
- **D-17 (DECIDE & STATE)** Step 1's chain-forward rule says "the resolved bundle (from `latest` or `<dir>`)" instead of "the named bundle". Reason: the walkthrough showed a same-minute tie sorts the older revision first, so `latest` alone can pick a superseded bundle. Following `superseded_by` forward gives the right bundle with no Go change.
- **D-18 (DECIDE & STATE)** The Output Format path-last rule is written as one sentence with a semicolon joining the HTTP note, to keep it as "one sentence" per the brief.
- **D-19 (DECIDE & STATE)** Guide Input line: the command is 58 characters, longer than the column-43 description slot, so the description follows two spaces (per the brief).
- **D-20 (DECIDE & STATE)** No example relative link like `../<prior folder>/index.html` in shipped prose. "a relative link to the prior bundle that works from the new folder" is enough, and a `../` token risks the citation scan.
- **D-15 sweep verdicts:**
  - `ai-report.md` Step 6, "creates only the report bundle": false for revise (it also regenerates the catalog). Widened to "..., plus the regenerated catalog for the revise form."
  - Same paragraph, "the catalog forms above are the one exception to searching": false for revise (it reads the catalog to resolve `latest` and the chain). Widened to "the catalog forms and the revise form above are the only exceptions to searching".
  - `ai-report.md` Output Format, "The timestamped folder is the only stakeholder artifact this action publishes": still true. The catalog is derived navigation data, not a stakeholder artifact, and the catalog forms already wrote it under this same sentence. Unchanged. The new path-last sentence follows it in the same paragraph.
  - `ai-report.md` last checklist line, "(the catalog forms' `catalog.json` and `index.html` excepted)": false for revise. Widened to "(the catalog forms' and the revise form's ...)".
  - `ai-report-guide.md:21`, "This action produces only the report bundle": false for revise. Widened to "... (a revise also regenerates the report catalog)."

## Discovered Tasks

- impact-minor: `ai-report-index` breaks same-date ties by path ascending, so of two bundles written in the same minute the older `-rev1` sorts ahead of `-rev2` and becomes `latest`. The revise prose works around it by following `superseded_by` (D-17). The Go sort could instead put a bundle that supersedes another first on a tie, or sort ties by path descending → report only

## Lessons Read

- `_dev/primes/lessons-releases.md` (whole file, the REQ's `required_lessons`): family `canonical-link-outlives-its-target` applies to the integrator's UR-144 close (see the note below). No link to a UR-144 record was added by this change.
- `_dev/primes/prime-action-files.md` (whole) and its Traps family `alternate-writer-contract-drift` (it drove the D-15 sweep).
- `_dev/primes/prime-releases.md` (read; the release is the integrator's).
- `_dev/primes/lessons-action-files.md`: not read (dropped for budget). The `example-path-read-as-citation` rule was taken from the brief: every example path starts with `ai-reports/`.
- Crew members: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `anti-slop.md`.

## Anti-bloat Check

`git diff 7fbeff30 --stat`: 5 files, 26 insertions, 5 deletions (table above). Items the brief did not name:
- The step-4 fallback parenthetical for N when the catalog refused (D-16). Kept: without it an explicit-folder revise under a refused catalog has no N.
- "(from `latest` or `<dir>`)" in step 1 (D-17). Kept: it fixes the same-minute tie the walkthrough found.
- "(step 4 above names the folder)" in step 5. Kept: without it Step 2's slug rule would compete with step 4's naming.
No new heading beyond `### Revise form` and `## Revising a Report`. No new verb, flag, script, Go code, test, Just recipe, queue field or status.

## Proposed CHANGELOG Entry

**ai-report Revise Writes a New Report Revision and Keeps the Old One**

You can now update an existing report without breaking the rule that reports never change. `ai-report revise <dir|latest> [what changed]` writes a new report next to the old one and links the two in the catalog.

- New `ai-report revise <dir|latest> [what changed]` form: re-checks the old report's claims, then writes a fresh `-rev<N>` sibling bundle through the normal collision-safe path. The new report opens with a "rev-N (date)" block (Changed, Still to do) and names the old one in an `ai-report-supersedes` meta.
- The catalog is regenerated, so the old report shows "superseded by" the new one. The old bundle's bytes never change.
- Report Design Rules: any report with more than six top-level sections or about 1,500 words gets a table of contents.
- Every run that writes a report now ends with the report's path and a `file://` link.
- Guide, help menu and routing phrases (`revise the report`, `update the report`) name the new form.

## Proposed Lesson Bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (it already holds the REQ-655 catalog lessons).

- [family: catalog-tie-order] [REQ-656: "newest first" by folder-name date is only minute-precise, so two revisions written in the same minute tie and sort by path ascending, which puts the older one first; any prose that treats the first catalog entry as the newest must also follow `superseded_by` forward](../../../../do-work/archive/UR-144/REQ-656-ai-report-revise-supersedes-sibling-bundle.md)

## Integration Seams

- On base 7fbeff30: **no** REQ-657 (ai-report judge) edits and **no** REQ-687 (`--kind proposal|root-cause`) edits were present (no `ai-report judge`, no `--kind` in the toolbox).
- Expected conflicts at merge, all resolved by keeping every side's text:
  - `ai-report.md` checklist end: I modified the last checklist line (D-15) and appended one line after it. REQ-687 adds a line after the render-judged line and REQ-657 edits that line, so these hunks are adjacent and git will likely show a conflict.
  - `SKILL.md` ai-report row: REQ-687 appends three phrases to the same row. Union both.
  - Guide end of file: REQ-687 appends `## Proposal and Root-Cause Reports` after my `## Revising a Report`. Union both. The guide Input fence and the `help.md` list each get one line from each sibling next to mine.
  - `ai-report.md` Step 6 paragraph (the D-15 widening): neither sibling touches it, per the brief.
- Test wall times: probe 1.2 s, reference contract 1.1 s, contract-regressions 24 s (single runs, no reruns needed).

## Note for the Integrator

REQ-656 closes UR-144 when it finalizes (REQ-654 cancelled; REQ-655 and REQ-657 archived by then). Closing the UR moves its records into `do-work/archive/UR-144/`, so any shipped link to a UR-144 record path must be swept in the same change (`lessons-releases.md`, family `canonical-link-outlives-its-target`). Known today: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:5` links `blob/main/do-work/archive/REQ-655-ai-report-index-and-find-catalog.md`, which becomes `do-work/archive/UR-144/REQ-655-...` at that close. I did not edit that file. The proposed lesson link above already uses the `UR-144/` path. Fix the relative depth when you place it.
