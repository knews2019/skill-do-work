```
do-work run --fan-out 2
This command is sufficient; everything below it is context.
```

---

## Reference

- Queue: 16 pending, 0 pending-answers, 0 blocked, 0 in progress; no clarification command is needed.
- Initial wave: REQ-390 and REQ-406 are dependency-ready and safe to build concurrently. REQ-390 is isolated to queue-kanban Timeline code; REQ-406 establishes the new core CLI module. Integration remains serial.
- Critical path: REQ-406 → REQ-407 → REQ-408 → REQ-409 → REQ-410 → REQ-411 → REQ-412 → REQ-413 → REQ-414 → REQ-415 → REQ-416 → REQ-417 → REQ-418 → REQ-419 → REQ-420. These dependency gates are already encoded in the queued REQs; after the initial wave, the Go-platform batch advances serially.
- REQ-406 partial foundation: commit `329c55a9`; focused Go tests, vet, formatting, launcher behavior, and ShellCheck passed. The REQ remains pending because the full acceptance suite and maintainer gate have not run.
- Capture trail: UR-081 and REQ-406–420 were committed as `fc049247`; capture verification found no gaps and queue verification passed.
- Current checkout verdict: ACTIVE — `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` on `main`; clean at survey time, with no claimed REQs, builder worktrees, or unmerged builder branches.
- Uncommitted files at survey time: none. The ignored local build output `skills/do-work/tools/do-work-cli/do-work-cli` is reproducible and is not source state.
- Heads-up for the next session: inspect the preserved REQ-406 foundation instead of recreating it, then run its remaining fixture matrix and full gate before accepting the REQ.
