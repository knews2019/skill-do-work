---
id: REQ-579
title: 'Render verify findings and skipped probes as compact rows in one list'
status: completed
created_at: 2026-09-05T00:19:58Z
user_request: UR-118
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T12:22:40Z
  basis:
    - trivial short-circuit
related: [REQ-580, REQ-578]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
claimed_at: 2026-09-05T12:21:55Z
route: B
dispatch_at: 2026-09-05T12:23:39Z
builder_handback_at: 2026-09-05T12:46:22Z
integration_at: 2026-09-05T12:46:22Z
review_at: 2026-09-05T14:23:29Z
heavy_verified_at: 2026-09-05T14:23:29Z
heavy_verified_revision: 7b2673b690a671ccb360c26b0c19c56ecc7356b5
commit: b169396e
completed_at: 2026-09-05T14:24:17Z
release_at: 2026-09-05T14:24:17Z
---

# Render Verify Findings and Skipped Probes as Compact Rows in One List

## What

The Verify Findings strip (`#board-findings`, added by REQ-285) renders each finding as a bordered card in a multi-column grid and lists skipped probes as bullets inside a collapsed disclosure below. Replace both with one flat list of compact rows: one row per finding and one row per skipped probe, no cards, no columns, no separate disclosure. A finding and a skipped probe are the same kind of thing to the reader ("verify has something to tell you"), so they get one visual language.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Five steps settled before code: the producer subject field and its payload mirror, the client's flat row list with subject grouping, the markup deletion, the CSS replacement, then RED in both lanes before implementing and a real render afterwards. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-120117/REQ-579-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`1dc13ef7`) touching exactly the seven Scope files.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `gofmt -l .` empty; `go vet ./...` clean; debug-artifact scan over added lines empty; per-file checks listed in the hand-back.

## Why

The user's words: "for a warning I like the small lines, the big boxes are using up too much space for something not that important". The cards also borrow the REQ-card shape, so they read as work items, which REQ-482 (stack verify-finding cards full width, now cancelled in favour of this REQ) had already noticed. The contradiction being removed is two visual languages for one idea: a bordered card for a finding, a bullet in a disclosure for a skipped probe.

## Detailed Requirements

