# Builder Brief — REQ-739 (the second-opinion header keeps a side margin on a phone)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-739-header-phone-side-margin`.
- Worktree (the ONLY tree you write, with the three exceptions below): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-739-header-phone-side-margin`. The coordinator creates it from main HEAD (`49e54593` or later, VERSION 0.324.2) right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-07, Exploration and Scope: read `do-work/working/REQ-739-second-opinion-header-side-margin.md` from the MAIN tree (your worktree may hold an older copy).
- Main-tree paths you may write, by absolute path, never staged and never committed by you:
  1. Hand-back file: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-202422/REQ-739-handback.md`
  2. Captures folder: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-202422/REQ-739-captures/`
  3. Deleting, by name, the files your own browser session created in the main tree's `.playwright-mcp/`.
- You are the only builder in this run. No sibling REQ shares a file with you.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-releases.md` and the required lesson `_dev/primes/lessons-releases.md`.
3. Dropped only for budget, read them with a fixed-string grep, not as whole files: the three `[family: writer-anchored-marker-position]` bullets in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (lines 105-107: no whitespace outside the markers, strip on the masked surface, decorating twice gives the same bytes) and the `[family: session-is-the-backend]` bullets (lines 17-20: measure a header with the page's own CSS at the real width, and look at a phone capture, not only the numbers).
4. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
5. REQ-738, which shipped the header you change: `do-work/archive/UR-180/REQ-738-second-opinion-direction-block-first-screen.md` (Decisions, Review F1-F4, Lessons Learned).
6. The code: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` lines 1565-1589 (`secondOpinionStripBlocks`), 1473 (`secondOpinionRawHTML`, which masks `<style>` bodies) and 1597-1648 (`secondOpinionDecorate`; the block is built on line 1643). The tests that pin the header: `report_second_opinion_test.go` lines 943-1003 and 1005-1077 (`TestSecondOpinionHeaderFoldsTheFullDirection`; its check 7 at line 1074 asserts the exact `aside` open tag), and `report_second_opinion_tree_test.go` lines 43-90.
7. The kit rules the header meets: `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` lines 62, 66-68, 79 and 196-199, and `presentation.js` lines 44-50.

## Decided (do not reopen)
- **D-01 (recommended form):** inside the markers, right before the `aside`, add one `<style>` element: `:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}`. Its zero specificity lets any page rule for the same element win. The `aside` open tag stays byte-identical: `<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">`. No whitespace outside the markers. If your measurement shows this form fails GREEN (2) or (3), pick another form, record it as a decision (D-08 onward) and name the forms you rejected. An inline padding or margin on the `aside` is rejected: it beats the tstack page's own rule.
- **D-02:** the text's left edge is the body margin plus 16 px: 16 px on the kit and on `body{margin:0}` pages, 24 px on a page with the default 8 px body margin. Both count as "about 16 px". Report exact numbers.
- **D-03:** the rule applies at every width, no media query.
- **D-04:** the RED test is one new function named exactly `TestSecondOpinionHeaderKeepsAPhoneSideMargin` in `report_second_opinion_test.go`, beside `TestSecondOpinionHeaderFoldsTheFullDirection`. The GREEN probe looks for that name.
- **D-05:** no contract text changes; `ai-report.md` line 89 stays true.
- **D-06:** the integrator picks the version bump (a patch, 0.324.2 to 0.324.3, is suggested). You only propose the changelog entry.
- **D-07:** the maintainer sees your phone captures before integration starts. Save them as described under "Captures".

## What to build
Two commits.

