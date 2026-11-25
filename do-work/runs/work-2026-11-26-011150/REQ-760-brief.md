# REQ-760 builder brief

Implement the request below in your T3-assigned checkout and branch. Record their absolute path and name in the handback. Read CLAUDE.md, the applicable primes and the implementation crew rules in skills/do-work/actions/work-implementation-reference.md, including required lesson entries. Consult do-work/lessons-index.md and the linked user request. Preserve sibling changes. Only this request is yours. Use the REQ's own validation, capture the baseline before edits, follow red/green where specified, run focused checks, review the complete diff, and commit the implementation with [REQ-760] at the beginning of the subject. Do not publish, merge or modify queue records, shared version files or CHANGELOG.

Write the completed handback to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-760-handback.md in the main checkout. Provide operative branch, SHA, base, file manifest, changed behavior, focused commands and results, baseline/red evidence, decisions and discovered tasks (explicitly say when none), integration seams and any limitations. Append UTC start/end lines to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-760-progress.log for long operations. These two main-tree paths are the only shared writes. No other agent owns your implementation. Full gate and release are serial integration work. No project run-policy file.

---
id: REQ-760
title: 'Scan skill instructions for new maintenance opportunities'
status: claimed
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-25T23:27:30Z
  basis:
  - Route B
  - 3-file write set
  - 2 subsystems involved
  - 5 acceptance criteria
route: B
created_at: 2026-11-25T22:57:02Z
user_request: UR-190
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
batch: skill-instruction-maintenance
related: [REQ-758, REQ-759]
claimed_at: 2026-11-25T23:27:19Z
---
# Scan skill instructions for new maintenance opportunities
## What
Perform a fresh, read-only maintenance review of the do-work skill instructions. Trace one straightforward request through capture, verification, implementation, testing, and review. Present the findings in a short, ranked `do-work-toolbox ai-report` covering improvements beyond the five findings reviewed at `efcd3618`.

## Review brief
“Tell the model what only I or this project can tell it. Let it exercise judgment on the rest.” Look for an instruction an experienced engineer can infer, two instructions that conflict or could both be deleted, or a rule that adds a question, test run, or write-up without changing the outcome. Check what broke before a rule was added and what a second check actually caught.

For each new candidate give the current file and line, the exact deletion or narrowing, what could go wrong, and a recorded failure if one exists. Distinguish an observed incident from a plausible risk. Preserve the reason the user asked, consequential decisions, and evidence the work succeeded. For capture quality, use corrections to already-stated requirements and avoidable interruptions as the signals. Later changes of preference are separate.

## Constraints
The prior five topics are the baseline: passing mentions and repeated capture checks, routine verification and scores, factual-question formatting, the line-count heuristic, and reuse of a focused check at review. The requested third item is a new scan. Report any fresh findings for the maintainer to decide on. Do not apply the proposed cuts or queue them during the scan. If nothing additional earns a recommendation, say so. If command usage matters, people run `just run-kanban` and `just do-work-update` in their own project after installing the skill.

## Full Context
See `do-work/user-requests/UR-190/input.md` for the verbatim capture instruction.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: User request, third item. but also capture in 3 a new scan for new improvements*

## Triage

Route B. A fresh audit needs exploration of historical failures and the current workflow.

## Plan

**Planning not required** — Route B, exploration-guided implementation.

Use the toolbox ai-report proposal format. Trace a real straightforward request and show what each step produced. The previous five topics are exclusions. Treat the two current capture changes as the agreed direction. Validate line citations against the integrated files if they move before finalization. Place findings only in the report. Use canonical render judging for its four viewport/theme combinations.
