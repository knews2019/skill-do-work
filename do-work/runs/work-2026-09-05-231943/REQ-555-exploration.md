# REQ-555 Exploration — Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form

**Explored at:** HEAD = `fce57fc` (`[work run] REQ-592 pre-flight: green gate at 15e2ec3, baseline recorded, builder dispatched`)
**REQ file:** `/home/user/skill-do-work/do-work/queue/REQ-555-rewrite-the-prescribed-shell-guide-executable-homes-table-to-the-cli-route-form.md`
**Prior exploration:** none. `do-work/runs/work-2026-09-05-170806/` contains explorations for REQ-486, 552, 554, 556, 583, 587 and 591 — there is no `REQ-555-exploration.md`. Nothing to re-verify as stale from a prior pass; the staleness in this REQ comes from its own capture-time baseline, covered in §2.

**Route: A (direct) — blocked on REQ-554.** The work is mechanical and fully specified. The lock-in assertion is written and proven below. Do not start until REQ-554 is committed.

---

## 1. What the REQ asks for

Two edits to one shipped guide, plus one test assertion:

1. `skills/do-work/docs/prescribed-shell-primitives.md:9-17` — nine table rows whose "Canonical executable route" cell names a `*.sh` shim. Reword each to the `tools/do-work-cli.sh … <subcommand>` form that rows 18-22 already use.
2. `skills/do-work/docs/prescribed-shell-primitives.md:41` — delete the false clause claiming `scripts/protected-inventory.sh` orchestrates two check scripts.
3. `_dev/tests/audit-lockins.sh` — one assertion pinning shim rows in that table at 0.

---

## 2. Claim verification against HEAD

### C1 — Nine shim rows, 6-11 lines each, all exec into `do-work-cli.sh`. **HOLDS.**

`skills/do-work/docs/prescribed-shell-primitives.md:9-17`:

| Guide line | Route cell | Script lines | `do-work-cli` matches |
|---|---|---:|---:|
| 9 | `scripts/show-commit-diff.sh` | 9 | 2 |
| 10 | `scripts/add-local-git-exclude.sh` | 10 | 3 |
| 11 | `scripts/atomic-download.sh` | 9 | 2 |
| 12 | `scripts/capture-screenshot.sh` | 9 | 2 |
| 13 | `scripts/run-blocked-check.sh` | 11 | 2 |
| 14 | `scripts/protected-inventory.sh` | 6 | 2 |
| 15 | `scripts/stage-exact-deletion.sh` | 9 | 2 |
| 16 | `../../do-work-knowledge/scripts/lexical-memory-recall.sh` | 6 | 2 |
| 17 | `../../do-work-knowledge/scripts/install-memory-hooks.sh` | 6 | 2 |

Every one is a launcher. Two representative bodies:

- `skills/do-work/scripts/show-commit-diff.sh:9` — `exec bash "$script_directory/../tools/do-work-cli.sh" --format text show-commit-diff "$@"`
- `skills/do-work-knowledge/scripts/install-memory-hooks.sh:6` — `exec bash "$script_directory/../../do-work/tools/do-work-cli.sh" --format text install-memory-hooks "$@"`

Each carries the header comment `# do-work-cli compatibility launcher: <subcommand>` on line 2, which is the suite's own label for what these are.

### C2 — "seven rows … (9-11 lines each)". **DOES NOT HOLD.**

`skills/do-work/scripts/protected-inventory.sh` is **6** lines, not 9-11. The REQ's Detailed Requirements bullet (REQ line 38) is wrong; the REQ's What section (line 24) says "6-to-11-line", which is right. Cosmetic, but do not use "9-11" as a check.

### C3 — The orchestration clause is false. **HOLDS, and is more false than the REQ states.**

`skills/do-work/docs/prescribed-shell-primitives.md:41`:

> The complete secret-aware inventory and REQ association ship behind `scripts/protected-inventory.sh`, which orchestrates `tools/checks/uncommitted-inventory.sh` and `tools/checks/associate-files.sh` without duplicating their low-level logic.

Three facts kill it:

- `skills/do-work/scripts/protected-inventory.sh:6` is one `exec` line. It invokes nothing else.
- Both named scripts are themselves shims into the same CLI: `skills/do-work/tools/checks/uncommitted-inventory.sh:7-9` and `skills/do-work/tools/checks/associate-files.sh:19`. There is no orchestration relationship in either direction anywhere in shell.
- The composition lives in Go: `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go:331` (`handleProtectedInventory`), dispatched from `internal/corehelpers/commands.go:53`.

