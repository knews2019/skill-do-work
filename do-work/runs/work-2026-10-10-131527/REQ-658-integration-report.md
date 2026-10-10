# REQ-658 integration report (finalize --auto-manifest)

- Release 0.305.104, "do-work-cli Builds the Finalization Manifest From the Judged Inputs and Checks the Tree First". Archive `do-work/archive/UR-145/REQ-658-finalize-auto-manifest.md`. Finalization commit `90a3680b`: cleanup_complete, no blocked paths or reason codes, first try. UR-145 is closed: REQ-659, REQ-660 and REQ-661 moved into `do-work/archive/UR-145/` with `input.md`.
- Dogfood: the manifest came from the merged `finalize --auto-manifest ... --emit`. Its commit_paths matched the helper-built manifest exactly (22 paths), and `advance --finalization-manifest` accepted it.
- Range: pre 84d56197. First merge 457608cc. Review-fix re-merge fe788bf4 (recorded as `commit:`). Integrator fix e09c25d8 on the builder branch. No conflicts.
- Scope: the pre-dispatch `## Scope` list used bare paths, so scope-drift flagged all 8 files. I added backticks only (D-20).
- Review 88%, Approve. I fixed F1: without `--emit` the command finalized past the `advance` phase gate and archived a Triage-only REQ. `--emit` is now required (D-21). This departs from the REQ text, which implied a direct mode. I also fixed F2 (now uses `ResolveTarget`), F3 (prose says which inputs are judged and which fields the command fills) and F4 (index check runs first). Delta re-review 93%, Pass, F1-F4 closed. Report only: F5 (refusal codes differ from D-16, corrected in D-21), F6 (no repo test for emit then advance), F7 (REQ-660 F3 dirt readers still stale), N1 (F4 not pinned by a test), N2 (D-17 wording), N3 (`work-reference.md:719` omits the failure inputs).
- Gate (`DO_WORK_FAST_STAGE_REUSE=off`) exit 0 at 457608cc (173 s, load 4.65 before) and at fe788bf4 (153 s, load 2.44 before). Probe exit 0 both times. Green record at fe788bf4 via `record-green-gate`.
- Heavy lanes at fe788bf4, all exit 0 and executed: do-work-cli-integrations 81 s, staged-skills 41 s, updater 73 s, installer 35 s. The earlier drain at 457608cc was also green.
- Timing events: handback-merge 1, verification-gate 4, review 2, folded into `## Timing`. Builder-work skipped (the hand-back had already landed).
- Lessons: one `guard-on-one-entry-path` bullet in `lessons-do-work-cli.md`. REQ-660's two links and REQ-659's one link now point to `archive/UR-145/`. Index row is 20000 tokens. `contract-regressions.sh` passes after finalization.
- Cleanup: builder worktree removed, branch deleted with `-d`, pruned. Both drain checkouts removed.
- Left dirty: the siblings' untracked `REQ-*-handback.md` files (not mine) and this report (untracked, for the coordinator).
- Not mine, for the coordinator: `do-work/working/REQ-688-capture-files-example-and-fence-fix.md:65` still cites the flat `do-work/archive/REQ-661-...` path, which now lives under `archive/UR-145/`.
- Report-only discovered tasks: `queueWriterLabel` duplicates `DefaultWriterLabel` (negligible). A refused `finalize --manifest` deletes a stray payload folder (pre-existing, negligible).
