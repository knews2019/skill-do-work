---
id: UR-057
title: Four defects in the do-work suite — upstream report from a consumer repo
created_at: 2026-08-19T13:42:45Z
requests: [REQ-279, REQ-280, REQ-281, REQ-282, REQ-283]
word_count: 1827
---

# Four Defects in the Do-Work Suite — Upstream Report From a Consumer Repo

## Summary

A consumer repo (`g1w-game-find-the-difference`, running v0.212.25 vendored at `.claude/skills/`) ran the suite's own actions against its live queue and reported four defects upstream. The user routed the report through `do-work-toolbox validate-feedback` with instructions to verify each claim against this repo before queueing.

All four verified as real. Five REQs were captured from them; one self-raised finding was dropped as a duplicate of queued REQ-272, and three sub-remedies inside otherwise-accepted findings were declined with reasons (see Batch Constraints).

## Extracted Requests

| Finding | Severity (as reported) | Verdict | REQ |
|---|---|---|---|
| D2 — `blanked-req-scan.sh` reports intact asset sidecars as destroyed content | high | Accept | REQ-279 |
| D1a — no `created_at <= claimed_at <= completed_at` read-side probe; Check 12's remedy predates the shipped repair | high | Accept | REQ-280 |
| D1b — nothing reconciles `calibration-log.tsv` against the frontmatter it was derived from | high | Accept | REQ-281 |
| D3 — `verify`'s release probes can never run in a consumer repo — and, found during capture, do not run in the maintainer repo either | medium (raised) | Accept (narrowed, widened) | REQ-282 |
| D4 — `verify` is not routable through `do-work-board` | low | Accept | REQ-283 |
| D1 remedy 2 — fabricated-stamp heuristic (`claimed_at == calculated_at` + `:00`) | high | Declined | — |
| D1 remedy 4 — "give the repair a home" | high | Already done | — |
| D2 sub-note — shebang-honoring note for `sh` invocation | high | Declined | — |
| D3 remedy A — discover a suite root and run release probes against it | medium | Declined | — |
| F5 (self-raised) — stale forensics "Check 11" citations in `repair-req-timestamps.sh` | low | Duplicate of queued REQ-272 | — |

## Batch Constraints

