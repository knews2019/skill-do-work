# Approach Directives

<!-- JIT_CONTEXT: This file is loaded by the work action or pipeline action when dispatching multiple sub-agents for parallel or sequential work on related REQs. It assigns each agent a distinct implementation lens to improve solution diversity. Not loaded for single-REQ processing. -->

> When dispatching sub-agents for parallel or sequential work, assign each agent
> a distinct approach directive. This shapes their implementation lens without
> changing the requirements. The REQ defines *what* to build; the directive
> shapes *how* to think about building it.

## Directive Pool

Assign one directive per sub-agent. Do not repeat directives within the same
wave or parallel batch.

### Implementation Lenses

1. **Correctness-First**: Prioritize exhaustive edge case handling, defensive
   coding, and comprehensive validation. Ask "what could go wrong?" at every step.
2. **Simplicity-First**: Prioritize the most minimal, readable implementation.
   Fewer lines, fewer abstractions, fewer dependencies. Ask "can this be simpler?"
3. **Performance-First**: Prioritize efficient algorithms, minimal allocations,
   and fast paths. Profile-aware choices. Ask "what's the hot path?"
4. **Extensibility-First**: Prioritize clean interfaces, separation of concerns,
   and future-proof structure. Ask "what changes next?"
5. **User-First**: Prioritize UX impact — error messages, loading states,
   accessibility, discoverability. Ask "what does the user experience?"
6. **Resilience-First**: Prioritize graceful degradation, retry logic, fallback
   paths, and observability. Ask "what happens when this fails?"
7. **Test-First**: Prioritize testability — write tests before implementation,
   design for dependency injection, ensure every behavior is verifiable.
8. **Security-First**: Prioritize threat modeling, input sanitization, least
   privilege, and audit trails. Ask "how could this be exploited?"

## Assignment Rules

- **Single agent**: No directive needed (agent uses its own judgment).
- **2 agents on related REQs**: Assign contrasting directives
  (e.g., Simplicity-First + Correctness-First).
- **3+ agents in a wave**: Assign diverse directives — avoid clustering
  similar lenses (e.g., don't pair Performance-First with Resilience-First
  in the same wave if a creativity lens is available).
- **Review phase**: Note which directive was used — reviewers should evaluate
  whether the lens was applied effectively and whether it introduced blind spots.

## Usage in Action Files

When dispatching a sub-agent, include in its context:

```
Your approach directive for this task is: [DIRECTIVE NAME]
This means you should prioritize [one-line summary] in your implementation
decisions. This does not change the requirements — it shapes your lens.
```

The directive is advisory. If following the directive would produce a clearly
worse outcome (e.g., Performance-First leading to unreadable code for a simple
CRUD endpoint), the agent should note the tension and choose pragmatically.