### C4 — Pointer target from 16 shipped files, ratchet-enforced. **HOLDS.**

`_dev/tests/prescribed-shell-canonicalization.sh:87-93` names 6 core sites; `:101-111` names 10 sibling sites. 6 + 10 = 16.

### C5 — "re-baseline the ratchet's heading and pointer counts". **DOES NOT HOLD for this REQ.**

The REQ's `write_set` includes `_dev/tests/prescribed-shell-canonicalization.sh`. It needs **no change**. That file holds no counts of anything:

- `:9-25` — 10 named scripts must exist and be executable. This REQ does not touch the shims (REQ Constraints, line 41), so unaffected.
- `:66-83` — a fixed list of 11 required headings. `## Shipped executable homes` is **not** in it.
- `:87-117` — two fixed pointer-site lists. A route reword changes no pointer.
- `:119-162` — stale-phrase and old-implementation scans that explicitly skip the canonical guide at `:146`.

Both tests are green at HEAD (`Prescribed shell primitive canonicalization checks passed.` / `Audit lock-in regressions passed.`, exit 0 each). The REQ's own Constraints say "no other test file changes" beyond the lock-in, which contradicts its `write_set`. **Follow the Constraints.**

### C6 — `audit-lockins.sh` exists, is executable, is registered. **HOLDS.**

156 lines, mode `-rwxr-xr-x`, registered at `_dev/tests/contracts/probe-lanes.sh:29` as `audit_lockins_probe`. It already carries four appended finding blocks — `:8` (Finding 10 / REQ-550), `:43` (Finding 5 / REQ-551), `:74` (Finding 8 / REQ-549), `:133` (Finding 2 / REQ-553). Appending a fifth before the exit check at `:151` is the established pattern.

### C7 — Expected net line delta −5. **DOES NOT HOLD.**

Rewording text inside 9 existing rows is net 0 lines. Deleting a clause inside the sentence on line 41 is net 0 lines. The realistic guide delta is **0**, plus about +22 lines in `_dev/tests/audit-lockins.sh`. The −5 comes from `do-work/audits/audit-2026-09-03.md:549` (R7) and is unreachable without deleting rows that the same audit's Remedy (`:384`) says to *reword*.

### C8 — Lock-in limit 0 (today 9). **Today's 9 confirmed; use the REQ's 0, not the audit's ≤9.**

The source audit says "shim rows ≤ 9 (red at 10)" at `do-work/audits/audit-2026-09-03.md:386` and `:549`. That threshold stays green after a no-op and proves nothing. The REQ file's own limit (REQ line 48: "0 after this REQ (today 9)") is the correct post-fix pin.

### C9 — All nine replacement subcommands exist. **HOLDS.**

`skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:28-38` registers `uncommitted-inventory`, `associate-files`, `protected-inventory`, `capture-screenshot`, `atomic-download`, `add-local-git-exclude`, `run-blocked-check`, `show-commit-diff`, `stage-exact-deletion`. `internal/knowledgecommands/commands.go:32-33` registers `install-memory-hooks` and `lexical-memory-recall`. Every new route cell names a real subcommand.

### C10 — The Reproduce command's `NR>=9 && NR<=22` range. **Correct at HEAD, unsafe after REQ-554.**

It reproduces exactly today. REQ-554 inserts a new shared section into this guide and will shift the table. Anchor on the heading instead (§5).

### C11 — The write_set is complete. **DOES NOT HOLD.**

`skills/do-work/docs/prescribed-shell-primitives.md` ships under `skills/`. Per `_dev/primes/prime-releases.md`: "A commit that changes shipped files under `skills/`, `tools/`, or `suite/` is a release." That requires a root `CHANGELOG.md` entry and a byte-identical copy to `skills/do-work/CHANGELOG.md`, enforced by `_dev/tests/shipped-package-reference-contract.sh`. Neither is in the write_set. Current version is 0.303.10 (`CHANGELOG.md:5`); root and mirror are identical at HEAD (`diff -q` clean).

### C12 — Will removing shims from the table orphan them? **No.**

