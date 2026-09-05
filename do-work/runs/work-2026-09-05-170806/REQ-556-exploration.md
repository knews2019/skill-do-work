# Pre-explore: REQ-556 — cut the debug-artifact rule prose `qualify` already enforces

HEAD at exploration time: `71eb49f3562217dddaad63d345746dd9ee243c6f` (branch `main`).
Audited commit named by the REQ: `dc8a64e3` ("[REQ-531] claim request lifecycle", 2026-09-03).

---

## 1. Re-verify the claim at HEAD — THE COUNT HAS MOVED

**The REQ's baseline of 9 is stale. At HEAD it is 7.** work.md dropped from 5 to 3; the other two files are unchanged.

Reproduce command run verbatim:

```
rg -n -e 'console\.log' -e 'debug artifacts' \
  skills/do-work/actions/work.md \
  skills/do-work/actions/work-reference.md \
  skills/do-work/actions/review-work.md
```

| File | REQ claim (at dc8a64e3) | At HEAD | Delta |
|---|---|---|---|
| `skills/do-work/actions/work.md` | 5 | **3** | −2 |
| `skills/do-work/actions/review-work.md` | 3 | **3** | 0 |
| `skills/do-work/actions/work-reference.md` | 1 | **1** | 0 |
| **total** | **9** | **7** | **−2** |

### The 7 hits at HEAD, with line numbers and verbatim sentences

**S1 — `skills/do-work/actions/work.md:292`** (bullet in the "All routes include these instructions to the agent" block)

> `- **P-A-U phasing is mandatory:** Edit the REQ's "AI Execution State (P-A-U Loop)" checkboxes in real time. [PLAN] writes a brief technical approach. [APPLY] stays in declared scope. [UNIFY] runs \`git diff --stat\`, runs native linters, verifies no debug artifacts, and lists each file checked (the orchestrator audits this during Qualification and Testing Judgment).`

**S2 — `skills/do-work/actions/work.md:574`** (Common Rationalizations table row)

> `| "P-A-U is bookkeeping — I'll just tick the boxes" | Do each phase; qualification audits the diff against the checked boxes | A checked \`[UNIFY]\` over a diff containing \`console.log\` is a false claim the qualifier will catch |`

**S3 — `skills/do-work/actions/work.md:583`** (Red Flags bullet)

> `- All P-A-U checkboxes marked complete but diff contains \`console.log\`, \`debugger\`, or \`TODO\` (debug artifacts)`

**S4 — `skills/do-work/actions/review-work.md:106`** (Step 6 Code Review → Code Quality bullet)

> `- Diff hygiene — no debug artifacts — console.log/print lines no contract reads (a check's own reporting is contract output), or temporary files left behind. **Protect lessons learned** — comments that document *why* something was done, what was tried and didn't work, or architectural reasoning are valuable and should stay. Strip noise, keep knowledge.`

**S5 — `skills/do-work/actions/review-work.md:374`** (the Review-Fix REQ **template** body, inside the fenced follow-up-REQ block under Step 10)

> `- [ ] **[UNIFY]:** (Agent: Run \`git diff --stat\` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)`

**S6 — `skills/do-work/actions/review-work.md:494`** (Red Flags bullet)

> `- Builder checked all P-A-U boxes but the diff contains \`console.log\`, \`debugger\`, or TODO/FIXME`

**S7 — `skills/do-work/actions/work-reference.md:604`** (row in `## Qualification Anti-Rationalization Table (Step 6.3)`)

> `| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |`

### Which two work.md sites disappeared, and why it matters

At `dc8a64e3` work.md carried two extra hits that no longer exist:

- `work.md:503` (old) — `It verifies checklist items **1 (files exist / show in diff)**, **4 (P-A-U box audit + debug artifacts in the diff)**, and the grep half of **5 (wiring)** ...`
- `work.md:515` (old) — `Apply the qualification anti-rationalization table in \`actions/work-reference.md\` → **Qualification Anti-Rationalization Table (Step 6.3)** (e.g., "the summary says files changed" → check the file system; "the builder checked UNIFY" → read the diff for debug artifacts).`

