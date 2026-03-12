# do-work: A Task Queue for Agentic Coding Tools

## What We Built

do-work is an operating system for AI-assisted software development. It turns the messy, ad-hoc way developers interact with coding agents — one-off prompts, lost context, no paper trail — into a structured pipeline with full traceability. You throw ideas at it in plain language (feature requests, bug reports, meeting notes, screenshots), and it captures, queues, triages, plans, builds, tests, reviews, and archives each one — with a git commit and a paper trail for every completed task.

Think of it as the difference between shouting tasks across a room and running a production sprint. The AI still does the work. do-work makes sure nothing falls through the cracks, every decision is documented, and the client can see exactly what was delivered and why.

## How It Works

### Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                         SKILL.md (Router)                           │
│  Parses user input → routes to the right action file                │
│  10 priority levels, keyword matching, payload preservation         │
└─────────┬────────┬────────┬────────┬────────┬────────┬──────────────┘
          │        │        │        │        │        │
          ▼        ▼        ▼        ▼        ▼        ▼
     ┌────────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌───────┐ ┌───────┐
     │Capture │ │ Work │ │Verify│ │Review│ │Present│ │Cleanup│
     │Requests│ │      │ │  Req │ │ Work │ │ Work  │ │       │
     └───┬────┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬────┘ └──┬────┘
         │         │        │        │        │        │
         ▼         ▼        ▼        ▼        ▼        ▼
┌──────────────────────────────────────────────────────────────────────┐
│                        do-work/ (File System)                       │
│                                                                     │
│  REQ-*.md (queue)    user-requests/UR-NNN/    working/    archive/  │
│  Pending tasks       Verbatim input +         In-flight   Completed │
│  ready to build      assets per request       work        history   │
└──────────────────────────────────────────────────────────────────────┘
```

**Seven actions**, each a standalone markdown file that any AI agent can follow:

| Action | What it does | File |
|--------|-------------|------|
| **Capture Requests** | Turns natural language into structured REQ files | `capture.md` |
| **Work** | Triages, plans, explores, builds, tests, reviews, archives | `work.md` |
| **Verify Requests** | Scores REQs against original input for completeness | `verify-requests.md` |
| **Review Work** | Post-build code review with requirements tracing | `review-work.md` |
| **Present Work** | Generates client-facing briefs, diagrams, video scripts | `present-work.md` |
| **Cleanup** | Consolidates archive, closes completed UR folders | `cleanup.md` |
| **Version** | Version check, update check, changelog display | `version.md` |

### Data Flow

```
User Input                    Queue                           Archive
(natural language)            (structured files)              (completed history)

  "add dark mode    ┌─────────────────────────┐    ┌──────────────────────┐
   and fix the      │                         │    │  archive/UR-001/     │
   search"          │  UR-003/input.md ◄──┐   │    │  ├── input.md        │
       │            │  (verbatim input)   │   │    │  ├── REQ-010.md      │
       │            │                     │   │    │  │   (Triage → Plan  │
       ▼            │  REQ-018-dark.md ───┘   │    │  │    → Explore →    │
   ┌────────┐       │  REQ-019-search.md ─┘   │    │  │    Implement →   │
   │Capture │       │                         │    │  │    Test → Review  │
   │Requests│──────►│  (status: pending)      │    │  │    → Lessons)    │
   └────────┘       └────────────┬────────────┘    │  └── commit: a1b2c3 │
                                 │                 └──────────────────────┘
                                 ▼                            ▲
                          ┌─────────────┐                     │
                          │    Work     │─────────────────────┘
                          │  Pipeline   │  triage → plan → explore
                          │             │  → build → test → review
                          └─────────────┘  → lessons → archive → commit
