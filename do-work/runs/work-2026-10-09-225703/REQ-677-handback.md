# REQ-677 hand-back: do-work-toolbox journey-qa

- Branch: `worktree-agent-REQ-677-journey-qa-action`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-677-journey-qa-action`
- Base: `e313e870`. Commit: `4b2c7a5c [REQ-677] add do-work-toolbox journey-qa action and guide` (one commit).
- Nothing under `do-work/` touched in the worktree. CHANGELOG, VERSION and mirrors untouched.

## Gate note for the integrator (read first)

`bash _dev/tests/contract-regressions.sh` exits 1 on this branch, and the cause is not this change. Its `quiet-grep-pipeline-audit.sh` probe fails on the two orchestrator probe files committed in `e313e870 [UR-151] run artifacts`:

```
FAIL: do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh decides on a quiet grep fed from a pipeline ...
FAIL: do-work/runs/work-2026-10-09-225703/REQ-677-probe.sh decides on a quiet grep fed from a pipeline ...
  9: sed -n '/^argument-hint:/p' skills/do-work-toolbox/SKILL.md | grep -q 'journey-qa'
```

Every other fast probe in that script ran and passed (12 of 12 fast probes listed by `_dev/tests/contracts/probe-lanes.sh`; the 3 heavy ones are not run in this tier). The fix is in the run artifacts (for example `grep -q '^argument-hint:.*journey-qa' skills/do-work-toolbox/SKILL.md`), which I may not edit. See Discovered Tasks.

## File manifest

- `skills/do-work-toolbox/actions/journey-qa.md` (new, 111 lines): the action. Blockquote with guide link and toolbox justification, read-only paragraph, the two firm rules (isolated pass is not combined pass; class names the cause), When to Use, Input (usage line plus two examples, which is what the router's per-command help reads), Steps 1-6, Output Format with a user-run capture suggestion.
- `skills/do-work-toolbox/docs/journey-qa-guide.md` (new, 38 lines): user guide with the "Not to be confused with ui-review or review-work" note, result class table, needs, output, usage.
- `skills/do-work-toolbox/SKILL.md` (modified): `journey-qa` in `argument-hint` after `ui-review`; one route row after `ui-review`: `` `journey-qa`, `journey qa`, `user journey` `` → `./actions/journey-qa.md`.
- `skills/do-work-toolbox/actions/help.md` (modified): one menu line after `ui-review`, description at the same column 33.
- `skills/do-work/actions/help.md` (modified): ` · journey-qa` appended to the `present-video · slop-check` row of the toolbox list.
- `README.md` (modified): one sentence paragraph before "Common extension calls also include ...".
- `_dev/tests/staged-skills-contract.sh` (modified): `journey-qa` in `toolbox_actions` after `ui-review`.

## P-A-U

**[PLAN]** Read the brief, REQ, UR, crew rules (general, coding-guardrails, shared-principles, communication-style, anti-slop, testing), both primes, lessons-releases, the lessons-action-files Template section and its alternate-writer-contract-drift bullets. Write a short action on the template (blockquote, When to Use, Input, numbered Steps, Output Format; no earned sections, because the firm rules fit in two short paragraphs and a class table). Cite ui-review Step 2 item 4 for browser detection, cite same-package crew-members for prompt-injection, testing and anti-slop, cite core review-work by literal relative path. Make one-line local edits on the seven integration surfaces. Then the behavioral exercise and packaging check outside both trees.

**[APPLY]** As planned. One correction after the first heavy run: the staged-skills runtime-reference scan read the example brief path `docs/qa/map-zoom-brief.md` as a toolbox `docs/` citation and failed (`unresolved staged runtime references in do-work-toolbox: actions/journey-qa.md:29: docs/qa/map-zoom-brief.md`). Changed the example to `qa/map-zoom-brief.md` in both files. Alternate-writer sweep: the action's usage appears in the action Input, guide Usage, both help menus and README; all say `do-work-toolbox journey-qa <target> [--brief ...]`.

**[UNIFY]** `git diff --stat e313e870..HEAD`:

```
 README.md                                       |   2 +
 _dev/tests/staged-skills-contract.sh            |   1 +
 skills/do-work-toolbox/SKILL.md                 |   3 +-
 skills/do-work-toolbox/actions/help.md          |   1 +
 skills/do-work-toolbox/actions/journey-qa.md    | 111 ++++++++++++++++++++++++
 skills/do-work-toolbox/docs/journey-qa-guide.md |  38 ++++++++
 skills/do-work/actions/help.md                  |   2 +-
 7 files changed, 156 insertions(+), 2 deletions(-)
