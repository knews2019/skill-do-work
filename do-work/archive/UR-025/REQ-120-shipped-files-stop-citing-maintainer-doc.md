---
id: REQ-120
title: Shipped files stop citing the export-ignored maintainer doc
status: completed
created_at: 2026-08-06T11:26:17Z
claimed_at: 2026-08-06T11:34:00Z
completed_at: 2026-08-06T11:38:00Z
route: A
kb_status: pending
user_request: UR-025
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: false
depends_on: []
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/frontmatter_cli.go, tools/queue-kanban/prime-do-kanban.md]
maintenance: false
related: [REQ-116, REQ-118, REQ-119]
batch: codex-pr133-findings
---

# Shipped Files Stop Citing the Export-Ignored Maintainer Doc

## What

Four locations in shipped `tools/queue-kanban/` files cite this repo's own `CLAUDE.md`, which is `export-ignore`d — so the citation dangles in every consumer install and every clone that installs from the tarball. `_dev/tests/contract-regressions.sh`'s maintainer-document probe fails on all four. Restate each rule inline or point at a shipped home.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Four hits, three files, each re-grounded rather than deleted. Two have a shipped home to point at (`actions/board.md` for the degrade-gracefully behaviour, `actions/work-reference.md` → Timestamp rule for `now`); two are better stated as the obligation itself, with no pointer at all — the write-surface count and the goes-stale reasoning are each one clause. The prime lesson keeps its content and loses only the filename.
- [x] **[APPLY]:** Comment and prose edits only. Verified each file reaches zero mentions after the edit, then re-ran the probe.
- [x] **[UNIFY]:** `git diff --stat` → 3 files, comments/prose only, no statements changed. `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` green (4.4s) — expected, since no behaviour was touched. The probe's maintainer-document check now reports **no hits**; the remaining 7 FAIL lines are the update-script behaviour probes, which need a non-root runner and fail identically on `main`. No debug artifacts. Confirmed `maintainer_doc_mention_allowlist` was **not** touched.

## Why (if provided)

The rule and its enforcement already exist; only compliance is missing. Two of the four hits (`frontmatter_cli.go:25,33`) came from REQ-112 and are already on `main`, so the probe has been red there — which is why fixing only the two this branch added would leave the check failing and teach the next reader that this probe is expected to be red. The probe was deliberately written to flag **any** mention rather than citation idioms, because idiom-matching caught 0 of 8 real occurrences.

## Detailed Requirements

Fix all four hits the probe names:

- `tools/queue-kanban/model.go` — the `hasSchemaFieldContract` comment's "(CLAUDE.md → Closed Enumerations Go Stale)". Restate the reasoning in one clause; the point is that the exempt set is "whatever has no row", so an enumeration of exempt names goes stale when a row is added.
- `tools/queue-kanban/frontmatter_cli.go:~25` — "(CLAUDE.md → Shipped Tooling, "Toolchain exception to design for the floor")". The accelerator rule has a shipped home: `actions/board.md` documents the precondition-check-and-degrade behaviour, and `actions/work-reference.md` → Timestamp rule covers `now`. Point there, or restate the rule in one sentence.
- `tools/queue-kanban/frontmatter_cli.go:~33` — "CLAUDE.md requires that sentence to be amended in the same commit as any third [write surface]". Restate as the obligation itself without naming the file that carries it.
- `tools/queue-kanban/prime-do-kanban.md:68` — REQ-116's lesson bullet, which names "`CLAUDE.md`'s lock-step sentence". The lesson's *content* is worth keeping; only the pointer has to change. State it as "the maintainer doc's list of fields the board parses for display" or similar, without the filename.

Additional requirements:

- **Do not widen `maintainer_doc_mention_allowlist`.** That allowlist exists for files whose mention really is about a *consumer project's* `CLAUDE.md` (prime routing, tidy-repo, KB schema). None of these four qualify — they cite this repo's maintainer doc — so allowlisting them would silence the probe on exactly the case it was inverted to catch.
- **Preserve the knowledge, not the pointer.** Each of these comments encodes a real reason a reader needs; the failure mode to avoid is deleting the sentence instead of re-grounding it. `crew-members/anti-slop.md`'s "protect lessons learned" applies: strip the dangling reference, keep what it was explaining.
- Re-run the probe and confirm the maintainer-document check passes with zero hits. The 7 update-script failures are a separate, pre-existing, runner-dependent baseline (that suite needs a non-root runner) and are out of scope.

## Constraints

- Comments and prose only — no behaviour change, no test change. `go test ./...` must stay green but this REQ should not need a new test; the probe in `_dev/tests/contract-regressions.sh` *is* the test.
- `_dev/` is not shipped, so the probe script itself is not in scope to edit.

## Dependencies

None.

## Builder Guidance

**Firm on what to remove, open on wording.** Every replacement sentence has to carry the same reason the citation was carrying — the goal is a reader in a consumer install understanding the rule, not a passing grep.

