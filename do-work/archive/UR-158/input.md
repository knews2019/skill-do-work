---
id: UR-158
title: 'Second opinion from the session model must run as a visible in-session prompt, not a hidden shell'
created_at: 2026-10-19T14:43:26Z
requests: [REQ-700, REQ-701]
word_count: 56
---
## Summary

The user ran `/do-work-toolbox ai-report second-opinion <bundle>/index.html` inside a Claude Code session. The action's Session form made the agent launch one shell command (`do-work-cli.sh ... ai-report-second-opinion --backend claude --model claude-fable-5-1 ...`) that spawned a second `claude -p` process. The shell sat for over 10 minutes with "No output available" and the user could not see progress. The user wants the clarity-editor pass for the session's own model to run as a normal prompt inside the agent, or as a background subagent or workflow the agent starts, never as a hidden shell subprocess.

A second screenshot, sent while this capture was being written, shows what happened next in that session: the run finished (exit code 0 for the background shell) but the launcher's model probe found two models in the Claude reply (the editor had spawned an Opus-based Explore subagent), called the result ambiguous, and deleted the finished page after 12 minutes of work. The agent recovered the page from the transcript and reran with the Agent tool denied through a temporary local setting.

## Extracted Requests

| REQ | Request |
| --- | --- |
| REQ-700 | The session-model second opinion runs as a visible in-session prompt, subagent or workflow, not a hidden shell subprocess |
| REQ-701 | A failed post-run model check keeps the finished page and reports it; a subagent's model never makes the backend's model ambiguous |

## Assets

- `do-work/user-requests/UR-158/assets/REQ-700-screenshot-1-second-opinion-shell-run.png` — the session transcript and shell-details pane described in REQ-700 (second opinion from the session model runs in-session, not as a hidden shell).
- `do-work/user-requests/UR-158/assets/REQ-701-screenshot-2-second-opinion-probe-deleted-page.png` — the follow-up transcript described in REQ-701 (an ambiguous model probe deleted the finished page).

## Full Verbatim Input
> ```
> do-work capture-request
> [Image #1] <- let's update the second opinion so we have a feedback, this should have just triggered a normal prompt inside the ai-coder that way I can see the progress, alternatively it's fine to start a background subagent, or a workflow, but not a shell like this
> [Image #3] <- see additional feedback
> ```
