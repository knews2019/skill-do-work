# REQ-748 round-two handback

Five second-round variations build on the five exact reports the reader marked with “Like all five”. Each generation plan contains one saved whole-report like with an empty figure ID and no note. The five original report bundles, six source files, source manifest and tracked feedback file were verified against their saved snapshots.

**Source ID** `do-work/test-efficiency.tsv@a1b622fb2fb3b6fe4a7b5c0d8ed0d2e300f7e46a#standalone-20260906T134958Z-1358`

## Final reports

| Variation | Frozen parent | Generation | Independent factual review |
| --- | --- | --- | --- |
| [Child CPU seconds per elapsed second](../../../ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/index.html) | [Two clocks at the command line](../../../ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html) | [T3 receipt](../../../ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-generation.json) | [Pass](../../../ai-reports/2026-11-19_1559_req748-r2-cpu-elapsed-audit/gallery-factual-review.json) |
| [Finalization Across Individual Runs](../../../ai-reports/2026-11-19_1559_req748-r2-latency-ranking/index.html) | [The Finalization Gap](../../../ai-reports/2026-11-19_0154_req748-latency-ranking/index.html) | [T3 receipt](../../../ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-generation.json) | [Pass](../../../ai-reports/2026-11-19_1559_req748-r2-latency-ranking/gallery-factual-review.json) |
| [Where the median falls](../../../ai-reports/2026-11-19_1559_req748-r2-replicate-variability/index.html) | [The range behind the middle](../../../ai-reports/2026-11-19_0154_req748-replicate-variability/index.html) | [T3 receipt](../../../ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-generation.json) | [Pass](../../../ai-reports/2026-11-19_1559_req748-r2-replicate-variability/gallery-factual-review.json) |
| [Tool shares with a case ledger](../../../ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/index.html) | [Case shares by tool](../../../ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html) | [T3 receipt](../../../ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-generation.json) | [Pass](../../../ai-reports/2026-11-19_1559_req748-r2-tool-invocation-counts/gallery-factual-review.json) |
| [Signed gaps against cold medians](../../../ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/index.html) | [The warm-label split](../../../ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html) | [T3 receipt](../../../ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-generation.json) | [Pass](../../../ai-reports/2026-11-19_1559_req748-r2-warm-cold-labels/gallery-factual-review.json) |

## Authorship and reviewed bytes

| Variation | Astra author commit | Final bundle SHA-256 | T3 thread |
| --- | --- | --- | --- |
| Child CPU seconds per elapsed second | `b7b2ce295170a190cfd22fc098bd958fd6c5d9c2` | `ecc48ae213d0d393b678d33554172fb1c2be8cb336facce2b6065cb4347048b6` | [Author REQ-748 elapsed and CPU variation](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A2fca33b0-20d1-4f55-a05d-3a03094e1b5e%3Adelegate-task%3A25535d35-67f0-4af9-b63a-3c3ac032dd5e) |
| Finalization Across Individual Runs | `877e1551fd7c6a9fb618844cd49badcc2cde3706` | `bf99d4520993908150d0ee8e772f59589bc7e6f33ade22118864fdad0e8fed7a` | [Author REQ-748 latency ranking variation](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A2fca33b0-20d1-4f55-a05d-3a03094e1b5e%3Adelegate-task%3A3a5d83d0-d288-47e9-85ca-f4771655d3f8) |
| Where the median falls | `6f6e7a184410323fe36b626708bc8bfdbe1f2983` | `3f214aa92beef8aed44aff9eac687125220cdafa553e946f5fa8716ff1d9f312` | [Author REQ-748 replicate variability variation](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A2fca33b0-20d1-4f55-a05d-3a03094e1b5e%3Adelegate-task%3A89e15045-185c-4234-bbe0-88301466653c) |
| Tool shares with a case ledger | `8bd309abae26ee26f5b578657a298ec945d5071f` | `64fd78b6ad1dbab434edf076ce44348f12544b32e051f8eca8bc7e33d9677c37` | [Author REQ-748 tool counts variation](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A2fca33b0-20d1-4f55-a05d-3a03094e1b5e%3Adelegate-task%3Adb98dc77-b0f8-4f4d-a2ec-179ce6780c8c) |
| Signed gaps against cold medians | `ca0218123e0a04f65c8f8e58fad7d1560e1c9cfa` | `5f7a43f5c66762d95da2dc4c703de30bd1d77b4e4137c4b55ae510cd31e51bfb` | [Author REQ-748 warm and cold variation](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A2fca33b0-20d1-4f55-a05d-3a03094e1b5e%3Adelegate-task%3A438ada81-c0a8-4011-8d45-8ffed9487304) |