- **D1, one list of rows.** The header line stays: "Verify findings", the finding count, and the hint. Add the skipped count to the header when it is non-zero (for example "2 findings · 2 probes not checked"). Below the header, one plain list. Each finding row carries: a small uppercase category chip (as today's `.board-finding-category`), the detail sentence, and the remedy in muted text on the same line after the detail. Long text wraps. No border, no card background, at most a thin divider between rows. Each skipped-probe row carries a "not checked" chip and the probe text in muted weight, in the same list, after the findings. Delete the `<details>` disclosure and the `board-anomalies-cards` grid from the strip's markup; the `.board-finding*` and `.board-findings-skipped*` classes are replaced by the row classes, not kept beside them.
- **D2, two weights, from the producer only.** A finding with `fixable: true` renders with a muted "cleanup can fix" tag and lighter text, because a command resolves it. Every other finding renders at normal weight. A skipped probe renders muted. Do not add any severity, colour scale or ordering in the client that the payload does not carry; `fixable` is the only weight (the same rule as today's comment in `renderVerifyFindingsStrip`: never inferred here, the producer sets the flag).
- **D3, grouped by subject.** Add a `Subject` string to `VerifyFinding` in verify.go and a `subject` field to the board payload in generate.go. The producer sets it to the thing the finding is about: the worktree name for the worktree probes, the REQ id for REQ-level probes, the changelog path or version for the release probes. A probe with no natural subject leaves it empty. The client groups rows with the same non-empty subject together and prints the subject once as a small heading line above its rows; rows with an empty subject render as single rows in producer order, after the grouped ones. Grouping is by exact string match on the payload field, never by parsing the detail text.
- **D5, hide rules unchanged.** The strip stays hidden when there are no findings and no skipped probes. REQ-578 (hide the strip on the Activity view) keeps working against the same `#board-findings` element and `hidden` attribute; do not restructure the strip in a way that breaks it.
- The strip is not a REQ card. Do not reuse `board-anomalies-cards` or any `.board-request*` class for finding rows.
- Board changes follow `_dev/primes/prime-kanban-board.md` (versioning, parser lock-step, build outputs). The generated payload gains one optional field, so the snapshot fixtures under the queue-kanban tests need regenerating per that prime.

## Constraints

- REQ-580 (stop probing committed queue state for an undetermined worktree) edits the same probe in verify.go and its test. Both REQs may add a field or a branch around `appendWorktreeFindings`; declare the overlap, do not serialise on it.
- The completion-anomalies strip above the findings strip is not part of this request and keeps its cards.

## Builder Guidance

The user is certain about the outcome (rows, not cards) and approved D1 to D3 and D5 as written. Latitude: spacing, divider style, chip typography and the exact header wording are the builder's call, as long as the strip reads as a list of warnings and not as a set of cards. Load `crew-members/anti-slop.md` before styling; the result is a human-facing artifact.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane (`javascript_behavior_c_test.go` pattern), render with two verify findings sharing `subject: "worktree-agent-REQ-506-focused-evidence"` (one `fixable: true`, one not) and one skipped probe in `boardData`. Assert: `#board-findings` contains no `.board-finding` card and no `details` element; it contains exactly three row elements in one list; the two findings sit under one subject heading carrying that worktree name; the fixable row and the skipped row carry the muted class and the other finding row does not. In Go, assert that `generatedVerifyFinding` serialises a `subject` field for a worktree leftover finding.
**Why RED now:** `renderVerifyFindingsStrip` in board-cards.js builds one `.board-finding` card per finding inside the `board-anomalies-cards` grid and puts skipped probes in a `<details>` list; `VerifyFinding` has no subject and the payload carries none.
**GREEN when:** The assertions above pass, the existing strip hide-when-empty test still passes, and the board screenshot shows the strip as a flat list of rows with the skipped probes in it.
**Validation:** User confirmed: "for a warning I like the small lines" and "ok, do it" to the D1 to D5 proposal; screenshot 1 shows the RED state.

## Assets

- `do-work/user-requests/UR-118/assets/REQ-579-screenshot-1-verify-findings-two-cards-and-skipped-probes.png`: queue-kanban board, light theme, Board view, 2026-09-05. Strip header "VERIFY FINDINGS 2". Two white bordered cards side by side, both labelled WORKTREE-MERGE-STATE-UNDETERMINED, for `worktree-agent-REQ-506-focused-evidence` and `worktree-agent-REQ-577-launcher-fixture`, each with the detail "git could not say whether this is merged ... inspect it by hand; cleanup Pass 5 deletes nothing it cannot establish a merge target for". Below them an expanded disclosure "2 probe(s) could not run — unverified, not clean" with two bullets, each "committed-queue-state probe for <worktree>: `git diff main...<branch> -- do-work/` failed (no such branch, or unrelated histories)". The cards use roughly a third of the strip's width each and most of its height; the two bullets take two lines.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views" and "static output". Over the budget on its own.

*Source: "for a warning I like the small lines, the big boxes are using up too much space for something not that important" / "how to make it good, beautiful and non-contradictory?" / "ok, do it"*

---

## Triage

**Route: B** - Medium

**Reasoning:** Substantive effort across Go producer, assembled client, CSS and two test lanes, with real styling judgment in how the strip reads. `effort_estimate: effort-substantive`.

**Planning:** Not required — the REQ carries requirements D1, D2, D3 and D5 as an approved design, so a plan agent would restate them.

**Exploration:** Satisfied by the request's own capture-time findings, which name the exact producer (`VerifyFinding` in `verify.go`), payload site (`generatedVerifyFinding` in `generate.go`), renderer (`renderVerifyFindingsStrip` in `board-cards.js`), the classes being replaced, and the Node-lane test pattern — everything a separate Explore agent would have returned. The builder re-derives anything the capture missed and reports it under `## Exploration` in its hand-back. Recorded as an orchestrator judgment so the skipped step is visible rather than silent.

## Plan

**Planning not required** - Route B: exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modify)

The producer gains a subject field set by each probe and mirrored into the board payload; the client replaces the card grid and the disclosure with one flat row list; the stylesheet replaces the card and disclosure rules with row rules; both test lanes gain a case.

**Files I will NOT touch:** the completion-anomalies strip and its cards, `web/board-controls.js` (REQ-578 owns the Activity hide rule and it must keep working against the same `#board-findings` element), `CHANGELOG.md`, `VERSION`, everything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] One flat list of rows: no `.board-finding` card, no `details` disclosure, header carries the finding count and the skipped count when non-zero
- [ ] Two weights only, both from the producer: `fixable: true` renders muted with a "cleanup can fix" tag, a skipped probe renders muted, everything else normal
- [ ] `Subject` exists on `VerifyFinding` and `subject` in the payload; the client groups rows by exact string match and prints each non-empty subject once as a heading, with empty-subject rows after the grouped ones in producer order
- [ ] The strip still hides when there is nothing to report, and REQ-578's Activity-view hide still works against the same element and attribute

