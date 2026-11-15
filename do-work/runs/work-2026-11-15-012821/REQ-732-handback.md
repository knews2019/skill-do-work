# REQ-732 hand-back (allow docs/run-prompts/ in the staged-skills contract and move the run prompts there)

Branch: `worktree-agent-REQ-732-allow-docs-run-prompts`, tip `342a36a6e6fda972312774f1900d18c0600a4de0`.

## Commits
- `342a36a6` [REQ-732] allow docs/run-prompts/ in the staged-skills contract and move the run prompts there (one commit; amended once on the unpushed branch to reflow two comment lines)

## Files
- modified `_dev/tests/staged-skills-contract.sh` (docs comment, `maintainer_docs_path` list, `find` exclusion)
- modified `.gitattributes` (header sentence, new comment paragraph, `/docs/run-prompts export-ignore`)
- renamed (100%, byte for byte) six files `_dev/run-prompts/2026-11-15-*.md` -> `docs/run-prompts/` (catalog-lists-second-opinion-pages, judge-renders-one-named-page, move-global-rules-into-skill, record-no-review-canvas-decision, relation-map-preset, run-prompts-folder-rule)

## Scope changes
None. Marker lines `for maintainer_docs_path in` and `board_files=(` keep their start. I also reflowed the pre-existing over-long comment line in `.gitattributes` (old line 18) because my edit touched it.

## Red-green evidence
- RED (base, with `docs/run-prompts/probe.md`): contract exit 1, `FAIL: legacy root runtime must be retired at modular cutover: docs (.../docs/run-prompts/probe.md)`.
- GREEN probe at base: `REQ-732 probe FAIL: the contract's docs checks do not name docs/run-prompts`, exit 1.
- GREEN (tip): clean tree exit 0 `staged skills contract: PASS`; with `docs/run-prompts/probe.md` exit 0 PASS; with `docs/other/probe.md` exit 1, `FAIL: legacy root runtime must be retired at modular cutover: docs (.../docs/other/probe.md)`.
- GREEN probe at tip: `REQ-732 probe PASS`. (A first probe run before my `rmdir` failed with `_dev/run-prompts still exists`: git mv had left the empty folder on disk. Removed with `rmdir`; probe then passed.)

## Verification
- `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`: clean 31 s exit 0; run-prompts probe 32 s exit 0; other probe 36 s exit 1 (expected); final rerun after the comment-reflow amend exit 0 PASS. No budget reruns needed. Each run prints `finding GIT-MUTATION-FAILED [error]: context canceled` and `finding REPORT-IMAGE-GENERATION-FAILED [warning]` from a sub-check; the same lines print at base and the contract still passes.
- Probe files and `docs/other` deleted; `git status --short` empty.
- `git check-attr export-ignore -- docs/run-prompts` -> `docs/run-prompts: export-ignore: set`.
- `git archive HEAD | tar -t | grep -c run-prompts` -> 0.
- `git grep -n '_dev/run-prompts' -- . ':(exclude)do-work' ':(exclude)docs/run-prompts'` -> no output.
- `ls _dev/run-prompts` -> No such file or directory.
- `shellcheck _dev/tests/staged-skills-contract.sh` -> exit 1 with 7 SC2016 info notes only, the same 7 at base (none on changed lines); `shellcheck -S warning` exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` -> passed (208 tracked shell files).
- `bash .../REQ-732-probe.sh` from the worktree -> `REQ-732 probe PASS`.

## P-A-U
- PLAN: two list entries + comment in the contract, one export-ignore line + comment in `.gitattributes`, `git mv` of the six prompts, one commit.
- APPLY: done as planned in `342a36a6`; plus `rmdir` of the empty `_dev/run-prompts`.
- UNIFY: `git show --stat HEAD`: 8 files, 15+/9-; `.gitattributes` and contract hunks reviewed line by line, six renames show 0 changed lines; no probe files or debug artifacts left.

## Decisions
- D-04 DECIDE & STATE: reflow the two comments I touched to the neighbours' ~95-column width (including the pre-existing long line in `.gitattributes` header) instead of leaving a long line. Comment-only.
- D-05 DECIDE & STATE: the find exclusion went onto the existing `docs/changelog` line with the `-print -quit` moved to its own line, keeping each line short; behaviour is unchanged.

## Discovered Tasks
- Heavy-lane `staged-skills` coverage in `_dev/tests/heavy-lanes.json` does not include root `docs/`, so the fast gate does not rerun the contract when a file lands under `docs/`. That is why the 2026-11-14 `docs/run-prompts/` file broke the lane unnoticed for five commits (D-03).
- `git mv` leaves an empty source folder on disk; any probe or check that tests "folder is gone" with `-d` sees it until `rmdir`. Minor; no task needed unless repeated.

## Integration seams
- `_dev/tests/staged-skills-contract.sh`: lines 70-73 (comment), 75 (loop list), 107-108 (find exclusion; `-print -quit` now on 108). Everything below shifted down by one line.
- `.gitattributes`: lines 17-19 (header sentence), 36-38 plus blank 39 (new paragraph after the /docs/changelog paragraph), 57 (`/docs/run-prompts export-ignore`, after `/docs/decisions`).

## Proposed CHANGELOG entry
none (no shipped file changes). Confirmed: `.gitattributes` export-ignores `/_dev`, `/docs/run-prompts` and `/.gitattributes`; `suite/modules.tsv` ships only `skills/do-work`, `skills/do-work-board`, `skills/do-work-knowledge`, `skills/do-work-toolbox`. `git archive HEAD` lists no run-prompts path.

## Proposed lesson
none

## Lessons read
- `_dev/primes/lessons-shell-commands.md` `[family: assertion-passes-when-feature-is-dead]` bullets (REQ-591, REQ-592, REQ-716): applied as the negative `docs/other/probe.md` case.
