# Run Manifest — work-2026-10-20-200259

Run dir: do-work/runs/work-2026-10-20-200259/
Concurrency: `do-work run UR-160 --coordinate` (bare `--fan-out`, bound 2). Main session coordinates only; pre-dispatch, builder and integrator run as background agents. REQ-704 depends on REQ-703 (same Go file), so members run one at a time.
Status: in-flight
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------------------------------|------------|----------------|
| REQ-703 | (pending pre-dispatch) | (from brief) | do-work/runs/work-2026-10-20-200259/REQ-703-handback.md | claimed a8c0cf11; pre-dispatch | | | | |
| REQ-704 | (waits on REQ-703) | | do-work/runs/work-2026-10-20-200259/REQ-704-handback.md | queued; dependency-deferred | | | | |