Nothing else reads the table — grepping `Shipped executable homes` across `*.sh`, `*.go`, `*.md` hits only `prescribed-shell-primitives.md:5` and two changelog history lines. The caller-less-shim check at `_dev/tests/audit-lockins.sh:44-51` scans only `skills/do-work-toolbox/scripts/*.sh`, so the nine shims are outside its scope. They also keep live prose callers: `skills/do-work/actions/commit.md:56,97` and `skills/do-work-toolbox/actions/inspect.md:64,129`.

---

## 3. The REQ-554 dependency — what must land first

**REQ-554 is `status: pending`** (`do-work/queue/REQ-554-move-the-46-lines-commit-md-and-inspect-md-share-into-the-prescribed-shell-guide.md`). Of the batch, REQ-549/550/551/553 are already archived; 552, 554, 555, 556, 557, 558 are all still pending.

### What REQ-554 must have landed

1. **Its new shared section committed into `skills/do-work/docs/prescribed-shell-primitives.md`.** Until then the table's line numbers and the paragraph at line 41 are both unsettled.
2. **Its own assertion appended to `_dev/tests/audit-lockins.sh`.** Both REQs append before the same exit check at `:151`. Landing 554 first means 555 appends after 554's block with no textual conflict.
3. **Whatever it does to `_dev/tests/prescribed-shell-canonicalization.sh`.** REQ-554 may add its new heading to the required-heading list at `:66-83`. REQ-555 must not touch that file either way.

### The real overlap is a paragraph, not line numbers

This is the part the REQ understates. REQ-554 moves the manual fallbacks out of `commit.md` and `inspect.md` into this guide. The sentence REQ-555 must edit is the one that governs those fallbacks:

- `skills/do-work/docs/prescribed-shell-primitives.md:41` ends with "…their manual fallback must preserve `-uall`, NUL parsing, rename/copy consumption, secret quarantine, and the scripts' documented exit meanings."
- The fallbacks it points at are `skills/do-work/actions/commit.md:79,113` and `skills/do-work-toolbox/actions/inspect.md:88,145` — exactly the four sites REQ-554 relocates.

REQ-554 will very likely rewrite this whole paragraph. Running REQ-555 first guarantees a conflict on it.

### What is contingent on 554's outcome

- **Edit (2) may become a no-op.** If REQ-554 rewrites line 41 and drops the orchestration clause as a side effect, REQ-555 has nothing to delete there. After 554 lands, run `grep -n 'which orchestrates' skills/do-work/docs/prescribed-shell-primitives.md`. Empty result means skip edit (2). **Do not re-add the clause in order to have something to delete.**
- **The replacement wording for line 41 depends on what 554 leaves.** The text proposed in §4 assumes the paragraph survives roughly as-is. If 554 restructured it, adapt: the requirement is only that no sentence claims a shim orchestrates other scripts.
- **Edit (1) is not contingent.** The table rows at 9-17 are untouched by REQ-554's subject matter. Only their line numbers may move, which the heading-anchored lock-in already handles.

---

## 4. Exact changes

### 4.1 `skills/do-work/docs/prescribed-shell-primitives.md` lines 9-17 — route column

Replace the route cell only. Leave the "Mechanics owned" column byte-identical. Use U+2026 for the ellipsis (bytes `E2 80 A6`, confirmed against line 18) so the new rows match rows 18-22 exactly.

| Old route cell | New route cell |
|---|---|
| `scripts/show-commit-diff.sh` | `tools/do-work-cli.sh … show-commit-diff` |
| `scripts/add-local-git-exclude.sh` | `tools/do-work-cli.sh … add-local-git-exclude` |
| `scripts/atomic-download.sh` | `tools/do-work-cli.sh … atomic-download` |
| `scripts/capture-screenshot.sh` | `tools/do-work-cli.sh … capture-screenshot` |
| `scripts/run-blocked-check.sh` | `tools/do-work-cli.sh … run-blocked-check` |
| `scripts/protected-inventory.sh` | `tools/do-work-cli.sh … protected-inventory` |
| `scripts/stage-exact-deletion.sh` | `tools/do-work-cli.sh … stage-exact-deletion` |
| `../../do-work-knowledge/scripts/lexical-memory-recall.sh` | `tools/do-work-cli.sh … lexical-memory-recall` |
| `../../do-work-knowledge/scripts/install-memory-hooks.sh` | `tools/do-work-cli.sh … install-memory-hooks` |

The two knowledge rows collapse to the same `tools/do-work-cli.sh` prefix as the rest, and that is correct: there is one CLI. `skills/do-work-knowledge/scripts/install-memory-hooks.sh:6` reaches it as `../../do-work/tools/do-work-cli.sh`, the same binary the guide addresses as `tools/do-work-cli.sh` from its own `skills/do-work/` root.

