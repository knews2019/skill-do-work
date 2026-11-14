# Capture prompt: ai-report suggests a Tinka canvas, with pros and cons

Written 2026-11-14. This file is the full text of one capture request. It stands alone for a session that never saw the reports it names.

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: if something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked.

## The request

Make the ordinary `ai-report` end with one more suggestion: review this report's lineage on a Tinka canvas. The suggestion lists the pros and cons of building that canvas for this report, and gives the line that builds it. The user reads the pros and cons and decides per report. There is no pilot counting rule and no new command.

### Context

- Tinka (github.com/bladnman/tinka) is a local board tool. It is driven by `node cli.mjs tools call NAME --config PATH --json JSON`. do-work says **canvas** for a Tinka board, because "board" already means queue-kanban.
- A lineage is one report bundle under `ai-reports/`, every later revision of it, and their second-opinion pages.
- Two proposals tested this on 2026-11-14 at Tinka commit 77ba29c: `ai-reports/2026-11-14_1452_proposal-tinka-review-surface` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface`. Both recommend a pilot by hand. The second one holds the hand recipe under "Check it yourself", and its page `index-v1c01-claude-claude-fable-5-2-custom-worked-example-tally.html` counts the actions.
- Decisions the user made on 2026-11-14: pilot by hand; keep it simple, with the suggestion at the end of the default report in place of a stop rule; Tinka's files live in `.tinka/` at the project root, locally excluded.

### Behavior

1. **Condition.** The suggestion appears only when a `.tinka/` folder exists at the project root. Without that folder a report ends exactly as it does today, and nothing mentions Tinka.
2. **Where.** It is a third block in the closing band "Try another layout or teaching method", after the layout and the teaching method.
3. **What it says.** Two or three pros and two or three cons, written for this report from facts the agent can count without Tinka: the number of pages and reply files in the lineage (each becomes one tile), the number of folded `<details>` sections (a capture shows them closed), and that the canvas opens on this computer only. Then one line an agent can follow to build the canvas.
4. **The line.** A plain instruction that names this bundle and the guide section, for example: put `ai-reports/<this bundle>` and its revisions on a Tinka canvas by following "Review on a Tinka Canvas" in the ai-report guide.
5. **The recipe.** One new section in the guide, "Review on a Tinka Canvas". It tells an agent to: use a source checkout of Tinka pinned to a commit; create a manual layout with `init --config` and `--boards`, both inside `.tinka/`; set update policy off; make one PNG capture per page with `ai-report-judge --out` (a second-opinion page needs a scratch copy of the bundle with that page copied to `index.html`, because the judge renders `index.html` only); copy the captures and each `-reply.md` into one folder inside `.tinka/` and register only that folder as a resource root; add one image tile per capture and one Markdown tile per reply with `tools call item.create`; read feedback with `tools call board.get` and `feedback.list`; hand it to the user as quoted data for `ai-report revise` or `ai-report second-opinion --editor-direction`.
6. **Step 8 summary.** When the block exists, the printed summary names it with its line, beside the two existing suggestions.

### Files to change

- `skills/do-work-toolbox/actions/ai-report.md`: step 9 and the Step 8 summary.
- `skills/do-work-toolbox/actions/ai-report-reference.md`: the rule "End with another layout and a teaching method".
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html`: the `#another-layout` band, with the new block marked as conditional.
- `skills/do-work-toolbox/docs/ai-report-guide.md`: the new section, after "A Second Opinion".
- `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires for shipped files.

### What must not change

- No Go code and no new command.
- No existing `ai-reports/` bundle. Bundles are immutable.
- `git status --porcelain` prints the same before and after the recipe runs. `.tinka/` is excluded through `.git/info/exclude`, never through the tracked `.gitignore`.
- No Tinka portable project (`project init`), because it writes a block into the root `AGENTS.md`.
- No HTML file and no bundle folder as a Tinka resource root. A Tinka HTML item can call Tinka's own write API. The resource root holds only `.png` and `.md` files.
- Creating a report never starts Tinka and never reads Tinka's files. It only checks that `.tinka/` exists. Unattended runs print the suggestion as text and do nothing else.

### How an absent Tinka behaves

- No `.tinka/` folder: no suggestion, no mention.
- `.tinka/` exists but the pinned Tinka checkout is missing: the suggestion is still only text. The recipe's first step checks for the checkout and stops with one line that says what is missing. Every `ai-report` form works unchanged.

### Check first, and park if it fails

Both proposals cite Tinka refusing a resource root whose path has a dot segment (`lib/resources.mjs`, lines 29 to 35 and 126 to 132 at 77ba29c). Tinka also creates its own `.tinka/` folder inside its boards folder. Before writing the recipe, test at the pinned commit whether a resource root inside the project's `.tinka/` is accepted. If Tinka refuses it, park one pending answer for the user: keep the layout in `.tinka/` and move the captures to a locally excluded folder without a leading dot, or choose another place. Do not pick for the user.

### Verification

RED first:

- `grep -c 'Tinka' skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` prints 0.
- `grep -c 'Review on a Tinka Canvas' skills/do-work-toolbox/docs/ai-report-guide.md` prints 0.

GREEN:

- Both greps print 1 or more.
- A report created in a project without `.tinka/` has no Tinka text in its closing band.
- A report created in a project with `.tinka/` has the third block, with pros, cons and the line.
- An agent follows the recipe once on a real lineage. `git status --porcelain` is identical before and after, `git diff --quiet HEAD -- AGENTS.md` passes, and the resource root holds only `.png` and `.md` files.
