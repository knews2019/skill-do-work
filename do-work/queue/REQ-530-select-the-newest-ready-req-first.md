---
id: REQ-530
title: 'Select the newest ready REQ first'
status: pending
created_at: 2026-09-03T10:59:01Z
user_request: UR-101
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
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
---

# Select the Newest Ready REQ First

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
  - D5 The Kanban board's display order is out of scope.

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
**Validation:** User confirmed the prerequisite-forward rule at capture (chosen from two options); the RED case is inferred during capture.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2927 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "frozen target projection and fan-out bounds" (`projection-before-bounding`).
- `_dev/primes/lessons-action-files.md` — 3887 tokens, over the 2000-token budget; `slugged: partial`. Matched on "pipeline fields, status contracts, downstream readers" through the `work-reference.md` sentence.

---
*Source: "when picking up the next REQuest it should be the latest LIFO not the oldest REQuest. Of course don't break the dependencies. The principle is that the one that I just captured I'm more interested in it."*
