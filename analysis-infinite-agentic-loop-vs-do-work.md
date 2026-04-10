# Comparative Analysis: infinite-agentic-loop vs. do-work

## Executive Summary

Both projects orchestrate AI agents to do work iteratively, but they occupy fundamentally different positions on the complexity spectrum. **infinite-agentic-loop** is a single-purpose proof-of-concept that parallelizes content generation in waves. **do-work** is a full production workflow system that manages the entire software development lifecycle — from capturing user intent through implementation, review, and archival.

They share one core insight: *AI agents work better when you give them structured loops with clear phases.* How they apply that insight diverges dramatically.

---

## 1. Architecture Comparison

### infinite-agentic-loop

```
.claude/
  commands/infinite.md    ← The entire system (one prompt file)
  commands/prime.md       ← Context priming (list files + read docs)
  settings.json           ← Permissions (Write, Edit, Bash)
specs/
  invent_new_ui_v*.md     ← Specification files (4 versions)
src/                      ← Generated output (35 HTML files)
src_infinite/             ← Infinite-mode output (25 HTML files)
ai_docs/                  ← Reference docs
CLAUDE.md                 ← Project context
```

**Total system:** ~2 files of logic (infinite.md + CLAUDE.md). Everything else is input specs or generated output.

### do-work

```
SKILL.md                  ← Master router (23-priority dispatch table)
actions/                  ← 20+ action files (capture, work, review, pipeline, etc.)
crew-members/             ← 8 domain-specific rule files (JIT-loaded)
hooks/                    ← Shell hooks (pipeline-guard)
next-steps.md             ← Post-action suggestions
do-work/                  ← Runtime queue structure (REQs, URs, archive)
```

**Total system:** ~300+ KB of structured prompts across 30+ files, with a file-based queue, state machine, and crash recovery.

### Verdict

infinite-agentic-loop is a **single command**. do-work is a **workflow engine**.

---

## 2. The "Loop" Pattern

### infinite-agentic-loop — Generation Loop

```
Analyze Spec → Scan Existing Output → Deploy Parallel Agents → Quality Check → Next Wave
                                          ↑                                        |
                                          └────────────────────────────────────────┘
```

- **Loop unit:** One "wave" of 3-5 parallel sub-agents, each generating a unique HTML file
- **Termination:** Fixed count (1, 5, 20) or "infinite" (until context exhaustion)
- **State between iterations:** Directory scan (what files exist already)
- **Progressive sophistication:** Wave 1 = basic, Wave 2 = multi-dimensional, Wave N = revolutionary
- **No persistence:** No state files, no crash recovery, no queue

### do-work — Work Processing Loop

```
Capture → Verify → Triage → Plan → Explore → Implement → Test → Review → Archive → Commit
    ↑                                                                                  |
    └──────────────────────── Next REQ (context wipe) ─────────────────────────────────┘
```

- **Loop unit:** One REQ (request) through the full development lifecycle
- **Termination:** Queue empty, or session context limit
- **State between iterations:** REQ files with YAML frontmatter (status, route, timestamps), CHECKPOINT.md
- **Progressive depth:** Route A (simple) → B (medium) → C (complex) with proportional effort
- **Full persistence:** File-based queue, working directory, crash recovery, session resumption

### Key Difference

infinite-agentic-loop loops over *generation iterations* of the same type of output. do-work loops over *distinct work items* that each go through a multi-phase lifecycle.

---

## 3. Agent Orchestration

### infinite-agentic-loop

| Aspect | Approach |
|--------|----------|
| **Parallelism** | 3-5 sub-agents launched simultaneously per wave |
| **Agent identity** | Anonymous — each gets spec + iteration number |
| **Coordination** | Directory scan prevents duplicates; uniqueness directives prevent similarity |
| **Context sharing** | Each agent gets: spec file + existing files snapshot + iteration assignment |
| **Failure handling** | "Handle failures gracefully" (no specific mechanism) |

### do-work

| Aspect | Approach |
|--------|----------|
| **Parallelism** | Sequential by default; pipeline can run foreground for synchronous steps |
| **Agent identity** | Role-based: Planner, Explorer, Builder, Reviewer (each with JIT-loaded domain rules) |
| **Coordination** | Orchestrator manages all state; agents never touch queue files |
| **Context sharing** | REQ file is the shared artifact; each agent appends its section |
| **Failure handling** | Typed errors (intent/spec/code/environment), remediation retry, completed-with-issues fallback |

### Key Difference

infinite-agentic-loop uses **homogeneous parallel agents** (same task, different inputs). do-work uses **heterogeneous sequential agents** (different roles, same work item).

---

## 4. State Management

### infinite-agentic-loop

- **State:** The filesystem IS the state (count files in output directory)
- **Recovery:** None — restart from scratch
- **History:** Output files only — no record of what was attempted or why

### do-work

- **State:** YAML frontmatter in REQ files (status machine: pending → claimed → completed)
- **Recovery:** `working/` directory inspection, CHECKPOINT.md, pipeline.json
- **History:** Each REQ is a living document: Triage → Plan → Exploration → Implementation Summary → Testing → Review → Lessons Learned → Decisions

### Key Difference

infinite-agentic-loop is **stateless** (directory = state). do-work is **stateful** (explicit state machine with typed transitions and crash recovery).

