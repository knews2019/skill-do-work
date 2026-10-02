---
id: UR-131
title: 'Add a low-disk-space probe to the board VERIFY band'
created_at: 2026-10-02T13:58:27Z
requests: [REQ-625]
word_count: 5
---
# Add a low-disk-space probe to the board VERIFY band

## Summary
The user pasted a proposal for a free-disk-space probe and asked `do-work validate-feedback` to triage it. The triage accepted the probe, the thresholds, the finding shape, the path reduction, the read-only contract, the tests and the docs, each with amendments. It pushed back on three parts, which are NOT captured here: the top-3 largest-directories walk inside the finding, the always-visible free-space chip in the board header, and the run-loop sentence that pauses fan-out dispatch on a critical finding. The user then asked to capture the accepted items. This invocation captures only; it does not start implementation.

## Extracted Requests
- REQ-625 (Add a low-disk-space probe to the board VERIFY band): one probe in `verify.go`, a build-tagged free-space helper, fixed thresholds, focused tests, prime and lessons entries.

## Batch Constraints
- Read-only probe. It measures and reports; nothing is deleted. Cleanup stays human-consented.
- No new Go module dependency. Stdlib `syscall` on unix, `kernel32.dll` through `syscall.NewLazyDLL` on Windows, as `atomic_replace_windows.go` already does.
- No header chip, no directory walk, no run-loop prose. Those were pushed back in triage and are out of scope.
- Shipped files change, so the REQ is a release: changelog entry plus version bump per the release rules.

## Triage amendments carried into the REQ
- Thresholds are two fixed constants, warn at 10 GiB and critical at 3 GiB. The "10 % whichever is smaller" rule was dropped because on small disks critical would fire before warn.
- Subject is the measured directory (repo root or worktree), not the mount point. A mount point is almost never inside the repo, so the existing path reduction would render it as `<path outside this repository>`.
- The unix build tag list must be narrower than `atomic_replace_unix.go`: aix, illumos and solaris have no `syscall.Statfs`.
- Dedupe by `Stat_t.Dev` (cast to uint64), not by `Statfs_t.Fsid`, whose field names differ per OS. Use `Bavail`, not `Bfree`.

## Source proposal (pasted by the user in the validate-feedback turn, verbatim)
> ````text
> # Add a free-disk-space probe to the do-work board's VERIFY band
>
> ## Why
> During a long `do-work run --fan-out 2` (Oct 2026, a game repo), browser-QA runs wrote 400–800 MB of
> screenshots per before/after comparison into `tests/output/`. The repo grew to about 20 GB in a few hours and
> the disk nearly filled, but nothing in the suite noticed. Builders, full test gates and git all fail
> in confusing ways at 0 bytes free, and an unattended run can corrupt a half-written finalization or
> checkpoint. The board's VERIFY band is already where the run's operator looks, so low disk space
> belongs there.
>
> ## What to build
> A new verify probe in `do-work-board/tools/queue-kanban/verify.go`, wired into
> `collectVerifyFindings` next to `appendWorktreeFindings`. It lands in both the CLI `verify` report and
> the board's VERIFY band through `attachVerifyFindings`, with no separate UI path.
>
> 1. **Measure free space** on the filesystem holding the repo root, and on each worktree directory the
>    worktree probe already enumerates (`git worktree list --porcelain`) when it is on a different
>    filesystem. Dedupe by device id. Use `syscall.Statfs` on unix and `GetDiskFreeSpaceEx` on Windows,
>    behind a build-tagged helper like the existing `atomic_replace_*` split. On an unsupported
>    platform, add the probe to `SkippedProbes` and do not fail.
> 2. **Thresholds:** warn below **10 GiB free or 10 % free, whichever is smaller**, and treat below
>    **3 GiB** as critical. Keep both as named constants. A static threshold beats a growth-rate model
>    (YAGNI).
> 3. **Finding shape:** `Category: "low-disk-space"`, `Subject: <mount point, reduced>`, `Fixable: false`.
>    Detail: `"<free> free of <total> (<pct>%) — below <threshold>; largest repo-local dirs: <top 3 with sizes>"`.
>    The top-3 list is a shallow `du`-style walk of the repo root's direct children, including gitignored
>    ones, capped at about 200 ms or a fixed entry budget so `serve` stays cheap. When the cap is hit, say
>    "(partial)" rather than guessing. Remedy:
>    `"free space: clear regenerable QA output (e.g. tests/output), finished builder worktrees (do-work cleanup), browser caches"`.
>    Keep the remedy generic: the suite must not know any project's output folders.
> 4. **Paths:** the mount point goes through `reduceAbsolutePaths` like every other finding, so static
>    snapshots stay shareable.
> 5. **Read-only:** like every verify probe, it measures and reports and never deletes anything.
>    Cleanup stays human-consented (`actions/cleanup.md`).
> 6. **Header indicator:** also show a small always-visible free-space chip in the board header (for
>    example "disk 27 GiB free"), even when nothing is wrong, coloured by the same thresholds. A run
>    operator should see the trend before it becomes a finding. The data comes from the same helper,
>    carried in `generatedBoardData`; serve recomputes it per request outside the mtime cache, as claim
>    age already is.
> 7. **Run-loop hook (optional, small):** in `do-work/actions/work.md` → Step 10 (Loop or Exit), and in
>    `work-reference.md` → Fan-Out Dispatch, add one sentence: before dispatching a new wave, a critical
>    `low-disk-space` finding pauses new dispatch, and integration continues. Do not stop the run.
>    Finishing integrations frees worktrees, so pausing only new builds is the self-healing choice.
>
> ## Tests (focused, each naming the failure it pins)
> - Probe reports nothing above the thresholds and a finding below them, using a fake statfs.
> - Critical vs warning boundary at exactly 3 GiB and at 10 % / 10 GiB.
> - Two worktrees on the same device produce one finding (dedupe by device id).
> - An unsupported platform puts the probe in `SkippedProbes`, never silently clean.
> - The board payload carries the finding with the mount point reduced (no absolute path), extending
>   `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths`.
> - The top-dirs walk respects its budget and marks "(partial)".
> - A JS behaviour test shows the header chip in all three states.
>
> ## Docs
> Add the probe to `prime-do-kanban.md`'s probe list and a one-line lesson in `lessons-do-kanban.md`
> (family `disk-space-blind-spot`): a long fan-out run with browser QA can fill a disk in hours, and the
> board is where the run operator looks.
> ````

## Full Verbatim Input
> ```
> capture the ones that are accepted
> ```
