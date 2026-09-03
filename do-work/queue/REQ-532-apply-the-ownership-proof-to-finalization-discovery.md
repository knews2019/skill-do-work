---
id: REQ-532
title: '[impact-user-visible] Apply the ownership proof to finalization discovery'
status: pending
created_at: 2026-09-03T11:10:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-461]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-461, REQ-460]
sweep: true
sweep_key: closed-enumeration-for-a-condition
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery_test.go
---

# Apply the Ownership Proof to Finalization Discovery

## What

REQ-461 replaced the release planner's directory-name denylist with proof of project ownership. Finalization discovery carries a second copy of that denylist, and a narrower one. It admits the same paths REQ-461 just closed, and loses some the old release predicate caught.

## Instances

- `internal/finalization/finalization_discovery.go:1045` — `installedReleasePathForDiscovery` tests a prefix denylist: `.claude/skills/`, `.codex/skills/`, `node_modules/`, `vendor/`, `vendored/`, `generated/`, `.generated/`. Called at `:725` while enumerating configured release members.
- It admits the three REQ-461 fixed: `third_party/do-work/VERSION`, `dist/skills/do-work/VERSION`, and any arbitrarily named cache tree.
- Because it is **prefix**-anchored rather than per-segment, it also admits `packages/vendor/do-work/VERSION` — which the old release predicate *did* catch, so this copy was already the weaker of the two.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- Finding **F1** from REQ-461's implementation shape-grep, which was requested because that defect has a class. This is the `[family: closed-enumeration-for-a-condition]` family's third appearance (REQ-460, REQ-461, now here).
- **Half the fix is already present:** `enumerateTrackedReleasePaths` only walks tracked paths, which is REQ-461's Git-index half. The missing half is exactly the `SKILL.md` package marker.
- Supporting precedent in the same file: its `suite/modules.tsv`-at-HEAD test for "is this the maintainer repo" is already an affirmative repository-owned declaration rather than a name test.

## Detailed Requirements

- Discovery must classify a configured release member by proof of project ownership, using the same invariant REQ-461 established, not by directory name.
- Reuse REQ-461's proof rather than writing a second one; if the package-direction rule prevents importing it, say so explicitly and state where the shared predicate should live.
- The per-segment loss must close too: a marker or dependency directory at any depth must be recognized, not only as a leading prefix.
- Preserve every existing discovery outcome for genuinely project-owned members, and every typed finding and recovery argv.

## Constraints

- Do not reintroduce a directory-name test, larger or prefix-anchored; REQ-461's constraint binds here.
- Do not weaken `singleInsertion`'s bounding or any other existing guard in this file.

## Dependencies

Depends on REQ-461, which establishes the proof this applies.

## Red-Green Proof

**RED prompt/case:** Configure release members at `third_party/do-work/VERSION`, `dist/skills/do-work/VERSION`, an arbitrarily named cache tree, and `packages/vendor/do-work/VERSION`, then run discovery.
**Why RED now:** the prefix denylist admits all four — the first three because their spellings are unlisted, the fourth because `vendor/` is matched only as a leading prefix.
**GREEN when:** all four are excluded by the ownership proof rather than by name, a genuinely project-owned member is still discovered, and the existing typed findings are unchanged.

---
*Source: REQ-461 implementation shape-grep, finding F1.*
