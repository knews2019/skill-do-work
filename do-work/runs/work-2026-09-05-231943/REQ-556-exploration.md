# REQ-556 exploration — cut the debug-artifact rule prose that `do-work-cli qualify` already enforces

REQ: `/home/user/skill-do-work/do-work/queue/REQ-556-cut-the-debug-artifact-rule-prose-that-qualify-already-enforces.md`
HEAD at exploration time: `fce57fccb19338491fea9d01bb0721a71f6d988b` (clean tree).
Prior exploration re-checked: `do-work/runs/work-2026-09-05-170806/REQ-556-exploration.md` (written at `71eb49f3`, 63 commits back).

---

## 1. The count at HEAD, and how stale the REQ's baseline is

The REQ's Red-Green Proof says the Reproduce command prints **nine** mentions (work.md 5, review-work.md 3, work-reference.md 1). It prints **seven**.

```
rg -n -e 'console\.log' -e 'debug artifacts' \
  skills/do-work/actions/work.md \
  skills/do-work/actions/work-reference.md \
  skills/do-work/actions/review-work.md
```

| File | REQ claim | At HEAD | Delta |
|---|---|---|---|
| `skills/do-work/actions/work.md` | 5 | **3** | −2 |
| `skills/do-work/actions/review-work.md` | 3 | **3** | 0 |
| `skills/do-work/actions/work-reference.md` | 1 | **1** | 0 |
| **total** | **9** | **7** | **−2** |

Per-file `grep -c -e 'console\.log' -e 'debug artifacts'` returns 3 / 3 / 1.

**The 9-count cannot be reproduced in this repository at all.** The REQ names audited commit `dc8a64e3`; `git cat-file -t dc8a64e3` returns `fatal: Not a valid object name`. This clone holds 99 commits and does not contain it. Only the HEAD count of 7 is verifiable evidence.

### The seven sites, verbatim

| Code | Anchor | Text |
|---|---|---|
| **S1** | `skills/do-work/actions/work.md:292` | `- **P-A-U phasing is mandatory:** ... [UNIFY] runs \`git diff --stat\`, runs native linters, verifies no debug artifacts, and lists each file checked (the orchestrator audits this during Qualification and Testing Judgment).` |
| **S2** | `skills/do-work/actions/work.md:574` | `\| "P-A-U is bookkeeping — I'll just tick the boxes" \| Do each phase; qualification audits the diff against the checked boxes \| A checked \`[UNIFY]\` over a diff containing \`console.log\` is a false claim the qualifier will catch \|` |
| **S3** | `skills/do-work/actions/work.md:583` | `- All P-A-U checkboxes marked complete but diff contains \`console.log\`, \`debugger\`, or \`TODO\` (debug artifacts)` |
| **S4** | `skills/do-work/actions/review-work.md:106` | `- Diff hygiene — no debug artifacts — console.log/print lines no contract reads (a check's own reporting is contract output), or temporary files left behind. **Protect lessons learned** — ...` |
| **S5** | `skills/do-work/actions/review-work.md:374` | `- [ ] **[UNIFY]:** (Agent: Run \`git diff --stat\` ... Verify no debug artifacts in diff. ...)` |
| **S6** | `skills/do-work/actions/review-work.md:494` | `- Builder checked all P-A-U boxes but the diff contains \`console.log\`, \`debugger\`, or TODO/FIXME` |
| **S7** | `skills/do-work/actions/work-reference.md:604` | `\| "The builder checked the UNIFY box" \| Read the actual diff for debug artifacts \| A checked box is a claim, not a fact \|` |

### What holds and what does not, in the prior exploration

`git diff --stat 71eb49f3..HEAD` over the three action files, `_dev/tests/audit-lockins.sh`, `_dev/tests/contract-regressions.sh`, `checks.go` and `_dev/primes/prime-action-files.md` is **empty** — 63 commits of churn, none of it touching this REQ's surface. So the prior exploration's recount of 7 and every one of its line numbers still hold at HEAD, and so do its `checks.go` anchors (24, 261, 347, 353-360, 377, 386-405, 561).

