# REQ-733 integration report (relation map page rebuilt left to right; no release)

- Release: none (NO_RELEASE=1; only an export-ignored ai-reports bundle changed)
- Archive path: do-work/archive/UR-177/REQ-733-rebuild-relation-map-page.md (UR-177 closed, input.md moved to do-work/archive/UR-177/)
- Merge: pre c2f21c0d. First merge fd210a91; two review-fix re-merges with the same pre: 567903a5, then final b6a1e1a7 (b6a1e1a72f5338729c2e8acaead6f711a14212e6). All merges clean.
- Finalization commit: 2017fb80 (2017fb8052f01b8f78f449a80a88d9125ed60f82), phase cleanup_complete, no blocked paths or reason codes
- Review (opus, three rounds): R1 74% (F1 labels over box titles after "Show everything" and on a phone; F2 Tab focus off-screen) → integrator fix 4be874a1. R2 72% (G1 words broken in narrow labels; G2 zoom floor lost related boxes) → fix 3eb2ab4b. R3 89% Approve. Rendered page in R3: 1440x1000 light/dark, first view 9 boxes at 88% (15.8 px titles); O2 17 boxes at 58% with 13/13 labels; Show everything 28 at 46%, then O2/K1 zoom to 57% with all labels; 0 labels on boxes, 0 overlaps, 0 broken words from 20% to 240%. 390x844: opens on D+O1-O4 at 73%; a click lands at 42% with 0-3 whole labels plus "N of M relation words in view" (D-16/D-18). Open, report only: H1 uneven zoom-in after manual zoom-out, H2 fit formula repeated, F5 no-JS lines have no words, F7 inherited duplicate label text.
- Judge: finalize pass covers builder bytes ed6a3a8a; the fixed page cannot be re-judged per page until REQ-728 lands (D-17). Integrator and R3 browser checks stand in.
- Gate: maintainer-verify exit 0 at each merge: 133 s (load 4.28), 137 s (4.82), 131 s (4.44)
- Heavy lanes (all executed, exit 0) at b6a1e1a7: queue-kanban-javascript 7 s, queue-kanban-browser 86 s, do-work-cli-integrations 72 s, staged-skills 43 s, updater 77 s, installer 39 s. Also all green at fd210a91 and 567903a5.
- Timing events: builder-work, handback-merge, verification-gate x4, review; folded into ## Timing (1h 16m 26s). The final heavy-drain event reads 9m 07s because the CLI has no end time; the drain took 5m 28s.
- Lessons: no satellite bullet (no satellite covers hand-built report pages, D-20); corrected focusin/scroll lesson kept in the REQ's Lessons Learned
- Helper fix: helpers/make-manifests.py now skips release payloads under NO_RELEASE=1 (it crashed before); committed in the finalization
- Discovered (report only): carry label rules + focusin reset + H1/H2 into REQ-727; .playwright-mcp files appear in the main tree on each browser navigation (every session here deleted its own)
- Left dirty in main checkout, not mine: untracked sibling run files (REQ-723/724/725/726/728/730/731/732 progress logs, REQ-724-handback.md)