## Exploration

The builder produced this section in its hand-back (`do-work/runs/work-2026-09-05-120117/REQ-579-handback.md` → `## Exploration`) in place of a separate Explore agent, per the Triage note above. Read it there; the material findings it changed are recorded in Qualification below.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)

**What was done:** Findings now carry a subject set by the probe that knows it — the worktree name, the request id, the changelog path — mirrored into the board payload. The card grid and the collapsed disclosure are gone, replaced by one flat list where a finding row and a skipped-probe row share a shape, rows with the same subject group under one heading, and subjectless rows follow in producer order. Two weights only, both from the payload: a fixable finding is muted and tagged, a skipped probe is muted. Merge range `4362ac0d..c78a0d3d`; builder branch head `1dc13ef7`. Builder-authored `## Decisions` (D-01 to D-06), `## Discovered Tasks` and `## Exploration` live in the hand-back.

## Qualification

**Passed.** Read from the merge range `4362ac0d..c78a0d3d`, and the module suite re-run on the merged tree.

- **The trap this request was warned about was not walked into.** REQ-578's rule reads `#board-findings-cards` and `#board-findings-skipped-list` children to decide whether the strip has anything to say, and this request's own D1 said to delete the disclosure that holds the second one. The builder built the honest single-host version first, watched REQ-578's test fail, and reverted to nesting both hosts (as `display: contents`) inside one `#board-findings-rows` list rather than handing back a knowingly-red test with a seam someone had to remember. The page is one list to the reader; the two ids survive as pass-through hosts with the naming wart filed as a follow-up. It also wrote the `display: contents` trap into the template comment beside the ids, since that property outranks the `hidden` attribute's own display rule.
- Two weights, both from the producer. Nothing in the client invents a severity, an ordering or a colour scale the payload did not carry.
- Grouping is exact string match on the payload field, never parsed out of the detail text, and the lookup uses a null-prototype object so a subject spelled `constructor` cannot be mistaken for an already-open group.
- The subject is set at every probe that has one, including categories the board currently suppresses, so the rule is "the probe that knows, sets it" rather than a list to maintain. Findings lifted from warning prose with no parsed id stay empty, as designed.
- One spacing rule (`.board-findings-row-detached`) exists because the builder looked at the render and saw a subjectless row reading as the previous group's third row. That is the kind of defect only a render finds.

Requirements traced: one flat list with no card and no disclosure (D1), header carrying both counts with nouns on each half (D1), two producer-only weights (D2), subject on the model, in the payload and grouped in the client (D3), and the strip still hiding when there is nothing to report with REQ-578's Activity rule intact (D5).

Scope: the seven touched files are exactly the declared write set. No drift.

*Checked by work action*

## Testing

