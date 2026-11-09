```
Continue the do-work run of the folder-layout queue. Two requests are released, one is claimed and merged, two are queued.

advance REQ-715, then do the phase it names (its build is merged, its gate is green and its review passed; do not rebuild it, finish its tail: testing record, review record, lessons, heavy lanes, finalization)
do-work run --fan-out 2

These commands are sufficient; everything below them is context.
```

---

## Reference

Written 2026-11-09T23:22:05Z by the session that ran `do-work audit-current-folder-structure` and then `do-work run` on its five captured requests. The user typed `phandoff` while REQ-715's second gate run finished.

### Queue state

| REQ | What it is | State |
|---|---|---|
| REQ-713 | folder-structure-check sees justfile recipes with default parameter values | released 0.314.1, archived (UR-165 closed) |
| REQ-714 | suite writes its knowledge base and ADR log under docs/ | released 0.315.0, archived (UR-166 closed) |
| REQ-715 | move the dated changelog archives to docs/changelog/ | **claimed**, in `do-work/working/`, merged, tail not done |
| REQ-716 | move the root installer scripts from tools/ to cmd/ | pending, no dependency |
| REQ-717 | suite-managed justfile recipes call a script in one line | pending, `depends_on: [REQ-713]` (satisfied) |

No hold needed a new queue field: the only ordering constraints were the two `depends_on` edges written at capture, and both prerequisites are archived.

### REQ-715 in detail (claimed by this session)

- **Merged.** Branch `worktree-agent-REQ-715-changelog-archives-to-docs-changelog`, commits 742256f3 (the move) and ab39e61c (review fix F1, F2). First merge f2ccb2a7, final merge **1a000d08**. Full merge range: `2b48f7db5b50602c6743a5e58e0c56c88c2a9d85..1a000d08baa4c836f7f378bca03a0165232eebb8`.
- **Done on the record:** Triage, estimate, Plan skip note, Exploration, Scope, Pre-Flight, Decisions (D-01 to D-04), Discovered Tasks, Implementation Summary, Qualification (qualify and scope-drift satisfied on the first merge).
- **Evidence gathered but not yet written into the REQ:**
  - Repository gate `bash _dev/tests/maintainer-verify.sh`: exit 0 at f2ccb2a7 (130s) and exit 0 at 1a000d08 (141s, both Go stages EXECUTING), run from the detached checkout.
  - Review: first pass 82% Partial (F1 broken `Current` link in each archive, F2 stale comment), re-check **96% Pass** after the fix. The final `## Review` body is `do-work/runs/2026-11-09-layout-queue/REQ-715-review.md`; replace the literal `REVIEW_AT` with your own `date -u` instant, then append it.
  - Heavy plan for the cumulative range: all six lanes (coverage is uncertain for `.gitattributes`). **No heavy lane has run yet.**
- **What remains, in order:**
  1. `advance REQ-715 --request-path do-work/working/REQ-715-changelog-archives-to-docs-changelog.md --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file <abs path>/do-work/runs/2026-11-09-layout-queue/REQ-715-probe.sh` (the green-gate record for 1a000d08 is not written yet; if advance wants the gate at the current HEAD, rerun it there first, because HEAD moved by the handoff commits, which touch only `do-work/`).
  2. Append `## Testing`, then `## Review` (stamp `review_at`), then `## Lessons Learned` and `## Orientation` (`kb_status: pending`).
  3. Heavy drain at the merge revision: six lanes, about 5.5 minutes, from the detached checkout. Then `record-commit-hash`, `heavy_verified_at`, `heavy_verified_revision`, and the two Heavy sections between Review and Lessons Learned.
  4. `fold-timing-summary`, release payloads, finalization manifest, `advance ... --finalization-manifest`. UR-167 closes (REQ-715 is its only member), so the archive path is `do-work/archive/UR-167/`.
  5. `worktree cleanup REQ-715`.
- **Release for REQ-715 may be refused.** The drafted entry is `do-work/runs/2026-11-09-layout-queue/REQ-715-entry.md` (0.315.1, patch). The only shipped file in the range is `skills/do-work/CHANGELOG.md` (header links), which the finalizer may count as release metadata and answer `RELEASE-WITHOUT-SHIPPED-CHANGE`. If so, finalize with `NO_RELEASE=1` in `make-manifests.py`; the next release (REQ-716 or REQ-717) carries the header-link change.
- **Uncommitted files at handoff, by name, now held by the handoff commit:** `do-work/working/REQ-715-changelog-archives-to-docs-changelog.md`, `do-work/working/baseline.json`, `do-work/CHECKPOINT.md`, `do-work/runs/2026-11-09-layout-queue/REQ-715-review.md`, `do-work/runs/2026-11-09-layout-queue/REQ-715-entry.md`, and the nine files under `do-work/runs/2026-11-09-layout-queue/helpers/`.
- **Recover result for the claim:** `recover` reported `FINALIZATION-NONE` and `RECOVERY-NONE` at the last queue boundary; bare `advance REQ-715` reports phase `test-gate`, missing the `Testing` section. Continue command: `advance REQ-715`.