## Red-Green Proof

**RED prompt/case:** `bash _dev/tests/contract-regressions.sh` and read the maintainer-document check — it currently prints `FAIL: shipped files must not mention the skill's own CLAUDE.md/AGENTS.md` followed by four paths (two in `frontmatter_cli.go`, one in `model.go`, one in `prime-do-kanban.md`).
**Why RED now:** `tools/` ships, `CLAUDE.md` is export-ignored, and four shipped comments name it.
**GREEN when:** that check reports no hits, `go test ./...` is still green, and each replaced comment still states the rule it used to cite.
**Validation:** Inferred during capture — derived from Codex's P1 finding on PR #133 and confirmed by running the probe. `tdd: false` because the proof is an existing repo check, not a runnable unit test to write first.

## Assets

None.

## Triage

**Route: A** - Simple

**Reasoning:** Four named locations, an existing check as the pass/fail oracle, and no behaviour to change.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

---

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/frontmatter_cli.go` (modified)
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** All four citations of the export-ignored maintainer doc are gone from shipped paths, with each rule restated rather than dropped. The floor-exception comment now points at `actions/board.md` and `actions/work-reference.md` → Timestamp rule, both of which ship. The write-surface comment states the amend-in-the-same-commit obligation without naming the file that carries it. `hasSchemaFieldContract`'s comment keeps the goes-stale reasoning as its own clause. The prime's REQ-116 lesson keeps its full content and refers to "the maintainer doc's lock-step rule" instead of the filename — and gained the actionable half, that a field joins the list when the board starts parsing it.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (the oracle for this REQ) and `go test ./...`
**Result:** ✓ maintainer-document check reports no hits; Go suite green

**Red-green validation:**
- `_dev/tests/contract-regressions.sh`, maintainer-document check: ✗ before — `FAIL: shipped files must not mention the skill's own CLAUDE.md/AGENTS.md`, listing `model.go:975`, `frontmatter_cli.go:25`, `frontmatter_cli.go:33`, `prime-do-kanban.md:68` → ✓ after, check absent from the FAIL list entirely
- Non-regression: `go test ./...` green before and after; this REQ changes no statements

**New tests added:** None — the repo check *is* the test, as the REQ specified.

*Verified by work action*

## Lessons Learned

**What worked:** Running the probe and reading **every** FAIL line rather than the tail. This finding was in the output of a suite I had already run and reported on earlier in the session; I had counted the update-script sub-suite's own summary ("7 failure(s)") and the trailing line, and missed a separate check failing above them. The reviewer found it in the same output I had.

**What didn't:** Two things. First, I wrote two of these citations while holding the rule that forbids them in context — the rule's own file was loaded, and the violation still went in, because "cite where this rule lives" is a strong writing instinct and the rule's whole point is that this particular file cannot be cited. Second, the earlier session's REQ-112 introduced the other two, meaning the probe had been red on `main` and nobody noticed — so a check that fails for an unrelated reason (here, 7 root-runner probes) provides cover for real failures in the same output.

**Worth knowing:** The probe deliberately flags **any** mention rather than matching citation idioms, because idiom-matching caught 0 of 8 real occurrences before it was inverted. The consequence for a writer: there is no phrasing that legitimately names this file from a shipped path — restate the rule or point at a shipped home. `maintainer_doc_mention_allowlist` is only for mentions of a *consumer project's* CLAUDE.md, so reaching for it to quiet a hit is nearly always the wrong fix.

## Orientation

Consumer installs no longer receive four comments pointing at a file that isn't in their tarball; each rule is now stated where a consumer can read it. Lives in the queue-kanban tool's comments plus its prime (`tools/queue-kanban/prime-do-kanban.md`). Not `[MAP CHANGED]` — prose only, no statements touched. This also takes the repo's contract-regression suite from 8 distinct failures to its true baseline of 7 (all root-runner artifacts).

## Review

**Overall: 96%** | 2026-08-06T11:37:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — the check reports no hits; Go suite unchanged and green.
**Follow-up REQs created:** None

**Minor:** two of the four fixes drop the pointer entirely rather than redirecting it, because the rules involved (the write-surface count, the enumeration-goes-stale reasoning) have no shipped home to point at — they live only in the maintainer doc. That is correct under the rule as written, but it means a consumer reading those comments gets the rule without its rationale's source. Promoting either into a shipped `decisions/` record is a larger call than this REQ.

**Restatement sweep:** ran, and it is the substance of this REQ rather than a side-check — the sweep target was "who else in a shipped path names the maintainer doc", answered by the probe itself with all four hits. Confirmed zero remaining mentions in each edited file, and confirmed the allowlist was not widened, which would have silenced the probe instead of satisfying it.

*Reviewed by review-work action*

---
*Source: `do-work/user-requests/UR-025/input.md` — Codex P1 finding on PR #133: "Remove maintainer-doc references from shipped files"*
