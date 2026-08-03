---
id: REQ-082
title: The fan-out hand-back file has no legal write location
status: pending
created_at: 2026-08-03T17:09:21Z
user_request: UR-016
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-073, REQ-084]
batch: audit-remediation-external
addendum_to: REQ-073
write_set: [actions/work-reference.md, _dev/tests/contract-regressions.sh]
---

# The Fan-Out Hand-Back File Has No Legal Write Location

## What

Fan-Out Dispatch makes a per-builder output file mandatory — `REQ-NNN-handback.md` inside
`do-work/runs/work-<timestamp>/` — and `crew-members/background-agents.md` is explicit that the
sub-agent writes that file itself. Worktree Dispatch Mode is equally explicit that a builder **never
writes the main tree**, and `do-work/` exists in the main tree only. There is no location satisfying
both rules, so the mandatory hand-back has no legal execution.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

An agent reaching this contradiction has three moves and two of them corrupt the run. It can write the
main tree (violating sole-integrator, the rule that makes worktree isolation mean anything); it can
write the worktree's own `do-work/runs/…`, where the file lands in the builder's branch and gets swept
into the merge as committed scratch while the orchestrator reads nothing; or it can skip the file and
return findings in conversation, which is the exact failure the durability pattern exists to prevent.
Nothing in the shipped prose tells it which.

This is the output-direction twin of a trap REQ-073 already closed for inputs. Its requirement 8 fixed
brief *delivery* — "the brief must reach the builder as prompt content or an absolute main-tree path.
A repo-relative path resolves inside the worktree, against its own stale tracked copy of `do-work/`" —
and the return path has the identical resolution problem with no equivalent sentence.

## Context

The three statements, verbatim:

- `actions/work-reference.md:313` (Fan-Out Dispatch, the mandatory run-directory table):
  > | per-builder output | `REQ-NNN-handback.md` — branch, file manifest, integration seams |

  introduced by "**The run directory is mandatory here, not optional.**" (`:307`).

- `crew-members/background-agents.md:44-48`:
  > **Each sub-agent writes its own findings file; returns only a one-line status.** Give every
  > sub-agent an output path inside the run directory … The agent writes its *full* findings to that
  > file and returns **only a one-line status** to the orchestrator — never the full findings inline.

- `actions/work-reference.md:275` (Sole integrator):
  > The builder never writes the main tree or its branch.

  reinforced by `:273` (State stays home): "`do-work/` — the queue, `working/`, `CHECKPOINT.md` —
  exists in the main tree only."

**Two secondary problems in the same neighbourhood.**

1. *State stays home* enumerates three things (`the queue, working/, CHECKPOINT.md`). `do-work/runs/`
   did not exist when that sentence was written and is not in the list, so a reader can argue the run
   directory is out of scope — which is exactly the closed-enumeration failure `CLAUDE.md` warns
   about. The sentence needs its condition stated, not a fourth item appended.
2. `crew-members/background-agents.md:33-41` makes the run directory an ordinary **committable** path
   under `do-work/`, deliberately, so a mid-run commit carries it. Combined with a builder writing
   into it from inside a worktree, that turns run scratch into branch content — and REQ-084's new
   merge-base probe would then correctly flag the builder for writing queue state. The two REQs must
   not disagree about whether that file is a violation.

**Why no check caught it.** The contradiction is between two files that no assertion compares, and
`_dev/tests/contract-regressions.sh` has no assertion touching the hand-back file at all. It is also
unreachable in practice today: everything since REQ-073 shipped has been built serially (see
REQ-085), so no run has ever produced a hand-back file.

## Detailed Requirements

1. **State the one exception explicitly in `actions/work-reference.md`**, at *Sole integrator* — the
   sentence the exception modifies — and reference it from the Fan-Out Dispatch table row rather than
   restating it. A builder may write **exactly one** path: its own
   `do-work/runs/work-<timestamp>/REQ-NNN-handback.md`.
2. **The path reaches the builder the same way its brief does** — as an absolute main-tree path, never
   repo-relative. Requirement 8's existing trap sentence covers the mechanism; say that it applies in
   both directions instead of writing a second copy of the reasoning.
3. **The builder never stages, commits, or merges that file.** It is a main-tree working file owned by
   the orchestrator's run directory, not branch content. State this as part of the exception, because
   the natural reading of "you may write it" includes committing it.
4. **The exception is bounded to that one filename.** Not "the run directory", not "files under
   `do-work/runs/`" — one path per builder, derived from its own REQ id. A builder writing a sibling's
   hand-back file, the manifest, or anything else remains a sole-integrator violation.
5. **Restate *State stays home* as a condition, not a list.** Replace the three-item enumeration with
   the rule it is trying to express — every path under `do-work/` is main-tree-only and
   orchestrator-owned — and mark any examples as illustrative. Per `CLAUDE.md` → Closed Enumerations
   Go Stale, grep for other enumerations of the same set and generalize each.