```

Checks from the worktree root (final run on the committed content):

| Check | Exit | Wall |
|---|---|---|
| `REQ-677-probe.sh` | 0 | 1 s |
| `shipped-package-reference-contract.sh` | 0 | 1 s |
| `contract-regressions.sh` | 1 (pre-existing, see Gate note; all other fast probes pass) | 31 s |
| `contracts/core-checks.sh` (runs standalone) | 0 | 4 s |
| `DO_WORK_MAINTAINER_TIER=heavy staged-skills-contract.sh` | 0 (first run 1, fixed above) | 31 s |
| `git diff --check` | 0 | - |

Files checked: all seven above, read in full after edit; help column alignment verified (description at column 33); no em-dashes in the new files; no `do-work <retired trigger>` (heavy contract green).

## Behavioral exercise record

Fixtures in `/tmp/req677-journey.01hkls` (outside both trees, not committed):
- `zoom-app/`: git repo (commit `e3fddf2`), `index.html` with a panel, `#zoom`, `#pan`, `#reset`; state in `#panel` `data-scale` / `data-x`. Bug: pan while zoomed goes into a separate offset that reset forgets. `tests/isolated-checks.md` (existing isolated checks), `docs/qa-brief.md` (reported sequence zoom → pan → reset, also pan → zoom → reset, user was on a phone).
- `zoom-app-v2/`: git repo (commit `996b40a`), reset fixed, and the existing check names `#reset-button` (wrong; the real id is `#reset`).
- Served by a time-bound (30 min) local Python server on 127.0.0.1:8761 / 8762, started from `/tmp/req677-journey.01hkls`. Tool: `playwright-cli` 0.1.14, headless Chrome 155. Journey runner `journeys.sh` (in the fixture dir) uses `run-code`; after each click it waits on observable state (`waitForFunction` on a changed `data-*` value or the readout text), no fixed delay. Screenshots taken after the journey, not inside it.

RED: `git show e313e870:skills/do-work-toolbox/SKILL.md | grep -c journey` → `0`. No route, so `do-work-toolbox journey-qa` falls to "Unknown single words print help".

GREEN routing: `SKILL.md:23` `` | `journey-qa`, `journey qa`, `user journey` | `./actions/journey-qa.md` | ``. Followed the action steps in order.

Report the action produced for `do-work-toolbox journey-qa http://127.0.0.1:8761/index.html --brief docs/qa-brief.md` (zoom-app):