---

## 5. Quality Assurance

### infinite-agentic-loop

- Quality enforcement is **embedded in the generation prompt** ("each iteration must be genuinely novel")
- No separate review phase
- No testing
- No acceptance criteria beyond spec compliance
- Quality relies on the spec file being detailed enough

### do-work

- **Multi-layered QA:**
  - Verify (capture QA — did we extract requirements correctly?)
  - Qualification (orchestrator self-check — does Implementation Summary match Scope?)
  - Testing (run actual tests)
  - Review (requirements check + code review + acceptance testing + risk assessment)
  - Remediation (one retry on failure, then completed-with-issues)
- Follow-up REQs auto-created for review findings
- RED-GREEN proof validation for behavioral changes

### Key Difference

infinite-agentic-loop trusts the spec and the agent. do-work has **4 separate quality gates** with remediation loops.

---

## 6. Human Interaction Model

### infinite-agentic-loop

- **One touch:** User provides spec file + output dir + count, then walks away
- **No clarification:** Agent interprets spec as-is
- **No feedback loop:** No mechanism to redirect mid-run

### do-work

- **Two attention windows:**
  1. **Capture:** Interactive clarification while user is present (resolve ambiguities early)
  2. **Clarify:** Batch review of builder decisions after work completes
- **Open Questions:** Builder proceeds with best judgment; user reviews later
- **Addenda:** Can add to in-flight work via new URs linked to existing REQs

### Key Difference

infinite-agentic-loop is **fire-and-forget**. do-work has **structured human-in-the-loop** touchpoints.

---

## 7. Scope and Applicability

### infinite-agentic-loop

- **Domain:** Content generation (specifically UI components)
- **Output type:** Static HTML files with embedded CSS/JS
- **Adaptability:** Swap the spec file to generate different content types
- **Scale model:** More agents = more output (horizontal scaling)

### do-work

- **Domain:** Any software development task
- **Output type:** Code changes, commits, documentation, reviews
- **Adaptability:** Domain-specific crew-member rules, prime files for project context
- **Scale model:** Queue depth = work backlog; complexity routing = proportional effort

### Key Difference

infinite-agentic-loop solves **"generate many variations of X"**. do-work solves **"manage a backlog of diverse development work"**.

---

## 8. What Each Could Learn from the Other

### What do-work could learn from infinite-agentic-loop

1. **Aggressive parallelism:** do-work processes REQs sequentially. For independent Route A tasks, wave-based parallel execution could dramatically increase throughput.
2. **Simplicity as a feature:** infinite-agentic-loop's entire system fits in one command file. do-work's 300KB+ prompt corpus is powerful but heavy — there may be opportunities to simplify.
3. **Progressive sophistication across waves:** The "Wave 1 = basic, Wave N = revolutionary" pattern could apply to iterative refinement within do-work's build phase.
4. **Low barrier to entry:** Anyone can fork infinite-agentic-loop and adapt it in minutes. do-work requires understanding the full system before modifying it.

### What infinite-agentic-loop could learn from do-work

1. **State persistence:** Adding a simple state file would enable crash recovery and session resumption.
2. **Quality gates:** A review phase after generation would catch spec violations and quality regressions.
3. **Traceability:** Logging what each agent was asked to do, what it produced, and what was kept/rejected would make the system debuggable.
4. **Structured input processing:** A capture phase that validates the spec and resolves ambiguities before generating would improve output quality.
5. **Domain-specific rules:** JIT-loaded rules files could adapt generation behavior to different content types without changing the core loop.
6. **Human checkpoints:** A mechanism to review partial output before the next wave would prevent wasted compute on wrong directions.

---

## 9. Summary Matrix

| Dimension | infinite-agentic-loop | do-work |
|-----------|----------------------|---------|
| **Complexity** | ~2 files, 1 command | 30+ files, 20+ actions |
| **Loop type** | Generation waves | Work queue processing |
| **Agent model** | Homogeneous parallel | Heterogeneous sequential |
| **State** | Filesystem (stateless) | File-based state machine |
| **Recovery** | None | Full crash recovery |
| **QA** | Spec compliance only | 4 quality gates + remediation |
| **Human interaction** | Fire-and-forget | 2 structured attention windows |
| **Domain** | Content generation | Software development |
| **Parallelism** | 3-5 agents per wave | Sequential (mostly) |
| **Traceability** | Output files only | Full intent-to-commit trail |
| **Platform** | Claude Code only | Any agentic tool |
| **Learning curve** | Minutes | Hours |

---

## 10. Conclusion

These projects represent two valid philosophies of agentic orchestration:

- **infinite-agentic-loop** says: *"Give the agent a clear spec and let it rip in parallel. Simplicity and speed over process."*
- **do-work** says: *"Structure the work into phases with quality gates, state tracking, and human checkpoints. Process and traceability over speed."*

Neither is universally better. infinite-agentic-loop is the right choice for bulk content generation where individual quality matters less than volume and variety. do-work is the right choice for software development where each change must be correct, reviewed, tested, and traceable.

The most interesting future direction would be a hybrid: do-work's structured lifecycle for managing *what* to build, with infinite-agentic-loop's parallel wave pattern for *how* to explore solutions within the build phase.
