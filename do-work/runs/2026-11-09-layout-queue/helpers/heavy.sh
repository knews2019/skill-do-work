#!/usr/bin/env bash
# heavy.sh <revision> <log> <lane>...: run heavy lanes at a revision in the detached drain-head checkout.
set -uo pipefail
main=/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
head="$main/.git/work-run-2026-11-09-layout/drain-head"
rev="$1"; log="$2"; shift 2
git -C "$head" checkout -q --detach "$rev" || exit 97
[ -f "$head/do-work/test-durations.tsv" ] || printf 'run_id\tfile\tseconds\tother_gate_processes\n' > "$head/do-work/test-durations.tsv"
lanes=()
for lane in "$@"; do lanes+=(--lane "$lane"); done
export QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
cd "$head" || exit 98
bash skills/do-work/tools/do-work-cli.sh --repo-root "$head" --format json run-heavy-verification --manifest _dev/tests/heavy-lanes.json "${lanes[@]}" > "$log" 2> "$log.err"
echo "HEAVY-EXIT=$?"
python3 - "$log" <<'PY'
import json,sys
raw=open(sys.argv[1]).read(); dec=json.JSONDecoder(); i=raw.find("{")
while i!=-1 and i<len(raw):
    try: obj,end=dec.raw_decode(raw,i)
    except Exception: i=raw.find("{",i+1); continue
    run=obj.get("heavy_verification_run")
    if run:
        print("execution_revision:", run.get("execution_revision") or run.get("revision"), "| keys:", list(run.keys()))
        for lane in run.get("lanes",[]):
            print(" lane", lane.get("lane_id"), "exit", lane.get("exit_status"), "wall", lane.get("wall_seconds"), "disp", lane.get("disposition"), lane.get("reason",""))
        print(" findings:", [f.get("code") for f in obj.get("findings",[])])
    i=raw.find("{",end)
PY
