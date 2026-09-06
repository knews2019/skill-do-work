# REQ-554 hand-back — commit/inspect shared body moved into the prescribed-shell guide

**Branch:** `worktree-agent-REQ-554-move-commit-inspect-shared-body`
**Branch head:** `d26093785dc8396863d56db62a7bca8309547b91` (based on `2876d96`)
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-554-move-commit-inspect-shared-body`
**Nothing pushed. No worktree removed. Nothing staged or committed in the main checkout.**

## Files changed (exactly the five in scope)

- `skills/do-work/docs/prescribed-shell-primitives.md` — new `## Protected inventory fallbacks`
  section inserted after the `## Per-file untracked inventory` pointer paragraph, before
  `## Merge-aware commit diff`. Line 41 rewritten to point at it instead of restating the contract;
  the `orchestrates` sentence REQ-555 will remove is left untouched.
- `skills/do-work/actions/commit.md` — moved prose deleted, four pointer citations at
  `../docs/prescribed-shell-primitives.md#protected-inventory-fallbacks`. 260 → 236 lines.
- `skills/do-work-toolbox/actions/inspect.md` — the mirror deletion, citations at
  `../../do-work/docs/prescribed-shell-primitives.md#protected-inventory-fallbacks`. 448 → 424 lines.
- `_dev/tests/prescribed-shell-canonicalization.sh` — one line: `'## Protected inventory fallbacks'`
  added to the required-heading membership list. Nothing re-baselined, because nothing there counts.
- `_dev/tests/audit-lockins.sh` — new `# Finding 6: commit-inspect-shared-body (REQ-554)` block
  inserted immediately before the final `failure_count` gate, beside REQ-552's Finding 9 block,
  which is untouched.

## The measured ceiling: 30

The request asked for 10. That is below the floor and would have been red forever. The measured
count after the move is **30**, and here is every one of the 30 lines, by kind:

| Count | Kind | Why it cannot move |
|---|---|---|
| **17** | Template scaffold | `_dev/primes/prime-action-files.md:70-80` requires both actions to carry them: `## When to Use`, `**Use when:**`, `## When This Runs`, `## Steps`, three bare code fences, `### Step 1: Preflight`, `### Step 4: Group Unassociated Files`, `### Step 6: Report`, `## Error Handling`, `\| Situation \| Action \|`, `\|-----------\|--------\|`, `## What This Action Does NOT Do`, `## Rules`, `## Red Flags`, `## Verification Checklist` |
| **4** | Structural rows | Three rows of each action's own ASCII flow diagram (two bare `│`, plus `├── Group Unassociated ── semantic clustering (1-5 files per group)`) and one row of each action's own Error Handling table (`\| Not a git repo \| Report "Not a git repository" and exit \|`) |
| **8** | Step 4 semantic clustering algorithm | Caller policy, not a shell primitive. The guide's charter at `prescribed-shell-primitives.md:3` says the guide owns the shared primitive, not the action's policy. Moving it would buy 8 lines on the metric by weakening the guide |
| **1** | `Files that come back \`-\` remain unassociated and move to Step 4.` | Routes to each action's own step number — action policy for the same reason |

17 + 4 + 8 + 1 = **30**. The reason it is not zero is written into the assertion's comment, so the
next reader does not try to ratchet it down into the scaffold.

This is one line above the exploration's suggested S2 of 29. The difference is the single routing
line above, which the exploration counted as moved; it stays local because "move to Step 4" is a
statement about the caller's own step numbering.

## RED proof, before the assertion was accepted

The new assertion block was extracted verbatim from `_dev/tests/audit-lockins.sh` into a standalone
script and run with `repo_root` pointed at a tree holding `git show 2876d96:` copies of both action
files. Output:

```
FAIL: commit.md and inspect.md share 46 identical lines; ceiling is 30. Move the shared body into skills/do-work/docs/prescribed-shell-primitives.md.
FAIL: manual "do it by hand" fallback remains in a shipped action: skills/do-work-toolbox/actions/inspect.md:88:...
FAIL: manual "do it by hand" fallback remains in a shipped action: skills/do-work-toolbox/actions/inspect.md:145:...
FAIL: manual "do it by hand" fallback remains in a shipped action: skills/do-work/actions/commit.md:79:...
FAIL: manual "do it by hand" fallback remains in a shipped action: skills/do-work/actions/commit.md:113:...
RED: 5 failures
exit=1
```