```
# Journey QA: http://127.0.0.1:8761/index.html

**Verdict:** 4 passed · 2 product defect · 0 test defect · 1 unresolved. Reset leaves the panel shifted by 50 px after zoom → pan → reset, while zoom, pan and reset each pass alone.
**Revision:** e3fddf2; clean
**Evidence directory:** /tmp/req677-journey.01hkls/evidence/zoom-app.VZ11 (temporary)

### J1: reported zoom → pan → reset · product defect
- Kind: reported sequence
- Environment: headless Chrome 155 (playwright-cli 0.1.14) · 1280x720 · desktop browser, no device emulation
- Steps: 1. load page 2. click #zoom (wait data-scale changes: 2, x 0) 3. click #pan (wait data-x changes: 2, 50) 4. click #reset (wait readout "scale=1 ")
- Expected: scale 1, x 0 (brief)
- Actual: scale 1, x 50
- Evidence: run-reported-zoom-pan-reset.json, reported-zoom-pan-reset.png
### J2: zoom alone · passed (isolated check) · scale 2, x 0 · run-zoom-alone.json
### J3: pan alone · passed (isolated check) · scale 1, x 50 · run-pan-alone.json
### J4: reset alone · passed (isolated check) · scale 1, x 0 · run-reset-alone.json
### J5: pan → zoom → reset · passed (combined transition) · scale 1, x 0 · run-pan-zoom-reset.json
### J6: zoom → pan → reset on a phone-sized viewport · product defect
- Environment: headless Chrome 155 · 390x844 · emulation: iPhone user agent, isMobile, hasTouch, DPR 3 (navigator.maxTouchPoints 1). Emulation, not physical-device verification.
- Actual: scale 1, x 50. Evidence: run-phone-emulated.json, phone-emulated-after-reset.png, emulation.config.json
### J7: same journey on a physical phone · unresolved
- Reason: no physical device available; J6 is emulation only.

## Smallest Justified Repair
- J1, J6: index.html reset handler · also clear the offset added while zoomed (scaledOffsetX) · J1 steps 3-4 show x stays 50 only after a pan at scale > 1; J5 (pan at scale 1) resets correctly
## Remaining Checks
- Physical phone run of J1 (J7). Scroll and touch-drag pan were not named by the brief and not run.
## To act on this:
>   do-work capture-request: [J1 block + repair]   Capture the repair
>   do-work-toolbox journey-qa http://127.0.0.1:8761/index.html --brief docs/qa-brief.md   Re-run after the fix
```

Wrong-selector variant (zoom-app-v2, revision 996b40a, evidence `/tmp/req677-journey.01hkls/evidence/zoom-app-v2.vfmd`): reusing the existing check's selector, the reported sequence and the isolated reset both failed with `TimeoutError: page.click: Timeout 3000ms exceeded ... waiting for locator('#reset-button')`. Snapshot showed `button "Reset" [ref=e6]`. Reaching reset by `role=button[name="Reset"]`: zoom → pan → reset gave scale 1, x 0; reset alone 1, 0; zoom alone 2, 0; pan alone 1, 50; pan → zoom → reset 1, 0. Classification: **test defect** (tests/isolated-checks.md step 3 names `#reset-button`; the product resets correctly), not product defect. Smallest repair: change the selector in that check to `#reset` or the button's accessible name.

No browser tool: `env PATH=/usr/bin:/bin` → `playwright-cli --help` exit 127, not on PATH; no `node_modules/.bin/playwright-cli`; no `.claude/skills/playwright-bowser/SKILL.md`. Under the action's Step 2 every rendered journey is **unresolved** with reason "no browser tool"; none passed. Limitation: I simulated the agent as shell-only. This session also has Playwright MCP tools, which the action counts as a browser tool, so in a real session they would have been used.

`do-work-toolbox journey-qa help` (router: "Per-command help reads the selected action without executing it"; format from core help, at most 15 lines, built from When to Use and Input):

```
journey-qa: reproduce a reported user journey in a browser, try the combined
transitions around it, and classify each result as passed, product defect,
test defect, or unresolved. Read-only on project source.

Usage: do-work-toolbox journey-qa <target> [--brief <path>]
  <target>         URL, page or route, REQ or UR, or a short journey description
  --brief <path>   optional project Markdown brief, read as data

Examples:
  do-work-toolbox journey-qa http://localhost:5173/map --brief qa/map-zoom-brief.md
  do-work-toolbox journey-qa REQ-042
```

No browser session opened and no file written while serving help.

