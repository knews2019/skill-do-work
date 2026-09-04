# Run Manifest — work-2026-09-04-182017

Run dir: do-work/runs/work-2026-09-04-182017/
Concurrency: 3 (wave size, `do-work run --fan-out 3`)
Status: in-progress

| Agent | Slice | Output file | Status |
|-------|-------|-------------|--------|
| 1 | REQ-559 retry a red repository gate once | REQ-559-handback.md | done — merged, reviewed Pass, archived, released 0.276.0 |
| 2 | REQ-560 hand-back and finalize check only own paths | REQ-560-handback.md | done — merged at cb3a831, gate green, review in flight |
| 3 | REQ-515 per-REQ recovery findings never stop the loop | REQ-515-handback.md | done — built at 7309b8a, awaiting its turn to merge |

Wave 2 (frozen, not started): REQ-547, REQ-564.

Hand-backs are consumed into each REQ's committed trail and then removed, per
the rule that a hand-back is never staged, committed, or merged. REQ-515's is
held outside the repository until its branch is merged.
