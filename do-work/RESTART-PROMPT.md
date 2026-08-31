```
do-work clarify
do-work run --fan-out 2
This command is sufficient; everything below it is context.
```

---

## Reference

- **Integration state:** `main` at `d090af93` before this handoff commit. REQ-426 is complete: builder commit `f06a60c1`, merge `73bd4a6f`, full merge range `29f07ae0..73bd4a6f`, release/bookkeeping `7e0112bc`, metadata `d090af93`, suite version 0.251.1.
- **Verification:** REQ-426's focused Go tests passed on the merged tree. The canonical `bash _dev/tests/maintainer-verify.sh` gate then passed twice with direct exit 0, including once after final version/changelog/archive bookkeeping; only the optional browser lane was skipped because no browser was available.
- **In-flight REQs:** none. `do-work/working/` is empty, `do-work/CHECKPOINT.md` has no interrupted claim, and there are no builder branches or worktrees.
- **Questions first:** REQ-427 is `pending-answers` because exact Go 1.23 and 1.24 toolchains fail while exact Go 1.25 passes all 16 core packages. Recommend 1.25.0; Go 1.23 requires a substantive traversal-resistant filesystem backport. REQ-436 is `pending-answers` for consent to audit two remaining low-reach special-mode publication paths. `do-work clarify` is sufficient to resolve both.
- **Queue state after clarification:** REQ-428 and REQ-429 are the first numeric dependency-ready pair. REQ-430 can pair with REQ-434 or REQ-435; REQ-431 → REQ-432 → REQ-433 must remain serial by their existing `depends_on` gates. REQ-411 remains gated by REQ-428, REQ-429, and REQ-435; REQ-412 remains gated by REQ-411 and REQ-433.
- **Critical path:** REQ-428 + REQ-429 + REQ-435 → REQ-411 → REQ-412 → REQ-413 → REQ-414 → REQ-415 → REQ-416 → REQ-417 → REQ-418 → REQ-419 → REQ-420. The cleanup chain REQ-430 → REQ-431 → REQ-432 → REQ-433 joins at REQ-412; REQ-434 is independent.
- **Parallelism:** fan-out 2 remains the useful bound. Every computed wave still uses isolated builder worktrees and serial integration. Existing dependency fields encode every must-not-run-together constraint; `write_set` remains display-only.
- **Worktree verdict:** ACTIVE — `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` on `main` is the sole checkout. No worktree is removable and no foreign claim exists.
- **Scratch lifecycle:** `do-work/runs/work-2026-08-31-103010/` was consumed and removed after its REQ-426 and REQ-427 evidence was promoted into the durable REQs. The handoff commit records the tracked deletions.
- **First-ten-minutes heads-up:** do not reinstate the old REQ-427 answer. Its earlier 1.23 result changed only `go.mod` while compiling with a newer toolchain; exact toolchains proved the current minimum is Go 1.25. The optional board and maintainer gate remain at their own Go 1.26 floor.
