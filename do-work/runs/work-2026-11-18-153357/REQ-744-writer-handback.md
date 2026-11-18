# REQ-744 writer handback

Implemented and committed on `t3/implement-req-744-ai-report` at `5bdc38b2dd2c86e92658f73ffc73d86fb97881f8`, based on `82a0eefe05f5cefee388151560ec56aeb378d6de`. The worktree is clean. The commit contains the four implementation files and no queue or release metadata.

## Files changed

- [skills/do-work-toolbox/actions/ai-report.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-c608b8cc-2ee895741fc0d8aa/skills/do-work-toolbox/actions/ai-report.md:170) runs the installed gallery listing after Step 3 evidence and before Step 4 visual work, including optional generation. Ordinary creation, both brief kinds, and revise use the shared step. Step 8 and the checklist verify the choice and repeat each closing suggestion and command.
- [skills/do-work-toolbox/actions/ai-report-reference.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-c608b8cc-2ee895741fc0d8aa/skills/do-work-toolbox/actions/ai-report-reference.md:75) explains how to borrow moves, preserve the kit, record the sample name or `none`, and choose runner-up directions and verified presets. It preserves the content-specific fallback and teaching suggestion.
- [skills/do-work-toolbox/docs/ai-report-kit/skeleton.html](/Users/t2/.t3/worktrees/skill-do-work2/t3-c608b8cc-2ee895741fc0d8aa/skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:181) adds the head marker and a repeatable layout block under `#another-layout`. Its teaching block remains singular.
- [skills/do-work-toolbox/docs/ai-report-guide.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-c608b8cc-2ee895741fc0d8aa/skills/do-work-toolbox/docs/ai-report-guide.md:107) documents matching, choice metadata, and the reader's follow-up commands.

## Verification

- Passed `git diff --check` and the repository's `contract-regressions.sh`, including shipped-package references and shell-block lint with ShellCheck.
- Passed six focused Go tests in `internal/toolboxcommands` covering the read-only gallery listing, installed override precedence, second-opinion prepare, and preset resolution from the bundle and shipped files.
- Installed a local candidate archive in a separate scratch Git consumer. Ran the exact installed Step 3 listing, confirmed four TSV rows and no gallery or Git-status changes, and read the two selected Borrow lists. Compared the second-opinion dispatcher byte-for-byte with the base revision.
- Materialized reduced skeleton fixtures for one runner-up, two runner-ups, and the `none` fallback. Parsed the HTML and command arguments to verify metadata, command parity with a fixture summary, one direction argument per layout, exact bundle targets, and the teaching line.
- Sent the parsed runner-up command to the installed `ai-report-second-opinion --prepare` path. The typed row returned `prepared` for the fixture's `index.html` with shipped `first-principles` and the exact direction in the prompt. The source page stayed byte-identical and no sibling was published.

## Limitations

The reduced fixtures exercise the writer's instructions, template shape, and CLI handoff. Independent acceptance covers complete report generation for each route and the light/dark render checks. The presentation CSS header still describes one alternative layout in a comment outside this REQ's four-file scope.