**Post-merge, main tree at `c78a0d3d`:**
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` — exit 0, 453 tests, wall 57s, slowest file `generate_test.go` at 11.54s against the 30s per-file budget. That run covers both new tests and REQ-573's, which met this branch in a merge conflict.

**Red-green validation** (traced to `## Red-Green Proof`): RED in both lanes on assertions, not on compile errors — the Go lane reported every finding's subject as empty and the payload carrying none, the Node lane reported the header as a bare count and the row list empty. Each lane's first attempt died on something that was not an assertion (an undeclared field, a stub lookup) and was fixed before the RED was taken, which is the right discipline. GREEN in both lanes afterwards, including REQ-578's hide-on-Activity test driving the real view switcher against the new markup.

**Render evidence (a browser was actually driven):** the generated board was served over HTTP and read in Chromium, twice — once against this repository's real queue (five findings, five distinct subjects, one heading each) and once with a crafted payload, which is the only way to see grouping, muting and skipped rows together. The crafted render measured the header as "6 findings · 2 probes not checked", 8 rows under 3 subject headings, zero `details` elements, zero card elements, both hosts as direct children of the list, and a strip 398px tall for 8 rows where the old cards used most of the strip's height for two.


## Review

**Overall: 96%** | 2026-09-05T12:58:03Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Reviewed over the full range `4362ac0d..b169396e`, which is two merges: `c78a0d3d` (builder commit `1dc13ef7`) and `b169396e` (builder commit `1de64c4c`). An earlier version of this section covered only the first merge and reported the second commit as unintegrated; it is now in main, so that finding is resolved and is not carried below. The follow-up merge also carried unrelated UR-123 capture commits on the main side; those are not this REQ's work and are not scored here. The two REQ-579 source files in the merged tree are byte-identical to the builder's `1de64c4c`.

**Conflict resolution (checked first).** The merge `c78a0d3d` resolved `javascript_behavior_c_test.go` without losing either side. REQ-573's `TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest` is byte-identical to the main-side version at `4362ac0d` (207 lines, every assertion intact); REQ-579's `TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList` is byte-identical to the builder-side version at `1dc13ef7` (218 lines); REQ-578's `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` took the builder's edited version, which is the intended cross-REQ change (the `skippedDisclosureOnTimeline` assertion read the deleted `<details>` element and would have passed vacuously against the stub `document` forever). The second merge added only the follow-up's 39 test lines and 22 renderer lines; nothing was re-resolved.

**The follow-up, judged.**

*Does clearing before the empty check cover every path?* Yes. The renderer has three exits. The `!strip` guard returns before the clears, which is the only safe order: if `#board-findings` is absent from the template the two hosts are absent too, so clearing first would throw. The empty-payload exit now clears both hosts and then hides the strip. The populated path clears both once and refills them — the follow-up also deleted the two per-host clears further down, so each host is cleared exactly once on every path that can reach it. One behavior change the commit message does not name: the empty path now dereferences both hosts where it previously returned before touching them, so a template carrying `#board-findings` without its two hosts would throw at boot on a board with nothing to report. That state cannot ship — `generate.go:23` embeds `web/template.html`, `web/board.css` and `web/*.js` together and inlines the client into the page, with only `board-data.js` served separately, so template and renderer always travel as one artifact. Fail-fast here is safe.

*Is the markup-id assertion strong enough?* Yes, and it fires on exactly the blind spot it was written for. Verified by mutation in a scratch copy of the module (`do-work/` fixture root, module copied out of the main checkout, main tree never edited):

- Dropping `skippedHost.textContent = ""` → FAIL `javascript_behavior_c_test.go:3112: stale rows survived the empty re-render (findings=0, skipped=1)`.
- Dropping `findingsHost.textContent = ""` → FAIL, same assertion, `(findings=4, skipped=0)`.
- Renaming both host ids in `web/template.html` only, leaving the renderer and the probe on the old ids — the case where the stub invents a node for the old id and every DOM assertion passes — → FAIL `javascript_behavior_c_test.go:3139: the strip dropped id="board-findings-cards" — applyView reads it to decide visibility and would dereference null`, and the same for the skipped host. Nothing else in the test fails. This is the assertion doing the exact job the stub cannot.
- Renaming a host id consistently in both `template.html` and `board-cards.js` → FAIL, but at `javascript_behavior_c_test.go:3050` ("the list holds N children"), which is a `t.Fatalf` and therefore kills the test before the markup-id checks run. See the Minor finding on message shadowing.

