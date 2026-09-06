# REQ-554 Remediation Hand-Back

REQ-554 is the request that moved the body `skills/do-work/actions/commit.md` and
`skills/do-work-toolbox/actions/inspect.md` shared — the inventory tag legend, the
secret-shaped matching rule, the four file-reading bullets, the association semantics and
both by-hand fallbacks — into `skills/do-work/docs/prescribed-shell-primitives.md`. The
review scored it 73%, Partial: the prose move was accepted, the lock-in that protects it was
not. This hand-back closes the six findings the remediation brief named.

**Branch:** `worktree-agent-REQ-554-move-commit-inspect-shared-body`
**Head:** `cfcf63e4c78200ddf8a2bf501d315c2475b4fb6e`
**Files changed:** `_dev/tests/audit-lockins.sh`, `skills/do-work/docs/prescribed-shell-primitives.md`
(2 files, 68 insertions, 28 deletions). No action file was touched; the accepted prose move
stands as merged. No release bump — that belongs to finalization.

## The choice F2 asked for: the similarity count is deleted, not re-baselined

F2 offered three shapes: honest headroom, a non-scaffold measure, or dropping the count.
Dropped, for three reasons.

1. **Headroom hides the false trips instead of removing them.** The four proven trips move
   the number by 2 or 3. Any headroom that absorbs them also absorbs the same amount of real
   duplication, and REQ-555 through REQ-558 edit these two files next.
2. **Excluding scaffold lines means hand-maintaining a list of template lines.** That is the
   exact shape `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale bans.
3. **The count caught nothing the sentence guard misses.** It never saw a one-sided return
   (F1). It never saw a paraphrase — the review proved that with its own case, which scored
   30 and exited 0. What was left was verbatim two-sided duplication of text nobody named,
   and finding that is the maintainability audit's job. This file pins findings the audit
   already made.

In its place: each action file is scanned on its own for one load-bearing sentence out of
each moved passage — two tag-legend rows, the case-insensitivity rule, the first and last
reading bullets, and two association rules. Seven phrases, `rg --fixed-strings`, exit status
read rather than emptiness. Because each file is scanned separately, a paste into either one
alone fails.

The comment above the block now records what the check does not catch (a paraphrase) instead
of claiming, as it did, that "Returning shared prose trips this".

## Findings closed

### F1 — the assertion missed the one-sided return it was named for

Fixed by the per-file phrase scan. Proved against the shipped
`_dev/tests/audit-lockins.sh`, mutating the real files and restoring them:

| Case | Old count | New result |
|---|---|---|
| Untouched tree | 30, exit 0 | `Audit lock-in regressions passed.` exit 0 |
| Legend + four reading bullets pasted into `commit.md` alone | 30, exit 0 | 4 FAIL lines naming `commit.md:237,241,242,245`, exit 1 |
| Same passage pasted into `inspect.md` alone | not covered | 4 FAIL lines naming `inspect.md:425,429,430,433`, exit 1 |
| Whole pre-move `commit.md` restored (`git show b2ba3ea2:…`) | 29, exit 0 | 8 FAIL lines — 6 phrase sites plus both by-hand fallbacks, exit 1 |
| Restored | 30, exit 0 | exit 0 |

### F2 — zero headroom, and the metric rose when a line was deleted

Fixed by deletion. Fuzz run as the brief required, over the REQ-554 block extracted verbatim
by line range from the shipped file and pointed at a scratch copy of both actions:

**Every single-line deletion across both action files: 660 runs, 0 false trips.**

That covers all four deletions the reviewer found red — `inspect.md:46`, and the fenced-block
opening lines at `inspect.md:113`, `commit.md:47` and `commit.md:55`.

Two structural additions the reviewers used were run separately, both green:
- an identical `### Step 7: Cleanup` heading plus one body line appended to both actions
  (scored 32 under the old ceiling) — exit 0
- an identical Error Handling table row appended to both actions — exit 0

### F3 — "glob both directories" named no directory