Statements in it that **do not** hold:

- **C1.** "`_dev/tests/audit-lockins.sh` (executable, 143 lines)". It is **156 lines**, at HEAD and at `71eb49f3` alike.
- **C2.** "lines 422 and 433 still cite 'Step 6.3' / 'Step 6.25' step numbers that work.md no longer has". `### Step 6.25: Implementation Summary` **does** exist at `work.md:314`. Only Step 6.3 is gone.
- **C3.** It calls the work-reference table heading "already stale" while recommending the row be cut — but missed that the heading is **mechanically pinned** (see §4).
- **C4.** Its paths are from another machine (`/Users/t2/Desktop/...`).
- **C5.** Five of the commit hashes it cites are not objects here: `dc8a64e3`, `0e1b44ff`, `6eeeee8d`, `0e599236`, `83594c5e`. `a296ee9e` and `59fe3e3a` exist but are **not ancestors of HEAD**, so its account of which chain removed the two vanished work.md sites is unverifiable on this branch. Its site judgments stand on their own evidence; its provenance claims should not be repeated as fact.
- **C6.** Its paste-ready lock-in snippet has a masking-fallback defect (see §5).

---

## 2. Retention or addition? — settled: **addition**, and it is earned

```
rg -n 'QUALIFY-' skills/ --glob '!*.go'      ->  zero hits
```

`QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and `QUALIFY-UNIFY-DISARMED` appear in **no shipped file**. The REQ's "Keep one sentence in `work.md` Step 6.3 naming the three finding codes" describes prose that does not exist. It is an addition, exactly as the prior exploration said.

The prior exploration then concluded the addition should be dropped for a command-level pointer, on `skills/do-work/crew-members/maintenance.md:26` ("No replay case, no addition"). **That conclusion does not hold — the replay case exists:**

- `work.md:439` and `work.md:473` invoke the CLI with `--format json`.
- `checks.go:451-473` is the **text** renderer; it converts codes to prose (`FAIL: debug artifacts in %s`). In JSON the raw `finding.Code` is what the agent reads.
- So an orchestrator working the Mechanical Evidence-Gate Loop (`work.md:165`, "judge qualification warnings") sees the literal token `QUALIFY-DEBUG-ARTIFACT` and today has **zero** shipped prose routing it to an owner.

Fails without the sentence, passes with it. That satisfies `maintenance.md:26` in substance, and it is also what the maintainer marked firm in Builder Guidance. **Write the sentence.**

Two hard constraints on where and how:

- **Not under a "Step 6.3" heading.** `_dev/tests/contracts/core-checks.sh:701-706` FAILs if `### Step 6.3:` is restored to `work.md`. The live home is `### Qualification and Testing Judgment` (`work.md:333`).
- **The REQ's code-to-rule mapping is wrong twice**, so do not copy it into the sentence:
  - `checks.go:24` matches only `\b(debugger|TODO|FIXME)\b`. `console.log` and `print(` go through a separate pattern at `checks.go:377` producing `QUALIFY-LIBRARY-OUTPUT` (410), `QUALIFY-OUTPUT-RELOCATED` (413) or `QUALIFY-REPORTER-OUTPUT` (415). Every prose site that equates `console.log` with "debug artifact" is already slightly wrong about the taxonomy.
  - The "judge entry-point or dynamic-wiring exceptions" deferral the REQ tells the builder to preserve belongs to **`QUALIFY-NEW-FILE-UNWIRED`** at `checks.go:347`, a different check. Its prose home is `work.md:335`, which is **not one of the seven sites** and is not in scope.

### The sentence to add

Append to the paragraph under `### Qualification and Testing Judgment` that currently ends at `work.md:335`:

