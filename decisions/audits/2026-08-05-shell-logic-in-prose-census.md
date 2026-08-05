# Shell-Logic-in-Prose Census — `actions/` and `prompts/`

**Date:** 2026-08-05 · **Scope:** every file in `actions/` (45 files) and `prompts/` (18 files) · **Read-only audit.**

Every claim below is **VERIFIED**: each row cites the file and line number(s) whose exact text was read, and every coverage
verdict was checked against the shipped source (`tools/checks/*.sh` headers, `tools/queue-kanban/main.go` subcommand
switch, `verify.go` probe constants, `model.go` field readers). No row is inferred from a filename or a summary.

---

## 1. What mechanical coverage exists today

The complete inventory of shipped executables, read from source:

| Executable | What it mechanizes | Verified at |
| --- | --- | --- |
| `tools/checks/archive-collision.sh` | REQ id already archived? | header L2–13; glob L25 |
| `tools/checks/preflight.sh` | work Step 5.75 — dirty tree (`-uall`), test baseline, deps; writes `baseline.json` | header L2–28 |
| `tools/checks/qualify.sh` | work Step 6.3 items 1, 4, and the grep half of 5; honors `DO_WORK_DIFF_RANGE` | header L2–19 |
| `tools/checks/scope-drift.sh` | Scope ↔ Implementation Summary set difference | header L2–13; `extract_section_paths` L23 |
| `tools/checks/record-commit-hash.sh` | the `commit:` frontmatter write-back + `--verify` | header L2–30 |
| `tools/checks/blanked-req-scan.sh` | 0-byte / unparseable REQ detection + `--restore` | header L2–30 |
| `queue-kanban summary` | column counts, completion anomalies, warnings | `main.go` L93–128 |
| `queue-kanban generate` / `serve` | static + live board (only surface that emits per-REQ ready/waiting) | `main.go` L131–155 |
| `queue-kanban next-req` | next free REQ **number** (not next REQ to work) | `main.go` L163–177 |
| `queue-kanban next-version` | rewrites one `**Current version**: X.Y.Z` line | `main.go` L23, L180–196 |
| `queue-kanban verify` | 13 probe categories | `verify.go` L16–31 |
| `queue-kanban now` | the UTC instant, one shape | `main.go` L28, L47–49 |

`verify`'s 13 probes, verbatim from `verify.go` L17–30: `version-changelog-mismatch`, `changelog-version-not-ahead`,
`reused-changelog-title`, `duplicate-req-id`, `merged-worktree-leftover`, `unmerged-worktree-leftover`,
`worktree-merge-state-undetermined`, `worktree-wrote-queue-state`, `worktree-committed-queue-state`,
`checkpoint-names-missing-req`, `claim-needs-attention`, `stranded-finished-req`,
`assigned-elsewhere-claimed-here`, `ur-archived-with-live-member`.

**Two structural facts that shape every verdict below.**

1. **`frontmatter.go` has no CLI surface.** The subcommand switch (`main.go` L60–76) exposes exactly
   `summary | generate | serve | next-req | next-version | verify | now`. `splitFrontmatter` /
   `parseFrontmatterFields` / `lenientFrontmatterFields` (`frontmatter.go` L28, L82, L118) are reachable **only**
   through those seven. There is no `queue-kanban frontmatter get <file> <field>`. So *no prose step can call
   `frontmatter.go`*, and Section 3's flag list is, necessarily, every frontmatter read in the two directories.
2. **Three independent frontmatter parsers already ship.** Go (`frontmatter.go` L28), awk inside
   `record-commit-hash.sh` (`frontmatter_line_for` L108–121), and a third in `blanked-req-scan.sh`
   (`has_parseable_frontmatter` L88). Every prose read is a fourth-and-onward reimplementation.

**`prompts/` is almost entirely clean.** 17 of 18 files contain zero shell commands, zero frontmatter reads, and zero
output-parsing steps (verified: `grep -c '^```\(bash\|sh\|shell\|console\)'` returns 0 for all 18, and the inline
command grep returns 0 for 17). The single exception is `architecture-decisions-log_create-or-expand.md`; its rows are
in the table. `prompts/README.md` is an index table only.

---

## 2. Census table

Coverage is **FULL** (a shipped executable does this step's mechanical work), **PARTIAL** (the logic exists in shipped
source but is not reachable as a command, or only the violation-detection half exists), or **NONE**.

### `actions/work.md` — the highest-frequency file in the skill

