---
id: UR-160
title: 'Directed second opinions: presets, free-text direction, and a linked iteration tree'
created_at: 2026-10-20T19:59:08Z
requests: [REQ-703, REQ-704]
word_count: 569
---
## Summary

The maintainer wants to steer `ai-report second-opinion` like a director: start a pass with a named preset (for example `bold` = "add motion, add more sections, go further") or free text, build the next pass on any earlier page, and walk every iteration as a tree from the original report. The intent was settled over four question rounds and an approved plan (decisions D1-D10 below). The user asked for two chained REQs.

## Extracted Requests

| REQ | What | Depends on |
| --- | --- | --- |
| REQ-703 | Presets and free-text editor direction: `--preset`, `--editor-direction`, shipped + project preset files, direction precedence, direction shown in filename and on the page | none |
| REQ-704 | Linked iteration tree: base page selection, `v1cNN` tree filenames, tree list on `index.html`, per-page parent header, `#rev-*` change marks | REQ-703 |

## Decisions (D1-D10, from the question rounds)

- D1: trigger is `--preset <name>` plus an optional `--editor-direction '<text>'`; both may be combined (preset text first, then the free text).
- D2: shipped presets live in `skills/do-work-toolbox/docs/second-opinion-presets/` (first set: `bold`, `visual`); a project file `ai-reports/second-opinion-presets/<name>.md` with the same name wins; the result row names the file used.
- D3: the direction overrides every editor prompt rule except the file BOUNDARY (the editor writes only its own page and prefixed assets).
- D4: the base is the original `index.html` by default; any earlier tree page can be the base. The user usually names the HTML file; when it is unclear which page is the base, the agent asks with AskUserQuestion, best guess first plus a free-text entry.
- D5: the editor rewrites the base page and reads the original `index.html` and the repository as fact sources, plus the base page's `#rev-*` change list.
- D6: filenames are `index-v1cNN[cNN...]-<cli>-<model>[-unverified]-<preset>[-<editor-slug>].html`; `index.html` is `v1`; child numbers are two digits; depth is unlimited; a plain sort lists pages in tree order.
- D7: the readable slug is the preset name, then a 2-4 word kebab slug the editor proposes; the launcher cleans and length-limits it.
- D8: each new page starts with a "Changed in this iteration" list whose items link to `#rev-<n>` anchors on badged changed sections; the launcher checks every listed anchor exists.
- D9: one launcher-owned list on the original `index.html` shows the whole tree (one-to-many, no depth limit); each tree page carries a launcher-owned header with its id, preset, direction and parent link.
- D10: the three existing second-opinion pages (two in the REQ-632 bundle, one in the REQ-698 bundle) are ignored by the tree code and left as they are.

## Batch Constraints

- Both REQs rewrite `report_second_opinion.go` and the same prompt and action files, so REQ-704 depends on REQ-703 (user chose two chained REQs).
- A plain run (no preset, no direction, base `index.html`) must keep today's output names and behavior.
- Keep REQ-701's rule: a finished page is never deleted by a later check failure.
- Each REQ is a release (shipped files under `skills/`).

## Full Verbatim Input
> ```
> Message 1: can I get a seccond oppinion with: "add motion, add more sections, go further"?
> 
> Message 2: what do I need to change for that?
> 
> Message 3: before you show me the plan, give me several several options on how to trigger this behaviour, and use the ask tool as many times as needed to help me think about it and answer it
> 
> Question round 1 (answers verbatim):
> - How should you start a second opinion that carries your own direction? -> Named presets (preview: ai-report second-opinion --preset bold latest; skills/do-work-toolbox/docs/second-opinion-presets/ bold.md add motion, more sections, go further; deep-dive.md trace every claim to code, add flows; minimal.md shorter, fewer sections; "second opinion, bold" -> --preset bold)
> - When your direction disagrees with the editor's fixed rules, which wins? -> direction wins all, because I'm the director and I want to steer reports as well while keeping them linked togheter so I can walk through all the iterations that where made
> - Where should the page show that it was made with a custom direction? -> 2 and 3, this will allow for easier iteration and filling  (option 2 = also on the page: the launcher-owned banner shows the direction; option 3 = also in filename)
> 
> Question round 2 (answers verbatim):
> - Besides named presets, should you also be able to type a one-off direction? -> Presets + free text (Recommended)
> - Where do preset files live? -> Shipped + project (Recommended)
> - What does the next iteration edit? -> Choose per run (Recommended). Notes: I will usually specify the HTML file and ask a second opinion with that base, when it's not clear which HTML file is the original use the ask tool with best guess and a freeform entry to clarify what is the base that we are building upon
> - How do you walk through the iterations? -> One list on the original. Notes: the tree format is one to many with no limit on depth.
> 
> Question round 3 (answers verbatim):
> - How should the filename show the direction? -> Depth number + preset. Notes: they should be sortable v1c1 (version1 , children 1,2,3..., and then v1c1 can have other children as well v1c1c2 (so first we had v1, it had a children c1, which had 2 children c2 (so this way we have sortable files, in the correct tree order))) it should also have text after that is not an id but human readable slug about what the variation is about.
> - Which presets should ship first? -> visual, bold
> - When building on a base page, what does the editor read? -> base + original + a #rev-* tag showing what changed where
> 
> Question round 4 (answers verbatim):
> - Is this the numbering you mean? -> Two-digit: v1c01c02 (Recommended)
> - Where does the readable slug after the ID come from? -> Preset, then editor slug (Recommended)
> - How should each iteration mark what it changed (the #rev-* tags)? -> Change list + anchors (Recommended)
> - What happens to the three second-opinion pages that already exist? -> Ignore them
> 
> Plan approved by the user (plan file /Users/t2/.claude/plans/dreamy-yawning-starlight.md, decisions D1-D10).
> 
> Message 4: use do-work capture-request to capture all the intent as well + coordinate the execution
> 
> Capture round (answers verbatim):
> - Is this the right proof target for the new request? -> Yes, use it (Recommended)
> - How should this be queued? -> Two REQs, chained
> ```
