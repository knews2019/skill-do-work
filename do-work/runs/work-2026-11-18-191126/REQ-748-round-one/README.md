# REQ-748 round-one gallery handoff

The five reviewed reports from `160a96b5` are ready for reader feedback. They share the same six [evidence files](../trial-source/gallery-input.json). The [verification manifest](manifest.json) records the matching bundle digests, passing factual and render verdicts, T3 execution records, figure IDs, and exact copied starters. The original warm-label revision failed factual review because the evidence could not establish statistical independence. The manifest records its failed digest and how to reconstruct that wording. Its corrected revision is the fifth report below.

## Mark the five reports

Open the live gallery provided by the coordinator and find these cards under **Generated report designs**. Use **Open report** to read each page. On each card, expand **Review report and figures**, choose **Like** or **Dislike** on the whole report or on at least one of its figures, add an optional note, and click **Save feedback**. Wait for **Saved like.** or **Saved dislike.** before moving to the next card. Let the coordinator know when all five cards have a saved mark.

| Gallery card | Figure feedback targets |
| --- | --- |
| [Two clocks at the command line](../../../../ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html) | [Elapsed by selection](../../../../ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html#elapsed-by-selection), [CPU by selection](../../../../ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html#cpu-by-selection) |
| [The Finalization Gap](../../../../ai-reports/2026-11-19_0154_req748-latency-ranking/index.html) | [Elapsed ranking](../../../../ai-reports/2026-11-19_0154_req748-latency-ranking/index.html#elapsed-ranking) |
| [The range behind the middle](../../../../ai-reports/2026-11-19_0154_req748-replicate-variability/index.html) | [Finalization ranges](../../../../ai-reports/2026-11-19_0154_req748-replicate-variability/index.html#finalization-ranges), [Heavy-build ties](../../../../ai-reports/2026-11-19_0154_req748-replicate-variability/index.html#heavy-build-ties) |
| [Case shares by tool](../../../../ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html) | [Tool mix](../../../../ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html#tool-mix) |
| [The warm-label split](../../../../ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html) | [Cold and warm elapsed medians](../../../../ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html#cold-warm-wall-medians) |

**Use this design** opens the starter. **Copy starter prompt** copies the command for that exact card. The report links in this table open the pages for reading. Save marks from the live gallery cards.

## Resume after reader feedback

1. Confirm the user saved the marks. Read `do-work/report-gallery-feedback.json` and check for at least one current, genuine report or figure entry per exact `report_path` in the manifest. Use the same six evidence pins and run the following read-only command once for each of the five manifest report paths. Each one must return a nonempty `marks` array for its parent.

   ```text
   skills/do-work/tools/do-work-cli.sh --format text ai-report-gallery --prepare-generation do-work/runs/work-2026-11-18-191126/trial-source/gallery-input.json <report_path>
   ```

2. After all five qualify, clear the manual hold and reclaim this request. These are coordinator actions after feedback.

   ```text
   skills/do-work/tools/do-work-cli.sh --format json unblock REQ-748 --source user-via-clarify --confirmed --commit
   skills/do-work/tools/do-work-cli.sh --format json advance REQ-748 --dispatch-bound 1 --frozen-member REQ-748 do-work/queue/REQ-748-first-two-round-report-trial.md explicit-req false
   skills/do-work/tools/do-work-cli.sh --format json advance REQ-748
   ```

3. Restore the Route B triage, plan, exploration, and scope evidence as applicable from `160a96b5` at `do-work/working/REQ-748-first-two-round-report-trial.md`. The hold transaction resets those generated sections. Continue using the live claim and the verified round-one receipt. In the next wave, generate five Astra children from the five reviewed parents using the same source pins. Record saved marks, before-and-after figure links, and fresh factual and render verdicts for each child before completing the trial.

## Queue exit receipt

- Before the hold, the scoped consumed-member continuation and the overall `advance --fan-out 3` returned zero claims and empty continuations. The scoped selector returned `TARGET-NOT-FOUND` for the queued target because REQ-748 was in `working/`. Per-request `advance` independently reported the preserved Route B claim.
- The live feedback GET returned an empty list, and all five read-only parent preparations reported no current feedback. The authoritative feedback file was absent at the transition. The initial recovery call refused the uncommitted status edit as ambiguous. After committing that exact edit, `recover --take-over REQ-748` moved the request to the [blocked queue entry](../../../queue/REQ-748-first-two-round-report-trial.md) and removed its checkpoint claim in `470940ff`. A later `recover` reported no working claims and no pending finalization.
- After the move, scoped and overall `next` each reported zero pending and one blocked. Both `advance` modes and one replay of each returned continuation exited 1 with `BLOCKED`, zero claims, zero writes, and no automated probe. The scoped continuation appears above. The overall continuation was `skills/do-work/tools/do-work-cli.sh --format json advance --dispatch-bound 3 --frozen-member REQ-748 do-work/queue/REQ-748-first-two-round-report-trial.md default false`.
- The selector's generic remedy is `do-work clarify`. The trial unblocks only after the user's five saved marks are verified, then the explicit unblock and scoped continuation above apply. The session-end queue summary in `CHECKPOINT.md` belongs to the coordinator's run.
