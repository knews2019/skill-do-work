# Pre-explore: REQ-554 — move the commit.md/inspect.md shared body into the prescribed-shell guide

Repo at HEAD `71eb49f3` ("[work run] checkpoint after REQ-588"). Read-only; nothing edited.

---

## 1. Claim re-verified at HEAD — both numbers reproduce exactly

```
python3 -c "import difflib;a=[l.rstrip() for l in open('skills/do-work/actions/commit.md')];b=[l.rstrip() for l in open('skills/do-work-toolbox/actions/inspect.md')];print(sum(1 for i,j,s in difflib.SequenceMatcher(None,a,b).get_matching_blocks() if s>=3 for k in range(s) if a[i+k].strip()))"
→ 46
rg -n 'If the script is missing or will not run' skills --glob '*.md'
→ 4 sites: commit.md:79, commit.md:113, inspect.md:88, inspect.md:145
```

The REQ was captured against `dc8a64e3`; the claim is unchanged at HEAD. File sizes now: `commit.md` 260 lines, `inspect.md` 448 lines, `prescribed-shell-primitives.md` 119 lines.

### F1 (BLOCKER) — the lock-in limit of ≤10 is mathematically unreachable

**20 of the 46 shared lines are template-mandated scaffolding, not prose.** They are section headings, code fences, and table header rows that `_dev/primes/prime-action-files.md` § Template *requires* both action files to have. Moving prose cannot remove them; only diverging the two files' required section structure could, and that would violate the template.

I measured the floor by deleting *every* shared content line from both files and re-running the REQ's own metric:

| Scenario | difflib count |
|---|---|
| Today (HEAD) | **46** |
| Move exactly the three classes the REQ names (legend + 4 bullets + both fallbacks) | **32** |
| Move *every* shared content line, including the ones the REQ never names | **17** |

The floor is 17, not 10, and the REQ's own scope only reaches 32. Whoever builds this must renegotiate the number before writing the assertion. The residue at the floor is exactly this, and every line is required by the action-file template:

```
## When to Use | **Use when:**
## When This Runs
## Steps | ```
``` | ### Step 1: Preflight
### Step 4: Group Unassociated Files
### Step 6: Report
``` | ## Error Handling | | Situation | Action | | |-----------|--------|
## What This Action Does NOT Do
## Rules
## Red Flags
## Verification Checklist
```

**Recommendation (D1):** keep the REQ's exact difflib metric so the number stays comparable to 46, and pin `≤ 20`. That is achievable by moving the REQ's three named classes *plus* the unnamed content listed in §2.5, leaves headroom of 3 over the measured floor of 17, and still catches any regrowth of real prose. Pinning at 17 exactly would be brittle: any future template change to either action's heading set moves it. Do **not** invent a filtered metric — a metric that strips headings is a second definition of "shared", and the REQ's Reproduce command is the one the audit and the report both quote.

`fallback sentences in actions: 0` is achievable as written and needs no renegotiation.

---

## 2. The exact shared blocks

Every run of 3+ byte-identical lines (span includes blanks; "content" counts non-blank, non-scaffold lines). 17 runs, 46 non-blank lines total = **20 scaffold + 26 content**.

### 2.1 Class A — the M/A/D/X/XD legend (REQ names it) — 5 content lines

| Block | commit.md | inspect.md | Content |
|---|---|---|---|
| B6 | 57–63 | 65–71 | `It gates on \`git rev-parse --git-dir\`, enumerates every uncommitted path, and prints one \`<tag>\t<path>\` row per file:` + the `- **M**`, `- **A**`, `- **D**` bullets (4 content lines; the 5th non-blank in the run is a ``` fence) |
| B7 | 67 | 75 | `Secret-shaped matching is case-insensitive and applies to the basename only: …` — byte-identical, 1 line |

**What differs (must stay local):** the `X` and `XD` bullets are *not* in the matching block because they are reworded for read-only mode.

```
commit.md:64  - **X** — … an ambiguous addition beside `X`/`XD`; fully excluded
inspect.md:72 - **X** — … an ambiguous addition beside `X`/`XD`; fully excluded from analysis

commit.md:65  - **XD** — deleted secret-shaped path; eligible for a deletion-only commit
inspect.md:73 - **XD** — deleted secret-shaped path; path/deletion state only
```

### 2.2 Class B — the four file-reading bullets (REQ names it) — 5 content lines

| Block | commit.md | inspect.md | Content |
|---|---|---|---|
| B8 | 80–89 | 89–98 | `Build a semantic understanding of each uncommitted file:` + `- **Modified files**`, `- **New/untracked files**`, `- **Deleted files (\`D\`)**`, `- **Deleted secret-shaped files (\`XD\`)**` |

All five byte-identical. The 6th non-blank in the run is the `### Step 2: Read Changes` heading (scaffold). Nothing differs between the two copies.

### 2.3 Class C — the two "do it by hand" fallbacks (REQ names it)

**Fallback 1 — the 207-word inventory algorithm.** `commit.md:79` / `inspect.md:88`. **This line contributes 0 to the 46** — the two copies differ, so it breaks the matching run. It is caught only by the `rg` count. Confirmed exactly two deltas, both relative-path fixups, exactly as the REQ states:

```
commit.md:79   [Per-file untracked inventory](../docs/prescribed-shell-primitives.md#per-file-untracked-inventory)
inspect.md:88  [Per-file untracked inventory](../../do-work/docs/prescribed-shell-primitives.md#per-file-untracked-inventory)

commit.md:79   … and `../../do-work-toolbox/actions/stray-check.md`'s Red Flags record that it has been hit.
inspect.md:88  … and `actions/stray-check.md`'s Red Flags record that it has been hit.
```

Both fixups resolve *identically* from the guide's own directory (`skills/do-work/docs/`), because that directory is two levels below `skills/` just like `skills/do-work/actions/`. So the moved copy takes commit.md's spellings, with the guide link collapsing to an in-file anchor:
- `[Per-file untracked inventory](../docs/…#per-file-untracked-inventory)` → `[Per-file untracked inventory](#per-file-untracked-inventory)` (or drop the link, since the text would then live inside that very section)
- `` `../../do-work-toolbox/actions/stray-check.md` `` — keep verbatim, correct from the guide.

**Fallback 2 — the association fallback.** `commit.md:113` / `inspect.md:145`. **Byte-identical** (verified with `diff`), inside block B11. Contributes 1 line to the 46.

```
If the script is missing or will not run, do it by hand: glob both directories, read each REQ's `status` (accepting every alias above) and `## Implementation Summary` list, path-match, and tie-break on the latest `completed_at`.
```

Note the phrase "accepting every alias above" — it back-references the "What the script settles" bullets. If the sentence moves to the guide alone, that reference dangles. The guide section must either carry the alias list or reword to `accepting the Schema Read Contract's terminal-success aliases`.

### 2.4 Scaffold blocks — 20 lines, immovable

B1 (`## When to Use`, `**Use when:**`), B2 (`## When This Runs`), B3 (`## Steps`, ```` ``` ````), B5 (```` ``` ````, `### Step 1: Preflight`), one fence in B6, `### Step 2: Read Changes` in B8, B9 (`### Step 3: Associate with REQs`), `### Step 4: Group Unassociated Files` in B11, B12 (`### Step 6: Report`), 4 lines in B13 (fence, `## Error Handling`, `| Situation | Action |`, `|-----------|--------|`), B14 (`## What This Action Does NOT Do`), B15 (`## Rules`), B16 (`## Red Flags`), B17 (`## Verification Checklist`).

### 2.5 F2 — 15 shared content lines the REQ never mentions

These are real duplicated prose and are **required** to reach the ≤20 target. The REQ's "Detailed Requirements" list is incomplete.

| Block | commit.md | inspect.md | Lines | What it is | Differs? |
|---|---|---|---|---|---|
| B4 | 35–37 | 50–52 | 3 | ASCII flow-diagram rows inside `## Steps`: `│`, `├── Group Unassociated ── semantic clustering (1-5 files per group)`, `│` | identical |
| B10 | 102 | 134 | 1 | `What the script settles, so this prose no longer has to:` (the 4 bullets under it differ: commit says "is the bug in the Red Flags below", and the `scope-drift.sh` path differs) | lead-in identical |
| B11 | 109 | 141 | 1 | `**Partial matches count.** If 3 out of 5 files …` | identical |
| B11 | 111 | 143 | 1 | `Files that come back \`-\` remain unassociated and move to Step 4.` | identical |
| B11 | 117, 119–125 | 149, 151–157 | 8 | The whole Step 4 clustering algorithm: `Cluster the remaining files into semantic groups of 1-5 files each:`, numbered items 1/2/3 and item 2's four sub-bullets | identical (item 4 and "When uncertain…" differ and are not counted) |
| B13 | 221 | 368 | 1 | `\| Not a git repo \| Report "Not a git repository" and exit \|` | identical |

**Scope call (D2):** the REQ's Constraints say "Scope is exactly this finding class: do not fix nearby code". These 15 lines *are* the same finding class — `commit-inspect-shared-body` — so moving them is in scope, not scope creep. Moving them is also the only way to reach any threshold below 32. The Step 4 clustering algorithm and the ASCII diagram are the two biggest chunks and are generic grouping guidance, a natural fit for one guide section.

---

## 3. The destination — `skills/do-work/docs/prescribed-shell-primitives.md`

### 3.1 Structure (119 lines, all `##`, zero `###`)

| Line | Heading |
|---|---|
| 1 | `# Prescribed Shell Primitives` |
| 5 | `## Shipped executable homes` |
| 26 | `## Lifecycle timing` |
| **32** | **`## Per-file untracked inventory`** ← the protected-inventory heading |
| 43 | `## Merge-aware commit diff` |
| 49 | `## Commit file listing` |
| 59 | `## Local Git ignore` |
| 69 | `## Verified exact publication` |
| 75 | `## Atomic download publication` |
| 91 | `## Portfolio summary publication` |
| 99 | `## Report image batch publication` |
| 109 | `## Raw text before shell quoting` |
| 113 | `## Diff output filtering` |
| 117 | `## State across command blocks` |

### 3.2 The natural home is already half-written

`prescribed-shell-primitives.md:41` is the last paragraph of `## Per-file untracked inventory` and already names both action files and summarises the fallback contract:

> The complete secret-aware inventory and REQ association ship behind `scripts/protected-inventory.sh`, which orchestrates `tools/checks/uncommitted-inventory.sh` and `tools/checks/associate-files.sh` without duplicating their low-level logic. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` invoke the wrapper; their manual fallback must preserve `-uall`, NUL parsing, rename/copy consumption, secret quarantine, and the scripts' documented exit meanings.

That sentence is the summary the two 207-word fallbacks expand. The clean move is: new `##` section immediately after line 41 (i.e. before `## Merge-aware commit diff` at line 43) carrying the legend, the bullets, the algorithms and both fallbacks, and line 41 rewritten to point at it. Suggested title: `## Protected inventory fallbacks` (anchor `#protected-inventory-fallbacks`).

### 3.3 House style

- Every section is `##` — no `###` anywhere. Keep it that way.
- Prose paragraphs; a fenced `bash` block only when there is a literal invocation (lines 53–55, 63–65, 79–81); one table at the top only.
- Sections cite each other with a Markdown link to a lowercase-hyphen anchor. Exact form, from line 89:
  ```
  Publication here makes the [Verified exact publication](#verified-exact-publication) check, …
  ```
- Guide voice: states the primitive and *why*, never the caller's policy. Line 3 is the contract — "Callers keep one line of local intent plus an invocation. A caller-specific gate always wins; this guide owns the shared primitive, not the action's policy."

### 3.4 Exact citation form from the action files

The suite uses `[Section Title](<relative path>#anchor)` — link text is the heading verbatim. 25 such citations exist. The two depths:

```
from skills/do-work/actions/*.md :
  [Per-file untracked inventory](../docs/prescribed-shell-primitives.md#per-file-untracked-inventory)

from skills/do-work-toolbox/actions/*.md :
  [Per-file untracked inventory](../../do-work/docs/prescribed-shell-primitives.md#per-file-untracked-inventory)
```

Both files already carry exactly this citation (commit.md:73 and :79; inspect.md:84 and :88), so the new pointer follows the same shape with the new anchor. `prime-action-files.md` § Cross-Referencing requires the literal relative path from the citing file's own directory — never `../<package>/…` as shorthand.

---

## 4. The ratchet — `_dev/tests/prescribed-shell-canonicalization.sh` (168 lines)

### F3 — the REQ is wrong about this file: **there are no counts to re-baseline**

The REQ says "The canonicalization ratchet counts headings in the guide and pointers to it from named files: re-baseline those counts". It does not count anything. Every check is a **membership** assertion over a hardcoded list — `grep -Fqx` for each required heading, `grep -Fq` for each pointer site. Adding a new heading or a new pointer **cannot break it** and requires no re-baselining. There is no numeric constant in the file at all.

What the file actually does, in five parts:

1. **Lines 9–25** — 10 named prescribed scripts must exist and be executable.
2. **Lines 40–59** — no `skills/*/tools/*.sh` may call `curl` directly outside a quoted heredoc.
3. **Lines 66–83** — the guide must contain each of 11 required headings, matched with `grep -Fqx` (exact whole line):
   ```
   '## Per-file untracked inventory' '## Merge-aware commit diff' '## Commit file listing'
   '## Local Git ignore' '## Verified exact publication' '## Atomic download publication'
   '## Portfolio summary publication' '## Report image batch publication'
   '## Raw text before shell quoting' '## Diff output filtering' '## State across command blocks'
   ```
   (`## Shipped executable homes` and `## Lifecycle timing` are deliberately *not* in the list.)
4. **Lines 85–117** — 16 shipped files must contain the guide pointer string (see §4.2).
5. **Lines 119–162** — the stale-pattern / old-implementation scan (see F4).

### 4.1 The only edit worth making here

Optional but recommended: pin the new section by adding one line to the required-heading list. Insert after line 67 (`'## Per-file untracked inventory' \`), matching the existing indentation of two spaces:

```sh
  '## Protected inventory fallbacks' \
```

That is the entire ratchet change. Nothing else in this file needs to move.

### 4.2 The 16 shipped files that point at the guide

`core_pointer='../docs/prescribed-shell-primitives.md'` (line 85), 6 sites (lines 87–94):

1. `skills/do-work/actions/commit.md`
2. `skills/do-work/actions/capture.md`
3. `skills/do-work/actions/review-work.md`
4. `skills/do-work/actions/work.md`
5. `skills/do-work/actions/work-reference.md`
6. `skills/do-work/crew-members/background-agents.md`

`sibling_pointer='../../do-work/docs/prescribed-shell-primitives.md'` (line 86), 10 sites (lines 101–111):

7. `skills/do-work-board/actions/board.md`
8. `skills/do-work-knowledge/actions/memory-reference.md`
9. `skills/do-work-knowledge/actions/setup-memory.md`
10. `skills/do-work-knowledge/crew-members/background-agents.md`
11. `skills/do-work-toolbox/actions/ai-report.md`
12. `skills/do-work-toolbox/actions/inspect.md`
13. `skills/do-work-toolbox/actions/install.md`
14. `skills/do-work-toolbox/actions/present-work.md`
15. `skills/do-work-toolbox/actions/stray-check.md`
16. `skills/do-work-toolbox/crew-members/background-agents.md`

Both `commit.md` and `inspect.md` are already on these lists and keep their pointers after the move, so no list edit is needed.

### F4 — trap: the stale-pattern scan (lines 145–162)

It walks every `.md` under `skills/`, skipping only the canonical guide itself and any `CHANGELOG.md`, and **fails** if a file contains any of 9 stale phrases or 7 old implementation fragments. The guide is exempt, so moving prose *into* it is safe. I checked the moved text against all 16 patterns — no match, so no new failures either way. The relevant point for the builder is the opposite direction: this is the mechanism that would prevent the text coming back, and adding a distinctive moved sentence to `stale_patterns_file` (lines 121–131) would be a stronger guard than a difflib count. The REQ mandates the lock-in land in `audit-lockins.sh`, so treat that as an option to raise, not to act on unilaterally.

---

## 5. The lock-in — `_dev/tests/audit-lockins.sh` (156 lines)

### 5.1 Structure

Standard preamble (`#!/usr/bin/env bash`, `set -uo pipefail`, `repo_root=`, `failure_count=0`), then four blocks, each headed by a comment naming the audit finding, its sweep key and its REQ. Each block computes a findings string in a `$( … )`, then a uniform `if [ -n "$var" ]; then while IFS= read -r … printf 'FAIL: …' >&2; failure_count=$((failure_count + 1)); done <<< "$var"; fi`. Ends at lines 151–155 with `if [ "$failure_count" -gt 0 ]; then exit 1; fi` and `printf 'Audit lock-in regressions passed.\n'`.

Existing blocks: Finding 10 exported one-line delegates (REQ-550, lines 8–41); Finding 5 toolbox shims with no callers (REQ-551, lines 43–58); shipped shell delegating check (lines 60–72); Finding 8 dead path pointers in records (REQ-549, lines 74–131); Finding 2 CLI launcher preamble copied (REQ-553, lines 133–149).

### 5.2 Two verbatim existing assertions

```sh
if [ -n "$callerless_shims" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: caller-less toolbox shell shim found: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$callerless_shims"
fi
```

```sh
if [ -n "$hand_rolled_preambles" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$hand_rolled_preambles"
fi
```

### 5.3 Registration in `_dev/tests/contracts/probe-lanes.sh`

Already registered in the fast tier, lines 29–30. Do not change it:

```sh
register_probe audit_lockins_probe "$repo_root/_dev/tests/audit-lockins.sh" \
  'audit lock-in regressions failed (see the attributed FAIL lines above).'
```

### 5.4 Paste-ready assertion, in the file's own style

Insert before the `if [ "$failure_count" -gt 0 ]` block at line 151. `python3` + `difflib` are already used by the fast tier (`shipped-package-reference-contract.sh:6`, `contracts/core-checks.sh:768`), so this adds no new dependency. Note the threshold is the renegotiated 20 from F1, not the REQ's unreachable 10 — the builder must record that decision in the REQ file per CLAUDE.md § Communication.

```sh
# Finding 6: commit-inspect-shared-body (REQ-554)
# The shared inventory legend, file-reading bullets, grouping algorithm and both
# manual fallbacks live once in skills/do-work/docs/prescribed-shell-primitives.md.
# The ceiling is not zero and cannot be: prime-action-files.md requires both actions to
# carry the same section headings, code fences and Error Handling table scaffold, which
# is 17 identical lines on its own. The ceiling sits just above that floor, so real
# prose coming back trips it while a template-shaped heading change does not.
shared_action_prose_ceiling=20
shared_action_prose_lines="$(cd "$repo_root" && python3 -c "
import difflib
a=[l.rstrip() for l in open('skills/do-work/actions/commit.md')]
b=[l.rstrip() for l in open('skills/do-work-toolbox/actions/inspect.md')]
print(sum(1 for i,j,s in difflib.SequenceMatcher(None,a,b).get_matching_blocks()
          if s>=3 for k in range(s) if a[i+k].strip()))
")"
if [ "${shared_action_prose_lines:-999}" -gt "$shared_action_prose_ceiling" ]; then
  printf 'FAIL: commit.md and inspect.md share %s identical lines; ceiling is %s. Move the shared body into skills/do-work/docs/prescribed-shell-primitives.md.\n' \
    "$shared_action_prose_lines" "$shared_action_prose_ceiling" >&2
  failure_count=$((failure_count + 1))
fi

manual_fallback_sites="$(
  rg -c --fixed-strings 'If the script is missing or will not run' \
    "$repo_root/skills" --glob '*/actions/*.md' 2>/dev/null \
    | awk -F: '{ total += $NF } END { print total + 0 }'
)"
if [ "${manual_fallback_sites:-999}" -ne 0 ]; then
  printf 'FAIL: %s manual "do it by hand" fallback sentences remain in shipped actions; the algorithms belong in the prescribed-shell guide, not in the actions.\n' \
    "$manual_fallback_sites" >&2
  failure_count=$((failure_count + 1))
fi
```

The `:-999` defaults are deliberate — `prime-shell-commands.md` § Unchecked Exit Status Reads as Content bans a fallback a later guard cannot tell from a real value; an empty capture here must fail, not pass.

---

## 6. Cross-check `_dev/tests/contract-regressions.sh` (77 lines)

### F5 — the REQ's suspicion is unfounded: **zero predicates pin any of the moved text**

The file is a runner, not an assertion set. It validates the tier, checks its own line count against a ratchet, sources `test-duration-log.sh`, runs four contract files (`core-checks.sh`, `queue-kanban.sh`, `replace-text-section.sh`, `recovery-set-aside.sh`) under a 30-second per-file budget, then sources `probe-batch.sh` and `probe-lanes.sh` and calls `collect_probes`. **No grep, no phrase, no reference to `commit.md` or `inspect.md`.** Nothing to delete.

I also grepped all of `_dev/` (excluding `.md`) for twelve distinctive phrases from the moved text — `If the script is missing`, `Secret-shaped`, `Modified files`, `New/untracked`, `Deleted secret-shaped`, `do it by hand`, `Partial matches count`, `renamed path is tagged M`, `staged-new and untracked`, `What the script settles`, `Group Unassociated`, `uall`. **Zero hits in any test.**

### F6 — but there is a separate trap in this file

```
_dev/tests/contract-regressions.sh:9   fast_contract_line_ceiling=77
_dev/tests/contract-regressions.sh:16-21  fails when the file exceeds it
```

The file is **exactly 77 lines** — at its ceiling, with zero headroom. It cannot grow by a single line. The REQ already routes the lock-in to `audit-lockins.sh`, which is why this is not a problem, but a builder who improvises and adds anything here breaks the gate immediately.

### F7 — one more test the REQ's write_set does not name

`_dev/tests/defensive-surface-audit.sh` (also a fast probe) holds ten `assert_phrase_absent` calls, four of them against `inspect.md` and two against `commit.md`. They are *absence* assertions for generic guidance deleted by REQ-168 (`'No REQ matches — just commit everything together'`, `'Single commit with >20 files'`, `'This looks ready to commit'`, and three more). None of those phrases appear in the shared blocks, so the move does not trip them — but the builder must not reintroduce any of those sentences while rewording. No change needed to this file.

Also present, both harmless: `_dev/tests/staged-skills-contract.sh:38` lists `commit.md` in a mirror set (heavy tier only), and `_dev/tests/shipped-package-reference-contract.sh:1712,1725` uses the string `do-work-toolbox/actions/inspect.md` inside a self-contained fixture, not as a pin on the real file.

---

## 7. Primes — what governs this change

### `_dev/primes/prime-action-files.md` (120 lines)

**Cross-reference form (line 91)** — this is the rule the moved fallback's two path fixups turn on:

> Cross-reference same-package actions by their local path (for example `actions/work.md`). Cross a package boundary with the **literal relative path from the citing file's own directory** — the spelling a reader could paste into a terminal there and have resolve, in both the source tree and an install. The depth is per-file, not a fixed prefix: from a package-root `SKILL.md` the sibling is `../do-work-board/...`, from `actions/` or `crew-members/` it is `../../do-work-board/...`, and from `tools/queue-kanban/` it is `../../../do-work/...`. Never write `../<package>/...` as shorthand for "up to the skills folder" — that skill-root-relative reading was retired by REQ-249, and from anywhere below a package root it points at nothing.

**Line 93** — the citation class is drawn by meaning, not punctuation, and `_dev/tests/shipped-package-reference-contract.sh` enforces it in both topologies. The moved text carries `` `../../do-work-toolbox/actions/stray-check.md` `` — correct from the guide's directory, unchanged from commit.md's spelling.

**Line 95** — shipped files must never cite this repo's `CLAUDE.md` or `AGENTS.md`; the guide is shipped, so this binds the new section.

**Template, lines 70–80** — this is the source of F1's immovable 20 lines:

> **Required:** Description blockquote, Steps (numbered). **Common:** Input, Output Format, When to Use.
> **Section order when present:** Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.

**Line 76** — the standard for what stays behind in each action after the move:

> **State intent, not a directive rule, when a capable model can infer the rest.** "Report drift, don't fix it inline" gives the model this action's boundary in one line — a five-line Rules section re-deriving why inline fixes are bad adds nothing a capable model didn't already know.

**Trap, line 107** — the direct warning for this REQ:

> [family: alternate-writer-contract-drift] Changing an emitted artifact or routing rule only at its primary writer leaves alternate modes and action-bearing readers silently following the old contract; grep every writer and downstream reader against the recorded scope before declaring it shipped — including sibling-package actions, board labels, and prime sentences that restate the rule (REQ-566 found six such restatements outside the three files that defined the contract).

Applied here: after the move, grep the moved sentences across all of `skills/` — `stray-check.md` also carries a `#per-file-untracked-inventory` citation (line 43) and may restate part of the same rule.

### `_dev/primes/prime-shell-commands.md` (62 lines)

**Line 13** — the guide's `## Per-file untracked inventory` section is the canonical home of exactly the rule the 207-word fallback restates:

> **`git status --porcelain` collapses wholly-untracked directories** into a single `?? dir/` row — it does not list the files inside. Any step that enumerates untracked files per-item (read each, check extension/size/name) must use `git status --porcelain --untracked-files=all` (`-uall`) or `git ls-files --others --exclude-standard`.

**Line 27** — the "one home" reasoning behind this REQ:

> When a review finds a bug in prescribed-command logic, **grep the same primitive across all actions before calling it fixed** — these patterns are usually copy-pasted, so the fix is rarely local. (The first trap above had been copy-pasted into four action files; the audit only flagged one of them.)

**§ Closed Enumerations Go Stale, line 51** — governs the ratchet's heading list and the new section's wording:

> When a rule applies "whenever X happens" (load a guardrail, honor an enum, keep a guide in sync), state the trigger _condition_ in the rule's canonical home and mark any caller/value list as illustrative, not exhaustive.

**§ Unchecked Exit Status Reads as Content, lines 33–47** — governs the lock-in's own shell (hence the `:-999` defaults in §5.4).

Both primes end with a `## Lessons` pointer. The REQ's frontmatter records both lesson satellites (`lessons-action-files.md`, `lessons-shell-commands.md`) as dropped for budget, so the builder works from the primes alone.

---

## Summary of findings

- **F1 (blocker)** — the ≤10 lock-in limit is unreachable. Floor is 17 (template scaffolding); the REQ's own named scope only reaches 32. Recommend pinning ≤20 and recording the renegotiation in the REQ file.
- **F2** — the REQ's Detailed Requirements list is incomplete: 15 further shared content lines (ASCII diagram, Step 4 clustering algorithm, "Partial matches count", "Files that come back `-`", "What the script settles" lead-in, the Error Handling row). Same finding class, and required to reach any threshold under 32.
- **F3** — the REQ is wrong that the ratchet "counts headings … re-baseline those counts". It is a membership list with no numeric constants; adding a heading cannot break it. One optional one-line addition after line 67.
- **F4** — the stale-pattern scan exempts the guide, so the move is safe; no moved text matches any of the 16 patterns.
- **F5** — `contract-regressions.sh` pins nothing from the moved text; it is a runner. Nothing to delete there or anywhere in `_dev/`.
- **F6** — `contract-regressions.sh` is exactly at its 77-line ceiling. It cannot grow.
- **F7** — `defensive-surface-audit.sh` holds six absence assertions on these two files. Not tripped by the move, but the builder must not reintroduce those sentences.
- **D1** — keep the REQ's exact difflib metric, pin ≤20.
- **D2** — moving the 15 unnamed lines is in scope, not scope creep.

## Key line numbers

- `skills/do-work/actions/commit.md` — 35-37, 59-67, 79, 83-88, 102, 109-125, 221
- `skills/do-work-toolbox/actions/inspect.md` — 50-52, 67-75, 88, 92-97, 134, 141-157, 368
- `skills/do-work/docs/prescribed-shell-primitives.md` — `## Per-file untracked inventory` at 32, pointer paragraph at 41, insert new section before 43
- `_dev/tests/prescribed-shell-canonicalization.sh` — heading list 66-77 (insert after 67), pointer sites 87-94 and 101-111, stale scan 119-162
- `_dev/tests/audit-lockins.sh` — insert new block before 151
- `_dev/tests/contracts/probe-lanes.sh` — 29-30 (already registered, do not touch)
- `_dev/tests/contract-regressions.sh` — 9 (`fast_contract_line_ceiling=77`, file is at 77)
