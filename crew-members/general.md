# The Compass — General Crew Member

<!-- JIT_CONTEXT: This file is always loaded during implementation (Step 6) regardless of domain. It provides foundational conventions (PRIME file discipline, code hygiene, commit etiquette) that apply to every REQ. No domain tag gates it. -->

## PRIME Files Philosophy

When asked to read or create a Prime file (`prime-*.md`), adhere strictly to these rules:

- **Purpose:** Prime files are semantic indexes for a specific utility or folder. They prevent the AI from having to reinvestigate the entire architecture from scratch.
- **Low Noise, High Value:** Keep them concise.
- **Pointers, Not Copies:** Do not copy-paste large blocks of code into the prime file. Instead, point to the code (`See src/utils/parser.ts for the core regex loop`) which acts as the source of truth.
- **No Volatile Metrics:** Do NOT include volatile data like test counts, exact line numbers, or pending invoice totals. These go stale immediately and create noise.
- **Multiple Aspects:** It is perfectly valid to have multiple prime files in the same folder if they describe different aspects (e.g., `prime-checkout-speed.md` and `prime-checkout-consolidation.md`).

## Lessons Discipline

When you load a prime file, scan it for a `## Lessons` section and read every linked REQ lesson **before implementing**. These encode prior mistakes and discoveries from this exact area of the codebase. Pay particular attention to "What didn't work" entries — they prevent repeating failed approaches. If a lesson directly contradicts your planned approach, note the conflict in your `[PLAN]` phase and explain why you're proceeding differently (or adjust the plan).

## Test-Writing Posture

- **Pragmatic over ceremonial.** For bug fixes and new features, prefer red-green validation: write or identify a test, run it to confirm it fails, then implement until it passes. For refactors, config changes, documentation, and cleanup, red-green may not apply — targeted regression tests, lint/build validation, or non-regression evidence is sufficient. The goal is proof the change works, not ritual.
- **Honor captured proof first.** If the REQ contains a `## Red-Green Proof` section, its RED prompt/case and GREEN outcome are the primary behavior your tests must prove. Treat it as a valuable artifact; only adapt it when the codebase genuinely requires a nearby equivalent, and document why.
- **Check the prime file for a testing section.** If the prime maps code areas to specific test commands (e.g., "changes to `lib/inpainting.js` → run `npm run test:api`"), follow that mapping — it takes precedence over generic test detection.
- **Always write tests for new functionality and regression tests for bug fixes.** A new feature without a test is undocumented behavior; a bug fix without a regression test will regress.

## Cross-REQ Test-Break Rules

When your changes cause tests from a prior REQ to fail:

1. Determine whether the behavior change is intentional.
2. **If intentional:** update the failing tests to match the new behavior, and document in the Testing section which REQ's tests changed and why. This creates traceability for which request altered which other request's behavior.
3. **If unintentional:** fix your implementation to preserve the existing behavior — the prior REQ's tests are the contract.

Never delete a failing test without doing this analysis. Deleted tests without justification = silently regressed functionality.

## Discovered-Tasks Contract

If you discover unrelated bugs, technical debt, or missing prerequisites during implementation, **do not fix them inline**. Append a `## Discovered Tasks` section to the REQ (a top-level section, not nested inside Implementation Summary) and list each finding as a bullet. The orchestrator (Step 8) classifies and queues them as follow-up REQs.

Inline fixes blow the declared scope, break traceability, and inflate the diff so it can't be safely reviewed.
