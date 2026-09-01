---
id: REQ-178
title: Build the audit-metrics tool for mechanical audit measurement
status: completed
created_at: 2026-08-13T22:35:10Z
claimed_at: 2026-08-14T09:19:43Z
completed_at: 2026-08-14T09:49:31Z
commit: 1afe780
kb_status: promoted
kb_entry: REQ-178-build-the-audit-metrics-tool-for-mechani.md
user_request: UR-040
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
related: [REQ-176, REQ-177]
batch: maintainability-audit
write_set: [skills/do-work-toolbox/tools/audit-metrics/]
maintenance: false
---

# Build the Audit-Metrics Tool for Mechanical Audit Measurement

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read prime-shell-commands.md, general/coding-guardrails/testing crew members, and the queue-kanban convention sources (main.go dispatch comment, open_work.go renderer split, frontmatter_test.go table style, verify_test.go real-git fixture). Approach:
  1. Files: `go.mod` (module github.com/knews2019/skill-do-work/audit-metrics, go 1.26), `.gitignore` (`/audit-metrics`), `main.go` (dispatch + shared flag helpers), `distribution.go`(+test) (nearest-rank percentiles + band evaluation), `inventory.go`(+test) (`inventory` + `folders` subcommands over `git ls-files -z`), `churn.go`(+test) (`churn` + `hotspots` over `git log -M --name-status`), `prime-audit-metrics.md`.
  2. Subcommands: `inventory` | `folders` | `churn` | `hotspots` (D-01). Flags all two-word kebab: `--repo-root`, repeatable `--exclude-path`, `--top-count`, `--since-window` (churn/hotspots), band thresholds `--watch-lines/--flag-lines/--watch-words/--flag-words` (inventory), `--watch-files/--flag-files` (folders). Bands only from flags; `>` flags, `==` does not.
  3. Renderer functions take `io.Writer` + computed data; `runXCommand` wrappers own FlagSet + `os.Exit` via `exitOnLeftoverArguments` (mirrored from queue-kanban). All git via `os/exec` with `git -C <repoRoot>`, three-valued exit handling.
  4. TDD: RED = distribution-math test against a stubbed percentile function → capture failing `go test ./...`; GREEN = implement; then lock-in tests for band edges, exclude honoring, shallow reporting, and rename attribution on a real `t.TempDir()` git fixture (commit A → rename A→B + touch → assert A-era touches land on B).
  5. Verify: `go build`, `go vet ./...`, `gofmt -l` empty, `go test ./...` green, spot-check `inventory` + `churn` against this repo (12-month window, five ceremony excludes).
