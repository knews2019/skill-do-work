# REQ-554 exploration — move the commit.md / inspect.md shared body into the prescribed-shell guide

Repo at HEAD `15e2ec3` ("[work run] REQ-592 pipeline: triage, estimate, exploration, scope"). Read-only; nothing edited, staged or committed.

REQ-554 is the request to stop `skills/do-work/actions/commit.md` and `skills/do-work-toolbox/actions/inspect.md` from restating the same 46 lines of inventory and association prose, by giving that prose one home in `skills/do-work/docs/prescribed-shell-primitives.md`.

---

## 0. Baseline freshness — the numbers did not move

A prior exploration exists at `do-work/runs/work-2026-09-05-170806/REQ-554-exploration.md`, written at HEAD `71eb49f3`. Every measurement in it was re-taken at `15e2ec3` and every one is unchanged:

| Measurement | Prior (71eb49f3) | HEAD (15e2ec3) | Moved? |
|---|---|---|---|
| difflib shared non-blank lines | 46 | **46** | no |
| fallback sentence sites | 4 | **4** | no |
| `skills/do-work/actions/commit.md` | 260 lines | **260** | no |
| `skills/do-work-toolbox/actions/inspect.md` | 448 lines | **448** | no |
| `skills/do-work/docs/prescribed-shell-primitives.md` | 119 lines | **119** | no |
| `_dev/tests/prescribed-shell-canonicalization.sh` | 168 lines | **168** | no |
| `_dev/tests/audit-lockins.sh` | 156 lines | **156** | no |

The REQ was captured at `dc8a64e3`; its own Reproduce command still prints 46. **The baseline is not stale.**

Two things in the prior exploration *are* wrong, and both would produce a broken lock-in if copied. They are corrected in §5 and §6 and marked **C1** and **C2**.

---

## 1. The claim, re-verified

```
$ python3 -c "import difflib;a=[l.rstrip() for l in open('skills/do-work/actions/commit.md')];b=[l.rstrip() for l in open('skills/do-work-toolbox/actions/inspect.md')];print(sum(1 for i,j,s in difflib.SequenceMatcher(None,a,b).get_matching_blocks() if s>=3 for k in range(s) if a[i+k].strip()))"
46

$ rg -n 'If the script is missing or will not run' skills --glob '*.md'
skills/do-work-toolbox/actions/inspect.md:88
skills/do-work-toolbox/actions/inspect.md:145
skills/do-work/actions/commit.md:79
skills/do-work/actions/commit.md:113
```

Both halves of the RED state reproduce exactly as the REQ states.

### 1.1 The 46 lines, split by kind

17 matching runs of 3+ lines. Splitting them into template scaffolding and real prose:

| Kind | Lines | Can it move to the guide? |
|---|---|---|
| Section headings, code fences, table header rows | **20** | No — `_dev/primes/prime-action-files.md:70-80` requires both actions to carry them |
| Duplicated prose | **26** | Mostly yes; 4 of the 26 are structural (see §2.4) |

### 1.2 The full block map

