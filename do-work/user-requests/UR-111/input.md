---
id: UR-111
title: 'Run held heavy lanes at queue exhaustion without asking, and record per-lane wall time'
created_at: 2026-09-04T13:19:11Z
requests: [REQ-566]
word_count: 103
---

# Run Held Heavy Lanes At Queue Exhaustion Without Asking, And Record Per-Lane Wall Time

## Full Verbatim Input

> ```
> "stopping on pending-heavy-testing because that needs your permission." <- I asked for this, but it turns out it's not very good, because the work just stops, what I need is an inteligent way to run it (either faster by 80/20 rule), or groupped with multiple tasks (but that can get complicated). Any suggestions?
> 
> I would start by monitoring the duration and lifting the human blocking, after all the goal is to do the queue faster.
> 
> [Answer to "When should the heavy-gate change (record per-lane durations, then run the held batch at queue exhaustion without asking) be captured as a REQ?"] Capture it now
> ```

---
*Captured: 2026-09-04T13:19:11Z*
