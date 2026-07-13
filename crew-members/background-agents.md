# Background Agent Durability

<!-- JIT_CONTEXT: This file is loaded by any action that fans work out to background or parallel sub-agents — that condition is the contract; illustrative callers today include code-review, work (multi-REQ), pipeline, deep-explore, and file-reorg (Step 3 reference mapping). It prescribes a disk-durable run-directory pattern so fan-out work survives an interrupted, compacted, or corrupted orchestrator session. Not loaded for single-agent in-context work that returns one result. -->

> When you fan work out to background or parallel sub-agents, the chat transcript
> is the worst possible place to keep the results. Make a directory on disk the
> source of truth instead. Sub-agents write their findings to files; the
> orchestrator synthesizes from those files, not from what came back into the
> conversation. The same files give you crash recovery in a fresh session.
>
> **Be honest about the ceiling.** This pattern does not *prevent* the failures
> below — some of them are harness- or API-level faults a markdown skill cannot
> reach. It makes them **survivable and recoverable**. Don't write or speak about
> it as a fix.

## Why This Matters

A fan-out where each sub-agent returns its findings only into the orchestrator's
chat has no durability: if the session is interrupted, compacted, or hits a provider
error mid-run, every finding that came back into the conversation is gone, there is
nothing on disk to recover from, and the whole fan-out must be re-run from scratch.
Disk-as-source-of-truth fixes that *regardless of why the session died*.

## The Durability Pattern

1. **Create the run directory before any spawn.** Make
   `do-work/runs/<action>-<YYYY-MM-DD-HHMMSS>/` first — this directory is the source
   of truth for the entire run. Derive the timestamp from the shell (e.g.
   `date +%Y-%m-%d-%H%M%S`); seconds resolution keeps two runs of the same action
   from colliding on one directory (if it somehow already exists, append a short
   numeric suffix). Nothing should be spawned before this directory exists.

   The run directory is an ordinary **committable** path under `do-work/` (the
   Trail of Intent) — not gitignored, not added to `.git/info/exclude`. Keeping it
   trackable is what lets a run be inspected and survive across sessions while it is
   live: if the user commits their `do-work/` tree mid-run (common in a multi-session
   pipeline), the run travels with the commit instead of sitting as untracked noise or
   being lost. It is not permanent, though — see the lifecycle in step 5: once its
   findings are consumed, the run directory is deleted, and the permanent record is the
   promoted output, never the raw scratch.

2. **Each sub-agent writes its own findings file; returns only a one-line
   status.** Give every sub-agent an output path inside the run directory (e.g.
   `<slice>.md`). The agent writes its *full* findings to that file and returns
   **only a one-line status** to the orchestrator — never the full findings
   inline. This keeps the orchestrator's assembly turn small, which both keeps
   context cheap and shrinks the window in which a large, long-lived assistant
   turn can be corrupted (see Known Failure Mode).

3. **Write a manifest per wave; spawn in bounded waves.** Maintain a
   `manifest.md` in the run directory recording each agent, its assigned slice,
   its expected output filename, and its landed status. Spawn in **bounded waves**
   sized to the harness concurrency limit — not one unbounded fan-out. Update the
   manifest as each wave's files land before launching the next wave.

4. **Synthesize from the files on disk, not from the conversation.** When all
   waves are done, read the findings files from the run directory and assemble the
   final output from them. Never synthesize from what agents "said" in chat. This
   is the property that makes the run recoverable: synthesis behaves identically in
   the original session and in a fresh recovery session that never saw the spawns.
   Once synthesis succeeds, **mark the run complete** — write `Status: complete` to
   the manifest. A completed run must never be offered for resume (see Known Failure
   Mode); it is no longer live state.

5. **Delete the run directory once its findings are consumed.** The run directory is
   working memory, not the deliverable. The moment its findings have been *synthesized
   and promoted to the run's permanent output* — a report, captured REQs, or artifacts
   copied into `do-work/deliverables/` — that scratch is redundant, so the action that
   owns the run **deletes the directory as its final step**. This is the normal end of
   the information's lifecycle, not an optional tidy-up: it keeps `do-work/runs/` from
   growing without bound (nothing else prunes it) while the promoted output carries the
   durable record. State the trigger as this *condition* — "findings consumed and
   promoted" — for whatever action owns the run; do not tie it to a fixed list of
   callers. The delete is a filesystem change that rides the normal commit flow; no
   action force-commits on its own. `Status: complete` is the bridge for the window
   between synthesis and deletion — if a run is abandoned after finishing but before its
   owner deletes it, `actions/cleanup.md` sweeps any `complete`-marked directory as a
   safety net (an incomplete one is left alone, because it may still be resumable).

## Local-ignore snippet (for genuinely-transient paths — NOT run dirs)

Run directories are committable (step 1) and must **never** be added to any ignore list.
But some sibling paths *are* genuinely transient and must stay out of git regardless of
install layout — `do-work/pipeline.json` (live pipeline state), a vendored engine install,
a `build/` artifact. The shipped `.gitignore` can't reach a project-root path from a nested
`.claude/skills/do-work/` install, so those paths append to the enclosing repo's
`.git/info/exclude` (local-only — never committed, never shipped). This is the canonical
snippet they reference; substitute the path being ignored:

```bash
exclude=$(git rev-parse --git-path info/exclude 2>/dev/null) || exclude=""
if [ -n "$exclude" ]; then
  git check-ignore -q <path> 2>/dev/null || echo '**/<path>' >> "$exclude"
fi
```

