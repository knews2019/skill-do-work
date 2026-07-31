---
id: UR-004
title: just run-kanban should replace a stale server and open the browser
created_at: 2026-07-01T21:55:23Z
requests: [REQ-017]
word_count: 55
---

# just run-kanban should replace a stale server and open the browser

## Full Verbatim Input

just run-kanban should not only start the server (kill the previous one if it is on the same port) but also use the open command to open the default browser with the target kanban, ohterwise I run the command and I don't see the dashboard

## Conversation Context

Sent immediately after the user ran `just run-kanban` in a consumer project (`~/Desktop/n1-weekly-signal-diff`) and saw the server start (`queue-kanban: live board at http://localhost:8090`) with no browser opening. The preceding assistant turn had explained that neither the `run-kanban` just recipe (`actions/install.md`) nor `queue-kanban serve` (`tools/queue-kanban/serve.go`) contains any browser-open logic, and had recommended an `--open` flag on `serve` (platform-aware, off by default) wired into the recipe as the better fix over a shell `sleep && open` hack. A screenshot of the terminal output was shown but adds nothing beyond the transcript above (no error visible — the server started fine; the complaint is the missing browser open and the port-collision behavior on re-runs).

---
*Captured: 2026-07-01T21:55:23Z*
