---
id: REQ-113
title: "Confirm: migrate the first prose read site onto queue-kanban frontmatter"
status: completed
status_changed_at: 2026-08-05T19:59:16Z
created_at: 2026-08-05T19:24:00Z
user_request: UR-021
addendum_to: REQ-112
domain: general
prime_files: []
tdd: false
claimed_at: 2026-08-05T19:59:16Z
completed_at: 2026-08-05T20:01:30Z
route: A
depends_on: [REQ-112]
maintenance: false
builder_decided: true
---

# Confirm: Migrate the First Prose Read Site

## What

While building REQ-112, the builder had to decide whether to also change one action file to actually *use* the new command, or to ship the command with tests only. It chose tests only. This follow-up confirms whether that matches your intent.

The background, since this will read cold later: the repo now has a `queue-kanban frontmatter get <file> <field>` command that reads a field out of a REQ file. Roughly 95 places across `actions/` still read REQ fields by hand instead — usually a few lines of `awk` or `grep` written inline in the prose of an action. The new command exists so those places *can* stop hand-rolling it, but nothing has been switched over yet, so today the command has no user.

**Why this is your call rather than the builder's:** switching one of those places over means editing an action file, and every action file has to keep working for an agent that has no compiler available. So the real question isn't "does the command work" — it does — but "what should an action say to do when the compiled tool isn't there?" That's a judgment about how much the skill leans on optional tooling, which is a standing preference of yours, not a technical detail.

## What the Builder Chose

Ship the command with tests only. REQ-112's diff stayed entirely inside `tools/queue-kanban/`, so it could be reviewed as a self-contained tooling change.

## What Would Change

If you'd rather see a site migrated: one low-risk action would gain a line naming the command as the preferred way to read a field, plus a documented fallback for when the binary isn't built. The candidate the builder had in mind is `actions/commit.md` Step 3, which checks whether an archived REQ has a terminal-success status — a check the new `--in-set terminal-success` flag does exactly, and one of the five places where a documented Red Flag warns that hand-rolling it drops `completed-with-issues`.

## Open Questions

- [x] Should one action file now be switched over to use the new command, as a worked example? → **Yes, migrate `actions/commit.md` Step 3 as the worked example.** The user delegated the decision ("use common sense") rather than answering, and the builder's recommended "No" does not survive that test: a command with zero callers closes none of the finding it was built for, and the per-action floor judgment I deferred on is small for a single site. `commit.md` Step 3 is the right first one because its terminal-success check is one of the five documented Red Flag sites for the `completed`-vs-`completed-with-issues` bug, so the migration removes a real hazard rather than demonstrating a pattern. Scope stays one site: the remaining ~94 are not in scope here.
  Recommended: No — leave the command unused for now. It's available when any action wants it, and each switch-over can be reviewed on its own.
  Value: keeps tooling changes and prose changes in separate reviews, and avoids committing the whole skill to a pattern before one real use has been seen.
  Risk: a command with no users can turn out to fit no real call site, and the hand-rolled reads the census flagged all stay as they are until someone starts. Cheap to reverse — adopting it later is a one-line prose edit per action, with nothing to undo.
  Also: Yes, migrate `actions/commit.md` Step 3 as the worked example; or yes, but pick a different site.

## Full Context

See `do-work/user-requests/UR-021/input.md`. The originating decision is D-01 in `REQ-112`'s `## Decisions`.

---

## Triage

**Route: A** - Simple

**Reasoning:** One named file, one named step, and the shape to follow (`next-req`/`now`'s accelerator form) already exists. No exploration.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add the preferred command plus a documented shell floor to `actions/commit.md` Step 3's terminal-success check. Accelerator shape only: gated on an already-built binary, explicit "never build it for this", floor spelled out. Verify by running *both* prescribed commands against a real archived REQ.
- [x] **[APPLY]:** One edit to `actions/commit.md` Step 3. No other action touched; the remaining read sites are out of scope.
- [x] **[UNIFY]:** `git diff --stat` → 2 files (`actions/commit.md`, this REQ). Ran both prescribed commands against `do-work/archive/UR-023/REQ-115-*.md`: the floor `awk` prints `completed`, and `frontmatter get … --in-set terminal-success` exits 0. No debug artifacts. `_dev/tests/contract-regressions.sh` at its 7 pre-existing failures, unchanged.

## Implementation Summary

**What was done:** Migrated the first prose frontmatter read site onto `queue-kanban frontmatter`, giving the new command its first caller.

**Files changed:**
- `actions/commit.md` (modified) — Step 3's terminal-success check now names `frontmatter get <req-file> status --in-set terminal-success` as preferred when the binary is already built, with an `awk` floor and an explicit "never build the tool for this"

**Scope:** exactly one site. The other ~94 prose reads the census counted are untouched and remain candidates, not commitments.

## Decisions

- **D-02: `--in-set` rather than `get … --normalize` plus a string compare.** DECIDE & STATE. `--in-set` answers by exit code, so the prose needs no value comparison at all — which is the half a hand-rolled version gets wrong (testing the literal `completed`). It also normalizes aliases first, so a `status: done` REQ is not skipped.
- **D-03: the floor is stated as a full working command, not "read the field yourself".** DECIDE & STATE. This repo's own rule is that prescribed commands must actually emit what the following logic consumes; a vague floor is how an action reads fine and fails on a real repo. Both commands were run against a real archived REQ before this shipped.

## Qualification

Passed — 1 shipped file in the diff, the change traces to the resolved question, and both prescribed commands were executed rather than assumed.

## Testing

**Tests run:** both prescribed commands against `do-work/archive/UR-023/REQ-115-status-normalize-warning.md`; `_dev/tests/contract-regressions.sh`.

**Result:** floor `awk` prints `completed`; preferred command exits 0. Contract suite unchanged at its 7 pre-existing failures.

**Red-green validation:** not applicable — this is a prose change to an action file, with no runnable assertion. The proof is that both prescribed commands execute correctly against a real REQ, which is the checkable property for prescribed-command prose.

## Review

**Approve** — gives the new command its first caller, removes a documented hazard, and preserves the shell floor.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Minor:** `actions/inspect.md` Step 3 carries a near-verbatim copy of this same check and was deliberately not migrated — scope was one site. That copy is now inconsistent with `commit.md`, which is a known and recorded state rather than an oversight; REQ-114's Candidate B covers consolidating the pair.
**Acceptance:** Pass — both commands verified against a real archived REQ.
**Follow-ups created:** None (REQ-114 Candidate B already covers the sibling site).

## Lessons Learned

**What worked:** Answering the question as the user's proxy when they delegated it, rather than re-asking. The builder's original "No" was defensible in isolation and wrong in context — a command with no callers closes none of the finding it was built for.

**What didn't:** Nothing failed here, but the split that produced this REQ is worth noting: REQ-112 shipped a surface and deferred its only proof of usefulness to a follow-up. That is defensible for review isolation and it did mean the command sat unused through two commits and a merge. A surface plus one caller in the same REQ would have been reviewable too.

**Worth knowing:** The accelerator shape has three required parts, and the third is the one that gets dropped — preferred command, documented floor, **and an explicit prohibition on building the tool to get the value**. Without the third, "preferred" quietly becomes "required" the first time an agent runs `go build` to satisfy it.

## Orientation

`actions/commit.md` Step 3 now reads a REQ's terminal-success status through `queue-kanban frontmatter --in-set` when the binary is available, falling back to a stated `awk` command otherwise. This is the first prose caller of the CLI added in REQ-112, so the census's central finding now has one site actually fixed rather than only made fixable. Leaf change to one action; no contract moved.
