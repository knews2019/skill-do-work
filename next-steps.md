# Suggest Next Steps

After every action completes, suggest the next logical prompts the user might want to run. Use fully qualified action names so the user can copy-paste directly.

**After pipeline (completed — queue fully processed):**
```
Next steps:
  do-work present all         Generate portfolio summary across all completed URs
  do-work commit              Commit any uncommitted changes
  do-work capture request: [describe]  Capture new requests
```

**After pipeline (interrupted — active pipeline still exists):**
```
Next steps:
  do-work pipeline            Resume the active pipeline
  do-work pipeline status     Check pipeline progress
```

**After capture requests:**
```
Next steps:
  do-work verify requests     Check capture quality before building
  do-work run                 Start processing the queue
  do-work note "[next hint]"  Jot a lightweight follow-up thought (not a REQ)
```

**After work (queue processing):**
```
Next steps:
  do-work review work         Review the completed work
  do-work present work        Generate client-facing deliverables
  do-work clarify             Answer any pending questions
  do-work roadmap             Survey what's left in the queue (if any REQs remain)
  do-work bkb triage          Sort promoted lessons into the KB (if any REQ has kb_status: promoted)
```

**After verify requests:**
```
Next steps:
  do-work run                 Start processing the queue
  do-work roadmap             Survey feasibility + TDD posture before picking work up
  do-work capture request: [describe changes]  Capture additional requests
```

**After review work:**
```
Next steps:
  do-work present work        Generate client-facing deliverables
  do-work ui-review [scope]   Validate UI quality (if domain: ui-design)
  do-work run                 Process follow-up REQs (if any were created)
  do-work bkb triage          Sort the promoted lesson into the KB (if kb_status: promoted)
  do-work bkb init            Initialize the KB (if kb_status: pending because no kb/ exists)
```

**After validate-feedback:**
```
Next steps:
  do-work capture request: [paste an accepted finding]  Capture an accepted finding as a request
  do-work run                                           Process the captured fixes
  do-work note "[a discuss item]"                       Park a Discuss item for later
```

**After code-review:**
```
Next steps:
  do-work run                   Process follow-up REQs (if any were created)
  do-work quick-wins [dir]      Scan for additional refactoring opportunities
  do-work capture request: [describe fix]  Capture a finding as a request
```

**After ui-review:**
```
Next steps:
  do-work capture request: [describe fix]  Capture findings as requests
  do-work run                   Process follow-up REQs (if any were created)
  do-work install bowser        Install Playwright CLI + Bowser skill for visual verification (if not installed)
```

**After ai-report:**
```
Next steps:
  do-work slop-check ai-reports/<slug>/index.html  Validate the report against the anti-slop principles before sharing
  do-work present work UR-NNN                Generate the complementary client brief / explainer (if not already done)
  do-work inspect                            Review uncommitted changes (the report folder)
  do-work commit                             Commit the report (and any code already staged)
```

**After present work:**
```
Next steps:
  do-work present all         Generate portfolio summary (if multiple URs completed)
  do-work slop-check          Validate the just-generated brief before sending
  do-work capture request: [describe]  Capture new requests
```

**After slop-check:**
```
Next steps:
  do-work slop-check [other-path]   Check another artifact
  do-work present work               Regenerate the brief with the fixes applied
  do-work capture request: [describe]  Capture a follow-up if a flagged issue needs deeper work
```

**After forensics:**
```
Next steps:
  do-work cleanup               Fix orphaned URs and misplaced files
  do-work run                   Process stuck or pending REQs
  do-work roadmap               Survey what's actionable next (pending feasibility + TDD posture)
  do-work capture request: [describe fix]  Capture a specific finding as a request
```

**After roadmap:**
```
Next steps:
  do-work run REQ-NNN           Pick up the top "Ready" REQ
  do-work clarify               Work through pending-answers REQs (if any flagged Needs Clarification)
  do-work bkb triage            Sort staged lessons (only if any REQ has kb_status: promoted)
  do-work review REQ-NNN        Re-run KB handoff (only if any REQ has kb_status: pending; run do-work bkb init first if no kb/)
  do-work forensics             Investigate further if any pending REQ looked suspicious
```

**After note:**
```
Next steps:
  do-work roadmap             See the note surfaced at the top of the queue survey
  do-work note "[next hint]"  Jot another lightweight next-step note
  do-work capture request: [describe]  Promote it to a real task if it warrants building
```

**After cleanup:**
```
Next steps:
  do-work commit                Commit the cleanup changes (archive moves, frontmatter)
  do-work run                   Process any remaining pending REQs
  do-work recap                 Summary of recently completed work
```

**After stray-check:**
```
Next steps:
  do-work commit                Commit the removals (git rm --cached, deletions)
  do-work cleanup               Tidy do-work's own files (loose REQs, misplaced do-work/)
  do-work inspect               Explain the resulting uncommitted changes before committing
```

**After install (any target):**
```
Next steps:
  do-work ui-review [scope]                    Validate UI quality (now with skill/visual verification)
  do-work install [other target]               Install the companion piece (ui-design ↔ bowser)
  do-work capture request: [describe UI work]  Capture a UI-design request
```