- **Verification context.** Every finding was checked against this repo's code and git history before capture. Two of the four reproduce here, in the maintainer repo, independently of the consumer's data — recorded in each REQ's Red-Green Proof.
- **Declined: the fabricated-stamp heuristic (D1 remedy 2).** `claimed_at == estimate.calculated_at` fires legitimately on 20 of this repo's 251 stamped REQs (8%), because Step 2's claim and Step 3.6's estimate can read the same instant. The `:00`-seconds co-signal narrows it to 4, but the output names a suspicion no path can act on mechanically. Under `crew-members/coding-guardrails.md` § 2's earned-defense rubric, the provable cases are already repaired by the two shipped repairers, and a warning apparatus on top is the surface the rubric exists to refuse. Do not build it as part of REQ-280 or REQ-281.
- **Already done: "give the repair a home" (D1 remedy 4).** `scripts/audit-archive-timestamps.sh` already repairs archive stamps from `git blame`, sharing the repairer's predicate by sourcing `scripts/repair-req-timestamps.sh`, deliberately never hook-wired. The residue is that no diagnostic points at it — that residue is REQ-280's second half, not a new script.
- **Declined: the `sh` shebang note (D2).** `blanked-req-scan.sh` is `#!/usr/bin/env bash`, mode 0755, and both documented call sites invoke it bare. Running it through `sh` by hand is a caller error, not a missing note.
- **Declined: suite-root discovery for release probes (D3 remedy A).** In a consumer repo the root `CHANGELOG.md` is the consumer's, so pointing the probes at a discovered suite root would check install integrity — an invariant `adr-019-four-skill-suite-contract.md` already assigns to the updater's all-or-recover contract. REQ-282 takes the reporting half only.
- **Write-set overlap.** REQ-280 and REQ-281 both write `skills/do-work-board/tools/queue-kanban/verify.go`; REQ-281 `depends_on` REQ-280 so they serialize. REQ-279 writes `skills/do-work/tools/checks/blanked-req-scan.sh`, which queued REQ-276 also touches (line 91's fence guard, a different concern).
- **Widened during capture: D3.** The upstream report assumed the release probes run correctly "in the skill's own development repo." Running `queue-kanban verify --repo-root "$PWD"` here during capture showed they do not: since the modular four-skill split the version file is `skills/do-work/actions/version.md` and no `actions/version.md` exists at the root, so all three probes are skipped in the maintainer repo too. REQ-282 therefore carries two halves — make them run here, and report not-applicable in a consumer install — with an explicit constraint that the second must not silence the first.
- **Ordering.** REQ-279 first: it is a one-line change, and until it lands, `cleanup` Pass 6 asks users to approve a destructive repair against undamaged files.
- **Path convention.** The report's paths are relative to each skill's own root (`do-work/`, `do-work-board/`), per the consuming repo's convention. They map onto `skills/do-work/…` and `skills/do-work-board/…` here.

## Full Verbatim Input

# Four defects in the do-work suite — upstream report

**For:** `knews2019/skill-do-work`
**Found in:** v0.212.25, installed into `g1w-game-find-the-difference` at `.claude/skills/`
**Date:** 2026-08-19
**Paths** are relative to each skill's own root (`do-work/`, `do-work-board/`), per the consuming repo's convention for upstream suggestions.

All four were found by running the suite's own actions against a live queue: `do-work cleanup`, and the board tool's `verify` subcommand as `actions/forensics.md` Check 14 documents it. Every claim below is reproducible from the commands given.

---

## D1 — A fabricated `claimed_at` is undetectable, unrepairable, and feeds the estimator

**Skill:** `do-work` (detection + repair), `do-work-board` (the one existing check)
**Severity:** high — it silently corrupts the estimator's calibration corpus

### What happens

`actions/work.md` Step 2 is the only writer of `claimed_at`, and it requires the current UTC instant. When a stamp is instead written after the fact from a guess, nothing in the suite notices unless the guess happens to land *after* `completed_at`.

In this repo, 26 archived REQs carry a `claimed_at` byte-identical to their `estimate.calculated_at`, with both rounded to `:00` seconds. Five of them share one instant (`2026-08-17T22:40:00Z`), which no sequence of real claims produces. Git proves two of them wrong outright:

| REQ | git truth | recorded |
| --- | --- | --- |
| REQ-1533 | claim commit `4d3c5f1a` 20:11:21, merge `09ee6176` 20:16:12 | `claimed_at 21:05:00Z` |
| REQ-1529 | claim commit `3c8b3e10` 19:37:09, merge `bfa63001` 19:50:08 | `claimed_at 20:25:00Z`, `completed_at 20:07:00Z` |

Only those two surface anywhere, because only they invert. The other 24 pass every check the suite has.

### Why it matters beyond the board

`actions/work.md` Step 7.5 appends `completed_at − claimed_at` to `do-work/calibration-log.tsv`, and `actions/estimate-reference.md` line 92 states the scoring table is fit from that corpus. So a fabricated claim stamp does not merely mis-render a card, it biases future estimates.

**25 of this repo's 54 calibration rows (46%) come from REQs in the fabricated set.** The read-time outlier rule (exclude spans over 4h or negative) drops the two inverted ones and keeps all 23 others.

Worse, the log is an independent third record that nothing reconciles against frontmatter. Comparing every row to its REQ:

```
rows compared: 54   agree: 39   DISAGREE: 15
```

Nine disagree materially, not by rounding: REQ-1475 (logged 70, frontmatter 10), REQ-1514 (95 / 33), REQ-1515 (85 / 41), REQ-1527 (68 / 24), REQ-1528 (72 / 26), REQ-1523 (28 / 25), REQ-1524 (105 / 91), REQ-1529 (95 / −18), REQ-1533 (20 / −47). For REQ-1529 and REQ-1533, neither record matches git.

### What the suite checks today

- `do-work/actions/forensics.md` Check 12 compares each `*_at` against **now**. A past-dated fabrication passes.
- `do-work-board/tools/queue-kanban/model.go:1261` is a single comparison, `completed_at < claimed_at`. That one line is the suite's entire time-consistency surface, shared by the board and by `verify`.
- Nothing checks `created_at ≤ claimed_at ≤ completed_at`.
- Nothing compares `calibration-log.tsv` against the frontmatter it was derived from.
- Nothing repairs. `actions/forensics.md` Core Rules say "Read-only… Report, don't fix," and `verify` marks the anomaly not `[fixable]`. Check 12's own remedy says to recover the true instant "from the REQ file's git history" and leaves it to the reader by hand.

### Suggested change

1. **Add an ordering check** to `actions/forensics.md` and to the tool's probe set: `created_at ≤ claimed_at ≤ completed_at`, reported per violated pair. Cheap, and it catches the class rather than the inverted subset.
2. **Add a fabrication heuristic.** `claimed_at == estimate.calculated_at` with `:00` seconds is a strong tell; so is one instant shared by three or more REQs. Report as a warning, not an error.
3. **Add a calibration-log reconciliation probe.** Recompute each row's `wall_minutes` from the REQ's frontmatter and report rows that disagree by more than a minute. A corpus nothing audits is a corpus that quietly decays.
4. **Give the repair a home.** Git holds the truth for every affected REQ here — the claim commit and the merge commit are both findable by REQ id. A `--fix-timestamps` mode on the scanner, or a `cleanup` pass following the Pass 6 consent pattern, would close the loop that Check 12 currently opens and abandons.

### Reproduce

```bash
# the 26 with claimed_at == calculated_at, both rounded to :00
cd do-work
for f in $(find archive queue working -name 'REQ-*.md'); do
  cl=$(grep -m1 '^claimed_at:' "$f" | sed 's/^claimed_at: *//')
  ca=$(grep -m1 '^  calculated_at:' "$f" | sed 's/.*calculated_at: *//')
  [ -n "$cl" ] && [ "$cl" = "$ca" ] && case "$cl" in *:00Z) echo "$(basename $f) $cl";; esac
done

# git truth for one of them
git log --format='%h %ad %s' --date=iso --diff-filter=A -- 'do-work/**/REQ-1533*'
```

---

## D2 — `blanked-req-scan.sh` reports intact asset sidecars as destroyed content

**Skill:** `do-work`
**File:** `tools/checks/blanked-req-scan.sh`, line 285
**Severity:** high — it invites a destructive repair against undamaged files

### What happens

The candidate enumeration is unscoped by depth:

```bash
find do-work/archive -type f \( -name 'REQ-*.md' -o -name 'UR-*.md' \) 2>/dev/null
```

`do-work/archive/UR-NNN/assets/` holds screenshot-description sidecars written by the capture flow, named after the REQ they illustrate. They are prose documents and correctly have no frontmatter. The scanner classifies "no parseable frontmatter" as destroyed content, finds no git version with frontmatter (there never was one), and reports permanent data loss.

In this repo it reports four, every one intact:

```
1545 bytes  # UR-334 Screenshot Descriptions            .../UR-334/assets/REQ-1340-screenshot-descriptions.md
1176 bytes  # UI-review screenshot descriptions         .../UR-347/assets/REQ-1386-ui-review-screenshot-descriptions.md
1089 bytes  # Designer workspace modal bug screenshots  .../UR-351/assets/REQ-1396-1397-designer-workspace-modal-bugs.md
1554 bytes  # REQ-1398 Screenshot Descriptions          .../UR-352/assets/REQ-1398-reject-dialog-screenshot-descriptions.md
```

Every finding in this repo is an `assets/` sidecar. There is no true positive to hide behind them.

### Why it matters

`actions/cleanup.md` Pass 6 runs this scanner, then asks the user to approve a restore that "overwrites each file's content with the recovered blob." The output carries a `git gc` urgency warning, which pressures toward yes. Here the recovered blob is nothing, and the scanner's own refusal to write empty content is the only thing standing between the user and four deleted files.

`actions/forensics.md` Check 13 reports the same findings as **Critical**.

### Suggested change

A `REQ-*.md` is a REQ record only where REQ records live. Restrict the enumeration to those locations rather than any depth under `archive/`:

```bash
find do-work/archive -maxdepth 2 -type f \( -name 'REQ-*.md' -o -name 'UR-*.md' \) 2>/dev/null
```

That covers `archive/REQ-*.md`, `archive/UR-NNN/REQ-*.md`, and `archive/legacy/REQ-*.md`, and excludes `archive/UR-NNN/assets/`. A `-not -path '*/assets/*'` filter also works but denylists one known subdir rather than stating where records live.

Second, separate the two states the scanner currently merges. "This file has no frontmatter and no version ever did" is not the same finding as "this file was truncated and git holds the content." The first is almost always a non-REQ file; only the second is data loss.

### Reproduce

```bash
bash .claude/skills/do-work/tools/checks/blanked-req-scan.sh   # exit 1, four findings
head -3 do-work/archive/UR-334/assets/REQ-1340-screenshot-descriptions.md
```

Note: the script needs `bash`. Running it with `sh` fails at line 113 on process substitution — worth a shebang-honoring note in Check 13 and Pass 6, which both show it invoked bare.

---

## D3 — `verify`'s release probes can never run in a consumer repo

**Skill:** `do-work-board`
**Files:** `tools/queue-kanban/release.go:26,52-57`, `tools/queue-kanban/verify.go:152`
**Severity:** medium — three invariants report as skipped rather than as inapplicable

### What happens

```go
const defaultVersionFileRelativePath = "actions/version.md"          // release.go:26
versionFilePath := resolveVersionFilePath(repoRoot, "")              // verify.go:152
```

The default resolves against `--repo-root`. That path is correct only when the repo root *is* the suite root, as in the skill's own development repo. In a consumer project the suite lives at `.claude/skills/do-work/actions/version.md`, so the probe looks in a place the file cannot be and reports:

```
- skipped version-vs-changelog probes: no version file readable at <project-root>/actions/version.md
```

`actions/forensics.md` Check 14 instructs consumers to run exactly this command with the project root, so in every consumer install these three probes are permanently off: version file against newest changelog entry, version numbers strictly increasing, titles not reused.

Its sibling subcommand has the escape hatch it lacks — `main.go` documents `next-version [--version-file PATH]`, while `verify` takes only `--repo-root`.

### Why it matters

Check 14 tells the reader that "a skipped probe is an unverified invariant, not a clean one," and to report it. That advice is right for a probe that could have run. Here it produces a permanent, unactionable warning in every consumer repo, which trains readers to ignore the skipped line — including on the day it means something.

### Suggested change

Make the release probes suite-aware rather than repo-root-relative, and distinguish **not applicable** from **skipped**:

- If a suite root is detectable (walk up from the executable, or accept `--suite-root`), run the probes against it.
- If the install is a consumer vendoring, report `- not applicable: release probes verify the suite's own release ritual`.

Consider whether they should run in a consumer repo at all. Comparing the vendored skill's version to the vendored skill's changelog detects a corrupted install, not a release-ritual violation the consumer can act on. Reporting that honestly is better than a path error either way.

### Reproduce

```bash
(cd .claude/skills/do-work-board/tools/queue-kanban && go build -o queue-kanban .) \
  && .claude/skills/do-work-board/tools/queue-kanban/queue-kanban verify --repo-root "$PWD"
```

---

## D4 — `verify` is unreachable through the skill that owns it

**Skill:** `do-work-board`
**Files:** `SKILL.md` routing table, `actions/board.md` mode table
**Severity:** low — discoverability, but it made a documented diagnostic unrunnable as typed

### What happens

`SKILL.md`'s routing table has two rows, `help` and `board`, and passes `serve`, `static`, `summary`, `cli`, `--port N`, `--out DIR` to the board action. `actions/board.md`'s mode table maps `cli` to `open-work`. Neither mentions `verify`, and `grep -n -i verify actions/board.md` returns nothing.

So `do-work-board verify` hits the rule "an unknown command prints board help and stops" — for a real subcommand of the tool this package owns, which the sibling core skill's `actions/forensics.md` Check 14 documents in full and depends on.

### Suggested change

Add a routing row and a board-action mode:

| Trigger | Mode | Effect |
| --- | --- | --- |
| `verify`, `check invariants`, `probes` | verify | Runs the tool's `verify` subcommand. Read-only. Exit 1 means findings, not an error. |

Keep Check 14 as the canonical description of the probe set so the two do not drift; the board action only needs the build-then-run contract it already uses for every other mode.

---

## Summary

| | Defect | Skill | Severity |
| --- | --- | --- | --- |
| D1 | Fabricated `claimed_at` undetectable, unrepairable, corrupts the estimator corpus | `do-work` + `do-work-board` | high |
| D2 | `blanked-req-scan.sh` reports intact asset sidecars as data loss | `do-work` | high |
| D3 | `verify`'s release probes never run in a consumer repo | `do-work-board` | medium |
| D4 | `verify` is not routable through `do-work-board` | `do-work-board` | low |

D2 is the one to fix first: it is a one-line change, and until it lands, Pass 6 asks users to approve a destructive repair against files that are not damaged.

---
*Captured: 2026-08-19T13:42:45Z*
