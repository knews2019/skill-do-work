---
id: UR-126
title: 'Cap the path list in a verify finding'
created_at: 2026-09-05T18:29:25Z
requests: [REQ-590]
word_count: 93
---

# Cap the Path List in a Verify Finding

## Summary

The user opened the board and found the VERIFY FINDINGS strip filled with several hundred comma-separated paths in one finding, pushing the whole board off the screen. The investigation that followed found the cause in the producer, not in the board: `queue-kanban verify` joins every path a probe collected into one detail sentence with no cap, and the board prints the producer's detail verbatim by design. The same unbounded join sits at three call sites, and the terminal report prints the same sentence on one line. The user then asked whether anything else needed fixing, was given the three sites plus the shared-snapshot consequence, and said to capture it as a REQ.

## Full Verbatim Input

> ```
> [Screenshot: the board at 127.0.0.1:8090, Board view. The VERIFY FINDINGS strip says "5 findings". Four rows read normally. The fifth row, chip WORKTREE-WROTE-QUEUE-STATE under subject worktree-agent-REQ-589-m4-slim-band, is a wall of text that fills the whole page: "has uncommitted changes under do-work/ in .git/work-run-20260905-1820/worktree-agent-REQ-589-m4-slim-band (do-work/.req-reservations/REQ-000391, ...)" followed by several hundred comma-separated paths — every reservation marker and every archived REQ file. The board columns are pushed far below the fold.]
> 
> <- why is the kanban board broken with such a long text?
> 
> anything else that needs to be fixed?
> 
> capture it as a req
> ```