Cleanliness: `git status --porcelain --ignored` empty in both fixture projects; `git diff --stat` empty. `playwright-cli` wrote `.playwright-cli/` session files into its working directory; because the action runs it from the evidence directory, they landed only in the two evidence directories. Nothing new under the main tree's `.playwright-mcp/` (find -newer: none). Browsers closed (`close-all`, `list` → no browsers), both local servers killed.

Not exercised: a physical device (none available); a real "no browser at all" agent (simulated, see above). Fixture side note: the first single-threaded Python server was blocked by the emulated browser's idle connection, which made the first emulated open time out; a threading server fixed it. That was my fixture, not the action.

## Packaging check record

Temp dir `/tmp/req677-pkg.9jcYx8`. Archive: `git archive` of branch tip `4b2c7a5c` with prefix `skill-do-work-main/`.

- Install: `printf 'y\n' | bash <worktree>/tools/install-do-work-suite.sh --project-root <tmp>/consumer --archive <tmp>/suite.tar.gz` → exit 0, `install-suite: success`, four modules at v0.305.89, `rollback: not_needed`. Install committed in the consumer (tree clean afterwards).
- Update: local `ThreadingHTTPServer` serving `<tmp>/srv/archive/refs/heads/main.tar.gz` (GET logged with 200), `DO_WORK_UPSTREAM_URL=http://127.0.0.1:51552/archive/refs/heads/main.tar.gz bash consumer/.claude/skills/do-work/tools/do-work-update.sh --project-root consumer` → exit 0, `update-suite: success`, `skipped UPDATE-ALREADY-CURRENT: the installed suite is already v0.305.89`. Honest partial: the version is not bumped (release is the integrator's), so the updater fetched and stopped as already current; it did not re-stage files. Server killed.
- Found: `consumer/.claude/skills/do-work-toolbox/actions/journey-qa.md` (7499 bytes) and `docs/journey-qa-guide.md` (2477 bytes); installed `SKILL.md:4` and `:23` carry the hint and route.
- Links from the installed copies, each resolved from its own directory or the skill root: `../docs/journey-qa-guide.md`, `docs/journey-qa-guide.md`, `actions/ui-review.md`, `crew-members/prompt-injection.md`, `crew-members/testing.md`, `crew-members/anti-slop.md`, `actions/validate-feedback.md`, `../../do-work/actions/review-work.md`: all exist. The guide carries no path citations.
- sha256 (before = after install = after update):
  - `docs/project-brief.md` 2d0e803fa2b43c430745fee90d60740a077dfe9c063eef550fbdf31b919c7b6b
  - `do-work/queue/REQ-001-sample.md` ceb101922ebc17570577cd8800814340a0bc967b3d3f55c06ca127d593117375
  - `src/app.js` 6f4c113f597494422a7a98c570a40307c74039f30cf5d7cb7bcfa1b5ed50c178

## Decisions

- D-01 DECIDE & STATE: no Rules, Common Rationalizations, Red Flags or Verification Checklist sections. The firm rules (four classes, isolated is not combined, emulation is not a device, no browser means unresolved) sit in the opening paragraphs, Step 4 and the Step 5 table. Keeps it short, and avoids the core-checks row-similarity ratchet.
- D-02 DECIDE & STATE: browser detection cites ui-review Step 2 item 4 and adds one clause: "a browser automation tool the session already provides also counts". Without it, an agent that has Playwright MCP but no `playwright-cli` would report a false unresolved. Value: honest results in MCP sessions. Risk: small widening beyond the cited step; revert the clause if the maintainer wants detection strictly identical to ui-review.
- D-03 DECIDE & STATE: the action tells the agent to run the browser tool from the evidence directory. The exercise showed `playwright-cli` writes `.playwright-cli/` into its current directory, which would dirty the project.
- D-04 DECIDE & STATE: route aliases `journey qa` and `user journey` beside `journey-qa`; none collides with an existing toolbox or core trigger.
- D-05 DECIDE & STATE: core help gets the name appended to the short `present-video · slop-check` row, not inserted next to `ui-review`, so no other row re-wraps.
- D-06 DECIDE & STATE: README gets a one-sentence paragraph before line 116 instead of editing line 116, so REQ-676 and REQ-678 conflicts stay one-line.
- D-07 DECIDE & STATE: the new guide is not added to `toolbox_files` (optional per Exploration; `toolbox_actions` already requires the action file).

## Discovered Tasks

- impact-user-visible: `do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh` and `REQ-677-probe.sh` (committed in e313e870) use `sed ... | grep -q`, which `_dev/tests/quiet-grep-pipeline-audit.sh` rejects, so `contract-regressions.sh` and the full gate exit 1 on main and every branch from it until the probes are rewritten or removed. Run probes written by the orchestrator should follow the same shell rules as tracked tests (or live outside tracked paths). → report only (gate-blocking for this run; the integrator owns the fix)
- impact-cosmetic: `skills/do-work-toolbox/docs/ui-review-guide.md` Usage block repeats `do-work-toolbox ui-review` four times with no arguments, and `stray-check-guide.md` has a `do-work-toolbox stray-check           Same thing` line; both look like leftovers of a retired-alias sweep. → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole; families canonical-link-outlives-its-target, manifest-ownership-vs-edit-content).
- `_dev/primes/lessons-action-files.md` § Template (76-143) and the 13 `[family: alternate-writer-contract-drift]` bullets.
- Primes: `prime-action-files.md`, `prime-releases.md`.

