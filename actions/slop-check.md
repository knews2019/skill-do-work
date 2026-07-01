# Slop Check Action

> **Part of the do-work skill.** Validates a human-facing artifact against the anti-slop principles before it ships. Read-only by default — flags findings, offers a rewrite, never auto-applies. User-facing walkthrough: [`docs/slop-check-guide.md`](../docs/slop-check-guide.md).

**Source-of-truth principles** live in `crew-members/anti-slop.md`. This action loads that file and runs each principle as an explicit check against the target artifact. The crew-member is *what* to enforce; this action is *how* to inspect a specific draft.

## When to Use

**Use when:**
- About to send a report, brief, summary, or message and want a sanity check first
- Received a long AI-generated draft and want to know what to cut
- A pipeline produced a deliverable that feels bloated — check it before the user reads it
- Pre-handoff validation of any human-facing artifact (client brief, review report, completion summary)

**Do NOT use when:**
- Reviewing code — use `do-work code-review` (or `do-work review-work` for REQ-scoped work)
- Reviewing UI quality — use `do-work ui-review`
- Internal agent status updates — caveman.md handles that
- Commit messages or PR titles — already short by convention; principles apply implicitly

## Input

`$ARGUMENTS` is one of:

1. **File path** — `do-work slop-check do-work/deliverables/UR-003-client-brief.md`. Read the file directly.
2. **REQ or UR reference** — `do-work slop-check REQ-042` or `do-work slop-check UR-003`. Resolve to the most relevant artifact: the deliverable if one exists, otherwise the REQ's Implementation Summary or review report.
3. **"most recent"** or no argument — find the newest **authored** artifact under `do-work/deliverables/`. Glob `*.md` and `*.single.html` (authored prose). Skip `*.marp.html` (mechanical Marp-CLI export of the `.marp.md` source) and contents of `*-video/` directories (Remotion TSX source, not prose). Among the surviving candidates, pick the newest by mtime. If no authored artifacts exist, ask the user to specify a path or paste the draft.
4. **Pasted text** — if the input is multi-paragraph prose rather than a path/ID, treat it as the artifact directly.

If the target cannot be resolved, ask the user to specify a path or paste the draft.

## Steps

### Step 1: Load the Principles

Read `crew-members/anti-slop.md`. The principles in that file are the checklist — **all of them, however many the file currently carries** (eight as of this writing); do not paraphrase, do not skip any. If the crew file has grown a principle this action's table doesn't list yet, add the row — the crew file is canonical, the table below is illustrative.

### Step 2: Resolve the Artifact

Apply the input resolution from the Input section above. Record:

- **Source** — file path, REQ/UR ID, or "pasted input"
- **Word count** of the artifact (rough — `wc -w` or equivalent)
- **Format** — markdown / plain text / HTML / slide deck source

If the artifact is huge (>5,000 words), confirm with the user before proceeding — the check itself shouldn't produce a longer report than the artifact.

### Step 3: Run Each Principle as a Check

For each principle, produce a PASS or FLAG with one-line evidence:

| # | Principle | Status | Evidence |
|---|-----------|--------|----------|
| 1 | Worth reading | PASS / FLAG | quote or line ref |
| 2 | Verified | PASS / FLAG / N-A | quote or line ref |
| 3 | Compressed | PASS / FLAG | word count + concrete inflation pattern |
| 4 | Conclusion first | PASS / FLAG | quote of opening line(s) |
| 5 | AI honesty | PASS / FLAG / N-A | quote or "no disclosure tag" |
| 6 | Needs to exist | PASS / FLAG | one-sentence argument |
| 7 | Medium matches stakes | PASS / FLAG | judgment call with rationale |
| 8 | Decision first, not self-grade | PASS / FLAG / N-A | quote of what leads — the verdict in words, or the score table that displaced it |

Evidence rules:

- For FLAGs, quote the specific phrase or cite the line number — never editorialize generically.
- N-A is only valid for #2 (no factual claims to verify), #5 (no AI was used at all), and #8 (the artifact surfaces no decision, question, or verdict). Document the N-A reason.
- A PASS without evidence is a skipped check — every row needs something concrete.

### Step 4: Summarize

Produce a top-line verdict + the single most important fix:

```
Verdict: {Slop / Borderline / Clean}
Top fix: {one concrete change, not generic advice}

Flags: {N}/{7}
```

**Verdict thresholds:**
- **Clean** — 0–1 FLAGs, none on principles 1, 3, or 4
- **Borderline** — 2–3 FLAGs, or a single FLAG on principles 1, 3, or 4
- **Slop** — 4+ FLAGs, or FLAGs on both 1 and 3

