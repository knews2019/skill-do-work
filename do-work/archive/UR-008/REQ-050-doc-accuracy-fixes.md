---
id: REQ-050
title: "Three doc-accuracy fixes: session-start 'self-clears' comment, malformed-glob phrasing, directory-entry miss-class"
status: completed
route: B
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T13:45:26Z
completed_at: 2026-07-29T13:51:45Z
commit: 548e487
user_request: UR-008
domain: general
prime_files: []
tdd: false
depends_on: []
related: []
batch: deep-review-followups
write_set: [hooks/memory-session-start.sh, tools/queue-kanban/model.go, actions/board.md, actions/work-reference.md, tools/queue-kanban/prime-do-kanban.md]
maintenance: false
---

# Doc-Accuracy Fixes from the Deep Review

## What

Three confirmed documentation inaccuracies — all comment/prose only, no behavior changes. Grouped because each is a one-to-four-line wording fix.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files: []` — read `crew-members/general.md` + `crew-members/coding-guardrails.md`. Approach: trace all three behaviors against the real code (awk state machine; `writeSetPatternsIntersect`'s literal-equality short-circuit; `path.Match` on directory entries) BEFORE editing, then reword `model.go`'s doc comment as the source of truth and align the three downstream restatements to it in the same change. Comment/prose only — no executable line touched.
- [x] **[APPLY]:** Five files edited, all comment/prose. `model.go` doc comment (malformed-vs-literal-equality interaction + illustrative miss-class list incl. directory entry); `hooks/memory-session-start.sh` legacy-suppression comment (to-EOF, self-clears at next UTC day); `actions/board.md`, `actions/work-reference.md`, `tools/queue-kanban/prime-do-kanban.md` aligned to `model.go`'s wording.
- [x] **[UNIFY]:** `git diff --stat` = exactly the 5 `write_set` files, 18 insertions / 12 deletions, no others (`git status --porcelain -uall` confirms; no build outputs, no debug artifacts). **No behavior change proved mechanically**: `git diff -U0` on the `.sh` and `.go` filtered for lines not starting with `#`/`//` returns empty — every changed line in both is a comment, so the awk block and `writeSetPatternsIntersect`'s body are byte-identical. Linters, in `tools/queue-kanban/`: `go build ./...` clean; `go test ./...` → `ok github.com/knews2019/skill-do-work/queue-kanban 1.285s`; `go vet ./...` clean; `gofmt -l .` empty. `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` Per-file review: `model.go` (doc comment only, gofmt-clean at existing comment width), `memory-session-start.sh` (comment block only, `set -euo pipefail` and awk untouched), `board.md`/`work-reference.md`/`prime-do-kanban.md` (one sentence/clause each, no structural or schema change).

## Requirements

