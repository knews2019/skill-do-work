# Performance Audit Action

> **Part of the do-work skill.** Evidence-based performance diagnosis for a target scope — identifies bottlenecks, quantifies impact, and ranks solutions by effort vs improvement. Read-only audit.

**Read-only** — this action does NOT modify any files. It produces a structured performance report only.

## Philosophy

Performance problems almost never occur where you think. Optimize based on evidence, not intuition. Quick wins first — a 5-minute index change that cuts query time by 80% beats a 3-day rewrite that saves 10ms.

## Input

`$TARGET` comes from `$ARGUMENTS` — a directory path, prime file reference, specific file, or a description of the performance concern. If empty, defaults to the current working directory.

### Targeting Modes

- **Directory/prime target**: `do work perf-audit src/api/` — audit all source files in scope
- **Specific concern**: `do work perf-audit the checkout page is slow` — focused investigation
- **Default (no args)**: Broad scan of cwd for common performance anti-patterns

## Steps

### Step 1: Baseline and Scope

1. Resolve `$TARGET` to concrete files or a specific concern
2. Load `prime-*.md` files for project context (architecture, data flow, known hotspots)
3. **Identify the performance domain**:

| Domain | Indicators | Key Metrics |
|--------|-----------|-------------|
| **API/Backend** | Route handlers, database queries, middleware, background jobs | Response time (p50/p95/p99), throughput, error rate |
| **Frontend** | React/Vue/Angular components, bundle size, rendering | Time to Interactive (TTI), Largest Contentful Paint (LCP), bundle size |
| **Database** | Queries, schema, indexes, migrations | Query execution time, row scans, lock contention |
| **Build/CI** | Webpack/Vite config, test suite, CI pipeline | Build time, test execution time |

4. If the user specified a concern, focus on that domain. Otherwise, scan broadly.

### Step 2: Anti-Pattern Scan

Scan source files for known performance anti-patterns. Check what's relevant to the detected stack:

#### Backend / API

| Pattern | What to look for |
|---------|-----------------|
| **N+1 queries** | Loops that execute a query per iteration. ORM `.find()` inside a `.map()` or `for` loop |
| **Missing indexes** | Queries filtering/sorting on columns without indexes. `WHERE` clauses on unindexed fields |
| **Unbounded queries** | `SELECT *` without `LIMIT`, or queries that fetch entire tables into memory |
| **Synchronous blocking** | Blocking I/O in async contexts. `fs.readFileSync` in a request handler. CPU-heavy work on the event loop |
| **Missing caching** | Repeated expensive computations or DB queries for data that changes infrequently |
| **Overfetching** | API endpoints returning full objects when clients need 2-3 fields. No field selection or projection |
| **Sequential I/O** | Multiple independent async operations run sequentially instead of concurrently (`Promise.all`, `asyncio.gather`) |
| **Large payloads** | JSON responses exceeding ~1MB. Missing pagination. Transferring binary data as base64 |

#### Frontend

| Pattern | What to look for |
|---------|-----------------|
| **Unnecessary re-renders** | Components re-rendering on every parent render. Missing `React.memo`, `useMemo`, `useCallback` where justified by profiling |
| **Bundle bloat** | Large dependencies imported for small features. No code splitting. No tree shaking. `moment.js` when `date-fns` suffices |
| **Unoptimized images** | Large images without lazy loading, responsive sizes, or modern formats (WebP/AVIF) |
| **Layout thrashing** | Reading DOM geometry (offsetHeight) then writing (style changes) in a loop |
| **No virtualization** | Rendering 1000+ DOM nodes in a list when only 20 are visible |
| **Render-blocking resources** | Synchronous scripts in `<head>`, large CSS files without critical CSS extraction |

#### Database

| Pattern | What to look for |
|---------|-----------------|
| **Full table scans** | Queries without `WHERE` on indexed columns. `LIKE '%term%'` on large tables |
| **Missing composite indexes** | Multi-column `WHERE` clauses that could benefit from a composite index |
| **Expensive joins** | Joining large tables without appropriate indexes. Cross-join patterns |
| **Lock contention** | Long-running transactions holding locks. `SELECT FOR UPDATE` on hot rows |
| **Schema issues** | Storing JSON blobs that need to be queried. Missing foreign key indexes |

For each finding, record:
- **File** and **line range** (specific)
- **Pattern** (which anti-pattern)
- **Evidence** (what you observed — quote the actual code)
- **Estimated impact** (High / Medium / Low — based on likely execution frequency and data size)

### Step 3: Hotspot Identification

Beyond anti-patterns, identify code that is likely to be performance-critical:

