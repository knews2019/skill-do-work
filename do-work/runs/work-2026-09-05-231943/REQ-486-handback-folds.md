# Hand-back — REQ-486 (collapsible UR progress summaries), increment 1 of 2: the folds

**Branch:** `worktree-agent-REQ-486-collapsible-ur-progress-summaries`
**Head:** `1bb3f162534b74b619145f0d352067fd345fc4c6`
**Base:** `a2497c6`
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-486-collapsible-ur-progress-summaries`
**Scope done:** T1 and T2 only. Nothing from T3, T4 or T5 was started. No Go change, no payload
field, no rollup, no clock change, no summary-strip CSS, no doc edit, no release.

## Commits

| Code | Commit | What it is |
| --- | --- | --- |
| C1 | `e4168dd` | T1 — By UR groups fold; the head's two-branch shape is deleted |
| C2 | `1bb3f16` | T2 — the UR drawer's "REQ ids" list folds and is height-capped |

Both commits end with the required `Co-Authored-By` line. Nothing is pushed. No worktree was
removed. The only file written in the main checkout is this hand-back, unstaged and uncommitted.

## Files changed

- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` — one shape of group head for both
  readings.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` — `appendFoldableMetaRow` beside
  `appendMetaRow`; the "REQ ids" row uses it.
- `skills/do-work-board/tools/queue-kanban/web/board.css` — `.detail-fold`, `.detail-fold-marker`,
  and the 40vh cap scoped to `.detail-foldable-value .detail-dep-list`.
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` — new, two probes.
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` — the inverted assertion.
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` — forced stub repair.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` — `renderedUserRequestRow` gains
  `DetailButtonId`; the citation probe's selector is repaired.
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` — the
  drawer trigger selector is repaired.

## Decisions carried out

- **D2** — the fold is a net deletion in `renderUserRequestLens`. The two head shapes are gone.
  `userRequestCardsFolded` now decides exactly two things: the initial `aria-expanded` value, and
  whether the card grid is built eagerly or on first open. Same wrapper, same sibling Details
  button, same click listener in both readings.
- **D3** — `javascript_behavior_a_test.go`'s assertion that the By UR head carries **no**
  `aria-expanded` was **inverted in place**, in the same commit as the markup change, never deleted.
  It now pins both defaults in one assertion: By UR starts `"true"` with cards present, URs only
  starts `"false"` with none. The failure message names REQ-486 as the inversion.
- **D9** — the drawer id list starts expanded and gets `max-height: 40vh; overflow-y: auto`. The cap
  is scoped to the foldable row, so "Depends on" and "Blocked by" in the REQ drawer keep their
  natural height. Unscoping it to `.detail-dep-list` would have changed two rows nobody complained
  about.
- **D15** — new probes went into the new file `javascript_behavior_d_test.go`. Only the one
  assertion named by D3 was edited in place.

## TDD evidence

Green baselines taken **before the first edit**, both heavy lanes at `a2497c6`:

```
JavaScript lane:  ok  github.com/knews2019/skill-do-work/queue-kanban  6.011s   63 PASS, 0 SKIP
browser lane:     ok  github.com/knews2019/skill-do-work/queue-kanban  96.696s  34 PASS, 0 SKIP
```

**RED for T1** (tests written first, production untouched), exit status 1:

```
--- FAIL: TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened (3.40s)
    javascript_behavior_a_test.go:1664: fold defaults = by-UR "" with 1 cards, URs only "false"
        with 0 cards; REQ-486 wants by-UR open ("true", cards present) and URs only collapsed
        ("false", no cards)
--- FAIL: TestJavaScriptBehaviorByUserRequestLensFoldsAndRestoresItsCards (0.03s)
    javascript_behavior_d_test.go:167: by-UR group UR-401 starts aria-expanded="" with cards
        []string{"REQ-601", "REQ-602"}; want "true" and its cards present
```

**RED for T2**, exit status 1:

```
--- FAIL: TestJavaScriptBehaviorUserRequestDrawerFoldsItsRequestIdList (3.16s)
    javascript_behavior_d_test.go:215: anchor "function appendFoldableMetaRow(" not found in the
        generated page
```

That second red is an anchor red — it proves the helper did not exist, not that the fold behaves.
The behavioural strength sits in the assertions that ran after the helper landed: the list's ids,
the `hidden` flip, the neighbouring rows staying visible in the same pass, and the body still
carrying its text.

**GREEN**, each lane's own exit line, taken at head `1bb3f16`:

```
fast stage:       ok  github.com/knews2019/skill-do-work/queue-kanban  60.824s
JavaScript lane:  ok  github.com/knews2019/skill-do-work/queue-kanban  7.691s   65 PASS, 0 SKIP
browser lane:     ok  github.com/knews2019/skill-do-work/queue-kanban  94.899s  34 PASS, 0 SKIP
```