`git check-ignore -q` already succeeds when *any* ignore source covers the path, so the
append only fires when genuinely needed and never duplicates. The appended pattern carries
a `**/` prefix because a pattern with an interior slash is root-anchored, while
`check-ignore` tests the cwd-relative path — without the prefix, a run from a repo
subdirectory would never see its own append, re-appending on every run while the path
stayed unignored; `**/` makes the pattern match at any depth, so the check and the pattern
agree. The `|| exclude=""` fallback keeps the guard a clean no-op outside a git repo —
without it the failed command substitution leaves the assignment nonzero, which aborts the
script under `set -e`. Resolve the exclude file with `git rev-parse --git-path info/exclude`
— **not** `$(git rev-parse --show-toplevel)/.git/info/exclude`, which breaks in linked
worktrees and submodules where `.git` is a file, not a directory (the redirect fails with
"Not a directory" and the path is left un-ignored).

## Known Failure Mode & Recovery

**The reasoning-block corruption (reasoning-model harnesses).** On harnesses that
use a signed *thinking* / reasoning block (e.g. Claude with extended thinking), an
assistant turn that holds an open signed thinking block *and* long-running
background spawns can be corrupted if it is interrupted or re-stitched while still
open. Once the malformed turn is written to the session log, every resume replays
it and re-throws the same error — typically `HTTP 400 — "thinking blocks cannot be
modified"`. The session becomes **permanently un-resumable**. The corruption is
per-transcript: other sessions on the machine are unaffected.

This is a harness/API-level fault. This pattern cannot prevent it — it makes it
recoverable. The one-line-status rule (step 2) also shrinks the corruption window
by keeping the assembly turn small, but that is mitigation, not prevention.

**Recovery procedure:**

1. **Do NOT resume the poisoned conversation.** Resuming replays the corrupt turn
   and re-throws the error every time.
2. Start a **fresh session** and re-invoke the same action.
3. Let the action **detect the most recent incomplete run directory** (glob
   `do-work/runs/<action>-*`; if several match, take the newest by timestamp) and
   read its manifest. A manifest marked `Status: complete` means that run already
   finished — skip it; only a run *without* that marker is resumable.
4. **Re-spawn every agent whose findings file is absent on disk.** Verify against
   the filesystem — do not trust the manifest's per-row label, because a crashed
   orchestrator may never have updated it. Agents whose findings file already exists
   are done; do not re-run them.
5. **Synthesize from disk** as normal.
6. The poisoned transcript can be deleted once recovery succeeds.

## Match the Pattern to the Harness

One invariant holds no matter how the fan-out runs: **disk is the source of truth —
sub-agents write findings to files, the orchestrator synthesizes from those files,
and recovery reads from disk.** What changes between harnesses is only *how much of
the machinery you hand-roll*. Use the highest rung the harness supports; the
invariant above carries down all three.

1. **Native orchestration engine.** The harness exposes a deterministic fan-out
   primitive with journaled resume — a `workflow` / `pipeline`-style API that caps
   concurrency, returns structured per-agent output, and replays cached results when
   re-run. Prefer it: bounded waves (step 3) and the structured findings hand-off
   (step 2) come for free, and its journaled resume covers the *orchestration* —
   re-running replays already-completed agents instead of re-spawning them. It also
   usually runs detached from the orchestrator turn, so the reasoning-block
   corruption above is less likely to strand you in the first place. That resume is
   not a substitute for the disk files, though: still write each slice's findings to
   the run directory. The engine's journal recovers the *run*; the on-disk files
   recover the *synthesis* and keep it identical across harnesses — belt and
   suspenders, not an either/or.

2. **Manual parallel/background spawns.** The harness can spawn concurrent
   sub-agents but offers no orchestration engine. Hand-roll the pattern exactly as
   steps 1–4 describe: run directory, per-slice findings files, manifest, bounded
   waves, synthesize from disk.

3. **Sequential in-context.** No parallel or background support at all. Do not skip
   the pattern — run the slices **one at a time in the current context**, but still
   create the run directory, still write each slice's findings to its file as you
   complete it, still update the manifest, and still synthesize from disk. A
   sequential run that crashes halfway is recoverable because the completed slices
   are already on disk.

## Manifest Format

Keep it small and append-friendly. A minimal `manifest.md` (this example uses the
`code-review` action's six dimensions; your slices will differ):

```markdown
# Run Manifest — code-review-2026-05-28-143012

Run dir: do-work/runs/code-review-2026-05-28-143012/
Concurrency: 4 (wave size)
Status: in-progress   # flips to `complete` after synthesis succeeds

| Agent | Slice | Output file | Status |
|-------|-------|-------------|--------|
| 1 | Consistency | consistency.md | done |
| 2 | Architecture | architecture.md | done |
| 3 | Security | security.md | pending |
| 4 | Performance | performance.md | done |
| 5 | Test Coverage | test-coverage.md | pending |
| 6 | Automated Checks | automated-checks.md | pending |
```

Per-row status is just `pending` (not yet confirmed on disk) and `done` (findings
file written and present); the happy path moves rows `pending → done` only. **There
is no `missing` status to write** — a crashed orchestrator can't be relied on to
record one. Recovery is derived from the filesystem instead: re-spawn any row whose
findings file is **absent on disk**, regardless of its label. The run-level
`Status:` line is the completion signal — `in-progress` until synthesis succeeds,
then `complete`; a `complete` run is never offered for resume.