| Block | commit.md | inspect.md | Non-blank | What it is |
|---|---|---|---|---|
| B1 | 8-11 | 6-9 | 2 | `## When to Use`, `**Use when:**` — scaffold |
| B2 | 19-21 | 16-18 | 1 | `## When This Runs` — scaffold |
| B3 | 23-26 | 36-39 | 2 | `## Steps`, fence — scaffold |
| B4 | 35-37 | 50-52 | 3 | ASCII flow-diagram rows — **structural, see §2.4** |
| B5 | 41-44 | 57-60 | 2 | fence, `### Step 1: Preflight` — scaffold |
| B6 | 57-63 | 65-71 | 5 | fence + tag legend lead-in + M/A/D bullets — **REQ class A** |
| B7 | 66-68 | 74-76 | 1 | `Secret-shaped matching is case-insensitive…` — **REQ class A** |
| B8 | 80-89 | 89-98 | 6 | `### Step 2: Read Changes` + lead-in + four file-reading bullets — **REQ class B** |
| B9 | 91-93 | 113-115 | 1 | `### Step 3: Associate with REQs` — scaffold |
| B10 | 101-103 | 133-135 | 1 | `What the script settles, so this prose no longer has to:` — association semantics |
| B11 | 108-125 | 140-157 | 12 | `### Step 4` heading + association fallback + `Partial matches count` + `Files that come back` + the whole Step 4 clustering algorithm |
| B12 | 183-185 | 239-241 | 1 | `### Step 6: Report` — scaffold |
| B13 | 215-221 | 362-368 | 5 | fence, `## Error Handling`, table header ×2, `\| Not a git repo \| …` — **structural** |
| B14 | 229-231 | 377-379 | 1 | `## What This Action Does NOT Do` — scaffold |
| B15 | 238-240 | 429-431 | 1 | `## Rules` — scaffold |
| B16 | 246-248 | 436-438 | 1 | `## Red Flags` — scaffold |
| B17 | 253-255 | 441-443 | 1 | `## Verification Checklist` — scaffold |

---

## 2. F1 (blocker) — the `<= 10` lock-in limit is unreachable

`Constraints` in the REQ file (`do-work/queue/REQ-554-move-the-commit-and-inspect-shared-body-into-the-prescribed-shell-guide.md:49`) says:

> Lock-in limit: commit.md/inspect.md shared lines <= 10 after this REQ (today 46)

That number is below the floor. Measured by deleting lines from both files and re-running the REQ's own metric:

| Scenario | difflib count |
|---|---|
| HEAD today | **46** |
| Move exactly the three classes the REQ names | **32** |
| Move every shared prose line, structural rows included | **17** |

**17 is the absolute floor**, and every survivor is template scaffolding that `_dev/primes/prime-action-files.md:70-80` requires both actions to carry:

```
## When to Use              ### Step 1: Preflight            ## Error Handling
**Use when:**               ### Step 4: Group Unassociated   | Situation | Action |
## When This Runs           ### Step 6: Report               |-----------|--------|
## Steps                    ``` ``` ```                      ## What This Action Does NOT Do
                                                             ## Rules
                                                             ## Red Flags
                                                             ## Verification Checklist
