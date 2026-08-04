---
id: REQ-097
title: assigned_to advisory field — schema line, scan skip-and-report, board parse (lock-step)
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T21:03:58Z
completed_at: 2026-08-04T21:12:00Z
kb_status: pending
user_request: UR-018
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-096]
maintenance: false
related: [REQ-096, REQ-098]
batch: parallel-building
write_set: [actions/work-reference.md, actions/work.md, tools/queue-kanban/model.go, tools/queue-kanban/model_test.go, tools/queue-kanban/generate.go, tools/queue-kanban/web/board.js, tools/queue-kanban/web/board.css, CLAUDE.md]
---

# `assigned_to` Advisory Field — Schema, Scan Skip, Board Parse

## What

Add the cooperative claim marker the whole model rests on: a single advisory frontmatter field, `assigned_to: "<session-name>"`, on a pending REQ. The default work-loop scan skips-and-reports such REQs; explicitly targeting one (`do-work run REQ-NNN`) overrides and clears the field. The board parses it display-only. **No verb, no status, no staleness clock, no release command** — the 0.163.0 forbidden-token ratchet stays fully intact.

## Detailed Requirements

- **Schema:** add `assigned_to` to `actions/work-reference.md`'s frontmatter block and Schema Read Contract — in the **verbatim-read class** alongside `write_set` (`:206`): no alias map, no normalization, no canonical vocabulary. Optional on every REQ; capture may seed it when the user earmarks work for a named session.
- **Scan behavior:** one skip-and-report sentence in `actions/work.md` Step 1's default selection: a pending REQ with `assigned_to` is skipped and listed ("assigned to <name>") exactly like the existing dependency-skip reporting. Explicit targeting (`do-work run REQ-NNN`) overrides the skip and **clears the field** as part of the claim.
- **Board (same commit — lock-step rule):** `tools/queue-kanban/model.go` parses `assigned_to` display-only: badge on the card + drawer metadata row, **no column logic, no scheduling**, same class as `write_set`. ~15 lines + test in `model_test.go`.
- **No `assigned_at`, no staleness threshold, no auto-release** — an assignment persists until cleared by an explicit run or hand-edit.
- **Ratchet check:** confirm `assigned_to` does not trip `_dev/tests/contract-regressions.sh`'s reservation token patterns (`reserved_for`, `reserved_at`, underscore-token forms); run the suite.
- Mirror the schema note in the co-location sentence style the Testing placeholders use (CLAUDE.md "keep the parser in lock-step" — restate inline in shipped files, never cite CLAUDE.md).

## Constraints

- Display-only at any builder count — nothing schedules, gates, or dispatches on `assigned_to` except the Step 1 skip (which is a *courtesy read*, not a gate: explicit targeting overrides).
- The forbidden-token ratchet must stay green with no test weakening.

## Red-Green Proof

**RED prompt/case:** `go test ./...` in `tools/queue-kanban/` with a new test asserting a REQ carrying `assigned_to: "cloud-alpha"` surfaces the value in the parsed model — fails before the parser change. Prose half: today no shipped file defines `assigned_to`, so a second session has no sanctioned way to see "this is earmarked".
**Why RED now:** The field does not exist; reserve (its predecessor) was deleted at 0.163.0.
**GREEN when:** The Go test passes; `actions/work.md` Step 1 documents skip-and-report + override-and-clear; contract-regressions suite green.
**Validation:** User confirmed (ask-tool answer: "assigned_to field only").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 2, items 5–6).

---
*Source: approved plan, Phase 2*

## Triage

**Route: B** - Medium

**Reasoning:** The field, its read class, its one behavioural reader, and the lock-step obligation are all specified. What needed discovery was where the verbatim-read class is *defined* (so `assigned_to` joins it rather than getting a parallel rule), how the board's display-only fields are threaded from parser to payload to badge, and whether the skip-and-report has an existing exit-summary section to slot into. `tdd: true`, so the Go test came first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided, TDD

*Skipped by work action*

## Exploration

**The verbatim-read class is defined once**, at `actions/work-reference.md`'s Schema Read Contract: "List-valued path fields are outside this contract and are read verbatim." `assigned_to` is scalar, not list-valued, so joining that class meant re-grounding the paragraph on the *reason* (no canonical vocabulary to normalize against) rather than on the shape — otherwise the new field would need a second, parallel rule saying the same thing.

