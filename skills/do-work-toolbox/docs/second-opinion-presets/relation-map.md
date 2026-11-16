Make a relation map: one zoomable inline SVG map whose boxes are the report's own items, such as options, findings, risks, limits, files and requests, and whose lines are the relations the report states between them. A mind map is a tree that grows from one center, while a relation map links any two items, so a reader sees what supports, tests or limits what. Say on the page that the lines are your reading of the report. Every fact still comes from the report or the repository.

1. No library: inline SVG and plain JavaScript only, no Mermaid, no remote file.
2. Columns run left to right in the order of the argument: what is known (evidence, findings, sources), then what is proposed (the subject, the decision, the options), then what tests or limits it (risks, objections, open questions, limits).
3. One map holds about 30 boxes at most; with more, group items into one box per group and say so on the page.
4. The map uses at most four kinds of line, each with one color, a plain name and an entry in a key beside the map.
5. The first view shows about nine main boxes: the subject or decision, the items the report ranks or recommends, and what tests them most directly.
6. With no click, the first view shows the words of the relations between its own boxes.
7. Zoomed out, a box shows its code and title; zoomed in, it also shows its summary.
8. A click on a box opens its detail panel with its sources and relations, shows the boxes related to it, and dims unrelated boxes only lightly, so their titles stay readable.
9. Escape, a second click on the selected box or a click on empty map deselects it and returns the panel to its first content.
10. A "Show everything" button shows every box.
11. Drag moves the map, pinch or Control-scroll zooms it, buttons and keys do the same, and plain scrolling still scrolls the page.
12. A relation word stays readable at every zoom and never covers a box, because shown boxes count as taken space.
13. A word with no free place is left out and counted as "N of M relation words in view", while the panel still lists every relation.
14. A label is never narrower than its longest word, and no word in a label is broken or hyphenated.
15. A click never leaves the clicked box or a related shown box outside the frame, and it zooms in only as far as all of them still fit. Below 600 px wide, readable titles win: after a click the clicked title and every related title in the frame stay at about 12 px or more, and the panel lists each related box with a button that brings that box into view on the map.
16. On keyboard focus, the script resets the map frame's own scroll to 0, then moves the map so the focused box is fully in view.
17. Below 600 px wide, the first view frames the main item, its titles read without zoom or sideways panning, and the page never scrolls sideways.
18. Below the map, every box is also given as text and every relation as one table row, and on a phone each row stacks so it needs no sideways swipe.
19. Without JavaScript, the map shows at full size and scrolls, each relation's words sit on its line as SVG text, and the tool buttons are hidden by a `[hidden]` rule that beats their display rule.
20. When you can, open the page in a browser at desktop and phone width and fix every check that fails.
