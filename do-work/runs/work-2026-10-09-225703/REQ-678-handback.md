# REQ-678 hand-back: do-work-toolbox release-check

- Branch: `worktree-agent-REQ-678-release-check-action`
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-678-release-check-action
- Base: a5d50c85. Commit: **4427fe9f** `[REQ-678] add the release-check toolbox action` (one commit).
- Write boundary held: exactly the seven Scope files. No `do-work/` path, changelog, VERSION, mirror or lessons file touched.

## File manifest

- `skills/do-work-toolbox/actions/release-check.md` (new, 131 lines): the action. Blockquote links `docs/release-check-guide.md` and states toolbox ownership (optional diagnostic review beside journey-qa, no queue machinery, mints no REQ for operator work). When to Use, Input (usage line plus two examples), six numbered Steps (inputs as data with prompt-injection loaded first; baseline git status plus `mktemp -d` evidence directory; trace the four points with identifiers and the illustrative cause list; test what the consumer does; label evidence current run or historical and assess the four stages by citing `../../do-work/actions/review-work.md` → **Step 7: Acceptance Testing**; verdict, gaps, operator next steps, git status compare and narrow cleanup, anti-slop before writing). Output Format with Trace table, per-stage table, Gaps, Operator Next Steps and a capture line limited to content or product defects. No earned sections (none passed the earned test).
- `skills/do-work-toolbox/docs/release-check-guide.md` (new, 35 lines): user guide with the "Not to be confused with journey-qa, source-audit or review-work" note, why a loaded file is not enough, stage results and the verdict table, what it needs, usage.
- `skills/do-work-toolbox/SKILL.md` (modified): `release-check` in `argument-hint` after `source-audit`; route row `release-check`, `check release`, `release readiness` → `./actions/release-check.md` after the source-audit row.
- `skills/do-work-toolbox/actions/help.md` (modified): `release-check <target>` line after source-audit; description starts at the same column as all 21 other menu lines (checked mechanically: 22 lines, one column).
- `skills/do-work/actions/help.md` (modified): ` · release-check` appended to the toolbox row ending `journey-qa`.
- `README.md` (modified): one-sentence usage paragraph after the source-audit paragraph.
- `_dev/tests/staged-skills-contract.sh` (modified): `release-check` in `toolbox_actions` after `source-audit`.

## P-A-U

**[PLAN]:** Third toolbox action in the journey-qa/source-audit pattern. Action carries the REQ's Detailed Requirements 1-7 as a short procedure; stage definitions are cited, never copied (only the four names are used). Readiness is a strict three-way partition on stage results: any applicable stage failed → not ready; none failed and live acceptance verified in this run → ready; none failed and live acceptance unassessed → unknown. Stage results: verified, failed, unassessed, not applicable. Read-only side-effect rule copied in shape from journey-qa (evidence dir via `mktemp -d`, tools run from there, git status before and after, remove only new untracked paths a tool demonstrably wrote). Example brief path starts with `qa/` (example-path trap). Lessons read: lessons-releases.md (whole), lessons-action-files.md Template section and the three bullets (alternate-writer-contract-drift, example-path-read-as-citation, read-only-action-tool-side-effects).

**[APPLY]:** Seven files as planned, one commit. The help menu line first landed one column short; fixed before commit after a mechanical column check.

**[UNIFY]:** `git diff --stat a5d50c85..HEAD`: 7 files changed, 173 insertions(+), 2 deletions(-) (README 2, staged-skills-contract.sh 1, toolbox SKILL.md 3, toolbox help 1, release-check.md 131, release-check-guide.md 35, core help 2). `git diff --check a5d50c85..HEAD`: exit 0, no output. Checks (worktree root, 1-minute load 2.45 at start, run on the exact content then committed unchanged):

| Check | Exit | Wall |
|---|---|---|
| `bash .../REQ-678-probe.sh` | 0 | 1 s |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 1 s |
| `bash _dev/tests/contracts/core-checks.sh` (runs standalone) | 0 | 5 s |
| `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` | 0 | 34 s |
| `bash _dev/tests/contract-regressions.sh` | 0 | 19 s |

