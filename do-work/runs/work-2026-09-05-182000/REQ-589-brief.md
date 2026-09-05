# Builder brief — REQ-589

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1820/worktree-agent-REQ-589-m4-slim-band`
- **Your branch (already checked out there):** `worktree-agent-REQ-589-m4-slim-band`
- **Route:** A
- **Base commit:** the branch's HEAD (main at dispatch; it contains REQ-588's merge ab251f24 and release 0.303.2)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work-board/tools/queue-kanban/VERSION` — release paths owned by finalization.
- Any file outside the write set below. If you need one, stop and report it in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). Run only the focused tests named below.
- Do not build or serve the board on port 8090: a live board owned by the user is running there.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/frontend.md`

Also read `_dev/primes/prime-kanban-board.md`, the shipped `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, and in the two lesson satellites (`skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, `_dev/primes/lessons-kanban-board.md`) the bullets touching the web assets and the Node behaviour lane, including the newest bullet `[family: subject-not-restated-in-detail]`.

## The REQ

Read it in full (read only): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-589-render-the-verify-findings-strip-as-the-m4-slim-band.md`. Its `## Prior Implementation` says what REQ-588 built; its `## Requirements` D1–D6 are the contract.

## The specification: mock-up M4

The approved design is the M4 pages in `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/ai-reports/2026-09-05_1800_REQ-588-verify-findings-slim-band-gallery/mockups/`: `m4-closed.html`, `m4-open.html`, `m4-open-remedy.html`, sharing `shared.css` (a copy of the board's tokens and the pre-REQ-588 strip rules, plus the top bar). Open the gallery's `index.html` for the reasoning and the captures.

- The page's own `<style>` block is the CSS you ship: the `.board-findings` band, `.vf-head`, `.vf-label`, `.vf-count`, `.vf-dot*`, `.vf-subject`, `.vf-cat`, `.vf-fix`, `.vf-detail`, the `.vf-rows`/`.vf-row`/`.vf-chev`/`.vf-remedy` rules, and the `.vf-strip`/`.vf-subjects`/`.vf-toggle` rules. Port them into `web/board.css` in place of the REQ-579/REQ-588 `.board-findings-*` row rules (keep the `.board-findings` selector name for the section and the two host ids; you may keep or rename the row class names, but delete the rules the new ones replace — no dead CSS).
- The page's markup is the structure `renderVerifyFindingsStrip` must produce: an outer `<details>` (the strip) whose `<summary>` holds the header line (warning glyph SVG, VERIFY label, counts, the subject list with weight dots, the Show/Hide toggle); inside it the rows, each a `<details>` whose `<summary>` holds dot, subject, category words, detail, optional "cleanup can fix" pill and chevron SVG, and whose content is the "What to do:" remedy block. Skipped probes are rows with the grey dot and category "not checked", no remedy.
- `web/template.html`: the strip section keeps `id="board-findings"`, its `hidden` attribute, `#board-findings-count`, and the two hosts `#board-findings-cards` and `#board-findings-skipped-list` (REQ-578's `applyView` in `web/board-controls.js` reads their children to decide whether the strip has content; do not change that reader). The hint sentence goes. Decide whether the outer `<details>` is static in the template with the hosts inside it (preferred: the template owns the shell, the renderer fills subjects and rows) or built by the renderer, and log it as a decision. Keep the two hosts as `display: contents` pass-throughs inside the rows container.
- Grouping by subject stays exactly as `groupFindingsBySubject` does it (exact string match on the payload field); rows render in that order. A subject heading element is no longer needed because every row carries its subject; the summary's subject list shows each finding's subject once per finding (not once per group) with that finding's dot.
- Category words: `category.replace(/-/g, " ")` — no list, no mapping.
- Remembered state (D4): one `localStorage` key (name it like `detailPanelWidthStorageKey` in `web/board-detail.js`; look at how that file guards reads and writes and copy the guard shape), written on the outer details' `toggle` event, read at render to set the initial `open`. Default closed. Row state not persisted.
- Both SVGs (warning glyph, chevron) are inline markup created with the DOM, never `innerHTML` with payload text; every payload string goes through `textContent`/`createElement` exactly as the strip does today.

## Write set (your write boundary)

- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go`

If the existing Node-lane cases from REQ-579/REQ-588 (`TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList`, `TestJavaScriptBehaviorVerifyFindingRemedyIsItsOwnLineAfterTheDetail`) pin structure this REQ replaces, rewrite them to pin the new structure rather than deleting them; `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` (REQ-578) and the hide-when-empty case must pass unchanged. If any other file turns out to be required by D1–D6 (for example a stub in the Node harness that lacks `localStorage` or `toggle` events), that is the REQ's own requirement demanding the file class: flag it in the hand-back and proceed.

Board version and lock-step follow `_dev/primes/prime-kanban-board.md`; do not bump `VERSION`.

## P-A-U phasing (mandatory, reported in the hand-back)

Report your P-A-U record under a `## P-A-U` heading in the hand-back: **[PLAN]** brief approach before code; **[APPLY]** code exactly as planned, inside the write set; **[UNIFY]** `git diff --stat`, `gofmt -l .`, `go vet ./...` in the queue-kanban module, a debug-artifact scan over added lines, and the list of files checked.

## Focused tests and proof

Every test-file invocation must finish in under 30 seconds. From the repo root of your worktree:
- Node lane: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...`
- Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`

`tdd: false`, but the REQ's `## Red-Green Proof` names the assertions and the Node lane can run them: write the new Node-lane case first (outer details closed by default whose summary names every subject; each row a details with dot, subject, category words, detail in its summary and the remedy in its content; the skipped probe a "not checked" row; the stored key rendering the strip open), watch it fail on the shipped renderer, then make it pass.

The look is a render fact: generate the board into a scratch directory, serve it on a port other than 8090, screenshot the strip closed, open, and with one remedy open, at 1600 px in light and dark (`"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless=new --disable-gpu --hide-scrollbars --window-size=1600,360 --screenshot=<png> <url>`; force a theme with `--blink-settings=preferredColorScheme=1` for light and `=2` for dark), compare each against the corresponding gallery capture in `ai-reports/2026-09-05_1800_REQ-588-verify-findings-slim-band-gallery/screenshots/`, look at the images, and record what you saw and the PNG paths in the hand-back. The live queue has at least one finding right now.

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-182000/REQ-589-handback.md`** using that absolute path — the one main-tree path you may write; never stage or commit it. Sections, each under its own `##` heading: `Branch`, `File manifest`, `P-A-U` (with `### [PLAN]`, `### [APPLY]`, `### [UNIFY]`), `Test evidence` (commands, exit codes, RED then GREEN excerpts, screenshots and what they show), `Lesson evidence`, `Decisions` (D-XX from D-01, with reasoning; Value/Risk lines on any you would escalate), `Discovered Tasks`, `Integration seams`.

Commit on your branch with a message starting `[REQ-589] `, ending with:
Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01RfyDKukKGpyRY6F2UFonVg
Then write the hand-back and stop. Do not merge, do not touch main.
