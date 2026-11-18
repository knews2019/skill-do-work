# REQ-743 CLI handback

- Branch `t3/implement-cli-report-gallery`
- Commit `683dcf1da005ae0af7167ea424160817be4b81c1`
- The worktree is clean.

## CLI contract

- `do-work-cli --repo-root ROOT --format text ai-report-gallery --list` prints sorted tab-separated `name`, `use_when` and repository-relative `sample.md` fields. It writes no files. Warnings identify skipped records on stderr. A valid project override replaces its shipped counterpart. A broken project override leaves the valid shipped record available.
- `do-work-cli --repo-root ROOT --format text ai-report-gallery` prints `docs/report-gallery/index.html` after generation. It reuses the catalog walk for the latest report of each kind, links the actual kit and preset files, and provides a copyable agent command. The gallery groups the catalog's legacy `unknown` kind under Completed work and shows explicit custom kinds when present.
- Installed toolbox files in `.claude/skills/do-work-toolbox/docs/` take precedence over source toolbox files in `skills/do-work-toolbox/docs/` when both exist. The generator refuses an existing page without its generator header and refuses linked output paths. A checked replacement revalidates the same preimage used for the ownership decision.

## Files

- Command registration and assertion. [commands.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539/skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go), [commands_test.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539/skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go)
- Source and installed discovery, sample validation, list and publish. [report_gallery.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539/skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go)
- Gallery projection, escaped links, HTML and Copy control. [report_gallery_page.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539/skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go)
- Shared rooted publish preimage check. [mutation.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539/skills/do-work/tools/do-work-cli/internal/toolboxcommands/mutation.go)
- Fixtures for discovery, catalog selection, path handling and overwrite refusal. [report_gallery_test.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539/skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go)

## RED and GREEN

- RED before implementation. `go test ./internal/toolboxcommands -run '^TestAIReportGallery' -count=1` exited 1. Handler fixtures reported `ai-report-gallery is not registered` and the runtime returned `UNKNOWN-COMMAND`.
- A parser fixture also failed on an unterminated quoted `use_when`. A page fixture failed on a missing link to a same-volume absolute HTML source. Both passed after the respective fixes.
- GREEN. `go test ./internal/toolboxcommands -run '^(TestAIReportGallery|TestHandlersRegisterCanonicalToolboxCommands)' -count=1 -v` passed. The full test run had 33 passing packages and one package with no tests. `go vet ./...`, the focused Go 1.24.0 gallery and index tests, gallery tests with the race detector, and a Windows toolbox test cross-compile passed.
- Built CLI smoke. A temporary installed project used the real kit and preset directories with four fixture samples. Listing returned four rows, then five after adding a project sample. A broken record was named on stderr. The updater-form invocation generated the page and all 23 local links resolved. A Node harness exercised clipboard success, file-origin fallback and manual selection when copying was denied.

## Open verification

- First install and update with the actual shipped samples and workflow hook remain integration verification. The copy control was exercised with simulated clipboard behavior and awaits a real local-browser check.
- Absolute local source links are relative to the generated page on the same filesystem volume. Same-machine paths on another Windows drive fall back to the linked sample record.