Both were removed by the REQ-504→REQ-506 "advance" chain (`a296ee9e`, `59fe3e3a`, `0e1b44ff`), which replaced the old `### Step 3.6` / `### Step 6.3` prose with the `### Mechanical Evidence-Gate Loop` and renamed `Step 6.3` to `### Qualification and Testing Judgment` (work.md:333).

**Two structural consequences the builder must handle:**

1. **`Step 6.3` no longer exists as a heading in `work.md`.** The REQ says "Keep one sentence in `work.md` Step 6.3 naming the three finding codes". The nearest live home is `### Qualification and Testing Judgment` (work.md:333) or the `### Mechanical Evidence-Gate Loop` (work.md:335 region).
2. **None of the three finding codes appear anywhere in the shipped prose today.** `rg 'QUALIFY-DEBUG-ARTIFACT|QUALIFY-PAU-UNCHECKED|QUALIFY-UNIFY-DISARMED'` outside `.go` files returns **zero hits**. So "keep one sentence naming the three finding codes" is an **addition**, not a keep. Under `maintenance.md` § 3 that addition needs a replay-pack justification, or it should be dropped and the pointer written to the command instead of to codes.

Stale-reference side note (out of scope, worth a discovered-task line): `work-reference.md:598` still titles the table `## Qualification Anti-Rationalization Table (Step 6.3)`, and lines 422 and 433 still cite "Step 6.3" / "Step 6.25" step numbers that work.md no longer has.

---

## 2. What the Go check actually enforces

File: `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go`. Entry point `handleQualify` at line 261.

### The artifact pattern (line 24)

```go
24: var qualificationDebugArtifactPattern = regexp.MustCompile(`\b(` + "debug" + `ger|` + "TO" + `DO|` + "FIX" + `ME)\b`)
```

Word-boundary match on three literals only: `debugger`, `TODO`, `FIXME`. (The source splits each literal across `+` concatenations so the file does not trip its own scan.)

A second, separate pattern covers output primitives (line 377):