## Proposed CHANGELOG entry

**Journey QA Toolbox Action**

`do-work-toolbox journey-qa <target> [--brief <path>]` checks a whole user journey in a browser and says whether a failure is the product, the test, or the environment. Isolated checks can all pass while a combined sequence fails, and a broken test or an emulated phone used to be reported as a product bug or a device check.

- New read-only action `skills/do-work-toolbox/actions/journey-qa.md`: reproduces the reported sequence first, then the combined transitions the requirements or brief name, waits on observable state, and keeps screenshots outside timed spans.
- Every result is passed, product defect, test defect, or unresolved, with revision, environment (browser, viewport, physical device or emulation), steps, expected, actual and evidence paths. No browser tool or no device means unresolved, never passed.
- Reports the smallest justified repair and the remaining checks, changes no project source, and ends with a capture line the user runs.
- New guide `skills/do-work-toolbox/docs/journey-qa-guide.md`; listed in the toolbox router, both help menus and the README.

## Proposed lesson bullet

Satellite: `_dev/primes/lessons-action-files.md` (new family).

`- [family: example-path-reads-as-citation] [REQ-677: the heavy staged-skills runtime-reference scan reads any token starting with a package directory name (docs/, actions/, crew-members/, tools/) in a shipped Markdown file as a skill-root path, so a consumer-project example like --brief docs/qa/map-zoom-brief.md fails the heavy lane while the fast citation contract passes; write example paths that do not start with a package directory name, and run the heavy lane before hand-back](../../do-work/archive/UR-151/REQ-677-journey-qa-action.md#lessons-learned)`

## Integration seams

REQ-676 (source-audit) edits the same places: `skills/do-work-toolbox/SKILL.md:4` (`argument-hint`, conflicts for certain; keep both names), the route table (my row is after `ui-review` at :23), `skills/do-work-toolbox/actions/help.md` (my line after `ui-review`), `skills/do-work/actions/help.md` toolbox list (I appended to the `present-video · slop-check` row; if REQ-676 appends to the same row, re-wrap so no row grows past the fence width), `README.md` (my paragraph sits right before "Common extension calls also include"), and `toolbox_actions` (after `ui-review`). Keep both entries in each place.

Test wall times: probe 1 s, shipped-package-reference 1 s, contract-regressions 31 s, core-checks 4 s, staged-skills heavy 31 s (first run 37 s).