The consumer half of the contract — `applyView` reading those two exact ids in `web/board-controls.js:57-58` — is covered by REQ-578's own test, whose findings-only and skipped-only scenarios each leave one host empty, so a wrong id read there flips the visibility assertion rather than passing vacuously.

**Restatement Sweep.** Swept every element the diff redefines: the deleted classes (`.board-finding*`, `.board-findings-skipped*`), the deleted element ids (`board-findings-skipped`, `board-findings-skipped-count`), the deleted `board-anomalies-cards` host in this strip, the retired "probe(s) could not run" footer wording, and the payload shape (`generatedVerifyFinding` gaining `subject`). No stale restatement survives in a shipped file. The one shipped sentence that did go stale — `routeWorktreeLeftover`'s `worktreeLeftoverStateUnknown` remedy, which pointed the reader at the deleted footer — was rewritten in the same commit and now names the `not checked` row. Everything else that still describes the old strip is a historical record that is never rewritten: archived REQ-285/REQ-482, `kb/wiki/sources/`, and `skills/do-work/CHANGELOG.md:1920` (the entry recording what REQ-285 shipped). The REQ's own line about regenerating snapshot fixtures was moot — the queue-kanban module has no `testdata/` or golden files, and `subject` is `omitempty`. REQ-580's overlap in `appendWorktreeFindings` (the undetermined-merge-state short-circuit and its `default` remedy) survived both merges intact with REQ-579's `Subject` fields added beside it.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- Light theme, `.board-findings-fixable` (`web/board.css:662`): the "cleanup can fix" tag is `--accent-done` `#3c875e` on the strip's `--surface-2` `#f1f4f8`, measuring 3.95:1 at 0.82rem — below 4.5:1 for small text. It was 4.36:1 on the old white card, so the move to the tinted strip cost 0.41 and it was already short. Dark theme is fine at 6.92:1. No information is lost either way: the words "cleanup can fix" carry the meaning, so the tag is distinguishable without colour. — impact-user-visible → report only
- The skipped probes lost their list semantics: they were `<ul>`/`<li>` and are now `<div class="board-findings-row">` with no `role="list"`/`role="listitem"`, and the subject heading (`board-findings-subject`) is a `<div>` rather than a heading element. A screen reader no longer announces "list, N items" nor offers heading navigation over the subjects. The section keeps `aria-label="Verify findings"`, so nothing is unreachable. Adding roles needs a browser check, because the two hosts carry `display: contents` and that changes how the list/listitem relation is exposed. — impact-user-visible → report only
- The two markup-id assertions (`javascript_behavior_c_test.go:3131-3145`) sit after a `t.Fatalf` on the row count (`:3050`), so on the most likely future edit — a consistent rename across `template.html` and `board-cards.js`, which is exactly the follow-up the builder filed under Discovered Tasks — the test fails with "the list holds N children" and a node dump, and the message that names `applyView` never prints. Coverage is not affected; the diagnostic is. Both checks read only `indexHtml`, which is in hand at the top of the test, so moving them above the DOM read-back costs nothing. — impact-negligible → report only
- The `display: contents` versus `hidden` trap is real and correctly stated (an author `display` rule outranks the UA stylesheet's `[hidden] { display: none }`, so `hidden` on either host would not hide anything), and it is documented three times at the point of use — `web/template.html:216-218`, `web/board.css:608-610`, and now the renderer comment. It is pinned by nothing: the shipped test asserts the `.board-findings-group { display: contents }` rule exists, but nothing fails if a future editor reintroduces `skippedHost.hidden = true`, and the Node lane's stub sets `hidden` as a plain property with no CSS to lose against. Documentation-only guard on a trap the previous renderer walked into. — impact-negligible → report only
- The hand-back's file manifest (`do-work/runs/work-2026-09-05-120117/REQ-579-handback.md:19`) cites "(see D-07)" but its Decisions section stops at D-06, so the decision behind the follow-up commit is named but never recorded. Now that the commit is merged, that decision is part of the delivered work and has no entry. — impact-negligible → report only

**Nit findings:**
- `board-findings-detail` is applied to every row's main text (`web/board-cards.js:757`, `:775`) but has no CSS rule and no test consumer — a styling hook that styles nothing. The builder recorded the reasoning (an unclassed span among classed ones reads as an oversight), which is a fair call; noted only so the next editor knows it is deliberate. — impact-negligible → report only

**Acceptance:** Pass — on the merged tree at HEAD `b169396e`, `TestVerifyNamesTheSubjectEachFindingIsAbout`, `TestVerifyNamesTheChangelogAsTheReleaseFindingSubject`, `TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList`, REQ-578's `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` and REQ-573's `TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest` all pass with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1` (exit 0). The row-list probe drives the shipped renderer over the generated page and asserts the markup half against the generated HTML, so the "no cards, no disclosure, one list, hosts still present" claim is checked against what ships. Five mutations were run against a scratch copy to confirm the new assertions fail for the right reasons. Contrast was computed from the palette tokens rather than measured in a browser; the browser render was not independently reproduced (the brief forbids driving one).
**Suggested testing:** 5 items — (1) open a real board and confirm the strip reads as one list at the widths the board targets; (2) a board with 20+ findings, since the collapse the `<details>` gave is gone and the strip now grows without limit; (3) dark theme by eye, where the muted weight and the chip carry more of the distinction; (4) a screen-reader pass over the strip, which is where the lost list semantics show; (5) after the host-id rename follow-up lands, re-run REQ-578's hide test in a browser, not only in the Node lane.
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Heavy Verification Result

- **Target revision:** b169396e
- **Execution revision:** 7b2673b690a671ccb360c26b0c19c56ecc7356b5
- **Run at:** 2026-09-05T14:23:29Z, from a detached worktree (the shared main tree carried other sessions' uncommitted work, which a lane result must not be attributed to)

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| `queue-kanban-javascript` | 0 | 9s | executed |
| `queue-kanban-browser` | 0 | 141s | executed |
| `staged-skills` | 0 | 44s | executed |

Every lane this request selected was present in the run, exited 0, and none was skipped.

## Lessons Learned

**What worked:** Building the honest version first and letting the test say no. The single-host markup this request's own requirements implied broke the rule another request had just shipped, and the builder found that by running the suite rather than by reasoning about it — then chose the naming wart over handing back a knowingly-red test with a seam someone had to remember to apply. Looking at the rendered page also earned its keep: the extra spacing rule exists because a subjectless row read as the previous group's third row on screen, which no assertion would have caught.

**What didn't:** Deleting an element id that another file reads. The disclosure this request was told to delete held one of two ids a view rule dereferences on every view switch, and the Node lane could not have caught it — its stub `getElementById` manufactures a node for any id it is asked for, so a deleted element reads as present-with-zero-children and the tests stay green while the browser throws.

**Worth knowing:** Two follow-ups landed because the neighbouring request's review looked at this one's plan before it was built. The renderer's early return left stale rows under a hidden strip, which only matters because the view rule now asks those hosts whether there is content — a latent defect that a second reader turned into a caught one. Reviews that read the queue, not only the diff, are where that comes from.

## Orientation

The board's Verify Findings strip is a list of warnings again rather than a set of work items: one compact row per finding and per skipped probe, grouped under the thing each is about — a worktree, a request, the changelog. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`), with the subject set by the Go probe that knows it and carried through the board payload. [MAP CHANGED] — a finding now carries a subject, which is a new field on the producer's contract and the first thing the client groups on. On a real board the strip is 398px for eight rows where two cards used most of its height.
