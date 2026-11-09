---
title: "Lessons from REQ-040: Board overlap badge: use path.Match and document the glob dialect"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-040-board-overlap-badge-use-path-match-and-d.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-040: Board overlap badge: use path.Match and document the glob dialect

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

REQ-034's `writeSetPatternsIntersect` uses `filepath.Match`, which is OS-dependent (`*` would cross `/` on Windows because the separator is `\`); `write_set` entries are slash-separated repo-relative paths, so `path.Match` is the correct drop-in primitive. Additionally the glob dialect is documented nowhere: `**` silently under-matches (`filepath.Match("src/**/x.ts", "src/a/b/x.ts")` is false) and a malformed pattern (`ErrBadPattern`) silently degrades to literal-only matching — while the pipeline gate treats an unexpandable glob as overlapping. Fix the primitive and state the dialect where readers meet the field.

## Solution summary

Swapped `writeSetPatternsIntersect`'s two `filepath.Match` calls to `path.Match` (added a `"path"` import alongside the retained `"path/filepath"`, which is still used by `filepath.Base`/`ToSlash`/`Dir` elsewhere in the file) so the board's overlap annotation is OS-independent — `write_set` entries are slash-separated repo-relative paths, and `filepath.Match`'s `\` separator on Windows would let `*` wrongly cross `/`. Rewrote the doc comment to name `path.Match`, give the rationale, and state the dialect (`*` never crosses `/`; `**` is not recursive — no superglob; malformed ⇒ `ErrBadPattern` ⇒ no-match), plus a display-only note that the pipeline's dispatch gate still serializes an unexpandable/overlapping glob, so a board false-negative never loosens the gate. Added two tests (`…GlobStarNeverCrossesSlash` with a same-segment positive control, and `…MalformedPatternMatchesNothing`) in the existing fixture style. Documented the same dialect in `actions/board.md`'s badge paragraph, `tools/queue-kanban/prime-do-kanban.md`'s lesson, and one clause on `actions/work-reference.md`'s `write_set` schema line. No frontmatter/schema shape change, so the board parser lock-step is unaffected.

## What worked

- Keeping `path/filepath` (still used by `.Base`/`.ToSlash`/`.Dir`) and adding `"path"` alongside avoided the obvious build-break of "swap the import." The positive-control assertion in the slash test (`web/*` DOES match `web/app.css`) makes the test prove the boundary is the slash, not the glob silently failing.

## What didn't work

- Nothing surprising — the anchors and the import situation were confirmed before the build, so it landed first try.

## Worth knowing

- `path.Match` has no recursive `**` superglob and returns `ErrBadPattern` on malformed input — so the board's overlap badge is a best-effort *visualizer*, not the dispatch gate. The safety-relevant reader is the work pipeline's gate, which treats an unexpandable/overlapping glob as overlapping (serialize); a board false-negative never loosens it. `write_set` entries must stay slash-separated repo-relative paths for `path.Match` to be correct.

## Back-reference

See `do-work/archive/UR-007/REQ-040-overlap-badge-glob-dialect.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `acc4722`.