1. **High-traffic paths** — request handlers, middleware chains, frequently called utilities. Check route definitions, middleware registration, and import frequency
2. **Data-heavy operations** — bulk processing, report generation, data exports, migrations
3. **Computation-heavy paths** — sorting large arrays, recursive algorithms, regex on large strings, serialization of complex objects
4. **Startup paths** — module initialization, config loading, connection pooling setup

For each hotspot, assess whether the current implementation is appropriate for its expected load.

### Step 4: Impact Quantification

For each finding, estimate the improvement potential:

| Impact Level | Criteria |
|-------------|----------|
| **High** | 50%+ improvement in the affected metric. Affects user-facing latency or a hot path |
| **Medium** | 20-50% improvement. Noticeable but not dramatic. Or high improvement on a cold path |
| **Low** | <20% improvement. Marginal gain. Or affects a rarely-exercised code path |

Pair with effort estimate:

| Effort | Criteria |
|--------|----------|
| **Trivial** | Add an index, add a `LIMIT`, add `Promise.all` — under 15 minutes |
| **Small** | Refactor a query, add caching, add pagination — 15-60 minutes |
| **Medium** | Restructure data fetching, add connection pooling, implement virtualization — 1-3 hours |
| **Large** | Architectural change, major refactor, new caching layer — 3+ hours |

### Step 5: Solution Ranking

Rank findings by **impact-to-effort ratio**:

1. High impact + Trivial effort = **do this first**
2. High impact + Small effort = **do this second**
3. Medium impact + Trivial effort = **quick win**
4. High impact + Medium effort = **worthwhile investment**
5. Low impact + any effort = **deprioritize**
6. Any impact + Large effort = **capture as a separate REQ for planning**

## Output Format

```markdown
# Performance Audit Report

**Scope**: {description}
**Files analyzed**: {N} files
**Performance domain**: {Backend / Frontend / Database / Build / Mixed}
**Date**: {today}

## Summary

{2-3 sentences — most impactful finding and overall performance health. Lead with the biggest win.}

## Findings

| # | File | Pattern | Evidence | Impact | Effort | Priority |
|---|------|---------|----------|--------|--------|----------|
| 1 | `src/api/orders.ts:34-52` | N+1 query | `orders.map(o => db.findUser(o.userId))` in loop | High | Small | 1 |
| 2 | `src/api/products.ts:88` | Unbounded query | `SELECT * FROM products` with no LIMIT (50k rows) | High | Trivial | 2 |

## Detailed Analysis

### Finding 1: N+1 Query in Order Listing

**File**: `src/api/orders.ts:34-52`
**Pattern**: N+1 query — fetches user for each order in a loop
**Current code**: {quote the problematic code}
**Impact**: High — this endpoint is called on every page load. With 50 orders, it executes 51 queries
**Recommended fix**: Batch fetch users with `WHERE id IN (...)` or use an ORM eager-loading method
**Effort**: Small — single query refactor
**Validation**: Measure query count before/after. Should drop from N+1 to 2 queries

### Finding 2: ...

## Hotspots

{List identified hotspots that aren't anti-patterns but deserve monitoring or profiling.}

| Hotspot | File | Why it matters | Current state |
|---------|------|---------------|---------------|
| Auth middleware | `src/middleware/auth.ts` | Runs on every request | Acceptable — but verify DB lookup is cached |

## Quick Wins Summary

{Top 3-5 findings sorted by impact/effort ratio — the "do these today" list.}

1. {Finding N}: {one-line fix description} — {expected improvement}
2. ...

## Recommended Next Steps

1. {Highest priority fix}
2. {Second priority}
3. {Third priority}

> To act on these findings:
>   do work capture request: [describe the fix]
>   do work run
```

## Rules

- **Do NOT modify any files.** This action is read-only. Report only.
- **Evidence over intuition.** Every finding must cite specific code with file path and line numbers. "The API feels slow" is not a finding — "`src/api/orders.ts:34` executes a query per loop iteration over an unbounded result set" is.
- **Quick wins first.** Lead with findings that have the best impact-to-effort ratio. The user wants to know what to fix today, not what requires a 3-week refactor.
- **Don't premature-optimize.** If a function is called once during startup and takes 200ms, that's probably fine. Focus on hot paths and user-facing latency.
- **Respect the architecture.** Recommend fixes that work within the existing stack. Don't suggest "switch to Redis" if the project doesn't use Redis. Suggest in-process caching first.
- **Be honest about uncertainty.** Without runtime profiling data, findings are based on static analysis. Say "likely bottleneck" not "confirmed bottleneck" when you're inferring from code structure.
- **Proportional depth.** Small scope gets line-by-line analysis. Large scope gets pattern-focused scanning with hotspot sampling. State your approach.
- **Skip vendored and generated files.** Same exclusions as quick-wins.
