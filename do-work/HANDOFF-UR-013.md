# Handoff — UR-013 (Parallel builds without coordination overhead)

**Written:** 2026-08-03T15:22Z · **Stopped at:** user request, after REQ-073's commit + metadata commit
**Tree state:** clean, nothing uncommitted, nothing in `do-work/working/` (except `baseline.json`, a preflight artifact — gitignored)
**Version:** 0.166.0 · **HEAD:** `9ba2cda`-plus-metadata-commit

---

## What landed

All three of UR-013's captured REQs are **completed, reviewed, committed, and archived**. Each got its own version bump and changelog entry.

| REQ | Title | Route | Commit | Review |
| --- | --- | --- | --- | --- |
| REQ-071 | Crash recovery must respect a live claim before stripping and re-queueing | B | `5c39899` | 94% Approve |
| REQ-072 | Go utility allocates REQ ids and version numbers and verifies release consistency | C | `5db22ea` | 91% Approve |
| REQ-073 | Fan-out dispatch: N concurrent builders under one queue owner | C | `9ba2cda` | 89% Approve with follow-ups |

Releases: **0.164.0** (REQ-071), **0.165.0** (REQ-072), **0.166.0** (REQ-073).

### REQ-071 — recovery stopped being unconditional
Crash recovery used to strip thirteen generated sections from every `REQ-*.md` in `do-work/working/` and re-queue it. Nothing is committed before Step 9, so a crash-and-restart destroyed a finished Plan/Exploration/Scope. Now a classification gate sits **above** the existing substeps (substeps 1–3 are byte-identical): a REQ named in `CHECKPOINT.md`'s `## In Progress (interrupted)` record is an own crash and recovers as before; anything else — **including the common case of no checkpoint at all** — is a foreign claim, left byte-identical and reported. Takeover is offered only past three hours, or immediately when `claimed_at` is unparseable/future-dated/absent, and **only a human authorizes it**. Unattended runs report and continue.

### REQ-072 — three new subcommands on the shipped board tool
`queue-kanban next-req` (next free REQ number), `next-version <patch|minor|major>` (bumps the `**Current version**: ` line, reads it back to confirm, prints it), `verify` (eight read-only invariant probes). Wired as **optional accelerators** into `actions/capture.md`, `actions/work.md` Step 9, and `actions/forensics.md` as Check 14 — each with a stated missing-`go` fallback. The tool never writes `CHANGELOG.md`. New files: `tools/queue-kanban/{allocate,release,verify}.go` + three `_test.go` (21 tests).

### REQ-073 — the builder cap was two sentences
The invariant is now **one queue owner per checkout** (was "one active REQ, one coder context"). Builders aren't owners, so any number may build concurrently; two queue owners stays banned. New **Fan-Out Dispatch** subsection inside Worktree Dispatch Mode. Net cost in the reference file: 29 lines — because everything there was already written per REQ.

---

## Queue state: two follow-ups, both belonging to UR-013

```
do-work/queue/
├── REQ-074-recovered-req-loses-its-status-change-timestamp.md   (pending-answers)
└── REQ-075-write-set-display-only-reason-went-stale.md          (pending)
```

**UR-013 is intentionally still open** in `do-work/user-requests/UR-013/` — both follow-ups carry `user_request: UR-013`, and neither is terminally resolved, so the UR folder must not be consolidated yet. That is correct per `actions/work.md` Step 8's UR-final check, not an oversight.

### REQ-074 — `pending-answers`, needs a human answer
Crash recovery resets a REQ to `pending` without stamping `status_changed_at`; the manual reset in `actions/forensics.md` Check 1 does stamp it. The board's state timer then dates a just-recovered REQ from its creation day. Three options are written out in its `## Open Questions` (stamp it / narrow the field's rule / do nothing). **Run `do-work clarify` to answer it.** Pre-existing, low severity, no data loss.

### REQ-075 — `pending`, ready to build
Five files still justify "nothing schedules on `write_set`" with the premise REQ-073 falsified ("one REQ runs at a time"): `actions/board.md` ×2, `docs/board-guide.md`, `tools/queue-kanban/prime-do-kanban.md`, `actions/capture-reference.md`. The conclusion survives, the premise doesn't — which is the dangerous shape, because a reader reasoning forward from the dead premise concludes `write_set` should now gate, which is forbidden. Rated **Important** for that reason. The corrected wording already exists in three places to copy from (`actions/work.md` § Rules, `actions/work.md` Step 5.5, `CLAUDE.md` § Shipped Tooling). Requirement 4 asks for an assertion so a sixth copy can't land quietly. `tdd: true`, `maintenance: true`, `depends_on: []` — buildable immediately.

---

## Not done — deliberately, and it matters

