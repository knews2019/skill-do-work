---
id: UR-103
title: 'Review fixes: blocked probe cwd and interruptions, collision ambiguity, lifecycle postimage modes'
created_at: 2026-09-03T12:20:21Z
requests: [REQ-534, REQ-535, REQ-536]
word_count: 844
---

# Review Fixes: Blocked Probe Cwd and Interruptions, Collision Ambiguity, Lifecycle Postimage Modes

## Summary

Capture the five findings accepted by the `do-work validate-feedback` triage of 2026-09-03 (external review of the installed do-work-cli; the cited files are identical to this repo's source). Three findings share one file and one REQ; the other two get one REQ each.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-534 | Run `blocked_check` probes with the selected repository root as working directory; treat SIGINT/SIGHUP/SIGTERM during a probe as an interruption that stops selection and exits `128+signal`; install the signal handler before the child starts (Findings 1, 2, 3). |
| REQ-535 | `duplicateStatusesSatisfied` never satisfies dependents through a filename/frontmatter collision, even when every colliding file declares the same frontmatter id and is terminal-successful (Finding 4). |
| REQ-536 | Lifecycle postimages record the mode the file will really have after `ApplyPlan` instead of a hardcoded 0644, so crash recovery converges for non-0644 files (Finding 5). |

## Batch Constraints

- Capture only; implementation belongs to a later `do-work run`.
- Finding provenance (verbatim claim, severity, triage evidence, Surface-cost N/A) travels with each REQ per the Finding-Closure Ratchet; every REQ names its lock-in test.
- The three REQs are independent; no `depends_on` edges. REQ-535 is related to queued REQ-490 (wave depth from satisfied duplicate records), which consumes the duplicate-satisfied verdict but does not define it.
- Findings 6, 7, 8 of the same triage were Discuss/Push back and are not captured here.

## Full Verbatim Input

> ```
> do-work capture-request: the 5 accepted issues
> 
> Context: the five findings the `do-work validate-feedback` triage of 2026-09-03 accepted from an external review of the installed do-work-cli (sa2-sentence-aligner2 at 0.266.3; the five cited files are byte-identical to this repo's source at 0.266.7).
> 
> ### Finding 1: Run blocked probes from the selected repository  ·  P1
> - **Verbatim claim:** [P1] Run blocked probes from the selected repository — .claude/skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go:15-15. When next is invoked from outside the target repository using --repo-root, a relative blocked_check inherits the caller's working directory because this command has no Dir. Checks such as test -e config/ready can therefore produce the opposite result and incorrectly unblock or exclude work from the selected repository.
> - **Verdict:** Accept
> - **Evidence:** `blocked_probe_unix.go:15` builds `sh -c` with no `Dir`; `command_runtime.go:95-145` only records `--repo-root`, never chdirs. Both shipped callers pass `--repo-root`: `actions/work.md:143` and `do-work-board/justfile.template:32`. Reproduced with a scratch build: a REQ with `blocked_check: 'test -e config/ready'` probes exit 0 when run from the repo root and exit 1 (`BLOCKED-PROBE-FAILED`) when run from another directory with `--repo-root`.
> - **Surface-cost:** N/A (direct fix)
> - **Remedy:** Thread the repository root into `RunBlockedProbe` and set `command.Dir` before `Start`. The `ProbeRunner` signature changes, so update `Select` callers and the fake runner in tests.
> 
> ### Finding 2: Propagate probe interruptions out of next  ·  P1
> - **Verbatim claim:** [P1] Propagate probe interruptions out of next — .claude/skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go:49-50. When SIGINT, SIGHUP, or SIGTERM arrives during a blocked probe, this returns 128+signal as an ordinary probe failure with no error. evaluateCandidate consequently records BLOCKED-PROBE-FAILED and Select continues with other candidates and returns success, so Ctrl-C can be ignored and another request may still be selected.
> - **Verdict:** Accept
> - **Evidence:** `blocked_probe_unix.go:44-50` returns `128+signal, nil`; `next_selection.go:184-194` treats any non-zero as an ordinary `BLOCKED-PROBE-FAILED` and continues. Reproduced: SIGINT sent to `next` during a `sleep 20` probe gave exit code 0, outcome `success`, the blocked REQ excluded with probe exit 130, and the next pending REQ selected. The child tree was cleaned up correctly (no leftover `sleep`), so only the propagation is wrong. Toolbox commands already have the right shape: `report_image.go:159,229` use `mutationSignalContext` and set `ExitCodeOverride`.
> - **Surface-cost:** N/A (direct fix)
> - **Remedy:** Return a typed interruption error from `runOwnedProbe`; have `Select` stop evaluating, report a failure finding, and set the exit override to `128+signal` (or re-raise after cleanup), matching the toolbox pattern.
> 
> ### Finding 3: Install probe signal handling before launch  ·  P2
> - **Verbatim claim:** [P2] Install probe signal handling before launch — .claude/skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go:33-35. If an interrupt arrives after command.Start() but before this signal.Notify, the CLI takes the default signal action without terminating or reaping the newly isolated process group. Because the child is in a different process group, the probe and its descendants can continue after the parent exits.
> - **Verdict:** Accept (fold into Finding 2)
> - **Evidence:** `blocked_probe_unix.go:17` starts the child; `signal.Notify` is at line 34. Between them Go's default SIGINT action kills the parent, and the child sits in its own process group (line 16), so it survives.
> - **Surface-cost:** N/A
> - **Remedy:** Move `signal.Notify` and its `defer signal.Stop` above `command.Start()`.
> 
> ### Finding 4: Keep filename/frontmatter collisions ambiguous  ·  P1
> - **Verbatim claim:** [P1] Keep filename/frontmatter collisions ambiguous — .claude/skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go:156-158. A filename/frontmatter collision can have the same number of claim paths and frontmatter-indexed records—for example, REQ-020-first.md and REQ-021-second.md both declaring id: REQ-021. If both statuses are successful, this count comparison returns true, causing BuildGraph to suppress the ambiguity and mark dependents ready despite the mismatched identity.
> - **Verdict:** Accept
> - **Evidence:** `dependency_graph.go:156-158` compares only counts. Reproduced with a scratch test: `REQ-020-first.md` and `REQ-021-second.md` both declaring `id: REQ-021`, both `completed`, gave the dependent `IsReady=true`, `AmbiguousTargets=[]`. The existing lock-in `TestFilenameFrontmatterCollisionMakesDependencyAmbiguous` (`dependency_graph_test.go:100`) only covers the case where the two frontmatter ids differ, so this shape is untested. Duplicate readiness landed in `2c82ef12` with no recorded rationale for mismatched filenames.
> - **Surface-cost:** N/A (direct fix)
> - **Remedy:** In `duplicateStatusesSatisfied`, return false if any file in `RequestsByID[id]` has `FilenameID != TypedRecord.RequestID`. Add the scratch case as the lock-in test.
> 
> ### Finding 5: Record lifecycle postimage modes accurately  ·  P1
> - **Verbatim claim:** [P1] Record lifecycle postimage modes accurately — .claude/skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:146-150. When a lifecycle source, moved UR, checkpoint, or calibration file has a mode other than 0644, ApplyPlan preserves that existing mode while these journal postimages claim 0644. A crash after lifecycle mutation but before the phase write then leaves files matching neither the preimage nor postimage, so recovery and rollback both refuse and the finalization journal becomes stuck.
> - **Verdict:** Accept
> - **Evidence:** `state_apply.go:147-160` hardcodes `0o644` for every postimage. The mutation path preserves the existing mode: `atomic_file.go:55` chmods the temp file to the original's mode, and `state_apply.go:64,79` use `os.Rename`. Preimages capture the real mode (`finalization_prepare.go:222`), `equalImage` compares modes (`finalization_apply.go:419`), and `imageSetState` refuses when a file matches neither (`finalization_apply.go:311-312`). Verified by reading the chain; not reproduced live. Release postimages already solve this by copying the preimage mode (`finalization_prepare.go:236-238`).
> - **Surface-cost:** N/A (direct fix)
> - **Remedy:** Have `PlannedPostimages` stat the source and carry its real mode for retained and moved files, or fill `Mode` from the preimage in `finalization_prepare.go:95` when the path existed.
> ```

---
*Captured: 2026-09-03T12:20:21Z*