### 4.2 `skills/do-work/docs/prescribed-shell-primitives.md` line 41 — orchestration clause

**Do not do the bare minimal deletion.** Cutting only "which orchestrates … low-level logic" leaves the trailing phrase "the scripts' documented exit meanings" with no antecedent — the two check scripts were named *only* in the deleted clause.

Proposed replacement (net 0 lines):

> The complete secret-aware inventory and REQ association ship behind `tools/do-work-cli.sh … protected-inventory`, which `scripts/protected-inventory.sh` still launches for compatibility. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` invoke that route; their manual fallback must preserve `-uall`, NUL parsing, rename/copy consumption, secret quarantine, and the command's documented exit meanings.

This keeps the retained shim visible (the 0.260.1 decision retains it and `commit.md:56,97` still invokes it), states the true owner, and fixes the dangling plural.

### 4.3 `_dev/tests/audit-lockins.sh` — one assertion

Append immediately before line 151 (`if [ "$failure_count" -gt 0 ]; then`). Text in §5.

### 4.4 `CHANGELOG.md` + `skills/do-work/CHANGELOG.md` — release entry

Not in the write_set; required. Bump from 0.303.10 (or from whatever REQ-554 leaves), newest on top, title says what was delivered with no codename, verified not already used. Then copy root to mirror byte-identically.

### 4.5 `_dev/tests/prescribed-shell-canonicalization.sh` — drop from write_set

No change needed (§2, C5).

---

## 5. Lock-in assertion — paste-ready

Anchored on the heading, not on line numbers, because REQ-554 shifts this table. Reads field 2 of the pipe split so only the route cell is judged; the Mechanics column may legitimately name a file later. This is the "state conditions, not lists" form — it keys on "a route cell ending in `` .sh` ``", so a tenth shim row added tomorrow is caught without anyone maintaining a list.

```bash

# Finding 7: stale-shell-ownership-prose (REQ-555)
# The "Shipped executable homes" table names the owning executable route. Every mechanic
# listed there moved into Go at 0.260.1, so a route whose backticked token ends in `.sh`
# points at a retained compatibility shim rather than at the owner. Anchored on the
# heading rather than line numbers, because neighbouring REQs edit this guide and shift
# the table. Only the route cell is read; the mechanics cell may legitimately name a file.
executable_homes_guide="$repo_root/skills/do-work/docs/prescribed-shell-primitives.md"
shim_route_rows="$(awk '
  /^## Shipped executable homes$/ { inside_homes_table = 1; next }
  inside_homes_table && /^## / { inside_homes_table = 0 }
  inside_homes_table && /^\|/ {
    split($0, table_cells, "|")
    if (table_cells[2] ~ /\.sh`/) printf "%d:%s\n", NR, table_cells[2]
  }
' "$executable_homes_guide")"
if [ -n "$shim_route_rows" ]; then
  while IFS= read -r shim_route_row; do
    [ -z "$shim_route_row" ] && continue
    printf 'FAIL: executable-homes table routes to a compatibility shim instead of a do-work-cli subcommand: %s\n' \
      "$shim_route_row" >&2
    failure_count=$((failure_count + 1))
  done <<< "$shim_route_rows"
