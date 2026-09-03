---
id: REQ-530
title: 'Order ready work by the newest REQ it unblocks'
status: cancelled
created_at: 2026-09-03T10:59:01Z
user_request: UR-101
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work/tools/do-work-cli/internal/nextselection/
  - skills/do-work/tools/do-work-cli/internal/dependencygraph/
  - skills/do-work/actions/work-reference.md
  - skills/do-work-board/tools/queue-kanban/
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
completed_at: 2026-09-03T20:40:53Z
---

# Order Ready Work by the Newest REQ It Unblocks

## What

Default and UR-expanded selection in `do-work run` picks the oldest ready REQ (lowest number) first. It should pick the newest first without breaking dependencies: ready work is ordered by the newest REQ it unblocks, so an old prerequisite of the newest capture runs before newer-but-independent REQs.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"The principle is that the one that I just captured I'm more interested in it."

## Context

- The selector is `do-work-cli next` (`skills/do-work/tools/do-work-cli/internal/nextselection/`). `next_selection.go:43` stable-sorts eligible records by the three priority classes only; inside a class the order is the repository snapshot's ascending REQ-number order (`internal/repositorymodel/repository_model.go:298-311`, `requestIDLess`), and `next_targets.go` `queueCandidates` walks that snapshot in order. So the oldest ready REQ wins today.
- `skills/do-work/actions/work-reference.md:295` states the current contract: "Existing queue order is preserved inside each class." That sentence changes with this REQ.
- Decisions taken at capture:
  - D1 "Newest" means the highest REQ number. Ids are allocated monotonically at capture, so id order is capture order; `created_at` is not consulted.
  - D2 The three priority classes stay (`repository_gate_repair` first, ready `gate_deferred` parents second, ordinary work third). Newest-first applies inside each class.
  - D3 Dependency readiness is unchanged: a REQ with an unmet `depends_on` is still excluded. Among ready REQs the sort key is the highest REQ number in the REQ's transitive dependent closure (itself included), descending; ties break on the REQ's own number, descending. User chose this over plain newest-ready-first at capture.
  - D4 Explicit REQ tokens keep caller order, as today. UR expansion, `--wave`, `--fan-out`, `--skip-impact-negligible`, and `run-simple-reqs` all flow through the same selector and inherit the new order.
  - D5 The Kanban board's pending order follows this REQ's order, read from the core CLI (Addendum 2026-09-03, D6–D10). *Resolved conflict: board display order out of scope → in scope, consuming the Go tool's order.*

## Detailed Requirements

- Default and UR-expanded selection order ready work newest-first with prerequisite promotion (D3), inside each priority class (D2).
- Ordering happens before the fan-out bound is applied, so `FAN-OUT-LIMIT` exclusions are the oldest-unblocking work, never the newest.
- Rewrite the one sentence at `work-reference.md:295` to state the new rule. `work.md` keeps "process selected REQs in the returned stable order" and gains no procedural prose.
- Lock-in tests in `internal/nextselection`, each naming its failure: (a) two ready ordinary REQs, the newer is selected; (b) the newest REQ depends on an old pending REQ and a newer independent REQ is also ready, the old prerequisite is selected first; (c) an older repair still beats a newer ordinary REQ; (d) explicit token order is unchanged.

## Constraints

- Mechanics stay in Go. No new sentence pins in `_dev/tests/contract-regressions.sh` (UR-100 keeps that file at its current line count).
- Overlap declared, not a dependency: REQ-505 (Move selection and claim behind advance) also writes `work-reference.md`; REQ-490 (Compute wave depth from satisfied duplicate records) also writes `internal/nextselection`.

## Red-Green Proof

**RED prompt/case:** A queue with REQ-811 (pending, no dependencies), REQ-815 (pending, no dependencies), and REQ-816 (pending, `depends_on: [REQ-815]`). Run `do-work-cli next` with no targeting tokens.
**Why RED now:** REQ-811 is selected because the snapshot is walked in ascending number order and nothing reorders inside the ordinary class.
**GREEN when:** REQ-815 is selected (it unblocks REQ-816, the newest), REQ-811 is a `FAN-OUT-LIMIT` exclusion; after REQ-815 completes, REQ-816 is selected before REQ-811. A second case with REQ-816 depending on REQ-811 instead selects REQ-811 over the newer REQ-815, which is what distinguishes prerequisite promotion from plain newest-ready-first.
**Validation:** User confirmed. The prerequisite-forward rule was chosen from two options at capture; the RED/GREEN case was confirmed at verify-requests on 2026-09-03.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2927 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "frozen target projection and fan-out bounds" (`projection-before-bounding`).
- `_dev/primes/lessons-action-files.md` — 3887 tokens, over the 2000-token budget; `slugged: partial`. Matched on "pipeline fields, status contracts, downstream readers" through the `work-reference.md` sentence.