1. **`hooks/memory-session-start.sh` (~:60-64):** the legacy-suppression comment claims "it self-clears as soon as the next capture is written in the new format" — false. Traced through the awk: once `capture_format` is `legacy`, no later heading (including a new quoted capture's) can be a boundary, so suppression runs to end-of-file for that day; it actually self-clears at the next UTC day's fresh log. Reword the sentence to match the behavior (which is correct and deliberate — only the comment is wrong).
2. **Malformed-glob phrasing at the four board-dialect sites** (`actions/board.md` ~:92, `actions/work-reference.md` ~:100, `tools/queue-kanban/prime-do-kanban.md` ~:51, and the `CHANGELOG.md` 0.146.3 bullet is history — leave it): "a malformed pattern matches nothing on the board" overstates. `writeSetPatternsIntersect` tests literal equality *before* `path.Match`, so two REQs declaring the *identical* malformed pattern do badge each other. The `model.go` doc comment is already precise ("treated here as no-match for that direction"); align the downstream restatements with it.
3. **Directory-entry miss-class** (`tools/queue-kanban/model.go` ~:1139-1142 doc comment and the same doc sites): `path.Match("actions/", "actions/board.md")` and `path.Match("actions", "actions/board.md")` are both false, so a `write_set` entry naming a directory never badges against a file inside it. The generic "the badge can miss real contention" hedge covers it, but the *named* miss-class list omits it — add it as an illustrative example (keeping the list marked illustrative per Closed Enumerations Go Stale). The dispatch gate is unaffected (absent/unexpandable ⇒ serialize).

## Constraints

- Prose/comments only — no changes to the awk, the Go matcher, or any behavior. If the builder is tempted to "fix" the behaviors instead, that is out of scope: the suppression-to-EOF and the conservative badge behavior are both deliberate.
- Item 2/3 doc sites ship to consumers; keep wording consistent across all of them in the same commit (the lock-step rule for `tools/queue-kanban` doc claims).

## Red-Green Proof
**RED prompt/case:** Read the session-start comment against the awk state machine (legacy → no boundary ever); read "matches nothing on the board" against `writeSetPatternsIntersect`'s literal-equality short-circuit with two identical malformed patterns; check the miss-class list for the directory-entry case.
**Why RED now:** Three doc claims describe behavior the code does not have (or omit behavior it does).
**GREEN when:** Each doc statement matches the traced behavior; the miss-class list includes the directory-entry example and stays marked illustrative.
**Validation:** User confirmed (approved capture of the reviewed finding set)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** B (Explore then Build)
**Reasoning:** Three independent doc fixes, but item 2/3 span multiple lock-step doc sites (`model.go` doc comment as the precise source + three downstream restatements in board.md / work-reference.md / prime-do-kanban.md) that must be verified against the actual code behavior and kept consistent — an exploration/trace pass before editing, not a blind edit. Route B.
**Complexity indicators:** 3 requirements across 5 files; a lock-step consistency requirement (align 3 restatements to model.go's precise wording); a code-behavior trace (awk state machine; `writeSetPatternsIntersect` literal-equality short-circuit; `path.Match` directory-entry falsity); a hard "prose/comments only — no behavior change" constraint; one `.go` file touched (doc comment only → `go test`/`vet`/`gofmt` must stay green).
**Rigor:** Standard independent review (main-context) + trace each reworded claim against the actual code, confirm the lock-step sites now agree, and confirm `go test`/`go vet`/`gofmt` still pass (comment-only Go change).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Exploration

### Behavior trace 1 — legacy suppression runs to EOF (`hooks/memory-session-start.sh:67-84`)

The awk carries two state variables, `in_capture_section` and `capture_format`. A heading line (`^## HH:MM UTC `) becomes a boundary only under two conditions (`:72-75`): `!in_capture_section`, or `capture_format == "quoted" && !is_quoted`. When `capture_format` is `legacy`, **neither** can hold — so no later heading is ever a boundary, `in_capture_section` never clears, and the `if (!in_capture_section) print` at `:82` suppresses every remaining line. Confirmed empirically against a synthetic legacy log (curated note A, then a legacy `session capture`, then curated note B): only note A printed; note B was suppressed. The scope is *that day's file* — `TODAY_LOG` is `logs/$(date -u +%F).md` (`:37`), so the next UTC day opens a fresh log with no legacy section, which is where the self-clear actually happens. The old comment's "self-clears as soon as the next capture is written in the new format" is exactly the case the state machine cannot reach.

### Behavior trace 2 — malformed patterns and the literal-equality short-circuit (`tools/queue-kanban/model.go:1152-1163`)

`writeSetPatternsIntersect` returns `true` on `leftPattern == rightPattern` at `:1153`, **before** either `path.Match` call. Probed with Go directly: `path.Match("src/[unclosed", "src/[unclosed")` → `false, syntax error in pattern` — so `path.Match` alone would *not* match a malformed pattern even against a byte-identical string; the literal-equality short-circuit is the sole reason two REQs declaring the same malformed entry badge each other. "A malformed pattern matches nothing on the board" is therefore an overstatement at the three downstream sites; `model.go`'s own "treated here as no-match for that direction" was already correct and is the wording the others were aligned to.

### Behavior trace 3 — directory entries never badge against files inside them

Probed: `path.Match("actions/", "actions/board.md")` → `false, <nil>`; `path.Match("actions", "actions/board.md")` → `false, <nil>`; positive control `path.Match("actions/*", "actions/board.md")` → `true`. Both directory forms also fail literal equality, so a `write_set` entry naming a directory renders no badge against any file inside it. The generic "the badge can miss real contention" hedge covered this, but the *named* miss-class list did not. The dispatch gate is unaffected: `actions/work.md` Step 1 reads absent/unexpandable as overlapping ⇒ serialize.

### Doc sites found (repo-wide grep for `malformed` and `self-clear`, excluding `do-work/`)

| Site | Status |
| --- | --- |
| `tools/queue-kanban/model.go:1139-1142` (doc comment) | edited — source of truth for items 2 & 3 |
| `hooks/memory-session-start.sh:59-63` (legacy-format comment) | edited — item 1 |
| `actions/board.md:92` (badge paragraph) | edited — items 2 & 3 |
| `actions/work-reference.md:100` (`write_set` schema line) | edited — items 2 & 3 |
| `tools/queue-kanban/prime-do-kanban.md:51` (REQ-032/034 lesson) | edited — items 2 & 3 |
| `CHANGELOG.md:65` (0.146.3 bullet) | **not touched** — history, and out of `write_set` per the REQ |
| `CHANGELOG.md:178` (0.139.4-era bullet, carries the same "self-clears with the next capture" claim) | **not touched** — same history rule; see Decisions D-02 |
| `tools/queue-kanban/model_test.go:596,605` (test comment + failure message: "must match nothing") | **not touched** — test file, out of `write_set`; see Decisions D-03 |

No other files in the repo restate either claim.

## Scope

**Files I will touch:** (all five `write_set` entries needed an edit; none dropped)

- `hooks/memory-session-start.sh`
- `tools/queue-kanban/model.go`
- `actions/board.md`
- `actions/work-reference.md`
- `tools/queue-kanban/prime-do-kanban.md`

## Implementation Summary

**Files changed:**

- `hooks/memory-session-start.sh` (modified) — reworded the `legacy` branch of the format-decides comment block. It now states the mechanism (once `capture_format` is `legacy`, no later heading can be a boundary, not even a new quoted capture's), the extent (suppression runs to the end of that day's log), and the real self-clear point (the next UTC day's fresh log). The awk block, `set -euo pipefail`, and every other executable line are byte-identical.
- `tools/queue-kanban/model.go` (modified) — doc comment on `writeSetPatternsIntersect` only. Item 2: the malformed-pattern clause now notes that the literal-equality check short-circuits before either `path.Match` call, so two entries carrying the same malformed text still intersect. Item 3: the closing display-only paragraph's parenthetical `(or a **/malformed pattern)` was expanded into an explicitly **ILLUSTRATIVE, not closed** miss-class list — glob-vs-glob, a `**` pattern, a malformed pattern against anything but its own twin, and an entry naming a directory, with the `actions/` vs `actions/board.md` falsity spelled out. The gate-unaffected sentence is retained. Function body unchanged.
- `actions/board.md` (modified) — one sentence in the `overlaps` badge paragraph (`:92`), aligned to `model.go`: malformed ⇒ no-match *for that direction*, literal equality first so identical malformed patterns still badge, illustrative-not-closed miss-class list including the directory entry, gate still serializes.
- `actions/work-reference.md` (modified) — the glob clause on the `write_set` schema line (`:100`), same three claims in the schema line's terser register.
- `tools/queue-kanban/prime-do-kanban.md` (modified) — the glob-dialect sentence in the REQ-032/034 lesson (`:51`), same three claims, phrased in the lesson's `model.go`-facing vocabulary ("short-circuits before either `path.Match` call").

**What was done:** Traced all three claims against the running code (awk exercised on a synthetic legacy log; `path.Match` probed directly for the malformed and directory-entry cases) and then corrected only the words. `model.go`'s comment was treated as the source of truth for items 2 and 3, with the three downstream restatements brought into lock-step in the same change. Zero behavior change: the suppression-to-EOF, the conservative badge, and the literal-equality-before-`path.Match` ordering are all deliberate and untouched — mechanically verified by filtering the `.sh`/`.go` diff for non-comment lines (empty).

## Decisions

- **D-01 — added the literal-equality note to `model.go` too, not only downstream.** `model.go`'s opening sentence already listed "identical text" as one of three intersecting cases, so it was not *wrong*; but a reader arriving at the malformed clause five lines later has no cue that the two interact. Since item 3 required editing this comment anyway, one clause was added there so the source of truth states the interaction the three downstream sites now state. Keeps the lock-step literal rather than inferential. (DECIDE & STATE: reversible, comment-only.)
- **D-02 — `CHANGELOG.md:178` carries the same false "self-clears with the next capture" claim and was left alone.** The REQ excludes the 0.146.3 bullet as history and `CHANGELOG.md` is not in `write_set`; the same rule applies to the older 0.139.4-era bullet. Flagged here rather than as a Discovered Task because the REQ's Constraints already settle the policy (changelog entries are a record of what shipped when, not live documentation).
- **D-03 — `tools/queue-kanban/model_test.go` still says "a malformed pattern must match nothing".** It is out of `write_set`, and the assertion itself is *correct* for what it tests (two **different** patterns, one malformed). Only the message's generality is loose. Left untouched to honor the declared write set; noted under Discovered Tasks.

## Discovered Tasks

- `tools/queue-kanban/model_test.go:596,605` — the malformed-glob test's comment and `t.Fatalf` message say "a malformed pattern must match nothing", the same overstatement REQ-050 corrected in the prose sites. The assertion is right; only the wording is loose, and a companion case pinning that two *identical* malformed patterns DO intersect would pin the literal-equality short-circuit that REQ-050's new wording now documents (currently nothing in the suite covers it). Out of this REQ's `write_set`.
- `CHANGELOG.md:178` (0.139.4-era entry) restates the "self-clears with the next capture" claim corrected here. History, so likely won't-do — recorded only so a future grep-based audit finds the decision instead of re-flagging it.

## Review

**Acceptance: Pass — overall ~96%.** Main-context review: traced each reworded claim against the actual code, confirmed lock-step consistency, verified the Go toolchain stays green.

**Requirements (all 3 met, all comment/prose):**
1. Hook legacy-suppression comment now says suppression runs to EOF for that day's log and self-clears at the next UTC day's fresh log — matches the awk (once `capture_format=legacy`, no later heading is a boundary). The false "self-clears as soon as the next capture" claim is gone.
2. Malformed-glob phrasing aligned to `model.go`'s precise wording across all three downstream sites (board.md, work-reference.md, prime-do-kanban.md): no-match for that direction, BUT literal equality short-circuits first, so two identical malformed patterns still badge. Verified against `writeSetPatternsIntersect`'s `if leftPattern == rightPattern` short-circuit. CHANGELOG 0.146.3 history left alone.
3. Directory-entry miss-class added as an ILLUSTRATIVE example (`path.Match` false for both `actions/` and `actions` vs `actions/board.md`) to model.go's doc comment + the 3 doc sites; list kept marked illustrative.

**Verified:** comment/prose-only (filtered git diff shows no non-comment line changed in the awk or the Go matcher); `go build`/`test`/`vet`/`gofmt` all clean; contract-regressions PASS; lock-step wording consistent across model.go (source) + 3 restatements.

No Important/Critical findings. No follow-ups. This fix is itself a restatement-sweep (aligning downstream restatements to the model.go source) — the discipline REQ-049 just institutionalized.

## Lessons Learned
**Worth knowing:** `writeSetPatternsIntersect`'s literal-equality short-circuit (`if left == right return true`) means "malformed patterns match nothing" was always an overstatement — identical malformed text badges. And a directory entry (`actions/`) never badges a file inside it (path.Match has no directory semantics). Both are deliberate conservative-badge behavior; only the docs were wrong. The `model.go` doc comment is the source of truth for the glob dialect — three doc files restate it and must move with it.

## Orientation
Three doc-accuracy fixes, comment/prose only: the memory hook's legacy-suppression comment (now: to-EOF, clears next UTC day), and the board's malformed-glob + directory-entry miss-class descriptions aligned lock-step across `model.go` (source), `board.md`, `work-reference.md`, `prime-do-kanban.md`. No behavior change; Go toolchain green. No map change.
