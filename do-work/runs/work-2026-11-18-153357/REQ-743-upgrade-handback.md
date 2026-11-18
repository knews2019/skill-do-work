# REQ-743 first-upgrade handback

Both Just recipes installed by 0.324.3 apply the gallery-capable suite successfully and leave `docs/report-gallery/index.html` absent. The first-upgrade acceptance criterion remains open for those callers. I found no supported way for that invocation to generate the page and made no code change.

## What runs during the upgrade

- Both historical recipes start the installed `.claude/skills/do-work/tools/do-work-cli.sh` with `update-suite`. The launcher selects the binary before the archive is fetched. The Go path execs the binary built from the installed sources in the Go cache. The prebuilt path selects the binary for the installed version and execs its per-version cache entry. The 0.324.3 Go process still runs after the new files are written. [Just recipes](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/skills/do-work-board/justfile.template:58) · [Launcher](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/skills/do-work/tools/do-work-cli.sh:48)
- The old process validates the incoming archive, checks its version, calls its own `RunInstall` in-process, verifies the result and returns. Its post-update report check only computes the existing read-only organize hint. The workflow slice `a8031e04` adds the gallery call to the new process's handlers after a confirmed install or a successful update that wrote files. [Old update transaction](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go:92) · [Old update handler](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go:165) · [Workflow handback](/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-18-153357/REQ-743-workflow-handback.md)
- The old manifest accepts exactly four fixed `.claude/skills/` destinations. It copies samples into the installed toolbox, while the requested project page lies outside those destinations. The other writes publish the managed Justfile, Claude settings and instructions. Just is called with `--list` for validation, and the settings hooks are composed as JSON. `--upstream-url` and `DO_WORK_UPSTREAM_URL` select archive data. Even supplying the candidate through the old installed installer's `--archive` uses the historical CLI. [Manifest](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go:16) · [Managed writes](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go:875)

## Verified rollout behavior

| Invocation | Observed result |
| --- | --- |
| Clean install through the candidate's extracted bootstrap installer | Installed the suite and generated the page with all four shipped samples. |
| Either old Just recipe, or the old installed `do-work-update.sh` | Exited 0, installed fixture v0.324.4 and its toolbox samples, and left the project page absent. |
| Either Just recipe again against the same archive | Exited 0 with `You're up to date (v0.324.4)`. The new CLI ran, the version check skipped the write, and the page stayed absent. |
| Explicit generation after the old update | The new CLI generated the page with all four samples. `--list` returned four rows and wrote nothing. |
| A later writing update from the gallery-aware install | Generated the page from the newly installed sample text. A hand-made page survived byte-for-byte and produced exactly one gallery warning while the update succeeded. The first historical update preserved the hand-made page silently. |
| Explicitly run the candidate bootstrap over 0.324.3 | The documented bootstrap flow reconciled the existing installation with one confirmation and generated the page. The candidate installer also generated it on an explicitly confirmed reinstall at the same version. [README bootstrap](/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c/README.md:9) |

After an old Just update, use the installed `ai-report gallery` action or run this command from the consumer Git root:

```bash
bash .claude/skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text ai-report-gallery
```

The [consumer probe](/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-18-153357/REQ-743-consumer-probe.sh) ran unchanged for fresh install and both old Just routes on Go 1.27.4 and Just 1.46.0. Fresh passed. Each old Just command succeeded and its probe failed at the missing-page assertion. Disposable follow-up checks used base `fe0b06b7` and a combined tree of workflow `a8031e04`, CLI `683dcf1d`, and samples `813c0517`. Version numbers v0.324.4 and v0.324.5 were fixture metadata only. The prebuilt path was traced in source.

The unchanged probe's hand-made branch also counts path mentions in the review diff. Its fresh hand-made run preserved the page and emitted one actual warning, yet the probe exited 1. The follow-up counted only stderr lines beginning `Report gallery: could not regenerate docs/report-gallery/index.html`.
