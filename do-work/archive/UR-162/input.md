---
id: UR-162
title: 'recover should not treat untracked do-work notes as unfinished finalization work'
created_at: 2026-11-02T19:41:34Z
requests: [REQ-709]
word_count: 62
---
## Full Verbatim Input
> ```
> do-work capture request: recover should not treat untracked do-work/inbox notes as unfinished finalization work, so do-work commit can commit them
> 
> Answers chosen during capture:
> - Which loose files under do-work/ should the recovery check stop treating as unfinished work? -> Any file no record can own (Recommended)
> - Is this the right test for the fix? -> Yes, use this test (Recommended)
> ```
