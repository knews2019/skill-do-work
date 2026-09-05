# Builder brief — REQ-579

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-579-finding-rows`
- **Your branch (already checked out there):** `worktree-agent-REQ-579-finding-rows`
- **Route:** B
- **Base commit:** 5f4821ab

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ below. If you need one, stop and report it in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (the REQ is `tdd: true`)
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/anti-slop.md` (the strip is a human-facing artifact and the REQ says to load it before styling)

Also read `_dev/primes/prime-kanban-board.md` and, where your change touches what they name, the lesson satellites `_dev/primes/lessons-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the declared write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes, `node --check` for changed client files), verify no debug artifacts in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. Use:
- Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`
- Node lane: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...`

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-120117/REQ-579-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, the RED observation (test name + failure text) and the GREEN observation.
- `## Lesson evidence` — each lesson satellite you read and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN`, each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings, each stamped with one of exactly these impact tokens: `impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible`. Do not invent a token outside that set and do not fix the items inline.
- `## Integration seams` — any exact line that belongs in a file outside your write set, with where it goes. The orchestrator applies it.
- `## Exploration` — what you found before writing code: the current renderer's shape, the classes you are replacing, and anything the request's own findings got wrong. The orchestrator folds this into the REQ.

This REQ was triaged Route B. Its `## Exploration` was not produced by a separate agent, because the request's own capture-time findings already name the producer, the payload site, the renderer and the classes. **You produce that section**, under `## Exploration` in your hand-back, and you correct anything the capture got wrong rather than trusting it.

REQ-578 landed on main just before your base: `applyView` in `web/board-controls.js` now hides `#board-findings` on the Activity view by reading the renderer's own output — `#board-findings-cards` children and `#board-findings-skipped-list` children — to decide whether the strip has anything to say. **Requirement D5 depends on this.** If your new markup renames or removes either of those element ids, `applyView` stops seeing content and the strip will misbehave on every view. Either keep both ids as the hosts of your new rows, or hand back the exact `board-controls.js` replacement lines as an integration seam — that file is outside your write set, so you must not edit it yourself.

Also generate a board and look at the result before you hand back. The prime asks for it and this is a styling change: `go run . generate --out <a scratch directory> --repo-root <your worktree>` from the queue-kanban module, then read the produced page. Report what you saw.

---

# The request

---
id: REQ-579
title: 'Render verify findings and skipped probes as compact rows in one list'
status: claimed
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
---

# Render Verify Findings and Skipped Probes as Compact Rows in One List

## What

The Verify Findings strip (`#board-findings`, added by REQ-285) renders each finding as a bordered card in a multi-column grid and lists skipped probes as bullets inside a collapsed disclosure below. Replace both with one flat list of compact rows: one row per finding and one row per skipped probe, no cards, no columns, no separate disclosure. A finding and a skipped probe are the same kind of thing to the reader ("verify has something to tell you"), so they get one visual language.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — `Subject` on `VerifyFinding`, set by each probe
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — subject assertions per probe
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — `subject` in the board payload
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — one flat row list replacing the card grid and the disclosure
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — row classes replacing the `.board-finding*` and `.board-findings-skipped*` classes
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — strip markup
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modify) — the captured RED case

**Files I will NOT touch:** the completion-anomalies strip and its cards, `web/board-controls.js` (REQ-578 owns the Activity hide rule and it must keep working against the same `#board-findings` element), `CHANGELOG.md`, `VERSION`, everything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] One flat list of rows: no `.board-finding` card, no `details` disclosure, header carries the finding count and the skipped count when non-zero
- [ ] Two weights only, both from the producer: `fixable: true` renders muted with a "cleanup can fix" tag, a skipped probe renders muted, everything else normal
- [ ] `Subject` exists on `VerifyFinding` and `subject` in the payload; the client groups rows by exact string match and prints each non-empty subject once as a heading, with empty-subject rows after the grouped ones in producer order
- [ ] The strip still hides when there is nothing to report, and REQ-578's Activity-view hide still works against the same element and attribute

