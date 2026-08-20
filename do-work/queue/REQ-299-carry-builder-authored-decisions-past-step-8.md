---
id: REQ-299
title: "[impact-rule-change] Review fix: carry builder-authored sections past Step 8, starting with ## Decisions"
status: pending
created_at: 2026-08-19T20:03:19Z
status_changed_at: 2026-08-20T08:22:51Z
user_request: UR-055
addendum_to: REQ-270
domain: general
review_generated: true
sweep: true
sweep_key: builder-authored-sections-unread-outside-step-8
impact: impact-rule-change
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/review-work.md
- _dev/tests/contract-regressions.sh
---

# Review Fix: Carry Builder-Authored Sections Past Step 8, Starting with `## Decisions`

## What

REQ-270 closed the case where a worktree builder's `## Discovered Tasks` never reached
Step 8, and keyed its new rule on the condition — but scoped that condition's home to
**Step 8's substeps** (`actions/work.md` → *Where a builder-authored section is read from*,
which opens "Some substeps **below**"). Its independent review then found a second
builder-authored section with the identical defect, read from **outside** Step 8, where the
new rule structurally cannot reach:

`## Decisions` is written by the builder (`actions/work.md` Step 6, "Log Decisions as
D-XX" — unqualified, so a worktree builder is told to write a file it may not write), and
read at two sites: `actions/review-work.md` Step 4's traceability check ("If a
`## Decisions` section exists in the REQ: verify that significant implementation
choices … are documented") and the end-of-run Decision Brief's HANDLED block
(`actions/work-reference.md` → **Decision Brief (hand-back format)**).

**The failure is silent in both directions.** Review's traceability check finds no section
and reports clean rather than flagging an absence it cannot distinguish from "the builder
made no decisions"; the Decision Brief renders an empty HANDLED list, so a DECIDE & STATE
choice the builder actually made never reaches the user. Under fan-out that is every
builder's decisions, every run.

**Done means the class cannot recur:** the rule lives somewhere every reader of a
builder-authored section inherits it, not in Step 8's preamble, and the sections Step 6
tells the builder to author are the same set the hand-back contract names.

## Instances

- [ ] `## Decisions` — builder-authored at `actions/work.md` Step 6; read by
  `actions/review-work.md` Step 4 and by the Decision Brief's HANDLED block. The instance
  that fired this REQ.
- [ ] `actions/work.md` Step 6's `## Decisions` instruction itself — it must say where the
  section goes when the builder may not write the REQ, exactly as REQ-270 fixed the
  `## Discovered Tasks` bullet beside it.
- [ ] The rule's home — REQ-270 put it under Step 8 because that was the only reader then
  known. Move or restate it so a reader outside Step 8 inherits it, and mark the reader
  list illustrative rather than closed.

## Context

Found during the independent review of REQ-270 (finding 2, Important, gate:
`impact-rule-change` — it changes where a rule lives and which readers inherit it, across
several sites). REQ-270's other Important finding was a stale restatement in
`crew-members/general.md` that actively defeated its own fix; that one was repaired inside
REQ-270 rather than deferred, because the fix did not work end to end without it. This one
is genuinely a different scope: it needs a reader outside Step 8 and a check that did not
exist.

Created `pending-answers` per the generation-≥2 cascade depth stop: REQ-270 carries
`review_generated: true`.

## Requirements

- Every reader of a builder-authored REQ section — at Step 8 or anywhere else — reads the
  hand-back when the REQ file does not carry the section, and the rule states that as its
  condition rather than listing the readers it knows about today.
- `actions/work.md` Step 6's `## Decisions` instruction names the hand-back for a builder
  that may not write the main tree.
- The Decision Brief and review-work's traceability check distinguish "no section anywhere"
  from "the builder recorded nothing", and say which.
- **Every change holds in a consumer install**, where the suite is vendored under
  `.claude/skills/` and only `do-work/` is at the project root — see **Consumer-Install
  Constraint** below. Proven against a consumer-shaped fixture, not by reading.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** The property check REQ-270's reviewer specified, which pins the
property rather than any wording: extract every `## <Name>` section that `actions/work.md`
Step 6 instructs the **builder** to author, and assert each one is named in
`actions/work-reference.md`'s per-builder-output hand-back contents. It fails today, on
`## Decisions`.
**Why RED now:** `work-reference.md`'s hand-back row names `## Discovered Tasks` only —
REQ-270 added it — while Step 6 tells the builder to author `## Decisions` too.
**GREEN when:** The same check passes, and the full suite still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure
Ratchet (Step 6.5)**.

## Open Questions

- [x] REQ-270 made Step 8 read a worktree builder's `## Discovered Tasks` from its hand-back
  when the REQ file cannot carry it. Its review found the same gap for `## Decisions` — the
  numbered record of the judgment calls a builder made — which is read by the code review and
  by the end-of-run report you see, both outside Step 8, so under parallel building those
  decisions silently never reach you. Closing it means moving where that rule lives so
  readers outside Step 8 inherit it, and adding a check that every section the builder is
  told to write is one the hand-back is told to carry. Should I process this as a new task?
  → **Yes, add to queue — the full scope**, with one added constraint: it must work from the
  perspective of the installed skill in another repo, not only in this maintainer checkout.
  See **Consumer-Install Constraint** below, which the answer added to the Requirements.

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the technical fix is clear, but it decides how much of
  the pipeline's instruction set reorganizes around parallel building — REQ-270 deliberately
  drew the boundary at Step 8, and widening it touches the review action and the hand-back
  format you read at the end of every run. It is also only worth paying for if you intend to
  keep using worktree fan-out; in a purely serial run nothing here ever fires.

**Answered [2026-08-20]:** User approved via `do-work clarify` at full scope, after asking
how the fix would work and being shown the two reader sites verbatim. The user added the
consumer-install constraint in the same answer; it is a Requirement, not a nice-to-have.

## Consumer-Install Constraint (added by the user at clarify, 2026-08-20)

**Everything this REQ changes must hold in a consumer install, where the suite is vendored
under `.claude/skills/` and only `do-work/` sits at the project root — not just in this
maintainer checkout, where the two coincide.** The user asked for this explicitly. It is the
same class of defect REQ-282 shipped a fix for the day before: a path that resolves in the
maintainer layout and silently resolves to nothing in a consumer's.

Concretely, before this REQ is done:

- **The hand-back path must resolve for an installed skill.** The rule points a reader at
  `do-work/runs/work-<stamp>/REQ-NNN-handback.md`. That is relative to the **project root**,
  where `do-work/` lives in both layouts, so it should hold — **verify it rather than assume
  it**, and state in the rule which root the path is relative to, since a builder resolving
  it against the vendored skill directory finds nothing.
- **No shipped instruction may cite a maintainer-only path.** The property check lives in
  `_dev/tests/`, which is export-ignored and never installed; the *rule* it enforces ships.
  Nothing in `work.md`, `work-reference.md`, or `review-work.md` may reference `_dev/`, the
  check, or this repo's layout. `_dev/tests/shipped-package-reference-contract.sh` already
  enforces that class — make sure it still passes and that it actually covers the new text.
- **Cross-package citations stay resolvable from the citing file's own directory.** The
  three touched action files are all in the core package, so no `../` hop should be needed;
  if one appears, it must be correct at the installed depth, not the repo depth.
- **Verification is an execution, not a reading.** Build a consumer-shaped fixture — suite
  vendored under `.claude/skills/`, `do-work/` at the root — and confirm the instructions a
  builder and an orchestrator would follow resolve there. REQ-282's review used exactly this
  fixture shape and it is cheap to rebuild.
