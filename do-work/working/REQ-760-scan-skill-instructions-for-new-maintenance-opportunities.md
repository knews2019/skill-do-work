---
id: REQ-760
title: 'Scan skill instructions for new maintenance opportunities'
status: claimed
dispatch_at: 2026-11-25T23:29:58Z
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
write_set: [ai-reports/2026-11-26_0128_proposal-instruction-maintenance, docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png]
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

## Exploration

The capture and verify actions define the original-clause comparison and repair steps. The work testing and review companions govern the later handoffs. The archive includes small direct-route requests with documented acceptance evidence, including REQ-318's Timeline ordering. The five preceding audit topics are excluded by the request. The toolbox proposal report can carry ranked recommendations, source-line evidence and historical receipts without changing the source instructions.

## Scope

**Files I will touch:**
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/` (new) — one self-contained proposal report, its design note and local assets.
- `docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png` (new) — preview of the rendered report.

**Acceptance criteria:**
- Trace one straightforward request from capture through review and distinguish recorded failures from hypothetical risks.
- Rank only new, evidenced recommendations, or report that none qualify. Give current citations, exact cuts and consequences, preserving the original reason, decisions and success evidence.
- Deliver and render the report. Leave the recommendations for the maintainer to decide.
