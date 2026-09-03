# Run With Recovery Action

> **Part of the do-work skill.** Recovers this queue under the user's explicit sole-writer assertion, then completes `actions/work.md` with the original arguments. It belongs in core because it resolves core finalization and claim ownership before handing control back to core's work pipeline.

## When to Use

**Use when:**

- A previous `do-work run` was interrupted and the user asserts that this checkout is now the queue's only writer and releaser.
- Ordinary `run` refused ambiguous shared finalization metadata or left stale claims for a human ownership decision.

**Do NOT use when:**

- Another checkout may still be capturing, claiming, building, or releasing against this queue — use `do-work run`, which preserves foreign and unattributed claims.
- The archive itself needs diagnosis or repair — use `do-work cleanup`.

## Input

Accept exactly the targeting tokens and selection flags documented by `actions/work.md` → **Input**. Validate that grammar before recovery and reject any unrecognized residue with `actions/work.md`'s usage error. After validation, preserve `$ARGUMENTS` verbatim for Step 1; this action does not reinterpret, add, or remove a work-pipeline argument.

Choosing this verb is the user's deliberate assertion that this checkout has sole queue authority for this invocation. It is not persisted and never changes plain `do-work run`.

## Steps

This verb is where the orchestrator's judgment does the most work. Every canonical command below is deterministic and refuses what it cannot prove; when one refuses, the run is stuck and the user has already asserted sole authority, so the next move is to reason about the cause and clear it, not to hand the refusal back. `actions/work-reference.md` → **Stuck Runs Hand Off to Judgment (any step)** is the rule.

### Step 0.1: Recover finalization with sole-releaser authority

From the project root, invoke the canonical launcher exactly once:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json recover-finalization --discover --assume-sole-releaser
```

Continue only on typed `success` with every ordered `finalizations` record carrying terminal phase `cleanup_complete` and empty `blocked_paths` and `reason_codes`. The singular `finalization` field is only the one-record compatibility projection. If the canonical launcher is missing, failed, or malformed, stop with its finding; there is no prose, manual, helper, or free-form fallback.

The assertion widens only the shared metadata classes owned by `--assume-sole-releaser`. It never widens recovery to secret-classified or project paths. Dirt the pipeline itself wrote earlier in the run — a working REQ or checkpoint whose bytes still hash to the journal's recorded preimage — is finalization input under every verb, so a remaining refusal names a path this run did not write. Clear it with judgment, not with another mechanical rule: the assertion has already answered the ownership question, so decide what the blocking bytes are, take the least destructive action that clears them, and re-run the exact command above (`actions/work-reference.md` → **Stuck Runs Hand Off to Judgment (any step)**). Under this verb, stop only for an action that would be destructive or irreversible.

### Step 0.2: Continue an interrupted archive, release, or commit tail

Treat each successful recovered finalization as a resumed `actions/work.md` Step 9, not as new work. When the interrupted tail left its implementation diff uncommitted or never wrote its release, the recovery command resumes that tail; its returned exact lifecycle and commit paths stand in for the canonical completion result wherever Step 9 consumes that result. Do not rerun `complete`, reconstruct staging, or infer provenance from the filesystem.

This is the only mid-step continuation this verb provides. A stopped build restarts from claim in Step 0.3; only the archive/release/commit tail resumes where it stopped.

### Step 0.3: Recover every working REQ under the authority assertion

Read `do-work/CHECKPOINT.md`, then classify every `REQ-*.md` in `do-work/working/` as this checkout's own crash, without a prompt or the three-hour takeover ladder. For each REQ, select exactly one checkpoint-evidence argument: `--checkpoint-writer '<exact label>'` when its entry carries a `writer:` label, `--checkpoint-unlabeled` when its entry has no label, or `--checkpoint-absent` when no entry names it. Report the observed label, or `no writer label`, then invoke the canonical ownership-transfer boundary from the project root:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json recover-claim REQ-NNN --request-path do-work/working/REQ-NNN-<slug>.md <checkpoint-evidence-argument> --assume-sole-writer --commit
```

The command consumes one snapshot, resets and moves the exact REQ, removes only the asserted checkpoint entry and its indented detail, and commits only its ownership-transfer postimages. Continue only on typed `success`. If the canonical launcher is missing, failed, or malformed, stop with its finding; there is no prose, manual, helper, or free-form fallback. Do not run Crash Recovery's mutation substeps separately.

Recovery returns the REQ to claimable state and strips incomplete generated sections, but its exact-path commit does not stage or discard an uncommitted implementation diff elsewhere in the project tree. `actions/work.md` pre-flight reports that diff, and the fresh builder may reuse it after judging it against the REQ. Mid-step build resumption remains out of scope.

### Step 1: Hand off to the work pipeline

Hand off with the validated original arguments unchanged:

```text
do-work run $ARGUMENTS
```

From here, follow `actions/work.md` unchanged. Do not reimplement selection, claiming, building, testing, review, or finalization in this action.

## Output Format

Report recovered finalization records and working-REQ takeovers first, then use `actions/work.md`'s ordinary progress and final hand-back formats.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
| --- | --- | --- |
| "The command refused and the instructions say no fallback, so the run ends here" | Judge why it refused, clear the cause, re-run the exact command | The no-fallback rule forbids reproducing a mutation by hand, not thinking; a Finder `.DS_Store` under `do-work/` once stopped a run whose only real obstacle was that one file |
| "The checkpoint says another writer, so I'll leave that REQ claimed" | Recover it and report the checkpoint's label | The user asserted sole queue authority by choosing `run-with-recovery`; preserving that claim would silently discard the verb's purpose |

## Verification Checklist

- [ ] Recovery used the canonical launcher with `--discover --assume-sole-releaser`, with no fallback
- [ ] Every working REQ was reset through canonical `recover-claim`, its exact checkpoint evidence was supplied, and its prior label was reported
- [ ] The work handoff received the original `$ARGUMENTS` verbatim