`skills/do-work/docs/prescribed-shell-primitives.md:74` now reads
`glob \`do-work/archive/**/REQ-*.md\` and \`do-work/working/REQ-*.md\``. Those are the two
roots `AssociateProjectPaths` walks (`internal/corehelpers/inventory.go:253`), and they are
the same two globs both actions already print in their Step 3 prose.

### F4 — the `git rev-parse --git-dir` claim was wrong

Verified rather than assumed. Built the CLI and ran both modes in a directory that is not a
Git repository:

```
mode=start     exit=2  finding HELPER-USAGE [error]: git rev-parse --git-path do-work-commit-secret-quarantine: exit status 128
mode=associate exit=2  finding HELPER-USAGE [error]: git rev-parse --git-path do-work-commit-secret-quarantine: exit status 128
```

`--git-dir` appears nowhere on that path. Guide line 45 now reads: "Both modes resolve the
run-level quarantine file with `git rev-parse --git-path <quarantine name>` and exit 2 when
that resolution fails, which is what happens outside a Git repository." That matches the code
(`inventory.go:357`, `usageResult` → `OutcomeFailure` → `ExitCode` 2 at
`resultmodel/result_model.go:658`), the executed run above, and what both actions already tell
their readers about exit 2.

### F5 — the FAIL message gave the wrong remedy

Resolved by removal. The message that offered "move the shared body into the guide" for a
legitimately added template heading is gone with the count. Every remaining failure branch
names a remedy that fits the case that fired:

- a moved passage came back → `<file>:<line> restates prose that is canonical in
  skills/do-work/docs/prescribed-shell-primitives.md#protected-inventory-fallbacks; cite that
  section instead.`
- a by-hand fallback came back → the same message as before, now naming where the
  protected-inventory by-hand procedures live.
- an action file moved → `shared-body lock-in cannot scan a missing action file: <path>`
- the scanner could not run → `could not scan <path> for moved shared prose (rg exit N)`

### F6 — the 999 sentinel printed a count nobody measured

Resolved by removal. The block no longer uses `python3` at all, so there is no sentinel and no
share count to misreport. Proved by running the block with `rg` off `PATH`: 15 FAIL lines,
every one naming the tooling failure and its exit status, exit 1. With one action file renamed
away: `FAIL: shared-body lock-in cannot scan a missing action file:
skills/do-work/actions/commit.md`, exit 1.

## Verification

| Check | Result |
|---|---|
| `bash _dev/tests/maintainer-verify.sh` (canonical gate, run once) | exit 0, 107s wall |
| `bash _dev/tests/audit-lockins.sh` | exit 0 — `Audit lock-in regressions passed.` |
| `bash _dev/tests/prescribed-shell-canonicalization.sh` | exit 0 |
| `bash _dev/tests/action-shell-blocks.sh` | exit 0 — 74 fenced blocks, 33 shipped shell files |
| `bash _dev/tests/shipped-package-reference-contract.sh` | exit 0 |
| `shellcheck --severity=warning _dev/tests/audit-lockins.sh` | exit 0 |
| Single-line-deletion fuzz, both actions | 660 runs, 0 false trips |

Every exit status was read from `$?` directly.

## Left for the reviewer

- **The version bump is still owed** and is deliberately not in this commit. It lands in the
  `[REQ-554] complete` finalization commit, at whatever the patch is by then.
- **Not touched, and still open from the review:** the paraphrase gap (now written into the
  comment rather than fixed), the "Read each tag class this way, and no other way" absolute in
  the guide, `commit.md`'s dropped "is the bug in the Red Flags below" clause, and the
  hand-back's wrong `near_identical_cross_file_pairs` citation. All were rated Minor or Nit
  and report-only, and none is in the five-file scope this remediation was given.
- **One correction to the prior hand-back, for the record:** its Qualification section cited
  `near_identical_cross_file_pairs` being 0 as independent proof the duplication was gone.
  The reviewer is right that the metric lives in `_dev/tests/contracts/core-checks.sh` and only
  compares `## Common Rationalizations` rows, so it proved nothing about this change. The
  evidence that stands is the red-before/green-after table under F1 above.