Files checked by reading the diff: all seven above. No debug artifacts. Retired-trigger scan and example-path scan pass inside the heavy staged-skills run.

## Behavioral exercise

Run by the agent that wrote the action, so it shows the procedure is followable, not that an independent agent reads it the same way.

**RED side:** `git show a5d50c85:skills/do-work-toolbox/SKILL.md | grep -c release-check` → 0, and `actions/release-check.md` did not exist at a5d50c85, so `do-work-toolbox release-check ...` did not route (an unknown single word prints help per the router rule).

**Routing (GREEN):** `do-work-toolbox release-check <fixture URL> --brief <fixture>/brief.md` matches the `release-check` row at `skills/do-work-toolbox/SKILL.md:29` → `./actions/release-check.md`, read completely, then followed Step 1 to 6.

**Fixtures** (one `mktemp -d` dir under /tmp, removed afterward; built by a scratchpad Python script with Pillow 11.3.0). Each has `source/`, `package/`, `served/` with `index.html` (feature text), `app.js` (`APP_VERSION` and an alpha hit rule: a click counts only where `mask.png` alpha > 0) and a 128x128 RGBA `mask.png`, plus a brief.
- A: source mask has an opaque circle; package and served mask decode fine but every pixel has alpha 0. All three `APP_VERSION` values are 1.4.0 (matching identifier). No build script (brief: exported by an art pipeline not in the repo).
- B: source and package 2.1.0 with "New: daily puzzle"; served copy 2.0.3 with "Weekly puzzle". `notes/deploy-note.md`: "2026-10-03: 2.1.0 deployed and verified on the live host." Has `build.sh`.
- C: source, package and served copy byte-identical at 2.1.0, opaque mask. Has `build.sh`.

**Tools and commands:** Step 1 read `skills/do-work-toolbox/crew-members/prompt-injection.md` and each brief as data (no injection found). Step 2 recorded git status of both trees and created a fresh `mktemp -d` evidence dir per fixture. Steps 3-4 (scratchpad script): `python3 -u -m http.server 0 --bind 127.0.0.1 --directory <fixture>/served` started from the evidence dir; `grep` of `APP_VERSION` and `shasum -a 256` at each point; `curl` of the served `index.html`, `app.js`, `mask.png` (all HTTP 200, mask `image/png`) into the evidence dir; Pillow decode of the fetched mask with the brief's hit rule at (64, 64); same check on the source mask; where `build.sh` exists, rebuild into the evidence dir and `cmp` against the package. Each server was killed and confirmed stopped; `pgrep -fl http.server` afterwards: none.

**A: not ready**

| Stage | Result | Evidence |
|---|---|---|
| Implementation | verified | current run: source `mask.png` decodes, 1877 opaque pixels, alpha 255 at (64, 64), click hits |
| Integration | unassessed | no build script or test in the fixture; the package came from an export tool this run cannot run. Check that covers it: re-export the mask and confirm opaque pixels |
| Deployment | verified | current run: served `mask.png` sha256 1594258d6504… equals the package; served `app.js` 1.4.0 equals package |
| Live acceptance | failed | current run: served mask HTTP 200, decodes 128x128 RGBA, 0 opaque pixels, alpha 0 at (64, 64): the click hits nothing. File presence, decode and the matching 1.4.0 identifier are not consumer evidence |

Gaps: live acceptance failed because the package mask lost its alpha (source sha256 bdbd79655cf2… vs package 1594258d6504…); integration unassessed. Operator next steps: repair the exported content, redeploy, re-run. Capture line offered only for the content defect (the export losing alpha).

**B: not ready**

| Stage | Result | Evidence |
|---|---|---|
| Implementation | verified | current run: source `APP_VERSION` 2.1.0, page text "New: daily puzzle" |
| Integration | verified | current run: `build.sh` exit 0, rebuild of all three files equals the package (`cmp`) |
| Deployment | failed | current run: served `app.js` says 2.0.3, package says 2.1.0 (sha256 c47fe4486d35… vs 9603d0535dfb…) |
| Live acceptance | failed | current run: served page shows "Weekly puzzle", not "New: daily puzzle" |