The browser lane ran with `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium` on **Chromium
141.0.7390.37** (`--version`; the page's own user agent reports the reduced `Chrome/141.0.0.0`).
It printed no SKIP line. Chromium 141 is deprecated per the board prime, so this green is evidence
about this run, not a compatibility claim.

## Three forced repairs — the plan named two

| Code | Site | Repair |
| --- | --- | --- |
| F1 | `user_request_clipboard_browser_probe_test.go:142` | selects `.ur-group-detail[data-detail-id=…]` instead of the head. Exactly as planned. |
| F2 | `generate_test.go:1037` (citation probe) | **diverges from the plan**, see below. |
| F3 | `javascript_behavior_c_test.go:518` and `:535` | **the plan missed this one entirely.** |

**F2 — why the plan's literal repair would have failed.** The plan says both selectors become
`.ur-group-detail[data-detail-id="…"]`. That is right for F1, which clicks the drawer trigger. It is
wrong for the citation probe: that probe wants the **head**, because the "cites" badge is appended
inside the head, and it then asserts the badge's rect is contained in the selected element's rect
and that the element's accessible name mentions the citation. Pointing it at the Details button
would make `button.querySelector('.citation-match')` return null and fail three assertions that
have nothing to do with this REQ. The head is still the right node; it just no longer carries
`data-detail-*`. It is now named by the button beside it:

```
.ur-group-row:has(.ur-group-detail[data-detail-id="UR-075"]) .ur-group-head
```

`:has()` has been in Chromium since 105, and the lane pins the browser at 141.

**F3 — the failure the plan did not predict.** `javascript_behavior_c_test.go`'s by-UR probe had a
DOM stub with only `childNodes`, `dataset` and `appendChild`, because the always-open reading
previously set no attribute and attached no listener. After T1 it does both, and the probe died
with `TypeError: head.setAttribute is not a function`. It also read the UR id from
`groupNode.childNodes[0].dataset.detailId`, which is now the `div.ur-group-row` and would have
reported `null`. The stub gained `setAttribute`, `addEventListener` and `removeChild`, and the id
now comes off the Details button. Both changes are in C1, so the JavaScript lane never had a
window in which it was red. Only the fast stage would have stayed green through it, which is
exactly the trap the brief warned about — it just had a third site.

## One planned change I deliberately did not make

The plan puts `.ur-count` losing its total ("12 of 43 shown" only under a filter, D11) inside T1,
with the total moving to the summary strip in T4. The strip is T4 and is not in this increment.
Making the change now would delete the group's REQ count from the board and put nothing in its
place until the second builder's work lands. `.ur-count` is therefore untouched and still renders
`43 REQ` unfiltered. **The second builder should make the D11 change together with the strip**,
in the same commit — my brief scoped T1 to the branch deletion and did not name `.ur-count`.

## Canonical gate

Run twice, as budgeted.

- **Run 1 exited 1.** One failure, in `_dev/tests/prescribed-shell-cases/qualify.sh`: `qualify
  untracked case missed the TODO — the unfinished-marker scan still does not read untracked files`.
  Nothing in this increment touches `skills/do-work/scripts/checks/`.
- Investigated rather than assumed. `_dev/tests/staged-skills-contract.sh` run alone with
  `DO_WORK_MAINTAINER_TIER=heavy` reports `qualify: 26 cases, 0 failures` **both** in a clean clone
  of the base commit `a2497c6` and in this worktree at `1bb3f16`. The two trees behave identically,
  so the failure is not this branch's.
- **Run 2 exited 0**: `Maintainer verification passed.`, gate wall 304s, with both heavy lanes
  present rather than skipped (`queue-kanban uncached tests with strict JavaScript behavior probes`,
  467 tests; `queue-kanban strict browser behavior lane`, 34 tests).

Read this as: the qualify case is flaky under the full parallel gate on this 4-CPU machine, and it
is a real thing somebody should look at — it is not caused by REQ-486 and I did not chase it
further. It is worth a discovered-task note if the same case flakes again in the second increment.

## What nothing here proves

- **Tab order.** Focus movement is a trusted-input default action. Both fold controls and the
  Details control are real `<button>` elements, which is what delivers it; the JavaScript lane only
  proves `aria-expanded` flips and the value node's `hidden` toggles. The browser probe that would
  drive real Tab keys belongs to T5.
- **Layout at narrow widths.** No pixel assertion was added in this increment. The 40vh cap and the
  drawer containment claim are T5's browser probe.
- **The drawer fold's visual affordance.** I added a `.detail-fold-marker` glyph (rotates on
  collapse) so the row does not read as a plain label. The plan asked only for the appearance reset
  and the focus outline; the marker is one span and two CSS rules, and it is trivially removable if
  the maintainer disagrees.

## Notes for the second builder

1. Rebase or branch from `1bb3f16`. The markup T4's strip attaches to is final: the strip is a
   sibling of `div.ur-group-row` inside `section.ur-group`, and that wrapper now exists in both
   readings, so `renderUserRequestLens` has one place to append it.
2. `appendFoldableMetaRow(label, valueNode, valueNodeId)` is generic. T4's drawer metric rows can
   use `appendMetaRow` as they are; only pass through the foldable one if a metric row should fold.
3. `javascript_behavior_d_test.go` exports `foldProbeDomStub`, a Go `const` holding the shared
   `makeNode()` / `collectByClassName()` stub with `hidden`, `id`, `classList` and a `textContent`
   setter that really empties a node. T4's probes can concatenate it instead of writing a third
   copy.
4. The drawer probe stubs `linkifyDetailBody`, `renderDetailGlossary`, `showDrawer` and
   `syncActivitySelectionToDrawer`, and drives the **shipped** `setDetailTarget`. Keep that split
   if you extend it.

Both folds ship, all three lanes are green at `1bb3f16` with their own exit lines recorded, and the
canonical gate passes; the only thing the second builder inherits from T1 is the `.ur-count`
denominator, which must move to the summary strip in the same commit that adds the strip.
