---
id: REQ-427
title: 'Confirm the Go version floor for installing and updating do-work'
status: pending
created_at: 2026-08-30T17:40:00Z
status_changed_at: 2026-08-30T18:24:45Z
user_request: UR-081
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
addendum_to: REQ-407
review_generated: true
---

# Confirm the Go Version Floor for Installing and Updating Do-Work

## What
UR-081 specified Go 1.26.1+ as the prerequisite, and REQ-406 and REQ-407 implemented exactly that. Measurement during REQ-407's review found the module does not need it. This asks whether the floor should stay where you put it.

## Why This Is A Question And Not A Fix
The floor is not a builder's choice. `REQ-406` says "require Go 1.26.1+" and `REQ-407` says "document the Go 1.26.1+ prerequisite", both taken from UR-081. A reviewer raised the height of the floor as a defect and it was **refuted 3-0** on exactly that ground: implementing your stated requirement is not a bug. But the requirement now decides who can install do-work at all, so it is worth confirming rather than inheriting.

## Measurement
Taken on the merged tree, on a copy of the module with only the `go` directive in `go.mod` changed:

| `go` directive | Result with the local Go 1.24.7 toolchain |
|---|---|
| 1.26.1, 1.26.0, 1.25.0 | refuses — `go.mod requires go >= …` |
| **1.24.0** | builds |
| **1.23.0** | builds; `go vet ./...` exit 0; `go test -count=1 ./...` **all six packages pass** |

So no language or standard-library feature in `do-work-cli` requires 1.26.1. As it stands, anyone on Go 1.23, 1.24 or 1.25 cannot install or update do-work.

One thing that is genuinely 1.26: `skills/do-work-board/tools/queue-kanban/go.mod` declares `go 1.26`. The board tool is optional; the installer is not. Lowering the installer's floor would not lower the board's.

## Open Questions
- [x] Should the Go floor for installing and updating stay at 1.26.1, or drop to the lowest version the module actually builds and tests clean on? → Confirmed: lower to `1.23.0`
  **Recommended:** lower it to `1.23.0` in `skills/do-work/tools/do-work-cli/go.mod`, the launcher's `minimum_go_version`, `README.md` and `skills/do-work/actions/version.md`, so the prerequisite matches what the code needs. **Also:** keep 1.26.1 deliberately, if you want one toolchain across the whole suite and are content to exclude older installs; or pick an intermediate floor such as 1.24 to match the toolchain most CI images ship today.
  Value: a floor that matches the code stops excluding installs for no technical reason.
  Risk: low and reversible either way — it is four literals and a `go.mod` directive. Choosing to *keep* 1.26.1 costs nothing to implement, since that is what already ships.

  **Answered 2026-08-30** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommendation via `do-work clarify`: lower the installer and updater floor
  to Go 1.23.0, the lowest version on which the module was built, vetted, and tested successfully.
  The optional board tool's Go 1.26 requirement remains unchanged and is outside this REQ's scope.

## Notes
If the answer is to lower it, the change is mechanical and belongs in one small REQ: the `go` directive, `minimum_go_version` in `skills/do-work/tools/do-work-cli.sh`, the README prerequisite line, and the version-action prerequisite line, plus a check that the launcher's refusal message quotes the new number.

---
*Source: REQ-407 review (UR-081). Answer with `do-work clarify`.*
