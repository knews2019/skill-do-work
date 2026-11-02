#!/usr/bin/env python3
"""Print the Testing section's heavy-plan lines from a plan-heavy-verification JSON file.
usage: heavy-lines.py <plan.json>"""
import json, sys
plan = json.load(open(sys.argv[1]))["heavy_verification"]
print(f"- Range: {plan['base_revision']}..{plan['target_revision']}")
for lane in plan["selected_lanes"]:
    reasons = lane["reasons"]
    reason = reasons[0] + (f" (and {len(reasons) - 1} more paths)" if len(reasons) > 1 else "")
    print(f"- {lane['lane_id']}: `{' '.join(lane['command_argv'])}` — {reason}")
