```
do-work run REQ-624
This command is sufficient; everything below it is context.

Continue the already-claimed request from its recorded review phase. The user requested this
handoff and authorizes the next session to continue it. Its implementation is already merged
and verified green. Finish independent review, the selected heavy verification, and canonical
finalization. The other two requests in UR-130 are complete; finish this request and stop.
```

---

## Reference

Written 2026-09-09 for the user's “phandoff and commit everything” request. Integration branch:
`main`. Current version: `0.305.37`. All remaining implementation and evidence are committed;
this handoff commit also includes the canonical checkpoint refresh. No push was requested.
The old restart prompt described a different, completed run and has been replaced.

### Current state and remaining order

**REQ-624 — Trim the shipped changelog to 50 entries.** ACTIVE, claimed, merged and verified;
independent review has not started. Canonical `advance REQ-624` returns `outcome: success`,
`phase: agent judgment: review`, `status: claimed`. This is an in-place continuation of the
existing claim. Its preserved workflow sections are the resume state, not a plan held only here.

- Request: `do-work/working/REQ-624-trim-shipped-changelog-to-50-entries.md`.
- Source merge: `959cb107d4b95395dcd54d5c44f7509197b97114`.
- **Full merge range:** `49c61ba76932abd5d6a40f7509ac3f56193dfea8..959cb107d4b95395dcd54d5c44f7509197b97114`.
- Builder commit: `f6048c580bd1330f89ffbe9a04b5f2890d81ab4a`.
- Evidence commit: `edd5df24c88865cc02bafc928a7bf4828776bda3`.
- Uncommitted files at handover: none. Before the handoff commit, the only existing dirty path
  was `do-work/CHECKPOINT.md`, changed by canonical `advance --checkpoint`; it is included here.

Remaining, in order:

1. Run `actions/review-work.md` in orchestrated mode against the fixed merge range. Read the
   current REQ and UR-130 input. Resolve actual findings; record Review, Lessons Learned and
   Orientation. Do not treat the moved release-note examples as executable code.
2. Reproduce `do-work/runs/work-2026-09-09-014954/req624-heavy-plan.json`, then run its sole
   selected lane with `skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json
   run-heavy-verification --manifest _dev/tests/heavy-lanes.json --lane staged-skills`.
   Record the real execution revision and lane result; it has not run for REQ-624 yet.
3. Finalize via canonical `advance` with an exact action-authored manifest, `supplied_commit`
   provenance and the merge hash above. **Omit release_manifest_path:** this history-only trim
   is release metadata under the existing guard and keeps 0.305.37. Do not create a synthetic
   release entry or alter release policy. Use a fresh completion timestamp and byte hashes.
4. Let canonical finalization close UR-130 and consolidate REQ-622, REQ-623 and REQ-624 with
   its input/assets. Then clean up the builder normally, consume run scratch after promoting
   the remaining evidence, refresh the checkpoint and run canonical cleanup. Stop at UR-130.

### Verified evidence

All artifacts below are committed under `do-work/runs/work-2026-09-09-014954/`:

- `REQ-624-handback.md`: full builder report, required lesson reads, file manifest and decisions.
- `req624-history-proof.json` and `verify-req624-history.py`: reproduce all 1,012 release blocks
  byte-for-byte from the fixed git baseline; 50 live, 962 archived, identical live copies,
  23 valid header links, five unchanged recent releases, archive export exclusion.
- `req624-merged-gate.json`: direct unpiped maintainer gate exited 0 in 108 seconds; 402 board
  and 815 CLI tests, every per-file budget below 30 seconds.
- `req624-test-gate.json`: focused shipped-reference and canonical green-gate records satisfied.
- `req624-qualification.json`: six relocated-marker warnings and four library-output findings
  in historical release prose. Every observed line was checked against pre-change bytes;
  the REQ records the qualification judgment. No source-code output was introduced.
- `req624-heavy-plan.json`: only staged-skills selected; still pending after review.
- `handoff-phase.json` and `handoff-recover.json`: exact lifecycle and ownership observations.

The two blank-at-EOF diff warnings preserve the original release separator bytes in the live
files. This is intentional, recorded as D-03; all other whitespace checks pass. New archive URLs
name committed local paths. They become available on remote main when these commits are pushed.

### Recovery and worktrees

Canonical `recover` succeeded: no unfinished finalizations. It reported
`RECOVERY-TAKEOVER-AVAILABLE` for REQ-624 and preserved the claim, with structural evidence:
`do-work/CHECKPOINT.md` line 21, writer
`t2s-Virtual-Machine.local:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`,
claimed at `2026-09-08T22:45:30Z`.

The exact offered reset command is:
`skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json recover --take-over REQ-624`.
It was **not run**: ordinary recover-claim returns a request to the queue and strips generated
workflow sections. This handoff retains the merged, green, review-ready state instead. The next
session has the user's authority to continue the existing claim and can inspect its phase with
`skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json advance REQ-624`.

Survey worktree verdicts:

- **ACTIVE integration checkout:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`,
  branch `main`; no uncommitted files after this commit.
- **ACTIVE retained builder:** `/private/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/do-work-ur130-3hy_uf65/worktree-agent-REQ-624-history-trim`,
  branch `worktree-agent-REQ-624-history-trim`, clean and merged. Its claim is still working,
  so it is not removable yet. No worktree was removed for this handoff.
- No foreign claims or unmerged builder branches were found. No subagent is still running.

### Completed work and parallelism

- REQ-622: completed and archived at `do-work/archive/REQ-622-correct-verified-prose-drift.md`,
  finalization commit `6a81154d007e4833c55474bb940ac5528a156b98`, release 0.305.36. Reconciled
  all 19 captured backlog entries: 16 corrected, 3 obsolete; runtime behavior unchanged.
- REQ-623: completed and archived at
  `do-work/archive/REQ-623-reduce-orchestration-instruction-duplication.md`, finalization
  commit `0ca8d37070bb9c2dbbfc881fda953e0cc6c95b7e`, release 0.305.37. Seven proven duplicate
  instruction cuts removed 8,369 source bytes / 1,272 words with unchanged owners and behavior.
- Those flat archive paths remain until REQ-624 closes UR-130; no shipped lesson links point
  to them, so the closure sweep found no link repair to add.
- Serial continuation: one unfinished request, no remaining builder dispatch. The original
  REQ-622 → REQ-623 → REQ-624 ordering is already encoded in depends_on. There is no reason
  to add fan-out to the review/finalization tail.
- The seven queued requests REQ-615–REQ-621 are unrelated and excluded by the explicit resume
  target. R4 stays held for the gate-policy decision; R5/R6 remain withdrawn. No new request,
  release-policy change, recipe-card expansion or KB promotion was authorized by this run.
- The approved report remains
  `ai-reports/2026-09-09_0013_do-work-improvement-review-rev-v2/index.html`.