> `advance`'s qualification gate owns these mechanics and reports them as typed findings — `QUALIFY-DEBUG-ARTIFACT` (unfinished-work markers added by the diff), `QUALIFY-PAU-UNCHECKED` (a P-A-U box left open), `QUALIFY-UNIFY-DISARMED` (no `[UNIFY]` box at all), plus a separate pair for leftover output primitives whose reporter-output case it hands back to you — so judge its findings rather than restating the rule here.

It deliberately contains neither `console.log` nor `debug artifacts`, so it consumes no ceiling headroom, and neither `TODO`, `FIXME` nor `debugger` as whole words, so this REQ's own qualify gate stays green (`checks.go:392` scans **added** lines).

---

## 3. Per-site verdicts

What `qualify` actually owns, from `handleQualify` (`checks.go:261`):

- `QUALIFY-UNIFY-DISARMED` (`checks.go:353, 356`) — warning; reads **the REQ file only**; fires when no `[UNIFY]` line exists.
- `QUALIFY-PAU-UNCHECKED` (`checks.go:354, 359`) — error; REQ file only; counts literal `- [ ]` boxes. Pure bookkeeping; says nothing about the diff.
- `QUALIFY-DEBUG-ARTIFACT` (`checks.go:386-405`) — error on added lines matching `\b(debugger|TODO|FIXME)\b`, downgraded to `-RELOCATED` warning when the identical line was also removed.

What it never sees: unchanged context (`qualificationChangedLines`, `checks.go:561`, runs `-U0`); everything under `do-work/` (`checks.go:387`); temporary files left behind (no code path anywhere); whether a marker's comment documents a lesson; and whether the `[UNIFY]` claim is truthful — the P-A-U codes read the REQ file and the artifact scan reads the diff, and nothing correlates them.

| Site | Verdict | Why |
|---|---|---|
| **S1** `work.md:292` | **Trim four words** | The bullet's job is the P-A-U phasing contract and the routing tail; its preamble at `work.md:285` already says this block is "pointers — the underlying rules live in the loaded crew-members files". The real UNIFY instruction ships in the template payload inside every REQ file, so "verifies no debug artifacts," here is a third copy. Keep the bullet. |
| **S2** `work.md:574` | **Cut the row** | `prime-action-files.md:74` — "can I name the specific failure this row prevents, and where it happened?" It cannot: the claim "the qualifier will catch a checked `[UNIFY]` over a diff containing `console.log`" describes `QUALIFY-LIBRARY-OUTPUT`, not the debug-artifact code, and that finding downgrades to a warning when output was relocated or the file owns its own exit (`checks.go:413, 415`). |
| **S3** `work.md:583` | **Cut the bullet** | Pure restatement of an error-severity check, in the orchestrator's own Red Flags list — the orchestrator is exactly who runs `qualify`. `prime-action-files.md:72` calls "generic engineering advice a capable model already follows" an explicit non-reason for a section. |
| **S4** `review-work.md:106` | **KEEP unchanged** | Three things the check never sees: temporary files left behind (no Go equivalent at all); "print lines no contract reads (a check's own reporting is contract output)", which is the human side of the `QUALIFY-REPORTER-OUTPUT` warning the code hands to judgment by design (`checks.go:415`); and **Protect lessons learned**, a judgment about comment intent no regex can form. |
| **S5** `review-work.md:374` | **KEEP unchanged — not a prose site** | Inside the ```` ```markdown ```` fence opened at `review-work.md:355` under `### Step 10: Create Follow-up REQs` (line 340). The identical sentence ships in **four** files: `review-work.md:374`, `capture-reference.md:104`, `do-work-toolbox/actions/code-review.md:312`, `sample-archived-req.md:33`. `prime-action-files.md:91` already treats a fenced block's payload as content landing elsewhere, not a citation from here. Cutting one copy desynchronizes the emitted template from the other three and from every REQ on disk, for zero reader-attention saving. |
| **S6** `review-work.md:494` | **Cut the bullet** | Duplicates S4's weakest half and carries none of its extra content. |
| **S7** `work-reference.md:604` | **Cut the row only** | Pure restatement of `QUALIFY-PAU-UNCHECKED` + `QUALIFY-DEBUG-ARTIFACT`. Five rows remain (`work-reference.md:602-607`). **The heading at `work-reference.md:598` must not be touched** — see §4. |