Each T3 run completed using GPT-6.2-astra in its own T3 worktree. No more than three authors ran concurrently. The integrator clarified the latency row order and parent design metadata, renamed the variability table heading to “Label”, removed trailing whitespace from the tool counts page, and gave the warm comparison a distinct page title and figure label. The touched bundles were rendered again and independently reviewed at their resulting digests.

## Acceptance evidence

- The source verifier matched the historical 27 TSV rows and nine groups. Independent recomputation checked timing medians, both rate denominators, the elapsed gap, intervals and ties, and intercepted invocation arithmetic.
- Twenty wide/phone and light/dark captures showed the five reports without horizontal overflow. The bundle-only judge flagged sibling links as 404. The live gallery opened all 12 linked before and after targets, establishing the `parent-links-only` render classification.
- The live gallery showed ten 1040×740 previews, eleven child figure controls and five exact copied starters that resolved to their named design. Ten report and figure save checks reloaded from a disposable store. The tracked feedback file still contains the five genuine likes.
- The selector excluded an unmarked parent and the preserved statistical-independence error in isolated replays.
- Canonical qualification and scope-drift checks passed for `4ffc9b99061e239f0f0c2ccb7bd6279b4d482fa7..c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb`. The integrated maintainer gate passed in 133 seconds with 427 board and 1,004 CLI tests. The focused test probe returned exit 0.

Local receipts are under `.git/do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-two/`. They include the five native T3 terminal records, 20 final captures, source recomputation, gallery browser results, exclusion tests, canonical gate records and finalization output.

## Canonical completion

- REQ-748 is [archived as completed](../../../do-work/archive/UR-184/REQ-748-first-two-round-report-trial.md). The canonical result contains one finalization, `cleanup_complete`, with empty blocked paths and reason codes.
- The supplied implementation commit is `c5a54ebbdbc8f48eeadaf6f6ef93cdafe3ce72eb`. Step 9 created archive commit `b8094babcb8932bce1794f7f5ddd6633f0897799`. The finalization journal was removed.
- The UR-184 closure gathered the three already-archived sibling records and the original input alongside REQ-748. Their bytes match the pre-finalization versions. The claims and artifacts of other requests were left alone.
- All five clean, merged author worktrees were removed. All commits are local. The scoped action stopped at Step 9.

[Step 9 receipt](../../../.git/do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-two/step9-finalization.json) · [Finalization manifest](../../../.git/do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-two/finalization-manifest.json) · [Integrated checks](../../../.git/do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-two/integrated-verification.json) · [Browser checks](../../../.git/do-work-run-scratch/work-2026-11-18-191126/REQ-748-round-two/gallery-browser-verification.json) · [Progress log](REQ-748-progress.log)

## Heavy verification

All six selected lanes executed at the integrated revision. The board browser lane initially skipped on PATH discovery and passed its rerun with the installed Chrome executable, running 39 browser tests.

| Lane | Seconds | Exit |
| --- | ---: | ---: |
| queue-kanban-javascript | 7 | 0 |
| queue-kanban-browser | 78 | 0 |
| do-work-cli-integrations | 68 | 0 |
| staged-skills | 35 | 0 |
| updater | 69 | 0 |
| installer | 31 | 0 |
