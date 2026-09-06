```
do-work run --fan-out 3
This command is sufficient; everything below it is context.

You are the agent working on do-work itself, on branch claude/do-work-queue-drain-4ee2xl (the
integration branch; never main). Push often with git push -u origin claude/do-work-queue-drain-4ee2xl
and open no pull request. End every commit message with the attribution trailers your session
prescribes. Four REQs are in do-work/working/ with claims from the previous session; take each over
with the canonical command the oracle prints (do-work-cli recover --take-over REQ-NNN) and then let
advance drive the phase. The maintainer gate in this container needs an environment scrub: run
do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh (add --heavy for the heavy tier), never
bare _dev/tests/maintainer-verify.sh. Finalize one REQ with handoff-tools/finalize-req.sh; its
neighbour FINALIZATION-RECIPE.md says why each step is in its order. Every run artifact you need
(review syntheses, the judged plan and its verified patch, the workflow scripts, a changelog draft)
is under do-work/runs/work-2026-09-05-231943/. Reviews are three independent lenses plus a
synthesizer that reproduces every finding; remediate rather than defend, and correct record claims
in place with strikethrough. Sections in a REQ record must carry repo-relative paths, one per line,
and no backticked token in Scope that is not a path.
```

---

## Reference

Written 2026-09-06 after the previous session hit its usage limit mid-run. Head of the integration
branch when this was written: see the commit that adds this file. VERSION is 0.305.8. Nothing is
dirty in any checkout. Four background workflows died to the session limit (REQ-598 build, REQ-602
verify, REQ-597 review) and one landed (REQ-600 review); nothing else went wrong.

### In-flight REQs, in the order to work them