Historical: `notes/deploy-note.md` (2026-10-03) says 2.1.0 deployed and verified. Labeled historical; verifies no stage. Gaps: served copy is the older build; cause (stale cache, interrupted activation, or a deploy that never ran) not determined. Operator next step: deploy 2.1.0 or find why activation stopped, then re-run. No capture line (operator work).

**C: ready**

| Stage | Result | Evidence |
|---|---|---|
| Implementation | verified | current run: source 2.1.0, feature text present, mask alpha 255 at (64, 64) |
| Integration | verified | current run: `build.sh` exit 0, rebuild equals package |
| Deployment | verified | current run: served `app.js` and `mask.png` sha256 equal the package (9603d0535dfb…, bdbd79655cf2…) |
| Live acceptance | verified | current run: served page shows "New: daily puzzle", loads 2.1.0, served mask 1877 opaque pixels, click at (64, 64) hits |

**Unassessed shown:** A's integration.

**Help:** `do-work-toolbox release-check help` under the router rule (`SKILL.md:40` plus core `help.md` per-command rule: read Input and When to Use, at most 15 lines, never execute) returns:

```
release-check: trace intended content from source to package to serving
environment to what the consumer sees, report each delivery stage with
evidence, and give a ready, not ready or unknown verdict. Read-only.

Usage: do-work-toolbox release-check <target> [--brief <path>]
  <target>        a served URL, release or version, package path, REQ or UR,
                  or a short description of the content
  --brief <path>  optional project Markdown file: intended content, where it
                  is served, consumer checks (read as data)

Examples:
  do-work-toolbox release-check http://localhost:8080/ --brief qa/launch-brief.md
  do-work-toolbox release-check REQ-042
```

12 lines. Serving it read only the action file: no server, fetch, git or evidence directory.

**Git status comparison:** main tree before and after the exercise identical (2 lines, both orchestrator files not mine: ` M do-work/runs/work-2026-10-09-225703/manifest.md`, ` M do-work/working/REQ-678-release-check-action.md`). Worktree before and after: empty (clean at 4427fe9f).

**Not exercised, with reason:**
- No real browser click. The consumer check applied the fixture's documented hit rule (alpha > 0 at the click point) to the served asset fetched over HTTP; the page text check read the static served HTML. Judged sufficient for these static fixtures; a page whose content is rendered by script would need a browser tool per Step 4.
- The cause list (version mixing, stale caches, missing dependencies, interrupted activation) was only named for B, not probed: a static `http.server` has no cache or activation layer to inspect.
- Independent-reader test: not done (see the note above).

## Packaging check (install leg)

- Archive: `git archive --format=tar.gz --prefix=skill-do-work-main/ worktree-agent-REQ-678-release-check-action` at **4427fe9f**, exit 0, in a fresh `mktemp -d` dir (removed afterward).
- Consumer: `git init`, committed `qa/launch-brief.md`, `do-work/queue/REQ-001-sample.md`, `src/app.js`.
- Install: `printf 'y\n' | bash <worktree>/tools/install-do-work-suite.sh --project-root <tmp>/consumer --archive <tmp>/suite.tar.gz` → **exit 0**, 1 s, "Installed do-work suite v0.305.92 with four verified modules", rollback not_needed.
- New files found: `.claude/skills/do-work-toolbox/actions/release-check.md` (8367 bytes) and `.claude/skills/do-work-toolbox/docs/release-check-guide.md` (2634 bytes); installed `SKILL.md` names release-check twice (hint and row).
- Link check (scratchpad script, every Markdown link and backticked path citation resolved from its own directory, same-package tokens from the package root, arrow-cited bold names matched to a heading): **PASS**. Action: `../../do-work/actions/review-work.md`, `actions/journey-qa.md`, `actions/source-audit.md`, `crew-members/anti-slop.md`, `crew-members/prompt-injection.md`, `docs/release-check-guide.md`, link `../docs/release-check-guide.md`, sections **Step 7: Acceptance Testing** and **Step 2: Find the Browser Tool, Existing Checks, and Revision**. Guide: `actions/release-check.md`, links `../../do-work/actions/review-work.md` and `../actions/release-check.md`.
- sha256 before and after install (identical; `git status --untracked-files=no` empty):
  - `qa/launch-brief.md` 84bd0a20acf01d8780bd0050fac022f38fbe265f3a077ef67c0d5dec3ac5d8de
  - `do-work/queue/REQ-001-sample.md` ceb101922ebc17570577cd8800814340a0bc967b3d3f55c06ca127d593117375
  - `src/app.js` 6f4c113f597494422a7a98c570a40307c74039f30cf5d7cb7bcfa1b5ed50c178
