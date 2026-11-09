#!/usr/bin/env bash
# dw <args...>: run do-work-cli in JSON mode against the main repo; print a compact view and save raw output.
root=/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
out="$(dirname "$0")/last.json"
bash "$root/skills/do-work/tools/do-work-cli.sh" --repo-root "$root" --format json "$@" > "$out" 2> "$out.err"
status=$?
echo "exit=$status"
python3 - "$out" <<'PY'
import json,sys
raw=open(sys.argv[1]).read()
try: d=json.loads(raw)
except Exception as e:
    print("non-json:", raw[:1500]); sys.exit(0)
print("outcome:", d.get("outcome"))
for f in d.get("findings",[]):
    print(" -", f.get("code"), "|", (f.get("message") or "")[:400], "|", f.get("affected_paths") or "")
    for k in ("automation_stop_reason","next_step"):
        if f.get(k): print("    ", k, ":", str(f[k])[:300])
adv=d.get("advance") or {}
if adv:
    print("phase:", adv.get("phase"), "| kind:", adv.get("phase_kind"), "| missing:", adv.get("missing_evidence"), "| next:", adv.get("next_argv"))
    for g in adv.get("gate_records") or []:
        print(" gate:", g.get("gate_id"), g.get("state"), g.get("outcome"), "| out:", g.get("output_lines"))
        for f in g.get("findings") or []:
            print("    -", f.get("code"), "|", (f.get("message") or "")[:500], "|", f.get("affected_paths") or "")
            if f.get("automation_stop_reason"): print("       stop:", f["automation_stop_reason"][:300])
for k in ("queue_advance","finalizations","recovery","worktree","heavy_verification","changes"):
    if d.get(k): print(k+":", json.dumps(d[k])[:1800])
if "--full" in sys.argv: print(json.dumps(d)[:6000])
PY
[ -s "$out.err" ] && head -c 600 "$out.err"
exit 0
