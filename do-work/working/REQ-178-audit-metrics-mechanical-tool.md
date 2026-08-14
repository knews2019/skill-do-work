---
id: REQ-178
title: Build the audit-metrics tool for mechanical audit measurement
status: claimed
created_at: 2026-08-13T22:35:10Z
claimed_at: 2026-08-14T09:19:43Z
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
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---
*Source: UR-040 — user follow-up during capture*