The four named sites are the proof the glob is live. The prior exploration's `--glob '*/actions/*.md'`
printed `0` here with the same four violations present, because a single `*` does not cross a
directory separator. This assertion uses `**/actions/*.md`, and reads `rg`'s own exit status rather
than a piped `awk` total, per `_dev/primes/prime-shell-commands.md` § Unchecked Exit Status Reads as
Content — an `awk` total prints `0` both when nothing matched and when the scan never ran.

## GREEN evidence

Every command wrapped in the sanitized env from the brief, exit status read from `$?` directly, never
piped.

```
Audit lock-in regressions passed.                                   EXIT[audit-lockins]=0
Prescribed shell primitive canonicalization checks passed.          EXIT[prescribed-shell-canonicalization]=0
Shell-block lint passed: 74 fenced blocks and 33 shipped shell files; ShellCheck enabled.
                                                                    EXIT[action-shell-blocks]=0
Defensive-surface exact deletion regressions passed.                EXIT[defensive-surface-audit]=0
Contract regression checks passed.                                  EXIT[contract-regressions]=0
Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).
                                                                    EXIT[prescribed-shell-scripts-behavior]=0   (DO_WORK_MAINTAINER_TIER=heavy)
Maintainer verification passed.                                     EXIT[maintainer-verify]=0   (gate wall 87s)
```

The canonical gate was run exactly once, as the brief asked.

`contract-regressions.sh` also reported `near_identical_cross_file_pairs 0`.

## Decisions taken, and the three request errors resolved

- **D-01 — ceiling is the measured 30, not the request's 10.** Ten is below the 17-line scaffold
  floor. Recorded in the commit message and in the assertion's own comment.
- **D-02 — nothing re-baselined in `prescribed-shell-canonicalization.sh`.** It has zero numeric
  constants; every check is a `grep -Fqx` / `grep -Fq` membership assertion. Adding a heading or a
  pointer cannot fail it. One heading added to the membership list so the new home cannot be quietly
  deleted.
- **D-03 — the Step 4 clustering algorithm stays in both actions.** It is semantic file grouping, not
  a shell primitive; the guide's charter excludes it. This is the deliberate reason the ceiling is 30
  rather than 22.
- **D-04 — the X and XD wordings that differ between the two actions are preserved.** The tag legend
  in the guide states classification only (`X` — non-deleted excluded path; `XD` — deleted
  secret-shaped path). What each action *does* with those rows stays in each action's own paragraph,
  which already carried the full mode-specific rule: commit.md keeps "only the deletion may proceed",
  inspect.md keeps "this read-only action inspects only the path and deletion state". The guide says
  in one sentence that this is caller policy it does not own.
- **F6 from the exploration is resolved, not worked around.** The association fallback's dangling
  `accepting every alias above` back-reference works again: the terminal-success alias bullet moved
  into the same guide section, so "above" now resolves. The wording was tightened to
  `accepting every terminal-success alias listed above`.

## Verification of no collateral drift

Each moved sentence now has exactly one home in the whole shipped tree. Grepped across `skills/`:
`Secret-shaped matching is case-insensitive`, `non-deleted excluded path`, `Partial matches count`,
`Terminal-success matching honors`, and ``Read the `git diff` for each file`` each return exactly one
file, `skills/do-work/docs/prescribed-shell-primitives.md`.

None of the six phrases `_dev/tests/defensive-surface-audit.sh` pins on these two files was
reintroduced (that test exits 0). The canonicalization stale-pattern scan skips the canonical guide
at line 146, so the moved prose is exempt, and none of its 16 patterns matches a moved line anyway.

## Discovered, not fixed

- The Step 4 semantic clustering algorithm is still duplicated between the two actions — 8 identical
  lines, the largest remaining block. It needs a home that is not the prescribed-shell guide, and
  that is a separate request.
- The two actions' `Start a run-level quarantine` paragraphs are near-identical, differing only in
  the trailing example sentence. Not in this REQ's scope and not counted by the difflib metric
  (the difference breaks the matching run).
- `_dev/tests/contract-regressions.sh` is exactly 77 lines against its own
  `fast_contract_line_ceiling=77`. Zero headroom. Untouched here, but the next person who adds a line
  to it breaks the fast gate.
