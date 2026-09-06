---
id: REQ-599
status: claimed
domain: backend
created_at: 2026-09-06T06:31:11Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
route: A
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-06T07:26:22Z
  basis:
    - Route A
    - 1-file write set plus its test
    - 3 acceptance criteria
maintenance: false
depends_on: []
related: [REQ-596]
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go]
title: 'Decide a REQ in-flight from the root being walked, not from a substring of its absolute path'
claimed_at: 2026-09-06T07:26:22Z
---

# Decide In-Flight-ness From the Walked Root, Not a Path Substring

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

`AssociateProjectPaths` decides whether a REQ is in-flight with
`active := strings.Contains(filepath.ToSlash(path), "/working/")` on the **absolute** walked path
(`internal/corehelpers/inventory.go:281`). The loop already knows which of the two roots it is walking;
it should ask that instead.

## Why

A repository checked out anywhere beneath a directory named `working` — a common name for a scratch or
workspace directory — makes every archived REQ satisfy the test. The terminal-success status filter on
the next line is then skipped, so a blocked, cancelled or abandoned archived REQ claims paths it should
not, and the wrong REQ is named as a file's owner during a commit.

## Context

Found during the independent three-lens review of REQ-596, which corrected the guide's description of
this rule. The prose is right and the tool is wrong, which is why it is its own request rather than part
of that one. One reviewer raised it; the synthesizer reproduced it.

## Detailed Requirements

- Decide in-flight-ness from the root the walk is currently in, not from a substring of the path.
- The observable rule must not change for a repository checked out anywhere else: a `working/` REQ
  counts whatever its status says, an `archive/` REQ only on a terminal-success alias.
- Add a test that fails on the current code: a fixture repository whose absolute path contains a
  `working` component, holding an archived REQ with a non-terminal status that claims a path. Today it
  wins the path; after the change it does not.

## Constraints

- One file plus its test. No prose change: the guide already describes the intended rule correctly.
- The package's existing inventory and association tests must stay green unchanged.

## Red-Green Proof

**RED case:** a checkout under a directory named `working`, an archived REQ with `status: blocked`
naming a path in its Implementation Summary, and that path uncommitted.
**Why RED now:** the substring test marks the archived REQ active, the status filter is skipped, and the
blocked REQ is reported as the path's owner.
**GREEN when:** the same fixture reports the path unassociated, and the existing tests are unchanged.

## Open Questions

None.

## Triage

**Route: A** — Build directly.

**Reasoning:** The defect is one expression at a known line, the intended rule is already stated
correctly in the guide, and the reproduction is a fixture whose absolute path contains a `working`
component. There is nothing to discover: the walk loop already knows which of its two roots it is in,
and the fix is to ask that instead of the path. The test is the point as much as the fix — no test in the
package reaches the archived-with-non-terminal-status case today, which is why the substring could sit
there.

**Planning:** Skipped.

## Plan

**Planning not required** — Route A: one expression, one test that fails on the current code.

*Skipped by work action*
