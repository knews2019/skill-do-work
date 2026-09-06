# Builder brief — REQ-552, replace two coreutils exec sites with pure Go

**Worktree:** created for you; the branch is checked out there. Work only in it.
**Route:** B. **TDD: yes.** **Impact: negligible.** **Estimate P50: 35 active minutes.**

## Read first

1. `do-work/runs/work-2026-09-05-231943/REQ-552-exploration.md` in the MAIN checkout — the evidence, with file:line
   anchors, the measured behaviour of both replacements, and ten risks. It is the authority.
2. `do-work/working/REQ-552-replace-two-coreutils-exec-sites-with-the-pure-go-the-package-already-has.md`
   in your worktree — the request, its `## Triage`, `## Exploration` and `## Scope`.
3. `_dev/primes/prime-shell-commands.md` — binding before you touch `_dev/tests/audit-lockins.sh` or
   either prescribed-shell case. The trap list and Closed Enumerations Go Stale.
4. Crew rules in `skills/do-work/crew-members/`: general.md, shared-principles.md,
   coding-guardrails.md, communication-style.md, testing.md, backend.md.

## The two code edits

**`internal/corehelpers/commands.go:724-728`** — replace the `exec.Command("find", archiveRoot,
"-name", "REQ-*.md", "-print0")` + `CombinedOutput` block with a `filepath.WalkDir` readability probe
that records the first traversal error and returns it as the evidence string. Keep the `os.Stat` guard
at 717-722 exactly as it is. **Do not filter to `REQ-*.md`** — the probe is about readability, not
about matching.

**`internal/toolboxcommands/architecture.go:133-135`** — replace the
`exec.Command("cp", draftPath, stagedPath).CombinedOutput()` block with
`os.WriteFile(stagedPath, data, 0o600)`, keeping the `draft copy failed: ` evidence prefix. Measured:
`cp src dst` onto a pre-created 0600 file leaves the mode at 0600 and copies the bytes, and
`os.CreateTemp` creates at 0600, so the draft's own mode is never copied today. `os.WriteFile` at
0o600 reproduces that exactly.

Neither replacement the request names is usable as named — `inventory.go`'s walk parses and builds an
ownership map, and `copyLast30DaysTree` cannot be called for a regular file with a pre-created target.
Write both fresh from the standard library. The exploration has the details.

## The lock-in — paste it, do not improvise it

Insert after line 149 of `_dev/tests/audit-lockins.sh` (the closing `fi` of the Finding 2 block) and
before line 151. The exact block, with its comment, is in the exploration's `red_green_recipe`. Two
things about it are decided and must not be changed: it keeps `--glob '!*_test.go'` (without it the
pattern matches `suiteinstall/update_transaction_test.go:25`, which spawns `cp -R` in fixture setup
and is not this finding, so the lock-in would be red on day one), and its command list stays
byte-identical to the audit's own Reproduce pattern.

## The two fixture cases you MUST also fix — this is the whole reason this is not a five-minute job

Both were proven inert against a patched binary. Rewrite each to drive the in-process failure
directly. **Never delete them** — deleting the first leaves the archive-walk failure path with no
coverage at all.

- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh:220-236` installs a fake `find` that
  exits 3 on `PATH` and asserts the run exits non-zero and never prints "audit clean". Measured: base
  exit 1 with "the archive walk failed"; patched exit 0 with "archive audit clean (1 file(s) scanned)".
  Both assertions at 234 and 236 flip.
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-215` writes a failing `cp`
  into a directory prefixed onto `PATH`. Measured: base exit 1 with no `index.html`; patched exit 0
  with one created. Three `fail_case` lines fire at 209, 210 and 214.

An in-process failure is producible without a PATH shim — make the target unreadable, or point at a
path the process cannot write. Choose the narrowest mechanism that makes the case assert what it names.

## Testing

These fixture cases are HEAVY tier only: `_dev/tests/prescribed-shell-harness.sh:11-14` refuses unless
`DO_WORK_MAINTAINER_TIER=heavy`. **A green fast gate is not evidence this change is clean.** Run the
prescribed-shell lane directly at heavy tier as your iteration loop, and record its exit line.

```
env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-scripts-behavior.sh
bash _dev/tests/audit-lockins.sh
go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/corehelpers/ ./internal/toolboxcommands/
```

**The corehelpers tests only pass in tree.** `inventory_test.go:136` resolves a sibling
`checks/uncommitted-inventory.sh` relative to the module directory, so a copy of the module built
elsewhere fails 59 subtests for reasons unrelated to any change. Test in the worktree, never in a copy.

## RED then GREEN

RED, at your base revision, prints two lines and exits 0:
```
rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
  skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'
```
GREEN: the same command prints nothing; `bash _dev/tests/audit-lockins.sh` prints
"Audit lock-in regressions passed." and exits 0; and the lock-in block run standalone against your
base revision must print two FAIL lines first — show that before you accept it.

## Environment

Fresh cloud container. Wrap every gate and every Go test:

```
env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL=/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gitconfig-gate \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  <command>
```

Capture exit status from `$?` directly — never pipe a gate to `tail`. This machine has 4 CPUs and
another builder is running: run the full canonical gate at most once, and prefer the targeted commands
above.

## Out of scope, record but do not fix

`commands.go:721-722` — when `do-work/archive` exists but is not a directory, `err` is nil and the
function returns the evidence string `<nil>`. Non-empty, so the gate still fires, but the text is
wrong. Report it as a discovered task.

## Hand-back

Commit on your branch in the worktree, message ending
`Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>`. Write your report to the ABSOLUTE path
`/home/user/skill-do-work/do-work/runs/work-2026-09-05-231943/REQ-552-handback.md` in the MAIN checkout. Do not stage or
commit it, do not push, do not touch any other worktree.
