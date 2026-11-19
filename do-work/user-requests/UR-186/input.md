---
id: UR-186
title: 'Show an actual preview for every gallery report'
created_at: 2026-11-19T12:02:21Z
requests: [REQ-753]
word_count: 187
---
## Full Verbatim Input
> ```
> Act as the implementation sub-agent for this task.
> 
> You are sole main-checkout integrator for the user correction that every gallery report needs an actual preview. Read CLAUDE.md, _dev/primes/prime-releases.md and .git/gallery-preview-fix-handback.md. Integrate completed t3/gallery-report-previews commit 7ca9382a into main (main also has reader preference commit 0cf8821d to preserve). Use existing prepared implementation and evidence, capture a small dedicated correction REQ via do-work if necessary for canonical release lifecycle, count work already done instead of rebuilding or rerunning earlier baselines. Complete required focused/integration verification and release metadata canonically using appropriate 0.331.1 patch. Review scope & real defect regressions, e.g. no report mutation, all 55 previews, links & existing controls, installed path, browser screenshot error handling. Before finalization check /tmp/gallery-preview-before-hashes.json source files all unchanged unless explicit catalog derived changes permitted. No user marks, no trial round-two generation, leave REQ-748 waiting. No nested workers. Do not start/restart the user's viewer (coordinator owns session port50932 and will restart after receipt). Ordinary direct commands, no approval-triggering compound wrappers; honor permissions. No push. Return clean main/head, release & integration hashes, exact test evidence and actual limitations. Preserve source worktree until fully verified, no force cleanup.
> ```
