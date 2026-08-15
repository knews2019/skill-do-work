---
id: REQ-188
title: Hotspot output silently drops unavailable tracked paths
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: backend
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-183, REQ-184, REQ-185, REQ-186, REQ-187]
batch: audit-findings-2026-08-14
write_set: [skills/do-work-toolbox/tools/audit-metrics/churn.go]
---

# Hotspot Output Silently Drops Unavailable Tracked Paths

## What

Keep unreadable or otherwise unavailable tracked paths visible in hotspot output as `NOT-MEASURED`, while preserving valid measured rows and warning that the ranking is incomplete.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`computeHotspotEntries` currently continues after a tracked-path measurement error and emits a successful numeric ranking without the path or a `NOT-MEASURED` marker. The omission makes unavailable evidence look like low risk and contradicts REQ-178's visible-warning decision.

## Context

- Audit priority: P3; impact 2; effort normal.
- Root-cause key: `hotspot-unavailable-evidence-visible`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 8.
- Reproduce: `cd skills/do-work-toolbox/tools/audit-metrics && probe_dir="$(mktemp -d /tmp/audit-metrics-unreadable.XXXXXX)" && git -C "$probe_dir" init -q --initial-branch=main && git -C "$probe_dir" config user.email probe@example.test && git -C "$probe_dir" config user.name Probe && ln -s missing-target "$probe_dir/unreadable.md" && git -C "$probe_dir" add unreadable.md && git -C "$probe_dir" commit -qm seed && go run . inventory --repo-root "$probe_dir" --top-count 20 && go run . hotspots --repo-root "$probe_dir" --since-window '10 years' --top-count 20`

## Detailed Requirements

- Carry a sorted collection of unavailable tracked paths through the churn/hotspot join instead of discarding measurement errors.
- Preserve every valid measured hotspot row.
- Render each unavailable path with current lines and score shown as `NOT-MEASURED`.
- Emit a visible warning that the numeric ranking is incomplete when unavailable paths exist.
- Add a real-Git fixture containing a tracked path missing from the worktree; prove valid rows survive while the unavailable path remains visible.

## Constraints

- Do not fail the entire hotspot command solely because one tracked path cannot be measured.
- Do not add a generic diagnostics subsystem.
- Lock-in limit: zero churn-bearing tracked paths silently absent from both measured and `NOT-MEASURED` output.

## Dependencies

None.

## Builder Guidance

Firm intent. One unavailable-path result field, one compact renderer section, and one real-Git test are the earned surface.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** In a temporary Git repository, commit a symlink to a missing target, then run `audit-metrics inventory` and `audit-metrics hotspots`.
**Why RED now:** Inventory exposes the unreadable tracked path, while hotspot output succeeds and silently omits it.
**GREEN when:** Valid hotspot rows still render, the missing-worktree path appears in a sorted `NOT-MEASURED` section with lines/score unavailable, and the output warns that ranking is incomplete.
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 08, labeled P3, impact 2, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 8 in the canonical audit for the complete batch constraints and validated evidence record.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