```

Reaching 10 would require deleting required sections from one of the two actions. That is not a refactor, it is a template violation.

**The `fallback sentences in actions: 0` half of the limit is achievable exactly as written** and needs no renegotiation.

### 2.1 C1 — the prior exploration's recommended ceiling of 20 is also wrong

`work-2026-09-05-170806/REQ-554-exploration.md` recommends pinning `<= 20`. That is below what is reachable without damaging the two actions.

Two of the shared blocks are not prose and cannot move:

**B4 — the ASCII flow diagram** (`skills/do-work/actions/commit.md:35-37`, `skills/do-work-toolbox/actions/inspect.md:50-52`). Two of the three shared lines are bare `│` characters; the third is `├── Group Unassociated ── semantic clustering (1-5 files per group)`. They sit inside each action's own flow diagram, which otherwise differs completely — commit's has a `Commit` branch, inspect's has `Assess Readiness` and a `REQ/UR scope?` sub-line. Deleting them breaks the diagram.

**B13's content row** — `| Not a git repo | Report "Not a git repository" and exit |` at `commit.md:221` / `inspect.md:368`. It is one row of each action's own Error Handling table, a template-required section. A guide does not own an action's error table.

With those two kept intact, the lowest reachable count is **21**, not 20.

### 2.2 The realistic landing values

| Scenario | Count | In the guide's charter? |
|---|---|---|
| **S1** — the REQ's three named classes only | **32** | Yes |
| **S2** — S1 + association semantics (`What the script settles` lead-in, `Partial matches count`, `Files that come back -`) | **29** | Yes |
| **S3** — S2 + the Step 4 clustering algorithm | **21** | **No** — see §2.3 |
| Floor incl. structural rows | 17 | n/a |

### 2.3 Why S3 is the wrong answer even though it scores best

`skills/do-work/docs/prescribed-shell-primitives.md:3` states the guide's charter:

> This is the canonical shipped rationale and executable-home contract for deterministic mechanics used across do-work actions. … Callers keep one line of local intent plus an invocation. A caller-specific gate always wins; this guide owns the shared primitive, not the action's policy.

The Step 4 clustering algorithm (`commit.md:117,119-125` / `inspect.md:149,151-157`) is semantic file grouping — "A component and its test file", "Config file changes that go together". It is not a shell primitive and not an executable-home contract. Putting it in this guide to win 8 lines on a metric buys the number by weakening the guide.

**D1 — recommendation: land at S2 (29) and pin the ceiling at 29.** S2 moves everything that is genuinely about `scripts/protected-inventory.sh` — the tag legend, the reading bullets, the association semantics, both manual fallbacks — and leaves alone what belongs to each action. The remaining clustering duplication is a real finding, but it needs a different home and is a separate REQ.

If the maintainer prefers to hold the line at the REQ's literal Detailed Requirements, land at S1 and pin 32. Either way, **the `<= 10` must be renegotiated in writing in the REQ file and the commit message** before the assertion is written, per `CLAUDE.md` § Communication.

### 2.4 The metric under-counts the real win

The single biggest duplication in the pair is the 207-word inventory algorithm at `commit.md:79` / `inspect.md:88`. It is duplicated in full. It contributes **0** to the count of 46, because the two copies differ by two path fixups, which breaks the matching run.

Anyone judging this REQ by the difflib delta alone (46 → 29) will undervalue it. The `fallback sentences: 0` assertion is what actually pins the largest win.

---

## 3. The two fallbacks, verified byte by byte

### Fallback 1 — the 207-word inventory algorithm

`skills/do-work/actions/commit.md:79` vs `skills/do-work-toolbox/actions/inspect.md:88`. A difflib opcode diff returns exactly three edits, all one change of relative-path depth:

```
replace: commit='docs'                     inspect='..'
insert:  commit=''                         inspect='do-work/docs/'
delete:  commit='../../do-work-toolbox/'   inspect=''
```

That is, the two deltas the REQ names at line 38 and nothing else:

```
commit.md:79   [Per-file untracked inventory](../docs/prescribed-shell-primitives.md#per-file-untracked-inventory)
inspect.md:88  [Per-file untracked inventory](../../do-work/docs/prescribed-shell-primitives.md#per-file-untracked-inventory)

commit.md:79   … and `../../do-work-toolbox/actions/stray-check.md`'s Red Flags record that it has been hit.
inspect.md:88  … and `actions/stray-check.md`'s Red Flags record that it has been hit.
```

**Both fixups resolve to commit.md's spelling when the text moves to the guide.** `skills/do-work/docs/` sits at the same depth below `skills/` as `skills/do-work/actions/`, so:

- `../../do-work-toolbox/actions/stray-check.md` — keep verbatim, correct from the guide. The guide already uses this exact depth at line 41 for `../../do-work-toolbox/actions/inspect.md`.
- The `[Per-file untracked inventory](../docs/…)` link becomes the in-file anchor `#per-file-untracked-inventory`, or drops entirely since the moved text sits directly under that section.

### Fallback 2 — the association fallback

`commit.md:113` vs `inspect.md:145`. `diff` returns no output — byte-identical:

```
If the script is missing or will not run, do it by hand: glob both directories, read each REQ's `status` (accepting every alias above) and `## Implementation Summary` list, path-match, and tie-break on the latest `completed_at`.
```

**Trap: `accepting every alias above` dangles after the move.** It back-references the `What the script settles` bullets at `commit.md:104-107` / `inspect.md:136-139`, which stay in the actions because they differ — `commit.md:104` adds "is the bug in the Red Flags below", and `commit.md:107` writes `tools/checks/scope-drift.sh` where `inspect.md:139` writes `../../do-work/tools/checks/scope-drift.sh`. The moved sentence must name the Schema Read Contract's terminal-success aliases directly instead of pointing "above".

### What stays local — the read-only delta

The `X` and `XD` bullets are reworded per mode and are correctly not in any matching block:

```
commit.md:64  - **X** — … an ambiguous addition beside `X`/`XD`; fully excluded
inspect.md:72 - **X** — … an ambiguous addition beside `X`/`XD`; fully excluded from analysis

commit.md:65  - **XD** — deleted secret-shaped path; eligible for a deletion-only commit
inspect.md:73 - **XD** — deleted secret-shaped path; path/deletion state only
```

This is exactly the delta the REQ says to keep local (`REQ-554…md:41`).

---

## 4. The destination guide

`skills/do-work/docs/prescribed-shell-primitives.md`, 119 lines, all `##`, zero `###`.

| Line | Heading |
|---|---|
| 1 | `# Prescribed Shell Primitives` |
| 5 | `## Shipped executable homes` |
| 26 | `## Lifecycle timing` |
| **32** | **`## Per-file untracked inventory`** |
| **41** | *(the pointer paragraph — not a heading)* |
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

### 4.1 The home is already half written

`prescribed-shell-primitives.md:41` is the closing paragraph of `## Per-file untracked inventory` and already names both action files and summarises the contract the two 207-word fallbacks expand:

> The complete secret-aware inventory and REQ association ship behind `scripts/protected-inventory.sh`, which orchestrates `tools/checks/uncommitted-inventory.sh` and `tools/checks/associate-files.sh` without duplicating their low-level logic. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` invoke the wrapper; their manual fallback must preserve `-uall`, NUL parsing, rename/copy consumption, secret quarantine, and the scripts' documented exit meanings.

**The clean move:** one new `##` section inserted after line 41, before `## Merge-aware commit diff` at line 43, carrying the legend, the bullets and both fallbacks. Line 41 is rewritten to point at it instead of restating the contract.

Suggested title `## Protected inventory fallbacks` (anchor `#protected-inventory-fallbacks`). The REQ's Builder Guidance grants latitude on the title.

**Do not delete the `orchestrates` sentence.** REQ-555 (rewriting the guide's executable-homes table to the do-work-cli route form) removes it; its frontmatter already declares `depends_on: [REQ-554]`, so REQ-554 lands first and leaves that sentence for REQ-555.

### 4.2 House style to match

- Every heading is `##`. No `###` exists in the file — keep it that way.
- Prose paragraphs; a fenced ```bash``` block only for a literal invocation (lines 53-55, 63-65, 79-81). One table, at the top only.
- Cross-section links use the heading verbatim as link text against a lowercase-hyphen anchor. Exact existing form, `prescribed-shell-primitives.md:89`:
  ```
  Publication here makes the [Verified exact publication](#verified-exact-publication) check, …
  ```

### 4.3 Citation depth from each action

`_dev/primes/prime-action-files.md:91` requires the literal relative path from the citing file's own directory, and explicitly retires `../<package>/…` as shorthand:

```
from skills/do-work/actions/*.md :
  [<Section Title>](../docs/prescribed-shell-primitives.md#<anchor>)

from skills/do-work-toolbox/actions/*.md :
  [<Section Title>](../../do-work/docs/prescribed-shell-primitives.md#<anchor>)
```

Both files already carry exactly these strings (`commit.md:73` and `:79`; `inspect.md:84` and `:88`), so the new pointer reuses the shape with the new anchor. Each file must retain at least one occurrence of its pointer string or `prescribed-shell-canonicalization.sh:88` / `:107` fails.

---

## 5. Every ratchet constant — the answer is that there are none to re-baseline

### 5.1 F2 — the REQ's ratchet claim is false

`REQ-554…md:42` says:

> The canonicalization ratchet counts headings in the guide and pointers to it from named files: re-baseline those counts in `_dev/tests/prescribed-shell-canonicalization.sh` in the same commit.

`_dev/tests/prescribed-shell-canonicalization.sh` contains **zero numeric constants**. Every check is a membership assertion over a hardcoded list:

| Lines | Check | Form |
|---|---|---|
| 9-25 | 10 prescribed scripts exist and are executable | `[ ! -x … ]` per name |
| 40-59 | no `skills/*/tools/*.sh` calls `curl` outside a quoted heredoc | awk scan |
| **66-83** | the guide contains each of 11 required headings | **`grep -Fqx` per heading** |
| **85-117** | 16 shipped files contain the guide pointer string | **`grep -Fq` per site** |
| 119-162 | 9 stale phrases and 7 old-implementation fragments absent from shipped `.md` | `grep -Fq` per pattern |

Adding a heading to the guide, or a pointer to a file, **cannot fail this test**. There is nothing to re-baseline. The test passes at HEAD today (exit 0, "Prescribed shell primitive canonicalization checks passed.").

### 5.2 The complete ratchet-constant inventory, repo-wide

`rg -n '^[a-z_]*(ceiling|limit|max|baseline|budget|threshold|cap)[a-z_]*=[0-9]+'` across `_dev/` and `skills/` returns exactly two live constants:

| Constant | File:line | Value | Needs re-baselining for REQ-554? |
|---|---|---|---|
| `fast_contract_line_ceiling` | `_dev/tests/contract-regressions.sh:9` | **77** | **No** — not in the write_set. See §5.3. |
| `test_file_budget_seconds` | `_dev/tests/maintainer-verify.sh:12` | 30 | No — per-file timing budget, unrelated |

**No existing ratchet constant anywhere in the repo requires a new value for this REQ.** The only numbers this change sets are the *new* constants it introduces in `_dev/tests/audit-lockins.sh` (§6).

### 5.3 F3 (trap) — `contract-regressions.sh` has zero headroom

`_dev/tests/contract-regressions.sh:9` sets `fast_contract_line_ceiling=77`, lines 16-21 fail when the file exceeds it, and `wc -l` reports the file is **exactly 77 lines**. It cannot grow by one line.

This file is not in the write_set and needs no change. It matters only because a builder who improvises an assertion there instead of in `audit-lockins.sh` breaks the fast gate on the first line added.

### 5.4 The one worthwhile edit to the canonicalization ratchet

Optional but recommended: add the new section to the required-heading membership list, so a later edit cannot quietly delete the new home. Insert after `prescribed-shell-canonicalization.sh:67`, matching the existing two-space indentation:

```sh
  '## Protected inventory fallbacks' \
```

That is the entire change to this file. Note that `## Shipped executable homes` and `## Lifecycle timing` are deliberately absent from the list, so membership is a choice, not an obligation.

---

## 6. The lock-in

### 6.1 Where it lands

`_dev/tests/audit-lockins.sh` — 156 lines, mode `-rwxr-xr-x`, already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh:29`:

```sh
register_probe audit_lockins_probe "$repo_root/_dev/tests/audit-lockins.sh" \
  'audit lock-in regressions failed (see the attributed FAIL lines above).'
```

Do not change that registration. The file passes at HEAD today (exit 0).

Existing blocks, each headed by a comment naming its audit finding, sweep key and REQ:

| Lines | Finding |
|---|---|
| 8-41 | Finding 10, exported one-line delegates (REQ-550) |
| 43-58 | Finding 5, toolbox shims with no callers (REQ-551) |
| 60-72 | shipped shell delegating check (REQ-551 companion) |
| 74-131 | Finding 8, dead path pointers in records (REQ-549) |
| 133-149 | Finding 2, CLI launcher preamble copied (REQ-553) |
| 151-155 | `if [ "$failure_count" -gt 0 ]; then exit 1; fi` + pass line |

Every block uses the same shape: compute findings into a `$( … )`, then `if [ -n "$var" ]; then while IFS= read -r … printf 'FAIL: …' >&2; failure_count=$((failure_count + 1)); done <<< "$var"; fi`.

### 6.2 C2 — the prior exploration's assertion does not work

The prior exploration's paste-ready fallback assertion uses:

```sh
rg -c --fixed-strings 'If the script is missing or will not run' \
  "$repo_root/skills" --glob '*/actions/*.md' 2>/dev/null | awk -F: '{ total += $NF } END { print total + 0 }'
```

Run at HEAD, that prints **`0`** — with four real violations present. Ripgrep's `*` does not cross a directory separator, so `*/actions/*.md` matches nothing under `skills/`. Verified:

```
--glob '*/actions/*.md'   → (no output)
--glob '**/actions/*.md'  → skills/do-work-toolbox/actions/inspect.md:2
                            skills/do-work/actions/commit.md:2
--glob '*.md'             → same two lines
```

An assertion built on that glob is green on day one for the wrong reason and never fires again. It must use `**/actions/*.md`.

### 6.3 Paste-ready assertion

Insert immediately before `_dev/tests/audit-lockins.sh:151`. Substitute `29` with the value measured after the move (29 for scenario S2, 32 for S1).

```sh
# Finding 6: commit-inspect-shared-body (REQ-554)
# The inventory tag legend, the four file-reading bullets, the association semantics and
# both manual "do it by hand" fallbacks live once, in
# skills/do-work/docs/prescribed-shell-primitives.md.
# The ceiling is not zero and cannot be: prime-action-files.md requires both actions to
# carry the same section headings, code fences and Error Handling table scaffold, which is
# 17 identical lines once every shared sentence is gone. The ceiling sits at the value this
# REQ landed on, so returning prose trips it while the template scaffold never does.
shared_action_prose_ceiling=29
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

# The two 207-word inventory fallbacks differ by two relative-path fixups, so the difflib
# count above scores them at 0. This is the assertion that actually pins them at zero.
manual_fallback_matches="$(rg -n --fixed-strings 'If the script is missing or will not run' \
  "$repo_root/skills" --glob '**/actions/*.md' 2>/dev/null)"
manual_fallback_scan_status=$?
if [ "$manual_fallback_scan_status" -gt 1 ]; then
  printf 'FAIL: could not scan shipped actions for manual fallback sentences (rg exit %s).\n' \
    "$manual_fallback_scan_status" >&2
  failure_count=$((failure_count + 1))
elif [ -n "$manual_fallback_matches" ]; then
  while IFS= read -r fallback_site; do
    [ -z "$fallback_site" ] && continue
    printf 'FAIL: manual "do it by hand" fallback remains in a shipped action: %s\n' \
      "${fallback_site#"$repo_root/"}" >&2
    failure_count=$((failure_count + 1))
  done <<< "$manual_fallback_matches"
fi
```

**Verified today, before any edit:** the difflib block measures 46 against 29 and prints its FAIL line. The fallback block prints four FAIL lines, naming `commit.md` twice and `inspect.md` twice. Both are RED, as the Red-Green Proof requires.

Design notes:

- The second block reads `rg`'s own exit status rather than piping into `awk`, because `awk` prints `0` on empty input whether that means "no matches" or "the scan never ran". `_dev/primes/prime-shell-commands.md` § Unchecked Exit Status Reads as Content bans a fallback a later guard cannot tell from a real value. `rg` exit 1 is no matches; 2 or higher is a scan failure and is reported as one.
- The `:-999` default on the difflib count guards the same way: an empty capture must fail, not pass.
- `python3` is already a fast-tier dependency (`_dev/tests/shipped-package-reference-contract.sh:6`, `_dev/tests/contracts/core-checks.sh:769`) and `rg` is already used three times inside `audit-lockins.sh` itself (lines 28, 48, 138). No new dependency.
- Names satisfy `skills/do-work/crew-members/coding-guardrails.md` § 5 Naming for Reach: `shared_action_prose_ceiling`, `shared_action_prose_lines`, `manual_fallback_matches`, `manual_fallback_scan_status`, `fallback_site` — all two words or more, all findable by plain-text search.

---

## 7. Other tests that touch these files

| File | Relationship | Action needed |
|---|---|---|
| `_dev/tests/defensive-surface-audit.sh:34,37` (commit.md) and `:40,43,46,49` (inspect.md) | Six `assert_phrase_absent` calls pinning generic guidance deleted by REQ-168 | **None** — no pinned phrase appears in any shared block. But rewording must not reintroduce `'No REQ matches — just commit everything together'`, `'Single commit with >20 files'`, `'This looks ready to commit'`, `'Files listed with no REQ association attempt'`, `'Debug artifacts (console.log, debugger, commented-out blocks)'`, or the `git status` / `git diff` reading phrase. |
| `_dev/tests/prescribed-shell-canonicalization.sh:145-162` | Stale-pattern scan over every shipped `.md` | **None** — line 146 skips the canonical guide, so text moved into it is exempt. No moved line matches any of the 16 patterns. |
| `_dev/tests/contract-regressions.sh` | Runner, not an assertion set | **None** — no reference to either action. See §5.3 for the line-ceiling trap. |
| `_dev/tests/staged-skills-contract.sh:38` | Lists `commit.md` in a heavy-tier mirror set | **None** |
| `_dev/tests/shipped-package-reference-contract.sh:1759,1772` | Uses the string `do-work-toolbox/actions/inspect.md` inside a self-contained citation fixture | **None** — a fixture, not a pin on the real file |

---

## 8. Primes that bind this change

### `_dev/primes/prime-action-files.md`

- **Line 91, Cross-Referencing** — the literal relative path from the citing file's own directory; `../<package>/…` shorthand was retired by REQ-249. This is the rule the fallback's two path fixups turn on.
- **Line 93** — the citation class is drawn by meaning, not punctuation, and `shipped-package-reference-contract.sh` enforces it in both topologies. The moved `../../do-work-toolbox/actions/stray-check.md` is correct from the guide's directory.
- **Line 95** — shipped files must never cite this repo's `CLAUDE.md` or `AGENTS.md`. The guide is shipped, so this binds the new section.
- **Lines 70-80, the Template** — the source of the 17-line floor. Required sections plus the fixed section order.
- **Line 76** — the standard for what stays behind: "State intent, not a directive rule, when a capable model can infer the rest." Each action keeps one line of local intent plus the citation, not a summary of what moved.
- **Line 107, trap `alternate-writer-contract-drift`** — after the move, grep the moved sentences across all of `skills/`. `skills/do-work-toolbox/actions/stray-check.md:43` carries its own `#per-file-untracked-inventory` citation and restates part of the same untracked-inventory rule.

### `_dev/primes/prime-shell-commands.md`

- **Line 13** — the guide's `## Per-file untracked inventory` section already owns the `-uall` rule the 207-word fallback restates. This is why the guide is the right home.
- **Line 27** — "grep the same primitive across all actions before calling it fixed"; these patterns are copy-pasted, so the fix is rarely local.
- **§ Closed Enumerations Go Stale, line 51** — governs the canonicalization heading list and the new section's wording: state the trigger condition, mark example lists illustrative.
- **§ Unchecked Exit Status Reads as Content, lines 33-47** — governs the lock-in's own shell, which is why the fallback assertion checks `rg`'s exit status instead of trusting an `awk` total.

Both lesson satellites (`lessons-action-files.md` 3968 tokens, `lessons-shell-commands.md` 3385 tokens) were dropped for budget per the REQ frontmatter, so the builder works from the primes alone.

---

## 9. Route: B — explore-then-build

The edit is mechanical prose movement across five named files with no design uncertainty left. But three of the REQ's own statements are wrong, and a builder following the REQ literally ships a broken test:

- The `<= 10` limit is below the mathematical floor of 17 (F1).
- The claim that the canonicalization ratchet "counts headings … re-baseline those counts" is false; that file has no numeric constants (F2).
- The Detailed Requirements list is incomplete — it names 11 of the 26 shared prose lines.

Route A would produce a red-forever assertion. Route C is unwarranted — this document already names every file, every line and the exact assertion text.

**The renegotiated ceiling must be written into the REQ file and the commit message** before the assertion is written, per `CLAUDE.md` § Communication ("write the challenge and the chosen resolution into the REQ file or commit message and continue").

---

## 10. Findings and decisions

| Code | Statement |
|---|---|
| **F1** | The `<= 10` lock-in limit is unreachable. Floor is 17 template-scaffold lines; the REQ's own named scope reaches only 32. Must be renegotiated. |
| **F2** | The REQ is wrong that `prescribed-shell-canonicalization.sh` counts headings and pointers. It has zero numeric constants; every check is a membership assertion. Nothing to re-baseline. |
| **F3** | `contract-regressions.sh` is exactly 77 lines against its own ceiling of 77. Zero headroom. Not in the write_set — a trap, not a change. |
| **F4** | The difflib metric scores the 207-word duplicated fallback at 0, because two path fixups break the run. The metric under-counts the real win; the fallback-count assertion captures it. |
| **F5** | `defensive-surface-audit.sh` pins six phrases on these two files. Not tripped by the move, but rewording must avoid all six. |
| **F6** | `accepting every alias above` in the association fallback dangles once moved; the bullets it references stay in the actions because they differ. |
| **C1** | Correction to the prior exploration: its recommended ceiling of 20 is below the reachable value of 21, and would force damage to each action's ASCII flow diagram and Error Handling table. |
| **C2** | Correction to the prior exploration: its paste-ready `rg --glob '*/actions/*.md'` matches zero files and prints 0 today with four violations present. Use `**/actions/*.md`. |
| **D1** | Land at scenario S2 (29) — move everything about `protected-inventory.sh`, leave the Step 4 clustering algorithm alone because it is not a shell primitive and the guide's charter at line 3 says so. Pin the ceiling at the measured value. |
| **D2** | Keep the REQ's exact difflib metric. Do not invent a filtered metric that strips headings — that is a second definition of "shared", and the REQ's Reproduce command is the one the audit report quotes. |

---

## Key line numbers

- `/home/user/skill-do-work/skills/do-work/actions/commit.md` — 35-37 (diagram, keep), 59/61-63 (legend), 64-65 (X/XD, keep local), 67 (secret-shaped), 79 (fallback 1), 83/85-88 (reading bullets), 102 (settles lead-in), 104-107 (settles bullets, differ, keep), 109/111 (association semantics), 113 (fallback 2), 117/119-125 (clustering, out of charter), 221 (error row, keep)
- `/home/user/skill-do-work/skills/do-work-toolbox/actions/inspect.md` — 50-52, 67/69-71, 72-73, 75, 88, 92/94-97, 134, 136-139, 141/143, 145, 149/151-157, 368
- `/home/user/skill-do-work/skills/do-work/docs/prescribed-shell-primitives.md` — charter at 3, `## Per-file untracked inventory` at 32, pointer paragraph at 41, insert the new section before 43, citation form example at 89
- `/home/user/skill-do-work/_dev/tests/prescribed-shell-canonicalization.sh` — heading list 66-77 (insert after 67), core pointer sites 87-94, sibling pointer sites 101-111, guide exemption at 146
- `/home/user/skill-do-work/_dev/tests/audit-lockins.sh` — insert the new block before 151
- `/home/user/skill-do-work/_dev/tests/contracts/probe-lanes.sh` — 29 (already registered, do not touch)
- `/home/user/skill-do-work/_dev/tests/contract-regressions.sh` — 9 (`fast_contract_line_ceiling=77`, file is at 77)
- `/home/user/skill-do-work/_dev/tests/defensive-surface-audit.sh` — 34, 37, 40, 43, 46, 49

---
*Generated by Explore agent*