1. **`do-work cleanup` was never run.** Step 10 normally runs it at the end of the loop. The queue holds two unresolved REQs so nothing is stranded, but the cleanup pass hasn't swept.
2. **`do-work/CHECKPOINT.md` is stale.** It still describes the session that ended 2026-08-01 (`last_completed: REQ-070`, "Queue empty"). **This is now load-bearing** — REQ-071 made the checkpoint crash recovery's input. It is not wrong in a way that loses work (a stale checkpoint means claimed REQs get treated as foreign claims and left alone, the safe direction), but it should be rewritten.
3. **KB handoff deferred on all three REQs.** Each carries `kb_status: pending`. `kb/` exists. The offer was batched rather than interrupting the run three times — see the restart prompt.
4. **REQ-073's live two-builder acceptance test was never run.** This is the honest gap in the batch and REQ-073's review says so: the whole run was serial, so the capability the REQ *enables* is unverified. Grep assertions prove the prose is right; they cannot prove two builders compose. Its review lists the exact procedure under Suggested additional testing — two non-overlapping REQs in two worktrees, then a deliberately overlapping pair that must **fail** at `git merge --no-ff --no-commit`.

---

## Things a fresh session should know

- **A polarity bug was found by running the tool, not by the tests.** REQ-072's version probe was built literally from "the current version is strictly greater than the newest `CHANGELOG.md` entry" — true only *while composing* a release. After the entry is written they are equal, so the probe fired on this repo's own clean tree while all 21 tests were green (code and fixtures encoded the same misreading). It now asserts *agreement*, with the direction naming the cause, and the strictly-greater check moved to ordering *within* the changelog — which is what actually catches the duplicate version numbers `CLAUDE.md` records. Recorded as REQ-072 D-02.
- **A contract assertion can fail silently.** `grep -roh <pattern> | wc -l` under `set -euo pipefail` aborts the whole suite when the pattern is absent — so REQ-073's first "the invariant exists" check exited 1 with **no output at all**, in exactly the case it existed to catch. Fixed with `{ … || true; }`. `three_actions_count`'s sibling counter (`three_attempt_count`) has the same shape and is filed as a Discovered Task on REQ-073.
- **The restatement sweep earned its keep twice.** REQ-071 found two self-contradicting lines in `actions/work.md`; REQ-073 found eight sites arguing from a premise it had just falsified. The productive sweep question was *"what was justified by this premise?"* — not *"where else is this string?"* The five REQ-075 sites share no phrase with anything that was edited.
- **`maintenance: false` on two instruction-editing REQs.** REQ-071 and REQ-073 were both removal/narrowing passes on the skill's own instructions, and both had the marker false, so `crew-members/maintenance.md` never loaded. That was honored on purpose (marker-only by design — `actions/work.md` Step 6 5a explains why inferring it is worse). Filed as a Discovered Task on REQ-073: either capture's assessment under-fires on this shape, or the marker is narrower than the crew file's JIT_CONTEXT.
- **Scope-drift check caught me once.** REQ-072: I logged the decision to touch `docs/forensics-guide.md` (D-03) but never added it to `## Scope`. Writing the decision down is not declaring the file — different checks read different sections.
- **`verify` currently passes clean on this repo** (`./tools/queue-kanban/queue-kanban verify --repo-root .` → exit 0). Useful as a first sanity check in a fresh session.

---

## Restart prompt for a fresh session

Paste this verbatim:

```
Resume the do-work skill repo at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2.

Read do-work/HANDOFF-UR-013.md first — it has full context. UR-013's three REQs
(REQ-071/072/073) are done, committed, and archived at v0.166.0; the tree is clean.

Do these in order, and stop to show me the result of each before moving on:

1. Rewrite do-work/CHECKPOINT.md — it still describes the 2026-08-01 session and
   claims the queue is empty. It's now load-bearing: REQ-071 made the checkpoint
   crash recovery's input, so a stale one changes recovery behavior. Use the
   Session Checkpoint Template in actions/work-reference.md, session_depth
   moderate (3 REQs), and record REQ-074 + REQ-075 under Still Queued.

2. Run do-work cleanup. It was skipped when the run stopped. Expect it to leave
   UR-013 open in do-work/user-requests/ — REQ-074 and REQ-075 both belong to it
   and neither is terminally resolved, so consolidating would be wrong.

3. Offer the KB handoff for all three REQs at once (they carry kb_status: pending,
   and kb/ exists). Follow actions/kb-lessons-handoff.md and ask before writing.
   The lessons worth promoting: the premise-goes-stale-not-the-conclusion pattern
   from REQ-073's restatement sweep, and REQ-072's "run a prose-derived probe
   against a healthy tree before trusting its fixtures."

4. Then ask me which of these I want next, don't pick for me:
   - do-work clarify  → answer REQ-074 (three options are already written out)
   - do-work run REQ-075 → the five stale write_set justifications (Important
     finding from REQ-073's sweep; corrected wording already exists to copy)
   - the live two-builder test of REQ-073's fan-out dispatch — the one thing the
     batch shipped unverified. Procedure is in REQ-073's ## Review under
     Suggested additional testing.
```