### Step 5: Offer a Rewrite (Optional)

If the verdict is Borderline or Slop, ask the user:

```
Apply suggested fixes? (compress, lead with conclusion, add disclosure tags)
[y/N]
```

Only rewrite on explicit confirmation. The rewrite must:

- Cut word count by at least 30% (target the inflated patterns: throat-clearing, hedge words, unearned headers, redundant bullets)
- Move the conclusion to the first sentence
- Add an explicit disclosure tag if AI-generated and unverified
- Preserve every factual claim verbatim — compression is not paraphrasing

If the user declines, save the report and exit.

### Step 6: Report

Print the table from Step 3, the verdict block from Step 4, and either the rewrite path (if applied) or the unchanged path. Do not pad with summary prose — the table is the report.

## Output Format

```
# Slop Check: {filename or "pasted input"}

**Source:** {path or description}
**Word count:** {N}
**Format:** {markdown / plain text / HTML / ...}

## Findings

| # | Principle | Status | Evidence |
|---|-----------|--------|----------|
| 1 | Worth reading | PASS | opening line states the decision; no filler intro |
| 2 | Verified | FLAG | "Q3 conversion up 47%" — no source cited |
| 3 | Compressed | FLAG | 1,847 words; "It's worth noting that..." appears 6× |
| 4 | Conclusion first | FLAG | First 4 paragraphs are context-setting; verdict in para 5 |
| 5 | AI honesty | FLAG | No disclosure tag; claims framed as considered analysis |
| 6 | Needs to exist | PASS | Decision document for a real deadline |
| 7 | Medium matches stakes | FLAG | Multi-page memo for a question that needs a yes/no |

## Verdict

**Slop** — 5/7 flags, including #1 and #3.

**Top fix:** Cut to 200 words with the verdict as the first sentence. If the supporting reasoning matters, link to a longer appendix; don't lead with it.

## Rewrite

{Optional — only if user confirmed.}
{Path to rewritten file, or inline rewrite if input was pasted.}
```

## Rules

- **Read-only by default.** Never overwrite the source artifact without explicit user confirmation. If applying a rewrite, write to a sibling path (`{name}.compressed.md`) and let the user decide whether to replace.
- **Cite, don't editorialize.** Every FLAG quotes the artifact or gives a line number. "Tone is bloated" is not evidence; "`it's worth noting` appears 6×" is.
- **Don't run on code.** This action checks prose artifacts. If the input is a source file, refuse and redirect to `do-work code-review`.
- **The report is shorter than the artifact.** If your slop-check report is longer than the thing it's checking, you've just produced more slop. Tighten ruthlessly.
- **Don't penalize length when length is earned.** A multi-REQ pipeline completion report can legitimately be long; a one-question memo cannot. Principle #7 (medium matches stakes) does the work here.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "It's already short, no need to check" | Run all 7 checks anyway — short artifacts can still bury the conclusion or lack disclosure | Length isn't the only slop axis; conclusion-burial and unchecked claims happen at any length |
| "I wrote most of this myself" | Check the AI scaffolding patterns anyway — they leak into prose even when the author thinks they're driving | Hedge words, unearned headers, and bullet-lists-for-everything are the tells |
| "The reader can skim" | Note this as a slop flag — making the reader skim IS the slop tax | Skimming is the reader paying the cost the writer should have paid |
| "FLAGging principle #6 ('needs to exist') feels harsh" | Flag it anyway when warranted — a 2-line answer beats a polished 2-page memo | The fact that documents are free to generate doesn't mean they're free to read |
| "Just rewrite it without showing the findings" | Show the table first; only rewrite on confirmation | The user needs to see *why* the artifact failed, not just a replacement |

## Red Flags

- The slop-check report is longer than the artifact being checked
- Every row is PASS but the artifact is obviously long-winded — checks were rubber-stamped
- Every row is FLAG but the evidence column is empty or generic — findings without proof
- A rewrite was applied without user confirmation — the action is read-only by default
- The rewrite paraphrased factual claims instead of preserving them verbatim — compression bled into editing

## Verification Checklist

- [ ] Loaded `crew-members/anti-slop.md` and applied every principle in the file (no skips except documented N-A for #2, #5, and #8)
- [ ] Every row in the findings table has concrete evidence (quote or line ref), not generic prose
- [ ] Verdict matches the threshold rules in Step 4 — no "Clean" with 3 FLAGs, no "Slop" with 0 FLAGs
- [ ] Top fix is specific (`cut to 200 words, conclusion first`) — not generic (`be more concise`)
- [ ] If a rewrite was applied, it was on explicit user confirmation and preserved factual claims verbatim
- [ ] The report itself is shorter than the artifact it checks