**REQ-600 — put the SIGPIPE trap in the shell prime, fix the one shipped block.**
Merged as `a25c7522`, full range `9e00a092cf29842506bea920137b52c952a62638..a25c7522566bea9d9d29c382e159b6a10157a9f1`.
Qualify and scope-drift satisfied; Qualification and Testing written; review recorded at 83 percent
(`do-work/runs/work-2026-09-05-231943/REQ-600-review.json`, and the record's Review section). Remaining,
in order: remediate the six findings the Review lists (the important one is the false safe-zone
sentence at `_dev/primes/prime-shell-commands.md:53`; the review gives the replacement; also the
scanner's stale comment at `_dev/tests/quiet-grep-pipeline-scanner.sh:18` and `:11`, the block prose
"first hit wins", the prime's line 51 list split, and this record's two mis-stated sentences), run the
four guards, commit; write Lessons Learned and Orientation; finalize as a **release** at 0.305.9 with
`handoff-tools/changelog-req600.md` (already headed 0.305.9; re-check the heading is unused). Carry the
lesson satellite lines inside the finalization commit via `EXTRA_COMMIT_PATHS` (see below). Worktree:
`worktree-agent-REQ-600-sigpipe-prime`, clean, merged. No uncommitted files.

**REQ-597 — correct the stale claims across the rest of the shell guide and its two callers.**
Merged as `d5cf28b`, full range `804a8ba32129a3cd12a4aaa7e89346db1b95115c..d5cf28b996a6deb0a0df908cbe4aa722cf2a6ad8`
(three builder commits `a1e652f`, `6913dc4`, `7df6488`; three shipped prose files, +23/-23). Qualify and
scope-drift satisfied; Implementation Summary, Decisions, Discovered Tasks, Qualification and Testing
written; the oracle's next phase is **review** and the review never ran (session limit). Remaining: run
`handoff-tools/req597-three-lens-review.workflow.js` (edit nothing but the paths if your scratchpad
differs; it clones the repo and executes the prescribed blocks against fixtures), remediate, write
Lessons and Orientation, finalize as a **release** (next patch version after REQ-600's). The builders'
per-sentence evidence is in `REQ-597-handback.md` and `REQ-597-builders.json`; their fixtures lived in
the old session's scratchpad and are gone, but `REQ-597-verification.json` and the hand-back say how to
rebuild them. Worktree: `worktree-agent-REQ-597-guide-and-callers`, clean, merged. No uncommitted files.

**REQ-602 — repoint fifteen lesson-satellite links, add a satellite link check.**
Merged as `ef8274b`, full range `66e9992f559627a280b113eda4fd1ad476016f07..ef8274bef8ea83c6961b6ad3d1d12848c011c5e8`
(five files, +69/-18: three maintainer satellites, `_dev/tests/audit-lockins.sh`, `do-work/lessons-index.md`).
The orchestrator re-ran the red (15 FAIL lines at `a70a04f`) and the green (lock-ins pass at head and
on the merged tree); the independent verifier never ran. Remaining: canonical `advance REQ-602
--diff-range <range above>` for qualify and scope-drift, Qualification and Testing, a three-lens review
(model it on the REQ-600 script; a fixture-driven verifier is in `req602-build-and-verify.workflow.js`),
Lessons and Orientation, finalize (**not a release**: nothing under `skills/` changed). Its hand-back
also reports the shipped satellite's canonical URLs; act on that report if any did not resolve.
Worktree: `worktree-agent-REQ-602-satellite-links`, clean, merged. No uncommitted files.

**REQ-598 — close the nil-handle panic in transaction rollback.**
Not built. Route C; the judged plan (Plan B, decide the handle once at its open) is in the record's Plan,
Exploration, Scope and Pre-Flight, and in `REQ-598-judged-plan.json`. The verified tree is
`req598-final.patch` (557 lines; `git apply --check` passed on `521e4a7`, and the three files it touches
have not changed since), the exact test is `req598-final-test-func.go.txt`, and the rejected minimum
shape is `req598-minimum-shape.patch`. `depends_on` now includes REQ-602 because both edit
`_dev/tests/audit-lockins.sh` (REQ-602 is merged, so the gate clears when it archives). The builder
must take its own RED evidence (seam and test first, panic at `root.Mkdir`), then the restructure, then
the lock-in rewrite (pin at zero), then canary, `-U0` diff, differential, `-race`, `GOOS=windows go vet`,
gate; `req598-build-and-verify.workflow.js` is that brief. Worktree `worktree-agent-REQ-598-rollback-handle`
is stale at `9e00a09` and clean; reset its branch to the integration head before building
(`git -C <path> checkout -B worktree-agent-REQ-598-rollback-handle claude/do-work-queue-drain-4ee2xl`)
or remove it and let the builder create `worktree-agent-REQ-598-decide-once`. Release when finalized.

### Queue

- **REQ-605** (finalization `diff-tree` without `-m`): free; `do-work run` selects it first.
- **REQ-601** (phantom-script claims in seven shipped callers, plus the guide's tie-break sentence):
  waits on REQ-597.
- **REQ-603** (protected-inventory launcher and shim): waits on REQ-597 and REQ-601.
- **REQ-604** (atomic-download occupancy rule and unchecked stat): waits on REQ-601.

### Parallelism (mirrored into the gates above)

Safe together: REQ-600, REQ-597, REQ-602 and REQ-605 touch disjoint files (prime and scanner comment;
three prose files under review; `_dev` satellites and the lock-in; `internal/finalization`). Must not run
together, and gated: REQ-598 with REQ-602 (`audit-lockins.sh`); REQ-601, REQ-603 and REQ-604 with each
other (all edit `skills/do-work/docs/prescribed-shell-primitives.md`; REQ-603 also edits `commit.md`
and `inspect.md`, which REQ-597 just rewrote, hence its REQ-597 gate). Critical path: REQ-597 review
and finalize, then REQ-601, then REQ-603. `--fan-out 3` covers the three in-flight REQs at review or
finalize while REQ-605 builds.

### Canonical recover results (writer evidence)

`do-work-cli --format text recover` reports all four working claims as
`RECOVERY-TAKEOVER-AVAILABLE`, writer `vm:/home/user/skill-do-work` at `do-work/CHECKPOINT.md` lines
35 (REQ-597), 38 (REQ-600), 40 (REQ-598), 42 (REQ-602), each `takeover available; claim preserved`.
Takeover command per claim: `do-work-cli recover --take-over REQ-NNN`. `recover-claim` requires
`--assume-sole-writer`; there is no other session.

### Worktrees (none removed; re-check all three conditions before removing)

- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-597-guide-and-callers` — **ACTIVE**: merged
  (`d5cf28b`), clean, claim still in `working/`. After REQ-597 archives:
  `git worktree remove /home/user/skill-do-work-worktrees/worktree-agent-REQ-597-guide-and-callers && git branch -D worktree-agent-REQ-597-guide-and-callers`
- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-600-sigpipe-prime` — **ACTIVE**: merged
  (`a25c752`), clean, claim in `working/`. After REQ-600 archives:
  `git worktree remove /home/user/skill-do-work-worktrees/worktree-agent-REQ-600-sigpipe-prime && git branch -D worktree-agent-REQ-600-sigpipe-prime`
- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-602-satellite-links` — **ACTIVE**: merged
  (`ef8274b`), clean, claim in `working/`. After REQ-602 archives:
  `git worktree remove /home/user/skill-do-work-worktrees/worktree-agent-REQ-602-satellite-links && git branch -D worktree-agent-REQ-602-satellite-links`
- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-598-rollback-handle` — **ACTIVE**: no
  commits beyond the old base `9e00a09`, clean, claim in `working/`; reset or replace before building.

### Lesson satellites owed (work.md Step 8 substep 4)

The previous session skipped the satellite append for every REQ it archived, then backfilled the three
maintainer satellites in `a38a8c4`. Still owed, and to be carried inside REQ-600's finalization commit
with `EXTRA_COMMIT_PATHS="skills/do-work/tools/do-work-cli/lessons-do-work-cli.md _dev/primes/lessons-shell-commands.md do-work/lessons-index.md"`
(a change under `skills/` is a release, which is why they ride that commit):
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (canonical-URL form, family marker first):
  REQ-583 (keep each mutation applied while writing the test that catches it; two behaviours were inert,
  not dead), REQ-592 (a request's own statement of why something is safe is a claim to verify), REQ-593
  (three fixes in one change is where one covers for another's missing test), REQ-599 (a test that
  defeats the bug does not pin the rule; put the forbidden token where the most specific wrong fix would
  still see it). Each REQ's Lessons Learned section carries the wording.
- `_dev/primes/lessons-shell-commands.md`: REQ-600's own line, once its Lessons are written.
- Recompute the two index rows (tokens = (bytes + 3) / 4; families = sorted marker set).
- Fix the satellite headers that cite "Step 8 substep 7" (there is no such substep).

### Heads-up

- The four workflows that died at 08:5x UTC left nothing dirty; do not look for half-applied work.
- `git merge -F -` does not read stdin here; write the message to a file. `pkill -f <pattern>` kills
  your own shell when the pattern is in its command line. `set -o pipefail` plus `grep -q` after a pipe
  is the defect REQ-593/594/600 exist for; the guard will fail your shell if you write it.
- `finalize-req.sh` needs the full 40-character merge hash and an empty index; the changelog entry file's
  first line must be `## <version> — <Title> (<date>)`; the finalizer refuses a title already used.
- The `heavyverification` package shows two failures in reviewers' clones under a bare `go test`; under
  `gate.sh` on the main checkout the heavy tier passed (wall 356s). Record the pair, do not chase it.
- Backticked tokens in a record's Scope or Implementation Summary are read as paths: `-` and `.sh`
  both tripped scope-drift this run. One repo-relative path per line; no comma-grouped bullets.
- The maintainer prime lessons index (`do-work/lessons-index.md`) has no test; recompute rows by program.
- Two structural suggestions were offered to the user and not started (by-hand fallbacks as tested
  scripts; one whole-file fixture-driven audit REQ for the shell guide). Wait for the user's word.
