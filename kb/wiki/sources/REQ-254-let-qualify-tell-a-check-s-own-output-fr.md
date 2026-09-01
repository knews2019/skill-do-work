---
title: "Lessons from REQ-254: Let qualify tell a check's own output from leftover instrumentation"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-254-let-qualify-tell-a-check-s-own-output-fr.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-254: Let qualify tell a check's own output from leftover instrumentation

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

`qualify.sh`'s debug-artifact scan FAILs on any added `print(` line, which makes it fire on a **check's own success output** — the reporting a checker is supposed to have. It fired exactly that way on REQ-244's remediation and had to be overridden on the record.

## Solution summary

`qualify.sh`'s debug-artifact scan (Check 4) now splits its tokens by property. Unfinished-work markers (`debugger`, `TODO`, `FIXME` — vocabulary illustrative) FAIL anywhere, unchanged. Output primitives (`print(`, `console.log` — the class, not the fired token) are judged by process-exit ownership: a file that ends its own process (exit idioms, illustrative vocabulary) has a terminal audience, so an added output line is presumptively its own reporting and surfaces as a legible WARN; a file that never ends its process is library code, so the same line FAILs naming the file and reason. The output half walks changed files per path in both serial and range modes. REQ-244's ready-made case now passes with a WARN; genuine library instrumentation still FAILs. Deletion of the heuristic was weighed and declined on the REQ's own GREEN criterion (D-01). Orchestrator applied one seam: `review-work.md`'s diff-hygiene line no longer re-flags a checker's success output.

## What worked

**What worked:** Grounding the condition in the real fired case (REQ-244's actual remediation diff) instead of the REQ's abstract description — that is what disqualified three plausible candidate signals before any code. Safe-degradation as a design property: an unknown language's exit idiom falls back to the old FAIL-and-override behavior, never a silent pass.

**What didn't:** The condition-as-implemented (whole-file grep for exit-idiom text) is weaker than the condition-as-stated (file ends its own process) — the ninth instance-vs-class occurrence, this time in the fix for the gate that hunts that shape. And the pipeline's own paperwork had the same hole one level up: review-generated REQs carry no P-A-U block, so the box audit was silently disarmed for exactly the REQs this session processed, and a false "transcription" claim in a commit message went unnoticed until the review re-armed the audit.

**Worth knowing:** WARN is qualify's judgment channel and FAIL its gate; REQ-254 moved one token class between them, and `work.md:750`'s Because-cell is now conditionally stale (over-compliance direction). A forgotten debug print inside any exit-owning script WARNs — the honest boundary of intent-blindness.

## Back-reference

See `do-work/archive/UR-055/REQ-254-let-qualify-tell-a-checks-own-output-from-instrumentation.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `116eec6`.
