---
id: REQ-128
title: Secret rename quarantine survives re-inventory
status: completed
created_at: 2026-08-07T08:05:50Z
claimed_at: 2026-08-07T08:06:54Z
completed_at: 2026-08-07T08:29:17Z
commit: 7bb03d2
user_request: UR-028
addendum_to: REQ-121
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
related: [REQ-121]
effort_estimate: normal
route: B
kb_status: pending
write_set: [tools/checks/uncommitted-inventory.sh, actions/commit.md, actions/inspect.md, _dev/tests/contract-regressions.sh, actions/version.md, CHANGELOG.md]
---

# Secret Rename Quarantine Survives Re-inventory

## What

Correct REQ-121's uncommitted-inventory workflow so a secret rename remains excluded after an index reset degrades it into a deletion plus addition, and so an already-staged secret deletion can be committed without a failing restage.

## Context

This addendum corrects two accepted review findings against the inventory and its commit/inspect consumers. It keeps the existing `M`, `A`, `D`, `X`, and `XD` tags, forces rename detection independently of user Git configuration, preserves secret-derived destination quarantine across re-inventory, and rejects ambiguous secret-deletion/addition combinations before any addition is read or staged.

## Prior Implementation

REQ-121 introduced `tools/checks/uncommitted-inventory.sh` and rewired `actions/commit.md` and `actions/inspect.md` to consume it with manual fallbacks. The script classifies secret-shaped rename origins as `XD` and their destinations as `X`; its implementation commit is `167b0ae`. The current contract test proves only the initial staged-rename state, not the required reset-and-reinventory sequence or an already-staged deletion at commit Step 5.

## Requirements

- Explicitly enable Git rename detection in the inventory and both manual fallbacks, regardless of `status.renames` configuration.
- Preserve every destination excluded because of secret rename provenance across all later inventories in the same commit or inspect run.
- When only an ambiguous `XD` plus `A` state is available, stop before reading, associating, or staging any additions.
- At commit Step 5, accept an exact cached deletion for an `XD` path without calling `git add -u`; otherwise stage and verify a deletion-only cached change.
- Keep the existing inventory tag interface unchanged.
- Narrow the script's copy-protection claim; do not add content-based copy detection.
- Mirror all safety behavior in the documented manual fallbacks.

## Constraints

- Do not inspect secret contents or staged secret diffs.
- Do not alter or rewrite the sixteen commits already ahead of `origin/main`.
- Do not push.

## Red-Green Proof

**RED prompt/case:** Stage `git mv .env visible-config.txt`, reset the destination as commit Step 1 prescribes, and re-run inventory. The current output degrades from `XD .env` plus `X visible-config.txt` to `XD .env` plus `A visible-config.txt`, exposing the destination to Steps 2 and 3; Step 5 then calls `git add -u -- .env`, which exits 128 because the deletion is already staged.
**Why RED now:** The script relies on the current porcelain record for rename provenance, while the action prose discards that provenance between inventories and always restages `XD` paths.
**GREEN when:** Rename detection remains on despite user configuration, prior `X` destinations stay quarantined on every later inventory, untraceable `XD` plus `A` inventories fail closed before additions are consumed, and an exact already-staged deletion passes Step 5 without restaging.
**Validation:** User confirmed by accepting P1 and P2 and approving the implementation plan.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Buffer the inventory's NUL-delimited classifications so an `XD` plus otherwise-readable `A` state can be retagged fail-closed; force `--renames`; make both action consumers preserve a run-level union of `X` paths before every later read/association/stage; and teach commit Step 5 to verify cached deletion metadata before deciding whether `git add -u` is needed. Add the regression cases first and preserve the existing tag/exit interface.
- [x] **[APPLY]:** Added nine failing contract assertions/probes before implementation, then forced rename detection, buffered and quarantined ambiguous additions, retained a deterministic Git-private once-X-always-X set in both actions, and replaced unconditional deletion staging with cached metadata verification. Narrowed the copy claim without adding content inspection.
- [x] **[UNIFY]:** Reviewed `git diff --stat`, `git diff --check`, and every changed project file. Bash syntax and ShellCheck at warning severity pass; contract, hash-guard, update-script, queue-kanban Go, and queue verification checks pass. No debug artifacts or undeclared source changes remain.

## Triage

**Route: B** - Medium

**Reasoning:** The failure and acceptance behavior are specific, but the shared inventory primitive and both fallback consumers must be explored together before editing.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The inventory is a Bash 3.2-compatible, read-only reporter that already captures porcelain output through a temporary file and parses NUL-delimited rename origins. Its current direct emission prevents reconsidering an `A` row after an unexplained `XD` is found, while its `git status` call inherits `status.renames=false`. Both actions independently re-run the helper and filter only the current pass's `X` rows, so neither preserves earlier secret provenance. Commit Step 5 always calls `git add -u`, even when cached name/status already proves an exact deletion.