- [x] **[APPLY]:** Code written per plan, all inside `skills/do-work-toolbox/tools/audit-metrics/`: go.mod, .gitignore, main.go, distribution.go(+test), git_support.go, inventory.go(+test), churn.go(+test), prime-audit-metrics.md (25 lines). TDD honored: distribution-math test written first against a stub, RED captured (`nearestRankPercentile(...) = 0, want 50` etc.), then implemented to GREEN; lock-in tests added for band edges (== not flagged, > flagged), exclude honoring (inventory + churn), binary sniff, shallow-clone reporting (real `--depth 1` file:// clone fixture), rename attribution, and staged copy-migration attribution — all on real `git init` fixtures in `t.TempDir()` with pinned identity/default-branch. One requirement-driven addition beyond the plan (D-06): `-C --find-copies-harder` copy detection, because the skills/ restructure was a staged copy-then-delete that `-M` alone cannot see — requirement 4 names that restructure's history as the thing that must survive.
- [x] **[UNIFY]:** All files are new (untracked dir); `git ls-files --others --exclude-standard` lists exactly the 11 source files — the built binary is correctly ignored by the nested .gitignore, and `git diff --stat` is empty (no tracked file touched; do-work/ is git-excluded). Verified per file: main.go (dispatch comment matches the actual flag surface; every wrapper rejects leftovers and resolves --repo-root before computing), distribution.go (nearest-rank + band-edge rules match their tests), git_support.go (`git -C` on every call; non-zero exits reported, never folded), inventory.go / churn.go (renderers take io.Writer only; no writes anywhere in the tool), tests (no skips except git-absent, no debug prints), prime-audit-metrics.md (15-30 line budget: 25; sections Read first / Do not edit / Must build / Traps). Checks: `go build` OK, `go vet ./...` clean, `gofmt -l .` prints nothing, `go test ./...` green (uncached). Spot-checked against this repo: inventory totals (608 files / 108,612 lines), churn with the five ceremony excludes over 12 months — `skills/do-work/actions/work.md` = 214 commits, exactly matching a hand-run `git log --follow` count; hotspots and folders render with band flags as specified. No debug artifacts.

## What

A small Go tool, `skills/do-work-toolbox/tools/audit-metrics/`, that produces the maintainability audit's deterministic numbers mechanically — inventory, distributions, band flags, churn — so the action pastes tool output instead of prescribing fragile find/wc/awk pipelines to an LLM. Script what can be scripted; judgment stays in prose.

## Why (if provided)

User (verbatim): "since you also have a go tool, consider building tools for the audit, that will also output mechanicanically some flagged folders, files, etc... for the MVP it does not have too be too complex, but basically whatever we can script would be good to have it as script not as LLM call, becuase those are cheaper and more robust". This is CLAUDE.md's "Programs beat prose for anything mechanical" applied to the audit itself.

## Detailed Requirements

MVP scope — what wc/find/git can answer robustly; CCN and duplication stay with the external tools (lizard/jscpd) and their NOT-MEASURED path:

1. **Inventory:** tracked-file counts and line/word totals by extension, honoring an exclude list (flag or config; defaults per REQ-176 requirement 14).
2. **Distributions:** per metric — file lines, file words, folder file-counts — median / p90 / p95 / max, plus top-N largest offenders.
3. **Band flags:** apply WATCH/FLAG thresholds passed as flags; output the flagged folders and files mechanically (path, value, band). Bands are inputs, never hardcoded — calibration happens in the action's conversation.
4. **Churn:** `git log --since=<window>` aggregation with: shallow-clone detection (report it, never silently truncate), exclude patterns for release-ceremony files, and **rename normalization before any current-path filter** — detect renames (e.g. `git log -M --name-status`, or per-file `--follow`) and attribute pre-rename touches to the file's current path, so the 2026-08-08 skills/ restructure's history counts toward the live files instead of being discarded with the dead paths; only paths deleted outright are dropped. Top-N output.
5. **Hotspot join:** churn × size (size as the MVP complexity proxy), top-N.
6. **Output:** markdown tables suitable for pasting directly into the audit report; machine-readable TSV behind a flag if trivially cheap, else skip (YAGNI).
7. **Pattern match with queue-kanban:** vendored source, built on demand (`go build` then run), invoked by the action as an accelerator with the manual-fallback contract — if `go` is absent or the build fails, the action falls back to the manual commands in its reference file; the tool is never a dependency.
8. Go tests pin the contract: distribution math on a fixed fixture, band flagging edges (value == threshold is not flagged; > is), shallow-detection reporting, exclude-list honoring, and rename attribution (pre-rename touches count toward the current path, not a dead one). Focused lock-in tests, not smoke slop.

## Constraints

- Read-only: the tool writes nothing; it prints. (The action owns the `do-work/audits/` report write.)
- Naming for reach per coding-guardrails § 5 — `audit-metrics` and its flags need findable, two-word names; single-word subcommands are exempt by design.

## Dependencies

REQ-176's action consumes this tool (its Phase 0/1 steps invoke it with fallback); build this first or in the same wave.

## Builder Guidance

Certainty: Firm on the MVP boundary (items 1–5, nothing more); Exploratory on CLI shape (single run vs subcommands — pick what queue-kanban's conventions suggest). Keep it deliberately small; complexity added here is complexity the audit pays forever.

## Red-Green Proof
**RED prompt/case:** `ls skills/do-work-toolbox/tools/audit-metrics/` fails; a failing Go test for the distribution math can be written first in the new package.
**Why RED now:** No mechanical measurement exists; the draft spec prescribes hand-run shell pipelines for every number.
**GREEN when:** `go test ./...` passes in the tool directory; running the built tool against this repo prints inventory, distributions, band flags (given thresholds), and churn tables matching hand-computed spot checks.
**Validation:** User requested mid-capture; MVP boundary inferred from their words ("does not have to be too complex").

## Full Context

See `do-work/user-requests/UR-040/input.md` for complete verbatim input.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is firmly specified (MVP boundary, 8 numbered requirements) but the patterns to follow — queue-kanban's CLI/module/test conventions, the vendored build-on-demand contract — need discovery before building. Not Route C: one tool, one package, no architectural ambiguity.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Key findings (Explore agent, full report in session context):

- **Module conventions:** `module github.com/knews2019/skill-do-work/audit-metrics`, `go 1.26`, flat `package main`, one file per capability with a why-doc comment at top, hand-rolled subcommand switch (no CLI library), per-subcommand `flag.FlagSet`, `--repo-root DIR` universal flag, kebab-case two-word flags, `exitOnLeftoverArguments` rejection, renderer split from command wrapper (renderer takes `io.Writer` so output is assertable; wrapper owns flags + `os.Exit`).
- **Output:** plain text `key : value` blocks in queue-kanban; REQ-178 wants markdown tables — renderer-split still applies.
- **Tests:** table-driven `testCases := []struct{...}`, fixtures as real files via `t.TempDir()`; git-behavior tests build a REAL git repo fixture with `exec.Command("git", ...)` (lesson REQ-083: plain temp dirs make git probes silently skip).
- **Git from Go:** shell out via `os/exec`, always `git -C <repoRoot>` (lesson: HEAD-relative answers inside worktrees); treat exit codes three-valued — 0 / 1 = "no" / other (128) = git declining to answer, never folded into "no". Precedent: `model.go:1258` uses `--format=%cI`.
- **Binary hygiene:** nested one-line `.gitignore` (`/audit-metrics`) inside the tool dir, never commit build outputs (prime-kanban-board.md:17).
- **Invocation contract for actions (REQ-176 will consume):** `(cd <suite-root>/do-work-toolbox/tools/audit-metrics && go build -o audit-metrics .) 2>/dev/null && <suite-root>/do-work-toolbox/tools/audit-metrics/audit-metrics <sub> --repo-root <project-root>` — accelerator form with manual fallback (capture.md:78-84 pattern).
- **Ceremony exclude list (this repo):** VERSION, skills/do-work/VERSION, skills/do-work/actions/version.md, CHANGELOG.md, skills/do-work/CHANGELOG.md — ship as illustrative default behind a flag, condition stated ("files touched by the release ritual"), per prime-shell-commands.md:29 (closed enumerations go stale).
- **Versioning:** vendored tools have no independent version/changelog — folded into the skill release (prime-kanban-board.md:12).

*Generated by Explore agent*

## Scope

**Files I will touch (all new, all inside `skills/do-work-toolbox/tools/audit-metrics/`):**
- `go.mod` — module github.com/knews2019/skill-do-work/audit-metrics, go 1.26
- `.gitignore` — one line `/audit-metrics` (never commit the binary)
- `main.go` — subcommand switch + flag parsing (testable parse functions)
- `inventory.go` + `inventory_test.go` — tracked-file inventory + line/word counts + distributions (median/p90/p95/max) + band flagging
- `churn.go` + `churn_test.go` — git churn with shallow detection, ceremony excludes, rename normalization; hotspot join
- `prime-audit-metrics.md` — 15-30 line routing-index prime (queue-kanban's in-dir prime pattern)

Exact .go file split within the directory is the builder's choice (write_set is the directory); anything OUTSIDE the directory is out of scope.

**Files I will NOT touch:** queue-kanban (pattern source only), any action file (REQ-176's job), CHANGELOG/VERSION (integrator's, Step 9).

**Acceptance criteria (restated from REQ):**
- [ ] `go test ./...` passes in the tool directory; a distribution-math test was written first and shown RED (tdd: true)
- [ ] Inventory honors an exclude list; distributions report median/p90/p95/max plus top-N offenders
- [ ] Band flags are inputs (flags), never hardcoded; value == threshold is NOT flagged, > is
- [ ] Churn reports shallow clones rather than silently truncating; excludes ceremony files; attributes pre-rename touches to current paths; only outright-deleted paths drop
- [ ] Hotspot join = churn × size, top-N
- [ ] Output is pasteable markdown tables; tool writes nothing (prints only)
- [ ] Naming per coding-guardrails §5 (two-word reach names; single-word subcommands exempt by design)

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/tools/audit-metrics/go.mod` (new)
- `skills/do-work-toolbox/tools/audit-metrics/.gitignore` (new)
- `skills/do-work-toolbox/tools/audit-metrics/main.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/distribution.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/distribution_test.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/git_support.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/inventory.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/inventory_test.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/churn.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/churn_test.go` (new)
- `skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md` (new)

**What was done:** Built the vendored audit-metrics Go CLI (zero dependencies, queue-kanban conventions): four subcommands — `inventory`, `folders`, `churn`, `hotspots` — emitting pasteable markdown tables with flag-supplied WATCH/FLAG bands, exclude-prefix filtering, shallow-clone reporting, and rename+copy-normalized churn (`-M -C --find-copies-harder` with dead copy-source reassignment; verified to reproduce `git log --follow`'s 214-touch count for work.md across the 2026-08-08 restructure). 10 lock-in tests including real-git and real-shallow-clone fixtures.

## Decisions

- **D-01 — Four single-word subcommands: `inventory` | `folders` | `churn` | `hotspots`.** Folders stays separate from inventory rather than folded in: the file-band flags (`--watch-lines/--flag-lines/--watch-words/--flag-words`) and the folder-band flags (`--watch-files/--flag-files`) would otherwise share one FlagSet and invite passing a file threshold to a folder metric. Single-word subcommand names are exempt by design (coding-guardrails § 5); every flag is two-word kebab. DECIDE & STATE.
- **D-02 — `--exclude-path` is a repeatable plain prefix match** on the repo-relative slash path, default EMPTY. The caller owns the list (a built-in ceremony list would go stale — prime-shell-commands § Closed Enumerations Go Stale); prefix semantics cover both a file (`CHANGELOG.md`) and a tree (`skills/`). DECIDE & STATE.
- **D-03 — Binary handling: NUL-byte sniff in the first 8 KiB.** Binaries count as files in the extension table but contribute zero lines/words and are excluded from the distributions; the output states how many were excluded. Word counts on binaries are noise, not data. DECIDE & STATE (the REQ pre-approved the sniff).
- **D-04 — Tracked-but-unreadable files are skipped and counted**, with a visible WARNING line when any exist — not a fatal error (a file deleted from the worktree mid-work would otherwise kill the whole run) and not silence.
- **D-05 — "Deleted outright" is implemented as a final filter against `git ls-files`** rather than parsing D entries into a death registry: only paths that exist today can carry churn, which is the same predicate with far less machinery.
- **D-06 — Copy detection added: `git log -M -C --find-copies-harder`, with dead copy-sources reassigned to their surviving copy.** The spot-check exposed that the 2026-08-08 skills/ restructure was a *staged* migration (copy in REQ-139, original deleted at cutover), which `-M` alone cannot see — `skills/do-work/actions/work.md` measured 8 touches vs 214 from `git log --follow`. Requirement 4 names that restructure's history as the thing rename normalization must preserve, so this is the REQ's own intent, not scope creep. After the change the tool's count is exactly 214 (matches `--follow`). Cost: ~5s on this repo, accepted for a tool that runs per audit, not per keystroke. DECIDE & STATE.
- **D-07 — `--since-window` defaults to "12 months"** (passed verbatim to `git log --since`); the rendered header always names the window in use, so a defaulted run is never ambiguous. The TSV flag from requirement 6 was skipped: with band sections and multiple tables per subcommand it is not "trivially cheap", so YAGNI as the requirement itself instructs.
- **D-08 — `folders` counts direct children only, not recursive** — the audit asks "which folder is crowded", and recursive counts would charge every ancestor for its whole subtree, drowning the signal in `skills/`-level totals.
- **D-09 — `--repo-root` (default `.`) is canonicalized via `git rev-parse --show-toplevel`** before any measuring. Load-bearing, not cosmetic: `git log` prints toplevel-relative paths while `git ls-files` prints cwd-relative ones, so a root pointing inside the repo would make the churn join silently empty — the silently-wrong-answer class this repo's lessons exist to prevent. Also turns "not a git repo" into one clean reported error.

---
*Source: UR-040 — user follow-up during capture*

## Qualification

Passed — 11 files verified on disk (1,298 lines), 8/8 requirements traced (TSV skipped under requirement 6's own YAGNI clause, logged as D-07), P-A-U audit confirmed (3/3 ticked, [UNIFY] evidence in Decisions), no debug artifacts in diff. qualify.sh mechanical checks OK; its 5 unreferenced-file WARNs are all Step 6.3 exceptions (test files, same-package Go compilation units, in-dir prime doc). Orchestrator independently re-ran build/vet/gofmt/tests and hand-verified the churn spot-check (tool: 214 touches for work.md == `git log --follow` count).

## Testing

**Tests run:** `go test -count=1 ./...` (orchestrator re-run, uncached)
**Result:** ✓ All passing (10 tests, 2.1s)

**Red-green validation:**
- `TestNearestRankPercentile` (distribution_test.go): ✗ before implementation (7 subtests failing, verbatim RED output captured in builder hand-back) → ✓ after — anchors the REQ's Red-Green Proof (`ls` of the tool dir failed before; distribution-math test written first)
- `TestChurnRenameAttribution` + `TestChurnStagedCopyMigrationAttribution`: pin requirement 4's rename normalization (pre-rename touches attributed to current path; staged copy-then-delete migrations reassigned)
- `TestChurnShallowCloneReported`: pins shallow detection on a real `--depth 1` clone fixture
- `TestBandLabelForValueEdges`: pins value == threshold NOT flagged, > flagged
- `TestChurnExcludePathHonored` + exclude/binary-sniff/no-flag-no-band inventory tests

**New tests added:** distribution_test.go, inventory_test.go, churn_test.go (10 lock-ins, real-git + real-shallow-clone fixtures)

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-14T09:55:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):**
- `skills/do-work/tools/checks/scope-drift.sh:47` — the Step 5.5 mechanical check silently self-disabled for this REQ: the awk pattern requires the literal header `**Files I will touch:**`, but REQ-178's Scope header carried a parenthetical, so zero paths parsed AND the "header present but unparseable" FAIL guard (line 70, same literal) was bypassed — observed `SKIP ... exit 2` instead of a comparison or a FAIL. Any REQ phrasing the header with a parenthetical dodges the check; the script's own comment names "silently disable the check" as the defect this guard exists to prevent. Not this builder's scope drift (hand comparison ran clean; finding is in a pre-existing pipeline file, routed per Restatement Sweep rule 3) — gate: rule-change → follow-up REQ-179 (sweep).

**Minor findings:** 2 (report only)
- `prime-audit-metrics.md:9` — Read-first gloss said `-M --name-status` while the shipped command is `-M -C --find-copies-harder`; corrected by the integrator pre-commit (one-line gloss fix in this REQ's own new file).
- Environment, NOT attributable to REQ-178: `bash _dev/tests/contract-regressions.sh` exits 2 in this sandbox on the run-blocked-check process-tree probe (`_dev/tests/prescribed-shell-scripts-behavior.sh:132`). Integrator verified `git diff origin/main -- skills/ _dev/` is EMPTY — the failing surface is byte-identical to main, so the failure is environmental (sandbox process-tree semantics), pre-existing relative to this REQ. `shipped-package-reference-contract.sh` exits 0 here.

**Nit:** `folders` reuses `computeInventoryReport` and reads file contents it only needs paths from — negligible for a per-audit tool; deliberate reuse.

**Acceptance:** Pass — reviewer independently built the binary, ran all four subcommands with and without band flags; output contract held on every probe: bands only with flags, strict-greater edges live-confirmed (1756 == flag threshold stayed WATCH), churn top row 214 for work.md matching an independent `git log --follow` count, ceremony excludes honored, shallow warning correctly absent on the full clone, leftover-token and non-git-root error paths exit 2/1 cleanly. Read-only guarantee verified.

**Scope drift:** `scope-drift.sh` → SKIP exit 2 (the Important finding). Hand comparison: 3 file additions vs Scope, all inside the declared directory and its explicit builder's-choice clause; zero outside touches; zero declared-untouched. No drift.

**Restatement sweep:** No stale external restatements — new contracts' canonical home is the new files; quick-wins.md/tidy-repo.md heuristics predate and do not restate this tool's contract. One internal gloss inconsistency (fixed, above).

**TDD gate:** Credible — verbatim RED capture matches distribution_test.go:30's Fatalf format; GREEN reproduced uncached by both reviewer and orchestrator.

**Suggested testing:** re-run contract-regressions.sh on a non-sandboxed machine; exercise accelerator-with-fallback end-to-end when REQ-176 lands; someday pin the renamed-then-recreated-old-path edge.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Mirroring queue-kanban's conventions wholesale (renderer/io.Writer split, per-subcommand FlagSets, real-git fixtures in t.TempDir()) meant zero design churn; the real-repo spot-check during build caught the biggest correctness bug (staged-copy migration) before review.
**What didn't:** `-M` rename detection alone missed the 2026-08-08 skills/ restructure entirely (8 vs 214 touches) — it was a staged copy-then-delete, invisible to rename detection; only `-C --find-copies-harder` plus dead-copy-source reassignment reproduces `git log --follow`. Also: phrasing the Scope header with a parenthetical silently disabled scope-drift.sh (→ REQ-179).
**Worth knowing:** Churn numbers from this tool are only trustworthy because of copy detection — anyone replacing it with a plain `git log --name-only | sort | uniq -c` resurrects the dead-path split. Shallow clones are detected and reported, never silently truncated. The tool is a separate Go module — a repo-root `go build ./...` never reaches it.

## Orientation

Now you can measure a repo mechanically for the maintainability audit: `audit-metrics` (new vendored Go module in the do-work-toolbox skill, beside queue-kanban in shape) prints inventory, size distributions, WATCH/FLAG band tables, rename-normalized churn, and hotspot joins as pasteable markdown. [MAP CHANGED] — new module; REQ-176's action will consume it as an accelerator-with-fallback. In-dir `prime-audit-metrics.md` is its routing index; `_dev/primes/prime-shell-commands.md` spot-checked — still accurate, not stale.