**Commit 1 (RED, test only).** Add `TestSecondOpinionHeaderKeepsAPhoneSideMargin` to `report_second_opinion_test.go`. Call `secondOpinionDecorate` directly, so the test costs milliseconds (the file runs close to its 30-second budget). The comment above the test names the real failure it pins: the header text started at x = 0 on a 390 px phone on pages with no rule of their own (REQ-738 review F3, the REQ-737 final page), and an inline padding would have overridden the tstack page's own rule. Pin these, and nothing decorative:
1. The bytes between the two markers carry the side-margin rule in a zero-specificity form (a `:where(` selector for the header `aside`).
2. The `aside` open tag carries no inline `padding` or `margin` (assert the exact tag `<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">`).
3. Decorating twice gives the same bytes.
4. The old-block strip case: a page holding a literal 0.324.1 block (copy its exact shape from today's line 1636 and 1643 into the test as a string) decorates to the same bytes as the same page without it.
5. A plain decorate (no slot, no tree marks) leaves no `<style>` from the launcher behind.
Run the test and record the failure text. It must fail on check 1, not on a compile error.

**Commit 2 (GREEN, the launcher).** Change only the block string on line 1643 of `secondOpinionDecorate` and its function comment (lines 1597-1602), so the comment names the side-margin rule and why it is zero-specificity. Do not change the strip calls, the metas, the fold, the insert position, the markers, or any other function. Run the new test (passes) and the whole package (see Verify). Edit an existing assert only if your form breaks it, and say so under "Scope changes".

**Not yours:** `VERSION`, both `CHANGELOG.md` files, anything under `do-work/` (even though the REQ's `write_set` lists the release files: the integrator writes them); `report_second_opinion_tree_test.go`; the two contract docs; every preset; the report kit; the judge; every file under `ai-reports/`.

## Browser proof (scratch copies, never committed)
Published pages are not rewritten. Work in a scratch folder outside every Git tree. Decorate with a throwaway Go test file in your worktree (for example `internal/toolboxcommands/zz_req739_scratch_test.go`) that reads a page, calls `secondOpinionDecorate`, and writes the result to the scratch folder. Rebuild each page's direction from its own head metas: preset slot, id, base, model and date from the `ai-report-second-opinion-*` metas, and the direction text from its `ai-report-second-opinion-direction` meta, HTML-unescaped (a `custom` page passes it as the editor direction). Check that each result's six metas equal the published page's metas byte for byte. Delete the throwaway file before your last commit and check `git status` shows it gone.

Scratch pages, all decorated by YOUR new code:
- **P1 (kit page, no header rule):** a copy of `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` with `presentation.css` and `presentation.js` beside it, decorated with preset `bold` and the shipped `bold.md` text as a tree page (id `v1c01`, base `index.html`).
- **P2 (custom page, no header rule):** copy the bundle folder `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/` to scratch and decorate its `index-v1c01-codex-gpt-6.1-terra-custom-matrix-to-implementation.html` from its metas (`body{margin:0}`, six `details` of its own).
- **P3 (plain page, default body margin):** a small page with no CSS at all, decorated with an editor direction (`custom`).
- **P4 (the page that styles the header itself):** in the same scratch bundle, `index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html` decorated from its metas. Also keep an undecorated copy of the published page as the before-state for check 3 below.
- **P5 (the final REQ-737 page):** in the same scratch bundle, `index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html` re-decorated from its metas.

Serve the scratch folder with `python3 -I -m http.server <free port> --bind 127.0.0.1` (the browser tool blocks `file:` URLs). Measure with ONE script, light theme, at 390 x 844 and at 1440 x 1000. Put the script in scratch, and put its path and its raw output in the hand-back. Per page and per width it reports:
1. The header's left text edge x (the left of the first text box inside the `aside`, for example a Range over its first text node) and the right gap (viewport width minus the right of the `aside`'s content box).
2. The header height (the `aside` box plus its margins), fold closed.
3. The computed `margin` and `padding` of the `aside`. For P4, compare with the published before-state copy: the numbers must be equal (12 px top margin, at least 16 px sides, padding 8 px 14 px), which proves the page's own rule wins.
4. The map frame top y on P4 and P5 (the pages with a map). P4 at 390 must stay at y = 830 or near it (REQ-738 measured 830).
5. Sideways page scroll (`scrollWidth` against `clientWidth`) with the fold closed and with it open.
6. That the fold line ("▸ Preset: …" summary) starts at the same x as the iteration line.
7. P1 only (REQ-738 kit checks): the header's `details` and `summary` are not the kit card (no 2 px border, no 14 px radius, no bold 1.1 rem summary), and a click on "Open all details" opens the kit's `details` and leaves the header's fold closed.

Targets: on P1, P2 and P5 the left text edge is about 16 px at 390 (D-02); on P3 about 24 px; on P4 everything equals the before-state; no sideways scroll anywhere.

Run the render judge on P4 and P5: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge <scratch bundle> --page <page> --out <scratch>/judge-<page id>`. It must report `pass`. If the judge refuses a bundle outside the repository, copy the scratch bundle into an untracked folder of your worktree's `ai-reports/`, judge it there, and delete the folder before your last commit; record that as a decision.

**Captures (D-07).** Save final PNG captures at 390 x 844, fold closed, into the captures folder named above, by absolute path: P1 (kit page), P2 (custom page), P4 (the v1c03 page that styles the header itself) and P5 (the REQ-737 final page). Add P3 if it looks different from P1. List each file in the hand-back section `## Captures for the maintainer` with one line on what it shows, and name the form you chose and any forms you rejected. Delete every other capture before hand-back. Captures the Playwright tools write land in the main tree's gitignored `.playwright-mcp/`: delete the files your session created there by name, nothing else, and record that as a decision. Stop the server.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-739]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not ask the user a question: an open choice becomes a decision in the hand-back.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.
- Nothing under `/Users/t2/.claude/` is ever edited.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l internal/toolboxcommands/report_second_opinion.go internal/toolboxcommands/report_second_opinion_test.go` prints nothing; `go vet ./internal/toolboxcommands/` is clean; `go test -count=1 ./internal/toolboxcommands/` passes (the whole package, as the REQ's Builder Guidance asks, because the block is on every directed page; record its time).
2. `bash _dev/tests/shipped-package-reference-contract.sh`, `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/quiet-grep-pipeline-audit.sh`.
3. `git diff --stat main...HEAD` lists only the launcher file and its test file. `git status --short` shows no scratch test file and no untracked folder under `ai-reports/`.
4. `git diff main...HEAD -- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` touches only the block string and the function comment of `secondOpinionDecorate`.
5. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-202422/REQ-739-probe.sh`, run with your worktree as the current directory, prints `REQ-739 probe PASS (...)`. It fails at base by design: the new test does not exist, so its first check stops it.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb.
- `## Scope changes` — "none", or each change with its reason.
- `## Red-green evidence` — the new test failing after commit 1 (the failure text) and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; the measurement script's path and raw output; the matrix numbers per page at both widths; the judge verdicts; for a rerun under load, both runs.
- `## Captures for the maintainer` — each PNG with one line on what it shows; the form chosen and the forms rejected.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-08 onward (D-01 to D-07 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — see below.
- `## Proposed CHANGELOG entry` — a descriptive title that says what shipped (for example "Second-Opinion Header Keeps a Side Margin on Phones"; check that no earlier entry uses it) and the line "What changes for an installed copy". There: a page that `ai-report second-opinion` builds now gives its short header about 16 px of side room on pages that do not style the header themselves; a page with its own rule for the header keeps its own look; the head metas and the fold are unchanged; pages already published keep their old header until a new pass rebuilds them.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
No sibling REQ runs in this run. The seams are in behaviour:
| Seam | What you report |
| --- | --- |
| A published page may style `body>aside[aria-label="Second-opinion direction"]` (the tstack v1c03 page does). | P4's computed margin and padding before and after, equal, and its map frame top y at both widths. |
| The header sits above the map on relation-map pages, so a taller header moves the map down (REQ-738 kept the map top at y = 830 on the phone). | The header height and map top y on P4 and P5 at both widths, and whether the side padding wrapped any header line on the phone. |
| Published pages are not rewritten. | Confirm `git diff --stat main...HEAD -- ai-reports/` prints nothing. |
| The 30-second per-file budget of `report_second_opinion_test.go`. | The time of `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/` before and after your change. |
| The exact new block. | The bytes between the markers for a tree page with preset `relation-map` (the direction shortened). |

## Notes for the integrator (the builder does nothing about these)
- **Release:** the integrator writes `CHANGELOG.md`, its byte-identical mirror `skills/do-work/CHANGELOG.md`, and `VERSION` (D-06: a patch from 0.324.2 is suggested).
- **Maintainer look:** the coordinator shows the PNG files in `REQ-739-captures/` to the maintainer before integration starts (D-07).
- **Catalog:** no rebuild; no page under `ai-reports/` changes.
- **GREEN probe:** `do-work/runs/work-2026-11-16-202422/REQ-739-probe.sh` (`-- --probe-file`). Pre-flight probe: `do-work/runs/work-2026-11-16-202422/REQ-739-preflight-probe.sh`; it passed at `49e54593` on 2026-11-16 (2.3 s).

Your final message: under 2500 characters. State the branch tip hash, the form you chose, the left text edge x on P1, P3, P5 at 390, P4's margin and padding before and after, the map top y on P4 and P5 at both widths, the verification results, the capture file names, and anything the integrator must decide.
