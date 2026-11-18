# Gallery integration

Integrate REQ-743 (report gallery) on main. This wave has one claimed request and no set-asides. The separate feedback feature was explicitly captured under UR-184 on the main tree; its work begins after REQ-744 and adds no scope to this merge.

Read REQ-743-preparation.md and the three implementation handbacks here. The implementation tips are 683dcf1d (CLI), 813c0517 (samples) and a8031e04 (workflow). Each is based on the same claim. The revision between their base and main only captures a separate request. Preserve every new queue record.

An additional isolated researcher is checking first-upgrade behavior. Its handback is REQ-743-upgrade-handback.md when available. Add any independently verified repair for that issue before final acceptance. Never bypass the updater's trust or ownership contract. Test and disclose any genuine first-upgrade limitation. Check no-op recovery as well as the four clean consumer cases. Generate and inspect docs/report-gallery/index.html on this repository too so the user can open the shipped result from the final response.

Use the current work action through canonical finalization Step 9. No recover, force, push, PR, or Step 10. Handle the production integration and repeat checks for actual fixes. T3 manages the builders' checkouts. Leave them in place. Record handback as REQ-743-integration-report.md with the finalization hash, terminal status, opened gallery path, verification, remaining limitations and actual wall time. REQ-744 waits for completion here.