### Why the reviewer-side keep is earned

`review-work.md:74` (standalone mode) reads a historical commit through `scripts/show-commit-diff.sh`; `qualify` never runs in that path. Independently, `grep -n 'qualify\|advance' skills/do-work/actions/review-work.md` returns **zero hits** — the file nowhere assumes qualification happened, so a pointer written into it would have to introduce a concept the file otherwise does not carry. In orchestrated mode (`review-work.md:72`) the reviewer reads the same `<pre>..<merge_hash>` range `qualify` consumes, but reads **full files** (`review-work.md:77`) where `qualify` reads `-U0` hunks. Keeping exactly one reviewer sentence, and cutting the other, is the right cut line.

**End state:** cut S2, S3, S6, S7; trim S1; keep S4 and S5; add one sentence at `work.md:335`. Remaining matches: **2** (S4, S5). Net line delta roughly −6, not the REQ's forecast −15 — that forecast assumed 9 sites.

---

## 4. Tests that constrain the cut

Two pins the prior exploration missed, both in `_dev/tests/contracts/core-checks.sh`:

```
core-checks.sh:701   '### Step 6.3:'   -> work.md must NOT restore this retired heading
core-checks.sh:718   grep -qF '## Qualification Anti-Rationalization Table (Step 6.3)' "$work_reference"
core-checks.sh:721   FAIL: work-reference.md must retain the judgment-owned qualification table ...
```

So `work-reference.md:598`'s heading — stale step number and all — is deliberately retained by contract. **Cut the row, never the heading.**

Everything else checked and clear:

- **Common Rationalizations contract** (`core-checks.sh:788-816`) only counts near-identical cross-file row pairs; deleting a row can only reduce that count. `prime-action-files.md:78`'s do-work-noun requirement is satisfied by the surviving rows (queue, archive, REQ).
- **`_dev/tests/contract-regressions.sh` has nothing to delete.** grep for `console\.log`, `debug artifact`, `UNIFY`, `Anti-Rationalization`, `Red Flags`, and the three filenames returns **0 for every one**. The REQ's constraint "delete the matching predicates in the same commit" has no referent — say so in the hand-back. And do not add a line there: `contract-regressions.sh:9` sets `fast_contract_line_ceiling=77` and `wc -l` on the file is exactly **77**, so one added line fails the suite (`contract-regressions.sh:17-21`).
- **`_dev/tests/defensive-surface-audit.sh:50`** is an `assert_phrase_absent` on `do-work-toolbox/actions/inspect.md` for `'Debug artifacts (console.log, debugger, commented-out blocks)'` — a *different* file, pinning a prior deletion. It needs no change, and it is useful precedent: this repo has already deleted an equivalent debug-artifact Red Flags entry from an action file and pinned it absent.
- **`_dev/tests/prescribed-shell-cases/qualify.sh`** asserts against the CLI's own stdout, not the action prose. Untouched.
- No test anywhere pins any of the seven sentences (grep for `Diff hygiene`, `Protect lessons learned`, `P-A-U phasing is mandatory`, `Read the actual diff` across `_dev/tests/**/*.sh`).

---

## 5. The lock-in

`_dev/tests/audit-lockins.sh` — 156 lines, `-rwxr-xr-x`, `set -uo pipefail` (no `-e`, deliberately: blocks accumulate). Registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh:29`; **no registration change needed**. Baseline: `bash _dev/tests/audit-lockins.sh` → `Audit lock-in regressions passed.`, exit 0. Four blocks exist (`# Finding 10` at line 8, `# Finding 5` at 43, the REQ-551 companion at 60, `# Finding 8` at 74, `# Finding 2` at 133); the footer is at lines 151-156.

