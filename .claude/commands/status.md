Report the current status of the do-work task queue.

Check these locations and report what you find:

1. **Pending REQs**: List files in `do-work/queue/` (the pending queue)
2. **In-progress work**: List files in `do-work/working/` (currently being processed)
3. **Active pipeline**: Check `do-work/pipeline.json` for active pipeline state
4. **Recent completions**: List the 3 most recently archived URs in `do-work/archive/`
5. **Version**: Read the current version from `actions/version.md`

After reporting, suggest the most logical next action based on the current state.
