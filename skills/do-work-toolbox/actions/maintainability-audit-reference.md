# Maintainability-Audit Reference

> **Companion file for `maintainability-audit.md`.** Contains the default bands, the calibration procedure, the default exclude list, the `audit-metrics` command reference, the manual fallback commands, the lock-in-limit guidance, the finding-class template, and the report template. Not invoked directly — loaded by the audit action at the steps that name each section.

---

## Default Bands

Defaults only — every value is recalibrated at the calibration gate before it is used. Metrics use two bands, never hard cliffs:

- **WATCH** — abnormal enough to list in the metrics appendix; becomes a finding only if the file is also a hotspot.
- **FLAG** — finding-eligible on its own.

Bands decide *eligibility* only. Severity comes from the 1-5 `Impact` score (churn × complexity × blast radius): a 700-line file with zero churn and no importers is a WATCH curiosity; a 350-line file that changes weekly is a P1.

| metric | WATCH > | FLAG > | measured by |
|---|---:|---:|---|
| FOLDER_FILES | 15 | 30 | `audit-metrics folders` (manual fallback below) |
| FILE_LINES | 300 (tests: 450) | 600 (tests: 900) | `audit-metrics inventory` (manual fallback below) |
| FILE_WORDS (markdown/prose) | no universal default — propose from the repo p95 | no universal default — propose from the repo p95 | `audit-metrics inventory` (manual fallback below) |
| FN_LINES | 40 | 80 | `lizard` (or a language-native tool) — no manual fallback: declined/absent ⇒ NOT-MEASURED |
| FN_CCN | 10 | 15 | `lizard` — no manual fallback: declined/absent ⇒ NOT-MEASURED |
| FILE_CCN | 80 | 150 | `lizard` — no manual fallback: declined/absent ⇒ NOT-MEASURED |
| FN_PER_FILE | 15 | 30 | `lizard` — no manual fallback: declined/absent ⇒ NOT-MEASURED |
| DUP_PCT | 3 | 8 | `jscpd` — **explicitly no fallback**: `wc`/`grep` cannot measure duplication, so a declined or absent jscpd means DUP_PCT is NOT-MEASURED, stated in the report — never approximated |

