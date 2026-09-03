```
do-work roadmap
This command is sufficient; everything below it is context.
Re-plan the entire queue from the committed repository state before claiming or starting work.
```

---

## Reference

- REQ-461 was the only held request. It is archived as `completed-with-issues` at `do-work/archive/REQ-461-require-affirmative-project-owned-release-targets.md`.
- REQ-461 implementation commit: `ca5735402c873afdc58b4eb9ae8e4b61fe9af73b`. Lifecycle/bookkeeping archive commit: `3214eff8e1c7bf8a308a14f229a5e74320f5218d`.
- `bash _dev/tests/maintainer-verify.sh` exited 0 on the committed implementation state. Only archive bookkeeping, the pre-existing staged `CLAUDE.md` cleanup committed as `add9b04f4b20b55e3be52edb8e6430ca989a5b73`, and this handoff followed; no implementation source changed after the green run.
- Remaining REQ-461 issue: legacy recovery still treats arbitrary release-metadata descendants or manifests below a declared maintainer source root as owned through `pathWithinReleaseRoots`, without an explicit workspace-parent relationship. Narrow seeding to exact declared roots/mirrors, then propagate descendants only through proven workspace parents; add an arbitrary-descendant regression inside a maintainer root.
- No follow-up REQ or addendum was created for that issue. The next session should decide how it fits the re-planned queue.
- Queue state was not edited during wind-down. No claim remains in `do-work/working/`.
- Parallelism is intentionally undecided because the user requested a fresh queue re-plan before more work. Do not inherit an old fan-out assumption from this handoff.
- Worktree verdict: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` on `main` is ACTIVE as the primary checkout and was clean before this handoff file was written. The survey found no other worktrees and no uncommitted files.
