---
title: "Lessons from REQ-072: Go utility allocates REQ ids and version numbers and verifies release consistency"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-072-go-utility-allocates-req-ids-and-version.md]
related:
  - page: REQ-071-crash-recovery-must-respect-a-live-claim
    rel: complements
  - page: REQ-073-fan-out-dispatch-n-concurrent-builders-u
    rel: complements
  - page: REQ-081-next-version-ignores-flags-placed-after
    rel: complements
  - page: REQ-083-verify-reports-every-builder-worktree-as
    rel: complements
  - page: REQ-084-verify-s-queue-state-probe-misses-a-buil
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-072: Go utility allocates REQ ids and version numbers and verifies release consistency

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Add a subcommand to the shipped Go tool that allocates the next REQ number, allocates and writes the
next version, and verifies the release-ritual invariants that are currently checked by hand. It
**never** writes the changelog body.

## Solution summary

Added three subcommands to the shipped board tool. `next-req` prints the next free REQ number by running `actions/capture.md`'s own scan on top of the existing `enumerateDoWorkTree` walk (no second parser), consulting both filename and frontmatter id and tolerating gaps. `next-version <patch|minor|major>` rewrites the single `**Current version**: X.Y.Z` line, reads it back to confirm the write landed, prints the new number, and refuses without writing when the repo keeps no such line — the bump size is a required positional argument, never inferred. `verify` is read-only and runs eight probes: version/changelog agreement, changelog version ordering, entry-title reuse, duplicate REQ ids, a checkpoint naming a REQ that no longer exists, untrustworthy `claimed_at` stamps (stale past three hours, unparseable, future-dated past the shared 2-minute skew allowance, or absent), finished REQs stranded outside the archive, and `worktree-agent-*` leftovers plus any such worktree carrying uncommitted `do-work/` changes. Each finding carries a `Fixable` flag set only for what `do-work cleanup` can mechanically resolve, and the report ends with `N fixable: run do-work cleanup`; a probe that cannot run is reported as skipped rather than passing silently. Nothing in the tool writes `CHANGELOG.md`. Wired into the three existing actions as optional accelerators — each stating its missing-`go` fallback — plus a fourteenth forensics check and its guide row, `CLAUDE.md`'s ONE→two write-surface correction, the tool's prime file, and 21 new Go tests plus 8 contract-suite assertions.

## What worked

- **Running the tool against the live repo before believing the test suite.** All 21 tests were green *and* the version probe was wrong — the fixtures encoded the same misreading as the code, because I wrote both from the same sentence. Pointing the binary at this repo took ten seconds and exposed it immediately. A prose-derived probe should always be run against a known-healthy tree, since "fires on a healthy repo" is the failure a fixture built from the same premise cannot see.
- **Reusing `enumerateDoWorkTree` and `Board.Warnings` instead of re-deriving.** `next-req` inherits the walk's pruning of `deliverables/`, `runs/`, and `assets/` for free — and that pruning is load-bearing: a deliverable copy named `REQ-900-*.md` under `assets/` would otherwise push allocation 800 numbers into the future. The duplicate-id probe likewise inherits tree-section precedence. Neither behavior would have occurred to me to write from scratch.

## What didn't work

- **Taking requirement 3's "strictly greater" literally.** The rule is real (it's in `CLAUDE.md` § Before Every Commit) but it describes a *transient* state during release composition; the committed steady state is equality. Encoding a mid-process condition as a standing invariant produces a check that fails on every healthy repo — the most useless possible outcome, since it trains the user to ignore the tool. The generalizable lesson: when a requirement quotes a pre-commit rule, ask *at what moment is this true*, because a verifier runs at an arbitrary moment.
- **Assuming `actions/version.md` is a findable path.** It only exists at the repo root in a skill-development checkout. Ten minutes went into the design before noticing that a consumer install puts the skill under `.claude/skills/do-work/` while `do-work/` sits at the consumer root — so the tool's own repo-root resolution (which deliberately *skips* skill installs, per `walk.go`'s `isSkillInstallDirectory`) would never find it. Anchoring on the content marker rather than the path is what made the subcommand portable, and it was the REQ's own instruction all along.
- **Writing the decision down is not declaring the file.** D-03 recorded the intent to touch `docs/forensics-guide.md` and `## Scope` still didn't list it; only `tools/checks/scope-drift.sh` noticed. The Decisions section and the Scope declaration are read by different checks, and satisfying one does not satisfy the other.

## Worth knowing

- **Adding a forensics check that shells out to a compiler breaks that action's central promise.** `actions/forensics.md` opens with "never modifies files"; `go build` writes a binary. The carve-out is scoped (gitignored binary, inside the skill install, nothing in the project or `do-work/`), but any future check that runs a build needs the same treatment — or the promise quietly stops being true.
- **`CLAUDE.md` is export-ignored, so it cannot be the only home for a shipped contract.** The complete two-write-surface statement lives there *and* in `tools/queue-kanban/prime-do-kanban.md`'s Traps, because only the second one reaches a consumer. Anything a downstream reader needs has to live in a shipped file.
- **The version probe's direction is diagnostic, not decoration.** "Ahead" means a bump landed without its changelog entry (normal mid-release, a real finding afterwards); "behind" means an entry carries a version the version file never received. Collapsing them into one "mismatch" message would throw away the half of the signal that tells you which file to fix.

## Back-reference

See `do-work/archive/UR-013/REQ-072-allocator-and-verifier-subcommand.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5db22ea`.
