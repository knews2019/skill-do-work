# REQ-430 Exploration

## Finding

The defect is confined to cleanup's operation-group planner/applier seam. `BuildPlan` already emits terminal member moves as `ARCHIVE-REQ-NNN` groups and the active UR input move as a separate `CLOSE-UR-NNN` group. `ApplyPlan` preflights every group independently and immediately admits each passing group to the one guarded transaction. There is no field or resolver for dependencies between operation groups, so refusing a dirty member group has no effect on its independently clean closure group.

The existing safety pattern is worth retaining:

- each group is its own preflight conflict domain (`existingDestination` plus `gittransaction.PreflightTargets`);
- directly refused groups produce `CLEANUP-GROUP-REFUSED` evidence;
- all admitted tracked groups are applied as one rollback-capable transaction and optional commit;
- unrelated groups continue when another group is refused;
- any mutation-time failure rolls the whole admitted union back, which already prevents a member failure from leaving closure applied.

There is no existing operation-group dependency representation elsewhere in cleanup. REQ-409's archived review explicitly identifies that missing dependency graph as the root cause. The shared `internal/dependencygraph` package models REQ frontmatter readiness and is not the right dependency layer for cleanup operations.

## Exact files

Production:

- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go`
  - extend `OperationGroup` with deterministic prerequisite group codes;
  - record every terminal member archive/consolidation group required by a resolved live UR;
  - attach the sorted, unique member group codes to that UR's `CLOSE-UR-NNN` group.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go`
  - retain independent direct preflight;
  - resolve prerequisites after direct preflight and before building `eligibleGroups`/`proposedChanges`;
  - fail closed for a missing, directly refused, transitively blocked, duplicate, or cyclic prerequisite;
  - add dependent refusal evidence naming the blocking member group while leaving unrelated directly eligible groups admitted.

Test:

- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go`
  - add the end-to-end dirty-member regression and inspect the planned closure prerequisite list in the same fixture, keeping the implementation to the REQ's estimated three-file write set.

No repository-model, transaction, command-handler, action, guide, or queue-file change is needed. `repositorymodel.RequestFile` already exposes the live location, normalized terminal status, request ID, and `user_request:` membership needed by the planner. `gittransaction` already supplies exact-target refusal and rollback behavior.

## Minimal approach

1. Add a clearly named field such as `RequiredGroupCodes []string` to `OperationGroup`.
2. While constructing member groups, collect the group code under the member's `user_request:` only when that member is terminally resolved and is actually planned to move in this run. This covers terminal members stranded in `queue/` or `working/` and loose archive members that must consolidate into `archive/UR-NNN/`; members already in their final UR archive folder require no operation and therefore no prerequisite.
3. Attach the collected codes to `CLOSE-UR-NNN`, sorted and deduplicated. Do not depend on all known UR members indiscriminately: the closure dependency is on the concrete moves required in this cleanup plan.
4. Refactor application into two logical stages without changing the transaction boundary:
   - direct group preflight/classification;
   - dependency eligibility resolution over the directly eligible group codes.
   A small three-state DFS or equivalent fixed-point resolver preserves deterministic behavior and makes chained dependencies safe for later cleanup fixes. Missing codes and cycles must block rather than silently admit a dependent.
5. Only dependency-eligible groups contribute proposed changes and transaction targets. When closure is blocked, emit a refusal finding whose evidence contains both `CLOSE-UR-NNN` and the blocking `ARCHIVE-REQ-NNN`. Keep the original member refusal too.
6. Leave execution ordering and the single transaction intact. If an admitted member operation fails during mutation, rollback already removes/reverts the closure and every other admitted mutation, satisfying the no-partial-closure boundary. Independent progress is preserved for the preflight-refusal case that exposed this bug.

This reusable group-code dependency seam is also the natural prerequisite for REQ-431, which depends on REQ-430 and needs conditional document-rewrite ownership. It should remain generic rather than hard-code `CLOSE-`/`ARCHIVE-` string checks in the applier.

## Acceptance checks

- A plan with two terminal live members under one otherwise closable UR gives `CLOSE-UR-NNN` prerequisites for both corresponding member move groups, in stable order.
- If one terminal member source is dirty, cleanup leaves that member and `do-work/user-requests/UR-NNN/input.md` active, does not create either refused destination, and reports evidence naming the blocking member group.
- In the same run, an unrelated clean terminal REQ still archives successfully.
- A clean fixture still archives all required members and closes the UR.
- Existing rollback coverage continues to prove that a mutation-time operation failure cannot publish closure.
- Dry-run and apply derive the same eligible set; refused closure does not appear as a proposed/applied change.
- Existing independent-group and exact-target tests remain green.

## Red/green and verification

Focused RED before implementation:

```sh
cd skills/do-work/tools/do-work-cli
go test ./internal/cleanup -run 'TestURClosureWaitsForRequiredMemberArchival' -count=1
```

Likely preflight gate after GREEN:

```sh
cd skills/do-work/tools/do-work-cli
go test ./internal/cleanup
```

Then run the module gates named by the CLI prime (`go vet ./...` and `go test -count=1 ./...`); the integrating owner can run the repository baseline when the fan-out workflow reaches integration. The focused cleanup package currently passes unchanged.