## Addendum (2026-09-03)

User added:

> ````text
> update ur-101 so that the kanban report also is in sync with the choosen task ordering
> BTW: is there a go tool that is doing the picking of the next task? I imagine there should be one to do that and the LLM should use that.
> ````

and then:

> ````text
> also the kanban board should use that so that the go tool is the point of truth in the ordering
> ````

The picker already exists: `do-work-cli next` (`internal/nextselection/`), and `actions/work.md` consumes its returned order. This addendum extends the REQ so the Kanban board shows the same order, read from the core CLI rather than re-derived. It reverses D5 (recorded above). `next` itself is not the board's source: it runs every REQ's `blocked_check` probe, evaluates claims and earmarks, and applies the fan-out bound, none of which a board rebuild may trigger.

### Decisions

- D6 One comparator in `internal/nextselection` (e.g. `orderReadyWork`) implements D2+D3. `Select` and the new command both call it; the board never reimplements the rule.
- D7 New read-only subcommand `do-work-cli queue-order`, registered beside `next` in `Handlers()`. It discovers the repository, builds the dependency graph, ranks every REQ under `do-work/queue/` with the shared comparator, and emits JSON `queue_order: [{request_id, rank, selection_priority, unblocks_newest, ready, unmet_dependencies}]`; text format is one line per REQ. No probes, no claim or earmark evaluation, no fan-out, no writes. Ready and waiting REQs are ranked by the same key so both board groups can follow it.
- D8 The board (`skills/do-work-board/tools/queue-kanban`) orders `PendingReady` and `PendingWaiting` by `queue-order` rank. `buildBoard` gains an injected `queueOrderLookup` with the same shape as `gitCommitDateLookup` (`model.go:418`); the live path runs the core launcher `do-work-cli.sh --repo-root <root> --format json queue-order`. Launcher discovery: a `--do-work-cli PATH` flag when given, else the sibling core skill relative to the board tool's own directory (`../../../do-work/tools/do-work-cli.sh`, valid in-repo and installed per `suite/modules.tsv`). Ids the CLI did not return trail in numeric order.
- D9 The fallback is loud, never silent: launcher missing, non-zero exit, or unparseable output keeps numeric order and appends a `board.Warnings` entry naming the reason and the exact command. The lookup runs only on a `refreshBoardData` fingerprint miss (`serve.go:334`) and in `generate`, never per HTTP request.
- D10 This is a read, not a write surface: the root `CLAUDE.md` three-surface count is unchanged.

### Added Requirements

- `queue-order` exists per D7 and `next` orders through the same comparator (D6).
- The board's READY and WAITING groups render in `queue-order` rank (D8); the client keeps rendering received order (`web/board-cards.js`), so no client sort is added.
- Fallback per D9, with the warning visible on the board.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` § Traps gains one line: pending order comes from core `queue-order`; numeric order only as a warned fallback. `_dev/primes/prime-kanban-board.md`'s "deterministic scalar commands do not depend on [the board]" line stays true and is not edited.
- Additional lock-in tests, each naming its failure: (e) `queue-order` and `next --fan-out N` agree on ready order over one fixture (pins the shared comparator); (f) `queue-order` never invokes the probe runner when a `blocked_check` REQ is present; (g) board `PendingReady`/`PendingWaiting` follow the injected order, not numeric; (h) lookup failure gives numeric order plus a warning. Existing pins `generate_test.go:1043-1072` and `model_test.go:201` change only if their fixture order changes under the new key.

### Red-Green Proof, case 2 (board)

**RED prompt/case:** Same queue as case 1: REQ-811 (pending), REQ-815 (pending), REQ-816 (pending, `depends_on: [REQ-815]`). Open the board.
**Why RED now:** READY shows REQ-811, REQ-815 (numeric), and no board code consults the CLI.
**GREEN when:** READY shows REQ-815, REQ-811 and WAITING shows REQ-816. With the launcher removed, READY is numeric and the board carries a warning naming `queue-order`.
**Validation:** User confirmed (this addendum, 2026-09-03).

### Added Constraints

- Overlap declared, not a dependency: REQ-482, REQ-485, REQ-486, REQ-519, REQ-520, and REQ-522 also write `skills/do-work-board/tools/queue-kanban/`.
- The board must keep working with no Go toolchain and no core launcher present: D9's fallback is the contract, not an error exit.

---
*Source: "when picking up the next REQuest it should be the latest LIFO not the oldest REQuest. Of course don't break the dependencies. The principle is that the one that I just captured I'm more interested in it."*

## Cancelled

- **When:** 2026-09-03T20:40:53Z
- **Why:** superseded by REQ-561 (UR-107): a three-value priority field the maintainer sets replaces inferring order from REQ numbers. Maintainer's 2026-09-03 triage.
- **Decided by:** user, via `do-work abandon`
