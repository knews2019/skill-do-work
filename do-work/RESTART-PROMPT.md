```
do-work clarify
do-work run --fan-out 2
This command is sufficient; everything below it is context.
```

---

## Reference

### Queue

14 pending, 1 pending-answers, 0 blocked, 0 in progress. `do-work/working/` holds only `baseline.json`.

`do-work clarify` comes first because **REQ-427 is `pending-answers`** and asks one question worth answering before the batch continues: UR-081 set the Go floor at 1.26.1, and measurement during REQ-407's review found the module builds and passes all six packages' tests at `go 1.23.0`. As shipped, that floor excludes anyone on Go 1.23, 1.24 or 1.25 from installing or updating do-work at all. Answering it is cheap; leaving it costs every consumer without a current toolchain. The REQ carries the full measurement table and three options.

### Wave 1 is REQ-408 and REQ-426

Both are dependency-ready and touch disjoint files, which is why `--fan-out 2`:

- **REQ-408** — shared request, schema, dependency, atomic-file and repository packages. New packages under `skills/do-work/tools/do-work-cli/internal/`. Head of the serial chain: REQ-409 through REQ-420 each gate on their predecessor, so after this wave the batch advances one REQ at a time whatever fan-out is passed.
- **REQ-426** — restores setuid/setgid/sticky bits that REQ-407 strips from `Justfile`, `CLAUDE.md` and `.claude/settings.json`. Touches `managedsection` and `suiteinstall`. Mechanical, and its Red-Green Proof carries the measured A/B.

Every ordering constraint is a `depends_on` field; nothing here needs reading for `do-work run` to do the right thing.

### Completed this session

| REQ | Release | Implementation |
|---|---|---|
| REQ-406 | 0.245.0 | `2ca25d7` (merge range `ad354e2..2ca25d7`) |
| REQ-390 | 0.246.0 | `59105df` (merge range `b3ea43d..59105df`) |
| REQ-407 | 0.247.0 | `f45cdca` (merge range `0bc1480..f45cdca`) |
| REQ-425 | 0.248.0 | `04b8120` (merge range `b2cbe87..04b8120`) |

All four are archived with their hashes recorded and verified by `tools/checks/record-commit-hash.sh --verify`.

### Worktree verdicts

None. All four builder worktrees were removed and their branches deleted with `git branch -d` run from the integration branch — the refusal-or-accept that proves each merge landed. `git worktree list` shows only `/home/user/skill-do-work`. Nothing to re-check and nothing to remove.

### Heads-up — things that will bite in the first ten minutes

- **The container ships no maintainer toolchain.** A fresh container needs ShellCheck 0.11.0, `just` 1.43.0, and Go 1.26.1 (`go env -w GOTOOLCHAIN=go1.26.1+auto` is enough — the toolchain downloads from the module proxy, which is reachable; `go.dev` is not). Without these `_dev/tests/maintainer-verify.sh` fails on its version preconditions before running anything. **Next session: check this first.**
- **Never set `QUEUE_KANBAN_BROWSER` for the canonical gate.** The container's Chromium is 141, the build REQ-375 deprecated, and `TestBrowserBehaviorCompletionCompanionsKeepReadableContrast` fails here at HEAD on *both* 141 and 151 — a Linux-headless difference from the macOS pass REQ-375 recorded. The gate's default skipped browser lane is the correct configuration. For timeline probes specifically, fetch Chrome for Testing 151.0.7922.174 from `storage.googleapis.com/chrome-for-testing-public/` and run them individually; all 21 pass there. **Whoever owns the board next: this deserves its own REQ, because the strict lane is currently unrunnable on Linux.**
- **Commit each unit as its tests pass.** A usage limit killed a builder mid-run this session with eight modified files and five new packages uncommitted; only the worktree saved it. Every builder brief now says this and it is worth keeping.
- **`tools/checks/scope-drift.sh` misreads prose.** It treats every backticked token in a "Files I will touch" bullet as a declared path, so a code span in the rationale produces a phantom "declared but never touched" line. Keep Scope bullets to one backticked path each until REQ-414 ports the checks to Go, where the fix is folded.
- **Build the Implementation Summary's file list from the merge range, not from the hand-back's prose.** Parsing prose produced 32 of 39 files on REQ-407 and seven phantom drift lines with it.
- **No Routine is armed.** An hourly self-check-in ran the drain through three usage-limit interruptions this session and was deleted at handoff, because it fired into *that* session rather than a new one. If you want the same safety net, create an hourly Routine bound to the new session with a prompt that (a) no-ops in one line while a limit is in force, (b) re-reads `do-work/CHECKPOINT.md` and the queue rather than trusting context, and (c) deletes itself once the queue holds no claimable pending REQs and the gate exits 0. Those three properties are what made it useful rather than noisy.

### Review posture that earned its cost

Each REQ got five review dimensions and three refutation lenses per finding. The lenses matter as much as the dimensions: across four REQs they refuted the large majority of findings by measurement, and what survived was real every time — a 10% flaky test already inside the canonical gate, two edge-of-range navigation defects, a goroutine race in an install signal handler, and a BOM that hard-blocked installs. Two findings were also settled by direct orchestrator measurement rather than by another agent, which is worth doing whenever a claim is cheaply testable.

One workflow bug is fixed and worth not reintroducing: an errored verifier must never count as a refutation. Fewer than two live votes means *unverified*, which survives.
