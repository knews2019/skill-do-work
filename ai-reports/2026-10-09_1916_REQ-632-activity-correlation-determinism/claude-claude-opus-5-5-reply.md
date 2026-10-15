I filled the staging file: [index-claude-claude-opus-5-5.html](index-claude-claude-opus-5-5.html). It's the only file I wrote; `index.html` and the evidence files are untouched. `index.html` already showed as modified before I started.

I couldn't open the page in a browser because the Playwright tool wasn't permitted here. I checked it statically instead: every in-page anchor, every link into the original's sections, and every evidence-file link resolves. The page has not actually been opened, so rendering and the JavaScript are untested.

**What I changed to make it clearer**
- **It opens with what has changed since.** The original was written to support a keep-or-delete decision. That decision shipped the next day, so a banner and the 3-sentence TL;DR say so before anything else.
- **One clickable diagram of the current code.** It shows three inputs → one matching rule → a per-request merge → the card and drawer. The deleted code is drawn as a dashed box, and clicking any box shows its file path.
- **One card per part**, each covering what it does, inputs, outputs, file and lines, and status. A 7-step walkthrough follows, with deep detail in collapsible sections.
- **The decision has tabs:** the report's option 1, its option 2, what actually shipped, and a side-by-side table. There's also a timeline from 0.305.67 to the REQ-691 follow-up, and charts built from the evidence files (file size, prefix coverage, git log timing).
- **Plain-word glossary and hover tooltips** for terms like "claimed", "owned tip" and "hand-back merge". The page has a sticky table of contents, dark mode, and works without JavaScript.

**Corrections (all caused by changes after the report was written)**
1. **The decision is no longer open.** Commit `4cc2fb20` (0.305.88) deleted the merge-range code without writing stamp code, which was neither of the report's two options. The reason given is that a stamp can't be deterministic unless it owns the commit.
2. **Merges no longer pass credit to the commits they bring in.** A commit is credited only if it touches the request's file or its subject has a `[REQ-NNN]` tag (`activity_correlation.go:126-141`).
3. **The git log format** no longer includes `%P` (the parent hashes). It is now `%H%x00%cI%x00%s` (`:160-161`).
4. **Line counts:** the report predicted 54 lines deleted and 303 kept. The actual deletion was 58 net (357 → 299), and the file is 308 lines today after REQ-691. The test-file prediction of 337 lines kept was exact.
5. **The tagging rule moved.** It now lives in `skills/do-work/actions/fan-out-reference.md:39` and also covers the hand-back merge. That file says no hook enforces it.
6. **Tests:** the report's test pattern now selects 9 tests, not 11, because the two merge-range tests were deleted. I counted these from source and did not run them.
7. **Line numbers moved:** `readWorktreeAgentGitState` is now at `verify.go:1478-1503`; the shared caller is at `generate.go:682`, called from `serve.go:173` and `generate.go:506`.

**Open questions**
- Is the tagging rule actually followed since 0.305.88/89? All 34 untagged commits in the evidence predate it, and I found no later measurement.
- Do consumer repositories follow it? UR-135's consumer repo used untagged `Integrate …` merges. Its other tool's `req-nnnn-*` branch names aren't read either.
- Not re-verified: the timing numbers, the five-git-calls count, "verify uses only three of the shared reads", and the original's claim that it was rendered in a browser.
- Repository files appear as plain paths, not links, because I didn't check whether the remote is reachable.