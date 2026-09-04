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

### Step 1: Recover under the user authority assertion

Invoke one constant command from the project root:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json recover --assume-sole-authority
```

Continue only on typed `success`, then read the returned `finalizations` one record at a time: a record carrying `FINALIZATION-SET-ASIDE` in its `reason_codes` is one REQ this run will not select — recovery keeps that REQ's existing claim instead of releasing it, so this verb's claim reset cannot hand the same REQ back to selection — and every other REQ still runs (`actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)**). Report each set-aside with its reason codes and the verb that resolves it before handing off. The command settles finalization first, then atomically resets every recoverable working claim it did not set aside and removes all same-request checkpoint entries from structural evidence. It reports hostile, multiple, and unlabelled writer records as data; no observed writer text enters shell source. Clear a refusal through judgment and re-run its exact `verification_argv`; never reproduce the mutation by hand.

### Step 2: Hand off to the work pipeline

Hand off the validated original arguments unchanged: `do-work run $ARGUMENTS`. Follow `actions/work.md`; do not reimplement selection, claiming, building, testing, review, or finalization here.

## Output Format

Report recovered finalization records and working-REQ takeovers first, then use `actions/work.md`'s ordinary progress and final hand-back formats.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
| --- | --- | --- |
| "The command refused and the instructions say no fallback, so the run ends here" | Judge why it refused, clear the cause, re-run the exact command | The no-fallback rule forbids reproducing a mutation by hand, not thinking; a Finder `.DS_Store` under `do-work/` once stopped a run whose only real obstacle was that one file |
| "The checkpoint says another writer, so I'll leave that REQ claimed" | Trust the `recover --assume-sole-authority` result | The user's authority assertion covers every structurally observed working claim for this invocation |
| "One record came back set aside, so this run is over" | Report that REQ with its reason codes, then hand `$ARGUMENTS` to the work pipeline unchanged | A set-aside is one REQ's exclusion, not the queue's; REQ-456's stuck finalization tail parked 31 pending REQs by being read as a whole-run stop |

## Verification Checklist

- [ ] Recovery used the constant canonical `recover --assume-sole-authority` argv, with no writer interpolation or fallback
- [ ] Every working REQ and every same-request checkpoint entry appears in the typed recovery result
- [ ] Every set-aside record was reported with its reason codes and a resolving verb, and the REQs it did not name still went to the work pipeline
- [ ] The work handoff received the original `$ARGUMENTS` verbatim