### Worktrees

| Path | Verdict |
|---|---|
| `../skill-do-work2-worktrees/worktree-agent-REQ-715-changelog-archives-to-docs-changelog` (branch of the same name, tip ab39e61c, clean, merged) | **ACTIVE**. Its claim is not archived yet. After REQ-715 finalizes, run `worktree cleanup REQ-715` (the canonical command; it uses `git branch -d`). |
| `.git/work-run-2026-11-09-layout/drain-head` (detached at 1a000d08) | **ACTIVE** tool checkout for gate and heavy runs, no claim. Reuse it. It holds one untracked file this session seeded, `do-work/test-durations.tsv`. When the queue is done: `git worktree remove --force .git/work-run-2026-11-09-layout/drain-head` (the next session decides; not run here). |

No foreign session or foreign claim was seen during this run.

### Helpers (committed in the handoff commit)

`do-work/runs/2026-11-09-layout-queue/helpers/`:
- `dw.sh <args>`: runs the CLI in JSON mode and prints phase, gate records and findings. The main repository path is hard-coded.
- `gate.sh <revision> <log>` and `heavy.sh <revision> <log> <lane>...`: run the gate or heavy lanes in the detached checkout; both hard-code the checkout path and print `GATE-EXIT=` or `HEAVY-EXIT=` and the per-lane results.
- `set-estimate.py`, `pau.py`, `insert-section.py`: edit the working REQ (frozen estimate block, the three P-A-U boxes, a section inserted before a heading).
- `add-lesson.py`: adds one bullet to a lesson satellite and refreshes its `lessons-index.md` row.
- `prep-payloads.py` (`BUMP=minor` for a minor bump) and `make-manifests.py` (`UR_CLOSES=1`, `NO_RELEASE=1`): build the release and finalization manifests. Write their output to a scratch folder outside the repository.

The run manifest is `do-work/runs/2026-11-09-layout-queue/manifest.md`. Add REQ-715's final status and the rows for REQ-716 and REQ-717 there.

### Parallelism

- REQ-716 and REQ-717 have no dependency on each other and can be built at the same time: `--fan-out 2`. Integration stays serial.
- Expected overlap at merge, none blocking: both may edit `_dev/tests/staged-skills-contract.sh` (REQ-716 for its `legacy_runtime_paths` entries under `tools/`, REQ-717 if it adds scripts the contract lists) and the root `justfile` (REQ-717 rewrites the managed recipes; REQ-716 only if a recipe names `tools/`).
- Critical path: REQ-715's tail first (it is claimed), then REQ-716, which is the riskier one (it must prove that an update from an already-installed older version still works, or keep a `tools/` forwarder). REQ-717 is a behavior-preserving refactor of two recipes.
- Nothing is held back.

### Heads-up (first ten minutes)

- **Next session:** every REQ in this run needed one review-fix re-merge. Launch the reviewer right after qualification, while the gate runs, and expect a second gate.
- **Next session:** a lesson bullet in a satellite under `skills/` must write project paths as `<project-root>/docs/...`. A bare `docs/kb` in the REQ-714 bullet took `contract-regressions.sh` red after the release commit (fixed in 2ed0be64). Run `bash _dev/tests/contract-regressions.sh` after adding a bullet and before finalizing.
- **Next session:** commit 2ed0be64 changed a shipped file (`skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`) without a version bump. The next release carries it.
- **Next session:** the scope-drift gate compares the `## Scope` bullets with the `## Implementation Summary` bullets. Describe a folder move in a plain sentence in both sections, not as bullets (REQ-714 moved 544 files this way).
- **Next session:** `frontmatter set REQ-NNN write_set "[...]"` writes the list as a quoted string. Edit the line by hand or insert it with a script.
- **Next session:** `do-work/working/baseline.json` is rewritten by each Route B or C pre-flight and is not committed by finalization. Commit it by exact path before the next `recover`, or `recover` refuses.
- **Maintainer:** `main` is about 6 commits ahead of GitHub (releases 0.312.0 to 0.315.0 are not pushed). Other repositories get `do-work audit-current-folder-structure` and the `docs/kb` default only after `git push origin main`.
- **Maintainer:** review findings left as report only are listed in each archived REQ's `## Review` and `## Discovered Tasks` (for example the justfile body-line count still misreads a header continued with a backslash).
