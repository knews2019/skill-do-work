# Suggest Next Steps

After every action completes, suggest the next logical prompts the user might want to run. Use fully qualified action names so the user can copy-paste directly.

**After pipeline (completed — queue fully drained):**
```
Next steps:
  do work present work        Generate client-facing deliverables
  do work capture request: [describe]  Capture new requests
```

**After pipeline (interrupted — active pipeline still exists):**
```
Next steps:
  do work pipeline            Resume the active pipeline
  do work pipeline status     Check pipeline progress
```

**After capture requests:**
```
Next steps:
  do work verify requests     Check capture quality before building
  do work run                 Start processing the queue
```

**After work (queue processing):**
```
Next steps:
  do work review work         Review the completed work
  do work present work        Generate client-facing deliverables
  do work clarify             Answer any pending questions
```

**After verify requests:**
```
Next steps:
  do work run                 Start processing the queue
  do work capture request: [describe changes]  Capture additional requests
```

**After review work:**
```
Next steps:
  do work present work        Generate client-facing deliverables
  do work ui-review [scope]   Validate UI quality (if domain: ui-design)
  do work run                 Process follow-up REQs (if any were created)
```

**After code-review:**
```
Next steps:
  do work run                   Process follow-up REQs (if any were created)
  do work quick-wins [dir]      Scan for additional refactoring opportunities
  do work capture request: [describe fix]  Capture a finding as a request
```

**After ui-review:**
```
Next steps:
  do work capture request: [describe fix]  Capture findings as requests
  do work run                   Process follow-up REQs (if any were created)
  do work install-bowser        Install Playwright CLI + Bowser skill for visual verification (if not installed)
```

**After present work:**
```
Next steps:
  do work present all         Generate portfolio summary (if multiple URs completed)
  do work capture request: [describe]  Capture new requests
```

**After forensics:**
```
Next steps:
  do work cleanup               Fix orphaned URs and misplaced files
  do work run                   Process stuck or pending REQs
  do work capture request: [describe fix]  Capture a specific finding as a request
```

**After prime create:**
```
Next steps:
  do work code-review prime-{name}   Review the code scope the prime covers
  do work prime audit                Run a full audit to check the new prime
  do work run                        Process the queue
```

**After prime audit:**
```
Next steps:
  do work prime create <path>         Create primes for flagged utilities
  do work capture request: [fix]      Capture audit findings as requests
  do work run                         Process the queue
```

**After quick-wins:**
```
Next steps:
  do work capture request: [describe fix]  Capture a finding as a request
  do work code-review [scope]   Full code review for the same scope
  do work run                   Process the queue
```

**After scan-ideas:**
```
Next steps:
  do work capture request: [paste an idea]  Capture an idea as a request
  do work scan-ideas [different focus]      Brainstorm a different area
  do work deep-explore [concept]            Explore an idea in depth
  do work quick-wins [dir]                  Scan for quick refactoring wins
```

**After deep-explore:**
```
Next steps:
  do work capture request: [paste a direction]  Capture a direction as a request
  do work deep-explore continue [session]       Resume or extend the session
  do work scan-ideas [focus]                    Quick idea scan for a related area
```

**After inspect:**
```
Next steps:
  do work commit              Commit the ready changes
  do work capture request: [describe fix]  Capture issues as requests
  do work run                 Process the queue (if fixes were captured)
```

**After commit:**
```
Next steps:
  do work inspect             Review remaining uncommitted changes (if any)
  do work review work         Review the committed changes
  do work capture request: [describe]  Capture new requests
```

**After clarify questions:**
```
Next steps:
  do work run                 Process answered questions
  do work clarify             Continue answering (if skipped any)
```

**After bkb init:**
```
Next steps:
  Drop files into <path>/raw/inbox/
  do work bkb triage            Sort inbox items
  do work bkb status            Check KB state
```

**After bkb triage:**
```
Next steps:
  do work bkb ingest            Compile ready sources into wiki
  do work bkb ingest <file>     Ingest a specific source
```

**After bkb ingest:**
```
Next steps:
  do work bkb query [question]  Ask the wiki a question
  do work bkb lint              Health check after ingestion
  do work bkb close             Finalize the day
```

**After bkb query:**
```
Next steps:
  do work bkb ingest            Add more sources
  do work bkb lint              Health check
```

**After bkb lint:**
```
Next steps:
  do work bkb resolve           Resolve flagged contradictions
  do work bkb defrag            Optimize structure (weekly)
  do work bkb garden            Audit relationships and clusters
  do work bkb ingest            Address gaps with new sources
  do work bkb close             Finalize the day
```

**After bkb resolve:**
```
Next steps:
  do work bkb lint              Verify fixes
  do work bkb close             Finalize the day
```

**After bkb close:**
```
Next steps:
  do work bkb rollup            Generate monthly summary (if end of month)
  do work bkb defrag            Weekly structural maintenance
  do work bkb status            Review KB state
```

**After bkb rollup:**
```
Next steps:
  do work bkb lint full         Full integrity check
  do work bkb status            Review KB state
```

**After bkb defrag:**
```
Next steps:
  do work bkb lint              Verify structural integrity after changes
  do work bkb garden            Check relationship hygiene
  do work bkb close             Finalize the day
```

**After bkb garden:**
```
Next steps:
  do work bkb defrag            Apply structural changes if needed
  do work bkb resolve           Resolve any flagged issues
  do work bkb close             Finalize the day
```

**After bkb crew (any action):**
```
Next steps:
  do work bkb crew              Review the full agent roster
  do work bkb status            Check KB state
```

**After tutorial:**
```
Next steps:
  do work capture request: [describe]  Capture your first request
  do work tutorial [mode]              Try another tutorial mode
  do work help                         Full command reference
```

**After version / recap:**
```
Next steps:
  do work run                 Start processing the queue
  do work capture request: [describe]  Capture new requests
```

**Rules:**
- Only suggest prompts that provide value given the current state (e.g., don't suggest `do work run` if the queue is empty)
- Use the full action name (`verify requests`, not just `verify`; `review work`, not just `review`)
- Keep it to 2-3 suggestions max — don't overwhelm
- Format as a simple list the user can scan and copy
- Always include a reminder at the end: `do work help` to see all available commands
