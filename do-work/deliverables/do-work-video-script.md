# Video Script: do-work — A Task Queue for AI Coding

**Duration:** ~3 minutes
**Format:** Screen recording with narration (Remotion / Loom / manual)
**Audience:** Developers and technical leaders evaluating AI-assisted development workflows

---

## Scene 1: The Problem (~20s)

**Visual:** Split screen — left side shows a typical AI coding session (long chat thread, scrolling past pages of context, repeated prompts). Right side shows a messy project with no structure, scattered notes, and a git log full of "wip" and "fix" commits.

**Narration:** "AI coding tools are powerful, but they have a memory problem. Every conversation starts from scratch. You lose context between sessions. Complex requests get half-built. There's no paper trail — no way to know what was tried, what failed, or what was actually delivered. And when the client asks 'what did I get?' — you're scrolling through chat logs."

**Transition:** Clear screen to a clean terminal prompt.

---

## Scene 2: Capture (~30s)

**Visual:** Terminal. Type `do work add dark mode to the settings page, also the search is slow, and fix the header alignment`. Show three REQ files being created instantly, plus a UR folder with the verbatim input.

**Narration:** "do-work separates thinking from doing. Throw ideas at it as they come — one-liners, feature specs, even screenshots. Each input becomes structured request files in a queue. Three features in one sentence? Three separate, trackable tasks."

**Key moments:**
- The three files appearing in the file tree
- Quick peek at one REQ file showing the structured format

**Visual:** Type `do work the auth system needs OAuth, user profiles, session management, and password reset` — show a complex input being split into 4 REQs with a UR folder preserving the full verbatim text.

**Narration:** "Complex specs work the same way. Every word is preserved. Nothing gets lost."

---

## Scene 3: Build (~40s)

**Visual:** Type `do work run`. Show the progress output as it processes each REQ:

```
Processing REQ-018-dark-mode.md...
  Triage: Complex (Route C)
  Planning...     [done]
  Exploring...    [done]
  Implementing... [done]
  Testing...      [done] ✓ 12 tests passing
  Reviewing...    [done] 92% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → abc1234
```

**Narration:** "When you're ready, one command starts the queue. Each request gets triaged by complexity. Simple bug fix? Straight to implementation. Complex feature? It plans first, explores the codebase, builds it, runs the tests, does a code review with requirements tracing, captures lessons learned, and creates a clean git commit. Every step is logged in the request file."

**Key moments:**
- Triage routing — showing how Route A/B/C adjusts the pipeline
- The review score appearing (92%)
- The git commit hash appearing

**Visual:** Open the archived REQ file. Scroll through the sections: Triage → Plan → Exploration → Implementation Summary → Testing → Review → Lessons Learned.

**Narration:** "And when it's done, you have a complete history. Not in a chat log you'll never find again — in a file that lives with your code."

---

## Scene 4: Quality Gates (~25s)

**Visual:** Type `do work verify requests`. Show the verification report with per-REQ scores (Coverage, UX Detail, Intent). Then type `do work review work`. Show the review report with requirements checklist and acceptance testing.

**Narration:** "Two quality gates catch problems before they reach the client. Verify checks that the capture was faithful to the original input — did we drop any requirements? Review checks that the code matches those requirements — did we actually build what was asked? Both are automated. Both are optional. Both save hours of rework."

**Key moments:**
- Verification score breakdown
- Review acceptance result: "Pass — component renders correctly"

---

## Scene 5: Present (~25s)

**Visual:** Type `do work present work`. Show the client brief being generated — What We Built, Architecture diagram, Data Flow, Value Proposition. Then show the video script artifact.

**Narration:** "Here's where it gets interesting. The present action reads the archive — what was requested, what was built, how it works — and generates client-facing deliverables. Architecture diagrams, value propositions, even video scripts. No manual effort. The work sells itself."

**Visual:** Quick scroll through the generated client brief, pausing on the architecture diagram and the "Revenue & Growth" section.

---

## Scene 6: The Big Picture (~20s)

**Visual:** Show the archive folder structure — clean UR folders, each self-contained. Type `do work present all` and show the portfolio summary being generated.

**Narration:** "Over time, the archive becomes your project's memory. Every request, every decision, every lesson learned — searchable, traceable, and portable. Switch AI tools whenever you want — the archive travels with you. And when the client asks 'what have we shipped?' — you have an answer in one command."

**Visual:** The architecture diagram from the client brief, annotated with the data flow arrows.

---

## Scene 7: Call to Action (~10s)

**Visual:** Clean terminal. Show the install command:
```
npx skills add knews2019/skill-do-work
```

**Narration:** "do-work. One install. Seven actions. Every task captured, built, reviewed, and ready to present."

---

**Production notes:**
- Terminal should use a dark theme with good contrast for readability
- File tree should be visible in a sidebar (VS Code or similar) when showing file creation
- REQ file content should be readable — zoom in when showing structured format
- Use real-looking but fictional project names and file paths
- The architecture diagram and data flow from the client brief can be shown as visual overlays during Scene 5
- Test with both short (1-line) and long (multi-feature) inputs to show versatility