### Do not use the prior exploration's snippet as written

It computed the count as `rg -c … 2>/dev/null | awk -F: '{ total += $NF } END { print total + 0 }'`. If `rg` is absent, or any of the three files is renamed, that yields **0** — the GREEN branch. `_dev/primes/prime-shell-commands.md:31-45` ("Unchecked Exit Status Reads as Content") names exactly this defect: *a value that a later guard makes a safety decision on, where the collapsed value silently satisfies the safe branch.* A lock-in that dies silently on a rename is worse than none.

### Paste-ready assertion

Append before the footer at `_dev/tests/audit-lockins.sh:151`. Names follow `coding-guardrails.md` § 5 Naming for Reach (two words minimum, plain-text findable) and match the file's existing `snake_case` shell style.

```bash
# Finding 1: qualify-debug-artifact-prose-restated (REQ-556)
# do-work-cli qualify owns the debug-artifact and P-A-U-honesty rule
# (QUALIFY-DEBUG-ARTIFACT, QUALIFY-PAU-UNCHECKED, QUALIFY-UNIFY-DISARMED in
# skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go), so the action files
# carry one pointer instead of a second copy of the rule. Counted, not name-listed: a new
# restatement is the regression whatever words it uses. The two mentions the ceiling allows
# are review-work.md's standalone-review hygiene bullet (a read qualify never makes) and the
# emitted P-A-U template payload, which is byte-identical in four shipped files -- neither is
# a restatement, so do not chase the count below the ceiling by cutting them.
# A missing target file FAILs instead of counting zero: a rename must not retire this lock-in
# silently (_dev/primes/prime-shell-commands.md -> Unchecked Exit Status Reads as Content).
debug_rule_mention_ceiling=2
debug_rule_scanned_files=(
  "$repo_root/skills/do-work/actions/work.md"
  "$repo_root/skills/do-work/actions/review-work.md"
  "$repo_root/skills/do-work/actions/work-reference.md"
)
debug_rule_mention_count=0
for debug_rule_file in "${debug_rule_scanned_files[@]}"; do
  if [ ! -f "$debug_rule_file" ]; then
    printf 'FAIL: debug-artifact prose lock-in cannot read %s; the file moved and the lock-in is dead\n' \
      "$debug_rule_file" >&2
    failure_count=$((failure_count + 1))
    continue
  fi
  debug_rule_file_hits=$(grep -c -e 'console\.log' -e 'debug artifacts' "$debug_rule_file")
  debug_rule_mention_count=$((debug_rule_mention_count + debug_rule_file_hits))
done
if [ "$debug_rule_mention_count" -gt "$debug_rule_mention_ceiling" ]; then
  printf 'FAIL: %s debug-artifact rule mentions across work.md, review-work.md and work-reference.md; ceiling is %s (do-work-cli qualify owns the rule)\n' \
    "$debug_rule_mention_count" "$debug_rule_mention_ceiling" >&2
  failure_count=$((failure_count + 1))
fi
```

Executed inline against HEAD, it printed `count=7` and both the ceiling FAIL and (with a deliberately bad path added) the missing-file FAIL. `grep -c` counts matching **lines**, not matches — the same unit the `rg -n` listing shows — and its exit 1 on no-match is harmless under `set -uo pipefail` without `-e` because the assignment still yields `0`, with the existence guard above covering the case that actually matters.

Set `debug_rule_mention_ceiling=2`. If the builder decides to keep S1's four words, it is `3` — still inside the REQ's ≤ 3 constraint, but with no headroom left.

---

## 6. REQ-510 overlap

`do-work/archive/UR-098/REQ-510-sweep-work-reference-sections-owned-by-cli-tests.md:4` → `status: completed`, `claimed_at: 2026-09-05T00:33:24Z`, route B. `work-reference.md:604` **still carries the row**, verbatim. The REQ's "skip it if REQ-510 already removed it" escape does not apply: **REQ-556 owns that row.**

