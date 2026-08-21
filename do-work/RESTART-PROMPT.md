```
do-work clarify
do-work run

This command is sufficient; everything below it is context.

Run `do-work clarify` first: two REQs sit at pending-answers and cannot be selected until you
answer them. Then `do-work run` drains the queue in dependency order.

Serial is the recommendation, not a limitation. 29 REQs are pending and only four have
verified write-set relationships (encoded as depends_on gates). The other 25 have not been
checked against each other, so `--fan-out` would be dispatching on unverified disjointness.
If you want concurrency, `do-work run --fan-out 2` is the honest ceiling until someone
audits the remaining write sets.
```

---

## Reference

### Session state

- **Checkout:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`, branch `main`, tree **clean**.
- **Worktrees:** none. No `worktree-agent-*` branches exist, merged or unmerged. Nothing to remove.
- **`do-work/working/`:** empty. No claim in flight, no foreign claim, `## In Progress (interrupted)` is empty.
- **Version:** 0.216.2. `bash _dev/tests/maintainer-verify.sh` exits 0 at `2314327`.

### What this session finished

**REQ-258 — Split the prescribed shell behavior suite per script.** Route B, review 94%, Acceptance Pass. Merged? Not applicable — serial mode, no worktree, no merge range. Committed at **`1cc1836`**, hash recorded at `35fb513` and confirmed with `record-commit-hash.sh --verify`. Archived at `do-work/archive/UR-056/REQ-258-split-the-shell-behavior-suite-per-script.md`. Nothing remains on it.

`_dev/tests/prescribed-shell-scripts-behavior.sh` keeps its path and exit-status contract but is now a 35-line runner. The 76 cases live one file per script in `_dev/tests/prescribed-shell-cases/`, over `_dev/tests/prescribed-shell-harness.sh`. **A REQ that adds a case now writes that script's case file, not the runner.**

### Heads-up list — things that will bite in the first ten minutes

1. **A second session was writing this repo during this one, and may still be.** Four commits landed from outside this session: `031c546` (clarify), `1311300` (duration rounding, 0.216.0 → 0.216.1), and `2314327` (UR-062 capture: REQ-303, REQ-304, REQ-305, plus an addendum to REQ-263). None collided with REQ-258 — verified with `git show --stat` on each. **Whoever resumes should confirm no other session is mid-write before running,** because this checkout has no lock and none is coming.
2. **`RESTART-PROMPT.md`'s previous version was wrong about the bottleneck, and REQ-300 exists to fix that class.** The old text said `prescribed-shell-scripts-behavior.sh` is written by five REQs at "at most one per wave." That file is no longer written by case-adding REQs. This file replaces that text; REQ-300 sweeps the rest.
3. **`tools/checks/qualify.sh` will FAIL any REQ that relocates code.** It reads `git diff` with no `-M`/`-C`, so a moved line is an added line and every pre-existing `TODO`/`console.log` inside moved text trips the debug-artifact gate. REQ-258 hit it and overrode it with evidence (`git show HEAD:<file> | grep` proving the lines pre-exist). **Do not un-check `[UNIFY]` when this happens** — prove the lines pre-exist and record the override. REQ-301 is the fix, awaiting your approval.
4. **Six reservation markers exist for REQ-300 through REQ-305** under `do-work/.req-reservations/`, all committed and all matched by real REQ files. Nothing to reap.

### Ordering, and where it is encoded

Every constraint below is a `depends_on` field, not prose. `do-work run` honors them with no reading.

| REQ | `depends_on` | Why |
|---|---|---|
| REQ-263 | `[REQ-300]` | REQ-300 rewrites REQ-263's own frontmatter; a builder must not hold it while that happens |
| REQ-271 | `[REQ-300]` | same |
| REQ-264 | `[REQ-300, REQ-263]` | same, **plus** REQ-263 and REQ-264 both write `skills/do-work/tools/checks/qualify.sh` and both will write `_dev/tests/prescribed-shell-cases/qualify.sh` — they must not run concurrently |
| REQ-301 | `[REQ-263, REQ-264]` | same two files again; the gate is pre-set so it holds whenever you approve it out of `pending-answers` |

Pre-existing gates, unchanged: REQ-281 → REQ-280, REQ-285 → REQ-284, REQ-292 → REQ-291.

**Not gated on purpose:** nothing forces REQ-300 to run before its own instance list goes stale. REQ-263/264/271's `write_set` values self-heal at Step 5.5, which overwrites the field from the fresh Scope declaration — so those three instances are cosmetic. REQ-300's durable value is this file and the eleven Coverage rows in `decisions/audits/2026-08-11-defensive-surface.md`.

### Awaiting your answer (`do-work clarify`)

- **REQ-301 — let qualify tell a moved line from an added one.** The heads-up above is the case for it. The risk if declined is habituation: a gate that cries wolf on a whole category of change trains builders to wave it away, and that gate is the one that catches real leftover instrumentation.
- **REQ-302 — check whether capture under-sizes reorganization REQs.** `effort_estimate: trivial` produced a 5-minute P50 for REQ-258's 19-file restructure. One data point, so the REQ asks the question before proposing a fix. Cheap to decline.

### Parallelism analysis

- **Safe to run concurrently:** unknown for 25 of the 29 pending REQs. Their write sets have not been checked against each other, and `write_set` is display-only — it gates nothing and the merge, not the field, is what proves non-interference.
- **Must not run concurrently:** REQ-263 with REQ-264 (and REQ-301 with either) — shared `qualify.sh` files. REQ-300 with REQ-263/264/271 — REQ-300 edits their files. All four encoded above.
- **Critical path:** REQ-300 → REQ-263 → REQ-264 → REQ-301, four deep. Everything else is a leaf or a two-deep pair. Starting there is starting on the longest chain.
- **Nothing held back.** No REQ is `blocked`, and no `blocked_check` probe exists in the queue.
