# REQ-624 Builder Handback — Trim the Shipped Changelog to 50 Entries

Complete and committed. Both live changelogs contain exactly 50 entries, and all 1,012 historical release blocks retain their exact bytes and original order.

- Branch: `worktree-agent-REQ-624-history-trim`
- Commit: `f6048c580bd1330f89ffbe9a04b5f2890d81ab4a`
- Base: `43a64ffe82f4ea020c32c7ae54a44c94448cbf65`
- Worktree: `/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/do-work-ur130-3hy_uf65/worktree-agent-REQ-624-history-trim`
- Source working tree is clean. No worktree `do-work/` paths were read, written, staged, or committed. This handback is the sole main-tree write.

## Changed Paths

| Path | Action | Result |
|---|---|---|
| `CHANGELOG.md` | Modified | Retains the newest 50 entries, with canonical GitHub links to every dated archive in newest-to-oldest order. |
| `skills/do-work/CHANGELOG.md` | Modified | Exact byte copy of the root live changelog. |
| `CHANGELOG-2026-09-05-up-to-v0.303.6.md` | Added | Preserves the 610 moved entries verbatim, with older/current navigation. Naming uses the date and version of its newest release, matching existing archives. |
| `CHANGELOG-2026-07-13-up-to-v0.121.0.md` | Modified | Adds the newer archive link and removes the now-stale `0.121.1 onward` current-range label. Its release blocks are unchanged. |

No version file or release entry was added or changed. Version remains `0.305.37` under brief decision D-01, because this change only packages release metadata.

## P-A-U Evidence

- **PLAN:** Read the authoritative brief, repository instructions, required crew contracts, release prime, and required lessons. Inspected the archive conventions, release ownership rules, actual changelog consumers, and the 50-entry boundary.
- **APPLY:** Split at release `0.303.6`, copied the root live file to its installed mirror, and updated archive navigation. Changes are limited to the four paths above.
- **UNIFY:** Reviewed every changed path and header. Reconciled complete release blocks against the parent commit before and after committing. Verified archive packaging, recent output, version agreement, and focused existing release fixtures. The intended result holds: smaller shipped history, complete navigable history, unchanged release text and cadence.

## History Preservation Proof

| Logical history file | Entries | Newest → oldest |
|---|---:|---|
| `CHANGELOG.md` | 50 | 0.305.37 → 0.303.7 |
| `CHANGELOG-2026-09-05-up-to-v0.303.6.md` | 610 | 0.303.6 → 0.121.1 |
| `CHANGELOG-2026-07-13-up-to-v0.121.0.md` | 16 | 0.121.0 → 0.110.0 |
| `CHANGELOG-2026-07-07-up-to-v0.109.0.md` | 144 | 0.109.0 → 0.65.0 |
| `CHANGELOG-2026-04-13-up-to-v0.64.1.md` | 56 | 0.64.1 → 0.50.0 |
| `CHANGELOG-2026-04-07-up-to-v0.49.0.md` | 136 | 0.49.0 → 0.1.0 |

Before: 660 live + 352 archived = 1,012. After: 50 live + 962 archived = 1,012 unique versions. The installed mirror is the required duplicate and excluded from the logical count.

Disposable Python reconciliation split each file at `^## X.Y.Z`, preserving all bytes through the next release heading or EOF. The ordered `(version, full block bytes)` sequence is exactly equal before and after. All existing archive release blocks match; the three archives without navigation changes remain whole-file byte-identical.

SHA-256 of the ordered concatenated release blocks, before and after:
`b33f120eba66aafc9701a172f91e68f758b18030c0484defcc271133601094e3`.

## Tests and Results

- Exact history reconciliation, unique version count, 50 live entries, and live mirror equality: **PASS**, including the committed state against `HEAD^`.
- Every live/archive header link resolves to an existing source path. All five installed history links use `https://github.com/knews2019/skill-do-work/blob/main/`: **PASS**. HTTP publication of the new target follows parent integration/push.
- `git check-attr export-ignore` marks all five dated archives excluded and leaves both live files included: **PASS**. Actual `git archive` of the staged tree, restricted to all seven changelog paths, includes only the two live changelogs with exact working-tree bytes. `.gitattributes` is unchanged.
- Version action replay reads the first 80 lines, extracts the complete newest five release blocks, and reverses them. Output before/after is byte-identical: `0.305.33`, `0.305.34`, `0.305.35`, `0.305.36`, `0.305.37`. Root VERSION, installed VERSION, action version and newest entry all remain `0.305.37`.
- Existing board release/parser fixtures: **PASS**, exit 0 (`github.com/knews2019/skill-do-work/queue-kanban`, 0.711s), using the command below from `skills/do-work-board/tools/queue-kanban`.
- Plain `git diff --check` reports only the intentional preserved blank EOF separator in each live file. `git -c core.whitespace=-blank-at-eof diff --cached --check`: **PASS**, exit 0. No whitespace normalization was performed on history.

```sh
go test -run 'Test(ReadChangelogEntries|ReadCurrentVersion.*|VerifyFlagsVersion.*|VerifyFlagsDuplicateChangelogVersion|VerifyFlagsReusedChangelogTitle|VerifySkipsReleaseProbesOnAForeignChangelogFormat|ReleaseProbesRunInASuiteCheckoutAndAreNotApplicableElsewhere|GenuineReleaseProbeSkipsStaySkipped|VerifyNamesTheChangelogAsTheReleaseFindingSubject)$' -count=1 ./...
```

No permanent tests were added for this mechanical packaging change. The full shipped-reference contract and maintainer/heavy gates are deferred to the parent because they inspect repository `do-work/` targets that this builder must treat as absent; the parent confirmed ownership of those checks.

## Consumer Review and Integration Seams

`skills/do-work/actions/version.md` consumes only the newest five blocks within its initial 80-line read. The added header leaves all five complete blocks within that limit. Board release probes inspect the live changelog for newest-version agreement and duplicate current titles/versions; their existing fixtures pass. They already excluded older archives before this trim.

`skills/do-work-knowledge/prompts/architecture-decisions-log_create-or-expand.md` uses `CHANGELOG.md` as an agent-read fallback history spine. The live header now explicitly directs complete-history readers through every archive in order, preserving access to older entries without modifying the generic prompt. Other inspected changelog references describe release publication, project discovery or separate interview histories and require no change.

Parent integration must retain exact block bytes, maintain the two live files' byte equality, and keep this metadata-only change at version `0.305.37`. No other builder dependency or source seam remains.

## Required Reads

Read `CLAUDE.md`, `AGENTS.md`, and crew files `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `background-agents.md`, plus `anti-slop.md` before this artifact. Read `_dev/primes/prime-releases.md`, the full `_dev/primes/lessons-releases.md` including both `canonical-link-outlives-its-target` bullets, and all of `skills/do-work/tools/lessons-do-work-update.md`. All required lesson files were present. Canonical targets and predecessor links were verified rather than assuming URL shape proves validity.

## Decisions

- **D-02 — Preserve complete-history consumers through the live header.** Explicit archive traversal and direct links to every range preserve the ADR fallback's access to the full history within this packaging scope.
- **D-03 — Preserve the former inter-entry blank separator at live EOF.** Exact historical block bytes take priority over the two harmless blank-at-EOF warnings. The parent explicitly confirmed this choice; all other whitespace checks pass.

## Discovered Tasks

None.
