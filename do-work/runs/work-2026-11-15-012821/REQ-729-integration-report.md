# REQ-729 integration report (ai-report-judge documents its relative --out rule)

- Release: 0.318.3 — ai-report-judge States Where a Relative --out Starts
- Archive path: do-work/archive/REQ-729-judge-relative-out-rule.md (flat; UR-176 stays open)
- Merge: pre cfa45037, merge 572e8b94 (572e8b940f1d98e3d3524e4309cc691f5a16a6e1), clean merge, no seams needed
- Finalization commit: the `[REQ-729] complete` commit that carries this file (hash in the follow-up commit below, if any)
- Gate: maintainer-verify exit 0, 141 s, single run, load 4.90 at start
- Review: 98%, Pass (opus reviewer). F1 "not the current folder" exact only when --repo-root is passed; F2 next ai-report.md paragraph says "a fresh directory", true only without --out; F3 stakeholder-report.md:44 lacks the rule sentence. All → report only.
- Heavy lanes (all executed, exit 0): do-work-cli-integrations 68 s, staged-skills 49 s, updater 73 s, installer 39 s
- Timing events: builder-work, handback-merge, verification-gate (gate), review, verification-gate (heavy drain); folded into ## Timing (26m 41s)
- Lessons: none (builder proposed none); no satellite or lessons-index change
- Discovered task (report only): ai-report-guide.md was checked; it never mentions --out, so nothing to add
- Left dirty in main checkout, not mine: do-work/working/REQ-724-move-global-standing-rules-into-skill.md (sibling REQ, dirty at start) and untracked sibling run files (REQ-72x handbacks, progress logs, REQ-724 brief/probes)