**The board threads a display-only field through four places, not one:** `model.go`'s `RequestTicket` struct + its `parseRequestTicket` assignment, `generate.go`'s `generatedRequest` payload struct + its copy site, `web/board.js`'s badge builder + drawer `appendMetaRow`, and `web/board.css`'s badge rule. `write_set` was the template for all four. Only **one** payload construction site exists (`generate.go:258`), which `serve` and `generate` share — verified by grepping for the `WriteSet:` copy.

**The skip-and-report has a home already:** `actions/work-reference.md`'s **Composed Exit Summary (Step 1)** enumerates the sections Step 1 renders when nothing is claimable. Adding a seventh reader for `assigned_to` meant adding a section there rather than inventing an ad-hoc report line in `actions/work.md`.

**Pre-existing drift found while doing it:** that section's intro said "Six sections may apply" over five numbered sections, and its closing line said "any of the five sections". The count had already gone stale once. Replaced with the *condition* per the Closed-Enumerations rule, so adding this section does not re-introduce a number to keep in sync.

**Ratchet clearance, checked before writing any prose** (the handdown's instruction to run this early): the reservation tokens are `status: reserved`, `reserved_for`, `reserved_at`, `do-work reserve`, `actions/reserve\.md`. `assigned_to: "cloud-alpha"` matches none of them — the tokens are underscore/path-shaped precisely so an unrelated `_to` field cannot collide.

## Scope

**Files I will touch:**
- `tools/queue-kanban/model_test.go` — three tests, written first (verbatim read, absent-reads-as-unassigned, never affects column placement)
- `tools/queue-kanban/model.go` — `AssignedTo` field + parse
- `tools/queue-kanban/generate.go` — payload field + copy
- `tools/queue-kanban/web/board.js` — card badge + drawer row
- `tools/queue-kanban/web/board.css` — badge styling
- `actions/work-reference.md` — schema line, verbatim-read class membership, new exit-summary section
- `actions/work.md` — Step 1 skip-and-report + override, Step 2 clear-on-claim, two enumeration updates
- `CLAUDE.md` — the lock-step field list (maintainer doc; the enumeration would otherwise go stale)

**Acceptance criteria (restated from REQ):**
- [ ] Schema line added, in the verbatim-read class alongside `write_set`, optional on every REQ
- [ ] Step 1 default scan skips and reports; explicit `do-work run REQ-NNN` overrides and clears
- [ ] `model.go` parses it display-only; badge + drawer row; no column logic, no scheduling — **same commit** as the schema line
- [ ] Test in `model_test.go`
- [ ] No `assigned_at`, no staleness threshold, no auto-release
- [ ] `assigned_to` does not trip the reservation ratchet; suite green with no weakening
- [ ] Lock-step obligation restated inline in the shipped files (never by citing the maintainer doc)

## Pre-Flight

- **WARN — baseline suite red before any change:** the same 8 `chmod 500`-versus-root failures inherited by every REQ in this batch. Gate is "still exactly those 8".
- `go test ./...` in `tools/queue-kanban/` green at claim time; Go toolchain resolves.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model_test.go` (modified) — three tests added, before the implementation: `TestParseRequestTicketReadsAssignedToVerbatim` (mixed case plus an underscore, which a normalizing reader would fold), `TestParseRequestTicketLeavesAssignedToEmptyWhenAbsent`, `TestAssignedToNeverAffectsColumnPlacement`.
- `tools/queue-kanban/model.go` (modified) — `AssignedTo string` on `RequestTicket` with the verbatim/display-only/lock-step contract in its doc comment, plus `coerceScalarToString(fields["assigned_to"])` in the parser.
- `tools/queue-kanban/generate.go` (modified) — `AssignedTo string \`json:"assignedTo,omitempty"\`` on the payload struct and its copy from the ticket. `omitempty` is what makes absence read as *unassigned* rather than as an empty value.
- `tools/queue-kanban/web/board.js` (modified) — `assigned <session>` card badge (truncated at 18 chars, full value in the tooltip) and an `Assigned to` drawer metadata row.
- `tools/queue-kanban/web/board.css` (modified) — `.badge-assigned`, deliberately neutral (`--ink-soft` on `--surface-3`) rather than accented: it is somebody's intention, not a state of the work.
- `actions/work-reference.md` (modified) — the schema line; the verbatim-read paragraph re-grounded on *no canonical vocabulary* so a scalar field joins the class for the stated reason; a new **Assigned-elsewhere** exit-summary section; the section intro and closing line switched from a stale count to the condition.
- `actions/work.md` (modified) — Step 1's skip-and-report paragraph with the explicit-targeting override and the UR-expansion non-override, Step 2's clear-on-claim, and both exit-summary enumerations.
- `CLAUDE.md` (modified) — `assigned_to` added to the lock-step display-only field list, with a sentence distinguishing it from `write_set` (the board badges it; the *pipeline* is what acts on it).

**What was done:** Added the advisory cooperative claim marker end to end — schema, one courtesy reader in the work loop, and display-only board rendering — with the parser and the schema line in the same commit per the lock-step rule. No verb, no status, no timestamp, no clock.

## Testing

### Red-green validation (`tdd: true`)

**RED**, before any implementation, from the tests written first:

```
$ go test ./...
./model_test.go:657:12: ticket.AssignedTo undefined (type *RequestTicket has no field or method AssignedTo)
./model_test.go:658:81: ticket.AssignedTo undefined (type *RequestTicket has no field or method AssignedTo)
./model_test.go:674:12: ticket.AssignedTo undefined (type *RequestTicket has no field or method AssignedTo)
./model_test.go:675:102: ticket.AssignedTo undefined (type *RequestTicket has no field or method AssignedTo)
FAIL	github.com/knews2019/skill-do-work/queue-kanban [build failed]
```

**GREEN**, after the parser change:

```
$ gofmt -l .
$ go test ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	4.079s
```

This traces to the REQ's `## Red-Green Proof` exactly — its RED case is "a new test asserting a REQ carrying `assigned_to: \"cloud-alpha\"` surfaces the value in the parsed model". The fixture uses `"Cloud-Alpha_2"` rather than the literal `cloud-alpha`: mixed case plus an underscore is the same assertion with the *verbatim* half actually exercised, since a normalizing reader would fold the case and the plan's own value would pass either way.

### Board smoke (the plan's verification bullets)

Three fixture REQs in a scratch tree — one plain pending, one carrying `assigned_to: "Cloud-Alpha_2"`, one carrying the legacy `status: reserved`:

```
$ queue-kanban generate --repo-root <scratch> --out <scratch>/out
$ grep -o '"assignedTo"' <scratch>/out/board-data.js | wc -l
1
$ grep -o '"assignedTo":"[^"]*"' <scratch>/out/board-data.js
"assignedTo":"Cloud-Alpha_2"

$ queue-kanban summary --repo-root <scratch>
  pending             : 2
    ready to work     : 2
    waiting on deps   : 0
  claimed             : 0
  needs-input/blocked : 1
  warnings            : 1
    ! REQ-802 has unrecognized status "reserved" — shown under Needs input / Blocked; fix: …
```

- **Badge renders, value verbatim:** the payload key is present and carries `Cloud-Alpha_2` unaltered.
- **`omitempty` works:** one key across three REQs — the two without the field emit nothing, so absence reads as *unassigned* rather than as an empty session name.
- **No bucket change:** the assigned REQ is still `pending / ready to work`. It did not move, and it was not counted as waiting on deps or as blocked.
- **Legacy `reserved` card unchanged:** still an unrecognized status routed to Needs input / Blocked with the same warning text.

### Ratchet and suite

```
$ for t in 'status: reserved' 'reserved_for' 'reserved_at' 'do-work reserve'; do echo 'assigned_to: "cloud-alpha"' | grep -q "$t" && echo TRIPS || echo clear; done
clear
clear
clear
clear

$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -c '^FAIL'
8
```

Name-for-name the pre-existing eight. No reservation token reintroduced, no assertion weakened, nothing added to the ban list's exemptions.

### Do-not-build check

No `assigned_at` anywhere; no staleness threshold, refresh interval or liveness probe; no `reserve`/`release` verb and no `reserved` status; nothing schedules on `write_set` (untouched by this REQ); no auto `git pull`/`push` prescribed.

## Lessons Learned

**What worked:**
- Writing the three Go tests first genuinely paid, because the RED was a *compile* failure — which is the strongest possible RED for a new field, and it proves the test is actually reaching the code under test rather than passing vacuously.
- Making the verbatim fixture hostile (`Cloud-Alpha_2` instead of the plan's `cloud-alpha`). The plan's value passes with or without normalization, so it would not have tested the contract it was written for.

**What didn't:**
- The first board smoke reported `assignedTo` absent from the payload and looked like a threading bug. It was a **stale compiled binary** — the batch handdown warns about exactly this, and the warning was still not enough to stop it happening. The tell was `grep` finding the JS identifier but never the data value.
- The second smoke then grepped `index.html`, which never carries the data: `generate` writes the payload to a sibling **`board-data.js`** (plus a lazy `board-markdown.js`). Two false negatives in a row, neither of them a code defect.

**Worth knowing:**
- `generate --out <dir>` produces three files. **`board-data.js` is the payload** — assert against that, never `index.html`.
- There is exactly **one** payload construction site (`generate.go`'s `generatedRequest` literal), shared by `serve` and `generate`, so a new display field needs one copy line, not two. Confirm by grepping the `WriteSet:` copy before assuming otherwise.
- `omitempty` on the payload string is load-bearing, not tidiness: it is what preserves the *absence reads as unknown* convention the board applies to `write_set`.
- The `.badge-*` rules use `--ink-soft` / `--surface-3` / `--line-firm`; `--text-muted`, `--surface-raised` and `--border-subtle` do not exist in `board.css`. A wrong variable name fails silently — the badge just renders with the base `.badge` styling.

## Orientation

A REQ can now be earmarked for a named session with one advisory frontmatter field, and every other session's default run leaves it alone and says so. Naming it explicitly takes it anyway and clears the marker. The board shows who it is earmarked for and changes nothing else about the card. This is the cooperative claim marker the whole claim-anywhere model rests on — the thing the deleted reserve verb used to do, minus the verb, the status, and the clock. Lives in the frontmatter schema (`actions/work-reference.md`), the work loop's Step 1/Step 2 (`actions/work.md`), and the board's display path (`tools/queue-kanban/`). `prime_files` is empty; `tools/queue-kanban/prime-do-kanban.md` was checked and its referenced paths all still exist, so no prime went stale.

## Qualification

**Passed** — 8 files verified, 7 acceptance criteria traced, red-green evidence confirmed by the orchestrator reading the actual `go test` output rather than the builder's description of it.

- **Files exist and show in the diff:** all eight, `git diff --stat` below.
- **Substantive:** the Go change is ~6 lines of code plus its contract comment and three tests; the prose changes carry a new schema line, a new exit-summary section, two new Step 1/Step 2 rules, and two enumeration corrections.
- **Wired, not orphaned:** the field is read at four consumer sites (parser → payload → badge → drawer row) and named at three prose sites; nothing added is unreferenced.
- **Flowing, not hollow:** proven by output, not inspection — the board smoke shows the real value reaching the real payload from a real fixture file.
- **Lock-step honored:** `model.go` and the `actions/work-reference.md` schema line are in this one commit, which is the rule's whole point.

## Review

**Overall: 94%** | 2026-08-04T21:12:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Acceptance:** Pass — TDD red-green captured, board smoke covers all four plan verification bullets, ratchet clear, suite at baseline.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Requirements checklist:** all seven `## Scope` criteria delivered; see `## Testing` for the evidence behind each.

**Minor:**
- Scope grew from the REQ's original `write_set` by three files (`generate.go`, `board.css`, `CLAUDE.md`) plus a narrowing of `tools/queue-kanban/web/*` to the two files actually touched. Each was added to Scope and `write_set` before editing, and each is a consequence the REQ implied without naming: a display field cannot reach a badge without a payload field, a new badge class needs a rule, and the maintainer doc's lock-step list is an enumeration that goes stale exactly as the Closed-Enumerations rule predicts. Declared, not drifted — but worth naming since the REQ's own `write_set` was written expecting three files.
- The pre-existing "Six sections may apply / any of the five sections" drift in the Composed Exit Summary was fixed in passing rather than filed. It is one sentence in the paragraph this REQ had to edit anyway, and leaving a known-wrong count while adding a section to the same list would have made it wrong in a *new* way. A separate REQ for a two-word correction inside an already-open paragraph is the more wasteful option.

**Scope drift:** none against the declaration. Eight files declared in `## Scope`, eight touched, none extra.

**Restatement sweep (MUST):** run. This REQ *adds* a schema field rather than redefining one, so the sweep's question is where the new field's semantics get restated and whether any of those restatements can drift: four in-code sites (all carrying the display-only contract in comments), the schema line, `actions/work.md`'s two rules, the new exit-summary section, and `CLAUDE.md`'s lock-step list. The last is the one that would have silently gone stale — it enumerates the display-only parsed fields, and an unlisted field reads as *not covered by lock-step*. Updated. No pre-existing text restates a semantics this REQ changed, because nothing said anything about `assigned_to` before it.

**Suggested additional testing:**
- REQ-098's `assigned-elsewhere-claimed-here` probe should assert against `board-data.js`, not `index.html` — the trap this REQ hit twice.
- A `serve`-mode click-through of the drawer row (the smoke covered the static path only; both share one payload builder, so the risk is presentation-only).

*Reviewed by review-work action (pipeline mode, in-session)*