6. **`manifest.md` stays the orchestrator's.** `background-agents.md:51-56` has the orchestrator
   maintain it per wave; make sure the amended prose cannot be read as licensing a builder to update
   it. This is the natural over-reach from requirement 1.
7. **Reconcile with REQ-084 in whichever order they land.** REQ-084 adds a probe for a builder branch
   carrying `do-work/` changes. The hand-back file must not trip it — which it cannot, if requirements
   2 and 3 hold, because the file is written to the main tree and never committed. Whichever REQ lands
   second must confirm that in its Qualification rather than assuming it.
8. **Add a contract assertion pinning the exception.** The failure mode is a later maintenance pass
   reading "the builder never writes the main tree" as absolute and deleting the carve-out as
   redundant, which silently restores the contradiction. Assert that the exception exists and that it
   names a single path.

## Constraints

- **`maintenance: true`.** The candidate fix narrows two shipped instructions (`Sole integrator`'s
  absolute prohibition; *State stays home*'s enumeration), so `crew-members/maintenance.md`'s
  delete-before-you-add rule governs. Requirement 5 is a replacement, not an addition, and
  requirement 1 should cost one sentence — if the exception needs a paragraph, the shape is wrong.
- **Do not weaken sole-integrator anywhere else.** The builder still never touches the queue, a
  status, an archive move, `CHECKPOINT.md`, `actions/version.md`, `CHANGELOG.md`, an integration seam,
  or a sibling's anything. This REQ opens one file, by name.
- **Do not describe the durability pattern as preventing failures.** `background-agents.md:11-14` and
  Fan-Out Dispatch's closing line (`actions/work-reference.md:317`) both require the
  survivable-not-prevented framing to be carried, not softened.
- **No new durable coordination state.** The forbidden-token sweep
  (`_dev/tests/contract-regressions.sh:132-137`) must stay green: no lock, heartbeat, claim registry,
  or liveness probe enters via the hand-back path.
- The run directory's lifecycle (created before any spawn, deleted once consumed) is
  `background-agents.md`'s, and is cited rather than restated.

## Dependencies

`addendum_to: REQ-073`, which introduced the mandatory hand-back. `related: REQ-084` for requirement 7.
No `depends_on`: buildable immediately in either order, and REQ-085's live run is the thing that would
have surfaced this, not a prerequisite for fixing it.

## Builder Guidance

**Certainty: Firm on the diagnosis and on the chosen direction; open on wording and placement.** The
direction was decided by the user at capture time — see the resolved question below — so do not
re-open it or re-derive the trade-off.

Prefer the shortest wording that survives a maintenance pass. This section of
`actions/work-reference.md` is already the file's largest, and REQ-073's Builder Guidance rule still
applies: anything restating what Worktree Dispatch Mode already says gets cut rather than written.

## Open Questions

- [x] The hand-back file must be written somewhere, and both candidate homes are currently forbidden.
  Should the builder be granted a narrow main-tree write, or should the file be dropped in favour of
  the builder reporting its manifest in its reply?
  → **Grant the narrow write.** Resolved by the user at capture time (2026-08-03, ask-tool prompt
  during `do-work capture-request`). The builder may write exactly one path — its own
  `REQ-NNN-handback.md`, by absolute main-tree path — and never commits it. Reasoning: the file exists
  *because* the transcript is not durable (`crew-members/background-agents.md` § Why This Matters), so
  dropping it to preserve an absolute prohibition trades a real recovery property for a tidier rule.
  **Out of scope as a result:** the return-the-manifest-in-the-reply shape, and any broadening of the
  exception past that one filename.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Red-Green Proof

**RED prompt/case:** A contract-suite assertion (`_dev/tests/contract-regressions.sh`) that fails
while `actions/work-reference.md` mandates a per-builder hand-back file without naming a legal write
location for it. Concretely: assert that the *Sole integrator* paragraph contains a bounded exception
naming `handback`, and that *State stays home* no longer expresses its scope as a three-item list.

**Why RED now:** `grep -c "handback" actions/work-reference.md` finds the mandate at `:313` and
nothing at `:275`; `_dev/tests/contract-regressions.sh` has no assertion mentioning the hand-back file.

**GREEN when:** The assertions pass; the full contract suite stays green; and a reader can answer
"where does the builder write its hand-back file, and who commits it?" from `actions/work-reference.md`
alone, without opening `crew-members/background-agents.md`.

**Manual proof, for the part an assertion cannot reach:** the amended prose must let a human dispatch
two builders and receive two hand-back files without either builder writing anything the orchestrator
did not authorize. That check is REQ-085's run, not this REQ's — note the dependency in the review
rather than claiming it here.

**Validation:** User confirmed — the direction was chosen from an explicit two-option prompt at capture
time; the contradiction itself was verified by reading all three cited statements against the tree.

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the
external audit, and the batch constraints.

---
*Source: external audit finding F2, second claim (P1) — "background-agent rules require builders to
write into the main-tree run directory, while worktree rules forbid builders from writing the main
tree" — accepted by `do-work validate-feedback` triage as the sharpest finding of the six.*