```go
377: outputPattern := regexp.MustCompile(`console\.log|(^|[^[:alnum:]_])print\s*\(`)
```

`console.log` and `print(` are **not** `QUALIFY-DEBUG-ARTIFACT`. They produce `QUALIFY-LIBRARY-OUTPUT` (error, line 410), or `QUALIFY-OUTPUT-RELOCATED` / `QUALIFY-REPORTER-OUTPUT` (warnings, line 413). Every prose site that writes "`console.log` ⇒ debug artifact" is therefore already slightly wrong about the code's taxonomy.

### `QUALIFY-UNIFY-DISARMED` (lines 353, 355-357)

```go
353: unifyPattern := regexp.MustCompile(`(?m)^\s*-\s*\[[ x~]\]\s*\*\*\[UNIFY\]`)
355: if !unifyPattern.Match(contents) {
356:   findings = append(findings, helperFinding("QUALIFY-UNIFY-DISARMED", resultmodel.SeverityWarning, []string{requestPath}, "no [UNIFY] box exists", resultmodel.FixabilityManual, "P-A-U audit is not armed", nil, nil))
357: }
```

Reads **the REQ file only** (`contents` = the file at `--request-path`). Fires when there is no `[UNIFY]` line at all. Severity **warning**. Note the regex does not care what text follows `**[UNIFY]` — so trimming the explanatory tail off the template line at review-work.md:374 cannot disarm this check.

### `QUALIFY-PAU-UNCHECKED` (lines 354, 358-360)

```go
354: uncheckedPattern := regexp.MustCompile(`(?m)^\s*-\s*\[ \]\s*\*\*\[(PLAN|APPLY|UNIFY)\]`)
358: if count := len(uncheckedPattern.FindAll(contents, -1)); count > 0 {
359:   findings = append(findings, helperFinding("QUALIFY-PAU-UNCHECKED", resultmodel.SeverityError, []string{requestPath}, fmt.Sprintf("%d P-A-U checkbox(es) remain unchecked", count), resultmodel.FixabilityManual, "the builder has not completed every phase", nil, nil))
360: }
```

Also the REQ file only. Counts literal `- [ ]` boxes. Severity **error**. It is pure bookkeeping: it proves a box is ticked, and says nothing whatever about whether the diff is clean.

### `QUALIFY-DEBUG-ARTIFACT` (lines 386-405)

```go
386: for _, path := range changedPaths {
387:   if path == "do-work" || strings.HasPrefix(path, "do-work/") {
388:     continue
389:   }
390:   changes := lineChanges[path]
391:   for _, line := range changes.Added {
392:     if !qualificationDebugArtifactPattern.MatchString(line) {
393:       continue
394:     }
395:     code, severity := "QUALIFY-DEBUG-ARTIFACT", resultmodel.SeverityError
396:     evidence := line
397:     stop := "unfinished/debug-only code is newly introduced by the change"
398:     if removedMarkerLines[line] > 0 {
399:       removedMarkerLines[line]--
400:       code, severity = "QUALIFY-DEBUG-ARTIFACT-RELOCATED", resultmodel.SeverityWarning
401:       evidence += " — relocated, not added"
402:       stop = "the marker was relocated rather than introduced; inspect its retained intent"
403:     }
404:     findings = append(findings, helperFinding(code, severity, []string{path}, evidence, resultmodel.FixabilityManual, stop, nil, nil))
405:   }
```

### Which files it looks at

`qualificationChangedLines` (line 561):

- With `--diff-range`: `git diff --no-ext-diff --no-color --unified=0 <range>` — **added and removed lines only** (`-U0`).
- Without a range: the working diff **plus** `git diff --cached`, **plus** every untracked non-binary file from `git ls-files --others --exclude-standard`, whose entire contents are treated as added lines (lines 597-616).
- `do-work` and `do-work/**` are skipped for artifact scanning (line 387) and for the wiring check (line 328, 545, 551).

### What it does NOT check

- **Unchanged context.** `-U0` means a `TODO` that already existed on a line the diff did not touch is invisible.
- **`console.log` / `print(` as a "debug artifact".** Different code path, different codes, and the relocation and reporter carve-outs (lines 412-417) make them warnings in the common cases.
- **Commented-out code blocks**, dead branches, stray fixtures.
- **Temporary files left behind** — nothing in `handleQualify` looks for them. (`review-work.md:106` is the only place that rule exists.)
- **Whether the `[UNIFY]` claim is truthful.** The two P-A-U codes read only the REQ file; the artifact scan reads only the diff. Nothing correlates the two. The `work.md:574` row's claim — "the qualifier will catch a checked `[UNIFY]` over a diff containing `console.log`" — is true only by coincidence of both checks running, and only for the `QUALIFY-LIBRARY-OUTPUT` path, not `QUALIFY-DEBUG-ARTIFACT`.
- **Whether the comment carrying a marker documents a lesson worth keeping.**

### The explicit deferral to the agent — and a mis-citation in the REQ

Line 347:

```go
347: findings = append(findings, helperFinding("QUALIFY-NEW-FILE-UNWIRED", resultmodel.SeverityWarning, []string{path}, "new file has no static reference outside itself", resultmodel.FixabilityManual, "judge entry-point or dynamic-wiring exceptions", nil, nil))
```

That is the only "judge …" stop-reason in `handleQualify`, and **it belongs to `QUALIFY-NEW-FILE-UNWIRED` (dead-code detection), not to any of the three debug-artifact codes the REQ names.** The prose that carries this deferral is `work.md:335` ("conventions, dynamic imports, side-effect imports, and standalone entry points are not dead code merely because a grep cannot see their consumers"), which is **not one of the 7 sites** and is not in scope.

The debug-artifact trio's own deferrals are softer, carried in stop-reason strings rather than a "judge …" instruction: `"the marker was relocated rather than introduced; inspect its retained intent"` (line 402) and `"reporter output requires human judgment"` (line 415, `QUALIFY-REPORTER-OUTPUT`). **No prose site currently states either.** `review-work.md:106`'s "console.log/print lines no contract reads (a check's own reporting is contract output)" is the closest live restatement of the reporter carve-out — which argues for keeping it.

---

## 3. Per-site judgment

Classes: **(a)** pure restatement of what the Go check enforces; **(b)** judgment the check explicitly defers; **(c)** a reviewer-side read the check never sees.

| Site | File:line | Class | Assessment |
|---|---|---|---|
| S1 | `work.md:292` | **(a) mostly, with a live routing tail** | "verifies no debug artifacts" is a restatement, but the bullet's job is the P-A-U phasing contract, and the parenthetical "(the orchestrator audits this during Qualification and Testing Judgment)" is live routing. **Cut the four words "verifies no debug artifacts," and keep the bullet.** This is the cheapest of the 7 and does not need a pointer of its own. |
| S2 | `work.md:574` | **(a) — and factually wrong** | Pure restatement. Worse, the claim "a diff containing `console.log` … the qualifier will catch" describes `QUALIFY-LIBRARY-OUTPUT`, not the debug-artifact code, and `QUALIFY-LIBRARY-OUTPUT` downgrades to a warning when output was relocated or the file owns its own exit. Under prime-action-files §"The test, not a vibe" this row names no failure the code does not already stop. **Cut the row.** |
| S3 | `work.md:583` | **(a)** | Pure restatement of `QUALIFY-DEBUG-ARTIFACT` + `QUALIFY-PAU-UNCHECKED`, and it is the orchestrator's Red Flags list — the orchestrator is exactly who runs `qualify`. **Cut the bullet; replace with nothing, or fold the finding-code pointer here instead of into Step 6.3.** |
| S4 | `review-work.md:106` | **(c) — KEEP** | Three things the check never sees: (i) "temporary files left behind" has no Go equivalent at all; (ii) "print lines no contract reads (a check's own reporting is contract output)" is the human side of the `QUALIFY-REPORTER-OUTPUT` warning, which the code hands to judgment by design; (iii) "**Protect lessons learned**" is a judgment about comment intent the regex cannot form. Only the bare words "no debug artifacts" duplicate. **Keep the bullet; at most drop the three duplicated words.** |
| S5 | `review-work.md:374` | **template payload, not an instruction — KEEP** | This line is inside the fenced Review-Fix REQ template that gets written into a *new* REQ file. It is the same text as the P-A-U block in every queued REQ (including REQ-556's own lines 27-29). Cutting it here desynchronizes the emitted template from `capture.md`'s and from every REQ on disk, for zero reader-attention saving — the reader of `review-work.md` is not being instructed by it. Prime-action-files §Cross-Referencing already treats a fenced block's payload as content that "lands in some other file," not a citation from here; the same reasoning applies. **Do not count this as a prose site.** |
| S6 | `review-work.md:494` | **(a) if the loop ran qualify; (c) in standalone mode** | See §3b below. Weakest of the keeps: it is a bare restatement with none of S4's extra content. **Cut it** and let S4 carry the reviewer-side read. |
| S7 | `work-reference.md:604` | **(a)** | Pure restatement of `QUALIFY-PAU-UNCHECKED` + `QUALIFY-DEBUG-ARTIFACT`, in a table whose own heading is already stale ("Step 6.3"). **Cut the row.** |

### 3b. Does the review pass run on a diff `qualify` never saw? — Yes, in standalone mode.

Read from `review-work.md` Step 4 (line 70) and Step 6 (line 97):

- **Step 4, orchestrated mode** (line 72): "Run `git diff` … If the working tree is clean … use `git diff --staged`. **In worktree dispatch mode** … read the diff from the merge range the orchestrator passes: `git diff <pre>..<merge_hash>`". That is byte-for-byte the same input `qualify` consumes: `qualificationChangedLines` runs `git diff` + `git diff --cached` with no range (lines 573-580), or `git diff <range>` with one (line 565). **In orchestrated mode the reviewer sees exactly what `qualify` saw** — except that the reviewer reads full files ("New files created (read them fully)", line 77) while `qualify` reads `-U0` hunks only.
- **Step 4, standalone mode** (line 74): "Read the commit with `<skill-root>/scripts/show-commit-diff.sh <commit>`." **`qualify` never runs in standalone review at all** — `handleQualify` is invoked by the work loop's evidence-gate path, and standalone review is entered from `review-work.md` Step 1 against an already-archived REQ. So in standalone mode every debug-artifact read is a first read.
- Independently of mode, `qualify` skips `do-work/**` (checks.go:387) and skips unchanged context (`-U0`), and it has no notion of temporary files or of comment intent. Step 6's bullet covers all three.

**Conclusion:** the Builder Guidance's latitude is earned, and it is earned by exactly one sentence. **Keep `review-work.md:106` (S4).** Cut `review-work.md:494` (S6) — it duplicates S4's weakest half and adds nothing standalone mode needs. Leave `review-work.md:374` (S5) alone as template payload.

### 3c. Recommended end state

Cut S2, S3, S6, S7 outright; trim four words from S1; keep S4; leave S5 untouched.

Remaining mentions of the `rg` pattern after that: **S1 (if the four words go, S1 stops matching), S4, S5** → **2 matches** (S4 and S5), comfortably under the ≤ 3 ceiling. If S1's four words are kept, the count is 3 — still green. Expected net line delta roughly −6 to −8, against the REQ's forecast of −15 (the forecast assumed 9 sites, two of which are already gone).

**The one open judgment for the builder:** whether to add a sentence naming `QUALIFY-DEBUG-ARTIFACT` / `QUALIFY-PAU-UNCHECKED` / `QUALIFY-UNIFY-DISARMED`. The REQ calls this firm, but §1 shows it is an addition, not a retention, and `maintenance.md` § 3 demands a replay case for any addition. A pointer naming the command rather than the codes ("`advance`'s qualification gate audits the diff for unfinished markers and the P-A-U boxes — do not restate its rules here") is the cheaper form and matches the house style the REQ itself cites from `capture.md`.

---

## 4. The REQ-510 overlap

- **Location:** `do-work/archive/UR-098/REQ-510-sweep-work-reference-sections-owned-by-cli-tests.md`
- **Status:** `completed` (frontmatter line 4). `claimed_at: 2026-09-05T00:33:24Z`, route B.
- **Commits:** `6eeeee8d` "[REQ-510] sweep work-reference sections a CLI behavior test now owns" and `0e599236` "[REQ-510] restore two unowned condensations and close the stale references review found".
- **Did it touch the site?** **No.** The `## Qualification Anti-Rationalization Table (Step 6.3)` row survives at `work-reference.md:604`, verbatim as the REQ describes it: `| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |`.

**So the "skip it if REQ-510 already removed it" escape does not apply. REQ-556 owns that row.** Note that `0e599236`'s title says the review found REQ-510 had over-condensed and two sections were *restored* — a caution that this file's table rows have already survived one deletion pass on the grounds that the CLI does not own them. The builder should state, in the REQ trail, why this row differs: `QUALIFY-PAU-UNCHECKED` (an error-severity check on the REQ file) is a strictly stronger enforcement of "a checked box is a claim, not a fact" than the prose row is.

---

## 5. The lock-in — `_dev/tests/audit-lockins.sh`

**Path:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/tests/audit-lockins.sh` (executable, 143 lines).

### Structure

Header `#!/usr/bin/env bash`, comment `# Audit lock-in regressions. Pinned after maintainability audits.`, `set -uo pipefail` (note: **no `-e`**, deliberately — every block is expected to run to completion and accumulate). Then `repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"` and `failure_count=0`.

Each assertion is one block in the same shape:

1. a comment naming `# Finding N: <sweep-key> (REQ-NNN)`,
2. a `<name>="$( … )"` command substitution that emits zero or more offender lines,
3. `if [ -n "$<name>" ]; then while IFS= read -r … printf 'FAIL: …' >&2; failure_count=$((failure_count + 1)); done <<< "$<name>"; fi`.

Footer: `if [ "$failure_count" -gt 0 ]; then exit 1; fi` then `printf 'Audit lock-in regressions passed.\n'`.

Four blocks exist today, in order: Finding 10 exported one-line delegates (REQ-550), Finding 5 caller-less toolbox shims (REQ-551), the shipped-shell-delegating companion (REQ-551), Finding 8 dead-path-pointers-in-records (REQ-549), Finding 2 cli-launcher-preamble-copied (REQ-553).

### Two verbatim existing assertions

```bash
# Finding 5: toolbox-shims-no-callers (REQ-551)
callerless_shims="$(
  for f in "$repo_root"/skills/do-work-toolbox/scripts/*.sh; do
    [ -e "$f" ] || continue
    b=$(basename "$f")
    n=$(rg -l --fixed-strings "$b" "$repo_root/skills" "$repo_root/tools" "$repo_root/suite" "$repo_root/README.md" "$repo_root/CLAUDE.md" "$repo_root/_dev/primes" --glob '!*CHANGELOG*' | grep -v "/$b$" | wc -l | tr -d ' ')
    [ "$n" -eq 0 ] && echo "$f"
  done
)"
if [ -n "$callerless_shims" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: caller-less toolbox shell shim found: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$callerless_shims"
fi
```

```bash
# Finding 2: cli-launcher-preamble-copied (REQ-553)
# The two exempt paths are named in full, not by trailing path shape: a third copy at
# skills/do-work-toolbox/tools/do-work-cli-preamble.sh is exactly the hand-rolled preamble
# this pins at zero, and a suffix filter would read it as the exempt file.
hand_rolled_preambles="$(
  rg -l --glob '*.sh' 'for cli_candidate in|^launcher_arguments=\(--format text\)$' \
    "$repo_root/skills" "$repo_root/tools" 2>/dev/null \
    | grep -vxF -e "$repo_root/tools/do-work-cli-preamble.sh" \
                -e "$repo_root/skills/do-work/tools/do-work-cli-preamble.sh"
)"
if [ -n "$hand_rolled_preambles" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$hand_rolled_preambles"
fi
```

### Registration in `_dev/tests/contracts/probe-lanes.sh`

Line 29-30, in the unconditional (fast-tier) block, via the file's own `register_probe` helper (defined at line 5, which fails the run if the script is missing or not executable, then hands it to `launch_probe`):

```bash
register_probe audit_lockins_probe "$repo_root/_dev/tests/audit-lockins.sh" \
  'audit lock-in regressions failed (see the attributed FAIL lines above).'
```

`probe-lanes.sh` is sourced by `_dev/tests/contract-regressions.sh:73`, after `probe-batch.sh`, and `collect_probes` runs them. **No registration change is needed.**

### Paste-ready assertion, in the file's own style

Append before the `if [ "$failure_count" -gt 0 ]` footer. Written as a count so it reads as the ceiling the REQ names; the `|| true` guard is there because `rg` exits 1 on zero matches and the script runs under `pipefail`.

```bash
# Finding 1: qualify-debug-artifact-prose-restated (REQ-556)
# do-work-cli qualify owns the debug-artifact and P-A-U rules (QUALIFY-DEBUG-ARTIFACT,
# QUALIFY-PAU-UNCHECKED, QUALIFY-UNIFY-DISARMED in internal/corehelpers/checks.go), so the
# action files carry a pointer, not a second copy. Counted, not name-listed: the point is
# that the restatements do not regrow, and any new one is the regression whatever it says.
debug_rule_mention_ceiling=3
debug_rule_mentions="$(
  rg -c -e 'console\.log' -e 'debug artifacts' \
    "$repo_root/skills/do-work/actions/work.md" \
    "$repo_root/skills/do-work/actions/review-work.md" \
    "$repo_root/skills/do-work/actions/work-reference.md" 2>/dev/null \
    | awk -F: '{ total += $NF } END { print total + 0 }'
)"
if [ "${debug_rule_mentions:-0}" -gt "$debug_rule_mention_ceiling" ]; then
  printf 'FAIL: %s debug-artifact rule mentions across work.md, review-work.md and work-reference.md; ceiling is %s (do-work-cli qualify owns the rule)\n' \
    "$debug_rule_mentions" "$debug_rule_mention_ceiling" >&2
  failure_count=$((failure_count + 1))
fi
```

Trap to be aware of when pinning the number: `rg -c` prints one `path:count` line per file **and omits files with zero matches**, which is why the `awk -F:` sums the last field and defaults to `0`. Do not use `rg -c | wc -l` — that counts files, not mentions.

**Set `debug_rule_mention_ceiling` to whatever the post-cut count actually is** (2 or 3 under the §3c end state), not to a round number. The REQ's constraint says "pinned at today's value so it is green on day one and red the moment the number regrows".

---

## 6. `_dev/tests/contract-regressions.sh` — no predicate to delete

Greps run: `console\.log`, `debug artifact`, `UNIFY box`, `Anti-Rationalization`, `Red Flags`, and the three filenames `work.md` / `review-work.md` / `work-reference.md`. **Zero hits for all of them.** The file is 77 lines and is now a pure driver: a tier check, a self-size ratchet, `run_contract_file` over four owner contracts under `_dev/tests/contracts/`, then `probe-batch.sh` + `probe-lanes.sh` + `collect_probes`.

**So the REQ's constraint "delete the matching predicates in the same commit" has nothing to act on.** Say so in the hand-back rather than leaving it looking unaddressed.

**A hard blocker on touching that file anyway** (lines 9 and 17-22):

```bash
fast_contract_line_ceiling=77
…
if [ "$actual_contract_lines" -gt "$fast_contract_line_ceiling" ]; then
  printf 'FAIL: contract-regressions.sh grew to %s lines; ratchet ceiling is %s.\n' …
```

`wc -l` on the file returns exactly 77. **Adding a single line to `contract-regressions.sh` fails the suite.** The lock-in must go in `audit-lockins.sh`, as the REQ says.

### Other tests that touch these strings (checked, none block the cut)

- `_dev/tests/prescribed-shell-cases/qualify.sh` — 28 hits, all fixture repos and assertions against the **CLI's own output** (`grep -q 'debug artifacts'` against `qualify` stdout at lines 70, 179, 343). It pins the Go behavior, not the action prose. **Untouched by this REQ.**
- `_dev/tests/defensive-surface-audit.sh:50` — `assert_phrase_absent 'skills/do-work-toolbox/actions/inspect.md' 'Debug artifacts (console.log, debugger, commented-out blocks)'`. This is an **absence** assertion on a *different* file (`do-work-toolbox/actions/inspect.md`), pinning a prior deletion. It does not pin any of the 7 sites and needs no change.

---

## 7. Prime and maintenance rules

### `_dev/primes/prime-action-files.md` (120 lines) — the entries that govern this cut

**What may be cut, and the test for it (line 72):**

> **Earned, not mandatory: Rules, Common Rationalizations, Red Flags, Verification Checklist.** Add one only when the file has something a capable model would otherwise get wrong — do-work machinery (a queue/pipeline mechanic, a frontmatter or schema contract) or a hard-won failure mode with a traceable origin (a real REQ or incident this stops from recurring). "This is generic engineering advice a capable model already follows" is an explicit *non*-reason — true or not, it doesn't earn a section.

Line 74:

> **The test, not a vibe:** before adding a Common Rationalizations row, ask *can I name the specific failure this row prevents, and where it happened?* No → don't add the row. If every row in a table fails that test, omit the whole section — a generic table is worse than no table: it teaches the reader the section is decorative, so they stop reading the ones that aren't. Apply the same test to Rules and Red Flags — specific to this action, not restated hygiene ("write tests," "don't skip validation"). When a file has nothing that passes, omit the section entirely; don't ship it empty or generic to satisfy the template.

This is the direct warrant for cutting S2 (a Common Rationalizations row), S3 and S6 (Red Flags bullets), and S7 (an Anti-Rationalization row): each restates a rule an error-severity check already enforces, so none names a failure the row prevents.

Line 76 — what a pointer should look like:

> **State intent, not a directive rule, when a capable model can infer the rest.** "Report drift, don't fix it inline" gives the model this action's boundary in one line — a five-line Rules section re-deriving why inline fixes are bad adds nothing a capable model didn't already know.

**What may never be cut (line 70):**

> **Required:** Description blockquote, Steps (numbered). **Common:** Input, Output Format, When to Use.

None of the 7 sites is a description blockquote or a numbered step heading, so nothing here is protected. Section order (line 80) is unaffected — no section is being removed wholesale.

**Cross-referencing rules that constrain any pointer written (line 91, 93):** same-package actions by local path (`actions/work.md`); a cross-package citation uses the literal relative path from the citing file's own directory, and any citation whose first segment names a sibling package is checked by `_dev/tests/shipped-package-reference-contract.sh` whether backticked or not. A pointer to the Go source should therefore be spelled as a same-package tools path from `actions/`, e.g. `../tools/do-work-cli/internal/corehelpers/checks.go`, or better as a behavior pointer to the command, avoiding a path that a refactor can dangle.

**Relevant trap (line 107):**

> `[family: alternate-writer-contract-drift]` Changing an emitted artifact or routing rule only at its primary writer leaves alternate modes and action-bearing readers silently following the old contract; grep every writer and downstream reader against the recorded scope before declaring it shipped — including sibling-package actions, board labels, and prime sentences that restate the rule (REQ-566 found six such restatements outside the three files that defined the contract).

Applies directly: the builder should re-grep after the cut, including `skills/do-work-toolbox/` and `skills/do-work-board/`, for further restatements outside the three declared files. (The lock-in as written only watches three files — that is a deliberate scoping choice matching the REQ's ceiling, not an oversight.)

**Lessons satellite (line 120):** `_dev/primes/lessons-action-files.md` — read before changing what **Read first** or **Traps** name. REQ-556's frontmatter records it as dropped for budget (3968 tokens vs a 2000 ceiling); since this REQ does not change Traps, the drop is fine, but say so in the hand-back.

### `skills/do-work/crew-members/maintenance.md` — "The Subtractor" — what it demands

Auto-loads because REQ-556 carries `maintenance: true`.

1. **Delete before you add** (§1) — subtraction is the first move; "a new rule is the most expensive fix: it is permanent, it competes for attention with every other rule, and it compounds."
2. **Ask the deletion questions** (§2) — is a stale source feeding the drift, is a bad example teaching it, is a tool too broad, is the job too vague? A "yes" to any is a fix by removal. Here the answer is "stale source": five prose copies of a rule the Go check owns. That is the diagnosis to record.
3. **Prove any addition against a replay pack** (§3) — "An addition is justified only when a concrete case **fails without it and passes with it** … No replay case, no addition." **This is the rule that bites the REQ's "keep one sentence naming the three finding codes" instruction**, since §1 of this report shows no such sentence exists today, so writing one is an addition.
4. **Boundaries** — "Subtraction is not vandalism … if removing something would have to be restored next week, it was foundation, not bloat … **When in doubt, narrow rather than delete — and record the call.**" That is the warrant for narrowing S1 (four words) and S4 (three words) instead of cutting the whole bullets, and the reason S4 stays.
5. Loaded **alongside** general / coding-guardrails / anti-slop, not instead of them.

---

## Summary of what the builder must decide

- **D1.** Re-baseline the REQ from 9 to 7 in its trail before starting. Two work.md sites are already gone (REQ-504→506 advance chain), so the −15 line forecast is now roughly −6 to −8.
- **D2.** "work.md Step 6.3" no longer exists; the live heading is `### Qualification and Testing Judgment` (work.md:333).
- **D3.** "Keep one sentence naming the three finding codes" is an **addition**, not a retention — no finding code appears in any shipped file today. `maintenance.md` §3 wants a replay case, or a command-level pointer instead of a code list.
- **D4.** The "judge entry-point or dynamic-wiring exceptions" deferral the REQ cites belongs to `QUALIFY-NEW-FILE-UNWIRED`, a different check, and its prose home (work.md:335) is not one of the sites. The debug-artifact trio's real deferrals are the relocation and reporter-output carve-outs — currently only half-stated, at `review-work.md:106`.
- **D5.** Keep `review-work.md:106`; standalone review runs `qualify` never made, and the bullet's temporary-files and protect-lessons clauses have no code equivalent at all.
- **D6.** `review-work.md:374` is emitted template payload, not reader instruction — leave it, and say why.
- **D7.** REQ-510 is `completed` and did **not** remove `work-reference.md:604`. REQ-556 owns it.
- **D8.** Nothing in `contract-regressions.sh` pins these sentences, and its 77-line self-ratchet is already at the ceiling — do not add a line there.
