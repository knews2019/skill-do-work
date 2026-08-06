---
id: REQ-120
title: Shipped files stop citing the export-ignored maintainer doc
status: pending
created_at: 2026-08-06T11:26:17Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---
*Source: `do-work/user-requests/UR-025/input.md` — Codex P1 finding on PR #133: "Remove maintainer-doc references from shipped files"*
