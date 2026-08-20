# Suggest Next Steps

**Intent:** After any do-work action completes, infer 2-3 next actions from what the action actually did and the current queue state, then suggest them as copy-pasteable `do-work <verb>` commands. Judge from the real outcome — an empty queue doesn't get a `run` suggestion, a clean tree doesn't get a `commit` suggestion, a domain of `ui-design` earns a `ui-review` suggestion where another domain wouldn't — rather than pattern-matching a fixed action-to-suggestion table.

## Format Constraints

- Render as a fenced block starting with `Next steps:`, one command per line, command left-aligned, a short gloss on the right, columns aligned.
- Suggest core commands from `SKILL.md`'s Routing table and extension commands from the owning sibling router. A capture suggestion is written `do-work capture-request: <text>`; extension commands use `do-work-board`, `do-work-knowledge`, or `do-work-toolbox` directly.
- Cap at 2-3 suggestions, ranked by relevance to what just happened and what's outstanding.
- Always close with: `do-work help` — full command reference.

**Example** (anchors the format only — not a template to reuse verbatim):

```
Next steps:
  do-work run                 Start processing the queue
  do-work clarify             Answer pending questions
  do-work roadmap             Survey what's left in the queue
```

## Non-Obvious Cases

Most actions have a next step inferable from what just ran (after `code-review`, suggest `run` to process findings; after `commit`, suggest `inspect` or `review-work`). The rows below are the ones where the *same* action leaves the queue in genuinely different states, and the right suggestion depends on which — getting it wrong here isn't a missed nicety, it actively misleads (e.g. suggesting `commit` before anything ran):

| After... | State | Suggest |
| --- | --- | --- |
| `capture-requests` | New REQs captured | `do-work verify-requests` before `do-work run` — check capture quality before building |
| `capture-requests` | Every new REQ is `effort_estimate: effort-mechanical` | `do-work run-simple-reqs` — the batch is cheap enough to hand a smaller model, and it lists and estimates before running |
| `clarify` | Questions still pending | `do-work clarify` again, not `do-work run` — unanswered REQs won't be picked up |
| `clarify` | All answered | `do-work run` to process them |