The minimum compatible change is to add explicit rename detection, buffer rows in a second NUL-safe temporary file, retag all `A` rows as `X` when the finished inventory contains `XD` plus `A`, and document an action-level once-`X`-always-`X` set that is overlaid before each consumer step. Cached deletion verification will use name/status only with rename detection disabled. The existing `C` parser remains defensive, but the header will stop claiming general copy protection because porcelain status does not discover ordinary copies.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `tools/checks/uncommitted-inventory.sh` (modify) — force rename detection and quarantine ambiguous additions without changing tags
- `actions/commit.md` (modify) — preserve quarantine across re-inventory and accept verified cached deletions
- `actions/inspect.md` (modify) — preserve the same quarantine in the read-only flow and fallback
- `_dev/tests/contract-regressions.sh` (modify) — prove the accepted rename and deletion cases
- `actions/version.md` (modify) — patch release version
- `CHANGELOG.md` (modify) — release note for the completed fix

**Files I will NOT touch:** `tools/checks/associate-files.sh` and content-based copy detection; the tag and association interfaces remain unchanged.

**Acceptance criteria (restated from REQ):**
- [x] Staged secret renames stay `XD`/`X` with rename detection disabled and after the destination reset.
- [x] Any unexplained `XD` plus `A` state quarantines additions before reads, association, or staging.
- [x] Both actions preserve every earlier `X` path for their full run and mirror the rule in manual fallback prose.
- [x] Exact already-staged secret deletions skip `git add -u`; unstaged deletions still stage and verify as deletion-only.
- [x] Existing tags, exit codes, ordinary additions without `XD`, and all prior regressions remain compatible.

## Decisions

- **D-01 — DECIDE & STATE:** Buffer tag/path pairs with NUL separators and retag every `A` as `X` only after the complete inventory proves `XD` plus `A`. This preserves the existing tags and path parser while making the ambiguous state fail closed.
- **D-02 — DECIDE & STATE:** Persist each action's quarantine in a deterministic file returned by `git rev-parse --git-path`, then re-derive that path in every command block and remove it at action exit. A shell variable or random temporary name cannot survive the action's separate invocations; Git-private scratch keeps it outside the working tree and worktree-safe.
- **D-03 — DECIDE & STATE:** Verify one exact cached `D` entry with `git diff --cached --name-status --no-renames -z`. This reads only path/status metadata and distinguishes an already-staged deletion from the unstaged case without displaying former contents.
- **D-04 — DECIDE & STATE:** Retain defensive parsing for any `C` record Git emits, but remove the general copy-protection claim. Discovering pure copies would require content-based work outside this accepted fix and would create a new secret-content inspection surface.

## Implementation Summary

**Files changed:**
- `tools/checks/uncommitted-inventory.sh` (modified)
- `actions/commit.md` (modified)
- `actions/inspect.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `actions/version.md` (modified)
- `CHANGELOG.md` (modified)

**What was done:** Inventory now overrides disabled rename detection, buffers classifications, and converts ambiguous additions beside a secret-shaped deletion to `X`. Commit and inspect keep excluded paths quarantined across re-inventory and mirror the rule in their manual fallbacks. Commit accepts an exact already-staged secret deletion and otherwise stages then verifies deletion-only metadata. Version 0.183.4 records the fix.

## Testing

**Tests run:** Bash syntax checks; ShellCheck at warning severity; `_dev/tests/contract-regressions.sh`; `_dev/tests/record-commit-hash-guards.sh`; `_dev/tests/update-script-behavior.sh`; `go test ./...` in `tools/queue-kanban`; `queue-kanban verify --repo-root .`; `git diff --check`.
**Result:** ✓ All passing.

**Red-green validation:**
- Contract inventory/action probes: ✗ nine expected failures before implementation → ✓ full contract suite after implementation.
- Rename detection disabled: destination was `A` before → `XD .env` plus `X visible-config.txt` after.
- Reset and re-inventory: destination was readable `A` before → quarantined `X` after.
- Cached deletion: prescribed `git add -u -- .env` failed before → exact cached deletion is accepted without restaging after.
- Unstaged deletion: still stages with `git add -u` and verifies as one exact cached `D` entry.
- Ordinary addition without `XD`: remains `A`.

**New tests added:**
- Rename-disabled and reset/re-inventory secret-rename probes
- Ambiguous `XD` plus `A` quarantine and ordinary-addition isolation probes
- Already-staged and unstaged secret-deletion metadata probes
- Action-contract checks for quarantine persistence, manual fallback parity, and both Step 5 branches

## Qualification

Passed — six project files verified, all seven requirements traced, P-A-U confirmed, and Scope matches the Implementation Summary. The script changes are substantive and preserve their read-only data flow; both action changes consume the quarantined inventory before content, association, or staging operations.

## Review

**Overall: 100%** | 2026-08-07T08:29:17Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — both accepted findings are covered end to end, and no quarantined destination can reach content inspection, REQ association, or staging through the prescribed flows.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Buffering the complete NUL-delimited inventory made the ambiguous `XD` plus `A` rule fail closed without changing the public tag or exit-code interface.
- Cached name/status metadata was enough to distinguish an exact already-staged deletion without reading secret content.

**What didn't:**
- Relying on Git's current rename record loses provenance after an index reset; later inventories need both an ambiguity rule and a run-level quarantine.
- The first action draft assumed a shell variable could survive separate prescribed command blocks; re-deriving a deterministic Git-private path is required.

**Worth knowing:** `git add -u -- <path>` can reject a rename source whose deletion is already present in the index, so deletion-only workflows must verify cached metadata before deciding whether to stage.

## Orientation

Secret rename handling now remains fail-closed across commit and inspect re-inventories, while verified deletion-only commits work whether the deletion was already staged or still unstaged.