**After prime create:**
```
Next steps:
  do-work code-review prime-{name}   Review the code scope the prime covers
  do-work prime audit                Run a full audit to check the new prime
  do-work run                        Process the queue
```

**After prime audit:**
```
Next steps:
  do-work prime create <path>         Create primes for flagged utilities
  do-work capture request: [fix]      Capture audit findings as requests
  do-work run                         Process the queue
```

**After quick-wins:**
```
Next steps:
  do-work capture request: [describe fix]  Capture a finding as a request
  do-work code-review [scope]   Full code review for the same scope
  do-work run                   Process the queue
```

**After scan-ideas:**
```
Next steps:
  do-work capture request: [paste an idea]  Capture an idea as a request
  do-work scan-ideas [different focus]      Brainstorm a different area
  do-work deep-explore [concept]            Explore an idea in depth
  do-work quick-wins [dir]                  Scan for quick refactoring wins
```

**After deep-explore:**
```
Next steps:
  do-work capture request: [paste a direction]  Capture a direction as a request
  do-work deep-explore continue [session]       Resume or extend the session
  do-work scan-ideas [focus]                    Quick idea scan for a related area
```

**After inspect:**
```
Next steps:
  do-work commit              Commit the ready changes
  do-work capture request: [describe fix]  Capture issues as requests
  do-work run                 Process the queue (if fixes were captured)
```

**After commit:**
```
Next steps:
  do-work inspect             Review remaining uncommitted changes (if any)
  do-work review work         Review the committed changes
  do-work capture request: [describe]  Capture new requests
```

**After clarify questions:**
```
Next steps:
  do-work run                 Process answered questions
  do-work clarify             Continue answering (if skipped any)
```

**After bkb init:**
```
Next steps:
  Drop files into <path>/raw/inbox/
  do-work bkb triage            Sort inbox items
  do-work bkb status            Check KB state
```

**After bkb triage:**
```
Next steps:
  do-work bkb ingest            Compile ready sources into wiki
  do-work bkb ingest <file>     Ingest a specific source
```

**After bkb ingest:**
```
Next steps:
  do-work bkb query [question]  Ask the wiki a question
  do-work bkb lint              Health check after ingestion
  do-work bkb close             Finalize the day
```

**After bkb query:**
```
Next steps:
  do-work bkb ingest            Add more sources
  do-work bkb lint              Health check
```

**After bkb lint:**
```
Next steps:
  do-work bkb resolve           Resolve flagged contradictions
  do-work bkb defrag            Optimize structure (weekly)
  do-work bkb garden            Audit relationships and clusters
  do-work bkb ingest            Address gaps with new sources
  do-work bkb close             Finalize the day
```

**After bkb (maintenance subcommands — `close`, `rollup`, `defrag`, `garden`, `crew`, `resolve`):**
```
Next steps:
  do-work bkb lint              Verify integrity after structural / content changes
  do-work bkb status            Review KB state
  do-work bkb close             Finalize the day (skip if you just ran close)
```

**After dream:**
```
Next steps:
  do-work commit                Commit the consolidated memory diff
  do-work bkb lint              (If memory is a bkb wiki) verify it's still structurally sound
  do-work dream [other-path]    Consolidate another memory directory
```

**After prompts list:**
```
Next steps:
  do-work prompts show <name>           Inspect a prompt before running it
  do-work prompts run <name> [args]     Execute a prompt
```

**After prompts show:**
```
Next steps:
  do-work prompts run <name> [args]     Run the prompt you just inspected
  do-work prompts list                  Browse other prompts
```

**After prompts run:**
```
Next steps:
  do-work inspect             Review any uncommitted changes the prompt produced
  do-work prompts list        Browse other prompts in the library
  do-work capture request: [describe]  Capture follow-up work as a request
```

**After interview (session in progress):**
```
Next steps:
  do-work interview <template>         Resume the session (next checkpoint)
  do-work interview <template> review  Cross-layer contradiction pass (once all layers complete)
```

**After interview (all layers complete):**
```
Next steps:
  do-work interview <template> review  Cross-layer contradiction pass before export
  do-work interview <template> export  Produce agent-ready artifacts
```

**After interview export:**
```
Next steps:
  do-work interview <template> ingest  Feed exports into the BKB
  do-work capture request: [describe]  Capture follow-up work the interview surfaced
```

**After interview list:**
```
Next steps:
  do-work interview <template>         Start or resume a listed template
  do-work help                         Full command reference
```

**After tutorial:**
```
Next steps:
  do-work capture request: [describe]  Capture your first request
  do-work tutorial [mode]              Try another tutorial mode
  do-work help                         Full command reference
```

**After version / recap:**
```
Next steps:
  do-work run                 Start processing the queue
  do-work capture request: [describe]  Capture new requests
```

**Rules:**
- Only suggest prompts that provide value given the current state (e.g., don't suggest `do-work run` if the queue is empty)
- Use the full action name (`verify requests`, not just `verify`; `review work`, not just `review`)
- Keep it to 2-3 suggestions max — don't overwhelm
- Format as a simple list the user can scan and copy
- Always include a reminder at the end: `do-work help` to see all available commands
