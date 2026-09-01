```
do-work clarify
do-work run REQ-488
do-work run
This command is sufficient; everything below it is context.
```

---

## Reference

- **Why three lines:** `do-work clarify` answers the one `pending-answers` REQ (REQ-489, consent only). `do-work run REQ-488` builds the selector fix first, named explicitly because explicit naming bypasses dependency gating. `do-work run` then drains the queue serially; until REQ-488 lands, a plain run wrongly excludes 20 REQs (below), so do not skip the middle line.
- **Parallelism:** serial (fan-out 1). The selector's ready set is wrong until REQ-488 lands, and after that most ready REQs edit `do-work-cli` internals (`requeststate`, `nextselection`, `doctor`) or shared action prose (`work.md`, `work-reference.md`) and would collide at merge. Board-only REQs (REQ-439, REQ-456, REQ-482, REQ-486) could pair with one CLI REQ under `--fan-out 2`; not encoded, so plain `do-work run` stays serial. Critical path: REQ-488 → REQ-468 → REQ-469 → REQ-470/REQ-471 → REQ-472.
- **Queue-state edits made by this handoff:** removed the literal `depends_on: []` line from `do-work/queue/REQ-488-keep-empty-inline-frontmatter-lists-empty.md` so the current selector can see it (the bug it fixes is exactly that line being read as a dependency named `[]`). Nothing else in the queue was touched.
- **Selector bug in force (REQ-488):** `do-work-cli --format json next` reports `DEPENDENCY-MISSING: missing dependencies: []` for every REQ carrying `depends_on: []`: REQ-437, 438, 439, 442, 443, 444, 448, 449, 450, 451, 452, 453, 456, 457, 458, 461, 468, 482, 484, 485, 486. Root cause: `FieldValue` in `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` copies `ListValues` with `append([]string(nil), ...)`, which turns an empty slice into `nil`, and `listValue` then falls back to the scalar `[]`. Explicit `do-work run REQ-NNN` bypasses the gate if anyone needs one of these before REQ-488 lands.
- **REQ-440 (this session):** completed and archived at `do-work/archive/REQ-440-refuse-non-file-static-board-targets.md`. Serial mode, no merge range. Implementation commit `cdf1732c`, metadata commit `8fde4f48`, release 0.260.3. Review Approve 96%, Acceptance Pass, three Minor findings in the report only; the stale doc-comment finding is on `do-work/prose-backlog.md`. UR-083 stays open (other members pending). No uncommitted files.
- **Gate fix (this session):** `_dev/tests/shipped-shell-thinness.sh` SC2034 fixed standalone as `2d140f63` (0.260.2). `bash _dev/tests/maintainer-verify.sh` exited 0 at 2026-09-01T19:4x against the tree containing REQ-440. REQ-469, REQ-470, REQ-471, REQ-472 in the queue address the unrelated-gate-failure hold shape that stopped the previous session.
- **REQ-489 (pending-answers):** canonical `complete`/`fail` removal of a checkpoint entry deletes only the `- REQ-NNN:` header and leaves the indented detail lines orphaned (`checkpointWithoutClaim` in `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`). Orphans from REQ-418 and REQ-440 were removed by hand in this session's checkpoint rewrite.
- **Worktree verdicts:**
  - `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` (`main`, `8fde4f48` + this handoff commit): ACTIVE, clean after the handoff commit. No `worktree-agent-*` branches exist.
  - `/private/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/tmp.2qQUpYmVN0` (detached at `4b52ae4d`, the REQ-420 claim commit): FOREIGN, not a `worktree-agent-*` name, dirty. Modified: `skills/do-work-toolbox/crew-members/general.md`, `skills/do-work/actions/review-work.md`, `skills/do-work/actions/work.md`, `skills/do-work/crew-members/general.md`, `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`; untracked `do-work/lessons-index.md`. Looks like a REQ-478/REQ-479 lessons-era scratch tree left under `$TMPDIR`. Left byte-identical. The user decides: inspect with `git -C <path> diff`, then `git worktree remove --force <path>` only if it is scratch.
  - `/private/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/tmp.gYIGToUyFs` (detached at `4b52ae4d`): FOREIGN, same dirty file set as above. Same disposition.
- **Leftovers in `do-work/working/` (untracked, not claims):** `baseline.json` (REQ-440 preflight record, `launched: false`), empty `.req478-release/` and `.req479-release/`. Safe to delete by hand; the work loop ignores them.
- **Heads-up for the first ten minutes:** the previous `RESTART-PROMPT.md` (REQ-420 handoff) is superseded by this file. `do-work/CHECKPOINT.md` has an empty `## In Progress (interrupted)` section and no foreign entries. Version is 0.260.3; the next release must be strictly greater.
