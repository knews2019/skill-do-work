# REQ-747 integration receipt

`do-work-toolbox ai-report gallery generate` creates one evidence-pinned design per T3 Code invocation. A follow-up names the exact parent page and snapshots its current report and figure marks. The gallery still copies a named starter for each design. A factual failure, stale review or missing review excludes a report from parent selection.

| Event | Revision or T3 execution |
| --- | --- |
| Clean preflight and builder base | `48c6a5532b1d93e1003fb322714fe37735380775` |
| Integration boundary | `d60c991e4001e6f31703ac9a2e3c3f6e33b5dd7d` |
| Isolated GPT-6.2-sol builder | `5e3f325fa49b6bba0f6c5b939834abcb683d0fa6`. [Implement REQ-747 report generation](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A735a479f-f725-4846-93e8-aa110cfd2dfc%3Adelegate-task%3A51962102-e894-4dbc-9bb3-643bff85c09e) |
| Merge and Step 9 implementation provenance | `4f9410eff23b3a42dde67cb18ee9ab8e75df274d` |
| Initial disposable generation | Completed GPT-6.2-astra on Codex. [Generate initial REQ-747 acceptance report](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A735a479f-f725-4846-93e8-aa110cfd2dfc%3Adelegate-task%3A5deb84e6-93ff-44e5-9f75-825c130e49ea) |
| Child disposable generation | Completed GPT-6.2-astra on Codex. [Generate marked REQ-747 acceptance child](t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A735a479f-f725-4846-93e8-aa110cfd2dfc%3Adelegate-task%3A99c1d1fe-0389-4545-8184-6ecbe02950f2) |

## Verification

- The request-bound preflight passed eight existing gallery checks and the direct maintainer gate. Merged scope and qualification passed. The two filename warnings resolve through the native dispatcher and Go test discovery.
- The builder's test-first `TestAIReportGalleryGenerationLineage` failed at the original `--prepare-generation` boundary with `TOOLBOX-USAGE` and passed after implementation. In a fresh install of the preflight tree the command also returned `TOOLBOX-USAGE`.
- On the merge, the maintainer gate passed in 141 seconds with 428 board and 1,003 CLI tests. Canonical `advance` recorded that exact SHA as green and the focused baseline comparison was satisfied.
- All four selected heavy lanes executed at the merge SHA and exited 0. CLI integrations took 70 seconds, staged skills 35 seconds, updater 71 seconds and installer 32 seconds.
- An archive of the merge installed all four suite modules into a fresh consumer. The installed action used actual T3 author delegation twice. The raw terminal task results provide the provider instance, actual executed model, task, thread and terminal run identity. The action persisted those fields into the two generation records.
- The trial source verifier matched 27 historical rows at `a1b622fb2fb3b6fe4a7b5c0d8ed0d2e300f7e46a`. The authors used identical pins for the six file paths in [the September 6 manifest](trial-source/gallery-input.json). The initial plan recorded an empty parent and no marks.
- Chrome saved a separate report like and dislike on **Median elapsed wall time (seconds)**, with a clearly synthetic note asking for labelled median and range rows. The child showed those rows under `#finalization-comparison`, kept the Git count panel, and linked the actual before and after figures. The stored mark snapshots matched both the version-1 gallery API and saved JSON. Chrome copied the two starters and the resolver reached each exact `gallery-design.md`, including by its qualified identity.
- Four desktop and phone captures in light and dark supported the factual review of each valid candidate. The child judge recorded one expected sibling parent 404 from its bundle-only origin. Chrome verified the actual before and after URLs and visible figures on the live gallery origin, so the child factual review passed with `parent-links-only`.
- A separate disposable copy changed only the visible cold label from `48.606 s` to `48.606 ms`, keeping its report like. Parent selection first rejected the stale review. The layout judge still passed, but independent rendered review recorded the unit error on `#finalization-wall`. The same child invocation then exited at selection before delegation with the specific factual finding. Chrome could still copy the disqualified design and the named resolver found its sidecar.

## Local acceptance artifacts

Root `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/do-work-run-scratch/work-2026-11-18-191126/REQ-747-acceptance.nrAtmY`

| What | Path under the root |
| --- | --- |
| Native execution records | `builder-t3.json`, `initial-t3.json`, `child-t3.json` |
| Builder audit | `builder-handback.md`, `builder-progress.log` |
| Positive report bundles | `consumer/ai-reports/2026-11-19_0140_req747-acceptance-wall-and-counts/`, `consumer/ai-reports/2026-11-19_0141_req747-acceptance-labelled-observations/` |
| Rejected report bundle | `negative/ai-reports/2026-11-19_0140_req747-acceptance-wall-and-counts/` |
| Browser proof and exclusion | `initial-gallery-marks-and-copy.json`, `child-starter.json`, `copied-starter-resolutions.json`, `live-lineage.json`, `negative-selection.json` |
| Render and fact records | `initial-render/`, `child-render/`, `negative-render/`, and each candidate's `gallery-factual-review.json` |
| Canonical evidence | `preflight.json`, `qualification.json`, `test-gate.json`, `merged-maintainer.log`, `heavy-plan.json`, `heavy-result.json` |

## Trial invocation and limits

Paste this command into the T3 Code agent in the main project:

```text
do-work-toolbox ai-report gallery generate --evidence do-work/runs/work-2026-11-18-191126/trial-source/gallery-input.json --design two-measures-side-by-side
```

This invocation is for the first candidate. Follow-up invocations take one exact `--parent <report/index.html>` and the same manifest. For the user trial, start the second round only after the user saves genuine first-round marks. The version-1 feedback file records targets, marks and notes, so the caller distinguishes genuine reader feedback from fixture input. The caller reviews the rendered facts against the pinned evidence, and the native selector requires that review to match the current bundle bytes. Keep trial source and parent paths stable between rounds.

The work remained local to main and disposable acceptance projects. REQ-748 and first-round trial generation remain with the coordinator. No genuine permission approval was pending at any point.