| Step | Mechanic | Coverage | Candidate script |
| --- | --- | --- | --- |
| Step 1, L127 | Glob `do-work/queue/REQ-*.md`, sort by number, read each frontmatter `status` | **PARTIAL** — `model.go` `normalizeStatus` L718 + bucketing L1006–1008 do this, but no subcommand emits a per-REQ selection; `summary` prints counts only | `queue-kanban next-work` |
| Step 1, L129–166 | Write `blocked_check` to a scratch file; run it under a 30 s bound via `timeout`/`gtimeout`/a hand-rolled `kill -0` poll loop; read `probe_exit`; fail closed on non-zero/124/unlaunchable | **NONE** | `tools/checks/blocked-probe.sh` |
| Step 1, L168 | On exit 0: set `status: pending`, stamp `status_changed_at`, **remove** `blocked_by`+`blocked_at`, keep `blocked_check`, append a `## Blocked` history line | **NONE** | `tools/checks/unblock-req.sh` |
| Step 1, L170–172 | Dependency-ready: resolve each `depends_on` id by globbing 4 locations, accept `completed`/`completed-with-issues`, cache per invocation, honor the `dependencies:` alias | **PARTIAL** — `model.go` `resolveDependsOn` L786 + `PendingReady`/`PendingWaiting` L1006–1008 implement exactly this; not exposed | `queue-kanban ready-set` |
| Step 1, L174 | Walk the `depends_on` graph collecting a seen set; on self-hit write `status: blocked-dependency-cycle` | **PARTIAL** — the status string exists in `model.go` L764 and `annotateDependencyState` L1088 warns on unresolvable ids; nothing computes the cycle or writes the status | `queue-kanban ready-set --cycles` |
| Step 1, L176–182 | `--wave N` depth math: depth 0 = no list or all-archived deps; depth K = max(dep depth)+1; a non-pending non-archived member contributes 0 | **NONE** | `queue-kanban wave-depth` |
| Step 1, L188–198 | Categorize every REQ by status and print `Queue: N pending \| N completed/done \| …`, counting `completed`/`completed-with-issues`/`cancelled`/`done` together and `blocked`/`blocked-archive-collision` separately | **PARTIAL** — `writeBoardSummary` (`main.go` L107–128) prints board columns, not these status categories | `queue-kanban summary --queue-line` |
| Step 1, L200 | Expand `UR-NNN` → member REQs by scanning `user_request:` frontmatter (never the UR's `requests:` array), dedupe the union, tag each with named/UR-expanded provenance | **NONE** | `queue-kanban resolve-targets` |
| Step 1, L211 | Read `assigned_to` verbatim; skip-and-report in default mode, override + clear on explicit targeting | **PARTIAL** — `verify.go` `assigned-elsewhere-claimed-here` L29 catches only the post-claim violation | `queue-kanban ready-set` |
| Step 1, L219 | Verify required fields `id`, `status`, `title`; skip + report a REQ with missing/unparseable frontmatter | **PARTIAL** — `frontmatter.go` L22 returns a skip-gracefully signal and `model.go` warns; no required-field gate anywhere | `queue-kanban lint-queue` |
| Step 2.0, L227–231 | Glob `do-work/archive/**/REQ-NNN-*.md` and `…/REQ-NNN.md` for a collision | **FULL** | `tools/checks/archive-collision.sh` |
| Step 2, L237–239 | `mkdir -p do-work/working`, move the file, flip `status`/`claimed_at`, remove `assigned_to`, append a `CHECKPOINT.md` entry carrying `writer: <hostname>:<abs-path>` | **NONE** | `tools/checks/claim-req.sh` |
| Step 5.5, L343 | Compare declared `## Scope` paths against Implementation Summary paths (both set-differences) | **FULL** | `tools/checks/scope-drift.sh` |
| Step 5.75, L353–357 | `git status --porcelain --untracked-files=all` outside `do-work/`, run the test command, check deps; write `baseline.json` | **FULL** | `tools/checks/preflight.sh` |
| Step 6, L367–371 | Normalize `domain`, `tdd`, `caveman`, `maintenance` per the Schema Read Contract to decide which `crew-members/*.md` load | **NONE** — verified: `model.go` reads `domain` verbatim via `coerceScalarToString` L626; only `status` (L718) and `testing_status` (`testing.go` L59) have normalizers, and neither covers these four | `queue-kanban normalize-fields` |
| Step 6, L396–403 | Hand-back merge: path-limited bookkeeping commit, `git rev-parse --short HEAD` ×2, `git diff --name-only <pre>...<name> -- do-work/` queue guard, `git merge --no-ff --no-commit`, detect the literal `Already up to date.` | **PARTIAL** — `verify.go` `worktree-committed-queue-state` L559–601 reads the same three-dot diff, but only after the fact and blind post-merge (`work-reference.md` L318 says so explicitly) | `tools/checks/handback-merge.sh` |
| Step 6.3, L429–441 | Files exist / show in diff; P-A-U box audit vs debug artifacts; wiring grep; the "only `do-work/` paths ⇒ not implemented" rule | **FULL** (items 1/4/5-grep). Items 2, 3, 6 are judgment by design | `tools/checks/qualify.sh` |
| Step 6.5 item 4, L469 | Read `baseline.json`'s `launched` field; compare each failing test against `baseline-failures.txt` "mechanically" to separate regressions from pre-existing failures | **PARTIAL** — `preflight.sh` *writes* both artifacts; nothing reads or diffs them | `tools/checks/baseline-diff.sh` |
| Step 8 substep 5, L563 | Walk the `addendum_to` chain (honoring `amends`/`parent`/`amendment_to`) collecting a seen set; abort follow-ups on a cycle | **NONE** | `queue-kanban chain-check` |
| Step 8 substep 6, L568 | UR closure: check `do-work/queue/`, `working/`, `archive/` root, `archive/UR-NNN/` for members and test each against the terminal-resolved set | **PARTIAL** — `verify.go` `ur-archived-with-live-member` L30 flags the violation after a wrong close; nothing computes the predicate that makes the decision | `queue-kanban ur-status` |
| Step 8 substep 7, L572–585 (+ `work-reference.md` L580–587) | Count the prime file's directory depth, prepend that many `../`, verify the resolved path exists before writing the link | **NONE** | `tools/checks/prime-link.sh` |
| Step 8 substep 8, L587 | `git worktree remove <path>` (no `--force`), `git branch -d <operative_name>` from the integration branch, `git worktree prune`; treat each refusal as signal | **PARTIAL** — `verify.go` classifies leftovers by merge state (L21–23); nothing performs the happy-path removal | `tools/checks/worktree-cleanup.sh` |
| Step 8 blocked-flip, L595–600 | `git status --porcelain -- . ':(exclude)do-work/'`; `git rev-parse --verify -q '<operative_name>'`; `git rev-list --count HEAD..<operative_name> > 0` | **NONE** | `tools/checks/landed-edits.sh` |
| Step 9, L608–623 | `git rev-parse --git-dir` gate; resolve the version source across `package.json`/`Cargo.toml`/`pyproject.toml`/`VERSION`/tags/changelog counter (`work-reference.md` L718–724); grep the changelog for title reuse; assert strictly-greater version; stage explicit paths; validate the staged list against the Implementation Summary or `git diff --name-only <pre>..<merge_hash>` | **PARTIAL** — `next-version` writes one line in one file (`main.go` L23); `verify` checks version-vs-changelog, not-ahead, and title reuse (`verify.go` L17–19). The multi-source resolution, the bump-size decision, and the staged-list validation are all prose | `queue-kanban next-version --resolve-source`; `tools/checks/validate-staged.sh` |
| Step 9, `work-reference.md` L794–809 | Write the hash into `commit:`, then `--verify` what the metadata commit actually landed | **FULL** | `tools/checks/record-commit-hash.sh` |
| Step 9 fallback, `work-reference.md` L821 | `git diff --numstat HEAD -- <req-file>` must read `1\t1` (or `1\t0`); `test -s <req-file>` | **FULL as a fallback** — this *is* the prose reimplementation of the shipped guard | `tools/checks/record-commit-hash.sh` |
| Step 10, L627 | Fresh re-check of `do-work/queue/` for `REQ-*.md` | **PARTIAL** — same as Step 1 row 1 | `queue-kanban next-work` |
| Step 10, L642–658 | Rewrite `CHECKPOINT.md` wholesale while carrying every entry this checkout did not write through verbatim; delete only once no foreign/label-less entry remains | **NONE** | `tools/checks/checkpoint-merge.sh` |

### `actions/work-reference.md`

| Step | Mechanic | Coverage | Candidate script |
| --- | --- | --- | --- |
| Timestamp rule, L101–107 | Obtain the UTC instant: `queue-kanban now`, else `date -u +%Y-%m-%dT%H:%M:%SZ`, else the PowerShell `.ToUniversalTime()` form | **FULL** | `queue-kanban now` (`main.go` L28) |
| Timestamp rule, L111 | Date-only stamp: `date -u +%F` / PowerShell equivalent — explicitly *no* tool subcommand | **NONE** (deliberate; documented) | — |
| Schema Read Contract, L200–210 | Nine enum/boolean fields, each with an alias map, a canonical enum, and a documented default + warning | **PARTIAL** — 2 of 9 have a Go normalizer (`status`, `testing_status`); 7 (`domain`, `route`, `caveman`, `maintenance`, `tdd`, `error_type`, `kb_status`) have none anywhere, and none of the 9 is CLI-reachable | `queue-kanban normalize-fields` |
| Terminal-success / -resolved sets, L216–228 | Membership test over `{completed, completed-with-issues}` and `{…, cancelled}` — the predicate ~8 named readers restate | **PARTIAL** — `normalizeStatus` L718 canonicalizes but exposes no membership predicate | `queue-kanban normalize-fields --set terminal-success` |
| Target ID Resolution, L230–238 | Canonicalize case-insensitive `REQ-`/`UR-` + digits, match digits by **numeric value** against the zero-padded stored id, expand UR by frontmatter scan | **NONE** | `queue-kanban resolve-targets` |
| Crash Recovery, L242–283 | Read `CHECKPOINT.md`; derive this checkout's `writer:` label from `hostname -s` (or plain `hostname`) + `git rev-parse --show-toplevel`; classify each `working/` file own/foreign/label-less; compute age from `claimed_at` against a 3 h threshold with 2 min skew | **PARTIAL** — `verify.go` `claim-needs-attention` uses the same `staleClaimThreshold = 3 * time.Hour` (L38) and `checkpoint-names-missing-req` reads checkpoint ids via `requestIdMentionPattern` L49; neither derives the writer label nor classifies own-vs-foreign | `tools/checks/classify-claims.sh` |
| Crash Recovery worktree sweep, L277–281 | `git worktree prune`, then enumerate `git worktree list --porcelain` and `git branch --list 'worktree-agent-*'`, pair them, read the REQ id out of each name | **FULL for classification** — `verify.go` L715 runs the identical `branch --list` with `worktreeAgentNamePrefix` L45; the mechanical removal half is uncovered | `tools/checks/worktree-cleanup.sh` |
| Crash Recovery substeps 1–3, L273–275 | Reset frontmatter (status→`pending`/`pending-answers`, preserve `blocked`, remove `claimed_at`+`route`, stamp `status_changed_at`, conditionally clear `write_set`), strip 13 named `##` sections, move back to `queue/` | **NONE** | `tools/checks/recover-req.sh` |
| Fan-Out auto-wave, L345–356 | Four-condition ready set (`pending` + dependency-ready + unclaimed + not assigned elsewhere), bounded to N, take the first N in numeric id order | **PARTIAL** — conditions 1–2 are `model.go`'s `PendingReady`; 3–4 are uncovered | `queue-kanban ready-set --bound N` |
| Changelog Entry Procedure, L716–736 | Grep for title reuse; assert strictly-greater version; resolve the version source across 3 tiers with two disagreement guards; apply the bump-size table | **PARTIAL** — see the Step 9 row above | `queue-kanban next-version --resolve-source` |
| Commit Procedure, L790 | Compare the Implementation Summary file list against the staged set (or `git diff --name-only <pre>..<merge_hash>`), excluding `do-work/` | **NONE** | `tools/checks/validate-staged.sh` |

### `actions/forensics.md` — 14 checks, 12 of them mechanical

| Check | Mechanic | Coverage | Candidate script |
| --- | --- | --- | --- |
| 1, L29–39 | Read `claimed_at` from every `working/REQ-*.md`, compute age, band at >1 h / >24 h | **PARTIAL** — `verify.go` `claim-needs-attention` (3 h band, L38) covers the same field with different thresholds | `tools/checks/classify-claims.sh` |
| 2, L43–50 | Scan `archive/` (incl. `UR-*/`) for terminal-success status, test for `## Implementation Summary` + non-`do-work/` paths in `**Files changed:**` | **NONE** | `tools/checks/hollow-scan.sh` |
| 3, L54–59 | Scan archived REQs for a missing `## Qualification`; date-gate on the presence of `## Scope`/`## Pre-Flight` | **NONE** | `tools/checks/section-audit.sh` |
| 4, L63–73 | UR closure predicate across 4 locations against the terminal-resolved set | **PARTIAL** — same as work Step 8 substep 6 | `queue-kanban ur-status` |
| 5, L77–83 | Parse all Implementation Summary file lists, build `path → [REQ ids]`, flag 3+ unrelated `user_request` values | **NONE** | `tools/checks/hotspot-map.sh` |
| 6, L87–91 | Scan `status: failed` REQs for `error_type`; look for a follow-up whose `addendum_to` points back | **NONE** | `queue-kanban chain-check` |
| 7, L95–103 | Age `pending-answers` from `created_at` (3–7 d / >7 d) and `blocked` from `blocked_at` with `created_at` fallback (7–14 d / >14 d) | **NONE** | `tools/checks/queue-age.sh` |
| 8, L107–113 | For the last 10 REQs with `commit:`, read the file list and `git log --since` each path; flag missing `(new)` files | **NONE** | `tools/checks/divergence-scan.sh` |
| 9, L117–126 | Scan `queue/` + `working/` for any terminal status incl. non-standard variants; group by `user_request` | **FULL** | `queue-kanban verify` (`stranded-finished-req`, `verify.go` L28) |
| 10, L130–141 | `find do-work/archive -name 'REQ-*.md'`, read each `## Lessons Learned`, group by theme, count **distinct** REQs, band at 2 / 3+ | **NONE** (the grouping is judgment; the enumeration + per-REQ dedupe is not) | `tools/checks/lessons-collect.sh` |
| 11, L145–152 | Read every REQ's `status:` across 3 locations, judge against the Schema Read Contract vocabulary + alias set, **skip files with no parseable frontmatter** | **PARTIAL** — `normalizeStatus` L718 is the same judgment in Go, and `model.go`'s `bucketColumns` raises the board warning this check mirrors (L154 says so); not CLI-reachable | `queue-kanban lint-queue` |
| 12, L156–162 | Parse every `*_at` field across 3 locations, compare against now + ~2 min skew | **PARTIAL** — `model.go` `detectFutureTimestampFields` is the same logic (cited at L163); not CLI-reachable, and `verify` only checks `claimed_at` | `queue-kanban lint-queue` |
| 13, L166–179 | 0-byte / unparseable-frontmatter scan with git recovery-point resolution | **FULL** | `tools/checks/blanked-req-scan.sh` |
| 14, L181–208 | `go build` then `queue-kanban verify --repo-root` | **FULL** | `queue-kanban verify` |

### `actions/cleanup.md` — 7 passes

| Pass | Mechanic | Coverage | Candidate script |
| --- | --- | --- | --- |
| 0, L33–40 | Glob `queue/REQ-*.md`, read `status`, rewrite 6 non-standard terminal statuses to canonical, move to `archive/`, drop the own-label `CHECKPOINT.md` entry | **PARTIAL** — `verify`'s `stranded-finished-req` finds them; the normalize + move + checkpoint edit is prose | `tools/checks/sweep-finished.sh` |
| 1, L50–65 | Collect UR members by scanning `user_request:` across 4 locations; test against the terminal-resolved set; detect the same id in both `archive/` root and `archive/UR-NNN/`; cross-check the `requests:` array report-only | **PARTIAL** — `verify`'s `ur-archived-with-live-member` + `duplicate-req-id` cover fragments | `queue-kanban ur-status` |
| 2, L70–80 | For each loose `archive/REQ-*.md`, read `user_request`, test which UR folder exists, move or warn | **NONE** | `tools/checks/consolidate-loose.sh` |
| 3a, L86–95 | Find any directory named `do-work/` outside the project root; per-subtree relocate with same-number conflict detection | **NONE** | `tools/checks/find-misplaced-dowork.sh` |
| 3b, L101–107 | Test for `archive/user-requests/`, move or merge each `UR-NNN/`, remove if empty | **NONE** | `tools/checks/consolidate-loose.sh` |
| 4, L115–118 | Glob `do-work/runs/*/`, read each `manifest.md`'s `Status:` line, delete only on `consumed` | **NONE** | `tools/checks/sweep-runs.sh` |
| 5, L128–131 | `git worktree prune`; enumerate `worktree list --porcelain` + `branch --list 'worktree-agent-*'`; try `worktree remove` + `branch -d` from the integration branch; refusal ⇒ consent gate | **PARTIAL** — `verify.go` classifies (L21–23); the removal is prose | `tools/checks/worktree-cleanup.sh` |
| 6, L145–161 | `blanked-req-scan.sh`, `--restore --dry-run`, `--restore`; read exit 0/1 and the `SKIP:`/`FAIL:` lines | **FULL** | `tools/checks/blanked-req-scan.sh` |
| Repoint, L169–181 | Record old→new for every move; `git grep -l -F '<filename>' -- '*.md' ':(exclude)do-work'`; rewrite path occurrences while preserving `#fragment` and skipping bare prose mentions | **NONE** | `tools/checks/repoint-links.sh` |
| Commit, L237–275 | `git rev-parse --git-dir` gate; stage exact paths incl. `git add -u -- <deleted-run-path>` | **NONE** | — (judgment-adjacent; low value) |

### Remaining `actions/` files

| File · step | Mechanic | Coverage | Candidate script |
| --- | --- | --- | --- |
| `stray-check.md` Step 1, L36–39 | `git rev-parse --git-dir` gate; build a noise skip-list scoped to untracked/ignored only | **NONE** | `tools/checks/stray-scan.sh` |
| `stray-check.md` Step 2, L43–47 | `git ls-files` + `git ls-files --others --exclude-standard`; `git ls-files \| git check-ignore --no-index --stdin`; `find <root> -type d -empty` with prunes; binary-extension skip; >500-line sampling | **NONE** | `tools/checks/stray-scan.sh` |
| `stray-check.md` Step 3, L55–64 | Ten detection categories, each a glob/extension/size/tracked-state test (temp files, ignored-but-tracked, build artifacts, secrets, dupes, empties, >5 MB blobs, AI scratch, dead-code grep) | **NONE** | `tools/checks/stray-scan.sh` |
| `version.md` Version request, L35–38 | Read the first ~80 lines of `CHANGELOG.md`, split at `## ` headings, take the first 5, reverse | **NONE** | `queue-kanban changelog-tail` |
| `version.md` Update check, L44–47 | Fetch upstream, extract `**Current version**:`, semver-compare | **PARTIAL** — `next-version` writes that line but never reads a remote or compares | `tools/do-work-update.sh` (exists; not wired to this step) |
| `version.md` Step 2, L54–70 | Resolve `<skill-root>` and `<project-root>` (`git -C … rev-parse --show-toplevel 2>/dev/null \|\| pwd`); test descendant-ness; refuse under any user-wide skills path | **NONE** | `tools/do-work-update.sh` |
| `version.md` Step 3, L73–95 | `git status --porcelain -- <11 shipped paths>`; download the tarball to deterministic paths; `tar xzf` with 6 `--exclude`s; loop `diff -ru --new-file` per shipped path filtered by `grep -v 'tools/queue-kanban/queue-kanban'` | **PARTIAL** — `tools/do-work-update.sh` ships and does this class of work; this action's prose is a parallel implementation | `tools/do-work-update.sh` |
| `version.md` Step 4–6, L98–118 | `cp -R` snapshot; two `find … -delete` pre-cleans; `test -s` on the reviewed tarball; `tar xzf`; two `rm -f` stale-file lists; re-run the diff loop | **PARTIAL** — same | `tools/do-work-update.sh` |
| `version.md` Recap, L151–157 | Read `input.md` titles; scan `queue/REQ-*.md` for matching `user_request:`; merge + dedupe by UR id (archive wins); sort descending; take 5; label by member status | **NONE** | `queue-kanban recap` |
| `board.md` Step 2–5, L43–90 | `go version` gate; `git rev-parse --show-toplevel 2>/dev/null \|\| pwd`; `go build`; `git check-ignore -q` + `git rev-parse --git-path info/exclude` idempotent append | **FULL** for the board itself; the wrapper is thin by design | `queue-kanban` |
| `commit.md` Step 1, L45–61 | `git rev-parse --git-dir`; `git status --porcelain --untracked-files=all`; categorize M/A/D; exclude 4 secret-shaped glob classes | **NONE** | `tools/checks/uncommitted-inventory.sh` |
| `commit.md` Step 3, L78–88 | Glob `archive/**/REQ-*.md`; read `commit:` + terminal-success `status`; parse `## Implementation Summary` file lists; path-match; tie-break on latest `completed_at` | **NONE** | `tools/checks/associate-files.sh` |
| `inspect.md` Step 1, L61–78 | Identical `git rev-parse --git-dir` + `git status --porcelain -uall` + M/A/D + secret-exclusion block as `commit.md` L45–61 (copy-paste pair) | **NONE** | `tools/checks/uncommitted-inventory.sh` |
| `inspect.md` Step 2, L92–101 | `git show <commit>:<path>` per Implementation Summary entry, with `(deleted)` special-case and `HEAD` fallback | **NONE** | `tools/checks/associate-files.sh` |
| `inspect.md` Step 3, L119–127 | Glob `archive/**/REQ-*.md`; read `status`; parse Implementation Summary; path-match; tie-break on `completed_at` (copy of `commit.md` L78–88) | **NONE** | `tools/checks/associate-files.sh` |
| `inspect.md` Step 5, L156–196 | Regex scans of added lines: 6 WIP markers, 3+ consecutive comment lines, empty bodies, 6 debug statements, 5 placeholder values, test-file existence patterns, 5 secret-prefix classes | **NONE** | `tools/checks/readiness-scan.sh` |
| `review-work.md` Step 1, L44–46 | Find the highest REQ number with terminal-success status across `archive/` root + all `UR-NNN/`; `git diff --stat <pre>..<merge_hash>` emptiness test | **NONE** | `queue-kanban latest-completed` |
| `review-work.md` Step 4, L66–68 | `git diff` → `git diff --staged` → merge range; `git rev-parse --verify -q '<commit>^2'` to detect a merge commit, then `git show --first-parent -m` | **NONE** | `tools/checks/req-diff.sh` |
| `review-work.md` Step 6, L115 | Compare Implementation Summary file list against `## Scope` | **FULL** | `tools/checks/scope-drift.sh` |
| `review-work.md` Step 6, L124–133 | Restatement sweep — grep every other statement/consumer of a redefined token | **NONE** (irreducibly judgment; the grep enumeration is not) | — |
| `ai-report.md` Step 3a, L70 | `git diff-tree --no-commit-id --name-only -r -m <commit> \| grep -E '\.(png\|jpg\|gif)$' \| sort -u` | **NONE** | `tools/checks/req-diff.sh` |
| `ai-report.md` Step 3b, L81–97 | `date +%Y-%m-%d_%H%M`; `mkdir -p`; `playwright-cli --help` probe; a 3-port `curl -s -o /dev/null -w '%{http_code}'` loop | **NONE** | `tools/checks/probe-dev-server.sh` |
| `ai-report.md` Step 4d, L174–176 | `mkdir -p`; canonicalize to an absolute `$GEN`; background jobs + `wait`; `[ -s "$f" ]` verify per expected file | **NONE** | — (report-local) |
| `ai-report.md` Step 5/7, L212, L249–258 | `git rev-parse --verify -q '<sha>^2'` merge test; `python3 -m http.server`; `lsof -ti :8123 \| xargs kill` | **NONE** | `tools/checks/req-diff.sh` |
| `present-work.md` L42, L62, L154–156, L397 | Highest terminal-success UR/REQ scan; `git show <commit>`; the same `'<sha>^2'` merge-commit test repeated 4× in one file | **NONE** | `queue-kanban latest-completed`; `tools/checks/req-diff.sh` |
| `present-work.md` Step 1, L429–433 | Scan the archive, read `title`/`route`/`commit`/`completed_at` frontmatter + parse the `## Review` block's `Overall: X%` | **NONE** | `queue-kanban recap` |
| `pipeline.md` Step 3, L84 | `git ls-files -- do-work/pipeline.json` (a `check-ignore` pass is powerless over an indexed file) + local-ignore append | **NONE** | `tools/checks/local-ignore.sh` |
| `pipeline.md` Step 4, L91, L101–127 | Find the first step whose `status` is `pending`/`in-progress`/`failed`; write `status`/`completed_at`; parse action output for REQ/UR ids | **NONE** | `tools/checks/pipeline-state.sh` |
| `pipeline.md` Step 5, L137–142 | Read each REQ's `id`/`title`/`commit`/`domain` frontmatter + `## Implementation Summary` + `## Testing`; group rows by domain; `git rev-parse --verify -q '{sha}^2'` again | **NONE** | `queue-kanban recap`; `tools/checks/req-diff.sh` |
| `pipeline.md` Step 5, L140 | Scan for `status: pending-answers` REQs, deferred items in `## Lessons Learned`, TODO/FIXME introduced by this run's commits | **NONE** | `tools/checks/carry-forward.sh` |
| `pipeline.md` Step 5/5a, L156, L168 | Scan `queue/REQ-*.md` for `status: pending`, excluding ids in the pipeline's `artifacts` array; loop | **PARTIAL** — same gap as work Step 1 | `queue-kanban ready-set` |
| `pipeline-reference.md` L41 | Every `pipeline.json` write must be multi-line pretty-printed, because `hooks/pipeline-guard.sh`'s non-`jq` fallback counts pending steps with line-oriented `grep -c` | **NONE** — an invariant enforced only by prose, consumed by a shipped hook | `tools/checks/pipeline-state.sh` |
| `pipeline-reference.md` L154, L254 | Glob **both** `ai-reports/*{UR-NNN}*/index.html` and legacy `ai-reports/*{UR-NNN}*.html`, pick the newest by mtime | **NONE** | `tools/checks/find-report.sh` |
| `capture.md` L83 | `queue-kanban next-req` as the preferred REQ-number source, with a hand-scan fallback; UR numbering stays manual | **PARTIAL** — `next-req` (`main.go` L163) covers REQ numbers only, by its own docs | `queue-kanban next-ur` |
| `capture.md` L111–117 | Read every `queue/REQ-*.md`'s `title` + heading + `## What` for duplicate-intent matching; filename-scan `working/` + `archive/` (incl. `UR-*/`) | **NONE** (the intent match is judgment; the enumeration is not) | `queue-kanban lint-queue` |
| `capture.md` L223 | The UR-closure predicate restated: scan `user_request:` across 4 locations, never the `requests:` array | **PARTIAL** — same as work Step 8 | `queue-kanban ur-status` |
| `capture.md` L243–269 | `git rev-parse --git-dir` gate; stage exact created paths | **NONE** | — |
| `capture-reference.md` L130 | Capture-side normalize-and-warn over 8 shared enum fields before emitting a REQ | **NONE** — see the Schema Read Contract row | `queue-kanban normalize-fields` |
| `roadmap.md` L59–61 | Read 13 frontmatter fields per REQ + detect which `##` sections exist; resolve 4 field-name aliases; normalize 7 enum fields; count non-canonical survivors for the header | **NONE** | `queue-kanban normalize-fields`; `queue-kanban lint-queue` |
| `roadmap.md` L67–69 | Ready/Needs-clarification/Blocked classification from `depends_on`, `addendum_to` (+3 aliases), `status: blocked`/`blocked-dependency-cycle` | **PARTIAL** — `resolveDependsOn` L786 + `PendingReady` L1006 cover the dependency half | `queue-kanban ready-set` |
| `roadmap.md` L76–82 | TDD posture from `tdd:` + presence of `## Red-Green Proof` | **NONE** | `tools/checks/section-audit.sh` |
| `roadmap.md` L90–108 | Locate `kb_entry` recursively under 3 `raw/` branches; two `find` name patterns incl. a `[0-9]{6}-` prefix; bucket by `processed/YYYY-MM-DD/` path vs mtime, with a documented mtime-resets-on-clone caveat | **NONE** | `tools/checks/kb-entry-locate.sh` |
| `clarify.md` Step 1, L29 | Find all `queue/REQ-*.md` with `status: pending-answers`; collect `status: blocked` for Step 5.5 | **PARTIAL** — the board buckets both; no command lists them | `queue-kanban ready-set --status` |
| `clarify.md` L104–160 | Six distinct frontmatter transitions (`pending`, `completed`, `cancelled`, unblock) each stamping `completed_at`/`status_changed_at` and, on unblock, removing `blocked_by`+`blocked_at` while keeping `blocked_check` | **NONE** | `tools/checks/set-status.sh` |
| `abandon.md` L33–41 | Resolve targeting tokens + UR expansion by `user_request:` scan; glob 4 path shapes per id; gate on where a `failed` REQ was found (`archive/` root vs `legacy/` vs inside a closed `UR-NNN/`) | **NONE** | `queue-kanban resolve-targets` |
| `abandon.md` L58–74 | Set `status: cancelled` + `completed_at`; retain `error`/`error_type`; conditionally emit the `Previously: failed` line by field presence | **NONE** | `tools/checks/set-status.sh` |
| `verify-requests.md` L165, L179 | Detect REQs lacking `user_request`; assert verify never writes `status: pending-answers` | **NONE** | `queue-kanban lint-queue` |
| `kb-lessons-handoff.md` L20, L58–63, L141–150 | Gate on `kb_status` absent-or-`pending`; read `title`/`completed_at`/`id`/`domain`/`prime_files` with a today's-date fallback; write back `kb_status` + `kb_entry` touching no other field | **NONE** | `tools/checks/set-frontmatter-field.sh` |
| `memory.md` L49, L60, L77 | `wc -c` against a 2,500-char cap; update the `updated:` frontmatter date; report cap usage, log-file count, newest date, `tail -5` of the ledger | **NONE** | `tools/checks/memory-status.sh` |
| `memory-reference.md` L16, L70, L104–110 | `PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null \|\| pwd)"`; consolidation loop re-checking `wc -c` ≤ 2,500 | **NONE** | `tools/checks/memory-status.sh` |
| `memory-reference.md` L132–136, L149, L158 | `printf … >> usage-ledger.jsonl 2>/dev/null \|\| true`; `grep -q "session capture $hash8"` idempotence gate; `sed -E` credential redaction over 7 token shapes **before** truncation | **NONE** | `hooks/` (partly shipped as hook scripts) |
| `memory-reference.md` L188 | `jq . "$settings_file" >/dev/null` parse check; assert both hook filenames present; compare entry counts against the backup | **NONE** | `tools/checks/merge-settings.sh` |
| `memory-value.md` L29–33 | `find <kb-root>/wiki -name '*.md' \| wc -l`; `git log --oneline -- <kb-root>/ \| wc -l`; `git log -1 --format=%ci`; `git log --format=%an \| sort -u`; grep date headings in `log.md`; `grep -rl 'wiki/'` inbound-reference scan | **NONE** | `tools/checks/engine-probe.sh` |
| `memory-value.md` L37–40 | `memory/working-memory.md` char count vs cap; `ls memory/logs/*.md \| wc -l`; heading-ratio count; grep `.claude/settings.json` for both hook names | **NONE** | `tools/checks/engine-probe.sh` |
| `memory-value.md` L44–48 | Bucket `usage-ledger.jsonl` events per week over 4 weeks; compute `hit_cited ÷ retrievals`; newest-event age | **NONE** | `tools/checks/ledger-stats.sh` |
| `dream.md` L71, L100–101 | Parse YAML frontmatter between the first two `---`; prefer `last_updated` → `updated` → `created`; parse `%Y-%m-%d` or `%Y/%m/%d`; compute `(today − parsed).days > 90` | **NONE** | `tools/checks/kb-scan.sh` |
| `dream.md` L119–126 | Pairwise title comparison (lowercase + whitespace-normalize), ≥80% shared-token test over the shorter title; `ls -la --time-style=long-iso` mtime vs `last_updated` | **NONE** | `tools/checks/kb-scan.sh` |
| `bkb.md` L135, L164–175, L286–300 | `git init` when the KB is not in a repo; route inbox files by `source_type:`/URL-in-frontmatter; scan the first 500 chars for a topic cluster; normalize-and-warn over 3 enum + 1 namespace field across 6 sub-commands | **NONE** | `tools/checks/kb-scan.sh` |
| `bkb.md` L410, L436, L446, L467, L507 | Assert index article counts match actual file counts; derive the last-lint date by scanning `log.md`; match `[RESOLVED] contradiction:` against open flags; count `raw/inbox/` files and "ready" queue items | **NONE** | `tools/checks/kb-scan.sh` |
| `prompts.md` L64, L100, L115, L132 | Glob `prompts/*.md` minus `README.md`; parse the `**Runnable:**` / `**Aliases:**` header lines, first token only, lowercased, stopping at the first `---`; resolve `<cwd>` and `<skill-root>` to absolute paths and test descendant-ness | **NONE** | `tools/checks/prompt-index.sh` |
| `note.md` L40–41 | `date +%F`; `mkdir -p do-work/`; create `notes.md` if absent; append one `- [YYYY-MM-DD] <text>` bullet | **PARTIAL** — `queue-kanban now` gives the instant, not the `%F` date (`work-reference.md` L111 states there is deliberately no subcommand) | — (2 lines; not worth extracting) |
| `quick-wins.md` L32–34, L114–115 | Language detection by marker file; glob by extension; approximate line counts; `git log --format='' --name-only \| sort \| uniq -c \| sort -rn` churn ranking | **NONE** | `tools/checks/repo-inventory.sh` |
| `code-review.md` L91–95, L109 | Parse `$ARGUMENTS` into primes + dirs; glob source files by detected language; build a file list with line counts; `mkdir` a run dir named from `date +%Y-%m-%d-%H%M%S` and seed `manifest.md` | **NONE** | `tools/checks/repo-inventory.sh` |
| `validate-feedback.md` L66 | Inspect `git diff` / `git log` / `git show` / staged changes for evidence an item was already handled | **NONE** (judgment-heavy) | — |
| `slop-check.md` L27, L43 | Glob `do-work/deliverables/*.md` + `*.single.html`, exclude `*.marp.html` and `*-video/`, pick newest by mtime; `wc -w` | **NONE** | `tools/checks/find-deliverable.sh` |
| `tidy-repo.md` L32–33 | Run the test suite and record pass/fail counts; `git status --short --untracked-files=all` and record dirty paths | **NONE** | `tools/checks/repo-inventory.sh` |
| `tidy-repo.md` L38–43 | `git ls-files \| awk -F/ '{print $1}' \| sort \| uniq -c \| sort -rn`; `git ls-files --others --exclude-standard` | **NONE** | `tools/checks/repo-inventory.sh` |
| `tidy-repo.md` L100–101 | Straggler grep for old paths; `git diff --check`; `git diff --summary --find-renames` rename verification | **NONE** | `tools/checks/repoint-links.sh` |
| `prime.md` L181, L192, L226–227 | `glob **/prime-*.md`; extract referenced file paths from a prime and verify each exists; `glob **/known-bugs-*.md`, `**/lessons-learned/**/*.md` | **NONE** | `tools/checks/prime-audit.sh` |
| `scan-ideas.md` L39 · `deep-explore.md` L91 | `Glob **/prime-*.md` (third and fourth copies of the same enumeration) | **NONE** | `tools/checks/prime-audit.sh` |
| `deep-explore.md` L76, L250 | Read `manifest.md` + `session/state.json`; branch on `synthesized`/`consumed`/legacy `complete`; write `writer_status`/`status`/`completed_at`/two counts; verify all Writer outputs exist | **NONE** | `tools/checks/run-state.sh` |
| `install.md` L44, L49, L53–57, L76, L110 | Per-target `detect_cmd`/`install_cmd`/`verify_cmd` triples: `test -s`, `mkdir -p`, `curl -fsSL -o …download && mv … \|\| { rm -f …; false; }`, `PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null \|\| pwd)"` | **NONE** — the manifest table *is* an executable spec written as prose | `tools/checks/install-target.sh` |
| `install.md` L264–269, L281–296 | `exclude=$(git rev-parse --git-path info/exclude 2>/dev/null) \|\| exclude=""`; local-ignore append; `git -C "$PROJECT_ROOT" rev-parse --git-dir` gate | **NONE** | `tools/checks/local-ignore.sh` |
| `install.md` L326–338 | `just-kanban`: compare an installed recipe block against the shipped version and offer a consent-gated upgrade | **PARTIAL** — `tools/do-work-update.sh` ships and covers the updater half | `tools/do-work-update.sh` |
| `sample-archived-req.md` | Example artifact only — no prescribed mechanics | n/a | — |
| `help.md`, `tutorial.md`, `ui-review.md`, `interview.md`, `interview-reference.md`, `bkb-reference.md`, `deep-explore-reference.md`, `ai-report-reference.md`, `memory.md` sub-command table, `dream.md` Phase 3 | No shell/frontmatter mechanics beyond those listed above | n/a | — |

### `prompts/`

| File · step | Mechanic | Coverage | Candidate script |
| --- | --- | --- | --- |
| `architecture-decisions-log_create-or-expand.md` L36–40 | `git status --porcelain` dirty gate; `git rev-parse --abbrev-ref HEAD` branch test against `main`/`master`; `git switch -c` | **NONE** | `tools/checks/adr-preflight.sh` |
| same, L221–223 (Phase 1) | Test for `decisions/_master_index.md`/`_progress.md`; on RESUME, scan `decisions/records/adr-*.md` for the highest number and every `req:`/`ur:` referenced | **NONE** | `tools/checks/adr-state.sh` |
| same, L235 | Grep every `[[…]]` reference in `decisions/` and verify each target exists | **NONE** | `tools/checks/wikilink-check.sh` |
| all 17 other `prompts/*.md` | **None** — zero shell commands, zero frontmatter reads, zero output parsing | n/a | — |

---

## 3. Frontmatter read or filtered WITHOUT `frontmatter.go` — flagged in full

**Every one of them.** This is not a sampling result: `frontmatter.go` is unexported Go with no CLI surface
(`main.go` L60–76 lists all seven subcommands; none takes a file-and-field argument), so a prose step *cannot* call it
even in principle. The flag list is therefore the complete set of frontmatter read sites in `actions/` and `prompts/`:

**Reads/filters `status`** — `work.md` L127, L188–198, L204, L206, L215, L219, L568; `work-reference.md` L203, L218,
L224, L246–249, L273, L347; `forensics.md` L43, L87, L95, L101, L117, L145–152; `cleanup.md` L34–40, L55, L72;
`clarify.md` L29, L104–160; `abandon.md` L35–41, L58; `review-work.md` L43–44; `ai-report.md` L47–49;
`present-work.md` L42; `commit.md` L80; `inspect.md` L121; `roadmap.md` L61–69; `pipeline.md` L156, L168;
`capture.md` L107; `capture-reference.md` L26, L130; `verify-requests.md` L179.

**Reads/filters `depends_on` (+ `dependencies:` alias)** — `work.md` L170–182, L186; `work-reference.md` L127, L348,
L667; `roadmap.md` L67–69.

**Reads/filters `addendum_to` (+ `amends`/`parent`/`amendment_to`)** — `work.md` L248, L563; `work-reference.md` L126,
L666; `forensics.md` L83, L91; `roadmap.md` L61, L67–69.

**Reads/filters `user_request`** — `work.md` L200, L568; `work-reference.md` L235, L384; `cleanup.md` L50, L72;
`forensics.md` L67, L81, L119; `capture.md` L223; `review-work.md` L60; `roadmap.md` L59; `version.md` L152.

**Reads timestamps (`claimed_at`, `completed_at`, `created_at`, `blocked_at`, `testing_updated_at`)** —
`work-reference.md` L253–256; `forensics.md` L33, L97–103, L158; `commit.md` L86; `inspect.md` L127;
`roadmap.md` L59; `kb-lessons-handoff.md` L59; `abandon.md` L74; `present-work.md` L433.

**Reads `commit:`** — `forensics.md` L109; `commit.md` L80; `inspect.md` L98; `review-work.md` L31, L68;
`ai-report.md` L56; `present-work.md` L62, L397; `pipeline.md` L137.

**Reads the seven un-normalized enum fields** (`domain`, `route`, `caveman`, `maintenance`, `tdd`, `error_type`,
`kb_status`) — `work.md` L291, L367–371; `work-reference.md` L202–209; `roadmap.md` L61, L78; `forensics.md` L90;
`kb-lessons-handoff.md` L20, L62; `capture-reference.md` L130. **No implementation of these seven exists anywhere in
the repo** — verified: `model.go` L626 reads `domain` verbatim via `coerceScalarToString`, and the only two normalizers
that ship are `normalizeStatus` (`model.go` L718) and `normalizeTestingStatus` (`testing.go` L59).

**Reads other fields** — `assigned_to`: `work.md` L211, L238; `work-reference.md` L130, L350. `write_set`:
`work.md` L335; `work-reference.md` L128, L273. `blocked_by`/`blocked_check`: `work.md` L129–168;
`work-reference.md` L138–140. `prime_files`: `work.md` L291, L313, L537; `kb-lessons-handoff.md` L63.
`kb_entry`: `roadmap.md` L90–108. Non-REQ frontmatter: `dream.md` L71, L101; `bkb.md` L164, L207, L236, L276;
`prompts.md` L100 (`**Runnable:**` header, not YAML); `pipeline-reference.md` L196 (Marp `marp: true`).

**Consequence, stated once.** The `status` vocabulary alone is read at ~35 prose sites. Each is an independent
opportunity to filter on the literal `completed` and silently drop `completed-with-issues` — which is exactly why five
separate Red Flags already exist warning about that one bug (`cleanup.md` L306, `commit.md` L224, `review-work.md`
L479, `ai-report.md` L341, `present-work.md` L527). Five documented instances of one bug class is the census's
strongest signal: the contract is correct and centralized; its *enforcement* is 35 hand copies.

---

## 4. Top 5 extraction candidates

Ranked by (execution frequency × bug risk of the prose reimplementation). Frequency is per work-loop iteration unless
noted; bug risk cites observed or documented failures, not speculation.

### 1. `queue-kanban frontmatter` — a field-read/normalize/membership CLI

**Frequency: highest in the skill.** ~35 `status` sites plus ~60 more across the other fields, and every work-loop
iteration hits several. **Bug risk: highest, and already realized five times.** The five Red Flags cited above are all
the same defect — a prose filter on literal `completed` — and seven of the nine normalize-and-warn fields have no
mechanical implementation at all, so `domain: back-end` silently mis-selects a crew file with nothing to catch it.
`frontmatter.go` already handles the hard parts prose cannot (CRLF bodies, duplicate top-level keys, lenient recovery
of block lists — L70–81, L109–117); the only thing missing is a way to *call* it. Smallest possible surface:
`queue-kanban frontmatter get <file> <field> [--normalize] [--in-set terminal-success|terminal-resolved]`.
This one script retires more prose than the other four combined.

### 2. `queue-kanban ready-set` — the queue-selection scan

**Frequency: twice per iteration minimum** (`work.md` Step 1 L127 and Step 10 L627), plus `pipeline.md` L156/L168,
`clarify.md` L29, `roadmap.md` L67. **Bug risk: high, and the logic is already written.** `model.go` computes
`PendingReady`/`PendingWaiting` (L1006–1008) using `resolveDependsOn` (L786) with the `dependencies:` alias — the exact
predicate `work.md` L170–172 spells out in prose. Shipping it as a subcommand also closes the four uncovered pieces the
prose owns alone: the `--wave N` depth math (L176–182), `depends_on` cycle detection (L174), the `assigned_to` skip
(L211), and the unclaimed test. Prose currently re-derives a computation the repo can already do, in the file that runs
most often.

### 3. `tools/checks/req-diff.sh` — resolve a REQ's diff, merge-commit-aware

**Frequency: once per review, once per report, and on every archived-REQ read.** **Bug risk: high, with a named
failure mode that is silent.** The `git rev-parse --verify -q '<sha>^2'` → `git show --first-parent -m` idiom is
copy-pasted at **seven** sites: `review-work.md` L68, `ai-report.md` L212, `present-work.md` L62, L154, L397,
`pipeline.md` L142, `pipeline-reference.md` L170. Get it wrong and `git show` on a merge prints a usually-empty
combined diff — the reviewer reads an empty diff as an empty REQ and passes it. The repo's own `CLAUDE.md` records this
trap and the rule that a prescribed-command fix must be grepped across every action; seven copies is the standing
violation of that rule. One script, one behavior, and the `<pre>..<merge_hash>` range handling comes along with it.

### 4. `tools/checks/uncommitted-inventory.sh` + `associate-files.sh`

**Frequency: every `commit`, `inspect`, `validate-feedback`, and `tidy-repo` run.** **Bug risk: high, and the trap is
already documented in-repo.** `commit.md` L45–61 and `inspect.md` L61–78 are a near-verbatim copy-pasted pair
(`git rev-parse --git-dir`, `git status --porcelain -uall`, M/A/D categorization, four secret-shaped exclusion globs),
and the REQ-association logic is a *third* copy across `commit.md` L78–88 and `inspect.md` L119–127. The `-uall` flag
is the load-bearing detail: without it `git status --porcelain` collapses a new directory to `?? dir/` and every file
inside escapes the secret-exclusion scan — a secret-leak path, not a cosmetic bug. Both files currently carry a
paragraph explaining the flag, which is prose doing a script's job. `stray-check.md` L159 lists the same miss as a Red
Flag, confirming it has been hit.

### 5. `tools/checks/classify-claims.sh` — writer-label claim classification

**Frequency: once per work-loop *session*, and on every forensics run.** Lower frequency than the four above, but
**the highest per-occurrence blast radius in the census.** `work-reference.md` L242–283 prescribes deriving this
checkout's identity (`hostname -s` with a plain-`hostname` fallback, plus `git rev-parse --show-toplevel`), comparing it
against each `CHECKPOINT.md` entry's `writer:` label, and classifying own-crash / foreign / label-less — where guessing
wrong runs recovery substeps 1–3, which strip thirteen generated sections from a REQ another checkout may be building
right now. `verify.go` already holds half the pieces (`staleClaimThreshold = 3 * time.Hour` L38, `requestIdMentionPattern`
L49, `claim-needs-attention` and `checkpoint-names-missing-req` probes) but never derives the label. The prose itself
argues the case: L248 records that a previous reading of local checkpoint state as authorship "strips a live foreign
claim." A destructive branch selected by hand-compared hostname strings is the one place in this census where a prose
slip destroys work rather than misreporting it.

**Honorable mention, excluded on frequency:** `tools/checks/stray-scan.sh` (`stray-check.md` L36–64) is the single
largest block of uncovered mechanics in the two directories — ten detection categories, the tracked-vs-untracked
skip-list distinction, and three documented traps — but it runs only when a user asks for a hygiene sweep.

---

## 5. What this census does not claim

- **It does not propose the extractions.** Every candidate name above is a name for a gap, not an approved design.
  Per this repo's own gate, an extraction that adds a `queue-kanban` subcommand has to justify itself against the
  compiled-tooling exception (`actions/board.md` is the only capability allowed to *need* a compiler), which means each
  one needs a documented shell floor or must be gated on the binary already existing.
- **It does not judge which mechanics *should* stay prose.** Some are judgment wearing shell clothing —
  `review-work.md`'s Restatement Sweep (L124–133), `forensics.md` Check 10's theme grouping (L136),
  `capture.md`'s duplicate-intent match (L111). Those are marked NONE because they contain a mechanical enumeration,
  not because the whole step should be a script.
