---
id: UR-192
title: 'Let capture-files create the UR asset directory'
created_at: 2026-11-26T11:18:47Z
requests: [REQ-762]
word_count: 13
---
## Referenced review finding

**[P2] Let `capture-files` create the UR asset directory — [capture.md:204](/Users/t2/Desktop/sa2-sentence-aligner2/.claude/skills/do-work/actions/capture.md:204)**

The complex-capture instructions now create `UR-NNN/assets/` before the transaction runs. If the manifest is refused, that directory survives rollback. A subsequent `verify-requests` selects the highest-numbered UR and then fails to find its `input.md`. Stage the asset payloads privately and let the transaction create their destination directories.

## Full Verbatim Input
> ```
> Capture F2 from this review with its original P2 claim and verification evidence.
> ```
