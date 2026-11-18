# REQ-743 integration handback

Released **0.325.0** locally on main. Terminal status **completed-with-issues**, review **Partial, 83%**. The previous updater installs the new samples but cannot create the gallery during the first invocation. It also preserves a hand-made page silently. Both old Just recipes reproduced this. Running the same recipe again creates a missing page with the new binary. Explicit `ai-report gallery` also works.

[Open this checkout’s gallery](../../../docs/report-gallery/index.html) · [Archived request](../../../do-work/archive/REQ-743-report-gallery-samples-command-add-form.md) · [Finalization receipt](REQ-743-finalization.json)

| Revision | Commit |
| --- | --- |
| Integration base | `8ea2e7585d9e627f0c0d3150f00fdb2bdf0fa667` |
| CLI, samples, workflow tips | `683dcf1d`, `813c0517`, `a8031e04` |
| Three merges in order | `42b65d33`, `f73425fd`, `68353189` |
| Retry repair and generated-page cleanup | `f0d67478`, `7a3fdd37` |
| Final implementation and heavy execution | `764f2d45dc9fc51a63e0f10a5ac73abffb7669f0` |
| Canonical finalization | `52fa1482ec27df8ad34492745cfeac8d47df5c6a` |

Step 9 returned `cleanup_complete` with empty blocked paths and reason codes. The claim ran 38m 36s through canonical finalization (2026-11-18T15:33:40Z to 2026-11-18T16:12:16Z). From recorded handback observation to finalization was 20m 16s. Integration through post-release verification took 23m 02s from that observation.

- Qualification, scope, focused tests and canonical gate passed. Twelve named focused tests passed. The final implementation gate ran 428 board tests and 995 CLI tests in 124 seconds. The same gate passed again on the finalized release in 129 seconds. [Qualification](REQ-743-packaging-qualification.json) · [Testing](REQ-743-packaging-testing.json)
- All six heavy lanes executed and passed on the final target. Recorded durations were queue-kanban-javascript 6s, queue-kanban-browser 70s, do-work-cli-integrations 54s, staged-skills 30s, updater 67s, installer 27s. An earlier browser skip was cleared by naming Chrome. The initial packaging failure led to an archive contract correction. [Heavy results](REQ-743-packaging-heavy-run.json)
- Chrome opened the main gallery and all four local examples. Keyboard Copy matched the exact macOS pasteboard text. The installed consumer added the inflation comparison as a fifth sample and opened its saved image. Override, collision, malformed record, newest nested report, empty kind, asset links and hand-made refusal checks passed. [Browser check](REQ-743-packaging-browser.log)
- Fresh install and later gallery-aware writes generate the page. A hand-made page survives with one warning on those writes. The immediate historical transition remains open. For a local HTML source on another Windows drive, the gallery falls back to the sample record, which retains the source path. This follows the renderer’s cross-volume path handling and was not tested on a Windows desktop. [Released-consumer checks](REQ-743-released-consumers.log)
- The checkout gallery is excluded from the consumer archive. The generated page itself is committed. UR-184 and REQ-744’s queue record retain their original bytes. The live focused baseline and run scratch remain on disk. The source worktrees remain available.

A [copy of the accepted sample record](REQ-743-accepted-sample.md) accompanies this handback. The saved consumer sample for the writer acceptance case is at `.git/REQ-743-checks/consumer/docs/report-gallery/pace-and-level/sample.md` in this checkout (consumer commit `c71d5c864113c17bef669b667e9ba1be65841381`). Its source image remains in the UR-182 archive. The new updater lesson links the current flat archive path. When UR-183 closes, its archive-link sweep should follow the record to its new location.
