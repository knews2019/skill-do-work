# Report gallery preparation

Status: in-progress  
First wave: both requests finalized on local main  
Main branch and dispatch base: `main` at `f3840b087913d8501267dda1c1395580304f2ba4`  
Wave size: 2

| Request | Route | P50 active | Builder | Operative ref and worktree | Brief | Expected handback | State | Dispatch UTC | Takeover to finalization | Full gates run | Stall restarts |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| REQ-745 | C | 60 min, low | GPT-6.2-sol | t3/persistent-gallery-annotations (t3-354d0a63-342e63c86e1e705e) | [Brief](REQ-745-brief.md) | `REQ-745-handback.md` | Finalized | Worktree created 2026-11-18T21:18:42+02:00 | See REQ-745 integration receipt | 1 pre-build plus integration checks | 0 |
| REQ-746 | C | 50 min, low | GPT-6.2-sol | t3/copyable-gallery-design-starters (t3-6937a03b-eefe202ad9d08d74) | [Brief](REQ-746-brief.md) | `REQ-746-handback.md` | Finalized | Worktree created 2026-11-18T21:18:51+02:00 | See REQ-746 integration receipt | Reused pre-build plus integration checks | 0 |

Total estimated effort is 110 active minutes and the estimated critical path is 60 active minutes. REQ-744 is already complete. [Estimator output](estimate-critical-path.txt)

- The full maintainer gate ran once under `full-gate.lock` and passed on the first attempt at the dispatch base. It ran 428 board and 998 CLI tests in 158 seconds. REQ-746 reused that exact revision and argv. The lock was removed after the gate finished. [Status](canonical-gate-status.log) · [REQ-745 preflight](REQ-745-advance-preflight.json) · [REQ-746 preflight](REQ-746-advance-preflight.json)
- Both focused executions used the same eight-test gallery/catalog command and passed. The request-bound snapshots are byte-identical. The active baseline matches both snapshots. [Probe](gallery-focused-baseline.sh) · [REQ-745 snapshot](REQ-745-focused-baseline.json) · [REQ-746 snapshot](REQ-746-focused-baseline.json)
- Both feature RED probes reproduced at the base. The fresh 0.328.1 ordinary install has 10 shipped starters and one project starter that resolve. Curated rendering keeps the 10 worked cards and name-copy hooks, hides the shipped starters, and omits generated starters and feedback controls. The feedback Just entry point and generated design resolver are absent. [REQ-745 RED](REQ-745-red.log) · [REQ-746 RED](REQ-746-red.log) · [Probe](gallery-red-baseline.py)
- Prepare the worktrees from the recorded base. Bind each brief to the actual operative ref and path at dispatch. Both typed readbacks keep the live claims and show the matching green record. The phase classifier stays at preflight until implementation is summarized. [REQ-745 readback](REQ-745-advance-stop.json) · [REQ-746 readback](REQ-746-advance-stop.json)
- The handoff contract records template ownership, generated design and feedback paths, installed and curated routes, a common real-evidence candidate, and proposed serial integration order. [Shared contract](shared-gallery-contract.md) · [Coordinated preflight](coordinator-preflight.md)

The checkpoint still carries both original claims. The index is empty and the preparation artifacts are unstaged. Both requests finalized. Completed builder scratch is retained under the local Git directory in do-work-run-scratch/work-2026-11-18-191126. The integration receipts in this run are the durable references.
