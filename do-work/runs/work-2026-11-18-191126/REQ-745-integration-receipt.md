# REQ-745 integration receipt

The merged gallery restores a report like and a separate figure dislike with its note. The browser, feedback API and saved project JSON agreed after restarting both Just launches. This receipt covers REQ-745.

## Provenance

| Event | Revision |
| --- | --- |
| Builder base | `f3840b087913d8501267dda1c1395580304f2ba4` |
| Integration boundary on main | `e6c2229eba1605dd4cf35212d51edb10d651c288` |
| Builder head, `t3/persistent-gallery-annotations` | `18e88b053872c352fefb9ecba5555703abefa64c` |
| Builder merge | `e90015a27bf09d72705d9355c9980ab395f17c6d` |
| Installer collision test fix | `f780672d` |
| Installer collision corrective merge | `aa2203c95402ecd055d5a4c5283422791c0f958f` |
| Bootstrap recipe verification fix | `b2407772` |
| Bootstrap corrective merge and supplied implementation SHA | `8567dede1b2c30ed5b6867088700159cbedaf17a` |

Qualification and review used `e6c2229e..8567dede`. The builder queue guards and the cumulative merge diff showed no paths under `do-work/`. Each merge had an empty pre-merge index. Preparation, the trial-source commits and REQ-746’s untracked handback and progress sit outside the implementation range.

## Effective implementation manifest

| Change | Path |
| --- | --- |
| Modified | `_dev/tests/install-suite-behavior.sh` |
| Modified | `justfile` |
| Modified | `skills/do-work-board/justfile.template` |
| Modified | `skills/do-work-toolbox/docs/ai-report-guide.md` |
| Modified | `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` |
| Modified | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` |
| Added | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` |
| Added | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback_test.go` |
| Modified | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` |
| Modified | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` |

## Verification and dispositions

- The prepared Scope label was missing a colon, and a backticked recipe name parsed as an extra file. Fixed both in the active REQ. Typed qualification and scope then passed. The two filename-based warnings were checked against the dispatcher call and Go test discovery.
- The first full gate at `e90015a2` exited 1 in 162 seconds. `TestEveryShippedManagedDefinitionRejectsAnExteriorCollision` expected 41 recipes and found 42. The isolated failure reproduced. The corrective test update runs the existing collision assertion for all 42 definitions. The next full maintainer gate at `aa2203c9` passed in 147 seconds. The installer heavy lane at that revision exited 1 after 33 seconds on a second old count. The other five heavy lanes passed. The final correction keeps the existing byte-for-byte template assertion and verifies that the fresh Justfile contains the named gallery recipe. The Go collision test also checks that the parser includes the named recipe. The preceding timing wrapper rejection used an invalid category and launched no gate.
- The request-bound focused baseline passed all eight checks. Ten focused gallery tests passed. The builder’s test-first `TestAIReportGalleryServeRoundTrip` RED was `TOOLBOX-USAGE`; integrated GREEN covers independent updates, restart, report-scoped figure identity, and rejected writes.
- At `8567dede`, the full maintainer gate exited 0 in 158 seconds. The request-bound `advance` recorded that exact revision as green and the focused test returned `satisfied`.
- The final heavy verification ran all six lanes at the same revision under the serial lock. Each exited 0. Durations were `queue-kanban-javascript` 7s, `queue-kanban-browser` 86s, `do-work-cli-integrations` 74s, `staged-skills` 36s, `updater` 75s, `installer` 38s.
- A disposable source fixture and fresh suite install ran the actual Just recipes. The source preview opened the real phase-by-phase report with its procedure diagram and proof. The installed consumer discovered the same bundle under `ai-reports/2026-11_1-report/` without a preview manifest. Both reopened tabs restored the same `how-diagram-title` figure alongside the whole-report mark. The recorded API responses matched the independently parsed files. An installed lock conflict kept the draft visible and the saved file byte-identical.
- Static regeneration was byte-identical. The existing browser gate passed all 28 checks in its own checkout across desktop and mobile, light and dark. The earlier scratch-path attempt hit G2 because the runner resolves trusted capture URLs relative to its own checkout.
- The first fixture installation returned a cancellation after receiving EOF at the confirmation prompt. The confirmed rerun installed the suite into the scratch repository. The integration fixtures recorded the system browser launch URL and the gallery itself was driven in T3 browser tabs. The builder also verified OS browser launch in both disposable projects. All marks were fixtures. The live inventory refreshes on restart. Unauthored figure IDs can depend on figure content or labeling. The integrator performed the review in-session.

The suite publishes this additive feature as 0.329.0 through Step 9, using the root version and changelog plus the three installed mirrors. The board executable carries its own version.

## Shell approval reports

The earlier compound shell submissions wrapped `run-timed-command … -- bash _dev/tests/maintainer-verify.sh` and `run-heavy-verification` with log redirection, lock acquisition and exit-trap cleanup. The first submission used an invalid timing category, the following two launched the initial maintainer gate and its retry, and the final compound submission sequenced the last maintainer gate and heavy rerun under the lock. The user reported repeated dialogs during this work. The shell results expose command outcomes but do not identify which submitted blocks displayed a dialog. After that report, the already-started verification ran to completion and the remaining lifecycle steps used individual documented commands.