```

1. **User says something** — could be 3 words or 3 pages of meeting notes
2. **Capture** creates a UR folder (preserving verbatim input) and one or more REQ files (the actionable queue items). Each REQ links back to its UR.
3. **Work** picks up pending REQs one at a time. It triages complexity (Simple/Medium/Complex), routes through the right pipeline (direct build, explore-then-build, or full plan-explore-build), runs tests, conducts a code review, captures lessons learned, archives the completed REQ, and creates a git commit.
4. **Archive** becomes a self-contained history. Each UR folder holds the original input + all completed REQs with their full lifecycle (triage decisions, plans, exploration findings, implementation notes, test results, review scores, lessons learned).

### Key Design Decisions

- **Capture and execute are strictly separated.** The user decides when to build — the system never "helpfully" starts working after capturing a request. This prevents runaway builds and keeps the user in control.
- **Markdown files as the database.** No external dependencies, no servers, no databases. Everything is plain text files in a `do-work/` directory. Any AI agent that can read and write files can run this system.
- **Agent-agnostic by design.** Action files work as standalone prompts. Originally built for Claude Code, but designed so any agentic coding tool (Cursor, Copilot Workspace, Aider, etc.) can follow the instructions. No tool-specific APIs in the action files.
- **Living documents over separate logs.** Each REQ file accumulates its entire history as appended sections rather than maintaining separate log files. One file tells the whole story.

## Why This Works

do-work solves the fundamental problem with AI-assisted development: context loss. Without a system, every conversation with a coding agent starts from scratch — the agent doesn't know what was tried before, what failed, what the user actually asked for, or what decisions were made along the way. do-work creates a persistent, file-based memory that survives between sessions, tracks every decision, and ensures nothing gets lost between "I have an idea" and "it's in production."

## Value Delivered

### Immediate Impact

- **Nothing falls through the cracks.** Every request is captured, queued, and tracked to completion. No more "I mentioned that three conversations ago."
- **Full traceability.** Every completed task has a paper trail: what was requested, what was triaged, what was planned, what was built, what was tested, what was reviewed, what lessons were learned. Timestamps on everything.
- **Quality gates built in.** Verify checks capture quality. Review checks code quality. Acceptance testing checks it actually works. These run automatically — no discipline required.
- **Institutional memory.** Lessons Learned sections capture what worked, what didn't, and gotchas for the next person (or agent) touching the same code.

### Revenue & Growth Opportunities

- **Faster client delivery.** Structured capture means fewer misunderstandings. Triage-based routing means simple tasks don't get over-engineered. Review catches issues before the client sees them.
- **Sellable artifacts.** The present work action turns completed engineering into client-ready deliverables — briefs, architecture diagrams, value propositions, and video scripts — without manual effort.
- **Scalable operations.** The system handles complex, multi-feature requests the same way it handles one-liners. Queue it, triage it, build it. The queue scales; the process stays consistent.
- **Reduced rework.** Open Questions with recommended defaults mean ambiguity gets resolved before code is written, not after. Review creates follow-up REQs that re-enter the queue — issues don't disappear into "we'll fix it later."

### Competitive Advantage

- **Process as infrastructure.** Most teams using AI coding tools have no process — they prompt and pray. do-work provides the missing layer between "AI can write code" and "AI can deliver software."
- **Works with any AI agent.** Not locked into one platform. Switch from Claude Code to Cursor to Aider — the skill files and the archive travel with you.
- **Self-documenting.** The archive is the documentation. No separate effort needed to explain what was built or why.

## Key Files

- `SKILL.md` — the router: parses user input and dispatches to the right action
- `actions/capture.md` — fast-capture system for turning ideas into structured requests
- `actions/work.md` — the build orchestrator: triage, plan, explore, implement, test, review, archive
- `actions/verify-requests.md` — quality gate comparing REQs against original input
- `actions/review-work.md` — post-build code review with requirements tracing and acceptance testing
- `actions/present-work.md` — generates client-facing deliverables from completed work
- `actions/cleanup.md` — archive consolidation
- `actions/version.md` — version management and changelog display
- `CHANGELOG.md` — 20 versions of release history

## What's Next

- **Dashboard view** — a summary command that shows queue status, in-flight work, archive stats, and pending questions at a glance
- **Priority and dependency ordering** — process REQs in dependency order rather than strictly by number
- **Estimation tracking** — compare triage predictions against actual implementation time to improve future estimates
- **Multi-agent parallelism** — process independent REQs concurrently when the environment supports it
- **Template system** — reusable templates for common request patterns (API endpoint, UI component, migration, etc.)