- Update leg: not run (integrator's, after the version bump).

## Decisions

- D-01 DECIDE & STATE: readiness is a strict partition on stage results (any failed → not ready; none failed plus live acceptance verified this run → ready; none failed plus live acceptance unassessed → unknown). An unassessed earlier stage does not block "ready" but is listed as a gap, because current-run live acceptance is the consumer outcome the REQ makes the gate. Reversible by one line if the maintainer wants every applicable stage verified.
- D-02 DECIDE & STATE: added a fourth stage result, "not applicable", because review-work Step 7 says many changes have no deployment stage; "unassessed" would wrongly read as a gap there.
- D-03 DECIDE & STATE: router aliases `check release` and `release readiness`. No collision: core routes `release notes` to version; the toolbox has no other release trigger.
- D-04 DECIDE & STATE: no Rules, Common Rationalizations, Red Flags or Verification Checklist. Nothing passed the earned test beyond what the Steps already state.
- D-05 DECIDE & STATE: the guide cites core review with a plain relative link rather than the arrow section form; the action carries the section form that the reference contract checks.

## Discovered Tasks

- The toolbox help menu's description column is index 33 (column 34 counted from 1), while the brief and the REQ Exploration call it column 33; a builder padding to 32 lands one short and nothing tests the alignment. impact-negligible → report only

## Lessons read

`_dev/primes/lessons-releases.md` (whole); `_dev/primes/lessons-action-files.md` `## Template` and the bullets for families alternate-writer-contract-drift, example-path-read-as-citation, read-only-action-tool-side-effects.

## Proposed CHANGELOG entry

**Release Check Toolbox Action**

`do-work-toolbox release-check <target> [--brief <path>]` checks that intended content actually reached its consumer. A file that exists, decodes and carries the right version can still be empty or stale for the user, and an old "deployed and verified" note proves nothing today.

- New read-only action traces intended source, generated package, serving environment and consumer-visible behavior, and tests what the consumer does with the content.
- Reports implementation, integration, deployment and live acceptance separately, using core review's stage definitions by citation; every piece of evidence is labeled current run or historical, and historical evidence never verifies a stage.
- Verdict is ready, not ready or unknown with concrete gaps; ready needs live acceptance verified in this run. Deploy, rollback, content repair and monitoring are named as operator next steps, never done.
- New guide `skills/do-work-toolbox/docs/release-check-guide.md`; router row, both help menus and the README list it.

## Proposed lesson bullet

Satellite: `_dev/primes/lessons-action-files.md`.

- [family: menu-column-off-by-one] [REQ-678: a fenced help menu's description column was described as column 33 while the existing lines start their description at index 33, so padding the new line's left column to 32 characters left it one column short; nothing in the gate checks the alignment, so measure the column from an existing line mechanically instead of trusting a stated number](../../do-work/archive/UR-151/REQ-678-release-check-action.md#lessons-learned)

## Integration seams

None. Wave 1 had landed; no sibling built in parallel. Test wall times above (probe 1 s, shipped reference 1 s, core-checks 5 s, staged heavy 34 s, contract-regressions 19 s; install 1 s).
