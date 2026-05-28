# Background Agent Durability

<!-- JIT_CONTEXT: This file is loaded by any action that fans work out to background or parallel sub-agents — code-review, work (multi-REQ), pipeline, and deep-explore. It prescribes a disk-durable run-directory pattern so fan-out work survives an interrupted, compacted, or corrupted orchestrator session. Not loaded for single-agent in-context work that returns one result. -->

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
chat has no durability. If the orchestrator session is interrupted, compacted, or
hits a provider error mid-run, every finding that came back into the conversation
is gone — there is nothing on disk to recover from, and the whole fan-out has to
be re-run from scratch.

Moving the source of truth to disk fixes the *recoverability* problem regardless
of why the session died.

## The Durability Pattern

1. **Create the run directory before any spawn.** Make
   `do-work/runs/<action>-<YYYY-MM-DD-HHMM>/` first — this directory is the source
   of truth for the entire run. Derive the timestamp from the shell (e.g.
   `date +%Y-%m-%d-%H%M`) so reruns and recovery can find it. Nothing should be
   spawned before this directory exists.

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
3. Let the action **detect the existing run directory** (`do-work/runs/<action>-*`)
   and read the manifest.
4. **Re-spawn only the agents whose output files are missing.** Agents that already
   wrote their findings file are done — do not re-run them.
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
   re-run. Prefer it: it gives you the bounded waves (step 3), the structured
   findings hand-off (step 2), and the manifest-plus-re-spawn recovery (Known Failure
   Mode) for free, and it usually runs detached from the orchestrator turn, so the
   reasoning-block corruption above is less likely to strand you in the first place.
   Still write each slice's findings to the run directory — the engine's journal and
   the on-disk files are belt-and-suspenders, not an either/or, and the on-disk files
   are what keep synthesis recovery-identical across harnesses.

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

Keep it small and append-friendly. A minimal `manifest.md`:

```markdown
# Run Manifest — code-review-2026-05-28-1430

Run dir: do-work/runs/code-review-2026-05-28-1430/
Concurrency: 4 (wave size)

| Agent | Slice | Output file | Status |
|-------|-------|-------------|--------|
| 1 | Consistency | consistency.md | done |
| 2 | Architecture | architecture.md | done |
| 3 | Security | security.md | missing |
| 4 | Performance | performance.md | done |
| 5 | Test Coverage | test-coverage.md | pending |
| 6 | Automated Checks | automated-checks.md | pending |
```

Status values: `pending` (not yet launched), `done` (findings file written), `missing`
(launched but no file landed — the recovery target). On recovery, re-spawn only the
`missing` rows.