fi
```

### Why this shape

- Matches the house pattern of all four existing blocks: compute into a variable, `if [ -n … ]`, `while IFS= read -r` over a `<<<` herestring, `printf 'FAIL: …' >&2`, increment `failure_count`.
- Names are two words and grep-findable per `skills/do-work/crew-members/coding-guardrails.md` §5: `executable_homes_guide`, `shim_route_rows`, `shim_route_row`, `inside_homes_table`, `table_cells`. No cryptic locals.
- Safe under the file's `set -uo pipefail` (line 3, no `-e`): no pipelines whose early exit could be misread, and the `awk` is the only command in the substitution.
- Emits the offending line number and cell, so a failure is actionable without re-deriving anything.

### Proof (run read-only in the scratchpad, no repo file touched)

| Check | Result |
|---|---|
| RED — block spliced into a copy of `audit-lockins.sh`, `repo_root` pinned, run against HEAD's guide | 9 `FAIL:` lines naming guide lines 9-17, **exit 1** |
| GREEN — same copy with `executable_homes_guide` pointed at a generated post-fix guide | `Audit lock-in regressions passed.`, **exit 0** |
| `shellcheck` on the spliced copy | same three pre-existing findings as the untouched baseline (SC2126 ×2 at lines 28/48, SC2016 at line 118), **zero new** |
| Standalone counter at HEAD | prints `9` |
| Standalone counter on simulated fix | prints `0` |

Standalone counter, for the claim-time re-verification (use this, not the REQ's `NR>=9 && NR<=22` awk):

```bash
awk '
  /^## Shipped executable homes$/ { inside_homes_table = 1; next }
  inside_homes_table && /^## / { inside_homes_table = 0 }
  inside_homes_table && /^\|/ { split($0, table_cells, "|"); if (table_cells[2] ~ /\.sh`/) shim_row_count++ }
  END { print shim_row_count + 0 }
' skills/do-work/docs/prescribed-shell-primitives.md
```

### Should the assertion also pin the orchestration sentence?

No. The REQ Constraints say one assertion, and the GREEN criterion names only the shim-row count. The sentence deletion is proved by `grep -n 'which orchestrates' …` returning nothing. Note that the canonicalization ratchet's `stale_patterns_file` (`_dev/tests/prescribed-shell-canonicalization.sh:121-131`) is *not* a home for it — that scan skips the canonical guide at `:146`, so a pattern added there would never look at this file.

---

## 6. Full verification commands

```bash
# Before (RED)
bash _dev/tests/audit-lockins.sh                        # 9 FAILs, exit 1

# After (GREEN)
bash _dev/tests/audit-lockins.sh                        # passes, exit 0
bash _dev/tests/prescribed-shell-canonicalization.sh    # passes, exit 0, file unchanged
grep -n 'which orchestrates' skills/do-work/docs/prescribed-shell-primitives.md   # no match
shellcheck _dev/tests/audit-lockins.sh                  # baseline findings only
diff -q CHANGELOG.md skills/do-work/CHANGELOG.md        # identical
git diff --stat
```

Baseline at HEAD: both probes green, root and mirror changelogs identical.

---

## 7. Risks and scope boundaries

**R1 — REQ-554 is not landed.** Blocking. The overlap is the paragraph at line 41, not just line numbers (§3).

**R2 — Edit (2) may be a no-op after 554.** Re-check with `grep -n 'which orchestrates'`. Do not re-add the clause to have something to delete.

**R3 — Stale line range in the REQ's Reproduce.** `NR>=9 && NR<=22` is correct today and wrong after 554. Use the heading-anchored form.

**R4 — Missing changelog in the write_set.** The guide is shipped; this commit is a release. Root entry plus byte-identical mirror.

**R5 — write_set contradicts Constraints on `prescribed-shell-canonicalization.sh`.** The Constraints win; that file needs no change.

**R6 — Table and actions will name different routes.** After the rewrite, the table says `tools/do-work-cli.sh … protected-inventory` while `commit.md:56,97` and `inspect.md:64,129` still prescribe `scripts/protected-inventory.sh`. Already reconciled by the guide's own intro at line 3: "retained scripts remain compatibility/parity surfaces only where stated." Point at that line; write no new reconciliation prose and do not touch the actions. The REQ explicitly retains the shims (REQ line 41).

**R7 — Table and guide body will differ in form.** Body sections still name shim paths as the invocation: line 45 (`scripts/show-commit-diff.sh <commit>`), line 64 (`scripts/add-local-git-exclude.sh`), line 80 (`scripts/atomic-download.sh`). Out of scope — Constraints forbid nearby fixes. Record as a discovered task.

**R8 — The audit's weaker lock-in.** `do-work/audits/audit-2026-09-03.md:386,549` says "≤ 9 (red at 10)", which stays green after a no-op. Follow the REQ file's 0.

**R9 — Temptation to collapse the column.** After the rewrite all 14 route cells share the identical `tools/do-work-cli.sh … ` prefix, inviting a bare-subcommand column with a lead-in sentence. That restructures 5 rows this finding does not name, against the audit Remedy's "reword the route column … the form the toolbox rows already use". Keep the literal form; log the collapse as a discovered task.

---

## 8. Effort

Matches `effort-mechanical`. Nine cell rewrites, one sentence rewrite, one ~22-line test block, one changelog entry plus mirror. All decisions are made above; the assertion is written and proven. The only real cost is waiting for REQ-554.

---
*Generated by Explore agent*