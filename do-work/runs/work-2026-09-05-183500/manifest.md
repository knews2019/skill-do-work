# Run manifest — work-2026-09-05-183500

Orchestrator: main checkout at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
Mode: `do-work run REQ-590` (serial, explicit target)

## Builders

| REQ | Route | Operative name | Worktree | Handback |
| --- | --- | --- | --- | --- |
| REQ-590 | A | worktree-agent-REQ-590-cap-the-path-list | ../skill-do-work2-worktrees/worktree-agent-REQ-590-cap-the-path-list | REQ-590-handback.md |

## Notes

- `recover` refused with FINALIZATION-DISCOVERY-AMBIGUOUS over uncommitted `do-work/` paths belonging to sibling sessions (REQ-583, REQ-587, REQ-588, REQ-589 evidence and their run directories). `uncommitted-inventory` returned those paths at info severity only. None is owned by REQ-590 and none is inside its write set, so the run proceeds and stages only its own declared paths.
- Sibling sessions are active in this checkout; REQ-589 is changing the same board strip on the client side only, and this REQ touches no client file.