Worth recording in the trail why this row goes when REQ-510's pass deliberately left it: `QUALIFY-PAU-UNCHECKED` is an error-severity check on the REQ file, a strictly stronger enforcement of "a checked box is a claim, not a fact" than the prose row is.

---

## 7. Scope notes the REQ does not cover

- **This is a release.** `_dev/primes/prime-releases.md:5` — a commit changing shipped files under `skills/` is a release. It needs a `CHANGELOG.md` entry, a `skills/do-work/VERSION` bump (currently `0.303.10`), and the byte-identical mirror to `skills/do-work/CHANGELOG.md` (`prime-releases.md:10`, enforced by `shipped-package-reference-contract.sh`). None of those paths is in the REQ's `write_set`. Orchestrator Step 9 owns it; the builder should name it in the hand-back rather than write it.
- **Re-grep after the cut** for restatements outside the three declared files (`prime-action-files.md:107`, `[family: alternate-writer-contract-drift]`). Done here: the only shipped hits outside the three files are the four-way template payload, `do-work-toolbox` quick-wins/inspect/debugging material (different actions, different contract), and `sample-archived-req.md:33`. None is in scope.
- **Dropped lesson satellite.** The REQ's frontmatter records `_dev/primes/lessons-action-files.md` dropped for budget (3968 tokens vs 2000). This change does not touch **Read first** or **Traps** in `prime-action-files.md`, so the drop is fine — say so in the hand-back.
- **Discovered task (do not fix here):** `work-reference.md:422` and `:433` still cite a "Step 6.3" step number `work.md` no longer has, and the table heading at `:598` carries it too. The heading is contract-pinned so it cannot move without changing `core-checks.sh:718` in the same commit — a separate REQ, not this one.

---

## 8. Decisions for the builder

- **D1.** Re-baseline the REQ from 9 to 7 before starting, and note that the audited commit `dc8a64e3` is not in this repository so the captured RED cannot be replayed. Expected line delta ≈ −6, not −15. Do not manufacture old lines to obtain the captured RED.
- **D2.** "work.md Step 6.3" does not exist and must not be recreated (`core-checks.sh:701`). Put the sentence in `### Qualification and Testing Judgment` (`work.md:333`).
- **D3.** The three-code sentence is an **addition**, and it is earned: `advance` runs with `--format json` (`work.md:439`), so the agent reads the raw codes, and no shipped file maps any of them to an owner. Record that replay case in the REQ trail as the `maintenance.md:26` justification.
- **D4.** Correct the REQ's mapping in the sentence you write: `console.log`/`print(` are `QUALIFY-LIBRARY-OUTPUT` / `-OUTPUT-RELOCATED` / `-REPORTER-OUTPUT` (`checks.go:377, 410-415`), and "judge entry-point or dynamic-wiring exceptions" belongs to `QUALIFY-NEW-FILE-UNWIRED` (`checks.go:347`), whose prose home `work.md:335` is out of scope.
- **D5.** Keep `review-work.md:106`. Standalone review (`review-work.md:74`) runs on a diff `qualify` never saw, the file never mentions `qualify` at all, and temporary-files and comment-intent have no Go equivalent.
- **D6.** Leave `review-work.md:374` alone — emitted template payload, identical in four shipped files.
- **D7.** REQ-510 is completed and did not remove `work-reference.md:604`. Cut the row, keep the heading at `:598`.
- **D8.** Nothing in `contract-regressions.sh` pins these sentences, and its 77-line self-ratchet is at the ceiling. The lock-in goes in `audit-lockins.sh`, as the REQ says.
- **D9.** No added line in this diff may contain `TODO`, `FIXME` or `debugger` as a whole word — `checks.go:392` scans added lines only, so the deletions are safe but a careless replacement would turn this REQ's own gate red.

---
*Generated by Explore agent*