Run configuration values (not bands): `CHURN_WINDOW` defaults to `12 months` (also `audit-metrics`' built-in default), `HOTSPOT_COUNT` to 20, `FINDINGS_MAX` to 12 root-cause classes per run — the loop handles the rest. `EXCLUDE` defaults are below.

## Calibration Procedure

**FLAG is finalized at the gate as `max(absolute floor, repo p95)`.** Percentile alone fails on uniformly bad repos (p95 of a bad distribution blesses the badness); the floor alone fails on repos whose normal is legitimately different (a prose-heavy repo's median file dwarfs a Go repo's). Both, always. Show the measured median / p90 / p95 / max next to each proposed WATCH and FLAG value, and justify any deviation: "median file is 120 lines — the 600 floor flags nothing; p95 is 340, propose FLAG > 340."

**One bundled, editable proposal.** The gate is a single approval covering three things the user can edit line-by-line before approving:

1. **Installs** — which external tools, the exact commands, why each earns its place for this stack. **User-local locations only**: `pipx` or `uv tool` for Python tools, `npx` for Node tools, `go install` for Go tools. Never install into the repo tree, and never propose system package managers (no brew/apt) — a system-level install is not user-local and contradicts the read-only posture. A declined item means its metric is NOT-MEASURED, honestly.
2. **Bands** — the calibrated table, per the rule above.
3. **Scope** — directories in, directories out (the EXCLUDE list below plus anything repo-specific), and why.

Plus **at most three focused domain questions** ("which area hurts most when you touch it?", "is the test suite trusted or tolerated?"). Then stop and wait for the go-ahead. On approval: run the installs, record versions, and lock the agreed bands as this run's config.

**Repeat runs:** when the latest report in `do-work/audits/` records an agreed config, present a one-paragraph state delta instead of the full proposal and ask one question — "reuse last run's calibration, or recalibrate?" — then wait.

## Default EXCLUDE (do-work repos)

`audit-metrics --exclude-path` takes a repeatable, repo-relative path **PREFIX**; the tool's own default is **empty** — the caller owns the exclude list, and this section is that caller-side list. Every entry is confirmed at the calibration gate, not assumed.

For any repo running do-work, default-exclude:

- `do-work/` — the do-work record (queue, working, archive, user-requests) **including `do-work/audits/` — the audit must not measure its own output**
- `kb/` — knowledge stores
- `ai-reports/` — generated report output
- Changelog archives (e.g. `CHANGELOG.md` and any mirrored copies) — release history, not living code
- The generic layer: vendored code, generated files, migrations, lockfiles

For churn and hotspots, additionally exclude the **release-ceremony files named by the project's commit ritual** — files touched by every release commit (changelog, version markers) drown real churn; in one measured do-work repo they were 24% of all file-touches. The condition is "touched by the commit ritual", not a fixed list — read the project's ritual and name its files.

## audit-metrics Command Reference

The shipped tool lives at `tools/audit-metrics/` inside this skill (own Go module; see its `prime-audit-metrics.md`). Build on demand and run it from anywhere:

```bash
# Build once (cached after the first run), then invoke by absolute path.
(cd <suite-root>/do-work-toolbox/tools/audit-metrics && go build -o audit-metrics .) 2>/dev/null \
  && <suite-root>/do-work-toolbox/tools/audit-metrics/audit-metrics inventory --repo-root <project-root>
```

Four subcommands, all read-only (they print markdown tables and write nothing):

| subcommand | measures | flags beyond the common set |
|---|---|---|
| `inventory` | tracked files by extension, file-lines/file-words distributions (median/p90/p95/max), largest files | `--watch-lines N` `--flag-lines N` `--watch-words N` `--flag-words N` |
| `folders` | files per folder (direct children, not recursive) | `--watch-files N` `--flag-files N` |
| `churn` | file touches in the history window, rename/copy-normalized | `--since-window WINDOW` |
| `hotspots` | churn × current lines join — the Phase-2 reading list | `--since-window WINDOW` |

Common flags: `--repo-root DIR` (default `.`), `--exclude-path PREFIX` (repeatable, default none), `--top-count N` (default 10).

Caller guidance (these are the tool's contract, not optional habits):

- **Bands come only from flags.** No flag, no band section in the output. A value strictly **greater than** a threshold is flagged; **equal is not**. Never hardcode a threshold — calibration lives in the gate conversation, so pass the *agreed* values.
- **Shallow clones are reported by the tool**, never silently truncated: churn/hotspots output carries a warning line when the history is shallow. Treat that warning as "churn undercounts" and say so in the report.
- **The exclude list is yours.** The tool defaults to empty on purpose; pass every prefix from the agreed scope on every invocation.
- **Quote multi-word windows**: `--since-window '12 months'`.
- **Separate Go module** — build inside `tools/audit-metrics/`; a repo-root `go build ./...` never reaches it.

## Manual Fallback Commands

For when `go` is absent or the build fails. Each block is self-contained — run it as-is from the repo root, adjusting excludes and the window to the agreed config. These approximate; the report must say which path produced each number.

**Inventory fallback** — tracked files by line count, largest first (includes binaries; ignore obvious binary rows):

```bash
git ls-files -z -- . ':!do-work/' ':!kb/' ':!ai-reports/' \
  | xargs -0 wc -l 2>/dev/null | sort -rn | head -25
```

**Distribution fallback** — the calibration gate derives FLAG from repo p95, so the manual path must produce the same nearest-rank distributions the tool prints. One file per `wc` invocation (`-n1`) so no `total` rows pollute the numbers; run once for lines and once with `wc -w` for words:

```bash
git ls-files -z -- . ':!do-work/' ':!kb/' ':!ai-reports/' \
  | xargs -0 -n1 wc -l 2>/dev/null | awk '{print $1}' | sort -n \
  | awk '{ v[NR] = $1 } END { if (NR == 0) exit 1;
      printf "file lines: count %d  median %d  p90 %d  p95 %d  max %d\n", \
        NR, v[int((NR*50+99)/100)], v[int((NR*90+99)/100)], v[int((NR*95+99)/100)], v[NR] }'
```

```bash
git ls-files -z -- . ':!do-work/' ':!kb/' ':!ai-reports/' \
  | xargs -0 -n1 wc -w 2>/dev/null | awk '{print $1}' | sort -n \
  | awk '{ v[NR] = $1 } END { if (NR == 0) exit 1;
      printf "file words: count %d  median %d  p90 %d  p95 %d  max %d\n", \
        NR, v[int((NR*50+99)/100)], v[int((NR*90+99)/100)], v[int((NR*95+99)/100)], v[NR] }'
```

Binaries are included in these approximations (the tool's NUL-sniff exclusion has no cheap shell equivalent) — say so in the report when the fallback produced the numbers.

**Folder-shape fallback** — files per folder (direct children), crowded first (root-level files count under `.`):

```bash
git ls-files -- . ':!do-work/' ':!kb/' ':!ai-reports/' \
  | sed -e 's|/[^/]*$||' -e 's|^[^/]*$|.|' | sort | uniq -c | sort -rn | head -25
```

**Churn fallback — run the shallow check FIRST.** A shallow clone undercounts churn silently; the tool detects this for you, the manual path must check by hand:

```bash
git rev-parse --is-shallow-repository
```

If it prints `true`, deepen/unshallow the clone first, or mark churn NOT-MEASURED — never present shallow counts as full-window churn. Then:

```bash
git log -M --since='12 months' --format= --name-only -- . ':!do-work/' ':!kb/' ':!ai-reports/' \
  | sed '/^$/d' | sort | uniq -c | sort -rn | head -25
```

**Rename-fidelity caveat:** this manual pipeline counts touches against the path name *at the time of each commit* — a renamed or restructured file's history splits across dead paths (the tool normalizes renames and copies; this pipeline cannot). Treat manual churn as approximate, and **cross-check the resulting hotspot list against the do-work record**: archived REQ frontmatter `write_set:` entries (and body `## Scope` / `**Files I will touch:**` lists) are rename-immune, shallow-proof evidence of where work actually landed. Count archived REQs touching each candidate hotspot and note the caveat that legacy REQs may lack both fields, so the count is a floor.

**Duplication has no fallback.** There is no `wc`/`grep` approximation of duplication. jscpd declined or absent ⇒ `DUP_PCT: NOT-MEASURED`, stated in the report.

**Function length / complexity have no manual fallback** either — lizard (or the agreed language-native tool: `gocyclo`, `radon`, an eslint complexity rule) measures them or they are NOT-MEASURED.

## Lock-In Limits

A **lock-in limit** is a single number or zero-hit assertion — never a band — because a CI check needs a binary verdict. It is **pinned at the current worst observed value** ("no folder over today's max of 34"): green on day one, blocks regression immediately, and tightens as fixes land. Never set a lock-in limit at the aspirational band — that turns day one red and teaches everyone to ignore it.

**The audit only proposes limits; it never installs enforcement.** Each proposed limit ships inside its finding class with its **red case** — the Reproduce command returning hits — so the proposal already contains its own regression test. Accepted proposals flow through `do-work-toolbox validate-feedback`, then the capture handoff, and land as lock-in tests in the project's own test suite or CI, built as ordinary REQs. Nothing under the audit's control ever writes a CI config or a test file.

A lock-in limit is added surface, so its proposal must carry the surface-cost evidence inline (see the template's Surface-cost field): incident = the finding class and its instances; replay = the Reproduce command; cost judgment = one CI line versus the class recurring; regression test = the red case itself.

**Dimension-5 exception:** discoverability/context-coverage findings (missing or stale prime coverage) are judgment findings only — **no lock-in limit is ever proposed for this class**. REQ-touch counts serve as `Impact`-score evidence, the remedy is consolidation, and `Surface-cost: N/A` (a prime is consolidated context, not guard apparatus).

## Finding-Class Template

Never emit per-instance findings. Group instances under root-cause classes; emit at most FINDINGS_MAX classes. One class = one finding = one verdict downstream. Every field below feeds `do-work-toolbox validate-feedback` — claims only, remedy in its own field, zero imperatives.

`sweep_key` is **provenance only**: it travels in the finding block and, if the finding is captured, verbatim inside the resulting REQ body — it is not a review-flow sweep marker and changes no capture schema.

```markdown
### Finding N: [short title]  ·  P1|P2|P3
- sweep_key: [root-cause slug]
- Claim: [one refutable sentence, specifics verbatim]
- Label: VERIFIED | INFERRED
- impact: impact-critical | impact-user-visible | impact-rule-change | impact-negligible
- Impact: [1-5] — churn: [pasted number], complexity: [pasted number], blast radius: [pasted count]
- effort_estimate: effort-mechanical | effort-substantive
- Reproduce: [exact command, runnable as-is from the repo root]
- Challenges-decision: [path]   <!-- only when the finding collides with a documented decision -->

#### Instances
- [ ] [path] — [greppable pattern]  (line [N] — optional, valid only at the audited SHA)

- Remedy: [one line]
- Surface-cost: N/A | [incident/replay case; the surface that will remain; cost judgment; regression test]
- Lock-in limit: [single number or zero-hit assertion at the current worst, with its red case] | none — [one sentence why]
```

Field rules:

- **Label**: `VERIFIED` = tool output or grep hit demonstrates it; `INFERRED` = judgment after reading code. The validator scrutinizes INFERRED hardest — correct and expected.
- **Instances coordinates are grep-pattern-first**: `path — greppable pattern`. A line number is an optional annotation valid **only at the audited SHA** — path:line tables rot within hours on an active repo, grep patterns survive. (The checklist heading is `#### Instances`, deliberately below `###` so the pasted `## Findings` section stays one self-contained block.)
- **Severity comes from the 1-5 `Impact` score alone**: 4-5 → P1, 3 → P2, 1-2 → P3. The lowercase `impact:` line is a different thing: it is the REQ frontmatter token any follow-up carries (`../../do-work/actions/work-reference.md` → Request File Schema, via `../../do-work/actions/review-work.md` Step 10's two questions), it routes how the fix lands, and it never doubles as a severity.
- **Impact shows its three inputs** — the churn, complexity, and blast-radius numbers behind the 1-5 call, pasted, not narrated.
- **`effort_estimate`** uses the canonical enum `effort-mechanical | effort-substantive` only, and means SIZE — never derive it from the `impact:` token.
- **Ranking** (decided, not optional): order classes by the 1-5 `Impact` score **descending**; break ties by `effort_estimate` (`effort-mechanical` before `effort-substantive`). There is no Score field and no `Impact`-score/effort division.
- **Surface-cost**: `N/A` for deletions, simplifications, renames, and direct fixes. Any remedy that adds a guard, rule, or warning apparatus — and a lock-in limit is added surface — supplies all four elements inline: the concrete incident or replay case, **the surface that will remain**, the cost judgment (is the fix still cheaper than the surface it adds?), and the regression test that keeps it live.

## Report Template

Write `do-work/audits/audit-YYYY-MM-DD.md` (today's date; create `do-work/audits/` on first use):

```markdown
# Maintainability Audit — [project] — [YYYY-MM-DD]

## Header
- Audited commit: [full SHA]  ([clean | DIRTY — findings may not reproduce at this SHA])
- Agreed bands: [the locked table from the calibration gate]
- Tool availability: [tool → version | NOT-MEASURED (declined/absent)] per tool used or skipped
- Measurement path: [audit-metrics | manual fallback] per metric family

## Metrics Summary
[Per-metric distribution stats and band results — pasted tool output.]
[Delta table vs the previous report: better / worse / unchanged per metric — the loop's convergence signal.
 First run: state "no baseline" and skip deltas.]

## Findings — paste this section into: do-work-toolbox validate-feedback
[The numbered finding-class blocks, ranked by the 1-5 `Impact` score descending, `effort_estimate` tie-break.
 Fully self-contained — the validator may see this section alone.]

## Metrics Appendix
[WATCH items that did not become findings — visible, not actionable.]

## Pre-empted
[Candidates dropped because a documented decision, recorded lesson, or waiver covers them —
 one line each, with the path to what covered it.]

## NOT-MEASURED
[Each unmeasured metric with the declined or missing tool that closes it next run.]

## Next steps
Paste the Findings section above into: do-work-toolbox validate-feedback
Capture its accepted findings through the validator's capture handoff, then: do-work run
Re-audit after the fixes land — the repeat run asks "reuse or recalibrate?" and the delta table must move.
A class you have decided to live with goes in do-work/audits/waivers.md, not into another round of triage.
```

The waivers file `do-work/audits/waivers.md` is a flat list the user maintains: one line per waived `sweep_key` with the reason. The audit reads it every run and never re-flags a waived class.
