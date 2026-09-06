# Builder brief — REQ-556, cut the debug-artifact rule prose that qualify already enforces

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-556-cut-debug-artifact-prose`
**Branch:** `worktree-agent-REQ-556-cut-debug-artifact-prose`, based on `d3ceca3`
**Route:** B. **TDD: yes.** **Impact: negligible.** **Estimate P50: 20 active minutes.**

## Read first

1. `do-work/runs/work-2026-09-05-231943/REQ-556-exploration.md` in the MAIN checkout — **the authority.** It has the
   site-by-site list, the exact edit for each, the paste-ready lock-in, and ten claims from the request
   checked against HEAD with the verdict for each.
2. `do-work/working/REQ-556-cut-the-debug-artifact-rule-prose-that-qualify-already-enforces.md` in your
   worktree — the request, its `## Triage`, `## Exploration` and `## Scope`.
3. `_dev/primes/prime-action-files.md` (you are editing three shipped action files) and
   `_dev/primes/prime-shell-commands.md` (you are adding a shell assertion).
4. Crew rules in `skills/do-work/crew-members/`: general.md, shared-principles.md,
   coding-guardrails.md, communication-style.md, anti-slop.md.

## Do not follow the request literally

Its baseline does not survive contact with HEAD. Its headline site count is wrong, the commit its
Reproduce line names cannot be checked out, and the sentence it says to *keep* does not exist in
shipped prose today — so keeping it is an addition, not a retention. The exploration lists all ten
checked claims with verdicts. Where they disagree, the exploration wins.

**Line numbers in the exploration may have moved.** `work-reference.md` was edited by REQ-486 earlier
in this run. Locate every edit by its text, never by its line number, and say in the hand-back which
anchors had moved.

## The edits

Four in `skills/do-work/actions/work.md`, one in `skills/do-work/actions/review-work.md`, one row in
`skills/do-work/actions/work-reference.md`, one assertion block in `_dev/tests/audit-lockins.sh`. The
exploration gives each one exactly.

**Two mentions must survive, and cutting either to reach a smaller number is the failure mode here:**
`review-work.md`'s standalone-review hygiene bullet is a read the canonical `qualify` never makes, and
the emitted P-A-U template payload is byte-identical across four shipped files. Neither is a
restatement of the rule.

**The rule itself lives in code** — `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and
`QUALIFY-UNIFY-DISARMED` in `internal/corehelpers/checks.go`. That file is out of scope and does not
change; it is what makes the prose copies restatements.

## The lock-in

Paste it from the exploration. Two properties are decided and must not be changed: it **counts** rather
than name-lists, because a new restatement is the regression whatever words it uses; and a missing or
renamed target file **fails loudly** instead of counting zero. `_dev/tests/audit-lockins.sh` already
carries blocks from REQ-552 and REQ-554 — add yours beside them and disturb neither.

Prove it red before you accept it: run your block standalone against your base revision and show it
reporting the pre-cut count above the ceiling. Then show it green after, and show it red again when a
restatement is pasted back into any one of the three files.

## Testing

```
bash _dev/tests/audit-lockins.sh
bash _dev/tests/action-shell-blocks.sh
bash _dev/tests/contract-regressions.sh
bash _dev/tests/prescribed-shell-canonicalization.sh
env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh
```

Record each exit line. `contract-regressions.sh` matters because it pins content in these very action
files; if a deletion trips it, that is a signal the prose was load-bearing, not noise.

## Environment

```
env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL=/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gitconfig-gate \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  <command>
```

Capture exit status from `$?` directly. Other builders are running on this 4-CPU machine: prefer the
targeted probes above and run the full canonical gate at most once.

## Do not do the release

The version bump and changelog belong to finalization, and the number depends on work still in flight.

## Hand-back

Commit on your branch, message ending `Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>`. Write
your report to `/home/user/skill-do-work/do-work/runs/work-2026-09-05-231943/REQ-556-handback.md` in the MAIN checkout.
Do not stage or commit it, do not push, do not touch any other worktree.
